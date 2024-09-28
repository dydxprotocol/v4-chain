package keeper_test

import (
	"math"
	"math/big"
	"strconv"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdaiservertypes "github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/server/types/sdaioracle"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/lib"
	ratelimitkeeper "github.com/StreamFinance-Protocol/stream-chain/protocol/x/ratelimit/keeper"
	ratelimittypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/ratelimit/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	bank_testutil "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/bank"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/dtypes"
	indexerevents "github.com/StreamFinance-Protocol/stream-chain/protocol/indexer/events"
	big_testutil "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/big"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/constants"
	testutil "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/keeper"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/nullify"
	perptest "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/perpetuals"
	pricestest "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/prices"
	asstypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/assets/types"
	perptypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/perpetuals/types"
	pricestypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices/types"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/subaccounts/keeper"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/subaccounts/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func createNSubaccount(keeper *keeper.Keeper, ctx sdk.Context, n int, tdaiBalance *big.Int) []types.Subaccount {
	items := make([]types.Subaccount, n)
	for i := range items {
		items[i].Id = &types.SubaccountId{
			Owner:  strconv.Itoa(i),
			Number: uint32(i),
		}
		items[i].AssetPositions = testutil.CreateTDaiAssetPosition(tdaiBalance)
		items[i].AssetYieldIndex = "1/1"

		keeper.SetSubaccount(ctx, items[i])
	}
	return items
}

// assertSubaccountUpdateEventsNotInIndexerBlock checks that no subaccount update events were
// included in the Indexer block kafka message
func assertSubaccountUpdateEventsNotInIndexerBlock(
	t *testing.T,
	k *keeper.Keeper,
	ctx sdk.Context,
) {
	subaccountUpdates := testutil.GetSubaccountUpdateEventsFromIndexerBlock(ctx, k)
	require.Empty(t, subaccountUpdates)
}

// assertSubaccountUpdateEventsInIndexerBlock checks that the correct subaccount update events were
// included in the Indexer block kafka message, given details of the updates applied,
// the expected return values of the update subaccount function.
func assertSubaccountUpdateEventsInIndexerBlock(
	t *testing.T,
	k *keeper.Keeper,
	ctx sdk.Context,
	expectedErr error,
	expectedSuccess bool,
	updates []types.Update,
	expectedSuccessPerUpdates []types.UpdateResult,
	expectedUpdatedPerpetualPositions map[types.SubaccountId][]*types.PerpetualPosition,
	expectedSubaccoundIdToFundingPayments map[types.SubaccountId]map[uint32]dtypes.SerializableInt,
	expectedUpdatedAssetPositions map[types.SubaccountId][]*types.AssetPosition,
	expectedAssetYieldIndex string,
) {
	subaccountUpdates := testutil.GetSubaccountUpdateEventsFromIndexerBlock(ctx, k)

	// No subaccount update events included in the case of an error or failure to update subaccounts.
	if expectedErr != nil || !expectedSuccess {
		require.Empty(t, subaccountUpdates)
		return
	}

	numSuccessfulUpdates := 0
	for idx := range updates {
		updateResult := expectedSuccessPerUpdates[idx]
		if updateResult != types.Success {
			continue
		}
		numSuccessfulUpdates += 1
	}

	// There should be exactly as many subaccount update events included as there were successful
	// subaccount updates.
	require.Len(t, subaccountUpdates, numSuccessfulUpdates)

	// For each update, verify that the expected SubaccountUpdateEvent is emitted.
	for _, update := range updates {
		expectedSubaccountUpdateEvent := indexerevents.NewSubaccountUpdateEvent(
			&update.SubaccountId,
			expectedUpdatedPerpetualPositions[update.SubaccountId],
			expectedUpdatedAssetPositions[update.SubaccountId],
			expectedSubaccoundIdToFundingPayments[update.SubaccountId],
			expectedAssetYieldIndex,
		)
		for _, gotUpdate := range subaccountUpdates {
			if gotUpdate.SubaccountId.Owner == expectedSubaccountUpdateEvent.SubaccountId.Owner &&
				gotUpdate.SubaccountId.Number == expectedSubaccountUpdateEvent.SubaccountId.Number {
				require.Equal(t,
					expectedSubaccountUpdateEvent,
					gotUpdate,
				)
			}
		}
	}
}

func TestGetCollateralPool(t *testing.T) {
	tests := map[string]struct {
		// state
		perpetuals         []perptypes.Perpetual
		perpetualPositions []*types.PerpetualPosition

		expectedAddress sdk.AccAddress
	}{
		"collateral pool with cross margin markets": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_SmallMarginRequirement,
			},
			perpetualPositions: []*types.PerpetualPosition{
				&constants.PerpetualPosition_OneBTCLong,
			},
			expectedAddress: authtypes.NewModuleAddress(types.ModuleName),
		},
		"collateral pool with isolated margin markets": {
			perpetuals: []perptypes.Perpetual{
				constants.IsoUsd_IsolatedMarket,
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId: constants.IsoUsd_IsolatedMarket.GetId(),
					Quantums:    dtypes.NewInt(100_000_000),
				},
			},
			expectedAddress: authtypes.NewModuleAddress(
				types.ModuleName + ":" + lib.UintToString(constants.IsoUsd_IsolatedMarket.GetId()),
			),
		},
		"collateral pool with no positions": {
			perpetualPositions: make([]*types.PerpetualPosition, 0),
			expectedAddress:    authtypes.NewModuleAddress(types.ModuleName),
		},
	}
	for name, tc := range tests {
		t.Run(
			name, func(t *testing.T) {
				ctx, keeper, pricesKeeper, perpetualsKeeper, _, _, assetsKeeper, rateLimitKeeper, _, _ := testutil.SubaccountsKeepers(
					t,
					true,
				)

				testutil.CreateTestMarkets(t, ctx, pricesKeeper)
				testutil.CreateTestLiquidityTiers(t, ctx, perpetualsKeeper)

				rateLimitKeeper.SetAssetYieldIndex(ctx, big.NewRat(1, 1))

				require.NoError(t, testutil.CreateTDaiAsset(ctx, assetsKeeper))
				for _, p := range tc.perpetuals {
					_, err := perpetualsKeeper.CreatePerpetual(
						ctx,
						p.Params.Id,
						p.Params.Ticker,
						p.Params.MarketId,
						p.Params.AtomicResolution,
						p.Params.DefaultFundingPpm,
						p.Params.LiquidityTier,
						p.Params.MarketType,
						p.Params.DangerIndexPpm,
						p.Params.IsolatedMarketMaxCumulativeInsuranceFundDeltaPerBlock,
						p.YieldIndex,
					)
					require.NoError(t, err)
				}

				subaccount := createNSubaccount(keeper, ctx, 1, big.NewInt(1_000))[0]
				subaccount.PerpetualPositions = tc.perpetualPositions
				keeper.SetSubaccount(ctx, subaccount)
				collateralPoolAddr, err := keeper.GetCollateralPoolForSubaccount(ctx, *subaccount.Id)
				require.NoError(t, err)
				require.Equal(t, tc.expectedAddress, collateralPoolAddr)
			},
		)
	}
}

func TestSubaccountGet(t *testing.T) {
	ctx, keeper, _, _, _, _, _, rateLimitKeeper, _, _ := testutil.SubaccountsKeepers(t, true)
	rateLimitKeeper.SetAssetYieldIndex(ctx, big.NewRat(1, 1))
	items := createNSubaccount(keeper, ctx, 10, big.NewInt(1_000))

	for _, item := range items {
		rst := keeper.GetSubaccount(ctx,
			*item.Id,
		)
		require.Equal(t,
			nullify.Fill(&item), //nolint:staticcheck
			nullify.Fill(&rst),  //nolint:staticcheck
		)
	}
}

func TestSubaccountSet_Empty(t *testing.T) {
	ctx, keeper, _, _, _, _, _, _, _, _ := testutil.SubaccountsKeepers(t, true)
	keeper.SetSubaccount(ctx, types.Subaccount{
		Id: &constants.Alice_Num0,
	})

	require.Len(t, keeper.GetAllSubaccount(ctx), 0)

	keeper.SetSubaccount(ctx, types.Subaccount{
		Id:             &constants.Alice_Num0,
		AssetPositions: testutil.CreateTDaiAssetPosition(big.NewInt(1_000)),
	})
	keeper.SetSubaccount(ctx, types.Subaccount{
		Id: &constants.Alice_Num0,
	})
	require.Len(t, keeper.GetAllSubaccount(ctx), 0)
}

func TestSubaccountGetNonExistent(t *testing.T) {
	ctx, keeper, _, _, _, _, _, rateLimitKeeper, _, _ := testutil.SubaccountsKeepers(t, true)
	rateLimitKeeper.SetAssetYieldIndex(ctx, big.NewRat(1, 1))
	id := types.SubaccountId{
		Owner:  "non-existent",
		Number: uint32(123),
	}
	acct := keeper.GetSubaccount(ctx, id)
	require.Equal(t, &id, acct.Id)
	require.Equal(t, new(big.Int), acct.GetTDaiPosition())
	require.Empty(t, acct.AssetPositions)
	require.Empty(t, acct.PerpetualPositions)
	require.False(t, acct.MarginEnabled)
}

func TestGetAllSubaccount(t *testing.T) {
	ctx, keeper, _, _, _, _, _, _, _, _ := testutil.SubaccountsKeepers(t, true)
	items := createNSubaccount(keeper, ctx, 10, big.NewInt(1_000))
	require.Equal(
		t,
		items,
		keeper.GetAllSubaccount(ctx),
	)
}

func TestForEachSubaccount(t *testing.T) {
	tests := map[string]struct {
		numSubaccountsInState int
		iterationCount        int
	}{
		"No subaccounts in state": {
			numSubaccountsInState: 0,
			iterationCount:        0,
		},
		"one subaccount in state, one iteration": {
			numSubaccountsInState: 1,
			iterationCount:        1,
		},
		"two subaccount in state, one iteration": {
			numSubaccountsInState: 2,
			iterationCount:        1,
		},
		"ten subaccount in state, one iteration": {
			numSubaccountsInState: 10,
			iterationCount:        1,
		},
		"ten subaccount in state, partial iteration": {
			numSubaccountsInState: 10,
			iterationCount:        8,
		},
		"ten subaccount in state, full iteration": {
			numSubaccountsInState: 10,
			iterationCount:        10,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctx, keeper, _, _, _, _, _, _, _, _ := testutil.SubaccountsKeepers(t, true)
			items := createNSubaccount(keeper, ctx, tc.numSubaccountsInState, big.NewInt(1_000))
			collectedSubaccounts := make([]types.Subaccount, 0)
			i := 0
			keeper.ForEachSubaccount(ctx, func(subaccount types.Subaccount) bool {
				i++
				collectedSubaccounts = append(collectedSubaccounts, subaccount)
				return i == tc.iterationCount
			})
			require.Equal(
				t,
				items[:tc.iterationCount],
				collectedSubaccounts,
			)
		})
	}
}

