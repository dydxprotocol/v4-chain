package v1

import (
	"fmt"

	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	v1types "github.com/dydxprotocol/v4-chain/protocol/indexer/protocol/v1/types"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	vaulttypes "github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
)

func SubaccountIdToIndexerSubaccountId(
	subaccountId satypes.SubaccountId,
) v1types.IndexerSubaccountId {
	return v1types.IndexerSubaccountId{
		Owner:  subaccountId.Owner,
		Number: subaccountId.Number,
	}
}

func PerpetualPositionToIndexerPerpetualPosition(
	perpetualPosition *satypes.PerpetualPosition,
	fundingPayment dtypes.SerializableInt,
) *v1types.IndexerPerpetualPosition {
	return &v1types.IndexerPerpetualPosition{
		PerpetualId:    perpetualPosition.PerpetualId,
		Quantums:       perpetualPosition.Quantums,
		FundingIndex:   perpetualPosition.FundingIndex,
		FundingPayment: fundingPayment,
	}
}

func PerpetualPositionsToIndexerPerpetualPositions(
	perpetualPositions []*satypes.PerpetualPosition,
	fundingPayments map[uint32]dtypes.SerializableInt,
) []*v1types.IndexerPerpetualPosition {
	if perpetualPositions == nil {
		return nil
	}
	indexerPerpetualPositions := make([]*v1types.IndexerPerpetualPosition, 0, len(perpetualPositions))
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
) *v1types.IndexerAssetPosition {
	return &v1types.IndexerAssetPosition{
		AssetId:  assetPosition.AssetId,
		Quantums: assetPosition.Quantums,
		Index:    assetPosition.Index,
	}
}

func AssetPositionsToIndexerAssetPositions(
	assetPositions []*satypes.AssetPosition,
) []*v1types.IndexerAssetPosition {
	if assetPositions == nil {
		return nil
	}
	indexerAssetPositions := make([]*v1types.IndexerAssetPosition, 0, len(assetPositions))
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
) v1types.IndexerOrderId {
	return v1types.IndexerOrderId{
		SubaccountId: SubaccountIdToIndexerSubaccountId(orderId.SubaccountId),
		ClientId:     orderId.ClientId,
		OrderFlags:   orderId.OrderFlags,
		ClobPairId:   orderId.ClobPairId,
	}
}

func OrderBuilderCodeParamsToIndexerOrderBuilderCodeParams(
	builderCodeParams *clobtypes.BuilderCodeParameters,
) *v1types.BuilderCodeParameters {
	if builderCodeParams == nil {
		return nil
	}
	return &v1types.BuilderCodeParameters{
		BuilderAddress: builderCodeParams.BuilderAddress,
		FeePpm:         builderCodeParams.FeePpm,
	}
}
func OrderSideToIndexerOrderSide(
	orderSide clobtypes.Order_Side,
) v1types.IndexerOrder_Side {
	return v1types.IndexerOrder_Side(orderSide)
}

func OrderTimeInForceToIndexerOrderTimeInForce(
	orderTimeInForce clobtypes.Order_TimeInForce,
) v1types.IndexerOrder_TimeInForce {
	return v1types.IndexerOrder_TimeInForce(orderTimeInForce)
}

func OrderConditionTypeToIndexerOrderConditionType(
	orderConditionType clobtypes.Order_ConditionType,
) v1types.IndexerOrder_ConditionType {
	return v1types.IndexerOrder_ConditionType(orderConditionType)
}

func OrderTwapParametersToIndexerOrderTwapParameters(
	orderTwapParameters *clobtypes.TwapParameters,
) *v1types.TwapParameters {
	if orderTwapParameters == nil {
		return nil
	}
	return &v1types.TwapParameters{
		Duration:       orderTwapParameters.Duration,
		Interval:       orderTwapParameters.Interval,
		PriceTolerance: orderTwapParameters.PriceTolerance,
	}
}

