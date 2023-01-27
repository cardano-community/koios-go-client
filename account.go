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
	"fmt"
	"io"

	"github.com/shopspring/decimal"
)

type (
	// AccountInfo data returned by `/account_info`.
	AccountInfo struct {
		Status           string          `json:"status"`
		DelegatedPool    *PoolID         `json:"delegated_pool"`
		StakeAddress     Address         `json:"stake_address"`
		TotalBalance     decimal.Decimal `json:"total_balance"`
		UTxO             decimal.Decimal `json:"utxo"`
		Rewards          decimal.Decimal `json:"rewards"`
		Withdrawals      decimal.Decimal `json:"withdrawals"`
		RewardsAvailable decimal.Decimal `json:"rewards_available"`
		Reserves         decimal.Decimal `json:"reserves"`
		Treasury         decimal.Decimal `json:"treasury"`
	}

	// AccountRewardsInfo data returned by `/account_rewards`.
	AccountRewardsInfo struct {
		StakeAddress Address          `json:"stake_address"`
		Rewards      []AccountRewards `json:"rewards"`
	}

	AccountRewards struct {
		PoolID         PoolID          `json:"pool_id"`
		EarnedEpoch    EpochNo         `json:"earned_epoch"`
		SpendableEpoch EpochNo         `json:"spendable_epoch"`
		Amount         decimal.Decimal `json:"amount"`
		Type           string          `json:"type"`
	}

	// AccountHistoryEntry history entry list item.
	AccountHistory struct {
		StakeAddress Address               `json:"stake_address"`
		History      []AccountHistoryEntry `json:"history"`
	}
	AccountHistoryEntry struct {
		PoolID      PoolID          `json:"pool_id"`
		EpochNo     EpochNo         `json:"epoch_no"`
		ActiveStake decimal.Decimal `json:"active_stake"`
	}

	// AccountListResponse represents response from `/account_list` endpoint.
	AccountListResponse struct {
		Response
		Data []Address `json:"data"`
	}

	// AccountInfoResponse represents response from `/account_info` endpoint.
	AccountInfoResponse struct {
		Response
		Data *AccountInfo `json:"data"`
	}

	AccountsInfoResponse struct {
		Response
		Data []AccountInfo `json:"data"`
	}

	// AccountRewardsResponse represents response from `/account_rewards` endpoint.
	AccountRewardsResponse struct {
		Response
		Data *AccountRewardsInfo `json:"data"`
	}
	AccountsRewardsResponse struct {
		Response
		Data []AccountRewardsInfo `json:"data"`
	}

	// AccountAction data entry for `/account_updates`.
	AccountUpdate struct {
		ActionType   string    `json:"action_type"`
		TxHash       TxHash    `json:"tx_hash"`
		EpochNo      EpochNo   `json:"epoch_no"`
		EpochSlot    Slot      `json:"epoch_slot"`
		AbsoluteSlot Slot      `json:"absolute_slot"`
		BlockTime    Timestamp `json:"block_time"`
	}

	AccountUpdates struct {
		StakeAddress Address         `json:"stake_address"`
		Updates      []AccountUpdate `json:"updates"`
	}

	// AccountUpdatesResponse represents response from `/account_rewards` endpoint.
	AccountUpdatesResponse struct {
		Response
		Data *AccountUpdates `json:"data"`
	}

	AccountsUpdatesResponse struct {
		Response
		Data []AccountUpdates `json:"data"`
	}

	AccountAddresses struct {
		StakeAddress Address   `json:"stake_address"`
		Addresses    []Address `json:"addresses"`
	}
	// AccountAddressesResponse represents response from `/account_addresses` endpoint.
	AccountAddressesResponse struct {
		Response
		Data *AccountAddresses `json:"data"`
	}
	AccountsAddressesResponse struct {
		Response
		Data []AccountAddresses `json:"data"`
	}

	// AccountAssetsResponse represents response from `/account_assets` endpoint.
	AccountAssetsResponse struct {
		Response
		Data *AccountAssets `json:"data"`
	}

	AccountsAssetsResponse struct {
		Response
		Data []AccountAssets `json:"data"`
	}

	AccountAssets struct {
		StakeAddress Address `json:"stake_address"`
		Assets       []Asset `json:"asset_list"`
	}

	// AccountHistoryResponse represents response from `/account_history` endpoint.
	AccountHistoryResponse struct {
		Response
		Data *AccountHistory `json:"data"`
	}
	AccountsHistoryResponse struct {
		Response
		Data []AccountHistory `json:"data"`
	}
)

