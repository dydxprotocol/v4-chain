package canonicalusdc

import (
	"github.com/dydxprotocol/v4-chain/protocol/x/canonicalusdc/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/canonicalusdc/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func InitGenesis(ctx sdk.Context, k keeper.Keeper, genesis types.GenesisState) {
	if err := genesis.Validate(); err != nil {
		panic(err)
	}
	if err := k.ValidateRouteBindings(ctx, genesis.Controls); err != nil {
		panic(err)
	}
	if err := k.SetControls(ctx, genesis.Controls); err != nil {
		panic(err)
	}
	if err := k.SetLedger(ctx, genesis.Ledger); err != nil {
		panic(err)
	}
	if err := k.ReplaceParticipants(ctx, genesis.Participants); err != nil {
		panic(err)
	}
	for _, pending := range genesis.PendingSettlements {
		if err := k.SetPending(ctx, pending); err != nil {
			panic(err)
		}
	}
}

func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	return &types.GenesisState{
		Controls:           k.GetControls(ctx),
		Ledger:             k.GetLedger(ctx),
		Participants:       k.GetAllParticipants(ctx),
		PendingSettlements: k.GetAllPending(ctx),
	}
}
