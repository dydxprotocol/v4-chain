//go:build all || integration_test

package client_test

import (
	"fmt"
	"net"
	"sync"
	"time"

	"cosmossdk.io/log"
	appflags "github.com/dydxprotocol/v4-chain/protocol/app/flags"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/flags"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/constants/exchange_common"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/price_function/bitfinex"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/price_function/testexchange"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/types"
	pricefeed_types "github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/types"
	daemonserver "github.com/dydxprotocol/v4-chain/protocol/daemons/server"
	servertypes "github.com/dydxprotocol/v4-chain/protocol/daemons/server/types"
	pricefeedserver_types "github.com/dydxprotocol/v4-chain/protocol/daemons/server/types/pricefeed"
	daemontypes "github.com/dydxprotocol/v4-chain/protocol/daemons/types"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/appoptions"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/daemons/pricefeed/exchange_config"
	grpc_util "github.com/dydxprotocol/v4-chain/protocol/testutil/grpc"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/pricefeed"
	pricetypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"

	"testing"
)

var (
	testExchangeToQueryDetails = map[types.ExchangeId]types.ExchangeQueryDetails{
		exchange_common.EXCHANGE_ID_TEST_EXCHANGE: {
			Exchange:      testexchange.TestExchangeDetails.Exchange,
			PriceFunction: testexchange.TestExchangeDetails.PriceFunction,
			Url:           fmt.Sprintf("http://127.0.0.1:%s/ticker?symbol=$", testexchange.TestExchangePort),
		},
		exchange_common.EXCHANGE_ID_BITFINEX: {
			Exchange:      exchange_common.EXCHANGE_ID_BITFINEX,
			Url:           fmt.Sprintf("http://127.0.0.1:%s/bitfinex-ticker?symbols=$", testexchange.TestExchangePort),
			PriceFunction: bitfinex.BitfinexPriceFunction,
			IsMultiMarket: true,
		},
	}

	// Initialize the daemon client with Bitfinex and TestExchange exchanges. Shorten intervals for testing
	// since we're using a mock exchange server on localhost with no rate limits.
	testExchangeQueryConfigs = map[types.ExchangeId]*types.ExchangeQueryConfig{
		exchange_common.EXCHANGE_ID_TEST_EXCHANGE: {
			ExchangeId: exchange_common.EXCHANGE_ID_TEST_EXCHANGE,
			IntervalMs: 100,
			TimeoutMs:  3_000,
			MaxQueries: 2,
		},
		exchange_common.EXCHANGE_ID_BITFINEX: {
			ExchangeId: exchange_common.EXCHANGE_ID_BITFINEX,
			IntervalMs: 100,
			TimeoutMs:  3_000,
			MaxQueries: 2,
		},
	}

	// defaultMarketParams: BTC on Bitfinex and TestExchange.
	defaultMarketParams = []pricetypes.MarketParam{
		{
			Id:       0,
			Pair:     "BTC-USD",
			Exponent: -5,
			ExchangeConfigJson: `{"exchanges":[{"exchangeName":"Bitfinex","ticker":"tBTCUSD"},` +
				`{"exchangeName":"TestExchange","ticker":"BTC-USD"}]}`,
			MinExchanges:      2,
			MinPriceChangePpm: 1,
		},
	}

	// marketParams_AddMarkets: adds ETH on Bitfinex and TestExchange, LINK on test exchange, to the
	// defaultMarketParams.
	marketParams_AddMarkets = []pricetypes.MarketParam{
		{
			Id:       0,
			Pair:     "BTC-USD",
			Exponent: -5,
			ExchangeConfigJson: `{"exchanges":[{"exchangeName":"Bitfinex","ticker":"tBTCUSD"},` +
				`{"exchangeName":"TestExchange","ticker":"BTC-USD"}]}`,
			MinExchanges:      2,
			MinPriceChangePpm: 1,
		},
		{
			Id:       1,
			Pair:     "ETH-USD",
			Exponent: -6,
			ExchangeConfigJson: `{"exchanges":[{"exchangeName":"Bitfinex","ticker":"tETHUSD"},` +
				`{"exchangeName":"TestExchange","ticker":"ETH-USD"}]}`,
			MinExchanges:      2,
			MinPriceChangePpm: 1,
		},
		{
			Id:                 2,
			Pair:               "LINK-USD",
			Exponent:           -8,
			ExchangeConfigJson: `{"exchanges":[{"exchangeName":"TestExchange","ticker":"LINK-USD"}]}`,
			MinExchanges:       1,
			MinPriceChangePpm:  1,
		},
	}

	// marketParams_AddMarketsWithAdjustments: adds ETH on Bitfinex and TestExchange, LINK on test exchange, and
	// USDT on the test exchange. ETH and LINK prices are all adjusted by USDT, and should be 90% of the un-adjusted
	// price from the non-converting test case.
	marketParams_AddMarketsWithAdjustments = []pricetypes.MarketParam{
		{
			Id:       0,
			Pair:     "BTC-USD",
			Exponent: -5,
			ExchangeConfigJson: `{"exchanges":[{"exchangeName":"Bitfinex","ticker":"tBTCUSD"},` +
				`{"exchangeName":"TestExchange","ticker":"BTC-USD"}]}`,
			MinExchanges:      2,
			MinPriceChangePpm: 1,
		},
		{
			Id:       1,
			Pair:     "ETH-USD",
			Exponent: -6,
			ExchangeConfigJson: `{"exchanges":[` +
				`{"exchangeName":"Bitfinex","ticker":"tETHUSD","adjustByMarket":"USDT-USD"},` +
				`{"exchangeName":"TestExchange","ticker":"ETH-USD","adjustByMarket":"USDT-USD"}]}`,
			MinExchanges:      2,
			MinPriceChangePpm: 1,
		},
		{
			Id:       2,
			Pair:     "LINK-USD",
			Exponent: -8,
			ExchangeConfigJson: `{"exchanges":[{"exchangeName":"TestExchange","ticker":"LINK-USD",` +
				`"adjustByMarket":"USDT-USD"}]}`,
			MinExchanges:      1,
			MinPriceChangePpm: 1,
		},
		{
			Id:                 33,
			Pair:               "USDT-USD",
			Exponent:           -9,
			ExchangeConfigJson: `{"exchanges":[{"exchangeName":"TestExchange","ticker":"USDT-USD"}]}`,
			MinExchanges:       1,
			MinPriceChangePpm:  1,
		},
	}

	marketParams_PartialUpdate = []pricetypes.MarketParam{
		{
			Id:       0,
			Pair:     "BTC-USD",
			Exponent: -5,
			// Invalid exchange name.
			ExchangeConfigJson: `{"exchanges":[{"exchangeName":"Nonexistent","ticker":"tBTCUSD"},` +
				`{"exchangeName":"TestExchange","ticker":"BTC-USD"}]}`,
			MinExchanges:      2,
			MinPriceChangePpm: 1,
		},
		{
			Id:       1,
			Pair:     "ETH-USD",
			Exponent: -6,
			ExchangeConfigJson: `{"exchanges":[{"exchangeName":"Bitfinex","ticker":"tETHUSD"},` +
				`{"exchangeName":"TestExchange","ticker":"ETH-USD"}]}`,
			MinExchanges:      2,
			MinPriceChangePpm: 1,
		},
	}

	// The test exchange adds 100 to the set prices for the Bitfinex exchange response, so median prices will
	// fall halfway between the set price and set price + 100.
	expectedMedianBtcPrice = uint64(100_005_000_000)
	expectedMedianEthPrice = uint64(2_000_050_000_000)
	testExchangeLinkPrice  = uint64(300_000_000_000_000) // Link not available on Bitfinex.

	// USDT is set to $.90, so expect 90% of the expected median price after applying USDT conversion.
	expectedAdjustedMedianEthPrice = uint64(1_800_045_000_000)
	expectedAdjustedLinkPrice      = uint64(270_000_000_000_000)

	expectedPrices1Market = map[types.MarketId]uint64{
		exchange_config.MARKET_BTC_USD: expectedMedianBtcPrice,
	}

	// expectedPricesPartialUpdate preserves the expected price of BTC, ignoring the invalid update params, and also
	// updates to expect the median price of ETH.
	expectedPricesPartialUpdate = map[types.MarketId]uint64{
		exchange_config.MARKET_BTC_USD: expectedMedianBtcPrice,
		exchange_config.MARKET_ETH_USD: expectedMedianEthPrice,
	}

	expectedPrices3Markets = map[types.MarketId]uint64{
		exchange_config.MARKET_BTC_USD:  expectedMedianBtcPrice,
		exchange_config.MARKET_ETH_USD:  expectedMedianEthPrice,
		exchange_config.MARKET_LINK_USD: testExchangeLinkPrice,
	}

	expectedPrices3MarketsWithConversions = map[types.MarketId]uint64{
		exchange_config.MARKET_BTC_USD:  expectedMedianBtcPrice,
		exchange_config.MARKET_ETH_USD:  expectedAdjustedMedianEthPrice,
		exchange_config.MARKET_LINK_USD: expectedAdjustedLinkPrice,
	}

	// 5s is chosen to give us a comfortable margin of error for prices to make it through the
	// encoder and price updater go routines and into the prices cache through a gRPC call to the
	// daemon server. We wouldn't want a price to expire before it makes it off of the daemon.
	testPriceCacheExpirationDuration = 5 * time.Second
)

