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

package main

import (
	"github.com/howijd/koios-rest-go-client"
	"github.com/urfave/cli/v2"
)

func addPoolCommands(app *cli.App, api *koios.Client) {
	app.Commands = append(app.Commands, []*cli.Command{
		{
			Name:     "pool-list",
			Category: "POOL",
			Usage:    "A list of all currently registered/retiring (not retired) pools.",
		},
		{
			Name:     "pool-info",
			Category: "POOL",
			Usage:    "Current pool statuses and details for a specified list of pool ids.",
		},
		{
			Name:     "pool-delegators",
			Category: "POOL",
			Usage:    "Return information about delegators by a given pool and optional epoch (current if omitted).",
		},
		{
			Name:     "pool-blocks",
			Category: "POOL",
			Usage:    "Return information about blocks minted by a given pool in current epoch (or _epoch_no if provided).",
		},
		{
			Name:     "pool-updates",
			Category: "POOL",
			Usage:    "Return all pool updates for all pools or only updates for specific pool if specified.",
		},
		{
			Name:     "pool-relays",
			Category: "POOL",
			Usage:    "A list of registered relays for all currently registered/retiring (not retired) pools.",
		},
		{
			Name:     "pool-metadata",
			Category: "POOL",
			Usage:    "Metadata(on & off-chain) for all currently registered/retiring (not retired) pools.",
		},
	}...)
}
