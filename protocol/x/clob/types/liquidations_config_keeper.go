package types

import (
	"github.com/StreamFinance-Protocol/stream-chain/protocol/lib"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	LiquidationsConfig_Default = LiquidationsConfig{
		InsuranceFundFeePpm: 5_000,
		ValidatorFeePpm:     200_000,
		LiquidityFeePpm:     800_000,
		FillablePriceConfig: FillablePriceConfig{
			BankruptcyAdjustmentPpm:           lib.OneMillion,
			SpreadToMaintenanceMarginRatioPpm: 100_000,
		},
		MaxCumulativeInsuranceFundDelta: uint64(1_000_000_000_000),
	}
)

// LiquidationsConfigKeeper is an interface that encapsulates all reads and writes to the
// liquidation configuration values written to state.
type LiquidationsConfigKeeper interface {
	GetLiquidationsConfig(
		ctx sdk.Context,
	) LiquidationsConfig
}
