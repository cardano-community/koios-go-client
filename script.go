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
	// ScriptHash defines type for _script_hash.
	ScriptHash string

	DatumHash string

	// ScriptRedeemers defines model for script_redeemers.
	ScriptRedeemers struct {
		// Hash of Transaction for which details are being shown
		ScriptHash ScriptHash `json:"script_hash"`

		// Redeemers list
		Redeemers []ScriptRedeemer `json:"redeemers"`
	}

	// ScriptRedeemer model.
	ScriptRedeemer struct {
		// The Hash of the Plutus Data
		DatumHash string `json:"datum_hash"`

		// The actual data in json format
		DatumValue map[string]any `json:"datum_value"`

		// The budget in fees to run a script - the fees depend on the
		// ExUnits and the current prices.
		Fee decimal.Decimal `json:"fee,omitempty"`

		// What kind pf validation this redeemer is used for,
		// it can be one of 'spend', 'mint', 'cert', 'reward'.
		Purpose string `json:"purpose"`

		// TxHash of Transaction containing the redeemer.
		TxHash TxHash `json:"tx_hash"`

		// TxIndex The index of the redeemer pointer in the transaction.
		TxIndex uint32 `json:"tx_index"`

		// The budget in Memory to run a script.
		UnitMem int `json:"unit_mem"`

		// The budget in Cpu steps to run a script.
		UnitSteps int `json:"unit_steps"`
	}

	// DatumInfo datum information for given datum hash
	DatumInfo struct {
		Hash  DatumHash        `json:"hash"`
		Bytes string           `json:"bytes"`
		Value *json.RawMessage `json:"value"`
	}

	PlutusContract struct {
		Address       Address    `json:"address"`
		ScriptHash    ScriptHash `json:"script_hash"`
		ByteCode      string     `json:"bytecode"`
		Size          int        `json:"size"`
		ValidContract bool       `json:"valid_contract"`
	}

	// NativeScript item of native script list.
	NativeScript struct {
		// Hash of the script creation transaction
		CreationTxHash TxHash `json:"creation_tx_hash"`

		// Hash of a script
		ScriptHash string `json:"script_hash"`
		Type       string `json:"type"`
		Script     struct {
			Type    string            `json:"type"`
			Scripts []json.RawMessage `json:"scripts"`
		} `json:"script"`
	}

	// NativeScriptListResponse represents response from `/native_script_list` endpoint.
	NativeScriptListResponse struct {
		Response
		Data []NativeScript `json:"data"`
	}

	// PlutusScriptListItem item of plutus script list.
	PlutusScriptListItem struct {
		// Hash of the script creation transaction
		CreationTxHash TxHash `json:"creation_tx_hash"`

		// Hash of a script
		ScriptHash string `json:"script_hash"`
	}

	// PlutusScriptListResponse represents response from `/plutus_script_list` endpoint.
	PlutusScriptListResponse struct {
		Response
		Data []PlutusScriptListItem `json:"data"`
	}

	// ScriptRedeemersResponse represents response from `/script_redeemers` endpoint.
	ScriptRedeemersResponse struct {
		Response
		Data *ScriptRedeemers `json:"data"`
	}

	// DatumInfosResponse represents response from `/datum_info` endpoint.
	DatumInfosResponse struct {
		Response
		Data []DatumInfo `json:"data"`
	}
	DatumInfoResponse struct {
		Response
		Data *DatumInfo `json:"data"`
	}
)

// String returns ScriptHash as string.
func (v ScriptHash) String() string {
	return string(v)
}

// String returns DatumHash as string.
func (v DatumHash) String() string {
	return string(v)
}

// GetNativeScriptList returns list of all existing native script hashes
// along with their creation transaction hashes.
func (c *Client) GetNativeScripts(
	ctx context.Context,
	opts *RequestOptions,
) (res *NativeScriptListResponse, err error) {
	res = &NativeScriptListResponse{}
	rsp, err := c.request(ctx, &res.Response, "GET", "/native_script_list", nil, opts)
	if err != nil {
		return
	}
	err = ReadAndUnmarshalResponse(rsp, &res.Response, &res.Data)
	return
}

// GetPlutusScriptList returns all existing Plutus script
// hashes along with their creation transaction hashes.
func (c *Client) GetPlutusScripts(
	ctx context.Context,
	opts *RequestOptions,
) (res *PlutusScriptListResponse, err error) {
	res = &PlutusScriptListResponse{}
	rsp, err := c.request(ctx, &res.Response, "GET", "/plutus_script_list", nil, opts)
	if err != nil {
		return
	}
	err = ReadAndUnmarshalResponse(rsp, &res.Response, &res.Data)
	return
}

// GetScriptRedeemers returns a list of all redeemers for a given script hash.
func (c *Client) GetScriptRedeemers(
	ctx context.Context,
	sh ScriptHash,
	opts *RequestOptions,
) (res *ScriptRedeemersResponse, err error) {
	res = &ScriptRedeemersResponse{}

	if opts == nil {
		opts = c.NewRequestOptions()
	}
	opts.QuerySet("_script_hash", sh.String())

	rsp, err := c.request(ctx, &res.Response, "GET", "/script_redeemers", nil, opts)
	if err != nil {
		return
	}
	r := []ScriptRedeemers{}
	err = ReadAndUnmarshalResponse(rsp, &res.Response, &r)

	if len(r) == 1 {
		res.Data = &r[0]
	}
	return
}

// GetTxStatus returns status of transaction.
func (c *Client) GetDatumInfo(
	ctx context.Context,
	hash DatumHash,
	opts *RequestOptions,
) (res *DatumInfoResponse, err error) {
	res = &DatumInfoResponse{}
	rsp, err := c.GetDatumInfos(ctx, []DatumHash{hash}, opts)
	res.Response = rsp.Response
	if len(rsp.Data) == 1 {
		res.Data = &rsp.Data[0]
	} else {
		err = fmt.Errorf("%w: %s", ErrNoData, hash)
	}
	return
}

// GetTxsStatuses returns status of transaction(s).
func (c *Client) GetDatumInfos(
	ctx context.Context,
	hashes []DatumHash,
	opts *RequestOptions,
) (*DatumInfosResponse, error) {
	res := &DatumInfosResponse{}
	if len(hashes) == 0 {
		err := ErrNoDatumHash
		res.applyError(nil, err)
		return res, err
	}

	rsp, err := c.request(ctx, &res.Response, "POST", "/datum_info", datumHashesPL(hashes), opts)
	if err != nil {
		return res, err
	}
	return res, ReadAndUnmarshalResponse(rsp, &res.Response, &res.Data)
}

func datumHashesPL(hashes []DatumHash) io.Reader {
	var payload = struct {
		DatumHashes []DatumHash `json:"_datum_hashes"`
	}{hashes}
	rpipe, w := io.Pipe()
	go func() {
		_ = json.NewEncoder(w).Encode(payload)
		defer w.Close()
	}()
	return rpipe
}
