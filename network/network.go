package network

import (
	"context"

	libp2p "github.com/libp2p/go-libp2p"
	host "github.com/libp2p/go-libp2p-core/host"
	peer "github.com/libp2p/go-libp2p-core/peer"
	protocol "github.com/libp2p/go-libp2p-core/protocol"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

var baseID protocol.ID = "/fil/"
var marketsID protocol.ID = baseID + "markets"

type Network interface {
	Start() error
	Stop() error

	// Connect connects directly to a peer
	Connect(peer.AddrInfo) error

	// Publish publishes a
	Publish()

	// Subscribe
	Subscribe()
}

type Host struct {
	ctx    context.Context
	host   host.Host
	pubsub *pubsub.PubSub
	topics map[string]*pubsub.Topic
}

func NewHost() (*Host, error) {
	ctx := context.Background()

	h, err := libp2p.New(ctx)
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
		topics: make(map[string]*pubsub.Topic),
	}, nil
}

func (h *Host) Start() error {
	t, err := h.pubsub.Join(string(marketsID))
	if err != nil {
		return err
	}

	h.topics[string(marketsID)] = t

	return nil
}

func (h *Host) Stop() error {
	for _, t := range h.topics {
		err := t.Close()
		if err != nil {
			return err
		}
	}

	return h.host.Close()
}

func (h *Host) Connect(p peer.AddrInfo) error {
	return h.host.Connect(h.ctx, p)
}

func (h *Host) Publish() {

}

func (h *Host) Subscribe() {

}
