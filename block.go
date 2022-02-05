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
)

type (
	// Totals defines model for block.
	Block struct {
		// BlockHash in hex format.
		BlockHash `json:"hash"`

		// Epoch number.
		Epoch EpochNo `json:"epoch"`

		// AbsoluteSlot is overall slot number (slots from genesis block of chain).
		AbsoluteSlot int `json:"abs_slot"`

		// EpochSlot slot number within epoch.
		EpochSlot int `json:"epoch_slot"`

		// Height is block number on chain where transaction was included.
		Height int `json:"height"`

		// Size of block.
		Size int `json:"size"`

		// BlockTime time of the block.
		BlockTime string `json:"block_time"`

		// TxCount transactions count in block.
		TxCount int `json:"tx_count"`

		// VrfKey is pool VRF key.
		VrfKey string `json:"vrf_key"`

		// Pool ID.
		Pool string `json:"pool"`

		// OpCertCounter is pool latest operational certificate counter value.
		OpCertCounter int `json:"op_cert_counter"`
	}

	// BlocksResponse represents response from `/blocks` enpoint.
	BlocksResponse struct {
		Response
		Blocks []Block `json:"response"`
	}
)

// GetBlocks returns summarised details about all blocks (paginated - latest first).
func (c *Client) GetBlocks(ctx context.Context) (*BlocksResponse, error) {
	rsp, err := c.GET(ctx, "/blocks")
	if err != nil {
		return nil, err
	}
	res := &BlocksResponse{}
	res.setStatus(rsp)
	body, err := readResponseBody(rsp)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(body, &res.Blocks); err != nil {
		res.applyError(body, err)
		return res, nil
	}

	return res, nil
}