type PriceDaemonIntegrationTestSuite struct {
	suite.Suite
	daemonFlags        flags.DaemonFlags
	appFlags           appflags.Flags
	exchangeServer     *pricefeed.ExchangeServer
	daemonServer       *daemonserver.Server
	exchangePriceCache *pricefeedserver_types.MarketToExchangePrices
	healthMonitor      *servertypes.HealthMonitor

	pricesMockQueryServer *mocks.QueryServer
	pricesGrpcServer      *grpc.Server

	pricefeedDaemon *client.Client

	activeServers sync.WaitGroup
}

func TestPriceDaemonIntegrationTestSuite(t *testing.T) {
	suite.Run(t, &PriceDaemonIntegrationTestSuite{})
}

// mockAllMarketParamsResponse mocks the response from the prices query server for the AllMarketParams query.
// This endpoint is used by the daemon to detect changes to the market params and update the exchange queries
// appropriately.
func (s *PriceDaemonIntegrationTestSuite) mockAllMarketParamsResponse(
	response *pricetypes.QueryAllMarketParamsResponse,
) {
	s.pricesMockQueryServer.On("AllMarketParams", mock.Anything, mock.Anything).Return(
		response,
		nil,
	)
}

// mockAllMarketParamsResponseNTimes mocks the response from the prices query server for the AllMarketParams query
// n times. This endpoint is used by the daemon to detect changes to the market params and update the exchange queries
// appropriately. To change the endpoint response, call this function again with a different response.
func (s *PriceDaemonIntegrationTestSuite) mockAllMarketParamsResponseNTimes(
	response *pricetypes.QueryAllMarketParamsResponse,
	n int,
) {
	s.pricesMockQueryServer.On("AllMarketParams", mock.Anything, mock.Anything).Return(
		response,
		nil,
	).Times(n)
}

