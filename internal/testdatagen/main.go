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
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"sync"

	"github.com/howijd/koios-rest-go-client"
	"github.com/howijd/koios-rest-go-client/internal"
	"github.com/urfave/cli/v2"
)

var (
	callctx context.Context
	cancel  context.CancelFunc
)

func main() {
	api, err := koios.New(
		koios.RateLimit(1),
	)
	handleErr(err)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	callctx, cancel = context.WithCancel(context.Background())
	go func() {
		<-c
		cancel()
	}()

	app := &cli.App{
		Authors: []*cli.Author{
			&cli.Author{
				Name: "The Howijd.Network Authors",
			},
		},
		Copyright: "(c) 2022",
		Usage:     "Generetate testdata from testnet api endpoint.",
		Before: func(c *cli.Context) error {
			handleErr(koios.Host(koios.TestnetHost)(api))
			return nil
		},
		Action: func(ctx *cli.Context) error {
			if ctx.NArg() != 1 {
				return errors.New("provide path to directory to output test data")
			}
			dirpath := ctx.Args().Get(0)

			dir, err := os.Stat(dirpath)
			handleErr(err)
			if !dir.IsDir() {
				return errors.New("path is not a directory")
			}

			var wg sync.WaitGroup

			for _, spec := range specs() {
				wg.Add(1)
				go func() {
					defer wg.Done()
					var (
						res *http.Response
						err error
					)

					log.Printf("requesting: %s\n", spec.Endpoint)
					switch spec.Request.Method {
					case "GET":
						res, err = api.GET(callctx, spec.Endpoint, spec.Request.Query, nil)
						handleErr(err)
					case "HEAD":
						res, err = api.HEAD(callctx, spec.Endpoint, spec.Request.Query, nil)
						handleErr(err)
					case "POST":
						res, err = api.POST(callctx, spec.Endpoint, spec.Request.Body, spec.Request.Query, nil)
						handleErr(err)

					}

					defer res.Body.Close()
					body, err := ioutil.ReadAll(res.Body)
					handleErr(err)
					spec.Response.Header = res.Header
					spec.Response.Code = res.StatusCode
					handleErr(json.Unmarshal(body, &spec.Response.Body))

					outfile := filepath.Join(dirpath, spec.Filename)
					log.Printf("saving: %s\n", outfile)

					out, err := json.MarshalIndent(spec, "", " ")
					handleErr(err)
					handleErr(os.WriteFile(outfile, out, 0600))
				}()
			}
			wg.Wait()
			return nil
		},
	}

	handleErr(app.Run(os.Args))
}

func handleErr(err error) {
	if err == nil {
		return
	}
	cancel()
	trace := err
	for errors.Unwrap(trace) != nil {
		trace = errors.Unwrap(trace)
		log.Println(trace)
	}
	log.Fatal(err)
}

func specs() []internal.APITestSpec {
	return []internal.APITestSpec{
		{
			Filename: "endpoint_api_tip.json",
			Endpoint: "/tip",
			Request: internal.APITestRequestSpec{
				Method: "GET",
				Body:   nil, //  bytes.NewReader(data)
			},
		},
	}
}
