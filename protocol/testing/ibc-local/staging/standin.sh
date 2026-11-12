#!/bin/bash
# Injective stand-in for a remote dYdX chain: an ibc-go simd with chain id
# injective-1 minting native USDC's Injective denom, plus a hermes relayer, and
# one transfer channel to the remote chain. State lives under STAGING_DIR;
# binaries come from the local ibc-local build.
set -euo pipefail
HARNESS="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
STAGING_DIR="${STAGING_DIR:-$HOME/.ibc-staging}"
export IBC_LOCAL_STATE_DIR="$STAGING_DIR"
source "$HARNESS/lib.sh"
BIN="$HARNESS/../../build/ibc-local/bin"
DYDX_BIN="$BIN/dydxprotocold"; SIMD_BIN="$BIN/ibc-simd"; HERMES_BIN="$BIN/hermes"; DASEL_BIN="$BIN/dasel"
HERMES_HOME="$STAGING_DIR/hermes"; HERMES_CONFIG="$HERMES_HOME/config.toml"

REMOTE_CHAIN_ID="${REMOTE_CHAIN_ID:-dydx-testnet-4}"
REMOTE_RPC="${REMOTE_RPC:-http://3.22.142.107:26657}"
REMOTE_GRPC="${REMOTE_GRPC:-http://3.22.142.107:9090}"
REMOTE_WS="${REMOTE_WS:-ws://3.22.142.107:26657/websocket}"
FUNDER_HOME="${FUNDER_HOME:-$HOME/.dydxprotocol}"
FUNDER_KEY="${FUNDER_KEY:-dydx-1}"
INJECTIVE_HOME="$STAGING_DIR/injective"
INJECTIVE_RPC_PORT=29657; INJECTIVE_P2P_PORT=29656; INJECTIVE_GRPC_PORT=9390; INJECTIVE_API_PORT=1617; INJECTIVE_PPROF_PORT=6360
DYDX_KEYS="$STAGING_DIR/dydx-keys"
mkdir -p "$STAGING_DIR" "$KEY_DIR" "$LOG_DIR" "$PID_DIR" "$HERMES_HOME" "$DYDX_KEYS"

set_toml() { "$DASEL_BIN" put -f "$1" -t "$2" -v "$4" "$3"; }
hermes() { "$HERMES_BIN" --config "$HERMES_CONFIG" "$@"; }

if [ ! -f "$INJECTIVE_HOME/config/genesis.json" ]; then
	log "initializing $INJECTIVE_CHAIN_ID stand-in"
	"$SIMD_BIN" init injective-standin --chain-id "$INJECTIVE_CHAIN_ID" --home "$INJECTIVE_HOME" >/dev/null 2>&1
	for name in validator user relayer; do
		"$SIMD_BIN" keys add "$name" --keyring-backend test --home "$INJECTIVE_HOME" --output json 2>&1 |
			jq -r '.mnemonic' >"$KEY_DIR/injective-$name.mnemonic"
	done
	addr_inj() { "$SIMD_BIN" keys show "$1" -a --keyring-backend test --home "$INJECTIVE_HOME"; }
	"$SIMD_BIN" genesis add-genesis-account "$(addr_inj validator)" 1000000000000stake --home "$INJECTIVE_HOME"
	"$SIMD_BIN" genesis add-genesis-account "$(addr_inj user)" "1000000000000000$INJECTIVE_USDC,1000000000stake" --home "$INJECTIVE_HOME"
	"$SIMD_BIN" genesis add-genesis-account "$(addr_inj relayer)" 1000000000stake --home "$INJECTIVE_HOME"
	"$SIMD_BIN" genesis gentx validator 500000000000stake --keyring-backend test --chain-id "$INJECTIVE_CHAIN_ID" --home "$INJECTIVE_HOME" >/dev/null 2>&1
	"$SIMD_BIN" genesis collect-gentxs --home "$INJECTIVE_HOME" >/dev/null 2>&1
	c="$INJECTIVE_HOME/config"
	set_toml "$c/config.toml" string '.rpc.laddr' "tcp://127.0.0.1:$INJECTIVE_RPC_PORT"
	set_toml "$c/config.toml" string '.rpc.pprof_laddr' "localhost:$INJECTIVE_PPROF_PORT"
	set_toml "$c/config.toml" string '.p2p.laddr' "tcp://127.0.0.1:$INJECTIVE_P2P_PORT"
	set_toml "$c/config.toml" string '.consensus.timeout_commit' '1s'
	set_toml "$c/app.toml" string '.grpc.address' "127.0.0.1:$INJECTIVE_GRPC_PORT"
	set_toml "$c/app.toml" string '.api.address' "tcp://127.0.0.1:$INJECTIVE_API_PORT"
	set_toml "$c/app.toml" string '.minimum-gas-prices' '0stake'
fi
start_process injective "$SIMD_BIN" start --home "$INJECTIVE_HOME"
for i in $(seq 1 60); do
	h=$(curl -sf "http://127.0.0.1:$INJECTIVE_RPC_PORT/status" | jq -r '.result.sync_info.latest_block_height' 2>/dev/null || echo 0)
	[ "${h:-0}" != "null" ] && [ "${h:-0}" -ge 3 ] && break
	sleep 1
done
log "$INJECTIVE_CHAIN_ID at height $h"

# Relayer account on the remote chain, funded from FUNDER_KEY.
if [ ! -f "$KEY_DIR/dydx-relayer.mnemonic" ]; then
	"$DYDX_BIN" keys add relayer --keyring-backend test --home "$DYDX_KEYS" --output json 2>&1 |
		jq -r '.mnemonic' >"$KEY_DIR/dydx-relayer.mnemonic"
