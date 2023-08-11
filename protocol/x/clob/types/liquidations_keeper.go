package types

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	satypes "github.com/dydxprotocol/v4/x/subaccounts/types"
)

// LiquidationsKeeper is an interface that encapsulates all reads and writes to the
// in-memory data structures that store liquidation information.
type LiquidationsKeeper interface {
	DeliverTxPlacePerpetualLiquidation(
		ctx sdk.Context,
		liquidationOrder LiquidationOrder,
		memclob MemClob,
	) (
		orderSizeOptimisticallyFilledFromMatchingQuantums satypes.BaseQuantums,
		orderStatus OrderStatus,
		err error,
	)
	CheckTxPlacePerpetualLiquidation(
		ctx sdk.Context,
		liquidationOrder LiquidationOrder,
	) (
		orderSizeOptimisticallyFilledFromMatchingQuantums satypes.BaseQuantums,
		orderStatus OrderStatus,
		err error,
	)
	IsLiquidatable(
		ctx sdk.Context,
		subaccountId satypes.SubaccountId,
	) (
		bool,
		error,
	)
	GetBankruptcyPriceInQuoteQuantums(
		ctx sdk.Context,
		subaccountId satypes.SubaccountId,
		perpetualId uint32,
		deltaQuantums *big.Int,
	) (
		quoteQuantums *big.Int,
		err error,
	)
	GetFillablePrice(
		ctx sdk.Context,
		subaccountId satypes.SubaccountId,
		perpetualId uint32,
		deltaQuantums *big.Int,
	) (
		fillablePrice *big.Rat,
		err error,
	)
	GetInsuranceFundBalance(
		ctx sdk.Context,
	) (
		balance uint64,
	)
	GetLiquidationInsuranceFundDelta(
		ctx sdk.Context,
		subaccountId satypes.SubaccountId,
		perpetualId uint32,
		isBuy bool,
		fillAmount uint64,
		subticks Subticks,
	) (
		insuranceFundDeltaQuoteQuantums *big.Int,
		err error,
	)
	ConvertFillablePriceToSubticks(
		ctx sdk.Context,
		fillablePrice *big.Rat,
		isLiquidatingLong bool,
		clobPair ClobPair,
	) (
		subticks Subticks,
	)
	GetPerpetualPositionToLiquidate(
		ctx sdk.Context,
		subaccountId satypes.SubaccountId,
	) (
		clobPair ClobPair,
		quantums *big.Int,
		err error,
	)
	GetMaxLiquidatableNotionalAndInsuranceLost(
		ctx sdk.Context,
		subaccountId satypes.SubaccountId,
		perpetualId uint32,
	) (
		bigMaxNotionalLiquidatable *big.Int,
		bigMaxQuantumsInsuranceLost *big.Int,
		err error,
	)
	GetMaxAndMinPositionNotionalLiquidatable(
		ctx sdk.Context,
		positionToLiquidate *satypes.PerpetualPosition,
	) (
		bigMinNotionalLiquidatable *big.Int,
		bigMaxNotionalLiquidatable *big.Int,
		err error,
	)
	MaybeLiquidateSubaccount(
		ctx sdk.Context,
		subaccountId satypes.SubaccountId,
	) (
		err error,
	)
	GetSubaccountLiquidationInfo(
		ctx sdk.Context,
		subaccountId satypes.SubaccountId,
	) (
		liquidationInfo SubaccountLiquidationInfo,
	)
	MustUpdateSubaccountPerpetualLiquidated(
		ctx sdk.Context,
		subaccountId satypes.SubaccountId,
		perpetualId uint32,
	)
	UpdateSubaccountLiquidationInfo(
		ctx sdk.Context,
		subaccountId satypes.SubaccountId,
		notionalLiquidatedQuoteQuantums *big.Int,
		insuranceFundDeltaQuoteQuantums *big.Int,
	)
}
