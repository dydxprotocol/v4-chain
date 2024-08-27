package keeper_test

import (
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/x/listing/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/listing/types"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	oracletypes "github.com/skip-mev/slinky/pkg/types"
	marketmaptypes "github.com/skip-mev/slinky/x/marketmap/types"
	"github.com/skip-mev/slinky/x/marketmap/types/tickermetadata"
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
			ticker:      "TEST2-USD",
			hardCap:     300,
			expectedErr: nil,
		},
		"failure - hard cap reached": {
			ticker:      "TEST2-USD",
			hardCap:     0,
			expectedErr: types.ErrMarketsHardCapReached,
		},
		"failure - ticker not found": {
			ticker:      "INVALID-USD",
			hardCap:     300,
			expectedErr: types.ErrMarketNotFound,
		},
		"failure - market already listed": {
			ticker:      "BTC-USD",
			hardCap:     300,
			expectedErr: pricestypes.ErrMarketParamPairAlreadyExists,
		},
	}

	for name, tc := range tests {
		t.Run(
			name, func(t *testing.T) {
				tApp := testapp.NewTestAppBuilder(t).Build()
				ctx := tApp.InitChain()
				k := tApp.App.ListingKeeper
				ms := keeper.NewMsgServerImpl(k)

				// Set hard cap
				err := k.SetMarketsHardCap(ctx, tc.hardCap)
				require.NoError(t, err)

				// Add TEST2-USD market to market map
				dydxMetadata, err := tickermetadata.MarshalDyDx(
					tickermetadata.DyDx{
						ReferencePrice: 10000000,
						Liquidity:      0,
						AggregateIDs:   nil,
					},
				)

				require.NoError(t, err)
				market := marketmaptypes.Market{
					Ticker: marketmaptypes.Ticker{
						CurrencyPair:     oracletypes.CurrencyPair{Base: "TEST2", Quote: "USD"},
						Decimals:         6,
						MinProviderCount: 2,
						Enabled:          false,
						Metadata_JSON:    string(dydxMetadata),
					},
					ProviderConfigs: []marketmaptypes.ProviderConfig{
						{
							Name:           "binance_ws",
							OffChainTicker: "TEST2USDT",
						},
					},
				}
				err = k.MarketMapKeeper.CreateMarket(ctx, market)
				require.NoError(t, err)

				msg := types.MsgCreateMarketPermissionless{
					Ticker: tc.ticker,
					SubaccountId: &satypes.SubaccountId{
						Owner:  constants.AliceAccAddress.String(),
						Number: 0,
					},
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
