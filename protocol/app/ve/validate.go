package ve

import (
	"fmt"

	codec "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/codec"
	cometabci "github.com/cometbft/cometbft/abci/types"
	cometproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func PruneAndValidateExtendedCommitInfo(
	ctx sdk.Context,
	extCommitInfo cometabci.ExtendedCommitInfo,
	veCodec codec.VoteExtensionCodec,
	pk PreparePricesKeeper,
	validateVoteExtensionsFn func(ctx sdk.Context, extCommitInfo cometabci.ExtendedCommitInfo) error,
) (cometabci.ExtendedCommitInfo, error) {
	for i, vote := range extCommitInfo.Votes {
		if err := validateVoteExtension(ctx, vote, veCodec, pk); err != nil {
			ctx.Logger().Info(
				"failed to validate vote extension - pruning vote",
				"err", err,
				"validator", vote.Validator.Address,
			)
			fmt.Println("failed to validate vote extension - pruning vote", err, vote.Validator.Address)

			// failed to validate this vote-extension, mark it as absent in the original commit
			vote.BlockIdFlag = cometproto.BlockIDFlagAbsent
			vote.ExtensionSignature = nil
			vote.VoteExtension = nil
			extCommitInfo.Votes[i] = vote
		}
	}

	// validate after pruning
	if err := validateVoteExtensionsFn(ctx, extCommitInfo); err != nil {
		ctx.Logger().Error(
			"failed to validate vote extensions; vote extensions may not comprise a super-majority",
			"err", err,
		)

		return cometabci.ExtendedCommitInfo{}, err
	}

	return extCommitInfo, nil
}

func validateVoteExtension(
	ctx sdk.Context,
	vote cometabci.ExtendedVoteInfo,
	voteExtensionCodec codec.VoteExtensionCodec,
	pk PreparePricesKeeper,
) error {
	// vote is not voted for if VE is nil
	if vote.VoteExtension == nil && vote.ExtensionSignature == nil {
		return nil
	}

	voteExt, err := voteExtensionCodec.Decode(vote.VoteExtension)
	if err != nil {
		return err
	}

	// The vote extensions are from the previous block.
	if err := ValidateDaemonVoteExtension(ctx, voteExt, pk); err != nil {
		return err
	}

	return nil
}

func ValidateExtendedCommitInfo(
	ctx sdk.Context,
	height int64,
	extCommitInfo cometabci.ExtendedCommitInfo,
	veCodec codec.VoteExtensionCodec,
	pk PreparePricesKeeper,
	validateVoteExtensionFn func(ctx sdk.Context, extCommitInfo cometabci.ExtendedCommitInfo) error,
) error {
	if err := validateVoteExtensionFn(ctx, extCommitInfo); err != nil {
		ctx.Logger().Error(
			"failed to validate vote extension",
			"height", height,
			"err", err,
		)
		return err
	}

	for _, vote := range extCommitInfo.Votes {
		addr := sdk.ConsAddress(vote.Validator.Address)

		if err := validateVoteExtension(ctx, vote, veCodec, pk); err != nil {
			ctx.Logger().Error(
				"failed to validate vote extension",
				"height", height,
				"validator", addr,
				"err", err,
			)
			return err
		}
	}
	return nil
}
