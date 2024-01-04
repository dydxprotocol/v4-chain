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
	testutil_rand "github.com/dydxprotocol/v4-chain/protocol/testutil/rand"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/sim_helpers"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/simulation"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	perpetualssimulation "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/simulation"
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
		// `Perpetuals` genesis has a dependency on `Prices` genesis.
		pricessimulation.RandomizedGenState(&simState)

		// `CLOB` genesis has a dependency on `Perpetuals` genesis.
		perpetualssimulation.RandomizedGenState(&simState)

		simulation.RandomizedGenState(&simState)
		var clobGenesis types.GenesisState
		simState.Cdc.MustUnmarshalJSON(simState.GenState[types.ModuleName], &clobGenesis)

		require.True(t, len(clobGenesis.ClobPairs) >= sim_helpers.MinValidClobPairs)
		require.True(t, len(clobGenesis.ClobPairs) <= sim_helpers.MaxValidClobPairs)

		for _, clobPair := range clobGenesis.ClobPairs {
			// Note that we only validate the `MaxQuantumConversionExponent` field because all other
			// fields are either validated by `CreatePerpetualClobPair` or all values of that type are valid.
			require.GreaterOrEqual(t, int32(sim_helpers.MaxQuantumConversionExponent.Valid), clobPair.QuantumConversionExponent)
			require.LessOrEqual(t, int32(sim_helpers.MinQuantumConversionExponent.Valid), clobPair.QuantumConversionExponent)
		}

		liquidationConfig := clobGenesis.LiquidationsConfig

		// Validate minimum position notional is within the specified range since we don't do any
		// validation in LiquidationConfig.Validate()
		require.LessOrEqual(
			t,
			uint64(sim_helpers.MinPositionNotionalBuckets[0]),
			liquidationConfig.PositionBlockLimits.MinPositionNotionalLiquidated,
		)
		require.GreaterOrEqual(
			t,
			uint64(sim_helpers.MinPositionNotionalBuckets[len(sim_helpers.MinPositionNotionalBuckets)-1]),
			liquidationConfig.PositionBlockLimits.MinPositionNotionalLiquidated,
		)
		err := liquidationConfig.Validate()
		require.NoError(t, err)
	}
}
