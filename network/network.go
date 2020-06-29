package network

import (
	"context"

	libp2p "github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	peer "github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	ma "github.com/multiformats/go-multiaddr"
)

var baseID = "/fil/"
var marketsID = baseID + "markets"

// Network is the p2p level interface requires by the markets module
// TODO: is this a good name? would something like Service be better?
type Network interface {
	Start() error
	Stop() error
	AddrInfo() peer.AddrInfo

	// Connect connects directly to a peer
	Connect(peer.AddrInfo) error

	// Next returns the next message in the subscription to marketsID
	Next() (*pubsub.Message, error)

	// Publish publishes some data
	Publish(string, []byte) error
}

// Host wraps a libp2p host. It contains the current pubsub state.
// Host implements the Network interface
type Host struct {
	ctx          context.Context
	host         host.Host
	pubsub       *pubsub.PubSub
	topic        *pubsub.Topic
	subscription *pubsub.Subscription
}

type Config struct {
	Bootnodes []string
	Identity  crypto.PrivKey
}

// NewHost returns a Host
// TODO: bootnodes
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

	ps, err := pubsub.NewGossipSub(ctx, h)
	if err != nil {
		return nil, err
	}

	return &Host{
		ctx:    ctx,
		host:   h,
		pubsub: ps,
	}, nil
}

func bootstrap(ctx context.Context, h host.Host, bns []string) error {
	for _, bn := range bns {
		maddr, err := ma.NewMultiaddr(bn)
		if err != nil {
			return err
		}

		info, err := peer.AddrInfoFromP2pAddr(maddr)
		if err != nil {
			return err
		}

		err = h.Connect(ctx, *info)
		if err != nil {
			return err
		}
	}

	return nil
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

func (h *Host) Next() (*pubsub.Message, error) {
	return h.subscription.Next(h.ctx)
}
