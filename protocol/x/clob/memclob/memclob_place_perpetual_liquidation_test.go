package memclob

import (
	"testing"

	clobtest "github.com/dydxprotocol/v4-chain/protocol/testutil/clob"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	sdktest "github.com/dydxprotocol/v4-chain/protocol/testutil/sdk"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

func TestPlacePerpetualLiquidation_Success(t *testing.T) {
	ctx, _, _ := sdktest.NewSdkContextWithMultistore()
	ctx = ctx.WithIsCheckTx(true)
	tests := map[string]struct {
		// State.
		placedMatchableOrders []types.MatchableOrder

		// Parameters.
		order types.LiquidationOrder

		// Expectations.
		expectedFilledSize         satypes.BaseQuantums
		expectedOrderStatus        types.OrderStatus
		expectedCollatCheck        []expectedMatch
		expectedRemainingBids      []OrderWithRemainingSize
		expectedRemainingAsks      []OrderWithRemainingSize
		expectedMatches            []expectedMatch
		expectedOperations         []types.Operation
		expectedInternalOperations []types.InternalOperation
	}{
		`Matches a liquidation buy order when it overlaps the orderbook`: {
			placedMatchableOrders: []types.MatchableOrder{
				&constants.Order_Alice_Num0_Id1_Clob0_Sell10_Price15_GTB15,
			},

			order: constants.LiquidationOrder_Bob_Num0_Clob0_Buy100_Price20_BTC,

			expectedMatches: []expectedMatch{
				{
					makerOrder:      &constants.Order_Alice_Num0_Id1_Clob0_Sell10_Price15_GTB15,
					takerOrder:      &constants.LiquidationOrder_Bob_Num0_Clob0_Buy100_Price20_BTC,
					matchedQuantums: 10,
				},
			},
			expectedCollatCheck: []expectedMatch{
				{
					makerOrder:      &constants.Order_Alice_Num0_Id1_Clob0_Sell10_Price15_GTB15,
					takerOrder:      &constants.LiquidationOrder_Bob_Num0_Clob0_Buy100_Price20_BTC,
					matchedQuantums: 10,
				},
			},
			expectedRemainingBids: []OrderWithRemainingSize{},
			expectedRemainingAsks: []OrderWithRemainingSize{},
			expectedOperations: []types.Operation{
				clobtest.NewOrderPlacementOperation(constants.Order_Alice_Num0_Id1_Clob0_Sell10_Price15_GTB15),
				clobtest.NewMatchOperation(
					&constants.LiquidationOrder_Bob_Num0_Clob0_Buy100_Price20_BTC,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num0_Id1_Clob0_Sell10_Price15_GTB15.OrderId,
							FillAmount:   10,
						},
					},
				),
			},
			expectedInternalOperations: []types.InternalOperation{
				types.NewShortTermOrderPlacementInternalOperation(constants.Order_Alice_Num0_Id1_Clob0_Sell10_Price15_GTB15),
				types.NewMatchPerpetualLiquidationInternalOperation(
					&constants.LiquidationOrder_Bob_Num0_Clob0_Buy100_Price20_BTC,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num0_Id1_Clob0_Sell10_Price15_GTB15.OrderId,
							FillAmount:   10,
						},
					},
				),
			},
		},
		`Matches a liquidation sell order when it overlaps the orderbook`: {
			placedMatchableOrders: []types.MatchableOrder{
				&constants.Order_Alice_Num1_Id11_Clob1_Buy10_Price45_GTB20,
			},

			order: constants.LiquidationOrder_Bob_Num0_Clob1_Sell70_Price10_ETH,

			expectedMatches: []expectedMatch{
				{
					makerOrder:      &constants.Order_Alice_Num1_Id11_Clob1_Buy10_Price45_GTB20,
					takerOrder:      &constants.LiquidationOrder_Bob_Num0_Clob1_Sell70_Price10_ETH,
					matchedQuantums: 10,
				},
			},
			expectedCollatCheck: []expectedMatch{
				{
					makerOrder:      &constants.Order_Alice_Num1_Id11_Clob1_Buy10_Price45_GTB20,
					takerOrder:      &constants.LiquidationOrder_Bob_Num0_Clob1_Sell70_Price10_ETH,
					matchedQuantums: 10,
				},
			},
			expectedRemainingBids: []OrderWithRemainingSize{},
			expectedRemainingAsks: []OrderWithRemainingSize{},
			expectedOperations: []types.Operation{
				clobtest.NewOrderPlacementOperation(constants.Order_Alice_Num1_Id11_Clob1_Buy10_Price45_GTB20),
				clobtest.NewMatchOperation(
					&constants.LiquidationOrder_Bob_Num0_Clob1_Sell70_Price10_ETH,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num1_Id11_Clob1_Buy10_Price45_GTB20.OrderId,
							FillAmount:   10,
						},
					},
				),
			},
			expectedInternalOperations: []types.InternalOperation{
				types.NewShortTermOrderPlacementInternalOperation(constants.Order_Alice_Num1_Id11_Clob1_Buy10_Price45_GTB20),
				types.NewMatchPerpetualLiquidationInternalOperation(
					&constants.LiquidationOrder_Bob_Num0_Clob1_Sell70_Price10_ETH,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num1_Id11_Clob1_Buy10_Price45_GTB20.OrderId,
							FillAmount:   10,
						},
					},
				),
			},
		},
		`Matches a liquidation buy order multiple times when it overlaps the orderbook`: {
			placedMatchableOrders: []types.MatchableOrder{
				&constants.Order_Alice_Num1_Id0_Clob0_Sell10_Price15_GTB20,
				&constants.Order_Alice_Num0_Id7_Clob0_Sell25_Price15_GTB20,
				&constants.Order_Bob_Num0_Id9_Clob0_Sell20_Price1000,
			},

			order: constants.LiquidationOrder_Bob_Num0_Clob0_Buy100_Price20_BTC,

			expectedMatches: []expectedMatch{
				{
					makerOrder:      &constants.Order_Alice_Num1_Id0_Clob0_Sell10_Price15_GTB20,
					takerOrder:      &constants.LiquidationOrder_Bob_Num0_Clob0_Buy100_Price20_BTC,
					matchedQuantums: 10,
				},
				{
					makerOrder:      &constants.Order_Alice_Num0_Id7_Clob0_Sell25_Price15_GTB20,
					takerOrder:      &constants.LiquidationOrder_Bob_Num0_Clob0_Buy100_Price20_BTC,
					matchedQuantums: 25,
				},
			},
			expectedCollatCheck: []expectedMatch{
				{
					makerOrder:      &constants.Order_Alice_Num1_Id0_Clob0_Sell10_Price15_GTB20,
					takerOrder:      &constants.LiquidationOrder_Bob_Num0_Clob0_Buy100_Price20_BTC,
					matchedQuantums: 10,
				},
				{
					makerOrder:      &constants.Order_Alice_Num0_Id7_Clob0_Sell25_Price15_GTB20,
					takerOrder:      &constants.LiquidationOrder_Bob_Num0_Clob0_Buy100_Price20_BTC,
					matchedQuantums: 25,
				},
			},
			expectedRemainingBids: []OrderWithRemainingSize{},
			expectedRemainingAsks: []OrderWithRemainingSize{
				{
					Order:         constants.Order_Bob_Num0_Id9_Clob0_Sell20_Price1000,
					RemainingSize: 20,
				},
			},
			expectedOperations: []types.Operation{
				clobtest.NewOrderPlacementOperation(constants.Order_Alice_Num1_Id0_Clob0_Sell10_Price15_GTB20),
				clobtest.NewOrderPlacementOperation(constants.Order_Alice_Num0_Id7_Clob0_Sell25_Price15_GTB20),
				clobtest.NewMatchOperation(
					&constants.LiquidationOrder_Bob_Num0_Clob0_Buy100_Price20_BTC,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num1_Id0_Clob0_Sell10_Price15_GTB20.OrderId,
							FillAmount:   10,
						},
						{
							MakerOrderId: constants.Order_Alice_Num0_Id7_Clob0_Sell25_Price15_GTB20.OrderId,
							FillAmount:   25,
						},
					},
				),
			},
			expectedInternalOperations: []types.InternalOperation{
				types.NewShortTermOrderPlacementInternalOperation(constants.Order_Alice_Num1_Id0_Clob0_Sell10_Price15_GTB20),
				types.NewShortTermOrderPlacementInternalOperation(constants.Order_Alice_Num0_Id7_Clob0_Sell25_Price15_GTB20),
				types.NewMatchPerpetualLiquidationInternalOperation(
					&constants.LiquidationOrder_Bob_Num0_Clob0_Buy100_Price20_BTC,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num1_Id0_Clob0_Sell10_Price15_GTB20.OrderId,
							FillAmount:   10,
						},
						{
							MakerOrderId: constants.Order_Alice_Num0_Id7_Clob0_Sell25_Price15_GTB20.OrderId,
							FillAmount:   25,
						},
					},
				),
			},
		},
		`Matches a liquidation sell order multiple times when it overlaps the orderbook`: {
			placedMatchableOrders: []types.MatchableOrder{
				&constants.Order_Alice_Num1_Clob0_Id4_Buy10_Price45_GTB20,
				&constants.Order_Bob_Num0_Id5_Clob0_Buy20_Price10_GTB22,
				&constants.Order_Alice_Num1_Id8_Clob0_Buy15_Price25_GTB31,
			},

			order: constants.LiquidationOrder_Alice_Num0_Clob0_Sell20_Price25_BTC,

			expectedMatches: []expectedMatch{
				{
					makerOrder:      &constants.Order_Alice_Num1_Clob0_Id4_Buy10_Price45_GTB20,
					takerOrder:      &constants.LiquidationOrder_Alice_Num0_Clob0_Sell20_Price25_BTC,
					matchedQuantums: 10,
				},
				{
					makerOrder:      &constants.Order_Alice_Num1_Id8_Clob0_Buy15_Price25_GTB31,
					takerOrder:      &constants.LiquidationOrder_Alice_Num0_Clob0_Sell20_Price25_BTC,
					matchedQuantums: 10,
				},
			},
			expectedCollatCheck: []expectedMatch{
				{
					makerOrder:      &constants.Order_Alice_Num1_Clob0_Id4_Buy10_Price45_GTB20,
					takerOrder:      &constants.LiquidationOrder_Alice_Num0_Clob0_Sell20_Price25_BTC,
					matchedQuantums: 10,
				},
				{
					makerOrder:      &constants.Order_Alice_Num1_Id8_Clob0_Buy15_Price25_GTB31,
					takerOrder:      &constants.LiquidationOrder_Alice_Num0_Clob0_Sell20_Price25_BTC,
					matchedQuantums: 10,
				},
			},
			expectedRemainingBids: []OrderWithRemainingSize{
				{
					Order:         constants.Order_Bob_Num0_Id5_Clob0_Buy20_Price10_GTB22,
					RemainingSize: 20,
				},
				{
					Order:         constants.Order_Alice_Num1_Id8_Clob0_Buy15_Price25_GTB31,
					RemainingSize: 5,
				},
			},
			expectedRemainingAsks: []OrderWithRemainingSize{},
			expectedOperations: []types.Operation{
				clobtest.NewOrderPlacementOperation(constants.Order_Alice_Num1_Clob0_Id4_Buy10_Price45_GTB20),
				clobtest.NewOrderPlacementOperation(constants.Order_Alice_Num1_Id8_Clob0_Buy15_Price25_GTB31),
				clobtest.NewMatchOperation(
					&constants.LiquidationOrder_Alice_Num0_Clob0_Sell20_Price25_BTC,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num1_Clob0_Id4_Buy10_Price45_GTB20.OrderId,
							FillAmount:   10,
						},
						{
							MakerOrderId: constants.Order_Alice_Num1_Id8_Clob0_Buy15_Price25_GTB31.OrderId,
							FillAmount:   10,
						},
					},
				),
			},
			expectedInternalOperations: []types.InternalOperation{
				types.NewShortTermOrderPlacementInternalOperation(constants.Order_Alice_Num1_Clob0_Id4_Buy10_Price45_GTB20),
				types.NewShortTermOrderPlacementInternalOperation(constants.Order_Alice_Num1_Id8_Clob0_Buy15_Price25_GTB31),
				types.NewMatchPerpetualLiquidationInternalOperation(
					&constants.LiquidationOrder_Alice_Num0_Clob0_Sell20_Price25_BTC,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num1_Clob0_Id4_Buy10_Price45_GTB20.OrderId,
							FillAmount:   10,
						},
						{
							MakerOrderId: constants.Order_Alice_Num1_Id8_Clob0_Buy15_Price25_GTB31.OrderId,
							FillAmount:   10,
						},
					},
				),
			},
		},
		`Cancels resting maker orders from same subaccount when liquidation order overlaps the orderbook`: {
			placedMatchableOrders: []types.MatchableOrder{
				&constants.Order_Bob_Num0_Id3_Clob1_Buy10_Price10_GTB20,
				&constants.Order_Bob_Num0_Id4_Clob1_Buy20_Price35_GTB22,
				&constants.Order_Alice_Num1_Id11_Clob1_Buy10_Price45_GTB20,
				&constants.Order_Alice_Num1_Id2_Clob1_Buy67_Price5_GTB20,
			},

			order: constants.LiquidationOrder_Bob_Num0_Clob1_Sell70_Price10_ETH,

			expectedMatches: []expectedMatch{
				{
					makerOrder:      &constants.Order_Alice_Num1_Id11_Clob1_Buy10_Price45_GTB20,
					takerOrder:      &constants.LiquidationOrder_Bob_Num0_Clob1_Sell70_Price10_ETH,
					matchedQuantums: 10,
				},
			},
			expectedCollatCheck: []expectedMatch{
				{
					makerOrder:      &constants.Order_Alice_Num1_Id11_Clob1_Buy10_Price45_GTB20,
					takerOrder:      &constants.LiquidationOrder_Bob_Num0_Clob1_Sell70_Price10_ETH,
					matchedQuantums: 10,
				},
			},
			expectedRemainingBids: []OrderWithRemainingSize{
				{
					Order:         constants.Order_Alice_Num1_Id2_Clob1_Buy67_Price5_GTB20,
					RemainingSize: 67,
				},
			},
			expectedRemainingAsks: []OrderWithRemainingSize{},
			expectedOperations: []types.Operation{
				clobtest.NewOrderPlacementOperation(constants.Order_Alice_Num1_Id11_Clob1_Buy10_Price45_GTB20),
				clobtest.NewMatchOperation(
					&constants.LiquidationOrder_Bob_Num0_Clob1_Sell70_Price10_ETH,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num1_Id11_Clob1_Buy10_Price45_GTB20.OrderId,
							FillAmount:   10,
						},
					},
				),
			},
			expectedInternalOperations: []types.InternalOperation{
				types.NewShortTermOrderPlacementInternalOperation(constants.Order_Alice_Num1_Id11_Clob1_Buy10_Price45_GTB20),
				types.NewMatchPerpetualLiquidationInternalOperation(
					&constants.LiquidationOrder_Bob_Num0_Clob1_Sell70_Price10_ETH,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num1_Id11_Clob1_Buy10_Price45_GTB20.OrderId,
							FillAmount:   10,
						},
					},
				),
			},
		},
		`Cancels partially-filled maker orders from same subaccount when liquidation order overlaps
			the orderbook and doesn't match`: {
			placedMatchableOrders: []types.MatchableOrder{
				&constants.Order_Bob_Num0_Id3_Clob1_Buy10_Price10_GTB20,
				&constants.Order_Alice_Num0_Id2_Clob1_Sell5_Price10_GTB15,
			},

			order: constants.LiquidationOrder_Bob_Num0_Clob1_Sell70_Price10_ETH,

			expectedMatches: []expectedMatch{
				{
					makerOrder:      &constants.Order_Bob_Num0_Id3_Clob1_Buy10_Price10_GTB20,
					takerOrder:      &constants.Order_Alice_Num0_Id2_Clob1_Sell5_Price10_GTB15,
					matchedQuantums: 5,
				},
			},
			expectedRemainingBids: []OrderWithRemainingSize{},
			expectedRemainingAsks: []OrderWithRemainingSize{},
			expectedOperations: []types.Operation{
				clobtest.NewOrderPlacementOperation(constants.Order_Bob_Num0_Id3_Clob1_Buy10_Price10_GTB20),
				clobtest.NewOrderPlacementOperation(constants.Order_Alice_Num0_Id2_Clob1_Sell5_Price10_GTB15),
				clobtest.NewMatchOperation(
					&constants.Order_Alice_Num0_Id2_Clob1_Sell5_Price10_GTB15,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Bob_Num0_Id3_Clob1_Buy10_Price10_GTB20.OrderId,
							FillAmount:   5,
						},
					},
				),
				clobtest.NewMatchOperation(
					&constants.LiquidationOrder_Bob_Num0_Clob1_Sell70_Price10_ETH,
					[]types.MakerFill{},
				),
			},
			expectedInternalOperations: []types.InternalOperation{
				types.NewShortTermOrderPlacementInternalOperation(constants.Order_Bob_Num0_Id3_Clob1_Buy10_Price10_GTB20),
				types.NewShortTermOrderPlacementInternalOperation(constants.Order_Alice_Num0_Id2_Clob1_Sell5_Price10_GTB15),
				types.NewMatchOrdersInternalOperation(
					constants.Order_Alice_Num0_Id2_Clob1_Sell5_Price10_GTB15,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Bob_Num0_Id3_Clob1_Buy10_Price10_GTB20.OrderId,
							FillAmount:   5,
						},
					},
				),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup memclob state and test expectations.
			memclob,
				_,
				expectedNumCollateralizationChecks,
				numCollateralChecks := placeOrderTestSetup(
				t,
				ctx,
				tc.placedMatchableOrders,
				&tc.order,
				tc.expectedCollatCheck,
				tc.expectedOrderStatus,
				satypes.BaseQuantums(0),
				nil,
				map[int]map[satypes.SubaccountId]satypes.UpdateResult{},
				constants.GetStatePosition_ZeroPositionSize,
			)

			// Run the test case and verify expectations.
			placePerpetualLiquidationAndVerifyExpectationsOperations(
				t,
				ctx,
				memclob,
				tc.order,
				numCollateralChecks,
				expectedNumCollateralizationChecks,
				tc.expectedRemainingBids,
				tc.expectedRemainingAsks,
				tc.expectedOperations,
				tc.expectedInternalOperations,
			)

			// Verify the correct offchain update messages were returned.
			// TODO(DEC-1587): Update the indexer tests to perform assertions on the expected operations queue.
			// assertPlacePerpetualLiquidationOffchainMessages(
			// 	t,
			// 	offchainUpdates,
			// 	tc.order,
			// 	tc.placedMatchableOrders,
			// 	nil,
			// 	tc.expectedMatches,
			// )
		})
	}
}

