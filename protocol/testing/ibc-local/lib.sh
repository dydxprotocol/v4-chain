#!/bin/bash
# Shared configuration and helpers for the local three-chain IBC environment.
# Sourced by setup.sh, up.sh, test.sh and down.sh; not meant to be run directly.

IBC_LOCAL_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROTOCOL_DIR="$(cd "$IBC_LOCAL_DIR/../.." && pwd)"
STATE_DIR="${IBC_LOCAL_STATE_DIR:-$PROTOCOL_DIR/build/ibc-local}"
BIN_DIR="$STATE_DIR/bin"
LOG_DIR="$STATE_DIR/logs"
PID_DIR="$STATE_DIR/pids"
KEY_DIR="$STATE_DIR/keys"
CHANNELS_ENV="$STATE_DIR/channels.env"

HERMES_VERSION="v1.13.3"
DASEL_VERSION="v2.8.1"
# Matches the upstream ibc-go version in go.mod; the counterparties run its simd.
IBC_GO_VERSION="v8.5.1"

DYDX_BIN="$BIN_DIR/dydxprotocold"
SIMD_BIN="$BIN_DIR/ibc-simd"
HERMES_BIN="$BIN_DIR/hermes"
DASEL_BIN="$BIN_DIR/dasel"
HERMES_HOME="$STATE_DIR/hermes"
HERMES_CONFIG="$HERMES_HOME/config.toml"

# The logical USDC denom is the hash of transfer/channel-0/uusdc, so the Noble
# stand-in must own channel-0 on dYdX. up.sh opens it first and asserts this.
USDC_DENOM="ibc/8E27BA2D5493AF5636760E354E46004562C46AB7EC0CC4C1CA14E9E20E2545B5"
NATIVE_TOKEN="adv4tnt"
COUNTERPARTY_USDC="uusdc"
DYDX_NOBLE_CHANNEL="channel-0"
DYDX_INJECTIVE_CHANNEL="channel-1"

DYDX_CHAIN_ID="localdydxprotocol"
DYDX_HOME="$STATE_DIR/dydx"
DYDX_RPC_PORT=26657
DYDX_P2P_PORT=26656
DYDX_GRPC_PORT=9090
DYDX_API_PORT=1317
DYDX_PPROF_PORT=6060
DYDX_FEES="10000000000000000$NATIVE_TOKEN"
DYDX_GAS=400000

NOBLE_CHAIN_ID="noble-local"
NOBLE_HOME="$STATE_DIR/noble"
NOBLE_RPC_PORT=27657
NOBLE_P2P_PORT=27656
NOBLE_GRPC_PORT=9190
NOBLE_API_PORT=1417
NOBLE_PPROF_PORT=6160

INJECTIVE_CHAIN_ID="injective-local"
INJECTIVE_HOME="$STATE_DIR/injective"
INJECTIVE_RPC_PORT=28657
INJECTIVE_P2P_PORT=28656
INJECTIVE_GRPC_PORT=9290
INJECTIVE_API_PORT=1517
INJECTIVE_PPROF_PORT=6260

# Seconds to wait for a relayed packet (recv + ack) before failing.
RELAY_TIMEOUT_SECONDS="${RELAY_TIMEOUT_SECONDS:-90}"

log() { printf '\033[1;34m[ibc-local]\033[0m %s\n' "$*" >&2; }
pass() { printf '\033[1;32m[PASS]\033[0m %s\n' "$*" >&2; }
die() {
	printf '\033[1;31m[FAIL]\033[0m %s\n' "$*" >&2
	exit 1
}

require_cmd() {
	local cmd
	for cmd in "$@"; do
		command -v "$cmd" >/dev/null 2>&1 || die "missing required command: $cmd"
	done
}

# dydx/noble/injective are the chain handles used by every helper below.
chain_bin() {
	case "$1" in
	dydx) echo "$DYDX_BIN" ;;
	noble | injective) echo "$SIMD_BIN" ;;
	*) die "unknown chain: $1" ;;
	esac
}

chain_home() {
	case "$1" in
	dydx) echo "$DYDX_HOME" ;;
	noble) echo "$NOBLE_HOME" ;;
	injective) echo "$INJECTIVE_HOME" ;;
	*) die "unknown chain: $1" ;;
	esac
}

chain_id() {
	case "$1" in
	dydx) echo "$DYDX_CHAIN_ID" ;;
	noble) echo "$NOBLE_CHAIN_ID" ;;
	injective) echo "$INJECTIVE_CHAIN_ID" ;;
	*) die "unknown chain: $1" ;;
	esac
}

