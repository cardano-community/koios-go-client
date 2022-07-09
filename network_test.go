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
func TestNetworkSuite(t *testing.T) {
	testsuite := &networkTestSuite{}
	testsuite.LoadSpecs([]string{
		"endpoint_network_tip",
		"endpoint_network_genesis",
		"endpoint_network_totals",
	})

	suite.Run(t, testsuite)
}

type networkTestSuite struct {
	endpointsTestSuite
}

func (s *networkTestSuite) TestNetworkTipEndpoint() {
	res, err := s.api.GetTip(context.Background(), nil)
	if s.NoError(err) {
		s.Greater(res.Data.AbsSlot, uint64(0))
		s.Greater(res.Data.BlockNo, uint64(0))
		s.Greater(res.Data.EpochNo, uint64(0))
		s.Greater(res.Data.EpochSlot, uint64(0))
		s.NotEmpty(res.Data.Hash)
		s.NotZero(res.Data.BlockTime)
	}
}

func (s *networkTestSuite) TestNetworkGenesiEndpoint() {
	res, err := s.api.GetGenesis(context.Background(), nil)
	if s.NoError(err) {
		s.False(res.Data.Activeslotcoeff.IsZero())
		s.False(res.Data.Activeslotcoeff.IsZero())
		s.False(res.Data.Epochlength.IsZero())
		s.False(res.Data.Maxkesrevolutions.IsZero())
		s.False(res.Data.Slotlength.IsZero())
		s.False(res.Data.Slotsperkesperiod.IsZero())
		s.False(res.Data.Systemstart.IsZero())
		s.Equal(decimal.New(45000000000000000, 0), res.Data.Maxlovelacesupply)
	}
}

func (s *networkTestSuite) TestNetworkTotalsEndpoint() {
	spec := s.GetSpec("endpoint_network_totals")
	if s.NotNil(spec) {
		epochNo, err := strconv.ParseUint(spec.Request.Query.Get("_epoch_no"), 10, 64)
		s.NoError(err)
		epoch := koios.EpochNo(epochNo)

		res, err := s.api.GetTotals(context.Background(), &epoch, nil)
		if s.NoError(err) {
			s.Equal(epoch, res.Data[0].Epoch)
			s.False(res.Data[0].Circulation.IsZero())
			s.False(res.Data[0].Reserves.IsZero())
			s.False(res.Data[0].Reward.IsZero())
			s.False(res.Data[0].Supply.IsZero())
			s.False(res.Data[0].Treasury.IsZero())
		}

		res2, err := s.api.GetTotals(context.Background(), nil, nil)
		if s.NoError(err) {
			s.Equal(epoch, res2.Data[0].Epoch)
			s.False(res2.Data[0].Circulation.IsZero())
			s.False(res2.Data[0].Reserves.IsZero())
			s.False(res2.Data[0].Reward.IsZero())
			s.False(res2.Data[0].Supply.IsZero())
			s.False(res2.Data[0].Treasury.IsZero())
		}
	}
}
