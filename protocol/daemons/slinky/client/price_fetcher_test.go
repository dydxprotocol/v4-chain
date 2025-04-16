package client_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"cosmossdk.io/log"
	"github.com/dydxprotocol/slinky/service/servers/oracle/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	pricefeed_types "github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/types"
	pricefeedserver_types "github.com/dydxprotocol/v4-chain/protocol/daemons/server/types/pricefeed"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/slinky/client"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
)

func TestPriceFetcher(t *testing.T) {
	logger := log.NewTestLogger(t)
	mpf := mocks.NewMarketPairFetcher(t)
	slinky := mocks.NewOracleClient(t)
	slinky.On("Stop").Return(nil)
	var fetcher client.PriceFetcher

	t.Run("fetches prices on valid inputs", func(t *testing.T) {
		slinky.On("Start", mock.Anything).Return(nil).Once()
		slinky.On("Prices", mock.Anything, mock.Anything).
			Return(&types.QueryPricesResponse{
				Prices: map[string]string{
					"FOO/BAR": "100000000000",
				},
				Timestamp: time.Now(),
			}, nil).Once()
		mpf.On("GetIDForPair", mock.Anything).Return(uint32(1), nil).Once()

		fetcher = client.NewPriceFetcher(
			mpf,
			pricefeedserver_types.NewMarketToExchangePrices(pricefeed_types.MaxPriceAge),
			slinky,
			logger,
		)
		require.NoError(t, fetcher.Start(context.Background()))
		require.NoError(t, fetcher.FetchPrices(context.Background()))
		fetcher.Stop()
	})

	t.Run("errors on slinky.Prices failure", func(t *testing.T) {
		slinky.On("Start", mock.Anything).Return(nil).Once()
		slinky.On("Prices", mock.Anything, mock.Anything).
			Return(&types.QueryPricesResponse{}, fmt.Errorf("foobar")).Once()
		fetcher = client.NewPriceFetcher(
			mpf,
			pricefeedserver_types.NewMarketToExchangePrices(pricefeed_types.MaxPriceAge),
			slinky,
			logger,
		)

		require.NoError(t, fetcher.Start(context.Background()))
		require.Errorf(t, fetcher.FetchPrices(context.Background()), "foobar")
		fetcher.Stop()
	})

	t.Run("errors on slinky.Prices returning invalid currency pairs", func(t *testing.T) {
		slinky.On("Start", mock.Anything).Return(nil).Once()
		slinky.On("Prices", mock.Anything, mock.Anything).
			Return(&types.QueryPricesResponse{
				Prices: map[string]string{
					"FOOBAR": "100000000000",
				},
			}, nil).Once()
		fetcher = client.NewPriceFetcher(
			mpf,
			pricefeedserver_types.NewMarketToExchangePrices(pricefeed_types.MaxPriceAge),
			slinky,
			logger,
		)

		require.NoError(t, fetcher.Start(context.Background()))
		require.Errorf(t, fetcher.FetchPrices(context.Background()), "incorrectly formatted CurrencyPair")
		fetcher.Stop()
	})

	t.Run("no-ops on marketPairFetcher currency pair not found", func(t *testing.T) {
		slinky.On("Start", mock.Anything).Return(nil).Once()
		slinky.On("Prices", mock.Anything, mock.Anything).
			Return(&types.QueryPricesResponse{
				Prices: map[string]string{
					"FOO/BAR": "100000000000",
				},
				Timestamp: time.Now(),
			}, nil).Once()
		mpf.On("GetIDForPair", mock.Anything).Return(uint32(1), fmt.Errorf("not found")).Once()

		fetcher = client.NewPriceFetcher(
			mpf,
			pricefeedserver_types.NewMarketToExchangePrices(pricefeed_types.MaxPriceAge),
			slinky,
			logger,
		)
		require.NoError(t, fetcher.Start(context.Background()))
		require.NoError(t, fetcher.FetchPrices(context.Background()))
		fetcher.Stop()
	})

	t.Run("continues on non-parsable price data", func(t *testing.T) {
		slinky.On("Start", mock.Anything).Return(nil).Once()
		slinky.On("Prices", mock.Anything, mock.Anything).
			Return(&types.QueryPricesResponse{
				Prices: map[string]string{
					"FOO/BAR": "abc123",
				},
				Timestamp: time.Now(),
			}, nil).Once()
		mpf.On("GetIDForPair", mock.Anything).Return(uint32(1), nil).Once()

		fetcher = client.NewPriceFetcher(
			mpf,
			pricefeedserver_types.NewMarketToExchangePrices(pricefeed_types.MaxPriceAge),
			slinky,
			logger,
		)
		require.NoError(t, fetcher.Start(context.Background()))
		require.NoError(t, fetcher.FetchPrices(context.Background()))
		fetcher.Stop()
	})

	t.Run("no-ops on empty price response", func(t *testing.T) {
		slinky.On("Start", mock.Anything).Return(nil).Once()
		slinky.On("Prices", mock.Anything, mock.Anything).
			Return(&types.QueryPricesResponse{
				Prices:    map[string]string{},
				Timestamp: time.Now(),
			}, nil).Once()

		fetcher = client.NewPriceFetcher(
			mpf,
			pricefeedserver_types.NewMarketToExchangePrices(pricefeed_types.MaxPriceAge),
			slinky,
			logger,
		)
		require.NoError(t, fetcher.Start(context.Background()))
		require.NoError(t, fetcher.FetchPrices(context.Background()))
		fetcher.Stop()
	})
}
