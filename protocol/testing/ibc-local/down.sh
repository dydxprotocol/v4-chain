#!/bin/bash
# Stops the relayer and all three chains. Pass --clean to also delete chain
# state, keys and logs (binaries are kept).
set -euo pipefail
source "$(dirname "${BASH_SOURCE[0]}")/lib.sh"

stop_hermes
for chain in dydx noble injective; do
	stop_process "$chain"
done

if [ "${1:-}" = "--clean" ]; then
	command rm -rf "$DYDX_HOME" "$NOBLE_HOME" "$INJECTIVE_HOME" "$HERMES_HOME" "$KEY_DIR" "$LOG_DIR" "$PID_DIR" "$CHANNELS_ENV"
	log "removed chain state under $STATE_DIR"
fi
