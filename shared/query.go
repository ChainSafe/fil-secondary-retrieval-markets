package shared

import (
	"github.com/ipfs/go-cid"
	"github.com/libp2p/go-libp2p-core/peer"
)

// Query is submitted by clients and observed by providers
type Query struct {
	PayloadCID cid.Cid `json:"payloadCID"`
	Client     peer.ID `json:"client"`
}
