// SPDX-License-Identifier: Apache-2.0
//
// Copyright © 2022 The Cardano Community Authors

package koios_test

import (
	"context"
	"testing"

	koios "github.com/cardano-community/koios-go-client/v5"
)

func TestMockGetTip(t *testing.T) {
	client, capture := mockClient(t, "tip")
	res, err := client.GetTip(context.Background(), nil)
	mockNoErr(t, err)
	mockExpect(t, capture, "GET", "/tip", nil, "")
	mockEq(t, "era", res.Data.Era, "Conway")
	mockEq(t, "abs_slot", res.Data.AbsSlot, koios.Slot(49073930))
	mockEq(t, "block_no", res.Data.BlockNo, koios.BlockNo(1794506))
	mockEq(t, "hash", res.Data.Hash, koios.BlockHash("2fa5178c5be950c7f40d6c74efe9b3dd12c4836dc3f3dd052cd1c2a12edd477f"))
}

func TestMockGetGenesis(t *testing.T) {
	client, capture := mockClientRaw(t, `[{
		"networkmagic": "764824073", "networkid": "Mainnet", "epochlength": "432000",
		"slotlength": "1", "maxlovelacesupply": "45000000000000000", "systemstart": 1506203091,
		"activeslotcoeff": "0.05", "slotsperkesperiod": "129600", "maxkesrevolutions": "62",
		"securityparam": "2160", "updatequorum": "5",
		"alonzogenesis": "{\"lovelacePerUTxOWord\":34482}"
	}]`)
	res, err := client.GetGenesis(context.Background(), nil)
	mockNoErr(t, err)
	mockExpect(t, capture, "GET", "/genesis", nil, "")
	mockEq(t, "networkid", res.Data.NetworkID, "Mainnet")
	mockEq(t, "epochlength", res.Data.EpochLength.String(), "432000")
	mockEq(t, "activeslotcoeff", res.Data.ActiveSlotCoeff.String(), "0.05")
	mockEq(t, "systemstart", res.Data.SystemStart.Unix(), int64(1506203091))
	mockEq(t, "alonzogenesis", res.Data.AlonzoGenesis, `{"lovelacePerUTxOWord":34482}`)
	alonzo, err := res.Data.AlonzoGenesisMap()
	mockNoErr(t, err)
	mockEq(t, "lovelacePerUTxOWord", alonzo["lovelacePerUTxOWord"].(float64), 34482.0)
}

func TestMockGetTotals(t *testing.T) {
	client, capture := mockClient(t, "totals")
	epoch := koios.EpochNo(294)
	res, err := client.GetTotals(context.Background(), &epoch, nil)
	mockNoErr(t, err)
	mockExpect(t, capture, "GET", "/totals", map[string]string{"_epoch_no": "294"}, "")
	mockLen(t, "data", res.Data, 1)
	got := res.Data[0]
	mockEq(t, "treasury", got.Treasury.String(), "637024173474141")
	mockEq(t, "fees", got.Fees.String(), "92010214941")
	mockEq(t, "deposits_stake", got.DepositsStake.String(), "3338702000000")
	mockEq(t, "deposits_drep", got.DepositsDRep.String(), "0")
	mockEq(t, "deposits_proposal", got.DepositsProposal.String(), "0")
	mockEq(t, "treasury_donation", got.TreasuryDonation.String(), "0")
	mockEq(t, "treasury_withdrawal", got.TreasuryWithdrawal.String(), "0")
	mockEq(t, "reserves_withdrawal", got.ReservesWithdrawal.String(), "0")
}

