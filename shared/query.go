// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: Apache-2.0, MIT

package shared

import (
	"math/big"

	"github.com/ipfs/go-cid"
)

// Query is submitted by clients and observed by providers
type Query struct {
	PayloadCID cid.Cid  `json:"payloadCID"` // CID of data being requested
	Client     []string `json:"client"`     // List of multiaddrs of the client
}

// QueryResponse is returned from a provider to a client if the provider has the requested data
type QueryResponse struct {
	PayloadCID              cid.Cid  `json:"payloadCID"` // CID of data being requested
	Provider                []string `json:"provider"`   // List of multiaddrs of the provider
	Total                   *big.Int `json:"total"`      // Total cost
	PaymentInterval         uint64   `json:"paymentInterval"`
	PaymentIntervalIncrease uint64   `json:"paymentIntervalIncrease"`
}
