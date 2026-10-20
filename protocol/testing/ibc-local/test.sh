#!/bin/bash
# Runs canonical USDC scenarios against the environment started by up.sh.
# Every step relays real ICS20 packets through hermes and asserts balances on
# both sides plus the canonicalusdc ledger and backing equation.
#
# Run against a freshly started environment: the scenarios walk the module from
# MODE_DISABLED through GRADUAL, PAUSED and back, and expect the initial state.
set -euo pipefail
source "$(dirname "${BASH_SOURCE[0]}")/lib.sh"

require_cmd jq bc shasum curl
for chain in dydx noble injective; do
	pid_alive "$chain" || die "$chain is not running; run testing/ibc-local/up.sh"
done
pid_alive hermes || die "hermes is not running; run testing/ibc-local/up.sh"

INJECTIVE_DENOM=$(ibc_denom "transfer/$DYDX_INJECTIVE_CHANNEL/$COUNTERPARTY_USDC")
DYDX_USER=$(addr dydx user)
NOBLE_USER=$(addr noble user)
INJECTIVE_USER=$(addr injective user)
DYDX_HOLDER=$(addr dydx holder)
[ -f "$CHANNELS_ENV" ] || die "$CHANNELS_ENV missing; rerun testing/ibc-local/up.sh"
source "$CHANNELS_ENV"
: "${DYDX_INJECTIVE_ALT_A:?}" "${DYDX_INJECTIVE_ALT_B:?}" "${INJECTIVE_NOBLE_A:?}" "${INJECTIVE_NOBLE_B:?}"
WORK_DIR="$STATE_DIR/test"
mkdir -p "$WORK_DIR"

assert_eq "noble route trace" "$(ibc_denom "transfer/$DYDX_NOBLE_CHANNEL/$COUNTERPARTY_USDC")" "$USDC_DENOM"

module_address() {
	q dydx auth module-account "$1" |
		jq -r '.account.value.address // .account.base_account.address // .account.address'
}
CANONICAL_MODULE=$(module_address canonicalusdc)
GOV_MODULE=$(module_address gov)

counterparty_channel() {
	q dydx ibc channel end transfer "$1" | jq -r '.channel.counterparty.channel_id'
}
NOBLE_CHANNEL=$(counterparty_channel "$DYDX_NOBLE_CHANNEL")
INJECTIVE_CHANNEL=$(counterparty_channel "$DYDX_INJECTIVE_CHANNEL")
# Logical USDC that left over the Injective channel as a legacy ICS20 voucher.
LOGICAL_ON_INJECTIVE=$(ibc_denom "transfer/$INJECTIVE_CHANNEL/transfer/$DYDX_NOBLE_CHANNEL/$COUNTERPARTY_USDC")

state() { q dydx canonicalusdc state; }
ledger_field() { state | jq -r ".ledger.$1 // \"0\""; }
pending_count() { q dydx canonicalusdc pending-settlements | jq -r '.pending_settlements | length'; }

wait_for_no_pending() {
	local i n
	for i in $(seq 1 "$RELAY_TIMEOUT_SECONDS"); do
		n=$(pending_count)
		[ "$n" = "0" ] && return 0
		sleep 1
	done
	die "$n canonical settlements still pending after ${RELAY_TIMEOUT_SECONDS}s"
}

# Asserts the module's own invariant from outside: in GRADUAL, logical supply
# equals noble + injective + pending_injective, and module balances cover the
# physical reserve and the locked pending liability.
check_backing() {
	local s supply noble inj pend_inj phys locked
	s=$(state)
	[ "$(echo "$s" | jq -r '.controls.mode')" = "MODE_DISABLED" ] && return 0
	supply=$(supply_of dydx "$USDC_DENOM")
	noble=$(echo "$s" | jq -r '.ledger.noble_backing')
	inj=$(echo "$s" | jq -r '.ledger.injective_backing')
	pend_inj=$(echo "$s" | jq -r '.ledger.pending_injective')
	phys=$(balance dydx "$CANONICAL_MODULE" "$INJECTIVE_DENOM")
	locked=$(balance dydx "$CANONICAL_MODULE" "$USDC_DENOM")
	[ "$(bc_calc "$phys >= $inj")" = 1 ] || die "physical reserve $phys below injective backing $inj"
	[ "$(bc_calc "$locked >= $pend_inj")" = 1 ] || die "locked logical $locked below pending $pend_inj"
	if [ "$(echo "$s" | jq -r '.controls.mode')" = "MODE_GRADUAL" ]; then
		assert_eq "logical supply == accounted backing" "$supply" "$(bc_calc "$noble + $inj + $pend_inj")"
	fi
}

# transfer CHAIN FROM CHANNEL RECEIVER AMOUNT [extra flags]
transfer() {
	local chain=$1 from=$2 channel=$3 receiver=$4 amount=$5
	shift 5
	tx "$chain" "$from" ibc-transfer transfer transfer "$channel" "$receiver" "$amount" "$@" >/dev/null
}

# Like transfer, but prints the committed tx JSON (for relay_dydx_packet).
transfer_tx() {
	local chain=$1 from=$2 channel=$3 receiver=$4 amount=$5
	shift 5
	tx "$chain" "$from" ibc-transfer transfer transfer "$channel" "$receiver" "$amount" "$@"
}

# Per-packet cap and Injective backing ceiling for GRADUAL. Kept small enough
# that the limit scenarios can exceed them with test balances.
MAX_TRANSFER=1000000000
MIGRATION_CEILING=1000000000

