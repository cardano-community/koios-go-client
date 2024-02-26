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

	"github.com/shopspring/decimal"
)

type (
	// Tip defines model for tip.
	Tip struct {
		// Absolute Slot number (slots not divided into epochs)
		AbsSlot Slot `json:"abs_slot"`

		// Block Height number on chain
		BlockNo BlockNo `json:"block_no"`

		// Timestamp for when the block was created
		BlockTime Timestamp `json:"block_time"`

		// EpochNo number
		EpochNo EpochNo `json:"epoch_no"`

		// Slot number within Epoch
		EpochSlot Slot `json:"epoch_slot"`

		// Block Hash in hex
		Hash BlockHash `json:"hash"`
	}

	// TipResponse response of /tip.
	TipResponse struct {
		Response
		Data Tip `json:"data"`
	}
)

// GetTip returns the tip info about the latest block seen by chain.
func (c *Client) GetTip(
	ctx context.Context,
	opts *RequestOptions,
) (res *TipResponse, err error) {
	res = &TipResponse{}
	rsp, err := c.request(ctx, &res.Response, "GET", "/tip", nil, opts)
	if err != nil {
		return res, err
	}
	tips := []Tip{}
	err = ReadAndUnmarshalResponse(rsp, &res.Response, &tips)
	if len(tips) == 1 {
		res.Data = tips[0]
	}
	return res, err
}

type (
	// Genesis defines model for genesis.
	Genesis struct {
		// Active Slot Co-Efficient (f) - determines the _probability_ of number of
		// slots in epoch that are expected to have blocks
		// (so mainnet, this would be: 432000 * 0.05 = 21600 estimated blocks).
		ActiveSlotCoeff decimal.Decimal `json:"activeslotcoeff"`

		// A JSON dump of Alonzo Genesis.
		AlonzoGenesis json.RawMessage `json:"alonzogenesis"`

		// Number of slots in an epoch.
		EpochLength decimal.Decimal `json:"epochlength"`

		// Number of KES key evolutions that will automatically occur before a KES
		// (hot) key is expired. This parameter is for security of a pool,
		// in case an operator had access to his hot(online) machine compromised.
		MaxKesRevolutions decimal.Decimal `json:"maxkesrevolutions"`

		// Maximum smallest units (lovelaces) supply for the blockchain.
		MaxLovelaceSupply decimal.Decimal `json:"maxlovelacesupply"`

		// Network ID used at various CLI identification to distinguish between
		// Mainnet and other networks.
		NetworkID string `json:"networkid"`

		// Unique network identifier for chain.
		NetworkMagic decimal.Decimal `json:"networkmagic"`

		// A unit (k) used to divide epochs to determine stability window
		// (used in security checks like ensuring atleast 1 block was
		// created in 3*k/f period, or to finalize next epoch's nonce
		// at 4*k/f slots before end of epoch).
		SecurityParam decimal.Decimal `json:"securityparam"`

		// Duration of a single slot (in seconds).
		SlotLength decimal.Decimal `json:"slotlength"`

		// Number of slots that represent a single KES period
		// (a unit used for validation of KES key evolutions).
		SlotsPerKesPeriod decimal.Decimal `json:"slotsperkesperiod"`

		// Timestamp for first block (genesis) on chain.
		SystemStart Timestamp `json:"systemstart"`

		// Number of BFT members that need to approve
		// (via vote) a Protocol Update Proposal.
		UpdateQuorum decimal.Decimal `json:"updatequorum"`
	}

	// GenesisResponse response of /genesis.
	GenesisResponse struct {
		Response
		Data Genesis `json:"data"`
	}
)

// GetGenesis returns the Genesis parameters used to start specific era on chain.
func (c *Client) GetGenesis(
	ctx context.Context,
	opts *RequestOptions,
) (*GenesisResponse, error) {
	res := &GenesisResponse{}
	rsp, err := c.request(ctx, &res.Response, "GET", "/genesis", nil, opts)
	if err != nil {
		return res, err
	}
	genesisres := []Genesis{}
	err = ReadAndUnmarshalResponse(rsp, &res.Response, &genesisres)
	if len(genesisres) == 1 {
		res.Data = genesisres[0]
	}
	return res, err
}

type (

	// Totals defines model for totals.
	Totals struct {

		// Circulating UTxOs for given epoch (in lovelaces).
		Circulation decimal.Decimal `json:"circulation"`

		// Epoch number.
		EpochNo EpochNo `json:"epoch_no"`

		// Total Reserves yet to be unlocked on chain.
		Reserves decimal.Decimal `json:"reserves"`

		// Rewards accumulated as of given epoch (in lovelaces).
		Reward decimal.Decimal `json:"reward"`

		// Total Active Supply (sum of treasury funds, rewards,
		// UTxOs, deposits and fees) for given epoch (in lovelaces).
		Supply decimal.Decimal `json:"supply"`

		// Funds in treasury for given epoch (in lovelaces).
		Treasury decimal.Decimal `json:"treasury"`
	}

	// TotalsResponse represents response from `/totals` endpoint.
	TotalsResponse struct {
		Response
		Data []Totals `json:"data"`
	}
)

