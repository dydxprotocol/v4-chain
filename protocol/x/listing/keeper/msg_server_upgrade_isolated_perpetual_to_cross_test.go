package keeper_test

import (
	"math/big"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	indexerevents "github.com/dydxprotocol/v4-chain/protocol/indexer/events"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"
	v1 "github.com/dydxprotocol/v4-chain/protocol/indexer/protocol/v1"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	bank_testutil "github.com/dydxprotocol/v4-chain/protocol/testutil/bank"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	asstypes "github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/listing/keeper"
	listingkeeper "github.com/dydxprotocol/v4-chain/protocol/x/listing/keeper"
	types "github.com/dydxprotocol/v4-chain/protocol/x/listing/types"
	perpetualtypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

var (
	validAuthority = lib.GovModuleAddress.String()
)

func TestMsgUpgradeIsolatedPerpetualToCross(t *testing.T) {
	tests := map[string]struct {
		msg                           *types.MsgUpgradeIsolatedPerpetualToCross
		isolatedInsuranceFundBalance  *big.Int
		isolatedCollateralPoolBalance *big.Int
		crossInsuranceFundBalance     *big.Int
		crossCollateralPoolBalance    *big.Int

		expectedErr string
	}{
		"Success": {
			msg: &types.MsgUpgradeIsolatedPerpetualToCross{
				Authority:   validAuthority,
				PerpetualId: 3, // isolated
			},
			isolatedInsuranceFundBalance:  big.NewInt(1),
			isolatedCollateralPoolBalance: big.NewInt(1),
			crossInsuranceFundBalance:     big.NewInt(1),
			crossCollateralPoolBalance:    big.NewInt(1),
			expectedErr:                   "",
		},
		"Success - empty isolated insurance fund": {
			msg: &types.MsgUpgradeIsolatedPerpetualToCross{
				Authority:   validAuthority,
				PerpetualId: 3, // isolated
			},
			isolatedInsuranceFundBalance:  big.NewInt(0),
			isolatedCollateralPoolBalance: big.NewInt(1),
			crossInsuranceFundBalance:     big.NewInt(1),
			crossCollateralPoolBalance:    big.NewInt(1),
			expectedErr:                   "",
		},
		"Success - empty isolated collateral fund": {
			msg: &types.MsgUpgradeIsolatedPerpetualToCross{
				Authority:   validAuthority,
				PerpetualId: 3, // isolated
			},
			isolatedInsuranceFundBalance:  big.NewInt(1),
			isolatedCollateralPoolBalance: big.NewInt(0),
			crossInsuranceFundBalance:     big.NewInt(1),
			crossCollateralPoolBalance:    big.NewInt(1),
			expectedErr:                   "",
		},
		"Success - empty isolated insurance fund + empty isolated collateral fund": {
			msg: &types.MsgUpgradeIsolatedPerpetualToCross{
				Authority:   validAuthority,
				PerpetualId: 3, // isolated
			},
			isolatedInsuranceFundBalance:  big.NewInt(0),
			isolatedCollateralPoolBalance: big.NewInt(0),
			crossInsuranceFundBalance:     big.NewInt(1),
			crossCollateralPoolBalance:    big.NewInt(1),
			expectedErr:                   "",
		},
		"Failure: Empty authority": {
			msg: &types.MsgUpgradeIsolatedPerpetualToCross{
				Authority:   "",
				PerpetualId: 3, // isolated
			},
			isolatedInsuranceFundBalance:  big.NewInt(1),
			isolatedCollateralPoolBalance: big.NewInt(1),
			crossInsuranceFundBalance:     big.NewInt(1),
			crossCollateralPoolBalance:    big.NewInt(1),
			expectedErr:                   "invalid authority",
		},
		"Failure: Invalid authority": {
			msg: &types.MsgUpgradeIsolatedPerpetualToCross{
				Authority:   "invalid",
				PerpetualId: 3, // isolated
			},
			isolatedInsuranceFundBalance:  big.NewInt(1),
			isolatedCollateralPoolBalance: big.NewInt(1),
			crossInsuranceFundBalance:     big.NewInt(1),
			crossCollateralPoolBalance:    big.NewInt(1),
			expectedErr:                   "invalid authority",
		},
		"Failure: Invalid perpetual ID": {
			msg: &types.MsgUpgradeIsolatedPerpetualToCross{
				Authority:   validAuthority,
				PerpetualId: 99999999, // invalid
			},
			expectedErr: "Perpetual does not exist",
		},
		"Failure - perpetual already has cross market type": {
			msg: &types.MsgUpgradeIsolatedPerpetualToCross{
				Authority:   validAuthority,
				PerpetualId: 1, // cross
			},
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

				var isolatedInsuranceFundAddr, crossInsuranceFundAddr, isolatedCollateralPoolAddr,
					crossCollateralPoolAddr sdk.AccAddress
				if tc.isolatedInsuranceFundBalance != nil {
					// Get addresses for isolated/cross insurance funds and collateral pools.
					isolatedInsuranceFundAddr, err = perpetualsKeeper.GetInsuranceFundModuleAddress(ctx, tc.msg.PerpetualId)
					require.NoError(t, err)

					isolatedCollateralPoolAddr, err = subaccountsKeeper.GetCollateralPoolFromPerpetualId(ctx, tc.msg.PerpetualId)
					require.NoError(t, err)

					crossInsuranceFundAddr = perpetualtypes.InsuranceFundModuleAddress

					crossCollateralPoolAddr = satypes.ModuleAddress

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
				}

				perpetual, err := perpetualsKeeper.GetPerpetual(ctx, tc.msg.PerpetualId)
				require.NoError(t, err)
				expectedIndexerEvent := &indexerevents.UpdatePerpetualEventV2{
					Id:               perpetual.Params.Id,
					Ticker:           perpetual.Params.Ticker,
					MarketId:         perpetual.Params.MarketId,
					AtomicResolution: perpetual.Params.AtomicResolution,
					LiquidityTier:    perpetual.Params.LiquidityTier,
					MarketType:       v1.ConvertToPerpetualMarketType(perpetualtypes.PerpetualMarketType_PERPETUAL_MARKET_TYPE_CROSS),
				}

				// Upgrade perpetual from isolated to cross.
				ms := listingkeeper.NewMsgServerImpl(*keeper)
				_, err = ms.UpgradeIsolatedPerpetualToCross(ctx, tc.msg)
				if tc.expectedErr != "" {
					require.ErrorContains(t, err, tc.expectedErr)
					return
				}
				require.NoError(t, err)

				// Check perpetual market type has been upgraded to cross.
				perpetual, err = perpetualsKeeper.GetPerpetual(ctx, tc.msg.PerpetualId)
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

				// Verify that expected indexer event was emitted.
				emittedIndexerEvents := getUpdatePerpetualEventsFromIndexerBlock(ctx, keeper)
				require.Equal(t, emittedIndexerEvents, []*indexerevents.UpdatePerpetualEventV2{expectedIndexerEvent})
			},
		)
	}
}

// getUpdatePerpetualEventsFromIndexerBlock returns all UpdatePerpetual events from the indexer block.
func getUpdatePerpetualEventsFromIndexerBlock(
	ctx sdk.Context,
	listingKeeper *keeper.Keeper,
) []*indexerevents.UpdatePerpetualEventV2 {
	block := listingKeeper.GetIndexerEventManager().ProduceBlock(ctx)
	var updatePerpetualEvents []*indexerevents.UpdatePerpetualEventV2
	for _, event := range block.Events {
		if event.Subtype != indexerevents.SubtypeUpdatePerpetual {
			continue
		}
		if _, ok := event.OrderingWithinBlock.(*indexer_manager.IndexerTendermintEvent_TransactionIndex); ok {
			var updatePerpetualEvent indexerevents.UpdatePerpetualEventV2
			err := proto.Unmarshal(event.DataBytes, &updatePerpetualEvent)
			if err != nil {
				panic(err)
			}
			updatePerpetualEvents = append(updatePerpetualEvents, &updatePerpetualEvent)
		}
	}
	return updatePerpetualEvents
}
