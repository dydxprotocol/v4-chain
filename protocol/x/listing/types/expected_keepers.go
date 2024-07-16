package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	marketmaptypes "github.com/skip-mev/slinky/x/marketmap/types"
)

type PricesKeeper interface {
	CreateMarket(
		ctx sdk.Context,
		marketParam pricestypes.MarketParam,
		marketPrice pricestypes.MarketPrice,
	) (pricestypes.MarketParam, error)
	AcquireNextMarketID(ctx sdk.Context) uint32
}

type ClobKeeper interface {
	CreatePerpetualClobPair(
		ctx sdk.Context,
		clobPairId uint32,
		perpetualId uint32,
		stepSizeBaseQuantums satypes.BaseQuantums,
		quantumConversionExponent int32,
		subticksPerTick uint32,
		status clobtypes.ClobPair_Status,
	) (clobtypes.ClobPair, error)
	AcquireNextClobPairID(ctx sdk.Context) uint32
}

type MarketMapKeeper interface {
	GetMarket(
		ctx sdk.Context,
		ticker string,
	) (marketmaptypes.Market, error)
}
