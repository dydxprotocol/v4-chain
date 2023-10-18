package ibc_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/testutil/ibc"
	assettypes "github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
	"github.com/stretchr/testify/require"
)

func TestDenomTraceToIBCDenom_Success(t *testing.T) {
	tests := []struct {
		denomTrace string
		expected   string
	}{
		// Check `transfer/channel-0/uusdc` results in expected ibc hash.
		{
			denomTrace: "transfer/channel-0/uusdc",
			expected:   assettypes.UusdcDenom,
		},
		// The following test cases and results are obtained from the private testnet.
		{
			denomTrace: "transfer/channel-8/uusdc",
			expected:   "ibc/39549F06486BACA7494C9ACDD53CDD30AA9E723AB657674DBD388F867B61CA7B",
		},
		{
			denomTrace: "transfer/channel-3/uosmo",
			expected:   "ibc/47BD209179859CDE4A2806763D7189B6E6FE13A17880FE2B42DE1E6C1E329E23",
		},
		{
			denomTrace: "transfer/channel-8/transfer/channel-6/uosmo", // two hops
			expected:   "ibc/7DF7B90D2F1FC60D07C866F660FC352CEAFFE0F9AA762EF98ECF6D32675938CB",
		},
		{
			denomTrace: "transfer/channel-0/uosmo",
			expected:   "ibc/ED07A3391A112B175915CD8FAF43A2DA8E4790EDE12566649D0C2F97716B8518",
		},
	}

	for _, tc := range tests {
		t.Run(tc.denomTrace, func(t *testing.T) {
			result, err := ibc.DenomTraceToIBCDenom(tc.denomTrace)
			require.NoError(t, err)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestDenomTraceToIBCDenom_Failure(t *testing.T) {
	tests := map[string]struct {
		denomTrace  string
		expectedErr string
	}{
		"invalid channel format with parenthesis": {
			denomTrace:  "transfer/(channel-1)/uusdc",
			expectedErr: "invalid denom trace 'transfer/(channel-1)/uusdc' is parsed into empty path or base denom",
		},
		"invalid channel id (not a number)": {
			denomTrace:  "transfer/channel-ABC/uusdc",
			expectedErr: "invalid denom trace 'transfer/channel-ABC/uusdc' is parsed into empty path or base denom",
		},
		"missing channel identifier": {
			denomTrace:  "transfer/uusdc",
			expectedErr: "invalid denom trace 'transfer/uusdc' is parsed into empty path or base denom",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := ibc.DenomTraceToIBCDenom(tc.denomTrace)
			require.ErrorContains(t, err, tc.expectedErr)
		})
	}
}
