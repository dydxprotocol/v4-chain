package price_fetcher

import (
	"errors"
	"testing"

	daemontypes "github.com/dydxprotocol/v4-chain/protocol/daemons/types"

	"cosmossdk.io/math"
	pricefeed_cosntants "github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/constants"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/daemons/pricefeed"

	"cosmossdk.io/log"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const (
	taskLoopIterations = 2
)

var (
	exchangeQueryHandlerFailure = errors.New("Failed to query exchange")
	tickerNotAvailable          = errors.New("Ticker not listed")
)

// TestRunTaskLoop tests that different exchange configurations results in the expected queries being made, and prices
// produced.
func TestRunTaskLoop(t *testing.T) {
	tests := map[string]struct {
		// parameters
		queryConfig           types.ExchangeQueryConfig
		queryDetails          types.ExchangeQueryDetails
		mutableExchangeConfig types.MutableExchangeMarketConfig
		mutableMarketConfigs  []*types.MutableMarketConfig

		// expectations
		expectedMarketIdsCalled []types.MarketId
	}{
		"No markets": {
			queryConfig:           constants.Exchange1_0MaxQueries_QueryConfig,
			queryDetails:          constants.SingleMarketExchangeQueryDetails,
			mutableExchangeConfig: constants.Exchange1_NoMarkets_MutableExchangeMarketConfig,
			mutableMarketConfigs:  constants.MutableMarketConfigs_0Markets,
		},
		"Num markets equals max query markets where there is only 1 market": {
			queryConfig:           constants.Exchange1_1MaxQueries_QueryConfig,
			queryDetails:          constants.SingleMarketExchangeQueryDetails,
			mutableExchangeConfig: constants.Exchange1_1Markets_MutableExchangeMarketConfig,
			mutableMarketConfigs:  constants.MutableMarketConfigs_1Markets,
			expectedMarketIdsCalled: []types.MarketId{
				constants.MarketId7,
				constants.MarketId7,
			},
		},
		"Num markets < max query markets": {
			queryConfig:           constants.Exchange1_2MaxQueries_QueryConfig,
			queryDetails:          constants.SingleMarketExchangeQueryDetails,
			mutableExchangeConfig: constants.Exchange1_1Markets_MutableExchangeMarketConfig,
			mutableMarketConfigs:  constants.MutableMarketConfigs_1Markets,
			expectedMarketIdsCalled: []types.MarketId{
				constants.MarketId7,
				constants.MarketId7,
			},
		},
		"Num markets equals max query markets": {
			queryConfig:           constants.Exchange1_2MaxQueries_QueryConfig,
			queryDetails:          constants.SingleMarketExchangeQueryDetails,
			mutableExchangeConfig: constants.Exchange1_2Markets_MutableExchangeMarketConfig,
			mutableMarketConfigs:  constants.MutableMarketConfigs_2Markets,
			expectedMarketIdsCalled: []types.MarketId{
				constants.MarketId7,
				constants.MarketId8,
				constants.MarketId7,
				constants.MarketId8,
			},
		},
		"Multi-market, 2 markets": {
			queryConfig:           constants.Exchange1_1MaxQueries_QueryConfig,
			queryDetails:          constants.MultiMarketExchangeQueryDetails,
			mutableExchangeConfig: constants.Exchange1_2Markets_MutableExchangeMarketConfig,
			mutableMarketConfigs:  constants.MutableMarketConfigs_2Markets,
			expectedMarketIdsCalled: []types.MarketId{
				constants.MarketId7,
				constants.MarketId8,
				constants.MarketId7,
				constants.MarketId8,
			},
		},
		"Num markets greater than max query markets": {
			queryConfig:           constants.Exchange1_2MaxQueries_QueryConfig,
			queryDetails:          constants.SingleMarketExchangeQueryDetails,
			mutableExchangeConfig: constants.Exchange1_3Markets_MutableExchangeMarketConfig,
			mutableMarketConfigs:  constants.MutableMarketConfigs_3Markets,
			expectedMarketIdsCalled: []types.MarketId{
				constants.MarketId7,
				constants.MarketId8,
				constants.MarketId9,
				constants.MarketId7,
			},
		},
		"Multi-market, 5 markets": {
			queryConfig:           constants.Exchange1_1MaxQueries_QueryConfig,
			queryDetails:          constants.MultiMarketExchangeQueryDetails,
			mutableExchangeConfig: constants.Exchange1_5Markets_MutableExchangeMarketConfig,
			mutableMarketConfigs:  constants.MutableMarketConfigs_5Markets,
			expectedMarketIdsCalled: []types.MarketId{
				constants.MarketId7,
				constants.MarketId8,
				constants.MarketId9,
				constants.MarketId10,
				constants.MarketId11,
				constants.MarketId7,
				constants.MarketId8,
				constants.MarketId9,
				constants.MarketId10,
				constants.MarketId11,
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup for testing.
			bCh := newTestPriceFetcherBufferedChannel()

			queryHandler := generateMockExchangeQueryHandler()

			pf, err := NewPriceFetcher(
				tc.queryConfig,
				tc.queryDetails,
				&tc.mutableExchangeConfig,
				tc.mutableMarketConfigs,
				queryHandler,
				log.NewNopLogger(),
				bCh,
			)
			require.NoError(t, err)

			// Run sub-task a specified number of iterations.
			for i := 0; i < taskLoopIterations; i++ {
				pf.RunTaskLoop(&daemontypes.RequestHandlerImpl{})
			}

			// Will hang until tests timeout if bCh is not full.
			var bufferedResponses []*types.MarketPriceTimestamp
			for len(bufferedResponses) < len(tc.expectedMarketIdsCalled) {
				// Will block until test times out if bCh is not written to `tc.expectedMarketIdsCalled`
				// times.
				val := <-bCh
				bufferedResponses = append(bufferedResponses, val.Price)
				require.NoError(t, val.Err)
			}
			close(bCh)

			expectedBChValues := make([]*types.MarketPriceTimestamp, 0, len(tc.expectedMarketIdsCalled))
			for _, market := range tc.expectedMarketIdsCalled {
				expectedBChValues = append(expectedBChValues, constants.CanonicalMarketPriceTimestampResponses[market])
			}

			// Verify contents of buffered channel.
			require.ElementsMatch(
				t,
				expectedBChValues,
				bufferedResponses,
			)

			// Verify each go routine was called as expected.
			// NOTE: ordering of calls is not checked in `AssertCalled`.
			expectedQueries := len(tc.expectedMarketIdsCalled)
			marketsPerCall := 1
			if tc.queryDetails.IsMultiMarket {
				expectedQueries = taskLoopIterations
				marketsPerCall = len(tc.mutableExchangeConfig.MarketToMarketConfig)
			}
			queryHandler.AssertNumberOfCalls(t, "Query", expectedQueries)
			for i := 0; i < len(tc.expectedMarketIdsCalled); i = i + marketsPerCall {
				assertQueryHandlerCalledWithMarkets(
					t,
					queryHandler,
					tc.expectedMarketIdsCalled[i:i+marketsPerCall],
					tc.mutableMarketConfigs,
				)
			}
		})
	}
}

// TestGetTaskLoopDefinition_SingleMarketExchange tests that the `getTaskLoopDefinition` method returns the
// expected task loop definition for a single market exchange, and that it properly advances through the
// market ids for an exchange on every call.
func TestGetTaskLoopDefinition_SingleMarketExchange(t *testing.T) {
	pf, err := NewPriceFetcher(
		constants.Exchange1_2MaxQueries_QueryConfig,
		constants.SingleMarketExchangeQueryDetails,
		&constants.Exchange1_3Markets_MutableExchangeMarketConfig,
		constants.MutableMarketConfigs_3Markets,
		&mocks.ExchangeQueryHandler{},
		log.NewNopLogger(),
		newTestPriceFetcherBufferedChannel(),
	)
	require.NoError(t, err)

	expectedExchangeConfig := &constants.Exchange1_3Markets_MutableExchangeMarketConfig
	expectedMarketExponents := generateMarketExponentsMap(constants.MutableMarketConfigs_3Markets)

	taskLoopDefinition := pf.getTaskLoopDefinition()

	// Expect that the definition uses a copy of the mutableExchangeConfig for synchronization purposes.
	require.NotSame(t, pf.mutableState.mutableExchangeConfig, taskLoopDefinition.mutableExchangeConfig)
	require.Equal(t, expectedExchangeConfig, taskLoopDefinition.mutableExchangeConfig)
	require.Equal(t, expectedMarketExponents, taskLoopDefinition.marketExponents)
	require.Equal(t, []types.MarketId{constants.MarketId7, constants.MarketId8}, taskLoopDefinition.marketIds)

	// Expect that the market ids ring has been advanced by 2.
	taskLoopDefinition = pf.getTaskLoopDefinition()

	// Sanity checks:

	require.NotSame(t, pf.mutableState.mutableExchangeConfig, taskLoopDefinition.mutableExchangeConfig)
	require.Equal(t, expectedExchangeConfig, taskLoopDefinition.mutableExchangeConfig)
	require.Equal(t, expectedMarketExponents, taskLoopDefinition.marketExponents)

	// Test that the markets have changed as expected.
	require.Equal(t, []types.MarketId{constants.MarketId9, constants.MarketId7}, taskLoopDefinition.marketIds)
}

// TestGetTaskLoopDefinition_MultiMarketExchange tests that the `getTaskLoopDefinition` method returns the
// expected task loop definition for a multi-market exchange.
func TestGetTaskLoopDefinition_MultiMarketExchange(t *testing.T) {
	pf, err := NewPriceFetcher(
		constants.Exchange1_1MaxQueries_QueryConfig,
		constants.MultiMarketExchangeQueryDetails,
		&constants.Exchange1_3Markets_MutableExchangeMarketConfig,
		constants.MutableMarketConfigs_3Markets,
		&mocks.ExchangeQueryHandler{},
		log.NewNopLogger(),
		newTestPriceFetcherBufferedChannel(),
	)
	require.NoError(t, err)

	taskLoopDefinition := pf.getTaskLoopDefinition()
	expectedMarkets := constants.Exchange1_3Markets_MutableExchangeMarketConfig.GetMarketIds()
	expectedExponents := generateMarketExponentsMap(constants.MutableMarketConfigs_3Markets)
	expectedExchangeConfig := &constants.Exchange1_3Markets_MutableExchangeMarketConfig

	// Expect that the definition uses a copy of the mutableExchangeConfig for synchronization purposes.
	require.NotSame(t, expectedExchangeConfig, taskLoopDefinition.mutableExchangeConfig)
	require.Equal(t, expectedExchangeConfig, taskLoopDefinition.mutableExchangeConfig)
	require.Equal(t, expectedExponents, taskLoopDefinition.marketExponents)
	require.Equal(t, expectedMarkets, taskLoopDefinition.marketIds)
}

// TestUpdateMutableExchangeConfig_CorrectlyUpdatesTaskDefinition tests that the `updateMutableExchangeConfig` method
// changes the price fetcher's state correctly so that the next call to `getTaskLoopDefinition` returns the expected
// definition.
func TestUpdateMutableExchangeConfig_CorrectlyUpdatesTaskDefinition(t *testing.T) {
	tests := map[string]struct {
		// parameters
		queryConfig  types.ExchangeQueryConfig
		queryDetails types.ExchangeQueryDetails

		initialMutableExchangeConfig types.MutableExchangeMarketConfig
		initialMarketConfig          []*types.MutableMarketConfig

		updateMutableExchangeConfig types.MutableExchangeMarketConfig
		updateMarketConfig          []*types.MutableMarketConfig

		isMultiMarket bool

		// expectations
		initialExpectedExponents map[types.MarketId]types.Exponent
		updateExpectedExponents  map[types.MarketId]types.Exponent
	}{
		"Multimarket: No markets to markets": {
			queryConfig:   constants.Exchange1_1MaxQueries_QueryConfig,
			queryDetails:  constants.MultiMarketExchangeQueryDetails,
			isMultiMarket: true,

			initialMutableExchangeConfig: constants.Exchange1_NoMarkets_MutableExchangeMarketConfig,
			initialMarketConfig:          constants.MutableMarketConfigs_0Markets,
			initialExpectedExponents:     map[types.MarketId]types.Exponent{},

			updateMutableExchangeConfig: constants.Exchange1_3Markets_MutableExchangeMarketConfig,
			updateMarketConfig:          constants.MutableMarketConfigs_3Markets,
			updateExpectedExponents:     constants.MutableMarketConfigs_3Markets_ExpectedExponents,
		},
		"Multimarket: Add markets": {
			queryConfig:   constants.Exchange1_1MaxQueries_QueryConfig,
			queryDetails:  constants.MultiMarketExchangeQueryDetails,
			isMultiMarket: true,

			initialMutableExchangeConfig: constants.Exchange1_3Markets_MutableExchangeMarketConfig,
			initialMarketConfig:          constants.MutableMarketConfigs_3Markets,
			initialExpectedExponents:     constants.MutableMarketConfigs_3Markets_ExpectedExponents,

			updateMutableExchangeConfig: constants.Exchange1_5Markets_MutableExchangeMarketConfig,
			updateMarketConfig:          constants.MutableMarketConfigs_5Markets,
			updateExpectedExponents:     constants.MutableMarketConfigs_5Markets_ExpectedExponents,
		},
		"Single market: No markets to markets": {
			queryConfig:  constants.Exchange1_2MaxQueries_QueryConfig,
			queryDetails: constants.SingleMarketExchangeQueryDetails,

			initialMutableExchangeConfig: constants.Exchange1_NoMarkets_MutableExchangeMarketConfig,
			initialMarketConfig:          constants.MutableMarketConfigs_0Markets,
			initialExpectedExponents:     map[types.MarketId]types.Exponent{},

			updateMutableExchangeConfig: constants.Exchange1_3Markets_MutableExchangeMarketConfig,
			updateMarketConfig:          constants.MutableMarketConfigs_3Markets,
			updateExpectedExponents:     constants.MutableMarketConfigs_3Markets_ExpectedExponents,
		},
		"Single market: Add markets": {
			queryConfig:  constants.Exchange1_2MaxQueries_QueryConfig,
			queryDetails: constants.SingleMarketExchangeQueryDetails,

			initialMutableExchangeConfig: constants.Exchange1_3Markets_MutableExchangeMarketConfig,
			initialMarketConfig:          constants.MutableMarketConfigs_3Markets,
			initialExpectedExponents:     constants.MutableMarketConfigs_3Markets_ExpectedExponents,

			updateMutableExchangeConfig: constants.Exchange1_5Markets_MutableExchangeMarketConfig,
			updateMarketConfig:          constants.MutableMarketConfigs_5Markets,
			updateExpectedExponents:     constants.MutableMarketConfigs_5Markets_ExpectedExponents,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup for testing.
			bCh := newTestPriceFetcherBufferedChannel()

			queryHandler := generateMockExchangeQueryHandler()

			pf, err := NewPriceFetcher(
				tc.queryConfig,
				tc.queryDetails,
				&tc.initialMutableExchangeConfig,
				tc.initialMarketConfig,
				queryHandler,
				log.NewNopLogger(),
				bCh,
			)
			require.NoError(t, err)

			taskLoopDefinition := pf.getTaskLoopDefinition()

			require.Equal(t, tc.initialExpectedExponents, taskLoopDefinition.marketExponents)
			require.Equal(t, tc.initialMutableExchangeConfig, *taskLoopDefinition.mutableExchangeConfig)

			if tc.isMultiMarket || len(tc.initialMarketConfig) == 0 {
				require.Equal(t, tc.initialMutableExchangeConfig.GetMarketIds(), taskLoopDefinition.marketIds)
			} else {
				numMarkets := lib.Min(len(tc.initialMarketConfig), int(tc.queryConfig.MaxQueries))
				require.Len(t, taskLoopDefinition.marketIds, numMarkets)

				for i := 0; i < numMarkets; i++ {
					require.Equal(t, tc.initialMarketConfig[i].Id, taskLoopDefinition.marketIds[i])
				}
			}

			err = pf.UpdateMutableExchangeConfig(&tc.updateMutableExchangeConfig, tc.updateMarketConfig)
			require.NoError(t, err)

			taskLoopDefinition = pf.getTaskLoopDefinition()

			require.Equal(t, tc.updateExpectedExponents, taskLoopDefinition.marketExponents)
			require.Equal(t, tc.updateMutableExchangeConfig, *taskLoopDefinition.mutableExchangeConfig)

			if tc.isMultiMarket {
				require.Equal(t, tc.updateMutableExchangeConfig.GetMarketIds(), taskLoopDefinition.marketIds)
			} else {
				numMarkets := lib.Min(len(tc.initialMarketConfig), int(tc.queryConfig.MaxQueries))
				require.Len(t, taskLoopDefinition.marketIds, int(tc.queryConfig.MaxQueries))

				for i := 0; i < numMarkets; i++ {
					require.Equal(t, tc.updateMarketConfig[i].Id, taskLoopDefinition.marketIds[i])
				}
			}
		})
	}
}

