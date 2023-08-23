package keeper

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"

	tmdb "github.com/cometbft/cometbft-db"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	epochskeeper "github.com/dydxprotocol/v4-chain/protocol/x/epochs/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/perpetuals"
	"github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	priceskeeper "github.com/dydxprotocol/v4-chain/protocol/x/prices/keeper"
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
	return PerpetualsKeepersWithPricePremiumGetter(
		t,
		nil,
	)
}

func PerpetualsKeepersWithPricePremiumGetter(
	t testing.TB,
	pricePremiumGetter types.PricePremiumGetter,
) (
	ctx sdk.Context,
	keeper *keeper.Keeper,
	pricesKeeper *priceskeeper.Keeper,
	epochsKeeper *epochskeeper.Keeper,
	storeKey storetypes.StoreKey,
) {
	ctx = initKeepers(t, func(
		db *tmdb.MemDB,
		registry codectypes.InterfaceRegistry,
		cdc *codec.ProtoCodec,
		stateStore storetypes.CommitMultiStore,
		transientStoreKey storetypes.StoreKey,
	) []GenesisInitializer {
		// Define necessary keepers here for unit tests
		pricesKeeper, _, _, _, _ = createPricesKeeper(stateStore, db, cdc, transientStoreKey)
		epochsKeeper, _ = createEpochsKeeper(stateStore, db, cdc)
		keeper, storeKey = createPerpetualsKeeperWithPricePremiumGetter(
			stateStore,
			db,
			cdc,
			pricesKeeper,
			epochsKeeper,
			pricePremiumGetter,
			transientStoreKey,
		)

		return []GenesisInitializer{pricesKeeper, keeper}
	})

	// Initialize perpetuals module parameters to default genesis values.
	perpetuals.InitGenesis(ctx, *keeper, constants.Perpetuals_GenesisState_ParamsOnly)

	return ctx, keeper, pricesKeeper, epochsKeeper, storeKey
}

func createPerpetualsKeeperWithPricePremiumGetter(
	stateStore storetypes.CommitMultiStore,
	db *tmdb.MemDB,
	cdc *codec.ProtoCodec,
	pk *priceskeeper.Keeper,
	ek *epochskeeper.Keeper,
	ppg types.PricePremiumGetter,
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
	)

	k.SetPricePremiumGetter(ppg)

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
	return createPerpetualsKeeperWithPricePremiumGetter(stateStore, db, cdc, pk, ek, nil, transientStoreKey)
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
			newPremiums[i] = *types.NewFundingPremium(p.Id, premiumPpm)
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
