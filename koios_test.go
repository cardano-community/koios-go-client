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
	"errors"
	"fmt"
	"os"

	"github.com/cardano-community/koios-go-client/v2"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

// TestingT is an interface wrapper around *testing.T
type TestingT interface {
	Errorf(format string, args ...any)
}

func networkEpoch() koios.EpochNo {
	var epoch koios.EpochNo
	switch os.Getenv("KOIOS_NETWORK") {
	case "guild":
		epoch = koios.EpochNo(1950)
	case "testnet":
		epoch = koios.EpochNo(185)
	default:
		// mainnet
		epoch = koios.EpochNo(320)
	}
	return epoch
}

func networkBlockHash() koios.BlockHash {
	var hash koios.BlockHash
	switch os.Getenv("KOIOS_NETWORK") {
	case "guild":
		hash = koios.BlockHash("af2f6f7dd4e4ea6765103a1e38e023da3edd2b3c7fea2aa367222564dbe01cfd")
	case "testnet":
		hash = koios.BlockHash("f75fea40852ed7d7f539d008e45255725daef8553ae7162750836f279570813a")
	default:
		// mainnet
		hash = koios.BlockHash("fb9087c9f1408a7bbd7b022fd294ab565fec8dd3a8ef091567482722a1fa4e30")
	}
	return hash
}

func networkTxHashes() []koios.TxHash {
	var hash []koios.TxHash
	switch os.Getenv("KOIOS_NETWORK") {
	case "guild":
		hash = []koios.TxHash{
			"bf04578d452dd3acb7c70fbac32dc972cb69f932f804171cfb4268f5af0228e7",
			"63b716064012f858450731cb5f960c100c6cb639ec1ec999b898c604451f116a",
		}
	case "testnet":
		hash = []koios.TxHash{
			"928052b80bfc23801da525a6bf8f805da36f22fa0fd5fec2198b0746eb82b72b",
			"c7e96e4cd6aa9e3afbc7b32d1e8023daf4197931f1ea61d2bdfc7a2e5e017cf1",
		}
	default:
		// mainnet
		hash = []koios.TxHash{
			"f144a8264acf4bdfe2e1241170969c930d64ab6b0996a4a45237b623f1dd670e",
			"0b8ba3bed976fa4913f19adc9f6dd9063138db5b4dd29cecde369456b5155e94",
		}
	}
	return hash
}

func getClient() (client *koios.Client, err error) {
	net, ok := os.LookupEnv("KOIOS_NETWORK")
	if !ok {
		return nil, errors.New("KOIOS_NETWORK not set")
	}
	var host string
	switch net {
	case "mainnet":
		host = koios.MainnetHost
	case "guild":
		host = koios.GuildHost
	case "testnet":
		host = koios.TestnetHost
	default:
		return nil, fmt.Errorf("invalid KOIOS_NETWORK=%q", net)
	}
	return koios.New(koios.Host(host))
}

func assertEqual[V comparable](t TestingT, want, got V, tag string) bool {
	msg := fmt.Sprintf("%s: want(%v) got(%v)", tag, want, got)
	return assert.Equal(t, want, got, msg)
}

func assertIsPositive(t TestingT, in decimal.Decimal, tag string) bool {
	msg := fmt.Sprintf("%s: should be positive got  %s", tag, in.String())
	return assert.True(t, in.IsPositive(), msg)
}

func assertGreater[V any](t TestingT, want, got V, tag string) bool {
	msg := fmt.Sprintf("%s: val(%v) should be greater than %v", tag, got, want)
	return assert.Greater(t, want, got, msg)
}

func assertNotEmpty(t TestingT, in any, tag string) bool {
	msg := fmt.Sprintf("%s: in(%v)", tag, in)
	return assert.NotEmpty(t, in, msg)
}

func assertTimeNotZero(t TestingT, in koios.Timestamp, tag string) bool {
	msg := fmt.Sprintf("%s: time is empty or not parsed from return value", tag)
	return assert.False(t, in.Time.IsZero(), msg)
}

func assertEUTxO(t TestingT, eutxo koios.EUTxO, tag string) {
	assertNotEmpty(t, eutxo.TxHash, fmt.Sprintf("eutxo[%s].tx_hash", eutxo.TxHash))
	for i, utxo := range eutxo.Inputs {
		assertUTxO(t, utxo, fmt.Sprintf("eutxo[%s].inputs[%d]", eutxo.TxHash, i))
	}
	for i, utxo := range eutxo.Inputs {
		assertUTxO(t, utxo, fmt.Sprintf("eutxo[%s].outputs[%d]", eutxo.TxHash, i))
	}
}

