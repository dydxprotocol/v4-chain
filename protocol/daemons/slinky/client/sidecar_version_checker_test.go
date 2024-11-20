package client_test

import (
	"context"
	"testing"

	"cosmossdk.io/log"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/dydxprotocol/v4-chain/protocol/daemons/slinky/client"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	"github.com/skip-mev/connect/v2/service/servers/oracle/types"
)

func TestSidecarVersionChecker(t *testing.T) {
	logger := log.NewTestLogger(t)
	var fetcher client.SidecarVersionChecker

	t.Run("Checks sidecar version passes", func(t *testing.T) {
		slinky := mocks.NewOracleClient(t)
		slinky.On("Stop").Return(nil)
		slinky.On("Start", mock.Anything).Return(nil).Once()
		slinky.On("Version", mock.Anything, mock.Anything).
			Return(&types.QueryVersionResponse{
				Version: client.MinSidecarVersion,
			}, nil)
		fetcher = client.NewSidecarVersionChecker(slinky, logger)
		require.NoError(t, fetcher.Start(context.Background()))
		require.NoError(t, fetcher.CheckSidecarVersion(context.Background()))
		fetcher.Stop()
	})

	t.Run("Checks sidecar version less than minimum version", func(t *testing.T) {
		slinky := mocks.NewOracleClient(t)
		slinky.On("Stop").Return(nil)
		slinky.On("Start", mock.Anything).Return(nil).Once()
		slinky.On("Version", mock.Anything, mock.Anything).
			Return(&types.QueryVersionResponse{
				Version: "v0.0.0",
			}, nil)
		fetcher = client.NewSidecarVersionChecker(slinky, logger)
		require.NoError(t, fetcher.Start(context.Background()))
		require.ErrorContains(t, fetcher.CheckSidecarVersion(context.Background()),
			"Sidecar version 0.0.0 is less than minimum required version")
		fetcher.Stop()
	})

	t.Run("Checks sidecar version incorrectly formatted", func(t *testing.T) {
		slinky := mocks.NewOracleClient(t)
		slinky.On("Stop").Return(nil)
		slinky.On("Start", mock.Anything).Return(nil).Once()
		slinky.On("Version", mock.Anything, mock.Anything).
			Return(&types.QueryVersionResponse{
				Version: "a.b.c",
			}, nil)
		fetcher = client.NewSidecarVersionChecker(slinky, logger)
		require.NoError(t, fetcher.Start(context.Background()))
		require.ErrorContains(t, fetcher.CheckSidecarVersion(context.Background()), "Malformed version: a.b.c")
		fetcher.Stop()
	})
}
