package simulation

import (
	"math/rand"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/clob/keeper"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/clob/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
)

func SimulateMsgProposedOperations(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		msg := simtypes.NoOpMsg(
			types.ModuleName,
			types.TypeMsgProposedOperations,
			"ProposedOperations simulation not implemented",
		)
		return msg, nil, nil
	}
}
