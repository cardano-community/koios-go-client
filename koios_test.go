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

package koios_test

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/cardano-community/koios-go-client/v2"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

// TestingT is an interface wrapper around *testing.T
type TestingT interface {
	Errorf(format string, args ...any)
}

func networkEpoch() koios.EpochNo {
	var epoch koios.EpochNo
	switch os.Getenv("KOIOS_NETWORK") {
	case "guild":
		epoch = koios.EpochNo(1950)
	case "testnet":
		epoch = koios.EpochNo(185)
	default:
		// mainnet
		epoch = koios.EpochNo(320)
	}
	return epoch
}

func networkBlockHash() koios.BlockHash {
	var hash koios.BlockHash
	switch os.Getenv("KOIOS_NETWORK") {
	case "guild":
		hash = koios.BlockHash("af2f6f7dd4e4ea6765103a1e38e023da3edd2b3c7fea2aa367222564dbe01cfd")
	case "testnet":
		hash = koios.BlockHash("f75fea40852ed7d7f539d008e45255725daef8553ae7162750836f279570813a")
	default:
		// mainnet
		hash = koios.BlockHash("fb9087c9f1408a7bbd7b022fd294ab565fec8dd3a8ef091567482722a1fa4e30")
	}
	return hash
}

func getClient() (client *koios.Client, err error) {
	net, ok := os.LookupEnv("KOIOS_NETWORK")
	if !ok {
		return nil, errors.New("KOIOS_NETWORK not set")
	}
	var host string
	switch net {
	case "mainnet":
		host = koios.MainnetHost
	case "guild":
		host = koios.GuildHost
	case "testnet":
		host = koios.TestnetHost
	default:
		return nil, fmt.Errorf("invalid KOIOS_NETWORK=%q", net)
	}
	return koios.New(koios.Host(host))
}

func assertEqual[V comparable](t TestingT, want, got V, tag string) {
	msg := fmt.Sprintf("%s: want(%v) got(%v)", tag, want, got)
	assert.Equal(t, want, got, msg)
}

func assertIsPositive(t TestingT, in decimal.Decimal, tag string) {
	msg := fmt.Sprintf("%s(should be positice): got  %s", tag, in.String())
	assert.True(t, in.IsPositive(), msg)
}

func assertGreater[V any](t TestingT, want, got V, tag string) {
	msg := fmt.Sprintf("%s: val(%v) should be greater than %v", tag, got, want)
	assert.Greater(t, want, got, msg)
}

func assertNotEmpty(t TestingT, in any, tag string) {
	msg := fmt.Sprintf("%s: in(%v)", tag, in)
	assert.NotEmpty(t, in, msg)
}

func assertTimeNotZero(t TestingT, in time.Time, tag string) {
	msg := fmt.Sprintf("%s: time is empty or not parsed from return value", tag)
	assert.False(t, in.IsZero(), msg)
}
