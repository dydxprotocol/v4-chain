package lib_test

import (
	"math/big"
	"testing"

	assetslib "github.com/dydxprotocol/v4-chain/protocol/x/assets/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
	"github.com/stretchr/testify/require"
)

func TestGetNetCollateralAndMarginRequirements(t *testing.T) {
	tests := map[string]struct {
		assetId     uint32
		bigQuantums *big.Int
		expectedNC  *big.Int
		expectedIMR *big.Int
		expectedMMR *big.Int
		expectedErr error
	}{
		"USDC asset. Positive Balance": {
			assetId:     types.AssetUsdc.Id,
			bigQuantums: big.NewInt(100),
			expectedNC:  big.NewInt(100),
			expectedIMR: big.NewInt(0),
			expectedMMR: big.NewInt(0),
			expectedErr: nil,
		},
		"USDC asset. Negative Balance": {
			assetId:     types.AssetUsdc.Id,
			bigQuantums: big.NewInt(-100),
			expectedNC:  big.NewInt(-100),
			expectedIMR: big.NewInt(0),
			expectedMMR: big.NewInt(0),
			expectedErr: nil,
		},
		"USDC asset. Zero Balance": {
			assetId:     types.AssetUsdc.Id,
			bigQuantums: big.NewInt(0),
			expectedNC:  big.NewInt(0),
			expectedIMR: big.NewInt(0),
			expectedMMR: big.NewInt(0),
			expectedErr: nil,
		},
		"Non USDC asset. Positive Balance": {
			assetId:     uint32(1),
			bigQuantums: big.NewInt(100),
			expectedNC:  big.NewInt(0),
			expectedIMR: big.NewInt(0),
			expectedMMR: big.NewInt(0),
			expectedErr: types.ErrNotImplementedMulticollateral,
		},
		"Non USDC asset. Negative Balance": {
			assetId:     uint32(1),
			bigQuantums: big.NewInt(-100),
			expectedNC:  big.NewInt(0),
			expectedIMR: big.NewInt(0),
			expectedMMR: big.NewInt(0),
			expectedErr: types.ErrNotImplementedMargin,
		},
		"Non USDC asset. Zero Balance": {
			assetId:     uint32(1),
			bigQuantums: big.NewInt(0),
			expectedNC:  big.NewInt(0),
			expectedIMR: big.NewInt(0),
			expectedMMR: big.NewInt(0),
			expectedErr: nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			risk, err := assetslib.GetNetCollateralAndMarginRequirements(
				tc.assetId,
				tc.bigQuantums,
			)

			require.Equal(t, tc.expectedNC, risk.NC)
			require.Equal(t, tc.expectedIMR, risk.IMR)
			require.Equal(t, tc.expectedMMR, risk.MMR)

			if tc.expectedErr != nil {
				require.ErrorIs(t, err, tc.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
