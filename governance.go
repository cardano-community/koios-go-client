// SPDX-License-Identifier: Apache-2.0
//
// Copyright © 2022 The Cardano Community Authors

package koios

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/shopspring/decimal"
)

var (
	// ErrNoDRepID is returned when a DRep ID is required but was not provided.
	ErrNoDRepID = errors.New("missing drep id")
	// ErrNoProposalID is returned when a governance proposal ID is required
	// but was not provided.
	ErrNoProposalID = errors.New("missing proposal id")
)

type (
	// DRepID is a delegated representative ID in CIP-129 bech32 format.
	DRepID string

	// ProposalID is a governance action ID in CIP-129 bech32 format.
	ProposalID string

	// DRepEpochSummary is the voting power and DRep count of an epoch,
	// returned by `/drep_epoch_summary`.
	DRepEpochSummary struct {
		// EpochNo is the epoch number.
		EpochNo EpochNo `json:"epoch_no"`
		// Amount is the total voting power between all DReps including
		// pre-defined roles for the epoch.
		Amount decimal.Decimal `json:"amount"`
		// DReps is the total number of DReps with vote power for the epoch.
		DReps uint `json:"dreps"`
	}

	// DRepListEntry is an item of the DRep list returned by `/drep_list`.
	DRepListEntry struct {
		// DRepID is the DRep ID in CIP-129 bech32 format.
		DRepID DRepID `json:"drep_id"`
		// Hex is the DRep ID in hex format.
		Hex string `json:"hex"`
		// HasScript tells if the credential is a script hash.
		HasScript bool `json:"has_script"`
		// DRepStatus is the DRep's current registration status.
		DRepStatus string `json:"drep_status"`
	}

	// DRepInfo is detailed information about a DRep, returned by `/drep_info`.
	DRepInfo struct {
		// DRepID is the DRep ID in CIP-129 bech32 format.
		DRepID DRepID `json:"drep_id"`
		// Hex is the DRep ID in hex format.
		Hex string `json:"hex"`
		// HasScript tells if the credential is a script hash.
		HasScript bool `json:"has_script"`
		// DRepStatus is the DRep's current registration status.
		DRepStatus string `json:"drep_status"`
		// Deposit is the DRep's registration deposit in lovelace.
		Deposit *decimal.Decimal `json:"deposit"`
		// Active tells if the DRep is active (i.e. not expired).
		Active bool `json:"active"`
		// ExpiresEpochNo is the epoch after which the DRep is considered
		// inactive.
		ExpiresEpochNo *EpochNo `json:"expires_epoch_no"`
		// Amount is the total voting power delegated to this DRep.
		Amount decimal.Decimal `json:"amount"`
		// MetaURL is the URL of the DRep's metadata JSON payload.
		MetaURL *string `json:"meta_url"`
		// MetaHash is the hash of the contents of the metadata URL.
		MetaHash *string `json:"meta_hash"`
		// LiveDelegatorCount is the number of delegators whose last voting
		// power delegation was to this DRep.
		LiveDelegatorCount uint `json:"live_delegator_count"`
	}

	// DRepMetadata is the metadata of a DRep, returned by `/drep_metadata`.
	DRepMetadata struct {
		// DRepID is the DRep ID in CIP-129 bech32 format.
		DRepID DRepID `json:"drep_id"`
		// Hex is the DRep ID in hex format.
		Hex string `json:"hex"`
		// HasScript tells if the credential is a script hash.
		HasScript bool `json:"has_script"`
		// MetaURL is the URL of the metadata JSON payload.
		MetaURL *string `json:"meta_url"`
		// MetaHash is the hash of the contents of the metadata URL.
		MetaHash *string `json:"meta_hash"`
		// MetaJSON is the metadata payload as JSON.
		MetaJSON json.RawMessage `json:"meta_json"`
		// Bytes is the raw bytes (hex) of the metadata payload.
		Bytes *string `json:"bytes"`
		// Warning is a warning raised while validating the metadata.
		Warning *string `json:"warning"`
		// Language is the language of the metadata as per CIP-100.
		Language *string `json:"language"`
		// Comment is the comment attached to the metadata.
		Comment *string `json:"comment"`
		// IsValid tells whether the metadata is valid.
		IsValid *bool `json:"is_valid"`
	}

	// DRepUpdate is a DRep update certificate, returned by `/drep_updates`.
	DRepUpdate struct {
		// DRepID is the DRep ID in CIP-129 bech32 format.
		DRepID DRepID `json:"drep_id"`
		// Hex is the DRep ID in hex format.
		Hex string `json:"hex"`
		// HasScript tells if the credential is a script hash.
		HasScript bool `json:"has_script"`
		// UpdateTxHash is the hash of the transaction carrying the update.
		UpdateTxHash TxHash `json:"update_tx_hash"`
		// CertIndex is the index of the certificate within the transaction.
		CertIndex uint `json:"cert_index"`
		// BlockTime is the time of the block.
		BlockTime Timestamp `json:"block_time"`
		// Action is the effective action of the update certificate.
		Action string `json:"action"`
		// Deposit is the DRep's registration deposit in lovelace.
		Deposit *decimal.Decimal `json:"deposit"`
		// MetaURL is the URL of the metadata JSON payload.
		MetaURL *string `json:"meta_url"`
		// MetaHash is the hash of the contents of the metadata URL.
		MetaHash *string `json:"meta_hash"`
		// MetaJSON is the metadata payload as JSON.
		MetaJSON json.RawMessage `json:"meta_json"`
	}

	// DRepVotingPower is the voting power of a DRep in an epoch, returned by
	// `/drep_history` and `/drep_voting_power_history`.
	DRepVotingPower struct {
		// DRepID is the DRep ID in CIP-129 bech32 format.
		DRepID DRepID `json:"drep_id"`
		// EpochNo is the epoch number.
		EpochNo EpochNo `json:"epoch_no"`
		// Amount is the voting power of the DRep.
		Amount *decimal.Decimal `json:"amount"`
	}

	// PoolVotingPower is the voting power of a pool in an epoch, returned by
	// `/pool_voting_power_history`.
	PoolVotingPower struct {
		// PoolID is the pool ID in bech32 format.
		PoolID PoolID `json:"pool_id_bech32"`
		// EpochNo is the epoch number.
		EpochNo EpochNo `json:"epoch_no"`
		// Amount is the voting power of the pool for the epoch.
		Amount *decimal.Decimal `json:"amount"`
	}

	// GovernanceVote is a vote cast by a single voter on a governance
	// proposal, returned by `/drep_votes`, `/committee_votes` and
	// `/pool_votes`.
	GovernanceVote struct {
		// ProposalID is the proposal action ID in CIP-129 format.
		ProposalID ProposalID `json:"proposal_id"`
		// ProposalTxHash is the hash of the transaction that created the
		// proposal.
		ProposalTxHash TxHash `json:"proposal_tx_hash"`
		// ProposalIndex is the index of the proposal in its transaction.
		ProposalIndex uint `json:"proposal_index"`
		// VoteTxHash is the hash of the transaction carrying the vote.
		VoteTxHash TxHash `json:"vote_tx_hash"`
		// BlockTime is the time of the block.
		BlockTime Timestamp `json:"block_time"`
		// Vote is the vote cast.
		Vote string `json:"vote"`
		// MetaURL is the URL of the vote's metadata JSON payload.
		MetaURL *string `json:"meta_url"`
		// MetaHash is the hash of the contents of the metadata URL.
		MetaHash *string `json:"meta_hash"`
	}

	// DRepDelegator is a delegator of a DRep, returned by `/drep_delegators`.
	DRepDelegator struct {
		// StakeAddress is the stake address (reward account) in bech32 format.
		StakeAddress Address `json:"stake_address"`
		// StakeAddressHex is the stake address in hex format.
		StakeAddressHex string `json:"stake_address_hex"`
		// ScriptHash is set when the stake address is locked by a script.
		ScriptHash *string `json:"script_hash"`
		// EpochNo is the epoch in which the vote delegation was made.
		EpochNo *EpochNo `json:"epoch_no"`
		// Amount is the total balance of the account including UTxO, rewards
		// and MIRs, in lovelace.
		Amount decimal.Decimal `json:"amount"`
	}

	// CommitteeMember is a member of the constitutional committee.
	CommitteeMember struct {
		// Status is the member authentication status.
		Status string `json:"status"`
		// CCColdHex is the cold key hash of the member in hex format.
		CCColdHex string `json:"cc_cold_hex"`
		// CCColdHasScript tells if the cold credential is a script hash.
		CCColdHasScript bool `json:"cc_cold_has_script"`
		// CCHotHex is the key hash from the member's last valid hot key
		// authorization certificate in hex format.
		CCHotHex *string `json:"cc_hot_hex"`
		// CCHotHasScript tells if the hot credential is a script hash.
		CCHotHasScript *bool `json:"cc_hot_has_script"`
		// ExpirationEpoch is the epoch in which the member's vote rights
		// expire.
		ExpirationEpoch EpochNo `json:"expiration_epoch"`
	}

	// CommitteeInfo is information about the active committee, returned by
	// `/committee_info`.
	CommitteeInfo struct {
		// ProposalID is the proposal action ID in CIP-129 format.
		ProposalID ProposalID `json:"proposal_id"`
		// ProposalTxHash is the hash of the transaction that created the
		// proposal.
		ProposalTxHash TxHash `json:"proposal_tx_hash"`
		// ProposalIndex is the index of the proposal in its transaction.
		ProposalIndex uint `json:"proposal_index"`
		// QuorumNumerator is the numerator of the committee quorum threshold.
		QuorumNumerator uint `json:"quorum_numerator"`
		// QuorumDenominator is the denominator of the committee quorum
		// threshold.
		QuorumDenominator uint `json:"quorum_denominator"`
		// Members of the committee.
		Members []CommitteeMember `json:"members"`
	}

	// ProposalWithdrawal is a treasury withdrawal of a governance proposal.
	ProposalWithdrawal struct {
		// StakeAddress is the stake address (reward account) in bech32 format.
		StakeAddress Address `json:"stake_address"`
		// Amount is the withdrawn amount in lovelace.
		Amount decimal.Decimal `json:"amount"`
	}

	// GovernanceProposal is a governance proposal, returned by
	// `/proposal_list` and `/voter_proposal_list`.
	GovernanceProposal struct {
		// BlockTime is the time of the block.
		BlockTime Timestamp `json:"block_time"`
		// ProposalID is the proposal action ID in CIP-129 format.
		ProposalID ProposalID `json:"proposal_id"`
		// ProposalTxHash is the hash of the transaction that created the
		// proposal.
		ProposalTxHash TxHash `json:"proposal_tx_hash"`
		// ProposalIndex is the index of the proposal in its transaction.
		ProposalIndex uint `json:"proposal_index"`
		// ProposalType is the proposal action type.
		ProposalType string `json:"proposal_type"`
		// ProposalDescription is the description of the proposal action.
		ProposalDescription json.RawMessage `json:"proposal_description"`
		// PreviousGovActionProposalID is the CIP-129 formatted ID of the
		// previous governance action in a chain of actions, if any.
		PreviousGovActionProposalID *ProposalID `json:"previous_gov_action_proposal_id"`
		// Deposit is the proposal deposit in lovelace.
		Deposit *decimal.Decimal `json:"deposit"`
		// ReturnAddress is the reward address that receives the deposit when
		// it is repaid.
		ReturnAddress Address `json:"return_address"`
		// ProposedEpoch is the epoch at which the action was proposed.
		ProposedEpoch EpochNo `json:"proposed_epoch"`
		// RatifiedEpoch is the epoch at which the proposal was ratified, if any.
		RatifiedEpoch *EpochNo `json:"ratified_epoch"`
		// EnactedEpoch is the epoch at which the proposal was enacted, if any.
		EnactedEpoch *EpochNo `json:"enacted_epoch"`
		// DroppedEpoch is the epoch at which the proposal was dropped
		// (expired or enacted), if any.
		DroppedEpoch *EpochNo `json:"dropped_epoch"`
		// ExpiredEpoch is the epoch at which the proposal expired, if any.
		ExpiredEpoch *EpochNo `json:"expired_epoch"`
		// Expiration is the epoch at which the action is expected to expire.
		Expiration *EpochNo `json:"expiration"`
		// MetaURL is the URL of the metadata JSON payload.
		MetaURL *string `json:"meta_url"`
		// MetaHash is the hash of the contents of the metadata URL.
		MetaHash *string `json:"meta_hash"`
		// MetaJSON is the metadata payload as JSON.
		MetaJSON json.RawMessage `json:"meta_json"`
		// MetaComment is the comment attached to the metadata.
		MetaComment *string `json:"meta_comment"`
		// MetaLanguage is the language of the metadata as per CIP-100.
		MetaLanguage *string `json:"meta_language"`
		// MetaIsValid tells whether the metadata is valid.
		MetaIsValid *bool `json:"meta_is_valid"`
		// Withdrawal lists the treasury withdrawals of the proposal.
		Withdrawal []ProposalWithdrawal `json:"withdrawal"`
		// ParamProposal is the proposed new parameter set, if any.
		ParamProposal json.RawMessage `json:"param_proposal"`
	}

	// ProposalVotingSummary is the summary of votes for a proposal, returned
	// by `/proposal_voting_summary`.
	ProposalVotingSummary struct {
		// ProposalType is the proposal action type.
		ProposalType string `json:"proposal_type"`
		// EpochNo is the epoch for which the data was collated.
		EpochNo EpochNo `json:"epoch_no"`

		// DRepYesVotesCast is the number of 'yes' votes cast by DReps.
		DRepYesVotesCast uint `json:"drep_yes_votes_cast"`
		// DRepActiveYesVotePower is the power of explicitly cast 'yes' votes.
		DRepActiveYesVotePower decimal.Decimal `json:"drep_active_yes_vote_power"`
		// DRepYesVotePower is the power of 'yes' votes from DReps, including
		// explicit and inferred votes.
		DRepYesVotePower decimal.Decimal `json:"drep_yes_vote_power"`
		// DRepYesPct is the percentage of 'yes' votes from DReps.
		DRepYesPct float64 `json:"drep_yes_pct"`
		// DRepNoVotesCast is the number of 'no' votes cast by DReps.
		DRepNoVotesCast uint `json:"drep_no_votes_cast"`
		// DRepActiveNoVotePower is the power of explicitly cast 'no' votes.
		DRepActiveNoVotePower decimal.Decimal `json:"drep_active_no_vote_power"`
		// DRepNoVotePower is the power of 'no' votes from DReps, including
		// explicit and inferred votes.
		DRepNoVotePower decimal.Decimal `json:"drep_no_vote_power"`
		// DRepNoPct is the percentage of 'no' votes from DReps.
		DRepNoPct float64 `json:"drep_no_pct"`
		// DRepAbstainVotesCast is the number of active 'abstain' votes from
		// DReps.
		DRepAbstainVotesCast uint `json:"drep_abstain_votes_cast"`
		// DRepActiveAbstainVotePower is the power of explicitly cast
		// 'abstain' votes.
		DRepActiveAbstainVotePower decimal.Decimal `json:"drep_active_abstain_vote_power"`
		// DRepAlwaysNoConfidenceVotePower is the power delegated to the
		// predefined 'always_no_confidence' DRep.
		DRepAlwaysNoConfidenceVotePower decimal.Decimal `json:"drep_always_no_confidence_vote_power"`
		// DRepAlwaysAbstainVotePower is the power delegated to the
		// predefined 'always_abstain' DRep.
		DRepAlwaysAbstainVotePower decimal.Decimal `json:"drep_always_abstain_vote_power"`

		// PoolYesVotesCast is the number of 'yes' votes cast by pools.
		PoolYesVotesCast uint `json:"pool_yes_votes_cast"`
		// PoolActiveYesVotePower is the power of explicitly cast 'yes' pool
		// votes.
		PoolActiveYesVotePower decimal.Decimal `json:"pool_active_yes_vote_power"`
		// PoolYesVotePower is the power of 'yes' votes from pools, including
		// explicit and inferred votes.
		PoolYesVotePower decimal.Decimal `json:"pool_yes_vote_power"`
		// PoolYesPct is the percentage of 'yes' votes from pools.
		PoolYesPct float64 `json:"pool_yes_pct"`
		// PoolNoVotesCast is the number of 'no' votes cast by pools.
		PoolNoVotesCast uint `json:"pool_no_votes_cast"`
		// PoolActiveNoVotePower is the power of explicitly cast 'no' pool
		// votes.
		PoolActiveNoVotePower decimal.Decimal `json:"pool_active_no_vote_power"`
		// PoolNoVotePower is the power of 'no' votes from pools, including
		// explicit and inferred votes.
		PoolNoVotePower decimal.Decimal `json:"pool_no_vote_power"`
		// PoolNoPct is the percentage of 'no' votes from pools.
		PoolNoPct float64 `json:"pool_no_pct"`
		// PoolAbstainVotesCast is the number of 'abstain' votes cast by pools.
		PoolAbstainVotesCast uint `json:"pool_abstain_votes_cast"`
		// PoolActiveAbstainVotePower is the power of explicitly cast
		// 'abstain' pool votes.
		PoolActiveAbstainVotePower decimal.Decimal `json:"pool_active_abstain_vote_power"`
		// PoolPassiveAlwaysAbstainVotesAssigned is the number of non-voting
		// pool reward addresses delegating to 'always_abstain'.
		PoolPassiveAlwaysAbstainVotesAssigned uint `json:"pool_passive_always_abstain_votes_assigned"`
		// PoolPassiveAlwaysAbstainVotePower is the combined power of
		// non-voting pools whose reward addresses delegate to
		// 'always_abstain'.
		PoolPassiveAlwaysAbstainVotePower decimal.Decimal `json:"pool_passive_always_abstain_vote_power"`
		// PoolPassiveAlwaysNoConfidenceVotesAssigned is the number of
		// non-voting pool reward addresses delegating to
		// 'always_no_confidence'.
		PoolPassiveAlwaysNoConfidenceVotesAssigned uint `json:"pool_passive_always_no_confidence_votes_assigned"`
		// PoolPassiveAlwaysNoConfidenceVotePower is the combined power of
		// non-voting pools whose reward addresses delegate to
		// 'always_no_confidence'.
		PoolPassiveAlwaysNoConfidenceVotePower decimal.Decimal `json:"pool_passive_always_no_confidence_vote_power"`

		// CommitteeYesVotesCast is the number of 'yes' votes cast by the
		// committee.
		CommitteeYesVotesCast uint `json:"committee_yes_votes_cast"`
		// CommitteeYesPct is the percentage of 'yes' votes from the committee.
		CommitteeYesPct float64 `json:"committee_yes_pct"`
		// CommitteeNoVotesCast is the number of 'no' votes cast by the
		// committee.
		CommitteeNoVotesCast uint `json:"committee_no_votes_cast"`
		// CommitteeNoPct is the percentage of 'no' votes from the committee.
		CommitteeNoPct float64 `json:"committee_no_pct"`
		// CommitteeAbstainVotesCast is the number of 'abstain' votes cast by
		// the committee.
		CommitteeAbstainVotesCast uint `json:"committee_abstain_votes_cast"`
	}

	// ProposalVote is a vote cast on a governance proposal, returned by
	// `/proposal_votes`.
	ProposalVote struct {
		// BlockTime is the time of the block.
		BlockTime Timestamp `json:"block_time"`
		// VoterRole is the role of the voter.
		VoterRole string `json:"voter_role"`
		// VoterID is the voter's DRep ID (CIP-129 bech32), pool ID (bech32)
		// or committee hot ID (CIP-129 bech32).
		VoterID string `json:"voter_id"`
		// VoterHex is the voter's DRep ID, pool ID or committee hash in hex
		// format.
		VoterHex string `json:"voter_hex"`
		// VoterHasScript tells if the voter credential is a script hash.
		VoterHasScript bool `json:"voter_has_script"`
		// Vote is the vote cast.
		Vote string `json:"vote"`
		// MetaURL is the URL of the vote's metadata JSON payload.
		MetaURL *string `json:"meta_url"`
		// MetaHash is the hash of the contents of the metadata URL.
		MetaHash *string `json:"meta_hash"`
	}

	// VoteListEntry is an on-chain vote, returned by `/vote_list`.
	VoteListEntry struct {
		// VoteTxHash is the hash of the transaction in which the vote was
		// posted.
		VoteTxHash TxHash `json:"vote_tx_hash"`
		// VoterRole is the role of the voter.
		VoterRole string `json:"voter_role"`
		// VoterID is the voter's DRep ID (CIP-129 bech32), pool ID (bech32)
		// or committee hot ID (CIP-129 bech32).
		VoterID string `json:"voter_id"`
		// ProposalID is the proposal action ID in CIP-129 format.
		ProposalID ProposalID `json:"proposal_id"`
		// ProposalTxHash is the hash of the transaction that created the
		// proposal.
		ProposalTxHash TxHash `json:"proposal_tx_hash"`
		// ProposalIndex is the index of the proposal in its transaction.
		ProposalIndex uint `json:"proposal_index"`
		// ProposalType is the proposal action type.
		ProposalType string `json:"proposal_type"`
		// EpochNo is the epoch number of the block.
		EpochNo EpochNo `json:"epoch_no"`
		// BlockHeight is the block height.
		BlockHeight BlockNo `json:"block_height"`
		// BlockTime is the time of the block.
		BlockTime Timestamp `json:"block_time"`
		// Vote is the vote cast.
		Vote string `json:"vote"`
		// MetaURL is the URL of the vote's metadata JSON payload.
		MetaURL *string `json:"meta_url"`
		// MetaHash is the hash of the contents of the metadata URL.
		MetaHash *string `json:"meta_hash"`
		// MetaJSON is the metadata payload as JSON.
		MetaJSON json.RawMessage `json:"meta_json"`
	}

	// DRepEpochSummaryResponse represents response from `/drep_epoch_summary` endpoint.
	DRepEpochSummaryResponse struct {
		Response
		Data []DRepEpochSummary `json:"data"`
	}
	// DRepListResponse represents response from `/drep_list` endpoint.
	DRepListResponse struct {
		Response
		Data []DRepListEntry `json:"data"`
	}
	// DRepInfoResponse represents response from `/drep_info` endpoint
	// for a single DRep.
	DRepInfoResponse struct {
		Response
		Data DRepInfo `json:"data"`
	}
	// DRepsInfoResponse represents response from `/drep_info` endpoint.
	DRepsInfoResponse struct {
		Response
		Data []DRepInfo `json:"data"`
	}
	// DRepMetadataResponse represents response from `/drep_metadata` endpoint
	// for a single DRep.
	DRepMetadataResponse struct {
		Response
		Data DRepMetadata `json:"data"`
	}
	// DRepsMetadataResponse represents response from `/drep_metadata` endpoint.
	DRepsMetadataResponse struct {
		Response
		Data []DRepMetadata `json:"data"`
	}
	// DRepUpdatesResponse represents response from `/drep_updates` endpoint.
	DRepUpdatesResponse struct {
		Response
		Data []DRepUpdate `json:"data"`
	}
	// DRepVotingPowerResponse represents response from `/drep_history` and
	// `/drep_voting_power_history` endpoints.
	DRepVotingPowerResponse struct {
		Response
		Data []DRepVotingPower `json:"data"`
	}
	// DRepVotesResponse represents response from `/drep_votes` endpoint.
	DRepVotesResponse struct {
		Response
		Data []GovernanceVote `json:"data"`
	}
	// DRepDelegatorsResponse represents response from `/drep_delegators` endpoint.
	DRepDelegatorsResponse struct {
		Response
		Data []DRepDelegator `json:"data"`
	}
	// CommitteeInfoResponse represents response from `/committee_info` endpoint.
	CommitteeInfoResponse struct {
		Response
		Data CommitteeInfo `json:"data"`
	}
	// CommitteeVotesResponse represents response from `/committee_votes` endpoint.
	CommitteeVotesResponse struct {
		Response
		Data []GovernanceVote `json:"data"`
	}
	// ProposalListResponse represents response from `/proposal_list` and
	// `/voter_proposal_list` endpoints.
	ProposalListResponse struct {
		Response
		Data []GovernanceProposal `json:"data"`
	}
	// ProposalVotingSummaryResponse represents response from
	// `/proposal_voting_summary` endpoint.
	ProposalVotingSummaryResponse struct {
		Response
		Data []ProposalVotingSummary `json:"data"`
	}
	// ProposalVotesResponse represents response from `/proposal_votes` endpoint.
	ProposalVotesResponse struct {
		Response
		Data []ProposalVote `json:"data"`
	}
	// VoteListResponse represents response from `/vote_list` endpoint.
	VoteListResponse struct {
		Response
		Data []VoteListEntry `json:"data"`
	}
	// PoolVotingPowerResponse represents response from
	// `/pool_voting_power_history` endpoint.
	PoolVotingPowerResponse struct {
		Response
		Data []PoolVotingPower `json:"data"`
	}
	// PoolVotesResponse represents response from `/pool_votes` endpoint.
	PoolVotesResponse struct {
		Response
		Data []GovernanceVote `json:"data"`
	}
)

