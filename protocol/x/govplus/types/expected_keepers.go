package types

import (
	context "context"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
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

	GetAllValidators(
		ctx context.Context,
	) ([]stakingtypes.Validator, error)

	PowerReduction(ctx context.Context) math.Int
}
