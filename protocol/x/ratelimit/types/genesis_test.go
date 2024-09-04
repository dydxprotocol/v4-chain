package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/ratelimit/types"
)

func TestDefaultGenesis(t *testing.T) {
	genesis := types.DefaultGenesis()

	require.NotNil(t, genesis)
	require.Equal(t, len(genesis.LimitParamsList), 1)
	require.Equal(t, types.DefaultSDaiRateLimitParams(), genesis.LimitParamsList[0])
}

func TestGenesisState_Validate(t *testing.T) {
	tests := []struct {
		name      string
		genesis   types.GenesisState
		expectErr bool
	}{
		{
			name: "valid genesis",
			genesis: types.GenesisState{
				LimitParamsList: []types.LimitParams{
					types.DefaultSDaiRateLimitParams(),
				},
			},
			expectErr: false,
		},
		{
			name: "invalid limit params",
			genesis: types.GenesisState{
				LimitParamsList: []types.LimitParams{
					{Denom: ""},
				},
			},
			expectErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.genesis.Validate()
			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
