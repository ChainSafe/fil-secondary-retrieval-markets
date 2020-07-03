// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: Apache-2.0, MIT

package provider

import (
	"context"

	"github.com/libp2p/go-libp2p-core/peer"
)

// Network defines the libp2p network interface used by the Provider
type Network interface {
	Start() error
	Stop() error
	Messages() <-chan []byte
	MultiAddrs() []string
	Connect(p peer.AddrInfo) error
	Send(context.Context, peer.ID, []byte) error
	PeerID() peer.ID
}
