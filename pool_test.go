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
// 	"encoding/json"
// 	"io"
// 	"strconv"
// 	"testing"

// 	"github.com/cardano-community/koios-go-client/v2"
// 	"github.com/stretchr/testify/suite"
// )

// // In order for 'go test' to run this suite, we need to create
// // a normal test function and pass our suite to suite.Run.
// func TestPoolSuite(t *testing.T) {
// 	testsuite := &poolTestSuite{}
// 	testsuite.LoadSpecs([]string{
// 		"endpoint_pool_blocks",
// 		"endpoint_pool_delegators",
// 		"endpoint_pool_history",
// 		"endpoint_pool_info",
// 		"endpoint_pool_list",
// 		"endpoint_pool_metadata",
// 		"endpoint_pool_relays",
// 		"endpoint_pool_updates",
// 	})

// 	suite.Run(t, testsuite)
// }

// type poolTestSuite struct {
// 	endpointsTestSuite
// }

// func (s *poolTestSuite) TestGetPoolBlocksEndpoint() {
// 	spec := s.GetSpec("endpoint_pool_blocks")
// 	if s.NotNil(spec) {
// 		epochNo, err := strconv.ParseUint(spec.Request.Query.Get("_epoch_no"), 10, 64)
// 		s.NoError(err)
// 		epoch := koios.EpochNo(epochNo)

// 		res, err := s.api.GetPoolBlocks(
// 			context.Background(),
// 			koios.PoolID(spec.Request.Query.Get("_pool_bech32")),
// 			&epoch,
// 			nil,
// 		)
// 		if s.NoError(err) {
// 			s.Greater(len(res.Data), 0)

// 			s.Greater(res.Data[0].AbsSlot, uint64(0))
// 			s.NotEmpty(res.Data[0].Hash)
// 			s.Greater(res.Data[0].Height, uint64(0))
// 			s.False(res.Data[0].Time.IsZero())
// 			s.Greater(res.Data[0].EpochSlot, uint64(0))
// 			s.Equal(epoch, res.Data[0].Epoch)
// 		}
// 	}
// }
// func (s *poolTestSuite) TestGetPoolDelegatorsEndpoint() {
// 	spec := s.GetSpec("endpoint_pool_delegators")
// 	if s.NotNil(spec) {
// 		epochNo, err := strconv.ParseUint(spec.Request.Query.Get("_epoch_no"), 10, 64)
// 		s.NoError(err)
// 		epoch := koios.EpochNo(epochNo)

// 		res, err := s.api.GetPoolDelegators(
// 			context.Background(),
// 			koios.PoolID(spec.Request.Query.Get("_pool_bech32")),
// 			&epoch,
// 			nil,
// 		)
// 		if s.NoError(err) {
// 			s.NotEmpty(res.Data[0].StakeAddress)
// 			s.Equal(epoch, res.Data[0].Epoch)
// 			s.True(res.Data[0].Amount.IsPositive())
// 		}
// 	}
// }

// func (s *poolTestSuite) TestGetPoolHistoryEndpoint() {
// 	spec := s.GetSpec("endpoint_pool_history")
// 	if s.NotNil(spec) {
// 		epochNo, err := strconv.ParseUint(spec.Request.Query.Get("_epoch_no"), 10, 64)
// 		s.NoError(err)
// 		epoch := koios.EpochNo(epochNo)

// 		res, err := s.api.GetPoolHistory(
// 			context.Background(),
// 			koios.PoolID(spec.Request.Query.Get("_pool_bech32")),
// 			&epoch,
// 			nil,
// 		)
// 		if s.NoError(err) {
// 			s.True(res.Data[0].ActiveStake.IsPositive())
// 			s.Greater(res.Data[0].ActiveStakePCT, float64(0))
// 			s.Greater(res.Data[0].BlockCNT, 0)
// 			s.Greater(res.Data[0].DelegatorCNT, 0)
// 			s.Greater(res.Data[0].SaturationPCT, float64(0))
// 			rew, _ := res.Data[0].DelegRewards.Float64()
// 			s.GreaterOrEqual(rew, float64(0))
// 			s.Equal(epoch, res.Data[0].Epoch)
// 			s.True(res.Data[0].FixedCost.IsPositive())
// 			s.True(res.Data[0].PoolFees.IsPositive())
// 		}
// 	}
// }

// func (s *poolTestSuite) TestGetPoolInfoEndpoint() {
// 	spec := s.GetSpec("endpoint_pool_info")
// 	if s.NotNil(spec) {
// 		var payload = struct {
// 			PoolIDs []koios.PoolID `json:"_pool_bech32_ids"`
// 		}{}
// 		err := json.Unmarshal(spec.Request.Body, &payload)
// 		s.NoError(err)

