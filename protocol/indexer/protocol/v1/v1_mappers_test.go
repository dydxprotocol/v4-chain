package v1_test

import (
	"fmt"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	v1 "github.com/dydxprotocol/v4-chain/protocol/indexer/protocol/v1"
	v1types "github.com/dydxprotocol/v4-chain/protocol/indexer/protocol/v1/types"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	vaulttypes "github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
	"github.com/stretchr/testify/require"
)

func TestSubaccountIdToIndexerSubaccountId(t *testing.T) {
	subaccountId := constants.Alice_Num1
	expectedSubaccountId := v1types.IndexerSubaccountId{
		Owner:  subaccountId.Owner,
		Number: subaccountId.Number,
	}

	require.Equal(
		t,
		expectedSubaccountId,
		v1.SubaccountIdToIndexerSubaccountId(subaccountId),
	)
}

func TestPerpetualPositionToIndexerPerpetualPosition(t *testing.T) {
	position := &constants.Short_Perp_1ETH_NegativeFunding
	fundingPayments := map[uint32]dtypes.SerializableInt{
		position.PerpetualId: dtypes.NewInt(100),
	}
	expectedPerpetualPosition := &v1types.IndexerPerpetualPosition{
		PerpetualId:    position.PerpetualId,
		Quantums:       position.Quantums,
		FundingIndex:   position.FundingIndex,
		FundingPayment: dtypes.NewInt(100),
	}

	require.Equal(
		t,
		expectedPerpetualPosition,
		v1.PerpetualPositionToIndexerPerpetualPosition(
			position,
			fundingPayments[position.PerpetualId],
		),
	)
}

func TestPerpetualPositionsToIndexerPerpetualPositions(t *testing.T) {
	position := &constants.Short_Perp_1ETH_NegativeFunding
	position2 := &constants.Long_Perp_1BTC_PositiveFunding

	tests := map[string]struct {
		// Input
		positions       []*satypes.PerpetualPosition
		fundingPayments map[uint32]dtypes.SerializableInt

		// Expectations
		expectedPerpetualPositions []*v1types.IndexerPerpetualPosition
	}{
		"Maps slice of PerpetualPosition to slice of IndexerPerpetualPosition with no funding payments": {
			positions: []*satypes.PerpetualPosition{
				position,
				position2,
			},
			expectedPerpetualPositions: []*v1types.IndexerPerpetualPosition{
				{
					PerpetualId:    position.PerpetualId,
					Quantums:       position.Quantums,
					FundingIndex:   position.FundingIndex,
					FundingPayment: dtypes.ZeroInt(),
				},
				{
					PerpetualId:    position2.PerpetualId,
					Quantums:       position2.Quantums,
					FundingIndex:   position2.FundingIndex,
					FundingPayment: dtypes.ZeroInt(),
				},
			},
		},
		"Maps slice of PerpetualPosition to slice of IndexerPerpetualPosition with non-zero funding payments": {
			positions: []*satypes.PerpetualPosition{
				position,
				position2,
			},
			fundingPayments: map[uint32]dtypes.SerializableInt{
				position.PerpetualId:  dtypes.NewInt(100),
				position2.PerpetualId: dtypes.NewInt(-100),
			},
			expectedPerpetualPositions: []*v1types.IndexerPerpetualPosition{
				{
					PerpetualId:    position.PerpetualId,
					Quantums:       position.Quantums,
					FundingIndex:   position.FundingIndex,
					FundingPayment: dtypes.NewInt(100),
				},
				{
					PerpetualId:    position2.PerpetualId,
					Quantums:       position2.Quantums,
					FundingIndex:   position2.FundingIndex,
					FundingPayment: dtypes.NewInt(-100),
				},
			},
		},
		"Maps empty slice to empty slice": {
			positions:                  []*satypes.PerpetualPosition{},
			expectedPerpetualPositions: []*v1types.IndexerPerpetualPosition{},
		},
		"Maps nil to nil slice": {
			positions:                  nil,
			expectedPerpetualPositions: nil,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			require.Equal(
				t,
				tc.expectedPerpetualPositions,
				v1.PerpetualPositionsToIndexerPerpetualPositions(
					tc.positions,
					tc.fundingPayments,
				),
			)
		})
	}
}

