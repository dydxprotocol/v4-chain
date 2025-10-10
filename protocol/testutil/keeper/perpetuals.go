package keeper

import (
	"fmt"
	"testing"

	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/gogoproto/proto"

	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	marketmapkeeper "github.com/dydxprotocol/slinky/x/marketmap/keeper"
	pricefeedserver_types "github.com/dydxprotocol/v4-chain/protocol/daemons/server/types/pricefeed"
	indexerevents "github.com/dydxprotocol/v4-chain/protocol/indexer/events"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	assetskeeper "github.com/dydxprotocol/v4-chain/protocol/x/assets/keeper"
	delaymsgmoduletypes "github.com/dydxprotocol/v4-chain/protocol/x/delaymsg/types"
	epochskeeper "github.com/dydxprotocol/v4-chain/protocol/x/epochs/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/perpetuals"
	"github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	priceskeeper "github.com/dydxprotocol/v4-chain/protocol/x/prices/keeper"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	"github.com/stretchr/testify/require"
)

type PerpKeepersTestContext struct {
	Ctx               sdk.Context
	PricesKeeper      *priceskeeper.Keeper
	IndexPriceCache   *pricefeedserver_types.MarketToExchangePrices
	AssetsKeeper      *assetskeeper.Keeper
	EpochsKeeper      *epochskeeper.Keeper
	PerpetualsKeeper  *keeper.Keeper
	MarketMapKeeper   *marketmapkeeper.Keeper
	StoreKey          storetypes.StoreKey
	MemKey            storetypes.StoreKey
	Cdc               *codec.ProtoCodec
	MockTimeProvider  *mocks.TimeProvider
	TransientStoreKey storetypes.StoreKey
}

func PerpetualsKeepers(
	t testing.TB,
) (pc PerpKeepersTestContext) {
	return PerpetualsKeepersWithClobHelpers(
		t,
		nil,
	)
}

func PerpetualsKeepersWithClobHelpers(
	t testing.TB,
	clobKeeper types.PerpetualsClobKeeper,
) (pc PerpKeepersTestContext) {
	pc.Ctx = initKeepers(
		t, func(
			db *dbm.MemDB,
			registry codectypes.InterfaceRegistry,
			cdc *codec.ProtoCodec,
			stateStore storetypes.CommitMultiStore,
			transientStoreKey storetypes.StoreKey,
		) []GenesisInitializer {
			// Define necessary keepers here for unit tests
			pc.EpochsKeeper, _ = createEpochsKeeper(stateStore, db, cdc)

			accountsKeeper, _ := createAccountKeeper(
				stateStore,
				db,
				cdc,
				registry)
			bankKeeper, _ := createBankKeeper(stateStore, db, cdc, accountsKeeper)
			stakingKeeper, _ := createStakingKeeper(
				stateStore,
				db,
				cdc,
				accountsKeeper,
				bankKeeper,
			)
			statsKeeper, _ := createStatsKeeper(
				stateStore,
				pc.EpochsKeeper,
				db,
				cdc,
				stakingKeeper,
			)
			affiliatesKeeper, _ := createAffiliatesKeeper(stateStore, db, cdc, statsKeeper, transientStoreKey, true)
			vaultKeeper, _ := createVaultKeeper(stateStore, db, cdc, transientStoreKey)
			feetiersKeeper, _ := createFeeTiersKeeper(stateStore, statsKeeper, vaultKeeper, affiliatesKeeper, db, cdc)
			revShareKeeper, _, _ := createRevShareKeeper(stateStore, db, cdc, affiliatesKeeper, feetiersKeeper, statsKeeper)
			pc.MarketMapKeeper, _ = createMarketMapKeeper(stateStore, db, cdc)
			pc.PricesKeeper, _, pc.IndexPriceCache, pc.MockTimeProvider = createPricesKeeper(
				stateStore,
				db,
				cdc,
				transientStoreKey,
				revShareKeeper,
				pc.MarketMapKeeper,
			)
			pc.PerpetualsKeeper, pc.StoreKey = createPerpetualsKeeperWithClobHelpers(
				stateStore,
				db,
				cdc,
				pc.PricesKeeper,
				pc.EpochsKeeper,
				clobKeeper,
				transientStoreKey,
			)
			pc.TransientStoreKey = transientStoreKey
			return []GenesisInitializer{pc.PricesKeeper, pc.PerpetualsKeeper, pc.MarketMapKeeper}
		},
	)

	// Mock time provider response for market creation.
	pc.MockTimeProvider.On("Now").Return(constants.TimeT)

	// Initialize perpetuals module parameters to default genesis values.
	perpetuals.InitGenesis(pc.Ctx, *pc.PerpetualsKeeper, constants.Perpetuals_GenesisState_ParamsOnly)

	return pc
}

