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
	"fmt"
	"net/http"
	"net/url"
)

type (
	// AccountInfo data returned by `/account_info`.
	AccountInfo struct {
		Status           string     `json:"status"`
		DelegatedPool    PoolBech32 `json:"delegated_pool"`
		TotalBalance     Lovelace   `json:"total_balance"`
		UTxO             Lovelace   `json:"utxo"`
		Rewards          Lovelace   `json:"rewards"`
		Withdrawals      Lovelace   `json:"withdrawals"`
		RewardsAvailable Lovelace   `json:"rewards_available"`
		Reserves         Lovelace   `json:"reserves"`
		Treasury         Lovelace   `json:"treasury"`
	}

	// AccountRewards data returned by `/account_rewards`.
	AccountRewards struct {
		PoolID         PoolBech32 `json:"pool_id"`
		EarnedEpoch    EpochNo    `json:"earned_epoch"`
		SpendableEpoch EpochNo    `json:"spendable_epoch"`
		Amount         Lovelace   `json:"amount"`
		Type           string     `json:"type"`
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
)

// GetTip returns the tip info about the latest block seen by chain.
func (c *Client) GetAccountList(ctx context.Context) (res *AccountListResponse, err error) {
	res = &AccountListResponse{}
	rsp, err := c.request(ctx, &res.Response, "GET", nil, "/account_list", nil, nil)
	if err != nil {
		res.applyError(nil, err)
		return
	}

	body, err := readResponseBody(rsp)
	if err != nil {
		res.applyError(body, err)
		return
	}

	accs := []struct {
		ID StakeAddress `json:"id"`
	}{}

	if err = json.Unmarshal(body, &accs); err != nil {
		res.applyError(body, err)
		return
	}

	if rsp.StatusCode != http.StatusOK {
		res.applyError(body, err)
		return
	}
	if len(accs) > 0 {
		for _, a := range accs {
			res.Data = append(res.Data, a.ID)
		}
	}
	return res, nil
}

// GetAccountInfo returns the account info of any (payment or staking) address.
//nolint: dupl
func (c *Client) GetAccountInfo(ctx context.Context, addr Address) (res *AccountInfoResponse, err error) {
	res = &AccountInfoResponse{}
	if len(addr) == 0 {
		err = ErrNoAddress
		res.applyError(nil, err)
		return
	}
	params := url.Values{}
	params.Set("_address", string(addr))

	rsp, err := c.request(ctx, &res.Response, "GET", nil, "/account_info", params, nil)
	if err != nil {
		return
	}
	body, err := readResponseBody(rsp)
	if err != nil {
		res.applyError(nil, err)
		return
	}

	addrs := []AccountInfo{}

	if err = json.Unmarshal(body, &addrs); err != nil {
		res.applyError(body, err)
		return
	}

	if rsp.StatusCode != http.StatusOK {
		res.applyError(body, err)
		return
	}
	if len(addrs) == 1 {
		res.Data = &addrs[0]
	}
	res.ready()
	return res, nil
}

// GetAccountRewards retruns the full rewards history (including MIR)
// for a stake address, or certain epoch if specified.
func (c *Client) GetAccountRewards(
	ctx context.Context,
	addr StakeAddress,
	epochNo *EpochNo,
) (res *AccountRewardsResponse, err error) {
	res = &AccountRewardsResponse{}
	params := url.Values{}
	params.Set("_stake_address", string(addr))
	if epochNo != nil {
		params.Set("_epoch_no", fmt.Sprint(*epochNo))
	}
	rsp, err := c.request(ctx, &res.Response, "GET", nil, "/account_rewards", params, nil)
	if err != nil {
		res.applyError(nil, err)
		return
	}
	body, err := readResponseBody(rsp)
	if err != nil {
		res.applyError(body, err)
		return
	}

	if err = json.Unmarshal(body, &res.Data); err != nil {
		res.applyError(body, err)
		return
	}
	res.ready()
	return res, nil
}

// AccountUpdatesResponse (History) retruns the account updates
// (registration, deregistration, delegation and withdrawals).
func (c *Client) GetAccountUpdates(
	ctx context.Context,
	addr StakeAddress,
) (res *AccountUpdatesResponse, err error) {
	res = &AccountUpdatesResponse{}
	params := url.Values{}
	params.Set("_stake_address", string(addr))

	rsp, err := c.request(ctx, &res.Response, "GET", nil, "/account_updates", params, nil)
	if err != nil {
		res.applyError(nil, err)
		return
	}
	body, err := readResponseBody(rsp)
	if err != nil {
		res.applyError(body, err)
		return
	}

	if err = json.Unmarshal(body, &res.Data); err != nil {
		res.applyError(body, err)
		return
	}
	res.ready()
	return res, nil
}
