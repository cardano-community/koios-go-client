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
	"net/url"
)

type (
	// Totals defines model for block.
	Block struct {
		// Hash block hash
		Hash BlockHash `json:"hash"`

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

		// Time time of the block.
		Time string `json:"block_time"`

		// TxCount transactions count in block.
		TxCount int `json:"tx_count"`

		// VrfKey is pool VRF key.
		VrfKey string `json:"vrf_key"`

		// OpCert latest ool operational certificate hash
		OpCert string `json:"op_cert,omitempty"`

		// Pool ID.
		Pool string `json:"pool"`

		// OpCertCounter is pool latest operational certificate counter value.
		OpCertCounter int `json:"op_cert_counter"`

		// ParentHash parent block hash
		ParentHash BlockHash `json:"parent_hash,omitempty"`

		// ChildHash child block hash
		ChildHash BlockHash `json:"child_hash,omitempty"`
	}

	// BlocksResponse represents response from `/blocks` endpoint.
	BlocksResponse struct {
		Response
		Blocks []Block `json:"response,omitempty"`
	}
	// BlockInfoResponse represents response from `/block_info` endpoint.
	BlockInfoResponse struct {
		Response
		Block *Block `json:"response,omitempty"`
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

// GetBlockInfo returns detailed information about a specific block
func (c *Client) GetBlockInfo(ctx context.Context, hash BlockHash) (*BlockInfoResponse, error) {
	params := url.Values{}
	params.Set("_block_hash", string(hash))

	rsp, err := c.GET(ctx, "/block_info", params)
	if err != nil {
		return nil, err
	}
	res := &BlockInfoResponse{}

	res.setStatus(rsp)
	body, err := readResponseBody(rsp)
	if err != nil {
		return nil, err
	}

	blockpl := []Block{}

	if err := json.Unmarshal(body, &blockpl); err != nil {
		res.applyError(body, err)
		return res, nil
	}

	if rsp.StatusCode != http.StatusOK {
		res.applyError(body, err)
		return res, nil
	}
	if len(blockpl) == 1 {
		res.Block = &blockpl[0]
	}
	return res, nil
}
