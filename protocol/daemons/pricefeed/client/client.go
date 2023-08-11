package client

import (
	"context"
	"fmt"

	"github.com/dydxprotocol/v4/daemons/pricefeed/api"
	"github.com/dydxprotocol/v4/daemons/pricefeed/client/handler"
	"github.com/dydxprotocol/v4/daemons/pricefeed/client/types"
	"github.com/dydxprotocol/v4/lib"

	"github.com/cometbft/cometbft/libs/log"
)

const (
	// 5K is chosen to be >> than the number of messages an exchange could send in any period before the
	// price encoder is able to read the messages from the buffer, even if we add O(10-100) markets dynamically,
	// but not large enough to allow more than at most a few minutes of price messages to accumulate.
	FixedBufferSize = 1024 * 5
)

// Start begins a job that:
// 1) periodically queries prices from external data sources and saves the retrieved prices in an
// in-memory datastore
// 2) periodically sends the most recent prices to a gRPC server
func Start(
	ctx context.Context,
	socketAddress string,
	logger log.Logger,
	grpcClient lib.GrpcClient,
	exchangeFeedIdToStartupConfig map[types.ExchangeFeedId]*types.ExchangeStartupConfig,
	exchangeFeedIdToMarkets map[types.ExchangeFeedId][]types.MarketId,
	exchangeFeedIdToExchangeDetails map[types.ExchangeFeedId]types.ExchangeQueryDetails,
	priceUpdaterLoopDelayMs uint32,
	subTaskRunner SubTaskRunner,
) (err error) {
	conn, err := grpcClient.NewGrpcConnection(ctx, socketAddress)

	if err != nil {
		logger.Error("Failed to establish gRPC connection to socket address", "error", err)
		return err
	}

	// Defer closing gRPC connection until job completes.
	defer func() {
		if connErr := grpcClient.CloseConnection(conn); connErr != nil {
			err = connErr
		}
	}()

	exchangeFeedIds := make([]types.ExchangeFeedId, 0, len(exchangeFeedIdToStartupConfig))
	for exchangeId := range exchangeFeedIdToStartupConfig {
		exchangeFeedIds = append(exchangeFeedIds, exchangeId)
	}

	exchangeToMarketPrices, err := types.NewExchangeToMarketPrices(exchangeFeedIds)
	if err != nil {
		return err
	}

	// Start PriceEncoder and PriceFetcher per exchange.
	timeProvider := &lib.TimeProviderImpl{}
	// Iterate through all exchanges and call `StartPriceEncoder` and `StartPriceFetcher` respectively.
	for exchangeFeedId, exchangeConfig := range exchangeFeedIdToStartupConfig {
		// Instantiate shared buffered channel to be written to by the price fetcher and read from
		// by the price encoder.
		bCh := make(chan *PriceFetcherSubtaskResponse, FixedBufferSize)
		exchangeMarkets, exists := exchangeFeedIdToMarkets[exchangeFeedId]
		if !exists || len(exchangeMarkets) == 0 {
			return fmt.Errorf("no exchange information exists for exchangeFeedId: %v", exchangeFeedId)
		}
		exchangeDetails, exists := exchangeFeedIdToExchangeDetails[exchangeFeedId]
		if !exists {
			return fmt.Errorf("no exchange details exists for exchangeFeedId: %v", exchangeFeedId)
		}

		go subTaskRunner.StartPriceEncoder(
			exchangeFeedId,
			exchangeToMarketPrices,
			logger,
			bCh,
		)

		go subTaskRunner.StartPriceFetcher(
			types.ExchangeConfig{
				Markets:               exchangeMarkets,
				ExchangeStartupConfig: *exchangeConfig,
				IsMultiMarket:         exchangeDetails.IsMultiMarket,
			},
			&handler.ExchangeQueryHandlerImpl{TimeProvider: timeProvider},
			logger,
			bCh,
		)
	}

	// Start PriceUpdater to begin broadcasting prices.
	client := api.NewPriceFeedServiceClient(conn)
	// `StartPriceUpdater` does not run in a go-routine since it is used to block indefinitely
	// until the pricefeed daemon ends.
	// The price updater will read from an in-memory cache and send updates over gRPC for the
	// server to read.
	subTaskRunner.StartPriceUpdater(
		ctx,
		exchangeToMarketPrices,
		client,
		priceUpdaterLoopDelayMs,
		logger,
	)
	return nil
}
