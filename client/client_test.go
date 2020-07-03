// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: Apache-2.0, MIT

package client

import (
	"context"
	"encoding/json"
	"math/big"
	"os"
	"testing"

	"github.com/ChainSafe/fil-secondary-retrieval-markets/shared"
	"github.com/ipfs/go-cid"
	logging "github.com/ipfs/go-log/v2"
	core "github.com/libp2p/go-libp2p-core"
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

func (h *mockHost) RegisterStreamHandler(id core.ProtocolID, handler network.StreamHandler) {}

func TestMain(m *testing.M) {
	lvl, err := logging.LevelFromString("debug")
	if err != nil {
		panic(err)
	}
	logging.SetAllLoggers(lvl)

	os.Exit(m.Run())
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

func TestClient_SubscribeToQueryResponses(t *testing.T) {
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

	// First, setup two subscribers

	// Use buffered channel to avoid blocking
	responsesA := make(chan shared.QueryResponse, 1)
	responsesB := make(chan shared.QueryResponse, 1)

	subscriberA := func(resp shared.QueryResponse) {
		responsesA <- resp
	}
	subscriberB := func(resp shared.QueryResponse) {
		responsesB <- resp
	}

	unsubA := client.SubscribeToQueryResponses(subscriberA, testCid)
	unsubB := client.SubscribeToQueryResponses(subscriberB, testCid)
	defer unsubB()

	// Process response and wait for result
	client.HandleProviderResponse(bz)

	select {
	case actual := <-responsesA:
		require.Equal(t, response, actual)
	default:
		t.Fatal("no response received for subscriberA")
	}

	select {
	case actual := <-responsesB:
		require.Equal(t, response, actual)
	default:
		t.Fatal("no response received for subscriberB")
	}

	// Now lets unsub A and verify no response is received
	unsubA()
	client.HandleProviderResponse(bz)

	select {
	case <-responsesA:
		t.Fatal("expected no response for subscriberA")
	default:
	}

	select {
	case actual := <-responsesB:
		require.Equal(t, response, actual)
	default:
		t.Fatal("no response received for subscriberB")
	}
}
