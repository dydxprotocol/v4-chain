package bridge

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/bridge/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
)

// InitGenesis initializes the bridge module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	k.InitializeForGenesis(ctx)

	if err := k.UpdateEventParams(ctx, genState.EventParams); err != nil {
		panic(err)
	}
	if err := k.UpdateProposeParams(ctx, genState.ProposeParams); err != nil {
		panic(err)
	}
	if err := k.UpdateSafetyParams(ctx, genState.SafetyParams); err != nil {
		panic(err)
	}
	if err := k.SetAcknowledgedEventInfo(ctx, genState.AcknowledgedEventInfo); err != nil {
		panic(err)
	}
}

// ExportGenesis returns the bridge module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	return &types.GenesisState{
		EventParams:           k.GetEventParams(ctx),
		ProposeParams:         k.GetProposeParams(ctx),
		SafetyParams:          k.GetSafetyParams(ctx),
		AcknowledgedEventInfo: k.GetAcknowledgedEventInfo(ctx),
	}
}
