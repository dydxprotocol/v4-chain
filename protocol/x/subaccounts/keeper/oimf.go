package keeper

import (
	"fmt"

	"github.com/dydxprotocol/v4-chain/protocol/lib/int256"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

// Helper function to compute the delta long for a single settled update on a perpetual.
func getDeltaLongFromSettledUpdate(
	u SettledUpdate,
	updatedPerpId uint32,
) (
	deltaLong *int256.Int,
) {
	var perpPosition *types.PerpetualPosition
	for _, p := range u.SettledSubaccount.PerpetualPositions {
		// TODO use a pre-populated map
		if p.PerpetualId == updatedPerpId {
			perpPosition = p
		}
	}

	prevQuantums := perpPosition.GetQuantums()
	afterQuantums := new(int256.Int).Add(
		prevQuantums,
		u.PerpetualUpdates[0].GetQuantums(),
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

// For `Match` updates:
//   - returns a struct `OpenInterest` if input updates results in OI delta.
//   - returns nil if OI delta is zero.
//   - panics if update format is invalid.
//
// For other update types, returns nil.
func GetDeltaOpenInterestFromUpdates(
	settledUpdates []SettledUpdate,
	updateType types.UpdateType,
) (ret *perptypes.OpenInterestDelta) {
	if updateType != types.Match {
		return nil
	}

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

	perpUpdate0 := settledUpdates[0].PerpetualUpdates[0]
	perpUpdate1 := settledUpdates[1].PerpetualUpdates[0]

	if perpUpdate0.PerpetualId != perpUpdate1.PerpetualId {
		panic(
			fmt.Sprintf(
				types.ErrMatchUpdatesMustBeSamePerpId,
				settledUpdates,
			),
		)
	}

	updatedPerpId := settledUpdates[0].PerpetualUpdates[0].PerpetualId

	if (perpUpdate0.QuantumsDelta.Sign()*perpUpdate1.QuantumsDelta.Sign() > 0) ||
		perpUpdate0.QuantumsDelta.CmpAbs(perpUpdate1.QuantumsDelta) != 0 {
		panic(
			fmt.Sprintf(
				types.ErrMatchUpdatesInvalidSize,
				settledUpdates,
			),
		)
	}

	baseQuantumsDelta := int256.NewInt(0)
	for _, u := range settledUpdates {
		deltaLong := getDeltaLongFromSettledUpdate(u, updatedPerpId)
		baseQuantumsDelta.Add(
			baseQuantumsDelta,
			deltaLong,
		)
	}

	if baseQuantumsDelta.Sign() == 0 {
		return nil
	}

	return &perptypes.OpenInterestDelta{
		PerpetualId:  updatedPerpId,
		BaseQuantums: baseQuantumsDelta,
	}
}
