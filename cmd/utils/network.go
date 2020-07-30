// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: Apache-2.0, MIT

package utils

import (
	"context"
	"strings"

	"github.com/ChainSafe/fil-secondary-retrieval-markets/network"
	"github.com/ChainSafe/fil-secondary-retrieval-markets/shared"
	libp2p "github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	peer "github.com/libp2p/go-libp2p-core/peer"
)

func NewNetwork(bootnodesStr string) (*network.Network, error) {
	ctx := context.Background()
	h, err := libp2p.New(ctx)
	if err != nil {
		return nil, err
	}

	n, err := network.NewNetwork(h)
	if err != nil {
		return nil, err
	}

	// bootstrap to network
	if bootnodesStr != "" {
		strs := strings.Split(bootnodesStr, ",")
		addrs, err := shared.StringsToAddrInfos(strs)
		if err != nil {
			return nil, err
		}

		err = bootstrap(h, addrs)
		if err != nil {
			return nil, err
		}
	}

	return n, nil
}

func bootstrap(h host.Host, bns []peer.AddrInfo) error {
	ctx := context.Background()
	for _, bn := range bns {
		err := h.Connect(ctx, bn)
		if err != nil {
			return err
		}
	}
	return nil
}
