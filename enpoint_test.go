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
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/cardano-community/koios-go-client"
	"github.com/cardano-community/koios-go-client/internal"
)

func TestNetworkTipEndpoint(t *testing.T) {
	expected := []koios.Tip{}

	spec := loadEndpointTestSpec(t, "endpoint_network_tip.json.gz", &expected)

	ts, api := setupTestServerAndClient(t, spec)

	defer ts.Close()

	res, err := api.GetTip(context.Background(), nil)

	assert.NoError(t, err)
	testHeaders(t, spec, res.Response)

	assert.Len(t, expected, 1)
	assert.Equal(t, &expected[0], res.Data)

	c, err := api.WithOptions(koios.Host("127.0.0.2:80"))
	assert.NoError(t, err)
	_, err = c.GetTip(context.Background(), nil)
	assert.EqualError(t, err, "dial tcp: lookup 127.0.0.2:80: no such host")
}

func TestRequestContext(t *testing.T) {
	expected := []koios.Tip{}
	spec := loadEndpointTestSpec(t, "endpoint_network_tip.json.gz", &expected)

	ts, api := setupTestServerAndClient(t, spec)

	defer ts.Close()

	res, err := api.GetTip(nil, nil) //nolint: staticcheck
	assert.EqualError(t, err, "net/http: nil Context")
	assert.Equal(t, "net/http: nil Context", res.Error.Message)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(0))
	defer cancel()

	res2, err := api.GetTip(ctx, nil)

	assert.EqualError(t, err, "context deadline exceeded")
	assert.Nil(t, res2.Data)

	ctx2, cancel2 := context.WithTimeout(context.Background(), time.Microsecond*201)
	defer cancel2()

	res3, err3 := api.GetTip(ctx2, nil)

	var edgeerr bool
	if err3.Error() == "context deadline exceeded" || strings.Contains(err3.Error(), "i/o timeout") {
		edgeerr = true
	}
	assert.True(t, edgeerr, "expected: context deadline exceeded or i/o timeout")
	assert.Nil(t, res3.Data)
}

func TestNetworkGenesiEndpoint(t *testing.T) {
	expected := []koios.Genesis{}

	spec := loadEndpointTestSpec(t, "endpoint_network_genesis.json.gz", &expected)

	ts, api := setupTestServerAndClient(t, spec)

	defer ts.Close()

	res, err := api.GetGenesis(context.Background(), nil)

	assert.NoError(t, err)
	testHeaders(t, spec, res.Response)

	assert.Len(t, expected, 1)
	assert.Equal(t, &expected[0], res.Data)

	c, err := api.WithOptions(koios.Host("127.0.0.2:80"))
	assert.NoError(t, err)
	_, err = c.GetGenesis(context.Background(), nil)
	assert.EqualError(t, err, "dial tcp: lookup 127.0.0.2:80: no such host")
}

func Test404s(t *testing.T) {
	// test invalid path
	tipspec := loadEndpointTestSpec(t, "endpoint_network_tip.json.gz", nil)
	ts, api := setupTestServerAndClient(t, tipspec)
	defer ts.Close()

	res, err := api.GetGenesis(context.Background(), nil)
	assert.Error(t, err)
	assert.Nil(t, res.Data)
	assert.Equal(t, "got unexpected response: 404 Not Found", res.Error.Message)
	assert.Equal(t, http.StatusNotFound, res.StatusCode)

	// errors with stats should be same
	res2, err := api.GetGenesis(context.Background(), nil)
	assert.Error(t, err)
	assert.Nil(t, res2.Data)
	assert.Equal(t, "got unexpected response: 404 Not Found", res2.Error.Message)
	assert.Equal(t, http.StatusNotFound, res2.StatusCode)
}

func TestNetworkTotalsEndpoint(t *testing.T) {
	expected := []koios.Totals{}

	spec := loadEndpointTestSpec(t, "endpoint_network_totals.json.gz", &expected)

	ts, api := setupTestServerAndClient(t, spec)

	defer ts.Close()

	epochNo, err := strconv.ParseUint(spec.Request.Query.Get("_epoch_no"), 10, 64)
	assert.NoError(t, err)
	epoch := koios.EpochNo(epochNo)

	res, err := api.GetTotals(context.Background(), &epoch, nil)
	assert.NoError(t, err)
	testHeaders(t, spec, res.Response)
	assert.Len(t, expected, 1)
	assert.Equal(t, expected[0], res.Data[0])

	// test data without epoch
	res2, err := api.GetTotals(context.Background(), nil, nil)
	assert.NoError(t, err)
	testHeaders(t, spec, res2.Response)
	assert.Len(t, expected, 1)
	assert.Equal(t, expected[0], res2.Data[0])

	c, err := api.WithOptions(koios.Host("127.0.0.2:80"))
	assert.NoError(t, err)
	_, err = c.GetTotals(context.Background(), nil, nil)
	assert.EqualError(t, err, "dial tcp: lookup 127.0.0.2:80: no such host")
}

func TestEpochInfoEndpoint(t *testing.T) {
	expected := []koios.EpochInfo{}

	spec := loadEndpointTestSpec(t, "endpoint_epoch_info.json.gz", &expected)

	ts, api := setupTestServerAndClient(t, spec)

	defer ts.Close()

	epochNo, err := strconv.ParseUint(spec.Request.Query.Get("_epoch_no"), 10, 64)
	assert.NoError(t, err)
	epoch := koios.EpochNo(epochNo)

	res, err := api.GetEpochInfo(context.Background(), &epoch, nil)

	assert.NoError(t, err)
	testHeaders(t, spec, res.Response)

	assert.Len(t, expected, 1)
	assert.Equal(t, expected[0], res.Data[0])

	c, err := api.WithOptions(koios.Host("127.0.0.2:80"))
	assert.NoError(t, err)
	_, err = c.GetEpochInfo(context.Background(), &epoch, nil)
	assert.EqualError(t, err, "dial tcp: lookup 127.0.0.2:80: no such host")
}

func TestEpochParamsEndpoint(t *testing.T) {
	expected := []koios.EpochParams{}

	spec := loadEndpointTestSpec(t, "endpoint_epoch_params.json.gz", &expected)

	ts, api := setupTestServerAndClient(t, spec)

	defer ts.Close()

	epochNo, err := strconv.ParseUint(spec.Request.Query.Get("_epoch_no"), 10, 64)
	assert.NoError(t, err)
	epoch := koios.EpochNo(epochNo)

	res, err := api.GetEpochParams(context.Background(), &epoch, nil)

	assert.NoError(t, err)
	testHeaders(t, spec, res.Response)

	assert.Len(t, expected, 1)
	assert.Equal(t, expected[0], res.Data[0])

	c, err := api.WithOptions(koios.Host("127.0.0.2:80"))
	assert.NoError(t, err)
	_, err = c.GetEpochParams(context.Background(), &epoch, nil)
	assert.EqualError(t, err, "dial tcp: lookup 127.0.0.2:80: no such host")
}