// TestUpdateMutableExchangeConfig_ProducesExpectedPrices tests that updating the price fetcher's mutable
// config produces the expected prices before and after the update, with no errors. This test validates
// that the query handler receives all the expected inputs to resolve the expected market prices.
func TestUpdateMutableExchangeConfig_ProducesExpectedPrices(t *testing.T) {
	tests := map[string]struct {
		// parameters
		queryConfig  types.ExchangeQueryConfig
		queryDetails types.ExchangeQueryDetails

		initialMutableExchangeConfig types.MutableExchangeMarketConfig
		initialMarketConfigs         []*types.MutableMarketConfig

		updateMutableExchangeConfig types.MutableExchangeMarketConfig
		updateMarketConfigs         []*types.MutableMarketConfig

		isMultiMarket bool

		// expectations
		expectedMarketIdsCalled []types.MarketId
		expectedNumQueryCalls   int
	}{
		"Multimarket: No markets to markets": {
			queryConfig:   constants.Exchange1_1MaxQueries_QueryConfig,
			queryDetails:  constants.MultiMarketExchangeQueryDetails,
			isMultiMarket: true,

			initialMutableExchangeConfig: constants.Exchange1_NoMarkets_MutableExchangeMarketConfig,
			initialMarketConfigs:         constants.MutableMarketConfigs_0Markets,

			updateMutableExchangeConfig: constants.Exchange1_3Markets_MutableExchangeMarketConfig,
			updateMarketConfigs:         constants.MutableMarketConfigs_3Markets,

			// Expect to loop through both sets of markets twice. The first set is empty.
			expectedMarketIdsCalled: []types.MarketId{
				constants.MarketId7,
				constants.MarketId8,
				constants.MarketId9,
				constants.MarketId7,
				constants.MarketId8,
				constants.MarketId9,
			},
			// We don't expect query calls for the initial empty set of markets.
			expectedNumQueryCalls: 2,
		},
		"Multimarket: Add markets": {
			queryConfig:   constants.Exchange1_1MaxQueries_QueryConfig,
			queryDetails:  constants.MultiMarketExchangeQueryDetails,
			isMultiMarket: true,

			initialMutableExchangeConfig: constants.Exchange1_3Markets_MutableExchangeMarketConfig,
			initialMarketConfigs:         constants.MutableMarketConfigs_3Markets,

			updateMutableExchangeConfig: constants.Exchange1_5Markets_MutableExchangeMarketConfig,
			updateMarketConfigs:         constants.MutableMarketConfigs_5Markets,

			// Expect to loop through both sets of markets twice.
			expectedMarketIdsCalled: []types.MarketId{
				constants.MarketId7,
				constants.MarketId8,
				constants.MarketId9,
				constants.MarketId7,
				constants.MarketId8,
				constants.MarketId9,
				constants.MarketId7,
				constants.MarketId8,
				constants.MarketId9,
				constants.MarketId10,
				constants.MarketId11,
				constants.MarketId7,
				constants.MarketId8,
				constants.MarketId9,
				constants.MarketId10,
				constants.MarketId11,
			},
			// Multi-market exchanges are expected to call 1x per task loop run.
			expectedNumQueryCalls: 4,
		},
		"Single market: No markets to markets": {
			queryConfig:  constants.Exchange1_2MaxQueries_QueryConfig,
			queryDetails: constants.SingleMarketExchangeQueryDetails,

			initialMutableExchangeConfig: constants.Exchange1_NoMarkets_MutableExchangeMarketConfig,
			initialMarketConfigs:         constants.MutableMarketConfigs_0Markets,

			updateMutableExchangeConfig: constants.Exchange1_3Markets_MutableExchangeMarketConfig,
			updateMarketConfigs:         constants.MutableMarketConfigs_3Markets,

			expectedMarketIdsCalled: []types.MarketId{
				constants.MarketId7,
				constants.MarketId8,
				constants.MarketId9,
				constants.MarketId7,
			},
			// Single-market exchanges are expected to call up to max queries per task loop run.
			expectedNumQueryCalls: 4,
		},
		"Single market: Add markets": {
			queryConfig:  constants.Exchange1_2MaxQueries_QueryConfig,
			queryDetails: constants.SingleMarketExchangeQueryDetails,

			initialMutableExchangeConfig: constants.Exchange1_3Markets_MutableExchangeMarketConfig,
			initialMarketConfigs:         constants.MutableMarketConfigs_3Markets,

			updateMutableExchangeConfig: constants.Exchange1_5Markets_MutableExchangeMarketConfig,
			updateMarketConfigs:         constants.MutableMarketConfigs_5Markets,

			expectedMarketIdsCalled: []types.MarketId{
				constants.MarketId7,
				constants.MarketId8,
				constants.MarketId9,
				constants.MarketId7,
				constants.MarketId7,
				constants.MarketId8,
				constants.MarketId9,
				constants.MarketId10,
			},
			// Single-market exchanges are expected to call up to max queries per task loop run.
			expectedNumQueryCalls: 8,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup for testing.
			bCh := newTestPriceFetcherBufferedChannel()
			queryHandler := generateMockExchangeQueryHandler()
			pf, err := NewPriceFetcher(
				tc.queryConfig,
				tc.queryDetails,
				&tc.initialMutableExchangeConfig,
				tc.initialMarketConfigs,
				queryHandler,
				log.NewNopLogger(),
				bCh,
			)
			require.NoError(t, err)

			// Run sub-task a specified number of iterations.
			for i := 0; i < taskLoopIterations; i++ {
				pf.RunTaskLoop(&daemontypes.RequestHandlerImpl{})
			}

			// No race conditions should affect the market output of the previous or following task loops.
			err = pf.UpdateMutableExchangeConfig(&tc.updateMutableExchangeConfig, tc.updateMarketConfigs)
			require.NoError(t, err)

			// Run sub-task a specified number of iterations.
			for i := 0; i < taskLoopIterations; i++ {
				go pf.RunTaskLoop(&daemontypes.RequestHandlerImpl{})
			}

			// Will hang until tests timeout if bCh is not full.
			var bufferedResponses []*types.MarketPriceTimestamp
			for len(bufferedResponses) < len(tc.expectedMarketIdsCalled) {
				// Will block until test times out if bCh is not written to `tc.expectedMarketIdsCalled`
				// times.
				val := <-bCh
				bufferedResponses = append(bufferedResponses, val.Price)
				require.NoError(t, val.Err)
			}
			close(bCh)

			expectedBChValues := make([]*types.MarketPriceTimestamp, 0, len(tc.expectedMarketIdsCalled))
			for _, market := range tc.expectedMarketIdsCalled {
				expectedBChValues = append(expectedBChValues, constants.CanonicalMarketPriceTimestampResponses[market])
			}

			// Verify contents of buffered channel.
			require.ElementsMatch(
				t,
				expectedBChValues,
				bufferedResponses,
			)

			queryHandler.AssertNumberOfCalls(t, "Query", tc.expectedNumQueryCalls)
			if tc.isMultiMarket {
				// For multi market exchanges, the query handler should be called once per task loop, and each
				// query should be for all markets supported by the exchange at that time.
				initialNumMarkets := len(tc.initialMutableExchangeConfig.MarketToMarketConfig)
				updateNumMarkets := len(tc.updateMutableExchangeConfig.MarketToMarketConfig)

				indexPtr := 0
				if initialNumMarkets > 0 {
					for i := 0; i < taskLoopIterations; i++ {
						assertQueryHandlerCalledWithMarkets(
							t,
							queryHandler,
							tc.expectedMarketIdsCalled[indexPtr:indexPtr+initialNumMarkets],
							tc.initialMarketConfigs,
						)
						indexPtr += initialNumMarkets
					}
				}
				for i := 0; i < taskLoopIterations; i++ {
					assertQueryHandlerCalledWithMarkets(
						t,
						queryHandler,
						tc.expectedMarketIdsCalled[indexPtr:indexPtr+updateNumMarkets],
						tc.updateMarketConfigs,
					)
					indexPtr += updateNumMarkets
				}
			} else {
				// For single market exchanges, the query handler should be called once per market per task loop.
				marketsPerInitialConfig := math.Min(
					len(tc.initialMarketConfigs),
					int(tc.queryConfig.MaxQueries),
				) * 2
				for i, marketId := range tc.expectedMarketIdsCalled {
					marketConfigs := tc.initialMarketConfigs
					if i >= marketsPerInitialConfig {
						marketConfigs = tc.updateMarketConfigs
					}
					assertQueryHandlerCalledWithMarkets(
						t,
						queryHandler,
						[]types.MarketId{marketId},
						marketConfigs,
					)
				}
			}
		})
	}
}

