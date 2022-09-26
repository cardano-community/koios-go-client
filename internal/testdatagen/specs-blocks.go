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

func blocksEndpointSpecs() []internal.APITestSpec {
	return []internal.APITestSpec{
		// {
		// 	Network:  "mainnet",
		// 	Filename: "blocks_endpoint_blocks",
		// 	Endpoint: "/blocks",
		// 	Request: internal.APITestRequestSpec{
		// 		Method: "GET",
		// 	},
		// },
		// {
		// 	Network:  "testnet",
		// 	Filename: "blocks_endpoint_blocks",
		// 	Endpoint: "/blocks",
		// 	Request: internal.APITestRequestSpec{
		// 		Method: "GET",
		// 	},
		// },
		// {
		// 	Network:  "mainnet",
		// 	Filename: "blocks_endpoint_block_info",
		// 	Endpoint: "/block_info",
		// 	Request: internal.APITestRequestSpec{
		// 		Method: "POST",
		// 		Body:   []byte("{\"_block_hashes\": [\"fb9087c9f1408a7bbd7b022fd294ab565fec8dd3a8ef091567482722a1fa4e30\",\"60188a8dcb6db0d80628815be2cf626c4d17cb3e826cebfca84adaff93ad492a\",\"c6646214a1f377aa461a0163c213fc6b86a559a2d6ebd647d54c4eb00aaab015\"]}"),
		// 	},
		// },
		// {
		// 	Network:  "testnet",
		// 	Filename: "blocks_endpoint_block_info",
		// 	Endpoint: "/block_info",
		// 	Request: internal.APITestRequestSpec{
		// 		Method: "POST",
		// 		Body:   []byte("{\"_block_hashes\": [\"f75fea40852ed7d7f539d008e45255725daef8553ae7162750836f279570813a\",\"ff9f0c7fb1136de2cd6f10c9a140af2887f1d3614cc949bfeb262266d4c202b7\",\"5ef645ee519cde94a82f0aa880048c37978374f248f11e408ac0571a9054d9d3\"]}"),
		// 	},
		// },
		// {
		// 	Network:  "mainnet",
		// 	Filename: "blocks_endpoint_block_txs",
		// 	Endpoint: "/block_txs",
		// 	Request: internal.APITestRequestSpec{
		// 		Method: "GET",
		// 		Query: url.Values{
		// 			"_block_hash": []string{"f6192a1aaa6d3d05b4703891a6b66cd757801c61ace86cbe5ab0d66e07f601ab"},
		// 		},
		// 	},
		// },
		// {
		// 	Network:  "testnet",
		// 	Filename: "blocks_endpoint_block_txs",
		// 	Endpoint: "/block_txs",
		// 	Request: internal.APITestRequestSpec{
		// 		Method: "GET",
		// 		Query: url.Values{
		// 			"_block_hash": []string{"f75fea40852ed7d7f539d008e45255725daef8553ae7162750836f279570813a"},
		// 		},
		// 	},
		// },
	}
}
