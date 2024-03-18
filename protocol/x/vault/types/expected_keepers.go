package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
)

type ClobKeeper interface {
	// Clob Pair.
	GetAllClobPairs(ctx sdk.Context) (list []clobtypes.ClobPair)

	// Order.
	GetLongTermOrderPlacement(
		ctx sdk.Context,
		orderId clobtypes.OrderId,
	) (val clobtypes.LongTermOrderPlacement, found bool)
	HandleMsgCancelOrder(
		ctx sdk.Context,
		msg *clobtypes.MsgCancelOrder,
	) (err error)
	HandleMsgPlaceOrder(
		ctx sdk.Context,
		msg *clobtypes.MsgPlaceOrder,
	) (err error)
}

type PerpetualsKeeper interface {
	GetAllPerpetuals(ctx sdk.Context) (list []perptypes.Perpetual)
}

type PricesKeeper interface {
	GetAllMarketPrices(ctx sdk.Context) (marketPrices []pricestypes.MarketPrice)
}

type SubaccountsKeeper interface {
}
