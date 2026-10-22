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

const mockTxHash = koios.TxHash("f144a8264acf4bdfe2e1241170969c930d64ab6b0996a4a45237b623f1dd670e")

func TestMockGetTxInfo(t *testing.T) {
	client, capture := mockClient(t, "tx_info")
	flags := &koios.TxInfoFlags{Inputs: true, Scripts: true, Governance: true}
	res, err := client.GetTxInfo(context.Background(), []koios.TxHash{mockTxHash}, flags, nil)
	mockNoErr(t, err)
	mockExpect(t, capture, "POST", "/tx_info", nil,
		`{"_tx_hashes": ["`+string(mockTxHash)+`"], "_inputs": true, "_scripts": true, "_governance": true}`)
	mockLen(t, "data", res.Data, 1)
	tx := res.Data[0]
	mockEq(t, "tx_hash", tx.TxHash, mockTxHash)
	mockEq(t, "treasury_donation", tx.TreasuryDonation.String(), "0")
	if tx.InvalidBefore == nil || tx.InvalidAfter == nil {
		t.Fatal("invalid_before/invalid_after must be set")
	}
	mockEq(t, "invalid_before", *tx.InvalidBefore, koios.Slot(42331172))
	mockEq(t, "invalid_after", *tx.InvalidAfter, koios.Slot(42332172))

	mockLen(t, "inputs", tx.Inputs, 1)
	in := tx.Inputs[0]
	mockEq(t, "payment_addr.cred", in.PaymentAddr.Cred, koios.PaymentCredential("de3c1c527e8826b9cd2030f88f75fc44cd4ce519b9ded9eb794b3794"))
	mockEq(t, "stake_addr", in.StakeAddress, koios.Address("stake1u8yxtugdv63wxafy9d00nuz6hjyyp4qnggvc9a3vxh8yl0ckml2uz"))
	mockEq(t, "value", in.Value.String(), "157832856")
	if in.InlineDatum == nil || in.ReferenceScript == nil {
		t.Fatal("inline_datum and reference_script must be set")
	}
	mockEq(t, "inline_datum.bytes", in.InlineDatum.Bytes, "19029a")
	mockJSONEq(t, "inline_datum.value", in.InlineDatum.Value, `{"int": 666}`)
	mockEq(t, "reference_script.type", in.ReferenceScript.Type, "plutusV1")
	mockEq(t, "reference_script.size", in.ReferenceScript.Size, uint(14))
	mockLen(t, "asset_list", in.AssetList, 1)
	mockEq(t, "asset_list.fingerprint", in.AssetList[0].Fingerprint, koios.AssetFingerprint("asset1ua6pz3yd5mdka946z8jw2fld3f8d0mmxt75gv9"))

	mockLen(t, "collateral_inputs", tx.CollateralInputs, 1)
	if tx.CollateralOutput == nil {
		t.Fatal("collateral_output must be set")
	}
	mockEq(t, "collateral_output.tx_index", tx.CollateralOutput.TxIndex, 0)
	if len(tx.CollateralOutput.AssetList) == 0 || tx.CollateralOutput.AssetList[0] != '"' {
		t.Fatalf("collateral_output.asset_list must be kept as the raw string, got %s", tx.CollateralOutput.AssetList)
	}
	mockLen(t, "reference_inputs", tx.ReferenceInputs, 1)
	mockLen(t, "outputs", tx.Outputs, 1)
	mockLen(t, "withdrawals", tx.Withdrawals, 1)
	mockEq(t, "withdrawal amount", tx.Withdrawals[0].Amount.String(), "9845162")

	mockLen(t, "native_scripts", tx.NativeScripts, 1)
	mockEq(t, "native script hash", tx.NativeScripts[0].ScriptHash, koios.ScriptHash("bfa7ffa9b2e164873db6ac6d0528c82e212963bc62e10fd1d81da4af"))
	if len(tx.NativeScripts[0].ScriptJSON) == 0 {
		t.Error("native script_json is empty")
	}

	mockLen(t, "plutus_contracts", tx.PlutusContracts, 1)
	pc := tx.PlutusContracts[0]
	if pc.SpendsInput == nil {
		t.Fatal("spends_input must be set")
	}
	mockEq(t, "spends_input.tx_hash", pc.SpendsInput.TxHash, mockTxHash)
	mockEq(t, "redeemer purpose", pc.Input.Redeemer.Purpose, "spend")
	mockEq(t, "redeemer unit mem", pc.Input.Redeemer.Unit.Mem.String(), "520448")
	mockEq(t, "redeemer unit steps", pc.Input.Redeemer.Unit.Steps.String(), "211535239")

	mockLen(t, "voting_procedures", tx.VotingProcedures, 1)
	mockEq(t, "voter_role", tx.VotingProcedures[0].VoterRole, "DRep")
	mockEq(t, "vote", tx.VotingProcedures[0].Vote, "Yes")
	mockLen(t, "proposal_procedures", tx.ProposalProcedures, 1)
	mockEq(t, "proposal type", tx.ProposalProcedures[0].Type, "ParameterChange")
	mockEq(t, "proposal deposit", tx.ProposalProcedures[0].Deposit.String(), "500000000")
	if tx.ProposalProcedures[0].Expiration != nil {
		t.Errorf("expiration = %v, want nil", *tx.ProposalProcedures[0].Expiration)
	}

	client, capture = mockClient(t, "tx_info")
	_, err = client.GetTxInfo(context.Background(), []koios.TxHash{mockTxHash}, nil, nil)
	mockNoErr(t, err)
	mockExpect(t, capture, "POST", "/tx_info", nil, `{"_tx_hashes": ["`+string(mockTxHash)+`"]}`)

	_, err = client.GetTxInfo(context.Background(), nil, nil, nil)
	if !errors.Is(err, koios.ErrNoTxHash) {
		t.Errorf("err = %v, want ErrNoTxHash", err)
	}
}

