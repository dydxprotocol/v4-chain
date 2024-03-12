package types

import (
	context "context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
)

type BankKeeper interface {
	GetBalance(ctx context.Context, addr sdk.AccAddress, denom string) sdk.Coin
}

type ClobKeeper interface {
	GetAllClobPairs(ctx sdk.Context) (list []clobtypes.ClobPair)
	GetLongTermOrderPlacement(
		ctx sdk.Context,
		orderId clobtypes.OrderId,
	) (val clobtypes.LongTermOrderPlacement, found bool)

	PlaceStatefulOrder(ctx sdk.Context, msg *clobtypes.MsgPlaceOrder) error
	CancelStatefulOrder(
		ctx sdk.Context,
		msg *clobtypes.MsgCancelOrder,
	) (err error)

	GetProcessProposerMatchesEvents(ctx sdk.Context) clobtypes.ProcessProposerMatchesEvents
	MustSetProcessProposerMatchesEvents(
		ctx sdk.Context,
		processProposerMatchesEvents clobtypes.ProcessProposerMatchesEvents,
	)
}

type PerpetualsKeeper interface {
	GetAllPerpetuals(ctx sdk.Context) (list []perptypes.Perpetual)
}

type PricesKeeper interface {
	GetAllMarketPrices(ctx sdk.Context) (marketPrices []pricestypes.MarketPrice)
}
