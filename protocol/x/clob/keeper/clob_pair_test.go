package keeper_test

import (
	"fmt"
	"strconv"
	"testing"

	indexerevents "github.com/dydxprotocol/v4-chain/protocol/indexer/events"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	clobtest "github.com/dydxprotocol/v4-chain/protocol/testutil/clob"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/nullify"
	perptest "github.com/dydxprotocol/v4-chain/protocol/testutil/perpetuals"
	pricestest "github.com/dydxprotocol/v4-chain/protocol/testutil/prices"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/memclob"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/perpetuals"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/prices"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func TestCreatePerpetualClobPair_MultiplePerpetual(t *testing.T) {
	memClob := memclob.NewMemClobPriceTimePriority(false)
	mockIndexerEventManager := &mocks.IndexerEventManager{}
	ks := keepertest.NewClobKeepersTestContext(t, memClob, &mocks.BankKeeper{}, mockIndexerEventManager)

	prices.InitGenesis(ks.Ctx, *ks.PricesKeeper, constants.Prices_DefaultGenesisState)
	perpetuals.InitGenesis(ks.Ctx, *ks.PerpetualsKeeper, constants.Perpetuals_DefaultGenesisState)

	clobPairs := []types.ClobPair{
		constants.ClobPair_Btc,
		constants.ClobPair_Eth,
	}

	for i, clobPair := range clobPairs {
		mockIndexerEventManager.On("AddTxnEvent",
			ks.Ctx,
			indexerevents.SubtypePerpetualMarket,
			indexer_manager.GetB64EncodedEventMessage(
				indexerevents.NewPerpetualMarketCreateEvent(
					clobPair.MustGetPerpetualId(),
					clobPair.Id,
					constants.Perpetuals_DefaultGenesisState.Perpetuals[i].Params.Ticker,
					constants.Perpetuals_DefaultGenesisState.Perpetuals[i].Params.MarketId,
					clobPair.Status,
					clobPair.QuantumConversionExponent,
					constants.Perpetuals_DefaultGenesisState.Perpetuals[i].Params.AtomicResolution,
					clobPair.SubticksPerTick,
					clobPair.StepBaseQuantums,
					constants.Perpetuals_DefaultGenesisState.Perpetuals[i].Params.LiquidityTier,
				),
			),
		).Once().Return()
		//nolint: errcheck
		ks.ClobKeeper.CreatePerpetualClobPair(
			ks.Ctx,
			clobPair.Id,
			clobtest.MustPerpetualId(clobPair),
			satypes.BaseQuantums(clobPair.StepBaseQuantums),
			clobPair.QuantumConversionExponent,
			clobPair.SubticksPerTick,
			clobPair.Status,
		)
	}

	require.Equal(
		t,
		map[uint32][]types.ClobPairId{
			0: {constants.ClobPair_Btc.GetClobPairId()},
			1: {constants.ClobPair_Eth.GetClobPairId()},
		},
		ks.ClobKeeper.PerpetualIdToClobPairId,
	)
}

func TestCreatePerpetualClobPair_FailsWithPerpetualAssociatedWithExistingClobPair(t *testing.T) {
	memClob := memclob.NewMemClobPriceTimePriority(false)
	// Set up mock indexer event manager that accepts anything.
	mockIndexerEventManager := &mocks.IndexerEventManager{}
	mockIndexerEventManager.On("AddTxnEvent",
		mock.Anything, mock.Anything, mock.Anything,
	).Return()
	ks := keepertest.NewClobKeepersTestContext(
		t,
		memClob,
		&mocks.BankKeeper{},
		mockIndexerEventManager,
	)
	// Create test perpetual and market (id 0).
	keepertest.CreateTestPricesAndPerpetualMarkets(
		t,
		ks.Ctx,
		ks.PerpetualsKeeper,
		ks.PricesKeeper,
		[]perptypes.Perpetual{
			*perptest.GeneratePerpetual(perptest.WithId(0), perptest.WithMarketId(0)),
		},
		[]pricestypes.MarketParamPrice{
			*pricestest.GenerateMarketParamPrice(pricestest.WithId(0)),
		},
	)
	// Create test clob pair id 0, associated with perpetual id 0.
	keepertest.CreateTestClobPairs(
		t,
		ks.Ctx,
		ks.ClobKeeper,
		[]types.ClobPair{
			*clobtest.GenerateClobPair(clobtest.WithId(0), clobtest.WithPerpetualId(0)),
		},
	)

	// Create test clob pair id 1, associated with perpetual id 0.
	newClobPair := clobtest.GenerateClobPair(clobtest.WithId(1), clobtest.WithPerpetualId(0))
	_, err := ks.ClobKeeper.CreatePerpetualClobPair(
		ks.Ctx,
		newClobPair.Id,
		newClobPair.MustGetPerpetualId(),
		satypes.BaseQuantums(newClobPair.StepBaseQuantums),
		newClobPair.QuantumConversionExponent,
		newClobPair.SubticksPerTick,
		newClobPair.Status,
	)
	require.ErrorContains(
		t,
		err,
		"perpetual ID is already associated with an existing CLOB pair",
	)
}

