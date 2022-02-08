// Copyright 2022 The Howijd.Network Authors
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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/howijd/koios-rest-go-client"
	"github.com/howijd/koios-rest-go-client/internal"
)

func TestNewDefaults(t *testing.T) {
	api, err := koios.New()
	assert.NoError(t, err)
	if assert.NotNil(t, api) {
		assert.Equal(t, uint64(0), api.TotalRequests(), "total requests should be 0 by default")

		raw := fmt.Sprintf(
			"%s://%s/api/%s/",
			koios.DefaultSchema,
			koios.MainnetHost,
			koios.DefaultAPIVersion,
		)
		u, err := url.ParseRequestURI(raw)
		assert.NoError(t, err, "default url can not be constructed")
		assert.Equal(t, u.String(), api.BaseURL(), "invalid default base url")
	}
}

// testHeaders universal header tester.
// Currently testing only headers we care about.
func testHeaders(t *testing.T, spec *internal.APITestSpec, res koios.Response) {
	assert.Equalf(t, res.RequestMethod, spec.Request.Method, "%s: invalid request method", spec.Request.Method)
	assert.Equalf(t, res.StatusCode, spec.Response.Code, "%s: invalid response code", spec.Request.Method)
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

// loadEndpointTestSpec load specs for endpoint.
func loadEndpointTestSpec(t *testing.T, filename string, exp interface{}) *internal.APITestSpec {
	spec := &internal.APITestSpec{}
	spec.Response.Body = exp
	gzfile, err := os.Open(filepath.Join("testdata", filename))
	assert.NoErrorf(t, err, "failed to open test compressed spec: %s", filename)
	defer gzfile.Close()

	gzr, err := gzip.NewReader(gzfile)
	assert.NoErrorf(t, err, "failed create reader for test spec: %s", filename)

	specb, err := ioutil.ReadAll(gzr)
	assert.NoErrorf(t, err, "failed to read test spec: %s", filename)
	gzr.Close()

	assert.NoErrorf(t, json.Unmarshal(specb, &spec), "failed to Unmarshal test spec: %s", filename)
	return spec
}

// createTestServerAndClient httptest server and api client based on specs.
func createTestServerAndClient(t *testing.T, spec *internal.APITestSpec) (*httptest.Server, *koios.Client) {
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
		w.WriteHeader(spec.Response.Code)

		// Add response payload
		res, err := json.Marshal(spec.Response.Body)
		if err != nil {
			http.Error(w, "failed to marshal response", http.StatusInternalServerError)
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
	)
	assert.NoError(t, err, "failed to create api client")
	return ts, c
}
