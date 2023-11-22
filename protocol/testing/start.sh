#!/bin/bash
set -eo pipefail

# For legacy reasons, our internal dev environment runs `/dydxprotocol/start.sh` as the entrypoint for
# cosmovisor images. This file serves a stub we can provide in those images.

if [[ ! -z "${DYDXPROTOCOLD_PATH}" ]]; then
    rm "$DAEMON_HOME/cosmovisor/genesis/bin/dydxprotocold"
	ln -s $DYDXPROTOCOLD_PATH "$DAEMON_HOME/cosmovisor/genesis/bin/dydxprotocold"
fi

cosmovisor "$@"
