package testutil

import (
	pricescli "github.com/dydxprotocol/v4/x/prices/client/cli"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/testutil"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
)

// MsgQueryAllMarketExec lists all markets in `Prices`.
func MsgQueryAllMarketExec(clientCtx client.Context) (testutil.BufferWriter, error) {
	return clitestutil.ExecTestCLICmd(
		clientCtx,
		pricescli.CmdListMarket(),
		[]string{},
	)
}
