// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: Apache-2.0, MIT

package cache

import (
	"github.com/ipfs/go-cid"
)

// MockCache stores up to size elements, randomly evicting when size is reached
type MockCache struct {
	items map[cid.Cid]struct{}
	size  int
}

// NewMockCache returns a MockCache with the given size
func NewMockCache(size int) *MockCache {
	return &MockCache{
		items: make(map[cid.Cid]struct{}),
		size:  size,
	}
}

func (c *MockCache) Put(cid cid.Cid) {
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