func TestGetExchangeId(t *testing.T) {
	pf, err := NewPriceFetcher(
		constants.Exchange1_1MaxQueries_QueryConfig,
		constants.MultiMarketExchangeQueryDetails,
		&constants.Exchange1_3Markets_MutableExchangeMarketConfig,
		constants.MutableMarketConfigs_3Markets,
		&mocks.ExchangeQueryHandler{},
		log.NewNopLogger(),
		newTestPriceFetcherBufferedChannel(),
	)
	require.NoError(t, err)
	require.Equal(t, constants.Exchange1_1MaxQueries_QueryConfig.ExchangeId, pf.GetExchangeId())
}

// Test runSubTask behavior with different query handler responses
func TestRunSubTask_Mixed(t *testing.T) {
	tests := map[string]struct {
		responsePriceTimestamps    []*types.MarketPriceTimestamp
		responseUnavailableMarkets map[types.MarketId]error
		responseError              error

		expectedPrices []*types.MarketPriceTimestamp
		expectedErrors []error
	}{
		"Failure - failed to query exchange": {
			responseError:  exchangeQueryHandlerFailure,
			expectedErrors: []error{exchangeQueryHandlerFailure},
		},
		"Mixed - returned prices have a 0": {
			responsePriceTimestamps: []*types.MarketPriceTimestamp{
				{
					MarketId: 7,
					Price:    0,
				},
				constants.Market8_TimeT_Price1,
			},
			expectedPrices: []*types.MarketPriceTimestamp{
				constants.Market8_TimeT_Price1,
			},
			expectedErrors: []error{
				errors.New("Invalid price of 0 for exchange: 'Exchange1' and market: 7"),
			},
		},
		"Mixed - unavailable tickers": {
			responsePriceTimestamps: []*types.MarketPriceTimestamp{
				constants.Market8_TimeT_Price1,
			},
			responseUnavailableMarkets: map[types.MarketId]error{
				constants.MarketId8: tickerNotAvailable,
			},
			expectedPrices: []*types.MarketPriceTimestamp{
				constants.Market8_TimeT_Price1,
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			exchangeQueryConfig := constants.Exchange1_1MaxQueries_QueryConfig
			mutableExchangeMarketConfig := constants.Exchange1_1Markets_MutableExchangeMarketConfig
			mutableMarketConfigs := constants.MutableMarketConfigs_1Markets
			mockExchangeQueryHandler := &mocks.ExchangeQueryHandler{}
			rh := &daemontypes.RequestHandlerImpl{}

			mockExchangeQueryHandler.On(
				"Query",
				mock.AnythingOfType("*context.timerCtx"),
				mock.AnythingOfType("*types.ExchangeQueryDetails"),
				&mutableExchangeMarketConfig,
				mutableExchangeMarketConfig.GetMarketIds(),
				rh,
				generateMarketExponentsMap(mutableMarketConfigs),
			).
				Return(tc.responsePriceTimestamps, tc.responseUnavailableMarkets, tc.responseError)

			// Setup for sub-task iterations.
			bCh := newTestPriceFetcherBufferedChannel()

			pf, err := NewPriceFetcher(
				exchangeQueryConfig,
				constants.MultiMarketExchangeQueryDetails,
				&mutableExchangeMarketConfig,
				mutableMarketConfigs,
				mockExchangeQueryHandler,
				log.NewNopLogger(),
				bCh,
			)
			require.NoError(t, err)

			// We just need a valid input that matches the mock signature.
			pf.runSubTask(
				&daemontypes.RequestHandlerImpl{},
				mutableExchangeMarketConfig.GetMarketIds(),
				pf.getTaskLoopDefinition(),
			)

			actualErrors := make([]error, 0, len(tc.expectedErrors))
			var actualPrices []*types.MarketPriceTimestamp
			if len(tc.expectedPrices) > 0 {
				actualPrices = make([]*types.MarketPriceTimestamp, 0, len(tc.expectedPrices))
			}

			for i := 0; i < len(tc.expectedErrors)+len(tc.expectedPrices); i++ {
				value := <-bCh
				if value.Err != nil {
					actualErrors = append(actualErrors, value.Err)
				} else {
					actualPrices = append(actualPrices, value.Price)
				}
			}
			require.Equal(t, tc.expectedPrices, actualPrices)
			pricefeed.ErrorsEqual(t, tc.expectedErrors, actualErrors)
		})
	}
}

