package simulation_test

import (
	"encoding/json"
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/dydxprotocol/v4/lib"
	testutil_rand "github.com/dydxprotocol/v4/testutil/rand"
	"github.com/dydxprotocol/v4/x/perpetuals/simulation"
	"github.com/dydxprotocol/v4/x/perpetuals/types"
	pricessimulation "github.com/dydxprotocol/v4/x/prices/simulation"
	"github.com/stretchr/testify/require"
)

func TestRandomizedGenState(t *testing.T) {
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(interfaceRegistry)

	r := testutil_rand.NewRand()

	simState := module.SimulationState{
		AppParams:    make(simtypes.AppParams),
		Cdc:          cdc,
		Rand:         r,
		NumBonded:    3,
		Accounts:     simtypes.RandomAccounts(r, 3),
		InitialStake: sdkmath.NewInt(1000),
		GenState:     make(map[string]json.RawMessage),
	}

	for i := 0; i < 100; i++ {
		// `Perpetuals` module has a dependency on `Prices` module.
		pricessimulation.RandomizedGenState(&simState)

		simulation.RandomizedGenState(&simState)
		var perpetualsGenesis types.GenesisState
		simState.Cdc.MustUnmarshalJSON(simState.GenState[types.ModuleName], &perpetualsGenesis)

		require.True(t, len(perpetualsGenesis.Perpetuals) >= 1)
		require.True(t, len(perpetualsGenesis.Perpetuals) <= 2_000)

		require.True(t, len(perpetualsGenesis.LiquidityTiers) >= 1)

		require.True(t, perpetualsGenesis.Params.FundingRateClampFactorPpm > 0)
		require.True(t, perpetualsGenesis.Params.PremiumVoteClampFactorPpm > 0)

		for _, lt := range perpetualsGenesis.LiquidityTiers {
			require.True(t, len(lt.Name) >= 1)

			require.True(t, lt.InitialMarginPpm <= lib.OneMillion)

			require.True(t, lt.MaintenanceFractionPpm <= lib.OneMillion)
		}

		for _, perp := range perpetualsGenesis.Perpetuals {
			require.True(t, len(perp.Ticker) >= 1)

			require.True(t, perp.MarketId <= 1_000)

			require.True(t, perp.AtomicResolution >= -10)
			require.True(t, perp.AtomicResolution <= 10)

			require.True(t, perp.DefaultFundingPpm > -int32(lib.OneMillion))
			require.True(t, perp.DefaultFundingPpm < int32(lib.OneMillion))

			require.True(t, perp.FundingIndex.BigInt().Sign() == 0)

			require.True(t, perp.OpenInterest == 0)
		}
	}
}
