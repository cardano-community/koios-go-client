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

// import (
// 	"context"
// 	"fmt"
// 	"strconv"
// 	"testing"

// 	"github.com/cardano-community/koios-go-client/v2"
// 	"github.com/shopspring/decimal"
// 	"github.com/stretchr/testify/suite"
// )

// // In order for 'go test' to run this suite, we need to create
// // a normal test function and pass our suite to suite.Run.
// func TestAccountSuite(t *testing.T) {
// 	testsuite := &accountTestSuite{}
// 	testsuite.LoadSpecs([]string{
// 		"account_endpoint_account_addresses",
// 		"account_endpoint_account_assets",
// 		"account_endpoint_account_history",
// 		"account_endpoint_account_info",
// 		"account_endpoint_account_list",
// 		"account_endpoint_account_rewards",
// 		"account_endpoint_account_updates",
// 	})

// 	suite.Run(t, testsuite)
// }

// type accountTestSuite struct {
// 	endpointsTestSuite
// }

// func (s *accountTestSuite) TestAccountListEndpoint() {
// 	mres, err := s.mainnet.GetAccountList(context.Background(), nil)
// 	if s.NoError(err) {
// 		s.Len(mres.Data, 1000)
// 	}

// 	tres, err := s.testnet.GetAccountList(context.Background(), nil)
// 	if s.NoError(err) {
// 		s.Len(tres.Data, 1000)
// 	}
// }

// func (s *accountTestSuite) TestAccountInfoEndpoint() {
// 	specfile := "account_endpoint_account_info"
// 	for _, network := range s.networks {
// 		s.Run(fmt.Sprintf("%s %s", network, specfile), func() {
// 			// actual test
// 			spec := s.GetSpec(specfile, network)
// 			if s.NotNil(spec) {
// 				res, err := s.Client(network).GetAccountInfo(
// 					context.Background(),
// 					koios.Address(spec.Request.Query.Get("_address")),
// 					nil,
// 				)
// 				if s.NoError(err) {
// 					s.Contains(res.Data.DelegatedPool.String(), "pool")
// 					s.True(res.Data.Reserves.Equal(koios.ZeroLovelace))
// 					s.True(res.Data.Rewards.IsPositive())
// 					s.True(res.Data.RewardsAvailable.IsPositive())
// 					s.Equal("registered", res.Data.Status)
// 					s.True(res.Data.TotalBalance.IsPositive())
// 					s.True(res.Data.Treasury.Equal(koios.ZeroLovelace))
// 					s.True(res.Data.UTxO.IsPositive())
// 					s.True(res.Data.Withdrawals.IsPositive())
// 				}
// 			}
// 			res2, err2 := s.Client(network).GetAccountInfo(
// 				context.Background(),
// 				koios.Address(""),
// 				nil,
// 			)
// 			if s.Error(err2) && s.NotNil(res2) {
// 				s.ErrorIs(res2.Error, koios.ErrNoAddress)
// 			}
// 		})
// 	}
// }

// func (s *accountTestSuite) TestAccountRewardsEndpoint() {
// 	specfile := "account_endpoint_account_rewards"
// 	for _, network := range s.networks {
// 		s.Run(fmt.Sprintf("%s %s", network, specfile), func() {
// 			// actual test
// 			spec := s.GetSpec(specfile, network)
// 			if s.NotNil(spec) {
// 				epochNo, err := strconv.ParseUint(spec.Request.Query.Get("_epoch_no"), 10, 64)
// 				s.NoError(err)
// 				epoch := koios.EpochNo(epochNo)

