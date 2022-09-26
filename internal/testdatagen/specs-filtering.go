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

func filtersSpecs() []internal.APITestSpec {
	return []internal.APITestSpec{
		// {
		// 	Network:  "mainnet",
		// 	Filename: "filtering_vertical",
		// 	Endpoint: "/blocks",
		// 	Request: internal.APITestRequestSpec{
		// 		Method: "GET",
		// 		Query: url.Values{
		// 			"select": []string{"epoch,epoch_slot,block_height"},
		// 		},
		// 	},
		// },
		// {
		// 	Network:  "testnet",
		// 	Filename: "filtering_vertical",
		// 	Endpoint: "/blocks",
		// 	Request: internal.APITestRequestSpec{
		// 		Method: "GET",
		// 		Query: url.Values{
		// 			"select": []string{"epoch,epoch_slot,block_height"},
		// 		},
		// 	},
		// },
		// {
		// 	Network:  "mainnet",
		// 	Filename: "filtering_horizontal",
		// 	Endpoint: "/blocks",
		// 	Request: internal.APITestRequestSpec{
		// 		Method: "GET",
		// 		Query: url.Values{
		// 			"epoch":      []string{"eq.250"},
		// 			"epoch_slot": []string{"lt.180"},
		// 		},
		// 	},
		// },
		// {
		// 	Network:  "testnet",
		// 	Filename: "filtering_horizontal",
		// 	Endpoint: "/blocks",
		// 	Request: internal.APITestRequestSpec{
		// 		Method: "GET",
		// 		Query: url.Values{
		// 			"epoch":      []string{"eq.250"},
		// 			"epoch_slot": []string{"lt.180"},
		// 		},
		// 	},
		// },
		// {
		// 	Network:  "mainnet",
		// 	Filename: "filtering_pagination_page_1",
		// 	Endpoint: "/blocks",
		// 	Request: internal.APITestRequestSpec{
		// 		Method: "GET",
		// 		Query: url.Values{
		// 			"select": []string{"epoch,epoch_slot,block_height"},
		// 			"order":  []string{"block_height.asc"},
		// 		},
		// 		Header: http.Header{
		// 			"Range": []string{"0-9"},
		// 		},
		// 	},
		// },
		// {
		// 	Network:  "testnet",
		// 	Filename: "filtering_pagination_page_1",
		// 	Endpoint: "/blocks",
		// 	Request: internal.APITestRequestSpec{
		// 		Method: "GET",
		// 		Query: url.Values{
		// 			"select": []string{"epoch,epoch_slot,block_height"},
		// 			"order":  []string{"block_height.asc"},
		// 		},
		// 		Header: http.Header{
		// 			"Range": []string{"0-9"},
		// 		},
		// 	},
		// },
		// {
		// 	Network:  "mainnet",
		// 	Filename: "filtering_pagination_page_2",
		// 	Endpoint: "/blocks",
		// 	Request: internal.APITestRequestSpec{
		// 		Method: "GET",
		// 		Query: url.Values{
		// 			"select": []string{"epoch,epoch_slot,block_height"},
		// 			"order":  []string{"block_height.asc"},
		// 		},
		// 		Header: http.Header{
		// 			"Range": []string{"10-19"},
		// 		},
		// 	},
		// },
		// {
		// 	Network:  "testnet",
		// 	Filename: "filtering_pagination_page_2",
		// 	Endpoint: "/blocks",
		// 	Request: internal.APITestRequestSpec{
		// 		Method: "GET",
		// 		Query: url.Values{
		// 			"select": []string{"epoch,epoch_slot,block_height"},
		// 			"order":  []string{"block_height.asc"},
		// 		},
		// 		Header: http.Header{
		// 			"Range": []string{"10-19"},
		// 		},
		// 	},
		// },
	}
}