chain_rpc_port() {
	case "$1" in
	dydx) echo "$DYDX_RPC_PORT" ;;
	noble) echo "$NOBLE_RPC_PORT" ;;
	injective) echo "$INJECTIVE_RPC_PORT" ;;
	*) die "unknown chain: $1" ;;
	esac
}

chain_node() { echo "tcp://127.0.0.1:$(chain_rpc_port "$1")"; }

# Runs a query against a chain and prints JSON.
q() {
	local chain=$1
	shift
	"$(chain_bin "$chain")" query "$@" --node "$(chain_node "$chain")" --output json
}

# Prints the bech32 address of a key in a chain's test keyring.
addr() {
	"$(chain_bin "$1")" keys show "$2" -a --keyring-backend test --home "$(chain_home "$1")"
}

balance() {
	local out
	out=$(q "$1" bank balances "$2") || die "balance query failed on $1 for $2"
	echo "$out" | jq -r --arg d "$3" '[.balances[] | select(.denom == $d) | .amount][0] // "0"'
}

supply_of() {
	local out
	out=$(q "$1" bank total --denom "$2" 2>/dev/null) || out=$(q "$1" bank total-supply-of "$2")
	echo "$out" | jq -r '.amount.amount // .amount // "0"'
}

# Signs and broadcasts a tx, waits for inclusion and fails on a non-zero code.
# Prints the committed tx JSON.
tx() {
	local chain=$1 from=$2
	shift 2
	local bin home fees res hash code
	bin=$(chain_bin "$chain")
	home=$(chain_home "$chain")
	fees="0stake"
	if [ "$chain" = dydx ]; then
		fees="$DYDX_FEES"
	fi
	res=$("$bin" tx "$@" --from "$from" --keyring-backend test --home "$home" \
		--chain-id "$(chain_id "$chain")" --node "$(chain_node "$chain")" \
		--gas "$DYDX_GAS" --fees "$fees" --yes --output json) || die "tx failed on $chain: $*"
	code=$(echo "$res" | jq -r '.code')
	[ "$code" = "0" ] || die "tx rejected on $chain (code $code): $(echo "$res" | jq -r '.raw_log')"
	hash=$(echo "$res" | jq -r '.txhash')
	wait_for_tx "$chain" "$hash"
}

wait_for_tx() {
	local chain=$1 hash=$2 i out code
	for i in $(seq 1 30); do
		if out=$(q "$chain" tx "$hash" 2>/dev/null); then
			code=$(echo "$out" | jq -r '.code')
			[ "$code" = "0" ] || die "tx $hash failed on $chain (code $code): $(echo "$out" | jq -r '.raw_log')"
			echo "$out"
			return 0
		fi
		sleep 1
	done
	die "tx $hash not committed on $chain after 30s"
}

height() {
	curl -sf "http://127.0.0.1:$(chain_rpc_port "$1")/status" | jq -r '.result.sync_info.latest_block_height'
}

wait_for_height() {
	local chain=$1 target=$2 i h
	for i in $(seq 1 120); do
		h=$(height "$chain" 2>/dev/null || echo 0)
		if [ -n "$h" ] && [ "$h" != "null" ] && [ "$h" -ge "$target" ]; then
			return 0
		fi
		sleep 1
	done
	die "$chain did not reach height $target within 120s (see $LOG_DIR/$chain.log)"
}

# Polls until `balance chain addr denom` equals want.
wait_for_balance() {
	local chain=$1 address=$2 denom=$3 want=$4 i got
	for i in $(seq 1 "$RELAY_TIMEOUT_SECONDS"); do
		got=$(balance "$chain" "$address" "$denom")
		if [ "$got" = "$want" ]; then
			return 0
		fi
		sleep 1
	done
	die "$chain balance of $denom for $address is $got, want $want after ${RELAY_TIMEOUT_SECONDS}s"
}

assert_eq() {
	[ "$2" = "$3" ] || die "$1: got '$2', want '$3'"
}

# Integer arithmetic on amounts that overflow 64-bit shell math.
bc_calc() { echo "$1" | bc; }

ibc_denom() {
	# ICS20 voucher denom: ibc/ + upper-case hex SHA256 of the trace path.
	printf '%s' "$1" | shasum -a 256 | awk '{print "ibc/" toupper($1)}'
}

