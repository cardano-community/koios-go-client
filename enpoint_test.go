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

func TestNetworkTipEndpoint(t *testing.T) {
	expected := []koios.Tip{}

	spec := loadEndpointTestSpec(t, "endpoint_network_tip.json.gz", &expected)

	ts, api := createTestServerAndClient(t, spec)

	defer ts.Close()

	res, err := api.GetTip(context.TODO())

	assert.NoError(t, err)
	testHeaders(t, spec, res.Response)

	assert.Len(t, expected, 1)
	assert.Equal(t, expected[0], *res.Data)
}

func TestNetworkGenesiEndpoint(t *testing.T) {
	expected := []koios.Genesis{}

	spec := loadEndpointTestSpec(t, "endpoint_network_genesis.json.gz", &expected)

	ts, api := createTestServerAndClient(t, spec)

	defer ts.Close()

	res, err := api.GetGenesis(context.TODO())

	assert.NoError(t, err)
	testHeaders(t, spec, res.Response)

	assert.Len(t, expected, 1)
	assert.Equal(t, expected[0], *res.Data)
}

func TestNetworkTotalsEndpoint(t *testing.T) {
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
	assert.Equal(t, expected[0], res.Data[0])
}

func TestEpochInfoEndpoint(t *testing.T) {
	expected := []koios.EpochInfo{}

	spec := loadEndpointTestSpec(t, "endpoint_epoch_info.json.gz", &expected)

	ts, api := createTestServerAndClient(t, spec)

	defer ts.Close()

	epochNo, err := strconv.ParseUint(spec.Request.Query.Get("_epoch_no"), 10, 64)
	assert.NoError(t, err)
	epoch := koios.EpochNo(epochNo)

	res, err := api.GetEpochInfo(context.TODO(), &epoch)

	assert.NoError(t, err)
	testHeaders(t, spec, res.Response)

	assert.Len(t, expected, 1)
	assert.Equal(t, expected[0], res.Data[0])
}

func TestEpochParamsEndpoint(t *testing.T) {
	expected := []koios.EpochParams{}

	spec := loadEndpointTestSpec(t, "endpoint_epoch_params.json.gz", &expected)

	ts, api := createTestServerAndClient(t, spec)

	defer ts.Close()

	epochNo, err := strconv.ParseUint(spec.Request.Query.Get("_epoch_no"), 10, 64)
	assert.NoError(t, err)
	epoch := koios.EpochNo(epochNo)

	res, err := api.GetEpochParams(context.TODO(), &epoch)

	assert.NoError(t, err)
	testHeaders(t, spec, res.Response)

	assert.Len(t, expected, 1)
	assert.Equal(t, expected[0], res.Data[0])
}

func TestAccountListEndpoint(t *testing.T) {
	expected := []struct {
		StakeAddress koios.StakeAddress `json:"id"`
	}{}

	spec := loadEndpointTestSpec(t, "endpoint_account_list.json.gz", &expected)

	ts, api := createTestServerAndClient(t, spec)

	defer ts.Close()

	res, err := api.GetAccountList(context.TODO())
	assert.NoError(t, err)

	testHeaders(t, spec, res.Response)

	for _, e := range expected {
		assert.Contains(t, res.Data, e.StakeAddress)
	}
}

func TestAccountInfoEndpoint(t *testing.T) {
	expected := []koios.AccountInfo{}

	spec := loadEndpointTestSpec(t, "endpoint_account_info.json.gz", &expected)

	ts, api := createTestServerAndClient(t, spec)

	defer ts.Close()

	res, err := api.GetAccountInfo(context.TODO(), koios.Address(spec.Request.Query.Get("_address")))

	assert.NoError(t, err)
	testHeaders(t, spec, res.Response)

	assert.Len(t, expected, 1)
	assert.Equal(t, &expected[0], res.Data)
}

func TestAccountRewardsEndpoint(t *testing.T) {
	expected := []koios.AccountRewards{}

	spec := loadEndpointTestSpec(t, "endpoint_account_rewards.json.gz", &expected)

	ts, api := createTestServerAndClient(t, spec)

	defer ts.Close()

	epochNo, err := strconv.ParseUint(spec.Request.Query.Get("_epoch_no"), 10, 64)
	assert.NoError(t, err)
	epoch := koios.EpochNo(epochNo)

	res, err := api.GetAccountRewards(context.TODO(), koios.StakeAddress(spec.Request.Query.Get("_address")), &epoch)

	assert.NoError(t, err)
	testHeaders(t, spec, res.Response)

	assert.Len(t, expected, 1)
	assert.Equal(t, expected, res.Data)
}

func TestAccountUpdatesEndpoint(t *testing.T) {
	expected := []koios.AccountAction{}

	spec := loadEndpointTestSpec(t, "endpoint_account_updates.json.gz", &expected)

	ts, api := createTestServerAndClient(t, spec)

	defer ts.Close()

	res, err := api.GetAccountUpdates(context.TODO(), koios.StakeAddress(spec.Request.Query.Get("_address")))

	assert.NoError(t, err)
	testHeaders(t, spec, res.Response)

	assert.Equal(t, expected, res.Data)
}

func TestAccountAddressesEndpoint(t *testing.T) {
	expected := []struct {
		Address koios.Address `json:"address"`
	}{}

	spec := loadEndpointTestSpec(t, "endpoint_account_addresses.json.gz", &expected)

	ts, api := createTestServerAndClient(t, spec)

	defer ts.Close()

	res, err := api.GetAccountAddresses(context.TODO(), koios.StakeAddress(spec.Request.Query.Get("_address")))

	assert.NoError(t, err)
	testHeaders(t, spec, res.Response)

	for _, e := range expected {
		assert.Contains(t, res.Data, e.Address)
	}
}
func TestAccountAssetsEndpoint(t *testing.T) {
	expected := []koios.AccountAsset{}

	spec := loadEndpointTestSpec(t, "endpoint_account_assets.json.gz", &expected)

	ts, api := createTestServerAndClient(t, spec)

	defer ts.Close()

	res, err := api.GetAccountAssets(context.TODO(), koios.StakeAddress(spec.Request.Query.Get("_address")))

	assert.NoError(t, err)
	testHeaders(t, spec, res.Response)

	assert.Equal(t, expected, res.Data)
}

func TestAccountHistoryEndpoint(t *testing.T) {
	expected := []koios.AccountHistoryEntry{}

	spec := loadEndpointTestSpec(t, "endpoint_account_history.json.gz", &expected)

	ts, api := createTestServerAndClient(t, spec)

	defer ts.Close()

	res, err := api.GetAccountHistory(context.TODO(), koios.StakeAddress(spec.Request.Query.Get("_address")))

	assert.NoError(t, err)
	testHeaders(t, spec, res.Response)

	assert.Equal(t, expected, res.Data)
}
