package keeper_test

import (
	"errors"
	"math/big"
	"testing"

	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	bank_testutil "github.com/dydxprotocol/v4-chain/protocol/testutil/bank"
	testutil "github.com/dydxprotocol/v4-chain/protocol/testutil/util"

	vaulttypes "github.com/dydxprotocol/v4-chain/protocol/x/vault/types"

	asstypes "github.com/dydxprotocol/v4-chain/protocol/x/assets/types"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"

	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"

	"github.com/stretchr/testify/mock"

	"github.com/dydxprotocol/v4-chain/protocol/lib"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"

	perpetualtypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	oracletypes "github.com/skip-mev/connect/v2/pkg/types"
	marketmaptypes "github.com/skip-mev/connect/v2/x/marketmap/types"
	"github.com/skip-mev/connect/v2/x/marketmap/types/tickermetadata"

	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"

	comettypes "github.com/cometbft/cometbft/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/listing/types"

	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	"github.com/stretchr/testify/require"
)

func TestCreateMarket(t *testing.T) {
	tests := map[string]struct {
		ticker          string
		duplicateMarket bool

		expectedErr error
	}{
		"success": {
			ticker:      "TEST-USD",
			expectedErr: nil,
		},
		"failure - market not found": {
			ticker:      "INVALID-USD",
			expectedErr: types.ErrMarketNotFound,
		},
		"failure - duplicate market": {
			ticker:      "BTC-USD",
			expectedErr: pricestypes.ErrMarketParamPairAlreadyExists,
		},
	}

	for name, tc := range tests {
		t.Run(
			name, func(t *testing.T) {
				mockIndexerEventManager := &mocks.IndexerEventManager{}
				ctx, keeper, _, _, pricesKeeper, _, _, marketMapKeeper, _, _, _ := keepertest.ListingKeepers(
					t,
					&mocks.BankKeeper{},
					mockIndexerEventManager,
				)

				keepertest.CreateTestMarkets(t, ctx, pricesKeeper)

				testMarketParams := pricestypes.MarketParam{
					Pair:               "TEST-USD",
					Exponent:           int32(-6),
					ExchangeConfigJson: `{"test_config_placeholder":{}}`,
					MinExchanges:       2,
					MinPriceChangePpm:  uint32(800),
				}

				keepertest.CreateMarketsInMarketMapFromParams(
					t,
					ctx,
					marketMapKeeper,
					[]pricestypes.MarketParam{
						testMarketParams,
					},
				)

				marketId, err := keeper.CreateMarket(ctx, tc.ticker)
				if tc.expectedErr != nil {
					require.ErrorContains(t, err, tc.expectedErr.Error())
				} else {
					require.NoError(t, err)

					// Check if the market was created
					market, exists := pricesKeeper.GetMarketParam(ctx, marketId)
					require.True(t, exists)
					require.Equal(t, testMarketParams.Pair, market.Pair)
					require.Equal(t, testMarketParams.MinPriceChangePpm, types.MinPriceChangePpm_LongTail)
				}
			},
		)
	}
}

