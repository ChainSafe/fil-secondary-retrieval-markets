// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: Apache-2.0, MIT

package provider

import (
	"context"
	"math/big"

	"github.com/ChainSafe/fil-secondary-retrieval-markets/shared"
	"github.com/ipfs/go-cid"
	logging "github.com/ipfs/go-log/v2"
)

var log = logging.Logger("provider")

// Provider ...
type Provider struct {
	host Host
	msgs <-chan []byte
}

// NewProvider returns a new Provider
func NewProvider(host Host) *Provider {
	return &Provider{
		host: host,
	}
}

// Start starts the provider
func (p *Provider) Start() error {
	err := p.host.Start()
	if err != nil {
		return err
	}

	p.msgs = p.host.Messages()
	go p.handleMessages()
	return nil
}

// Stop stops the provider
func (p *Provider) Stop() error {
	return p.host.Stop()
}

func (p *Provider) handleMessages() {
	for msg := range p.msgs {
		query := new(shared.Query)
		err := query.Unmarshal(msg)
		if err != nil {
			log.Error("cannot unmarshal query; error:", err)
			continue
		}

		log.Info("received query!", query)

		if p.hasData(query.PayloadCID) {
			err = p.sendResponse(query)
			if err != nil {
				log.Error("cannot send response; error: ", err)
			}
		}
	}
}

func (p *Provider) sendResponse(query *shared.Query) error {
	if len(query.Client) == 0 {
		return ErrNoAddrsProvided
	}

	resp := &shared.QueryResponse{
		PayloadCID:              query.PayloadCID,
		Provider:                p.host.MultiAddrs(),
		Total:                   big.NewInt(0),
		PaymentInterval:         0,
		PaymentIntervalIncrease: 0,
	}

	addrs, err := shared.StringsToAddrInfos(query.Client)
	if err != nil {
		log.Error("cannot convert client addrs to multiaddrs; error: ", err)
		return err
	}

	for i, addr := range addrs {
		err = p.host.Connect(addr)
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
	return p.host.Send(context.Background(), addrs[0].ID, bz)
}

func (p *Provider) hasData(data cid.Cid) bool {
	// TODO: implement this using actual data store
	return true
}
