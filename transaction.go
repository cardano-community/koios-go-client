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
	"io"
)

type (
	// UTxO model holds inputs and outputs for given UTxO.
	UTxO struct {
		/// TxHash is hash of transaction.
		TxHash TxHash `json:"tx_hash"`

		// Inputs An array with details about inputs used in a transaction.
		Inputs []TxInput `json:"inputs"`
		// Outputs An array with details about outputs from the transaction.
		Outputs []TxOutput `json:"outputs"`
	}

	TxInput struct {
		// An array of assets contained on input UTxO.
		AssetList []Asset `json:"asset_list,omitempty"`

		// input UTxO.
		PaymentAddr PaymentAddr `json:"payment_addr"`

		// StakeAddress for transaction's input UTxO.
		StakeAddress StakeAddress `json:"stake_addr,omitempty"`

		// Hash of Transaction for input UTxO.
		TxHash TxHash `json:"tx_hash"`

		// Index of input UTxO on the mentioned address used for input.
		TxIndex int `json:"tx_index"`

		// Balance on the selected input transaction.
		Value Lovelace `json:"value"`
	}

	TxOutput struct {
		// An array of assets to be included in output UTxO.
		AssetList []Asset `json:"asset_list,omitempty"`

		// where funds were sent or change to be returned.
		PaymentAddr PaymentAddr `json:"payment_addr"`

		// StakeAddress for transaction's output UTxO.
		StakeAddress StakeAddress `json:"stake_addr,omitempty"`

		// Hash of this transaction.
		TxHash TxHash `json:"tx_hash"`

		// Index of output UTxO.
		TxIndex int `json:"tx_index"`

		// Total sum on the output address.
		Value Lovelace `json:"value"`
	}

	TxMetadata struct {
		// JSON containing details about metadata within transaction.
		JSON map[string]interface{} `json:"json"`

		// Key is metadata (index).
		Key int `json:"key"`
	}

	TxsWithdrawal struct {
		// Amount is withdrawal amount in lovelaces.
		Amount Lovelace `json:"amount,omitempty"`
		// StakeAddress fo withdrawal.
		StakeAddress StakeAddress `json:"stake_addr,omitempty"`
	}

	TxInfo struct {
		// TxHash is hash of transaction.
		TxHash TxHash `json:"tx_hash"`

		// BlockHash is hash of the block in which transaction was included.
		BlockHash BlockHash `json:"block_hash"`

		// BlockHeight is block number on chain where transaction was included.
		BlockHeight int `json:"block_height"`

		// Epoch number.
		Epoch EpochNo `json:"epoch"`

		// EpochSlot is slot number within epoch.
		EpochSlot int `json:"epoch_slot"`

		// AbsoluteSlot is overall slot number (slots from genesis block of chain).
		AbsoluteSlot int `json:"absolute_slot"`

		// TxTimestamp is timestamp when block containing transaction was created.
		TxTimestamp string `json:"tx_timestamp"`

		// TxBlockIndex is index of transaction within block.
		TxBlockIndex int `json:"tx_block_index"`

		// TxSize is transaction size in bytes.
		TxSize int `json:"tx_size"`

		// TotalOutput is total sum of all transaction outputs (in lovelaces).
		TotalOutput Lovelace `json:"total_output"`

		// Fee is total transaction fee (in lovelaces).
		Fee Lovelace `json:"fee"`

		// Deposit is total deposits included in transaction (for example,
		// if it is registering a pool/key).
		Deposit Lovelace `json:"deposit"`

		// InvalidAfter is slot number after which transaction cannot be validated.
		InvalidAfter int `json:"invalid_after,omitempty"`

		// InvalidBefore is slot number before which transaction cannot be validated.
		// (if supplied, else 0)
		InvalidBefore int `json:"invalid_before,omitempty"`

		// Inputs An array with details about inputs used in a transaction
		Inputs []TxInput `json:"inputs"`

		// Outputs An array with details about outputs from the transaction.
		Outputs []TxOutput `json:"outputs,omitempty"`

		// AssetsMinted An array of minted assets with-in a transaction (if any).
		AssetsMinted []Asset `json:"assets_minted,omitempty"`

		// Collaterals An array of collateral inputs needed when dealing with smart contracts.
		Collaterals []TxInput `json:"collaterals,omitempty"`

		// Metadata present with-in a transaction (if any)
		Metadata []TxMetadata `json:"metadata,omitempty"`

		// Array of withdrawals with-in a transaction (if any)
		Withdrawals []TxsWithdrawal `json:"withdrawals,omitempty"`

		// Certificates present with-in a transaction (if any)
		Certificates []Certificate `json:"certificates,omitempty"`
	}

	// TxsInfoResponse represents response from `/tx_info` endpoint.
	TxsInfosResponse struct {
		Response
		TXS []TxInfo `json:"response,omitempty"`
	}

	// TxsInfoResponse represents response from `/tx_info` endpoint.
	// when requesting info about single transaction.
	TxInfoResponse struct {
		Response
		TX TxInfo `json:"response,omitempty"`
	}

	// TxUTxOsResponse represents response from `/tx_utxos` endpoint.
	TxUTxOsResponse struct {
		Response
		UTxOs []UTxO `json:"response,omitempty"`
	}
)

// GetTxInfo returns detailed information about transaction.
func (c *Client) GetTxInfo(ctx context.Context, tx TxHash) (res *TxInfoResponse, err error) {
	res = &TxInfoResponse{}
	rsp, err := c.GetTxsInfos(ctx, []TxHash{tx})
	res.Response = rsp.Response
	if len(rsp.TXS) == 1 {
		res.TX = rsp.TXS[0]
	}
	return
}

// GetTxsInfos returns detailed information about transaction(s).
func (c *Client) GetTxsInfos(ctx context.Context, txs []TxHash) (res *TxsInfosResponse, err error) {
	res = &TxsInfosResponse{}
	if len(txs) == 0 {
		err = ErrNoTxHash
		res.applyError(nil, err)
		return
	}

	rsp, err := c.request(ctx, &res.Response, "POST", txHashesPL(txs), "/tx_info")
	if err != nil {
		res.applyError(nil, err)
		return
	}

	body, err := readResponseBody(rsp)
	if err != nil {
		res.applyError(body, err)
		return
	}
	if err = json.Unmarshal(body, &res.TXS); err != nil {
		res.applyError(body, err)
		return
	}
	res.ready()
	return res, nil
}

// GetTxsUTxOs returns UTxO set (inputs/outputs) of transactions.
func (c *Client) GetTxsUTxOs(ctx context.Context, txs []TxHash) (res *TxUTxOsResponse, err error) {
	res = &TxUTxOsResponse{}
	if len(txs) == 0 {
		err = ErrNoTxHash
		res.applyError(nil, err)
		return
	}

	rsp, err := c.request(ctx, &res.Response, "POST", txHashesPL(txs), "/tx_utxos")
	if err != nil {
		res.applyError(nil, err)
		return
	}

	body, err := readResponseBody(rsp)
	if err != nil {
		res.applyError(body, err)
		return
	}
	if err = json.Unmarshal(body, &res.UTxOs); err != nil {
		res.applyError(body, err)
		return
	}
	res.ready()
	return res, nil
}

func txHashesPL(txs []TxHash) io.Reader {
	var payload = struct {
		TxHashes []TxHash `json:"_tx_hashes"`
	}{txs}
	rpipe, w := io.Pipe()
	go func() {
		_ = json.NewEncoder(w).Encode(payload)
		w.Close()
	}()
	return rpipe
}
