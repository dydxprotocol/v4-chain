package keeper_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	btkeeper "github.com/dydxprotocol/v4-chain/protocol/x/blocktime/keeper"
	blocktimetypes "github.com/dydxprotocol/v4-chain/protocol/x/blocktime/types"
	sakeeper "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

func TestQueryWithdrawalAndTransfersBlockedInfo(t *testing.T) {
	for testName, tc := range map[string]struct {
		// Setup.
		setup func(ctx sdktypes.Context, sk sakeeper.Keeper, bk btkeeper.Keeper) error

		// Parameters.
		request *types.QueryGetWithdrawalAndTransfersBlockedInfoRequest

		// Expectations.
		response *types.QueryGetWithdrawalAndTransfersBlockedInfoResponse
		err      error
	}{
		"Nil request returns an error": {
			setup: func(ctx sdktypes.Context, sk sakeeper.Keeper, bk btkeeper.Keeper) error {
				return nil
			},

			err: status.Error(codes.InvalidArgument, "invalid request"),
		},
		`No negative TNC subaccount or chain outage in state returns withdrawals and transfers unblocked
            at block 0`: {
			setup: func(ctx sdktypes.Context, sk sakeeper.Keeper, bk btkeeper.Keeper) error {
				return nil
			},

			request: &types.QueryGetWithdrawalAndTransfersBlockedInfoRequest{
				PerpetualId: uint32(constants.BtcUsd_NoMarginRequirement.Params.Id),
			},

			response: &types.QueryGetWithdrawalAndTransfersBlockedInfoResponse{
				NegativeTncSubaccountSeenAtBlock:        0,
				ChainOutageSeenAtBlock:                  0,
				WithdrawalsAndTransfersUnblockedAtBlock: 0,
			},
		},
		`Negative TNC subaccount seen in state returns withdrawals and transfers unblocked
            after the delay`: {
			setup: func(ctx sdktypes.Context, sk sakeeper.Keeper, bk btkeeper.Keeper) error {
				return sk.SetNegativeTncSubaccountSeenAtBlock(
					ctx,
					constants.BtcUsd_NoMarginRequirement.Params.Id,
					7,
				)
			},

			request: &types.QueryGetWithdrawalAndTransfersBlockedInfoRequest{
				PerpetualId: uint32(constants.BtcUsd_NoMarginRequirement.Params.Id),
			},

			response: &types.QueryGetWithdrawalAndTransfersBlockedInfoResponse{
				NegativeTncSubaccountSeenAtBlock: 7,
				ChainOutageSeenAtBlock:           0,
				WithdrawalsAndTransfersUnblockedAtBlock: 7 +
					types.WITHDRAWAL_AND_TRANSFERS_BLOCKED_AFTER_NEGATIVE_TNC_SUBACCOUNT_SEEN_BLOCKS,
			},
		},
		`Negative TNC subaccount seen in state returns withdrawals and transfers unblocked
            after the delay (for isolated perpetual)`: {
			setup: func(ctx sdktypes.Context, sk sakeeper.Keeper, bk btkeeper.Keeper) error {
				return sk.SetNegativeTncSubaccountSeenAtBlock(
					ctx,
					constants.IsoUsd_IsolatedMarket.Params.Id,
					5,
				)
			},

			request: &types.QueryGetWithdrawalAndTransfersBlockedInfoRequest{
				PerpetualId: uint32(constants.IsoUsd_IsolatedMarket.Params.Id),
			},

			response: &types.QueryGetWithdrawalAndTransfersBlockedInfoResponse{
				NegativeTncSubaccountSeenAtBlock: 5,
				ChainOutageSeenAtBlock:           0,
				WithdrawalsAndTransfersUnblockedAtBlock: 5 +
					types.WITHDRAWAL_AND_TRANSFERS_BLOCKED_AFTER_NEGATIVE_TNC_SUBACCOUNT_SEEN_BLOCKS,
			},
		},
		`Chain outage seen in state returns withdrawals and transfers unblocked after the delay`: {
			setup: func(ctx sdktypes.Context, k sakeeper.Keeper, bk btkeeper.Keeper) error {
				bk.SetAllDowntimeInfo(
					ctx,
					&blocktimetypes.AllDowntimeInfo{
						Infos: []*blocktimetypes.AllDowntimeInfo_DowntimeInfo{
							{
								Duration: 10 * time.Second,
								BlockInfo: blocktimetypes.BlockInfo{
									Height:    30,
									Timestamp: time.Unix(300, 0).UTC(),
								},
							},
							{
								Duration: 5 * time.Minute,
								BlockInfo: blocktimetypes.BlockInfo{
									Height:    25,
									Timestamp: time.Unix(300, 0).UTC(),
								},
							},
						},
					})
				return nil
			},

			request: &types.QueryGetWithdrawalAndTransfersBlockedInfoRequest{
				PerpetualId: constants.BtcUsd_NoMarginRequirement.Params.Id,
			},

			response: &types.QueryGetWithdrawalAndTransfersBlockedInfoResponse{
				NegativeTncSubaccountSeenAtBlock: 0,
				ChainOutageSeenAtBlock:           25,
				WithdrawalsAndTransfersUnblockedAtBlock: 25 +
					types.WITHDRAWAL_AND_TRANSFERS_BLOCKED_AFTER_NEGATIVE_TNC_SUBACCOUNT_SEEN_BLOCKS,
			},
		},
		`Negative TNC subaccount and chain outage seen in state returns withdrawals and transfers
			unblocked after the max block number + delay (negative TNC subaccount block greater)`: {
			setup: func(ctx sdktypes.Context, sk sakeeper.Keeper, bk btkeeper.Keeper) error {
				err := sk.SetNegativeTncSubaccountSeenAtBlock(
					ctx,
					constants.BtcUsd_NoMarginRequirement.Params.Id,
					27,
				)
				if err != nil {
					return nil
				}
				bk.SetAllDowntimeInfo(
					ctx,
					&blocktimetypes.AllDowntimeInfo{
						Infos: []*blocktimetypes.AllDowntimeInfo_DowntimeInfo{
							{
								Duration: 10 * time.Second,
								BlockInfo: blocktimetypes.BlockInfo{
									Height:    30,
									Timestamp: time.Unix(300, 0).UTC(),
								},
							},
							{
								Duration: 5 * time.Minute,
								BlockInfo: blocktimetypes.BlockInfo{
									Height:    25,
									Timestamp: time.Unix(300, 0).UTC(),
								},
							},
						},
					})
				return nil
			},

			request: &types.QueryGetWithdrawalAndTransfersBlockedInfoRequest{
				PerpetualId: constants.BtcUsd_NoMarginRequirement.Params.Id,
			},

			response: &types.QueryGetWithdrawalAndTransfersBlockedInfoResponse{
				NegativeTncSubaccountSeenAtBlock: 27,
				ChainOutageSeenAtBlock:           25,
				WithdrawalsAndTransfersUnblockedAtBlock: 27 +
					types.WITHDRAWAL_AND_TRANSFERS_BLOCKED_AFTER_NEGATIVE_TNC_SUBACCOUNT_SEEN_BLOCKS,
			},
		},
		`Negative TNC subaccount and chain outage seen in state returns withdrawals and transfers
			unblocked after the max block number + delay (chain outage block greater)`: {
			setup: func(ctx sdktypes.Context, sk sakeeper.Keeper, bk btkeeper.Keeper) error {
				err := sk.SetNegativeTncSubaccountSeenAtBlock(
					ctx,
					constants.BtcUsd_NoMarginRequirement.Params.Id,
					37,
				)
				if err != nil {
					return err
				}
				bk.SetAllDowntimeInfo(
					ctx,
					&blocktimetypes.AllDowntimeInfo{
						Infos: []*blocktimetypes.AllDowntimeInfo_DowntimeInfo{
							{
								Duration: 10 * time.Second,
								BlockInfo: blocktimetypes.BlockInfo{
									Height:    50,
									Timestamp: time.Unix(300, 0).UTC(),
								},
							},
							{
								Duration: 5 * time.Minute,
								BlockInfo: blocktimetypes.BlockInfo{
									Height:    47,
									Timestamp: time.Unix(300, 0).UTC(),
								},
							},
						},
					})
				return nil
			},

			request: &types.QueryGetWithdrawalAndTransfersBlockedInfoRequest{
				PerpetualId: constants.BtcUsd_NoMarginRequirement.Params.Id,
			},

			response: &types.QueryGetWithdrawalAndTransfersBlockedInfoResponse{
				NegativeTncSubaccountSeenAtBlock: 37,
				ChainOutageSeenAtBlock:           47,
				WithdrawalsAndTransfersUnblockedAtBlock: 47 +
					types.WITHDRAWAL_AND_TRANSFERS_BLOCKED_AFTER_NEGATIVE_TNC_SUBACCOUNT_SEEN_BLOCKS,
			},
		},
		`Negative TNC subaccount and chain outage seen in state returns withdrawals and transfers
			unblocked after the max block number + delay (both blocks equal)`: {
			setup: func(ctx sdktypes.Context, sk sakeeper.Keeper, bk btkeeper.Keeper) error {
				err := sk.SetNegativeTncSubaccountSeenAtBlock(
					ctx,
					constants.BtcUsd_NoMarginRequirement.Params.Id,
					3,
				)
				if err != nil {
					return err
				}
				bk.SetAllDowntimeInfo(
					ctx,
					&blocktimetypes.AllDowntimeInfo{
						Infos: []*blocktimetypes.AllDowntimeInfo_DowntimeInfo{
							{
								Duration: 10 * time.Second,
								BlockInfo: blocktimetypes.BlockInfo{
									Height:    50,
									Timestamp: time.Unix(300, 0).UTC(),
								},
							},
							{
								Duration: 5 * time.Minute,
								BlockInfo: blocktimetypes.BlockInfo{
									Height:    3,
									Timestamp: time.Unix(300, 0).UTC(),
								},
							},
						},
					})
				return nil
			},

			request: &types.QueryGetWithdrawalAndTransfersBlockedInfoRequest{
				PerpetualId: constants.BtcUsd_NoMarginRequirement.Params.Id,
			},

			response: &types.QueryGetWithdrawalAndTransfersBlockedInfoResponse{
				NegativeTncSubaccountSeenAtBlock: 3,
				ChainOutageSeenAtBlock:           3,
				WithdrawalsAndTransfersUnblockedAtBlock: 3 +
					types.WITHDRAWAL_AND_TRANSFERS_BLOCKED_AFTER_NEGATIVE_TNC_SUBACCOUNT_SEEN_BLOCKS,
			},
		},
	} {
		t.Run(testName, func(t *testing.T) {
			ctx, keeper, pricesKeeper, perpetualsKeeper, _, _, _, blocktimeKeeper, _, _, _ := keepertest.SubaccountsKeepers(
				t,
				true,
			)
			keepertest.CreateTestMarkets(t, ctx, pricesKeeper)
			keepertest.CreateTestLiquidityTiers(t, ctx, perpetualsKeeper)
			keepertest.CreateTestPerpetuals(t, ctx, perpetualsKeeper)
			err := tc.setup(ctx, *keeper, *blocktimeKeeper)
			require.NoError(t, err)
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
