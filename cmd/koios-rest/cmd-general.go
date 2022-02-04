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
	"io/ioutil"

	"github.com/howijd/koios-rest-go-client"
	"github.com/urfave/cli/v2"
)

func addGeneralCommands(app *cli.App, api *koios.Client) {
	app.Commands = append(app.Commands, []*cli.Command{
		{
			Name:  "get",
			Usage: "Send get request to API endpoint",
			Action: func(ctx *cli.Context) error {
				endpoint := ctx.Args().Get(0)
				if len(endpoint) == 0 {
					return errors.New("provide endpoint as argument e.g. /tip")
				}
				res, err := api.GET(context.Background(), endpoint)
				handleErr(err)
				defer res.Body.Close()
				body, err := ioutil.ReadAll(res.Body)
				printResponseBody(body, ctx)
				return nil
			},
		},
	}...)
}
