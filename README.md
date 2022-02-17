<h1>Koios API Client Library for Go</h1>

**:warning: UNTIL v1.0.0 RELEASE THIS LIBRARY WILL GET BREAKING CHANGES :warning:**

- before v1 every `v0.x.x` MINOR semver update most likely has breaking change. 
- before updating e.g. `go get -u` check for changes to prevent inconveniences. 

**[Koios API] is Elastic Cardano Query Layer!**

> A consistent query layer for developers to build upon Cardano, with   
> multiple, redundant endpoints that allow $$for easy scalability.

**[Koios API] Client Library for Go**

[![PkgGoDev](https://pkg.go.dev/badge/github.com/cardano-community/koios-go-client)](https://pkg.go.dev/github.com/cardano-community/koios-go-client)

```shell
go get github.com/cardano-community/koios-go-client
```

```go
...
import (
  "github.com/cardano-community/koios-go-client" // imports as package "koios"
)
...
```

**There is CLI Application to interact with [Koios API] from Command-line see:**

[koios cli repository for souce and installing instruction of koios cli][koios-cli]

**Build Status**

[![linux](https://github.com/cardano-community/koios-go-client/workflows/linux/badge.svg)](https://github.com/cardano-community/koios-go-client/actions/workflows/linux.yml)
[![macos](https://github.com/cardano-community/koios-go-client/workflows/macos/badge.svg)](https://github.com/cardano-community/koios-go-client/actions/workflows/macos.yml)
[![windows](https://github.com/cardano-community/koios-go-client/workflows/windows/badge.svg)](https://github.com/cardano-community/koios-go-client/actions/workflows/windows.yml)

**Development Status**

![GitHub last commit](https://img.shields.io/github/last-commit/cardano-community/koios-go-client)
[![codecov](https://codecov.io/gh/cardano-community/koios-go-client/branch/main/graph/badge.svg?token=FA1KGG6ZQ5)](https://codecov.io/gh/cardano-community/koios-go-client)
[![codeql](https://github.com/cardano-community/koios-go-client/workflows/codeql/badge.svg)](https://github.com/cardano-community/koios-go-client/actions/workflows/codeql.yml)
[![misspell](https://github.com/cardano-community/koios-go-client/workflows/misspell/badge.svg)](https://github.com/cardano-community/koios-go-client/actions/workflows/misspell.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/cardano-community/koios-go-client)](https://goreportcard.com/report/github.com/cardano-community/koios-go-client)

---

- [Usage](#usage)
  - [Basic usage](#basic-usage)
  - [Concurrency using goroutines](#concurrency-using-goroutines)
- [Lovelace (math on ada, assets and tokens).](#lovelace-math-on-ada-assets-and-tokens)
- [Implemented Endpoints](#implemented-endpoints)
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
- [| **markdown** | Markdown files |](#-markdown--markdown-files-)
    - [Subject](#subject)
    - [Body](#body)
    - [Footer](#footer)
  - [Development Documentation](#development-documentation)
    - [Setup your machine](#setup-your-machine)
- [Credits](#credits)

---

## Usage

See Godoc [![PkgGoDev](https://pkg.go.dev/badge/github.com/cardano-community/koios-go-client)](https://pkg.go.dev/github.com/cardano-community/koios-go-client)

Additionally you can find all usecases by looking source of `koio-cli` Command-line application [koios-cli] which utilizes entire API of this library.

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

	koios "github.com/cardano-community/koios-go-client"
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

Library uses for most cases to represent lovelace using [`Lovelace`](https://pkg.go.dev/github.com/cardano-community/koios-go-client#Lovelace) data type.

**For decimal package API see**

[![](https://pkg.go.dev/badge/github.com/shopspring/decimal)](https://pkg.go.dev/github.com/shopspring/decimal)


## Implemented Endpoints

| **endpoint** | Go API | CLI command | API Doc |
| --- | --- | --- | --- | 
| **NETWORK** | | | |
| `/tip` | [`*.GetTip(...) *TipResponse`](https://pkg.go.dev/github.com/cardano-community/koios-go-client#Client.GetTip) | `tip` | [![](https://img.shields.io/badge/API-doc-%2349cc90)](https://api.koios.rest/#get-/tip) |
| `/genesis` | [`*.GetGenesis(...) *GenesisResponse`](https://pkg.go.dev/github.com/cardano-community/koios-go-client#Client.GetGenesis) | `genesis` | [![](https://img.shields.io/badge/API-doc-%2349cc90)](https://api.koios.rest/#get-/genesis) |
| `/totals` | [`*.GetTotals(...) *TotalsResponse`](https://pkg.go.dev/github.com/cardano-community/koios-go-client#Client.GetTotals) | `totals` | [![](https://img.shields.io/badge/API-doc-%2349cc90)](https://api.koios.rest/#get-/totals) |
| **EPOCH** | | | |
| `/epoch_info` | [`*.GetEpochInfo(...) *EpochInfoResponse`](https://pkg.go.dev/github.com/cardano-community/koios-go-client#Client.GetEpochInfo) | `epoch-info` | [![](https://img.shields.io/badge/API-doc-%2349cc90)](https://api.koios.rest/#get-/epoch_info) |
| `/epoch_params` | [`*.GetEpochParams(...) *EpochParamsResponse`](https://pkg.go.dev/github.com/cardano-community/koios-go-client#Client.GetEpochParams) | `epoch-params` | [![](https://img.shields.io/badge/API-doc-%2349cc90)](https://api.koios.rest/#get-/epoch_params) |
| **BLOCK** | | | |
| `/blocks` | [`*.GetBlocks(...) *BlocksResponse`](https://pkg.go.dev/github.com/cardano-community/koios-go-client#Client.GetBlocks) | `blocks` | [![](https://img.shields.io/badge/API-doc-%2349cc90)](https://api.koios.rest/#get-/blocks) |
| `/block_info` | [`*.GetBlockInfo(...) *BlockInfoResponse`](https://pkg.go.dev/github.com/cardano-community/koios-go-client#Client.GetBlockInfo) | `block-info` | [![](https://img.shields.io/badge/API-doc-%2349cc90)](https://api.koios.rest/#get-/block_info) |
| `/block_txs` | [`*.GetBlockTxs(...) *BlockTxsResponse`](https://pkg.go.dev/github.com/cardano-community/koios-go-client#Client.GetBlockTxs) | `block-txs` | [![](https://img.shields.io/badge/API-doc-%2349cc90)](https://api.koios.rest/#get-/block_txs) |
| **TRANSACTIONS** | | | |
| `/tx_info` | [`*.GetTxsInfos(...) *TxsInfoResponse`](https://pkg.go.dev/github.com/cardano-community/koios-go-client#Client.GetTxsInfos) | `txs-infos` | [![](https://img.shields.io/badge/API-doc-%2349cc90)](https://api.koios.rest/#get-/tx_info) |
| | [`*.GetTxInfo(...) *TxInfoResponse`](https://pkg.go.dev/github.com/cardano-community/koios-go-client#Client.GetTxInfo) | `tx-info` | |
| `/tx_utxos` | [`*.GetTxsUTxOs(...) *TxUTxOsResponse`](https://pkg.go.dev/github.com/cardano-community/koios-go-client#Client.GetTxsUTxOs) | `tx-utxos` | [![](https://img.shields.io/badge/API-doc-%2349cc90)](https://api.koios.rest/#get-/tx_utxos) |
| `/tx_metadata` | [`*.GetTxsMetadata(...) *TxsMetadataResponse`](https://pkg.go.dev/github.com/cardano-community/koios-go-client#Client.GetTxsMetadata) | `txs-metadata` | [![](https://img.shields.io/badge/API-doc-%2349cc90)](https://api.koios.rest/#get-/tx_metadata) |
| | [`*.GetTxMetadata(...) *TxMetadataResponse`](https://pkg.go.dev/github.com/cardano-community/koios-go-client#Client.GetTxMetadata) | `tx-metadata` | [![](https://img.shields.io/badge/API-doc-%2349cc90)](https://api.koios.rest/#get-/tx_metadata) |
| `/tx_metalabels` | [`*.GetTxMetaLabels(...) *TxMetaLabelsResponse`](https://pkg.go.dev/github.com/cardano-community/koios-go-client#Client.GetTxMetaLabels) | `tx-metalabels` | [![](https://img.shields.io/badge/API-doc-%2349cc90)](https://api.koios.rest/#get-/tx_metalabels) |
| `/submittx` | [`*.SubmitSignedTx(...) *SubmitSignedTxResponse`](https://pkg.go.dev/github.com/cardano-community/koios-go-client#Client.SubmitSignedTx) | `tx-submit` | [![](https://img.shields.io/badge/API-doc-%2349cc90)](https://api.koios.rest/#get-/submittx) |
| `/tx_status` | [`*.GetTxsStatuses(...) *TxsStatusesResponse`](https://pkg.go.dev/github.com/cardano-community/koios-go-client#Client.GetTxsStatuses) | `txs-statuses` | [![](https://img.shields.io/badge/API-doc-%2349cc90)](https://api.koios.rest/#get-/tx_status) |
|  | [`*.GetTxStatus(...) *TxStatusResponse`](https://pkg.go.dev/github.com/cardano-community/koios-go-client#Client.GetTxStatus) | `tx-status` | [![](https://img.shields.io/badge/API-doc-%2349cc90)](https://api.koios.rest/#get-/tx_status) |
| **ADDRESS** | | | |
| `/address_info` | [`*.GetAddressInfo(...) *AddressInfoResponse`](https://pkg.go.dev/github.com/cardano-community/koios-go-client#Client.GetAddressInfo) | `address-info` | [![](https://img.shields.io/badge/API-doc-%2349cc90)](https://api.koios.rest/#get-/address_info) |
| `/address_txs` | [`*.GetAddressTxs(...) *AddressTxsResponse`](https://pkg.go.dev/github.com/cardano-community/koios-go-client#Client.GetAddressTxs) | `address-txs` | [![](https://img.shields.io/badge/API-doc-%2349cc90)](https://api.koios.rest/#get-/address_txs) |
| `/address_assets` | [`*.GetAddressAssets(...) *AddressAssetsResponse`](https://pkg.go.dev/github.com/cardano-community/koios-go-client#Client.GetAddressAssets) | `address-assets` | [![](https://img.shields.io/badge/API-doc-%2349cc90)](https://api.koios.rest/#get-/address_assets) |
| `/credential_txs` | [`*.GetCredentialTxs(...) *CredentialTxsResponse`](https://pkg.go.dev/github.com/cardano-community/koios-go-client#Client.GetCredentialTxs) | `credential-txs` | [![](https://img.shields.io/badge/API-doc-%2349cc90)](https://api.koios.rest/#get-/credential_txs) |
| **ACCOUNT** | | | |
| `/account_list` | [`*.GetAccountList(...) *AccountListResponse`](https://pkg.go.dev/github.com/cardano-community/koios-go-client#Client.GetAccountList) | `account-list` | [![](https://img.shields.io/badge/API-doc-%2349cc90)](https://api.koios.rest/#get-/account_list) |
| `/account_info` | [`*.GetAccountInfo(...) *AccountListResponse`](https://pkg.go.dev/github.com/cardano-community/koios-go-client#Client.GetAccountInfo) | `account-info` | [![](https://img.shields.io/badge/API-doc-%2349cc90)](https://api.koios.rest/#get-/account_info) |
| `/account_rewards` | [`*.GetAccountRewards(...) *AccountRewardsResponse`](https://pkg.go.dev/github.com/cardano-community/koios-go-client#Client.GetAccountRewards) | `account-rewards` | [![](https://img.shields.io/badge/API-doc-%2349cc90)](https://api.koios.rest/#get-/account_info) |
| `/account_updates` | [`*.GetAccountUpdates(...) *AccountUpdatesResponse`](https://pkg.go.dev/github.com/cardano-community/koios-go-client#Client.GetAccountUpdates) | `account-updates` | [![](https://img.shields.io/badge/API-doc-%2349cc90)](https://api.koios.rest/#get-/account_updates) |
| `/account_addresses` | [`*.GetAccountAddresses(...) *AccountAddressesResponse`](https://pkg.go.dev/github.com/cardano-community/koios-go-client#Client.GetAccountAddresses) | `account-addresses` | [![](https://img.shields.io/badge/API-doc-%2349cc90)](https://api.koios.rest/#get-/account_addresses) |
| `/account_assets` | [`*.GetAccountAssets(...) *AccountAssetsResponse`](https://pkg.go.dev/github.com/cardano-community/koios-go-client#Client.GetAccountAssets) | `account-assets` | [![](https://img.shields.io/badge/API-doc-%2349cc90)](https://api.koios.rest/#get-/account_assets) |
| `/account_history` | [`*.GetAccountHistory(...) *AccountHistoryResponse`](https://pkg.go.dev/github.com/cardano-community/koios-go-client#Client.GetAccountHistory) | `account-history` | [![](https://img.shields.io/badge/API-doc-%2349cc90)](https://api.koios.rest/#get-/account_history) |
| **ASSET** | | | |
| `/asset_list` | [`*.GetAssetList(...) *AssetAddressListResponse`](https://pkg.go.dev/github.com/cardano-community/koios-go-client#Client.GetAssetList) | `asset-list` | [![](https://img.shields.io/badge/API-doc-%2349cc90)](https://api.koios.rest/#get-/asset_address_list) |
| `/asset_address_list` | [`*.GetAssetAddressList(...) *AssetAddressListResponse`](https://pkg.go.dev/github.com/cardano-community/koios-go-client#Client.GetAssetAddressList) | `asset-address-list` | [![](https://img.shields.io/badge/API-doc-%2349cc90)](https://api.koios.rest/#get-/asset_address_list) |
| `/asset_info` | [`*.GetAssetInfo(...) *AssetInfoResponse`](https://pkg.go.dev/github.com/cardano-community/koios-go-client#Client.GetAssetInfo) | `asset-info` | [![](https://img.shields.io/badge/API-doc-%2349cc90)](https://api.koios.rest/#get-/asset_info) |
| `/asset_summary` | [`*.GetAssetSummary(...) *AssetSummaryResponse`](https://pkg.go.dev/github.com/cardano-community/koios-go-client#Client.GetAssetSummary) | `asset-summary` | [![](https://img.shields.io/badge/API-doc-%2349cc90)](https://api.koios.rest/#get-/asset_summary) |
| `/asset_txs` | [`*.GetAssetTxs(...) *AssetTxsResponse`](https://pkg.go.dev/github.com/cardano-community/koios-go-client#Client.GetAssetTxs) | `asset-txs` | [![](https://img.shields.io/badge/API-doc-%2349cc90)](https://api.koios.rest/#get-/asset_txs) |
| **POOL** | | | |
| `/pool_list` | [`*.GetPoolList(...) *PoolListResponse`](https://pkg.go.dev/github.com/cardano-community/koios-go-client#Client.GetPoolList) | `pool-list` | [![](https://img.shields.io/badge/API-doc-%2349cc90)](https://api.koios.rest/#get-/pool_list) |
| `/pool_info` | [`*.GetPoolInfos(...) *PoolInfoResponse`](https://pkg.go.dev/github.com/cardano-community/koios-go-client#Client.GetPoolInfos) | `pool-infos` | [![](https://img.shields.io/badge/API-doc-%2349cc90)](https://api.koios.rest/#get-/pool_info) |
| | [`*.GetPoolInfo(...) *PoolInfoResponse`](https://pkg.go.dev/github.com/cardano-community/koios-go-client#Client.GetPoolInfo) | `pool-info` | |
| `/pool_delegators` | [`*.GetPoolDelegators(...) *PoolDelegatorsResponse`](https://pkg.go.dev/github.com/cardano-community/koios-go-client#Client.GetPoolDelegators) | `pool-delegators` | [![](https://img.shields.io/badge/API-doc-%2349cc90)](https://api.koios.rest/#get-/pool_delegators) |
| `/pool_blocks` | [`*.GetPoolBlocks(...) *PoolBlocksResponse`](https://pkg.go.dev/github.com/cardano-community/koios-go-client#Client.GetPoolBlocks) | `pool-blocks` | [![](https://img.shields.io/badge/API-doc-%2349cc90)](https://api.koios.rest/#get-/pool_blocks) |
| `/pool_updates` | [`*.GetPoolUpdates(...) *PoolUpdatesResponse`](https://pkg.go.dev/github.com/cardano-community/koios-go-client#Client.GetPoolUpdates) | `pool-updates` | [![](https://img.shields.io/badge/API-doc-%2349cc90)](https://api.koios.rest/#get-/pool_updates) |
| `/pool_relays` | [`*.GetPoolRelays(...) *PoolRelaysResponse`](https://pkg.go.dev/github.com/cardano-community/koios-go-client#Client.GetPoolRelays) | `pool-relays` | [![](https://img.shields.io/badge/API-doc-%2349cc90)](https://api.koios.rest/#get-/pool_relays) |
| `/pool_metadata` | [`*.GetPoolMetadata(...) *PoolMetadataResponse`](https://pkg.go.dev/github.com/cardano-community/koios-go-client#Client.GetPoolMetadata) | `pool-metadata` | [![](https://img.shields.io/badge/API-doc-%2349cc90)](https://api.koios.rest/#get-/pool_metadata) |
| **SCRIPT** | | | |
| `/script_list` | [`*.GetScriptList(...) *ScriptRedeemersResponse`](https://pkg.go.dev/github.com/cardano-community/koios-go-client#Client.GetScriptList) | `script-list` | [![](https://img.shields.io/badge/API-doc-%2349cc90)](https://api.koios.rest/#get-/script_list) |
| `/script_redeemers` | [`*.GetScriptRedeemers(...) *ScriptListResponse`](https://pkg.go.dev/github.com/cardano-community/koios-go-client#Client.GetScriptRedeemers) | `script-redeemers` | [![](https://img.shields.io/badge/API-doc-%2349cc90)](https://api.koios.rest/#get-/script_redeemers) |

---

## Contributing

We would love for you to contribute to [Koios API Client Library for Go][github] and help make it even better than it is today! As a contributor, here are the guidelines we would like you to follow:

 - [Code of Conduct](#code-of-conduct)
 - [Question or Problem?](#got-a-question-or-problem)
 - [Found a Bug?](#issues-and-bugs)
 - [Missing a Feature?](#feature-requests)
 - [Submission Guidelines](#submission-guidelines)
 - [Coding Rules](#coding-rules)
 - [Commit Message Guidelines](#commit-message-guidelines)
 - [Development Documentation](#development-documentation)

### Code of Conduct

Help us keep [Koios API Client Library for Go][github] open and inclusive. Please read and follow our [Code of Conduct][coc]

---

### Got a Question or Problem?

Do not open issues for general support questions as we want to keep GitHub issues for bug reports and feature requests. You've got much better chances of getting your question answered on [Koios Telegram Group](https://t.me/joinchat/+zE4Lce_QUepiY2U1)

---

### Issues and Bugs

If you find a bug in the source code, you can help us by
[submitting an issue](#submitting-an-issue) to our [GitHub Repository][github]. Even better, you can
[submit a Pull Request](#submitting-a-pull-request-pr) with a fix.

---

### Feature Requests

You can *request* a new feature by [submitting an issue](#submitting-an-issue) to our GitHub
Repository. If you would like to *implement* a new feature, please submit an issue with
a proposal for your work first, to be sure that we can use it.
Please consider what kind of change it is:

* For a **Major Feature**, first open an issue and outline your proposal so that it can be
discussed. This will also allow us to better coordinate our efforts, prevent duplication of work,
and help you to craft the change so that it is successfully accepted into the project.
* **Small Features** can be crafted and directly [submitted as a Pull Request](#submitting-a-pull-request-pr).

---

### Submission Guidelines

#### Submitting an Issue

Before you submit an issue, please search the issue tracker, maybe an issue for your problem already exists and the discussion might inform you of workarounds readily available.

You can file new issues by filling out our [new issue form](https://github.com/cardano-community/koios-go-client/issues/new).

---

#### Submitting a Pull Request (PR)

Before you submit your Pull Request (PR) consider the following guidelines:

1. Search [GitHub](https://github.com/cardano-community/koios-go-client/pulls) for an open or closed PR that relates to your submission. You don't want to duplicate effort.
2. Fork the [cardano-community/koios-go-client][github] repo.
3. Setup you local repository

    ```shell
    git@github.com:<your-github-username>/koios-go-client.git
    cd koios-go-client
    git remote add upstream git@github.com:cardano-community/koios-go-client.git
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
6. Follow our [Coding Rules](#coding-rules).
7. If changes are in source code except documentations then run the full test suite, as described in the [developer documentation](#development-documentation), and ensure that all tests pass.
8.  Commit your changes using a descriptive commit message that follows our
  [commit message conventions](#commit-message-guidelines). Adherence to these conventions
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
fix(endpoint): update Tip endpoint to latest specs.

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
| **client** | API client related changes |
| **endpoint** | Changes related to api endpoints |
| **godoc** | Go documentation |
| **markdown** | Markdown files |
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
* Fork the [cardano-community/koios-go-client][github] repo.
* Setup you local repository

    ```shell
    git@github.com:<your-github-username>/koios-go-client.git
    cd koios-go-client
    git remote add upstream git@github.com:cardano-community/koios-go-client.git
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

## Credits

[![GitHub contributors](https://img.shields.io/github/contributors/cardano-community/koios-go-client?style=flat-square)](https://github.com/cardano-community/koios-go-client/graphs/contributors)

<sub>**Original author.**</sub>  
<sup>koios-go-client was moved under Cardano Community from <a href="https://github.com/howijd/koios-rest-go-client">howijd/koios-rest-go-client</a></sup>

<!-- LINKS -->
[Koios API]: https://koios.rest "Koios API"
[github.com/shopspring/decimal]: https://github.com/shopspring/decimal
[coc]: https://github.com/cardano-community/.github/blob/main/CODE_OF_CONDUCT.md
[github]: https://github.com/cardano-community/koios-go-client
[koios-cli]: https://github.com/cardano-community/koios-cli "cardano-community/koios-cli"


