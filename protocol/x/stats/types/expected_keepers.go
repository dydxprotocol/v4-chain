package types

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	epochtypes "github.com/dydxprotocol/v4-chain/protocol/x/epochs/types"
)

// EpochsKeeper defines the expected epochs keeper to get epoch info.
type EpochsKeeper interface {
	MustGetStatsEpochInfo(ctx sdk.Context) epochtypes.EpochInfo
}

type StakingKeeper interface {
	GetDelegatorDelegations(ctx context.Context,
		delegator sdk.AccAddress, maxRetrieve uint16) ([]stakingtypes.Delegation, error)
}

// StatsExpirationHook is called when stats are expired from the rolling window
type StatsExpirationHook interface {
	OnStatsExpired(ctx sdk.Context, userAddress string, resultingUserStats *UserStats) error
}