func TestAccountListEndpoint(t *testing.T) {
	expected := []struct {
		StakeAddress koios.StakeAddress `json:"id"`
	}{}

	spec := loadEndpointTestSpec(t, "endpoint_account_list.json.gz", &expected)

	ts, api := setupTestServerAndClient(t, spec)

	defer ts.Close()

	res, err := api.GetAccountList(context.Background(), nil)
	assert.NoError(t, err)

	testHeaders(t, spec, res.Response)

	for _, e := range expected {
		assert.Contains(t, res.Data, e.StakeAddress)
	}

	c, err := api.WithOptions(koios.Host("127.0.0.2:80"))
	assert.NoError(t, err)
	_, err = c.GetAccountList(context.Background(), nil)
	assert.EqualError(t, err, "dial tcp: lookup 127.0.0.2:80: no such host")
}

func TestAccountInfoEndpoint(t *testing.T) {
	expected := []koios.AccountInfo{}

	spec := loadEndpointTestSpec(t, "endpoint_account_info.json.gz", &expected)

	ts, api := setupTestServerAndClient(t, spec)

	defer ts.Close()

	res, err := api.GetAccountInfo(context.Background(), koios.Address(spec.Request.Query.Get("_address")), nil)

	assert.NoError(t, err)
	testHeaders(t, spec, res.Response)

	assert.Len(t, expected, 1)
	assert.Equal(t, &expected[0], res.Data)

	res2, err := api.GetAccountInfo(context.Background(), koios.Address(""), nil)
	assert.ErrorIs(t, err, koios.ErrNoAddress)
	assert.Nil(t, res2.Data, "response data should be nil if arg is invalid")
	assert.Equal(t, res2.Error.Message, "missing address")

	c, err := api.WithOptions(koios.Host("127.0.0.2:80"))
	assert.NoError(t, err)
	_, err = c.GetAccountInfo(context.Background(), koios.Address(spec.Request.Query.Get("_address")), nil)
	assert.EqualError(t, err, "dial tcp: lookup 127.0.0.2:80: no such host")
}

func TestAccountRewardsEndpoint(t *testing.T) {
	expected := []koios.AccountRewards{}

	spec := loadEndpointTestSpec(t, "endpoint_account_rewards.json.gz", &expected)

	ts, api := setupTestServerAndClient(t, spec)

	defer ts.Close()

	epochNo, err := strconv.ParseUint(spec.Request.Query.Get("_epoch_no"), 10, 64)
	assert.NoError(t, err)
	epoch := koios.EpochNo(epochNo)

	res, err := api.GetAccountRewards(
		context.Background(),
		koios.StakeAddress(spec.Request.Query.Get("_address")),
		&epoch,
		nil,
	)

	assert.NoError(t, err)
	testHeaders(t, spec, res.Response)

	assert.Len(t, expected, 1)
	assert.Equal(t, expected, res.Data)

	c, err := api.WithOptions(koios.Host("127.0.0.2:80"))
	assert.NoError(t, err)
	_, err = c.GetAccountRewards(context.Background(), koios.StakeAddress(spec.Request.Query.Get("_address")), &epoch, nil)
	assert.EqualError(t, err, "dial tcp: lookup 127.0.0.2:80: no such host")
}

func TestAccountUpdatesEndpoint(t *testing.T) {
	expected := []koios.AccountAction{}

	spec := loadEndpointTestSpec(t, "endpoint_account_updates.json.gz", &expected)

	ts, api := setupTestServerAndClient(t, spec)

	defer ts.Close()

	res, err := api.GetAccountUpdates(context.Background(), koios.StakeAddress(spec.Request.Query.Get("_address")), nil)

	assert.NoError(t, err)
	testHeaders(t, spec, res.Response)

	assert.Equal(t, expected, res.Data)

	c, err := api.WithOptions(koios.Host("127.0.0.2:80"))
	assert.NoError(t, err)
	_, err = c.GetAccountUpdates(context.Background(), koios.StakeAddress(spec.Request.Query.Get("_address")), nil)
	assert.EqualError(t, err, "dial tcp: lookup 127.0.0.2:80: no such host")
}

func TestAccountAddressesEndpoint(t *testing.T) {
	expected := []struct {
		Address koios.Address `json:"address"`
	}{}

	spec := loadEndpointTestSpec(t, "endpoint_account_addresses.json.gz", &expected)

	ts, api := setupTestServerAndClient(t, spec)

	defer ts.Close()

	res, err := api.GetAccountAddresses(context.Background(), koios.StakeAddress(spec.Request.Query.Get("_address")), nil)

	assert.NoError(t, err)
	testHeaders(t, spec, res.Response)

	for _, e := range expected {
		assert.Contains(t, res.Data, e.Address)
	}

	c, err := api.WithOptions(koios.Host("127.0.0.2:80"))
	assert.NoError(t, err)
	_, err = c.GetAccountAddresses(context.Background(), koios.StakeAddress(spec.Request.Query.Get("_address")), nil)
	assert.EqualError(t, err, "dial tcp: lookup 127.0.0.2:80: no such host")
}
func TestAccountAssetsEndpoint(t *testing.T) {
	expected := []koios.AccountAsset{}

	spec := loadEndpointTestSpec(t, "endpoint_account_assets.json.gz", &expected)

	ts, api := setupTestServerAndClient(t, spec)

	defer ts.Close()

	res, err := api.GetAccountAssets(context.Background(), koios.StakeAddress(spec.Request.Query.Get("_address")), nil)

	assert.NoError(t, err)
	testHeaders(t, spec, res.Response)

	assert.Equal(t, expected, res.Data)

	c, err := api.WithOptions(koios.Host("127.0.0.2:80"))
	assert.NoError(t, err)
	_, err = c.GetAccountAssets(context.Background(), koios.StakeAddress(spec.Request.Query.Get("_address")), nil)
	assert.EqualError(t, err, "dial tcp: lookup 127.0.0.2:80: no such host")
}

func TestAccountHistoryEndpoint(t *testing.T) {
	expected := []koios.AccountHistoryEntry{}

	spec := loadEndpointTestSpec(t, "endpoint_account_history.json.gz", &expected)

	ts, api := setupTestServerAndClient(t, spec)

	defer ts.Close()

	res, err := api.GetAccountHistory(context.Background(), koios.StakeAddress(spec.Request.Query.Get("_address")), nil)

	assert.NoError(t, err)
	testHeaders(t, spec, res.Response)

	assert.Equal(t, expected, res.Data)

	c, err := api.WithOptions(koios.Host("127.0.0.2:80"))
	assert.NoError(t, err)
	_, err = c.GetAccountHistory(context.Background(), koios.StakeAddress(spec.Request.Query.Get("_address")), nil)
	assert.EqualError(t, err, "dial tcp: lookup 127.0.0.2:80: no such host")
}

func TestGetAddressInfoEndpoint(t *testing.T) {
	expected := []koios.AddressInfo{}

	spec := loadEndpointTestSpec(t, "endpoint_address_info.json.gz", &expected)

	ts, api := setupTestServerAndClient(t, spec)

	defer ts.Close()

	res, err := api.GetAddressInfo(context.Background(), koios.Address(spec.Request.Query.Get("_address")), nil)

	assert.NoError(t, err)
	testHeaders(t, spec, res.Response)

	assert.Equal(t, &expected[0], res.Data)

	res2, err := api.GetAddressInfo(context.Background(), koios.Address(""), nil)
	assert.ErrorIs(t, err, koios.ErrNoAddress)
	assert.Nil(t, res2.Data, "response data should be nil if arg is invalid")
	assert.Equal(t, res2.Error.Message, "missing address")

	c, err := api.WithOptions(koios.Host("127.0.0.2:80"))
	assert.NoError(t, err)
	_, err = c.GetAddressInfo(context.Background(), koios.Address(spec.Request.Query.Get("_address")), nil)
	assert.EqualError(t, err, "dial tcp: lookup 127.0.0.2:80: no such host")
}

