package prices

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4/x/prices/keeper"
	"github.com/dydxprotocol/v4/x/prices/types"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	k.InitializeForGenesis(ctx)

	// Set all the feeds.
	for _, elem := range genState.ExchangeFeeds {
		if _, err := k.CreateExchangeFeed(
			ctx,
			elem.Name,
			elem.Memo,
		); err != nil {
			panic(err)
		}
	}

	// Set all the markets
	for i, elem := range genState.Markets {
		if _, err := k.CreateMarket(
			ctx,
			elem.Pair,
			elem.Exponent,
			elem.Exchanges,
			elem.MinExchanges,
			elem.MinPriceChangePpm,
		); err != nil {
			panic(err)
		}

		if err := k.UpdateMarketPrices(
			ctx,
			[]*types.MsgUpdateMarketPrices_MarketPrice{
				{
					MarketId: uint32(i),
					Price:    elem.Price,
				},
			},
			// Do not emit market price update events during genesis, as these are arbitrarily set.
			false,
		); err != nil {
			panic(err)
		}
	}
}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()

	genesis.ExchangeFeeds = k.GetAllExchangeFeeds(ctx)
	genesis.Markets = k.GetAllMarkets(ctx)

	return genesis
}