func (s *PriceDaemonIntegrationTestSuite) SetupTest() {
	// Configure test to use test exchange.
	s.exchangeServer = pricefeed.NewExchangeServer()
	s.exchangeServer.SetPrice(exchange_config.MARKET_BTC_USD, 1_000_000)
	s.exchangeServer.SetPrice(exchange_config.MARKET_ETH_USD, 2_000_000)
	s.exchangeServer.SetPrice(exchange_config.MARKET_LINK_USD, 3_000_000)

	// Set USDT to 90 cents.
	s.exchangeServer.SetPrice(exchange_config.MARKET_USDT_USD, .9)

	// Save daemon flags to use for client startup.
	s.daemonFlags = flags.GetDefaultDaemonFlags()
	s.appFlags = appflags.GetFlagValuesFromOptions(appoptions.GetDefaultTestAppOptions("", nil))

	// Configure mock daemon server with prices cache.
	s.daemonServer = daemonserver.NewServer(
		log.NewNopLogger(),
		grpc.NewServer(),
		&daemontypes.FileHandlerImpl{},
		s.daemonFlags.Shared.SocketAddress,
	)

	s.healthMonitor = servertypes.NewHealthMonitor(
		servertypes.DaemonStartupGracePeriod,
		servertypes.HealthCheckPollFrequency,
		log.NewNopLogger(),
		flags.GetDefaultDaemonFlags().Shared.PanicOnDaemonFailureEnabled, // Use default behavior for testing
	)

	s.exchangePriceCache = pricefeedserver_types.NewMarketToExchangePrices(pricefeed_types.MaxPriceAge)
	s.daemonServer.WithPriceFeedMarketToExchangePrices(s.exchangePriceCache)

	// Create a gRPC server running on the default port and attach the mock prices query service.
	s.pricesMockQueryServer = &mocks.QueryServer{}
	s.pricesGrpcServer = grpc.NewServer()
	pricetypes.RegisterQueryServer(s.pricesGrpcServer, s.pricesMockQueryServer)

	// Start daemon server and the gRPC server, with a wait group to ensure the test doesn't exit before they're
	// terminated.
	s.activeServers.Add(1)
	go func() {
		defer s.activeServers.Done()
		s.daemonServer.Start()
	}()

	s.activeServers.Add(1)
	go func() {
		defer s.activeServers.Done()
		ls, err := net.Listen("tcp", s.appFlags.GrpcAddress)
		s.Require().NoError(err)
		err = s.pricesGrpcServer.Serve(ls)
		s.Require().NoError(err)
	}()
}

