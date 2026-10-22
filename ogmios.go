// SPDX-License-Identifier: Apache-2.0
//
// Copyright © 2022 The Cardano Community Authors

package koios

import (
	"context"
	"encoding/json"
	"errors"
)

// ErrNoOgmiosMethod is returned when an Ogmios request has no method.
var ErrNoOgmiosMethod = errors.New("missing ogmios method")

type (
	// OgmiosMethod is the Ogmios JSON-RPC method accepted by `/ogmios`.
	OgmiosMethod string

	// OgmiosRequest is the JSON-RPC 2.0 request body of `/ogmios`.
	OgmiosRequest struct {
		// JSONRPC is the JSON-RPC version identifier. It defaults to "2.0".
		JSONRPC string `json:"jsonrpc"`

		// Method is the Ogmios method to call.
		Method OgmiosMethod `json:"method"`

		// Params are the parameters of the method. They are omitted when nil.
		Params any `json:"params,omitempty"`
	}

	// OgmiosResponse represents response from `/ogmios` endpoint.
	OgmiosResponse struct {
		Response

		// Data is the raw JSON-RPC response document, which carries the
		// jsonrpc version, the method that was called and its result.
		Data json.RawMessage `json:"data"`
	}
)

// Ogmios methods supported by `/ogmios`.
const (
	OgmiosQueryNetworkBlockHeight                OgmiosMethod = "queryNetwork/blockHeight"
	OgmiosQueryNetworkGenesisConfiguration       OgmiosMethod = "queryNetwork/genesisConfiguration"
	OgmiosQueryNetworkStartTime                  OgmiosMethod = "queryNetwork/startTime"
	OgmiosQueryNetworkTip                        OgmiosMethod = "queryNetwork/tip"
	OgmiosQueryLedgerStateEpoch                  OgmiosMethod = "queryLedgerState/epoch"
	OgmiosQueryLedgerStateEraStart               OgmiosMethod = "queryLedgerState/eraStart"
	OgmiosQueryLedgerStateEraSummaries           OgmiosMethod = "queryLedgerState/eraSummaries"
	OgmiosQueryLedgerStateLiveStakeDistribution  OgmiosMethod = "queryLedgerState/liveStakeDistribution"
	OgmiosQueryLedgerStateProtocolParameters     OgmiosMethod = "queryLedgerState/protocolParameters"
	OgmiosQueryLedgerStateProposedProtocolParams OgmiosMethod = "queryLedgerState/proposedProtocolParameters"
	OgmiosQueryLedgerStateStakePools             OgmiosMethod = "queryLedgerState/stakePools"
	OgmiosSubmitTransaction                      OgmiosMethod = "submitTransaction"
	OgmiosEvaluateTransaction                    OgmiosMethod = "evaluateTransaction"
)

// Ogmios sends a JSON-RPC 2.0 request to the Ogmios endpoint of the API and
// returns the raw JSON-RPC response document.
func (c *Client) Ogmios(
	ctx context.Context,
	req OgmiosRequest,
	opts *RequestOptions,
) (*OgmiosResponse, error) {
	res := &OgmiosResponse{}
	if req.Method == "" {
		res.applyError(nil, ErrNoOgmiosMethod)
		return res, ErrNoOgmiosMethod
	}
	if req.JSONRPC == "" {
		req.JSONRPC = "2.0"
	}
	body, err := jsonBody(req)
	if err != nil {
		res.applyError(nil, err)
		return res, err
	}

	rsp, err := c.request(ctx, &res.Response, "POST", "/ogmios", body, opts)
	if err != nil {
		return res, err
	}
	return res, ReadAndUnmarshalResponse(rsp, &res.Response, &res.Data)
}
