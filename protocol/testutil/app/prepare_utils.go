package app

import (
	cometabci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NoOpValidateVoteExtensionsFn(
	_ sdk.Context,
	_ cometabci.ExtendedCommitInfo,
) error {
	return nil
}
