// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: Apache-2.0, MIT

package provider

import (
	"github.com/ipfs/go-cid"
)

// RequestCache si the interface for the provider's cache of requests
type RequestCache interface {
	// Put adds the cid to the cache or updates it if it already exists
	Put(cid.Cid)
	// Get gets the top n cids, as determined by the implementation
	Get(n int) []cid.Cid
}
