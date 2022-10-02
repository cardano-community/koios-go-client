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
	"io"

	"github.com/shopspring/decimal"
)

// introduces breaking change since v1.3.0

type (

	// PoolListItem defines model for pool list item.
	PoolListItem struct {
		// PoolID Bech32 representation of pool ID.
		PoolID PoolID `json:"pool_id_bech32"`

		// Ticker of Pool.
		Ticker string `json:"ticker,omitempty"`
	}

	// PoolMetaJSON pool meadata json.
	PoolMetaJSON struct {
		// Pool description
		Description *string `json:"description"`

		// Pool homepage URL
		Homepage *string `json:"homepage"`

		// Pool name
		Name *string `json:"name"`

		// Pool ticker
		Ticker *string `json:"ticker"`
	}

	// PoolMetadata metadata list item.
	PoolMetadata struct {
		// ID (bech32 format)
		PoolID PoolID `json:"pool_id_bech32"`
		// MetaUrl Pool metadata URL
		MetaURL string `json:"meta_url"`

		// MetaHash Pool metadata hash
		MetaHash string `json:"meta_hash"`

		// MetaJson pool meta json
		MetaJSON *PoolMetaJSON `json:"meta_json,omitempty"`
	}

	// Relay defines model for pool relay.
	Relay struct {
		// DNS name of the relay (nullable)
		DNS string `json:"dns"`

		// IPv4 address of the relay (nullable)
		Ipv4 string `json:"ipv4,"`

		// IPv6 address of the relay (nullable)
		Ipv6 string `json:"ipv6,"`

		// Port number of the relay (nullable)
		Port uint16 `json:"port"`

		// DNS service name of the relay (nullable)
		Srv string `json:"srv"`
	}

	// PoolInfo defines model for pool_info.
	PoolInfo struct {
		// ActiveEpochNo Block number on chain where transaction was included.
		ActiveEpoch EpochNo `json:"active_epoch_no"`

		// ActiveStake Pool active stake.
		ActiveStake decimal.Decimal `json:"active_stake"`

		// Total pool blocks on chain
		BlockCount uint64 `json:"block_count"`

		// FixedCost Pool fixed cost per epoch
		FixedCost decimal.Decimal `json:"fixed_cost"`

		// LiveDelegators Pool live delegator count
		LiveDelegators uint64 `json:"live_delegators"`

		// LiveSaturation Pool live saturation (decimal format)
		LiveSaturation float32 `json:"live_saturation"`

		// LiveStake Pool live stake
		LiveStake decimal.Decimal `json:"live_stake"`

		// LivePledge Pool live pledge
		LivePledge decimal.Decimal `json:"live_pledge"`

		// Margin (decimal format)
		Margin float32 `json:"margin"`

		// MetaHash Pool metadata hash
		MetaHash string `json:"meta_hash"`

		// MetaJson pool meta json
		MetaJSON PoolMetaJSON `json:"meta_json"`

		// MetaUrl Pool metadata URL
		MetaURL string `json:"meta_url"`

		// OpCert Pool latest operational certificate hash
		OpCert string `json:"op_cert"`

		// OpCertCounter Pool latest operational certificate counter value
		OpCertCounter int `json:"op_cert_counter"`

		// Owners of the pool
		Owners []Address `json:"owners"`

		// Pledge pledge in lovelace
		Pledge decimal.Decimal `json:"pledge"`

		// ID (bech32 format)
		PoolID PoolID `json:"pool_id_bech32"`

		// IDHex Pool ID (Hex format)
		PoolIdHex string `json:"pool_id_hex"`

		// Pool status (registered | retiring | retired)
		PoolStatus string `json:"pool_status"`

		// Announced retiring epoch (nullable)
		RetiringEpoch *EpochNo `json:"retiring_epoch"`

		// Pool reward address
		RewardAddr Address `json:"reward_addr"`

		// Pool VRF key hash
		VrfKeyHash string `json:"vrf_key_hash"`

		// Relays of the pool
		Relays []Relay `json:"relays"`

		Sigma float64 `json:"sigma"`
	}

	// PoolUpdateInfo response item from `/pool_updates`.
	PoolUpdateInfo struct {
		// TxHash update transaction
		TxHash TxHash `json:"tx_hash"`

		// Time time of the block.
		BlockTime Timestamp `json:"block_time"`

		// ID (bech32 format)
		PoolID PoolID `json:"pool_id_bech32"`

		// IDHex Pool ID (Hex format)
		PoolIDHex string `json:"pool_id_hex"`

		// ActiveEpochNo Block number on chain where transaction was included.
		ActiveEpoch EpochNo `json:"active_epoch_no"`

		// // ActiveStake Pool active stake.
		// ActiveStake Lovelace `json:"active_stake"`

		// // Total pool blocks on chain
		// BlockCount uint64 `json:"block_count"`

		// FixedCost Pool fixed cost per epoch
		FixedCost decimal.Decimal `json:"fixed_cost"`

		// // LiveDelegators Pool live delegator count
		// LiveDelegators uint64 `json:"live_delegators"`

		// // LiveSaturation Pool live saturation (decimal format)
		// LiveSaturation float32 `json:"live_saturation"`

		// // LiveStake Pool live stake
		// LiveStake Lovelace `json:"live_stake"`

		// Margin (decimal format)
		Margin float32 `json:"margin"`

		// MetaHash Pool metadata hash
		MetaHash string `json:"meta_hash"`

		// // MetaJson pool meta json
		// MetaJSON PoolMeta `json:"meta_json"`

		// MetaUrl Pool metadata URL
		MetaURL string `json:"meta_url"`

		// // OpCert Pool latest operational certificate hash
		// OpCert string `json:"op_cert"`

		// // OpCertCounter Pool latest operational certificate counter value
		// OpCertCounter int `json:"op_cert_counter"`

		// Owners of the pool.
		Owners []Address `json:"owners"`

		// Pledge pledge in lovelace.
		Pledge decimal.Decimal `json:"pledge"`

		// Pool status (registered | retiring | retired).
		PoolStatus string `json:"pool_status"`

		// Announced retiring epoch (nullable).
		RetiringEpoch *EpochNo `json:"retiring_epoch.omitempty"`

		// Pool reward address.
		RewardAddr Address `json:"reward_addr"`

		// Pool VRF key hash.
		VrfKeyHash string `json:"vrf_key_hash"`

		// Relays of the pool.
		Relays []Relay `json:"relays"`
	}

	// PoolDelegator info.
	PoolDelegator struct {
		StakeAddress           Address         `json:"stake_address"`
		Amount                 decimal.Decimal `json:"amount"`
		ActiveEpochNo          EpochNo         `json:"active_epoch_no"`
		LatestDelegationTxHash TxHash          `json:"latest_delegation_tx_hash"`
	}

	PoolDelegatorHistory struct {
		StakeAddress Address         `json:"stake_address"`
		Amount       decimal.Decimal `json:"amount"`
		EpochNo      EpochNo         `json:"epoch_no"`
	}

	// PoolRelays list item.
	PoolRelays struct {
		PoolID PoolID `json:"pool_id_bech32"`
		// Relays of the pool.
		Relays []Relay `json:"relays"`
	}

	// PoolBlockInfo block info.
	PoolBlockInfo struct {
		// Slot is overall slot number (slots from genesis block of chain).
		AbsSlot Slot `json:"abs_slot"`

		// Hash block hash
		BlockHash BlockHash `json:"block_hash"`

		// BlockHeight ogf the block
		BlockHeight uint64 `json:"block_height"`

		// Time time of the block.
		BlockTime Timestamp `json:"block_time"`

		// Epoch number.
		EpochNo EpochNo `json:"epoch_no"`

		// EpochSlot slot number within epoch.
		EpochSlot Slot `json:"epoch_slot"`
	}

	// PoolHistory entry.
	PoolHistory struct {
		// Epoch number.
		EpochNo EpochNo `json:"epoch_no"`
		// ActiveStake Pool active stake.
		ActiveStake    decimal.Decimal `json:"active_stake"`
		ActiveStakePCT float64         `json:"active_stake_pct"`
		SaturationPCT  float64         `json:"saturation_pct"`
		BlockCNT       int             `json:"block_cnt"`
		DelegatorCNT   int             `json:"delegator_cnt"`
		Margin         float64         `json:"margin"`
		FixedCost      decimal.Decimal `json:"fixed_cost"`
		PoolFees       decimal.Decimal `json:"pool_fees"`
		DelegRewards   decimal.Decimal `json:"deleg_rewards"`
		EpochROS       float64         `json:"epoch_ros"`
	}

	PoolSnapshot struct {
		Snapshot    string          `json:"snapshot"`
		EpochNo     EpochNo         `json:"epoch_no"`
		Nonce       string          `json:"nonce"`
		PoolStake   decimal.Decimal `json:"pool_stake"`
		ActiveStake decimal.Decimal `json:"active_stake"`
	}

	PoolSnapshotResponse struct {
		Response
		Data []PoolSnapshot `json:"data"`
	}

	// PoolListResponse represents response from `/pool_list` endpoint.
	PoolListResponse struct {
		Response
		Data []PoolListItem `json:"data"`
	}

	// PoolInfosResponse represents response from `/pool_info` endpoint.
	PoolInfosResponse struct {
		Response
		Data []PoolInfo `json:"data"`
	}

	// PoolInfoResponse represents response from `/pool_info` endpoint.
	// when requesting info about single pool.
	PoolInfoResponse struct {
		Response
		Data *PoolInfo `json:"data"`
	}

	// PoolDelegatorsResponse represents response from `/pool_delegators` endpoint.
	PoolDelegatorsResponse struct {
		Response
		Data []PoolDelegator `json:"data"`
	}

	PoolDelegatorsHistoryResponse struct {
		Response
		Data []PoolDelegatorHistory `json:"data"`
	}

	// PoolBlocksResponse represents response from `/pool_blocks` endpoint.
	PoolBlocksResponse struct {
		Response
		Data []PoolBlockInfo `json:"data"`
	}

	// PoolUpdatesResponse represents response from `/pool_updates` endpoint.
	PoolUpdatesResponse struct {
		Response
		Data []PoolUpdateInfo `json:"data"`
	}

	// PoolRelaysResponse represents response from `/pool_relays` endpoint.
	PoolRelaysResponse struct {
		Response
		Data []PoolRelays `json:"data"`
	}

	// PoolMetadataResponse represents response from `/pool_metadata` endpoint.
	PoolMetadataResponse struct {
		Response
		Data []PoolMetadata `json:"data"`
	}

	// PoolHistoryResponse represents response from `/pool_history` endpoint.
	PoolHistoryResponse struct {
		Response
		Data []PoolHistory `json:"data"`
	}
)