func TestCreatePerpetual(t *testing.T) {
	tests := map[string]struct {
		ticker         string
		referencePrice uint64

		expectedErr error
	}{
		"success": {
			ticker:         "TEST-USD",
			referencePrice: 1000000000, // $1000
			expectedErr:    nil,
		},
		"failure - reference price 0": {
			ticker:         "TEST-USD",
			referencePrice: 0,
			expectedErr:    types.ErrReferencePriceZero,
		},
		"failure - market not found": {
			ticker:      "INVALID-USD",
			expectedErr: types.ErrMarketNotFound,
		},
	}

	for name, tc := range tests {
		t.Run(
			name, func(t *testing.T) {
				mockIndexerEventManager := &mocks.IndexerEventManager{}
				ctx, keeper, _, _, pricesKeeper, perpetualsKeeper, _, marketMapKeeper, _, _, _ := keepertest.ListingKeepers(
					t,
					&mocks.BankKeeper{},
					mockIndexerEventManager,
				)
				keepertest.CreateLiquidityTiersAndNPerpetuals(t, ctx, perpetualsKeeper, pricesKeeper, 10)

				// Create a marketmap with a single market
				dydxMetadata, err := tickermetadata.MarshalDyDx(
					tickermetadata.DyDx{
						ReferencePrice: tc.referencePrice,
						Liquidity:      0,
						AggregateIDs:   nil,
					},
				)
				require.NoError(t, err)

				market := marketmaptypes.Market{
					Ticker: marketmaptypes.Ticker{
						CurrencyPair:     oracletypes.CurrencyPair{Base: "TEST", Quote: "USD"},
						Decimals:         10,
						MinProviderCount: 2,
						Enabled:          false,
						Metadata_JSON:    string(dydxMetadata),
					},
					ProviderConfigs: []marketmaptypes.ProviderConfig{
						{
							Name:           "binance_ws",
							OffChainTicker: "TESTUSDT",
						},
					},
				}
				err = marketMapKeeper.CreateMarket(ctx, market)
				require.NoError(t, err)

				marketId, err := keeper.CreateMarket(ctx, tc.ticker)
				if errors.Is(tc.expectedErr, types.ErrMarketNotFound) {
					require.ErrorContains(t, err, tc.expectedErr.Error())
					return
				}

				perpetualId, err := keeper.CreatePerpetual(ctx, marketId, tc.ticker)
				if tc.expectedErr != nil {
					require.Error(t, err)
				} else {
					require.NoError(t, err)

					// Check if the perpetual was created
					perpetual, err := perpetualsKeeper.GetPerpetual(ctx, perpetualId)
					require.NoError(t, err)
					require.Equal(t, uint32(10), perpetual.GetId())
					require.Equal(t, marketId, perpetual.Params.MarketId)
					require.Equal(t, tc.ticker, perpetual.Params.Ticker)
					// Expected resolution = -6 - (Floor(log10(1000000000))-10) = -5
					require.Equal(t, int32(-5), perpetual.Params.AtomicResolution)
					require.Equal(t, int32(types.DefaultFundingPpm), perpetual.Params.DefaultFundingPpm)
					require.Equal(t, uint32(types.LiquidityTier_Isolated), perpetual.Params.LiquidityTier)
					require.Equal(
						t, perpetualtypes.PerpetualMarketType_PERPETUAL_MARKET_TYPE_ISOLATED,
						perpetual.Params.MarketType,
					)
				}
			},
		)
	}
}

