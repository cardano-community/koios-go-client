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
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"

	"github.com/fxamacker/cbor/v2"
)

type (
	TX struct {
		TxInfo
	}

	// UTxO model holds inputs and outputs for given UTxO.
	UTxO struct {
		/// TxHash is hash of transaction.
		TxHash TxHash `json:"tx_hash,omitempty"`

		// Inputs An array with details about inputs used in a transaction.
		Inputs []TxInput `json:"inputs" cbor:"0,keyasint"`
		// Outputs An array with details about outputs from the transaction.
		Outputs []TxOutput `json:"outputs" cbor:"1,keyasint"`
	}

	// TxMetalabel defines model for tx_metalabels.
	TxMetalabel struct {
		// A distinct known metalabel
		Metalabel string `json:"metalabel"`
	}

	// TxInput an transaxtion input.
	TxInput struct {
		// An array of assets contained on input UTxO.
		AssetList []Asset `json:"asset_list,omitempty"`

		// input UTxO.
		PaymentAddr PaymentAddr `json:"payment_addr,omitempty"`

		// StakeAddress for transaction's input UTxO.
		StakeAddress StakeAddress `json:"stake_addr,omitempty"`

		// Hash of Transaction for input UTxO.
		TxHash TxHash `json:"tx_hash"`

		// Index of input UTxO on the mentioned address used for input.
		TxIndex uint32 `json:"tx_index"`

		// Balance on the selected input transaction.
		Value Lovelace `json:"value"`
	}

	// TxOutput an transaxtion output.
	TxOutput struct {
		// An array of assets to be included in output UTxO.
		AssetList []Asset `json:"asset_list,omitempty"`

		// where funds were sent or change to be returned.
		PaymentAddr PaymentAddr `json:"payment_addr,omitempty"`

		// StakeAddress for transaction's output UTxO.
		StakeAddress StakeAddress `json:"stake_addr,omitempty"`

		// Hash of this transaction.
		TxHash TxHash `json:"tx_hash"`

		// Index of output UTxO.
		TxIndex uint32 `json:"tx_index"`

		// Total sum on the output address.
		Value Lovelace `json:"value"`
	}

	// TxInfoMetadata metadata in transaction info.
	TxInfoMetadata struct {
		// JSON containing details about metadata within transaction.
		JSON map[string]interface{} `json:"json"`

		// Key is metadata (index).
		Key string `json:"key"`
	}

	// TxsWithdrawal withdrawal record in transaction.
	TxsWithdrawal struct {
		// Amount is withdrawal amount in lovelaces.
		Amount Lovelace `json:"amount,omitempty"`
		// StakeAddress fo withdrawal.
		StakeAddress StakeAddress `json:"stake_addr,omitempty"`
	}

	// TxInfo transaction info.
	TxInfo struct {
		UTxO

		// BlockHash is hash of the block in which transaction was included.
		BlockHash BlockHash `json:"block_hash,omitempty"`

		// BlockHeight is block number on chain where transaction was included.
		BlockHeight uint64 `json:"block_height,omitempty"`

		// Epoch number.
		Epoch EpochNo `json:"epoch,omitempty"`

		// EpochSlot is slot number within epoch.
		EpochSlot uint32 `json:"epoch_slot,omitempty"`

		// AbsoluteSlot is overall slot number (slots from genesis block of chain).
		AbsoluteSlot uint64 `json:"absolute_slot,omitempty"`

		// TxTimestamp is timestamp when block containing transaction was created.
		TxTimestamp string `json:"tx_timestamp,omitempty"`

		// TxBlockIndex is index of transaction within block.
		TxBlockIndex uint32 `json:"tx_block_index,omitempty"`

		// TxSize is transaction size in bytes.
		TxSize uint32 `json:"tx_size,omitempty"`

		// TotalOutput is total sum of all transaction outputs (in lovelaces).
		TotalOutput Lovelace `json:"total_output,omitempty"`

		// Fee is total transaction fee (in lovelaces).
		Fee Lovelace `json:"fee,omitempty" cbor:"2,keyasint"`

		// Deposit is total deposits included in transaction (for example,
		// if it is registering a pool/key).
		Deposit Lovelace `json:"deposit,omitempty"`

		// InvalidAfter is slot number after which transaction cannot be validated.
		InvalidAfter uint64 `json:"invalid_after,omitempty" cbor:"3,keyasint,omitempty"`

		// InvalidBefore is slot number before which transaction cannot be validated.
		// (if supplied, else 0)
		InvalidBefore uint64 `json:"invalid_before,omitempty" cbor:"8,keyasint,omitempty"`

		// AssetsMinted An array of minted assets with-in a transaction (if any).
		AssetsMinted []Asset `json:"assets_minted,omitempty"`

		// Collaterals An array of collateral inputs needed when dealing with smart contracts.
		Collaterals []TxInput `json:"collaterals,omitempty"`

		// Metadata present with-in a transaction (if any)
		Metadata []TxInfoMetadata `json:"metadata,omitempty"`

		// Array of withdrawals with-in a transaction (if any)
		Withdrawals []TxsWithdrawal `json:"withdrawals,omitempty"`

		// Certificates present with-in a transaction (if any)
		Certificates []Certificate `json:"certificates,omitempty"`
	}

	// TxsInfosResponse represents response from `/tx_info` endpoint.
	TxsInfosResponse struct {
		Response
		Data []TxInfo `json:"response"`
	}

	// TxInfoResponse represents response from `/tx_info` endpoint.
	// when requesting info about single transaction.
	TxInfoResponse struct {
		Response
		Data *TxInfo `json:"response"`
	}

	// TxUTxOsResponse represents response from `/tx_utxos` endpoint.
	TxUTxOsResponse struct {
		Response
		Data *UTxO `json:"data"`
	}

	// TxsUTxOsResponse represents response from `/tx_utxos` endpoint.
	TxsUTxOsResponse struct {
		Response
		Data []UTxO `json:"data"`
	}

	// TxMetadata transaction metadata lookup res for `/tx_metadata` endpoint.
	TxMetadata struct {
		// TxHash is hash of transaction.
		TxHash TxHash `json:"tx_hash"`
		// Metadata present with-in a transaction (if any)
		Metadata map[string]interface{} `json:"metadata"`
	}

	// SubmitSignedTxResponse represents response from `/submittx` endpoint.
	SubmitSignedTxResponse struct {
		Response
		Data TxHash `json:"data"`
	}

	// TxBodyJSON used to Unmarshal built transactions.
	TxBodyJSON struct {
		Type        string `json:"type"`
		Description string `json:"description"`
		CborHex     string `json:"cborHex"`
	}
	// TxMetadataResponse represents response from `/tx_metadata` endpoint.
	TxMetadataResponse struct {
		Response
		Data *TxMetadata `json:"data"`
	}

	// TxsMetadataResponse represents response from `/tx_metadata` endpoint.
	TxsMetadataResponse struct {
		Response
		Data []TxMetadata `json:"data"`
	}

	// TxMetaLabelsResponse represents response from `/tx_metalabels` endpoint.
	TxMetaLabelsResponse struct {
		Response
		Data []TxMetalabel `json:"data"`
	}

	// TxStatus is tx_status enpoint response.
	TxStatus struct {
		TxHash        TxHash `json:"tx_hash"`
		Confirmations uint64 `json:"num_confirmations"`
	}

	// TxsStatusesResponse represents response from `/tx_status` endpoint.
	TxsStatusesResponse struct {
		Response
		Data []TxStatus `json:"response"`
	}

	// TxStatusResponse represents response from `/tx_status` endpoint.
	TxStatusResponse struct {
		Response
		Data *TxStatus `json:"response"`
	}
)

