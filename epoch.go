// SPDX-License-Identifier: Apache-2.0
//
// Copyright © 2022 The Cardano Community Authors

package koios

import (
	"context"
	"fmt"

	"github.com/shopspring/decimal"
)

// introduces breaking change since v1.3.0

type (
	// EpochNo defines type for _epoch_no.
	EpochNo uint

	// EpochInfo defines model for epoch_info.
	EpochInfo struct {
		// Epoch number
		EpochNo EpochNo `json:"epoch_no"`

		// OutSum total output value across all transactions in epoch.
		OutSum decimal.Decimal `json:"out_sum"`

		// Number of blocks created in epoch
		BlkCount int `json:"blk_count"`

		// Total fees incurred by transactions in epoch
		Fees decimal.Decimal `json:"fees"`

		// Timestamp for first block created in epoch
		FirstBlockTime Timestamp `json:"first_block_time"`

		// Timestamp for last block created in epoch
		LastBlockTime Timestamp `json:"last_block_time"`

		// Number of transactions submitted in epoch
		TxCount int `json:"tx_count"`

		// EndTime of epoch
		EndTime Timestamp `json:"end_time"`

		// StartTime of epoch
		StartTime Timestamp `json:"start_time"`

		// ActiveStake Total active stake in epoch stake snapshot
		// (null for pre-Shelley epochs)
		ActiveStake decimal.Decimal `json:"active_stake,omitempty"`

		// TotalRewards earned in epoch (null for pre-Shelley epochs)
		TotalRewards decimal.Decimal `json:"total_rewards,omitempty"`

		// AvgBlkReward Average block reward for epoch (null for pre-Shelley epochs)
		AvgBlkReward decimal.Decimal `json:"avg_blk_reward,omitempty"`
	}

	// EpochInfoResponse response of /epoch_info.
	EpochInfoResponse struct {
		Response
		Data []EpochInfo `json:"data"`
	}

	// EpochParams defines model for epoch_params.
	EpochParams struct {
		// Epoch number
		EpochNo EpochNo `json:"epoch_no"`
		// The 'a' parameter to calculate the minimum transaction fee
		MinFeeA decimal.Decimal `json:"min_fee_a"`
		// The 'b' parameter to calculate the minimum transaction fee
		MinFeeB decimal.Decimal `json:"min_fee_b"`
		// The maximum block size (in bytes)
		MaxBlockSize int `json:"max_block_size"`
		// The maximum transaction size (in bytes)
		MaxTxSize int `json:"max_tx_size"`
		// The maximum block header size (in bytes)
		MaxBhSize int `json:"max_bh_size"`
		// The amount (in lovelace) required for a deposit to register a stake address
		KeyDeposit decimal.Decimal `json:"key_deposit"`
		// The amount (in lovelace) required for a deposit to register a stake pool
		PoolDeposit decimal.Decimal `json:"pool_deposit"`
		// The maximum number of epochs in the future that a pool retirement
		// is allowed to be scheduled for
		MaxEpoch int `json:"max_epoch"`
		// The optimal number of stake pools
		OptimalPoolCount int `json:"optimal_pool_count"`
		// The pledge influence on pool rewards
		Influence float64 `json:"influence"`
		// The monetary expansion rate
		MonetaryExpandRate float64 `json:"monetary_expand_rate"`
		// The treasury growth rate
		TreasuryGrowthRate float64 `json:"treasury_growth_rate"`
		// The decentralisation parameter (1 fully centralised, 0 fully decentralised)
		Decentralisation float64 `json:"decentralisation"`
		// The hash of 32-byte string of extra random-ness added into
		// the protocol's entropy pool
		ExtraEntropy string `json:"extra_entropy"`
		// The protocol major version
		ProtocolMajor int `json:"protocol_major"`
		// The protocol minor version
		ProtocolMinor int `json:"protocol_minor"`
		// The minimum value of a UTxO entry
		MinUtxoValue decimal.Decimal `json:"min_utxo_value"`
		// The minimum pool cost
		MinPoolCost decimal.Decimal `json:"min_pool_cost"`
		// The nonce value for this epoch
		Nonce string `json:"nonce"`
		// The hash of the first block where these parameters are valid
		BlockHash BlockHash `json:"block_hash"`
		// The per language cost models
		CostModels map[string][]int `json:"cost_models"`
		// The per word cost of script memory usage
		PriceMem decimal.Decimal `json:"price_mem"`
		// The cost of script execution step usage
		PriceStep decimal.Decimal `json:"price_step"`
		// The maximum number of execution memory allowed to be used in a single transaction
		MaxTxExMem float32 `json:"max_tx_ex_mem"`
		// The maximum number of execution steps allowed to be used in a single transaction
		MaxTxExSteps float32 `json:"max_tx_ex_steps"`
		// The maximum number of execution memory allowed to be used in a single block
		MaxBlockExMem float32 `json:"max_block_ex_mem"`
		// The maximum number of execution steps allowed to be used in a single block
		MaxBlockExSteps float32 `json:"max_block_ex_steps"`
		// The maximum Val size
		MaxValSize float64 `json:"max_val_size"`
		// The percentage of the tx fee which must be provided as collateral
		// when including non-native scripts
		CollateralPercent int `json:"collateral_percent"`
		// The maximum number of collateral inputs allowed in a transaction
		MaxCollateralInputs int `json:"max_collateral_inputs"`
		// The cost per UTxO word
		CoinsPerUtxoSize decimal.Decimal `json:"coins_per_utxo_size"`
	}

	// EpochParamsResponse response of /epoch_params.
	EpochParamsResponse struct {
		Response
		Data []EpochParams `json:"data"`
	}

	BlockProtocol struct {
		// The protocol major version
		ProtoMajor int `json:"proto_major"`

		// The protocol minor version
		ProtoMinor int `json:"proto_minor"`

		Blocks int `json:"blocks"`
	}

	EpochBlockProtocolsResponse struct {
		Response
		Data []BlockProtocol `json:"data"`
	}
)

