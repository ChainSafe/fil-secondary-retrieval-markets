package provider

import (
	"github.com/ChainSafe/fil-secondary-retrieval-markets/shared"
)

type RetrievalProviderStore interface {
	Has(params shared.Params) (bool, error)
}
