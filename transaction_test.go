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

func TestTxInfo(t *testing.T) {
	client, err := getClient()
	if !assert.NoError(t, err) {
		return
	}
	txInfoTest(t, networkTxHashes(), client)
}

func txInfoTest(t TestingT, hashes []koios.TxHash, client *koios.Client) {
	res, err := client.GetTxsInfo(context.Background(), hashes, nil)
	if !assert.NoError(t, err) {
		return
	}

	for _, tx := range res.Data {
		assertNotEmpty(t, tx.TxHash, fmt.Sprintf("tx[%s].tx_hash", tx.TxHash))
		assertNotEmpty(t, tx.BlockHash, fmt.Sprintf("tx[%s].block_hash", tx.TxHash))
		assertGreater(t, tx.BlockHeight, 0, fmt.Sprintf("tx[%s].block_height", tx.TxHash))
		assertGreater(t, tx.EpochNo, koios.EpochNo(0), fmt.Sprintf("tx[%s].epoch_no", tx.TxHash))
		assertGreater(t, tx.EpochSlot, koios.Slot(0), fmt.Sprintf("tx[%s].epoch_slot", tx.TxHash))
		assertGreater(t, tx.AbsoluteSlot, koios.Slot(0), fmt.Sprintf("tx[%s].absolute_slot", tx.TxHash))
		assertTimeNotZero(t, tx.TxTimestamp, fmt.Sprintf("tx[%s].tx_timestamp", tx.TxHash))
		// assertGreater(t, tx.TxBlockIndex, 0, fmt.Sprintf("tx[%s].tx_block_index", tx.TxHash))
		assertGreater(t, tx.TxSize, 0, fmt.Sprintf("tx[%s].tx_size", tx.TxHash))
		assertIsPositive(t, tx.TotalOutput, fmt.Sprintf("tx[%s].total_output", tx.TxHash))
		assertIsPositive(t, tx.Fee, fmt.Sprintf("tx[%s].fee", tx.TxHash))

		// assertTimeNotZero(t, tx.InvalidBefore, "invalid_before")
		// assertTimeNotZero(t, tx.InvalidAfter, fmt.Sprintf("tx[%s].invalid_after", tx.TxHash))

		if !tx.Deposit.IsZero() {
			assertIsPositive(t, tx.Deposit, fmt.Sprintf("tx[%s].deposit", tx.TxHash))
		}
		// DATA
		if len(tx.CollateralInputs) > 0 {
			for _, utxo := range tx.CollateralInputs {
				assertUTxO(t, utxo, fmt.Sprintf("tx[%s].tx.collateral_inputs", tx.TxHash))
			}
		}

		if tx.CollateralOutput != nil {
			// assertNotEmpty(t, , "collateral_inputs")
			assertUTxO(t, *tx.CollateralOutput, fmt.Sprintf("tx[%s].tx.collateral_output", tx.TxHash))
		}

		if len(tx.ReferenceInputs) > 0 {
			for i, utxo := range tx.ReferenceInputs {
				assertUTxO(t, utxo, fmt.Sprintf("tx[%s].reference_inputs[%d]", tx.TxHash, i))
			}
		}

		if assertGreater(t, len(tx.Inputs), 0, "inputs") {
			for i, utxo := range tx.Inputs {
				assertUTxO(t, utxo, fmt.Sprintf("tx[%s].inputs[%d]", tx.TxHash, i))
			}
		}
		if assertGreater(t, len(tx.Outputs), 0, "outputs") {
			for i, utxo := range tx.Outputs {
				assertUTxO(t, utxo, fmt.Sprintf("tx[%s].outputs[%d]", tx.TxHash, i))
			}
		}

		if len(tx.Withdrawals) > 0 {
			for i, withdrawal := range tx.Withdrawals {
				assertNotEmpty(t, withdrawal.StakeAddress, fmt.Sprintf("tx[%s].withdrawals[%d].stake_addr", tx.TxHash, i))
				assertIsPositive(t, withdrawal.Amount, fmt.Sprintf("tx[%s].withdrawals[%d].amount", tx.TxHash, i))
			}
		}

		if len(tx.AssetsMinted) > 0 {
			for i, asset := range tx.AssetsMinted {
				assertAsset(t, asset, fmt.Sprintf("tx[%s].assets_minted[%d]", tx.TxHash, i))
			}
		}

		if len(tx.Metadata) > 0 {
			assertTxMetadata(t, tx.Metadata, fmt.Sprintf("tx[%s].metadata", tx.TxHash))
		}

		if len(tx.Certificates) > 0 {
			assertCertificates(t, tx.Certificates, fmt.Sprintf("tx[%s].certificates", tx.TxHash))
		}
		if len(tx.NativeScripts) > 0 {
			assertNativeScripts(t, tx.NativeScripts, fmt.Sprintf("tx[%s].native_scripts", tx.TxHash))
		}
		if len(tx.PlutusContracts) > 0 {
			assertPlutusContracts(t, tx.PlutusContracts, fmt.Sprintf("tx[%s].plutus_contracts", tx.TxHash))
		}
	}
}

func TestTxUTxO(t *testing.T) {
	client, err := getClient()
	if !assert.NoError(t, err) {
		return
	}
	txUTxOsTest(t, networkTxHashes(), client)
}

func txUTxOsTest(t TestingT, hashes []koios.TxHash, client *koios.Client) {
	res, err := client.GetTxsUTxOs(context.Background(), hashes, nil)
	if !assert.NoError(t, err) {
		return
	}
	for _, eutxo := range res.Data {
		assertEUTxO(t, eutxo, fmt.Sprintf("tx[%s]", eutxo.TxHash))
	}
}
