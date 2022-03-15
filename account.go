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
	"fmt"
	"net/url"
)

type (
	// AccountInfo data returned by `/account_info`.
	AccountInfo struct {
		Status           string   `json:"status"`
		DelegatedPool    PoolID   `json:"delegated_pool"`
		TotalBalance     Lovelace `json:"total_balance"`
		UTxO             Lovelace `json:"utxo"`
		Rewards          Lovelace `json:"rewards"`
		Withdrawals      Lovelace `json:"withdrawals"`
		RewardsAvailable Lovelace `json:"rewards_available"`
		Reserves         Lovelace `json:"reserves"`
		Treasury         Lovelace `json:"treasury"`
	}

	// AccountRewards data returned by `/account_rewards`.
	AccountRewards struct {
		PoolID         PoolID   `json:"pool_id"`
		EarnedEpoch    EpochNo  `json:"earned_epoch"`
		SpendableEpoch EpochNo  `json:"spendable_epoch"`
		Amount         Lovelace `json:"amount"`
		Type           string   `json:"type"`
	}

	// AccountHistoryEntry history entry list item.
	AccountHistoryEntry struct {
		StakeAddress StakeAddress `json:"stake_address"`
		PoolID       PoolID       `json:"pool_id"`
		Epoch        EpochNo      `json:"epoch_no"`
		ActiveStake  Lovelace     `json:"active_stake"`
	}

	// AccountAsset asset list item when requesting account info.
	AccountAsset struct {
		// Name Asset Name (hex).
		Name string `json:"asset_name"`

		// PolicyID Asset Policy ID (hex).
		PolicyID PolicyID `json:"asset_policy"`

		// Quantity of assets
		Quantity Lovelace `json:"quantity"`
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
func (c *Client) GetAccountList(ctx context.Context) (res *AccountListResponse, err error) {
	res = &AccountListResponse{}
	rsp, err := c.request(ctx, &res.Response, "GET", "/account_list", nil, nil, nil)

	if err != nil {
		return
	}

	accs := []struct {
		ID StakeAddress `json:"id"`
	}{}

	err = readAndUnmarshalResponse(rsp, &res.Response, &accs)

	if len(accs) > 0 {
		for _, a := range accs {
			res.Data = append(res.Data, a.ID)
		}
	}
	return
}

// GetAccountInfo returns the account info of any (payment or staking) address.

func (c *Client) GetAccountInfo(ctx context.Context, addr Address) (res *AccountInfoResponse, err error) {
	res = &AccountInfoResponse{}
	if len(addr) == 0 {
		err = ErrNoAddress
		res.applyError(nil, err)
		return
	}
	params := url.Values{}
	params.Set("_address", string(addr))

	rsp, err := c.request(ctx, &res.Response, "GET", "/account_info", nil, params, nil)
	if err != nil {
		return
	}

	addrs := []AccountInfo{}
	err = readAndUnmarshalResponse(rsp, &res.Response, &addrs)

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
) (res *AccountRewardsResponse, err error) {
	res = &AccountRewardsResponse{}
	params := url.Values{}
	params.Set("_stake_address", string(addr))
	if epoch != nil {
		params.Set("_epoch_no", fmt.Sprint(*epoch))
	}

	rsp, err := c.request(ctx, &res.Response, "GET", "/account_rewards", nil, params, nil)
	if err != nil {
		return
	}
	err = readAndUnmarshalResponse(rsp, &res.Response, &res.Data)
	return
}

// GetAccountUpdates (History) retruns the account updates
// (registration, deregistration, delegation and withdrawals).
func (c *Client) GetAccountUpdates(
	ctx context.Context,
	addr StakeAddress,
) (res *AccountUpdatesResponse, err error) {
	res = &AccountUpdatesResponse{}
	params := url.Values{}
	params.Set("_stake_address", string(addr))

	rsp, err := c.request(ctx, &res.Response, "GET", "/account_updates", nil, params, nil)
	if err != nil {
		return
	}
	err = readAndUnmarshalResponse(rsp, &res.Response, &res.Data)
	return
}

// GetAccountAddresses retruns all addresses associated with an account.
func (c *Client) GetAccountAddresses(
	ctx context.Context,
	addr StakeAddress,
) (res *AccountAddressesResponse, err error) {
	res = &AccountAddressesResponse{}
	params := url.Values{}
	params.Set("_address", string(addr))

	rsp, err := c.request(ctx, &res.Response, "GET", "/account_addresses", nil, params, nil)
	if err != nil {
		return
	}
	addrs := []struct {
		Addr Address `json:"address"`
	}{}

	err = readAndUnmarshalResponse(rsp, &res.Response, &addrs)

	if len(addrs) > 0 {
		for _, a := range addrs {
			res.Data = append(res.Data, a.Addr)
		}
	}
	return
}

// GetAccountAssets retruns all the native asset balance of an account.
func (c *Client) GetAccountAssets(
	ctx context.Context,
	addr StakeAddress,
) (res *AccountAssetsResponse, err error) {
	res = &AccountAssetsResponse{}
	params := url.Values{}
	params.Set("_address", string(addr))

	rsp, err := c.request(ctx, &res.Response, "GET", "/account_assets", nil, params, nil)
	if err != nil {
		return
	}
	err = readAndUnmarshalResponse(rsp, &res.Response, &res.Data)
	return
}

// GetAccountHistory retruns the staking history of an account.
func (c *Client) GetAccountHistory(
	ctx context.Context,
	addr StakeAddress,
) (res *AccountHistoryResponse, err error) {
	res = &AccountHistoryResponse{}
	params := url.Values{}
	params.Set("_address", string(addr))

	rsp, err := c.request(ctx, &res.Response, "GET", "/account_history", nil, params, nil)
	if err != nil {
		return
	}
	err = readAndUnmarshalResponse(rsp, &res.Response, &res.Data)
	return
}
