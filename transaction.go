// SPDX-License-Identifier: Apache-2.0
//
// Copyright © 2022 The Cardano Community Authors

package koios

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/shopspring/decimal"
)

// introduces breaking change since v1.3.0

type (
	TX struct {
		EUTxO
		TxInfo
	}

	// UTxO model holds inputs and outputs for given UTxO.
	EUTxO struct {
		/// TxHash is hash of transaction.
		TxHash TxHash `json:"tx_hash,omitempty"`

		// Inputs An array with details about inputs used in a transaction.
		Inputs []UTxO `json:"inputs" cbor:"0,keyasint"`
		// Outputs An array with details about outputs from the transaction.
		Outputs []UTxO `json:"outputs" cbor:"1,keyasint"`
	}

	// TxMetalabel
	TxMetalabel struct {
		// A distinct known metalabel
		Key string `json:"key"`
	}

	UTxO struct {
		// Hash of this transaction.
		TxHash TxHash `json:"tx_hash"`

		// Index of output UTxO.
		TxIndex int `json:"tx_index"`

		// where funds were sent or change to be returned.
		PaymentAddr *PaymentAddr `json:"payment_addr,omitempty"`

		// StakeAddress for transaction's output UTxO.
		StakeAddress *Address `json:"stake_addr,omitempty"`

		// BlockHeight int `json:"block_height"`
		// BlockTime   Timestamp `json:"block_time"`

		// Total sum on the output address.
		Value decimal.Decimal `json:"value"`

		DatumHash       DatumHash `json:"datum_hash,omitempty"`
		InlineDatum     any       `json:"inline_datum,omitempty"`
		ReferenceScript any       `json:"reference_script,omitempty"`
		// An array of assets to be included in output UTxO.
		AssetList []Asset `json:"asset_list,omitempty"`
	}

	// TxsWithdrawal withdrawal record in transaction.
	TxsWithdrawal struct {
		// Amount is withdrawal amount in lovelaces.
		Amount decimal.Decimal `json:"amount,omitempty"`
		// StakeAddress fo withdrawal.
		StakeAddress Address `json:"stake_addr,omitempty"`
	}

	// TxInfo transaction info.
	TxInfo struct {
		// BlockHash is hash of the block in which transaction was included.
		BlockHash BlockHash `json:"block_hash"`

		// BlockHeight is block number on chain where transaction was included.
		BlockHeight int `json:"block_height"`

		// Epoch number.
		EpochNo EpochNo `json:"epoch_no"`

		// EpochSlot is slot number within epoch.
		EpochSlot Slot `json:"epoch_slot"`

		// AbsoluteSlot is overall slot number (slots from genesis block of chain).
		AbsoluteSlot Slot `json:"absolute_slot"`

		// TxTimestamp is timestamp when block containing transaction was created.
		TxTimestamp Timestamp `json:"tx_timestamp"`

		// TxBlockIndex is index of transaction within block.
		TxBlockIndex int `json:"tx_block_index"`

		// TxSize is transaction size in bytes.
		TxSize int `json:"tx_size"`

		// TotalOutput is total sum of all transaction outputs (in lovelaces).
		TotalOutput decimal.Decimal `json:"total_output"`

		// Fee is total transaction fee (in lovelaces).
		Fee decimal.Decimal `json:"fee" cbor:"2,keyasint"`

		// Deposit is total deposits included in transaction (for example,
		// if it is registering a pool/key).
		Deposit decimal.Decimal `json:"deposit"`

		// InvalidAfter is slot number after which transaction cannot be validated.
		InvalidAfter Timestamp `json:"invalid_after,omitempty" cbor:"3,keyasint,omitempty"`

		// InvalidBefore is slot number before which transaction cannot be validated.
		// (if supplied, else 0)
		InvalidBefore Timestamp `json:"invalid_before,omitempty" cbor:"8,keyasint,omitempty"`

		// CollateralInputs An array of collateral inputs needed when dealing with smart contracts.
		CollateralInputs []UTxO `json:"collateral_inputs,omitempty"`

		// CollateralOutput
		CollateralOutput *UTxO `json:"collateral_output,omitempty"`

		// CollateralInputs An array of collateral inputs needed when dealing with smart contracts.
		ReferenceInputs []UTxO `json:"reference_inputs,omitempty"`

		// AssetsMinted An array of minted assets with-in a transaction (if any).
		AssetsMinted []Asset `json:"assets_minted,omitempty"`

		// Metadata present with-in a transaction (if any)
		Metadata TxMetadata `json:"metadata,omitempty"`

		// Array of withdrawals with-in a transaction (if any)
		Withdrawals []TxsWithdrawal `json:"withdrawals,omitempty"`

		// Certificates present with-in a transaction (if any)
		Certificates []Certificate `json:"certificates,omitempty"`

		NativeScripts   []NativeScript   `json:"native_scripts,omitempty"`
		PlutusContracts []PlutusContract `json:"plutus_contracts,omitempty"`
	}

	// TxsInfosResponse represents response from `/tx_info` endpoint.
	TxsInfosResponse struct {
		Response
		Data []TX `json:"data"`
	}

	// TxInfoResponse represents response from `/tx_info` endpoint.
	// when requesting info about single transaction.
	TxInfoResponse struct {
		Response
		Data TX `json:"data"`
	}

	// TxUTxOsResponse represents response from `/tx_utxos` endpoint.
	TxUTxOsResponse struct {
		Response
		Data *EUTxO `json:"data"`
	}

	// TxsUTxOsResponse represents response from `/tx_utxos` endpoint.
	TxsUTxOsResponse struct {
		Response
		Data []EUTxO `json:"data"`
	}

	// TxMetadata transaction metadata lookup res for `/tx_metadata` endpoint.
	TxMetadata map[string]json.RawMessage

	TxMetadataOf struct {
		TxHash   TxHash     `json:"tx_hash"`
		Metadata TxMetadata `json:"metadata,omitempty"`
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
		Data *TxMetadataOf `json:"data"`
	}

	// TxsMetadataResponse represents response from `/tx_metadata` endpoint.
	TxsMetadataResponse struct {
		Response
		Data []TxMetadataOf `json:"data"`
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
		Data []TxStatus `json:"data"`
	}

	// TxStatusResponse represents response from `/tx_status` endpoint.
	TxStatusResponse struct {
		Response
		Data *TxStatus `json:"data"`
	}
)

