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

	"github.com/cardano-community/koios-go-client/v2"
	"github.com/shopspring/decimal"
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
			s.Equal("pool12fclephansjkz0qn339w7mu9jwef2ty439as08avaw7fuyk56j6", res.Data.DelegatedPool.String())
			s.True(res.Data.Reserves.Equal(koios.ZeroLovelace))
			s.True(res.Data.Rewards.IsPositive())
			s.True(res.Data.RewardsAvailable.IsPositive())
			s.Equal("registered", res.Data.Status)
			s.True(res.Data.TotalBalance.IsPositive())
			s.True(res.Data.Treasury.Equal(koios.ZeroLovelace))
			s.True(res.Data.UTxO.IsPositive())
			s.True(res.Data.Withdrawals.IsPositive())
		}
	}

	res2, err2 := s.api.GetAccountInfo(
		context.Background(),
		koios.Address(""),
		nil,
	)
	if s.Error(err2) && s.NotNil(res2) {
		s.ErrorIs(res2.Error, koios.ErrNoAddress)
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
			koios.StakeAddress(spec.Request.Query.Get("_stake_address")),
			&epoch,
			nil,
		)

		if s.NoError(err) {
			s.True(res.Data[0].Amount.IsPositive())
			s.Greater(res.Data[0].EarnedEpoch, koios.EpochNo(0))
			s.NotEmpty(res.Data[0].PoolID)
			s.Greater(res.Data[0].SpendableEpoch, koios.EpochNo(0))
			s.Equal("leader", res.Data[0].Type)
		}
	}

	res2, err2 := s.api.GetAccountRewards(
		context.Background(),
		koios.StakeAddress(""),
		nil,
		nil,
	)
	if s.Error(err2) && s.NotNil(res2) {
		s.ErrorIs(res2.Error, koios.ErrNoAddress)
	}
}

func (s *accountTestSuite) TestAccountUpdatesEndpoint() {
	spec := s.GetSpec("endpoint_account_updates")
	if s.NotNil(spec) {
		res, err := s.api.GetAccountUpdates(
			context.Background(),
			koios.StakeAddress(spec.Request.Query.Get("_stake_address")),
			nil,
		)
		if s.NoError(err) {
			s.NotEmpty(res.Data[0].ActionType)
			s.NotEmpty(res.Data[0].TxHash)
			s.Len(res.Data, 3)
		}
	}

	res2, err2 := s.api.GetAccountUpdates(
		context.Background(),
		koios.StakeAddress(""),
		nil,
	)
	if s.Error(err2) && s.NotNil(res2) {
		s.ErrorIs(res2.Error, koios.ErrNoAddress)
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
		s.Equal(
			koios.Address(
				"addr1qxeal3u4jlknqpmpkswyqmq58ggzrytu6v5rg73p64eppaxgvhcs6e4zud6jg267l8c940yggr2pxssestmzcdwwf7lsswrvnt",
			),
			res.Data[0],
		)
	}
	res2, err2 := s.api.GetAccountAddresses(
		context.Background(),
		koios.StakeAddress(""),
		nil,
	)
	if s.Error(err2) && s.NotNil(res2) {
		s.ErrorIs(res2.Error, koios.ErrNoAddress)
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
			s.Equal("DONTSPAM", res.Data[0].Name)
			s.Equal("d3501d9531fcc25e3ca4b6429318c2cc374dbdbcf5e99c1c1e5da1ff", res.Data[0].PolicyID.String())
			s.True(res.Data[0].Quantity.Equal(decimal.New(9900000, 0)))
		}
	}

	res2, err2 := s.api.GetAccountAssets(
		context.Background(),
		koios.StakeAddress(""),
		nil,
	)
	if s.Error(err2) && s.NotNil(res2) {
		s.ErrorIs(res2.Error, koios.ErrNoAddress)
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
			s.Equal(koios.EpochNo(213), res.Data[0].Epoch)
			s.Equal("pool12fclephansjkz0qn339w7mu9jwef2ty439as08avaw7fuyk56j6", res.Data[0].PoolID.String())
			s.NotEmpty("stake1u8yxtugdv63wxafy9d00nuz6hjyyp4qnggvc9a3vxh8yl0ckml2uz", res.Data[0].StakeAddress.String())
		}
	}

	res2, err2 := s.api.GetAccountHistory(
		context.Background(),
		koios.StakeAddress(""),
		nil,
	)
	if s.Error(err2) && s.NotNil(res2) {
		s.ErrorIs(res2.Error, koios.ErrNoAddress)
	}
}

