#!/bin/bash
# Starts the environment on an earlier dydxprotocold release with USDC in
# circulation and an open perpetual position, upgrades it in place to the
# working tree's binary through x/upgrade, checks that the canonicalusdc store
# appears DISABLED and empty, then runs every test.sh scenario on the upgraded
# chain. Subaccount collateral, positions, a resting long-term order and
# MegaVault shares must be untouched by the upgrade and by every canonical USDC
# flow, including the migration loop; afterwards the resting order must still
# fill and the MegaVault depositor must get their USDC back.
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
DYDX_HOLDER=$(addr dydx holder)
NOBLE_USER=$(addr noble user)
NOBLE_CHANNEL=$(q dydx ibc channel end transfer "$DYDX_NOBLE_CHANNEL" | jq -r '.channel.counterparty.channel_id')
GOV_MODULE=$(module_address gov)
WORK_DIR="$STATE_DIR/test"
mkdir -p "$WORK_DIR"

# Submits a tx without waiting for inclusion: short-term orders never appear
# in a block, only their fills do.
broadcast() {
	local chain=$1 from=$2
	shift 2
	local res code
	res=$("$(chain_bin "$chain")" tx "$@" --from "$from" --keyring-backend test --home "$(chain_home "$chain")" \
		--chain-id "$(chain_id "$chain")" --node "$(chain_node "$chain")" --gas 0 --fees "0$NATIVE_TOKEN" \
		--yes --output json) || die "broadcast failed on $chain: $*"
	code=$(echo "$res" | jq -r '.code')
	[ "$code" = "0" ] || die "tx rejected on $chain (code $code): $(echo "$res" | jq -r '.raw_log')"
}

# Every subaccount's collateral and perpetual positions plus each perpetual's
# open interest: the state the upgrade and the canonical USDC flows must not
# touch.
positions_snapshot() {
	q dydx subaccounts list-subaccount | jq -S -c '[.subaccount[] |
		select((.asset_positions | length) + (.perpetual_positions | length) > 0) |
		{id, asset_positions: [.asset_positions[] | {asset_id, quantums}],
		 perpetual_positions: [.perpetual_positions[] | {perpetual_id, quantums}]}]'
	q dydx perpetuals list-perpetual | jq -S -c '[.perpetual[] | {id: .params.id, open_interest}]'
}

# Funds two subaccounts with USDC and fills 0.01 BTC-USD between them.
open_position() {
	local subticks gtb i
	tx dydx user bank send "$DYDX_USER" "$DYDX_HOLDER" "300000000$USDC_DENOM" >/dev/null
	tx dydx user sending deposit-to-subaccount "$DYDX_USER" "$DYDX_USER" 0 300000000 >/dev/null
	tx dydx holder sending deposit-to-subaccount "$DYDX_HOLDER" "$DYDX_HOLDER" 0 300000000 >/dev/null
	# BTC-USD prices have exponent -5 and 1e5 subticks per dollar, so the
	# oracle price rounded down to a tick is the order price.
	subticks=$((price / 100000 * 100000))
	gtb=$(($(height dydx) + 15))
	broadcast dydx user clob place-order "$DYDX_USER" 0 1 0 1 100000000 "$subticks" "$gtb"
	broadcast dydx holder clob place-order "$DYDX_HOLDER" 0 1 0 2 100000000 "$subticks" "$gtb"
	for i in $(seq 1 30); do
		[ "$(q dydx subaccounts show-subaccount "$DYDX_USER" 0 |
			jq -r '.subaccount.perpetual_positions[0].quantums // "0"')" = "100000000" ] && return 0
		sleep 1
	done
	die "BTC-USD position did not fill"
}

# Places a long-term bid from user at the given subticks. The short-term CLI
# command generates the tx, which is patched to a long-term order and signed.
place_long_term_bid() {
	local client_id=$1 subticks=$2 unsigned res
	unsigned="$WORK_DIR/long-term-$client_id.json"
	"$DYDX_BIN" tx clob place-order "$DYDX_USER" 0 "$client_id" 0 1 100000000 "$subticks" 1 \
		--from user --keyring-backend test --home "$DYDX_HOME" --chain-id "$DYDX_CHAIN_ID" \
		--node "$(chain_node dydx)" --gas 0 --fees "0$NATIVE_TOKEN" --generate-only 2>/dev/null |
		jq --argjson t "$(($(date +%s) + 86400))" \
			'.body.messages[0].order |= (del(.good_til_block) | .order_id.order_flags = 64 | .good_til_block_time = $t)' \
			>"$unsigned"
	"$DYDX_BIN" tx sign "$unsigned" --from user --keyring-backend test --home "$DYDX_HOME" \
		--chain-id "$DYDX_CHAIN_ID" --node "$(chain_node dydx)" --output-document "$unsigned.signed" 2>/dev/null
	res=$("$DYDX_BIN" tx broadcast "$unsigned.signed" --node "$(chain_node dydx)" --output json)
	[ "$(echo "$res" | jq -r '.code')" = "0" ] || die "long-term order rejected: $(echo "$res" | jq -r '.raw_log')"
	wait_for_tx dydx "$(echo "$res" | jq -r '.txhash')" >/dev/null
}

