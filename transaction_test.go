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
	"net/http"
	"testing"

	"github.com/cardano-community/koios-go-client"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
)

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run.
func TestTransactionsSuite(t *testing.T) {
	testsuite := &transactionTestSuite{}
	testsuite.LoadSpecs([]string{
		"endpoint_tx_info",
		"endpoint_tx_utxos",
		"endpoint_tx_metadata",
		"endpoint_tx_metalabels",
		"endpoint_tx_status",
		"endpoint_tx_submit",
	})

	suite.Run(t, testsuite)
}

type transactionTestSuite struct {
	endpointsTestSuite
}

func (s *transactionTestSuite) TestGetTxInfoEndpoint() {
	spec := s.GetSpec("endpoint_tx_info")
	if s.NotNil(spec) {
		var payload = struct {
			TxHashes []koios.TxHash `json:"_tx_hashes"`
		}{}
		err := json.Unmarshal(spec.Request.Body, &payload)
		if s.NoError(err) {
			res, err := s.api.GetTxInfo(context.Background(), payload.TxHashes[0], nil)
			if s.NoError(err) {
				s.Greater(res.Data.AbsoluteSlot, uint64(0))
				s.Len(res.Data.AssetsMinted, 0)
				s.NotEmpty(res.Data.BlockHash)
				s.Greater(res.Data.BlockHeight, uint64(0))
				s.True(res.Data.Deposit.GreaterThanOrEqual(decimal.Zero))
				s.Greater(res.Data.Epoch, uint64(0))
				s.Greater(res.Data.EpochSlot, uint32(0))
				s.True(res.Data.Fee.IsPositive())
				s.Len(res.Data.Inputs, 1)
				s.Greater(res.Data.InvalidAfter, uint64(0))
				s.Equal(res.Data.InvalidBefore, uint64(0))
				s.True(res.Data.TotalOutput.IsPositive())
				s.Greater(res.Data.TxBlockIndex, uint32(0))
				s.Greater(res.Data.TxSize, uint32(0))
			}
		}

		res2, err := s.api.GetTxInfo(context.Background(), koios.TxHash(""), nil)
		s.ErrorIs(err, koios.ErrNoTxHash)
		s.Nil(res2.Data)
		s.Equal(koios.ErrNoTxHash.Error(), res2.Error.Message)
	}
}

func (s *transactionTestSuite) TestGetTxMetadataEndpoint() {
	spec := s.GetSpec("endpoint_tx_metadata")
	if s.NotNil(spec) {
		var payload = struct {
			TxHashes []koios.TxHash `json:"_tx_hashes"`
		}{}
		err := json.Unmarshal(spec.Request.Body, &payload)
		if s.NoError(err) {
			res, err := s.api.GetTxMetadata(context.Background(), payload.TxHashes[0], nil)
			if s.NoError(err) {
				s.Contains(payload.TxHashes, res.Data.TxHash)
				s.NotNil(res.Data.Metadata)
			}
		}

		res2, err := s.api.GetTxMetadata(context.Background(), koios.TxHash(""), nil)
		s.ErrorIs(err, koios.ErrNoTxHash)
		s.Nil(res2.Data)
		s.Equal(koios.ErrNoTxHash.Error(), res2.Error.Message)
	}
}

func (s *transactionTestSuite) TestGetTxMetaLabelsEndpoint() {
	res, err := s.api.GetTxMetaLabels(context.Background(), nil)
	if s.NoError(err) {
		s.Greater(len(res.Data), 500)
	}
}

func (s *transactionTestSuite) TestGetTxStatusEndpoint() {
	spec := s.GetSpec("endpoint_tx_status")
	if s.NotNil(spec) {
		var payload = struct {
			TxHashes []koios.TxHash `json:"_tx_hashes"`
		}{}
		err := json.Unmarshal(spec.Request.Body, &payload)
		if s.NoError(err) {
			res, err := s.api.GetTxStatus(context.Background(), payload.TxHashes[0], nil)
			if s.NoError(err) {
				s.Equal(payload.TxHashes[0], res.Data.TxHash)
				s.Greater(res.Data.Confirmations, uint64(0))
			}
		}
	}
}

func (s *transactionTestSuite) TestGetTxsUTxOsEndpoint() {
	spec := s.GetSpec("endpoint_tx_utxos")
	if s.NotNil(spec) {
		var payload = struct {
			TxHashes []koios.TxHash `json:"_tx_hashes"`
		}{}
		err := json.Unmarshal(spec.Request.Body, &payload)
		if s.NoError(err) {
			res, err := s.api.GetTxUTxOs(context.Background(), payload.TxHashes[0], nil)
			if s.NoError(err) {
				s.Equal(payload.TxHashes[0].String(), res.Data.TxHash.String())
				s.Greater(len(res.Data.Inputs), 0)
				s.Greater(len(res.Data.Outputs), 0)
			}
		}

		res2, err := s.api.GetTxUTxOs(context.Background(), koios.TxHash(""), nil)
		s.ErrorIs(err, koios.ErrNoTxHash)
		s.Nil(res2.Data)
		s.Equal(koios.ErrNoTxHash.Error(), res2.Error.Message)
	}
}

func (s *transactionTestSuite) TestTxSubmit() {
	spec := s.GetSpec("endpoint_tx_submit")
	if s.NotNil(spec) {
		payload := koios.TxBodyJSON{}
		err := json.Unmarshal(spec.Request.Body, &payload)
		if s.NoError(err) {
			res, err := s.api.SubmitSignedTx(context.Background(), payload, nil)
			if s.NoError(err) {
				s.Equal(202, spec.Response.Code)
				s.Equal(http.StatusAccepted, res.StatusCode)
				s.Equal(res.Status, "202 Accepted")
			}
		}

		res2, err := s.api.SubmitSignedTx(context.Background(), koios.TxBodyJSON{CborHex: "x"}, nil)
		s.Error(err, "submited tx should return error")
		s.Equal(res2.StatusCode, 400)
	}
}