// ----------------- Generate Mock Instances ----------------- //
func generateMockExchangeQueryHandler() *mocks.ExchangeQueryHandler {
	mockExchangeQueryHandler := &mocks.ExchangeQueryHandler{}
	mockSingleMarketCalls(mockExchangeQueryHandler)
	// Mock multi-market call for 2 market test exchanges.
	mockMultiMarketCall(
		constants.Exchange1_2Markets_MutableExchangeMarketConfig,
		constants.MutableMarketConfigs_2Markets,
		mockExchangeQueryHandler,
	)
	// Mock multi-market call for 3 market test exchanges.
	mockMultiMarketCall(
		constants.Exchange1_3Markets_MutableExchangeMarketConfig,
		constants.MutableMarketConfigs_3Markets,
		mockExchangeQueryHandler,
	)
	// Mock multi-market call for 5 market test exchanges.
	mockMultiMarketCall(
		constants.Exchange1_5Markets_MutableExchangeMarketConfig,
		constants.MutableMarketConfigs_5Markets,
		mockExchangeQueryHandler,
	)

	return mockExchangeQueryHandler
}

// generateMarketExponentsMap generates a map of marketId to exponent for the given marketConfigs. This is
// the exponents map we would expect to be passed to the exchange query handler.
func generateMarketExponentsMap(marketConfigs []*types.MutableMarketConfig) map[types.MarketId]types.Exponent {
	marketExponents := make(map[types.MarketId]types.Exponent, len(marketConfigs))
	for _, marketConfigs := range marketConfigs {
		marketExponents[marketConfigs.Id] = marketConfigs.Exponent
	}
	return marketExponents
}

