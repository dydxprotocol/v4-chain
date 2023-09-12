package keeper

import (
	errorsmod "cosmossdk.io/errors"
	"math/big"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

// ValidateLiquidationOrderAgainstProposedLiquidation performs stateless validation of a liquidation order
// against a proposed liquidation.
// An error is returned when
//   - The CLOB pair IDs of the order and proposed liquidation do not match.
//   - The perpetual IDs of the order and proposed liquidation do not match.
//   - The total size of the order and proposed liquidation do not match.
//   - The side of the order and proposed liquidation do not match.
func (k Keeper) ValidateLiquidationOrderAgainstProposedLiquidation(
	ctx sdk.Context,
	order *types.LiquidationOrder,
	proposedMatch *types.MatchPerpetualLiquidation,
) error {
	if order.GetClobPairId() != types.ClobPairId(proposedMatch.GetClobPairId()) {
		return errorsmod.Wrapf(
			types.ErrClobPairAndPerpetualDoNotMatch,
			"Order CLOB Pair ID: %v, Match CLOB Pair ID: %v",
			order.GetClobPairId(),
			proposedMatch.GetClobPairId(),
		)
	}

	if order.MustGetLiquidatedPerpetualId() != proposedMatch.GetPerpetualId() {
		return errorsmod.Wrapf(
			types.ErrClobPairAndPerpetualDoNotMatch,
			"Order Perpetual ID: %v, Match Perpetual ID: %v",
			order.MustGetLiquidatedPerpetualId(),
			proposedMatch.GetPerpetualId(),
		)
	}

	if order.GetBaseQuantums() != satypes.BaseQuantums(proposedMatch.TotalSize) {
		return errorsmod.Wrapf(
			types.ErrInvalidLiquidationOrderTotalSize,
			"Order Size: %v, Match Size: %v",
			order.GetBaseQuantums(),
			proposedMatch.TotalSize,
		)
	}

	if order.IsBuy() != proposedMatch.GetIsBuy() {
		return errorsmod.Wrapf(
			types.ErrInvalidLiquidationOrderSide,
			"Order Side: %v, Match Side: %v",
			order.IsBuy(),
			proposedMatch.GetIsBuy(),
		)
	}
	return nil
}

func (k Keeper) validateMatchedLiquidation(
	ctx sdk.Context,
	order types.MatchableOrder,
	perpetualId uint32,
	fillAmount satypes.BaseQuantums,
	makerSubticks types.Subticks,
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
	// in the insurance fund (such that the liquidation couldn't have possibly continued).
	if !k.IsValidInsuranceFundDelta(ctx, insuranceFundDelta) {
		k.Logger(ctx).Info("ProcessMatches: insurance fund has insufficient balance to process the liquidation.")
		return nil, errorsmod.Wrapf(
			types.ErrInsuranceFundHasInsufficientFunds,
			"Liquidation order %v, insurance fund delta %v",
			order,
			insuranceFundDelta.String(),
		)
	}

	// Validate that total notional liquidated and total insurance funds lost do not exceed subaccount block limits.
	if err := k.validateLiquidationAgainstSubaccountBlockLimits(
		ctx,
		liquidatedSubaccountId,
		perpetualId,
		fillAmount,
		insuranceFundDelta,
	); err != nil {
		return nil, err
	}

	return insuranceFundDelta, nil
}

// validateLiquidationAgainstSubaccountBlockLimits performs stateful validation
// against the subaccount block limits specified in liquidation configs.
// If validation fails, an error is returned.
//
// The following validation occurs in this method:
//   - The subaccount and perpetual ID pair has not been previously liquidated in the same block.
//   - The total notional liquidated does not exceed the maximum notional amount that a single subaccount
//     can have liquidated per block.
//   - The total insurance lost does not exceed the maximum insurance lost per block.
func (k Keeper) validateLiquidationAgainstSubaccountBlockLimits(
	ctx sdk.Context,
	subaccountId satypes.SubaccountId,
	perpetualId uint32,
	fillAmount satypes.BaseQuantums,
	insuranceFundDeltaQuoteQuantums *big.Int,
) (
	err error,
) {
	// Validate that this liquidation does not exceed the maximum notional amount that a single subaccount can have
	// liquidated per block.
	bigMaxNotionalLiquidatable, err := k.GetSubaccountMaxNotionalLiquidatable(
		ctx,
		subaccountId,
		perpetualId,
	)
	if err != nil {
		return err
	}

	bigNotionalLiquidated, err := k.perpetualsKeeper.GetNetNotional(ctx, perpetualId, fillAmount.ToBigInt())
	if err != nil {
		return err
	}

	if bigNotionalLiquidated.CmpAbs(bigMaxNotionalLiquidatable) > 0 {
		return errorsmod.Wrapf(
			types.ErrLiquidationExceedsSubaccountMaxNotionalLiquidated,
			"Subaccount ID: %v, Perpetual ID: %v, Max Notional Liquidatable: %v, Notional Liquidated: %v",
			subaccountId,
			perpetualId,
			bigMaxNotionalLiquidatable,
			bigNotionalLiquidated,
		)
	}

	// Validate that this liquidation does not exceed the maximum insurance fund payout amount for this
	// subaccount per block.
	if insuranceFundDeltaQuoteQuantums.Sign() == -1 {
		bigMaxQuantumsInsuranceLost, err := k.GetSubaccountMaxInsuranceLost(
			ctx,
			subaccountId,
			perpetualId,
		)
		if err != nil {
			return err
		}

		if insuranceFundDeltaQuoteQuantums.CmpAbs(bigMaxQuantumsInsuranceLost) > 0 {
			return errorsmod.Wrapf(
				types.ErrLiquidationExceedsSubaccountMaxInsuranceLost,
				"Subaccount ID: %v, Perpetual ID: %v, Max Insurance Lost: %v, Insurance Lost: %v",
				subaccountId,
				perpetualId,
				bigMaxQuantumsInsuranceLost,
				insuranceFundDeltaQuoteQuantums,
			)
		}
	}
	return nil
}
