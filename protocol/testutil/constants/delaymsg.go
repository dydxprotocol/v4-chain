package constants

import (
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bridgetypes "github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
)

var (
	TestMsg1 = &testdata.TestMsg{Signers: []string{"meh"}}
	TestMsg2 = &testdata.TestMsg{Signers: []string{"blah"}}
	TestMsg3 = &testdata.TestMsg{Signers: []string{"nerp"}}
	// Use a real bridge event as a test message that passes validation for a live network.
	TestMsg4 = &bridgetypes.MsgCompleteBridge{
		Authority: authtypes.NewModuleAddress(bridgetypes.ModuleName).String(),
		Event: bridgetypes.BridgeEvent{
			Id: 1,
		},
	}

	Msg1Bytes = []byte("\n\x0f\x2ftestpb.TestMsg\x12\x05\n\x03meh")
	Msg2Bytes = []byte("\n\x0f\x2ftestpb.TestMsg\x12\x06\n\x04blah")
	Msg3Bytes = []byte("\n\x0f\x2ftestpb.TestMsg\x12\x06\n\x04nerp")
	Msg4Bytes = []byte("\n&/dydxprotocol.bridge.MsgCompleteBridge\x126\n+dydx1zlefkpe3g0vvm9a4h0jf9000lmqutlh9jwjnsv\x12\x07\x08\x01\x12\x03\x12\x010")

	TestMsgAuthorities = []string{"meh", "blah", "nerp", authtypes.NewModuleAddress(bridgetypes.ModuleName).String()}

	AllMsgs = []sdk.Msg{TestMsg1, TestMsg2, TestMsg3}
)
