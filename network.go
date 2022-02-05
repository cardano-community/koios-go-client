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
	// Tip defines model for tip.
	Tip struct {
		// Absolute Slot number (slots not divided into epochs)
		AbsSlot int `json:"abs_slot"`

		// Block Height number on chain
		BlockNo int `json:"block_no"`

		// Timestamp for when the block was created
		BlockTime string `json:"block_time"`

		// Epoch number
		Epoch int `json:"epoch"`

		// Slot number within Epoch
		EpochSlot int `json:"epoch_slot"`

		// Block Hash in hex
		Hash string `json:"hash"`
	}

	// TipResponse response of /tip.
	TipResponse struct {
		Response
		Tip *Tip `json:"response,omitempty"`
	}

	// Genesis defines model for genesis.
	Genesis struct {
		// Active Slot Co-Efficient (f) - determines the _probability_ of number of
		// slots in epoch that are expected to have blocks
		// (so mainnet, this would be: 432000 * 0.05 = 21600 estimated blocks).
		Activeslotcoeff string `json:"activeslotcoeff"`

		// A JSON dump of Alonzo Genesis.
		Alonzogenesis string `json:"alonzogenesis"`

		// Number of slots in an epoch.
		Epochlength string `json:"epochlength"`

		// Number of KES key evolutions that will automatically occur before a KES
		// (hot) key is expired. This parameter is for security of a pool,
		// in case an operator had access to his hot(online) machine compromised.
		Maxkesrevolutions string `json:"maxkesrevolutions"`

		// Maximum smallest units (lovelaces) supply for the blockchain.
		Maxlovelacesupply string `json:"maxlovelacesupply"`

		// Network ID used at various CLI identification to distinguish between
		// Mainnet and other networks.
		Networkid string `json:"networkid"`

		// Unique network identifier for chain.
		Networkmagic string `json:"networkmagic"`

		// A unit (k) used to divide epochs to determine stability window
		// (used in security checks like ensuring atleast 1 block was
		// created in 3*k/f period, or to finalize next epoch's nonce
		// at 4*k/f slots before end of epoch).
		Securityparam string `json:"securityparam"`

		// Duration of a single slot (in seconds).
		Slotlength string `json:"slotlength"`

		// Number of slots that represent a single KES period
		// (a unit used for validation of KES key evolutions).
		Slotsperkesperiod string `json:"slotsperkesperiod"`

		// Timestamp for first block (genesis) on chain.
		Systemstart string `json:"systemstart"`

		// Number of BFT members that need to approve
		// (via vote) a Protocol Update Proposal.
		Updatequorum string `json:"updatequorum"`
	}

	// GenesisResponse response of /genesis.
	GenesisResponse struct {
		Response
		Genesis *Genesis `json:"response,omitempty"`
	}

	// Totals defines model for totals.
	Totals struct {

		// Circulating UTxOs for given epoch (in lovelaces).
		Circulation Lovelace `json:"circulation"`

		// Epoch number.
		EpochNo EpochNo `json:"epoch_no"`

		// Total Reserves yet to be unlocked on chain.
		Reserves Lovelace `json:"reserves"`

		// Rewards accumulated as of given epoch (in lovelaces).
		Reward Lovelace `json:"reward"`

		// Total Active Supply (sum of treasury funds, rewards,
		// UTxOs, deposits and fees) for given epoch (in lovelaces).
		Supply Lovelace `json:"supply"`

		// Funds in treasury for given epoch (in lovelaces).
		Treasury Lovelace `json:"treasury"`
	}

	// TotalsResponse represents response from `/totals` endpoint.
	TotalsResponse struct {
		Response
		Totals []Totals `json:"response,omitempty"`
	}
)

// GetTip returns the tip info about the latest block seen by chain.
func (c *Client) GetTip(ctx context.Context) (*TipResponse, error) {
	rsp, err := c.GET(ctx, "/tip")
	if err != nil {
		return nil, err
	}
	res := &TipResponse{}
	res.setStatus(rsp)
	body, err := readResponseBody(rsp)
	if err != nil {
		return nil, err
	}

	tips := []Tip{}
	if err := json.Unmarshal(body, &tips); err != nil {
		res.applyError(body, err)
		return res, nil
	}
	if len(tips) == 1 {
		res.Tip = &tips[0]
	}
	return res, nil
}

// GetGenesis returns the Genesis parameters used to start specific era on chain.
func (c *Client) GetGenesis(ctx context.Context) (*GenesisResponse, error) {
	rsp, err := c.GET(ctx, "/genesis")
	if err != nil {
		return nil, err
	}
	res := &GenesisResponse{}
	res.setStatus(rsp)
	body, err := readResponseBody(rsp)
	if err != nil {
		return nil, err
	}

	genesisres := []Genesis{}

	if err := json.Unmarshal(body, &genesisres); err != nil {
		res.applyError(body, err)
		return res, nil
	}

	if len(genesisres) == 1 {
		res.Genesis = &genesisres[0]
	}
	return res, nil
}

// GetTotals returns the circulating utxo, treasury, rewards, supply and
// reserves in lovelace for specified epoch, all epochs if empty.
func (c *Client) GetTotals(ctx context.Context, epochNo *EpochNo) (*TotalsResponse, error) {
	params := url.Values{}
	if epochNo != nil {
		params.Set("_epoch_no", fmt.Sprint(*epochNo))
	}

	rsp, err := c.GET(ctx, "/totals", params)
	if err != nil {
		return nil, err
	}
	res := &TotalsResponse{}
	res.setStatus(rsp)
	body, err := readResponseBody(rsp)
	if err != nil {
		return nil, err
	}

	totals := []Totals{}
	if err := json.Unmarshal(body, &totals); err != nil {
		res.applyError(body, err)
		return res, nil
	}
	if len(totals) > 0 {
		res.Totals = totals
	}
	return res, nil
}
