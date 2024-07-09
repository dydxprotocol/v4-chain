package keeper_test

import (
	"math/big"
	"testing"

	testapp "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/app"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/ratelimit/keeper"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/ratelimit/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/sDAIOracle/api"
)

func TestMsgUpdateSDAIConversionRate(t *testing.T) {

	testCases := []struct {
		name             string
		input            *types.MsgUpdateSDAIConversionRate
		expectedSDAIRate string
		expErr           bool
	}{
		{
			name: "Valid input",
			input: &types.MsgUpdateSDAIConversionRate{
				Sender:              "cosmos139f7kncmglres2nf3h4hc4tade85ekfr8sulz5",
				ConversionRate:      "1",
				EthereumBlockNumber: "1",
			},
			expectedSDAIRate: "1",
			expErr:           false,
		},
		{
			name: "Invalid conversion rate (empty)",
			input: &types.MsgUpdateSDAIConversionRate{
				Sender:              "cosmos139f7kncmglres2nf3h4hc4tade85ekfr8sulz5",
				ConversionRate:      "",
				EthereumBlockNumber: "1",
			},
			expectedSDAIRate: "",
			expErr:           true,
		},
		{
			name: "Invalid ethereum block number (empty)",
			input: &types.MsgUpdateSDAIConversionRate{
				Sender:              "cosmos139f7kncmglres2nf3h4hc4tade85ekfr8sulz5",
				ConversionRate:      "1",
				EthereumBlockNumber: "",
			},
			expectedSDAIRate: "",
			expErr:           true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).Build()
			ctx := tApp.InitChain()
			k := tApp.App.RatelimitKeeper
			ms := keeper.NewMsgServerImpl(k)

			k.SetSDAIPrice(ctx, big.NewInt(0))

			sDAIEventManager := k.GetSDAIEventManagerForTestingOnly()

			sDAIEventManager.AddsDAIEvent(&api.AddsDAIEventsRequest{
				ConversionRate:      "1",
				EthereumBlockNumber: "1",
			})

			_, err := ms.UpdateSDAIConversionRate(ctx, tc.input)
			if tc.expErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				sdkCtx := sdk.UnwrapSDKContext(ctx)
				price, found := k.GetSDAIPrice(sdkCtx)
				require.True(t, found)
				require.Equal(t,
					tc.expectedSDAIRate,
					price.String(),
				)
			}
		})
	}
}
