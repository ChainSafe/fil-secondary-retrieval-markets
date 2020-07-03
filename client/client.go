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
	"github.com/ipfs/go-cid"
	logging "github.com/ipfs/go-log/v2"
	"github.com/libp2p/go-libp2p-core/network"
)

var log = logging.Logger("client")

type ClientSubscriber func(resp shared.QueryResponse)

type Unsubscribe func()

type Client struct {
	net             Network
	subscribersLock *sync.Mutex
	subscribers     map[cid.Cid][]ClientSubscriber
}

func NewClient(net Network) *Client {
	c := &Client{
		net:             net,
		subscribersLock: &sync.Mutex{},
		subscribers:     make(map[cid.Cid][]ClientSubscriber),
	}

	// Register handler for provider responses
	c.net.RegisterStreamHandler(shared.RetrievalProtocolID, c.HandleProviderStream)

	return c
}

// SubmitQuery encodes a query a submits it to the network to be gossiped
func (c *Client) SubmitQuery(ctx context.Context, cid cid.Cid) error {
	query := shared.Query{
		PayloadCID:  cid,
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
func (c *Client) SubscribeToQueryResponses(subscriber ClientSubscriber, payloadCID cid.Cid) Unsubscribe {
	c.subscribersLock.Lock()
	c.subscribers[payloadCID] = append(c.subscribers[payloadCID], subscriber)
	c.subscribersLock.Unlock()

	return c.unsubscribeAt(subscriber, payloadCID)
}

// unsubscribeAt returns a function that removes an item from a CID's subscribers list by comparing
// their reflect.ValueOf before pulling the item out of the slice.  Does not preserve order.
// Subsequent, repeated calls to the func with the same Subscriber are a no-op.
// Modified from: https://github.com/filecoin-project/go-fil-markets/blob/6ca8089cea5477fd8539e70ca9b34a61ada6dc27/retrievalmarket/impl/provider.go#L139
func (c *Client) unsubscribeAt(sub ClientSubscriber, cid cid.Cid) Unsubscribe {
	return func() {
		c.subscribersLock.Lock()
		defer c.subscribersLock.Unlock()
		curLen := len(c.subscribers[cid])
		// Remove entry from map if last subscriber
		if curLen == 1 {
			delete(c.subscribers, cid)
			return
		}

		for i, el := range c.subscribers[cid] {
			if reflect.ValueOf(sub) == reflect.ValueOf(el) {
				c.subscribers[cid][i] = c.subscribers[cid][curLen-1]
				c.subscribers[cid] = c.subscribers[cid][:curLen-1]
				return
			}
		}
	}
}

// HandleProviderStream reads the first message and calls HandleProviderResponse
// Note: implements the libp2p StreamHandler interface
func (c *Client) HandleProviderStream(s network.Stream) {
	// Read message from stream
	buf := bufio.NewReader(s)
	bz, err := buf.ReadBytes('\n')
	if err != nil {
		log.Error(err)
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

	c.subscribersLock.Lock()
	defer c.subscribersLock.Unlock()
	if sub := c.subscribers[response.PayloadCID]; sub != nil {
		log.Info("Response received for requested CID: ", response.PayloadCID)
		for _, notifyFn := range sub {
			notifyFn(response)
		}
	} else {
		log.Debug("Provider response received for unknown CID: ", response.PayloadCID)
	}
}
