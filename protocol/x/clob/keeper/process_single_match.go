package keeper

import (
	"errors"
	"fmt"
	"math"
	"math/big"

	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/lib/log"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	affiliatetypes "github.com/dydxprotocol/v4-chain/protocol/x/affiliates/types"
	assettypes "github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	revsharetypes "github.com/dydxprotocol/v4-chain/protocol/x/revshare/types"
	statstypes "github.com/dydxprotocol/v4-chain/protocol/x/stats/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	gometrics "github.com/hashicorp/go-metrics"
)

// ProcessSingleMatch accepts a single match and its associated orders matched in the block,
// persists the resulting subaccount updates and state fill amounts.
// This function assumes that the provided match with orders has undergone stateless validations.
// If additional validation of the provided orders or match fails, an error is returned.
// The following validation occurs in this method:
//   - Order is for a valid ClobPair.
//   - Order is for a valid Perpetual.
//   - Validate the `fillAmount` of a match is divisible by the `ClobPair`'s `StepBaseQuantums`.
//   - Validate the new total fill amount of an order does not exceed the total quantums of the order given
//     the fill amounts present in the provided `matchOrders` and in state.
//   - Validate the subaccount updates resulting from the match are valid (before persisting the updates to state)
//   - For liquidation orders, stateful validations through
//     calling `validateMatchPerpetualLiquidationAgainstSubaccountBlockLimits`.
//   - Validating that deleveraging is not required for processing liquidation orders.
//
// This method returns `takerUpdateResult` and `makerUpdateResult` which can be used to determine whether the maker
// and/or taker failed collateralization checks. This information is particularly pertinent for the `memclob` which
// calls this method during matching.
// TODO(DEC-1282): Remove redundant checks from `ProcessSingleMatch` for matching.
// This method mutates matchWithOrders by setting the fee fields.
func (k Keeper) ProcessSingleMatch(
	ctx sdk.Context,
	matchWithOrders *types.MatchWithOrders,
	affiliateOverrides map[string]bool,
	affiliateParameters affiliatetypes.AffiliateParameters,
) (
	success bool,
	takerUpdateResult satypes.UpdateResult,
	makerUpdateResult satypes.UpdateResult,
	affiliateRevSharesQuoteQuantums *big.Int,
	err error,
) {
	if matchWithOrders.TakerOrder.IsLiquidation() {
		defer func() {
			if errors.Is(err, satypes.ErrFailedToUpdateSubaccounts) && !takerUpdateResult.IsSuccess() {
				takerSubaccount := k.subaccountsKeeper.GetSubaccount(ctx, matchWithOrders.TakerOrder.GetSubaccountId())
				riskTaker, _ := k.subaccountsKeeper.GetNetCollateralAndMarginRequirements(
					ctx,
					satypes.Update{SubaccountId: *takerSubaccount.Id},
				)
				log.ErrorLog(ctx,
					"collateralization check failed for liquidation",
					"takerSubaccount", fmt.Sprintf("%+v", takerSubaccount),
					"takerTNC", riskTaker.NC,
					"takerIMR", riskTaker.IMR,
					"takerMMR", riskTaker.MMR,
					"liquidationOrder", fmt.Sprintf("%+v", matchWithOrders.TakerOrder),
					"makerOrder", fmt.Sprintf("%+v", matchWithOrders.MakerOrder),
					"fillAmount", matchWithOrders.FillAmount,
					"result", takerUpdateResult,
				)
			}
		}()
	}

	// Perform stateless validation on the match.
	if err := matchWithOrders.Validate(); err != nil {
		return false, takerUpdateResult, makerUpdateResult, affiliateRevSharesQuoteQuantums, errorsmod.Wrapf(
			err,
			"ProcessSingleMatch: Invalid MatchWithOrders: %+v",
			matchWithOrders,
		)
	}

	makerMatchableOrder := matchWithOrders.MakerOrder
	takerMatchableOrder := matchWithOrders.TakerOrder
	fillAmount := matchWithOrders.FillAmount

	// Retrieve the ClobPair from state.
	clobPairId := makerMatchableOrder.GetClobPairId()
	clobPair, found := k.GetClobPair(ctx, clobPairId)
	if !found {
		return false, takerUpdateResult, makerUpdateResult, affiliateRevSharesQuoteQuantums, types.ErrInvalidClob
	}

	// Verify that the `fillAmount` is divisible by the `StepBaseQuantums` of the `clobPair`.
	if fillAmount.ToUint64()%clobPair.StepBaseQuantums != 0 {
		return false,
			takerUpdateResult,
			makerUpdateResult,
			affiliateRevSharesQuoteQuantums,
			types.ErrFillAmountNotDivisibleByStepSize
	}

	// Define local variable relevant to retrieving QuoteQuantums based on the fill amount.
	makerSubticks := makerMatchableOrder.GetOrderSubticks()

	// Calculate the number of quote quantums for the match based on the maker order subticks.
	bigFillQuoteQuantums := types.FillAmountToQuoteQuantums(
		makerSubticks,
		fillAmount,
		clobPair.QuantumConversionExponent,
	)

	if bigFillQuoteQuantums.Sign() == 0 {
		// Note: If `subticks`, `baseQuantums`, are small enough, `quantumConversionExponent` is negative,
		// it's possible to have zero `quoteQuantums` for a non-zero amount of `baseQuantums`.
		// This could mean that it's possible that a maker sell order on the book
		// at a very unfavorable price (subticks) could receive `0` `quoteQuantums` amount.
		log.ErrorLog(
			ctx,
			"Match resulted in zero quote quantums",
			"MakerOrder",
			fmt.Sprintf("%+v", matchWithOrders.MakerOrder),
			"TakerOrder",
			fmt.Sprintf("%+v", matchWithOrders.TakerOrder),
			"FillAmount",
			matchWithOrders.FillAmount.ToUint64(),
		)
	}

	// Retrieve the associated perpetual id for the `ClobPair`.
	// TODO(OTE-805): call this outside of ProcessSingleMatch to avoid duplicate calls.
	perpetualId, err := clobPair.GetPerpetualId()
	if err != nil {
		return false, takerUpdateResult, makerUpdateResult, affiliateRevSharesQuoteQuantums, err
	}

	// Fee tier for affiliates
	referreeIndexOverride := affiliateParameters.RefereeMinimumFeeTierIdx

	// Calculate taker and maker fee ppms.
	takerFeePpm := k.feeTiersKeeper.GetPerpetualFeePpm(
		ctx, matchWithOrders.TakerOrder.GetSubaccountId().Owner, true, referreeIndexOverride, clobPairId.ToUint32())
	makerFeePpm := k.feeTiersKeeper.GetPerpetualFeePpm(
		ctx, matchWithOrders.MakerOrder.GetSubaccountId().Owner, false, referreeIndexOverride, clobPairId.ToUint32())

	takerInsuranceFundDelta := new(big.Int)
	if takerMatchableOrder.IsLiquidation() {
		// Liquidation orders do not pay trading fees because they already pay a liquidation fee.
		takerFeePpm = 0
		// Temporarily cap maker rebates to 0 for liquidations. This is to prevent an issue where
		// the fee collector has insufficient funds to pay the maker rebate.
		// TODO(CLOB-812): find a longer term solution to handle maker rebates for liquidations.
		makerFeePpm = lib.Max(makerFeePpm, 0)
		takerInsuranceFundDelta, err = k.validateMatchedLiquidation(
			ctx,
			takerMatchableOrder,
			perpetualId,
			fillAmount,
			makerMatchableOrder.GetOrderSubticks(),
		)

		if err != nil {
			return false, takerUpdateResult, makerUpdateResult, affiliateRevSharesQuoteQuantums, err
		}
	}

	// Calculate the new fill amounts and pruneable block heights for the orders.
	var curTakerFillAmount satypes.BaseQuantums
	var curTakerPruneableBlockHeight uint32
	var newTakerTotalFillAmount satypes.BaseQuantums
	var curMakerFillAmount satypes.BaseQuantums
	var curMakerPruneableBlockHeight uint32
	var newMakerTotalFillAmount satypes.BaseQuantums

	// Liquidation orders can only be placed when a subaccount is liquidatable
	// and cannot be replayed, therefore we don't need to track their filled amount in state.
	if !takerMatchableOrder.IsLiquidation() {
		// Retrieve the current fillAmount and current pruneableBlockHeight for the taker order.
		// If the order has never been filled before, these will both be `0`.
		_,
			curTakerFillAmount,
			curTakerPruneableBlockHeight = k.GetOrderFillAmount(
			ctx,
			matchWithOrders.TakerOrder.MustGetOrder().OrderId,
		)

		// Verify the orders have sufficient remaining quantums, and calculate the new total fill amount.
		newTakerTotalFillAmount, err = getUpdatedOrderFillAmount(
			matchWithOrders.TakerOrder.MustGetOrder().OrderId,
			matchWithOrders.TakerOrder.GetBaseQuantums(),
			curTakerFillAmount,
			fillAmount,
		)

		if err != nil {
			return false, takerUpdateResult, makerUpdateResult, affiliateRevSharesQuoteQuantums, err
		}
	}

	// Retrieve the current fillAmount and current pruneableBlockHeight for the maker order.
	// If the order has never been filled before, these will both be `0`.
	_,
		curMakerFillAmount,
		curMakerPruneableBlockHeight = k.GetOrderFillAmount(
		ctx,
		matchWithOrders.MakerOrder.MustGetOrder().OrderId,
	)

	// Verify the orders have sufficient remaining quantums, and calculate the new total fill amount.
	newMakerTotalFillAmount, err = getUpdatedOrderFillAmount(
		matchWithOrders.MakerOrder.MustGetOrder().OrderId,
		matchWithOrders.MakerOrder.GetBaseQuantums(),
		curMakerFillAmount,
		fillAmount,
	)

	if err != nil {
		return false, takerUpdateResult, makerUpdateResult, affiliateRevSharesQuoteQuantums, err
	}

	// Update both subaccounts in the matched order atomically.
	takerUpdateResult, makerUpdateResult, affiliateRevSharesQuoteQuantums, err = k.persistMatchedOrders(
		ctx,
		matchWithOrders,
		perpetualId,
		takerFeePpm,
		makerFeePpm,
		bigFillQuoteQuantums,
		takerInsuranceFundDelta,
		affiliateOverrides,
		affiliateParameters,
	)

	if err != nil {
		return false, takerUpdateResult, makerUpdateResult, affiliateRevSharesQuoteQuantums, err
	}

	// Update subaccount total quantums liquidated and total insurance fund lost for liquidation orders.
	if matchWithOrders.TakerOrder.IsLiquidation() {
		notionalLiquidatedQuoteQuantums, err := k.perpetualsKeeper.GetNetNotional(
			ctx,
			perpetualId,
			fillAmount.ToBigInt(),
		)
		if err != nil {
			return false, takerUpdateResult, makerUpdateResult, affiliateRevSharesQuoteQuantums, err
		}

		k.UpdateSubaccountLiquidationInfo(
			ctx,
			matchWithOrders.TakerOrder.GetSubaccountId(),
			notionalLiquidatedQuoteQuantums,
			takerInsuranceFundDelta,
		)

		labels := []gometrics.Label{
			metrics.GetLabelForIntValue(metrics.PerpetualId, int(perpetualId)),
			metrics.GetLabelForBoolValue(metrics.CheckTx, ctx.IsCheckTx()),
		}
		if matchWithOrders.TakerOrder.IsBuy() {
			labels = append(labels, metrics.GetLabelForStringValue(metrics.OrderSide, metrics.Buy))
		} else {
			labels = append(labels, metrics.GetLabelForStringValue(metrics.OrderSide, metrics.Sell))
		}

		// Stat quote quantums liquidated.
		gometrics.AddSampleWithLabels(
			[]string{metrics.Liquidations, metrics.PlacePerpetualLiquidation, metrics.Filled, metrics.QuoteQuantums},
			metrics.GetMetricValueFromBigInt(notionalLiquidatedQuoteQuantums),
			labels,
		)
		// Stat insurance fund delta.
		gometrics.AddSampleWithLabels(
			[]string{metrics.Liquidations, metrics.InsuranceFundDelta},
			metrics.GetMetricValueFromBigInt(new(big.Int).Abs(takerInsuranceFundDelta)),
			append(labels, metrics.GetLabelForBoolValue(metrics.Positive, takerInsuranceFundDelta.Sign() == 1)),
		)
	}

	// Liquidation orders can only be placed when a subaccount is liquidatable
	// and cannot be replayed, therefore we don't need to track their filled amount in state.
	if !matchWithOrders.TakerOrder.IsLiquidation() {
		k.setOrderFillAmountsAndPruning(
			ctx,
			matchWithOrders.TakerOrder.MustGetOrder(),
			newTakerTotalFillAmount,
			curTakerPruneableBlockHeight,
		)
	}

	k.setOrderFillAmountsAndPruning(
		ctx,
		matchWithOrders.MakerOrder.MustGetOrder(),
		newMakerTotalFillAmount,
		curMakerPruneableBlockHeight,
	)

	// Check and update the remaining TWAP quantity for both maker and taker orders.
	makerOrder := matchWithOrders.MakerOrder.MustGetOrder()
	if err := k.checkAndUpdateTWAPOrderRemainingQuantity(ctx, makerOrder.OrderId, fillAmount); err != nil {
		return false, takerUpdateResult, makerUpdateResult, affiliateRevSharesQuoteQuantums, err
	}

	if !matchWithOrders.TakerOrder.IsLiquidation() {
		takerOrder := matchWithOrders.TakerOrder.MustGetOrder()
		if err := k.checkAndUpdateTWAPOrderRemainingQuantity(ctx, takerOrder.OrderId, fillAmount); err != nil {
			return false, takerUpdateResult, makerUpdateResult, affiliateRevSharesQuoteQuantums, err
		}
	}

	return true, takerUpdateResult, makerUpdateResult, affiliateRevSharesQuoteQuantums, nil
}

