package keeper

import (
	"bytes"
	"errors"
	"math"
	"math/big"
	"sort"
	"time"

	errorsmod "cosmossdk.io/errors"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/lib/log"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	perplib "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/lib"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

// subaccountToDeleverage is a struct containing a subaccount ID and perpetual ID to deleverage.
// This struct is used as a return type for the LiquidateSubaccountsAgainstOrderbook and
// GetSubaccountsWithOpenPositionsInFinalSettlementMarkets called in PrepareCheckState.
type subaccountToDeleverage struct {
	SubaccountId satypes.SubaccountId
	PerpetualId  uint32
}

// LiquidateSubaccountsAgainstOrderbook takes a list of subaccount IDs and liquidates them against
// the orderbook. It will liquidate as many subaccounts as possible up to the maximum number of
// liquidations per block. Subaccounts are selected with a pseudo-randomly generated offset. A slice
// of subaccounts to deleverage is returned from this function, derived from liquidation orders that
// failed to fill.
func (k Keeper) LiquidateSubaccountsAgainstOrderbook(
	ctx sdk.Context,
	subaccountIds []satypes.SubaccountId,
) (
	subaccountsToDeleverage []subaccountToDeleverage,
	err error,
) {
	lib.AssertCheckTxMode(ctx)

	metrics.AddSample(
		metrics.LiquidationsLiquidatableSubaccountIdsCount,
		float32(len(subaccountIds)),
	)

	// Early return if there are 0 subaccounts to liquidate.
	numSubaccounts := len(subaccountIds)
	if numSubaccounts == 0 {
		return nil, nil
	}

	defer telemetry.MeasureSince(
		time.Now(),
		types.ModuleName,
		metrics.ClobLiquidateSubaccountsAgainstOrderbook,
		metrics.Latency,
	)

	// Get the liquidation order for each subaccount.
	// Process at-most `MaxLiquidationAttemptsPerBlock` subaccounts, starting from a pseudorandom location
	// in the slice. Note `numSubaccounts` is guaranteed to be non-zero at this point, so `Intn` shouldn't panic.
	pseudoRand := k.GetPseudoRand(ctx)
	liquidationOrders := make([]types.LiquidationOrder, 0)
	numLiqOrders := lib.Min(numSubaccounts, int(k.Flags.MaxLiquidationAttemptsPerBlock))
	indexOffset := pseudoRand.Intn(numSubaccounts)

	startGetLiquidationOrders := time.Now()
	for i := 0; i < numLiqOrders; i++ {
		index := (i + indexOffset) % numSubaccounts
		subaccountId := subaccountIds[index]
		liquidationOrder, err := k.MaybeGetLiquidationOrder(ctx, subaccountId)
		if err != nil {
			// Subaccount might not always be liquidatable since liquidation daemon runs
			// in a separate goroutine and is not always in sync with the application.
			// Therefore, if subaccount is not liquidatable, continue.
			if errors.Is(err, types.ErrSubaccountNotLiquidatable) {
				telemetry.IncrCounter(1, metrics.MaybeGetLiquidationOrder, metrics.SubaccountsNotLiquidatable, metrics.Count)
				continue
			}

			// Return unexpected errors.
			return nil, err
		}

		liquidationOrders = append(liquidationOrders, *liquidationOrder)
	}
	telemetry.MeasureSince(
		startGetLiquidationOrders,
		types.ModuleName,
		metrics.LiquidateSubaccounts_GetLiquidations,
		metrics.Latency,
	)

	// Sort liquidation orders. The most underwater accounts should be liquidated first.
	// These orders are only used for sorting. When we match these orders here in PrepareCheckState,
	// liquidation matches will be put into the Operations Queue. However, when we process liquidations,
	// we will generate a new liquidation order for each subaccount because previous liquidation orders
	// can alter quantity sizes of subsequent liquidation orders.
	k.SortLiquidationOrders(ctx, liquidationOrders)

	subaccountIdsToLiquidate := lib.MapSlice(liquidationOrders, func(order types.LiquidationOrder) satypes.SubaccountId {
		return order.GetSubaccountId()
	})

	// Attempt to place each liquidation order and perform deleveraging if necessary.
	startPlaceLiquidationOrders := time.Now()
	for _, subaccountId := range subaccountIdsToLiquidate {
		// Generate a new liquidation order with the appropriate order size from the sorted subaccount ids.
		liquidationOrder, err := k.MaybeGetLiquidationOrder(ctx, subaccountId)
		if err != nil {
			// Subaccount might not always be liquidatable if previous liquidation orders
			// improves the net collateral of this subaccount.
			if errors.Is(err, types.ErrSubaccountNotLiquidatable) {
				continue
			}

			// Return unexpected errors.
			return nil, err
		}

		optimisticallyFilledQuantums, _, err := k.PlacePerpetualLiquidation(ctx, *liquidationOrder)
		// Exception for liquidation which conflicts with clob pair status. This is expected for liquidations generated
		// for subaccounts with open positions in final settlement markets.
		if err != nil && !errors.Is(err, types.ErrLiquidationConflictsWithClobPairStatus) {
			log.ErrorLogWithError(
				ctx,
				"Failed to liquidate subaccount",
				err,
				"liquidationOrder", *liquidationOrder,
			)
			return nil, err
		}

		if optimisticallyFilledQuantums == 0 {
			subaccountsToDeleverage = append(subaccountsToDeleverage, subaccountToDeleverage{
				SubaccountId: liquidationOrder.GetSubaccountId(),
				PerpetualId:  liquidationOrder.MustGetLiquidatedPerpetualId(),
			})
		}
	}
	telemetry.MeasureSince(
		startPlaceLiquidationOrders,
		types.ModuleName,
		metrics.LiquidateSubaccounts_PlaceLiquidations,
		metrics.Latency,
	)

	// Stat the number of subaccounts that require deleveraging.
	metrics.SetGaugeWithLabels(
		metrics.ClobSubaccountsRequiringDeleveragingCount,
		float32(len(subaccountsToDeleverage)),
	)

	return subaccountsToDeleverage, nil
}

// MaybeGetLiquidationOrder takes a subaccount ID and returns a liquidation order that can be used to
// liquidate the subaccount.
// If the subaccount is not currently liquidatable, it will do nothing. This function will return an error if calling
// `IsLiquidatable`, `GetPerpetualPositionToLiquidate` or `GetFillablePrice` returns an error.
func (k Keeper) MaybeGetLiquidationOrder(
	ctx sdk.Context,
	subaccountId satypes.SubaccountId,
) (
	liquidationOrder *types.LiquidationOrder,
	err error,
) {
	// If the subaccount is not liquidatable, do nothing.
	if err := k.EnsureIsLiquidatable(ctx, subaccountId); err != nil {
		return nil, err
	}

	defer telemetry.ModuleMeasureSince(
		types.ModuleName,
		time.Now(),
		metrics.ConstructLiquidationOrder,
	)

	// The subaccount is liquidatable. Get the perpetual position and position size to liquidate.
	perpetualId, err := k.GetPerpetualPositionToLiquidate(ctx, subaccountId)
	if err != nil {
		return nil, err
	}

	return k.GetLiquidationOrderForPerpetual(
		ctx,
		subaccountId,
		perpetualId,
	)
}

// GetLiquidationOrderForPerpetual returns a liquidation order for a subaccount
// given a perpetual ID.
func (k Keeper) GetLiquidationOrderForPerpetual(
	ctx sdk.Context,
	subaccountId satypes.SubaccountId,
	perpetualId uint32,
) (
	liquidationOrder *types.LiquidationOrder,
	err error,
) {
	deltaQuantums, err := k.GetLiquidatablePositionSizeDelta(
		ctx,
		subaccountId,
		perpetualId,
	)
	if err != nil {
		return nil, err
	}

	// Get the fillable price of the liquidation order in subticks.
	fillablePriceRat, err := k.GetFillablePrice(ctx, subaccountId, perpetualId, deltaQuantums)
	if err != nil {
		return nil, err
	}

	// Calculate the fillable price.
	isLiquidatingLong := deltaQuantums.Sign() == -1
	clobPair := k.mustGetClobPairForPerpetualId(ctx, perpetualId)
	fillablePriceSubticks := k.ConvertFillablePriceToSubticks(
		ctx,
		fillablePriceRat,
		isLiquidatingLong,
		clobPair,
	)

	// Create the liquidation order.
	absBaseQuantums := deltaQuantums.Abs(deltaQuantums)
	liquidationOrder = types.NewLiquidationOrder(
		subaccountId,
		clobPair,
		!isLiquidatingLong,
		satypes.BaseQuantums(absBaseQuantums.Uint64()),
		fillablePriceSubticks,
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

	orderSizeOptimisticallyFilledFromMatchingQuantums,
		orderStatus,
		offchainUpdates,
		err :=
		k.MemClob.PlacePerpetualLiquidation(
			ctx,
			liquidationOrder,
		)
	if err != nil {
		return 0, 0, err
	}

	// TODO(DEC-1323): Potentially allow liquidating the same perpetual + subaccount
	// multiple times in a block.
	perpetualId := liquidationOrder.MustGetLiquidatedPerpetualId()
	k.MustUpdateSubaccountPerpetualLiquidated(
		ctx,
		liquidationOrder.GetSubaccountId(),
		perpetualId,
	)

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

	k.SendOffchainMessages(offchainUpdates, nil, metrics.SendPlacePerpetualLiquidationOffchainUpdates)
	return orderSizeOptimisticallyFilledFromMatchingQuantums, orderStatus, err
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
	risk, err := k.subaccountsKeeper.GetNetCollateralAndMarginRequirements(
		ctx,
		satypes.Update{SubaccountId: subaccountId},
	)
	if err != nil {
		return false, err
	}

	return risk.IsLiquidatable(), nil
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

// Returns the bankruptcy-price of a subaccount’s position delta in quote quantums.
// Note that the result `deltaQuoteQuantums` is signed and always rounded towards
// positive infinity so that closing the position at the rounded bankruptcy price
// does not require any insurance fund payment.
//
// Also note that this function does not check whether the given subaccount is liquidatable,
// but validates that the provided deltaQuantums is valid with respect to the current position size.
func (k Keeper) GetBankruptcyPriceInQuoteQuantums(
	ctx sdk.Context,
	subaccountId satypes.SubaccountId,
	perpetualId uint32,
	deltaQuantums *big.Int,
) (
	deltaQuoteQuantums *big.Int,
	err error,
) {
	// The equation for calculating the bankruptcy price is the following:
	// `-DNNV - (TNC * (abs(DMMR) / TMMR))`.
	// To calculate this, we must first fetch the following values:
	// - DNNV (delta position net notional value).
	//   - Note that this is calculated from PNNV (position net notional value) and
	//     PNNVAD (position net notional value after delta).
	// - TNC (total net collateral).
	// - DMMR (delta maintenance margin requirement).
	// - TMMR (total maintenance margin requirement).

	riskTotal, err := k.subaccountsKeeper.GetNetCollateralAndMarginRequirements(
		ctx,
		satypes.Update{SubaccountId: subaccountId},
	)
	if err != nil {
		return nil, err
	}

	// Position size is necessary for calculating DNNV and DMMR.
	subaccount := k.subaccountsKeeper.GetSubaccount(ctx, subaccountId)
	position, _ := subaccount.GetPerpetualPositionForId(perpetualId)
	psBig := position.GetBigQuantums()

	// Validate that the provided deltaQuantums is valid with respect to
	// the current position size.
	if psBig.Sign()*deltaQuantums.Sign() != -1 || psBig.CmpAbs(deltaQuantums) == -1 {
		return nil, errorsmod.Wrapf(
			types.ErrInvalidPerpetualPositionSizeDelta,
			"Position size delta %v is invalid for %v and perpetual %v, outstanding position size is %v",
			deltaQuantums,
			subaccountId,
			perpetualId,
			psBig,
		)
	}

	perpetual,
		marketPrice,
		liquidityTier,
		err := k.perpetualsKeeper.
		GetPerpetualAndMarketPriceAndLiquidityTier(ctx, perpetualId)
	if err != nil {
		return nil, err
	}

	// `DNNV = PNNVAD - PNNV`, where `PNNVAD` is the perpetual's net notional
	// with a position size of `PS + deltaQuantums`.
	// Note that we are intentionally not calculating `DNNV` from `deltaQuantums`
	// directly to avoid rounding errors.
	riskPosOld := perplib.GetPositionNetNotionalValueAndMarginRequirements(
		perpetual,
		marketPrice,
		liquidityTier,
		psBig,
		0, // No custom IMF for liquidations
	)
	riskPosNew := perplib.GetPositionNetNotionalValueAndMarginRequirements(
		perpetual,
		marketPrice,
		liquidityTier,
		new(big.Int).Add(psBig, deltaQuantums),
		0, // No custom IMF for liquidations
	)
	// `DMMR = PMMRAD - PMMR`, where `PMMRAD` is the perpetual's maintenance margin requirement
	// with a position size of `PS + deltaQuantums`.
	// Note that we cannot directly calculate `DMMR` from `deltaQuantums` because the maintenance
	// margin requirement function could be non-linear.
	deltaNC := new(big.Int).Sub(riskPosNew.NC, riskPosOld.NC)
	deltaMMR := new(big.Int).Sub(riskPosNew.MMR, riskPosOld.MMR)
	// `deltaMMR` should never be positive if `| PS | >= | PS + deltaQuantums |`. If it is, panic.
	if deltaMMR.Sign() == 1 {
		panic("GetBankruptcyPriceInQuoteQuantums: DMMR is positive")
	}

	// Calculate `TNC * abs(DMMR) / TMMR`.
	tncMulDmmrBig := new(big.Int).Mul(riskTotal.NC, new(big.Int).Abs(deltaMMR))
	// This calculation is intentionally rounded down to negative infinity to ensure the
	// final result is rounded towards positive-infinity. This works because of the following:
	// - This is the only division in the equation.
	// - This calculation is subtracted from `-DNNV` to get the final result.
	// - The dividend `TNC * abs(DMMR)` is the only number that can be negative, and `Div` uses
	//   Euclidean division so even if `TNC < 0` this will still round towards negative infinity.
	quoteQuantumsBeforeBankruptcyBig := new(big.Int).Div(tncMulDmmrBig, riskTotal.MMR)

	// Calculate `-DNNV - TNC * abs(DMMR) / TMMR`.
	bankruptcyPriceQuoteQuantumsBig := new(big.Int).Sub(
		new(big.Int).Neg(deltaNC),
		quoteQuantumsBeforeBankruptcyBig,
	)

	return bankruptcyPriceQuoteQuantumsBig, nil
}

// GetFillablePrice returns the fillable-price of a subaccount’s position. It returns a rational
// number to avoid rounding errors.
func (k Keeper) GetFillablePrice(
	ctx sdk.Context,
	subaccountId satypes.SubaccountId,
	perpetualId uint32,
	deltaQuantums *big.Int,
) (
	fillablePrice *big.Rat,
	err error,
) {
	// The equation for calculating the fillable price is the following:
	// `(PNNV - ABR * SMMR * PMMR) / PS`, where `ABR = BA * (1 - (TNC / TMMR))`.
	// To calculate this, we must first fetch the following values:
	// - PS (The perpetual position size held by the subaccount, used for calculating the
	//   position net notional value and maintenance margin requirement).
	// - PNNV (position net notional value).
	// - PMMR (position maintenance margin requirement).
	// - TNC (total net collateral).
	// - TMMR (total maintenance margin requirement).
	// - BA (bankruptcy adjustment PPM).
	// - SMMR (spread to maintenance margin ratio)

	subaccount := k.subaccountsKeeper.GetSubaccount(ctx, subaccountId)
	position, _ := subaccount.GetPerpetualPositionForId(perpetualId)
	psBig := position.GetBigQuantums()

	// Validate that the provided deltaQuantums is valid with respect to
	// the current position size.
	if psBig.Sign()*deltaQuantums.Sign() != -1 || psBig.CmpAbs(deltaQuantums) == -1 {
		return nil, errorsmod.Wrapf(
			types.ErrInvalidPerpetualPositionSizeDelta,
			"Position size delta %v is invalid for %v and perpetual %v, outstanding position size is %v",
			deltaQuantums,
			subaccountId,
			perpetualId,
			psBig,
		)
	}

	perpetual,
		marketPrice,
		liquidityTier,
		err := k.perpetualsKeeper.GetPerpetualAndMarketPriceAndLiquidityTier(ctx, perpetualId)
	if err != nil {
		return nil, err
	}

	riskPos := perplib.GetPositionNetNotionalValueAndMarginRequirements(
		perpetual,
		marketPrice,
		liquidityTier,
		psBig,
		0, // No custom IMF for liquidations
	)

	riskTotal, err := k.subaccountsKeeper.GetNetCollateralAndMarginRequirements(
		ctx,
		satypes.Update{SubaccountId: subaccountId},
	)
	if err != nil {
		return nil, err
	}

	// stat liquidation order for negative TNC
	// TODO(CLOB-906) Prevent duplicated stat emissions for liquidation orders in PrepareCheckState.
	if riskTotal.NC.Sign() < 0 {
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
			riskTotal.NC,
		)
	}

	liquidationsConfig := k.GetLiquidationsConfig(ctx)
	ba := liquidationsConfig.FillablePriceConfig.BankruptcyAdjustmentPpm
	smmr := liquidationsConfig.FillablePriceConfig.SpreadToMaintenanceMarginRatioPpm

	// Calculate the ABR (adjusted bankruptcy rating).
	tncDivTmmrRat := new(big.Rat).SetFrac(riskTotal.NC, riskTotal.MMR)
	unboundedAbrRat := lib.BigRatMulPpm(
		new(big.Rat).Sub(
			lib.BigRat1(),
			tncDivTmmrRat,
		),
		ba,
	)

	// Bound the ABR between 0 and 1.
	abrRat := lib.BigRatClamp(unboundedAbrRat, lib.BigRat0(), lib.BigRat1())

	// Calculate `SMMR * PMMR` (the maximum liquidation spread in quote quantums).
	maxLiquidationSpreadQuoteQuantumsRat := lib.BigRatMulPpm(
		new(big.Rat).SetInt(riskPos.MMR),
		smmr,
	)

	fillablePriceOracleDeltaQuoteQuantumsRat := new(big.Rat).Mul(abrRat, maxLiquidationSpreadQuoteQuantumsRat)

	// Calculate `PNNV - ABR * SMMR * PMMR`, which represents the fillable price in quote quantums.
	// For longs, `pnnvRat > 0` meaning the fillable price in quote quantums will be lower than the
	// oracle price.
	// For shorts, `pnnvRat < 0` meaning the fillable price in quote quantums will be higher than
	// the oracle price (in this case the result will be negative, but dividing by `positionSize` below
	// will make it positive since `positionSize < 0` for shorts).
	pnnvRat := new(big.Rat).SetInt(riskPos.NC)
	fillablePriceQuoteQuantumsRat := new(big.Rat).Sub(pnnvRat, fillablePriceOracleDeltaQuoteQuantumsRat)

	// Calculate the fillable price by dividing by `PS`.
	// Note that `fillablePriceQuoteQuantumsRat` and `PS` should always have the same sign,
	// meaning the resulting fillable price should always be positive.
	fillablePrice = new(big.Rat).Quo(
		fillablePriceQuoteQuantumsRat,
		new(big.Rat).SetInt(psBig),
	)
	if fillablePrice.Sign() < 0 {
		panic("GetFillablePrice: Calculated fillable price is negative")
	}

	return fillablePrice, nil
}

// GetLiquidationInsuranceFundDelta returns the net payment value between the liquidated account
// and the insurance fund. Positive if the liquidated account pays fees to the insurance fund.
// Negative if the insurance fund covers losses from the subaccount.
func (k Keeper) GetLiquidationInsuranceFundDelta(
	ctx sdk.Context,
	subaccountId satypes.SubaccountId,
	perpetualId uint32,
	isBuy bool,
	fillAmount uint64,
	subticks types.Subticks,
) (
	insuranceFundDeltaQuoteQuantums *big.Int,
	err error,
) {
	// Verify that fill amount is not zero.
	if fillAmount == 0 {
		return nil, errorsmod.Wrapf(
			types.ErrInvalidQuantumsForInsuranceFundDeltaCalculation,
			"FillAmount is zero for subaccount %v and perpetual %v.",
			subaccountId,
			perpetualId,
		)
	}

	// Get the delta quantums and delta quote quantums.
	clobPair := k.mustGetClobPairForPerpetualId(ctx, perpetualId)
	deltaQuantums := new(big.Int).SetUint64(fillAmount)
	deltaQuoteQuantums := types.FillAmountToQuoteQuantums(
		subticks,
		satypes.BaseQuantums(fillAmount),
		clobPair.QuantumConversionExponent,
	)
	if isBuy {
		deltaQuoteQuantums.Neg(deltaQuoteQuantums)
	} else {
		deltaQuantums.Neg(deltaQuantums)
	}

	// To determine the liquidation insurance fund delta we need the bankruptcy price.
	bankruptcyPriceInQuoteQuantumsBig, err := k.GetBankruptcyPriceInQuoteQuantums(
		ctx,
		subaccountId,
		perpetualId,
		deltaQuantums,
	)
	if err != nil {
		return nil, err
	}

	// Determine the delta between the quote quantums received from closing the position and the
	// bankruptcy price in quote quantums.
	insuranceFundDeltaQuoteQuantumsBig := new(big.Int).Sub(
		deltaQuoteQuantums,
		bankruptcyPriceInQuoteQuantumsBig,
	)

	// If the insurance fund delta is less than or equal to zero, this means the insurance fund
	// needs to cover the difference between the quote quantums received from closing the position
	// and the bankruptcy price in quote quantums.
	if insuranceFundDeltaQuoteQuantumsBig.Sign() <= 0 {
		return insuranceFundDeltaQuoteQuantumsBig, nil
	}

	// The insurance fund delta is positive. We must read the liquidations config from state to
	// determine the max liquidation fee this user must pay.
	liquidationsConfig := k.GetLiquidationsConfig(ctx)

	// Calculate the max liquidation fee from the magnitude of quote quantums the subaccount
	// will receive from closing this position and the max liquidation fee PPM.
	maxLiquidationFeeQuoteQuantumsBig := lib.BigIntMulPpm(
		new(big.Int).Abs(deltaQuoteQuantums),
		liquidationsConfig.MaxLiquidationFeePpm,
	)

	// The liquidation fee paid by the user is the minimum of the max liquidation fee and the
	// leftover collateral from liquidating the position.
	return lib.BigMin(
		maxLiquidationFeeQuoteQuantumsBig,
		insuranceFundDeltaQuoteQuantumsBig,
	), nil
}

// GetPerpetualPositionToLiquidate determines which position to liquidate on the
// passed-in subaccount (after accounting for the `update`). It will return the perpetual id that
// will be used for liquidating the perpetual position.
// This function returns an error if the subaccount has no perpetual positions to liquidate.
func (k Keeper) GetPerpetualPositionToLiquidate(
	ctx sdk.Context,
	subaccountId satypes.SubaccountId,
) (
	perpetualId uint32,
	err error,
) {
	// Fetch the subaccount from state.
	subaccount := k.subaccountsKeeper.GetSubaccount(ctx, subaccountId)

	numPositions := len(subaccount.PerpetualPositions)
	if numPositions > 0 {
		subaccountLiquidationInfo := k.GetSubaccountLiquidationInfo(ctx, subaccountId)
		indexOffset := k.GetPseudoRand(ctx).Intn(numPositions)
		for i := 0; i < numPositions; i++ {
			position := subaccount.PerpetualPositions[(i+indexOffset)%numPositions]
			// Note that this could run in O(n^2) time. This is fine for now because we have less than a hundred
			// perpetuals and only liquidate once per subaccount per block. This means that the position with smallest
			// id will be liquidated first.
			if !subaccountLiquidationInfo.HasPerpetualBeenLiquidatedForSubaccount(position.PerpetualId) {
				return position.PerpetualId, nil
			}
		}
	}

	// Return an error if there are no perpetual positions to liquidate.
	return 0,
		errorsmod.Wrapf(
			types.ErrNoPerpetualPositionsToLiquidate,
			"Subaccount ID: %v",
			subaccount.Id,
		)
}

// GetLiquidatablePositionSizeDelta returns the max number of base quantums to liquidate
// from the perpetual position without exceeding the block and position limits.
func (k Keeper) GetLiquidatablePositionSizeDelta(
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

	clobPair := k.mustGetClobPairForPerpetualId(ctx, perpetualId)

	// Get the maximum notional liquidatable for this position.
	_, bigMaxPositionNotionalLiquidatable, err := k.GetMaxAndMinPositionNotionalLiquidatable(
		ctx,
		perpetualPosition,
	)
	if err != nil {
		return nil, err
	}

	// Get the maximum notional liquidatable for this subaccount.
	bigMaxSubaccountNotionalLiquidatable, err := k.GetSubaccountMaxNotionalLiquidatable(
		ctx,
		subaccountId,
		perpetualId,
	)
	if err != nil {
		return nil, err
	}

	// Take the minimum of the subaccount block limit and position block limit.
	bigMaxQuoteQuantumsLiquidatable := lib.BigMin(
		bigMaxPositionNotionalLiquidatable,
		bigMaxSubaccountNotionalLiquidatable,
	)

	bigQuoteQuantums, err := k.perpetualsKeeper.GetNetNotional(
		ctx,
		perpetualId,
		perpetualPosition.GetBigQuantums(),
	)
	if err != nil {
		panic(err)
	}

	// Return the full position to avoid any rounding errors.
	if bigQuoteQuantums.CmpAbs(bigMaxQuoteQuantumsLiquidatable) <= 0 ||
		perpetualPosition.GetBigQuantums().CmpAbs(
			new(big.Int).SetUint64(clobPair.StepBaseQuantums),
		) <= 0 {
		return new(big.Int).Neg(perpetualPosition.GetBigQuantums()), nil
	}

	// Convert the max notional liquidatable to base quantums.
	absDeltaQuantums, err := k.perpetualsKeeper.GetNotionalInBaseQuantums(
		ctx,
		perpetualId,
		bigMaxQuoteQuantumsLiquidatable,
	)
	if err != nil {
		panic(err)
	}

	// Round to the nearest step size.
	absDeltaQuantums = lib.BigIntRoundToMultiple(
		absDeltaQuantums,
		new(big.Int).SetUint64(clobPair.StepBaseQuantums),
		false,
	)

	// Clamp the base quantums to liquidate to the step size and the size of the position
	// in case there's rounding errors.
	absDeltaQuantums = lib.BigIntClamp(
		absDeltaQuantums,
		new(big.Int).SetUint64(clobPair.StepBaseQuantums),
		new(big.Int).Abs(perpetualPosition.GetBigQuantums()),
	)

	// Negate the position size if it's a long position to get the size delta.
	if perpetualPosition.GetIsLong() {
		return absDeltaQuantums.Neg(absDeltaQuantums), nil
	}

	return absDeltaQuantums, nil
}

// GetSubaccountMaxNotionalLiquidatable returns the maximum notional that the subaccount can liquidate
// without exceeding the subaccount block limits.
// This function takes into account any previous liquidations in the same block and returns an error if
// called with a previously liquidated perpetual id.
func (k Keeper) GetSubaccountMaxNotionalLiquidatable(
	ctx sdk.Context,
	subaccountId satypes.SubaccountId,
	perpetualId uint32,
) (
	bigMaxNotionalLiquidatable *big.Int,
	err error,
) {
	subaccountLiquidationInfo := k.GetSubaccountLiquidationInfo(ctx, subaccountId)

	// Make sure that this subaccount <> perpetual has not previously been liquidated in the same block.
	if subaccountLiquidationInfo.HasPerpetualBeenLiquidatedForSubaccount(perpetualId) {
		return nil, errorsmod.Wrapf(
			types.ErrSubaccountHasLiquidatedPerpetual,
			"Subaccount %v and perpetual %v have already been liquidated within the last block",
			subaccountId,
			perpetualId,
		)
	}

	liquidationConfig := k.GetLiquidationsConfig(ctx)

	// Calculate the maximum notional amount that the given subaccount can liquidate in this block.
	bigTotalNotionalLiquidated := new(big.Int).SetUint64(subaccountLiquidationInfo.NotionalLiquidated)
	bigNotionalLiquidatedBlockLimit := new(big.Int).SetUint64(
		liquidationConfig.SubaccountBlockLimits.MaxNotionalLiquidated,
	)
	if bigTotalNotionalLiquidated.Cmp(bigNotionalLiquidatedBlockLimit) > 0 {
		panic(
			errorsmod.Wrapf(
				types.ErrLiquidationExceedsSubaccountMaxNotionalLiquidated,
				"Subaccount %+v notional liquidated exceeds block limit. Current notional liquidated: %v, block limit: %v",
				subaccountId,
				bigTotalNotionalLiquidated,
				bigNotionalLiquidatedBlockLimit,
			),
		)
	}

	bigMaxNotionalLiquidatable = new(big.Int).Sub(
		bigNotionalLiquidatedBlockLimit,
		bigTotalNotionalLiquidated,
	)

	return bigMaxNotionalLiquidatable, nil
}

// GetSubaccountMaxInsuranceLost returns the maximum insurance fund payout that can be performed
// in this block without exceeding the subaccount block limits.
// This function takes into account any previous liquidations in the same block and returns an error if
// called with a previously liquidated perpetual id.
func (k Keeper) GetSubaccountMaxInsuranceLost(
	ctx sdk.Context,
	subaccountId satypes.SubaccountId,
	perpetualId uint32,
) (
	bigMaxQuantumsInsuranceLost *big.Int,
	err error,
) {
	subaccountLiquidationInfo := k.GetSubaccountLiquidationInfo(ctx, subaccountId)

	// Make sure that the subaccount has not previously liquidated this perpetual in the same block.
	if subaccountLiquidationInfo.HasPerpetualBeenLiquidatedForSubaccount(perpetualId) {
		return nil, errorsmod.Wrapf(
			types.ErrSubaccountHasLiquidatedPerpetual,
			"Subaccount %v and perpetual %v have already been liquidated within the last block",
			subaccountId,
			perpetualId,
		)
	}

	liquidationConfig := k.GetLiquidationsConfig(ctx)

	// Calculate the maximum insurance fund payout amount for the given subaccount in this block.
	bigCurrentInsuranceFundLost := new(big.Int).SetUint64(subaccountLiquidationInfo.QuantumsInsuranceLost)
	bigInsuranceFundLostBlockLimit := new(big.Int).SetUint64(
		liquidationConfig.SubaccountBlockLimits.MaxQuantumsInsuranceLost,
	)
	if bigCurrentInsuranceFundLost.Cmp(bigInsuranceFundLostBlockLimit) > 0 {
		panic(
			errorsmod.Wrapf(
				types.ErrLiquidationExceedsSubaccountMaxInsuranceLost,
				"Subaccount %+v insurance lost exceeds block limit. Current insurance lost: %v, block limit: %v",
				subaccountId,
				bigCurrentInsuranceFundLost,
				bigInsuranceFundLostBlockLimit,
			),
		)
	}

	bigMaxQuantumsInsuranceLost = new(big.Int).Sub(
		bigInsuranceFundLostBlockLimit,
		bigCurrentInsuranceFundLost,
	)
	return bigMaxQuantumsInsuranceLost, nil
}

// GetMaxAndMinPositionNotionalLiquidatable returns the maximum and minimum notional that can be liquidated
// without exceeding the position block limits.
// The minimum amount to liquidate is specified by the liquidation config and is overridden
// by the maximum size of the position.
// The maximum amount of quantums is calculated using max_position_portion_liquidated_ppm
// of the liquidation config, overridden by the minimum notional liquidatable.
func (k Keeper) GetMaxAndMinPositionNotionalLiquidatable(
	ctx sdk.Context,
	positionToLiquidate *satypes.PerpetualPosition,
) (
	bigMinPosNotionalLiquidatable *big.Int,
	bigMaxPosNotionalLiquidatable *big.Int,
	err error,
) {
	liquidationConfig := k.GetLiquidationsConfig(ctx)

	// Get the position size in quote quantums.
	bigNetNotionalQuoteQuantums, err := k.perpetualsKeeper.GetNetNotional(
		ctx,
		positionToLiquidate.PerpetualId,
		positionToLiquidate.GetBigQuantums(),
	)
	if err != nil {
		return nil, nil, err
	}

	bigAbsNetNotionalQuoteQuantums := new(big.Int).Abs(bigNetNotionalQuoteQuantums)
	// Get the mininum notional of this position that can be liquidated, which cannot exceed the size of the position.
	bigMinPosNotionalLiquidatable = lib.BigMin(
		new(big.Int).SetUint64(liquidationConfig.PositionBlockLimits.MinPositionNotionalLiquidated),
		bigAbsNetNotionalQuoteQuantums,
	)

	// Get the maximum notional of this position that can be liquidated,
	// which cannot be less than `minPositionNotionalLiquidated`.
	bigMaxPosNotionalLiquidatable = lib.BigIntMulPpm(
		bigAbsNetNotionalQuoteQuantums,
		liquidationConfig.PositionBlockLimits.MaxPositionPortionLiquidatedPpm,
	)
	bigMaxPosNotionalLiquidatable = lib.BigMax(
		bigMinPosNotionalLiquidatable,
		bigMaxPosNotionalLiquidatable,
	)
	return bigMinPosNotionalLiquidatable, bigMaxPosNotionalLiquidatable, nil
}

// ConvertFillablePriceToSubticks converts the fillable price of a liquidation order to subticks.
// The returned subticks will be rounded to the nearest tick (such that
// `subticks % clobPair.SubticksPerTick == 0`). This function will round up for sells
// that close longs, and round down for buys that close shorts.
//
// Note the returned `subticks` will be bounded (inclusive) between `clobPair.SubticksPerTick` and
// `math.MaxUint64 - math.MaxUint64 % clobPair.SubticksPerTick` (the maximum `uint64` that is a
// multiple of `clobPair.SubticksPerTick`).
func (k Keeper) ConvertFillablePriceToSubticks(
	ctx sdk.Context,
	fillablePrice *big.Rat,
	isLiquidatingLong bool,
	clobPair types.ClobPair,
) (
	subticks types.Subticks,
) {
	// The fillable price is invalid if it is negative.
	if fillablePrice.Sign() < 0 {
		panic("ConvertFillablePriceToSubticks: FillablePrice should not be negative")
	}

	// Assuming `fillablePrice` is in units of `quote quantums / base quantum`,  then dividing by
	// `quote quantums / (base quantum * subtick)` will give the resulting units of subticks.
	exponent := clobPair.QuantumConversionExponent
	p10, inverse := lib.BigPow10(exponent)
	subticksRat := new(big.Rat)
	if inverse {
		subticksRat.SetFrac(
			new(big.Int).Mul(p10, fillablePrice.Num()),
			fillablePrice.Denom(),
		)
	} else {
		subticksRat.SetFrac(
			fillablePrice.Num(),
			new(big.Int).Mul(p10, fillablePrice.Denom()),
		)
	}

	// If we are liquidating a long position with a sell order, then we round up to the nearest
	// subtick (and vice versa for liquidating shorts).
	roundUp := isLiquidatingLong

	// Round the subticks to the nearest `big.Int` in the correct direction.
	roundedSubticksBig := lib.BigRatRound(subticksRat, roundUp)

	// Ensure `roundedSubticksBig % clobPair.SubticksPerTick == 0`, rounding in the correct
	// direction if necessary.
	roundedAlignedSubticksBig := lib.BigIntRoundToMultiple(
		roundedSubticksBig,
		new(big.Int).SetUint64(uint64(clobPair.SubticksPerTick)),
		roundUp,
	)

	// Bound the result between `clobPair.SubticksPerTick` and
	// `math.MaxUint64 - math.MaxUint64 % clobPair.SubticksPerTick`.
	minSubticks := uint64(clobPair.SubticksPerTick)
	maxSubticks := uint64(math.MaxUint64 - (math.MaxUint64 % uint64(clobPair.SubticksPerTick)))
	boundedSubticks := lib.BigUint64Clamp(
		roundedAlignedSubticksBig,
		minSubticks,
		maxSubticks,
	)

	// Panic if the bounded subticks is zero or is not a multiple of `clobPair.SubticksPerTick`,
	// which would indicate the rounding or clamp logic failed.
	if boundedSubticks == 0 {
		panic("ConvertFillablePriceToSubticks: Bounded subticks is 0.")
	} else if boundedSubticks%uint64(clobPair.SubticksPerTick) != 0 {
		panic("ConvertFillablePriceToSubticks: Bounded subticks is not a multiple of SubticksPerTick.")
	}

	return types.Subticks(boundedSubticks)
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

	// Validate that processing the liquidation fill does not leave insufficient funds
	// in the insurance fund (such that the liquidation couldn't have possibly continued).
	if !k.IsValidInsuranceFundDelta(ctx, insuranceFundDelta, perpetualId) {
		log.DebugLog(ctx, "ProcessMatches: insurance fund has insufficient balance to process the liquidation.")
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

// SortLiquidationOrders deterministically sorts the liquidation orders in place.
// Orders are first ordered by their absolute percentage difference from the oracle price in descending order,
// followed by the their size in quote quantums in descending order, and finally by order hashes.
func (k Keeper) SortLiquidationOrders(
	ctx sdk.Context,
	liquidationOrders []types.LiquidationOrder,
) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), metrics.SortLiquidationOrders)

	sort.Slice(liquidationOrders, func(i, j int) bool {
		x, y := liquidationOrders[i], liquidationOrders[j]

		// First, sort by abs percentage difference from oracle price in descending order.
		xAbsPercentageDiffFromOraclePrice := k.getAbsPercentageDiffFromOraclePrice(ctx, x)
		yAbsPercentageDiffFromOraclePrice := k.getAbsPercentageDiffFromOraclePrice(ctx, y)
		if xAbsPercentageDiffFromOraclePrice.Cmp(yAbsPercentageDiffFromOraclePrice) != 0 {
			return xAbsPercentageDiffFromOraclePrice.Cmp(yAbsPercentageDiffFromOraclePrice) == 1
		}

		// Then sort by order quote quantums in descending order.
		xQuoteQuantums := k.getQuoteQuantumsForLiquidationOrder(ctx, x)
		yQuoteQuantums := k.getQuoteQuantumsForLiquidationOrder(ctx, y)
		if xQuoteQuantums.Cmp(yQuoteQuantums) != 0 {
			return xQuoteQuantums.Cmp(yQuoteQuantums) == 1
		}

		// Sort by order hash by default.
		xHash := x.GetOrderHash()
		yHash := y.GetOrderHash()
		return bytes.Compare(xHash[:], yHash[:]) == -1
	})
}

func (k Keeper) getAbsPercentageDiffFromOraclePrice(
	ctx sdk.Context,
	liquidationOrder types.LiquidationOrder,
) *big.Rat {
	clobPair := k.mustGetClobPair(ctx, liquidationOrder.GetClobPairId())
	oraclePriceRat := k.GetOraclePriceSubticksRat(ctx, clobPair)
	fillablePriceRat := liquidationOrder.GetOrderSubticks().ToBigRat()

	return new(big.Rat).Abs(
		new(big.Rat).Quo(
			new(big.Rat).Sub(fillablePriceRat, oraclePriceRat),
			oraclePriceRat,
		),
	)
}

func (k Keeper) getQuoteQuantumsForLiquidationOrder(
	ctx sdk.Context,
	liquidationOrder types.LiquidationOrder,
) *big.Int {
	quoteQuantums, err := k.perpetualsKeeper.GetNetNotional(
		ctx,
		liquidationOrder.MustGetLiquidatedPerpetualId(),
		liquidationOrder.GetBaseQuantums().ToBigInt(),
	)
	if err != nil {
		panic(err)
	}
	return quoteQuantums
}
