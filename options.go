// SPDX-License-Identifier: Apache-2.0
//
// Copyright © 2022 The Cardano Community Authors

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
			if port >= 80 && port != 443 {
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
			if c.locked {
				return ErrClientLocked
			}
			timeout := DefaultTimeout
			if client == nil {
				client = &http.Client{}
			}

			if client.Timeout == 0 && c.client != nil && c.client.Timeout > 0 {
				timeout = c.client.Timeout
			}
			client.Timeout = timeout

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
//
// Deprecated: Use EnableRequestsStats instead.
// CollectRequestsStats will be removed in v5.0.0
func CollectRequestsStats(enabled bool) Option {
	return EnableRequestsStats(enabled)
}

// EnableRequestsStats when enabled uses httptrace is used
// to collect detailed timing information about the request.
func EnableRequestsStats(enable bool) Option {
	return Option{
		apply: func(c *Client) error {
			c.reqStatsEnabled = enable
			return nil
		},
	}
}

func Timeout(timeout time.Duration) Option {
	return Option{
		apply: func(c *Client) error {
			if c.client == nil {
				return ErrHTTPClientNotSet
			}
			c.client.Timeout = timeout
			return nil
		},
	}
}
