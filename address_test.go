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
	"encoding/json"
	"testing"

	"github.com/cardano-community/koios-go-client"
	"github.com/stretchr/testify/suite"
)

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run.
func TestAddressSuite(t *testing.T) {
	testsuite := &addressTestSuite{}
	testsuite.LoadSpecs([]string{
		"endpoint_address_assets",
		"endpoint_address_info",
		"endpoint_address_txs",
		"endpoint_credential_txs",
	})

	suite.Run(t, testsuite)
}

type addressTestSuite struct {
	endpointsTestSuite
}

func (s *addressTestSuite) TestGetAddressInfoEndpoint() {
	spec := s.GetSpec("endpoint_address_info")
	if s.NotNil(spec) {
		res, err := s.api.GetAddressInfo(
			context.Background(),
			koios.Address(spec.Request.Query.Get("_address")),
			nil,
		)
		if s.NoError(err) {
			s.True(res.Data.Balance.IsPositive())
			s.NotEmpty(res.Data.StakeAddress)
			s.Greater(len(res.Data.UTxOs), 0)
		}
	}
}

func (s *addressTestSuite) TestGetAddressTxsEndpoint() {
	spec := s.GetSpec("endpoint_address_txs")
	if s.NotNil(spec) {
		var payload = struct {
			Adresses         []koios.Address `json:"_addresses"`
			AfterBlockHeight uint64          `json:"_after_block_height,omitempty"`
		}{}
		err := json.Unmarshal(spec.Request.Body, &payload)
		if s.NoError(err) {
			res, err := s.api.GetAddressTxs(context.Background(), payload.Adresses, payload.AfterBlockHeight, nil)
			if s.NoError(err) {
				s.Greater(res.Data[0].BlockHeight, uint64(0))
				s.False(res.Data[0].BlockTime.IsZero())
				s.NotEmpty(res.Data[0].TxHash)
			}
		}

		res2, err := s.api.GetAddressTxs(context.Background(), []koios.Address{}, 0, nil)
		s.ErrorIs(err, koios.ErrNoAddress)
		s.Nil(res2.Data)
		s.Equal("missing address", res2.Error.Message)
	}
}

func (s *addressTestSuite) TestGetAddressAssetsEndpoint() {
	spec := s.GetSpec("endpoint_address_assets")
	if s.NotNil(spec) {
		res, err := s.api.GetAddressAssets(
			context.Background(),
			koios.Address(spec.Request.Query.Get("_address")),
			nil,
		)
		if s.NoError(err) {
			s.NotEmpty(res.Data[0].Name)
			s.NotEmpty(res.Data[0].PolicyID)
			s.True(res.Data[0].Quantity.IsPositive())
		}
		res2, err := s.api.GetAddressAssets(
			context.Background(),
			koios.Address(""),
			nil,
		)
		s.ErrorIs(err, koios.ErrNoAddress)
		s.Nil(res2.Data)
		s.Equal("missing address", res2.Error.Message)
	}
}

func (s *addressTestSuite) TestGetCredentialTxsEndpoint() {
	spec := s.GetSpec("endpoint_credential_txs")
	if s.NotNil(spec) {
		var payload = struct {
			Credentials      []koios.PaymentCredential `json:"_payment_credentials"`
			AfterBlockHeight uint64                    `json:"_after_block_height,omitempty"`
		}{}
		err := json.Unmarshal(spec.Request.Body, &payload)
		if s.NoError(err) {
			res, err := s.api.GetCredentialTxs(
				context.Background(),
				payload.Credentials,
				payload.AfterBlockHeight,
				nil,
			)
			if s.NoError(err) {
				s.Greater(res.Data[0].BlockHeight, uint64(0))
				s.False(res.Data[0].BlockTime.IsZero())
				s.NotEmpty(res.Data[0].TxHash)
			}

			res2, err := s.api.GetCredentialTxs(
				context.Background(),
				[]koios.PaymentCredential{},
				payload.AfterBlockHeight,
				nil,
			)
			s.ErrorIs(err, koios.ErrNoAddress)
			s.Nil(res2.Data)
			s.Equal("missing address", res2.Error.Message)
		}
	}
}