func TestCreatePerpetualClobPair_FailsWithDuplicateClobPairId(t *testing.T) {
	memClob := memclob.NewMemClobPriceTimePriority(false)
	mockIndexerEventManager := &mocks.IndexerEventManager{}
	ks := keepertest.NewClobKeepersTestContext(
		t,
		memClob,
		&mocks.BankKeeper{},
		mockIndexerEventManager,
	)
	prices.InitGenesis(ks.Ctx, *ks.PricesKeeper, constants.Prices_DefaultGenesisState)
	perpetuals.InitGenesis(ks.Ctx, *ks.PerpetualsKeeper, constants.Perpetuals_DefaultGenesisState)

	// Read a new `ClobPair` and make sure it does not exist.
	_, err := ks.ClobKeeper.GetClobPairIdForPerpetual(ks.Ctx, 1)
	require.ErrorIs(t, err, types.ErrNoClobPairForPerpetual)

	// Write `ClobPair` to state, but don't call `keeper.createOrderbook`.
	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)
	store := prefix.NewStore(ks.Ctx.KVStore(ks.StoreKey), types.KeyPrefix(types.ClobPairKeyPrefix))

	// Write clob pair to state with clob pair id 0.
	b := cdc.MustMarshal(&constants.ClobPair_Btc)
	store.Set(types.ClobPairKey(
		types.ClobPairId(constants.ClobPair_Btc.Id),
	), b)

	clobPair := *clobtest.GenerateClobPair()

	mockIndexerEventManager.On("AddTxnEvent",
		ks.Ctx,
		indexerevents.SubtypePerpetualMarket,
		indexer_manager.GetB64EncodedEventMessage(
			indexerevents.NewPerpetualMarketCreateEvent(
				clobPair.MustGetPerpetualId(),
				clobPair.Id,
				constants.Perpetuals_DefaultGenesisState.Perpetuals[0].Params.Ticker,
				constants.Perpetuals_DefaultGenesisState.Perpetuals[0].Params.MarketId,
				clobPair.Status,
				clobPair.QuantumConversionExponent,
				constants.Perpetuals_DefaultGenesisState.Perpetuals[0].Params.AtomicResolution,
				clobPair.SubticksPerTick,
				clobPair.StepBaseQuantums,
				constants.Perpetuals_DefaultGenesisState.Perpetuals[0].Params.LiquidityTier,
			),
		),
	).Once().Return()

	_, err = ks.ClobKeeper.CreatePerpetualClobPair(
		ks.Ctx,
		clobPair.Id,
		clobtest.MustPerpetualId(clobPair),
		satypes.BaseQuantums(clobPair.StepBaseQuantums),
		clobPair.QuantumConversionExponent,
		clobPair.SubticksPerTick,
		clobPair.Status,
	)

	require.ErrorIs(
		t,
		err,
		types.ErrClobPairAlreadyExists,
	)
}

func TestCreatePerpetualClobPair(t *testing.T) {
	tests := map[string]struct {
		// CLOB pair.
		clobPair types.ClobPair

		// Expectations.
		expectedErr string
	}{
		"CLOB pair is valid": {
			clobPair: *clobtest.GenerateClobPair(),
		},
		"CLOB pair is invalid when the perpetual ID does not match an existing perpetual in the store": {
			clobPair: *clobtest.GenerateClobPair(clobtest.WithPerpetualMetadata(
				&types.ClobPair_PerpetualClobMetadata{
					PerpetualClobMetadata: &types.PerpetualClobMetadata{
						PerpetualId: 1000000,
					},
				},
			)),
			expectedErr: "has invalid perpetual.",
		},
		"CLOB pair is invalid when the step size is 0": {
			clobPair:    *clobtest.GenerateClobPair(clobtest.WithStepBaseQuantums(0)),
			expectedErr: "invalid ClobPair parameter: StepBaseQuantums must be > 0.",
		},
		"CLOB pair is invalid when the subticks per tick is 0": {
			clobPair:    *clobtest.GenerateClobPair(clobtest.WithSubticksPerTick(0)),
			expectedErr: "invalid ClobPair parameter: SubticksPerTick must be > 0.",
		},
		"CLOB pair is invalid when the status is unspecified": {
			clobPair:    *clobtest.GenerateClobPair(clobtest.WithStatus(types.ClobPair_STATUS_UNSPECIFIED)),
			expectedErr: "has unsupported status STATUS_UNSPECIFIED",
		},
		"CLOB pair status is not supported": {
			clobPair: *clobtest.GenerateClobPair(
				clobtest.WithStatus(types.ClobPair_STATUS_PAUSED),
			),
			expectedErr: "has unsupported status STATUS_PAUSED",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Boilerplate setup.
			memClob := memclob.NewMemClobPriceTimePriority(false)
			mockIndexerEventManager := &mocks.IndexerEventManager{}
			ks := keepertest.NewClobKeepersTestContext(t, memClob, &mocks.BankKeeper{}, mockIndexerEventManager)

			prices.InitGenesis(ks.Ctx, *ks.PricesKeeper, constants.Prices_DefaultGenesisState)
			perpetuals.InitGenesis(ks.Ctx, *ks.PerpetualsKeeper, constants.Perpetuals_DefaultGenesisState)
			// PerpetualMarketCreateEvents are emitted when initializing the genesis state, so we need to mock
			// the indexer event manager to expect these events.
			if tc.expectedErr == "" {
				perpetualId := clobtest.MustPerpetualId(tc.clobPair)
				perpetual := constants.Perpetuals_DefaultGenesisState.Perpetuals[perpetualId]
				mockIndexerEventManager.On("AddTxnEvent",
					ks.Ctx,
					indexerevents.SubtypePerpetualMarket,
					indexer_manager.GetB64EncodedEventMessage(
						indexerevents.NewPerpetualMarketCreateEvent(
							perpetualId,
							perpetualId,
							perpetual.Params.Ticker,
							perpetual.Params.MarketId,
							tc.clobPair.Status,
							tc.clobPair.QuantumConversionExponent,
							perpetual.Params.AtomicResolution,
							tc.clobPair.SubticksPerTick,
							tc.clobPair.StepBaseQuantums,
							perpetual.Params.LiquidityTier,
						),
					),
				).Return()
			}

			// Perform the method under test.
			createdClobPair, actualErr := ks.ClobKeeper.CreatePerpetualClobPair(
				ks.Ctx,
				tc.clobPair.Id,
				clobtest.MustPerpetualId(tc.clobPair),
				satypes.BaseQuantums(tc.clobPair.StepBaseQuantums),
				tc.clobPair.QuantumConversionExponent,
				tc.clobPair.SubticksPerTick,
				tc.clobPair.Status,
			)
			storedClobPair, found := ks.ClobKeeper.GetClobPair(ks.Ctx, types.ClobPairId(tc.clobPair.Id))

			if tc.expectedErr == "" {
				// A valid CLOB pair should not raise any validation errors.
				require.NoError(t, actualErr)

				// The CLOB pair returned should be identical to the test case.
				require.Equal(t, tc.clobPair, createdClobPair)

				// The CLOB pair should be able to be retrieved from the store.
				require.True(t, found)
				require.NotNil(t, storedClobPair)

				// The stored CLOB pair should be identical to the test case.
				require.Equal(t, tc.clobPair, storedClobPair)
			} else {
				// The create method should have returned a validation error matching the test case.
				require.Error(t, actualErr)
				require.ErrorContains(t, actualErr, tc.expectedErr)

				// The CLOB pair should not be able to be found in the store.
				require.False(t, found)
			}
		})
	}
}

