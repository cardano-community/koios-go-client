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

package koios

import (
	"context"
	"encoding/json"
	"net/url"
)

type (
	// Block defines model for block.
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
		Data []Block `json:"data"`
	}
	// BlockInfoResponse represents response from `/block_info` endpoint.
	BlockInfoResponse struct {
		Response
		Data *Block `json:"data"`
	}
	// BlockTxsHashesResponse represents response from `/block_txs` endpoint.
	BlockTxsHashesResponse struct {
		Response
		Data []TxHash `json:"data"`
	}
)

// GetBlocks returns summarised details about all blocks (paginated - latest first).
func (c *Client) GetBlocks(ctx context.Context) (res *BlocksResponse, err error) {
	res = &BlocksResponse{}
	rsp, err := c.request(ctx, &res.Response, "GET", "/blocks", nil, nil, nil)
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
	res.ready()
	return
}

// GetBlockInfo returns detailed information about a specific block.
func (c *Client) GetBlockInfo(ctx context.Context, hash BlockHash) (res *BlockInfoResponse, err error) {
	res = &BlockInfoResponse{}
	params := url.Values{}
	params.Set("_block_hash", string(hash))

	rsp, err := c.request(ctx, &res.Response, "GET", "/block_info", nil, params, nil)
	if err != nil {
		return
	}
	body, err := readResponseBody(rsp)
	if err != nil {
		res.applyError(nil, err)
		return
	}

	blockpl := []Block{}

	if err = json.Unmarshal(body, &blockpl); err != nil {
		res.applyError(body, err)
		return
	}

	if len(blockpl) == 1 {
		res.Data = &blockpl[0]
	}
	res.ready()
	return res, nil
}

// GetBlockTxHashes returns a list of all transactions hashes
// included in a provided block.
func (c *Client) GetBlockTxHashes(ctx context.Context, hash BlockHash) (res *BlockTxsHashesResponse, err error) {
	res = &BlockTxsHashesResponse{}
	params := url.Values{}
	params.Set("_block_hash", string(hash))

	rsp, err := c.request(ctx, &res.Response, "GET", "/block_txs", nil, params, nil)
	if err != nil {
		res.applyError(nil, err)
		return nil, err
	}

	body, err := readResponseBody(rsp)
	if err != nil {
		res.applyError(body, err)
		return
	}

	blockTxs := []struct {
		Hash TxHash `json:"tx_hash"`
	}{}

	if err = json.Unmarshal(body, &blockTxs); err != nil {
		res.applyError(body, err)
		return
	}

	if len(blockTxs) > 0 {
		for _, tx := range blockTxs {
			res.Data = append(res.Data, tx.Hash)
		}
	}
	res.ready()
	return res, nil
}
