package cli

import (
	"fmt"
	"testing"

	"github.com/dydxprotocol/v4/testutil/constants"
	"github.com/dydxprotocol/v4/testutil/network"
	"github.com/dydxprotocol/v4/x/prices/types"
	"github.com/stretchr/testify/require"
)

func NetworkWithMarketObjects(t *testing.T, n int) (*network.Network, []types.MarketParam, []types.MarketPrice) {
	t.Helper()
	cfg := network.DefaultConfig(nil)
	state := types.GenesisState{}
	require.NoError(t, cfg.Codec.UnmarshalJSON(cfg.GenesisState[types.ModuleName], &state))

	// Overwrite market params and prices in default genesis state.
	state.MarketParams = []types.MarketParam{}
	state.MarketPrices = []types.MarketPrice{}

	// Market params
	for i := 0; i < n; i++ {
		marketParam := types.MarketParam{
			Id:                uint32(i),
			Pair:              fmt.Sprint(constants.BtcUsdPair, i),
			MinExchanges:      uint32(1),
			MinPriceChangePpm: uint32((i + 1) * 2),
		}
		state.MarketParams = append(state.MarketParams, marketParam)
	}

	// Market prices
	for i := 0; i < n; i++ {
		marketPrice := types.MarketPrice{
			Id:    uint32(i),
			Price: constants.FiveBillion,
		}
		state.MarketPrices = append(state.MarketPrices, marketPrice)
	}

	buf, err := cfg.Codec.MarshalJSON(&state)
	require.NoError(t, err)
	cfg.GenesisState[types.ModuleName] = buf
	return network.New(t, cfg), state.MarketParams, state.MarketPrices
}
