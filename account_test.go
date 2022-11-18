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
	"strings"
	"testing"

	"github.com/cardano-community/koios-go-client/v3"
	"github.com/stretchr/testify/assert"
)

func TestAccounts(t *testing.T) {
	client, err := getLiveClient()
	if testIsLocal(t, err) {
		return
	}

	accountsTest(t, client)
}

func accountsTest(t TestingT, client *koios.Client) {

	opts := client.NewRequestOptions()
	opts.SetPageSize(10)
	res, err := client.GetAccounts(context.Background(), opts)
	if !assert.NoError(t, err) {
		return
	}
	assertEqual(t, 10, len(res.Data), "total blocks returned")

	for _, addr := range res.Data {
		assert.True(
			t,
			strings.HasPrefix(addr.String(), "stake"),
			fmt.Sprintf("account list returned %s", addr))
	}
}

func TestAccountInfo(t *testing.T) {
	client, err := getLiveClient()
	if testIsLocal(t, err) {
		return
	}
	accountInfoTest(t, networkAccounts(), false, client)
}

func TestAccountInfoCached(t *testing.T) {
	client, err := getLiveClient()
	if testIsLocal(t, err) {
		return
	}
	accountInfoTest(t, networkAccounts(), true, client)
}

func accountInfoTest(t TestingT, accs []koios.Address, cached bool, client *koios.Client) {

	res, err := client.GetAccountsInfo(context.Background(), accs, cached, nil)
	if !assert.NoError(t, err) {
		return
	}

	if !cached {
		assertEqual(t, len(accs), len(res.Data), "total account infos returned")
	}

	for i, account := range res.Data {
		label := fmt.Sprintf("account[%d]", i)
		assertNotEmpty(t, account.StakeAddress, label+".stake_address")
		assertNotEmpty(t, account.Status, label+".status")
		if account.DelegatedPool != nil {
			assertNotEmpty(t, account.DelegatedPool, label+".delegated_pool")
		}
		// assertIsPositive(t, account.TotalBalance, label+".total_balance")
		// assertIsPositive(t, account.UTxO, label+".utxo")
		// assertIsPositive(t, account.Rewards, label+".rewards")
		// assertIsPositive(t, account.Withdrawals, label+".withdrawals")
		// assertIsPositive(t, account.RewardsAvailable, label+".rewards_available")
		// assertIsPositive(t, account.Reserves, label+".reserves")
		// assertIsPositive(t, account.Treasury, label+".treasury")
	}
}

func TestAccountRewards(t *testing.T) {
	client, err := getLiveClient()
	if testIsLocal(t, err) {
		return
	}
	epoch := networkEpoch()
	accountRewardsTest(t, networkAccounts(), &epoch, client)
}

func accountRewardsTest(t TestingT, accs []koios.Address, epoch *koios.EpochNo, client *koios.Client) {
	res, err := client.GetAccountsRewards(context.Background(), accs, epoch, nil)
	if !assert.NoError(t, err) {
		return
	}

	for i, account := range res.Data {
		label := fmt.Sprintf("account[%d]", i)
		assertNotEmpty(t, account.StakeAddress, label+".stake_address")
		for _, reward := range account.Rewards {
			label := fmt.Sprintf("account[%s][%d]", account.StakeAddress, reward.EarnedEpoch)
			assertGreater(t, reward.EarnedEpoch, 0, label+".earned_epoch")
			assertGreater(t, reward.SpendableEpoch, 0, label+".spendable_epoch")
			assertIsPositive(t, reward.Amount, label+".amount")
			assertNotEmpty(t, reward.Type, label+".type")
			assertNotEmpty(t, reward.PoolID, label+".pool_id")
		}
	}
}

func TestAccountUpdates(t *testing.T) {
	client, err := getLiveClient()
	if testIsLocal(t, err) {
		return
	}
	accountUpdatesTest(t, networkAccounts(), client)
}

func accountUpdatesTest(t TestingT, accs []koios.Address, client *koios.Client) {
	res, err := client.GetAccountsUpdates(context.Background(), accs, nil)
	if !assert.NoError(t, err) {
		return
	}
	for i, account := range res.Data {
		label := fmt.Sprintf("account[%d]", i)
		assertNotEmpty(t, account.StakeAddress, label+".stake_address")
		for j, update := range account.Updates {
			label := fmt.Sprintf("account[%s].updates[%d]", account.StakeAddress, j)
			assertNotEmpty(t, update.TxHash, label+".tx_hash")
			assertGreater(t, update.EpochNo, 0, label+".epoch_no")
			assertGreater(t, update.EpochSlot, 0, label+".epoch_slot")
			assertGreater(t, update.AbsoluteSlot, 0, label+".absolute_slot")
			assertTimeNotZero(t, update.BlockTime, label+".block_time")
		}
	}
}

func TestAccountAddresses(t *testing.T) {
	client, err := getLiveClient()
	if testIsLocal(t, err) {
		return
	}
	accountAddressesTest(t, networkAccounts(), client)
}

func accountAddressesTest(t TestingT, accs []koios.Address, client *koios.Client) {
	res, err := client.GetAccountsAddresses(context.Background(), accs, nil)
	if !assert.NoError(t, err) {
		return
	}

	for i, account := range res.Data {
		label := fmt.Sprintf("stake_address[%d]", i)
		assertNotEmpty(t, account.StakeAddress, label)
		for j, addr := range account.Addresses {
			assertNotEmpty(t, addr, fmt.Sprintf("stake_address[%d].addresses[%d]", i, j))
		}
	}
}

func TestAccountAssets(t *testing.T) {
	client, err := getLiveClient()
	if testIsLocal(t, err) {
		return
	}
	accountAssetsTest(t, networkAccounts(), client)
}

func accountAssetsTest(t TestingT, accs []koios.Address, client *koios.Client) {
	res, err := client.GetAccountsAssets(context.Background(), accs, nil)
	if !assert.NoError(t, err) {
		return
	}

	for i, account := range res.Data {
		label := fmt.Sprintf("acc[%d]", i)
		assertNotEmpty(t, account.StakeAddress, label)
		for j, asset := range account.Assets {
			assertNotEmpty(t, asset.PolicyID, fmt.Sprintf("acc[%d].assets[%d].policy_id", i, j))
			// asset name can be empty
			if len(asset.AssetName) > 0 {
				assertNotEmpty(t, asset.AssetName, fmt.Sprintf("acc[%d].assets[%d].asset_name", i, j))
			}
			assertNotEmpty(t, asset.Fingerprint, fmt.Sprintf("acc[%d].assets[%d].fingerprint", i, j))
			assertIsPositive(t, asset.Quantity, fmt.Sprintf("acc[%d].assets[%d].quantity", i, j))
		}
	}
}

func TestAccountHistory(t *testing.T) {
	client, err := getLiveClient()
	if testIsLocal(t, err) {
		return
	}
	accountHistoryTest(t, networkAccounts(), client)
}

func accountHistoryTest(t TestingT, accs []koios.Address, client *koios.Client) {
	res, err := client.GetAccountsHistory(context.Background(), accs, nil)
	if !assert.NoError(t, err) {
		return
	}
	for _, account := range res.Data {
		label := fmt.Sprintf("acc[%s].stake_address", account.StakeAddress)
		assertNotEmpty(t, account.StakeAddress, label)
		for i, history := range account.History {
			label = fmt.Sprintf("%s.history[%d]", label, i)
			assertNotEmpty(t, history.PoolID, label+".pool_id")
			assertGreater(t, history.EpochNo, 0, label+".epoch_no")
			assertIsPositive(t, history.ActiveStake, label+".active_stake")
		}
	}
}
