//go:build all || integration_test

package cli_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/network"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/ratelimit/types"
	tmcli "github.com/cometbft/cometbft/libs/cli"
	"github.com/stretchr/testify/require"

	"testing"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/assert"

	store "github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/sDAIOracle/client/contract"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/sDAIOracle/client/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func TestGetSDAIPriceQuery(t *testing.T) {
	cfg := network.DefaultConfig(nil)

	// Test with real client
	client, err := ethclient.Dial(types.ETHRPC)
	if err != nil {
		t.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	chi, blockNumber, err := store.QueryDaiConversionRate(client)
	assert.Nil(t, err, "Expected no error with real client")

	// todo solal set the sDAI price to chi for block number

	param := fmt.Sprintf("--%s=json", tmcli.OutputFlag)

	rateQuery := "docker exec interchain-security-instance-setup interchain-security-cd" +
		" query ratelimit get-sdai-price " + param
	data, _, err := network.QueryCustomNetwork(rateQuery)

	require.NoError(t, err)
	var resp types.GetSDAIPriceQueryResponse
	require.NoError(t, cfg.Codec.UnmarshalJSON(data, &resp))
	require.Equal(t, chi, resp.Price)
}