// 				res, err := s.Client(network).GetAccountRewards(
// 					context.Background(),
// 					koios.StakeAddress(spec.Request.Query.Get("_stake_address")),
// 					&epoch,
// 					nil,
// 				)
// 				if s.NoError(err) {
// 					s.True(res.Data[0].Amount.IsPositive())
// 					s.Greater(res.Data[0].EarnedEpoch, koios.EpochNo(0))
// 					s.NotEmpty(res.Data[0].PoolID)
// 					s.Greater(res.Data[0].SpendableEpoch, koios.EpochNo(0))
// 					if network == "mainnet" {
// 						s.Equal("leader", res.Data[0].Type)
// 					} else {
// 						s.Equal("member", res.Data[0].Type)
// 					}
// 				}
// 			}
// 			res2, err2 := s.Client(network).GetAccountRewards(
// 				context.Background(),
// 				koios.StakeAddress(""),
// 				nil,
// 				nil,
// 			)
// 			if s.Error(err2) && s.NotNil(res2) {
// 				s.ErrorIs(res2.Error, koios.ErrNoAddress)
// 			}
// 		})
// 	}
// }

// func (s *accountTestSuite) TestAccountUpdatesEndpoint() {
// 	specfile := "account_endpoint_account_updates"
// 	for _, network := range s.networks {
// 		s.Run(fmt.Sprintf("%s %s", network, specfile), func() {
// 			// actual test
// 			spec := s.GetSpec(specfile, network)
// 			if s.NotNil(spec) {
// 				res, err := s.Client(network).GetAccountUpdates(
// 					context.Background(),
// 					koios.StakeAddress(spec.Request.Query.Get("_stake_address")),
// 					nil,
// 				)
// 				if s.NoError(err) {
// 					s.NotEmpty(res.Data[0].ActionType)
// 					s.NotEmpty(res.Data[0].TxHash)
// 					if network == "mainnet" {
// 						s.Len(res.Data, 3)
// 					} else {
// 						s.Len(res.Data, 2)
// 					}
// 				}
// 			}
// 			res2, err2 := s.Client(network).GetAccountUpdates(
// 				context.Background(),
// 				koios.StakeAddress(""),
// 				nil,
// 			)
// 			if s.Error(err2) && s.NotNil(res2) {
// 				s.ErrorIs(res2.Error, koios.ErrNoAddress)
// 			}
// 		})
// 	}
// }

// func (s *accountTestSuite) TestAccountAddressesEndpoint() {
// 	specfile := "account_endpoint_account_addresses"
// 	for _, network := range s.networks {
// 		s.Run(fmt.Sprintf("%s %s", network, specfile), func() {
// 			// actual test
// 			spec := s.GetSpec(specfile, network)
// 			if s.NotNil(spec) {
// 				res, err := s.Client(network).GetAccountAddresses(
// 					context.Background(),
// 					koios.StakeAddress(spec.Request.Query.Get("_address")),
// 					nil,
// 				)
// 				if s.NoError(err) {
// 					s.Greater(len(res.Data), 0)
// 				}
// 				if network == "mainnet" {
// 					s.Equal(
// 						koios.Address(
// 							"addr1qxeal3u4jlknqpmpkswyqmq58ggzrytu6v5rg73p64eppaxgvhcs6e4zud6jg267l8c940yggr2pxssestmzcdwwf7lsswrvnt",
// 						),
// 						res.Data[0],
// 					)
// 				} else {
// 					s.Equal(
// 						koios.Address(
// 							"addr_test1qqreqxry8atyrv0wzz8q90mzky782pt5dkhqjwuj4ccs560l0dw5r75vk42mv3ykq8vyjeaanvpytg79xqzymqy5acmq28qk2d",
// 						),
// 						res.Data[0],
// 					)
// 				}
// 			}

// 			res2, err2 := s.Client(network).GetAccountAddresses(
// 				context.Background(),
// 				koios.StakeAddress(""),
// 				nil,
// 			)
// 			if s.Error(err2) && s.NotNil(res2) {
// 				s.ErrorIs(res2.Error, koios.ErrNoAddress)
// 			}
// 		})
// 	}
// }

