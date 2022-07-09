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
	"github.com/stretchr/testify/suite"
)

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run.
func TestEpochSuite(t *testing.T) {
	testsuite := &epochTestSuite{}
	testsuite.LoadSpecs([]string{
		"endpoint_epoch_info",
		"endpoint_epoch_params",
	})
	suite.Run(t, testsuite)
}

type epochTestSuite struct {
	endpointsTestSuite
}

func (s *epochTestSuite) TestEpochInfoEndpoint() {
	spec := s.GetSpec("endpoint_epoch_info")
	if s.NotNil(spec) {
		epochNo, err := strconv.ParseUint(spec.Request.Query.Get("_epoch_no"), 10, 64)
		s.NoError(err)
		epoch := koios.EpochNo(epochNo)

		res, err := s.api.GetEpochInfo(context.Background(), &epoch, nil)

		if s.NoError(err) {
			s.True(res.Data[0].ActiveStake.IsPositive())
			s.Greater(res.Data[0].BlkCount, int64(0))
			s.NotZero(res.Data[0].EndTime)
			s.Equal(epoch, res.Data[0].Epoch)
			s.True(res.Data[0].Fees.IsPositive())
			s.NotZero(res.Data[0].FirstBlockTime)
			s.NotZero(res.Data[0].LastBlockTime)
			s.True(res.Data[0].OutSum.IsPositive())
			s.NotZero(res.Data[0].StartTime)
			s.Greater(res.Data[0].TxCount, int64(0))
		}
	}
}

func (s *epochTestSuite) TestEpochParamsEndpoint() {
	spec := s.GetSpec("endpoint_epoch_params")
	if s.NotNil(spec) {
		epochNo, err := strconv.ParseUint(spec.Request.Query.Get("_epoch_no"), 10, 64)
		s.NoError(err)
		epoch := koios.EpochNo(epochNo)

		res, err := s.api.GetEpochParams(context.Background(), &epoch, nil)

		if s.NoError(err) {
			s.Equal(epoch, res.Data[0].Epoch)
			s.True(res.Data[0].MinFeeA.IsPositive())
			s.True(res.Data[0].MinFeeB.IsPositive())
			s.Greater(res.Data[0].MaxBlockSize, 0)
			s.Greater(res.Data[0].MaxTxSize, 0)
			s.True(res.Data[0].KeyDeposit.IsPositive())
			s.True(res.Data[0].PoolDeposit.IsPositive())
			s.Greater(res.Data[0].MaxEpoch, 0)
			s.Greater(res.Data[0].OptimalPoolCount, 0)
			s.Greater(res.Data[0].Influence, float64(0))
			s.Greater(res.Data[0].MonetaryExpandRate, float64(0))
			s.Greater(res.Data[0].TreasuryGrowthRate, float64(0))
			s.Equal(float64(0), res.Data[0].Decentralisation)
			s.Equal("", res.Data[0].Entropy)
			s.Equal(6, res.Data[0].ProtocolMajor)
			s.Equal(0, res.Data[0].ProtocolMinor)
			s.True(res.Data[0].MinUtxoValue.IsPositive())
			s.True(res.Data[0].MinPoolCost.IsPositive())
			s.NotEmpty(res.Data[0].Nonce)
			s.NotEmpty(res.Data[0].BlockHash)
			s.NotEmpty(res.Data[0].CostModels)
			s.True(res.Data[0].PriceMem.IsPositive())
			s.True(res.Data[0].PriceStep.IsPositive())
			s.Greater(res.Data[0].MaxTxExMem, float32(0))
			s.Greater(res.Data[0].MaxTxExSteps, float32(0))
			s.Greater(res.Data[0].MaxBlockExMem, float32(0))
			s.Greater(res.Data[0].MaxBlockExSteps, float32(0))
			s.Greater(res.Data[0].MaxValSize, float32(0))
			s.Greater(res.Data[0].CollateralPercent, int64(0))
			s.Greater(res.Data[0].MaxCollateralInputs, 0)
			s.True(res.Data[0].CoinsPerUtxoWord.IsPositive())
		}
	}
}
