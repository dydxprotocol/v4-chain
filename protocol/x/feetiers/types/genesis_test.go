package types_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/x/feetiers/types"
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
				Params: types.PerpetualFeeParams{
					Tiers: []*types.PerpetualFeeTier{
						{},
						{
							AbsoluteVolumeRequirement: 10,
							MakerFeePpm:               1,
							TakerFeePpm:               1,
						},
						{
							AbsoluteVolumeRequirement:      15,
							TotalVolumeShareRequirementPpm: 1,
							MakerVolumeShareRequirementPpm: 1,
							MakerFeePpm:                    2,
							TakerFeePpm:                    2,
						},
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