# propose_controls MODE EXPECTED [JQ_FILTER]
# Submits MsgUpdateControls through governance with the live route bindings,
# optionally rewritten by JQ_FILTER (applied to the controls object), votes yes
# and waits for EXPECTED: PASSED or FAILED (message execution rejected).
propose_controls() {
	local mode=$1 expected=$2 filter=${3:-.} noble_conn noble_client inj_conn inj_client id i status
	noble_conn=$(q dydx ibc channel end transfer "$DYDX_NOBLE_CHANNEL" | jq -r '.channel.connection_hops[0]')
	inj_conn=$(q dydx ibc channel end transfer "$DYDX_INJECTIVE_CHANNEL" | jq -r '.channel.connection_hops[0]')
	noble_client=$(q dydx ibc connection end "$noble_conn" | jq -r '.connection.client_id')
	inj_client=$(q dydx ibc connection end "$inj_conn" | jq -r '.connection.client_id')
	jq -n --arg gov "$GOV_MODULE" --arg mode "$mode" --arg logical "$USDC_DENOM" \
		--arg nch "$DYDX_NOBLE_CHANNEL" --arg ncl "$noble_client" --arg nco "$noble_conn" \
		--arg ich "$DYDX_INJECTIVE_CHANNEL" --arg icl "$inj_client" --arg ico "$inj_conn" \
		--arg idenom "$INJECTIVE_DENOM" --arg pd "$COUNTERPARTY_USDC" --arg dep "10000000$NATIVE_TOKEN" \
		--arg max "$MAX_TRANSFER" --arg ceiling "$MIGRATION_CEILING" '{
		messages: [{
			"@type": "/dydxprotocol.canonicalusdc.MsgUpdateControls",
			authority: $gov,
			controls: ({
				mode: $mode, logical_denom: $logical,
				noble_channel: $nch, noble_client: $ncl, noble_connection: $nco,
				noble_packet_denom: $pd, noble_withdrawals_enabled: true,
				injective_channel: $ich, injective_client: $icl, injective_connection: $ico,
				injective_packet_denom: $pd, injective_denom: $idenom,
				max_transfer_amount: $max, migration_ceiling: $ceiling,
				max_pending_settlements: 100
			} | '"$filter"')
		}],
		deposit: $dep, title: ("canonical USDC " + $mode), summary: "ibc-local test", expedited: false
	}' >"$WORK_DIR/proposal.json"

	tx dydx validator gov submit-proposal "$WORK_DIR/proposal.json" >/dev/null
	id=$(q dydx gov proposals | jq -r '.proposals | map(.id | tonumber) | max')
	tx dydx validator gov vote "$id" yes >/dev/null
	status=""
	for i in $(seq 1 "$RELAY_TIMEOUT_SECONDS"); do
		status=$(q dydx gov proposal "$id" | jq -r '.proposal.status // .status')
		case "$status" in
		PROPOSAL_STATUS_PASSED | 3 | PROPOSAL_STATUS_REJECTED | 4 | PROPOSAL_STATUS_FAILED | 5) break ;;
		esac
		sleep 1
	done
	case "$expected:$status" in
	PASSED:PROPOSAL_STATUS_PASSED | PASSED:3 | FAILED:PROPOSAL_STATUS_FAILED | FAILED:5) ;;
	*)
		die "proposal $id ($mode) ended '$status', want $expected: $(q dydx gov proposal "$id" |
			jq -c '.proposal.failed_reason // .failed_reason // empty')"
		;;
	esac
	echo "$id"
}

submit_controls() {
	propose_controls "$1" PASSED >/dev/null
	assert_eq "controls mode" "$(state | jq -r '.controls.mode')" "$1"
}

ledger_json() { state | jq -c '.ledger'; }
controls_json() { state | jq -c '.controls'; }

scenario_disabled_noble_roundtrip() {
	log "DISABLED: Noble deposit and legacy withdrawal"
	local noble_before dydx_before
	noble_before=$(balance noble "$NOBLE_USER" "$COUNTERPARTY_USDC")
	dydx_before=$(balance dydx "$DYDX_USER" "$USDC_DENOM")
	transfer noble user "$NOBLE_CHANNEL" "$DYDX_USER" "1000000000$COUNTERPARTY_USDC"
	wait_for_balance dydx "$DYDX_USER" "$USDC_DENOM" "$(bc_calc "$dydx_before + 1000000000")"
	pass "Noble uusdc arrives on dYdX as $USDC_DENOM"

	transfer dydx user "$DYDX_NOBLE_CHANNEL" "$NOBLE_USER" "100000000$USDC_DENOM"
	wait_for_balance noble "$NOBLE_USER" "$COUNTERPARTY_USDC" "$(bc_calc "$noble_before - 900000000")"
	assert_eq "dYdX user USDC" "$(balance dydx "$DYDX_USER" "$USDC_DENOM")" "$(bc_calc "$dydx_before + 900000000")"
	pass "legacy withdrawal to Noble unescrows native uusdc"
}

scenario_enable_gradual() {
	log "governance: enable MODE_GRADUAL"
	submit_controls MODE_GRADUAL
	assert_eq "noble_backing snapshot" "$(ledger_field noble_backing)" "$(supply_of dydx "$USDC_DENOM")"
	check_backing
	pass "GRADUAL entered, Noble backing snapshotted to logical supply"
}

