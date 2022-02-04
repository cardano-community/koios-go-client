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
package koios

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// GET sends api http get request to provided relative path with query params
// and returns an HTTP response. When using these methods you are wxpected
// to handle the response according to net/http.Do documentation.
// e.g. Caller should close resp.Body when done reading from it.
func (c *Client) GET(ctx context.Context, path string, query ...url.Values) (*http.Response, error) {
	return c.request(ctx, "GET", nil, path, query...)
}

func (c *Client) request(
	ctx context.Context,
	m string,
	body io.Reader,
	p string,
	query ...url.Values) (*http.Response, error) {
	var (
		requrl *url.URL
	)

	p = strings.TrimLeft(p, "/")
	c.mux.RLock()
	switch len(query) {
	case 0:
		requrl = c.url.ResolveReference(&url.URL{Path: p})
	case 1:
		requrl = c.url.ResolveReference(&url.URL{Path: p, RawQuery: query[0].Encode()})
	default:
		c.mux.RUnlock()
		return nil, fmt.Errorf("%w: got %d", ErrURLValuesLenght, len(query))
	}
	c.mux.RUnlock()

	// optain lock to update last ts and total
	// request count. Lock will block is another request is already queued.
	// e.g. in other go routine.
	c.mux.Lock()

	// handle rate limit
	for !c.lastRequest.IsZero() && time.Since(c.lastRequest) < c.reqInterval {
	}

	c.lastRequest = time.Now()
	c.totalReq++

	// Release client so that other requests can use it.
	c.mux.Unlock()

	req, err := http.NewRequestWithContext(ctx, strings.ToUpper(m), requrl.String(), body)
	if err != nil {
		return nil, err
	}

	return c.client.Do(req)
}

// Get currently used base url e.g. https://api.koios.rest/api/v0
func (c *Client) BaseURL() string {
	c.mux.RLock()
	defer c.mux.RUnlock()
	return c.url.String()
}

// Get number of total requests made by this API client.
func (c *Client) TotalRequests() uint {
	c.mux.RLock()
	defer c.mux.RUnlock()
	return c.totalReq
}

func (c *Client) updateBaseURL() error {
	c.mux.Lock()
	defer c.mux.Unlock()
	raw := fmt.Sprintf("%s://%s", c.schema, c.host)
	if c.port != 80 && c.port != 443 {
		raw = fmt.Sprintf("%s:%d", raw, c.port)
	}
	raw += "/api/" + c.version + "/"
	u, err := url.ParseRequestURI(raw)
	if err != nil {
		return err
	}
	c.url = u
	return nil
}
