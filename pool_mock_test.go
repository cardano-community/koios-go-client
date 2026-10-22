// SPDX-License-Identifier: Apache-2.0
//
// Copyright © 2022 The Cardano Community Authors

package koios_test

import (
	"context"
	"errors"
	"testing"

	koios "github.com/cardano-community/koios-go-client/v5"
)

const mockPool = koios.PoolID("pool155efqn9xpcf73pphkk88cmlkdwx4ulkg606tne970qswczg3asc")

func TestMockGetPoolList(t *testing.T) {
	client, capture := mockClient(t, "pool_list")
	res, err := client.GetPoolList(context.Background(), nil)
	mockNoErr(t, err)
	mockExpect(t, capture, "GET", "/pool_list", nil, "")
	mockLen(t, "data", res.Data, 1)
	p := res.Data[0]
	mockEq(t, "margin", p.Margin, 0.1)
	mockEq(t, "deposit", p.Deposit.String(), "500000000")
	mockEq(t, "pool_group", p.PoolGroup, "IOG")
	mockEq(t, "active_stake", p.ActiveStake.String(), "64328627680963")
	mockEq(t, "pool_status", p.PoolStatus, "registered")
	mockLen(t, "relays", p.Relays, 1)
	mockEq(t, "relay ipv4", p.Relays[0].Ipv4, "54.220.20.40")
	mockEq(t, "relay ipv6", p.Relays[0].Ipv6, "2604:ed40:1000:1711:6082:78ff:fe0c:ebf")
	mockEq(t, "relay port", p.Relays[0].Port, uint16(6000))
	if p.RetiringEpoch != nil {
		t.Errorf("retiring_epoch = %v, want nil", *p.RetiringEpoch)
	}

	client, _ = mockClientRaw(t, `[{"pool_id_bech32": "`+string(mockPool)+`", "pool_status": "retiring", "retiring_epoch": 400, "relays": []}]`)
	res, err = client.GetPoolList(context.Background(), nil)
	mockNoErr(t, err)
	if res.Data[0].RetiringEpoch == nil {
		t.Fatal("retiring_epoch must decode")
	}
	mockEq(t, "retiring_epoch", *res.Data[0].RetiringEpoch, koios.EpochNo(400))
}

func TestMockGetPoolInfo(t *testing.T) {
	client, capture := mockClient(t, "pool_info")
	res, err := client.GetPoolInfo(context.Background(), mockPool, nil)
	mockNoErr(t, err)
	mockExpect(t, capture, "POST", "/pool_info", nil, `{"_pool_bech32_ids": ["`+string(mockPool)+`"]}`)
	if res.Data == nil {
		t.Fatal("data is nil")
	}
	p := res.Data
	if p.RewardAddrDelegatedDRep == nil {
		t.Fatal("reward_addr_delegated_drep must be set")
	}
	mockEq(t, "reward_addr_delegated_drep", *p.RewardAddrDelegatedDRep, "drep1yfhyq6tztjksqqpd5lglc3zr2tn8vylgjh9xzz7n2p4l4lgk3qam3")
	mockEq(t, "deposit", p.Deposit.String(), "500000000")
	mockEq(t, "voting_power", p.VotingPower.String(), "123456789")
	mockEq(t, "margin", p.Margin, 0.1)
	mockEq(t, "live_saturation", p.LiveSaturation, 94.52)
	mockEq(t, "block_count", p.BlockCount, uint64(4509))
	if p.MetaJSON == nil || p.MetaJSON.Ticker == nil {
		t.Fatal("meta_json must be set")
	}
	mockEq(t, "meta_json.ticker", *p.MetaJSON.Ticker, "IOGP")

	client, _ = mockClientRaw(t, `[{"pool_id_bech32": "`+string(mockPool)+`", "meta_json": null}]`)
	res, err = client.GetPoolInfo(context.Background(), mockPool, nil)
	mockNoErr(t, err)
	if res.Data.MetaJSON != nil {
		t.Error("null meta_json must decode to nil")
	}
}

func TestMockGetPoolUpdates(t *testing.T) {
	client, capture := mockClient(t, "pool_updates")
	res, err := client.GetPoolUpdates(context.Background(), mockPool, nil)
	mockNoErr(t, err)
	mockExpect(t, capture, "GET", "/pool_updates", map[string]string{"_pool_bech32": string(mockPool)}, "")
	mockLen(t, "data", res.Data, 1)
	mockEq(t, "update_type", res.Data[0].UpdateType, "registration")
	mockEq(t, "margin", res.Data[0].Margin, 0.1)
	if res.Data[0].MetaJSON == nil || res.Data[0].MetaJSON.Name == nil {
		t.Fatal("meta_json must be set")
	}
	mockEq(t, "meta_json.name", *res.Data[0].MetaJSON.Name, "Input Output Global (IOHK) - Private")
}

func TestMockGetPoolMetadata(t *testing.T) {
	client, capture := mockClient(t, "pool_metadata")
	res, err := client.GetPoolMetadata(context.Background(), nil, nil)
	mockNoErr(t, err)
	mockExpect(t, capture, "POST", "/pool_metadata", nil, `{}`)
	mockLen(t, "data", res.Data, 1)

	_, err = client.GetPoolMetadata(context.Background(), []koios.PoolID{mockPool}, nil)
	mockNoErr(t, err)
	mockExpect(t, capture, "POST", "/pool_metadata", nil, `{"_pool_bech32_ids": ["`+string(mockPool)+`"]}`)
}

