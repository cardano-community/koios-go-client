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

package koios_test

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"sync"
	"time"

	koios "github.com/howijd/koios-rest-go-client"
)

func ExampleRateLimit() {
	api, err := koios.New(
		// limit client request 1 per second
		koios.RateLimit(1),
	)
	handleErr(err)

	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Second*10,
	)
	defer cancel()

	var wg sync.WaitGroup

	servers := []string{
		"api.koios.rest",
		"guild.koios.rest",
		"testnet.koios.rest",
	}

	for _, host := range servers {
		wg.Add(1)
		go func(host string) {
			defer wg.Done()
			// switching host
			koios.Host(host)(api)
			res, err := api.GET(ctx, "/tip")
			handleErr(err)
			defer res.Body.Close()
			body, err := ioutil.ReadAll(res.Body)
			fmt.Println(string(body))
			handleErr(err)
		}(host)
	}

	wg.Wait()
	fmt.Println("requests processed: ", api.TotalRequests())
}

func handleErr(err error) {
	if err == nil {
		return
	}
	trace := err
	for errors.Unwrap(trace) != nil {
		trace = errors.Unwrap(trace)
		log.Println(trace)
	}
	log.Fatal(err)
}