func OrderToIndexerOrder(
	order clobtypes.Order,
) v1types.IndexerOrder {
	switch goodTil := order.GoodTilOneof.(type) {
	case *clobtypes.Order_GoodTilBlock:
		return orderToIndexerOrder_GoodTilBlock(
			order,
			v1types.IndexerOrder_GoodTilBlock{GoodTilBlock: goodTil.GoodTilBlock},
		)
	case *clobtypes.Order_GoodTilBlockTime:
		return orderToIndexerOrder_GoodTilBlockTime(
			order,
			v1types.IndexerOrder_GoodTilBlockTime{GoodTilBlockTime: goodTil.GoodTilBlockTime},
		)
	default:
		panic(fmt.Errorf("Unexpected GoodTilOneof in Order: %+v", order))
	}
}

func orderToIndexerOrder_GoodTilBlock(
	order clobtypes.Order,
	goodTilBlock v1types.IndexerOrder_GoodTilBlock,
) v1types.IndexerOrder {
	return v1types.IndexerOrder{
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
		BuilderCodeParams:               OrderBuilderCodeParamsToIndexerOrderBuilderCodeParams(order.BuilderCodeParameters),
		OrderRouterAddress:              order.GetOrderRouterAddress(),
	}
}

func orderToIndexerOrder_GoodTilBlockTime(
	order clobtypes.Order,
	goodTilBlockTime v1types.IndexerOrder_GoodTilBlockTime,
) v1types.IndexerOrder {
	return v1types.IndexerOrder{
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
		BuilderCodeParams:               OrderBuilderCodeParamsToIndexerOrderBuilderCodeParams(order.BuilderCodeParameters),
		OrderRouterAddress:              order.GetOrderRouterAddress(),
		TwapParameters:                  OrderTwapParametersToIndexerOrderTwapParameters(order.TwapParameters),
	}
}

func ConvertToClobPairStatus(status clobtypes.ClobPair_Status) v1types.ClobPairStatus {
	switch status {
	case clobtypes.ClobPair_STATUS_ACTIVE:
		return v1types.ClobPairStatus_CLOB_PAIR_STATUS_ACTIVE
	case clobtypes.ClobPair_STATUS_PAUSED:
		return v1types.ClobPairStatus_CLOB_PAIR_STATUS_PAUSED
	case clobtypes.ClobPair_STATUS_CANCEL_ONLY:
		return v1types.ClobPairStatus_CLOB_PAIR_STATUS_CANCEL_ONLY
	case clobtypes.ClobPair_STATUS_POST_ONLY:
		return v1types.ClobPairStatus_CLOB_PAIR_STATUS_POST_ONLY
	case clobtypes.ClobPair_STATUS_INITIALIZING:
		return v1types.ClobPairStatus_CLOB_PAIR_STATUS_INITIALIZING
	case clobtypes.ClobPair_STATUS_FINAL_SETTLEMENT:
		return v1types.ClobPairStatus_CLOB_PAIR_STATUS_FINAL_SETTLEMENT
	default:
		panic(
			fmt.Sprintf(
				"ConvertToClobPairStatus: invalid clob pair status: %+v",
				status,
			),
		)
	}
}

func ConvertToPerpetualMarketType(marketType perptypes.PerpetualMarketType) v1types.PerpetualMarketType {
	switch marketType {
	case perptypes.PerpetualMarketType_PERPETUAL_MARKET_TYPE_CROSS:
		return v1types.PerpetualMarketType_PERPETUAL_MARKET_TYPE_CROSS
	case perptypes.PerpetualMarketType_PERPETUAL_MARKET_TYPE_ISOLATED:
		return v1types.PerpetualMarketType_PERPETUAL_MARKET_TYPE_ISOLATED
	default:
		panic(
			fmt.Sprintf(
				"ConvertToPerpetualMarketType: invalid perpetual market type: %+v",
				marketType,
			),
		)
	}
}

func VaultStatusToIndexerVaultStatus(vaultStatus vaulttypes.VaultStatus) v1types.VaultStatus {
	return v1types.VaultStatus(vaultStatus)
}
