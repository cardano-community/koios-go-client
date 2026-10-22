// SPDX-License-Identifier: Apache-2.0
//
// Copyright © 2022 The Cardano Community Authors

package koios_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	koios "github.com/cardano-community/koios-go-client/v5"
)

type govCapture struct {
	method string
	path   string
	query  url.Values
	body   string
}

// govClient returns a client talking to a test server that records the
// request and answers with the given JSON body.
func govClient(t *testing.T, response string) (*koios.Client, *govCapture) {
	t.Helper()
	capture := &govCapture{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		capture.method = r.Method
		capture.path = r.URL.Path
		capture.query = r.URL.Query()
		capture.body = strings.TrimSpace(string(b))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(response))
	}))
	t.Cleanup(srv.Close)

	u, err := url.Parse(srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	port, err := strconv.ParseUint(u.Port(), 10, 16)
	if err != nil {
		t.Fatal(err)
	}
	client, err := koios.New(
		koios.Scheme("http"),
		koios.Host(u.Hostname()),
		koios.Port(uint16(port)),
		koios.RateLimit(255),
	)
	if err != nil {
		t.Fatal(err)
	}
	return client, capture
}

func govExpect(t *testing.T, c *govCapture, method, path string, query map[string]string, body string) {
	t.Helper()
	if c.method != method {
		t.Errorf("method = %q, want %q", c.method, method)
	}
	if !strings.HasSuffix(c.path, path) {
		t.Errorf("path = %q, want suffix %q", c.path, path)
	}
	if len(c.query) != len(query) {
		t.Errorf("query = %v, want %v", c.query, query)
	}
	for k, v := range query {
		if got := c.query.Get(k); got != v {
			t.Errorf("query %s = %q, want %q", k, got, v)
		}
	}
	if c.body != body {
		t.Errorf("body = %q, want %q", c.body, body)
	}
}

func govEq[T comparable](t *testing.T, name string, got, want T) {
	t.Helper()
	if got != want {
		t.Errorf("%s = %v, want %v", name, got, want)
	}
}

