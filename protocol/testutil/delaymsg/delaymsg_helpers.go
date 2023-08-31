package delaymsg

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/stretchr/testify/require"
	"testing"
)

// CreateTestAnyMsg returns an encoded Any object for an sdk.Msg for testing. This is useful
// when a valid message is needed, but the message will never be executed.
func CreateTestAnyMsg(t *testing.T) *codectypes.Any {
	any, err := codectypes.NewAnyWithValue(constants.TestMsg1)
	require.NoError(t, err)
	return any
}

// EncodeMessageToAny converts a message to an Any object for protobuf encoding.
func EncodeMessageToAny(t *testing.T, msg sdk.Msg) *codectypes.Any {
	any, err := codectypes.NewAnyWithValue(msg)
	require.NoError(t, err)
	return any
}