// GetPoolList returns the list of all currently registered/retiring (not retired) pools.
func (c *Client) GetPools(
	ctx context.Context,
	opts *RequestOptions,
) (res *PoolListResponse, err error) {
	res = &PoolListResponse{}
	rsp, err := c.request(ctx, &res.Response, "GET", "/pool_list", nil, opts)
	if err != nil {
		return
	}
	err = ReadAndUnmarshalResponse(rsp, &res.Response, &res.Data)
	return
}

// GetPoolInfo returns current pool status and details for a specified pool.
func (c *Client) GetPoolInfo(
	ctx context.Context,
	pid PoolID,
	opts *RequestOptions,
) (res *PoolInfoResponse, err error) {
	res = &PoolInfoResponse{}
	rsp, err := c.GetPoolInfos(ctx, []PoolID{pid}, opts)
	res.Response = rsp.Response
	if len(rsp.Data) == 1 {
		res.Data = &rsp.Data[0]
	}
	return
}

// GetPoolInfos returns current pool statuses and details for a specified list of pool ids.
func (c *Client) GetPoolInfos(
	ctx context.Context,
	pids []PoolID,
	opts *RequestOptions,
) (res *PoolInfosResponse, err error) {
	res = &PoolInfosResponse{}
	if len(pids) == 0 {
		err = ErrNoPoolID
		res.applyError(nil, err)
		return
	}

	rsp, err := c.request(ctx, &res.Response, "POST", "/pool_info", poolIdsPL(pids), opts)
	if err != nil {
		return
	}
	err = ReadAndUnmarshalResponse(rsp, &res.Response, &res.Data)
	return
}

