package price_fetcher

import (
	"context"
	"fmt"
	"github.com/cosmos/cosmos-sdk/telemetry"
	daemontypes "github.com/dydxprotocol/v4-chain/protocol/daemons/types"
	"math/rand"
	"sync"
	"time"

	"cosmossdk.io/log"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/constants"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/handler"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/types"
	pricefeedmetrics "github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/metrics"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	gometrics "github.com/hashicorp/go-metrics"
	"gopkg.in/typ.v4/lists"
)

// PriceFetcherSubtaskResponse represents a transformed exchange API response that contains price
// info or an error.
type PriceFetcherSubtaskResponse struct {
	Price *types.MarketPriceTimestamp
	Err   error
}

// PriceFetcher fetches prices from an exchange by making a query based on the
// `exchangeConfig` specifications and then encodes the price or any associated error.
type PriceFetcher struct {
	exchangeQueryConfig types.ExchangeQueryConfig
	exchangeDetails     types.ExchangeQueryDetails
	queryHandler        handler.ExchangeQueryHandler
	logger              log.Logger
	bCh                 chan<- *PriceFetcherSubtaskResponse

	// mutableState contains all mutable state on the price fetcher is consolidated into a single object with access
	// and update protected by a mutex.
	mutableState *mutableState
}

// NewPriceFetcher creates a new PriceFetcher struct. It manages querying markets via goroutine
// queries to an exchange and encodes the responses or related errors into the shared buffered
// channel `bCh`.
func NewPriceFetcher(
	exchangeQueryConfig types.ExchangeQueryConfig,
	exchangeDetails types.ExchangeQueryDetails,
	mutableExchangeConfig *types.MutableExchangeMarketConfig,
	mutableMarketConfigs []*types.MutableMarketConfig,
	queryHandler handler.ExchangeQueryHandler,
	logger log.Logger,
	bCh chan<- *PriceFetcherSubtaskResponse,
) (
	*PriceFetcher,
	error,
) {
	// Configure price fetcher logger to have fetcher-specific metadata.
	pfLogger := logger.With(
		constants.SubmoduleLogKey,
		constants.PriceFetcherSubmoduleName,
		constants.ExchangeIdLogKey,
		exchangeQueryConfig.ExchangeId,
	)

	pf := &PriceFetcher{
		exchangeQueryConfig: exchangeQueryConfig,
		exchangeDetails:     exchangeDetails,
		queryHandler:        queryHandler,
		logger:              pfLogger,
		bCh:                 bCh,
		mutableState:        &mutableState{},
	}

	// This will instantiate the price fetcher's mutable state.
	err := pf.UpdateMutableExchangeConfig(mutableExchangeConfig, mutableMarketConfigs)
	if err != nil {
		return nil, err
	}

	return pf, nil
}

// GetExchangeId returns the exchange id for the exchange queried by the price fetcher.
// This method is added to support the MutableExchangeConfigUpdater interface.
func (p *PriceFetcher) GetExchangeId() types.ExchangeId {
	return p.exchangeQueryConfig.ExchangeId
}

// UpdateMutableExchangeConfig updates the price fetcher with the most current copy of the exchange config, as
// well as all markets supported by the exchange.
// This method is added to support the ExchangeConfigUpdater interface.
func (p *PriceFetcher) UpdateMutableExchangeConfig(
	newConfig *types.MutableExchangeMarketConfig,
	newMarketConfigs []*types.MutableMarketConfig,
) error {
	// 1. Validate new config.
	if newConfig.Id != p.exchangeQueryConfig.ExchangeId {
		return fmt.Errorf("PriceFetcher.UpdateMutableExchangeConfig: exchange id mismatch")
	}

	if err := newConfig.Validate(newMarketConfigs); err != nil {
		return fmt.Errorf("PriceFetcher.UpdateMutableExchangeConfig: invalid exchange config update: %w", err)
	}

	// 2. Derive price fetcher mutable state.
	// 2.A Compute market exponents.
	marketExponents := make(map[types.MarketId]types.Exponent, len(newMarketConfigs))
	for _, marketConfig := range newMarketConfigs {
		marketExponents[marketConfig.Id] = marketConfig.Exponent
	}

	// 2.B Compute market ids ring.
	marketIdsRing := lists.NewRing[types.MarketId](len(newConfig.GetMarketIds()))
	for _, marketId := range newConfig.GetMarketIds() {
		marketIdsRing.Value = marketId
		marketIdsRing = marketIdsRing.Next()
	}

	// 3. Perform update.
	p.mutableState.Update(newConfig, marketExponents, marketIdsRing)
	return nil
}