// GetTxInfo returns detailed information about transaction.
func (c *Client) GetTxInfo(
	ctx context.Context,
	hash TxHash,
	opts *RequestOptions,
) (res *TxInfoResponse, err error) {
	res = &TxInfoResponse{}
	rsp, err := c.GetTxsInfo(ctx, []TxHash{hash}, opts)
	res.Response = rsp.Response
	if len(rsp.Data) == 1 {
		res.Data = rsp.Data[0]
	} else {
		err = fmt.Errorf("%w: %s", ErrNoData, hash)
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
	hash TxHash,
	opts *RequestOptions,
) (res *TxUTxOsResponse, err error) {
	res = &TxUTxOsResponse{}
	rsp, err := c.GetTxsUTxOs(ctx, []TxHash{hash}, opts)
	res.Response = rsp.Response
	if len(rsp.Data) == 1 {
		res.Data = &rsp.Data[0]
	} else {
		err = fmt.Errorf("%w: %s", ErrNoData, hash)
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
	hash TxHash,
	opts *RequestOptions,
) (res *TxMetadataResponse, err error) {
	res = &TxMetadataResponse{}
	rsp, err := c.GetTxsMetadata(ctx, []TxHash{hash}, opts)
	res.Response = rsp.Response
	if len(rsp.Data) == 1 {
		res.Data = &rsp.Data[0]
	} else {
		err = fmt.Errorf("%w: %s", ErrNoData, hash)
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
	opts.HeaderSet("Content-Type", "application/cbor")
	opts.HeaderSet("Content-Length", fmt.Sprint(len(cborb)))

	rsp, err := c.request(ctx, &res.Response, method, "/submittx", bytes.NewBuffer(cborb), opts)
	if err != nil {
		body, _ := ReadResponseBody(rsp)
		res.applyError(body, err)
		return res, err
	}
	body, err := ReadResponseBody(rsp)
	res.Data = TxHash(strings.Trim(string(body), "\""))
	return res, err
}

// GetTxStatus returns status of transaction.
func (c *Client) GetTxStatus(
	ctx context.Context,
	hash TxHash,
	opts *RequestOptions,
) (res *TxStatusResponse, err error) {
	res = &TxStatusResponse{}
	rsp, err := c.GetTxsStatuses(ctx, []TxHash{hash}, opts)
	res.Response = rsp.Response
	if len(rsp.Data) == 1 {
		res.Data = &rsp.Data[0]
	} else {
		err = fmt.Errorf("%w: %s", ErrNoData, hash)
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

type metaArrayItem struct {
	// Key is metadata (index).
	Key string `json:"key,omitempty"`
	// JSON containing details about metadata within transaction.
	JSON json.RawMessage `json:"json,omitempty"`
}

func (m *TxMetadata) UnmarshalJSON(b []byte) error {
	if len(b) == 0 || string(b) == "[]" {
		return nil
	}
	var txMetadata map[string]json.RawMessage
	var txMetadataArray []metaArrayItem
	if err2 := json.Unmarshal(b, &txMetadataArray); err2 != nil {
		if err3 := json.Unmarshal(b, &txMetadata); err3 == nil {
			*m = txMetadata
			return nil
		}
		return fmt.Errorf("unmarshal metadata: %w", err2)
	}
	if len(txMetadataArray) == 0 {
		return nil
	}
	if len(txMetadataArray) == 1 && len(txMetadataArray[0].Key) == 0 {
		return nil
	}
	(*m) = make(TxMetadata)
	for _, meta := range txMetadataArray {
		(*m)[meta.Key] = meta.JSON
	}
	return nil
}

type metaListItem struct {
	TxHash   TxHash                     `json:"tx_hash"`
	Metadata map[string]json.RawMessage `json:"metadata"`
}

func (m *TxMetadataOf) UnmarshalJSON(b []byte) error {
	if len(b) == 0 {
		return nil
	}
	var data metaListItem
	if err := json.Unmarshal(b, &data); err != nil {
		return fmt.Errorf("unmarshal metadata list: %w", err)
	}
	m.TxHash = data.TxHash
	m.Metadata = data.Metadata
	return nil
}
