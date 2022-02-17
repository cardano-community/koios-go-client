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

package koios_test

import (
	"context"
	"encoding/json"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cardano-community/koios-go-client"
)

func TestNetworkTipEndpoint(t *testing.T) {
	expected := []koios.Tip{}

	spec := loadEndpointTestSpec(t, "endpoint_network_tip.json.gz", &expected)

	ts, api := setupTestServerAndClient(t, spec)

	defer ts.Close()

	res, err := api.GetTip(context.TODO())

	assert.NoError(t, err)
	testHeaders(t, spec, res.Response)

	assert.Len(t, expected, 1)
	assert.Equal(t, &expected[0], res.Data)
}

func TestNetworkGenesiEndpoint(t *testing.T) {
	expected := []koios.Genesis{}

	spec := loadEndpointTestSpec(t, "endpoint_network_genesis.json.gz", &expected)

	ts, api := setupTestServerAndClient(t, spec)

	defer ts.Close()

	res, err := api.GetGenesis(context.TODO())

	assert.NoError(t, err)
	testHeaders(t, spec, res.Response)

	assert.Len(t, expected, 1)
	assert.Equal(t, &expected[0], res.Data)
}

func TestNetworkTotalsEndpoint(t *testing.T) {
	expected := []koios.Totals{}

	spec := loadEndpointTestSpec(t, "endpoint_network_totals.json.gz", &expected)

	ts, api := setupTestServerAndClient(t, spec)

	defer ts.Close()

	epochNo, err := strconv.ParseUint(spec.Request.Query.Get("_epoch_no"), 10, 64)
	assert.NoError(t, err)
	epoch := koios.EpochNo(epochNo)

	res, err := api.GetTotals(context.TODO(), &epoch)
	assert.NoError(t, err)
	testHeaders(t, spec, res.Response)
	assert.Len(t, expected, 1)
	assert.Equal(t, expected[0], res.Data[0])

	// test data without epoch
	res2, err := api.GetTotals(context.TODO(), nil)
	assert.NoError(t, err)
	testHeaders(t, spec, res2.Response)
	assert.Len(t, expected, 1)
	assert.Equal(t, expected[0], res2.Data[0])
}

func TestEpochInfoEndpoint(t *testing.T) {
	expected := []koios.EpochInfo{}

	spec := loadEndpointTestSpec(t, "endpoint_epoch_info.json.gz", &expected)

	ts, api := setupTestServerAndClient(t, spec)

	defer ts.Close()

	epochNo, err := strconv.ParseUint(spec.Request.Query.Get("_epoch_no"), 10, 64)
	assert.NoError(t, err)
	epoch := koios.EpochNo(epochNo)

	res, err := api.GetEpochInfo(context.TODO(), &epoch)

	assert.NoError(t, err)
	testHeaders(t, spec, res.Response)

	assert.Len(t, expected, 1)
	assert.Equal(t, expected[0], res.Data[0])
}

func TestEpochParamsEndpoint(t *testing.T) {
	expected := []koios.EpochParams{}

	spec := loadEndpointTestSpec(t, "endpoint_epoch_params.json.gz", &expected)

	ts, api := setupTestServerAndClient(t, spec)

	defer ts.Close()

	epochNo, err := strconv.ParseUint(spec.Request.Query.Get("_epoch_no"), 10, 64)
	assert.NoError(t, err)
	epoch := koios.EpochNo(epochNo)

	res, err := api.GetEpochParams(context.TODO(), &epoch)

	assert.NoError(t, err)
	testHeaders(t, spec, res.Response)

	assert.Len(t, expected, 1)
	assert.Equal(t, expected[0], res.Data[0])
}

func TestAccountListEndpoint(t *testing.T) {
	expected := []struct {
		StakeAddress koios.StakeAddress `json:"id"`
	}{}

	spec := loadEndpointTestSpec(t, "endpoint_account_list.json.gz", &expected)

	ts, api := setupTestServerAndClient(t, spec)

	defer ts.Close()

	res, err := api.GetAccountList(context.TODO())
	assert.NoError(t, err)

	testHeaders(t, spec, res.Response)

	for _, e := range expected {
		assert.Contains(t, res.Data, e.StakeAddress)
	}
}

