// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: Apache-2.0, MIT

package provider

import (
	"errors"
)

// ErrNoAddrsProvided is returned when a client message is received that has no client multiaddrs.
var ErrNoAddrsProvided = errors.New("no client multiaddrs provided")
