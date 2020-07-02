package client

import (
	"github.com/ChainSafe/fil-secondary-retrieval-markets/shared"
	"github.com/ipfs/go-cid"
)

type QuerySubscription struct {
	// Client will push responses here, listeners will read
	ch chan shared.QueryResponse
}

type QuerySubscriptionStore map[cid.Cid]*QuerySubscription
