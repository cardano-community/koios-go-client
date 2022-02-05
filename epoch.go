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

type (
	// EpochInfo defines model for epoch_info.
	EpochInfo []struct {
		// Rewards accumulated as of given epoch (in lovelaces)
		ActiveStake string `json:"active_stake"`

		// Number of blocks created in epoch
		BlkCount int `json:"blk_count"`

		// Epoch number
		EpochNo int `json:"epoch_no"`

		// Total fees incurred by transactions in epoch
		Fees string `json:"fees"`

		// Timestamp for first block created in epoch
		FirstBlockTime string `json:"first_block_time"`

		// Timestamp for last block created in epoch
		LastBlockTime string `json:"last_block_time"`

		// Total output value across all transactions in epoch
		OutSum string `json:"out_sum"`

		// Number of transactions submitted in epoch
		TxCount int `json:"tx_count"`
	}

	// EpochInfoResponse response of /epoch_info.
	EpochInfoResponse struct {
		Response
		EpochInfo EpochInfo `json:"response"`
	}
)

// Get the epoch information, all epochs if no epoch specified.
func (c *Client) GetEpochInfo(ctx context.Context, epochNo *EpochNo) (*EpochInfoResponse, error) {
	params := url.Values{}
	if epochNo != nil {
		params.Set("_epoch_no", string(*epochNo))
	}

	rsp, err := c.GET(ctx, "/epoch_info", params)
	if err != nil {
		return nil, err
	}
	res := &EpochInfoResponse{}
	res.setStatus(rsp)
	body, err := readResponseBody(rsp)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(body, &res.EpochInfo); err != nil {
		return nil, err
	}

	return res, nil
}
