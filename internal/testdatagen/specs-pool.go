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

func poolsEndpointSpecs() []internal.APITestSpec {
	return []internal.APITestSpec{
		// {
		// 	Network:  "mainnet",
		// 	Filename: "pools_endpoint_pool_list",
		// 	Endpoint: "/pool_list",
		// 	Request: internal.APITestRequestSpec{
		// 		Method: "GET",
		// 	},
		// },
		// {
		// 	Network:  "testnet",
		// 	Filename: "pools_endpoint_pool_list",
		// 	Endpoint: "/pool_list",
		// 	Request: internal.APITestRequestSpec{
		// 		Method: "GET",
		// 	},
		// },
		// {
		// 	Network:  "mainnet",
		// 	Filename: "pools_endpoint_pool_info",
		// 	Endpoint: "/pool_info",
		// 	Request: internal.APITestRequestSpec{
		// 		Method: "POST",
		// 		Body:   []byte("{\"_pool_bech32_ids\": [\"pool100wj94uzf54vup2hdzk0afng4dhjaqggt7j434mtgm8v2gfvfgp\",\"pool102s2nqtea2hf5q0s4amj0evysmfnhrn4apyyhd4azcmsclzm96m\",\"pool102vsulhfx8ua2j9fwl2u7gv57fhhutc3tp6juzaefgrn7ae35wm\"]}"),
		// 	},
		// },
		// {
		// 	Network:  "testnet",
		// 	Filename: "pools_endpoint_pool_info",
		// 	Endpoint: "/pool_info",
		// 	Request: internal.APITestRequestSpec{
		// 		Method: "POST",
		// 		Body:   []byte("{\"_pool_bech32_ids\": [\"pool102llj7e7a0mmxssjvjkv2d6lppuh6cz6q9xwc3tsksn0jqwz9eh\",\"pool102x86jz7uus6p6mlw02fdw2s805kng7g6ujs6s342t5msk36tch\",\"pool103qt58f9xlsr7y9anz3lnyq6cph4xh2yr4qrrtc356ldzz6ktqz\"]}"),
		// 	},
		// },
		// {
		// 	Network:  "mainnet",
		// 	Filename: "pools_endpoint_pool_delegators",
		// 	Endpoint: "/pool_delegators",
		// 	Request: internal.APITestRequestSpec{
		// 		Method: "GET",
		// 		Query: url.Values{
		// 			"_pool_bech32": []string{"pool155efqn9xpcf73pphkk88cmlkdwx4ulkg606tne970qswczg3asc"},
		// 			"_epoch_no":    []string{MainnetEpoch},
		// 		},
		// 	},
		// },
		// {
		// 	Network:  "testnet",
		// 	Filename: "pools_endpoint_pool_delegators",
		// 	Endpoint: "/pool_delegators",
		// 	Request: internal.APITestRequestSpec{
		// 		Method: "GET",
		// 		Query: url.Values{
		// 			"_pool_bech32": []string{"pool102llj7e7a0mmxssjvjkv2d6lppuh6cz6q9xwc3tsksn0jqwz9eh"},
		// 			"_epoch_no":    []string{TestnetEpoch},
		// 		},
		// 	},
		// },
		// {
		// 	Network:  "mainnet",
		// 	Filename: "pools_endpoint_pool_blocks",
		// 	Endpoint: "/pool_blocks",
		// 	Request: internal.APITestRequestSpec{
		// 		Method: "GET",
		// 		Query: url.Values{
		// 			"_pool_bech32": []string{"pool155efqn9xpcf73pphkk88cmlkdwx4ulkg606tne970qswczg3asc"},
		// 			"_epoch_no":    []string{MainnetEpoch},
		// 		},
		// 	},
		// },
		// {
		// 	Network:  "testnet",
		// 	Filename: "pools_endpoint_pool_blocks",
		// 	Endpoint: "/pool_blocks",
		// 	Request: internal.APITestRequestSpec{
		// 		Method: "GET",
		// 		Query: url.Values{
		// 			"_pool_bech32": []string{"pool102llj7e7a0mmxssjvjkv2d6lppuh6cz6q9xwc3tsksn0jqwz9eh"},
		// 			"_epoch_no":    []string{TestnetEpoch},
		// 		},
		// 	},
		// },
		// {
		// 	Network:  "mainnet",
		// 	Filename: "pools_endpoint_pool_history",
		// 	Endpoint: "/pool_history",
		// 	Request: internal.APITestRequestSpec{
		// 		Method: "GET",
		// 		Query: url.Values{
		// 			"_pool_bech32": []string{"pool155efqn9xpcf73pphkk88cmlkdwx4ulkg606tne970qswczg3asc"},
		// 			"_epoch_no":    []string{"320"},
		// 		},
		// 	},
		// },
		// {
		// 	Network:  "testnet",
		// 	Filename: "pools_endpoint_pool_history",
		// 	Endpoint: "/pool_history",
		// 	Request: internal.APITestRequestSpec{
		// 		Method: "GET",
		// 		Query: url.Values{
		// 			"_pool_bech32": []string{"pool102llj7e7a0mmxssjvjkv2d6lppuh6cz6q9xwc3tsksn0jqwz9eh"},
		// 			"_epoch_no":    []string{"185"},
		// 		},
		// 	},
		// },
		// {
		// 	Network:  "mainnet",
		// 	Filename: "pools_endpoint_pool_updates",
		// 	Endpoint: "/pool_updates",
		// 	Request: internal.APITestRequestSpec{
		// 		Method: "GET",
		// 		Query: url.Values{
		// 			"_pool_bech32": []string{"pool155efqn9xpcf73pphkk88cmlkdwx4ulkg606tne970qswczg3asc"},
		// 		},
		// 	},
		// },
		// {
		// 	Network:  "testnet",
		// 	Filename: "pools_endpoint_pool_updates",
		// 	Endpoint: "/pool_updates",
		// 	Request: internal.APITestRequestSpec{
		// 		Method: "GET",
		// 		Query: url.Values{
		// 			"_pool_bech32": []string{"pool102llj7e7a0mmxssjvjkv2d6lppuh6cz6q9xwc3tsksn0jqwz9eh"},
		// 		},
		// 	},
		// },
		// {
		// 	Network:  "mainnet",
		// 	Filename: "pools_endpoint_pool_relays",
		// 	Endpoint: "/pool_relays",
		// 	Request: internal.APITestRequestSpec{
		// 		Method: "GET",
		// 	},
		// },
		// {
		// 	Network:  "testnet",
		// 	Filename: "pools_endpoint_pool_relays",
		// 	Endpoint: "/pool_relays",
		// 	Request: internal.APITestRequestSpec{
		// 		Method: "GET",
		// 	},
		// },
		// {
		// 	Network:  "mainnet",
		// 	Filename: "pools_endpoint_pool_metadata",
		// 	Endpoint: "/pool_metadata",
		// 	Request: internal.APITestRequestSpec{
		// 		Method: "POST",
		// 		Body:   []byte("{\"_pool_bech32_ids\": [\"pool100wj94uzf54vup2hdzk0afng4dhjaqggt7j434mtgm8v2gfvfgp\",\"pool102s2nqtea2hf5q0s4amj0evysmfnhrn4apyyhd4azcmsclzm96m\",\"pool102vsulhfx8ua2j9fwl2u7gv57fhhutc3tp6juzaefgrn7ae35wm\"]}"),
		// 	},
		// },
		// {
		// 	Network:  "testnet",
		// 	Filename: "pools_endpoint_pool_metadata",
		// 	Endpoint: "/pool_metadata",
		// 	Request: internal.APITestRequestSpec{
		// 		Method: "POST",
		// 		Body:   []byte("{\"_pool_bech32_ids\": [\"pool102llj7e7a0mmxssjvjkv2d6lppuh6cz6q9xwc3tsksn0jqwz9eh\",\"pool102x86jz7uus6p6mlw02fdw2s805kng7g6ujs6s342t5msk36tch\",\"pool103qt58f9xlsr7y9anz3lnyq6cph4xh2yr4qrrtc356ldzz6ktqz\"]}"),
		// 	},
		// },
	}
}
