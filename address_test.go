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

func TestAddressInfo(t *testing.T) {
	client, err := getClient()
	if !assert.NoError(t, err) {
		return
	}
	addressInfoTest(t, networkAddresses(), client)
}

func addressInfoTest(t TestingT, addrs []koios.Address, client *koios.Client) {
	res, err := client.GetAddressesInfo(context.Background(), addrs, nil)
	if !assert.NoError(t, err) {
		return
	}
	for i, addr := range res.Data {
		assertNotEmpty(t, addr.Address, fmt.Sprintf("info[%d].address", i))
		assertIsPositive(t, addr.Balance, fmt.Sprintf("status[%d].balance", i))
		assert.False(t, addr.ScriptAddress, fmt.Sprintf("status[%d].script_address", i))
		for i, utxo := range addr.UTxOs {
			assertUTxO(t, utxo, fmt.Sprintf("add[%s].utxo_set[%d]", addr.Address, i))
		}
	}
}

func TestAddressTxs(t *testing.T) {
	client, err := getClient()
	if !assert.NoError(t, err) {
		return
	}
	addressTxsTest(t, networkAddresses(), client)
}

func addressTxsTest(t TestingT, addrs []koios.Address, client *koios.Client) {
	res, err := client.GetAddressTxs(context.Background(), addrs, 0, nil)
	if !assert.NoError(t, err) {
		return
	}

	assertGreater(t, len(res.Data), 10, "expected transactions list")
	for i, tx := range res.Data {
		assertNotEmpty(t, tx.TxHash, fmt.Sprintf("tx[%d].tx_hash", i))
		assertNotEmpty(t, tx.EpochNo, fmt.Sprintf("tx[%d].epoch_no", i))
		assertNotEmpty(t, tx.BlockHeight, fmt.Sprintf("tx[%d].block_height", i))
		assertTimeNotZero(t, tx.BlockTime, fmt.Sprintf("tx[%d].block_time", i))
	}
}

func TestAddressAssets(t *testing.T) {
	client, err := getClient()
	if !assert.NoError(t, err) {
		return
	}
	addressAssetsTest(t, networkAddresses(), client)
}

func addressAssetsTest(t TestingT, addrs []koios.Address, client *koios.Client) {
	res, err := client.GetAddressesAssets(context.Background(), addrs, nil)
	if err != nil {
		if assert.ErrorIs(t, err, koios.ErrNoData) {
			githubActionWarning("AddressAssets", err.Error())
			return
		}
		assert.NoError(t, err)
		return
	}

	for _, addrcol := range res.Data {
		assertNotEmpty(t, addrcol.Address, "address")
		for i, col := range addrcol.Collections {
			label := fmt.Sprintf("address[%s].assets[%d]", addrcol.Address, i)
			assertNotEmpty(t, col.PolicyID, label+".ploicy_id")
			for j, asset := range col.Assets {
				label2 := fmt.Sprintf("%s.assets[%d]", label, j)
				if len(asset.AssetName) > 0 {
					assertNotEmpty(t, asset.AssetName, label2+"asset_name")
					assertNotEmpty(t, asset.AssetNameASCII, label2+"asset_name_ascii")
				}
				assertIsPositive(t, asset.Balance, label2+"balance")
			}
		}
	}
}

func TestCredentialTxs(t *testing.T) {
	client, err := getClient()
	if !assert.NoError(t, err) {
		return
	}
	credentialTxsTest(t, networkPaymentCredentials(), client)
}

func credentialTxsTest(t TestingT, creds []koios.PaymentCredential, client *koios.Client) {
	res, err := client.GetCredentialTxs(context.Background(), creds, 0, nil)
	if !assert.NoError(t, err) {
		return
	}

	assertGreater(t, len(res.Data), 2, "expected transactions list")
	for i, tx := range res.Data {
		assertNotEmpty(t, tx.TxHash, fmt.Sprintf("tx[%d].tx_hash", i))
		assertNotEmpty(t, tx.EpochNo, fmt.Sprintf("tx[%d].epoch_no", i))
		assertNotEmpty(t, tx.BlockHeight, fmt.Sprintf("tx[%d].block_height", i))
		assertTimeNotZero(t, tx.BlockTime, fmt.Sprintf("tx[%d].block_time", i))
	}
}
