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
	"fmt"
	"net/url"
)

type (
	// ScriptListItem item of script list.
	ScriptListItem struct {
		// Hash of the script creation transaction
		CreationTxHash TxHash `json:"creation_tx_hash"`

		// Hash of a script
		ScriptHash string `json:"script_hash"`
	}

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
		DatumValue map[string]interface{} `json:"datum_value"`

		// The budget in fees to run a script - the fees depend on the
		// ExUnits and the current prices.
		Fee Lovelace `json:"fee,omitempty"`

		// What kind pf validation this redeemer is used for,
		// it can be one of 'spend', 'mint', 'cert', 'reward'.
		Purpose string `json:"purpose"`

		// TxHash of Transaction containing the redeemer.
		TxHash TxHash `json:"tx_hash"`

		// TxIndex The index of the redeemer pointer in the transaction.
		TxIndex int `json:"tx_index"`

		// The budget in Memory to run a script.
		UnitMem int `json:"unit_mem"`

		// The budget in Cpu steps to run a script.
		UnitSteps int `json:"unit_steps"`
	}

	// ScriptListResponse represents response from `/script_list` endpoint.
	ScriptListResponse struct {
		Response
		Data []ScriptListItem `json:"response"`
	}

	// ScriptRedeemersResponse represents response from `/script_redeemers` endpoint.
	ScriptRedeemersResponse struct {
		Response
		Data *ScriptRedeemers `json:"response"`
	}
)

// GetScriptList returns the list of all existing script
// hashes along with their creation transaction hashes.
func (c *Client) GetScriptList(ctx context.Context) (res *ScriptListResponse, err error) {
	res = &ScriptListResponse{}
	rsp, _ := c.request(ctx, &res.Response, "GET", "/script_list", nil, nil, nil)
	body, err := readResponseBody(rsp)
	if err != nil {
		res.applyError(body, err)
		return
	}

	if err = json.Unmarshal(body, &res.Data); err != nil {
		res.applyError(body, err)
		return
	}

	return res, nil
}

// GetScriptRedeemers returns a list of all redeemers for a given script hash.
func (c *Client) GetScriptRedeemers(
	ctx context.Context,
	sh ScriptHash,
) (res *ScriptRedeemersResponse, err error) {
	res = &ScriptRedeemersResponse{}

	params := url.Values{}
	params.Set("_script_hash", fmt.Sprint(sh))

	rsp, _ := c.request(ctx, &res.Response, "GET", "/script_redeemers", nil, params, nil)
	body, err := readResponseBody(rsp)
	if err != nil {
		res.applyError(body, err)
		return
	}
	r := []ScriptRedeemers{}

	if err = json.Unmarshal(body, &r); err != nil {
		res.applyError(body, err)
		return
	}

	if len(r) == 1 {
		res.Data = &r[0]
	}
	res.ready()
	return res, nil
}