func TestAccountInfoEndpoint(t *testing.T) {
	expected := []koios.AccountInfo{}

	spec := loadEndpointTestSpec(t, "endpoint_account_info.json.gz", &expected)

	ts, api := setupTestServerAndClient(t, spec)

	defer ts.Close()

	res, err := api.GetAccountInfo(context.TODO(), koios.Address(spec.Request.Query.Get("_address")))

	assert.NoError(t, err)
	testHeaders(t, spec, res.Response)

	assert.Len(t, expected, 1)
	assert.Equal(t, &expected[0], res.Data)
}

func TestAccountRewardsEndpoint(t *testing.T) {
	expected := []koios.AccountRewards{}

	spec := loadEndpointTestSpec(t, "endpoint_account_rewards.json.gz", &expected)

	ts, api := setupTestServerAndClient(t, spec)

	defer ts.Close()

	epochNo, err := strconv.ParseUint(spec.Request.Query.Get("_epoch_no"), 10, 64)
	assert.NoError(t, err)
	epoch := koios.EpochNo(epochNo)

	res, err := api.GetAccountRewards(context.TODO(), koios.StakeAddress(spec.Request.Query.Get("_address")), &epoch)

	assert.NoError(t, err)
	testHeaders(t, spec, res.Response)

	assert.Len(t, expected, 1)
	assert.Equal(t, expected, res.Data)
}

func TestAccountUpdatesEndpoint(t *testing.T) {
	expected := []koios.AccountAction{}

	spec := loadEndpointTestSpec(t, "endpoint_account_updates.json.gz", &expected)

	ts, api := setupTestServerAndClient(t, spec)

	defer ts.Close()

	res, err := api.GetAccountUpdates(context.TODO(), koios.StakeAddress(spec.Request.Query.Get("_address")))

	assert.NoError(t, err)
	testHeaders(t, spec, res.Response)

	assert.Equal(t, expected, res.Data)
}

func TestAccountAddressesEndpoint(t *testing.T) {
	expected := []struct {
		Address koios.Address `json:"address"`
	}{}

	spec := loadEndpointTestSpec(t, "endpoint_account_addresses.json.gz", &expected)

	ts, api := setupTestServerAndClient(t, spec)

	defer ts.Close()

	res, err := api.GetAccountAddresses(context.TODO(), koios.StakeAddress(spec.Request.Query.Get("_address")))

	assert.NoError(t, err)
	testHeaders(t, spec, res.Response)

	for _, e := range expected {
		assert.Contains(t, res.Data, e.Address)
	}
}
func TestAccountAssetsEndpoint(t *testing.T) {
	expected := []koios.AccountAsset{}

	spec := loadEndpointTestSpec(t, "endpoint_account_assets.json.gz", &expected)

	ts, api := setupTestServerAndClient(t, spec)

	defer ts.Close()

	res, err := api.GetAccountAssets(context.TODO(), koios.StakeAddress(spec.Request.Query.Get("_address")))

	assert.NoError(t, err)
	testHeaders(t, spec, res.Response)

	assert.Equal(t, expected, res.Data)
}

func TestAccountHistoryEndpoint(t *testing.T) {
	expected := []koios.AccountHistoryEntry{}

	spec := loadEndpointTestSpec(t, "endpoint_account_history.json.gz", &expected)

	ts, api := setupTestServerAndClient(t, spec)

	defer ts.Close()

	res, err := api.GetAccountHistory(context.TODO(), koios.StakeAddress(spec.Request.Query.Get("_address")))

	assert.NoError(t, err)
	testHeaders(t, spec, res.Response)

	assert.Equal(t, expected, res.Data)
}

func TestGetAddressInfoEndpoint(t *testing.T) {
	expected := []koios.AddressInfo{}

	spec := loadEndpointTestSpec(t, "endpoint_address_info.json.gz", &expected)

	ts, api := setupTestServerAndClient(t, spec)

	defer ts.Close()

	res, err := api.GetAddressInfo(context.TODO(), koios.Address(spec.Request.Query.Get("_address")))

	assert.NoError(t, err)
	testHeaders(t, spec, res.Response)

	assert.Equal(t, &expected[0], res.Data)
}

