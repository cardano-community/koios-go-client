// Copyright 2022 The Howijd.Network Authors
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

package koios

import (
	"context"
	"encoding/json"
	"net/http"
)

type (
	// Asset info.
	Asset struct {
		// Asset Name (hex).
		Name string `json:"asset_name"`

		// Asset Policy ID (hex).
		PolicyID PolicyID `json:"policy_id"`

		// Quantity
		// Input: asset balance on the selected input transaction.
		// Output: sum of assets for output UTxO.
		// Mint: sum of minted assets (negative on burn).
		Quantity Lovelace `json:"quantity"`
	}

	// AssetListItem used to represent response from /asset_list`.
	AssetListItem struct {
		PolicyID   PolicyID `json:"policy_id"`
		AssetNames struct {
			HEX   []string `json:"hex"`
			ASCII []string `json:"ascii"`
		} `json:"asset_names"`
	}

	// AssetListResponse represents response from `/asset_list` endpoint.
	AssetListResponse struct {
		Response
		Data []AssetListItem `json:"response"`
	}
)

// GetTip returns the list of all native assets (paginated).
func (c *Client) GetAssetList(ctx context.Context) (res *AssetListResponse, err error) {
	res = &AssetListResponse{}
	rsp, err := c.request(ctx, &res.Response, "GET", nil, "/asset_list", nil, nil)
	if err != nil {
		res.applyError(nil, err)
		return
	}

	body, err := readResponseBody(rsp)
	if err != nil {
		res.applyError(body, err)
		return
	}

	if err = json.Unmarshal(body, &res.Data); err != nil {
		res.applyError(body, err)
		return
	}

	if rsp.StatusCode != http.StatusOK {
		res.applyError(body, err)
		return
	}
	return res, nil
}
