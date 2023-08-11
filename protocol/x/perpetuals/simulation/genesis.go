package simulation

// DONTCOVER

import (
	"fmt"
	"math/rand"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"

	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"

	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/dydxprotocol/v4/dtypes"
	"github.com/dydxprotocol/v4/testutil/sim_helpers"
	"github.com/dydxprotocol/v4/x/perpetuals/types"
	pricestypes "github.com/dydxprotocol/v4/x/prices/types"
)

// genNumPerpetuals returns randomized number of perpetuals.
func genNumPerpetuals(r *rand.Rand, isReasonableGenesis bool, numMarkets int) int {
	return simtypes.RandIntBetween(
		r,
		sim_helpers.PickGenesisParameter(sim_helpers.MinNumPerpetuals, isReasonableGenesis),
		sim_helpers.PickGenesisParameter(sim_helpers.MaxNumPerpetuals, isReasonableGenesis)+1,
	)
}

// genNumLiquidityTiers returns a randomized number of liquidity tiers.
func genNumLiquidityTiers(r *rand.Rand, isReasonableGenesis bool) int {
	return simtypes.RandIntBetween(
		r,
		sim_helpers.PickGenesisParameter(sim_helpers.MinNumLiquidityTiers, isReasonableGenesis),
		sim_helpers.PickGenesisParameter(sim_helpers.MaxNumLiquidityTiers, isReasonableGenesis),
	)
}

// genPerpetualToMarketMap returns a list of `Market.Id` that should correspond to each`Perpetual.Id`.
func genPerpetualToMarketMap(r *rand.Rand, numPerpetuals, numMarkets int) []uint32 {
	markets := sim_helpers.MakeRange(uint32(numMarkets))

	// Pad more markets if there are more perpetuals.
	if numPerpetuals > numMarkets {
		diff := numPerpetuals - numMarkets
		extraMarkets := make([]uint32, diff)
		for i := 0; i < diff; i++ {
			randomIdx := simtypes.RandIntBetween(r, 0, numMarkets)
			extraMarkets[i] = markets[randomIdx]
		}
		markets = append(markets, extraMarkets...)
	}

	// Shuffle markets, so we randomize which `Perpetual` gets matched with which `Market`.
	r.Shuffle(numPerpetuals, func(i, j int) { markets[i], markets[j] = markets[j], markets[i] })

	return markets
}

// genTicker returns a randomized string used for `Perpetual.Ticker`.
func genTicker(r *rand.Rand) string {
	return simtypes.RandStringOfLength(r, simtypes.RandIntBetween(r, 3, 6)) + "-USD"
}

// genAtomicResolution returns a randomized int used for `Perpetual.AtomicResolution`.
func genAtomicResolution(r *rand.Rand, isReasonableGenesis bool) int32 {
	return int32(simtypes.RandIntBetween(
		r,
		sim_helpers.PickGenesisParameter(sim_helpers.MinAtomicResolution, isReasonableGenesis),
		sim_helpers.PickGenesisParameter(sim_helpers.MaxAtomicResolution, isReasonableGenesis)+1,
	))
}

// genDefaultFundingPpm returns a randomized int used for `Perpetual.DefaultFundingPpm`.
func genDefaultFundingPpm(r *rand.Rand) int32 {
	defaultFundingPpmAbs := sim_helpers.GetRandomBucketValue(r, sim_helpers.DefaultFundingPpmAbsBuckets)
	if sim_helpers.RandBool(r) {
		return -int32(defaultFundingPpmAbs)
	}
	return int32(defaultFundingPpmAbs)
}

// genInitialAndMaintenanceMargin returns a randomized set of ints used for Initial and Maintenance margins.
func genInitialAndMaintenanceFraction(r *rand.Rand) (uint32, uint32) {
	initialMargin := sim_helpers.GetRandomBucketValue(r, sim_helpers.InitialMarginBuckets)
	// MaintenanceFraction must be less than or equal to 100%.
	maintenanceFraction := simtypes.RandIntBetween(r, 0, 1_000_000)
	return uint32(initialMargin), uint32(maintenanceFraction)
}

// genBasePositionNotional returns a randomized notional value of base position.
func genBasePositionNotional(r *rand.Rand, isReasonableGenesis bool) uint64 {
	return uint64(
		simtypes.RandIntBetween(
			r,
			sim_helpers.PickGenesisParameter(sim_helpers.MinBasePositionNotional, isReasonableGenesis),
			sim_helpers.PickGenesisParameter(sim_helpers.MaxBasePositionNotional, isReasonableGenesis)+1,
		),
	)
}