func TestAssetPositionToIndexerAssetPosition(t *testing.T) {
	position := &constants.Long_Asset_1BTC
	expectedAssetPosition := &v1types.IndexerAssetPosition{
		AssetId:  position.AssetId,
		Quantums: position.Quantums,
		Index:    position.Index,
	}

	require.Equal(
		t,
		v1.AssetPositionToIndexerAssetPosition(position),
		expectedAssetPosition,
	)
}

func TestAssetPositionsToIndexerAssetPositions(t *testing.T) {
	position := &constants.Long_Asset_1BTC
	position2 := &constants.Usdc_Asset_100_000

	tests := map[string]struct {
		// Input
		positions []*satypes.AssetPosition

		// Expectations
		expectedAssetPositions []*v1types.IndexerAssetPosition
	}{
		"Maps slice of AssetPosition to slice of IndexerAssetPosition": {
			positions: []*satypes.AssetPosition{
				position,
				position2,
			},
			expectedAssetPositions: []*v1types.IndexerAssetPosition{
				{
					AssetId:  position.AssetId,
					Quantums: position.Quantums,
					Index:    position.Index,
				},
				{
					AssetId:  position2.AssetId,
					Quantums: position2.Quantums,
					Index:    position2.Index,
				},
			},
		},
		"Maps empty slice to empty slice": {
			positions:              []*satypes.AssetPosition{},
			expectedAssetPositions: []*v1types.IndexerAssetPosition{},
		},
		"Maps nil to nil slice": {
			positions:              nil,
			expectedAssetPositions: nil,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			require.Equal(
				t,
				tc.expectedAssetPositions,
				v1.AssetPositionsToIndexerAssetPositions(tc.positions),
			)
		})
	}
}

func TestOrderIdToIndexerOrderId(t *testing.T) {
	orderId := constants.LongTermOrderId_Alice_Num1_ClientId3_Clob1
	expectedOrderId := v1types.IndexerOrderId{
		SubaccountId: v1types.IndexerSubaccountId{
			Owner:  orderId.SubaccountId.Owner,
			Number: orderId.SubaccountId.Number,
		},
		ClientId:   orderId.ClientId,
		ClobPairId: orderId.ClobPairId,
		OrderFlags: orderId.OrderFlags,
	}

	require.Equal(
		t,
		expectedOrderId,
		v1.OrderIdToIndexerOrderId(orderId),
	)
}

func TestOrderSideToIndexerOrderSide(t *testing.T) {
	tests := map[string]struct {
		// Input
		side clobtypes.Order_Side

		// Expectations
		expectedSide v1types.IndexerOrder_Side
	}{}
	// Iterate through all the values for Order_Side to create test cases.
	for name, value := range clobtypes.Order_Side_value {
		testName := fmt.Sprintf("Converts Order_Side %s to IndexerOrderV1_Side", name)
		tests[testName] = struct {
			side         clobtypes.Order_Side
			expectedSide v1types.IndexerOrder_Side
		}{
			side:         clobtypes.Order_Side(value),
			expectedSide: v1types.IndexerOrder_Side(v1types.IndexerOrder_Side_value[name]),
		}
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			require.Equal(
				t,
				tc.expectedSide,
				v1.OrderSideToIndexerOrderSide(tc.side),
			)
		})
	}
}

func TestOrderTimeInForceToIndexerOrderTimeInForce(t *testing.T) {
	tests := map[string]struct {
		// Input
		timeInForce clobtypes.Order_TimeInForce

		// Expectations
		expectedTimeInForce v1types.IndexerOrder_TimeInForce
	}{}
	// Iterate through all the values for Order_TimeInForce to create test cases.
	for name, value := range clobtypes.Order_TimeInForce_value {
		testName := fmt.Sprintf("Converts Order_TimeInForce %s to IndexerOrderV1_TimeInForce", name)
		tests[testName] = struct {
			timeInForce         clobtypes.Order_TimeInForce
			expectedTimeInForce v1types.IndexerOrder_TimeInForce
		}{
			timeInForce:         clobtypes.Order_TimeInForce(value),
			expectedTimeInForce: v1types.IndexerOrder_TimeInForce(v1types.IndexerOrder_TimeInForce_value[name]),
		}
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			require.Equal(
				t,
				tc.expectedTimeInForce,
				v1.OrderTimeInForceToIndexerOrderTimeInForce(tc.timeInForce),
			)
		})
	}
}

