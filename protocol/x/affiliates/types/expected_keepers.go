package types

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	revsharetypes "github.com/dydxprotocol/v4-chain/protocol/x/revshare/types"
	stattypes "github.com/dydxprotocol/v4-chain/protocol/x/stats/types"
)

type StatsKeeper interface {
	GetStakedAmount(ctx sdk.Context, delegatorAddr string) *big.Int
	GetBlockStats(ctx sdk.Context) *stattypes.BlockStats
}

type RevShareKeeper interface {
	GetUnconditionalRevShareConfigParams(ctx sdk.Context) (revsharetypes.UnconditionalRevShareConfig, error)
	GetMarketMapperRevenueShareParams(
		ctx sdk.Context,
	) revsharetypes.MarketMapperRevenueShareParams
	ValidateRevShareSafety(
		ctx sdk.Context,
		affiliateTiers AffiliateTiers,
		unconditionalRevShareConfig revsharetypes.UnconditionalRevShareConfig,
		marketMapperRevShareParams revsharetypes.MarketMapperRevenueShareParams,
	) bool
}
