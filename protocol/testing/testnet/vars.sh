#!/bin/bash
set -eo pipefail

source "./version.sh"

# Full node home directories will be set up for indices 0 to LAST_FULL_NODE_INDEX
LAST_FULL_NODE_INDEX=2

# Define monikers for each validator. These are made up strings and can be anything.
# This also controls in which directory the validator's home will be located. i.e. `/dydxprotocol/chain/.alice`
MONIKERS=(
	"dydx-1"
	"dydx-2"
)