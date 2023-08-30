package delaymsg_test

import (
	testutildelaymsg "github.com/dydxprotocol/v4-chain/protocol/testutil/delaymsg"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/delaymsg"
	"github.com/dydxprotocol/v4-chain/protocol/x/delaymsg/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestInitGenesis(t *testing.T) {
	tests := map[string]struct {
		genesisState *types.GenesisState
	}{
		"default genesis": {
			genesisState: types.DefaultGenesis(),
		},
		"non-default genesis (e.g. network restart)": {
			genesisState: &types.GenesisState{
				NumMessages: 20,
				DelayedMessages: []*types.DelayedMessage{
					{
						Id:          3,
						Msg:         testutildelaymsg.CreateTestAnyMsg(t),
						BlockHeight: 10,
					},
					{
						Id:          7,
						Msg:         testutildelaymsg.CreateTestAnyMsg(t),
						BlockHeight: 15,
					},
					{
						Id:          11,
						Msg:         testutildelaymsg.CreateTestAnyMsg(t),
						BlockHeight: 10,
					},
				},
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctx, delaymsgKeeper, _, _, _, _ := keeper.DelayMsgKeepers(t)
			delaymsgKeeper.InitializeForGenesis(ctx)
			delaymsg.InitGenesis(ctx, *delaymsgKeeper, *tc.genesisState)
			got := delaymsg.ExportGenesis(ctx, *delaymsgKeeper)
			require.NotNil(t, got)
			require.Equal(t, tc.genesisState, got)
		})
	}
}