// func (s *accountTestSuite) TestAccountAssetsEndpoint() {
// 	specfile := "account_endpoint_account_assets"
// 	for _, network := range s.networks {
// 		s.Run(fmt.Sprintf("%s %s", network, specfile), func() {
// 			// actual test
// 			spec := s.GetSpec(specfile, network)
// 			if s.NotNil(spec) {
// 				res, err := s.Client(network).GetAccountAssets(
// 					context.Background(),
// 					koios.StakeAddress(spec.Request.Query.Get("_address")),
// 					nil,
// 				)
// 				if s.NoError(err) {
// 					if network == "mainnet" {
// 						if s.Len(res.Data, 1) {
// 							s.Equal("DONTSPAM", res.Data[0].Name)
// 							s.Equal("d3501d9531fcc25e3ca4b6429318c2cc374dbdbcf5e99c1c1e5da1ff", res.Data[0].PolicyID.String())
// 							s.True(res.Data[0].Quantity.Equal(decimal.New(9900000, 0)))
// 						}
// 					} else {
// 						if s.Len(res.Data, 506) {
// 							s.Equal("SpacePet1505", res.Data[0].Name)
// 							s.Equal("03e2ed6c0c9de902a156e7dac279dd1cc76f75607fcd9cd1d218cd92", res.Data[0].PolicyID.String())
// 							s.True(res.Data[0].Quantity.Equal(decimal.New(1, 0)))
// 						}
// 					}
// 				}
// 			}

// 			res2, err2 := s.Client(network).GetAccountAssets(
// 				context.Background(),
// 				koios.StakeAddress(""),
// 				nil,
// 			)
// 			if s.Error(err2) && s.NotNil(res2) {
// 				s.ErrorIs(res2.Error, koios.ErrNoAddress)
// 			}
// 		})
// 	}
// }

// func (s *accountTestSuite) TestAccountHistoryEndpoint() {
// 	specfile := "account_endpoint_account_history"
// 	for _, network := range s.networks {
// 		s.Run(fmt.Sprintf("%s %s", network, specfile), func() {
// 			// actual test
// 			spec := s.GetSpec(specfile, network)
// 			if s.NotNil(spec) {
// 				res, err := s.Client(network).GetAccountHistory(
// 					context.Background(),
// 					koios.StakeAddress(spec.Request.Query.Get("_address")),
// 					nil,
// 				)
// 				if s.NoError(err) {
// 					if network == "mainnet" {
// 						s.Equal(decimal.NewFromInt(30097808275), res.Data[0].ActiveStake)
// 						s.Equal(koios.EpochNo(213), res.Data[0].Epoch)
// 						s.Equal("pool12fclephansjkz0qn339w7mu9jwef2ty439as08avaw7fuyk56j6", res.Data[0].PoolID.String())
// 						s.NotEmpty("stake1u8yxtugdv63wxafy9d00nuz6hjyyp4qnggvc9a3vxh8yl0ckml2uz", res.Data[0].StakeAddress.String())
// 					} else {
// 						s.Equal(decimal.NewFromInt(1939988859), res.Data[0].ActiveStake)
// 						s.Equal(koios.EpochNo(81), res.Data[0].Epoch)
// 						s.Equal("pool1frevxe70aqw2ce58c0muyesnahl88nfjjsp25h85jwakzgd2g2l", res.Data[0].PoolID.String())
// 						s.NotEmpty("stake_test1urlhkh2pl2xt24dkgjtqrkzfv77ekqj950znqpzdsz2wuds0xlsk6", res.Data[0].StakeAddress.String())
// 					}
// 				}
// 			}

// 			res2, err2 := s.Client(network).GetAccountHistory(
// 				context.Background(),
// 				koios.StakeAddress(""),
// 				nil,
// 			)
// 			if s.Error(err2) && s.NotNil(res2) {
// 				s.ErrorIs(res2.Error, koios.ErrNoAddress)
// 			}
// 		})
// 	}
// }

// // nolint: funlen
// func (s *accountTestSuite) TestAccountErrors() {