// GetEpochInfo returns the epoch information, all epochs if no epoch specified.
func (c *Client) GetEpochInfo(
	ctx context.Context,
	epoch EpochNo,
	includeNextEpoch bool,
	opts *RequestOptions,
) (res *EpochInfoResponse, err error) {
	res = &EpochInfoResponse{}
	if opts == nil {
		opts = c.NewRequestOptions()
	}
	if !opts.query.Has("order") {
		opts.query.Set("order", "epoch_no.desc")
	}
	if epoch > 0 {
		opts.QuerySet("_epoch_no", epoch.String())
	}
	if includeNextEpoch {
		opts.QuerySet("_include_next_epoch", "true")
	}

	rsp, err := c.request(ctx, &res.Response, "GET", "/epoch_info", nil, opts)
	if err != nil {
		return
	}
	err = ReadAndUnmarshalResponse(rsp, &res.Response, &res.Data)
	return
}

// GetEpochParams returns the protocol parameters for specific epoch,
// and information about all epochs if no epoch specified.
func (c *Client) GetEpochParams(
	ctx context.Context,
	epoch EpochNo,
	opts *RequestOptions,
) (res *EpochParamsResponse, err error) {
	res = &EpochParamsResponse{}
	if opts == nil {
		opts = c.NewRequestOptions()
	}
	if epoch > 0 {
		opts.QuerySet("_epoch_no", epoch.String())
	}

	rsp, err := c.request(ctx, &res.Response, "GET", "/epoch_params", nil, opts)
	if err != nil {
		return
	}

	err = ReadAndUnmarshalResponse(rsp, &res.Response, &res.Data)

	if len(res.Data) == 0 {
		return nil, fmt.Errorf("%w: could not get epoch params %s", ErrResponse, epoch)
	}
	return
}

// GetEpochBlockProtocols returns the information about block protocol distribution in epoch.
func (c *Client) GetEpochBlockProtocols(
	ctx context.Context,
	epoch EpochNo,
	opts *RequestOptions,
) (res *EpochBlockProtocolsResponse, err error) {
	res = &EpochBlockProtocolsResponse{}
	if opts == nil {
		opts = c.NewRequestOptions()
	}
	if epoch > 0 {
		opts.QuerySet("_epoch_no", epoch.String())
	}

	rsp, err := c.request(ctx, &res.Response, "GET", "/epoch_block_protocols", nil, opts)
	if err != nil {
		return
	}

	err = ReadAndUnmarshalResponse(rsp, &res.Response, &res.Data)

	if len(res.Data) == 0 {
		return nil, fmt.Errorf("%w: could not get epoch block protocols %s", ErrResponse, epoch)
	}
	return
}
