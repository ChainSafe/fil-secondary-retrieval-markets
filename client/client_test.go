package client

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/ChainSafe/fil-secondary-retrieval-markets/shared"
	"github.com/ipfs/go-cid"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/stretchr/testify/require"
)

type mockHost struct{ queries []shared.Query }

func (h *mockHost) Publish(ctx context.Context, data []byte) error {
	var query shared.Query
	err := json.Unmarshal(data, &query)
	if err != nil {
		return err
	}

	h.queries = append(h.queries, query)
	return nil
}

func TestClient_SubmitQuery(t *testing.T) {
	host := &mockHost{queries: []shared.Query{}}
	client := NewClient(host)

	testCid, err := cid.Decode("bafybeierhgbz4zp2x2u67urqrgfnrnlukciupzenpqpipiz5nwtq7uxpx4")
	require.NoError(t, err)

	testPeer, err := peer.Decode("QmYyQSo1c1Ym7orWxLYvCrM2EmxFTANf8wXmmE7DWjhx5N")
	require.NoError(t, err)

	query := shared.Query{
		PayloadCID: testCid,
		Client:     testPeer,
	}

	err = client.SubmitQuery(context.Background(), query)
	require.NoError(t, err)

	require.ElementsMatch(t, []shared.Query{query}, host.queries)
}
