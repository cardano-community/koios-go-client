// SPDX-License-Identifier: Apache-2.0
//
// Copyright © 2022 The Cardano Community Authors

package koios_test

import (
	"context"
	"encoding/json"
	"testing"

	koios "github.com/cardano-community/koios-go-client/v5"
)

const (
	mockPolicy    = koios.PolicyID("d3501d9531fcc25e3ca4b6429318c2cc374dbdbcf5e99c1c1e5da1ff")
	mockAssetName = koios.AssetName("444f4e545350414d")
)

func TestMockGetAssets(t *testing.T) {
	client, capture := mockClient(t, "asset_list")
	res, err := client.GetAssets(context.Background(), mockPolicy, mockAssetName, nil)
	mockNoErr(t, err)
	mockExpect(t, capture, "GET", "/asset_list", map[string]string{
		"_asset_policy": string(mockPolicy),
		"_asset_name":   string(mockAssetName),
	}, "")
	mockLen(t, "data", res.Data, 1)
	mockEq(t, "asset_name_ascii", res.Data[0].AssetNameASCII, "DONTSPAM")
	mockEq(t, "fingerprint", res.Data[0].Fingerprint, koios.AssetFingerprint("asset1ua6pz3yd5mdka946z8jw2fld3f8d0mmxt75gv9"))

	client, capture = mockClient(t, "asset_list")
	_, err = client.GetAssets(context.Background(), mockPolicy, "", nil)
	mockNoErr(t, err)
	mockExpect(t, capture, "GET", "/asset_list", map[string]string{"_asset_policy": string(mockPolicy)}, "")
}

func TestMockGetAssetTokenRegistry(t *testing.T) {
	client, capture := mockClient(t, "asset_token_registry")
	res, err := client.GetAssetTokenRegistry(context.Background(), nil)
	mockNoErr(t, err)
	mockExpect(t, capture, "GET", "/asset_token_registry", nil, "")
	mockLen(t, "data", res.Data, 1)
	mockEq(t, "policy_id", res.Data[0].PolicyID, mockPolicy)
	mockEq(t, "asset_name_ascii", res.Data[0].AssetNameASCII, "DONTSPAM")
	mockEq(t, "ticker", res.Data[0].Ticker, "MOB")
	mockEq(t, "url", res.Data[0].URL, "https://www.rackmob.com/")
}

func TestMockGetAssetInfo(t *testing.T) {
	client, capture := mockClient(t, "asset_info")
	res, err := client.GetAssetInfo(context.Background(), []koios.Asset{{PolicyID: mockPolicy}}, nil)
	mockNoErr(t, err)
	mockExpect(t, capture, "POST", "/asset_info", nil, `{"_asset_list": [["`+string(mockPolicy)+`", ""]]}`)
	mockLen(t, "data", res.Data, 1)
	if res.Data[0].TokenRegistryMetadata == nil {
		t.Fatal("token_registry_metadata is nil")
	}
	mockEq(t, "token_registry_metadata.name", res.Data[0].TokenRegistryMetadata.Name, "Rackmob")
	mockEq(t, "token_registry_metadata.ticker", res.Data[0].TokenRegistryMetadata.Ticker, "MOB")
	mockEq(t, "mint_cnt", res.Data[0].MintCnt, 1)
	mockEq(t, "burn_cnt", res.Data[0].BurnCnt, 5)
}

func TestMockGetAssetSummary(t *testing.T) {
	client, capture := mockClient(t, "asset_summary")
	res, err := client.GetAssetSummary(context.Background(), mockPolicy, mockAssetName, nil)
	mockNoErr(t, err)
	mockExpect(t, capture, "GET", "/asset_summary", map[string]string{
		"_asset_policy": string(mockPolicy),
		"_asset_name":   string(mockAssetName),
	}, "")
	if res.Data == nil {
		t.Fatal("data is nil")
	}
	mockEq(t, "policy_id", res.Data.PolicyID, mockPolicy)
	mockEq(t, "total_transactions", res.Data.TotalTransactions, uint64(89416))
	mockEq(t, "staked_wallets", res.Data.StakedWallets, uint64(548))
	mockEq(t, "unstaked_addresses", res.Data.UnstakedAddresses, uint64(245))
	mockEq(t, "addresses", res.Data.Addresses, uint64(812))
}