// 	checkErrorRes := func(s *accountTestSuite, res koios.Response) {
// 		s.Equal("http error: 400 Bad Request", res.Error.Message)
// 		s.Equal(koios.ErrorCode("400"), res.Error.Code)
// 		s.Equal(400, res.Error.Code.Int())
// 		s.ErrorIs(res.Error, koios.ErrResponse)
// 	}

// 	for _, network := range s.networks {
// 		var stakeaddr string
// 		if network == "mainnet" {
// 			stakeaddr = "stake1u8yxtugdv63wxafy9d00nuz6hjyyp4qnggvc9a3vxh8yl0ckml2uz"
// 		} else {
// 			stakeaddr = "stake_test1urlhkh2pl2xt24dkgjtqrkzfv77ekqj950znqpzdsz2wuds0xlsk6"
// 		}

// 		opts := s.Client(network).NewRequestOptions()
// 		opts.QueryAdd("test-http", "400")

// 		s.Run(fmt.Sprintf("%s GetAccountList", network), func() {
// 			res, err := s.Client(network).GetAccountList(context.Background(), opts)
// 			if s.Error(err) && s.NotNil(res) {
// 				s.Nil(res.Data)
// 				checkErrorRes(s, res.Response)
// 			}
// 		})

// 		s.Run(fmt.Sprintf("%s GetAccountInfo", network), func() {
// 			res, err := s.Client(network).GetAccountInfo(
// 				context.Background(),
// 				koios.Address(stakeaddr),
// 				opts.Clone(),
// 			)
// 			if s.Error(err) && s.NotNil(res) {
// 				s.Nil(res.Data)
// 				checkErrorRes(s, res.Response)
// 			}
// 		})

// 		s.Run(fmt.Sprintf("%s GetAccountRewards", network), func() {
// 			res, err := s.Client(network).GetAccountRewards(
// 				context.Background(),
// 				koios.StakeAddress(stakeaddr),
// 				nil,
// 				opts.Clone(),
// 			)
// 			if s.Error(err) && s.NotNil(res) {
// 				s.Nil(res.Data)
// 				checkErrorRes(s, res.Response)
// 			}
// 		})

// 		s.Run(fmt.Sprintf("%s GetAccountUpdates", network), func() {
// 			res, err := s.Client(network).GetAccountUpdates(
// 				context.Background(),
// 				koios.StakeAddress(stakeaddr),
// 				opts.Clone(),
// 			)
// 			if s.Error(err) && s.NotNil(res) {
// 				s.Nil(res.Data)
// 				checkErrorRes(s, res.Response)
// 			}
// 		})

// 		s.Run(fmt.Sprintf("%s GetAccountAddresses", network), func() {
// 			res, err := s.Client(network).GetAccountAddresses(
// 				context.Background(),
// 				koios.StakeAddress(stakeaddr),
// 				opts.Clone(),
// 			)
// 			if s.Error(err) && s.NotNil(res) {
// 				s.Nil(res.Data)
// 				checkErrorRes(s, res.Response)
// 			}
// 		})

// 		s.Run(fmt.Sprintf("%s GetAccountAssets", network), func() {
// 			res, err := s.Client(network).GetAccountAssets(
// 				context.Background(),
// 				koios.StakeAddress(stakeaddr),
// 				opts.Clone(),
// 			)
// 			if s.Error(err) && s.NotNil(res) {
// 				s.Nil(res.Data)
// 				checkErrorRes(s, res.Response)
// 			}
// 		})

// 		s.Run(fmt.Sprintf("%s GetAccountHistory", network), func() {
// 			res, err := s.Client(network).GetAccountHistory(
// 				context.Background(),
// 				koios.StakeAddress(stakeaddr),
// 				opts.Clone(),
// 			)
// 			if s.Error(err) && s.NotNil(res) {
// 				s.Nil(res.Data)
// 				checkErrorRes(s, res.Response)
// 			}
// 		})
// 	}
// }
