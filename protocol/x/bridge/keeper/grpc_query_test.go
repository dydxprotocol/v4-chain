package keeper_test

import (
	"testing"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
	delaymsgtypes "github.com/dydxprotocol/v4-chain/protocol/x/delaymsg/types"
)

func TestEventParams(t *testing.T) {
	tApp := testapp.NewTestAppBuilder().WithTesting(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.BridgeKeeper

	for name, tc := range map[string]struct {
		req *types.QueryEventParamsRequest
		res *types.QueryEventParamsResponse
		err error
	}{
		"Success": {
			req: &types.QueryEventParamsRequest{},
			res: &types.QueryEventParamsResponse{
				Params: types.DefaultGenesis().EventParams,
			},
			err: nil,
		},
		"Nil": {
			req: nil,
			res: nil,
			err: status.Error(codes.InvalidArgument, "invalid request"),
		},
	} {
		t.Run(name, func(t *testing.T) {
			res, err := k.EventParams(ctx, tc.req)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.res, res)
			}
		})
	}
}

func TestProposeParams(t *testing.T) {
	tApp := testapp.NewTestAppBuilder().WithTesting(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.BridgeKeeper

	for name, tc := range map[string]struct {
		req *types.QueryProposeParamsRequest
		res *types.QueryProposeParamsResponse
		err error
	}{
		"Success": {
			req: &types.QueryProposeParamsRequest{},
			res: &types.QueryProposeParamsResponse{
				Params: types.DefaultGenesis().ProposeParams,
			},
			err: nil,
		},
		"Nil": {
			req: nil,
			res: nil,
			err: status.Error(codes.InvalidArgument, "invalid request"),
		},
	} {
		t.Run(name, func(t *testing.T) {
			res, err := k.ProposeParams(ctx, tc.req)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.res, res)
			}
		})
	}
}

func TestSafetyParams(t *testing.T) {
	tApp := testapp.NewTestAppBuilder().WithTesting(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.BridgeKeeper

	for name, tc := range map[string]struct {
		req *types.QuerySafetyParamsRequest
		res *types.QuerySafetyParamsResponse
		err error
	}{
		"Success": {
			req: &types.QuerySafetyParamsRequest{},
			res: &types.QuerySafetyParamsResponse{
				Params: types.DefaultGenesis().SafetyParams,
			},
			err: nil,
		},
		"Nil": {
			req: nil,
			res: nil,
			err: status.Error(codes.InvalidArgument, "invalid request"),
		},
	} {
		t.Run(name, func(t *testing.T) {
			res, err := k.SafetyParams(ctx, tc.req)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.res, res)
			}
		})
	}
}

func TestAcknowledgedEventInfo(t *testing.T) {
	tApp := testapp.NewTestAppBuilder().WithTesting(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.BridgeKeeper

	for name, tc := range map[string]struct {
		req *types.QueryAcknowledgedEventInfoRequest
		res *types.QueryAcknowledgedEventInfoResponse
		err error
	}{
		"Success": {
			req: &types.QueryAcknowledgedEventInfoRequest{},
			res: &types.QueryAcknowledgedEventInfoResponse{
				Info: types.DefaultGenesis().AcknowledgedEventInfo,
			},
			err: nil,
		},
		"Nil": {
			req: nil,
			res: nil,
			err: status.Error(codes.InvalidArgument, "invalid request"),
		},
	} {
		t.Run(name, func(t *testing.T) {
			res, err := k.AcknowledgedEventInfo(ctx, tc.req)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.res, res)
			}
		})
	}
}

func TestRecognizedEventInfo(t *testing.T) {
	tApp := testapp.NewTestAppBuilder().WithTesting(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.BridgeKeeper

	for name, tc := range map[string]struct {
		req *types.QueryRecognizedEventInfoRequest
		res *types.QueryRecognizedEventInfoResponse
		err error
	}{
		"Success": {
			req: &types.QueryRecognizedEventInfoRequest{},
			res: &types.QueryRecognizedEventInfoResponse{
				Info: types.BridgeEventInfo{
					NextId:         0,
					EthBlockHeight: 0,
				},
			},
			err: nil,
		},
		"Nil": {
			req: nil,
			res: nil,
			err: status.Error(codes.InvalidArgument, "invalid request"),
		},
	} {
		t.Run(name, func(t *testing.T) {
			res, err := k.RecognizedEventInfo(ctx, tc.req)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.res, res)
			}
		})
	}
}

func TestInFlightCompleteBridgeMessages(t *testing.T) {
	for name, tc := range map[string]struct {
		events []types.BridgeEvent
		res    []types.MsgCompleteBridge
	}{
		"Success - two bridge events": {
			events: []types.BridgeEvent{
				constants.BridgeEvent_Id0_Height0,
				constants.BridgeEvent_Id1_Height0,
			},
		},
		"Success - no bridge event": {
			events: []types.BridgeEvent{},
		},
	} {
		t.Run(name, func(t *testing.T) {
			// Initialize test app.
			tApp := testapp.NewTestAppBuilder().WithTesting(t).Build()
			ctx := tApp.InitChain()
			k := tApp.App.BridgeKeeper
			delayMsgKeeper := tApp.App.DelayMsgKeeper

			// Acknowledge bridge events, for each of which there should be a delayed `MsgCompleteBridge`.
			k.AcknowledgeBridges(ctx, tc.events)
			// Also delay some other types of messages, which should not show up in the result.
			delayMsgKeeper.DelayMessageByBlocks(
				ctx,
				&constants.MsgDepositToSubaccount_Alice_To_Alice_Num0_500,
				1,
			)
			delayMsgKeeper.DelayMessageByBlocks(
				ctx,
				&constants.MsgWithdrawFromSubaccount_Carl_Num0_To_Alice_750,
				1,
			)
			delayMsgKeeper.DelayMessageByBlocks(
				ctx,
				constants.ValidMsgUpdateMarketPrices,
				1,
			)

			// Get in flight complete bridge messages.
			got := k.GetInFlightCompleteBridgeMessages(ctx)

			// Verify response contains and only contains expected delayed `MsgCompleteBridge`s.
			delayMsgAuthority := authtypes.NewModuleAddress(delaymsgtypes.ModuleName).String()
			expectedCompleteBridgeMsgs := make([]types.MsgCompleteBridge, len(tc.events))
			for i, event := range tc.events {
				expectedCompleteBridgeMsgs[i] = types.MsgCompleteBridge{
					Authority: delayMsgAuthority,
					Event:     event,
				}
			}
			require.Equal(t, expectedCompleteBridgeMsgs, got)
		})
	}
}
