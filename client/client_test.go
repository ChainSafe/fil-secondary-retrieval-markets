// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: Apache-2.0, MIT

package client

import (
	"context"
	"encoding/json"
	"math/big"
	"testing"

	"github.com/ChainSafe/fil-secondary-retrieval-markets/shared"
	"github.com/ipfs/go-cid"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/multiformats/go-multiaddr"
	"github.com/stretchr/testify/require"
)

var testMultiAddr = multiaddr.StringCast("/ip4/1.2.3.4/tcp/5678/p2p/QmYyQSo1c1Ym7orWxLYvCrM2EmxFTANf8wXmmE7DWjhx5N")

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

func (h *mockHost) MultiAddrs() []string {
	return []string{
		testMultiAddr.String(),
	}
}

func (h *mockHost) RegisterStreamHandler(handler network.StreamHandler) {
	panic("not implemented")
}

func TestClient_SubmitQuery(t *testing.T) {
	host := &mockHost{queries: []shared.Query{}}
	client := NewClient(host)

	testCid, err := cid.Decode("bafybeierhgbz4zp2x2u67urqrgfnrnlukciupzenpqpipiz5nwtq7uxpx4")
	require.NoError(t, err)

	query := shared.Query{
		PayloadCID:  testCid,
		ClientAddrs: []string{testMultiAddr.String()},
	}

	err = client.SubmitQuery(context.Background(), testCid)
	require.NoError(t, err)

	require.ElementsMatch(t, []shared.Query{query}, host.queries)
}

func TestClient_HandleProviderResponse(t *testing.T) {
	t.Skip("test incomplete")

	host := &mockHost{queries: []shared.Query{}}
	client := NewClient(host)

	testCid, err := cid.Decode("bafybeierhgbz4zp2x2u67urqrgfnrnlukciupzenpqpipiz5nwtq7uxpx4")
	require.NoError(t, err)

	testPeerId, err := peer.Decode("QmYyQSo1c1Ym7orWxLYvCrM2EmxFTANf8wXmmE7DWjhx5N")
	require.NoError(t, err)

	response := shared.QueryResponse{
		PayloadCID:              testCid,
		Provider:                testPeerId,
		Total:                   big.NewInt(10),
		PaymentInterval:         0,
		PaymentIntervalIncrease: 0,
	}

	bz, err := json.Marshal(&response)
	require.NoError(t, err)

	client.HandleProviderResponse(bz)

	// TODO: Verify message is properly handled
}
