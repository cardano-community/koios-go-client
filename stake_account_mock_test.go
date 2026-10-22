// SPDX-License-Identifier: Apache-2.0
//
// Copyright © 2022 The Cardano Community Authors

package koios_test

import (
	"context"
	"testing"

	koios "github.com/cardano-community/koios-go-client/v5"
)

const mockStake = koios.Address("stake1u8yxtugdv63wxafy9d00nuz6hjyyp4qnggvc9a3vxh8yl0ckml2uz")

func TestMockGetAccountList(t *testing.T) {
	client, capture := mockClient(t, "account_list")
	res, err := client.GetAccountList(context.Background(), nil)
	mockNoErr(t, err)
	mockExpect(t, capture, "GET", "/account_list", nil, "")
	mockLen(t, "data", res.Data, 1)
	mockEq(t, "stake_address", res.Data[0].StakeAddress, mockStake)
	mockEq(t, "stake_address_hex", res.Data[0].StakeAddressHex, "e1c865f10d66a2e375242b5ef9f05abc8840d413421982f62c35ce4fbf")
	mockEq(t, "script_hash", res.Data[0].ScriptHash, koios.ScriptHash("bf357f5de69e4aad71954bebd64357a4813ea5233df12fce4a9de582"))
}

func TestMockGetAccountInfo(t *testing.T) {
	for name, call := range map[string]func(*koios.Client) (*koios.AccountsInfoResponse, error){
		"account_info": func(c *koios.Client) (*koios.AccountsInfoResponse, error) {
			return c.GetAccountInfo(context.Background(), []koios.Address{mockStake}, nil)
		},
		"account_info_cached": func(c *koios.Client) (*koios.AccountsInfoResponse, error) {
			return c.GetAccountInfoCached(context.Background(), []koios.Address{mockStake}, nil)
		},
	} {
		t.Run(name, func(t *testing.T) {
			client, capture := mockClient(t, name)
			res, err := call(client)
			mockNoErr(t, err)
			mockExpect(t, capture, "POST", "/"+name, nil, `{"_stake_addresses": ["`+string(mockStake)+`"]}`)
			mockLen(t, "data", res.Data, 1)
			a := res.Data[0]
			mockEq(t, "status", a.Status, "registered")
			if a.DelegatedDRep == nil || a.DelegatedPool == nil {
				t.Fatal("delegated_drep and delegated_pool must be set")
			}
			mockEq(t, "delegated_drep", *a.DelegatedDRep, "drep1yfhyq6tztjksqqpd5lglc3zr2tn8vylgjh9xzz7n2p4l4lgk3qam3")
			mockEq(t, "delegated_pool", *a.DelegatedPool, koios.PoolID("pool155efqn9xpcf73pphkk88cmlkdwx4ulkg606tne970qswczg3asc"))
			mockEq(t, "total_balance", a.TotalBalance.String(), "207116800428")
			mockEq(t, "deposit", a.Deposit.String(), "2000000")
			mockEq(t, "proposal_refund", a.ProposalRefund.String(), "0")
		})
	}

	client, _ := mockClientRaw(t, `[{"stake_address": "`+string(mockStake)+`", "status": "not registered",
		"delegated_drep": null, "delegated_pool": null, "total_balance": "0", "utxo": "0", "rewards": "0",
		"withdrawals": "0", "rewards_available": "0", "deposit": "0", "reserves": "0", "treasury": "0",
		"proposal_refund": "125000000000"}]`)
	res, err := client.GetAccountInfo(context.Background(), []koios.Address{mockStake}, nil)
	mockNoErr(t, err)
	if res.Data[0].DelegatedDRep != nil || res.Data[0].DelegatedPool != nil {
		t.Error("null delegations must decode to nil")
	}
	mockEq(t, "proposal_refund", res.Data[0].ProposalRefund.String(), "125000000000")
}

