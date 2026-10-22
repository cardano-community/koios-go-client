// SPDX-License-Identifier: Apache-2.0
//
// Copyright © 2022 The Cardano Community Authors

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

		// Era is the Cardano era of the tip block.
		Era string `json:"era"`
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

		// AlonzoGenesis is the Alonzo genesis delivered by the API as a
		// JSON document encoded in a string. See AlonzoGenesisMap.
		AlonzoGenesis string `json:"alonzogenesis"`

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

		// Fees is the amount (in lovelace) in the fee pot.
		Fees decimal.Decimal `json:"fees"`

		// DepositsStake is the amount (in lovelace) in the obligation pot
		// coming from stake key and pool deposits.
		DepositsStake decimal.Decimal `json:"deposits_stake"`

		// DepositsDRep is the amount (in lovelace) in the obligation pot
		// coming from DRep registration deposits.
		DepositsDRep decimal.Decimal `json:"deposits_drep"`

		// DepositsProposal is the amount (in lovelace) in the obligation pot
		// coming from governance proposal deposits.
		DepositsProposal decimal.Decimal `json:"deposits_proposal"`

		// TreasuryDonation is the amount (in lovelace) donated to the treasury.
		TreasuryDonation decimal.Decimal `json:"treasury_donation"`

		// TreasuryWithdrawal is the amount (in lovelace) withdrawn from the treasury.
		TreasuryWithdrawal decimal.Decimal `json:"treasury_withdrawal"`

		// ReservesWithdrawal is the amount (in lovelace) withdrawn from reserves.
		ReservesWithdrawal decimal.Decimal `json:"reserves_withdrawal"`
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

		// EarnedEpoch is the epoch where the amount is earned.
		EarnedEpoch EpochNo `json:"earned_epoch"`

		// SpendableEpoch is the epoch where the earned amount can be spent.
		SpendableEpoch EpochNo `json:"spendable_epoch"`
	}

	// ReserveWithdrawalsResponse represents response from `/reserve_withdrawals` endpoint.
	ReserveWithdrawalsResponse struct {
		Response
		Data []ReserveWithdrawal `json:"data"`
	}
)

// GetReserveWithdrawals returns all withdrawals from reserves against stake accounts.
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

// AlonzoGenesisMap decodes the JSON document held in AlonzoGenesis.
func (g *Genesis) AlonzoGenesisMap() (map[string]any, error) {
	var data map[string]any
	if err := json.Unmarshal([]byte(g.AlonzoGenesis), &data); err != nil {
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

		// EarnedEpoch is the epoch where the amount is earned.
		EarnedEpoch EpochNo `json:"earned_epoch"`

		// SpendableEpoch is the epoch where the earned amount can be spent.
		SpendableEpoch EpochNo `json:"spendable_epoch"`
	}
	TreasuryWithdrawalsResponse struct {
		Response
		Data []TreasuryWithdrawal `json:"data"`
	}
)

// GetTreasuryWithdrawals returns all withdrawals from treasury against stake accounts.
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

type (
	// EraSummary defines model for era_summaries.
	EraSummary struct {
		// Era is the Cardano era name.
		Era string `json:"era"`

		// ProtocolMajor is the protocol major version (nullable).
		ProtocolMajor *int `json:"protocol_major"`

		// ProtocolMinor is the protocol minor version (nullable).
		ProtocolMinor *int `json:"protocol_minor"`

		// LedgerProtocol names the collection of features sitting between
		// the consensus and ledger layers.
		LedgerProtocol string `json:"ledger_protocol"`

		// ConsensusMechanism names the collection of consensus features
		// introduced at a hard fork.
		ConsensusMechanism string `json:"consensus_mechanism"`

		// Notes is additional information about when this protocol update
		// was pushed, such as an associated fork name.
		Notes string `json:"notes"`

		// EpochNo is the epoch in which the era starts.
		EpochNo EpochNo `json:"epoch_no"`

		// FirstBlockTime is the time of the first block of the era.
		FirstBlockTime Timestamp `json:"first_block_time"`

		// FirstBlockHash is the hash of the first block of the era.
		FirstBlockHash BlockHash `json:"first_block_hash"`
	}

	// EraSummariesResponse represents response from `/era_summaries` endpoint.
	EraSummariesResponse struct {
		Response
		Data []EraSummary `json:"data"`
	}

	// CliProtocolParamsResponse represents response from `/cli_protocol_params` endpoint.
	CliProtocolParamsResponse struct {
		Response

		// Data is the raw protocol parameters JSON as published by cardano-cli.
		// Its schema varies with the cardano-node version and era.
		Data json.RawMessage `json:"data"`
	}
)

// GetEraSummaries returns information about each era on the network
// along with its start epoch.
func (c *Client) GetEraSummaries(
	ctx context.Context,
	opts *RequestOptions,
) (*EraSummariesResponse, error) {
	res := &EraSummariesResponse{}
	rsp, err := c.request(ctx, &res.Response, "GET", "/era_summaries", nil, opts)
	if err != nil {
		return res, err
	}
	return res, ReadAndUnmarshalResponse(rsp, &res.Response, &res.Data)
}

// GetCliProtocolParams returns the current protocol parameters as published
// by cardano-cli. The output schema depends on the cardano-node version and
// era, so it is returned as raw JSON.
func (c *Client) GetCliProtocolParams(
	ctx context.Context,
	opts *RequestOptions,
) (*CliProtocolParamsResponse, error) {
	res := &CliProtocolParamsResponse{}
	rsp, err := c.request(ctx, &res.Response, "GET", "/cli_protocol_params", nil, opts)
	if err != nil {
		return res, err
	}
	return res, ReadAndUnmarshalResponse(rsp, &res.Response, &res.Data)
}