func TestMockGetUTxOInfo(t *testing.T) {
	client, capture := mockClient(t, "utxo_info")
	res, err := client.GetUTxOInfo(context.Background(), []koios.UTxORef{koios.UTxORef(string(mockTxHash) + "#0")}, true, nil)
	mockNoErr(t, err)
	mockExpect(t, capture, "POST", "/utxo_info", nil, `{"_utxo_refs": ["`+string(mockTxHash)+`#0"], "_extended": true}`)
	mockLen(t, "data", res.Data, 1)
	u := res.Data[0]
	mockEq(t, "tx_hash", u.TxHash, mockTxHash)
	mockEq(t, "address", u.Address, koios.Address("addr1qxkfe8s6m8qt5436lec3f0320hrmpppwqgs2gah4360krvyssntpwjcz303mx3h4avg7p29l3zd8u3jyglmewds9ezrqdc3cxp"))
	mockEq(t, "stake_address", u.StakeAddress, koios.Address("stake1u8yxtugdv63wxafy9d00nuz6hjyyp4qnggvc9a3vxh8yl0ckml2uz"))
	mockEq(t, "payment_cred", u.PaymentCred, koios.PaymentCredential("de3c1c527e8826b9cd2030f88f75fc44cd4ce519b9ded9eb794b3794"))
	mockEq(t, "block_time", u.BlockTime.Unix(), int64(1704757130))
	mockEq(t, "datum_hash", u.DatumHash, koios.DatumHash("5a595ce795815e81d22a1a522cf3987d546dc5bb016de61b002edd63a5413ec4"))
	mockEq(t, "is_spent", u.IsSpent, true)
	if u.InlineDatum == nil || u.ReferenceScript == nil {
		t.Fatal("inline_datum and reference_script must be set")
	}
	mockEq(t, "inline_datum.bytes", u.InlineDatum.Bytes, "19029a")
	mockEq(t, "reference_script.hash", u.ReferenceScript.Hash, koios.ScriptHash("67f33146617a5e61936081db3b2117cbf59bd2123748f58ac9678656"))
	mockEq(t, "reference_script.bytes", u.ReferenceScript.Bytes, "4e4d01000033222220051200120011")
}

func TestMockGetTxCbor(t *testing.T) {
	client, capture := mockClientRaw(t, `[{
		"tx_hash": "f144a8264acf4bdfe2e1241170969c930d64ab6b0996a4a45237b623f1dd670e",
		"block_hash": "2fa5178c5be950c7f40d6c74efe9b3dd12c4836dc3f3dd052cd1c2a12edd477f",
		"block_height": 1794506, "epoch_no": 117, "absolute_slot": 49073930,
		"tx_timestamp": 1704757130, "cbor": "84a40081825820", "valid_contract": false
	}]`)
	res, err := client.GetTxCbor(context.Background(), []koios.TxHash{mockTxHash}, nil)
	mockNoErr(t, err)
	mockExpect(t, capture, "POST", "/tx_cbor", nil, `{"_tx_hashes": ["`+string(mockTxHash)+`"]}`)
	mockLen(t, "data", res.Data, 1)
	mockEq(t, "cbor", res.Data[0].Cbor, "84a40081825820")
	mockEq(t, "valid_contract", res.Data[0].ValidContract, false)
	mockEq(t, "block_height", res.Data[0].BlockHeight, uint64(1794506))

	_, err = client.GetTxCbor(context.Background(), nil, nil)
	if !errors.Is(err, koios.ErrNoTxHash) {
		t.Errorf("err = %v, want ErrNoTxHash", err)
	}
}

