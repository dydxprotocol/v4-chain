package time_test

import (
	"testing"
	"time"

	libtime "github.com/dydxprotocol/v4-chain/protocol/lib/time"
	"github.com/stretchr/testify/require"
)

func TestMustParseDuration(t *testing.T) {
	tests := map[string]struct {
		input          string
		expectPanic    bool
		expectedOutput time.Duration
	}{
		"valid: zero": {
			input:          "0",
			expectedOutput: time.Duration(0),
		},
		"valid: zero with unit": {
			input:          "0ms",
			expectedOutput: time.Duration(0),
		},
		"valid: positive decimal": {
			input:          "1.137120913s",
			expectedOutput: time.Duration(1137120913) * time.Nanosecond,
		},
		"valid: negative decimal": {
			input:          "-0.123457913s",
			expectedOutput: time.Duration(-123457913) * time.Nanosecond,
		},
		"invalid: empty string": {
			input:       "", // empty input
			expectPanic: true,
		},
		"invalid: invalid format": {
			input:       "1.137120913", // no unit
			expectPanic: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if tc.expectPanic {
				require.Panics(
					t,
					func() {
						v := libtime.MustParseDuration(tc.input)
						require.Nil(t, v)
					},
				)
				return
			}

			v := libtime.MustParseDuration(tc.input)
			require.Equal(t, tc.expectedOutput, v)
		})
	}
}