resting_order() { q dydx clob stateful-order "$DYDX_USER" 0 2 0 64 | jq -c '.'; }

vault_snapshot() {
	q dydx vault owner-shares "$DYDX_USER" | jq -c '.'
	q dydx vault total-shares | jq -c '.'
}

subaccount_quote() {
	q dydx subaccounts show-subaccount "$1" 0 |
		jq -r '[.subaccount.asset_positions[] | select(.asset_id == 0) | .quantums][0] // "0"'
}

log "pre-upgrade: Noble deposit and withdrawal on the old binary"
noble_before=$(balance noble "$NOBLE_USER" "$COUNTERPARTY_USDC")
tx noble user ibc-transfer transfer transfer "$NOBLE_CHANNEL" "$DYDX_USER" "1000000000$COUNTERPARTY_USDC" >/dev/null
wait_for_balance dydx "$DYDX_USER" "$USDC_DENOM" 1000000000
tx dydx user ibc-transfer transfer transfer "$DYDX_NOBLE_CHANNEL" "$NOBLE_USER" "100000000$USDC_DENOM" >/dev/null
wait_for_balance noble "$NOBLE_USER" "$COUNTERPARTY_USDC" "$(bc_calc "$noble_before - 900000000")"
pass "USDC circulates on $FROM_REF without the canonicalusdc module"

log "pre-upgrade: open a BTC-USD position collateralised with that USDC"
price=$(q dydx prices list-market-price | jq -r '.market_prices[] | select(.id == 0) | .price')
open_position
positions_before=$(positions_snapshot)
assert_eq "long and short filled" \
	"$(echo "$positions_before" | head -1 | jq -c '[.[] | .perpetual_positions[]?.quantums] | sort')" \
	'["-100000000","100000000"]'
user_before=$(balance dydx "$DYDX_USER" "$USDC_DENOM")
supply_before=$(supply_of dydx "$USDC_DENOM")
assert_eq "USDC in circulation before the upgrade" "$supply_before" 900000000
pass "0.01 BTC-USD open between two subaccounts before the upgrade"

log "pre-upgrade: rest a long-term bid at half the oracle price"
resting_subticks=$((price / 2 / 100000 * 100000))
place_long_term_bid 2 "$resting_subticks"
resting_before=$(resting_order)
assert_eq "long-term bid in state" "$(echo "$resting_before" | jq -r '.. | .order_flags? // empty' | head -1)" 64
pass "long-term bid resting in state before the upgrade"

log "pre-upgrade: deposit 100 USDC into MegaVault"
tx dydx user vault deposit-to-megavault "$DYDX_USER" 0 100000000 >/dev/null
vault_before=$(vault_snapshot)
assert_eq "MegaVault shares" "$(q dydx vault total-shares | jq -r '.total_shares.num_shares')" 100000000
positions_before=$(positions_snapshot)
pass "100 USDC deposited into MegaVault before the upgrade"

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
assert_eq "positions after upgrade" "$(positions_snapshot)" "$positions_before"
assert_eq "MegaVault shares after upgrade" "$(vault_snapshot)" "$vault_before"
assert_eq "resting order after upgrade" "$(resting_order)" "$resting_before"
pass "subaccount collateral, positions, resting order and MegaVault shares unchanged by the upgrade"

"$IBC_LOCAL_DIR/test.sh"
assert_eq "positions after every scenario" "$(positions_snapshot)" "$positions_before"
assert_eq "MegaVault shares after every scenario" "$(vault_snapshot)" "$vault_before"
assert_eq "resting order after every scenario" "$(resting_order)" "$resting_before"
pass "subaccount collateral, positions, resting order and MegaVault shares unchanged by every canonical USDC flow"

gtb=$(($(height dydx) + 15))
broadcast dydx holder clob place-order "$DYDX_HOLDER" 0 3 0 2 100000000 "$resting_subticks" "$gtb"
for i in $(seq 1 30); do
	[ "$(q dydx subaccounts show-subaccount "$DYDX_USER" 0 |
		jq -r '.subaccount.perpetual_positions[0].quantums // "0"')" = "200000000" ] && break
	sleep 1
done
assert_eq "resting bid filled" "$(q dydx subaccounts show-subaccount "$DYDX_USER" 0 |
	jq -r '.subaccount.perpetual_positions[0].quantums // "0"')" 200000000
pass "the long-term bid placed before the upgrade still fills after the migration"

quote_before=$(subaccount_quote "$DYDX_USER")
tx dydx user vault withdraw-from-megavault "$DYDX_USER" 0 100000000 100000000 >/dev/null
assert_eq "MegaVault withdrawal returns the deposit" "$(subaccount_quote "$DYDX_USER")" "$(bc_calc "$quote_before + 100000000")"
assert_eq "MegaVault shares redeemed" "$(q dydx vault total-shares | jq -r '.total_shares.num_shares')" 0
pass "MegaVault depositor withdraws the full 100 USDC after the migration"
log "upgrade run from $FROM_REF passed"
