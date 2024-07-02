<p align="center"><img src="https://dydx.exchange/icon.svg?" width="256" /></p>

<h1 align="center">v4 Indexer</h1>

<div align="center">
  <a href="https://github.com/dydxprotocol/v4-chain/actions/workflows/indexer-build-test-coverage.yml?query=branch%3Amain" style="text-decoration:none;">
    <img src="https://github.com/dydxprotocol/v4-chain/actions/workflows/indexer-build-test-coverage.yml/badge.svg?branch=main" />
  </a>
</div>

This monorepo holds a set of packages and services used to deploy the V4 Indexer.

The V4 Indexer is a set of microservices which provides a set of rich APIs and websocket channels
for a dYdX V4 application.

The services that make up the V4 Indexer are:
- ender: Receives on-chain events from V4 node via Kafka, archives them to a PostgreSQL database
and sends websocket events for the on-chain events to Kafka
- vulcan: Receives off-chain events from V4 node via Kafka, archives them to a PostgreSQL database
and sends websocket events for the off-chain events to Kafka
- socks: Websocket server that accepts subscriptions for websocket channels, and forwards websocket
events from Kafka to websocket clients that are subscribed
- comlink: Serves all HTTP API calls made to the V4 Indexer
- roundtable: Periodic job service
- bazooka: Deploys Indexer services via AWS Lambda

TODO(DEC-784): Add comprehensive descriptions/documentation for each service and the mono-repo as 
a whole.
## Setup

