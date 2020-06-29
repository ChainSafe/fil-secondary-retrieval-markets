// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: Apache-2.0, MIT

package client

type Client struct {
	host Host
}

func NewClient(host Host) *Client {
	return &Client{host: host}
}
