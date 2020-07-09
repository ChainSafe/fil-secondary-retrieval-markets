// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: Apache-2.0, MIT

package cache

import (
	"time"

	"github.com/ChainSafe/go-lfu"
	"github.com/ipfs/go-cid"
)

// LFUCache is a least frequentry used cache
type LFUCache struct {
	cache          *lfu.Cache
	size           int
	insertTime     map[cid.Cid]time.Time
	lastAccessTime map[cid.Cid]time.Time
}

// NewLFUCache returns a LFUCache with the given size
func NewLFUCache(size int) *LFUCache {
	return &LFUCache{
		cache:          lfu.New(),
		size:           size,
		insertTime:     make(map[cid.Cid]time.Time),
		lastAccessTime: make(map[cid.Cid]time.Time),
	}
}

// Put adds a cid to the cache
func (c *LFUCache) Put(cid cid.Cid) {
	has := c.cache.Has(cid.String())
	if !has && c.cache.Len() == c.size {
		c.cache.Evict(1)
	}

	c.cache.Set(cid.String(), cid)
	c.lastAccessTime[cid] = time.Now()
	if !has {
		c.insertTime[cid] = time.Now()
	}
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
		Frequency:     freq,
		LastAccessed:  c.lastAccessTime[cid],
		InsertionTime: c.insertTime[cid],
	}
	return r
}