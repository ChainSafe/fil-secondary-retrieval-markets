// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: Apache-2.0, MIT

package provider

import (
	"github.com/ChainSafe/fil-secondary-retrieval-markets/shared"
)

type RetrievalProviderStore interface {
	Has(params shared.Params) (bool, error)
}
