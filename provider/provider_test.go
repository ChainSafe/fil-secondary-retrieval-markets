// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: Apache-2.0, MIT

package provider

import (
	"testing"

	"github.com/stretchr/testify/require"
)

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
