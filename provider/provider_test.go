// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: Apache-2.0, MIT

package provider

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/ChainSafe/fil-secondary-retrieval-markets/cache"
	"github.com/ChainSafe/fil-secondary-retrieval-markets/shared"
	block "github.com/ipfs/go-block-format"
	ds "github.com/ipfs/go-datastore"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	core "github.com/libp2p/go-libp2p-core"
	"github.com/libp2p/go-libp2p-core/peer"

	"github.com/stretchr/testify/require"
)

var testMultiAddrStr = "/ip4/1.2.3.4/tcp/5678/p2p/QmYyQSo1c1Ym7orWxLYvCrM2EmxFTANf8wXmmE7DWjhx5N"
var testCacheSize = 64

type mockNetwork struct {
	msgs chan []byte
	sent []byte
}

func newMockNetwork() *mockNetwork {
	return &mockNetwork{
		msgs: make(chan []byte),
	}
}

func (n *mockNetwork) Start() error {
	return nil
}

func (n *mockNetwork) Stop() error {
	return nil
}

func (n *mockNetwork) Messages() <-chan []byte {
	return n.msgs
}

func (n *mockNetwork) MultiAddrs() []string {
	return []string{testMultiAddrStr}
}

func (n *mockNetwork) Connect(p peer.AddrInfo) error {
	return nil
}

func (n *mockNetwork) Send(ctx context.Context, protocol core.ProtocolID, id peer.ID, msg []byte) error {
	n.sent = msg
	return nil
}

func (n *mockNetwork) PeerID() peer.ID {
	id, err := peer.Decode("QmYyQSo1c1Ym7orWxLYvCrM2EmxFTANf8wXmmE7DWjhx5N")
	if err != nil {
		panic(err)
	}
	return id
}

func newTestBlockstore() blockstore.Blockstore {
	nds := ds.NewMapDatastore()
	return blockstore.NewBlockstore(nds)
}

func TestProvider(t *testing.T) {
	n := newMockNetwork()
	p := NewProvider(n, newTestBlockstore(), cache.NewMockCache(testCacheSize))
	err := p.Start()
	require.NoError(t, err)

	defer func() {
		err = p.Stop()
		require.NoError(t, err)
	}()

	msg := []byte("bork")
	n.msgs <- msg
}

func TestProvider_Response(t *testing.T) {
	n := newMockNetwork()
	p := NewProvider(n, newTestBlockstore(), cache.NewMockCache(testCacheSize))
	err := p.Start()
	require.NoError(t, err)

	defer func() {
		err = p.Stop()
		require.NoError(t, err)
	}()

	b := block.NewBlock([]byte("noot"))
	testCid := b.Cid()

	err = p.blockstore.Put(b)
	require.NoError(t, err)

	query := &shared.Query{
		PayloadCID:  testCid,
		ClientAddrs: []string{testMultiAddrStr},
	}

	bz, err := query.Marshal()
	require.NoError(t, err)

	n.msgs <- bz

	resp := &shared.QueryResponse{
		PayloadCID:              query.PayloadCID,
		Provider:                n.PeerID(),
		Total:                   big.NewInt(0),
		PaymentInterval:         0,
		PaymentIntervalIncrease: 0,
	}

	expected, err := resp.Marshal()
	require.NoError(t, err)
	time.Sleep(time.Millisecond * 10)
	require.Equal(t, expected, n.sent)
}
