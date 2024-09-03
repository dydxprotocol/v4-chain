package keeper

import (
	"errors"
	"math"
	"math/big"
	"time"

	errorsmod "cosmossdk.io/errors"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/dtypes"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/lib"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/lib/log"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/lib/metrics"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/clob/types"
	perpkeeper "github.com/StreamFinance-Protocol/stream-chain/protocol/x/perpetuals/keeper"
	perptypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/perpetuals/types"
	pricestypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices/types"
	satypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/subaccounts/types"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
	subaccountIds *LiquidationPriorityHeap,
) (
	subaccountsToDeleverage []subaccountToDeleverage,
	err error,
) {
	lib.AssertCheckTxMode(ctx)

	metrics.AddSample(
		metrics.LiquidationsLiquidatableSubaccountIdsCount,
		float32(subaccountIds.Len()),
	)

	// Early return if there are 0 subaccounts to liquidate.
	numSubaccounts := subaccountIds.Len()
	if numSubaccounts == 0 {
		return nil, nil
	}

	defer telemetry.MeasureSince(
		time.Now(),
		types.ModuleName,
		metrics.ClobLiquidateSubaccountsAgainstOrderbook,
		metrics.Latency,
	)

	// Process at-most `MaxLiquidationAttemptsPerBlock` in order of priority.
	perpetuals := k.perpetualsKeeper.GetAllPerpetuals(ctx)
	perpetualsMap := lib.UniqueSliceToMap(perpetuals, func(p perptypes.Perpetual) uint32 {
		return p.Params.Id
	})

	numIsolatedLiquidations := 0
	isolatedPositionsPriorityHeap := NewLiquidationPriorityHeap()

	startGetLiquidationOrders := time.Now()
	for i := 0; i < int(k.Flags.MaxLiquidationAttemptsPerBlock); i++ {

		if subaccountIds.Len() == 0 && isolatedPositionsPriorityHeap.Len() > 0 {
			subaccountIds = isolatedPositionsPriorityHeap
			numIsolatedLiquidations = -1000000
		}

		subaccountId := subaccountIds.PopHighestPriority()

		subaccount := k.subaccountsKeeper.GetSubaccount(ctx, subaccountId.SubaccountId)
		if len(subaccount.PerpetualPositions) == 0 {
			i--
			continue
		} else {
			perpetual := perpetualsMap[subaccount.PerpetualPositions[0].PerpetualId]
			if perpetual.Params.MarketType == perptypes.PerpetualMarketType_PERPETUAL_MARKET_TYPE_ISOLATED {
				if numIsolatedLiquidations < int(k.Flags.MaxIsolatedLiquidationAttemptsPerBlock) {
					numIsolatedLiquidations++
				} else {
					isolatedPositionsPriorityHeap.AddSubaccount(subaccountId.SubaccountId, subaccountId.Priority)
					i--
					continue
				}
			}
		}

		// Generate a new liquidation order with the appropriate order size from the sorted subaccount ids.
		liquidationOrder, err := k.MaybeGetLiquidationOrder(ctx, subaccountId.SubaccountId)
		if err != nil {
			// Subaccount might not always be liquidatable if previous liquidation orders
			// improves the net collateral of this subaccount.
			if errors.Is(err, types.ErrSubaccountNotLiquidatable) {
				i--
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

		// the subaccount has been updated to reflect the order match
		subaccount = k.subaccountsKeeper.GetSubaccount(ctx, liquidationOrder.GetSubaccountId())
		isLiquidatable, priority, err := k.GetSubaccountPriority(ctx, subaccount)
		if err != nil {
			return nil, err
		}

		if isLiquidatable {
			subaccountIds.AddSubaccount(liquidationOrder.GetSubaccountId(), priority)
		}

	}
	telemetry.MeasureSince(
		startGetLiquidationOrders,
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

func (k Keeper) GetSubaccountPriority(
	ctx sdk.Context,
	subaccount satypes.Subaccount,
) (
	isLiquidatable bool,
	priority *big.Float,
	err error,
) {

	_, marketPricesMap, perpetualsMap, liquidityTiersMap := k.FetchInformationForLiquidations(ctx)
	isLiquidatable, _, priority, err = k.CheckSubaccountCollateralization(
		subaccount,
		marketPricesMap,
		perpetualsMap,
		liquidityTiersMap,
	)

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

	orderQuantums, err := k.GetNegativePositionSize(
		ctx,
		subaccountId,
		perpetualId,
	)
	if err != nil {
		return nil, err
	}

	bankruptcyPriceQuoteQuantums, err := k.GetBankruptcyPriceInQuoteQuantums(ctx, subaccountId, perpetualId, orderQuantums)
	if err != nil {
		return nil, err
	}
	bankruptcyPriceQuoteQuantumsRat := new(big.Rat).SetInt(bankruptcyPriceQuoteQuantums)
	bankruptcyPriceRat := new(big.Rat).Quo(
		bankruptcyPriceQuoteQuantumsRat,
		new(big.Rat).Neg(new(big.Rat).SetInt(orderQuantums)),
	)

	fillablePriceRat, err := k.GetFillablePrice(ctx, subaccountId, perpetualId, orderQuantums)
	if err != nil {
		return nil, err
	}

	// Calculate the bankruptcy price.
	isLiquidatingLong := orderQuantums.Sign() == -1

	// take the most aggresive price
	liquidationPriceRat := new(big.Rat)
	if isLiquidatingLong {
		if bankruptcyPriceRat.Cmp(fillablePriceRat) < 0 {
			liquidationPriceRat = bankruptcyPriceRat
		} else {
			liquidationPriceRat = fillablePriceRat
		}
	} else {
		if bankruptcyPriceRat.Cmp(fillablePriceRat) > 0 {
			liquidationPriceRat = bankruptcyPriceRat
		} else {
			liquidationPriceRat = fillablePriceRat
		}
	}

	clobPair := k.mustGetClobPairForPerpetualId(ctx, perpetualId)
	liquidationPriceSubticks := k.ConvertLiquidationPriceToSubticks(
		ctx,
		liquidationPriceRat,
		isLiquidatingLong,
		clobPair,
	)

	// Create the liquidation order.
	absBaseQuantums := orderQuantums.Abs(orderQuantums)
	liquidationOrder = types.NewLiquidationOrder(
		subaccountId,
		clobPair,
		!isLiquidatingLong,
		satypes.BaseQuantums(absBaseQuantums.Uint64()),
		liquidationPriceSubticks,
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
	// SOLAL uses price & seems like would change once we change logic
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

// getHealth returns the ratio of collateral to maintenance margin.
// If the net collateral is negative, it returns 0.
// If the maintenance margin is less than or equal to zero, it returns a large number.
//
// This is a stateless function.
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

	// stat liquidation order for negative TNC
	// TODO(CLOB-906) Prevent duplicated stat emissions for liquidation orders in PrepareCheckState.
	if tncBig.Sign() < 0 {
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
	subaccountLiquidationInfo := k.GetSubaccountLiquidationInfo(ctx, subaccountId)
	marketPrices := k.pricesKeeper.GetAllMarketPrices(ctx)
	marketPricesMap := lib.UniqueSliceToMap(marketPrices, func(m pricestypes.MarketPrice) uint32 {
		return m.Id
	})

	bestPriority := big.NewFloat(-1)
	bestPerpetualId := uint32(0)

	for _, position := range subaccount.PerpetualPositions {
		if !subaccountLiquidationInfo.HasPerpetualBeenLiquidatedForSubaccount(position.PerpetualId) {
			closedSubaccount, err := k.simulateClosePerpetualPosition(ctx, subaccount, position, marketPricesMap[position.PerpetualId])
			if err != nil {
				return 0, err
			}

			_, priority, err := k.GetSubaccountPriority(ctx, closedSubaccount)
			if err != nil {
				return 0, err
			}

			if priority.Cmp(bestPriority) > 0 {
				bestPriority = priority
				bestPerpetualId = position.PerpetualId
			}
		}
	}

	if bestPriority.Sign() > 0 {
		return bestPerpetualId, nil
	}

	// Return an error if there are no perpetual positions to liquidate.
	return 0,
		errorsmod.Wrapf(
			types.ErrNoPerpetualPositionsToLiquidate,
			"Subaccount ID: %v",
			subaccount.Id,
		)
}

func (k Keeper) simulateClosePerpetualPosition(
	ctx sdk.Context,
	subaccount satypes.Subaccount,
	position *satypes.PerpetualPosition,
	price pricestypes.MarketPrice,
) (
	closedSubaccount satypes.Subaccount,
	err error,
) {
	// Copy the subaccount to avoid modifying the original
	closedSubaccount = subaccount

	// Find and remove the perpetual position
	for i, pos := range closedSubaccount.PerpetualPositions {
		if pos.PerpetualId == position.PerpetualId {
			// Remove the position from the slice
			closedSubaccount.PerpetualPositions = append(
				closedSubaccount.PerpetualPositions[:i],
				closedSubaccount.PerpetualPositions[i+1:]...,
			)
			break
		}
	}

	// Get the perpetual details
	perpetual, err := k.perpetualsKeeper.GetPerpetual(ctx, position.PerpetualId)
	if err != nil {
		return satypes.Subaccount{}, err
	}

	// Calculate the net notional in quote quantums
	bigQuantums := position.GetBigQuantums()
	bigNetCollateralQuoteQuantums := perpkeeper.GetNetNotionalInQuoteQuantums(perpetual, price, bigQuantums)

	// Add the net notional to the USDC balance
	usdcPosition := closedSubaccount.AssetPositions[0]
	usdcPosition.Quantums = dtypes.NewIntFromBigInt(new(big.Int).Add(usdcPosition.Quantums.BigInt(), bigNetCollateralQuoteQuantums))

	// Update the USDC position in the subaccount
	closedSubaccount.AssetPositions = []*satypes.AssetPosition{usdcPosition}

	return closedSubaccount, nil
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

// ConvertLiquidationPriceToSubticks converts the liquidation price of a liquidation order to subticks.
// The returned subticks will be rounded to the nearest tick (such that
// `subticks % clobPair.SubticksPerTick == 0`). This function will round up for sells
// that close longs, and round down for buys that close shorts.
//
// Note the returned `subticks` will be bounded (inclusive) between `clobPair.SubticksPerTick` and
// `math.MaxUint64 - math.MaxUint64 % clobPair.SubticksPerTick` (the maximum `uint64` that is a
// multiple of `clobPair.SubticksPerTick`).
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

	exponent := clobPair.QuantumConversionExponent
	absExponentiatedValueBig := lib.BigPow10(uint64(lib.AbsInt32(exponent)))
	quoteQuantumsPerBaseQuantumAndSubtickRat := new(big.Rat).SetInt(absExponentiatedValueBig)
	// If `exponent` is negative, invert the fraction to set the result to `1 / 10^exponent`.
	if exponent < 0 {
		quoteQuantumsPerBaseQuantumAndSubtickRat.Inv(quoteQuantumsPerBaseQuantumAndSubtickRat)
	}

	// Assuming `liquidationPrice` is in units of `quote quantums / base quantum`,  then dividing by
	// `quote quantums / (base quantum * subtick)` will give the resulting units of subticks.
	subticksRat := new(big.Rat).Quo(
		liquidationPrice,
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
		panic("ConvertLiquidationPriceToSubticks: Bounded subticks is 0.")
	} else if boundedSubticks%uint64(clobPair.SubticksPerTick) != 0 {
		panic("ConvertLiquidationPriceToSubticks: Bounded subticks is not a multiple of SubticksPerTick.")
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
	if err := k.validateLiquidationParams(
		ctx,
		liquidatedSubaccountId,
		perpetualId,
		insuranceFundDelta,
	); err != nil {
		return nil, err
	}

	return insuranceFundDelta, nil
}

// validateLiquidationParams performs stateful validation
// against the subaccount block limits specified in liquidation configs.
// If validation fails, an error is returned.
//
// The following validation occurs in this method:
//   - The subaccount and perpetual ID pair has not been previously liquidated in the same block.
//   - The total insurance lost does not exceed the maximum insurance lost per block.
func (k Keeper) validateLiquidationParams(
	ctx sdk.Context,
	subaccountId satypes.SubaccountId,
	perpetualId uint32,
	insuranceFundDeltaQuoteQuantums *big.Int,
) (
	err error,
) {
	// Validate that this liquidation does not exceed the maximum notional amount that a single subaccount can have
	// liquidated per block.

	subaccountLiquidationInfo := k.GetSubaccountLiquidationInfo(ctx, subaccountId)

	// Make sure that this subaccount <> perpetual has not previously been liquidated in the same block.
	if subaccountLiquidationInfo.HasPerpetualBeenLiquidatedForSubaccount(perpetualId) {
		return errorsmod.Wrapf(
			types.ErrSubaccountHasLiquidatedPerpetual,
			"Subaccount %v and perpetual %v have already been liquidated within the last block",
			subaccountId,
			perpetualId,
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
