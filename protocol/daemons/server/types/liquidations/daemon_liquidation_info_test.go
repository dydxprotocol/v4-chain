package types_test

import (
	"testing"

	liquidationstypes "github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/server/types/liquidations"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/constants"
	clobtypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/clob/types"
	satypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/require"
)

func TestNewDaemonLiquidationInfo(t *testing.T) {
	ls := liquidationstypes.NewDaemonLiquidationInfo()
	require.Empty(t, ls.GetSubaccountsWithOpenPositions(0))
}

func TestSubaccountsWithOpenPositions_Multiple_Reads(t *testing.T) {
	ls := liquidationstypes.NewDaemonLiquidationInfo()

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
	ls.UpdateSubaccountsWithPositions(input)

	expected := []satypes.SubaccountId{
		constants.Alice_Num1,
		constants.Bob_Num0,
	}
	require.Equal(t, expected, ls.GetSubaccountsWithOpenPositions(0))
	require.Equal(t, expected, ls.GetSubaccountsWithOpenPositions(0))
	require.Equal(t, expected, ls.GetSubaccountsWithOpenPositions(0))
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
	ls.UpdateSubaccountsWithPositions(input)
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
	ls.UpdateSubaccountsWithPositions(input2)
	expected = []satypes.SubaccountId{
		constants.Carl_Num0,
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
	ls.UpdateSubaccountsWithPositions(input3)
	expected = []satypes.SubaccountId{
		constants.Dave_Num1,
		constants.Alice_Num1,
	}
	require.Equal(t, expected, ls.GetSubaccountsWithOpenPositions(0))
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
	ls.UpdateSubaccountsWithPositions(input)
	expected := []satypes.SubaccountId{
		constants.Alice_Num1,
		constants.Bob_Num0,
	}
	require.Equal(t, expected, ls.GetSubaccountsWithOpenPositions(0))

	input2 := []clobtypes.SubaccountOpenPositionInfo{}
	ls.UpdateSubaccountsWithPositions(input2)
	require.Empty(t, ls.GetSubaccountsWithOpenPositions(0))
}
