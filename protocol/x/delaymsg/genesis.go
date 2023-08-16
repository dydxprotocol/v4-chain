package delaymsg

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/delaymsg/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/delaymsg/types"
)

// InitGenesis initializes the delaymsg module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	k.InitializeForGenesis(ctx)

	k.SetNumMessages(ctx, genState.NumMessages)
	for _, msg := range genState.DelayedMessages {
		// panic if the module cannot be initialized by the genesis state.
		if err := k.SetDelayedMessage(ctx, msg); err != nil {
			panic(err)
		}
	}
}

// ExportGenesis returns the delaymsg module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	return &types.GenesisState{
		NumMessages:     k.GetNumMessages(ctx),
		DelayedMessages: k.GetAllDelayedMessages(ctx),
	}
}
