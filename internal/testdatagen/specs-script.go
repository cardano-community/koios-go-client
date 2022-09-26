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
	"github.com/cardano-community/koios-go-client/v2/internal"
)

func scriptsEndpointSpecs() []internal.APITestSpec {
	return []internal.APITestSpec{
		// {
		// 	Network:  "mainnet",
		// 	Filename: "scripts_endpoint_native_script_list",
		// 	Endpoint: "/native_script_list",
		// 	Request: internal.APITestRequestSpec{
		// 		Method: "GET",
		// 	},
		// },
		// {
		// 	Network:  "testnet",
		// 	Filename: "scripts_endpoint_native_script_list",
		// 	Endpoint: "/native_script_list",
		// 	Request: internal.APITestRequestSpec{
		// 		Method: "GET",
		// 	},
		// },
		// {
		// 	Network:  "mainnet",
		// 	Filename: "scripts_endpoint_plutus_script_list",
		// 	Endpoint: "/plutus_script_list",
		// 	Request: internal.APITestRequestSpec{
		// 		Method: "GET",
		// 	},
		// },
		// {
		// 	Network:  "testnet",
		// 	Filename: "scripts_endpoint_plutus_script_list",
		// 	Endpoint: "/plutus_script_list",
		// 	Request: internal.APITestRequestSpec{
		// 		Method: "GET",
		// 	},
		// },
		// {
		// 	Network:  "mainnet",
		// 	Filename: "scripts_endpoint_script_redeemers",
		// 	Endpoint: "/script_redeemers",
		// 	Request: internal.APITestRequestSpec{
		// 		Method: "GET",
		// 		Query: url.Values{
		// 			"_script_hash": []string{"d8480dc869b94b80e81ec91b0abe307279311fe0e7001a9488f61ff8"},
		// 		},
		// 	},
		// },
		// {
		// 	Network:  "testnet",
		// 	Filename: "scripts_endpoint_script_redeemers",
		// 	Endpoint: "/script_redeemers",
		// 	Request: internal.APITestRequestSpec{
		// 		Method: "GET",
		// 		Query: url.Values{
		// 			"_script_hash": []string{"9a3910acc1e1d49a25eb5798d987739a63f65eb48a78462ffae21e6f"},
		// 		},
		// 	},
		// },
	}
}
