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

const mockScript = koios.ScriptHash("bfa7ffa9b2e164873db6ac6d0528c82e212963bc62e10fd1d81da4af")

func TestMockGetScriptRedeemers(t *testing.T) {
	client, capture := mockClientRaw(t, `[{"script_hash": "`+string(mockScript)+`", "redeemers": [{
		"tx_hash": "f144a8264acf4bdfe2e1241170969c930d64ab6b0996a4a45237b623f1dd670e", "tx_index": 0,
		"unit_mem": "520448", "unit_steps": 211535239, "fee": "45282", "purpose": "spend",
		"datum_hash": "5a595ce795815e81d22a1a522cf3987d546dc5bb016de61b002edd63a5413ec4",
		"datum_value": {"int": 666}}]}]`)
	res, err := client.GetScriptRedeemers(context.Background(), mockScript, nil)
	mockNoErr(t, err)
	mockExpect(t, capture, "GET", "/script_redeemers", map[string]string{"_script_hash": string(mockScript)}, "")
	if res.Data == nil {
		t.Fatal("data is nil")
	}
	mockLen(t, "redeemers", res.Data.Redeemers, 1)
	r := res.Data.Redeemers[0]
	mockEq(t, "unit_mem", r.UnitMem.String(), "520448")
	mockEq(t, "unit_steps", r.UnitSteps.String(), "211535239")
	mockEq(t, "fee", r.Fee.String(), "45282")
	mockEq(t, "purpose", r.Purpose, "spend")
	mockJSONEq(t, "datum_value", r.DatumValue, `{"int": 666}`)
}

func TestMockGetScriptInfo(t *testing.T) {
	client, capture := mockClient(t, "script_info")
	res, err := client.GetScriptInfo(context.Background(), []koios.ScriptHash{mockScript}, nil)
	mockNoErr(t, err)
	mockExpect(t, capture, "POST", "/script_info", nil, `{"_script_hashes": ["`+string(mockScript)+`"]}`)
	mockLen(t, "data", res.Data, 1)
	mockEq(t, "script_hash", res.Data[0].ScriptHash, mockScript)
	mockEq(t, "type", res.Data[0].Type, "plutusV1")
	mockEq(t, "creation_tx_hash", res.Data[0].CreationTxHash, koios.TxHash("255f061502ad83230351fbcf2d9fade1b5d118d332f92c9861075010a1fe3fbe"))
	mockEq(t, "size", res.Data[0].Size, uint(2039))
}

func TestMockGetReferenceScriptUTxOs(t *testing.T) {
	client, capture := mockClient(t, "reference_script_utxos")
	res, err := client.GetReferenceScriptUTxOs(context.Background(), []koios.ScriptHash{mockScript}, nil)
	mockNoErr(t, err)
	mockExpect(t, capture, "POST", "/reference_script_utxos", nil, `{"_script_hashes": ["`+string(mockScript)+`"]}`)
	mockLen(t, "data", res.Data, 1)
	mockEq(t, "script_hash", res.Data[0].ScriptHash, mockScript)
	mockEq(t, "tx_hash", res.Data[0].TxHash, koios.TxHash("f144a8264acf4bdfe2e1241170969c930d64ab6b0996a4a45237b623f1dd670e"))
	mockEq(t, "tx_index", res.Data[0].TxIndex, 0)

	_, err = client.GetReferenceScriptUTxOs(context.Background(), nil, nil)
	if !errors.Is(err, koios.ErrNoScriptHash) {
		t.Errorf("err = %v, want ErrNoScriptHash", err)
	}
}

func TestMockGetScriptUtxos(t *testing.T) {
	client, capture := mockClient(t, "script_utxos")
	res, err := client.GetScriptUtxos(context.Background(), mockScript, true, nil)
	mockNoErr(t, err)
	mockExpect(t, capture, "GET", "/script_utxos", map[string]string{"_script_hash": string(mockScript), "_extended": "true"}, "")
	mockLen(t, "data", res.Data, 1)
	mockEq(t, "stake_address", res.Data[0].StakeAddress, mockStake)
	mockEq(t, "reference_script.type", res.Data[0].ReferenceScript.Type, "plutusV1")
}