// GetDRepEpochSummary returns the voting power and DRep count of the
// requested epoch, or of every epoch when epoch is nil.
func (c *Client) GetDRepEpochSummary(
	ctx context.Context,
	epoch *EpochNo,
	opts *RequestOptions,
) (res *DRepEpochSummaryResponse, err error) {
	res = &DRepEpochSummaryResponse{}
	opts = c.governanceQuery(opts, "_epoch_no", epochQuery(epoch))
	err = c.governanceRequest(ctx, &res.Response, "GET", "/drep_epoch_summary", nil, opts, &res.Data)
	return
}

// GetDRepList returns the list of all active delegated representatives (DReps).
func (c *Client) GetDRepList(
	ctx context.Context,
	opts *RequestOptions,
) (res *DRepListResponse, err error) {
	res = &DRepListResponse{}
	err = c.governanceRequest(ctx, &res.Response, "GET", "/drep_list", nil, opts, &res.Data)
	return
}

// GetDRepInfo returns detailed information about a delegated representative.
func (c *Client) GetDRepInfo(
	ctx context.Context,
	id DRepID,
	opts *RequestOptions,
) (res *DRepInfoResponse, err error) {
	res = &DRepInfoResponse{}
	if id == "" {
		return res, ErrNoDRepID
	}
	rsp, err := c.GetDRepInfos(ctx, []DRepID{id}, opts)
	res.Response = rsp.Response
	if err != nil {
		return
	}
	if len(rsp.Data) == 0 {
		return res, fmt.Errorf("%w: drep_info response was empty", ErrNoData)
	}
	res.Data = rsp.Data[0]
	return
}

