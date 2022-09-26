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
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"sync"

	"github.com/cardano-community/koios-go-client/v2"
	"github.com/cardano-community/koios-go-client/v2/internal"
	"github.com/go-openapi/loads"
)

var (
	ErrInvalidCommand = errors.New("unknown command, valid commands [unpack|generate]")
	ErrMissingCommand = errors.New("missing command, valid commands [unpack|generate]")
	ErrInvalidNetwork = errors.New("invalid network, valid networks [all|mainnet|testnet|guildnet]")
	ErrTestID         = errors.New("invalid test id(s)")

	callctx context.Context
	cancel  context.CancelFunc

	mainnet = NetworkDefaults{
		EpochNo: koios.EpochNo(320),
	}
	testnet = NetworkDefaults{
		EpochNo: koios.EpochNo(185),
	}
	guildnet = NetworkDefaults{
		EpochNo: koios.EpochNo(1950),
	}
)

type NetworkDefaults struct {
	EpochNo koios.EpochNo
}

const (
	MainnetEpoch  = "349"
	TestnetEpoch  = "216"
	GuildnetEpoch = "1950"
)

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	callctx, cancel = context.WithCancel(context.Background())
	go func() {
		<-c
		cancel()
	}()

	var (
		cmd       string
		args      []string
		networkID string
	)

	generateCmd := flag.NewFlagSet("generate", flag.ExitOnError)
	generateCmd.Usage = func() {
		fmt.Printf(" %s [flags|arg] optional arg endpoint id (default: all).\n\n", generateCmd.Name())
		generateCmd.PrintDefaults()
		fmt.Println("")
	}
	generateCmd.StringVar(&networkID, "n", "all", "network(s) to use [all|mainnet|testnet|guildnet]")

	// unpackCmd := flag.NewFlagSet("unpack", flag.ExitOnError)
	// unpackCmd.Usage = func() {
	// 	fmt.Printf(" %s [arg] optional arg endpoint id (default: all).\n\n", unpackCmd.Name())
	// 	unpackCmd.PrintDefaults()
	// 	fmt.Println("")
	// }

	if len(os.Args) >= 2 {
		cmd = os.Args[1]
		args = os.Args[2:]
	}

	switch cmd {
	// case "unpack":
	// 	handleErr(unpackCmd.Parse(os.Args[2:]))
	// 	handleErr(unpack(networkID, args))
	case "generate":
		handleErr(generateCmd.Parse(os.Args[2:]))
		handleErr(generate(networkID, args))
	default:
		fmt.Print("TESTDATA GENERATOR HELP\n" +
			"The Cardano Community Authors (c) 2022\n\n" +
			"Helper program to generate testdata for Koios Go Client.\n\n" +
			"USAGE\n\n")
		generateCmd.Usage()
		// unpackCmd.Usage()
	}
}

func generate(networkID string, args []string) error {
	networks, err := useNetorks(networkID)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup

	for _, networkID := range networks {
		fmt.Println("")
		log.Println("NETWORK: ", networkID)

		wg.Add(1)
		go func() {
			defer wg.Done()

			api, err := getClient(networkID)
			if err != nil {
				logErr(err)
				return
			}

			if err := updateAPISpec(networkID, api); err != nil {
				logErr(err)
				return
			}

			if err := genTestdata(context.WithValue(callctx, "networkID", networkID), api, args); err != nil {
				logErr(err)
				return
			}

		}()
		wg.Wait()
	}

	return nil
}

