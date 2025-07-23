package keeper

import (
	"errors"
	"fmt"
	"math/big"
	"time"

	errorsmod "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/telemetry"

	sdk "github.com/cosmos/cosmos-sdk/types"
	indexerevents "github.com/dydxprotocol/v4-chain/protocol/indexer/events"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/lib/log"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	assettypes "github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	perplib "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/lib"
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

	ctx = log.AddPersistentTagsToLogger(ctx,
		log.PerpetualId, perpetualId,
		log.Subaccount, subaccountId,
	)

	shouldDeleverageAtBankruptcyPrice, shouldDeleverageAtOraclePrice, err := k.CanDeleverageSubaccount(
		ctx,
		subaccountId,
		perpetualId,
	)
	if err != nil {
		return new(big.Int), err
	}

	// Early return to skip deleveraging if the subaccount doesn't have negative equity or a position in a final
	// settlement market.
	if !shouldDeleverageAtBankruptcyPrice && !shouldDeleverageAtOraclePrice {
		metrics.IncrCounter(
			metrics.ClobPrepareCheckStateCannotDeleverageSubaccount,
			1,
		)
		return new(big.Int), nil
	}

	// Deleverage the entire position for the given perpetual id.
	subaccount := k.subaccountsKeeper.GetSubaccount(ctx, subaccountId)
	position, exists := subaccount.GetPerpetualPositionForId(perpetualId)
	if !exists {
		// Early return to skip deleveraging if the subaccount does not have an open position for the perpetual.
		// This could happen if the subaccount's position was closed by other liquidation matches.
		log.DebugLog(ctx, "Subaccount does not have an open position for the perpetual that is being deleveraged")
		return new(big.Int), nil
	}

	deltaQuantums := new(big.Int).Neg(position.GetBigQuantums())
	quantumsDeleveraged, err = k.MemClob.DeleverageSubaccount(
		ctx,
		subaccountId,
		perpetualId,
		deltaQuantums,
		shouldDeleverageAtOraclePrice,
	)

	labels := []metrics.Label{
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
	metrics.IncrCounterWithLabels(
		metrics.ClobDeleverageSubaccount,
		1,
		labels...,
	)

	if quoteQuantums, err := k.perpetualsKeeper.GetNetNotional(
		ctx,
		perpetualId,
		new(big.Int).Abs(deltaQuantums),
	); err == nil {
		metrics.IncrCounterWithLabels(
			metrics.ClobDeleverageSubaccountTotalQuoteQuantums,
			metrics.GetMetricValueFromBigInt(quoteQuantums),
			labels...,
		)

		metrics.AddSampleWithLabels(
			metrics.ClobDeleverageSubaccountTotalQuoteQuantumsDistribution,
			metrics.GetMetricValueFromBigInt(quoteQuantums),
			labels...,
		)
	}

	// Record the percent filled of the deleveraging operation as a distribution.
	percentFilled, _ := new(big.Float).Quo(
		new(big.Float).SetInt(new(big.Int).Abs(quantumsDeleveraged)),
		new(big.Float).SetInt(new(big.Int).Abs(deltaQuantums)),
	).Float32()

	metrics.AddSampleWithLabels(
		metrics.DeleveragingPercentFilledDistribution,
		percentFilled,
		labels...,
	)

	return quantumsDeleveraged, err
}

// CanDeleverageSubaccount returns true if a subaccount can be deleveraged.
// This function returns two booleans, shouldDeleverageAtBankruptcyPrice and shouldDeleverageAtOraclePrice.
// - shouldDeleverageAtBankruptcyPrice is true if the subaccount has negative TNC.
// - shouldDeleverageAtOraclePrice is true if the subaccount has non-negative TNC and the market is in final settlement.
// This function returns an error if `GetNetCollateralAndMarginRequirements` returns an error or if there is
// an error when fetching the clob pair for the provided perpetual.
func (k Keeper) CanDeleverageSubaccount(
	ctx sdk.Context,
	subaccountId satypes.SubaccountId,
	perpetualId uint32,
) (shouldDeleverageAtBankruptcyPrice bool, shouldDeleverageAtOraclePrice bool, err error) {
	risk, err := k.subaccountsKeeper.GetNetCollateralAndMarginRequirements(
		ctx,
		satypes.Update{SubaccountId: subaccountId},
	)
	if err != nil {
		return false, false, err
	}

	// Negative TNC, deleverage at bankruptcy price.
	if risk.NC.Sign() == -1 {
		return true, false, nil
	}

	clobPairId, err := k.GetClobPairIdForPerpetual(ctx, perpetualId)
	if err != nil {
		return false, false, err
	}
	clobPair := k.mustGetClobPair(ctx, clobPairId)

	// Non-negative TNC, deleverage at oracle price if market is in final settlement. Deleveraging at oracle price
	// is always a valid state transition when TNC is non-negative. This is because the TNC/TMMR ratio is improving;
	// TNC is staying constant while TMMR is decreasing.
	return false, clobPair.Status == types.ClobPair_STATUS_FINAL_SETTLEMENT, nil
}

// GateWithdrawalsIfNegativeTncSubaccountSeen gates withdrawals if a negative TNC subaccount exists.
// It does this by inserting a zero-fill deleveraging operation into the operations queue iff any of
// the provided negative TNC subaccounts are still negative TNC.
func (k Keeper) GateWithdrawalsIfNegativeTncSubaccountSeen(
	ctx sdk.Context,
	negativeTncSubaccountIds []satypes.SubaccountId,
) (err error) {
	defer metrics.ModuleMeasureSince(
		types.ModuleName,
		metrics.GateWithdrawalsIfNegativeTncSubaccountSeenLatency,
		time.Now(),
	)
	metrics.IncrCounter(
		metrics.GateWithdrawalsIfNegativeTncSubaccountSeen,
		1,
	)

	foundNegativeTncSubaccount := false
	var negativeTncSubaccountId satypes.SubaccountId
	for _, subaccountId := range negativeTncSubaccountIds {
		risk, err := k.subaccountsKeeper.GetNetCollateralAndMarginRequirements(
			ctx,
			satypes.Update{SubaccountId: subaccountId},
		)
		if err != nil {
			return err
		}

		// If the subaccount has negative TNC, mark that a negative TNC subaccount was found.
		if risk.NC.Sign() == -1 {
			foundNegativeTncSubaccount = true
			negativeTncSubaccountId = subaccountId
			break
		}
	}

	if !foundNegativeTncSubaccount {
		return nil
	}

	// A negative TNC subaccount was found, therefore insert a zero-fill deleveraging operation into
	// the operations queue to indicate withdrawals should be gated.
	subaccount := k.subaccountsKeeper.GetSubaccount(ctx, negativeTncSubaccountId)
	perpetualPositions := subaccount.GetPerpetualPositions()
	if len(perpetualPositions) == 0 {
		return errorsmod.Wrapf(
			types.ErrNoPerpetualPositionsToLiquidate,
			"GateWithdrawalsIfNegativeTncSubaccountSeen: subaccount has no open positions: (%s)",
			lib.MaybeGetJsonString(subaccount),
		)
	}
	perpetualId := subaccount.PerpetualPositions[0].PerpetualId
	k.MemClob.InsertZeroFillDeleveragingIntoOperationsQueue(negativeTncSubaccountId, perpetualId)
	metrics.IncrCountMetricWithLabels(
		types.ModuleName,
		metrics.SubaccountsNegativeTncSubaccountSeen,
		metrics.GetLabelForIntValue(metrics.PerpetualId, int(perpetualId)),
		metrics.GetLabelForBoolValue(metrics.IsLong, subaccount.PerpetualPositions[0].GetIsLong()),
		metrics.GetLabelForBoolValue(metrics.DeliverTx, false),
	)

	return nil
}

// IsValidInsuranceFundDelta returns true if the insurance fund has enough funds to cover the insurance
// fund delta. Specifically, this function returns true if either of the following are true:
// - The `insuranceFundDelta` is non-negative.
// - The insurance fund balance + `insuranceFundDelta` is greater-than-or-equal-to 0.
func (k Keeper) IsValidInsuranceFundDelta(ctx sdk.Context, insuranceFundDelta *big.Int, perpetualId uint32) bool {
	// Non-negative insurance fund deltas are valid.
	if insuranceFundDelta.Sign() >= 0 {
		return true
	}

	// The insurance fund delta is valid if the insurance fund balance is non-negative after adding
	// the delta.
	currentInsuranceFundBalance := k.subaccountsKeeper.GetInsuranceFundBalance(ctx, perpetualId)
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
	isFinalSettlement bool,
) (
	fills []types.MatchPerpetualDeleveraging_Fill,
	deltaQuantumsRemaining *big.Int,
) {
	defer metrics.ModuleMeasureSince(
		types.ModuleName,
		metrics.ClobOffsettingSubaccountPerpetualPosition,
		time.Now(),
	)

	numSubaccountsIterated := uint32(0)
	numSubaccountsWithNonOverlappingBankruptcyPrices := uint32(0)
	numSubaccountsWithNoOpenPositionOnOppositeSide := uint32(0)
	deltaQuantumsRemaining = new(big.Int).Set(deltaQuantumsTotal)
	fills = make([]types.MatchPerpetualDeleveraging_Fill, 0)

	// Find subaccounts with open positions on the opposite side of the liquidated subaccount.
	isDeleveragingLong := deltaQuantumsTotal.Sign() == -1
	subaccountsWithOpenPositions := k.DaemonLiquidationInfo.GetSubaccountsWithOpenPositionsOnSide(
		perpetualId,
		!isDeleveragingLong,
	)

	numSubaccounts := len(subaccountsWithOpenPositions)
	if numSubaccounts == 0 {
		liquidatedSubaccount := k.subaccountsKeeper.GetSubaccount(ctx, liquidatedSubaccountId)
		log.ErrorLog(
			ctx,
			"Failed to find subaccounts with open positions on opposite side of liquidated subaccount",
			"deltaQuantumsTotal", deltaQuantumsTotal,
			"liquidatedSubaccount", liquidatedSubaccount,
		)
		return fills, deltaQuantumsRemaining
	}

	// Start from a random subaccount.
	pseudoRand := k.GetPseudoRand(ctx)
	indexOffset := pseudoRand.Intn(numSubaccounts)

	// Iterate at most `MaxDeleveragingSubaccountsToIterate` subaccounts.
	numSubaccountsToIterate := lib.Min(numSubaccounts, int(k.Flags.MaxDeleveragingSubaccountsToIterate))

	for i := 0; i < numSubaccountsToIterate && deltaQuantumsRemaining.Sign() != 0; i++ {
		index := (i + indexOffset) % numSubaccounts
		subaccountId := subaccountsWithOpenPositions[index]

		numSubaccountsIterated++
		offsettingSubaccount := k.subaccountsKeeper.GetSubaccount(ctx, subaccountId)
		offsettingPosition, _ := offsettingSubaccount.GetPerpetualPositionForId(perpetualId)
		bigOffsettingPositionQuantums := offsettingPosition.GetBigQuantums()

		// Skip subaccounts that do not have a position in the opposite direction as the liquidated subaccount.
		if deltaQuantumsRemaining.Sign() != bigOffsettingPositionQuantums.Sign() {
			numSubaccountsWithNoOpenPositionOnOppositeSide++
			continue
		}

		// TODO(DEC-1495): Determine max amount to offset per offsetting subaccount.
		var deltaBaseQuantums *big.Int
		if deltaQuantumsRemaining.CmpAbs(bigOffsettingPositionQuantums) > 0 {
			deltaBaseQuantums = new(big.Int).Set(bigOffsettingPositionQuantums)
		} else {
			deltaBaseQuantums = new(big.Int).Set(deltaQuantumsRemaining)
		}

		// Fetch delta quote quantums. Calculated at bankruptcy price for standard
		// deleveraging and at oracle price for final settlement deleveraging.
		deltaQuoteQuantums, err := k.getDeleveragingQuoteQuantumsDelta(
			ctx,
			perpetualId,
			liquidatedSubaccountId,
			deltaBaseQuantums,
			isFinalSettlement,
		)
		if err != nil {
			liquidatedSubaccount := k.subaccountsKeeper.GetSubaccount(ctx, liquidatedSubaccountId)
			log.ErrorLogWithError(ctx, "Encountered error when getting quote quantums for deleveraging",
				err,
				"deltaBaseQuantums", deltaBaseQuantums,
				"liquidatedSubaccount", liquidatedSubaccount,
				"offsettingSubaccount", offsettingSubaccount,
			)
			continue
		}

		// Try to process the deleveraging operation for both subaccounts.
		if err := k.ProcessDeleveraging(
			ctx,
			liquidatedSubaccountId,
			*offsettingSubaccount.Id,
			perpetualId,
			deltaBaseQuantums,
			deltaQuoteQuantums,
		); err == nil {
			// Update the remaining liquidatable quantums.
			deltaQuantumsRemaining.Sub(deltaQuantumsRemaining, deltaBaseQuantums)
			fills = append(fills, types.MatchPerpetualDeleveraging_Fill{
				OffsettingSubaccountId: *offsettingSubaccount.Id,
				FillAmount:             new(big.Int).Abs(deltaBaseQuantums).Uint64(),
			})
			// Send on-chain update for the deleveraging. The events are stored in a TransientStore which should be rolled-back
			// if the branched state is discarded, so batching is not necessary.
			k.GetIndexerEventManager().AddTxnEvent(
				ctx,
				indexerevents.SubtypeDeleveraging,
				indexerevents.DeleveragingEventVersion,
				indexer_manager.GetBytes(
					indexerevents.NewDeleveragingEvent(
						liquidatedSubaccountId,
						*offsettingSubaccount.Id,
						perpetualId,
						satypes.BaseQuantums(new(big.Int).Abs(deltaBaseQuantums).Uint64()),
						// TODO(CT-641): Use the actual unit price rather than the total quote quantums.
						satypes.BaseQuantums(deltaQuoteQuantums.Uint64()),
						deltaBaseQuantums.Sign() > 0,
						isFinalSettlement,
					),
				),
			)
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
			log.DebugLog(ctx, "Encountered error when processing deleveraging",
				err,
				"blockHeight", ctx.BlockHeight(),
				"checkTx", ctx.IsCheckTx(),
				"perpetualId", perpetualId,
				"deltaBaseQuantums", deltaBaseQuantums,
				"liquidatedSubaccount", liquidatedSubaccount,
				"offsettingSubaccount", offsettingSubaccount,
			)
			numSubaccountsWithNonOverlappingBankruptcyPrices++
		}
	}

	labels := []metrics.Label{
		metrics.GetLabelForIntValue(metrics.PerpetualId, int(perpetualId)),
	}

	metrics.AddSampleWithLabels(
		metrics.ClobDeleveragingNumSubaccountsIteratedCount,
		float32(numSubaccountsIterated),
		labels...,
	)

	metrics.AddSampleWithLabels(
		metrics.ClobDeleveragingNonOverlappingBankrupcyPricesCount,
		float32(numSubaccountsWithNonOverlappingBankruptcyPrices),
		labels...,
	)
	metrics.AddSampleWithLabels(
		metrics.ClobDeleveragingNoOpenPositionOnOppositeSideCount,
		float32(numSubaccountsWithNoOpenPositionOnOppositeSide),
		labels...,
	)
	return fills, deltaQuantumsRemaining
}

// getDeleveragingQuoteQuantums returns the quote quantums delta to apply to a deleveraging operation.
// This returns the bankruptcy price for standard deleveraging operations, and the oracle price for
// final settlement deleveraging operations.
func (k Keeper) getDeleveragingQuoteQuantumsDelta(
	ctx sdk.Context,
	perpetualId uint32,
	subaccountId satypes.SubaccountId,
	deltaQuantums *big.Int,
	isFinalSettlement bool,
) (deltaQuoteQuantums *big.Int, err error) {
	// If market is in final settlement and the subaccount has non-negative TNC, use the oracle price.
	if isFinalSettlement {
		return k.perpetualsKeeper.GetNetNotional(ctx, perpetualId, new(big.Int).Neg(deltaQuantums))
	}

	// For standard deleveraging, use the bankruptcy price.
	return k.GetBankruptcyPriceInQuoteQuantums(
		ctx,
		subaccountId,
		perpetualId,
		deltaQuantums,
	)
}

// ProcessDeleveraging processes a deleveraging operation by closing both the liquidated subaccount's
// position and the offsetting subaccount's position at the bankruptcy price of the _liquidated_ position.
// This function takes a `deltaQuantums` argument, which is the delta with respect to the liquidated subaccount's
// position, to allow for partial deleveraging. This function emits a cometbft event if the deleveraging match
// is successfully written to state.
//
// This function returns an error if:
// - `deltaBaseQuantums` is not valid with respect to either of the subaccounts.
// - `GetBankruptcyPriceInQuoteQuantums` returns an error.
// - subaccount updates cannot be applied when the bankruptcy prices of both subaccounts don't overlap.
func (k Keeper) ProcessDeleveraging(
	ctx sdk.Context,
	liquidatedSubaccountId satypes.SubaccountId,
	offsettingSubaccountId satypes.SubaccountId,
	perpetualId uint32,
	deltaBaseQuantums *big.Int,
	deltaQuoteQuantums *big.Int,
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
	if liquidatedPositionQuantums.Sign()*deltaBaseQuantums.Sign() != -1 ||
		liquidatedPositionQuantums.CmpAbs(deltaBaseQuantums) == -1 ||
		offsettingPositionQuantums.Sign()*deltaBaseQuantums.Sign() != 1 ||
		offsettingPositionQuantums.CmpAbs(deltaBaseQuantums) == -1 {
		return errorsmod.Wrapf(
			types.ErrInvalidPerpetualPositionSizeDelta,
			"ProcessDeleveraging: liquidated = (%s), offsetting = (%s), perpetual id = (%d), deltaQuantums = (%+v)",
			lib.MaybeGetJsonString(liquidatedSubaccount),
			lib.MaybeGetJsonString(offsettingSubaccount),
			perpetualId,
			deltaBaseQuantums,
		)
	}

	deleveragedSubaccountQuoteBalanceDelta := deltaQuoteQuantums
	offsettingSubaccountQuoteBalanceDelta := new(big.Int).Neg(deltaQuoteQuantums)
	deleveragedSubaccountPerpetualQuantumsDelta := deltaBaseQuantums
	offsettingSubaccountPerpetualQuantumsDelta := new(big.Int).Neg(deltaBaseQuantums)

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
	success, successPerUpdate, err := k.subaccountsKeeper.UpdateSubaccounts(ctx, updates, satypes.Match)
	if err != nil {
		return err
	}

	// If not successful, return error indicating why.
	if updateErr := satypes.GetErrorFromUpdateResults(success, successPerUpdate, updates); updateErr != nil {
		return updateErr
	}

	// Stat quantums deleveraged in quote quantums.
	if perpetual, marketPrice, err := k.perpetualsKeeper.GetPerpetualAndMarketPrice(ctx, perpetualId); err == nil {
		deleveragedQuoteQuantums := perplib.GetNetNotionalInQuoteQuantums(
			perpetual,
			marketPrice,
			new(big.Int).Abs(deltaBaseQuantums),
		)
		labels := []metrics.Label{
			metrics.GetLabelForIntValue(metrics.PerpetualId, int(perpetualId)),
			metrics.GetLabelForBoolValue(metrics.CheckTx, ctx.IsCheckTx()),
			metrics.GetLabelForBoolValue(metrics.IsLong, deltaBaseQuantums.Sign() == -1),
		}

		metrics.AddSampleWithLabels(
			metrics.ClobDeleverageSubaccountFilledQuoteQuantums,
			metrics.GetMetricValueFromBigInt(deleveragedQuoteQuantums),
			labels...,
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
			// Builder Code Params
			"",
			"",
			big.NewInt(0),
			big.NewInt(0),
			// Order Router Rev Share Params
			"",
			"",
			big.NewInt(0),
			big.NewInt(0),
		),
	)

	return nil
}

// GetSubaccountsWithPositionsInFinalSettlementMarkets uses the subaccountOpenPositionInfo returned from the
// liquidations daemon to fetch subaccounts with open positions in final settlement markets. These subaccounts
// will be deleveraged in either at the oracle price if non-negative TNC or at bankruptcy price if negative TNC. This
// function is called in PrepareCheckState during the deleveraging step.
func (k Keeper) GetSubaccountsWithPositionsInFinalSettlementMarkets(
	ctx sdk.Context,
) (subaccountsToDeleverage []subaccountToDeleverage) {
	defer telemetry.MeasureSince(
		time.Now(),
		types.ModuleName,
		metrics.ClobGetSubaccountsWithPositionsInFinalSettlementMarkets,
		metrics.Latency,
	)

	for _, clobPair := range k.GetAllClobPairs(ctx) {
		if clobPair.Status != types.ClobPair_STATUS_FINAL_SETTLEMENT {
			continue
		}

		finalSettlementPerpetualId := clobPair.MustGetPerpetualId()
		subaccountsWithPosition := k.DaemonLiquidationInfo.GetSubaccountsWithOpenPositions(
			finalSettlementPerpetualId,
		)
		for _, subaccountId := range subaccountsWithPosition {
			subaccountsToDeleverage = append(subaccountsToDeleverage, subaccountToDeleverage{
				SubaccountId: subaccountId,
				PerpetualId:  finalSettlementPerpetualId,
			})
		}
	}

	metrics.AddSample(
		metrics.ClobSubaccountsWithFinalSettlementPositionsCount,
		float32(len(subaccountsToDeleverage)),
	)

	return subaccountsToDeleverage
}

// DeleverageSubaccounts deleverages a slice of subaccounts paired with a perpetual position to deleverage with.
// Returns an error if a deleveraging attempt returns an error.
func (k Keeper) DeleverageSubaccounts(
	ctx sdk.Context,
	subaccountsToDeleverage []subaccountToDeleverage,
) error {
	defer telemetry.MeasureSince(
		time.Now(),
		types.ModuleName,
		metrics.LiquidateSubaccounts_Deleverage,
		metrics.Latency,
	)

	// For each unfilled liquidation, attempt to deleverage the subaccount.
	for i := 0; i < int(k.Flags.MaxDeleveragingAttemptsPerBlock) && i < len(subaccountsToDeleverage); i++ {
		subaccountId := subaccountsToDeleverage[i].SubaccountId
		perpetualId := subaccountsToDeleverage[i].PerpetualId
		_, err := k.MaybeDeleverageSubaccount(ctx, subaccountId, perpetualId)
		if err != nil {
			log.ErrorLogWithError(
				ctx,
				"Failed to deleverage subaccount.",
				err,
				"subaccount", subaccountId,
				"perpetualId", perpetualId,
			)
			return err
		}
	}

	return nil
}
