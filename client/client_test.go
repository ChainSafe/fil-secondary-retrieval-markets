// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: Apache-2.0, MIT

package client

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/ChainSafe/fil-secondary-retrieval-markets/provider"
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

var testCid, _ = cid.Decode("bafybeierhgbz4zp2x2u67urqrgfnrnlukciupzenpqpipiz5nwtq7uxpx4")
var testParams = shared.Params{
	PayloadCID: testCid,
}

type mockNetwork struct{ queries []shared.Query }

func (n *mockNetwork) Start() error {
	return nil
}

func (n *mockNetwork) Stop() error {
	return nil
}

func (n *mockNetwork) Publish(ctx context.Context, data []byte) error {
	var query shared.Query
	err := json.Unmarshal(data, &query)
	if err != nil {
		return err
	}

	n.queries = append(n.queries, query)
	return nil
}

func (n *mockNetwork) MultiAddrs() []string {
	return []string{
		testMultiAddr.String(),
	}
}

func (n *mockNetwork) RegisterStreamHandler(id core.ProtocolID, handler network.StreamHandler) {}

func TestMain(m *testing.M) {
	lvl, err := logging.LevelFromString("debug")
	if err != nil {
		panic(err)
	}
	logging.SetAllLoggers(lvl)

	os.Exit(m.Run())
}

func TestClient_SubmitQuery(t *testing.T) {
	host := &mockNetwork{queries: []shared.Query{}}
	client := NewClient(host)

	query := shared.Query{
		Params:      testParams,
		ClientAddrs: []string{testMultiAddr.String()},
	}

	err := client.SubmitQuery(context.Background(), testParams)
	require.NoError(t, err)

	require.ElementsMatch(t, []shared.Query{query}, host.queries)
}

func TestClient_SubscribeToQueryResponses(t *testing.T) {
	host := &mockNetwork{queries: []shared.Query{}}
	client := NewClient(host)

	testPeerId, err := peer.Decode("QmYyQSo1c1Ym7orWxLYvCrM2EmxFTANf8wXmmE7DWjhx5N")
	require.NoError(t, err)

	response := shared.QueryResponse{
		Params:                  testParams,
		Provider:                testPeerId,
		PricePerByte:            provider.DefaultPricePerByte,
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

	unsubA := client.SubscribeToQueryResponses(subscriberA, testParams)
	unsubB := client.SubscribeToQueryResponses(subscriberB, testParams)
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
