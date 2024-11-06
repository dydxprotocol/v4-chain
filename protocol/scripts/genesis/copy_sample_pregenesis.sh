#!/bin/bash

./scripts/genesis/prod_pregenesis.sh klyraprotocold
cp /tmp/prod-chain/.klyraprotocol/config/sorted_genesis.json ./scripts/genesis/sample_pregenesis.json