func genParams(r *rand.Rand, isReasonableGenesis bool) types.Params {
	return types.Params{
		FundingRateClampFactorPpm: genFundingRateClampFactorPpm(r, isReasonableGenesis),
		PremiumVoteClampFactorPpm: genPremiumVoteClampFactorPpm(r, isReasonableGenesis),
	}
}

func genFundingRateClampFactorPpm(r *rand.Rand, isReasonableGenesis bool) uint32 {
	return uint32(
		simtypes.RandIntBetween(
			r,
			sim_helpers.PickGenesisParameter(sim_helpers.MinFundingRateClampFactorPpm, isReasonableGenesis),
			sim_helpers.PickGenesisParameter(sim_helpers.MaxFundingRateClampFactorPpm, isReasonableGenesis)+1,
		),
	)
}

// genPremiumVoteClampFactorPpm returns a randomized uint32 for premium vote clamp factor ppm.
func genPremiumVoteClampFactorPpm(r *rand.Rand, isReasonableGenesis bool) uint32 {
	return uint32(
		simtypes.RandIntBetween(
			r,
			sim_helpers.PickGenesisParameter(sim_helpers.MinPremiumVoteClampFactorPpm, isReasonableGenesis),
			sim_helpers.PickGenesisParameter(sim_helpers.MaxPremiumVoteClampFactorPpm, isReasonableGenesis)+1,
		),
	)
}

// RandomizedGenState generates a random GenesisState for `Perpetuals`.
func RandomizedGenState(simState *module.SimulationState) {
	r := simState.Rand
	isReasonableGenesis := sim_helpers.ShouldGenerateReasonableGenesis(r, simState.GenTimestamp)

	// Generate `Params`.
	params := genParams(r, isReasonableGenesis)

	// Generate `LiquidityTier`s.
	numLiquidityTiers := genNumLiquidityTiers(r, isReasonableGenesis)
	liquidityTiers := make([]types.LiquidityTier, numLiquidityTiers)
	for i := 0; i < numLiquidityTiers; i++ {
		initialMarginPpm, maintenanceFractionPpm := genInitialAndMaintenanceFraction(r)
		basePositionNotional := genBasePositionNotional(r, isReasonableGenesis)
		liquidityTiers[i] = types.LiquidityTier{
			Name:                   fmt.Sprintf("%d", i),
			InitialMarginPpm:       initialMarginPpm,
			MaintenanceFractionPpm: maintenanceFractionPpm,
			BasePositionNotional:   basePositionNotional,
		}
	}

	// Get number of `Prices.Markets`.
	cdc := codec.NewProtoCodec(codectypes.NewInterfaceRegistry())
	pricesGenesisBytes := simState.GenState[pricestypes.ModuleName]
	var pricesGenesis pricestypes.GenesisState
	if err := cdc.UnmarshalJSON(pricesGenesisBytes, &pricesGenesis); err != nil {
		panic(fmt.Sprintf("Could not unmarshal Prices GenesisState %s", err))
	}
	numMarkets := len(pricesGenesis.GetMarkets())
	if numMarkets == 0 {
		panic("Number of Markets cannot be zero")
	}

	// Generate number of `Perpetuals`.
	numPerpetuals := genNumPerpetuals(r, isReasonableGenesis, numMarkets)

	// Generate `Market`s for each `Perpetual`.
	marketsForPerp := genPerpetualToMarketMap(r, numPerpetuals, numMarkets)

	// Generate `Perpetuals`.
	perpetuals := make([]types.Perpetual, numPerpetuals)
	for i := 0; i < numPerpetuals; i++ {
		marketId := marketsForPerp[i]

		perpetuals[i] = types.Perpetual{
			Ticker:            genTicker(r),
			MarketId:          marketId,
			AtomicResolution:  genAtomicResolution(r, isReasonableGenesis),
			DefaultFundingPpm: genDefaultFundingPpm(r),
			FundingIndex:      dtypes.ZeroInt(),
			OpenInterest:      types.DefaultOpenInterest,
			LiquidityTier:     uint32(simtypes.RandIntBetween(r, 0, numLiquidityTiers)),
		}
	}

	perpetualsGenesis := types.GenesisState{
		Perpetuals:     perpetuals,
		LiquidityTiers: liquidityTiers,
		Params:         params,
	}

	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&perpetualsGenesis)
}
