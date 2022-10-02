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

func TestPools(t *testing.T) {
	client, err := getClient()
	if !assert.NoError(t, err) {
		return
	}
	poolsTest(t, client)
}

func poolsTest(t TestingT, client *koios.Client) {
	opts := client.NewRequestOptions()
	opts.SetPageSize(10)
	res, err := client.GetPools(context.Background(), opts)
	if !assert.NoError(t, err) {
		return
	}
	assertEqual(t, 10, len(res.Data), "total blocks returned")

	for _, pool := range res.Data {
		assertNotEmpty(t, pool.PoolID, "pool_id_bech32")
		// assertNotEmpty(t, pool.Ticker, "ticker")
	}
}

func TestPoolSnapshot(t *testing.T) {
	client, err := getClient()
	if !assert.NoError(t, err) {
		return
	}
	poolID := networkPoolID()
	poolSnapshotTest(t, client, poolID)
}

func poolSnapshotTest(t TestingT, client *koios.Client, poolID koios.PoolID) {
	opts := client.NewRequestOptions()
	opts.SetPageSize(10)
	res, err := client.GetPoolSnapshot(context.Background(), poolID, opts)
	if !assert.NoError(t, err) {
		return
	}
	assertEqual(t, 3, len(res.Data), "snapshots returned")

	for i, pool := range res.Data {
		assertNotEmpty(t, pool.Snapshot, "snapshot")
		if i < 2 {
			assertNotEmpty(t, pool.Nonce, "nonce")
		}
		assertIsPositive(t, pool.PoolStake, "pool_stake")
		assertIsPositive(t, pool.ActiveStake, "active_stake")
		assertGreater(t, pool.EpochNo, 0, "epoch_no")
	}
}

func TestPoolInfo(t *testing.T) {
	client, err := getClient()
	if !assert.NoError(t, err) {
		return
	}
	poolID := networkPoolID()
	poolInfoTest(t, client, poolID)
}

func poolInfoTest(t TestingT, client *koios.Client, poolID koios.PoolID) {

	res, err := client.GetPoolInfo(context.Background(), poolID, nil)
	if !assert.NoError(t, err) {
		return
	}
	assertNotEmpty(t, res.Data.PoolID, "pool_id_bech32")
	assertNotEmpty(t, res.Data.PoolIdHex, "pool_id_hex")
	assertGreater(t, res.Data.ActiveEpoch, 0, "active_epoch_no")
	assertNotEmpty(t, res.Data.VrfKeyHash, "vrf_key_hash")
	assertGreater(t, res.Data.Margin, 0, "margin")
	assertIsPositive(t, res.Data.FixedCost, "fixed_cost")
	assertIsPositive(t, res.Data.Pledge, "pledge")
	assertNotEmpty(t, res.Data.RewardAddr, "reward_addr")
	for _, owner := range res.Data.Owners {
		assertNotEmpty(t, owner, "owners")
	}
	// assertNotEmpty(t, res.Data.Relays, "relays")
	assertNotEmpty(t, res.Data.MetaURL, "meta_url")
	assertNotEmpty(t, res.Data.MetaHash, "meta_hash")
	// assertNotEmpty(t, res.Data.MetaJSON, "meta_json")
	assertNotEmpty(t, res.Data.PoolStatus, "pool_status")
	// assertNotEmpty(t, res.Data.RetiringEpoch, "retiring_epoch")
	// assertNotEmpty(t, res.Data.OpCert, "op_cert")
	// assertNotEmpty(t, res.Data.OpCertCounter, "op_cert_counter")
	assertNotEmpty(t, res.Data.ActiveStake, "active_stake")
	// assertNotEmpty(t, res.Data.Sigma, "sigma")
	// assertNotEmpty(t, res.Data.BlockCount, "block_count")
	assertNotEmpty(t, res.Data.LivePledge, "live_pledge")
	assertNotEmpty(t, res.Data.LiveStake, "live_stake")
	assertNotEmpty(t, res.Data.LiveDelegators, "live_delegators")
	assertNotEmpty(t, res.Data.LiveSaturation, "live_saturation")
}

