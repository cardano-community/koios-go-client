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

// import (
// 	"context"
// 	"testing"

// 	"github.com/cardano-community/koios-go-client/v2"
// 	"github.com/stretchr/testify/suite"
// )

// // In order for 'go test' to run this suite, we need to create
// // a normal test function and pass our suite to suite.Run.
// func TestScriptSuite(t *testing.T) {
// 	testsuite := &sciptTestSuite{}
// 	testsuite.LoadSpecs([]string{
// 		"endpoint_script_redeemers",
// 		"endpoint_plutus_script_list",
// 		"endpoint_native_script_list",
// 	})

// 	suite.Run(t, testsuite)
// }

// type sciptTestSuite struct {
// 	endpointsTestSuite
// }

// func (s *sciptTestSuite) TestGetNativeScriptListEndpoint() {
// 	res, err := s.api.GetNativeScriptList(context.Background(), nil)
// 	if s.NoError(err) && s.Len(res.Data, 1000) {
// 		s.NotEmpty(res.Data[0].CreationTxHash)
// 		s.NotEmpty(res.Data[0].ScriptHash)
// 		s.NotEmpty(res.Data[0].Type)
// 	}
// }
// func (s *sciptTestSuite) TestGetPlutusScriptNativeListEndpoint() {
// 	res, err := s.api.GetPlutusScriptList(context.Background(), nil)
// 	if s.NoError(err) && s.Len(res.Data, 1000) {
// 		s.NotEmpty(res.Data[0].CreationTxHash)
// 		s.NotEmpty(res.Data[0].ScriptHash)
// 	}
// }

// func (s *sciptTestSuite) TestGetScriptRedeemersEndpoint() {
// 	spec := s.GetSpec("endpoint_script_redeemers")
// 	if s.NotNil(spec) {
// 		res, err := s.api.GetScriptRedeemers(
// 			context.Background(),
// 			koios.ScriptHash(spec.Request.Query.Get("_script_hash")),
// 			nil,
// 		)
// 		if s.NoError(err) {
// 			s.NotEmpty(res.Data.Redeemers[0].DatumHash)
// 			s.NotEmpty(res.Data.Redeemers[0].TxHash)
// 			s.Greater(res.Data.Redeemers[0].UnitMem, 0)
// 			s.Greater(res.Data.Redeemers[0].UnitSteps, 0)
// 			s.True(res.Data.Redeemers[0].Fee.IsPositive())
// 			s.NotEmpty(res.Data.ScriptHash)
// 		}
// 	}
// }
