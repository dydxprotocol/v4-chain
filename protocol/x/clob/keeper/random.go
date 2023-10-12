package keeper

import (
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetPseudoRand returns a random number generator seeded with a pseudorandom seed.
// The seed is based on the previous block timestamp.
func (k *Keeper) GetPseudoRand(ctx sdk.Context) *rand.Rand {
	s := rand.NewSource(
		k.blockTimeKeeper.GetPreviousBlockInfo(ctx).Timestamp.Unix(),
	)
	return rand.New(s)
}
