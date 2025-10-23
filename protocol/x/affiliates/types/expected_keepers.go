package types

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stattypes "github.com/dydxprotocol/v4-chain/protocol/x/stats/types"
)

type StatsKeeper interface {
	GetStakedAmount(ctx sdk.Context, delegatorAddr string) *big.Int
	GetBlockStats(ctx sdk.Context) *stattypes.BlockStats
	GetUserStats(ctx sdk.Context, address string) *stattypes.UserStats
	SetUserStats(ctx sdk.Context, address string, userStats *stattypes.UserStats)
}

type FeetiersKeeper interface {
	GetAffiliateRefereeLowestTakerFee(ctx sdk.Context) int32
	GetLowestMakerFee(ctx sdk.Context) int32
}
