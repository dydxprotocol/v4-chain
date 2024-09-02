package constants

import (
	"math"

	clobtypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/clob/types"
)

var (
	// Block limits.
	SubaccountBlockLimits_Default = clobtypes.SubaccountBlockLimits{
		MaxQuantumsInsuranceLost: 100_000_000_000_000,
	}
	SubaccountBlockLimits_No_Limit = clobtypes.SubaccountBlockLimits{
		MaxQuantumsInsuranceLost: math.MaxUint64,
	}
	// Liquidation Configs.
	LiquidationsConfig_No_Limit = clobtypes.LiquidationsConfig{
		MaxLiquidationFeePpm:  5_000,
		SubaccountBlockLimits: SubaccountBlockLimits_No_Limit,
	}
	LiquidationsConfig_FillablePrice_Max_Smmr = clobtypes.LiquidationsConfig{
		MaxLiquidationFeePpm:  5_000,
		SubaccountBlockLimits: SubaccountBlockLimits_No_Limit,
	}
	LiquidationsConfig_Position_Min10m_Max05mPpm = clobtypes.LiquidationsConfig{
		MaxLiquidationFeePpm:  5_000,
		SubaccountBlockLimits: SubaccountBlockLimits_No_Limit,
	}
	LiquidationsConfig_Subaccount_Max10bNotionalLiquidated_Max10bInsuranceLost = clobtypes.LiquidationsConfig{
		MaxLiquidationFeePpm: 5_000,
		SubaccountBlockLimits: clobtypes.SubaccountBlockLimits{
			MaxQuantumsInsuranceLost: 10_000_000_000, // $10,000
		},
	}
)
