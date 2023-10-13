package constants

import (
	"math"

	"github.com/dydxprotocol/v4-chain/protocol/lib"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
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
	PositionBlockLimits_Default = clobtypes.PositionBlockLimits{
		MinPositionNotionalLiquidated:   1_000,
		MaxPositionPortionLiquidatedPpm: 1_000_000,
	}
	SubaccountBlockLimits_Default = clobtypes.SubaccountBlockLimits{
		MaxNotionalLiquidated:    100_000_000_000_000,
		MaxQuantumsInsuranceLost: 100_000_000_000_000,
	}
	PositionBlockLimits_No_Limit = clobtypes.PositionBlockLimits{
		MinPositionNotionalLiquidated:   1,
		MaxPositionPortionLiquidatedPpm: lib.OneMillion,
	}
	SubaccountBlockLimits_No_Limit = clobtypes.SubaccountBlockLimits{
		MaxNotionalLiquidated:    math.MaxUint64,
		MaxQuantumsInsuranceLost: math.MaxUint64,
	}
	// Liquidation Configs.
	LiquidationsConfig_No_Limit = clobtypes.LiquidationsConfig{
		MaxLiquidationFeePpm:  5_000,
		FillablePriceConfig:   FillablePriceConfig_Default,
		PositionBlockLimits:   PositionBlockLimits_No_Limit,
		SubaccountBlockLimits: SubaccountBlockLimits_No_Limit,
	}
	LiquidationsConfig_FillablePrice_Max_Smmr = clobtypes.LiquidationsConfig{
		MaxLiquidationFeePpm:  5_000,
		FillablePriceConfig:   FillablePriceConfig_Max_Smmr,
		PositionBlockLimits:   PositionBlockLimits_No_Limit,
		SubaccountBlockLimits: SubaccountBlockLimits_No_Limit,
	}
	LiquidationsConfig_Position_Min10m_Max05mPpm = clobtypes.LiquidationsConfig{
		MaxLiquidationFeePpm: 5_000,
		FillablePriceConfig:  FillablePriceConfig_Default,
		PositionBlockLimits: clobtypes.PositionBlockLimits{
			MinPositionNotionalLiquidated:   10_000_000, // $10
			MaxPositionPortionLiquidatedPpm: 500_000,
		},
		SubaccountBlockLimits: SubaccountBlockLimits_No_Limit,
	}
	LiquidationsConfig_Subaccount_Max10bNotionalLiquidated_Max10bInsuranceLost = clobtypes.LiquidationsConfig{
		MaxLiquidationFeePpm: 5_000,
		FillablePriceConfig:  FillablePriceConfig_Default,
		PositionBlockLimits:  PositionBlockLimits_No_Limit,
		SubaccountBlockLimits: clobtypes.SubaccountBlockLimits{
			MaxNotionalLiquidated:    10_000_000_000, // $10,000
			MaxQuantumsInsuranceLost: 10_000_000_000, // $10,000
		},
	}
)
