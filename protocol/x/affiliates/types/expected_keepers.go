package types

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type StatsKeeper interface {
	GetStakedAmount(ctx sdk.Context, delegatorAddr string) big.Int
}
