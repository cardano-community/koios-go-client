// Copyright 2022 The Cardano Community Authors
// SPDX-License-Identifier: Apache-2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//
//	http://www.apache.org/licenses/LICENSE-2.0
//	or LICENSE file in repository root.
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
	"time"

	"golang.org/x/time/rate"
)

type (
	// Option is callback function to apply
	// configurations options of API Client.
	Option struct {
		apply func(*Client) error
	}
)

// Host returns option apply func which can be used to change the
// baseurl hostname https://<host>/api/v0/
func Host(host string) Option {
	return Option{
		apply: func(c *Client) error {
			if c.url.Port() == "" || c.url.Port() == "80" || c.url.Port() == "443" {
				c.url.Host = host
			} else {
				c.url.Host = fmt.Sprint(host, ":", c.url.Port())
			}
			return nil
		},
	}
}

// APIVersion returns option to apply change of the
// baseurl api version https://api.koios.rest/api/<version>/
func APIVersion(version string) Option {
	return Option{
		apply: func(c *Client) error {
			url, err := c.url.Parse("/api/" + version + "/")
			c.url = url
			return err
		},
	}
}

// Port returns option apply func which can be used to change the
// baseurl port https://api.koios.rest:<port>/api/v0/
func Port(port uint16) Option {
	return Option{
		apply: func(c *Client) error {
			if port != 80 && port != 443 {
				c.url.Host = fmt.Sprint(c.url.Hostname(), ":", port)
			}
			return nil
		},
	}
}

// Scheme returns option apply func which can be used to change the
// baseurl scheme <scheme>://api.koios.rest/api/v0/.
func Scheme(scheme string) Option {
	return Option{
		apply: func(c *Client) error {
			c.url.Scheme = scheme
			if scheme != "http" && scheme != "https" {
				return ErrSchema
			}
			return nil
		},
	}
}

// HTTPClient enables to set htt.Client to be used for requests.
func HTTPClient(client *http.Client) Option {
	return Option{
		apply: func(c *Client) error {
			if c.client != nil {
				return ErrHTTPClientChange
			}
			if client == nil {
				client = &http.Client{
					Timeout: time.Second * 60,
				}
			}
			if client.Timeout == 0 {
				return ErrHTTPClientTimeoutSetting
			}
			c.client = client
			if c.client.Transport == nil {
				t := http.DefaultTransport.(*http.Transport).Clone()
				t.MaxIdleConns = 100
				t.MaxConnsPerHost = 100
				t.MaxIdleConnsPerHost = 100
				c.client.Transport = t
			}
			return nil
		},
	}
}

// RateLimit sets requests per second this client is allowed to create
// and effectievely rate limits outgoing requests.
// Let's respect usage of the community provided resources.
func RateLimit(reqps int) Option {
	return Option{
		apply: func(c *Client) error {
			if reqps == 0 {
				return ErrRateLimitRange
			}
			c.r = rate.NewLimiter(rate.Every(time.Second), reqps)
			return nil
		},
	}
}

// Origin sets Origin header for outgoing api requests.
// Recomoended is to set it to URL or FQDN of your project using this library.
//
// In case you appliation goes rouge this could help to keep api.koios.rest
// service stable and up and running while temporary limiting requests
// it accepts from your application.
//
// It's not required, but considered as good practice so that Cardano Community
// can provide HA services for Cardano ecosystem.
func Origin(origin string) Option {
	return Option{
		apply: func(c *Client) error {
			o, err := url.ParseRequestURI(origin)
			if err != nil {
				return err
			}
			c.commonHeaders.Set("Origin", o.String())
			return nil
		},
	}
}

// CollectRequestsStats when enabled uses httptrace is used
// to collect detailed timing information about the request.
func CollectRequestsStats(enabled bool) Option {
	return Option{
		apply: func(c *Client) error {
			c.reqStatsEnabled = enabled
			return nil
		},
	}
}
