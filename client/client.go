// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: Apache-2.0, MIT

package client

import (
	"bufio"
	"context"

	"encoding/json"

	"github.com/ChainSafe/fil-secondary-retrieval-markets/shared"
	"github.com/ipfs/go-cid"
	logging "github.com/ipfs/go-log/v2"
	"github.com/libp2p/go-libp2p-core/network"
)

var log = logging.Logger("client")

type Client struct {
	host          Host
	activeQueries QuerySubscriptionStore
}

func NewClient(host Host) *Client {
	c := &Client{
		host:          host,
		activeQueries: make(QuerySubscriptionStore),
	}

	// Register handler for provider responses
	c.host.RegisterStreamHandler(shared.RetrievalProtocolID, c.HandleProviderStream)

	return c
}

// SubmitQuery encodes a query a submits it to the network to be gossiped
func (c *Client) SubmitQuery(ctx context.Context, cid cid.Cid) (*QuerySubscription, error) {
	query := shared.Query{
		PayloadCID:  cid,
		ClientAddrs: c.host.MultiAddrs(),
	}
	bz, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}

	// Setup subscription to responses
	c.activeQueries[cid] = &QuerySubscription{ch: make(chan shared.QueryResponse)}

	err = c.host.Publish(ctx, bz)
	if err != nil {
		return nil, err
	}

	return c.activeQueries[cid], nil
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

	if sub := c.activeQueries[response.PayloadCID]; sub != nil {
		sub.ch <- response
		log.Info("Response received for requested CID", response.PayloadCID)
	} else {
		log.Debug("Provider response received for unknown CID: ", response.PayloadCID)
	}
}
