package prices

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	indexerevents "github.com/dydxprotocol/v4-chain/protocol/indexer/events"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"
	"github.com/dydxprotocol/v4-chain/protocol/x/prices/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	k.InitializeForGenesis(ctx)

	if len(genState.MarketPrices) != len(genState.MarketParams) {
		panic("Expected the same number of market prices and market params")
	}

	// Set all the market params and prices.
	for i, elem := range genState.MarketParams {
		if _, err := k.CreateMarket(ctx, elem, genState.MarketPrices[i]); err != nil {
			panic(err)
		}
	}

	marketPriceUpdates := keeper.GenerateMarketPriceUpdateEvents(genState.MarketPrices)
	for _, update := range marketPriceUpdates {
		k.GetIndexerEventManager().AddTxnEvent(
			ctx,
			indexerevents.SubtypeMarket,
			indexer_manager.GetB64EncodedEventMessage(
				update,
			),
			indexerevents.MarketEventVersion,
			indexer_manager.GetBytes(
				update,
			),
		)
	}
}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()

	genesis.MarketParams = k.GetAllMarketParams(ctx)
	genesis.MarketPrices = k.GetAllMarketPrices(ctx)

	return genesis
}
