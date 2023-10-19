package keeper_test

import (
	"testing"

	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
	"github.com/stretchr/testify/require"
)

func TestGetAcknowledgedEventInfo(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.BridgeKeeper

	require.Equal(
		t,
		types.BridgeEventInfo{
			NextId:         0,
			EthBlockHeight: 0,
		},
		k.GetAcknowledgedEventInfo(ctx),
	)
}

func TestSetAcknowledgedEventInfo(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.BridgeKeeper

	info1 := types.BridgeEventInfo{
		NextId:         111,
		EthBlockHeight: 0,
	}
	info2 := types.BridgeEventInfo{
		NextId:         222,
		EthBlockHeight: 111,
	}

	err := k.SetAcknowledgedEventInfo(ctx, info1)
	require.NoError(t, err)
	require.Equal(t, info1, k.GetAcknowledgedEventInfo(ctx))
	err = k.SetAcknowledgedEventInfo(ctx, info2)
	require.NoError(t, err)
	require.Equal(t, info2, k.GetAcknowledgedEventInfo(ctx))
}
