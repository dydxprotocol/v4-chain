package store

import (
	"testing"
	"time"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/sdaioracle/client/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
)

func TestQueryDaiConversionRate_Integration(t *testing.T) {

	client, err := ethclient.Dial(types.ETHRPC)
	require.NoError(t, err, "Failed to connect to Ethereum client")
	defer client.Close()

	rate, err := QueryDaiConversionRate(client)
	require.NoError(t, err, "Failed to query DAI conversion rate")

	instance, err := NewStore(types.MakerContractAddress, client)
	require.NoError(t, err, "Failed to create contract instance")

	directRate, err := instance.Chi(&bind.CallOpts{})
	require.NoError(t, err, "Failed to query Chi directly from contract")

	require.Equal(t, directRate.String(), rate, "Rate from QueryDaiConversionRate should match direct contract query")
}

func TestQueryDaiConversionRateWithRetries_Integration(t *testing.T) {
	client, err := ethclient.Dial(types.ETHRPC)
	require.NoError(t, err, "Failed to connect to Ethereum client")
	defer client.Close()

	startTime := time.Now()
	rate, err := QueryDaiConversionRateWithRetries(client, 3)
	endTime := time.Now()
	require.NoError(t, err, "Failed to query DAI conversion rate with retries")
	require.NotEmpty(t, rate, "Rate should not be empty")

	instance, err := NewStore(types.MakerContractAddress, client)
	require.NoError(t, err, "Failed to create contract instance")

	directRate, err := instance.Chi(&bind.CallOpts{})
	require.NoError(t, err, "Failed to query Chi directly from contract")

	require.Equal(t, directRate.String(), rate, "Rate from QueryDaiConversionRateWithRetries should match direct contract query")

	require.True(t, endTime.Sub(startTime) < time.Second, "Successful query should be quick")

	// Test retry logic by temporarily using an invalid RPC endpoint
	invalidClient, err := ethclient.Dial("http://invalid-endpoint:8545")
	require.NoError(t, err, "Failed to create invalid client")
	defer invalidClient.Close()

	startTime = time.Now()
	_, err = QueryDaiConversionRateWithRetries(invalidClient, 3)
	endTime = time.Now()

	require.Error(t, err, "Should fail to query with invalid client")
	require.Equal(t, "failed to query DAI conversion rate", err.Error())

	retryDuration := endTime.Sub(startTime)
	require.True(t, retryDuration >= 2*time.Second, "Retry logic should have introduced delays")
}
