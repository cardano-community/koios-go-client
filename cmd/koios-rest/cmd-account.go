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

func addAccountCommands(app *cli.App, api *koios.Client) {
	app.Commands = append(app.Commands, []*cli.Command{
		{
			Name:     "account-list",
			Category: "ACCOUNT",
			Usage:    "Get a list of all accounts returns array of stake addresses.",
			Action: func(ctx *cli.Context) error {
				res, err := api.GetAccountList(callctx)
				output(ctx, res, err)
				return nil
			},
		},
		{
			Name:      "account-info",
			Category:  "ACCOUNT",
			Usage:     "Get the account info of any (payment or staking) address.",
			ArgsUsage: "[account]",
			Action: func(ctx *cli.Context) error {
				if ctx.NArg() != 1 {
					return errors.New("account-info requires single address")
				}
				res, err := api.GetAccountInfo(callctx, koios.Address(ctx.Args().Get(0)))
				output(ctx, res, err)
				return nil
			},
		},
		{
			Name:     "account-rewards",
			Category: "ACCOUNT",
			Usage:    "Get the full rewards history (including MIR) for a stake address, or certain epoch if specified.",
		},
		{
			Name:     "account-updates",
			Category: "ACCOUNT",
			Usage:    "Get the account updates (registration, deregistration, delegation and withdrawals).",
		},
		{
			Name:     "account-addresses",
			Category: "ACCOUNT",
			Usage:    "Get all addresses associated with an account.",
		},
		{
			Name:     "account-assets",
			Category: "ACCOUNT",
			Usage:    "Get the native asset balance of an account.",
		},
		{
			Name:     "account-history",
			Category: "ACCOUNT",
			Usage:    "Get the staking history of an account.",
		},
	}...)
}
