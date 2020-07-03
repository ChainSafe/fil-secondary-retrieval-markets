// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: Apache-2.0, MIT

package network

import (
	"context"
	"fmt"

	"github.com/ChainSafe/fil-secondary-retrieval-markets/shared"
	logging "github.com/ipfs/go-log/v2"
	core "github.com/libp2p/go-libp2p-core"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	peer "github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

var log = logging.Logger("network")

// Host wraps a libp2p host. It contains the current pubsub state.
// Host implements the Network interface
type Network struct {
	host         host.Host
	pubsub       *pubsub.PubSub
	topic        *pubsub.Topic
	subscription *pubsub.Subscription
	msgs         chan []byte
}

// NewNetwork returns a Network
func NewNetwork(h host.Host) (*Network, error) {
	if h == nil {
		return nil, ErrNilHost
	}

	ctx := context.Background()

	psOpts := []pubsub.Option{
		pubsub.WithFloodPublish(true),
	}

	ps, err := pubsub.NewGossipSub(ctx, h, psOpts...)
	if err != nil {
		return nil, err
	}

	return &Network{
		host:   h,
		pubsub: ps,
		msgs:   make(chan []byte),
	}, nil
}

// AddrInfo returns the host's AddrInfo
func (n *Network) AddrInfo() peer.AddrInfo {
	maddrs := n.host.Addrs()
	id := n.host.ID()

	return peer.AddrInfo{
		ID:    id,
		Addrs: maddrs,
	}
}

func (n *Network) MultiAddrs() []string {
	addrs := n.host.Addrs()
	multiaddrs := []string{}

	for _, addr := range addrs {
		multiaddr := fmt.Sprintf("%s/p2p/%s", addr, n.host.ID())
		multiaddrs = append(multiaddrs, multiaddr)
	}

	return multiaddrs
}

func (n *Network) PeerID() peer.ID {
	return n.host.ID()
}

// Start begins pubsub by subscribing to the markets topic
func (n *Network) Start() error {
	var err error
	n.topic, err = n.pubsub.Join(string(shared.RetrievalProtocolID))
	if err != nil {
		return err
	}

	n.subscription, err = n.topic.Subscribe()
	if err != nil {
		return err
	}

	go n.handleMessages()
	return nil
}

// Stop cancels all subscriptions
func (n *Network) Stop() error {
	n.subscription.Cancel()
	return n.topic.Close()
}

// RegisterStreamHandler registers a handler and protocol ID on the libp2p host
func (n *Network) RegisterStreamHandler(id core.ProtocolID, handler network.StreamHandler) {
	n.host.SetStreamHandler(id, handler)
}

// Connect connects directly to a peer
func (n *Network) Connect(p peer.AddrInfo) error {
	ctx := context.Background()
	return n.host.Connect(ctx, p)
}

// Send opens a stream and sends data to the given peer
// TODO: should the protocol ID be passed into this function?
func (n *Network) Send(ctx context.Context, protocol core.ProtocolID, p peer.ID, data []byte) error {
	// TODO: check for existing stream
	s, err := n.host.NewStream(ctx, p, protocol)
	if err != nil {
		return err
	}

	// TODO: add length encoding to msg? or add terminal char?
	_, err = s.Write(data)
	return err
}

// Publish publishes some data
func (n *Network) Publish(ctx context.Context, data []byte) error {
	return n.topic.Publish(ctx, data)
}

// Messages returns the receive-only pubsub message channel
func (n *Network) Messages() <-chan []byte {
	return n.msgs
}

// handleMessages puts each message received through the host's subscription into the host's msgs channel
func (n *Network) handleMessages() {
	for {
		msg, err := n.next()
		if err != nil {
			log.Warn("failed to get next message from subscription")
			continue
		}

		if msg != nil {
			n.msgs <- msg.Data
		}
	}
}

// next returns the next message in the subscription
func (n *Network) next() (*pubsub.Message, error) {
	ctx := context.Background()
	return n.subscription.Next(ctx)
}