func TestUpdateSubaccounts(t *testing.T) {
	// default subaccount id, the first subaccount id generated when calling createNSubaccount
	defaultSubaccountId := types.SubaccountId{
		Owner:  "0",
		Number: 0,
	}

	tests := map[string]struct {
		// state
		perpetuals        []perptypes.Perpetual
		newFundingIndices []*big.Int // 1:1 mapped to perpetuals list
		assets            []*asstypes.Asset
		marketParamPrices []pricestypes.MarketParamPrice
		// If not specified, default to `CollatCheck`
		updateType                types.UpdateType
		additionalTestSubaccounts []types.Subaccount
		// Only set when specified. Defaults to 0/1.
		// Set perpYieldIndex in the perpetuals state.
		globalAssetYieldIndex *big.Rat
		fundsInTDaiPool       *big.Int

		// subaccount state
		perpetualPositions        []*types.PerpetualPosition
		assetPositions            []*types.AssetPosition
		subaccountAssetYieldIndex string

		// collateral pool state
		collateralPoolTDaiBalances map[string]int64

		// updates
		updates []types.Update

		// expectations
		expectedCollateralPoolTDaiBalances map[string]int64
		expectedQuoteBalance               *big.Int
		expectedPerpetualPositions         []*types.PerpetualPosition
		expectedAssetPositions             []*types.AssetPosition
		expectedSuccess                    bool
		expectedSuccessPerUpdate           []types.UpdateResult
		expectedTDaiYieldPoolBalance       *big.Int
		expectedErr                        error
		// List of expected open interest.
		// If not specified, this means OI is default value.
		expectedOpenInterest map[uint32]*big.Int

		// Only contains the updated perpetual positions, to assert against the events included.
		expectedUpdatedPerpetualPositions     map[types.SubaccountId][]*types.PerpetualPosition
		expectedSubaccoundIdToFundingPayments map[types.SubaccountId]map[uint32]dtypes.SerializableInt
		expectedUpdatedAssetPositions         map[types.SubaccountId][]*types.AssetPosition
		expectedAssetYieldIndex               string
		msgSenderEnabled                      bool
	}{
		"one update to TDai asset position": {
			expectedQuoteBalance:     big.NewInt(100),
			expectedSuccess:          true,
			expectedSuccessPerUpdate: []types.UpdateResult{types.Success},
			updates: []types.Update{
				{
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(100)),
				},
			},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(100), // 100 TDai
				},
			},
			expectedUpdatedAssetPositions: map[types.SubaccountId][]*types.AssetPosition{
				defaultSubaccountId: {
					{
						AssetId:  uint32(0),
						Quantums: dtypes.NewInt(100), // 100 TDai
					},
				},
			},
			expectedAssetYieldIndex: constants.AssetYieldIndex_Zero,
			msgSenderEnabled:        true,
		},
		"one negative update to TDai asset position": {
			expectedQuoteBalance:     big.NewInt(-100),
			expectedSuccess:          true,
			expectedSuccessPerUpdate: []types.UpdateResult{types.Success},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_SmallMarginRequirement,
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000), // 1 BTC
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000), // 1 BTC
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(-100), // 100 TDai
				},
			},
			expectedUpdatedAssetPositions: map[types.SubaccountId][]*types.AssetPosition{
				defaultSubaccountId: {
					{
						AssetId:  uint32(0),
						Quantums: dtypes.NewInt(-100), // 100 TDai
					},
				},
			},
			updates: []types.Update{
				{
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(-100)),
				},
			},
			expectedAssetYieldIndex: constants.AssetYieldIndex_Zero,
			msgSenderEnabled:        true,
		},
		"one negative update to TDai asset position + persist unsettled negative funding": {
			expectedQuoteBalance:     big.NewInt(-2100),
			expectedSuccess:          true,
			expectedSuccessPerUpdate: []types.UpdateResult{types.Success},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_SmallMarginRequirement,
			},
			newFundingIndices: []*big.Int{big.NewInt(-10)},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000), // 1 BTC
					FundingIndex: dtypes.NewInt(-30),         // indexDelta=20, settlement=-20*100
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000), // 1 BTC
					FundingIndex: dtypes.NewInt(-10),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(-2100), // 2100 TDai
				},
			},
			expectedUpdatedPerpetualPositions: map[types.SubaccountId][]*types.PerpetualPosition{
				defaultSubaccountId: {
					{
						PerpetualId:  uint32(0),
						Quantums:     dtypes.NewInt(100_000_000),
						FundingIndex: dtypes.NewInt(-10),
					},
				},
			},
			expectedSubaccoundIdToFundingPayments: map[types.SubaccountId]map[uint32]dtypes.SerializableInt{
				defaultSubaccountId: {
					uint32(0): dtypes.NewInt(2_000), // negated settlement
				},
			},
			expectedUpdatedAssetPositions: map[types.SubaccountId][]*types.AssetPosition{
				defaultSubaccountId: {
					{
						AssetId:  uint32(0),
						Quantums: dtypes.NewInt(-2100), // 2100 TDai
					},
				},
			},
			updates: []types.Update{
				{
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(-100)),
				},
			},
			expectedAssetYieldIndex: constants.AssetYieldIndex_Zero,
			msgSenderEnabled:        true,
		},
		"one negative update to TDai asset position + persist unsettled positive funding": {
			expectedQuoteBalance:     big.NewInt(-92),
			expectedSuccess:          true,
			expectedSuccessPerUpdate: []types.UpdateResult{types.Success},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_SmallMarginRequirement,
			},
			newFundingIndices: []*big.Int{big.NewInt(-17)},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId: uint32(0),
					Quantums:    dtypes.NewInt(500_000), // 0.005 BTC
					// indexDelta=-17, settlement=17*500_000/1_000_000=8
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(500_000), // 1 BTC
					FundingIndex: dtypes.NewInt(-17),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(-92),
				},
			},
			expectedUpdatedPerpetualPositions: map[types.SubaccountId][]*types.PerpetualPosition{
				defaultSubaccountId: {
					{
						PerpetualId:  uint32(0),
						Quantums:     dtypes.NewInt(500_000),
						FundingIndex: dtypes.NewInt(-17),
					},
				},
			},
			expectedSubaccoundIdToFundingPayments: map[types.SubaccountId]map[uint32]dtypes.SerializableInt{
				defaultSubaccountId: {
					uint32(0): dtypes.NewInt(-8), // negated settlement
				},
			},
			expectedUpdatedAssetPositions: map[types.SubaccountId][]*types.AssetPosition{
				defaultSubaccountId: {
					{
						AssetId:  uint32(0),
						Quantums: dtypes.NewInt(-92),
					},
				},
			},
			updates: []types.Update{
				{
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(-100)),
				},
			},
			expectedAssetYieldIndex: constants.AssetYieldIndex_Zero,
			msgSenderEnabled:        true,
		},
		"multiple updates for same position not allowed": {
			expectedQuoteBalance:     big.NewInt(0),
			expectedSuccess:          false,
			expectedSuccessPerUpdate: nil,
			expectedErr:              types.ErrNonUniqueUpdatesPosition,
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_SmallMarginRequirement,
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000), // 1 BTC
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000), // 1 BTC
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			updates: []types.Update{
				{
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(0),
							BigQuantumsDelta: big.NewInt(9_900_000_000), // 99 BTC
						},
						{
							PerpetualId:      uint32(0),
							BigQuantumsDelta: big.NewInt(9_900_000_000), // 99 BTC
						},
					},
				},
			},
			expectedAssetYieldIndex: constants.AssetYieldIndex_Zero,
			msgSenderEnabled:        true,
		},
		"multiple updates to same account not allowed": {
			expectedQuoteBalance:     big.NewInt(0),
			expectedErr:              types.ErrNonUniqueUpdatesSubaccount,
			expectedSuccess:          false,
			expectedSuccessPerUpdate: nil,
			updates: []types.Update{
				{
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(-100)),
				},
				{
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(-100)),
				},
			},
			expectedAssetYieldIndex: constants.AssetYieldIndex_Zero,
			msgSenderEnabled:        true,
		},
		"update increases position size": {
			assetPositions:           testutil.CreateTDaiAssetPosition(big.NewInt(25_000_000_000)), // $25,000
			expectedQuoteBalance:     big.NewInt(0),
			expectedSuccess:          true,
			expectedSuccessPerUpdate: []types.UpdateResult{types.Success},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_NoMarginRequirement,
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000), // 1 BTC
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(150_000_000), // 1.5 BTC
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			expectedUpdatedPerpetualPositions: map[types.SubaccountId][]*types.PerpetualPosition{
				defaultSubaccountId: {
					{
						PerpetualId:  uint32(0),
						Quantums:     dtypes.NewInt(150_000_000), // 1.5 BTC
						FundingIndex: dtypes.NewInt(0),
					},
				},
			},
			expectedAssetPositions: []*types.AssetPosition{},
			expectedUpdatedAssetPositions: map[types.SubaccountId][]*types.AssetPosition{
				defaultSubaccountId: {
					{
						AssetId:  uint32(0),
						Quantums: dtypes.NewInt(0),
					},
				},
			},
			updates: []types.Update{
				{
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(-25_000_000_000)), // -$25,000
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(0),
							BigQuantumsDelta: big.NewInt(50_000_000), // .5 BTC
						},
					},
				},
			},
			msgSenderEnabled: false,
		},
		"update decreases position size": {
			assetPositions:           testutil.CreateTDaiAssetPosition(big.NewInt(25_000_000_000)), // $25,000
			expectedQuoteBalance:     big.NewInt(50_000_000_000),                                   // $50,000
			expectedSuccess:          true,
			expectedSuccessPerUpdate: []types.UpdateResult{types.Success},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_NoMarginRequirement,
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000), // 1 BTC
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(50_000_000), // .50 BTC
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			expectedUpdatedPerpetualPositions: map[types.SubaccountId][]*types.PerpetualPosition{
				defaultSubaccountId: {
					{
						PerpetualId:  uint32(0),
						Quantums:     dtypes.NewInt(50_000_000), // .50 BTC
						FundingIndex: dtypes.NewInt(0),
					},
				},
			},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(50_000_000_000),
				},
			},
			expectedUpdatedAssetPositions: map[types.SubaccountId][]*types.AssetPosition{
				defaultSubaccountId: {
					{
						AssetId:  uint32(0),
						Quantums: dtypes.NewInt(50_000_000_000),
					},
				},
			},
			updates: []types.Update{
				{
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(25_000_000_000)), // $25,000
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(0),
							BigQuantumsDelta: big.NewInt(-50_000_000), // -.5 BTC
						},
					},
				},
			},
			msgSenderEnabled: false,
		},
		"update closes long position": {
			assetPositions:           testutil.CreateTDaiAssetPosition(big.NewInt(25_000_000_000)), // $25,000
			expectedQuoteBalance:     big.NewInt(75_000_000_000),                                   // $75,000
			expectedSuccess:          true,
			expectedSuccessPerUpdate: []types.UpdateResult{types.Success},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_NoMarginRequirement,
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000), // 1 BTC
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			expectedPerpetualPositions: []*types.PerpetualPosition{},
			expectedUpdatedPerpetualPositions: map[types.SubaccountId][]*types.PerpetualPosition{
				defaultSubaccountId: {
					// Position closed update.
					{
						PerpetualId:  uint32(0),
						Quantums:     dtypes.NewInt(0),
						FundingIndex: dtypes.NewInt(0),
					},
				},
			},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(75_000_000_000),
				},
			},
			expectedUpdatedAssetPositions: map[types.SubaccountId][]*types.AssetPosition{
				defaultSubaccountId: {
					{
						AssetId:  uint32(0),
						Quantums: dtypes.NewInt(75_000_000_000),
					},
				},
			},
			updates: []types.Update{
				{
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(50_000_000_000)), // $50,000
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(0),
							BigQuantumsDelta: big.NewInt(-100_000_000), // -1 BTC
						},
					},
				},
			},
			expectedAssetYieldIndex: constants.AssetYieldIndex_Zero,
			msgSenderEnabled:        true,
		},
		"update closes short position": {
			assetPositions:           testutil.CreateTDaiAssetPosition(big.NewInt(100_000_000_000)), // $100,000
			expectedQuoteBalance:     big.NewInt(50_000_000_000),                                    // $50,000
			expectedSuccess:          true,
			expectedSuccessPerUpdate: []types.UpdateResult{types.Success},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_NoMarginRequirement,
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(-100_000_000), // -1 BTC
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			expectedPerpetualPositions: []*types.PerpetualPosition{},
			expectedUpdatedPerpetualPositions: map[types.SubaccountId][]*types.PerpetualPosition{
				defaultSubaccountId: {
					// Position closed update.
					{
						PerpetualId:  uint32(0),
						Quantums:     dtypes.NewInt(0),
						FundingIndex: dtypes.NewInt(0),
					},
				},
			},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(50_000_000_000),
				},
			},
			expectedUpdatedAssetPositions: map[types.SubaccountId][]*types.AssetPosition{
				defaultSubaccountId: {
					{
						AssetId:  uint32(0),
						Quantums: dtypes.NewInt(50_000_000_000),
					},
				},
			},
			updates: []types.Update{
				{
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(-50_000_000_000)), // -$50,000
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(0),
							BigQuantumsDelta: big.NewInt(100_000_000), // 1 BTC
						},
					},
				},
			},
			expectedAssetYieldIndex: constants.AssetYieldIndex_Zero,
			msgSenderEnabled:        true,
		},
		"update closes 2nd position and updates 1st": {
			assetPositions:           testutil.CreateTDaiAssetPosition(big.NewInt(100_000_000_000)), // $100,000
			expectedSuccess:          true,
			expectedSuccessPerUpdate: []types.UpdateResult{types.Success},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_NoMarginRequirement,
				constants.EthUsd_NoMarginRequirement,
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(-100_000_000), // -1 BTC
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
				{
					PerpetualId:  uint32(1),
					Quantums:     dtypes.NewInt(-1_000_000_000_000_000_000), // -1 ETH
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(-200_000_000), // -2 BTC
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			expectedUpdatedPerpetualPositions: map[types.SubaccountId][]*types.PerpetualPosition{
				defaultSubaccountId: {
					{
						PerpetualId:  uint32(0),
						Quantums:     dtypes.NewInt(-200_000_000), // -2 BTC
						FundingIndex: dtypes.NewInt(0),
					},
					// Position closed update.
					{
						PerpetualId:  uint32(1),
						Quantums:     dtypes.NewInt(0),
						FundingIndex: dtypes.NewInt(0),
					},
				},
			},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(100_000_000_000),
				},
			},
			updates: []types.Update{
				{
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(0),
							BigQuantumsDelta: big.NewInt(-100_000_000), // -1 BTC
						},
						{
							PerpetualId:      uint32(1),
							BigQuantumsDelta: big.NewInt(1_000_000_000_000_000_000), // 1 ETH
						},
					},
				},
			},
			expectedAssetYieldIndex: constants.AssetYieldIndex_Zero,
			msgSenderEnabled:        true,
		},
		"update closes first asset position and updates 2nd": {
			assets: []*asstypes.Asset{
				constants.BtcUsd,
			},
			assetPositions: append(
				testutil.CreateTDaiAssetPosition(big.NewInt(100_000_000_000)), // $100,000
				&types.AssetPosition{
					AssetId:  constants.BtcUsd.Id,
					Quantums: dtypes.NewInt(50_000),
				},
			),
			expectedQuoteBalance:     big.NewInt(200_000_000_000), // $200,000
			expectedSuccess:          true,
			expectedSuccessPerUpdate: []types.UpdateResult{types.Success},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(200_000_000_000),
				},
			},
			expectedUpdatedAssetPositions: map[types.SubaccountId][]*types.AssetPosition{
				defaultSubaccountId: {
					{
						AssetId:  uint32(0),
						Quantums: dtypes.NewInt(200_000_000_000),
					},
					// Asset position closed
					{
						AssetId:  uint32(1),
						Quantums: dtypes.NewInt(0),
					},
				},
			},
			updates: []types.Update{
				{
					AssetUpdates: []types.AssetUpdate{
						{
							AssetId:          asstypes.AssetTDai.Id,
							BigQuantumsDelta: big.NewInt(100_000_000_000),
						},
						{
							AssetId:          constants.BtcUsd.Id,
							BigQuantumsDelta: big.NewInt(-50_000),
						},
					},
				},
			},
			expectedAssetYieldIndex: constants.AssetYieldIndex_Zero,
			msgSenderEnabled:        true,
		},
		"update closes first 1 positions and updates 2nd": {
			assetPositions:           testutil.CreateTDaiAssetPosition(big.NewInt(100_000_000_000)), // $100,000
			expectedQuoteBalance:     big.NewInt(50_000_000_000),                                    // $50,000
			expectedSuccess:          true,
			expectedSuccessPerUpdate: []types.UpdateResult{types.Success},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_NoMarginRequirement,
				constants.EthUsd_NoMarginRequirement,
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(-100_000_000), // -1 BTC
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
				{
					PerpetualId:  uint32(1),
					Quantums:     dtypes.NewInt(-1_000_000_000), // -1 ETH
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(1),
					Quantums:     dtypes.NewInt(-2_000_000_000), // -2 ETH
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			expectedUpdatedPerpetualPositions: map[types.SubaccountId][]*types.PerpetualPosition{
				defaultSubaccountId: {
					// Position closed update.
					{
						PerpetualId:  uint32(0),
						Quantums:     dtypes.NewInt(0),
						FundingIndex: dtypes.NewInt(0),
					},
					{
						PerpetualId:  uint32(1),
						Quantums:     dtypes.NewInt(-2_000_000_000), // -2 ETH
						FundingIndex: dtypes.NewInt(0),
					},
				},
			},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(50_000_000_000),
				},
			},
			expectedUpdatedAssetPositions: map[types.SubaccountId][]*types.AssetPosition{
				defaultSubaccountId: {
					{
						AssetId:  uint32(0),
						Quantums: dtypes.NewInt(50_000_000_000),
					},
				},
			},
			updates: []types.Update{
				{
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(-50_000_000_000)), // -$50,000
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(0),
							BigQuantumsDelta: big.NewInt(100_000_000), // 1 BTC
						},
						{
							PerpetualId:      uint32(1),
							BigQuantumsDelta: big.NewInt(-1_000_000_000), // -1 ETH
						},
					},
				},
			},
			expectedAssetYieldIndex: constants.AssetYieldIndex_Zero,
			msgSenderEnabled:        true,
		},
		"update opens new long position, uses current perpetual funding index": {
			assetPositions:           testutil.CreateTDaiAssetPosition(big.NewInt(100_000_000_000)), // $100,000
			expectedQuoteBalance:     big.NewInt(50_000_000_000),                                    // $50,000
			expectedSuccess:          true,
			expectedSuccessPerUpdate: []types.UpdateResult{types.Success},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_NoMarginRequirement,
			},
			newFundingIndices:  []*big.Int{big.NewInt(-15)},
			perpetualPositions: []*types.PerpetualPosition{},
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000), // 1 BTC
					FundingIndex: dtypes.NewInt(-15),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			expectedUpdatedPerpetualPositions: map[types.SubaccountId][]*types.PerpetualPosition{
				defaultSubaccountId: {
					{
						PerpetualId:  uint32(0),
						Quantums:     dtypes.NewInt(100_000_000), // 1 BTC
						FundingIndex: dtypes.NewInt(-15),
					},
				},
			},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(50_000_000_000),
				},
			},
			expectedUpdatedAssetPositions: map[types.SubaccountId][]*types.AssetPosition{
				defaultSubaccountId: {
					{
						AssetId:  uint32(0),
						Quantums: dtypes.NewInt(50_000_000_000),
					},
				},
			},
			updates: []types.Update{
				{
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(-50_000_000_000)), // -$50,000
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(0),
							BigQuantumsDelta: big.NewInt(100_000_000), // 1 BTC
						},
					},
				},
			},
			msgSenderEnabled: false,
		},
		"update opens new short position": {
			assetPositions:           testutil.CreateTDaiAssetPosition(big.NewInt(100_000_000_000)), // $100,000
			expectedQuoteBalance:     big.NewInt(150_000_000_000),                                   // $50,000
			expectedSuccess:          true,
			expectedSuccessPerUpdate: []types.UpdateResult{types.Success},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_NoMarginRequirement,
			},
			perpetualPositions: []*types.PerpetualPosition{},
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(-100_000_000), // 1 BTC
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			expectedUpdatedPerpetualPositions: map[types.SubaccountId][]*types.PerpetualPosition{
				defaultSubaccountId: {
					{
						PerpetualId:  uint32(0),
						Quantums:     dtypes.NewInt(-100_000_000), // 1 BTC
						FundingIndex: dtypes.NewInt(0),
					},
				},
			},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(150_000_000_000),
				},
			},
			expectedUpdatedAssetPositions: map[types.SubaccountId][]*types.AssetPosition{
				defaultSubaccountId: {
					{
						AssetId:  uint32(0),
						Quantums: dtypes.NewInt(150_000_000_000),
					},
				},
			},
			updates: []types.Update{
				{
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(50_000_000_000)), // $50,000
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(0),
							BigQuantumsDelta: big.NewInt(-100_000_000), // -1 BTC
						},
					},
				},
			},
			msgSenderEnabled: false,
		},
		"update opens new long eth position with existing btc position": {
			assetPositions:           testutil.CreateTDaiAssetPosition(big.NewInt(100_000_000_000)), // $100,000
			expectedQuoteBalance:     big.NewInt(100_000_000_000),                                   // $100,000
			expectedSuccess:          true,
			expectedSuccessPerUpdate: []types.UpdateResult{types.Success},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_NoMarginRequirement,
				constants.EthUsd_NoMarginRequirement,
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000), // 1 BTC
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000), // 1 BTC
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
				{
					PerpetualId:  uint32(1),
					Quantums:     dtypes.NewInt(1_000_000_000), // 1 ETH
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			expectedUpdatedPerpetualPositions: map[types.SubaccountId][]*types.PerpetualPosition{
				defaultSubaccountId: {
					{
						PerpetualId:  uint32(1),
						Quantums:     dtypes.NewInt(1_000_000_000), // 1 ETH
						FundingIndex: dtypes.NewInt(0),
					},
				},
			},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(100_000_000_000),
				},
			},
			updates: []types.Update{
				{
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(1),
							BigQuantumsDelta: big.NewInt(1_000_000_000), // 1 ETH
						},
					},
				},
			},
			expectedAssetYieldIndex: constants.AssetYieldIndex_Zero,
			msgSenderEnabled:        true,
		},
		// TODO(DEC-581): add similar test case for multi-collateral asset support.
		"update eth position from long to short with existing btc position": {
			assetPositions:           testutil.CreateTDaiAssetPosition(big.NewInt(100_000_000_000)), // $100,000
			expectedQuoteBalance:     big.NewInt(100_000_000_000),                                   // $100,000
			expectedSuccess:          true,
			expectedSuccessPerUpdate: []types.UpdateResult{types.Success},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_NoMarginRequirement,
				constants.EthUsd_NoMarginRequirement,
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000), // 1 BTC
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
				{
					PerpetualId:  uint32(1),
					Quantums:     dtypes.NewInt(500_000_000), // 5 ETH
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000), // 1 BTC
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
				{
					PerpetualId:  uint32(1),
					Quantums:     dtypes.NewInt(-500_000_000), // -5 ETH
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			expectedUpdatedPerpetualPositions: map[types.SubaccountId][]*types.PerpetualPosition{
				defaultSubaccountId: {
					{
						PerpetualId:  uint32(1),
						Quantums:     dtypes.NewInt(-500_000_000), // -5 ETH
						FundingIndex: dtypes.NewInt(0),
					},
				},
			},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(100_000_000_000),
				},
			},
			updates: []types.Update{
				{
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(1),
							BigQuantumsDelta: big.NewInt(-1_000_000_000), // -10 ETH
						},
					},
				},
			},
			expectedAssetYieldIndex: constants.AssetYieldIndex_Zero,
			msgSenderEnabled:        true,
		},
		"update opens new long eth position with existing btc and sol position": {
			assetPositions:           testutil.CreateTDaiAssetPosition(big.NewInt(100_000_000_000)), // $100,000
			expectedQuoteBalance:     big.NewInt(100_000_000_000),                                   // $100,000
			expectedSuccess:          true,
			expectedSuccessPerUpdate: []types.UpdateResult{types.Success},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_NoMarginRequirement,
				constants.EthUsd_20PercentInitial_10PercentMaintenance,
				constants.SolUsd_20PercentInitial_10PercentMaintenance,
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000), // 1 BTC
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
				{
					PerpetualId:  uint32(2),
					Quantums:     dtypes.NewInt(1_000_000_000), // 1 SOL
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000), // 1 BTC
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
				{
					PerpetualId:  uint32(1),
					Quantums:     dtypes.NewInt(1_000_000_000), // 1 ETH
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
				{
					PerpetualId:  uint32(2),
					Quantums:     dtypes.NewInt(1_000_000_000), // 1 SOL
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			expectedUpdatedPerpetualPositions: map[types.SubaccountId][]*types.PerpetualPosition{
				defaultSubaccountId: {
					{
						PerpetualId:  uint32(1),
						Quantums:     dtypes.NewInt(1_000_000_000), // 1 ETH
						FundingIndex: dtypes.NewInt(0),
					},
				},
			},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(100_000_000_000),
				},
			},
			updates: []types.Update{
				{
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(1),
							BigQuantumsDelta: big.NewInt(1_000_000_000), // 1 ETH
						},
					},
				},
			},
			expectedAssetYieldIndex: constants.AssetYieldIndex_Zero,
			msgSenderEnabled:        true,
		},
		"update opens new long btc position with existing eth and sol position": {
			assetPositions:           testutil.CreateTDaiAssetPosition(big.NewInt(100_000_000_000)), // $100,000
			expectedQuoteBalance:     big.NewInt(100_000_000_000),                                   // $100,000
			expectedSuccess:          true,
			expectedSuccessPerUpdate: []types.UpdateResult{types.Success},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_NoMarginRequirement,
				constants.EthUsd_20PercentInitial_10PercentMaintenance,
				constants.SolUsd_20PercentInitial_10PercentMaintenance,
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(1),
					Quantums:     dtypes.NewInt(1_000_000_000), // 1 ETH
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
				{
					PerpetualId:  uint32(2),
					Quantums:     dtypes.NewInt(1_000_000_000), // 1 SOL
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000), // 1 BTC
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
				{
					PerpetualId:  uint32(1),
					Quantums:     dtypes.NewInt(1_000_000_000), // 1 ETH
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
				{
					PerpetualId:  uint32(2),
					Quantums:     dtypes.NewInt(1_000_000_000), // 1 SOL
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			expectedUpdatedPerpetualPositions: map[types.SubaccountId][]*types.PerpetualPosition{
				defaultSubaccountId: {
					{
						PerpetualId:  uint32(0),
						Quantums:     dtypes.NewInt(100_000_000), // 1 BTC
						FundingIndex: dtypes.NewInt(0),
					},
				},
			},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(100_000_000_000),
				},
			},
			updates: []types.Update{
				{
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(0),
							BigQuantumsDelta: big.NewInt(100_000_000), // 1 BTC
						},
					},
				},
			},
			expectedAssetYieldIndex: constants.AssetYieldIndex_Zero,
			msgSenderEnabled:        true,
		},
		"update opens new long eth position with existing unsettled sol position": {
			assetPositions:           testutil.CreateTDaiAssetPosition(big.NewInt(100_000_000_000)), // $100,000
			expectedQuoteBalance:     big.NewInt(100_000_000_000),                                   // $100,000
			expectedSuccess:          true,
			expectedSuccessPerUpdate: []types.UpdateResult{types.Success},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_NoMarginRequirement,
				constants.EthUsd_20PercentInitial_10PercentMaintenance,
				constants.SolUsd_20PercentInitial_10PercentMaintenance,
			},
			newFundingIndices: []*big.Int{
				big.NewInt(1234),  // btc
				big.NewInt(-5000), // eth
				big.NewInt(2000),  // sol
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(2),
					Quantums:     dtypes.NewInt(1_000_000_000), // 1 SOL
					FundingIndex: dtypes.NewInt(1700),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(1),
					Quantums:     dtypes.NewInt(1_000_000_000), // 1 ETH
					FundingIndex: dtypes.NewInt(-5000),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
				{
					PerpetualId:  uint32(2),
					Quantums:     dtypes.NewInt(1_000_000_000), // 1 SOL
					FundingIndex: dtypes.NewInt(2000),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			expectedUpdatedPerpetualPositions: map[types.SubaccountId][]*types.PerpetualPosition{
				defaultSubaccountId: {
					{
						PerpetualId:  uint32(1),
						Quantums:     dtypes.NewInt(1_000_000_000), // 1 ETH
						FundingIndex: dtypes.NewInt(-5000),
					},
					{
						PerpetualId:  uint32(2),
						Quantums:     dtypes.NewInt(1_000_000_000), // 1 SOL
						FundingIndex: dtypes.NewInt(2000),
					},
				},
			},
			expectedSubaccoundIdToFundingPayments: map[types.SubaccountId]map[uint32]dtypes.SerializableInt{
				defaultSubaccountId: {
					uint32(2): dtypes.NewInt(300_000),
				},
			},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(99_999_700_000),
				},
			},
			updates: []types.Update{
				{
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(1),
							BigQuantumsDelta: big.NewInt(1_000_000_000), // 1 ETH
						},
					},
				},
			},
			expectedAssetYieldIndex: constants.AssetYieldIndex_Zero,
			msgSenderEnabled:        true,
		},
		"provides out-of-order updates (not ordered by PerpetualId)": {
			assetPositions:           testutil.CreateTDaiAssetPosition(big.NewInt(100_000_000_000)), // $100,000
			expectedQuoteBalance:     big.NewInt(100_000_000_000),                                   // $100,000
			expectedSuccess:          true,
			expectedSuccessPerUpdate: []types.UpdateResult{types.Success},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_NoMarginRequirement,
				constants.EthUsd_20PercentInitial_10PercentMaintenance,
				constants.SolUsd_20PercentInitial_10PercentMaintenance,
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000), // 1 BTC
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
				{
					PerpetualId:  uint32(1),
					Quantums:     dtypes.NewInt(1_000_000_000), // 1 ETH
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
				{
					PerpetualId:  uint32(2),
					Quantums:     dtypes.NewInt(1_000_000_000), // 1 SOL
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(200_000_000), // 2 BTC
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
				{
					PerpetualId:  uint32(1),
					Quantums:     dtypes.NewInt(2_000_000_000), // 2 ETH
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
				{
					PerpetualId:  uint32(2),
					Quantums:     dtypes.NewInt(2_000_000_000), // 2 SOL
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			expectedUpdatedPerpetualPositions: map[types.SubaccountId][]*types.PerpetualPosition{
				defaultSubaccountId: {
					{
						PerpetualId:  uint32(0),
						Quantums:     dtypes.NewInt(200_000_000), // 2 BTC
						FundingIndex: dtypes.NewInt(0),
					},
					{
						PerpetualId:  uint32(1),
						Quantums:     dtypes.NewInt(2_000_000_000), // 2 ETH
						FundingIndex: dtypes.NewInt(0),
					},
					{
						PerpetualId:  uint32(2),
						Quantums:     dtypes.NewInt(2_000_000_000), // 2 SOL
						FundingIndex: dtypes.NewInt(0),
					},
				},
			},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(100_000_000_000),
				},
			},
			updates: []types.Update{
				{
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(2),
							BigQuantumsDelta: big.NewInt(1_000_000_000), // 1 SOL
						},
						{
							PerpetualId:      uint32(1),
							BigQuantumsDelta: big.NewInt(1_000_000_000), // 1 ETH
						},
						{
							PerpetualId:      uint32(0),
							BigQuantumsDelta: big.NewInt(100_000_000), // 1 BTC
						},
					},
				},
			},
			expectedAssetYieldIndex: constants.AssetYieldIndex_Zero,
			msgSenderEnabled:        true,
		},
		"updates multiple subaccounts with new perpetual and asset positions": {
			expectedQuoteBalance:     big.NewInt(100_000_000), // $100
			expectedSuccess:          true,
			expectedSuccessPerUpdate: []types.UpdateResult{types.Success, types.Success},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_NoMarginRequirement,
				constants.EthUsd_20PercentInitial_10PercentMaintenance,
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(1),
					Quantums:     dtypes.NewInt(1_000_000_000), // 1 ETH
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			expectedPerpetualPositions: []*types.PerpetualPosition{},
			expectedUpdatedPerpetualPositions: map[types.SubaccountId][]*types.PerpetualPosition{
				defaultSubaccountId: {
					// Position closed update.
					{
						PerpetualId:  uint32(1),
						Quantums:     dtypes.NewInt(0),
						FundingIndex: dtypes.NewInt(0),
					},
				},
				{
					Owner:  "non-existent account",
					Number: uint32(12),
				}: {
					{
						PerpetualId:  uint32(0),
						Quantums:     dtypes.NewInt(100_000_000), // 1 BTC
						FundingIndex: dtypes.NewInt(0),
					},
				},
			},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(100_000_000),
				},
			},
			expectedUpdatedAssetPositions: map[types.SubaccountId][]*types.AssetPosition{
				defaultSubaccountId: {
					{
						AssetId:  uint32(0),
						Quantums: dtypes.NewInt(100_000_000),
					},
				},
				{
					Owner:  "non-existent account",
					Number: uint32(12),
				}: {
					{
						AssetId:  uint32(0),
						Quantums: dtypes.NewInt(500_000_000),
					},
				},
			},
			updates: []types.Update{
				{
					SubaccountId: types.SubaccountId{
						Owner:  "non-existent account",
						Number: uint32(12),
					},
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(500_000_000)), // $500
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(0),
							BigQuantumsDelta: big.NewInt(100_000_000), // 1 BTC
						},
					},
				},
				{
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(100_000_000)), // $100
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(1),
							BigQuantumsDelta: big.NewInt(-1_000_000_000), // -1 ETH
						},
					},
				},
			},
			expectedAssetYieldIndex: constants.AssetYieldIndex_Zero,
			msgSenderEnabled:        true,
		},
		"update would make account undercollateralized": {
			expectedQuoteBalance:     big.NewInt(0),
			expectedSuccess:          false,
			expectedSuccessPerUpdate: []types.UpdateResult{types.NewlyUndercollateralized},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_SmallMarginRequirement,
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(1_000_000), // 0.01 BTC
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(1_000_000), // 0.01 BTC
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			updates: []types.Update{
				{
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(-50_000_000_000)), // -$50,000
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(0),
							BigQuantumsDelta: big.NewInt(100_000_000), // 1 BTC
						},
					},
				},
			},
			expectedAssetYieldIndex: constants.AssetYieldIndex_Zero,
			msgSenderEnabled:        true,
		},
		"updates new TDai asset position which exceeds max uint64": {
			assetPositions:           testutil.CreateTDaiAssetPosition(new(big.Int).SetUint64(math.MaxUint64)),
			expectedQuoteBalance:     new(big.Int).SetUint64(math.MaxUint64),
			expectedSuccess:          true,
			expectedSuccessPerUpdate: []types.UpdateResult{types.Success},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId: uint32(0),
					Quantums: dtypes.NewIntFromBigInt(
						new(big.Int).Add(
							new(big.Int).SetUint64(math.MaxUint64),
							new(big.Int).SetUint64(1),
						),
					),
				},
			},
			updates: []types.Update{
				{
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(1)),
				},
			},
			expectedUpdatedAssetPositions: map[types.SubaccountId][]*types.AssetPosition{
				defaultSubaccountId: {
					{
						AssetId: uint32(0),
						Quantums: dtypes.NewIntFromBigInt(
							new(big.Int).Add(
								new(big.Int).SetUint64(math.MaxUint64),
								new(big.Int).SetUint64(1),
							),
						),
					},
				},
			},
			expectedAssetYieldIndex: constants.AssetYieldIndex_Zero,
			msgSenderEnabled:        true,
		},
		"new TDai asset position (including unsettled funding) size exceeds max uint64": {
			assetPositions: testutil.CreateTDaiAssetPosition(new(big.Int).SetUint64(math.MaxUint64 - 5)),
			expectedQuoteBalance: new(big.Int).Add(
				new(big.Int).SetUint64(math.MaxUint64),
				new(big.Int).SetInt64(1),
			),
			expectedSuccess:          true,
			expectedSuccessPerUpdate: []types.UpdateResult{types.Success},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_SmallMarginRequirement,
			},
			newFundingIndices: []*big.Int{big.NewInt(-10)},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(1_000_000), // 0.01 BTC
					FundingIndex: dtypes.NewInt(-7),        // indexDelta=-3, settlement=3
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			updates: []types.Update{
				{
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(3)),
				},
			},
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(1_000_000), // 0.01 BTC
					FundingIndex: dtypes.NewInt(-10),       // indexDelta=-3, settlement=3
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId: uint32(0),
					Quantums: dtypes.NewIntFromBigInt(new(big.Int).Add(
						new(big.Int).SetUint64(math.MaxUint64),
						new(big.Int).SetInt64(1),
					)),
				},
			},
			expectedUpdatedPerpetualPositions: map[types.SubaccountId][]*types.PerpetualPosition{
				defaultSubaccountId: {
					{
						PerpetualId:  uint32(0),
						Quantums:     dtypes.NewInt(1_000_000),
						FundingIndex: dtypes.NewInt(-10),
					},
				},
			},
			expectedSubaccoundIdToFundingPayments: map[types.SubaccountId]map[uint32]dtypes.SerializableInt{
				defaultSubaccountId: {
					uint32(0): dtypes.NewInt(-3), // negated settlement
				},
			},
			expectedUpdatedAssetPositions: map[types.SubaccountId][]*types.AssetPosition{
				defaultSubaccountId: {
					{
						AssetId: uint32(0),
						Quantums: dtypes.NewIntFromBigInt(new(big.Int).Add(
							new(big.Int).SetUint64(math.MaxUint64),
							new(big.Int).SetInt64(1),
						)),
					},
				},
			},
			expectedAssetYieldIndex: constants.AssetYieldIndex_Zero,
			msgSenderEnabled:        true,
		},
		"new position size exceeds max uint64": {
			expectedQuoteBalance:     big.NewInt(0),
			expectedSuccess:          true,
			expectedSuccessPerUpdate: []types.UpdateResult{types.Success},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_SmallMarginRequirement,
			},
			updates: []types.Update{
				{
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(0),
							BigQuantumsDelta: big_testutil.MustFirst(new(big.Int).SetString("18446744073709551616", 10)),
						},
					},
				},
			},
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId: uint32(0),
					Quantums: dtypes.NewIntFromBigInt(
						big_testutil.MustFirst(new(big.Int).SetString("18446744073709551616", 10)), // 1 BTC
					),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			expectedUpdatedPerpetualPositions: map[types.SubaccountId][]*types.PerpetualPosition{
				defaultSubaccountId: {
					{
						PerpetualId: uint32(0),
						Quantums: dtypes.NewIntFromBigInt(
							big_testutil.MustFirst(new(big.Int).SetString("18446744073709551616", 10)), // 1 BTC
						),
						FundingIndex: dtypes.NewInt(0),
					},
				},
			},
			expectedAssetYieldIndex: constants.AssetYieldIndex_Zero,
			msgSenderEnabled:        true,
		},
		"existing position size + update exceeds max uint64": {
			expectedQuoteBalance:     big.NewInt(0),
			expectedSuccess:          true,
			expectedSuccessPerUpdate: []types.UpdateResult{types.Success},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_SmallMarginRequirement,
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewIntFromUint64(math.MaxUint64),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  0,
					Quantums: dtypes.NewInt(1),
				},
			},
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId: uint32(0),
					Quantums: dtypes.NewIntFromBigInt(
						new(big.Int).Add(
							new(big.Int).SetUint64(math.MaxUint64),
							new(big.Int).SetUint64(1),
						),
					),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			updates: []types.Update{
				{
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(1)),
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(0),
							BigQuantumsDelta: big.NewInt(1),
						},
					},
				},
			},
			expectedUpdatedPerpetualPositions: map[types.SubaccountId][]*types.PerpetualPosition{
				defaultSubaccountId: {
					{
						PerpetualId: uint32(0),
						Quantums: dtypes.NewIntFromBigInt(
							new(big.Int).Add(
								new(big.Int).SetUint64(math.MaxUint64),
								new(big.Int).SetUint64(1),
							),
						),
						FundingIndex: dtypes.NewInt(0),
					},
				},
			},
			expectedUpdatedAssetPositions: map[types.SubaccountId][]*types.AssetPosition{
				defaultSubaccountId: {
					{
						AssetId:  uint32(0),
						Quantums: dtypes.NewInt(1),
					},
				},
			},
			msgSenderEnabled: false,
		},
		"perpetual does not exist": {
			expectedQuoteBalance: big.NewInt(0),
			expectedErr:          perptypes.ErrPerpetualDoesNotExist,
			updates: []types.Update{
				{
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(999),
							BigQuantumsDelta: big.NewInt(1),
						},
					},
				},
			},
			expectedAssetYieldIndex: constants.AssetYieldIndex_Zero,
			msgSenderEnabled:        true,
		},
		"update ETH position; start with BTC and ETH positions; both BTC and ETH positions have unsettled funding": {
			assetPositions:           testutil.CreateTDaiAssetPosition(big.NewInt(100_000_000_000)), // $100,000
			expectedSuccess:          true,
			expectedSuccessPerUpdate: []types.UpdateResult{types.Success},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_NoMarginRequirement,
				constants.EthUsd_NoMarginRequirement,
			},
			newFundingIndices: []*big.Int{big.NewInt(-10), big.NewInt(-8)},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId: uint32(0),
					Quantums:    dtypes.NewInt(-100_000_000), // -1 BTC
					// indexDelta=-5
					FundingIndex: dtypes.NewInt(-5),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
				{
					PerpetualId: uint32(1),
					Quantums:    dtypes.NewInt(-2_000_000_000), // -2 ETH
					// indexDelta=-2
					FundingIndex: dtypes.NewInt(-6),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(-100_000_000), // -1 BTC
					FundingIndex: dtypes.NewInt(-10),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
				{
					PerpetualId:  uint32(1),
					Quantums:     dtypes.NewInt(-1_000_000_000), // -1 ETH
					FundingIndex: dtypes.NewInt(-8),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			expectedUpdatedPerpetualPositions: map[types.SubaccountId][]*types.PerpetualPosition{
				defaultSubaccountId: {
					{
						PerpetualId:  uint32(0),
						Quantums:     dtypes.NewInt(-100_000_000), // -1 BTC
						FundingIndex: dtypes.NewInt(-10),
					},
					{
						PerpetualId:  uint32(1),
						Quantums:     dtypes.NewInt(-1_000_000_000), // -1 ETH
						FundingIndex: dtypes.NewInt(-8),
					},
				},
			},
			expectedSubaccoundIdToFundingPayments: map[types.SubaccountId]map[uint32]dtypes.SerializableInt{
				defaultSubaccountId: {
					// indexDelta=-5, settlement=5*-100_000_000/1_000_000=-500
					uint32(0): dtypes.NewInt(500),
					// indexDelta=-2, settlement=2*-2_000_000_000/1_000_000=-4_000
					uint32(1): dtypes.NewInt(4_000),
				},
			},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId: uint32(0),
					// Original Asset Position - Funding Payments
					// = 100_000_000_000 - 4_000 - 500
					// = 99_999_995_500
					Quantums: dtypes.NewInt(99_999_995_500),
				},
			},
			updates: []types.Update{
				{
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(1),
							BigQuantumsDelta: big.NewInt(1_000_000_000), // 1 ETH
						},
					},
				},
			},
			expectedAssetYieldIndex: constants.AssetYieldIndex_Zero,
			msgSenderEnabled:        true,
		},
		"update ETH position; start with BTC and ETH positions; only ETH position has unsettled funding": {
			assetPositions:           testutil.CreateTDaiAssetPosition(big.NewInt(100_000_000_000)), // $100,000
			expectedSuccess:          true,
			expectedSuccessPerUpdate: []types.UpdateResult{types.Success},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_NoMarginRequirement,
				constants.EthUsd_NoMarginRequirement,
			},
			newFundingIndices: []*big.Int{big.NewInt(0), big.NewInt(-8)},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId: uint32(0),
					Quantums:    dtypes.NewInt(-100_000_000), // -1 BTC
					// indexDelta=0
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
				{
					PerpetualId: uint32(1),
					Quantums:    dtypes.NewInt(-2_000_000_000), // -2 ETH
					// indexDelta=-2
					FundingIndex: dtypes.NewInt(-6),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(-100_000_000), // -1 BTC
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
				{
					PerpetualId:  uint32(1),
					Quantums:     dtypes.NewInt(-1_000_000_000), // -1 ETH
					FundingIndex: dtypes.NewInt(-8),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			expectedUpdatedPerpetualPositions: map[types.SubaccountId][]*types.PerpetualPosition{
				defaultSubaccountId: {
					// Only ETH position is emitted here.
					{
						PerpetualId:  uint32(1),
						Quantums:     dtypes.NewInt(-1_000_000_000), // -1 ETH
						FundingIndex: dtypes.NewInt(-8),
					},
				},
			},
			expectedSubaccoundIdToFundingPayments: map[types.SubaccountId]map[uint32]dtypes.SerializableInt{
				defaultSubaccountId: {
					// indexDelta=-2, settlement=2*-2_000_000_000/1_000_000=-4_000
					uint32(1): dtypes.NewInt(4_000),
				},
			},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId: uint32(0),
					// Original Asset Position - Funding Payments
					// = 100_000_000_000 - 4_000
					// = 99_999_996_000
					Quantums: dtypes.NewInt(99_999_996_000),
				},
			},
			updates: []types.Update{
				{
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(1),
							BigQuantumsDelta: big.NewInt(1_000_000_000), // 1 ETH
						},
					},
				},
			},
			expectedAssetYieldIndex: constants.AssetYieldIndex_Zero,
			msgSenderEnabled:        true,
		},
		"update closes ETH position; start with BTC and ETH positions; both BTC and ETH positions have unsettled funding": {
			assetPositions:           testutil.CreateTDaiAssetPosition(big.NewInt(100_000_000_000)), // $100,000
			expectedSuccess:          true,
			expectedSuccessPerUpdate: []types.UpdateResult{types.Success},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_NoMarginRequirement,
				constants.EthUsd_NoMarginRequirement,
			},
			newFundingIndices: []*big.Int{big.NewInt(-10), big.NewInt(-8)},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId: uint32(0),
					Quantums:    dtypes.NewInt(-100_000_000), // -1 BTC
					// indexDelta=-5
					FundingIndex: dtypes.NewInt(-5),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
				{
					PerpetualId: uint32(1),
					Quantums:    dtypes.NewInt(-1_000_000_000), // -1 ETH
					// indexDelta=-2
					FundingIndex: dtypes.NewInt(-6),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(-100_000_000), // -1 BTC
					FundingIndex: dtypes.NewInt(-10),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			expectedUpdatedPerpetualPositions: map[types.SubaccountId][]*types.PerpetualPosition{
				defaultSubaccountId: {
					{
						PerpetualId:  uint32(0),
						Quantums:     dtypes.NewInt(-100_000_000), // -1 BTC
						FundingIndex: dtypes.NewInt(-10),
					},
					// Position closed update.
					{
						PerpetualId:  uint32(1),
						Quantums:     dtypes.NewInt(0),
						FundingIndex: dtypes.NewInt(0),
					},
				},
			},
			expectedSubaccoundIdToFundingPayments: map[types.SubaccountId]map[uint32]dtypes.SerializableInt{
				defaultSubaccountId: {
					// indexDelta=-5, settlement=5*-100_000_000/1_000_000=-500
					uint32(0): dtypes.NewInt(500),
					// indexDelta=-2, settlement=2*-1_000_000_000/1_000_000=-2_000
					uint32(1): dtypes.NewInt(2_000),
				},
			},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId: uint32(0),
					// Original Asset Position - Funding Payments
					// = 100_000_000_000 - 2_000 - 500
					// = 99_999_997_500
					Quantums: dtypes.NewInt(99_999_997_500),
				},
			},
			updates: []types.Update{
				{
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(1),
							BigQuantumsDelta: big.NewInt(1_000_000_000), // 1 ETH
						},
					},
				},
			},
			expectedAssetYieldIndex: constants.AssetYieldIndex_Zero,
			msgSenderEnabled:        true,
		},
		"update closes ETH position; start with ETH position; ETH position has no unsettled funding": {
			assetPositions:           testutil.CreateTDaiAssetPosition(big.NewInt(100_000_000_000)), // $100,000
			expectedSuccess:          true,
			expectedSuccessPerUpdate: []types.UpdateResult{types.Success},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_NoMarginRequirement,
				constants.EthUsd_NoMarginRequirement,
			},
			newFundingIndices: []*big.Int{big.NewInt(0)},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId: uint32(1),
					Quantums:    dtypes.NewInt(-1_000_000_000), // -1 ETH
					// indexDelta=0
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			expectedPerpetualPositions: []*types.PerpetualPosition{},
			expectedUpdatedPerpetualPositions: map[types.SubaccountId][]*types.PerpetualPosition{
				defaultSubaccountId: {
					// Position closed update.
					{
						PerpetualId:  uint32(1),
						Quantums:     dtypes.NewInt(0),
						FundingIndex: dtypes.NewInt(0),
					},
				},
			},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(100_000_000_000),
				},
			},
			updates: []types.Update{
				{
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(1),
							BigQuantumsDelta: big.NewInt(1_000_000_000), // 1 ETH
						},
					},
				},
			},
			expectedAssetYieldIndex: constants.AssetYieldIndex_Zero,
			msgSenderEnabled:        true,
		},
		"2 updates, 1 update involves not-updatable perp": {
			assetPositions: testutil.CreateTDaiAssetPosition(big.NewInt(1_000_000_000_000)),
			expectedErr:    types.ErrProductPositionNotUpdatable,
			perpetuals: []perptypes.Perpetual{
				*perptest.GeneratePerpetual(
					perptest.WithId(100),
					perptest.WithMarketId(100),
				),
				*perptest.GeneratePerpetual(
					perptest.WithId(101),
					perptest.WithMarketId(101),
				),
			},
			marketParamPrices: []pricestypes.MarketParamPrice{
				*pricestest.GenerateMarketParamPrice(pricestest.WithId(100)),
				*pricestest.GenerateMarketParamPrice(
					pricestest.WithId(101),
					pricestest.WithSpotPriceValue(0),
					pricestest.WithPnlPriceValue(0),
				),
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(100),
					Quantums:     dtypes.NewInt(1_000_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
				{
					PerpetualId:  uint32(101),
					Quantums:     dtypes.NewInt(1_000_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(100),
					Quantums:     dtypes.NewInt(1_000_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
				{
					PerpetualId:  uint32(101),
					Quantums:     dtypes.NewInt(1_000_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			expectedUpdatedPerpetualPositions: map[types.SubaccountId][]*types.PerpetualPosition{},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(1_000_000_000_000),
				},
			},
			updates: []types.Update{
				{
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(100),
							BigQuantumsDelta: big.NewInt(-1_000),
						},
						{
							PerpetualId:      uint32(101),
							BigQuantumsDelta: big.NewInt(1_000),
						},
					},
				},
			},
			expectedAssetYieldIndex: constants.AssetYieldIndex_Zero,
			msgSenderEnabled:        true,
		},
		"Isolated subaccounts - has update for both an isolated perpetual and non-isolated perpetual": {
			assetPositions:           testutil.CreateTDaiAssetPosition(big.NewInt(1_000_000_000_000)),
			expectedSuccess:          false,
			expectedSuccessPerUpdate: []types.UpdateResult{types.ViolatesIsolatedSubaccountConstraints},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_NoMarginRequirement,
				constants.IsoUsd_IsolatedMarket,
			},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(1_000_000_000_000),
				},
			},
			updates: []types.Update{
				{
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(0),
							BigQuantumsDelta: big.NewInt(-100_000_000), // -1 BTC
						},
						{
							PerpetualId:      uint32(3),
							BigQuantumsDelta: big.NewInt(1_000_000_000), // 1 ISO
						},
					},
				},
			},
			expectedAssetYieldIndex: constants.AssetYieldIndex_Zero,
			msgSenderEnabled:        true,
		},
		"Isolated subaccounts - has update for both 2 isolated perpetuals": {
			assetPositions:           testutil.CreateTDaiAssetPosition(big.NewInt(1_000_000_000_000)),
			expectedSuccess:          false,
			expectedSuccessPerUpdate: []types.UpdateResult{types.ViolatesIsolatedSubaccountConstraints},
			perpetuals: []perptypes.Perpetual{
				constants.IsoUsd_IsolatedMarket,
				constants.Iso2Usd_IsolatedMarket,
			},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(1_000_000_000_000),
				},
			},
			updates: []types.Update{
				{
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(3),
							BigQuantumsDelta: big.NewInt(-1_000_000_000), // 1 ISO
						},
						{
							PerpetualId:      uint32(4),
							BigQuantumsDelta: big.NewInt(10_000_000), // 1 ISO2
						},
					},
				},
			},
			expectedAssetYieldIndex: constants.AssetYieldIndex_Zero,
			msgSenderEnabled:        true,
		},
		"Isolated subaccounts - subaccount with isolated perpetual position has update for non-isolated perpetual": {
			assetPositions:           testutil.CreateTDaiAssetPosition(big.NewInt(1_000_000_000_000)),
			expectedSuccess:          false,
			expectedSuccessPerUpdate: []types.UpdateResult{types.ViolatesIsolatedSubaccountConstraints},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_NoMarginRequirement,
				constants.IsoUsd_IsolatedMarket,
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(3),
					Quantums:     dtypes.NewInt(1_000_000_000), // 1 ISO
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(3),
					Quantums:     dtypes.NewInt(1_000_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(1_000_000_000_000),
				},
			},
			updates: []types.Update{
				{
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(0),
							BigQuantumsDelta: big.NewInt(-100_000_000), // -1 BTC
						},
					},
				},
			},
			expectedAssetYieldIndex: constants.AssetYieldIndex_Zero,
			msgSenderEnabled:        true,
		},
		"Isolated subaccounts - subaccount with isolated perpetual position has update for another isolated perpetual": {
			assetPositions:           testutil.CreateTDaiAssetPosition(big.NewInt(1_000_000_000_000)),
			expectedSuccess:          false,
			expectedSuccessPerUpdate: []types.UpdateResult{types.ViolatesIsolatedSubaccountConstraints},
			perpetuals: []perptypes.Perpetual{
				constants.IsoUsd_IsolatedMarket,
				constants.Iso2Usd_IsolatedMarket,
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(3),
					Quantums:     dtypes.NewInt(1_000_000_000), // 1 ISO
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(3),
					Quantums:     dtypes.NewInt(1_000_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(1_000_000_000_000),
				},
			},
			updates: []types.Update{
				{
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(4),
							BigQuantumsDelta: big.NewInt(-10_000_000), // -1 ISO2
						},
					},
				},
			},
			expectedAssetYieldIndex: constants.AssetYieldIndex_Zero,
			msgSenderEnabled:        true,
		},
		"Isolated subaccounts - subaccount with non-isolated perpetual position has update for isolated perpetual": {
			assetPositions:           testutil.CreateTDaiAssetPosition(big.NewInt(1_000_000_000_000)),
			expectedSuccess:          false,
			expectedSuccessPerUpdate: []types.UpdateResult{types.ViolatesIsolatedSubaccountConstraints},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_NoMarginRequirement,
				constants.IsoUsd_IsolatedMarket,
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000), // 1 BTC
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(1_000_000_000_000),
				},
			},
			updates: []types.Update{
				{
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(3),
							BigQuantumsDelta: big.NewInt(-1_000_000_000), // -1 ISO
						},
					},
				},
			},
			expectedAssetYieldIndex: constants.AssetYieldIndex_Zero,
			msgSenderEnabled:        true,
		},
		`Isolated - subaccounts - empty subaccount has update to open position for isolated perpetual,
		collateral is moved from cross-perpetual collateral pool to isolated perpetual collateral pool`: {
			assetPositions: testutil.CreateTDaiAssetPosition(big.NewInt(1_000_000_000_000)),
			collateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 1_500_000_000_000, // $1,500,000 TDai
			},
			expectedCollateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 500_000_000_000, // $500,000 TDai
				authtypes.NewModuleAddress(
					types.ModuleName + ":" + lib.UintToString(constants.PerpetualPosition_OneISOLong.PerpetualId),
				).String(): 1_000_000_000_000, // $1,000,000 TDai
			},
			expectedSuccess:          true,
			expectedSuccessPerUpdate: []types.UpdateResult{types.Success},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_NoMarginRequirement,
				constants.IsoUsd_IsolatedMarket,
			},
			perpetualPositions: []*types.PerpetualPosition{},
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(3),
					Quantums:     dtypes.NewInt(1_000_000_000), // 1 ISO
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			expectedUpdatedPerpetualPositions: map[types.SubaccountId][]*types.PerpetualPosition{
				defaultSubaccountId: {
					{
						PerpetualId:  uint32(3),
						Quantums:     dtypes.NewInt(1_000_000_000), // 1 ISO
						FundingIndex: dtypes.NewInt(0),
					},
				},
			},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(999_900_000_000),
				},
			},
			expectedUpdatedAssetPositions: map[types.SubaccountId][]*types.AssetPosition{
				defaultSubaccountId: {
					{
						AssetId:  uint32(0),
						Quantums: dtypes.NewInt(999_900_000_000),
					},
				},
			},
			updates: []types.Update{
				{
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(-100_000_000)), // -$100
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(3),
							BigQuantumsDelta: big.NewInt(1_000_000_000), // 1 ISO
						},
					},
				},
			},
			expectedAssetYieldIndex: constants.AssetYieldIndex_Zero,
			msgSenderEnabled:        true,
		},
		`Isolated - subaccounts - subaccount has update to close position for isolated perpetual,
		collateral is moved from isolated perpetual collateral pool to cross perpetual collateral pool`: {
			assetPositions: testutil.CreateTDaiAssetPosition(big.NewInt(999_900_000_000)), // $999,900 TDai
			collateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 2_000_000_000_000, // $500,000 TDai
				authtypes.NewModuleAddress(
					types.ModuleName + ":" + lib.UintToString(constants.PerpetualPosition_OneISOLong.PerpetualId),
				).String(): 1_500_000_000_000, // $1,500,000 TDai
			},
			expectedCollateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 3_000_000_000_000, // $3,000,000 TDai
				authtypes.NewModuleAddress(
					types.ModuleName + ":" + lib.UintToString(constants.PerpetualPosition_OneISOLong.PerpetualId),
				).String(): 500_000_000_000, // $500,000 TDai
			},
			expectedSuccess:          true,
			expectedSuccessPerUpdate: []types.UpdateResult{types.Success},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_NoMarginRequirement,
				constants.IsoUsd_IsolatedMarket,
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(3),
					Quantums:     dtypes.NewInt(1_000_000_000), // 1 ISO
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			expectedPerpetualPositions: []*types.PerpetualPosition{},
			expectedUpdatedPerpetualPositions: map[types.SubaccountId][]*types.PerpetualPosition{
				defaultSubaccountId: {
					{
						PerpetualId:  uint32(3),
						Quantums:     dtypes.NewInt(0),
						FundingIndex: dtypes.NewInt(0),
					},
				},
			},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(1_000_000_000_000),
				},
			},
			expectedUpdatedAssetPositions: map[types.SubaccountId][]*types.AssetPosition{
				defaultSubaccountId: {
					{
						AssetId:  uint32(0),
						Quantums: dtypes.NewInt(1_000_000_000_000),
					},
				},
			},
			updates: []types.Update{
				{
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(100_000_000)), // $100
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(3),
							BigQuantumsDelta: big.NewInt(-1_000_000_000), // -1 ISO
						},
					},
				},
			},
			expectedAssetYieldIndex: constants.AssetYieldIndex_Zero,
			msgSenderEnabled:        true,
		},
		`Isolated subaccounts - empty subaccount has update to open position for isolated perpetual, 
		errors out when collateral pool for cross perpetuals has no funds`: {
			assetPositions:           testutil.CreateTDaiAssetPosition(big.NewInt(1_000_000_000_000)),
			expectedSuccess:          false,
			expectedSuccessPerUpdate: []types.UpdateResult{},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_NoMarginRequirement,
				constants.IsoUsd_IsolatedMarket,
			},
			perpetualPositions:         []*types.PerpetualPosition{},
			expectedPerpetualPositions: []*types.PerpetualPosition{},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(1_000_000_000_000),
				},
			},
			updates: []types.Update{
				{
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(-100_000_000)), // -$100
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(3),
							BigQuantumsDelta: big.NewInt(1_000_000_000), // 1 ISO
						},
					},
				},
			},
			expectedAssetYieldIndex: constants.AssetYieldIndex_Zero,
			expectedErr:             sdkerrors.ErrInsufficientFunds,
			msgSenderEnabled:        true,
		},
		`Isolated subaccounts - isolated subaccount has update to close position for isolated perpetual, 
		errors out when collateral pool for isolated perpetual has no funds`: {
			assetPositions:           testutil.CreateTDaiAssetPosition(big.NewInt(1_000_000_000_000)),
			expectedSuccess:          false,
			expectedSuccessPerUpdate: []types.UpdateResult{},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_NoMarginRequirement,
				constants.IsoUsd_IsolatedMarket,
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(3),
					Quantums:     dtypes.NewInt(1_000_000_000), // 1 ISO
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(3),
					Quantums:     dtypes.NewInt(1_000_000_000), // 1 ISO
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(1_000_000_000_000),
				},
			},
			updates: []types.Update{
				{
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(100_000_000)), // $100
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(3),
							BigQuantumsDelta: big.NewInt(-1_000_000_000), // -1 ISO
						},
					},
				},
			},
			expectedAssetYieldIndex: constants.AssetYieldIndex_Zero,
			expectedErr:             sdkerrors.ErrInsufficientFunds,
			msgSenderEnabled:        true,
		},
		"Match updates increase OI: 0 -> 0.9, 0 -> -0.9": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance_OpenInterest1,
			},
			updates: []types.Update{
				{
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(0),
							BigQuantumsDelta: big.NewInt(9_000_000_000), // 90 BTC
						},
					},
					AssetUpdates: []types.AssetUpdate{
						{
							AssetId:          uint32(0),
							BigQuantumsDelta: big.NewInt(-4_500_000_000_000), // -4,500,000 TDai
						},
					},
					SubaccountId: constants.Bob_Num0,
				},
				{
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(0),
							BigQuantumsDelta: big.NewInt(-9_000_000_000), // 9 BTC
						},
					},
					AssetUpdates: []types.AssetUpdate{
						{
							AssetId:          uint32(0),
							BigQuantumsDelta: big.NewInt(4_500_000_000_000), // 4,500,000 TDai
						},
					},
				},
			},
			assetPositions: testutil.CreateTDaiAssetPosition(big.NewInt(900_000_000_000)), // 900_000 TDai
			additionalTestSubaccounts: []types.Subaccount{
				{
					Id: &constants.Bob_Num0,
					AssetPositions: testutil.CreateTDaiAssetPosition(big.NewInt(
						900_000_000_000,
					)), // 900_000 TDai
					AssetYieldIndex: "1/1",
				},
			},
			updateType: types.Match,
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(5_400_000_000_000),
				},
			},
			expectedUpdatedAssetPositions: map[types.SubaccountId][]*types.AssetPosition{
				defaultSubaccountId: {
					{
						AssetId:  uint32(0),
						Quantums: dtypes.NewInt(5_400_000_000_000),
					},
				},
				constants.Bob_Num0: {
					{
						AssetId:  uint32(0),
						Quantums: dtypes.NewInt(-3_600_000_000_000),
					},
				},
			},
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(-9_000_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			expectedUpdatedPerpetualPositions: map[types.SubaccountId][]*types.PerpetualPosition{
				defaultSubaccountId: {
					{
						PerpetualId:  uint32(0),
						Quantums:     dtypes.NewInt(-9_000_000_000),
						FundingIndex: dtypes.NewInt(0),
					},
				},
				constants.Bob_Num0: {
					{
						PerpetualId:  uint32(0),
						Quantums:     dtypes.NewInt(9_000_000_000),
						FundingIndex: dtypes.NewInt(0),
					},
				},
			},
			expectedAssetYieldIndex:  constants.AssetYieldIndex_Zero,
			expectedSuccess:          true,
			expectedSuccessPerUpdate: []types.UpdateResult{types.Success, types.Success},
			expectedOpenInterest: map[uint32]*big.Int{
				0: big.NewInt(9_100_000_000), // 1 + 90 = 91 BTC
			},
			msgSenderEnabled: true,
		},
		"Match updates decreases OI: 1 -> 0.1, -2 -> -1.1": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance_OpenInterest2,
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId: uint32(0),
					Quantums:    dtypes.NewInt(100_000_000), // 1 BTC
					YieldIndex:  big.NewRat(0, 1).String(),
				},
			},
			assetPositions: testutil.CreateTDaiAssetPosition(big.NewInt(-40_000_000_000)), // -40_000 TDai
			updates: []types.Update{
				{
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(0),
							BigQuantumsDelta: big.NewInt(90_000_000), // 0.9 BTC
						},
					},
					AssetUpdates: []types.AssetUpdate{
						{
							AssetId:          uint32(0),
							BigQuantumsDelta: big.NewInt(-45_000_000_000), // -45,000 TDai
						},
					},
					SubaccountId: constants.Bob_Num0,
				},
				{
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(0),
							BigQuantumsDelta: big.NewInt(-90_000_000), // -0.9 BTC
						},
					},
					AssetUpdates: []types.AssetUpdate{
						{
							AssetId:          uint32(0),
							BigQuantumsDelta: big.NewInt(45_000_000_000), // 45,000 TDai
						},
					},
				},
			},
			additionalTestSubaccounts: []types.Subaccount{
				{
					Id: &constants.Bob_Num0,
					AssetPositions: testutil.CreateTDaiAssetPosition(big.NewInt(
						120_000_000_000,
					)), // 120_000 TDai
					PerpetualPositions: []*types.PerpetualPosition{
						{
							PerpetualId: uint32(0),
							Quantums:    dtypes.NewInt(-200_000_000), // -2 BTC
						},
					},
					AssetYieldIndex: "1/1",
				},
			},
			updateType: types.Match,
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(5_000_000_000), // 5_000 TDai
				},
			},
			expectedUpdatedAssetPositions: map[types.SubaccountId][]*types.AssetPosition{
				defaultSubaccountId: {
					{
						AssetId:  uint32(0),
						Quantums: dtypes.NewInt(5_000_000_000), // 5_000 TDai
					},
				},
				constants.Bob_Num0: {
					{
						AssetId:  uint32(0),
						Quantums: dtypes.NewInt(75_000_000_000), // 75_000 TDai
					},
				},
			},
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(10_000_000), // 0.1 BTC
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			expectedUpdatedPerpetualPositions: map[types.SubaccountId][]*types.PerpetualPosition{
				defaultSubaccountId: {
					{
						PerpetualId:  uint32(0),
						Quantums:     dtypes.NewInt(10_000_000), // 0.1 BTC
						FundingIndex: dtypes.NewInt(0),
					},
				},
				constants.Bob_Num0: {
					{
						PerpetualId:  uint32(0),
						Quantums:     dtypes.NewInt(-110_000_000), // -1.1 BTC
						FundingIndex: dtypes.NewInt(0),
					},
				},
			},
			expectedSuccess:          true,
			expectedSuccessPerUpdate: []types.UpdateResult{types.Success, types.Success},
			expectedOpenInterest: map[uint32]*big.Int{
				0: big.NewInt(110_000_000), // 2 - 0.9 = 1.1 BTC
			},
			expectedAssetYieldIndex: constants.AssetYieldIndex_Zero,
			msgSenderEnabled:        true,
		},
		"Match updates does not change OI: 1 -> 0.1, 0.1 -> 1": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance_OpenInterest1,
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId: uint32(0),
					Quantums:    dtypes.NewInt(100_000_000), // 1 BTC
					YieldIndex:  big.NewRat(0, 1).String(),
				},
			},
			assetPositions: testutil.CreateTDaiAssetPosition(big.NewInt(-40_000_000_000)), // -40_000 TDai
			updates: []types.Update{
				{
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(0),
							BigQuantumsDelta: big.NewInt(90_000_000), // 0.9 BTC
						},
					},
					AssetUpdates: []types.AssetUpdate{
						{
							AssetId:          uint32(0),
							BigQuantumsDelta: big.NewInt(-45_000_000_000), // -45,000 TDai
						},
					},
					SubaccountId: constants.Bob_Num0,
				},
				{
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(0),
							BigQuantumsDelta: big.NewInt(-90_000_000), // -0.9 BTC
						},
					},
					AssetUpdates: []types.AssetUpdate{
						{
							AssetId:          uint32(0),
							BigQuantumsDelta: big.NewInt(45_000_000_000), // 45,000 TDai
						},
					},
				},
			},
			additionalTestSubaccounts: []types.Subaccount{
				{
					Id:             &constants.Bob_Num0,
					AssetPositions: testutil.CreateTDaiAssetPosition(big.NewInt(5_000_000_000)), // 5000 TDai
					PerpetualPositions: []*types.PerpetualPosition{
						{
							PerpetualId: uint32(0),
							Quantums:    dtypes.NewInt(10_000_000), // 0.1 BTC
						},
					},
					AssetYieldIndex: "1/1",
				},
			},
			updateType: types.Match,
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(5_000_000_000), // 5_000 TDai
				},
			},
			expectedUpdatedAssetPositions: map[types.SubaccountId][]*types.AssetPosition{
				defaultSubaccountId: {
					{
						AssetId:  uint32(0),
						Quantums: dtypes.NewInt(5_000_000_000), // 5_000 TDai
					},
				},
				constants.Bob_Num0: {
					{
						AssetId:  uint32(0),
						Quantums: dtypes.NewInt(-40_000_000_000), // -40_000 TDai
					},
				},
			},
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(10_000_000), // 0.1 BTC
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			expectedUpdatedPerpetualPositions: map[types.SubaccountId][]*types.PerpetualPosition{
				defaultSubaccountId: {
					{
						PerpetualId:  uint32(0),
						Quantums:     dtypes.NewInt(10_000_000), // 0.1 BTC
						FundingIndex: dtypes.NewInt(0),
					},
				},
				constants.Bob_Num0: {
					{
						PerpetualId:  uint32(0),
						Quantums:     dtypes.NewInt(100_000_000), // 1 BTC
						FundingIndex: dtypes.NewInt(0),
					},
				},
			},
			expectedSuccess:          true,
			expectedSuccessPerUpdate: []types.UpdateResult{types.Success, types.Success},
			expectedOpenInterest: map[uint32]*big.Int{
				0: big.NewInt(100_000_000), // 1 BTC
			},
			expectedAssetYieldIndex: constants.AssetYieldIndex_Zero,
			msgSenderEnabled:        true,
		},
		"Successfully claims yield for tDai asset": {
			assetPositions:            testutil.CreateTDaiAssetPosition(big.NewInt(100_000_000_000)),
			subaccountAssetYieldIndex: constants.AssetYieldIndex_Zero,
			globalAssetYieldIndex:     big.NewRat(2, 1),
			fundsInTDaiPool:           big.NewInt(200_000_000_000),
			collateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 100_000_000_000,
			},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_NoMarginRequirement,
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000), // 1 BTC
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			expectedSuccess:          true,
			expectedSuccessPerUpdate: []types.UpdateResult{types.Success},
			expectedAssetYieldIndex:  big.NewRat(2, 1).String(),
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(150_000_000), // 1.5 BTC
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			expectedUpdatedPerpetualPositions: map[types.SubaccountId][]*types.PerpetualPosition{
				defaultSubaccountId: {
					{
						PerpetualId:  uint32(0),
						Quantums:     dtypes.NewInt(150_000_000), // 1.5 BTC
						FundingIndex: dtypes.NewInt(0),
					},
				},
			},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(175_000_000_000),
				},
			},
			expectedUpdatedAssetPositions: map[types.SubaccountId][]*types.AssetPosition{
				defaultSubaccountId: {
					{
						AssetId:  uint32(0),
						Quantums: dtypes.NewInt(175_000_000_000),
					},
				},
			},
			updates: []types.Update{
				{
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(-25_000_000_000)), // -$25,000
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(0),
							BigQuantumsDelta: big.NewInt(50_000_000), // .5 BTC
						},
					},
				},
			},
			expectedTDaiYieldPoolBalance: big.NewInt(100_000_000_000),
			expectedCollateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 200_000_000_000,
			},
			msgSenderEnabled: true,
		},
		"Successfully claims yield for one perp": {
			assetPositions:            testutil.CreateTDaiAssetPosition(big.NewInt(100_000_000_000)), // $100,000
			subaccountAssetYieldIndex: constants.AssetYieldIndex_Zero,
			globalAssetYieldIndex:     big.NewRat(1, 1),
			fundsInTDaiPool:           big.NewInt(200_000_000_000),
			collateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 100_000_000_000,
			},
			perpetuals: []perptypes.Perpetual{
				{
					Params:       constants.BtcUsd_NoMarginRequirement.Params,
					FundingIndex: constants.BtcUsd_NoMarginRequirement.FundingIndex,
					OpenInterest: constants.BtcUsd_NoMarginRequirement.OpenInterest,
					YieldIndex:   big.NewRat(1, 1).String(),
				},
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(1_000_000_000), // 10 BTC
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			expectedSuccess:          true,
			expectedSuccessPerUpdate: []types.UpdateResult{types.Success},
			expectedAssetYieldIndex:  big.NewRat(1, 1).String(),
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(1_000_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(1, 1).String(),
				},
			},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(101_000_000_001),
				},
			},
			expectedUpdatedAssetPositions: map[types.SubaccountId][]*types.AssetPosition{
				defaultSubaccountId: {
					{
						AssetId:  uint32(0),
						Quantums: dtypes.NewInt(101_000_000_001),
					},
				},
			},
			updates: []types.Update{
				{
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(1)),
				},
			},
			expectedTDaiYieldPoolBalance: big.NewInt(199_000_000_000),
			expectedCollateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 101_000_000_000,
			},
			msgSenderEnabled: true,
		},
		"Successfully claims yield for tDai asset and one perp": {
			assetPositions:            testutil.CreateTDaiAssetPosition(big.NewInt(100_000_000_000)),
			subaccountAssetYieldIndex: big.NewRat(1, 1).String(),
			globalAssetYieldIndex:     big.NewRat(3, 2),
			fundsInTDaiPool:           big.NewInt(200_000_000_000),
			collateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 100_000_000_000,
			},
			perpetuals: []perptypes.Perpetual{
				{
					Params:       constants.BtcUsd_NoMarginRequirement.Params,
					FundingIndex: constants.BtcUsd_NoMarginRequirement.FundingIndex,
					OpenInterest: constants.BtcUsd_NoMarginRequirement.OpenInterest,
					YieldIndex:   big.NewRat(1, 1).String(),
				},
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000), // 1 BTC
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			expectedSuccess:          true,
			expectedSuccessPerUpdate: []types.UpdateResult{types.Success},
			expectedAssetYieldIndex:  big.NewRat(3, 2).String(),
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(150_000_000), // 1.5 BTC
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(1, 1).String(),
				},
			},
			expectedUpdatedPerpetualPositions: map[types.SubaccountId][]*types.PerpetualPosition{
				defaultSubaccountId: {
					{
						PerpetualId:  uint32(0),
						Quantums:     dtypes.NewInt(150_000_000), // 1.5 BTC
						FundingIndex: dtypes.NewInt(0),
						YieldIndex:   big.NewRat(1, 1).String(),
					},
				},
			},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(125_100_000_000),
				},
			},
			expectedUpdatedAssetPositions: map[types.SubaccountId][]*types.AssetPosition{
				defaultSubaccountId: {
					{
						AssetId:  uint32(0),
						Quantums: dtypes.NewInt(125_100_000_000),
					},
				},
			},
			updates: []types.Update{
				{
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(-25_000_000_000)), // -$25,000
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(0),
							BigQuantumsDelta: big.NewInt(50_000_000), // .5 BTC
						},
					},
				},
			},
			expectedTDaiYieldPoolBalance: big.NewInt(149_900_000_000),
			expectedCollateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 150_100_000_000,
			},
			msgSenderEnabled: true,
		},
		"Successfully claims yield for tDai deposit": {
			assetPositions: testutil.CreateTDaiAssetPosition(big.NewInt(100_000_000_000)),
			collateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 100_000_000_000,
			},
			subaccountAssetYieldIndex: constants.AssetYieldIndex_Zero,
			globalAssetYieldIndex:     big.NewRat(2, 1),
			fundsInTDaiPool:           big.NewInt(200_000_000_000),
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_NoMarginRequirement,
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000), // 1 BTC
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			expectedSuccess:          true,
			expectedSuccessPerUpdate: []types.UpdateResult{types.Success},
			expectedAssetYieldIndex:  big.NewRat(2, 1).String(),
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(210_000_000_000),
				},
			},
			expectedUpdatedAssetPositions: map[types.SubaccountId][]*types.AssetPosition{
				defaultSubaccountId: {
					{
						AssetId:  uint32(0),
						Quantums: dtypes.NewInt(210_000_000_000),
					},
				},
			},
			updates: []types.Update{
				{
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(10_000_000_000)),
				},
			},
			expectedTDaiYieldPoolBalance: big.NewInt(100_000_000_000),
			expectedCollateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 200_000_000_000,
			},
			msgSenderEnabled: true,
		},
		"Successfully claims yield for tDai withdrawal": {
			assetPositions: testutil.CreateTDaiAssetPosition(big.NewInt(100_000_000_000)),
			collateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 100_000_000_000,
			},
			subaccountAssetYieldIndex: constants.AssetYieldIndex_Zero,
			globalAssetYieldIndex:     big.NewRat(2, 1),
			fundsInTDaiPool:           big.NewInt(200_000_000_000),
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_NoMarginRequirement,
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000), // 1 BTC
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			expectedSuccess:          true,
			expectedSuccessPerUpdate: []types.UpdateResult{types.Success},
			expectedAssetYieldIndex:  big.NewRat(2, 1).String(),
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(190_000_000_000),
				},
			},
			expectedUpdatedAssetPositions: map[types.SubaccountId][]*types.AssetPosition{
				defaultSubaccountId: {
					{
						AssetId:  uint32(0),
						Quantums: dtypes.NewInt(190_000_000_000),
					},
				},
			},
			updates: []types.Update{
				{
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(-10_000_000_000)),
				},
			},
			expectedTDaiYieldPoolBalance: big.NewInt(100_000_000_000),
			expectedCollateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 200_000_000_000,
			},
			msgSenderEnabled: true,
		},
		"Successfully claims yield and correctly sets AssetYieldIndex when only perp position is open": {
			subaccountAssetYieldIndex: constants.AssetYieldIndex_Zero,
			globalAssetYieldIndex:     big.NewRat(1, 1),
			fundsInTDaiPool:           big.NewInt(222_000_000_000),
			perpetuals: []perptypes.Perpetual{
				{
					Params:       constants.BtcUsd_NoMarginRequirement.Params,
					FundingIndex: constants.BtcUsd_NoMarginRequirement.FundingIndex,
					OpenInterest: constants.BtcUsd_NoMarginRequirement.OpenInterest,
					YieldIndex:   big.NewRat(3, 2).String(),
				},
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000), // 1 BTC
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(3, 4).String(),
				},
			},
			expectedSuccess:          true,
			expectedSuccessPerUpdate: []types.UpdateResult{types.Success},
			expectedAssetYieldIndex:  big.NewRat(1, 1).String(),
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_001),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(3, 2).String(),
				},
			},
			expectedUpdatedPerpetualPositions: map[types.SubaccountId][]*types.PerpetualPosition{
				defaultSubaccountId: {
					{
						PerpetualId:  uint32(0),
						Quantums:     dtypes.NewInt(100_000_001),
						FundingIndex: dtypes.NewInt(0),
						YieldIndex:   big.NewRat(3, 2).String(),
					},
				},
			},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(75_000_000),
				},
			},
			updates: []types.Update{
				{
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(0),
							BigQuantumsDelta: big.NewInt(1),
						},
					},
				},
			},
			expectedTDaiYieldPoolBalance: big.NewInt(221_925_000_000),
			expectedCollateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 75_000_000,
			},
			msgSenderEnabled: true,
		},
		"Successfully claims yield and correctly sets AssetYieldIndex when only asset position open": {
			assetPositions: testutil.CreateTDaiAssetPosition(big.NewInt(50_000_000_000)),
			collateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 50_000_000_000,
			},
			subaccountAssetYieldIndex: constants.AssetYieldIndex_Zero,
			globalAssetYieldIndex:     big.NewRat(15, 7),
			fundsInTDaiPool:           big.NewInt(123_000_000_000),
			expectedSuccess:           true,
			expectedSuccessPerUpdate:  []types.UpdateResult{types.Success},
			expectedAssetYieldIndex:   big.NewRat(15, 7).String(),
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(107_142_857_142),
				},
			},
			expectedUpdatedAssetPositions: map[types.SubaccountId][]*types.AssetPosition{
				defaultSubaccountId: {
					{
						AssetId:  uint32(0),
						Quantums: dtypes.NewInt(107_142_857_142),
					},
				},
			},
			updates: []types.Update{
				{
					AssetUpdates: []types.AssetUpdate{
						{
							AssetId:          uint32(0),
							BigQuantumsDelta: big.NewInt(0),
						},
					},
				},
			},
			expectedTDaiYieldPoolBalance: big.NewInt(65_857_142_858),
			expectedCollateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 107_142_857_142,
			},
			msgSenderEnabled: true,
		},
		"Claims yield when multiple perp positions are open and no tDai position open": {
			subaccountAssetYieldIndex: constants.AssetYieldIndex_Zero,
			globalAssetYieldIndex:     big.NewRat(12, 11),
			fundsInTDaiPool:           big.NewInt(222_000_000_000),
			perpetuals: []perptypes.Perpetual{
				{
					Params:       constants.BtcUsd_NoMarginRequirement.Params,
					FundingIndex: constants.BtcUsd_NoMarginRequirement.FundingIndex,
					OpenInterest: constants.BtcUsd_NoMarginRequirement.OpenInterest,
					YieldIndex:   big.NewRat(11, 3).String(),
				},
				{
					Params:       constants.EthUsd_NoMarginRequirement.Params,
					FundingIndex: constants.EthUsd_NoMarginRequirement.FundingIndex,
					OpenInterest: constants.EthUsd_NoMarginRequirement.OpenInterest,
					YieldIndex:   big.NewRat(4, 3).String(),
				},
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(1, 2).String(),
				},
				{
					PerpetualId:  uint32(1),
					Quantums:     dtypes.NewInt(-2_000_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(11, 9).String(),
				},
			},
			expectedSuccess:          true,
			expectedSuccessPerUpdate: []types.UpdateResult{types.Success},
			expectedAssetYieldIndex:  big.NewRat(12, 11).String(),
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(11, 3).String(),
				},
				{
					PerpetualId:  uint32(1),
					Quantums:     dtypes.NewInt(-2_000_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(4, 3).String(),
				},
			},
			expectedUpdatedPerpetualPositions: map[types.SubaccountId][]*types.PerpetualPosition{
				defaultSubaccountId: {
					{
						PerpetualId:  uint32(0),
						Quantums:     dtypes.NewInt(100_000_000),
						FundingIndex: dtypes.NewInt(0),
						YieldIndex:   big.NewRat(11, 3).String(),
					},
				},
			},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(94444443),
				},
			},
			updates: []types.Update{
				{
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(0),
							BigQuantumsDelta: big.NewInt(0),
						},
					},
				},
			},
			expectedTDaiYieldPoolBalance: big.NewInt(221_905_555_557),
			expectedCollateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 94_444_443,
			},
			msgSenderEnabled: true,
		},
		"Successfully claims yield when multiple perp positions are open and tDai position open": {
			assetPositions:            testutil.CreateTDaiAssetPosition(big.NewInt(50_000_000_000)),
			subaccountAssetYieldIndex: big.NewRat(13, 11).String(),
			globalAssetYieldIndex:     big.NewRat(26, 11),
			fundsInTDaiPool:           big.NewInt(222_000_000_000),
			collateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 50_000_000_000,
			},
			perpetuals: []perptypes.Perpetual{
				{
					Params:       constants.BtcUsd_NoMarginRequirement.Params,
					FundingIndex: constants.BtcUsd_NoMarginRequirement.FundingIndex,
					OpenInterest: constants.BtcUsd_NoMarginRequirement.OpenInterest,
					YieldIndex:   big.NewRat(11, 3).String(),
				},
				{
					Params:       constants.EthUsd_NoMarginRequirement.Params,
					FundingIndex: constants.EthUsd_NoMarginRequirement.FundingIndex,
					OpenInterest: constants.EthUsd_NoMarginRequirement.OpenInterest,
					YieldIndex:   big.NewRat(4, 3).String(),
				},
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(1, 2).String(),
				},
				{
					PerpetualId:  uint32(1),
					Quantums:     dtypes.NewInt(-2_000_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(11, 9).String(),
				},
			},
			expectedSuccess:          true,
			expectedSuccessPerUpdate: []types.UpdateResult{types.Success},
			expectedAssetYieldIndex:  big.NewRat(26, 11).String(),
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(150_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(11, 3).String(),
				},
				{
					PerpetualId:  uint32(1),
					Quantums:     dtypes.NewInt(-5_000_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(4, 3).String(),
				},
			},
			expectedUpdatedPerpetualPositions: map[types.SubaccountId][]*types.PerpetualPosition{
				defaultSubaccountId: {
					{
						PerpetualId:  uint32(0),
						Quantums:     dtypes.NewInt(150_000_000),
						FundingIndex: dtypes.NewInt(0),
						YieldIndex:   big.NewRat(11, 3).String(),
					},
					{
						PerpetualId:  uint32(1),
						Quantums:     dtypes.NewInt(-5_000_000_000),
						FundingIndex: dtypes.NewInt(0),
						YieldIndex:   big.NewRat(4, 3).String(),
					},
				},
			},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(75_094_444_443), // Total Yield: 50_094_444_443
				},
			},
			expectedUpdatedAssetPositions: map[types.SubaccountId][]*types.AssetPosition{
				defaultSubaccountId: {
					{
						AssetId:  uint32(0),
						Quantums: dtypes.NewInt(75_094_444_443),
					},
				},
			},
			updates: []types.Update{
				{
					AssetUpdates: []types.AssetUpdate{
						{
							AssetId:          uint32(0),
							BigQuantumsDelta: big.NewInt(-25_000_000_000),
						},
					},
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(0),
							BigQuantumsDelta: big.NewInt(50_000_000),
						},
						{
							PerpetualId:      uint32(1),
							BigQuantumsDelta: big.NewInt(-3_000_000_000),
						},
					},
				},
			},
			expectedTDaiYieldPoolBalance: big.NewInt(171_905_555_557),
			expectedCollateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 100_094_444_443,
			},
			msgSenderEnabled: true,
		},
		"Successfully claims all yield in tDaiPool for tDai asset": {
			assetPositions:            testutil.CreateTDaiAssetPosition(big.NewInt(100_000_000_000)),
			subaccountAssetYieldIndex: constants.AssetYieldIndex_Zero,
			globalAssetYieldIndex:     big.NewRat(3, 1),
			fundsInTDaiPool:           big.NewInt(200_000_000_000),
			collateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 100_000_000_000,
			},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_NoMarginRequirement,
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000), // 1 BTC
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			expectedSuccess:          true,
			expectedSuccessPerUpdate: []types.UpdateResult{types.Success},
			expectedAssetYieldIndex:  big.NewRat(3, 1).String(),
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(150_000_000), // 1.5 BTC
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			expectedUpdatedPerpetualPositions: map[types.SubaccountId][]*types.PerpetualPosition{
				defaultSubaccountId: {
					{
						PerpetualId:  uint32(0),
						Quantums:     dtypes.NewInt(150_000_000), // 1.5 BTC
						FundingIndex: dtypes.NewInt(0),
					},
				},
			},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(275_000_000_000),
				},
			},
			expectedUpdatedAssetPositions: map[types.SubaccountId][]*types.AssetPosition{
				defaultSubaccountId: {
					{
						AssetId:  uint32(0),
						Quantums: dtypes.NewInt(275_000_000_000),
					},
				},
			},
			updates: []types.Update{
				{
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(-25_000_000_000)), // -$25,000
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(0),
							BigQuantumsDelta: big.NewInt(50_000_000), // .5 BTC
						},
					},
				},
			},
			expectedTDaiYieldPoolBalance: big.NewInt(0),
			expectedCollateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 300_000_000_000,
			},
			msgSenderEnabled: true,
		},
		"Successfully claims yield for isolated market": {
			assetPositions:            testutil.CreateTDaiAssetPosition(big.NewInt(100_000_000_000)),
			subaccountAssetYieldIndex: constants.AssetYieldIndex_Zero,
			globalAssetYieldIndex:     big.NewRat(9, 4),
			fundsInTDaiPool:           big.NewInt(200_000_000_000),
			collateralPoolTDaiBalances: map[string]int64{
				authtypes.NewModuleAddress(
					types.ModuleName + ":" + lib.UintToString(constants.PerpetualPosition_OneISOLong.PerpetualId),
				).String(): 100_000_000_000,
			},
			perpetuals: []perptypes.Perpetual{
				{
					Params:       constants.IsoUsd_IsolatedMarket.Params,
					FundingIndex: constants.IsoUsd_IsolatedMarket.FundingIndex,
					OpenInterest: constants.IsoUsd_IsolatedMarket.OpenInterest,
					YieldIndex:   big.NewRat(4, 5).String(),
				},
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(3),
					Quantums:     dtypes.NewInt(-100_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			expectedSuccess:          true,
			expectedSuccessPerUpdate: []types.UpdateResult{types.Success},
			expectedAssetYieldIndex:  big.NewRat(9, 4).String(),
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(3),
					Quantums:     dtypes.NewInt(-150_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(4, 5).String(),
				},
			},
			expectedUpdatedPerpetualPositions: map[types.SubaccountId][]*types.PerpetualPosition{
				defaultSubaccountId: {
					{
						PerpetualId:  uint32(3),
						Quantums:     dtypes.NewInt(-150_000_000),
						FundingIndex: dtypes.NewInt(0),
						YieldIndex:   big.NewRat(4, 5).String(),
					},
				},
			},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(249920000000),
				},
			},
			expectedUpdatedAssetPositions: map[types.SubaccountId][]*types.AssetPosition{
				defaultSubaccountId: {
					{
						AssetId:  uint32(0),
						Quantums: dtypes.NewInt(249920000000),
					},
				},
			},
			updates: []types.Update{
				{
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(25_000_000_000)),
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(3),
							BigQuantumsDelta: big.NewInt(-50_000_000),
						},
					},
				},
			},
			expectedTDaiYieldPoolBalance: big.NewInt(75080000000),
			expectedCollateralPoolTDaiBalances: map[string]int64{
				authtypes.NewModuleAddress(
					types.ModuleName + ":" + lib.UintToString(constants.PerpetualPosition_OneISOLong.PerpetualId),
				).String(): 224_920_000_000,
				types.ModuleAddress.String(): 0,
			},
			msgSenderEnabled: true,
		},
		"Successfully claims yield for multiple subaccounts": {
			globalAssetYieldIndex: big.NewRat(9, 4),
			fundsInTDaiPool:       big.NewInt(1_200_000_000_000),
			collateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 1_000_000_000_000,
			},
			subaccountAssetYieldIndex: constants.AssetYieldIndex_Zero,
			assetPositions:            testutil.CreateTDaiAssetPosition(big.NewInt(100_000_000_000)),
			additionalTestSubaccounts: []types.Subaccount{
				{
					Id: &constants.Bob_Num0,
					AssetPositions: testutil.CreateTDaiAssetPosition(big.NewInt(
						900_000_000_000,
					)),
					AssetYieldIndex: big.NewRat(3, 2).String(),
				},
			},
			perpetuals: []perptypes.Perpetual{
				{
					Params:       constants.BtcUsd_NoMarginRequirement.Params,
					FundingIndex: constants.BtcUsd_NoMarginRequirement.FundingIndex,
					OpenInterest: constants.BtcUsd_NoMarginRequirement.OpenInterest,
					YieldIndex:   big.NewRat(4, 5).String(),
				},
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(-100_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			expectedSuccess:          true,
			expectedSuccessPerUpdate: []types.UpdateResult{types.Success, types.Success},
			expectedAssetYieldIndex:  big.NewRat(9, 4).String(),
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(-9_100_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(4, 5).String(),
				},
			},
			expectedUpdatedPerpetualPositions: map[types.SubaccountId][]*types.PerpetualPosition{
				defaultSubaccountId: {
					{
						PerpetualId:  uint32(0),
						Quantums:     dtypes.NewInt(-9_100_000_000),
						FundingIndex: dtypes.NewInt(0),
						YieldIndex:   big.NewRat(4, 5).String(),
					},
				},
				constants.Bob_Num0: {
					{
						PerpetualId:  uint32(0),
						Quantums:     dtypes.NewInt(9_000_000_000),
						FundingIndex: dtypes.NewInt(0),
						YieldIndex:   big.NewRat(4, 5).String(),
					},
				},
			},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(4_724_920_000_000), // Yield Collected: 124_920_000_000 tDAI
				},
			},
			expectedUpdatedAssetPositions: map[types.SubaccountId][]*types.AssetPosition{
				defaultSubaccountId: {
					{
						AssetId:  uint32(0),
						Quantums: dtypes.NewInt(4_724_920_000_000),
					},
				},
				constants.Bob_Num0: {
					{
						AssetId:  uint32(0),
						Quantums: dtypes.NewInt(-3_150_000_000_000), // Bob Yield: 450_000_000_000 tDAI
					},
				},
			},
			updates: []types.Update{
				{
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(0),
							BigQuantumsDelta: big.NewInt(9_000_000_000), // 90 BTC
						},
					},
					AssetUpdates: []types.AssetUpdate{
						{
							AssetId:          uint32(0),
							BigQuantumsDelta: big.NewInt(-4_500_000_000_000), // -4,500,000 TDai
						},
					},
					SubaccountId: constants.Bob_Num0,
				},
				{
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(0),
							BigQuantumsDelta: big.NewInt(-9_000_000_000), // 9 BTC
						},
					},
					AssetUpdates: []types.AssetUpdate{
						{
							AssetId:          uint32(0),
							BigQuantumsDelta: big.NewInt(4_500_000_000_000), // 4,500,000 TDai
						},
					},
				},
			},
			expectedTDaiYieldPoolBalance: big.NewInt(625_080_000_000),
			expectedCollateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 1_574_920_000_000,
			},
			msgSenderEnabled: true,
		},
		"Successfully claims yield from tdai position when there is unsettled funding": {
			globalAssetYieldIndex: big.NewRat(9, 4),
			fundsInTDaiPool:       big.NewInt(1_200_000_000_000),
			collateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 1_000_000_000_000,
			},
			subaccountAssetYieldIndex: constants.AssetYieldIndex_Zero,
			assetPositions:            testutil.CreateTDaiAssetPosition(big.NewInt(100_000_000_000)), // $100,000
			expectedSuccess:           true,
			expectedSuccessPerUpdate:  []types.UpdateResult{types.Success},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_NoMarginRequirement,
				constants.EthUsd_NoMarginRequirement,
			},
			newFundingIndices: []*big.Int{big.NewInt(-10), big.NewInt(-8)},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId: uint32(0),
					Quantums:    dtypes.NewInt(-100_000_000), // -1 BTC
					// indexDelta=-5
					FundingIndex: dtypes.NewInt(-5),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
				{
					PerpetualId: uint32(1),
					Quantums:    dtypes.NewInt(-2_000_000_000), // -2 ETH
					// indexDelta=-2
					FundingIndex: dtypes.NewInt(-6),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(-100_000_000), // -1 BTC
					FundingIndex: dtypes.NewInt(-10),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
				{
					PerpetualId:  uint32(1),
					Quantums:     dtypes.NewInt(-1_000_000_000), // -1 ETH
					FundingIndex: dtypes.NewInt(-8),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			expectedUpdatedPerpetualPositions: map[types.SubaccountId][]*types.PerpetualPosition{
				defaultSubaccountId: {
					{
						PerpetualId:  uint32(0),
						Quantums:     dtypes.NewInt(-100_000_000), // -1 BTC
						FundingIndex: dtypes.NewInt(-10),
					},
					{
						PerpetualId:  uint32(1),
						Quantums:     dtypes.NewInt(-1_000_000_000), // -1 ETH
						FundingIndex: dtypes.NewInt(-8),
					},
				},
			},
			expectedSubaccoundIdToFundingPayments: map[types.SubaccountId]map[uint32]dtypes.SerializableInt{
				defaultSubaccountId: {
					// indexDelta=-5, settlement=5*-100_000_000/1_000_000=-500
					uint32(0): dtypes.NewInt(500),
					// indexDelta=-2, settlement=2*-2_000_000_000/1_000_000=-4_000
					uint32(1): dtypes.NewInt(4_000),
				},
			},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId: uint32(0),
					// Original Asset Position - Funding Payments
					// = 100_000_000_000 - 4_000 - 500
					// = 99_999_995_500
					Quantums: dtypes.NewInt(224_999_995_500),
				},
			},
			updates: []types.Update{
				{
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(1),
							BigQuantumsDelta: big.NewInt(1_000_000_000), // 1 ETH
						},
					},
				},
			},
			expectedAssetYieldIndex:      big.NewRat(9, 4).String(),
			expectedTDaiYieldPoolBalance: big.NewInt(1_075_000_000_000),
			expectedCollateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 1_125_000_000_000,
			},
			msgSenderEnabled: true,
		},
		"Successfully claims yield from tdai position and perp positions when there is unsettled funding": {
			globalAssetYieldIndex: big.NewRat(5, 4),
			fundsInTDaiPool:       big.NewInt(1_200_000_000_000),
			collateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 1_000_000_000_000,
			},
			subaccountAssetYieldIndex: constants.AssetYieldIndex_Zero,
			assetPositions:            testutil.CreateTDaiAssetPosition(big.NewInt(100_000_000_000)),
			expectedSuccess:           true,
			expectedSuccessPerUpdate:  []types.UpdateResult{types.Success},
			perpetuals: []perptypes.Perpetual{
				{
					Params:       constants.BtcUsd_NoMarginRequirement.Params,
					FundingIndex: constants.BtcUsd_NoMarginRequirement.FundingIndex,
					OpenInterest: constants.BtcUsd_NoMarginRequirement.OpenInterest,
					YieldIndex:   big.NewRat(3, 4).String(),
				},
				{
					Params:       constants.EthUsd_NoMarginRequirement.Params,
					FundingIndex: constants.EthUsd_NoMarginRequirement.FundingIndex,
					OpenInterest: constants.EthUsd_NoMarginRequirement.OpenInterest,
					YieldIndex:   big.NewRat(1, 7).String(),
				},
			},
			newFundingIndices: []*big.Int{big.NewInt(-10), big.NewInt(-8)},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId: uint32(0),
					Quantums:    dtypes.NewInt(-100_000_000), // -1 BTC
					// indexDelta=-5
					FundingIndex: dtypes.NewInt(-5),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
				{
					PerpetualId: uint32(1),
					Quantums:    dtypes.NewInt(-2_000_000_000), // -2 ETH
					// indexDelta=-2
					FundingIndex: dtypes.NewInt(-6),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(-100_000_000), // -1 BTC
					FundingIndex: dtypes.NewInt(-10),
					YieldIndex:   big.NewRat(3, 4).String(),
				},
				{
					PerpetualId:  uint32(1),
					Quantums:     dtypes.NewInt(-1_000_000_000), // -1 ETH
					FundingIndex: dtypes.NewInt(-8),
					YieldIndex:   big.NewRat(1, 7).String(),
				},
			},
			expectedUpdatedPerpetualPositions: map[types.SubaccountId][]*types.PerpetualPosition{
				defaultSubaccountId: {
					{
						PerpetualId:  uint32(0),
						Quantums:     dtypes.NewInt(-100_000_000), // -1 BTC
						FundingIndex: dtypes.NewInt(-10),
						YieldIndex:   big.NewRat(3, 4).String(),
					},
					{
						PerpetualId:  uint32(1),
						Quantums:     dtypes.NewInt(-1_000_000_000), // -1 ETH
						FundingIndex: dtypes.NewInt(-8),
						YieldIndex:   big.NewRat(1, 7).String(),
					},
				},
			},
			expectedSubaccoundIdToFundingPayments: map[types.SubaccountId]map[uint32]dtypes.SerializableInt{
				defaultSubaccountId: {
					// indexDelta=-5, settlement=5*-100_000_000/1_000_000=-500
					uint32(0): dtypes.NewInt(500),
					// indexDelta=-2, settlement=2*-2_000_000_000/1_000_000=-4_000
					uint32(1): dtypes.NewInt(4_000),
				},
			},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId: uint32(0),
					// Original Asset Position - Funding Payments
					// = 100_000_000_000 - 4_000 - 500
					// = 99_999_995_500
					Quantums: dtypes.NewInt(124_639_281_214),
				},
			},
			updates: []types.Update{
				{
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(1),
							BigQuantumsDelta: big.NewInt(1_000_000_000), // 1 ETH
						},
					},
				},
			},
			expectedAssetYieldIndex:      big.NewRat(5, 4).String(),
			expectedTDaiYieldPoolBalance: big.NewInt(1_175_360_714_286),
			expectedCollateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 1_024_639_285_714,
			},
			msgSenderEnabled: true,
		},
		"Successfully does not claim yield when asset yield index is already updated": {
			globalAssetYieldIndex: big.NewRat(5, 4),
			fundsInTDaiPool:       big.NewInt(1_200_000_000_000),
			collateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 1_000_000_000_000,
			},
			subaccountAssetYieldIndex: big.NewRat(5, 4).String(),
			assetPositions:            testutil.CreateTDaiAssetPosition(big.NewInt(100_000_000_000)),
			additionalTestSubaccounts: []types.Subaccount{
				{
					Id: &constants.Bob_Num0,
					AssetPositions: testutil.CreateTDaiAssetPosition(big.NewInt(
						900_000_000_000,
					)),
					AssetYieldIndex: big.NewRat(5, 4).String(),
				},
			},
			perpetuals: []perptypes.Perpetual{
				{
					Params:       constants.BtcUsd_NoMarginRequirement.Params,
					FundingIndex: constants.BtcUsd_NoMarginRequirement.FundingIndex,
					OpenInterest: constants.BtcUsd_NoMarginRequirement.OpenInterest,
					YieldIndex:   big.NewRat(4, 5).String(),
				},
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(-100_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(4, 5).String(),
				},
			},
			expectedSuccess:          true,
			expectedSuccessPerUpdate: []types.UpdateResult{types.Success, types.Success},
			expectedAssetYieldIndex:  big.NewRat(5, 4).String(),
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(-9_100_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(4, 5).String(),
				},
			},
			expectedUpdatedPerpetualPositions: map[types.SubaccountId][]*types.PerpetualPosition{
				defaultSubaccountId: {
					{
						PerpetualId:  uint32(0),
						Quantums:     dtypes.NewInt(-9_100_000_000),
						FundingIndex: dtypes.NewInt(0),
						YieldIndex:   big.NewRat(4, 5).String(),
					},
				},
				constants.Bob_Num0: {
					{
						PerpetualId:  uint32(0),
						Quantums:     dtypes.NewInt(9_000_000_000),
						FundingIndex: dtypes.NewInt(0),
						YieldIndex:   big.NewRat(4, 5).String(),
					},
				},
			},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(4_600_000_000_000), // Yield Collected: 0 tDAI
				},
			},
			expectedUpdatedAssetPositions: map[types.SubaccountId][]*types.AssetPosition{
				defaultSubaccountId: {
					{
						AssetId:  uint32(0),
						Quantums: dtypes.NewInt(4_600_000_000_000),
					},
				},
				constants.Bob_Num0: {
					{
						AssetId:  uint32(0),
						Quantums: dtypes.NewInt(-3_600_000_000_000), // Bob Yield: 0 tDAI
					},
				},
			},
			updates: []types.Update{
				{
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(0),
							BigQuantumsDelta: big.NewInt(9_000_000_000), // 90 BTC
						},
					},
					AssetUpdates: []types.AssetUpdate{
						{
							AssetId:          uint32(0),
							BigQuantumsDelta: big.NewInt(-4_500_000_000_000), // -4,500,000 TDai
						},
					},
					SubaccountId: constants.Bob_Num0,
				},
				{
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(0),
							BigQuantumsDelta: big.NewInt(-9_000_000_000), // 9 BTC
						},
					},
					AssetUpdates: []types.AssetUpdate{
						{
							AssetId:          uint32(0),
							BigQuantumsDelta: big.NewInt(4_500_000_000_000), // 4,500,000 TDai
						},
					},
				},
			},
			expectedTDaiYieldPoolBalance: big.NewInt(1_200_000_000_000),
			expectedCollateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 1_000_000_000_000,
			},
			msgSenderEnabled: true,
		},
		"Successfully does not claim yield when perp position is created": {
			assetPositions:            testutil.CreateTDaiAssetPosition(big.NewInt(100_000_000_000)), // $100,000
			subaccountAssetYieldIndex: big.NewRat(1, 1).String(),
			globalAssetYieldIndex:     big.NewRat(1, 1),
			fundsInTDaiPool:           big.NewInt(200_000_000_000),
			collateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 100_000_000_000,
			},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_NoMarginRequirement,
			},
			expectedSuccess:          true,
			expectedSuccessPerUpdate: []types.UpdateResult{types.Success},
			expectedAssetYieldIndex:  big.NewRat(1, 1).String(),
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(50_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			expectedUpdatedPerpetualPositions: map[types.SubaccountId][]*types.PerpetualPosition{
				defaultSubaccountId: {
					{
						PerpetualId:  uint32(0),
						Quantums:     dtypes.NewInt(50_000_000),
						FundingIndex: dtypes.NewInt(0),
					},
				},
			},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(75_000_000_000),
				},
			},
			expectedUpdatedAssetPositions: map[types.SubaccountId][]*types.AssetPosition{
				defaultSubaccountId: {
					{
						AssetId:  uint32(0),
						Quantums: dtypes.NewInt(75_000_000_000),
					},
				},
			},
			updates: []types.Update{
				{
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(-25_000_000_000)), // -$25,000
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(0),
							BigQuantumsDelta: big.NewInt(50_000_000), // .5 BTC
						},
					},
				},
			},
			expectedTDaiYieldPoolBalance: big.NewInt(200_000_000_000),
			expectedCollateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 100_000_000_000,
			},
			msgSenderEnabled: true,
		},
		"Successfully does not claim yield when negative positions cancel out positive position yield claims": {
			assetPositions:            testutil.CreateTDaiAssetPosition(big.NewInt(100_000_000_000)), // $100,000
			subaccountAssetYieldIndex: big.NewRat(1, 1).String(),
			globalAssetYieldIndex:     big.NewRat(1, 1),
			fundsInTDaiPool:           big.NewInt(200_000_000_000),
			collateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 100_000_000_000,
			},
			perpetuals: []perptypes.Perpetual{
				{
					Params:       constants.BtcUsd_NoMarginRequirement.Params,
					FundingIndex: constants.BtcUsd_NoMarginRequirement.FundingIndex,
					OpenInterest: constants.BtcUsd_NoMarginRequirement.OpenInterest,
					YieldIndex:   big.NewRat(1000, 1).String(),
				},
			},
			expectedSuccess:          true,
			expectedSuccessPerUpdate: []types.UpdateResult{types.Success},
			expectedAssetYieldIndex:  big.NewRat(1, 1).String(),
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(-100_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(-50_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(1000, 1).String(),
				},
			},
			expectedUpdatedPerpetualPositions: map[types.SubaccountId][]*types.PerpetualPosition{
				defaultSubaccountId: {
					{
						PerpetualId:  uint32(0),
						Quantums:     dtypes.NewInt(-50_000_000),
						FundingIndex: dtypes.NewInt(0),
						YieldIndex:   big.NewRat(1000, 1).String(),
					},
				},
			},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(75_000_000_000),
				},
			},
			expectedUpdatedAssetPositions: map[types.SubaccountId][]*types.AssetPosition{
				defaultSubaccountId: {
					{
						AssetId:  uint32(0),
						Quantums: dtypes.NewInt(75_000_000_000),
					},
				},
			},
			updates: []types.Update{
				{
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(-25_000_000_000)), // -$25,000
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(0),
							BigQuantumsDelta: big.NewInt(50_000_000), // .5 BTC
						},
					},
				},
			},
			expectedTDaiYieldPoolBalance: big.NewInt(200_000_000_000),
			expectedCollateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 100_000_000_000,
			},
			msgSenderEnabled: true,
		},
		"Fails yield claim: Negative general asset yield index": {
			assetPositions:            testutil.CreateTDaiAssetPosition(big.NewInt(100_000_000_000)), // $100,000
			subaccountAssetYieldIndex: constants.AssetYieldIndex_Zero,
			globalAssetYieldIndex:     big.NewRat(-1, 1),
			fundsInTDaiPool:           big.NewInt(200_000_000_000),
			collateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 100_000_000_000,
			},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_NoMarginRequirement,
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			expectedSuccess:         false,
			expectedErr:             types.ErrGlobalYieldIndexNegative,
			expectedAssetYieldIndex: big.NewRat(-1, 1).String(),
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			expectedUpdatedPerpetualPositions: map[types.SubaccountId][]*types.PerpetualPosition{
				defaultSubaccountId: {
					{
						PerpetualId:  uint32(0),
						Quantums:     dtypes.NewInt(100_000_000),
						FundingIndex: dtypes.NewInt(0),
					},
				},
			},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(100_000_000_000),
				},
			},
			expectedUpdatedAssetPositions: map[types.SubaccountId][]*types.AssetPosition{
				defaultSubaccountId: {
					{
						AssetId:  uint32(0),
						Quantums: dtypes.NewInt(100_000_000_000),
					},
				},
			},
			updates: []types.Update{
				{
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(-25_000_000_000)), // -$25,000
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(0),
							BigQuantumsDelta: big.NewInt(50_000_000), // .5 BTC
						},
					},
				},
			},
			expectedTDaiYieldPoolBalance: big.NewInt(200_000_000_000),
			expectedCollateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 100_000_000_000,
			},
			msgSenderEnabled: true,
		},
		"Fails yield claim: Asset yield index in account higher than in general ": {
			assetPositions:            testutil.CreateTDaiAssetPosition(big.NewInt(100_000_000_000)), // $100,000
			subaccountAssetYieldIndex: big.NewRat(1, 1).String(),
			globalAssetYieldIndex:     big.NewRat(1, 2),
			fundsInTDaiPool:           big.NewInt(200_000_000_000),
			collateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 100_000_000_000,
			},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_NoMarginRequirement,
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			expectedSuccess:         false,
			expectedErr:             types.ErrGeneralYieldIndexSmallerThanYieldIndexInSubaccount,
			expectedAssetYieldIndex: big.NewRat(1, 2).String(),
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			expectedUpdatedPerpetualPositions: map[types.SubaccountId][]*types.PerpetualPosition{
				defaultSubaccountId: {
					{
						PerpetualId:  uint32(0),
						Quantums:     dtypes.NewInt(100_000_000), // 1.5 BTC
						FundingIndex: dtypes.NewInt(0),
					},
				},
			},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(100_000_000_000),
				},
			},
			expectedUpdatedAssetPositions: map[types.SubaccountId][]*types.AssetPosition{
				defaultSubaccountId: {
					{
						AssetId:  uint32(0),
						Quantums: dtypes.NewInt(100_000_000_000),
					},
				},
			},
			updates: []types.Update{
				{
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(-25_000_000_000)), // -$25,000
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(0),
							BigQuantumsDelta: big.NewInt(50_000_000), // .5 BTC
						},
					},
				},
			},
			expectedTDaiYieldPoolBalance: big.NewInt(200_000_000_000),
			expectedCollateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 100_000_000_000,
			},
			msgSenderEnabled: true,
		},
		"Fails yield claim: Negative general perp yield index": {
			assetPositions:            testutil.CreateTDaiAssetPosition(big.NewInt(100_000_000_000)), // $100,000
			subaccountAssetYieldIndex: constants.AssetYieldIndex_Zero,
			globalAssetYieldIndex:     big.NewRat(1, 1),
			fundsInTDaiPool:           big.NewInt(200_000_000_000),
			collateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 100_000_000_000,
			},
			perpetuals: []perptypes.Perpetual{
				{
					Params:       constants.BtcUsd_NoMarginRequirement.Params,
					FundingIndex: constants.BtcUsd_NoMarginRequirement.FundingIndex,
					OpenInterest: constants.BtcUsd_NoMarginRequirement.OpenInterest,
					YieldIndex:   big.NewRat(-1, 1).String(),
				},
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			expectedSuccess:         false,
			expectedErr:             types.ErrGlobalYieldIndexNegative,
			expectedAssetYieldIndex: big.NewRat(-1, 1).String(),
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			expectedUpdatedPerpetualPositions: map[types.SubaccountId][]*types.PerpetualPosition{
				defaultSubaccountId: {
					{
						PerpetualId:  uint32(0),
						Quantums:     dtypes.NewInt(100_000_000),
						FundingIndex: dtypes.NewInt(0),
					},
				},
			},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(100_000_000_000),
				},
			},
			expectedUpdatedAssetPositions: map[types.SubaccountId][]*types.AssetPosition{
				defaultSubaccountId: {
					{
						AssetId:  uint32(0),
						Quantums: dtypes.NewInt(100_000_000_000),
					},
				},
			},
			updates: []types.Update{
				{
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(-25_000_000_000)), // -$25,000
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(0),
							BigQuantumsDelta: big.NewInt(50_000_000), // .5 BTC
						},
					},
				},
			},
			expectedTDaiYieldPoolBalance: big.NewInt(200_000_000_000),
			expectedCollateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 100_000_000_000,
			},
			msgSenderEnabled: true,
		},
		"Fails yield claim: Perp yield index in subaccount higher than in general": {
			assetPositions:            testutil.CreateTDaiAssetPosition(big.NewInt(100_000_000_000)), // $100,000
			subaccountAssetYieldIndex: constants.AssetYieldIndex_Zero,
			globalAssetYieldIndex:     big.NewRat(1, 1),
			fundsInTDaiPool:           big.NewInt(200_000_000_000),
			collateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 100_000_000_000,
			},
			perpetuals: []perptypes.Perpetual{
				{
					Params:       constants.BtcUsd_NoMarginRequirement.Params,
					FundingIndex: constants.BtcUsd_NoMarginRequirement.FundingIndex,
					OpenInterest: constants.BtcUsd_NoMarginRequirement.OpenInterest,
					YieldIndex:   big.NewRat(1, 2).String(),
				},
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(1, 1).String(),
				},
			},
			expectedSuccess:         false,
			expectedErr:             types.ErrGeneralYieldIndexSmallerThanYieldIndexInSubaccount,
			expectedAssetYieldIndex: big.NewRat(1, 1).String(),
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(1, 1).String(),
				},
			},
			expectedUpdatedPerpetualPositions: map[types.SubaccountId][]*types.PerpetualPosition{
				defaultSubaccountId: {
					{
						PerpetualId:  uint32(0),
						Quantums:     dtypes.NewInt(100_000_000),
						FundingIndex: dtypes.NewInt(0),
					},
				},
			},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(100_000_000_000),
				},
			},
			expectedUpdatedAssetPositions: map[types.SubaccountId][]*types.AssetPosition{
				defaultSubaccountId: {
					{
						AssetId:  uint32(0),
						Quantums: dtypes.NewInt(100_000_000_000),
					},
				},
			},
			updates: []types.Update{
				{
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(-25_000_000_000)), // -$25,000
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(0),
							BigQuantumsDelta: big.NewInt(50_000_000), // .5 BTC
						},
					},
				},
			},
			expectedTDaiYieldPoolBalance: big.NewInt(200_000_000_000),
			expectedCollateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 100_000_000_000,
			},
			msgSenderEnabled: true,
		},
		"Fails yield claim: Perp yield index in subaccount badly initialized": {
			assetPositions:            testutil.CreateTDaiAssetPosition(big.NewInt(100_000_000_000)), // $100,000
			subaccountAssetYieldIndex: constants.AssetYieldIndex_Zero,
			globalAssetYieldIndex:     big.NewRat(1, 1),
			fundsInTDaiPool:           big.NewInt(200_000_000_000),
			collateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 100_000_000_000,
			},
			perpetuals: []perptypes.Perpetual{
				{
					Params:       constants.BtcUsd_NoMarginRequirement.Params,
					FundingIndex: constants.BtcUsd_NoMarginRequirement.FundingIndex,
					OpenInterest: constants.BtcUsd_NoMarginRequirement.OpenInterest,
					YieldIndex:   big.NewRat(1, 2).String(),
				},
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   "",
				},
			},
			expectedSuccess:         false,
			expectedErr:             types.ErrYieldIndexUninitialized,
			expectedAssetYieldIndex: big.NewRat(1, 1).String(),
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   "",
				},
			},
			expectedUpdatedPerpetualPositions: map[types.SubaccountId][]*types.PerpetualPosition{
				defaultSubaccountId: {
					{
						PerpetualId:  uint32(0),
						Quantums:     dtypes.NewInt(100_000_000),
						FundingIndex: dtypes.NewInt(0),
					},
				},
			},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(100_000_000_000),
				},
			},
			expectedUpdatedAssetPositions: map[types.SubaccountId][]*types.AssetPosition{
				defaultSubaccountId: {
					{
						AssetId:  uint32(0),
						Quantums: dtypes.NewInt(100_000_000_000),
					},
				},
			},
			updates: []types.Update{
				{
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(-25_000_000_000)), // -$25,000
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(0),
							BigQuantumsDelta: big.NewInt(50_000_000), // .5 BTC
						},
					},
				},
			},
			expectedTDaiYieldPoolBalance: big.NewInt(200_000_000_000),
			expectedCollateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 100_000_000_000,
			},
			msgSenderEnabled: true,
		},
		"Successful yield claim: not enough yield in tdai pool so we take what's available": {
			assetPositions:            testutil.CreateTDaiAssetPosition(big.NewInt(100_000_000_000)), // $100,000
			subaccountAssetYieldIndex: constants.AssetYieldIndex_Zero,
			globalAssetYieldIndex:     big.NewRat(1, 1),
			fundsInTDaiPool:           big.NewInt(1),
			collateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 100_000_000_000,
			},
			perpetuals: []perptypes.Perpetual{
				{
					Params:       constants.BtcUsd_NoMarginRequirement.Params,
					FundingIndex: constants.BtcUsd_NoMarginRequirement.FundingIndex,
					OpenInterest: constants.BtcUsd_NoMarginRequirement.OpenInterest,
					YieldIndex:   big.NewRat(1, 2).String(),
				},
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			expectedSuccess:          true,
			expectedSuccessPerUpdate: []types.UpdateResult{types.Success},
			expectedAssetYieldIndex:  big.NewRat(1, 1).String(),
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(150_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(1, 2).String(),
				},
			},
			expectedUpdatedPerpetualPositions: map[types.SubaccountId][]*types.PerpetualPosition{
				defaultSubaccountId: {
					{
						PerpetualId:  uint32(0),
						Quantums:     dtypes.NewInt(150_000_000),
						FundingIndex: dtypes.NewInt(0),
					},
				},
			},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(75_000_000_001),
				},
			},
			expectedUpdatedAssetPositions: map[types.SubaccountId][]*types.AssetPosition{
				defaultSubaccountId: {
					{
						AssetId:  uint32(0),
						Quantums: dtypes.NewInt(75_000_000_001),
					},
				},
			},
			updates: []types.Update{
				{
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(-25_000_000_000)), // -$25,000
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(0),
							BigQuantumsDelta: big.NewInt(50_000_000), // .5 BTC
						},
					},
				},
			},
			expectedTDaiYieldPoolBalance: big.NewInt(0),
			expectedCollateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 100_000_000_001,
			},
			msgSenderEnabled: true,
		},
		"Successfully claims 0 yield when subaccount's yield is negative": {
			assetPositions:            testutil.CreateTDaiAssetPosition(big.NewInt(100_000_000_000)), // $100,000
			subaccountAssetYieldIndex: constants.AssetYieldIndex_Zero,
			globalAssetYieldIndex:     big.NewRat(1, 1),
			fundsInTDaiPool:           big.NewInt(200_000_000_000),
			collateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 100_000_000_000,
			},
			perpetuals: []perptypes.Perpetual{
				{
					Params:       constants.BtcUsd_NoMarginRequirement.Params,
					FundingIndex: constants.BtcUsd_NoMarginRequirement.FundingIndex,
					OpenInterest: constants.BtcUsd_NoMarginRequirement.OpenInterest,
					YieldIndex:   big.NewRat(10_000, 1).String(),
				},
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(-100_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			expectedSuccess:          true,
			expectedSuccessPerUpdate: []types.UpdateResult{types.Success},
			expectedAssetYieldIndex:  big.NewRat(1, 1).String(),
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(-50_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(10_000, 1).String(),
				},
			},
			expectedUpdatedPerpetualPositions: map[types.SubaccountId][]*types.PerpetualPosition{
				defaultSubaccountId: {
					{
						PerpetualId:  uint32(0),
						Quantums:     dtypes.NewInt(-50_000_000),
						FundingIndex: dtypes.NewInt(0),
						YieldIndex:   big.NewRat(10_000, 1).String(),
					},
				},
			},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(75_000_000_000),
				},
			},
			expectedUpdatedAssetPositions: map[types.SubaccountId][]*types.AssetPosition{
				defaultSubaccountId: {
					{
						AssetId:  uint32(0),
						Quantums: dtypes.NewInt(75_000_000_000),
					},
				},
			},
			updates: []types.Update{
				{
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(-25_000_000_000)), // -$25,000
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(0),
							BigQuantumsDelta: big.NewInt(50_000_000), // .5 BTC
						},
					},
				},
			},
			expectedTDaiYieldPoolBalance: big.NewInt(200_000_000_000),
			expectedCollateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 100_000_000_000,
			},
			msgSenderEnabled: true,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctx, keeper, pricesKeeper, perpetualsKeeper, accountKeeper, bankKeeper, assetsKeeper, rateLimitKeeper, _, _ := testutil.SubaccountsKeepers(
				t,
				tc.msgSenderEnabled,
			)
			ctx = ctx.WithTxBytes(constants.TestTxBytes)
			testutil.CreateTestMarkets(t, ctx, pricesKeeper)
			testutil.CreateTestLiquidityTiers(t, ctx, perpetualsKeeper)

			// Set up initial sdai price
			rateString := sdaiservertypes.TestSDAIEventRequest.ConversionRate
			rate, conversionErr := ratelimitkeeper.ConvertStringToBigInt(rateString)
			require.NoError(t, conversionErr)

			rateLimitKeeper.SetSDAIPrice(ctx, rate)
			globalAssetYieldIndex := big.NewRat(1, 1)
			if tc.globalAssetYieldIndex != nil {
				globalAssetYieldIndex = tc.globalAssetYieldIndex
			}
			rateLimitKeeper.SetAssetYieldIndex(ctx, globalAssetYieldIndex)

			for _, m := range tc.marketParamPrices {
				_, err := pricesKeeper.CreateMarket(
					ctx,
					m.Param,
					m.Price,
				)
				require.NoError(t, err)
			}

			// Always creates TDai asset first
			require.NoError(t, testutil.CreateTDaiAsset(ctx, assetsKeeper))
			for _, a := range tc.assets {
				_, err := assetsKeeper.CreateAsset(
					ctx,
					a.Id,
					a.Symbol,
					a.Denom,
					a.DenomExponent,
					a.HasMarket,
					a.MarketId,
					a.AtomicResolution,
					a.AssetYieldIndex,
				)
				require.NoError(t, err)
			}

			for i, p := range tc.perpetuals {
				perpetualsKeeper.SetPerpetualForTest(ctx, p)

				// Update FundingIndex for testing settlements.
				if i < len(tc.newFundingIndices) {
					err := perpetualsKeeper.ModifyFundingIndex(
						ctx,
						p.Params.Id,
						tc.newFundingIndices[i],
					)
					require.NoError(t, err)
				}
			}

			for collateralPoolAddr, TDaiBal := range tc.collateralPoolTDaiBalances {
				err := bank_testutil.FundAccount(
					ctx,
					sdk.MustAccAddressFromBech32(collateralPoolAddr),
					sdk.Coins{
						sdk.NewCoin(asstypes.AssetTDai.Denom, sdkmath.NewInt(TDaiBal)),
					},
					*bankKeeper,
				)
				require.NoError(t, err)
			}

			if tc.fundsInTDaiPool != nil {
				err := bank_testutil.FundModuleAccount(
					ctx,
					ratelimittypes.TDaiPoolAccount,
					sdk.Coins{
						sdk.NewCoin(asstypes.AssetTDai.Denom, sdkmath.NewIntFromBigInt(tc.fundsInTDaiPool)),
					},
					*bankKeeper,
				)
				require.NoError(t, err)
			}

			subaccount := createNSubaccount(keeper, ctx, 1, big.NewInt(1_000))[0]
			subaccount.PerpetualPositions = tc.perpetualPositions
			subaccount.AssetPositions = tc.assetPositions
			subaccountYieldIndex := constants.AssetYieldIndex_Zero
			if tc.subaccountAssetYieldIndex != "" {
				subaccountYieldIndex = tc.subaccountAssetYieldIndex
			}
			subaccount.AssetYieldIndex = subaccountYieldIndex
			keeper.SetSubaccount(ctx, subaccount)
			subaccountId := *subaccount.Id

			for _, sa := range tc.additionalTestSubaccounts {
				keeper.SetSubaccount(ctx, sa)
			}

			for i, u := range tc.updates {
				if u.SubaccountId == (types.SubaccountId{}) {
					u.SubaccountId = subaccountId
				}
				tc.updates[i] = u
			}

			updateType := types.CollatCheck
			if tc.updateType != types.UpdateTypeUnspecified {
				updateType = tc.updateType
			}
			success, successPerUpdate, err := keeper.UpdateSubaccounts(ctx, tc.updates, updateType)
			if tc.expectedErr != nil {
				require.ErrorIs(t, err, tc.expectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedSuccessPerUpdate, successPerUpdate)
				require.Equal(t, tc.expectedSuccess, success)
			}
			newSubaccount := keeper.GetSubaccount(ctx, subaccountId)
			require.Equal(t, len(newSubaccount.PerpetualPositions), len(tc.expectedPerpetualPositions))
			for i, ep := range tc.expectedPerpetualPositions {
				require.Equal(t, *ep, *newSubaccount.PerpetualPositions[i])
			}
			require.Equal(t, len(newSubaccount.AssetPositions), len(tc.expectedAssetPositions))
			for i, ep := range tc.expectedAssetPositions {
				require.Equal(t, *ep, *newSubaccount.AssetPositions[i])
			}
			if tc.expectedErr == nil {
				require.Equal(t, 0, globalAssetYieldIndex.Cmp(ratelimitkeeper.ConvertStringToBigRatWithPanicOnErr(newSubaccount.AssetYieldIndex)),
					"Expected AssetYieldIndex %v. Got %v.", globalAssetYieldIndex, newSubaccount.AssetYieldIndex,
				)
			}

			if tc.msgSenderEnabled {
				assertSubaccountUpdateEventsInIndexerBlock(
					t,
					keeper,
					ctx,
					tc.expectedErr,
					tc.expectedSuccess,
					tc.updates,
					tc.expectedSuccessPerUpdate,
					tc.expectedUpdatedPerpetualPositions,
					tc.expectedSubaccoundIdToFundingPayments,
					tc.expectedUpdatedAssetPositions,
					tc.expectedAssetYieldIndex,
				)
			} else {
				assertSubaccountUpdateEventsNotInIndexerBlock(
					t,
					keeper,
					ctx,
				)
			}

			for collateralPoolAddr, expectedTDaiBal := range tc.expectedCollateralPoolTDaiBalances {
				TDaiBal := bankKeeper.GetBalance(
					ctx,
					sdk.MustAccAddressFromBech32(collateralPoolAddr),
					asstypes.AssetTDai.Denom,
				)
				require.Equal(t,
					sdk.NewCoin(asstypes.AssetTDai.Denom, sdkmath.NewInt(expectedTDaiBal)),
					TDaiBal,
				)
			}

			if tc.expectedTDaiYieldPoolBalance != nil {
				TDaiBal := bankKeeper.GetBalance(
					ctx,
					accountKeeper.GetModuleAddress(ratelimittypes.TDaiPoolAccount),
					asstypes.AssetTDai.Denom,
				)
				require.Equal(t,
					sdk.NewCoin(asstypes.AssetTDai.Denom, sdkmath.NewIntFromBigInt(tc.expectedTDaiYieldPoolBalance)),
					TDaiBal,
				)

			}

			for _, perp := range tc.perpetuals {
				gotPerp, err := perpetualsKeeper.GetPerpetual(ctx, perp.GetId())
				require.NoError(t, err)

				if expectedOI, exists := tc.expectedOpenInterest[perp.GetId()]; exists {
					require.Equal(t, expectedOI, gotPerp.OpenInterest.BigInt())
				} else {
					// If no specified expected OI, then check OI is unchanged.
					require.Zero(t, perp.OpenInterest.BigInt().Cmp(
						gotPerp.OpenInterest.BigInt(),
					))
				}
			}
		})
	}
}

