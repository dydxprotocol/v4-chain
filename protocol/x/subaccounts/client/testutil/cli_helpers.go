package testutil

import (
	"fmt"

	sacli "github.com/StreamFinance-Protocol/stream-chain/protocol/x/subaccounts/client/cli"
	sdk "github.com/cosmos/cosmos-sdk/types"

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