// GetAccountList returns a list of all accounts.
func (c *Client) GetAccounts(
	ctx context.Context,
	opts *RequestOptions,
) (res *AccountListResponse, err error) {
	res = &AccountListResponse{}
	rsp, err := c.request(ctx, &res.Response, "GET", "/account_list", nil, opts)
	if err != nil {
		return
	}

	accs := []struct {
		ID Address `json:"id"`
	}{}

	err = ReadAndUnmarshalResponse(rsp, &res.Response, &accs)

	if len(accs) > 0 {
		for _, a := range accs {
			res.Data = append(res.Data, a.ID)
		}
	}
	return
}

// GetAccountInfo returns the account info of any (payment or staking) address.
func (c *Client) GetAccountInfo(
	ctx context.Context,
	acc Address,
	cached bool,
	opts *RequestOptions,
) (res *AccountInfoResponse, err error) {
	res = &AccountInfoResponse{}
	res2, err := c.GetAccountsInfo(ctx, []Address{acc}, cached, opts)
	if err != nil {
		return
	}
	if len(res2.Data) == 1 {
		res.Data = &res2.Data[0]
	} else {
		return nil, fmt.Errorf("%w: no account info found for address %s", ErrNoData, acc)
	}
	return
}

func (c *Client) GetAccountsInfo(
	ctx context.Context,
	accs []Address,
	cached bool,
	opts *RequestOptions,
) (res *AccountsInfoResponse, err error) {
	res = &AccountsInfoResponse{}
	if len(accs) == 0 {
		err = ErrNoAddress
		res.applyError(nil, err)
		return
	}
	endpoint := "/account_info"
	if cached {
		endpoint = "/account_info_cached"
	}

	rsp, err := c.request(ctx, &res.Response, "POST", endpoint, stakeAddressesPL(accs, nil), opts)
	if err != nil {
		return
	}
	err = ReadAndUnmarshalResponse(rsp, &res.Response, &res.Data)
	return
}

// GetAccountRewards retruns the full rewards history (including MIR)
// for a stake address, or certain epoch if specified.
func (c *Client) GetAccountRewards(
	ctx context.Context,
	acc Address,
	epoch *EpochNo,
	opts *RequestOptions,
) (res *AccountRewardsResponse, err error) {
	res = &AccountRewardsResponse{}

	res2, err := c.GetAccountsRewards(ctx, []Address{acc}, epoch, opts)
	if err != nil {
		return
	}
	if len(res2.Data) == 1 {
		res.Data = &res2.Data[0]
	} else {
		return nil, fmt.Errorf("%w: no rewards found for account %s", ErrNoData, acc)
	}
	return
}

func (c *Client) GetAccountsRewards(
	ctx context.Context,
	accs []Address,
	epoch *EpochNo,
	opts *RequestOptions,
) (res *AccountsRewardsResponse, err error) {
	res = &AccountsRewardsResponse{}
	if len(accs) == 0 {
		err = ErrNoAddress
		res.applyError(nil, err)
		return
	}
	rsp, err := c.request(ctx, &res.Response, "POST", "/account_rewards", stakeAddressesPL(accs, epoch), opts)
	if err != nil {
		return
	}
	err = ReadAndUnmarshalResponse(rsp, &res.Response, &res.Data)
	return
}

// GetAccountUpdates (History) retruns the account updates
// (registration, deregistration, delegation and withdrawals).

func (c *Client) GetAccountUpdates(
	ctx context.Context,
	acc Address,
	epoch *EpochNo,
	opts *RequestOptions,
) (res *AccountUpdatesResponse, err error) {
	res = &AccountUpdatesResponse{}

	res2, err := c.GetAccountsUpdates(ctx, []Address{acc}, opts)
	if err != nil {
		return
	}
	if len(res2.Data) == 1 {
		res.Data = &res2.Data[0]
	} else {
		return nil, fmt.Errorf("%w: no updates found for account %s", ErrNoData, acc)
	}
	return
}

