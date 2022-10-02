// Copyright 2022 The Cardano Community Authors
// SPDX-License-Identifier: Apache-2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//
//   http://www.apache.org/licenses/LICENSE-2.0
//   or LICENSE file in repository root.
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package koios

import (
	"context"
	"encoding/json"

	"github.com/shopspring/decimal"
)

type (
	// ScriptRedeemers defines model for script_redeemers.
	ScriptRedeemers struct {
		// Hash of Transaction for which details are being shown
		ScriptHash ScriptHash `json:"script_hash"`

		// Redeemers list
		Redeemers []ScriptRedeemer `json:"redeemers"`
	}

	// ScriptRedeemer model.
	ScriptRedeemer struct {
		// The Hash of the Plutus Data
		DatumHash string `json:"datum_hash"`

		// The actual data in json format
		DatumValue map[string]any `json:"datum_value"`

		// The budget in fees to run a script - the fees depend on the
		// ExUnits and the current prices.
		Fee decimal.Decimal `json:"fee,omitempty"`

		// What kind pf validation this redeemer is used for,
		// it can be one of 'spend', 'mint', 'cert', 'reward'.
		Purpose string `json:"purpose"`

		// TxHash of Transaction containing the redeemer.
		TxHash TxHash `json:"tx_hash"`

		// TxIndex The index of the redeemer pointer in the transaction.
		TxIndex uint32 `json:"tx_index"`

		// The budget in Memory to run a script.
		UnitMem int `json:"unit_mem"`

		// The budget in Cpu steps to run a script.
		UnitSteps int `json:"unit_steps"`
	}

	PlutusContract struct {
		Address       Address    `json:"address"`
		ScriptHash    ScriptHash `json:"script_hash"`
		ByteCode      string     `json:"bytecode"`
		Size          int        `json:"size"`
		ValidContract bool       `json:"valid_contract"`
	}

	// NativeScript item of native script list.
	NativeScript struct {
		// Hash of the script creation transaction
		CreationTxHash TxHash `json:"creation_tx_hash"`

		// Hash of a script
		ScriptHash string `json:"script_hash"`
		Type       string `json:"type"`
		Script     struct {
			Type    string            `json:"type"`
			Scripts []json.RawMessage `json:"scripts"`
		} `json:"script"`
	}

	// NativeScriptListResponse represents response from `/native_script_list` endpoint.
	NativeScriptListResponse struct {
		Response
		Data []NativeScript `json:"response"`
	}

	// PlutusScriptListItem item of plutus script list.
	PlutusScriptListItem struct {
		// Hash of the script creation transaction
		CreationTxHash TxHash `json:"creation_tx_hash"`

		// Hash of a script
		ScriptHash string `json:"script_hash"`
	}

	// PlutusScriptListResponse represents response from `/plutus_script_list` endpoint.
	PlutusScriptListResponse struct {
		Response
		Data []PlutusScriptListItem `json:"response"`
	}

	// ScriptRedeemersResponse represents response from `/script_redeemers` endpoint.
	ScriptRedeemersResponse struct {
		Response
		Data *ScriptRedeemers `json:"response"`
	}
)

// GetNativeScriptList returns list of all existing native script hashes
// along with their creation transaction hashes.
func (c *Client) GetNativeScriptList(
	ctx context.Context,
	opts *RequestOptions,
) (res *NativeScriptListResponse, err error) {
	res = &NativeScriptListResponse{}
	rsp, err := c.request(ctx, &res.Response, "GET", "/native_script_list", nil, opts)
	if err != nil {
		return
	}
	err = ReadAndUnmarshalResponse(rsp, &res.Response, &res.Data)
	return
}

// GetPlutusScriptList returns all existing Plutus script
// hashes along with their creation transaction hashes.
func (c *Client) GetPlutusScriptList(
	ctx context.Context,
	opts *RequestOptions,
) (res *PlutusScriptListResponse, err error) {
	res = &PlutusScriptListResponse{}
	rsp, err := c.request(ctx, &res.Response, "GET", "/plutus_script_list", nil, opts)
	if err != nil {
		return
	}
	err = ReadAndUnmarshalResponse(rsp, &res.Response, &res.Data)
	return
}

// GetScriptRedeemers returns a list of all redeemers for a given script hash.
func (c *Client) GetScriptRedeemers(
	ctx context.Context,
	sh ScriptHash,
	opts *RequestOptions,
) (res *ScriptRedeemersResponse, err error) {
	res = &ScriptRedeemersResponse{}

	if opts == nil {
		opts = c.NewRequestOptions()
	}
	opts.QuerySet("_script_hash", sh.String())

	rsp, err := c.request(ctx, &res.Response, "GET", "/script_redeemers", nil, opts)
	if err != nil {
		return
	}
	r := []ScriptRedeemers{}
	err = ReadAndUnmarshalResponse(rsp, &res.Response, &r)

	if len(r) == 1 {
		res.Data = &r[0]
	}
	return
}
