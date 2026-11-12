#!/bin/bash
# Starts a single-validator dYdX chain plus two ibc-go simd chains standing in
# for Noble and Injective, opens dYdX channel-0 to Noble and channel-1 to
# Injective with hermes, and leaves the relayer running.
#
# Any previous chain state under build/ibc-local is wiped; binaries are kept.
set -euo pipefail
source "$(dirname "${BASH_SOURCE[0]}")/lib.sh"

require_cmd jq curl bc shasum
for bin in "$DYDX_BIN" "$SIMD_BIN" "$HERMES_BIN" "$DASEL_BIN"; do
	[ -x "$bin" ] || die "$bin missing; run testing/ibc-local/setup.sh first"
done
export PATH="$BIN_DIR:$PATH"

"$IBC_LOCAL_DIR/down.sh"
command rm -rf "$DYDX_HOME" "$NOBLE_HOME" "$INJECTIVE_HOME" "$HERMES_HOME" "$KEY_DIR" "$LOG_DIR"
mkdir -p "$KEY_DIR" "$LOG_DIR" "$HERMES_HOME"

# Creates a key on the chain and saves its mnemonic so the same identity can be
# recovered on the other chains and in hermes.
new_key() {
	local chain=$1 name=$2
	"$(chain_bin "$chain")" keys add "$name" --keyring-backend test --home "$(chain_home "$chain")" \
		--output json 2>&1 | jq -r '.mnemonic' >"$KEY_DIR/$name.mnemonic"
	[ -s "$KEY_DIR/$name.mnemonic" ] || die "could not create key $name"
}

recover_key() {
	local chain=$1 name=$2
	"$(chain_bin "$chain")" keys add "$name" --recover --keyring-backend test \
		--home "$(chain_home "$chain")" <"$KEY_DIR/$name.mnemonic" >/dev/null 2>&1
}

set_toml() {
	local file=$1 type=$2 selector=$3 value=$4
	"$DASEL_BIN" put -f "$file" -t "$type" -v "$value" "$selector"
}

init_dydx() {
	log "initializing $DYDX_CHAIN_ID"
	"$DYDX_BIN" init alice --chain-id "$DYDX_CHAIN_ID" --home "$DYDX_HOME" >/dev/null 2>&1
	local config="$DYDX_HOME/config"

	# genesis.sh resolves its helpers and env from the testing dir.
	(
		set +u
		cd "$PROTOCOL_DIR/testing"
		source ./genesis.sh
		edit_genesis "$config" "" "" "" "" \
			"$PROTOCOL_DIR/daemons/pricefeed/client/constants/testdata" \
			"$PROTOCOL_DIR/testing/delaymsg_config" "" "" >"$LOG_DIR/dydx-genesis.log" 2>&1
	) || die "edit_genesis failed, see $LOG_DIR/dydx-genesis.log"

	local genesis="$config/genesis.json"
	set_toml "$genesis" string '.app_state.gov.params.voting_period' '20s'
	set_toml "$genesis" string '.app_state.gov.params.expedited_voting_period' '10s'

	new_key dydx validator
	new_key dydx user
	new_key dydx relayer
	new_key dydx holder
	local native="1000000000000000000000000$NATIVE_TOKEN"
	local name
	for name in validator user relayer holder; do
		"$DYDX_BIN" genesis add-genesis-account "$(addr dydx "$name")" "$native" --home "$DYDX_HOME" 2>/dev/null ||
			"$DYDX_BIN" add-genesis-account "$(addr dydx "$name")" "$native" --home "$DYDX_HOME"
	done
	"$DYDX_BIN" gentx validator "500000000000000000000000$NATIVE_TOKEN" --moniker alice \
		--keyring-backend test --chain-id "$DYDX_CHAIN_ID" --home "$DYDX_HOME" >/dev/null 2>&1
	"$DYDX_BIN" collect-gentxs --home "$DYDX_HOME" >/dev/null 2>&1 ||
		"$DYDX_BIN" genesis collect-gentxs --home "$DYDX_HOME" >/dev/null 2>&1

	# No slinky sidecar, price daemon, or ETH RPC locally; prices stay at genesis.
	set_toml "$config/app.toml" string '.oracle.enabled' 'false'
	set_toml "$config/config.toml" string '.rpc.pprof_laddr' "localhost:$DYDX_PPROF_PORT"
}

init_simd() {
	local chain=$1 rpc=$2 p2p=$3 grpc=$4 api=$5 pprof=$6 usdc=$7
	local home id config
	home=$(chain_home "$chain")
	id=$(chain_id "$chain")
	config="$home/config"
	log "initializing $id"
	"$SIMD_BIN" init "$chain-val" --chain-id "$id" --home "$home" >/dev/null 2>&1

	"$SIMD_BIN" keys add validator --keyring-backend test --home "$home" >/dev/null 2>&1
	recover_key "$chain" user
	recover_key "$chain" relayer
	"$SIMD_BIN" genesis add-genesis-account "$(addr "$chain" validator)" 1000000000000stake --home "$home"
	"$SIMD_BIN" genesis add-genesis-account "$(addr "$chain" user)" \
		"1000000000000000$usdc,1000000000stake" --home "$home"
	"$SIMD_BIN" genesis add-genesis-account "$(addr "$chain" relayer)" 1000000000stake --home "$home"
	"$SIMD_BIN" genesis gentx validator 500000000000stake --keyring-backend test \
		--chain-id "$id" --home "$home" >/dev/null 2>&1
	"$SIMD_BIN" genesis collect-gentxs --home "$home" >/dev/null 2>&1
	set_toml "$config/genesis.json" string '.app_state.gov.params.voting_period' '20s'

	set_toml "$config/config.toml" string '.rpc.laddr' "tcp://127.0.0.1:$rpc"
	set_toml "$config/config.toml" string '.rpc.pprof_laddr' "localhost:$pprof"
	set_toml "$config/config.toml" string '.p2p.laddr' "tcp://127.0.0.1:$p2p"
	set_toml "$config/config.toml" string '.consensus.timeout_commit' '1s'
	set_toml "$config/app.toml" string '.grpc.address' "127.0.0.1:$grpc"
	set_toml "$config/app.toml" string '.api.address' "tcp://127.0.0.1:$api"
	set_toml "$config/app.toml" string '.minimum-gas-prices' '0stake'
}