// GetDRepInfos returns detailed information about a set of delegated
// representatives.
func (c *Client) GetDRepInfos(
	ctx context.Context,
	ids []DRepID,
	opts *RequestOptions,
) (res *DRepsInfoResponse, err error) {
	res = &DRepsInfoResponse{}
	if len(ids) == 0 {
		return res, ErrNoDRepID
	}
	err = c.governanceRequest(ctx, &res.Response, "POST", "/drep_info", drepIDsPL(ids), opts, &res.Data)
	return
}

// GetDRepMetadata returns the metadata of a delegated representative.
func (c *Client) GetDRepMetadata(
	ctx context.Context,
	id DRepID,
	opts *RequestOptions,
) (res *DRepMetadataResponse, err error) {
	res = &DRepMetadataResponse{}
	if id == "" {
		return res, ErrNoDRepID
	}
	rsp, err := c.GetDRepMetadatas(ctx, []DRepID{id}, opts)
	res.Response = rsp.Response
	if err != nil {
		return
	}
	if len(rsp.Data) == 0 {
		return res, fmt.Errorf("%w: drep_metadata response was empty", ErrNoData)
	}
	res.Data = rsp.Data[0]
	return
}

// GetDRepMetadatas returns the metadata of a set of delegated representatives.
func (c *Client) GetDRepMetadatas(
	ctx context.Context,
	ids []DRepID,
	opts *RequestOptions,
) (res *DRepsMetadataResponse, err error) {
	res = &DRepsMetadataResponse{}
	if len(ids) == 0 {
		return res, ErrNoDRepID
	}
	err = c.governanceRequest(ctx, &res.Response, "POST", "/drep_metadata", drepIDsPL(ids), opts, &res.Data)
	return
}