func genTestdata(ctx context.Context, api *koios.Client, args []string) error {

	networkID := fmt.Sprint(ctx.Value("networkID"))
	specs, err := specs(networkID, args)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	for _, spec := range specs {
		gzfilename := filepath.Join("../../testdata", networkID, spec.Filename+".json.gz")
		jsonfilename := filepath.Join("../../testdata", networkID, spec.Filename+".json")
		if err := os.Remove(gzfilename); err != nil && !errors.Is(err, fs.ErrNotExist) {
			logErr(err)
			continue
		}
		if err := os.Remove(jsonfilename); err != nil && !errors.Is(err, fs.ErrNotExist) {
			logErr(err)
			continue
		}

		log.Printf("update specs(%s): %s - %s\n", networkID, spec.Filename, spec.Endpoint)
		wg.Add(1)
		go func(gzfilename, jsonfilename string, spec internal.APITestSpec) {
			defer wg.Done()
			var (
				res *http.Response
				err error
			)
			log.Printf("requesting: %s - %s", spec.Network, spec.Endpoint)
			opts := api.NewRequestOptions()
			opts.HeadersApply(spec.Request.Header)
			opts.QueryApply(spec.Request.Query)
			switch spec.Request.Method {
			case "GET":
				res, err = api.GET(callctx, spec.Endpoint, opts)
				if err != nil {
					logSpecErr(networkID, spec.Filename, fmt.Errorf("%w: %s - %s", err, spec.Network, spec.Endpoint))
					return
				}
			case "HEAD":
				res, err = api.HEAD(callctx, spec.Endpoint, opts)
				if err != nil {
					logSpecErr(networkID, spec.Filename, fmt.Errorf("%w: %s - %s", err, spec.Network, spec.Endpoint))
					return
				}
			case "POST":
				res, err = api.POST(callctx, spec.Endpoint, bytes.NewReader(spec.Request.Body), opts)
				if err != nil {
					logSpecErr(networkID, spec.Filename, fmt.Errorf("%w: %s - %s", err, spec.Network, spec.Endpoint))
					return
				}
			}

			defer res.Body.Close()
			body, err := io.ReadAll(res.Body)
			if err != nil {
				logSpecErr(networkID, spec.Filename, err)
				return
			}
			spec.Response.Header = res.Header
			spec.Response.Code = res.StatusCode
			if err := json.Unmarshal(body, &spec.Response.Body); err != nil {
				logSpecErr(networkID, spec.Filename, err)
				return
			}
			var jsongz bytes.Buffer
			z := gzip.NewWriter(&jsongz)
			jsonbytes, err := json.Marshal(spec)
			if err != nil {
				logSpecErr(networkID, spec.Filename, err)
				return
			}
			jsonbytesp, err := json.MarshalIndent(spec, "", "  ")
			if err != nil {
				logSpecErr(networkID, spec.Filename, err)
				return
			}
			if _, err := z.Write(jsonbytes); err != nil {
				logSpecErr(networkID, spec.Filename, err)
				return
			}
			z.Close()
			logSpecErr(networkID, spec.Filename, os.WriteFile(gzfilename, jsongz.Bytes(), 0644))
			logSpecErr(networkID, spec.Filename, os.WriteFile(jsonfilename, jsonbytesp, 0644))

		}(gzfilename, jsonfilename, spec)
	}
	wg.Wait()

	return nil
}

func updateAPISpec(networkID string, client *koios.Client) error {
	src := client.ServerURL().ResolveReference(&url.URL{Path: "/koiosapi.yaml"}).String()
	dest := filepath.Join("../../testdata", networkID, "koiosapi.json.gz")
	jsondest := filepath.Join("../../testdata", networkID, "koiosapi.json")

	log.Println("update spec(src):", src)
	doc, err := loads.Spec(src)
	if err != nil {
		return err
	}

	var jsongz bytes.Buffer
	z := gzip.NewWriter(&jsongz)

	jsonout, err := json.Marshal(doc.Raw())
	if err != nil {
		return err
	}

	if _, err := z.Write(jsonout); err != nil {
		return err
	}

	if err := os.WriteFile(dest, jsongz.Bytes(), 0600); err != nil {
		return err
	}
	log.Println("wrote spec(dest):", dest)

	if err := os.WriteFile(jsondest, jsonout, 0600); err != nil {
		return err
	}
	log.Println("wrote spec(json):", jsondest)

	return nil
}

func getClient(networkID string) (*koios.Client, error) {
	var host string
	switch networkID {
	case "mainnet":
		host = koios.MainnetHost
	case "testnet":
		host = koios.TestnetHost
	case "guildnet":
		host = koios.GuildnetHost
	}
	return koios.New(koios.Host(host))
}

func unpack(networkID string, args []string) error {
	networks, err := useNetorks(networkID)
	if err != nil {
		return err
	}
	var wg sync.WaitGroup

	for _, network := range networks {
		log.Println("network: ", network)
		wg.Add(1)
		go func() {
			defer wg.Done()
		}()
		wg.Wait()
	}

	return nil
}

func useNetorks(networkID string) ([]string, error) {
	log.Println("use networks", networkID)
	if networkID == "all" {
		return []string{"mainnet", "testnet", "guildnet"}, nil
	}
	switch networkID {
	case "mainnet", "testnet", "guildnet":
		return []string{networkID}, nil
	}
	return nil, fmt.Errorf("%w: %s", ErrInvalidNetwork, networkID)
}

func handleErr(err error) {
	if err == nil {
		return
	}
	cancel()
	log.Fatal("ERROR: ", err)
}

func logErr(err error) {
	if err == nil {
		return
	}
	log.Println("ERROR: ", err)
}

func logSpecErr(networkID, specID string, err error) {
	if err == nil {
		return
	}
	log.Printf("ERROR(%s-%s): %s", networkID, specID, err.Error())
}

func specs(networkID string, testids []string) ([]internal.APITestSpec, error) {
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

	if (networkID == "all") && (len(testids) == 0 || testids[0] == "all") {
		return specs, nil
	}

	var usespecs []internal.APITestSpec
	for _, spec := range specs {
		if len(testids) > 0 {
			for _, testid := range testids {
				if (networkID == "all" || networkID == spec.Network) && testid == spec.Filename {
					usespecs = append(usespecs, spec)
				}
			}
		} else {
			if networkID == "all" || networkID == spec.Network {
				usespecs = append(usespecs, spec)
			}
		}
	}
	if len(usespecs) == 0 {
		return nil, fmt.Errorf("%w: %s - %v", ErrTestID, networkID, testids)
	}
	return usespecs, nil
}