func TestOrderConditionTypeToIndexerOrderConditionType(t *testing.T) {
	tests := map[string]struct {
		// Input
		conditionType clobtypes.Order_ConditionType

		// Expectations
		expectedConditionType v1types.IndexerOrder_ConditionType
	}{}
	// Iterate through all the values for Order_ConditionType to create test cases.
	for name, value := range clobtypes.Order_ConditionType_value {
		testName := fmt.Sprintf("Converts Order_ConditionType %s to IndexerOrderV1_ConditionType", name)
		tests[testName] = struct {
			conditionType         clobtypes.Order_ConditionType
			expectedConditionType v1types.IndexerOrder_ConditionType
		}{
			conditionType:         clobtypes.Order_ConditionType(value),
			expectedConditionType: v1types.IndexerOrder_ConditionType(v1types.IndexerOrder_ConditionType_value[name]),
		}
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			require.Equal(
				t,
				tc.expectedConditionType,
				v1.OrderConditionTypeToIndexerOrderConditionType(tc.conditionType),
			)
		})
	}
}

func TestOrderToIndexerOrderV1(t *testing.T) {
	shortTermOrder := constants.Order_Alice_Num1_Id2_Clob1_Buy67_Price5_GTB20
	statefulOrder := constants.LongTermOrder_Bob_Num0_Id1_Clob0_Sell50_Price10_GTBT15

	tests := map[string]struct {
		// Input
		order clobtypes.Order

		// Expectations
		expectedOrder v1types.IndexerOrder
	}{
		"Maps short term order to IndexerOrderV1": {
			order: shortTermOrder,
			expectedOrder: v1types.IndexerOrder{
				OrderId: v1types.IndexerOrderId{
					SubaccountId: v1types.IndexerSubaccountId{
						Owner:  shortTermOrder.OrderId.SubaccountId.Owner,
						Number: shortTermOrder.OrderId.SubaccountId.Number,
					},
					ClientId:   shortTermOrder.OrderId.ClientId,
					ClobPairId: shortTermOrder.OrderId.ClobPairId,
					OrderFlags: shortTermOrder.OrderId.OrderFlags,
				},
				Side:     v1.OrderSideToIndexerOrderSide(shortTermOrder.Side),
				Quantums: shortTermOrder.Quantums,
				Subticks: shortTermOrder.Subticks,
				GoodTilOneof: &v1types.IndexerOrder_GoodTilBlock{
					GoodTilBlock: shortTermOrder.GoodTilOneof.(*clobtypes.Order_GoodTilBlock).GoodTilBlock,
				},
				TimeInForce:                     v1.OrderTimeInForceToIndexerOrderTimeInForce(shortTermOrder.TimeInForce),
				ReduceOnly:                      shortTermOrder.ReduceOnly,
				ClientMetadata:                  shortTermOrder.ClientMetadata,
				ConditionType:                   v1.OrderConditionTypeToIndexerOrderConditionType(shortTermOrder.ConditionType),
				ConditionalOrderTriggerSubticks: shortTermOrder.ConditionalOrderTriggerSubticks,
			},
		},
		"Maps stateful order to IndexerOrderV1": {
			order: statefulOrder,
			expectedOrder: v1types.IndexerOrder{
				OrderId: v1types.IndexerOrderId{
					SubaccountId: v1types.IndexerSubaccountId{
						Owner:  statefulOrder.OrderId.SubaccountId.Owner,
						Number: statefulOrder.OrderId.SubaccountId.Number,
					},
					ClientId:   statefulOrder.OrderId.ClientId,
					ClobPairId: statefulOrder.OrderId.ClobPairId,
					OrderFlags: statefulOrder.OrderId.OrderFlags,
				},
				Side:     v1.OrderSideToIndexerOrderSide(statefulOrder.Side),
				Quantums: statefulOrder.Quantums,
				Subticks: statefulOrder.Subticks,
				GoodTilOneof: &v1types.IndexerOrder_GoodTilBlockTime{
					GoodTilBlockTime: statefulOrder.GoodTilOneof.(*clobtypes.Order_GoodTilBlockTime).GoodTilBlockTime,
				},
				TimeInForce:                     v1.OrderTimeInForceToIndexerOrderTimeInForce(statefulOrder.TimeInForce),
				ReduceOnly:                      statefulOrder.ReduceOnly,
				ClientMetadata:                  statefulOrder.ClientMetadata,
				ConditionType:                   v1.OrderConditionTypeToIndexerOrderConditionType(statefulOrder.ConditionType),
				ConditionalOrderTriggerSubticks: statefulOrder.ConditionalOrderTriggerSubticks,
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			require.Equal(
				t,
				tc.expectedOrder,
				v1.OrderToIndexerOrder(tc.order),
			)
		})
	}
}

