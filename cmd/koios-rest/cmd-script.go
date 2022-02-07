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

func addScriptCommands(app *cli.App, api *koios.Client) {
	app.Commands = append(app.Commands, []*cli.Command{
		{
			Name:     "script-list",
			Category: "SCRIPT",
			Usage:    "List of all existing script hashes along with their creation transaction hashes.",
			Action: func(ctx *cli.Context) error {
				res, err := api.GetScriptList(callctx)
				output(ctx, res, err)
				return nil
			},
		},
		{
			Name:      "script-redeemers",
			Category:  "SCRIPT",
			Usage:     "List of all redeemers for a given script hash.",
			ArgsUsage: "[script_hash]",
			Action: func(ctx *cli.Context) error {
				if ctx.NArg() != 1 {
					return errors.New("script-redeemers requires single script-hash as arg")
				}
				res, err := api.GetScriptRedeemers(callctx, koios.ScriptHash(ctx.Args().Get(0)))
				output(ctx, res, err)
				return nil
			},
		},
	}...)
}
