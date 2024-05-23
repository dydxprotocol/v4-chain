package simulation_test

import (
	"encoding/json"
	sdk "github.com/cosmos/cosmos-sdk/types"
	v4module "github.com/dydxprotocol/v4-chain/protocol/app/module"
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	testutil_rand "github.com/dydxprotocol/v4-chain/protocol/testutil/rand"
	"github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/simulation"
	"github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	pricessimulation "github.com/dydxprotocol/v4-chain/protocol/x/prices/simulation"
	"github.com/stretchr/testify/require"
)

func TestRandomizedGenState(t *testing.T) {
	cdc := codec.NewProtoCodec(v4module.InterfaceRegistry)

	r := testutil_rand.NewRand()

	simState := module.SimulationState{
		AppParams:    make(simtypes.AppParams),
		Cdc:          cdc,
		Rand:         r,
		NumBonded:    3,
		Accounts:     simtypes.RandomAccounts(r, 3),
		InitialStake: sdkmath.NewInt(1000),
		GenState:     make(map[string]json.RawMessage),
		BondDenom:    sdk.DefaultBondDenom,
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
			require.True(t, len(perp.Params.Ticker) >= 1)

			require.True(t, perp.Params.MarketId <= 1_000)

			require.True(t, perp.Params.AtomicResolution >= -10)
			require.True(t, perp.Params.AtomicResolution <= 10)

			require.True(t, perp.Params.DefaultFundingPpm > -int32(lib.OneMillion))
			require.True(t, perp.Params.DefaultFundingPpm < int32(lib.OneMillion))

			require.True(t, perp.FundingIndex.Sign() == 0)
		}
	}
}
