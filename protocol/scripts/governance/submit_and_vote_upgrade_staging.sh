#!/bin/bash
set -eo pipefail

# Submits a MsgSoftwareUpgrade proposal (via `tx upgrade software-upgrade`) and votes yes from
# every validator key staging.sh provisions. Mirrors vote_in_staging.sh's conventions: same binary
# path, fees, and reliance on the local client already being configured to reach staging with
# these keys imported.
#
# Usage: ./scripts/governance/submit_and_vote_upgrade_staging.sh <upgrade-name> [height-buffer-blocks]
#
# height-buffer-blocks defaults to 300, i.e. ~5 minutes assuming ~1s blocks -- staging's
# voting_period is 120s (see testing/genesis.sh), so this leaves room for the deposit/vote txs
# and tally to land before the upgrade height arrives. Adjust if staging's block time differs.

UPGRADE_NAME="${1:?Usage: $0 <upgrade-name> [height-buffer-blocks]}"
HEIGHT_BUFFER="${2:-300}"

FEES="5000000000000000adv4tnt"
DEPOSIT="10000000000000000000000adv4tnt" # 10,000 native tokens, matches enable_all_clob_pairs.sh

CURRENT_HEIGHT=$(./build/dydxprotocold status | jq -r '.sync_info.latest_block_height')
UPGRADE_HEIGHT=$((CURRENT_HEIGHT + HEIGHT_BUFFER))

echo "Current height: $CURRENT_HEIGHT. Scheduling '$UPGRADE_NAME' upgrade at height $UPGRADE_HEIGHT."

./build/dydxprotocold tx upgrade software-upgrade "$UPGRADE_NAME" \
	--upgrade-height "$UPGRADE_HEIGHT" \
	--no-validate \
	--title "Upgrade to $UPGRADE_NAME" \
	--summary "Upgrade to $UPGRADE_NAME" \
	--deposit "$DEPOSIT" \
	--from alice \
	--broadcast-mode sync \
	--fees "$FEES" \
	--yes

sleep 6 # let the submit-proposal tx land before querying for its proposal id.

PROPOSAL_ID=$(./build/dydxprotocold query gov proposals --output json | jq -r '.proposals[-1].id')
echo "Voting yes on proposal $PROPOSAL_ID from all staging keys."

for key in alice bob carl dave emily fiona greg henry ian jeff; do
	./build/dydxprotocold tx gov vote "$PROPOSAL_ID" yes --from "$key" --broadcast-mode sync --fees "$FEES" --yes
done

echo "Voted yes from all keys on proposal $PROPOSAL_ID. Upgrade '$UPGRADE_NAME' is scheduled for" \
	"height $UPGRADE_HEIGHT once the vote tallies -- watch cosmovisor logs around then."
