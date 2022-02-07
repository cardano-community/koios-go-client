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
	"fmt"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/howijd/koios-rest-go-client"
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
