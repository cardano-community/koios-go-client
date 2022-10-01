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

package testsuite

import (
	"fmt"
	"os"

	"github.com/cardano-community/koios-go-client/v2"
)

// TestingT is an interface wrapper around *testing.T
type TestingT interface {
	Errorf(format string, args ...any)
}

func getClient() (*koios.Client, error) {
	var host string
	switch os.Getenv("KOIOS_NETWORK") {
	case "mainnet":
		host = koios.MainnetHost
	case "guild":
		host = koios.GuildHost
	case "testnet":
		host = koios.TestnetHost
	default:
		return nil, fmt.Errorf("invalid network or KOIOS_NETWORK not set: %s", os.Getenv("KOIOS_NETWORK"))
	}
	return koios.New(koios.Host(host))
}
