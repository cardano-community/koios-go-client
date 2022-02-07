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
	"errors"

	"github.com/howijd/koios-rest-go-client"
	"github.com/urfave/cli/v2"
)

func addPoolCommands(app *cli.App, api *koios.Client) {
	app.Commands = append(app.Commands, []*cli.Command{
		{
			Name:     "pool-list",
			Category: "POOL",
			Usage:    "A list of all currently registered/retiring (not retired) pools.",
			Action: func(ctx *cli.Context) error {
				res, err := api.GetPoolList(callctx)
				output(ctx, res, err)
				return nil
			},
		},
		{
			Name:      "pool-infos",
			Category:  "POOL",
			Usage:     "Current pool statuses and details for a specified list of pool ids.",
			ArgsUsage: "[pool-id...]",
			Action: func(ctx *cli.Context) error {
				var pids []koios.PoolID
				for _, pid := range ctx.Args().Slice() {
					pids = append(pids, koios.PoolID(pid))
				}
				res, err := api.GetPoolInfos(callctx, pids)
				output(ctx, res, err)
				return nil
			},
		},
		{
			Name:      "pool-info",
			Category:  "POOL",
			Usage:     "Current pool status and details for a specified pool by pool id.",
			ArgsUsage: "[pool-id]",
			Action: func(ctx *cli.Context) error {
				if ctx.NArg() != 1 {
					return errors.New("pool-info requires single pool id")
				}
				res, err := api.GetPoolInfo(callctx, koios.PoolID(ctx.Args().Get(0)))
				output(ctx, res, err)
				return nil
			},
		},
		{
			Name:     "pool-delegators",
			Category: "POOL",
			Usage:    "Return information about delegators by a given pool and optional epoch (current if omitted).",
			Flags: []cli.Flag{
				&cli.Uint64Flag{
					Name:  "epoch",
					Usage: "Epoch Number to fetch details for",
					Value: uint64(0),
				},
			},
			Action: func(ctx *cli.Context) error {
				if ctx.NArg() != 1 {
					return errors.New("pool-delegators requires single pool id")
				}
				var epoch *koios.EpochNo
				if ctx.Uint("epoch") > 0 {
					v := koios.EpochNo(ctx.Uint64("epoch"))
					epoch = &v
				}

				res, err := api.GetPoolDelegators(callctx, koios.PoolID(ctx.Args().Get(0)), epoch)
				output(ctx, res, err)
				return nil
			},
		},
		{
			Name:     "pool-blocks",
			Category: "POOL",
			Usage:    "Return information about blocks minted by a given pool in current epoch (or _epoch_no if provided).",
			Flags: []cli.Flag{
				&cli.Uint64Flag{
					Name:  "epoch",
					Usage: "Epoch Number to fetch details for",
					Value: uint64(0),
				},
			},
			Action: func(ctx *cli.Context) error {
				if ctx.NArg() != 1 {
					return errors.New("pool-blocks requires single pool id")
				}
				var epoch *koios.EpochNo
				if ctx.Uint("epoch") > 0 {
					v := koios.EpochNo(ctx.Uint64("epoch"))
					epoch = &v
				}

				res, err := api.GetPoolBlocks(callctx, koios.PoolID(ctx.Args().Get(0)), epoch)
				output(ctx, res, err)
				return nil
			},
		},
		{
			Name:     "pool-updates",
			Category: "POOL",
			Usage:    "Return all pool updates for all pools or only updates for specific pool if specified.",
			Action: func(ctx *cli.Context) error {
				var pool *koios.PoolID
				if ctx.NArg() == 1 {
					v := koios.PoolID(ctx.Args().Get(0))
					pool = &v
				}

				res, err := api.GetPoolUpdates(callctx, pool)
				output(ctx, res, err)
				return nil
			},
		},
		{
			Name:     "pool-relays",
			Category: "POOL",
			Usage:    "A list of registered relays for all currently registered/retiring (not retired) pools.",
			Action: func(ctx *cli.Context) error {
				res, err := api.GetPoolRelays(callctx)
				output(ctx, res, err)
				return nil
			},
		},
		{
			Name:     "pool-metadata",
			Category: "POOL",
			Usage:    "Metadata(on & off-chain) for all currently registered/retiring (not retired) pools.",
			Action: func(ctx *cli.Context) error {
				res, err := api.GetPoolMetadata(callctx)
				output(ctx, res, err)
				return nil
			},
		},
	}...)
}