// getTaskLoopDefinition returns a snapshot of the current price fetcher mutable state.
func (p *PriceFetcher) getTaskLoopDefinition() *taskLoopDefinition {
	return p.mutableState.getTaskLoopDefinition(
		p.exchangeDetails.IsMultiMarket,
		p.getNumQueriesPerTaskLoop(),
	)
}

// isMultiMarketAndHasMarkets returns true if the price fetcher is a multi-market fetcher
// and is currently configured to query for 1 or more markets. In this case, the fetcher
// should execute a single subtask query for all markets. For multi-market exchanges, this
// will still be false if the price fetcher has no supported markets.
func (pf *PriceFetcher) isMultiMarketAndHasMarkets() bool {
	return pf.exchangeDetails.IsMultiMarket && len(pf.mutableState.GetMarketIds()) > 0
}

// getNumQueriesPerTaskLoop returns the number of queries that the price fetcher should execute
// on each task loop execution. For multi-market exchanges, this will always be 1.
// Otherwise, it will be the minimum of the number of markets supported by the exchange and
// the query limit specified in the exchange config to prevent exceeding the exchange's rate
// limit.
func (p *PriceFetcher) getNumQueriesPerTaskLoop() int {
	if p.exchangeDetails.IsMultiMarket {
		return 1
	}
	return lib.Min(
		int(p.exchangeQueryConfig.MaxQueries),
		len(p.mutableState.GetMarketIds()),
	)
}

// RunTaskLoop queries the exchange for market prices.
// Each goroutine makes a single exchange query for a specific set of one or more markets.
// RunTaskLoop blocks until all spawned goroutines have completed.
func (pf *PriceFetcher) RunTaskLoop(requestHandler daemontypes.RequestHandler) {
	taskLoopDefinition := pf.getTaskLoopDefinition()

	if pf.isMultiMarketAndHasMarkets() {
		pf.runSubTask(
			requestHandler,
			taskLoopDefinition.marketIds,
			taskLoopDefinition,
		)
	} else {
		// Run all subtasks in parallel and wait for each to complete.
		var waitGroup sync.WaitGroup
		for i := 0; i < len(taskLoopDefinition.marketIds); i++ {
			market := taskLoopDefinition.marketIds[i]
			waitGroup.Add(1)
			go func() {
				defer waitGroup.Done()
				pf.runSubTask(
					requestHandler,
					[]types.MarketId{market},
					taskLoopDefinition,
				)
			}()
		}
		waitGroup.Wait()
	}
}

// emitMarketAvailabilityMetrics emits telemetry that tracks whether a market was available when queried on an exchange.
// Success is tracked by (market, exchange) so that we can track the availability of each market on each exchange.
func emitMarketAvailabilityMetrics(exchangeId types.ExchangeId, id types.MarketId, available bool) {
	success := metrics.Success
	if !available {
		success = metrics.Error
	}
	telemetry.IncrCounterWithLabels(
		[]string{
			metrics.PricefeedDaemon,
			metrics.PriceFetcherQueryForMarket,
			success,
		},
		1,
		[]gometrics.Label{
			pricefeedmetrics.GetLabelForExchangeId(exchangeId),
			pricefeedmetrics.GetLabelForMarketId(id),
		},
	)
}

