// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: Apache-2.0, MIT

package provider

import (
	"context"
	"encoding/json"
	"log"
	"math/big"

	"github.com/ChainSafe/fil-secondary-retrieval-markets/shared"
	"github.com/ipfs/go-cid"
)

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
		err := json.Unmarshal(msg, query)
		if err != nil {
			log.Println("cannot unmarshal query; error:", err)
			continue
		}

		log.Println("received query!", query)
		if p.hasData(query.PayloadCID) {
			err = p.sendResponse(query)
			if err != nil {
				log.Println("cannot send response; error:", err)
			}
		}
	}
}

func (p *Provider) sendResponse(query *shared.Query) error {
	resp := &shared.QueryResponse{
		PayloadCID:              query.PayloadCID,
		Provider:                p.host.MultiAddrs(),
		Total:                   big.NewInt(0),
		PaymentInterval:         0,
		PaymentIntervalIncrease: 0,
	}

	addrs, err := shared.AddrsToAddrInfos(query.Client)
	if err != nil {
		return err
	}

	for i, addr := range addrs {
		err = p.host.Connect(addr)
		if err != nil {
			// couldn't connect using any addrs
			if i == len(addrs)-1 {
				return err
			}

			continue
		}
	}

	bz, err := json.Marshal(resp)
	if err != nil {
		return err
	}

	return p.host.Send(context.Background(), addrs[0].ID, bz)
}

func (p *Provider) hasData(data cid.Cid) bool {
	// TODO: implement this using actual data store
	return true
}
