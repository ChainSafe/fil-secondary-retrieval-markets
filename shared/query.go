// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: Apache-2.0, MIT

package shared

import (
	"encoding/json"
	"math/big"

	"github.com/ipfs/go-cid"
	"github.com/libp2p/go-libp2p-core/peer"
)

// Query is submitted by clients and observed by providers
type Query struct {
	PayloadCID  cid.Cid  `json:"payloadCID"`  // CID of data being requested
	ClientAddrs []string `json:"clientAddrs"` // List of multiaddrs of the client
}

// Marshal returns the JSON marshalled Query
func (q *Query) Marshal() ([]byte, error) {
	return json.Marshal(q)
}

// Unmarshal JSON unmarshals the input into a Query
func (q *Query) Unmarshal(bz []byte) error {
	return json.Unmarshal(bz, q)
}

// QueryResponse is returned from a provider to a client if the provider has the requested data
type QueryResponse struct {
	PayloadCID cid.Cid `json:"payloadCID"` // CID of data being requested
	// TODO: Do we need their FIL address as well?
	Provider                peer.ID  `json:"provider"` // List of multiaddrs of the provider
	Total                   *big.Int `json:"total"`    // Total cost
	PaymentInterval         uint64   `json:"paymentInterval"`
	PaymentIntervalIncrease uint64   `json:"paymentIntervalIncrease"`
}

// Marshal returns the JSON marshalled QueryResponse
func (q *QueryResponse) Marshal() ([]byte, error) {
	return json.Marshal(q)
}

// Unmarshal JSON unmarshals the input into a QueryResponse
func (q *QueryResponse) Unmarshal(bz []byte) error {
	return json.Unmarshal(bz, q)
}