// GetTxInfo returns detailed information about transaction.
func (c *Client) GetTxInfo(
	ctx context.Context,
	tx TxHash,
	opts *RequestOptions,
) (res *TxInfoResponse, err error) {
	res = &TxInfoResponse{}
	rsp, err := c.GetTxsInfo(ctx, []TxHash{tx}, opts)
	res.Response = rsp.Response
	if len(rsp.Data) > 0 {
		res.Data = &rsp.Data[0]
	}
	return
}

// GetTxsInfo returns detailed information about transaction(s).
func (c *Client) GetTxsInfo(
	ctx context.Context,
	txs []TxHash,
	opts *RequestOptions,
) (*TxsInfosResponse, error) {
	res := &TxsInfosResponse{}
	if len(txs) == 0 || len(txs[0]) == 0 {
		err := ErrNoTxHash
		res.applyError(nil, err)
		return res, err
	}

	rsp, err := c.request(ctx, &res.Response, "POST", "/tx_info", txHashesPL(txs), opts)
	if err != nil {
		return res, err
	}
	return res, ReadAndUnmarshalResponse(rsp, &res.Response, &res.Data)
}

// GetTxUTxOs returns UTxO set (inputs/outputs) of transaction.
func (c *Client) GetTxUTxOs(
	ctx context.Context,
	tx TxHash,
	opts *RequestOptions,
) (res *TxUTxOsResponse, err error) {
	res = &TxUTxOsResponse{}
	rsp, err := c.GetTxsUTxOs(ctx, []TxHash{tx}, opts)
	res.Response = rsp.Response
	if len(rsp.Data) > 0 {
		res.Data = &rsp.Data[0]
	}
	return
}

