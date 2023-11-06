#!/bin/bash

# This script takes in the result of a `/dydxprotocol/clob/clob_pair` query, and generates 
# a sample governance proposal to enable trading on all CLOB pairs. 
# Example usage:
# 1. Get all clob pairs from a REST endpoint:
#   % curl -X GET "https://dydx-testnet-archive.allthatnode.com:1317/dydxprotocol/clob/clob_pair" -H "accept: application/json" > /tmp/clob_pairs.json
# 2. Generate proposal JSON file:
#   % ./scripts/governance/enable_all_clob_pairs.sh /tmp/clob_pairs.json > /tmp/proposal_enable_trading_all_markets.json
# 3. Submit proposal:
#   % dydxprotocold tx gov submit-proposal /tmp/proposal_enable_trading_all_markets.json --from alice --gas auto --fees 400000000000000000adv4tnt

# Constants
NINE_ZEROS="000000000"
AUTHORITY="dydx10d07y265gmmuvt4z0w9aw880jnsr700jnmapky"

# Customizable proposal fields
TITLE="Enable trading on all markets"
NATIVE_TOKEN_DENOM="adv4tnt"
DEPOSIT="10000${NINE_ZEROS}${NINE_ZEROS}${NATIVE_TOKEN_DENOM}" # 10,000 native tokens
SUMMARY="Use MsgUpdateClobPair to change the status of all CLOB pairs to ACTIVE. All other fields remain unchanged."

if [ -z "$1" ]; then
  echo "Usage: $0 <input_json_file>"
  exit 1
fi

INPUT_JSON="$1"


# Use jq to construct the messages array from the input JSON
MESSAGES=$(jq --arg authority "$AUTHORITY" '
  .clob_pair | map({
    "@type": "/dydxprotocol.clob.MsgUpdateClobPair",
    "authority": $authority,
    "clob_pair": {
      "id": .id,
      "perpetual_clob_metadata": {
        "perpetual_id": .perpetual_clob_metadata.perpetual_id
      },
      "quantum_conversion_exponent": .quantum_conversion_exponent,
      "status": .status,
      "step_base_quantums": (.step_base_quantums | tonumber),
      "subticks_per_tick": (.subticks_per_tick | tonumber)
    }
  })
' "$INPUT_JSON")

FINAL_JSON=$(jq -n --argjson messages "$MESSAGES" --arg title "$TITLE" --arg deposit "$DEPOSIT" --arg summary "$SUMMARY" '
  {
    "title": $title,
    "deposit": $deposit,
    "summary": $summary,
    "messages": $messages
  }
')

echo "$FINAL_JSON"
