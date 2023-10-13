package types_test

import (
	"fmt"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/testutil/delaymsg"
	"github.com/dydxprotocol/v4-chain/protocol/x/delaymsg/types"
	"github.com/stretchr/testify/require"
)

func TestGenesisState_Validate(t *testing.T) {
	tests := map[string]struct {
		genState    *types.GenesisState
		expectedErr error
	}{
		"default is valid": {
			genState: types.DefaultGenesis(),
		},
		"invalid delayed message id": {
			genState: &types.GenesisState{
				DelayedMessages: []*types.DelayedMessage{
					{
						Id:          2,
						Msg:         delaymsg.CreateTestAnyMsg(t),
						BlockHeight: 1,
					},
				},
				NextDelayedMessageId: 2,
			},
			expectedErr: fmt.Errorf(
				"delayed message id cannot be greater than or equal to next id: Invalid genesis state",
			),
		},
		"invalid delayed message - no message bytes": {
			genState: &types.GenesisState{
				DelayedMessages: []*types.DelayedMessage{
					{
						Id:          1,
						BlockHeight: 1,
					},
				},
				NextDelayedMessageId: 2,
			},
			expectedErr: fmt.Errorf(
				"invalid delayed message at index 0 with id 1: Delayed msg is nil: Invalid genesis state",
			),
		},
		"invalid delayed message - empty message": {
			genState: &types.GenesisState{
				DelayedMessages: []*types.DelayedMessage{
					{
						Id:          1,
						BlockHeight: 1,
						Msg:         nil,
					},
				},
				NextDelayedMessageId: 2,
			},
			expectedErr: fmt.Errorf(
				"invalid delayed message at index 0 with id 1: Delayed msg is nil: Invalid genesis state",
			),
		},
		"invalid genesis state - duplicate message id": {
			genState: &types.GenesisState{
				DelayedMessages: []*types.DelayedMessage{
					{
						Id:          1,
						BlockHeight: 1,
						Msg:         delaymsg.CreateTestAnyMsg(t),
					},
					{
						Id:          1,
						BlockHeight: 2,
						Msg:         delaymsg.CreateTestAnyMsg(t),
					},
				},
				NextDelayedMessageId: 2,
			},
			expectedErr: fmt.Errorf("duplicate delayed message id: Invalid genesis state"),
		},
		"valid genesis state - multiple noncontiguous delayed messages": {
			genState: &types.GenesisState{
				DelayedMessages: []*types.DelayedMessage{
					{
						Id:          0,
						BlockHeight: 34,
						Msg:         delaymsg.CreateTestAnyMsg(t),
					},
					{
						Id:          2,
						BlockHeight: 5,
						Msg:         delaymsg.CreateTestAnyMsg(t),
					},
					{
						Id:          4,
						BlockHeight: 88,
						Msg:         delaymsg.CreateTestAnyMsg(t),
					},
				},
				NextDelayedMessageId: 5,
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := tc.genState.Validate()
			if tc.expectedErr == nil {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tc.expectedErr.Error())
			}
		})
	}
}
