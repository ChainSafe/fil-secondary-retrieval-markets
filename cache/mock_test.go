// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: Apache-2.0, MIT

package cache

import (
	"sort"
	"testing"

	"github.com/ipfs/go-cid"
	"github.com/stretchr/testify/require"
)

func TestMockCache(t *testing.T) {
	c := NewMockCache(2)

	cid0, err := cid.Decode("QmWATWQ7fVPP2EFGu71UkfnqhYXDYH566qy47CnJDgvs8u")
	require.NoError(t, err)
	cid1, err := cid.Decode("QmSnuWmxptJZdLJpKRarxBMS2Ju2oANVrgbr2xWbie9b2D")
	require.NoError(t, err)

	c.Put(cid0)
	c.Put(cid1)

	res := c.Get(2)
	sort.Slice(res, func(i, j int) bool {
		return res[i].String() < res[j].String()
	})
	require.Equal(t, []cid.Cid{cid1, cid0}, res)

	res = c.Get(3)
	sort.Slice(res, func(i, j int) bool {
		return res[i].String() < res[j].String()
	})
	require.Equal(t, []cid.Cid{cid1, cid0, cid.Cid{}}, res)
}
