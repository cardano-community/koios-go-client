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

package koios

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDefaults(t *testing.T) {
	api, err := New()
	assert.NoError(t, err)
	if assert.NotNil(t, api) {
		raw := fmt.Sprintf(
			"%s://%s/api/%s/",
			DefaultSchema,
			MainnetHost,
			DefaultAPIVersion,
		)
		u, err := url.ParseRequestURI(raw)
		assert.NoError(t, err, "default url can not be constructed")
		assert.Equal(t, u.String(), api.BaseURL(), "invalid default base url")
	}
}

func TestOptions(t *testing.T) {
	api, err := New(
		Host("localhost"),
		APIVersion("v1"),
		Port(8080),
		Schema("http"),
		RateLimit(100),
		Origin("http://localhost.localdomain"),
		CollectRequestsStats(true),
	)
	assert.NoError(t, err)
	if assert.NotNil(t, api) {
		assert.Equal(t, "http://localhost:8080/api/v1/", api.BaseURL(), "invalid default base url")
	}
}

func TestOptionErrs(t *testing.T) {
	client, _ := New()
	assert.Error(t, HTTPClient(http.DefaultClient)(client),
		"should not allow changing http client.")
	assert.Error(t, RateLimit(0)(client),
		"should not unlimited requests p/s")
	assert.Error(t, Origin("localhost")(client),
		"origin should be valid http origin")
	_, err := New(Origin("localhost.localdomain"))
	assert.Error(t, err, "New should return err when option is invalid")
}

func TestHTTPClient(t *testing.T) {
	client, err := New(HTTPClient(http.DefaultClient))
	assert.Nil(t, client, "client should be nil if there was error")
	assert.Error(t, err, "should not accept default http client")
}

func TestReadResponseBody(t *testing.T) {
	// enure that readResponseBody behaves consistently
	nil1, nil2 := readResponseBody(nil)
	assert.Nil(t, nil1)
	assert.Nil(t, nil2)
}