# A Noble withdrawal sent while DISABLED has no pending record, and the GRADUAL
# snapshot excludes its burned amount. Its timeout lands after GRADUAL, so the
# ICS20 re-mint must be credited back to noble_backing. Enables GRADUAL itself.
scenario_disabled_noble_refunded_after_gradual() {
	log "DISABLED Noble withdrawal times out after GRADUAL is enabled"
	local amount=70000000 timeout_seconds=15 dydx_before noble_before supply_before sent sent_at backing remaining
	assert_eq "mode" "$(state | jq -r '.controls.mode')" MODE_DISABLED
	dydx_before=$(balance dydx "$DYDX_USER" "$USDC_DENOM")
	noble_before=$(balance noble "$NOBLE_USER" "$COUNTERPARTY_USDC")
	supply_before=$(supply_of dydx "$USDC_DENOM")

	stop_hermes
	sent_at=$(date +%s)
	sent=$(transfer_tx dydx user "$DYDX_NOBLE_CHANNEL" "$NOBLE_USER" "$amount$USDC_DENOM" \
		--packet-timeout-height 0-0 --packet-timeout-timestamp "${timeout_seconds}000000000")
	assert_eq "legacy send burned" "$(supply_of dydx "$USDC_DENOM")" "$(bc_calc "$supply_before - $amount")"
	assert_eq "no pending record" "$(pending_count)" 0
	assert_eq "legacy packet in flight" "$(legacy_noble_in_flight)" 1

	scenario_enable_gradual
	backing=$(ledger_field noble_backing)
	assert_eq "snapshot excludes the burned amount" "$backing" "$(bc_calc "$supply_before - $amount")"

	log "waiting for the packet to expire on Noble"
	remaining=$((sent_at + timeout_seconds + 10 - $(date +%s)))
	if [ "$remaining" -gt 0 ]; then sleep "$remaining"; fi
	start_hermes
	relay_dydx_packet "$sent"
	wait_for_balance dydx "$DYDX_USER" "$USDC_DENOM" "$dydx_before"
	assert_eq "noble_backing credited" "$(ledger_field noble_backing)" "$(bc_calc "$backing + $amount")"
	assert_eq "supply restored" "$(supply_of dydx "$USDC_DENOM")" "$supply_before"
	assert_eq "Noble user unchanged" "$(balance noble "$NOBLE_USER" "$COUNTERPARTY_USDC")" "$noble_before"
	assert_eq "legacy packet settled" "$(legacy_noble_in_flight)" 0
	assert_eq "still GRADUAL after EndBlock" "$(state | jq -r '.controls.mode')" MODE_GRADUAL
	check_backing
	pass "legacy Noble refund re-mints to the user and restores Noble backing"
}

# Counts packets awaiting ack or timeout on the Noble channel that have no
# pending settlement: Noble withdrawals sent while DISABLED. Each one re-adds
# Noble backing if it is refunded, so migration finalization requires zero.
legacy_noble_in_flight() {
	local commitments pending
	commitments=$(q dydx ibc channel packet-commitments transfer "$DYDX_NOBLE_CHANNEL" |
		jq -c '[.commitments[].sequence | tostring]')
	pending=$(q dydx canonicalusdc pending-settlements |
		jq -c --arg ch "$DYDX_NOBLE_CHANNEL" '[.pending_settlements[] | select(.source_channel == $ch) | .sequence | tostring]')
	jq -n --argjson c "$commitments" --argjson p "$pending" '$c - $p | length'
}

scenario_gradual_noble_deposit_rejected() {
	log "GRADUAL: Noble deposit is refused with an error ack"
	local noble_before dydx_before
	noble_before=$(balance noble "$NOBLE_USER" "$COUNTERPARTY_USDC")
	dydx_before=$(balance dydx "$DYDX_USER" "$USDC_DENOM")
	transfer noble user "$NOBLE_CHANNEL" "$DYDX_USER" "50000000$COUNTERPARTY_USDC"
	assert_eq "Noble escrowed" "$(balance noble "$NOBLE_USER" "$COUNTERPARTY_USDC")" "$(bc_calc "$noble_before - 50000000")"
	wait_for_balance noble "$NOBLE_USER" "$COUNTERPARTY_USDC" "$noble_before"
	assert_eq "dYdX user USDC" "$(balance dydx "$DYDX_USER" "$USDC_DENOM")" "$dydx_before"
	check_backing
	pass "Noble deposit refunded on Noble, nothing minted on dYdX"
}

scenario_gradual_injective_roundtrip() {
	log "GRADUAL: Injective deposit mints logical USDC, withdrawal returns native uusdc"
	local inj_before dydx_before supply_before
	inj_before=$(balance injective "$INJECTIVE_USER" "$COUNTERPARTY_USDC")
	dydx_before=$(balance dydx "$DYDX_USER" "$USDC_DENOM")
	supply_before=$(supply_of dydx "$USDC_DENOM")

	transfer injective user "$INJECTIVE_CHANNEL" "$DYDX_USER" "500000000$COUNTERPARTY_USDC"
	wait_for_balance dydx "$DYDX_USER" "$USDC_DENOM" "$(bc_calc "$dydx_before + 500000000")"
	assert_eq "user holds no Injective voucher" "$(balance dydx "$DYDX_USER" "$INJECTIVE_DENOM")" 0
	assert_eq "module physical reserve" "$(balance dydx "$CANONICAL_MODULE" "$INJECTIVE_DENOM")" 500000000
	assert_eq "injective_backing" "$(ledger_field injective_backing)" 500000000
	assert_eq "logical supply" "$(supply_of dydx "$USDC_DENOM")" "$(bc_calc "$supply_before + 500000000")"
	check_backing
	pass "Injective deposit credited as logical USDC, physical voucher held by module"

	transfer dydx user "$DYDX_INJECTIVE_CHANNEL" "$INJECTIVE_USER" "200000000$USDC_DENOM"
	wait_for_balance injective "$INJECTIVE_USER" "$COUNTERPARTY_USDC" "$(bc_calc "$inj_before - 300000000")"
	wait_for_no_pending
	assert_eq "injective_backing" "$(ledger_field injective_backing)" 300000000
	assert_eq "pending_injective" "$(ledger_field pending_injective)" 0
	assert_eq "module locked logical" "$(balance dydx "$CANONICAL_MODULE" "$USDC_DENOM")" 0
	assert_eq "logical supply" "$(supply_of dydx "$USDC_DENOM")" "$(bc_calc "$supply_before + 300000000")"
	check_backing
	pass "logical withdrawal to Injective delivers native uusdc and burns on ack"
}

