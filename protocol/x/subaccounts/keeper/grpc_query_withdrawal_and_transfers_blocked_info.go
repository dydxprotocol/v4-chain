package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

func (k Keeper) GetWithdrawalAndTransfersBlockedInfo(
	c context.Context,
	req *types.QueryGetWithdrawalAndTransfersBlockedInfoRequest,
) (*types.QueryGetWithdrawalAndTransfersBlockedInfoResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdktypes.UnwrapSDKContext(c)

	downtimeInfo := k.blocktimeKeeper.GetDowntimeInfoFor(
		ctx,
		types.WITHDRAWAL_AND_TRANSFERS_BLOCKED_AFTER_CHAIN_OUTAGE_DURATION,
	)
	chainOutageSeenAtBlock, chainOutageExists := downtimeInfo.BlockInfo.Height,
		downtimeInfo.BlockInfo.Height > 0 && downtimeInfo.Duration > 0
	negativeTncSubaccountSeenAtBlock, negativeTncSubaccountSeenAtBlockExists, err := k.GetNegativeTncSubaccountSeenAtBlock(
		ctx,
		req.PerpetualId,
	)
	if err != nil {
		return nil, err
	}

	// Withdrawals and transfers are blocked at non-zero block iff a chain outage or negative TNC subaccount exists.
	withdrawalsAndTransfersBlockedUntilBlock := uint32(0)
	if chainOutageExists || negativeTncSubaccountSeenAtBlockExists {
		withdrawalsAndTransfersBlockedUntilBlock = max(
			chainOutageSeenAtBlock,
			negativeTncSubaccountSeenAtBlock,
		) + types.WITHDRAWAL_AND_TRANSFERS_BLOCKED_AFTER_NEGATIVE_TNC_SUBACCOUNT_SEEN_BLOCKS
	}

	return &types.QueryGetWithdrawalAndTransfersBlockedInfoResponse{
		NegativeTncSubaccountSeenAtBlock:        negativeTncSubaccountSeenAtBlock,
		ChainOutageSeenAtBlock:                  chainOutageSeenAtBlock,
		WithdrawalsAndTransfersUnblockedAtBlock: withdrawalsAndTransfersBlockedUntilBlock,
	}, nil
}
