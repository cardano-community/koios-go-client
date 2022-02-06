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

func addNetworkCommands(app *cli.App, api *koios.Client) {
	app.Commands = append(app.Commands, []*cli.Command{
		{
			Name:     "tip",
			Category: "NETWORK",
			Usage:    "Get the tip info about the latest block seen by chain.",
			Action: func(ctx *cli.Context) error {
				res, err := api.GetTip(callctx)
				output(ctx, res, err)
				return nil
			},
		},
		{
			Name:     "genesis",
			Category: "NETWORK",
			Usage:    "Get the Genesis parameters used to start specific era on chain.",
			Action: func(ctx *cli.Context) error {
				res, err := api.GetGenesis(callctx)
				output(ctx, res, err)
				return nil
			},
		},
		{
			Name:     "totals",
			Category: "NETWORK",
			Usage:    "Get the circulating utxo, treasury, rewards, supply and reserves in lovelace for specified epoch, all epochs if empty.",
			Flags: []cli.Flag{
				&cli.UintFlag{
					Name:  "epoch-no",
					Usage: "Epoch Number to fetch details for",
				},
			},
			Action: func(ctx *cli.Context) error {
				var epochNo *koios.EpochNo
				if ctx.Uint("epoch-no") > 0 {
					v := koios.EpochNo(ctx.Uint("epoch-no"))
					epochNo = &v
				}

				res, err := api.GetTotals(callctx, epochNo)
				output(ctx, res, err)
				return nil
			},
		},
	}...)
}
