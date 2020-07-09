// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: Apache-2.0, MIT

package cache

import (
	"time"
)

type Record struct {
	Frequency     int
	LastAccessed  time.Time
	InsertionTime time.Time
}