func TestUpdateSubaccounts_WithdrawalsBlocked(t *testing.T) {
	// default subaccount id, the first subaccount id generated when calling createNSubaccount
	firstSubaccountId := types.SubaccountId{
		Owner:  "0",
		Number: 0,
	}
	secondSubaccountId := types.SubaccountId{
		Owner:  "1",
		Number: 1,
	}

	tests := map[string]struct {
		// state
		perpetuals        []perptypes.Perpetual
		newFundingIndices []*big.Int // 1:1 mapped to perpetuals list
		assets            []*asstypes.Asset
		marketParamPrices []pricestypes.MarketParamPrice

		// subaccount state
		perpetualPositions map[types.SubaccountId][]*types.PerpetualPosition
		assetPositions     map[types.SubaccountId][]*types.AssetPosition

		// updates
		updates []types.Update

		// expectations
		expectedQuoteBalance       *big.Int
		expectedPerpetualPositions map[types.SubaccountId][]*types.PerpetualPosition
		expectedAssetPositions     map[types.SubaccountId][]*types.AssetPosition
		expectedSuccess            bool
		expectedSuccessPerUpdate   []types.UpdateResult
		expectedErr                error

		// Only contains the updated perpetual positions, to assert against the events included.
		expectedUpdatedPerpetualPositions     map[types.SubaccountId][]*types.PerpetualPosition
		expectedSubaccoundIdToFundingPayments map[types.SubaccountId]map[uint32]dtypes.SerializableInt
		expectedUpdatedAssetPositions         map[types.SubaccountId][]*types.AssetPosition
		expectedAssetYieldIndex               string
		msgSenderEnabled                      bool

		// Negative TNC subaccount state
		currentBlock                     uint32
		negativeTncSubaccountSeenAtBlock map[uint32]uint32

		// Update type
		updateType types.UpdateType
	}{
		"deposits are not blocked if negative TNC subaccount was seen at current block": {
			expectedQuoteBalance:     big.NewInt(100),
			expectedSuccess:          true,
			expectedSuccessPerUpdate: []types.UpdateResult{types.Success},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_SmallMarginRequirement,
			},
			updates: []types.Update{
				{
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(100)),
				},
			},
			expectedAssetPositions: map[types.SubaccountId][]*types.AssetPosition{
				firstSubaccountId: {
					{
						AssetId:  uint32(0),
						Quantums: dtypes.NewInt(100), // 100 TDai
					},
				},
			},
			expectedUpdatedAssetPositions: map[types.SubaccountId][]*types.AssetPosition{
				firstSubaccountId: {
					{
						AssetId:  uint32(0),
						Quantums: dtypes.NewInt(100), // 100 TDai
					},
				},
			},
			expectedAssetYieldIndex: constants.AssetYieldIndex_Zero,
			msgSenderEnabled:        true,

			currentBlock: 100,
			negativeTncSubaccountSeenAtBlock: map[uint32]uint32{
				constants.BtcUsd_NoMarginRequirement.Params.Id: 100,
			},

			updateType: types.Deposit,
		},
		`deposits are not blocked if current block is within
			WITHDRAWAL_AND_TRANSFERS_BLOCKED_AFTER_NEGATIVE_TNC_SUBACCOUNT_SEEN_BLOCKS`: {
			expectedQuoteBalance:     big.NewInt(100),
			expectedSuccess:          true,
			expectedSuccessPerUpdate: []types.UpdateResult{types.Success},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_SmallMarginRequirement,
			},
			updates: []types.Update{
				{
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(100)),
				},
			},
			expectedAssetPositions: map[types.SubaccountId][]*types.AssetPosition{
				firstSubaccountId: {
					{
						AssetId:  uint32(0),
						Quantums: dtypes.NewInt(100), // 100 TDai
					},
				},
			},
			expectedUpdatedAssetPositions: map[types.SubaccountId][]*types.AssetPosition{
				firstSubaccountId: {
					{
						AssetId:  uint32(0),
						Quantums: dtypes.NewInt(100), // 100 TDai
					},
				},
			},
			expectedAssetYieldIndex: constants.AssetYieldIndex_Zero,
			msgSenderEnabled:        true,

			currentBlock: 100,
			negativeTncSubaccountSeenAtBlock: map[uint32]uint32{
				constants.BtcUsd_NoMarginRequirement.Params.Id: 100 -
					types.WITHDRAWAL_AND_TRANSFERS_BLOCKED_AFTER_NEGATIVE_TNC_SUBACCOUNT_SEEN_BLOCKS + 1,
			},

			updateType: types.Deposit,
		},
		"deposits are not blocked if negative TNC subaccount was never seen": {
			expectedQuoteBalance:     big.NewInt(100),
			expectedSuccess:          true,
			expectedSuccessPerUpdate: []types.UpdateResult{types.Success},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_SmallMarginRequirement,
			},
			updates: []types.Update{
				{
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(100)),
				},
			},
			expectedAssetPositions: map[types.SubaccountId][]*types.AssetPosition{
				firstSubaccountId: {
					{
						AssetId:  uint32(0),
						Quantums: dtypes.NewInt(100), // 100 TDai
					},
				},
			},
			expectedUpdatedAssetPositions: map[types.SubaccountId][]*types.AssetPosition{
				firstSubaccountId: {
					{
						AssetId:  uint32(0),
						Quantums: dtypes.NewInt(100), // 100 TDai
					},
				},
			},
			expectedAssetYieldIndex: constants.AssetYieldIndex_Zero,
			msgSenderEnabled:        true,

			currentBlock: 100,
			negativeTncSubaccountSeenAtBlock: map[uint32]uint32{
				constants.BtcUsd_NoMarginRequirement.Params.Id: 0,
			},

			updateType: types.Deposit,
		},
		"withdrawals are blocked if negative TNC subaccount was seen at current block": {
			expectedQuoteBalance:     big.NewInt(-100),
			expectedSuccess:          false,
			expectedSuccessPerUpdate: []types.UpdateResult{types.WithdrawalsAndTransfersBlocked},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_SmallMarginRequirement,
			},
			perpetualPositions: map[types.SubaccountId][]*types.PerpetualPosition{
				firstSubaccountId: {&constants.PerpetualPosition_OneBTCLong},
			},
			expectedPerpetualPositions: map[types.SubaccountId][]*types.PerpetualPosition{
				firstSubaccountId: {&constants.PerpetualPosition_OneBTCLong},
			},
			expectedAssetPositions: map[types.SubaccountId][]*types.AssetPosition{},
			expectedUpdatedAssetPositions: map[types.SubaccountId][]*types.AssetPosition{
				firstSubaccountId: {},
			},
			updates: []types.Update{
				{
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(-100)),
				},
			},
			expectedAssetYieldIndex: constants.AssetYieldIndex_Zero,
			msgSenderEnabled:        true,

			currentBlock: 100,
			negativeTncSubaccountSeenAtBlock: map[uint32]uint32{
				constants.BtcUsd_NoMarginRequirement.Params.Id: 100,
			},

			updateType: types.Withdrawal,
		},
		`withdrawals are blocked if negative TNC subaccount was seen within
			WITHDRAWAL_AND_TRANSFERS_BLOCKED_AFTER_NEGATIVE_TNC_SUBACCOUNT_SEEN_BLOCKS`: {
			expectedQuoteBalance:     big.NewInt(-100),
			expectedSuccess:          false,
			expectedSuccessPerUpdate: []types.UpdateResult{types.WithdrawalsAndTransfersBlocked},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_SmallMarginRequirement,
			},
			perpetualPositions: map[types.SubaccountId][]*types.PerpetualPosition{
				firstSubaccountId: {&constants.PerpetualPosition_OneBTCLong},
			},
			expectedPerpetualPositions: map[types.SubaccountId][]*types.PerpetualPosition{
				firstSubaccountId: {&constants.PerpetualPosition_OneBTCLong},
			},
			expectedAssetPositions: map[types.SubaccountId][]*types.AssetPosition{},
			expectedUpdatedAssetPositions: map[types.SubaccountId][]*types.AssetPosition{
				firstSubaccountId: {},
			},
			updates: []types.Update{
				{
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(-100)),
				},
			},
			expectedAssetYieldIndex: constants.AssetYieldIndex_Zero,
			msgSenderEnabled:        true,

			currentBlock: 100,
			negativeTncSubaccountSeenAtBlock: map[uint32]uint32{
				constants.BtcUsd_NoMarginRequirement.Params.Id: 100 -
					types.WITHDRAWAL_AND_TRANSFERS_BLOCKED_AFTER_NEGATIVE_TNC_SUBACCOUNT_SEEN_BLOCKS + 1,
			},

			updateType: types.Withdrawal,
		},
		`withdrawals are not blocked if negative TNC subaccount was seen after
			WITHDRAWAL_AND_TRANSFERS_BLOCKED_AFTER_NEGATIVE_TNC_SUBACCOUNT_SEEN_BLOCKS`: {
			expectedQuoteBalance:     big.NewInt(-100),
			expectedSuccess:          true,
			expectedSuccessPerUpdate: []types.UpdateResult{types.Success},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_SmallMarginRequirement,
			},
			perpetualPositions: map[types.SubaccountId][]*types.PerpetualPosition{
				firstSubaccountId: {&constants.PerpetualPosition_OneBTCLong},
			},
			expectedPerpetualPositions: map[types.SubaccountId][]*types.PerpetualPosition{
				firstSubaccountId: {&constants.PerpetualPosition_OneBTCLong},
			},
			expectedAssetPositions: map[types.SubaccountId][]*types.AssetPosition{
				firstSubaccountId: {
					{
						AssetId:  uint32(0),
						Quantums: dtypes.NewInt(-100), // 100 TDai
					},
				},
			},
			expectedUpdatedAssetPositions: map[types.SubaccountId][]*types.AssetPosition{
				firstSubaccountId: {
					{
						AssetId:  uint32(0),
						Quantums: dtypes.NewInt(-100), // 100 TDai
					},
				},
			},
			updates: []types.Update{
				{
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(-100)),
				},
			},
			expectedAssetYieldIndex: constants.AssetYieldIndex_Zero,
			msgSenderEnabled:        true,

			currentBlock: 100,
			negativeTncSubaccountSeenAtBlock: map[uint32]uint32{
				constants.BtcUsd_NoMarginRequirement.Params.Id: 100 -
					types.WITHDRAWAL_AND_TRANSFERS_BLOCKED_AFTER_NEGATIVE_TNC_SUBACCOUNT_SEEN_BLOCKS,
			},

			updateType: types.Withdrawal,
		},
		"withdrawals are not blocked if negative TNC subaccount was never seen": {
			expectedQuoteBalance:     big.NewInt(-100),
			expectedSuccess:          true,
			expectedSuccessPerUpdate: []types.UpdateResult{types.Success},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_SmallMarginRequirement,
			},
			perpetualPositions: map[types.SubaccountId][]*types.PerpetualPosition{
				firstSubaccountId: {&constants.PerpetualPosition_OneBTCLong},
			},
			expectedPerpetualPositions: map[types.SubaccountId][]*types.PerpetualPosition{
				firstSubaccountId: {&constants.PerpetualPosition_OneBTCLong},
			},
			expectedAssetPositions: map[types.SubaccountId][]*types.AssetPosition{
				firstSubaccountId: {
					{
						AssetId:  uint32(0),
						Quantums: dtypes.NewInt(-100), // 100 TDai
					},
				},
			},
			expectedUpdatedAssetPositions: map[types.SubaccountId][]*types.AssetPosition{
				firstSubaccountId: {
					{
						AssetId:  uint32(0),
						Quantums: dtypes.NewInt(-100), // 100 TDai
					},
				},
			},
			updates: []types.Update{
				{
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(-100)),
				},
			},
			expectedAssetYieldIndex: constants.AssetYieldIndex_Zero,
			msgSenderEnabled:        true,

			currentBlock: 100,
			negativeTncSubaccountSeenAtBlock: map[uint32]uint32{
				constants.BtcUsd_NoMarginRequirement.Params.Id: 0,
			},

			updateType: types.Withdrawal,
		},
		`withdrawals are not blocked if negative TNC subaccount was seen within
		WITHDRAWAL_AND_TRANSFERS_BLOCKED_AFTER_NEGATIVE_TNC_SUBACCOUNT_SEEN_BLOCKS for a different
		collateral pool`: {
			expectedQuoteBalance:     big.NewInt(-100),
			expectedSuccess:          true,
			expectedSuccessPerUpdate: []types.UpdateResult{types.Success},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_SmallMarginRequirement,
				constants.IsoUsd_IsolatedMarket,
			},
			perpetualPositions: map[types.SubaccountId][]*types.PerpetualPosition{
				firstSubaccountId: {&constants.PerpetualPosition_OneBTCLong},
			},
			expectedPerpetualPositions: map[types.SubaccountId][]*types.PerpetualPosition{
				firstSubaccountId: {&constants.PerpetualPosition_OneBTCLong},
			},
			expectedAssetPositions: map[types.SubaccountId][]*types.AssetPosition{
				firstSubaccountId: {
					{
						AssetId:  uint32(0),
						Quantums: dtypes.NewInt(-100), // 100 TDai
					},
				},
			},
			expectedUpdatedAssetPositions: map[types.SubaccountId][]*types.AssetPosition{
				firstSubaccountId: {
					{
						AssetId:  uint32(0),
						Quantums: dtypes.NewInt(-100), // 100 TDai
					},
				},
			},
			updates: []types.Update{
				{
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(-100)),
				},
			},
			expectedAssetYieldIndex: constants.AssetYieldIndex_Zero,
			msgSenderEnabled:        true,

			currentBlock: 100,
			negativeTncSubaccountSeenAtBlock: map[uint32]uint32{
				constants.IsoUsd_IsolatedMarket.Params.Id: 100 -
					types.WITHDRAWAL_AND_TRANSFERS_BLOCKED_AFTER_NEGATIVE_TNC_SUBACCOUNT_SEEN_BLOCKS + 1,
			},

			updateType: types.Withdrawal,
		},
		`withdrawals are blocked if negative TNC subaccount was seen within
		WITHDRAWAL_AND_TRANSFERS_BLOCKED_AFTER_NEGATIVE_TNC_SUBACCOUNT_SEEN_BLOCKS for an isolated
		perpetual collateral pool`: {
			expectedQuoteBalance:     big.NewInt(-100),
			expectedSuccess:          false,
			expectedSuccessPerUpdate: []types.UpdateResult{types.WithdrawalsAndTransfersBlocked},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_SmallMarginRequirement,
				constants.IsoUsd_IsolatedMarket,
			},
			perpetualPositions: map[types.SubaccountId][]*types.PerpetualPosition{
				firstSubaccountId: {&constants.PerpetualPosition_OneISOLong},
			},
			expectedPerpetualPositions: map[types.SubaccountId][]*types.PerpetualPosition{
				firstSubaccountId: {&constants.PerpetualPosition_OneISOLong},
			},
			expectedAssetPositions: map[types.SubaccountId][]*types.AssetPosition{},
			expectedUpdatedAssetPositions: map[types.SubaccountId][]*types.AssetPosition{
				firstSubaccountId: {},
			},
			updates: []types.Update{
				{
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(-100)),
				},
			},
			expectedAssetYieldIndex: constants.AssetYieldIndex_Zero,
			msgSenderEnabled:        true,

			currentBlock: 100,
			negativeTncSubaccountSeenAtBlock: map[uint32]uint32{
				constants.IsoUsd_IsolatedMarket.Params.Id: 100 -
					types.WITHDRAWAL_AND_TRANSFERS_BLOCKED_AFTER_NEGATIVE_TNC_SUBACCOUNT_SEEN_BLOCKS + 1,
			},

			updateType: types.Withdrawal,
		},
		`withdrawals are blocked if negative TNC subaccount was seen within
		WITHDRAWAL_AND_TRANSFERS_BLOCKED_AFTER_NEGATIVE_TNC_SUBACCOUNT_SEEN_BLOCKS for one isolated
		perpetual collateral pool and negative TNC subaccount was never seen for the cross-perpetual
		collateral pool, both of which are associated with subaccounts being updated`: {
			expectedQuoteBalance: big.NewInt(-100),
			expectedSuccess:      false,
			expectedSuccessPerUpdate: []types.UpdateResult{
				types.WithdrawalsAndTransfersBlocked,
				types.WithdrawalsAndTransfersBlocked,
			},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_SmallMarginRequirement,
				constants.IsoUsd_IsolatedMarket,
			},
			perpetualPositions: map[types.SubaccountId][]*types.PerpetualPosition{
				firstSubaccountId: {&constants.PerpetualPosition_OneISOLong},
			},
			expectedPerpetualPositions: map[types.SubaccountId][]*types.PerpetualPosition{
				firstSubaccountId: {&constants.PerpetualPosition_OneISOLong},
			},
			expectedAssetPositions: map[types.SubaccountId][]*types.AssetPosition{},
			expectedUpdatedAssetPositions: map[types.SubaccountId][]*types.AssetPosition{
				firstSubaccountId:  {},
				secondSubaccountId: {},
			},
			updates: []types.Update{
				{
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(-100)),
				},
				{
					SubaccountId: secondSubaccountId,
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(-100)),
				},
			},
			expectedAssetYieldIndex: constants.AssetYieldIndex_Zero,
			msgSenderEnabled:        true,

			currentBlock: 100,
			negativeTncSubaccountSeenAtBlock: map[uint32]uint32{
				constants.IsoUsd_IsolatedMarket.Params.Id: 100 -
					types.WITHDRAWAL_AND_TRANSFERS_BLOCKED_AFTER_NEGATIVE_TNC_SUBACCOUNT_SEEN_BLOCKS + 1,
				constants.BtcUsd_NoMarginRequirement.Params.Id: 0,
			},

			updateType: types.Withdrawal,
		},
		`withdrawals are blocked if negative TNC subaccount was seen within
		WITHDRAWAL_AND_TRANSFERS_BLOCKED_AFTER_NEGATIVE_TNC_SUBACCOUNT_SEEN_BLOCKS for one isolated
		perpetual collateral pool and negative TNC subaccount was seen for the cross-perpetual
		collateral pool after WITHDRAWAL_AND_TRANSFERS_BLOCKED_AFTER_NEGATIVE_TNC_SUBACCOUNT_SEEN_BLOCKS,
		both of which are associated with subaccounts being updated`: {
			expectedQuoteBalance: big.NewInt(-100),
			expectedSuccess:      false,
			expectedSuccessPerUpdate: []types.UpdateResult{
				types.WithdrawalsAndTransfersBlocked,
				types.WithdrawalsAndTransfersBlocked,
			},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_SmallMarginRequirement,
				constants.IsoUsd_IsolatedMarket,
			},
			perpetualPositions: map[types.SubaccountId][]*types.PerpetualPosition{
				firstSubaccountId: {&constants.PerpetualPosition_OneISOLong},
			},
			expectedPerpetualPositions: map[types.SubaccountId][]*types.PerpetualPosition{
				firstSubaccountId: {&constants.PerpetualPosition_OneISOLong},
			},
			expectedAssetPositions: map[types.SubaccountId][]*types.AssetPosition{},
			expectedUpdatedAssetPositions: map[types.SubaccountId][]*types.AssetPosition{
				firstSubaccountId:  {},
				secondSubaccountId: {},
			},
			updates: []types.Update{
				{
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(-100)),
				},
				{
					SubaccountId: secondSubaccountId,
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(-100)),
				},
			},
			expectedAssetYieldIndex: constants.AssetYieldIndex_Zero,
			msgSenderEnabled:        true,

			currentBlock: 100,
			negativeTncSubaccountSeenAtBlock: map[uint32]uint32{
				constants.IsoUsd_IsolatedMarket.Params.Id: 100 -
					types.WITHDRAWAL_AND_TRANSFERS_BLOCKED_AFTER_NEGATIVE_TNC_SUBACCOUNT_SEEN_BLOCKS + 1,
				constants.BtcUsd_NoMarginRequirement.Params.Id: 100 -
					types.WITHDRAWAL_AND_TRANSFERS_BLOCKED_AFTER_NEGATIVE_TNC_SUBACCOUNT_SEEN_BLOCKS,
			},

			updateType: types.Withdrawal,
		},
		`withdrawals are blocked if negative TNC subaccount was seen within
		WITHDRAWAL_AND_TRANSFERS_BLOCKED_AFTER_NEGATIVE_TNC_SUBACCOUNT_SEEN_BLOCKS for one isolated
		perpetual collateral pool and negative TNC subaccount was never seen for another isolated
		collateral pool, both of which are associated with subaccounts being updated`: {
			expectedQuoteBalance: big.NewInt(-100),
			expectedSuccess:      false,
			expectedSuccessPerUpdate: []types.UpdateResult{
				types.WithdrawalsAndTransfersBlocked,
				types.WithdrawalsAndTransfersBlocked,
			},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_SmallMarginRequirement,
				constants.IsoUsd_IsolatedMarket,
				constants.Iso2Usd_IsolatedMarket,
			},
			perpetualPositions: map[types.SubaccountId][]*types.PerpetualPosition{
				firstSubaccountId:  {&constants.PerpetualPosition_OneISOLong},
				secondSubaccountId: {&constants.PerpetualPosition_OneISO2Long},
			},
			expectedPerpetualPositions: map[types.SubaccountId][]*types.PerpetualPosition{
				firstSubaccountId:  {&constants.PerpetualPosition_OneISOLong},
				secondSubaccountId: {&constants.PerpetualPosition_OneISO2Long},
			},
			expectedAssetPositions: map[types.SubaccountId][]*types.AssetPosition{},
			expectedUpdatedAssetPositions: map[types.SubaccountId][]*types.AssetPosition{
				firstSubaccountId:  {},
				secondSubaccountId: {},
			},
			updates: []types.Update{
				{
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(-100)),
				},
				{
					SubaccountId: secondSubaccountId,
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(-100)),
				},
			},
			expectedAssetYieldIndex: constants.AssetYieldIndex_Zero,
			msgSenderEnabled:        true,

			currentBlock: 100,
			negativeTncSubaccountSeenAtBlock: map[uint32]uint32{
				constants.IsoUsd_IsolatedMarket.Params.Id: 0,
				constants.Iso2Usd_IsolatedMarket.Params.Id: 100 -
					types.WITHDRAWAL_AND_TRANSFERS_BLOCKED_AFTER_NEGATIVE_TNC_SUBACCOUNT_SEEN_BLOCKS + 1,
			},

			updateType: types.Withdrawal,
		},
		`withdrawals are blocked if negative TNC subaccount was seen within
		WITHDRAWAL_AND_TRANSFERS_BLOCKED_AFTER_NEGATIVE_TNC_SUBACCOUNT_SEEN_BLOCKS for one isolated
		perpetual collateral pool and negative TNC subaccount was seen for another isolated perpetual
		collateral pool after WITHDRAWAL_AND_TRANSFERS_BLOCKED_AFTER_NEGATIVE_TNC_SUBACCOUNT_SEEN_BLOCKS,
		both of which are associated with subaccounts being updated`: {
			expectedQuoteBalance: big.NewInt(-100),
			expectedSuccess:      false,
			expectedSuccessPerUpdate: []types.UpdateResult{
				types.WithdrawalsAndTransfersBlocked,
				types.WithdrawalsAndTransfersBlocked,
			},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_SmallMarginRequirement,
				constants.IsoUsd_IsolatedMarket,
				constants.Iso2Usd_IsolatedMarket,
			},
			perpetualPositions: map[types.SubaccountId][]*types.PerpetualPosition{
				firstSubaccountId:  {&constants.PerpetualPosition_OneISOLong},
				secondSubaccountId: {&constants.PerpetualPosition_OneISO2Long},
			},
			expectedPerpetualPositions: map[types.SubaccountId][]*types.PerpetualPosition{
				firstSubaccountId:  {&constants.PerpetualPosition_OneISOLong},
				secondSubaccountId: {&constants.PerpetualPosition_OneISO2Long},
			},
			expectedAssetPositions: map[types.SubaccountId][]*types.AssetPosition{},
			expectedUpdatedAssetPositions: map[types.SubaccountId][]*types.AssetPosition{
				firstSubaccountId:  {},
				secondSubaccountId: {},
			},
			updates: []types.Update{
				{
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(-100)),
				},
				{
					SubaccountId: secondSubaccountId,
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(-100)),
				},
			},
			expectedAssetYieldIndex: constants.AssetYieldIndex_Zero,
			msgSenderEnabled:        true,

			currentBlock: 100,
			negativeTncSubaccountSeenAtBlock: map[uint32]uint32{
				constants.IsoUsd_IsolatedMarket.Params.Id: 100 -
					types.WITHDRAWAL_AND_TRANSFERS_BLOCKED_AFTER_NEGATIVE_TNC_SUBACCOUNT_SEEN_BLOCKS,
				constants.Iso2Usd_IsolatedMarket.Params.Id: 100 -
					types.WITHDRAWAL_AND_TRANSFERS_BLOCKED_AFTER_NEGATIVE_TNC_SUBACCOUNT_SEEN_BLOCKS + 1,
			},

			updateType: types.Withdrawal,
		},
		"well-collateralized matches are not blocked if negative TNC subaccount was seen at current block": {
			assetPositions: map[types.SubaccountId][]*types.AssetPosition{
				firstSubaccountId: testutil.CreateTDaiAssetPosition(big.NewInt(25_000_000_000)), // $25,000
			},
			expectedQuoteBalance:     big.NewInt(0),
			expectedSuccess:          true,
			expectedSuccessPerUpdate: []types.UpdateResult{types.Success, types.Success},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_NoMarginRequirement,
			},
			perpetualPositions: map[types.SubaccountId][]*types.PerpetualPosition{
				firstSubaccountId: {&constants.PerpetualPosition_OneBTCLong},
			},
			expectedPerpetualPositions: map[types.SubaccountId][]*types.PerpetualPosition{
				firstSubaccountId: {&constants.PerpetualPosition_OneAndHalfBTCLong},
			},
			expectedUpdatedPerpetualPositions: map[types.SubaccountId][]*types.PerpetualPosition{
				firstSubaccountId: {
					{
						PerpetualId:  uint32(0),
						Quantums:     dtypes.NewInt(150_000_000), // 1.5 BTC
						FundingIndex: dtypes.NewInt(0),
					},
				},
			},
			expectedAssetPositions: map[types.SubaccountId][]*types.AssetPosition{},
			expectedUpdatedAssetPositions: map[types.SubaccountId][]*types.AssetPosition{
				firstSubaccountId: {
					{
						AssetId:  uint32(0),
						Quantums: dtypes.NewInt(0),
					},
				},
			},
			updates: []types.Update{
				{
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(-25_000_000_000)), // -$25,000
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(0),
							BigQuantumsDelta: big.NewInt(50_000_000), // .5 BTC
						},
					},
				},
				{
					SubaccountId: secondSubaccountId,
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(25_000_000_000)), // $25,000
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(0),
							BigQuantumsDelta: big.NewInt(-50_000_000), // .5 BTC
						},
					},
				},
			},
			msgSenderEnabled: false,

			currentBlock: 100,
			negativeTncSubaccountSeenAtBlock: map[uint32]uint32{
				constants.BtcUsd_NoMarginRequirement.Params.Id: 100,
			},

			updateType: types.Match,
		},
		`well-collateralized matches are not blocked if current block is within
			WITHDRAWAL_AND_TRANSFERS_BLOCKED_AFTER_NEGATIVE_TNC_SUBACCOUNT_SEEN_BLOCKS`: {
			assetPositions: map[types.SubaccountId][]*types.AssetPosition{
				firstSubaccountId: testutil.CreateTDaiAssetPosition(big.NewInt(25_000_000_000)), // $25,000
			},
			expectedQuoteBalance:     big.NewInt(0),
			expectedSuccess:          true,
			expectedSuccessPerUpdate: []types.UpdateResult{types.Success, types.Success},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_NoMarginRequirement,
			},
			perpetualPositions: map[types.SubaccountId][]*types.PerpetualPosition{
				firstSubaccountId: {&constants.PerpetualPosition_OneBTCLong},
			},
			expectedPerpetualPositions: map[types.SubaccountId][]*types.PerpetualPosition{
				firstSubaccountId: {&constants.PerpetualPosition_OneAndHalfBTCLong},
			},
			expectedUpdatedPerpetualPositions: map[types.SubaccountId][]*types.PerpetualPosition{
				firstSubaccountId: {
					{
						PerpetualId:  uint32(0),
						Quantums:     dtypes.NewInt(150_000_000), // 1.5 BTC
						FundingIndex: dtypes.NewInt(0),
					},
				},
			},
			expectedAssetPositions: map[types.SubaccountId][]*types.AssetPosition{},
			expectedUpdatedAssetPositions: map[types.SubaccountId][]*types.AssetPosition{
				firstSubaccountId: {
					{
						AssetId:  uint32(0),
						Quantums: dtypes.NewInt(0),
					},
				},
			},
			updates: []types.Update{
				{
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(-25_000_000_000)), // -$25,000
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(0),
							BigQuantumsDelta: big.NewInt(50_000_000), // .5 BTC
						},
					},
				},
				{
					SubaccountId: secondSubaccountId,
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(25_000_000_000)), // $25,000
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(0),
							BigQuantumsDelta: big.NewInt(-50_000_000), // .5 BTC
						},
					},
				},
			},
			msgSenderEnabled: false,

			currentBlock: 100,
			negativeTncSubaccountSeenAtBlock: map[uint32]uint32{
				constants.BtcUsd_NoMarginRequirement.Params.Id: 100 -
					types.WITHDRAWAL_AND_TRANSFERS_BLOCKED_AFTER_NEGATIVE_TNC_SUBACCOUNT_SEEN_BLOCKS + 1,
			},

			updateType: types.Match,
		},
		"well-collateralized matches are not blocked if negative TNC subaccount was never seen": {
			assetPositions: map[types.SubaccountId][]*types.AssetPosition{
				firstSubaccountId: testutil.CreateTDaiAssetPosition(big.NewInt(25_000_000_000)), // $25,000
			},
			expectedQuoteBalance:     big.NewInt(0),
			expectedSuccess:          true,
			expectedSuccessPerUpdate: []types.UpdateResult{types.Success, types.Success},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_NoMarginRequirement,
			},
			perpetualPositions: map[types.SubaccountId][]*types.PerpetualPosition{
				firstSubaccountId: {&constants.PerpetualPosition_OneBTCLong},
			},
			expectedPerpetualPositions: map[types.SubaccountId][]*types.PerpetualPosition{
				firstSubaccountId: {&constants.PerpetualPosition_OneAndHalfBTCLong},
			},
			expectedUpdatedPerpetualPositions: map[types.SubaccountId][]*types.PerpetualPosition{
				firstSubaccountId: {
					{
						PerpetualId:  uint32(0),
						Quantums:     dtypes.NewInt(150_000_000), // 1.5 BTC
						FundingIndex: dtypes.NewInt(0),
					},
				},
			},
			expectedAssetPositions: map[types.SubaccountId][]*types.AssetPosition{},
			expectedUpdatedAssetPositions: map[types.SubaccountId][]*types.AssetPosition{
				firstSubaccountId: {
					{
						AssetId:  uint32(0),
						Quantums: dtypes.NewInt(0),
					},
				},
			},
			updates: []types.Update{
				{
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(-25_000_000_000)), // -$25,000
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(0),
							BigQuantumsDelta: big.NewInt(50_000_000), // .5 BTC
						},
					},
				},
				{
					SubaccountId: secondSubaccountId,
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(25_000_000_000)), // $25,000
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(0),
							BigQuantumsDelta: big.NewInt(-50_000_000), // .5 BTC
						},
					},
				},
			},
			msgSenderEnabled: false,

			currentBlock: 100,
			negativeTncSubaccountSeenAtBlock: map[uint32]uint32{
				constants.BtcUsd_NoMarginRequirement.Params.Id: 0,
			},

			updateType: types.Match,
		},
		"undercollateralized matches are not blocked if negative TNC subaccount was seen at current block": {
			expectedQuoteBalance: big.NewInt(0),
			expectedSuccess:      false,
			expectedSuccessPerUpdate: []types.UpdateResult{
				types.NewlyUndercollateralized,
				types.NewlyUndercollateralized,
			},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_SmallMarginRequirement,
			},
			perpetualPositions: map[types.SubaccountId][]*types.PerpetualPosition{
				firstSubaccountId: {&constants.PerpetualPosition_OneHundredthBTCLong},
			},
			expectedPerpetualPositions: map[types.SubaccountId][]*types.PerpetualPosition{
				firstSubaccountId: {&constants.PerpetualPosition_OneHundredthBTCLong},
			},
			updates: []types.Update{
				{
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(-50_000_000_000)), // -$50,000
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(0),
							BigQuantumsDelta: big.NewInt(100_000_000), // 1 BTC
						},
					},
				},
				{
					SubaccountId: secondSubaccountId,
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(50_000_000_000)), // $50,000
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(0),
							BigQuantumsDelta: big.NewInt(-100_000_000), // -1 BTC
						},
					},
				},
			},
			expectedAssetYieldIndex: constants.AssetYieldIndex_Zero,
			msgSenderEnabled:        true,

			currentBlock: 100,
			negativeTncSubaccountSeenAtBlock: map[uint32]uint32{
				constants.BtcUsd_NoMarginRequirement.Params.Id: 100,
			},

			updateType: types.Match,
		},
		`undercollateralized matches are not blocked if current block is within
			WITHDRAWAL_AND_TRANSFERS_BLOCKED_AFTER_NEGATIVE_TNC_SUBACCOUNT_SEEN_BLOCKS`: {
			expectedQuoteBalance: big.NewInt(0),
			expectedSuccess:      false,
			expectedSuccessPerUpdate: []types.UpdateResult{
				types.NewlyUndercollateralized,
				types.NewlyUndercollateralized,
			},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_SmallMarginRequirement,
			},
			perpetualPositions: map[types.SubaccountId][]*types.PerpetualPosition{
				firstSubaccountId: {&constants.PerpetualPosition_OneHundredthBTCLong},
			},
			expectedPerpetualPositions: map[types.SubaccountId][]*types.PerpetualPosition{
				firstSubaccountId: {&constants.PerpetualPosition_OneHundredthBTCLong},
			},
			updates: []types.Update{
				{
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(-50_000_000_000)), // -$50,000
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(0),
							BigQuantumsDelta: big.NewInt(100_000_000), // 1 BTC
						},
					},
				},
				{
					SubaccountId: secondSubaccountId,
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(50_000_000_000)), // $50,000
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(0),
							BigQuantumsDelta: big.NewInt(-100_000_000), // 1 BTC
						},
					},
				},
			},
			expectedAssetYieldIndex: constants.AssetYieldIndex_Zero,
			msgSenderEnabled:        true,

			currentBlock: 100,
			negativeTncSubaccountSeenAtBlock: map[uint32]uint32{
				constants.BtcUsd_NoMarginRequirement.Params.Id: 100 -
					types.WITHDRAWAL_AND_TRANSFERS_BLOCKED_AFTER_NEGATIVE_TNC_SUBACCOUNT_SEEN_BLOCKS + 1,
			},

			updateType: types.Match,
		},
		"undercollateralized matches are not blocked if negative TNC subaccount was never seen": {
			expectedQuoteBalance: big.NewInt(0),
			expectedSuccess:      false,
			expectedSuccessPerUpdate: []types.UpdateResult{
				types.NewlyUndercollateralized,
				types.NewlyUndercollateralized,
			},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_SmallMarginRequirement,
			},
			perpetualPositions: map[types.SubaccountId][]*types.PerpetualPosition{
				firstSubaccountId: {&constants.PerpetualPosition_OneHundredthBTCLong},
			},
			expectedPerpetualPositions: map[types.SubaccountId][]*types.PerpetualPosition{
				firstSubaccountId: {&constants.PerpetualPosition_OneHundredthBTCLong},
			},
			updates: []types.Update{
				{
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(-50_000_000_000)), // -$50,000
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(0),
							BigQuantumsDelta: big.NewInt(100_000_000), // 1 BTC
						},
					},
				},
				{
					SubaccountId: secondSubaccountId,
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(50_000_000_000)), // $50,000
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(0),
							BigQuantumsDelta: big.NewInt(-100_000_000), // -1 BTC
						},
					},
				},
			},
			expectedAssetYieldIndex: constants.AssetYieldIndex_Zero,
			msgSenderEnabled:        true,

			currentBlock: 100,
			negativeTncSubaccountSeenAtBlock: map[uint32]uint32{
				constants.BtcUsd_NoMarginRequirement.Params.Id: 0,
			},

			updateType: types.Match,
		},
		"transfers are blocked if negative TNC subaccount was seen at current block": {
			expectedQuoteBalance: big.NewInt(-100),
			expectedSuccess:      false,
			expectedSuccessPerUpdate: []types.UpdateResult{
				types.WithdrawalsAndTransfersBlocked,
				types.WithdrawalsAndTransfersBlocked,
			},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_SmallMarginRequirement,
			},
			perpetualPositions:         map[types.SubaccountId][]*types.PerpetualPosition{},
			expectedPerpetualPositions: map[types.SubaccountId][]*types.PerpetualPosition{},
			expectedAssetPositions:     map[types.SubaccountId][]*types.AssetPosition{},
			expectedUpdatedAssetPositions: map[types.SubaccountId][]*types.AssetPosition{
				firstSubaccountId:  {},
				secondSubaccountId: {},
			},
			updates: []types.Update{
				{
					SubaccountId: firstSubaccountId,
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(-100)),
				},
				{
					SubaccountId: secondSubaccountId,
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(100)),
				},
			},
			expectedAssetYieldIndex: constants.AssetYieldIndex_Zero,
			msgSenderEnabled:        true,

			currentBlock: 100,
			negativeTncSubaccountSeenAtBlock: map[uint32]uint32{
				constants.BtcUsd_NoMarginRequirement.Params.Id: 100,
			},

			updateType: types.Transfer,
		},
		`transfers are blocked if negative TNC subaccount was seen within
			WITHDRAWAL_AND_TRANSFERS_BLOCKED_AFTER_NEGATIVE_TNC_SUBACCOUNT_SEEN_BLOCKS`: {
			expectedQuoteBalance: big.NewInt(-100),
			expectedSuccess:      false,
			expectedSuccessPerUpdate: []types.UpdateResult{
				types.WithdrawalsAndTransfersBlocked,
				types.WithdrawalsAndTransfersBlocked,
			},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_SmallMarginRequirement,
			},
			perpetualPositions:         map[types.SubaccountId][]*types.PerpetualPosition{},
			expectedPerpetualPositions: map[types.SubaccountId][]*types.PerpetualPosition{},
			expectedAssetPositions:     map[types.SubaccountId][]*types.AssetPosition{},
			expectedUpdatedAssetPositions: map[types.SubaccountId][]*types.AssetPosition{
				firstSubaccountId:  {},
				secondSubaccountId: {},
			},
			updates: []types.Update{
				{
					SubaccountId: firstSubaccountId,
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(-100)),
				},
				{
					SubaccountId: secondSubaccountId,
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(100)),
				},
			},
			expectedAssetYieldIndex: constants.AssetYieldIndex_Zero,
			msgSenderEnabled:        true,

			currentBlock: 100,
			negativeTncSubaccountSeenAtBlock: map[uint32]uint32{
				constants.BtcUsd_NoMarginRequirement.Params.Id: 100 -
					types.WITHDRAWAL_AND_TRANSFERS_BLOCKED_AFTER_NEGATIVE_TNC_SUBACCOUNT_SEEN_BLOCKS + 1,
			},

			updateType: types.Transfer,
		},
		`transfers are not blocked if negative TNC subaccount was seen after
			WITHDRAWAL_AND_TRANSFERS_BLOCKED_AFTER_NEGATIVE_TNC_SUBACCOUNT_SEEN_BLOCKS`: {
			expectedQuoteBalance: big.NewInt(-100),
			expectedSuccess:      false,
			expectedSuccessPerUpdate: []types.UpdateResult{
				types.NewlyUndercollateralized,
				types.Success,
			},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_SmallMarginRequirement,
			},
			perpetualPositions:         map[types.SubaccountId][]*types.PerpetualPosition{},
			expectedPerpetualPositions: map[types.SubaccountId][]*types.PerpetualPosition{},
			expectedAssetPositions:     map[types.SubaccountId][]*types.AssetPosition{},
			expectedUpdatedAssetPositions: map[types.SubaccountId][]*types.AssetPosition{
				firstSubaccountId:  {},
				secondSubaccountId: {},
			},
			updates: []types.Update{
				{
					SubaccountId: firstSubaccountId,
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(-100)),
				},
				{
					SubaccountId: secondSubaccountId,
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(100)),
				},
			},
			expectedAssetYieldIndex: constants.AssetYieldIndex_Zero,
			msgSenderEnabled:        true,

			currentBlock: 100,
			negativeTncSubaccountSeenAtBlock: map[uint32]uint32{
				constants.BtcUsd_NoMarginRequirement.Params.Id: 100 -
					types.WITHDRAWAL_AND_TRANSFERS_BLOCKED_AFTER_NEGATIVE_TNC_SUBACCOUNT_SEEN_BLOCKS,
			},

			updateType: types.Transfer,
		},
		"transfers are not blocked if negative TNC subaccount was never seen": {
			expectedQuoteBalance: big.NewInt(-100),
			expectedSuccess:      false,
			expectedSuccessPerUpdate: []types.UpdateResult{
				types.NewlyUndercollateralized,
				types.Success,
			},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_SmallMarginRequirement,
			},
			perpetualPositions:         map[types.SubaccountId][]*types.PerpetualPosition{},
			expectedPerpetualPositions: map[types.SubaccountId][]*types.PerpetualPosition{},
			expectedAssetPositions:     map[types.SubaccountId][]*types.AssetPosition{},
			expectedUpdatedAssetPositions: map[types.SubaccountId][]*types.AssetPosition{
				firstSubaccountId:  {},
				secondSubaccountId: {},
			},
			updates: []types.Update{
				{
					SubaccountId: firstSubaccountId,
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(-100)),
				},
				{
					SubaccountId: secondSubaccountId,
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(100)),
				},
			},
			expectedAssetYieldIndex: constants.AssetYieldIndex_Zero,
			msgSenderEnabled:        true,

			currentBlock: 100,
			negativeTncSubaccountSeenAtBlock: map[uint32]uint32{
				constants.BtcUsd_NoMarginRequirement.Params.Id: 0,
			},

			updateType: types.Transfer,
		},
		`transfers are not blocked if negative TNC subaccount was seen within
		WITHDRAWAL_AND_TRANSFERS_BLOCKED_AFTER_NEGATIVE_TNC_SUBACCOUNT_SEEN_BLOCKS for a different
		collateral pool from the ones associated with the subaccounts being updated`: {
			expectedQuoteBalance: big.NewInt(-100),
			expectedSuccess:      false,
			expectedSuccessPerUpdate: []types.UpdateResult{
				types.NewlyUndercollateralized,
				types.Success,
			},
			perpetuals: []perptypes.Perpetual{
				constants.IsoUsd_IsolatedMarket,
				constants.Iso2Usd_IsolatedMarket,
			},
			perpetualPositions: map[types.SubaccountId][]*types.PerpetualPosition{
				secondSubaccountId: {&constants.PerpetualPosition_OneISOLong},
			},
			expectedPerpetualPositions: map[types.SubaccountId][]*types.PerpetualPosition{
				secondSubaccountId: {&constants.PerpetualPosition_OneISOLong},
			},
			expectedAssetPositions: map[types.SubaccountId][]*types.AssetPosition{},
			expectedUpdatedAssetPositions: map[types.SubaccountId][]*types.AssetPosition{
				firstSubaccountId:  {},
				secondSubaccountId: {},
			},
			updates: []types.Update{
				{
					SubaccountId: firstSubaccountId,
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(-100)),
				},
				{
					SubaccountId: secondSubaccountId,
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(100)),
				},
			},
			expectedAssetYieldIndex: constants.AssetYieldIndex_Zero,
			msgSenderEnabled:        true,

			currentBlock: 100,
			negativeTncSubaccountSeenAtBlock: map[uint32]uint32{
				constants.Iso2Usd_IsolatedMarket.Params.Id: 100 -
					types.WITHDRAWAL_AND_TRANSFERS_BLOCKED_AFTER_NEGATIVE_TNC_SUBACCOUNT_SEEN_BLOCKS + 1,
			},

			updateType: types.Transfer,
		},
		`transfers are blocked if negative TNC subaccount was seen within
		WITHDRAWAL_AND_TRANSFERS_BLOCKED_AFTER_NEGATIVE_TNC_SUBACCOUNT_SEEN_BLOCKS for one isolated
		perpetual collateral pool and negative TNC subaccount was never seen for the cross-perpetual
		collateral pool, both of which are associated with subaccounts being updated`: {
			expectedQuoteBalance: big.NewInt(-100),
			expectedSuccess:      false,
			expectedSuccessPerUpdate: []types.UpdateResult{
				types.WithdrawalsAndTransfersBlocked,
				types.WithdrawalsAndTransfersBlocked,
			},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_SmallMarginRequirement,
				constants.IsoUsd_IsolatedMarket,
			},
			perpetualPositions: map[types.SubaccountId][]*types.PerpetualPosition{
				secondSubaccountId: {&constants.PerpetualPosition_OneISOLong},
			},
			expectedPerpetualPositions: map[types.SubaccountId][]*types.PerpetualPosition{
				secondSubaccountId: {&constants.PerpetualPosition_OneISOLong},
			},
			expectedAssetPositions: map[types.SubaccountId][]*types.AssetPosition{},
			expectedUpdatedAssetPositions: map[types.SubaccountId][]*types.AssetPosition{
				firstSubaccountId:  {},
				secondSubaccountId: {},
			},
			updates: []types.Update{
				{
					SubaccountId: firstSubaccountId,
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(-100)),
				},
				{
					SubaccountId: secondSubaccountId,
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(100)),
				},
			},
			expectedAssetYieldIndex: constants.AssetYieldIndex_Zero,
			msgSenderEnabled:        true,

			currentBlock: 100,
			negativeTncSubaccountSeenAtBlock: map[uint32]uint32{
				constants.IsoUsd_IsolatedMarket.Params.Id: 100 -
					types.WITHDRAWAL_AND_TRANSFERS_BLOCKED_AFTER_NEGATIVE_TNC_SUBACCOUNT_SEEN_BLOCKS + 1,
				constants.BtcUsd_NoMarginRequirement.Params.Id: 0,
			},

			updateType: types.Transfer,
		},
		`transfers are blocked if negative TNC subaccount was seen within
		WITHDRAWAL_AND_TRANSFERS_BLOCKED_AFTER_NEGATIVE_TNC_SUBACCOUNT_SEEN_BLOCKS for one isolated
		perpetual collateral pool and negative TNC subaccount was seen for the cross-perpetual
		collateral pool after WITHDRAWAL_AND_TRANSFERS_BLOCKED_AFTER_NEGATIVE_TNC_SUBACCOUNT_SEEN_BLOCKS,
		both of which are associated with subaccounts being updated`: {
			expectedQuoteBalance: big.NewInt(-100),
			expectedSuccess:      false,
			expectedSuccessPerUpdate: []types.UpdateResult{
				types.WithdrawalsAndTransfersBlocked,
				types.WithdrawalsAndTransfersBlocked,
			},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_SmallMarginRequirement,
				constants.IsoUsd_IsolatedMarket,
			},
			perpetualPositions: map[types.SubaccountId][]*types.PerpetualPosition{
				secondSubaccountId: {&constants.PerpetualPosition_OneISOLong},
			},
			expectedPerpetualPositions: map[types.SubaccountId][]*types.PerpetualPosition{
				secondSubaccountId: {&constants.PerpetualPosition_OneISOLong},
			},
			expectedAssetPositions: map[types.SubaccountId][]*types.AssetPosition{},
			expectedUpdatedAssetPositions: map[types.SubaccountId][]*types.AssetPosition{
				firstSubaccountId:  {},
				secondSubaccountId: {},
			},
			updates: []types.Update{
				{
					SubaccountId: firstSubaccountId,
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(-100)),
				},
				{
					SubaccountId: secondSubaccountId,
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(100)),
				},
			},
			expectedAssetYieldIndex: constants.AssetYieldIndex_Zero,
			msgSenderEnabled:        true,

			currentBlock: 100,
			negativeTncSubaccountSeenAtBlock: map[uint32]uint32{
				constants.IsoUsd_IsolatedMarket.Params.Id: 100 -
					types.WITHDRAWAL_AND_TRANSFERS_BLOCKED_AFTER_NEGATIVE_TNC_SUBACCOUNT_SEEN_BLOCKS + 1,
				constants.BtcUsd_NoMarginRequirement.Params.Id: 100 -
					types.WITHDRAWAL_AND_TRANSFERS_BLOCKED_AFTER_NEGATIVE_TNC_SUBACCOUNT_SEEN_BLOCKS,
			},

			updateType: types.Transfer,
		},
		`transfers are blocked if negative TNC subaccount was seen within
		WITHDRAWAL_AND_TRANSFERS_BLOCKED_AFTER_NEGATIVE_TNC_SUBACCOUNT_SEEN_BLOCKS for one isolated
		perpetual collateral pool and negative TNC subaccount was never seen for another isolated perpetual
		collateral pool, both of which are associated with subaccounts being updated`: {
			expectedQuoteBalance: big.NewInt(-100),
			expectedSuccess:      false,
			expectedSuccessPerUpdate: []types.UpdateResult{
				types.WithdrawalsAndTransfersBlocked,
				types.WithdrawalsAndTransfersBlocked,
			},
			perpetuals: []perptypes.Perpetual{
				constants.IsoUsd_IsolatedMarket,
				constants.Iso2Usd_IsolatedMarket,
			},
			perpetualPositions: map[types.SubaccountId][]*types.PerpetualPosition{
				firstSubaccountId:  {&constants.PerpetualPosition_OneISOLong},
				secondSubaccountId: {&constants.PerpetualPosition_OneISO2Long},
			},
			expectedPerpetualPositions: map[types.SubaccountId][]*types.PerpetualPosition{
				firstSubaccountId:  {&constants.PerpetualPosition_OneISOLong},
				secondSubaccountId: {&constants.PerpetualPosition_OneISO2Long},
			},
			expectedAssetPositions: map[types.SubaccountId][]*types.AssetPosition{},
			expectedUpdatedAssetPositions: map[types.SubaccountId][]*types.AssetPosition{
				firstSubaccountId:  {},
				secondSubaccountId: {},
			},
			updates: []types.Update{
				{
					SubaccountId: firstSubaccountId,
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(-100)),
				},
				{
					SubaccountId: secondSubaccountId,
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(100)),
				},
			},
			expectedAssetYieldIndex: constants.AssetYieldIndex_Zero,
			msgSenderEnabled:        true,

			currentBlock: 100,
			negativeTncSubaccountSeenAtBlock: map[uint32]uint32{
				constants.IsoUsd_IsolatedMarket.Params.Id: 0,
				constants.Iso2Usd_IsolatedMarket.Params.Id: 100 -
					types.WITHDRAWAL_AND_TRANSFERS_BLOCKED_AFTER_NEGATIVE_TNC_SUBACCOUNT_SEEN_BLOCKS + 1,
			},

			updateType: types.Transfer,
		},
		`transferss are blocked if negative TNC subaccount was seen within
		WITHDRAWAL_AND_TRANSFERS_BLOCKED_AFTER_NEGATIVE_TNC_SUBACCOUNT_SEEN_BLOCKS for one isolated
		perpetual collateral pool and negative TNC subaccount was seen for another the cross-perpetual
		collateral pool after WITHDRAWAL_AND_TRANSFERS_BLOCKED_AFTER_NEGATIVE_TNC_SUBACCOUNT_SEEN_BLOCKS,
		both of which are associated with subaccounts being updated`: {
			expectedQuoteBalance: big.NewInt(-100),
			expectedSuccess:      false,
			expectedSuccessPerUpdate: []types.UpdateResult{
				types.WithdrawalsAndTransfersBlocked,
				types.WithdrawalsAndTransfersBlocked,
			},
			perpetuals: []perptypes.Perpetual{
				constants.IsoUsd_IsolatedMarket,
				constants.Iso2Usd_IsolatedMarket,
			},
			perpetualPositions: map[types.SubaccountId][]*types.PerpetualPosition{
				firstSubaccountId:  {&constants.PerpetualPosition_OneISOLong},
				secondSubaccountId: {&constants.PerpetualPosition_OneISO2Long},
			},
			expectedPerpetualPositions: map[types.SubaccountId][]*types.PerpetualPosition{
				firstSubaccountId:  {&constants.PerpetualPosition_OneISOLong},
				secondSubaccountId: {&constants.PerpetualPosition_OneISO2Long},
			},
			expectedAssetPositions: map[types.SubaccountId][]*types.AssetPosition{},
			expectedUpdatedAssetPositions: map[types.SubaccountId][]*types.AssetPosition{
				firstSubaccountId:  {},
				secondSubaccountId: {},
			},
			updates: []types.Update{
				{
					SubaccountId: firstSubaccountId,
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(-100)),
				},
				{
					SubaccountId: secondSubaccountId,
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(100)),
				},
			},
			expectedAssetYieldIndex: constants.AssetYieldIndex_Zero,
			msgSenderEnabled:        true,

			currentBlock: 100,
			negativeTncSubaccountSeenAtBlock: map[uint32]uint32{
				constants.IsoUsd_IsolatedMarket.Params.Id: 100 -
					types.WITHDRAWAL_AND_TRANSFERS_BLOCKED_AFTER_NEGATIVE_TNC_SUBACCOUNT_SEEN_BLOCKS,
				constants.Iso2Usd_IsolatedMarket.Params.Id: 100 -
					types.WITHDRAWAL_AND_TRANSFERS_BLOCKED_AFTER_NEGATIVE_TNC_SUBACCOUNT_SEEN_BLOCKS + 1,
			},

			updateType: types.Transfer,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctx, keeper, pricesKeeper, perpetualsKeeper, _, _, assetsKeeper, ratelimitKeeper, _, _ := testutil.SubaccountsKeepers(
				t,
				tc.msgSenderEnabled,
			)
			ctx = ctx.WithTxBytes(constants.TestTxBytes)
			testutil.CreateTestMarkets(t, ctx, pricesKeeper)
			testutil.CreateTestLiquidityTiers(t, ctx, perpetualsKeeper)

			// Set up initial sdai price
			rateString := sdaiservertypes.TestSDAIEventRequest.ConversionRate
			rate, conversionErr := ratelimitkeeper.ConvertStringToBigInt(rateString)
			require.NoError(t, conversionErr)

			ratelimitKeeper.SetSDAIPrice(ctx, rate)
			ratelimitKeeper.SetAssetYieldIndex(ctx, big.NewRat(1, 1))

			// ratelimitKeeper.SetCurrentDaiYieldEpochNumber(ctx, 0)
			for _, m := range tc.marketParamPrices {
				_, err := pricesKeeper.CreateMarket(
					ctx,
					m.Param,
					m.Price,
				)
				require.NoError(t, err)
			}

			// Always creates TDai asset first
			require.NoError(t, testutil.CreateTDaiAsset(ctx, assetsKeeper))
			for _, a := range tc.assets {
				_, err := assetsKeeper.CreateAsset(
					ctx,
					a.Id,
					a.Symbol,
					a.Denom,
					a.DenomExponent,
					a.HasMarket,
					a.MarketId,
					a.AtomicResolution,
					a.AssetYieldIndex,
				)
				require.NoError(t, err)
			}

			for i, p := range tc.perpetuals {
				perp, err := perpetualsKeeper.CreatePerpetual(
					ctx,
					p.Params.Id,
					p.Params.Ticker,
					p.Params.MarketId,
					p.Params.AtomicResolution,
					p.Params.DefaultFundingPpm,
					p.Params.LiquidityTier,
					p.Params.MarketType,
					p.Params.DangerIndexPpm,
					p.Params.IsolatedMarketMaxCumulativeInsuranceFundDeltaPerBlock,
					p.YieldIndex,
				)
				require.NoError(t, err)

				// Update FundingIndex for testing settlements.
				if i < len(tc.newFundingIndices) {
					err = perpetualsKeeper.ModifyFundingIndex(
						ctx,
						perp.Params.Id,
						tc.newFundingIndices[i],
					)
					require.NoError(t, err)
				}
			}

			subaccounts := createNSubaccount(keeper, ctx, 2, big.NewInt(1_000))
			for _, subaccount := range subaccounts {
				if perpetualPositions, exists := tc.perpetualPositions[*subaccount.Id]; exists {
					subaccount.PerpetualPositions = perpetualPositions
				} else {
					subaccount.PerpetualPositions = []*types.PerpetualPosition{}
				}
				if assetPositions, exists := tc.assetPositions[*subaccount.Id]; exists {
					subaccount.AssetPositions = assetPositions
				} else {
					subaccount.AssetPositions = []*types.AssetPosition{}
				}
				keeper.SetSubaccount(ctx, subaccount)
			}
			subaccountId := *subaccounts[0].Id

			// Set the negative TNC subaccount seen at block in state if it's greater than 0.
			for perpetualId, negativeTncSubaccountSeenAtBlock := range tc.negativeTncSubaccountSeenAtBlock {
				if negativeTncSubaccountSeenAtBlock != 0 {
					err := keeper.SetNegativeTncSubaccountSeenAtBlock(
						ctx,
						perpetualId,
						negativeTncSubaccountSeenAtBlock,
					)
					require.NoError(t, err)
				}
			}

			// Set the current block number on the context.
			ctx = ctx.WithBlockHeight(int64(tc.currentBlock))

			for i, u := range tc.updates {
				if u.SubaccountId == (types.SubaccountId{}) {
					u.SubaccountId = subaccountId
				}
				tc.updates[i] = u
			}

			success, successPerUpdate, err := keeper.UpdateSubaccounts(ctx, tc.updates, tc.updateType)
			if tc.expectedErr != nil {
				require.ErrorIs(t, tc.expectedErr, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedSuccessPerUpdate, successPerUpdate)
				require.Equal(t, tc.expectedSuccess, success)
			}

			if tc.msgSenderEnabled {
				assertSubaccountUpdateEventsInIndexerBlock(
					t,
					keeper,
					ctx,
					tc.expectedErr,
					tc.expectedSuccess,
					tc.updates,
					tc.expectedSuccessPerUpdate,
					tc.expectedUpdatedPerpetualPositions,
					tc.expectedSubaccoundIdToFundingPayments,
					tc.expectedUpdatedAssetPositions,
					tc.expectedAssetYieldIndex,
				)
			} else {
				assertSubaccountUpdateEventsNotInIndexerBlock(
					t,
					keeper,
					ctx,
				)
			}

			for subaccountIdToCheck, expectedPerpetualPositions := range tc.expectedPerpetualPositions {
				newSubaccount := keeper.GetSubaccount(ctx, subaccountIdToCheck)
				require.Equal(t, len(expectedPerpetualPositions), len(newSubaccount.PerpetualPositions))
				for i, ep := range expectedPerpetualPositions {
					require.Equal(t, *ep, *newSubaccount.PerpetualPositions[i])
				}
			}
			for subaccountIdToCheck, expectedAssetPositions := range tc.expectedAssetPositions {
				newSubaccount := keeper.GetSubaccount(ctx, subaccountIdToCheck)
				require.Equal(t, len(expectedAssetPositions), len(newSubaccount.AssetPositions))
				for i, ap := range expectedAssetPositions {
					require.Equal(t, *ap, *newSubaccount.AssetPositions[i])
				}
			}
		})
	}
}

