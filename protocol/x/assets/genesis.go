package assets

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/assets/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	k.InitializeForGenesis(ctx)

	for _, asset := range genState.Assets {
		_, err := k.CreateAsset(
			ctx,
			asset.Id,
			asset.Symbol,
			asset.Denom,
			asset.DenomExponent,
			asset.HasMarket,
			asset.MarketId,
			asset.AtomicResolution,
		)
		if err != nil {
			panic(err)
		}
	}
}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Assets = k.GetAllAssets(ctx)
	return genesis
}
