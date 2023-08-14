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
	liquidatableSubaccountIds := liquidationtypes.NewLiquidatableSubaccountIds()

	s := createServerWithMocks(
		mockGrpcServer,
		mockFileHandler,
	).WithLiquidatableSubaccountIds(
		liquidatableSubaccountIds,
	)
	_, err := s.LiquidateSubaccounts(grpc.Ctx, &api.LiquidateSubaccountsRequest{
		SubaccountIds: []satypes.SubaccountId{},
	})
	require.NoError(t, err)
	require.Empty(t, liquidatableSubaccountIds.GetSubaccountIds())
}

func TestLiquidateSubaccounts_Multiple_Subaccount_Ids(t *testing.T) {
	mockGrpcServer := &mocks.GrpcServer{}
	mockFileHandler := &mocks.FileHandler{}
	liquidatableSubaccountIds := liquidationtypes.NewLiquidatableSubaccountIds()

	s := createServerWithMocks(
		mockGrpcServer,
		mockFileHandler,
	).WithLiquidatableSubaccountIds(
		liquidatableSubaccountIds,
	)

	expectedSubaccountIds := []satypes.SubaccountId{
		constants.Alice_Num1,
		constants.Bob_Num0,
		constants.Carl_Num0,
	}
	_, err := s.LiquidateSubaccounts(grpc.Ctx, &api.LiquidateSubaccountsRequest{
		SubaccountIds: expectedSubaccountIds,
	})
	require.NoError(t, err)

	actualSubaccountIds := liquidatableSubaccountIds.GetSubaccountIds()
	require.Equal(t, expectedSubaccountIds, actualSubaccountIds)
}