func TestCreateClobPair(t *testing.T) {
	tests := map[string]struct {
		ticker      string
		isDeliverTx bool
	}{
		"deliverTx - true": {
			ticker:      "TEST-USD",
			isDeliverTx: true,
		},
		"deliverTx - false": {
			ticker:      "TEST-USD",
			isDeliverTx: false,
		},
	}

	for name, tc := range tests {
		t.Run(
			name, func(t *testing.T) {
				mockIndexerEventManager := &mocks.IndexerEventManager{}
				ctx, keeper, _, _, pricesKeeper, perpetualsKeeper, clobKeeper, marketMapKeeper,
					_, _, _ := keepertest.ListingKeepers(
					t,
					&mocks.BankKeeper{},
					mockIndexerEventManager,
				)
				mockIndexerEventManager.On(
					"AddTxnEvent",
					mock.Anything,
					mock.Anything,
					mock.Anything,
					mock.Anything,
				).Return()
				keepertest.CreateLiquidityTiersAndNPerpetuals(t, ctx, perpetualsKeeper, pricesKeeper, 10)

				// Set deliverTx mode
				if tc.isDeliverTx {
					ctx = ctx.WithIsCheckTx(false).WithIsReCheckTx(false).WithExecMode(sdk.ExecModeFinalize)
					lib.AssertDeliverTxMode(ctx)
				} else {
					ctx = ctx.WithIsCheckTx(true)
					lib.AssertCheckTxMode(ctx)
				}

				// Create a marketmap with a single market
				dydxMetadata, err := tickermetadata.MarshalDyDx(
					tickermetadata.DyDx{
						ReferencePrice: 1000000000,
						Liquidity:      0,
						AggregateIDs:   nil,
					},
				)
				require.NoError(t, err)

				market := marketmaptypes.Market{
					Ticker: marketmaptypes.Ticker{
						CurrencyPair:     oracletypes.CurrencyPair{Base: "TEST", Quote: "USD"},
						Decimals:         6,
						MinProviderCount: 2,
						Enabled:          false,
						Metadata_JSON:    string(dydxMetadata),
					},
					ProviderConfigs: []marketmaptypes.ProviderConfig{
						{
							Name:           "binance_ws",
							OffChainTicker: "TESTUSDT",
						},
					},
				}
				err = marketMapKeeper.CreateMarket(ctx, market)
				require.NoError(t, err)

				marketId, err := keeper.CreateMarket(ctx, tc.ticker)
				require.NoError(t, err)

				perpetualId, err := keeper.CreatePerpetual(ctx, marketId, tc.ticker)
				require.NoError(t, err)

				clobPairId, err := keeper.CreateClobPair(ctx, perpetualId)
				require.NoError(t, err)

				clobPair, found := clobKeeper.GetClobPair(ctx, clobtypes.ClobPairId(clobPairId))
				require.True(t, found)
				require.Equal(t, clobtypes.ClobPair_STATUS_ACTIVE, clobPair.Status)
				require.Equal(
					t,
					clobtypes.SubticksPerTick(types.SubticksPerTick_LongTail),
					clobPair.GetClobPairSubticksPerTick(),
				)
				require.Equal(
					t,
					types.DefaultStepBaseQuantums,
					clobPair.GetClobPairMinOrderBaseQuantums().ToUint64(),
				)
				require.Equal(t, perpetualId, clobPair.MustGetPerpetualId())

				// Should not modify in-memory object right away
				_, found = clobKeeper.PerpetualIdToClobPairId[perpetualId]
				require.False(t, found)

				// Check the corresponding ClobPair creation was staged.
				stagedEvents := clobKeeper.GetStagedClobFinalizeBlockEvents(ctx)

				if tc.isDeliverTx {
					require.Equal(t, 1, len(stagedEvents))
					require.Equal(t,
						stagedEvents[0].GetCreateClobPair().GetPerpetualClobMetadata().PerpetualId,
						perpetualId,
					)
				} else {
					require.Equal(t, 0, len(stagedEvents))
				}
			},
		)
	}
}

