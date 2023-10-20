//go:build all || integration_test

package cli_test

import (
	"fmt"
	tmcli "github.com/cometbft/cometbft/libs/cli"
	"github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	clobcli "github.com/dydxprotocol/v4-chain/protocol/x/clob/client/cli"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
)

func TestCmdListStatefulOrders(t *testing.T) {
	net, _ := networkWithClobPairObjects(t, 2)

	ctx := net.Validators[0].ClientCtx

	args := []string{
		fmt.Sprintf("--%s=json", tmcli.OutputFlag),
	}

	out, err := cli.ExecTestCLICmd(ctx, clobcli.CmdListStatefulOrders(), args)
	require.NoError(t, err)

	var res types.QueryAllStatefulOrdersResponse
	require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &res))

	// Expect no stateful orders.
	require.Equal(t, types.QueryAllStatefulOrdersResponse{
		StatefulOrders: []types.Order{},
	}, res)
}

func TestCmdGetStatefulOrderCount(t *testing.T) {
	net, _ := networkWithClobPairObjects(t, 2)
	ctx := net.Validators[0].ClientCtx

	for name, tc := range map[string]struct {
		owner       string
		number      uint32
		count       uint32
		expectedErr error
	}{
		"valid": {
			owner:  constants.AliceAccAddress.String(),
			number: 0,
			count:  0,
		},
		"invalid owner": {
			owner:  "invalid",
			number: 0,
			expectedErr: status.Error(
				codes.InvalidArgument,
				"Invalid owner address: decoding bech32 failed: invalid bech32 string length 7",
			),
		},
	} {
		t.Run(name, func(t *testing.T) {
			args := []string{
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
				tc.owner,
				fmt.Sprintf("%d", tc.number),
			}
			out, err := cli.ExecTestCLICmd(ctx, clobcli.CmdGetStatefulOrderCount(), args)
			if tc.expectedErr != nil {
				require.ErrorIs(t, err, tc.expectedErr)
			} else {
				require.NoError(t, err)

				var res types.QueryStatefulOrderCountResponse
				require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &res))

				// Expect no stateful orders.
				require.Equal(t, types.QueryStatefulOrderCountResponse{
					Count: tc.count,
				}, res)
			}
		})
	}
}