func TestMockGetAccountTxs(t *testing.T) {
	client, capture := mockClient(t, "account_txs")
	res, err := client.GetAccountTxs(context.Background(), mockStake, 1000, nil)
	mockNoErr(t, err)
	mockExpect(t, capture, "GET", "/account_txs", map[string]string{
		"_stake_address":      string(mockStake),
		"_after_block_height": "1000",
	}, "")
	mockLen(t, "data", res.Data, 1)
	mockEq(t, "tx_hash", res.Data[0].TxHash, koios.TxHash("f144a8264acf4bdfe2e1241170969c930d64ab6b0996a4a45237b623f1dd670e"))
}

func TestMockGetAccountHistory(t *testing.T) {
	client, capture := mockClient(t, "account_history")
	res, err := client.GetAccountHistory(context.Background(), []koios.Address{mockStake}, nil)
	mockNoErr(t, err)
	mockExpect(t, capture, "POST", "/account_history", nil, `{"_stake_addresses": ["`+string(mockStake)+`"]}`)
	mockLen(t, "data", res.Data, 1)
	mockLen(t, "history", res.Data[0].History, 1)
}

func TestMockGetAccountRewardHistory(t *testing.T) {
	client, capture := mockClient(t, "account_reward_history")
	res, err := client.GetAccountRewardHistory(context.Background(), []koios.Address{mockStake}, mockEpochPtr(294), nil)
	mockNoErr(t, err)
	mockExpect(t, capture, "POST", "/account_reward_history", nil,
		`{"_stake_addresses": ["`+string(mockStake)+`"], "_epoch_no": 294}`)
	mockLen(t, "data", res.Data, 1)
	r := res.Data[0]
	mockEq(t, "stake_address", r.StakeAddress, mockStake)
	mockEq(t, "earned_epoch", r.EarnedEpoch, koios.EpochNo(294))
	mockEq(t, "spendable_epoch", r.SpendableEpoch, koios.EpochNo(294))
	mockEq(t, "type", r.Type, "member")
	mockEq(t, "pool_id_bech32", r.PoolIDBech32, koios.PoolID("pool155efqn9xpcf73pphkk88cmlkdwx4ulkg606tne970qswczg3asc"))

	client, capture = mockClient(t, "account_reward_history")
	_, err = client.GetAccountRewardHistory(context.Background(), []koios.Address{mockStake}, nil, nil)
	mockNoErr(t, err)
	mockExpect(t, capture, "POST", "/account_reward_history", nil, `{"_stake_addresses": ["`+string(mockStake)+`"]}`)
}

func TestMockGetAccountUpdateHistory(t *testing.T) {
	client, capture := mockClient(t, "account_update_history")
	res, err := client.GetAccountUpdateHistory(context.Background(), []koios.Address{mockStake}, nil)
	mockNoErr(t, err)
	mockExpect(t, capture, "POST", "/account_update_history", nil, `{"_stake_addresses": ["`+string(mockStake)+`"]}`)
	mockLen(t, "data", res.Data, 1)
	u := res.Data[0]
	mockEq(t, "stake_address", u.StakeAddress, mockStake)
	mockEq(t, "action_type", u.ActionType, "registration")
	mockEq(t, "epoch_slot", u.EpochSlot, koios.Slot(171530))
	mockEq(t, "absolute_slot", u.AbsoluteSlot, koios.Slot(49073930))
	mockEq(t, "block_time", u.BlockTime.Unix(), int64(1704757130))
}

func TestMockGetAccountStakeHistory(t *testing.T) {
	client, capture := mockClient(t, "account_stake_history")
	res, err := client.GetAccountStakeHistory(context.Background(), []koios.Address{mockStake}, nil)
	mockNoErr(t, err)
	mockExpect(t, capture, "POST", "/account_stake_history", nil, `{"_stake_addresses": ["`+string(mockStake)+`"]}`)
	mockLen(t, "data", res.Data, 1)
	mockEq(t, "pool_id_bech32", res.Data[0].PoolIDBech32, koios.PoolID("pool1z5uqdk7dzdxaae5633fqfcu2eqzy3a3rgtuvy087fdld7yws0xt"))
	mockEq(t, "epoch_no", res.Data[0].EpochNo, koios.EpochNo(301))
	mockEq(t, "active_stake", res.Data[0].ActiveStake.String(), "682334162")
}
