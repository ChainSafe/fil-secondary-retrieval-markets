// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: Apache-2.0, MIT

package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBootstrap(t *testing.T) {
	net0, err := NewNetwork("")
	require.NoError(t, err)

	maddrs := net0.MultiAddrs()
	str := ""
	for i, maddr := range maddrs {
		if i < len(maddrs)-1 {
			str = str + maddr + ","
		} else {
			str = str + maddr
		}
	}

	net1, err := NewNetwork(str)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(net1.Peers()), 1)
}
