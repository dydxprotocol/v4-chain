package types_test

import (
	"sort"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	"github.com/stretchr/testify/require"
)

func TestSortLongTermOrderPlacements(t *testing.T) {
	tests := map[string]struct {
		// Parameters.
		longTermOrderPlacements         []types.StatefulOrderPlacement
		expectedLongTermOrderPlacements []types.StatefulOrderPlacement
	}{
		`sorts long term order placements with different block heights and same transaction index in
        ascending order`: {
			longTermOrderPlacements: []types.StatefulOrderPlacement{
				&types.LongTermOrderPlacement{
					Order: constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB5,
					PlacementIndex: types.TransactionOrdering{
						BlockHeight:      5,
						TransactionIndex: 5,
					},
				},
				&types.LongTermOrderPlacement{
					Order: constants.LongTermOrder_Alice_Num0_Id1_Clob1_Sell65_Price15_GTBT25,
					PlacementIndex: types.TransactionOrdering{
						BlockHeight:      7,
						TransactionIndex: 5,
					},
				},
				&types.LongTermOrderPlacement{
					Order: constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10,
					PlacementIndex: types.TransactionOrdering{
						BlockHeight:      6,
						TransactionIndex: 5,
					},
				},
			},
			expectedLongTermOrderPlacements: []types.StatefulOrderPlacement{
				&types.LongTermOrderPlacement{
					Order: constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB5,
					PlacementIndex: types.TransactionOrdering{
						BlockHeight:      5,
						TransactionIndex: 5,
					},
				},
				&types.LongTermOrderPlacement{
					Order: constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10,
					PlacementIndex: types.TransactionOrdering{
						BlockHeight:      6,
						TransactionIndex: 5,
					},
				},
				&types.LongTermOrderPlacement{
					Order: constants.LongTermOrder_Alice_Num0_Id1_Clob1_Sell65_Price15_GTBT25,
					PlacementIndex: types.TransactionOrdering{
						BlockHeight:      7,
						TransactionIndex: 5,
					},
				},
			},
		},
		`sorts long term order placements with same block heights and different transaction indices in
        ascending order`: {
			longTermOrderPlacements: []types.StatefulOrderPlacement{
				&types.LongTermOrderPlacement{
					Order: constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB5,
					PlacementIndex: types.TransactionOrdering{
						BlockHeight:      5,
						TransactionIndex: 5,
					},
				},
				&types.LongTermOrderPlacement{
					Order: constants.LongTermOrder_Alice_Num0_Id1_Clob1_Sell65_Price15_GTBT25,
					PlacementIndex: types.TransactionOrdering{
						BlockHeight:      5,
						TransactionIndex: 7,
					},
				},
				&types.LongTermOrderPlacement{
					Order: constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10,
					PlacementIndex: types.TransactionOrdering{
						BlockHeight:      5,
						TransactionIndex: 6,
					},
				},
			},
			expectedLongTermOrderPlacements: []types.StatefulOrderPlacement{
				&types.LongTermOrderPlacement{
					Order: constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB5,
					PlacementIndex: types.TransactionOrdering{
						BlockHeight:      5,
						TransactionIndex: 5,
					},
				},
				&types.LongTermOrderPlacement{
					Order: constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10,
					PlacementIndex: types.TransactionOrdering{
						BlockHeight:      5,
						TransactionIndex: 6,
					},
				},
				&types.LongTermOrderPlacement{
					Order: constants.LongTermOrder_Alice_Num0_Id1_Clob1_Sell65_Price15_GTBT25,
					PlacementIndex: types.TransactionOrdering{
						BlockHeight:      5,
						TransactionIndex: 7,
					},
				},
			},
		},
		`sorts long term order placements with different block heights and different transaction indices in
        ascending order`: {
			longTermOrderPlacements: []types.StatefulOrderPlacement{
				&types.LongTermOrderPlacement{
					Order: constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB5,
					PlacementIndex: types.TransactionOrdering{
						BlockHeight:      4,
						TransactionIndex: 8,
					},
				},
				&types.LongTermOrderPlacement{
					Order: constants.LongTermOrder_Alice_Num0_Id1_Clob1_Sell65_Price15_GTBT25,
					PlacementIndex: types.TransactionOrdering{
						BlockHeight:      8,
						TransactionIndex: 2,
					},
				},
				&types.LongTermOrderPlacement{
					Order: constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20,
					PlacementIndex: types.TransactionOrdering{
						BlockHeight:      4,
						TransactionIndex: 7,
					},
				},
				&types.LongTermOrderPlacement{
					Order: constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10,
					PlacementIndex: types.TransactionOrdering{
						BlockHeight:      8,
						TransactionIndex: 3,
					},
				},
			},
			expectedLongTermOrderPlacements: []types.StatefulOrderPlacement{
				&types.LongTermOrderPlacement{
					Order: constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20,
					PlacementIndex: types.TransactionOrdering{
						BlockHeight:      4,
						TransactionIndex: 7,
					},
				},
				&types.LongTermOrderPlacement{
					Order: constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB5,
					PlacementIndex: types.TransactionOrdering{
						BlockHeight:      4,
						TransactionIndex: 8,
					},
				},
				&types.LongTermOrderPlacement{
					Order: constants.LongTermOrder_Alice_Num0_Id1_Clob1_Sell65_Price15_GTBT25,
					PlacementIndex: types.TransactionOrdering{
						BlockHeight:      8,
						TransactionIndex: 2,
					},
				},
				&types.LongTermOrderPlacement{
					Order: constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10,
					PlacementIndex: types.TransactionOrdering{
						BlockHeight:      8,
						TransactionIndex: 3,
					},
				},
			},
		},
		`sorts empty list of long term order placements`: {
			longTermOrderPlacements:         []types.StatefulOrderPlacement{},
			expectedLongTermOrderPlacements: []types.StatefulOrderPlacement{},
		},
		`sorts list of long term order placements with one element`: {
			longTermOrderPlacements: []types.StatefulOrderPlacement{
				&types.LongTermOrderPlacement{
					Order: constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10,
					PlacementIndex: types.TransactionOrdering{
						BlockHeight:      4,
						TransactionIndex: 8,
					},
				},
			},
			expectedLongTermOrderPlacements: []types.StatefulOrderPlacement{
				&types.LongTermOrderPlacement{
					Order: constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10,
					PlacementIndex: types.TransactionOrdering{
						BlockHeight:      4,
						TransactionIndex: 8,
					},
				},
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			statefulOrderPlacements := tc.longTermOrderPlacements
			sort.Sort(types.SortedStatefulOrderPlacement(statefulOrderPlacements))
			require.Equal(t, tc.expectedLongTermOrderPlacements, statefulOrderPlacements)
		})
	}
}

