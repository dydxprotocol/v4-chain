package delaymsg_test

import (
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
						Msg:         []byte("foo"),
						BlockHeight: 10,
					},
					{
						Id:          7,
						Msg:         []byte("bar"),
						BlockHeight: 15,
					},
					{
						Id:          11,
						Msg:         []byte("baz"),
						BlockHeight: 10,
					},
				},
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctx, delaymsgKeeper, _, _ := keeper.DelayMsgKeepers(t)
			delaymsgKeeper.InitializeForGenesis(ctx)
			delaymsg.InitGenesis(ctx, *delaymsgKeeper, *tc.genesisState)
			got := delaymsg.ExportGenesis(ctx, *delaymsgKeeper)
			require.NotNil(t, got)
			require.Equal(t, tc.genesisState, got)
		})
	}
}
