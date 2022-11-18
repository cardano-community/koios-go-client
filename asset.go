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

	"github.com/shopspring/decimal"
)

// introduces breaking change since v1.3.0

type (
	// AssetName defines type for _asset_name.
	AssetName string

	// AssetFingerprint defines type for asset_fingerprint.
	// The CIP14 fingerprint of the asset,
	// This specification defines a user-facing asset fingerprint
	// as a bech32-encoded blake2b-160 digest of the concatenation
	// of the policy id and the asset name.
	AssetFingerprint string

	// Asset represents Cardano Asset.
	Asset struct {
		// Asset Name (hex).
		AssetName AssetName `json:"asset_name,omitempty"`

		Fingerprint AssetFingerprint `json:"fingerprint,omitempty"`

		// Asset Policy ID (hex).
		PolicyID PolicyID `json:"policy_id"`

		// Quantity
		// Input: asset balance on the selected input transaction.
		// Output: sum of assets for output UTxO.
		// Mint: sum of minted assets (negative on burn).
		Quantity decimal.Decimal `json:"quantity,omitempty"`
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
		Asset

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
		AssetName AssetName `json:"asset_name"`

		// Asset Name (ASCII)
		AssetNameASCII string `json:"asset_name_ascii"`

		// The CIP14 fingerprint of the asset
		Fingerprint AssetFingerprint `json:"fingerprint"`

		// MintingTxMetadata minting Tx JSON payload if it can be decoded as JSON
		// MintingTxMetadata *TxInfoMetadata `json:"minting_tx_metadata"`
		MintingTxMetadata *json.RawMessage `json:"minting_tx_metadata,omitempty"`

		// Asset metadata registered on the Cardano Token Registry
		TokenRegistryMetadata *TokenRegistryMetadata `json:"token_registry_metadata,omitempty"`

		// Asset Policy ID (hex).
		PolicyID PolicyID `json:"policy_id,omitempty"`

		// TotalSupply of Asset
		TotalSupply decimal.Decimal `json:"total_supply"`

		// CreationTime of Asset
		CreationTime Timestamp `json:"creation_time"`

		// MintCnt count of mint transactions
		MintCnt int `json:"mint_cnt"`

		// BurnCnt count of burn transactions
		BurnCnt int `json:"burn_cnt"`

		// MintingTxHash mint tx
		MintingTxHash TxHash `json:"minting_tx_hash"`
	}

	// AssetListResponse represents response from `/asset_list` endpoint.
	AssetListResponse struct {
		Response
		Data []Asset `json:"data"`
	}

	// AssetHolder payment addresses holding the given token (including balance).
	AssetHolder struct {
		PaymentAddress Address         `json:"payment_address"`
		Quantity       decimal.Decimal `json:"quantity"`
	}

	// AssetAddressListResponse represents response from `/asset_address_list` endpoint.
	AssetAddressListResponse struct {
		Response
		Data []AssetHolder `json:"data"`
	}

	// AssetInfoResponse represents response from `/asset_info` endpoint.
	AssetInfoResponse struct {
		Data *AssetInfo `json:"data"`
		Response
	}

	// AssetSummaryResponse represents response from `/asset_summary` endpoint.
	AssetSummaryResponse struct {
		Response
		Data []AssetSummary `json:"data"`
	}

	// AssetTxsResponse represents response from `/asset_txs` endpoint.
	AssetTxsResponse struct {
		Response
		Data []AddressTx `json:"data"`
	}

	// AssetPolicyInfoResponse represents response from `/asset_policy_info` endpoint.
	AssetPolicyInfoResponse struct {
		Response
		Data []AssetInfo `json:"data"`
	}

	// AssetMintTX holds specific mint tx hash and amount.
	AssetMintTX struct {
		TxHash    TxHash          `json:"tx_hash"`
		Quantity  decimal.Decimal `json:"quantity"`
		BlockTime Timestamp       `json:"block_time"`
		Metadata  TxMetadata      `json:"metadata,omitempty"`
	}

	// AssetHistory holds given asset mint/burn tx's.
	AssetHistory struct {
		Asset
		MintingTXs []AssetMintTX `json:"minting_txs"`
	}

	// AssetHistoryResponse represents response from `/asset_history` endpoint.
	AssetHistoryResponse struct {
		Response
		Data *AssetHistory `json:"data"`
	}
)

// String returns AssetName as string.
func (v AssetName) String() string {
	return string(v)
}

// String returns AssetFingerprint as string.
func (v AssetFingerprint) String() string {
	return string(v)
}

// GetAssetList returns the list of all native assets (paginated).
func (c *Client) GetAssets(
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
func (c *Client) GetAssetAddresses(
	ctx context.Context,
	policy PolicyID,
	assetName AssetName,
	opts *RequestOptions,
) (res *AssetAddressListResponse, err error) {
	res = &AssetAddressListResponse{}

	if opts == nil {
		opts = c.NewRequestOptions()
	}
	opts.QuerySet("_asset_policy", policy.String())
	opts.QuerySet("_asset_name", assetName.String())

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

	err = ReadAndUnmarshalResponse(rsp, &res.Response, &res.Data)

	return
}

// GetAssetTxs returns the list of all asset transaction hashes (newest first).
func (c *Client) GetAssetTxs(
	ctx context.Context,
	policy PolicyID,
	name AssetName,
	afterBlockHeight int,
	history bool,
	opts *RequestOptions,
) (res *AssetTxsResponse, err error) {
	res = &AssetTxsResponse{}

	if opts == nil {
		opts = c.NewRequestOptions()
	}
	opts.QuerySet("_asset_policy", policy.String())
	if len(name) > 0 {
		opts.QuerySet("_asset_name", name.String())
	}
	if afterBlockHeight > 0 {
		opts.QuerySet("_after_block_height", fmt.Sprint(afterBlockHeight))
	}
	if history {
		opts.QuerySet("_history", fmt.Sprint(history))
	}
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

	for i := range res.Data {
		res.Data[i].PolicyID = policy
	}
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
	if len(name) > 0 {
		opts.QuerySet("_asset_name", name.String())
	}

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
