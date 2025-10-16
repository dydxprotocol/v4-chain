package feetiers

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/feetiers/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/feetiers/types"
)

// InitGenesis initializes the feetiers module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	k.InitializeForGenesis(ctx)

	if err := k.SetPerpetualFeeParams(ctx, genState.Params); err != nil {
		panic(err)
	}

	// Set staking tiers
	if err := k.SetStakingTiers(ctx, genState.StakingTiers); err != nil {
		panic(err)
	}
}

// ExportGenesis returns the feetiers module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	return &types.GenesisState{
		Params:       k.GetPerpetualFeeParams(ctx),
		StakingTiers: k.GetAllStakingTiers(ctx),
	}
}
