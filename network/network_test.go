// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: Apache-2.0, MIT

package network

import (
	"context"
	"testing"

	"github.com/ChainSafe/fil-secondary-retrieval-markets/shared"
	libp2p "github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"

	"github.com/stretchr/testify/require"
)

func newTestHost(t *testing.T) host.Host {
	ctx := context.Background()
	h, err := libp2p.New(ctx)
	require.NoError(t, err)
	return h
}

func TestStartAndStop(t *testing.T) {
	h := newTestHost(t)
	n, err := NewNetwork(h)
	require.NoError(t, err)

	err = n.Start()
	require.NoError(t, err)

	err = n.Stop()
	require.NoError(t, err)
}

func TestPubSubTopics(t *testing.T) {
	h := newTestHost(t)
	n, err := NewNetwork(h)
	require.NoError(t, err)

	err = n.Start()
	require.NoError(t, err)

	defer func() {
		err = n.Stop()
		require.NoError(t, err)
	}()

	topics := n.pubsub.GetTopics()
	require.Equal(t, 1, len(topics))
	require.Equal(t, string(shared.RetrievalProtocolID), topics[0])
}