// GetDRepUpdates returns the update certificates of the requested delegated
// representative, or of all of them when id is empty.
func (c *Client) GetDRepUpdates(
	ctx context.Context,
	id DRepID,
	opts *RequestOptions,
) (res *DRepUpdatesResponse, err error) {
	res = &DRepUpdatesResponse{}
	opts = c.governanceQuery(opts, "_drep_id", string(id))
	err = c.governanceRequest(ctx, &res.Response, "GET", "/drep_updates", nil, opts, &res.Data)
	return
}

// GetDRepHistory returns the history of DRep voting power for each epoch, or
// only for the requested epoch and DRep when epoch is not nil or id is not
// empty.
func (c *Client) GetDRepHistory(
	ctx context.Context,
	epoch *EpochNo,
	id DRepID,
	opts *RequestOptions,
) (res *DRepVotingPowerResponse, err error) {
	res = &DRepVotingPowerResponse{}
	opts = c.governanceQuery(opts, "_epoch_no", epochQuery(epoch), "_drep_id", string(id))
	err = c.governanceRequest(ctx, &res.Response, "GET", "/drep_history", nil, opts, &res.Data)
	return
}

// GetDRepVotingPowerHistory returns the history of DRep voting power for each
// epoch, or only for the requested epoch and DRep when epoch is not nil or id
// is not empty.
func (c *Client) GetDRepVotingPowerHistory(
	ctx context.Context,
	epoch *EpochNo,
	id DRepID,
	opts *RequestOptions,
) (res *DRepVotingPowerResponse, err error) {
	res = &DRepVotingPowerResponse{}
	opts = c.governanceQuery(opts, "_epoch_no", epochQuery(epoch), "_drep_id", string(id))
	err = c.governanceRequest(ctx, &res.Response, "GET", "/drep_voting_power_history", nil, opts, &res.Data)
	return
}

