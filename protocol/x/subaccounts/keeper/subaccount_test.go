package keeper_test

import (
	"math"
	"math/big"
	"math/rand"
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	indexerevents "github.com/dydxprotocol/v4-chain/protocol/indexer/events"
	big_testutil "github.com/dydxprotocol/v4-chain/protocol/testutil/big"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	testutil "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/nullify"
	perptest "github.com/dydxprotocol/v4-chain/protocol/testutil/perpetuals"
	pricestest "github.com/dydxprotocol/v4-chain/protocol/testutil/prices"
	asstypes "github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/require"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func createNSubaccount(keeper *keeper.Keeper, ctx sdk.Context, n int, usdcBalance *big.Int) []types.Subaccount {
	items := make([]types.Subaccount, n)
	for i := range items {
		items[i].Id = &types.SubaccountId{
			Owner:  strconv.Itoa(i),
			Number: uint32(i),
		}
		items[i].AssetPositions = testutil.CreateUsdcAssetPosition(usdcBalance)

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
		expecetedSubaccountUpdateEvent := indexerevents.NewSubaccountUpdateEvent(
			&update.SubaccountId,
			expectedUpdatedPerpetualPositions[update.SubaccountId],
			expectedUpdatedAssetPositions[update.SubaccountId],
			expectedSubaccoundIdToFundingPayments[update.SubaccountId],
		)
		require.Contains(t, subaccountUpdates, expecetedSubaccountUpdateEvent)
	}
}

func TestSubaccountGet(t *testing.T) {
	ctx, keeper, _, _, _, _, _, _ := testutil.SubaccountsKeepers(t, true)
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
	ctx, keeper, _, _, _, _, _, _ := testutil.SubaccountsKeepers(t, true)
	keeper.SetSubaccount(ctx, types.Subaccount{
		Id: &constants.Alice_Num0,
	})

	require.Len(t, keeper.GetAllSubaccount(ctx), 0)

	keeper.SetSubaccount(ctx, types.Subaccount{
		Id:             &constants.Alice_Num0,
		AssetPositions: testutil.CreateUsdcAssetPosition(big.NewInt(1_000)),
	})
	keeper.SetSubaccount(ctx, types.Subaccount{
		Id: &constants.Alice_Num0,
	})
	require.Len(t, keeper.GetAllSubaccount(ctx), 0)
}

func TestSubaccountGetNonExistent(t *testing.T) {
	ctx, keeper, _, _, _, _, _, _ := testutil.SubaccountsKeepers(t, true)
	id := types.SubaccountId{
		Owner:  "non-existent",
		Number: uint32(123),
	}
	acct := keeper.GetSubaccount(ctx, id)
	require.Equal(t, &id, acct.Id)
	require.Equal(t, new(big.Int), acct.GetUsdcPosition())
	require.Empty(t, acct.AssetPositions)
	require.Empty(t, acct.PerpetualPositions)
	require.False(t, acct.MarginEnabled)
}

func TestGetAllSubaccount(t *testing.T) {
	ctx, keeper, _, _, _, _, _, _ := testutil.SubaccountsKeepers(t, true)
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
			ctx, keeper, _, _, _, _, _, _ := testutil.SubaccountsKeepers(t, true)
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

func TestForEachSubaccountRandomStart(t *testing.T) {
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
			rand := rand.New(rand.NewSource(53))
			ctx, keeper, _, _, _, _, _, _ := testutil.SubaccountsKeepers(t, true)
			_ = createNSubaccount(keeper, ctx, tc.numSubaccountsInState, big.NewInt(1_000))
			collectedSubaccounts := make([]types.Subaccount, 0)
			i := 0
			keeper.ForEachSubaccountRandomStart(
				ctx,
				func(subaccount types.Subaccount) bool {
					i++
					collectedSubaccounts = append(collectedSubaccounts, subaccount)
					return i == tc.iterationCount
				},
				rand,
			)

			require.Len(t, collectedSubaccounts, tc.iterationCount)

			if tc.iterationCount > 0 {
				subaccounts := keeper.GetAllSubaccount(ctx)

				offset := 0
				for i, subaccount := range subaccounts {
					if *subaccount.Id == *collectedSubaccounts[0].Id {
						offset = i
						break
					}
				}

				for i := 0; i < tc.iterationCount; i++ {
					require.Equal(
						t,
						subaccounts[(i+offset)%(tc.numSubaccountsInState)],
						collectedSubaccounts[i],
					)
				}
			}
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

		// subaccount state
		perpetualPositions []*types.PerpetualPosition
		assetPositions     []*types.AssetPosition

		// updates
		updates []types.Update

		// expectations
		expectedQuoteBalance       *big.Int
		expectedPerpetualPositions []*types.PerpetualPosition
		expectedAssetPositions     []*types.AssetPosition
		expectedSuccess            bool
		expectedSuccessPerUpdate   []types.UpdateResult
		expectedErr                error
		// Only contains the updated perpetual positions, to assert against the events included.
		expectedUpdatedPerpetualPositions     map[types.SubaccountId][]*types.PerpetualPosition
		expectedSubaccoundIdToFundingPayments map[types.SubaccountId]map[uint32]dtypes.SerializableInt
		expectedUpdatedAssetPositions         map[types.SubaccountId][]*types.AssetPosition
		msgSenderEnabled                      bool
	}{
		"one update to USDC asset position": {
			expectedQuoteBalance:     big.NewInt(100),
			expectedSuccess:          true,
			expectedSuccessPerUpdate: []types.UpdateResult{types.Success},
			updates: []types.Update{
				{
					AssetUpdates: testutil.CreateUsdcAssetUpdate(big.NewInt(100)),
				},
			},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(100), // 100 USDC
				},
			},
			expectedUpdatedAssetPositions: map[types.SubaccountId][]*types.AssetPosition{
				defaultSubaccountId: {
					{
						AssetId:  uint32(0),
						Quantums: dtypes.NewInt(100), // 100 USDC
					},
				},
			},
			msgSenderEnabled: true,
		},
		"one negative update to USDC asset position": {
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
				},
			},
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000), // 1 BTC
					FundingIndex: dtypes.NewInt(0),
				},
			},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(-100), // 100 USDC
				},
			},
			expectedUpdatedAssetPositions: map[types.SubaccountId][]*types.AssetPosition{
				defaultSubaccountId: {
					{
						AssetId:  uint32(0),
						Quantums: dtypes.NewInt(-100), // 100 USDC
					},
				},
			},
			updates: []types.Update{
				{
					AssetUpdates: testutil.CreateUsdcAssetUpdate(big.NewInt(-100)),
				},
			},
			msgSenderEnabled: true,
		},
		"one negative update to USDC asset position + persist unsettled negative funding": {
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
				},
			},
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000), // 1 BTC
					FundingIndex: dtypes.NewInt(-10),
				},
			},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(-2100), // 2100 USDC
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
						Quantums: dtypes.NewInt(-2100), // 2100 USDC
					},
				},
			},
			updates: []types.Update{
				{
					AssetUpdates: testutil.CreateUsdcAssetUpdate(big.NewInt(-100)),
				},
			},
			msgSenderEnabled: true,
		},
		"one negative update to USDC asset position + persist unsettled positive funding": {
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
				},
			},
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(500_000), // 1 BTC
					FundingIndex: dtypes.NewInt(-17),
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
					AssetUpdates: testutil.CreateUsdcAssetUpdate(big.NewInt(-100)),
				},
			},
			msgSenderEnabled: true,
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
				},
			},
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000), // 1 BTC
					FundingIndex: dtypes.NewInt(0),
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
			msgSenderEnabled: true,
		},
		"multiple updates to same account not allowed": {
			expectedQuoteBalance:     big.NewInt(0),
			expectedErr:              types.ErrNonUniqueUpdatesSubaccount,
			expectedSuccess:          false,
			expectedSuccessPerUpdate: nil,
			updates: []types.Update{
				{
					AssetUpdates: testutil.CreateUsdcAssetUpdate(big.NewInt(-100)),
				},
				{
					AssetUpdates: testutil.CreateUsdcAssetUpdate(big.NewInt(-100)),
				},
			},
			msgSenderEnabled: true,
		},
		"update increases position size": {
			assetPositions:           testutil.CreateUsdcAssetPosition(big.NewInt(25_000_000_000)), // $25,000
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
				},
			},
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(150_000_000), // 1.5 BTC
					FundingIndex: dtypes.NewInt(0),
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
					AssetUpdates: testutil.CreateUsdcAssetUpdate(big.NewInt(-25_000_000_000)), // -$25,000
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
			assetPositions:           testutil.CreateUsdcAssetPosition(big.NewInt(25_000_000_000)), // $25,000
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
				},
			},
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(50_000_000), // .50 BTC
					FundingIndex: dtypes.NewInt(0),
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
					AssetUpdates: testutil.CreateUsdcAssetUpdate(big.NewInt(25_000_000_000)), // $25,000
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
			assetPositions:           testutil.CreateUsdcAssetPosition(big.NewInt(25_000_000_000)), // $25,000
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
					AssetUpdates: testutil.CreateUsdcAssetUpdate(big.NewInt(50_000_000_000)), // $50,000
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(0),
							BigQuantumsDelta: big.NewInt(-100_000_000), // -1 BTC
						},
					},
				},
			},
			msgSenderEnabled: true,
		},
		"update closes short position": {
			assetPositions:           testutil.CreateUsdcAssetPosition(big.NewInt(100_000_000_000)), // $100,000
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
					AssetUpdates: testutil.CreateUsdcAssetUpdate(big.NewInt(-50_000_000_000)), // -$50,000
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(0),
							BigQuantumsDelta: big.NewInt(100_000_000), // 1 BTC
						},
					},
				},
			},
			msgSenderEnabled: true,
		},
		"update closes 2nd position and updates 1st": {
			assetPositions:           testutil.CreateUsdcAssetPosition(big.NewInt(100_000_000_000)), // $100,000
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
				},
				{
					PerpetualId:  uint32(1),
					Quantums:     dtypes.NewInt(-1_000_000_000_000_000_000), // -1 ETH
					FundingIndex: dtypes.NewInt(0),
				},
			},
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(-200_000_000), // -2 BTC
					FundingIndex: dtypes.NewInt(0),
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
			msgSenderEnabled: true,
		},
		"update closes first asset position and updates 2nd": {
			assets: []*asstypes.Asset{
				constants.BtcUsd,
			},
			assetPositions: append(
				testutil.CreateUsdcAssetPosition(big.NewInt(100_000_000_000)), // $100,000
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
							AssetId:          asstypes.AssetUsdc.Id,
							BigQuantumsDelta: big.NewInt(100_000_000_000),
						},
						{
							AssetId:          constants.BtcUsd.Id,
							BigQuantumsDelta: big.NewInt(-50_000),
						},
					},
				},
			},
			msgSenderEnabled: true,
		},
		"update closes first 1 positions and updates 2nd": {
			assetPositions:           testutil.CreateUsdcAssetPosition(big.NewInt(100_000_000_000)), // $100,000
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
				},
				{
					PerpetualId:  uint32(1),
					Quantums:     dtypes.NewInt(-1_000_000_000), // -1 ETH
					FundingIndex: dtypes.NewInt(0),
				},
			},
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(1),
					Quantums:     dtypes.NewInt(-2_000_000_000), // -2 ETH
					FundingIndex: dtypes.NewInt(0),
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
					AssetUpdates: testutil.CreateUsdcAssetUpdate(big.NewInt(-50_000_000_000)), // -$50,000
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
			msgSenderEnabled: true,
		},
		"update opens new long position, uses current perpetual funding index": {
			assetPositions:           testutil.CreateUsdcAssetPosition(big.NewInt(100_000_000_000)), // $100,000
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
					AssetUpdates: testutil.CreateUsdcAssetUpdate(big.NewInt(-50_000_000_000)), // -$50,000
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
			assetPositions:           testutil.CreateUsdcAssetPosition(big.NewInt(100_000_000_000)), // $100,000
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
					AssetUpdates: testutil.CreateUsdcAssetUpdate(big.NewInt(50_000_000_000)), // $50,000
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
			assetPositions:           testutil.CreateUsdcAssetPosition(big.NewInt(100_000_000_000)), // $100,000
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
				},
			},
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000), // 1 BTC
					FundingIndex: dtypes.NewInt(0),
				},
				{
					PerpetualId:  uint32(1),
					Quantums:     dtypes.NewInt(1_000_000_000), // 1 ETH
					FundingIndex: dtypes.NewInt(0),
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
			msgSenderEnabled: true,
		},
		// TODO(DEC-581): add similar test case for multi-collateral asset support.
		"update eth position from long to short with existing btc position": {
			assetPositions:           testutil.CreateUsdcAssetPosition(big.NewInt(100_000_000_000)), // $100,000
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
				},
				{
					PerpetualId:  uint32(1),
					Quantums:     dtypes.NewInt(500_000_000), // 5 ETH
					FundingIndex: dtypes.NewInt(0),
				},
			},
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000), // 1 BTC
					FundingIndex: dtypes.NewInt(0),
				},
				{
					PerpetualId:  uint32(1),
					Quantums:     dtypes.NewInt(-500_000_000), // -5 ETH
					FundingIndex: dtypes.NewInt(0),
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
			msgSenderEnabled: true,
		},
		"update opens new long eth position with existing btc and sol position": {
			assetPositions:           testutil.CreateUsdcAssetPosition(big.NewInt(100_000_000_000)), // $100,000
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
				},
				{
					PerpetualId:  uint32(2),
					Quantums:     dtypes.NewInt(1_000_000_000), // 1 SOL
					FundingIndex: dtypes.NewInt(0),
				},
			},
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000), // 1 BTC
					FundingIndex: dtypes.NewInt(0),
				},
				{
					PerpetualId:  uint32(1),
					Quantums:     dtypes.NewInt(1_000_000_000), // 1 ETH
					FundingIndex: dtypes.NewInt(0),
				},
				{
					PerpetualId:  uint32(2),
					Quantums:     dtypes.NewInt(1_000_000_000), // 1 SOL
					FundingIndex: dtypes.NewInt(0),
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
			msgSenderEnabled: true,
		},
		"update opens new long btc position with existing eth and sol position": {
			assetPositions:           testutil.CreateUsdcAssetPosition(big.NewInt(100_000_000_000)), // $100,000
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
				},
				{
					PerpetualId:  uint32(2),
					Quantums:     dtypes.NewInt(1_000_000_000), // 1 SOL
					FundingIndex: dtypes.NewInt(0),
				},
			},
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000), // 1 BTC
					FundingIndex: dtypes.NewInt(0),
				},
				{
					PerpetualId:  uint32(1),
					Quantums:     dtypes.NewInt(1_000_000_000), // 1 ETH
					FundingIndex: dtypes.NewInt(0),
				},
				{
					PerpetualId:  uint32(2),
					Quantums:     dtypes.NewInt(1_000_000_000), // 1 SOL
					FundingIndex: dtypes.NewInt(0),
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
			msgSenderEnabled: true,
		},
		"update opens new long eth position with existing unsettled sol position": {
			assetPositions:           testutil.CreateUsdcAssetPosition(big.NewInt(100_000_000_000)), // $100,000
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
				},
			},
			expectedPerpetualPositions: []*types.PerpetualPosition{
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
			msgSenderEnabled: true,
		},
		"provides out-of-order updates (not ordered by PerpetualId)": {
			assetPositions:           testutil.CreateUsdcAssetPosition(big.NewInt(100_000_000_000)), // $100,000
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
				},
				{
					PerpetualId:  uint32(1),
					Quantums:     dtypes.NewInt(1_000_000_000), // 1 ETH
					FundingIndex: dtypes.NewInt(0),
				},
				{
					PerpetualId:  uint32(2),
					Quantums:     dtypes.NewInt(1_000_000_000), // 1 SOL
					FundingIndex: dtypes.NewInt(0),
				},
			},
			expectedPerpetualPositions: []*types.PerpetualPosition{
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
			msgSenderEnabled: true,
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
					AssetUpdates: testutil.CreateUsdcAssetUpdate(big.NewInt(500_000_000)), // $500
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(0),
							BigQuantumsDelta: big.NewInt(100_000_000), // 1 BTC
						},
					},
				},
				{
					AssetUpdates: testutil.CreateUsdcAssetUpdate(big.NewInt(100_000_000)), // $100
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(1),
							BigQuantumsDelta: big.NewInt(-1_000_000_000), // -1 ETH
						},
					},
				},
			},
			msgSenderEnabled: true,
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
				},
			},
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(1_000_000), // 0.01 BTC
					FundingIndex: dtypes.NewInt(0),
				},
			},
			updates: []types.Update{
				{
					AssetUpdates: testutil.CreateUsdcAssetUpdate(big.NewInt(-50_000_000_000)), // -$50,000
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(0),
							BigQuantumsDelta: big.NewInt(100_000_000), // 1 BTC
						},
					},
				},
			},
			msgSenderEnabled: true,
		},
		"updates new USDC asset position which exceeds max uint64": {
			assetPositions:           testutil.CreateUsdcAssetPosition(new(big.Int).SetUint64(math.MaxUint64)),
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
					AssetUpdates: testutil.CreateUsdcAssetUpdate(big.NewInt(1)),
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
			msgSenderEnabled: true,
		},
		"new USDC asset position (including unsettled funding) size exceeds max uint64": {
			assetPositions: testutil.CreateUsdcAssetPosition(new(big.Int).SetUint64(math.MaxUint64 - 5)),
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
				},
			},
			updates: []types.Update{
				{
					AssetUpdates: testutil.CreateUsdcAssetUpdate(big.NewInt(3)),
				},
			},
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(1_000_000), // 0.01 BTC
					FundingIndex: dtypes.NewInt(-10),       // indexDelta=-3, settlement=3
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
			msgSenderEnabled: true,
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
			msgSenderEnabled: true,
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
				},
			},
			updates: []types.Update{
				{
					AssetUpdates: testutil.CreateUsdcAssetUpdate(big.NewInt(1)),
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
			msgSenderEnabled: true,
		},
		"update ETH position; start with BTC and ETH positions; both BTC and ETH positions have unsettled funding": {
			assetPositions:           testutil.CreateUsdcAssetPosition(big.NewInt(100_000_000_000)), // $100,000
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
				},
				{
					PerpetualId: uint32(1),
					Quantums:    dtypes.NewInt(-2_000_000_000), // -2 ETH
					// indexDelta=-2
					FundingIndex: dtypes.NewInt(-6),
				},
			},
			expectedPerpetualPositions: []*types.PerpetualPosition{
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
			msgSenderEnabled: true,
		},
		"update ETH position; start with BTC and ETH positions; only ETH position has unsettled funding": {
			assetPositions:           testutil.CreateUsdcAssetPosition(big.NewInt(100_000_000_000)), // $100,000
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
				},
				{
					PerpetualId: uint32(1),
					Quantums:    dtypes.NewInt(-2_000_000_000), // -2 ETH
					// indexDelta=-2
					FundingIndex: dtypes.NewInt(-6),
				},
			},
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(-100_000_000), // -1 BTC
					FundingIndex: dtypes.NewInt(0),
				},
				{
					PerpetualId:  uint32(1),
					Quantums:     dtypes.NewInt(-1_000_000_000), // -1 ETH
					FundingIndex: dtypes.NewInt(-8),
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
			msgSenderEnabled: true,
		},
		"update closes ETH position; start with BTC and ETH positions; both BTC and ETH positions have unsettled funding": {
			assetPositions:           testutil.CreateUsdcAssetPosition(big.NewInt(100_000_000_000)), // $100,000
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
				},
				{
					PerpetualId: uint32(1),
					Quantums:    dtypes.NewInt(-1_000_000_000), // -1 ETH
					// indexDelta=-2
					FundingIndex: dtypes.NewInt(-6),
				},
			},
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(-100_000_000), // -1 BTC
					FundingIndex: dtypes.NewInt(-10),
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
			msgSenderEnabled: true,
		},
		"update closes ETH position; start with ETH position; ETH position has no unsettled funding": {
			assetPositions:           testutil.CreateUsdcAssetPosition(big.NewInt(100_000_000_000)), // $100,000
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
			msgSenderEnabled: true,
		},
		"2 updates, 1 update involves not-updatable perp": {
			assetPositions: testutil.CreateUsdcAssetPosition(big.NewInt(1_000_000_000_000)),
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
					pricestest.WithPriceValue(0),
				),
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(100),
					Quantums:     dtypes.NewInt(1_000_000_000),
					FundingIndex: dtypes.NewInt(0),
				},
				{
					PerpetualId:  uint32(101),
					Quantums:     dtypes.NewInt(1_000_000_000),
					FundingIndex: dtypes.NewInt(0),
				},
			},
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(100),
					Quantums:     dtypes.NewInt(1_000_000_000),
					FundingIndex: dtypes.NewInt(0),
				},
				{
					PerpetualId:  uint32(101),
					Quantums:     dtypes.NewInt(1_000_000_000),
					FundingIndex: dtypes.NewInt(0),
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
			msgSenderEnabled: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctx, keeper, pricesKeeper, perpetualsKeeper, _, _, assetsKeeper, _ := testutil.SubaccountsKeepers(
				t,
				tc.msgSenderEnabled,
			)
			ctx = ctx.WithTxBytes(constants.TestTxBytes)
			testutil.CreateTestMarkets(t, ctx, pricesKeeper)
			testutil.CreateTestLiquidityTiers(t, ctx, perpetualsKeeper)

			for _, m := range tc.marketParamPrices {
				_, err := pricesKeeper.CreateMarket(
					ctx,
					m.Param,
					m.Price,
				)
				require.NoError(t, err)
			}

			// Always creates USDC asset first
			require.NoError(t, testutil.CreateUsdcAsset(ctx, assetsKeeper))
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

			subaccount := createNSubaccount(keeper, ctx, 1, big.NewInt(1_000))[0]
			subaccount.PerpetualPositions = tc.perpetualPositions
			subaccount.AssetPositions = tc.assetPositions
			keeper.SetSubaccount(ctx, subaccount)
			subaccountId := *subaccount.Id

			for i, u := range tc.updates {
				if u.SubaccountId == (types.SubaccountId{}) {
					u.SubaccountId = subaccountId
				}
				tc.updates[i] = u
			}

			success, successPerUpdate, err := keeper.UpdateSubaccounts(ctx, tc.updates)
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
				)
			} else {
				assertSubaccountUpdateEventsNotInIndexerBlock(
					t,
					keeper,
					ctx,
				)
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
		})
	}
}

