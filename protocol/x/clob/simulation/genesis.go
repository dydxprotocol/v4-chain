package simulation

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"

	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/dydxprotocol/v4/lib"
	"github.com/dydxprotocol/v4/testutil/sim_helpers"
	"github.com/dydxprotocol/v4/x/clob/types"
	perptypes "github.com/dydxprotocol/v4/x/perpetuals/types"
)

// genNumClobPairs returns a randomized number of CLOB pairs.
func genNumClobPairs(r *rand.Rand, isReasonableGenesis bool, numPerpetuals int) int {
	if isReasonableGenesis {
		return numPerpetuals
	}
	return simtypes.RandIntBetween(
		r,
		sim_helpers.MinValidClobPairs,
		sim_helpers.MaxValidClobPairs+1,
	)
}

// genRandomClob returns a CLOB pair with randomized parameters.
func genRandomClob(
	r *rand.Rand,
	isReasonableGenesis bool,
	clobPairId types.ClobPairId,
	perpetualId uint32,
) types.ClobPair {
	var clobPair types.ClobPair

	clobPair.QuantumConversionExponent = int32(simtypes.RandIntBetween(
		r,
		sim_helpers.PickGenesisParameter(sim_helpers.MinQuantumConversionExponent, isReasonableGenesis),
		sim_helpers.PickGenesisParameter(sim_helpers.MaxQuantumConversionExponent, isReasonableGenesis)+1,
	))
	clobPair.TakerFeePpm = uint32(simtypes.RandIntBetween(
		r,
		sim_helpers.PickGenesisParameter(sim_helpers.MinFeePpm, isReasonableGenesis),
		sim_helpers.PickGenesisParameter(sim_helpers.MaxFeePpm, isReasonableGenesis)+1,
	))
	clobPair.StepBaseQuantums = uint64(simtypes.RandIntBetween(
		r,
		sim_helpers.PickGenesisParameter(sim_helpers.MinStepBaseQuantums, isReasonableGenesis),
		sim_helpers.PickGenesisParameter(sim_helpers.MaxStepBaseQuantums, isReasonableGenesis)+1,
	))
	clobPair.MinOrderBaseQuantums = alignToStepBaseQuantums(
		clobPair.StepBaseQuantums,
		uint64(simtypes.RandIntBetween(
			r,
			sim_helpers.PickGenesisParameter(sim_helpers.MinOrderBaseQuantums, isReasonableGenesis),
			sim_helpers.PickGenesisParameter(sim_helpers.MaxOrderBaseQuantums, isReasonableGenesis)+1,
		)),
	)
	clobPair.SubticksPerTick = uint32(simtypes.RandIntBetween(
		r,
		sim_helpers.PickGenesisParameter(sim_helpers.MinSubticksPerTick, isReasonableGenesis),
		sim_helpers.PickGenesisParameter(sim_helpers.MaxSubticksPerTick, isReasonableGenesis)+1,
	))

	clobPair.Id = clobPairId.ToUint32()
	clobPair.MakerFeePpm = uint32(
		simtypes.RandIntBetween(
			r,
			sim_helpers.PickGenesisParameter(sim_helpers.MinFeePpm, isReasonableGenesis),
			int(clobPair.TakerFeePpm)+1,
		),
	)

	perpetualClobMetadata := createPerpetualClobMetadata(perpetualId)
	clobPair.Metadata = &perpetualClobMetadata

	// TODO(DEC-977): Specify `Status` in `RandomizedGenState`.
	clobPair.Status = types.ClobPair_STATUS_ACTIVE

	return clobPair
}

// alignToStepBaseQuantums takes a value and aligns it such that `n % stepBaseQuantums == 0`.
// It will do this by rounding towards the smallest multiple of `stepBaseQuantums` that is
// greater than zero.
func alignToStepBaseQuantums(
	stepBaseQuantums uint64,
	seedMinOrderBaseQuantums uint64,
) uint64 {
	// `MinOrderBaseQuantums` cannot be smaller than `StepBaseQuantums`.
	if seedMinOrderBaseQuantums < stepBaseQuantums {
		return stepBaseQuantums
	}

	return seedMinOrderBaseQuantums - (seedMinOrderBaseQuantums % stepBaseQuantums)
}

// createPerpetualClobMetadata returns a `PerpetualClobMetadata`.
func createPerpetualClobMetadata(perpetualId uint32) types.ClobPair_PerpetualClobMetadata {
	perpetualClobMetadata := types.ClobPair_PerpetualClobMetadata{
		PerpetualClobMetadata: &types.PerpetualClobMetadata{
			PerpetualId: perpetualId,
		},
	}

	return perpetualClobMetadata
}

