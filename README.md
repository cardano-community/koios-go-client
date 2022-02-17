<h1>Koios API Client Library for Go</h1>

**:warning: UNTIL v1.0.0 RELEASE THIS LIBRARY WILL GET BREAKING CHANGES :warning:**

- before v1 every `v0.x.x` MINOR semver update most likely has breaking change. 
- before updating e.g. `go get -u` check for changes to prevent inconveniences. 
- `v1.0.0` enhancements are tracked under following [issue](https://github.com/howijd/koios-rest-go-client/issues/1)

**[Koios API] is Elastic Cardano Query Layer!**

> A consistent query layer for developers to build upon Cardano, with   
> multiple, redundant endpoints that allow for easy scalability.

**[Koios API] Client Library for Go**

[![PkgGoDev](https://pkg.go.dev/badge/github.com/howijd/koios-rest-go-client)](https://pkg.go.dev/github.com/howijd/koios-rest-go-client)

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

**CLI Application to interact with [Koios API] from Command-line see**

[See installing instruction](#cli-application)

**Build Status**

[![linux](https://github.com/howijd/koios-rest-go-client/workflows/linux/badge.svg)](https://github.com/howijd/koios-rest-go-client/actions/workflows/linux.yaml)
[![macos](https://github.com/howijd/koios-rest-go-client/workflows/macos/badge.svg)](https://github.com/howijd/koios-rest-go-client/actions/workflows/macos.yaml)
[![windows](https://github.com/howijd/koios-rest-go-client/workflows/windows/badge.svg)](https://github.com/howijd/koios-rest-go-client/actions/workflows/windows.yaml)

**Development Status**

![GitHub last commit](https://img.shields.io/github/last-commit/howijd/koios-rest-go-client)
[![coverage](https://coveralls.io/repos/github/howijd/koios-rest-go-client/badge.svg?branch=main)](https://coveralls.io/github/howijd/koios-rest-go-client?branch=main)
[![codeql](https://github.com/howijd/koios-rest-go-client/workflows/codeql/badge.svg)](https://github.com/howijd/koios-rest-go-client/actions/workflows/codeql.yaml)
[![misspell](https://github.com/howijd/koios-rest-go-client/workflows/misspell/badge.svg)](https://github.com/howijd/koios-rest-go-client/actions/workflows/misspell.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/howijd/koios-rest-go-client)](https://goreportcard.com/report/github.com/howijd/koios-rest-go-client)

---

- [Usage](#usage)
  - [Basic usage](#basic-usage)
  - [Concurrency using goroutines](#concurrency-using-goroutines)
- [Lovelace (math on ada, assets and tokens).](#lovelace-math-on-ada-assets-and-tokens)
- [Implemented Endpoints](#implemented-endpoints)
- [CLI Application](#cli-application)
  - [List of all commands](#list-of-all-commands)
  - [Example Usage](#example-usage)
    - [Example to query mainnet tip from cli](#example-to-query-mainnet-tip-from-cli)
    - [Example to query testnet tip from cli](#example-to-query-testnet-tip-from-cli)
  - [Install](#install)
    - [Install from Source](#install-from-source)
- [Contributing](#contributing)
  - [Code of Conduct](#code-of-conduct)
  - [Got a Question or Problem?](#got-a-question-or-problem)
  - [Issues and Bugs](#issues-and-bugs)
  - [Feature Requests](#feature-requests)
  - [Submission Guidelines](#submission-guidelines)
    - [Submitting an Issue](#submitting-an-issue)
    - [Submitting a Pull Request (PR)](#submitting-a-pull-request-pr)
    - [After your pull request is merged](#after-your-pull-request-is-merged)
  - [Coding Rules](#coding-rules)
  - [Commit Message Guidelines](#commit-message-guidelines)
    - [Commit Message Format](#commit-message-format)
    - [Revert](#revert)
    - [Type](#type)
    - [Scope](#scope)
    - [Subject](#subject)
    - [Body](#body)
    - [Footer](#footer)
  - [Development Documentation](#development-documentation)
    - [Setup your machine](#setup-your-machine)

---

## Usage

See Godoc [![PkgGoDev](https://pkg.go.dev/badge/github.com/howijd/koios-rest-go-client)](https://pkg.go.dev/github.com/howijd/koios-rest-go-client)
Additionally you can find all usecases by looking source of `koio-rest` Command-line application [source](./cmd/koios-rest) which utilizes entire API of this library.

**NOTE**

Library normalizes some of the API responses and constructs Typed response for each end point.
If you wish to work with `*http.Response` directly you can do so by using api client `GET,POST, HEAD` methods.

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

  fmt.Println("abs_slot: ", res.Data.AbsSlot)
  fmt.Println("block_no: ", res.Data.BlockNo)
  fmt.Println("block_time: ", res.Data.BlockTime)
  fmt.Println("epoch: ", res.Data.Epoch)
  fmt.Println("epoch_slot: ", res.Data.EpochSlot)
  fmt.Println("hash: ", res.Data.Hash)
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
      res, _ := api.GET(ctx, "/tip", nil, nil)
      defer res.Body.Close()
      body, _ := ioutil.ReadAll(res.Body)
      fmt.Println("Host: ", host)
      fmt.Println("Response: ", string(body))
    }(ctx, host)
  }

  wg.Wait()
  fmt.Println("requests done: ", api.TotalRequests())
}
```

## Lovelace (math on ada, assets and tokens).

Library uses for most cases to represent lovelace using [`Lovelace`](https://pkg.go.dev/github.com/howijd/koios-rest-go-client#Lovelace) data type.

This library uses forked snapshot of [github.com/shopspring/decimal] package to provide  
JSON and XML serialization/deserialization and make it ease to work with calculations  
and deciimal precisions of ADA lovelace and native assets.

**For decimal package API see**

[![](https://pkg.go.dev/badge/github.com/shopspring/decimal)](https://pkg.go.dev/github.com/shopspring/decimal)

FORK: https://github.com/howijd/decimal  
issues and bug reports are welcome to: https://github.com/howijd/decimal/issue.
So that we can ensure that [github.com/shopspring/decimal] repository is not polluted with 
issues which are not problems with upstream library.

## Implemented Endpoints

| **endpoint** | Go API | CLI command | API Doc |
| --- | --- | --- | --- | 
| **NETWORK** | | | |
| `/tip` | [`*.GetTip(...) *TipResponse`](https://pkg.go.dev/github.com/howijd/koios-rest-go-client#Client.GetTip) | `tip` | [![](https://img.shields.io/badge/API-doc-%2349cc90)](https://api.koios.rest/#get-/tip) |
| `/genesis` | [`*.GetGenesis(...) *GenesisResponse`](https://pkg.go.dev/github.com/howijd/koios-rest-go-client#Client.GetGenesis) | `genesis` | [![](https://img.shields.io/badge/API-doc-%2349cc90)](https://api.koios.rest/#get-/genesis) |
| `/totals` | [`*.GetTotals(...) *TotalsResponse`](https://pkg.go.dev/github.com/howijd/koios-rest-go-client#Client.GetTotals) | `totals` | [![](https://img.shields.io/badge/API-doc-%2349cc90)](https://api.koios.rest/#get-/totals) |
| **EPOCH** | | | |
| `/epoch_info` | [`*.GetEpochInfo(...) *EpochInfoResponse`](https://pkg.go.dev/github.com/howijd/koios-rest-go-client#Client.GetEpochInfo) | `epoch-info` | [![](https://img.shields.io/badge/API-doc-%2349cc90)](https://api.koios.rest/#get-/epoch_info) |
| `/epoch_params` | [`*.GetEpochParams(...) *EpochParamsResponse`](https://pkg.go.dev/github.com/howijd/koios-rest-go-client#Client.GetEpochParams) | `epoch-params` | [![](https://img.shields.io/badge/API-doc-%2349cc90)](https://api.koios.rest/#get-/epoch_params) |
| **BLOCK** | | | |
| `/blocks` | [`*.GetBlocks(...) *BlocksResponse`](https://pkg.go.dev/github.com/howijd/koios-rest-go-client#Client.GetBlocks) | `blocks` | [![](https://img.shields.io/badge/API-doc-%2349cc90)](https://api.koios.rest/#get-/blocks) |
| `/block_info` | [`*.GetBlockInfo(...) *BlockInfoResponse`](https://pkg.go.dev/github.com/howijd/koios-rest-go-client#Client.GetBlockInfo) | `block-info` | [![](https://img.shields.io/badge/API-doc-%2349cc90)](https://api.koios.rest/#get-/block_info) |
| `/block_txs` | [`*.GetBlockTxs(...) *BlockTxsResponse`](https://pkg.go.dev/github.com/howijd/koios-rest-go-client#Client.GetBlockTxs) | `block-txs` | [![](https://img.shields.io/badge/API-doc-%2349cc90)](https://api.koios.rest/#get-/block_txs) |
| **TRANSACTIONS** | | | |
| `/tx_info` | [`*.GetTxsInfos(...) *TxsInfoResponse`](https://pkg.go.dev/github.com/howijd/koios-rest-go-client#Client.GetTxsInfos) | `txs-infos` | [![](https://img.shields.io/badge/API-doc-%2349cc90)](https://api.koios.rest/#get-/tx_info) |
| | [`*.GetTxInfo(...) *TxInfoResponse`](https://pkg.go.dev/github.com/howijd/koios-rest-go-client#Client.GetTxInfo) | `tx-info` | |
| `/tx_utxos` | [`*.GetTxsUTxOs(...) *TxUTxOsResponse`](https://pkg.go.dev/github.com/howijd/koios-rest-go-client#Client.GetTxsUTxOs) | `tx-utxos` | [![](https://img.shields.io/badge/API-doc-%2349cc90)](https://api.koios.rest/#get-/tx_utxos) |
| `/tx_metadata` | [`*.GetTxsMetadata(...) *TxsMetadataResponse`](https://pkg.go.dev/github.com/howijd/koios-rest-go-client#Client.GetTxsMetadata) | `txs-metadata` | [![](https://img.shields.io/badge/API-doc-%2349cc90)](https://api.koios.rest/#get-/tx_metadata) |
| | [`*.GetTxMetadata(...) *TxMetadataResponse`](https://pkg.go.dev/github.com/howijd/koios-rest-go-client#Client.GetTxMetadata) | `tx-metadata` | [![](https://img.shields.io/badge/API-doc-%2349cc90)](https://api.koios.rest/#get-/tx_metadata) |
| `/tx_metalabels` | [`*.GetTxMetaLabels(...) *TxMetaLabelsResponse`](https://pkg.go.dev/github.com/howijd/koios-rest-go-client#Client.GetTxMetaLabels) | `tx-metalabels` | [![](https://img.shields.io/badge/API-doc-%2349cc90)](https://api.koios.rest/#get-/tx_metalabels) |
| `/submittx` | [`*.SubmitSignedTx(...) *SubmitSignedTxResponse`](https://pkg.go.dev/github.com/howijd/koios-rest-go-client#Client.SubmitSignedTx) | `tx-submit` | [![](https://img.shields.io/badge/API-doc-%2349cc90)](https://api.koios.rest/#get-/submittx) |
| `/tx_status` | [`*.GetTxsStatuses(...) *TxsStatusesResponse`](https://pkg.go.dev/github.com/howijd/koios-rest-go-client#Client.GetTxsStatuses) | `txs-statuses` | [![](https://img.shields.io/badge/API-doc-%2349cc90)](https://api.koios.rest/#get-/tx_status) |
|  | [`*.GetTxStatus(...) *TxStatusResponse`](https://pkg.go.dev/github.com/howijd/koios-rest-go-client#Client.GetTxStatus) | `tx-status` | [![](https://img.shields.io/badge/API-doc-%2349cc90)](https://api.koios.rest/#get-/tx_status) |
| **ADDRESS** | | | |
| `/address_info` | [`*.GetAddressInfo(...) *AddressInfoResponse`](https://pkg.go.dev/github.com/howijd/koios-rest-go-client#Client.GetAddressInfo) | `address-info` | [![](https://img.shields.io/badge/API-doc-%2349cc90)](https://api.koios.rest/#get-/address_info) |
| `/address_txs` | [`*.GetAddressTxs(...) *AddressTxsResponse`](https://pkg.go.dev/github.com/howijd/koios-rest-go-client#Client.GetAddressTxs) | `address-txs` | [![](https://img.shields.io/badge/API-doc-%2349cc90)](https://api.koios.rest/#get-/address_txs) |
| `/address_assets` | [`*.GetAddressAssets(...) *AddressAssetsResponse`](https://pkg.go.dev/github.com/howijd/koios-rest-go-client#Client.GetAddressAssets) | `address-assets` | [![](https://img.shields.io/badge/API-doc-%2349cc90)](https://api.koios.rest/#get-/address_assets) |
| `/credential_txs` | [`*.GetCredentialTxs(...) *CredentialTxsResponse`](https://pkg.go.dev/github.com/howijd/koios-rest-go-client#Client.GetCredentialTxs) | `credential-txs` | [![](https://img.shields.io/badge/API-doc-%2349cc90)](https://api.koios.rest/#get-/credential_txs) |
| **ACCOUNT** | | | |
| `/account_list` | [`*.GetAccountList(...) *AccountListResponse`](https://pkg.go.dev/github.com/howijd/koios-rest-go-client#Client.GetAccountList) | `account-list` | [![](https://img.shields.io/badge/API-doc-%2349cc90)](https://api.koios.rest/#get-/account_list) |
| `/account_info` | [`*.GetAccountInfo(...) *AccountListResponse`](https://pkg.go.dev/github.com/howijd/koios-rest-go-client#Client.GetAccountInfo) | `account-info` | [![](https://img.shields.io/badge/API-doc-%2349cc90)](https://api.koios.rest/#get-/account_info) |
| `/account_rewards` | [`*.GetAccountRewards(...) *AccountRewardsResponse`](https://pkg.go.dev/github.com/howijd/koios-rest-go-client#Client.GetAccountRewards) | `account-rewards` | [![](https://img.shields.io/badge/API-doc-%2349cc90)](https://api.koios.rest/#get-/account_info) |
| `/account_updates` | [`*.GetAccountUpdates(...) *AccountUpdatesResponse`](https://pkg.go.dev/github.com/howijd/koios-rest-go-client#Client.GetAccountUpdates) | `account-updates` | [![](https://img.shields.io/badge/API-doc-%2349cc90)](https://api.koios.rest/#get-/account_updates) |
| `/account_addresses` | [`*.GetAccountAddresses(...) *AccountAddressesResponse`](https://pkg.go.dev/github.com/howijd/koios-rest-go-client#Client.GetAccountAddresses) | `account-addresses` | [![](https://img.shields.io/badge/API-doc-%2349cc90)](https://api.koios.rest/#get-/account_addresses) |
| `/account_assets` | [`*.GetAccountAssets(...) *AccountAssetsResponse`](https://pkg.go.dev/github.com/howijd/koios-rest-go-client#Client.GetAccountAssets) | `account-assets` | [![](https://img.shields.io/badge/API-doc-%2349cc90)](https://api.koios.rest/#get-/account_assets) |
| `/account_history` | [`*.GetAccountHistory(...) *AccountHistoryResponse`](https://pkg.go.dev/github.com/howijd/koios-rest-go-client#Client.GetAccountHistory) | `account-history` | [![](https://img.shields.io/badge/API-doc-%2349cc90)](https://api.koios.rest/#get-/account_history) |
| **ASSET** | | | |
| `/asset_list` | [`*.GetAssetList(...) *AssetAddressListResponse`](https://pkg.go.dev/github.com/howijd/koios-rest-go-client#Client.GetAssetList) | `asset-list` | [![](https://img.shields.io/badge/API-doc-%2349cc90)](https://api.koios.rest/#get-/asset_address_list) |
| `/asset_address_list` | [`*.GetAssetAddressList(...) *AssetAddressListResponse`](https://pkg.go.dev/github.com/howijd/koios-rest-go-client#Client.GetAssetAddressList) | `asset-address-list` | [![](https://img.shields.io/badge/API-doc-%2349cc90)](https://api.koios.rest/#get-/asset_address_list) |
| `/asset_info` | [`*.GetAssetInfo(...) *AssetInfoResponse`](https://pkg.go.dev/github.com/howijd/koios-rest-go-client#Client.GetAssetInfo) | `asset-info` | [![](https://img.shields.io/badge/API-doc-%2349cc90)](https://api.koios.rest/#get-/asset_info) |
| `/asset_summary` | [`*.GetAssetSummary(...) *AssetSummaryResponse`](https://pkg.go.dev/github.com/howijd/koios-rest-go-client#Client.GetAssetSummary) | `asset-summary` | [![](https://img.shields.io/badge/API-doc-%2349cc90)](https://api.koios.rest/#get-/asset_summary) |
| `/asset_txs` | [`*.GetAssetTxs(...) *AssetTxsResponse`](https://pkg.go.dev/github.com/howijd/koios-rest-go-client#Client.GetAssetTxs) | `asset-txs` | [![](https://img.shields.io/badge/API-doc-%2349cc90)](https://api.koios.rest/#get-/asset_txs) |
| **POOL** | | | |
| `/pool_list` | [`*.GetPoolList(...) *PoolListResponse`](https://pkg.go.dev/github.com/howijd/koios-rest-go-client#Client.GetPoolList) | `pool-list` | [![](https://img.shields.io/badge/API-doc-%2349cc90)](https://api.koios.rest/#get-/pool_list) |
| `/pool_info` | [`*.GetPoolInfos(...) *PoolInfoResponse`](https://pkg.go.dev/github.com/howijd/koios-rest-go-client#Client.GetPoolInfos) | `pool-infos` | [![](https://img.shields.io/badge/API-doc-%2349cc90)](https://api.koios.rest/#get-/pool_info) |
| | [`*.GetPoolInfo(...) *PoolInfoResponse`](https://pkg.go.dev/github.com/howijd/koios-rest-go-client#Client.GetPoolInfo) | `pool-info` | |
| `/pool_delegators` | [`*.GetPoolDelegators(...) *PoolDelegatorsResponse`](https://pkg.go.dev/github.com/howijd/koios-rest-go-client#Client.GetPoolDelegators) | `pool-delegators` | [![](https://img.shields.io/badge/API-doc-%2349cc90)](https://api.koios.rest/#get-/pool_delegators) |
| `/pool_blocks` | [`*.GetPoolBlocks(...) *PoolBlocksResponse`](https://pkg.go.dev/github.com/howijd/koios-rest-go-client#Client.GetPoolBlocks) | `pool-blocks` | [![](https://img.shields.io/badge/API-doc-%2349cc90)](https://api.koios.rest/#get-/pool_blocks) |
| `/pool_updates` | [`*.GetPoolUpdates(...) *PoolUpdatesResponse`](https://pkg.go.dev/github.com/howijd/koios-rest-go-client#Client.GetPoolUpdates) | `pool-updates` | [![](https://img.shields.io/badge/API-doc-%2349cc90)](https://api.koios.rest/#get-/pool_updates) |
| `/pool_relays` | [`*.GetPoolRelays(...) *PoolRelaysResponse`](https://pkg.go.dev/github.com/howijd/koios-rest-go-client#Client.GetPoolRelays) | `pool-relays` | [![](https://img.shields.io/badge/API-doc-%2349cc90)](https://api.koios.rest/#get-/pool_relays) |
| `/pool_metadata` | [`*.GetPoolMetadata(...) *PoolMetadataResponse`](https://pkg.go.dev/github.com/howijd/koios-rest-go-client#Client.GetPoolMetadata) | `pool-metadata` | [![](https://img.shields.io/badge/API-doc-%2349cc90)](https://api.koios.rest/#get-/pool_metadata) |
| **SCRIPT** | | | |
| `/script_list` | [`*.GetScriptList(...) *ScriptRedeemersResponse`](https://pkg.go.dev/github.com/howijd/koios-rest-go-client#Client.GetScriptList) | `script-list` | [![](https://img.shields.io/badge/API-doc-%2349cc90)](https://api.koios.rest/#get-/script_list) |
| `/script_redeemers` | [`*.GetScriptRedeemers(...) *ScriptListResponse`](https://pkg.go.dev/github.com/howijd/koios-rest-go-client#Client.GetScriptRedeemers) | `script-redeemers` | [![](https://img.shields.io/badge/API-doc-%2349cc90)](https://api.koios.rest/#get-/script_redeemers) |

## CLI Application

source of cli: [./cmd/koios-rest](./cmd/koios-rest).

[Installation instructions](#install)

### List of all commands

<details>
  <summary><code>koios-rest --help</code></summary>

```
NAME:
   koios-rest - CLI Client to consume Koios API https://api.koios.rest

USAGE:
   koios-rest [global options] command [command options] [arguments...]

VERSION:
   (devel)

AUTHOR:
   The Cardano Community Authors

COMMANDS:
   help, h  Shows a list of commands or help for one command
   ACCOUNT:
     account-list       Get a list of all accounts returns array of stake addresses.
     account-info       Get the account info of any (payment or staking) address.
     account-rewards    Get the full rewards history (including MIR) for a stake address, or certain epoch if specified.
     account-updates    Get the account updates (registration, deregistration, delegation and withdrawals).
     account-addresses  Get all addresses associated with an account payment or staking address
     account-assets     Get the native asset balance of an account.
     account-history    Get the staking history of an account.
   ADDRESS:
     address-info    Get address info - balance, associated stake address (if any) and UTxO set.
     address-txs     Get the transaction hash list of input address array, optionally filtering after specified block height (inclusive).
     address-assets  Get the list of all the assets (policy, name and quantity) for a given address.
     credential-txs  Get the transaction hash list of input payment credential array, optionally filtering after specified block height (inclusive).
   ASSET:
     asset-list          Get the list of all native assets (paginated).
     asset-address-list  Get the list of all addresses holding a given asset.
     asset-info          Get the information of an asset including first minting & token registry metadata.
     asset-summary       Get the summary of an asset (total transactions exclude minting/total wallets include only wallets with asset balance).
     asset-txs           Get the list of all asset transaction hashes (newest first).
   BLOCK:
     blocks      Get summarised details about all blocks (paginated - latest first).
     block-info  Get detailed information about a specific block.
     block-txs   Get a list of all transactions included in a provided block.
   EPOCH:
     epoch-info    Get the epoch information, all epochs if no epoch specified.
     epoch-params  Get the protocol parameters for specific epoch, returns information about all epochs if no epoch specified.
   NETWORK:
     tip      Get the tip info about the latest block seen by chain.
     genesis  Get the Genesis parameters used to start specific era on chain.
     totals   Get the circulating utxo, treasury, rewards, supply and reserves in lovelace for specified epoch, all epochs if empty.
   POOL:
     pool-list        A list of all currently registered/retiring (not retired) pools.
     pool-infos       Current pool statuses and details for a specified list of pool ids.
     pool-info        Current pool status and details for a specified pool by pool id.
     pool-delegators  Return information about delegators by a given pool and optional epoch (current if omitted).
     pool-blocks      Return information about blocks minted by a given pool in current epoch (or _epoch_no if provided).
     pool-updates     Return all pool updates for all pools or only updates for specific pool if specified.
     pool-relays      A list of registered relays for all currently registered/retiring (not retired) pools.
     pool-metadata    Metadata(on & off-chain) for all currently registered/retiring (not retired) pools.
   SCRIPT:
     script-list       List of all existing script hashes along with their creation transaction hashes.
     script-redeemers  List of all redeemers for a given script hash.
   TRANSACTIONS:
     txs-infos      Get detailed information about transaction(s).
     tx-info        Get detailed information about single transaction.
     tx-utxos       Get UTxO set (inputs/outputs) of transactions.
     txs-metadata   Get metadata information (if any) for given transaction(s).
     tx-metadata    Get metadata information (if any) for given transaction.
     tx-metalabels  Get a list of all transaction metalabels.
     tx-submit      Submit signed transaction to the network.
     txs-statuses   Get the number of block confirmations for a given transaction hash list
     tx-status      Get the number of block confirmations for a given transaction hash
   UTILS:
     get   get issues a GET request to the specified API endpoint
     head  head issues a HEAD request to the specified API endpoint

GLOBAL OPTIONS:
   --port value, -p value  Set port (default: 443)
   --host value            Set host (default: "api.koios.rest")
   --api-version value     Set API version (default: "v0")
   --schema value          Set URL schema (default: "https")
   --origin value          Set Origin header for requests. (default: "https://github.com/howijd/koios-rest-go-client")
   --rate-limit value      Set API Client rate limit for outgoing requests (default: 5)
   --no-format             prints response json strings directly without calling json pretty. (default: false)
   --enable-req-stats      Enable request stats. (default: false)
   --testnet               use default testnet as host. (default: false)
   --help, -h              show help (default: false)
   --version, -v           print the version (default: false)

COPYRIGHT:
   (c) 2022
```

</details>

---

### Example Usage

#### Example to query mainnet tip from cli

```cli
koios-rest --enable-req-stats tip
```

response

```json
{
  "request_url": "https://api.koios.rest/api/v0/tip",
  "request_method": "GET",
  "status_code": 200,
  "status": "200 OK",
  "date": "Mon, 07 Feb 2022 12:49:49 GMT",
  "content_range": "0-0/*",
  "stats": {
    "req_started_at": "2022-02-07T12:49:48.565834833Z",
    "req_dns_lookup_dur": 1284269, // dns lookup duration in nanosecons.
    "tls_hs_dur": 208809082, // handshake duration in nanosecons.
    "est_cxn_dur": 159857626, // time it took to establish connection with server in nanosecons.
    "ttfb": 998874037, // time since start of the request it took to recieve first byte.
    "req_dur": 999186595, // total request duration in nanoseconds
    "req_dur_str": "999.186595ms" // string of req_dur
  },
  "data": {
    "abs_slot": 52671876,
    "block_no": 6852764,
    "block_time": "2022-02-07T12:49:27",
    "epoch": 319,
    "epoch_slot": 227076,
    "hash": "1dad134750188460dd48068e655b5935403d2f51afaf53a39337a4c89771754a"
  }

```

---

#### Example to query testnet tip from cli

```cli
koios-rest --enable-req-stats --testnet tip
# OR
koios-rest --enable-req-stats --host testnet.koios.rest tip
```

response

```json
{
  "request_url": "https://testnet.koios.rest/api/v0/tip",
  "request_method": "GET",
  "status_code": 200,
  "status": "200 OK",
  "date": "Mon, 07 Feb 2022 12:50:04 GMT",
  "content_range": "0-0/*",
  "stats": {
    "req_started_at": "2022-02-07T12:50:03.98615637Z",
    "req_dns_lookup_dur": 1383437,
    "tls_hs_dur": 69093093,
    "est_cxn_dur": 43733700,
    "ttfb": 167423049,
    "req_dur": 167738287,
    "req_dur_str": "167.738287ms"
  },
  "data": {
    "abs_slot": 49868948,
    "block_no": 3300758,
    "block_time": "2022-02-07T12:49:24",
    "epoch": 185,
    "epoch_slot": 318548,
    "hash": "d7623e68cb78f450f42ba4b5a169124b26677f08f676ca4241b27edb6dbf0071"
  }
}
```

### Install

It's highly recommended installing a latest version of koios-rest available on the [releases page](https://github.com/howijd/koios-rest-go-client/releases/latest).

#### Install from Source

Installing from source requires a working Go environment. [See the install instructions for Go](http://golang.org/doc/install.html).

Since `koios-rest` cli application uses `replace` in [./cmd/koiso-rest/go.mod](https://github.com/howijd/koios-rest-go-client/blob/main/cmd/koios-rest/go.mod).
Then `go install won't work`. To install it from source use following commands.

1. `git clone git@github.com:howijd/koios-rest-go-client.git`
2. `cd ./koios-rest-go-client/cmd/koios-rest`
3. `go install .`

verify installation

`koios-rest --version`

---

## Contributing

We would love for you to contribute to [Koios API Client Library for Go][github] and help make it even better than it is today! As a contributor, here are the guidelines we would like you to follow:

 - [Code of Conduct](#coc)
 - [Question or Problem?](#got-a-question-or-problem)
 - [Found a Bug?](#issues-and-bugs)
 - [Missing a Feature?](#feature-requests)
 - [Submission Guidelines](#submission-guidelines)
 - [Coding Rules](#coding-rules)
 - [Commit Message Guidelines](#commit-message-guidelines)
 - [Development Documentation](#development-documentation)

### Code of Conduct

Help us keep [Koios API Client Library for Go][github] open and inclusive. Please read and follow our [Code of Conduct][code-of-conduc]

---

### Got a Question or Problem?

Do not open issues for general support questions as we want to keep GitHub issues for bug reports and feature requests. You've got much better chances of getting your question answered on [Koios Telegram Group](https://t.me/joinchat/+zE4Lce_QUepiY2U1)

---

### Issues and Bugs

If you find a bug in the source code, you can help us by
[submitting an issue](#submit-issue) to our [GitHub Repository][github]. Even better, you can
[submit a Pull Request](#submit-pr) with a fix.

---

### Feature Requests
You can *request* a new feature by [submitting an issue](#submit-issue) to our GitHub
Repository. If you would like to *implement* a new feature, please submit an issue with
a proposal for your work first, to be sure that we can use it.
Please consider what kind of change it is:

* For a **Major Feature**, first open an issue and outline your proposal so that it can be
discussed. This will also allow us to better coordinate our efforts, prevent duplication of work,
and help you to craft the change so that it is successfully accepted into the project.
* **Small Features** can be crafted and directly [submitted as a Pull Request](#submit-pr).

---

### Submission Guidelines

#### Submitting an Issue

Before you submit an issue, please search the issue tracker, maybe an issue for your problem already exists and the discussion might inform you of workarounds readily available.

You can file new issues by filling out our [new issue form](https://github.com/howijd/koios-rest-go-client/issues/new).

---

#### Submitting a Pull Request (PR)

Before you submit your Pull Request (PR) consider the following guidelines:

1. Search [GitHub](https://github.com/howijd/koios-rest-go-client/pulls) for an open or closed PR that relates to your submission. You don't want to duplicate effort.
2. Fork the [howijd/koios-rest-go-client][github] repo.
3. Setup you local repository

    ```shell
    git@github.com:<your-github-username>/koios-rest-go-client.git
    cd koios-rest-go-client
    git remote add upstream git@github.com:howijd/koios-rest-go-client.git
    ```
4. Make your changes in a new git branch and ensure that you always start from up to date main branch. **Repeat this step every time you are about to start woking on new PR**.

    e.g. Start new change work to update readme:
    ```shell
    # if you are not in main branch e.g. still on previous work branch
    git checkout main
    git pull --ff upstream main
    git checkout -b update-readme main
    ```
5. Create your patch, **including appropriate test cases**.
6. Follow our [Coding Rules](#rules).
7. If changes are in source code except documentations then run the full test suite, as described in the [developer documentation](#dev-doc), and ensure that all tests pass.
8.  Commit your changes using a descriptive commit message that follows our
  [commit message conventions](#commit). Adherence to these conventions
  is necessary because release notes are automatically generated from these messages.

     ```shell
     git add -A
     git commit --signoff
     # or in short
     git commit -sm"docs(markdown): update readme examples"
     ```
9. Push your branch to GitHub:

    ```shell
    git push -u origin update-readme
    ```
10. In GitHub, send a pull request to `main` branch.
* If we suggest changes then:
  * Make the required updates.
  * Re-run the test suites to ensure tests are still passing.
  * Rebase your branch and force push to your GitHub repository (this will update your Pull Request):

     ```shell
    git fetch --all
    git rebase upstream main
    git push -uf origin update-readme
    ```
That's it! Thank you for your contribution!

---

#### After your pull request is merged

After your pull request is merged, you can safely delete your branch and pull the changes from the main (upstream) repository:

* Delete the remote branch on GitHub either through the GitHub web UI or your local shell as follows:
  
    ```shell
    git push origin --delete update-readme
    ```
* Check out the main branch:
  
    ```shell
    git checkout main -f
    ```

* Delete the local branch:

    ```shell
    git branch -D update-readme
    ```
* Update your master with the latest upstream version:

    ```shell
    git pull --ff upstream main
    ```
---

### Coding Rules

To ensure consistency throughout the source code, keep these rules in mind as you are working:

* All features or bug fixes **must be tested** by one or more specs (unit-tests).
* All public API methods **must be documented**.

---

### Commit Message Guidelines

[![Conventional Commits](https://img.shields.io/badge/Conventional%20Commits-1.0.0-yellow.svg)](https://conventionalcommits.org)

We have very precise rules over how our git commit messages can be formatted. This leads to **more readable messages** that are easy to follow when looking through the **project history**. Commit messages should be well formatted, and to make that "standardized", we are using Conventional Commits. Our release workflow uses these rules to generate changelogs.

---

#### Commit Message Format

Each commit message consists of a **header**, a **body** and a **footer**.  The header has a special format that includes a **type**, a **scope** and a **subject**:

```
<type>(<scope>): <subject>
<BLANK LINE>
<body>
<BLANK LINE>
<footer>
```

*When maintainers are merging PR merge commit should be edited:*

```
<type>(<scope>): <subject> (#pr)
<BLANK LINE>
<body>
<BLANK LINE>
<footer>
```

The **header** is mandatory and the **scope** of the header is optional.

Any line of the commit message cannot be longer 100 characters! This allows the message to be easier to read on GitHub as well as in various git tools.

The footer should contain a [closing reference to an issue](https://help.github.com/articles/closing-issues-via-commit-messages/) if any.

Samples:

```
docs(markdown): update readme examples
```

```
fix(cli): fix cli command get action

description of your change.
```

```
refactor(client): change Client GET function signature

change order of client GET method arguments.

BREAKING CHANGE: Clien.Get signature has changed
```


---

#### Revert

If the commit reverts a previous commit, it should begin with `revert: `, followed by the header of the reverted commit. In the body it should say: `This reverts commit <hash>.`, where the hash is the SHA of the commit being reverted.

---

#### Type

Must be one of the following:

* **build**: Changes that affect the build system or external dependencies (example scopes: goreleaser, taskfile)
* **chore**: Other changes that don't modify src or test files.
* **ci**: Changes to our CI configuration files and scripts.
* **dep**: Changes related to dependecies e.g. `go.mod`
* **docs**: Documentation only changes (example scopes: markdown, godoc)
* **feat**: A new feature
* **fix**: A bug fix
* **perf**: A code change that improves performance
* **refactor**: A code change that neither fixes a bug nor adds a feature
* **revert**: Reverts a previous commit
* **style**: Changes that do not affect the meaning of the code (white-space, formatting, missing semi-colons, etc)
* **test**: Adding missing tests or correcting existing tests

---

#### Scope

The following is the list of supported scopes:

| scope | description |
| --- | --- |
| **cli** | CLI app related changes |
| **client** | API client related changes |
| **endpoint** | Changes related to api endpoints |
| **godoc** | Go documentation |
| **markdown** | Markdown files |
| **packaging** | Used for changes that change the release packages |

---

#### Subject

The subject contains a succinct description of the change:

* use the imperative, present tense: "change" not "changed" nor "changes"
* don't capitalize the first letter
* no dot (.) at the end
  
#### Body
Just as in the **subject**, use the imperative, present tense: "change" not "changed" nor "changes".
The body should include the motivation for the change and contrast this with previous behavior.

#### Footer
The footer should contain any information about **Breaking Changes** and is also the place to
reference GitHub issues that this commit **Closes**.

**Breaking Changes** should start with the word `BREAKING CHANGE:` with a space or two newlines. The rest of the commit message is then used for this.

A detailed explanation can be found in this [document][commit-message-format].

---

### Development Documentation

#### Setup your machine

**Prerequisites:**

* Working Go environment. [See the install instructions for Go](http://golang.org/doc/install.html).
* [golangci-lint](https://golangci-lint.run/usage/install/#local-installation) - Go linters aggregator should be installed
* [taskfile](https://taskfile.dev/#/installation) - task runner / build tool should be installed
* [svu](https://github.com/caarlos0/svu#install) - Semantic Version Util tool should be installed
* Fork the [howijd/koios-rest-go-client][github] repo.
* Setup you local repository

    ```shell
    git@github.com:<your-github-username>/koios-rest-go-client.git
    cd koios-rest-go-client
    git remote add upstream git@github.com:howijd/koios-rest-go-client.git
    ```

**Setup local env**

```shell
task setup
```

**Lint your code**

```shell
task lint
```

**Test your change**

```shell
task test
```


**View code coverage report from in browser (results from `task test`)**

```shell
task cover
```

**Build snapshot binaries to ./cmd/koios-rest/dist.**

> use it if you want to test release packages

```shell
task build:snapshot
```

<!-- LINKS -->
[Koios API]: https://koios.rest "Koios API"
[github.com/shopspring/decimal]: https://github.com/shopspring/decimal
[coc]: https://github.com/howijd/.github/blob/main/CODE_OF_CONDUCT.md
[github]: https://github.com/howijd/koios-rest-go-client
