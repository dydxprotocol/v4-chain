package client_test

import (
	"context"
	"net"
	"testing"
	"time"

	"cosmossdk.io/log"
	"github.com/dydxprotocol/slinky/service/servers/oracle/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
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

func TestClient(t *testing.T) {
	var cli *client.Client
	slinky := mocks.NewOracleClient(t)
	logger := log.NewTestLogger(t)

	daemonFlags := daemonflags.GetDefaultDaemonFlags()
	appFlags := appflags.GetFlagValuesFromOptions(appoptions.GetDefaultTestAppOptions("", nil))

	grpcServer := grpc.NewServer()
	pricesMockQueryServer := &mocks.QueryServer{}
	pricetypes.RegisterQueryServer(grpcServer, pricesMockQueryServer)
	pricesMockQueryServer.On("AllMarketParams", mock.Anything, mock.Anything).
		Return(
			&pricetypes.QueryAllMarketParamsResponse{
				MarketParams: []pricetypes.MarketParam{
					{Id: 0, Pair: "FOO-BAR"},
					{Id: 1, Pair: "BAR-FOO"},
				}},
			nil,
		)

	defer grpcServer.Stop()
	go func() {
		ls, err := net.Listen("tcp", appFlags.GrpcAddress)
		require.NoError(t, err)
		err = grpcServer.Serve(ls)
		require.NoError(t, err)
	}()

	slinky.On("Stop").Return(nil)
	slinky.On("Start", mock.Anything).Return(nil).Twice()
	slinky.On("Prices", mock.Anything, mock.Anything).
		Return(&types.QueryPricesResponse{
			Prices: map[string]string{
				"FOO/BAR": "100000000000",
			},
			Timestamp: time.Now(),
		}, nil)
	slinky.On("Version", mock.Anything, mock.Anything).
		Return(&types.QueryVersionResponse{
			Version: client.MinSidecarVersion,
		}, nil)

	client.SlinkyPriceFetchDelay = time.Millisecond
	client.SlinkyMarketParamFetchDelay = time.Millisecond
	client.SlinkySidecarCheckDelay = time.Millisecond
	cli = client.StartNewClient(
		context.Background(),
		slinky,
		pricefeedtypes.NewMarketToExchangePrices(5*time.Second),
		&daemontypes.GrpcClientImpl{},
		daemonFlags,
		appFlags,
		logger,
	)
	waitTime := time.Second * 5
	require.Eventually(t, func() bool {
		return cli.GetMarketPairHC().HealthCheck() == nil &&
			cli.GetPriceHC().HealthCheck() == nil &&
			cli.GetSidecarVersionHC().HealthCheck() == nil
	}, waitTime, time.Millisecond*500, "Slinky daemon failed to become healthy within %s", waitTime)
	cli.Stop()
}
