package keeper_test

import (
	"math/big"

	comettypes "github.com/cometbft/cometbft/types"
	oracletypes "github.com/dydxprotocol/slinky/pkg/types"
	marketmaptypes "github.com/dydxprotocol/slinky/x/marketmap/types"
	"github.com/dydxprotocol/slinky/x/marketmap/types/tickermetadata"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	testutil "github.com/dydxprotocol/v4-chain/protocol/testutil/util"
	"github.com/dydxprotocol/v4-chain/protocol/x/listing/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/listing/types"
	perpetualstypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	vaulttypes "github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
	"github.com/stretchr/testify/require"

	"testing"
)

func TestMsgCreateMarketPermissionless(t *testing.T) {
	tests := map[string]struct {
		ticker  string
		hardCap uint32
		balance *big.Int

		expectedErr error
	}{
		"success": {
			ticker:      "TEST2-USD",
			hardCap:     300,
			balance:     big.NewInt(10_000_000_000),
			expectedErr: nil,
		},
		"failure - hard cap reached": {
			ticker:  "TEST2-USD",
			hardCap: 0,
			balance: big.NewInt(10_000_000_000),

			expectedErr: types.ErrMarketsHardCapReached,
		},
		"failure - ticker not found": {
			ticker:  "INVALID-USD",
			hardCap: 300,
			balance: big.NewInt(10_000_000_000),

			expectedErr: types.ErrMarketNotFound,
		},
		"failure - market already listed": {
			ticker:  "BTC-USD",
			hardCap: 300,
			balance: big.NewInt(10_000_000_000),

			expectedErr: pricestypes.ErrMarketParamPairAlreadyExists,
		},
	}

	for name, tc := range tests {
		t.Run(
			name, func(t *testing.T) {
				tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(
					func() (genesis comettypes.GenesisDoc) {
						genesis = testapp.DefaultGenesis()
						// Initialize vault with its existing equity.
						testapp.UpdateGenesisDocWithAppStateForModule(
							&genesis,
							func(genesisState *satypes.GenesisState) {
								genesisState.Subaccounts = []satypes.Subaccount{
									{
										Id: &vaulttypes.MegavaultMainSubaccount,
										AssetPositions: []*satypes.AssetPosition{
											testutil.CreateSingleAssetPosition(
												0,
												big.NewInt(1_000_000),
											),
										},
									},
									{
										Id: &constants.Alice_Num0,
										AssetPositions: []*satypes.AssetPosition{
											testutil.CreateSingleAssetPosition(
												0,
												tc.balance,
											),
										},
									},
								}
							},
						)

						testapp.UpdateGenesisDocWithAppStateForModule(
							&genesis,
							func(genesisState *perpetualstypes.GenesisState) {
								genesisState.LiquidityTiers = append(genesisState.LiquidityTiers, perpetualstypes.LiquidityTier{
									Id:                     7,
									Name:                   "IML 5x",
									InitialMarginPpm:       200000,
									MaintenanceFractionPpm: 500000,
									ImpactNotional:         1000000000,
									OpenInterestLowerCap:   500000000000,
									OpenInterestUpperCap:   1000000000000,
								})
							},
						)
						return genesis
					},
				).Build()
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
