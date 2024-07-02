package keeper_test

import (
	"testing"

	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/x/revshare/keeper"
	"github.com/stretchr/testify/require"

	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/x/revshare/types"
)

func TestSetMarketMapperRevenueShareDetailsForMarket(t *testing.T) {
	tests := map[string]struct {
		// Msg
		msg *types.MsgSetMarketMapperRevShareDetailsForMarket
		// Expected error
		expectedErr string
	}{
		"Success - Set revenue share details for market": {
			msg: &types.MsgSetMarketMapperRevShareDetailsForMarket{
				Authority: lib.GovModuleAddress.String(),
				MarketId:  constants.MarketId0,
				Params: types.MarketMapperRevShareDetails{
					ExpirationTs: 100,
				},
			},
			expectedErr: "",
		},
		"Failure - Invalid Authority": {
			msg: &types.MsgSetMarketMapperRevShareDetailsForMarket{
				Authority: constants.AliceAccAddress.String(),
				MarketId:  constants.MarketId0,
				Params: types.MarketMapperRevShareDetails{
					ExpirationTs: 100,
				},
			},
			expectedErr: "invalid authority",
		},
		"Failure - Empty Authority": {
			msg: &types.MsgSetMarketMapperRevShareDetailsForMarket{
				MarketId: constants.MarketId0,
				Params: types.MarketMapperRevShareDetails{
					ExpirationTs: 100,
				},
			},
			expectedErr: "invalid authority",
		},
	}

	for name, tc := range tests {
		t.Run(
			name, func(t *testing.T) {
				tApp := testapp.NewTestAppBuilder(t).Build()
				ctx := tApp.InitChain()
				k := tApp.App.RevShareKeeper
				ms := keeper.NewMsgServerImpl(k)
				_, err := ms.SetMarketMapperRevShareDetailsForMarket(ctx, tc.msg)
				if tc.expectedErr != "" {
					require.Error(t, err)
					require.Contains(t, err.Error(), tc.expectedErr)
				} else {
					require.NoError(t, err)
					revShareDetails, _ := k.GetMarketMapperRevShareDetails(ctx, tc.msg.MarketId)
					require.Equal(t, tc.msg.Params, revShareDetails)
				}
			},
		)
	}
}
