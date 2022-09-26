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

	"golang.org/x/time/rate"
)

type (
	// Client is api client instance.
	Client struct {
		r               *rate.Limiter
		reqStatsEnabled bool
		url             *url.URL
		client          *http.Client
		commonHeaders   http.Header
	}
)

// WithOptions returns new light clone of client with modified options applied.
func (c *Client) WithOptions(opts ...Option) (*Client, error) {
	nc := &Client{
		r:               c.r,
		reqStatsEnabled: c.reqStatsEnabled,
		commonHeaders:   c.commonHeaders.Clone(),
	}
	u, uerr := url.Parse(c.url.String())
	nc.url = u

	// Apply provided options
	for _, opt := range opts {
		if err := opt.apply(nc); err != nil {
			return nil, err
		}
	}

	if nc.client == nil {
		nc.client = c.client
	}

	return nc, uerr
}

// HEAD sends api http HEAD request to provided relative path with query params
// and returns an HTTP response.
func (c *Client) HEAD(
	ctx context.Context,
	path string,
	opts *RequestOptions,
) (*http.Response, error) {
	return c.request(ctx, nil, "HEAD", path, nil, opts)
}

// POST sends api http POST request to provided relative path with query params
// and returns an HTTP response. When using POST method you are expected
// to handle the response according to net/http.Do documentation.
// e.g. Caller should close resp.Body when done reading from it.
func (c *Client) POST(
	ctx context.Context,
	path string,
	body io.Reader,
	opts *RequestOptions,
) (*http.Response, error) {
	return c.request(ctx, nil, "POST", path, body, opts)
}

// GET sends api http GET request to provided relative path with query params
// and returns an HTTP response. When using GET method you are expected
// to handle the response according to net/http.Do documentation.
// e.g. Caller should close resp.Body when done reading from it.
func (c *Client) GET(
	ctx context.Context,
	path string,
	opts *RequestOptions,
) (*http.Response, error) {
	return c.request(ctx, nil, "GET", path, nil, opts)
}

// BaseURL returns currently used base url e.g. https://api.koios.rest/api/v0
func (c *Client) BaseURL() string {
	return c.url.String()
}

// ServerURL returns currently used server url e.g. https://api.koios.rest/
func (c *Client) ServerURL() *url.URL {
	return c.url.ResolveReference(&url.URL{Path: "/"})
}

func (c *Client) NewRequestOptions() *RequestOptions {
	return &RequestOptions{
		pageSize: PageSize,
		page:     1,
		query:    url.Values{},
		headers:  c.commonHeaders.Clone(),
	}
}

func (c *Client) request(
	ctx context.Context,
	res *Response,
	method string,
	path string,
	body io.Reader,
	opts *RequestOptions,
) (*http.Response, error) {
	if opts == nil {
		opts = c.NewRequestOptions()
	}

	if err := opts.lock(); err != nil {
		return nil, err
	}

	path = strings.TrimLeft(path, "/")
	requrl := c.url.ResolveReference(&url.URL{Path: path, RawQuery: opts.query.Encode()}).String()

	if res != nil {
		res.RequestURL = requrl
		res.RequestMethod = method
	}

	req, err := http.NewRequestWithContext(ctx, strings.ToUpper(method), requrl, body)
	if err != nil {
		if res != nil {
			res.applyError(nil, err)
		}
		return nil, err
	}

	// handle rate limit
	if err := c.r.Wait(ctx); err != nil {
		return nil, err
	}

	c.applyReqHeaders(req, opts.headers)

	var (
		eqerr error
		rsp   *http.Response
	)
	if res != nil && c.reqStatsEnabled {
		rsp, eqerr = c.requestWithStats(req, res)
	} else {
		rsp, eqerr = c.client.Do(req)
	}

	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	if eqerr != nil {
		if res != nil {
			res.applyError(nil, eqerr)
		}
		return nil, eqerr
	}

	if res != nil {
		res.applyRsp(rsp)
	}

	if rsp.StatusCode > http.StatusAccepted {
		rerr := fmt.Errorf("%w: %s", ErrResponse, rsp.Status)
		if res != nil {
			res.applyError(nil, rerr)
		}
		return rsp, rerr
	}

	return rsp, nil
}

func (c *Client) applyReqHeaders(req *http.Request, headers http.Header) {
	for name, values := range headers {
		for _, value := range values {
			req.Header.Add(name, value)
		}
	}
	if req.Method == "POST" && len(headers.Get("Content-Type")) == 0 {
		req.Header.Set("Content-Type", "application/json")
	}
}

func (c *Client) requestWithStats(req *http.Request, res *Response) (*http.Response, error) {
	res.Stats = &RequestStats{}
	var dns, tlshs, connect time.Time
	req = req.WithContext(
		httptrace.WithClientTrace(
			req.Context(),
			&httptrace.ClientTrace{
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
					res.Stats.TLSHSDur = time.Since(tlshs)
				},
				ConnectStart: func(network, addr string) {
					connect = time.Now().UTC()
				},
				ConnectDone: func(network, addr string, err error) {
					res.Stats.ESTCXNDur = time.Since(connect)
				},
				GotFirstResponseByte: func() {
					res.Stats.TTFB = time.Since(res.Stats.ReqStartedAt)
				},
			},
		),
	)
	res.Stats.ReqStartedAt = time.Now().UTC()
	rsp, err := c.client.Transport.RoundTrip(req)

	if err != nil {
		res.applyError(nil, err)
		return nil, err
	}

	res.applyRsp(rsp)
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
