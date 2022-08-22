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
)

type (
	// Asset represents Cardano Asset.
	Asset struct {
		// Asset Name (hex).
		Name string `json:"asset_name"`

		// Asset Policy ID (hex).
		PolicyID PolicyID `json:"policy_id"`

		// Quantity
		// Input: asset balance on the selected input transaction.
		// Output: sum of assets for output UTxO.
		// Mint: sum of minted assets (negative on burn).
		Quantity Lovelace `json:"quantity"`
	}

	// TokenRegistryMetadata metadata registered on the Cardano Token Registry.
	TokenRegistryMetadata struct {
		Decimals    int    `json:"decimals"`
		Description string `json:"description"`

		// A PNG image file as a byte string
		Logo   string `json:"logo"`
		Name   string `json:"name"`
		Ticker string `json:"ticker"`
		URL    string `json:"url"`
	}

	// AssetSummary aggregated asset summary.
	AssetSummary struct {
		// Asset Name (hex)
		AssetName AssetName `json:"asset_name"`

		// Asset Policy ID (hex)
		PolicyID PolicyID `json:"policy_id"`

		// Total number of registered wallets holding the given asset
		StakedWallets uint64 `json:"staked_wallets"`

		// Total number of transactions including the given asset
		TotalTransactions uint64 `json:"total_transactions"`

		// Total number of payment addresses (not belonging
		// to registered wallets) holding the given asset
		UnstakedAddresses uint64 `json:"unstaked_addresses"`
	}

	// AssetInfo info about the asset.
	AssetInfo struct {
		// Asset Name (hex).
		Name AssetName `json:"asset_name"`

		// Asset Name (ASCII)
		NameASCII string `json:"asset_name_ascii"`

		// The CIP14 fingerprint of the asset
		Fingerprint string `json:"fingerprint"`

		// MintingTxMetadata minting Tx JSON payload if it can be decoded as JSON
		MintingTxMetadata []TxInfoMetadata `json:"minting_tx_metadata"`

		// Asset metadata registered on the Cardano Token Registry
		TokenRegistryMetadata *TokenRegistryMetadata `json:"token_registry_metadata"`

		// Asset Policy ID (hex).
		PolicyID PolicyID `json:"policy_id,omitempty"`

		// TotalSupply of Asset
		TotalSupply Lovelace `json:"total_supply"`

		// CreationTime of Asset
		CreationTime Time `json:"creation_time"`

		// MintCnt count of mint transactions
		MintCnt int `json:"mint_cnt"`

		// BurnCnt count of burn transactions
		BurnCnt int `json:"burn_cnt"`

		// MintingTxHash mint tx
		MintingTxHash TxHash `json:"minting_tx_hash"`
	}

	// AssetListItem used to represent response from /asset_list`.
	AssetListItem struct {
		PolicyID   PolicyID `json:"policy_id"`
		AssetNames struct {
			HEX   []string `json:"hex"`
			ASCII []string `json:"ascii"`
		} `json:"asset_names"`
	}

	// AssetListResponse represents response from `/asset_list` endpoint.
	AssetListResponse struct {
		Response
		Data []AssetListItem `json:"response"`
	}

	// AssetHolder payment addresses holding the given token (including balance).
	AssetHolder struct {
		PaymentAddress Address  `json:"payment_address"`
		Quantity       Lovelace `json:"quantity"`
	}

	// AssetAddressListResponse represents response from `/asset_address_list` endpoint.
	AssetAddressListResponse struct {
		Response
		Data []AssetHolder `json:"response"`
	}

	// AssetInfoResponse represents response from `/asset_info` endpoint.
	AssetInfoResponse struct {
		Data *AssetInfo `json:"response"`
		Response
	}

	// AssetSummaryResponse represents response from `/asset_summary` endpoint.
	AssetSummaryResponse struct {
		Response
		Data *AssetSummary `json:"response"`
	}

	// AssetTxsResponse represents response from `/asset_txs` endpoint.
	AssetTxsResponse struct {
		Response
		Data []TX `json:"response"`
	}

	// AssetPolicyInfoResponse represents response from `/asset_policy_info` endpoint.
	AssetPolicyInfoResponse struct {
		Response
		Data []AssetInfo `json:"response"`
	}

	// AssetMintTX holds specific mint tx hash and amount.
	AssetMintTX struct {
		TxHash   TxHash   `json:"tx_hash"`
		Quantity Lovelace `json:"quantity"`
	}

	// AssetHistory holds given asset mint/burn tx's.
	AssetHistory struct {
		PolicyID   PolicyID      `json:"policy_id"`
		AssetName  AssetName     `json:"asset_name"`
		MintingTXs []AssetMintTX `json:"minting_txs"`
	}

	// AssetHistoryResponse represents response from `/asset_history` endpoint.
	AssetHistoryResponse struct {
		Response
		Data *AssetHistory `json:"response"`
	}
)