We use [pnpm](https://pnpm.io/) in this repo because it is a "faster and more disk efficient package manager" and it works better than `npm` with a monorepo.

```sh
nvm install
nvm use
nvm alias default $(nvm version) # optional
npm i -g pnpm@6
```

### Installation

Now, you can install dependencies for Indexer. This should also be run anytime packages are updated.

```
pnpm install
```

### Build

To build all services and packages, run:

```
pnpm run build:all
```
This should be run whenever code is changed, and you need to deploy or run the updated code, including running unit tests, deploying locally, or deploying to AWS.

## Adding Packages

Use `packages/example-package` as a template:
```
cp -r packages/example-package packages/<package-name>
```
- Update `package.json` name to `@dydxprotocol-indexer/<package-name>`, update the `README.md`, and run `pnpm i` to install dependencies.
- Add in `Dockerfile.service.local`, `Dockerfile.service.remote`, and `Dockerfile.postgres-package.local` to copy the `package.json` file and `build/` files.
- Add in `Dockerfile.bazooka.remote` to copy the package directory

## Adding Services

Use `services/example-service` as a template:
```
cp -r services/example-service services/<service-name>
```
- Update `package.json` name to `<package-name>`, update the `README.md`, and run `pnpm i` to install dependencies.
Add service deployment config to `docker-compose-local-deployment.yml`.
- Add Service to `v4-terraform` repo (TODO(DEC-838): Add link to the specific files)
- Add Github action in `.github/workflows/build-and-push.yml` to push image to ECR
- Add Service to be deployed to Orb so it will be shut down, updated, and redeployed
- Update Auxo to deploy service

## Running package.json scripts

`package.json` scripts for any service or package can be run with `pnpm run <script-name>` in the respective directory.

NOTE: `pnpm` allows running scripts across all directories using the `-r` flag or running them in parallel using the `--parallel` flag. Due to our unit tests sharing a test database, they cannot be run in parallel, so do not use either flag.

## Protos

Protos can be found in `proto/` [here](https://github.com/dydxprotocol/v4-chain/tree/main/proto).

## Running unit tests

First, make sure all services and packages are built with the latest code by running:

```
pnpm run build:all
```

Open up 2 terminals (or have another `tmux` or `screen` session) and run:

```
# In session / terminal 1
docker-compose up

# In session / terminal 2
pnpm run test:all
```

If you change any logic, you'll have to re-build the services and packages before running unit tests.

### To run a single test file:
`cd services/{service_name} && pnpm build && pnpm test -- {test_name}`

# Running Dockerfile locally
TODO(DEC-671): Add e2e tests
Deploying Indexer locally serves to run e2e tests (DEC-671) and to test new messages from a local deployment from v4-protocol.

## Deploying

Add to your ~/.zshrc file:
`export DD_API_KEY=<INSERT_DD_API_KEY_HERE>`

See https://app.datadoghq.com/organization-settings/api-keys for API keys. API Key is "Key", NOT "Key Id".

Indexer can be deployed locally with:
`docker-compose -f docker-compose-local-deployment.yml up`

If you want to export stats to Datadog, add the following flag to above command:
`--profile export-to-datadog`

Follow steps under [redeploying](#redeploying) if any changes have been made to services since the
last time indexer was deployed locally (e.g. after pulling from master/new branch).

NOTE: The kafka container may run into issues when starting up/initialization, if the kafka container is not
healthy after several minutes, Ctrl-C and run `docker-compose -f docker-compose-local-deployment.yml down` to
remove all the containers, then try running `docker-compose -f docker-compose-local-deployment.yml up` again.

### Running local V4 node
By default the Indexer services connect to a local V4 node. To run a local V4 node alongside the indexer, 
follow the instructions [here](https://github.com/dydxprotocol/v4#running-the-chain-locally).

NOTE: The local V4 node needs to be started up before deploying the Indexer locally.
### Connecting to remote V4 node
To change the locally deployed indexer to connect to a remote V4 node, update the environment variable
`TENDERMINT_WS_URL` for both the `comlink` service in `docker-compose-local-deployment.yml`.

e.g.
```
comlink:
    build:
      dockerfile: Dockerfile.service.local
      args:
        service: comlink
    ports:
      - 3002:3002
    links:
      - postgres
    depends_on:
      postgres-package:
        condition: service_completed_successfully
```

### Deploying to AWS Dev Environment
We use [Terraform](https://github.com/dydxprotocol/v4-infrastructure) to describe our infrastructure. On merging to
master branch, we automatically build ECR images of each service for both the dev/staging environments (see .github/workflows directory).
However, if you want to test and deploy local changes, you can use the `scripts/deploy-commit-to-dev.sh` script.

Example usage:
scripts/deploy-commit-to-dev.sh <service> <commit_hash>, e.g.
```
scripts/deploy-commit-to-dev.sh comlink 875aecd
```

You can get the 7 character commit hash like so:
```
git add .
git commit -m "<commit_msg>"
git rev-parse --short=7 HEAD
```

## Redeploying
The docker image for each container must be rebuilt for each code change or update to `Dockerfile.local`.

Images can be deleted individually:
```
# List all docker images
docker images
# Delete docker image
docker rmi {docker_image_id} -f
```
All Images can be deleted with this command:
```
docker rmi $(docker images 'indexer_*' -a -q) -f
```
The docker container can then be redeployed with the command in the previous section.

## Querying API of locally run Indexer
`comlink` will serving API requests from `http://localhost:3002`.

For example, to get all `BTC-USD` trades:
```
curl http://localhost:3002/v4/trades/perpetualMarket/BTC-USD
```

## Subscribing to websocket messages from Indexer
`socks` will accepting websocket connections from `http://localhost:3003`.

To connect to the dev/staging endpoints, use the following commands:
```
# dev
wscli connect ws://dev-indexer-apne1-lb-public-890774175.ap-northeast-1.elb.amazonaws.com/v4/ws

# staging
wscli connect wss://indexer.v4staging.dydx.exchange/v4/ws
```

Use a command-line websocket client such as 
[interactive-websocket-cli](https://www.npmjs.com/package/interactive-websocket-cli) to connect and
subscribe to channels.

Example (with `interactive-websocket-cli`)
```
wscli connect http://localhost:3003
<output from ws-cli>
<type 's' to send> { "type": "subscribe", "channel": "v4_trades", "id": "BTC-USD" }
```

Other example subscription events:
```
{ "type": "subscribe", "channel": "v4_candles", "id": "BTC-USD/1MIN" }
{ "type": "subscribe", "channel": "v4_markets" }
{ "type": "subscribe", "channel": "v4_orderbook", "id": "BTC-USD" }
{ "type": "subscribe", "channel": "v4_subaccounts", "id": "address/0" }
{ "type": "subscribe", "channel": "v4_block_height" }
```
