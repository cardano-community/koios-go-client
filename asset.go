// SPDX-License-Identifier: Apache-2.0
//
// Copyright © 2022 The Cardano Community Authors

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
		AssetName AssetName `json:"asset_name"`

		Fingerprint AssetFingerprint `json:"fingerprint"`

		// Asset Policy ID (hex).
		PolicyID PolicyID `json:"policy_id"`

		// Quantity
		// Input: asset balance on the selected input transaction.
		// Output: sum of assets for output UTxO.
		// Mint: sum of minted assets (negative on burn).
		Quantity decimal.Decimal `json:"quantity,omitempty"`

		Decimals uint8 `json:"decimals,omitempty"`
	}

	// AssetListItem is an asset as returned by `/asset_list`.
	AssetListItem struct {
		// Asset Name (hex).
		AssetName AssetName `json:"asset_name"`

		// AssetNameASCII is the asset name in ASCII (only reported by `/asset_list`).
		AssetNameASCII string `json:"asset_name_ascii,omitempty"`

		Fingerprint AssetFingerprint `json:"fingerprint"`

		// Asset Policy ID (hex).
		PolicyID PolicyID `json:"policy_id"`
	}

	PolicyAssetListItem struct {
		// Asset Name (hex).
		AssetName AssetName `json:"asset_name"`

		Fingerprint AssetFingerprint `json:"fingerprint"`
		// TotalSupply of Asset
		TotalSupply decimal.Decimal `json:"total_supply"`

		Decimals uint8 `json:"decimals"`
	}

	// TokenRegistryMetadata is the asset metadata registered on the Cardano
	// Token Registry, as nested in `/asset_info` and `/policy_asset_info`.
	TokenRegistryMetadata struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Ticker      string `json:"ticker"`
		URL         string `json:"url"`
		// Logo is a PNG image file as a byte string.
		Logo     string `json:"logo"`
		Decimals int    `json:"decimals"`
	}

	// TokenRegistryAsset is an asset registered via the token registry on
	// GitHub, as returned by `/asset_token_registry`.
	TokenRegistryAsset struct {
		PolicyID       PolicyID  `json:"policy_id"`
		AssetName      AssetName `json:"asset_name"`
		AssetNameASCII string    `json:"asset_name_ascii"`
		Ticker         string    `json:"ticker"`
		Description    string    `json:"description"`
		URL            string    `json:"url"`
		Decimals       int       `json:"decimals"`
		// Logo is a PNG image file as a byte string.
		Logo string `json:"logo"`
	}

	// AssetSummary aggregated asset summary.
	AssetSummary struct {
		// Asset Policy ID (hex).
		PolicyID PolicyID `json:"policy_id"`

		// Asset Name (hex).
		AssetName AssetName `json:"asset_name"`

		Fingerprint AssetFingerprint `json:"fingerprint"`

		// Total number of registered wallets holding the given asset
		StakedWallets uint64 `json:"staked_wallets"`

		// Total number of transactions including the given asset
		TotalTransactions uint64 `json:"total_transactions"`

		// Total number of payment addresses (not belonging
		// to registered wallets) holding the given asset
		UnstakedAddresses uint64 `json:"unstaked_addresses"`

		// Total number of unique addresses holding the given asset
		Addresses uint64 `json:"addresses"`
	}

	// AssetInfo info about the asset.
	AssetInfo struct {
		// Asset Policy ID (hex).
		PolicyID PolicyID `json:"policy_id,omitempty"`
		// Asset Name (hex).
		AssetName AssetName `json:"asset_name"`
		// Asset Name (ASCII)
		AssetNameASCII string `json:"asset_name_ascii,omitempty"`
		// The CIP14 fingerprint of the asset
		Fingerprint AssetFingerprint `json:"fingerprint"`
		// MintingTxHash mint tx
		MintingTxHash TxHash `json:"minting_tx_hash,omitempty"`
		// TotalSupply of Asset
		TotalSupply decimal.Decimal `json:"total_supply"`
		// MintCnt count of mint transactions
		MintCnt int `json:"mint_cnt,omitempty"`
		// BurnCnt count of burn transactions
		BurnCnt int `json:"burn_cnt,omitempty"`
		// CreationTime of Asset
		CreationTime Timestamp `json:"creation_time,omitempty"`
		// MintingTxMetadata latest minting transaction metadata (aligns with CIP-25).
		MintingTxMetadata *json.RawMessage `json:"minting_tx_metadata,omitempty"`
		// Asset metadata registered on the Cardano Token Registry
		TokenRegistryMetadata *TokenRegistryMetadata `json:"token_registry_metadata,omitempty"`
		CIP68Metadata         *json.RawMessage       `json:"cip68_metadata,omitempty"`
	}

	// AssetListResponse represents response from `/asset_list` endpoint.
	AssetListResponse struct {
		Response
		Data []AssetListItem `json:"data"`
	}

	// AssetHolder payment addresses holding the given token (including balance).
	AssetHolder struct {
		// AssetName is only reported by `/policy_asset_addresses`.
		AssetName      AssetName `json:"asset_name,omitempty"`
		PaymentAddress Address   `json:"payment_address"`
		// StakeAddress is the staking address of the holder (nullable).
		StakeAddress Address `json:"stake_address"`
		// Quantity is only reported by `/asset_addresses` and `/policy_asset_addresses`.
		Quantity decimal.Decimal `json:"quantity"`
	}

	// AssetAddressListResponse represents response from `/policy_asset_addresses` endpoint.
	AssetAddressListResponse struct {
		Response
		Data []AssetHolder `json:"data"`
	}
	AssetNFTAddressResponse struct {
		Response
		Data *Address `json:"data"`
	}

	// AssetInfoResponse represents response from `/asset_info` endpoint.
	AssetInfoResponse struct {
		Data []AssetInfo `json:"data"`
		Response
	}

	// AssetSummaryResponse represents response from `/asset_summary` endpoint.
	AssetSummaryResponse struct {
		Response
		Data *AssetSummary `json:"data"`
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
		// Metadata is the array of transaction metadata (CIP-25) of the transaction.
		Metadata []json.RawMessage `json:"metadata,omitempty"`
	}

	// AssetHistory holds given asset mint/burn tx's.
	AssetHistory struct {
		AssetListItem
		MintingTXs []AssetMintTX `json:"minting_txs"`
	}

	// AssetHistoryResponse represents response from `/asset_history` endpoint.
	AssetHistoryResponse struct {
		Response
		Data *AssetHistory `json:"data"`
	}

	// AssetPolicyAssetListResponse represents response from `/policy_asset_list` endpoint.
	AssetPolicyAssetListResponse struct {
		Response
		Data []PolicyAssetListItem `json:"data"`
	}

	// AssetTokenRegistryResponse represents response from `/asset_token_registry` endpoint.
	AssetTokenRegistryResponse struct {
		Response
		Data []TokenRegistryAsset `json:"data"`
	}

	AssetUTxOsResponse struct {
		Response
		Data []UTxO `json:"data"`
	}

	PolicyAssetMintsResponse struct {
		Response
		Data []PolicyAssetMint `json:"data"`
	}

	PolicyAssetMint struct {
		AssetName      AssetName        `json:"asset_name"`
		AssetNameASCII string           `json:"asset_name_ascii"`
		Fingerprint    AssetFingerprint `json:"fingerprint"`
		MintingTxHash  TxHash           `json:"minting_tx_hash"`
		TotalSupply    decimal.Decimal  `json:"total_supply"`
		MintCNT        uint             `json:"mint_cnt"`
		BurnCNT        uint             `json:"burn_cnt"`
		CreationTime   Timestamp        `json:"creation_time"`
		// MintingTxMetadata is the latest minting transaction metadata (CIP-25).
		MintingTxMetadata json.RawMessage `json:"minting_tx_metadata"`
		Decimals          uint8           `json:"decimals"`
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

// GetAssets returns the list of native assets (paginated). The policy
// and optional asset name filters bypass the cost of filtering the output;
// an empty name returns all assets of the policy.
func (c *Client) GetAssets(
	ctx context.Context,
	policy PolicyID,
	name AssetName,
	opts *RequestOptions,
) (res *AssetListResponse, err error) {
	res = &AssetListResponse{}
	if len(policy) == 0 {
		err = ErrNoPolicyID
		res.applyError(nil, err)
		return
	}
	if opts == nil {
		opts = c.NewRequestOptions()
	}
	opts.QuerySet("_asset_policy", policy.String())
	if len(name) > 0 {
		opts.QuerySet("_asset_name", name.String())
	}
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
	if len(assetName) > 0 {
		opts.QuerySet("_asset_name", assetName.String())
	}

	rsp, err := c.request(ctx, &res.Response, "GET", "/asset_addresses", nil, opts)
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
	assets []Asset,
	opts *RequestOptions,
) (res *AssetInfoResponse, err error) {
	res = &AssetInfoResponse{}

	if opts == nil {
		opts = c.NewRequestOptions()
	}

	if len(assets) == 0 {
		return nil, fmt.Errorf("%w: atleast one asset must be provided", ErrAsset)
	}

	var payload = struct {
		Assets [][]string `json:"_asset_list"`
	}{}

	for _, asset := range assets {
		if asset.PolicyID == "" {
			return nil, fmt.Errorf("%w: policy_id must be provided", ErrAsset)
		}
		payload.Assets = append(payload.Assets, []string{asset.PolicyID.String(), asset.AssetName.String()})
	}

	rpipe := payloadReader(payload)

	rsp, err := c.request(ctx, &res.Response, "POST", "/asset_info", rpipe, opts)
	if err != nil {
		return
	}
	err = ReadAndUnmarshalResponse(rsp, &res.Response, &res.Data)
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
	if len(name) > 0 {
		opts.QuerySet("_asset_name", name.String())
	}

	rsp, err := c.request(ctx, &res.Response, "GET", "/asset_summary", nil, opts)
	if err != nil {
		return
	}

	var assetSummary []AssetSummary
	err = ReadAndUnmarshalResponse(rsp, &res.Response, &assetSummary)
	if len(assetSummary) == 1 {
		data := assetSummary[0]
		res.Data = &data
	}

	return
}

// GetAssetTxs returns the list of all asset transaction hashes (newest first).
func (c *Client) GetAssetTxs(
	ctx context.Context,
	policy PolicyID,
	name AssetName,
	afterBlockHeight uint,
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
func (c *Client) GetPolicyAssetInfo(
	ctx context.Context,
	policy PolicyID,
	opts *RequestOptions,
) (res *AssetPolicyInfoResponse, err error) {
	res = &AssetPolicyInfoResponse{}

	if opts == nil {
		opts = c.NewRequestOptions()
	}
	opts.QuerySet("_asset_policy", policy.String())

	rsp, err := c.request(ctx, &res.Response, "GET", "/policy_asset_info", nil, opts)
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

func (c *Client) GetPolicyAssetList(
	ctx context.Context,
	policy PolicyID,
	opts *RequestOptions,
) (res *AssetPolicyAssetListResponse, err error) {
	res = &AssetPolicyAssetListResponse{}

	if opts == nil {
		opts = c.NewRequestOptions()
	}
	opts.QuerySet("_asset_policy", policy.String())

	rsp, err := c.request(ctx, &res.Response, "GET", "/policy_asset_list", nil, opts)
	if err != nil {
		return
	}

	err = ReadAndUnmarshalResponse(rsp, &res.Response, &res.Data)
	return
}

func (c *Client) GetAssetTokenRegistry(
	ctx context.Context,
	opts *RequestOptions,
) (res *AssetTokenRegistryResponse, err error) {
	res = &AssetTokenRegistryResponse{}

	rsp, err := c.request(ctx, &res.Response, "GET", "/asset_token_registry", nil, opts)
	if err != nil {
		return
	}
	err = ReadAndUnmarshalResponse(rsp, &res.Response, &res.Data)
	return
}

func (c *Client) GetAssetUTxOs(
	ctx context.Context,
	assets []Asset,
	opts *RequestOptions,
) (res *AssetUTxOsResponse, err error) {
	res = &AssetUTxOsResponse{}

	if opts == nil {
		opts = c.NewRequestOptions()
	}

	var payload = struct {
		Assets [][]string `json:"_asset_list"`
	}{}

	for _, asset := range assets {
		if asset.PolicyID == "" {
			return nil, fmt.Errorf("%w: policy_id must be provided", ErrAsset)
		}
		payload.Assets = append(payload.Assets, []string{asset.PolicyID.String(), asset.AssetName.String()})
	}

	rpipe := payloadReader(payload)

	rsp, err := c.request(ctx, &res.Response, "POST", "/asset_utxos", rpipe, opts)
	if err != nil {
		return
	}
	err = ReadAndUnmarshalResponse(rsp, &res.Response, &res.Data)
	return
}

func (c *Client) GetAssetNftAddress(
	ctx context.Context,
	policy PolicyID,
	name AssetName,
	opts *RequestOptions,
) (res *AssetNFTAddressResponse, err error) {
	res = &AssetNFTAddressResponse{}

	if opts == nil {
		opts = c.NewRequestOptions()
	}
	opts.QuerySet("_asset_policy", policy.String())
	if len(name) > 0 {
		opts.QuerySet("_asset_name", name.String())
	}

	rsp, err := c.request(ctx, &res.Response, "GET", "/asset_nft_address", nil, opts)
	if err != nil {
		return
	}

	var holders []AssetHolder

	err = ReadAndUnmarshalResponse(rsp, &res.Response, &holders)
	if len(holders) > 0 {
		res.Data = &holders[0].PaymentAddress
	}
	return
}

func (c *Client) GetPolicyAssetAddresses(
	ctx context.Context,
	policy PolicyID,
	opts *RequestOptions,
) (res *AssetAddressListResponse, err error) {
	res = &AssetAddressListResponse{}

	if opts == nil {
		opts = c.NewRequestOptions()
	}
	opts.QuerySet("_asset_policy", policy.String())

	rsp, err := c.request(ctx, &res.Response, "GET", "/policy_asset_addresses", nil, opts)
	if err != nil {
		return
	}
	err = ReadAndUnmarshalResponse(rsp, &res.Response, &res.Data)

	return
}

func (c *Client) GetPolicyAssetMints(
	ctx context.Context,
	policy PolicyID,
	opts *RequestOptions,
) (res *PolicyAssetMintsResponse, err error) {
	res = &PolicyAssetMintsResponse{}

	if opts == nil {
		opts = c.NewRequestOptions()
	}
	opts.QuerySet("_asset_policy", policy.String())

	rsp, err := c.request(ctx, &res.Response, "GET", "/policy_asset_mints", nil, opts)
	if err != nil {
		return
	}

	err = ReadAndUnmarshalResponse(rsp, &res.Response, &res.Data)
	return
}