// GetDRepVotes returns all votes cast by the requested delegated representative.
func (c *Client) GetDRepVotes(
	ctx context.Context,
	id DRepID,
	opts *RequestOptions,
) (res *DRepVotesResponse, err error) {
	res = &DRepVotesResponse{}
	if id == "" {
		return res, ErrNoDRepID
	}
	opts = c.governanceQuery(opts, "_drep_id", string(id))
	err = c.governanceRequest(ctx, &res.Response, "GET", "/drep_votes", nil, opts, &res.Data)
	return
}

// GetDRepDelegators returns all delegators of the requested delegated
// representative.
func (c *Client) GetDRepDelegators(
	ctx context.Context,
	id DRepID,
	opts *RequestOptions,
) (res *DRepDelegatorsResponse, err error) {
	res = &DRepDelegatorsResponse{}
	if id == "" {
		return res, ErrNoDRepID
	}
	opts = c.governanceQuery(opts, "_drep_id", string(id))
	err = c.governanceRequest(ctx, &res.Response, "GET", "/drep_delegators", nil, opts, &res.Data)
	return
}

// GetCommitteeInfo returns information about the active constitutional
// committee and its members.
func (c *Client) GetCommitteeInfo(
	ctx context.Context,
	opts *RequestOptions,
) (res *CommitteeInfoResponse, err error) {
	res = &CommitteeInfoResponse{}
	err = c.governanceRequest(ctx, &res.Response, "GET", "/committee_info", nil, opts, &res.Data)
	if err == nil && len(res.Data.ProposalTxHash) == 0 && len(res.Data.Members) == 0 {
		err = fmt.Errorf("%w: committee_info response was empty", ErrNoData)
		res.applyError(nil, err)
	}
	return
}