func mockSingleMarketCalls(mockExchangeQueryHandler *mocks.ExchangeQueryHandler) {
	// Support single market calls for all possible marketConfigs paired with the test exchange configs.
	initialMarketConfigsList := [][]*types.MutableMarketConfig{
		constants.MutableMarketConfigs_1Markets,
		constants.MutableMarketConfigs_2Markets,
		constants.MutableMarketConfigs_3Markets,
		constants.MutableMarketConfigs_5Markets,
	}
	for marketId, priceTimestamp := range constants.CanonicalMarketPriceTimestampResponses {
		for _, initialMarketConfigs := range initialMarketConfigsList {
			mockExchangeQueryHandler.On(
				"Query",
				mock.AnythingOfType("*context.timerCtx"),
				mock.AnythingOfType("*types.ExchangeQueryDetails"),
				mock.AnythingOfType("*types.MutableExchangeMarketConfig"),
				[]types.MarketId{marketId},
				&daemontypes.RequestHandlerImpl{},
				generateMarketExponentsMap(initialMarketConfigs),
			).Return([]*types.MarketPriceTimestamp{priceTimestamp}, nil, nil)
		}
	}
}

func mockMultiMarketCall(
	mutableExchangeMarketConfig types.MutableExchangeMarketConfig,
	mutableMarketConfigs []*types.MutableMarketConfig,
	mockExchangeQueryHandler *mocks.ExchangeQueryHandler,
) {
	markets := mutableExchangeMarketConfig.GetMarketIds()
	prices := make([]*types.MarketPriceTimestamp, 0, len(markets))
	for _, market := range markets {
		prices = append(prices, constants.CanonicalMarketPriceTimestampResponses[market])
	}

	mockExchangeQueryHandler.On(
		"Query",
		mock.AnythingOfType("*context.timerCtx"),
		mock.AnythingOfType("*types.ExchangeQueryDetails"),
		mock.AnythingOfType("*types.MutableExchangeMarketConfig"),
		markets,
		&daemontypes.RequestHandlerImpl{},
		generateMarketExponentsMap(mutableMarketConfigs),
	).Return(prices, nil, nil)
}

