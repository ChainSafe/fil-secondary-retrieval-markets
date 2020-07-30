// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: Apache-2.0, MIT

package main

import (
	"github.com/ChainSafe/fil-secondary-retrieval-markets/shared"
	"github.com/ipfs/go-cid"
)

type ProviderStore struct {
	cids map[cid.Cid]struct{}
}

func (s *ProviderStore) Has(params shared.Params) (bool, error) {
	if _, has := s.cids[params.PayloadCID]; has {
		return true, nil
	}

	if _, has := s.cids[*params.PieceCID]; has {
		return true, nil
	}

	return false, nil
}

type ProviderStoreJSON struct {
	cids []string
}

func (s *ProviderStoreJSON) ToProviderStore() *ProviderStore {
	ps := &ProviderStore{
		cids: make(map[cid.Cid]struct{}),
	}

	for _, s := range s.cids {
		cid, err := cid.Decode(s)
		if err != nil {
			continue
		}

		ps.cids[cid] = struct{}{}
	}

	return ps
}
