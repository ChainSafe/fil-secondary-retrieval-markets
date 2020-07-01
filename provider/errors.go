// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: Apache-2.0, MIT

package provider

import (
	"errors"
	"fmt"
)

// ErrNoAddrsProvided is returned when a client message is received that has no client multiaddrs.
var ErrNoAddrsProvided = errors.New("no client multiaddrs provided")

// ErrCannotConnect is returned when a provider is unable to connect using any of the client's multiaddrs
func ErrCannotConnect(err error) error {
	return fmt.Errorf("cannot connect to any provided multiaddrs: %s", err.Error())
}
