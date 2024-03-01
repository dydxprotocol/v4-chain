//go:build all || integration_test

package cli_test

import (
	"fmt"
	"strconv"
	"testing"

	tmcli "github.com/cometbft/cometbft/libs/cli"
	"github.com/cosmos/cosmos-sdk/client/flags"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/network"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/nullify"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/client/cli"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	perpetualstypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func networkWithClobPairObjects(t *testing.T, n int) (*network.Network, []types.ClobPair) {
	t.Helper()
	cfg := network.DefaultConfig(nil)

	// Init Prices genesis state.
	pricesState := constants.Prices_DefaultGenesisState
	pricesBuf, pricesErr := cfg.Codec.MarshalJSON(&pricesState)
	require.NoError(t, pricesErr)
	cfg.GenesisState[pricestypes.ModuleName] = pricesBuf

	// Init Perpetuals genesis state.
	// Add additional perps for objects exceeding the default perpetual count.
	// ClobPairs and Perpetuals should be one to one.
	perpetualsState := constants.Perpetuals_DefaultGenesisState
	for i := 2; i < n; i++ {
		perpetualsState.Perpetuals = append(
			perpetualsState.Perpetuals,
			perpetualstypes.Perpetual{
				Params: perpetualstypes.PerpetualParams{
					Id:            uint32(i),
					Ticker:        fmt.Sprintf("genesis_test_ticker_%d", i),
					LiquidityTier: 0,
					MarketType:    perpetualstypes.PerpetualMarketType_PERPETUAL_MARKET_TYPE_CROSS,
				},
				FundingIndex: dtypes.ZeroInt(),
			},
		)
	}
	perpetualsBuf, perpetualsErr := cfg.Codec.MarshalJSON(&perpetualsState)
	require.NoError(t, perpetualsErr)
	cfg.GenesisState[perpetualstypes.ModuleName] = perpetualsBuf

	// Init Clob State.
	state := types.GenesisState{}
	require.NoError(t, cfg.Codec.UnmarshalJSON(cfg.GenesisState[types.ModuleName], &state))

	for i := 0; i < n; i++ {
		clobPair := types.ClobPair{
			Id: uint32(i),
			Metadata: &types.ClobPair_PerpetualClobMetadata{
				PerpetualClobMetadata: &types.PerpetualClobMetadata{PerpetualId: uint32(i)},
			},
			SubticksPerTick:  5,
			StepBaseQuantums: 5,
			Status:           types.ClobPair_STATUS_ACTIVE,
		}
		nullify.Fill(&clobPair) //nolint:staticcheck
		state.ClobPairs = append(state.ClobPairs, clobPair)
	}
	buf, err := cfg.Codec.MarshalJSON(&state)
	require.NoError(t, err)
	cfg.GenesisState[types.ModuleName] = buf

	return network.New(t, cfg), state.ClobPairs
}

func TestShowClobPair(t *testing.T) {
	net, objs := networkWithClobPairObjects(t, 2)

	ctx := net.Validators[0].ClientCtx
	common := []string{
		fmt.Sprintf("--%s=json", tmcli.OutputFlag),
	}
	for _, tc := range []struct {
		desc string
		id   uint32

		args []string
		err  error
		obj  types.ClobPair
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
				lib.UintToString(tc.id),
			}
			args = append(args, tc.args...)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdShowClobPair(), args)
			if tc.err != nil {
				stat, ok := status.FromError(tc.err)
				require.True(t, ok)
				require.ErrorIs(t, stat.Err(), tc.err)
			} else {
				require.NoError(t, err)
				var resp types.QueryClobPairResponse
				require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
				require.NotNil(t, resp.ClobPair)
				require.Equal(t,
					nullify.Fill(&tc.obj),        //nolint:staticcheck
					nullify.Fill(&resp.ClobPair), //nolint:staticcheck
				)
			}
		})
	}
}

func TestListClobPair(t *testing.T) {
	net, objs := networkWithClobPairObjects(t, 5)

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
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListClobPair(), args)
			require.NoError(t, err)
			var resp types.QueryClobPairAllResponse
			require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
			require.LessOrEqual(t, len(resp.ClobPair), step)
			require.Subset(t,
				nullify.Fill(objs),          //nolint:staticcheck
				nullify.Fill(resp.ClobPair), //nolint:staticcheck
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(objs); i += step {
			args := request(next, 0, uint64(step), false)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListClobPair(), args)
			require.NoError(t, err)
			var resp types.QueryClobPairAllResponse
			require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
			require.LessOrEqual(t, len(resp.ClobPair), step)
			require.Subset(t,
				nullify.Fill(objs),          //nolint:staticcheck
				nullify.Fill(resp.ClobPair), //nolint:staticcheck
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		args := request(nil, 0, uint64(len(objs)), true)
		out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListClobPair(), args)
		require.NoError(t, err)
		var resp types.QueryClobPairAllResponse
		require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
		require.NoError(t, err)
		require.Equal(t, len(objs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(objs),          //nolint:staticcheck
			nullify.Fill(resp.ClobPair), //nolint:staticcheck
		)
	})
}
