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

const mockAddr = koios.Address("addr1qxkfe8s6m8qt5436lec3f0320hrmpppwqgs2gah4360krvyssntpwjcz303mx3h4avg7p29l3zd8u3jyglmewds9ezrqdc3cxp")

func TestMockGetAddressList(t *testing.T) {
	client, capture := mockClient(t, "address_list")
	res, err := client.GetAddressList(context.Background(), nil)
	mockNoErr(t, err)
	mockExpect(t, capture, "GET", "/address_list", nil, "")
	mockLen(t, "data", res.Data, 1)
	a := res.Data[0]
	mockEq(t, "address", a.Address, mockAddr)
	mockEq(t, "address_hex", a.AddressHex, "60f879af535c02e37bf77fb3b5f51d017229fe47995c37b9aca9e0091c")
	mockEq(t, "has_script", a.HasScript, false)
	mockEq(t, "payment_cred", a.PaymentCred, koios.PaymentCredential("de3c1c527e8826b9cd2030f88f75fc44cd4ce519b9ded9eb794b3794"))
	mockEq(t, "stake_address_hex", a.StakeAddressHex, "e1c865f10d66a2e375242b5ef9f05abc8840d413421982f62c35ce4fbf")
}

func TestMockGetAddressOutputs(t *testing.T) {
	client, capture := mockClient(t, "address_outputs")
	res, err := client.GetAddressOutputs(context.Background(), []koios.Address{mockAddr}, 1000, nil)
	mockNoErr(t, err)
	mockExpect(t, capture, "POST", "/address_outputs", nil,
		`{"_addresses": ["`+string(mockAddr)+`"], "_after_block_height": 1000}`)
	mockLen(t, "data", res.Data, 1)
	mockEq(t, "address", res.Data[0].Address, mockAddr)
	mockEq(t, "value", res.Data[0].Value.String(), "157832856")
	mockEq(t, "is_spent", res.Data[0].IsSpent, true)
	mockEq(t, "block_height", res.Data[0].BlockHeight, uint64(1794506))

	client, capture = mockClient(t, "address_outputs")
	_, err = client.GetAddressOutputs(context.Background(), []koios.Address{mockAddr}, 0, nil)
	mockNoErr(t, err)
	mockExpect(t, capture, "POST", "/address_outputs", nil, `{"_addresses": ["`+string(mockAddr)+`"]}`)

	_, err = client.GetAddressOutputs(context.Background(), nil, 0, nil)
	if !errors.Is(err, koios.ErrNoAddress) {
		t.Errorf("err = %v, want ErrNoAddress", err)
	}
}

func TestMockGetAddressInfoKeepsResponse(t *testing.T) {
	client, capture := mockClient(t, "address_info")
	res, err := client.GetAddressInfo(context.Background(), mockAddr, nil)
	mockNoErr(t, err)
	mockExpect(t, capture, "POST", "/address_info", nil, `{"_addresses": ["`+string(mockAddr)+`"]}`)
	mockEq(t, "status_code", res.StatusCode, 200)
	mockEq(t, "request_method", res.RequestMethod, "POST")
	if res.Data == nil {
		t.Fatal("data is nil")
	}
	mockEq(t, "balance", res.Data.Balance.String(), "10723473983")
	mockEq(t, "script_address", res.Data.ScriptAddress, true)
	mockLen(t, "utxo_set", res.Data.UTxOs, 1)
	mockEq(t, "inline_datum.bytes", res.Data.UTxOs[0].InlineDatum.Bytes, "19029a")
}

func TestMockGetAddressUTxOs(t *testing.T) {
	client, capture := mockClient(t, "address_utxos")
	res, err := client.GetAddressUTxOs(context.Background(), []koios.Address{mockAddr}, true, nil)
	mockNoErr(t, err)
	mockExpect(t, capture, "POST", "/address_utxos", nil, `{"_addresses": ["`+string(mockAddr)+`"], "_extended": true}`)
	mockLen(t, "data", res.Data, 1)
	mockEq(t, "stake_address", res.Data[0].StakeAddress, koios.Address("stake1u8yxtugdv63wxafy9d00nuz6hjyyp4qnggvc9a3vxh8yl0ckml2uz"))
	mockEq(t, "payment_cred", res.Data[0].PaymentCred, koios.PaymentCredential("de3c1c527e8826b9cd2030f88f75fc44cd4ce519b9ded9eb794b3794"))
	mockLen(t, "asset_list", res.Data[0].AssetList, 1)
}