func TestCanUpdateSubaccounts(t *testing.T) {
	tests := map[string]struct {
		// State.
		perpetuals        []perptypes.Perpetual
		assets            []*asstypes.Asset
		marketParamPrices []pricestypes.MarketParamPrice
		openInterests     []perptypes.OpenInterestDelta

		// Subaccount state.
		useEmptySubaccount        bool
		perpetualPositions        []*types.PerpetualPosition
		assetPositions            []*types.AssetPosition
		additionalTestSubaccounts []types.Subaccount

		// Updates.
		updates    []types.Update
		updateType types.UpdateType

		// Expectations.
		expectedSuccess          bool
		expectedSuccessPerUpdate []types.UpdateResult
		expectedErr              error
	}{
		"(OIMF) OI increased, still at base IMF, match is collateralized": {
			expectedSuccess:          true,
			expectedSuccessPerUpdate: []types.UpdateResult{types.Success, types.Success},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},
			assetPositions: []*types.AssetPosition{
				{
					AssetId: uint32(0),
					// 900_000 TDai (just enough to colalteralize 90 BTC at $50_000 and 20% IMF)
					Quantums: dtypes.NewInt(900_000_000_000),
				},
			},
			additionalTestSubaccounts: []types.Subaccount{
				{
					Id: &constants.Bob_Num0,
					AssetPositions: []*types.AssetPosition{
						{
							AssetId: uint32(0),
							// 900_000 TDai (just enough to colalteralize 90 BTC at $50_000 and 20% IMF)
							Quantums: dtypes.NewInt(900_000_000_000),
						},
					},
					AssetYieldIndex: "1/1",
				},
			},
			updates: []types.Update{
				{
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(0),
							BigQuantumsDelta: big.NewInt(9_000_000_000), // 90 BTC
						},
					},
					AssetUpdates: []types.AssetUpdate{
						{
							AssetId:          uint32(0),
							BigQuantumsDelta: big.NewInt(-4_500_000_000_000), // -4,500,000 TDai
						},
					},
					SubaccountId: constants.Bob_Num0,
				},
				{
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(0),
							BigQuantumsDelta: big.NewInt(-9_000_000_000), // 9 BTC
						},
					},
					AssetUpdates: []types.AssetUpdate{
						{
							AssetId:          uint32(0),
							BigQuantumsDelta: big.NewInt(4_500_000_000_000), // 4,500,000 TDai
						},
					},
				},
			},
			updateType: types.Match,
		},
		"(OIMF) current OI soft lower cap, match collateralized at base IMF but not OIMF": {
			expectedSuccess: false,
			expectedSuccessPerUpdate: []types.UpdateResult{
				types.NewlyUndercollateralized,
				types.NewlyUndercollateralized,
			},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance_25mmLowerCap_50mmUpperCap,
			},
			assetPositions: []*types.AssetPosition{
				{
					AssetId: uint32(0),
					// 900_000 TDai (just enough to colalteralize 90 BTC at $50_000 and 20% IMF)
					Quantums: dtypes.NewInt(900_000_000_000),
				},
			},
			additionalTestSubaccounts: []types.Subaccount{
				{
					Id: &constants.Bob_Num0,
					AssetPositions: []*types.AssetPosition{
						{
							AssetId: uint32(0),
							// 900_000 TDai (just enough to colalteralize 90 BTC at $50_000 and 20% IMF)
							Quantums: dtypes.NewInt(900_000_000_000),
						},
					},
					AssetYieldIndex: "1/1",
				},
			},
			openInterests: []perptypes.OpenInterestDelta{
				{
					PerpetualId: uint32(0),
					// 500 BTC. At $50,000, this is $25,000,000 of OI.
					BaseQuantums: big.NewInt(50_000_000_000),
				},
			},
			updates: []types.Update{
				{
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(0),
							BigQuantumsDelta: big.NewInt(9_000_000_000), // 90 BTC
						},
					},
					AssetUpdates: []types.AssetUpdate{
						{
							AssetId:          uint32(0),
							BigQuantumsDelta: big.NewInt(-4_500_000_000_000), // -4,500,000 TDai
						},
					},
					SubaccountId: constants.Bob_Num0,
				},
				{
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(0),
							BigQuantumsDelta: big.NewInt(-9_000_000_000), // 9 BTC
						},
					},
					AssetUpdates: []types.AssetUpdate{
						{
							AssetId:          uint32(0),
							BigQuantumsDelta: big.NewInt(4_500_000_000_000), // 4,500,000 TDai
						},
					},
				},
			},
			updateType: types.Match,
		},
		"(OIMF) match collateralized at base IMF and just collateralized at OIMF": {
			expectedSuccess: true,
			expectedSuccessPerUpdate: []types.UpdateResult{
				types.Success,
				types.Success,
			},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance_25mmLowerCap_50mmUpperCap,
			},
			assetPositions: []*types.AssetPosition{
				{
					AssetId: uint32(0),
					// 900_000 TDai (just enough to colalteralize 90 BTC at $50_000 and 20% IMF)
					Quantums: dtypes.NewInt(900_000_000_000),
				},
			},
			additionalTestSubaccounts: []types.Subaccount{
				{
					Id: &constants.Bob_Num0,
					AssetPositions: []*types.AssetPosition{
						{
							AssetId: uint32(0),
							// 900_000 TDai (just enough to colalteralize 90 BTC at $50_000 and 20% IMF)
							Quantums: dtypes.NewInt(900_000_000_000),
						},
					},
					AssetYieldIndex: "1/1",
				},
			},
			openInterests: []perptypes.OpenInterestDelta{
				{
					PerpetualId: uint32(0),
					// (Only difference from prevoius test case)
					// 410 BTC. At $50,000, this is $20,500,000 of OI.
					// OI would be $25,000,000 after the Match updates, so OIMF is still at base IMF.
					BaseQuantums: big.NewInt(41_000_000_000),
				},
			},
			updates: []types.Update{
				{
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(0),
							BigQuantumsDelta: big.NewInt(9_000_000_000), // 90 BTC
						},
					},
					AssetUpdates: []types.AssetUpdate{
						{
							AssetId:          uint32(0),
							BigQuantumsDelta: big.NewInt(-4_500_000_000_000), // -4,500,000 TDai
						},
					},
					SubaccountId: constants.Bob_Num0,
				},
				{
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(0),
							BigQuantumsDelta: big.NewInt(-9_000_000_000), // 9 BTC
						},
					},
					AssetUpdates: []types.AssetUpdate{
						{
							AssetId:          uint32(0),
							BigQuantumsDelta: big.NewInt(4_500_000_000_000), // 4,500,000 TDai
						},
					},
				},
			},
			updateType: types.Match,
		},
		"(OIMF) match collateralized at base IMF and just failed collateralization at OIMF": {
			expectedSuccess: false,
			expectedSuccessPerUpdate: []types.UpdateResult{
				types.NewlyUndercollateralized,
				types.NewlyUndercollateralized,
			},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance_25mmLowerCap_50mmUpperCap,
			},
			assetPositions: []*types.AssetPosition{
				{
					AssetId: uint32(0),
					// 900_000 TDai (just enough to colalteralize 90 BTC at $50_000 and 20% IMF)
					Quantums: dtypes.NewInt(900_000_000_000),
				},
			},
			additionalTestSubaccounts: []types.Subaccount{
				{
					Id: &constants.Bob_Num0,
					AssetPositions: []*types.AssetPosition{
						{
							AssetId: uint32(0),
							// 900_000 TDai (just enough to colalteralize 90 BTC at $50_000 and 20% IMF)
							Quantums: dtypes.NewInt(900_000_000_000),
						},
					},
					AssetYieldIndex: "1/1",
				},
			},
			openInterests: []perptypes.OpenInterestDelta{
				{
					PerpetualId: uint32(0),
					// (Only difference from prevoius test case)
					// 410 BTC + 1 base quantum. At $50,000, this is > $20,500,000 of OI.
					// OI would be just past $25,000,000 after the Match updates, so OIMF > IMF = 20%
					BaseQuantums: big.NewInt(41_000_000_001),
				},
			},
			updates: []types.Update{
				{
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(0),
							BigQuantumsDelta: big.NewInt(9_000_000_000), // 90 BTC
						},
					},
					AssetUpdates: []types.AssetUpdate{
						{
							AssetId:          uint32(0),
							BigQuantumsDelta: big.NewInt(-4_500_000_000_000), // -4,500,000 TDai
						},
					},
					SubaccountId: constants.Bob_Num0,
				},
				{
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(0),
							BigQuantumsDelta: big.NewInt(-9_000_000_000), // 9 BTC
						},
					},
					AssetUpdates: []types.AssetUpdate{
						{
							AssetId:          uint32(0),
							BigQuantumsDelta: big.NewInt(4_500_000_000_000), // 4,500,000 TDai
						},
					},
				},
			},
			updateType: types.Match,
		},
		"(OIMF) OIMF caps at 100%, un-leveraged trade always succeeds": {
			expectedSuccess: true,
			expectedSuccessPerUpdate: []types.UpdateResult{
				types.Success,
				types.Success,
			},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance_25mmLowerCap_50mmUpperCap,
			},
			assetPositions: []*types.AssetPosition{
				{
					AssetId: uint32(0),
					// 4_500_000 TDai (just enough to collateralize 90 BTC at $50_000 and 100% IMF)
					Quantums: dtypes.NewInt(4_500_000_000_000)},
			},
			additionalTestSubaccounts: []types.Subaccount{
				{
					Id: &constants.Bob_Num0,
					AssetPositions: []*types.AssetPosition{
						{
							AssetId: uint32(0),
							// 4_500_000 TDai (just enough to collateralize 90 BTC at $50_000 and 100% IMF)
							Quantums: dtypes.NewInt(4_500_000_000_000),
						},
					},
					AssetYieldIndex: "1/1",
				},
			},
			openInterests: []perptypes.OpenInterestDelta{
				{
					PerpetualId: uint32(0),
					// 10_000 BTC. At $50,000, this is $500mm of OI which way past upper cap
					BaseQuantums: big.NewInt(1_000_000_000_000),
				},
			},
			updates: []types.Update{
				{
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(0),
							BigQuantumsDelta: big.NewInt(9_000_000_000), // 90 BTC
						},
					},
					AssetUpdates: []types.AssetUpdate{
						{
							AssetId:          uint32(0),
							BigQuantumsDelta: big.NewInt(-4_500_000_000_000), // -4,500,000 TDai
						},
					},
					SubaccountId: constants.Bob_Num0,
				},
				{
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(0),
							BigQuantumsDelta: big.NewInt(-9_000_000_000), // 9 BTC
						},
					},
					AssetUpdates: []types.AssetUpdate{
						{
							AssetId:          uint32(0),
							BigQuantumsDelta: big.NewInt(4_500_000_000_000), // 4,500,000 TDai
						},
					},
				},
			},
			updateType: types.Match,
		},
		"one update with no existing position and no margin requirements": {
			expectedSuccess:          true,
			expectedSuccessPerUpdate: []types.UpdateResult{types.Success},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_NoMarginRequirement,
			},
			updates: []types.Update{
				{
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(0),
							BigQuantumsDelta: big.NewInt(9_900_000_000), // 99 BTC
						},
					},
				},
			},
		},
		"new TDai asset position exceeds max uint64": {
			assetPositions: testutil.CreateTDaiAssetPosition(new(big.Int).SetUint64(math.MaxUint64)),
			updates: []types.Update{
				{
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(1)),
				},
			},
			updateType:               types.Deposit,
			expectedSuccess:          true,
			expectedSuccessPerUpdate: []types.UpdateResult{types.Success},
		},
		"perpetual does not exist (should never happen)": {
			expectedErr: perptypes.ErrPerpetualDoesNotExist,
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(999999),
					Quantums:     dtypes.NewIntFromUint64(math.MaxUint64),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			updates: []types.Update{
				{},
			},
		},
		"new position quantums exceeds max uint64": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_SmallMarginRequirement,
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewIntFromUint64(math.MaxUint64),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			updates: []types.Update{
				{
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(0),
							BigQuantumsDelta: big.NewInt(1),
						},
					},
				},
			},
			updateType:               types.Deposit,
			expectedSuccess:          true,
			expectedSuccessPerUpdate: []types.UpdateResult{types.Success},
		},
		"update refers to the same position twice": {
			expectedErr: types.ErrNonUniqueUpdatesPosition,
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_SmallMarginRequirement,
			},
			updates: []types.Update{
				{
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(-50_000_000_000)), // -$50,000
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(0),
							BigQuantumsDelta: big.NewInt(100_000_000), // 1 BTC
						},
						{
							PerpetualId:      uint32(0),
							BigQuantumsDelta: big.NewInt(100_000_000), // 1 BTC
						},
					},
				},
			},
		},
		"multiple updates are considered independently for same account": {
			expectedSuccess:          false,
			expectedSuccessPerUpdate: []types.UpdateResult{types.NewlyUndercollateralized, types.Success, types.Success},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_SmallMarginRequirement,
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(1_000_000), // 0.01 BTC
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			updates: []types.Update{
				{
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(-50_000_000_000)), // -$50,000
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(0),
							BigQuantumsDelta: big.NewInt(100_000_000), // 1 BTC
						},
					},
				},
				{
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(-49_999_000_000)), // -$49,999
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(0),
							BigQuantumsDelta: big.NewInt(100_000_000), // 1 BTC
						},
					},
				},
				{
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(0),
							BigQuantumsDelta: big.NewInt(0), // 0 BTC
						},
					},
				},
			},
		},
		"Undercollateralized: " +
			"First update makes account less collateralized, " +
			"Second update results in no change, " +
			"Third update makes account _more_ collateralized," +
			"Fourth update makes it collateralized": {
			assetPositions:  testutil.CreateTDaiAssetPosition(big.NewInt(-496_000_000)), // -$496
			expectedSuccess: false,
			expectedSuccessPerUpdate: []types.UpdateResult{
				types.StillUndercollateralized,
				types.Success,
				types.Success,
				types.Success,
			},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_SmallMarginRequirement,
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(1_000_000), // 0.01 BTC
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			updates: []types.Update{
				{
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(-1)), // -$0.000001
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(0),
							BigQuantumsDelta: big.NewInt(0), // 0 BTC
						},
					},
				},
				{
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(0),
							BigQuantumsDelta: big.NewInt(0), // 0 BTC
						},
					},
				},
				{
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(500_000)), // $.50
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(0),
							BigQuantumsDelta: big.NewInt(0), // 0 BTC
						},
					},
				},
				{
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(2_000_000)), // $2
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(0),
							BigQuantumsDelta: big.NewInt(0), // 0 BTC
						},
					},
				},
			},
		},
		"TDai asset position is negative but increasing when no positions are open": {
			assetPositions: testutil.CreateTDaiAssetPosition(big.NewInt(-10)),
			updates: []types.Update{
				{
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(1)), // $.000001
				},
			},
			expectedSuccess: true,
			expectedSuccessPerUpdate: []types.UpdateResult{
				types.Success,
			},
		},
		"TDai asset position is negative but unchanging when no positions are open": {
			assetPositions: testutil.CreateTDaiAssetPosition(big.NewInt(-10)),
			updates: []types.Update{
				{
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(0)), // $0
				},
			},
			expectedSuccess: false,
			expectedSuccessPerUpdate: []types.UpdateResult{
				types.StillUndercollateralized,
			},
		},
		"TDai asset position decreases below zero when no positions are open": {
			updates: []types.Update{
				{
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(-1)), // -$0.000001
				},
			},
			expectedSuccess: false,
			expectedSuccessPerUpdate: []types.UpdateResult{
				types.NewlyUndercollateralized,
			},
		},
		"TDai asset position decreases further below zero when no positions are open": {
			assetPositions: testutil.CreateTDaiAssetPosition(big.NewInt(-1)),
			updates: []types.Update{
				{
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(-1)), // -$0.000001
				},
			},
			expectedSuccess: false,
			expectedSuccessPerUpdate: []types.UpdateResult{
				types.StillUndercollateralized,
			},
		},
		"two updates on different accounts, second account is new account": {
			assetPositions:           testutil.CreateTDaiAssetPosition(big.NewInt(50_000_000_000)), // $50,000
			expectedSuccess:          false,
			expectedSuccessPerUpdate: []types.UpdateResult{types.Success, types.NewlyUndercollateralized},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000), // 1 BTC
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			updates: []types.Update{
				{
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(-50_000_000_000)), // -$50,000
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(0),
							BigQuantumsDelta: big.NewInt(100_000_000), // 1 BTC
						},
					},
				},
				{
					SubaccountId: types.SubaccountId{
						Owner:  "non-existent-acount",
						Number: uint32(0),
					},
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(-50_000_000_000)), // -$50,000
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(0),
							BigQuantumsDelta: big.NewInt(100_000_000), // 1 BTC
						},
					},
				},
			},
		},
		"unsettled funding reduces TDai asset position to 1; further decrease TDai asset position, still collateralized": {
			assetPositions: testutil.CreateTDaiAssetPosition(big.NewInt(100)),
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(1_000_000), // 0.01 BTC,
					FundingIndex: dtypes.NewInt(-99),       // indexDelta=99, net settlement=-99
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			updates: []types.Update{
				{
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(-1)), // -$0.000001
				},
			},
			expectedSuccess: true,
			expectedSuccessPerUpdate: []types.UpdateResult{
				types.Success,
			},
		},
		"unsettled funding reduces TDai asset position to zero; further decrease TDai asset position, undercollateralized": {
			assetPositions: testutil.CreateTDaiAssetPosition(big.NewInt(100)),
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(1_000_000), // 0.01 BTC,
					FundingIndex: dtypes.NewInt(-100),      // indexDelta=100, net settlement=-100
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			updates: []types.Update{
				{
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(-1)), // -$0.000001
				},
			},
			expectedSuccess: false,
			expectedSuccessPerUpdate: []types.UpdateResult{
				types.NewlyUndercollateralized,
			},
		},
		"unsettled funding makes position undercollateralized": {
			assetPositions: testutil.CreateTDaiAssetPosition(big.NewInt(200)),
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(1_000_000), // 0.01 BTC,
					FundingIndex: dtypes.NewInt(-200),      // indexDelta=200, net settlement=-200
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			updates: []types.Update{
				{
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(-1)), // -$0.000001
				},
			},
			expectedSuccess: false,
			expectedSuccessPerUpdate: []types.UpdateResult{
				types.NewlyUndercollateralized,
			},
		},
		"position was undercollateralized before update due to funding and still undercollateralized" +
			"after due to funding": {
			assetPositions: testutil.CreateTDaiAssetPosition(big.NewInt(199)),
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(1_000_000), // 0.01 BTC,
					FundingIndex: dtypes.NewInt(-200),      // indexDelta=200, net settlement=-200
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			updates: []types.Update{
				{
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(-1)), // -$0.000001
				},
			},
			expectedSuccess: false,
			expectedSuccessPerUpdate: []types.UpdateResult{
				types.StillUndercollateralized,
			},
		},
		"unsettled funding makes position with negative TDai asset position collateralized before update": {
			assetPositions: testutil.CreateTDaiAssetPosition(big.NewInt(-100)),
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(1_000_000), // 0.01 BTC,
					FundingIndex: dtypes.NewInt(100),       // indexDelta=-100, net settlement=100
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			updates: []types.Update{
				{},
			},
			expectedSuccess: true,
			expectedSuccessPerUpdate: []types.UpdateResult{
				types.Success,
			},
		},
		"adding unsettled funding to TDai asset position exceeds max uint64": {
			assetPositions: testutil.CreateTDaiAssetPosition(new(big.Int).SetUint64(math.MaxUint64 - 1)),
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(1_000_000), // 0.01 BTC,
					FundingIndex: dtypes.NewInt(100),       // indexDelta=-100, net settlement=100
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			updates: []types.Update{
				{},
			},
			updateType:               types.Deposit,
			expectedSuccess:          true,
			expectedSuccessPerUpdate: []types.UpdateResult{types.Success},
		},
		"adding unsettled funding to TDai asset position exceeds negative max uint64": {
			assetPositions: testutil.CreateTDaiAssetPosition(
				new(big.Int).Neg(new(big.Int).SetUint64(math.MaxUint64 - 1)),
			),
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(1_000_000), // 0.01 BTC,
					FundingIndex: dtypes.NewInt(-100),      // indexDelta=100, net settlement=-100
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			updates: []types.Update{
				{},
			},
			updateType:               types.Deposit,
			expectedSuccess:          true,
			expectedSuccessPerUpdate: []types.UpdateResult{types.Success},
		},
		"adding unsettled funding, original TDai asset position and TDai asset position delta exceeds max int64": {
			assetPositions: testutil.CreateTDaiAssetPosition(
				new(big.Int).SetUint64(math.MaxUint64 - 5),
			),
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(1_000_000), // 0.01 BTC,
					FundingIndex: dtypes.NewInt(3),         // indexDelta=-3, net settlement=3
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			updates: []types.Update{
				{
					AssetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(3)), // $3
				},
			},
			updateType:               types.Deposit,
			expectedSuccess:          true,
			expectedSuccessPerUpdate: []types.UpdateResult{types.Success},
		},
		"2 updates, 1 update involves not-updatable perp": {
			assetPositions: testutil.CreateTDaiAssetPosition(big.NewInt(1_000_000_000_000)),
			expectedErr:    types.ErrProductPositionNotUpdatable,
			perpetuals: []perptypes.Perpetual{
				*perptest.GeneratePerpetual(
					perptest.WithId(100),
					perptest.WithMarketId(100),
				),
				*perptest.GeneratePerpetual(
					perptest.WithId(101),
					perptest.WithMarketId(101),
				),
			},
			marketParamPrices: []pricestypes.MarketParamPrice{
				*pricestest.GenerateMarketParamPrice(pricestest.WithId(100)),
				*pricestest.GenerateMarketParamPrice(
					pricestest.WithId(101),
					pricestest.WithSpotPriceValue(0),
					pricestest.WithPnlPriceValue(0),
				),
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(100),
					Quantums:     dtypes.NewInt(1_000_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
				{
					PerpetualId:  uint32(101),
					Quantums:     dtypes.NewInt(1_000_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			updates: []types.Update{
				{
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(100),
							BigQuantumsDelta: big.NewInt(-1_000),
						},
						{
							PerpetualId:      uint32(101),
							BigQuantumsDelta: big.NewInt(1_000),
						},
					},
				},
			},
		},
		"Isolated subaccounts - has update for both an isolated perpetual and non-isolated perpetual": {
			assetPositions:           testutil.CreateTDaiAssetPosition(big.NewInt(1_000_000_000_000)),
			expectedSuccess:          false,
			expectedSuccessPerUpdate: []types.UpdateResult{types.ViolatesIsolatedSubaccountConstraints},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_NoMarginRequirement,
				constants.IsoUsd_IsolatedMarket,
			},
			updates: []types.Update{
				{
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(0),
							BigQuantumsDelta: big.NewInt(-100_000_000), // -1 BTC
						},
						{
							PerpetualId:      uint32(3),
							BigQuantumsDelta: big.NewInt(1_000_000_000), // 1 ISO
						},
					},
				},
			},
		},
		"Isolated subaccounts - has update for both 2 isolated perpetuals": {
			assetPositions:           testutil.CreateTDaiAssetPosition(big.NewInt(1_000_000_000_000)),
			expectedSuccess:          false,
			expectedSuccessPerUpdate: []types.UpdateResult{types.ViolatesIsolatedSubaccountConstraints},
			perpetuals: []perptypes.Perpetual{
				constants.IsoUsd_IsolatedMarket,
				constants.Iso2Usd_IsolatedMarket,
			},
			updates: []types.Update{
				{
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(3),
							BigQuantumsDelta: big.NewInt(-1_000_000_000), // 1 ISO
						},
						{
							PerpetualId:      uint32(4),
							BigQuantumsDelta: big.NewInt(10_000_000), // 1 ISO2
						},
					},
				},
			},
		},
		"Isolated subaccounts - subaccount with isolated perpetual position has update for non-isolated perpetual": {
			assetPositions:           testutil.CreateTDaiAssetPosition(big.NewInt(1_000_000_000_000)),
			expectedSuccess:          false,
			expectedSuccessPerUpdate: []types.UpdateResult{types.ViolatesIsolatedSubaccountConstraints},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_NoMarginRequirement,
				constants.IsoUsd_IsolatedMarket,
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(3),
					Quantums:     dtypes.NewInt(1_000_000_000), // 1 ISO
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			updates: []types.Update{
				{
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(0),
							BigQuantumsDelta: big.NewInt(-100_000_000), // -1 BTC
						},
					},
				},
			},
		},
		"Isolated subaccounts - subaccount with isolated perpetual position has update for another isolated perpetual": {
			assetPositions:           testutil.CreateTDaiAssetPosition(big.NewInt(1_000_000_000_000)),
			expectedSuccess:          false,
			expectedSuccessPerUpdate: []types.UpdateResult{types.ViolatesIsolatedSubaccountConstraints},
			perpetuals: []perptypes.Perpetual{
				constants.IsoUsd_IsolatedMarket,
				constants.Iso2Usd_IsolatedMarket,
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(3),
					Quantums:     dtypes.NewInt(1_000_000_000), // 1 ISO
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			updates: []types.Update{
				{
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(4),
							BigQuantumsDelta: big.NewInt(-10_000_000), // -1 ISO2
						},
					},
				},
			},
		},
		"Isolated subaccounts - subaccount with non-isolated perpetual position has update for isolated perpetual": {
			assetPositions:           testutil.CreateTDaiAssetPosition(big.NewInt(1_000_000_000_000)),
			expectedSuccess:          false,
			expectedSuccessPerUpdate: []types.UpdateResult{types.ViolatesIsolatedSubaccountConstraints},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_NoMarginRequirement,
				constants.IsoUsd_IsolatedMarket,
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000), // 1 BTC
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			updates: []types.Update{
				{
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(3),
							BigQuantumsDelta: big.NewInt(-1_000_000_000), // -1 ISO
						},
					},
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctx, keeper, pricesKeeper, perpetualsKeeper, _, _, assetsKeeper, ratelimitKeeper, _, _ := testutil.SubaccountsKeepers(t, true)
			testutil.CreateTestMarkets(t, ctx, pricesKeeper)
			testutil.CreateTestLiquidityTiers(t, ctx, perpetualsKeeper)

			// Set up initial sdai price
			rateString := sdaiservertypes.TestSDAIEventRequest.ConversionRate
			rate, conversionErr := ratelimitkeeper.ConvertStringToBigInt(rateString)
			require.NoError(t, conversionErr)

			ratelimitKeeper.SetSDAIPrice(ctx, rate)
			ratelimitKeeper.SetAssetYieldIndex(ctx, big.NewRat(1, 1))

			// ratelimitKeeper.SetCurrentDaiYieldEpochNumber(ctx, 0)

			require.NoError(t, testutil.CreateTDaiAsset(ctx, assetsKeeper))
			for _, a := range tc.assets {
				_, err := assetsKeeper.CreateAsset(
					ctx,
					a.Id,
					a.Symbol,
					a.Denom,
					a.DenomExponent,
					a.HasMarket,
					a.MarketId,
					a.AtomicResolution,
					a.AssetYieldIndex,
				)
				require.NoError(t, err)
			}

			for _, m := range tc.marketParamPrices {
				_, err := pricesKeeper.CreateMarket(
					ctx,
					m.Param,
					m.Price,
				)
				require.NoError(t, err)
			}

			for _, p := range tc.perpetuals {
				_, err := perpetualsKeeper.CreatePerpetual(
					ctx,
					p.Params.Id,
					p.Params.Ticker,
					p.Params.MarketId,
					p.Params.AtomicResolution,
					p.Params.DefaultFundingPpm,
					p.Params.LiquidityTier,
					p.Params.MarketType,
					p.Params.DangerIndexPpm,
					p.Params.IsolatedMarketMaxCumulativeInsuranceFundDeltaPerBlock,
					p.YieldIndex,
				)
				require.NoError(t, err)
			}

			for _, openInterest := range tc.openInterests {
				// Update open interest for each perpetual from default `0`.
				require.NoError(t, perpetualsKeeper.ModifyOpenInterest(
					ctx,
					openInterest.PerpetualId,
					openInterest.BaseQuantums,
				))
			}

			subaccountId := types.SubaccountId{Owner: "foo", Number: 0}
			if !tc.useEmptySubaccount {
				subaccount := createNSubaccount(keeper, ctx, 1, big.NewInt(1_000))[0]
				subaccount.PerpetualPositions = tc.perpetualPositions
				subaccount.AssetPositions = tc.assetPositions
				keeper.SetSubaccount(ctx, subaccount)
				subaccountId = *subaccount.Id
			}

			for i, u := range tc.updates {
				if u.SubaccountId == (types.SubaccountId{}) {
					u.SubaccountId = subaccountId
				}
				tc.updates[i] = u
			}

			for _, sa := range tc.additionalTestSubaccounts {
				keeper.SetSubaccount(ctx, sa)
			}

			// If test case has unspecified update type, use `CollatCheck` as default.
			updateType := tc.updateType
			if updateType == types.UpdateTypeUnspecified {
				updateType = types.CollatCheck
			}
			success, successPerUpdate, err := keeper.CanUpdateSubaccounts(ctx, tc.updates, updateType)
			if tc.expectedErr != nil {
				require.ErrorIs(t, tc.expectedErr, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedSuccessPerUpdate, successPerUpdate)
				require.Equal(t, tc.expectedSuccess, success)
			}

			for _, openInterest := range tc.openInterests {
				// Check open interest for each perpetual did not change after the check.
				perp, err := perpetualsKeeper.GetPerpetual(ctx, openInterest.PerpetualId)
				require.NoError(t, err)
				require.Zerof(t,
					openInterest.BaseQuantums.Cmp(perp.OpenInterest.BigInt()),
					"expected: %s, got: %s",
					openInterest.BaseQuantums.String(),
					perp.OpenInterest.String(),
				)
			}
		})
	}
}

