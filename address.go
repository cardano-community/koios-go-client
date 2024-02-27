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
	// Address defines type for _address.
	Address string

	// AddressInfo esponse for `/address_info`.
	AddressInfo struct {
		// Balance ADA Lovelace balance of address
		Balance decimal.Decimal `json:"balance"`

		Address Address `json:"address"`

		// StakeAddress associated with address
		StakeAddress Address `json:"stake_address,omitempty"`

		ScriptAddress bool `json:"script_address"`

		UTxOs []UTxO `json:"utxo_set"`
	}

	AddressTx struct {
		TxHash      TxHash    `json:"tx_hash"`
		EpochNo     EpochNo   `json:"epoch_no"`
		BlockTime   Timestamp `json:"block_time"`
		BlockHeight uint64    `json:"block_height"`
	}

	AddressAssets struct {
		Address   Address `json:"address"`
		AssetList []Asset `json:"asset_list"`
	}

	// AddressInfoResponse represents response from `/address_info` endpoint.
	AddressInfoResponse struct {
		Response
		Data *AddressInfo `json:"data"`
	}
	AddressesInfoResponse struct {
		Response
		Data []AddressInfo `json:"data"`
	}

	// AddressTxsResponse represents response from `/address_txs` endpoint.
	AddressTxsResponse struct {
		Response
		Data []AddressTx `json:"data"`
	}

	// CredentialTxsResponse represents response from `/credential_txs` endpoint.
	CredentialTxsResponse struct {
		Response
		Data []AddressTx `json:"data"`
	}

	// AddressAssetsResponse represents response from `/address_info` endpoint.
	AddressAssetsResponse struct {
		Response
		Data *AddressAssets `json:"data"`
	}

	AddressesAssetsResponse struct {
		Response
		Data []AddressAssets `json:"data"`
	}

	AddressAsset struct {
		AssetName      AssetName       `json:"asset_name"`
		AssetNameASCII string          `json:"asset_name_ascii"`
		Balance        decimal.Decimal `json:"balance"`
	}
)

// Valid validates address and returns false and error
// if address is invalid otherwise it returns true, nil.
func (a Address) Valid() (bool, error) {
	if len(a) == 0 {
		return false, ErrNoAddress
	}
	return true, nil
}

// String returns StakeAddress as string.
func (a Address) String() string {
	return string(a)
}

// Bytes returns address bytes.
func (a Address) Bytes() []byte {
	return []byte(a)
}

// GetAddressInfo returns address info - balance,
// associated stake address (if any) and UTxO set.
func (c *Client) GetAddressInfo(
	ctx context.Context,
	addr Address,
	opts *RequestOptions,
) (res *AddressInfoResponse, err error) {
	res = &AddressInfoResponse{}
	res2, err := c.GetAddressesInfo(ctx, []Address{addr}, opts)
	if err != nil {
		return
	}
	if len(res2.Data) == 1 {
		res.Data = &res2.Data[0]
	}
	return
}

func (c *Client) GetAddressesInfo(
	ctx context.Context,
	addr []Address,
	opts *RequestOptions,
) (res *AddressesInfoResponse, err error) {
	res = &AddressesInfoResponse{}
	if len(addr) == 0 {
		err = ErrNoAddress
		res.applyError(nil, err)
		return
	}
	rsp, err := c.request(ctx, &res.Response, "POST", "/address_info", addressesPL(addr), opts)
	if err != nil {
		return
	}
	err = ReadAndUnmarshalResponse(rsp, &res.Response, &res.Data)
	return
}

// GetAddressTxs returns the transaction hash list of input address array,
// optionally filtering after specified block height (inclusive).

func (c *Client) GetAddressTxs(
	ctx context.Context,
	addrs []Address,
	h uint64,
	opts *RequestOptions,
) (*AddressTxsResponse, error) {
	res := &AddressTxsResponse{}
	if len(addrs) == 0 {
		err := ErrNoAddress
		res.applyError(nil, err)
		return res, err
	}

	var payload = struct {
		Adresses         []Address `json:"_addresses"`
		AfterBlockHeight uint64    `json:"_after_block_height,omitempty"`
	}{
		Adresses:         addrs,
		AfterBlockHeight: h,
	}

	rpipe, w := io.Pipe()
	go func() {
		_ = json.NewEncoder(w).Encode(payload)
		defer w.Close()
	}()

	rsp, err := c.request(ctx, &res.Response, "POST", "/address_txs", rpipe, opts)
	if err != nil {
		return res, err
	}

	err = ReadAndUnmarshalResponse(rsp, &res.Response, &res.Data)
	return res, err
}

// GetAddressAssets returns the list of all the assets (policy, name and quantity)
// for a given address.
func (c *Client) GetAddressAssets(
	ctx context.Context,
	addr Address,
	opts *RequestOptions,
) (res *AddressAssetsResponse, err error) {
	res = &AddressAssetsResponse{}
	if len(addr) == 0 {
		err = ErrNoAddress
		res.applyError(nil, err)
		return
	}

	rsp, err := c.GetAddressesAssets(ctx, []Address{addr}, opts)
	if err != nil {
		return
	}
	if len(rsp.Data) == 1 {
		res.Data = &rsp.Data[0]
	} else {
		return nil, fmt.Errorf("%w: no assets on address %s", ErrNoData, addr)
	}
	return
}

func (c *Client) GetAddressesAssets(
	ctx context.Context,
	addrs []Address,
	opts *RequestOptions,
) (*AddressesAssetsResponse, error) {
	res := &AddressesAssetsResponse{}
	if len(addrs) == 0 {
		err := ErrNoAddress
		res.applyError(nil, err)
		return res, err
	}

	rsp, err := c.request(ctx, &res.Response, "POST", "/address_assets", addressesPL(addrs), opts)
	if err != nil {
		return res, err
	}
	err = ReadAndUnmarshalResponse(rsp, &res.Response, &res.Data)
	if err != nil {
		return res, err
	}
	if len(res.Data) == 0 {
		return nil, fmt.Errorf("%w: no assets on %d addresses", ErrNoData, len(addrs))
	}
	return res, err
}

// GetCredentialTxs returns the transaction hash list of input
// payment credential array, optionally filtering after specified block height (inclusive).

func (c *Client) GetCredentialTxs(
	ctx context.Context,
	creds []PaymentCredential,
	h uint64,
	opts *RequestOptions,
) (res *CredentialTxsResponse, err error) {
	res = &CredentialTxsResponse{}
	if len(creds) == 0 || len(creds[0]) == 0 {
		err = ErrNoAddress
		res.applyError(nil, err)
		return
	}

	var payload = struct {
		Credentials      []PaymentCredential `json:"_payment_credentials"`
		AfterBlockHeight uint64              `json:"_after_block_height,omitempty"`
	}{
		Credentials:      creds,
		AfterBlockHeight: h,
	}

	rpipe, w := io.Pipe()
	go func() {
		_ = json.NewEncoder(w).Encode(payload)
		defer w.Close()
	}()

	rsp, err := c.request(ctx, &res.Response, "POST", "/credential_txs", rpipe, opts)
	res.applyError(nil, err)

	return res, ReadAndUnmarshalResponse(rsp, &res.Response, &res.Data)
}

func addressesPL(addrs []Address) io.Reader {
	var payload = struct {
		Adresses []Address `json:"_addresses"`
	}{addrs}
	rpipe, w := io.Pipe()
	go func() {
		_ = json.NewEncoder(w).Encode(payload)
		defer w.Close()
	}()
	return rpipe
}
