<p align="center"><img src="https://dydx.exchange/icon.svg?" width="256" /></p>

<h1 align="center">v4 Protocol</h1>

<div align="center">
  <a href="https://github.com/dydxprotocol/v4-chain/actions/workflows/protocol-test.yml?query=branch%3Amain" style="text-decoration:none;">
    <img src="https://github.com/dydxprotocol/v4-chain/actions/workflows/protocol-test.yml/badge.svg?branch=main" />
  </a>
</div>

Sovereign Blockchain built using Cosmos SDK & CometBFT for v4 of the dYdX Protocol.

TODO(CORE-512): add info/resources around V4. [Doc](https://www.notion.so/dydx/V4-36a9f30eee1d478cb88e0c50860fdbee)

## Get started

### Installation

1. Install Go `1.19` [here](https://go.dev/dl/).
2. Install Docker Desktop [here](https://www.docker.com/products/docker-desktop/)
3. Run `make install` to install project dependencies.

### Helpful Development Commands

* `make test` runs unit tests.
* `make lint` lints source code (`make lint-fix` to also fix).
* `make proto-all` formats, lints, and generates Go from proto files.
* `make build` builds source.
* `make mock-gen` generates mocks for files listed in [mocks/Makefile](https://github.com/dydxprotocol/v4/tree/main/mocks/Makefile). More info about mocking [here](https://github.com/dydxprotocol/v4/tree/main/mocks/README.md).


### Running the chain locally

Requirements: Ensure you are running docker-compose version `1.30.x` or newer, and Docker engine version `20.10.x` or newer.

You can quickly test your changes to DYDX with just a few commands:

1. Make any change to the DYDX code that you want to test

2. Once ready to test, run `make localnet-start` (or `make localnet-startd` to run the network headlessly)
    - This first compiles all your changes to docker image called `dydxprotocol-base` (~90 seconds)
    - You will then be running a local network with your changes!
    - Note, these commands will **reset the chain** to genesis

3. To remove all block history and start from scratch, re-run `make localnet-start` or `make localnet-startd`.

4. To stop the chain but keep the state, run `make localnet-stop`. To restart the protocol with the previous state, run `make localnet-continue`.


#### Deployment to AWS testnets

Merges to the `main` branch automatically trigger a new Docker container image to be built and pushed to ECR. After the image has been pushed to ECR, a Terraform Cloud run is currently necessary to deploy the new container to ECS.

The following commands can be used to locally build and push containers to ECR.

* `make aws-push-dev` locally build and push a container image to the "dev" environment.
* `make aws-push-dev2` locally build and push a container image to the "dev2" environment.
* `make aws-push-dev3` locally build and push a container image to the "dev3" environment.
* `make aws-push-dev4` locally build and push a container image to the "dev4" environment.
* `make aws-push-dev5` locally build and push a container image to the "dev5" environment.
* `make aws-push-staging` locally build and push a container image to the "staging" environment.

#### Linting

We use [`yamllint`](https://github.com/adrienverge/yamllint) for linting YAML files. Instructions for using specific
`yamllint` actions are linked below:
- [Installing `yamllint`](https://yamllint.readthedocs.io/en/latest/quickstart.html#installing-yamllint).
- [Running `yamllint`](https://yamllint.readthedocs.io/en/latest/quickstart.html#running-yamllint).
- [Configuring `yamllint`](https://yamllint.readthedocs.io/en/latest/configuration.html).

We currently lint the following YAML files in the [`Lint` CI job](https://github.com/dydxprotocol/v4/blob/c5ec83f074b4ff997d71a6f5dc486579ea112600/.github/workflows/lint.yml):
- `.golangci.yml`.
- `.github/workflows/*.yml`.
  - Note this includes all files that end in the `yml` file extension in the `.github/workflows` directory.
- `buf.work.yaml`.

#### Protos

Protos can be found in `proto/` [here](https://github.com/dydxprotocol/v4-chain/tree/main/proto).

#### Genesis

You can find the local chain's genesis data in `testing/genesis.sh`. This dictates the starting app state of your chain when running `make localnet-start`. We currently start the chain with BTC and ETH perpetuals and prices but could easily add another perpetual and market like so:

```
...
	dasel put string -f "$GENESIS" '.app_state.perpetuals.perpetuals.[2].ticker' 'LINK-USD'
	dasel put int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[2].market_id' '1'
	dasel put int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[2].atomic_resolution' -v '-9'
	dasel put int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[2].default_funding_ppm' '0'
	dasel put int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[2].initial_margin_ppm' '50000'     # 5 %
	dasel put int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[2].maintenance_fraction_ppm' '600000' # 3 % (60% of IM)
...
	dasel put string -f "$GENESIS" '.app_state.prices.markets.[2].pair' 'LINK-USD'
	dasel put int -f "$GENESIS" '.app_state.prices.markets.[2].exponent' -v '-6'
	dasel put int -f "$GENESIS" '.app_state.prices.markets.[2].min_exchanges' '1'
	dasel put int -f "$GENESIS" '.app_state.prices.markets.[2].min_price_change_ppm' '50'
	dasel put int -f "$GENESIS" '.app_state.prices.markets.[2].price' '3000000000' # $3,000 = 1 ETH.
	dasel put int -f "$GENESIS" '.app_state.prices.markets.[2].exchanges.[0]' '0'
```

Another module that can be modified similarly is the subaccounts. We can add another subaccount by updating the `TEST_ACCOUNTS` array in `./testing/local.sh`.

## Debugging Tips

### Setting up keychain
To run the below commands, you'll want to import the private keys of the test accounts specified in [testnet-local/local.sh](https://github.com/dydxprotocol/v4/blob/main/testing/testnet-local/local.sh). Run the following commands and input the corresponding 12-word string from `MNEMONICS`. The resulting address should match those in `TEST_ACCOUNTS`.

```sh
./build/dydxprotocold keys add alice --recover

./build/dydxprotocold keys add bob --recover
```

### Send a test transaction locally
It's occasionally helpful to send a transaction to the local chain to observe Cosmos behavior through the API such as events. Until `clob` `v0.1` is complete, you can use the default Cosmos `bank` module to transfer assets between two accounts defined at genesis in the `genesis.sh` file.

```sh
./build/dydxprotocold tx bank send dydx199tqg4wdlnu4qjlxchpd7seg454937hjrknju4 dydx10fx7sy6ywd5senxae9dwytf8jxek3t2gcen2vs 100usdc
```

### Placing a test order locally

It's occasionally helpful to send a transaction to the local chain to test order placement and matching. Run the following two commands in succession in order to match an order between two accounts.

```sh
./build/dydxprotocold tx clob place-order dydx199tqg4wdlnu4qjlxchpd7seg454937hjrknju4 0 0 0 1 10 10000 20 --from alice --chain-id localdydxprotocol
./build/dydxprotocold tx clob place-order dydx10fx7sy6ywd5senxae9dwytf8jxek3t2gcen2vs 0 0 0 2 10 10000 20 --from bob --chain-id localdydxprotocol
```

Run the following command to cancel an order.

```sh
./build/dydxprotocold tx clob cancel-order dydx199tqg4wdlnu4qjlxchpd7seg454937hjrknju4 10 0 20 --from alice
```

### Querying the chain locally

While running the development server via `make localnet-start`, you can make queries locally using the Tendermint API. All endpoints listed [here](https://docs.tendermint.com/v0.37/rpc/#/Info/block) are supported. For example to get the block at height 2: `curl -X GET "localhost:26657/block?height=2"`.

### Updating local flags
When debugging or inspecting behavior of the chain locally, you may wish modify the flags passed to `dydxprotocold`. You can achieve this by modifying your `docker-compose.yml` file locally in the `entrypoint` section to change these passed in flags.

### Enabled more verbose logging locally
Refer to the section above and change the `log_level` to `trace`. Note that `trace` can be pretty noisy as it logs every block proposal, message, and committed block to stdout.

### Debugging behavior in Cosmos SDK
It's occasionally useful to be able to output logs or modify behavior in `cosmos-sdk` itself. To do this, check out [cosmos-sdk](https://github.com/cosmos/cosmos-sdk) locally at the branch which represents the version specified in the [`go.mod` file](https://github.com/dydxprotocol/v4/blob/main/go.mod) in this repository.

```sh
git clone git@github.com:cosmos/cosmos-sdk.git
git checkout v0.47.0-alpha2
```

After you've cloned the repo, you can modify the `go.mod` file in `v4` locally to include a [_replace directive_](https://go.dev/ref/mod#go-mod-file-replace) which locally points `cosmos-sdk` to your local version of `cosmos-sdk`. Example:

```diff
replace (
+	github.com/cosmos/cosmos-sdk v0.47.0-alpha2 => /Users/bryce/projects/cosmos-sdk
	github.com/keybase/go-keychain => github.com/99designs/go-keychain v0.0.0-20191008050251-8e49817e8af4
	google.golang.org/grpc => google.golang.org/grpc v1.50.1
)
```

Now running `make localnet-start` will include any changes you've made to your `cosmos-sdk` repository locally.

### Slowing down blocks
It may be useful for debugging purposes to change the speed at which blocks are committed. This can be useful for more easily sending multiple transactions within the same block, or just for reducing the noise of the `--verbose` output (which includes output every time a new block is committed).

For example, for a block time of `1 minute`, add the following lines to your `local.sh` file locally in this repository:

```yaml
# How long we wait for a proposal block before prevoting nil
dasel put string -f "$CONFIG_FOLDER"/config.toml '.consensus.timeout_propose' '60s'

# How long we wait after committing a block, before starting on the new
# height (this gives us a chance to receive some more precommits, even
# though we already have +2/3).
dasel put string -f "$CONFIG_FOLDER"/config.toml '.consensus.timeout_commit' '60s'
```

## Bootstrapping a new module
Though we do not use `ignite` to run our chain locally, you can still use [ignite's](https://docs.ignite.com) scaffolding tools to bootstrap a new module. You will need to install `ignite` first. This is good as a starting point, but will require you to make some manual changes. Example:

```sh
ignite scaffold module prices --require-registration
ignite scaffold message UpdatePrices prices --module prices
ignite scaffold map price name:string exponent:int min_exchanges:uint min_price_change:uint value:uint  --no-message --module prices
```

* You will likely need to manually modify protos as the `ignite scaffold` command doesn't allow complex types such as repeated fields or references to other structs.
* Do not check the openapi, or vue-generated code.
* The `ignite` tool may place your generated protos in the wrong directory. Ensure they are in the `proto/dydxprotocol/yourmodule` directory, and the package in the proto is `dydxprotocol.yourmodule`;
* The `ignite` tool will generate your protos with linting errors, so you will need to run `make proto-lint` and fix any errors.
* The `ignite` tool will generate your code with linting errors, so you will need to run `make lint` and fix any errors.

See [this](https://github.com/dydxprotocol/v4/pull/30) PR for a commit-by-commit example of bootstrapping a new module.

## CometBFT fork

Our current implementation contains a light fork of CometBFT. The fork can be found [here](https://github.com/dydxprotocol/cometbft). Instructions to update the fork are included there.

## CosmosSDK fork

Our current implementation contains a light fork of CosmosSDK. The fork can be found [here](https://github.com/dydxprotocol/cosmos-sdk). Instructions to update the fork are included there.

## Daemons

Daemons are background processes that run in go-routines to do asyncronous work. Daemons can be configured with their respective flags, e.g. `price-daemon-enabled` or `price-daemon-delay-loop-ms`.

### Pricefeed Daemon

The Pricefeed Daemon is responsible for ingesting prices from 3rd party exchanges like Binance and sending these prices to the application where they are then used by the Prices module. The Pricefeed daemon is started by default when the application starts.

The path to the configuration file used for querying exchanges is dictated by the flag `pricefeed-exchange-config-file`. The default configuration file is located within `v4` at the path `daemons/configs/default_pricefeed_exchange_config.go`. The file utilizes a template that is populated from `daemons/pricefeed/client/constants/static_exchange_startup_config.go`. The first exchange within the default configuration file additionally has inline comments explaining what each field is for.

## Learn more

- [Tutorials](https://docs.ignite.com/guide)
- [Cosmos SDK docs](https://docs.cosmos.network)
- [Developer Chat](https://discord.gg/H6wGTY8sxw)
