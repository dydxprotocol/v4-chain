package types_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/x/feetiers/types"
	"github.com/stretchr/testify/require"
)

func TestPerpetualFeeParams_Validate(t *testing.T) {
	tests := map[string]struct {
		params *types.PerpetualFeeParams
		err    error
	}{
		"no tiers is invalid": {
			params: &types.PerpetualFeeParams{
				Tiers: []*types.PerpetualFeeTier{},
			},
			err: types.ErrNoTiersExist,
		},
		"first tier with requirements is invalid": {
			params: &types.PerpetualFeeParams{
				Tiers: []*types.PerpetualFeeTier{
					{
						AbsoluteVolumeRequirement: 10,
						MakerFeePpm:               1,
						TakerFeePpm:               1,
					},
				},
			},
			err: types.ErrInvalidFirstTierRequirements,
		},
		"tiers out of order": {
			params: &types.PerpetualFeeParams{
				Tiers: []*types.PerpetualFeeTier{
					{},
					{
						AbsoluteVolumeRequirement: 10,
					},
					{
						AbsoluteVolumeRequirement: 5,
					},
				},
			},
			err: types.ErrTiersOutOfOrder,
		},
		"maker rebate exceeds taker fee": {
			params: &types.PerpetualFeeParams{
				Tiers: []*types.PerpetualFeeTier{
					{},
					{
						AbsoluteVolumeRequirement: 5,
						MakerFeePpm:               -2,
						TakerFeePpm:               3,
					},
					{
						AbsoluteVolumeRequirement: 10,
						TakerFeePpm:               1,
					},
				},
			},
			err: types.ErrInvalidFee,
		},
		"valid maker rebate": {
			params: &types.PerpetualFeeParams{
				Tiers: []*types.PerpetualFeeTier{
					{
						MakerFeePpm: -2,
						TakerFeePpm: 4,
					},
					{
						AbsoluteVolumeRequirement: 5,
						MakerFeePpm:               -2,
						TakerFeePpm:               4,
					},
					{
						AbsoluteVolumeRequirement: 10,
						TakerFeePpm:               3,
					},
				},
			},
			err: nil,
		},
		"maker rebate cannot coexist with no taker fee": {
			params: &types.PerpetualFeeParams{
				Tiers: []*types.PerpetualFeeTier{
					{
						MakerFeePpm: -2,
						TakerFeePpm: 0,
					},
				},
			},
			err: types.ErrInvalidFee,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := tc.params.Validate()
			if tc.err == nil {
				require.NoError(t, err)
			} else {
				require.ErrorIs(t, tc.err, err)
			}
		})
	}
}
