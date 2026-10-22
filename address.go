// SPDX-License-Identifier: Apache-2.0
//
// Copyright © 2022 The Cardano Community Authors

package koios

import (
	"context"
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

	// AddressListItem is an address as returned by `/address_list`.
	AddressListItem struct {
		// Address is the Cardano payment/base address (bech32 encoded).
		Address Address `json:"address"`
		// AddressHex is the address in hex format.
		AddressHex string `json:"address_hex"`
		// HasScript is true if the payment credential is a script hash.
		HasScript bool `json:"has_script"`
		// PaymentCred is the payment credential (nullable).
		PaymentCred PaymentCredential `json:"payment_cred"`
		// StakeAddress is the staking address (nullable).
		StakeAddress Address `json:"stake_address"`
		// StakeAddressHex is the staking address in hex format (nullable).
		StakeAddressHex string `json:"stake_address_hex"`
	}

	// AddressOutput is a basic transaction output as returned by `/address_outputs`.
	AddressOutput struct {
		Address      Address           `json:"address"`
		TxHash       TxHash            `json:"tx_hash"`
		TxIndex      int               `json:"tx_index"`
		Value        decimal.Decimal   `json:"value"`
		StakeAddress Address           `json:"stake_address"`
		PaymentCred  PaymentCredential `json:"payment_cred"`
		EpochNo      EpochNo           `json:"epoch_no"`
		BlockHeight  uint64            `json:"block_height"`
		BlockTime    Timestamp         `json:"block_time"`
		IsSpent      bool              `json:"is_spent"`
	}

	// AddressListResponse represents response from `/address_list` endpoint.
	AddressListResponse struct {
		Response
		Data []AddressListItem `json:"data"`
	}

	// AddressOutputsResponse represents response from `/address_outputs` endpoint.
	AddressOutputsResponse struct {
		Response
		Data []AddressOutput `json:"data"`
	}

	// AddressTx is a transaction reference with its block position.
	AddressTx struct {
		TxHash      TxHash    `json:"tx_hash"`
		EpochNo     EpochNo   `json:"epoch_no"`
		BlockTime   Timestamp `json:"block_time"`
		BlockHeight uint64    `json:"block_height"`
	}

	AddressAssets struct {
		Asset
		Address Address `json:"address"`
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

	UTxOsResponse struct {
		Response
		Data []UTxO `json:"data"`
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
	res.Response = res2.Response
	if err != nil {
		return
	}
	if len(res2.Data) == 1 {
		res.Data = &res2.Data[0]
	}
	return
}

// GetAddressesInfo returns address info - balance, associated stake address
// (if any) and UTxO set for the given addresses.
func (c *Client) GetAddressesInfo(
	ctx context.Context,
	addr []Address,
	opts *RequestOptions,
) (res *AddressesInfoResponse, err error) {
	res = &AddressesInfoResponse{}
	if len(addr) == 0 {
		res.applyError(nil, ErrNoAddressesProvided)
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

	rpipe := payloadReader(payload)

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

	rpipe := payloadReader(payload)

	rsp, err := c.request(ctx, &res.Response, "POST", "/credential_txs", rpipe, opts)
	res.applyError(nil, err)

	return res, ReadAndUnmarshalResponse(rsp, &res.Response, &res.Data)
}

func (c *Client) GetAddressUTxOs(
	ctx context.Context,
	addrs []Address,
	extended bool,
	opts *RequestOptions,
) (*UTxOsResponse, error) {
	res := &UTxOsResponse{}
	if len(addrs) == 0 {
		err := ErrNoAddress
		res.applyError(nil, err)
		return res, err
	}

	var payload = struct {
		Adresses []Address `json:"_addresses"`
		Extended bool      `json:"_extended,omitempty"`
	}{
		Adresses: addrs,
		Extended: extended,
	}

	rpipe := payloadReader(payload)

	rsp, err := c.request(ctx, &res.Response, "POST", "/address_utxos", rpipe, opts)
	res.applyError(nil, err)
	err = ReadAndUnmarshalResponse(rsp, &res.Response, &res.Data)
	return res, err
}

func (c *Client) GetCredentialUTxOs(
	ctx context.Context,
	creds []PaymentCredential,
	extended bool,
	opts *RequestOptions,
) (*UTxOsResponse, error) {
	res := &UTxOsResponse{}
	if len(creds) == 0 || len(creds[0]) == 0 {
		err := ErrNoCredentialsProvided
		res.applyError(nil, err)
		return res, err
	}
	var payload = struct {
		Creds    []PaymentCredential `json:"_payment_credentials"`
		Extended bool                `json:"_extended,omitempty"`
	}{
		Creds:    creds,
		Extended: extended,
	}
	rpipe := payloadReader(payload)

	rsp, err := c.request(ctx, &res.Response, "POST", "/credential_utxos", rpipe, opts)
	res.applyError(nil, err)
	err = ReadAndUnmarshalResponse(rsp, &res.Response, &res.Data)
	return res, err
}

// GetAddressList returns the list of all used addresses on the blockchain (paginated).
func (c *Client) GetAddressList(
	ctx context.Context,
	opts *RequestOptions,
) (*AddressListResponse, error) {
	res := &AddressListResponse{}
	rsp, err := c.request(ctx, &res.Response, "GET", "/address_list", nil, opts)
	if err != nil {
		return res, err
	}
	return res, ReadAndUnmarshalResponse(rsp, &res.Response, &res.Data)
}

// GetAddressOutputs returns basic transaction output info for the given
// addresses, optionally only after the given block height (0 for no filter).
func (c *Client) GetAddressOutputs(
	ctx context.Context,
	addrs []Address,
	afterBlockHeight uint64,
	opts *RequestOptions,
) (*AddressOutputsResponse, error) {
	res := &AddressOutputsResponse{}
	if len(addrs) == 0 {
		err := ErrNoAddress
		res.applyError(nil, err)
		return res, err
	}

	payload := struct {
		Addresses        []Address `json:"_addresses"`
		AfterBlockHeight uint64    `json:"_after_block_height,omitempty"`
	}{addrs, afterBlockHeight}
	body, err := jsonBody(payload)
	if err != nil {
		res.applyError(nil, err)
		return res, err
	}

	rsp, err := c.request(ctx, &res.Response, "POST", "/address_outputs", body, opts)
	if err != nil {
		return res, err
	}
	return res, ReadAndUnmarshalResponse(rsp, &res.Response, &res.Data)
}

func addressesPL(addrs []Address) io.Reader {
	var payload = struct {
		Adresses []Address `json:"_addresses"`
	}{addrs}
	rpipe := payloadReader(payload)
	return rpipe
}
