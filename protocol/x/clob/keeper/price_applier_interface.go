package keeper

import (
	"math/big"

	aggregator "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/aggregator"
	voteweighted "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type VEApplierInterface interface {
	WritePricesToStoreAndMaybeCache(ctx sdk.Context, prices map[string]voteweighted.AggregatorPricePair, round int32, writeToCache bool) error
	WriteSDaiConversionRateToStoreAndMaybeCache(ctx sdk.Context, sDaiConversionRate *big.Int, round int32, writeToCache bool) error
	VoteAggregator() aggregator.VoteAggregator
}
