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
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/http/httptrace"
	"net/url"
	"strings"
	"time"
)

// WithOptions returns new light clone of client with modified options applied.
func (c *Client) WithOptions(opts ...Option) (*Client, error) {
	nc := &Client{
		r:               c.r,
		reqStatsEnabled: c.reqStatsEnabled,
		url:             c.url,
		client:          c.client,
		commonHeaders:   c.commonHeaders.Clone(),
	}
	// Apply provided options
	for _, opt := range opts {
		if err := opt.apply(c); err != nil {
			return nil, err
		}
	}
	return nc, nil
}

// HEAD sends api http HEAD request to provided relative path with query params
// and returns an HTTP response.
func (c *Client) HEAD(
	ctx context.Context,
	path string,
	query url.Values,
	headers http.Header,
) (*http.Response, error) {
	return c.request(ctx, nil, "HEAD", path, nil, query, headers)
}

// POST sends api http POST request to provided relative path with query params
// and returns an HTTP response. When using POST method you are expected
// to handle the response according to net/http.Do documentation.
// e.g. Caller should close resp.Body when done reading from it.
func (c *Client) POST(
	ctx context.Context,
	path string,
	body io.Reader,
	query url.Values,
	headers http.Header,
) (*http.Response, error) {
	return c.request(ctx, nil, "POST", path, body, query, headers)
}

// GET sends api http GET request to provided relative path with query params
// and returns an HTTP response. When using GET method you are expected
// to handle the response according to net/http.Do documentation.
// e.g. Caller should close resp.Body when done reading from it.
func (c *Client) GET(
	ctx context.Context,
	path string,
	query url.Values,
	headers http.Header,
) (*http.Response, error) {
	return c.request(ctx, nil, "GET", path, nil, query, headers)
}

// BaseURL returns currently used base url e.g. https://api.koios.rest/api/v0
func (c *Client) BaseURL() string {
	return c.url.String()
}

func (c *Client) request(
	ctx context.Context,
	res *Response,
	method string,
	path string,
	body io.Reader,
	query url.Values,
	headers http.Header) (*http.Response, error) {
	var (
		requrl string
	)

	path = strings.TrimLeft(path, "/")

	if query == nil {
		requrl = c.url.ResolveReference(&url.URL{Path: path}).String()
	} else {
		requrl = c.url.ResolveReference(&url.URL{Path: path, RawQuery: query.Encode()}).String()
	}

	if res != nil {
		res.RequestURL = requrl
	}

	// handle rate limit
	if err := c.r.Wait(ctx); err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, strings.ToUpper(method), requrl, body)
	if err != nil {
		if res != nil {
			res.applyError(nil, err)
		}
		return nil, err
	}

	c.applyReqHeaders(req, headers)

	if res != nil && c.reqStatsEnabled {
		return c.requestWithStats(req, res)
	}

	rsp, err := c.client.Do(req)
	if err != nil {
		if res != nil {
			res.applyError(nil, err)
		}
		return nil, err
	}

	if res != nil {
		res.applyRsp(rsp)
	}

	if rsp.StatusCode > http.StatusCreated {
		if res != nil {
			res.applyError(nil, ErrResponse)
		}
		return rsp, nil
	}

	return rsp, nil
}

func (c *Client) applyReqHeaders(req *http.Request, headers http.Header) {
	req.Header = c.commonHeaders.Clone()
	if headers != nil {
		for name, values := range headers {
			for _, value := range values {
				req.Header.Add(name, value)
			}
		}
		return
	}
	// only apply if originally there were no headers defined.
	switch req.Method {
	case "POST":
	case "PATCH":
	case "PUT":
		req.Header.Set("Content-Type", "application/json")
	}
}

func (c *Client) requestWithStats(req *http.Request, res *Response) (*http.Response, error) {
	res.Stats = &RequestStats{}
	var dns, tlshs, connect time.Time

	trace := &httptrace.ClientTrace{
		DNSStart: func(dsi httptrace.DNSStartInfo) {
			dns = time.Now().UTC()
		},
		DNSDone: func(ddi httptrace.DNSDoneInfo) {
			res.Stats.DNSLookupDur = time.Since(dns)
		},
		TLSHandshakeStart: func() {
			tlshs = time.Now().UTC()
		},
		TLSHandshakeDone: func(cs tls.ConnectionState, err error) {
			if err != nil {
				res.applyError(nil, err)
			}
			res.Stats.TLSHSDur = time.Since(tlshs)
		},
		ConnectStart: func(network, addr string) {
			connect = time.Now().UTC()
		},
		ConnectDone: func(network, addr string, err error) {
			if err != nil {
				res.applyError(nil, err)
			}
			res.Stats.ESTCXNDur = time.Since(connect)
		},
		GotFirstResponseByte: func() {
			res.Stats.TTFB = time.Since(res.Stats.ReqStartedAt)
		},
	}
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))

	res.Stats.ReqStartedAt = time.Now().UTC()
	rsp, err := c.client.Transport.RoundTrip(req)
	if err != nil {
		res.applyError(nil, err)
		return nil, err
	}

	res.applyRsp(rsp)

	if rsp.StatusCode > http.StatusCreated {
		res.applyError(nil, ErrResponse)
		return rsp, nil
	}

	return rsp, nil
}

func (c *Client) setBaseURL(schema, host, version string, port uint16) error {
	raw := fmt.Sprintf("%s://%s", schema, host)
	if port != 80 && port != 443 {
		raw = fmt.Sprintf("%s:%d", raw, port)
	}
	raw += "/api/" + version + "/"
	u, err := url.ParseRequestURI(raw)
	if err != nil {
		return err
	}
	c.url = u
	return nil
}
