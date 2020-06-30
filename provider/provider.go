// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: Apache-2.0, MIT

package provider

import (
	"log"
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
		log.Println("received message!", msg)
	}
}
