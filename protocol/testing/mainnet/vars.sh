#!/bin/bash
set -eo pipefail

source "./version.sh"

# Full node home directories will be set up for indices 0 to LAST_FULL_NODE_INDEX
LAST_FULL_NODE_INDEX=5
