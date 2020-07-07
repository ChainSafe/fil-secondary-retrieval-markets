// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: Apache-2.0, MIT

package provider

import (
	"github.com/ipfs/go-cid"
)

type Cache interface {
	// Put adds the cid to the cache or updates it if it already exists
	Put(cid.Cid)
	// Get gets the top n most popular cids
	Get(n int) []cid.Cid
}
