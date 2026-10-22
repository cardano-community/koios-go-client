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

type (
	// EUTxO holds the inputs and outputs of a transaction as returned by
	// `/tx_utxos`.
	EUTxO struct {
		// TxHash is hash of transaction.
		TxHash TxHash `json:"tx_hash"`

		// Inputs An array with details about inputs used in a transaction.
		Inputs []TxUTxO `json:"inputs"`

		// Outputs An array with details about outputs from the transaction.
		Outputs []TxUTxO `json:"outputs"`
	}

	// TxMetalabel is a known transaction metadata label.
	TxMetalabel struct {
		// A distinct known metalabel
		Key string `json:"key"`
	}

	// InlineDatum is a datum attached directly to a UTxO (CIP-32).
	InlineDatum struct {
		// Bytes is the datum bytes (hex).
		Bytes string `json:"bytes"`
		// Value is the datum in JSON format.
		Value json.RawMessage `json:"value"`
	}

	// ReferenceScript is a script attached to a UTxO that can be used to
	// satisfy script requirements without including the script in the
	// spending transaction.
	ReferenceScript struct {
		// Hash of the referenced script.
		Hash ScriptHash `json:"hash"`
		// Size in bytes.
		Size uint `json:"size"`
		// Type of the script.
		Type string `json:"type"`
		// Bytes of the script (hex).
		Bytes string `json:"bytes"`
		// Value is the script in JSON format (nullable).
		Value json.RawMessage `json:"value"`
	}

	// UTxO is a complete unspent or spent transaction output as returned by
	// `/utxo_info`, `/address_utxos`, `/credential_utxos`, `/account_utxos`,
	// `/asset_utxos` and `/script_utxos`.
	UTxO struct {
		// TxHash is the hash of the transaction that created this output.
		TxHash TxHash `json:"tx_hash"`
		// TxIndex is the index of the output in the transaction.
		TxIndex int `json:"tx_index"`
		// Address is the Cardano payment/base address (bech32 encoded).
		Address Address `json:"address"`
		// Value is the total sum of ADA on the UTxO (in lovelace).
		Value decimal.Decimal `json:"value"`
		// StakeAddress is the staking address of the output (nullable).
		StakeAddress Address `json:"stake_address"`
		// PaymentCred is the payment credential of the address (nullable).
		PaymentCred PaymentCredential `json:"payment_cred"`
		// EpochNo is the epoch of the block that created the output.
		EpochNo EpochNo `json:"epoch_no"`
		// BlockHeight is the height of the block that created the output.
		BlockHeight uint64 `json:"block_height"`
		// BlockTime is the time of the block that created the output.
		BlockTime Timestamp `json:"block_time"`
		// DatumHash is the hash of the Plutus data attached to the output (nullable).
		DatumHash DatumHash `json:"datum_hash"`
		// InlineDatum is the datum attached to the output (nullable).
		InlineDatum *InlineDatum `json:"inline_datum"`
		// ReferenceScript is the script attached to the output (nullable).
		ReferenceScript *ReferenceScript `json:"reference_script"`
		// AssetList is the list of assets on the output (nullable).
		AssetList []Asset `json:"asset_list"`
		// IsSpent is true if the UTxO has been spent.
		IsSpent bool `json:"is_spent"`
	}

	// TxUTxO is a transaction input or output as nested in `/tx_info`
	// and `/tx_utxos` results.
	TxUTxO struct {
		// PaymentAddr is the payment address and its credential.
		PaymentAddr PaymentAddr `json:"payment_addr"`
		// StakeAddress is the staking address of the output (nullable).
		StakeAddress Address `json:"stake_addr"`
		// TxHash is the hash of the transaction that created this output.
		TxHash TxHash `json:"tx_hash"`
		// TxIndex is the index of the output in the transaction.
		TxIndex int `json:"tx_index"`
		// Value is the total sum of ADA on the UTxO (in lovelace).
		Value decimal.Decimal `json:"value"`
		// DatumHash is the hash of the datum connected to the output (nullable).
		DatumHash DatumHash `json:"datum_hash"`
		// InlineDatum is the datum attached to the output (nullable).
		InlineDatum *InlineDatum `json:"inline_datum"`
		// ReferenceScript is the script attached to the output (nullable).
		ReferenceScript *ReferenceScript `json:"reference_script"`
		// AssetList is the list of assets on the output (nullable).
		AssetList []Asset `json:"asset_list"`
	}

	// TxCollateralOutput is the collateral return output of a transaction.
	//
	// The spec describes the field as an array of outputs, but a transaction
	// has at most one collateral return and the API sends a single object, with
	// asset_list as a Haskell-formatted string rather than an array; decoding
	// accepts either form of the wrapper and keeps asset_list undecoded.
	TxCollateralOutput struct {
		PaymentAddr     PaymentAddr      `json:"payment_addr"`
		StakeAddress    Address          `json:"stake_addr"`
		TxHash          TxHash           `json:"tx_hash"`
		TxIndex         int              `json:"tx_index"`
		Value           decimal.Decimal  `json:"value"`
		DatumHash       DatumHash        `json:"datum_hash"`
		InlineDatum     *InlineDatum     `json:"inline_datum"`
		ReferenceScript *ReferenceScript `json:"reference_script"`
		// AssetList is the raw asset_list value as sent by the API.
		AssetList json.RawMessage `json:"asset_list"`
	}

	// TxsWithdrawal withdrawal record in transaction.
	TxsWithdrawal struct {
		// Amount is withdrawal amount in lovelaces.
		Amount decimal.Decimal `json:"amount"`
		// StakeAddress fo withdrawal.
		StakeAddress Address `json:"stake_addr"`
	}

	// TxListItem is a transaction reference with its block position.
	TxListItem struct {
		TxHash      TxHash    `json:"tx_hash"`
		BlockHeight uint64    `json:"block_height"`
		BlockTime   Timestamp `json:"block_time"`
		EpochNo     EpochNo   `json:"epoch_no"`
	}

	// TxCbor is a raw transaction as returned by `/tx_cbor` and `/block_tx_cbor`.
	TxCbor struct {
		TxHash TxHash `json:"tx_hash"`
		// BlockHash is the hash of the block containing the transaction.
		BlockHash BlockHash `json:"block_hash"`
		// BlockHeight is the height of the block containing the transaction.
		BlockHeight uint64 `json:"block_height"`
		// EpochNo is the epoch of the block containing the transaction.
		EpochNo EpochNo `json:"epoch_no"`
		// AbsoluteSlot is the absolute slot of the block containing the transaction.
		AbsoluteSlot Slot `json:"absolute_slot"`
		// TxTimestamp is the time of the transaction.
		TxTimestamp Timestamp `json:"tx_timestamp"`
		// Cbor is the CBOR encoded raw transaction (hex).
		Cbor string `json:"cbor"`
		// ValidContract is false when the transaction contains a failed contract.
		ValidContract bool `json:"valid_contract"`
	}

	// TxByMetalabel is a transaction carrying a metadata label as returned
	// by `/tx_by_metalabel`.
	TxByMetalabel struct {
		TxHash       TxHash    `json:"tx_hash"`
		BlockHash    BlockHash `json:"block_hash"`
		BlockHeight  uint64    `json:"block_height"`
		EpochNo      EpochNo   `json:"epoch_no"`
		AbsoluteSlot Slot      `json:"absolute_slot"`
		TxTimestamp  Timestamp `json:"tx_timestamp"`
	}

	// EpochTxOut is a basic transaction output as returned by `/tx_outs_epoch`.
	EpochTxOut struct {
		TxHash       TxHash            `json:"tx_hash"`
		TxIndex      int               `json:"tx_index"`
		Address      Address           `json:"address"`
		StakeAddress Address           `json:"stake_address"`
		PaymentCred  PaymentCredential `json:"payment_cred"`
		EpochNo      EpochNo           `json:"epoch_no"`
		BlockHeight  uint64            `json:"block_height"`
		BlockTime    Timestamp         `json:"block_time"`
		Value        decimal.Decimal   `json:"value"`
		DatumHash    DatumHash         `json:"datum_hash"`
		IsSpent      bool              `json:"is_spent"`
	}

	// TxNativeScript is a native script present in a transaction.
	TxNativeScript struct {
		ScriptHash ScriptHash `json:"script_hash"`
		// ScriptJSON is the JSON representation of the timelock script.
		ScriptJSON json.RawMessage `json:"script_json"`
	}

	// TxVotingProcedure is a governance vote cast in a transaction.
	TxVotingProcedure struct {
		// ProposalTxHash is the hash of the transaction that created the proposal.
		ProposalTxHash TxHash `json:"proposal_tx_hash"`
		// ProposalIndex is the index of the governance proposal in its transaction.
		ProposalIndex int `json:"proposal_index"`
		// VoterRole is one of ConstitutionalCommittee, DRep or SPO.
		VoterRole string `json:"voter_role"`
		// Voter is the DRep ID (CIP-129), pool ID (bech32) or committee hot ID (CIP-129).
		Voter string `json:"voter"`
		// VoterHex is the voter ID in hex format.
		VoterHex string `json:"voter_hex"`
		// Vote is one of Yes, No or Abstain.
		Vote string `json:"vote"`
	}

	// TxProposalWithdrawal is an amount withdrawn from the treasury by a proposal.
	TxProposalWithdrawal struct {
		StakeAddress Address         `json:"stake_address"`
		Amount       decimal.Decimal `json:"amount"`
	}

	// TxProposalProcedure is a governance proposal submitted in a transaction.
	TxProposalProcedure struct {
		// Index of the governance proposal in the transaction.
		Index int `json:"index"`
		// Type is the governance action type, for example ParameterChange or InfoAction.
		Type string `json:"type"`
		// Description of the proposal action.
		Description json.RawMessage `json:"description"`
		// Deposit is the proposal deposit (in lovelace, nullable).
		Deposit decimal.Decimal `json:"deposit"`
		// ReturnAddress receives the deposit when it is repaid.
		ReturnAddress Address `json:"return_address"`
		// Expiration is the epoch at which the action is expected to expire (nullable).
		Expiration *EpochNo `json:"expiration"`
		// MetaURL is the URL of the proposal metadata JSON payload (nullable).
		MetaURL string `json:"meta_url"`
		// MetaHash is the hash of the contents of the metadata URL (nullable).
		MetaHash string `json:"meta_hash"`
		// Withdrawal lists the treasury withdrawals requested by the proposal.
		Withdrawal []TxProposalWithdrawal `json:"withdrawal"`
		// ParamProposal is the proposed parameter set, if any.
		ParamProposal json.RawMessage `json:"param_proposal"`
	}

	// TxInfo transaction info.
	TxInfo struct {
		// TxHash is hash of transaction.
		TxHash TxHash `json:"tx_hash"`
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
		Fee decimal.Decimal `json:"fee"`
		// TreasuryDonation is the donation to the on-chain treasury (in lovelaces).
		TreasuryDonation decimal.Decimal `json:"treasury_donation"`
		// Deposit is total deposits included in transaction (for example,
		// if it is registering a pool/key).
		Deposit decimal.Decimal `json:"deposit"`
		// InvalidBefore is slot before which transaction cannot be validated
		// (nil if not supplied).
		InvalidBefore *Slot `json:"invalid_before"`
		// InvalidAfter is slot after which transaction cannot be validated
		// (nil if not supplied).
		InvalidAfter *Slot `json:"invalid_after"`
		// CollateralInputs An array of collateral inputs needed when dealing with smart contracts.
		CollateralInputs []TxUTxO `json:"collateral_inputs"`
		// CollateralOutput is the collateral output returned as change if
		// the smart contract fails to execute (CIP-40).
		CollateralOutput *TxCollateralOutput `json:"collateral_output"`
		// ReferenceInputs An array of reference inputs (CIP-31).
		ReferenceInputs []TxUTxO `json:"reference_inputs"`
		// Inputs An array with details about inputs used in a transaction.
		Inputs []TxUTxO `json:"inputs"`
		// Outputs An array with details about outputs from the transaction.
		Outputs []TxUTxO `json:"outputs"`
		// Array of withdrawals with-in a transaction (if any)
		Withdrawals []TxsWithdrawal `json:"withdrawals"`
		// AssetsMinted An array of minted assets with-in a transaction (if any).
		AssetsMinted []Asset `json:"assets_minted"`
		// Metadata present with-in a transaction (if any)
		Metadata TxMetadata `json:"metadata"`
		// Certificates present with-in a transaction (if any)
		Certificates []Certificate `json:"certificates"`
		// NativeScripts present in a transaction (if any).
		NativeScripts []TxNativeScript `json:"native_scripts"`
		// PlutusContracts present in a transaction (if any).
		PlutusContracts []PlutusContract `json:"plutus_contracts"`
		// VotingProcedures are the governance votes in a transaction (if any).
		VotingProcedures []TxVotingProcedure `json:"voting_procedures"`
		// ProposalProcedures are the governance proposals in a transaction (if any).
		ProposalProcedures []TxProposalProcedure `json:"proposal_procedures"`
	}

	// TxInfoFlags selects the optional sections of a transaction info result.
	// The API omits every section that is not requested.
	TxInfoFlags struct {
		// Inputs includes transaction inputs.
		Inputs bool
		// Metadata includes transaction metadata.
		Metadata bool
		// Assets includes the assets involved in the transaction.
		Assets bool
		// Withdrawals includes stake account reward withdrawals.
		Withdrawals bool
		// Certs includes transaction certificates.
		Certs bool
		// Scripts includes collateral, reference, datum and script details.
		Scripts bool
		// Bytecode includes bytecode of reference and Plutus scripts.
		Bytecode bool
		// Governance includes governance certificates, votes and proposals.
		Governance bool
	}

	txInfoFlagsPayload struct {
		Inputs      bool `json:"_inputs,omitempty"`
		Metadata    bool `json:"_metadata,omitempty"`
		Assets      bool `json:"_assets,omitempty"`
		Withdrawals bool `json:"_withdrawals,omitempty"`
		Certs       bool `json:"_certs,omitempty"`
		Scripts     bool `json:"_scripts,omitempty"`
		Bytecode    bool `json:"_bytecode,omitempty"`
		Governance  bool `json:"_governance,omitempty"`
	}

	// TxsInfosResponse represents response from `/tx_info` endpoint.
	TxsInfosResponse struct {
		Response
		Data []TxInfo `json:"data"`
	}

	// TxInfoResponse represents response from `/tx_info` endpoint.
	// when requesting info about single transaction.
	TxInfoResponse struct {
		Response
		Data TxInfo `json:"data"`
	}

	// TxsUTxOsResponse represents response from `/tx_utxos` endpoint.
	TxsUTxOsResponse struct {
		Response
		Data []EUTxO `json:"data"`
	}

	// TxCborResponse represents response from `/tx_cbor` endpoint.
	TxCborResponse struct {
		Response
		Data []TxCbor `json:"data"`
	}

	// TxByMetalabelResponse represents response from `/tx_by_metalabel` endpoint.
	TxByMetalabelResponse struct {
		Response
		Data []TxByMetalabel `json:"data"`
	}

	// TxOutsEpochResponse represents response from `/tx_outs_epoch` endpoint.
	TxOutsEpochResponse struct {
		Response
		Data []EpochTxOut `json:"data"`
	}

	// TxMetadata transaction metadata lookup res for `/tx_metadata` endpoint.
	TxMetadata map[string]json.RawMessage

	// TxMetadataOf is the metadata of a single transaction.
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
		TxHash TxHash `json:"tx_hash"`
		// Confirmations is the number of block confirmations (nil when unknown).
		Confirmations *uint64 `json:"num_confirmations"`
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

	// UTxORef is a UTxO reference in the form "hash#index".
	UTxORef string

	// UTxOInfoResponse represents response from `/utxo_info` endpoint.
	UTxOInfoResponse struct {
		Response
		Data []UTxO `json:"data"`
	}
)