func TestDepositToMegavaultforPML(t *testing.T) {
	tests := map[string]struct {
		address    string
		balance    *big.Int
		asset      asstypes.Asset
		clobPairId uint32

		expectedErr string
	}{
		"success": {
			address:    constants.AliceAccAddress.String(),
			balance:    big.NewInt(10_000_000_000), // 10k USDC
			asset:      *constants.Usdc,
			clobPairId: 1,

			expectedErr: "",
		},
		"failure - insufficient balance": {
			address:    constants.AliceAccAddress.String(),
			balance:    big.NewInt(0),
			asset:      *constants.Usdc,
			clobPairId: 1,

			expectedErr: "NewlyUndercollateralized",
		},
		"failure - invalid clob pair id": {
			address:    constants.AliceAccAddress.String(),
			balance:    big.NewInt(10_000_000_000), // 10k USDC
			asset:      *constants.Usdc,
			clobPairId: 100, // non existent clob pair id

			expectedErr: "ClobPair not found",
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
										Id: &satypes.SubaccountId{
											Owner:  tc.address,
											Number: 0,
										},
										AssetPositions: []*satypes.AssetPosition{
											testutil.CreateSingleAssetPosition(
												tc.asset.Id,
												tc.balance,
											),
										},
									},
								}
							},
						)
						return genesis
					},
				).Build()

				ctx := tApp.InitChain()

				// Set existing total shares.
				err := tApp.App.VaultKeeper.SetTotalShares(
					ctx,
					vaulttypes.BigIntToNumShares(big.NewInt(1_000_000)),
				)
				require.NoError(t, err)

				err = tApp.App.ListingKeeper.DepositToMegavaultforPML(
					ctx,
					satypes.SubaccountId{
						Owner:  tc.address,
						Number: 0,
					},
					tc.clobPairId,
				)
				if tc.expectedErr != "" {
					require.ErrorContains(t, err, tc.expectedErr)
				} else {
					require.NoError(t, err)
					vaultParams, exists := tApp.App.VaultKeeper.GetVaultParams(
						ctx,
						vaulttypes.VaultId{
							Type:   vaulttypes.VaultType_VAULT_TYPE_CLOB,
							Number: tc.clobPairId,
						},
					)
					require.True(t, exists)
					require.Equal(t, vaulttypes.VaultStatus_VAULT_STATUS_QUOTING, vaultParams.Status)
				}
			},
		)
	}
}

