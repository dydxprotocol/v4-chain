package keeper

import (
	errorsmod "cosmossdk.io/errors"
	"math"
	"math/big"

	gometrics "github.com/armon/go-metrics"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

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
	isLiquidatable, err := k.IsLiquidatable(ctx, subaccountId)
	if err != nil {
		return nil, err
	}
	if !isLiquidatable {
		return nil, types.ErrSubaccountNotLiquidatable
	}

	// The subaccount is liquidatable. Get the perpetual position and position size to liquidate.
	clobPair, positionSizeBig, err := k.GetPerpetualPositionToLiquidate(ctx, subaccountId)
	if err != nil {
		return nil, err
	}
	perpetualId := clobPair.GetPerpetualClobMetadata().PerpetualId

	// Get the fillable price of the liquidation order in subticks.
	deltaQuantumsBig := positionSizeBig.Neg(positionSizeBig)
	fillablePriceRat, err := k.GetFillablePrice(ctx, subaccountId, perpetualId, deltaQuantumsBig)
	if err != nil {
		return nil, err
	}
	isLiquidatingLong := deltaQuantumsBig.Sign() == -1
	fillablePriceSubticks := k.ConvertFillablePriceToSubticks(
		ctx,
		fillablePriceRat,
		isLiquidatingLong,
		clobPair,
	)

	// Create the liquidation order.
	positionSize := deltaQuantumsBig.Abs(deltaQuantumsBig).Uint64()
	liquidationOrder = types.NewLiquidationOrder(
		subaccountId,
		clobPair,
		!isLiquidatingLong,
		satypes.BaseQuantums(positionSize),
		fillablePriceSubticks,
	)
	return liquidationOrder, nil
}

