// SPDX-License-Identifier: Apache-2.0
//
// Copyright © 2022 The Cardano Community Authors

package koios_test

// func TestNetworkTip(t *testing.T) {
// 	client, err := getLiveClient()
// 	if testIsLocal(t, err) {
// 		return
// 	}
// 	networkTipTest(t, client)
// }

// func networkTipTest(t TestingT, client *koios.Client) {
// 	tip, err := client.GetTip(context.Background(), nil)
// 	if !assert.NoError(t, err) {
// 		return
// 	}
// 	assertGreater(t, tip.Data.AbsSlot, koios.Slot(10000), "abs_slot")
// 	assertGreater(t, tip.Data.BlockNo, koios.BlockNo(10000), "block_no")
// 	assertGreater(t, tip.Data.EpochNo, koios.EpochNo(1), "epoch_no")
// 	assertGreater(t, tip.Data.EpochSlot, koios.Slot(1), "epoch_slot")
// 	assertNotEmpty(t, tip.Data.Hash, "hash")
// 	assertTimeNotZero(t, tip.Data.BlockTime, "block_time")
// }

// func TestNetworkGenesis(t *testing.T) {
// 	client, err := getLiveClient()
// 	if testIsLocal(t, err) {
// 		return
// 	}
// 	networkGenesisTest(t, client)
// }

// func networkGenesisTest(t TestingT, client *koios.Client) {
// 	res, err := client.GetGenesis(context.Background(), nil)
// 	if !assert.NoError(t, err) {
// 		return
// 	}
// 	assertIsPositive(t, res.Data.NetworkMagic, "networkmagic")
// 	assertNotEmpty(t, res.Data.NetworkID, "networkid")
// 	assertIsPositive(t, res.Data.EpochLength, "epochlength")
// 	assertIsPositive(t, res.Data.SlotLength, "slotlength")
// 	assertIsPositive(t, res.Data.MaxLovelaceSupply, "maxlovelacesupply")
// 	assertTimeNotZero(t, res.Data.SystemStart, "systemstart")
// 	assertIsPositive(t, res.Data.ActiveSlotCoeff, "activeslotcoeff")
// 	assertIsPositive(t, res.Data.SlotsPerKesPeriod, "slotsperkesperiod")
// 	assertIsPositive(t, res.Data.MaxKesRevolutions, "maxkesrevolutions")
// 	assertIsPositive(t, res.Data.SecurityParam, "securityparam")
// 	assertIsPositive(t, res.Data.UpdateQuorum, "updatequorum")
// 	assertNotEmpty(t, res.Data.AlonzoGenesis, "alonzogenesis")
// }

// func TestNetworkTotals(t *testing.T) {
// 	client, err := getLiveClient()
// 	if testIsLocal(t, err) {
// 		return
// 	}
// 	networkTotalsTest(t, networkEpoch(), client)
// }

// func networkTotalsTest(t TestingT, epoch koios.EpochNo, client *koios.Client) {
// 	res, err := client.GetTotals(context.Background(), &epoch, nil)
// 	if !assert.NoError(t, err) {
// 		return
// 	}
// 	zero := int64(0)
// 	assertEqual(t, epoch, res.Data[0].EpochNo, "epoch_no")
// 	assertGreater(t, res.Data[0].Circulation.IntPart(), zero, "circulation")
// 	assertGreater(t, res.Data[0].Reserves.IntPart(), zero, "reserves")
// 	assertGreater(t, res.Data[0].Reward.IntPart(), zero, "reward")
// 	assertGreater(t, res.Data[0].Supply.IntPart(), zero, "supply")
// 	assertGreater(t, res.Data[0].Treasury.IntPart(), zero, "treasury")
// }