func createPerpetualsKeeperWithClobHelpers(
	stateStore storetypes.CommitMultiStore,
	db *dbm.MemDB,
	cdc *codec.ProtoCodec,
	pk *priceskeeper.Keeper,
	ek *epochskeeper.Keeper,
	pck types.PerpetualsClobKeeper,
	transientStoreKey storetypes.StoreKey,
) (*keeper.Keeper, storetypes.StoreKey) {
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)

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
			lib.GovModuleAddress.String(),
			delaymsgmoduletypes.ModuleAddress.String(),
		},
		transientStoreKey,
	)

	k.SetClobKeeper(pck)

	return k, storeKey
}

func createPerpetualsKeeper(
	stateStore storetypes.CommitMultiStore,
	db *dbm.MemDB,
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
	t testing.TB,
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

func CreateTestPerpetuals(t testing.TB, ctx sdk.Context, k *keeper.Keeper) {
	for _, p := range constants.TestMarketPerpetuals {
		_, err := k.CreatePerpetual(
			ctx,
			p.Params.Id,
			p.Params.Ticker,
			p.Params.MarketId,
			p.Params.AtomicResolution,
			p.Params.DefaultFundingPpm,
			p.Params.LiquidityTier,
			p.Params.MarketType,
		)
		require.NoError(t, err)
	}
}

func CreateTestLiquidityTiers(t testing.TB, ctx sdk.Context, k *keeper.Keeper) {
	for _, l := range constants.LiquidityTiers {
		_, err := k.SetLiquidityTier(
			ctx,
			l.Id,
			l.Name,
			l.InitialMarginPpm,
			l.MaintenanceFractionPpm,
			l.ImpactNotional,
			l.OpenInterestLowerCap,
			l.OpenInterestUpperCap,
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
		var liquidityTierEvent indexerevents.LiquidityTierUpsertEventV1
		err := proto.Unmarshal(event.DataBytes, &liquidityTierEvent)
		if err != nil {
			panic(err)
		}
		liquidityTierEvents = append(liquidityTierEvents, &liquidityTierEvent)
	}
	return liquidityTierEvents
}

func CreateNPerpetuals(
	t testing.TB,
	ctx sdk.Context,
	keeper *keeper.Keeper,
	pricesKeeper *priceskeeper.Keeper,
	n int,
) ([]types.Perpetual, error) {
	items := make([]types.Perpetual, n)
	allLiquidityTiers := keeper.GetAllLiquidityTiers(ctx)
	require.Greater(t, len(allLiquidityTiers), 0)

	CreateNMarkets(t, ctx, pricesKeeper, n)
	for i := range items {
		var defaultFundingPpm int32
		marketType := types.PerpetualMarketType_PERPETUAL_MARKET_TYPE_CROSS
		if i%3 == 0 {
			defaultFundingPpm = 1
			marketType = types.PerpetualMarketType_PERPETUAL_MARKET_TYPE_ISOLATED
		} else if i%3 == 1 {
			defaultFundingPpm = -1
		} else {
			defaultFundingPpm = 0
		}

		perpetual, err := keeper.CreatePerpetual(
			ctx,
			uint32(i),            // Id
			fmt.Sprintf("%v", i), // Ticker
			uint32(i),            // MarketId
			int32(i),             // AtomicResolution
			defaultFundingPpm,    // DefaultFundingPpm
			allLiquidityTiers[i%len(allLiquidityTiers)].Id, // LiquidityTier
			marketType, // MarketType
		)
		if err != nil {
			return items, err
		}

		items[i] = perpetual
	}
	return items, nil
}

func CreateLiquidityTiersAndNPerpetuals(
	t testing.TB,
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
	t testing.TB,
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
			perp.Params.MarketType,
		)
		require.NoError(t, err)
	}
}