// GetTotals returns the circulating utxo, treasury, rewards, supply and
// reserves in lovelace for specified epoch, all epochs if empty.
func (c *Client) GetTotals(
	ctx context.Context,
	epoch *EpochNo,
	opts *RequestOptions,
) (*TotalsResponse, error) {
	if opts == nil {
		opts = c.NewRequestOptions()
	}
	if epoch != nil {
		opts.QuerySet("_epoch_no", epoch.String())
	}
	res := &TotalsResponse{}
	rsp, err := c.request(ctx, &res.Response, "GET", "/totals", nil, opts)
	if err != nil {
		return res, err
	}
	return res, ReadAndUnmarshalResponse(rsp, &res.Response, &res.Data)
}

type (
	// ParamUpdate defines model for param update.
	ParamUpdate struct {
		// Epoch number.
		TxHash      TxHash          `json:"tx_hash"`
		BlockHeight BlockNo         `json:"block_height"`
		BlockTime   Timestamp       `json:"block_time"`
		EpochNo     EpochNo         `json:"epoch_no"`
		Data        json.RawMessage `json:"data"`
	}

	// ParamUpdatesResponse represents response from `/param_updates` endpoint.
	ParamUpdatesResponse struct {
		Response
		Data []ParamUpdate `json:"data"`
	}
)

// GetParamUpdates returns the parameter updates for the network.
func (c *Client) GetParamUpdates(
	ctx context.Context,
	opts *RequestOptions,
) (*ParamUpdatesResponse, error) {
	if opts == nil {
		opts = c.NewRequestOptions()
	}
	res := &ParamUpdatesResponse{}
	rsp, err := c.request(ctx, &res.Response, "GET", "/param_updates", nil, opts)
	if err != nil {
		return res, err
	}
	return res, ReadAndUnmarshalResponse(rsp, &res.Response, &res.Data)
}

type (
	// ReserveWithdrawal defines model for reserve withdrawal.
	ReserveWithdrawal struct {
		EpochNo      EpochNo         `json:"epoch_no"`
		EpochSlot    Slot            `json:"epoch_slot"`
		TxHash       TxHash          `json:"tx_hash"`
		BlockHash    BlockHash       `json:"block_hash"`
		BlockHeight  BlockNo         `json:"block_height"`
		Amount       decimal.Decimal `json:"amount"`
		StakeAddress Address         `json:"stake_address"`
	}

	// ReserveWithdrawalsResponse represents response from `/reserve_withdrawals` endpoint.
	ReserveWithdrawalsResponse struct {
		Response
		Data []ReserveWithdrawal `json:"data"`
	}
)

func (c *Client) GetReserveWithdrawals(
	ctx context.Context,
	opts *RequestOptions,
) (*ReserveWithdrawalsResponse, error) {
	if opts == nil {
		opts = c.NewRequestOptions()
	}

	res := &ReserveWithdrawalsResponse{}
	rsp, err := c.request(ctx, &res.Response, "GET", "/reserve_withdrawals", nil, opts)
	if err != nil {
		return res, err
	}
	return res, ReadAndUnmarshalResponse(rsp, &res.Response, &res.Data)
}

func (g *Genesis) AlonzoGenesisMap() (map[string]any, error) {
	var data map[string]any
	if err := json.Unmarshal(g.AlonzoGenesis, &data); err != nil {
		return nil, err
	}
	return data, nil
}

type (
	// TreasuryWithdrawal defines model for treasury withdrawal.
	TreasuryWithdrawal struct {
		EpochNo      EpochNo         `json:"epoch_no"`
		EpochSlot    Slot            `json:"epoch_slot"`
		TxHash       TxHash          `json:"tx_hash"`
		BlockHash    BlockHash       `json:"block_hash"`
		BlockHeight  BlockNo         `json:"block_height"`
		Amount       decimal.Decimal `json:"amount"`
		StakeAddress Address         `json:"stake_address"`
	}
	TreasuryWithdrawalsResponse struct {
		Response
		Data []TreasuryWithdrawal `json:"data"`
	}
)

func (c *Client) GetTreasuryWithdrawals(
	ctx context.Context,
	opts *RequestOptions,
) (*TreasuryWithdrawalsResponse, error) {
	if opts == nil {
		opts = c.NewRequestOptions()
	}

	res := &TreasuryWithdrawalsResponse{}
	rsp, err := c.request(ctx, &res.Response, "GET", "/treasury_withdrawals", nil, opts)
	if err != nil {
		return res, err
	}
	return res, ReadAndUnmarshalResponse(rsp, &res.Response, &res.Data)
}