func (f *TxInfoFlags) payload() txInfoFlagsPayload {
	return txInfoFlagsPayload(*f)
}

// GetTxInfo returns detailed information about transaction(s). The optional
// sections of the result are selected with flags; a nil flags value requests
// none of them.
func (c *Client) GetTxInfo(
	ctx context.Context,
	txs []TxHash,
	flags *TxInfoFlags,
	opts *RequestOptions,
) (*TxsInfosResponse, error) {
	res := &TxsInfosResponse{}
	if len(txs) == 0 || len(txs[0]) == 0 {
		err := ErrNoTxHash
		res.applyError(nil, err)
		return res, err
	}

	payload := struct {
		TxHashes []TxHash `json:"_tx_hashes"`
		txInfoFlagsPayload
	}{TxHashes: txs}
	if flags != nil {
		payload.txInfoFlagsPayload = flags.payload()
	}
	body, err := jsonBody(payload)
	if err != nil {
		res.applyError(nil, err)
		return res, err
	}

	rsp, err := c.request(ctx, &res.Response, "POST", "/tx_info", body, opts)
	if err != nil {
		return res, err
	}
	return res, ReadAndUnmarshalResponse(rsp, &res.Response, &res.Data)
}

// GetTxMetadata returns metadata information (if any) for given transaction.
func (c *Client) GetTxMetadata(
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

// GetTxCbor returns the raw CBOR of the requested transactions.
func (c *Client) GetTxCbor(
	ctx context.Context,
	txs []TxHash,
	opts *RequestOptions,
) (*TxCborResponse, error) {
	res := &TxCborResponse{}
	if len(txs) == 0 || len(txs[0]) == 0 {
		err := ErrNoTxHash
		res.applyError(nil, err)
		return res, err
	}

	rsp, err := c.request(ctx, &res.Response, "POST", "/tx_cbor", txHashesPL(txs), opts)
	if err != nil {
		return res, err
	}
	return res, ReadAndUnmarshalResponse(rsp, &res.Response, &res.Data)
}

// GetTxByMetalabel returns the transactions that include the metadata label,
// optionally only those after the given block height (0 for no filter).
func (c *Client) GetTxByMetalabel(
	ctx context.Context,
	label uint64,
	afterBlockHeight uint64,
	opts *RequestOptions,
) (*TxByMetalabelResponse, error) {
	res := &TxByMetalabelResponse{}
	if opts == nil {
		opts = c.NewRequestOptions()
	}
	opts.QuerySet("_label", fmt.Sprint(label))
	if afterBlockHeight > 0 {
		opts.QuerySet("_after_block_height", fmt.Sprint(afterBlockHeight))
	}

	rsp, err := c.request(ctx, &res.Response, "GET", "/tx_by_metalabel", nil, opts)
	if err != nil {
		return res, err
	}
	return res, ReadAndUnmarshalResponse(rsp, &res.Response, &res.Data)
}

// GetTxOutsEpoch returns the basic details of all transaction outputs
// of the requested epoch.
func (c *Client) GetTxOutsEpoch(
	ctx context.Context,
	epoch EpochNo,
	opts *RequestOptions,
) (*TxOutsEpochResponse, error) {
	res := &TxOutsEpochResponse{}
	if opts == nil {
		opts = c.NewRequestOptions()
	}
	opts.QuerySet("_epoch_no", epoch.String())

	rsp, err := c.request(ctx, &res.Response, "GET", "/tx_outs_epoch", nil, opts)
	if err != nil {
		return res, err
	}
	return res, ReadAndUnmarshalResponse(rsp, &res.Response, &res.Data)
}

// GetTxUTxOs returns the inputs and outputs of the requested transactions.
//
// Deprecated: the API deprecates `/tx_utxos`; use GetUTxOInfo instead.
func (c *Client) GetTxUTxOs(
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

// GetTxStatus returns status of transaction(s).
func (c *Client) GetTxStatus(
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

// GetUTxOInfo returns the UTxO set for the requested UTxO references.
// The extended flag populates optional fields supported by the endpoint.
func (c *Client) GetUTxOInfo(
	ctx context.Context,
	refs []UTxORef,
	extended bool,
	opts *RequestOptions,
) (*UTxOInfoResponse, error) {
	res := &UTxOInfoResponse{}
	if len(refs) == 0 {
		err := ErrNoUTxORef
		res.applyError(nil, err)
		return res, err
	}

	if opts == nil {
		opts = c.NewRequestOptions()
	}

	rsp, err := c.request(ctx, &res.Response, "POST", "/utxo_info", utxoRefsPL(refs, extended), opts)
	if err != nil {
		res.applyError(nil, err)
		return res, err
	}

	return res, ReadAndUnmarshalResponse(rsp, &res.Response, &res.Data)
}

func txHashesPL(txs []TxHash) io.Reader {
	var payload = struct {
		TxHashes []TxHash `json:"_tx_hashes"`
	}{txs}
	rpipe := payloadReader(payload)
	return rpipe
}

func utxoRefsPL(refs []UTxORef, extended bool) io.Reader {
	var payload = struct {
		UtxORefs []UTxORef `json:"_utxo_refs"`
		Extended bool      `json:"_extended"`
	}{refs, extended}
	rpipe := payloadReader(payload)
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

// UnmarshalJSON accepts the collateral output as a single object or as a
// one-element array.
func (o *TxCollateralOutput) UnmarshalJSON(b []byte) error {
	type plain TxCollateralOutput
	if t := bytes.TrimSpace(b); len(t) > 0 && t[0] == '[' {
		var list []plain
		if err := json.Unmarshal(t, &list); err != nil {
			return err
		}
		if len(list) == 0 {
			*o = TxCollateralOutput{}
			return nil
		}
		*o = TxCollateralOutput(list[0])
		return nil
	}
	return json.Unmarshal(b, (*plain)(o))
}
