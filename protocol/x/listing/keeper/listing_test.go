package keeper

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	"github.com/stretchr/testify/require"
)

func TestCreateMarket(t *testing.T) {

	tests := map[string]struct {
		ticker string

		expectedErr error
	}{
		"success": {
			ticker:      "TEST-USD",
			expectedErr: nil,
		},
		"invalid market": {
			ticker:      "INVALID-USD",
			expectedErr: nil, // TODO: Add expected error
		},
	}

	for name, tc := range tests {
		t.Run(
			name, func(t *testing.T) {
				mockIndexerEventManager := &mocks.IndexerEventManager{}
				ctx, keeper, _, _, _, _, _, _ := keeper.ListingKeepers(
					t,
					&mocks.BankKeeper{},
					mockIndexerEventManager,
				)

				_, err := keeper.CreateMarket(ctx, tc.ticker)
				if tc.expectedErr != nil {
					require.Error(t, err)
				} else {
					require.NoError(t, err)
				}
			},
		)
	}
}
