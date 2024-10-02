package app

// This file includes clob helpers used in the end-to-end test suites. Functions here cannot live in
// protocol/testutil/clob because they depend on the TestApp struct, and would create an import cycle.

import (
	"fmt"

	"github.com/cometbft/cometbft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/app"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	clobtest "github.com/dydxprotocol/v4-chain/protocol/testutil/clob"
	perptest "github.com/dydxprotocol/v4-chain/protocol/testutil/perpetuals"
	pricestest "github.com/dydxprotocol/v4-chain/protocol/testutil/prices"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
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

	exponent, err := app.PricesKeeper.GetExponent(ctx, marketParams.Pair)
	if err != nil {
		panic(err)
	}

	marketPrice := pricestest.MustHumanPriceToMarketPrice(humanPrice, exponent)
	subticks := clobtypes.PriceToSubticks(
		pricestypes.MarketPrice{
			Price:    marketPrice,
			Exponent: exponent,
		},
		clobPair,
		perp.Params.AtomicResolution,
		lib.QuoteCurrencyAtomicResolution,
	)
	order.Subticks = subticks.Num().Uint64()
	return order
}

// MustScaleOrder scales clobtypes.Order and clobtypes.MsgPlaceorder based upon the clob information provided.
// Will panic if:
//   - OrderT is an unknown type.
//   - ClobPairSrcT is an unknown type.
//   - The clob pair id can't be used to look up the clob pair from the ClobPairSrcT.
func MustScaleOrder[
	OrderT clobtypes.Order | clobtypes.MsgPlaceOrder,
	ClobPairSrcT clobtypes.ClobPair | types.GenesisDoc](
	order OrderT,
	clobPairSrc ClobPairSrcT,
) OrderT {
	var msgPlaceOrder clobtypes.MsgPlaceOrder

	// Find the clob pair id based upon the type of order passed in.
	var clobPairId clobtypes.ClobPairId
	switch v := any(order).(type) {
	case clobtypes.MsgPlaceOrder:
		clobPairId = v.Order.GetClobPairId()
		msgPlaceOrder = v
	case clobtypes.Order:
		clobPairId = v.GetClobPairId()
		msgPlaceOrder = *clobtypes.NewMsgPlaceOrder(v)
	default:
		panic(fmt.Errorf("Unknown order type %T to get order", order))
	}

	// Find the clob pair based upon the clobPairSrc of the clob information passed in.
	var clobPair clobtypes.ClobPair
	switch v := any(clobPairSrc).(type) {
	case clobtypes.ClobPair:
		clobPair = v
	case types.GenesisDoc:
		clobPairs := MustGetClobPairsFromGenesis(v)
		if hasClobPair, ok := clobPairs[clobPairId]; ok {
			clobPair = hasClobPair
		} else {
			panic(fmt.Errorf("Clob not found in genesis doc for clob id %d", clobPairId))
		}
	default:
		panic(fmt.Errorf("Unknown source type %T to get clob pair", clobPairSrc))
	}

	// Scale the order based upon the quantums and subticks passed in.
	msgPlaceOrder.Order.Quantums = msgPlaceOrder.Order.Quantums * clobPair.StepBaseQuantums
	msgPlaceOrder.Order.Subticks = msgPlaceOrder.Order.Subticks * uint64(clobPair.SubticksPerTick)
	msgPlaceOrder.Order.ConditionalOrderTriggerSubticks = msgPlaceOrder.Order.ConditionalOrderTriggerSubticks *
		uint64(clobPair.SubticksPerTick)

	// Return a type that matches what the user passed in for the order type.
	switch any(order).(type) {
	case clobtypes.MsgPlaceOrder:
		return any(msgPlaceOrder).(OrderT)
	case clobtypes.Order:
		return any(msgPlaceOrder.Order).(OrderT)
	default:
		panic(fmt.Errorf("Unable to convert to %T to %T", clobtypes.MsgPlaceOrder{}, order))
	}
}

// MustGetClobPairsFromGenesis unmarshals the initial genesis state and returns a map from clob pair id to clob pair.
func MustGetClobPairsFromGenesis(genesisDoc types.GenesisDoc) map[clobtypes.ClobPairId]clobtypes.ClobPair {
	var genesisState clobtypes.GenesisState
	UpdateGenesisDocWithAppStateForModule(&genesisDoc, func(genesisStatePtr *clobtypes.GenesisState) {
		genesisState = *genesisStatePtr
	})

	clobPairs := make(map[clobtypes.ClobPairId]clobtypes.ClobPair, len(genesisState.ClobPairs))
	for _, clobPair := range genesisState.ClobPairs {
		clobPairs[clobPair.GetClobPairId()] = clobPair
	}
	return clobPairs
}
