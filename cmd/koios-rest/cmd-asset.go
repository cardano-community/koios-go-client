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

func addAssetCommands(app *cli.App, api *koios.Client) {
	app.Commands = append(app.Commands, []*cli.Command{
		{
			Name:     "asset-list",
			Category: "ASSET",
			Usage:    "Get the list of all native assets (paginated).",
			Action: func(ctx *cli.Context) error {
				res, err := api.GetAssetList(callctx)
				output(ctx, res, err)
				return nil
			},
		},
		{
			Name:     "asset-address-list",
			Category: "ASSET",
			Usage:    "Get the list of all addresses holding a given asset.",
		},
		{
			Name:     "asset-info",
			Category: "ASSET",
			Usage:    "Get the information of an asset including first minting & token registry metadata.",
		},
		{
			Name:     "asset-summary",
			Category: "ASSET",
			Usage:    "Get the summary of an asset (total transactions exclude minting/total wallets include only wallets with asset balance).",
		},
		{
			Name:     "asset-txs",
			Category: "ASSET",
			Usage:    "Get the list of all asset transaction hashes (newest first).",
		},
	}...)
}