func TestMockReserveAndTreasuryWithdrawals(t *testing.T) {
	client, capture := mockClient(t, "reserve_withdrawals")
	reserves, err := client.GetReserveWithdrawals(context.Background(), nil)
	mockNoErr(t, err)
	mockExpect(t, capture, "GET", "/reserve_withdrawals", nil, "")
	mockLen(t, "data", reserves.Data, 1)
	mockEq(t, "earned_epoch", reserves.Data[0].EarnedEpoch, koios.EpochNo(294))
	mockEq(t, "spendable_epoch", reserves.Data[0].SpendableEpoch, koios.EpochNo(294))
	mockEq(t, "amount", reserves.Data[0].Amount.String(), "64328591517480")

	client, capture = mockClient(t, "treasury_withdrawals")
	treasury, err := client.GetTreasuryWithdrawals(context.Background(), nil)
	mockNoErr(t, err)
	mockExpect(t, capture, "GET", "/treasury_withdrawals", nil, "")
	mockLen(t, "data", treasury.Data, 1)
	mockEq(t, "earned_epoch", treasury.Data[0].EarnedEpoch, koios.EpochNo(294))
	mockEq(t, "spendable_epoch", treasury.Data[0].SpendableEpoch, koios.EpochNo(294))
}

func TestMockGetEraSummaries(t *testing.T) {
	client, capture := mockClientRaw(t, `[
		{"era": "Byron", "protocol_major": null, "protocol_minor": null, "ledger_protocol": "Classic",
		 "consensus_mechanism": "Ouroboros Classic", "notes": "", "epoch_no": 0,
		 "first_block_time": 1506203091, "first_block_hash": "f0f7892b5c333cffc4b3c4344de48af4cc63f55e44936196f365a9ef2244134f"},
		{"era": "Conway", "protocol_major": 9, "protocol_minor": 0, "ledger_protocol": "Praos",
		 "consensus_mechanism": "Ouroboros Genesis/Praos", "notes": "Chang HF", "epoch_no": 507,
		 "first_block_time": 1506635091, "first_block_hash": "f9dc2a2fc3a2db09a71af007a740261de585afc9e3022b8e30535592ff4dd9e5"}
	]`)
	res, err := client.GetEraSummaries(context.Background(), nil)
	mockNoErr(t, err)
	mockExpect(t, capture, "GET", "/era_summaries", nil, "")
	mockLen(t, "data", res.Data, 2)
	if res.Data[0].ProtocolMajor != nil || res.Data[0].ProtocolMinor != nil {
		t.Errorf("byron protocol version should be nil, got %v.%v", res.Data[0].ProtocolMajor, res.Data[0].ProtocolMinor)
	}
	mockEq(t, "era", res.Data[1].Era, "Conway")
	if res.Data[1].ProtocolMajor == nil {
		t.Fatal("conway protocol_major is nil")
	}
	mockEq(t, "protocol_major", *res.Data[1].ProtocolMajor, 9)
	mockEq(t, "ledger_protocol", res.Data[1].LedgerProtocol, "Praos")
	mockEq(t, "consensus_mechanism", res.Data[1].ConsensusMechanism, "Ouroboros Genesis/Praos")
	mockEq(t, "notes", res.Data[1].Notes, "Chang HF")
	mockEq(t, "epoch_no", res.Data[1].EpochNo, koios.EpochNo(507))
	mockEq(t, "first_block_time", res.Data[1].FirstBlockTime.Unix(), int64(1506635091))
	mockEq(t, "first_block_hash", res.Data[1].FirstBlockHash, koios.BlockHash("f9dc2a2fc3a2db09a71af007a740261de585afc9e3022b8e30535592ff4dd9e5"))
}

func TestMockGetCliProtocolParams(t *testing.T) {
	client, capture := mockClientRaw(t, `{"collateralPercentage": 150, "protocolVersion": {"major": 8, "minor": 0}, "txFeeFixed": 155381}`)
	res, err := client.GetCliProtocolParams(context.Background(), nil)
	mockNoErr(t, err)
	mockExpect(t, capture, "GET", "/cli_protocol_params", nil, "")
	mockJSONEq(t, "data", res.Data, `{"collateralPercentage": 150, "protocolVersion": {"major": 8, "minor": 0}, "txFeeFixed": 155381}`)
}
