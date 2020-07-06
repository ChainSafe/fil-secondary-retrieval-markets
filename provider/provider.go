// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: Apache-2.0, MIT

package provider

import (
	"context"
	"math/big"

	"github.com/ChainSafe/fil-secondary-retrieval-markets/shared"
	"github.com/ipfs/go-cid"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	logging "github.com/ipfs/go-log/v2"
)

var log = logging.Logger("provider")

// Provider ...
type Provider struct {
	net        Network
	blockstore blockstore.Blockstore
	msgs       <-chan []byte
}

// NewProvider returns a new Provider
func NewProvider(net Network, bs blockstore.Blockstore) *Provider {
	return &Provider{
		net:        net,
		blockstore: bs,
	}
}

// Start starts the provider
func (p *Provider) Start() error {
	err := p.net.Start()
	if err != nil {
		return err
	}

	p.msgs = p.net.Messages()
	go p.handleMessages()
	return nil
}

// Stop stops the provider
func (p *Provider) Stop() error {
	return p.net.Stop()
}

func (p *Provider) handleMessages() {
	for msg := range p.msgs {
		query := new(shared.Query)
		err := query.Unmarshal(msg)
		if err != nil {
			log.Error("cannot unmarshal query; error:", err)
			continue
		}

		log.Debug("received query for CID", query.PayloadCID)
		has, err := p.hasData(query.PayloadCID)
		if err != nil {
			log.Error("failed to check for data in blockstore; error:", err)
			continue
		}

		if has {
			err = p.sendResponse(query)
			if err != nil {
				log.Error("cannot send response; error: ", err)
			}
		}
	}
}

func (p *Provider) sendResponse(query *shared.Query) error {
	if len(query.ClientAddrs) == 0 {
		return ErrNoAddrsProvided
	}

	resp := &shared.QueryResponse{
		PayloadCID:              query.PayloadCID,
		Provider:                p.net.PeerID(),
		Total:                   big.NewInt(0),
		PaymentInterval:         0,
		PaymentIntervalIncrease: 0,
	}

	addrs, err := shared.StringsToAddrInfos(query.ClientAddrs)
	if err != nil {
		log.Error("cannot convert client addrs to multiaddrs; error: ", err)
		return err
	}

	for i, addr := range addrs {
		// TODO: check if already connected using client's peer ID
		err = p.net.Connect(addr)
		if err != nil {
			log.Error("failed to connect to addr: ", err)

			// couldn't connect using any addrs
			if i == len(addrs)-1 {
				return ErrConnectFailed
			}

			continue
		}
	}

	bz, err := resp.Marshal()
	if err != nil {
		return err
	}

	// TODO: if we open up a substream with the client, what protocol ID do we use?
	// or do we use the existing /fil/markets stream?
	return p.net.Send(context.Background(), shared.ResponseProtocolID, addrs[0].ID, bz)
}

func (p *Provider) hasData(data cid.Cid) (bool, error) {
	return p.blockstore.Has(data)
}
