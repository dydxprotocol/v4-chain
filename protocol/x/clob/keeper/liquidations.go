package keeper

import (
	"errors"
	"fmt"
	"math"
	"math/big"
	"time"

	errorsmod "cosmossdk.io/errors"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/dtypes"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/lib"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/lib/log"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/lib/metrics"
	assetstypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/assets/types"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/clob/heap"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/clob/types"
	perpkeeper "github.com/StreamFinance-Protocol/stream-chain/protocol/x/perpetuals/keeper"
	pricestypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices/types"
	satypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/subaccounts/types"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// LiquidateSubaccountsAgainstOrderbook takes a list of subaccount IDs and liquidates them against
// the orderbook. It will liquidate as many subaccounts as possible up to the maximum number of
// liquidations per block. Subaccounts are selected with a pseudo-randomly generated offset. A slice
// of subaccounts to deleverage is returned from this function, derived from liquidation orders that
// failed to fill.
func (k Keeper) LiquidateSubaccountsAgainstOrderbook(
	ctx sdk.Context,
	subaccountIds *heap.LiquidationPriorityHeap,
) (
	subaccountsToDeleverage []heap.SubaccountToDeleverage,
	err error,
) {

	lib.AssertCheckTxMode(ctx)

	defer telemetry.MeasureSince(
		time.Now(),
		types.ModuleName,
		metrics.ClobLiquidateSubaccountsAgainstOrderbook,
		metrics.Latency,
	)

	metrics.AddSample(
		metrics.LiquidationsLiquidatableSubaccountIdsCount,
		float32(subaccountIds.Len()),
	)
	startGetLiquidationOrders := time.Now()

	if subaccountIds.Len() == 0 {
		return nil, nil
	}
	subaccountsToDeleverage, err = k.LiquidateSubaccountsAgainstOrderbookInternal(ctx, subaccountIds, heap.NewLiquidationPriorityHeap())
	if err != nil {
		return nil, err
	}

	telemetry.MeasureSince(
		startGetLiquidationOrders,
		types.ModuleName,
		metrics.LiquidateSubaccounts_PlaceLiquidations,
		metrics.Latency,
	)
	metrics.SetGaugeWithLabels(
		metrics.ClobSubaccountsRequiringDeleveragingCount,
		float32(len(subaccountsToDeleverage)),
	)

	return subaccountsToDeleverage, nil
}

// LiquidateSubaccountsAgainstOrderbookInternal is a helper function that performs the core logic of
// liquidating subaccounts against the orderbook.
func (k Keeper) LiquidateSubaccountsAgainstOrderbookInternal(
	ctx sdk.Context,
	subaccountIds *heap.LiquidationPriorityHeap,
	isolatedPositionsPriorityHeap *heap.LiquidationPriorityHeap,
) (
	subaccountsToDeleverage []heap.SubaccountToDeleverage,
	err error,
) {
	numIsolatedLiquidations := 0
	for i := 0; i < int(k.Flags.MaxLiquidationAttemptsPerBlock); i++ {

		subaccount, subaccountId := k.GetNextSubaccountToLiquidate(ctx, subaccountIds, isolatedPositionsPriorityHeap, &numIsolatedLiquidations)
		if subaccountId == nil {
			break
		}

		// will always have at least one perpetual position
		isIsolated, err := k.perpetualsKeeper.IsIsolatedPerpetual(ctx, subaccount.PerpetualPositions[0].PerpetualId)
		if err != nil {
			return nil, err
		}
		if isIsolated {
			if numIsolatedLiquidations < int(k.Flags.MaxIsolatedLiquidationAttemptsPerBlock) {
				numIsolatedLiquidations++
			} else {
				isolatedPositionsPriorityHeap.AddSubaccount(subaccountId.SubaccountId, subaccountId.Priority)
				i--
				continue
			}
		}

		// Generate a new liquidation order with the appropriate order size from the sorted subaccount ids.
		liquidationOrder, err := k.MaybeGetLiquidationOrder(ctx, subaccountId.SubaccountId)
		if err == types.ErrNoPerpetualPositionsToLiquidate {
			i--
			continue
		} else if err != nil {
			return nil, err
		}

		optimisticallyFilledQuantums, _, err := k.PlacePerpetualLiquidation(ctx, *liquidationOrder)
		// Exception for liquidation which conflicts with clob pair status. This is expected for liquidations generated
		// for subaccounts with open positions in final settlement markets.
		if err != nil {
			if !errors.Is(err, types.ErrLiquidationConflictsWithClobPairStatus) {
				return nil, err
			} else {
				i--
				continue
			}
		}

		err = k.handleLiquidationOrderPlacementResult(ctx, liquidationOrder, optimisticallyFilledQuantums, &subaccountsToDeleverage, subaccountIds)
		if err != nil {
			return nil, err
		}
	}

	return subaccountsToDeleverage, nil
}

func (k Keeper) GetNextSubaccountToLiquidate(
	ctx sdk.Context,
	subaccountIds *heap.LiquidationPriorityHeap,
	isolatedPositionsPriorityHeap *heap.LiquidationPriorityHeap,
	numIsolatedLiquidations *int,
) (
	subaccount satypes.Subaccount,
	subaccountId *heap.LiquidationPriority,
) {
	// If we have exceeded the max numIsolatedLiquidations and there are no more non-isolated subaccounts to liquidate
	if subaccountIds.Len() == 0 {
		if isolatedPositionsPriorityHeap.Len() > 0 {
			*subaccountIds = *isolatedPositionsPriorityHeap
			*isolatedPositionsPriorityHeap = *heap.NewLiquidationPriorityHeap()
			*numIsolatedLiquidations = -1000000
		} else {
			return satypes.Subaccount{}, nil
		}
	}

	subaccountId = subaccountIds.PopLowestPriority()
	subaccount = k.subaccountsKeeper.GetSubaccount(ctx, subaccountId.SubaccountId)

	return subaccount, subaccountId
}

func (k Keeper) handleLiquidationOrderPlacementResult(
	ctx sdk.Context,
	liquidationOrder *types.LiquidationOrder,
	optimisticallyFilledQuantums satypes.BaseQuantums,
	subaccountsToDeleverage *[]heap.SubaccountToDeleverage,
	subaccountIds *heap.LiquidationPriorityHeap,
) error {
	if optimisticallyFilledQuantums == 0 {
		*subaccountsToDeleverage = append(*subaccountsToDeleverage, heap.SubaccountToDeleverage{
			SubaccountId: liquidationOrder.GetSubaccountId(),
			PerpetualId:  liquidationOrder.MustGetLiquidatedPerpetualId(),
		})
		return nil
	}

	return k.insertIntoLiquidationHeapIfUnhealthy(ctx, liquidationOrder.GetSubaccountId(), subaccountIds)
}

func (k Keeper) insertIntoLiquidationHeapIfUnhealthy(
	ctx sdk.Context,
	subaccountId satypes.SubaccountId,
	subaccountIds *heap.LiquidationPriorityHeap,
) error {

	subaccount := k.subaccountsKeeper.GetSubaccount(ctx, subaccountId)
	isLiquidatable, priority, err := k.GetSubaccountPriority(ctx, subaccount)
	if err != nil {
		return err
	}

	fmt.Printf("insertIntoLiquidationHeapIfUnhealthy subaccountId: %v, priority: %v\n",
		subaccountId, priority)

	if isLiquidatable {
		subaccountIds.AddSubaccount(subaccountId, priority)
	}

	return nil
}

func (k Keeper) GetSubaccountPriority(
	ctx sdk.Context,
	subaccount satypes.Subaccount,
) (
	isLiquidatable bool,
	priority *big.Float,
	err error,
) {

	_, marketPricesMap, perpetualsMap, liquidityTiersMap := k.FetchInformationForLiquidations(ctx)

	isLiquidatable, _, priority, err = k.GetSubaccountCollateralizationInfo(ctx, subaccount, marketPricesMap, perpetualsMap, liquidityTiersMap)

	return isLiquidatable, priority, err

}

// MaybeGetLiquidationOrder takes a subaccount ID and returns a liquidation order that can be used to
// liquidate the subaccount.
// The account is assumed to be liquidatable
func (k Keeper) MaybeGetLiquidationOrder(
	ctx sdk.Context,
	subaccountId satypes.SubaccountId,
) (
	liquidationOrder *types.LiquidationOrder,
	err error,
) {

	defer telemetry.ModuleMeasureSince(
		types.ModuleName,
		time.Now(),
		metrics.ConstructLiquidationOrder,
	)

	// The subaccount is liquidatable. Get the perpetual position and position size to liquidate.
	perpetualId, err := k.GetBestPerpetualPositionToLiquidate(ctx, subaccountId)
	if err != nil {
		return nil, err
	}

	return k.GetLiquidationOrderForPerpetual(ctx, subaccountId, perpetualId)
}

func (k Keeper) GetLiquidationOrderForPerpetual(
	ctx sdk.Context,
	subaccountId satypes.SubaccountId,
	perpetualId uint32,
) (
	liquidationOrder *types.LiquidationOrder,
	err error,
) {
	clobPair := k.mustGetClobPairForPerpetualId(ctx, perpetualId)
	orderQuantums, err := k.GetNegativePositionSize(ctx, subaccountId, perpetualId)
	if err != nil {
		return nil, err
	}
	isPositionLong := orderQuantums.Sign() == -1

	liquidationPrice, err := k.getLiquidationPrice(ctx, subaccountId, perpetualId, orderQuantums, isPositionLong, clobPair)
	if err != nil {
		return nil, err
	}

	liquidationOrder = types.NewLiquidationOrder(
		subaccountId,
		clobPair,
		!isPositionLong,
		satypes.BaseQuantums(orderQuantums.Abs(orderQuantums).Uint64()),
		liquidationPrice,
	)
	return liquidationOrder, nil
}

// PlacePerpetualLiquidation places an IOC liquidation order onto the book that results in fills of type
// `PerpetualLiquidation`. This function will return an error if attempting to place a liquidation order
// in a non-active market.
func (k Keeper) PlacePerpetualLiquidation(
	ctx sdk.Context,
	liquidationOrder types.LiquidationOrder,
) (
	orderSizeOptimisticallyFilledFromMatchingQuantums satypes.BaseQuantums,
	orderStatus types.OrderStatus,
	err error,
) {
	lib.AssertCheckTxMode(ctx)

	defer telemetry.ModuleMeasureSince(
		types.ModuleName,
		time.Now(),
		metrics.PlacePerpetualLiquidation,
	)

	if err := k.validateLiquidationAgainstClobPairStatus(ctx, liquidationOrder); err != nil {
		return 0, 0, err
	}

	orderSizeOptimisticallyFilledFromMatchingQuantums, orderStatus, offchainUpdates, err := k.MemClob.PlacePerpetualLiquidation(ctx, liquidationOrder)
	if err != nil {
		return 0, 0, err
	}

	perpetualId := liquidationOrder.MustGetLiquidatedPerpetualId()
	k.MustUpdateSubaccountPerpetualLiquidated(ctx, liquidationOrder.GetSubaccountId(), perpetualId)

	k.handleLiquidationMetrics(ctx, liquidationOrder, orderSizeOptimisticallyFilledFromMatchingQuantums, perpetualId)
	k.SendOffchainMessages(offchainUpdates, nil, metrics.SendPlacePerpetualLiquidationOffchainUpdates)
	return orderSizeOptimisticallyFilledFromMatchingQuantums, orderStatus, err
}

func (k Keeper) handleLiquidationMetrics(
	ctx sdk.Context,
	liquidationOrder types.LiquidationOrder,
	orderSizeOptimisticallyFilledFromMatchingQuantums satypes.BaseQuantums,
	perpetualId uint32,
) {

	labels := []metrics.Label{
		metrics.GetLabelForIntValue(metrics.PerpetualId, int(perpetualId)),
	}
	if liquidationOrder.IsBuy() {
		labels = append(labels, metrics.GetLabelForStringValue(metrics.OrderSide, metrics.Buy))
	} else {
		labels = append(labels, metrics.GetLabelForStringValue(metrics.OrderSide, metrics.Sell))
	}

	// Record the percent filled of the liquidation as a distribution.
	percentFilled, _ := new(big.Float).Quo(
		new(big.Float).SetUint64(orderSizeOptimisticallyFilledFromMatchingQuantums.ToUint64()),
		new(big.Float).SetUint64(liquidationOrder.GetBaseQuantums().ToUint64()),
	).Float32()

	metrics.AddSampleWithLabels(
		metrics.LiquidationsPercentFilledDistribution,
		percentFilled,
		labels...,
	)

	if orderSizeOptimisticallyFilledFromMatchingQuantums == 0 {
		labels = append(labels, metrics.GetLabelForStringValue(metrics.Status, metrics.Unfilled))
	} else if orderSizeOptimisticallyFilledFromMatchingQuantums == liquidationOrder.GetBaseQuantums() {
		labels = append(labels, metrics.GetLabelForStringValue(metrics.Status, metrics.FullyFilled))
	} else {
		labels = append(labels, metrics.GetLabelForStringValue(metrics.Status, metrics.PartiallyFilled))
	}
	// Stat the number of liquidation orders placed.
	telemetry.IncrCounterWithLabels(
		[]string{metrics.Liquidations, metrics.PlacePerpetualLiquidation, metrics.Count},
		1,
		labels,
	)

	// Stat the volume of liquidation orders placed.
	if totalQuoteQuantums, err := k.perpetualsKeeper.GetNetNotional(
		ctx,
		perpetualId,
		liquidationOrder.GetBaseQuantums().ToBigInt(),
	); err == nil {
		metrics.IncrCounterWithLabels(
			metrics.LiquidationsPlacePerpetualLiquidationQuoteQuantums,
			metrics.GetMetricValueFromBigInt(totalQuoteQuantums),
			labels...,
		)

		metrics.AddSampleWithLabels(
			metrics.LiquidationsPlacePerpetualLiquidationQuoteQuantumsDistribution,
			metrics.GetMetricValueFromBigInt(totalQuoteQuantums),
			labels...,
		)
	}
}

// IsLiquidatable returns true if the subaccount is able to be liquidated; that is,
// if-and-only-if the maintenance margin requirement is non-zero and greater than the net collateral
// of the subaccount.
// If `GetNetCollateralAndMarginRequirements` returns an error, this function will return that
// error to the caller.
func (k Keeper) IsLiquidatable(
	ctx sdk.Context,
	subaccountId satypes.SubaccountId,
) (
	bool,
	error,
) {
	bigNetCollateral,
		_,
		bigMaintenanceMargin,
		err := k.subaccountsKeeper.GetNetCollateralAndMarginRequirements(
		ctx,
		satypes.Update{SubaccountId: subaccountId},
	)
	if err != nil {
		return false, err
	}

	return CanLiquidateSubaccount(bigNetCollateral, bigMaintenanceMargin), nil
}

// CanLiquidateSubaccount returns true if a subaccount is liquidatable given its total net collateral and
// maintenance margin requirement.
//
// The subaccount is liquidatable if both of the following are true:
// - The maintenance margin requirements are greater than zero (note that they can never be negative).
// - The maintenance margin requirements are greater than the subaccount's net collateral.
//
// Note that this is a stateless function.
func CanLiquidateSubaccount(
	bigNetCollateral *big.Int,
	bigMaintenanceMargin *big.Int,
) bool {
	return bigMaintenanceMargin.Sign() > 0 && bigMaintenanceMargin.Cmp(bigNetCollateral) == 1
}

func GetHealth(
	bigNetCollateral *big.Int,
	bigMaintenanceMargin *big.Int,
) *big.Float {
	// If net collateral is less than 0, return 0
	if bigNetCollateral.Sign() < 0 {
		return big.NewFloat(0)
	}

	// If maintenance margin is less than or equal to 0, return a large number
	if bigMaintenanceMargin.Sign() <= 0 {
		return big.NewFloat(math.MaxFloat64)
	}

	// Calculate the collateral/maintenance margin ratio
	health := new(big.Float).Quo(new(big.Float).SetInt(bigNetCollateral), new(big.Float).SetInt(bigMaintenanceMargin))

	return health
}

func CalculateLiquidationPriority(
	bigTotalNetCollateral *big.Int,
	bigTotalMaintenanceMargin *big.Int,
	bigWeightedMaintenanceMargin *big.Int,
) (
	liquidationPriority *big.Float,
) {

	if bigWeightedMaintenanceMargin.Sign() <= 0 {
		return big.NewFloat(math.MaxFloat64)
	}

	health := GetHealth(bigTotalNetCollateral, bigTotalMaintenanceMargin)
	return new(big.Float).Quo(health, new(big.Float).SetInt(bigWeightedMaintenanceMargin))
}

// EnsureIsLiquidatable returns an error if the subaccount is not liquidatable.
func (k Keeper) EnsureIsLiquidatable(
	ctx sdk.Context,
	subaccountId satypes.SubaccountId,
) (
	err error,
) {
	isLiquidatable, err := k.IsLiquidatable(ctx, subaccountId)
	if err != nil {
		return err
	}
	if !isLiquidatable {
		return errorsmod.Wrapf(
			types.ErrSubaccountNotLiquidatable,
			"SubaccountId %v is not liquidatable",
			subaccountId,
		)
	}
	return nil
}

// getLiquidationPrice returns the liquidation price for a given subaccount and perpetual.
// It calculates the most aggressive price between the bankruptcy price and the fillable price.
func (k Keeper) getLiquidationPrice(
	ctx sdk.Context,
	subaccountId satypes.SubaccountId,
	perpetualId uint32,
	orderQuantums *big.Int,
	isPositionLong bool,
	clobPair types.ClobPair,
) (
	liquidationPrice types.Subticks,
	err error,
) {

	bankruptcyPriceRat, err := k.GetBankruptcyPrice(ctx, subaccountId, perpetualId, orderQuantums)
	if err != nil {
		return 0, err
	}

	fillablePriceRat, err := k.GetFillablePrice(ctx, subaccountId, perpetualId)
	if err != nil {
		return 0, err
	}

	liquidationPriceRat := GetMostAggressivePrice(bankruptcyPriceRat, fillablePriceRat, isPositionLong)
	return k.ConvertLiquidationPriceToSubticks(ctx, liquidationPriceRat, isPositionLong, clobPair), nil
}

// For a long position (isLong == true), it returns the lower price.
// For a short position (isLong == false), it returns the higher price.
func GetMostAggressivePrice(bankruptcyPriceRat *big.Rat, fillablePriceRat *big.Rat, isLong bool) (liquidationPrice *big.Rat) {
	if (isLong && bankruptcyPriceRat.Cmp(fillablePriceRat) < 0) || (!isLong && bankruptcyPriceRat.Cmp(fillablePriceRat) > 0) {
		return bankruptcyPriceRat
	}
	return fillablePriceRat
}

// GetFillablePrice returns the fillable-price of a subaccount’s position. It returns a rational
// number to avoid rounding errors.
//
// The equation for calculating the fillable price is the following:
// `(PNNV - ABR * SMMR * PMMR) / PS`, where `ABR = BA * (1 - (TNC / TMMR))`.
// To calculate this, we must first fetch the following values:
//   - PS (The perpetual position size held by the subaccount, used for calculating the
//     position net notional value and maintenance margin requirement).
//   - PNNV (position net notional value).
//   - PMMR (position maintenance margin requirement).
//   - TNC (total net collateral).
//   - TMMR (total maintenance margin requirement).
//   - BA (bankruptcy adjustment PPM).
//   - SMMR (spread to maintenance margin ratio)
func (k Keeper) GetFillablePrice(
	ctx sdk.Context,
	subaccountId satypes.SubaccountId,
	perpetualId uint32,
) (
	fillablePrice *big.Rat,
	err error,
) {

	pnnvBig, pmmrBig, tncBig, tmmrBig, ba, smmr, bigPositionSizeQuantums, err := k.getFillablePriceCalculationInputs(ctx, subaccountId, perpetualId)
	if err != nil {
		return nil, err
	}

	fillablePrice = calculateFillablePrice(
		pnnvBig,
		pmmrBig,
		tncBig,
		tmmrBig,
		ba,
		smmr,
		bigPositionSizeQuantums,
	)

	if fillablePrice.Sign() < 0 {
		panic("GetFillablePrice: Calculated fillable price is negative")
	}

	return fillablePrice, nil
}

// calculateFillablePrice calculates the fillable price for a liquidation order.
// It uses the formula: (PNNV - ABR * SMMR * PMMR) / PS
func calculateFillablePrice(
	pnnvBig *big.Int,
	pmmrBig *big.Int,
	tncBig *big.Int,
	tmmrBig *big.Int,
	ba uint32,
	smmr uint32,
	bigPositionSizeQuantums *big.Int,
) (
	fillablePrice *big.Rat,
) {
	// Calculate ABR (Adjusted Bankruptcy Rating)
	adjustedBankruptcyRating := calculateAdjustedBankruptcyRating(tncBig, tmmrBig, ba)

	// Calculate SMMR * PMMR
	maxLiquidationSpreadQuoteQuantumsRat := lib.BigRatMulPpm(new(big.Rat).SetInt(pmmrBig), smmr)

	// Calculate ABR * SMMR * PMMR
	fillablePriceOracleDeltaQuoteQuantumsRat := new(big.Rat).Mul(adjustedBankruptcyRating, maxLiquidationSpreadQuoteQuantumsRat)

	// Calculate PNNV - (ABR * SMMR * PMMR)
	pnnvRat := new(big.Rat).SetInt(pnnvBig)
	fillablePriceQuoteQuantumsRat := new(big.Rat).Sub(pnnvRat, fillablePriceOracleDeltaQuoteQuantumsRat)

	// Calculate the fillable price by dividing by PS
	return new(big.Rat).Quo(fillablePriceQuoteQuantumsRat, new(big.Rat).SetInt(bigPositionSizeQuantums))
}

// It uses the formula: ABR = BA * (1 - (TNC / TMMR))
func calculateAdjustedBankruptcyRating(
	tncBig *big.Int,
	tmmrBig *big.Int,
	ba uint32,
) (
	adjustedBankruptcyRating *big.Rat,
) {
	// Calculate TNC / TMMR
	tncDivTmmrRat := new(big.Rat).SetFrac(tncBig, tmmrBig)

	// Calculate 1 - (TNC / TMMR)
	oneMinusTncDivTmmrRat := new(big.Rat).Sub(lib.BigRat1(), tncDivTmmrRat)

	// Calculate (1 - (TNC / TMMR)) * BA
	unboundedAbrRat := lib.BigRatMulPpm(oneMinusTncDivTmmrRat, ba)

	// Bound the ABR between 0 and 1
	return lib.BigRatClamp(unboundedAbrRat, lib.BigRat0(), lib.BigRat1())
}

// getFillablePriceCalculationInputs retrieves and calculates the necessary inputs for the fillable price calculation.
func (k Keeper) getFillablePriceCalculationInputs(
	ctx sdk.Context,
	subaccountId satypes.SubaccountId,
	perpetualId uint32,
) (
	pnnvBig *big.Int,
	pmmrBig *big.Int,
	tncBig *big.Int,
	tmmrBig *big.Int,
	ba uint32,
	smmr uint32,
	bigPositionSizeQuantums *big.Int,
	err error,
) {
	bigPositionSizeQuantums = k.getPositionSize(ctx, subaccountId, perpetualId)

	if bigPositionSizeQuantums.Sign() == 0 {
		return nil, nil, nil, nil, 0, 0, nil, types.ErrInvalidPerpetualPositionSizeDelta
	}

	pnnvBig, err = k.perpetualsKeeper.GetNetCollateral(ctx, perpetualId, bigPositionSizeQuantums)
	if err != nil {
		return nil, nil, nil, nil, 0, 0, nil, err
	}

	_, pmmrBig, err = k.perpetualsKeeper.GetMarginRequirements(ctx, perpetualId, bigPositionSizeQuantums)
	if err != nil {
		return nil, nil, nil, nil, 0, 0, nil, err
	}

	tncBig, _, tmmrBig, err = k.subaccountsKeeper.GetNetCollateralAndMarginRequirements(ctx, satypes.Update{SubaccountId: subaccountId})
	if err != nil {
		return nil, nil, nil, nil, 0, 0, nil, err
	}

	// Emit metrics for negative TNC
	if tncBig.Sign() < 0 {
		k.handleMetricsForNegativeTNC(ctx, subaccountId, perpetualId, tncBig)
	}

	liquidationsConfig := k.GetLiquidationsConfig(ctx)
	ba = liquidationsConfig.FillablePriceConfig.BankruptcyAdjustmentPpm
	smmr = liquidationsConfig.FillablePriceConfig.SpreadToMaintenanceMarginRatioPpm

	return pnnvBig, pmmrBig, tncBig, tmmrBig, ba, smmr, bigPositionSizeQuantums, nil
}

// handleMetricsForNegativeTNC handles metrics and logging for cases where a subaccount has negative TNC.
func (k Keeper) handleMetricsForNegativeTNC(
	ctx sdk.Context,
	subaccountId satypes.SubaccountId,
	perpetualId uint32,
	tncBig *big.Int,
) {
	callback := metrics.PrepareCheckState
	if !ctx.IsCheckTx() && !ctx.IsReCheckTx() {
		callback = metrics.DeliverTx
	}

	metrics.IncrCounterWithLabels(
		metrics.LiquidationsLiquidationMatchNegativeTNC,
		1,
		metrics.GetLabelForIntValue(
			metrics.PerpetualId,
			int(perpetualId),
		),
		metrics.GetLabelForStringValue(
			metrics.Callback,
			callback,
		),
	)

	ctx.Logger().Info(
		"GetFillablePrice: Subaccount has negative TNC. SubaccountId: %+v, TNC: %+v",
		subaccountId,
		tncBig,
	)
}

func (k Keeper) getPositionSize(
	ctx sdk.Context,
	subaccountId satypes.SubaccountId,
	perpetualId uint32,
) (
	bigPositionSizeQuantums *big.Int,
) {
	subaccount := k.subaccountsKeeper.GetSubaccount(ctx, subaccountId)
	position, _ := subaccount.GetPerpetualPositionForId(perpetualId)
	return position.GetBigQuantums()
}

// GetBankruptcyPrice calculates the bankruptcy price for a given subaccount's position in a perpetual.
// It returns the bankruptcy price as a big.Rat.
func (k Keeper) GetBankruptcyPrice(
	ctx sdk.Context,
	subaccountId satypes.SubaccountId,
	perpetualId uint32,
	orderQuantums *big.Int,
) (
	bankruptcyPrice *big.Rat,
	err error,
) {
	bankruptcyPriceQuoteQuantums, err := k.GetBankruptcyPriceInQuoteQuantums(ctx, subaccountId, perpetualId, orderQuantums)
	if err != nil {
		return nil, err
	}

	bankruptcyPriceQuoteQuantumsRat := new(big.Rat).SetInt(bankruptcyPriceQuoteQuantums)
	bankruptcyPrice = new(big.Rat).Quo(
		bankruptcyPriceQuoteQuantumsRat,
		new(big.Rat).Neg(new(big.Rat).SetInt(orderQuantums)),
	)

	return bankruptcyPrice, nil
}

// Returns the bankruptcy-price of a subaccount’s position delta in quote quantums.
// Note that the result `deltaQuoteQuantums` is signed and always rounded towards
// positive infinity so that closing the position at the rounded bankruptcy price
// does not require any insurance fund payment.
//
// Also note that this function does not check whether the given subaccount is liquidatable,
// but validates that the provided deltaQuantums is valid with respect to the current position size.
//
// The equation for calculating the bankruptcy price is the following:
// `-DNNV - (TNC * (abs(DMMR) / TMMR))`.
// To calculate this, we must first fetch the following values:
// - DNNV (delta position net notional value).
//   - Note that this is calculated from PNNV (position net notional value) and
//     PNNVAD (position net notional value after delta).
//
// - TNC (total net collateral).
// - DMMR (delta maintenance margin requirement).
// - TMMR (total maintenance margin requirement).
func (k Keeper) GetBankruptcyPriceInQuoteQuantums(
	ctx sdk.Context,
	subaccountId satypes.SubaccountId,
	perpetualId uint32,
	deltaQuantums *big.Int,
) (
	bankruptcyPriceQuoteQuantumsBig *big.Int,
	err error,
) {

	tncBig, tmmrBig, pnnvBig, pnnvadBig, pmmrBig, pmmradBig, err := k.getBankruptcyPriceCalculationInputs(ctx, subaccountId, perpetualId, deltaQuantums)
	if err != nil {
		return nil, err
	}

	bankruptcyPriceQuoteQuantumsBig = calculateBankruptcyPrice(pnnvBig, pnnvadBig, pmmrBig, pmmradBig, tncBig, tmmrBig)

	return bankruptcyPriceQuoteQuantumsBig, nil
}

func calculateBankruptcyPrice(
	pnnvBig *big.Int,
	pnnvadBig *big.Int,
	pmmrBig *big.Int,
	pmmradBig *big.Int,
	tncBig *big.Int,
	tmmrBig *big.Int,
) (
	bankruptcyPriceQuoteQuantumsBig *big.Int,
) {

	// `DNNV = PNNVAD - PNNV`, where `PNNVAD` is the perpetual's net notional
	// with a position size of `PS + deltaQuantums`.
	// Note that we are intentionally not calculating `DNNV` from `deltaQuantums`
	// directly to avoid rounding errors.
	dnnvBig := new(big.Int).Sub(pnnvadBig, pnnvBig)

	// `DMMR = PMMRAD - PMMR`, where `PMMRAD` is the perpetual's maintenance margin requirement
	// with a position size of `PS + deltaQuantums`.
	// Note that we cannot directly calculate `DMMR` from `deltaQuantums` because the maintenance
	// margin requirement function could be non-linear.
	dmmrBig := new(big.Int).Sub(pmmradBig, pmmrBig)
	if dmmrBig.Sign() == 1 {
		panic("calculateBankruptcyPrice: DMMR is positive")
	}

	// Calculate TNC * abs(DMMR) / TMMR
	tncMulDmmrBig := new(big.Int).Mul(tncBig, new(big.Int).Abs(dmmrBig))
	// This calculation is intentionally rounded down to negative infinity to ensure the
	// final result is rounded towards positive-infinity. This works because of the following:
	// - This is the only division in the equation.
	// - This calculation is subtracted from `-DNNV` to get the final result.
	// - The dividend `TNC * abs(DMMR)` is the only number that can be negative, and `Div` uses
	//   Euclidean division so even if `TNC < 0` this will still round towards negative infinity.
	quoteQuantumsBeforeBankruptcyBig := new(big.Int).Div(tncMulDmmrBig, tmmrBig)

	// Calculate -DNNV - (TNC * abs(DMMR) / TMMR)
	return new(big.Int).Sub(new(big.Int).Neg(dnnvBig), quoteQuantumsBeforeBankruptcyBig)
}

// getBankruptcyPriceCalculationInputs retrieves and calculates the necessary inputs for the bankruptcy price calculation.
func (k Keeper) getBankruptcyPriceCalculationInputs(
	ctx sdk.Context,
	subaccountId satypes.SubaccountId,
	perpetualId uint32,
	deltaQuantums *big.Int,
) (
	tncBig *big.Int,
	tmmrBig *big.Int,
	pnnvBig *big.Int,
	pnnvadBig *big.Int,
	pmmrBig *big.Int,
	pmmradBig *big.Int,
	err error,
) {
	tncBig, _, tmmrBig, err = k.subaccountsKeeper.GetNetCollateralAndMarginRequirements(ctx, satypes.Update{SubaccountId: subaccountId})
	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	bigPositionSizeQuantums, err := k.getAndValidatePositionSize(ctx, subaccountId, perpetualId, deltaQuantums)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	pnnvBig, err = k.perpetualsKeeper.GetNetNotional(
		ctx,
		perpetualId,
		bigPositionSizeQuantums,
	)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	pnnvadBig, err = k.perpetualsKeeper.GetNetNotional(
		ctx,
		perpetualId,
		new(big.Int).Add(bigPositionSizeQuantums, deltaQuantums),
	)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	_, pmmrBig, err = k.perpetualsKeeper.GetMarginRequirements(ctx, perpetualId, bigPositionSizeQuantums)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	_, pmmradBig, err = k.perpetualsKeeper.GetMarginRequirements(
		ctx,
		perpetualId,
		new(big.Int).Add(
			bigPositionSizeQuantums,
			deltaQuantums,
		),
	)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	return tncBig, tmmrBig, pnnvBig, pnnvadBig, pmmrBig, pmmradBig, nil
}

// getAndValidatePositionSize retrieves the position size for a given subaccount and perpetual,
// and validates it against the provided delta quantums.
func (k Keeper) getAndValidatePositionSize(
	ctx sdk.Context,
	subaccountId satypes.SubaccountId,
	perpetualId uint32,
	deltaQuantums *big.Int,
) (
	bigPositionSizeQuantums *big.Int,
	err error,
) {
	bigPositionSizeQuantums = k.getPositionSize(ctx, subaccountId, perpetualId)

	// Validate that the provided deltaQuantums is valid with respect to
	// the current position size.
	if bigPositionSizeQuantums.Sign()*deltaQuantums.Sign() != -1 || bigPositionSizeQuantums.CmpAbs(deltaQuantums) == -1 {
		return nil, errorsmod.Wrapf(
			types.ErrInvalidPerpetualPositionSizeDelta,
			"Position size delta %v is invalid for %v and perpetual %v, outstanding position size is %v",
			deltaQuantums,
			subaccountId,
			perpetualId,
			bigPositionSizeQuantums,
		)
	}

	return bigPositionSizeQuantums, nil
}

func (k Keeper) GetLiquidationInsuranceFundFeeAndRemainingAvailableCollateral(
	ctx sdk.Context,
	subaccountId satypes.SubaccountId,
	perpetualId uint32,
	isBuy bool,
	fillAmount uint64,
	subticks types.Subticks,
) (
	remainingQuoteQuantumsBig *big.Int,
	insuranceFundFeeQuoteQuantums *big.Int,
	err error,
) {
	// Verify that fill amount is not zero.
	if fillAmount == 0 {
		return nil, nil, errorsmod.Wrapf(
			types.ErrInvalidQuantumsForInsuranceFundDeltaCalculation,
			"FillAmount is zero for subaccount %v and perpetual %v.",
			subaccountId,
			perpetualId,
		)
	}

	positionChangeQuoteQuantums, bankruptcyPriceInQuoteQuantums, err := k.getPositionChangeAndBankruptcyQuoteQuantums(ctx, perpetualId, subaccountId, isBuy, fillAmount, subticks)
	if err != nil {
		return nil, nil, err
	}

	// Determine the delta between the quote quantums received from closing the position and the
	// bankruptcy price in quote quantums.
	bankruptcyDeltaQuoteQuantums := new(big.Int).Sub(positionChangeQuoteQuantums, bankruptcyPriceInQuoteQuantums)

	// If the insurance fund delta is less than or equal to zero, this means the insurance fund
	// needs to cover the difference between the quote quantums received from closing the position
	// and the bankruptcy price in quote quantums.
	if bankruptcyDeltaQuoteQuantums.Sign() <= 0 {
		return big.NewInt(0), bankruptcyDeltaQuoteQuantums, nil
	}

	insuranceFundFeeQuoteQuantums = k.getMaxInsuranceFundFee(ctx, positionChangeQuoteQuantums, bankruptcyDeltaQuoteQuantums)
	remainingQuoteQuantumsBig = new(big.Int).Sub(bankruptcyDeltaQuoteQuantums, insuranceFundFeeQuoteQuantums)

	return remainingQuoteQuantumsBig, insuranceFundFeeQuoteQuantums, nil

}

func (k Keeper) getMaxInsuranceFundFee(
	ctx sdk.Context,
	positionChangeQuoteQuantums *big.Int,
	bankruptcyDeltaQuoteQuantums *big.Int,
) (
	insuranceFundFeeQuoteQuantums *big.Int,
) {

	// The insurance fund delta is positive. We must read the liquidations config from state to
	// determine the max liquidation fee this user must pay.
	liquidationsConfig := k.GetLiquidationsConfig(ctx)

	// Calculate the max liquidation fee from the magnitude of quote quantums the subaccount
	// will receive from closing this position and the max liquidation fee PPM.
	maxLiquidationFeeQuoteQuantumsBig := lib.BigIntMulPpm(
		new(big.Int).Abs(positionChangeQuoteQuantums),
		liquidationsConfig.InsuranceFundFeePpm,
	)

	// The liquidation fee paid by the user is the minimum of the max liquidation fee and the
	// leftover collateral from liquidating the position.
	return lib.BigMin(
		maxLiquidationFeeQuoteQuantumsBig,
		bankruptcyDeltaQuoteQuantums,
	)
}

func (k Keeper) getPositionChangeAndBankruptcyQuoteQuantums(
	ctx sdk.Context,
	perpetualId uint32,
	subaccountId satypes.SubaccountId,
	isBuy bool,
	fillAmount uint64,
	subticks types.Subticks,
) (
	positionChangeQuoteQuantums *big.Int,
	bankruptcyPriceInQuoteQuantums *big.Int,
	err error,
) {

	clobPair := k.mustGetClobPairForPerpetualId(ctx, perpetualId)
	liquidationOrderQuantums := new(big.Int).SetUint64(fillAmount)
	positionChangeQuoteQuantums, err = getFillQuoteQuantums(clobPair, subticks, satypes.BaseQuantums(fillAmount))
	if err != nil {
		return nil, nil, err
	}

	// If the fill is a buy, the position was a short, so we will remove from the asset struct
	if isBuy {
		positionChangeQuoteQuantums.Neg(positionChangeQuoteQuantums)
	} else {
		liquidationOrderQuantums.Neg(liquidationOrderQuantums)
	}

	bankruptcyPriceInQuoteQuantums, err = k.GetBankruptcyPriceInQuoteQuantums(ctx, subaccountId, perpetualId, liquidationOrderQuantums)
	if err != nil {
		return nil, nil, err
	}

	return positionChangeQuoteQuantums, bankruptcyPriceInQuoteQuantums, nil
}

func (k Keeper) GetValidatorAndLiquidityFee(
	ctx sdk.Context,
	remainingQuoteQuantumsBig *big.Int,
) (
	validatorFeeQuoteQuantums *big.Int,
	liquidityFeeQuoteQuantums *big.Int,
	err error,
) {

	if remainingQuoteQuantumsBig.Cmp(big.NewInt(0)) < 0 {
		return nil, nil, errorsmod.Wrapf(
			types.ErrInvalidQuantumsForInsuranceFundDeltaCalculation,
			"Remaining quote quantums %v is negative",
			remainingQuoteQuantumsBig,
		)
	}

	if remainingQuoteQuantumsBig.Cmp(big.NewInt(0)) == 0 {
		return big.NewInt(0), big.NewInt(0), nil
	}

	liquidationsConfig := k.GetLiquidationsConfig(ctx)
	validatorFeeQuoteQuantums = lib.BigIntMulPpm(remainingQuoteQuantumsBig, liquidationsConfig.ValidatorFeePpm)
	liquidityFeeQuoteQuantums = lib.BigIntMulPpm(remainingQuoteQuantumsBig, liquidationsConfig.LiquidityFeePpm)

	err = k.validateValidatorAndLiquidityFee(remainingQuoteQuantumsBig, validatorFeeQuoteQuantums, liquidityFeeQuoteQuantums)
	if err != nil {
		return nil, nil, err
	}

	return validatorFeeQuoteQuantums, liquidityFeeQuoteQuantums, nil

}

func (k Keeper) validateValidatorAndLiquidityFee(
	remainingQuoteQuantumsBig *big.Int,
	validatorFeeQuoteQuantums *big.Int,
	liquidityFeeQuoteQuantums *big.Int,
) error {

	totalFees := new(big.Int).Add(validatorFeeQuoteQuantums, liquidityFeeQuoteQuantums)
	if totalFees.Cmp(remainingQuoteQuantumsBig) > 0 {
		return errorsmod.Wrapf(
			types.ErrInvalidQuantumsForInsuranceFundDeltaCalculation,
			"Total fees %v exceed remaining quote quantums %v",
			totalFees,
			remainingQuoteQuantumsBig,
		)
	}
	return nil
}

func (k Keeper) GetBestPerpetualPositionToLiquidate(
	ctx sdk.Context,
	subaccountId satypes.SubaccountId,
) (
	perpetualId uint32,
	err error,
) {

	subaccount := k.subaccountsKeeper.GetSubaccount(ctx, subaccountId)
	subaccountLiquidationInfo := k.GetSubaccountLiquidationInfo(ctx, subaccountId)

	bestPriority := big.NewFloat(-1)
	bestPerpetualId := uint32(0)

	if len(subaccount.PerpetualPositions) == 1 {
		if subaccountLiquidationInfo.HasPerpetualBeenLiquidatedForSubaccount(subaccount.PerpetualPositions[0].PerpetualId) {
			return 0, types.ErrNoPerpetualPositionsToLiquidate
		} else {
			return subaccount.PerpetualPositions[0].PerpetualId, nil
		}
	}

	for _, position := range subaccount.PerpetualPositions {
		err := k.SimulatePriorityWithClosedPosition(ctx, subaccount, subaccountLiquidationInfo, position, bestPriority, &bestPerpetualId)
		if err != nil {
			return 0, err
		}
	}

	if bestPriority.Sign() >= 0 {
		return bestPerpetualId, nil
	}
	return 0, types.ErrNoPerpetualPositionsToLiquidate
}

func (k Keeper) SimulatePriorityWithClosedPosition(
	ctx sdk.Context,
	subaccount satypes.Subaccount,
	subaccountLiquidationInfo types.SubaccountLiquidationInfo,
	position *satypes.PerpetualPosition,
	bestPriority *big.Float,
	bestPerpetualId *uint32,
) error {

	perpetual, err := k.perpetualsKeeper.GetPerpetual(ctx, position.PerpetualId)
	if err != nil {
		return err
	}
	price, err := k.pricesKeeper.GetMarketPrice(ctx, perpetual.Params.MarketId)
	if err != nil {
		return err
	}

	if !subaccountLiquidationInfo.HasPerpetualBeenLiquidatedForSubaccount(position.PerpetualId) {
		err := k.simulatePriorityWithClosedPosition(ctx, subaccount, position, price, bestPriority, bestPerpetualId)
		if err != nil {
			return err
		}
	}
	return nil
}

func deepCopySubaccount(subaccount satypes.Subaccount) satypes.Subaccount {

	copySubaccount := satypes.Subaccount{
		Id:              subaccount.Id,
		MarginEnabled:   subaccount.MarginEnabled,
		AssetYieldIndex: subaccount.AssetYieldIndex,
	}

	// Deep copy AssetPositions if not nil
	if subaccount.AssetPositions != nil {
		copySubaccount.AssetPositions = make([]*satypes.AssetPosition, len(subaccount.AssetPositions))
		for i, ap := range subaccount.AssetPositions {
			newAp := *ap // Dereference and copy the AssetPosition
			copySubaccount.AssetPositions[i] = &newAp
		}
	}

	// Deep copy PerpetualPositions if not nil
	if subaccount.PerpetualPositions != nil {
		copySubaccount.PerpetualPositions = make([]*satypes.PerpetualPosition, len(subaccount.PerpetualPositions))
		for i, pp := range subaccount.PerpetualPositions {
			newPp := *pp // Dereference and copy the PerpetualPosition
			copySubaccount.PerpetualPositions[i] = &newPp
		}
	}

	return copySubaccount
}

func (k Keeper) simulatePriorityWithClosedPosition(
	ctx sdk.Context,
	subaccount satypes.Subaccount,
	position *satypes.PerpetualPosition,
	price pricestypes.MarketPrice,
	bestPriority *big.Float,
	bestPerpetualId *uint32,
) error {

	closedSubaccount := deepCopySubaccount(subaccount)

	closedSubaccount, err := k.SimulateClosePerpetualPosition(ctx, closedSubaccount, position, price)
	if err != nil {
		return err
	}

	_, priority, err := k.GetSubaccountPriority(ctx, closedSubaccount)
	if err != nil {
		return err
	}

	if priority.Cmp(bestPriority) > 0 {
		*bestPriority = *priority
		*bestPerpetualId = position.PerpetualId
	}
	return nil
}

func (k Keeper) SimulateClosePerpetualPosition(
	ctx sdk.Context,
	subaccount satypes.Subaccount,
	position *satypes.PerpetualPosition,
	price pricestypes.MarketPrice,
) (
	closedSubaccount satypes.Subaccount,
	err error,
) {

	RemovePerpetualPosition(&subaccount, position.PerpetualId)

	perpetual, err := k.perpetualsKeeper.GetPerpetual(ctx, position.PerpetualId)
	if err != nil {
		return satypes.Subaccount{}, err
	}
	bigNetCollateralQuoteQuantums := perpkeeper.GetNetNotionalInQuoteQuantums(perpetual, price, position.GetBigQuantums())

	err = UpdateTDaiPosition(&subaccount, bigNetCollateralQuoteQuantums)
	if err != nil {
		return satypes.Subaccount{}, err
	}
	return subaccount, nil
}

func UpdateTDaiPosition(subaccount *satypes.Subaccount, quantumsDelta *big.Int) (err error) {

	assetPosition := subaccount.AssetPositions[0]
	if assetPosition.AssetId != assetstypes.AssetTDai.Id {
		return errors.New("first asset position must be TDai")
	}

	assetPosition.Quantums = dtypes.NewIntFromBigInt(new(big.Int).Add(assetPosition.Quantums.BigInt(), quantumsDelta))

	if assetPosition.Quantums.BigInt().Sign() == 0 {
		subaccount.AssetPositions = []*satypes.AssetPosition{}
	} else {
		subaccount.AssetPositions = []*satypes.AssetPosition{assetPosition}
	}
	return nil
}

func RemovePerpetualPosition(subaccount *satypes.Subaccount, perpetualId uint32) {
	for i, pos := range subaccount.PerpetualPositions {
		if pos.PerpetualId == perpetualId {
			// Remove the position from the slice
			subaccount.PerpetualPositions = append(
				subaccount.PerpetualPositions[:i],
				subaccount.PerpetualPositions[i+1:]...,
			)
			return
		}
	}
}

// GetNegativePositionSize returns the number of base quantums needed to liquidate
func (k Keeper) GetNegativePositionSize(
	ctx sdk.Context,
	subaccountId satypes.SubaccountId,
	perpetualId uint32,
) (
	deltaQuantums *big.Int,
	err error,
) {
	subaccount := k.subaccountsKeeper.GetSubaccount(ctx, subaccountId)
	perpetualPosition, exists := subaccount.GetPerpetualPositionForId(perpetualId)
	if !exists {
		return nil,
			errorsmod.Wrapf(
				types.ErrNoPerpetualPositionsToLiquidate,
				"SubaccountId: %v, perpetualId: %d",
				subaccount.Id,
				perpetualId,
			)
	}

	return new(big.Int).Neg(perpetualPosition.GetBigQuantums()), nil
}

func (k Keeper) GetMaxQuantumsInsuranceDelta(
	ctx sdk.Context,
	perpetualId uint32,
) (
	bigMaxQuantumsInsuranceLost *big.Int,
	err error,
) {

	bigInsuranceFundLostBlockLimit, err := k.GetInsuranceFundDeltaBlockLimit(ctx, perpetualId)
	if err != nil {
		return nil, err
	}
	bigCurrentInsuranceFundLost, err := k.GetCumulativeInsuranceFundDelta(ctx, perpetualId)
	if err != nil {
		return nil, err
	}

	if bigCurrentInsuranceFundLost.Cmp(bigInsuranceFundLostBlockLimit) > 0 {
		return nil, errorsmod.Wrapf(
			types.ErrLiquidationExceedsMaxInsuranceLost,
			"Insurance lost exceeds block limit. Current insurance lost: %v, block limit: %v",
			bigCurrentInsuranceFundLost,
			bigInsuranceFundLostBlockLimit,
		)
	}

	return new(big.Int).Sub(bigInsuranceFundLostBlockLimit, bigCurrentInsuranceFundLost), nil
}

func (k Keeper) GetInsuranceFundDeltaBlockLimit(ctx sdk.Context, perpetualId uint32) (*big.Int, error) {
	isIsolated, err := k.perpetualsKeeper.IsIsolatedPerpetual(ctx, perpetualId)
	if err != nil {
		return big.NewInt(0), err
	}

	if isIsolated {
		perpetual, _ := k.perpetualsKeeper.GetPerpetual(ctx, perpetualId)
		return new(big.Int).SetUint64(perpetual.Params.IsolatedMarketMaxCumulativeInsuranceFundDeltaPerBlock), nil
	}

	return new(big.Int).SetUint64(k.GetLiquidationsConfig(ctx).MaxCumulativeInsuranceFundDelta), nil
}

// ConvertLiquidationPriceToSubticks converts the liquidation price of a liquidation order to subticks.
// The returned subticks will be rounded to the nearest tick (such that
// `subticks % clobPair.SubticksPerTick == 0`). This function will round up for sells
// that close longs, and round down for buys that close shorts.
//
// Note the returned `subticks` will be bounded (inclusive) between `clobPair.SubticksPerTick` and
// `math.MaxUint64 - math.MaxUint64 % clobPair.SubticksPerTick` (the maximum `uint64` that is a
// multiple of `clobPair.SubticksPerTick`).
//
// If we are liquidating a long position with a sell order, then we round up to the nearest
// subtick (and vice versa for liquidating shorts).
func (k Keeper) ConvertLiquidationPriceToSubticks(
	ctx sdk.Context,
	liquidationPrice *big.Rat,
	isLiquidatingLong bool,
	clobPair types.ClobPair,
) (
	subticks types.Subticks,
) {
	// The liquidation price is invalid if it is negative.
	if liquidationPrice.Sign() < 0 {
		panic("ConvertLiquidationPriceToSubticks: liquidationPrice should not be negative")
	}

	// Assuming `liquidationPrice` is in units of `quote quantums / base quantum`,  then dividing by
	// `quote quantums / (base quantum * subtick)` will give the resulting units of subticks.
	subticksRat := new(big.Rat).Quo(
		liquidationPrice,
		getQuoteQuantumsPerBaseQuantumAndSubtickRat(clobPair),
	)

	// Round the subticks to the nearest `big.Int` in the correct direction.
	roundedSubticksBig := lib.BigRatRound(subticksRat, isLiquidatingLong)

	// Ensure `roundedSubticksBig % clobPair.SubticksPerTick == 0`, rounding in the correct
	// direction if necessary.
	roundedAlignedSubticksBig := lib.BigIntRoundToMultiple(
		roundedSubticksBig,
		new(big.Int).SetUint64(uint64(clobPair.SubticksPerTick)),
		isLiquidatingLong,
	)
	return boundSubticks(roundedAlignedSubticksBig, clobPair)
}

func boundSubticks(subticksBig *big.Int, clobPair types.ClobPair) types.Subticks {

	// Bound the result between `clobPair.SubticksPerTick` and
	// `math.MaxUint64 - math.MaxUint64 % clobPair.SubticksPerTick`.
	minSubticks := uint64(clobPair.SubticksPerTick)
	maxSubticks := uint64(math.MaxUint64 - (math.MaxUint64 % uint64(clobPair.SubticksPerTick)))
	boundedSubticks := lib.BigUint64Clamp(
		subticksBig,
		minSubticks,
		maxSubticks,
	)

	// Panic if the bounded subticks is zero or is not a multiple of `clobPair.SubticksPerTick`,
	// which would indicate the rounding or clamp logic failed.
	if boundedSubticks == 0 {
		panic("boundSubticks: Bounded subticks is 0.")
	} else if boundedSubticks%uint64(clobPair.SubticksPerTick) != 0 {
		panic("boundSubticks: Bounded subticks is not a multiple of SubticksPerTick.")
	}
	return types.Subticks(boundedSubticks)
}

func getQuoteQuantumsPerBaseQuantumAndSubtickRat(
	clobPair types.ClobPair,
) (
	quoteQuantumsPerBaseQuantumAndSubtickRat *big.Rat,
) {
	exponent := clobPair.QuantumConversionExponent
	absExponentiatedValueBig := lib.BigPow10(uint64(lib.AbsInt32(exponent)))
	quoteQuantumsPerBaseQuantumAndSubtickRat = new(big.Rat).SetInt(absExponentiatedValueBig)
	// If `exponent` is negative, invert the fraction to set the result to `1 / 10^exponent`.
	if exponent < 0 {
		quoteQuantumsPerBaseQuantumAndSubtickRat.Inv(quoteQuantumsPerBaseQuantumAndSubtickRat)
	}
	return quoteQuantumsPerBaseQuantumAndSubtickRat
}

func (k Keeper) validateMatchedLiquidationAndGetFees(
	ctx sdk.Context,
	order types.MatchableOrder,
	perpetualId uint32,
	fillAmount satypes.BaseQuantums,
	makerSubticks types.Subticks,
) (
	insuranceFundDelta *big.Int,
	validatorFeeQuoteQuantums *big.Int,
	liquidityFeeQuoteQuantums *big.Int,
	err error,
) {

	remainingQuoteQuantumsBig, insuranceFundDelta, err := k.GetLiquidationInsuranceFundFeeAndRemainingAvailableCollateral(ctx, order.GetSubaccountId(), perpetualId, order.IsBuy(), fillAmount.ToUint64(), makerSubticks)
	if err != nil {
		return nil, nil, nil, err
	}

	if err := k.validateLiquidationParams(ctx, order.GetSubaccountId(), perpetualId, insuranceFundDelta); err != nil {
		return nil, nil, nil, err
	}

	validatorFeeQuoteQuantums, liquidityFeeQuoteQuantums, err = k.GetValidatorAndLiquidityFee(ctx, remainingQuoteQuantumsBig)
	if err != nil {
		return nil, nil, nil, err
	}

	return insuranceFundDelta, validatorFeeQuoteQuantums, liquidityFeeQuoteQuantums, nil
}

func (k Keeper) validateLiquidationParams(
	ctx sdk.Context,
	subaccountId satypes.SubaccountId,
	perpetualId uint32,
	insuranceFundDelta *big.Int,
) (
	err error,
) {

	err = k.EnsurePerpetualNotAlreadyLiquidated(ctx, subaccountId, perpetualId)
	if err != nil {
		return err
	}

	err = k.CheckInsuranceFundLimits(ctx, perpetualId, insuranceFundDelta)
	if err != nil {
		return err
	}

	return k.verifyInsuranceFundHasSufficientBalance(ctx, perpetualId, insuranceFundDelta)
}

func (k Keeper) EnsurePerpetualNotAlreadyLiquidated(
	ctx sdk.Context,
	subaccountId satypes.SubaccountId,
	perpetualId uint32,
) error {
	subaccountLiquidationInfo := k.GetSubaccountLiquidationInfo(ctx, subaccountId)
	if subaccountLiquidationInfo.HasPerpetualBeenLiquidatedForSubaccount(perpetualId) {
		return errorsmod.Wrapf(
			types.ErrSubaccountHasLiquidatedPerpetual,
			"Subaccount %v and perpetual %v have already been liquidated within the last block",
			subaccountId,
			perpetualId,
		)
	}
	return nil
}

func (k Keeper) CheckInsuranceFundLimits(
	ctx sdk.Context,
	perpetualId uint32,
	insuranceFundDelta *big.Int,
) error {
	if insuranceFundDelta.Sign() == -1 {

		bigMaxQuantumsInsuranceLost, err := k.GetMaxQuantumsInsuranceDelta(ctx, perpetualId)
		if err != nil {
			return err
		}
		if insuranceFundDelta.CmpAbs(bigMaxQuantumsInsuranceLost) > 0 {
			return errorsmod.Wrapf(
				types.ErrLiquidationExceedsMaxInsuranceLost,
				"Max Insurance Lost: %v, Insurance Lost: %v",
				bigMaxQuantumsInsuranceLost,
				insuranceFundDelta,
			)
		}
	}
	return nil
}

func (k Keeper) verifyInsuranceFundHasSufficientBalance(
	ctx sdk.Context,
	perpetualId uint32,
	insuranceFundDelta *big.Int,
) error {

	if !k.IsValidInsuranceFundDelta(ctx, insuranceFundDelta, perpetualId) {
		log.DebugLog(ctx, "ProcessMatches: insurance fund has insufficient balance to process the liquidation.")
		return errorsmod.Wrapf(
			types.ErrInsuranceFundHasInsufficientFunds,
			"Insurance fund delta %v",
			insuranceFundDelta.String(),
		)
	}

	return nil
}