func TestCreateMultipleClobPairs(t *testing.T) {
	type CreationExpectation struct {
		// CLOB pair.
		clobPair types.ClobPair

		// Expectations.
		expectedErr string
	}
	tests := map[string]struct {
		// The CLOB pairs to attempt to make.
		clobPairs []CreationExpectation

		// The expected number of created CLOB pairs.
		expectedNumClobPairs uint32

		// The expected mapping of ID -> CLOB pair.
		expectedStoredClobPairs map[types.ClobPairId]types.ClobPair
	}{
		"Successfully makes multiple CLOB pairs": {
			clobPairs: []CreationExpectation{
				{clobPair: constants.ClobPair_Btc},
				{clobPair: constants.ClobPair_Eth},
			},
			expectedNumClobPairs: 2,
			expectedStoredClobPairs: map[types.ClobPairId]types.ClobPair{
				0: constants.ClobPair_Btc,
				1: constants.ClobPair_Eth,
			},
		},
		"Can create a CLOB pair and then fail validation": {
			clobPairs: []CreationExpectation{
				{clobPair: constants.ClobPair_Btc},
				{
					clobPair: *clobtest.GenerateClobPair(
						clobtest.WithStatus(types.ClobPair_STATUS_UNSPECIFIED),
						clobtest.WithId(99999), // unused id
						clobtest.WithPerpetualId(1),
					),
					expectedErr: "has unsupported status STATUS_UNSPECIFIED",
				},
			},
			expectedNumClobPairs: 1,
			expectedStoredClobPairs: map[types.ClobPairId]types.ClobPair{
				0: constants.ClobPair_Btc,
			},
		},
		"Can create a CLOB pair after failing to create one": {
			clobPairs: []CreationExpectation{
				{
					clobPair:    *clobtest.GenerateClobPair(clobtest.WithStatus(types.ClobPair_STATUS_UNSPECIFIED)),
					expectedErr: "has unsupported status STATUS_UNSPECIFIED",
				},
				{clobPair: constants.ClobPair_Btc},
			},
			expectedNumClobPairs: 1,
			expectedStoredClobPairs: map[types.ClobPairId]types.ClobPair{
				0: constants.ClobPair_Btc,
			},
		},
		"Can alternate between passing/failing CLOB pair validation with no issues": {
			clobPairs: []CreationExpectation{
				{clobPair: constants.ClobPair_Btc},
				{
					clobPair: *clobtest.GenerateClobPair(
						clobtest.WithStatus(types.ClobPair_STATUS_UNSPECIFIED),
						clobtest.WithId(99999), // unused id
						clobtest.WithPerpetualId(1),
					),
					expectedErr: "has unsupported status STATUS_UNSPECIFIED",
				},
				{clobPair: constants.ClobPair_Eth},
			},
			expectedNumClobPairs: 2,
			expectedStoredClobPairs: map[types.ClobPairId]types.ClobPair{
				0: constants.ClobPair_Btc,
				1: constants.ClobPair_Eth,
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Boilerplate setup.
			memClob := memclob.NewMemClobPriceTimePriority(false)
			mockIndexerEventManager := &mocks.IndexerEventManager{}
			ks := keepertest.NewClobKeepersTestContext(t, memClob, &mocks.BankKeeper{}, mockIndexerEventManager)

			prices.InitGenesis(ks.Ctx, *ks.PricesKeeper, constants.Prices_DefaultGenesisState)
			perpetuals.InitGenesis(ks.Ctx, *ks.PerpetualsKeeper, constants.Perpetuals_DefaultGenesisState)

			// Perform the method under test.
			for _, make := range tc.clobPairs {
				if make.expectedErr == "" {
					perpetualId := clobtest.MustPerpetualId(make.clobPair)
					perpetual := constants.Perpetuals_DefaultGenesisState.Perpetuals[perpetualId]
					mockIndexerEventManager.On("AddTxnEvent",
						ks.Ctx,
						indexerevents.SubtypePerpetualMarket,
						indexer_manager.GetB64EncodedEventMessage(
							indexerevents.NewPerpetualMarketCreateEvent(
								perpetualId,
								perpetualId,
								perpetual.Params.Ticker,
								perpetual.Params.MarketId,
								make.clobPair.Status,
								make.clobPair.QuantumConversionExponent,
								perpetual.Params.AtomicResolution,
								make.clobPair.SubticksPerTick,
								make.clobPair.StepBaseQuantums,
								perpetual.Params.LiquidityTier,
							),
						),
					).Return()
				}

				_, err := ks.ClobKeeper.CreatePerpetualClobPair(
					ks.Ctx,
					make.clobPair.Id,
					clobtest.MustPerpetualId(make.clobPair),
					satypes.BaseQuantums(make.clobPair.StepBaseQuantums),
					make.clobPair.QuantumConversionExponent,
					make.clobPair.SubticksPerTick,
					make.clobPair.Status,
				)
				if make.expectedErr == "" {
					require.NoError(t, err)
				} else {
					require.Error(t, err)
					require.ErrorContains(t, err, make.expectedErr)
				}
			}

			for key, expectedClobPair := range tc.expectedStoredClobPairs {
				actual, found := ks.ClobKeeper.GetClobPair(ks.Ctx, key)
				require.True(t, found)
				require.Equal(t, expectedClobPair, actual)
			}

			_, found := ks.ClobKeeper.GetClobPair(ks.Ctx, types.ClobPairId(tc.expectedNumClobPairs))
			require.False(t, found)
		})
	}
}

