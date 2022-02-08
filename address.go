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
	"io"
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
		StakeAddress StakeAddress `json:"stake_address,omitempty"`

		UTxOs []AddressUTxO `json:"utxo_set"`
	}

	// AddressInfoResponse represents response from `/address_info` endpoint.
	AddressInfoResponse struct {
		Response
		Data *AddressInfo `json:"response"`
	}

	// AddressTxsResponse represents response from `/address_txs` endpoint.
	AddressTxsResponse struct {
		Response
		Data []TxHash `json:"response"`
	}

	// CredentialTxsResponse represents response from `/credential_txs` endpoint.
	CredentialTxsResponse struct {
		Response
		Data []TxHash `json:"response"`
	}

	// AddressAsset payload item returned by.
	AddressAsset struct {
		// Asset Name (hex).
		NameHEX string `json:"asset_name_hex"`

		// Asset Policy ID (hex).
		PolicyID PolicyID `json:"asset_policy_hex"`

		// Quantity of asset accoiated to the given address.
		Quantity Lovelace `json:"quantity"`
	}

	// AddressAssetsResponse represents response from `/address_info` endpoint.
	AddressAssetsResponse struct {
		Response
		Data []AddressAsset `json:"response"`
	}
)

// GetAddressInfo returns address info - balance,
// associated stake address (if any) and UTxO set.
//nolint: dupl
func (c *Client) GetAddressInfo(ctx context.Context, addr Address) (res *AddressInfoResponse, err error) {
	res = &AddressInfoResponse{}
	if len(addr) == 0 {
		err = ErrNoAddress
		res.applyError(nil, err)
		return
	}
	params := url.Values{}
	params.Set("_address", string(addr))

	rsp, err := c.request(ctx, &res.Response, "GET", "/address_info", nil, params, nil)
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

// GetAddressTxs returns the transaction hash list of input address array,
// optionally filtering after specified block height (inclusive).
//nolint: dupl
func (c *Client) GetAddressTxs(ctx context.Context, addrs []Address, h uint64) (res *AddressTxsResponse, err error) {
	res = &AddressTxsResponse{}
	if len(addrs) == 0 {
		err = ErrNoAddress
		res.applyError(nil, err)
		return
	}

	var payload = struct {
		Adresses         []Address `json:"_addresses"`
		AfterBlockHeight uint64    `json:"_after_block_height,omitempty"`
	}{
		Adresses:         addrs,
		AfterBlockHeight: h,
	}

	rpipe, w := io.Pipe()
	go func() {
		_ = json.NewEncoder(w).Encode(payload)
		defer w.Close()
	}()

	rsp, err := c.request(ctx, &res.Response, "POST", "/address_txs", rpipe, nil, nil)
	if err != nil {
		return
	}
	body, err := readResponseBody(rsp)
	if err != nil {
		res.applyError(nil, err)
		return
	}

	atxs := []struct {
		Hash TxHash `json:"tx_hash"`
	}{}

	if err = json.Unmarshal(body, &atxs); err != nil {
		res.applyError(body, err)
		return
	}

	if rsp.StatusCode != http.StatusOK {
		res.applyError(body, err)
		return
	}
	if len(atxs) > 0 {
		for _, tx := range atxs {
			res.Data = append(res.Data, tx.Hash)
		}
	}
	res.ready()
	return res, nil
}

// GetAddressAssets returns the list of all the assets (policy, name and quantity)
// for a given address.
func (c *Client) GetAddressAssets(ctx context.Context, addr Address) (res *AddressAssetsResponse, err error) {
	res = &AddressAssetsResponse{}
	if len(addr) == 0 {
		err = ErrNoAddress
		res.applyError(nil, err)
		return
	}
	params := url.Values{}
	params.Set("_address", string(addr))

	rsp, err := c.request(ctx, &res.Response, "GET", "/address_assets", nil, params, nil)
	if err != nil {
		return
	}
	body, err := readResponseBody(rsp)
	if err != nil {
		res.applyError(nil, err)
		return
	}

	if err = json.Unmarshal(body, &res.Data); err != nil {
		res.applyError(body, err)
		return
	}

	if rsp.StatusCode != http.StatusOK {
		res.applyError(body, err)
		return
	}

	res.ready()
	return res, nil
}

// GetCredentialTxs returns the transaction hash list of input
// payment credential array, optionally filtering after specified block height (inclusive).
//nolint: dupl
func (c *Client) GetCredentialTxs(
	ctx context.Context,
	creds []PaymentCredential,
	h uint64,
) (res *CredentialTxsResponse, err error) {
	res = &CredentialTxsResponse{}
	if len(creds) == 0 {
		err = ErrNoAddress
		res.applyError(nil, err)
		return
	}

	var payload = struct {
		Credentials      []PaymentCredential `json:"_payment_credentials"`
		AfterBlockHeight uint64              `json:"_after_block_height,omitempty"`
	}{
		Credentials:      creds,
		AfterBlockHeight: h,
	}

	rpipe, w := io.Pipe()
	go func() {
		_ = json.NewEncoder(w).Encode(payload)
		defer w.Close()
	}()

	rsp, err := c.request(ctx, &res.Response, "POST", "/credential_txs", rpipe, nil, nil)
	if err != nil {
		return
	}
	body, err := readResponseBody(rsp)
	if err != nil {
		res.applyError(nil, err)
		return
	}

	atxs := []struct {
		Hash TxHash `json:"tx_hash"`
	}{}

	if err = json.Unmarshal(body, &atxs); err != nil {
		res.applyError(body, err)
		return
	}

	if rsp.StatusCode != http.StatusOK {
		res.applyError(body, err)
		return
	}
	if len(atxs) > 0 {
		for _, tx := range atxs {
			res.Data = append(res.Data, tx.Hash)
		}
	}
	res.ready()
	return res, nil
}