// ----------------- Helper Functions ----------------- //

// newTestPriceFetcherBufferedChannel returns a buffered channel with the default fixed size, suitably
// large enough to hold all the responses from multiple sub-task runs.
func newTestPriceFetcherBufferedChannel() chan *PriceFetcherSubtaskResponse {
	bCh := make(
		chan *PriceFetcherSubtaskResponse,
		pricefeed_cosntants.FixedBufferSize,
	)

	return bCh
}

// asserQueryHandlerCalledWithMarkets asserts that the query handler was called with the expected markets.
func assertQueryHandlerCalledWithMarkets(
	t *testing.T,
	queryHandler *mocks.ExchangeQueryHandler,
	markets []types.MarketId,
	marketConfigs []*types.MutableMarketConfig,
) {
	marketExponents := make(map[types.MarketId]types.Exponent)
	for _, market := range markets {
		marketExponents[market] = constants.CanonicalMarketExponents[market]
	}
	queryHandler.AssertCalled(
		t,
		"Query",
		mock.Anything,
		mock.AnythingOfType("*types.ExchangeQueryDetails"),
		mock.AnythingOfType("*types.MutableExchangeMarketConfig"),
		markets,
		&daemontypes.RequestHandlerImpl{},
		generateMarketExponentsMap(marketConfigs),
	)
}
