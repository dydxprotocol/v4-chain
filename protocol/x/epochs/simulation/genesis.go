package simulation

import (
	"math"
	"math/rand"

	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"

	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/sim_helpers"
	"github.com/dydxprotocol/v4-chain/protocol/x/epochs/types"
)

const (
	minNumEpochs = 1
	maxNumEpochs = 1000

	numSecInOneMinute   = 60
	numSecInFiveMinutes = numSecInOneMinute * 5
	numSecInOneHour     = numSecInOneMinute * 60 // = 60 * 60
	numSecInOneDay      = numSecInOneHour * 24   // = 60 * 60 * 24
)

var (
	durationBuckets = []int{
		1, // min
		numSecInOneMinute,
		numSecInFiveMinutes,
		numSecInOneHour,
		numSecInOneDay,
		math.MaxUint32 + 1, // max
	}
)

// genNumEpochs returns randomized num of epochs.
func genNumEpochs(r *rand.Rand) int {
	return simtypes.RandIntBetween(r, minNumEpochs, maxNumEpochs+1)
}

// genEpochName returns randomized epoch name.
func genEpochName(r *rand.Rand) string {
	return simtypes.RandStringOfLength(r, simtypes.RandIntBetween(r, 5, 20))
}

// genDuration returns randomized epoch duration.
func genDuration(r *rand.Rand) uint32 {
	return uint32(sim_helpers.GetRandomBucketValue(r, durationBuckets))
}

// TODO(DEC-957): improve how we pick the random value to better test "NextTick is before genesis time"
// vs. "NextTick is after the genesis time".
// genNextTick returns randomized epoch next tick.
func genNextTick(r *rand.Rand) uint32 {
	return uint32(simtypes.RandIntBetween(r, 0, math.MaxUint32+1))
}

// genCurrEpochAndStartBlock returns randomized current epoch and its start block.
// There are two rules:
// 1. `CurrentEpochStartBlock` == 0 if `CurrentEpoch` == 0
// 2. `CurrentEpochStartBlock` != 0 if `CurrentEpoch` != 0
func genCurrEpochAndStartBlock(r *rand.Rand) (currEpoch, currEpochStartBlock uint32) {
	shouldReturnZero := sim_helpers.RandBool(r)
	if shouldReturnZero {
		return 0, 0
	}
	epoch := uint32(simtypes.RandIntBetween(r, 1, math.MaxUint32+1))
	// Use 1 as the start block for current epoch so that `currentBlockHeight - currentEpochStartBlock` >= 0.
	// Otherwise, `NumBlockSinceEpochStart` will panic.
	// Currently all our simulations start at the default, block 1. Even if we configure simulations
	// to start at a later block, there's no easy to access the `initialBlockHeight` in this function.
	// TODO(DEC-1745): Generate random `currentEpochStartBlock` based on `initialBlockHeight` of simulation.
	epochStartBlock := uint32(1)
	return epoch, epochStartBlock
}

// genFastForwardNextTick returns randomized fast forward next tick.
func genFastForwardNextTick(r *rand.Rand) bool {
	return sim_helpers.RandBool(r)
}

// genIsInitialized returns randomized IsInitialized boolean.
func genIsInitialized(r *rand.Rand) bool {
	return sim_helpers.RandBool(r)
}

// RandomizedGenState generates a random GenesisState for `Epochs`.
func RandomizedGenState(simState *module.SimulationState) {
	r := simState.Rand
	numEpochs := genNumEpochs(r)

	defaultGenesis := types.DefaultGenesis()
	allEpochNames := make(map[string]bool)
	allEpochGenInfo := make([]types.EpochInfo, 0, numEpochs+len(defaultGenesis.EpochInfoList))

	// The `epoch`s defined in the default genesis are required by other modules (i.e. `Perpetuals`).
	for _, defaultEpochGenInfo := range defaultGenesis.GetEpochInfoList() {
		epochName := defaultEpochGenInfo.GetName()
		allEpochNames[epochName] = true

		currEpoch, currEpochStartBlock := genCurrEpochAndStartBlock(r)

		epochInfo := types.EpochInfo{
			Name:                   epochName,
			NextTick:               genNextTick(r),
			Duration:               genDuration(r),
			CurrentEpoch:           currEpoch,
			CurrentEpochStartBlock: currEpochStartBlock,
			FastForwardNextTick:    genFastForwardNextTick(r),
			IsInitialized:          genIsInitialized(r),
		}
		allEpochGenInfo = append(allEpochGenInfo, epochInfo)
	}

	for i := 0; i < numEpochs; i++ {
		// Generate a new unique epochName and add it to set for tracking.
		var epochName = genEpochName(r)
		for _, exists := allEpochNames[epochName]; exists; {
			epochName = genEpochName(r)
		}
		allEpochNames[epochName] = true

		currEpoch, currEpochStartBlock := genCurrEpochAndStartBlock(r)

		// Add epoch genesis info.
		epochInfoGen := types.EpochInfo{
			Name:                   epochName,
			NextTick:               genNextTick(r),
			Duration:               genDuration(r),
			CurrentEpoch:           currEpoch,
			CurrentEpochStartBlock: currEpochStartBlock,
			FastForwardNextTick:    genFastForwardNextTick(r),
		}
		allEpochGenInfo = append(allEpochGenInfo, epochInfoGen)
	}

	epochsGenesis := types.GenesisState{
		EpochInfoList: allEpochGenInfo,
	}

	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&epochsGenesis)
}
