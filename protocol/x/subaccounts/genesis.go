package subaccounts

import (
	"math/big"

	indexerevents "github.com/StreamFinance-Protocol/stream-chain/protocol/indexer/events"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/indexer/indexer_manager"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/subaccounts/keeper"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/subaccounts/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes the subaccounts module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	k.InitializeForGenesis(ctx)

	// Set all the subaccounts
	for _, elem := range genState.Subaccounts {
		if elem.AssetYieldIndex == "" {
			elem.AssetYieldIndex = big.NewRat(1, 1).String()
		}

		k.SetSubaccount(ctx, elem)
		k.GetIndexerEventManager().AddTxnEvent(
			ctx,
			indexerevents.SubtypeSubaccountUpdate,
			indexerevents.SubaccountUpdateEventVersion,
			indexer_manager.GetBytes(
				indexerevents.NewSubaccountUpdateEvent(
					elem.Id,
					elem.PerpetualPositions,
					elem.AssetPositions,
					nil,
					elem.AssetYieldIndex,
				),
			),
		)
	}
}

// ExportGenesis returns the subaccounts module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()

	genesis.Subaccounts = k.GetAllSubaccount(ctx)

	return genesis
}
