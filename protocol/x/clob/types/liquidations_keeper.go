package types

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

// LiquidationsKeeper is an interface that encapsulates all reads and writes to the
// in-memory data structures that store liquidation information.
type LiquidationsKeeper interface {
	PlacePerpetualLiquidation(
		ctx sdk.Context,
		liquidationOrder LiquidationOrder,
	) (
		orderSizeOptimisticallyFilledFromMatchingQuantums satypes.BaseQuantums,
		orderStatus OrderStatus,
		err error,
	)
	MaybeDeleverageSubaccount(
		ctx sdk.Context,
		subaccountId satypes.SubaccountId,
		perpetualId uint32,
	) (
		quantumsDeleveraged *big.Int,
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
		perpetualId uint32,
		err error,
	)
	GetSubaccountMaxNotionalLiquidatable(
		ctx sdk.Context,
		subaccountId satypes.SubaccountId,
		perpetualId uint32,
	) (
		bigMaxNotionalLiquidatable *big.Int,
		err error,
	)
	GetSubaccountMaxInsuranceLost(
		ctx sdk.Context,
		subaccountId satypes.SubaccountId,
		perpetualId uint32,
	) (
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
	MaybeGetLiquidationOrder(
		ctx sdk.Context,
		subaccountId satypes.SubaccountId,
	) (
		liquidationOrder *LiquidationOrder,
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
