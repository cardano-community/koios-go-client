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

func epochEndpointSpecs() []internal.APITestSpec {
	return []internal.APITestSpec{
		// {
		// 	Network:  "mainnet",
		// 	Filename: "epoch_info",
		// 	Endpoint: "/epoch_info",
		// 	Request: internal.APITestRequestSpec{
		// 		Query: url.Values{
		// 			"_epoch_no": []string{MainnetEpoch},
		// 		},
		// 		Method: "GET",
		// 	},
		// },
		// {
		// 	Network:  "testnet",
		// 	Filename: "epoch_endpoint_epoch_info",
		// 	Endpoint: "/epoch_info",
		// 	Request: internal.APITestRequestSpec{
		// 		Query: url.Values{
		// 			"_epoch_no": []string{TestnetEpoch},
		// 		},
		// 		Method: "GET",
		// 	},
		// },
		// {
		// 	Network:  "mainnet",
		// 	Filename: "epoch_endpoint_epoch_params",
		// 	Endpoint: "/epoch_params",
		// 	Request: internal.APITestRequestSpec{
		// 		Query: url.Values{
		// 			"_epoch_no": []string{MainnetEpoch},
		// 		},
		// 		Method: "GET",
		// 	},
		// },
		// {
		// 	Network:  "testnet",
		// 	Filename: "epoch_endpoint_epoch_params",
		// 	Endpoint: "/epoch_params",
		// 	Request: internal.APITestRequestSpec{
		// 		Query: url.Values{
		// 			"_epoch_no": []string{TestnetEpoch},
		// 		},
		// 		Method: "GET",
		// 	},
		// },
	}
}