// PlacePerpetualLiquidation places an IOC liquidation order onto the book that results in fills of type
// `PerpetualLiquidation`.
func (k Keeper) PlacePerpetualLiquidation(
	ctx sdk.Context,
	liquidationOrder types.LiquidationOrder,
) (
	orderSizeOptimisticallyFilledFromMatchingQuantums satypes.BaseQuantums,
	orderStatus types.OrderStatus,
	err error,
) {
	lib.AssertCheckTxMode(ctx)

	orderSizeOptimisticallyFilledFromMatchingQuantums,
		orderStatus,
		offchainUpdates,
		err :=
		k.MemClob.PlacePerpetualLiquidation(
			ctx,
			liquidationOrder,
		)

	// TODO(DEC-1323): Potentially allow liquidating the same perpetual + subaccount
	// multiple times in a block.
	k.MustUpdateSubaccountPerpetualLiquidated(
		ctx,
		liquidationOrder.GetSubaccountId(),
		liquidationOrder.MustGetLiquidatedPerpetualId(),
	)

	telemetry.IncrCounter(
		1,
		metrics.Liquidations,
		metrics.PlacePerpetualLiquidation,
		metrics.Count,
	)

	telemetry.IncrCounterWithLabels(
		[]string{metrics.Liquidations, metrics.PlacePerpetualLiquidation, metrics.BaseQuantums},
		metrics.GetMetricValueFromBigInt(liquidationOrder.GetBaseQuantums().ToBigInt()),
		[]gometrics.Label{
			metrics.GetLabelForIntValue(metrics.PerpetualId, int(liquidationOrder.MustGetLiquidatedPerpetualId())),
		},
	)

	telemetry.IncrCounterWithLabels(
		[]string{metrics.Liquidations, metrics.PlacePerpetualLiquidation, metrics.Filled, metrics.BaseQuantums},
		metrics.GetMetricValueFromBigInt(orderSizeOptimisticallyFilledFromMatchingQuantums.ToBigInt()),
		[]gometrics.Label{
			metrics.GetLabelForIntValue(metrics.PerpetualId, int(liquidationOrder.MustGetLiquidatedPerpetualId())),
		},
	)

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

	// The subaccount is liquidatable if both of the following are true:
	// - The maintenance margin requirements are greater than zero (note that they can never be negative).
	// - The maintenance margin requirements are greater than the subaccount's net collateral.
	isLiquidatable := bigMaintenanceMargin.Sign() > 0 && bigMaintenanceMargin.Cmp(bigNetCollateral) == 1
	return isLiquidatable, nil
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

	tncBig, _, tmmrBig, err := k.subaccountsKeeper.GetNetCollateralAndMarginRequirements(
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

	// `DNNV = PNNVAD - PNNV`, where `PNNVAD` is the perpetual's net notional
	// with a position size of `PS + deltaQuantums`.
	// Note that we are intentionally not calculating `DNNV` from `deltaQuantums`
	// directly to avoid rounding errors.
	pnnvBig, err := k.perpetualsKeeper.GetNetNotional(
		ctx,
		perpetualId,
		psBig,
	)
	if err != nil {
		return nil, err
	}

	pnnvadBig, err := k.perpetualsKeeper.GetNetNotional(
		ctx,
		perpetualId,
		new(big.Int).Add(psBig, deltaQuantums),
	)
	if err != nil {
		return nil, err
	}

	dnnvBig := new(big.Int).Sub(pnnvadBig, pnnvBig)

	// `DMMR = PMMRAD - PMMR`, where `PMMRAD` is the perpetual's maintenance margin requirement
	// with a position size of `PS + deltaQuantums`.
	// Note that we cannot directly calculate `DMMR` from `deltaQuantums` because the maintenance
	// margin requirement function could be non-linear.
	_, pmmrBig, err := k.perpetualsKeeper.GetMarginRequirements(ctx, perpetualId, psBig)
	if err != nil {
		return nil, err
	}

	_, pmmradBig, err := k.perpetualsKeeper.GetMarginRequirements(
		ctx,
		perpetualId,
		new(big.Int).Add(
			psBig,
			deltaQuantums,
		),
	)
	if err != nil {
		return nil, err
	}

	dmmrBig := new(big.Int).Sub(pmmradBig, pmmrBig)
	// `dmmrBig` should never be positive if `| PS | >= | PS + deltaQuantums |`. If it is, panic.
	if dmmrBig.Sign() == 1 {
		panic("GetBankruptcyPriceInQuoteQuantums: DMMR is positive")
	}

	// Calculate `TNC * abs(DMMR) / TMMR`.
	tncMulDmmrBig := new(big.Int).Mul(tncBig, new(big.Int).Abs(dmmrBig))
	// This calculation is intentionally rounded down to negative infinity to ensure the
	// final result is rounded towards positive-infinity. This works because of the following:
	// - This is the only division in the equation.
	// - This calculation is subtracted from `-DNNV` to get the final result.
	// - The dividend `TNC * abs(DMMR)` is the only number that can be negative, and `Div` uses
	//   Euclidean division so even if `TNC < 0` this will still round towards negative infinity.
	quoteQuantumsBeforeBankruptcyBig := new(big.Int).Div(tncMulDmmrBig, tmmrBig)

	// Calculate `-DNNV - TNC * abs(DMMR) / TMMR`.
	bankruptcyPriceQuoteQuantumsBig := new(big.Int).Sub(
		new(big.Int).Neg(dnnvBig),
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

	pnnvBig, err := k.perpetualsKeeper.GetNetCollateral(ctx, perpetualId, psBig)
	if err != nil {
		return nil, err
	}

	_, pmmrBig, err := k.perpetualsKeeper.GetMarginRequirements(ctx, perpetualId, psBig)
	if err != nil {
		return nil, err
	}

	tncBig,
		_,
		tmmrBig,
		err := k.subaccountsKeeper.GetNetCollateralAndMarginRequirements(
		ctx,
		satypes.Update{SubaccountId: subaccountId},
	)
	if err != nil {
		return nil, err
	}

	liquidationsConfig := k.GetLiquidationsConfig(ctx)
	ba := liquidationsConfig.FillablePriceConfig.BankruptcyAdjustmentPpm
	smmr := liquidationsConfig.FillablePriceConfig.SpreadToMaintenanceMarginRatioPpm

	// Calculate the ABR (adjusted bankruptcy rating).
	tncDivTmmrRat := new(big.Rat).SetFrac(tncBig, tmmrBig)
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
		new(big.Rat).SetInt(pmmrBig),
		smmr,
	)

	fillablePriceOracleDeltaQuoteQuantumsRat := new(big.Rat).Mul(abrRat, maxLiquidationSpreadQuoteQuantumsRat)

	// Calculate `PNNV - ABR * SMMR * PMMR`, which represents the fillable price in quote quantums.
	// For longs, `pnnvRat > 0` meaning the fillable price in quote quantums will be lower than the
	// oracle price.
	// For shorts, `pnnvRat < 0` meaning the fillable price in quote quantums will be higher than
	// the oracle price (in this case the result will be negative, but dividing by `positionSize` below
	// will make it positive since `positionSize < 0` for shorts).
	pnnvRat := new(big.Rat).SetInt(pnnvBig)
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
	deltaQuoteQuantums, err := getFillQuoteQuantums(
		clobPair,
		subticks,
		satypes.BaseQuantums(fillAmount),
	)
	if err != nil {
		return nil, err
	}
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

// GetPerpetualPositionToLiquidate determines which position and position size to liquidate on the
// passed-in subaccount (after accounting for the `update`). It will return the `ClobPair` that
// will be used for liquidating the perpetual position and the number of quantums to liquidate
// from the perpetual position (positive if long, negative if short).
// This function returns an error if the subaccount has no perpetual positions to liquidate.
func (k Keeper) GetPerpetualPositionToLiquidate(
	ctx sdk.Context,
	subaccountId satypes.SubaccountId,
) (
	clobPair types.ClobPair,
	quantums *big.Int,
	err error,
) {
	// Fetch the subaccount from state.
	subaccount := k.subaccountsKeeper.GetSubaccount(ctx, subaccountId)
	subaccountLiquidationInfo := k.GetSubaccountLiquidationInfo(ctx, subaccountId)

	var perpetualPosition *satypes.PerpetualPosition

	for _, position := range subaccount.PerpetualPositions {
		// Note that this could run in O(n^2) time. This is fine for now because we have less than a hundred
		// perpetuals and only liquidate once per subaccount per block. This means that the position with smallest
		// id will be liquidated first.
		if !subaccountLiquidationInfo.HasPerpetualBeenLiquidatedForSubaccount(position.PerpetualId) {
			perpetualPosition = position
			break
		}
	}

	// Return an error if there are no perpetual positions to liquidate.
	if perpetualPosition == nil {
		return types.ClobPair{},
			nil,
			errorsmod.Wrapf(
				types.ErrNoPerpetualPositionsToLiquidate,
				"Subaccount ID: %v",
				subaccount.Id,
			)
	}

	clobPair = k.mustGetClobPairForPerpetualId(ctx, perpetualPosition.PerpetualId)

	// Get the maximum notional liquidatable for this position.
	_, bigMaxPositionNotionalLiquidatable, err := k.GetMaxAndMinPositionNotionalLiquidatable(
		ctx,
		perpetualPosition,
	)
	if err != nil {
		panic(err)
	}

	// Get the maximum notional liquidatable for this subaccount.
	bigMaxSubaccountNotionalLiquidatable, err := k.GetSubaccountMaxNotionalLiquidatable(
		ctx,
		subaccountId,
		perpetualPosition.PerpetualId,
	)
	if err != nil {
		panic(err)
	}

	// Take the minimum of the subaccount block limit and position block limit.
	bigMaxQuoteQuantumsLiquidatable := lib.BigMin(
		bigMaxPositionNotionalLiquidatable,
		bigMaxSubaccountNotionalLiquidatable,
	)

	bigQuoteQuantums, err := k.perpetualsKeeper.GetNetNotional(
		ctx,
		perpetualPosition.PerpetualId,
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
		return clobPair, perpetualPosition.GetBigQuantums(), nil
	}

	// Convert the max notional liquidatable to base quantums.
	bigBaseQuantumsToLiquidate, err := k.perpetualsKeeper.GetNotionalInBaseQuantums(
		ctx,
		perpetualPosition.PerpetualId,
		bigMaxQuoteQuantumsLiquidatable,
	)
	if err != nil {
		panic(err)
	}

	// Round to the nearest step size.
	bigBaseQuantumsToLiquidate = lib.BigIntRoundToMultiple(
		bigBaseQuantumsToLiquidate,
		new(big.Int).SetUint64(clobPair.StepBaseQuantums),
		false,
	)

	// Clamp the base quantums to liquidate to the step size and the size of the position
	// in case there's rounding errors.
	bigBaseQuantumsToLiquidate = lib.BigIntClamp(
		bigBaseQuantumsToLiquidate,
		new(big.Int).SetUint64(clobPair.StepBaseQuantums),
		new(big.Int).Abs(perpetualPosition.GetBigQuantums()),
	)

	// Negate the position size if it's short.
	if !perpetualPosition.GetIsLong() {
		bigBaseQuantumsToLiquidate.Neg(bigBaseQuantumsToLiquidate)
	}

	return clobPair, bigBaseQuantumsToLiquidate, nil
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

	exponent := clobPair.QuantumConversionExponent
	absExponentiatedValueBig := lib.BigPow10(uint64(lib.AbsInt32(exponent)))
	quoteQuantumsPerBaseQuantumAndSubtickRat := new(big.Rat).SetInt(absExponentiatedValueBig)
	// If `exponent` is negative, invert the fraction to set the result to `1 / 10^exponent`.
	if exponent < 0 {
		quoteQuantumsPerBaseQuantumAndSubtickRat.Inv(quoteQuantumsPerBaseQuantumAndSubtickRat)
	}

	// Assuming `fillablePrice` is in units of `quote quantums / base quantum`,  then dividing by
	// `quote quantums / (base quantum * subtick)` will give the resulting units of subticks.
	subticksRat := new(big.Rat).Quo(
		fillablePrice,
		quoteQuantumsPerBaseQuantumAndSubtickRat,
	)

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
