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

	"github.com/cardano-community/koios-go-client/v3"
	"github.com/stretchr/testify/assert"
)

func TestEpochInfo(t *testing.T) {
	client, err := getLiveClient()
	if testIsLocal(t, err) {
		return
	}
	epochInfoTest(t, networkEpoch(), client)
}

func epochInfoTest(t TestingT, epoch koios.EpochNo, client *koios.Client) {
	res, err := client.GetEpochInfo(context.Background(), &epoch, nil)
	if !assert.NoError(t, err) {
		return
	}
	assertEqual(t, epoch, res.Data[0].EpochNo, "epoch_no")

	if res.Data[0].TxCount == 0 {
		githubActionWarning("/epoch_info", fmt.Sprintf("epoch(%d) tx count is 0", res.Data[0].EpochNo))
	} else {
		assertIsPositive(t, res.Data[0].Fees, "fees")
		assertIsPositive(t, res.Data[0].OutSum, "out_sum")
	}

	assertGreater(t, res.Data[0].BlkCount, 0, "blk_count")
	assertTimeNotZero(t, res.Data[0].StartTime, "start_time")
	assertTimeNotZero(t, res.Data[0].EndTime, "end_time")
	assertTimeNotZero(t, res.Data[0].FirstBlockTime, "first_block_time")
	assertTimeNotZero(t, res.Data[0].LastBlockTime, "last_block_time")
	assertIsPositive(t, res.Data[0].ActiveStake, "active_stake")
	assertIsPositive(t, res.Data[0].TotalRewards, "total_rewards")
	assertIsPositive(t, res.Data[0].AvgBlkReward, "avg_blk_reward")
}

func TestEpochParams(t *testing.T) {
	client, err := getLiveClient()
	if testIsLocal(t, err) {
		return
	}
	epochParamsTest(t, networkEpoch(), client)

}

func epochParamsTest(t TestingT, epoch koios.EpochNo, client *koios.Client) {
	res, err := client.GetEpochParams(context.Background(), &epoch, nil)
	if !assert.NoError(t, err) {
		return
	}

	// assertGreater(t, res.Data[0].Decentralisation, 0, "decentralisation")
	// assertNotEmpty(t, res.Data[0].ExtraEntropy, "extra_entropy")
	// assertGreater(t, res.Data[0].ProtocolMinor, 0, "protocol_minor")

	assertEqual(t, epoch, res.Data[0].EpochNo, "epoch_no")
	assertIsPositive(t, res.Data[0].MinFeeA, "min_fee_a")
	assertIsPositive(t, res.Data[0].MinFeeB, "min_fee_b")
	assertGreater(t, res.Data[0].MaxBlockSize, 0, "max_block_size")
	assertGreater(t, res.Data[0].MaxTxSize, 0, "max_tx_size")
	assertGreater(t, res.Data[0].MaxBhSize, 0, "max_bh_size")
	assertNotEmpty(t, res.Data[0].KeyDeposit, "key_deposit")
	assertNotEmpty(t, res.Data[0].PoolDeposit, "pool_deposit")
	assertGreater(t, res.Data[0].MaxEpoch, 0, "max_epoch")
	assertGreater(t, res.Data[0].OptimalPoolCount, 0, "optimal_pool_count")
	assertGreater(t, res.Data[0].Influence, 0, "influence")
	assertGreater(t, res.Data[0].MonetaryExpandRate, 0, "monetary_expand_rate")
	assertGreater(t, res.Data[0].TreasuryGrowthRate, 0, "treasury_growth_rate")
	assertGreater(t, res.Data[0].ProtocolMajor, 0, "protocol_major")
	assertIsPositive(t, res.Data[0].MinUtxoValue, "min_utxo_value")
	assertIsPositive(t, res.Data[0].MinPoolCost, "min_pool_cost")
	assertNotEmpty(t, res.Data[0].Nonce, "nonce")
	assertNotEmpty(t, res.Data[0].BlockHash.String(), "block_hash")
	assertNotEmpty(t, res.Data[0].CostModels, "cost_models")
	assertIsPositive(t, res.Data[0].PriceMem, "price_mem")
	assertIsPositive(t, res.Data[0].PriceStep, "price_step")
	assertGreater(t, res.Data[0].MaxTxExMem, 0, "max_tx_ex_mem")
	assertGreater(t, res.Data[0].MaxTxExSteps, 0, "max_tx_ex_steps")
	assertGreater(t, res.Data[0].MaxBlockExMem, 0, "max_block_ex_mem")
	assertGreater(t, res.Data[0].MaxBlockExSteps, 0, "max_block_ex_steps")
	assertGreater(t, res.Data[0].MaxValSize, 0, "max_val_size")
	assertGreater(t, res.Data[0].CollateralPercent, 0, "collateral_percent")
	assertGreater(t, res.Data[0].MaxCollateralInputs, 0, "max_collateral_inputs")
	assertIsPositive(t, res.Data[0].CoinsPerUtxoSize, "coins_per_utxo_size")
	assertIsPositive(t, res.Data[0].CoinsPerUtxoSize, "coins_per_utxo_size")
}

func TestEpochBlockProtocols(t *testing.T) {
	client, err := getLiveClient()
	if testIsLocal(t, err) {
		return
	}
	epochBlockProtocolsTest(t, networkEpoch(), client)

}

func epochBlockProtocolsTest(t TestingT, epoch koios.EpochNo, client *koios.Client) {
	res, err := client.GetEpochBlockProtocols(context.Background(), &epoch, nil)
	if !assert.NoError(t, err) {
		return
	}
	assertGreater(t, res.Data[0].ProtoMajor, 0, "proto_major")
	// assertGreater(t, res.Data[0].ProtoMinor, 0, "proto_minor")
	assertGreater(t, res.Data[0].Blocks, 0, "blocks")
}