func TestOrderToIndexerOrder_Panic(t *testing.T) {
	invalidOrder := constants.Order_Alice_Num1_Id2_Clob1_Buy67_Price5_GTB20
	invalidOrder.GoodTilOneof = nil

	require.Panics(t, func() {
		v1.OrderToIndexerOrder(invalidOrder)
	})
}

func TestConvertToClobPairStatus(t *testing.T) {
	type convertToClobPairStatusTestCase struct {
		status         clobtypes.ClobPair_Status
		expectedStatus v1types.ClobPairStatus
		expectedPanic  string
	}

	tests := make(map[string]convertToClobPairStatusTestCase)
	// Iterate through all the values for ClobPair_Status to create test cases.
	for name, value := range clobtypes.ClobPair_Status_value {
		testName := fmt.Sprintf("Converts ClobPair_Status %s to v1.ClobPairStatus", name)
		testCase := convertToClobPairStatusTestCase{
			status:         clobtypes.ClobPair_Status(value),
			expectedStatus: v1types.ClobPairStatus(clobtypes.ClobPair_Status_value[name]),
		}
		if value == int32(clobtypes.ClobPair_STATUS_UNSPECIFIED) {
			testCase.expectedPanic = fmt.Sprintf(
				"ConvertToClobPairStatus: invalid clob pair status: %+v",
				clobtypes.ClobPair_STATUS_UNSPECIFIED,
			)
		}
		tests[testName] = testCase
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if tc.expectedPanic != "" {
				require.PanicsWithValue(
					t,
					tc.expectedPanic,
					func() {
						v1.ConvertToClobPairStatus(tc.status)
					},
				)
			} else {
				require.Equal(
					t,
					tc.expectedStatus,
					v1.ConvertToClobPairStatus(tc.status),
				)
			}
		})
	}
}

func TestConvertToPerpetualMarketType(t *testing.T) {
	type convertToPerpetualMarketTypeTestCase struct {
		status         perptypes.PerpetualMarketType
		expectedStatus v1types.PerpetualMarketType
		expectedPanic  string
	}

	tests := make(map[string]convertToPerpetualMarketTypeTestCase)
	// Iterate through all the values for PerpetualMarketType to create test cases.
	for name, value := range perptypes.PerpetualMarketType_value {
		testName := fmt.Sprintf("Converts PerpetualMarketType %s to v1.PerpetualMarketType", name)
		testCase := convertToPerpetualMarketTypeTestCase{
			status:         perptypes.PerpetualMarketType(value),
			expectedStatus: v1types.PerpetualMarketType(perptypes.PerpetualMarketType_value[name]),
		}
		if value == int32(perptypes.PerpetualMarketType_PERPETUAL_MARKET_TYPE_UNSPECIFIED) {
			testCase.expectedPanic = fmt.Sprintf(
				"ConvertToPerpetualMarketType: invalid perpetual market type: %+v",
				perptypes.PerpetualMarketType_PERPETUAL_MARKET_TYPE_UNSPECIFIED,
			)
		}
		tests[testName] = testCase
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if tc.expectedPanic != "" {
				require.PanicsWithValue(
					t,
					tc.expectedPanic,
					func() {
						v1.ConvertToPerpetualMarketType(tc.status)
					},
				)
			} else {
				require.Equal(
					t,
					tc.expectedStatus,
					v1.ConvertToPerpetualMarketType(tc.status),
				)
			}
		})
	}
}

func TestVaultStatusToIndexerVaultStatus(t *testing.T) {
	tests := map[string]struct {
		// Input
		vaultStatus vaulttypes.VaultStatus

		// Expectations
		expectedVaultStatus v1types.VaultStatus
	}{}
	// Iterate through all the values for VaultStatus to create test cases.
	for name, value := range vaulttypes.VaultStatus_value {
		testName := fmt.Sprintf("Converts VaultStatus %s to IndexerVaultStatus", name)
		tests[testName] = struct {
			vaultStatus         vaulttypes.VaultStatus
			expectedVaultStatus v1types.VaultStatus
		}{
			vaultStatus:         vaulttypes.VaultStatus(value),
			expectedVaultStatus: v1types.VaultStatus(v1types.VaultStatus_value[name]),
		}
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			require.Equal(
				t,
				tc.expectedVaultStatus,
				v1.VaultStatusToIndexerVaultStatus(tc.vaultStatus),
			)
		})
	}
}
