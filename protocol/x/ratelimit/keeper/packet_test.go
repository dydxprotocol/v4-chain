package keeper_test

import (
	"testing"

	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/stretchr/testify/require"
)

func TestPendingPacket(t *testing.T) {
	testChannelId := "channel-0"
	testSequence := uint64(20)
	testSequence2 := uint64(22)

	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.RatelimitKeeper

	// Set pending packet in state
	k.SetPendingSendPacket(ctx, testChannelId, testSequence)
	k.SetPendingSendPacket(ctx, testChannelId, testSequence2)

	// Test HasPendingSendPacket
	require.True(t, k.HasPendingSendPacket(ctx, testChannelId, testSequence))
	require.True(t, k.HasPendingSendPacket(ctx, testChannelId, testSequence2))
	require.False(t, k.HasPendingSendPacket(ctx, "non-existent-channel", testSequence))
	require.False(t,
		k.HasPendingSendPacket(
			ctx, testChannelId,
			42, // non-existent sequence number
		),
	)

	// Remove pending packet from state
	k.RemovePendingSendPacket(
		ctx,
		testChannelId,
		testSequence,
	)

	require.False(t, k.HasPendingSendPacket(ctx, testChannelId, testSequence)) // Removed
	require.True(t, k.HasPendingSendPacket(ctx, testChannelId, testSequence2))
	require.False(t, k.HasPendingSendPacket(ctx, "non-existent-channel", testSequence))
	require.False(t,
		k.HasPendingSendPacket(
			ctx, testChannelId,
			42, // non-existent sequence number
		),
	)
}

// TODO(CORE-856): Improve coverage for remaining functions in packet.go