func govNoErr(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

const (
	govDRep     = "drep1ytc5jz9mu2crjrqugvqsvqe8jdu9cv5hlzgf3lxjpsd2dsq5hcpu6"
	govProposal = "gov_action1k2jertppnnndejjcglszfqq4yzw8evzrd2nj6g4kpxuupmugz0lqq9kwsdv"
	govPool     = "pool1pu5jlj4q9w9jlxeu370a3c9myx47md5j5m2str0naunn2q3lkdy"
	govTx       = "6ecf8b5f6ea47c8a5ab06a4fbaa5e4d1d9a7b7e0c4a1a0f7a8e1a2d0e8c5b9a4"
)

func TestGovernanceDRepEpochSummary(t *testing.T) {
	client, cap := govClient(t, `[{"epoch_no":530,"amount":"1066331327890419","dreps":4256}]`)
	epoch := koios.EpochNo(530)
	res, err := client.GetDRepEpochSummary(context.Background(), &epoch, nil)
	govNoErr(t, err)
	govExpect(t, cap, "GET", "/drep_epoch_summary", map[string]string{"_epoch_no": "530"}, "")
	govEq(t, "len", len(res.Data), 1)
	govEq(t, "epoch", res.Data[0].EpochNo, koios.EpochNo(530))
	govEq(t, "amount", res.Data[0].Amount.String(), "1066331327890419")
	govEq(t, "dreps", res.Data[0].DReps, uint(4256))

	_, err = client.GetDRepEpochSummary(context.Background(), nil, nil)
	govNoErr(t, err)
	govExpect(t, cap, "GET", "/drep_epoch_summary", nil, "")
}

func TestGovernanceDRepList(t *testing.T) {
	client, cap := govClient(t, `[{"drep_id":"`+govDRep+`","hex":"2f1490bbe2b0390c1c4430104c1e49b78b8651","has_script":false,"drep_status":"registered"}]`)
	res, err := client.GetDRepList(context.Background(), nil)
	govNoErr(t, err)
	govExpect(t, cap, "GET", "/drep_list", nil, "")
	govEq(t, "drep_id", res.Data[0].DRepID, koios.DRepID(govDRep))
	govEq(t, "has_script", res.Data[0].HasScript, false)
	govEq(t, "status", res.Data[0].DRepStatus, "registered")
}

const govDRepInfoJSON = `[{"drep_id":"` + govDRep + `","hex":"2f14","has_script":false,"drep_status":"registered","deposit":"500000000","active":true,"expires_epoch_no":410,"amount":"8696873","meta_url":"https://example.org/drep.jsonld","meta_hash":"ab12","live_delegator_count":7}]`

func TestGovernanceDRepInfo(t *testing.T) {
	client, cap := govClient(t, govDRepInfoJSON)
	res, err := client.GetDRepInfo(context.Background(), koios.DRepID(govDRep), nil)
	govNoErr(t, err)
	govExpect(t, cap, "POST", "/drep_info", nil, `{"_drep_ids":["`+govDRep+`"]}`)
	d := res.Data
	govEq(t, "deposit", d.Deposit.String(), "500000000")
	govEq(t, "active", d.Active, true)
	govEq(t, "expires", *d.ExpiresEpochNo, koios.EpochNo(410))
	govEq(t, "amount", d.Amount.String(), "8696873")
	govEq(t, "meta_url", *d.MetaURL, "https://example.org/drep.jsonld")
	govEq(t, "meta_hash", *d.MetaHash, "ab12")
	govEq(t, "delegators", d.LiveDelegatorCount, uint(7))

	_, err = client.GetDRepInfo(context.Background(), "", nil)
	if !errors.Is(err, koios.ErrNoDRepID) {
		t.Errorf("err = %v, want ErrNoDRepID", err)
	}
}

func TestGovernanceDRepInfosNullable(t *testing.T) {
	client, cap := govClient(t, `[{"drep_id":"a","hex":"b","has_script":true,"drep_status":"retired","deposit":null,"active":false,"expires_epoch_no":null,"amount":"0","meta_url":null,"meta_hash":null,"live_delegator_count":0}]`)
	res, err := client.GetDRepInfos(context.Background(), []koios.DRepID{"a", "c"}, nil)
	govNoErr(t, err)
	govExpect(t, cap, "POST", "/drep_info", nil, `{"_drep_ids":["a","c"]}`)
	d := res.Data[0]
	if d.Deposit != nil || d.ExpiresEpochNo != nil || d.MetaURL != nil || d.MetaHash != nil {
		t.Errorf("expected nullable fields to be nil: %+v", d)
	}
	_, err = client.GetDRepInfos(context.Background(), nil, nil)
	if !errors.Is(err, koios.ErrNoDRepID) {
		t.Errorf("err = %v, want ErrNoDRepID", err)
	}
}

func TestGovernanceDRepMetadata(t *testing.T) {
	client, cap := govClient(t, `[{"drep_id":"`+govDRep+`","hex":"2f14","has_script":false,"meta_url":"https://example.org/m.jsonld","meta_hash":"ab12","meta_json":{"body":{"title":"x"}},"bytes":"7b7d","warning":null,"language":"en-us","comment":null,"is_valid":true}]`)
	res, err := client.GetDRepMetadata(context.Background(), koios.DRepID(govDRep), nil)
	govNoErr(t, err)
	govExpect(t, cap, "POST", "/drep_metadata", nil, `{"_drep_ids":["`+govDRep+`"]}`)
	d := res.Data
	govEq(t, "meta_json", string(d.MetaJSON), `{"body":{"title":"x"}}`)
	govEq(t, "bytes", *d.Bytes, "7b7d")
	govEq(t, "language", *d.Language, "en-us")
	govEq(t, "is_valid", *d.IsValid, true)
	if d.Warning != nil || d.Comment != nil {
		t.Errorf("warning/comment should be nil")
	}

	multi, err := client.GetDRepMetadatas(context.Background(), []koios.DRepID{"a", "b"}, nil)
	govNoErr(t, err)
	govExpect(t, cap, "POST", "/drep_metadata", nil, `{"_drep_ids":["a","b"]}`)
	govEq(t, "len", len(multi.Data), 1)
}

func TestGovernanceDRepUpdates(t *testing.T) {
	client, cap := govClient(t, `[{"drep_id":"`+govDRep+`","hex":"2f14","has_script":false,"update_tx_hash":"`+govTx+`","cert_index":0,"block_time":1506203091,"action":"updated","deposit":null,"meta_url":"https://example.org/m","meta_hash":"ab","meta_json":null}]`)
	res, err := client.GetDRepUpdates(context.Background(), koios.DRepID(govDRep), nil)
	govNoErr(t, err)
	govExpect(t, cap, "GET", "/drep_updates", map[string]string{"_drep_id": govDRep}, "")
	u := res.Data[0]
	govEq(t, "tx", u.UpdateTxHash, koios.TxHash(govTx))
	govEq(t, "action", u.Action, "updated")
	govEq(t, "block_time", u.BlockTime.Unix(), int64(1506203091))
	if u.Deposit != nil {
		t.Errorf("deposit should be nil")
	}

	_, err = client.GetDRepUpdates(context.Background(), "", nil)
	govNoErr(t, err)
	govExpect(t, cap, "GET", "/drep_updates", nil, "")
}

func TestGovernanceDRepHistories(t *testing.T) {
	const body = `[{"drep_id":"` + govDRep + `","epoch_no":500,"amount":"14231445553"},{"drep_id":"x","epoch_no":501,"amount":null}]`
	epoch := koios.EpochNo(500)
	for path, call := range map[string]func(*koios.Client) (*koios.DRepVotingPowerResponse, error){
		"/drep_history": func(c *koios.Client) (*koios.DRepVotingPowerResponse, error) {
			return c.GetDRepHistory(context.Background(), &epoch, koios.DRepID(govDRep), nil)
		},
		"/drep_voting_power_history": func(c *koios.Client) (*koios.DRepVotingPowerResponse, error) {
			return c.GetDRepVotingPowerHistory(context.Background(), &epoch, koios.DRepID(govDRep), nil)
		},
	} {
		client, cap := govClient(t, body)
		res, err := call(client)
		govNoErr(t, err)
		govExpect(t, cap, "GET", path, map[string]string{"_epoch_no": "500", "_drep_id": govDRep}, "")
		govEq(t, path+" len", len(res.Data), 2)
		govEq(t, path+" amount", res.Data[0].Amount.String(), "14231445553")
		govEq(t, path+" epoch", res.Data[1].EpochNo, koios.EpochNo(501))
		if res.Data[1].Amount != nil {
			t.Errorf("%s: null amount should decode to nil", path)
		}
	}
}

const govVotesJSON = `[{"proposal_id":"` + govProposal + `","proposal_tx_hash":"` + govTx + `","proposal_index":2,"vote_tx_hash":"abcd","block_time":1506203091,"vote":"Yes","meta_url":null,"meta_hash":"ff"}]`

func govCheckVote(t *testing.T, v koios.GovernanceVote) {
	t.Helper()
	govEq(t, "proposal_id", v.ProposalID, koios.ProposalID(govProposal))
	govEq(t, "proposal_tx_hash", v.ProposalTxHash, koios.TxHash(govTx))
	govEq(t, "proposal_index", v.ProposalIndex, uint(2))
	govEq(t, "vote_tx_hash", v.VoteTxHash, koios.TxHash("abcd"))
	govEq(t, "block_time", v.BlockTime.Unix(), int64(1506203091))
	govEq(t, "vote", v.Vote, "Yes")
	if v.MetaURL != nil {
		t.Errorf("meta_url should be nil")
	}
	govEq(t, "meta_hash", *v.MetaHash, "ff")
}

func TestGovernanceDRepVotes(t *testing.T) {
	client, cap := govClient(t, govVotesJSON)
	res, err := client.GetDRepVotes(context.Background(), koios.DRepID(govDRep), nil)
	govNoErr(t, err)
	govExpect(t, cap, "GET", "/drep_votes", map[string]string{"_drep_id": govDRep}, "")
	govCheckVote(t, res.Data[0])
	_, err = client.GetDRepVotes(context.Background(), "", nil)
	if !errors.Is(err, koios.ErrNoDRepID) {
		t.Errorf("err = %v, want ErrNoDRepID", err)
	}
}

func TestGovernanceCommitteeVotes(t *testing.T) {
	client, cap := govClient(t, govVotesJSON)
	res, err := client.GetCommitteeVotes(context.Background(), "cc_hot1abc", nil)
	govNoErr(t, err)
	govExpect(t, cap, "GET", "/committee_votes", map[string]string{"_cc_hot_id": "cc_hot1abc"}, "")
	govCheckVote(t, res.Data[0])
	_, err = client.GetCommitteeVotes(context.Background(), "", nil)
	govNoErr(t, err)
	govExpect(t, cap, "GET", "/committee_votes", nil, "")
}

func TestGovernancePoolVotes(t *testing.T) {
	client, cap := govClient(t, govVotesJSON)
	res, err := client.GetPoolVotes(context.Background(), koios.PoolID(govPool), nil)
	govNoErr(t, err)
	govExpect(t, cap, "GET", "/pool_votes", map[string]string{"_pool_bech32": govPool}, "")
	govCheckVote(t, res.Data[0])
	_, err = client.GetPoolVotes(context.Background(), "", nil)
	if !errors.Is(err, koios.ErrNoPoolID) {
		t.Errorf("err = %v, want ErrNoPoolID", err)
	}
}

func TestGovernanceDRepDelegators(t *testing.T) {
	client, cap := govClient(t, `[{"stake_address":"stake1u8abc","stake_address_hex":"e1aa","script_hash":null,"epoch_no":490,"amount":"1234567"},{"stake_address":"stake1u9def","stake_address_hex":"e1bb","script_hash":"cc","epoch_no":null,"amount":"0"}]`)
	res, err := client.GetDRepDelegators(context.Background(), koios.DRepID(govDRep), nil)
	govNoErr(t, err)
	govExpect(t, cap, "GET", "/drep_delegators", map[string]string{"_drep_id": govDRep}, "")
	govEq(t, "stake", res.Data[0].StakeAddress, koios.Address("stake1u8abc"))
	govEq(t, "hex", res.Data[0].StakeAddressHex, "e1aa")
	govEq(t, "epoch", *res.Data[0].EpochNo, koios.EpochNo(490))
	govEq(t, "amount", res.Data[0].Amount.String(), "1234567")
	govEq(t, "script", *res.Data[1].ScriptHash, "cc")
	if res.Data[0].ScriptHash != nil || res.Data[1].EpochNo != nil {
		t.Errorf("nullable fields should be nil")
	}
	_, err = client.GetDRepDelegators(context.Background(), "", nil)
	if !errors.Is(err, koios.ErrNoDRepID) {
		t.Errorf("err = %v, want ErrNoDRepID", err)
	}
}

const govCommitteeObj = `{"proposal_id":"` + govProposal + `","proposal_tx_hash":"` + govTx + `","proposal_index":1,"quorum_numerator":67,"quorum_denominator":100,"members":[{"status":"authorized","cc_cold_hex":"6095","cc_cold_has_script":false,"cc_hot_hex":"65d4","cc_hot_has_script":false,"expiration_epoch":324},{"status":"not_authorized","cc_cold_hex":"7aaa","cc_cold_has_script":true,"cc_hot_hex":null,"cc_hot_has_script":null,"expiration_epoch":330}]}`

func TestGovernanceCommitteeInfo(t *testing.T) {
	for name, body := range map[string]string{"object": govCommitteeObj, "array": "[" + govCommitteeObj + "]"} {
		client, cap := govClient(t, body)
		res, err := client.GetCommitteeInfo(context.Background(), nil)
		govNoErr(t, err)
		govExpect(t, cap, "GET", "/committee_info", nil, "")
		d := res.Data
		govEq(t, name+" proposal", d.ProposalID, koios.ProposalID(govProposal))
		govEq(t, name+" index", d.ProposalIndex, uint(1))
		govEq(t, name+" num", d.QuorumNumerator, uint(67))
		govEq(t, name+" den", d.QuorumDenominator, uint(100))
		govEq(t, name+" members", len(d.Members), 2)
		govEq(t, name+" status", d.Members[0].Status, "authorized")
		govEq(t, name+" hot", *d.Members[0].CCHotHex, "65d4")
		govEq(t, name+" exp", d.Members[0].ExpirationEpoch, koios.EpochNo(324))
		govEq(t, name+" cold script", d.Members[1].CCColdHasScript, true)
		if d.Members[1].CCHotHex != nil || d.Members[1].CCHotHasScript != nil {
			t.Errorf("%s: unauthorized member hot key should be nil", name)
		}
	}
}

const govProposalJSON = `[{"block_time":1506203091,"proposal_id":"` + govProposal + `","proposal_tx_hash":"` + govTx + `","proposal_index":0,"proposal_type":"TreasuryWithdrawals","proposal_description":{"tag":"TreasuryWithdrawals"},"previous_gov_action_proposal_id":null,"deposit":"100000000000","return_address":"stake1u8abc","proposed_epoch":500,"ratified_epoch":505,"enacted_epoch":null,"dropped_epoch":null,"expired_epoch":null,"expiration":510,"meta_url":"https://example.org/p.json","meta_hash":"ab","meta_json":{"body":{"title":"t"}},"meta_comment":null,"meta_language":"en-us","meta_is_valid":true,"withdrawal":[{"stake_address":"stake1u9def","amount":"5000000"}],"param_proposal":null}]`

func govCheckProposal(t *testing.T, p koios.GovernanceProposal) {
	t.Helper()
	govEq(t, "proposal_type", p.ProposalType, "TreasuryWithdrawals")
	govEq(t, "description", string(p.ProposalDescription), `{"tag":"TreasuryWithdrawals"}`)
	govEq(t, "deposit", p.Deposit.String(), "100000000000")
	govEq(t, "return", p.ReturnAddress, koios.Address("stake1u8abc"))
	govEq(t, "proposed", p.ProposedEpoch, koios.EpochNo(500))
	govEq(t, "ratified", *p.RatifiedEpoch, koios.EpochNo(505))
	govEq(t, "expiration", *p.Expiration, koios.EpochNo(510))
	govEq(t, "meta_json", string(p.MetaJSON), `{"body":{"title":"t"}}`)
	govEq(t, "meta_is_valid", *p.MetaIsValid, true)
	govEq(t, "withdrawals", len(p.Withdrawal), 1)
	govEq(t, "withdrawal amount", p.Withdrawal[0].Amount.String(), "5000000")
	if p.PreviousGovActionProposalID != nil || p.EnactedEpoch != nil || p.DroppedEpoch != nil || p.ExpiredEpoch != nil || p.MetaComment != nil {
		t.Errorf("nullable fields should be nil: %+v", p)
	}
	if s := string(p.ParamProposal); s != "null" && s != "" {
		t.Errorf("param_proposal = %q, want null", s)
	}
}

func TestGovernanceProposalList(t *testing.T) {
	client, cap := govClient(t, govProposalJSON)
	res, err := client.GetProposalList(context.Background(), nil)
	govNoErr(t, err)
	govExpect(t, cap, "GET", "/proposal_list", nil, "")
	govCheckProposal(t, res.Data[0])
}

func TestGovernanceVoterProposalList(t *testing.T) {
	client, cap := govClient(t, govProposalJSON)
	res, err := client.GetVoterProposalList(context.Background(), govDRep, nil)
	govNoErr(t, err)
	govExpect(t, cap, "GET", "/voter_proposal_list", map[string]string{"_voter_id": govDRep}, "")
	govCheckProposal(t, res.Data[0])
	_, err = client.GetVoterProposalList(context.Background(), "", nil)
	govNoErr(t, err)
	govExpect(t, cap, "GET", "/voter_proposal_list", nil, "")
}

func TestGovernanceProposalVotingSummary(t *testing.T) {
	client, cap := govClient(t, `[{"proposal_type":"ParameterChange","epoch_no":520,
"drep_yes_votes_cast":10,"drep_active_yes_vote_power":"100","drep_yes_vote_power":"110","drep_yes_pct":55.5,
"drep_no_votes_cast":3,"drep_active_no_vote_power":"30","drep_no_vote_power":"33","drep_no_pct":16.5,
"drep_abstain_votes_cast":2,"drep_active_abstain_vote_power":"20","drep_always_no_confidence_vote_power":"5","drep_always_abstain_vote_power":"6",
"pool_yes_votes_cast":4,"pool_active_yes_vote_power":"40","pool_yes_vote_power":"44","pool_yes_pct":60.25,
"pool_no_votes_cast":1,"pool_active_no_vote_power":"10","pool_no_vote_power":"11","pool_no_pct":10.5,
"pool_abstain_votes_cast":1,"pool_active_abstain_vote_power":"12","pool_passive_always_abstain_votes_assigned":7,"pool_passive_always_abstain_vote_power":"70",
"pool_passive_always_no_confidence_votes_assigned":8,"pool_passive_always_no_confidence_vote_power":"80",
"committee_yes_votes_cast":5,"committee_yes_pct":71.4,"committee_no_votes_cast":1,"committee_no_pct":14.3,"committee_abstain_votes_cast":1}]`)
	res, err := client.GetProposalVotingSummary(context.Background(), koios.ProposalID(govProposal), nil)
	govNoErr(t, err)
	govExpect(t, cap, "GET", "/proposal_voting_summary", map[string]string{"_proposal_id": govProposal}, "")
	s := res.Data[0]
	govEq(t, "type", s.ProposalType, "ParameterChange")
	govEq(t, "epoch", s.EpochNo, koios.EpochNo(520))
	govEq(t, "drep yes cast", s.DRepYesVotesCast, uint(10))
	govEq(t, "drep active yes", s.DRepActiveYesVotePower.String(), "100")
	govEq(t, "drep yes power", s.DRepYesVotePower.String(), "110")
	govEq(t, "drep yes pct", s.DRepYesPct, 55.5)
	govEq(t, "drep no cast", s.DRepNoVotesCast, uint(3))
	govEq(t, "drep no pct", s.DRepNoPct, 16.5)
	govEq(t, "drep abstain cast", s.DRepAbstainVotesCast, uint(2))
	govEq(t, "drep abstain power", s.DRepActiveAbstainVotePower.String(), "20")
	govEq(t, "drep no confidence", s.DRepAlwaysNoConfidenceVotePower.String(), "5")
	govEq(t, "drep always abstain", s.DRepAlwaysAbstainVotePower.String(), "6")
	govEq(t, "pool yes cast", s.PoolYesVotesCast, uint(4))
	govEq(t, "pool yes pct", s.PoolYesPct, 60.25)
	govEq(t, "pool no power", s.PoolNoVotePower.String(), "11")
	govEq(t, "pool abstain cast", s.PoolAbstainVotesCast, uint(1))
	govEq(t, "pool passive abstain n", s.PoolPassiveAlwaysAbstainVotesAssigned, uint(7))
	govEq(t, "pool passive abstain power", s.PoolPassiveAlwaysAbstainVotePower.String(), "70")
	govEq(t, "pool passive noconf n", s.PoolPassiveAlwaysNoConfidenceVotesAssigned, uint(8))
	govEq(t, "pool passive noconf power", s.PoolPassiveAlwaysNoConfidenceVotePower.String(), "80")
	govEq(t, "cc yes cast", s.CommitteeYesVotesCast, uint(5))
	govEq(t, "cc yes pct", s.CommitteeYesPct, 71.4)
	govEq(t, "cc no pct", s.CommitteeNoPct, 14.3)
	govEq(t, "cc abstain", s.CommitteeAbstainVotesCast, uint(1))

	_, err = client.GetProposalVotingSummary(context.Background(), "", nil)
	if !errors.Is(err, koios.ErrNoProposalID) {
		t.Errorf("err = %v, want ErrNoProposalID", err)
	}
}

func TestGovernanceProposalVotes(t *testing.T) {
	client, cap := govClient(t, `[{"block_time":1506203091,"voter_role":"DRep","voter_id":"`+govDRep+`","voter_hex":"2f14","voter_has_script":false,"vote":"No","meta_url":null,"meta_hash":null}]`)
	res, err := client.GetProposalVotes(context.Background(), koios.ProposalID(govProposal), nil)
	govNoErr(t, err)
	govExpect(t, cap, "GET", "/proposal_votes", map[string]string{"_proposal_id": govProposal}, "")
	v := res.Data[0]
	govEq(t, "role", v.VoterRole, "DRep")
	govEq(t, "voter", v.VoterID, govDRep)
	govEq(t, "hex", v.VoterHex, "2f14")
	govEq(t, "script", v.VoterHasScript, false)
	govEq(t, "vote", v.Vote, "No")
	govEq(t, "block_time", v.BlockTime.Unix(), int64(1506203091))
	_, err = client.GetProposalVotes(context.Background(), "", nil)
	if !errors.Is(err, koios.ErrNoProposalID) {
		t.Errorf("err = %v, want ErrNoProposalID", err)
	}
}

func TestGovernanceVoteList(t *testing.T) {
	client, cap := govClient(t, `[{"vote_tx_hash":"`+govTx+`","voter_role":"SPO","voter_id":"`+govPool+`","proposal_id":"`+govProposal+`","proposal_tx_hash":"abcd","proposal_index":3,"proposal_type":"InfoAction","epoch_no":512,"block_height":10500000,"block_time":1506203091,"vote":"Abstain","meta_url":"https://example.org/v","meta_hash":"ee","meta_json":{"a":1}}]`)
	res, err := client.GetVoteList(context.Background(), nil)
	govNoErr(t, err)
	govExpect(t, cap, "GET", "/vote_list", nil, "")
	v := res.Data[0]
	govEq(t, "tx", v.VoteTxHash, koios.TxHash(govTx))
	govEq(t, "role", v.VoterRole, "SPO")
	govEq(t, "voter", v.VoterID, govPool)
	govEq(t, "index", v.ProposalIndex, uint(3))
	govEq(t, "type", v.ProposalType, "InfoAction")
	govEq(t, "epoch", v.EpochNo, koios.EpochNo(512))
	govEq(t, "height", v.BlockHeight, koios.BlockNo(10500000))
	govEq(t, "vote", v.Vote, "Abstain")
	govEq(t, "meta_url", *v.MetaURL, "https://example.org/v")
	govEq(t, "meta_json", string(v.MetaJSON), `{"a":1}`)
}

func TestGovernancePoolVotingPowerHistory(t *testing.T) {
	client, cap := govClient(t, `[{"pool_id_bech32":"`+govPool+`","epoch_no":510,"amount":"2200000000"},{"pool_id_bech32":"`+govPool+`","epoch_no":511,"amount":null}]`)
	epoch := koios.EpochNo(510)
	res, err := client.GetPoolVotingPowerHistory(context.Background(), &epoch, koios.PoolID(govPool), nil)
	govNoErr(t, err)
	govExpect(t, cap, "GET", "/pool_voting_power_history", map[string]string{"_epoch_no": "510", "_pool_bech32": govPool}, "")
	govEq(t, "pool", res.Data[0].PoolID, koios.PoolID(govPool))
	govEq(t, "epoch", res.Data[0].EpochNo, koios.EpochNo(510))
	govEq(t, "amount", res.Data[0].Amount.String(), "2200000000")
	if res.Data[1].Amount != nil {
		t.Errorf("null amount should decode to nil")
	}
	_, err = client.GetPoolVotingPowerHistory(context.Background(), nil, "", nil)
	govNoErr(t, err)
	govExpect(t, cap, "GET", "/pool_voting_power_history", nil, "")
}

func TestGovernanceRequestOptionsPreserved(t *testing.T) {
	client, cap := govClient(t, `[]`)
	opts := client.NewRequestOptions()
	opts.QuerySet("order", "block_time.desc")
	_, err := client.GetDRepVotes(context.Background(), koios.DRepID(govDRep), opts)
	govNoErr(t, err)
	govExpect(t, cap, "GET", "/drep_votes", map[string]string{"order": "block_time.desc", "_drep_id": govDRep}, "")
}
