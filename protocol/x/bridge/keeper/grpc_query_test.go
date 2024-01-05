package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	pricestest "github.com/dydxprotocol/v4-chain/protocol/testutil/prices"
	"github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
	delaymsgtypes "github.com/dydxprotocol/v4-chain/protocol/x/delaymsg/types"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	sendingtypes "github.com/dydxprotocol/v4-chain/protocol/x/sending/types"
)

var (
	DelayMsgAuthority = delaymsgtypes.ModuleAddress.String()
)

func TestEventParams(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
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
	tApp := testapp.NewTestAppBuilder(t).Build()
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
	tApp := testapp.NewTestAppBuilder(t).Build()
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
	tApp := testapp.NewTestAppBuilder(t).Build()
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
	tApp := testapp.NewTestAppBuilder(t).Build()
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

func TestDelayedCompleteBridgeMessages(t *testing.T) {
	for name, tc := range map[string]struct {
		events []types.BridgeEvent
		res    []types.MsgCompleteBridge
	}{
		"Success - no bridge event": {
			events: []types.BridgeEvent{},
		},
		"Success - two bridge events": {
			events: []types.BridgeEvent{
				constants.BridgeEvent_Id0_Height0,
				constants.BridgeEvent_Id1_Height0,
			},
		},
		"Success - five bridge events": {
			events: []types.BridgeEvent{
				constants.BridgeEvent_Id0_Height0,
				constants.BridgeEvent_Id1_Height0,
				constants.BridgeEvent_Id2_Height1,
				constants.BridgeEvent_Id3_Height3,
				constants.BridgeEvent_Id55_Height15,
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			// Initialize test app.
			tApp := testapp.NewTestAppBuilder(t).Build()
			ctx := tApp.InitChain()
			k := tApp.App.BridgeKeeper
			delayMsgKeeper := tApp.App.DelayMsgKeeper

			// Acknowledge bridge events, for each of which there should be a delayed `MsgCompleteBridge`.
			err := k.AcknowledgeBridges(ctx, tc.events)
			require.NoError(t, err)
			// Also delay some other types of messages, which should not show up in the result.
			_, err = delayMsgKeeper.DelayMessageByBlocks(
				ctx,
				sendingtypes.NewMsgSendFromModuleToAccount(
					DelayMsgAuthority,
					types.ModuleName,
					constants.AliceAccAddress.String(),
					sdk.NewCoin("adv4tnt", sdkmath.NewInt(100)),
				),
				100,
			)
			require.NoError(t, err)
			_, err = delayMsgKeeper.DelayMessageByBlocks(
				ctx,
				&pricestypes.MsgUpdateMarketParam{
					Authority:   DelayMsgAuthority,
					MarketParam: pricestest.GenerateMarketParamPrice().Param,
				},
				123,
			)
			require.NoError(t, err)

			// Construct expected responses.
			delayMsgAuthority := DelayMsgAuthority
			blockOfExecution := k.GetSafetyParams(ctx).DelayBlocks + uint32(ctx.BlockHeight())
			expectedMsgs := make([]types.DelayedCompleteBridgeMessage, 0)
			expectedMsgsByAddress := make(map[string][]types.DelayedCompleteBridgeMessage)
			for _, event := range tc.events {
				DelayedMsg := types.DelayedCompleteBridgeMessage{
					Message: types.MsgCompleteBridge{
						Authority: delayMsgAuthority,
						Event:     event,
					},
					BlockHeight: blockOfExecution,
				}

				expectedMsgs = append(expectedMsgs, DelayedMsg)

				expectedMsgsByAddress[event.Address] = append(
					expectedMsgsByAddress[event.Address],
					DelayedMsg,
				)
			}

			// Get all delayed complete bridge messages and verify they are as expected.
			msgs := k.GetDelayedCompleteBridgeMessages(ctx, "")
			require.Equal(t, expectedMsgs, msgs)

			// Get delayed complete bridge messages for each address and verify they are as expected.
			for address, expectedMsgsForAddr := range expectedMsgsByAddress {
				msgs := k.GetDelayedCompleteBridgeMessages(ctx, address)
				require.Equal(t, expectedMsgsForAddr, msgs)
			}
		})
	}
}
