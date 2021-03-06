// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: Apache-2.0, MIT

package main

import (
	"github.com/ChainSafe/fil-secondary-retrieval-markets/shared"
)

type responseHandler struct {
	respCh chan *shared.QueryResponse
}

func newResponseHandler() *responseHandler {
	return &responseHandler{
		respCh: make(chan *shared.QueryResponse),
	}
}

func (h *responseHandler) handleResponse(resp shared.QueryResponse) {
	h.respCh <- &resp
}
