package types_test

import (
	"testing"
	time "time"

	"github.com/dydxprotocol/v4-chain/protocol/x/blocktime/types"
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
			genState: &types.GenesisState{
				Params: types.DowntimeParams{
					Durations: []time.Duration{
						2 * time.Minute,
					},
				},
			},
			err: nil,
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