func TestMockGetPoolRelays(t *testing.T) {
	client, capture := mockClient(t, "pool_relays")
	res, err := client.GetPoolRelays(context.Background(), nil)
	mockNoErr(t, err)
	mockExpect(t, capture, "GET", "/pool_relays", nil, "")
	mockLen(t, "data", res.Data, 1)
	mockEq(t, "pool_id_bech32", res.Data[0].PoolIDBech32, mockPool)
	mockLen(t, "relays", res.Data[0].Relays, 1)
}

func TestMockGetPoolHistory(t *testing.T) {
	client, capture := mockClient(t, "pool_history")
	res, err := client.GetPoolHistory(context.Background(), mockPool, 312, nil)
	mockNoErr(t, err)
	mockExpect(t, capture, "GET", "/pool_history", map[string]string{"_pool_bech32": string(mockPool), "_epoch_no": "312"}, "")
	mockLen(t, "data", res.Data, 1)
	mockEq(t, "epoch_ros", res.Data[0].EpochROS, 3.000340466)
	mockEq(t, "member_rewards", res.Data[0].MemberRewards.String(), "123456780123")
}

func TestMockGetPoolInvalidDelegators(t *testing.T) {
	client, capture := mockClient(t, "pool_invalid_delegators")
	res, err := client.GetPoolInvalidDelegators(context.Background(), mockPool, mockEpochPtr(294), nil)
	mockNoErr(t, err)
	mockExpect(t, capture, "GET", "/pool_invalid_delegators", map[string]string{"_pool_bech32": string(mockPool), "_epoch_no": "294"}, "")
	mockLen(t, "data", res.Data, 1)
	mockEq(t, "stake_address", res.Data[0].StakeAddress, mockStake)
	mockEq(t, "epoch_no", res.Data[0].EpochNo, koios.EpochNo(294))

	client, capture = mockClient(t, "pool_invalid_delegators")
	_, err = client.GetPoolInvalidDelegators(context.Background(), mockPool, nil, nil)
	mockNoErr(t, err)
	mockExpect(t, capture, "GET", "/pool_invalid_delegators", map[string]string{"_pool_bech32": string(mockPool)}, "")
}

func TestMockGetPoolOwnerHistory(t *testing.T) {
	client, capture := mockClient(t, "pool_owner_history")
	res, err := client.GetPoolOwnerHistory(context.Background(), []koios.PoolID{mockPool}, nil)
	mockNoErr(t, err)
	mockExpect(t, capture, "POST", "/pool_owner_history", nil, `{"_pool_bech32_ids": ["`+string(mockPool)+`"]}`)
	mockLen(t, "data", res.Data, 1)
	mockEq(t, "pool_id_bech32", res.Data[0].PoolIDBech32, mockPool)
	mockEq(t, "declared_pledge", res.Data[0].DeclaredPledge.String(), "64000000000000")
	mockEq(t, "epoch", res.Data[0].Epoch, koios.EpochNo(312))
	mockEq(t, "active_stake", res.Data[0].ActiveStake.String(), "64328627680963")

	_, err = client.GetPoolOwnerHistory(context.Background(), nil, nil)
	if !errors.Is(err, koios.ErrNoPoolID) {
		t.Errorf("err = %v, want ErrNoPoolID", err)
	}
}

func TestMockGetPoolGroups(t *testing.T) {
	client, capture := mockClient(t, "pool_groups")
	res, err := client.GetPoolGroups(context.Background(), nil)
	mockNoErr(t, err)
	mockExpect(t, capture, "GET", "/pool_groups", nil, "")
	mockLen(t, "data", res.Data, 1)
	mockEq(t, "pool_group", res.Data[0].PoolGroup, "IOG")
	mockEq(t, "ticker", res.Data[0].Ticker, "IOGP")
	mockEq(t, "adastat_group", res.Data[0].AdastatGroup, "iohk.io")
	mockEq(t, "balanceanalytics_group", res.Data[0].BalanceanalyticsGroup, "IOG")
}

func TestMockGetPoolCalidusKeys(t *testing.T) {
	client, capture := mockClient(t, "pool_calidus_keys")
	res, err := client.GetPoolCalidusKeys(context.Background(), nil)
	mockNoErr(t, err)
	mockExpect(t, capture, "GET", "/pool_calidus_keys", nil, "")
	mockLen(t, "data", res.Data, 1)
	k := res.Data[0]
	mockEq(t, "pool_status", k.PoolStatus, "registered")
	mockEq(t, "calidus_nonce", k.CalidusNonce, uint64(150720822))
	mockEq(t, "calidus_pub_key", k.CalidusPubKey, "dea2d81424a285f7fd2846beb39320543ba1df27b890789b2c152b64a352e2df")
	mockEq(t, "calidus_id_bech32", k.CalidusIDBech32, "calidus158w8ft8k4endw5q6m0rx2myl2ldvw5f5d8xmwyem6529fcgn575zh")
	mockEq(t, "registered", k.Registered, true)
	mockEq(t, "bytes", k.Bytes, "19029a324512")
	mockEq(t, "block_height", k.BlockHeight, uint64(1794506))
}
