package network

import (
	"context"

	libp2p "github.com/libp2p/go-libp2p"
	host "github.com/libp2p/go-libp2p-core/host"
	peer "github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

var baseID = "/fil/"
var marketsID = baseID + "markets"

// Network is the p2p level interface requires by the markets module
// TODO: is this a good name? would something like Service be better?
type Network interface {
	Start() error
	Stop() error

	// Connect connects directly to a peer
	Connect(peer.AddrInfo) error

	// Publish publishes some data to a topic
	Publish(string, []byte) error

	// Subscribe returns a subscription to a topic
	Subscribe(string) (*pubsub.Subscription, error)

	// Unsubscribe cancels a subscription to a topic, if there is one
	Unsubscribe(string)
}

// Host wraps a libp2p host. It contains the current pubsub state.
// Host implements the Network interface
type Host struct {
	ctx           context.Context
	host          host.Host
	pubsub        *pubsub.PubSub
	topics        map[string]*pubsub.Topic
	subscriptions map[string]*pubsub.Subscription
}

// NewHost returns a Host
// TODO: bootnodes
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
		ctx:           ctx,
		host:          h,
		pubsub:        ps,
		topics:        make(map[string]*pubsub.Topic),
		subscriptions: make(map[string]*pubsub.Subscription),
	}, nil
}

// Start begins pubsub by subscribing to the markets topic
// TODO: hello protocol
func (h *Host) Start() error {
	// TODO: determine what protocol we are using
	t, err := h.pubsub.Join(marketsID)
	if err != nil {
		return err
	}

	h.topics[marketsID] = t

	s, err := t.Subscribe()
	if err != nil {
		return err
	}

	h.subscriptions[marketsID] = s

	return nil
}

// Stop cancels all subscriptions and shuts down the host.
func (h *Host) Stop() error {
	for _, s := range h.subscriptions {
		s.Cancel()
	}

	for _, t := range h.topics {
		err := t.Close()
		if err != nil {
			return err
		}
	}

	return h.host.Close()
}

// Connect connects directly to a peer
func (h *Host) Connect(p peer.AddrInfo) error {
	return h.host.Connect(h.ctx, p)
}

// Publish publishes some data to a topic
// TODO: is there a need for sub-topics under the main marketsID?
func (h *Host) Publish(topic string, data []byte) error {
	t, err := h.join(topic)
	if err != nil {
		return err
	}

	return t.Publish(h.ctx, data)
}

// Subscribe returns a subscription to a topic
// TODO: do we want to wrap the pubsub.Subscription type with our own type to allow for a channel?
func (h *Host) Subscribe(topic string) (*pubsub.Subscription, error) {
	t, err := h.join(topic)
	if err != nil {
		return nil, err
	}

	s, err := t.Subscribe()
	if err != nil {
		return nil, err
	}

	h.subscriptions[topic] = s
	return s, nil
}

// Unsubscribe cancels a subscription to a topic, if there is one
func (h *Host) Unsubscribe(topic string) {
	mt := marketsTopic(topic)
	if h.subscriptions[mt] != nil {
		h.subscriptions[mt].Cancel()
	}
}

// join joins a topic, if the host hasn't already joined it
func (h *Host) join(topic string) (*pubsub.Topic, error) {
	if h.topics[topic] == nil {
		var err error
		h.topics[topic], err = h.pubsub.Join(marketsTopic(topic))
		if err != nil {
			return nil, err
		}
	}

	return h.topics[topic], nil
}

func marketsTopic(topic string) string {
	return marketsID + topic
}
