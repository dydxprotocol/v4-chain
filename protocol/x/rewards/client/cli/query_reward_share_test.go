package cli_test

import (
	"fmt"
	tmcli "github.com/cometbft/cometbft/libs/cli"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/x/rewards/client/cli"
	"github.com/dydxprotocol/v4-chain/protocol/x/rewards/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestQueryRewardShare(t *testing.T) {
	net, ctx := setupNetwork(t)

	out, err := clitestutil.ExecTestCLICmd(
		ctx,
		cli.CmdQueryRewardShare(),
		[]string{
			constants.AliceAccAddress.String(),
			fmt.Sprintf("--%s=json", tmcli.OutputFlag),
		},
	)

	require.NoError(t, err)
	var resp types.QueryRewardShareResponse
	require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
	require.Equal(t, types.RewardShare{
		Address: constants.AliceAccAddress.String(),
		Weight:  dtypes.NewInt(0),
	}, resp.RewardShare)
}