scenario_gradual_noble_withdrawal() {
	log "GRADUAL: logical withdrawal over Noble"
	local noble_before backing_before
	noble_before=$(balance noble "$NOBLE_USER" "$COUNTERPARTY_USDC")
	backing_before=$(ledger_field noble_backing)
	transfer dydx user "$DYDX_NOBLE_CHANNEL" "$NOBLE_USER" "300000000$USDC_DENOM"
	wait_for_balance noble "$NOBLE_USER" "$COUNTERPARTY_USDC" "$(bc_calc "$noble_before + 300000000")"
	wait_for_no_pending
	assert_eq "noble_backing" "$(ledger_field noble_backing)" "$(bc_calc "$backing_before - 300000000")"
	assert_eq "pending_noble" "$(ledger_field pending_noble)" 0
	check_backing
	pass "Noble withdrawal settles and reduces Noble backing"
}

scenario_gradual_noble_timeout() {
	log "GRADUAL: Noble withdrawal times out and is refunded"
	local dydx_before noble_before backing_before supply_before sent
	dydx_before=$(balance dydx "$DYDX_USER" "$USDC_DENOM")
	noble_before=$(balance noble "$NOBLE_USER" "$COUNTERPARTY_USDC")
	backing_before=$(ledger_field noble_backing)
	supply_before=$(supply_of dydx "$USDC_DENOM")

	stop_hermes
	# Relative timeout: 15s in nanoseconds, no height timeout.
	sent=$(transfer_tx dydx user "$DYDX_NOBLE_CHANNEL" "$NOBLE_USER" "80000000$USDC_DENOM" \
		--packet-timeout-height 0-0 --packet-timeout-timestamp 15000000000)
	assert_eq "pending_noble in flight" "$(ledger_field pending_noble)" 80000000
	assert_eq "noble_backing reduced" "$(ledger_field noble_backing)" "$(bc_calc "$backing_before - 80000000")"
	assert_eq "dYdX user debited" "$(balance dydx "$DYDX_USER" "$USDC_DENOM")" "$(bc_calc "$dydx_before - 80000000")"
	check_backing
	log "waiting for the packet to expire on Noble"
	sleep 25
	start_hermes
	relay_dydx_packet "$sent"

	wait_for_no_pending
	wait_for_balance dydx "$DYDX_USER" "$USDC_DENOM" "$dydx_before"
	assert_eq "noble_backing restored" "$(ledger_field noble_backing)" "$backing_before"
	assert_eq "pending_noble" "$(ledger_field pending_noble)" 0
	assert_eq "supply restored" "$(supply_of dydx "$USDC_DENOM")" "$supply_before"
	assert_eq "Noble user unchanged" "$(balance noble "$NOBLE_USER" "$COUNTERPARTY_USDC")" "$noble_before"
	assert_eq "mode" "$(state | jq -r '.controls.mode')" MODE_GRADUAL
	check_backing
	pass "timed-out Noble withdrawal refunds logical USDC and restores Noble backing"
}

scenario_gradual_noble_invalid_receiver() {
	log "GRADUAL: Noble withdrawal to an invalid receiver is refunded"
	local dydx_before ledger noble_before supply_before
	dydx_before=$(balance dydx "$DYDX_USER" "$USDC_DENOM")
	noble_before=$(balance noble "$NOBLE_USER" "$COUNTERPARTY_USDC")
	supply_before=$(supply_of dydx "$USDC_DENOM")
	ledger=$(ledger_json)
	transfer dydx user "$DYDX_NOBLE_CHANNEL" "not-a-bech32-address" "40000000$USDC_DENOM"
	wait_for_no_pending
	wait_for_balance dydx "$DYDX_USER" "$USDC_DENOM" "$dydx_before"
	assert_eq "ledger" "$(ledger_json)" "$ledger"
	assert_eq "supply restored" "$(supply_of dydx "$USDC_DENOM")" "$supply_before"
	assert_eq "Noble user unchanged" "$(balance noble "$NOBLE_USER" "$COUNTERPARTY_USDC")" "$noble_before"
	assert_eq "mode" "$(state | jq -r '.controls.mode')" MODE_GRADUAL
	check_backing
	pass "error ack refunds logical USDC and restores Noble backing"
}

scenario_gradual_injective_timeout() {
	log "GRADUAL: Injective withdrawal times out and is refunded"
	local dydx_before inj_before backing_before sent
	dydx_before=$(balance dydx "$DYDX_USER" "$USDC_DENOM")
	inj_before=$(balance injective "$INJECTIVE_USER" "$COUNTERPARTY_USDC")
	backing_before=$(ledger_field injective_backing)

	stop_hermes
	# Relative timeout: 15s in nanoseconds, no height timeout.
	sent=$(transfer_tx dydx user "$DYDX_INJECTIVE_CHANNEL" "$INJECTIVE_USER" "100000000$USDC_DENOM" \
		--packet-timeout-height 0-0 --packet-timeout-timestamp 15000000000)
	assert_eq "pending_injective in flight" "$(ledger_field pending_injective)" 100000000
	assert_eq "dYdX user debited" "$(balance dydx "$DYDX_USER" "$USDC_DENOM")" "$(bc_calc "$dydx_before - 100000000")"
	check_backing
	log "waiting for the packet to expire on Injective"
	sleep 25
	start_hermes
	relay_dydx_packet "$sent"

	wait_for_no_pending
	wait_for_balance dydx "$DYDX_USER" "$USDC_DENOM" "$dydx_before"
	assert_eq "injective_backing restored" "$(ledger_field injective_backing)" "$backing_before"
	assert_eq "pending_injective" "$(ledger_field pending_injective)" 0
	assert_eq "Injective user unchanged" "$(balance injective "$INJECTIVE_USER" "$COUNTERPARTY_USDC")" "$inj_before"
	check_backing
	pass "timed-out Injective withdrawal refunds logical USDC and restores backing"
}

