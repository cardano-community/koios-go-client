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
	"testing"

	"github.com/cardano-community/koios-go-client/v2"
	"github.com/stretchr/testify/assert"
)

func TestNativeScriptList(t *testing.T) {
	client, err := getLiveClient()
	if testIsLocal(t, err) {
		return
	}
	nativeScriptListTest(t, client)
}

func nativeScriptListTest(t TestingT, client *koios.Client) {
	opts := client.NewRequestOptions()
	opts.SetPageSize(10)

	res, err := client.GetNativeScripts(context.Background(), opts)
	if !assert.NoError(t, err) {
		return
	}
	assertEqual(t, 10, len(res.Data), "native script returned")

	for _, script := range res.Data {
		assertNotEmpty(t, script.ScriptHash, "script_hash")
		assertNotEmpty(t, script.Type, "type")
		assertNotEmpty(t, script.CreationTxHash, "creation_tx_hash")
	}
}

func TestPlutusScriptList(t *testing.T) {
	client, err := getLiveClient()
	if testIsLocal(t, err) {
		return
	}
	plutusScriptListTest(t, client)
}

func plutusScriptListTest(t TestingT, client *koios.Client) {
	opts := client.NewRequestOptions()
	opts.SetPageSize(10)

	res, err := client.GetPlutusScripts(context.Background(), opts)
	if !assert.NoError(t, err) {
		return
	}
	assertEqual(t, 10, len(res.Data), "plutus script list returned")

	for _, script := range res.Data {
		assertNotEmpty(t, script.ScriptHash, "script_hash")
		assertNotEmpty(t, script.CreationTxHash, "creation_tx_hash")
	}
}

func TestScriptRedeemers(t *testing.T) {
	client, err := getLiveClient()
	if testIsLocal(t, err) {
		return
	}
	scripthash := networkScriptHash()
	scriptRedeemersTest(t, client, scripthash)
}

func scriptRedeemersTest(t TestingT, client *koios.Client, scripthash koios.ScriptHash) {
	res, err := client.GetScriptRedeemers(context.Background(), scripthash, nil)
	if !assert.NoError(t, err) {
		return
	}

	assertNotEmpty(t, res.Data.ScriptHash, "script_hash")
	for _, redeemer := range res.Data.Redeemers {
		assertNotEmpty(t, redeemer.TxHash, "tx_hash")
		assertNotEmpty(t, redeemer.Purpose, "purpose")
		assertNotEmpty(t, redeemer.DatumHash, "datum_hash")
		// assertGreater(t, redeemer.TxIndex,0, "tx_index")
		assertGreater(t, redeemer.UnitMem, 0, "unit_mem")
		assertGreater(t, redeemer.UnitSteps, 0, "unit_steps")
		assertIsPositive(t, redeemer.Fee, "fee")
	}
}

func TestDatumInfo(t *testing.T) {
	client, err := getLiveClient()
	if testIsLocal(t, err) {
		return
	}
	datumhash := networkDatumHash()
	datumInfoTest(t, client, datumhash)
}

func datumInfoTest(t TestingT, client *koios.Client, datumhash koios.DatumHash) {
	res, err := client.GetDatumInfo(context.Background(), datumhash, nil)
	if !assert.NoError(t, err) {
		return
	}

	assertNotEmpty(t, res.Data.Hash, "hash")
	assertNotEmpty(t, res.Data.Bytes, "bytes")
}