// nolint: funlen
func (s *accountTestSuite) TestAccountErrors() {
	stakeaddr := "stake1u8yxtugdv63wxafy9d00nuz6hjyyp4qnggvc9a3vxh8yl0ckml2uz"
	opts := s.api.NewRequestOptions()
	opts.QueryAdd("test-http", "400")

	res1, err1 := s.api.GetAccountList(context.Background(), opts)
	if s.Error(err1) && s.NotNil(res1) {
		s.Nil(res1.Data)
		s.Equal("http error: 400 Bad Request", res1.Error.Message)
		s.Equal(koios.ErrorCode("400"), res1.Error.Code)
		s.Equal(400, res1.Error.Code.Int())
		s.ErrorIs(res1.Error, koios.ErrResponse)
	}

	res2, err2 := s.api.GetAccountInfo(
		context.Background(),
		koios.Address(stakeaddr),
		opts.Clone(),
	)
	if s.Error(err2) && s.NotNil(res2) {
		s.Nil(res2.Data)
		s.Equal("http error: 400 Bad Request", res2.Error.Message)
		s.Equal(koios.ErrorCode("400"), res2.Error.Code)
		s.Equal(400, res2.Error.Code.Int())
		s.ErrorIs(res2.Error, koios.ErrResponse)
	}

	res3, err3 := s.api.GetAccountRewards(
		context.Background(),
		koios.StakeAddress(stakeaddr),
		nil,
		opts.Clone(),
	)
	if s.Error(err3) && s.NotNil(res3) {
		s.Nil(res3.Data)
		s.Equal("http error: 400 Bad Request", res3.Error.Message)
		s.Equal(koios.ErrorCode("400"), res3.Error.Code)
		s.Equal(400, res3.Error.Code.Int())
		s.ErrorIs(res3.Error, koios.ErrResponse)
	}

	res4, err4 := s.api.GetAccountUpdates(
		context.Background(),
		koios.StakeAddress(stakeaddr),
		opts.Clone(),
	)
	if s.Error(err4) && s.NotNil(res4) {
		s.Nil(res4.Data)
		s.Equal("http error: 400 Bad Request", res4.Error.Message)
		s.Equal(koios.ErrorCode("400"), res4.Error.Code)
		s.Equal(400, res4.Error.Code.Int())
		s.ErrorIs(res4.Error, koios.ErrResponse)
	}

	res5, err5 := s.api.GetAccountAddresses(
		context.Background(),
		koios.StakeAddress(stakeaddr),
		opts.Clone(),
	)
	if s.Error(err5) && s.NotNil(res5) {
		s.Nil(res5.Data)
		s.Equal("http error: 400 Bad Request", res5.Error.Message)
		s.Equal(koios.ErrorCode("400"), res5.Error.Code)
		s.Equal(400, res5.Error.Code.Int())
		s.ErrorIs(res5.Error, koios.ErrResponse)
	}

	res6, err6 := s.api.GetAccountAssets(
		context.Background(),
		koios.StakeAddress(stakeaddr),
		opts.Clone(),
	)
	if s.Error(err6) && s.NotNil(res6) {
		s.Nil(res6.Data)
		s.Equal("http error: 400 Bad Request", res6.Error.Message)
		s.Equal(koios.ErrorCode("400"), res6.Error.Code)
		s.Equal(400, res6.Error.Code.Int())
		s.ErrorIs(res6.Error, koios.ErrResponse)
	}

	res7, err7 := s.api.GetAccountHistory(
		context.Background(),
		koios.StakeAddress(stakeaddr),
		opts.Clone(),
	)
	if s.Error(err7) && s.NotNil(res7) {
		s.Nil(res7.Data)
		s.Equal("http error: 400 Bad Request", res7.Error.Message)
		s.Equal(koios.ErrorCode("400"), res7.Error.Code)
		s.Equal(400, res7.Error.Code.Int())
		s.ErrorIs(res7.Error, koios.ErrResponse)
	}
}
