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

	"github.com/shopspring/decimal"
)

type (
	// AccountInfo data returned by `/account_info`.
	AccountInfo struct {
		Status           string          `json:"status"`
		DelegatedPool    PoolID          `json:"delegated_pool"`
		TotalBalance     decimal.Decimal `json:"total_balance"`
		UTxO             decimal.Decimal `json:"utxo"`
		Rewards          decimal.Decimal `json:"rewards"`
		Withdrawals      decimal.Decimal `json:"withdrawals"`
		RewardsAvailable decimal.Decimal `json:"rewards_available"`
		Reserves         decimal.Decimal `json:"reserves"`
		Treasury         decimal.Decimal `json:"treasury"`
	}

	// AccountRewards data returned by `/account_rewards`.
	AccountRewards struct {
		PoolID         PoolID          `json:"pool_id"`
		EarnedEpoch    EpochNo         `json:"earned_epoch"`
		SpendableEpoch EpochNo         `json:"spendable_epoch"`
		Amount         decimal.Decimal `json:"amount"`
		Type           string          `json:"type"`
	}

	// AccountHistoryEntry history entry list item.
	AccountHistoryEntry struct {
		StakeAddress StakeAddress    `json:"stake_address"`
		PoolID       PoolID          `json:"pool_id"`
		Epoch        EpochNo         `json:"epoch_no"`
		ActiveStake  decimal.Decimal `json:"active_stake"`
	}

	// AccountAsset asset list item when requesting account info.
	AccountAsset struct {
		// Name Asset Name (hex).
		Name string `json:"asset_name"`

		// PolicyID Asset Policy ID (hex).
		PolicyID PolicyID `json:"asset_policy"`

		// Quantity of assets
		Quantity decimal.Decimal `json:"quantity"`
	}

	// AccountAction data entry for `/account_updates`.
	AccountAction struct {
		ActionType string `json:"action_type"`
		TxHash     TxHash `json:"tx_hash"`
	}

	// AccountListResponse represents response from `/account_list` endpoint.
	AccountListResponse struct {
		Response
		Data []StakeAddress `json:"response"`
	}

	// AccountInfoResponse represents response from `/account_info` endpoint.
	AccountInfoResponse struct {
		Response
		Data *AccountInfo `json:"response"`
	}

	// AccountRewardsResponse represents response from `/account_rewards` endpoint.
	AccountRewardsResponse struct {
		Response
		Data []AccountRewards `json:"response"`
	}

	// AccountUpdatesResponse represents response from `/account_rewards` endpoint.
	AccountUpdatesResponse struct {
		Response
		Data []AccountAction `json:"response"`
	}

	// AccountAddressesResponse represents response from `/account_addresses` endpoint.
	AccountAddressesResponse struct {
		Response
		Data []Address `json:"response"`
	}

	// AccountAssetsResponse represents response from `/account_assets` endpoint.
	AccountAssetsResponse struct {
		Response
		Data []AccountAsset `json:"response"`
	}

	// AccountHistoryResponse represents response from `/account_history` endpoint.
	AccountHistoryResponse struct {
		Response
		Data []AccountHistoryEntry `json:"response"`
	}
)

