package testutil

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	clobcli "github.com/dydxprotocol/v4-chain/protocol/x/clob/client/cli"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/testutil"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
)

// MsgPlaceOrderExec broadcasts a place order message.
func MsgPlaceOrderExec(
	clientCtx client.Context,
	owner sdk.AccAddress,
	number uint32,
	clientId uint64,
	clobPairId uint32,
	side types.Order_Side,
	quantums satypes.BaseQuantums,
	subticks uint64,
	goodTilBlock uint32,
) (testutil.BufferWriter, error) {
	sideNum := 1
	if side == types.Order_SIDE_SELL {
		sideNum = 2
	}
	args := []string{
		owner.String(),
		fmt.Sprint(number),
		fmt.Sprint(clientId),
		fmt.Sprint(clobPairId),
		fmt.Sprint(sideNum),
		fmt.Sprint(quantums),
		fmt.Sprint(subticks),
		fmt.Sprint(goodTilBlock),
	}

	args = append(args,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, "node0"),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
	)

	return clitestutil.ExecTestCLICmd(clientCtx, clobcli.CmdPlaceOrder(), args)
}

// MsgCancelOrderExec broadcasts a cancel order message.
func MsgCancelOrderExec(
	clientCtx client.Context,
	owner sdk.AccAddress,
	number uint32,
	clientId uint64,
	clobPairId uint32,
	goodTilBlock uint32,
) (testutil.BufferWriter, error) {
	args := []string{
		owner.String(),
		fmt.Sprint(number),
		fmt.Sprint(clientId),
		fmt.Sprint(clobPairId),
		fmt.Sprint(goodTilBlock),
	}

	args = append(args,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, "node0"),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
	)

	return clitestutil.ExecTestCLICmd(clientCtx, clobcli.CmdCancelOrder(), args)
}