// runSubTask makes a single query to an exchange for market prices. This query can be for 1 or
// n markets.
// For single market exchanges, a task loop execution will execute multiple runSubTask goroutines, where
// each goroutine will query for a single market. To support this, we explicitly define the set of markets
// to query for in the `marketIds` parameter, even though in some cases it may be redundantly defined on
// the taskLoopDefinition. For multi-market exchanges, the taskLoop will execute exactly one subtask, and
// that subtask will query all markets defined in the taskLoopDefinition.
func (pf *PriceFetcher) runSubTask(
	requestHandler daemontypes.RequestHandler,
	marketIds []types.MarketId,
	taskLoopDefinition *taskLoopDefinition,
) {
	exchangeId := pf.exchangeQueryConfig.ExchangeId

	// Measure total latency for subtask to run for one API call and creating a context with timeout.
	defer metrics.ModuleMeasureSinceWithLabels(
		metrics.PricefeedDaemon,
		[]string{
			metrics.PricefeedDaemon,
			metrics.PriceFetcherSubtaskLoopAndSetCtxTimeout,
			metrics.Latency,
		},
		time.Now(),
		[]gometrics.Label{pricefeedmetrics.GetLabelForExchangeId(exchangeId)},
	)

	ctxWithTimeout, cancelFunc := context.WithTimeout(
		context.Background(),
		time.Duration(pf.exchangeQueryConfig.TimeoutMs)*time.Millisecond,
	)

	defer cancelFunc()

	// Measure total latency for subtask to run for one API call.
	defer metrics.ModuleMeasureSinceWithLabels(
		metrics.PricefeedDaemon,
		[]string{
			metrics.PricefeedDaemon,
			metrics.PriceFetcherSubtaskLoop,
			metrics.Latency,
		},
		time.Now(),
		[]gometrics.Label{pricefeedmetrics.GetLabelForExchangeId(exchangeId)},
	)

	prices, _, err := pf.queryHandler.Query(
		ctxWithTimeout,
		&pf.exchangeDetails,
		taskLoopDefinition.mutableExchangeConfig,
		marketIds,
		requestHandler,
		taskLoopDefinition.marketExponents,
	)

	// Emit metrics at the `AvailableMarketsSampleRate`.
	emitMetricsSample := rand.Float64() < metrics.AvailableMarketsSampleRate

	if err != nil {
		pf.writeToBufferedChannel(exchangeId, nil, err)

		// Since the query failed, report all markets as unavailable, according to the sampling rate.
		if emitMetricsSample {
			for _, marketId := range marketIds {
				emitMarketAvailabilityMetrics(exchangeId, marketId, false)
			}
		}

		return
	}

	// Track which markets were available when queried, and which were not, for telemetry.
	availableMarkets := make(map[types.MarketId]bool, len(marketIds))
	for _, marketId := range marketIds {
		availableMarkets[marketId] = false
	}

	for _, price := range prices {
		// No price should validly be zero. A price of zero points to an error in the API queried.
		if price.Price == uint64(0) {
			pf.writeToBufferedChannel(
				exchangeId,
				nil,
				fmt.Errorf(
					"Invalid price of 0 for exchange: '%v' and market: %v",
					exchangeId,
					price.MarketId,
				),
			)

			continue
		}

		// Log each new price (per-market per-exchange).
		pf.logger.Debug(
			"price_fetcher: Adding new price for market.",
			constants.PriceLogKey,
			price.Price,
			constants.MarketIdLogKey,
			price.MarketId,
			"LastUpdatedAt",
			price.LastUpdatedAt,
		)

		// Report market as available.
		availableMarkets[price.MarketId] = true

		pf.writeToBufferedChannel(exchangeId, price, err)
	}

	// Emit metrics on this exchange's market availability according to the sampling rate.
	if emitMetricsSample {
		for marketId, available := range availableMarkets {
			emitMarketAvailabilityMetrics(exchangeId, marketId, available)
		}
	}
}

// writeToBufferedChannel writes the (price, error) generated during querying to the price fetcher's
// buffered channel, which outputs the query result to the price encoder.
func (pf *PriceFetcher) writeToBufferedChannel(
	exchangeId types.ExchangeId,
	price *types.MarketPriceTimestamp,
	err error,
) {
	// Sanity check that the channel is not full already.
	if len(pf.bCh) == constants.FixedBufferSize {
		// Log if writing to buffered channel failed.
		pf.logger.Error("Pricefeed daemon's shared buffer is full.")
	}

	pf.bCh <- &PriceFetcherSubtaskResponse{
		Err:   err,
		Price: price,
	}
}
