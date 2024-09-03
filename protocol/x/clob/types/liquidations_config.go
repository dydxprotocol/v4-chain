package types

import (
	errorsmod "cosmossdk.io/errors"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/lib"
)

// Validate validates each individual field of the liquidations config for validity.
// It returns an error if any of the liquidation config fields fail the following validation:
// - `maxLiquidationFee == 0 || maxLiquidationFee > 1_000_000`.
// - `bankruptcyAdjustmentPpm < 1_000_000`.
// - `spreadToMaintenanceMarginRatioPpm == 0.
// - `maxQuantumsInsuranceLost == 0`.

func (lc *LiquidationsConfig) Validate() error {

	// Validate the MaxLiquidationFeePpm.
	if lc.MaxLiquidationFeePpm == 0 || lc.MaxLiquidationFeePpm > lib.OneMillion {
		return errorsmod.Wrapf(
			ErrInvalidLiquidationsConfig,
			"%v is not a valid MaxLiquidationFeePpm",
			lc.MaxLiquidationFeePpm,
		)
	}

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