func TestGetAddressTxsEndpoint(t *testing.T) {
	expected := []struct {
		TxHash koios.TxHash `json:"tx_hash"`
	}{}

	spec := loadEndpointTestSpec(t, "endpoint_address_txs.json.gz", &expected)

	ts, api := setupTestServerAndClient(t, spec)

	defer ts.Close()

	var payload = struct {
		Adresses         []koios.Address `json:"_addresses"`
		AfterBlockHeight uint64          `json:"_after_block_height,omitempty"`
	}{}
	err := json.Unmarshal(spec.Request.Body, &payload)
	assert.NoError(t, err)

	res, err := api.GetAddressTxs(context.TODO(), payload.Adresses, payload.AfterBlockHeight)

	assert.NoError(t, err)
	testHeaders(t, spec, res.Response)

	for _, e := range expected {
		assert.Contains(t, res.Data, e.TxHash)
	}
}

func TestGetAddressAssetsEndpoint(t *testing.T) {
	expected := []koios.AddressAsset{}

	spec := loadEndpointTestSpec(t, "endpoint_address_assets.json.gz", &expected)

	ts, api := setupTestServerAndClient(t, spec)

	defer ts.Close()

	res, err := api.GetAddressAssets(context.TODO(), koios.Address(spec.Request.Query.Get("_address")))

	assert.NoError(t, err)
	testHeaders(t, spec, res.Response)

	assert.Equal(t, expected, res.Data)
}

func TestGetCredentialTxsEndpoint(t *testing.T) {
	expected := []struct {
		TxHash koios.TxHash `json:"tx_hash"`
	}{}

	spec := loadEndpointTestSpec(t, "endpoint_credential_txs.json.gz", &expected)

	ts, api := setupTestServerAndClient(t, spec)

	defer ts.Close()

	var payload = struct {
		Credentials      []koios.PaymentCredential `json:"_payment_credentials"`
		AfterBlockHeight uint64                    `json:"_after_block_height,omitempty"`
	}{}
	err := json.Unmarshal(spec.Request.Body, &payload)
	assert.NoError(t, err)

	res, err := api.GetCredentialTxs(context.TODO(), payload.Credentials, payload.AfterBlockHeight)

	assert.NoError(t, err)
	testHeaders(t, spec, res.Response)

	for _, e := range expected {
		assert.Contains(t, res.Data, e.TxHash)
	}
}

func TestAssetListEndpoint(t *testing.T) {
	expected := []koios.AssetListItem{}

	spec := loadEndpointTestSpec(t, "endpoint_asset_list.json.gz", &expected)

	ts, api := setupTestServerAndClient(t, spec)

	defer ts.Close()

	res, err := api.GetAssetList(context.TODO())

	assert.NoError(t, err)
	testHeaders(t, spec, res.Response)

	assert.Equal(t, expected, res.Data)
}

func TestGetAssetAddressListEndpoint(t *testing.T) {
	expected := []koios.AssetHolder{}

	spec := loadEndpointTestSpec(t, "endpoint_asset_address_list.json.gz", &expected)

	ts, api := setupTestServerAndClient(t, spec)

	defer ts.Close()

	res, err := api.GetAssetAddressList(
		context.TODO(),
		koios.PolicyID(spec.Request.Query.Get("_asset_policy")),
		koios.AssetName(spec.Request.Query.Get("_asset_name")),
	)

	assert.NoError(t, err)
	testHeaders(t, spec, res.Response)

	assert.Equal(t, expected, res.Data)
}

func TestGetAssetInfoEndpoint(t *testing.T) {
	expected := []koios.AssetInfo{}

	spec := loadEndpointTestSpec(t, "endpoint_asset_info.json.gz", &expected)

	ts, api := setupTestServerAndClient(t, spec)

	defer ts.Close()

	res, err := api.GetAssetInfo(
		context.TODO(),
		koios.PolicyID(spec.Request.Query.Get("_asset_policy")),
		koios.AssetName(spec.Request.Query.Get("_asset_name")),
	)

	assert.NoError(t, err)
	testHeaders(t, spec, res.Response)

	assert.Equal(t, &expected[0], res.Data)
}

