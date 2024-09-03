package types

import (
	"math/big"

	satypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/subaccounts/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
	GetInsuranceFundBalanceInQuoteQuantums(
		ctx sdk.Context,
		perpetualId uint32,
	) (
		balance *big.Int,
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
	ConvertLiquidationPriceToSubticks(
		ctx sdk.Context,
		liquidationPrice *big.Rat,
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
	GetSubaccountMaxInsuranceLost(
		ctx sdk.Context,
		subaccountId satypes.SubaccountId,
		perpetualId uint32,
	) (
		bigMaxQuantumsInsuranceLost *big.Int,
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
		insuranceFundDeltaQuoteQuantums *big.Int,
	)
}
