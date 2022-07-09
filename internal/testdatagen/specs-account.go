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

func accountEndpointSpecs() []internal.APITestSpec {
	return []internal.APITestSpec{
		{
			Network:  "mainnet",
			Filename: "account_endpoint_account_list",
			Endpoint: "/account_list",
			Request: internal.APITestRequestSpec{
				Method: "GET",
			},
		},
		{
			Network:  "testnet",
			Filename: "account_endpoint_account_list",
			Endpoint: "/account_list",
			Request: internal.APITestRequestSpec{
				Method: "GET",
			},
		},
		{
			Network:  "mainnet",
			Filename: "account_endpoint_account_info",
			Endpoint: "/account_info",
			Request: internal.APITestRequestSpec{
				Method: "GET",
				Query: url.Values{
					"_address": []string{"stake1u8yxtugdv63wxafy9d00nuz6hjyyp4qnggvc9a3vxh8yl0ckml2uz"},
				},
			},
		},
		{
			Network:  "testnet",
			Filename: "account_endpoint_account_info",
			Endpoint: "/account_info",
			Request: internal.APITestRequestSpec{
				Method: "GET",
				Query: url.Values{
					"_address": []string{"addr_test1qqn5x72yymml6aka0cmkew3jynqgld7xlnwtlsen9ln5tfll0dw5r75vk42mv3ykq8vyjeaanvpytg79xqzymqy5acmqqhx2n7"},
				},
			},
		},
		{
			Network:  "mainnet",
			Filename: "account_endpoint_account_rewards",
			Endpoint: "/account_rewards",
			Request: internal.APITestRequestSpec{
				Method: "GET",
				Query: url.Values{
					"_stake_address": []string{"stake1u8yxtugdv63wxafy9d00nuz6hjyyp4qnggvc9a3vxh8yl0ckml2uz"},
					"_epoch_no":      []string{"320"},
				},
			},
		},
		{
			Network:  "testnet",
			Filename: "account_endpoint_account_rewards",
			Endpoint: "/account_rewards",
			Request: internal.APITestRequestSpec{
				Method: "GET",
				Query: url.Values{
					"_stake_address": []string{"stake_test1uqqzl36c3vk850wk22yqgum9l0upy0y8588hcvsjq9m6j4cxw3qau"},
					"_epoch_no":      []string{"185"},
				},
			},
		},
		{
			Network:  "mainnet",
			Filename: "account_endpoint_account_updates",
			Endpoint: "/account_updates",
			Request: internal.APITestRequestSpec{
				Method: "GET",
				Query: url.Values{
					"_stake_address": []string{"stake1u8yxtugdv63wxafy9d00nuz6hjyyp4qnggvc9a3vxh8yl0ckml2uz"},
				},
			},
		},
		{
			Network:  "testnet",
			Filename: "account_endpoint_account_updates",
			Endpoint: "/account_updates",
			Request: internal.APITestRequestSpec{
				Method: "GET",
				Query: url.Values{
					"_stake_address": []string{"stake_test1uqqzl36c3vk850wk22yqgum9l0upy0y8588hcvsjq9m6j4cxw3qau"},
				},
			},
		},
		{
			Network:  "mainnet",
			Filename: "account_endpoint_account_addresses",
			Endpoint: "/account_addresses",
			Request: internal.APITestRequestSpec{
				Method: "GET",
				Query: url.Values{
					"_address": []string{"stake1u8yxtugdv63wxafy9d00nuz6hjyyp4qnggvc9a3vxh8yl0ckml2uz"},
				},
			},
		},
		{
			Network:  "testnet",
			Filename: "account_endpoint_account_addresses",
			Endpoint: "/account_addresses",
			Request: internal.APITestRequestSpec{
				Method: "GET",
				Query: url.Values{
					"_address": []string{"addr_test1qqn5x72yymml6aka0cmkew3jynqgld7xlnwtlsen9ln5tfll0dw5r75vk42mv3ykq8vyjeaanvpytg79xqzymqy5acmqqhx2n7"},
				},
			},
		},
		{
			Network:  "mainnet",
			Filename: "account_endpoint_account_assets",
			Endpoint: "/account_assets",
			Request: internal.APITestRequestSpec{
				Method: "GET",
				Query: url.Values{
					"_address": []string{"stake1u8yxtugdv63wxafy9d00nuz6hjyyp4qnggvc9a3vxh8yl0ckml2uz"},
				},
			},
		},
		{
			Network:  "testnet",
			Filename: "account_endpoint_account_assets",
			Endpoint: "/account_assets",
			Request: internal.APITestRequestSpec{
				Method: "GET",
				Query: url.Values{
					"_address": []string{"addr_test1qqn5x72yymml6aka0cmkew3jynqgld7xlnwtlsen9ln5tfll0dw5r75vk42mv3ykq8vyjeaanvpytg79xqzymqy5acmqqhx2n7"},
				},
			},
		},
		{
			Network:  "mainnet",
			Filename: "account_endpoint_account_history",
			Endpoint: "/account_history",
			Request: internal.APITestRequestSpec{
				Method: "GET",
				Query: url.Values{
					"_address": []string{"stake1u8yxtugdv63wxafy9d00nuz6hjyyp4qnggvc9a3vxh8yl0ckml2uz"},
				},
			},
		},
		{
			Network:  "testnet",
			Filename: "account_endpoint_account_history",
			Endpoint: "/account_history",
			Request: internal.APITestRequestSpec{
				Method: "GET",
				Query: url.Values{
					"_address": []string{"addr_test1qqn5x72yymml6aka0cmkew3jynqgld7xlnwtlsen9ln5tfll0dw5r75vk42mv3ykq8vyjeaanvpytg79xqzymqy5acmqqhx2n7"},
				},
			},
		},
	}
}