func TestInitMemClobOrderbooks(t *testing.T) {
	memClob := memclob.NewMemClobPriceTimePriority(false)
	ks := keepertest.NewClobKeepersTestContext(
		t,
		memClob,
		&mocks.BankKeeper{},
		&mocks.IndexerEventManager{},
	)

	// Read a new `ClobPair` and make sure it does not exist.
	_, err := ks.ClobKeeper.GetClobPairIdForPerpetual(ks.Ctx, 1)
	require.ErrorIs(t, err, types.ErrNoClobPairForPerpetual)

	// Write multiple `ClobPairs` to state, but don't call `MemClob.CreateOrderbook`.
	store := prefix.NewStore(ks.Ctx.KVStore(ks.StoreKey), types.KeyPrefix(types.ClobPairKeyPrefix))
	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)

	b := cdc.MustMarshal(&constants.ClobPair_Eth)
	store.Set(types.ClobPairKey(
		types.ClobPairId(constants.ClobPair_Eth.Id),
	), b)

	b = cdc.MustMarshal(&constants.ClobPair_Btc)
	store.Set(types.ClobPairKey(
		types.ClobPairId(constants.ClobPair_Btc.Id),
	), b)

	// Read the new `ClobPairs` and make sure they do not exist.
	_, err = ks.ClobKeeper.GetClobPairIdForPerpetual(ks.Ctx, 1)
	require.ErrorIs(t, err, types.ErrNoClobPairForPerpetual)

	// Initialize the `ClobPairs` from Keeper state.
	ks.ClobKeeper.InitMemClobOrderbooks(ks.Ctx)
}

func TestHydrateClobPairAndPerpetualMapping(t *testing.T) {
	memClob := memclob.NewMemClobPriceTimePriority(false)
	ks := keepertest.NewClobKeepersTestContext(
		t,
		memClob,
		&mocks.BankKeeper{},
		&mocks.IndexerEventManager{},
	)

	// Read a new `ClobPair` and make sure it does not exist.
	_, err := ks.ClobKeeper.GetClobPairIdForPerpetual(ks.Ctx, 1)
	require.ErrorIs(t, err, types.ErrNoClobPairForPerpetual)

	// Write multiple `ClobPairs` to state, but don't call `MemClob.CreateOrderbook`.
	store := prefix.NewStore(ks.Ctx.KVStore(ks.StoreKey), types.KeyPrefix(types.ClobPairKeyPrefix))
	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)

	b := cdc.MustMarshal(&constants.ClobPair_Eth)
	store.Set(types.ClobPairKey(
		types.ClobPairId(constants.ClobPair_Eth.Id),
	), b)

	b = cdc.MustMarshal(&constants.ClobPair_Btc)
	store.Set(types.ClobPairKey(
		types.ClobPairId(constants.ClobPair_Btc.Id),
	), b)

	// Read the new `ClobPairs` and make sure they do not exist.
	_, err = ks.ClobKeeper.GetClobPairIdForPerpetual(ks.Ctx, 1)
	require.ErrorIs(t, err, types.ErrNoClobPairForPerpetual)

	// Initialize the `ClobPairs` from Keeper state.
	ks.ClobKeeper.HydrateClobPairAndPerpetualMapping(ks.Ctx)

	// Read the new `ClobPairs` and make sure they exist.
	_, err = ks.ClobKeeper.GetClobPairIdForPerpetual(ks.Ctx, 0)
	require.NoError(t, err)

	_, err = ks.ClobKeeper.GetClobPairIdForPerpetual(ks.Ctx, 1)
	require.NoError(t, err)
}