func TestCanUpdateSubaccounts(t *testing.T) {
	tests := map[string]struct {
		// State.
		perpetuals        []perptypes.Perpetual
		assets            []*asstypes.Asset
		marketParamPrices []pricestypes.MarketParamPrice

		// Subaccount state.
		useEmptySubaccount bool
		perpetualPositions []*types.PerpetualPosition
		assetPositions     []*types.AssetPosition

		// Updates.
		updates []types.Update

		// Expectations.
		expectedSuccess          bool
		expectedSuccessPerUpdate []types.UpdateResult
		expectedErr              error
	}{
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
		"new USDC asset position exceeds max uint64": {
			assetPositions: testutil.CreateUsdcAssetPosition(new(big.Int).SetUint64(math.MaxUint64)),
			updates: []types.Update{
				{
					AssetUpdates: testutil.CreateUsdcAssetUpdate(big.NewInt(1)),
				},
			},
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
					AssetUpdates: testutil.CreateUsdcAssetUpdate(big.NewInt(-50_000_000_000)), // -$50,000
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
				},
			},
			updates: []types.Update{
				{
					AssetUpdates: testutil.CreateUsdcAssetUpdate(big.NewInt(-50_000_000_000)), // -$50,000
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(0),
							BigQuantumsDelta: big.NewInt(100_000_000), // 1 BTC
						},
					},
				},
				{
					AssetUpdates: testutil.CreateUsdcAssetUpdate(big.NewInt(-49_999_000_000)), // -$49,999
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
			assetPositions:  testutil.CreateUsdcAssetPosition(big.NewInt(-496_000_000)), // -$496
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
				},
			},
			updates: []types.Update{
				{
					AssetUpdates: testutil.CreateUsdcAssetUpdate(big.NewInt(-1)), // -$0.000001
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
					AssetUpdates: testutil.CreateUsdcAssetUpdate(big.NewInt(500_000)), // $.50
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(0),
							BigQuantumsDelta: big.NewInt(0), // 0 BTC
						},
					},
				},
				{
					AssetUpdates: testutil.CreateUsdcAssetUpdate(big.NewInt(2_000_000)), // $2
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(0),
							BigQuantumsDelta: big.NewInt(0), // 0 BTC
						},
					},
				},
			},
		},
		"USDC asset position is negative but increasing when no positions are open": {
			assetPositions: testutil.CreateUsdcAssetPosition(big.NewInt(-10)),
			updates: []types.Update{
				{
					AssetUpdates: testutil.CreateUsdcAssetUpdate(big.NewInt(1)), // $.000001
				},
			},
			expectedSuccess: true,
			expectedSuccessPerUpdate: []types.UpdateResult{
				types.Success,
			},
		},
		"USDC asset position is negative but unchanging when no positions are open": {
			assetPositions: testutil.CreateUsdcAssetPosition(big.NewInt(-10)),
			updates: []types.Update{
				{
					AssetUpdates: testutil.CreateUsdcAssetUpdate(big.NewInt(0)), // $0
				},
			},
			expectedSuccess: false,
			expectedSuccessPerUpdate: []types.UpdateResult{
				types.StillUndercollateralized,
			},
		},
		"USDC asset position decreases below zero when no positions are open": {
			updates: []types.Update{
				{
					AssetUpdates: testutil.CreateUsdcAssetUpdate(big.NewInt(-1)), // -$0.000001
				},
			},
			expectedSuccess: false,
			expectedSuccessPerUpdate: []types.UpdateResult{
				types.NewlyUndercollateralized,
			},
		},
		"USDC asset position decreases further below zero when no positions are open": {
			assetPositions: testutil.CreateUsdcAssetPosition(big.NewInt(-1)),
			updates: []types.Update{
				{
					AssetUpdates: testutil.CreateUsdcAssetUpdate(big.NewInt(-1)), // -$0.000001
				},
			},
			expectedSuccess: false,
			expectedSuccessPerUpdate: []types.UpdateResult{
				types.StillUndercollateralized,
			},
		},
		"two updates on different accounts, second account is new account": {
			assetPositions:           testutil.CreateUsdcAssetPosition(big.NewInt(50_000_000_000)), // $50,000
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
				},
			},
			updates: []types.Update{
				{
					AssetUpdates: testutil.CreateUsdcAssetUpdate(big.NewInt(-50_000_000_000)), // -$50,000
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
					AssetUpdates: testutil.CreateUsdcAssetUpdate(big.NewInt(-50_000_000_000)), // -$50,000
					PerpetualUpdates: []types.PerpetualUpdate{
						{
							PerpetualId:      uint32(0),
							BigQuantumsDelta: big.NewInt(100_000_000), // 1 BTC
						},
					},
				},
			},
		},
		"unsettled funding reduces USDC asset position to 1; further decrease USDC asset position, still collateralized": {
			assetPositions: testutil.CreateUsdcAssetPosition(big.NewInt(100)),
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(1_000_000), // 0.01 BTC,
					FundingIndex: dtypes.NewInt(-99),       // indexDelta=99, net settlement=-99
				},
			},
			updates: []types.Update{
				{
					AssetUpdates: testutil.CreateUsdcAssetUpdate(big.NewInt(-1)), // -$0.000001
				},
			},
			expectedSuccess: true,
			expectedSuccessPerUpdate: []types.UpdateResult{
				types.Success,
			},
		},
		"unsettled funding reduces USDC asset position to zero; further decrease USDC asset position, undercollateralized": {
			assetPositions: testutil.CreateUsdcAssetPosition(big.NewInt(100)),
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(1_000_000), // 0.01 BTC,
					FundingIndex: dtypes.NewInt(-100),      // indexDelta=100, net settlement=-100
				},
			},
			updates: []types.Update{
				{
					AssetUpdates: testutil.CreateUsdcAssetUpdate(big.NewInt(-1)), // -$0.000001
				},
			},
			expectedSuccess: false,
			expectedSuccessPerUpdate: []types.UpdateResult{
				types.NewlyUndercollateralized,
			},
		},
		"unsettled funding makes position undercollateralized": {
			assetPositions: testutil.CreateUsdcAssetPosition(big.NewInt(200)),
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(1_000_000), // 0.01 BTC,
					FundingIndex: dtypes.NewInt(-200),      // indexDelta=200, net settlement=-200
				},
			},
			updates: []types.Update{
				{
					AssetUpdates: testutil.CreateUsdcAssetUpdate(big.NewInt(-1)), // -$0.000001
				},
			},
			expectedSuccess: false,
			expectedSuccessPerUpdate: []types.UpdateResult{
				types.NewlyUndercollateralized,
			},
		},
		"position was undercollateralized before update due to funding and still undercollateralized" +
			"after due to funding": {
			assetPositions: testutil.CreateUsdcAssetPosition(big.NewInt(199)),
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(1_000_000), // 0.01 BTC,
					FundingIndex: dtypes.NewInt(-200),      // indexDelta=200, net settlement=-200
				},
			},
			updates: []types.Update{
				{
					AssetUpdates: testutil.CreateUsdcAssetUpdate(big.NewInt(-1)), // -$0.000001
				},
			},
			expectedSuccess: false,
			expectedSuccessPerUpdate: []types.UpdateResult{
				types.StillUndercollateralized,
			},
		},
		"unsettled funding makes position with negative USDC asset position collateralized before update": {
			assetPositions: testutil.CreateUsdcAssetPosition(big.NewInt(-100)),
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(1_000_000), // 0.01 BTC,
					FundingIndex: dtypes.NewInt(100),       // indexDelta=-100, net settlement=100
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
		"adding unsettled funding to USDC asset position exceeds max uint64": {
			assetPositions: testutil.CreateUsdcAssetPosition(new(big.Int).SetUint64(math.MaxUint64 - 1)),
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(1_000_000), // 0.01 BTC,
					FundingIndex: dtypes.NewInt(100),       // indexDelta=-100, net settlement=100
				},
			},
			updates: []types.Update{
				{},
			},
			expectedSuccess:          true,
			expectedSuccessPerUpdate: []types.UpdateResult{types.Success},
		},
		"adding unsettled funding to USDC asset position exceeds negative max uint64": {
			assetPositions: testutil.CreateUsdcAssetPosition(
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
				},
			},
			updates: []types.Update{
				{},
			},
			expectedSuccess:          true,
			expectedSuccessPerUpdate: []types.UpdateResult{types.Success},
		},
		"adding unsettled funding, original USDC asset position and USDC asset position delta exceeds max int64": {
			assetPositions: testutil.CreateUsdcAssetPosition(
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
				},
			},
			updates: []types.Update{
				{
					AssetUpdates: testutil.CreateUsdcAssetUpdate(big.NewInt(3)), // $3
				},
			},
			expectedSuccess:          true,
			expectedSuccessPerUpdate: []types.UpdateResult{types.Success},
		},
		"2 updates, 1 update involves not-updatable perp": {
			assetPositions: testutil.CreateUsdcAssetPosition(big.NewInt(1_000_000_000_000)),
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
					pricestest.WithPriceValue(0),
				),
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(100),
					Quantums:     dtypes.NewInt(1_000_000_000),
					FundingIndex: dtypes.NewInt(0),
				},
				{
					PerpetualId:  uint32(101),
					Quantums:     dtypes.NewInt(1_000_000_000),
					FundingIndex: dtypes.NewInt(0),
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
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctx, keeper, pricesKeeper, perpetualsKeeper, _, _, assetsKeeper, _ := testutil.SubaccountsKeepers(t, true)
			testutil.CreateTestMarkets(t, ctx, pricesKeeper)
			testutil.CreateTestLiquidityTiers(t, ctx, perpetualsKeeper)

			require.NoError(t, testutil.CreateUsdcAsset(ctx, assetsKeeper))
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

			for i, u := range tc.updates {
				if u.SubaccountId == (types.SubaccountId{}) {
					u.SubaccountId = subaccountId
				}
				tc.updates[i] = u
			}

			success, successPerUpdate, err := keeper.CanUpdateSubaccounts(ctx, tc.updates)
			if tc.expectedErr != nil {
				require.ErrorIs(t, tc.expectedErr, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedSuccessPerUpdate, successPerUpdate)
				require.Equal(t, tc.expectedSuccess, success)
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
		"non-negative USDC asset position": {
			assetPositions:        testutil.CreateUsdcAssetPosition(big.NewInt(123_456)),
			expectedNetCollateral: big.NewInt(123_456),
		},
		"negative USDC asset position": {
			assetPositions:        testutil.CreateUsdcAssetPosition(big.NewInt(-123_456)),
			expectedNetCollateral: big.NewInt(-123_456),
		},
		"USDC asset position with update": {
			assetPositions:        testutil.CreateUsdcAssetPosition(big.NewInt(-123_456)),
			expectedNetCollateral: big.NewInt(0),
			assetUpdates:          testutil.CreateUsdcAssetUpdate(big.NewInt(123_456)),
		},
		"single perpetual and USDC asset position": {
			assetPositions:        testutil.CreateUsdcAssetPosition(big.NewInt(10_000_000_001)), // $10,000.000001
			expectedNetCollateral: big.NewInt(60_000_000_001),                                   // $60,000.000001
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_NoMarginRequirement,
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000), // 1 BTC
					FundingIndex: dtypes.NewInt(0),
				},
			},
		},
		"single perpetual, USDC asset position and unsettled funding (long)": {
			assetPositions:        testutil.CreateUsdcAssetPosition(big.NewInt(10_000_000_001)), // $10,000.000001
			expectedNetCollateral: big.NewInt(60_006_250_001),                                   // $60,006.250001
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_NoMarginRequirement,
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000), // 1 BTC
					FundingIndex: dtypes.NewInt(62500),       // 0.0125% rate at BTC=50,000 USDC
				},
			},
		},
		"single perpetual, USDC asset position and unsettled funding (short)": {
			assetPositions:        testutil.CreateUsdcAssetPosition(big.NewInt(-10_000_000_001)), // -$10,000.000001
			expectedNetCollateral: big.NewInt(-60_006_250_001),                                   // -$60,006.250001
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_NoMarginRequirement,
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(-100_000_000), // 1 BTC
					FundingIndex: dtypes.NewInt(62500),        // 0.0125% rate at BTC=50,000 USDC
				},
			},
		},
		"non-existing perpetual heled by subaccount (should never happen)": {
			assetPositions: testutil.CreateUsdcAssetPosition(
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
					FundingIndex: dtypes.NewInt(62500),        // 0.0125% rate at BTC=50,000 USDC
				},
			},
			expectedErr: perptypes.ErrPerpetualDoesNotExist,
		},
		"USDC asset position update underflows uint64": {
			assetPositions: testutil.CreateUsdcAssetPosition(
				constants.BigNegMaxUint64(),
			),
			assetUpdates: testutil.CreateUsdcAssetUpdate(big.NewInt(-1)),
		},
		"USDC asset position update overflows uint64": {
			assetPositions: testutil.CreateUsdcAssetPosition(
				new(big.Int).SetUint64(math.MaxUint64),
			),
			assetUpdates: testutil.CreateUsdcAssetUpdate(big.NewInt(1)),
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
			assetPositions:        testutil.CreateUsdcAssetPosition(big.NewInt(10_000_000_001)), // $10,000.000001
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
		"single perpetual with USDC asset position and positive update to perpetual": {
			assetPositions:        testutil.CreateUsdcAssetPosition(big.NewInt(10_000_000_001)), // $10,000.000001
			expectedNetCollateral: big.NewInt(110_000_000_001),                                  // $110,000.000001
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_NoMarginRequirement,
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000), // 1 BTC
					FundingIndex: dtypes.NewInt(0),
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
				},
			},
			perpetualUpdates: []types.PerpetualUpdate{
				{
					PerpetualId:      uint32(0),
					BigQuantumsDelta: big.NewInt(-1),
				},
			},
		},
		"single perpetual with USDC asset position and negative update to perpetual": {
			assetPositions:        testutil.CreateUsdcAssetPosition(big.NewInt(10_000_000_001)), // $10,000.000001
			expectedNetCollateral: big.NewInt(10_000_000_001),                                   // $10,000.000001
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_NoMarginRequirement,
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000), // 1 BTC
					FundingIndex: dtypes.NewInt(0),
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
				&constants.Usdc_Asset_100_000,
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
			assetUpdates:              testutil.CreateUsdcAssetUpdate(big.NewInt(1_000_000)),
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
			assetPositions: testutil.CreateUsdcAssetPosition(big.NewInt(1000000)),
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
				},
				{
					PerpetualId:  uint32(1),
					Quantums:     dtypes.NewInt(500_000_000), // .5 ETH
					FundingIndex: dtypes.NewInt(0),
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
			expectedErr: asstypes.ErrNotImplementedMargin,
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
			ctx, keeper, pricesKeeper, perpetualsKeeper, _, _, assetsKeeper, _ := testutil.SubaccountsKeepers(t, true)
			testutil.CreateTestMarkets(t, ctx, pricesKeeper)
			testutil.CreateTestLiquidityTiers(t, ctx, perpetualsKeeper)

			require.NoError(t, testutil.CreateUsdcAsset(ctx, assetsKeeper))
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
