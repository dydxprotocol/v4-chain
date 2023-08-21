package keeper

import (
	"math/big"
	"math/rand"

	gometrics "github.com/armon/go-metrics"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

// GetInsuranceFundBalance returns the current balance of the insurance fund (in quote quantums).
// This calls the Bank Keeper’s GetBalance() function for the Module Address of the insurance fund.
func (k Keeper) GetInsuranceFundBalance(
	ctx sdk.Context,
) (
	balance uint64,
) {
	usdcAsset, err := k.assetsKeeper.GetAsset(ctx, lib.UsdcAssetId)
	if err != nil {
		panic("GetInsuranceFundBalance: Usdc asset not found in state")
	}

	insuranceFundBalance := k.bankKeeper.GetBalance(
		ctx,
		authtypes.NewModuleAddress(types.InsuranceFundName),
		usdcAsset.Denom,
	)

	floatBalance, _ := new(big.Float).SetUint64(insuranceFundBalance.Amount.Uint64()).Float32()
	telemetry.ModuleSetGauge(
		types.ModuleName,
		floatBalance,
		metrics.InsuranceFundBalance,
	)
	// Return the amount as uint64. `Uint64` panics if amount
	// cannot be represented in a uint64.
	return insuranceFundBalance.Amount.Uint64()
}

// ShouldPerformDeleveraging returns true if deleveraging needs to occur.
// Specifically, this function returns true if both of the following are true:
// - The `insuranceFundDelta` is negative.
// - The insurance fund balance is less than `MaxInsuranceFundQuantumsForDeleveraging` or `abs(insuranceFundDelta)`.
func (k Keeper) ShouldPerformDeleveraging(
	ctx sdk.Context,
	insuranceFundDelta *big.Int,
) (
	shouldPerformDeleveraging bool,
) {
	if insuranceFundDelta.Sign() >= 0 {
		return false
	}

	currentInsuranceFundBalance := new(big.Int).SetUint64(k.GetInsuranceFundBalance(ctx))

	liquidationConfig := k.GetLiquidationsConfig(ctx)
	bigMaxInsuranceFundForDeleveraging := new(big.Int).SetUint64(liquidationConfig.MaxInsuranceFundQuantumsForDeleveraging)

	return new(big.Int).Add(currentInsuranceFundBalance, insuranceFundDelta).Sign() < 0 ||
		currentInsuranceFundBalance.Cmp(bigMaxInsuranceFundForDeleveraging) < 0
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
	numSubaccountsIterated := 0
	deltaQuantumsRemaining = new(big.Int).Set(deltaQuantumsTotal)
	fills = make([]types.MatchPerpetualDeleveraging_Fill, 0)

	s := rand.NewSource(k.MustGetBlockTimeForLastCommittedBlock(ctx).Unix())
	rand := rand.New(s)

	// TODO(DEC-1487): Determine how offsetting subaccounts should be selected.
	k.subaccountsKeeper.ForEachSubaccountFromRandomStart(
		ctx,
		func(offsettingSubaccount satypes.Subaccount) (finished bool) {
			numSubaccountsIterated++
			offsettingPosition, _ := offsettingSubaccount.GetPerpetualPositionForId(perpetualId)
			bigOffsettingPositionQuantums := offsettingPosition.GetBigQuantums()

			// Skip subaccounts that do not have a position in the opposite direction as the liquidated subaccount.
			if deltaQuantumsRemaining.Sign() != bigOffsettingPositionQuantums.Sign() {
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
			} else {
				// If an error is returned, then the subaccounts' bankruptcy prices do not overlap.
				telemetry.IncrCounterWithLabels(
					[]string{types.ModuleName, metrics.Deleveraging, metrics.NonOverlappingBankruptcyPrices, metrics.Count},
					1,
					[]gometrics.Label{metrics.GetLabelForIntValue(metrics.BlockHeight, int(ctx.BlockHeight()))},
				)
			}
			return deltaQuantumsRemaining.Sign() == 0
		},
		rand,
	)

	telemetry.ModuleSetGauge(
		types.ModuleName,
		float32(numSubaccountsIterated),
		metrics.NumSubaccountsIterated,
		metrics.Count,
	)

	mode := metrics.DeliverTx
	blockHeight := ctx.BlockHeight()
	if ctx.IsCheckTx() {
		mode = metrics.CheckTx
		blockHeight += 1
	}

	if deltaQuantumsRemaining.Sign() == 0 {
		// Deleveraging was successful.
		telemetry.IncrCounterWithLabels(
			[]string{types.ModuleName, mode, metrics.Deleveraging, metrics.Success, metrics.Count},
			1,
			[]gometrics.Label{
				metrics.GetLabelForIntValue(metrics.BlockHeight, int(blockHeight)),
			},
		)
	} else {
		// Not enough offsetting subaccounts to fully offset the liquidated subaccount's position.
		telemetry.IncrCounterWithLabels(
			[]string{types.ModuleName, mode, metrics.Deleveraging, metrics.NotEnoughPositionToFullyOffset, metrics.Count},
			1,
			[]gometrics.Label{
				metrics.GetLabelForIntValue(metrics.BlockHeight, int(blockHeight)),
			},
		)
		ctx.Logger().Error(
			sdkerrors.Wrapf(
				types.ErrPositionCannotBeFullyOffset,
				"OffsetSubaccountPerpetualPosition: Not enough position to fully offset position, "+
					"subaccount = (%+v), perpetual = (%d), quantums remaining = (%+v)",
				liquidatedSubaccountId,
				perpetualId,
				deltaQuantumsRemaining.String(),
			).Error(),
		)
		// TODO(CLOB-75): Support deleveraging subaccounts with non overlapping bankruptcy prices.
	}

	return fills, deltaQuantumsRemaining
}

// ProcessDeleveraging processes a deleveraging operation by closing both the liquidated subaccount's
// position and the offsetting subaccount's position at the bankruptcy price of the _liquidated_ position.
// This function takes a `deltaQuantums` argument, which is the delta with respect to the liquidated subaccount's
// position, to allow for partial deleveraging.
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
		return sdkerrors.Wrapf(
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

	updates := []satypes.Update{
		// Liquidated subaccount update.
		{
			AssetUpdates: []satypes.AssetUpdate{
				{
					AssetId:          lib.UsdcAssetId,
					BigQuantumsDelta: bankruptcyPriceQuoteQuantums,
				},
			},
			PerpetualUpdates: []satypes.PerpetualUpdate{
				{
					PerpetualId:      perpetualId,
					BigQuantumsDelta: deltaQuantums,
				},
			},
			SubaccountId: liquidatedSubaccountId,
		},
		// Offsetting subaccount update.
		{
			AssetUpdates: []satypes.AssetUpdate{
				{
					AssetId:          lib.UsdcAssetId,
					BigQuantumsDelta: new(big.Int).Neg(bankruptcyPriceQuoteQuantums),
				},
			},
			PerpetualUpdates: []satypes.PerpetualUpdate{
				{
					PerpetualId:      perpetualId,
					BigQuantumsDelta: new(big.Int).Neg(deltaQuantums),
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
	return satypes.GetErrorFromUpdateResults(success, successPerUpdate, updates)
}