func TestGetAddressTxsEndpoint(t *testing.T) {
	expected := []struct {
		TxHash koios.TxHash `json:"tx_hash"`
	}{}

	spec := loadEndpointTestSpec(t, "endpoint_address_txs.json.gz", &expected)

	ts, api := setupTestServerAndClient(t, spec)

	defer ts.Close()

	var payload = struct {
		Adresses         []koios.Address `json:"_addresses"`
		AfterBlockHeight uint64          `json:"_after_block_height,omitempty"`
	}{}
	err := json.Unmarshal(spec.Request.Body, &payload)
	assert.NoError(t, err)

	res, err := api.GetAddressTxs(context.Background(), payload.Adresses, payload.AfterBlockHeight, nil)

	assert.NoError(t, err)
	testHeaders(t, spec, res.Response)

	for _, e := range expected {
		assert.Contains(t, res.Data, e.TxHash)
	}

	res2, err := api.GetAddressTxs(context.Background(), []koios.Address{}, 0, nil)
	assert.ErrorIs(t, err, koios.ErrNoAddress)
	assert.Nil(t, res2.Data, "response data should be nil if arg is invalid")
	assert.Equal(t, res2.Error.Message, "missing address")

	c, err := api.WithOptions(koios.Host("127.0.0.2:80"))
	assert.NoError(t, err)
	_, err = c.GetAddressTxs(
		context.Background(),
		[]koios.Address{koios.Address(spec.Request.Query.Get("_address"))},
		0,
		nil,
	)
	assert.EqualError(t, err, "dial tcp: lookup 127.0.0.2:80: no such host")
}

func TestGetAddressAssetsEndpoint(t *testing.T) {
	expected := []koios.AddressAsset{}

	spec := loadEndpointTestSpec(t, "endpoint_address_assets.json.gz", &expected)

	ts, api := setupTestServerAndClient(t, spec)

	defer ts.Close()

	res, err := api.GetAddressAssets(context.Background(), koios.Address(spec.Request.Query.Get("_address")), nil)

	assert.NoError(t, err)
	testHeaders(t, spec, res.Response)

	assert.Equal(t, expected, res.Data)

	res2, err := api.GetAddressAssets(context.Background(), koios.Address(""), nil)
	assert.ErrorIs(t, err, koios.ErrNoAddress)
	assert.Nil(t, res2.Data, "response data should be nil if arg is invalid")
	assert.Equal(t, res2.Error.Message, "missing address")

	c, err := api.WithOptions(koios.Host("127.0.0.2:80"))
	assert.NoError(t, err)
	_, err = c.GetAddressAssets(context.Background(), koios.Address(spec.Request.Query.Get("_address")), nil)
	assert.EqualError(t, err, "dial tcp: lookup 127.0.0.2:80: no such host")
}

func TestGetCredentialTxsEndpoint(t *testing.T) {
	expected := []struct {
		TxHash koios.TxHash `json:"tx_hash"`
	}{}

	spec := loadEndpointTestSpec(t, "endpoint_credential_txs.json.gz", &expected)

	ts, api := setupTestServerAndClient(t, spec)

	defer ts.Close()

	var payload = struct {
		Credentials      []koios.PaymentCredential `json:"_payment_credentials"`
		AfterBlockHeight uint64                    `json:"_after_block_height,omitempty"`
	}{}
	err := json.Unmarshal(spec.Request.Body, &payload)
	assert.NoError(t, err)

	res, err := api.GetCredentialTxs(context.Background(), payload.Credentials, payload.AfterBlockHeight, nil)

	assert.NoError(t, err)
	testHeaders(t, spec, res.Response)

	for _, e := range expected {
		assert.Contains(t, res.Data, e.TxHash)
	}

	res2, err := api.GetCredentialTxs(context.Background(), []koios.PaymentCredential{}, 0, nil)
	assert.ErrorIs(t, err, koios.ErrNoAddress)
	assert.Nil(t, res2.Data, "response data should be nil if arg is invalid")
	assert.Equal(t, res2.Error.Message, "missing address")

	c, err := api.WithOptions(koios.Host("127.0.0.2:80"))
	assert.NoError(t, err)
	_, err = c.GetCredentialTxs(context.Background(), payload.Credentials, payload.AfterBlockHeight, nil)
	assert.EqualError(t, err, "dial tcp: lookup 127.0.0.2:80: no such host")
}

func TestGetAssetListEndpoint(t *testing.T) {
	expected := []koios.AssetListItem{}

	spec := loadEndpointTestSpec(t, "endpoint_asset_list.json.gz", &expected)

	ts, api := setupTestServerAndClient(t, spec)

	defer ts.Close()

	res, err := api.GetAssetList(context.Background(), nil)

	assert.NoError(t, err)
	testHeaders(t, spec, res.Response)

	assert.Equal(t, expected, res.Data)

	c, err := api.WithOptions(koios.Host("127.0.0.2:80"))
	assert.NoError(t, err)
	_, err = c.GetAssetList(context.Background(), nil)
	assert.EqualError(t, err, "dial tcp: lookup 127.0.0.2:80: no such host")
}

func TestGetAssetAddressListEndpoint(t *testing.T) {
	expected := []koios.AssetHolder{}

	spec := loadEndpointTestSpec(t, "endpoint_asset_address_list.json.gz", &expected)

	ts, api := setupTestServerAndClient(t, spec)

	defer ts.Close()

	res, err := api.GetAssetAddressList(
		context.Background(),
		koios.PolicyID(spec.Request.Query.Get("_asset_policy")),
		koios.AssetName(spec.Request.Query.Get("_asset_name")),
		nil,
	)

	assert.NoError(t, err)
	testHeaders(t, spec, res.Response)

	assert.Equal(t, expected, res.Data)

	c, err := api.WithOptions(koios.Host("127.0.0.2:80"))
	assert.NoError(t, err)
	_, err = c.GetAssetAddressList(
		context.Background(),
		koios.PolicyID(spec.Request.Query.Get("_asset_policy")),
		koios.AssetName(spec.Request.Query.Get("_asset_name")),
		nil,
	)
	assert.EqualError(t, err, "dial tcp: lookup 127.0.0.2:80: no such host")
}

func TestGetAssetInfoEndpoint(t *testing.T) {
	expected := []koios.AssetInfo{}

	spec := loadEndpointTestSpec(t, "endpoint_asset_info.json.gz", &expected)

	ts, api := setupTestServerAndClient(t, spec)

	defer ts.Close()

	res, err := api.GetAssetInfo(
		context.Background(),
		koios.PolicyID(spec.Request.Query.Get("_asset_policy")),
		koios.AssetName(spec.Request.Query.Get("_asset_name")),
		nil,
	)

	assert.NoError(t, err)
	testHeaders(t, spec, res.Response)

	assert.Equal(t, &expected[0], res.Data)

	c, err := api.WithOptions(koios.Host("127.0.0.2:80"))
	assert.NoError(t, err)
	_, err = c.GetAssetInfo(
		context.Background(),
		koios.PolicyID(spec.Request.Query.Get("_asset_policy")),
		koios.AssetName(spec.Request.Query.Get("_asset_name")),
		nil,
	)
	assert.EqualError(t, err, "dial tcp: lookup 127.0.0.2:80: no such host")
}

