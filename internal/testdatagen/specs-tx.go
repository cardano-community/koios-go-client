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
	"net/http"

	"github.com/cardano-community/koios-go-client/v2/internal"
)

func txEndpointSpecs() []internal.APITestSpec {
	return []internal.APITestSpec{
		{
			Network:  "mainnet",
			Filename: "tx_endpoint_tx_info",
			Endpoint: "/tx_info",
			Request: internal.APITestRequestSpec{
				Method: "POST",
				Body:   []byte("{\"_tx_hashes\": [\"f144a8264acf4bdfe2e1241170969c930d64ab6b0996a4a45237b623f1dd670e\",\"0b8ba3bed976fa4913f19adc9f6dd9063138db5b4dd29cecde369456b5155e94\"]}"),
			},
		},
		{
			Network:  "testnet",
			Filename: "tx_endpoint_tx_info",
			Endpoint: "/tx_info",
			Request: internal.APITestRequestSpec{
				Method: "POST",
				Body:   []byte("{\"_tx_hashes\": [\"928052b80bfc23801da525a6bf8f805da36f22fa0fd5fec2198b0746eb82b72b\",\"c7e96e4cd6aa9e3afbc7b32d1e8023daf4197931f1ea61d2bdfc7a2e5e017cf1\"]}"),
			},
		},
		{
			Network:  "mainnet",
			Filename: "tx_endpoint_tx_utxos",
			Endpoint: "/tx_utxos",
			Request: internal.APITestRequestSpec{
				Method: "POST",
				Body:   []byte("{\"_tx_hashes\": [\"f144a8264acf4bdfe2e1241170969c930d64ab6b0996a4a45237b623f1dd670e\", \"0b8ba3bed976fa4913f19adc9f6dd9063138db5b4dd29cecde369456b5155e94\"]}"),
			},
		},
		{
			Network:  "testnet",
			Filename: "tx_endpoint_tx_utxos",
			Endpoint: "/tx_utxos",
			Request: internal.APITestRequestSpec{
				Method: "POST",
				Body:   []byte("{\"_tx_hashes\": [\"928052b80bfc23801da525a6bf8f805da36f22fa0fd5fec2198b0746eb82b72b\", \"c7e96e4cd6aa9e3afbc7b32d1e8023daf4197931f1ea61d2bdfc7a2e5e017cf1\"]}"),
			},
		},
		{
			Network:  "mainnet",
			Filename: "tx_endpoint_tx_metadata",
			Endpoint: "/tx_metadata",
			Request: internal.APITestRequestSpec{
				Method: "POST",
				Body:   []byte("{\"_tx_hashes\": [\"f144a8264acf4bdfe2e1241170969c930d64ab6b0996a4a45237b623f1dd670e\",\"0b8ba3bed976fa4913f19adc9f6dd9063138db5b4dd29cecde369456b5155e94\"]}"),
			},
		},
		{
			Network:  "testnet",
			Filename: "tx_endpoint_tx_metadata",
			Endpoint: "/tx_metadata",
			Request: internal.APITestRequestSpec{
				Method: "POST",
				Body:   []byte("{\"_tx_hashes\": [\"928052b80bfc23801da525a6bf8f805da36f22fa0fd5fec2198b0746eb82b72b\",\"c7e96e4cd6aa9e3afbc7b32d1e8023daf4197931f1ea61d2bdfc7a2e5e017cf1\"]}"),
			},
		},
		{
			Network:  "mainnet",
			Filename: "tx_endpoint_tx_metalabels",
			Endpoint: "/tx_metalabels",
			Request: internal.APITestRequestSpec{
				Method: "GET",
			},
		},
		{
			Network:  "testnet",
			Filename: "tx_endpoint_tx_metalabels",
			Endpoint: "/tx_metalabels",
			Request: internal.APITestRequestSpec{
				Method: "GET",
			},
		},
		{
			Network:  "mainnet",
			Filename: "tx_endpoint_tx_status",
			Endpoint: "/tx_status",
			Request: internal.APITestRequestSpec{
				Method: "POST",
				Body:   []byte("{\"_tx_hashes\": [\"f144a8264acf4bdfe2e1241170969c930d64ab6b0996a4a45237b623f1dd670e\",\"0b8ba3bed976fa4913f19adc9f6dd9063138db5b4dd29cecde369456b5155e94\"]}"),
			},
		},
		{
			Network:  "testnet",
			Filename: "tx_endpoint_tx_status",
			Endpoint: "/tx_status",
			Request: internal.APITestRequestSpec{
				Method: "POST",
				Body:   []byte("{\"_tx_hashes\": [\"928052b80bfc23801da525a6bf8f805da36f22fa0fd5fec2198b0746eb82b72b\",\"c7e96e4cd6aa9e3afbc7b32d1e8023daf4197931f1ea61d2bdfc7a2e5e017cf1\"]}"),
			},
		},
		{
			Network:  "mainnet",
			Filename: "tx_endpoint_submittx",
			Endpoint: "/submittx",
			Request: internal.APITestRequestSpec{
				Header: http.Header{
					"Content-Type":   []string{"application/cbor"},
					"Content-Length": []string{"585"},
				},
				Method: "POST",
				Body:   []byte("{\"type\":\"Tx AlonzoEra\",\"description\":\"\",\"cborHex\":\"84a60081825820bf9b23cdd9bff2b1a802da7b527a0c6dd0378efa73c0800e8875f9c37930f7ef010d800182825839011f56a82c4c006289171fced204a37a2806e15c88a98872ef9626d3ddc5e778ead6d4d614c64ec8475c8b3dee4d2b8613fa1f3adee95581151a001e848082581d61e1eabc77c631f9dffa24b4c938bf09458d384764ede698d13bb3957f1a00563386021a0002acfd031a0322b0aa0e80a10081825820112bb18afb7f33b90ad1be59accfc7bcc4784c47fde6a5a10d2c932119df16bb584033642286d7805776288655000e2cebbac069def2e1735b91fa53fc5e6650b5921d54c5c5492dc97d8dce9e3539691ca4e45ae9ed4573f6d691adac8aae345001f5f6\"}"),
			},
		},
		{
			Network:  "testnet",
			Filename: "tx_endpoint_submittx",
			Endpoint: "/submittx",
			Request: internal.APITestRequestSpec{
				Header: http.Header{
					"Content-Type":   []string{"application/cbor"},
					"Content-Length": []string{"585"},
				},
				Method: "POST",
				Body:   []byte("{\"type\":\"Tx AlonzoEra\",\"description\":\"\",\"cborHex\":\"84a60081825820bf9b23cdd9bff2b1a802da7b527a0c6dd0378efa73c0800e8875f9c37930f7ef010d800182825839011f56a82c4c006289171fced204a37a2806e15c88a98872ef9626d3ddc5e778ead6d4d614c64ec8475c8b3dee4d2b8613fa1f3adee95581151a001e848082581d61e1eabc77c631f9dffa24b4c938bf09458d384764ede698d13bb3957f1a00563386021a0002acfd031a0322b0aa0e80a10081825820112bb18afb7f33b90ad1be59accfc7bcc4784c47fde6a5a10d2c932119df16bb584033642286d7805776288655000e2cebbac069def2e1735b91fa53fc5e6650b5921d54c5c5492dc97d8dce9e3539691ca4e45ae9ed4573f6d691adac8aae345001f5f6\"}"),
			},
		},
	}
}
