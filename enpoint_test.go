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

	"github.com/cardano-community/koios-go-client/v2"
	"github.com/cardano-community/koios-go-client/v2/internal"
	"github.com/stretchr/testify/suite"
)

type endpointsTestSuite struct {
	suite.Suite
	specNames []string
	mnsrv     *httptest.Server
	mainnet   *koios.Client
	tnsrv     *httptest.Server
	testnet   *koios.Client
	specs     []internal.APITestSpec
	networks  []string
}

func (s *endpointsTestSuite) LoadSpecs(specNames []string) {
	s.specNames = specNames
}

func (s *endpointsTestSuite) Client(network string) *koios.Client {
	switch network {
	case "testnet":
		return s.testnet
	default:
		return s.mainnet
	}
}

func (s *endpointsTestSuite) GetSpec(specName, network string) *internal.APITestSpec {
	for _, spec := range s.specs {
		if spec.Filename == specName && spec.Network == network {
			return &spec
		}
	}
	return nil
}

func (s *endpointsTestSuite) TearDownSuite() {
	s.mnsrv.Close()
	s.tnsrv.Close()
}

func (s *endpointsTestSuite) SetupSuite() {
	s.networks = []string{"mainnet", "testnet"}
	s.T().Parallel()
	s.setupMainnet()
	s.setupTestnet()
}

func (s *endpointsTestSuite) setupMainnet() {
	// mainnet
	mnmux := http.NewServeMux()
	for _, specName := range s.specNames {
		spec := internal.APITestSpec{}
		gzfile, err := os.Open(filepath.Join("testdata", "mainnet", specName+".json.gz"))
		s.NoErrorf(err, "failed to open test compressed spec: %s", specName)
		defer gzfile.Close()

		gzr, err := gzip.NewReader(gzfile)
		s.NoErrorf(err, "failed create reader for test spec: %s", specName)

		specb, err := io.ReadAll(gzr)
		s.NoErrorf(err, "failed to read test spec: %s", specName)
		gzr.Close()

		s.NoErrorf(json.Unmarshal(specb, &spec), "failed to Unmarshal test spec: %s", specName)

		endpoint := fmt.Sprintf("/api/%s%s", "mainnet", spec.Endpoint)
		mnmux.HandleFunc(endpoint, s.newHandleFunc(spec))
		s.specs = append(s.specs, spec)
	}

	s.mnsrv = httptest.NewUnstartedServer(mnmux)
	s.mnsrv.EnableHTTP2 = true
	s.mnsrv.StartTLS()
	mnu, err := url.Parse(s.mnsrv.URL)
	s.NoErrorf(err, "failed to parse test server url: %s", s.mnsrv.URL)
	mnport, err := strconv.ParseUint(mnu.Port(), 0, 16)
	s.NoError(err, "failed to parse port from server url %s", s.mnsrv.URL)
	mnclient := s.mnsrv.Client()
	mnclient.Timeout = time.Second * 10

	c, err := koios.New(
		koios.HTTPClient(mnclient),
		koios.Port(uint16(mnport)),
		koios.Host(mnu.Hostname()),
		koios.CollectRequestsStats(true),
		koios.APIVersion("mainnet"),
	)
	s.NoError(err, "failed to create mainnet api client")

	c2, err := c.WithOptions(koios.Host("127.0.0.2:80"))
	s.NoError(err)
	_, err = c2.GetTip(context.Background(), nil)
	s.EqualError(err, "dial tcp: lookup 127.0.0.2:80: no such host")
	s.mainnet = c
}

func (s *endpointsTestSuite) setupTestnet() {
	// testnet
	tnmux := http.NewServeMux()
	for _, specName := range s.specNames {
		spec := internal.APITestSpec{}
		gzfile, err := os.Open(filepath.Join("testdata", "testnet", specName+".json.gz"))
		s.NoErrorf(err, "failed to open test compressed spec: %s", specName)
		defer gzfile.Close()

		gzr, err := gzip.NewReader(gzfile)
		s.NoErrorf(err, "failed create reader for test spec: %s", specName)

		specb, err := io.ReadAll(gzr)
		s.NoErrorf(err, "failed to read test spec: %s", specName)
		gzr.Close()

		s.NoErrorf(json.Unmarshal(specb, &spec), "failed to Unmarshal test spec: %s", specName)

		endpoint := fmt.Sprintf("/api/%s%s", "testnet", spec.Endpoint)
		tnmux.HandleFunc(endpoint, s.newHandleFunc(spec))
		s.specs = append(s.specs, spec)
	}

	s.tnsrv = httptest.NewUnstartedServer(tnmux)
	s.tnsrv.EnableHTTP2 = true
	s.tnsrv.StartTLS()
	tnu, err := url.Parse(s.tnsrv.URL)
	s.NoErrorf(err, "failed to parse test server url: %s", s.tnsrv.URL)
	tnport, err := strconv.ParseUint(tnu.Port(), 0, 16)
	s.NoError(err, "failed to parse port from server url %s", s.tnsrv.URL)
	tnclient := s.tnsrv.Client()
	tnclient.Timeout = time.Second * 10

	c3, err := koios.New(
		koios.HTTPClient(tnclient),
		koios.Port(uint16(tnport)),
		koios.Host(tnu.Hostname()),
		koios.CollectRequestsStats(true),
		koios.APIVersion("testnet"),
	)
	s.NoError(err, "failed to create testnet api client")

	c4, err := c3.WithOptions(koios.Host("127.0.0.2:80"))
	s.NoError(err)
	_, err = c4.GetTip(context.Background(), nil)
	s.EqualError(err, "dial tcp: lookup 127.0.0.2:80: no such host")
	s.testnet = c3
}

func (s *endpointsTestSuite) newHandleFunc(spec internal.APITestSpec) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Has("test-http") {
			switch r.URL.Query().Get("test-http") {
			case "400":
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("[]"))
			default:
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("[]"))
			}
		}
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