func TestGetAssetPolicyInfoEndpoint(t *testing.T) {
	expected := []koios.AssetPolicyInfo{}

	spec := loadEndpointTestSpec(t, "endpoint_asset_policy_info.json.gz", &expected)

	ts, api := setupTestServerAndClient(t, spec)

	defer ts.Close()

	res, err := api.GetAssetPolicyInfo(
		context.Background(),
		koios.PolicyID(spec.Request.Query.Get("_asset_policy")),
		nil,
	)

	assert.NoError(t, err)
	testHeaders(t, spec, res.Response)

	assert.Equal(t, &expected[0], res.Data)

	c, err := api.WithOptions(koios.Host("127.0.0.2:80"))
	assert.NoError(t, err)
	_, err = c.GetAssetPolicyInfo(
		context.Background(),
		koios.PolicyID(spec.Request.Query.Get("_asset_policy")),
		nil,
	)
	assert.EqualError(t, err, "dial tcp: lookup 127.0.0.2:80: no such host")
}

func TestGetAssetSummaryEndpoint(t *testing.T) {
	expected := []koios.AssetSummary{}

	spec := loadEndpointTestSpec(t, "endpoint_asset_summary.json.gz", &expected)

	ts, api := setupTestServerAndClient(t, spec)

	defer ts.Close()

	res, err := api.GetAssetSummary(
		context.Background(),
		koios.PolicyID(spec.Request.Query.Get("_asset_policy")),
		koios.AssetName(spec.Request.Query.Get("_asset_name")),
		nil,
	)

	assert.NoError(t, err)
	testHeaders(t, spec, res.Response)

	assert.Equal(t, &expected[0], res.Data)

	c, err := api.WithOptions(koios.Host("127.0.0.2:80"))
	assert.NoError(t, err)
	_, err = c.GetAssetSummary(
		context.Background(),
		koios.PolicyID(spec.Request.Query.Get("_asset_policy")),
		koios.AssetName(spec.Request.Query.Get("_asset_name")),
		nil,
	)
	assert.EqualError(t, err, "dial tcp: lookup 127.0.0.2:80: no such host")
}

func TestGetAssetTxsEndpoint(t *testing.T) {
	expected := []koios.AssetTxs{}

	spec := loadEndpointTestSpec(t, "endpoint_asset_txs.json.gz", &expected)

	ts, api := setupTestServerAndClient(t, spec)

	defer ts.Close()

	res, err := api.GetAssetTxs(
		context.Background(),
		koios.PolicyID(spec.Request.Query.Get("_asset_policy")),
		koios.AssetName(spec.Request.Query.Get("_asset_name")),
		nil,
	)

	assert.NoError(t, err)
	testHeaders(t, spec, res.Response)

	assert.Equal(t, &expected[0], res.Data)

	c, err := api.WithOptions(koios.Host("127.0.0.2:80"))
	assert.NoError(t, err)
	_, err = c.GetAssetTxs(
		context.Background(),
		koios.PolicyID(spec.Request.Query.Get("_asset_policy")),
		koios.AssetName(spec.Request.Query.Get("_asset_name")),
		nil,
	)
	assert.EqualError(t, err, "dial tcp: lookup 127.0.0.2:80: no such host")
}

func TestGetAssetHistoryEndpoint(t *testing.T) {
	expected := []koios.AssetHistory{}

	spec := loadEndpointTestSpec(t, "endpoint_asset_history.json.gz", &expected)

	ts, api := setupTestServerAndClient(t, spec)

	defer ts.Close()

	res, err := api.GetAssetHistory(
		context.Background(),
		koios.PolicyID(spec.Request.Query.Get("_asset_policy")),
		koios.AssetName(spec.Request.Query.Get("_asset_name")),
		nil,
	)

	assert.NoError(t, err)
	testHeaders(t, spec, res.Response)

	assert.Equal(t, &expected[0], res.Data)

	c, err := api.WithOptions(koios.Host("127.0.0.2:80"))
	assert.NoError(t, err)
	_, err = c.GetAssetHistory(
		context.Background(),
		koios.PolicyID(spec.Request.Query.Get("_asset_policy")),
		koios.AssetName(spec.Request.Query.Get("_asset_name")),
		nil,
	)
	assert.EqualError(t, err, "dial tcp: lookup 127.0.0.2:80: no such host")
}

func TestGetBlockInfoEndpoint(t *testing.T) {
	expected := []koios.Block{}

	spec := loadEndpointTestSpec(t, "endpoint_block_info.json.gz", &expected)

	ts, api := setupTestServerAndClient(t, spec)

	defer ts.Close()

	var payload = struct {
		BlockHashes []koios.BlockHash `json:"_block_hashes"`
	}{}

	err := json.Unmarshal(spec.Request.Body, &payload)
	assert.NoError(t, err)

	res, err := api.GetBlockInfo(
		context.Background(),
		payload.BlockHashes[0],
		nil,
	)

	assert.NoError(t, err)
	testHeaders(t, spec, res.Response)

	assert.Equal(t, &expected[0], res.Data)

	c, err := api.WithOptions(koios.Host("127.0.0.2:80"))
	assert.NoError(t, err)
	_, err = c.GetBlockInfo(
		context.Background(),
		payload.BlockHashes[0],
		nil,
	)
	assert.EqualError(t, err, "dial tcp: lookup 127.0.0.2:80: no such host")
}

func TestGetBlocksInfoEndpoint(t *testing.T) {
	expected := []koios.Block{}

	spec := loadEndpointTestSpec(t, "endpoint_block_info.json.gz", &expected)

	ts, api := setupTestServerAndClient(t, spec)

	defer ts.Close()

	var payload = struct {
		BlockHashes []koios.BlockHash `json:"_block_hashes"`
	}{}
	err := json.Unmarshal(spec.Request.Body, &payload)
	assert.NoError(t, err)

	res, err := api.GetBlocksInfo(
		context.Background(),
		payload.BlockHashes,
		nil,
	)

	assert.NoError(t, err)
	testHeaders(t, spec, res.Response)

	assert.Equal(t, expected, res.Data)

	c, err := api.WithOptions(koios.Host("127.0.0.2:80"))
	assert.NoError(t, err)
	_, err = c.GetBlocksInfo(
		context.Background(),
		payload.BlockHashes,
		nil,
	)
	assert.EqualError(t, err, "dial tcp: lookup 127.0.0.2:80: no such host")
}

func TestGetBlockTxsEndpoint(t *testing.T) {
	expected := []struct {
		TxHash koios.TxHash `json:"tx_hash"`
	}{}

	spec := loadEndpointTestSpec(t, "endpoint_block_txs.json.gz", &expected)

	ts, api := setupTestServerAndClient(t, spec)

	defer ts.Close()

	res, err := api.GetBlockTxHashes(
		context.Background(),
		koios.BlockHash(spec.Request.Query.Get("_block_hash")),
		nil,
	)

	assert.NoError(t, err)
	testHeaders(t, spec, res.Response)

	for _, e := range expected {
		assert.Contains(t, res.Data, e.TxHash)
	}

	c, err := api.WithOptions(koios.Host("127.0.0.2:80"))
	assert.NoError(t, err)
	_, err = c.GetBlockTxHashes(
		context.Background(),
		koios.BlockHash(spec.Request.Query.Get("_block_hash")),
		nil,
	)
	assert.EqualError(t, err, "dial tcp: lookup 127.0.0.2:80: no such host")
}

