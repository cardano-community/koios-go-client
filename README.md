<h1>Koios API Client Library for Go</h1>

**[Koios API] is Elastic Cardano Query Layer!**

> A consistent query layer for developers to build upon Cardano, with multiple, redundant endpoints that allow for easy scalability.

**Repository contains**

1. **[Koios API] Client Library for Go**

```
go get github.com/howijd/koios-rest-go-client
```
```go
...
import (
  "github.com/howijd/koios-rest-go-client" // imports as package "koios"
)
...
```

2. **CLI Application to interact with [Koios API] from Command-line see [./cmd/koios-rest](./cmd/koios-rest)**

```sh
# provides command `koios-rest` installed into ~/go/bin/koios-rest
go install github.com/howijd/koios-rest-go-client/cmd/koios-rest@latest
```

[![PkgGoDev](https://pkg.go.dev/badge/github.com/howijd/koios-rest-go-client)](https://pkg.go.dev/github.com/howijd/koios-rest-go-client)
![license](https://img.shields.io/github/license/howijd/koios-rest-go-client)

![tests](https://github.com/howijd/koios-rest-go-client/workflows/tests/badge.svg)
[![Coverage Status](https://coveralls.io/repos/github/howijd/koios-rest-go-client/badge.svg?branch=main)](https://coveralls.io/github/howijd/koios-rest-go-client?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/howijd/koios-rest-go-client)](https://goreportcard.com/report/github.com/howijd/koios-rest-go-client)
![GitHub last commit](https://img.shields.io/github/last-commit/howijd/koios-rest-go-client)

--- 

- [Install](#install)
- [Usage](#usage)
  - [Basic usage](#basic-usage)
  - [Concurrency using goroutines](#concurrency-using-goroutines)
- [Lovelace (math on ada, assets and tokens).](#lovelace-math-on-ada-assets-and-tokens)
- [Implemented Endpoints](#implemented-endpoints)


---

## Install

Using this package requires a working Go environment. [See the install instructions for Go](http://golang.org/doc/install.html).

```
go get -u github.com/howijd/koios-rest-go-client
```

## Usage

Godoc includes many examples how to use the library see: [![PkgGoDev](https://pkg.go.dev/badge/github.com/howijd/koios-rest-go-client)](https://pkg.go.dev/github.com/howijd/koios-rest-go-client)
Additionally you can find all usecases by looking source of `koio-rest` Command-line application [source](./cmd/koios-rest) which utilizes entire API of this library.

### Basic usage

```go
package main

import (
	"context"
	"fmt"
	"log"

	koios "github.com/howijd/koios-rest-go-client"
)

func main() {
  // Call to koios.New without options is same as calling it with default opts.
  // See godoc for available configuration options.
  // api, err := koios.New(
  // 	koios.Host(koios.MainnetHost),
  // 	koios.APIVersion(koios.DefaultAPIVersion),
  // 	koios.Port(koios.DefaultPort),
  // 	koios.Schema(koios.DefaultSchema),
  // 	koios.HttpClient(koios.DefaultHttpClient),
  // ).
  api, err := koios.New()
  if err != nil {
    log.Fatal(err)
  }

  res, err := api.GetTip(context.Background())
  if err != nil {
	  log.Fatal(err)
  }
  fmt.Println("status: ", res.Status)
  fmt.Println("statu_code: ", res.StatusCode)

  fmt.Println("abs_slot: ", res.Tip.AbsSlot)
  fmt.Println("block_no: ", res.Tip.BlockNo)
  fmt.Println("block_time: ", res.Tip.BlockTime)
  fmt.Println("epoch: ", res.Tip.Epoch)
  fmt.Println("epoch_slot: ", res.Tip.EpochSlot)
  fmt.Println("hash: ", res.Tip.Hash)
}
```

### Concurrency using goroutines

This library is thread-safe so you can freerly use same api client instance passing it to your goroutines.

**Following example uses goroutines to query chain tip from different endpoints.**

```go
func main() {
  api, _ := koios.New(
    // limit client request 1 per second even though 
    // this example will send requests in goroutines.
    koios.RateLimit(1),
  )
  ctx := context.Background()
  defer cancel()
  var wg sync.WaitGroup
  servers := []string{
    "api.koios.rest",
    "guild.koios.rest",
    "testnet.koios.rest",
  }

  // Thanks to rate limit option requests will be made
  // once in a second.
  for _, host := range servers {
    wg.Add(1)
    go func(ctx context.Context, host string) {
      defer wg.Done()
      // switching host. all options changes are safe to call from goroutines.
      koios.Host(host)(api)
      res, _ := api.GET(ctx, "/tip")
      defer res.Body.Close()
      body, _ := ioutil.ReadAll(res.Body)
      fmt.Println(string(body))
    }(ctx, host)
  }
  
  wg.Wait()
  fmt.Println("requests done: ", api.TotalRequests())
}
```

## Lovelace (math on ada, assets and tokens).

Liprary uses for most cases to represent lovelace using [`Lovelace`](https://pkg.go.dev/github.com/howijd/koios-rest-go-client#Lovelace) data type.

This library uses forked snapshot of [github.com/shopspring/decimal] package to provide  
JSON and XML serialization/deserialization and make it ease to work with calculations  
and deciimal precisions of ADA lovelace and native assets.


**For decimal package API see**

[![](https://pkg.go.dev/badge/github.com/shopspring/decimal)](https://pkg.go.dev/github.com/shopspring/decimal) 

FORK: https://github.com/howijd/decimal  
issues and bug reports are welcome to: https://github.com/howijd/decimal/issues.

## Implemented Endpoints

> WORK IN PROGRESS

| **endpoint** | Method | CLI command | Godoc | API Doc |
| --- | --- | --- | --- | --- |
| NETWORK | | | | |
| `/tip` | `*.GetTip(...) *TipResponse` | `tip` | [![](https://pkg.go.dev/badge/github.com/howijd/koios-rest-go-client)](https://pkg.go.dev/github.com/howijd/koios-rest-go-client#Client.GetTip) | [![](https://img.shields.io/badge/API-doc-%2349cc90)](https://api.koios.rest/#get-/tip) |
| `/genesis` | `*.GetGenesis(...) *GenesisResponse` | `genesis` | [![PkgGoDev](https://pkg.go.dev/badge/github.com/howijd/koios-rest-go-client)](https://pkg.go.dev/github.com/howijd/koios-rest-go-client#Client.GetGenesis) | [![](https://img.shields.io/badge/API-doc-%2349cc90)](https://api.koios.rest/#get-/genesis) |
| `/totals` | `*.GetTotals(...) *TotalsResponse` | `totals` | [![PkgGoDev](https://pkg.go.dev/badge/github.com/howijd/koios-rest-go-client)](https://pkg.go.dev/github.com/howijd/koios-rest-go-client#Client.GetTotals) | [![](https://img.shields.io/badge/API-doc-%2349cc90)](https://api.koios.rest/#get-/totals) |
| EPOCH | | | | |
| `/epoch_info` | `*.GetEpochInfo(...) *EpochInfoResponse` | `epoch-info` | [![PkgGoDev](https://pkg.go.dev/badge/github.com/howijd/koios-rest-go-client)](https://pkg.go.dev/github.com/howijd/koios-rest-go-client#Client.GetEpochInfo) | [![](https://img.shields.io/badge/API-doc-%2349cc90)](https://api.koios.rest/#get-/epoch_info) |
| `/epoch_params` | `*.GetEpochParams(...) *EpochParamsResponse` | `epoch-params` | [![PkgGoDev](https://pkg.go.dev/badge/github.com/howijd/koios-rest-go-client)](https://pkg.go.dev/github.com/howijd/koios-rest-go-client#Client.GetEpochParams) | [![](https://img.shields.io/badge/API-doc-%2349cc90)](https://api.koios.rest/#get-/epoch_params) |
| BLOCK | | | | |
| `/blocks` | `*.GetBlocks(...) *BlocksResponse` | `blocks` | [![PkgGoDev](https://pkg.go.dev/badge/github.com/howijd/koios-rest-go-client)](https://pkg.go.dev/github.com/howijd/koios-rest-go-client#Client.GetBlocks) | [![](https://img.shields.io/badge/API-doc-%2349cc90)](https://api.koios.rest/#get-/blocks) |
| `/block_info` | `*.GetBlockInfo(...) *BlockInfoResponse` | `block-info` | [![PkgGoDev](https://pkg.go.dev/badge/github.com/howijd/koios-rest-go-client)](https://pkg.go.dev/github.com/howijd/koios-rest-go-client#Client.GetBlockInfo) | [![](https://img.shields.io/badge/API-doc-%2349cc90)](https://api.koios.rest/#get-/block_info) |
| `/block_txs` | `*.GetBlockTxs(...) *BlockTxsResponse` | `block-txs` | [![PkgGoDev](https://pkg.go.dev/badge/github.com/howijd/koios-rest-go-client)](https://pkg.go.dev/github.com/howijd/koios-rest-go-client#Client.GetBlockTxs) | [![](https://img.shields.io/badge/API-doc-%2349cc90)](https://api.koios.rest/#get-/block_txs) |
| TRANSACTIONS | | | | |
| `/tx_info` | `*.GetTxsInfos(...) *TxsInfoResponse` | `txs-infos` | [![PkgGoDev](https://pkg.go.dev/badge/github.com/howijd/koios-rest-go-client)](https://pkg.go.dev/github.com/howijd/koios-rest-go-client#Client.GetTxsInfos) | [![](https://img.shields.io/badge/API-doc-%2349cc90)](https://api.koios.rest/#get-/tx_info) |
| | `*.GetTxInfo(...) *TxInfoResponse` | `tx-info` | [![PkgGoDev](https://pkg.go.dev/badge/github.com/howijd/koios-rest-go-client)](https://pkg.go.dev/github.com/howijd/koios-rest-go-client#Client.GetTxInfo) | |
| `/tx_utxos` | `*.GetTxsUTxOs(...) *TxUTxOsResponse` | `tx-utxos` | [![PkgGoDev](https://pkg.go.dev/badge/github.com/howijd/koios-rest-go-client)](https://pkg.go.dev/github.com/howijd/koios-rest-go-client#Client.GetTxsUTxOs) | [![](https://img.shields.io/badge/API-doc-%2349cc90)](https://api.koios.rest/#get-/tx_utxos) |
| `/tx_metadata` | `*.GetTxsMetadata(...) *TxsMetadataResponse` | `txs-metadata` | [![PkgGoDev](https://pkg.go.dev/badge/github.com/howijd/koios-rest-go-client)](https://pkg.go.dev/github.com/howijd/koios-rest-go-client#Client.GetTxsMetadata) | [![](https://img.shields.io/badge/API-doc-%2349cc90)](https://api.koios.rest/#get-/tx_metadata) |
| | `*.GetTxMetadata(...) *TxMetadataResponse` | `tx-metadata` | [![PkgGoDev](https://pkg.go.dev/badge/github.com/howijd/koios-rest-go-client)](https://pkg.go.dev/github.com/howijd/koios-rest-go-client#Client.GetTxMetadata) | [![](https://img.shields.io/badge/API-doc-%2349cc90)](https://api.koios.rest/#get-/tx_metadata) |
| ADDRESS | | | | |
| ACCOUNT | | | | |
| POOL | | | | |
| SCRIPT | | | | |

<!-- 
[![PkgGoDev](https://pkg.go.dev/badge/github.com/howijd/koios-rest-go-client)]() | [![](https://img.shields.io/badge/API-doc-%2349cc90)]()
-->
<!-- LINKS -->
[Koios API]: https://koios.rest "Koios API"
[github.com/shopspring/decimal]: https://github.com/shopspring/decimal
