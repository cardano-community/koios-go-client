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
func TestBlockSuite(t *testing.T) {
	testsuite := &blockTestSuite{}
	testsuite.LoadSpecs([]string{
		"endpoint_blocks",
		"endpoint_block_info",
		"endpoint_block_txs",
	})

	suite.Run(t, testsuite)
}

type blockTestSuite struct {
	endpointsTestSuite
}

func (s *blockTestSuite) TestGetBlocksEndpoint() {
	res, err := s.api.GetBlocks(context.Background(), nil)
	if s.NoError(err) {
		s.Len(res.Data, 1000)
		s.Greater(res.Data[0].AbsoluteSlot, 0)
		s.Greater(res.Data[0].Height, 0)
		s.Greater(res.Data[0].Size, 0)
		s.False(res.Data[0].Time.IsZero(), 0)
		s.Greater(res.Data[0].Epoch, uint64(300))
		s.Greater(res.Data[0].EpochSlot, 0)
		s.NotEmpty(res.Data[0].Hash)
		s.NotEmpty(res.Data[0].Pool)
		s.NotEmpty(res.Data[0].VrfKey)
		s.Greater(res.Data[0].OpCertCounter, 0)
		s.Greater(res.Data[0].TxCount, 0)
	}
}

func (s *blockTestSuite) TestGetBlockInfoEndpoint() {
	spec := s.GetSpec("endpoint_block_info")
	if s.NotNil(spec) {
		var payload = struct {
			BlockHashes []koios.BlockHash `json:"_block_hashes"`
		}{}
		err := json.Unmarshal(spec.Request.Body, &payload)
		s.NoError(err)

		res, err := s.api.GetBlockInfo(
			context.Background(),
			payload.BlockHashes[0],
			nil,
		)

		if s.NoError(err) {
			s.Greater(res.Data.AbsoluteSlot, 0)
			s.Greater(res.Data.Height, 0)
			s.Greater(res.Data.Size, 0)
			s.False(res.Data.Time.IsZero(), 0)
			s.Greater(res.Data.Epoch, uint64(300))
			s.Greater(res.Data.EpochSlot, 0)
			s.NotEmpty(res.Data.Hash)
			s.NotEmpty(res.Data.Pool)
			s.NotEmpty(res.Data.VrfKey)
			s.Greater(res.Data.OpCertCounter, 0)
			s.Greater(res.Data.TxCount, 0)
		}
	}
}

func (s *blockTestSuite) TestGetBlocksInfoEndpoint() {
	spec := s.GetSpec("endpoint_block_info")
	if s.NotNil(spec) {
		var payload = struct {
			BlockHashes []koios.BlockHash `json:"_block_hashes"`
		}{}
		err := json.Unmarshal(spec.Request.Body, &payload)
		s.NoError(err)

		res, err := s.api.GetBlocksInfo(
			context.Background(),
			payload.BlockHashes,
			nil,
		)
		if s.NoError(err) {
			s.Equal(len(payload.BlockHashes), len(res.Data))
		}
	}
}
func (s *blockTestSuite) TestGetBlockTxsEndpoint() {
	spec := s.GetSpec("endpoint_block_txs")
	if s.NotNil(spec) {
		res, err := s.api.GetBlockTxHashes(
			context.Background(),
			koios.BlockHash(spec.Request.Query.Get("_block_hash")),
			nil,
		)
		if s.NoError(err) {
			s.NotEmpty(res.Data[0])
			s.Greater(len(res.Data), 10)
		}
	}
}