func (c *Client) GetPoolSnapshot(
	ctx context.Context,
	pid PoolID,
	opts *RequestOptions,
) (res *PoolSnapshotResponse, err error) {
	res = &PoolSnapshotResponse{}
	if opts == nil {
		opts = c.NewRequestOptions()
	}
	opts.QuerySet("_pool_bech32", pid.String())

	rsp, err := c.request(ctx, &res.Response, "GET", "/pool_stake_snapshot", nil, opts)
	if err != nil {
		return
	}
	err = ReadAndUnmarshalResponse(rsp, &res.Response, &res.Data)
	return
}

// GetPoolDelegators returns information about delegators
// by a given pool and optional epoch (current if omitted).
func (c *Client) GetPoolDelegators(
	ctx context.Context,
	pid PoolID,
	opts *RequestOptions,
) (res *PoolDelegatorsResponse, err error) {
	res = &PoolDelegatorsResponse{}

	if opts == nil {
		opts = c.NewRequestOptions()
	}
	opts.QuerySet("_pool_bech32", pid.String())

	rsp, err := c.request(ctx, &res.Response, "GET", "/pool_delegators", nil, opts)
	if err != nil {
		return
	}
	err = ReadAndUnmarshalResponse(rsp, &res.Response, &res.Data)

	return
}
func (c *Client) GetPoolDelegatorsHistory(
	ctx context.Context,
	pid PoolID,
	epoch *EpochNo,
	opts *RequestOptions,
) (res *PoolDelegatorsHistoryResponse, err error) {
	res = &PoolDelegatorsHistoryResponse{}

	if opts == nil {
		opts = c.NewRequestOptions()
	}
	opts.QuerySet("_pool_bech32", pid.String())
	if epoch != nil {
		opts.QuerySet("_epoch_no", epoch.String())
	}

	rsp, err := c.request(ctx, &res.Response, "GET", "/pool_delegators_history", nil, opts)
	if err != nil {
		return
	}
	err = ReadAndUnmarshalResponse(rsp, &res.Response, &res.Data)

	return
}

