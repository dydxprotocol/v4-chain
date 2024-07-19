package ve

import (
	"fmt"

	constants "github.com/StreamFinance-Protocol/stream-chain/protocol/app/constants"
	codec "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/codec"
	vetypes "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/types"
	veutils "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/utils"
	priceskeeper "github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices/keeper"
	cometabci "github.com/cometbft/cometbft/abci/types"
	cometproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type ValidateVoteExtensionsFn func(ctx sdk.Context, extCommitInfo cometabci.ExtendedCommitInfo) error

func CleanAndValidateExtCommitInfo(
	ctx sdk.Context,
	extCommitInfo cometabci.ExtendedCommitInfo,
	veCodec codec.VoteExtensionCodec,
	pricesKeeper PreparePricesKeeper,
	validateVoteExtensionsFn ValidateVoteExtensionsFn,
) (cometabci.ExtendedCommitInfo, error) {
	for i, vote := range extCommitInfo.Votes {
		if err := validateVoteExtension(ctx, vote, veCodec, pricesKeeper); err != nil {
			ctx.Logger().Info(
				"failed to validate vote extension - pruning vote",
				"err", err,
				"validator", vote.Validator.Address,
			)

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
	pricesKeeper PreparePricesKeeper,
) error {
	if vote.VoteExtension == nil && vote.ExtensionSignature == nil {
		return nil
	}

	voteExt, err := voteExtensionCodec.Decode(vote.VoteExtension)
	if err != nil {
		return err
	}

	// The vote extensions are from the previous block.
	if err := ValidateDaemonVoteExtension(ctx, voteExt, pricesKeeper); err != nil {
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

func ValidateDaemonVoteExtension(
	ctx sdk.Context,
	ve vetypes.DaemonVoteExtension,
	pricesKeeper PreparePricesKeeper,
) error {
	if err := ValidateMarketCountInVE(ctx, ve, pricesKeeper); err != nil {
		return err
	}

	if err := ValidatePricesBytesSizeInVE(ctx, ve); err != nil {
		return err
	}

	return nil
}

func ValidateMarketCountInVE(
	ctx sdk.Context,
	ve vetypes.DaemonVoteExtension,
	pricesKeeper PreparePricesKeeper,
) error {
	maxPairs := veutils.GetMaxMarketPairs(ctx, pricesKeeper.(priceskeeper.Keeper))
	if uint32(len(ve.Prices)) > maxPairs {
		return fmt.Errorf(
			"number of oracle vote extension pairs of %d greater than maximum expected pairs of %d",
			uint64(len(ve.Prices)),
			uint64(maxPairs),
		)
	}
	return nil
}

func ValidatePricesBytesSizeInVE(
	ctx sdk.Context,
	ve vetypes.DaemonVoteExtension,
) error {
	for _, priceBytes := range ve.Prices {

		// Ensure that the price bytes are not too long.
		if len(priceBytes) > constants.MaximumPriceSizeInBytes {
			return fmt.Errorf("price bytes are too long: %d", len(priceBytes))
		}
	}
	return nil
}
