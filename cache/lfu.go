// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: Apache-2.0, MIT

package cache

import (
	"github.com/ChainSafe/go-lfu"
	"github.com/ipfs/go-cid"
)

// LFUCache is a least frequentry used cache
type LFUCache struct {
	cache *lfu.Cache
}

// NewLFUCache returns a LFUCache with the given size
func NewLFUCache(size int) *LFUCache {
	cache := lfu.New()
	cache.UpperBound = size
	cache.LowerBound = size
	return &LFUCache{
		cache: cache,
	}
}

// Put adds a cid to the cache
func (c *LFUCache) Put(cid cid.Cid) {
	c.cache.Set(cid.String(), cid)
}

// Keys returns all the cids in the cache
func (c *LFUCache) Keys() []cid.Cid {
	strs := c.cache.Keys()
	cids := make([]cid.Cid, len(strs))
	var err error

	for i, s := range strs {
		cids[i], err = cid.Decode(s)
		if err != nil {
			continue
		}
	}

	return cids
}

// GetRecord returns the Record for the given cid
func (c *LFUCache) GetRecord(cid cid.Cid) *Record {
	freq := c.cache.GetFrequency(cid.String())
	r := &Record{
		Frequency: freq,
	}
	return r
}
