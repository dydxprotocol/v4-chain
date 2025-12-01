#!/bin/bash
# Vote on a governance proposal in staging
# Usage: ./vote_in_staging.sh <proposal_id> <vote_option>
#   $1 = proposal ID (e.g., 5)
#   $2 = vote option (yes, no, abstain, no_with_veto)
# Example: ./vote_in_staging.sh 5 yes

dydxprotocold tx gov vote "$1" $2 --from alice --broadcast-mode sync --fees 5000000000000000adv4tnt --yes --keyring-backend test
dydxprotocold tx gov vote "$1" $2 --from bob --broadcast-mode sync --fees 5000000000000000adv4tnt --yes --keyring-backend test
dydxprotocold tx gov vote "$1" $2 --from carl --broadcast-mode sync --fees 5000000000000000adv4tnt --yes --keyring-backend test
dydxprotocold tx gov vote "$1" $2 --from dave --broadcast-mode sync --fees 5000000000000000adv4tnt --yes --keyring-backend test
dydxprotocold tx gov vote "$1" $2 --from emily --broadcast-mode sync --fees 5000000000000000adv4tnt --yes --keyring-backend test
dydxprotocold tx gov vote "$1" $2 --from fiona --broadcast-mode sync --fees 5000000000000000adv4tnt --yes --keyring-backend test
dydxprotocold tx gov vote "$1" $2 --from greg --broadcast-mode sync --fees 5000000000000000adv4tnt --yes --keyring-backend test
dydxprotocold tx gov vote "$1" $2 --from henry --broadcast-mode sync --fees 5000000000000000adv4tnt --yes --keyring-backend test
dydxprotocold tx gov vote "$1" $2 --from ian --broadcast-mode sync --fees 5000000000000000adv4tnt --yes --keyring-backend test
dydxprotocold tx gov vote "$1" $2 --from jeff --broadcast-mode sync --fees 5000000000000000adv4tnt --yes --keyring-backend test
