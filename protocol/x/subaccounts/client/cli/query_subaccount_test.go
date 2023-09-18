//go:build all || integration_test

package cli_test

import (
	"fmt"
	"github.com/dydxprotocol/v4-chain/protocol/app/stoppable"
	"strconv"
	"testing"

	tmcli "github.com/cometbft/cometbft/libs/cli"
	"github.com/cosmos/cosmos-sdk/client/flags"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/dydxprotocol/v4-chain/protocol/testutil/network"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/nullify"
	"github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/client/cli"
	"github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func networkWithSubaccountObjects(t *testing.T, n int) (*network.Network, []types.Subaccount) {
	t.Helper()
	cfg := network.DefaultConfig(nil)
	state := types.GenesisState{}
	require.NoError(t, cfg.Codec.UnmarshalJSON(cfg.GenesisState[types.ModuleName], &state))

	for i := 0; i < n; i++ {
		subaccount := types.Subaccount{
			Id: &types.SubaccountId{
				Owner:  strconv.Itoa(i),
				Number: uint32(n),
			},
		}
		nullify.Fill(&subaccount) //nolint:staticcheck
		state.Subaccounts = append(state.Subaccounts, subaccount)
	}
	buf, err := cfg.Codec.MarshalJSON(&state)
	require.NoError(t, err)
	cfg.GenesisState[types.ModuleName] = buf

	t.Cleanup(func() {
		stoppable.StopServices(t, cfg.GRPCAddress)
	})

	return network.New(t, cfg), state.Subaccounts
}

func TestShowSubaccount(t *testing.T) {
	net, objs := networkWithSubaccountObjects(t, 2)

	ctx := net.Validators[0].ClientCtx
	common := []string{
		fmt.Sprintf("--%s=json", tmcli.OutputFlag),
	}
	for _, tc := range []struct {
		desc   string
		owner  string
		number uint32

		args []string
		err  error
		obj  types.Subaccount
	}{
		{
			desc:   "found",
			owner:  objs[0].Id.Owner,
			number: objs[0].Id.Number,
			args:   common,
			obj:    objs[0],
		},
		{
			desc:   "not found owner",
			owner:  "abdefg",
			number: objs[0].Id.Number,

			args: common,
			err:  status.Error(codes.NotFound, "not found"),
		},
		{
			desc:   "not found number",
			owner:  objs[0].Id.Owner,
			number: uint32(0),

			args: common,
			err:  status.Error(codes.NotFound, "not found"),
		},
	} {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			args := []string{
				tc.owner,
				strconv.Itoa(int(tc.number)),
			}
			args = append(args, tc.args...)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdShowSubaccount(), args)
			if tc.err != nil {
				stat, ok := status.FromError(tc.err)
				require.True(t, ok)
				require.ErrorIs(t, stat.Err(), tc.err)
			} else {
				require.NoError(t, err)
				var resp types.QuerySubaccountResponse
				require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
				require.NotNil(t, resp.Subaccount)
				require.Equal(t,
					nullify.Fill(&tc.obj),          //nolint:staticcheck
					nullify.Fill(&resp.Subaccount), //nolint:staticcheck
				)
			}
		})
	}
}

func TestListSubaccount(t *testing.T) {
	net, objs := networkWithSubaccountObjects(t, 5)

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
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListSubaccount(), args)
			require.NoError(t, err)
			var resp types.QuerySubaccountAllResponse
			require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
			require.LessOrEqual(t, len(resp.Subaccount), step)
			require.Subset(t,
				nullify.Fill(objs),            //nolint:staticcheck
				nullify.Fill(resp.Subaccount), //nolint:staticcheck
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(objs); i += step {
			args := request(next, 0, uint64(step), false)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListSubaccount(), args)
			require.NoError(t, err)
			var resp types.QuerySubaccountAllResponse
			require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
			require.LessOrEqual(t, len(resp.Subaccount), step)
			require.Subset(t,
				nullify.Fill(objs),            //nolint:staticcheck
				nullify.Fill(resp.Subaccount), //nolint:staticcheck
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		args := request(nil, 0, uint64(len(objs)), true)
		out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListSubaccount(), args)
		require.NoError(t, err)
		var resp types.QuerySubaccountAllResponse
		require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
		require.NoError(t, err)
		require.Equal(t, len(objs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(objs),            //nolint:staticcheck
			nullify.Fill(resp.Subaccount), //nolint:staticcheck
		)
	})
}
