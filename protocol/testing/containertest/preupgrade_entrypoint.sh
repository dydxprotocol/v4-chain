#!/bin/bash
set -eo pipefail

# This entrypoint configures the chain to start in a preupgrade state before forwarding args to cosmovisor.
# The following need to be done:
#   - Set the preupgrade binary (downloaded in containertest.sh) to be the genesis binary
#   - Set the binary at the current commit to  be the ugprade binary
#   - Set the genesis to be the preupgrade genesis

if [[ -z "${UPGRADE_TO_VERSION}" ]]; then
    echo >&2 "UPGRADE_TO_VERSION must be set"
    exit 1
fi

MONIKERS=(
	"alice"
	"bob"
	"carl"
	"dave"
)

for i in "${!MONIKERS[@]}"; do
	DAEMON_NAME="dydxprotocold"
	DAEMON_HOME="$HOME/chain/.${MONIKERS[$i]}"

    rm "$DAEMON_HOME/cosmovisor/genesis/bin/dydxprotocold"
    ln -s /bin/dydxprotocold_preupgrade "$DAEMON_HOME/cosmovisor/genesis/bin/dydxprotocold"
    mkdir -p "$DAEMON_HOME/cosmovisor/upgrades/$UPGRADE_TO_VERSION/bin/"
    ln -s /bin/dydxprotocold "$DAEMON_HOME/cosmovisor/upgrades/$UPGRADE_TO_VERSION/bin/dydxprotocold"

    rm "$DAEMON_HOME/config/genesis.json"
    cp "$HOME/preupgrade_genesis.json" "$DAEMON_HOME/config/genesis.json"
done

cosmovisor run "$@"
