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

// Package koios provides api client library to interact with Koios API
// endpoints and Cardano Blockchain. Sub package ./cmd/koios-rest
// provides cli application.
//
// Koios is best described as a Decentralized and Elastic RESTful query layer
// for exploring data on Cardano blockchain to consume within
// applications/wallets/explorers/etc.
package koios // imports as package "koios"

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"runtime"
	"time"

	"github.com/shopspring/decimal"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"golang.org/x/time/rate"
)

// MainnetHost       : is primay and default api host.
// GuildHost         : is Guild network host.
// TestnetHost       : is api host for testnet.
// DefaultAPIVersion : is openapi spec version e.g. /v0.
// DefaultPort       : default port used by api client.
// DefaultSchema     : default schema used by api client.
// LibraryVersion    : koios go library version.
// DefaultRateLimit  : is default rate limit used by api client.
// DefaultOrigin     : is default origin header used by api client.
const (
	MainnetHost              = "api.koios.rest"
	GuildHost                = "guild.koios.rest"
	TestnetHost              = "testnet.koios.rest"
	DefaultAPIVersion        = "v0"
	DefaultPort       uint16 = 443
	DefaultScheme            = "https"
	LibraryVersion           = "v0"
	DefaultRateLimit  int    = 10 // https://api.koios.rest/#overview--limits
	DefaultOrigin            = "https://github.com/cardano-community/koios-go-client"
	PageSize          uint   = 1000
)

// Predefined errors used by the library.
var (
	ErrURLValuesLenght          = errors.New("if presenent then only single url.Values should be provided")
	ErrHTTPClientTimeoutSetting = errors.New("http.Client.Timeout should never be 0 in production")
	ErrHTTPClientChange         = errors.New("http.Client can only be set as option to koios.New")
	ErrOriginSet                = errors.New("origin can only be set as option to koios.New")
	ErrRateLimitRange           = errors.New("rate limit must be between 1-255 requests per sec")
	ErrResponseIsNotJSON        = errors.New("got non json response")
	ErrNoTxHash                 = errors.New("missing transaxtion hash(es)")
	ErrNoAddress                = errors.New("missing address")
	ErrNoPoolID                 = errors.New("missing pool id")
	ErrResponse                 = errors.New("got unexpected response")
	ErrSchema                   = errors.New("scheme must be http or https")
	ErrReqOptsAlreadyUsed       = errors.New("request options can only be used once")
	ErrUnexpectedResponseField  = errors.New("unexpected response field")
	ZeroLovelace                = NewLovelace(0, 1)
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

	// Option is callback function to apply
	// configurations options of API Client.
	Option struct {
		apply func(*Client) error
	}

	// Address defines type for _address.
	Address string

	// PaymentCredential type def.
	PaymentCredential string

	// AssetName defines type for _asset_name.
	AssetName string

	// BlockHash defines type for _block_hash.
	BlockHash string

	// TxHash defines type for tx_hash.
	TxHash string

	// EpochNo defines type for _epoch_no.
	EpochNo uint64

	// PoolID type def.
	PoolID string

	// PolicyID type def.
	PolicyID string

	// ScriptHash defines type for _script_hash.
	ScriptHash string

	// StakeAddress is Cardano staking address (reward account, bech32 encoded).
	StakeAddress string

	// Time extends time to fix time format anomalies turing Unmarshal and Marshal
	Time struct {
		time.Time
	}

	// Lovelace defines type for ADA lovelaces. This library uses forked snapshot
	// of github.com/shopspring/decimal package to provide. JSON and XML
	// serialization/deserialization and make it ease to work with calculations
	// and deciimal precisions of ADA lovelace and native assets.
	//
	// For API of decimal package see
	// https://pkg.go.dev/github.com/shopspring/decimal
	Lovelace struct {
		decimal.Decimal
	}

	// PaymentAddr info.
	PaymentAddr struct {
		// Bech32 is Cardano payment/base address (bech32 encoded)
		// for transaction's or change to be returned.
		Bech32 string `json:"bech32"`

		// Payment credential.
		Cred PaymentCredential `json:"cred"`
	}

	// Certificate information.
	Certificate struct {
		// Index of the certificate
		Index int `json:"index"`

		// Info is A JSON object containing information from the certificate.
		Info map[string]interface{} `json:"info"`

		// Type of certificate could be:
		// delegation, stake_registration, stake_deregistraion, pool_update,
		// pool_retire, param_proposal, reserve_MIR, treasury_MIR).
		Type string `json:"type"`
	}

	// Response wraps API responses.
	Response struct {
		// RequestURL is full request url.
		RequestURL string `json:"request_url"`

		// RequestMethod is HTTP method used for request.
		RequestMethod string `json:"request_method"`

		// StatusCode of the HTTP response.
		StatusCode int `json:"status_code"`

		// Status of the HTTP response header if present.
		Status string `json:"status"`

		// Date response header.
		Date string `json:"date,omitempty"`

		// ContentLocation response header if present.
		ContentLocation string `json:"content_location,omitempty"`

		// ContentRange response header if present.
		ContentRange string `json:"content_range,omitempty"`

		// Error response body if present.
		Error *ResponseError `json:"error,omitempty"`

		// Stats of the request if stats are enabled.
		Stats *RequestStats `json:"stats,omitempty"`
	}

	// RequestOptions for the request.
	RequestOptions struct {
		page     uint
		pageSize uint
		locked   bool
		query    url.Values
		headers  http.Header
	}

	// RequestStats represent collected request stats if collecting
	// request stats is enabled.
	RequestStats struct {
		// ReqStartedAt time when request was started.
		ReqStartedAt time.Time `json:"req_started_at,omitempty"`

		// DNSLookupDur DNS lookup duration.
		DNSLookupDur time.Duration `json:"req_dns_lookup_dur,omitempty"`

		// TLSHSDur time it took to perform TLS handshake.
		TLSHSDur time.Duration `json:"tls_hs_dur,omitempty"`

		// ESTCXNDur time it took to establish connection.
		ESTCXNDur time.Duration `json:"est_cxn_dur,omitempty"`

		// TTFB time it took to get the first byte of the response
		// after connextion was established.
		TTFB time.Duration `json:"ttfb,omitempty"`

		// ReqDur total time it took to peform the request.
		ReqDur time.Duration `json:"req_dur,omitempty"`

		// ReqDurStr String representation of ReqDur.
		ReqDurStr string `json:"req_dur_str,omitempty"`
	}

	// ResponseError represents api error messages.
	ResponseError struct {
		// Hint of the error reported by server.
		Hint string `json:"hint,omitempty"`

		// Details of the error reported by server.
		Details string `json:"details,omitempty"`

		// Code is error code reported by server.
		Code string `json:"code,omitempty"`

		// Message is error message reported by server.
		Message string `json:"message,omitempty"`
	}
)

