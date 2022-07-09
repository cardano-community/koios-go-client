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

func assetsEndpointSpecs() []internal.APITestSpec {
	return []internal.APITestSpec{
		{
			Network:  "mainnet",
			Filename: "assets_endpoint_asset_list",
			Endpoint: "/asset_list",
			Request: internal.APITestRequestSpec{
				Method: "GET",
			},
		},
		{
			Network:  "testnet",
			Filename: "assets_endpoint_asset_list",
			Endpoint: "/asset_list",
			Request: internal.APITestRequestSpec{
				Method: "GET",
			},
		},
		{
			Network:  "mainnet",
			Filename: "assets_endpoint_asset_address_list",
			Endpoint: "/asset_address_list",
			Request: internal.APITestRequestSpec{
				Method: "GET",
				Query: url.Values{
					"_asset_policy": []string{"d3501d9531fcc25e3ca4b6429318c2cc374dbdbcf5e99c1c1e5da1ff"},
					"_asset_name":   []string{"444f4e545350414d"},
				},
			},
		},
		{
			Network:  "testnet",
			Filename: "assets_endpoint_asset_address_list",
			Endpoint: "/asset_address_list",
			Request: internal.APITestRequestSpec{
				Method: "GET",
				Query: url.Values{
					"_asset_policy": []string{"000327a9e427a3a3256eb6212ae26b7f53f7969b8e62d37ea9138a7b"},
					"_asset_name":   []string{"54735465737431"},
				},
			},
		},
		{
			Network:  "mainnet",
			Filename: "assets_endpoint_asset_info",
			Endpoint: "/asset_info",
			Request: internal.APITestRequestSpec{
				Method: "GET",
				Query: url.Values{
					"_asset_policy": []string{"d3501d9531fcc25e3ca4b6429318c2cc374dbdbcf5e99c1c1e5da1ff"},
					"_asset_name":   []string{"444f4e545350414d"},
				},
			},
		},
		{
			Network:  "testnet",
			Filename: "assets_endpoint_asset_info",
			Endpoint: "/asset_info",
			Request: internal.APITestRequestSpec{
				Method: "GET",
				Query: url.Values{
					"_asset_policy": []string{"000327a9e427a3a3256eb6212ae26b7f53f7969b8e62d37ea9138a7b"},
					"_asset_name":   []string{"54735465737431"},
				},
			},
		},
		{
			Network:  "mainnet",
			Filename: "assets_endpoint_asset_history",
			Endpoint: "/asset_history",
			Request: internal.APITestRequestSpec{
				Method: "GET",
				Query: url.Values{
					"_asset_policy": []string{"d3501d9531fcc25e3ca4b6429318c2cc374dbdbcf5e99c1c1e5da1ff"},
					"_asset_name":   []string{"444f4e545350414d"},
				},
			},
		},
		{
			Network:  "testnet",
			Filename: "assets_endpoint_asset_history",
			Endpoint: "/asset_history",
			Request: internal.APITestRequestSpec{
				Method: "GET",
				Query: url.Values{
					"_asset_policy": []string{"000327a9e427a3a3256eb6212ae26b7f53f7969b8e62d37ea9138a7b"},
					"_asset_name":   []string{"54735465737431"},
				},
			},
		},
		{
			Network:  "mainnet",
			Filename: "assets_endpoint_asset_policy_info",
			Endpoint: "/asset_policy_info",
			Request: internal.APITestRequestSpec{
				Method: "GET",
				Query: url.Values{
					"_asset_policy": []string{"a8102151506a8a81dc1763ee05cdd01d787f50dfeb6f843071e1c6a0"},
				},
			},
		},
		{
			Network:  "testnet",
			Filename: "assets_endpoint_asset_policy_info",
			Endpoint: "/asset_policy_info",
			Request: internal.APITestRequestSpec{
				Method: "GET",
				Query: url.Values{
					"_asset_policy": []string{"000327a9e427a3a3256eb6212ae26b7f53f7969b8e62d37ea9138a7b"},
				},
			},
		},
		{
			Network:  "mainnet",
			Filename: "assets_endpoint_asset_summary",
			Endpoint: "/asset_summary",
			Request: internal.APITestRequestSpec{
				Method: "GET",
				Query: url.Values{
					"_asset_policy": []string{"d3501d9531fcc25e3ca4b6429318c2cc374dbdbcf5e99c1c1e5da1ff"},
					"_asset_name":   []string{"444f4e545350414d"},
				},
			},
		},
		{
			Network:  "testnet",
			Filename: "assets_endpoint_asset_summary",
			Endpoint: "/asset_summary",
			Request: internal.APITestRequestSpec{
				Method: "GET",
				Query: url.Values{
					"_asset_policy": []string{"000327a9e427a3a3256eb6212ae26b7f53f7969b8e62d37ea9138a7b"},
					"_asset_name":   []string{"54735465737431"},
				},
			},
		},
		{
			Network:  "mainnet",
			Filename: "assets_endpoint_asset_txs",
			Endpoint: "/asset_txs",
			Request: internal.APITestRequestSpec{
				Method: "GET",
				Query: url.Values{
					"_asset_policy": []string{"d3501d9531fcc25e3ca4b6429318c2cc374dbdbcf5e99c1c1e5da1ff"},
					"_asset_name":   []string{"444f4e545350414d"},
				},
			},
		},
		{
			Network:  "testnet",
			Filename: "assets_endpoint_asset_txs",
			Endpoint: "/asset_txs",
			Request: internal.APITestRequestSpec{
				Method: "GET",
				Query: url.Values{
					"_asset_policy": []string{"000327a9e427a3a3256eb6212ae26b7f53f7969b8e62d37ea9138a7b"},
					"_asset_name":   []string{"54735465737431"},
				},
			},
		},
	}
}
