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

const mockBlockHash = koios.BlockHash("2fa5178c5be950c7f40d6c74efe9b3dd12c4836dc3f3dd052cd1c2a12edd477f")

func TestMockGetBlockInfos(t *testing.T) {
	client, capture := mockClient(t, "block_info")
	res, err := client.GetBlockInfos(context.Background(), []koios.BlockHash{mockBlockHash}, nil)
	mockNoErr(t, err)
	mockExpect(t, capture, "POST", "/block_info", nil, `{"_block_hashes": ["`+string(mockBlockHash)+`"]}`)
	mockLen(t, "data", res.Data, 1)
	b := res.Data[0]
	mockEq(t, "era", b.Era, "Conway")
	mockEq(t, "epoch_no", b.EpochNo, koios.EpochNo(117))
	mockEq(t, "height", b.Height, 1794506)
	mockEq(t, "pool", b.Pool, "pool13m26ky08vz205232k20u8ft5nrg8u68klhn0xfsk9m4gsqsc44v")
	mockEq(t, "total_fees", b.TotalFees.String(), "2346834")
	mockEq(t, "confirmations", b.Confirmations, 664275)
}

func TestMockGetBlockTxCbor(t *testing.T) {
	client, capture := mockClientRaw(t, `[{
		"tx_hash": "f144a8264acf4bdfe2e1241170969c930d64ab6b0996a4a45237b623f1dd670e",
		"block_hash": "2fa5178c5be950c7f40d6c74efe9b3dd12c4836dc3f3dd052cd1c2a12edd477f",
		"block_height": 1794506, "epoch_no": 117, "absolute_slot": 49073930,
		"tx_timestamp": 1704757130, "cbor": "84a40081825820", "valid_contract": true
	}]`)
	res, err := client.GetBlockTxCbor(context.Background(), []koios.BlockHash{mockBlockHash}, nil)
	mockNoErr(t, err)
	mockExpect(t, capture, "POST", "/block_tx_cbor", nil, `{"_block_hashes": ["`+string(mockBlockHash)+`"]}`)
	mockLen(t, "data", res.Data, 1)
	mockEq(t, "tx_hash", res.Data[0].TxHash, koios.TxHash("f144a8264acf4bdfe2e1241170969c930d64ab6b0996a4a45237b623f1dd670e"))
	mockEq(t, "block_hash", res.Data[0].BlockHash, mockBlockHash)
	mockEq(t, "absolute_slot", res.Data[0].AbsoluteSlot, koios.Slot(49073930))
	mockEq(t, "tx_timestamp", res.Data[0].TxTimestamp.Unix(), int64(1704757130))
	mockEq(t, "cbor", res.Data[0].Cbor, "84a40081825820")
	mockEq(t, "valid_contract", res.Data[0].ValidContract, true)

	_, err = client.GetBlockTxCbor(context.Background(), nil, nil)
	if !errors.Is(err, koios.ErrNoBlockHash) {
		t.Errorf("err = %v, want ErrNoBlockHash", err)
	}
}

func TestMockGetBlockTxInfo(t *testing.T) {
	client, capture := mockClient(t, "block_tx_info")
	flags := &koios.TxInfoFlags{Inputs: true, Metadata: true, Governance: true}
	res, err := client.GetBlockTxInfo(context.Background(), []koios.BlockHash{mockBlockHash}, flags, nil)
	mockNoErr(t, err)
	mockExpect(t, capture, "POST", "/block_tx_info", nil,
		`{"_block_hashes": ["`+string(mockBlockHash)+`"], "_inputs": true, "_metadata": true}`)
	mockLen(t, "data", res.Data, 1)
	tx := res.Data[0]
	mockEq(t, "tx_hash", tx.TxHash, koios.TxHash("f144a8264acf4bdfe2e1241170969c930d64ab6b0996a4a45237b623f1dd670e"))
	mockEq(t, "tx_size", tx.TxSize, 391)
	mockEq(t, "treasury_donation", tx.TreasuryDonation.String(), "0")
	mockLen(t, "outputs", tx.Outputs, 1)
	mockEq(t, "payment cred", tx.Outputs[0].PaymentAddr.Cred, koios.PaymentCredential("de3c1c527e8826b9cd2030f88f75fc44cd4ce519b9ded9eb794b3794"))

	client, capture = mockClient(t, "block_tx_info")
	_, err = client.GetBlockTxInfo(context.Background(), []koios.BlockHash{mockBlockHash}, nil, nil)
	mockNoErr(t, err)
	mockExpect(t, capture, "POST", "/block_tx_info", nil, `{"_block_hashes": ["`+string(mockBlockHash)+`"]}`)
}
