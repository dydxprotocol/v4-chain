package v1

import (
	"fmt"

	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

func SubaccountIdToIndexerSubaccountId(
	subaccountId satypes.SubaccountId,
) IndexerSubaccountId {
	return IndexerSubaccountId{
		Owner:  subaccountId.Owner,
		Number: subaccountId.Number,
	}
}

func PerpetualPositionToIndexerPerpetualPosition(
	perpetualPosition *satypes.PerpetualPosition,
	fundingPayment dtypes.SerializableInt,
) *IndexerPerpetualPosition {
	return &IndexerPerpetualPosition{
		PerpetualId:    perpetualPosition.PerpetualId,
		Quantums:       perpetualPosition.Quantums,
		FundingIndex:   perpetualPosition.FundingIndex,
		FundingPayment: fundingPayment,
	}
}

func PerpetualPositionsToIndexerPerpetualPositions(
	perpetualPositions []*satypes.PerpetualPosition,
	fundingPayments map[uint32]dtypes.SerializableInt,
) []*IndexerPerpetualPosition {
	if perpetualPositions == nil {
		return nil
	}
	indexerPerpetualPositions := make([]*IndexerPerpetualPosition, 0, len(perpetualPositions))
	for _, perpetualPosition := range perpetualPositions {
		// Retrieve funding payment for this perpetual position (0 by default).
		fundingPayment, exists := fundingPayments[perpetualPosition.PerpetualId]
		if !exists {
			fundingPayment = dtypes.ZeroInt()
		}
		indexerPerpetualPositions = append(
			indexerPerpetualPositions,
			PerpetualPositionToIndexerPerpetualPosition(
				perpetualPosition,
				fundingPayment,
			),
		)
	}
	return indexerPerpetualPositions
}

func AssetPositionToIndexerAssetPosition(
	assetPosition *satypes.AssetPosition,
) *IndexerAssetPosition {
	return &IndexerAssetPosition{
		AssetId:  assetPosition.AssetId,
		Quantums: assetPosition.Quantums,
		Index:    assetPosition.Index,
	}
}

func AssetPositionsToIndexerAssetPositions(
	assetPositions []*satypes.AssetPosition,
) []*IndexerAssetPosition {
	if assetPositions == nil {
		return nil
	}
	indexerAssetPositions := make([]*IndexerAssetPosition, 0, len(assetPositions))
	for _, assetPosition := range assetPositions {
		indexerAssetPositions = append(
			indexerAssetPositions,
			AssetPositionToIndexerAssetPosition(assetPosition),
		)
	}
	return indexerAssetPositions
}

func OrderIdToIndexerOrderId(
	orderId clobtypes.OrderId,
) IndexerOrderId {
	return IndexerOrderId{
		SubaccountId: SubaccountIdToIndexerSubaccountId(orderId.SubaccountId),
		ClientId:     orderId.ClientId,
		OrderFlags:   orderId.OrderFlags,
		ClobPairId:   orderId.ClobPairId,
	}
}

func OrderSideToIndexerOrderSide(
	orderSide clobtypes.Order_Side,
) IndexerOrder_Side {
	return IndexerOrder_Side(orderSide)
}

func OrderTimeInForceToIndexerOrderTimeInForce(
	orderTimeInForce clobtypes.Order_TimeInForce,
) IndexerOrder_TimeInForce {
	return IndexerOrder_TimeInForce(orderTimeInForce)
}

func OrderConditionTypeToIndexerOrderConditionType(
	orderConditionType clobtypes.Order_ConditionType,
) IndexerOrder_ConditionType {
	return IndexerOrder_ConditionType(orderConditionType)
}

func OrderToIndexerOrder(
	order clobtypes.Order,
) IndexerOrder {
	switch goodTil := order.GoodTilOneof.(type) {
	case *clobtypes.Order_GoodTilBlock:
		return orderToIndexerOrder_GoodTilBlock(
			order,
			IndexerOrder_GoodTilBlock{GoodTilBlock: goodTil.GoodTilBlock},
		)
	case *clobtypes.Order_GoodTilBlockTime:
		return orderToIndexerOrder_GoodTilBlockTime(
			order,
			IndexerOrder_GoodTilBlockTime{GoodTilBlockTime: goodTil.GoodTilBlockTime},
		)
	default:
		panic(fmt.Errorf("Unexpected GoodTilOneof in Order: %+v", order))
	}
}

func orderToIndexerOrder_GoodTilBlock(
	order clobtypes.Order,
	goodTilBlock IndexerOrder_GoodTilBlock,
) IndexerOrder {
	return IndexerOrder{
		OrderId:                         OrderIdToIndexerOrderId(order.OrderId),
		Side:                            OrderSideToIndexerOrderSide(order.Side),
		Quantums:                        order.Quantums,
		Subticks:                        order.Subticks,
		GoodTilOneof:                    &goodTilBlock,
		TimeInForce:                     OrderTimeInForceToIndexerOrderTimeInForce(order.TimeInForce),
		ReduceOnly:                      order.ReduceOnly,
		ClientMetadata:                  order.ClientMetadata,
		ConditionType:                   OrderConditionTypeToIndexerOrderConditionType(order.ConditionType),
		ConditionalOrderTriggerSubticks: order.ConditionalOrderTriggerSubticks,
	}
}

func orderToIndexerOrder_GoodTilBlockTime(
	order clobtypes.Order,
	goodTilBlockTime IndexerOrder_GoodTilBlockTime,
) IndexerOrder {
	return IndexerOrder{
		OrderId:                         OrderIdToIndexerOrderId(order.OrderId),
		Side:                            OrderSideToIndexerOrderSide(order.Side),
		Quantums:                        order.Quantums,
		Subticks:                        order.Subticks,
		GoodTilOneof:                    &goodTilBlockTime,
		TimeInForce:                     OrderTimeInForceToIndexerOrderTimeInForce(order.TimeInForce),
		ReduceOnly:                      order.ReduceOnly,
		ClientMetadata:                  order.ClientMetadata,
		ConditionType:                   OrderConditionTypeToIndexerOrderConditionType(order.ConditionType),
		ConditionalOrderTriggerSubticks: order.ConditionalOrderTriggerSubticks,
	}
}

func ConvertToClobPairStatus(status clobtypes.ClobPair_Status) ClobPairStatus {
	switch status {
	case clobtypes.ClobPair_STATUS_ACTIVE:
		return ClobPairStatus_CLOB_PAIR_STATUS_ACTIVE
	case clobtypes.ClobPair_STATUS_PAUSED:
		return ClobPairStatus_CLOB_PAIR_STATUS_PAUSED
	case clobtypes.ClobPair_STATUS_CANCEL_ONLY:
		return ClobPairStatus_CLOB_PAIR_STATUS_CANCEL_ONLY
	case clobtypes.ClobPair_STATUS_POST_ONLY:
		return ClobPairStatus_CLOB_PAIR_STATUS_POST_ONLY
	case clobtypes.ClobPair_STATUS_INITIALIZING:
		return ClobPairStatus_CLOB_PAIR_STATUS_INITIALIZING
	default:
		panic("invalid clob pair status")
	}
}
