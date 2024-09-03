package constants

import (
	"math"

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
	SubaccountBlockLimits_Default = clobtypes.SubaccountBlockLimits{
		MaxQuantumsInsuranceLost: 100_000_000_000_000,
	}
	SubaccountBlockLimits_No_Limit = clobtypes.SubaccountBlockLimits{
		MaxQuantumsInsuranceLost: math.MaxUint64,
	}
	// Liquidation Configs.
	LiquidationsConfig_No_Limit = clobtypes.LiquidationsConfig{
		InsuranceFundFeePpm:   5_000,
		ValidatorFeePpm:       200_000,
		LiquidityFeePpm:       800_000,
		FillablePriceConfig:   FillablePriceConfig_Default,
		SubaccountBlockLimits: SubaccountBlockLimits_No_Limit,
	}
	LiquidationsConfig_FillablePrice_Max_Smmr = clobtypes.LiquidationsConfig{
		InsuranceFundFeePpm:   5_000,
		ValidatorFeePpm:       200_000,
		LiquidityFeePpm:       800_000,
		FillablePriceConfig:   FillablePriceConfig_Max_Smmr,
		SubaccountBlockLimits: SubaccountBlockLimits_No_Limit,
	}
	LiquidationsConfig_Position_Min10m_Max05mPpm = clobtypes.LiquidationsConfig{
		InsuranceFundFeePpm:   5_000,
		ValidatorFeePpm:       200_000,
		LiquidityFeePpm:       800_000,
		FillablePriceConfig:   FillablePriceConfig_Default,
		SubaccountBlockLimits: SubaccountBlockLimits_No_Limit,
	}
	LiquidationsConfig_Subaccount_Max10bNotionalLiquidated_Max10bInsuranceLost = clobtypes.LiquidationsConfig{
		InsuranceFundFeePpm: 5_000,
		ValidatorFeePpm:     200_000,
		LiquidityFeePpm:     800_000,
		FillablePriceConfig: FillablePriceConfig_Default,
		SubaccountBlockLimits: clobtypes.SubaccountBlockLimits{
			MaxQuantumsInsuranceLost: 10_000_000_000, // $10,000
		},
	}
)
