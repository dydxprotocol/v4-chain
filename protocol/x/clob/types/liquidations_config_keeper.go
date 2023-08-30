package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
)

var (
	LiquidationsConfig_Default = LiquidationsConfig{
		MaxLiquidationFeePpm: 5_000,
		FillablePriceConfig: FillablePriceConfig{
			BankruptcyAdjustmentPpm:           lib.OneMillion,
			SpreadToMaintenanceMarginRatioPpm: 100_000,
		},
		PositionBlockLimits: PositionBlockLimits{
			MinPositionNotionalLiquidated:   1_000,
			MaxPositionPortionLiquidatedPpm: 1_000_000,
		},
		SubaccountBlockLimits: SubaccountBlockLimits{
			MaxNotionalLiquidated:    100_000_000_000_000,
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
