package handler

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	gometrics "github.com/armon/go-metrics"
	"github.com/cosmos/cosmos-sdk/telemetry"
	"github.com/dydxprotocol/v4/daemons/pricefeed/client/constants"
	"github.com/dydxprotocol/v4/daemons/pricefeed/client/types"
	pricefeedmetrics "github.com/dydxprotocol/v4/daemons/pricefeed/metrics"
	"github.com/dydxprotocol/v4/lib"
	"github.com/dydxprotocol/v4/lib/metrics"
)

// ExchangeQueryHandlerImpl is the struct that implements the `ExchangeQueryHandler` interface.
type ExchangeQueryHandlerImpl struct {
	lib.TimeProvider
}

// Ensure the `ExchangeQueryHandlerImpl` struct is implemented at compile time
var _ ExchangeQueryHandler = (*ExchangeQueryHandlerImpl)(nil)

// ExchangeQueryHandler is an interface that encapsulates querying an exchange for price info.
type ExchangeQueryHandler interface {
	lib.TimeProvider
	Query(
		ctx context.Context,
		exchangeQueryDetails *types.ExchangeQueryDetails,
		marketIds []types.MarketId,
		requestHandler lib.RequestHandler,
		marketPriceExponent map[types.MarketId]types.Exponent,
	) (marketPriceTimestamps []*types.MarketPriceTimestamp, unavailableMarkets map[types.MarketId]error, err error)
}