write_hermes_config() {
	cat >"$HERMES_CONFIG" <<EOF
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
port = 3000

[telemetry]
enabled = false
host = '127.0.0.1'
port = 3001
EOF
	hermes_chain "$DYDX_CHAIN_ID" "$DYDX_RPC_PORT" "$DYDX_GRPC_PORT" dydx 25000000000 "$NATIVE_TOKEN"
	hermes_chain "$NOBLE_CHAIN_ID" "$NOBLE_RPC_PORT" "$NOBLE_GRPC_PORT" cosmos 0 stake
	hermes_chain "$INJECTIVE_CHAIN_ID" "$INJECTIVE_RPC_PORT" "$INJECTIVE_GRPC_PORT" cosmos 0 stake
}

hermes_chain() {
	local id=$1 rpc=$2 grpc=$3 prefix=$4 price=$5 denom=$6
	cat >>"$HERMES_CONFIG" <<EOF

[[chains]]
id = '$id'
type = 'CosmosSdk'
rpc_addr = 'http://127.0.0.1:$rpc'
grpc_addr = 'http://127.0.0.1:$grpc'
event_source = { mode = 'push', url = 'ws://127.0.0.1:$rpc/websocket', batch_delay = '200ms' }
rpc_timeout = '10s'
account_prefix = '$prefix'
key_name = 'relayer'
store_prefix = 'ibc'
default_gas = 500000
max_gas = 10000000
gas_price = { price = $price, denom = '$denom' }
gas_multiplier = 1.5
max_msg_num = 30
max_tx_size = 2097152
clock_drift = '10s'
max_block_time = '30s'
trust_threshold = '2/3'
address_type = { derivation = 'cosmos' }
EOF
}

open_channel() {
	local counterparty=$1 want=$2 out got
	log "opening transfer channel $DYDX_CHAIN_ID <-> $counterparty"
	out=$(hermes --json create channel --a-chain "$DYDX_CHAIN_ID" --b-chain "$counterparty" \
		--a-port transfer --b-port transfer --new-client-connection --yes 2>>"$LOG_DIR/hermes-setup.log" | tail -1)
	got=$(echo "$out" | jq -r '.result.a_side.channel_id')
	[ "$got" = "$want" ] || die "expected dYdX $want for $counterparty, got '$got': $out"
}

init_dydx
init_simd noble "$NOBLE_RPC_PORT" "$NOBLE_P2P_PORT" "$NOBLE_GRPC_PORT" "$NOBLE_API_PORT" "$NOBLE_PPROF_PORT" \
	"$NOBLE_USDC"
init_simd injective "$INJECTIVE_RPC_PORT" "$INJECTIVE_P2P_PORT" "$INJECTIVE_GRPC_PORT" \
	"$INJECTIVE_API_PORT" "$INJECTIVE_PPROF_PORT" "$INJECTIVE_USDC"

start_dydx
start_process noble "$SIMD_BIN" start --home "$NOBLE_HOME"
start_process injective "$SIMD_BIN" start --home "$INJECTIVE_HOME"
wait_for_height dydx 3
wait_for_height noble 3
wait_for_height injective 3
log "all chains producing blocks"

write_hermes_config
for id in "$DYDX_CHAIN_ID" "$NOBLE_CHAIN_ID" "$INJECTIVE_CHAIN_ID"; do
	hermes keys add --chain "$id" --key-name relayer --mnemonic-file "$KEY_DIR/relayer.mnemonic" \
		--overwrite >>"$LOG_DIR/hermes-setup.log" 2>&1 || die "hermes key import failed for $id"
done

# Order matters: channel-0 must be Noble so the trace matches USDC_DENOM.
open_channel "$NOBLE_CHAIN_ID" "$DYDX_NOBLE_CHANNEL"
open_channel "$INJECTIVE_CHAIN_ID" "$DYDX_INJECTIVE_CHANNEL"
: >"$CHANNELS_ENV"
# upgrade.sh defers the alternate Injective channel until after the upgrade,
# so the handler finds exactly one channel to Injective.
[ -n "${IBC_LOCAL_DEFER_ALT_CHANNEL:-}" ] || open_alt_injective_channel
# Injective<->Noble, to route Noble-origin USDC to dYdX through Injective.
open_extra_channel INJECTIVE_NOBLE "$INJECTIVE_CHAIN_ID" "$NOBLE_CHAIN_ID" ""
start_hermes

log "environment up"
log "  dydx      rpc :$DYDX_RPC_PORT  user $(addr dydx user)"
log "  noble     rpc :$NOBLE_RPC_PORT  user $(addr noble user)"
log "  injective rpc :$INJECTIVE_RPC_PORT  user $(addr injective user)"
log "next: testing/ibc-local/test.sh   stop: testing/ibc-local/down.sh"
