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

func addEpochCommands(app *cli.App, api *koios.Client) {
	app.Commands = append(app.Commands, []*cli.Command{
		{
			Name:     "epoch-info",
			Category: "EPOCH",
			Usage:    "Get the epoch information, all epochs if no epoch specified.",
			Flags: []cli.Flag{
				&cli.Uint64Flag{
					Name:  "epoch",
					Usage: "Epoch Number to fetch details for",
					Value: uint64(0),
				},
			},
			Action: func(ctx *cli.Context) error {
				var epoch *koios.EpochNo
				if ctx.Uint("epoch") > 0 {
					v := koios.EpochNo(ctx.Uint64("epoch"))
					epoch = &v
				}

				res, err := api.GetEpochInfo(callctx, epoch)
				output(ctx, res, err)
				return nil
			},
		},
		{
			Name:     "epoch-params",
			Category: "EPOCH",
			Usage:    "Get the protocol parameters for specific epoch, returns information about all epochs if no epoch specified.",
			Flags: []cli.Flag{
				&cli.Uint64Flag{
					Name:  "epoch",
					Usage: "Epoch Number to fetch details for",
					Value: uint64(0),
				},
			},
			Action: func(ctx *cli.Context) error {
				var epoch *koios.EpochNo
				if ctx.Uint("epoch") > 0 {
					v := koios.EpochNo(ctx.Uint64("epoch"))
					epoch = &v
				}

				res, err := api.GetEpochParams(callctx, epoch)
				output(ctx, res, err)
				return nil
			},
		},
	}...)
}