// Query makes an API call to a specific exchange and returns the transformed response, including both valid prices
// and any unavailable markets with specific errors.
// 1) Validate `marketIds` contains at least one id.
// 2) Convert the list of `marketIds` to market symbols that are specific for a given exchange. Create a mapping of
// market symbols to price exponents and a reverse mapping of market symbol back to `MarketId`.
// 3) Make API call to an exchange and verify the response status code is not an error status code.
// 4) Transform the API response to market prices, while tracking unavailable symbols.
// 5) Return dual values:
// - a slice of `MarketPriceTimestamp`s that contains resolved market prices
// - a map of marketIds that could not be resolved with corresponding specific errors.
func (eqh *ExchangeQueryHandlerImpl) Query(
	ctx context.Context,
	exchangeQueryDetails *types.ExchangeQueryDetails,
	marketIds []types.MarketId,
	requestHandler lib.RequestHandler,
	marketPriceExponent map[types.MarketId]types.Exponent,
) (marketPriceTimestamps []*types.MarketPriceTimestamp, unavailableMarkets map[types.MarketId]error, err error) {
	// Measure latency to run query function per exchange.
	defer metrics.ModuleMeasureSinceWithLabels(
		metrics.PricefeedDaemon,
		[]string{
			metrics.PricefeedDaemon,
			metrics.PriceFetcherQueryExchange,
			metrics.Latency,
		},
		time.Now(),
		[]gometrics.Label{pricefeedmetrics.GetLabelForExchangeFeedId(exchangeQueryDetails.Exchange)},
	)
	// 1) Validate `marketIds` contains at least one id.
	if len(marketIds) == 0 {
		return nil, nil, errors.New("At least one marketId must be queried")
	}

	// 2) Convert the list of `marketIds` to market symbols that are specific for a given exchange. Create a mapping
	// of market symbols to price exponents and a reverse mapping of market symbol back to `MarketId`.
	marketSymbols := make([]string, 0, len(marketIds))
	marketSymbolPriceExponentMap := make(map[string]int32, len(marketIds))
	marketSymbolToMarketIdMap := make(map[string]types.MarketId, len(marketIds))
	for _, marketId := range marketIds {
		marketSymbol, ok := exchangeQueryDetails.MarketSymbols[marketId]
		if !ok {
			return nil, nil, fmt.Errorf("No market symbol for id: %v", marketId)
		}
		priceExponent, ok := marketPriceExponent[marketId]
		if !ok {
			return nil, nil, fmt.Errorf("No market price exponent for id: %v", marketId)
		}

		marketSymbols = append(marketSymbols, marketSymbol)
		marketSymbolPriceExponentMap[marketSymbol] = priceExponent
		marketSymbolToMarketIdMap[marketSymbol] = marketId

		// Measure count of requests sent.
		telemetry.IncrCounterWithLabels(
			[]string{
				metrics.PricefeedDaemon,
				metrics.HttpGetRequest,
				metrics.Count,
			},
			1,
			[]gometrics.Label{
				pricefeedmetrics.GetLabelForMarketId(marketId),
				pricefeedmetrics.GetLabelForExchangeFeedId(exchangeQueryDetails.Exchange),
			},
		)
	}

	// 3) Make API call to an exchange and verify the response status code is not an error status code.
	url := CreateRequestUrl(exchangeQueryDetails.Url, marketSymbols)

	beforeRequest := time.Now()
	response, err := requestHandler.Get(ctx, url)
	// Measure time to make API request for exchange.
	metrics.ModuleMeasureSinceWithLabels(
		metrics.PricefeedDaemon,
		[]string{
			metrics.PricefeedDaemon,
			metrics.ExchangeQueryHandlerApiRequest,
			metrics.Latency,
		},
		beforeRequest,
		[]gometrics.Label{
			pricefeedmetrics.GetLabelForExchangeFeedId(exchangeQueryDetails.Exchange),
		},
	)
	if err != nil {
		return nil, nil, err
	}

	// Measure count of exchange API calls as well as what the status code is.
	telemetry.IncrCounterWithLabels(
		[]string{
			metrics.PricefeedDaemon,
			metrics.HttpGetResponse,
			metrics.Count,
		},
		1,
		[]gometrics.Label{
			metrics.GetLabelForIntValue(metrics.StatusCode, response.StatusCode),
			pricefeedmetrics.GetLabelForExchangeFeedId(exchangeQueryDetails.Exchange),
		},
	)

	// Verify response is not 4xx or 5xx.
	if response.StatusCode < 200 || response.StatusCode > 299 {
		return nil, nil, fmt.Errorf("%s %v", constants.UnexpectedResponseStatusMessage, response.StatusCode)
	}

	// 4) Transform the API response to market prices, while tracking unavailable symbols.
	prices, unavailableSymbols, err := exchangeQueryDetails.PriceFunction(
		response,
		marketSymbolPriceExponentMap,
		&lib.MedianizerImpl{},
	)
	if err != nil {
		return nil, nil, err
	}

	// 5) Insert prices into MarketPriceTimestamp struct slice, convert unavailable symbols back into marketIds,
	// and return.
	marketPriceTimestamps = make([]*types.MarketPriceTimestamp, 0, len(prices))
	now := eqh.Now()

	for marketSymbol, price := range prices {
		marketId, ok := marketSymbolToMarketIdMap[marketSymbol]
		if !ok {
			return nil, nil, fmt.Errorf("Severe unexpected error: no market id for symbol: %v", marketSymbol)
		}

		marketPriceTimestamp := &types.MarketPriceTimestamp{
			MarketId:      marketId,
			Price:         price,
			LastUpdatedAt: now,
		}

		marketPriceTimestamps = append(marketPriceTimestamps, marketPriceTimestamp)
	}

	unavailableMarkets = make(map[types.MarketId]error, len(unavailableSymbols))
	for marketSymbol, error := range unavailableSymbols {
		marketId, ok := marketSymbolToMarketIdMap[marketSymbol]
		if !ok {
			return nil, nil, fmt.Errorf("Severe unexpected error: no market id for symbol: %v", marketSymbol)
		}
		unavailableMarkets[marketId] = error
	}

	return marketPriceTimestamps, unavailableMarkets, nil
}

func CreateRequestUrl(baseUrl string, marketSymbols []string) string {
	return strings.Replace(baseUrl, "$", strings.Join(marketSymbols, ","), -1)
}