func TestUpgradeIsolatedPerpetualToCross(t *testing.T) {
	tests := map[string]struct {
		perpetualId                   uint32
		isolatedInsuranceFundBalance  *big.Int
		isolatedCollateralPoolBalance *big.Int
		crossInsuranceFundBalance     *big.Int
		crossCollateralPoolBalance    *big.Int

		expectedErr string
	}{
		"success": {
			perpetualId:                   3, // isolated
			isolatedInsuranceFundBalance:  big.NewInt(1),
			isolatedCollateralPoolBalance: big.NewInt(1),
			crossInsuranceFundBalance:     big.NewInt(1),
			crossCollateralPoolBalance:    big.NewInt(1),
			expectedErr:                   "",
		},
		"success - empty isolated insurance fund": {
			perpetualId:                   3, // isolated
			isolatedInsuranceFundBalance:  big.NewInt(0),
			isolatedCollateralPoolBalance: big.NewInt(1),
			crossInsuranceFundBalance:     big.NewInt(1),
			crossCollateralPoolBalance:    big.NewInt(1),
			expectedErr:                   "",
		},
		"success - empty isolated collateral fund": {
			perpetualId:                   3, // isolated
			isolatedInsuranceFundBalance:  big.NewInt(1),
			isolatedCollateralPoolBalance: big.NewInt(0),
			crossInsuranceFundBalance:     big.NewInt(1),
			crossCollateralPoolBalance:    big.NewInt(1),
			expectedErr:                   "",
		},
		"success - empty isolated insurance fund + empty isolated collateral fund": {
			perpetualId:                   3, // isolated
			isolatedInsuranceFundBalance:  big.NewInt(0),
			isolatedCollateralPoolBalance: big.NewInt(0),
			crossInsuranceFundBalance:     big.NewInt(1),
			crossCollateralPoolBalance:    big.NewInt(1),
			expectedErr:                   "",
		},
		"failure - perpetual already has cross market type": {
			perpetualId:                   1, // cross
			isolatedInsuranceFundBalance:  big.NewInt(1),
			isolatedCollateralPoolBalance: big.NewInt(1),
			crossInsuranceFundBalance:     big.NewInt(1),
			crossCollateralPoolBalance:    big.NewInt(1),
			expectedErr:                   "perpetual 1 is not an isolated perpetual and cannot be upgraded to cross",
		},
	}

	for name, tc := range tests {
		t.Run(
			name, func(t *testing.T) {
				mockIndexerEventManager := &mocks.IndexerEventManager{}

				ctx, keeper, _, _, pricesKeeper, perpetualsKeeper, _, _, assetsKeeper,
					bankKeeper, subaccountsKeeper := keepertest.ListingKeepers(
					t,
					&mocks.BankKeeper{},
					mockIndexerEventManager,
				)

				// Create the default markets.
				keepertest.CreateTestMarkets(t, ctx, pricesKeeper)

				// Create liquidity tiers.
				keepertest.CreateTestLiquidityTiers(t, ctx, perpetualsKeeper)

				// Create USDC asset.
				err := keepertest.CreateUsdcAsset(ctx, assetsKeeper)
				require.NoError(t, err)

				// Create test perpetuals.
				keepertest.CreateTestPerpetuals(t, ctx, perpetualsKeeper)

				// Get addresses for isolated/cross insurance funds and collateral pools.
				isolatedInsuranceFundAddr, err := perpetualsKeeper.GetInsuranceFundModuleAddress(ctx, tc.perpetualId)
				require.NoError(t, err)

				isolatedCollateralPoolAddr, err := subaccountsKeeper.GetCollateralPoolFromPerpetualId(ctx, tc.perpetualId)
				require.NoError(t, err)

				crossInsuranceFundAddr := perpetualtypes.InsuranceFundModuleAddress

				crossCollateralPoolAddr := satypes.ModuleAddress

				// Fund the isolated insurance account, cross insurance account,
				// isolated collateral pool, and cross collateral pool.
				fundingData := [][]interface{}{
					{isolatedInsuranceFundAddr, tc.isolatedInsuranceFundBalance},
					{crossInsuranceFundAddr, tc.crossInsuranceFundBalance},
					{isolatedCollateralPoolAddr, tc.isolatedCollateralPoolBalance},
					{crossCollateralPoolAddr, tc.crossCollateralPoolBalance},
				}

				for _, data := range fundingData {
					addr := data[0].(sdk.AccAddress)
					amount := data[1].(*big.Int)

					if amount.Cmp(big.NewInt(0)) != 0 {
						err = bank_testutil.FundAccount(
							ctx,
							addr,
							sdk.Coins{
								sdk.NewCoin(constants.Usdc.Denom, sdkmath.NewIntFromBigInt(amount)),
							},
							*bankKeeper,
						)
						require.NoError(t, err)
					}
				}

				// Upgrade perpetual from isolated to cross.
				err = keeper.UpgradeIsolatedPerpetualToCross(ctx, tc.perpetualId)
				if tc.expectedErr != "" {
					require.ErrorContains(t, err, tc.expectedErr)
					return
				}
				require.NoError(t, err)

				// Check perpetual market type has been upgraded to cross.
				perpetual, err := perpetualsKeeper.GetPerpetual(ctx, tc.perpetualId)
				require.NoError(t, err)
				require.Equal(
					t,
					perpetualtypes.PerpetualMarketType_PERPETUAL_MARKET_TYPE_CROSS,
					perpetual.Params.MarketType,
				)

				// Check expected balance for isolated/cross insurance funds and collateral pools.
				expectedBalances := [][]interface{}{
					{isolatedInsuranceFundAddr, big.NewInt(0)},
					{crossInsuranceFundAddr, big.NewInt(0).Add(tc.isolatedInsuranceFundBalance, tc.crossInsuranceFundBalance)},
					{isolatedCollateralPoolAddr, big.NewInt(0)},
					{crossCollateralPoolAddr, big.NewInt(0).Add(tc.isolatedCollateralPoolBalance, tc.crossCollateralPoolBalance)},
				}

				for _, data := range expectedBalances {
					addr := data[0].(sdk.AccAddress)
					amount := data[1].(*big.Int)

					require.Equal(
						t,
						sdk.NewCoin(
							asstypes.AssetUsdc.Denom,
							sdkmath.NewIntFromBigInt(amount),
						),
						bankKeeper.GetBalance(ctx, addr, asstypes.AssetUsdc.Denom),
					)
				}
			},
		)
	}
}
