package network

import (
	"errors"
)

// ErrNoConfig is returned when no configuration is provided for a host
var ErrNoConfig = errors.New("no configuration provided")
