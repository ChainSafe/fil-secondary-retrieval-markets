// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: Apache-2.0, MIT

package test

import (
	"context"
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

func newTestNetwork(t *testing.T) *network.Network {
	ctx := context.Background()
	h, err := libp2p.New(ctx)
	require.NoError(t, err)

	net, err := network.NewNetwork(h)
	require.NoError(t, err)

	t.Cleanup(func() {
		require.NoError(t, net.Stop())
		require.NoError(t, h.Close())
	})
	return net
}

func newTestBlockstore() blockstore.Blockstore {
	nds := ds.NewMapDatastore()
	return blockstore.NewBlockstore(nds)
}

type mockRetrievalProviderStore struct {
	bs blockstore.Blockstore
}

func newTestRetrievalProviderStore() *mockRetrievalProviderStore {
	return &mockRetrievalProviderStore{
		bs: newTestBlockstore(),
	}
}

func (s *mockRetrievalProviderStore) Has(params shared.Params) (bool, error) {
	return s.bs.Has(params.PayloadCID)
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

	time.Sleep(time.Second * 10)
	os.Exit(m.Run())
}

func TestBasic(t *testing.T) {
	pnet := newTestNetwork(t)
	cnet := newTestNetwork(t)
	s := newTestRetrievalProviderStore()

	err := pnet.Connect(cnet.AddrInfo())
	require.NoError(t, err)

	p := provider.NewProvider(pnet, s, cache.NewMockCache(0))
	c := client.NewClient(cnet)

	// add data block to blockstore
	b := block.NewBlock([]byte("noot"))
	testCid := b.Cid()
	err = s.bs.Put(b)
	require.NoError(t, err)

	params := shared.Params{PayloadCID: testCid}

	// start provider
	err = p.Start()
	require.NoError(t, err)

	// start client
	err = c.Start()
	require.NoError(t, err)

	// subscribe to responses
	bt := newBasicTester()
	unsubscribe := c.SubscribeToQueryResponses(bt.handleResponse, params)
	defer unsubscribe()

	// submit query
	err = c.SubmitQuery(context.Background(), params)
	require.NoError(t, err)

	// assert response was received
	expected := &shared.QueryResponse{
		Params:                  params,
		Provider:                pnet.PeerID(),
		PricePerByte:            provider.DefaultPricePerByte,
		PaymentInterval:         provider.DefaultPaymentInterval,
		PaymentIntervalIncrease: provider.DefaultPaymentIntervalIncrease,
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
		net := newTestNetwork(t)
		c := client.NewClient(net)

		err := c.Start()
		require.NoError(t, err)

		clients[i] = c
		cnets[i] = net
	}

	// create and start providers
	for i := 0; i < numProviders; i++ {
		net := newTestNetwork(t)
		s := newTestRetrievalProviderStore()
		p := provider.NewProvider(net, s, cache.NewMockCache(0))

		err := p.Start()
		require.NoError(t, err)

		providers[i] = p
		blockstores[i] = s.bs
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
		params := shared.Params{
			PayloadCID: cids[i],
		}

		// subscribe to responses
		bt := newBasicTester()
		unsubscribe := c.SubscribeToQueryResponses(bt.handleResponse, params)
		defer unsubscribe()

		// submit query
		err := c.SubmitQuery(context.Background(), params)
		require.NoError(t, err)

		// assert response was received
		expected := &shared.QueryResponse{
			Params:                  params,
			Provider:                pnets[i].PeerID(),
			PricePerByte:            provider.DefaultPricePerByte,
			PaymentInterval:         provider.DefaultPaymentInterval,
			PaymentIntervalIncrease: provider.DefaultPaymentIntervalIncrease,
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
	if testing.Short() {
		t.Skip("skipping TestMultiProvider")
	}

	pnet0 := newTestNetwork(t)
	pnet1 := newTestNetwork(t)
	cnet := newTestNetwork(t)
	s0 := newTestRetrievalProviderStore()
	s1 := newTestRetrievalProviderStore()

	err := pnet0.Connect(cnet.AddrInfo())
	require.NoError(t, err)
	err = pnet1.Connect(cnet.AddrInfo())
	require.NoError(t, err)
	err = pnet1.Connect(pnet0.AddrInfo())
	require.NoError(t, err)

	require.GreaterOrEqual(t, len(pnet0.Peers()), 2)
	require.GreaterOrEqual(t, len(pnet1.Peers()), 2)
	require.GreaterOrEqual(t, len(cnet.Peers()), 2)

	p0 := provider.NewProvider(pnet0, s0, cache.NewMockCache(0))
	p1 := provider.NewProvider(pnet1, s1, cache.NewMockCache(0))
	c := client.NewClient(cnet)

	// add data to both providers's blockstores
	b := block.NewBlock([]byte("noot"))
	testCid := b.Cid()
	err = s0.bs.Put(b)
	require.NoError(t, err)
	err = s1.bs.Put(b)
	require.NoError(t, err)

	// start providers and client
	err = p0.Start()
	require.NoError(t, err)

	err = p1.Start()
	require.NoError(t, err)

	err = c.Start()
	require.NoError(t, err)

	params := shared.Params{
		PayloadCID: testCid,
	}

	// query for CID, should receive responses from both providers
	bt := newBasicTester()
	unsubscribe := c.SubscribeToQueryResponses(bt.handleResponse, params)
	defer unsubscribe()

	// submit query
	err = c.SubmitQuery(context.Background(), params)
	require.NoError(t, err)

	// assert response was received
	expected := &shared.QueryResponse{
		Params:                  params,
		PricePerByte:            provider.DefaultPricePerByte,
		PaymentInterval:         provider.DefaultPaymentInterval,
		PaymentIntervalIncrease: provider.DefaultPaymentIntervalIncrease,
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