func TestSortLongTermOrderPlacements_PanicsOnDuplicateBlockHeightAndTransactionIndex(t *testing.T) {
	longTermOrderPlacementsDuplicate := []types.StatefulOrderPlacement{
		&types.LongTermOrderPlacement{
			Order: constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10,
			PlacementIndex: types.TransactionOrdering{
				BlockHeight:      4,
				TransactionIndex: 8,
			},
		},
		&types.LongTermOrderPlacement{
			Order: constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
			PlacementIndex: types.TransactionOrdering{
				BlockHeight:      4,
				TransactionIndex: 8,
			},
		},
	}
	require.Panics(
		t,
		func() {
			sort.Sort(types.SortedStatefulOrderPlacement(longTermOrderPlacementsDuplicate))
		},
	)
}

func TestSortConditionalOrderPlacements(t *testing.T) {
	tests := map[string]struct {
		// Parameters.
		conditionalOrderPlacements         []types.StatefulOrderPlacement
		expectedConditionalOrderPlacements []types.StatefulOrderPlacement
	}{
		`sorts conditional order placements with different block heights and same transaction index in
		ascending order`: {
			conditionalOrderPlacements: []types.StatefulOrderPlacement{
				&types.ConditionalOrderPlacement{
					Order: constants.ConditionalOrder_Alice_Num0_Id2_Clob0_Buy20_Price10_GTBT15_TakeProfit10,
					PlacementIndex: types.TransactionOrdering{
						BlockHeight:      5,
						TransactionIndex: 5,
					},
				},
				&types.ConditionalOrderPlacement{
					Order: constants.ConditionalOrder_Alice_Num0_Id1_Clob0_Buy15_Price10_GTBT15_StopLoss20,
					PlacementIndex: types.TransactionOrdering{
						BlockHeight:      7,
						TransactionIndex: 5,
					},
				},
				&types.ConditionalOrderPlacement{
					Order: constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_TakeProfit20,
					PlacementIndex: types.TransactionOrdering{
						BlockHeight:      6,
						TransactionIndex: 5,
					},
				},
			},
			expectedConditionalOrderPlacements: []types.StatefulOrderPlacement{
				&types.ConditionalOrderPlacement{
					Order: constants.ConditionalOrder_Alice_Num0_Id2_Clob0_Buy20_Price10_GTBT15_TakeProfit10,
					PlacementIndex: types.TransactionOrdering{
						BlockHeight:      5,
						TransactionIndex: 5,
					},
				},
				&types.ConditionalOrderPlacement{
					Order: constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_TakeProfit20,
					PlacementIndex: types.TransactionOrdering{
						BlockHeight:      6,
						TransactionIndex: 5,
					},
				},
				&types.ConditionalOrderPlacement{
					Order: constants.ConditionalOrder_Alice_Num0_Id1_Clob0_Buy15_Price10_GTBT15_StopLoss20,
					PlacementIndex: types.TransactionOrdering{
						BlockHeight:      7,
						TransactionIndex: 5,
					},
				},
			},
		},
		`sorts conditional order placements with same block heights and different transaction indices in
		ascending order`: {
			conditionalOrderPlacements: []types.StatefulOrderPlacement{
				&types.ConditionalOrderPlacement{
					Order: constants.ConditionalOrder_Alice_Num0_Id2_Clob0_Buy20_Price10_GTBT15_TakeProfit10,
					PlacementIndex: types.TransactionOrdering{
						BlockHeight:      5,
						TransactionIndex: 5,
					},
				},
				&types.ConditionalOrderPlacement{
					Order: constants.ConditionalOrder_Alice_Num0_Id1_Clob0_Buy15_Price10_GTBT15_StopLoss20,
					PlacementIndex: types.TransactionOrdering{
						BlockHeight:      5,
						TransactionIndex: 7,
					},
				},
				&types.ConditionalOrderPlacement{
					Order: constants.ConditionalOrder_Alice_Num0_Id1_Clob0_Buy15_Price10_GTBT15_StopLoss20,
					PlacementIndex: types.TransactionOrdering{
						BlockHeight:      5,
						TransactionIndex: 6,
					},
				},
			},
			expectedConditionalOrderPlacements: []types.StatefulOrderPlacement{
				&types.ConditionalOrderPlacement{
					Order: constants.ConditionalOrder_Alice_Num0_Id2_Clob0_Buy20_Price10_GTBT15_TakeProfit10,
					PlacementIndex: types.TransactionOrdering{
						BlockHeight:      5,
						TransactionIndex: 5,
					},
				},
				&types.ConditionalOrderPlacement{
					Order: constants.ConditionalOrder_Alice_Num0_Id1_Clob0_Buy15_Price10_GTBT15_StopLoss20,
					PlacementIndex: types.TransactionOrdering{
						BlockHeight:      5,
						TransactionIndex: 6,
					},
				},
				&types.ConditionalOrderPlacement{
					Order: constants.ConditionalOrder_Alice_Num0_Id1_Clob0_Buy15_Price10_GTBT15_StopLoss20,
					PlacementIndex: types.TransactionOrdering{
						BlockHeight:      5,
						TransactionIndex: 7,
					},
				},
			},
		},
		`sorts conditional order placements with different block heights and different transaction indices in
		ascending order`: {
			conditionalOrderPlacements: []types.StatefulOrderPlacement{
				&types.ConditionalOrderPlacement{
					Order: constants.ConditionalOrder_Alice_Num0_Id2_Clob0_Buy20_Price10_GTBT15_TakeProfit10,
					PlacementIndex: types.TransactionOrdering{
						BlockHeight:      4,
						TransactionIndex: 8,
					},
				},
				&types.ConditionalOrderPlacement{
					Order: constants.ConditionalOrder_Alice_Num0_Id1_Clob0_Buy15_Price10_GTBT15_StopLoss20,
					PlacementIndex: types.TransactionOrdering{
						BlockHeight:      8,
						TransactionIndex: 2,
					},
				},
				&types.ConditionalOrderPlacement{
					Order: constants.ConditionalOrder_Alice_Num0_Id1_Clob0_Buy15_Price10_GTBT15_StopLoss20,
					PlacementIndex: types.TransactionOrdering{
						BlockHeight:      4,
						TransactionIndex: 7,
					},
				},
				&types.ConditionalOrderPlacement{
					Order: constants.ConditionalOrder_Alice_Num0_Id1_Clob0_Buy15_Price10_GTBT15_TakeProfit20,
					PlacementIndex: types.TransactionOrdering{
						BlockHeight:      8,
						TransactionIndex: 3,
					},
				},
			},
			expectedConditionalOrderPlacements: []types.StatefulOrderPlacement{
				&types.ConditionalOrderPlacement{
					Order: constants.ConditionalOrder_Alice_Num0_Id1_Clob0_Buy15_Price10_GTBT15_StopLoss20,
					PlacementIndex: types.TransactionOrdering{
						BlockHeight:      4,
						TransactionIndex: 7,
					},
				},
				&types.ConditionalOrderPlacement{
					Order: constants.ConditionalOrder_Alice_Num0_Id2_Clob0_Buy20_Price10_GTBT15_TakeProfit10,
					PlacementIndex: types.TransactionOrdering{
						BlockHeight:      4,
						TransactionIndex: 8,
					},
				},
				&types.ConditionalOrderPlacement{
					Order: constants.ConditionalOrder_Alice_Num0_Id1_Clob0_Buy15_Price10_GTBT15_StopLoss20,
					PlacementIndex: types.TransactionOrdering{
						BlockHeight:      8,
						TransactionIndex: 2,
					},
				},
				&types.ConditionalOrderPlacement{
					Order: constants.ConditionalOrder_Alice_Num0_Id1_Clob0_Buy15_Price10_GTBT15_TakeProfit20,
					PlacementIndex: types.TransactionOrdering{
						BlockHeight:      8,
						TransactionIndex: 3,
					},
				},
			},
		},
		`sorts empty list of conditional order placements`: {
			conditionalOrderPlacements:         []types.StatefulOrderPlacement{},
			expectedConditionalOrderPlacements: []types.StatefulOrderPlacement{},
		},
		`sorts list of conditional order placements with one element`: {
			conditionalOrderPlacements: []types.StatefulOrderPlacement{
				&types.ConditionalOrderPlacement{
					Order: constants.ConditionalOrder_Alice_Num0_Id1_Clob0_Buy15_Price10_GTBT15_TakeProfit20,
					PlacementIndex: types.TransactionOrdering{
						BlockHeight:      4,
						TransactionIndex: 8,
					},
				},
			},
			expectedConditionalOrderPlacements: []types.StatefulOrderPlacement{
				&types.ConditionalOrderPlacement{
					Order: constants.ConditionalOrder_Alice_Num0_Id1_Clob0_Buy15_Price10_GTBT15_TakeProfit20,
					PlacementIndex: types.TransactionOrdering{
						BlockHeight:      4,
						TransactionIndex: 8,
					},
				},
			},
		},

		// Sorting tests with triggered conditional orders.
		`Same placement index with one triggered index`: {
			conditionalOrderPlacements: []types.StatefulOrderPlacement{
				&types.ConditionalOrderPlacement{
					Order: constants.ConditionalOrder_Alice_Num0_Id1_Clob0_Buy15_Price10_GTBT15_TakeProfit20,
					PlacementIndex: types.TransactionOrdering{
						BlockHeight:      4,
						TransactionIndex: 8,
					},
				},
				&types.ConditionalOrderPlacement{
					Order: constants.ConditionalOrder_Alice_Num0_Id1_Clob0_Buy15_Price10_GTBT15_TakeProfit20,
					PlacementIndex: types.TransactionOrdering{
						BlockHeight:      4,
						TransactionIndex: 7,
					},
					TriggerIndex: &types.TransactionOrdering{
						BlockHeight:      10,
						TransactionIndex: 2,
					},
				},
			},
			expectedConditionalOrderPlacements: []types.StatefulOrderPlacement{
				&types.ConditionalOrderPlacement{
					Order: constants.ConditionalOrder_Alice_Num0_Id1_Clob0_Buy15_Price10_GTBT15_TakeProfit20,
					PlacementIndex: types.TransactionOrdering{
						BlockHeight:      4,
						TransactionIndex: 8,
					},
				},
				&types.ConditionalOrderPlacement{
					Order: constants.ConditionalOrder_Alice_Num0_Id1_Clob0_Buy15_Price10_GTBT15_TakeProfit20,
					PlacementIndex: types.TransactionOrdering{
						BlockHeight:      4,
						TransactionIndex: 7,
					},
					TriggerIndex: &types.TransactionOrdering{
						BlockHeight:      10,
						TransactionIndex: 2,
					},
				},
			},
		},
		`Same placement index with one triggered index, opposite order`: {
			conditionalOrderPlacements: []types.StatefulOrderPlacement{
				&types.ConditionalOrderPlacement{
					Order: constants.ConditionalOrder_Alice_Num0_Id1_Clob0_Buy15_Price10_GTBT15_TakeProfit20,
					PlacementIndex: types.TransactionOrdering{
						BlockHeight:      4,
						TransactionIndex: 8,
					},
					TriggerIndex: &types.TransactionOrdering{
						BlockHeight:      10,
						TransactionIndex: 2,
					},
				},
				&types.ConditionalOrderPlacement{
					Order: constants.ConditionalOrder_Alice_Num0_Id1_Clob0_Buy15_Price10_GTBT15_TakeProfit20,
					PlacementIndex: types.TransactionOrdering{
						BlockHeight:      4,
						TransactionIndex: 9,
					},
				},
			},
			expectedConditionalOrderPlacements: []types.StatefulOrderPlacement{
				&types.ConditionalOrderPlacement{
					Order: constants.ConditionalOrder_Alice_Num0_Id1_Clob0_Buy15_Price10_GTBT15_TakeProfit20,
					PlacementIndex: types.TransactionOrdering{
						BlockHeight:      4,
						TransactionIndex: 9,
					},
				},
				&types.ConditionalOrderPlacement{
					Order: constants.ConditionalOrder_Alice_Num0_Id1_Clob0_Buy15_Price10_GTBT15_TakeProfit20,
					PlacementIndex: types.TransactionOrdering{
						BlockHeight:      4,
						TransactionIndex: 8,
					},
					TriggerIndex: &types.TransactionOrdering{
						BlockHeight:      10,
						TransactionIndex: 2,
					},
				},
			},
		},
		`sort triggered orders based on triggered transaction ordering`: {
			conditionalOrderPlacements: []types.StatefulOrderPlacement{
				&types.ConditionalOrderPlacement{
					Order: constants.ConditionalOrder_Alice_Num0_Id1_Clob0_Buy15_Price10_GTBT15_TakeProfit20,
					PlacementIndex: types.TransactionOrdering{
						BlockHeight:      4,
						TransactionIndex: 8,
					},
					TriggerIndex: &types.TransactionOrdering{
						BlockHeight:      10,
						TransactionIndex: 2,
					},
				},
				&types.ConditionalOrderPlacement{
					Order: constants.ConditionalOrder_Alice_Num0_Id1_Clob0_Buy15_Price10_GTBT15_TakeProfit20,
					PlacementIndex: types.TransactionOrdering{
						BlockHeight:      4,
						TransactionIndex: 7,
					},
					TriggerIndex: &types.TransactionOrdering{
						BlockHeight:      10,
						TransactionIndex: 3,
					},
				},
				&types.ConditionalOrderPlacement{
					Order: constants.ConditionalOrder_Alice_Num0_Id2_Clob0_Buy20_Price10_GTBT15_StopLoss20,
					PlacementIndex: types.TransactionOrdering{
						BlockHeight:      4,
						TransactionIndex: 6,
					},
					TriggerIndex: &types.TransactionOrdering{
						BlockHeight:      5,
						TransactionIndex: 1,
					},
				},
				&types.ConditionalOrderPlacement{
					Order: constants.ConditionalOrder_Alice_Num0_Id2_Clob0_Buy20_Price10_GTBT15_TakeProfit10,
					PlacementIndex: types.TransactionOrdering{
						BlockHeight:      1,
						TransactionIndex: 2,
					},
					TriggerIndex: &types.TransactionOrdering{
						BlockHeight:      6,
						TransactionIndex: 1,
					},
				},
			},
			expectedConditionalOrderPlacements: []types.StatefulOrderPlacement{
				&types.ConditionalOrderPlacement{
					Order: constants.ConditionalOrder_Alice_Num0_Id2_Clob0_Buy20_Price10_GTBT15_StopLoss20,
					PlacementIndex: types.TransactionOrdering{
						BlockHeight:      4,
						TransactionIndex: 6,
					},
					TriggerIndex: &types.TransactionOrdering{
						BlockHeight:      5,
						TransactionIndex: 1,
					},
				},
				&types.ConditionalOrderPlacement{
					Order: constants.ConditionalOrder_Alice_Num0_Id2_Clob0_Buy20_Price10_GTBT15_TakeProfit10,
					PlacementIndex: types.TransactionOrdering{
						BlockHeight:      1,
						TransactionIndex: 2,
					},
					TriggerIndex: &types.TransactionOrdering{
						BlockHeight:      6,
						TransactionIndex: 1,
					},
				},
				&types.ConditionalOrderPlacement{
					Order: constants.ConditionalOrder_Alice_Num0_Id1_Clob0_Buy15_Price10_GTBT15_TakeProfit20,
					PlacementIndex: types.TransactionOrdering{
						BlockHeight:      4,
						TransactionIndex: 8,
					},
					TriggerIndex: &types.TransactionOrdering{
						BlockHeight:      10,
						TransactionIndex: 2,
					},
				},
				&types.ConditionalOrderPlacement{
					Order: constants.ConditionalOrder_Alice_Num0_Id1_Clob0_Buy15_Price10_GTBT15_TakeProfit20,
					PlacementIndex: types.TransactionOrdering{
						BlockHeight:      4,
						TransactionIndex: 7,
					},
					TriggerIndex: &types.TransactionOrdering{
						BlockHeight:      10,
						TransactionIndex: 3,
					},
				},
			},
		},
		`same placement index, sort differing trigger index based on block height and transaction index with
		untriggered orders`: {
			conditionalOrderPlacements: []types.StatefulOrderPlacement{
				&types.ConditionalOrderPlacement{
					Order: constants.ConditionalOrder_Alice_Num0_Id1_Clob0_Buy15_Price10_GTBT15_TakeProfit20,
					PlacementIndex: types.TransactionOrdering{
						BlockHeight:      4,
						TransactionIndex: 7,
					},
					TriggerIndex: &types.TransactionOrdering{
						BlockHeight:      10,
						TransactionIndex: 2,
					},
				},
				&types.ConditionalOrderPlacement{
					Order: constants.ConditionalOrder_Alice_Num0_Id1_Clob0_Buy15_Price10_GTBT15_TakeProfit20,
					PlacementIndex: types.TransactionOrdering{
						BlockHeight:      5,
						TransactionIndex: 7,
					},
				},
				&types.ConditionalOrderPlacement{
					Order: constants.ConditionalOrder_Alice_Num0_Id1_Clob0_Buy15_Price10_GTBT15_TakeProfit20,
					PlacementIndex: types.TransactionOrdering{
						BlockHeight:      5,
						TransactionIndex: 6,
					},
					TriggerIndex: &types.TransactionOrdering{
						BlockHeight:      7,
						TransactionIndex: 7,
					},
				},
				&types.ConditionalOrderPlacement{
					Order: constants.ConditionalOrder_Alice_Num0_Id1_Clob0_Buy15_Price10_GTBT15_TakeProfit20,
					PlacementIndex: types.TransactionOrdering{
						BlockHeight:      4,
						TransactionIndex: 20,
					},
					TriggerIndex: &types.TransactionOrdering{
						BlockHeight:      10,
						TransactionIndex: 3,
					},
				},
				&types.ConditionalOrderPlacement{
					Order: constants.ConditionalOrder_Alice_Num0_Id2_Clob0_Buy20_Price10_GTBT15_StopLoss20,
					PlacementIndex: types.TransactionOrdering{
						BlockHeight:      4,
						TransactionIndex: 6,
					},
					TriggerIndex: &types.TransactionOrdering{
						BlockHeight:      5,
						TransactionIndex: 1,
					},
				},
				&types.ConditionalOrderPlacement{
					Order: constants.ConditionalOrder_Alice_Num0_Id1_Clob0_Buy15_Price10_GTBT15_TakeProfit20,
					PlacementIndex: types.TransactionOrdering{
						BlockHeight:      4,
						TransactionIndex: 1000,
					},
				},
				&types.ConditionalOrderPlacement{
					Order: constants.ConditionalOrder_Alice_Num0_Id2_Clob0_Buy20_Price10_GTBT15_TakeProfit10,
					PlacementIndex: types.TransactionOrdering{
						BlockHeight:      4,
						TransactionIndex: 8,
					},
					TriggerIndex: &types.TransactionOrdering{
						BlockHeight:      6,
						TransactionIndex: 1,
					},
				},
			},
			expectedConditionalOrderPlacements: []types.StatefulOrderPlacement{
				&types.ConditionalOrderPlacement{
					Order: constants.ConditionalOrder_Alice_Num0_Id1_Clob0_Buy15_Price10_GTBT15_TakeProfit20,
					PlacementIndex: types.TransactionOrdering{
						BlockHeight:      4,
						TransactionIndex: 1000,
					},
				},
				&types.ConditionalOrderPlacement{
					Order: constants.ConditionalOrder_Alice_Num0_Id2_Clob0_Buy20_Price10_GTBT15_StopLoss20,
					PlacementIndex: types.TransactionOrdering{
						BlockHeight:      4,
						TransactionIndex: 6,
					},
					TriggerIndex: &types.TransactionOrdering{
						BlockHeight:      5,
						TransactionIndex: 1,
					},
				},
				&types.ConditionalOrderPlacement{
					Order: constants.ConditionalOrder_Alice_Num0_Id1_Clob0_Buy15_Price10_GTBT15_TakeProfit20,
					PlacementIndex: types.TransactionOrdering{
						BlockHeight:      5,
						TransactionIndex: 7,
					},
				},
				&types.ConditionalOrderPlacement{
					Order: constants.ConditionalOrder_Alice_Num0_Id2_Clob0_Buy20_Price10_GTBT15_TakeProfit10,
					PlacementIndex: types.TransactionOrdering{
						BlockHeight:      4,
						TransactionIndex: 8,
					},
					TriggerIndex: &types.TransactionOrdering{
						BlockHeight:      6,
						TransactionIndex: 1,
					},
				},
				&types.ConditionalOrderPlacement{
					Order: constants.ConditionalOrder_Alice_Num0_Id1_Clob0_Buy15_Price10_GTBT15_TakeProfit20,
					PlacementIndex: types.TransactionOrdering{
						BlockHeight:      5,
						TransactionIndex: 6,
					},
					TriggerIndex: &types.TransactionOrdering{
						BlockHeight:      7,
						TransactionIndex: 7,
					},
				},
				&types.ConditionalOrderPlacement{
					Order: constants.ConditionalOrder_Alice_Num0_Id1_Clob0_Buy15_Price10_GTBT15_TakeProfit20,
					PlacementIndex: types.TransactionOrdering{
						BlockHeight:      4,
						TransactionIndex: 7,
					},
					TriggerIndex: &types.TransactionOrdering{
						BlockHeight:      10,
						TransactionIndex: 2,
					},
				},
				&types.ConditionalOrderPlacement{
					Order: constants.ConditionalOrder_Alice_Num0_Id1_Clob0_Buy15_Price10_GTBT15_TakeProfit20,
					PlacementIndex: types.TransactionOrdering{
						BlockHeight:      4,
						TransactionIndex: 20,
					},
					TriggerIndex: &types.TransactionOrdering{
						BlockHeight:      10,
						TransactionIndex: 3,
					},
				},
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			conditionalOrderPlacements := tc.conditionalOrderPlacements
			sort.Sort(types.SortedStatefulOrderPlacement(conditionalOrderPlacements))
			require.Equal(t, tc.expectedConditionalOrderPlacements, conditionalOrderPlacements)
		})
	}
}

