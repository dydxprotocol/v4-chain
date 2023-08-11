package keeper_test

import (
	"testing"

	testapp "github.com/dydxprotocol/v4/testutil/app"
	"github.com/stretchr/testify/require"
)

func TestGetNextAcknowledgedEventId(t *testing.T) {
	tApp := testapp.NewTestAppBuilder().WithTesting(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.BridgeKeeper

	require.Equal(t, uint32(0), k.GetNextAcknowledgedEventId(ctx))
}

func TestSetNextAcknowledgedEventId(t *testing.T) {
	tApp := testapp.NewTestAppBuilder().WithTesting(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.BridgeKeeper

	id1 := uint32(111)
	id2 := uint32(222)

	k.SetNextAcknowledgedEventId(ctx, id1)
	require.Equal(t, id1, k.GetNextAcknowledgedEventId(ctx))
	k.SetNextAcknowledgedEventId(ctx, id2)
	require.Equal(t, id2, k.GetNextAcknowledgedEventId(ctx))
}