// GetCommitteeVotes returns all votes cast by the given committee member
// (hot key ID in bech32, CIP-5 or CIP-129), or by the whole committee when
// ccHotID is empty.
func (c *Client) GetCommitteeVotes(
	ctx context.Context,
	ccHotID string,
	opts *RequestOptions,
) (res *CommitteeVotesResponse, err error) {
	res = &CommitteeVotesResponse{}
	opts = c.governanceQuery(opts, "_cc_hot_id", ccHotID)
	err = c.governanceRequest(ctx, &res.Response, "GET", "/committee_votes", nil, opts, &res.Data)
	return
}

// UnmarshalJSON decodes the committee information from either a JSON object,
// as documented by the API specification, or a single-element array, which
// some Koios deployments return for `/committee_info`.
func (ci *CommitteeInfo) UnmarshalJSON(b []byte) error {
	type plain CommitteeInfo
	if trimmed := bytes.TrimSpace(b); len(trimmed) > 0 && trimmed[0] == '[' {
		var list []plain
		if err := json.Unmarshal(trimmed, &list); err != nil {
			return err
		}
		if len(list) == 0 {
			*ci = CommitteeInfo{}
			return nil
		}
		*ci = CommitteeInfo(list[0])
		return nil
	}
	return json.Unmarshal(b, (*plain)(ci))
}

