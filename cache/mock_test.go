// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: Apache-2.0, MIT

package cache

import (
	"sort"
	"testing"

	"github.com/ipfs/go-cid"
	"github.com/stretchr/testify/require"
)

var cid0, _ = cid.Decode("QmWATWQ7fVPP2EFGu71UkfnqhYXDYH566qy47CnJDgvs8u")
var cid1, _ = cid.Decode("QmSnuWmxptJZdLJpKRarxBMS2Ju2oANVrgbr2xWbie9b2D")
var cid2, _ = cid.Decode("QmdmQXB2mzChmMeKY47C43LxUdg1NDJ5MWcKMKxDu7RgQm")

func TestMockCache(t *testing.T) {
	c := NewMockCache(2)
	c.Put(cid0)
	c.Put(cid1)

	res := c.Keys()
	sort.Slice(res, func(i, j int) bool {
		return res[i].String() < res[j].String()
	})
	require.Equal(t, []cid.Cid{cid1, cid0}, res)
}