func TestClobPairGet(t *testing.T) {
	memClob := memclob.NewMemClobPriceTimePriority(false)
	mockIndexerEventManager := &mocks.IndexerEventManager{}
	ks := keepertest.NewClobKeepersTestContext(t, memClob, &mocks.BankKeeper{}, mockIndexerEventManager)
	items := keepertest.CreateNClobPair(t,
		ks.ClobKeeper,
		ks.PerpetualsKeeper,
		ks.PricesKeeper,
		ks.Ctx,
		10,
		mockIndexerEventManager,
	)
	for _, item := range items {
		rst, found := ks.ClobKeeper.GetClobPair(ks.Ctx,
			types.ClobPairId(item.Id),
		)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item), //nolint:staticcheck
			nullify.Fill(&rst),  //nolint:staticcheck
		)
	}
}
func TestClobPairRemove(t *testing.T) {
	memClob := memclob.NewMemClobPriceTimePriority(false)
	mockIndexerEventManager := &mocks.IndexerEventManager{}
	ks := keepertest.NewClobKeepersTestContext(t, memClob, &mocks.BankKeeper{}, mockIndexerEventManager)
	items := keepertest.CreateNClobPair(t,
		ks.ClobKeeper,
		ks.PerpetualsKeeper,
		ks.PricesKeeper,
		ks.Ctx,
		10,
		mockIndexerEventManager,
	)
	for _, item := range items {
		ks.ClobKeeper.RemoveClobPair(ks.Ctx,
			types.ClobPairId(item.Id),
		)
		_, found := ks.ClobKeeper.GetClobPair(ks.Ctx,
			types.ClobPairId(item.Id),
		)
		require.False(t, found)
	}
}

func TestClobPairGetAll(t *testing.T) {
	memClob := memclob.NewMemClobPriceTimePriority(false)
	mockIndexerEventManager := &mocks.IndexerEventManager{}
	ks := keepertest.NewClobKeepersTestContext(t, memClob, &mocks.BankKeeper{}, mockIndexerEventManager)
	items := keepertest.CreateNClobPair(t,
		ks.ClobKeeper,
		ks.PerpetualsKeeper,
		ks.PricesKeeper,
		ks.Ctx,
		10,
		mockIndexerEventManager,
	)
	require.ElementsMatch(t,
		nullify.Fill(items), //nolint:staticcheck
		nullify.Fill(ks.ClobKeeper.GetAllClobPairs(ks.Ctx)), //nolint:staticcheck
	)
}

