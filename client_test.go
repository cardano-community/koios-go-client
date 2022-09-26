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
// 	"encoding/json"
// 	"io"
// 	"net/http"
// 	"strings"
// 	"testing"
// 	"time"

// 	"github.com/cardano-community/koios-go-client/v2"
// 	"github.com/stretchr/testify/suite"
// )

// func TestClientSuite(t *testing.T) {
// 	testsuite := &clientTestSuite{}
// 	testsuite.LoadSpecs([]string{
// 		"endpoint_network_tip",
// 	})
// 	suite.Run(t, testsuite)
// }

// type clientTestSuite struct {
// 	endpointsTestSuite
// }

// func (s *clientTestSuite) TestRequestContext() {
// 	res, err := s.api.GetTip(nil, nil) //nolint: staticcheck
// 	s.EqualError(err, "net/http: nil Context")
// 	s.Equal("net/http: nil Context", res.Error.Message)

// 	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(0))
// 	defer cancel()

// 	res2, err := s.api.GetTip(ctx, nil)
// 	s.EqualError(err, "context deadline exceeded")

// 	s.Nil(res2.Data)

// 	ctx2, cancel2 := context.WithTimeout(context.Background(), time.Microsecond*150)
// 	defer cancel2()

// 	res3, err3 := s.api.GetTip(ctx2, nil)
// 	var edgeerr bool
// 	if err3.Error() == "context deadline exceeded" || strings.Contains(err3.Error(), "i/o timeout") {
// 		edgeerr = true
// 	}
// 	s.True(edgeerr, "expected: context deadline exceeded or i/o timeout")
// 	s.Nil(res3.Data)
// }

// func (s *clientTestSuite) Test404s() {
// 	res, err := s.api.GetGenesis(context.Background(), nil)
// 	s.Error(err)
// 	s.Nil(res.Data)
// 	s.Equal("http error: 404 Not Found", res.Error.Message)
// 	s.Equal(http.StatusNotFound, res.StatusCode)

// 	// errors with stats should be same
// 	res2, err := s.api.GetGenesis(context.Background(), nil)
// 	s.Error(err)
// 	s.Nil(res2.Data)
// 	s.Equal("http error: 404 Not Found", res2.Error.Message)
// 	s.Equal(http.StatusNotFound, res2.StatusCode)
// }

// //nolint: funlen
// func (s *clientTestSuite) TestHTTP() {
// 	spec := s.GetSpec("endpoint_network_tip")
// 	if s.NotNil(spec) {
// 		// GET
// 		opts := s.api.NewRequestOptions()
// 		for k, vv := range spec.Request.Header {
// 			for _, v := range vv {
// 				opts.HeadersAdd(k, v)
// 			}
// 		}
// 		for k, vv := range spec.Request.Query {
// 			for _, v := range vv {
// 				opts.QueryAdd(k, v)
// 			}
// 		}
// 		extheaders := http.Header{}
// 		extheaders.Set("X-Ext-Header", "ext-val")
// 		opts.HeadersAdd("X-Custom-Header", "header-val")
// 		opts.QueryAdd("prop", "val")
// 		opts.HeadersApply(extheaders)
// 		opts2 := opts.Clone()
// 		s.NotNil(opts2)
// 		opts3 := s.api.NewRequestOptions()

// 		res, err := s.api.GET(context.Background(), "/tip", opts2)
// 		if s.NoError(err) {
// 			body, err := io.ReadAll(res.Body)
// 			defer func() { _ = res.Body.Close() }()
// 			s.NoError(err)
// 			data := []koios.Tip{}
// 			err = json.Unmarshal(body, &data)
// 			s.NoError(err)
// 		}

// 		// HEAD
// 		res2, err2 := s.api.HEAD(context.Background(), "/tip", opts3)
// 		defer func() { _ = res2.Body.Close() }()
// 		data2 := []koios.Tip{}
// 		s.NoError(koios.ReadAndUnmarshalResponse(res2, &koios.Response{}, &data2))
// 		if s.NoError(err2) {
// 			s.Equal("application/json; charset=utf-8", res2.Header.Get("Content-Type"))
// 		}
// 		res21, err21 := s.api.HEAD(context.Background(), "/tip", opts3)
// 		s.Nil(res21)
// 		if res21 != nil {
// 			res21.Body.Close()
// 		}
// 		s.ErrorIs(err21, koios.ErrReqOptsAlreadyUsed)
// 		s.NotNil(opts3.Clone())