// genClobPairToPerpetualSlice returns a slice of length `numClobPairs`, where each index
// corresponds to a `ClobPair.Id` and entry is the `Perpetual.Id` that should be assigned to the
// CLOB pair.
func genClobPairToPerpetualSlice(r *rand.Rand, numClobPairs, numPerpetuals int) []uint32 {
	perpetuals := sim_helpers.MakeRange(uint32(numPerpetuals))

	// Add additional perpetuals if there are more CLOB pairs than perpetuals.
	if numClobPairs > numPerpetuals {
		diff := numClobPairs - numPerpetuals
		extraPerpetuals := make([]uint32, diff)
		for i := 0; i < diff; i++ {
			randomIdx := simtypes.RandIntBetween(r, 0, numPerpetuals)
			extraPerpetuals[i] = perpetuals[randomIdx]
		}
		perpetuals = append(perpetuals, extraPerpetuals...)
	}

	// Shuffle perpetuals, so we randomize which `ClobPair` gets matched with which `Perpetual`.
	r.Shuffle(numPerpetuals, func(i, j int) { perpetuals[i], perpetuals[j] = perpetuals[j], perpetuals[i] })

	return perpetuals
}

// genRandomPositivePpm returns a random positive parts-per-million value.
func genRandomPositivePpm(r *rand.Rand, skewTowardsLower bool) uint32 {
	if skewTowardsLower {
		return uint32(sim_helpers.GetRandomBucketValue(r, sim_helpers.PpmSkewedTowardLowerBuckets))
	}
	return uint32(sim_helpers.GetRandomBucketValue(r, sim_helpers.PpmSkewedTowardLargerBuckets))
}

// RandomizedGenState generates a random GenesisState for `CLOB`.
func RandomizedGenState(simState *module.SimulationState) {
	r := simState.Rand

	isReasonableGenesis := sim_helpers.ShouldGenerateReasonableGenesis(r, simState.GenTimestamp)
	clobGenesis := types.GenesisState{}

	// Get number of perpetuals.
	cdc := codec.NewProtoCodec(codectypes.NewInterfaceRegistry())
	perpGenesisBytes := simState.GenState[perptypes.ModuleName]
	var perpetualsGenesis perptypes.GenesisState
	if err := cdc.UnmarshalJSON(perpGenesisBytes, &perpetualsGenesis); err != nil {
		panic(fmt.Sprintf("Could not unmarshal Perpetuals GenesisState %s", err))
	}
	numPerpetuals := len(perpetualsGenesis.Perpetuals)
	if numPerpetuals == 0 {
		panic("Number of Perpetuals cannot be zero")
	}

	// Generate number of CLOB pairs.
	numClobPairs := genNumClobPairs(r, isReasonableGenesis, numPerpetuals)

	// Generate `ClobPair` to `Perpetual` slice.
	clobPairToPerpetual := genClobPairToPerpetualSlice(r, numClobPairs, numPerpetuals)

	clobPairs := make([]types.ClobPair, numClobPairs)
	for i := 0; i < numClobPairs; i++ {
		clobPairId := types.ClobPairId(i)
		// TODO(DEC-1039): Allow generating a random spot CLOB pair.
		clobPair := genRandomClob(r, isReasonableGenesis, clobPairId, clobPairToPerpetual[clobPairId])
		clobPairs[i] = clobPair
	}

	clobGenesis.ClobPairs = clobPairs

	clobGenesis.LiquidationsConfig = types.LiquidationsConfig{
		MaxInsuranceFundQuantumsForDeleveraging: uint64(
			sim_helpers.GetRandomBucketValue(r, sim_helpers.MaxInsuranceFundQuantumsForDeleveragingBuckets),
		),
		// MaxLiquidationFeePpm determines the fee that subaccount usually pays for liquidating a position.
		// This is typically a very small percentage, so skewing towards lower values here.
		MaxLiquidationFeePpm: genRandomPositivePpm(r, true),
		FillablePriceConfig: types.FillablePriceConfig{
			BankruptcyAdjustmentPpm: uint32(
				simtypes.RandIntBetween(r, int(lib.OneMillion), int(math.MaxUint32)),
			),
			// SpreadToMaintenanceMarginRatioPpm represents the maximum liquidation spread
			// in the fillable price calculation.
			// This is typically also a small percentage to protect against MEV,
			// so skewing towards lower values here.
			SpreadToMaintenanceMarginRatioPpm: genRandomPositivePpm(r, true),
		},
		PositionBlockLimits: types.PositionBlockLimits{
			MinPositionNotionalLiquidated: uint64(sim_helpers.GetRandomBucketValue(r, sim_helpers.MinPositionNotionalBuckets)),
			// MaxPositionPortionLiquidatedPpm determines the maximum portion of a position
			// that can be liquidated in a block.
			// Since we may want to liquidate as quickly as possible to avoid losing any insurance fund,
			// skewing towards larger values here.
			MaxPositionPortionLiquidatedPpm: genRandomPositivePpm(r, false),
		},
		SubaccountBlockLimits: types.SubaccountBlockLimits{
			MaxNotionalLiquidated:    uint64(sim_helpers.GetRandomBucketValue(r, sim_helpers.SubaccountBlockLimitsBuckets)),
			MaxQuantumsInsuranceLost: uint64(sim_helpers.GetRandomBucketValue(r, sim_helpers.SubaccountBlockLimitsBuckets)),
		},
	}

	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&clobGenesis)
}
