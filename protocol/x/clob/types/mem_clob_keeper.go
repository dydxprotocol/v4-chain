package types

import (
	"math/big"

	"cosmossdk.io/log"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"
	"github.com/dydxprotocol/v4-chain/protocol/x/affiliates/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

type MemClobKeeper interface {
	GetOrderFillAmount(
		ctx sdk.Context,
		orderId OrderId,
	) (
		exists bool,
		fillAmount satypes.BaseQuantums,
		prunableBlockHeight uint32,
	)
	ProcessSingleMatch(
		ctx sdk.Context,
		matchWithOrders *MatchWithOrders,
		affiliateOverrides map[string]bool,
		affiliateParameters types.AffiliateParameters,
	) (
		success bool,
		takerUpdateResult satypes.UpdateResult,
		makerUpdateResult satypes.UpdateResult,
		affiliateRevSharesQuoteQuantums *big.Int,
		err error,
	)
	CanDeleverageSubaccount(
		ctx sdk.Context,
		subaccountId satypes.SubaccountId,
		perpetualId uint32,
	) (bool, bool, error)
	GetStatePosition(
		ctx sdk.Context,
		subaccountId satypes.SubaccountId,
		clobPairId ClobPairId,
	) (
		positionSizeQuantums *big.Int,
	)
	ReplayPlaceOrder(
		ctx sdk.Context,
		msg *MsgPlaceOrder,
	) (
		orderSizeOptimisticallyFilledFromMatchingQuantums satypes.BaseQuantums,
		orderStatus OrderStatus,
		offchainUpdates *OffchainUpdates,
		err error,
	)
	AddPreexistingStatefulOrder(
		ctx sdk.Context,
		order *Order,
		memclob MemClob,
	) (
		orderSizeOptimisticallyFilledFromMatchingQuantums satypes.BaseQuantums,
		orderStatus OrderStatus,
		offchainUpdates *OffchainUpdates,
		err error,
	)
	CancelShortTermOrder(
		ctx sdk.Context,
		msgCancelOrder *MsgCancelOrder,
	) error
	GetLongTermOrderPlacement(
		ctx sdk.Context,
		orderId OrderId,
	) (val LongTermOrderPlacement, found bool)
	SetLongTermOrderPlacement(
		ctx sdk.Context,
		order Order,
		blockHeight uint32,
	)
	OffsetSubaccountPerpetualPosition(
		ctx sdk.Context,
		liquidatedSubaccountId satypes.SubaccountId,
		perpetualId uint32,
		deltaQuantumsTotal *big.Int,
		isFinalSettlement bool,
	) (
		fills []MatchPerpetualDeleveraging_Fill,
		deltaQuantumsRemaining *big.Int,
	)
	GetIndexerEventManager() indexer_manager.IndexerEventManager
	IsLiquidatable(
		ctx sdk.Context,
		subaccountId satypes.SubaccountId,
	) (
		bool,
		error,
	)
	ValidateSubaccountEquityTierLimitForStatefulOrder(
		ctx sdk.Context,
		order Order,
	) error
	Logger(
		ctx sdk.Context,
	) log.Logger
	SendOrderbookUpdates(
		ctx sdk.Context,
		offchainUpdates *OffchainUpdates,
	)
	SendOrderbookFillUpdate(
		ctx sdk.Context,
		orderbookFill StreamOrderbookFill,
	)
	SendTakerOrderStatus(
		ctx sdk.Context,
		takerOrder StreamTakerOrder,
	)
	AddOrderToOrderbookSubaccountUpdatesCheck(
		ctx sdk.Context,
		subaccountId satypes.SubaccountId,
		order PendingOpenOrder,
	) satypes.UpdateResult
	MaybeValidateAuthenticators(ctx sdk.Context, txBytes []byte) error
}
