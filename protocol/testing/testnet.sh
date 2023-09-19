# Creates the genesis configuration and keys for a testnet validator node.
# This function expects the following global variables:
#  * MONIKERS
#  * CHAIN_ID
#  * MNEMONICS
#  * TEST_ACCOUNTS
#  * FAUCET_ACCOUNTS
#  * NATIVE_TOKEN
#  * USDC_DENOM
#
# Args:
#   i - the index of the validator to create.
create_validator() {
  i="$1"
  local VAL_HOME_DIR="$HOME/chain/.${MONIKERS[$i]}"
  local VAL_CONFIG_DIR="$VAL_HOME_DIR/config"

  # Initialize the chain and validator files.
  dydxprotocold init "${MONIKERS[$i]}" -o --chain-id=$CHAIN_ID --home "$VAL_HOME_DIR"

  # Overwrite the randomly generated `priv_validator_key.json` with a key generated deterministically from the mnemonic.
  dydxprotocold tendermint gen-priv-key --home "$VAL_HOME_DIR" --mnemonic "${MNEMONICS[$i]}"

  # Note: `dydxprotocold init` non-deterministically creates `node_id.json` for each validator.
  # This is inconvenient for persistent peering during testing in Terraform configuration as the `node_id`
  # would change with every build of this container.
  #
  # For that reason we overwrite the non-deterministically generated one with a deterministic key defined in this file here.
  local new_file=$(jq ".priv_key.value = \"${NODE_KEYS[$i]}\"" "$VAL_CONFIG_DIR"/node_key.json)
  cat <<<"$new_file" >"$VAL_CONFIG_DIR"/node_key.json

  edit_config "$VAL_CONFIG_DIR"

  echo "${MNEMONICS[$i]}" | dydxprotocold keys add "${MONIKERS[$i]}" --recover --keyring-backend=test --home "$VAL_HOME_DIR"

  # Using "*" as a subscript results in a single arg: "dydx1... dydx1... dydx1..."
  # Using "@" as a subscript results in separate args: "dydx1..." "dydx1..." "dydx1..."
  # Note: `edit_genesis` must be called before `add-genesis-account`.
  edit_genesis "$VAL_CONFIG_DIR" "${TEST_ACCOUNTS[*]}" "${FAUCET_ACCOUNTS[*]}" "" ""
  update_genesis_use_test_volatile_market "$VAL_CONFIG_DIR"
  update_genesis_complete_bridge_delay "$VAL_CONFIG_DIR" "600"

  local acct
  for acct in "${TEST_ACCOUNTS[@]}"; do
    dydxprotocold add-genesis-account "$acct" 100000000000000000$USDC_DENOM,100000000000$NATIVE_TOKEN --home "$VAL_HOME_DIR"
  done
  for acct in "${FAUCET_ACCOUNTS[@]}"; do
    dydxprotocold add-genesis-account "$acct" 900000000000000000$USDC_DENOM,100000000000$NATIVE_TOKEN --home "$VAL_HOME_DIR"
  done

  dydxprotocold gentx "${MONIKERS[$i]}" 500000000$NATIVE_TOKEN --moniker="${MONIKERS[$i]}" --keyring-backend=test --chain-id=$CHAIN_ID --home "$VAL_HOME_DIR"

  # Copy the gentx to a shared directory.
  cp -a "$VAL_CONFIG_DIR/gentx/." /tmp/gentx
}

# Creates the genesis configuration for a testnet full node.
# This function expects the following global variables:
#  * CHAIN_ID
#  * FULL_NODE_KEYS
#
# Args:
#   i - the index of the validator to create.
create_full_node() {
  local i="$1"
  local FULL_NODE_HOME_DIR="$HOME/chain/.full-node-$i"
  local FULL_NODE_CONFIG_DIR="$FULL_NODE_HOME_DIR/config"
  dydxprotocold init "full-node" -o --chain-id=$CHAIN_ID --home "$FULL_NODE_HOME_DIR"

  # Note: `dydxprotocold init` non-deterministically creates `node_id.json` for each validator.
  # This is inconvenient for persistent peering during testing in Terraform configuration as the `node_id`
  # would change with every build of this container.
  #
  # For that reason we overwrite the non-deterministically generated one with a deterministic key defined in this file here.
  local new_file=$(jq ".priv_key.value = \"${FULL_NODE_KEYS[$i]}\"" "$FULL_NODE_CONFIG_DIR"/node_key.json)
  cat <<<"$new_file" >"$FULL_NODE_CONFIG_DIR"/node_key.json

  edit_config "$FULL_NODE_CONFIG_DIR"
}