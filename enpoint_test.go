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
	"github.com/go-openapi/loads"
	"github.com/stretchr/testify/suite"
)

// introduces breaking change since v1.3.0

type endpointsTestSuite struct {
	suite.Suite
	specNames []string
	mnsrv     *httptest.Server
	mainnet   *koios.Client
	tnsrv     *httptest.Server
	testnet   *koios.Client
	gnsrv     *httptest.Server
	guildnet  *koios.Client
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
	case "guildnet":
		return s.guildnet
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
	s.gnsrv.Close()
}

func (s *endpointsTestSuite) SetupSuite() {
	s.networks = []string{"mainnet", "testnet", "guildnet"}
	s.T().Parallel()

	s.setupForNetwork("testnet")
	s.setupForNetwork("guildnet")
	s.setupForNetwork("mainnet")
}

func (s *endpointsTestSuite) KoiosHttpError(res koios.Response) {
	if s.NotNil(res, "expected response") && s.NotNil(res.Error, "expected response error") {
		s.Equal("http error: 400 Bad Request", res.Error.Message)
		s.Equal(koios.ErrorCode("400"), res.Error.Code)
		s.Equal(400, res.Error.Code.Int())
		s.ErrorIs(res.Error, koios.ErrResponse)
	}
}

func (s *endpointsTestSuite) setupForNetwork(network string) {
	mux := http.NewServeMux()
	for _, specName := range s.specNames {
		spec := internal.APITestSpec{}
		gzfile, err := os.Open(filepath.Join("testdata", network, specName+".json.gz"))
		s.NoErrorf(err, "failed to open compressed spec: %s - %s", network, specName)
		defer gzfile.Close()

		gzr, err := gzip.NewReader(gzfile)
		s.NoErrorf(err, "failed create reader for test spec: %s - %s", network, specName)

		specb, err := io.ReadAll(gzr)
		s.NoErrorf(err, "failed to read test spec: %s - %s", network, specName)
		gzr.Close()

		s.NoErrorf(json.Unmarshal(specb, &spec), "failed to Unmarshal spec: %s - %s", network, specName)
		endpoint := fmt.Sprintf("/api/%s%s", network, spec.Endpoint)
		mux.HandleFunc(endpoint, s.newHandleFunc(spec))
		s.specs = append(s.specs, spec)
	}

	srv := httptest.NewUnstartedServer(mux)
	srv.EnableHTTP2 = true
	srv.StartTLS()

	u, err := url.Parse(srv.URL)
	s.NoErrorf(err, "failed to parse server url: %s - %s", network, srv.URL)

	port, err := strconv.ParseUint(u.Port(), 0, 16)
	s.NoError(err, "failed to parse port from server url %s - %s", network, srv.URL)

	client := srv.Client()
	client.Timeout = time.Second * 10

	c, err := koios.New(
		koios.HTTPClient(client),
		koios.Port(uint16(port)),
		koios.Host(u.Hostname()),
		koios.CollectRequestsStats(true),
		koios.APIVersion(network),
	)

	s.NoError(err, "failed to create mainnet api client")

	c2, err := c.WithOptions(koios.Host("127.0.0.2:80"))
	s.NoError(err)
	_, err = c2.GetTip(context.Background(), nil)
	s.EqualError(err, "dial tcp: lookup 127.0.0.2:80: no such host")

	switch network {
	case "mainnet":
		s.mainnet = c
		s.mnsrv = srv
		s.mnspec, err = loads.JSONSpec("testdata/mainnet/koiosapi")

	case "testnet":
		s.testnet = c
		s.tnsrv = srv

	case "guildnet":
		s.guildnet = c
		s.gnsrv = srv

	default:
		s.Failf("invalid network passed to setupForNetwork %s ", network)
	}
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
		_, _ = w.Write(res)
	}
}
