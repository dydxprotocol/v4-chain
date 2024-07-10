package types

import sdk "github.com/cosmos/cosmos-sdk/types"

type RevShareKeeper interface {
	// MarketMapperRevenueShareParams
	SetMarketMapperRevenueShareParams(
		ctx sdk.Context,
		params MarketMapperRevenueShareParams,
	) (err error)
	GetMarketMapperRevenueShareParams(
		ctx sdk.Context,
	) (params MarketMapperRevenueShareParams)

	// MarketMapperRevShareDetails
	SetMarketMapperRevShareDetails(
		ctx sdk.Context,
		marketId uint32,
		params MarketMapperRevShareDetails,
	)
	GetMarketMapperRevShareDetails(
		ctx sdk.Context,
		marketId uint32,
	) (params MarketMapperRevShareDetails, err error)
	CreateNewMarketRevShare(ctx sdk.Context, marketId uint32)
}
