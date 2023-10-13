package types_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/rewards/types"
	"github.com/stretchr/testify/require"
)

func TestDefaultGenesis(t *testing.T) {
	genState := types.DefaultGenesis()

	expectedGenesisState := &types.GenesisState{
		Params: types.Params{
			TreasuryAccount:  "rewards_treasury",
			Denom:            "adv4tnt",
			DenomExponent:    lib.BaseDenomExponent,
			MarketId:         1,
			FeeMultiplierPpm: 990_000, // 0.99
		},
	}

	require.Equal(t, expectedGenesisState, genState)
}

func TestGenesisState_Validate(t *testing.T) {
	tests := []struct {
		desc        string
		genState    *types.GenesisState
		expectedErr string
	}{
		{
			desc:        "default is valid",
			genState:    types.DefaultGenesis(),
			expectedErr: "",
		},
		{
			desc: "invalid: empty denom",
			genState: &types.GenesisState{
				Params: types.Params{
					TreasuryAccount: "rewards_treasury",
					Denom:           "",
				},
			},
			expectedErr: "invalid denom",
		},
		{
			desc: "invalid: invalid denom",
			genState: &types.GenesisState{
				Params: types.Params{
					TreasuryAccount: "rewards_treasury",
					Denom:           "!!!!",
				},
			},
			expectedErr: "invalid denom",
		},
		{
			desc: "invalid: empty rewards treasury",
			genState: &types.GenesisState{
				Params: types.Params{
					TreasuryAccount: "",
					Denom:           "dummy",
				},
			},
			expectedErr: "treasury account cannot have empty name",
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.genState.Validate()
			if tc.expectedErr == "" {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, tc.expectedErr)
			}
		})
	}
}