// GetAssetList returns the list of all native assets (paginated).
func (c *Client) GetAssetList(
	ctx context.Context,
	opts *RequestOptions,
) (res *AssetListResponse, err error) {
	res = &AssetListResponse{}
	rsp, err := c.request(ctx, &res.Response, "GET", "/asset_list", nil, opts)
	if err != nil {
		return
	}
	err = ReadAndUnmarshalResponse(rsp, &res.Response, &res.Data)
	return
}

// GetAssetAddressList returns the list of all addresses holding a given asset.
func (c *Client) GetAssetAddressList(
	ctx context.Context,
	policy PolicyID,
	name AssetName,
	opts *RequestOptions,
) (res *AssetAddressListResponse, err error) {
	res = &AssetAddressListResponse{}

	if opts == nil {
		opts = c.NewRequestOptions()
	}
	opts.QuerySet("_asset_policy", policy.String())
	opts.QuerySet("_asset_name", name.String())

	rsp, err := c.request(ctx, &res.Response, "GET", "/asset_address_list", nil, opts)
	if err != nil {
		return
	}
	err = ReadAndUnmarshalResponse(rsp, &res.Response, &res.Data)

	return
}

// GetAssetInfo returns the information of an asset including
// first minting & token registry metadata.
func (c *Client) GetAssetInfo(
	ctx context.Context,
	policy PolicyID,
	name AssetName,
	opts *RequestOptions,
) (res *AssetInfoResponse, err error) {
	res = &AssetInfoResponse{}

	if opts == nil {
		opts = c.NewRequestOptions()
	}
	opts.QuerySet("_asset_policy", policy.String())
	opts.QuerySet("_asset_name", name.String())

	rsp, err := c.request(ctx, &res.Response, "GET", "/asset_info", nil, opts)
	if err != nil {
		return
	}
	info := []AssetInfo{}
	err = ReadAndUnmarshalResponse(rsp, &res.Response, &info)

	if len(info) == 1 {
		res.Data = &info[0]
	}
	return
}

// GetAssetSummary returns the summary of an asset
// (total transactions exclude minting/total wallets
// include only wallets with asset balance).
func (c *Client) GetAssetSummary(
	ctx context.Context,
	policy PolicyID,
	name AssetName,
	opts *RequestOptions,
) (res *AssetSummaryResponse, err error) {
	res = &AssetSummaryResponse{}

	if opts == nil {
		opts = c.NewRequestOptions()
	}
	opts.QuerySet("_asset_policy", policy.String())
	opts.QuerySet("_asset_name", name.String())

	rsp, err := c.request(ctx, &res.Response, "GET", "/asset_summary", nil, opts)
	if err != nil {
		return
	}
	summary := []AssetSummary{}
	err = ReadAndUnmarshalResponse(rsp, &res.Response, &summary)

	if len(summary) == 1 {
		res.Data = &summary[0]
	}
	return
}

// GetAssetTxs returns the list of all asset transaction hashes (newest first).
func (c *Client) GetAssetTxs(
	ctx context.Context,
	policy PolicyID,
	name AssetName,
	opts *RequestOptions,
) (res *AssetTxsResponse, err error) {
	res = &AssetTxsResponse{}

	if opts == nil {
		opts = c.NewRequestOptions()
	}
	opts.QuerySet("_asset_policy", policy.String())
	opts.QuerySet("_asset_name", name.String())

	rsp, err := c.request(ctx, &res.Response, "GET", "/asset_txs", nil, opts)
	if err != nil {
		return
	}

	err = ReadAndUnmarshalResponse(rsp, &res.Response, &res.Data)

	return
}

// GetAssetPolicyInfo returns information for all assets under the same policy.
func (c *Client) GetAssetPolicyInfo(
	ctx context.Context,
	policy PolicyID,
	opts *RequestOptions,
) (res *AssetPolicyInfoResponse, err error) {
	res = &AssetPolicyInfoResponse{}

	if opts == nil {
		opts = c.NewRequestOptions()
	}
	opts.QuerySet("_asset_policy", policy.String())

	rsp, err := c.request(ctx, &res.Response, "GET", "/asset_policy_info", nil, opts)
	if err != nil {
		return
	}
	err = ReadAndUnmarshalResponse(rsp, &res.Response, &res.Data)

	return
}

// GetAssetHistory returns mint/burn history of an asset.
func (c *Client) GetAssetHistory(
	ctx context.Context,
	policy PolicyID,
	name AssetName,
	opts *RequestOptions,
) (res *AssetHistoryResponse, err error) {
	res = &AssetHistoryResponse{}
	if opts == nil {
		opts = c.NewRequestOptions()
	}
	opts.QuerySet("_asset_policy", policy.String())
	opts.QuerySet("_asset_name", name.String())

	rsp, err := c.request(ctx, &res.Response, "GET", "/asset_history", nil, opts)
	if err != nil {
		return
	}
	info := []AssetHistory{}
	err = ReadAndUnmarshalResponse(rsp, &res.Response, &info)

	if len(info) == 1 {
		res.Data = &info[0]
	}
	return
}