// persistMatchedOrders persists a matched order to the subaccount state,
// by updating the quoteBalance and perpetual position size of the
// affected subaccounts.
// This method also transfers fees to the fee collector module, and
// transfers insurance fund payments to the insurance fund.
// This method mutates matchWithOrders by setting the fee fields.
func (k Keeper) persistMatchedOrders(
	ctx sdk.Context,
	matchWithOrders *types.MatchWithOrders,
	perpetualId uint32,
	takerFeePpm int32,
	makerFeePpm int32,
	bigFillQuoteQuantums *big.Int,
	insuranceFundDelta *big.Int,
	affiliateOverrides map[string]bool,
	affiliateParameters affiliatetypes.AffiliateParameters,
) (
	takerUpdateResult satypes.UpdateResult,
	makerUpdateResult satypes.UpdateResult,
	affiliateRevSharesQuoteQuantums *big.Int,
	err error,
) {
	isTakerLiquidation := matchWithOrders.TakerOrder.IsLiquidation()
	affiliateRevSharesQuoteQuantums = big.NewInt(0)

	// Taker fees and maker fees/rebates are rounded towards positive infinity.
	bigTakerFeeQuoteQuantums := lib.BigMulPpm(bigFillQuoteQuantums, lib.BigI(takerFeePpm), true)
	bigMakerFeeQuoteQuantums := lib.BigMulPpm(bigFillQuoteQuantums, lib.BigI(makerFeePpm), true)

	matchWithOrders.MakerFee = bigMakerFeeQuoteQuantums.Int64()
	// Liquidation orders pay the liquidation fee instead of the standard taker fee
	if matchWithOrders.TakerOrder.IsLiquidation() {
		matchWithOrders.TakerFee = insuranceFundDelta.Int64()
	} else {
		matchWithOrders.TakerFee = bigTakerFeeQuoteQuantums.Int64()
	}

	// If the taker is a liquidation order, it should never pay fees.
	if isTakerLiquidation && bigTakerFeeQuoteQuantums.Sign() != 0 {
		panic(fmt.Sprintf(
			`Taker order is liquidation and should never pay taker fees.
      TakerOrder: %v
      bigTakerFeeQuoteQuantums: %v`,
			matchWithOrders.TakerOrder,
			bigTakerFeeQuoteQuantums,
		))
	}

	bigTakerQuoteBalanceDelta := new(big.Int).Set(bigFillQuoteQuantums)
	bigMakerQuoteBalanceDelta := new(big.Int).Set(bigFillQuoteQuantums)

	bigTakerPerpetualQuantumsDelta := matchWithOrders.FillAmount.ToBigInt()
	bigMakerPerpetualQuantumsDelta := matchWithOrders.FillAmount.ToBigInt()

	if matchWithOrders.TakerOrder.IsBuy() {
		bigTakerQuoteBalanceDelta.Neg(bigTakerQuoteBalanceDelta)
		bigMakerPerpetualQuantumsDelta.Neg(bigMakerPerpetualQuantumsDelta)
	} else {
		bigMakerQuoteBalanceDelta.Neg(bigMakerQuoteBalanceDelta)
		bigTakerPerpetualQuantumsDelta.Neg(bigTakerPerpetualQuantumsDelta)
	}

	// Subtract quote balance delta with fees paid.
	bigTakerQuoteBalanceDelta.Sub(bigTakerQuoteBalanceDelta, bigTakerFeeQuoteQuantums)
	bigMakerQuoteBalanceDelta.Sub(bigMakerQuoteBalanceDelta, bigMakerFeeQuoteQuantums)

	// Subtract quote balance delta with insurance fund payments.
	if matchWithOrders.TakerOrder.IsLiquidation() {
		bigTakerQuoteBalanceDelta.Sub(bigTakerQuoteBalanceDelta, insuranceFundDelta)
	}

	// apply broker fees for taker and maker separately

	if matchWithOrders.MakerOrder.IsLiquidation() {
		panic("maker order can not be a liquidation order")
	}

	makerBuilderCodeParams := matchWithOrders.MakerOrder.MustGetOrder().BuilderCodeParameters
	makerBuilderFeeQuantums := makerBuilderCodeParams.GetBuilderFee(bigFillQuoteQuantums)
	matchWithOrders.MakerBuilderFee = makerBuilderFeeQuantums.Uint64()

	bigMakerQuoteBalanceDelta.Sub(bigMakerQuoteBalanceDelta, makerBuilderFeeQuantums)
	makerBuilderAddress := makerBuilderCodeParams.GetBuilderAddress()

	takerBuilderFeeQuantums := big.NewInt(0)
	var takerBuilderAddress string
	if !matchWithOrders.TakerOrder.IsLiquidation() {
		takerBuilderCodeParams := matchWithOrders.TakerOrder.MustGetOrder().BuilderCodeParameters
		takerBuilderFeeQuantums = takerBuilderCodeParams.GetBuilderFee(bigFillQuoteQuantums)
		bigTakerQuoteBalanceDelta.Sub(bigTakerQuoteBalanceDelta, takerBuilderFeeQuantums)

		matchWithOrders.TakerBuilderFee = takerBuilderFeeQuantums.Uint64()
		takerBuilderAddress = takerBuilderCodeParams.GetBuilderAddress()
	}

	// Do this before subaccount updates so that bank sends are always valid between different
	// module accounts.
	if err := k.subaccountsKeeper.TransferInsuranceFundPayments(ctx, insuranceFundDelta, perpetualId); err != nil {
		return takerUpdateResult, makerUpdateResult, affiliateRevSharesQuoteQuantums, err
	}

	// Create the subaccount update.
	updates := []satypes.Update{
		// Taker update
		{
			AssetUpdates: []satypes.AssetUpdate{
				{
					AssetId:          assettypes.AssetUsdc.Id,
					BigQuantumsDelta: bigTakerQuoteBalanceDelta,
				},
			},
			PerpetualUpdates: []satypes.PerpetualUpdate{
				{
					PerpetualId:      perpetualId,
					BigQuantumsDelta: bigTakerPerpetualQuantumsDelta,
				},
			},
			SubaccountId: matchWithOrders.TakerOrder.GetSubaccountId(),
		},
		// Maker update
		{
			AssetUpdates: []satypes.AssetUpdate{
				{
					AssetId:          assettypes.AssetUsdc.Id,
					BigQuantumsDelta: bigMakerQuoteBalanceDelta,
				},
			},
			PerpetualUpdates: []satypes.PerpetualUpdate{
				{
					PerpetualId:      perpetualId,
					BigQuantumsDelta: bigMakerPerpetualQuantumsDelta,
				},
			},
			SubaccountId: matchWithOrders.MakerOrder.GetSubaccountId(),
		},
	}

	// Apply the update.
	success, successPerUpdate, err := k.subaccountsKeeper.UpdateSubaccounts(
		ctx,
		updates,
		satypes.Match,
	)
	if err != nil {
		return satypes.UpdateCausedError, satypes.UpdateCausedError, affiliateRevSharesQuoteQuantums, err
	}

	takerUpdateResult = successPerUpdate[0]
	makerUpdateResult = successPerUpdate[1]

	// If not successful, return error indicating why.
	if updateResultErr := satypes.GetErrorFromUpdateResults(
		success,
		successPerUpdate,
		updates,
	); updateResultErr != nil {
		return takerUpdateResult, makerUpdateResult, affiliateRevSharesQuoteQuantums, updateResultErr
	}

	if !success {
		panic(
			fmt.Sprintf(
				"persistMatchedOrders: UpdateSubaccounts failed but err == nil and no error returned"+
					"from successPerUpdate but success was false. Error: %v, Updates: %+v, SuccessPerUpdate: %+v",
				err,
				updates,
				successPerUpdate,
			),
		)
	}

	// TODO: get perpetual from perpetualId once and pass it to the functions that need the full
	// perpetual object. This will reduce the number of times we need to get the perpetual from the
	// keeper.

	perpetual, err := k.perpetualsKeeper.GetPerpetual(ctx, perpetualId)
	if err != nil {
		return takerUpdateResult, makerUpdateResult, affiliateRevSharesQuoteQuantums, err
	}

	// Transfer builder fees for taker and maker builders if they exist
	// Builder code fees are tranferred directly from the collateral pool to the
	// builder address because the builder fee is always taken out from
	// the trader's subaccount quote balance.
	if takerBuilderFeeQuantums.Sign() > 0 {
		if err := k.subaccountsKeeper.TransferBuilderFees(ctx,
			perpetualId,
			takerBuilderFeeQuantums,
			takerBuilderAddress,
		); err != nil {
			return takerUpdateResult, makerUpdateResult, affiliateRevSharesQuoteQuantums, err
		}
	}
	if makerBuilderFeeQuantums.Sign() > 0 {
		if err := k.subaccountsKeeper.TransferBuilderFees(ctx,
			perpetualId,
			makerBuilderFeeQuantums,
			makerBuilderAddress,
		); err != nil {
			return takerUpdateResult, makerUpdateResult, affiliateRevSharesQuoteQuantums, err
		}
	}

	fillForProcess := types.FillForProcess{
		TakerAddr:             matchWithOrders.TakerOrder.GetSubaccountId().Owner,
		TakerFeeQuoteQuantums: bigTakerFeeQuoteQuantums,
		MakerAddr:             matchWithOrders.MakerOrder.GetSubaccountId().Owner,
		MakerFeeQuoteQuantums: bigMakerFeeQuoteQuantums,
		FillQuoteQuantums:     bigFillQuoteQuantums,
		ProductId:             perpetualId,
		MarketId:              perpetual.Params.MarketId,
		MonthlyRollingTakerVolumeQuantums: k.statsKeeper.GetUserStats(
			ctx,
			matchWithOrders.TakerOrder.GetSubaccountId().Owner,
		).TakerNotional,
		TakerOrderRouterAddr: matchWithOrders.TakerOrder.GetOrderRouterAddress(),
		MakerOrderRouterAddr: matchWithOrders.MakerOrder.GetOrderRouterAddress(),
	}

	// Distribute the fee amount from subacounts module to fee collector and rev share accounts
	bigTotalFeeQuoteQuantums := new(big.Int).Add(bigTakerFeeQuoteQuantums, bigMakerFeeQuoteQuantums)
	revSharesForFill, err := k.revshareKeeper.GetAllRevShares(
		ctx,
		fillForProcess,
		affiliateOverrides,
		affiliateParameters,
	)
	if err != nil {
		revSharesForFill = revsharetypes.RevSharesForFill{}
		log.ErrorLogWithError(ctx, "error getting rev shares for fill", err)
	}
	if revSharesForFill.AffiliateRevShare != nil {
		affiliateRevSharesQuoteQuantums = revSharesForFill.AffiliateRevShare.QuoteQuantums
	}
	if err := k.subaccountsKeeper.DistributeFees(
		ctx,
		assettypes.AssetUsdc.Id,
		revSharesForFill,
		fillForProcess,
	); err != nil {
		return takerUpdateResult, makerUpdateResult, affiliateRevSharesQuoteQuantums, errorsmod.Wrapf(
			types.ErrSubaccountFeeTransferFailed,
			"persistMatchedOrders: subaccounts (%v, %v) updated, but fee transfer (bigFeeQuoteQuantums: %v)"+
				" to fee-collector failed. Err: %v",
			matchWithOrders.MakerOrder.GetSubaccountId(),
			matchWithOrders.TakerOrder.GetSubaccountId(),
			bigTotalFeeQuoteQuantums,
			err,
		)
	}

	// Update the last trade price for the perpetual.
	k.SetTradePricesForPerpetual(ctx, perpetualId, matchWithOrders.MakerOrder.GetOrderSubticks())

	// Process fill in x/stats and x/rewards.
	k.rewardsKeeper.AddRewardSharesForFill(
		ctx,
		fillForProcess,
		revSharesForFill,
	)

	attributableVolumeAttributions := k.buildAttributableVolumeAttributions(
		ctx,
		revSharesForFill,
		bigFillQuoteQuantums,
		matchWithOrders,
		affiliateParameters,
	)

	k.statsKeeper.RecordFill(
		ctx,
		matchWithOrders.TakerOrder.GetSubaccountId().Owner,
		matchWithOrders.MakerOrder.GetSubaccountId().Owner,
		bigFillQuoteQuantums,
		affiliateRevSharesQuoteQuantums,
		attributableVolumeAttributions,
	)

	takerOrderRouterFeeQuoteQuantums := big.NewInt(0)
	makerOrderRouterFeeQuoteQuantums := big.NewInt(0)
	for _, revShare := range revSharesForFill.AllRevShares {
		if revShare.Recipient == matchWithOrders.TakerOrder.GetOrderRouterAddress() &&
			revShare.RevShareType == revsharetypes.REV_SHARE_TYPE_ORDER_ROUTER {
			takerOrderRouterFeeQuoteQuantums.Add(takerOrderRouterFeeQuoteQuantums, revShare.QuoteQuantums)
		}
		if revShare.Recipient == matchWithOrders.MakerOrder.GetOrderRouterAddress() &&
			revShare.RevShareType == revsharetypes.REV_SHARE_TYPE_ORDER_ROUTER {
			makerOrderRouterFeeQuoteQuantums.Add(makerOrderRouterFeeQuoteQuantums, revShare.QuoteQuantums)
		}
	}

	matchWithOrders.MakerOrderRouterFee = makerOrderRouterFeeQuoteQuantums.Uint64()
	matchWithOrders.TakerOrderRouterFee = takerOrderRouterFeeQuoteQuantums.Uint64()

	// Emit an event indicating a match occurred.
	ctx.EventManager().EmitEvent(
		types.NewCreateMatchEvent(
			matchWithOrders.TakerOrder.GetSubaccountId(),
			matchWithOrders.MakerOrder.GetSubaccountId(),
			bigTakerFeeQuoteQuantums,
			bigMakerFeeQuoteQuantums,
			bigTakerQuoteBalanceDelta,
			bigMakerQuoteBalanceDelta,
			bigTakerPerpetualQuantumsDelta,
			bigMakerPerpetualQuantumsDelta,
			insuranceFundDelta,
			isTakerLiquidation,
			false,
			perpetualId,
			takerBuilderAddress,
			makerBuilderAddress,
			takerBuilderFeeQuantums,
			makerBuilderFeeQuantums,
			matchWithOrders.TakerOrder.GetOrderRouterAddress(),
			matchWithOrders.MakerOrder.GetOrderRouterAddress(),
			takerOrderRouterFeeQuoteQuantums,
			makerOrderRouterFeeQuoteQuantums,
		),
	)

	return takerUpdateResult, makerUpdateResult, affiliateRevSharesQuoteQuantums, nil
}

