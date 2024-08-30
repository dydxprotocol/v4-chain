package delaymsg

import (
	"testing"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	delaymsgtypes "github.com/dydxprotocol/v4-chain/protocol/x/delaymsg/types"
	"github.com/stretchr/testify/require"
)

// CreateTestAnyMsg returns an encoded Any object for an sdk.Msg for testing. This is useful
// when a valid message is needed, but the message will never be executed.
func CreateTestAnyMsg(t testing.TB) *codectypes.Any {
	any, err := codectypes.NewAnyWithValue(constants.TestMsg1)
	require.NoError(t, err)
	return any
}

// FilterDelayedMsgsByType filters a slice of DelayedMessage by a specific sdk.Msg type.
func FilterDelayedMsgsByType[T sdk.Msg](
	t testing.TB,
	delayedmsgs []*delaymsgtypes.DelayedMessage,
) []*delaymsgtypes.DelayedMessage {
	var filtered []*delaymsgtypes.DelayedMessage
	for _, delayedmsg := range delayedmsgs {
		msg, err := delayedmsg.GetMessage()
		require.NoError(t, err)
		if _, ok := msg.(T); ok {
			filtered = append(filtered, delayedmsg)
		}
	}
	return filtered
}
