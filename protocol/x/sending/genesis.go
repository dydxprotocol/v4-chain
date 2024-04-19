package sending

import (
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/sending/keeper"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/sending/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes the sending module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	k.InitializeForGenesis(ctx)
}

// ExportGenesis returns the sending module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()

	return genesis
}