func TestGetAssetSummaryEndpoint(t *testing.T) {
	expected := []koios.AssetSummary{}

	spec := loadEndpointTestSpec(t, "endpoint_asset_summary.json.gz", &expected)

	ts, api := setupTestServerAndClient(t, spec)

	defer ts.Close()

	res, err := api.GetAssetSummary(
		context.TODO(),
		koios.PolicyID(spec.Request.Query.Get("_asset_policy")),
		koios.AssetName(spec.Request.Query.Get("_asset_name")),
	)

	assert.NoError(t, err)
	testHeaders(t, spec, res.Response)

	assert.Equal(t, &expected[0], res.Data)
}

func TestGetAssetTxsEndpoint(t *testing.T) {
	expected := []koios.AssetTxs{}

	spec := loadEndpointTestSpec(t, "endpoint_asset_txs.json.gz", &expected)

	ts, api := setupTestServerAndClient(t, spec)

	defer ts.Close()

	res, err := api.GetAssetTxs(
		context.TODO(),
		koios.PolicyID(spec.Request.Query.Get("_asset_policy")),
		koios.AssetName(spec.Request.Query.Get("_asset_name")),
	)

	assert.NoError(t, err)
	testHeaders(t, spec, res.Response)

	assert.Equal(t, &expected[0], res.Data)
}

func TestGetBlockInfoEndpoint(t *testing.T) {
	expected := []koios.Block{}

	spec := loadEndpointTestSpec(t, "endpoint_block_info.json.gz", &expected)

	ts, api := setupTestServerAndClient(t, spec)

	defer ts.Close()

	res, err := api.GetBlockInfo(
		context.TODO(),
		koios.BlockHash(spec.Request.Query.Get("_block_hash")),
	)

	assert.NoError(t, err)
	testHeaders(t, spec, res.Response)

	assert.Equal(t, &expected[0], res.Data)
}

func TestGetBlockTxsEndpoint(t *testing.T) {
	expected := []struct {
		TxHash koios.TxHash `json:"tx_hash"`
	}{}

	spec := loadEndpointTestSpec(t, "endpoint_block_txs.json.gz", &expected)

	ts, api := setupTestServerAndClient(t, spec)

	defer ts.Close()

	res, err := api.GetBlockTxHashes(
		context.TODO(),
		koios.BlockHash(spec.Request.Query.Get("_block_hash")),
	)

	assert.NoError(t, err)
	testHeaders(t, spec, res.Response)

	for _, e := range expected {
		assert.Contains(t, res.Data, e.TxHash)
	}
}

func TestGetBlocksEndpoint(t *testing.T) {
	expected := []koios.Block{}

	spec := loadEndpointTestSpec(t, "endpoint_blocks.json.gz", &expected)

	ts, api := setupTestServerAndClient(t, spec)

	defer ts.Close()

	res, err := api.GetBlocks(context.TODO())

	assert.NoError(t, err)
	testHeaders(t, spec, res.Response)

	assert.Equal(t, expected, res.Data)
}

func TestGetPoolBlocksEndpoint(t *testing.T) {
	expected := []koios.PoolBlockInfo{}

	spec := loadEndpointTestSpec(t, "endpoint_pool_blocks.json.gz", &expected)

	ts, api := setupTestServerAndClient(t, spec)

	defer ts.Close()

	epochNo, err := strconv.ParseUint(spec.Request.Query.Get("_epoch_no"), 10, 64)
	assert.NoError(t, err)
	epoch := koios.EpochNo(epochNo)

	res, err := api.GetPoolBlocks(
		context.TODO(),
		koios.PoolID(spec.Request.Query.Get("_pool_bech32")),
		&epoch,
	)

	assert.NoError(t, err)
	testHeaders(t, spec, res.Response)

	assert.Equal(t, expected, res.Data)
}

