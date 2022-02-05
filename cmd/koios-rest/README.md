# Koios API CLI Application

```
go install github.com/howijd/koios-rest-go-client/cmd/koios-rest@latest
```

**see cmd help for available commands and flags**

```
koios-rest -h
```

```
NAME:
   koios-rest - CLI Client to consume Koios API https://api.koios.rest

USAGE:
   koios-rest [global options] command [command options] [arguments...]

VERSION:
   0.1.0

AUTHOR:
   The Howijd.Network Authors
COMMANDS:
   help, h  Shows a list of commands or help for one command
   ACCOUNT:
     account-list       Get a list of all accounts.
     account-info       Get the account info of any (payment or staking) address.
     account-rewards    Get the full rewards history (including MIR) for a stake address, or certain epoch if specified.
     account-updates    Get the account updates (registration, deregistration, delegation and withdrawals).
     account-addresses  Get all addresses associated with an account.
     account-assets     Get the native asset balance of an account.
     account-history    Get the staking history of an account.
   ADDRESS:
     address-info    Get address info - balance, associated stake address (if any) and UTxO set.
     address-txs     Get the transaction hash list of input address array, optionally filtering after specified block height (inclusive).
     address-assets  Get the list of all the assets (policy, name and quantity) for a given address.
     credential-txs  Get the transaction hash list of input payment credential array, optionally filtering after specified block height (inclusive)
   BLOCK:
     blocks      Get summarised details about all blocks (paginated - latest first).
     block-info  Get detailed information about a specific block.
     block-txs   Get a list of all transactions included in a provided block.
   DEVELOPMENT COMMANDS:
     dev  koios-rest-go-client development commands.
   EPOCH:
     epoch-info    Get the epoch information, all epochs if no epoch specified.
     epoch-params  Get the protocol parameters for specific epoch, returns information about all epochs if no epoch specified.
   NETWORK:
     tip      Get the tip info about the latest block seen by chain.
     genesis  Get the Genesis parameters used to start specific era on chain.
     totals   Get the circulating utxo, treasury, rewards, supply and reserves in lovelace for specified epoch, all epochs if empty.
   POOL:
     pool-list        A list of all currently registered/retiring (not retired) pools.
     pool-info        Current pool statuses and details for a specified list of pool ids.
     pool-delegators  Return information about delegators by a given pool and optional epoch (current if omitted).
     pool-blocks      Return information about blocks minted by a given pool in current epoch (or _epoch_no if provided).
     pool-updates     Return all pool updates for all pools or only updates for specific pool if specified.
     pool-relays      A list of registered relays for all currently registered/retiring (not retired) pools.
     pool-metadata    Metadata(on & off-chain) for all currently registered/retiring (not retired) pools.
   SCRIPT:
     script-list       List of all existing script hashes along with their creation transaction hashes.
     script-redeemers  List of all redeemers for a given script hash.
   TRANSACTIONS:
     tx-info        Get detailed information about transaction(s).
     tx-utxos       Get UTxO set (inputs/outputs) of transactions.
     tx-metadata    Get metadata information (if any) for given transaction(s)..
     tx-metalabels  Get a list of all transaction metalabels.
     submittx       Submit an already serialized transaction to the network.
     tx-status      Get the number of block confirmations for a given transaction hash list

GLOBAL OPTIONS:
   --port value, -p value  Set port (default: 443)
   --host value            Set host (default: "api.koios.rest")
   --api-version value     Set API version (default: "v0")
   --schema value          Set URL schema (default: "https")
   --origin value          Set Origin header for requests. (default: "https://github.com/howijd/koios-rest-go-client")
   --rate-limit value      Set API Client rate limit for outgoing requests (default: 5)
   --ugly                  Ugly prints response json strings directly without calling json pretty. (default: false)
   --help, -h              show help (default: false)
   --version, -v           print the version (default: false)

COPYRIGHT:
   (c) 2022

```
