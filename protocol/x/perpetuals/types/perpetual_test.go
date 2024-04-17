package types_test

import (
	"testing"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/perpetuals/types"
	"github.com/stretchr/testify/require"
)

func TestPerpetualParams_Validate(t *testing.T) {
	tests := []struct {
		desc        string
		params      types.PerpetualParams
		expectedErr string
	}{
		{
			desc: "Valid param",
			params: types.PerpetualParams{
				Ticker:            "test",
				DefaultFundingPpm: 1_000_000,
			},
			expectedErr: "",
		},
		{
			desc: "Empty ticker",
			params: types.PerpetualParams{
				Ticker:            "",
				DefaultFundingPpm: 1_000_000,
			},
			expectedErr: "Ticker must be non-empty string",
		},
		{
			desc: "Invalid DefaultFundingPpm",
			params: types.PerpetualParams{
				Ticker:            "test",
				DefaultFundingPpm: 100_000_000,
			},
			expectedErr: "DefaultFundingPpm magnitude exceeds maximum value",
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.params.Validate()
			if tc.expectedErr == "" {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, tc.expectedErr)
			}
		})
	}
}
