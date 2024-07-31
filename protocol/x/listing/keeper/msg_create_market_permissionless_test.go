package keeper_test

import (
	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/x/listing/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/listing/types"
	subaccounttypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/require"

	"testing"
)

func TestMsgCreateMarketPermissionless(t *testing.T) {
	tests := map[string]struct {
		ticker  string
		hardCap uint32

		expectedErr error
	}{
		"success": {
			ticker:      "TEST-USD",
			hardCap:     300,
			expectedErr: nil,
		},
		"failure - hard cap set to 0": {
			ticker:      "TEST-USD",
			hardCap:     0,
			expectedErr: types.ErrMarketsHardCapReached,
		},
	}

	for name, tc := range tests {
		t.Run(
			name, func(t *testing.T) {
				tApp := testapp.NewTestAppBuilder(t).Build()
				ctx := tApp.InitChain()
				k := tApp.App.ListingKeeper
				ms := keeper.NewMsgServerImpl(k)

				err := k.SetMarketsHardCap(ctx, tc.hardCap)
				require.NoError(t, err)

				msg := types.MsgCreateMarketPermissionless{
					Ticker: tc.ticker,
					SubaccountId: &subaccounttypes.SubaccountId{
						Owner:  constants.AliceAccAddress.String(),
						Number: 0,
					},
					QuoteQuantums: dtypes.SerializableInt{},
				}

				_, err = ms.CreateMarketPermissionless(ctx, &msg)
				if tc.expectedErr != nil {
					require.ErrorContains(t, err, tc.expectedErr.Error())
				} else {
					require.NoError(t, err)
				}
			},
		)
	}
}
