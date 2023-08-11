package handler

import (
	"context"
	"errors"
	"fmt"
	"github.com/dydxprotocol/v4/testutil/daemons/pricefeed"
	"net/http"
	"testing"
	"time"

	pf_constants "github.com/dydxprotocol/v4/daemons/pricefeed/client/constants"
	"github.com/dydxprotocol/v4/daemons/pricefeed/client/constants/exchange_common"
	"github.com/dydxprotocol/v4/daemons/pricefeed/client/types"
	"github.com/dydxprotocol/v4/lib"
	"github.com/dydxprotocol/v4/mocks"
	"github.com/dydxprotocol/v4/testutil/constants"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const (
	successStatus           = 200
	failStatus400           = 400
	failStatus500           = 500
	dummyPrice              = uint64(1)
	noPriceExponentMarketId = 100000
	FAKEUSD_ID              = 100001
	unavailableId           = 100002
	// No price exponent exists for this fake pair.
	noPriceExponentMarketSymbol = "INVALID-USD"
	noMarketSymbol              = "NO-MARKET-SYMBOL"
	unavailableSymbol           = "UNAVAILABLE"
	unavailableExponent         = -6
)

var (
	queryError              = errors.New("Failed to query exchange")
	priceFuncError          = errors.New("Failed to get Price")
	tickerNotAvailableError = errors.New("Ticker not available")
	baseEqd                 = &types.ExchangeQueryDetails{
		Url: "https://api.binance.us/api/v3/ticker/24hr?symbol=$",
		MarketSymbols: map[types.MarketId]string{
			exchange_common.MARKET_BTC_USD: constants.BtcUsdPair,
			exchange_common.MARKET_ETH_USD: constants.EthUsdPair,
			noPriceExponentMarketId:        noPriceExponentMarketSymbol,
			unavailableId:                  unavailableSymbol,
		},
	}
	testMarketExponentMap = generateTestMarketPriceExponentMap()
)

func TestQuery(t *testing.T) {
	lastUpdatedAt := time.Unix(0, 0)
	eqh := ExchangeQueryHandlerImpl{generateMockTimeProvider(lastUpdatedAt)}

	tests := map[string]struct {
		// parameters
		priceFunc func(
			response *http.Response,
			marketSymbolPriceExponentMap map[string]int32,
			medianizer lib.Medianizer,
		) (prices map[string]uint64, unavailable map[string]error, err error)
		marketIds      []types.MarketId
		requestHandler *mocks.RequestHandler

		// expectations
		expectedPrices      []*types.MarketPriceTimestamp
		expectedUnavailable map[types.MarketId]error
		expectApiRequest    bool
		expectedError       error
	}{
		"Success - single market": {
			priceFunc: priceFunc,
			marketIds: []types.MarketId{exchange_common.MARKET_BTC_USD},
			requestHandler: generateMockRequestHandler(
				CreateRequestUrl(baseEqd.Url, []string{constants.BtcUsdPair}),
				successStatus,
				nil,
			),
			expectApiRequest: true,
			expectedPrices: []*types.MarketPriceTimestamp{
				{
					Price:         dummyPrice,
					MarketId:      exchange_common.MARKET_BTC_USD,
					LastUpdatedAt: lastUpdatedAt,
				},
			},
		},
		"Success - multiple markets": {
			priceFunc: priceFunc,
			marketIds: []types.MarketId{exchange_common.MARKET_BTC_USD, exchange_common.MARKET_ETH_USD},
			requestHandler: generateMockRequestHandler(
				CreateRequestUrl(baseEqd.Url, []string{
					constants.BtcUsdPair,
					constants.EthUsdPair,
				}),
				successStatus,
				nil,
			),
			expectApiRequest: true,
			expectedPrices: []*types.MarketPriceTimestamp{
				{
					Price:         dummyPrice,
					MarketId:      exchange_common.MARKET_BTC_USD,
					LastUpdatedAt: lastUpdatedAt,
				},
				{
					Price:         dummyPrice,
					MarketId:      exchange_common.MARKET_ETH_USD,
					LastUpdatedAt: lastUpdatedAt,
				},
			},
		},
		"Success - multiple markets and unavailable symbol": {
			priceFunc: priceFuncWithValidAndUnavailableSymbols,
			marketIds: []types.MarketId{exchange_common.MARKET_BTC_USD, exchange_common.MARKET_ETH_USD, unavailableId},
			requestHandler: generateMockRequestHandler(
				CreateRequestUrl(baseEqd.Url, []string{
					constants.BtcUsdPair,
					constants.EthUsdPair,
					unavailableSymbol,
				}),
				successStatus,
				nil,
			),
			expectApiRequest: true,
			expectedPrices: []*types.MarketPriceTimestamp{
				{
					Price:         dummyPrice,
					MarketId:      exchange_common.MARKET_BTC_USD,
					LastUpdatedAt: lastUpdatedAt,
				},
				{
					Price:         dummyPrice,
					MarketId:      exchange_common.MARKET_ETH_USD,
					LastUpdatedAt: lastUpdatedAt,
				},
			},
			expectedUnavailable: map[types.MarketId]error{
				unavailableId: tickerNotAvailableError,
			},
		},
		"Failure - price function returns non-existent unavailable symbol": {
			priceFunc: priceFuncReturnsInvalidUnavailableSymbol,
			marketIds: []types.MarketId{exchange_common.MARKET_BTC_USD},
			requestHandler: generateMockRequestHandler(
				CreateRequestUrl(baseEqd.Url, []string{
					constants.BtcUsdPair,
				}),
				successStatus,
				nil,
			),
			expectApiRequest: true,
			expectedError:    fmt.Errorf("Severe unexpected error: no market id for symbol: %s", noMarketSymbol),
		},
		"Failure - no marketIds queried": {
			marketIds: []types.MarketId{},
			requestHandler: generateMockRequestHandler(
				CreateRequestUrl(baseEqd.Url, []string{constants.BtcUsdPair}),
				successStatus,
				nil,
			),
			expectApiRequest: false,
			expectedError:    errors.New("At least one marketId must be queried"),
		},
		"Failure - market symbol not defined for market": {
			marketIds: []types.MarketId{FAKEUSD_ID},
			requestHandler: generateMockRequestHandler(
				CreateRequestUrl(baseEqd.Url, []string{}),
				successStatus,
				nil,
			),
			expectApiRequest: false,
			expectedError:    fmt.Errorf("No market symbol for id: %v", FAKEUSD_ID),
		},
		"Failure - market price exponent not defined for market": {
			marketIds: []types.MarketId{noPriceExponentMarketId},
			requestHandler: generateMockRequestHandler(
				CreateRequestUrl(baseEqd.Url, []string{constants.BtcUsdPair}),
				successStatus,
				nil,
			),
			expectApiRequest: false,
			expectedError:    fmt.Errorf("No market price exponent for id: %v", noPriceExponentMarketId),
		},
		"Failure - query fails": {
			marketIds: []types.MarketId{exchange_common.MARKET_BTC_USD},
			requestHandler: generateMockRequestHandler(
				CreateRequestUrl(baseEqd.Url, []string{constants.BtcUsdPair}),
				successStatus,
				queryError,
			),
			expectApiRequest: true,
			expectedError:    queryError,
		},
		"Failure - unexpected API response code: 400": {
			priceFunc: priceFunc,
			marketIds: []types.MarketId{exchange_common.MARKET_BTC_USD},
			requestHandler: generateMockRequestHandler(
				CreateRequestUrl(baseEqd.Url, []string{constants.BtcUsdPair}),
				failStatus400,
				nil,
			),
			expectApiRequest: true,
			expectedError:    fmt.Errorf("%s %v", pf_constants.UnexpectedResponseStatusMessage, 400),
		},
		"Failure - unexpected API response code: 500": {
			priceFunc: priceFunc,
			marketIds: []types.MarketId{exchange_common.MARKET_BTC_USD},
			requestHandler: generateMockRequestHandler(
				CreateRequestUrl(baseEqd.Url, []string{constants.BtcUsdPair}),
				failStatus500,
				nil,
			),
			expectApiRequest: true,
			expectedError:    fmt.Errorf("%s %v", pf_constants.UnexpectedResponseStatusMessage, 500),
		},
		"Failure - PriceFunction fails": {
			priceFunc: priceFuncWithErr,
			marketIds: []types.MarketId{exchange_common.MARKET_BTC_USD},
			requestHandler: generateMockRequestHandler(
				CreateRequestUrl(baseEqd.Url, []string{constants.BtcUsdPair}),
				successStatus,
				nil,
			),
			expectApiRequest: true,
			expectedError:    priceFuncError,
		},
		"Failure - PriceFunction returns invalid response": {
			priceFunc: priceFuncWithInvalidResponse,
			marketIds: []types.MarketId{exchange_common.MARKET_BTC_USD, exchange_common.MARKET_ETH_USD},
			requestHandler: generateMockRequestHandler(
				CreateRequestUrl(baseEqd.Url, []string{
					constants.BtcUsdPair,
					constants.EthUsdPair,
				}),
				successStatus,
				nil,
			),
			expectApiRequest: true,
			expectedError: fmt.Errorf(
				"Severe unexpected error: no market id for symbol: %v",
				noMarketSymbol,
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			baseEqd.PriceFunction = tc.priceFunc

			prices, unavailableMarkets, err := eqh.Query(
				context.Background(),
				baseEqd,
				tc.marketIds,
				tc.requestHandler,
				testMarketExponentMap,
			)

			if tc.expectApiRequest {
				// Request argument is already tested in `generateMockRequestHandler`.
				tc.requestHandler.AssertCalled(t, "Get", context.Background(), mock.Anything)
			} else {
				tc.requestHandler.AssertNotCalled(t, "Get")
			}

			if tc.expectedError != nil {
				require.EqualError(t, err, tc.expectedError.Error())
				require.Nil(t, prices)
				require.Nil(t, unavailableMarkets)
			} else {
				require.NoError(t, err)
				require.ElementsMatch(t, tc.expectedPrices, prices)
				pricefeed.ErrorMapsEqual(t, tc.expectedUnavailable, unavailableMarkets)
			}
		})
	}
}

func generateMockTimeProvider(time time.Time) *mocks.TimeProvider {
	mockTimeProvider := &mocks.TimeProvider{}
	mockTimeProvider.On("Now").Return(time)

	return mockTimeProvider
}

func generateMockRequestHandler(url string, statusCode int, err error) *mocks.RequestHandler {
	mockRequestHandler := &mocks.RequestHandler{}
	mockRequestHandler.On("Get", context.Background(), url).Return(&http.Response{StatusCode: statusCode}, err)

	return mockRequestHandler
}

func generateTestMarketPriceExponentMap() map[types.MarketId]types.Exponent {
	marketExponents := make(map[types.MarketId]types.Exponent, len(pf_constants.StaticMarketPriceExponent)+1)
	for market, exponent := range pf_constants.StaticMarketPriceExponent {
		marketExponents[market] = exponent
	}
	marketExponents[unavailableId] = unavailableExponent
	return marketExponents
}

func priceFunc(
	response *http.Response,
	marketSymbolPriceExponentMap map[string]int32,
	medianizer lib.Medianizer,
) (prices map[string]uint64, unavailable map[string]error, err error) {
	prices = make(map[string]uint64, len(marketSymbolPriceExponentMap))
	for symbol := range marketSymbolPriceExponentMap {
		prices[symbol] = dummyPrice
	}
	return prices, nil, nil
}

func priceFuncWithInvalidResponse(
	response *http.Response,
	marketSymbolPriceExponentMap map[string]int32,
	medianizer lib.Medianizer,
) (prices map[string]uint64, unavailable map[string]error, err error) {
	prices = make(map[string]uint64, len(marketSymbolPriceExponentMap))
	for range marketSymbolPriceExponentMap {
		prices[noMarketSymbol] = dummyPrice
	}
	return prices, nil, nil
}

func priceFuncWithValidAndUnavailableSymbols(
	response *http.Response,
	marketSymbolPriceExponentMap map[string]int32,
	medianizer lib.Medianizer,
) (prices map[string]uint64, unavailable map[string]error, err error) {
	prices = make(map[string]uint64, len(marketSymbolPriceExponentMap))
	for symbol := range marketSymbolPriceExponentMap {
		if symbol != unavailableSymbol {
			prices[symbol] = dummyPrice
		}
	}
	return prices, map[string]error{unavailableSymbol: tickerNotAvailableError}, nil
}

func priceFuncReturnsInvalidUnavailableSymbol(
	response *http.Response,
	marketSymbolPriceExponentMap map[string]int32,
	medianizer lib.Medianizer,
) (prices map[string]uint64, unavailable map[string]error, err error) {
	return nil, map[string]error{noMarketSymbol: tickerNotAvailableError}, nil
}

func priceFuncWithErr(
	response *http.Response,
	marketSymbolPriceExponentMap map[string]int32,
	medianizer lib.Medianizer,
) (prices map[string]uint64, unavailable map[string]error, err error) {
	return nil, nil, priceFuncError
}
