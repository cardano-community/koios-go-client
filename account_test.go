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

package koios_test

import (
	"context"
	"strconv"
	"testing"

	"github.com/cardano-community/koios-go-client"
	"github.com/stretchr/testify/suite"
)

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run.
func TestAccountSuite(t *testing.T) {
	testsuite := &accountTestSuite{}
	testsuite.LoadSpecs([]string{
		"endpoint_account_addresses",
		"endpoint_account_assets",
		"endpoint_account_history",
		"endpoint_account_info",
		"endpoint_account_list",
		"endpoint_account_rewards",
		"endpoint_account_updates",
	})

	suite.Run(t, testsuite)
}

type accountTestSuite struct {
	endpointsTestSuite
}

func (s *accountTestSuite) TestAccountListEndpoint() {
	res, err := s.api.GetAccountList(context.Background(), nil)
	if s.NoError(err) {
		s.Len(res.Data, 1000)
	}
}

func (s *accountTestSuite) TestAccountInfoEndpoint() {
	spec := s.GetSpec("endpoint_account_info")
	if s.NotNil(spec) {
		res, err := s.api.GetAccountInfo(
			context.Background(),
			koios.Address(spec.Request.Query.Get("_address")),
			nil,
		)
		if s.NoError(err) {
			s.NotEmpty(res.Data.DelegatedPool)
			s.NotEmpty(res.Data.Reserves.String())
			s.NotEmpty(res.Data.Treasury.String())
			s.True(res.Data.Rewards.IsPositive())
			s.True(res.Data.RewardsAvailable.IsPositive())
			s.NotEmpty(res.Data.Status)
			s.True(res.Data.TotalBalance.IsPositive())
			s.True(res.Data.UTxO.IsPositive())
			s.True(res.Data.Withdrawals.IsPositive())
		}
	}
}

func (s *accountTestSuite) TestAccountRewardsEndpoint() {
	spec := s.GetSpec("endpoint_account_rewards")
	if s.NotNil(spec) {
		epochNo, err := strconv.ParseUint(spec.Request.Query.Get("_epoch_no"), 10, 64)
		s.NoError(err)
		epoch := koios.EpochNo(epochNo)

		res, err := s.api.GetAccountRewards(
			context.Background(),
			koios.StakeAddress(spec.Request.Query.Get("_address")),
			&epoch,
			nil,
		)

		if s.NoError(err) {
			s.True(res.Data[0].Amount.IsPositive())
			s.Greater(res.Data[0].EarnedEpoch, koios.EpochNo(0))
			s.NotEmpty(res.Data[0].PoolID)
			s.Greater(res.Data[0].SpendableEpoch, koios.EpochNo(0))
			s.NotEmpty(res.Data[0].Type)
		}
	}
}

func (s *accountTestSuite) TestAccountUpdatesEndpoint() {
	spec := s.GetSpec("endpoint_account_updates")
	if s.NotNil(spec) {
		res, err := s.api.GetAccountUpdates(
			context.Background(),
			koios.StakeAddress(spec.Request.Query.Get("_address")),
			nil,
		)
		if s.NoError(err) {
			s.NotEmpty(res.Data[0].ActionType)
			s.NotEmpty(res.Data[0].TxHash)
		}
	}
}

func (s *accountTestSuite) TestAccountAddressesEndpoint() {
	spec := s.GetSpec("endpoint_account_addresses")
	if s.NotNil(spec) {
		res, err := s.api.GetAccountAddresses(
			context.Background(),
			koios.StakeAddress(spec.Request.Query.Get("_address")),
			nil,
		)
		if s.NoError(err) {
			s.Greater(len(res.Data), 0)
		}
	}
}

func (s *accountTestSuite) TestAccountAssetsEndpoint() {
	spec := s.GetSpec("endpoint_account_assets")
	if s.NotNil(spec) {
		res, err := s.api.GetAccountAssets(
			context.Background(),
			koios.StakeAddress(spec.Request.Query.Get("_address")),
			nil,
		)
		if s.NoError(err) && s.Greater(len(res.Data), 0) {
			s.NotEmpty(res.Data[0].Name)
			s.NotEmpty(res.Data[0].PolicyID)
			s.True(res.Data[0].Quantity.IsPositive())
		}
	}
}

func (s *accountTestSuite) TestAccountHistoryEndpoint() {
	spec := s.GetSpec("endpoint_account_history")
	if s.NotNil(spec) {
		res, err := s.api.GetAccountHistory(
			context.Background(),
			koios.StakeAddress(spec.Request.Query.Get("_address")),
			nil,
		)
		if s.NoError(err) {
			s.True(res.Data[0].ActiveStake.IsPositive())
			s.Greater(res.Data[0].Epoch, koios.EpochNo(0))
			s.NotEmpty(res.Data[0].PoolID)
			s.NotEmpty(res.Data[0].StakeAddress)
		}
	}
}
