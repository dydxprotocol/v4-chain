package keeper

import (
	"fmt"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/indexer/common"
	indexerevents "github.com/dydxprotocol/v4-chain/protocol/indexer/events"

	"github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"

	tmdb "github.com/cometbft/cometbft-db"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	delaymsgmoduletypes "github.com/dydxprotocol/v4-chain/protocol/x/delaymsg/types"
	epochskeeper "github.com/dydxprotocol/v4-chain/protocol/x/epochs/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/perpetuals"
	"github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	priceskeeper "github.com/dydxprotocol/v4-chain/protocol/x/prices/keeper"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	"github.com/stretchr/testify/require"
)

func PerpetualsKeepers(
	t testing.TB,
) (
	ctx sdk.Context,
	keeper *keeper.Keeper,
	pricesKeeper *priceskeeper.Keeper,
	epochsKeeper *epochskeeper.Keeper,
	storeKey storetypes.StoreKey,
) {
	return PerpetualsKeepersWithClobHelpers(
		t,
		nil,
	)
}

func PerpetualsKeepersWithClobHelpers(
	t testing.TB,
	clobKeeper types.PerpetualsClobKeeper,
) (
	ctx sdk.Context,
	keeper *keeper.Keeper,
	pricesKeeper *priceskeeper.Keeper,
	epochsKeeper *epochskeeper.Keeper,
	storeKey storetypes.StoreKey,
) {
	var mockTimeProvider *mocks.TimeProvider
	ctx = initKeepers(t, func(
		db *tmdb.MemDB,
		registry codectypes.InterfaceRegistry,
		cdc *codec.ProtoCodec,
		stateStore storetypes.CommitMultiStore,
		transientStoreKey storetypes.StoreKey,
	) []GenesisInitializer {
		// Define necessary keepers here for unit tests
		pricesKeeper, _, _, _, mockTimeProvider = createPricesKeeper(stateStore, db, cdc, transientStoreKey)
		epochsKeeper, _ = createEpochsKeeper(stateStore, db, cdc)
		keeper, storeKey = createPerpetualsKeeperWithClobHelpers(
			stateStore,
			db,
			cdc,
			pricesKeeper,
			epochsKeeper,
			clobKeeper,
			transientStoreKey,
		)

		return []GenesisInitializer{pricesKeeper, keeper}
	})

	// Mock time provider response for market creation.
	mockTimeProvider.On("Now").Return(constants.TimeT)

	// Initialize perpetuals module parameters to default genesis values.
	perpetuals.InitGenesis(ctx, *keeper, constants.Perpetuals_GenesisState_ParamsOnly)

	return ctx, keeper, pricesKeeper, epochsKeeper, storeKey
}

func createPerpetualsKeeperWithClobHelpers(
	stateStore storetypes.CommitMultiStore,
	db *tmdb.MemDB,
	cdc *codec.ProtoCodec,
	pk *priceskeeper.Keeper,
	ek *epochskeeper.Keeper,
	pck types.PerpetualsClobKeeper,
	transientStoreKey storetypes.StoreKey,
) (*keeper.Keeper, storetypes.StoreKey) {
	storeKey := sdk.NewKVStoreKey(types.StoreKey)

	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)

	mockMsgSender := &mocks.IndexerMessageSender{}
	mockMsgSender.On("Enabled").Return(true)
	mockIndexerEventsManager := indexer_manager.NewIndexerEventManager(mockMsgSender, transientStoreKey, true)

	k := keeper.NewKeeper(
		cdc,
		storeKey,
		pk,
		ek,
		mockIndexerEventsManager,
		[]string{
			authtypes.NewModuleAddress(govtypes.ModuleName).String(),
			authtypes.NewModuleAddress(delaymsgmoduletypes.ModuleName).String(),
		},
	)

	k.SetClobKeeper(pck)

	return k, storeKey
}

func createPerpetualsKeeper(
	stateStore storetypes.CommitMultiStore,
	db *tmdb.MemDB,
	cdc *codec.ProtoCodec,
	pk *priceskeeper.Keeper,
	ek *epochskeeper.Keeper,
	transientStoreKey storetypes.StoreKey,
) (*keeper.Keeper, storetypes.StoreKey) {
	return createPerpetualsKeeperWithClobHelpers(stateStore, db, cdc, pk, ek, nil, transientStoreKey)
}

// PopulateTestPremiumStore populates either `PremiumVotes` (`isVote` is true) or
// `PremiumSamples` (`isVote` is false) for test.
// For each perpetual in the given perpetuals, insert the same list of testFundingSamples
// into state.
func PopulateTestPremiumStore(
	t *testing.T,
	ctx sdk.Context,
	k *keeper.Keeper,
	perpetuals []types.Perpetual,
	testFundingPremiums []int32,
	isVote bool,
) {
	for _, premiumPpm := range testFundingPremiums {
		newPremiums := make([]types.FundingPremium, len(perpetuals))
		for i, p := range perpetuals {
			newPremiums[i] = *types.NewFundingPremium(p.Params.Id, premiumPpm)
		}

		if isVote {
			err := k.AddPremiumVotes(ctx, newPremiums)
			require.NoError(t, err)
			return
		}

		err := k.AddPremiumSamples(ctx, newPremiums)
		require.NoError(t, err)
	}
}