scenario_disabled_legacy_positions() {
	log "DISABLED: legacy Injective voucher and logical USDC escrowed toward Injective"
	local dydx_before inj_voucher_before
	transfer injective user "$INJECTIVE_CHANNEL" "$DYDX_HOLDER" "50000000$COUNTERPARTY_USDC"
	wait_for_balance dydx "$DYDX_HOLDER" "$INJECTIVE_DENOM" 50000000
	pass "DISABLED Injective deposit leaves the holder with the physical voucher"

	dydx_before=$(balance dydx "$DYDX_USER" "$USDC_DENOM")
	inj_voucher_before=$(balance injective "$INJECTIVE_USER" "$LOGICAL_ON_INJECTIVE")
	transfer dydx user "$DYDX_INJECTIVE_CHANNEL" "$INJECTIVE_USER" "10000000$USDC_DENOM"
	wait_for_balance injective "$INJECTIVE_USER" "$LOGICAL_ON_INJECTIVE" "$(bc_calc "$inj_voucher_before + 10000000")"
	assert_eq "dYdX user USDC" "$(balance dydx "$DYDX_USER" "$USDC_DENOM")" "$(bc_calc "$dydx_before - 10000000")"
	pass "DISABLED logical USDC leaves over Injective as a legacy escrowed voucher"
}

scenario_gradual_legacy_positions_unwind() {
	log "GRADUAL: legacy positions unwind without touching the ledger"
	local ledger reserve inj_before supply dydx_before
	ledger=$(ledger_json)
	reserve=$(balance dydx "$CANONICAL_MODULE" "$INJECTIVE_DENOM")
	inj_before=$(balance injective "$INJECTIVE_USER" "$COUNTERPARTY_USDC")
	transfer dydx holder "$DYDX_INJECTIVE_CHANNEL" "$INJECTIVE_USER" "20000000$INJECTIVE_DENOM"
	wait_for_balance injective "$INJECTIVE_USER" "$COUNTERPARTY_USDC" "$(bc_calc "$inj_before + 20000000")"
	assert_eq "holder voucher" "$(balance dydx "$DYDX_HOLDER" "$INJECTIVE_DENOM")" 30000000
	assert_eq "pending" "$(pending_count)" 0
	assert_eq "module reserve" "$(balance dydx "$CANONICAL_MODULE" "$INJECTIVE_DENOM")" "$reserve"
	assert_eq "ledger" "$(ledger_json)" "$ledger"
	check_backing
	pass "existing Injective voucher holder unwinds to native uusdc"

	supply=$(supply_of dydx "$USDC_DENOM")
	dydx_before=$(balance dydx "$DYDX_USER" "$USDC_DENOM")
	transfer injective user "$INJECTIVE_CHANNEL" "$DYDX_USER" "10000000$LOGICAL_ON_INJECTIVE"
	wait_for_balance dydx "$DYDX_USER" "$USDC_DENOM" "$(bc_calc "$dydx_before + 10000000")"
	assert_eq "logical supply" "$(supply_of dydx "$USDC_DENOM")" "$supply"
	assert_eq "ledger" "$(ledger_json)" "$ledger"
	check_backing
	pass "legacy escrowed logical USDC returns over Injective from escrow"
}

scenario_gradual_noncanonical_channel() {
	log "GRADUAL: logical USDC cannot use a second channel to Injective"
	local out ledger dydx_before alt_denom
	out=$(tx_expect_fail dydx user ibc-transfer transfer transfer "$DYDX_INJECTIVE_ALT_A" "$INJECTIVE_USER" \
		"1000000$USDC_DENOM")
	assert_contains "direct send over $DYDX_INJECTIVE_ALT_A" "$out" "unsupported canonical USDC channel"
	pass "direct MsgTransfer of logical USDC over $DYDX_INJECTIVE_ALT_A rejected"

	tx dydx user authz grant "$DYDX_HOLDER" generic \
		--msg-type=/ibc.applications.transfer.v1.MsgTransfer >/dev/null
	"$DYDX_BIN" tx ibc-transfer transfer transfer "$DYDX_INJECTIVE_ALT_A" "$INJECTIVE_USER" \
		"1000000$USDC_DENOM" --from "$DYDX_USER" --chain-id "$DYDX_CHAIN_ID" --node "$(chain_node dydx)" \
		--generate-only >"$WORK_DIR/authz-transfer.json"
	out=$(tx_expect_fail dydx holder authz exec "$WORK_DIR/authz-transfer.json")
	assert_contains "authz exec over $DYDX_INJECTIVE_ALT_A" "$out" "unsupported canonical USDC channel"
	pass "authz MsgExec of logical USDC over $DYDX_INJECTIVE_ALT_A rejected"

	ledger=$(ledger_json)
	dydx_before=$(balance dydx "$DYDX_USER" "$USDC_DENOM")
	alt_denom=$(ibc_denom "transfer/$DYDX_INJECTIVE_ALT_A/$COUNTERPARTY_USDC")
	transfer injective user "$DYDX_INJECTIVE_ALT_B" "$DYDX_USER" "5000000$COUNTERPARTY_USDC"
	wait_for_balance dydx "$DYDX_USER" "$alt_denom" 5000000
	assert_eq "dYdX user USDC" "$(balance dydx "$DYDX_USER" "$USDC_DENOM")" "$dydx_before"
	assert_eq "ledger" "$(ledger_json)" "$ledger"
	check_backing
	pass "Injective uusdc over $DYDX_INJECTIVE_ALT_A stays a plain voucher, nothing minted"
}

scenario_gradual_noble_via_injective() {
	log "GRADUAL: Noble-origin USDC routed through Injective"
	local noble_on_inj ledger dydx_before multihop
	noble_on_inj=$(ibc_denom "transfer/$INJECTIVE_NOBLE_A/$COUNTERPARTY_USDC")
	transfer noble user "$INJECTIVE_NOBLE_B" "$INJECTIVE_USER" "7000000$COUNTERPARTY_USDC"
	wait_for_balance injective "$INJECTIVE_USER" "$noble_on_inj" 7000000

	ledger=$(ledger_json)
	dydx_before=$(balance dydx "$DYDX_USER" "$USDC_DENOM")
	multihop=$(ibc_denom "transfer/$DYDX_INJECTIVE_CHANNEL/transfer/$INJECTIVE_NOBLE_A/$COUNTERPARTY_USDC")
	transfer injective user "$INJECTIVE_CHANNEL" "$DYDX_USER" "7000000$noble_on_inj"
	wait_for_balance dydx "$DYDX_USER" "$multihop" 7000000
	assert_eq "dYdX user USDC" "$(balance dydx "$DYDX_USER" "$USDC_DENOM")" "$dydx_before"
	assert_eq "module reserve of multihop voucher" "$(balance dydx "$CANONICAL_MODULE" "$multihop")" 0
	assert_eq "ledger" "$(ledger_json)" "$ledger"
	check_backing
	pass "Noble->Injective->dYdX USDC arrives as a multihop voucher, not logical USDC"
}

