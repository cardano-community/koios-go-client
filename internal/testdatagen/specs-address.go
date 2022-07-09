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

package main

import (
	"net/url"

	"github.com/cardano-community/koios-go-client/v2/internal"
)

func addressEndpointSpecs() []internal.APITestSpec {
	return []internal.APITestSpec{
		{
			Network:  "mainnet",
			Filename: "address_endpoint_address_info",
			Endpoint: "/address_info",
			Request: internal.APITestRequestSpec{
				Method: "GET",
				Query: url.Values{
					"_address": []string{"addr1qyp9kz50sh9c53hpmk3l4ewj9ur794t2hdqpngsjn3wkc5sztv9glpwt3frwrhdrltjaytc8ut2k4w6qrx3p98zad3fq07xe9g"},
				},
			},
		},
		{
			Network:  "testnet",
			Filename: "address_endpoint_address_info",
			Endpoint: "/address_info",
			Request: internal.APITestRequestSpec{
				Method: "GET",
				Query: url.Values{
					"_address": []string{"addr_test1qpqd4nqjepzdlh9zx5mx5ftnp4hecpgttcprtg0ur3ptpe9efftg48dy58fqwqwvatkn6pj877x858cr7peyr9466jmshglmne"},
				},
			},
		},
		{
			Network:  "mainnet",
			Filename: "address_endpoint_address_txs",
			Endpoint: "/address_txs",
			Request: internal.APITestRequestSpec{
				Method: "POST",
				Body:   []byte("{\"_addresses\": [\"addr1qyp9kz50sh9c53hpmk3l4ewj9ur794t2hdqpngsjn3wkc5sztv9glpwt3frwrhdrltjaytc8ut2k4w6qrx3p98zad3fq07xe9g\",\"addr1qyfldpcvte8nkfpyv0jdc8e026cz5qedx7tajvupdu2724tlj8sypsq6p90hl40ya97xamkm9fwsppus2ru8zf6j8g9sm578cu\"], \"_after_block_height\": 6238675}"),
			},
		},
		{
			Network:  "testnet",
			Filename: "address_endpoint_address_txs",
			Endpoint: "/address_txs",
			Request: internal.APITestRequestSpec{
				Method: "POST",
				Body:   []byte("{\"_addresses\": [\"addr_test1qzx9hu8j4ah3auytk0mwcupd69hpc52t0cw39a65ndrah86djs784u92a3m5w475w3w35tyd6v3qumkze80j8a6h5tuqq5xe8y\",\"addr_test1qrk7920v35zukhcch4kyydy6rxnhqdcvetkvngeqrvtgavw8tpzdklse3kwer7urhrlfg962m9fc8cznfcdpka5pd07sgf8n0w\"], \"_after_block_height\": 2342661}"),
			},
		},
		{
			Network:  "mainnet",
			Filename: "address_endpoint_address_assets",
			Endpoint: "/address_assets",
			Request: internal.APITestRequestSpec{
				Method: "GET",
				Query: url.Values{
					"_address": []string{"addr1q8h22z0n3zqecr9n4q9ysds2m2ms3dqesz575accjpc3jclw55yl8zypnsxt82q2fqmq4k4hpz6pnq9fafm33yr3r93sgnpdw6"},
				},
			},
		},
		{
			Network:  "testnet",
			Filename: "address_endpoint_address_assets",
			Endpoint: "/address_assets",
			Request: internal.APITestRequestSpec{
				Method: "GET",
				Query: url.Values{
					"_address": []string{"addr_test1qzd9gtea50mqv60k3mq9txxtq2ynqwsxcnlx9ltvv3lh0rk9q2x2l6wv8fcr4wpgwmrcwhucsp80ycfw5ensx038hlfsp6lsxj"},
				},
			},
		},
		{
			Network:  "mainnet",
			Filename: "address_endpoint_credential_txs",
			Endpoint: "/credential_txs",
			Request: internal.APITestRequestSpec{
				Method: "POST",
				Body:   []byte("{\"_payment_credentials\": [\"025b0a8f85cb8a46e1dda3fae5d22f07e2d56abb4019a2129c5d6c52\",\"13f6870c5e4f3b242463e4dc1f2f56b02a032d3797d933816f15e555\"], \"_after_block_height\": 6238675}"),
			},
		},
		{
			Network:  "testnet",
			Filename: "address_endpoint_credential_txs",
			Endpoint: "/credential_txs",
			Request: internal.APITestRequestSpec{
				Method: "POST",
				Body:   []byte("{\"_payment_credentials\": [\"00003fac863dc2267d0cd90768c4af653572d719a79ca3b01957fa79\",\"000056d48603bf7daada30c9c175be9c93172d36f82fba0ca972c245\"], \"_after_block_height\": 2342661}"),
			},
		},
	}
}
