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
	"net/url"
)

type (

	// PoolListItem defines model for pool list item.
	PoolListItem struct {
		// PoolID Bech32 representation of pool ID.
		PoolID PoolID `json:"pool_id_bech32"`

		// Ticker of Pool.
		Ticker *string `json:"ticker,omitempty"`
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
		MetaJSON PoolMetaJSON `json:"meta_json"`
	}

	// Relay defines model for pool relay.
	Relay struct {
		// DNS name of the relay (nullable)
		DNS *string `json:"dns"`

		// IPv4 address of the relay (nullable)
		Ipv4 *string `json:"ipv4,"`

		// IPv6 address of the relay (nullable)
		Ipv6 *string `json:"ipv6,"`

		// Port number of the relay (nullable)
		Port *uint16 `json:"port"`

		// DNS service name of the relay (nullable)
		Srv *string `json:"srv"`
	}

	// PoolInfo defines model for pool_info.
	PoolInfo struct {
		// ActiveEpochNo Block number on chain where transaction was included.
		ActiveEpoch EpochNo `json:"active_epoch_no"`

		// ActiveStake Pool active stake.
		ActiveStake Lovelace `json:"active_stake"`

		// Total pool blocks on chain
		BlockCount uint64 `json:"block_count"`

		// FixedCost Pool fixed cost per epoch
		FixedCost Lovelace `json:"fixed_cost"`

		// LiveDelegators Pool live delegator count
		LiveDelegators uint64 `json:"live_delegators"`

		// LiveSaturation Pool live saturation (decimal format)
		LiveSaturation float32 `json:"live_saturation"`

		// LiveStake Pool live stake
		LiveStake Lovelace `json:"live_stake"`

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
		Owners []StakeAddress `json:"owners"`

		// Pledge pledge in lovelace
		Pledge Lovelace `json:"pledge"`

		// ID (bech32 format)
		ID PoolID `json:"pool_id_bech32"`

		// IDHex Pool ID (Hex format)
		IDHex string `json:"pool_id_hex"`

		// Pool status (registered | retiring | retired)
		Status string `json:"pool_status"`

		// Announced retiring epoch (nullable)
		RetiringEpoch *EpochNo `json:"retiring_epoch"`

		// Pool reward address
		RewardAddr StakeAddress `json:"reward_addr"`

		// Pool VRF key hash
		VrfKeyHash string `json:"vrf_key_hash"`

		// Relays of the pool
		Relays []Relay `json:"relays"`
	}

	// PoolUpdateInfo response item from `/pool_updates`.
	PoolUpdateInfo struct {
		// TxHash update transaction
		TxHash TxHash `json:"tx_hash"`

		// Time time of the block.
		BlockTime string `json:"block_time"`

		// ID (bech32 format)
		ID PoolID `json:"pool_id_bech32"`

		// IDHex Pool ID (Hex format)
		IDHex string `json:"pool_id_hex"`

		// ActiveEpochNo Block number on chain where transaction was included.
		ActiveEpoch EpochNo `json:"active_epoch_no"`

		// // ActiveStake Pool active stake.
		// ActiveStake Lovelace `json:"active_stake"`

		// // Total pool blocks on chain
		// BlockCount uint64 `json:"block_count"`

		// FixedCost Pool fixed cost per epoch
		FixedCost Lovelace `json:"fixed_cost"`

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
		Owners []StakeAddress `json:"owners"`

		// Pledge pledge in lovelace.
		Pledge Lovelace `json:"pledge"`

		// Pool status (registered | retiring | retired).
		Status string `json:"pool_status"`

		// Announced retiring epoch (nullable).
		RetiringEpoch *EpochNo `json:"retiring_epoch"`

		// Pool reward address.
		RewardAddr StakeAddress `json:"reward_addr"`

		// Pool VRF key hash.
		VrfKeyHash string `json:"vrf_key_hash"`

		// Relays of the pool.
		Relays []PoolRelays `json:"relays"`
	}

	// PoolDelegator info.
	PoolDelegator struct {
		StakeAddress StakeAddress `json:"stake_address"`
		Amount       Lovelace     `json:"amount"`
		Epoch        EpochNo      `json:"epoch_no"`
	}

	// PoolRelays list item.
	PoolRelays struct {
		PoolID PoolID `json:"pool_id_bech32"`
		// Relays of the pool.
		Relays []Relay `json:"relays"`
	}

	// PoolBlockInfo block info.
	PoolBlockInfo struct {
		// Hash block hash
		Hash BlockHash `json:"block_hash"`

		// Epoch number.
		Epoch EpochNo `json:"epoch_no"`

		// EpochSlot slot number within epoch.
		EpochSlot int `json:"epoch_slot_no"`

		// Slot is overall slot number (slots from genesis block of chain).
		Slot int `json:"slot_no"`

		// Time time of the block.
		Time string `json:"block_time"`

		// BlockNo is block number on chain.
		BlockNo uint64 `json:"block_no"`
	}

	// PoolListResponse represents response from `/pool_list` endpoint.
	PoolListResponse struct {
		Response
		Data []PoolListItem `json:"response"`
	}

	// PoolInfosResponse represents response from `/pool_info` endpoint.
	PoolInfosResponse struct {
		Response
		Data []PoolInfo `json:"response"`
	}

	// PoolInfoResponse represents response from `/pool_info` endpoint.
	// when requesting info about single pool.
	PoolInfoResponse struct {
		Response
		Data *PoolInfo `json:"response"`
	}

	// PoolDelegatorsResponse represents response from `/pool_delegators` endpoint.
	PoolDelegatorsResponse struct {
		Response
		Data []PoolDelegator `json:"response"`
	}

	// PoolBlocksResponse represents response from `/pool_blocks` endpoint.
	PoolBlocksResponse struct {
		Response
		Data []PoolBlockInfo `json:"response"`
	}

	// PoolUpdatesResponse represents response from `/pool_updates` endpoint.
	PoolUpdatesResponse struct {
		Response
		Data []PoolUpdateInfo `json:"response"`
	}

	// PoolRelaysResponse represents response from `/pool_relays` endpoint.
	PoolRelaysResponse struct {
		Response
		Data []PoolRelays `json:"response"`
	}

	// PoolMetadataResponse represents response from `/pool_metadata` endpoint.
	PoolMetadataResponse struct {
		Response
		Data []PoolMetadata `json:"response"`
	}
)