# Sends from Injective to the dYdX user and asserts the packet is error-acked:
# Injective refunds the sender and nothing changes on dYdX.
expect_injective_deposit_refunded() {
	local amount=$1 inj_before dydx_before ledger
	inj_before=$(balance injective "$INJECTIVE_USER" "$COUNTERPARTY_USDC")
	dydx_before=$(balance dydx "$DYDX_USER" "$USDC_DENOM")
	ledger=$(ledger_json)
	transfer injective user "$INJECTIVE_CHANNEL" "$DYDX_USER" "$amount$COUNTERPARTY_USDC"
	assert_eq "Injective escrowed" "$(balance injective "$INJECTIVE_USER" "$COUNTERPARTY_USDC")" \
		"$(bc_calc "$inj_before - $amount")"
	wait_for_balance injective "$INJECTIVE_USER" "$COUNTERPARTY_USDC" "$inj_before"
	assert_eq "dYdX user USDC" "$(balance dydx "$DYDX_USER" "$USDC_DENOM")" "$dydx_before"
	assert_eq "ledger" "$(ledger_json)" "$ledger"
	check_backing
}

scenario_gradual_limits() {
	log "GRADUAL: per-transfer cap and migration ceiling"
	local over out backing room
	over=$(bc_calc "$MAX_TRANSFER + 1")
	for channel in "$DYDX_INJECTIVE_CHANNEL" "$DYDX_NOBLE_CHANNEL"; do
		out=$(tx_expect_fail dydx user ibc-transfer transfer transfer "$channel" "$INJECTIVE_USER" \
			"$over$USDC_DENOM")
		assert_contains "outbound over cap on $channel" "$out" "invalid canonical USDC amount"
	done
	pass "outbound logical transfers above max_transfer_amount rejected"

	expect_injective_deposit_refunded "$over"
	pass "inbound Injective deposit above max_transfer_amount refunded"

	backing=$(ledger_field injective_backing)
	room=$(bc_calc "$MIGRATION_CEILING - $backing")
	[ "$(bc_calc "$room + 1 <= $MAX_TRANSFER")" = 1 ] || die "ceiling room $room not testable under cap"
	expect_injective_deposit_refunded "$(bc_calc "$room + 1")"
	pass "inbound Injective deposit past migration_ceiling refunded"
}

scenario_gradual_invalid_receiver() {
	log "GRADUAL: Injective withdrawal to an invalid receiver is refunded"
	local dydx_before ledger reserve
	dydx_before=$(balance dydx "$DYDX_USER" "$USDC_DENOM")
	ledger=$(ledger_json)
	reserve=$(balance dydx "$CANONICAL_MODULE" "$INJECTIVE_DENOM")
	transfer dydx user "$DYDX_INJECTIVE_CHANNEL" "not-a-bech32-address" "50000000$USDC_DENOM"
	wait_for_no_pending
	wait_for_balance dydx "$DYDX_USER" "$USDC_DENOM" "$dydx_before"
	assert_eq "ledger" "$(ledger_json)" "$ledger"
	assert_eq "module reserve" "$(balance dydx "$CANONICAL_MODULE" "$INJECTIVE_DENOM")" "$reserve"
	assert_eq "module locked logical" "$(balance dydx "$CANONICAL_MODULE" "$USDC_DENOM")" 0
	check_backing
	pass "error ack refunds logical USDC and restores Injective backing"
}

scenario_gradual_injective_deposits_toggle() {
	log "governance: Injective deposits off refuses inflow, exits keep working, back on resumes"
	local backing inj_before dydx_before
	propose_controls MODE_GRADUAL PASSED '.injective_deposits_disabled = true' >/dev/null
	assert_eq "deposits disabled" "$(state | jq -r '.controls.injective_deposits_disabled')" true
	expect_injective_deposit_refunded 30000000
	pass "Injective deposit refunded on Injective while deposits are disabled"

	backing=$(ledger_field injective_backing)
	inj_before=$(balance injective "$INJECTIVE_USER" "$COUNTERPARTY_USDC")
	transfer dydx user "$DYDX_INJECTIVE_CHANNEL" "$INJECTIVE_USER" "20000000$USDC_DENOM"
	wait_for_balance injective "$INJECTIVE_USER" "$COUNTERPARTY_USDC" "$(bc_calc "$inj_before + 20000000")"
	wait_for_no_pending
	assert_eq "injective_backing" "$(ledger_field injective_backing)" "$(bc_calc "$backing - 20000000")"
	check_backing
	pass "Injective withdrawal still settles while deposits are disabled"

	propose_controls MODE_GRADUAL PASSED >/dev/null
	assert_eq "deposits enabled" "$(state | jq -r '.controls.injective_deposits_disabled')" false
	dydx_before=$(balance dydx "$DYDX_USER" "$USDC_DENOM")
	backing=$(ledger_field injective_backing)
	transfer injective user "$INJECTIVE_CHANNEL" "$DYDX_USER" "20000000$COUNTERPARTY_USDC"
	wait_for_balance dydx "$DYDX_USER" "$USDC_DENOM" "$(bc_calc "$dydx_before + 20000000")"
	assert_eq "injective_backing" "$(ledger_field injective_backing)" "$(bc_calc "$backing + 20000000")"
	check_backing
	pass "Injective deposits credited again after re-enabling"
}

