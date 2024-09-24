package constants

import (
	"github.com/StreamFinance-Protocol/stream-chain/protocol/lib"
	clobtypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/clob/types"
)

var (
	// Block limits.
	FillablePriceConfig_Default = clobtypes.FillablePriceConfig{
		BankruptcyAdjustmentPpm:           lib.OneMillion,
		SpreadToMaintenanceMarginRatioPpm: 100_000,
	}
	FillablePriceConfig_Max_Smmr = clobtypes.FillablePriceConfig{
		BankruptcyAdjustmentPpm:           lib.OneMillion,
		SpreadToMaintenanceMarginRatioPpm: lib.OneMillion,
	}
	// Liquidation Configs.
	LiquidationsConfig_No_Limit = clobtypes.LiquidationsConfig{
		InsuranceFundFeePpm:             5_000,
		ValidatorFeePpm:                 0,
		LiquidityFeePpm:                 0,
		FillablePriceConfig:             FillablePriceConfig_Default,
		MaxCumulativeInsuranceFundDelta: uint64(1_000_000_000_000),
	}
	LiquidationsConfig_FillablePrice_Max_Smmr = clobtypes.LiquidationsConfig{
		InsuranceFundFeePpm:             5_000,
		ValidatorFeePpm:                 0,
		LiquidityFeePpm:                 0,
		FillablePriceConfig:             FillablePriceConfig_Max_Smmr,
		MaxCumulativeInsuranceFundDelta: uint64(1_000_000_000_000),
	}
	LiquidationsConfig_FillablePrice_Max_Smmr_With_Fees = clobtypes.LiquidationsConfig{
		InsuranceFundFeePpm:             5_000,
		ValidatorFeePpm:                 100_000,
		LiquidityFeePpm:                 400_000,
		FillablePriceConfig:             FillablePriceConfig_Max_Smmr,
		MaxCumulativeInsuranceFundDelta: uint64(1_000_000_000_000),
	}
	LiquidationsConfig_Position_Min10m_Max05mPpm = clobtypes.LiquidationsConfig{
		InsuranceFundFeePpm:             5_000,
		ValidatorFeePpm:                 200_000,
		LiquidityFeePpm:                 800_000,
		FillablePriceConfig:             FillablePriceConfig_Default,
		MaxCumulativeInsuranceFundDelta: uint64(1_000_000_000_000),
	}
	LiquidationsConfig_Subaccount_Max10bNotionalLiquidated_Max10bInsuranceLost = clobtypes.LiquidationsConfig{
		InsuranceFundFeePpm:             5_000,
		ValidatorFeePpm:                 200_000,
		LiquidityFeePpm:                 800_000,
		FillablePriceConfig:             FillablePriceConfig_Default,
		MaxCumulativeInsuranceFundDelta: uint64(1_000_000_000_000),
	}
)
