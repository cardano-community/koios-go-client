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
	"testing"

	"github.com/cardano-community/koios-go-client"
	"github.com/stretchr/testify/suite"
)

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run.
func TestAssetSuite(t *testing.T) {
	testsuite := &assetTestSuite{}
	testsuite.LoadSpecs([]string{
		"endpoint_asset_address_list",
		"endpoint_asset_history",
		"endpoint_asset_info",
		"endpoint_asset_list",
		"endpoint_asset_policy_info",
		"endpoint_asset_summary",
		"endpoint_asset_txs",
	})

	suite.Run(t, testsuite)
}

type assetTestSuite struct {
	endpointsTestSuite
}

func (s *assetTestSuite) TestGetAssetListEndpoint() {
	res, err := s.api.GetAssetList(context.Background(), nil)
	if s.NoError(err) {
		s.Len(res.Data, 1000)
		s.Greater(len(res.Data[0].AssetNames.ASCII), 0)
		s.Greater(len(res.Data[0].AssetNames.HEX), 0)
		s.NotEmpty(res.Data[0].PolicyID)
	}
}

func (s *assetTestSuite) TestGetAssetAddressListEndpoint() {
	spec := s.GetSpec("endpoint_asset_address_list")
	if s.NotNil(spec) {
		res, err := s.api.GetAssetAddressList(
			context.Background(),
			koios.PolicyID(spec.Request.Query.Get("_asset_policy")),
			koios.AssetName(spec.Request.Query.Get("_asset_name")),
			nil,
		)
		if s.NoError(err) {
			s.NotEmpty(res.Data[0].PaymentAddress)
			s.True(res.Data[0].Quantity.IsPositive())
		}
	}
}

func (s *assetTestSuite) TestGetAssetInfoEndpoint() {
	spec := s.GetSpec("endpoint_asset_info")
	if s.NotNil(spec) {
		res, err := s.api.GetAssetInfo(
			context.Background(),
			koios.PolicyID(spec.Request.Query.Get("_asset_policy")),
			koios.AssetName(spec.Request.Query.Get("_asset_name")),
			nil,
		)
		if s.NoError(err) {
			s.NotEmpty(res.Data.Name)
			s.NotEmpty(res.Data.NameASCII)
			s.False(res.Data.CreationTime.IsZero())
			s.NotEmpty(res.Data.Fingerprint)
			s.Greater(res.Data.MintCnt, 0)
			s.NotEmpty(res.Data.MintingTxHash)
			s.NotEmpty(res.Data.PolicyID)
			s.True(res.Data.TotalSupply.IsPositive())
		}
	}
}

func (s *assetTestSuite) TestGetAssetPolicyInfoEndpoint() {
	spec := s.GetSpec("endpoint_asset_policy_info")
	if s.NotNil(spec) {
		res, err := s.api.GetAssetPolicyInfo(
			context.Background(),
			koios.PolicyID(spec.Request.Query.Get("_asset_policy")),
			nil,
		)
		if s.NoError(err) {
			s.NotEmpty(res.Data.PolicyID)
			s.NotEmpty(res.Data.Assets[0].Name)
			s.NotEmpty(res.Data.Assets[0].NameASCII)
			s.False(res.Data.Assets[0].CreationTime.IsZero())
			s.True(res.Data.Assets[0].TotalSupply.IsPositive())
		}
	}
}
func (s *assetTestSuite) TestGetAssetSummaryEndpoint() {
	spec := s.GetSpec("endpoint_asset_summary")
	if s.NotNil(spec) {
		res, err := s.api.GetAssetSummary(
			context.Background(),
			koios.PolicyID(spec.Request.Query.Get("_asset_policy")),
			koios.AssetName(spec.Request.Query.Get("_asset_name")),
			nil,
		)
		if s.NoError(err) {
			s.Equal(res.Data.PolicyID, koios.PolicyID(spec.Request.Query.Get("_asset_policy")))
			s.Equal(res.Data.AssetName, koios.AssetName(spec.Request.Query.Get("_asset_name")))
			s.Greater(res.Data.StakedWallets, uint64(0))
			// s.Greater(res.Data.UnstakedAddresses, uint64(0))
			s.Greater(res.Data.TotalTransactions, uint64(0))
		}
	}
}

func (s *assetTestSuite) TestGetAssetTxsEndpoint() {
	spec := s.GetSpec("endpoint_asset_txs")
	if s.NotNil(spec) {
		res, err := s.api.GetAssetTxs(
			context.Background(),
			koios.PolicyID(spec.Request.Query.Get("_asset_policy")),
			koios.AssetName(spec.Request.Query.Get("_asset_name")),
			nil,
		)
		if s.NoError(err) {
			s.Equal(res.Data.PolicyID, koios.PolicyID(spec.Request.Query.Get("_asset_policy")))
			s.Equal(res.Data.AssetName, koios.AssetName(spec.Request.Query.Get("_asset_name")))
			s.Greater(len(res.Data.TxHashes), 0)
		}
	}
}

func (s *assetTestSuite) TestGetAssetHistoryEndpoint() {
	spec := s.GetSpec("endpoint_asset_history")
	if s.NotNil(spec) {
		res, err := s.api.GetAssetHistory(
			context.Background(),
			koios.PolicyID(spec.Request.Query.Get("_asset_policy")),
			koios.AssetName(spec.Request.Query.Get("_asset_name")),
			nil,
		)
		if s.NoError(err) {
			s.Equal(res.Data.PolicyID, koios.PolicyID(spec.Request.Query.Get("_asset_policy")))
			s.Equal(res.Data.AssetName, koios.AssetName(spec.Request.Query.Get("_asset_name")))
			s.Greater(len(res.Data.MintingTXs), 0)
			s.True(res.Data.MintingTXs[0].Quantity.IsPositive())
			s.NotEmpty(res.Data.MintingTXs[0].TxHash)
		}
	}
}