scenario_gradual_relayer_backlog() {
	log "GRADUAL: relayer outage with mixed traffic, then catch-up"
	local dydx_before inj_before noble_before inj_backing noble_backing to_inj to_noble
	dydx_before=$(balance dydx "$DYDX_USER" "$USDC_DENOM")
	inj_before=$(balance injective "$INJECTIVE_USER" "$COUNTERPARTY_USDC")
	noble_before=$(balance noble "$NOBLE_USER" "$COUNTERPARTY_USDC")
	inj_backing=$(ledger_field injective_backing)
	noble_backing=$(ledger_field noble_backing)

	stop_hermes
	to_inj=$(transfer_tx dydx user "$DYDX_INJECTIVE_CHANNEL" "$INJECTIVE_USER" "30000000$USDC_DENOM")
	to_noble=$(transfer_tx dydx user "$DYDX_NOBLE_CHANNEL" "$NOBLE_USER" "20000000$USDC_DENOM")
	transfer injective user "$INJECTIVE_CHANNEL" "$DYDX_USER" "40000000$COUNTERPARTY_USDC"
	assert_eq "pending while relayer down" "$(pending_count)" 2
	assert_eq "pending_injective" "$(ledger_field pending_injective)" 30000000
	assert_eq "pending_noble" "$(ledger_field pending_noble)" 20000000
	check_backing
	start_hermes
	relay_dydx_packet "$to_inj"
	relay_dydx_packet "$to_noble"

	wait_for_balance injective "$INJECTIVE_USER" "$COUNTERPARTY_USDC" "$(bc_calc "$inj_before + 30000000 - 40000000")"
	wait_for_balance noble "$NOBLE_USER" "$COUNTERPARTY_USDC" "$(bc_calc "$noble_before + 20000000")"
	wait_for_balance dydx "$DYDX_USER" "$USDC_DENOM" "$(bc_calc "$dydx_before - 30000000 - 20000000 + 40000000")"
	wait_for_no_pending
	assert_eq "injective_backing" "$(ledger_field injective_backing)" "$(bc_calc "$inj_backing - 30000000 + 40000000")"
	assert_eq "noble_backing" "$(ledger_field noble_backing)" "$(bc_calc "$noble_backing - 20000000")"
	assert_eq "pending_injective" "$(ledger_field pending_injective)" 0
	assert_eq "pending_noble" "$(ledger_field pending_noble)" 0
	check_backing
	pass "backlog relayed after restart with ledger consistent"
}

scenario_governance_mismatched_route_rejected() {
	log "governance: controls bound to the wrong client fail and change nothing"
	local controls inj_client
	controls=$(controls_json)
	inj_client=$(echo "$controls" | jq -r '.injective_client')
	propose_controls MODE_GRADUAL FAILED ".noble_client = \"$inj_client\"" >/dev/null
	assert_eq "controls" "$(controls_json)" "$controls"
	check_backing
	pass "proposal with mismatched Noble client FAILED, controls unchanged"
}

scenario_pending_blocks_disable_then_resume() {
	log "governance: pending settlement blocks DISABLED, PAUSED settles, GRADUAL resumes"
	local inj_before backing dydx_before sent
	inj_before=$(balance injective "$INJECTIVE_USER" "$COUNTERPARTY_USDC")
	backing=$(ledger_field injective_backing)

	stop_hermes
	sent=$(transfer_tx dydx user "$DYDX_INJECTIVE_CHANNEL" "$INJECTIVE_USER" "25000000$USDC_DENOM" \
		--packet-timeout-height 0-0 --packet-timeout-timestamp 3600000000000)
	assert_eq "pending" "$(pending_count)" 1
	propose_controls MODE_DISABLED FAILED >/dev/null
	assert_eq "mode after failed disable" "$(state | jq -r '.controls.mode')" MODE_GRADUAL
	pass "DISABLED rejected while a settlement is pending"

	submit_controls MODE_PAUSED
	start_hermes
	relay_dydx_packet "$sent"
	wait_for_balance injective "$INJECTIVE_USER" "$COUNTERPARTY_USDC" "$(bc_calc "$inj_before + 25000000")"
	wait_for_no_pending
	assert_eq "injective_backing" "$(ledger_field injective_backing)" "$(bc_calc "$backing - 25000000")"
	check_backing
	pass "in-flight settlement completes while PAUSED"

	submit_controls MODE_GRADUAL
	assert_eq "re-snapshotted noble_backing" "$(ledger_field noble_backing)" \
		"$(bc_calc "$(supply_of dydx "$USDC_DENOM") - $(ledger_field injective_backing) - $(ledger_field pending_injective)")"
	check_backing
	dydx_before=$(balance dydx "$DYDX_USER" "$USDC_DENOM")
	transfer injective user "$INJECTIVE_CHANNEL" "$DYDX_USER" "5000000$COUNTERPARTY_USDC"
	wait_for_balance dydx "$DYDX_USER" "$USDC_DENOM" "$(bc_calc "$dydx_before + 5000000")"
	transfer dydx user "$DYDX_INJECTIVE_CHANNEL" "$INJECTIVE_USER" "5000000$USDC_DENOM"
	wait_for_no_pending
	wait_for_balance dydx "$DYDX_USER" "$USDC_DENOM" "$dydx_before"
	check_backing
	pass "GRADUAL resumes after PAUSED and routes both directions"
}

scenario_paused_rejects_transfers() {
	log "governance: MODE_PAUSED rejects logical transfers"
	local out
	submit_controls MODE_PAUSED
	out=$(tx_expect_fail dydx user ibc-transfer transfer transfer "$DYDX_INJECTIVE_CHANNEL" "$INJECTIVE_USER" \
		"1000000$USDC_DENOM")
	assert_contains "PAUSED transfer" "$out" "paused"
	check_backing
	pass "PAUSED blocks logical USDC transfers"
}

# Ceiling used while draining: every Noble-backed coin becomes Injective-backed.
MIGRATION_FINAL_CEILING=100000000000