func TestUpdateClobPair(t *testing.T) {
	testCases := map[string]struct {
		setup         func(t *testing.T, ks keepertest.ClobKeepersTestContext, manager *mocks.IndexerEventManager)
		status        types.ClobPair_Status
		expectedErr   string
		expectedPanic string
	}{
		"Succeeds with valid status transition": {
			setup: func(t *testing.T, ks keepertest.ClobKeepersTestContext, mockIndexerEventManager *mocks.IndexerEventManager) {
				clobPair := constants.ClobPair_Btc
				mockIndexerEventManager.On("AddTxnEvent",
					ks.Ctx,
					indexerevents.SubtypePerpetualMarket,
					indexer_manager.GetB64EncodedEventMessage(
						indexerevents.NewPerpetualMarketCreateEvent(
							0,
							0,
							constants.Perpetuals_DefaultGenesisState.Perpetuals[0].Params.Ticker,
							constants.Perpetuals_DefaultGenesisState.Perpetuals[0].Params.MarketId,
							types.ClobPair_STATUS_INITIALIZING,
							clobPair.QuantumConversionExponent,
							constants.Perpetuals_DefaultGenesisState.Perpetuals[0].Params.AtomicResolution,
							clobPair.SubticksPerTick,
							clobPair.StepBaseQuantums,
							constants.Perpetuals_DefaultGenesisState.Perpetuals[0].Params.LiquidityTier,
						),
					),
				).Once().Return()

				_, err := ks.ClobKeeper.CreatePerpetualClobPair(
					ks.Ctx,
					clobPair.Id,
					clobtest.MustPerpetualId(clobPair),
					satypes.BaseQuantums(clobPair.StepBaseQuantums),
					clobPair.QuantumConversionExponent,
					clobPair.SubticksPerTick,
					types.ClobPair_STATUS_INITIALIZING,
				)
				require.NoError(t, err)

				mockIndexerEventManager.On("AddTxnEvent",
					ks.Ctx,
					indexerevents.SubtypeUpdateClobPair,
					indexer_manager.GetB64EncodedEventMessage(
						indexerevents.NewUpdateClobPairEvent(
							clobPair.GetClobPairId(),
							types.ClobPair_STATUS_ACTIVE,
							clobPair.QuantumConversionExponent,
							types.SubticksPerTick(clobPair.GetSubticksPerTick()),
							satypes.BaseQuantums(clobPair.GetStepBaseQuantums()),
						),
					),
				).Once().Return()
			},
			status: types.ClobPair_STATUS_ACTIVE,
		},
		"Panics with missing clob pair": {
			setup: func(t *testing.T, ks keepertest.ClobKeepersTestContext, mockIndexerEventManager *mocks.IndexerEventManager) {
			},
			status:        types.ClobPair_STATUS_ACTIVE,
			expectedPanic: "mustGetClobPair: ClobPair with id 0 not found",
		},
		"Errors with unsupported transition to supported status": {
			setup: func(t *testing.T, ks keepertest.ClobKeepersTestContext, mockIndexerEventManager *mocks.IndexerEventManager) {
				clobPair := constants.ClobPair_Btc
				mockIndexerEventManager.On("AddTxnEvent",
					ks.Ctx,
					indexerevents.SubtypePerpetualMarket,
					indexer_manager.GetB64EncodedEventMessage(
						indexerevents.NewPerpetualMarketCreateEvent(
							0,
							0,
							constants.Perpetuals_DefaultGenesisState.Perpetuals[0].Params.Ticker,
							constants.Perpetuals_DefaultGenesisState.Perpetuals[0].Params.MarketId,
							clobPair.Status,
							clobPair.QuantumConversionExponent,
							constants.Perpetuals_DefaultGenesisState.Perpetuals[0].Params.AtomicResolution,
							clobPair.SubticksPerTick,
							clobPair.StepBaseQuantums,
							constants.Perpetuals_DefaultGenesisState.Perpetuals[0].Params.LiquidityTier,
						),
					),
				).Once().Return()

				_, err := ks.ClobKeeper.CreatePerpetualClobPair(
					ks.Ctx,
					clobPair.Id,
					clobtest.MustPerpetualId(clobPair),
					satypes.BaseQuantums(clobPair.StepBaseQuantums),
					clobPair.QuantumConversionExponent,
					clobPair.SubticksPerTick,
					clobPair.Status,
				)
				require.NoError(t, err)
			},
			status:      types.ClobPair_STATUS_INITIALIZING,
			expectedErr: "Cannot transition from status STATUS_ACTIVE to status STATUS_INITIALIZING",
		},
		"Errors with unsupported transition to unsupported status": {
			setup: func(t *testing.T, ks keepertest.ClobKeepersTestContext, mockIndexerEventManager *mocks.IndexerEventManager) {
				clobPair := constants.ClobPair_Btc
				mockIndexerEventManager.On("AddTxnEvent",
					ks.Ctx,
					indexerevents.SubtypePerpetualMarket,
					indexer_manager.GetB64EncodedEventMessage(
						indexerevents.NewPerpetualMarketCreateEvent(
							0,
							0,
							constants.Perpetuals_DefaultGenesisState.Perpetuals[0].Params.Ticker,
							constants.Perpetuals_DefaultGenesisState.Perpetuals[0].Params.MarketId,
							clobPair.Status,
							clobPair.QuantumConversionExponent,
							constants.Perpetuals_DefaultGenesisState.Perpetuals[0].Params.AtomicResolution,
							clobPair.SubticksPerTick,
							clobPair.StepBaseQuantums,
							constants.Perpetuals_DefaultGenesisState.Perpetuals[0].Params.LiquidityTier,
						),
					),
				).Once().Return()

				_, err := ks.ClobKeeper.CreatePerpetualClobPair(
					ks.Ctx,
					clobPair.Id,
					clobtest.MustPerpetualId(clobPair),
					satypes.BaseQuantums(clobPair.StepBaseQuantums),
					clobPair.QuantumConversionExponent,
					clobPair.SubticksPerTick,
					clobPair.Status,
				)
				require.NoError(t, err)
			},
			status:      types.ClobPair_Status(100),
			expectedErr: "Cannot transition from status STATUS_ACTIVE to status 100",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			memClob := memclob.NewMemClobPriceTimePriority(false)
			mockIndexerEventManager := &mocks.IndexerEventManager{}
			ks := keepertest.NewClobKeepersTestContext(t, memClob, &mocks.BankKeeper{}, mockIndexerEventManager)
			prices.InitGenesis(ks.Ctx, *ks.PricesKeeper, constants.Prices_DefaultGenesisState)
			perpetuals.InitGenesis(ks.Ctx, *ks.PerpetualsKeeper, constants.Perpetuals_DefaultGenesisState)

			tc.setup(t, ks, mockIndexerEventManager)
			clobPair := constants.ClobPair_Btc
			clobPair.Status = tc.status
			if tc.expectedPanic != "" {
				require.PanicsWithValue(
					t,
					tc.expectedPanic,
					func() {
						err := ks.ClobKeeper.UpdateClobPair(ks.Ctx, clobPair)
						require.NoError(t, err)
					},
				)
			} else {
				err := ks.ClobKeeper.UpdateClobPair(ks.Ctx, clobPair)
				mockIndexerEventManager.AssertExpectations(t)

				if tc.expectedErr != "" {
					require.ErrorContains(t, err, tc.expectedErr)
				} else {
					require.NoError(t, err)
				}
			}
		})
	}
}

