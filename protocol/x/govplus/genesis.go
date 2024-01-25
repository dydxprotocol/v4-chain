package govplus

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/govplus/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/govplus/types"
)

// InitGenesis initializes the govplus module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {}

// ExportGenesis returns the govplus module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	return &types.GenesisState{}
}