func (c *Client) GetAccountsUpdates(
	ctx context.Context,
	accs []Address,
	opts *RequestOptions,
) (res *AccountsUpdatesResponse, err error) {
	res = &AccountsUpdatesResponse{}
	if len(accs) == 0 {
		err = ErrNoAddress
		res.applyError(nil, err)
		return
	}
	rsp, err := c.request(ctx, &res.Response, "POST", "/account_updates", stakeAddressesPL(accs, nil), opts)
	if err != nil {
		return
	}
	err = ReadAndUnmarshalResponse(rsp, &res.Response, &res.Data)
	return
}

// GetAccountAddresses retruns all addresses associated with an account.
func (c *Client) GetAccountAddresses(
	ctx context.Context,
	acc Address,
	epoch *EpochNo,
	opts *RequestOptions,
) (res *AccountAddressesResponse, err error) {
	res = &AccountAddressesResponse{}

	res2, err := c.GetAccountsAddresses(ctx, []Address{acc}, opts)
	if err != nil {
		return
	}
	if len(res2.Data) == 1 {
		res.Data = &res2.Data[0]
	} else {
		return nil, fmt.Errorf("%w: no updates found for account %s", ErrNoData, acc)
	}
	return
}

func (c *Client) GetAccountsAddresses(
	ctx context.Context,
	accs []Address,
	opts *RequestOptions,
) (res *AccountsAddressesResponse, err error) {
	res = &AccountsAddressesResponse{}
	if len(accs) == 0 {
		err = ErrNoAddress
		res.applyError(nil, err)
		return
	}
	rsp, err := c.request(ctx, &res.Response, "POST", "/account_addresses", stakeAddressesPL(accs, nil), opts)
	if err != nil {
		return
	}
	err = ReadAndUnmarshalResponse(rsp, &res.Response, &res.Data)
	return
}

// GetAccountAssets retruns all the native asset balance of an account.
func (c *Client) GetAccountAssets(
	ctx context.Context,
	acc Address,
	opts *RequestOptions,
) (res *AccountAssetsResponse, err error) {
	res = &AccountAssetsResponse{}

	res2, err := c.GetAccountsAssets(ctx, []Address{acc}, opts)
	if err != nil {
		return
	}
	if len(res2.Data) == 1 {
		res.Data = &res2.Data[0]
	} else {
		return nil, fmt.Errorf("%w: no assets found for account %s", ErrNoData, acc)
	}
	return
}

func (c *Client) GetAccountsAssets(
	ctx context.Context,
	accs []Address,
	opts *RequestOptions,
) (res *AccountsAssetsResponse, err error) {
	res = &AccountsAssetsResponse{}
	if len(accs) == 0 {
		err = ErrNoAddress
		res.applyError(nil, err)
		return
	}
	rsp, err := c.request(ctx, &res.Response, "POST", "/account_assets", stakeAddressesPL(accs, nil), opts)
	if err != nil {
		return
	}
	err = ReadAndUnmarshalResponse(rsp, &res.Response, &res.Data)
	return
}

// GetAccountHistory retruns the staking history of an account.
func (c *Client) GetAccountHistory(
	ctx context.Context,
	acc Address,
	epoch *EpochNo,
	opts *RequestOptions,
) (res *AccountHistoryResponse, err error) {
	res = &AccountHistoryResponse{}

	res2, err := c.GetAccountsHistory(ctx, []Address{acc}, opts)
	if err != nil {
		return
	}
	if len(res2.Data) == 1 {
		res.Data = &res2.Data[0]
	} else {
		return nil, fmt.Errorf("%w: no history found for account %s", ErrNoData, acc)
	}
	return
}

func (c *Client) GetAccountsHistory(
	ctx context.Context,
	accs []Address,
	opts *RequestOptions,
) (res *AccountsHistoryResponse, err error) {
	res = &AccountsHistoryResponse{}
	if len(accs) == 0 {
		err = ErrNoAddress
		res.applyError(nil, err)
		return
	}
	rsp, err := c.request(ctx, &res.Response, "POST", "/account_history", stakeAddressesPL(accs, nil), opts)
	if err != nil {
		return
	}
	err = ReadAndUnmarshalResponse(rsp, &res.Response, &res.Data)
	return
}

func stakeAddressesPL(addrs []Address, epoch *EpochNo) io.Reader {
	var payload = struct {
		Adresses []Address `json:"_stake_addresses"`
		Epoch    *EpochNo  `json:"_epoch_no,omitempty"`
	}{addrs, epoch}
	rpipe, w := io.Pipe()
	go func() {
		_ = json.NewEncoder(w).Encode(payload)
		defer w.Close()
	}()
	return rpipe
}