func TestGetBlocksEndpoint(t *testing.T) {
	expected := []koios.Block{}

	spec := loadEndpointTestSpec(t, "endpoint_blocks.json.gz", &expected)

	ts, api := setupTestServerAndClient(t, spec)

	defer ts.Close()

	res, err := api.GetBlocks(context.Background(), nil)

	assert.NoError(t, err)
	testHeaders(t, spec, res.Response)

	assert.Equal(t, expected, res.Data)

	c, err := api.WithOptions(koios.Host("127.0.0.2:80"))
	assert.NoError(t, err)
	_, err = c.GetBlocks(context.Background(), nil)
	assert.EqualError(t, err, "dial tcp: lookup 127.0.0.2:80: no such host")
}

func TestGetPoolBlocksEndpoint(t *testing.T) {
	expected := []koios.PoolBlockInfo{}

	spec := loadEndpointTestSpec(t, "endpoint_pool_blocks.json.gz", &expected)

	ts, api := setupTestServerAndClient(t, spec)

	defer ts.Close()

	epochNo, err := strconv.ParseUint(spec.Request.Query.Get("_epoch_no"), 10, 64)
	assert.NoError(t, err)
	epoch := koios.EpochNo(epochNo)

	res, err := api.GetPoolBlocks(
		context.Background(),
		koios.PoolID(spec.Request.Query.Get("_pool_bech32")),
		&epoch,
		nil,
	)

	assert.NoError(t, err)
	testHeaders(t, spec, res.Response)

	assert.Equal(t, expected, res.Data)

	c, err := api.WithOptions(koios.Host("127.0.0.2:80"))
	assert.NoError(t, err)
	_, err = c.GetPoolBlocks(
		context.Background(),
		koios.PoolID(spec.Request.Query.Get("_pool_bech32")),
		&epoch,
		nil,
	)
	assert.EqualError(t, err, "dial tcp: lookup 127.0.0.2:80: no such host")
}

func TestGetPoolDelegatorsEndpoint(t *testing.T) {
	expected := []koios.PoolDelegator{}

	spec := loadEndpointTestSpec(t, "endpoint_pool_delegators.json.gz", &expected)

	ts, api := setupTestServerAndClient(t, spec)

	defer ts.Close()

	epochNo, err := strconv.ParseUint(spec.Request.Query.Get("_epoch_no"), 10, 64)
	assert.NoError(t, err)
	epoch := koios.EpochNo(epochNo)

	res, err := api.GetPoolDelegators(
		context.Background(),
		koios.PoolID(spec.Request.Query.Get("_pool_bech32")),
		&epoch,
		nil,
	)

	assert.NoError(t, err)
	testHeaders(t, spec, res.Response)

	assert.Equal(t, expected, res.Data)

	c, err := api.WithOptions(koios.Host("127.0.0.2:80"))
	assert.NoError(t, err)
	_, err = c.GetPoolDelegators(
		context.Background(),
		koios.PoolID(spec.Request.Query.Get("_pool_bech32")),
		&epoch,
		nil,
	)
	assert.EqualError(t, err, "dial tcp: lookup 127.0.0.2:80: no such host")
}

func TestGetPoolHistoryEndpoint(t *testing.T) {
	expected := []koios.PoolHistory{}

	spec := loadEndpointTestSpec(t, "endpoint_pool_history.json.gz", &expected)

	ts, api := setupTestServerAndClient(t, spec)

	defer ts.Close()

	epochNo, err := strconv.ParseUint(spec.Request.Query.Get("_epoch_no"), 10, 64)
	assert.NoError(t, err)
	epoch := koios.EpochNo(epochNo)

	res, err := api.GetPoolHistory(
		context.Background(),
		koios.PoolID(spec.Request.Query.Get("_pool_bech32")),
		&epoch,
		nil,
	)

	assert.NoError(t, err)
	testHeaders(t, spec, res.Response)

	assert.Equal(t, expected, res.Data)

	c, err := api.WithOptions(koios.Host("127.0.0.2:80"))
	assert.NoError(t, err)
	_, err = c.GetPoolHistory(
		context.Background(),
		koios.PoolID(spec.Request.Query.Get("_pool_bech32")),
		&epoch,
		nil,
	)
	assert.EqualError(t, err, "dial tcp: lookup 127.0.0.2:80: no such host")
}

func TestGetPoolInfoEndpoint(t *testing.T) {
	expected := []koios.PoolInfo{}

	spec := loadEndpointTestSpec(t, "endpoint_pool_info.json.gz", &expected)
	ts, api := setupTestServerAndClient(t, spec)
	defer ts.Close()

	var payload = struct {
		PoolIDs []koios.PoolID `json:"_pool_bech32_ids"`
	}{}
	err := json.Unmarshal(spec.Request.Body, &payload)
	assert.NoError(t, err)

	res, err := api.GetPoolInfo(
		context.Background(),
		payload.PoolIDs[0],
		nil,
	)

	assert.NoError(t, err)
	testHeaders(t, spec, res.Response)

	assert.Equal(t, &expected[0], res.Data)

	res2, err := api.GetPoolInfos(context.Background(), []koios.PoolID{}, nil)
	assert.ErrorIs(t, err, koios.ErrNoPoolID)
	assert.Nil(t, res2.Data, "response data should be nil if arg is invalid")
	assert.Equal(t, res2.Error.Message, "missing pool id")

	c, err := api.WithOptions(koios.Host("127.0.0.2:80"))
	assert.NoError(t, err)
	_, err = c.GetPoolInfo(
		context.Background(),
		payload.PoolIDs[0],
		nil,
	)
	assert.EqualError(t, err, "dial tcp: lookup 127.0.0.2:80: no such host")

	rpipe, w := io.Pipe()
	go func() {
		_ = json.NewEncoder(w).Encode(payload)
		defer w.Close()
	}()

	opts := api.NewRequestOptions()

	for k, vv := range spec.Request.Header {
		for _, v := range vv {
			opts.HeadersAdd(k, v)
		}
	}
	rsp3, err3 := api.POST(context.Background(), "/pool_info", rpipe, opts)
	assert.NoError(t, err3)

	defer func() { _ = rsp3.Body.Close() }()
	res3 := &koios.PoolInfosResponse{}
	assert.NoError(t, koios.ReadAndUnmarshalResponse(rsp3, &res3.Response, &res3.Data))
	assert.Equal(t, expected[0], res3.Data[0])
}

func TestGetPoolListEndpoint(t *testing.T) {
	expected := []koios.PoolListItem{}

	spec := loadEndpointTestSpec(t, "endpoint_pool_list.json.gz", &expected)

	ts, api := setupTestServerAndClient(t, spec)

	defer ts.Close()

	res, err := api.GetPoolList(context.Background(), nil)

	assert.NoError(t, err)
	testHeaders(t, spec, res.Response)

	assert.Equal(t, expected, res.Data)

	c, err := api.WithOptions(koios.Host("127.0.0.2:80"))
	assert.NoError(t, err)
	_, err = c.GetPoolList(context.Background(), nil)
	assert.EqualError(t, err, "dial tcp: lookup 127.0.0.2:80: no such host")
}