fi
RELAYER=$("$DYDX_BIN" keys show relayer -a --keyring-backend test --home "$DYDX_KEYS")
bal=$("$DYDX_BIN" query bank balances "$RELAYER" --node "$REMOTE_RPC" -o json | jq -r '[.balances[] | select(.denom == "adv4tnt") | .amount][0] // "0"')
if [ "$(echo "$bal < 50000000000000000000" | bc)" = 1 ]; then
	log "funding relayer $RELAYER with 100 DV4TNT from $FUNDER_KEY"
	"$DYDX_BIN" tx bank send "$FUNDER_KEY" "$RELAYER" 100000000000000000000adv4tnt --keyring-backend test --home "$FUNDER_HOME" \
		--node "$REMOTE_RPC" --chain-id "$REMOTE_CHAIN_ID" --gas 200000 --gas-prices 25000000000adv4tnt -y -o json | jq -r '.txhash'
	sleep 6
fi
log "relayer $RELAYER balance: $("$DYDX_BIN" query bank balances "$RELAYER" --node "$REMOTE_RPC" -o json | jq -c '.balances')"

cat >"$HERMES_CONFIG" <<CFG
[global]
log_level = 'info'
[mode.clients]
enabled = true
refresh = true
misbehaviour = false
[mode.connections]
enabled = false
[mode.channels]
enabled = false
[mode.packets]
enabled = true
clear_interval = 10
clear_on_start = true
tx_confirmation = true
[rest]
enabled = false
host = '127.0.0.1'
port = 3100
[telemetry]
enabled = false
host = '127.0.0.1'
port = 3101

[[chains]]
id = '$REMOTE_CHAIN_ID'
type = 'CosmosSdk'
rpc_addr = '$REMOTE_RPC'
grpc_addr = '$REMOTE_GRPC'
event_source = { mode = 'push', url = '$REMOTE_WS', batch_delay = '500ms' }
rpc_timeout = '15s'
account_prefix = 'dydx'
key_name = 'relayer'
key_store_folder = '$HERMES_HOME/keys'
store_prefix = 'ibc'
default_gas = 500000
max_gas = 10000000
gas_price = { price = 25000000000, denom = 'adv4tnt' }
gas_multiplier = 1.5
max_msg_num = 30
max_tx_size = 2097152
clock_drift = '15s'
max_block_time = '30s'
trust_threshold = '2/3'
address_type = { derivation = 'cosmos' }

[[chains]]
id = '$INJECTIVE_CHAIN_ID'
type = 'CosmosSdk'
rpc_addr = 'http://127.0.0.1:$INJECTIVE_RPC_PORT'
grpc_addr = 'http://127.0.0.1:$INJECTIVE_GRPC_PORT'
event_source = { mode = 'push', url = 'ws://127.0.0.1:$INJECTIVE_RPC_PORT/websocket', batch_delay = '200ms' }
rpc_timeout = '10s'
account_prefix = 'cosmos'
key_name = 'relayer'
key_store_folder = '$HERMES_HOME/keys'
store_prefix = 'ibc'
default_gas = 500000
max_gas = 10000000
gas_price = { price = 0, denom = 'stake' }
gas_multiplier = 1.5
max_msg_num = 30
max_tx_size = 2097152
clock_drift = '10s'
max_block_time = '30s'
trust_threshold = '2/3'
address_type = { derivation = 'cosmos' }
CFG
hermes keys add --chain "$REMOTE_CHAIN_ID" --key-name relayer --mnemonic-file "$KEY_DIR/dydx-relayer.mnemonic" --overwrite >>"$LOG_DIR/hermes-setup.log" 2>&1
hermes keys add --chain "$INJECTIVE_CHAIN_ID" --key-name relayer --mnemonic-file "$KEY_DIR/injective-relayer.mnemonic" --overwrite >>"$LOG_DIR/hermes-setup.log" 2>&1
hermes health-check 2>&1 | tail -3

if [ ! -f "$STAGING_DIR/channel.env" ]; then
	log "opening transfer channel $INJECTIVE_CHAIN_ID <-> $REMOTE_CHAIN_ID"
	out=$(hermes --json create channel --a-chain "$INJECTIVE_CHAIN_ID" --b-chain "$REMOTE_CHAIN_ID" \
		--a-port transfer --b-port transfer --new-client-connection --yes 2>>"$LOG_DIR/hermes-setup.log" | tail -1)
	echo "$out" | jq -c '.result | {a: .a_side, b: .b_side}' || die "channel creation failed: $out"
	{
		echo "INJECTIVE_SIDE_CHANNEL=$(echo "$out" | jq -r '.result.a_side.channel_id')"
		echo "INJECTIVE_SIDE_CONNECTION=$(echo "$out" | jq -r '.result.a_side.connection_id')"
		echo "INJECTIVE_SIDE_CLIENT=$(echo "$out" | jq -r '.result.a_side.client_id')"
		echo "DYDX_SIDE_CHANNEL=$(echo "$out" | jq -r '.result.b_side.channel_id')"
		echo "DYDX_SIDE_CONNECTION=$(echo "$out" | jq -r '.result.b_side.connection_id')"
		echo "DYDX_SIDE_CLIENT=$(echo "$out" | jq -r '.result.b_side.client_id')"
	} >"$STAGING_DIR/channel.env"
fi
cat "$STAGING_DIR/channel.env"
start_process hermes "$HERMES_BIN" --config "$HERMES_CONFIG" start
log "stand-in up; state in $STAGING_DIR"
