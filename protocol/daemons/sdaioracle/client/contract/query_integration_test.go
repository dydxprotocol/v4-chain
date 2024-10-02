package store

import (
	"testing"

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
