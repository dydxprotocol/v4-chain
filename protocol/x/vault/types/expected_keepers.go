package types

import (
	"context"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib/margin"
	assettypes "github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	sendingtypes "github.com/dydxprotocol/v4-chain/protocol/x/sending/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

type AssetsKeeper interface {
	GetAsset(
		ctx sdk.Context,
		assetId uint32,
	) (
		asset assettypes.Asset,
		exists bool,
	)
}

type BankKeeper interface {
	GetBalance(
		ctx context.Context,
		addr sdk.AccAddress,
		denom string,
	) sdk.Coin
}

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
		isInternalOrder bool,
	) (err error)
}

type DelayMsgKeeper interface {
	DelayMessageByBlocks(
		ctx sdk.Context,
		msg sdk.Msg,
		blockDelay uint32,
	) (
		id uint32,
		err error,
	)
}

type PerpetualsKeeper interface {
	GetPerpetual(
		ctx sdk.Context,
		id uint32,
	) (val perptypes.Perpetual, err error)
	GetLiquidityTier(
		ctx sdk.Context,
		id uint32,
	) (val perptypes.LiquidityTier, err error)
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

type SendingKeeper interface {
	ProcessTransfer(
		ctx sdk.Context,
		pendingTransfer *sendingtypes.Transfer,
	) (err error)
}

type SubaccountsKeeper interface {
	DepositFundsFromAccountToSubaccount(
		ctx sdk.Context,
		fromAccount sdk.AccAddress,
		toSubaccountId satypes.SubaccountId,
		assetId uint32,
		quantums *big.Int,
	) error
	GetNetCollateralAndMarginRequirements(
		ctx sdk.Context,
		update satypes.Update,
	) (
		risk margin.Risk,
		err error,
	)
	GetSubaccount(
		ctx sdk.Context,
		id satypes.SubaccountId,
	) satypes.Subaccount
}
