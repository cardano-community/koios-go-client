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

func networkAddresses() []koios.Address {
	var addrs []koios.Address
	switch os.Getenv("KOIOS_NETWORK") {
	case "guild":
		addrs = []koios.Address{
			"addr_test1qzmtfv43a8ncx6ve92ja6yy25npn9raz9pu5a2tfxsqv9gy9ktf0pu6yu4zjh9r37fzx3h4tsxqdjhu3t4d5ffdsfz9s6ska3z",
			"addr_test1vq67g5u8ls4vm4wdvs0r8xvsuej66nvaqedyrj2tcz6tuycz275pu",
		}
	case "testnet":
		addrs = []koios.Address{
			"addr_test1qzx9hu8j4ah3auytk0mwcupd69hpc52t0cw39a65ndrah86djs784u92a3m5w475w3w35tyd6v3qumkze80j8a6h5tuqq5xe8y",
			"addr_test1qrk7920v35zukhcch4kyydy6rxnhqdcvetkvngeqrvtgavw8tpzdklse3kwer7urhrlfg962m9fc8cznfcdpka5pd07sgf8n0w",
		}
	default:
		// mainnet
		addrs = []koios.Address{
			"addr1qyp9kz50sh9c53hpmk3l4ewj9ur794t2hdqpngsjn3wkc5sztv9glpwt3frwrhdrltjaytc8ut2k4w6qrx3p98zad3fq07xe9g",
			"addr1qyfldpcvte8nkfpyv0jdc8e026cz5qedx7tajvupdu2724tlj8sypsq6p90hl40ya97xamkm9fwsppus2ru8zf6j8g9sm578cu",
		}
	}
	return addrs
}

func networkPaymentCredentials() []koios.PaymentCredential {
	var creds []koios.PaymentCredential
	switch os.Getenv("KOIOS_NETWORK") {
	case "guild":
		creds = []koios.PaymentCredential{
			"b6b4b2b1e9e78369992aa5dd108aa4c3328fa228794ea9693400c2a0",
			"35e45387fc2acdd5cd641e339990e665ad4d9d065a41c94bc0b4be13",
		}
	case "testnet":
		creds = []koios.PaymentCredential{
			"00003fac863dc2267d0cd90768c4af653572d719a79ca3b01957fa79",
			"000056d48603bf7daada30c9c175be9c93172d36f82fba0ca972c245",
		}
	default:
		// mainnet
		creds = []koios.PaymentCredential{
			"025b0a8f85cb8a46e1dda3fae5d22f07e2d56abb4019a2129c5d6c52",
			"13f6870c5e4f3b242463e4dc1f2f56b02a032d3797d933816f15e555",
		}
	}
	return creds
}

func networkAccounts() []koios.Address {
	var accs []koios.Address
	switch os.Getenv("KOIOS_NETWORK") {
	case "guild":
		accs = []koios.Address{
			"stake_test17zt9x005zkd2usz2vhvktyzqsuwz25gmgnaqdka5hcj9m2qfg2py2",
			"stake_test1uzzm95hs7dzw23ftj3cly3rgm64crqxet7g46k6y5kcy3zcs3mpjd",
		}
	case "testnet":
		accs = []koios.Address{
			"stake_test1uqrw9tjymlm8wrwq7jk68n6v7fs9qz8z0tkdkve26dylmfc2ux2hj",
			"stake_test1uq7g7kqeucnqfweqzgxk3dw34e8zg4swnc7nagysug2mm4cm77jrx",
		}
	default:
		// mainnet
		accs = []koios.Address{
			"stake1uyfmzu5qqy70a8kq4c8rw09q0w0ktfcxppwujejnsh6tyrg5c774g",
			"stake1uydhlh7f2kkw9eazct5zyzlrvj32gjnkmt2v5qf6t8rut4qwch8ey",
		}
	}
	return accs
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

func assertGreater[V any](t TestingT, count, min V, tag string) bool {
	msg := fmt.Sprintf("%s: should be greater than %v", tag, min)
	return assert.Greater(t, count, min, msg)
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

	if utxo.PaymentAddr != nil {
		assertNotEmpty(t, utxo.PaymentAddr.Bech32, fmt.Sprintf("%s.payment_addr.bech32", tag))
		assertNotEmpty(t, utxo.PaymentAddr.Cred, fmt.Sprintf("%s.payment_addr.cred", tag))
	}
	if utxo.StakeAddress != nil {
		assertNotEmpty(t, utxo.StakeAddress, fmt.Sprintf("%s.stake_addr", tag))
	}
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