func TestGetClobPairIdForPerpetual_Success(t *testing.T) {
	memClob := memclob.NewMemClobPriceTimePriority(false)
	ks := keepertest.NewClobKeepersTestContext(t, memClob, &mocks.BankKeeper{}, &mocks.IndexerEventManager{})

	ks.ClobKeeper.PerpetualIdToClobPairId = map[uint32][]types.ClobPairId{
		0: {types.ClobPairId(0)},
	}

	clobPairId, err := ks.ClobKeeper.GetClobPairIdForPerpetual(ks.Ctx, 0)
	require.NoError(t, err)
	require.Equal(t, types.ClobPairId(0), clobPairId)
}

func TestGetAllClobPairs_Sorted(t *testing.T) {
	memClob := memclob.NewMemClobPriceTimePriority(false)
	mockIndexerEventManager := &mocks.IndexerEventManager{}
	ks := keepertest.NewClobKeepersTestContext(t, memClob, &mocks.BankKeeper{}, mockIndexerEventManager)
	keepertest.CreateLiquidityTiersAndNPerpetuals(t,
		ks.Ctx,
		ks.PerpetualsKeeper,
		ks.PricesKeeper,
		6,
	)

	clobPairs := []types.ClobPair{
		*clobtest.GenerateClobPair(clobtest.WithId(0), clobtest.WithPerpetualId(0)),
		*clobtest.GenerateClobPair(clobtest.WithId(5), clobtest.WithPerpetualId(1)),
		*clobtest.GenerateClobPair(clobtest.WithId(999), clobtest.WithPerpetualId(2)),
		*clobtest.GenerateClobPair(clobtest.WithId(900), clobtest.WithPerpetualId(3)),
		*clobtest.GenerateClobPair(clobtest.WithId(10), clobtest.WithPerpetualId(4)),
		*clobtest.GenerateClobPair(clobtest.WithId(1), clobtest.WithPerpetualId(5)),
	}

	mockIndexerEventManager.On("AddTxnEvent",
		ks.Ctx,
		indexerevents.SubtypePerpetualMarket,
		mock.Anything,
	).Return().Times(len(clobPairs))

	for _, clobPair := range clobPairs {
		_, err := ks.ClobKeeper.CreatePerpetualClobPair(
			ks.Ctx,
			clobPair.Id,
			clobtest.MustPerpetualId(clobPair),
			satypes.BaseQuantums(clobPair.StepBaseQuantums),
			clobPair.QuantumConversionExponent,
			clobPair.SubticksPerTick,
			clobPair.Status,
		)
		require.NoError(t, err)
	}

	expected := []types.ClobPair{
		*clobtest.GenerateClobPair(clobtest.WithId(0), clobtest.WithPerpetualId(0)),
		*clobtest.GenerateClobPair(clobtest.WithId(1), clobtest.WithPerpetualId(5)),
		*clobtest.GenerateClobPair(clobtest.WithId(5), clobtest.WithPerpetualId(1)),
		*clobtest.GenerateClobPair(clobtest.WithId(10), clobtest.WithPerpetualId(4)),
		*clobtest.GenerateClobPair(clobtest.WithId(900), clobtest.WithPerpetualId(3)),
		*clobtest.GenerateClobPair(clobtest.WithId(999), clobtest.WithPerpetualId(2)),
	}
	got := ks.ClobKeeper.GetAllClobPairs(ks.Ctx)
	require.Equal(t, expected, got)
}

func TestGetClobPairIdForPerpetual_ErrorNoClobPair(t *testing.T) {
	memClob := memclob.NewMemClobPriceTimePriority(false)
	ks := keepertest.NewClobKeepersTestContext(t, memClob, &mocks.BankKeeper{}, &mocks.IndexerEventManager{})

	_, err := ks.ClobKeeper.GetClobPairIdForPerpetual(ks.Ctx, 0)
	require.EqualError(
		t,
		err,
		"Perpetual ID 0 has no associated CLOB pairs: The provided perpetual ID "+
			"does not have any associated CLOB pairs",
	)
}

func TestGetClobPairIdForPerpetual_PanicsEmptyClobPair(t *testing.T) {
	memClob := memclob.NewMemClobPriceTimePriority(false)
	ks := keepertest.NewClobKeepersTestContext(t, memClob, &mocks.BankKeeper{}, &mocks.IndexerEventManager{})

	ks.ClobKeeper.PerpetualIdToClobPairId = map[uint32][]types.ClobPairId{
		0: {},
	}

	require.PanicsWithValue(
		t,
		"GetClobPairIdForPerpetual: Perpetual ID was created without a CLOB pair ID.",
		func() {
			if _, err := ks.ClobKeeper.GetClobPairIdForPerpetual(ks.Ctx, 0); err != nil {
				fmt.Printf("function should panic, not have error %+v", err)
			}
		},
	)
}

func TestGetClobPairIdForPerpetual_PanicsMultipleClobPairIds(t *testing.T) {
	memClob := memclob.NewMemClobPriceTimePriority(false)
	ks := keepertest.NewClobKeepersTestContext(t, memClob, &mocks.BankKeeper{}, &mocks.IndexerEventManager{})

	ks.ClobKeeper.PerpetualIdToClobPairId = map[uint32][]types.ClobPairId{
		0: {types.ClobPairId(0), types.ClobPairId(1)},
	}

	require.PanicsWithValue(
		t,
		"GetClobPairIdForPerpetual: Perpetual ID was created with multiple CLOB pair IDs.",
		func() {
			if _, err := ks.ClobKeeper.GetClobPairIdForPerpetual(ks.Ctx, 0); err != nil {
				fmt.Printf("function should panic, not have error %+v", err)
			}
		},
	)
}

