package keeper_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"

	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/x/revshare/types"
	"github.com/stretchr/testify/require"
)

func TestGetSetMarketMapperRevShareDetails(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()

	k := tApp.App.RevShareKeeper
	// Set the rev share details for a market
	marketId := uint32(42)
	setDetails := types.MarketMapperRevShareDetails{
		ExpirationTs: 1735707600,
	}
	k.SetMarketMapperRevShareDetails(ctx, marketId, setDetails)

	// Get the rev share details for the market
	getDetails, err := k.GetMarketMapperRevShareDetails(ctx, marketId)
	require.NoError(t, err)
	require.Equal(t, getDetails.ExpirationTs, setDetails.ExpirationTs)

	// Set expiration ts to 0
	setDetails.ExpirationTs = 0
	k.SetMarketMapperRevShareDetails(ctx, marketId, setDetails)

	getDetails, err = k.GetMarketMapperRevShareDetails(ctx, marketId)
	require.NoError(t, err)
	require.Equal(t, getDetails.ExpirationTs, setDetails.ExpirationTs)
}

func TestGetMarketMapperRevShareDetailsFailure(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.RevShareKeeper

	// Get the rev share details for non-existent market
	_, err := k.GetMarketMapperRevShareDetails(ctx, 42)
	require.ErrorContains(t, err, "MarketMapperRevShareDetails not found for marketId: 42")
}

func TestCreateNewMarketRevShare(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()

	k := tApp.App.RevShareKeeper

	// Set base rev share params
	err := k.SetMarketMapperRevenueShareParams(
		ctx, types.MarketMapperRevenueShareParams{
			Address:         constants.AliceAccAddress.String(),
			RevenueSharePpm: 100_000, // 10%
			ValidDays:       240,
		},
	)
	require.NoError(t, err)

	// Create a new market rev share
	marketId := uint32(42)
	k.CreateNewMarketRevShare(ctx, marketId)
	require.NoError(t, err)

	// Check if the market rev share exists
	details, err := k.GetMarketMapperRevShareDetails(ctx, marketId)
	require.NoError(t, err)

	// TODO: is this blocktime call deterministic?
	expectedExpirationTs := ctx.BlockTime().Unix() + 240*24*60*60
	require.Equal(t, details.ExpirationTs, uint64(expectedExpirationTs))
}

func TestGetMarketMapperRevenueShareForMarket(t *testing.T) {
	tests := map[string]struct {
		revShareParams     types.MarketMapperRevenueShareParams
		marketId           uint32
		expirationDelta    int64
		setRevShareDetails bool

		// expected
		expectedMarketMapperAddr sdk.AccAddress
		expectedRevenueSharePpm  uint32
		expectedErr              error
	}{
		"valid market": {
			revShareParams: types.MarketMapperRevenueShareParams{
				Address:         constants.AliceAccAddress.String(),
				RevenueSharePpm: 100_000, // 10%
				ValidDays:       0,
			},
			marketId:           42,
			expirationDelta:    10,
			setRevShareDetails: true,

			expectedMarketMapperAddr: constants.AliceAccAddress,
			expectedRevenueSharePpm:  100_000,
		},
		"invalid market": {
			revShareParams: types.MarketMapperRevenueShareParams{
				Address:         constants.AliceAccAddress.String(),
				RevenueSharePpm: 100_000, // 10%
				ValidDays:       0,
			},
			marketId:           42,
			setRevShareDetails: false,

			expectedErr: types.ErrMarketMapperRevShareDetailsNotFound,
		},
		// TODO: investigate why tApp blocktime doesn't translate to ctx.BlockTime()
		//"expired market rev share": {
		//	revShareParams: types.MarketMapperRevenueShareParams{
		//		Address:         constants.AliceAccAddress.String(),
		//		RevenueSharePpm: 100_000, // 10%
		//		ValidDays:       0,
		//	},
		//	marketId:           42,
		//	expirationDelta:    -10,
		//	setRevShareDetails: true,
		//
		//	expectedMarketMapperAddr: nil,
		//	expectedRevenueSharePpm:  0,
		//},
	}

	for name, tc := range tests {
		t.Run(
			name, func(t *testing.T) {
				tApp := testapp.NewTestAppBuilder(t).Build()
				ctx := tApp.InitChain()
				k := tApp.App.RevShareKeeper
				tApp.AdvanceToBlock(
					2, testapp.AdvanceToBlockOptions{BlockTime: time.Now()},
				)

				// Set base rev share params
				err := k.SetMarketMapperRevenueShareParams(ctx, tc.revShareParams)
				require.NoError(t, err)

				// Set market rev share details
				if tc.setRevShareDetails {
					k.SetMarketMapperRevShareDetails(
						ctx, tc.marketId, types.MarketMapperRevShareDetails{
							ExpirationTs: uint64(ctx.BlockTime().Unix() + tc.expirationDelta),
						},
					)
				}

				// Get the revenue share for the market
				marketMapperAddr, revenueSharePpm, err := k.GetMarketMapperRevenueShareForMarket(ctx, tc.marketId)
				if tc.expectedErr != nil {
					require.ErrorIs(t, err, tc.expectedErr)
				} else {
					require.NoError(t, err)
					require.Equal(t, tc.expectedMarketMapperAddr, marketMapperAddr)
					require.Equal(t, tc.expectedRevenueSharePpm, revenueSharePpm)
				}
			},
		)
	}
}
