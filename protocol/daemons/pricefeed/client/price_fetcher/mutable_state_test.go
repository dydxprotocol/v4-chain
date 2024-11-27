package price_fetcher

import (
	"testing"

	"cosmossdk.io/log"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/types"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/stretchr/testify/require"
)

func TestGetNextNMarkets(t *testing.T) {
	pf, err := NewPriceFetcher(
		constants.Exchange1_2MaxQueries_QueryConfig,
		types.ExchangeQueryDetails{},
		&constants.Exchange1_3Markets_MutableExchangeMarketConfig,
		constants.MutableMarketConfigs_3Markets,
		&mocks.ExchangeQueryHandler{},
		log.NewNopLogger(),
		newTestPriceFetcherBufferedChannel(),
	)
	require.NoError(t, err)

	markets := pf.mutableState.GetNextNMarkets(2)
	require.Equal(t, []types.MarketId{7, 8}, markets)
	markets = pf.mutableState.GetNextNMarkets(2)
	require.Equal(t, []types.MarketId{9, 7}, markets)

	// Expect the update to reset the index for the next n markets.
	err = pf.UpdateMutableExchangeConfig(
		&constants.Exchange1_5Markets_MutableExchangeMarketConfig,
		constants.MutableMarketConfigs_5Markets,
	)
	require.NoError(t, err)

	markets = pf.mutableState.GetNextNMarkets(2)
	require.Equal(t, []types.MarketId{7, 8}, markets)
	markets = pf.mutableState.GetNextNMarkets(2)
	require.Equal(t, []types.MarketId{9, 10}, markets)
}

func TestGetMarketIds(t *testing.T) {
	pf, err := NewPriceFetcher(
		constants.Exchange1_2MaxQueries_QueryConfig,
		types.ExchangeQueryDetails{},
		&constants.Exchange1_3Markets_MutableExchangeMarketConfig,
		constants.MutableMarketConfigs_3Markets,
		&mocks.ExchangeQueryHandler{},
		log.NewNopLogger(),
		newTestPriceFetcherBufferedChannel(),
	)
	require.NoError(t, err)

	markets := pf.mutableState.GetMarketIds()
	require.Equal(t, []types.MarketId{7, 8, 9}, markets)
}

func TestGetMarketExponents(t *testing.T) {
	pf, err := NewPriceFetcher(
		constants.Exchange1_2MaxQueries_QueryConfig,
		types.ExchangeQueryDetails{},
		&constants.Exchange1_3Markets_MutableExchangeMarketConfig,
		constants.MutableMarketConfigs_3Markets,
		&mocks.ExchangeQueryHandler{},
		log.NewNopLogger(),
		newTestPriceFetcherBufferedChannel(),
	)
	require.NoError(t, err)

	marketExponents := pf.mutableState.GetMarketExponents()
	// Check that the mutableState contains the correct set of marketExponents
	// and that it returns a copy of the map and not the original.
	require.NotSame(t, &marketExponents, &pf.mutableState.marketExponents)
	require.Equal(t, pf.mutableState.marketExponents, marketExponents)
}

func TestGetMutableExchangeConfig(t *testing.T) {
	pf, err := NewPriceFetcher(
		constants.Exchange1_2MaxQueries_QueryConfig,
		types.ExchangeQueryDetails{},
		&constants.Exchange1_3Markets_MutableExchangeMarketConfig,
		constants.MutableMarketConfigs_3Markets,
		&mocks.ExchangeQueryHandler{},
		log.NewNopLogger(),
		newTestPriceFetcherBufferedChannel(),
	)
	require.NoError(t, err)

	mutableExchangeConfig := pf.mutableState.GetMutableExchangeConfig()

	require.NotSame(t, mutableExchangeConfig, pf.mutableState.mutableExchangeConfig)
	require.Equal(t, pf.mutableState.mutableExchangeConfig, mutableExchangeConfig)
}

// TestGetTaskLoopDefinition asserts that the task loop definition is correctly
// set from mutable state and that it uses copies of all identical data structures.
func TestGetTaskLoopDefinition(t *testing.T) {
	pf, err := NewPriceFetcher(
		constants.Exchange1_2MaxQueries_QueryConfig,
		types.ExchangeQueryDetails{},
		&constants.Exchange1_3Markets_MutableExchangeMarketConfig,
		constants.MutableMarketConfigs_3Markets,
		&mocks.ExchangeQueryHandler{},
		log.NewNopLogger(),
		newTestPriceFetcherBufferedChannel(),
	)
	require.NoError(t, err)

	taskLoopDefinition := pf.getTaskLoopDefinition()

	// The taskLoopDefinition should use copies of shared state
	require.NotSame(t, pf.mutableState.mutableExchangeConfig, taskLoopDefinition.mutableExchangeConfig)
	require.NotSame(t, &pf.mutableState.marketExponents, &taskLoopDefinition.marketExponents)

	require.Equal(t, pf.mutableState.mutableExchangeConfig, taskLoopDefinition.mutableExchangeConfig)
	require.Equal(t, pf.mutableState.marketExponents, taskLoopDefinition.marketExponents)
	require.Equal(t, []types.MarketId{7, 8}, taskLoopDefinition.marketIds)
}
