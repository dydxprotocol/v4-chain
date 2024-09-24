package types

import (
	errorsmod "cosmossdk.io/errors"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/lib"
)

// Validate validates each individual field of the liquidations config for validity.
// It returns an error if any of the liquidation config fields fail the following validation:
// - `InsuranceFundFeePpm == 0 || InsuranceFundFeePpm > 1_000_000`.
// - `ValidatorFeePpm > 1_000_000`.
// - `LiquidityFeePpm > 1_000_000`.
// - `bankruptcyAdjustmentPpm < 1_000_000`.
// - `spreadToMaintenanceMarginRatioPpm == 0.

func (lc *LiquidationsConfig) Validate() error {

	// Validate the InsuranceFundFeePpm.
	if lc.InsuranceFundFeePpm == 0 || lc.InsuranceFundFeePpm > lib.OneMillion {
		return errorsmod.Wrapf(
			ErrInvalidLiquidationsConfig,
			"%v is not a valid InsuranceFundFeePpm",
			lc.InsuranceFundFeePpm,
		)
	}

	// Validate the ValidatorFeePpm.
	if lc.ValidatorFeePpm > lib.OneMillion {
		return errorsmod.Wrapf(
			ErrInvalidLiquidationsConfig,
			"%v is not a valid ValidatorFeePpm",
			lc.ValidatorFeePpm,
		)
	}

	// Validate the LiquidityFeePpm.
	if lc.LiquidityFeePpm > lib.OneMillion {
		return errorsmod.Wrapf(
			ErrInvalidLiquidationsConfig,
			"%v is not a valid LiquidityFeePpm",
			lc.LiquidityFeePpm,
		)
	}

	if lc.ValidatorFeePpm+lc.LiquidityFeePpm > lib.OneMillion {
		return errorsmod.Wrapf(
			ErrInvalidLiquidationsConfig,
			"%v is not a valid sum of ValidatorFeePpm and LiquidityFeePpm",
			lc.ValidatorFeePpm+lc.LiquidityFeePpm,
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

	// Validate the MaxCumulativeInsuranceFundDelta.
	if lc.MaxCumulativeInsuranceFundDelta == 0 {
		return errorsmod.Wrapf(
			ErrInvalidLiquidationsConfig,
			"%v is not a valid MaxCumulativeInsuranceFundDelta, it must be greater than 0",
			lc.MaxCumulativeInsuranceFundDelta,
		)
	}

	return nil
}
