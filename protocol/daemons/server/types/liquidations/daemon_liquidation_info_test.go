package types_test

import (
	"testing"

	liquidationstypes "github.com/dydxprotocol/v4-chain/protocol/daemons/server/types/liquidations"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/require"
)

func TestNewDaemonLiquidationInfo(t *testing.T) {
	ls := liquidationstypes.NewDaemonLiquidationInfo()
	require.Empty(t, ls.GetLiquidatableSubaccountIds())
	require.Empty(t, ls.GetNegativeTncSubaccountIds())
	require.Empty(t, ls.GetSubaccountsWithOpenPositions(0))
}

func TestLiquidatableSubaccountIds_Multiple_Reads(t *testing.T) {
	ls := liquidationstypes.NewDaemonLiquidationInfo()
	require.Empty(t, ls.GetLiquidatableSubaccountIds())

	expectedSubaccountIds := []satypes.SubaccountId{
		constants.Alice_Num1,
		constants.Bob_Num0,
	}
	ls.UpdateLiquidatableSubaccountIds(expectedSubaccountIds, 1)
	require.Equal(t, expectedSubaccountIds, ls.GetLiquidatableSubaccountIds())
	require.Equal(t, expectedSubaccountIds, ls.GetLiquidatableSubaccountIds())
	require.Equal(t, expectedSubaccountIds, ls.GetLiquidatableSubaccountIds())
}

func TestNegativeTncSubaccounts_Multiple_Reads(t *testing.T) {
	ls := liquidationstypes.NewDaemonLiquidationInfo()
	require.Empty(t, ls.GetNegativeTncSubaccountIds())

	expectedSubaccountIds := []satypes.SubaccountId{
		constants.Alice_Num1,
		constants.Bob_Num0,
	}
	ls.UpdateNegativeTncSubaccountIds(expectedSubaccountIds, 1)
	require.Equal(t, expectedSubaccountIds, ls.GetNegativeTncSubaccountIds())
	require.Equal(t, expectedSubaccountIds, ls.GetNegativeTncSubaccountIds())
	require.Equal(t, expectedSubaccountIds, ls.GetNegativeTncSubaccountIds())
}

func TestSubaccountsWithOpenPositions_Multiple_Reads(t *testing.T) {
	ls := liquidationstypes.NewDaemonLiquidationInfo()
	require.Empty(t, ls.GetNegativeTncSubaccountIds())

	info := clobtypes.SubaccountOpenPositionInfo{
		PerpetualId: 0,
		SubaccountsWithLongPosition: []satypes.SubaccountId{
			constants.Alice_Num1,
		},
		SubaccountsWithShortPosition: []satypes.SubaccountId{
			constants.Bob_Num0,
		},
	}

	input := []clobtypes.SubaccountOpenPositionInfo{info}
	ls.UpdateSubaccountsWithPositions(input, 1)

	expected := []satypes.SubaccountId{
		constants.Alice_Num1,
		constants.Bob_Num0,
	}
	require.Equal(t, expected, ls.GetSubaccountsWithOpenPositions(0))
	require.Equal(t, expected, ls.GetSubaccountsWithOpenPositions(0))
	require.Equal(t, expected, ls.GetSubaccountsWithOpenPositions(0))
}

func TestLiquidatableSubaccountIds_Multiple_Writes(t *testing.T) {
	ls := liquidationstypes.NewDaemonLiquidationInfo()
	require.Empty(t, ls.GetLiquidatableSubaccountIds())

	expectedSubaccountIds := []satypes.SubaccountId{
		constants.Alice_Num1,
	}
	ls.Update(1, expectedSubaccountIds, []satypes.SubaccountId{}, []clobtypes.SubaccountOpenPositionInfo{})
	require.Equal(t, expectedSubaccountIds, ls.GetLiquidatableSubaccountIds())

	expectedSubaccountIds = []satypes.SubaccountId{
		constants.Alice_Num1,
		constants.Bob_Num0,
	}
	ls.Update(
		1,
		[]satypes.SubaccountId{constants.Bob_Num0},
		[]satypes.SubaccountId{},
		[]clobtypes.SubaccountOpenPositionInfo{},
	)
	require.Equal(t, expectedSubaccountIds, ls.GetLiquidatableSubaccountIds())

	expectedSubaccountIds = []satypes.SubaccountId{
		constants.Carl_Num0,
	}
	ls.Update(2, expectedSubaccountIds, []satypes.SubaccountId{}, []clobtypes.SubaccountOpenPositionInfo{})
	require.Equal(t, expectedSubaccountIds, ls.GetLiquidatableSubaccountIds())
}

func TestNegativeTncSubaccounts_Multiple_Writes(t *testing.T) {
	ls := liquidationstypes.NewDaemonLiquidationInfo()
	require.Empty(t, ls.GetNegativeTncSubaccountIds())

	expectedSubaccountIds := []satypes.SubaccountId{
		constants.Alice_Num1,
	}
	ls.Update(1, []satypes.SubaccountId{}, expectedSubaccountIds, []clobtypes.SubaccountOpenPositionInfo{})
	require.Equal(t, expectedSubaccountIds, ls.GetNegativeTncSubaccountIds())

	expectedSubaccountIds = []satypes.SubaccountId{
		constants.Alice_Num1,
		constants.Bob_Num0,
	}
	ls.Update(
		1,
		[]satypes.SubaccountId{},
		[]satypes.SubaccountId{constants.Bob_Num0},
		[]clobtypes.SubaccountOpenPositionInfo{},
	)
	require.Equal(t, expectedSubaccountIds, ls.GetNegativeTncSubaccountIds())

	expectedSubaccountIds = []satypes.SubaccountId{
		constants.Carl_Num0,
	}
	ls.Update(2, []satypes.SubaccountId{}, expectedSubaccountIds, []clobtypes.SubaccountOpenPositionInfo{})
	require.Equal(t, expectedSubaccountIds, ls.GetNegativeTncSubaccountIds())
}

