package cli

import (
	"fmt"
	"github.com/dydxprotocol/v4-chain/protocol/app/stoppable"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/network"
	"github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
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
			Id:                 uint32(i),
			Pair:               fmt.Sprint(constants.BtcUsdPair, i),
			MinExchanges:       uint32(1),
			MinPriceChangePpm:  uint32((i + 1) * 2),
			ExchangeConfigJson: "{}",
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

	t.Cleanup(func() {
		stoppable.StopServices(t, cfg.GRPCAddress)
	})

	return network.New(t, cfg), state.MarketParams, state.MarketPrices
}
