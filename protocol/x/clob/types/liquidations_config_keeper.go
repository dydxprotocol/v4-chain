package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	LiquidationsConfig_Default = LiquidationsConfig{
		MaxLiquidationFeePpm: 5_000,
		SubaccountBlockLimits: SubaccountBlockLimits{
			MaxQuantumsInsuranceLost: 100_000_000_000_000,
		},
	}
)

// LiquidationsConfigKeeper is an interface that encapsulates all reads and writes to the
// liquidation configuration values written to state.
type LiquidationsConfigKeeper interface {
	GetLiquidationsConfig(
		ctx sdk.Context,
	) LiquidationsConfig
}
