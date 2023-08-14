package testutil

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sacli "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/client/cli"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/testutil"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
)

// MsgQuerySubaccountExec executes a query for the given subaccount id.
func MsgQuerySubaccountExec(
	clientCtx client.Context,
	owner sdk.AccAddress,
	number uint32,
) (testutil.BufferWriter, error) {
	return clitestutil.ExecTestCLICmd(
		clientCtx,
		sacli.CmdShowSubaccount(),
		[]string{
			owner.String(),
			fmt.Sprint(number),
		},
	)
}