// GetPoolList returns the list of all currently registered/retiring (not retired) pools.
func (c *Client) GetPoolList(ctx context.Context) (res *PoolListResponse, err error) {
	res = &PoolListResponse{}
	rsp, err := c.request(ctx, &res.Response, "GET", "/pool_list", nil, nil, nil)
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

	return res, nil
}

// GetPoolInfo returns current pool status and details for a specified pool.
func (c *Client) GetPoolInfo(ctx context.Context, pid PoolID) (res *PoolInfoResponse, err error) {
	res = &PoolInfoResponse{}
	rsp, err := c.GetPoolInfos(ctx, []PoolID{pid})
	res.Response = rsp.Response
	if len(rsp.Data) == 1 {
		res.Data = &rsp.Data[0]
	}
	return
}

// GetPoolInfos returns current pool statuses and details for a specified list of pool ids.
func (c *Client) GetPoolInfos(ctx context.Context, pids []PoolID) (res *PoolInfosResponse, err error) {
	res = &PoolInfosResponse{}
	if len(pids) == 0 {
		err = ErrNoPoolID
		res.applyError(nil, err)
		return
	}

	rsp, err := c.request(ctx, &res.Response, "POST", "/pool_info", poolIdsPL(pids), nil, nil)
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

// GetPoolDelegators returns information about delegators
// by a given pool and optional epoch (current if omitted).
func (c *Client) GetPoolDelegators(
	ctx context.Context,
	pid PoolID,
	epoch *EpochNo,
) (res *PoolDelegatorsResponse, err error) {
	res = &PoolDelegatorsResponse{}

	params := url.Values{}
	params.Set("_pool_bech32", string(pid))
	if epoch != nil {
		params.Set("_epoch_no", fmt.Sprint(*epoch))
	}
	rsp, err := c.request(ctx, &res.Response, "GET", "/pool_delegators", nil, params, nil)
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

// GetPoolBlocks returns information about blocks minted by a given pool
// in current epoch (or _epoch_no if provided).
func (c *Client) GetPoolBlocks(
	ctx context.Context,
	pid PoolID,
	epoch *EpochNo,
) (res *PoolBlocksResponse, err error) {
	res = &PoolBlocksResponse{}

	params := url.Values{}
	params.Set("_pool_bech32", string(pid))
	if epoch != nil {
		params.Set("_epoch_no", fmt.Sprint(*epoch))
	}
	rsp, err := c.request(ctx, &res.Response, "GET", "/pool_blocks", nil, params, nil)
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

// GetPoolUpdates returns all pool updates for all pools or
// only updates for specific pool if specified.
func (c *Client) GetPoolUpdates(
	ctx context.Context,
	pid *PoolID,
) (res *PoolUpdatesResponse, err error) {
	res = &PoolUpdatesResponse{}

	params := url.Values{}
	if pid != nil {
		params.Set("_pool_bech32", fmt.Sprint(*pid))
	}

	rsp, err := c.request(ctx, &res.Response, "GET", "/pool_updates", nil, params, nil)
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

// GetPoolRelays returns a list of registered relays
// for all currently registered/retiring (not retired) pools.
func (c *Client) GetPoolRelays(ctx context.Context) (res *PoolRelaysResponse, err error) {
	res = &PoolRelaysResponse{}

	rsp, err := c.request(ctx, &res.Response, "GET", "/pool_relays", nil, nil, nil)
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

// GetPoolMetadata returns Metadata(on & off-chain)
// for all currently registered/retiring (not retired) pools.
func (c *Client) GetPoolMetadata(ctx context.Context) (res *PoolMetadataResponse, err error) {
	res = &PoolMetadataResponse{}

	rsp, err := c.request(ctx, &res.Response, "GET", "/pool_metadata", nil, nil, nil)
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
