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
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"

	"github.com/cardano-community/koios-go-client"
	"github.com/cardano-community/koios-go-client/internal"
	"github.com/urfave/cli/v2"
)

var (
	callctx context.Context
	cancel  context.CancelFunc
)

const TestEpoch = "294"

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

					var wg sync.WaitGroup
					for _, filestat := range filestats {
						if filepath.Ext(filestat.Name()) == ".gz" &&
							(filename == "all" || filename+".gz" == filestat.Name()) {
							wg.Add(1)
							go func(filestat fs.FileInfo) {
								log.Print("reading: ", filestat.Name())
								defer wg.Done()
								gzfilename := filepath.Join(dirpath, filestat.Name())

								gzfile, err := os.Open(gzfilename)
								defer gzfile.Close()

								handleErr(err)
								gzr, err := gzip.NewReader(gzfile)
								handleErr(err)
								specb, err := io.ReadAll(gzr)
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
							)

							log.Println("requesting: ", spec.Endpoint)
							switch spec.Request.Method {
							case "GET":
								res, err = api.GET(callctx, spec.Endpoint, spec.Request.Query, spec.Request.Header)
								handleErr(err)
							case "HEAD":
								res, err = api.HEAD(callctx, spec.Endpoint, spec.Request.Query, spec.Request.Header)
								handleErr(err)
							case "POST":
								res, err = api.POST(callctx, spec.Endpoint, bytes.NewReader(spec.Request.Body), spec.Request.Query, spec.Request.Header)
								handleErr(err)

							}

							defer res.Body.Close()
							body, err := io.ReadAll(res.Body)
							handleErr(err)
							spec.Response.Header = res.Header
							spec.Response.Code = res.StatusCode
							handleErr(json.Unmarshal(body, &spec.Response.Body))

							outfile := filepath.Join(dirpath, spec.Filename+".gz")
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
			Filename: "endpoint_network_tip.json",
			Endpoint: "/tip",
			Request: internal.APITestRequestSpec{
				Method: "GET",
			},
		},
		{
			Filename: "endpoint_network_genesis.json",
			Endpoint: "/genesis",
			Request: internal.APITestRequestSpec{
				Method: "GET",
			},
		},
		{
			Filename: "endpoint_network_totals.json",
			Endpoint: "/totals",
			Request: internal.APITestRequestSpec{
				Query: url.Values{
					"_epoch_no": []string{TestEpoch},
				},
				Method: "GET",
			},
		},
		{
			Filename: "endpoint_epoch_info.json",
			Endpoint: "/epoch_info",
			Request: internal.APITestRequestSpec{
				Query: url.Values{
					"_epoch_no": []string{TestEpoch},
				},
				Method: "GET",
			},
		},
		{
			Filename: "endpoint_epoch_params.json",
			Endpoint: "/epoch_params",
			Request: internal.APITestRequestSpec{
				Query: url.Values{
					"_epoch_no": []string{TestEpoch},
				},
				Method: "GET",
			},
		},
		{
			Filename: "endpoint_blocks.json",
			Endpoint: "/blocks",
			Request: internal.APITestRequestSpec{
				Method: "GET",
			},
		},
		{
			Filename: "endpoint_block_info.json",
			Endpoint: "/block_info",
			Request: internal.APITestRequestSpec{
				Method: "GET",
				Query: url.Values{
					"_block_hash": []string{"f6192a1aaa6d3d05b4703891a6b66cd757801c61ace86cbe5ab0d66e07f601ab"},
				},
			},
		},
		{
			Filename: "endpoint_block_txs.json",
			Endpoint: "/block_txs",
			Request: internal.APITestRequestSpec{
				Method: "GET",
				Query: url.Values{
					"_block_hash": []string{"f6192a1aaa6d3d05b4703891a6b66cd757801c61ace86cbe5ab0d66e07f601ab"},
				},
			},
		},
		{
			Filename: "endpoint_tx_info.json",
			Endpoint: "/tx_info",
			Request: internal.APITestRequestSpec{
				Method: "POST",
				Body:   []byte("{\"_tx_hashes\": [\"f144a8264acf4bdfe2e1241170969c930d64ab6b0996a4a45237b623f1dd670e\"]}"),
			},
		},
		{
			Filename: "endpoint_tx_utxos.json",
			Endpoint: "/tx_utxos",
			Request: internal.APITestRequestSpec{
				Method: "POST",
				Body:   []byte("{\"_tx_hashes\": [\"f144a8264acf4bdfe2e1241170969c930d64ab6b0996a4a45237b623f1dd670e\"]}"),
			},
		},
		{
			Filename: "endpoint_tx_metadata.json",
			Endpoint: "/tx_metadata",
			Request: internal.APITestRequestSpec{
				Method: "POST",
				Body:   []byte("{\"_tx_hashes\": [\"f144a8264acf4bdfe2e1241170969c930d64ab6b0996a4a45237b623f1dd670e\"]}"),
			},
		},
		{
			Filename: "endpoint_tx_metalabels.json",
			Endpoint: "/tx_metalabels",
			Request: internal.APITestRequestSpec{
				Method: "GET",
			},
		},
		{
			Filename: "endpoint_tx_status.json",
			Endpoint: "/tx_status",
			Request: internal.APITestRequestSpec{
				Method: "POST",
				Body:   []byte("{\"_tx_hashes\": [\"f144a8264acf4bdfe2e1241170969c930d64ab6b0996a4a45237b623f1dd670e\"]}"),
			},
		},
		{
			Filename: "endpoint_tx_submit.json",
			Endpoint: "/submittx",
			Request: internal.APITestRequestSpec{
				Header: http.Header{
					"Content-Type":   []string{"application/cbor"},
					"Content-Length": []string{"585"},
				},
				Method: "POST",
				Body:   []byte("{\"type\":\"Tx AlonzoEra\",\"description\":\"\",\"cborHex\":\"84a60081825820bf9b23cdd9bff2b1a802da7b527a0c6dd0378efa73c0800e8875f9c37930f7ef010d800182825839011f56a82c4c006289171fced204a37a2806e15c88a98872ef9626d3ddc5e778ead6d4d614c64ec8475c8b3dee4d2b8613fa1f3adee95581151a001e848082581d61e1eabc77c631f9dffa24b4c938bf09458d384764ede698d13bb3957f1a00563386021a0002acfd031a0322b0aa0e80a10081825820112bb18afb7f33b90ad1be59accfc7bcc4784c47fde6a5a10d2c932119df16bb584033642286d7805776288655000e2cebbac069def2e1735b91fa53fc5e6650b5921d54c5c5492dc97d8dce9e3539691ca4e45ae9ed4573f6d691adac8aae345001f5f6\"}"),
			},
		},
		{
			Filename: "endpoint_address_info.json",
			Endpoint: "/address_info",
			Request: internal.APITestRequestSpec{
				Method: "GET",
				Query: url.Values{
					"_address": []string{"addr1qyp9kz50sh9c53hpmk3l4ewj9ur794t2hdqpngsjn3wkc5sztv9glpwt3frwrhdrltjaytc8ut2k4w6qrx3p98zad3fq07xe9g"},
				},
			},
		},
		{
			Filename: "endpoint_address_txs.json",
			Endpoint: "/address_txs",
			Request: internal.APITestRequestSpec{
				Method: "POST",
				Body:   []byte("{\"_addresses\": [\"addr1qyp9kz50sh9c53hpmk3l4ewj9ur794t2hdqpngsjn3wkc5sztv9glpwt3frwrhdrltjaytc8ut2k4w6qrx3p98zad3fq07xe9g\"], \"_after_block_height\": 6238675}"),
			},
		},
		{
			Filename: "endpoint_address_assets.json",
			Endpoint: "/address_assets",
			Request: internal.APITestRequestSpec{
				Method: "GET",
				Query: url.Values{
					"_address": []string{"addr1q8d5kc3u4lcu84g08apa9ckj7df605lahzzpcfy0m5tpkyk9uauw44k56c2vvnkggawgk00wf54cvyl6ruada624sy2se0snsj"},
				},
			},
		},
		{
			Filename: "endpoint_credential_txs.json",
			Endpoint: "/credential_txs",
			Request: internal.APITestRequestSpec{
				Method: "POST",
				Body:   []byte("{\"_payment_credentials\": [\"025b0a8f85cb8a46e1dda3fae5d22f07e2d56abb4019a2129c5d6c52\"], \"_after_block_height\": 6238675}"),
			},
		},
		{
			Filename: "endpoint_account_list.json",
			Endpoint: "/account_list",
			Request: internal.APITestRequestSpec{
				Method: "GET",
			},
		},
		{
			Filename: "endpoint_account_info.json",
			Endpoint: "/account_info",
			Request: internal.APITestRequestSpec{
				Method: "GET",
				Query: url.Values{
					"_address": []string{"stake1u8yxtugdv63wxafy9d00nuz6hjyyp4qnggvc9a3vxh8yl0ckml2uz"},
				},
			},
		},
		{
			Filename: "endpoint_account_rewards.json",
			Endpoint: "/account_rewards",
			Request: internal.APITestRequestSpec{
				Method: "GET",
				Query: url.Values{
					"_stake_address": []string{"stake1u8yxtugdv63wxafy9d00nuz6hjyyp4qnggvc9a3vxh8yl0ckml2uz"},
					"_epoch_no":      []string{TestEpoch},
				},
			},
		},
		{
			Filename: "endpoint_account_updates.json",
			Endpoint: "/account_updates",
			Request: internal.APITestRequestSpec{
				Method: "GET",
				Query: url.Values{
					"_stake_address": []string{"stake1u8yxtugdv63wxafy9d00nuz6hjyyp4qnggvc9a3vxh8yl0ckml2uz"},
				},
			},
		},
		{
			Filename: "endpoint_account_addresses.json",
			Endpoint: "/account_addresses",
			Request: internal.APITestRequestSpec{
				Method: "GET",
				Query: url.Values{
					"_address": []string{"stake1u8yxtugdv63wxafy9d00nuz6hjyyp4qnggvc9a3vxh8yl0ckml2uz"},
				},
			},
		},
		{
			Filename: "endpoint_account_assets.json",
			Endpoint: "/account_assets",
			Request: internal.APITestRequestSpec{
				Method: "GET",
				Query: url.Values{
					"_address": []string{"stake1u8yxtugdv63wxafy9d00nuz6hjyyp4qnggvc9a3vxh8yl0ckml2uz"},
				},
			},
		},
		{
			Filename: "endpoint_account_history.json",
			Endpoint: "/account_history",
			Request: internal.APITestRequestSpec{
				Method: "GET",
				Query: url.Values{
					"_address": []string{"stake1u8yxtugdv63wxafy9d00nuz6hjyyp4qnggvc9a3vxh8yl0ckml2uz"},
				},
			},
		},
		{
			Filename: "endpoint_asset_list.json",
			Endpoint: "/asset_list",
			Request: internal.APITestRequestSpec{
				Method: "GET",
			},
		},
		{
			Filename: "endpoint_asset_address_list.json",
			Endpoint: "/asset_address_list",
			Request: internal.APITestRequestSpec{
				Method: "GET",
				Query: url.Values{
					"_asset_policy": []string{"d3501d9531fcc25e3ca4b6429318c2cc374dbdbcf5e99c1c1e5da1ff"},
					"_asset_name":   []string{"444f4e545350414d"},
				},
			},
		},
		{
			Filename: "endpoint_asset_info.json",
			Endpoint: "/asset_info",
			Request: internal.APITestRequestSpec{
				Method: "GET",
				Query: url.Values{
					"_asset_policy": []string{"d3501d9531fcc25e3ca4b6429318c2cc374dbdbcf5e99c1c1e5da1ff"},
					"_asset_name":   []string{"444f4e545350414d"},
				},
			},
		},
		{
			Filename: "endpoint_asset_summary.json",
			Endpoint: "/asset_summary",
			Request: internal.APITestRequestSpec{
				Method: "GET",
				Query: url.Values{
					"_asset_policy": []string{"d3501d9531fcc25e3ca4b6429318c2cc374dbdbcf5e99c1c1e5da1ff"},
					"_asset_name":   []string{"444f4e545350414d"},
				},
			},
		},
		{
			Filename: "endpoint_asset_txs.json",
			Endpoint: "/asset_txs",
			Request: internal.APITestRequestSpec{
				Method: "GET",
				Query: url.Values{
					"_asset_policy": []string{"d3501d9531fcc25e3ca4b6429318c2cc374dbdbcf5e99c1c1e5da1ff"},
					"_asset_name":   []string{"444f4e545350414d"},
				},
			},
		},
		{
			Filename: "endpoint_pool_list.json",
			Endpoint: "/pool_list",
			Request: internal.APITestRequestSpec{
				Method: "GET",
			},
		},
		{
			Filename: "endpoint_pool_info.json",
			Endpoint: "/pool_info",
			Request: internal.APITestRequestSpec{
				Method: "POST",
				Body:   []byte("{\"_pool_bech32_ids\": [\"pool100wj94uzf54vup2hdzk0afng4dhjaqggt7j434mtgm8v2gfvfgp\"]}"),
			},
		},
		{
			Filename: "endpoint_pool_delegators.json",
			Endpoint: "/pool_delegators",
			Request: internal.APITestRequestSpec{
				Method: "GET",
				Query: url.Values{
					"_pool_bech32": []string{"pool155efqn9xpcf73pphkk88cmlkdwx4ulkg606tne970qswczg3asc"},
					"_epoch_no":    []string{TestEpoch},
				},
			},
		},
		{
			Filename: "endpoint_pool_blocks.json",
			Endpoint: "/pool_blocks",
			Request: internal.APITestRequestSpec{
				Method: "GET",
				Query: url.Values{
					"_pool_bech32": []string{"pool155efqn9xpcf73pphkk88cmlkdwx4ulkg606tne970qswczg3asc"},
					"_epoch_no":    []string{TestEpoch},
				},
			},
		},
		{
			Filename: "endpoint_pool_history.json",
			Endpoint: "/pool_history",
			Request: internal.APITestRequestSpec{
				Method: "GET",
				Query: url.Values{
					"_pool_bech32": []string{"pool155efqn9xpcf73pphkk88cmlkdwx4ulkg606tne970qswczg3asc"},
					"_epoch_no":    []string{TestEpoch},
				},
			},
		},
		{
			Filename: "endpoint_pool_updates.json",
			Endpoint: "/pool_updates",
			Request: internal.APITestRequestSpec{
				Method: "GET",
				Query: url.Values{
					"_pool_bech32": []string{"pool155efqn9xpcf73pphkk88cmlkdwx4ulkg606tne970qswczg3asc"},
				},
			},
		},
		{
			Filename: "endpoint_pool_relays.json",
			Endpoint: "/pool_relays",
			Request: internal.APITestRequestSpec{
				Method: "GET",
			},
		},
		{
			Filename: "endpoint_pool_metadata.json",
			Endpoint: "/pool_metadata",
			Request: internal.APITestRequestSpec{
				Method: "GET",
			},
		},
		{
			Filename: "endpoint_native_script_list.json",
			Endpoint: "/native_script_list",
			Request: internal.APITestRequestSpec{
				Method: "GET",
			},
		},
		{
			Filename: "endpoint_plutus_script_list.json",
			Endpoint: "/plutus_script_list",
			Request: internal.APITestRequestSpec{
				Method: "GET",
			},
		},
		{
			Filename: "endpoint_script_redeemers.json",
			Endpoint: "/script_redeemers",
			Request: internal.APITestRequestSpec{
				Method: "GET",
				Query: url.Values{
					"_script_hash": []string{"d8480dc869b94b80e81ec91b0abe307279311fe0e7001a9488f61ff8"},
				},
			},
		},
		{
			Filename: "endpoint_asset_policy_info.json",
			Endpoint: "/asset_policy_info",
			Request: internal.APITestRequestSpec{
				Method: "GET",
				Query: url.Values{
					"_asset_policy": []string{"a8102151506a8a81dc1763ee05cdd01d787f50dfeb6f843071e1c6a0"},
				},
			},
		},
		{
			Filename: "endpoint_asset_history.json",
			Endpoint: "/asset_history",
			Request: internal.APITestRequestSpec{
				Method: "GET",
				Query: url.Values{
					"_asset_policy": []string{"d3501d9531fcc25e3ca4b6429318c2cc374dbdbcf5e99c1c1e5da1ff"},
					"_asset_name":   []string{"444f4e545350414d"},
				},
			},
		},
	}
}
