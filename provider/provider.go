// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: Apache-2.0, MIT

package provider

type Provider struct {
	host Host
}

func NewClient(host Host) *Provider {
	return &Provider{host: host}
}