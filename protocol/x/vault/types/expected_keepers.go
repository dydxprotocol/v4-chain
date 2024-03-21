package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
)

type ClobKeeper interface {
	// Clob Pair.
	GetClobPair(ctx sdk.Context, id clobtypes.ClobPairId) (val clobtypes.ClobPair, found bool)

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
	GetPerpetual(
		ctx sdk.Context,
		id uint32,
	) (val perptypes.Perpetual, err error)
}

type PricesKeeper interface {
	GetMarketParam(
		ctx sdk.Context,
		id uint32,
	) (market pricestypes.MarketParam, exists bool)
	GetMarketPrice(
		ctx sdk.Context,
		id uint32,
	) (pricestypes.MarketPrice, error)
}

type SubaccountsKeeper interface {
}
