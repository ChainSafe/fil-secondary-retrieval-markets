// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: Apache-2.0, MIT

package provider

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/ChainSafe/fil-secondary-retrieval-markets/shared"
	"github.com/ipfs/go-cid"
	"github.com/libp2p/go-libp2p-core/peer"

	"github.com/stretchr/testify/require"
)

var testMultiAddrStr = "/ip4/1.2.3.4/tcp/5678/p2p/QmYyQSo1c1Ym7orWxLYvCrM2EmxFTANf8wXmmE7DWjhx5N"

type mockHost struct {
	msgs chan []byte
	sent []byte
}

func newMockHost() *mockHost {
	return &mockHost{
		msgs: make(chan []byte),
	}
}

func (h *mockHost) Start() error {
	return nil
}

func (h *mockHost) Stop() error {
	return nil
}

func (h *mockHost) Messages() <-chan []byte {
	return h.msgs
}

func (h *mockHost) MultiAddrs() []string {
	return []string{testMultiAddrStr}
}

func (h *mockHost) Connect(p peer.AddrInfo) error {
	return nil
}

func (h *mockHost) Send(ctx context.Context, id peer.ID, msg []byte) error {
	h.sent = msg
	return nil
}

func TestProvider(t *testing.T) {
	h := newMockHost()
	p := NewProvider(h)
	err := p.Start()
	require.NoError(t, err)

	defer func() {
		err = p.Stop()
		require.NoError(t, err)
	}()

	msg := []byte("bork")
	h.msgs <- msg
}

func TestProvider_Response(t *testing.T) {
	h := newMockHost()
	p := NewProvider(h)
	err := p.Start()
	require.NoError(t, err)

	defer func() {
		err = p.Stop()
		require.NoError(t, err)
	}()

	testCid, err := cid.Decode("bafybeierhgbz4zp2x2u67urqrgfnrnlukciupzenpqpipiz5nwtq7uxpx4")
	require.NoError(t, err)

	query := &shared.Query{
		PayloadCID: testCid,
		Client:     []string{testMultiAddrStr},
	}

	bz, err := query.Marshal()
	require.NoError(t, err)

	h.msgs <- bz

	resp := &shared.QueryResponse{
		PayloadCID:              query.PayloadCID,
		Provider:                h.MultiAddrs(),
		Total:                   big.NewInt(0),
		PaymentInterval:         0,
		PaymentIntervalIncrease: 0,
	}

	expected, err := resp.Marshal()
	require.NoError(t, err)
	time.Sleep(time.Millisecond * 10)
	require.Equal(t, expected, h.sent)
}
