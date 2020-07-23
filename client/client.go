// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: Apache-2.0, MIT

package client

import (
	"bufio"
	"context"
	"reflect"
	"sync"

	"encoding/json"

	"github.com/ChainSafe/fil-secondary-retrieval-markets/shared"
	logging "github.com/ipfs/go-log/v2"
	"github.com/libp2p/go-libp2p-core/network"
)

var log = logging.Logger("client")

type ClientSubscriber func(resp shared.QueryResponse)

type Unsubscribe func()

type Client struct {
	net             Network
	subscribersLock *sync.Mutex
	subscribers     map[string][]ClientSubscriber
}

func NewClient(net Network) *Client {
	c := &Client{
		net:             net,
		subscribersLock: &sync.Mutex{},
		subscribers:     make(map[string][]ClientSubscriber),
	}

	// Register handler for provider responses
	c.net.RegisterStreamHandler(shared.ResponseProtocolID, c.HandleProviderStream)

	return c
}

// Start starts the client's network
func (c *Client) Start() error {
	return c.net.Start()
}

// Stop stops the client's network
func (c *Client) Stop() error {
	return c.net.Stop()
}

// SubmitQuery encodes a query and submits it to the network to be gossiped
func (c *Client) SubmitQuery(ctx context.Context, params shared.Params) error {
	query := shared.Query{
		Params:      params,
		ClientAddrs: c.net.MultiAddrs(),
	}
	bz, err := json.Marshal(query)
	if err != nil {
		return err
	}

	err = c.net.Publish(ctx, bz)
	if err != nil {
		return err
	}

	return nil
}

// SubscribeQueryResponses registers a subscriber as a listener for a specific payload CID.
// It returns an unsubscribe method that can be called to terminate the subscription.
func (c *Client) SubscribeToQueryResponses(subscriber ClientSubscriber, params shared.Params) Unsubscribe {
	c.subscribersLock.Lock()
	str := params.MustString()
	c.subscribers[str] = append(c.subscribers[str], subscriber)
	c.subscribersLock.Unlock()

	return c.unsubscribeAt(subscriber, params)
}

// unsubscribeAt returns a function that removes an item from a CID's subscribers list by comparing
// their reflect.ValueOf before pulling the item out of the slice.  Does not preserve order.
// Subsequent, repeated calls to the func with the same Subscriber are a no-op.
// Modified from: https://github.com/filecoin-project/go-fil-markets/blob/6ca8089cea5477fd8539e70ca9b34a61ada6dc27/retrievalmarket/impl/provider.go#L139
func (c *Client) unsubscribeAt(sub ClientSubscriber, params shared.Params) Unsubscribe {
	return func() {
		str := params.MustString()
		c.subscribersLock.Lock()
		defer c.subscribersLock.Unlock()
		curLen := len(c.subscribers[str])
		// Remove entry from map if last subscriber
		if curLen == 1 {
			delete(c.subscribers, str)
			return
		}

		for i, el := range c.subscribers[str] {
			if reflect.ValueOf(sub) == reflect.ValueOf(el) {
				c.subscribers[str][i] = c.subscribers[str][curLen-1]
				c.subscribers[str] = c.subscribers[str][:curLen-1]
				return
			}
		}
	}
}

// HandleProviderStream reads the first message and calls HandleProviderResponse
// Note: implements the libp2p StreamHandler interface
func (c *Client) HandleProviderStream(s network.Stream) {
	log.Debug("got stream from peer ", s.Conn().RemotePeer())

	// Read message from stream
	buf := bufio.NewReader(s)
	bz, err := buf.ReadBytes('\n')
	if err != nil {
		return
	}

	c.HandleProviderResponse(bz)
}

// HandleProviderResponse is called to handle a QueryResponse from a provider
func (c *Client) HandleProviderResponse(msg []byte) {
	var response shared.QueryResponse
	err := json.Unmarshal(msg, &response)
	if err != nil {
		log.Error(err)
		return
	}

	log.Info("Response received for requested params: ", response.Params)

	c.subscribersLock.Lock()
	defer c.subscribersLock.Unlock()
	str := response.Params.MustString()
	if sub := c.subscribers[str]; sub != nil {
		for _, notifyFn := range sub {
			notifyFn(response)
		}
	} else {
		log.Debug("Provider response received for unknown params: ", response.Params)
	}
}
