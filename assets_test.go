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
	"fmt"
	"testing"

	"github.com/cardano-community/koios-go-client/v2"
	"github.com/stretchr/testify/assert"
)

func TestAssets(t *testing.T) {
	client, err := getLiveClient()
	if testIsLocal(t, err) {
		return
	}
	assetsTest(t, client)
}

func assetsTest(t TestingT, client *koios.Client) {
	opts := client.NewRequestOptions()
	opts.SetPageSize(10)
	res, err := client.GetAssets(context.Background(), opts)
	if !assert.NoError(t, err) {
		return
	}
	assertEqual(t, 10, len(res.Data), "total blocks returned")

	for _, item := range res.Data {
		assertNotEmpty(t, item.PolicyID, "policy_id")
		assertGreater(t, len(item.Fingerprint), 0, item.PolicyID.String()+" fingerprint")
	}
}

func TestAssetAddresses(t *testing.T) {
	client, err := getLiveClient()
	if testIsLocal(t, err) {
		return
	}
	policyID, assetName, _, _ := networkPolicyAsset()
	assetAddressesTest(t, client, policyID, assetName)
}

func assetAddressesTest(t TestingT, client *koios.Client, policyID koios.PolicyID, assetName koios.AssetName) {
	res, err := client.GetAssetAddresses(context.Background(), policyID, assetName, nil)
	if !assert.NoError(t, err) {
		return
	}

	for _, holder := range res.Data {
		assertNotEmpty(t, holder.PaymentAddress, "payment_address")
		assertIsPositive(t, holder.Quantity, "quantity")
	}
}

func TestAssetInfo(t *testing.T) {
	client, err := getLiveClient()
	if testIsLocal(t, err) {
		return
	}
	policyID, assetName, _, _ := networkPolicyAsset()
	assetInfoTest(t, client, policyID, assetName)
}

func assetInfoTest(t TestingT, client *koios.Client, policyID koios.PolicyID, assetName koios.AssetName) {
	res, err := client.GetAssetInfo(context.Background(), policyID, assetName, nil)
	if !assert.NoError(t, err) {
		return
	}

	assertNotEmpty(t, res.Data.PolicyID, "policy_id")
	if len(res.Data.AssetName) > 0 {
		assertNotEmpty(t, res.Data.AssetNameASCII, "asset_name_ascii")
	}
	assertNotEmpty(t, res.Data.Fingerprint, "fingerprint")
	assertNotEmpty(t, res.Data.MintingTxHash, "minting_tx_hash")
	assertIsPositive(t, res.Data.TotalSupply, "total_supply")
	assertGreater(t, res.Data.MintCnt, 0, "mint_cnt")
	// assertGreater(t, res.Data.BurnCnt, 0, "burn_cnt")
	assertTimeNotZero(t, res.Data.CreationTime, "creation_time")
	// assertNotEmpty(t, res.Data.MintingTxMetadata, "minting_tx_metadata")
	// assertNotEmpty(t, res.Data.TokenRegistryMetadata, "token_registry_metadata")
}

func TestAssetHistory(t *testing.T) {
	client, err := getLiveClient()
	if testIsLocal(t, err) {
		return
	}
	policyID, assetName, _, _ := networkPolicyAsset()
	assetHistoryTest(t, client, policyID, assetName)
}

func assetHistoryTest(t TestingT, client *koios.Client, policyID koios.PolicyID, assetName koios.AssetName) {
	res, err := client.GetAssetHistory(context.Background(), policyID, assetName, nil)
	if !assert.NoError(t, err) {
		return
	}
	assertNotEmpty(t, res.Data.PolicyID, "policy_id")
	assertNotEmpty(t, res.Data.Fingerprint, "fingerprint")
	for _, minttx := range res.Data.MintingTXs {
		assertNotEmpty(t, minttx.TxHash, "tx_hash")
		assertNotEmpty(t, minttx.BlockTime, "block_time")
		assertCoinNotZero(t, minttx.Quantity, "quantity")
		assertTxMetadata(t, minttx.Metadata, fmt.Sprintf("policy[%s].metadata", policyID))
	}
}