func (s *PriceDaemonIntegrationTestSuite) TearDownTest() {
	err := s.exchangeServer.CleanUp()
	s.Require().NoError(err)

	// Stop all running services.
	s.pricefeedDaemon.Stop()
	s.daemonServer.Stop()
	s.pricesGrpcServer.GracefulStop()

	s.activeServers.Wait()
}

// startClient starts the pricefeed daemon client.
func (s *PriceDaemonIntegrationTestSuite) startClient() {
	s.pricefeedDaemon = client.StartNewClient(
		grpc_util.Ctx,
		s.daemonFlags,
		s.appFlags,
		log.NewNopLogger(),
		&daemontypes.GrpcClientImpl{},
		testExchangeQueryConfigs,
		testExchangeToQueryDetails,
		&client.SubTaskRunnerImpl{},
	)
	err := s.healthMonitor.RegisterService(
		s.pricefeedDaemon,
		time.Duration(s.daemonFlags.Shared.MaxDaemonUnhealthySeconds)*time.Second,
	)
	s.Require().NoError(err)
}

// expectPricesWithTimeout waits for the exchange price cache to contain the expected prices, with a timeout.
// This is used to give the daemon time to update the exchange queries, receive and process responses, and send them
// to the daemon server.
func (s *PriceDaemonIntegrationTestSuite) expectPricesWithTimeout(
	expectedPrices map[types.MarketId]uint64,
	marketParams []pricetypes.MarketParam,
	timeout time.Duration,
) {
	start := time.Now()

	for {
		// Fail if we've timed out.
		if time.Since(start) > timeout {
			s.Require().Fail("timed out waiting for expected prices")
			return
		}

		// Poll every 100 milliseconds.
		time.Sleep(100 * time.Millisecond)

		// Check if the prices cache contains the expected prices.
		prices := s.exchangePriceCache.GetValidMedianPrices(log.NewNopLogger(), marketParams, time.Now())
		if len(prices) != len(expectedPrices) {
			continue
		}

		allPricesMatch := true

		for marketId, expectedPrice := range expectedPrices {
			actualPrice, ok := prices[marketId]
			if !ok || actualPrice != expectedPrice {
				allPricesMatch = false
				break
			}
		}
		if allPricesMatch {
			return
		}
	}
}

