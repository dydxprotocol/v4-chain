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
	require.Empty(t, ls.GetSubaccountsWithPositions())
}

func TestLiquidatableSubaccountIds_Multiple_Reads(t *testing.T) {
	ls := liquidationstypes.NewDaemonLiquidationInfo()
	require.Empty(t, ls.GetLiquidatableSubaccountIds())

	expectedSubaccountIds := []satypes.SubaccountId{
		constants.Alice_Num1,
		constants.Bob_Num0,
	}
	ls.UpdateLiquidatableSubaccountIds(expectedSubaccountIds)
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
	ls.UpdateNegativeTncSubaccountIds(expectedSubaccountIds)
	require.Equal(t, expectedSubaccountIds, ls.GetNegativeTncSubaccountIds())
	require.Equal(t, expectedSubaccountIds, ls.GetNegativeTncSubaccountIds())
	require.Equal(t, expectedSubaccountIds, ls.GetNegativeTncSubaccountIds())
}

func TestSubaccountsWithOpenPositions_Multiple_Reads(t *testing.T) {
	ls := liquidationstypes.NewDaemonLiquidationInfo()
	require.Empty(t, ls.GetNegativeTncSubaccountIds())

	expected := map[uint32]*clobtypes.SubaccountOpenPositionInfo{
		0: {
			PerpetualId: 0,
			SubaccountsWithLongPosition: []satypes.SubaccountId{
				constants.Alice_Num1,
			},
			SubaccountsWithShortPosition: []satypes.SubaccountId{
				constants.Bob_Num0,
			},
		},
	}
	ls.UpdateSubaccountsWithPositions(expected)
	require.Equal(t, expected, ls.GetSubaccountsWithPositions())
	require.Equal(t, expected, ls.GetSubaccountsWithPositions())
	require.Equal(t, expected, ls.GetSubaccountsWithPositions())
}

func TestLiquidatableSubaccountIds_Multiple_Writes(t *testing.T) {
	ls := liquidationstypes.NewDaemonLiquidationInfo()
	require.Empty(t, ls.GetLiquidatableSubaccountIds())

	expectedSubaccountIds := []satypes.SubaccountId{
		constants.Alice_Num1,
	}
	ls.UpdateLiquidatableSubaccountIds(expectedSubaccountIds)
	require.Equal(t, expectedSubaccountIds, ls.GetLiquidatableSubaccountIds())

	expectedSubaccountIds = []satypes.SubaccountId{
		constants.Bob_Num0,
	}
	ls.UpdateLiquidatableSubaccountIds(expectedSubaccountIds)
	require.Equal(t, expectedSubaccountIds, ls.GetLiquidatableSubaccountIds())

	expectedSubaccountIds = []satypes.SubaccountId{
		constants.Carl_Num0,
	}
	ls.UpdateLiquidatableSubaccountIds(expectedSubaccountIds)
	require.Equal(t, expectedSubaccountIds, ls.GetLiquidatableSubaccountIds())
}

func TestNegativeTncSubaccounts_Multiple_Writes(t *testing.T) {
	ls := liquidationstypes.NewDaemonLiquidationInfo()
	require.Empty(t, ls.GetNegativeTncSubaccountIds())

	expectedSubaccountIds := []satypes.SubaccountId{
		constants.Alice_Num1,
	}
	ls.UpdateNegativeTncSubaccountIds(expectedSubaccountIds)
	require.Equal(t, expectedSubaccountIds, ls.GetNegativeTncSubaccountIds())

	expectedSubaccountIds = []satypes.SubaccountId{
		constants.Bob_Num0,
	}
	ls.UpdateNegativeTncSubaccountIds(expectedSubaccountIds)
	require.Equal(t, expectedSubaccountIds, ls.GetNegativeTncSubaccountIds())

	expectedSubaccountIds = []satypes.SubaccountId{
		constants.Carl_Num0,
	}
	ls.UpdateNegativeTncSubaccountIds(expectedSubaccountIds)
	require.Equal(t, expectedSubaccountIds, ls.GetNegativeTncSubaccountIds())
}