func TestIsPerpetualClobPairActive(t *testing.T) {
	testCases := map[string]struct {
		clobPair                *types.ClobPair
		perpetualIdToClobPairId map[uint32][]types.ClobPairId
		resp                    bool
		expectedErr             error
	}{
		"Errors when perpetual has no clob pairs": {
			expectedErr: types.ErrNoClobPairForPerpetual,
		},
		"Errors when clob pair does not exist": {
			perpetualIdToClobPairId: map[uint32][]types.ClobPairId{
				0: {types.ClobPairId(0)},
			},
			expectedErr: types.ErrInvalidClob,
		},
		"Succeeds when clob pair is initializing": {
			perpetualIdToClobPairId: map[uint32][]types.ClobPairId{
				0: {types.ClobPairId(0)},
			},
			clobPair: &constants.ClobPair_Btc_Init,
			resp:     false,
		},
		"Succeeds when clob pair is active": {
			perpetualIdToClobPairId: map[uint32][]types.ClobPairId{
				0: {types.ClobPairId(0)},
			},
			clobPair: &constants.ClobPair_Btc,
			resp:     true,
		},
		"Succeeds when clob pair is paused": {
			perpetualIdToClobPairId: map[uint32][]types.ClobPairId{
				0: {types.ClobPairId(0)},
			},
			clobPair: &constants.ClobPair_Btc_Paused,
			resp:     false,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			memClob := memclob.NewMemClobPriceTimePriority(false)
			mockIndexerEventManager := &mocks.IndexerEventManager{}
			ks := keepertest.NewClobKeepersTestContext(t, memClob, &mocks.BankKeeper{}, mockIndexerEventManager)
			prices.InitGenesis(ks.Ctx, *ks.PricesKeeper, constants.Prices_DefaultGenesisState)
			perpetuals.InitGenesis(ks.Ctx, *ks.PerpetualsKeeper, constants.Perpetuals_DefaultGenesisState)

			if tc.clobPair != nil {
				// allows us to circumvent CreatePerpetualClobPair and write unsupported statuses to state to
				// test this function with unsupported statuses.
				registry := codectypes.NewInterfaceRegistry()
				cdc := codec.NewProtoCodec(registry)
				store := prefix.NewStore(ks.Ctx.KVStore(ks.StoreKey), types.KeyPrefix(types.ClobPairKeyPrefix))

				b := cdc.MustMarshal(tc.clobPair)
				store.Set(types.ClobPairKey(
					types.ClobPairId(tc.clobPair.Id),
				), b)
			}

			ks.ClobKeeper.PerpetualIdToClobPairId = tc.perpetualIdToClobPairId

			resp, err := ks.ClobKeeper.IsPerpetualClobPairActive(ks.Ctx, 0)

			if tc.expectedErr != nil {
				require.ErrorIs(t, err, tc.expectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.resp, resp)
			}
		})
	}
}

func TestClobPairValidate(t *testing.T) {
	tests := []struct {
		desc        string
		clobPair    types.ClobPair
		expectedErr string
	}{
		{
			desc: "Invalid Metadata (SpotClobMetadata)",
			clobPair: types.ClobPair{
				Metadata:         &types.ClobPair_SpotClobMetadata{},
				StepBaseQuantums: 1,
				SubticksPerTick:  1,
				Status:           types.ClobPair_STATUS_ACTIVE,
			},
			expectedErr: "is not a perpetual CLOB",
		},
		{
			desc: "Unsupported Status",
			clobPair: types.ClobPair{
				Metadata:         &types.ClobPair_PerpetualClobMetadata{},
				StepBaseQuantums: 1,
				SubticksPerTick:  1,
				Status:           types.ClobPair_STATUS_PAUSED,
			},
			expectedErr: "has unsupported status",
		},
		{
			desc: "StepBaseQuantums <= 0",
			clobPair: types.ClobPair{
				Metadata:         &types.ClobPair_PerpetualClobMetadata{},
				StepBaseQuantums: 0,
				SubticksPerTick:  1,
				Status:           types.ClobPair_STATUS_ACTIVE,
			},
			expectedErr: "StepBaseQuantums must be > 0.",
		},
		{
			desc: "SubticksPerTick <= 0",
			clobPair: types.ClobPair{
				Metadata:         &types.ClobPair_PerpetualClobMetadata{},
				StepBaseQuantums: 1,
				SubticksPerTick:  0,
				Status:           types.ClobPair_STATUS_ACTIVE,
			},
			expectedErr: "SubticksPerTick must be > 0",
		},
		{
			desc: "Valid ClobPair",
			clobPair: types.ClobPair{
				Metadata:         &types.ClobPair_PerpetualClobMetadata{},
				StepBaseQuantums: 1,
				SubticksPerTick:  1,
				Status:           types.ClobPair_STATUS_ACTIVE,
			},
			expectedErr: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.clobPair.Validate()
			if tc.expectedErr == "" {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, tc.expectedErr)
			}
		})
	}
}
