package test

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/ChainSafe/fil-secondary-retrieval-markets/client"
	"github.com/ChainSafe/fil-secondary-retrieval-markets/network"
	"github.com/ChainSafe/fil-secondary-retrieval-markets/provider"
	"github.com/ChainSafe/fil-secondary-retrieval-markets/shared"
	block "github.com/ipfs/go-block-format"
	ds "github.com/ipfs/go-datastore"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	logging "github.com/ipfs/go-log/v2"
	libp2p "github.com/libp2p/go-libp2p"

	"github.com/stretchr/testify/require"
)

var testTimeout = time.Second * 5

func newTestNetwork(t *testing.T) *network.Network {
	ctx := context.Background()
	h, err := libp2p.New(ctx)
	require.NoError(t, err)

	net, err := network.NewNetwork(h)
	require.NoError(t, err)
	return net
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

func TestBasic(t *testing.T) {
	logging.SetLogLevel("client", "debug")
	logging.SetLogLevel("provider", "debug")

	pnet := newTestNetwork(t)
	cnet := newTestNetwork(t)
	bs := newTestBlockstore()

	err := pnet.Connect(cnet.AddrInfo())
	require.NoError(t, err)

	p := provider.NewProvider(pnet, bs)
	c := client.NewClient(cnet)

	// ad data block to blockstore
	b := block.NewBlock([]byte("noot"))
	testCid := b.Cid()
	err = bs.Put(b)
	require.NoError(t, err)

	// start provider
	err = p.Start()
	require.NoError(t, err)
	defer func() {
		err = p.Stop()
		require.NoError(t, err)
	}()

	// start client
	err = c.Start()
	require.NoError(t, err)
	defer func() {
		err = c.Stop()
		require.NoError(t, err)
	}()

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
