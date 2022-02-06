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

func addBlockCommands(app *cli.App, api *koios.Client) {
	app.Commands = append(app.Commands, []*cli.Command{
		{
			Name:     "blocks",
			Category: "BLOCK",
			Usage:    "Get summarised details about all blocks (paginated - latest first).",
			Action: func(ctx *cli.Context) error {
				res, err := api.GetBlocks(callctx)
				output(ctx, res, err)
				return nil
			},
		},
		{
			Name:     "block-info",
			Category: "BLOCK",
			Usage:    "Get detailed information about a specific block.",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "block-hash",
					Usage:    "Block Hash in hex format to fetch details for",
					Required: true,
				},
			},
			Action: func(ctx *cli.Context) error {
				res, err := api.GetBlockInfo(
					callctx,
					koios.BlockHash(ctx.String("block-hash")),
				)
				output(ctx, res, err)
				return nil
			},
		},
		{
			Name:     "block-txs",
			Category: "BLOCK",
			Usage:    "Get a list of all transactions included in a provided block.",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "block-hash",
					Usage:    "Block Hash in hex format to fetch details for",
					Required: true,
				},
			},
			Action: func(ctx *cli.Context) error {
				res, err := api.GetBlockTxHashes(
					callctx,
					koios.BlockHash(ctx.String("block-hash")),
				)
				output(ctx, res, err)
				return nil
			},
		},
	}...)
}
