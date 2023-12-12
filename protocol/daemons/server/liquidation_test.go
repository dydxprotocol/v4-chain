package server_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/daemons/liquidation/api"
	liquidationtypes "github.com/dydxprotocol/v4-chain/protocol/daemons/server/types/liquidations"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/grpc"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/require"
)

func TestLiquidateSubaccounts_Empty_Update(t *testing.T) {
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
		LiquidatableSubaccountIds: []satypes.SubaccountId{},
	})
	require.NoError(t, err)
	require.Empty(t, daemonLiquidationInfo.GetLiquidatableSubaccountIds())
}

func TestLiquidateSubaccounts_Multiple_Subaccount_Ids(t *testing.T) {
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

	expectedSubaccountIds := []satypes.SubaccountId{
		constants.Alice_Num1,
		constants.Bob_Num0,
		constants.Carl_Num0,
	}
	_, err := s.LiquidateSubaccounts(grpc.Ctx, &api.LiquidateSubaccountsRequest{
		LiquidatableSubaccountIds: expectedSubaccountIds,
	})
	require.NoError(t, err)

	actualSubaccountIds := daemonLiquidationInfo.GetLiquidatableSubaccountIds()
	require.Equal(t, expectedSubaccountIds, actualSubaccountIds)
}
