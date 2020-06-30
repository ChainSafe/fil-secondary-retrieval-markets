// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: Apache-2.0, MIT

package client

import (
	"context"
	"encoding/json"

	"github.com/ChainSafe/fil-secondary-retrieval-markets/types"
)

type Client struct {
	host Host
}

func NewClient(host Host) *Client {
	return &Client{host: host}
}

// SubmitQuery encodes a query a submits it to the network to be gossiped
func (c *Client) SubmitQuery(ctx context.Context, query types.Query) error {
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
