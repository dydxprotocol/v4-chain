#!/bin/bash
# Submit a governance proposal from proposal.json
# Usage: 
#   ./submit_proposal.sh [network] [from_account] [proposal_file]
#   
#   network: local (default), testnet, staging
#   from_account: alice (default for local/staging), dydx-1 (default for testnet)
#   proposal_file: proposal.json (default)

set -e

NETWORK="${1:-local}"
FROM_ACCOUNT="${2:-alice}"
PROPOSAL_FILE="${3:-proposal.json}"

echo "============================================================"
echo "Submit Governance Proposal"
echo "============================================================"
echo "Network: $NETWORK"
echo "From Account: $FROM_ACCOUNT"
echo "Proposal File: $PROPOSAL_FILE"
echo "============================================================"
echo ""

if [ "$NETWORK" = "testnet" ]; then
    # Testnet uses dydxcli and adydx
    echo "Submitting to testnet..."
    dydxprotocold tx gov submit-proposal "$PROPOSAL_FILE" \
        --from "$FROM_ACCOUNT" \
        --yes \
        --broadcast-mode sync \
        --gas auto \
        --fees 200000000000adydx \
        --keyring-backend test
else
    # Local/staging uses dydxprotocold and adv4tnt
    echo "Submitting to $NETWORK..."
    dydxprotocold tx gov submit-proposal "$PROPOSAL_FILE" \
        --from "$FROM_ACCOUNT" \
        --keyring-backend test \
        --fees 500000000000000adv4tnt \
        --yes
fi

echo ""
echo "✓ Proposal submitted!"
echo ""
echo "============================================================"
echo "Examples:"
echo "  Local:     $0 local alice proposal.json"
echo "  Staging:   $0 staging alice proposal.json"
echo "  Testnet:   $0 testnet dydx-1 proposal.json"
echo "============================================================"
