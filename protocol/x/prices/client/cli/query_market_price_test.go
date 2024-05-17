//go:build all || integration_test

package cli_test

import (
	"encoding/base64"
	"fmt"
	"strings"
	"testing"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/network"
	cli_util "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/prices/cli"

	tmcli "github.com/cometbft/cometbft/libs/cli"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/stretchr/testify/require"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices/types"
)

func TestShowMarketPrice(t *testing.T) {
	_, _, objs := cli_util.NetworkWithMarketObjects(t, 2)
	common := []string{
		fmt.Sprintf("--%s=json", tmcli.OutputFlag),
	}
	cfg := network.DefaultConfig(nil)
	for _, tc := range []struct {
		desc string
		id   uint32

		args []string
		err  string
		obj  types.MarketPrice
	}{
		{
			desc: "found",
			id:   objs[0].Id,

			args: common,
			obj:  objs[0],
		},
		{
			desc: "not found",
			id:   uint32(100000),

			args: common,
			err:  "not found",
		},
	} {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			args := []string{
				fmt.Sprintf("%v", tc.id),
			}
			args = append(args, tc.args...)
			query := "docker exec interchain-security-instance interchain-security-cd query prices show-market-price " + fmt.Sprintf("%d", tc.id) + " --node tcp://7.7.8.4:26658 -o json"
			data, stderrOutput, err := network.QueryCustomNetwork(query)

			if tc.err != "" {
				require.Error(t, err)
				require.Contains(t, stderrOutput, tc.err)
			} else {
				require.NoError(t, err)
				var resp types.QueryMarketPriceResponse
				require.NoError(t, cfg.Codec.UnmarshalJSON(data, &resp))
				require.NotNil(t, resp.MarketPrice)
				require.Equal(t, tc.obj, resp.MarketPrice)
			}
		})
	}
	network.CleanupCustomNetwork()
}

func TestListMarketPrice(t *testing.T) {
	_, _, objs := cli_util.NetworkWithMarketObjects(t, 5)

	cfg := network.DefaultConfig(nil)
	request := func(next []byte, offset, limit uint64, total bool) []string {
		args := []string{
			fmt.Sprintf("--%s=json", tmcli.OutputFlag),
		}
		if next == nil {
			args = append(args, fmt.Sprintf("--%s=%d", flags.FlagOffset, offset))
		} else {
			args = append(args, fmt.Sprintf("--%s=%s", flags.FlagPageKey, next))
		}
		args = append(args, fmt.Sprintf("--%s=%d", flags.FlagLimit, limit))
		if total {
			args = append(args, fmt.Sprintf("--%s", flags.FlagCountTotal))
		}
		return args
	}
	t.Run("ByOffset", func(t *testing.T) {
		step := 2
		for i := 0; i < len(objs); i += step {
			args := request(nil, uint64(i), uint64(step), false)
			argsString := strings.Join(args, " ")
			commandString := "docker exec interchain-security-instance interchain-security-cd query prices list-market-price --node tcp://7.7.8.4:26658 " + argsString
			data, _, err := network.QueryCustomNetwork(commandString)
			require.NoError(t, err)
			var resp types.QueryAllMarketPricesResponse
			require.NoError(t, cfg.Codec.UnmarshalJSON(data, &resp))
			require.LessOrEqual(t, len(resp.MarketPrices), step)
			require.Subset(t, objs, resp.MarketPrices)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(objs); i += step {
			var nextKeyStr string
			if next != nil {
				nextKeyStr = base64.StdEncoding.EncodeToString(next)
			}
			args := request([]byte(nextKeyStr), 0, uint64(step), false)
			argsString := strings.Join(args, " ")
			commandString := "docker exec interchain-security-instance interchain-security-cd query prices list-market-price --node tcp://7.7.8.4:26658 " + argsString
			data, _, err := network.QueryCustomNetwork(commandString)
			require.NoError(t, err)
			var resp types.QueryAllMarketPricesResponse
			require.NoError(t, cfg.Codec.UnmarshalJSON(data, &resp))
			require.LessOrEqual(t, len(resp.MarketPrices), step)
			require.Subset(t, objs, resp.MarketPrices)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		args := request(nil, 0, uint64(len(objs)), true)
		argsString := strings.Join(args, " ")
		commandString := "docker exec interchain-security-instance interchain-security-cd query prices list-market-price --node tcp://7.7.8.4:26658 " + argsString
		data, _, err := network.QueryCustomNetwork(commandString)
		require.NoError(t, err)
		var resp types.QueryAllMarketPricesResponse
		require.NoError(t, cfg.Codec.UnmarshalJSON(data, &resp))
		require.NoError(t, err)
		require.Equal(t, len(objs), int(resp.Pagination.Total))
		require.ElementsMatch(t, objs, resp.MarketPrices)
	})
	network.CleanupCustomNetwork()
}