func TestGetPoolDelegatorsEndpoint(t *testing.T) {
	expected := []koios.PoolDelegator{}

	spec := loadEndpointTestSpec(t, "endpoint_pool_delegators.json.gz", &expected)

	ts, api := setupTestServerAndClient(t, spec)

	defer ts.Close()

	epochNo, err := strconv.ParseUint(spec.Request.Query.Get("_epoch_no"), 10, 64)
	assert.NoError(t, err)
	epoch := koios.EpochNo(epochNo)

	res, err := api.GetPoolDelegators(
		context.TODO(),
		koios.PoolID(spec.Request.Query.Get("_pool_bech32")),
		&epoch,
	)

	assert.NoError(t, err)
	testHeaders(t, spec, res.Response)

	assert.Equal(t, expected, res.Data)
}

func TestGetPoolInfoEndpoint(t *testing.T) {
	expected := []koios.PoolInfo{}

	spec := loadEndpointTestSpec(t, "endpoint_pool_info.json.gz", &expected)
	ts, api := setupTestServerAndClient(t, spec)
	defer ts.Close()

	res, err := api.GetPoolInfo(
		context.TODO(),
		koios.PoolID(spec.Request.Query.Get("_pool_bech32")),
	)

	assert.NoError(t, err)
	testHeaders(t, spec, res.Response)

	assert.Equal(t, &expected[0], res.Data)
}

func TestGetPoolListEndpoint(t *testing.T) {
	expected := []koios.PoolListItem{}

	spec := loadEndpointTestSpec(t, "endpoint_pool_list.json.gz", &expected)

	ts, api := setupTestServerAndClient(t, spec)

	defer ts.Close()

	res, err := api.GetPoolList(context.TODO())

	assert.NoError(t, err)
	testHeaders(t, spec, res.Response)

	assert.Equal(t, expected, res.Data)
}

func TestGetPoolMetadataEndpoint(t *testing.T) {
	expected := []koios.PoolMetadata{}

	spec := loadEndpointTestSpec(t, "endpoint_pool_metadata.json.gz", &expected)

	ts, api := setupTestServerAndClient(t, spec)

	defer ts.Close()

	res, err := api.GetPoolMetadata(context.TODO())

	assert.NoError(t, err)
	testHeaders(t, spec, res.Response)

	assert.Equal(t, expected, res.Data)
}

func TestGetPoolRelaysEndpoint(t *testing.T) {
	expected := []koios.PoolRelays{}

	spec := loadEndpointTestSpec(t, "endpoint_pool_relays.json.gz", &expected)

	ts, api := setupTestServerAndClient(t, spec)

	defer ts.Close()

	res, err := api.GetPoolRelays(context.TODO())

	assert.NoError(t, err)
	testHeaders(t, spec, res.Response)

	assert.Equal(t, expected, res.Data)
}

func TestGetPoolUpdatesEndpoint(t *testing.T) {
	expected := []koios.PoolUpdateInfo{}

	spec := loadEndpointTestSpec(t, "endpoint_pool_updates.json.gz", &expected)
	ts, api := setupTestServerAndClient(t, spec)
	defer ts.Close()

	poolID := koios.PoolID(spec.Request.Query.Get("_pool_bech32"))
	res, err := api.GetPoolUpdates(
		context.TODO(),
		&poolID,
	)

	assert.NoError(t, err)
	testHeaders(t, spec, res.Response)

	assert.Equal(t, expected, res.Data)
}

func TestGetScriptListEndpoint(t *testing.T) {
	expected := []koios.ScriptListItem{}

	spec := loadEndpointTestSpec(t, "endpoint_script_list.json.gz", &expected)

	ts, api := setupTestServerAndClient(t, spec)

	defer ts.Close()

	res, err := api.GetScriptList(context.TODO())

	assert.NoError(t, err)
	testHeaders(t, spec, res.Response)

	assert.Equal(t, expected, res.Data)
}

