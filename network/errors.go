// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: Apache-2.0, MIT

package network

import (
	"errors"
)

// ErrNoConfig is returned when no configuration is provided for a host
//var ErrNoConfig = errors.New("no configuration provided")

// ErrNilHost is returned when trying to instantiate a network with a nil host
var ErrNilHost = errors.New("host is nil")
