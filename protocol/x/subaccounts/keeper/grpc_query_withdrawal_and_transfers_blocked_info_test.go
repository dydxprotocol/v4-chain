package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	btkeeper "github.com/dydxprotocol/v4-chain/protocol/x/blocktime/keeper"
	sakeeper "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

func TestQueryWithdrawalAndTransfersBlockedInfo(t *testing.T) {
	for testName, tc := range map[string]struct {
		// Setup.
		setup func(ctx sdktypes.Context, sk sakeeper.Keeper, bk btkeeper.Keeper)

		// Parameters.
		request *types.QueryGetWithdrawalAndTransfersBlockedInfoRequest

		// Expectations.
		response *types.QueryGetWithdrawalAndTransfersBlockedInfoResponse
		err      error
	}{
		"Nil request returns an error": {
			setup: func(ctx sdktypes.Context, sk sakeeper.Keeper, bk btkeeper.Keeper) {},

			err: status.Error(codes.InvalidArgument, "invalid request"),
		},
		`No negative TNC subaccount or chain outage in state returns withdrawals and transfers unblocked
            at block 0`: {
			setup: func(ctx sdktypes.Context, sk sakeeper.Keeper, bk btkeeper.Keeper) {},

			request: &types.QueryGetWithdrawalAndTransfersBlockedInfoRequest{},

			response: &types.QueryGetWithdrawalAndTransfersBlockedInfoResponse{
				NegativeTncSubaccountSeenAtBlock:        0,
				ChainOutageSeenAtBlock:                  0,
				WithdrawalsAndTransfersUnblockedAtBlock: 0,
			},
		},
		`Negative TNC subaccount seen in state returns withdrawals and transfers unblocked
            after the delay`: {
			setup: func(ctx sdktypes.Context, sk sakeeper.Keeper, bk btkeeper.Keeper) {
				sk.SetNegativeTncSubaccountSeenAtBlock(ctx, 7)
			},

			request: &types.QueryGetWithdrawalAndTransfersBlockedInfoRequest{},

			response: &types.QueryGetWithdrawalAndTransfersBlockedInfoResponse{
				NegativeTncSubaccountSeenAtBlock:        7,
				ChainOutageSeenAtBlock:                  0,
				WithdrawalsAndTransfersUnblockedAtBlock: 7 + types.WITHDRAWAL_AND_TRANSFERS_BLOCKED_AFTER_NEGATIVE_TNC_SUBACCOUNT_SEEN_BLOCKS,
			},
		},
		`Chain outage seen in state returns withdrawals and transfers unblocked after the delay`: {
			setup: func(ctx sdktypes.Context, k sakeeper.Keeper) {
				k
				k.bloc(ctx, 7)
			},

			request: &types.QueryGetWithdrawalAndTransfersBlockedInfoRequest{},

			response: &types.QueryGetWithdrawalAndTransfersBlockedInfoResponse{
				NegativeTncSubaccountSeenAtBlock:        7,
				ChainOutageSeenAtBlock:                  0,
				WithdrawalsAndTransfersUnblockedAtBlock: 7 + types.WITHDRAWAL_AND_TRANSFERS_BLOCKED_AFTER_NEGATIVE_TNC_SUBACCOUNT_SEEN_BLOCKS,
			},
		},
	} {
		t.Run(testName, func(t *testing.T) {
			ctx, keeper, _, _, _, _, _, _, _ := keepertest.SubaccountsKeepers(t, true)
			tc.setup(ctx, *keeper)
			response, err := keeper.GetWithdrawalAndTransfersBlockedInfo(ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.response, response)
			}
		})
	}
}
