// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: Apache-2.0, MIT

package cache

import (
	"sort"
	"testing"

	"github.com/ipfs/go-cid"
	"github.com/stretchr/testify/require"
)

func TestKeys(t *testing.T) {
	c := NewLFUCache(2)
	c.Put(cid0)
	c.Put(cid1)

	res := c.Keys()
	sort.Slice(res, func(i, j int) bool {
		return res[i].String() < res[j].String()
	})
	require.Equal(t, []cid.Cid{cid1, cid0}, res)
}

func TestGetRecord(t *testing.T) {
	c := NewLFUCache(2)
	c.Put(cid0)
	c.Put(cid1)

	res := c.Keys()
	sort.Slice(res, func(i, j int) bool {
		return res[i].String() < res[j].String()
	})
	require.Equal(t, []cid.Cid{cid1, cid0}, res)

	r0 := c.GetRecord(cid0)
	require.Equal(t, &Record{
		Frequency: 1,
	}, r0)

	c.Put(cid0)
	r0 = c.GetRecord(cid0)
	require.Equal(t, &Record{
		Frequency: 2,
	}, r0)
}

func TestEvict(t *testing.T) {
	c := NewLFUCache(2)
	c.Put(cid0)
	c.Put(cid1)

	res := c.Keys()
	sort.Slice(res, func(i, j int) bool {
		return res[i].String() < res[j].String()
	})
	require.Equal(t, []cid.Cid{cid1, cid0}, res)

	c.Put(cid0) // freq 2
	c.Put(cid2) // should evict cid1

	res = c.Keys()
	sort.Slice(res, func(i, j int) bool {
		return res[i].String() < res[j].String()
	})
	require.Equal(t, []cid.Cid{cid0, cid2}, res)

}
