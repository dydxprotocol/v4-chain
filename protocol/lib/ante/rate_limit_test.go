package ante_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	libante "github.com/dydxprotocol/v4-chain/protocol/lib/ante"
	"github.com/stretchr/testify/require"
)

func TestShouldRateLimit(t *testing.T) {
	ctx := sdk.Context{}
	ctx.WithExecMode(
		sdk.ExecModeCheck,
	)

	tests := []struct {
		name     string
		expected bool
		ctx      sdk.Context
	}{
		{
			name:     "returns true if the context is CheckTx",
			expected: true,
			ctx:      ctx.WithIsCheckTx(true),
		},
		{
			name:     "returns false if the context is ReCheckTx",
			expected: false,
			ctx:      ctx.WithIsReCheckTx(true),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(
				t,
				tc.expected,
				libante.ShouldRateLimit(tc.ctx),
			)
		})
	}
}
