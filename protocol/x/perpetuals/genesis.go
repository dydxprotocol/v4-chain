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
		)

		if err != nil {
			panic(err)
		}
	}

	// Create all the perpetuals.
	for _, elem := range genState.Perpetuals {
		_, err := k.CreatePerpetual(
			ctx,
			elem.Params.Id,
			elem.Params.Ticker,
			elem.Params.MarketId,
			elem.Params.AtomicResolution,
			elem.Params.DefaultFundingPpm,
			elem.Params.LiquidityTier,
		)

		if err != nil {
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
