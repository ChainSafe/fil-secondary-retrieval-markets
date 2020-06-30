// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: Apache-2.0, MIT

package shared

import (
	"github.com/ipfs/go-cid"
)

// Query is submitted by clients and observed by providers
type Query struct {
	PayloadCID cid.Cid  `json:"payloadCID"` // COD of data being requested
	Client     []string `json:"client"`     // List of multiaddrs of the client
}
