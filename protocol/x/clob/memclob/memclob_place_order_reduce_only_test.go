package memclob

import (
	"math/big"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	sdktest "github.com/dydxprotocol/v4-chain/protocol/testutil/sdk"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/require"
)

func TestPlaceOrder_ReduceOnly(t *testing.T) {
	ctx, _, _ := sdktest.NewSdkContextWithMultistore()
	ctx = ctx.WithIsCheckTx(true)
	tests := map[string]struct {
		// State.
		placedMatchableOrders          []types.MatchableOrder
		collateralizationCheckFailures map[int]map[satypes.SubaccountId]satypes.UpdateResult
		statePositionSizes             map[types.ClobPairId]map[satypes.SubaccountId]*big.Int

		// Parameters.
		order types.Order

		// Expectations.
		expectedOrderStatus                    types.OrderStatus
		expectedFilledSize                     satypes.BaseQuantums
		expectedErr                            error
		expectedPendingMatches                 []expectedMatch
		expectedExistingMatches                []expectedMatch
		expectedNewMatches                     []expectedMatch
		expectedRemainingBids                  []OrderWithRemainingSize
		expectedRemainingAsks                  []OrderWithRemainingSize
		expectedSubaccountOpenReduceOnlyOrders map[types.ClobPairId]map[satypes.SubaccountId]map[types.OrderId]bool
		expectedCancelledReduceOnlyOrders      []types.OrderId
	}{
		`Can place a reduce-only sell order with a subaccount that doesn't have an open position or
						pending matches and it's canceled`: {
			placedMatchableOrders:          []types.MatchableOrder{},
			collateralizationCheckFailures: map[int]map[satypes.SubaccountId]satypes.UpdateResult{},
			statePositionSizes: map[types.ClobPairId]map[satypes.SubaccountId]*big.Int{
				0: {
					constants.Alice_Num1: big.NewInt(0),
				},
			},

			order: constants.Order_Alice_Num1_Id1_Clob0_Sell10_Price15_GTB20_RO,

			expectedErr:             types.ErrReduceOnlyWouldIncreasePositionSize,
			expectedRemainingBids:   []OrderWithRemainingSize{},
			expectedRemainingAsks:   []OrderWithRemainingSize{},
			expectedExistingMatches: []expectedMatch{},
			expectedNewMatches:      []expectedMatch{},
			expectedSubaccountOpenReduceOnlyOrders: map[types.ClobPairId]map[satypes.SubaccountId]map[types.OrderId]bool{
				0: {},
			},
			expectedCancelledReduceOnlyOrders: []types.OrderId{},
		},
		`Can place a reduce-only buy order with a subaccount that has a long position and
									it's canceled`: {
			placedMatchableOrders:          []types.MatchableOrder{},
			collateralizationCheckFailures: map[int]map[satypes.SubaccountId]satypes.UpdateResult{},
			statePositionSizes: map[types.ClobPairId]map[satypes.SubaccountId]*big.Int{
				0: {
					constants.Alice_Num1: big.NewInt(15),
				},
			},

			order: constants.Order_Alice_Num1_Id2_Clob0_Buy20_Price30_GTB20_RO,

			expectedErr:             types.ErrReduceOnlyWouldIncreasePositionSize,
			expectedRemainingBids:   []OrderWithRemainingSize{},
			expectedRemainingAsks:   []OrderWithRemainingSize{},
			expectedExistingMatches: []expectedMatch{},
			expectedNewMatches:      []expectedMatch{},
			expectedSubaccountOpenReduceOnlyOrders: map[types.ClobPairId]map[satypes.SubaccountId]map[types.OrderId]bool{
				0: {},
			},
			expectedCancelledReduceOnlyOrders: []types.OrderId{},
		},
		`Can place a reduce-only sell order with a subaccount that has a long position and it's
									added to the orderbook`: {
			placedMatchableOrders:          []types.MatchableOrder{},
			collateralizationCheckFailures: map[int]map[satypes.SubaccountId]satypes.UpdateResult{},
			statePositionSizes: map[types.ClobPairId]map[satypes.SubaccountId]*big.Int{
				0: {
					constants.Alice_Num1: big.NewInt(15),
				},
			},

			order: constants.Order_Alice_Num1_Id1_Clob0_Sell10_Price15_GTB20_RO,

			expectedOrderStatus:   types.Success,
			expectedRemainingBids: []OrderWithRemainingSize{},
			expectedRemainingAsks: []OrderWithRemainingSize{

				{
					Order:         constants.Order_Alice_Num1_Id1_Clob0_Sell10_Price15_GTB20_RO,
					RemainingSize: 10,
				},
			},
			expectedExistingMatches: []expectedMatch{},
			expectedNewMatches:      []expectedMatch{},
			expectedSubaccountOpenReduceOnlyOrders: map[types.ClobPairId]map[satypes.SubaccountId]map[types.OrderId]bool{
				0: {
					constants.Alice_Num1: {
						constants.Order_Alice_Num1_Id1_Clob0_Sell10_Price15_GTB20_RO.OrderId: true,
					},
				},
			},
			expectedCancelledReduceOnlyOrders: []types.OrderId{},
		},
		`Can place a reduce-only buy order with a subaccount that has a short position and it's
									added to the orderbook`: {
			placedMatchableOrders:          []types.MatchableOrder{},
			collateralizationCheckFailures: map[int]map[satypes.SubaccountId]satypes.UpdateResult{},
			statePositionSizes: map[types.ClobPairId]map[satypes.SubaccountId]*big.Int{
				0: {
					constants.Alice_Num1: big.NewInt(-35),
				},
			},

			order: constants.Order_Alice_Num1_Id2_Clob0_Buy20_Price30_GTB20_RO,

			expectedOrderStatus: types.Success,
			expectedRemainingBids: []OrderWithRemainingSize{

				{
					Order:         constants.Order_Alice_Num1_Id2_Clob0_Buy20_Price30_GTB20_RO,
					RemainingSize: 20,
				},
			},
			expectedRemainingAsks:   []OrderWithRemainingSize{},
			expectedExistingMatches: []expectedMatch{},
			expectedNewMatches:      []expectedMatch{},
			expectedSubaccountOpenReduceOnlyOrders: map[types.ClobPairId]map[satypes.SubaccountId]map[types.OrderId]bool{
				0: {
					constants.Alice_Num1: {
						constants.Order_Alice_Num1_Id2_Clob0_Buy20_Price30_GTB20_RO.OrderId: true,
					},
				},
			},
			expectedCancelledReduceOnlyOrders: []types.OrderId{},
		},
		`Can place a reduce-only sell order with a subaccount that has a pending matched sell order
									and it's canceled`: {
			placedMatchableOrders: []types.MatchableOrder{
				&constants.Order_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB20,
				&constants.Order_Bob_Num0_Id4_Clob1_Buy20_Price35_GTB22,
			},
			collateralizationCheckFailures: map[int]map[satypes.SubaccountId]satypes.UpdateResult{},
			statePositionSizes: map[types.ClobPairId]map[satypes.SubaccountId]*big.Int{
				1: {
					constants.Alice_Num1: big.NewInt(-10),
					constants.Bob_Num0:   big.NewInt(0),
				},
			},

			order: constants.Order_Alice_Num1_Id5_Clob1_Sell10_Price15_GTB20_RO,

			expectedErr: types.ErrReduceOnlyWouldIncreasePositionSize,
			expectedRemainingBids: []OrderWithRemainingSize{
				{
					Order:         constants.Order_Bob_Num0_Id4_Clob1_Buy20_Price35_GTB22,
					RemainingSize: 10,
				},
			},
			expectedRemainingAsks: []OrderWithRemainingSize{},
			expectedExistingMatches: []expectedMatch{
				// Pending matched sell order generated before test case.
				{
					makerOrder:      &constants.Order_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB20,
					takerOrder:      &constants.Order_Bob_Num0_Id4_Clob1_Buy20_Price35_GTB22,
					matchedQuantums: 10,
				},
			},
			expectedNewMatches: []expectedMatch{},
			expectedSubaccountOpenReduceOnlyOrders: map[types.ClobPairId]map[satypes.SubaccountId]map[types.OrderId]bool{
				1: {},
			},
			expectedCancelledReduceOnlyOrders: []types.OrderId{},
		},
		`Can place a reduce-only buy order with a subaccount that has a pending matched sell order
									and it's added to the orderbook`: {
			placedMatchableOrders: []types.MatchableOrder{
				&constants.Order_Alice_Num1_Id0_Clob0_Sell10_Price15_GTB20,
				&constants.Order_Bob_Num0_Id6_Clob0_Buy20_Price1000_GTB22,
			},
			collateralizationCheckFailures: map[int]map[satypes.SubaccountId]satypes.UpdateResult{},
			statePositionSizes: map[types.ClobPairId]map[satypes.SubaccountId]*big.Int{
				0: {
					constants.Alice_Num1: big.NewInt(0),
					constants.Bob_Num0:   big.NewInt(0),
				},
			},

			order: constants.Order_Alice_Num1_Id2_Clob0_Buy20_Price30_GTB20_RO,

			expectedOrderStatus: types.Success,
			expectedRemainingBids: []OrderWithRemainingSize{
				{
					Order:         constants.Order_Alice_Num1_Id2_Clob0_Buy20_Price30_GTB20_RO,
					RemainingSize: 20,
				},
				{
					Order:         constants.Order_Bob_Num0_Id6_Clob0_Buy20_Price1000_GTB22,
					RemainingSize: 10,
				},
			},
			expectedRemainingAsks: []OrderWithRemainingSize{},
			expectedExistingMatches: []expectedMatch{
				// Pending matched sell order generated before test case.
				{
					makerOrder:      &constants.Order_Alice_Num1_Id0_Clob0_Sell10_Price15_GTB20,
					takerOrder:      &constants.Order_Bob_Num0_Id6_Clob0_Buy20_Price1000_GTB22,
					matchedQuantums: 10,
				},
			},
			expectedNewMatches: []expectedMatch{},
			expectedSubaccountOpenReduceOnlyOrders: map[types.ClobPairId]map[satypes.SubaccountId]map[types.OrderId]bool{
				0: {
					constants.Alice_Num1: {
						constants.Order_Alice_Num1_Id2_Clob0_Buy20_Price30_GTB20_RO.OrderId: true,
					},
				},
			},
			expectedCancelledReduceOnlyOrders: []types.OrderId{},
		},
		`Can place a reduce-only buy order with a subaccount that has a smaller long position and a
									larger pending matched sell order, and it's added to the orderbook`: {
			placedMatchableOrders: []types.MatchableOrder{
				&constants.Order_Alice_Num1_Id0_Clob0_Sell10_Price15_GTB20,
				&constants.Order_Bob_Num0_Id6_Clob0_Buy20_Price1000_GTB22,
			},
			collateralizationCheckFailures: map[int]map[satypes.SubaccountId]satypes.UpdateResult{},
			statePositionSizes: map[types.ClobPairId]map[satypes.SubaccountId]*big.Int{
				0: {
					constants.Alice_Num1: big.NewInt(5),
					constants.Bob_Num0:   big.NewInt(0),
				},
			},

			order: constants.Order_Alice_Num1_Id2_Clob0_Buy20_Price30_GTB20_RO,

			expectedOrderStatus: types.Success,
			expectedRemainingBids: []OrderWithRemainingSize{
				{
					Order:         constants.Order_Alice_Num1_Id2_Clob0_Buy20_Price30_GTB20_RO,
					RemainingSize: 20,
				},
				{
					Order:         constants.Order_Bob_Num0_Id6_Clob0_Buy20_Price1000_GTB22,
					RemainingSize: 10,
				},
			},
			expectedRemainingAsks: []OrderWithRemainingSize{},
			expectedExistingMatches: []expectedMatch{
				{
					makerOrder:      &constants.Order_Alice_Num1_Id0_Clob0_Sell10_Price15_GTB20,
					takerOrder:      &constants.Order_Bob_Num0_Id6_Clob0_Buy20_Price1000_GTB22,
					matchedQuantums: 10,
				},
			},
			expectedNewMatches: []expectedMatch{},
			expectedSubaccountOpenReduceOnlyOrders: map[types.ClobPairId]map[satypes.SubaccountId]map[types.OrderId]bool{
				0: {
					constants.Alice_Num1: {
						constants.Order_Alice_Num1_Id2_Clob0_Buy20_Price30_GTB20_RO.OrderId: true,
					},
				},
			},
			expectedCancelledReduceOnlyOrders: []types.OrderId{},
		},
		`Can place a reduce-only buy order with a subaccount that has a short position and a larger
								pending matched buy order with another CLOB pair, and it's added to the orderbook`: {
			placedMatchableOrders: []types.MatchableOrder{
				&constants.Order_Alice_Num1_Id11_Clob1_Buy10_Price45_GTB20,
				&constants.Order_Bob_Num0_Id2_Clob1_Sell12_Price13_GTB20,
			},
			collateralizationCheckFailures: map[int]map[satypes.SubaccountId]satypes.UpdateResult{},
			statePositionSizes: map[types.ClobPairId]map[satypes.SubaccountId]*big.Int{
				0: {
					constants.Alice_Num1: big.NewInt(-5),
				},
				1: {
					constants.Alice_Num1: big.NewInt(0),
					constants.Bob_Num0:   big.NewInt(0),
				},
			},

			order: constants.Order_Alice_Num1_Id2_Clob0_Buy20_Price30_GTB20_RO,

			expectedOrderStatus: types.Success,
			expectedRemainingBids: []OrderWithRemainingSize{
				{
					Order:         constants.Order_Alice_Num1_Id2_Clob0_Buy20_Price30_GTB20_RO,
					RemainingSize: 20,
				},
			},
			expectedRemainingAsks: []OrderWithRemainingSize{
				{
					Order:         constants.Order_Bob_Num0_Id2_Clob1_Sell12_Price13_GTB20,
					RemainingSize: 2,
				},
			},
			expectedExistingMatches: []expectedMatch{
				// Pending matched buy order generated before the test case.
				{
					makerOrder:      &constants.Order_Alice_Num1_Id11_Clob1_Buy10_Price45_GTB20,
					takerOrder:      &constants.Order_Bob_Num0_Id2_Clob1_Sell12_Price13_GTB20,
					matchedQuantums: 10,
				},
			},
			expectedNewMatches: []expectedMatch{},
			expectedSubaccountOpenReduceOnlyOrders: map[types.ClobPairId]map[satypes.SubaccountId]map[types.OrderId]bool{
				0: {
					constants.Alice_Num1: {
						constants.Order_Alice_Num1_Id2_Clob0_Buy20_Price30_GTB20_RO.OrderId: true,
					},
				},
			},
			expectedCancelledReduceOnlyOrders: []types.OrderId{},
		},
		`Can place a reduce-only buy order with a subaccount that has a short position and it's
								matched against resting liquidity`: {
			placedMatchableOrders: []types.MatchableOrder{
				&constants.Order_Bob_Num0_Id8_Clob0_Sell20_Price10_GTB22,
			},
			collateralizationCheckFailures: map[int]map[satypes.SubaccountId]satypes.UpdateResult{},
			statePositionSizes: map[types.ClobPairId]map[satypes.SubaccountId]*big.Int{
				0: {
					constants.Alice_Num1: big.NewInt(-20),
					constants.Bob_Num0:   big.NewInt(0),
				},
			},

			order: constants.Order_Alice_Num1_Id2_Clob0_Buy20_Price30_GTB20_RO,

			expectedOrderStatus:   types.Success,
			expectedFilledSize:    20,
			expectedRemainingBids: []OrderWithRemainingSize{},
			expectedRemainingAsks: []OrderWithRemainingSize{},
			expectedPendingMatches: []expectedMatch{
				{
					makerOrder:      &constants.Order_Bob_Num0_Id8_Clob0_Sell20_Price10_GTB22,
					takerOrder:      &constants.Order_Alice_Num1_Id2_Clob0_Buy20_Price30_GTB20_RO,
					matchedQuantums: 20,
				},
			},
			expectedExistingMatches: []expectedMatch{},
			expectedNewMatches: []expectedMatch{
				{
					makerOrder:      &constants.Order_Bob_Num0_Id8_Clob0_Sell20_Price10_GTB22,
					takerOrder:      &constants.Order_Alice_Num1_Id2_Clob0_Buy20_Price30_GTB20_RO,
					matchedQuantums: 20,
				},
			},
			expectedSubaccountOpenReduceOnlyOrders: map[types.ClobPairId]map[satypes.SubaccountId]map[types.OrderId]bool{
				0: {},
			},
			expectedCancelledReduceOnlyOrders: []types.OrderId{},
		},
		`Can place a regular order that matches and closes the position, and the taker's resting
								reduce-only orders on that CLOB are canceled`: {
			placedMatchableOrders: []types.MatchableOrder{
				&constants.Order_Alice_Num1_Id1_Clob0_Sell10_Price15_GTB20_RO,
				&constants.Order_Alice_Num1_Id4_Clob0_Sell15_Price20_GTB20_RO,
				&constants.Order_Bob_Num0_Id5_Clob0_Buy20_Price10_GTB22,
			},
			collateralizationCheckFailures: map[int]map[satypes.SubaccountId]satypes.UpdateResult{},
			statePositionSizes: map[types.ClobPairId]map[satypes.SubaccountId]*big.Int{
				0: {
					constants.Alice_Num1: big.NewInt(20),
					constants.Bob_Num0:   big.NewInt(0),
				},
			},

			order: constants.Order_Alice_Num1_Id12_Clob0_Sell20_Price5_GTB25,

			expectedOrderStatus:   types.Success,
			expectedFilledSize:    20,
			expectedRemainingBids: []OrderWithRemainingSize{},
			expectedRemainingAsks: []OrderWithRemainingSize{},
			expectedPendingMatches: []expectedMatch{
				{
					makerOrder:      &constants.Order_Bob_Num0_Id5_Clob0_Buy20_Price10_GTB22,
					takerOrder:      &constants.Order_Alice_Num1_Id12_Clob0_Sell20_Price5_GTB25,
					matchedQuantums: 20,
				},
			},
			expectedExistingMatches: []expectedMatch{},
			expectedNewMatches: []expectedMatch{
				{
					makerOrder:      &constants.Order_Bob_Num0_Id5_Clob0_Buy20_Price10_GTB22,
					takerOrder:      &constants.Order_Alice_Num1_Id12_Clob0_Sell20_Price5_GTB25,
					matchedQuantums: 20,
				},
			},
			expectedSubaccountOpenReduceOnlyOrders: map[types.ClobPairId]map[satypes.SubaccountId]map[types.OrderId]bool{
				0: {},
			},
			expectedCancelledReduceOnlyOrders: []types.OrderId{
				constants.Order_Alice_Num1_Id1_Clob0_Sell10_Price15_GTB20_RO.OrderId,
				constants.Order_Alice_Num1_Id4_Clob0_Sell15_Price20_GTB20_RO.OrderId,
			},
		},
		`Can place a regular order that matches and does not close the position, and the taker's resting
						reduce-only orders on that CLOB remain open`: {
			placedMatchableOrders: []types.MatchableOrder{
				&constants.Order_Alice_Num1_Id1_Clob0_Sell10_Price15_GTB20_RO,
				&constants.Order_Alice_Num1_Id4_Clob0_Sell15_Price20_GTB20_RO,
				&constants.Order_Bob_Num0_Id5_Clob0_Buy20_Price10_GTB22,
			},
			collateralizationCheckFailures: map[int]map[satypes.SubaccountId]satypes.UpdateResult{},
			statePositionSizes: map[types.ClobPairId]map[satypes.SubaccountId]*big.Int{
				0: {
					constants.Alice_Num1: big.NewInt(30),
					constants.Bob_Num0:   big.NewInt(0),
				},
			},

			order: constants.Order_Alice_Num1_Id12_Clob0_Sell20_Price5_GTB25,

			expectedOrderStatus:   types.Success,
			expectedFilledSize:    20,
			expectedRemainingBids: []OrderWithRemainingSize{},
			expectedRemainingAsks: []OrderWithRemainingSize{
				{
					Order:         constants.Order_Alice_Num1_Id1_Clob0_Sell10_Price15_GTB20_RO,
					RemainingSize: 10,
				},
				{
					Order:         constants.Order_Alice_Num1_Id4_Clob0_Sell15_Price20_GTB20_RO,
					RemainingSize: 15,
				},
			},
			expectedPendingMatches: []expectedMatch{
				{
					makerOrder:      &constants.Order_Bob_Num0_Id5_Clob0_Buy20_Price10_GTB22,
					takerOrder:      &constants.Order_Alice_Num1_Id12_Clob0_Sell20_Price5_GTB25,
					matchedQuantums: 20,
				},
			},
			expectedExistingMatches: []expectedMatch{},
			expectedNewMatches: []expectedMatch{
				{
					makerOrder:      &constants.Order_Bob_Num0_Id5_Clob0_Buy20_Price10_GTB22,
					takerOrder:      &constants.Order_Alice_Num1_Id12_Clob0_Sell20_Price5_GTB25,
					matchedQuantums: 20,
				},
			},
			expectedSubaccountOpenReduceOnlyOrders: map[types.ClobPairId]map[satypes.SubaccountId]map[types.OrderId]bool{
				0: {
					constants.Alice_Num1: {
						constants.Order_Alice_Num1_Id1_Clob0_Sell10_Price15_GTB20_RO.OrderId: true,
						constants.Order_Alice_Num1_Id4_Clob0_Sell15_Price20_GTB20_RO.OrderId: true,
					},
				},
			},
			expectedCancelledReduceOnlyOrders: []types.OrderId{},
		},
		`Can place a regular order that has multiple fills and closes the position, and the taker's resting
						reduce-only orders on that CLOB are canceled`: {
			placedMatchableOrders: []types.MatchableOrder{
				&constants.Order_Alice_Num1_Id5_Clob1_Sell10_Price15_GTB20_RO,
				&constants.Order_Alice_Num1_Id1_Clob0_Sell10_Price15_GTB20_RO,
				&constants.Order_Alice_Num1_Id4_Clob0_Sell15_Price20_GTB20_RO,
				&constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15,
				&constants.Order_Alice_Num0_Id6_Clob0_Buy25_Price5_GTB20,
			},
			collateralizationCheckFailures: map[int]map[satypes.SubaccountId]satypes.UpdateResult{},
			statePositionSizes: map[types.ClobPairId]map[satypes.SubaccountId]*big.Int{
				0: {
					constants.Alice_Num1: big.NewInt(20),
					constants.Alice_Num0: big.NewInt(0),
				},
				1: {
					constants.Alice_Num1: big.NewInt(15),
				},
			},

			order: constants.Order_Alice_Num1_Id12_Clob0_Sell20_Price5_GTB25,

			expectedOrderStatus: types.Success,
			expectedFilledSize:  20,
			expectedRemainingBids: []OrderWithRemainingSize{
				{
					Order:         constants.Order_Alice_Num0_Id6_Clob0_Buy25_Price5_GTB20,
					RemainingSize: 10,
				},
			},
			expectedRemainingAsks: []OrderWithRemainingSize{
				{
					Order:         constants.Order_Alice_Num1_Id5_Clob1_Sell10_Price15_GTB20_RO,
					RemainingSize: 10,
				},
			},
			expectedPendingMatches: []expectedMatch{
				{
					makerOrder:      &constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15,
					takerOrder:      &constants.Order_Alice_Num1_Id12_Clob0_Sell20_Price5_GTB25,
					matchedQuantums: 5,
				},
				{
					makerOrder:      &constants.Order_Alice_Num0_Id6_Clob0_Buy25_Price5_GTB20,
					takerOrder:      &constants.Order_Alice_Num1_Id12_Clob0_Sell20_Price5_GTB25,
					matchedQuantums: 15,
				},
			},
			expectedExistingMatches: []expectedMatch{},
			expectedNewMatches: []expectedMatch{
				{
					makerOrder:      &constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15,
					takerOrder:      &constants.Order_Alice_Num1_Id12_Clob0_Sell20_Price5_GTB25,
					matchedQuantums: 5,
				},
				{
					makerOrder:      &constants.Order_Alice_Num0_Id6_Clob0_Buy25_Price5_GTB20,
					takerOrder:      &constants.Order_Alice_Num1_Id12_Clob0_Sell20_Price5_GTB25,
					matchedQuantums: 15,
				},
			},
			expectedSubaccountOpenReduceOnlyOrders: map[types.ClobPairId]map[satypes.SubaccountId]map[types.OrderId]bool{
				0: {},
				1: {
					constants.Alice_Num1: {
						constants.Order_Alice_Num1_Id5_Clob1_Sell10_Price15_GTB20_RO.OrderId: true,
					},
				},
			},
			expectedCancelledReduceOnlyOrders: []types.OrderId{
				constants.Order_Alice_Num1_Id1_Clob0_Sell10_Price15_GTB20_RO.OrderId,
				constants.Order_Alice_Num1_Id4_Clob0_Sell15_Price20_GTB20_RO.OrderId,
			},
		},
		`Can place a regular order that matches and changes the position side, and the taker's resting
						reduce-only orders on that CLOB are canceled`: {
			placedMatchableOrders: []types.MatchableOrder{
				&constants.Order_Alice_Num1_Id2_Clob0_Buy20_Price30_GTB20_RO,
				&constants.Order_Bob_Num0_Id13_Clob0_Sell35_Price35_GTB30,
			},
			collateralizationCheckFailures: map[int]map[satypes.SubaccountId]satypes.UpdateResult{},
			statePositionSizes: map[types.ClobPairId]map[satypes.SubaccountId]*big.Int{
				0: {
					constants.Alice_Num1: big.NewInt(-20),
					constants.Bob_Num0:   big.NewInt(0),
				},
			},

			order: constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25,

			expectedOrderStatus:   types.Success,
			expectedFilledSize:    30,
			expectedRemainingBids: []OrderWithRemainingSize{},
			expectedRemainingAsks: []OrderWithRemainingSize{

				{
					Order:         constants.Order_Bob_Num0_Id13_Clob0_Sell35_Price35_GTB30,
					RemainingSize: 5,
				},
			},
			expectedPendingMatches: []expectedMatch{
				{
					makerOrder:      &constants.Order_Bob_Num0_Id13_Clob0_Sell35_Price35_GTB30,
					takerOrder:      &constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25,
					matchedQuantums: 30,
				},
			},
			expectedExistingMatches: []expectedMatch{},
			expectedNewMatches: []expectedMatch{
				{
					makerOrder:      &constants.Order_Bob_Num0_Id13_Clob0_Sell35_Price35_GTB30,
					takerOrder:      &constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25,
					matchedQuantums: 30,
				},
			},
			expectedSubaccountOpenReduceOnlyOrders: map[types.ClobPairId]map[satypes.SubaccountId]map[types.OrderId]bool{
				0: {},
			},
			expectedCancelledReduceOnlyOrders: []types.OrderId{
				constants.Order_Alice_Num1_Id2_Clob0_Buy20_Price30_GTB20_RO.OrderId,
			},
		},
		`Can place a regular taker order that matches with a maker order and the taker and maker
						change position sides, and both taker and maker's resting reduce-only orders on that CLOB
						are canceled`: {
			placedMatchableOrders: []types.MatchableOrder{
				&constants.Order_Alice_Num1_Id2_Clob0_Buy20_Price30_GTB20_RO,
				&constants.Order_Bob_Num0_Id13_Clob0_Sell35_Price35_GTB30,
				&constants.Order_Bob_Num0_Id1_Clob0_Sell15_Price50_GTB20_RO,
			},
			collateralizationCheckFailures: map[int]map[satypes.SubaccountId]satypes.UpdateResult{},
			statePositionSizes: map[types.ClobPairId]map[satypes.SubaccountId]*big.Int{
				0: {
					constants.Alice_Num1: big.NewInt(-20),
					constants.Bob_Num0:   big.NewInt(15),
				},
			},

			order: constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25,

			expectedOrderStatus:   types.Success,
			expectedFilledSize:    30,
			expectedRemainingBids: []OrderWithRemainingSize{},
			expectedRemainingAsks: []OrderWithRemainingSize{

				{
					Order:         constants.Order_Bob_Num0_Id13_Clob0_Sell35_Price35_GTB30,
					RemainingSize: 5,
				},
			},
			expectedPendingMatches: []expectedMatch{
				{
					makerOrder:      &constants.Order_Bob_Num0_Id13_Clob0_Sell35_Price35_GTB30,
					takerOrder:      &constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25,
					matchedQuantums: 30,
				},
			},
			expectedExistingMatches: []expectedMatch{},
			expectedNewMatches: []expectedMatch{
				{
					makerOrder:      &constants.Order_Bob_Num0_Id13_Clob0_Sell35_Price35_GTB30,
					takerOrder:      &constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25,
					matchedQuantums: 30,
				},
			},
			expectedSubaccountOpenReduceOnlyOrders: map[types.ClobPairId]map[satypes.SubaccountId]map[types.OrderId]bool{
				0: {},
			},
			expectedCancelledReduceOnlyOrders: []types.OrderId{
				constants.Order_Bob_Num0_Id1_Clob0_Sell15_Price50_GTB20_RO.OrderId,
				constants.Order_Alice_Num1_Id2_Clob0_Buy20_Price30_GTB20_RO.OrderId,
			},
		},
		`Can place a regular taker order that matches with multiple maker orders and the maker
						change position sides after the first match, and the maker's resting reduce-only orders on
						that CLOB are canceled`: {
			placedMatchableOrders: []types.MatchableOrder{
				&constants.Order_Bob_Num0_Id8_Clob0_Sell20_Price10_GTB22,
				&constants.Order_Bob_Num0_Id13_Clob0_Sell35_Price35_GTB30,
				&constants.Order_Bob_Num0_Id2_Clob0_Sell10_Price35_GTB20_RO,
				&constants.Order_Bob_Num0_Id1_Clob0_Sell15_Price50_GTB20_RO,
			},
			collateralizationCheckFailures: map[int]map[satypes.SubaccountId]satypes.UpdateResult{},
			statePositionSizes: map[types.ClobPairId]map[satypes.SubaccountId]*big.Int{
				0: {
					constants.Alice_Num1: big.NewInt(0),
					constants.Bob_Num0:   big.NewInt(15),
				},
			},

			order: constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25,

			expectedOrderStatus:   types.Success,
			expectedFilledSize:    30,
			expectedRemainingBids: []OrderWithRemainingSize{},
			expectedRemainingAsks: []OrderWithRemainingSize{

				{
					Order:         constants.Order_Bob_Num0_Id13_Clob0_Sell35_Price35_GTB30,
					RemainingSize: 25,
				},
			},
			expectedPendingMatches: []expectedMatch{
				{
					makerOrder:      &constants.Order_Bob_Num0_Id8_Clob0_Sell20_Price10_GTB22,
					takerOrder:      &constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25,
					matchedQuantums: 20,
				},
				{
					makerOrder:      &constants.Order_Bob_Num0_Id13_Clob0_Sell35_Price35_GTB30,
					takerOrder:      &constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25,
					matchedQuantums: 10,
				},
			},
			expectedExistingMatches: []expectedMatch{},
			expectedNewMatches: []expectedMatch{
				{
					makerOrder:      &constants.Order_Bob_Num0_Id8_Clob0_Sell20_Price10_GTB22,
					takerOrder:      &constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25,
					matchedQuantums: 20,
				},
				{
					makerOrder:      &constants.Order_Bob_Num0_Id13_Clob0_Sell35_Price35_GTB30,
					takerOrder:      &constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25,
					matchedQuantums: 10,
				},
			},
			expectedSubaccountOpenReduceOnlyOrders: map[types.ClobPairId]map[satypes.SubaccountId]map[types.OrderId]bool{
				0: {},
			},
			expectedCancelledReduceOnlyOrders: []types.OrderId{
				constants.Order_Bob_Num0_Id1_Clob0_Sell15_Price50_GTB20_RO.OrderId,
				constants.Order_Bob_Num0_Id2_Clob0_Sell10_Price35_GTB20_RO.OrderId,
			},
		},
		`Can place a reduce-only taker order that matches and sets the taker's position size to zero,
						which cancels the remaining size of the taker order and the taker's resting reduce-only orders`: {
			placedMatchableOrders: []types.MatchableOrder{
				&constants.Order_Bob_Num0_Id8_Clob0_Sell20_Price10_GTB22,
				&constants.Order_Alice_Num1_Id6_Clob0_Buy10_Price5_GTB20_RO,
			},
			collateralizationCheckFailures: map[int]map[satypes.SubaccountId]satypes.UpdateResult{},
			statePositionSizes: map[types.ClobPairId]map[satypes.SubaccountId]*big.Int{
				0: {
					constants.Alice_Num1: big.NewInt(-15),
					constants.Bob_Num0:   big.NewInt(0),
				},
			},

			order: constants.Order_Alice_Num1_Id2_Clob0_Buy20_Price30_GTB20_RO,

			expectedOrderStatus:   types.ReduceOnlyResized,
			expectedFilledSize:    15,
			expectedRemainingBids: []OrderWithRemainingSize{},
			expectedRemainingAsks: []OrderWithRemainingSize{

				{
					Order:         constants.Order_Bob_Num0_Id8_Clob0_Sell20_Price10_GTB22,
					RemainingSize: 5,
				},
			},
			expectedPendingMatches: []expectedMatch{
				{
					makerOrder:      &constants.Order_Bob_Num0_Id8_Clob0_Sell20_Price10_GTB22,
					takerOrder:      &constants.Order_Alice_Num1_Id2_Clob0_Buy20_Price30_GTB20_RO,
					matchedQuantums: 15,
				},
			},
			expectedExistingMatches: []expectedMatch{},
			expectedNewMatches: []expectedMatch{
				{
					makerOrder:      &constants.Order_Bob_Num0_Id8_Clob0_Sell20_Price10_GTB22,
					takerOrder:      &constants.Order_Alice_Num1_Id2_Clob0_Buy20_Price30_GTB20_RO,
					matchedQuantums: 15,
				},
			},
			expectedSubaccountOpenReduceOnlyOrders: map[types.ClobPairId]map[satypes.SubaccountId]map[types.OrderId]bool{
				0: {},
			},
			expectedCancelledReduceOnlyOrders: []types.OrderId{
				constants.Order_Alice_Num1_Id6_Clob0_Buy10_Price5_GTB20_RO.OrderId,
				constants.Order_Alice_Num1_Id2_Clob0_Buy20_Price30_GTB20_RO.OrderId,
			},
		},
		`Can place a reduce-only taker order that matches multiple times and sets the taker's position
						size to zero, which cancels the remaining size of the taker order and the taker's resting
						reduce-only orders. More matchable orders remain on the book but matching is stopped when
		    the taker's position size reaches 0`: {
			placedMatchableOrders: []types.MatchableOrder{
				&constants.Order_Alice_Num0_Id1_Clob0_Sell5_Price15_GTB15,
				&constants.Order_Alice_Num0_Id7_Clob0_Sell25_Price15_GTB20,
				&constants.Order_Alice_Num0_Id10_Clob0_Sell25_Price15_GTB20,
				&constants.Order_Alice_Num1_Id6_Clob0_Buy10_Price5_GTB20_RO,
			},
			collateralizationCheckFailures: map[int]map[satypes.SubaccountId]satypes.UpdateResult{},
			statePositionSizes: map[types.ClobPairId]map[satypes.SubaccountId]*big.Int{
				0: {
					constants.Alice_Num0: big.NewInt(0),
					constants.Alice_Num1: big.NewInt(-15),
				},
			},

			order: constants.Order_Alice_Num1_Id2_Clob0_Buy20_Price30_GTB20_RO,

			expectedOrderStatus:   types.ReduceOnlyResized,
			expectedFilledSize:    15,
			expectedRemainingBids: []OrderWithRemainingSize{},
			expectedRemainingAsks: []OrderWithRemainingSize{

				{
					Order:         constants.Order_Alice_Num0_Id7_Clob0_Sell25_Price15_GTB20,
					RemainingSize: 15,
				},
				{
					Order:         constants.Order_Alice_Num0_Id10_Clob0_Sell25_Price15_GTB20,
					RemainingSize: 25,
				},
			},
			expectedPendingMatches: []expectedMatch{
				{
					makerOrder:      &constants.Order_Alice_Num0_Id1_Clob0_Sell5_Price15_GTB15,
					takerOrder:      &constants.Order_Alice_Num1_Id2_Clob0_Buy20_Price30_GTB20_RO,
					matchedQuantums: 5,
				},
				{
					makerOrder:      &constants.Order_Alice_Num0_Id7_Clob0_Sell25_Price15_GTB20,
					takerOrder:      &constants.Order_Alice_Num1_Id2_Clob0_Buy20_Price30_GTB20_RO,
					matchedQuantums: 10,
				},
			},
			expectedExistingMatches: []expectedMatch{},
			expectedNewMatches: []expectedMatch{
				{
					makerOrder:      &constants.Order_Alice_Num0_Id1_Clob0_Sell5_Price15_GTB15,
					takerOrder:      &constants.Order_Alice_Num1_Id2_Clob0_Buy20_Price30_GTB20_RO,
					matchedQuantums: 5,
				},
				{
					makerOrder:      &constants.Order_Alice_Num0_Id7_Clob0_Sell25_Price15_GTB20,
					takerOrder:      &constants.Order_Alice_Num1_Id2_Clob0_Buy20_Price30_GTB20_RO,
					matchedQuantums: 10,
				},
			},
			expectedSubaccountOpenReduceOnlyOrders: map[types.ClobPairId]map[satypes.SubaccountId]map[types.OrderId]bool{
				0: {},
			},
			expectedCancelledReduceOnlyOrders: []types.OrderId{
				constants.Order_Alice_Num1_Id6_Clob0_Buy10_Price5_GTB20_RO.OrderId,
				constants.Order_Alice_Num1_Id2_Clob0_Buy20_Price30_GTB20_RO.OrderId,
			},
		},
		`Can place a taker order that matches with a reduce-only maker order and sets the maker's
				        position size to zero, which cancels all remaining reduce-only maker orders`: {
			placedMatchableOrders: []types.MatchableOrder{
				&constants.Order_Alice_Num1_Id2_Clob0_Buy20_Price30_GTB20_RO,
				&constants.Order_Alice_Num1_Id6_Clob0_Buy10_Price5_GTB20_RO,
			},
			collateralizationCheckFailures: map[int]map[satypes.SubaccountId]satypes.UpdateResult{},
			statePositionSizes: map[types.ClobPairId]map[satypes.SubaccountId]*big.Int{
				0: {
					constants.Alice_Num1: big.NewInt(-15),
					constants.Bob_Num0:   big.NewInt(0),
				},
			},

			order: constants.Order_Bob_Num0_Id8_Clob0_Sell20_Price10_GTB22,

			expectedOrderStatus:   types.Success,
			expectedFilledSize:    15,
			expectedRemainingBids: []OrderWithRemainingSize{},
			expectedRemainingAsks: []OrderWithRemainingSize{

				{
					Order:         constants.Order_Bob_Num0_Id8_Clob0_Sell20_Price10_GTB22,
					RemainingSize: 5,
				},
			},
			expectedPendingMatches: []expectedMatch{
				{
					makerOrder:      &constants.Order_Alice_Num1_Id2_Clob0_Buy20_Price30_GTB20_RO,
					takerOrder:      &constants.Order_Bob_Num0_Id8_Clob0_Sell20_Price10_GTB22,
					matchedQuantums: 15,
				},
			},
			expectedExistingMatches: []expectedMatch{},
			expectedNewMatches: []expectedMatch{
				{
					makerOrder:      &constants.Order_Alice_Num1_Id2_Clob0_Buy20_Price30_GTB20_RO,
					takerOrder:      &constants.Order_Bob_Num0_Id8_Clob0_Sell20_Price10_GTB22,
					matchedQuantums: 15,
				},
			},
			expectedSubaccountOpenReduceOnlyOrders: map[types.ClobPairId]map[satypes.SubaccountId]map[types.OrderId]bool{
				0: {},
			},
			expectedCancelledReduceOnlyOrders: []types.OrderId{
				constants.Order_Alice_Num1_Id2_Clob0_Buy20_Price30_GTB20_RO.OrderId,
				constants.Order_Alice_Num1_Id6_Clob0_Buy10_Price5_GTB20_RO.OrderId,
			},
		},
		// TODO(DEC-1415): Uncomment reduce-only tests after the patch is removed.
		// `Can place a taker order that matches with a maker order that sets the maker's position size
		// 	        to zero, and all following matched reduce-only orders from the same maker are canceled`: {
		// 	placedMatchableOrders: []types.MatchableOrder{
		// 		&constants.Order_Alice_Num1_Id9_Clob0_Sell10_Price10_GTB31,
		// 		&constants.Order_Alice_Num1_Id1_Clob0_Sell10_Price15_GTB20_RO,
		// 		&constants.Order_Alice_Num1_Id4_Clob0_Sell15_Price20_GTB20_RO,
		// 		&constants.Order_Alice_Num0_Id1_Clob0_Sell5_Price15_GTB15,
		// 	},
		// 	collateralizationCheckFailures: map[int]map[satypes.SubaccountId]satypes.UpdateResult{},
		// 	statePositionSizes: map[types.ClobPairId]map[satypes.SubaccountId]*big.Int{
		// 		0: {
		// 			constants.Alice_Num0: big.NewInt(0),
		// 			constants.Alice_Num1: big.NewInt(10),
		// 			constants.Bob_Num0: big.NewInt(0),
		// 		},
		// 	},

		// 	order: constants.Order_Bob_Num0_Id6_Clob0_Buy20_Price1000_GTB22,

		// 	expectedOrderStatus: types.Success,
		// 	expectedFilledSize:  15,
		// 	expectedRemainingBids: []orderWithRemainingSize{
		// 		{
		// 			order:         constants.Order_Bob_Num0_Id6_Clob0_Buy20_Price1000_GTB22,
		// 			remainingSize: 5,
		// 		},
		// 	},
		// 	expectedRemainingAsks: []orderWithRemainingSize{},
		// 	expectedPendingMatches: []expectedMatch{
		// 		{
		// 			makerOrder:      &constants.Order_Alice_Num1_Id9_Clob0_Sell10_Price10_GTB31,
		// 			takerOrder:      &constants.Order_Bob_Num0_Id6_Clob0_Buy20_Price1000_GTB22,
		// 			matchedQuantums: 10,
		// 		},
		// 		{
		// 			makerOrder:      &constants.Order_Alice_Num0_Id1_Clob0_Sell5_Price15_GTB15,
		// 			takerOrder:      &constants.Order_Bob_Num0_Id6_Clob0_Buy20_Price1000_GTB22,
		// 			matchedQuantums: 5,
		// 		},
		// 	},
		// 	expectedExistingMatches: []expectedMatch{},
		// 	expectedNewMatches: []expectedMatch{
		// 		{
		// 			makerOrder:      &constants.Order_Alice_Num1_Id9_Clob0_Sell10_Price10_GTB31,
		// 			takerOrder:      &constants.Order_Bob_Num0_Id6_Clob0_Buy20_Price1000_GTB22,
		// 			matchedQuantums: 10,
		// 		},
		// 		{
		// 			makerOrder:      &constants.Order_Alice_Num0_Id1_Clob0_Sell5_Price15_GTB15,
		// 			takerOrder:      &constants.Order_Bob_Num0_Id6_Clob0_Buy20_Price1000_GTB22,
		// 			matchedQuantums: 5,
		// 		},
		// 	},
		// 	expectedSubaccountOpenReduceOnlyOrders: map[types.ClobPairId]map[satypes.SubaccountId]map[types.OrderId]bool{
		// 		0: {},
		// 	},
		// 	expectedCancelledReduceOnlyOrders: []types.OrderId{
		// 		constants.Order_Alice_Num1_Id1_Clob0_Sell10_Price15_GTB20_RO.OrderId,
		// 		constants.Order_Alice_Num1_Id4_Clob0_Sell15_Price20_GTB20_RO.OrderId,
		// 	},
		// },
		// `Can place a taker order that matches with a maker order that changes the maker's position
		// 		        side, and all following matched reduce-only orders from the same maker are canceled`: {
		// 	placedMatchableOrders: []types.MatchableOrder{
		// 		&constants.Order_Alice_Num1_Id9_Clob0_Sell10_Price10_GTB31,
		// 		&constants.Order_Alice_Num1_Id1_Clob0_Sell10_Price15_GTB20_RO,
		// 		&constants.Order_Alice_Num1_Id4_Clob0_Sell15_Price20_GTB20_RO,
		// 		&constants.Order_Alice_Num0_Id1_Clob0_Sell5_Price15_GTB15,
		// 	},
		// 	collateralizationCheckFailures: map[int]map[satypes.SubaccountId]satypes.UpdateResult{},
		// 	statePositionSizes: map[types.ClobPairId]map[satypes.SubaccountId]*big.Int{
		// 		0: {
		// 			constants.Alice_Num0: big.NewInt(0),
		// 			constants.Alice_Num1: big.NewInt(5),
		// 			constants.Bob_Num0: big.NewInt(0),
		// 		},
		// 	},

		// 	order: constants.Order_Bob_Num0_Id6_Clob0_Buy20_Price1000_GTB22,

		// 	expectedOrderStatus: types.Success,
		// 	expectedFilledSize:  15,
		// 	expectedRemainingBids: []orderWithRemainingSize{
		// 		{
		// 			order:         constants.Order_Bob_Num0_Id6_Clob0_Buy20_Price1000_GTB22,
		// 			remainingSize: 5,
		// 		},
		// 	},
		// 	expectedRemainingAsks: []orderWithRemainingSize{},
		// 	expectedPendingMatches: []expectedMatch{
		// 		{
		// 			makerOrder:      &constants.Order_Alice_Num1_Id9_Clob0_Sell10_Price10_GTB31,
		// 			takerOrder:      &constants.Order_Bob_Num0_Id6_Clob0_Buy20_Price1000_GTB22,
		// 			matchedQuantums: 10,
		// 		},
		// 		{
		// 			makerOrder:      &constants.Order_Alice_Num0_Id1_Clob0_Sell5_Price15_GTB15,
		// 			takerOrder:      &constants.Order_Bob_Num0_Id6_Clob0_Buy20_Price1000_GTB22,
		// 			matchedQuantums: 5,
		// 		},
		// 	},
		// 	expectedExistingMatches: []expectedMatch{},
		// 	expectedNewMatches: []expectedMatch{
		// 		{
		// 			makerOrder:      &constants.Order_Alice_Num1_Id9_Clob0_Sell10_Price10_GTB31,
		// 			takerOrder:      &constants.Order_Bob_Num0_Id6_Clob0_Buy20_Price1000_GTB22,
		// 			matchedQuantums: 10,
		// 		},
		// 		{
		// 			makerOrder:      &constants.Order_Alice_Num0_Id1_Clob0_Sell5_Price15_GTB15,
		// 			takerOrder:      &constants.Order_Bob_Num0_Id6_Clob0_Buy20_Price1000_GTB22,
		// 			matchedQuantums: 5,
		// 		},
		// 	},
		// 	expectedSubaccountOpenReduceOnlyOrders: map[types.ClobPairId]map[satypes.SubaccountId]map[types.OrderId]bool{
		// 		0: {},
		// 	},
		// 	expectedCancelledReduceOnlyOrders: []types.OrderId{
		// 		constants.Order_Alice_Num1_Id1_Clob0_Sell10_Price15_GTB20_RO.OrderId,
		// 		constants.Order_Alice_Num1_Id4_Clob0_Sell15_Price20_GTB20_RO.OrderId,
		// 	},
		// },
		`Can place a reduce-only taker order that matches with a reduce-only maker order, and both
						reduce-only orders change their subaccount's position side and need to be resized. Since the
						maker order is resized by a larger delta than the taker, the taker order's full size is
						placed on the book`: {
			placedMatchableOrders: []types.MatchableOrder{
				&constants.Order_Bob_Num0_Id3_Clob0_Sell20_Price10_GTB20_RO,
			},
			collateralizationCheckFailures: map[int]map[satypes.SubaccountId]satypes.UpdateResult{},
			statePositionSizes: map[types.ClobPairId]map[satypes.SubaccountId]*big.Int{
				0: {
					constants.Alice_Num1: big.NewInt(-15),
					constants.Bob_Num0:   big.NewInt(10),
				},
			},

			order: constants.Order_Alice_Num1_Id2_Clob0_Buy20_Price30_GTB20_RO,

			expectedOrderStatus: types.Success,
			expectedFilledSize:  10,
			expectedRemainingBids: []OrderWithRemainingSize{
				{
					Order:         constants.Order_Alice_Num1_Id2_Clob0_Buy20_Price30_GTB20_RO,
					RemainingSize: 10,
				},
			},
			expectedRemainingAsks: []OrderWithRemainingSize{},
			expectedPendingMatches: []expectedMatch{
				{
					makerOrder:      &constants.Order_Bob_Num0_Id3_Clob0_Sell20_Price10_GTB20_RO,
					takerOrder:      &constants.Order_Alice_Num1_Id2_Clob0_Buy20_Price30_GTB20_RO,
					matchedQuantums: 10,
				},
			},
			expectedExistingMatches: []expectedMatch{},
			expectedNewMatches: []expectedMatch{
				{
					makerOrder:      &constants.Order_Bob_Num0_Id3_Clob0_Sell20_Price10_GTB20_RO,
					takerOrder:      &constants.Order_Alice_Num1_Id2_Clob0_Buy20_Price30_GTB20_RO,
					matchedQuantums: 10,
				},
			},
			expectedSubaccountOpenReduceOnlyOrders: map[types.ClobPairId]map[satypes.SubaccountId]map[types.OrderId]bool{
				0: {
					constants.Alice_Num1: {
						constants.Order_Alice_Num1_Id2_Clob0_Buy20_Price30_GTB20_RO.OrderId: true,
					},
				},
			},
			expectedCancelledReduceOnlyOrders: []types.OrderId{
				constants.Order_Bob_Num0_Id3_Clob0_Sell20_Price10_GTB20_RO.OrderId,
			},
		},
		`Can place a reduce-only taker order that matches with a reduce-only maker order, and both
						reduce-only orders change their subaccount's position side and need to be resized. Since the
						taker order is resized by a larger delta than the maker, the maker order is partially filled
						and is not resized down`: {
			placedMatchableOrders: []types.MatchableOrder{
				&constants.Order_Bob_Num0_Id3_Clob0_Sell20_Price10_GTB20_RO,
			},
			collateralizationCheckFailures: map[int]map[satypes.SubaccountId]satypes.UpdateResult{},
			statePositionSizes: map[types.ClobPairId]map[satypes.SubaccountId]*big.Int{
				0: {
					constants.Alice_Num1: big.NewInt(-10),
					constants.Bob_Num0:   big.NewInt(15),
				},
			},

			order: constants.Order_Alice_Num1_Id2_Clob0_Buy20_Price30_GTB20_RO,

			expectedOrderStatus:   types.ReduceOnlyResized,
			expectedFilledSize:    10,
			expectedRemainingBids: []OrderWithRemainingSize{},
			expectedRemainingAsks: []OrderWithRemainingSize{

				{
					Order:         constants.Order_Bob_Num0_Id3_Clob0_Sell20_Price10_GTB20_RO,
					RemainingSize: 10,
				},
			},
			expectedPendingMatches: []expectedMatch{
				{
					makerOrder:      &constants.Order_Bob_Num0_Id3_Clob0_Sell20_Price10_GTB20_RO,
					takerOrder:      &constants.Order_Alice_Num1_Id2_Clob0_Buy20_Price30_GTB20_RO,
					matchedQuantums: 10,
				},
			},
			expectedExistingMatches: []expectedMatch{},
			expectedNewMatches: []expectedMatch{
				{
					makerOrder:      &constants.Order_Bob_Num0_Id3_Clob0_Sell20_Price10_GTB20_RO,
					takerOrder:      &constants.Order_Alice_Num1_Id2_Clob0_Buy20_Price30_GTB20_RO,
					matchedQuantums: 10,
				},
			},
			expectedSubaccountOpenReduceOnlyOrders: map[types.ClobPairId]map[satypes.SubaccountId]map[types.OrderId]bool{
				0: {
					constants.Bob_Num0: {
						constants.Order_Bob_Num0_Id3_Clob0_Sell20_Price10_GTB20_RO.OrderId: true,
					},
				},
			},
			expectedCancelledReduceOnlyOrders: []types.OrderId{
				constants.Order_Alice_Num1_Id2_Clob0_Buy20_Price30_GTB20_RO.OrderId,
			},
		},
		`Can place a reduce-only taker order that matches with a reduce-only maker order, and both
						reduce-only orders change their subaccount's position side and need to be resized. Since the
						taker order is resized by the same delta as the maker, both orders are resized and
						removed from the book`: {
			placedMatchableOrders: []types.MatchableOrder{
				&constants.Order_Bob_Num0_Id3_Clob0_Sell20_Price10_GTB20_RO,
			},
			collateralizationCheckFailures: map[int]map[satypes.SubaccountId]satypes.UpdateResult{},
			statePositionSizes: map[types.ClobPairId]map[satypes.SubaccountId]*big.Int{
				0: {
					constants.Alice_Num1: big.NewInt(-10),
					constants.Bob_Num0:   big.NewInt(10),
				},
			},

			order: constants.Order_Alice_Num1_Id2_Clob0_Buy20_Price30_GTB20_RO,

			expectedOrderStatus:   types.ReduceOnlyResized,
			expectedFilledSize:    10,
			expectedRemainingBids: []OrderWithRemainingSize{},
			expectedRemainingAsks: []OrderWithRemainingSize{},
			expectedPendingMatches: []expectedMatch{
				{
					makerOrder:      &constants.Order_Bob_Num0_Id3_Clob0_Sell20_Price10_GTB20_RO,
					takerOrder:      &constants.Order_Alice_Num1_Id2_Clob0_Buy20_Price30_GTB20_RO,
					matchedQuantums: 10,
				},
			},
			expectedExistingMatches: []expectedMatch{},
			expectedNewMatches: []expectedMatch{
				{
					makerOrder:      &constants.Order_Bob_Num0_Id3_Clob0_Sell20_Price10_GTB20_RO,
					takerOrder:      &constants.Order_Alice_Num1_Id2_Clob0_Buy20_Price30_GTB20_RO,
					matchedQuantums: 10,
				},
			},
			expectedSubaccountOpenReduceOnlyOrders: map[types.ClobPairId]map[satypes.SubaccountId]map[types.OrderId]bool{
				0: {},
			},
			expectedCancelledReduceOnlyOrders: []types.OrderId{
				constants.Order_Bob_Num0_Id3_Clob0_Sell20_Price10_GTB20_RO.OrderId,
				constants.Order_Alice_Num1_Id2_Clob0_Buy20_Price30_GTB20_RO.OrderId,
			},
		},
		`Can place a reduce-only taker order that matches with a reduce-only maker order, and both
						reduce-only orders change their subaccount's position side and need to be resized. However
						since the taker order fails collateralization checks, the full size of the maker order remains
						on the book`: {
			placedMatchableOrders: []types.MatchableOrder{
				&constants.Order_Bob_Num0_Id3_Clob0_Sell20_Price10_GTB20_RO,
			},
			collateralizationCheckFailures: map[int]map[satypes.SubaccountId]satypes.UpdateResult{
				0: {
					constants.Alice_Num1: satypes.NewlyUndercollateralized,
				},
			},
			statePositionSizes: map[types.ClobPairId]map[satypes.SubaccountId]*big.Int{
				0: {
					constants.Alice_Num1: big.NewInt(-10),
					constants.Bob_Num0:   big.NewInt(10),
				},
			},

			order: constants.Order_Alice_Num1_Id2_Clob0_Buy20_Price30_GTB20_RO,

			expectedOrderStatus:   types.Undercollateralized,
			expectedFilledSize:    0,
			expectedRemainingBids: []OrderWithRemainingSize{},
			expectedRemainingAsks: []OrderWithRemainingSize{
				{
					Order:         constants.Order_Bob_Num0_Id3_Clob0_Sell20_Price10_GTB20_RO,
					RemainingSize: 20,
				},
			},
			expectedPendingMatches: []expectedMatch{
				{
					makerOrder:      &constants.Order_Bob_Num0_Id3_Clob0_Sell20_Price10_GTB20_RO,
					takerOrder:      &constants.Order_Alice_Num1_Id2_Clob0_Buy20_Price30_GTB20_RO,
					matchedQuantums: 10,
				},
			},
			expectedExistingMatches: []expectedMatch{},
			expectedNewMatches:      []expectedMatch{},
			expectedSubaccountOpenReduceOnlyOrders: map[types.ClobPairId]map[satypes.SubaccountId]map[types.OrderId]bool{
				0: {
					constants.Bob_Num0: {
						constants.Order_Bob_Num0_Id3_Clob0_Sell20_Price10_GTB20_RO.OrderId: true,
					},
				},
			},
			expectedCancelledReduceOnlyOrders: []types.OrderId{},
		},
		`Can place a reduce-only taker order that matches with a reduce-only maker order, and both
						reduce-only orders change their subaccount's position side and need to be resized. However
						since the maker order fails collateralization checks, the full size of the taker order is
						added to the book`: {
			placedMatchableOrders: []types.MatchableOrder{
				&constants.Order_Bob_Num0_Id3_Clob0_Sell20_Price10_GTB20_RO,
			},
			collateralizationCheckFailures: map[int]map[satypes.SubaccountId]satypes.UpdateResult{
				0: {
					constants.Bob_Num0: satypes.NewlyUndercollateralized,
				},
			},
			statePositionSizes: map[types.ClobPairId]map[satypes.SubaccountId]*big.Int{
				0: {
					constants.Alice_Num1: big.NewInt(-10),
					constants.Bob_Num0:   big.NewInt(10),
				},
			},

			order: constants.Order_Alice_Num1_Id2_Clob0_Buy20_Price30_GTB20_RO,

			expectedOrderStatus: types.Success,
			expectedFilledSize:  0,
			expectedRemainingBids: []OrderWithRemainingSize{
				{
					Order:         constants.Order_Alice_Num1_Id2_Clob0_Buy20_Price30_GTB20_RO,
					RemainingSize: 20,
				},
			},
			expectedRemainingAsks: []OrderWithRemainingSize{},
			expectedPendingMatches: []expectedMatch{
				{
					makerOrder:      &constants.Order_Bob_Num0_Id3_Clob0_Sell20_Price10_GTB20_RO,
					takerOrder:      &constants.Order_Alice_Num1_Id2_Clob0_Buy20_Price30_GTB20_RO,
					matchedQuantums: 10,
				},
			},
			expectedExistingMatches: []expectedMatch{},
			expectedNewMatches:      []expectedMatch{},
			expectedSubaccountOpenReduceOnlyOrders: map[types.ClobPairId]map[satypes.SubaccountId]map[types.OrderId]bool{
				0: {
					constants.Alice_Num1: {
						constants.Order_Alice_Num1_Id2_Clob0_Buy20_Price30_GTB20_RO.OrderId: true,
					},
				},
			},
			expectedCancelledReduceOnlyOrders: []types.OrderId{},
		},
		`Can place a taker order that causes a maker reduce-only to fail collateralization checks,
			then matches with a regular maker order that changes the makers side, meaning the maker
			order will be canceled because it failed collateralization checks and not from the maker
			subaccount changing sides`: {
			placedMatchableOrders: []types.MatchableOrder{
				&constants.Order_Bob_Num0_Id3_Clob0_Sell20_Price10_GTB20_RO,
				&constants.Order_Bob_Num0_Id14_Clob0_Sell10_Price10_GTB25,
			},
			collateralizationCheckFailures: map[int]map[satypes.SubaccountId]satypes.UpdateResult{
				0: {
					constants.Bob_Num0: satypes.NewlyUndercollateralized,
				},
			},
			statePositionSizes: map[types.ClobPairId]map[satypes.SubaccountId]*big.Int{
				0: {
					constants.Alice_Num1: big.NewInt(0),
					constants.Bob_Num0:   big.NewInt(5),
				},
			},

			order: constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25,

			expectedOrderStatus: types.Success,
			expectedFilledSize:  10,
			expectedRemainingBids: []OrderWithRemainingSize{
				{
					Order:         constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25,
					RemainingSize: 20,
				},
			},
			expectedRemainingAsks: []OrderWithRemainingSize{},
			expectedPendingMatches: []expectedMatch{
				{
					makerOrder:      &constants.Order_Bob_Num0_Id3_Clob0_Sell20_Price10_GTB20_RO,
					takerOrder:      &constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25,
					matchedQuantums: 5,
				},
				{
					makerOrder:      &constants.Order_Bob_Num0_Id14_Clob0_Sell10_Price10_GTB25,
					takerOrder:      &constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25,
					matchedQuantums: 10,
				},
			},
			expectedExistingMatches: []expectedMatch{},
			expectedNewMatches: []expectedMatch{
				{
					makerOrder:      &constants.Order_Bob_Num0_Id14_Clob0_Sell10_Price10_GTB25,
					takerOrder:      &constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25,
					matchedQuantums: 10,
				},
			},
			expectedSubaccountOpenReduceOnlyOrders: map[types.ClobPairId]map[satypes.SubaccountId]map[types.OrderId]bool{
				0: {},
			},
			expectedCancelledReduceOnlyOrders: []types.OrderId{},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup memclob state and test expectations.
			require.NotNil(t, tc.expectedSubaccountOpenReduceOnlyOrders)
			order := tc.order
			addOrderToOrderbookSize := satypes.BaseQuantums(0)
			expectedFilledSize := tc.expectedFilledSize
			if tc.expectedErr == nil && tc.expectedOrderStatus.IsSuccess() {
				addOrderToOrderbookSize = order.GetBaseQuantums() - expectedFilledSize
			}

			getStatePosition := func(
				subaccountId satypes.SubaccountId,
				clobPairId types.ClobPairId,
			) (
				statePositionSize *big.Int,
			) {
				clobStatePositionSizes, exists := tc.statePositionSizes[clobPairId]
				require.True(
					t,
					exists,
					"Expected CLOB pair ID %d to exist in statePositionSizes",
					clobPairId,
				)

				statePositionSize, exists = clobStatePositionSizes[subaccountId]
				require.True(
					t,
					exists,
					"Expected subaccount ID %v to exist in statePositionSizes",
					subaccountId,
				)

				return statePositionSize
			}

			memclob, fakeMemClobKeeper, expectedNumCollateralizationChecks, numCollateralChecks := placeOrderTestSetup(
				t,
				ctx,
				tc.placedMatchableOrders,
				&order,
				tc.expectedPendingMatches,
				tc.expectedOrderStatus,
				addOrderToOrderbookSize,
				tc.expectedErr,
				tc.collateralizationCheckFailures,
				getStatePosition,
			)

			// Run the test case and verify expectations.
			offchainUpdates := placeOrderAndVerifyExpectations(
				t,
				ctx,
				memclob,
				order,
				numCollateralChecks,
				expectedFilledSize,
				expectedFilledSize,
				tc.expectedOrderStatus,
				tc.expectedErr,
				expectedNumCollateralizationChecks,
				tc.expectedRemainingBids,
				tc.expectedRemainingAsks,
				append(tc.expectedExistingMatches, tc.expectedNewMatches...),
				fakeMemClobKeeper,
			)

			// Verify that the expected reduce-only orders remain on each orderbook.
			for clobPairId, expectedOpenReduceOnlyOrders := range tc.expectedSubaccountOpenReduceOnlyOrders {
				orderbook, exists := memclob.orderbooks[clobPairId]
				require.True(
					t,
					exists,
					"Expected orderbook with CLOB pair ID %d to exist in memclob",
					clobPairId,
				)
				require.Equal(
					t,
					expectedOpenReduceOnlyOrders,
					orderbook.SubaccountOpenReduceOnlyOrders,
				)
			}

			// Verify the correct offchain update messages were returned.
			assertPlaceOrderOffchainMessages(
				t,
				ctx,
				offchainUpdates,
				order,
				tc.placedMatchableOrders,
				tc.collateralizationCheckFailures,
				tc.expectedErr,
				expectedFilledSize,
				tc.expectedOrderStatus,
				tc.expectedExistingMatches,
				tc.expectedNewMatches,
				tc.expectedCancelledReduceOnlyOrders,
				// TODO(IND-261): Add tests for replaced reduce-only orders.
				false,
			)
		})
	}
}