func assertUTxO(t TestingT, utxo koios.UTxO, tag string) {
	assertNotEmpty(t, utxo.TxHash, fmt.Sprintf("%s.tx_hash", tag))
	assertGreater(t, utxo.TxIndex, -1, fmt.Sprintf("%s.tx_index", tag))
	assertNotEmpty(t, utxo.PaymentAddr.Bech32, fmt.Sprintf("%s.payment_addr.bech32", tag))
	assertNotEmpty(t, utxo.PaymentAddr.Cred, fmt.Sprintf("%s.payment_addr.cred", tag))
	// assertNotEmpty(t, utxo.StakeAddress, fmt.Sprintf("%s.stake_addr", tag))
	// assertGreater(t, utxo.BlockHeight, 0, fmt.Sprintf("%s.block_height", tag))
	// assertTimeNotZero(t, utxo.BlockTime, fmt.Sprintf("%s.block_time", tag))
	assertIsPositive(t, utxo.Value, fmt.Sprintf("%s.value", tag))

	// assertNotEmpty(t, utxo.DatumHash, fmt.Sprintf("%s.datum_hash", tag))

	if utxo.InlineDatum != nil {
		assertNotEmpty(t, "", fmt.Sprintf("%s.inline_datum", tag))
	}
	if utxo.ReferenceScript != nil {
		assertNotEmpty(t, "", fmt.Sprintf("%s.reference_script", tag))
	}
	if len(utxo.AssetList) > 0 {
		for i, asset := range utxo.AssetList {
			assertAsset(t, asset, fmt.Sprintf("%s.asset_list[%d]", tag, i))
		}
	}
}

func assertAsset(t TestingT, asset koios.Asset, tag string) {
	// assertNotEmpty(t, asset.Name, fmt.Sprintf("%s.asset_name", tag))
	// assertNotEmpty(t, asset.Fingerprint, fmt.Sprintf("%s.fingerprint", tag))
	assertNotEmpty(t, asset.PolicyID, fmt.Sprintf("%s.policy_id", tag))
	assertIsPositive(t, asset.Quantity, fmt.Sprintf("%s.quantity", tag))
}

func assertTxMetadata(t TestingT, metadata koios.TxMetadata, tag string) {
	for key, json := range metadata {
		assertNotEmpty(t, key, fmt.Sprintf("%s[%s]", tag, key))
		assertNotEmpty(t, json, fmt.Sprintf("%s[%s]", tag, json))
	}
}

func assertCertificates(t TestingT, certs []koios.Certificate, tag string) {
	for i, cert := range certs {
		// assertGreater(t, cert.Index, 0, fmt.Sprintf("%s[%d].index", tag, i))
		assertNotEmpty(t, cert.Type, fmt.Sprintf("%s[%d].type", tag, i))
		assertNotEmpty(t, cert.Info, fmt.Sprintf("%s[%d].info", tag, i))
	}
}

func assertNativeScripts(t TestingT, nscripts []koios.NativeScript, tag string) {
	for i, nscript := range nscripts {
		assertNotEmpty(t, nscript.CreationTxHash, fmt.Sprintf("%s[%d].creation_tx_hash", tag, i))
		assertNotEmpty(t, nscript.ScriptHash, fmt.Sprintf("%s[%d].script_hash", tag, i))
		assertNotEmpty(t, nscript.Type, fmt.Sprintf("%s[%d].type", tag, i))
		assertNotEmpty(t, nscript.Script.Type, fmt.Sprintf("%s[%d].script.type", tag, i))

		for j, script := range nscript.Script.Scripts {
			assertNotEmpty(t, script, fmt.Sprintf("%s[%d].scripts[%d]", tag, i, j))
		}
	}
}

func assertPlutusContracts(t TestingT, contracts []koios.PlutusContract, tag string) {
	for i, contract := range contracts {
		assertNotEmpty(t, contract.Address, fmt.Sprintf("%s[%d].address", tag, i))
		assertNotEmpty(t, contract.ScriptHash, fmt.Sprintf("%s[%d].script_hash", tag, i))
		assertNotEmpty(t, contract.ByteCode, fmt.Sprintf("%s[%d].bytecode", tag, i))
		assertGreater(t, contract.Size, 0, fmt.Sprintf("%s[%d].size", tag, i))
		assert.True(t, contract.ValidContract, 0, fmt.Sprintf("%s[%d].valid_contract", tag, i))
	}
}
