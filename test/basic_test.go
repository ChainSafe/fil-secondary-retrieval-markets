package test

import (
	"context"
	"testing"

	"github.com/ChainSafe/fil-secondary-retrieval-markets/client"
	"github.com/ChainSafe/fil-secondary-retrieval-markets/network"
	"github.com/ChainSafe/fil-secondary-retrieval-markets/provider"
	//"github.com/ChainSafe/fil-secondary-retrieval-markets/shared"
	block "github.com/ipfs/go-block-format"
	ds "github.com/ipfs/go-datastore"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	libp2p "github.com/libp2p/go-libp2p"

	"github.com/stretchr/testify/require"
)

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

func TestBasic(t *testing.T) {
	pnet := newTestNetwork(t)
	cnet := newTestNetwork(t)
	bs := newTestBlockstore()

	p := provider.NewProvider(pnet, bs)
	c := client.NewClient(cnet)

	// ad data block to blockstore
	b := block.NewBlock([]byte("noot"))
	testCid := b.Cid()
	err := bs.Put(b)
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

	// submit query
	err = c.SubmitQuery(context.Background(), testCid)
	require.NoError(t, err)
}