func TestSubaccountsWithOpenPositions_Multiple_Writes(t *testing.T) {
	ls := liquidationstypes.NewDaemonLiquidationInfo()
	require.Empty(t, ls.GetSubaccountsWithOpenPositions(0))

	info := clobtypes.SubaccountOpenPositionInfo{
		PerpetualId: 0,
		SubaccountsWithLongPosition: []satypes.SubaccountId{
			constants.Alice_Num1,
		},
		SubaccountsWithShortPosition: []satypes.SubaccountId{
			constants.Bob_Num0,
		},
	}

	input := []clobtypes.SubaccountOpenPositionInfo{info}
	ls.Update(1, []satypes.SubaccountId{}, []satypes.SubaccountId{}, input)
	expected := []satypes.SubaccountId{
		constants.Alice_Num1,
		constants.Bob_Num0,
	}
	require.Equal(t, expected, ls.GetSubaccountsWithOpenPositions(0))

	info2 := clobtypes.SubaccountOpenPositionInfo{
		PerpetualId: 0,
		SubaccountsWithLongPosition: []satypes.SubaccountId{
			constants.Carl_Num0,
		},
		SubaccountsWithShortPosition: []satypes.SubaccountId{
			constants.Dave_Num0,
		},
	}

	input2 := []clobtypes.SubaccountOpenPositionInfo{info2}
	ls.Update(1, []satypes.SubaccountId{}, []satypes.SubaccountId{}, input2)
	expected = []satypes.SubaccountId{
		constants.Alice_Num1,
		constants.Carl_Num0,
		constants.Bob_Num0,
		constants.Dave_Num0,
	}
	require.Equal(t, expected, ls.GetSubaccountsWithOpenPositions(0))

	info3 := clobtypes.SubaccountOpenPositionInfo{
		PerpetualId: 0,
		SubaccountsWithLongPosition: []satypes.SubaccountId{
			constants.Dave_Num1,
		},
		SubaccountsWithShortPosition: []satypes.SubaccountId{
			constants.Alice_Num1,
		},
	}

	input3 := []clobtypes.SubaccountOpenPositionInfo{info3}
	ls.Update(2, []satypes.SubaccountId{}, []satypes.SubaccountId{}, input3)
	expected = []satypes.SubaccountId{
		constants.Dave_Num1,
		constants.Alice_Num1,
	}
	require.Equal(t, expected, ls.GetSubaccountsWithOpenPositions(0))
}

func TestLiquidatableSubaccountIds_Empty_Update(t *testing.T) {
	ls := liquidationstypes.NewDaemonLiquidationInfo()
	require.Empty(t, ls.GetLiquidatableSubaccountIds())

	expectedSubaccountIds := []satypes.SubaccountId{
		constants.Alice_Num1,
	}
	ls.Update(1, expectedSubaccountIds, []satypes.SubaccountId{}, []clobtypes.SubaccountOpenPositionInfo{})
	require.Equal(t, expectedSubaccountIds, ls.GetLiquidatableSubaccountIds())

	expectedSubaccountIds = []satypes.SubaccountId{}
	ls.Update(2, expectedSubaccountIds, []satypes.SubaccountId{}, []clobtypes.SubaccountOpenPositionInfo{})
	require.Empty(t, ls.GetLiquidatableSubaccountIds())
}

func TestNegativeTnc_Empty_Update(t *testing.T) {
	ls := liquidationstypes.NewDaemonLiquidationInfo()
	require.Empty(t, ls.GetNegativeTncSubaccountIds())

	expectedSubaccountIds := []satypes.SubaccountId{
		constants.Alice_Num1,
	}
	ls.Update(1, []satypes.SubaccountId{}, expectedSubaccountIds, []clobtypes.SubaccountOpenPositionInfo{})
	require.Equal(t, expectedSubaccountIds, ls.GetNegativeTncSubaccountIds())

	expectedSubaccountIds = []satypes.SubaccountId{}
	ls.Update(2, []satypes.SubaccountId{}, expectedSubaccountIds, []clobtypes.SubaccountOpenPositionInfo{})
	require.Empty(t, ls.GetNegativeTncSubaccountIds())
}

func TestSubaccountsWithOpenPosition_Empty_Update(t *testing.T) {
	ls := liquidationstypes.NewDaemonLiquidationInfo()
	require.Empty(t, ls.GetSubaccountsWithOpenPositions(0))

	info := clobtypes.SubaccountOpenPositionInfo{
		PerpetualId: 0,
		SubaccountsWithLongPosition: []satypes.SubaccountId{
			constants.Alice_Num1,
		},
		SubaccountsWithShortPosition: []satypes.SubaccountId{
			constants.Bob_Num0,
		},
	}
	input := []clobtypes.SubaccountOpenPositionInfo{info}
	ls.Update(1, []satypes.SubaccountId{}, []satypes.SubaccountId{}, input)
	expected := []satypes.SubaccountId{
		constants.Alice_Num1,
		constants.Bob_Num0,
	}
	require.Equal(t, expected, ls.GetSubaccountsWithOpenPositions(0))

	input2 := []clobtypes.SubaccountOpenPositionInfo{}
	ls.Update(2, []satypes.SubaccountId{}, []satypes.SubaccountId{}, input2)
	require.Empty(t, ls.GetSubaccountsWithOpenPositions(0))
}