func TestAssetPolicyInfo(t *testing.T) {
	client, err := getLiveClient()
	if testIsLocal(t, err) {
		return
	}
	policyID, _, _, _ := networkPolicyAsset()
	assetPolicyInfoTest(t, client, policyID)
}

func assetPolicyInfoTest(t TestingT, client *koios.Client, policyID koios.PolicyID) {
	res, err := client.GetAssetPolicyInfo(context.Background(), policyID, nil)
	if !assert.NoError(t, err) {
		return
	}

	for i, asset := range res.Data {
		label := fmt.Sprintf("policy[%s][%d]", policyID, i)
		assertNotEmpty(t, asset.PolicyID, label+".policy_id")
		if len(asset.AssetName) > 0 {
			assertNotEmpty(t, asset.AssetNameASCII, label+".asset_name_ascii")
		}
		assertNotEmpty(t, asset.Fingerprint, label+".fingerprint")
		// assertNotEmpty(t, asset.MintingTxHash, label+".minting_tx_hash")
		assertIsPositive(t, asset.TotalSupply, label+".total_supply")
		// assertGreater(t, asset.MintCnt, 0, label+".mint_cnt")
		// assertGreater(t, res.Data.BurnCnt, 0, "burn_cnt")
		assertTimeNotZero(t, asset.CreationTime, label+".creation_time")
		// assertNotEmpty(t, res.Data.MintingTxMetadata, "minting_tx_metadata")
		// assertNotEmpty(t, res.Data.TokenRegistryMetadata, "token_registry_metadata")
	}
}

func TestAssetSummary(t *testing.T) {
	client, err := getLiveClient()
	if testIsLocal(t, err) {
		return
	}
	policyID, asset, _, _ := networkPolicyAsset()
	testAssetSummaryTest(t, client, policyID, asset)
}

func testAssetSummaryTest(t TestingT, client *koios.Client, policyID koios.PolicyID, asset koios.AssetName) {
	res, err := client.GetAssetSummary(context.Background(), policyID, asset, nil)
	if !assert.NoError(t, err) {
		return
	}

	for _, summary := range res.Data {
		assertNotEmpty(t, summary.PolicyID, "policy_id")
		// assertNotEmpty(t, summary.AssetName, "asset_name")
		assertGreater(t, summary.TotalTransactions, 0, "total_transactions")
		if summary.StakedWallets == 0 {
			githubActionWarning("AssetSummary", "staked_wallets is 0")
		}
		// assertGreater(t, summary.UnstakedAddresses, 0, "unstaked_addresses")
	}
}

func TestAssetTxs(t *testing.T) {
	client, err := getLiveClient()
	if testIsLocal(t, err) {
		return
	}
	policyID, asset, abh, history := networkPolicyAsset()
	assetTxsTest(t, client, policyID, asset, abh, history)
}

func assetTxsTest(
	t TestingT,
	client *koios.Client,
	policyID koios.PolicyID,
	asset koios.AssetName,
	abh int,
	history bool,
) {
	res, err := client.GetAssetTxs(context.Background(), policyID, asset, abh, history, nil)
	if !assert.NoError(t, err) {
		return
	}

	for i, tx := range res.Data {
		assertNotEmpty(t, tx.TxHash, fmt.Sprintf("tx[%d].tx_hash", i))
		assertNotEmpty(t, tx.EpochNo, fmt.Sprintf("tx[%d].epoch_no", i))
		assertNotEmpty(t, tx.BlockHeight, fmt.Sprintf("tx[%d].block_height", i))
		assertTimeNotZero(t, tx.BlockTime, fmt.Sprintf("tx[%d].block_time", i))
	}
}
