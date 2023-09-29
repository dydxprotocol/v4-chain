package types_test

import (
	"sort"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"

	"github.com/stretchr/testify/require"
)

func TestMustMarshal(t *testing.T) {
	// Success
	b, _ := constants.Alice_Num0.Marshal()
	require.Equal(t, b, constants.Alice_Num0.MustMarshal())

	// Panic
	var sa types.SubaccountId
	require.Panics(t, func() { sa.MustMarshal() })
}

func TestSortSubaccountIds(t *testing.T) {
	tests := map[string]struct {
		// Parameters.
		ids         []types.SubaccountId
		expectedIds []types.SubaccountId
	}{
		"sorts with different owners": {
			ids: []types.SubaccountId{
				constants.Alice_Num0,
				constants.Bob_Num0,
			},
			expectedIds: []types.SubaccountId{
				constants.Bob_Num0,
				constants.Alice_Num0,
			},
		},
		"sorts with same owner different number": {
			ids: []types.SubaccountId{
				constants.Alice_Num0,
				constants.Alice_Num1,
			},
			expectedIds: []types.SubaccountId{
				constants.Alice_Num0,
				constants.Alice_Num1,
			},
		},
		"sorts with same owner and number": {
			ids: []types.SubaccountId{
				constants.Alice_Num0,
				constants.Alice_Num0,
			},
			expectedIds: []types.SubaccountId{
				constants.Alice_Num0,
				constants.Alice_Num0,
			},
		},
		"sorts with one subaccountId": {
			ids: []types.SubaccountId{
				constants.Alice_Num0,
			},
			expectedIds: []types.SubaccountId{
				constants.Alice_Num0,
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ids := tc.ids

			sort.Sort(types.SortedSubaccountIds(ids))

			require.Equal(t, tc.expectedIds, ids)
		})
	}
}
