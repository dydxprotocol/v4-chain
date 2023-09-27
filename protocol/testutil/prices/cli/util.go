package cli

import (
	"fmt"
	"github.com/dydxprotocol/v4-chain/protocol/app/stoppable"
	pricefeedtypes "github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/types"
	"sort"
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
func GetTickersSortedByMarketId(marketToMarketConfig map[uint32]pricefeedtypes.MarketConfig) []string {
	// Get all `marketId`s in `marketIdToTicker` as a sorted array.
	marketIds := make([]uint32, 0, len(marketToMarketConfig))
	for marketId := range marketToMarketConfig {
		marketIds = append(marketIds, marketId)
	}
	sort.Slice(marketIds, func(i, j int) bool {
		return marketIds[i] < marketIds[j]
	})

	// Get a list of tickers sorted by their corresponding `marketId`.
	tickers := make([]string, 0, len(marketToMarketConfig))
	for _, marketId := range marketIds {
		tickers = append(tickers, marketToMarketConfig[marketId].Ticker)
	}

	return tickers
}
