package client_test

import (
	"context"
	"fmt"
	"net"
	"sync"
	"testing"
	"time"

	"cosmossdk.io/log"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"

	"github.com/skip-mev/slinky/service/servers/oracle/types"

	appflags "github.com/dydxprotocol/v4-chain/protocol/app/flags"
	daemonflags "github.com/dydxprotocol/v4-chain/protocol/daemons/flags"
	pricefeed_types "github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/types"
	daemonserver "github.com/dydxprotocol/v4-chain/protocol/daemons/server"
	pricefeedserver_types "github.com/dydxprotocol/v4-chain/protocol/daemons/server/types/pricefeed"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/slinky/client"
	daemontypes "github.com/dydxprotocol/v4-chain/protocol/daemons/types"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/appoptions"
)

func TestPriceFetcherTestSuite(t *testing.T) {
	suite.Run(t, &PriceFetcherTestSuite{})
}

type PriceFetcherTestSuite struct {
	suite.Suite
	daemonFlags      daemonflags.DaemonFlags
	appFlags         appflags.Flags
	daemonServer     *daemonserver.Server
	pricesGrpcServer *grpc.Server
	wg               sync.WaitGroup
}

func (p *PriceFetcherTestSuite) SetupTest() {
	// Setup daemon and grpc servers.
	p.daemonFlags = daemonflags.GetDefaultDaemonFlags()
	p.appFlags = appflags.GetFlagValuesFromOptions(appoptions.GetDefaultTestAppOptions("", nil))

	// Configure and run daemon server.
	p.daemonServer = daemonserver.NewServer(
		log.NewNopLogger(),
		grpc.NewServer(),
		&daemontypes.FileHandlerImpl{},
		p.daemonFlags.Shared.SocketAddress,
	)
	p.daemonServer.WithPriceFeedMarketToExchangePrices(
		pricefeedserver_types.NewMarketToExchangePrices(5 * time.Second),
	)

	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		p.daemonServer.Start()
	}()

	// Create a gRPC server running on the default port and attach the mock prices query response.
	p.pricesGrpcServer = grpc.NewServer()

	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		ls, err := net.Listen("tcp", p.appFlags.GrpcAddress)
		p.Require().NoError(err)
		_ = p.pricesGrpcServer.Serve(ls)
	}()
}

func (p *PriceFetcherTestSuite) TearDownTest() {
	p.daemonServer.Stop()
	p.pricesGrpcServer.Stop()
	p.wg.Wait()
}

func (p *PriceFetcherTestSuite) TestPriceFetcher() {
	logger := log.NewTestLogger(p.T())
	mpf := mocks.NewMarketPairFetcher(p.T())
	slinky := mocks.NewOracleClient(p.T())
	slinky.On("Stop").Return(nil)
	var fetcher client.PriceFetcher

	p.Run("fetches prices on valid inputs", func() {
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
		p.Require().NoError(fetcher.Start(context.Background()))
		p.Require().NoError(fetcher.FetchPrices(context.Background()))
		fetcher.Stop()
	})

	p.Run("errors on slinky.Prices failure", func() {
		slinky.On("Start", mock.Anything).Return(nil).Once()
		slinky.On("Prices", mock.Anything, mock.Anything).
			Return(&types.QueryPricesResponse{}, fmt.Errorf("foobar")).Once()
		fetcher = client.NewPriceFetcher(
			mpf,
			pricefeedserver_types.NewMarketToExchangePrices(pricefeed_types.MaxPriceAge),
			slinky,
			logger,
		)

		p.Require().NoError(fetcher.Start(context.Background()))
		p.Require().Errorf(fetcher.FetchPrices(context.Background()), "foobar")
		fetcher.Stop()
	})

	p.Run("errors on slinky.Prices returning invalid currency pairs", func() {
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

		p.Require().NoError(fetcher.Start(context.Background()))
		p.Require().Errorf(fetcher.FetchPrices(context.Background()), "incorrectly formatted CurrencyPair")
		fetcher.Stop()
	})

	p.Run("no-ops on marketPairFetcher currency pair not found", func() {
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
		p.Require().NoError(fetcher.Start(context.Background()))
		p.Require().NoError(fetcher.FetchPrices(context.Background()))
		fetcher.Stop()
	})

	p.Run("continues on non-parsable price data", func() {
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
		p.Require().NoError(fetcher.Start(context.Background()))
		p.Require().NoError(fetcher.FetchPrices(context.Background()))
		fetcher.Stop()
	})

	p.Run("no-ops on empty price response", func() {
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
		p.Require().NoError(fetcher.Start(context.Background()))
		p.Require().NoError(fetcher.FetchPrices(context.Background()))
		fetcher.Stop()
	})
}
