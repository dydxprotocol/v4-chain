package keeper

import (
	"errors"
	"fmt"
	"math/big"
	"time"

	errorsmod "cosmossdk.io/errors"

	gometrics "github.com/armon/go-metrics"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	assettypes "github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

// MaybeDeleverageSubaccount is the main entry point to deleverage a subaccount. It attempts to find positions
// on the opposite side of deltaQuantums and use them to offset the liquidated subaccount's position at
// the bankruptcy price of the liquidated position.
// Note that the full position size will get deleveraged.
func (k Keeper) MaybeDeleverageSubaccount(
	ctx sdk.Context,
	subaccountId satypes.SubaccountId,
	perpetualId uint32,
) (
	quantumsDeleveraged *big.Int,
	err error,
) {
	lib.AssertCheckTxMode(ctx)

	canPerformDeleveraging, err := k.CanDeleverageSubaccount(ctx, subaccountId)
	if err != nil {
		return new(big.Int), err
	}

	// Early return to skip deleveraging if the subaccount can't be deleveraged.
	if !canPerformDeleveraging {
		telemetry.IncrCounter(
			1,
			types.ModuleName,
			metrics.PrepareCheckState,
			metrics.CannotDeleverageSubaccount,
		)
		return new(big.Int), nil
	}

	// Deleverage the entire position for the given perpetual id.
	subaccount := k.subaccountsKeeper.GetSubaccount(ctx, subaccountId)
	position, exists := subaccount.GetPerpetualPositionForId(perpetualId)
	if !exists {
		return nil, types.ErrNoOpenPositionForPerpetual
	}
	deltaQuantums := new(big.Int).Neg(position.GetBigQuantums())

	quantumsDeleveraged, err = k.MemClob.DeleverageSubaccount(ctx, subaccountId, perpetualId, deltaQuantums)

	labels := []gometrics.Label{
		metrics.GetLabelForIntValue(metrics.PerpetualId, int(perpetualId)),
		metrics.GetLabelForBoolValue(metrics.IsLong, deltaQuantums.Sign() == -1),
	}
	if quantumsDeleveraged.Sign() == 0 {
		labels = append(labels, metrics.GetLabelForStringValue(metrics.Status, metrics.Unfilled))
	} else if quantumsDeleveraged.CmpAbs(deltaQuantums) == 0 {
		labels = append(labels, metrics.GetLabelForStringValue(metrics.Status, metrics.FullyFilled))
	} else {
		labels = append(labels, metrics.GetLabelForStringValue(metrics.Status, metrics.PartiallyFilled))
	}
	// Record the status of the deleveraging operation.
	telemetry.IncrCounterWithLabels([]string{types.ModuleName, metrics.DeleverageSubaccount}, 1, labels)

	if quoteQuantums, err := k.perpetualsKeeper.GetNetNotional(
		ctx,
		perpetualId,
		new(big.Int).Abs(deltaQuantums),
	); err == nil {
		telemetry.IncrCounterWithLabels(
			[]string{types.ModuleName, metrics.DeleverageSubaccount, metrics.TotalQuoteQuantums},
			metrics.GetMetricValueFromBigInt(quoteQuantums),
			labels,
		)
		gometrics.AddSampleWithLabels(
			[]string{types.ModuleName, metrics.DeleverageSubaccount, metrics.TotalQuoteQuantums, metrics.Distribution},
			metrics.GetMetricValueFromBigInt(quoteQuantums),
			labels,
		)
	}

	// Record the percent filled of the deleveraging operation as a distribution.
	percentFilled, _ := new(big.Float).Quo(
		new(big.Float).SetInt(new(big.Int).Abs(quantumsDeleveraged)),
		new(big.Float).SetInt(new(big.Int).Abs(deltaQuantums)),
	).Float32()
	gometrics.AddSampleWithLabels(
		[]string{metrics.Deleveraging, metrics.PercentFilled, metrics.Distribution},
		percentFilled,
		labels,
	)

	return quantumsDeleveraged, err
}

// GetInsuranceFundBalance returns the current balance of the insurance fund (in quote quantums).
// This calls the Bank Keeper’s GetBalance() function for the Module Address of the insurance fund.
func (k Keeper) GetInsuranceFundBalance(
	ctx sdk.Context,
) (
	balance *big.Int,
) {
	usdcAsset, exists := k.assetsKeeper.GetAsset(ctx, assettypes.AssetUsdc.Id)
	if !exists {
		panic("GetInsuranceFundBalance: Usdc asset not found in state")
	}
	insuranceFundBalance := k.bankKeeper.GetBalance(
		ctx,
		authtypes.NewModuleAddress(types.InsuranceFundName),
		usdcAsset.Denom,
	)

	// Return as big.Int.
	return insuranceFundBalance.Amount.BigInt()
}

// CanDeleverageSubaccount returns true if a subaccount can be deleveraged.
// Specifically, this function returns true if both of the following are true:
// - The subaccount's total net collateral is negative.
// This function returns an error if `GetNetCollateralAndMarginRequirements` returns an error.
func (k Keeper) CanDeleverageSubaccount(
	ctx sdk.Context,
	subaccountId satypes.SubaccountId,
) (bool, error) {
	bigNetCollateral,
		_,
		_,
		err := k.subaccountsKeeper.GetNetCollateralAndMarginRequirements(
		ctx,
		satypes.Update{SubaccountId: subaccountId},
	)
	if err != nil {
		return false, err
	}

	// Deleveraging cannot be performed if the subaccounts net collateral is non-negative.
	if bigNetCollateral.Sign() >= 0 {
		return false, nil
	}

	// The subaccount's total net collateral is negative, so deleveraging can be performed.
	return true, nil
}

// IsValidInsuranceFundDelta returns true if the insurance fund has enough funds to cover the insurance
// fund delta. Specifically, this function returns true if either of the following are true:
// - The `insuranceFundDelta` is non-negative.
// - The insurance fund balance + `insuranceFundDelta` is greater-than-or-equal-to 0.
func (k Keeper) IsValidInsuranceFundDelta(
	ctx sdk.Context,
	insuranceFundDelta *big.Int,
) bool {
	// Non-negative insurance fund deltas are valid.
	if insuranceFundDelta.Sign() >= 0 {
		return true
	}

	// The insurance fund delta is valid if the insurance fund balance is non-negative after adding
	// the delta.
	currentInsuranceFundBalance := k.GetInsuranceFundBalance(ctx)
	return new(big.Int).Add(currentInsuranceFundBalance, insuranceFundDelta).Sign() >= 0
}

// OffsetSubaccountPerpetualPosition iterates over all subaccounts and use those with positions
// on the opposite side to offset the liquidated subaccount's position by `deltaQuantumsTotal`.
//
// This function returns the fills that were processed and the remaining amount to offset.
// Note that each deleveraging fill is being processed _optimistically_, and the state transitions are
// still persisted even if there are not enough subaccounts to offset the liquidated subaccount's position.
func (k Keeper) OffsetSubaccountPerpetualPosition(
	ctx sdk.Context,
	liquidatedSubaccountId satypes.SubaccountId,
	perpetualId uint32,
	deltaQuantumsTotal *big.Int,
) (
	fills []types.MatchPerpetualDeleveraging_Fill,
	deltaQuantumsRemaining *big.Int,
) {
	defer telemetry.ModuleMeasureSince(
		types.ModuleName,
		time.Now(),
		types.ModuleName,
		metrics.OffsettingSubaccountPerpetualPosition,
	)

	numSubaccountsIterated := uint32(0)
	numSubaccountsWithNonOverlappingBankruptcyPrices := uint32(0)
	numSubaccountsWithNoOpenPositionOnOppositeSide := uint32(0)
	deltaQuantumsRemaining = new(big.Int).Set(deltaQuantumsTotal)
	fills = make([]types.MatchPerpetualDeleveraging_Fill, 0)

	k.subaccountsKeeper.ForEachSubaccountRandomStart(
		ctx,
		func(offsettingSubaccount satypes.Subaccount) (finished bool) {
			// Iterate at most `MaxDeleveragingSubaccountsToIterate` subaccounts.
			if numSubaccountsIterated >= k.Flags.MaxDeleveragingSubaccountsToIterate {
				return true
			}

			numSubaccountsIterated++
			offsettingPosition, _ := offsettingSubaccount.GetPerpetualPositionForId(perpetualId)
			bigOffsettingPositionQuantums := offsettingPosition.GetBigQuantums()

			// Skip subaccounts that do not have a position in the opposite direction as the liquidated subaccount.
			if deltaQuantumsRemaining.Sign() != bigOffsettingPositionQuantums.Sign() {
				numSubaccountsWithNoOpenPositionOnOppositeSide++
				return false
			}

			// TODO(DEC-1495): Determine max amount to offset per offsetting subaccount.
			var deltaQuantums *big.Int
			if deltaQuantumsRemaining.CmpAbs(bigOffsettingPositionQuantums) > 0 {
				deltaQuantums = new(big.Int).Set(bigOffsettingPositionQuantums)
			} else {
				deltaQuantums = new(big.Int).Set(deltaQuantumsRemaining)
			}

			// Try to process the deleveraging operation for both subaccounts.
			if err := k.ProcessDeleveraging(
				ctx,
				liquidatedSubaccountId,
				*offsettingSubaccount.Id,
				perpetualId,
				deltaQuantums,
			); err == nil {
				// Update the remaining liquidatable quantums.
				deltaQuantumsRemaining = new(big.Int).Sub(
					deltaQuantumsRemaining,
					deltaQuantums,
				)
				fills = append(fills, types.MatchPerpetualDeleveraging_Fill{
					OffsettingSubaccountId: *offsettingSubaccount.Id,
					FillAmount:             new(big.Int).Abs(deltaQuantums).Uint64(),
				})
			} else if errors.Is(err, types.ErrInvalidPerpetualPositionSizeDelta) {
				panic(
					fmt.Sprintf(
						"Invalid perpetual position size delta when processing deleveraging. error: %v",
						err,
					),
				)
			} else {
				// If an error is returned, it's likely because the subaccounts' bankruptcy prices do not overlap.
				// TODO(CLOB-75): Support deleveraging subaccounts with non overlapping bankruptcy prices.
				liquidatedSubaccount := k.subaccountsKeeper.GetSubaccount(ctx, liquidatedSubaccountId)
				offsettingSubaccount := k.subaccountsKeeper.GetSubaccount(ctx, *offsettingSubaccount.Id)
				k.Logger(ctx).Debug(
					"Encountered error when processing deleveraging",
					"error", err,
					"blockHeight", ctx.BlockHeight(),
					"checkTx", ctx.IsCheckTx(),
					"perpetualId", perpetualId,
					"deltaQuantums", deltaQuantums,
					"liquidatedSubaccount", log.NewLazySprintf("%+v", liquidatedSubaccount),
					"offsettingSubaccount", log.NewLazySprintf("%+v", offsettingSubaccount),
				)
				numSubaccountsWithNonOverlappingBankruptcyPrices++
			}
			return deltaQuantumsRemaining.Sign() == 0
		},
		k.GetPseudoRand(ctx),
	)

	labels := []gometrics.Label{
		metrics.GetLabelForIntValue(metrics.PerpetualId, int(perpetualId)),
	}
	gometrics.AddSampleWithLabels(
		[]string{
			types.ModuleName, metrics.Deleveraging, metrics.NumSubaccountsIterated, metrics.Count,
		},
		float32(numSubaccountsIterated),
		labels,
	)
	gometrics.AddSampleWithLabels(
		[]string{
			types.ModuleName, metrics.Deleveraging, metrics.NonOverlappingBankruptcyPrices, metrics.Count,
		},
		float32(numSubaccountsWithNonOverlappingBankruptcyPrices),
		labels,
	)
	gometrics.AddSampleWithLabels(
		[]string{
			types.ModuleName, metrics.Deleveraging, metrics.NoOpenPositionOnOppositeSide, metrics.Count,
		},
		float32(numSubaccountsWithNoOpenPositionOnOppositeSide),
		labels,
	)
	return fills, deltaQuantumsRemaining
}

// ProcessDeleveraging processes a deleveraging operation by closing both the liquidated subaccount's
// position and the offsetting subaccount's position at the bankruptcy price of the _liquidated_ position.
// This function takes a `deltaQuantums` argument, which is the delta with respect to the liquidated subaccount's
// position, to allow for partial deleveraging. This function emits a cometbft event if the deleveraging match
// is successfully written to state.
//
// This function returns an error if:
// - `deltaQuantums` is not valid with respect to either of the subaccounts.
// - `GetBankruptcyPriceInQuoteQuantums` returns an error.
// - subaccount updates cannot be applied when the bankruptcy prices of both subaccounts don't overlap.
func (k Keeper) ProcessDeleveraging(
	ctx sdk.Context,
	liquidatedSubaccountId satypes.SubaccountId,
	offsettingSubaccountId satypes.SubaccountId,
	perpetualId uint32,
	deltaQuantums *big.Int,
) (
	err error,
) {
	// Get the liquidated subaccount.
	liquidatedSubaccount := k.subaccountsKeeper.GetSubaccount(ctx, liquidatedSubaccountId)
	liquidatedPosition, _ := liquidatedSubaccount.GetPerpetualPositionForId(perpetualId)
	liquidatedPositionQuantums := liquidatedPosition.GetBigQuantums()

	// Get the offsetting subaccount.
	offsettingSubaccount := k.subaccountsKeeper.GetSubaccount(ctx, offsettingSubaccountId)
	offsettingPosition, _ := offsettingSubaccount.GetPerpetualPositionForId(perpetualId)
	offsettingPositionQuantums := offsettingPosition.GetBigQuantums()

	// Make sure that `deltaQuantums` is valid with respect to the liquidated and offsetting subaccounts
	// by checking that `deltaQuantums` is on the opposite side of the liquidated position side,
	// the same side as the offsetting subaccount position side, and the magnitude of `deltaQuantums`
	// is not larger than both positions.
	if liquidatedPositionQuantums.Sign()*deltaQuantums.Sign() != -1 ||
		liquidatedPositionQuantums.CmpAbs(deltaQuantums) == -1 ||
		offsettingPositionQuantums.Sign()*deltaQuantums.Sign() != 1 ||
		offsettingPositionQuantums.CmpAbs(deltaQuantums) == -1 {
		return errorsmod.Wrapf(
			types.ErrInvalidPerpetualPositionSizeDelta,
			"ProcessDeleveraging: liquidated = (%+v), offsetting = (%+v), perpetual id = (%d), deltaQuantums = (%+v)",
			liquidatedSubaccount,
			offsettingSubaccount,
			perpetualId,
			deltaQuantums,
		)
	}

	// Calculate the bankruptcy price of the liquidated position. This is the price at which both positions
	// are closed.
	bankruptcyPriceQuoteQuantums, err := k.GetBankruptcyPriceInQuoteQuantums(
		ctx,
		liquidatedSubaccountId,
		perpetualId,
		deltaQuantums,
	)
	if err != nil {
		return err
	}

	deleveragedSubaccountQuoteBalanceDelta := bankruptcyPriceQuoteQuantums
	offsettingSubaccountQuoteBalanceDelta := new(big.Int).Neg(bankruptcyPriceQuoteQuantums)
	deleveragedSubaccountPerpetualQuantumsDelta := deltaQuantums
	offsettingSubaccountPerpetualQuantumsDelta := new(big.Int).Neg(deltaQuantums)

	updates := []satypes.Update{
		// Liquidated subaccount update.
		{
			AssetUpdates: []satypes.AssetUpdate{
				{
					AssetId:          assettypes.AssetUsdc.Id,
					BigQuantumsDelta: deleveragedSubaccountQuoteBalanceDelta,
				},
			},
			PerpetualUpdates: []satypes.PerpetualUpdate{
				{
					PerpetualId:      perpetualId,
					BigQuantumsDelta: deleveragedSubaccountPerpetualQuantumsDelta,
				},
			},
			SubaccountId: liquidatedSubaccountId,
		},
		// Offsetting subaccount update.
		{
			AssetUpdates: []satypes.AssetUpdate{
				{
					AssetId:          assettypes.AssetUsdc.Id,
					BigQuantumsDelta: offsettingSubaccountQuoteBalanceDelta,
				},
			},
			PerpetualUpdates: []satypes.PerpetualUpdate{
				{
					PerpetualId:      perpetualId,
					BigQuantumsDelta: offsettingSubaccountPerpetualQuantumsDelta,
				},
			},
			SubaccountId: offsettingSubaccountId,
		},
	}

	// Apply the update.
	success, successPerUpdate, err := k.subaccountsKeeper.UpdateSubaccounts(ctx, updates)
	if err != nil {
		return err
	}

	// If not successful, return error indicating why.
	if updateErr := satypes.GetErrorFromUpdateResults(success, successPerUpdate, updates); updateErr != nil {
		return updateErr
	}

	// Stat quantums deleveraged in quote quantums.
	if deleveragedQuoteQuantums, err := k.perpetualsKeeper.GetNetCollateral(
		ctx,
		perpetualId,
		new(big.Int).Abs(deltaQuantums),
	); err == nil {
		labels := []gometrics.Label{
			metrics.GetLabelForIntValue(metrics.PerpetualId, int(perpetualId)),
			metrics.GetLabelForBoolValue(metrics.CheckTx, ctx.IsCheckTx()),
			metrics.GetLabelForBoolValue(metrics.IsLong, deltaQuantums.Sign() == -1),
		}
		gometrics.AddSampleWithLabels(
			[]string{types.ModuleName, metrics.DeleverageSubaccount, metrics.Filled, metrics.QuoteQuantums},
			metrics.GetMetricValueFromBigInt(deleveragedQuoteQuantums),
			labels,
		)
	}

	// Deleveraging was successful, therefore emit a cometbft event indicating a deleveraging match occurred.
	ctx.EventManager().EmitEvent(
		types.NewCreateMatchEvent(
			liquidatedSubaccountId,
			offsettingSubaccountId,
			big.NewInt(0),
			big.NewInt(0),
			deleveragedSubaccountQuoteBalanceDelta,
			offsettingSubaccountQuoteBalanceDelta,
			deleveragedSubaccountPerpetualQuantumsDelta,
			offsettingSubaccountPerpetualQuantumsDelta,
			big.NewInt(0),
			false, // IsLiquidation is false since this isn't a liquidation match.
			true,  // IsDeleverage is true since this is a deleveraging match.
			perpetualId,
		),
	)

	return nil
}