pid_alive() {
	local pidfile="$PID_DIR/$1.pid"
	[ -f "$pidfile" ] && kill -0 "$(cat "$pidfile")" 2>/dev/null
}

start_process() {
	local name=$1
	shift
	mkdir -p "$LOG_DIR" "$PID_DIR"
	if pid_alive "$name"; then
		log "$name already running (pid $(cat "$PID_DIR/$name.pid"))"
		return 0
	fi
	nohup "$@" >"$LOG_DIR/$name.log" 2>&1 &
	echo $! >"$PID_DIR/$name.pid"
	log "started $name (pid $!), log: $LOG_DIR/$name.log"
}

stop_process() {
	local name=$1 pidfile="$PID_DIR/$1.pid" pid i
	[ -f "$pidfile" ] || return 0
	pid=$(cat "$pidfile")
	if kill -0 "$pid" 2>/dev/null; then
		kill "$pid"
		for i in $(seq 1 20); do
			kill -0 "$pid" 2>/dev/null || break
			sleep 0.5
		done
		kill -0 "$pid" 2>/dev/null && kill -9 "$pid"
		log "stopped $name (pid $pid)"
	fi
	command rm -f "$pidfile"
}

start_hermes() {
	start_process hermes "$HERMES_BIN" --config "$HERMES_CONFIG" start
}

stop_hermes() { stop_process hermes; }

# Like tx, but the tx must fail, either at submission or in the block.
# Prints the failure log so callers can check the reason.
tx_expect_fail() {
	local chain=$1 from=$2
	shift 2
	local fees="0stake" res code hash i out
	if [ "$chain" = dydx ]; then
		fees="$DYDX_FEES"
	fi
	if ! res=$("$(chain_bin "$chain")" tx "$@" --from "$from" --keyring-backend test \
		--home "$(chain_home "$chain")" --chain-id "$(chain_id "$chain")" --node "$(chain_node "$chain")" \
		--gas "$DYDX_GAS" --fees "$fees" --yes --output json 2>&1); then
		echo "$res"
		return 0
	fi
	code=$(echo "$res" | jq -r '.code' 2>/dev/null || echo "?")
	if [ "$code" != "0" ]; then
		echo "$res" | jq -r '.raw_log' 2>/dev/null || echo "$res"
		return 0
	fi
	hash=$(echo "$res" | jq -r '.txhash')
	for i in $(seq 1 30); do
		if out=$(q "$chain" tx "$hash" 2>/dev/null); then
			[ "$(echo "$out" | jq -r '.code')" != "0" ] || die "tx unexpectedly succeeded on $chain: $*"
			echo "$out" | jq -r '.raw_log'
			return 0
		fi
		sleep 1
	done
	die "tx $hash not committed on $chain after 30s"
}

# Fails unless the text contains the expected substring.
assert_contains() {
	case "$2" in
	*"$3"*) ;;
	*) die "$1: '$2' does not contain '$3'" ;;
	esac
}

# The dYdX node's tx_search returns at most one match per query: the kv tx
# indexer in the cometbft fork keeps goleveldb iterator buffers without copying
# them (upstream v0.38.21 copies). Hermes looks packet data up with tx_search
# when clearing after a restart, so it cannot find dYdX packets sent while it
# was down. This relays such a packet (recv or timeout) from block_results at
# the height of its MsgTransfer instead. Takes the committed tx JSON.
relay_dydx_packet() {
	local txjson=$1 height channel sequence dst
	height=$(echo "$txjson" | jq -r '.height')
	read -r channel sequence < <(echo "$txjson" | jq -r '[.events[] | select(.type == "send_packet") |
		.attributes | from_entries][0] | "\(.packet_src_channel) \(.packet_sequence)"')
	case "$channel" in
	"$DYDX_NOBLE_CHANNEL") dst=$NOBLE_CHAIN_ID ;;
	"$DYDX_INJECTIVE_CHANNEL") dst=$INJECTIVE_CHAIN_ID ;;
	*) die "no relay route for dYdX channel '$channel'" ;;
	esac
	"$HERMES_BIN" --config "$HERMES_CONFIG" tx packet-recv --dst-chain "$dst" --src-chain "$DYDX_CHAIN_ID" \
		--src-port transfer --src-channel "$channel" --packet-sequences "$sequence" \
		--packet-data-query-height "$height" >>"$LOG_DIR/hermes-manual.log" 2>&1 ||
		log "manual relay of $channel/$sequence failed (see $LOG_DIR/hermes-manual.log)"
}