func TestGetNetCollateralAndMarginRequirements(t *testing.T) {
	tests := map[string]struct {
		// state
		perpetuals []perptypes.Perpetual
		assets     []*asstypes.Asset

		// subaccount state
		useEmptySubaccount bool
		perpetualPositions []*types.PerpetualPosition
		assetPositions     []*types.AssetPosition

		// updates
		assetUpdates     []types.AssetUpdate
		perpetualUpdates []types.PerpetualUpdate

		// expectations
		expectedNetCollateral     *big.Int
		expectedInitialMargin     *big.Int
		expectedMaintenanceMargin *big.Int
		expectedErr               error
	}{
		"zero balance": {},
		"non-negative TDai asset position": {
			assetPositions:        testutil.CreateTDaiAssetPosition(big.NewInt(123_456)),
			expectedNetCollateral: big.NewInt(123_456),
		},
		"negative TDai asset position": {
			assetPositions:        testutil.CreateTDaiAssetPosition(big.NewInt(-123_456)),
			expectedNetCollateral: big.NewInt(-123_456),
		},
		"TDai asset position with update": {
			assetPositions:        testutil.CreateTDaiAssetPosition(big.NewInt(-123_456)),
			expectedNetCollateral: big.NewInt(0),
			assetUpdates:          testutil.CreateTDaiAssetUpdate(big.NewInt(123_456)),
		},
		"single perpetual and TDai asset position": {
			assetPositions:        testutil.CreateTDaiAssetPosition(big.NewInt(10_000_000_001)), // $10,000.000001
			expectedNetCollateral: big.NewInt(60_000_000_001),                                   // $60,000.000001
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_NoMarginRequirement,
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000), // 1 BTC
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
		},
		"single perpetual, TDai asset position and unsettled funding (long)": {
			assetPositions:        testutil.CreateTDaiAssetPosition(big.NewInt(10_000_000_001)), // $10,000.000001
			expectedNetCollateral: big.NewInt(60_006_250_001),                                   // $60,006.250001
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_NoMarginRequirement,
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000), // 1 BTC
					FundingIndex: dtypes.NewInt(62500),       // 0.0125% rate at BTC=50,000 TDai
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
		},
		"single perpetual, TDai asset position and unsettled funding (short)": {
			assetPositions:        testutil.CreateTDaiAssetPosition(big.NewInt(-10_000_000_001)), // -$10,000.000001
			expectedNetCollateral: big.NewInt(-60_006_250_001),                                   // -$60,006.250001
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_NoMarginRequirement,
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(-100_000_000), // 1 BTC
					FundingIndex: dtypes.NewInt(62500),        // 0.0125% rate at BTC=50,000 TDai
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
		},
		"non-existing perpetual heled by subaccount (should never happen)": {
			assetPositions: testutil.CreateTDaiAssetPosition(
				big.NewInt(-10_000_000_001), // -$10,000.000001
			),
			expectedNetCollateral: big.NewInt(-60_006_250_001), // -$60,006.250001
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_NoMarginRequirement,
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(999999999),
					Quantums:     dtypes.NewInt(-100_000_000), // 1 BTC
					FundingIndex: dtypes.NewInt(62500),        // 0.0125% rate at BTC=50,000 TDai
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			expectedErr: perptypes.ErrPerpetualDoesNotExist,
		},
		"TDai asset position update underflows uint64": {
			assetPositions: testutil.CreateTDaiAssetPosition(
				constants.BigNegMaxUint64(),
			),
			assetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(-1)),
		},
		"TDai asset position update overflows uint64": {
			assetPositions: testutil.CreateTDaiAssetPosition(
				new(big.Int).SetUint64(math.MaxUint64),
			),
			assetUpdates: testutil.CreateTDaiAssetUpdate(big.NewInt(1)),
		},
		"update for non-existent perpetual": {
			expectedErr: perptypes.ErrPerpetualDoesNotExist,
			perpetualUpdates: []types.PerpetualUpdate{
				{
					PerpetualId:      uint32(0),
					BigQuantumsDelta: big.NewInt(100_000_000), // 1 BTC
				},
			},
		},
		"update with no existing position": {
			assetPositions:        testutil.CreateTDaiAssetPosition(big.NewInt(10_000_000_001)), // $10,000.000001
			expectedNetCollateral: big.NewInt(60_000_000_001),                                   // $60,000.000001
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_NoMarginRequirement,
			},
			perpetualUpdates: []types.PerpetualUpdate{
				{
					PerpetualId:      uint32(0),
					BigQuantumsDelta: big.NewInt(100_000_000), // 1 BTC
				},
			},
		},
		"single perpetual with TDai asset position and positive update to perpetual": {
			assetPositions:        testutil.CreateTDaiAssetPosition(big.NewInt(10_000_000_001)), // $10,000.000001
			expectedNetCollateral: big.NewInt(110_000_000_001),                                  // $110,000.000001
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_NoMarginRequirement,
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000), // 1 BTC
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			perpetualUpdates: []types.PerpetualUpdate{
				{
					PerpetualId:      uint32(0),
					BigQuantumsDelta: big.NewInt(100_000_000), // 1 BTC
				},
			},
		},
		"single long perpetual with position and update which overflows uint64": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_NoMarginRequirement,
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewIntFromUint64(math.MaxUint64),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			perpetualUpdates: []types.PerpetualUpdate{
				{
					PerpetualId:      uint32(0),
					BigQuantumsDelta: big.NewInt(1),
				},
			},
		},
		"single short perpetual with position and update which overflows uint64": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_NoMarginRequirement,
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId: uint32(0),
					Quantums: dtypes.NewIntFromBigInt(
						new(big.Int).Neg(
							new(big.Int).SetUint64(math.MaxUint64),
						)),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			perpetualUpdates: []types.PerpetualUpdate{
				{
					PerpetualId:      uint32(0),
					BigQuantumsDelta: big.NewInt(-1),
				},
			},
		},
		"single perpetual with TDai asset position and negative update to perpetual": {
			assetPositions:        testutil.CreateTDaiAssetPosition(big.NewInt(10_000_000_001)), // $10,000.000001
			expectedNetCollateral: big.NewInt(10_000_000_001),                                   // $10,000.000001
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_NoMarginRequirement,
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000), // 1 BTC
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			perpetualUpdates: []types.PerpetualUpdate{
				{
					PerpetualId:      uint32(0),
					BigQuantumsDelta: big.NewInt(-100_000_000), // -1 BTC
				},
			},
		},
		"multiple asset updates for the same position": {
			expectedErr: types.ErrNonUniqueUpdatesPosition,
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_NoMarginRequirement,
			},
			assetPositions: []*types.AssetPosition{
				&constants.TDai_Asset_100_000,
				&constants.Short_Asset_1BTC,
			},
			assetUpdates: []types.AssetUpdate{
				{
					AssetId:          constants.BtcUsd.Id,
					BigQuantumsDelta: big.NewInt(-100_000_000), // -1 BTC
				},
				{
					AssetId:          constants.BtcUsd.Id,
					BigQuantumsDelta: big.NewInt(-100_000_000), // -1 BTC
				},
			},
		},
		"multiple perpetual updates for the same position": {
			expectedErr: types.ErrNonUniqueUpdatesPosition,
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_NoMarginRequirement,
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000), // 1 BTC
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			perpetualUpdates: []types.PerpetualUpdate{
				{
					PerpetualId:      uint32(0),
					BigQuantumsDelta: big.NewInt(-100_000_000), // -1 BTC
				},
				{
					PerpetualId:      uint32(0),
					BigQuantumsDelta: big.NewInt(-100_000_000), // -1 BTC
				},
			},
		},
		"speculative update to non-existent subaccount": {
			useEmptySubaccount:        true,
			assetUpdates:              testutil.CreateTDaiAssetUpdate(big.NewInt(1_000_000)),
			expectedNetCollateral:     big.NewInt(-99_249_000_000), // $1 - $100,000 (BTC update) + $750 (ETH update)
			expectedInitialMargin:     big.NewInt(50_150_000_000),  // $50,000 (BTC update) + $150 (ETH update)
			expectedMaintenanceMargin: big.NewInt(40_075_000_000),  // $40,000 (BTC update) + $75 (ETH update)
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_50PercentInitial_40PercentMaintenance,
				constants.EthUsd_20PercentInitial_10PercentMaintenance,
			},
			perpetualUpdates: []types.PerpetualUpdate{
				{
					PerpetualId:      uint32(0),
					BigQuantumsDelta: big.NewInt(-200_000_000), // -2 BTC
				},
				{
					PerpetualId:      uint32(1),
					BigQuantumsDelta: big.NewInt(250_000_000), // .25 ETH
				},
			},
		},
		"multiple perpetuals with margin requirements and updates": {
			// $1
			assetPositions: testutil.CreateTDaiAssetPosition(big.NewInt(1000000)),
			// $1 + $50,000 (BTC) + $1,500 (ETH) - $100,000 (BTC update) + $750 (ETH update)
			expectedNetCollateral: big.NewInt(-47_749_000_000),
			// abs($25,000 (BTC) - $50,000 (BTC update)) + $300 (ETH) + $150 (ETH update)
			expectedInitialMargin: big.NewInt(25_450_000_000),
			// abs($20,000 (BTC) - $40,000 (BTC update)) + $150 (ETH) + $75 (ETH update)
			expectedMaintenanceMargin: big.NewInt(20_225_000_000),
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_50PercentInitial_40PercentMaintenance,
				constants.EthUsd_20PercentInitial_10PercentMaintenance,
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000), // 1 BTC
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
				{
					PerpetualId:  uint32(1),
					Quantums:     dtypes.NewInt(500_000_000), // .5 ETH
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			perpetualUpdates: []types.PerpetualUpdate{
				{
					PerpetualId:      uint32(0),
					BigQuantumsDelta: big.NewInt(-200_000_000), // -2 BTC
				},
				{
					PerpetualId:      uint32(1),
					BigQuantumsDelta: big.NewInt(250_000_000), // .25 ETH
				},
			},
		},
		"single perpetual": {
			expectedNetCollateral: big.NewInt(50_000_000_000),
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_NoMarginRequirement,
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId: uint32(0),
					Quantums:    dtypes.NewInt(100_000_000), // 1 BTC
					YieldIndex:  big.NewRat(0, 1).String(),
				},
			},
		},
		"asset with no balance and update": {
			expectedErr: asstypes.ErrNotImplementedMulticollateral,
			assets: []*asstypes.Asset{
				constants.BtcUsd,
			},
			assetPositions: []*types.AssetPosition{
				&constants.Long_Asset_1BTC,
			},
			assetUpdates: []types.AssetUpdate{
				{
					AssetId:          constants.BtcUsd.Id,
					BigQuantumsDelta: big.NewInt(100_000_000), // 1 BTC
				},
			},
		},
		"single positive asset": {
			expectedErr: asstypes.ErrNotImplementedMulticollateral,
			assets: []*asstypes.Asset{
				constants.BtcUsd,
			},
			assetPositions: []*types.AssetPosition{
				&constants.Long_Asset_1BTC,
			},
		},
		"single negative asset": {
			expectedErr: asstypes.ErrNotImplementedMulticollateral,
			assets: []*asstypes.Asset{
				constants.BtcUsd,
			},
			assetPositions: []*types.AssetPosition{
				&constants.Short_Asset_1BTC,
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctx, keeper, pricesKeeper, perpetualsKeeper, _, _, assetsKeeper, ratelimitKeeper, _, _ := testutil.SubaccountsKeepers(t, true)
			testutil.CreateTestMarkets(t, ctx, pricesKeeper)
			testutil.CreateTestLiquidityTiers(t, ctx, perpetualsKeeper)

			// Set up initial sdai price
			rateString := sdaiservertypes.TestSDAIEventRequest.ConversionRate
			rate, conversionErr := ratelimitkeeper.ConvertStringToBigInt(rateString)
			require.NoError(t, conversionErr)

			ratelimitKeeper.SetSDAIPrice(ctx, rate)
			ratelimitKeeper.SetAssetYieldIndex(ctx, big.NewRat(1, 1))

			// ratelimitKeeper.SetCurrentDaiYieldEpochNumber(ctx, 0)

			require.NoError(t, testutil.CreateTDaiAsset(ctx, assetsKeeper))
			for _, a := range tc.assets {
				_, err := assetsKeeper.CreateAsset(
					ctx,
					a.Id,
					a.Symbol,
					a.Denom,
					a.DenomExponent,
					a.HasMarket,
					a.MarketId,
					a.AtomicResolution,
					a.AssetYieldIndex,
				)
				require.NoError(t, err)
			}

			for _, p := range tc.perpetuals {
				_, err := perpetualsKeeper.CreatePerpetual(
					ctx,
					p.Params.Id,
					p.Params.Ticker,
					p.Params.MarketId,
					p.Params.AtomicResolution,
					p.Params.DefaultFundingPpm,
					p.Params.LiquidityTier,
					p.Params.MarketType,
					p.Params.DangerIndexPpm,
					p.Params.IsolatedMarketMaxCumulativeInsuranceFundDeltaPerBlock,
					p.YieldIndex,
				)
				require.NoError(t, err)
			}

			subaccountId := types.SubaccountId{Owner: "foo", Number: 0}
			if !tc.useEmptySubaccount {
				subaccount := createNSubaccount(keeper, ctx, 1, big.NewInt(1_000))[0]
				subaccount.PerpetualPositions = tc.perpetualPositions
				subaccount.AssetPositions = tc.assetPositions
				keeper.SetSubaccount(ctx, subaccount)
				subaccountId = *subaccount.Id
			}

			update := types.Update{
				SubaccountId:     subaccountId,
				AssetUpdates:     tc.assetUpdates,
				PerpetualUpdates: tc.perpetualUpdates,
			}

			netCollateral, initialMargin, maintenanceMargin, err :=
				keeper.GetNetCollateralAndMarginRequirements(ctx, update)

			if tc.expectedErr != nil {
				require.ErrorIs(t, tc.expectedErr, err)
			} else {
				// Testify is bad at printing unsigned integers and prints values as hex
				// https://github.com/stretchr/testify/issues/1116
				// for that reason we convert to strings here to make the output more readable
				if tc.expectedNetCollateral != nil {
					require.Equal(t, tc.expectedNetCollateral.String(), netCollateral.String())
				}
				if tc.expectedInitialMargin != nil {
					require.Equal(t, tc.expectedInitialMargin.String(), initialMargin.String())
				}
				if tc.expectedMaintenanceMargin != nil {
					require.Equal(t, tc.expectedMaintenanceMargin.String(), maintenanceMargin.String())
				}
				require.NoError(t, err)
			}
		})
	}
}