// getAttributableVolume calculates the attributable volume for a referee based on their
// already-attributed volume in the last 30 days and the maximum attributable volume cap.
// This does not modify any state.
func (k Keeper) getAttributableVolume(
	ctx sdk.Context,
	referee string,
	volume uint64,
	affiliateParameters affiliatetypes.AffiliateParameters,
) uint64 {
	// Get the user stats from the referee
	refereeUserStats := k.statsKeeper.GetUserStats(ctx, referee)
	if refereeUserStats == nil {
		return 0
	}

	// Use the ATTRIBUTED volume (how much has already been attributed to their affiliate)
	// NOT total trading volume (TakerNotional + MakerNotional)
	previouslyAttributedVolume := refereeUserStats.Affiliate_30DAttributedVolumeQuoteQuantums

	// If parameter is 0 then no limit is applied
	cap := affiliateParameters.Maximum_30DAttributableVolumePerReferredUserQuoteQuantums
	if cap == 0 {
		return volume
	}

	if previouslyAttributedVolume >= cap {
		return 0
	} else if previouslyAttributedVolume+volume > cap {
		// Remainder of the volume to get them to the cap
		return cap - previouslyAttributedVolume
	}

	return volume
}

func (k Keeper) buildAttributableVolumeAttributions(
	ctx sdk.Context,
	revSharesForFill revsharetypes.RevSharesForFill,
	bigFillQuoteQuantums *big.Int,
	matchWithOrders *types.MatchWithOrders,
	affiliateParameters affiliatetypes.AffiliateParameters,
) []*statstypes.AffiliateAttribution {
	// Build affiliate revenue attributions array (can include both taker and maker)
	var affiliateRevenueAttributions []*statstypes.AffiliateAttribution

	// Add taker affiliate attribution if present
	if revSharesForFill.AffiliateRevShare != nil &&
		revSharesForFill.AffiliateRevShare.Recipient != "" &&
		bigFillQuoteQuantums.Sign() > 0 {
		// Calculate the attributable volume based on the taker's current 30-day volume
		// and the maximum attributable volume cap from affiliate parameters
		takerAttributableVolume := k.getAttributableVolume(
			ctx,
			matchWithOrders.TakerOrder.GetSubaccountId().Owner,
			bigFillQuoteQuantums.Uint64(),
			affiliateParameters,
		)
		if takerAttributableVolume > 0 {
			affiliateRevenueAttributions = append(affiliateRevenueAttributions, &statstypes.AffiliateAttribution{
				Role:                        statstypes.AffiliateAttribution_ROLE_TAKER,
				ReferrerAddress:             revSharesForFill.AffiliateRevShare.Recipient,
				ReferredVolumeQuoteQuantums: takerAttributableVolume,
			})
		}
	}

	// Add maker affiliate attribution if present
	// Check if maker has an affiliate referrer
	makerReferrer, makerHasReferrer := k.affiliatesKeeper.GetReferredBy(
		ctx,
		matchWithOrders.MakerOrder.GetSubaccountId().Owner,
	)
	if makerHasReferrer && makerReferrer != "" && bigFillQuoteQuantums.Sign() > 0 {
		// Calculate the attributable volume based on the maker's current 30-day volume
		// and the maximum attributable volume cap from affiliate parameters
		makerAttributableVolume := k.getAttributableVolume(
			ctx,
			matchWithOrders.MakerOrder.GetSubaccountId().Owner,
			bigFillQuoteQuantums.Uint64(),
			affiliateParameters,
		)
		if makerAttributableVolume > 0 {
			affiliateRevenueAttributions = append(affiliateRevenueAttributions, &statstypes.AffiliateAttribution{
				Role:                        statstypes.AffiliateAttribution_ROLE_MAKER,
				ReferrerAddress:             makerReferrer,
				ReferredVolumeQuoteQuantums: makerAttributableVolume,
			})
		}
	}

	return affiliateRevenueAttributions
}

