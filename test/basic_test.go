// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: Apache-2.0, MIT

package test

import (
	"context"
	"math/big"
	"os"
	"sort"
	"testing"
	"time"

	"github.com/ChainSafe/fil-secondary-retrieval-markets/cache"
	"github.com/ChainSafe/fil-secondary-retrieval-markets/client"
	"github.com/ChainSafe/fil-secondary-retrieval-markets/network"
	"github.com/ChainSafe/fil-secondary-retrieval-markets/provider"
	"github.com/ChainSafe/fil-secondary-retrieval-markets/shared"
	block "github.com/ipfs/go-block-format"
	"github.com/ipfs/go-cid"
	ds "github.com/ipfs/go-datastore"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	logging "github.com/ipfs/go-log/v2"
	libp2p "github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/peer"

	"github.com/stretchr/testify/require"
)

var testTimeout = time.Second * 30

func newTestNetwork(t *testing.T) (*network.Network, func()) {
	ctx := context.Background()
	h, err := libp2p.New(ctx)
	require.NoError(t, err)

	net, err := network.NewNetwork(h)
	require.NoError(t, err)

	stop := func() {
		require.NoError(t, net.Stop())
		require.NoError(t, h.Close())
	}
	return net, stop
}

func newTestBlockstore() blockstore.Blockstore {
	nds := ds.NewMapDatastore()
	return blockstore.NewBlockstore(nds)
}

type basicTester struct {
	respCh chan *shared.QueryResponse
}

func newBasicTester() *basicTester {
	return &basicTester{
		respCh: make(chan *shared.QueryResponse),
	}
}

func (bt *basicTester) handleResponse(resp shared.QueryResponse) {
	bt.respCh <- &resp
}

func TestMain(m *testing.M) {
	err := logging.SetLogLevel("client", "debug")
	if err != nil {
		panic(err)
	}
	err = logging.SetLogLevel("provider", "debug")
	if err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}

func TestBasic(t *testing.T) {
	pnet, pnetStop := newTestNetwork(t)
	cnet, cnetStop := newTestNetwork(t)
	bs := newTestBlockstore()

	err := pnet.Connect(cnet.AddrInfo())
	require.NoError(t, err)

	p := provider.NewProvider(pnet, bs, cache.NewMockCache(0))
	c := client.NewClient(cnet)

	// add data block to blockstore
	b := block.NewBlock([]byte("noot"))
	testCid := b.Cid()
	err = bs.Put(b)
	require.NoError(t, err)

	// start provider
	err = p.Start()
	require.NoError(t, err)
	defer pnetStop()

	// start client
	err = c.Start()
	require.NoError(t, err)
	defer cnetStop()

	// subscribe to responses
	bt := newBasicTester()
	unsubscribe := c.SubscribeToQueryResponses(bt.handleResponse, testCid)
	defer unsubscribe()

	// submit query
	err = c.SubmitQuery(context.Background(), testCid)
	require.NoError(t, err)

	// assert response was received
	expected := &shared.QueryResponse{
		PayloadCID:              testCid,
		Provider:                pnet.PeerID(),
		Total:                   big.NewInt(0),
		PaymentInterval:         0,
		PaymentIntervalIncrease: 0,
	}

	select {
	case resp := <-bt.respCh:
		require.NotNil(t, resp)
		require.Equal(t, expected, resp)
	case <-time.After(testTimeout):
		t.Fatal("did not receive response")
	}
}

