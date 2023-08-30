package testutil

import (
	pricescli "github.com/dydxprotocol/v4-chain/protocol/x/prices/client/cli"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/testutil"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
)

// MsgQueryAllMarketParamExec lists all markets params in `Prices`.
func MsgQueryAllMarketParamExec(clientCtx client.Context) (testutil.BufferWriter, error) {
	return clitestutil.ExecTestCLICmd(
		clientCtx,
		pricescli.CmdListMarketParam(),
		[]string{},
	)
}

// MsgQueryAllMarketPriceExec lists all markets prices in `Prices`.
func MsgQueryAllMarketPriceExec(clientCtx client.Context) (testutil.BufferWriter, error) {
	return clitestutil.ExecTestCLICmd(
		clientCtx,
		pricescli.CmdListMarketPrice(),
		[]string{},
	)
}
