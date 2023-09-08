package types

import (
	errorsmod "cosmossdk.io/errors"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
)

// Validate validates each individual field of the liquidations config for validity.
// It returns an error if any of the liquidation config fields fail the following validation:
// - `maxLiquidationFee == 0 || maxLiquidationFee > 1_000_000`.
// - `bankruptcyAdjustmentPpm < 1_000_000`.
// - `spreadToMaintenanceMarginRatioPpm == 0.
// - `maxPositionPortionLiquidatedPpm == 0 || maxPositionPortionLiquidatedPpm > 1_000_000`.
// - `maxNotionalLiquidated == 0`.
// - `maxQuantumsInsuranceLost == 0`.
//
// Note that `minPositionNotionalLiquidated` is intentionally not validated.

func (lc *LiquidationsConfig) Validate() error {
	// Validate the BankruptcyAdjustmentPpm.
	bankruptcyAdjustmentPpm := lc.FillablePriceConfig.BankruptcyAdjustmentPpm
	if bankruptcyAdjustmentPpm < lib.OneMillion {
		return errorsmod.Wrapf(
			ErrInvalidLiquidationsConfig,
			"%v is not a valid BankruptcyAdjustmentPpm",
			bankruptcyAdjustmentPpm,
		)
	}

	// Validate the SpreadToMaintenanceMarginRatioPpm.
	spreadToMaintenanceMarginRatioPpm := lc.FillablePriceConfig.SpreadToMaintenanceMarginRatioPpm
	if spreadToMaintenanceMarginRatioPpm == 0 {
		return errorsmod.Wrapf(
			ErrInvalidLiquidationsConfig,
			"%v is not a valid SpreadToMaintenanceMarginRatioPpm",
			spreadToMaintenanceMarginRatioPpm,
		)
	}

	// Validate the MaxLiquidationFeePpm.
	if lc.MaxLiquidationFeePpm == 0 || lc.MaxLiquidationFeePpm > lib.OneMillion {
		return errorsmod.Wrapf(
			ErrInvalidLiquidationsConfig,
			"%v is not a valid MaxLiquidationFeePpm",
			lc.MaxLiquidationFeePpm,
		)
	}

	// Validate the MaxPositionPortionLiquidatedPpm.
	maxPositionPortionLiquidatedPpm := lc.PositionBlockLimits.MaxPositionPortionLiquidatedPpm
	if maxPositionPortionLiquidatedPpm == 0 || maxPositionPortionLiquidatedPpm > lib.OneMillion {
		return errorsmod.Wrapf(
			ErrInvalidLiquidationsConfig,
			"%v is not a valid MaxPositionPortionLiquidatedPpm",
			maxPositionPortionLiquidatedPpm,
		)
	}

	// Validate the MaxNotionalLiquidated.
	maxNotionalLiquidated := lc.SubaccountBlockLimits.MaxNotionalLiquidated
	if maxNotionalLiquidated == 0 {
		return errorsmod.Wrapf(
			ErrInvalidLiquidationsConfig,
			"%v is not a valid MaxNotionalLiquidated",
			maxNotionalLiquidated,
		)
	}

	// Validate the MaxQuantumsInsuranceLost.
	maxQuantumsInsuranceLost := lc.SubaccountBlockLimits.MaxQuantumsInsuranceLost
	if maxQuantumsInsuranceLost == 0 {
		return errorsmod.Wrapf(
			ErrInvalidLiquidationsConfig,
			"%v is not a valid MaxQuantumsInsuranceLost",
			maxQuantumsInsuranceLost,
		)
	}

	return nil
}
