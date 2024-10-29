package testutil

import (
	"fmt"
	"math/big"
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/network"
	clobcli "github.com/StreamFinance-Protocol/stream-chain/protocol/x/clob/client/cli"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/clob/types"
	satypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/subaccounts/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stretchr/testify/require"

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

func CreateBankGenesisState(
	t testing.TB,
	cfg network.Config,
) []byte {
	bankGenState := banktypes.GenesisState{
		Balances: []banktypes.Balance{
			{
				Address: "dydx1v88c3xv9xyv3eetdx0tvcmq7ung3dywp5upwc6",
				Coins: []sdk.Coin{
					sdk.NewInt64Coin(
						"utdai",
						10000000000,
					),
				},
			},
			{
				Address: "dydx1r3fsd6humm0ghyq0te5jf8eumklmclya37zle0",
				Coins: []sdk.Coin{
					{
						Denom:  "ibc/DEEFE2DEFDC8EA8879923C4CCA42BB888C3CD03FF7ECFEFB1C2FEC27A732ACC8",
						Amount: sdkmath.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(22), nil)),
					},
				},
			},
		},
	}

	bankbuf, err := cfg.Codec.MarshalJSON(&bankGenState)
	require.NoError(t, err)

	return bankbuf
}
