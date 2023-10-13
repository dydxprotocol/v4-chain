# Overview

For development and testing purposes, it's useful to be able to deploy the chain remotely where the validators are running on cloud infra. This directory contains setup scripts and Dockerfile that can be used to build Docker container images for running validator nodes.

Overall, the steps involve:
1. Define validator node keys
2. Build the app binary (i.e. `dydxprotocold`)
3. Set up configuration files (i.e. `genesis.json`, `app.toml`, `app.toml`, and etc.)
4. Setup the app data and configs (i.e. moving data and configs files into correct directory)

Currently, step 1 is a manual process and steps 2, 3 and 4 are automated within the scripts (`local.sh`, `dev.sh`, and `staging.sh`).

TODO(CORE-512): add info/resources around cloud setup. [Doc](https://www.notion.so/dydx/Cloud-Setup-ccc68e7b6a4b4e83a8d0c35029a9997f)

# [Manual] Add new validator node keys to `dev.sh` or `staging.sh`

In order to add a new validator node to `dev.sh` or `staging.sh`, simply run `make localnet-new-validator` and copy the `mnemonic`, `address`, `node_key`, `node_id` in the their respective fields in `dev.sh` or `staging.sh`, and add a new `moniker` for this validator.

# Running the `local` chain

The local chain can be run using the `docker-compose.yml` file at the root of this repository. Simply run `make localnet-start` to start the validator and `make localnet-stop` to stop it.

# Running the `dev` or `staging` chains locally

You can run the testnet chain locally by running the following command from the `v4` repository root.

It's necessary to specify the `--home` flag as this is how the container knows which validator to run as. A new home directory is created for each `MONIKER` defined in `dev.sh` or `staging.sh`. There is also an additional home directory for running the chain as a full-node (see instructions below).

```sh
# dev
docker build . --progress=plain --no-cache -f ./testing/testnet-dev/Dockerfile -t testnet && docker run testnet start --home /dydxprotocol/chain/.alice

# staging
docker build . --progress=plain --no-cache -f ./testing/testnet-staging/Dockerfile -t testnet && docker run testnet start --home /dydxprotocol/chain/.alice
```

# Building and Pushing the Docker container image to ECR

Currently the `staging` container is pushed automatically to ECR when new changes are pushed to the `main` branch, however you can also manually deploy the container.

Run the `./build-push-ecr.sh` script in this directory from this repository's root directory. You can specify either the "dev", "dev2", "dev3", "dev4", "dev5", or "staging" environment.

Ensure that your `AWS_PROFILE` is correct for the environment you want to deploy to (this can be found in ~/.aws/credentials), and can be selected using the `AWS_PROFILE` env var.

example: `AWS_PROFILE=dydx-v4 ./scripts/build-push-ecr.sh staging`

# Opening a shell on Docker container images

It's useful to open up a shell on Docker container images to inspect or interact. Run the following command to do so:

```sh
$ docker run -it --entrypoint /bin/sh <image id>
```

# Running a full-node

If you wish to run as a full-node instead of a validator, specify the `--home` flag as ` /dydxprotocol/chain/.full-node`.

```sh
$ docker build . --progress=plain --no-cache -f ./testing/testnet-dev/Dockerfile -t testnet && docker run testnet start --home /dydxprotocol/chain/.full-node
```
