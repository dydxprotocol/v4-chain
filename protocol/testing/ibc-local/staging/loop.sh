#!/bin/bash
# Injective-leg loop on the live testnet through the injective-1 stand-in
# (standin.sh must be up). Noble withdrawals are never sent: the testnet
# Noble clients are expired on both sides.
#   loop.sh state          print ledger, balances and rate-limit capacity
#   loop.sh round [uusdc]  deposit from the stand-in, then withdraw the same amount
#   loop.sh drain          withdraw the whole injective_backing
#   loop.sh reject         withdraw 1 USDC with no backing; expects "insufficient route backing"
set -u
HARNESS="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
STAGING_DIR="${STAGING_DIR:-$HOME/.ibc-staging}"
BIN="$HARNESS/../../build/ibc-local/bin"
DYDX="$BIN/dydxprotocold"; SIMD="$BIN/ibc-simd"
NODE="${REMOTE_RPC:-tcp://3.22.142.107:26657}"; INJ=tcp://localhost:29657
INJ_HOME="$STAGING_DIR/injective"; DYDX_KEYS="$STAGING_DIR/dydx-keys"
. "$STAGING_DIR/channel.env"
INJECTIVE_USDC=erc20:0xa00C59fF5a080D2b954d0c75e46E22a0c371235a
LOGICAL=ibc/8E27BA2D5493AF5636760E354E46004562C46AB7EC0CC4C1CA14E9E20E2545B5
INJ_USER=$("$SIMD" keys show user -a --keyring-backend test --home "$INJ_HOME")
USER=$("$DYDX" keys show relayer -a --keyring-backend test --home "$DYDX_KEYS")
MODULE=$("$DYDX" q auth module-account canonicalusdc --node "$NODE" -o json | jq -r '.account.value.address')
ESCROW=$("$SIMD" q ibc-transfer escrow-address transfer "$INJECTIVE_SIDE_CHANNEL" --node "$INJ" | tr -d '[:space:]')
VOUCHER=$("$DYDX" q canonicalusdc state --node "$NODE" -o json | jq -r '.controls.injective_denom')

dydx_bal() { "$DYDX" q bank balance "$1" "$2" --node "$NODE" -o json 2>/dev/null | jq -r '.balance.amount // "0"'; }
inj_bal() { "$SIMD" q bank balance "$1" "$2" --node "$INJ" -o json 2>/dev/null | jq -r '.balance.amount // "0"'; }
ledger() { "$DYDX" q canonicalusdc state --node "$NODE" -o json | jq -r '.ledger | "nb=\(.noble_backing) ib=\(.injective_backing) pi=\(.pending_injective) pn=\(.pending_noble) rc=\(.refillable_capacity)"'; }
pending() { "$DYDX" q canonicalusdc state --node "$NODE" -o json | jq -r '.ledger.pending_injective'; }
backing() { "$DYDX" q canonicalusdc state --node "$NODE" -o json | jq -r '.ledger.injective_backing'; }
state() {
	echo "[$(date -u +%H:%M:%S)] $1"
	echo "   ledger: $(ledger)  supply=$("$DYDX" q bank total-supply-of "$LOGICAL" --node "$NODE" -o json | jq -r '.amount.amount')"
	echo "   user logical=$(dydx_bal "$USER" "$LOGICAL")  module voucher=$(dydx_bal "$MODULE" "$VOUCHER")  stand-in user=$(inj_bal "$INJ_USER" "$INJECTIVE_USDC")  stand-in escrow=$(inj_bal "$ESCROW" "$INJECTIVE_USDC")"
	echo "   capacity: $("$DYDX" q ratelimit capacity-by-denom "$LOGICAL" --node "$NODE" -o json | jq -c '[.limiter_capacity_list[].capacity]')"
}
wait_for() { # wait_for <desc> <want> <cmd...>
	local desc=$1 want=$2 i v; shift 2
	for i in $(seq 1 90); do v=$("$@"); [ "$v" = "$want" ] && { echo "   ok: $desc = $v after $((i * 2))s"; return 0; }; sleep 2; done
	echo "   TIMEOUT: $desc = $v (wanted $want)"; return 1
}
tx_result() { # tx_result <bin> <node> <hash>
	local i r; for i in $(seq 1 30); do r=$("$1" q tx "$3" --node "$2" -o json 2>/dev/null) && { echo "$r" | jq -r '"code=\(.code) \(.raw_log[0:120])"'; return; }; sleep 2; done; echo "not found"
}
deposit() {
	local amt=$1 before h; before=$(dydx_bal "$USER" "$LOGICAL")
	echo "== deposit $amt $INJECTIVE_USDC stand-in -> dydx"
	h=$("$SIMD" tx ibc-transfer transfer transfer "$INJECTIVE_SIDE_CHANNEL" "$USER" "$amt$INJECTIVE_USDC" --from user --keyring-backend test --home "$INJ_HOME" \
		--chain-id injective-1 --node "$INJ" --packet-timeout-height 0-0 --gas 300000 --fees 0stake -y -o json | jq -r .txhash)
	echo "   stand-in tx $h: $(tx_result "$SIMD" "$INJ" "$h")"
	wait_for "user logical" $((before + amt)) dydx_bal "$USER" "$LOGICAL"
	state "after deposit $amt"
}
withdraw() {
	local amt=$1 before h; before=$(inj_bal "$INJ_USER" "$INJECTIVE_USDC")
	echo "== withdraw $amt logical dydx -> stand-in"
	h=$("$DYDX" tx ibc-transfer transfer transfer "$DYDX_SIDE_CHANNEL" "$INJ_USER" "$amt$LOGICAL" --from relayer --keyring-backend test --home "$DYDX_KEYS" \
		--chain-id dydx-testnet-4 --node "$NODE" --gas 500000 --gas-prices 25000000000adv4tnt -y -o json | jq -r .txhash)
	echo "   dydx tx $h: $(tx_result "$DYDX" "$NODE" "$h")"
	echo "   ledger after send: $(ledger)"
	wait_for "stand-in user" $((before + amt)) inj_bal "$INJ_USER" "$INJECTIVE_USDC"
	wait_for "pending_injective" 0 pending
	state "after withdraw $amt"
}

state start
case "${1:-state}" in
	state) ;;
	round) deposit "${2:-25000000}"; withdraw "${2:-25000000}" ;;
	drain) withdraw "$(backing)" ;;
	reject)
		h=$("$DYDX" tx ibc-transfer transfer transfer "$DYDX_SIDE_CHANNEL" "$INJ_USER" "1000000$LOGICAL" --from relayer --keyring-backend test --home "$DYDX_KEYS" \
			--chain-id dydx-testnet-4 --node "$NODE" --gas 500000 --gas-prices 25000000000adv4tnt -y -o json | jq -r .txhash)
		echo "== withdraw 1 USDC with no backing, dydx tx $h: $(tx_result "$DYDX" "$NODE" "$h")"
		state "after the rejected withdrawal" ;;
	*) echo "usage: loop.sh state|round [uusdc]|drain|reject" >&2; exit 2 ;;
esac
