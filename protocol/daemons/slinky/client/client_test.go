package client_test

import (
	"context"
	"net"
	"sync"
	"testing"
	"time"

	"cosmossdk.io/log"
	"github.com/skip-mev/slinky/service/servers/oracle/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"

	appflags "github.com/dydxprotocol/v4-chain/protocol/app/flags"
	daemonflags "github.com/dydxprotocol/v4-chain/protocol/daemons/flags"
	pricefeedtypes "github.com/dydxprotocol/v4-chain/protocol/daemons/server/types/pricefeed"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/slinky/client"
	daemontypes "github.com/dydxprotocol/v4-chain/protocol/daemons/types"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/appoptions"
	pricetypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
)

func TestClientTestSuite(t *testing.T) {
	suite.Run(t, &ClientTestSuite{})
}

type ClientTestSuite struct {
	suite.Suite
	daemonFlags           daemonflags.DaemonFlags
	appFlags              appflags.Flags
	grpcServer            *grpc.Server
	pricesMockQueryServer *mocks.QueryServer
	wg                    sync.WaitGroup
}

func (c *ClientTestSuite) SetupTest() {
	// Setup grpc server.
	c.daemonFlags = daemonflags.GetDefaultDaemonFlags()
	c.appFlags = appflags.GetFlagValuesFromOptions(appoptions.GetDefaultTestAppOptions("", nil))
	c.grpcServer = grpc.NewServer()

	c.pricesMockQueryServer = &mocks.QueryServer{}
	pricetypes.RegisterQueryServer(c.grpcServer, c.pricesMockQueryServer)

	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		ls, err := net.Listen("tcp", c.appFlags.GrpcAddress)
		c.Require().NoError(err)
		_ = c.grpcServer.Serve(ls)
	}()
}

func (c *ClientTestSuite) TearDownTest() {
	c.grpcServer.Stop()
	c.wg.Wait()
}

func (c *ClientTestSuite) TestClient() {
	var cli *client.Client
	slinky := mocks.NewOracleClient(c.T())
	logger := log.NewTestLogger(c.T())

	c.pricesMockQueryServer.On("AllMarketParams", mock.Anything, mock.Anything).
		Return(
			&pricetypes.QueryAllMarketParamsResponse{
				MarketParams: []pricetypes.MarketParam{
					{Id: 0, Pair: "FOO-BAR"},
					{Id: 1, Pair: "BAR-FOO"},
				}},
			nil,
		)

	c.Run("services are all started and call their deps", func() {
		slinky.On("Stop").Return(nil)
		slinky.On("Start", mock.Anything).Return(nil).Once()
		slinky.On("Prices", mock.Anything, mock.Anything).
			Return(&types.QueryPricesResponse{
				Prices: map[string]string{
					"FOO/BAR": "100000000000",
				},
				Timestamp: time.Now(),
			}, nil)
		client.SlinkyPriceFetchDelay = time.Millisecond
		client.SlinkyMarketParamFetchDelay = time.Millisecond
		cli = client.StartNewClient(
			context.Background(),
			slinky,
			pricefeedtypes.NewMarketToExchangePrices(5*time.Second),
			&daemontypes.GrpcClientImpl{},
			c.daemonFlags,
			c.appFlags,
			logger,
		)
		waitTime := time.Second * 5
		c.Require().Eventually(func() bool {
			return cli.HealthCheck() == nil
		}, waitTime, time.Millisecond*500, "Slinky daemon failed to become healthy within %s", waitTime)
		// Need to wait until a single c
		cli.Stop()
		c.Require().NoError(cli.HealthCheck())
	})
}
