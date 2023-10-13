package types_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/x/delaymsg/types"
	"github.com/stretchr/testify/require"
)

func TestModuleKeys(t *testing.T) {
	require.Equal(t, "delaymsg", types.ModuleName)
	require.Equal(t, "delaymsg", types.StoreKey)
}

func TestStateKeys(t *testing.T) {
	require.Equal(t, "BlockMsgIds:", types.BlockMessageIdsPrefix)
	require.Equal(t, "Msg:", types.DelayedMessageKeyPrefix)
	require.Equal(t, "NextDelayedMessageId", types.NextDelayedMessageIdKey)
}
