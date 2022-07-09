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

func networkEndpointSpecs() []internal.APITestSpec {
	return []internal.APITestSpec{
		{
			Network:  "mainnet",
			Filename: "network_endpoint_tip",
			Endpoint: "/tip",
			Request: internal.APITestRequestSpec{
				Method: "GET",
			},
		},
		{
			Network:  "testnet",
			Filename: "network_endpoint_tip",
			Endpoint: "/tip",
			Request: internal.APITestRequestSpec{
				Method: "GET",
			},
		},
		{
			Network:  "mainnet",
			Filename: "network_endpoint_genesis",
			Endpoint: "/genesis",
			Request: internal.APITestRequestSpec{
				Method: "GET",
			},
		},
		{
			Network:  "testnet",
			Filename: "network_endpoint_genesis",
			Endpoint: "/genesis",
			Request: internal.APITestRequestSpec{
				Method: "GET",
			},
		},
		{
			Network:  "mainnet",
			Filename: "network_endpoint_totals",
			Endpoint: "/totals",
			Request: internal.APITestRequestSpec{
				Method: "GET",
				Query: url.Values{
					"_epoch_no": []string{MainnetEpoch},
				},
			},
		},
		{
			Network:  "testnet",
			Filename: "network_endpoint_totals",
			Endpoint: "/totals",
			Request: internal.APITestRequestSpec{
				Method: "GET",
				Query: url.Values{
					"_epoch_no": []string{TestnetEpoch},
				},
			},
		},
	}
}
