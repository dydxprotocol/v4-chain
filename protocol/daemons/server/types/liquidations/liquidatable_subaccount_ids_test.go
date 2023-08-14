package types_test

import (
	"testing"

	liquidationstypes "github.com/dydxprotocol/v4-chain/protocol/daemons/server/types/liquidations"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/require"
)

func TestNewLiquidatableSubaccountIds(t *testing.T) {
	ls := liquidationstypes.NewLiquidatableSubaccountIds()
	require.Empty(t, ls.GetSubaccountIds())
}

func TestLiquidatableSubaccountIds_Multiple_Reads(t *testing.T) {
	ls := liquidationstypes.NewLiquidatableSubaccountIds()
	require.Empty(t, ls.GetSubaccountIds())

	expectedSubaccountIds := []satypes.SubaccountId{
		constants.Alice_Num1,
		constants.Bob_Num0,
	}
	ls.UpdateSubaccountIds(expectedSubaccountIds)
	require.Equal(t, expectedSubaccountIds, ls.GetSubaccountIds())
	require.Equal(t, expectedSubaccountIds, ls.GetSubaccountIds())
	require.Equal(t, expectedSubaccountIds, ls.GetSubaccountIds())
}

func TestLiquidatableSubaccountIds_Multiple_Writes(t *testing.T) {
	ls := liquidationstypes.NewLiquidatableSubaccountIds()
	require.Empty(t, ls.GetSubaccountIds())

	expectedSubaccountIds := []satypes.SubaccountId{
		constants.Alice_Num1,
	}
	ls.UpdateSubaccountIds(expectedSubaccountIds)
	require.Equal(t, expectedSubaccountIds, ls.GetSubaccountIds())

	expectedSubaccountIds = []satypes.SubaccountId{
		constants.Bob_Num0,
	}
	ls.UpdateSubaccountIds(expectedSubaccountIds)
	require.Equal(t, expectedSubaccountIds, ls.GetSubaccountIds())

	expectedSubaccountIds = []satypes.SubaccountId{
		constants.Carl_Num0,
	}
	ls.UpdateSubaccountIds(expectedSubaccountIds)
	require.Equal(t, expectedSubaccountIds, ls.GetSubaccountIds())
}

func TestLiquidatableSubaccountIds_Empty_Update(t *testing.T) {
	ls := liquidationstypes.NewLiquidatableSubaccountIds()
	require.Empty(t, ls.GetSubaccountIds())

	expectedSubaccountIds := []satypes.SubaccountId{
		constants.Alice_Num1,
	}
	ls.UpdateSubaccountIds(expectedSubaccountIds)
	require.Equal(t, expectedSubaccountIds, ls.GetSubaccountIds())

	expectedSubaccountIds = []satypes.SubaccountId{}
	ls.UpdateSubaccountIds(expectedSubaccountIds)
	require.Empty(t, ls.GetSubaccountIds())
}
