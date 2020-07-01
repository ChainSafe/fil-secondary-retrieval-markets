// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: Apache-2.0, MIT

package client

import (
	"context"

	core "github.com/libp2p/go-libp2p-core"
	"github.com/libp2p/go-libp2p-core/network"
)

// Host defines the libp2p host used by the client
type Host interface {
	// Publish broadcasts a message over pub sub on the default topic
	Publish(ctx context.Context, msg []byte) error
	// Returns all the hosts multiaddrs
	MultiAddrs() []string

	RegisterStreamHandler(id core.ProtocolID, handler network.StreamHandler)
}
