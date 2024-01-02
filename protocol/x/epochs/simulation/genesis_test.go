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
	"github.com/dydxprotocol/v4-chain/protocol/x/epochs/simulation"
	"github.com/dydxprotocol/v4-chain/protocol/x/epochs/types"
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
		GenTimestamp: simtypes.RandTimestamp(r),
		BondDenom:    sdk.DefaultBondDenom,
	}

	for i := 0; i < 100; i++ {
		simulation.RandomizedGenState(&simState)
		var epochsGenesis types.GenesisState
		simState.Cdc.MustUnmarshalJSON(simState.GenState[types.ModuleName], &epochsGenesis)

		require.True(t, len(epochsGenesis.EpochInfoList) >= 1+3)     // +3 comes from default genesis
		require.True(t, len(epochsGenesis.EpochInfoList) <= 1_000+3) // +3 comes from default genesis

		for _, epochInfo := range epochsGenesis.EpochInfoList {
			require.True(t, len(epochInfo.Name) >= 5)
			require.True(t, len(epochInfo.Name) <= 19)

			require.True(t, epochInfo.Duration >= 1)
			// no need to check duration max value since it's already capped by uint32 type.

			// no need to check `NextTick`, `CurrentEpcoh` or `CurrentEpochStartBlock` value
			// since we randomly generate any uint32 value.

			// no need to check `FastForwardNextTick` value since we generate both true/false.
		}
	}
}
