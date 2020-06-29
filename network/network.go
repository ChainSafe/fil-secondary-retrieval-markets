package network

import (
	"context"
	"log"

	libp2p "github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	peer "github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

var baseID = "/fil/"
var marketsID = baseID + "markets"

// Network is the p2p level interface requires by the markets module
type Network interface {
	Start() error
	Stop() error
	AddrInfo() peer.AddrInfo
	Connect(peer.AddrInfo) error
	Messages() <-chan *pubsub.Message
	Publish([]byte) error
}

// Host wraps a libp2p host. It contains the current pubsub state.
// Host implements the Network interface
type Host struct {
	ctx          context.Context
	host         host.Host
	pubsub       *pubsub.PubSub
	topic        *pubsub.Topic
	subscription *pubsub.Subscription
	msgs         chan *pubsub.Message
}

// Config contains the configuration options for the host
type Config struct {
	Bootnodes []peer.AddrInfo
	Identity  crypto.PrivKey
}

// NewHost returns a Host
func NewHost(cfg *Config) (*Host, error) {
	if cfg == nil {
		return nil, ErrNoConfig
	}

	ctx := context.Background()

	h, err := libp2p.New(ctx)
	if err != nil {
		return nil, err
	}

	err = bootstrap(ctx, h, cfg.Bootnodes)
	if err != nil {
		return nil, err
	}

	psOpts := []pubsub.Option{
		pubsub.WithDirectPeers(cfg.Bootnodes),
		pubsub.WithFloodPublish(true),
	}

	ps, err := pubsub.NewGossipSub(ctx, h, psOpts...)
	if err != nil {
		return nil, err
	}

	return &Host{
		ctx:    ctx,
		host:   h,
		pubsub: ps,
		msgs:   make(chan *pubsub.Message),
	}, nil
}

// AddrInfo returns the host's AddrInfo
func (h *Host) AddrInfo() peer.AddrInfo {
	maddrs := h.host.Addrs()
	id := h.host.ID()

	return peer.AddrInfo{
		ID:    id,
		Addrs: maddrs,
	}
}

// Start begins pubsub by subscribing to the markets topic
// TODO: hello protocol
func (h *Host) Start() error {
	var err error
	h.topic, err = h.pubsub.Join(marketsID)
	if err != nil {
		return err
	}

	h.subscription, err = h.topic.Subscribe()
	if err != nil {
		return err
	}

	go h.handleMessages()
	return nil
}

// Stop cancels all subscriptions and shuts down the host.
func (h *Host) Stop() error {
	h.subscription.Cancel()
	err := h.topic.Close()
	if err != nil {
		return err
	}

	return h.host.Close()
}

// Connect connects directly to a peer
func (h *Host) Connect(p peer.AddrInfo) error {
	return h.host.Connect(h.ctx, p)
}

// Publish publishes some data
func (h *Host) Publish(data []byte) error {
	return h.topic.Publish(h.ctx, data)
}

// Messages returns the receive-only pubsub message channel
func (h *Host) Messages() <-chan *pubsub.Message {
	return h.msgs
}

// handleMessages puts each message received through the host's subscription into the host's msgs channel
func (h *Host) handleMessages() {
	for {
		msg, err := h.next()
		if err != nil {
			// TODO: logger
			log.Println("could not receive msg:", err)
		}

		h.msgs <- msg
	}
}

// next returns the next message in the subscription
func (h *Host) next() (*pubsub.Message, error) {
	return h.subscription.Next(h.ctx)
}

func bootstrap(ctx context.Context, h host.Host, bns []peer.AddrInfo) error {
	for _, bn := range bns {
		err := h.Connect(ctx, bn)
		if err != nil {
			return err
		}
	}

	return nil
}