func TestSortConditionalOrderPlacements_PanicsOnDuplicateBlockHeightAndTransactionIndex(t *testing.T) {
	tests := map[string]struct {
		conditionalOrderPlacementsDuplicate []types.StatefulOrderPlacement
	}{
		"two untriggered conditional orders": {
			conditionalOrderPlacementsDuplicate: []types.StatefulOrderPlacement{
				&types.ConditionalOrderPlacement{
					Order: constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10,
					PlacementIndex: types.TransactionOrdering{
						BlockHeight:      4,
						TransactionIndex: 8,
					},
				},
				&types.ConditionalOrderPlacement{
					Order: constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
					PlacementIndex: types.TransactionOrdering{
						BlockHeight:      4,
						TransactionIndex: 8,
					},
				},
			},
		},
		"two triggered conditional orders": {
			conditionalOrderPlacementsDuplicate: []types.StatefulOrderPlacement{
				&types.ConditionalOrderPlacement{
					Order: constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10,
					PlacementIndex: types.TransactionOrdering{
						BlockHeight:      4,
						TransactionIndex: 8,
					},
					TriggerIndex: &types.TransactionOrdering{
						BlockHeight:      5,
						TransactionIndex: 10,
					},
				},
				&types.ConditionalOrderPlacement{
					Order: constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
					PlacementIndex: types.TransactionOrdering{
						BlockHeight:      4,
						TransactionIndex: 8,
					},
					TriggerIndex: &types.TransactionOrdering{
						BlockHeight:      5,
						TransactionIndex: 10,
					},
				},
			},
		},
		"one triggered, one untriggered": {
			conditionalOrderPlacementsDuplicate: []types.StatefulOrderPlacement{
				&types.ConditionalOrderPlacement{
					Order: constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10,
					PlacementIndex: types.TransactionOrdering{
						BlockHeight:      4,
						TransactionIndex: 8,
					},
					TriggerIndex: &types.TransactionOrdering{
						BlockHeight:      5,
						TransactionIndex: 10,
					},
				},
				&types.ConditionalOrderPlacement{
					Order: constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
					PlacementIndex: types.TransactionOrdering{
						BlockHeight:      5,
						TransactionIndex: 10,
					},
				},
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			require.Panics(
				t,
				func() {
					sort.Sort(types.SortedStatefulOrderPlacement(tc.conditionalOrderPlacementsDuplicate))
				},
			)
		})
	}
}
