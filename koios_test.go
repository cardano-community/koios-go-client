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

	"github.com/cardano-community/koios-go-client/v4"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

// TestingT is an interface wrapper around *testing.T
type TestingT interface {
	Errorf(format string, args ...any)
	Error(args ...any)
	SkipNow()
}

var errLocalClient = errors.New("local client is used")

func networkEpoch() koios.EpochNo {
	var epoch koios.EpochNo
	switch os.Getenv("KOIOS_NETWORK") {
	case "guild":
		epoch = koios.EpochNo(1950)
	case "mainnet":
		epoch = koios.EpochNo(320)
	case "preprod":
		epoch = koios.EpochNo(31)
	case "preview":
		epoch = koios.EpochNo(12)
	default:
		// local
		epoch = koios.EpochNo(0)
	}
	return epoch
}

func networkBlockHash() koios.BlockHash {
	var hash koios.BlockHash
	switch os.Getenv("KOIOS_NETWORK") {
	case "guild":
		hash = koios.BlockHash("bddbbc6df0ad09567a513349bafd56d8ec5c8fcd9ee9db12173624b896350d57")
	case "mainnet":
		hash = koios.BlockHash("fb9087c9f1408a7bbd7b022fd294ab565fec8dd3a8ef091567482722a1fa4e30")
	case "preprod":
		hash = koios.BlockHash("2abeb8d1c1227139763be30ddb7a2fd79abd7d44195fca87a7c836a510b2802d")
	case "preview":
		hash = koios.BlockHash("a4504e2495ed03b48be36676f430c54dca0769d29f72ebf18d493abf42d2167b")
	default:
		// mainnet
		hash = koios.BlockHash("")
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
	case "mainnet":
		hash = []koios.TxHash{
			"f144a8264acf4bdfe2e1241170969c930d64ab6b0996a4a45237b623f1dd670e",
			"0b8ba3bed976fa4913f19adc9f6dd9063138db5b4dd29cecde369456b5155e94",
		}
	case "preprod":
		hash = []koios.TxHash{
			"d10133964da9e443b303917fd0b7644ae3d01c133deff85b4f59416c2d00f530",
			"145688d3619e7524510ea64c0ec6363b77a9b8da179ef9bb0273a0940d57d576",
		}
	case "preview":
		hash = []koios.TxHash{
			"f1592b29b79ae85d753913dd25644c60925a4a0683979faa33832fead4b4bd9c",
			"206f6da5b0b0de45605a95f5ce7e172be9674550f7dde3838c45cbf24bab8b00",
		}
	default:
		// local
		hash = []koios.TxHash{}
	}
	return hash
}

func networkPoolID() koios.PoolID {
	switch os.Getenv("KOIOS_NETWORK") {
	case "guild":
		return "pool1xc9eywck4e20tydz4yvh5vfe0ep8whawvwz8wqkc9k046a2ypp4"
	case "mainnet":
		return "pool155efqn9xpcf73pphkk88cmlkdwx4ulkg606tne970qswczg3asc"
	case "preprod":
		return "pool1ext7qrwjzaxcdfhdnkq5mth59ukuu2atcg6tgqpmevpt7ratkta"
	case "preview":
		return "pool1p90428kec03mjdya3k4gv5d20w7lmed7ca0snknef5j977l3y8l"
	default:
		// local
		return ""
	}
}

func networkScriptHash() koios.ScriptHash {
	switch os.Getenv("KOIOS_NETWORK") {
	case "guild":
		return "160301a01ee86d8e46cbe3aef1e3bf69bfa28c65d5be2dde56a37af8"
	case "mainnet":
		return "d8480dc869b94b80e81ec91b0abe307279311fe0e7001a9488f61ff8"
	case "preprod":
		return "590555d7b5760e98ae2bdd29b356247776251dfab0a207bfce98a485"
	case "preview":
		return "f758cf422ca0cbed7d9d6fad1eb5a3c70537d62e661ad450dd2a3810"
	default:
		// local
		return ""
	}
}

func networkDatumHash() koios.DatumHash {
	switch os.Getenv("KOIOS_NETWORK") {
	case "guild":
		return "45b0cfc220ceec5b7c1c62c4d4193d38e4eba48e8815729ce75f9c0ab0e4c1c0"
	case "preprod":
		return "5571e2c3549f15934a38382d1318707a86751fb70827f4cbd29b104480f1be9b"
	case "preview":
		return "6181b3dc623cd8812caf027a3507e9b3095388a7cf3db65983e1fddd3a84c88c"
	case "mainnet":
		return "818ee3db3bbbd04f9f2ce21778cac3ac605802a4fcb00c8b3a58ee2dafc17d46"
	default:
		// local
		return ""
	}
}

