//go:build all || integration_test

package cli_test

import (
	"encoding/base64"
	"fmt"
	"strings"
	"testing"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/network"
	cli_util "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/prices/cli"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/stretchr/testify/require"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices/types"
)

func TestShowMarketParam(t *testing.T) {
	_, objs, _ := cli_util.NetworkWithMarketObjects(t, 2)

	for _, tc := range []struct {
		desc string
		id   uint32

		err string
		obj types.MarketParam
	}{
		{
			desc: "found",
			id:   objs[0].Id,

			obj: objs[0],
		},
		{
			desc: "not found",
			id:   uint32(100000),

			err: "not found",
		},
	} {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			cfg := network.DefaultConfig(nil)
			query := "docker exec interchain-security-instance interchain-security-cd query prices show-market-param " + fmt.Sprintf("%d", tc.id)
			data, stderrOutput, err := network.QueryCustomNetwork(query)
			if tc.err != "" {
				require.Error(t, err)
				require.Contains(t, stderrOutput, tc.err)
			} else {
				require.NoError(t, err)
				var resp types.QueryMarketParamResponse
				require.NoError(t, cfg.Codec.UnmarshalJSON(data, &resp))
				require.NotNil(t, resp.MarketParam)
				require.Equal(t,
					&tc.obj,
					&resp.MarketParam,
				)
			}
		})
	}
	network.CleanupCustomNetwork()
}

func TestListMarketParam(t *testing.T) {
	_, objs, _ := cli_util.NetworkWithMarketObjects(t, 5)
	cfg := network.DefaultConfig(nil)
	request := func(next []byte, offset, limit uint64, total bool) []string {
		args := []string{}
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
			commandString := "docker exec interchain-security-instance interchain-security-cd query prices list-market-param " + argsString
			data, _, err := network.QueryCustomNetwork(commandString)

			require.NoError(t, err)
			var resp types.QueryAllMarketParamsResponse
			require.NoError(t, cfg.Codec.UnmarshalJSON(data, &resp))
			require.LessOrEqual(t, len(resp.MarketParams), step)
			require.Subset(t, objs, resp.MarketParams)
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
			commandString := "docker exec interchain-security-instance interchain-security-cd query prices list-market-param " + argsString
			data, _, err := network.QueryCustomNetwork(commandString)

			require.NoError(t, err)
			var resp types.QueryAllMarketParamsResponse
			require.NoError(t, cfg.Codec.UnmarshalJSON(data, &resp))
			require.LessOrEqual(t, len(resp.MarketParams), step)
			require.Subset(t, objs, resp.MarketParams)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		args := request(nil, 0, uint64(len(objs)), true)

		argsString := strings.Join(args, " ")
		commandString := "docker exec interchain-security-instance interchain-security-cd query prices list-market-param " + argsString
		data, _, err := network.QueryCustomNetwork(commandString)
		require.NoError(t, err)
		var resp types.QueryAllMarketParamsResponse
		require.NoError(t, cfg.Codec.UnmarshalJSON(data, &resp))
		require.NoError(t, err)
		require.Equal(t, len(objs), int(resp.Pagination.Total))
		require.ElementsMatch(t, objs, resp.MarketParams)
	})
	network.CleanupCustomNetwork()
}