// These tests aim to test two different scenarios regarding stateful reduce-only order removals:
//  1. A stateful reduce-only maker order is encountered during matching which would result in increasing position size.
//     This should result in the maker order being removed and an OrderRemoval operation being added to the ops queue.
//     This happens DURING the matching loop, in `mustPerformTakerOrderMatching`.
//  2. A taker order matches with a a maker order that changes the maker subaccount's position side.
//     The maker subaccount's resting stateful reduce-only orders should be removed and an OrderRemoval operation
//     should be added to the ops queue for each resting order.
//     This happens AFTER the matching loop, in `maybeCancelReduceOnlyOrders`.
func TestPlaceOrder_LongTermReduceOnlyRemovals(t *testing.T) {
	ctx, _, _ := sdktest.NewSdkContextWithMultistore()
	ctx = ctx.WithIsCheckTx(true)
	tests := map[string]struct {
		// State.
		placedMatchableOrders          []types.MatchableOrder
		collateralizationCheckFailures map[int]map[satypes.SubaccountId]satypes.UpdateResult
		statePositionSizes             map[types.ClobPairId]map[satypes.SubaccountId]*big.Int

		// Parameters.
		order types.Order

		// Expectations.
		expectedOrderStatus                    types.OrderStatus
		expectedFilledSize                     satypes.BaseQuantums
		expectedErr                            error
		expectedPendingMatches                 []expectedMatch
		expectedExistingMatches                []expectedMatch
		expectedNewMatches                     []expectedMatch
		expectedRemainingBids                  []OrderWithRemainingSize
		expectedRemainingAsks                  []OrderWithRemainingSize
		expectedSubaccountOpenReduceOnlyOrders map[types.ClobPairId]map[satypes.SubaccountId]map[types.OrderId]bool
		expectedCancelledReduceOnlyOrders      []types.OrderId
		expectedInternalOperations             []types.InternalOperation
	}{
		`Can place a regular taker order that partially matches with a regular order which changes the
			maker order's position side. The maker's resting stateful reduce-only order is canceled and an
			OrderRemoval operation is added while attempting to match remaining taker size against it`: {
			placedMatchableOrders: []types.MatchableOrder{
				&constants.Order_Bob_Num0_Id8_Clob0_Sell20_Price10_GTB22,
				&constants.LongTermOrder_Bob_Num0_Id2_Clob0_Sell10_Price35_GTB20_RO,
			},
			collateralizationCheckFailures: map[int]map[satypes.SubaccountId]satypes.UpdateResult{},
			statePositionSizes: map[types.ClobPairId]map[satypes.SubaccountId]*big.Int{
				0: {
					constants.Alice_Num1: big.NewInt(0),
					constants.Bob_Num0:   big.NewInt(15),
				},
			},

			order: constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25,

			expectedOrderStatus: types.Success,
			expectedFilledSize:  20,
			expectedRemainingBids: []OrderWithRemainingSize{
				{
					Order:         constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25,
					RemainingSize: 10,
				},
			},
			expectedRemainingAsks: []OrderWithRemainingSize{},
			expectedPendingMatches: []expectedMatch{
				{
					makerOrder:      &constants.Order_Bob_Num0_Id8_Clob0_Sell20_Price10_GTB22,
					takerOrder:      &constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25,
					matchedQuantums: 20,
				},
			},
			expectedExistingMatches: []expectedMatch{},
			expectedNewMatches: []expectedMatch{
				{
					makerOrder:      &constants.Order_Bob_Num0_Id8_Clob0_Sell20_Price10_GTB22,
					takerOrder:      &constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25,
					matchedQuantums: 20,
				},
			},
			expectedSubaccountOpenReduceOnlyOrders: map[types.ClobPairId]map[satypes.SubaccountId]map[types.OrderId]bool{
				0: {},
			},
			expectedCancelledReduceOnlyOrders: []types.OrderId{
				constants.LongTermOrder_Bob_Num0_Id2_Clob0_Sell10_Price35_GTB20_RO.OrderId,
			},
			expectedInternalOperations: []types.InternalOperation{
				types.NewOrderRemovalInternalOperation(
					constants.LongTermOrder_Bob_Num0_Id2_Clob0_Sell10_Price35_GTB20_RO.OrderId,
					types.OrderRemoval_REMOVAL_REASON_INVALID_REDUCE_ONLY,
				),
				types.NewShortTermOrderPlacementInternalOperation(
					constants.Order_Bob_Num0_Id8_Clob0_Sell20_Price10_GTB22,
				),
				types.NewShortTermOrderPlacementInternalOperation(
					constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25,
				),
				types.NewMatchOrdersInternalOperation(
					constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Bob_Num0_Id8_Clob0_Sell20_Price10_GTB22.OrderId,
							FillAmount:   20,
						},
					},
				),
			},
		},
		`Can place a regular taker order that fully-matches with a regular order which changes the
			maker order's position side. The maker's resting stateful reduce-only order is canceled by 
			maybeCancelReduceOnlyOrders after matching finishes. An Order Removal operation is added
			to the operations queue`: {
			placedMatchableOrders: []types.MatchableOrder{
				&constants.Order_Bob_Num0_Id13_Clob0_Sell35_Price35_GTB30,
				&constants.LongTermOrder_Bob_Num0_Id2_Clob0_Sell10_Price35_GTB20_RO,
			},
			collateralizationCheckFailures: map[int]map[satypes.SubaccountId]satypes.UpdateResult{},
			statePositionSizes: map[types.ClobPairId]map[satypes.SubaccountId]*big.Int{
				0: {
					constants.Alice_Num1: big.NewInt(0),
					constants.Bob_Num0:   big.NewInt(15),
				},
			},

			order: constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25,

			expectedOrderStatus:   types.Success,
			expectedFilledSize:    30,
			expectedRemainingBids: []OrderWithRemainingSize{},
			expectedRemainingAsks: []OrderWithRemainingSize{
				{
					Order:         constants.Order_Bob_Num0_Id13_Clob0_Sell35_Price35_GTB30,
					RemainingSize: 5,
				},
			},
			expectedPendingMatches: []expectedMatch{
				{
					makerOrder:      &constants.Order_Bob_Num0_Id13_Clob0_Sell35_Price35_GTB30,
					takerOrder:      &constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25,
					matchedQuantums: 30,
				},
			},
			expectedExistingMatches: []expectedMatch{},
			expectedNewMatches: []expectedMatch{
				{
					makerOrder:      &constants.Order_Bob_Num0_Id13_Clob0_Sell35_Price35_GTB30,
					takerOrder:      &constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25,
					matchedQuantums: 30,
				},
			},
			expectedSubaccountOpenReduceOnlyOrders: map[types.ClobPairId]map[satypes.SubaccountId]map[types.OrderId]bool{
				0: {},
			},
			expectedCancelledReduceOnlyOrders: []types.OrderId{
				constants.LongTermOrder_Bob_Num0_Id2_Clob0_Sell10_Price35_GTB20_RO.OrderId,
			},
			expectedInternalOperations: []types.InternalOperation{
				types.NewShortTermOrderPlacementInternalOperation(
					constants.Order_Bob_Num0_Id13_Clob0_Sell35_Price35_GTB30,
				),
				types.NewShortTermOrderPlacementInternalOperation(
					constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25,
				),
				types.NewMatchOrdersInternalOperation(
					constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Bob_Num0_Id13_Clob0_Sell35_Price35_GTB30.OrderId,
							FillAmount:   30,
						},
					},
				),
				types.NewOrderRemovalInternalOperation(
					constants.LongTermOrder_Bob_Num0_Id2_Clob0_Sell10_Price35_GTB20_RO.OrderId,
					types.OrderRemoval_REMOVAL_REASON_INVALID_REDUCE_ONLY,
				),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup memclob state and test expectations.
			require.NotNil(t, tc.expectedSubaccountOpenReduceOnlyOrders)
			order := tc.order
			addOrderToOrderbookSize := satypes.BaseQuantums(0)
			expectedFilledSize := tc.expectedFilledSize
			if tc.expectedErr == nil && tc.expectedOrderStatus.IsSuccess() {
				addOrderToOrderbookSize = order.GetBaseQuantums() - expectedFilledSize
			}

			getStatePosition := func(
				subaccountId satypes.SubaccountId,
				clobPairId types.ClobPairId,
			) (
				statePositionSize *big.Int,
			) {
				clobStatePositionSizes, exists := tc.statePositionSizes[clobPairId]
				require.True(
					t,
					exists,
					"Expected CLOB pair ID %d to exist in statePositionSizes",
					clobPairId,
				)

				statePositionSize, exists = clobStatePositionSizes[subaccountId]
				require.True(
					t,
					exists,
					"Expected subaccount ID %v to exist in statePositionSizes",
					subaccountId,
				)

				return statePositionSize
			}

			memclob, fakeMemClobKeeper, expectedNumCollateralizationChecks, numCollateralChecks := placeOrderTestSetup(
				t,
				ctx,
				tc.placedMatchableOrders,
				&order,
				tc.expectedPendingMatches,
				tc.expectedOrderStatus,
				addOrderToOrderbookSize,
				tc.expectedErr,
				tc.collateralizationCheckFailures,
				getStatePosition,
			)

			// Run the test case and verify expectations.
			offchainUpdates := placeOrderAndVerifyExpectations(
				t,
				ctx,
				memclob,
				order,
				numCollateralChecks,
				expectedFilledSize,
				expectedFilledSize,
				tc.expectedOrderStatus,
				tc.expectedErr,
				expectedNumCollateralizationChecks,
				tc.expectedRemainingBids,
				tc.expectedRemainingAsks,
				append(tc.expectedExistingMatches, tc.expectedNewMatches...),
				fakeMemClobKeeper,
			)

			require.Equal( // asserting internal operations to ensure OrderRemovals are included
				t,
				tc.expectedInternalOperations,
				memclob.operationsToPropose.OperationsQueue,
			)

			// Verify that the expected reduce-only orders remain on each orderbook.
			for clobPairId, expectedOpenReduceOnlyOrders := range tc.expectedSubaccountOpenReduceOnlyOrders {
				orderbook, exists := memclob.orderbooks[clobPairId]
				require.True(
					t,
					exists,
					"Expected orderbook with CLOB pair ID %d to exist in memclob",
					clobPairId,
				)
				require.Equal(
					t,
					expectedOpenReduceOnlyOrders,
					orderbook.SubaccountOpenReduceOnlyOrders,
				)
			}

			// Verify the correct offchain update messages were returned.
			assertPlaceOrderOffchainMessages(
				t,
				ctx,
				offchainUpdates,
				order,
				tc.placedMatchableOrders,
				tc.collateralizationCheckFailures,
				tc.expectedErr,
				expectedFilledSize,
				tc.expectedOrderStatus,
				tc.expectedExistingMatches,
				tc.expectedNewMatches,
				tc.expectedCancelledReduceOnlyOrders,
				false,
			)
		})
	}
}