func TestMulti(t *testing.T) {
	numClients := 3
	numProviders := 3
	data := [][]byte{
		[]byte("noot"),
		[]byte("was"),
		[]byte("here"),
	}
	cids := make([]cid.Cid, len(data))

	clients := make([]*client.Client, numClients)
	providers := make([]*provider.Provider, numProviders)
	blockstores := make([]blockstore.Blockstore, numProviders)
	cnets := make([]*network.Network, numClients)
	pnets := make([]*network.Network, numProviders)

	// create and start clients
	for i := 0; i < numClients; i++ {
		net, stop := newTestNetwork(t)
		c := client.NewClient(net)

		err := c.Start()
		require.NoError(t, err)
		defer stop()

		clients[i] = c
		cnets[i] = net
	}

	// create and start providers
	for i := 0; i < numProviders; i++ {
		net, stop := newTestNetwork(t)
		bs := newTestBlockstore()
		p := provider.NewProvider(net, bs, cache.NewMockCache(0))

		err := p.Start()
		require.NoError(t, err)
		defer stop()

		providers[i] = p
		blockstores[i] = bs
		pnets[i] = net
	}

	// connect clients to providers
	for _, cnet := range cnets {
		for _, pnet := range pnets {
			err := pnet.Connect(cnet.AddrInfo())
			require.NoError(t, err)
		}
	}

	// add data to blockstores
	for i, bs := range blockstores {
		// add data block to blockstore
		b := block.NewBlock(data[i])
		cids[i] = b.Cid()
		err := bs.Put(b)
		require.NoError(t, err)
	}

	// each client queries for a different cid
	for i, c := range clients {
		// subscribe to responses
		bt := newBasicTester()
		unsubscribe := c.SubscribeToQueryResponses(bt.handleResponse, cids[i])
		defer unsubscribe()

		// submit query
		err := c.SubmitQuery(context.Background(), cids[i])
		require.NoError(t, err)

		// assert response was received
		expected := &shared.QueryResponse{
			PayloadCID:              cids[i],
			Provider:                pnets[i].PeerID(),
			Total:                   big.NewInt(0),
			PaymentInterval:         0,
			PaymentIntervalIncrease: 0,
		}

		select {
		case resp := <-bt.respCh:
			require.NotNil(t, resp)
			require.Equal(t, expected, resp)
		case <-time.After(testTimeout):
			t.Fatal("did not receive response")
		}
	}
}

func TestMultiProvider(t *testing.T) {
	pnet0, p0Stop := newTestNetwork(t)
	pnet1, p1Stop := newTestNetwork(t)
	cnet, cStop := newTestNetwork(t)
	bs0 := newTestBlockstore()
	bs1 := newTestBlockstore()

	err := pnet0.Connect(cnet.AddrInfo())
	require.NoError(t, err)
	err = pnet1.Connect(cnet.AddrInfo())
	require.NoError(t, err)
	err = pnet1.Connect(pnet0.AddrInfo())
	require.NoError(t, err)

	time.Sleep(time.Second * 30)

	p0 := provider.NewProvider(pnet0, bs0, cache.NewMockCache(0))
	p1 := provider.NewProvider(pnet1, bs1, cache.NewMockCache(0))
	c := client.NewClient(cnet)

	// add data to both providers's blockstores
	b := block.NewBlock([]byte("noot"))
	testCid := b.Cid()
	err = bs0.Put(b)
	require.NoError(t, err)
	err = bs1.Put(b)
	require.NoError(t, err)

	// start providers and client
	err = p0.Start()
	require.NoError(t, err)
	defer p0Stop()

	err = p1.Start()
	require.NoError(t, err)
	defer p1Stop()

	err = c.Start()
	require.NoError(t, err)
	defer cStop()

	// query for CID, should receive responses from both providers
	bt := newBasicTester()
	unsubscribe := c.SubscribeToQueryResponses(bt.handleResponse, testCid)
	defer unsubscribe()

	// submit query
	err = c.SubmitQuery(context.Background(), testCid)
	require.NoError(t, err)

	// assert response was received
	expected := &shared.QueryResponse{
		PayloadCID:              testCid,
		Total:                   big.NewInt(0),
		PaymentInterval:         0,
		PaymentIntervalIncrease: 0,
	}

	receivedFrom := []peer.ID{}

	for i := 0; i < 2; i++ {
		select {
		case resp := <-bt.respCh:
			require.NotNil(t, resp)
			respProvider := resp.Provider
			resp.Provider = ""
			require.Equal(t, expected, resp)
			t.Log("received from", respProvider)
			receivedFrom = append(receivedFrom, respProvider)
		case <-time.After(testTimeout):
			t.Fatal("did not receive response")
		}
	}

	// assert response was received from providers 0 and 1
	expectedResponders := []peer.ID{pnet0.PeerID(), pnet1.PeerID()}
	sort.Slice(expectedResponders, func(i, j int) bool {
		return expectedResponders[i].String() < expectedResponders[j].String()
	})
	sort.Slice(receivedFrom, func(i, j int) bool {
		return receivedFrom[i].String() < receivedFrom[j].String()
	})
	require.Equal(t, expectedResponders, receivedFrom)
}
