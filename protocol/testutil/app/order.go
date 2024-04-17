package app

// This file includes clob helpers used in the end-to-end test suites. Functions here cannot live in
// protocol/testutil/clob because they depend on the TestApp struct, and would create an import cycle.

import (
	"fmt"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/app"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/lib"
	clobtest "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/clob"
	perptest "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/perpetuals"
	pricestest "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/prices"
	clobtypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/clob/types"
	pricestypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Subsitute quantums and subticks with value converted from human readable price and amount.
func MustMakeOrderFromHumanInput(
	ctx sdk.Context,
	app *app.App,
	order clobtypes.Order,
	humanPrice string,
	humanSize string,
) clobtypes.Order {
	clobPair, exists := app.ClobKeeper.GetClobPair(ctx, clobtypes.ClobPairId(order.OrderId.ClobPairId))
	if !exists {
		panic(fmt.Sprintf("clobPair does not exist: %v", order.OrderId.ClobPairId))
	}
	perp, err := app.PerpetualsKeeper.GetPerpetual(ctx, clobtest.MustPerpetualId(clobPair))
	if err != nil {
		panic(err)
	}
	baseQuantums := perptest.MustHumanSizeToBaseQuantums(humanSize, perp.Params.AtomicResolution)
	order.Quantums = baseQuantums

	marketParams, exists := app.PricesKeeper.GetMarketParam(ctx, perp.Params.MarketId)
	if !exists {
		panic(fmt.Sprintf("marketParam does not exist: %v", perp.Params.MarketId))
	}
	marketPrice := pricestest.MustHumanPriceToMarketPrice(humanPrice, marketParams.Exponent)
	subticks := clobtypes.PriceToSubticks(
		pricestypes.MarketPrice{
			Price:    marketPrice,
			Exponent: marketParams.Exponent,
		},
		clobPair,
		perp.Params.AtomicResolution,
		lib.QuoteCurrencyAtomicResolution,
	)
	order.Subticks = subticks.Num().Uint64()
	return order
}
