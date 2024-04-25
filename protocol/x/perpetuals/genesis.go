package perpetuals

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
)

// InitGenesis initializes the perpetual module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	k.InitializeForGenesis(ctx)

	// Set parameters in state.
	err := k.SetParams(ctx, genState.Params)
	if err != nil {
		panic(err)
	}

	// Create all liquidity tiers.
	for _, elem := range genState.LiquidityTiers {
		_, err := k.SetLiquidityTier(
			ctx,
			elem.Id,
			elem.Name,
			elem.InitialMarginPpm,
			elem.MaintenanceFractionPpm,
			elem.ImpactNotional,
			elem.OpenInterestLowerCap,
			elem.OpenInterestUpperCap,
		)

		if err != nil {
			panic(err)
		}
	}

	// Initialize all the perpetuals.
	// Instead of calling `CreatePerpetual`, we call `SetPerpetual` directly to
	// to allow initializing non-zero open interest/funding index at genesis.
	for _, elem := range genState.Perpetuals {
		if err := k.ValidateAndSetPerpetual(
			ctx,
			elem,
		); err != nil {
			panic(err)
		}
	}
}

// ExportGenesis returns the perpetual module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()

	genesis.Perpetuals = k.GetAllPerpetuals(ctx)
	genesis.LiquidityTiers = k.GetAllLiquidityTiers(ctx)
	genesis.Params = k.GetParams(ctx)

	return genesis
}
