package types_test

import (
	"sort"
	"testing"

	"github.com/dydxprotocol/v4/testutil/constants"
	"github.com/dydxprotocol/v4/x/clob/types"

	"github.com/stretchr/testify/require"
)

func TestSortStatefulOrderPlacements(t *testing.T) {
	tests := map[string]struct {
		// Parameters.
		statefulOrderPlacements         []types.StatefulOrderPlacement
		expectedStatefulOrderPlacements []types.StatefulOrderPlacement
	}{
		`sorts stateful order placements with different block heights and same transaction index in
        ascending order`: {
			statefulOrderPlacements: []types.StatefulOrderPlacement{
				{
					Order:            constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20,
					BlockHeight:      5,
					TransactionIndex: 5,
				},
				{
					Order:            constants.LongTermOrder_Alice_Num0_Id1_Clob1_Sell65_Price15_GTBT25,
					BlockHeight:      7,
					TransactionIndex: 5,
				},
				{
					Order:            constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTB15,
					BlockHeight:      6,
					TransactionIndex: 5,
				},
			},
			expectedStatefulOrderPlacements: []types.StatefulOrderPlacement{
				{
					Order:            constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20,
					BlockHeight:      5,
					TransactionIndex: 5,
				},
				{
					Order:            constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTB15,
					BlockHeight:      6,
					TransactionIndex: 5,
				},
				{
					Order:            constants.LongTermOrder_Alice_Num0_Id1_Clob1_Sell65_Price15_GTBT25,
					BlockHeight:      7,
					TransactionIndex: 5,
				},
			},
		},
		`sorts stateful order placements with same block heights and different transaction indices in
        ascending order`: {
			statefulOrderPlacements: []types.StatefulOrderPlacement{
				{
					Order:            constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20,
					BlockHeight:      5,
					TransactionIndex: 5,
				},
				{
					Order:            constants.LongTermOrder_Alice_Num0_Id1_Clob1_Sell65_Price15_GTBT25,
					BlockHeight:      5,
					TransactionIndex: 7,
				},
				{
					Order:            constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTB15,
					BlockHeight:      5,
					TransactionIndex: 6,
				},
			},
			expectedStatefulOrderPlacements: []types.StatefulOrderPlacement{
				{
					Order:            constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20,
					BlockHeight:      5,
					TransactionIndex: 5,
				},
				{
					Order:            constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTB15,
					BlockHeight:      5,
					TransactionIndex: 6,
				},
				{
					Order:            constants.LongTermOrder_Alice_Num0_Id1_Clob1_Sell65_Price15_GTBT25,
					BlockHeight:      5,
					TransactionIndex: 7,
				},
			},
		},
		`sorts stateful order placements with different block heights and different transaction indices in
        ascending order`: {
			statefulOrderPlacements: []types.StatefulOrderPlacement{
				{
					Order:            constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
					BlockHeight:      4,
					TransactionIndex: 8,
				},
				{
					Order:            constants.LongTermOrder_Alice_Num0_Id1_Clob1_Sell65_Price15_GTBT25,
					BlockHeight:      8,
					TransactionIndex: 2,
				},
				{
					Order:            constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20,
					BlockHeight:      4,
					TransactionIndex: 7,
				},
				{
					Order:            constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTB15,
					BlockHeight:      8,
					TransactionIndex: 3,
				},
			},
			expectedStatefulOrderPlacements: []types.StatefulOrderPlacement{
				{
					Order:            constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20,
					BlockHeight:      4,
					TransactionIndex: 7,
				},
				{
					Order:            constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
					BlockHeight:      4,
					TransactionIndex: 8,
				},
				{
					Order:            constants.LongTermOrder_Alice_Num0_Id1_Clob1_Sell65_Price15_GTBT25,
					BlockHeight:      8,
					TransactionIndex: 2,
				},
				{
					Order:            constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTB15,
					BlockHeight:      8,
					TransactionIndex: 3,
				},
			},
		},
		`sorts empty list of stateful order placements`: {
			statefulOrderPlacements:         []types.StatefulOrderPlacement{},
			expectedStatefulOrderPlacements: []types.StatefulOrderPlacement{},
		},
		`sorts list of stateful order placements with one element`: {
			statefulOrderPlacements: []types.StatefulOrderPlacement{
				{
					Order:            constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
					BlockHeight:      4,
					TransactionIndex: 8,
				},
			},
			expectedStatefulOrderPlacements: []types.StatefulOrderPlacement{
				{
					Order:            constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
					BlockHeight:      4,
					TransactionIndex: 8,
				},
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			statefulOrderPlacements := tc.statefulOrderPlacements
			sort.Sort(types.SortedStatefulOrderPlacements(statefulOrderPlacements))
			require.Equal(t, tc.expectedStatefulOrderPlacements, statefulOrderPlacements)
		})
	}
}

func TestSortStatefulOrderPlacements_PanicsOnDuplicateBlockHeightAndTransactionIndex(t *testing.T) {
	statefulOrderPlacementsDuplicate := []types.StatefulOrderPlacement{
		{
			Order:            constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
			BlockHeight:      4,
			TransactionIndex: 8,
		},
		{
			Order:            constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
			BlockHeight:      4,
			TransactionIndex: 8,
		},
	}
	require.Panics(
		t,
		func() {
			sort.Sort(types.SortedStatefulOrderPlacements(statefulOrderPlacementsDuplicate))
		},
	)
}
