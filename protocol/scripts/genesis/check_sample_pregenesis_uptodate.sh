#!/bin/bash
# This script checks `sample_pregenesis.json` is up to date with the result of running `prod_pregenesis.sh`.
# Usage: make check-sample-pregenesis-up-to-date

echo "Running prod_pregenesis.sh..."
./scripts/genesis/prod_pregenesis.sh dydxprotocold

echo "Diffing output against current sample_pregenesis.json..."
diff_output=$(diff "/tmp/prod-chain/.dydxprotocol/config/sorted_genesis.json" "./scripts/genesis/sample_pregenesis.json")

if [ -z "$diff_output" ]; then
    echo "./scripts/genesis/sample_pregenesis.json is up-to-date"
else
    echo "./scripts/genesis/sample_pregenesis.json is not up-to-date. Run `make update-sample-pregenesis` to update."
    echo "$diff_output"
    exit 1
fi