func TestPlacePerpetualLiquidation_CollatCheckFailure(t *testing.T) {
	ctx, _, _ := sdktest.NewSdkContextWithMultistore()
	ctx = ctx.WithIsCheckTx(true)
	tests := map[string]struct {
		// State.
		placedMatchableOrders          []types.MatchableOrder
		collateralizationCheckFailures map[int]map[satypes.SubaccountId]satypes.UpdateResult

		// Parameters.
		order types.LiquidationOrder

		// Expectations.
		expectedFilledSize    satypes.BaseQuantums
		expectedOrderStatus   types.OrderStatus
		expectedCollatCheck   []expectedMatch
		expectedRemainingBids []OrderWithRemainingSize
		expectedRemainingAsks []OrderWithRemainingSize
		expectedMatches       []expectedMatch
	}{
		`Matching a liquidation order continues when the only maker order fails collateralization
		checks`: {
			placedMatchableOrders: []types.MatchableOrder{
				&constants.Order_Alice_Num0_Id1_Clob0_Sell10_Price15_GTB15,
			},
			collateralizationCheckFailures: map[int]map[satypes.SubaccountId]satypes.UpdateResult{
				0: {
					constants.Alice_Num0: satypes.NewlyUndercollateralized,
				},
			},

			order: constants.LiquidationOrder_Bob_Num0_Clob0_Buy100_Price20_BTC,

			expectedMatches: []expectedMatch{},
			expectedCollatCheck: []expectedMatch{
				{
					makerOrder:      &constants.Order_Alice_Num0_Id1_Clob0_Sell10_Price15_GTB15,
					takerOrder:      &constants.LiquidationOrder_Bob_Num0_Clob0_Buy100_Price20_BTC,
					matchedQuantums: 10,
				},
			},
			expectedRemainingBids: []OrderWithRemainingSize{},
			expectedRemainingAsks: []OrderWithRemainingSize{},
		},
		`Matching a liquidation order continues if multiple maker orders fail collateralization
		checks`: {
			placedMatchableOrders: []types.MatchableOrder{
				&constants.Order_Alice_Num1_Clob0_Id4_Buy10_Price45_GTB20,
				&constants.Order_Bob_Num0_Id5_Clob0_Buy20_Price10_GTB22,
				&constants.Order_Alice_Num1_Id8_Clob0_Buy15_Price25_GTB31,
			},
			collateralizationCheckFailures: map[int]map[satypes.SubaccountId]satypes.UpdateResult{
				0: {
					constants.Alice_Num1: satypes.NewlyUndercollateralized,
				},
				1: {
					constants.Alice_Num1: satypes.StillUndercollateralized,
				},
			},

			order: constants.LiquidationOrder_Alice_Num0_Clob0_Sell20_Price25_BTC,

			expectedMatches: []expectedMatch{},
			expectedCollatCheck: []expectedMatch{
				{
					makerOrder:      &constants.Order_Alice_Num1_Clob0_Id4_Buy10_Price45_GTB20,
					takerOrder:      &constants.LiquidationOrder_Alice_Num0_Clob0_Sell20_Price25_BTC,
					matchedQuantums: 10,
				},
				{
					makerOrder:      &constants.Order_Alice_Num1_Id8_Clob0_Buy15_Price25_GTB31,
					takerOrder:      &constants.LiquidationOrder_Alice_Num0_Clob0_Sell20_Price25_BTC,
					matchedQuantums: 15,
				},
			},
			expectedRemainingBids: []OrderWithRemainingSize{
				{
					Order:         constants.Order_Bob_Num0_Id5_Clob0_Buy20_Price10_GTB22,
					RemainingSize: 20,
				},
			},
			expectedRemainingAsks: []OrderWithRemainingSize{},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup memclob state and test expectations.
			memclob,
				_,
				expectedNumCollateralizationChecks,
				numCollateralChecks := placeOrderTestSetup(
				t,
				ctx,
				tc.placedMatchableOrders,
				&tc.order,
				tc.expectedCollatCheck,
				tc.expectedOrderStatus,
				satypes.BaseQuantums(0),
				nil,
				tc.collateralizationCheckFailures,
				constants.GetStatePosition_ZeroPositionSize,
			)

			// Run the test case and verify expectations.
			offchainUpdates := placePerpetualLiquidationAndVerifyExpectations(
				t,
				ctx,
				memclob,
				tc.order,
				numCollateralChecks,
				expectedNumCollateralizationChecks,
				tc.expectedRemainingBids,
				tc.expectedRemainingAsks,
				tc.expectedMatches,
			)

			// Verify the correct offchain update messages were returned.
			assertPlacePerpetualLiquidationOffchainMessages(
				t,
				ctx,
				offchainUpdates,
				tc.order,
				tc.placedMatchableOrders,
				tc.collateralizationCheckFailures,
				tc.expectedMatches,
			)
		})
	}
}