// GetTxsUTxOs returns UTxO set (inputs/outputs) of transactions.
func (c *Client) GetTxsUTxOs(
	ctx context.Context,
	txs []TxHash,
	opts *RequestOptions,
) (*TxsUTxOsResponse, error) {
	res := &TxsUTxOsResponse{}
	if len(txs) == 0 || len(txs[0]) == 0 {
		err := ErrNoTxHash
		res.applyError(nil, err)
		return res, err
	}

	rsp, err := c.request(ctx, &res.Response, "POST", "/tx_utxos", txHashesPL(txs), opts)
	if err != nil {
		return res, err
	}

	return res, ReadAndUnmarshalResponse(rsp, &res.Response, &res.Data)
}

// GetTxMetadata returns metadata information (if any) for given transaction.
func (c *Client) GetTxMetadata(
	ctx context.Context,
	tx TxHash,
	opts *RequestOptions,
) (res *TxMetadataResponse, err error) {
	res = &TxMetadataResponse{}
	rsp, err := c.GetTxsMetadata(ctx, []TxHash{tx}, opts)
	res.Response = rsp.Response
	if len(rsp.Data) > 0 {
		res.Data = &rsp.Data[0]
	}
	return
}

// GetTxsMetadata returns metadata for requested transaction(s).
func (c *Client) GetTxsMetadata(
	ctx context.Context,
	txs []TxHash,
	opts *RequestOptions,
) (*TxsMetadataResponse, error) {
	res := &TxsMetadataResponse{}
	if len(txs) == 0 || len(txs[0]) == 0 {
		err := ErrNoTxHash
		res.applyError(nil, err)
		return res, err
	}

	rsp, err := c.request(ctx, &res.Response, "POST", "/tx_metadata", txHashesPL(txs), opts)
	if err != nil {
		return res, err
	}

	return res, ReadAndUnmarshalResponse(rsp, &res.Response, &res.Data)
}

// GetTxMetaLabels retruns a list of all transaction metalabels.
func (c *Client) GetTxMetaLabels(
	ctx context.Context,
	opts *RequestOptions,
) (*TxMetaLabelsResponse, error) {
	res := &TxMetaLabelsResponse{}
	rsp, err := c.request(ctx, &res.Response, "GET", "/tx_metalabels", nil, opts)
	if err != nil {
		return res, err
	}
	return res, ReadAndUnmarshalResponse(rsp, &res.Response, &res.Data)
}

