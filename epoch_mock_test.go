// SPDX-License-Identifier: Apache-2.0
//
// Copyright © 2022 The Cardano Community Authors

package koios_test

import (
	"context"
	"testing"

	koios "github.com/cardano-community/koios-go-client/v5"
)

func TestMockGetEpochInfo(t *testing.T) {
	client, capture := mockClient(t, "epoch_info")
	res, err := client.GetEpochInfo(context.Background(), 294, true, nil)
	mockNoErr(t, err)
	mockExpect(t, capture, "GET", "/epoch_info", map[string]string{
		"order":               "epoch_no.desc",
		"_epoch_no":           "294",
		"_include_next_epoch": "true",
	}, "")
	mockLen(t, "data", res.Data, 1)
	mockEq(t, "era", res.Data[0].Era, "Conway")
	mockEq(t, "out_sum", res.Data[0].OutSum.String(), "15432725054364942")
	mockEq(t, "active_stake", res.Data[0].ActiveStake.String(), "23395112387185880")
	mockEq(t, "blk_count", res.Data[0].BlkCount, 17321)
}

func TestMockGetEpochParams(t *testing.T) {
	client, capture := mockClient(t, "epoch_params")
	res, err := client.GetEpochParams(context.Background(), 294, nil)
	mockNoErr(t, err)
	mockExpect(t, capture, "GET", "/epoch_params", map[string]string{"_epoch_no": "294"}, "")
	mockLen(t, "data", res.Data, 1)
	p := res.Data[0]
	mockEq(t, "era", p.Era, "Conway")
	mockEq(t, "min_fee_a", p.MinFeeA.String(), "44")
	mockEq(t, "max_tx_ex_mem", p.MaxTxExMem, int64(10000000))
	mockEq(t, "max_tx_ex_steps", p.MaxTxExSteps, int64(10000000000))
	mockEq(t, "max_block_ex_steps", p.MaxBlockExSteps, int64(40000000000))
	mockEq(t, "max_val_size", p.MaxValSize, int64(5000))
	mockEq(t, "price_mem", p.PriceMem.String(), "0.0577")
	mockEq(t, "pvt_motion_no_confidence", p.PvtMotionNoConfidence, 0.6)
	mockEq(t, "pvt_committee_normal", p.PvtCommitteeNormal, 0.65)
	mockEq(t, "pvt_committee_no_confidence", p.PvtCommitteeNoConfidence, 0.65)
	mockEq(t, "pvt_hard_fork_initiation", p.PvtHardForkInitiation, 0.51)
	mockEq(t, "dvt_motion_no_confidence", p.DvtMotionNoConfidence, 0.67)
	mockEq(t, "dvt_committee_normal", p.DvtCommitteeNormal, 0.67)
	mockEq(t, "dvt_committee_no_confidence", p.DvtCommitteeNoConfidence, 0.65)
	mockEq(t, "dvt_update_to_constitution", p.DvtUpdateToConstitution, 0.75)
	mockEq(t, "dvt_hard_fork_initiation", p.DvtHardForkInitiation, 0.6)
	mockEq(t, "dvt_p_p_network_group", p.DvtPPNetworkGroup, 0.67)
	mockEq(t, "dvt_p_p_economic_group", p.DvtPPEconomicGroup, 0.67)
	mockEq(t, "dvt_p_p_technical_group", p.DvtPPTechnicalGroup, 0.67)
	mockEq(t, "dvt_p_p_gov_group", p.DvtPPGovGroup, 0.75)
	mockEq(t, "dvt_treasury_withdrawal", p.DvtTreasuryWithdrawal, 0.67)
	mockEq(t, "committee_min_size", p.CommitteeMinSize, 5)
	mockEq(t, "committee_max_term_length", p.CommitteeMaxTermLength, 146)
	mockEq(t, "gov_action_lifetime", p.GovActionLifetime, 14)
	mockEq(t, "gov_action_deposit", p.GovActionDeposit.String(), "100000000000")
	mockEq(t, "drep_deposit", p.DrepDeposit.String(), "500000000")
	mockEq(t, "drep_activity", p.DrepActivity, 20)
	mockEq(t, "pvtpp_security_group", p.PvtppSecurityGroup, 0.6)
	mockEq(t, "min_fee_ref_script_cost_per_byte", p.MinFeeRefScriptCostPerByte, 15.0)
}

func TestMockGetEpochBlockProtocols(t *testing.T) {
	client, capture := mockClient(t, "epoch_block_protocols")
	res, err := client.GetEpochBlockProtocols(context.Background(), koios.EpochNo(294), nil)
	mockNoErr(t, err)
	mockExpect(t, capture, "GET", "/epoch_block_protocols", map[string]string{"_epoch_no": "294"}, "")
	mockLen(t, "data", res.Data, 1)
	mockEq(t, "era", res.Data[0].Era, "Conway")
	mockEq(t, "proto_major", res.Data[0].ProtoMajor, 6)
	mockEq(t, "blocks", res.Data[0].Blocks, 2183)
}
