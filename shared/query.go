// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: Apache-2.0, MIT

package shared

import (
	"encoding/json"
	"math/big"

	"github.com/ipfs/go-cid"
	ipld "github.com/ipsn/go-ipfs/gxlibs/github.com/ipfs/go-ipld-format"
	"github.com/libp2p/go-libp2p-core/peer"
)

// Params is the query parameters
type Params struct {
	PayloadCID cid.Cid
	PieceCID   *cid.Cid
	Selector   ipld.Node
}

// Marshal returns the JSON marshalled Params
func (p *Params) Marshal() ([]byte, error) {
	return json.Marshal(p)
}

// MustString returns Params as a string
// It panics if it fails to marshal the Params
func (p *Params) MustString() string {
	bz, err := p.Marshal()
	if err != nil {
		panic(err)
	}
	return string(bz)
}

// Query is submitted by clients and observed by providers
type Query struct {
	Params      Params   `json:"params"`      // CID of data being requested
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
	Params Params `json:"params"` // CID of data being requested
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
