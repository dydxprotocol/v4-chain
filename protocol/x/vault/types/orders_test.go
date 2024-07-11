package types_test

import (
	math "math"
	"testing"

	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
	"github.com/stretchr/testify/require"
)

func TestGetVaultClobOrderClientId(t *testing.T) {
	tests := map[string]struct {
		/* --- Setup --- */
		// side.
		side clobtypes.Order_Side
		// layer.
		layer uint8

		/* --- Expectations --- */
		// Expected client ID.
		expectedClientId uint32
	}{
		"Buy, Layer 0": {
			side:             clobtypes.Order_SIDE_BUY, // 0<<31
			layer:            0,                        // 0<<23
			expectedClientId: 0<<31 | 0<<1,
		},
		"Sell, Layer 0": {
			side:             clobtypes.Order_SIDE_SELL, // 1<<31
			layer:            0,                         // 0<<23
			expectedClientId: 1<<31 | 0<<23,
		},
		"Buy, Layer 1": {
			side:             clobtypes.Order_SIDE_BUY, // 0<<31
			layer:            1,                        // 1<<23
			expectedClientId: 0<<31 | 1<<23,
		},
		"Sell, Layer 1": {
			side:             clobtypes.Order_SIDE_SELL, // 1<<31
			layer:            1,                         // 1<<23
			expectedClientId: 1<<31 | 1<<23,
		},
		"Buy, Layer 2": {
			side:             clobtypes.Order_SIDE_BUY, // 0<<31
			layer:            2,                        // 2<<23
			expectedClientId: 0<<31 | 2<<23,
		},
		"Sell, Layer 2": {
			side:             clobtypes.Order_SIDE_SELL, // 1<<31
			layer:            2,                         // 2<<23
			expectedClientId: 1<<31 | 2<<23,
		},
		"Buy, Layer 123": {
			side:             clobtypes.Order_SIDE_BUY, // 0<<31
			layer:            123,                      // 123<<23
			expectedClientId: 0<<31 | 123<<23,
		},
		"Sell, Layer 123": {
			side:             clobtypes.Order_SIDE_SELL, // 1<<31
			layer:            123,                       // 123<<23
			expectedClientId: 1<<31 | 123<<23,
		},
		"Buy, Layer Max Uint8": {
			side:             clobtypes.Order_SIDE_BUY, // 0<<31
			layer:            math.MaxUint8,            // 255<<23
			expectedClientId: 0<<31 | 255<<23,
		},
		"Sell, Layer Max Uint8": {
			side:             clobtypes.Order_SIDE_SELL, // 1<<31
			layer:            math.MaxUint8,             // 255<<23
			expectedClientId: 1<<31 | 255<<23,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			require.Equal(
				t,
				tc.expectedClientId,
				types.GetVaultClobOrderClientId(
					tc.side,
					tc.layer,
				),
			)
		})
	}
}
