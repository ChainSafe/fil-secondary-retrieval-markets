// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: Apache-2.0, MIT

package provider

import (
	"context"
	"testing"

	"github.com/libp2p/go-libp2p-core/peer"

	"github.com/stretchr/testify/require"
)

var testMultiAddrStr = "/ip4/1.2.3.4/tcp/5678/p2p/QmYyQSo1c1Ym7orWxLYvCrM2EmxFTANf8wXmmE7DWjhx5N"

type mockHost struct {
	msgs chan []byte
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

func (h *mockHost) Send(context.Context, peer.ID, []byte) error {
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
