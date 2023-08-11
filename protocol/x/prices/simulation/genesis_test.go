package simulation_test

import (
	"encoding/json"
	"strings"
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	testutil_rand "github.com/dydxprotocol/v4/testutil/rand"
	"github.com/dydxprotocol/v4/x/prices/simulation"
	"github.com/dydxprotocol/v4/x/prices/types"
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
		simulation.RandomizedGenState(&simState)
		var pricesGenesis types.GenesisState
		simState.Cdc.MustUnmarshalJSON(simState.GenState[types.ModuleName], &pricesGenesis)

		require.True(t, len(pricesGenesis.ExchangeFeeds) >= 1)
		require.True(t, len(pricesGenesis.ExchangeFeeds) <= 1_000)

		require.True(t, len(pricesGenesis.Markets) >= 1)
		require.True(t, len(pricesGenesis.Markets) <= 1_000)
		for _, market := range pricesGenesis.Markets {
			require.True(t, len(market.Pair) >= 7)
			require.True(t, strings.HasSuffix(market.Pair, "-USD"))

			require.True(t, market.Exponent >= -15)
			require.True(t, market.Exponent <= 15)

			require.True(t, market.MinPriceChangePpm >= 1)
			require.True(t, market.MinPriceChangePpm < 10_000)
		}
	}
}
