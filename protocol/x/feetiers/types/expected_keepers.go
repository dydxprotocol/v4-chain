package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/stats/types"
)

// StatsKeeper defines the expected stats keeper
type StatsKeeper interface {
	GetUserStats(ctx sdk.Context, address string) *types.UserStats
	GetGlobalStats(ctx sdk.Context) *types.GlobalStats
}