func TestSubaccountsWithOpenPositions_Multiple_Writes(t *testing.T) {
	ls := liquidationstypes.NewDaemonLiquidationInfo()
	require.Empty(t, ls.GetSubaccountsWithPositions())

	expected := map[uint32]*clobtypes.SubaccountOpenPositionInfo{
		0: {
			PerpetualId: 0,
			SubaccountsWithLongPosition: []satypes.SubaccountId{
				constants.Alice_Num1,
			},
			SubaccountsWithShortPosition: []satypes.SubaccountId{
				constants.Bob_Num0,
			},
		},
	}
	ls.UpdateSubaccountsWithPositions(expected)
	require.Equal(t, expected, ls.GetSubaccountsWithPositions())

	expected = map[uint32]*clobtypes.SubaccountOpenPositionInfo{
		0: {
			PerpetualId: 0,
			SubaccountsWithLongPosition: []satypes.SubaccountId{
				constants.Carl_Num0,
			},
			SubaccountsWithShortPosition: []satypes.SubaccountId{
				constants.Dave_Num0,
			},
		},
	}
	ls.UpdateSubaccountsWithPositions(expected)
	require.Equal(t, expected, ls.GetSubaccountsWithPositions())

	expected = map[uint32]*clobtypes.SubaccountOpenPositionInfo{
		0: {
			PerpetualId: 0,
			SubaccountsWithLongPosition: []satypes.SubaccountId{
				constants.Dave_Num1,
			},
			SubaccountsWithShortPosition: []satypes.SubaccountId{
				constants.Alice_Num1,
			},
		},
	}
	ls.UpdateSubaccountsWithPositions(expected)
	require.Equal(t, expected, ls.GetSubaccountsWithPositions())
}

func TestLiquidatableSubaccountIds_Empty_Update(t *testing.T) {
	ls := liquidationstypes.NewDaemonLiquidationInfo()
	require.Empty(t, ls.GetLiquidatableSubaccountIds())

	expectedSubaccountIds := []satypes.SubaccountId{
		constants.Alice_Num1,
	}
	ls.UpdateLiquidatableSubaccountIds(expectedSubaccountIds)
	require.Equal(t, expectedSubaccountIds, ls.GetLiquidatableSubaccountIds())

	expectedSubaccountIds = []satypes.SubaccountId{}
	ls.UpdateLiquidatableSubaccountIds(expectedSubaccountIds)
	require.Empty(t, ls.GetLiquidatableSubaccountIds())
}

func TestNegativeTnc_Empty_Update(t *testing.T) {
	ls := liquidationstypes.NewDaemonLiquidationInfo()
	require.Empty(t, ls.GetNegativeTncSubaccountIds())

	expectedSubaccountIds := []satypes.SubaccountId{
		constants.Alice_Num1,
	}
	ls.UpdateNegativeTncSubaccountIds(expectedSubaccountIds)
	require.Equal(t, expectedSubaccountIds, ls.GetNegativeTncSubaccountIds())

	expectedSubaccountIds = []satypes.SubaccountId{}
	ls.UpdateNegativeTncSubaccountIds(expectedSubaccountIds)
	require.Empty(t, ls.GetNegativeTncSubaccountIds())
}

func TestSubaccountsWithOpenPosition_Empty_Update(t *testing.T) {
	ls := liquidationstypes.NewDaemonLiquidationInfo()
	require.Empty(t, ls.GetSubaccountsWithPositions())

	expected := map[uint32]*clobtypes.SubaccountOpenPositionInfo{
		0: {
			PerpetualId: 0,
			SubaccountsWithLongPosition: []satypes.SubaccountId{
				constants.Alice_Num1,
			},
			SubaccountsWithShortPosition: []satypes.SubaccountId{
				constants.Bob_Num0,
			},
		},
	}
	ls.UpdateSubaccountsWithPositions(expected)
	require.Equal(t, expected, ls.GetSubaccountsWithPositions())

	expected = map[uint32]*clobtypes.SubaccountOpenPositionInfo{}
	ls.UpdateSubaccountsWithPositions(expected)
	require.Empty(t, ls.GetSubaccountsWithPositions())
}
