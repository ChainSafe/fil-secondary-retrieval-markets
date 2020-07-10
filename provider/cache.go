// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: Apache-2.0, MIT

package provider

import (
	"github.com/ipfs/go-cid"

	"github.com/ChainSafe/fil-secondary-retrieval-markets/cache"
)

// RequestCache is the interface for the provider's cache of requests
type RequestCache interface {
	// Put adds the cid to the cache or updates it if it already exists
	Put(cid.Cid)

	// Keys returns all the keys in the cache
	Keys() []cid.Cid

	// GetRecord returns the Record for the given cid
	GetRecord(cid.Cid) *cache.Record
}