func TestGetPoolMetadataEndpoint(t *testing.T) {
	expected := []koios.PoolMetadata{}

	spec := loadEndpointTestSpec(t, "endpoint_pool_metadata.json.gz", &expected)

	ts, api := setupTestServerAndClient(t, spec)

	defer ts.Close()

	var payload = struct {
		PoolIDs []koios.PoolID `json:"_pool_bech32_ids"`
	}{}
	err := json.Unmarshal(spec.Request.Body, &payload)
	assert.NoError(t, err)

	res, err := api.GetPoolMetadata(context.Background(), payload.PoolIDs, nil)

	assert.NoError(t, err)
	testHeaders(t, spec, res.Response)

	assert.Equal(t, expected, res.Data)

	c, err := api.WithOptions(koios.Host("127.0.0.2:80"))
	assert.NoError(t, err)
	_, err = c.GetPoolMetadata(context.Background(), payload.PoolIDs, nil)
	assert.EqualError(t, err, "dial tcp: lookup 127.0.0.2:80: no such host")
}

func TestGetPoolRelaysEndpoint(t *testing.T) {
	expected := []koios.PoolRelays{}

	spec := loadEndpointTestSpec(t, "endpoint_pool_relays.json.gz", &expected)

	ts, api := setupTestServerAndClient(t, spec)

	defer ts.Close()

	res, err := api.GetPoolRelays(context.Background(), nil)

	assert.NoError(t, err)
	testHeaders(t, spec, res.Response)

	assert.Equal(t, expected, res.Data)

	c, err := api.WithOptions(koios.Host("127.0.0.2:80"))
	assert.NoError(t, err)
	_, err = c.GetPoolRelays(context.Background(), nil)
	assert.EqualError(t, err, "dial tcp: lookup 127.0.0.2:80: no such host")
}

func TestGetPoolUpdatesEndpoint(t *testing.T) {
	expected := []koios.PoolUpdateInfo{}

	spec := loadEndpointTestSpec(t, "endpoint_pool_updates.json.gz", &expected)
	ts, api := setupTestServerAndClient(t, spec)
	defer ts.Close()

	poolID := koios.PoolID(spec.Request.Query.Get("_pool_bech32"))
	res, err := api.GetPoolUpdates(
		context.Background(),
		&poolID,
		nil,
	)

	assert.NoError(t, err)
	testHeaders(t, spec, res.Response)

	assert.Equal(t, expected, res.Data)

	c, err := api.WithOptions(koios.Host("127.0.0.2:80"))
	assert.NoError(t, err)
	_, err = c.GetPoolUpdates(
		context.Background(),
		&poolID,
		nil,
	)
	assert.EqualError(t, err, "dial tcp: lookup 127.0.0.2:80: no such host")
}

func TestGetNativeScriptListEndpoint(t *testing.T) {
	expected := []koios.NativeScriptListItem{}

	spec := loadEndpointTestSpec(t, "endpoint_native_script_list.json.gz", &expected)

	ts, api := setupTestServerAndClient(t, spec)

	defer ts.Close()

	res, err := api.GetNativeScriptList(context.Background(), nil)

	assert.NoError(t, err)
	testHeaders(t, spec, res.Response)

	assert.Equal(t, expected, res.Data)

	c, err := api.WithOptions(koios.Host("127.0.0.2:80"))
	assert.NoError(t, err)
	_, err = c.GetNativeScriptList(context.Background(), nil)
	assert.EqualError(t, err, "dial tcp: lookup 127.0.0.2:80: no such host")
}

func TestGetPlutusScriptNativeListEndpoint(t *testing.T) {
	expected := []koios.PlutusScriptListItem{}

	spec := loadEndpointTestSpec(t, "endpoint_plutus_script_list.json.gz", &expected)

	ts, api := setupTestServerAndClient(t, spec)

	defer ts.Close()

	res, err := api.GetPlutusScriptList(context.Background(), nil)

	assert.NoError(t, err)
	testHeaders(t, spec, res.Response)

	assert.Equal(t, expected, res.Data)

	c, err := api.WithOptions(koios.Host("127.0.0.2:80"))
	assert.NoError(t, err)
	_, err = c.GetPlutusScriptList(context.Background(), nil)
	assert.EqualError(t, err, "dial tcp: lookup 127.0.0.2:80: no such host")
}

func TestGetScriptRedeemersEndpoint(t *testing.T) {
	expected := []koios.ScriptRedeemers{}

	spec := loadEndpointTestSpec(t, "endpoint_script_redeemers.json.gz", &expected)
	ts, api := setupTestServerAndClient(t, spec)
	defer ts.Close()

	res, err := api.GetScriptRedeemers(
		context.Background(),
		koios.ScriptHash(spec.Request.Query.Get("_script_hash")),
		nil,
	)

	assert.NoError(t, err)
	testHeaders(t, spec, res.Response)

	assert.Equal(t, &expected[0], res.Data)

	c, err := api.WithOptions(koios.Host("127.0.0.2:80"))
	assert.NoError(t, err)
	_, err = c.GetScriptRedeemers(
		context.Background(),
		koios.ScriptHash(spec.Request.Query.Get("_script_hash")),
		nil,
	)
	assert.EqualError(t, err, "dial tcp: lookup 127.0.0.2:80: no such host")
}

func TestGetTxInfoEndpoint(t *testing.T) {
	expected := []koios.TxInfo{}

	spec := loadEndpointTestSpec(t, "endpoint_tx_info.json.gz", &expected)

	ts, api := setupTestServerAndClient(t, spec)

	defer ts.Close()

	var payload = struct {
		TxHashes []koios.TxHash `json:"_tx_hashes"`
	}{}
	err := json.Unmarshal(spec.Request.Body, &payload)
	assert.NoError(t, err)

	// Valid
	res, err := api.GetTxInfo(context.Background(), payload.TxHashes[0], nil)
	assert.NoError(t, err)
	testHeaders(t, spec, res.Response)
	assert.Equal(t, &expected[0], res.Data)

	// Empty payload
	res2, err := api.GetTxInfo(context.Background(), koios.TxHash(""), nil)
	assert.ErrorIs(t, err, koios.ErrNoTxHash)
	assert.Nil(t, res2.Data)
	if assert.NotNil(t, res2.Error) {
		assert.Equal(t, koios.ErrNoTxHash.Error(), res2.Error.Message)
	}

	c, err := api.WithOptions(koios.Host("127.0.0.2:80"))
	assert.NoError(t, err)
	_, err = c.GetTxInfo(context.Background(), payload.TxHashes[0], nil)
	assert.EqualError(t, err, "dial tcp: lookup 127.0.0.2:80: no such host")
}

