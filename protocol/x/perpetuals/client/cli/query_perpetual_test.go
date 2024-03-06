//go:build all || integration_test

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

	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/network"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/nullify"
	"github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/client/cli"
	"github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

var (
	serializableIntCmpOpts = cmp.Options{
		cmpopts.IgnoreFields(types.Perpetual{}, "FundingIndex"), // existing ignore option
		cmp.FilterValues(
			func(x, y dtypes.SerializableInt) bool {
				// This will apply the custom comparer only to SerializableInt fields
				return true // Apply this filter to all SerializableInt values
			},
			cmp.Comparer(func(x, y dtypes.SerializableInt) bool { return x.Cmp(y) == 0 }),
		),
	}
)

func networkWithLiquidityTierAndPerpetualObjects(
	t *testing.T,
	m int,
	n int,
) (
	*network.Network,
	[]types.LiquidityTier,
	[]types.Perpetual,
) {
	t.Helper()
	cfg := network.DefaultConfig(nil)

	// Init Prices state.
	pricesState := constants.Prices_DefaultGenesisState
	pricesBuf, pricesErr := cfg.Codec.MarshalJSON(&pricesState)
	require.NoError(t, pricesErr)
	cfg.GenesisState[pricestypes.ModuleName] = pricesBuf

	// Init Perpetuals state.
	state := types.GenesisState{}
	require.NoError(t, cfg.Codec.UnmarshalJSON(cfg.GenesisState[types.ModuleName], &state))

	// Generate `m` Liquidity Tiers.
	for i := 0; i < m; i++ {
		liquidityTier := types.LiquidityTier{
			Id:                     uint32(i),
			Name:                   fmt.Sprintf("test_liquidity_tier_name_%d", i),
			InitialMarginPpm:       uint32(1_000_000 / (i + 1)),
			MaintenanceFractionPpm: uint32(1_000_000 / (i + 1)),
			ImpactNotional:         uint64(500_000_000 * (i + 1)),
		}
		nullify.Fill(&liquidityTier) //nolint:staticcheck
		state.LiquidityTiers = append(state.LiquidityTiers, liquidityTier)
	}

	// Generate `n` Perpetuals.

	for i := 0; i < n; i++ {
		marketType := types.PerpetualMarketType_PERPETUAL_MARKET_TYPE_CROSS
		if i%2 == 1 {
			marketType = types.PerpetualMarketType_PERPETUAL_MARKET_TYPE_ISOLATED
		}

		perpetual := types.Perpetual{
			Params: types.PerpetualParams{
				Id:            uint32(i),
				Ticker:        fmt.Sprintf("test_query_ticker_%d", i),
				LiquidityTier: uint32(i % m),
				MarketType:    marketType,
			},
			FundingIndex: dtypes.ZeroInt(),
			OpenInterest: dtypes.ZeroInt(),
		}
		nullify.Fill(&perpetual) //nolint:staticcheck
		state.Perpetuals = append(state.Perpetuals, perpetual)
	}
	buf, err := cfg.Codec.MarshalJSON(&state)
	require.NoError(t, err)
	cfg.GenesisState[types.ModuleName] = buf

	return network.New(t, cfg), state.LiquidityTiers, state.Perpetuals
}

func TestShowPerpetual(t *testing.T) {
	net, _, objs := networkWithLiquidityTierAndPerpetualObjects(t, 2, 2)

	ctx := net.Validators[0].ClientCtx
	common := []string{
		fmt.Sprintf("--%s=json", tmcli.OutputFlag),
	}
	for _, tc := range []struct {
		desc string
		id   uint32

		args []string
		err  error
		obj  types.Perpetual
	}{
		{
			desc: "found",
			id:   objs[0].Params.Id,

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
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdShowPerpetual(), args)
			if tc.err != nil {
				stat, ok := status.FromError(tc.err)
				require.True(t, ok)
				require.ErrorIs(t, stat.Err(), tc.err)
			} else {
				require.NoError(t, err)
				var resp types.QueryPerpetualResponse
				require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
				require.NotNil(t, resp.Perpetual)
				checkExpectedPerp(t, tc.obj, resp.Perpetual)
			}
		})
	}
}

// Check the received perpetual matches with expected.
// FundingIndex field is ignored since it can vary depending on funding-tick epoch.
// TODO(DEC-606): Improve end-to-end testing related to ticking epochs.
func checkExpectedPerp(t *testing.T, expected types.Perpetual, received types.Perpetual) {
	if diff := cmp.Diff(expected, received, serializableIntCmpOpts...); diff != "" {
		t.Errorf("resp.Perpetual mismatch (-want +received):\n%s", diff)
	}
}

// Check the received perpetual object matches one of the expected perpetuals.
func expectedContainsReceived(t *testing.T, expectedPerps []types.Perpetual, received types.Perpetual) {
	for _, expected := range expectedPerps {
		if received.Params.Id == expected.Params.Id {
			checkExpectedPerp(t, expected, received)
			return
		}
	}
	t.Errorf("Received perp (%v) not found in expected perps (%v)", received, expectedPerps)
}

func TestListPerpetual(t *testing.T) {
	net, _, objs := networkWithLiquidityTierAndPerpetualObjects(t, 3, 5)

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
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListPerpetual(), args)
			require.NoError(t, err)
			var resp types.QueryAllPerpetualsResponse
			require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
			require.LessOrEqual(t, len(resp.Perpetual), step)
			for _, perp := range resp.Perpetual {
				expectedContainsReceived(t, objs, perp)
			}
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(objs); i += step {
			args := request(next, 0, uint64(step), false)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListPerpetual(), args)
			require.NoError(t, err)
			var resp types.QueryAllPerpetualsResponse
			require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
			require.LessOrEqual(t, len(resp.Perpetual), step)
			for _, perp := range resp.Perpetual {
				expectedContainsReceived(t, objs, perp)
			}
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		args := request(nil, 0, uint64(len(objs)), true)
		out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListPerpetual(), args)
		require.NoError(t, err)
		var resp types.QueryAllPerpetualsResponse
		require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
		require.NoError(t, err)
		require.Equal(t, len(objs), int(resp.Pagination.Total))
		cmpOptions := append(
			serializableIntCmpOpts,
			cmpopts.SortSlices(func(x, y types.Perpetual) bool {
				return x.Params.Id > y.Params.Id
			}),
		)

		if diff := cmp.Diff(objs, resp.Perpetual, cmpOptions...); diff != "" {
			t.Errorf("resp.Perpetual mismatch (-want +received):\n%s", diff)
		}
	})
}
