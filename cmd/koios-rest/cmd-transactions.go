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
	"context"
	"errors"

	"github.com/howijd/koios-rest-go-client"
	"github.com/urfave/cli/v2"
)

func addTransactionsCommands(app *cli.App, api *koios.Client) {
	app.Commands = append(app.Commands, []*cli.Command{
		{
			Name:      "tx-infos",
			Category:  "TRANSACTIONS",
			Usage:     "Get detailed information about transaction(s).",
			ArgsUsage: "[tx-hashes...]",
			Action: func(ctx *cli.Context) error {
				var txs []koios.TxHash
				for _, a := range ctx.Args().Slice() {
					txs = append(txs, koios.TxHash(a))
				}
				res, err := api.GetTxsInfos(context.Background(), txs)
				output(ctx, res, err)
				return nil
			},
		},
		{
			Name:      "tx-info",
			Category:  "TRANSACTIONS",
			Usage:     "Get detailed information about single transaction.",
			ArgsUsage: "[tx-hashes]",
			Action: func(ctx *cli.Context) error {
				if ctx.NArg() != 1 {
					return errors.New("tx-info requires single transaction hash")
				}
				res, err := api.GetTxsInfo(context.Background(), koios.TxHash(ctx.Args().Get(0)))
				output(ctx, res, err)
				return nil
			},
		},
		{
			Name:     "tx-utxos",
			Category: "TRANSACTIONS",
			Usage:    "Get UTxO set (inputs/outputs) of transactions.",
		},
		{
			Name:     "tx-metadata",
			Category: "TRANSACTIONS",
			Usage:    "Get metadata information (if any) for given transaction(s)..",
		},
		{
			Name:     "tx-metalabels",
			Category: "TRANSACTIONS",
			Usage:    "Get a list of all transaction metalabels.",
		},
		{
			Name:     "submittx",
			Category: "TRANSACTIONS",
			Usage:    "Submit an already serialized transaction to the network.",
		},
		{
			Name:     "tx-status",
			Category: "TRANSACTIONS",
			Usage:    "Get the number of block confirmations for a given transaction hash list",
		},
	}...)
}