func TestGetTxMetadataEndpoint(t *testing.T) {
	expected := []koios.TxMetadata{}

	spec := loadEndpointTestSpec(t, "endpoint_tx_metadata.json.gz", &expected)

	ts, api := setupTestServerAndClient(t, spec)

	defer ts.Close()

	var payload = struct {
		TxHashes []koios.TxHash `json:"_tx_hashes"`
	}{}
	err := json.Unmarshal(spec.Request.Body, &payload)
	assert.NoError(t, err)

	res, err := api.GetTxMetadata(context.Background(), payload.TxHashes[0], nil)

	assert.NoError(t, err)
	testHeaders(t, spec, res.Response)

	assert.Equal(t, &expected[0], res.Data)

	res2, err := api.GetTxsMetadata(context.Background(), []koios.TxHash{}, nil)
	assert.ErrorIs(t, err, koios.ErrNoTxHash)
	assert.Nil(t, res2.Data, "response data should be nil if arg is invalid")
	assert.Equal(t, res2.Error.Message, "missing transaxtion hash(es)")

	c, err := api.WithOptions(koios.Host("127.0.0.2:80"))
	assert.NoError(t, err)
	_, err = c.GetTxMetadata(context.Background(), payload.TxHashes[0], nil)
	assert.EqualError(t, err, "dial tcp: lookup 127.0.0.2:80: no such host")
}

func TestGetTxMetaLabelsEndpoint(t *testing.T) {
	expected := []koios.TxMetalabel{}

	spec := loadEndpointTestSpec(t, "endpoint_tx_metalabels.json.gz", &expected)

	ts, api := setupTestServerAndClient(t, spec)

	defer ts.Close()

	res, err := api.GetTxMetaLabels(context.Background(), nil)

	assert.NoError(t, err)
	testHeaders(t, spec, res.Response)

	assert.Equal(t, expected, res.Data)

	c, err := api.WithOptions(koios.Host("127.0.0.2:80"))
	assert.NoError(t, err)
	_, err = c.GetTxMetaLabels(context.Background(), nil)
	assert.EqualError(t, err, "dial tcp: lookup 127.0.0.2:80: no such host")
}

func TestGetTxStatusEndpoint(t *testing.T) {
	expected := []koios.TxStatus{}

	spec := loadEndpointTestSpec(t, "endpoint_tx_status.json.gz", &expected)

	ts, api := setupTestServerAndClient(t, spec)

	defer ts.Close()

	var payload = struct {
		TxHashes []koios.TxHash `json:"_tx_hashes"`
	}{}
	err := json.Unmarshal(spec.Request.Body, &payload)
	assert.NoError(t, err)

	res, err := api.GetTxStatus(context.Background(), payload.TxHashes[0], nil)

	assert.NoError(t, err)
	testHeaders(t, spec, res.Response)

	assert.Equal(t, &expected[0], res.Data)

	res2, err := api.GetTxsStatuses(context.Background(), []koios.TxHash{}, nil)
	assert.ErrorIs(t, err, koios.ErrNoTxHash)
	assert.Nil(t, res2.Data, "response data should be nil if arg is invalid")
	assert.Equal(t, res2.Error.Message, "missing transaxtion hash(es)")

	c, err := api.WithOptions(koios.Host("127.0.0.2:80"))
	assert.NoError(t, err)
	_, err = c.GetTxStatus(context.Background(), payload.TxHashes[0], nil)
	assert.EqualError(t, err, "dial tcp: lookup 127.0.0.2:80: no such host")
}

func TestGetTxsUTxOsEndpoint(t *testing.T) {
	expected := []koios.UTxO{}

	spec := loadEndpointTestSpec(t, "endpoint_tx_utxos.json.gz", &expected)

	ts, api := setupTestServerAndClient(t, spec)

	defer ts.Close()

	var payload = struct {
		TxHashes []koios.TxHash `json:"_tx_hashes"`
	}{}
	err := json.Unmarshal(spec.Request.Body, &payload)
	assert.NoError(t, err)

	res, err := api.GetTxsUTxOs(context.Background(), payload.TxHashes, nil)

	assert.NoError(t, err)
	testHeaders(t, spec, res.Response)

	assert.Equal(t, expected, res.Data)

	res2, err := api.GetTxsUTxOs(context.Background(), []koios.TxHash{}, nil)
	assert.ErrorIs(t, err, koios.ErrNoTxHash)
	assert.Nil(t, res2.Data, "response data should be nil if arg is invalid")
	assert.Equal(t, res2.Error.Message, "missing transaxtion hash(es)")

	c, err := api.WithOptions(koios.Host("127.0.0.2:80"))
	assert.NoError(t, err)
	_, err = c.GetTxsUTxOs(context.Background(), payload.TxHashes, nil)
	assert.EqualError(t, err, "dial tcp: lookup 127.0.0.2:80: no such host")
}

func TestTxSubmit(t *testing.T) {
	spec := loadEndpointTestSpec(t, "endpoint_tx_submit.json.gz", nil)

	ts, api := setupTestServerAndClient(t, spec)

	defer ts.Close()
	payload := koios.TxBodyJSON{}
	err := json.Unmarshal(spec.Request.Body, &payload)
	assert.NoError(t, err)

	res, err := api.SubmitSignedTx(context.Background(), payload, nil)

	assert.NoError(t, err)
	assert.Equal(t, spec.Response.Code, 202)
	assert.Equal(t, res.StatusCode, 202)
	assert.Equal(t, res.Status, "202 Accepted")
	testHeaders(t, spec, res.Response)

	res2, err := api.SubmitSignedTx(context.Background(), koios.TxBodyJSON{CborHex: "x"}, nil)

	assert.Error(t, err, "submited tx should return error")
	assert.Equal(t, res2.StatusCode, 400)

	c, err := api.WithOptions(koios.Host("127.0.0.2:80"))
	assert.NoError(t, err)
	_, err = c.SubmitSignedTx(context.Background(), payload, nil)
	assert.EqualError(t, err, "dial tcp: lookup 127.0.0.2:80: no such host")
}

func TestVerticalFiltering(t *testing.T) {
	expected := []koios.Block{}

	spec := loadEndpointTestSpec(t, "vertical_filtering.json.gz", &expected)

	ts, api := setupTestServerAndClient(t, spec)

	defer ts.Close()

	opts := api.NewRequestOptions()
	opts.HeadersApply(spec.Request.Header)
	opts.QueryApply(spec.Request.Query)

	res, err := api.GetBlocks(context.Background(), opts)

	assert.NoError(t, err)
	testHeaders(t, spec, res.Response)

	assert.Equal(t, expected, res.Data)

	c, err := api.WithOptions(koios.Host("127.0.0.2:80"))
	assert.NoError(t, err)

	opts2 := api.NewRequestOptions()
	opts2.HeadersApply(spec.Request.Header)
	opts2.QueryApply(spec.Request.Query)

	_, err = c.GetBlocks(context.Background(), opts2)
	assert.EqualError(t, err, "dial tcp: lookup 127.0.0.2:80: no such host")
}

func TestHorizontalFiltering(t *testing.T) {
	expected := []koios.Block{}

	spec := loadEndpointTestSpec(t, "horizontal_filtering.json.gz", &expected)

	ts, api := setupTestServerAndClient(t, spec)

	defer ts.Close()

	opts := api.NewRequestOptions()
	opts.HeadersApply(spec.Request.Header)
	opts.QueryApply(spec.Request.Query)
	res, err := api.GetBlocks(context.Background(), opts)

	assert.NoError(t, err)
	testHeaders(t, spec, res.Response)

	assert.Equal(t, expected, res.Data)

	c, err := api.WithOptions(koios.Host("127.0.0.2:80"))
	assert.NoError(t, err)

	opts2 := api.NewRequestOptions()
	opts2.HeadersApply(spec.Request.Header)
	opts2.QueryApply(spec.Request.Query)
	_, err = c.GetBlocks(context.Background(), opts2)
	assert.EqualError(t, err, "dial tcp: lookup 127.0.0.2:80: no such host")
}

