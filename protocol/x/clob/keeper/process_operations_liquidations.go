package keeper

import (
	"math/big"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

func (k Keeper) validateMatchedLiquidation(
	ctx sdk.Context,
	order types.MatchableOrder,
	perpetualId uint32,
	fillAmount satypes.BaseQuantums,
	makerSubticks types.Subticks,
	bigFillQuoteQuantums *big.Int,
) (
	insuranceFundDelta *big.Int,
	err error,
) {
	if !order.IsLiquidation() {
		panic("Expected validateMatchedLiquidation to be called with a liquidation order")
	}

	// Calculate the insurance fund delta for this fill.
	liquidatedSubaccountId := order.GetSubaccountId()
	insuranceFundDelta, err = k.GetLiquidationInsuranceFundDelta(
		ctx,
		liquidatedSubaccountId,
		perpetualId,
		order.IsBuy(),
		fillAmount.ToUint64(),
		makerSubticks,
	)
	if err != nil {
		return nil, err
	}

	sign := metrics.Positive
	if insuranceFundDelta.Sign() == -1 {
		sign = metrics.Negative
	}

	// Only increment this counter during `DeliverTx`.
	if !ctx.IsCheckTx() && !ctx.IsReCheckTx() {
		telemetry.IncrCounter(1, metrics.Liquidations, metrics.InsuranceFundDelta, sign)
	}

	// Validate that processing the liquidation fill does not leave insufficient funds
	// in the insurance fund (such that deleveraging is required and the liquidation couldn't
	// have possibly continued).
	if k.ShouldPerformDeleveraging(ctx, insuranceFundDelta) {
		ctx.Logger().Info("ProcessMatches: insurance fund does not have enough balance and deleveraging is required.")
		return nil, sdkerrors.Wrapf(
			types.ErrInsuranceFundHasInsufficientFunds,
			"Liquidation order %v",
			order,
		)
	}

	// Validate that total notional liquidated and total insurance funds lost do not exceed subaccount block limits.
	if err := k.validateMatchPerpetualLiquidationAgainstSubaccountBlockLimits(
		ctx,
		liquidatedSubaccountId,
		perpetualId,
		bigFillQuoteQuantums,
		insuranceFundDelta,
	); err != nil {
		return nil, err
	}

	return insuranceFundDelta, nil
}

// validateMatchPerpetualLiquidationAgainstSubaccountBlockLimits performs stateful validation
// against the subaccount block limits specified in liquidation configs.
// If validation fails, an error is returned.
//
// The following validation occurs in this method:
//   - The subaccount and perpetual ID pair has not been previously liquidated in the same block.
//   - The total notional liquidated does not exceed the maximum notional amount that a single subaccount
//     can have liquidated per block.
//   - The total insurance lost does not exceed the maximum insurance lost per block.
func (k Keeper) validateMatchPerpetualLiquidationAgainstSubaccountBlockLimits(
	ctx sdk.Context,
	subaccountId satypes.SubaccountId,
	perpetualId uint32,
	bigNotionalLiquidated *big.Int,
	insuranceFundDeltaQuoteQuantums *big.Int,
) (
	err error,
) {
	// Get the maximum liquidatable notional and insurance lost for the liquidated subaccount.
	bigMaxNotionalLiquidatable, bigMaxQuantumsInsuranceLost, err := k.GetMaxLiquidatableNotionalAndInsuranceLost(
		ctx,
		subaccountId,
		perpetualId,
	)
	if err != nil {
		return err
	}

	// Validate that this liquidation does not exceed the maximum notional amount that a single subaccount can have
	// liquidated per block.
	if bigNotionalLiquidated.CmpAbs(bigMaxNotionalLiquidatable) > 0 {
		return types.ErrLiquidationExceedsSubaccountMaxNotionalLiquidated
	}

	// Validate that this liquidation does not exceed the maximum insurance fund payout amount for this
	// subaccount per block.
	if insuranceFundDeltaQuoteQuantums.Sign() == -1 {
		bigAbsInsuranceFundDelta := new(big.Int).Abs(insuranceFundDeltaQuoteQuantums)
		if bigAbsInsuranceFundDelta.Cmp(bigMaxQuantumsInsuranceLost) > 0 {
			return types.ErrLiquidationExceedsSubaccountMaxInsuranceLost
		}
	}
	return nil
}

// ConstructTakerOrderFromMatchPerpetualLiquidation creates and returns the corresponding LiquidationOrder
// for the given match.
// An error is returned if:
//   - The clob pair is invalid or does not match the provided perpetual id.
//   - `GetFillablePrice` returns an error.
func (k Keeper) ConstructTakerOrderFromMatchPerpetualLiquidation(
	ctx sdk.Context,
	match *types.MatchPerpetualLiquidation,
) (
	takerOrder *types.LiquidationOrder,
	err error,
) {
	takerClobPair, found := k.GetClobPair(ctx, types.ClobPairId(match.ClobPairId))
	if !found {
		return nil, sdkerrors.Wrapf(
			types.ErrInvalidClob,
			"CLOB pair ID %d not found in state",
			match.ClobPairId,
		)
	}

	perpetualId, err := takerClobPair.GetPerpetualId()
	if err != nil || perpetualId != match.PerpetualId {
		return nil, sdkerrors.Wrapf(
			types.ErrClobPairAndPerpetualDoNotMatch,
			"Clob pair id: %v, perpetual id: %v",
			match.ClobPairId,
			perpetualId,
		)
	}

	deltaQuantumsBig := new(big.Int).SetUint64(match.TotalSize)
	if !match.IsBuy {
		deltaQuantumsBig.Neg(deltaQuantumsBig)
	}
	fillablePrice, err := k.GetFillablePrice(
		ctx,
		match.Liquidated,
		match.PerpetualId,
		deltaQuantumsBig,
	)
	if err != nil {
		return nil, err
	}
	fillablePriceSubticks := k.ConvertFillablePriceToSubticks(
		ctx,
		fillablePrice,
		!match.IsBuy,
		takerClobPair,
	)
	return types.NewLiquidationOrder(
		match.Liquidated,
		takerClobPair,
		match.IsBuy,
		satypes.BaseQuantums(match.TotalSize),
		fillablePriceSubticks,
	), nil
}
