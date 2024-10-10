package server_test

import (
	"testing"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/deleveraging/api"
	deleveragingtypes "github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/server/types/deleveraging"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/mocks"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/constants"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/grpc"
	clobtypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/clob/types"
	satypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/require"
)

func TestUpdateSubaccountsListForDeleveragingDaemonRequest_Empty_Update_Subaccount_Open_Positions(t *testing.T) {
	mockGrpcServer := &mocks.GrpcServer{}
	mockFileHandler := &mocks.FileHandler{}
	daemonDeleveragingInfo := deleveragingtypes.NewDaemonDeleveragingInfo()

	s := createServerWithMocks(
		t,
		mockGrpcServer,
		mockFileHandler,
	).WithDaemonDeleveragingInfo(
		daemonDeleveragingInfo,
	)
	_, err := s.UpdateSubaccountsListForDeleveragingDaemon(grpc.Ctx, &api.UpdateSubaccountsListForDeleveragingDaemonRequest{
		SubaccountOpenPositionInfo: []clobtypes.SubaccountOpenPositionInfo{},
	})
	require.NoError(t, err)
	require.Empty(t, daemonDeleveragingInfo.GetSubaccountsWithOpenPositions(0))
}

func TestUpdateSubaccountsListForDeleveragingDaemon_Multiple_Subaccount_Open_Positions(t *testing.T) {
	mockGrpcServer := &mocks.GrpcServer{}
	mockFileHandler := &mocks.FileHandler{}
	daemonDeleveragingInfo := deleveragingtypes.NewDaemonDeleveragingInfo()

	s := createServerWithMocks(
		t,
		mockGrpcServer,
		mockFileHandler,
	).WithDaemonDeleveragingInfo(
		daemonDeleveragingInfo,
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

	_, err := s.UpdateSubaccountsListForDeleveragingDaemon(grpc.Ctx, &api.UpdateSubaccountsListForDeleveragingDaemonRequest{
		SubaccountOpenPositionInfo: subaccountOpenPositions,
	})
	require.NoError(t, err)

	actualSubaccountIdsPerp1 := daemonDeleveragingInfo.GetSubaccountsWithOpenPositions(1)
	require.Equal(t, expectedSubaccountIdsPerp1, actualSubaccountIdsPerp1)

	actualSubaccountIdsPerp1Long := daemonDeleveragingInfo.GetSubaccountsWithOpenPositionsOnSide(1, true)
	require.Equal(t, expectedSubaccountIdsPerp1Long, actualSubaccountIdsPerp1Long)
}
