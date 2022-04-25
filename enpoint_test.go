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
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/cardano-community/koios-go-client"
	"github.com/cardano-community/koios-go-client/internal"
	"github.com/stretchr/testify/suite"
)

type endpointsTestSuite struct {
	suite.Suite
	specNames []string
	srv       *httptest.Server
	api       *koios.Client
	specs     []internal.APITestSpec
}

func (s *endpointsTestSuite) LoadSpecs(specNames []string) {
	s.specNames = specNames
}

func (s *endpointsTestSuite) GetSpec(specName string) *internal.APITestSpec {
	for _, spec := range s.specs {
		if spec.Filename == specName {
			return &spec
		}
	}
	return nil
}

func (s *endpointsTestSuite) TearDownSuite() {
	s.srv.Close()
}

func (s *endpointsTestSuite) SetupSuite() {
	mux := http.NewServeMux()

	for _, specName := range s.specNames {
		spec := internal.APITestSpec{}
		gzfile, err := os.Open(filepath.Join("testdata", specName+".json.gz"))
		s.NoErrorf(err, "failed to open test compressed spec: %s", specName)
		defer gzfile.Close()

		gzr, err := gzip.NewReader(gzfile)
		s.NoErrorf(err, "failed create reader for test spec: %s", specName)

		specb, err := io.ReadAll(gzr)
		s.NoErrorf(err, "failed to read test spec: %s", specName)
		gzr.Close()

		s.NoErrorf(json.Unmarshal(specb, &spec), "failed to Unmarshal test spec: %s", specName)

		endpoint := fmt.Sprintf("/api/%s%s", koios.DefaultAPIVersion, spec.Endpoint)
		mux.HandleFunc(endpoint, s.newHandleFunc(spec))
		s.specs = append(s.specs, spec)
	}

	s.srv = httptest.NewUnstartedServer(mux)
	s.srv.EnableHTTP2 = true
	s.srv.StartTLS()
	u, err := url.Parse(s.srv.URL)
	s.NoErrorf(err, "failed to parse test server url: %s", s.srv.URL)
	port, err := strconv.ParseUint(u.Port(), 0, 16)
	s.NoError(err, "failed to parse port from server url %s", s.srv.URL)
	client := s.srv.Client()
	client.Timeout = time.Second * 10

	c, err := koios.New(
		koios.HTTPClient(client),
		koios.Port(uint16(port)),
		koios.Host(u.Hostname()),
		koios.CollectRequestsStats(true),
	)
	s.NoError(err, "failed to create api client")

	c2, err := c.WithOptions(koios.Host("127.0.0.2:80"))
	s.NoError(err)
	_, err = c2.GetTip(context.Background(), nil)
	s.EqualError(err, "dial tcp: lookup 127.0.0.2:80: no such host")
	s.api = c
}

func (s *endpointsTestSuite) newHandleFunc(spec internal.APITestSpec) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
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
		w.WriteHeader(spec.Response.Code)
		res, err := json.Marshal(spec.Response.Body)
		if err != nil {
			http.Error(w, "failed to marshal response", spec.Response.Code)
			return
		}
		w.Write(res)
	}
}
