package types_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"

	"github.com/stretchr/testify/require"
)

func TestMsgAddPremiumVotes(t *testing.T) {
	sample := types.NewFundingPremium(uint32(0), int32(10))
	samples := []types.FundingPremium{*sample}
	msg := types.NewMsgAddPremiumVotes(samples)

	require.Equal(t, uint32(0), sample.PerpetualId)
	require.Equal(t, int32(10), sample.PremiumPpm)
	require.Equal(t, samples, msg.Votes)
}

func TestValidateBasic(t *testing.T) {
	errStr := "premium votes must be sorted by perpetual id in ascending order and cannot " +
		"contain duplicates: MsgAddPremiumVotes is invalid"

	tests := map[string]struct {
		samples []types.FundingPremium

		expectedErr bool
	}{
		"Error: duplicate perpetual ids": {
			samples: []types.FundingPremium{
				{PerpetualId: 1},
				{PerpetualId: 1},
			},
			expectedErr: true,
		},
		"Error: desending perpetual ids": {
			samples: []types.FundingPremium{
				{PerpetualId: 2},
				{PerpetualId: 3},
				{PerpetualId: 4},
				{PerpetualId: 1},
			},
			expectedErr: true,
		},
		"No error: empty samples": {
			samples: []types.FundingPremium{},
		},
		"No error: valid ordering": {
			samples: []types.FundingPremium{
				{PerpetualId: 1},
				{PerpetualId: 3},
				{PerpetualId: 99},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			msg := types.NewMsgAddPremiumVotes(tc.samples)
			err := msg.ValidateBasic()
			if tc.expectedErr {
				require.ErrorContains(t, err, errStr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