func TestMockGetTxByMetalabel(t *testing.T) {
	client, capture := mockClient(t, "tx_by_metalabel")
	res, err := client.GetTxByMetalabel(context.Background(), 721, 1000, nil)
	mockNoErr(t, err)
	mockExpect(t, capture, "GET", "/tx_by_metalabel", map[string]string{"_label": "721", "_after_block_height": "1000"}, "")
	mockLen(t, "data", res.Data, 1)
	mockEq(t, "tx_hash", res.Data[0].TxHash, mockTxHash)
	mockEq(t, "absolute_slot", res.Data[0].AbsoluteSlot, koios.Slot(49073930))
	mockEq(t, "tx_timestamp", res.Data[0].TxTimestamp.Unix(), int64(1704757130))

	client, capture = mockClient(t, "tx_by_metalabel")
	_, err = client.GetTxByMetalabel(context.Background(), 674, 0, nil)
	mockNoErr(t, err)
	mockExpect(t, capture, "GET", "/tx_by_metalabel", map[string]string{"_label": "674"}, "")
}

func TestMockGetTxOutsEpoch(t *testing.T) {
	client, capture := mockClient(t, "tx_outs_epoch")
	res, err := client.GetTxOutsEpoch(context.Background(), 117, nil)
	mockNoErr(t, err)
	mockExpect(t, capture, "GET", "/tx_outs_epoch", map[string]string{"_epoch_no": "117"}, "")
	mockLen(t, "data", res.Data, 1)
	o := res.Data[0]
	mockEq(t, "tx_index", o.TxIndex, 0)
	mockEq(t, "value", o.Value.String(), "157832856")
	mockEq(t, "stake_address", o.StakeAddress, koios.Address("stake1u8yxtugdv63wxafy9d00nuz6hjyyp4qnggvc9a3vxh8yl0ckml2uz"))
	mockEq(t, "is_spent", o.IsSpent, true)
}

func TestMockGetTxUTxOs(t *testing.T) {
	client, capture := mockClient(t, "tx_utxos")
	res, err := client.GetTxUTxOs(context.Background(), []koios.TxHash{mockTxHash}, nil)
	mockNoErr(t, err)
	mockExpect(t, capture, "POST", "/tx_utxos", nil, `{"_tx_hashes": ["`+string(mockTxHash)+`"]}`)
	mockLen(t, "data", res.Data, 1)
	mockEq(t, "tx_hash", res.Data[0].TxHash, mockTxHash)
	mockLen(t, "inputs", res.Data[0].Inputs, 1)
	mockLen(t, "outputs", res.Data[0].Outputs, 1)
	mockEq(t, "input cred", res.Data[0].Inputs[0].PaymentAddr.Cred, koios.PaymentCredential("de3c1c527e8826b9cd2030f88f75fc44cd4ce519b9ded9eb794b3794"))
}

func TestMockGetTxStatus(t *testing.T) {
	client, capture := mockClientRaw(t, `[
		{"tx_hash": "f144a8264acf4bdfe2e1241170969c930d64ab6b0996a4a45237b623f1dd670e", "num_confirmations": 17},
		{"tx_hash": "0b8ba3bed976fa4913f19adc9f6dd9063138db5b4dd29cecde369456b5155e94", "num_confirmations": null}
	]`)
	res, err := client.GetTxStatus(context.Background(), []koios.TxHash{mockTxHash}, nil)
	mockNoErr(t, err)
	mockExpect(t, capture, "POST", "/tx_status", nil, `{"_tx_hashes": ["`+string(mockTxHash)+`"]}`)
	mockLen(t, "data", res.Data, 2)
	if res.Data[0].Confirmations == nil {
		t.Fatal("confirmations must be set")
	}
	mockEq(t, "confirmations", *res.Data[0].Confirmations, uint64(17))
	if res.Data[1].Confirmations != nil {
		t.Errorf("confirmations = %d, want nil", *res.Data[1].Confirmations)
	}
}
