package types_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/require"
)

func TestPositionStateTransitionString(t *testing.T) {
	tests := map[string]struct {
		value          types.PositionStateTransition
		expectedResult string
	}{
		"Opened": {
			value:          types.Opened,
			expectedResult: "opened",
		},
		"Closed": {
			value:          types.Closed,
			expectedResult: "closed",
		},
		"UnexpectedError": {
			value:          types.PositionStateTransition(2),
			expectedResult: "UnexpectedStateTransitionError",
		},
	}

	// Run tests.
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := tc.value.String()
			require.Equal(t, result, tc.expectedResult)
		})
	}
}
