// Copyright 2022 The Howijd.Network Authors
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

	"github.com/test-go/testify/assert"

	"github.com/howijd/koios-rest-go-client"
)

func TestNetworkTip(t *testing.T) {
	expected := []koios.Tip{}

	spec := loadEndpointTestSpec(t, "endpoint_network_tip.json.gz", &expected)

	ts, api := createTestServerAndClient(t, spec)

	defer ts.Close()

	res, err := api.GetTip(context.TODO())

	assert.NoError(t, err)
	testHeaders(t, spec, res.Response)

	assert.Len(t, expected, 1)
	assert.Equal(t, expected[0].AbsSlot, res.Data.AbsSlot, "wrong AbsSlot")
	assert.Equal(t, expected[0].BlockNo, res.Data.BlockNo, "wrong BlockNo")
	assert.Equal(t, expected[0].BlockTime, res.Data.BlockTime, "wrong BlockTime")
	assert.Equal(t, expected[0].Epoch, res.Data.Epoch, "wrong Epoch")
	assert.Equal(t, expected[0].EpochSlot, res.Data.EpochSlot, "wrong EpochSlot")
	assert.Equal(t, expected[0].Hash, res.Data.Hash, "wrong Hash")
}

func TestNetworkGenesis(t *testing.T) {
	expected := []koios.Genesis{}

	spec := loadEndpointTestSpec(t, "endpoint_network_genesis.json.gz", &expected)

	ts, api := createTestServerAndClient(t, spec)

	defer ts.Close()

	res, err := api.GetGenesis(context.TODO())

	assert.NoError(t, err)
	testHeaders(t, spec, res.Response)

	assert.Len(t, expected, 1)
	assert.Equal(t, expected[0].Activeslotcoeff, res.Data.Activeslotcoeff, "wrong Activeslotcoeff")
	assert.Equal(t, expected[0].Alonzogenesis, res.Data.Alonzogenesis, "wrong Alonzogenesis")
	assert.Equal(t, expected[0].Epochlength, res.Data.Epochlength, "wrong Epochlength")
	assert.Equal(t, expected[0].Maxkesrevolutions, res.Data.Maxkesrevolutions, "wrong Maxkesrevolutions")
	assert.Equal(t, expected[0].Maxlovelacesupply, res.Data.Maxlovelacesupply, "wrong Maxlovelacesupply")
	assert.Equal(t, expected[0].Networkid, res.Data.Networkid, "wrong Networkid")
	assert.Equal(t, expected[0].Networkmagic, res.Data.Networkmagic, "wrong Networkmagic")
	assert.Equal(t, expected[0].Securityparam, res.Data.Securityparam, "wrong Securityparam")
	assert.Equal(t, expected[0].Slotlength, res.Data.Slotlength, "wrong Slotlength")
	assert.Equal(t, expected[0].Slotsperkesperiod, res.Data.Slotsperkesperiod, "wrong Slotsperkesperiod")
	assert.Equal(t, expected[0].Systemstart, res.Data.Systemstart, "wrong Systemstart")
	assert.Equal(t, expected[0].Updatequorum, res.Data.Updatequorum, "wrong Updatequorum")
}
func TestNetworkTotals(t *testing.T) {
	expected := []koios.Totals{}

	spec := loadEndpointTestSpec(t, "endpoint_network_totals.json.gz", &expected)

	ts, api := createTestServerAndClient(t, spec)

	defer ts.Close()

	epochNo, err := strconv.ParseUint(spec.Request.Query.Get("_epoch_no"), 10, 64)
	assert.NoError(t, err)
	epoch := koios.EpochNo(epochNo)
	res, err := api.GetTotals(context.TODO(), &epoch)

	assert.NoError(t, err)
	testHeaders(t, spec, res.Response)

	assert.Len(t, expected, 1)
	assert.Equal(t, expected[0].Circulation, res.Data[0].Circulation, "wrong Circulation")
	assert.Equal(t, expected[0].Epoch, res.Data[0].Epoch, "wrong Epoch")
	assert.Equal(t, epoch, res.Data[0].Epoch, "wrong Epoch")
	assert.Equal(t, expected[0].Reserves, res.Data[0].Reserves, "wrong Reserves")
	assert.Equal(t, expected[0].Reward, res.Data[0].Reward, "wrong Reward")
	assert.Equal(t, expected[0].Supply, res.Data[0].Supply, "wrong Supply")
	assert.Equal(t, expected[0].Treasury, res.Data[0].Treasury, "wrong Treasury")
}
