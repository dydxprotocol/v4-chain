package keeper

import (
	"fmt"
	"math/big"

	gometrics "github.com/armon/go-metrics"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/dydxprotocol/v4/lib"
	"github.com/dydxprotocol/v4/lib/metrics"
	"github.com/dydxprotocol/v4/x/clob/types"
	satypes "github.com/dydxprotocol/v4/x/subaccounts/types"
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

// MustGetOffsettingSubaccountsForDeleveraging returns a list of subaccounts that can be used to
// offset the liquidated subaccount's position.
// This function assumes `deltaQuantums` is valid with respect to the subaccount being deleveraged.
func (k Keeper) MustGetOffsettingSubaccountsForDeleveraging(
	ctx sdk.Context,
	liquidatedSubaccountId satypes.SubaccountId,
	perpetualId uint32,
	deltaQuantums *big.Int,
) (
	offsettingSubaccounts []satypes.SubaccountId,
) {
	if deltaQuantums.Sign() == 0 {
		panic(
			fmt.Sprintf(
				"MustGetOffsettingSubaccountsForDeleveraging: deltaQuantums is zero. SubaccountId: (%+v). perpetualId %d",
				liquidatedSubaccountId,
				perpetualId,
			),
		)
	}

	// Verify that the subaccount to be deleveraged has a non-negative total net collateral.
	// TODO(DEC-1543): Support deleveraging for subaccounts with negative total net collateral.
	totalNetCollateral, _, _, err := k.subaccountsKeeper.GetNetCollateralAndMarginRequirements(
		ctx,
		satypes.Update{SubaccountId: liquidatedSubaccountId},
	)
	if err != nil {
		panic(err)
	}
	if totalNetCollateral.Sign() < 0 {
		telemetry.IncrCounterWithLabels(
			[]string{types.ModuleName, metrics.DeleveragedSubaccountWithNegativeTNC, metrics.Count},
			1,
			[]gometrics.Label{metrics.GetLabelForIntValue(metrics.BlockHeight, int(ctx.BlockHeight()))},
		)
		ctx.Logger().Debug(
			fmt.Sprintf(
				"MustGetOffsettingSubaccountsForDeleveraging: Subaccount %+v has negative "+
					"total net collateral and cannot be deleveraged",
				liquidatedSubaccountId,
			),
		)
		return offsettingSubaccounts
	}

	// TODO(DEC-1487): Determine how offsetting subaccounts should be selected.
	for i, subaccount := range k.subaccountsKeeper.GetAllSubaccount(ctx) {
		offsettingPosition, _ := subaccount.GetPerpetualPositionForId(perpetualId)
		bigOffsettingPositionQuantums := offsettingPosition.GetBigQuantums()

		// Skip subaccounts that do not have a position in the opposite direction as the liquidated subaccount.
		if deltaQuantums.Sign() != bigOffsettingPositionQuantums.Sign() {
			continue
		}

		totalNetCollateral, _, _, err := k.subaccountsKeeper.GetNetCollateralAndMarginRequirements(
			ctx,
			satypes.Update{SubaccountId: *subaccount.Id},
		)
		if err != nil {
			panic(err)
		}

		// Skip subaccounts with negative total net collateral.
		if totalNetCollateral.Sign() < 0 {
			telemetry.IncrCounter(1, types.ModuleName, metrics.OffsettingSubaccountWithNegativeTNC, metrics.Count)
			continue
		}

		offsettingSubaccounts = append(offsettingSubaccounts, *subaccount.Id)

		// Return the offsetting subaccounts if the liquidatable position has been fully offset.
		if deltaQuantums.CmpAbs(bigOffsettingPositionQuantums) <= 0 {
			telemetry.ModuleSetGauge(types.ModuleName, float32(i)+1, metrics.NumSubaccountsIterated, metrics.Count)
			return offsettingSubaccounts
		}

		// Update the remaining liquidatable quantums.
		// TODO(DEC-1495): Determine how much to offset per subaccount.
		deltaQuantums = new(big.Int).Sub(
			deltaQuantums,
			bigOffsettingPositionQuantums,
		)
	}

	// Not enough offsetting subaccounts to fully offset the liquidated subaccount's position.
	telemetry.IncrCounter(1, types.ModuleName, metrics.NotEnoughPositionToFullyOffset, metrics.Count)
	return offsettingSubaccounts
}
