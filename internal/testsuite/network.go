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

package testsuite

import (
	"context"
	"os"

	"github.com/cardano-community/koios-go-client/v2"
	"github.com/stretchr/testify/assert"
)

func NetworkGenesisTest(t TestingT, client *koios.Client) {
	res, err := client.GetGenesis(context.Background(), nil)
	if !assert.NoError(t, err) {
		return
	}
	assert.True(t, res.Data.NetworkMagic.IsPositive())
	assert.NotEmpty(t, res.Data.NetworkID)
	assert.True(t, res.Data.EpochLength.IsPositive())
	assert.True(t, res.Data.SlotLength.IsPositive())
	assert.True(t, res.Data.MaxLovelaceSupply.IsPositive())
	assert.False(t, res.Data.SystemStart.IsZero())
	assert.True(t, res.Data.ActiveSlotCoeff.IsPositive())
	assert.True(t, res.Data.SlotsPerKesPeriod.IsPositive())
	assert.True(t, res.Data.MaxKesRevolutions.IsPositive())
	assert.True(t, res.Data.SecurityParam.IsPositive())
	assert.True(t, res.Data.UpdateQuorum.IsPositive())
	assert.NotEmpty(t, res.Data.AlonzoGenesis)
}

func NetworkTipTest(t TestingT, client *koios.Client) {
	tip, err := client.GetTip(context.Background(), nil)
	assert.NoError(t, err)

	assert.Greater(t, tip.Data.AbsSlot, int(100000))
	assert.Greater(t, tip.Data.BlockNo, int(100000))
	assert.Greater(t, tip.Data.EpochNo, int(230))
	assert.Greater(t, tip.Data.EpochSlot, int(1))
	assert.NotEmpty(t, tip.Data.BlockHash)
	assert.NotZero(t, tip.Data.BlockTime)
	assert.False(t, tip.Data.BlockTime.IsZero())
}

func NetworkTotalsTest(t TestingT, client *koios.Client) {
	var epoch koios.EpochNo
	switch os.Getenv("KOIOS_NETWORK") {
	case "mainnet":
		epoch = koios.EpochNo(320)
	case "guild":
		epoch = koios.EpochNo(1950)
	case "testnet":
		epoch = koios.EpochNo(185)
	}

	res, err := client.GetTotals(context.Background(), &epoch, nil)
	if !assert.NoError(t, err) {
		return
	}
	zero := int64(0)
	assert.Equal(t, epoch, res.Data[0].EpochNo)
	assert.Greater(t, res.Data[0].Circulation.IntPart(), zero)
	assert.Greater(t, res.Data[0].Reserves.IntPart(), zero)
	assert.Greater(t, res.Data[0].Reward.IntPart(), zero)
	assert.Greater(t, res.Data[0].Supply.IntPart(), zero)
	assert.Greater(t, res.Data[0].Treasury.IntPart(), zero)
}
