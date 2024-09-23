package keeper

import (
	veaggregator "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/aggregator"
	vecodec "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/codec"
	abcicomet "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) SetNextBlocksPricesFromExtendedCommitInfo(ctx sdk.Context, extendedCommitInfo *abcicomet.ExtendedCommitInfo) error {

	// from cometbft so is either nil or is valid and > 2/3
	if extendedCommitInfo != nil {
		veCodec := vecodec.NewDefaultVoteExtensionCodec()
		votes, err := veaggregator.FetchVotesFromExtCommitInfo(*extendedCommitInfo, veCodec)
		if err != nil {
			return err
		}

		if len(votes) > 0 {
			prices, err := k.PriceApplier.VoteAggregator().AggregateDaemonVEIntoFinalPrices(ctx, votes)
			if err == nil {
				err = k.PriceApplier.WritePricesToStoreAndMaybeCache(ctx, prices, 0, false)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
