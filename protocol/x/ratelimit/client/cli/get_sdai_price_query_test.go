//go:build all || integration_test

package cli_test

import (
	"math/big"
	"strconv"
	"testing"
	"time"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/network"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/ratelimit/types"
	"github.com/stretchr/testify/require"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/assert"

	store "github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/sdaioracle/client/contract"
	sdaiOracleTypes "github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/sdaioracle/client/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func TestGetSDAIPriceQuery(t *testing.T) {
	cfg := network.DefaultConfig(nil)

	// Test with real client
	client, err := ethclient.Dial(sdaiOracleTypes.ETHRPC)
	if err != nil {
		t.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	chi, err := store.QueryDaiConversionRate(client)
	assert.Nil(t, err, "Expected no error with real client")

	time.Sleep(15 * time.Second)

	rateQuery := "docker exec interchain-security-instance-setup interchain-security-cd" +
		" query ratelimit get-sdai-price "
	data, _, err := network.QueryCustomNetwork(rateQuery)

	require.NoError(t, err)
	var resp types.GetSDAIPriceQueryResponse
	require.NoError(t, cfg.Codec.UnmarshalJSON(data, &resp))

	chiFloat, success := new(big.Float).SetString(chi)
	require.True(t, success, "Failed to parse chi as big.Float")

	priceFloat, success := new(big.Float).SetString(resp.Price)
	require.True(t, success, "Failed to parse price as big.Float")

	// Compare the big.Float values directly
	comparison := new(big.Float).Quo(priceFloat, chiFloat)

	minThreshold := big.NewFloat(0.99)
	maxThreshold := big.NewFloat(1.16)

	require.True(t, comparison.Cmp(minThreshold) >= 0, "Price should be at least 99% of chi")
	require.True(t, comparison.Cmp(maxThreshold) <= 0, "Price should be at most 116% of chi")
}