func TestPoolDelegators(t *testing.T) {
	client, err := getClient()
	if !assert.NoError(t, err) {
		return
	}
	poolID := networkPoolID()
	poolDelegatorsTest(t, client, poolID)
}

func poolDelegatorsTest(t TestingT, client *koios.Client, poolID koios.PoolID) {
	res, err := client.GetPoolDelegators(context.Background(), poolID, nil)
	if !assert.NoError(t, err) {
		return
	}

	for _, d := range res.Data {
		assertNotEmpty(t, d.StakeAddress, "stake_address")
		assertIsPositive(t, d.Amount, "amount")
		assertGreater(t, d.ActiveEpochNo, 0, "active_epoch_no")
		assertNotEmpty(t, d.LatestDelegationTxHash, "latest_delegation_tx_hash")
	}
}

func TestPoolDelegatorsHistory(t *testing.T) {
	client, err := getClient()
	if !assert.NoError(t, err) {
		return
	}
	poolID := networkPoolID()
	epoch := networkEpoch()
	poolDelegatorsHistoryTest(t, client, poolID, &epoch)
}

func poolDelegatorsHistoryTest(t TestingT, client *koios.Client, poolID koios.PoolID, epoch *koios.EpochNo) {
	res, err := client.GetPoolDelegatorsHistory(context.Background(), poolID, epoch, nil)
	if !assert.NoError(t, err) {
		return
	}

	for _, d := range res.Data {
		assertNotEmpty(t, d.StakeAddress, "stake_address")
		assertIsPositive(t, d.Amount, "amount")
		assertGreater(t, d.EpochNo, 0, "epoch_no")
	}
}

func TestPoolBlocks(t *testing.T) {
	client, err := getClient()
	if !assert.NoError(t, err) {
		return
	}
	poolID := networkPoolID()
	epoch := networkEpoch()
	poolBlocksTest(t, client, poolID, &epoch)
}

func poolBlocksTest(t TestingT, client *koios.Client, poolID koios.PoolID, epoch *koios.EpochNo) {
	res, err := client.GetPoolBlocks(context.Background(), poolID, epoch, nil)
	if !assert.NoError(t, err) {
		return
	}

	for _, d := range res.Data {
		assertNotEmpty(t, d.BlockHash, "block_hash")
		assertGreater(t, d.AbsSlot, 0, "abs_slot")
		assertGreater(t, d.EpochNo, 0, "epoch_no")
		assertGreater(t, d.BlockHeight, 0, "block_height")
		assertTimeNotZero(t, d.BlockTime, "block_time")
	}
}

func TestPoolHistory(t *testing.T) {
	client, err := getClient()
	if !assert.NoError(t, err) {
		return
	}
	poolID := networkPoolID()
	epoch := networkEpoch()
	poolHistoryTest(t, client, poolID, &epoch)
}

func poolHistoryTest(t TestingT, client *koios.Client, poolID koios.PoolID, epoch *koios.EpochNo) {
	res, err := client.GetPoolHistory(context.Background(), poolID, epoch, nil)
	if !assert.NoError(t, err) {
		return
	}
	for _, d := range res.Data {
		assertGreater(t, d.EpochNo, 0, "block_hash")
		assertIsPositive(t, d.ActiveStake, "active_stake")
		assertGreater(t, d.ActiveStakePCT, 0, "active_stake_pct")
		assertGreater(t, d.SaturationPCT, 0, "saturation_pct")
		assertGreater(t, d.BlockCNT, 0, "block_cnt")
		assertGreater(t, d.DelegatorCNT, 0, "delegator_cnt")
		assertGreater(t, d.Margin, 0, "margin")
		assertIsPositive(t, d.PoolFees, "pool_fees")
		assertIsPositive(t, d.FixedCost, "fixed_cost")
		// assertIsPositive(t, d.DelegRewards, "deleg_rewards")
		// assertGreater(t, d.EpochROS, 0, "epoch_ros")
	}
}

