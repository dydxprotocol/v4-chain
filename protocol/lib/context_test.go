package lib_test

import (
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTxHash(t *testing.T) {
	tests := map[string]struct {
		// parameters
		value []byte

		// expectations
		expected lib.TxHash
	}{
		"Empty": {
			expected: lib.TxHash("E3B0C44298FC1C149AFBF4C8996FB92427AE41E4649B934CA495991B7852B855"),
		},
		"0xCAFEBABE": {
			value:    []byte{0xCA, 0xFE, 0xBA, 0xBE},
			expected: lib.TxHash("65AB12A8FF3263FBC257E5DDF0AA563C64573D0BAB1F1115B9B107834CFA6971"),
		},
		"0xDEADBEEF": {
			value:    []byte{0xDE, 0xAD, 0xBE, 0xEF},
			expected: lib.TxHash("5F78C33274E43FA9DE5659265C1D917E25C03722DCB0B8D27DB8D5FEAA813953"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			require.Equal(t, lib.GetTxHash(tc.value), tc.expected)
		})
	}
}
