package keeper_test

import (
	"testing"
	"time"

	testapp "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/app"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/ratelimit/keeper"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/ratelimit/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	store "github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/sDAIOracle/client/contract"
	oracletypes "github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/sDAIOracle/client/types"
)

func TestMsgUpdateSDAIConversionRate(t *testing.T) {

	// Test with real client
	client, err := ethclient.Dial(oracletypes.ETHRPC)
	if err != nil {
		t.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	chi, blockNumber, err := store.QueryDaiConversionRate(client)
	assert.Nil(t, err, "Expected no error with real client")

	time.Sleep(10 * time.Second) // to ensure other validators have queried the sdai rate at this block

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
				ConversionRate:      chi,
				EthereumBlockNumber: blockNumber,
			},
			expectedSDAIRate: chi,
			expErr:           false,
		},
		{
			name: "Invalid address (empty)",
			input: &types.MsgUpdateSDAIConversionRate{
				Sender:              "",
				ConversionRate:      "1",
				EthereumBlockNumber: "1",
			},
			expectedSDAIRate: "",
			expErr:           true,
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