// GetProposalList returns all governance proposals.
func (c *Client) GetProposalList(
	ctx context.Context,
	opts *RequestOptions,
) (res *ProposalListResponse, err error) {
	res = &ProposalListResponse{}
	err = c.governanceRequest(ctx, &res.Response, "GET", "/proposal_list", nil, opts, &res.Data)
	return
}

// GetVoterProposalList returns the governance proposals for the specified
// DRep, SPO or committee member credential (bech32, CIP-5 or CIP-129), or all
// proposals when voterID is empty.
func (c *Client) GetVoterProposalList(
	ctx context.Context,
	voterID string,
	opts *RequestOptions,
) (res *ProposalListResponse, err error) {
	res = &ProposalListResponse{}
	opts = c.governanceQuery(opts, "_voter_id", voterID)
	err = c.governanceRequest(ctx, &res.Response, "GET", "/voter_proposal_list", nil, opts, &res.Data)
	return
}

// GetProposalVotingSummary returns the summary of votes for the given proposal.
func (c *Client) GetProposalVotingSummary(
	ctx context.Context,
	id ProposalID,
	opts *RequestOptions,
) (res *ProposalVotingSummaryResponse, err error) {
	res = &ProposalVotingSummaryResponse{}
	if id == "" {
		return res, ErrNoProposalID
	}
	opts = c.governanceQuery(opts, "_proposal_id", string(id))
	err = c.governanceRequest(ctx, &res.Response, "GET", "/proposal_voting_summary", nil, opts, &res.Data)
	return
}