func networkAddresses() []koios.Address {
	var addrs []koios.Address
	switch os.Getenv("KOIOS_NETWORK") {
	case "guild":
		addrs = []koios.Address{
			"addr_test1qzmtfv43a8ncx6ve92ja6yy25npn9raz9pu5a2tfxsqv9gy9ktf0pu6yu4zjh9r37fzx3h4tsxqdjhu3t4d5ffdsfz9s6ska3z",
			"addr_test1vq67g5u8ls4vm4wdvs0r8xvsuej66nvaqedyrj2tcz6tuycz275pu",
		}
	case "mainnet":
		addrs = []koios.Address{
			"addr1qyp9kz50sh9c53hpmk3l4ewj9ur794t2hdqpngsjn3wkc5sztv9glpwt3frwrhdrltjaytc8ut2k4w6qrx3p98zad3fq07xe9g",
			"addr1qyfldpcvte8nkfpyv0jdc8e026cz5qedx7tajvupdu2724tlj8sypsq6p90hl40ya97xamkm9fwsppus2ru8zf6j8g9sm578cu",
		}
	case "preprod":
		addrs = []koios.Address{
			"addr_test1vzpwq95z3xyum8vqndgdd9mdnmafh3djcxnc6jemlgdmswcve6tkw",
			"addr_test1vpfwv0ezc5g8a4mkku8hhy3y3vp92t7s3ul8g778g5yegsgalc6gc",
		}
	case "preview":
		addrs = []koios.Address{
			"addr_test1vpfwv0ezc5g8a4mkku8hhy3y3vp92t7s3ul8g778g5yegsgalc6gc",
			"addr_test1vqneq3v0dqh3x3muv6ee3lt8e5729xymnxuavx6tndcjc2cv24ef9",
		}
	default:
		// mainnet
		addrs = []koios.Address{}
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
	case "mainnet":
		creds = []koios.PaymentCredential{
			"025b0a8f85cb8a46e1dda3fae5d22f07e2d56abb4019a2129c5d6c52",
			"13f6870c5e4f3b242463e4dc1f2f56b02a032d3797d933816f15e555",
		}
	case "preprod":
		creds = []koios.PaymentCredential{
			"b429738bd6cc58b5c7932d001aa2bd05cfea47020a556c8c753d4436",
			"82e016828989cd9d809b50d6976d9efa9bc5b2c1a78d4b3bfa1bb83b",
		}
	case "preview":
		creds = []koios.PaymentCredential{
			"33c378cee41b2e15ac848f7f6f1d2f78155ab12d93b713de898d855f",
			"52e63f22c5107ed776b70f7b92248b02552fd08f3e747bc745099441",
		}
	default:
		// local
		creds = []koios.PaymentCredential{}

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
	case "mainnet":
		accs = []koios.Address{
			"stake1uyfmzu5qqy70a8kq4c8rw09q0w0ktfcxppwujejnsh6tyrg5c774g",
			"stake1uydhlh7f2kkw9eazct5zyzlrvj32gjnkmt2v5qf6t8rut4qwch8ey",
		}
	case "preprod":
		accs = []koios.Address{
			"stake_test1urq4rcynzj4uxqc74c852zky7wa6epgmn9r6k3j3gv7502q8jks0l",
			"stake_test1ur4t5nhceyn2amfuj7z74uxmmj8jf9fmgd2egqw8c98ve3cp2g4wx",
		}
	case "preview":
		accs = []koios.Address{
			"stake_test1upv7n2x0lxepkyx8ux2gjt74ecaa39tjgaccxl6hw5fwzngpzf5zt",
			"stake_test1up6wqzrw2h9vvjy5zfkjn0dwtymy5r29zyhf8fyhm6fat9c2am5hl",
		}
	default:
		// local
		accs = []koios.Address{}
	}
	return accs
}

func networkPolicyAsset() (koios.PolicyID, koios.AssetName, int, bool) {
	switch os.Getenv("KOIOS_NETWORK") {
	case "guild":
		return "313534a537bc476c86ff7c57ec511bd7f24a9d15654091b24e9c606e", "41484c636f696e", 63487, true
	case "mainnet":
		return "d3501d9531fcc25e3ca4b6429318c2cc374dbdbcf5e99c1c1e5da1ff", "444f4e545350414d", 63487, true
	case "preprod":
		return "c6e65ba7878b2f8ea0ad39287d3e2fd256dc5c4160fc19bdf4c4d87e", "7447454e53", 50000, true
	case "preview":
		return "065270479316f1d92e00f7f9f095ebeaac9d009c878dc35ce36d3404", "433374", 50000, true
	default:
		// mainnet
		return "", "", 0, false
	}
}

func getLiveClient() (client *koios.Client, err error) {
	net, ok := os.LookupEnv("KOIOS_NETWORK")
	if !ok {
		return nil, fmt.Errorf("%w: KOIOS_NETWORK not set", errLocalClient)
	}
	var host string
	switch net {
	case "mainnet":
		host = koios.MainnetHost
	case "guild":
		host = koios.GuildHost
	case "preprod":
		host = koios.PreProdHost
	case "preview":
		host = koios.PreviewHost
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

func assertCoinNotZero(t TestingT, in decimal.Decimal, tag string) bool {
	msg := fmt.Sprintf("%s: should not be 0", tag)
	return assert.True(t, !in.IsZero(), msg)
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
	if len(metadata) == 0 {
		return
	}
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

func githubActionWarning(title, msg string) {
	fmt.Printf(
		"::warning title=%s::%q\n",
		title,
		msg,
	)
}

func testIsLocal(t TestingT, err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, errLocalClient) {
		t.SkipNow()
		return true
	}
	t.Error(err)
	return false
}
