// SPDX-License-Identifier: Apache-2.0
//
// Copyright © 2022 The Cardano Community Authors

// Package koios provides api client library to interact with Koios API
// endpoints and Cardano Blockchain. Sub package ./cmd/koios-rest
// provides cli application.
// Koios is best described as a Decentralized and Elastic RESTful query layer
// for exploring data on Cardano blockchain to consume within
// applications/wallets/explorers/etc.
package koios // imports as package "koios"

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/happy-sdk/happy/pkg/version"
	"github.com/shopspring/decimal"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const (
	// MainnetHost is primay and default api host.
	MainnetHost = "api.koios.rest"
	// MainnetHostEU is main net api host in EU.
	MainnetHostEU = "eu-api.koios.rest"
	// GuildnetHost is Guild network host.
	GuildHost = "guild.koios.rest"
	// PreviewHost is Preview network host.
	PreviewHost = "preview.koios.rest"
	// PreProdHost is Pre Production network host.
	PreProdHost = "preprod.koios.rest"
	// DefaultAPIVersion is openapi spec version e.g. /v1.
	DefaultAPIVersion = "v1"
	// DefaultPort default port used by api client.
	DefaultPort uint16 = 443
	// DefaultSchema default schema used by api client.
	DefaultScheme = "https"
	// DefaultRateLimit is default rate limit used by api client.
	DefaultRateLimit int = 10 // https://api.koios.rest/#overview--limits
	// DefaultOrigin is default origin header used by api client.
	DefaultOrigin = "https://github.com/cardano-community/koios-go-client/v4"
	// PageSize is default page size used by api client.
	PageSize uint = 1000
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
	ErrNoDatumHash              = errors.New("missing datum hash(es)")
	ErrNoAddress                = errors.New("missing address")
	ErrNoPoolID                 = errors.New("missing pool id")
	ErrResponse                 = errors.New("response error")
	ErrSchema                   = errors.New("scheme must be http or https")
	ErrReqOptsAlreadyUsed       = errors.New("request options can only be used once")
	ErrUnexpectedResponseField  = errors.New("unexpected response field")
	ErrUTxOInputAlreadyUsed     = errors.New("UTxO already used")
	ErrNoData                   = errors.New("no data")

	// ZeroLovelace is alias decimal.Zero.
	ZeroLovelace = decimal.Zero.Copy() //nolint: gochecknoglobals
	// ZeroCoin is alias decimal.Zero.
	ZeroCoin = decimal.Zero.Copy() //nolint: gochecknoglobals
)

type (
	Slot uint

	BlockNo uint

	// PaymentCredential type def.
	PaymentCredential string

	// BlockHash defines type for _block_hash.
	BlockHash string

	// TxHash defines type for tx_hash.
	TxHash string

	// PoolID type def.
	PoolID string

	// PolicyID type def.
	PolicyID string

	// Timestamp extends time to work with unix timestamps and
	// fix time format anomalies when Unmarshaling and Marshaling
	// Koios API times.
	Timestamp struct {
		time.Time
	}

	// PaymentAddr info.
	PaymentAddr struct {
		// Bech32 is Cardano payment/base address (bech32 encoded)
		// for transaction's or change to be returned.
		Bech32 Address `json:"bech32"`

		// Payment credential.
		Cred PaymentCredential `json:"cred"`
	}

	// Certificate information.
	Certificate struct {
		// Index of the certificate
		Index int `json:"index"`

		// Info is A JSON object containing information from the certificate.
		Info map[string]json.RawMessage `json:"info"`

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
)

// New creates thread-safe API client you can freerly pass this
// client to multiple go routines.
//
// Call to New without options is same as call with default options.
// e.g.
// api, err := koios.New(
//
//	koios.Host(koios.MainnetHost),
//	koios.APIVersion(koios.DefaultAPIVersion),
//	koios.Port(koios.DefaultPort),
//	koios.Schema(koios.DefaultSchema),
//	koios.HttpClient(koios.DefaultHttpClient),
//
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
			Version(),
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

	if c.commonHeaders.Get("Origin") == "" {
		// Sets default origin if option was not provided.
		_ = Origin(DefaultOrigin).apply(c)
	}

	// If HttpClient option was not provided
	// use default http.Client
	if c.client == nil {
		// there is really no point to check that error
		_ = HTTPClient(nil).apply(c)
	}

	return c, nil
}

// Version returns koios go library version
func Version() string {
	return libraryVersion
}

var libraryVersion string

func init() {
	v := version.Current()
	libraryVersion = v.String()
}

// String returns PaymentCredential as string.
func (v PaymentCredential) String() string {
	return string(v)
}

// String returns BlockHash as string.
func (v BlockHash) String() string {
	return string(v)
}

// String returns TxHash as string.
func (v TxHash) String() string {
	return strings.Trim(string(v), "\"")
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

func (t *Timestamp) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		return nil
	}
	q, err := strconv.ParseInt(string(b), 10, 64)
	if err != nil {
		return nil
	}
	t.Time = time.Unix(q, 0)
	return nil
}

// MarshalJSON turns our time.Time back into an int.
func (t Timestamp) MarshalJSON() ([]byte, error) {
	if t.IsZero() {
		return []byte("null"), nil
	}
	return []byte(strconv.FormatInt(t.Time.Unix(), 10)), nil
}

func (b BlockNo) String() string {
	return fmt.Sprintf("%d", b)
}

func (s Slot) String() string {
	return fmt.Sprintf("%d", s)
}

func (h BlockHash) MarshalJSON() ([]byte, error) {
	if len(h) == 0 {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("%q", h)), nil
}
