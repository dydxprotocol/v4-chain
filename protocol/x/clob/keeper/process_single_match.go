package keeper

import (
	errorsmod "cosmossdk.io/errors"
	"errors"
	"fmt"
	"math"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/off_chain_updates"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
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
) (
	success bool,
	takerUpdateResult satypes.UpdateResult,
	makerUpdateResult satypes.UpdateResult,
	offchainUpdates *types.OffchainUpdates,
	err error,
) {
	if matchWithOrders.TakerOrder.IsLiquidation() {
		defer func() {
			if errors.Is(err, satypes.ErrFailedToUpdateSubaccounts) && !takerUpdateResult.IsSuccess() {
				takerSubaccount := k.subaccountsKeeper.GetSubaccount(ctx, matchWithOrders.TakerOrder.GetSubaccountId())
				takerTnc, takerIMR, takerMMR, _ := k.subaccountsKeeper.GetNetCollateralAndMarginRequirements(
					ctx,
					satypes.Update{SubaccountId: *takerSubaccount.Id},
				)
				k.Logger(ctx).Error(
					"collateralization check failed for liquidation",
					"takerSubaccount", fmt.Sprintf("%+v", takerSubaccount),
					"takerTNC", takerTnc,
					"takerIMR", takerIMR,
					"takerMMR", takerMMR,
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
		return false, takerUpdateResult, makerUpdateResult, nil, errorsmod.Wrapf(
			err,
			"ProcessSingleMatch: Invalid MatchWithOrders: %+v",
			matchWithOrders,
		)
	}

	offchainUpdates = types.NewOffchainUpdates()
	makerMatchableOrder := matchWithOrders.MakerOrder
	takerMatchableOrder := matchWithOrders.TakerOrder
	fillAmount := matchWithOrders.FillAmount

	// Retrieve the ClobPair from state.
	clobPairId := makerMatchableOrder.GetClobPairId()
	clobPair, found := k.GetClobPair(ctx, clobPairId)
	if !found {
		return false, takerUpdateResult, makerUpdateResult, nil, types.ErrInvalidClob
	}

	// Verify that the `fillAmount` is divisible by the `StepBaseQuantums` of the `clobPair`.
	if fillAmount.ToUint64()%clobPair.StepBaseQuantums != 0 {
		return false,
			takerUpdateResult,
			makerUpdateResult,
			nil,
			types.ErrFillAmountNotDivisibleByStepSize
	}

	// Define local variable relevant to retrieving QuoteQuantums based on the fill amount.
	makerSubticks := makerMatchableOrder.GetOrderSubticks()

	// Calculate the number of quote quantums for the match based on the maker order subticks.
	bigFillQuoteQuantums, err := getFillQuoteQuantums(clobPair, makerSubticks, fillAmount)
	if err != nil {
		return false, takerUpdateResult, makerUpdateResult, nil, err
	}

	if bigFillQuoteQuantums.Sign() == 0 {
		// Note: If `subticks`, `baseQuantums`, are small enough, `quantumConversionExponent` is negative,
		// it's possible to have zero `quoteQuantums` for a non-zero amount of `baseQuantums`.
		// This could mean that it's possible that a maker sell order on the book
		// at a very unfavorable price (subticks) could receive `0` `quoteQuantums` amount.
		k.Logger(ctx).Error(
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
	perpetualId, err := clobPair.GetPerpetualId()
	if err != nil {
		return false, takerUpdateResult, makerUpdateResult, nil, err
	}

	// Calculate taker and maker fee ppms.
	takerFeePpm := k.feeTiersKeeper.GetPerpetualFeePpm(
		ctx, matchWithOrders.TakerOrder.GetSubaccountId().Owner, true)
	makerFeePpm := k.feeTiersKeeper.GetPerpetualFeePpm(
		ctx, matchWithOrders.MakerOrder.GetSubaccountId().Owner, false)

	takerInsuranceFundDelta := new(big.Int)
	if takerMatchableOrder.IsLiquidation() {
		// Liquidation orders do not take trading fees because they already pay a liquidation fee.
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
			return false, takerUpdateResult, makerUpdateResult, nil, err
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
			return false, takerUpdateResult, makerUpdateResult, nil, err
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
		return false, takerUpdateResult, makerUpdateResult, nil, err
	}

	// Update both subaccounts in the matched order atomically.
	takerUpdateResult, makerUpdateResult, err = k.persistMatchedOrders(
		ctx,
		matchWithOrders,
		perpetualId,
		takerFeePpm,
		makerFeePpm,
		bigFillQuoteQuantums,
		takerInsuranceFundDelta,
	)

	if err != nil {
		return false, takerUpdateResult, makerUpdateResult, nil, err
	}

	// Update subaccount total quantums liquidated and total insurance fund lost for liquidation orders.
	if matchWithOrders.TakerOrder.IsLiquidation() {
		notionalLiquidatedQuoteQuantums, err := k.perpetualsKeeper.GetNetNotional(
			ctx,
			perpetualId,
			fillAmount.ToBigInt(),
		)
		if err != nil {
			return false, takerUpdateResult, makerUpdateResult, nil, err
		}

		k.UpdateSubaccountLiquidationInfo(
			ctx,
			matchWithOrders.TakerOrder.GetSubaccountId(),
			notionalLiquidatedQuoteQuantums,
			takerInsuranceFundDelta,
		)
	}

	// Liquidation orders can only be placed when a subaccount is liquidatable
	// and cannot be replayed, therefore we don't need to track their filled amount in state.
	if !matchWithOrders.TakerOrder.IsLiquidation() {
		takerOffchainUpdates := k.setOrderFillAmountsAndPruning(
			ctx,
			matchWithOrders.TakerOrder.MustGetOrder(),
			newTakerTotalFillAmount,
			curTakerPruneableBlockHeight,
		)
		offchainUpdates.Append(takerOffchainUpdates)
	}

	makerOffchainUpdates := k.setOrderFillAmountsAndPruning(
		ctx,
		matchWithOrders.MakerOrder.MustGetOrder(),
		newMakerTotalFillAmount,
		curMakerPruneableBlockHeight,
	)
	offchainUpdates.Append(makerOffchainUpdates)

	return true, takerUpdateResult, makerUpdateResult, offchainUpdates, nil
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
) (
	takerUpdateResult satypes.UpdateResult,
	makerUpdateResult satypes.UpdateResult,
	err error,
) {
	isTakerLiquidation := matchWithOrders.TakerOrder.IsLiquidation()
	bigTakerFeeQuoteQuantums := lib.BigIntMulSignedPpm(bigFillQuoteQuantums, takerFeePpm)
	bigMakerFeeQuoteQuantums := lib.BigIntMulSignedPpm(bigFillQuoteQuantums, makerFeePpm)
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

	// Create the subaccount update.
	updates := []satypes.Update{
		// Taker update
		{
			AssetUpdates: []satypes.AssetUpdate{
				{
					AssetId:          lib.UsdcAssetId,
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
					AssetId:          lib.UsdcAssetId,
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
	)
	if err != nil {
		return satypes.UpdateCausedError, satypes.UpdateCausedError, err
	}

	takerUpdateResult = successPerUpdate[0]
	makerUpdateResult = successPerUpdate[1]

	// If not successful, return error indicating why.
	if updateResultErr := satypes.GetErrorFromUpdateResults(
		success,
		successPerUpdate,
		updates,
	); updateResultErr != nil {
		return takerUpdateResult, makerUpdateResult, updateResultErr
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

	if err := k.subaccountsKeeper.TransferInsuranceFundPayments(ctx, insuranceFundDelta); err != nil {
		return takerUpdateResult, makerUpdateResult, err
	}

	// Transfer the fee amount from subacounts module to fee collector module account.
	bigTotalFeeQuoteQuantums := new(big.Int).Add(bigTakerFeeQuoteQuantums, bigMakerFeeQuoteQuantums)
	if err := k.subaccountsKeeper.TransferFeesToFeeCollectorModule(
		ctx,
		lib.UsdcAssetId,
		bigTotalFeeQuoteQuantums,
	); err != nil {
		return takerUpdateResult, makerUpdateResult, errorsmod.Wrapf(
			types.ErrSubaccountFeeTransferFailed,
			"persistMatchedOrders: subaccounts (%v, %v) updated, but fee transfer (bigFeeQuoteQuantums: %v)"+
				" to fee-collector failed. Err: %v",
			matchWithOrders.MakerOrder.GetSubaccountId(),
			matchWithOrders.TakerOrder.GetSubaccountId(),
			bigTotalFeeQuoteQuantums,
			err,
		)
	}

	// Process fill in x/stats and x/rewards.
	k.rewardsKeeper.AddRewardSharesForFill(
		ctx,
		matchWithOrders.TakerOrder.GetSubaccountId().Owner,
		matchWithOrders.MakerOrder.GetSubaccountId().Owner,
		bigFillQuoteQuantums,
		bigTakerFeeQuoteQuantums,
		bigMakerFeeQuoteQuantums,
	)

	k.statsKeeper.RecordFill(
		ctx,
		matchWithOrders.TakerOrder.GetSubaccountId().Owner,
		matchWithOrders.MakerOrder.GetSubaccountId().Owner,
		bigFillQuoteQuantums,
	)

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
		),
	)

	return takerUpdateResult, makerUpdateResult, nil
}

func (k Keeper) setOrderFillAmountsAndPruning(
	ctx sdk.Context,
	order types.Order,
	newTotalFillAmount satypes.BaseQuantums,
	curPruneableBlockHeight uint32,
) *types.OffchainUpdates {
	// Note that stateful orders are never pruned by `BlockHeight`, so we set the value to `MaxUInt32` here.
	pruneableBlockHeight := uint32(math.MaxUint32)
	offchainUpdates := types.NewOffchainUpdates()

	if !order.IsStatefulOrder() {
		// Compute the block at which this state fill amount can be pruned. This is the greater of
		// `GoodTilBlock + ShortBlockWindow` and the existing `pruneableBlockHeight`.
		pruneableBlockHeight = lib.MaxUint32(
			order.GetGoodTilBlock()+types.ShortBlockWindow,
			curPruneableBlockHeight,
		)

		// Note: We should always prune out orders using the latest `GoodTilBlock` seen. It's possible there could be
		// multiple `GoodTilBlock`s for the same `OrderId` given order replacements. We would generally expect to see
		// the same `OrderId` with a lower `GoodTilBlock` first if the proposer is using this unmodified application,
		// but it's still not necessarily guaranteed due to MEV.
		if curPruneableBlockHeight > order.GetGoodTilBlock()+types.ShortBlockWindow {
			k.Logger(ctx).Info(
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

	if k.GetIndexerEventManager().Enabled() {
		if _, exists := k.MemClob.GetOrder(ctx, order.OrderId); exists {
			// Generate an off-chain update message updating the total filled amount of order.
			if message, success := off_chain_updates.CreateOrderUpdateMessage(
				k.Logger(ctx),
				order.OrderId,
				newTotalFillAmount,
			); success {
				offchainUpdates.AddUpdateMessage(order.OrderId, message)
			}
		}
	}

	return offchainUpdates
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
