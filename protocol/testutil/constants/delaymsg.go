package constants

import (
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bridgetypes "github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
)

var (
	// MsgCompleteBridge is an example of an expected Msg type in the delaymsg module.
	TestMsg1 = &bridgetypes.MsgCompleteBridge{
		Authority: authtypes.NewModuleAddress(bridgetypes.ModuleName).String(),
		Event: bridgetypes.BridgeEvent{
			Id: 1,
		},
	}
	TestMsg2 = &bridgetypes.MsgCompleteBridge{
		Authority: authtypes.NewModuleAddress(bridgetypes.ModuleName).String(),
		Event: bridgetypes.BridgeEvent{
			Id: 2,
		},
	}
	TestMsg3 = &bridgetypes.MsgCompleteBridge{
		Authority: authtypes.NewModuleAddress(bridgetypes.ModuleName).String(),
		Event: bridgetypes.BridgeEvent{
			Id: 3,
		},
	}
	InvalidMsg = &testdata.TestMsg{Signers: []string{"invalid - no module handles this message"}}

	// Msg1Bytes, Msg2Bytes and Msg3Bytes are left as long lines for ease of byte-by-byte comparison.
	Msg1Bytes = []byte("\n&/dydxprotocol.bridge.MsgCompleteBridge\x126\n+dydx1zlefkpe3g0vvm9a4h0jf9000lmqutlh9jwjnsv\x12\x07\x08\x01\x12\x03\x12\x010") // nolint:lll
	Msg2Bytes = []byte("\n&/dydxprotocol.bridge.MsgCompleteBridge\x126\n+dydx1zlefkpe3g0vvm9a4h0jf9000lmqutlh9jwjnsv\x12\x07\x08\x02\x12\x03\x12\x010") // nolint:lll
	Msg3Bytes = []byte("\n&/dydxprotocol.bridge.MsgCompleteBridge\x126\n+dydx1zlefkpe3g0vvm9a4h0jf9000lmqutlh9jwjnsv\x12\x07\x08\x03\x12\x03\x12\x010") // nolint:lll

	AllMsgs = []sdk.Msg{TestMsg1, TestMsg2, TestMsg3}
)
