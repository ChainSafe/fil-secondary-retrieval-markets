// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: Apache-2.0, MIT

package shared

import (
	"math/big"

	"github.com/ipfs/go-cid"
	"github.com/libp2p/go-libp2p-core/peer"
)

// Query is submitted by clients and observed by providers
type Query struct {
	PayloadCID  cid.Cid  `json:"payloadCID"`  // CID of data being requested
	ClientAddrs []string `json:"clientAddrs"` // List of multiaddrs of the client
}

type QueryResponse struct {
	PayloadCID cid.Cid
	// TODO: Do we need their FIL address as well?
	Provider                peer.ID
	Total                   *big.Int
	PaymentInterval         uint64
	PaymentIntervalIncrease uint64
}