# One replacement batch: Injective USDC in to the migrator, the same logical
# amount out over Noble, settled.
migration_batch() {
	local amount=$1 dydx_before noble_before
	dydx_before=$(balance dydx "$DYDX_USER" "$USDC_DENOM")
	noble_before=$(balance noble "$NOBLE_USER" "$COUNTERPARTY_USDC")
	transfer injective user "$INJECTIVE_CHANNEL" "$DYDX_USER" "$amount$COUNTERPARTY_USDC"
	wait_for_balance dydx "$DYDX_USER" "$USDC_DENOM" "$(bc_calc "$dydx_before + $amount")"
	transfer dydx user "$DYDX_NOBLE_CHANNEL" "$NOBLE_USER" "$amount$USDC_DENOM"
	wait_for_balance noble "$NOBLE_USER" "$COUNTERPARTY_USDC" "$(bc_calc "$noble_before + $amount")"
	wait_for_no_pending
	check_backing
}

scenario_migration_finalize() {
	log "migration: drain Noble backing into Injective backing, then disable Noble withdrawals"
	local ceiling_filter=".migration_ceiling = \"$MIGRATION_FINAL_CEILING\""
	local i noble batch supply dydx_before noble_before out
	propose_controls MODE_GRADUAL PASSED "$ceiling_filter" >/dev/null
	assert_eq "mode" "$(state | jq -r '.controls.mode')" MODE_GRADUAL
	propose_controls MODE_GRADUAL FAILED "$ceiling_filter | .noble_withdrawals_enabled = false" >/dev/null
	assert_eq "Noble withdrawals still enabled" "$(state | jq -r '.controls.noble_withdrawals_enabled')" true
	pass "Noble withdrawals cannot be disabled while Noble backing remains"

	for i in $(seq 1 20); do
		noble=$(ledger_field noble_backing)
		[ "$noble" = "0" ] && break
		batch=$noble
		[ "$(bc_calc "$batch > $MAX_TRANSFER")" = 1 ] && batch=$MAX_TRANSFER
		log "  batch $i: $batch (noble_backing $noble)"
		migration_batch "$batch"
	done
	assert_eq "noble_backing drained" "$(ledger_field noble_backing)" 0
	supply=$(supply_of dydx "$USDC_DENOM")
	assert_eq "supply fully Injective-backed" "$(ledger_field injective_backing)" "$supply"
	pass "Noble backing replaced by Injective backing"

	assert_eq "no legacy Noble packets in flight" "$(legacy_noble_in_flight)" 0
	propose_controls MODE_GRADUAL PASSED "$ceiling_filter | .noble_withdrawals_enabled = false" >/dev/null
	assert_eq "Noble withdrawals disabled" "$(state | jq -r '.controls.noble_withdrawals_enabled')" false
	check_backing

	out=$(tx_expect_fail dydx user ibc-transfer transfer transfer "$DYDX_NOBLE_CHANNEL" "$NOBLE_USER" \
		"1000000$USDC_DENOM")
	assert_contains "Noble withdrawal after finalization" "$out" "routing is disabled"
	noble_before=$(balance noble "$NOBLE_USER" "$COUNTERPARTY_USDC")
	dydx_before=$(balance dydx "$DYDX_USER" "$USDC_DENOM")
	transfer noble user "$NOBLE_CHANNEL" "$DYDX_USER" "1000000$COUNTERPARTY_USDC"
	wait_for_balance noble "$NOBLE_USER" "$COUNTERPARTY_USDC" "$noble_before"
	assert_eq "no mint from Noble" "$(balance dydx "$DYDX_USER" "$USDC_DENOM")" "$dydx_before"
	transfer dydx user "$DYDX_INJECTIVE_CHANNEL" "$INJECTIVE_USER" "1000000$USDC_DENOM"
	wait_for_no_pending
	assert_eq "exit via Injective burns backing" "$(ledger_field injective_backing)" "$(bc_calc "$supply - 1000000")"
	check_backing
	pass "after finalization USDC enters and exits only via Injective"
}

# Opt-in check that the dYdX node can serve tx_search, which relayers rely on to
# clear packets after downtime. Fails against the current cometbft fork; see
# relay_dydx_packet in lib.sh.
scenario_dydx_tx_index() {
	log "dYdX tx index: tx_search returns every matching MsgTransfer"
	local count
	count=$(curl -sf -G "http://127.0.0.1:$DYDX_RPC_PORT/tx_search" \
		--data-urlencode "query=\"message.action='/ibc.applications.transfer.v1.MsgTransfer'\"" |
		jq -r '.result.total_count')
	[ "$count" -gt 1 ] || die "tx_search found $count MsgTransfer txs after several transfers"
	pass "tx_search returned $count MsgTransfer txs"
}

# Not part of a default run; select by name.
DIAGNOSTICS="dydx_tx_index enable_gradual"

SCENARIOS="disabled_noble_roundtrip disabled_legacy_positions disabled_noble_refunded_after_gradual
gradual_noble_deposit_rejected gradual_injective_roundtrip gradual_noble_withdrawal gradual_noble_timeout
gradual_noble_invalid_receiver gradual_injective_timeout
gradual_legacy_positions_unwind gradual_noncanonical_channel gradual_noble_via_injective gradual_limits gradual_invalid_receiver
gradual_injective_deposits_toggle gradual_relayer_backlog governance_mismatched_route_rejected pending_blocks_disable_then_resume
paused_rejects_transfers migration_finalize"

# With arguments, runs only the named scenarios (in the given order), e.g. to
# resume after a failure without restarting the chains.
selected="${*:-$SCENARIOS}"
for name in $selected; do
	case " $(echo $SCENARIOS $DIAGNOSTICS) " in
	*" $name "*) "scenario_$name" ;;
	*) die "unknown scenario '$name'; available: $(echo $SCENARIOS $DIAGNOSTICS)" ;;
	esac
done
log "scenarios passed: $(echo $selected)"
