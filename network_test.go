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

// // In order for 'go test' to run this suite, we need to create
// // a normal test function and pass our suite to suite.Run.
// func TestNetworkSuite(t *testing.T) {
// 	testsuite := &networkTestSuite{}
// 	testsuite.LoadSpecs([]string{
// 		"network_endpoint_tip",
// 		"network_endpoint_genesis",
// 		"network_endpoint_totals",
// 	})

// 	suite.Run(t, testsuite)
// }

// type networkTestSuite struct {
// 	endpointsTestSuite
// }

// func (s *networkTestSuite) TestNetworkTipEndpoint() {
// 	for _, network := range s.networks {
// 		s.Run(network, func() {
// 			res, err := s.Client(network).GetTip(context.Background(), nil)
// 			if s.NoError(err) {
// 				s.Greater(res.Data.AbsSlot, uint64(0))
// 				s.Greater(res.Data.BlockNo, uint64(0))
// 				s.Greater(res.Data.EpochNo, uint64(0))
// 				s.Greater(res.Data.EpochSlot, uint64(0))
// 				s.NotEmpty(res.Data.BlockHash)
// 				s.NotZero(res.Data.BlockTime)
// 			}
// 		})
// 	}
// }

// func (s *networkTestSuite) TestNetworkGenesiEndpoint() {
// 	for _, network := range s.networks {
// 		s.Run(network, func() {
// 			res, err := s.Client(network).GetGenesis(context.Background(), nil)
// 			if s.NoError(err) {
// 				s.True(res.Data.NetworkMagic.IsPositive())
// 				s.Greater(res.Data.EpochLength, 0)
// 				s.Equal(decimal.NewFromInt(1), res.Data.Slotlength)
// 				s.Equal(decimal.NewFromInt(45000000000000000), res.Data.MaxLovelaceSupply)
// 				s.False(res.Data.Systemstart.IsZero())
// 				s.Equal(decimal.NewFromFloat(0.05), res.Data.Activeslotcoeff)
// 				s.Equal(decimal.NewFromInt(129600), res.Data.Slotsperkesperiod)
// 				s.Equal(decimal.NewFromInt(62), res.Data.Maxkesrevolutions)
// 				s.NotEmpty(res.Data.Alonzogenesis)

// 				kv, err := res.Data.AlonzoGenesisMap()
// 				if s.NoError(err) {
// 					s.Len(kv, 8)
// 				}

// 				if network == "guildnet" {
// 					s.Equal(decimal.NewFromInt(36), res.Data.Securityparam)
// 				} else {
// 					s.Equal(decimal.NewFromInt(2160), res.Data.Securityparam)
// 				}
// 			}
// 		})
// 	}
// }

// func (s *networkTestSuite) TestNetworkTotalsEndpoint() {
// 	specfile := "network_endpoint_totals"
// 	for _, network := range s.networks {
// 		s.Run(network, func() {
// 			// actual test
// 			spec := s.GetSpec(specfile, network)
// 			if s.NotNil(spec) {
// 				epochNo, err := strconv.ParseUint(spec.Request.Query.Get("_epoch_no"), 10, 64)
// 				s.NoError(err)
// 				epoch := koios.EpochNo(epochNo)
// 				res, err := s.Client(network).GetTotals(context.Background(), &epoch, nil)
// 				if s.NoError(err) {
// 					s.Equal(epoch, res.Data[0].Epoch)
// 					s.False(res.Data[0].Circulation.IsZero())
// 					s.False(res.Data[0].Reserves.IsZero())
// 					s.False(res.Data[0].Reward.IsZero())
// 					s.False(res.Data[0].Supply.IsZero())
// 					s.False(res.Data[0].Treasury.IsZero())
// 				}

// 				res2, err := s.Client(network).GetTotals(context.Background(), nil, nil)
// 				if s.NoError(err) {
// 					s.Equal(epoch, res2.Data[0].Epoch)
// 					s.False(res2.Data[0].Circulation.IsZero())
// 					s.False(res2.Data[0].Reserves.IsZero())
// 					s.False(res2.Data[0].Reward.IsZero())
// 					s.False(res2.Data[0].Supply.IsZero())
// 					s.False(res2.Data[0].Treasury.IsZero())
// 				}
// 			}
// 		})
// 	}
// }

// func (s *networkTestSuite) TestNetworkErrors() {
// 	for _, network := range s.networks {
// 		opts := s.Client(network).NewRequestOptions()
// 		opts.QueryAdd("test-http", "400")

// 		s.Run("GetTip", func() {
// 			res, err := s.Client(network).GetTip(context.Background(), opts.Clone())
// 			if s.Error(err) && s.NotNil(res) {
// 				s.Nil(res.Data)
// 				s.KoiosHttpError(res.Response)
// 			}
// 		})

// 		s.Run("GetGenesis", func() {
// 			res, err := s.Client(network).GetGenesis(context.Background(), opts.Clone())
// 			if s.Error(err) && s.NotNil(res) {
// 				s.Nil(res.Data)
// 				s.KoiosHttpError(res.Response)
// 			}
// 		})

// 		s.Run("GetTotals", func() {
// 			res, err := s.Client(network).GetTotals(context.Background(), nil, opts.Clone())
// 			if s.Error(err) && s.NotNil(res) {
// 				s.Nil(res.Data)
// 				s.KoiosHttpError(res.Response)
// 			}
// 		})
// 	}
// }

// func (s *networkTestSuite) TestAlonzoGenesisMapUnmarshal() {
// 	res, err := s.Client("mainnet").GetGenesis(context.Background(), nil)
// 	if s.NoError(err) {
// 		res.Data.Alonzogenesis += "-" // invalidate json
// 		kv, err := res.Data.AlonzoGenesisMap()
// 		s.Error(err)
// 		s.Nil(kv)
// 	}
// }
