// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: Apache-2.0, MIT

package network

import (
	"testing"
	"time"

	peer "github.com/libp2p/go-libp2p-core/peer"

	"github.com/stretchr/testify/require"
)

var timeout = time.Second * 5

func TestStartAndStop(t *testing.T) {
	host, err := NewHost(&Config{})
	require.NoError(t, err)

	err = host.Start()
	require.NoError(t, err)

	err = host.Stop()
	require.NoError(t, err)
}

func TestBootstrap(t *testing.T) {
	hostA, err := NewHost(&Config{})
	require.NoError(t, err)

	hostB, err := NewHost(&Config{
		Bootnodes: []peer.AddrInfo{hostA.AddrInfo()},
	})
	require.NoError(t, err)

	peersA := hostA.host.Peerstore().Peers()
	require.GreaterOrEqual(t, len(peersA), 1)

	peersB := hostB.host.Peerstore().Peers()
	require.GreaterOrEqual(t, len(peersB), 1)
}

func TestBootstrap_PubSub(t *testing.T) {
	t.Skip()
	hostA, err := NewHost(&Config{})
	require.NoError(t, err)

	hostB, err := NewHost(&Config{
		Bootnodes: []peer.AddrInfo{hostA.AddrInfo()},
	})
	require.NoError(t, err)

	err = hostA.Start()
	require.NoError(t, err)

	err = hostB.Start()
	require.NoError(t, err)

	defer func() {
		err = hostA.Stop()
		require.NoError(t, err)

		err = hostB.Stop()
		require.NoError(t, err)
	}()

	// TODO: this should probably be marketsID
	peersA := hostA.pubsub.ListPeers("")
	require.GreaterOrEqual(t, len(peersA), 1)

	peersB := hostB.pubsub.ListPeers("")
	require.GreaterOrEqual(t, len(peersB), 1)
}

func TestPubSub(t *testing.T) {
	t.Skip()
	hostA, err := NewHost(&Config{})
	require.NoError(t, err)

	err = hostA.Start()
	require.NoError(t, err)

	hostB, err := NewHost(&Config{
		Bootnodes: []peer.AddrInfo{hostA.AddrInfo()},
	})
	require.NoError(t, err)

	err = hostB.Start()
	require.NoError(t, err)

	defer func() {
		err = hostA.Stop()
		require.NoError(t, err)

		err = hostB.Stop()
		require.NoError(t, err)
	}()

	data := []byte("nootwashere")
	err = hostA.Publish(data)
	require.NoError(t, err)

	msgs := hostB.Messages()
	select {
	case msg := <-msgs:
		require.NoError(t, err)
		require.Equal(t, data, msg.From)
	case <-time.After(timeout):
		t.Fatal("did not receive msg")
	}
}