// GetAccountList returns a list of all accounts.
func (c *Client) GetAccountList(
	ctx context.Context,
	opts *RequestOptions,
) (res *AccountListResponse, err error) {
	res = &AccountListResponse{}
	rsp, err := c.request(ctx, &res.Response, "GET", "/account_list", nil, opts)
	if err != nil {
		return
	}

	accs := []struct {
		ID StakeAddress `json:"id"`
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
	addr StakeAddress,
	opts *RequestOptions,
) (res *AccountInfoResponse, err error) {
	res = &AccountInfoResponse{}

	if _, err = addr.Valid(); err != nil {
		res.applyError(nil, err)
		return
	}

	if opts == nil {
		opts = c.NewRequestOptions()
	}
	opts.QuerySet("_address", addr.String())

	rsp, err := c.request(ctx, &res.Response, "GET", "/account_info", nil, opts)
	if err != nil {
		return
	}

	addrs := []AccountInfo{}
	err = ReadAndUnmarshalResponse(rsp, &res.Response, &addrs)

	if len(addrs) == 1 {
		res.Data = &addrs[0]
	}
	return
}

// GetAccountRewards retruns the full rewards history (including MIR)
// for a stake address, or certain epoch if specified.
func (c *Client) GetAccountRewards(
	ctx context.Context,
	addr StakeAddress,
	epoch *EpochNo,
	opts *RequestOptions,
) (res *AccountRewardsResponse, err error) {
	res = &AccountRewardsResponse{}

	if _, err = addr.Valid(); err != nil {
		res.applyError(nil, err)
		return
	}

	if opts == nil {
		opts = c.NewRequestOptions()
	}
	opts.QuerySet("_stake_address", addr.String())

	if epoch != nil {
		opts.QuerySet("_epoch_no", epoch.String())
	}

	rsp, err := c.request(ctx, &res.Response, "GET", "/account_rewards", nil, opts)
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
	addr StakeAddress,
	opts *RequestOptions,
) (res *AccountUpdatesResponse, err error) {
	res = &AccountUpdatesResponse{}

	if _, err = addr.Valid(); err != nil {
		res.applyError(nil, err)
		return
	}

	if opts == nil {
		opts = c.NewRequestOptions()
	}
	opts.QuerySet("_stake_address", addr.String())

	rsp, err := c.request(ctx, &res.Response, "GET", "/account_updates", nil, opts)
	if err != nil {
		return
	}
	err = ReadAndUnmarshalResponse(rsp, &res.Response, &res.Data)
	return
}

// GetAccountAddresses retruns all addresses associated with an account.
func (c *Client) GetAccountAddresses(
	ctx context.Context,
	addr StakeAddress,
	opts *RequestOptions,
) (res *AccountAddressesResponse, err error) {
	res = &AccountAddressesResponse{}

	if _, err = addr.Valid(); err != nil {
		res.applyError(nil, err)
		return res, err
	}

	if opts == nil {
		opts = c.NewRequestOptions()
	}
	opts.QuerySet("_address", addr.String())

	rsp, err := c.request(ctx, &res.Response, "GET", "/account_addresses", nil, opts)
	if err != nil {
		return res, err
	}
	addrs := []struct {
		Addr Address `json:"address"`
	}{}

	err = ReadAndUnmarshalResponse(rsp, &res.Response, &addrs)

	if len(addrs) > 0 {
		for _, a := range addrs {
			res.Data = append(res.Data, a.Addr)
		}
	}
	return res, err
}

// GetAccountAssets retruns all the native asset balance of an account.
func (c *Client) GetAccountAssets(
	ctx context.Context,
	addr StakeAddress,
	opts *RequestOptions,
) (res *AccountAssetsResponse, err error) {
	res = &AccountAssetsResponse{}

	if _, err = addr.Valid(); err != nil {
		res.applyError(nil, err)
		return
	}

	if opts == nil {
		opts = c.NewRequestOptions()
	}
	opts.QuerySet("_address", addr.String())

	rsp, err := c.request(ctx, &res.Response, "GET", "/account_assets", nil, opts)
	if err != nil {
		return
	}
	err = ReadAndUnmarshalResponse(rsp, &res.Response, &res.Data)
	return
}

// GetAccountHistory retruns the staking history of an account.
func (c *Client) GetAccountHistory(
	ctx context.Context,
	addr StakeAddress,
	opts *RequestOptions,
) (res *AccountHistoryResponse, err error) {
	res = &AccountHistoryResponse{}

	if _, err = addr.Valid(); err != nil {
		res.applyError(nil, err)
		return
	}

	if opts == nil {
		opts = c.NewRequestOptions()
	}
	opts.QuerySet("_address", addr.String())

	rsp, err := c.request(ctx, &res.Response, "GET", "/account_history", nil, opts)
	if err != nil {
		return
	}
	err = ReadAndUnmarshalResponse(rsp, &res.Response, &res.Data)
	return
}
