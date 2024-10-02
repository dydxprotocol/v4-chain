package types

import (
	"testing"
	"time"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/sdaioracle/api"
	store "github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/sdaioracle/client/contract"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/sdaioracle/client/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
)

func TestGetInitialEvent(t *testing.T) {
	fetcher := &EthEventFetcher{}

	t.Run("Empty request", func(t *testing.T) {
		result, err := fetcher.GetInitialEvent(true)
		require.NoError(t, err)
		require.Equal(t, api.AddsDAIEventsRequest{}, result)
	})

	t.Run("Successful query", func(t *testing.T) {
		client, err := ethclient.Dial(types.ETHRPC)
		require.NoError(t, err, "Failed to connect to Ethereum client")
		defer client.Close()

		startTime := time.Now()
		result, err := fetcher.GetInitialEvent(false)
		endTime := time.Now()

		require.NoError(t, err, "Failed to get initial event")
		require.NotEmpty(t, result.ConversionRate, "Conversion rate should not be empty")

		directRate, err := store.QueryDaiConversionRate(client)
		require.NoError(t, err, "Failed to query DAI conversion rate directly")

		require.Equal(t, directRate, result.ConversionRate, "Rate from GetInitialEvent should match direct query")

		require.True(t, endTime.Sub(startTime) < time.Second, "Successful query should be quick")
	})

	t.Run("Invalid RPC endpoint", func(t *testing.T) {
		originalETHRPC := types.ETHRPC
		types.ETHRPC = "http://invalid-endpoint:8545"
		defer func() { types.ETHRPC = originalETHRPC }()

		startTime := time.Now()
		_, err := fetcher.GetInitialEvent(false)
		endTime := time.Now()

		require.Error(t, err, "Should fail to get initial event with invalid RPC endpoint")

		retryDuration := endTime.Sub(startTime)
		require.True(t, retryDuration >= 2*time.Second, "Retry logic should have introduced delays")
	})
}
