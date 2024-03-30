// SPDX-License-Identifier: Apache-2.0
//
// Copyright © 2022 The Cardano Community Authors

package koios

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"time"
)

func GetTokenAuthInfo(jwt string) (AuthInfo, error) {
	auth, err := decodeJWT(jwt)
	if err != nil {
		fmt.Printf("Error decoding JWT: %v\n", err)
		return AuthInfo{}, fmt.Errorf("%w: error decoding JWT %v", ErrAuth, err)
	}
	return *auth, nil
}

func (c *Client) SetAuth(jwt string) error {
	auth, err := GetTokenAuthInfo(jwt)
	if err != nil {
		return err
	}
	c.auth = &auth
	return nil
}

func (c *Client) getAuth() AuthInfo {
	if c.auth == nil {
		return AuthInfo{}
	}
	auth := *c.auth
	return auth
}

type AuthInfo struct {
	// AuthType is the type of authentication used.
	Tier    AuthTier    `json:"tier"`
	ProjID  string      `json:"projID"`
	Addr    string      `json:"addr"`
	Expires AuthExpires `json:"expires"`

	MaxRequests     uint          `json:"max_requests"`
	MaxRPS          uint          `json:"max_rps"`
	MaxQueryTimeout time.Duration `json:"query_timeout"`
	CORSRestricted  bool          `json:"cors_restricted"`
	token           string
}

func (a *AuthInfo) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Tier    string      `json:"tier"`
		ProjID  string      `json:"projID,omitempty"`
		Addr    string      `json:"addr,omitempty"`
		Expires AuthExpires `json:"expires"`

		MaxRequests     uint   `json:"max_requests"`
		MaxRPS          uint   `json:"max_rps"`
		MaxQueryTimeout string `json:"query_timeout"`
		CORSRestricted  bool   `json:"cors_restricted"`
	}{
		Tier:            a.Tier.String(),
		MaxRequests:     a.Tier.MaxRequest(),
		MaxRPS:          a.Tier.MaxRPS(),
		MaxQueryTimeout: a.Tier.MaxQueryTimeout().String(),
		CORSRestricted:  a.Tier.CORSRestricted(),
		Expires:         a.Expires,
		Addr:            a.Addr,
		ProjID:          a.ProjID,
	})
}

type AuthTier uint8

const (
	// AutTierPublic gives everyone Free access to our API services so everyone is ready to start building.
	AutTierPublic = iota
	// AutTierFree Once familiar with our platform, consider our Advanced solution for "production ready" deployment.
	AutTierFree
	// AuthTierPro Enhanced resource capabilities that enable you to develop your application with ease.
	AuthTierPro
	// AuthTierPremium Top-notch resource capabilities, allowing you to relax and not worry about resource consumption.
	AuthTierPremium
	// AuthTierCustom A tailor made infrastructure design around your requirements, for optimal resource management, giving you peace of mind.
	AuthTierCustom
)

func (at AuthTier) String() string {
	switch at {
	case AutTierPublic:
		return "public"
	case AutTierFree:
		return "free"
	case AuthTierPro:
		return "pro"
	case AuthTierPremium:
		return "premium"
	case AuthTierCustom:
		return "custom"
	default:
		return "unknown"
	}
}

func (at AuthTier) MaxRequest() uint {
	switch at {
	case AutTierPublic:
		return 5000
	case AutTierFree:
		return 50000
	case AuthTierPro:
		return 500000
	case AuthTierPremium:
		return 1200000
	case AuthTierCustom:
		return math.MaxUint
	default:
		return 0
	}
}

func (at AuthTier) MaxRPS() uint {
	switch at {
	case AutTierPublic, AutTierFree:
		return 10
	case AuthTierPro:
		return 25
	case AuthTierPremium:
		return 500
	case AuthTierCustom:
		return math.MaxUint
	default:
		return 0
	}
}

func (at AuthTier) MaxQueryTimeout() time.Duration {
	switch at {
	case AutTierPublic, AutTierFree:
		return 30 * time.Second
	case AuthTierPro:
		return time.Minute
	case AuthTierPremium:
		return 2 * time.Minute
	case AuthTierCustom:
		return time.Duration(math.MaxInt64)
	default:
		return 0
	}
}

func (at AuthTier) CORSRestricted() bool {
	switch at {
	case AutTierPublic:
		return true
	default:
		return false
	}
}

func (at AuthTier) MarshalJSON() ([]byte, error) {
	return []byte(`"` + at.String() + `"`), nil
}

type AuthExpires time.Time

func (ae AuthExpires) MarshalJSON() ([]byte, error) {
	if time.Time(ae).IsZero() {
		return []byte("\"No expiration date\""), nil
	}
	return json.Marshal(ae.String())
}

func (ae AuthExpires) Time() time.Time {
	return time.Time(ae)
}

func (ae AuthExpires) String() string {
	return time.Time(ae).String()
}

// decodeBase64URL decodes a Base64URL string.
func decodeBase64URL(s string) (string, error) {
	b, err := base64.RawURLEncoding.DecodeString(s)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// decodeJWT decodes the payload of a JWT.
func decodeJWT(token string) (*AuthInfo, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid JWT format")
	}

	decoded, err := decodeBase64URL(parts[1])
	if err != nil {
		return nil, fmt.Errorf("failed to decode payload: %w", err)
	}

	var payload struct {
		Addr   string    `json:"addr"`
		Tier   AuthTier  `json:"tier"`
		ProjID string    `json:"projID"`
		Exp    Timestamp `json:"exp"`
	}
	err = json.Unmarshal([]byte(decoded), &payload)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	var auth AuthInfo
	auth.Addr = payload.Addr
	auth.Tier = payload.Tier
	auth.ProjID = payload.ProjID
	auth.Expires = AuthExpires(payload.Exp.Time)
	auth.token = token

	auth.MaxRequests = payload.Tier.MaxRequest()
	auth.MaxRPS = payload.Tier.MaxRPS()
	auth.MaxQueryTimeout = payload.Tier.MaxQueryTimeout()

	return &auth, nil
}
