//go:build integration_test

package cli_test

import (
	"fmt"
	"testing"

	tmcli "github.com/cometbft/cometbft/libs/cli"
	"github.com/cosmos/cosmos-sdk/client/flags"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/dydxprotocol/v4/testutil/constants"
	"github.com/dydxprotocol/v4/testutil/network"
	"github.com/dydxprotocol/v4/testutil/nullify"
	"github.com/dydxprotocol/v4/x/prices/client/cli"
	"github.com/dydxprotocol/v4/x/prices/types"
)

func networkWithMarketObjects(t *testing.T, n int) (*network.Network, []types.Market) {
	t.Helper()
	cfg := network.DefaultConfig(nil)
	state := types.GenesisState{}
	require.NoError(t, cfg.Codec.UnmarshalJSON(cfg.GenesisState[types.ModuleName], &state))

	// ExchangeFeeds
	for i := 0; i < 2; i++ {
		exchangeFeed := types.ExchangeFeed{
			Id:   uint32(i),
			Name: fmt.Sprint("ExchangeFeedName", i),
			Memo: fmt.Sprint("TestMemo", i),
		}
		nullify.Fill(&exchangeFeed) //nolint:staticcheck
		state.ExchangeFeeds = append(state.ExchangeFeeds, exchangeFeed)
	}

	// Markets
	for i := 0; i < n; i++ {
		market := types.Market{
			Id:                uint32(i),
			Pair:              fmt.Sprint(constants.BtcUsdPair, i),
			Exchanges:         []uint32{0, 1},
			MinExchanges:      uint32(1),
			MinPriceChangePpm: uint32((i + 1) * 2),
		}
		nullify.Fill(&market) //nolint:staticcheck
		state.Markets = append(state.Markets, market)
	}

	buf, err := cfg.Codec.MarshalJSON(&state)
	require.NoError(t, err)
	cfg.GenesisState[types.ModuleName] = buf
	return network.New(t, cfg), state.Markets
}

func TestShowMarket(t *testing.T) {
	net, objs := networkWithMarketObjects(t, 2)

	ctx := net.Validators[0].ClientCtx
	common := []string{
		fmt.Sprintf("--%s=json", tmcli.OutputFlag),
	}
	for _, tc := range []struct {
		desc string
		id   uint32

		args []string
		err  error
		obj  types.Market
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
			err:  status.Error(codes.NotFound, "not found"),
		},
	} {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			args := []string{
				fmt.Sprintf("%v", tc.id),
			}
			args = append(args, tc.args...)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdShowMarket(), args)
			if tc.err != nil {
				stat, ok := status.FromError(tc.err)
				require.True(t, ok)
				require.ErrorIs(t, stat.Err(), tc.err)
			} else {
				require.NoError(t, err)
				var resp types.QueryMarketResponse
				require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
				require.NotNil(t, resp.Market)
				require.Equal(t,
					nullify.Fill(&tc.obj),      //nolint:staticcheck
					nullify.Fill(&resp.Market), //nolint:staticcheck
				)
			}
		})
	}
}

func TestListMarket(t *testing.T) {
	net, objs := networkWithMarketObjects(t, 5)

	ctx := net.Validators[0].ClientCtx
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
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListMarket(), args)
			require.NoError(t, err)
			var resp types.QueryAllMarketsResponse
			require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
			require.LessOrEqual(t, len(resp.Market), step)
			require.Subset(t,
				nullify.Fill(objs),        //nolint:staticcheck
				nullify.Fill(resp.Market), //nolint:staticcheck
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(objs); i += step {
			args := request(next, 0, uint64(step), false)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListMarket(), args)
			require.NoError(t, err)
			var resp types.QueryAllMarketsResponse
			require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
			require.LessOrEqual(t, len(resp.Market), step)
			require.Subset(t,
				nullify.Fill(objs),        //nolint:staticcheck
				nullify.Fill(resp.Market), //nolint:staticcheck
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		args := request(nil, 0, uint64(len(objs)), true)
		out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListMarket(), args)
		require.NoError(t, err)
		var resp types.QueryAllMarketsResponse
		require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
		require.NoError(t, err)
		require.Equal(t, len(objs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(objs),        //nolint:staticcheck
			nullify.Fill(resp.Market), //nolint:staticcheck
		)
	})
}