func (k Keeper) setOrderFillAmountsAndPruning(
	ctx sdk.Context,
	order types.Order,
	newTotalFillAmount satypes.BaseQuantums,
	curPruneableBlockHeight uint32,
) {
	// Note that stateful orders are never pruned by `BlockHeight`, so we set the value to `math.MaxUint32` here.
	pruneableBlockHeight := uint32(math.MaxUint32)

	if !order.IsStatefulOrder() {
		// Compute the block at which this state fill amount can be pruned. This is the greater of
		// `GoodTilBlock + ShortBlockWindow` and the existing `pruneableBlockHeight`.
		pruneableBlockHeight = lib.Max(
			order.GetGoodTilBlock()+types.ShortBlockWindow,
			curPruneableBlockHeight,
		)

		// Note: We should always prune out orders using the latest `GoodTilBlock` seen. It's possible there could be
		// multiple `GoodTilBlock`s for the same `OrderId` given order replacements. We would generally expect to see
		// the same `OrderId` with a lower `GoodTilBlock` first if the proposer is using this unmodified application,
		// but it's still not necessarily guaranteed due to MEV.
		if curPruneableBlockHeight > order.GetGoodTilBlock()+types.ShortBlockWindow {
			log.InfoLog(
				ctx,
				"Found an `orderId` in ProcessProposerMatches which had a lower GoodTilBlock than"+
					" a previous order in the list of fills. This could mean a lower priority order was allowed on the book.",
				"orderId",
				order.OrderId,
			)
		}

		// Add this order for pruning at the desired block height.
		k.AddOrdersForPruning(ctx, []types.OrderId{order.OrderId}, pruneableBlockHeight)
	}

	// Update the state with the new `fillAmount` for this `orderId`.
	// TODO(DEC-1219): Determine whether we should use `OrderFillState` proto for stateful order fill amounts.
	k.SetOrderFillAmount(
		ctx,
		order.OrderId,
		newTotalFillAmount,
		pruneableBlockHeight,
	)
}

