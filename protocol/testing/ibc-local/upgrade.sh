#!/bin/bash
# Starts the environment on an earlier dydxprotocold release with USDC already
# in circulation, upgrades it in place to the working tree's binary through
# x/upgrade, checks that the canonicalusdc store appears DISABLED and empty,
# then runs every test.sh scenario on the upgraded chain.
#
#   upgrade.sh [FROM_REF]    default: protocol/v9.8.2
set -euo pipefail
source "$(dirname "${BASH_SOURCE[0]}")/lib.sh"

FROM_REF="${1:-protocol/v9.8.2}"
UPGRADE_NAME="v9.9"
OLD_BIN="$BIN_DIR/dydxprotocold-$(echo "$FROM_REF" | tr '/' '-')"
OLD_SRC="$STATE_DIR/src/$(echo "$FROM_REF" | tr '/' '-')"
NEW_BIN="$BIN_DIR/dydxprotocold"

require_cmd git go jq curl bc shasum
"$IBC_LOCAL_DIR/setup.sh"
if [ ! -x "$OLD_BIN" ]; then
	log "building dydxprotocold at $FROM_REF"
	[ -d "$OLD_SRC" ] || git -C "$PROTOCOL_DIR" worktree add --detach "$OLD_SRC" "$FROM_REF^{commit}" >/dev/null
	(cd "$OLD_SRC/protocol" && go build -o "$OLD_BIN" ./cmd/dydxprotocold)
fi

IBC_LOCAL_DYDX_BIN="$OLD_BIN" "$IBC_LOCAL_DIR/up.sh"
DYDX_BIN="$OLD_BIN"
log "running $("$DYDX_BIN" version 2>&1 | tail -1) from $FROM_REF"

DYDX_USER=$(addr dydx user)
NOBLE_USER=$(addr noble user)
NOBLE_CHANNEL=$(q dydx ibc channel end transfer "$DYDX_NOBLE_CHANNEL" | jq -r '.channel.counterparty.channel_id')
GOV_MODULE=$(module_address gov)
WORK_DIR="$STATE_DIR/test"
mkdir -p "$WORK_DIR"

log "pre-upgrade: Noble deposit and withdrawal on the old binary"
noble_before=$(balance noble "$NOBLE_USER" "$COUNTERPARTY_USDC")
tx noble user ibc-transfer transfer transfer "$NOBLE_CHANNEL" "$DYDX_USER" "1000000000$COUNTERPARTY_USDC" >/dev/null
wait_for_balance dydx "$DYDX_USER" "$USDC_DENOM" 1000000000
tx dydx user ibc-transfer transfer transfer "$DYDX_NOBLE_CHANNEL" "$NOBLE_USER" "100000000$USDC_DENOM" >/dev/null
wait_for_balance noble "$NOBLE_USER" "$COUNTERPARTY_USDC" "$(bc_calc "$noble_before - 900000000")"
user_before=$(balance dydx "$DYDX_USER" "$USDC_DENOM")
supply_before=$(supply_of dydx "$USDC_DENOM")
assert_eq "USDC in circulation before the upgrade" "$supply_before" 900000000
pass "USDC circulates on $FROM_REF without the canonicalusdc module"

# 120 blocks at 500ms leaves the 20s voting period plus margin.
target=$(($(height dydx) + 120))
jq -n --arg gov "$GOV_MODULE" --arg name "$UPGRADE_NAME" --arg height "$target" --arg dep "10000000$NATIVE_TOKEN" '{
	messages: [{
		"@type": "/cosmos.upgrade.v1beta1.MsgSoftwareUpgrade",
		authority: $gov,
		plan: { name: $name, time: "0001-01-01T00:00:00Z", height: $height, info: "", upgraded_client_state: null }
	}],
	deposit: $dep, title: ("upgrade " + $name), summary: "ibc-local upgrade", expedited: false
}' >"$WORK_DIR/upgrade-proposal.json"
govern "$WORK_DIR/upgrade-proposal.json" PASSED "upgrade $UPGRADE_NAME" >/dev/null
[ "$(height dydx)" -lt "$target" ] || die "upgrade height $target already passed"
log "upgrade $UPGRADE_NAME scheduled at height $target"

stop_hermes
for i in $(seq 1 120); do
	if ! pid_alive dydx || grep -q "UPGRADE \"$UPGRADE_NAME\" NEEDED" "$LOG_DIR/dydx.log"; then
		break
	fi
	sleep 1
done
grep -q "UPGRADE \"$UPGRADE_NAME\" NEEDED" "$LOG_DIR/dydx.log" || die "old binary did not halt at height $target"
stop_process dydx
assert_eq "upgrade-info name" "$(jq -r '.name' "$DYDX_HOME/data/upgrade-info.json")" "$UPGRADE_NAME"
assert_eq "upgrade-info height" "$(jq -r '.height' "$DYDX_HOME/data/upgrade-info.json")" "$target"
pass "old binary halted at the upgrade height"

mv "$LOG_DIR/dydx.log" "$LOG_DIR/dydx-$(echo "$FROM_REF" | tr '/' '-').log"
DYDX_BIN="$NEW_BIN"
start_dydx
wait_for_height dydx $((target + 5))
grep -q "applying upgrade \"$UPGRADE_NAME\"" "$LOG_DIR/dydx.log" || die "new binary did not apply $UPGRADE_NAME"
start_hermes

state=$(q dydx canonicalusdc state)
assert_eq "mode after upgrade" "$(echo "$state" | jq -r '.controls.mode')" MODE_DISABLED
assert_eq "logical denom after upgrade" "$(echo "$state" | jq -r '.controls.logical_denom')" "$USDC_DENOM"
assert_eq "ledger after upgrade" "$(echo "$state" | jq -c '.ledger | map_values(tonumber) | add')" 0
assert_eq "pending after upgrade" "$(q dydx canonicalusdc pending-settlements | jq -r '.pending_settlements | length')" 0
assert_eq "user USDC after upgrade" "$(balance dydx "$DYDX_USER" "$USDC_DENOM")" "$user_before"
assert_eq "USDC supply after upgrade" "$(supply_of dydx "$USDC_DENOM")" "$supply_before"
pass "canonicalusdc store added DISABLED and empty, balances intact"

"$IBC_LOCAL_DIR/test.sh"
log "upgrade run from $FROM_REF passed"