// New creates thread-safe API client you can freerly pass this
// client to multiple go routines.
//
// Call to New without options is same as call with default options.
// e.g.
// api, err := koios.New(
// 	koios.Host(koios.MainnetHost),
// 	koios.APIVersion(koios.DefaultAPIVersion),
// 	koios.Port(koios.DefaultPort),
// 	koios.Schema(koios.DefaultSchema),
// 	koios.HttpClient(koios.DefaultHttpClient),
// ).
func New(opts ...Option) (*Client, error) {
	c := &Client{
		commonHeaders: make(http.Header),
	}
	// set default base url
	_ = c.setBaseURL(DefaultScheme, MainnetHost, DefaultAPIVersion, DefaultPort)

	// set default common headers
	c.commonHeaders.Set("Accept", "application/json")
	c.commonHeaders.Set("Accept-Encoding", "gzip, deflate")
	c.commonHeaders.Set(
		"User-Agent",
		fmt.Sprintf(
			"go-koios/%s (%s %s) %s/%s https://github.com/cardano-community/go-koios",
			LibraryVersion,
			cases.Title(language.English).String(runtime.GOOS),
			runtime.GOARCH,
			runtime.GOOS,
			runtime.GOARCH,
		),
	)

	// Apply provided options
	for _, opt := range opts {
		if err := opt.apply(c); err != nil {
			return nil, err
		}
	}

	if c.r == nil {
		// set default rate limit for outgoing requests if not configured.
		_ = RateLimit(DefaultRateLimit).apply(c)
	}

	// Sets default origin if option was not provided.
	_ = Origin(DefaultOrigin).apply(c)

	// If HttpClient option was not provided
	// use default http.Client
	if c.client == nil {
		// there is really no point to check that error
		_ = HTTPClient(nil).apply(c)
	}

	return c, nil
}

// String returns Address as string.
func (v Address) String() string {
	return string(v)
}

// String returns PaymentCredential as string.
func (v PaymentCredential) String() string {
	return string(v)
}

// String returns AssetName as string.
func (v AssetName) String() string {
	return string(v)
}

// String returns BlockHash as string.
func (v BlockHash) String() string {
	return string(v)
}

// String returns TxHash as string.
func (v TxHash) String() string {
	return string(v)
}

// String returns EpochNo as string.
func (v EpochNo) String() string {
	return fmt.Sprintf("%d", v)
}

// String returns PoolID as string.
func (v PoolID) String() string {
	return string(v)
}

// String returns PolicyID as string.
func (v PolicyID) String() string {
	return string(v)
}

// String returns ScriptHash as string.
func (v ScriptHash) String() string {
	return string(v)
}

// String returns StakeAddress as string.
func (v StakeAddress) String() string {
	return string(v)
}

func (t *Time) UnmarshalJSON(b []byte) error {
	str := string(b)
	p, err := time.Parse("\""+time.RFC3339+"\"", str)
	if err != nil {
		p, err = time.Parse("\"2006-01-02T15:04:05\"", str)
	}
	t.Time = p
	return err
}
