package keeper

import (
	errorsmod "cosmossdk.io/errors"
	"errors"
	"fmt"
	"math/big"
	"math/rand"
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
	balance *big.Int,
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

	// Return as big.Int.
	return insuranceFundBalance.Amount.BigInt()
}

// CanDeleverageSubaccount returns true if a subaccount can be deleveraged.
// Specifically, this function returns true if both of the following are true:
// - The insurance fund balance is less-than-or-equal to `MaxInsuranceFundQuantumsForDeleveraging`.
// - The subaccount's total net collateral is negative.
// This function returns an error if `GetNetCollateralAndMarginRequirements` returns an error.
func (k Keeper) CanDeleverageSubaccount(
	ctx sdk.Context,
	subaccountId satypes.SubaccountId,
) (bool, error) {
	currentInsuranceFundBalance := k.GetInsuranceFundBalance(ctx)
	liquidationConfig := k.GetLiquidationsConfig(ctx)
	bigMaxInsuranceFundForDeleveraging := new(big.Int).SetUint64(liquidationConfig.MaxInsuranceFundQuantumsForDeleveraging)

	// Deleveraging cannot be performed if the current insurance fund balance is greater than the
	// max insurance fund for deleveraging,
	if currentInsuranceFundBalance.Cmp(bigMaxInsuranceFundForDeleveraging) > 0 {
		return false, nil
	}

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

	// The insurance fund balance is less-than-or-equal to `MaxInsuranceFundQuantumsForDeleveraging`
	// and the subaccount's total net collateral is negative, so deleveraging can be performed.
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

	numSubaccountsIterated := 0
	deltaQuantumsRemaining = new(big.Int).Set(deltaQuantumsTotal)
	fills = make([]types.MatchPerpetualDeleveraging_Fill, 0)

	s := rand.NewSource(k.blockTimeKeeper.GetPreviousBlockInfo(ctx).Timestamp.Unix())
	rand := rand.New(s)

	k.subaccountsKeeper.ForEachSubaccountRandomStart(
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
			} else if errors.Is(err, types.ErrInvalidPerpetualPositionSizeDelta) {
				panic(
					fmt.Sprintf(
						"Invalid perpetual position size delta when processing deleveraging. error: %v",
						err,
					),
				)
			} else {
				// If an error is returned, it's likely because the subaccounts' bankruptcy prices do not overlap.
				liquidatedSubaccount := k.subaccountsKeeper.GetSubaccount(ctx, liquidatedSubaccountId)
				liquidatedBankruptcyPrice, bankruptcyPriceError := k.GetBankruptcyPriceInQuoteQuantums(
					ctx,
					liquidatedSubaccountId,
					perpetualId,
					deltaQuantums,
				)
				if bankruptcyPriceError != nil {
					k.Logger(ctx).Error(
						"error when getting bankruptcy price for liquidated subaccount",
						"error", bankruptcyPriceError,
						"blockHeight", ctx.BlockHeight(),
						"checkTx", ctx.IsCheckTx(),
						"perpetualId", perpetualId,
						"deltaQuantums", deltaQuantums,
					)
					return false
				}
				liquidatedTnc, _, _, tncErr := k.subaccountsKeeper.GetNetCollateralAndMarginRequirements(
					ctx, satypes.Update{SubaccountId: *liquidatedSubaccount.Id},
				)
				if tncErr != nil {
					k.Logger(ctx).Error(
						"error when getting TNC for liquidated subaccount",
						"error", tncErr,
						"blockHeight", ctx.BlockHeight(),
						"checkTx", ctx.IsCheckTx(),
						"perpetualId", perpetualId,
						"deltaQuantums", deltaQuantums,
					)
					return false
				}

				offsettingSubaccount := k.subaccountsKeeper.GetSubaccount(ctx, *offsettingSubaccount.Id)
				offsettingBankruptcyPrice, bankruptcyPriceError := k.GetBankruptcyPriceInQuoteQuantums(
					ctx,
					*offsettingSubaccount.Id,
					perpetualId,
					new(big.Int).Neg(deltaQuantums),
				)
				if bankruptcyPriceError != nil {
					k.Logger(ctx).Error(
						"error when getting bankruptcy price for offsetting subaccount",
						"error", bankruptcyPriceError,
						"blockHeight", ctx.BlockHeight(),
						"checkTx", ctx.IsCheckTx(),
						"perpetualId", perpetualId,
						"deltaQuantums", deltaQuantums,
					)
					return false
				}
				offsettingTnc, _, _, tncErr := k.subaccountsKeeper.GetNetCollateralAndMarginRequirements(
					ctx, satypes.Update{SubaccountId: *offsettingSubaccount.Id},
				)
				if tncErr != nil {
					k.Logger(ctx).Error(
						"error when getting TNC for offsetting subaccount",
						"error", tncErr,
						"blockHeight", ctx.BlockHeight(),
						"checkTx", ctx.IsCheckTx(),
						"perpetualId", perpetualId,
						"deltaQuantums", deltaQuantums,
					)
					return false
				}

				k.Logger(ctx).Info(
					"Encountered error when processing deleveraging",
					"error", err,
					"blockHeight", ctx.BlockHeight(),
					"checkTx", ctx.IsCheckTx(),
					"perpetualId", perpetualId,
					"deltaQuantums", deltaQuantums,
					"liquidatedSubaccount", fmt.Sprintf("%+v", liquidatedSubaccount),
					"liquidatedBankruptcyPriceQuoteQuantums", liquidatedBankruptcyPrice,
					"liquidatedTnc", liquidatedTnc,
					"offsettingSubaccount", fmt.Sprintf("%+v", offsettingSubaccount),
					"offsettingBankruptcyPriceQuoteQuantums", offsettingBankruptcyPrice,
					"offsettingTnc", offsettingTnc,
				)
				telemetry.IncrCounter(
					1,
					types.ModuleName, metrics.Deleveraging, metrics.NonOverlappingBankruptcyPrices, metrics.Count,
				)
			}
			return deltaQuantumsRemaining.Sign() == 0
		},
		rand,
	)

	telemetry.SetGauge(float32(numSubaccountsIterated), metrics.NumSubaccountsIterated, metrics.Count)

	if deltaQuantumsRemaining.Sign() == 0 {
		// Deleveraging was successful.
		telemetry.IncrCounter(1, types.ModuleName, metrics.CheckTx, metrics.Deleveraging, metrics.Success, metrics.Count)
	} else {
		// Not enough offsetting subaccounts to fully offset the liquidated subaccount's position.
		telemetry.IncrCounter(
			1,
			types.ModuleName, metrics.CheckTx, metrics.Deleveraging, metrics.NotEnoughPositionToFullyOffset, metrics.Count,
		)
		k.Logger(ctx).Error(
			errorsmod.Wrapf(
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
					AssetId:          lib.UsdcAssetId,
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
					AssetId:          lib.UsdcAssetId,
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
