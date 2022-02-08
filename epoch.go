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
	"net/url"
)

type (
	// EpochInfo defines model for epoch_info.
	EpochInfo struct {
		// Rewards accumulated as of given epoch (in lovelaces)
		ActiveStake string `json:"active_stake"`

		// Number of blocks created in epoch
		BlkCount int `json:"blk_count"`

		// Epoch number
		Epoch EpochNo `json:"epoch_no"`

		// Total fees incurred by transactions in epoch
		Fees Lovelace `json:"fees"`

		// Timestamp for first block created in epoch
		FirstBlockTime string `json:"first_block_time"`

		// Timestamp for last block created in epoch
		LastBlockTime string `json:"last_block_time"`

		// Total output value across all transactions in epoch
		OutSum Lovelace `json:"out_sum"`

		// Number of transactions submitted in epoch
		TxCount int `json:"tx_count"`
	}

	// EpochInfoResponse response of /epoch_info.
	EpochInfoResponse struct {
		Response
		Data []EpochInfo `json:"data"`
	}

	// EpochParams defines model for epoch_params.
	EpochParams struct {
		// The hash of the first block where these parameters are valid
		BlockHash string `json:"block_hash"`

		// The cost per UTxO word
		CoinsPerUtxoWord Lovelace `json:"coins_per_utxo_word"`

		// The percentage of the tx fee which must be provided as collateral
		// when including non-native scripts
		CollateralPercent int `json:"collateral_percent"`

		// The per language cost models
		CostModels string `json:"cost_models"`

		// The decentralisation parameter (1 fully centralised, 0 fully decentralised)
		Decentralisation float64 `json:"decentralisation"`

		// The hash of 32-byte string of extra random-ness added into
		// the protocol's entropy pool
		Entropy string `json:"entropy"`

		// Epoch number
		Epoch EpochNo `json:"epoch_no"`

		// The pledge influence on pool rewards
		Influence float64 `json:"influence"`

		// The amount (in lovelace) required for a deposit to register a stake address
		KeyDeposit Lovelace `json:"key_deposit"`

		// The maximum block header size (in bytes)
		MaxBhSize int `json:"max_bh_size"`

		// The maximum number of execution memory allowed to be used in a single block
		MaxBlockExMem float32 `json:"max_block_ex_mem"`

		// The maximum number of execution steps allowed to be used in a single block
		MaxBlockExSteps float32 `json:"max_block_ex_steps"`

		// The maximum block size (in bytes)
		MaxBlockSize int `json:"max_block_size"`

		// The maximum number of collateral inputs allowed in a transaction
		MaxCollateralInputs int `json:"max_collateral_inputs"`

		// The maximum number of epochs in the future that a pool retirement
		// is allowed to be scheduled for
		MaxEpoch int `json:"max_epoch"`

		// The maximum number of execution memory allowed to be used in a single transaction
		MaxTxExMem float32 `json:"max_tx_ex_mem"`

		// The maximum number of execution steps allowed to be used in a single transaction
		MaxTxExSteps float32 `json:"max_tx_ex_steps"`

		// The maximum transaction size (in bytes)
		MaxTxSize int `json:"max_tx_size"`

		// The maximum Val size
		MaxValSize float32 `json:"max_val_size"`

		// The 'a' parameter to calculate the minimum transaction fee
		MinFeeA int `json:"min_fee_a"`

		// The 'b' parameter to calculate the minimum transaction fee
		MinFeeB int `json:"min_fee_b"`

		// The minimum pool cost
		MinPoolCost Lovelace `json:"min_pool_cost"`

		// The minimum value of a UTxO entry
		MinUtxoValue int `json:"min_utxo_value"`

		// The monetary expansion rate
		MonetaryExpandRate float64 `json:"monetary_expand_rate"`

		// The nonce value for this epoch
		Nonce string `json:"nonce"`

		// The optimal number of stake pools
		OptimalPoolCount int `json:"optimal_pool_count"`

		// The amount (in lovelace) required for a deposit to register a stake pool
		PoolDeposit Lovelace `json:"pool_deposit"`

		// The per word cost of script memory usage
		PriceMem float64 `json:"price_mem"`

		// The cost of script execution step usage
		PriceStep float64 `json:"price_step"`

		// The protocol major version
		ProtocolMajor int `json:"protocol_major"`

		// The protocol minor version
		ProtocolMinor int `json:"protocol_minor"`

		// The treasury growth rate
		TreasuryGrowthRate float64 `json:"treasury_growth_rate"`
	}

	// EpochParamsResponse response of /epoch_params.
	EpochParamsResponse struct {
		Response
		Data []EpochParams `json:"data"`
	}
)

// GetEpochInfo returns the epoch information, all epochs if no epoch specified.
func (c *Client) GetEpochInfo(ctx context.Context, epoch *EpochNo) (res *EpochInfoResponse, err error) {
	res = &EpochInfoResponse{}
	params := url.Values{}
	if epoch != nil {
		params.Set("_epoch_no", fmt.Sprint(*epoch))
	}

	rsp, err := c.request(ctx, &res.Response, "GET", "/epoch_info", nil, params, nil)
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
	return
}

// GetEpochParams returns the protocol parameters for specific epoch,
// and information about all epochs if no epoch specified.
func (c *Client) GetEpochParams(ctx context.Context, epoch *EpochNo) (res *EpochParamsResponse, err error) {
	res = &EpochParamsResponse{}
	params := url.Values{}
	if epoch != nil {
		params.Set("_epoch_no", fmt.Sprint(*epoch))
	}

	rsp, err := c.request(ctx, &res.Response, "GET", "/epoch_params", nil, params, nil)
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
	return
}
