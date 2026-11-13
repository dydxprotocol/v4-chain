#!/bin/bash
# Governance helpers for x/canonicalusdc on a chain whose validator keys are local
# (testnet: dydx-1 and dydx-2 in ~/.dydxprotocol, keyring test, expedited = 10 min at 75%).
#   gov.sh controls <out.json> <title> <jq filter>   build MsgUpdateControls from the LIVE controls
#   gov.sh submit <proposal.json>                     submit, vote yes from both validators, wait
# The controls message is a full replacement, so always start from the live set, e.g.
#   gov.sh controls pause.json "canonicalusdc: PAUSED" '.mode = "MODE_PAUSED"'
#   gov.sh controls nodep.json "no Injective deposits" '.injective_deposits_disabled = true'
set -u
BIN="$(cd "$(dirname "${BASH_SOURCE[0]}")/../../../build/ibc-local/bin" && pwd)"
DYDX="$BIN/dydxprotocold"
NODE="${REMOTE_RPC:-tcp://3.22.142.107:26657}"; CHAIN_ID="${REMOTE_CHAIN_ID:-dydx-testnet-4}"
FUNDER_HOME="${FUNDER_HOME:-$HOME/.dydxprotocol}"; VOTERS="${VOTERS:-dydx-1 dydx-2}"
GOV=$("$DYDX" q auth module-account gov --node "$NODE" -o json | jq -r '.account.value.address')
TX="--home $FUNDER_HOME --keyring-backend test --chain-id $CHAIN_ID --node $NODE --gas-prices 25000000000adv4tnt -y -o json"

controls() {
	local out=$1 title=$2 filter=$3
	jq -n --arg gov "$GOV" --arg title "$title" \
		--argjson controls "$("$DYDX" q canonicalusdc state --node "$NODE" -o json | jq ".controls | $filter")" \
		'{messages: [{"@type": "/dydxprotocol.canonicalusdc.MsgUpdateControls", authority: $gov, controls: $controls}],
		  deposit: "10000000000000000000adv4tnt", title: $title, summary: $title, metadata: "", expedited: true}' >"$out"
	jq -c '.messages[0].controls | {mode, injective_deposits_disabled, noble_deposits_disabled, noble_withdrawals_enabled}' "$out"
}

submit() {
	local json=$1 h id st v i
	h=$("$DYDX" tx gov submit-proposal "$json" --from "${VOTERS%% *}" --gas auto --gas-adjustment 1.5 $TX | jq -r .txhash)
	for i in $(seq 1 20); do
		id=$("$DYDX" q tx "$h" --node "$NODE" -o json 2>/dev/null | jq -r '[.events[] | select(.type=="submit_proposal") | .attributes[] | select(.key=="proposal_id") | .value][0] // empty')
		[ -n "$id" ] && break; sleep 2
	done
	[ -n "$id" ] || { echo "no proposal id in tx $h"; return 1; }
	echo "proposal $id submitted: $(jq -r .title "$json")"
	for v in $VOTERS; do "$DYDX" tx gov vote "$id" yes --from "$v" --gas 300000 $TX | jq -r '"  vote '"$v"': code \(.code)"'; sleep 3; done
	# the CLI prints the status as a number: 2 voting, 3 passed, 4 rejected, 5 failed
	for i in $(seq 1 100); do
		st=$("$DYDX" q gov proposal "$id" --node "$NODE" -o json | jq -r .proposal.status)
		case "$st" in 3|4|5|PROPOSAL_STATUS_PASSED|PROPOSAL_STATUS_REJECTED|PROPOSAL_STATUS_FAILED) break;; esac; sleep 10
	done
	echo "proposal $id status $st"; [ "$st" = 3 ] || [ "$st" = PROPOSAL_STATUS_PASSED ]
}

case "${1:-}" in
	controls) controls "$2" "$3" "$4" ;;
	submit) submit "$2" ;;
	*) echo "usage: gov.sh controls <out.json> <title> <jq filter> | submit <proposal.json>" >&2; exit 2 ;;
esac
