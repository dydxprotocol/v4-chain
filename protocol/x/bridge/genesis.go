package bridge

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4/x/bridge/keeper"
	"github.com/dydxprotocol/v4/x/bridge/types"
)

// InitGenesis initializes the bridge module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	k.InitializeForGenesis(ctx)

	if err := k.SetEventParams(ctx, genState.EventParams); err != nil {
		panic(err)
	}
	if err := k.SetProposeParams(ctx, genState.ProposeParams); err != nil {
		panic(err)
	}
	if err := k.SetSafetyParams(ctx, genState.SafetyParams); err != nil {
		panic(err)
	}

	k.SetNextAcknowledgedEventId(ctx, genState.NextAcknowledgedEventId)
}

// ExportGenesis returns the bridge module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	return &types.GenesisState{
		EventParams:             k.GetEventParams(ctx),
		ProposeParams:           k.GetProposeParams(ctx),
		SafetyParams:            k.GetSafetyParams(ctx),
		NextAcknowledgedEventId: k.GetNextAcknowledgedEventId(ctx),
	}
}
