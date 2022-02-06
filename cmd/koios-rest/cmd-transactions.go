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
	"encoding/json"
	"errors"
	"io/ioutil"

	"github.com/howijd/koios-rest-go-client"
	"github.com/urfave/cli/v2"
)

func addTransactionsCommands(app *cli.App, api *koios.Client) {
	app.Commands = append(app.Commands, []*cli.Command{
		{
			Name:      "txs-infos",
			Category:  "TRANSACTIONS",
			Usage:     "Get detailed information about transaction(s).",
			ArgsUsage: "[tx-hashes...]",
			Action: func(ctx *cli.Context) error {
				var txs []koios.TxHash
				for _, a := range ctx.Args().Slice() {
					txs = append(txs, koios.TxHash(a))
				}
				res, err := api.GetTxsInfos(callctx, txs)
				output(ctx, res, err)
				return nil
			},
		},
		{
			Name:      "tx-info",
			Category:  "TRANSACTIONS",
			Usage:     "Get detailed information about single transaction.",
			ArgsUsage: "[tx-hash]",
			Action: func(ctx *cli.Context) error {
				if ctx.NArg() != 1 {
					return errors.New("tx-info requires single transaction hash")
				}
				res, err := api.GetTxInfo(callctx, koios.TxHash(ctx.Args().Get(0)))
				output(ctx, res, err)
				return nil
			},
		},
		{
			Name:      "tx-utxos",
			Category:  "TRANSACTIONS",
			Usage:     "Get UTxO set (inputs/outputs) of transactions.",
			ArgsUsage: "[tx-hashes...]",
			Action: func(ctx *cli.Context) error {
				var txs []koios.TxHash
				for _, a := range ctx.Args().Slice() {
					txs = append(txs, koios.TxHash(a))
				}
				res, err := api.GetTxsUTxOs(callctx, txs)
				output(ctx, res, err)
				return nil
			},
		},
		{
			Name:      "txs-metadata",
			Category:  "TRANSACTIONS",
			ArgsUsage: "[tx-hashes...]",
			Usage:     "Get metadata information (if any) for given transaction(s).",
			Action: func(ctx *cli.Context) error {
				var txs []koios.TxHash
				for _, a := range ctx.Args().Slice() {
					txs = append(txs, koios.TxHash(a))
				}
				res, err := api.GetTxsMetadata(callctx, txs)
				output(ctx, res, err)
				return nil
			},
		},
		{
			Name:      "tx-metadata",
			Category:  "TRANSACTIONS",
			ArgsUsage: "[tx-hash]",
			Usage:     "Get metadata information (if any) for given transaction.",
			Action: func(ctx *cli.Context) error {
				if ctx.NArg() != 1 {
					return errors.New("tx-metadata requires single transaction hash")
				}
				res, err := api.GetTxMetadata(callctx, koios.TxHash(ctx.Args().Get(0)))
				output(ctx, res, err)
				return nil
			},
		},
		{
			Name:     "tx-metalabels",
			Category: "TRANSACTIONS",
			Usage:    "Get a list of all transaction metalabels.",
			Action: func(ctx *cli.Context) error {
				res, err := api.GetTxMetaLabels(callctx)
				output(ctx, res, err)
				return nil
			},
		},
		{
			Name:     "tx-submit",
			Category: "TRANSACTIONS",
			Usage:    "Submit signed transaction to the network.",
			Action: func(ctx *cli.Context) error {
				if ctx.NArg() != 1 {
					return errors.New("submittx requires single arg (path to file tx.signed)")
				}

				stx := koios.TxBodyJSON{}

				txfile, err := ioutil.ReadFile(ctx.Args().Get(0))
				if err != nil {
					return err
				}

				if err = json.Unmarshal(txfile, &stx); err != nil {
					return err
				}
				res, err := api.SubmitSignedTx(callctx, stx)
				output(ctx, res, err)
				return nil
			},
		},
		{
			Name:      "txs-statuses",
			Category:  "TRANSACTIONS",
			Usage:     "Get the number of block confirmations for a given transaction hash list",
			ArgsUsage: "[tx-hashes...]",
			Action: func(ctx *cli.Context) error {
				var txs []koios.TxHash
				for _, a := range ctx.Args().Slice() {
					txs = append(txs, koios.TxHash(a))
				}
				res, err := api.GetTxsStatuses(callctx, txs)
				output(ctx, res, err)
				return nil
			},
		},
		{
			Name:      "tx-status",
			Category:  "TRANSACTIONS",
			Usage:     "Get the number of block confirmations for a given transaction hash",
			ArgsUsage: "[tx-hash]",
			Action: func(ctx *cli.Context) error {
				if ctx.NArg() != 1 {
					return errors.New("tx-status requires single transaction hash")
				}
				res, err := api.GetTxStatus(callctx, koios.TxHash(ctx.Args().Get(0)))
				output(ctx, res, err)
				return nil
			},
		},
	}...)
}
