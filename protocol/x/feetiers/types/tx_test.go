package types_test

import (
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	types "github.com/dydxprotocol/v4-chain/protocol/x/feetiers/types"
	"github.com/stretchr/testify/require"
	"testing"
)

var (
	validAuthority = constants.BobAccAddress.String()
)

func TestValidateBasic(t *testing.T) {
	tests := map[string]struct {
		msg         types.MsgUpdatePerpetualFeeParams
		expectedErr error
	}{
		"Success": {
			msg: types.MsgUpdatePerpetualFeeParams{
				Authority: validAuthority,
				Params: types.PerpetualFeeParams{
					Tiers: []*types.PerpetualFeeTier{
						{
							AbsoluteVolumeRequirement:      0,
							TotalVolumeShareRequirementPpm: 0,
							MakerVolumeShareRequirementPpm: 0,
						},
					},
				},
			},
		},
		"Failure: Invalid authority": {
			msg: types.MsgUpdatePerpetualFeeParams{
				Authority: "",
			},
			expectedErr: types.ErrInvalidAuthority,
		},
		"Failure: Invalid params": {
			msg: types.MsgUpdatePerpetualFeeParams{
				Authority: validAuthority,
				Params: types.PerpetualFeeParams{
					Tiers: []*types.PerpetualFeeTier{}, // invalid - empty
				},
			},
			expectedErr: types.ErrNoTiersExist,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := tc.msg.ValidateBasic()
			if tc.expectedErr == nil {
				require.NoError(t, err)
			} else {
				require.ErrorIs(t, err, tc.expectedErr)
			}
		})
	}
}
