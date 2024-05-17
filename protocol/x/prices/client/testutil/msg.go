package testutil

import (
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/network"
	pricescli "github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices/client/cli"
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
func MsgQueryAllMarketPriceExec() ([]byte, error) {
	query := "docker exec interchain-security-instance interchain-security-cd query prices list-market-price  --node tcp://7.7.8.4:26658 -o json"
	data, _, err := network.QueryCustomNetwork(query)
	return data, err
}