func TestMockGetAssetHistory(t *testing.T) {
	client, capture := mockClient(t, "asset_history")
	res, err := client.GetAssetHistory(context.Background(), mockPolicy, mockAssetName, nil)
	mockNoErr(t, err)
	mockExpect(t, capture, "GET", "/asset_history", map[string]string{
		"_asset_policy": string(mockPolicy),
		"_asset_name":   string(mockAssetName),
	}, "")
	if res.Data == nil {
		t.Fatal("data is nil")
	}
	mockLen(t, "minting_txs", res.Data.MintingTXs, 1)
	tx := res.Data.MintingTXs[0]
	mockEq(t, "quantity", tx.Quantity.String(), "-10")
	mockLen(t, "metadata", tx.Metadata, 1)
	var cip25 map[string]any
	if err := json.Unmarshal(tx.Metadata[0], &cip25); err != nil {
		t.Fatal(err)
	}
	if _, ok := cip25["721"]; !ok {
		t.Errorf("metadata %s has no 721 label", tx.Metadata[0])
	}
}

func TestMockGetAssetAddresses(t *testing.T) {
	client, capture := mockClient(t, "asset_addresses")
	res, err := client.GetAssetAddresses(context.Background(), mockPolicy, mockAssetName, nil)
	mockNoErr(t, err)
	mockExpect(t, capture, "GET", "/asset_addresses", map[string]string{
		"_asset_policy": string(mockPolicy),
		"_asset_name":   string(mockAssetName),
	}, "")
	mockLen(t, "data", res.Data, 1)
	mockEq(t, "payment_address", res.Data[0].PaymentAddress, mockAddr)
	mockEq(t, "stake_address", res.Data[0].StakeAddress, koios.Address("stake1u8yxtugdv63wxafy9d00nuz6hjyyp4qnggvc9a3vxh8yl0ckml2uz"))
	mockEq(t, "quantity", res.Data[0].Quantity.String(), "23")
}

func TestMockGetPolicyAssetAddresses(t *testing.T) {
	client, capture := mockClient(t, "policy_asset_addresses")
	res, err := client.GetPolicyAssetAddresses(context.Background(), mockPolicy, nil)
	mockNoErr(t, err)
	mockExpect(t, capture, "GET", "/policy_asset_addresses", map[string]string{"_asset_policy": string(mockPolicy)}, "")
	mockLen(t, "data", res.Data, 1)
	mockEq(t, "asset_name", res.Data[0].AssetName, mockAssetName)
	mockEq(t, "stake_address", res.Data[0].StakeAddress, koios.Address("stake1u8yxtugdv63wxafy9d00nuz6hjyyp4qnggvc9a3vxh8yl0ckml2uz"))
}

func TestMockGetPolicyAssetMints(t *testing.T) {
	client, capture := mockClient(t, "policy_asset_mints")
	res, err := client.GetPolicyAssetMints(context.Background(), mockPolicy, nil)
	mockNoErr(t, err)
	mockExpect(t, capture, "GET", "/policy_asset_mints", map[string]string{"_asset_policy": string(mockPolicy)}, "")
	mockLen(t, "data", res.Data, 1)
	mockEq(t, "mint_cnt", res.Data[0].MintCNT, uint(1))
	mockEq(t, "burn_cnt", res.Data[0].BurnCNT, uint(5))
	mockEq(t, "total_supply", res.Data[0].TotalSupply.String(), "35000")
	if len(res.Data[0].MintingTxMetadata) == 0 {
		t.Error("minting_tx_metadata is empty")
	}
}
