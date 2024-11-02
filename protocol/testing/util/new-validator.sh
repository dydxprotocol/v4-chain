#!/bin/bash

# This is a simple script that outputs the information for a new test validator.
# This information can be used to add a new validator to `dev.sh` or `staging.sh`.

echo "Generating new validator information.."

apk add dasel jq >/dev/null 2>&1

mkdir ./tmp-new-validator
klyraprotocold init tmp --chain-id klyraprotocol --home ./tmp-new-validator >/dev/null 2>&1

# Take the last line of the output which contains the new mnemonic
mnemonic=$(klyraprotocold keys add val --home ./tmp-new-validator --keyring-backend=test 2>&1 >/dev/null | tail -1)

# Overwrite the randomly generated `priv_validator_key.json` with a key generated deterministically from the mnemonic.
klyraprotocold tendermint gen-priv-key --home ./tmp-new-validator --mnemonic "$mnemonic"
address=$(klyraprotocold keys show val --address --keyring-backend=test --home ./tmp-new-validator)
node_key=$(jq -r ".priv_key.value" ./tmp-new-validator/config/node_key.json)
node_id=$(klyraprotocold tendermint show-node-id --home ./tmp-new-validator)
cons_address=$(klyraprotocold tendermint show-address --home ./tmp-new-validator)

echo "New Validator Information:"
echo "mnemonic: $mnemonic"
echo "account address: $address"
echo "node_key: $node_key"
echo "node_id: $node_id"
echo "consensus address: $cons_address"

rm -rf ./tmp-new-validator
