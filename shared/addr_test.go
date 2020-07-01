// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: Apache-2.0, MIT

package shared

import (
	"testing"

	"github.com/libp2p/go-libp2p-core/peer"

	ma "github.com/multiformats/go-multiaddr"
	"github.com/stretchr/testify/require"
)

func mustDecodePeerID(s string) peer.ID {
	id, err := peer.Decode(s)
	if err != nil {
		panic(err)
	}
	return id
}

func TestStringsToAddrInfos(t *testing.T) {
	testStrs := []string{
		"/ip4/178.62.158.247/tcp/4001/ipfs/QmSoLer265NRgSp2LA3dPaeykiS1J6DifTC88f5uVQKNAd",
		"/ip6/2604:a880:1:20::203:d001/tcp/4001/ipfs/QmSoLPppuBtQSGwKDZT2M73ULpjvfd3aZ6ha4oFGL1KrGM",
	}

	expected := []peer.AddrInfo{
		{
			ID: mustDecodePeerID("QmSoLer265NRgSp2LA3dPaeykiS1J6DifTC88f5uVQKNAd"),
			Addrs: []ma.Multiaddr{
				ma.StringCast("/ip4/178.62.158.247/tcp/4001"),
			},
		},
		{
			ID: mustDecodePeerID("QmSoLPppuBtQSGwKDZT2M73ULpjvfd3aZ6ha4oFGL1KrGM"),
			Addrs: []ma.Multiaddr{
				ma.StringCast("/ip6/2604:a880:1:20::203:d001/tcp/4001"),
			},
		},
	}

	res, err := StringsToAddrInfos(testStrs)
	require.NoError(t, err)
	require.Equal(t, expected, res)
}
