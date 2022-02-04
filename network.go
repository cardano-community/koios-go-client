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
	"net/url"
)

// GetTip returns the tip info about the latest block seen by chain.
func (c *Client) GetTip(ctx context.Context) (*TipResponse, error) {
	rsp, err := c.GET(ctx, "/tip")
	if err != nil {
		return nil, err
	}
	res := &TipResponse{}
	res.setStatus(rsp)
	body, err := readResponseBody(rsp)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(body, &res.Tip); err != nil {
		return nil, err
	}

	return res, nil
}

// GetGenesis returns the Genesis parameters used to start specific era on chain.
func (c *Client) GetGenesis(ctx context.Context) (*GenesisResponse, error) {
	rsp, err := c.GET(ctx, "/genesis")
	if err != nil {
		return nil, err
	}
	res := &GenesisResponse{}
	res.setStatus(rsp)
	body, err := readResponseBody(rsp)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(body, &res.Genesis); err != nil {
		return nil, err
	}

	return res, nil
}

// GetTotals returns the circulating utxo, treasury, rewards, supply and reserves in
// lovelace for specified epoch, all epochs if empty.
func (c *Client) GetTotals(ctx context.Context, epochNo *EpochNo) (*TotalsResponse, error) {
	params := url.Values{}
	if epochNo != nil {
		params.Set("_epoch_no", string(*epochNo))
	}

	rsp, err := c.GET(ctx, "/totals", params)
	if err != nil {
		return nil, err
	}
	res := &TotalsResponse{}
	res.setStatus(rsp)
	body, err := readResponseBody(rsp)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(body, &res.Totals); err != nil {
		return nil, err
	}

	return res, nil
}
