package client_test

import (
	"context"
	"testing"

	"cosmossdk.io/log"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/deleveraging/api"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/deleveraging/client"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/flags"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/mocks"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/constants"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/grpc"
	blocktimetypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/blocktime/types"
	clobtypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/clob/types"
	satypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestRunDeleveragingDaemonTaskLoop(t *testing.T) {
	tests := map[string]struct {
		// mocks
		setupMocks func(ctx context.Context, mck *mocks.QueryClient)

		// expectations
		expectedError error
	}{
		"Can get subaccount with short position": {
			setupMocks: func(ctx context.Context, mck *mocks.QueryClient) {
				// Block height.
				res := &blocktimetypes.QueryPreviousBlockInfoResponse{
					Info: &blocktimetypes.BlockInfo{
						Height:    uint32(50),
						Timestamp: constants.TimeTen,
					},
				}
				mck.On("PreviousBlockInfo", mock.Anything, mock.Anything).Return(res, nil)

				// Subaccount.
				res2 := &satypes.QuerySubaccountAllResponse{
					Subaccount: []satypes.Subaccount{
						constants.Carl_Num0_1BTC_Short_54999USD,
					},
				}
				mck.On("SubaccountAll", mock.Anything, mock.Anything).Return(res2, nil)

				// Sends liquidatable subaccount ids to the server.
				req := &api.UpdateSubaccountsListForDeleveragingDaemonRequest{
					SubaccountOpenPositionInfo: []clobtypes.SubaccountOpenPositionInfo{
						{
							PerpetualId:                 0,
							SubaccountsWithLongPosition: []satypes.SubaccountId{},
							SubaccountsWithShortPosition: []satypes.SubaccountId{
								constants.Carl_Num0,
							},
						},
					},
				}
				response3 := &api.UpdateSubaccountsListForDeleveragingDaemonResponse{}
				mck.On("UpdateSubaccountsListForDeleveragingDaemon", ctx, req).Return(response3, nil)
			},
		},
		"Can get subaccount with long position": {
			setupMocks: func(ctx context.Context, mck *mocks.QueryClient) {
				// Block height.
				res := &blocktimetypes.QueryPreviousBlockInfoResponse{
					Info: &blocktimetypes.BlockInfo{
						Height:    uint32(50),
						Timestamp: constants.TimeTen,
					},
				}
				mck.On("PreviousBlockInfo", mock.Anything, mock.Anything).Return(res, nil)

				// Subaccount.
				res2 := &satypes.QuerySubaccountAllResponse{
					Subaccount: []satypes.Subaccount{
						constants.Dave_Num0_1BTC_Long_45001USD_Short,
					},
				}
				mck.On("SubaccountAll", mock.Anything, mock.Anything).Return(res2, nil)

				// Sends liquidatable subaccount ids to the server.
				req := &api.UpdateSubaccountsListForDeleveragingDaemonRequest{
					SubaccountOpenPositionInfo: []clobtypes.SubaccountOpenPositionInfo{
						{
							PerpetualId: 0,
							SubaccountsWithLongPosition: []satypes.SubaccountId{
								constants.Dave_Num0,
							},
							SubaccountsWithShortPosition: []satypes.SubaccountId{},
						},
					},
				}
				response3 := &api.UpdateSubaccountsListForDeleveragingDaemonResponse{}
				mck.On("UpdateSubaccountsListForDeleveragingDaemon", ctx, req).Return(response3, nil)
			},
		},
		"Skip subaccounts with no open positions": {
			setupMocks: func(ctx context.Context, mck *mocks.QueryClient) {
				// Block height.
				res := &blocktimetypes.QueryPreviousBlockInfoResponse{
					Info: &blocktimetypes.BlockInfo{
						Height:    uint32(50),
						Timestamp: constants.TimeTen,
					},
				}
				mck.On("PreviousBlockInfo", mock.Anything, mock.Anything).Return(res, nil)

				// Subaccount.
				res2 := &satypes.QuerySubaccountAllResponse{
					Subaccount: []satypes.Subaccount{
						constants.Alice_Num0_100_000USD,
					},
				}
				mck.On("SubaccountAll", mock.Anything, mock.Anything).Return(res2, nil)

				// Sends liquidatable subaccount ids to the server.
				req := &api.UpdateSubaccountsListForDeleveragingDaemonRequest{
					SubaccountOpenPositionInfo: []clobtypes.SubaccountOpenPositionInfo{},
				}
				response3 := &api.UpdateSubaccountsListForDeleveragingDaemonResponse{}
				mck.On("UpdateSubaccountsListForDeleveragingDaemon", ctx, req).Return(response3, nil)
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			queryClientMock := &mocks.QueryClient{}
			tc.setupMocks(grpc.Ctx, queryClientMock)
			s := client.SubTaskRunnerImpl{}

			c := client.NewClient(log.NewNopLogger())
			c.SubaccountQueryClient = queryClientMock
			c.DeleveragingServiceClient = queryClientMock
			c.BlocktimeQueryClient = queryClientMock

			err := s.RunDeleveragingDaemonTaskLoop(
				grpc.Ctx,
				c,
				flags.GetDefaultDaemonFlags().Deleveraging,
			)
			if tc.expectedError != nil {
				require.EqualError(t, err, tc.expectedError.Error())
			} else {
				require.NoError(t, err)
				queryClientMock.AssertExpectations(t)
			}
		})
	}
}
