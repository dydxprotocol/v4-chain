package stats

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/stats/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/stats/types"
)

// InitGenesis initializes the stat module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	k.InitializeForGenesis(ctx)

	if err := k.SetParams(ctx, genState.Params); err != nil {
		panic(err)
	}

	for _, addressToUserStats := range genState.AddressToUserStats {
		k.SetUserStats(ctx, addressToUserStats.Address, addressToUserStats.UserStats)
	}
}

// ExportGenesis returns the stat module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)
	return genesis
}
