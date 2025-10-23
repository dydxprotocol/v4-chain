package memclob

import (
	"math/big"
	"time"

	"cosmossdk.io/log"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	affiliatetypes "github.com/dydxprotocol/v4-chain/protocol/x/affiliates/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

var _ types.MemClobKeeper = &FakeMemClobKeeper{}

// TODO(DEC-1629): Remove FakeMemClobKeeper struct.
type FakeMemClobKeeper struct {
	collatCheckFn                        types.AddOrderToOrderbookCollateralizationCheckFn
	fillAmounts                          map[types.OrderId]satypes.BaseQuantums
	dirtyFillAmounts                     map[types.OrderId]satypes.BaseQuantums
	positionSizes                        map[satypes.SubaccountId]map[types.ClobPairId]*big.Int
	dirtyPositionSizes                   map[satypes.SubaccountId]map[types.ClobPairId]*big.Int
	orderIdToLongTermOrderPlacement      map[types.OrderId]types.LongTermOrderPlacement
	dirtyOrderIdToLongTermOrderPlacement map[types.OrderId]types.LongTermOrderPlacement
	timeToStatefulOrdersExpiring         map[time.Time][]types.OrderId
	dirtyTimeToStatefulOrdersExpiring    map[time.Time][]types.OrderId
	nextTransactionIndex                 uint32
	subaccountsToDeleverage              map[satypes.SubaccountId]bool
	statePositionFn                      types.GetStatePositionFn
	useCollatCheckFnForSingleMatch       bool
	indexerEventManager                  indexer_manager.IndexerEventManager
}

func NewFakeMemClobKeeper() *FakeMemClobKeeper {
	return &FakeMemClobKeeper{
		collatCheckFn:                        constants.CollatCheck_EmptyUpdateResults_Success,
		dirtyFillAmounts:                     make(map[types.OrderId]satypes.BaseQuantums),
		fillAmounts:                          make(map[types.OrderId]satypes.BaseQuantums),
		statePositionFn:                      constants.GetStatePosition_ZeroPositionSize,
		useCollatCheckFnForSingleMatch:       false,
		positionSizes:                        make(map[satypes.SubaccountId]map[types.ClobPairId]*big.Int),
		dirtyPositionSizes:                   make(map[satypes.SubaccountId]map[types.ClobPairId]*big.Int),
		orderIdToLongTermOrderPlacement:      make(map[types.OrderId]types.LongTermOrderPlacement),
		dirtyOrderIdToLongTermOrderPlacement: make(map[types.OrderId]types.LongTermOrderPlacement),
		timeToStatefulOrdersExpiring:         make(map[time.Time][]types.OrderId),
		dirtyTimeToStatefulOrdersExpiring:    make(map[time.Time][]types.OrderId),
		nextTransactionIndex:                 0,
		subaccountsToDeleverage:              make(map[satypes.SubaccountId]bool),
	}
}

func (f *FakeMemClobKeeper) WithCollatCheckFnForProcessSingleMatch() *FakeMemClobKeeper {
	f.useCollatCheckFnForSingleMatch = true
	return f
}

func (f *FakeMemClobKeeper) WithIndexerEventManager(
	indexerEventManager indexer_manager.IndexerEventManager,
) *FakeMemClobKeeper {
	f.indexerEventManager = indexerEventManager
	return f
}

// Commit simulates the `checkState` being reset and uncommitted.
func (f *FakeMemClobKeeper) ResetState() {
	f.dirtyPositionSizes = make(map[satypes.SubaccountId]map[types.ClobPairId]*big.Int)
	f.dirtyFillAmounts = make(map[types.OrderId]satypes.BaseQuantums)
	f.dirtyOrderIdToLongTermOrderPlacement = make(map[types.OrderId]types.LongTermOrderPlacement)
	f.dirtyTimeToStatefulOrdersExpiring = make(map[time.Time][]types.OrderId)
	f.nextTransactionIndex = 0
}

func (f *FakeMemClobKeeper) GetIndexerEventManager() indexer_manager.IndexerEventManager {
	return f.indexerEventManager
}

func (f *FakeMemClobKeeper) ReplayPlaceOrder(
	ctx sdk.Context,
	msg *types.MsgPlaceOrder,
) (
	orderSizeOptimisticallyFilledFromMatchingQuantums satypes.BaseQuantums,
	orderStatus types.OrderStatus,
	offchainUpdates *types.OffchainUpdates,
	err error,
) {
	panic("PlaceShortTermOrder not currently implemented on FakeMemClobKeeper")
}

func (f *FakeMemClobKeeper) PlaceShortTermOrder(
	ctx sdk.Context,
	msg *types.MsgPlaceOrder,
) (
	orderSizeOptimisticallyFilledFromMatchingQuantums satypes.BaseQuantums,
	orderStatus types.OrderStatus,
	err error,
) {
	panic("PlaceShortTermOrder not currently implemented on FakeMemClobKeeper")
}

func (f *FakeMemClobKeeper) CancelShortTermOrder(
	ctx sdk.Context,
	msgCancelOrder *types.MsgCancelOrder,
) error {
	panic("CancelShortTermOrder not currently implemented on FakeMemClobKeeper")
}

func (f *FakeMemClobKeeper) CanDeleverageSubaccount(
	ctx sdk.Context,
	msg satypes.SubaccountId,
	perpetualId uint32,
) (
	bool,
	bool,
	error,
) {
	return f.subaccountsToDeleverage[msg], false, nil
}

// Commit simulates `checkState.Commit()`.
func (f *FakeMemClobKeeper) CommitState() {
	for orderId, quantums := range f.dirtyFillAmounts {
		_, exists := f.fillAmounts[orderId]

		if !exists {
			f.fillAmounts[orderId] = quantums
		} else {
			f.fillAmounts[orderId] = f.fillAmounts[orderId] + quantums
		}
	}

	for subaccountId := range f.dirtyPositionSizes {
		_, exists := f.positionSizes[subaccountId]

		if !exists {
			f.positionSizes[subaccountId] = f.dirtyPositionSizes[subaccountId]
		} else {
			for clobPairId, positionSize := range f.dirtyPositionSizes[subaccountId] {
				_, exists := f.positionSizes[subaccountId][clobPairId]

				if !exists {
					f.positionSizes[subaccountId][clobPairId] = f.dirtyPositionSizes[subaccountId][clobPairId]
				} else {
					f.positionSizes[subaccountId][clobPairId] = positionSize.Add(
						positionSize,
						f.positionSizes[subaccountId][clobPairId],
					)
				}
			}
		}
	}

	f.dirtyPositionSizes = make(map[satypes.SubaccountId]map[types.ClobPairId]*big.Int)
	f.dirtyFillAmounts = make(map[types.OrderId]satypes.BaseQuantums)
	for orderId, orderPlacement := range f.dirtyOrderIdToLongTermOrderPlacement {
		f.orderIdToLongTermOrderPlacement[orderId] = orderPlacement
	}

	for time, orderIds := range f.dirtyTimeToStatefulOrdersExpiring {
		currentOrderIds, exists := f.timeToStatefulOrdersExpiring[time]
		if !exists {
			f.timeToStatefulOrdersExpiring[time] = orderIds
		} else {
			f.timeToStatefulOrdersExpiring[time] = append(
				currentOrderIds,
				orderIds...,
			)
		}
	}

	f.ResetState()
}

func (f *FakeMemClobKeeper) WithCollatCheckFn(
	fn types.AddOrderToOrderbookCollateralizationCheckFn,
) *FakeMemClobKeeper {
	f.collatCheckFn = fn
	return f
}

func (f *FakeMemClobKeeper) WithStatePositionFn(fn types.GetStatePositionFn) *FakeMemClobKeeper {
	f.statePositionFn = fn
	return f
}

func (f *FakeMemClobKeeper) SetOrderFillAmount(
	ctx sdk.Context,
	orderId types.OrderId,
	fillAmount satypes.BaseQuantums,
) {
	f.fillAmounts[orderId] = fillAmount
}

func (f *FakeMemClobKeeper) GetOrderFillAmount(
	ctx sdk.Context,
	orderId types.OrderId,
) (
	exists bool,
	fillAmount satypes.BaseQuantums,
	prunableBlockHeight uint32,
) {
	fillAmount = f.fillAmounts[orderId] + f.dirtyFillAmounts[orderId]
	return true, fillAmount, uint32(0)
}

func (f *FakeMemClobKeeper) SetLongTermOrderPlacement(
	ctx sdk.Context,
	order types.Order,
	blockHeight uint32,
) {
	order.MustBeStatefulOrder()

	f.dirtyOrderIdToLongTermOrderPlacement[order.OrderId] = types.LongTermOrderPlacement{
		Order: order,
		PlacementIndex: types.TransactionOrdering{
			BlockHeight:      blockHeight,
			TransactionIndex: f.nextTransactionIndex,
		},
	}

	f.nextTransactionIndex += 1
}

func (f *FakeMemClobKeeper) MustAddOrderToStatefulOrdersTimeSlice(
	ctx sdk.Context,
	goodTilBlockTime time.Time,
	orderId types.OrderId,
) {
	ordersExpiringAtTime, exists := f.dirtyTimeToStatefulOrdersExpiring[goodTilBlockTime]
	if !exists {
		ordersExpiringAtTime = make([]types.OrderId, 0)
	}

	f.dirtyTimeToStatefulOrdersExpiring[goodTilBlockTime] = append(ordersExpiringAtTime, orderId)
}

func (f *FakeMemClobKeeper) DoesLongTermOrderExistInState(
	ctx sdk.Context,
	order types.Order,
) bool {
	order.MustBeStatefulOrder()

	if orderPlacementInDirtyState, exists := f.dirtyOrderIdToLongTermOrderPlacement[order.OrderId]; exists {
		return order.GetOrderHash() == orderPlacementInDirtyState.Order.GetOrderHash()
	}

	if orderPlacementInState, exists := f.orderIdToLongTermOrderPlacement[order.OrderId]; exists {
		return order.GetOrderHash() == orderPlacementInState.Order.GetOrderHash()
	}

	return false
}

func (f *FakeMemClobKeeper) GetLongTermOrderPlacement(
	ctx sdk.Context,
	orderId types.OrderId,
) (val types.LongTermOrderPlacement, found bool) {
	if orderPlacementInDirtyState, exists := f.dirtyOrderIdToLongTermOrderPlacement[orderId]; exists {
		return orderPlacementInDirtyState, true
	}

	if orderPlacementInState, exists := f.orderIdToLongTermOrderPlacement[orderId]; exists {
		return orderPlacementInState, true
	}

	return val, false
}

func (f *FakeMemClobKeeper) GetStatefulOrdersTimeSlice(
	ctx sdk.Context,
	goodTilBlockTime time.Time,
) (
	orderIds []types.OrderId,
) {
	orderIds = make([]types.OrderId, 0)
	if statefulOrdersExpiringDirtyState, exists := f.timeToStatefulOrdersExpiring[goodTilBlockTime]; exists {
		orderIds = append(orderIds, statefulOrdersExpiringDirtyState...)
	}

	if statefulOrdersExpiringState, exists := f.dirtyTimeToStatefulOrdersExpiring[goodTilBlockTime]; exists {
		orderIds = append(orderIds, statefulOrdersExpiringState...)
	}

	return orderIds
}

func (f *FakeMemClobKeeper) addFakePositionSize(
	ctx sdk.Context,
	clobPairId types.ClobPairId,
	subaccountId satypes.SubaccountId,
	isBuy bool,
	fillAmount satypes.BaseQuantums,
) {
	clobPairPositionSizes, exists := f.dirtyPositionSizes[subaccountId]
	if !exists {
		clobPairPositionSizes = make(map[types.ClobPairId]*big.Int)
	}

	curPositionSize, exists := clobPairPositionSizes[clobPairId]
	if !exists {
		curPositionSize = big.NewInt(0)
	}

	if isBuy {
		curPositionSize = curPositionSize.Add(curPositionSize, fillAmount.ToBigInt())
	} else {
		curPositionSize = curPositionSize.Sub(curPositionSize, fillAmount.ToBigInt())
	}

	clobPairPositionSizes[clobPairId] = curPositionSize
	f.dirtyPositionSizes[subaccountId] = clobPairPositionSizes
}

func (f *FakeMemClobKeeper) addFakeFillAmount(
	ctx sdk.Context,
	orderId types.OrderId,
	fillAmount satypes.BaseQuantums,
) {
	curFillAmount, exists := f.dirtyFillAmounts[orderId]
	if !exists {
		curFillAmount = 0
	}

	f.dirtyFillAmounts[orderId] = curFillAmount + fillAmount
}

func (f *FakeMemClobKeeper) ProcessSingleMatch(
	ctx sdk.Context,
	matchWithOrders *types.MatchWithOrders,
	affiliateOverrides map[string]bool,
	affiliateParameters affiliatetypes.AffiliateParameters,
) (
	success bool,
	takerUpdateResult satypes.UpdateResult,
	makerUpdateResult satypes.UpdateResult,
	affiliateRevSharesQuoteQuantums *big.Int,
	err error,
) {
	makerOrder := matchWithOrders.MakerOrder
	clobPairId := matchWithOrders.TakerOrder.GetClobPairId()
	makerOrderId := makerOrder.MustGetOrder().OrderId

	takerOrder := matchWithOrders.TakerOrder

	fillAmount := matchWithOrders.FillAmount

	if !f.useCollatCheckFnForSingleMatch {
		f.addFakeFillAmount(ctx, makerOrderId, fillAmount)
		f.addFakePositionSize(
			ctx,
			clobPairId,
			makerOrderId.SubaccountId,
			makerOrder.IsBuy(),
			fillAmount,
		)

		if !matchWithOrders.TakerOrder.IsLiquidation() {
			takerOrderId := takerOrder.MustGetOrder().OrderId
			f.addFakeFillAmount(ctx, takerOrderId, fillAmount)
			f.addFakePositionSize(
				ctx,
				clobPairId,
				takerOrderId.SubaccountId,
				takerOrder.IsBuy(),
				fillAmount,
			)
		}

		return true, satypes.Success, satypes.Success, big.NewInt(0), nil
	}

	subaccountMatchedOrders := make(map[satypes.SubaccountId][]types.PendingOpenOrder)

	subaccountMatchedOrders[matchWithOrders.MakerOrder.GetSubaccountId()] = []types.PendingOpenOrder{{
		RemainingQuantums: fillAmount,
		IsBuy:             makerOrder.IsBuy(),
		Subticks:          makerOrder.GetOrderSubticks(),
		ClobPairId:        makerOrder.GetClobPairId(),
	}}

	subaccountMatchedOrders[matchWithOrders.TakerOrder.GetSubaccountId()] = []types.PendingOpenOrder{{
		RemainingQuantums: fillAmount,
		IsBuy:             takerOrder.IsBuy(),
		Subticks:          makerOrder.GetOrderSubticks(),
		ClobPairId:        takerOrder.GetClobPairId(),
	}}

	success, successPerUpdate := f.collatCheckFn(subaccountMatchedOrders)

	takerUpdateResult = successPerUpdate[matchWithOrders.TakerOrder.GetSubaccountId()]
	makerUpdateResult = successPerUpdate[matchWithOrders.MakerOrder.GetSubaccountId()]

	if success {
		f.addFakeFillAmount(ctx, makerOrderId, fillAmount)
		f.addFakePositionSize(
			ctx,
			clobPairId,
			makerOrderId.SubaccountId,
			makerOrder.IsBuy(),
			fillAmount,
		)

		if !matchWithOrders.TakerOrder.IsLiquidation() {
			takerOrderId := takerOrder.MustGetOrder().OrderId
			f.addFakeFillAmount(ctx, takerOrderId, fillAmount)
			f.addFakePositionSize(
				ctx,
				clobPairId,
				takerOrderId.SubaccountId,
				takerOrder.IsBuy(),
				fillAmount,
			)
		}
	}

	return success, takerUpdateResult, makerUpdateResult, big.NewInt(0), nil
}

func (f *FakeMemClobKeeper) GetStatePosition(
	ctx sdk.Context,
	subaccountId satypes.SubaccountId,
	clobPairId types.ClobPairId,
) (
	positionSizeQuantums *big.Int,
) {
	mockAmount := f.statePositionFn(subaccountId, clobPairId)
	testStateFillAmount := big.NewInt(0).Set(mockAmount)

	var dirtyPositionSize = big.NewInt(0)

	clobPairPositionSize, exists := f.dirtyPositionSizes[subaccountId]
	if exists {
		dirtyPositionSize, exists = clobPairPositionSize[clobPairId]
		if !exists {
			dirtyPositionSize = big.NewInt(0)
		}
	}

	var positionSize = big.NewInt(0)
	clobPairPositionSize, exists = f.positionSizes[subaccountId]
	if exists {
		positionSize, exists = clobPairPositionSize[clobPairId]
		if !exists {
			positionSize = big.NewInt(0)
		}
	}

	testStateFillAmount = testStateFillAmount.Add(testStateFillAmount, positionSize)
	testStateFillAmount = testStateFillAmount.Add(testStateFillAmount, dirtyPositionSize)

	return testStateFillAmount
}

func (f *FakeMemClobKeeper) OffsetSubaccountPerpetualPosition(
	ctx sdk.Context,
	liquidatedSubaccountId satypes.SubaccountId,
	perpetualId uint32,
	deltaQuantumsTotal *big.Int,
	isFinalSettlement bool,
) (
	fills []types.MatchPerpetualDeleveraging_Fill,
	deltaQuantumsRemaining *big.Int,
) {
	panic("This function should not be implemented as FakeMemClobKeeper is getting deprecated (CLOB-175)")
}

func (f *FakeMemClobKeeper) AddPreexistingStatefulOrder(
	ctx sdk.Context,
	order *types.Order,
	memclob types.MemClob,
) (
	orderSizeOptimisticallyFilledFromMatchingQuantums satypes.BaseQuantums,
	orderStatus types.OrderStatus,
	offchainUpdates *types.OffchainUpdates,
	err error,
) {
	panic("This function should not be implemented as FakeMemClobKeeper is getting deprecated (CLOB-175)")
}

func (f *FakeMemClobKeeper) IsLiquidatable(
	ctx sdk.Context,
	subaccountId satypes.SubaccountId,
) (
	isLiquidatable bool,
	err error,
) {
	panic("This function should not be implemented as FakeMemClobKeeper is getting deprecated (CLOB-175)")
}

func (f *FakeMemClobKeeper) ValidateSubaccountEquityTierLimitForShortTermOrder(
	ctx sdk.Context,
	order types.Order) error {
	return nil
}

func (f *FakeMemClobKeeper) ValidateSubaccountEquityTierLimitForStatefulOrder(
	ctx sdk.Context,
	order types.Order) error {
	return nil
}

func (f *FakeMemClobKeeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger()
}

func (f *FakeMemClobKeeper) SendOrderbookUpdates(
	ctx sdk.Context,
	offchainUpdates *types.OffchainUpdates,
) {
}

func (f *FakeMemClobKeeper) SendOrderbookFillUpdate(
	ctx sdk.Context,
	orderbookFill types.StreamOrderbookFill,
) {
}

func (f *FakeMemClobKeeper) SendTakerOrderStatus(
	ctx sdk.Context,
	takerOrder types.StreamTakerOrder,
) {
}

// Placeholder to satisfy interface implementation of types.MemClobKeeper
func (f *FakeMemClobKeeper) AddOrderToOrderbookSubaccountUpdatesCheck(
	ctx sdk.Context,
	subaccountId satypes.SubaccountId,
	order types.PendingOpenOrder,
) satypes.UpdateResult {
	return satypes.Success
}

func (f *FakeMemClobKeeper) MaybeValidateAuthenticators(ctx sdk.Context, txBytes []byte) error {
	return nil
}
