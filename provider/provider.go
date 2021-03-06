// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: Apache-2.0, MIT

package provider

import (
	"context"
	"reflect"
	"sync"

	"github.com/ChainSafe/fil-secondary-retrieval-markets/shared"
	"github.com/filecoin-project/specs-actors/actors/abi"
	logging "github.com/ipfs/go-log/v2"
)

var log = logging.Logger("provider")

// DefaultPricePerByte is the charge per byte retrieved if the miner does
// not specifically set it
var DefaultPricePerByte = abi.NewTokenAmount(2)

// DefaultPaymentInterval is the baseline interval, set to 1Mb
// if the miner does not explicitly set it otherwise
var DefaultPaymentInterval = uint64(1 << 20)

// DefaultPaymentIntervalIncrease is the amount interval increases on each payment,
// set to to 1Mb if the miner does not explicitly set it otherwise
var DefaultPaymentIntervalIncrease = uint64(1 << 20)

type ProviderSubscriber func(query shared.Query)
type Unsubscribe func()

// Provider ...
type Provider struct {
	net             Network
	store           RetrievalProviderStore
	cache           RequestCache
	msgs            <-chan []byte
	subscribers     []ProviderSubscriber
	subscribersLock sync.Mutex

	pricePerByte            abi.TokenAmount
	paymentInterval         uint64
	paymentIntervalIncrease uint64
	priceLock               sync.Mutex
}

// NewProvider returns a new Provider
func NewProvider(net Network, s RetrievalProviderStore, cache RequestCache) *Provider {
	return &Provider{
		net:                     net,
		store:                   s,
		cache:                   cache,
		subscribers:             []ProviderSubscriber{},
		pricePerByte:            DefaultPricePerByte,
		paymentInterval:         DefaultPaymentInterval,
		paymentIntervalIncrease: DefaultPaymentIntervalIncrease,
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

// SetPricePerByte sets the provider's pricePerByte
func (p *Provider) SetPricePerByte(price abi.TokenAmount) {
	p.priceLock.Lock()
	defer p.priceLock.Unlock()
	p.pricePerByte = price
}

// SetPaymentInterval sets the provider's paymentInterval and paymentIntervalIncrease
func (p *Provider) SetPaymentInterval(interval, increase uint64) {
	p.priceLock.Lock()
	defer p.priceLock.Unlock()
	p.paymentInterval = interval
	p.paymentIntervalIncrease = increase
}

// SubscribeToQueries registers the given subscriber and calls it upon receiving queries
func (p *Provider) SubscribeToQueries(s ProviderSubscriber) Unsubscribe {
	p.subscribersLock.Lock()
	defer p.subscribersLock.Unlock()

	p.subscribers = append(p.subscribers, s)
	return p.unsubscribeAt(s)
}

func (p *Provider) unsubscribeAt(s ProviderSubscriber) Unsubscribe {
	return func() {
		p.subscribersLock.Lock()
		defer p.subscribersLock.Unlock()
		curLen := len(p.subscribers)
		for i, el := range p.subscribers {
			if reflect.ValueOf(s) == reflect.ValueOf(el) {
				p.subscribers[i] = p.subscribers[curLen-1]
				p.subscribers = p.subscribers[:curLen-1]
				return
			}
		}
	}
}

func (p *Provider) handleMessages() {
	for msg := range p.msgs {
		query := new(shared.Query)
		err := query.Unmarshal(msg)
		if err != nil {
			log.Error("cannot unmarshal query; error:", err)
			continue
		}

		p.notifySubscribers(*query)

		log.Info("received query for params", query.Params)
		has, err := p.hasData(query.Params)
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

		// TODO: update cache to accept params?
		p.cache.Put(query.Params.PayloadCID)
	}
}

func (p *Provider) notifySubscribers(query shared.Query) {
	p.subscribersLock.Lock()
	defer p.subscribersLock.Unlock()

	for _, s := range p.subscribers {
		s(query)
	}
}

func (p *Provider) sendResponse(query *shared.Query) error {
	if len(query.ClientAddrs) == 0 {
		return ErrNoAddrsProvided
	}

	resp := &shared.QueryResponse{
		Params:                  query.Params,
		Provider:                p.net.PeerID(),
		PricePerByte:            p.pricePerByte,
		PaymentInterval:         p.paymentInterval,
		PaymentIntervalIncrease: p.paymentIntervalIncrease,
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

func (p *Provider) hasData(params shared.Params) (bool, error) {
	return p.store.Has(params)
}