func TestIsValidStateTransitionForUndercollateralizedSubaccount_ZeroMarginRequirements(t *testing.T) {
	tests := map[string]struct {
		bigCurNetCollateral     *big.Int
		bigCurInitialMargin     *big.Int
		bigCurMaintenanceMargin *big.Int
		bigNewNetCollateral     *big.Int
		bigNewMaintenanceMargin *big.Int

		expectedResult types.UpdateResult
	}{
		// Tests when current margin requirement is zero and margin requirement increases.
		"fails when MMR increases and TNC decreases - negative TNC": {
			bigCurNetCollateral:     big.NewInt(-1),
			bigCurInitialMargin:     big.NewInt(0),
			bigCurMaintenanceMargin: big.NewInt(0),
			bigNewNetCollateral:     big.NewInt(-2),
			bigNewMaintenanceMargin: big.NewInt(1),
			expectedResult:          types.StillUndercollateralized,
		},
		"fails when MMR increases and TNC stays the same - negative TNC": {
			bigCurNetCollateral:     big.NewInt(-1),
			bigCurInitialMargin:     big.NewInt(0),
			bigCurMaintenanceMargin: big.NewInt(0),
			bigNewNetCollateral:     big.NewInt(-1),
			bigNewMaintenanceMargin: big.NewInt(1),
			expectedResult:          types.StillUndercollateralized,
		},
		"fails when MMR increases and TNC increases - negative TNC": {
			bigCurNetCollateral:     big.NewInt(-1),
			bigCurInitialMargin:     big.NewInt(0),
			bigCurMaintenanceMargin: big.NewInt(0),
			bigNewNetCollateral:     big.NewInt(100),
			bigNewMaintenanceMargin: big.NewInt(1),
			expectedResult:          types.StillUndercollateralized,
		},
		// Tests when both margin requirements are zero.
		"fails when both new and old MMR are zero and TNC stays the same": {
			bigCurNetCollateral:     big.NewInt(-1),
			bigCurInitialMargin:     big.NewInt(0),
			bigCurMaintenanceMargin: big.NewInt(0),
			bigNewNetCollateral:     big.NewInt(-1),
			bigNewMaintenanceMargin: big.NewInt(0),
			expectedResult:          types.StillUndercollateralized,
		},
		"fails when both new and old MMR are zero and TNC decrease from negative to negative": {
			bigCurNetCollateral:     big.NewInt(-1),
			bigCurInitialMargin:     big.NewInt(0),
			bigCurMaintenanceMargin: big.NewInt(0),
			bigNewNetCollateral:     big.NewInt(-2),
			bigNewMaintenanceMargin: big.NewInt(0),
			expectedResult:          types.StillUndercollateralized,
		},
		"succeeds when both new and old MMR are zero and TNC increases": {
			bigCurNetCollateral:     big.NewInt(-2),
			bigCurInitialMargin:     big.NewInt(0),
			bigCurMaintenanceMargin: big.NewInt(0),
			bigNewNetCollateral:     big.NewInt(-1),
			bigNewMaintenanceMargin: big.NewInt(0),
			expectedResult:          types.Success,
		},
		// Tests when new margin requirement is zero.
		"fails when MMR decreased to zero, and TNC increases but is still negative": {
			bigCurNetCollateral:     big.NewInt(-2),
			bigCurInitialMargin:     big.NewInt(1),
			bigCurMaintenanceMargin: big.NewInt(1),
			bigNewNetCollateral:     big.NewInt(-1),
			bigNewMaintenanceMargin: big.NewInt(0),
			expectedResult:          types.StillUndercollateralized,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			require.Equal(
				t,
				tc.expectedResult,
				keeper.IsValidStateTransitionForUndercollateralizedSubaccount(
					tc.bigCurNetCollateral,
					tc.bigCurInitialMargin,
					tc.bigCurMaintenanceMargin,
					tc.bigNewNetCollateral,
					tc.bigNewMaintenanceMargin,
				),
			)
		})
	}
}