// 		opts4 := s.api.NewRequestOptions()
// 		opts4.QueryApply(spec.Request.Query)
// 		opts4.HeadersApply(spec.Request.Header)

// 		// 404
// 		rsp3, err3 := s.api.HEAD(context.Background(), "/404", opts4)
// 		res3 := &koios.Response{}
// 		defer func() { _ = rsp3.Body.Close() }()
// 		s.EqualError(koios.ReadAndUnmarshalResponse(rsp3, res3, nil), "got non json response: ")
// 		s.EqualError(err3, "http error: 404 Not Found")
// 		s.Equal(rsp3.Header.Get("Content-Type"), "text/plain; charset=utf-8")
// 	}
// }

// func TestVerticalFiltering(t *testing.T) {
// 	s := &endpointsTestSuite{}
// 	s.LoadSpecs([]string{
// 		"vertical_filtering",
// 	})
// 	s.SetT(t)
// 	s.SetupSuite()
// 	defer s.TearDownSuite()

// 	spec := s.GetSpec("vertical_filtering")
// 	if s.NotNil(spec) {
// 		opts := s.api.NewRequestOptions()
// 		opts.HeadersApply(spec.Request.Header)
// 		opts.QueryApply(spec.Request.Query)

// 		res, err := s.api.GetBlocks(context.Background(), opts)
// 		if s.NoError(err) {
// 			s.Greater(res.Data[0].Height, 0)
// 			s.Greater(res.Data[0].EpochSlot, 0)
// 			s.Empty(res.Data[0].Hash)
// 		}
// 	}
// }

// func TestHorizontalFiltering(t *testing.T) {
// 	s := &endpointsTestSuite{}
// 	s.LoadSpecs([]string{
// 		"horizontal_filtering",
// 	})
// 	s.SetT(t)
// 	s.SetupSuite()
// 	defer s.TearDownSuite()

// 	spec := s.GetSpec("horizontal_filtering")
// 	if s.NotNil(spec) {
// 		opts := s.api.NewRequestOptions()
// 		opts.HeadersApply(spec.Request.Header)
// 		opts.QueryApply(spec.Request.Query)

// 		res, err := s.api.GetBlocks(context.Background(), opts)
// 		if s.NoError(err) && s.Equal(2, len(res.Data)) {
// 			s.Greater(res.Data[0].Height, 0)
// 			s.Greater(res.Data[0].EpochSlot, 0)
// 			s.NotEmpty(res.Data[0].Hash)
// 		}
// 	}
// }

// func TestPaginator(t *testing.T) {
// 	s := &endpointsTestSuite{}
// 	s.LoadSpecs([]string{
// 		"pagination-page-1",
// 	})
// 	s.SetT(t)
// 	s.SetupSuite()
// 	defer s.TearDownSuite()

// 	spec := s.GetSpec("pagination-page-1")
// 	if s.NotNil(spec) {
// 		opts := s.api.NewRequestOptions()
// 		opts.PageSize(10)
// 		opts.Page(1)
// 		opts.QueryApply(spec.Request.Query)
// 		res, err := s.api.GetBlocks(context.Background(), opts)

// 		if s.NoError(err) && s.Equal(10, len(res.Data)) {
// 			s.Equal(1, res.Data[0].Height)
// 			s.Equal(0, res.Data[0].EpochSlot)
// 			s.Empty(res.Data[0].Hash)
// 		}
// 	}

// 	s2 := &endpointsTestSuite{}
// 	s2.LoadSpecs([]string{
// 		"pagination-page-2",
// 	})
// 	s2.SetT(t)
// 	s2.SetupSuite()
// 	defer s2.TearDownSuite()

// 	spec2 := s2.GetSpec("pagination-page-2")
// 	if s2.NotNil(spec2) {
// 		opts := s.api.NewRequestOptions()
// 		opts.PageSize(10)
// 		opts.Page(1)
// 		opts.QueryApply(spec.Request.Query)
// 		res, err := s2.api.GetBlocks(context.Background(), opts)

// 		if s2.NoError(err) && s2.Equal(10, len(res.Data)) {
// 			s2.Equal(11, res.Data[0].Height)
// 			s2.Equal(10, res.Data[0].EpochSlot)
// 			s2.Empty(res.Data[0].Hash)
// 		}
// 	}
// }
