package keeper

import (
	"fmt"
	"math/big"

	"github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

// Helper function to compute the delta long for a single settled update on a perpetual.
func getDeltaLongFromSettledUpdate(
	u SettledUpdate,
	updatedPerpId uint32,
) (
	deltaLong *big.Int,
) {
	var perpPosition *types.PerpetualPosition
	for _, p := range u.SettledSubaccount.PerpetualPositions {
		// TODO use a pre-populated map
		if p.PerpetualId == updatedPerpId {
			perpPosition = p
		}
	}

	prevQuantums := perpPosition.GetBigQuantums()
	afterQuantums := new(big.Int).Add(
		prevQuantums,
		u.PerpetualUpdates[0].GetBigQuantums(),
	)

	prevLong := prevQuantums // re-use pointer for efficiency
	if prevLong.Sign() < 0 {
		prevLong.SetUint64(0)
	}
	afterLong := afterQuantums // re-use pointer for efficiency
	if afterLong.Sign() < 0 {
		afterLong.SetUint64(0)
	}

	return afterLong.Sub(
		afterLong,
		prevLong,
	)
}

// Returns the delta open_interest for a pair of Match updates if they were applied.
func GetDeltaOpenInterestFromPerpMatchUpdates(
	settledUpdates []SettledUpdate,
) (
	updatedPerpId uint32,
	deltaOpenInterest *big.Int,
) {
	if len(settledUpdates) != 2 {
		panic(
			fmt.Sprintf(
				types.ErrMatchUpdatesMustHaveTwoUpdates,
				settledUpdates,
			),
		)
	}

	if len(settledUpdates[0].PerpetualUpdates) != 1 || len(settledUpdates[1].PerpetualUpdates) != 1 {
		panic(
			fmt.Sprintf(
				types.ErrMatchUpdatesMustUpdateOnePerp,
				settledUpdates,
			),
		)
	}

	if settledUpdates[0].PerpetualUpdates[0].PerpetualId != settledUpdates[1].PerpetualUpdates[0].PerpetualId {
		panic(
			fmt.Sprintf(
				types.ErrMatchUpdatesMustBeSamePerpId,
				settledUpdates,
			),
		)
	} else {
		updatedPerpId = settledUpdates[0].PerpetualUpdates[0].PerpetualId
	}

	deltaOpenInterest = big.NewInt(0)
	for _, u := range settledUpdates {
		deltaLong := getDeltaLongFromSettledUpdate(u, updatedPerpId)
		deltaOpenInterest.Add(deltaOpenInterest, deltaLong)
	}

	return updatedPerpId, deltaOpenInterest
}