//nolint: funlen
func TestHTTP(t *testing.T) {
	expected := []koios.Tip{}

	spec := loadEndpointTestSpec(t, "endpoint_network_tip.json.gz", &expected)

	ts, api := setupTestServerAndClient(t, spec)

	defer ts.Close()

	// GET
	opts := api.NewRequestOptions()
	for k, vv := range spec.Request.Header {
		for _, v := range vv {
			opts.HeadersAdd(k, v)
		}
	}
	for k, vv := range spec.Request.Query {
		for _, v := range vv {
			opts.QueryAdd(k, v)
		}
	}
	opts2, o2err := opts.Clone()
	assert.NoError(t, o2err)
	opts3, o3err := opts.Clone()
	assert.NoError(t, o3err)

	res, err := api.GET(context.Background(), "/tip", opts2)
	assert.NoError(t, err)

	body, err := io.ReadAll(res.Body)
	defer func() { _ = res.Body.Close() }()
	assert.NoError(t, err)

	data := []koios.Tip{}
	err = json.Unmarshal(body, &data)
	assert.NoError(t, err)

	assert.Len(t, expected, 1)
	assert.Equal(t, expected, data)

	// HEAD
	res2, err2 := api.HEAD(context.Background(), "/tip", opts3)
	defer func() { _ = res2.Body.Close() }()
	data2 := []koios.Tip{}
	assert.NoError(t, koios.ReadAndUnmarshalResponse(res2, &koios.Response{}, &data2))

	assert.NoError(t, err2)
	assert.Equal(t, "application/json; charset=utf-8", res2.Header.Get("Content-Type"))

	_, o4err := opts3.Clone()
	assert.ErrorIs(t, o4err, koios.ErrReqOptsAlreadyUsed)
	opts4 := api.NewRequestOptions()
	opts4.QueryApply(spec.Request.Query)
	opts4.HeadersApply(spec.Request.Header)

	// 404
	rsp3, err3 := api.HEAD(context.Background(), "/404", opts4)
	res3 := &koios.Response{}
	defer func() { _ = rsp3.Body.Close() }()

	assert.EqualError(t, koios.ReadAndUnmarshalResponse(rsp3, res3, nil), "got non json response: ")
	assert.EqualError(t, err3, "got unexpected response: 404 Not Found")
	assert.Equal(t, rsp3.Header.Get("Content-Type"), "text/plain; charset=utf-8")
}

// loadEndpointTestSpec load specs for endpoint.
func loadEndpointTestSpec(t *testing.T, filename string, exp interface{}) *internal.APITestSpec {
	spec := &internal.APITestSpec{}
	spec.Response.Body = exp
	gzfile, err := os.Open(filepath.Join("testdata", filename))
	assert.NoErrorf(t, err, "failed to open test compressed spec: %s", filename)
	defer gzfile.Close()

	gzr, err := gzip.NewReader(gzfile)
	assert.NoErrorf(t, err, "failed create reader for test spec: %s", filename)

	specb, err := io.ReadAll(gzr)
	assert.NoErrorf(t, err, "failed to read test spec: %s", filename)
	gzr.Close()

	assert.NoErrorf(t, json.Unmarshal(specb, &spec), "failed to Unmarshal test spec: %s", filename)
	return spec
}

// INTENAL TEST UTILS

// setupTestServerAndClient httptest server and api client based on specs.
//nolint: gocognit, funlen
func setupTestServerAndClient(t *testing.T, spec *internal.APITestSpec) (*httptest.Server, *koios.Client) {
	mux := http.NewServeMux()
	endpoint := fmt.Sprintf("/api/%s%s", koios.DefaultAPIVersion, spec.Endpoint)
	mux.HandleFunc(endpoint, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != spec.Request.Method && r.Method != "HEAD" {
			http.Error(w, "Method Not Allowed.", http.StatusMethodNotAllowed)
			return
		}

		// Add response headers
		for header, values := range spec.Response.Header {
			for _, value := range values {
				w.Header().Add(header, value)
			}
		}

		//nolint: nestif
		if spec.Request.Method == "POST" {
			var expectedBody map[string]interface{}
			if err := json.Unmarshal(spec.Request.Body, &expectedBody); err != nil {
				http.Error(w, "failed to verify expected post body", spec.Response.Code)
				return
			}
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "failed to read post body", spec.Response.Code)
				return
			}
			var submitedBody map[string]interface{}

			if err := json.Unmarshal(body, &submitedBody); err != nil {
				http.Error(w, "failed to verify submitted post body", spec.Response.Code)
				return
			}
			for k, v := range expectedBody {
				val, ok := submitedBody[k]
				if !ok {
					http.Error(w, fmt.Sprintf("did not find expected post body: %s", k), spec.Response.Code)
					return
				}
				expected := fmt.Sprint(v)
				actual := fmt.Sprint(val)
				if expected != actual {
					http.Error(
						w,
						fmt.Sprintf(
							"post body: %s has invalid value(%v) expected(%v)",
							k,
							actual,
							expected,
						), spec.Response.Code)
					return
				}
			}
		}
		w.WriteHeader(spec.Response.Code)

		// Add response payload
		res, err := json.Marshal(spec.Response.Body)
		if err != nil {
			http.Error(w, "failed to marshal response", spec.Response.Code)
			return
		}

		w.Write(res)
	})

	ts := httptest.NewUnstartedServer(mux)
	ts.EnableHTTP2 = true
	ts.StartTLS()

	u, err := url.Parse(ts.URL)
	assert.NoErrorf(t, err, "failed to parse test server url: %s", ts.URL)
	port, err := strconv.ParseUint(u.Port(), 0, 16)
	assert.NoError(t, err, "failed to parse port from server url %s", ts.URL)

	client := ts.Client()
	client.Timeout = time.Second * 10
	c, err := koios.New(
		koios.HTTPClient(client),
		koios.Port(uint16(port)),
		koios.Host(u.Hostname()),
		koios.CollectRequestsStats(true),
	)
	assert.NoError(t, err, "failed to create api client")
	return ts, c
}

// testHeaders universal header tester.
// Currently testing only headers we care about.
func testHeaders(t *testing.T, spec *internal.APITestSpec, res koios.Response) {
	assert.Equalf(
		t,
		spec.Request.Method,
		res.RequestMethod,
		"%s: invalid request method (%s)",
		spec.Request.Method,
		res.Status,
	)
	assert.Equalf(t, spec.Response.Code, res.StatusCode, "%s: invalid response code (%s)", spec.Request.Method, res.Status)
	assert.Equalf(
		t,
		res.ContentRange,
		spec.Response.Header.Get("content-range"),
		"%s: has invalid content-range header", spec.Request.Method,
	)
	assert.Equalf(
		t,
		res.ContentLocation,
		spec.Response.Header.Get("content-location"),
		"%s: has invalid content-location header",
		spec.Request.Method,
	)
}