// GetPoolBlocks returns information about blocks minted by a given pool
// in current epoch (or _epoch_no if provided).
func (c *Client) GetPoolBlocks(
	ctx context.Context,
	pid PoolID,
	epoch *EpochNo,
	opts *RequestOptions,
) (res *PoolBlocksResponse, err error) {
	res = &PoolBlocksResponse{}

	if opts == nil {
		opts = c.NewRequestOptions()
	}
	opts.QuerySet("_pool_bech32", pid.String())
	if epoch != nil {
		opts.QuerySet("_epoch_no", epoch.String())
	}
	rsp, err := c.request(ctx, &res.Response, "GET", "/pool_blocks", nil, opts)
	if err != nil {
		return
	}
	err = ReadAndUnmarshalResponse(rsp, &res.Response, &res.Data)
	return
}

// GetPoolUpdates returns all pool updates for all pools or
// only updates for specific pool if specified.
func (c *Client) GetPoolUpdates(
	ctx context.Context,
	pid *PoolID,
	opts *RequestOptions,
) (res *PoolUpdatesResponse, err error) {
	res = &PoolUpdatesResponse{}

	if opts == nil {
		opts = c.NewRequestOptions()
	}
	if pid != nil {
		opts.QuerySet("_pool_bech32", pid.String())
	}

	rsp, err := c.request(ctx, &res.Response, "GET", "/pool_updates", nil, opts)
	if err != nil {
		return
	}
	err = ReadAndUnmarshalResponse(rsp, &res.Response, &res.Data)
	return
}

// GetPoolRelays returns a list of registered relays
// for all currently registered/retiring (not retired) pools.
func (c *Client) GetPoolRelays(
	ctx context.Context,
	opts *RequestOptions,
) (res *PoolRelaysResponse, err error) {
	res = &PoolRelaysResponse{}

	rsp, err := c.request(ctx, &res.Response, "GET", "/pool_relays", nil, opts)
	if err != nil {
		return
	}
	err = ReadAndUnmarshalResponse(rsp, &res.Response, &res.Data)
	return
}

// GetPoolMetadata returns Metadata(on & off-chain)
// for all currently registered/retiring (not retired) pools.
func (c *Client) GetPoolMetadata(
	ctx context.Context,
	pids []PoolID,
	opts *RequestOptions,
) (res *PoolMetadataResponse, err error) {
	res = &PoolMetadataResponse{}

	rsp, err := c.request(ctx, &res.Response, "POST", "/pool_metadata", poolIdsPL(pids), opts)
	if err != nil {
		return
	}
	err = ReadAndUnmarshalResponse(rsp, &res.Response, &res.Data)
	return
}

// GetPoolHistory returns information about pool stake, block and reward history
// in a given epoch _epoch_no (or all epochs that pool existed for, in descending
// order if no _epoch_no was provided).
func (c *Client) GetPoolHistory(
	ctx context.Context,
	pid PoolID,
	epoch *EpochNo,
	opts *RequestOptions,
) (res *PoolHistoryResponse, err error) {
	res = &PoolHistoryResponse{}

	if opts == nil {
		opts = c.NewRequestOptions()
	}
	opts.QuerySet("_pool_bech32", pid.String())
	if epoch != nil {
		opts.QuerySet("_epoch_no", epoch.String())
	}
	rsp, err := c.request(ctx, &res.Response, "GET", "/pool_history", nil, opts)
	if err != nil {
		return
	}
	err = ReadAndUnmarshalResponse(rsp, &res.Response, &res.Data)
	return
}

func poolIdsPL(pids []PoolID) io.Reader {
	var payload = struct {
		PIDS []PoolID `json:"_pool_bech32_ids"`
	}{pids}
	rpipe, w := io.Pipe()
	go func() {
		_ = json.NewEncoder(w).Encode(payload)
		defer w.Close()
	}()
	return rpipe
}
