//go:build all || integration_test

package cli_test

import (
	"fmt"
	"strconv"
	"testing"
	"time"

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

	time.Sleep(15 * time.Second) // to ensure other validators have queried the sdai rate at this block

	setTx := "docker exec interchain-security-instance-setup interchain-security-cd" +
		" tx ratelimit update-market-prices dydx1eeeggku6dzk3mv7wph3zq035rhtd890smfq5z6 " +
		chi + " " + blockNumber +
		" --chain-id consu --home /consu/validatoralice --keyring-backend test -y"
	_, _, err = network.QueryCustomNetwork(setTx)
	require.NoError(t, err)

	time.Sleep(10 * time.Second)

	param := fmt.Sprintf("--%s=json", tmcli.OutputFlag)

	rateQuery := "docker exec interchain-security-instance-setup interchain-security-cd" +
		" query ratelimit get-sdai-price " + param
	data, _, err := network.QueryCustomNetwork(rateQuery)

	require.NoError(t, err)
	var resp types.GetSDAIPriceQueryResponse
	require.NoError(t, cfg.Codec.UnmarshalJSON(data, &resp))
	require.Equal(t, chi, resp.Price)
}
