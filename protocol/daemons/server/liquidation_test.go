package server_test

import (
	"testing"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/liquidation/api"
	liquidationtypes "github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/server/types/liquidations"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/mocks"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/constants"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/grpc"
	clobtypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/clob/types"
	satypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/require"
)

func TestLiquidateSubaccounts_Empty_Update_Subaccount_Open_Positions(t *testing.T) {
	mockGrpcServer := &mocks.GrpcServer{}
	mockFileHandler := &mocks.FileHandler{}
	daemonLiquidationInfo := liquidationtypes.NewDaemonLiquidationInfo()

	s := createServerWithMocks(
		t,
		mockGrpcServer,
		mockFileHandler,
	).WithDaemonLiquidationInfo(
		daemonLiquidationInfo,
	)
	_, err := s.LiquidateSubaccounts(grpc.Ctx, &api.LiquidateSubaccountsRequest{
		SubaccountOpenPositionInfo: []clobtypes.SubaccountOpenPositionInfo{},
	})
	require.NoError(t, err)
	require.Empty(t, daemonLiquidationInfo.GetSubaccountsWithOpenPositions(0))
}

func TestLiquidateSubaccounts_Multiple_Subaccount_Open_Positions(t *testing.T) {
	mockGrpcServer := &mocks.GrpcServer{}
	mockFileHandler := &mocks.FileHandler{}
	daemonLiquidationInfo := liquidationtypes.NewDaemonLiquidationInfo()

	s := createServerWithMocks(
		t,
		mockGrpcServer,
		mockFileHandler,
	).WithDaemonLiquidationInfo(
		daemonLiquidationInfo,
	)

	subaccountOpenPositions := []clobtypes.SubaccountOpenPositionInfo{
		{
			PerpetualId: 0,
			SubaccountsWithLongPosition: []satypes.SubaccountId{
				constants.Alice_Num0,
			},
			SubaccountsWithShortPosition: []satypes.SubaccountId{}, // No short positions for perp 0
		},
		{
			PerpetualId: 1,
			SubaccountsWithLongPosition: []satypes.SubaccountId{
				constants.Alice_Num1,
				constants.Bob_Num0,
			},
			SubaccountsWithShortPosition: []satypes.SubaccountId{
				constants.Carl_Num0,
			},
		},
	}

	expectedSubaccountIdsPerp1 := []satypes.SubaccountId{
		constants.Alice_Num1,
		constants.Bob_Num0,
		constants.Carl_Num0,
	}

	expectedSubaccountIdsPerp1Long := []satypes.SubaccountId{
		constants.Alice_Num1,
		constants.Bob_Num0,
	}

	_, err := s.LiquidateSubaccounts(grpc.Ctx, &api.LiquidateSubaccountsRequest{
		SubaccountOpenPositionInfo: subaccountOpenPositions,
	})
	require.NoError(t, err)

	actualSubaccountIdsPerp1 := daemonLiquidationInfo.GetSubaccountsWithOpenPositions(0)
	require.Equal(t, expectedSubaccountIdsPerp1, actualSubaccountIdsPerp1)

	actualSubaccountIdsPerp1Long := daemonLiquidationInfo.GetSubaccountsWithOpenPositionsOnSide(0, true)
	require.Equal(t, expectedSubaccountIdsPerp1Long, actualSubaccountIdsPerp1Long)
}
