package types

import (
	epochstypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/epochs/types"
	pricestypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type PricesKeeper interface {
	GetMarketPrice(
		ctx sdk.Context,
		id uint32,
	) (marketPrice pricestypes.MarketPrice, err error)
	GetMarketIdToValidDaemonPrice(
		ctx sdk.Context,
	) map[uint32]pricestypes.MarketSpotPrice
}

// PerpetualsClobKeeper defines the expected interface for the clob keeper.
type PerpetualsClobKeeper interface {
	GetPricePremiumForPerpetual(
		ctx sdk.Context,
		perpetualId uint32,
		params GetPricePremiumParams,
	) (
		premiumPpm int32,
		err error,
	)
	IsPerpetualClobPairActive(
		ctx sdk.Context,
		perpetualId uint32,
	) (
		isActive bool,
		err error,
	)
}

// EpochsKeeper defines the expected epochs keeper to get epoch info.
type EpochsKeeper interface {
	NumBlocksSinceEpochStart(
		ctx sdk.Context,
		id epochstypes.EpochInfoName,
	) (uint32, error)
	MustGetFundingTickEpochInfo(
		ctx sdk.Context,
	) epochstypes.EpochInfo
	MustGetFundingSampleEpochInfo(
		ctx sdk.Context,
	) epochstypes.EpochInfo
}
