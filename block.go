// SPDX-License-Identifier: Apache-2.0
//
// Copyright © 2022 The Cardano Community Authors

package koios

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/shopspring/decimal"
)

type (
	// Block defines model for block.
	Block struct {
		// Hash block hash
		Hash BlockHash `json:"hash"`

		// EpochNo number.
		EpochNo EpochNo `json:"epoch_no"`

		// AbsSlot is overall slot number (slots from genesis block of chain).
		AbsSlot Slot `json:"abs_slot"`

		// EpochSlot slot number within epoch.
		EpochSlot Slot `json:"epoch_slot"`

		// Height is block number on chain where transaction was included.
		Height int `json:"block_height"`

		// Size of block.
		Size int `json:"block_size"`

		// Time of the block.
		Time Timestamp `json:"block_time"`

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

		// ProtoMajor is protocol major version
		ProtoMajor int `json:"proto_major,omitempty"`
		// ProtoMinor is protocol minor version
		ProtoMinor int `json:"proto_minor"`

		// TotalOutput output of the block (in lovelace)
		TotalOutput decimal.Decimal `json:"total_output,omitempty"`

		// TotalOutput Total fees of the block (in lovelace)
		TotalFees decimal.Decimal `json:"total_fees,omitempty"`

		// Confirmations is number of confirmations for the block
		Confirmations int `json:"num_confirmations,omitempty"`
	}

	// BlocksResponse represents response from `/blocks` endpoint.
	BlocksResponse struct {
		Response
		Data []Block `json:"data"`
	}
	// BlockInfoResponse represents response from `/block_info` endpoint.
	BlockInfoResponse struct {
		Response
		Data Block `json:"data"`
	}
	// BlockInfoResponse represents response from `/block_info` endpoint.
	BlocksInfoResponse struct {
		Response
		Data []Block `json:"data"`
	}
	// BlockTxsHashesResponse represents response from `/block_txs` endpoint.

	BlockTxs struct {
		BlockHash BlockHash `json:"block_hash"`
		TxHashes  []TxHash  `json:"tx_hashes"`
	}

	BlocksTxsResponse struct {
		Response
		Data []BlockTxs `json:"data"`
	}
	BlockTxsResponse struct {
		Response
		Data BlockTxs `json:"data"`
	}
)

// GetBlocks returns summarised details about all blocks (paginated - latest first).
func (c *Client) GetBlocks(
	ctx context.Context,
	opts *RequestOptions,
) (res *BlocksResponse, err error) {
	res = &BlocksResponse{}
	rsp, err := c.request(ctx, &res.Response, "GET", "/blocks", nil, opts)
	if err != nil {
		return
	}
	err = ReadAndUnmarshalResponse(rsp, &res.Response, &res.Data)
	return
}

// GetBlockInfo returns detailed information about a specific block.
func (c *Client) GetBlockInfo(
	ctx context.Context,
	hash BlockHash,
	opts *RequestOptions,
) (res *BlockInfoResponse, err error) {
	res = &BlockInfoResponse{}
	rsp, err := c.GetBlockInfos(ctx, []BlockHash{hash}, opts)
	res.Response = rsp.Response

	if len(rsp.Data) > 0 {
		res.Data = rsp.Data[0]
	} else {
		err = fmt.Errorf("%w: block_info response was empty", ErrResponse)
	}
	return
}

// GetBlocksInfo returns detailed information about a set of blocks.
func (c *Client) GetBlockInfos(
	ctx context.Context,
	hashes []BlockHash,
	opts *RequestOptions,
) (res *BlocksInfoResponse, err error) {
	res = &BlocksInfoResponse{}
	rsp, err := c.request(ctx, &res.Response, "POST", "/block_info", blockHashesPL(hashes), opts)
	if err != nil {
		return
	}

	err = ReadAndUnmarshalResponse(rsp, &res.Response, &res.Data)
	return
}

// GetBlocksTxs returns a list of all transactions included in a blocks.
func (c *Client) GetBlockTxs(
	ctx context.Context,
	hash BlockHash,
	opts *RequestOptions,
) (res *BlockTxsResponse, err error) {
	res = &BlockTxsResponse{}
	rsp, err := c.GetBlocksTxs(ctx, []BlockHash{hash}, opts)
	res.Response = rsp.Response

	if len(rsp.Data) > 0 {
		res.Data = rsp.Data[0]
	} else {
		err = fmt.Errorf("%w: block %s had no transactions", ErrNoData, hash)
	}
	return
}

// GetBlocksTxs returns a list of all transactions included in a blocks.
func (c *Client) GetBlocksTxs(
	ctx context.Context,
	hashes []BlockHash,
	opts *RequestOptions,
) (res *BlocksTxsResponse, err error) {
	res = &BlocksTxsResponse{}
	rsp, err := c.request(ctx, &res.Response, "POST", "/block_txs", blockHashesPL(hashes), opts)
	if err != nil {
		return
	}
	err = ReadAndUnmarshalResponse(rsp, &res.Response, &res.Data)
	return
}

func blockHashesPL(bhash []BlockHash) io.Reader {
	var payload = struct {
		BlockHashes []BlockHash `json:"_block_hashes"`
	}{bhash}
	rpipe, w := io.Pipe()
	go func() {
		_ = json.NewEncoder(w).Encode(payload)
		defer w.Close()
	}()
	return rpipe
}

// handle api json tags Block.epoch and Block.epoch_no.
// SEE: https://github.com/cardano-community/koios-artifacts/issues/102
func (block *Block) UnmarshalJSON(b []byte) error {
	type B Block
	if err := json.Unmarshal(b, (*B)(block)); err != nil {
		return err
	}
	if block.EpochNo == 0 {
		var fix = struct {
			EpochNo EpochNo `json:"epoch"`
		}{}
		if err := json.Unmarshal(b, &fix); err != nil {
			return err
		}
		block.EpochNo = fix.EpochNo
	}
	return nil
}
