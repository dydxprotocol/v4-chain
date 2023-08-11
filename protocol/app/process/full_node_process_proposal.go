package process

import (
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// FullNodeProcessProposalHandler is the `ProcessProposal` implementation for full-nodes.
// This implementation always returns `abci.ResponseProcessProposal_ACCEPT`.
// Validators within the validator set should never use this implementation.
func FullNodeProcessProposalHandler() sdk.ProcessProposalHandler {
	return func(ctx sdk.Context, req abci.RequestProcessProposal) abci.ResponseProcessProposal {
		return abci.ResponseProcessProposal{Status: abci.ResponseProcessProposal_ACCEPT}
	}
}