func TestGetScriptRedeemersEndpoint(t *testing.T) {
	expected := []koios.ScriptRedeemers{}

	spec := loadEndpointTestSpec(t, "endpoint_script_redeemers.json.gz", &expected)
	ts, api := setupTestServerAndClient(t, spec)
	defer ts.Close()

	res, err := api.GetScriptRedeemers(
		context.TODO(),
		koios.ScriptHash(spec.Request.Query.Get("_script_hash")),
	)

	assert.NoError(t, err)
	testHeaders(t, spec, res.Response)

	assert.Equal(t, &expected[0], res.Data)
}

func TestGetTxInfoEndpoint(t *testing.T) {
	expected := []koios.TxInfo{}

	spec := loadEndpointTestSpec(t, "endpoint_tx_info.json.gz", &expected)

	ts, api := setupTestServerAndClient(t, spec)

	defer ts.Close()

	var payload = struct {
		TxHashes []koios.TxHash `json:"_tx_hashes"`
	}{}
	err := json.Unmarshal(spec.Request.Body, &payload)
	assert.NoError(t, err)

	// Valid
	res, err := api.GetTxInfo(context.TODO(), payload.TxHashes[0])
	assert.NoError(t, err)
	testHeaders(t, spec, res.Response)
	assert.Equal(t, &expected[0], res.Data)

	// Empty payload
	res2, err := api.GetTxInfo(context.TODO(), koios.TxHash(""))
	assert.ErrorIs(t, err, koios.ErrNoTxHash)
	assert.Nil(t, res2.Data)
	if assert.NotNil(t, res2.Error) {
		assert.Equal(t, koios.ErrNoTxHash.Error(), res2.Error.Message)
	}
}

func TestGetTxMetadataEndpoint(t *testing.T) {
	expected := []koios.TxMetadata{}

	spec := loadEndpointTestSpec(t, "endpoint_tx_metadata.json.gz", &expected)

	ts, api := setupTestServerAndClient(t, spec)

	defer ts.Close()

	var payload = struct {
		TxHashes []koios.TxHash `json:"_tx_hashes"`
	}{}
	err := json.Unmarshal(spec.Request.Body, &payload)
	assert.NoError(t, err)

	res, err := api.GetTxMetadata(context.TODO(), payload.TxHashes[0])

	assert.NoError(t, err)
	testHeaders(t, spec, res.Response)

	assert.Equal(t, &expected[0], res.Data)
}

func TestGetTxMetaLabelsEndpoint(t *testing.T) {
	expected := []koios.TxMetalabel{}

	spec := loadEndpointTestSpec(t, "endpoint_tx_metalabels.json.gz", &expected)

	ts, api := setupTestServerAndClient(t, spec)

	defer ts.Close()

	res, err := api.GetTxMetaLabels(context.TODO())

	assert.NoError(t, err)
	testHeaders(t, spec, res.Response)

	assert.Equal(t, expected, res.Data)
}

func TestGetTxStatusEndpoint(t *testing.T) {
	expected := []koios.TxStatus{}

	spec := loadEndpointTestSpec(t, "endpoint_tx_status.json.gz", &expected)

	ts, api := setupTestServerAndClient(t, spec)

	defer ts.Close()

	var payload = struct {
		TxHashes []koios.TxHash `json:"_tx_hashes"`
	}{}
	err := json.Unmarshal(spec.Request.Body, &payload)
	assert.NoError(t, err)

	res, err := api.GetTxStatus(context.TODO(), payload.TxHashes[0])

	assert.NoError(t, err)
	testHeaders(t, spec, res.Response)

	assert.Equal(t, &expected[0], res.Data)
}
func TestGetTxsUTxOsEndpoint(t *testing.T) {
	expected := []koios.UTxO{}

	spec := loadEndpointTestSpec(t, "endpoint_tx_utxos.json.gz", &expected)

	ts, api := setupTestServerAndClient(t, spec)

	defer ts.Close()

	var payload = struct {
		TxHashes []koios.TxHash `json:"_tx_hashes"`
	}{}
	err := json.Unmarshal(spec.Request.Body, &payload)
	assert.NoError(t, err)

	res, err := api.GetTxsUTxOs(context.TODO(), payload.TxHashes)

	assert.NoError(t, err)
	testHeaders(t, spec, res.Response)

	assert.Equal(t, expected, res.Data)
}