func CreateTestLiquidityTiers(t *testing.T, ctx sdk.Context, k *keeper.Keeper) {
	for _, l := range constants.LiquidityTiers {
		_, err := k.CreateLiquidityTier(
			ctx,
			l.Name,
			l.InitialMarginPpm,
			l.MaintenanceFractionPpm,
			l.BasePositionNotional,
			l.ImpactNotional,
		)

		require.NoError(t, err)
	}
}

// GetLiquidityTierUpsertEventsFromIndexerBlock returns the liquidityTier upsert events in the
// Indexer Block event Kafka message.
// TODO(IND-365): Consider using generics here to reduce duplicated code.
func GetLiquidityTierUpsertEventsFromIndexerBlock(
	ctx sdk.Context,
	keeper *keeper.Keeper,
) []*indexerevents.LiquidityTierUpsertEventV1 {
	var liquidityTierEvents []*indexerevents.LiquidityTierUpsertEventV1
	block := keeper.GetIndexerEventManager().ProduceBlock(ctx)
	if block == nil {
		return liquidityTierEvents
	}
	for _, event := range block.Events {
		if event.Subtype != indexerevents.SubtypeLiquidityTier {
			continue
		}
		bytes := indexer_manager.GetBytesFromEventData(event.Data)
		unmarshaler := common.UnmarshalerImpl{}
		var liquidityTierEvent indexerevents.LiquidityTierUpsertEventV1
		err := unmarshaler.Unmarshal(bytes, &liquidityTierEvent)
		if err != nil {
			panic(err)
		}
		liquidityTierEvents = append(liquidityTierEvents, &liquidityTierEvent)
	}
	return liquidityTierEvents
}

func CreateNPerpetuals(
	t *testing.T,
	ctx sdk.Context,
	keeper *keeper.Keeper,
	pricesKeeper *priceskeeper.Keeper,
	n int,
) ([]types.Perpetual, error) {
	items := make([]types.Perpetual, n)
	numLiquidityTiers := keeper.GetNumLiquidityTiers(ctx)
	require.Greater(t, numLiquidityTiers, uint32(0))

	for i := range items {
		CreateNMarkets(t, ctx, pricesKeeper, n)

		var defaultFundingPpm int32
		if i%3 == 0 {
			defaultFundingPpm = 1
		} else if i%3 == 1 {
			defaultFundingPpm = -1
		} else {
			defaultFundingPpm = 0
		}

		perpetual, err := keeper.CreatePerpetual(
			ctx,
			uint32(i),                        // Id
			fmt.Sprintf("%v", i),             // Ticker
			uint32(i),                        // MarketId
			int32(i),                         // AtomicResolution
			defaultFundingPpm,                // DefaultFundingPpm
			uint32(i%int(numLiquidityTiers)), // LiquidityTier
		)
		if err != nil {
			return items, err
		}

		items[i] = perpetual
	}
	return items, nil
}

func CreateLiquidityTiersAndNPerpetuals(
	t *testing.T,
	ctx sdk.Context,
	keeper *keeper.Keeper,
	pricesKeeper *priceskeeper.Keeper,
	n int,
) []types.Perpetual {
	// Create liquidity tiers.
	CreateTestLiquidityTiers(t, ctx, keeper)
	// Create perpetuals.
	perpetuals, err := CreateNPerpetuals(t, ctx, keeper, pricesKeeper, n)
	require.NoError(t, err)
	return perpetuals
}

// CreateTestPricesAndPerpetualMarkets is a test utility function that creates list of given
// prices and perpetual markets in state.
func CreateTestPricesAndPerpetualMarkets(
	t *testing.T,
	ctx sdk.Context,
	perpKeeper *keeper.Keeper,
	pricesKeeper *priceskeeper.Keeper,
	perpetuals []types.Perpetual,
	markets []pricestypes.MarketParamPrice,
) {
	// Create liquidity tiers.
	CreateTestLiquidityTiers(t, ctx, perpKeeper)

	CreateTestPriceMarkets(t, ctx, pricesKeeper, markets)

	for _, perp := range perpetuals {
		_, err := perpKeeper.CreatePerpetual(
			ctx,
			perp.Params.Id,
			perp.Params.Ticker,
			perp.Params.MarketId,
			perp.Params.AtomicResolution,
			perp.Params.DefaultFundingPpm,
			perp.Params.LiquidityTier,
		)
		require.NoError(t, err)
	}
}
