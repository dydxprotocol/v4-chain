package types

import (
	context "context"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// StatsKeeper defines the expected stats keeper
type StakingKeeper interface {
	Slash(
		ctx context.Context,
		consAddr sdk.ConsAddress,
		infractionHeight,
		power int64,
		slashFactor math.LegacyDec,
	) (math.Int, error)
}
