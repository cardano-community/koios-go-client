// Copyright 2022 The Howijd.Network Authors
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
	"net/http"
	"net/url"
)

type (
	// AddressUTxO UTxO attached to address.
	AddressUTxO struct {
		// Hash of Transaction for input UTxO.
		TxHash TxHash `json:"tx_hash"`

		// Index of input UTxO on the mentioned address used for input.
		TxIndex int `json:"tx_index"`

		// Balance on the selected input transaction.
		Value Lovelace `json:"value"`

		// An array of assets contained on UTxO.
		AssetList []Asset `json:"asset_list,omitempty"`
	}

	// AddressInfo esponse for `/address_info`.
	AddressInfo struct {
		// Balance ADA Lovelace balance of address
		Balance Lovelace `json:"balance"`
		// StakeAddress associated with address
		StakeAddress StakeAddress  `json:"stake_address, omitempty"`
		UTxOs        []AddressUTxO `json:"utxo_set"`
	}

	// AddressInfoResponse represents response from `/address_info` endpoint.
	AddressInfoResponse struct {
		Response
		Data *AddressInfo `json:"response"`
	}
)

// GetAddressInfo returns address info - balance,
// associated stake address (if any) and UTxO set.
func (c *Client) GetAddressInfo(ctx context.Context, addr Address) (res *AddressInfoResponse, err error) {
	res = &AddressInfoResponse{}
	if len(addr) == 0 {
		err = ErrNoAddress
		res.applyError(nil, err)
		return
	}
	params := url.Values{}
	params.Set("_address", string(addr))

	rsp, err := c.request(ctx, &res.Response, "GET", nil, "/address_info", params, nil)
	if err != nil {
		return
	}
	body, err := readResponseBody(rsp)
	if err != nil {
		res.applyError(nil, err)
		return
	}

	addrs := []AddressInfo{}

	if err = json.Unmarshal(body, &addrs); err != nil {
		res.applyError(body, err)
		return
	}

	if rsp.StatusCode != http.StatusOK {
		res.applyError(body, err)
		return
	}
	if len(addrs) == 1 {
		res.Data = &addrs[0]
	}
	res.ready()
	return res, nil
}
