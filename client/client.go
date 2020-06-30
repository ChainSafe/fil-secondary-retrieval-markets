// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: Apache-2.0, MIT

package client

import (
	"context"
	"encoding/json"

	"github.com/ChainSafe/fil-secondary-retrieval-markets/shared"
	"github.com/ipfs/go-cid"
)

type Client struct {
	host Host
}

func NewClient(host Host) *Client {
	return &Client{host: host}
}

// SubmitQuery encodes a query a submits it to the network to be gossiped
func (c *Client) SubmitQuery(ctx context.Context, cid cid.Cid) error {
	query := shared.Query{
		PayloadCID: cid,
		Client:     c.host.MultiAddrs(),
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