func TestPoolUpdates(t *testing.T) {
	client, err := getClient()
	if !assert.NoError(t, err) {
		return
	}
	poolID := networkPoolID()
	poolUpdatesTest(t, client, &poolID)
}

func poolUpdatesTest(t TestingT, client *koios.Client, poolID *koios.PoolID) {
	res, err := client.GetPoolUpdates(context.Background(), poolID, nil)
	if !assert.NoError(t, err) {
		return
	}
	for _, update := range res.Data {
		assertNotEmpty(t, update.TxHash, "tx_hash")
		assertTimeNotZero(t, update.BlockTime, "block_time")
		assertNotEmpty(t, update.PoolID, "pool_id_bech32")
		assertNotEmpty(t, update.PoolIDHex, "pool_id_hex")
		assertGreater(t, update.ActiveEpoch, 0, "active_epoch_no")
		assertNotEmpty(t, update.VrfKeyHash, "vrf_key_hash")
		assertGreater(t, update.Margin, 0, "margin")
		assertIsPositive(t, update.FixedCost, "fixed_cost")
		assertIsPositive(t, update.Pledge, "pledge")
		assertNotEmpty(t, update.RewardAddr, "reward_addr")
		for _, owner := range update.Owners {
			assertNotEmpty(t, owner, "owners")
		}
		for _, relay := range update.Relays {
			assertNotEmpty(t, relay.DNS, "relays.dns")
		}
		assertNotEmpty(t, update.MetaURL, "meta_url")
		assertNotEmpty(t, update.MetaHash, "meta_hash")
		assertNotEmpty(t, update.PoolStatus, "pool_status")
		// assertGreater(t, update.RetiringEpoch, 0, "retiring_epoch")
	}
}

func TestPoolRelays(t *testing.T) {
	client, err := getClient()
	if !assert.NoError(t, err) {
		return
	}
	poolRelaysTest(t, client)
}

func poolRelaysTest(t TestingT, client *koios.Client) {
	opts := client.NewRequestOptions()
	opts.SetPageSize(10)
	res, err := client.GetPoolRelays(context.Background(), opts)
	if !assert.NoError(t, err) {
		return
	}
	assertEqual(t, 10, len(res.Data), "ppol relays returned")
	for _, pool := range res.Data {
		assertNotEmpty(t, pool.PoolID, "pool_id_bech32")
		for i, relay := range pool.Relays {
			if len(relay.Ipv4) == 0 {
				assertNotEmpty(t, relay.DNS, fmt.Sprintf("pool[%s].relays[%d]", pool.PoolID, i))
			}
		}
	}
}

func TestPoolMetadata(t *testing.T) {
	client, err := getClient()
	if !assert.NoError(t, err) {
		return
	}
	poolID := networkPoolID()
	poolMetadataTest(t, client, poolID)
}

func poolMetadataTest(t TestingT, client *koios.Client, poolID koios.PoolID) {
	res, err := client.GetPoolMetadata(context.Background(), []koios.PoolID{poolID}, nil)
	if !assert.NoError(t, err) {
		return
	}
	for _, meta := range res.Data {
		assertNotEmpty(t, meta.PoolID, "pool_id_bech32")
		assertNotEmpty(t, meta.MetaURL, "meta_url")
		assertNotEmpty(t, meta.MetaHash, "meta_hash")
		assertNotEmpty(t, meta.MetaJSON.Name, "meta_json.name")
		assertNotEmpty(t, meta.MetaJSON.Ticker, "meta_json.ticker")
		assertNotEmpty(t, meta.MetaJSON.Homepage, "meta_json.homepage")
		assertNotEmpty(t, meta.MetaJSON.Description, "meta_json.description")
	}
}