// GetProposalVotes returns all votes cast on the given governance proposal.
func (c *Client) GetProposalVotes(
	ctx context.Context,
	id ProposalID,
	opts *RequestOptions,
) (res *ProposalVotesResponse, err error) {
	res = &ProposalVotesResponse{}
	if id == "" {
		return res, ErrNoProposalID
	}
	opts = c.governanceQuery(opts, "_proposal_id", string(id))
	err = c.governanceRequest(ctx, &res.Response, "GET", "/proposal_votes", nil, opts, &res.Data)
	return
}

// GetVoteList returns all votes posted on-chain.
func (c *Client) GetVoteList(
	ctx context.Context,
	opts *RequestOptions,
) (res *VoteListResponse, err error) {
	res = &VoteListResponse{}
	err = c.governanceRequest(ctx, &res.Response, "GET", "/vote_list", nil, opts, &res.Data)
	return
}

// GetPoolVotingPowerHistory returns the history of pool voting power for each
// epoch, or only for the requested epoch and pool when epoch is not nil or
// pool is not empty.
func (c *Client) GetPoolVotingPowerHistory(
	ctx context.Context,
	epoch *EpochNo,
	pool PoolID,
	opts *RequestOptions,
) (res *PoolVotingPowerResponse, err error) {
	res = &PoolVotingPowerResponse{}
	opts = c.governanceQuery(opts, "_epoch_no", epochQuery(epoch), "_pool_bech32", string(pool))
	err = c.governanceRequest(ctx, &res.Response, "GET", "/pool_voting_power_history", nil, opts, &res.Data)
	return
}

// GetPoolVotes returns all votes cast by the given pool.
func (c *Client) GetPoolVotes(
	ctx context.Context,
	pool PoolID,
	opts *RequestOptions,
) (res *PoolVotesResponse, err error) {
	res = &PoolVotesResponse{}
	if pool == "" {
		return res, ErrNoPoolID
	}
	opts = c.governanceQuery(opts, "_pool_bech32", string(pool))
	err = c.governanceRequest(ctx, &res.Response, "GET", "/pool_votes", nil, opts, &res.Data)
	return
}

// governanceQuery sets the given key/value pairs on opts, creating opts when
// nil and skipping pairs whose value is empty.
func (c *Client) governanceQuery(opts *RequestOptions, kv ...string) *RequestOptions {
	if opts == nil {
		opts = c.NewRequestOptions()
	}
	for i := 0; i+1 < len(kv); i += 2 {
		if kv[i+1] != "" {
			opts.QuerySet(kv[i], kv[i+1])
		}
	}
	return opts
}

func epochQuery(epoch *EpochNo) string {
	if epoch == nil {
		return ""
	}
	return epoch.String()
}

func (c *Client) governanceRequest(
	ctx context.Context,
	res *Response,
	method, path string,
	body io.Reader,
	opts *RequestOptions,
	data any,
) error {
	rsp, err := c.request(ctx, res, method, path, body, opts)
	if err != nil {
		return err
	}
	return ReadAndUnmarshalResponse(rsp, res, data)
}

func drepIDsPL(ids []DRepID) io.Reader {
	b, _ := json.Marshal(struct {
		DRepIDs []DRepID `json:"_drep_ids"`
	}{ids})
	return bytes.NewReader(b)
}
