package types_test

import (
	"testing"
	"time"

	"github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
	"github.com/stretchr/testify/require"
)

func TestGenesisState_Validate(t *testing.T) {
	tests := map[string]struct {
		genState *types.GenesisState
		err      string
	}{
		"default is valid": {
			genState: types.DefaultGenesis(),
		},
		"invalid EventParams": {
			genState: &types.GenesisState{
				EventParams: types.EventParams{
					Denom:      "7coin",
					EthChainId: 2,
					EthAddress: "test",
				},
			},
			err: "invalid denom",
		},
		"invalid ProposeParams": {
			genState: &types.GenesisState{
				EventParams: types.EventParams{
					Denom:      "test-coin",
					EthChainId: 2,
					EthAddress: "test",
				},
				ProposeParams: types.ProposeParams{
					ProposeDelayDuration: time.Duration(-1),
				},
			},
			err: types.ErrNegativeDuration.Error(),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := tc.genState.Validate()
			if tc.err == "" {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, tc.err)
			}
		})
	}
}
