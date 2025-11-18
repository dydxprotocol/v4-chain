#!/bin/bash
# Vote on a governance proposal in testnet
# Usage: ./vote_in_testnet.sh <proposal_id> <vote_option>
#   $1 = proposal ID (e.g., 5)
#   $2 = vote option (yes, no, abstain, no_with_veto)
# Example: ./vote_in_testnet.sh 5 yes

dydxprotocold tx gov vote "$1" $2 --from dydx-1 --broadcast-mode sync --fees 5000000000000000adv4tnt --yes --keyring-backend test
dydxprotocold tx gov vote "$1" $2 --from dydx-2 --broadcast-mode sync --fees 5000000000000000adv4tnt --yes --keyring-backend test
