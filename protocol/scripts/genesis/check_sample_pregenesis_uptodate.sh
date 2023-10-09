#!/bin/bash
# This script checks `sample_pregenesis.json` is up to date with the result of running `prod_pregenesis.sh`.
# Usage: % ./scripts/genesis/check_sample_pregenesis_uptodate.sh
# Currently can only be run locally on macOS.
# TODO(CORE-632): Make this script work on CI.
echo "Installing dasel..."
apk add dasel jq
echo "Building binary..."
make build
echo "Running prod_pregenesis.sh..."
./scripts/genesis/prod_pregenesis.sh build/dydxprotocold

diff_output=$(diff "/tmp/prod-chain/.dydxprotocol/config/genesis.json" "./scripts/genesis/sample_pregenesis.json")

if [ -z "$diff_output" ]; then
    echo "./scripts/genesis/sample_pregenesis.json is up-to-date"
else
    echo "./scripts/genesis/sample_pregenesis.json is not up-to-date"
    echo "$diff_output"
    exit 1
fi