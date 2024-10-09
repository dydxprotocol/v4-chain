package keeper

import (
	veaggregator "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/aggregator"
	vecodec "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/codec"
	abcicomet "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) SetNextBlocksPricesAndSDAIRateFromExtendedCommitInfo(ctx sdk.Context, extendedCommitInfo *abcicomet.ExtendedCommitInfo) error {
	// from cometbft so is either nil or is valid and > 2/3
	if extendedCommitInfo == nil {
		return nil
	}

	veCodec := vecodec.NewDefaultVoteExtensionCodec()
	votes, err := veaggregator.FetchVotesFromExtCommitInfo(*extendedCommitInfo, veCodec)
	if err != nil {
		return err
	}

	if len(votes) == 0 {
		return nil
	}

	prices, conversionRate, err := k.VEApplier.VoteAggregator().AggregateDaemonVEIntoFinalPricesAndConversionRate(ctx, votes)

	// Handle aggregation errors gracefully
	if err != nil {
		return nil
	}

	dummyBytes := []byte{}

	err = k.VEApplier.WritePricesToStoreAndMaybeCache(ctx, prices, dummyBytes, false)
	if err != nil {
		return err
	}

	err = k.VEApplier.WriteSDaiConversionRateToStoreAndMaybeCache(ctx, conversionRate, dummyBytes, false)
	if err != nil {
		return err
	}

	return nil
}
