package types_test

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	feetierstypes "github.com/dydxprotocol/v4-chain/protocol/x/feetiers/types"

	"github.com/dydxprotocol/v4-chain/protocol/testutil/daemons/pricefeed"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/delaymsg"
	"github.com/dydxprotocol/v4-chain/protocol/x/delaymsg/types"
	"github.com/stretchr/testify/require"
	"testing"
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
				NumMessages: 2,
				DelayedMessages: []*types.DelayedMessage{
					{
						Id:          2,
						Msg:         delaymsg.CreateTestAnyMsg(t),
						BlockHeight: 1,
					},
				},
			},
			expectedErr: fmt.Errorf(
				"delayed message id exceeds total number of messages: Invalid genesis state",
			),
		},
		"invalid delayed message - no message bytes": {
			genState: &types.GenesisState{
				NumMessages: 2,
				DelayedMessages: []*types.DelayedMessage{
					{
						Id:          1,
						BlockHeight: 1,
					},
				},
			},
			expectedErr: fmt.Errorf(
				"invalid delayed message at index 0 with id 1: Delayed msg is nil: Invalid genesis state",
			),
		},
		"invalid delayed message - empty message": {
			genState: &types.GenesisState{
				NumMessages: 2,
				DelayedMessages: []*types.DelayedMessage{
					{
						Id:          1,
						BlockHeight: 1,
						Msg:         nil,
					},
				},
			},
			expectedErr: fmt.Errorf(
				"invalid delayed message at index 0 with id 1: Delayed msg is nil: Invalid genesis state",
			),
		},
		"invalid genesis state - duplicate message id": {
			genState: &types.GenesisState{
				NumMessages: 2,
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
			},
			expectedErr: fmt.Errorf("duplicate delayed message id: Invalid genesis state"),
		},
		"valid genesis state - multiple noncontiguous delayed messages": {
			genState: &types.GenesisState{
				NumMessages: 5,
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

// TestDefaultGenesis validates the literal contents of the default genesis.
func TestDefaultGenesis(t *testing.T) {
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(interfaceRegistry)
	feetierstypes.RegisterInterfaces(interfaceRegistry)

	defaultGenesisString := string(cdc.MustMarshalJSON(types.DefaultGenesis()))

	expected := pricefeed.ReadJsonTestFile(t, "default_genesis.json")
	require.Equal(t, expected, defaultGenesisString)
}