// TestPriceDaemon tests that the pricefeed daemon produces the expected price updates and sends them to the
// daemon server.
func (s *PriceDaemonIntegrationTestSuite) TestPriceDaemon() {
	// Set up the mock prices query server to return market params that reflect the same markets and exchanges the
	// daemon is initialized with.
	s.mockAllMarketParamsResponse(&pricetypes.QueryAllMarketParamsResponse{
		MarketParams: defaultMarketParams,
	})

	s.startClient()

	s.expectPricesWithTimeout(
		expectedPrices1Market,
		defaultMarketParams,
		testPriceCacheExpirationDuration,
	)
}

// TestUpdateMarkets_AddMarket tests that the pricefeed daemon produces prices for a new market after it is added.
func (s *PriceDaemonIntegrationTestSuite) TestUpdateMarkets_AddMarket() {
	// Start the daemon with a single market. Then, update the endpoint to return market params that have new markets.
	s.mockAllMarketParamsResponseNTimes(
		&pricetypes.QueryAllMarketParamsResponse{
			MarketParams: defaultMarketParams,
		},
		1,
	)
	s.mockAllMarketParamsResponseNTimes(
		&pricetypes.QueryAllMarketParamsResponse{
			MarketParams: marketParams_AddMarkets,
		},
		100,
	)

	s.startClient()

	// Expect prices for one market configuration first.
	s.expectPricesWithTimeout(
		expectedPrices1Market,
		defaultMarketParams,
		testPriceCacheExpirationDuration,
	)

	// Eventually, the daemon should update its market params and produce prices for the new markets.
	s.expectPricesWithTimeout(
		expectedPrices3Markets,
		marketParams_AddMarkets,
		30*time.Second,
	)
}

func (s *PriceDaemonIntegrationTestSuite) TestUpdateMarkets_AddMarketWithUSDTConversion() {
	// Start the daemon with a single market. Then, update the endpoint to return market params that have new markets.
	s.mockAllMarketParamsResponseNTimes(
		&pricetypes.QueryAllMarketParamsResponse{
			MarketParams: defaultMarketParams,
		},
		1,
	)
	s.mockAllMarketParamsResponseNTimes(
		&pricetypes.QueryAllMarketParamsResponse{
			MarketParams: marketParams_AddMarketsWithAdjustments,
		},
		100,
	)

	s.startClient()

	// Expect prices for one market configuration first.
	s.expectPricesWithTimeout(
		expectedPrices1Market,
		defaultMarketParams,
		testPriceCacheExpirationDuration,
	)

	// Eventually, the daemon should update its market params and produce prices for the new markets.
	s.expectPricesWithTimeout(
		expectedPrices3MarketsWithConversions,
		marketParams_AddMarkets,
		30*time.Second,
	)
}

// TestUpdateMarkets_PartialUpdates tests that the pricefeed daemon applies valid market params and discards invalid
// params whenever an update is partially valid.
func (s *PriceDaemonIntegrationTestSuite) TestUpdateMarkets_PartialUpdate() {
	// Start the daemon with a single market. Then, update the endpoint to return partially valid params.
	s.mockAllMarketParamsResponseNTimes(
		&pricetypes.QueryAllMarketParamsResponse{
			MarketParams: defaultMarketParams,
		},
		1,
	)
	s.mockAllMarketParamsResponseNTimes(
		&pricetypes.QueryAllMarketParamsResponse{
			MarketParams: marketParams_PartialUpdate,
		},
		100,
	)

	s.startClient()

	// The cache expires old prices after 10 seconds, so wait an extra 5 seconds to validate that
	// all prices sent from the daemon before the invalid config was seen are drained / expired
	// from the cache by the time we check the cache for the expected prices.
	time.Sleep(testPriceCacheExpirationDuration + 5*time.Second)

	s.expectPricesWithTimeout(
		// The 1st market param is invalid and should not have applied, so the BTC price does not change.
		// The 2nd market param is valid and should have applied, so the ETH price should be added.
		expectedPricesPartialUpdate,
		marketParams_PartialUpdate,
		1*time.Second,
	)
}
