package keeper_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/x/listing/types"
	"github.com/stretchr/testify/require"

	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
)

func TestQueryMarketsHardCap(t *testing.T) {
	tests := map[string]struct {
		hardCap uint32
	}{
		"Hard cap: 0": {
			hardCap: 0,
		},
		"Hard cap: 100": {
			hardCap: 100,
		},
	}

	for name, tc := range tests {
		t.Run(
			name, func(t *testing.T) {
				tApp := testapp.NewTestAppBuilder(t).Build()
				ctx := tApp.InitChain()
				k := tApp.App.ListingKeeper

				// set hard cap for markets for test
				err := k.SetMarketsHardCap(ctx, tc.hardCap)
				require.NoError(t, err)

				// query hard cap for markets
				resp, err := k.MarketsHardCap(ctx, &types.QueryMarketsHardCap{})
				require.NoError(t, err)
				require.Equal(t, resp.HardCap, tc.hardCap)
			},
		)
	}
}