// SubmitSignedTx Submit an transaction to the network.
func (c *Client) SubmitSignedTx(
	ctx context.Context,
	stx TxBodyJSON,
	opts *RequestOptions,
) (*SubmitSignedTxResponse, error) {
	res := &SubmitSignedTxResponse{}

	var method = "POST"
	cborb, err := hex.DecodeString(stx.CborHex)
	if err != nil {
		res.RequestMethod = method
		res.StatusCode = 400
		res.applyError(nil, err)
		return res, err
	}

	if opts == nil {
		opts = c.NewRequestOptions()
	}
	opts.HeadersSet("Content-Type", "application/cbor")
	opts.HeadersSet("Content-Length", fmt.Sprint(len(cborb)))

	rsp, err := c.request(ctx, &res.Response, method, "/submittx", bytes.NewBuffer(cborb), opts)
	if err != nil {
		body, _ := ReadResponseBody(rsp)
		res.applyError(body, err)
		return res, err
	}
	body, err := ReadResponseBody(rsp)
	res.Data = TxHash(body)
	return res, err
}

// GetTxStatus returns status of transaction.
func (c *Client) GetTxStatus(
	ctx context.Context,
	tx TxHash,
	opts *RequestOptions,
) (res *TxStatusResponse, err error) {
	res = &TxStatusResponse{}
	rsp, err := c.GetTxsStatuses(ctx, []TxHash{tx}, opts)
	res.Response = rsp.Response
	if len(rsp.Data) > 0 {
		res.Data = &rsp.Data[0]
	}
	return
}

// GetTxsStatuses returns status of transaction(s).
func (c *Client) GetTxsStatuses(
	ctx context.Context,
	txs []TxHash,
	opts *RequestOptions,
) (*TxsStatusesResponse, error) {
	res := &TxsStatusesResponse{}
	if len(txs) == 0 {
		err := ErrNoTxHash
		res.applyError(nil, err)
		return res, err
	}

	rsp, err := c.request(ctx, &res.Response, "POST", "/tx_status", txHashesPL(txs), opts)
	if err != nil {
		return res, err
	}
	return res, ReadAndUnmarshalResponse(rsp, &res.Response, &res.Data)
}

// NewTxWithAutoFee.
func NewTransaction() *TX {
	return &TX{}
}

func txHashesPL(txs []TxHash) io.Reader {
	var payload = struct {
		TxHashes []TxHash `json:"_tx_hashes"`
	}{txs}
	rpipe, w := io.Pipe()
	go func() {
		_ = json.NewEncoder(w).Encode(payload)
		defer w.Close()
	}()
	return rpipe
}

// AddUTxO adds utxo inputs and otputs to transaction
// Useful when composing batch transaction.
func (tx *TX) AddUTxO(utxo UTxO) error {
	for _, in := range utxo.Inputs {
		// check that UTxO inputs are not already used
		for _, tin := range tx.Inputs {
			if in.TxHash == tin.TxHash && in.TxIndex == tin.TxIndex {
				return fmt.Errorf("%w: %s#%d", ErrUTxOInputAlreadyUsed, in.TxHash, in.TxIndex)
			}
		}
		tx.Inputs = append(tx.Inputs, in)
	}

	tx.Outputs = append(tx.Outputs, utxo.Outputs...)

	return nil
}

// CalculateSharedFee calculates fee and shares fee between transaction outputs.
// Optionally you can provide destination addresses which will not share the fee.
// e.g
// most cases you would exclude your own address so that only
// external receivers pay the fee. Useful when using this lib in faucet.
func (tx *TX) CalculateSharedFee(params EpochParams, exclude ...Address) error {
	return nil
}

func (out *TxOutput) MarshalCBOR() ([]byte, error) {
	if len(out.PaymentAddr.Bech32) == 0 {
		return nil, fmt.Errorf("cbor: %w", ErrNoAddress)
	}

	// handle assets
	//   if len(o.Assets) > 0 {
	// 		return cbor.Marshal([]interface{}{o.Address.Bytes(), []interface{}{o.Amount, o.Assets}})
	// 	}
	return cbor.Marshal([]interface{}{out.PaymentAddr.Bech32.Bytes(), out.Value.IntPart()})
}