// getUpdatedOrderFillAmount accepts an order's current total fill amount, total base quantums, and a new fill amount,
// and returns an error if the new fill amount would cause the order to exceed its base quantums.
// Returns the new total fill amount of the order.
func getUpdatedOrderFillAmount(
	orderId types.OrderId,
	orderBaseQuantums satypes.BaseQuantums,
	currentFillAmount satypes.BaseQuantums,
	fillQuantums satypes.BaseQuantums,
) (satypes.BaseQuantums, error) {
	bigCurrentFillAmount := currentFillAmount.ToBigInt()
	bigNewFillAmount := bigCurrentFillAmount.Add(bigCurrentFillAmount, fillQuantums.ToBigInt())
	if bigNewFillAmount.Cmp(orderBaseQuantums.ToBigInt()) == 1 {
		return 0, errorsmod.Wrapf(
			types.ErrInvalidMsgProposedOperations,
			"Match with Quantums %v would exceed total Quantums %v of OrderId %v. New total filled quantums would be %v.",
			fillQuantums,
			orderBaseQuantums,
			orderId,
			bigNewFillAmount.String(),
		)
	}

	return satypes.BaseQuantums(bigNewFillAmount.Uint64()), nil
}

func (k Keeper) checkAndUpdateTWAPOrderRemainingQuantity(
	ctx sdk.Context,
	orderId types.OrderId,
	fillAmount satypes.BaseQuantums,
) error {
	if orderId.IsTwapSuborder() {
		parentOrderId := types.OrderId{
			SubaccountId: orderId.SubaccountId,
			ClientId:     orderId.ClientId,
			OrderFlags:   types.OrderIdFlags_Twap, // Set directly to TWAP
			ClobPairId:   orderId.ClobPairId,
		}
		if err := k.UpdateTWAPOrderRemainingQuantityOnFill(ctx, parentOrderId, fillAmount.ToUint64()); err != nil {
			return err
		}
	}
	return nil
}
