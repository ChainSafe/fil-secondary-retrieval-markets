// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: Apache-2.0, MIT

package network

import (
	"testing"

	peer "github.com/libp2p/go-libp2p-core/peer"

	"github.com/stretchr/testify/require"
)

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
	require.Equal(t, len(peersA), 2)

	peersB := hostB.host.Peerstore().Peers()
	require.Equal(t, len(peersB), 2)
}

func TestPubSubTopics(t *testing.T) {
	host, err := NewHost(&Config{})
	require.NoError(t, err)

	err = host.Start()
	require.NoError(t, err)

	defer func() {
		err = host.Stop()
		require.NoError(t, err)
	}()

	topics := host.pubsub.GetTopics()
	require.Equal(t, 1, len(topics))
	require.Equal(t, marketsID, topics[0])
}
