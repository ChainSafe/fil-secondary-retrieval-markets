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
	host Host
}

func NewClient(host Host) *Client {
	return &Client{host: host}
}

// SubmitQuery encodes a query a submits it to the network to be gossiped
func (c *Client) SubmitQuery(ctx context.Context, cid cid.Cid) error {
	query := shared.Query{
		PayloadCID:  cid,
		ClientAddrs: c.host.MultiAddrs(),
	}
	bz, err := json.Marshal(query)
	if err != nil {
		return err
	}

	err = c.host.Publish(ctx, bz)
	if err != nil {
		return err
	}

	return nil
}

// HandleProviderStream verifies the protocol ID of the stream, reads the first message, and calls HandleProviderResponse
// Note: implements the libp2p StreamHandler interface
func (c *Client) HandleProviderStream(s network.Stream) {
	// TODO: Is this necessary?
	// Verify expected protocol ID
	if s.Protocol() != shared.RetrievalProtocolID {
		log.Errorf("unexpected protocol ID %s", s.Protocol())
		_ = s.Close()
		return
	}

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

	// TODO: Make use of message
	log.Infof("Response received: %+v", msg)
}