// 		res, err := s.api.GetPoolInfo(
// 			context.Background(),
// 			payload.PoolIDs[0],
// 			nil,
// 		)
// 		if s.NoError(err) {
// 			s.Greater(res.Data.ActiveEpoch, koios.EpochNo(0))
// 			s.False(res.Data.ActiveStake.IsNegative())
// 			s.GreaterOrEqual(res.Data.BlockCount, uint64(0))
// 			s.True(res.Data.FixedCost.IsPositive())
// 			s.GreaterOrEqual(res.Data.LiveDelegators, uint64(0))
// 			s.NotEmpty(res.Data.VrfKeyHash)
// 			s.NotEmpty(res.Data.RewardAddr)
// 			s.NotEmpty(res.Data.Status)
// 			s.NotEmpty(res.Data.ID.String())
// 			s.NotEmpty(res.Data.IDHex)
// 		}

// 		res2, err := s.api.GetPoolInfos(context.Background(), []koios.PoolID{}, nil)
// 		s.ErrorIs(err, koios.ErrNoPoolID)
// 		s.Nil(res2.Data, "response data should be nil if arg is invalid")
// 		s.Equal(res2.Error.Message, "missing pool id")

// 		rpipe, w := io.Pipe()
// 		go func() {
// 			_ = json.NewEncoder(w).Encode(payload)
// 			defer w.Close()
// 		}()
// 		rsp3, err3 := s.api.POST(context.Background(), "/pool_info", rpipe, nil)

// 		if s.NoError(err3) {
// 			defer func() { _ = rsp3.Body.Close() }()
// 			res3 := &koios.PoolInfosResponse{}
// 			if s.NoError(koios.ReadAndUnmarshalResponse(rsp3, &res3.Response, &res3.Data)) {
// 				s.Greater(res3.Data[0].ActiveEpoch, koios.EpochNo(0))
// 				s.False(res3.Data[0].ActiveStake.IsNegative())
// 			}
// 		}
// 	}
// }

// func (s *poolTestSuite) TestGetPoolListEndpoint() {
// 	spec := s.GetSpec("endpoint_pool_list")
// 	if s.NotNil(spec) {
// 		res, err := s.api.GetPoolList(context.Background(), nil)
// 		if s.NoError(err) {
// 			s.Len(res.Data, 1000)
// 			s.NotEmpty(res.Data[0].PoolID)
// 			s.NotEmpty(res.Data[0].Ticker)
// 		}
// 	}
// }
// func (s *poolTestSuite) TestGetPoolMetadataEndpoint() {
// 	spec := s.GetSpec("endpoint_pool_metadata")
// 	if s.NotNil(spec) {
// 		var payload = struct {
// 			PoolIDs []koios.PoolID `json:"_pool_bech32_ids"`
// 		}{}
// 		err := json.Unmarshal(spec.Request.Body, &payload)
// 		s.NoError(err)

// 		res, err := s.api.GetPoolMetadata(context.Background(), payload.PoolIDs, nil)
// 		if s.NoError(err) {
// 			s.NotEmpty(res.Data[0].MetaHash)
// 			s.NotEmpty(res.Data[0].MetaURL)
// 			s.NotEmpty(res.Data[0].PoolID)
// 			s.NotEmpty(res.Data[0].MetaJSON.Description)
// 			s.NotEmpty(res.Data[0].MetaJSON.Homepage)
// 			s.NotEmpty(res.Data[0].MetaJSON.Name)
// 			s.NotEmpty(res.Data[0].MetaJSON.Ticker)
// 		}
// 	}
// }

// func (s *poolTestSuite) TestGetPoolRelaysEndpoint() {
// 	res, err := s.api.GetPoolRelays(context.Background(), nil)
// 	if s.NoError(err) {
// 		s.Len(res.Data, 1000)
// 		s.NotEmpty(res.Data[0].PoolID)
// 		s.Greater(len(res.Data[0].Relays), 0)
// 	}
// }

// func (s *poolTestSuite) TestGetPoolUpdatesEndpoint() {
// 	spec := s.GetSpec("endpoint_pool_metadata")
// 	if s.NotNil(spec) {
// 		poolID := koios.PoolID(spec.Request.Query.Get("_pool_bech32"))
// 		res, err := s.api.GetPoolUpdates(
// 			context.Background(),
// 			&poolID,
// 			nil,
// 		)
// 		if s.NoError(err) {
// 			s.Greater(res.Data[0].ActiveEpoch, koios.EpochNo(0))
// 			s.True(res.Data[0].FixedCost.IsPositive())
// 			s.NotEmpty(res.Data[0].VrfKeyHash)
// 			s.NotEmpty(res.Data[0].RewardAddr)
// 			s.NotEmpty(res.Data[0].Status)
// 			s.NotEmpty(res.Data[0].ID.String())
// 			s.NotEmpty(res.Data[0].IDHex)
// 			s.NotEmpty(res.Data[0].MetaHash)
// 			s.NotEmpty(res.Data[0].MetaURL)
// 			s.NotEmpty(res.Data[0].TxHash)
// 		}
// 	}
// }
