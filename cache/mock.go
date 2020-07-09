// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: Apache-2.0, MIT

package cache

import (
	"sync"

	"github.com/ipfs/go-cid"
)

// MockCache stores up to size elements, randomly evicting when size is reached
type MockCache struct {
	items   map[cid.Cid]struct{}
	itemsMu sync.Mutex
	size    int
}

// NewMockCache returns a MockCache with the given size
func NewMockCache(size int) *MockCache {
	return &MockCache{
		items: make(map[cid.Cid]struct{}),
		size:  size,
	}
}

func (c *MockCache) Put(cid cid.Cid) {
	c.itemsMu.Lock()
	defer c.itemsMu.Unlock()

	if len(c.items) == c.size {
		// if cache has reached capacity, delete an element
		for k := range c.items {
			delete(c.items, k)
			break
		}
	}

	c.items[cid] = struct{}{}
}

func (c *MockCache) Get(n int) []cid.Cid {
	c.itemsMu.Lock()
	defer c.itemsMu.Unlock()

	cids := make([]cid.Cid, n)
	i := 0

	for k := range c.items {
		cids[i] = k
		i++
		if i == n {
			break
		}
	}

	return cids
}

func (c *MockCache) Keys() []cid.Cid {
	return []cid.Cid{}
}

func (c *MockCache) GetRecord(cid.Cid) *Record {
	return &Record{}
}
