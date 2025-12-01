#!/bin/bash
# Submit and vote on a governance proposal from proposal.json
# Pure bash implementation

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Parse arguments
NETWORK="${1:-local}"
FROM_ACCOUNT="${2:-alice}"
PROPOSAL_FILE="${3:-proposal.json}"
VOTE_OPTION="${4:-yes}"

echo "============================================================"
echo "Submit and Vote on Governance Proposal"
echo "============================================================"
echo "Network: $NETWORK"
echo "From Account: $FROM_ACCOUNT"
echo "Proposal File: $PROPOSAL_FILE"
echo "Vote Option: $VOTE_OPTION"
echo "============================================================"
echo ""

# Submit the proposal
echo "Submitting proposal..."
if [ "$NETWORK" = "testnet" ]; then
    # Testnet uses dydxcli and adydx
    dydxprotocold tx gov submit-proposal "$PROPOSAL_FILE" \
        --from "$FROM_ACCOUNT" \
        --yes \
        --broadcast-mode sync \
        --gas auto \
        --fees 200000000000adydx \
        --keyring-backend test
else
    # Local/staging uses dydxprotocold and adv4tnt
    dydxprotocold tx gov submit-proposal "$PROPOSAL_FILE" \
        --from "$FROM_ACCOUNT" \
        --keyring-backend test \
        --fees 5000000000000000adv4tnt \
        --yes
fi

echo "✓ Proposal submitted!"
echo ""

# Always auto-vote
echo "Waiting 5 seconds for proposal to be indexed..."
sleep 5

# Get the latest proposal ID from remote node
echo "Getting latest proposal ID..."
if [ "$NETWORK" = "testnet" ]; then
    NODE="https://dydx-testnet-rpc.kingnodes.com"
    CHAIN_ID="dydx-testnet-4"
elif [ "$NETWORK" = "staging" ]; then
    NODE="https://dydx-ops-rpc.kingnodes.com"
    CHAIN_ID="dydx-staging-1"
else
    # Local
    NODE="http://localhost:26657"
    CHAIN_ID="dydxprotocol"
fi

PROPOSAL_ID=$(dydxprotocold query gov proposals \
    --node="$NODE" \
    --chain-id="$CHAIN_ID" \
    --output json \
    --limit 1 \
    --reverse 2>/dev/null | jq -r '.proposals[0].id // empty')

if [ -z "$PROPOSAL_ID" ]; then
    echo "⚠️  Could not retrieve proposal ID from $NODE"
    echo "    Please check the node is accessible or vote manually."
    exit 0
fi

echo "✓ Proposal ID: $PROPOSAL_ID"
echo ""
echo "============================================================"
echo "Voting '$VOTE_OPTION' on Proposal $PROPOSAL_ID"
echo "============================================================"
echo ""

# Reuse the vote_in_* scripts
if [ "$NETWORK" = "testnet" ]; then
    "$SCRIPT_DIR/vote_in_testnet.sh" "$PROPOSAL_ID" "$VOTE_OPTION"
elif [ "$NETWORK" = "staging" ]; then
    "$SCRIPT_DIR/vote_in_staging.sh" "$PROPOSAL_ID" "$VOTE_OPTION"
else
    # Local/dev
    "$SCRIPT_DIR/vote_in_dev.sh" "$PROPOSAL_ID" "$VOTE_OPTION"
fi

echo ""
echo "============================================================"
echo "✓ Voting complete!"
echo "============================================================"

echo ""
echo "✓ Done!"
echo ""
echo "============================================================"
echo "Usage examples:"
echo "  Local:      $0 local alice proposal.json"
echo "  Testnet:    $0 testnet dydx-1 proposal.json"
echo "  Staging:    $0 staging alice proposal.json"
echo "  Vote 'no':  $0 local alice proposal.json no"
echo "============================================================"
