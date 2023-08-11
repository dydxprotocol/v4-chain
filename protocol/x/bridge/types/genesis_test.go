package types_test

import (
	"testing"
	"time"

	"github.com/dydxprotocol/v4/x/bridge/types"
	"github.com/stretchr/testify/require"
)

func TestGenesisState_Validate(t *testing.T) {
	tests := map[string]struct {
		genState *types.GenesisState
		err      error
	}{
		"default is valid": {
			genState: types.DefaultGenesis(),
			err:      nil,
		},
		"valid genesis state": {
			genState: &types.GenesisState{},
			err:      nil,
		},
		"invalid ProposeParams": {
			genState: &types.GenesisState{
				ProposeParams: types.ProposeParams{
					ProposeDelayDuration: time.Duration(-1),
				},
			},
			err: types.ErrNegativeDuration,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := tc.genState.Validate()
			if tc.err == nil {
				require.NoError(t, err)
			} else {
				require.ErrorIs(t, tc.err, err)
			}
		})
	}
}
