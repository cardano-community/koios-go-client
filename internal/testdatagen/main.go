// Copyright 2022 The Cardano Community Authors
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
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"

	"github.com/cardano-community/koios-go-client/v2"
	"github.com/cardano-community/koios-go-client/v2/internal"
	"github.com/urfave/cli/v2"
)

var (
	callctx context.Context
	cancel  context.CancelFunc
)

const (
	MainnetEpoch = "349"
	TestnetEpoch = "216"
)

func main() {
	mainnet, err := koios.New(
		koios.Host(koios.MainnetHost),
	)
	handleErr(err)

	testnet, err := koios.New(
		koios.Host(koios.TestnetHost),
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
			{
				Name: "The Cardano Community Authors",
			},
		},
		Copyright: "(c) 2022",
		Usage:     "Generetate testdata from testnet api endpoint.",
		Before: func(c *cli.Context) error {
			// handleErr(koios.Host(koios.TestnetHost)(api))
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:      "unpack",
				Usage:     "unpack test data.",
				ArgsUsage: "[testdata dir path]",
				Action: func(ctx *cli.Context) error {
					if ctx.NArg() == 0 {
						return errors.New("provide path to directory to output test data")
					}
					dirpath := ctx.Args().Get(0)
					filename := ctx.Args().Get(1)

					dir, err := os.Stat(dirpath)
					handleErr(err)
					if !dir.IsDir() {
						return errors.New("path is not a directory")
					}
					filestats, err := ioutil.ReadDir(dirpath)
					if err != nil {
						return err
					}
					var wg sync.WaitGroup
					for _, filestat := range filestats {
						if filepath.Ext(filestat.Name()) == ".gz" &&
							(filename == "all" || filename+".json.gz" == filestat.Name()) {
							wg.Add(1)
							go func(filestat fs.FileInfo) {
								log.Print("reading: ", filestat.Name())
								defer wg.Done()
								gzfilename := filepath.Join(dirpath, filestat.Name())

								gzfile, err := os.Open(gzfilename)
								handleErr(err)
								defer func() {
									if err := gzfile.Close(); err != nil {
										log.Println(err)
									}
								}()

								gzr, err := gzip.NewReader(gzfile)
								handleErr(err)
								specb, err := io.ReadAll(gzr)
								handleErr(err)
								gzr.Close()

								spec := &internal.APITestSpec{}
								jsonfile := strings.TrimRight(gzfilename, ".gz")

								if err := json.Unmarshal(specb, spec); err != nil {
									log.Println("failed: ", jsonfile)
									log.Println("> ", string(specb))
									log.Fatal(err)
								}

								log.Println("saving: ", jsonfile)
								out, err := json.MarshalIndent(spec, "", " ")
								handleErr(err)
								handleErr(os.WriteFile(jsonfile, out, 0644))
							}(filestat)
						}
					}
					wg.Wait()
					return nil
				},
			},
			{
				Name:      "generate",
				Usage:     "generate or update test data.",
				ArgsUsage: "[testdata dir path]",
				Action: func(ctx *cli.Context) error {
					if ctx.NArg() == 0 {
						return errors.New("provide path to directory to output test data")
					}
					dirpath := ctx.Args().Get(0)
					filename := ctx.Args().Get(1)

					dir, err := os.Stat(dirpath)
					handleErr(err)
					if !dir.IsDir() {
						return errors.New("path is not a directory")
					}

					var wg sync.WaitGroup

					for _, spec := range specs() {
						if filename != "all" && filename != spec.Filename {
							continue
						}

						wg.Add(1)
						go func(spec internal.APITestSpec) {

							defer wg.Done()
							var (
								res *http.Response
								err error
								api *koios.Client
							)

							log.Printf("requesting: %s - %s", spec.Network, spec.Endpoint)
							switch spec.Network {
							case "testnet":
								api = testnet
							default:
								api = mainnet
							}

							opts := api.NewRequestOptions()

							opts.HeadersApply(spec.Request.Header)
							opts.QueryApply(spec.Request.Query)

							switch spec.Request.Method {
							case "GET":
								res, err = api.GET(callctx, spec.Endpoint, opts)
								if err != nil {
									handleErr(fmt.Errorf("%w: %s - %s", err, spec.Network, spec.Endpoint))
								}
							case "HEAD":
								res, err = api.HEAD(callctx, spec.Endpoint, opts)
								if err != nil {
									handleErr(fmt.Errorf("%w: %s - %s", err, spec.Network, spec.Endpoint))
								}
							case "POST":
								res, err = api.POST(callctx, spec.Endpoint, bytes.NewReader(spec.Request.Body), opts)
								if err != nil {
									handleErr(fmt.Errorf("%w: %s - %s", err, spec.Network, spec.Endpoint))
								}
							}

							defer res.Body.Close()
							body, err := io.ReadAll(res.Body)
							handleErr(err)
							spec.Response.Header = res.Header
							spec.Response.Code = res.StatusCode
							handleErr(json.Unmarshal(body, &spec.Response.Body))

							outfile := filepath.Join(dirpath, spec.Network, spec.Filename+".json.gz")
							_ = os.Remove(filepath.Join(dirpath, spec.Filename))
							_ = os.Remove(outfile)

							var jsongz bytes.Buffer
							z := gzip.NewWriter(&jsongz)

							gzout, err := json.Marshal(spec)
							handleErr(err)
							_, err = z.Write(gzout)
							handleErr(err)
							z.Close()
							handleErr(os.WriteFile(outfile, jsongz.Bytes(), 0644))
							log.Println("generating: ", outfile)
						}(spec)
					}
					wg.Wait()
					return nil
				},
			},
		},
	}

	handleErr(app.Run(os.Args))
}

func handleErr(err error) {
	if err == nil {
		return
	}
	log.Fatal(err)
	cancel()
}

func specs() []internal.APITestSpec {
	var specs []internal.APITestSpec
	specs = append(specs, accountEndpointSpecs()...)
	specs = append(specs, addressEndpointSpecs()...)
	specs = append(specs, assetsEndpointSpecs()...)
	specs = append(specs, blocksEndpointSpecs()...)
	specs = append(specs, epochEndpointSpecs()...)
	specs = append(specs, networkEndpointSpecs()...)
	specs = append(specs, poolsEndpointSpecs()...)
	specs = append(specs, scriptsEndpointSpecs()...)
	specs = append(specs, txEndpointSpecs()...)
	specs = append(specs, filtersSpecs()...)
	return specs
}
