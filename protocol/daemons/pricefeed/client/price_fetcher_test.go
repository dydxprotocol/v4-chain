package client

import (
	"errors"
	"fmt"
	"github.com/dydxprotocol/v4/testutil/daemons/pricefeed"
	"testing"

	"github.com/cometbft/cometbft/libs/log"
	client_constants "github.com/dydxprotocol/v4/daemons/pricefeed/client/constants"
	"github.com/dydxprotocol/v4/daemons/pricefeed/client/types"
	"github.com/dydxprotocol/v4/lib"
	"github.com/dydxprotocol/v4/mocks"
	"github.com/dydxprotocol/v4/testutil/constants"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gopkg.in/typ.v4/lists"
)

const (
	taskLoopIterations = 2
)

var (
	exchangeQueryHandlerFailure = errors.New("Failed to query exchange")
	symbolNotAvailable          = errors.New("Symbol not listed")
)

// Test different exchange configurations
func TestRunTaskLoop(t *testing.T) {
	tests := map[string]struct {
		// parameters
		exchangeConfig types.ExchangeConfig

		// expectations
		expectedMarketIdsCalled []types.MarketId
	}{
		"No markets": {
			exchangeConfig: constants.Exchange1_NoMarkets_0MaxQueries_Config,
		},
		"Num markets equals max query markets where there is only 1 market": {
			exchangeConfig: constants.Exchange1_1Markets_1MaxQueries_Config,
			expectedMarketIdsCalled: []types.MarketId{
				constants.MarketId7,
				constants.MarketId7,
			},
		},
		"Num markets < max query markets": {
			exchangeConfig: constants.Exchange1_1Markets_2MaxQueries_Config,
			expectedMarketIdsCalled: []types.MarketId{
				constants.MarketId7,
				constants.MarketId7,
			},
		},
		"Num markets equals max query markets": {
			exchangeConfig: constants.Exchange1_2Markets_2MaxQueries_Config,
			expectedMarketIdsCalled: []types.MarketId{
				constants.MarketId7,
				constants.MarketId8,
				constants.MarketId7,
				constants.MarketId8,
			},
		},
		"Multi-market, 2 markets": {
			exchangeConfig: constants.Exchange1_2Markets_Multimarket_Config,
			expectedMarketIdsCalled: []types.MarketId{
				constants.MarketId7,
				constants.MarketId8,
				constants.MarketId7,
				constants.MarketId8,
			},
		},
		"Num markets greater than max query markets": {
			exchangeConfig: constants.Exchange1_3Markets_2MaxQueries_Config,
			expectedMarketIdsCalled: []types.MarketId{
				constants.MarketId7,
				constants.MarketId8,
				constants.MarketId9,
				constants.MarketId7,
			},
		},
		"Multi-market, 5 markets": {
			exchangeConfig: constants.Exchange1_5Markets_Multimarket_Config,
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
			// Setup for sub-task iterations.
			marketIdsRing, bCh := setupForPriceFetcherTaskLoop(
				tc.exchangeConfig,
			)

			queryHandler := generateMockExchangeQueryHandler()

			pf := NewPriceFetcher(
				tc.exchangeConfig,
				queryHandler,
				marketIdsRing,
				log.NewNopLogger(),
				bCh,
			)

			// Run sub-task a specified number of iterations.
			for i := 0; i < taskLoopIterations; i++ {
				pf.RunTaskLoop(&lib.RequestHandlerImpl{})
			}

			// Will hang until tests timeout if bCh is not full.
			bufferedResponses := []*types.MarketPriceTimestamp{}
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
			if tc.exchangeConfig.IsMultiMarket {
				expectedQueries = taskLoopIterations
				marketsPerCall = len(tc.exchangeConfig.Markets)
			}
			queryHandler.AssertNumberOfCalls(t, "Query", expectedQueries)
			for i := 0; i < len(tc.expectedMarketIdsCalled); i = i + marketsPerCall {
				queryHandler.AssertCalled(
					t,
					"Query",
					mock.Anything,
					mock.AnythingOfType("*types.ExchangeQueryDetails"),
					tc.expectedMarketIdsCalled[i:i+marketsPerCall],
					&lib.RequestHandlerImpl{},
					client_constants.StaticMarketPriceExponent,
				)
			}
		})
	}
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
				errors.New("Invalid price of 0 for exchange: 1 and market: 7"),
			},
		},
		"Mixed - unavailable symbols": {
			responsePriceTimestamps: []*types.MarketPriceTimestamp{
				constants.Market8_TimeT_Price1,
			},
			responseUnavailableMarkets: map[types.MarketId]error{
				constants.MarketId8: symbolNotAvailable,
			},
			expectedPrices: []*types.MarketPriceTimestamp{
				constants.Market8_TimeT_Price1,
			},
			expectedErrors: []error{
				fmt.Errorf("Market 8 unavailable on exchange 1 (%w)", symbolNotAvailable),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			exchangeConfig := constants.Exchange1_1Markets_1MaxQueries_Config
			mockExchangeQueryHandler := &mocks.ExchangeQueryHandler{}
			rh := &lib.RequestHandlerImpl{}
			mockExchangeQueryHandler.On(
				"Query",
				mock.AnythingOfType("*context.timerCtx"),
				mock.AnythingOfType("*types.ExchangeQueryDetails"),
				exchangeConfig.Markets,
				rh,
				client_constants.StaticMarketPriceExponent,
			).
				Return(tc.responsePriceTimestamps, tc.responseUnavailableMarkets, tc.responseError)

			// Setup for sub-task iterations.
			marketIdsRing, bCh := setupForPriceFetcherTaskLoop(exchangeConfig)

			pf := NewPriceFetcher(
				exchangeConfig,
				mockExchangeQueryHandler,
				marketIdsRing,
				log.NewNopLogger(),
				bCh,
			)

			// We just need a valid input that matches the mock signature.
			pf.runSubTask(&lib.RequestHandlerImpl{}, exchangeConfig.Markets)

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
	mockMultiMarketCall(constants.Exchange1_5Markets_Multimarket_Config.Markets, mockExchangeQueryHandler)
	mockMultiMarketCall(constants.Exchange1_2Markets_Multimarket_Config.Markets, mockExchangeQueryHandler)
	return mockExchangeQueryHandler
}

func mockSingleMarketCalls(mockExchangeQueryHandler *mocks.ExchangeQueryHandler) {
	rh := &lib.RequestHandlerImpl{}
	for marketId, priceTimestamp := range constants.CanonicalMarketPriceTimestampResponses {
		mockExchangeQueryHandler.On(
			"Query",
			mock.AnythingOfType("*context.timerCtx"),
			mock.AnythingOfType("*types.ExchangeQueryDetails"),
			[]types.MarketId{marketId},
			rh,
			client_constants.StaticMarketPriceExponent,
		).Return([]*types.MarketPriceTimestamp{priceTimestamp}, nil, nil)
	}
}

func mockMultiMarketCall(
	markets []types.MarketId,
	mockExchangeQueryHandler *mocks.ExchangeQueryHandler,
) {
	prices := make([]*types.MarketPriceTimestamp, 0, len(markets))
	for _, market := range markets {
		prices = append(prices, constants.CanonicalMarketPriceTimestampResponses[market])
	}

	mockExchangeQueryHandler.On(
		"Query",
		mock.AnythingOfType("*context.timerCtx"),
		mock.AnythingOfType("*types.ExchangeQueryDetails"),
		markets,
		&lib.RequestHandlerImpl{},
		client_constants.StaticMarketPriceExponent,
	).Return(prices, nil, nil)
}

// ----------------- Helper Functions ----------------- //

func setupForPriceFetcherTaskLoop(
	exchangeConfig types.ExchangeConfig,
) (*lists.Ring[types.MarketId], chan *PriceFetcherSubtaskResponse) {
	// Create ring that holds all markets for an exchange.
	marketIds := exchangeConfig.Markets
	marketIdsRing := lists.NewRing[types.MarketId](len(marketIds))
	for _, marketId := range marketIds {
		marketIdsRing.Value = marketId
		marketIdsRing = marketIdsRing.Next()
	}
	bCh := make(
		chan *PriceFetcherSubtaskResponse,
		FixedBufferSize,
	)

	return marketIdsRing, bCh
}
