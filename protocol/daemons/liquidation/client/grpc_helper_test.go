package client_test

import (
	"context"
	"errors"
	"testing"

	"cosmossdk.io/log"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/flags"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/liquidation/api"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/liquidation/client"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/mocks"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/constants"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/grpc"
	blocktimetypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/blocktime/types"
	clobtypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/clob/types"
	satypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/subaccounts/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetPreviousBlockInfo(t *testing.T) {
	tests := map[string]struct {
		// mocks
		setupMocks func(
			ctx context.Context,
			mck *mocks.QueryClient,
		)

		// expectations
		expectedBlockHeight uint32
		expectedError       error
	}{
		"Success": {
			setupMocks: func(
				ctx context.Context,
				mck *mocks.QueryClient,
			) {
				response := &blocktimetypes.QueryPreviousBlockInfoResponse{
					Info: &blocktimetypes.BlockInfo{
						Height:    uint32(50),
						Timestamp: constants.TimeTen,
					},
				}
				mck.On("PreviousBlockInfo", ctx, mock.Anything).Return(response, nil)
			},
			expectedBlockHeight: 50,
		},
		"Errors are propagated": {
			setupMocks: func(
				ctx context.Context,
				mck *mocks.QueryClient,
			) {
				mck.On("PreviousBlockInfo", ctx, mock.Anything).Return(nil, errors.New("test error"))
			},
			expectedBlockHeight: 0,
			expectedError:       errors.New("test error"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			queryClientMock := &mocks.QueryClient{}
			tc.setupMocks(grpc.Ctx, queryClientMock)

			daemon := client.NewClient(log.NewNopLogger())
			daemon.BlocktimeQueryClient = queryClientMock
			actualBlockHeight, err := daemon.GetPreviousBlockInfo(grpc.Ctx)

			if err != nil {
				require.EqualError(t, err, tc.expectedError.Error())
			} else {
				require.Equal(t, tc.expectedBlockHeight, actualBlockHeight)
			}
		})
	}
}

func TestGetAllSubaccounts(t *testing.T) {
	df := flags.GetDefaultDaemonFlags()
	tests := map[string]struct {
		// mocks
		setupMocks func(ctx context.Context, mck *mocks.QueryClient)

		// expectations
		expectedSubaccounts []satypes.Subaccount
		expectedError       error
	}{
		"Success": {
			setupMocks: func(ctx context.Context, mck *mocks.QueryClient) {
				req := &satypes.QueryAllSubaccountRequest{
					Pagination: &query.PageRequest{
						Limit: df.Liquidation.QueryPageLimit,
					},
				}
				response := &satypes.QuerySubaccountAllResponse{
					Subaccount: []satypes.Subaccount{
						constants.Carl_Num0_599USD,
						constants.Dave_Num0_599USD,
					},
				}
				mck.On("SubaccountAll", ctx, req).Return(response, nil)
			},
			expectedSubaccounts: []satypes.Subaccount{
				constants.Carl_Num0_599USD,
				constants.Dave_Num0_599USD,
			},
		},
		"Success Paginated": {
			setupMocks: func(ctx context.Context, mck *mocks.QueryClient) {
				req := &satypes.QueryAllSubaccountRequest{
					Pagination: &query.PageRequest{
						Limit: df.Liquidation.QueryPageLimit,
					},
				}
				nextKey := []byte("next key")
				response := &satypes.QuerySubaccountAllResponse{
					Subaccount: []satypes.Subaccount{
						constants.Carl_Num0_599USD,
					},
					Pagination: &query.PageResponse{
						NextKey: nextKey,
					},
				}
				mck.On("SubaccountAll", ctx, req).Return(response, nil)
				req2 := &satypes.QueryAllSubaccountRequest{
					Pagination: &query.PageRequest{
						Key:   nextKey,
						Limit: df.Liquidation.QueryPageLimit,
					},
				}
				response2 := &satypes.QuerySubaccountAllResponse{
					Subaccount: []satypes.Subaccount{
						constants.Dave_Num0_599USD,
					},
				}
				mck.On("SubaccountAll", ctx, req2).Return(response2, nil)
			},
			expectedSubaccounts: []satypes.Subaccount{
				constants.Carl_Num0_599USD,
				constants.Dave_Num0_599USD,
			},
		},
		"Errors are propagated": {
			setupMocks: func(ctx context.Context, mck *mocks.QueryClient) {
				req := &satypes.QueryAllSubaccountRequest{
					Pagination: &query.PageRequest{
						Limit: df.Liquidation.QueryPageLimit,
					},
				}
				mck.On("SubaccountAll", ctx, req).Return(nil, errors.New("test error"))
			},
			expectedError: errors.New("test error"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			queryClientMock := &mocks.QueryClient{}
			tc.setupMocks(grpc.Ctx, queryClientMock)

			daemon := client.NewClient(log.NewNopLogger())
			daemon.SubaccountQueryClient = queryClientMock
			actual, err := daemon.GetAllSubaccounts(
				grpc.Ctx,
				df.Liquidation.QueryPageLimit,
			)
			if err != nil {
				require.EqualError(t, err, tc.expectedError.Error())
			} else {
				require.Equal(t, tc.expectedSubaccounts, actual)
			}
		})
	}
}

func TestSendLiquidatableSubaccountIds(t *testing.T) {
	tests := map[string]struct {
		// mocks
		setupMocks                 func(context.Context, *mocks.QueryClient)
		subaccountOpenPositionInfo map[uint32]*clobtypes.SubaccountOpenPositionInfo

		// expectations
		expectedError error
	}{
		"Success": {
			setupMocks: func(ctx context.Context, mck *mocks.QueryClient) {
				req := &api.LiquidateSubaccountsRequest{
					SubaccountOpenPositionInfo: []clobtypes.SubaccountOpenPositionInfo{
						{
							PerpetualId: 0,
							SubaccountsWithLongPosition: []satypes.SubaccountId{
								constants.Alice_Num0,
								constants.Carl_Num0,
							},
							SubaccountsWithShortPosition: []satypes.SubaccountId{
								constants.Bob_Num0,
								constants.Dave_Num0,
							},
						},
					},
				}
				response := &api.LiquidateSubaccountsResponse{}
				mck.On("LiquidateSubaccounts", ctx, req).Return(response, nil)
			},
			subaccountOpenPositionInfo: map[uint32]*clobtypes.SubaccountOpenPositionInfo{
				0: {
					PerpetualId: 0,
					SubaccountsWithLongPosition: []satypes.SubaccountId{
						constants.Alice_Num0,
						constants.Carl_Num0,
					},
					SubaccountsWithShortPosition: []satypes.SubaccountId{
						constants.Bob_Num0,
						constants.Dave_Num0,
					},
				},
			},
		},
		"Success Empty": {
			setupMocks: func(ctx context.Context, mck *mocks.QueryClient) {
				req := &api.LiquidateSubaccountsRequest{
					SubaccountOpenPositionInfo: []clobtypes.SubaccountOpenPositionInfo{},
				}
				response := &api.LiquidateSubaccountsResponse{}
				mck.On("LiquidateSubaccounts", ctx, req).Return(response, nil)
			},
			subaccountOpenPositionInfo: map[uint32]*clobtypes.SubaccountOpenPositionInfo{},
		},
		"Errors are propagated": {
			setupMocks: func(ctx context.Context, mck *mocks.QueryClient) {
				req := &api.LiquidateSubaccountsRequest{
					SubaccountOpenPositionInfo: []clobtypes.SubaccountOpenPositionInfo{},
				}
				mck.On("LiquidateSubaccounts", ctx, req).Return(nil, errors.New("test error"))
			},
			subaccountOpenPositionInfo: map[uint32]*clobtypes.SubaccountOpenPositionInfo{},
			expectedError:              errors.New("test error"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			queryClientMock := &mocks.QueryClient{}
			tc.setupMocks(grpc.Ctx, queryClientMock)

			daemon := client.NewClient(log.NewNopLogger())
			daemon.LiquidationServiceClient = queryClientMock

			err := daemon.SendLiquidatableSubaccountIds(
				grpc.Ctx,
				tc.subaccountOpenPositionInfo,
			)
			require.Equal(t, tc.expectedError, err)
		})
	}
}
