package keeper_test

import (
	"sort"
	"strconv"
	"testing"
	"time"

	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/nullify"
	"github.com/dydxprotocol/v4-chain/protocol/x/epochs/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/epochs/types"
	"github.com/stretchr/testify/require"
)

// Prevent strconv unused error
var _ = strconv.IntSize

// creates n EpochInfo with TestBlockTime and TestEpochDuration
func createNEpochInfo(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.EpochInfo {
	blockTime := time.Unix(keepertest.TestCreateEpochBlockTimeSec, 0)
	mockCtx := ctx.WithBlockTime(blockTime)

	items := make([]types.EpochInfo, n)
	for i := range items {
		items[i].Name = strconv.Itoa(i)
		items[i].Duration = keepertest.TestEpochDuration
		items[i].NextTick = keepertest.TestCreateEpochBlockTimeSec
		if err := keeper.CreateEpochInfo(mockCtx, items[i]); err != nil {
			panic(err)
		}
	}
	return items
}

func TestMaybeStartNextEpoch(t *testing.T) {
	tests := map[string]struct {
		epochToCreate     types.EpochInfo
		laterBlockTimeSec int64
		laterBlockHeight  int64
		// Expected new `EpochInfo` object (newly initialized and/or newly ticked)
		expectedEpoch *types.EpochInfo
		// Whether a new epoch has ticked
		expectNextEpochTicked bool
		expectedErr           error
		expectedEvents        []sdk.Event
	}{
		"already initialized, next epoch didn't tick": {
			epochToCreate: types.EpochInfo{
				Name:          keepertest.TestEpochInfoName,
				Duration:      60,
				NextTick:      1800000060,
				IsInitialized: true,
			},
			laterBlockTimeSec: 1800000059,
		},
		"already initialized, next epoch ticked": {
			epochToCreate: types.EpochInfo{
				Name:                   keepertest.TestEpochInfoName,
				Duration:               60,
				CurrentEpoch:           56,
				CurrentEpochStartBlock: 1,
				NextTick:               1800000060,
				IsInitialized:          true,
			},
			laterBlockTimeSec:     1800000060,
			laterBlockHeight:      1234,
			expectNextEpochTicked: true,
			expectedEpoch: &types.EpochInfo{
				Name:                   keepertest.TestEpochInfoName,
				Duration:               60,
				CurrentEpoch:           57,
				NextTick:               1800000120,
				CurrentEpochStartBlock: 1234,
				IsInitialized:          true,
			},
			expectedEvents: []sdk.Event{
				sdk.NewEvent(
					types.EventTypeNewEpoch,
					sdk.NewAttribute(types.AttributeKeyEpochInfoName, keepertest.TestEpochInfoName),
					sdk.NewAttribute(types.AttributeKeyEpochNumber, "57"),
					sdk.NewAttribute(types.AttributeKeyEpochStartTickTime, "1800000060"),
					sdk.NewAttribute(types.AttributeKeyEpochStartBlockTime, "1800000060"),
					sdk.NewAttribute(types.AttributeKeyEpochStartBlock, "1234"),
				),
			},
		},
		"already initialized, downtime recovery: next epoch reached, only catch up epoch per block": {
			epochToCreate: types.EpochInfo{
				Name:                   keepertest.TestEpochInfoName,
				Duration:               60,
				CurrentEpoch:           56,
				CurrentEpochStartBlock: 1,
				NextTick:               1800000060,
				IsInitialized:          true,
			},
			laterBlockTimeSec:     1800006660,
			laterBlockHeight:      1234,
			expectNextEpochTicked: true,
			expectedEpoch: &types.EpochInfo{
				Name:                   keepertest.TestEpochInfoName,
				Duration:               60,
				CurrentEpoch:           57,
				NextTick:               1800000120,
				CurrentEpochStartBlock: 1234,
				IsInitialized:          true,
			},
			expectedEvents: []sdk.Event{
				sdk.NewEvent(
					types.EventTypeNewEpoch,
					sdk.NewAttribute(types.AttributeKeyEpochInfoName, keepertest.TestEpochInfoName),
					sdk.NewAttribute(types.AttributeKeyEpochNumber, "57"),
					sdk.NewAttribute(types.AttributeKeyEpochStartTickTime, "1800000060"),
					sdk.NewAttribute(types.AttributeKeyEpochStartBlockTime, "1800006660"),
					sdk.NewAttribute(types.AttributeKeyEpochStartBlock, "1234"),
				),
			},
		},
		"not initialized, don't initialize epoch": {
			epochToCreate: types.EpochInfo{
				Name:                keepertest.TestEpochInfoName,
				Duration:            60,
				NextTick:            1900000000,
				IsInitialized:       false,
				FastForwardNextTick: true,
			},
			laterBlockTimeSec: 1800006660, // < `NextTick`
			laterBlockHeight:  1234,
		},
		"initialize epoch, don't tick": {
			epochToCreate: types.EpochInfo{
				Name:                keepertest.TestEpochInfoName,
				Duration:            60,
				NextTick:            1800000000,
				IsInitialized:       false,
				FastForwardNextTick: true,
			},
			laterBlockTimeSec: 1800000001, // > `NextTick`
			laterBlockHeight:  1234,
			expectedEpoch: &types.EpochInfo{
				Name:                keepertest.TestEpochInfoName,
				Duration:            60,
				IsInitialized:       true,
				FastForwardNextTick: true,
				NextTick:            1800000060,
			},
		},
		"initialize epoch, `NextTick` same as block time, don't tick": {
			epochToCreate: types.EpochInfo{
				Name:                keepertest.TestEpochInfoName,
				Duration:            60,
				NextTick:            1500000000,
				IsInitialized:       false,
				FastForwardNextTick: true,
			},
			laterBlockTimeSec: 1500000000, // = `NextTick`
			laterBlockHeight:  1234,
			expectedEpoch: &types.EpochInfo{
				Name:                keepertest.TestEpochInfoName,
				Duration:            60,
				IsInitialized:       true,
				FastForwardNextTick: true,
				NextTick:            1500000060, // fast-forwarded
			},
		},
		"initialize epoch and tick first epoch": {
			epochToCreate: types.EpochInfo{
				Name:          keepertest.TestEpochInfoName,
				Duration:      60,
				NextTick:      1800000000,
				IsInitialized: false,
				// must be false for the first epoch to tick at initialization
				FastForwardNextTick: false,
			},
			laterBlockTimeSec: 1800000001,
			laterBlockHeight:  1234,
			expectedEpoch: &types.EpochInfo{
				Name:                   keepertest.TestEpochInfoName,
				Duration:               60,
				CurrentEpoch:           1,
				NextTick:               1800000060,
				CurrentEpochStartBlock: 1234,
				IsInitialized:          true,
			},
			expectNextEpochTicked: true,
			expectedEvents: []sdk.Event{
				sdk.NewEvent(
					types.EventTypeNewEpoch,
					sdk.NewAttribute(types.AttributeKeyEpochInfoName, keepertest.TestEpochInfoName),
					sdk.NewAttribute(types.AttributeKeyEpochNumber, "1"),
					sdk.NewAttribute(types.AttributeKeyEpochStartTickTime, "1800000000"),
					sdk.NewAttribute(types.AttributeKeyEpochStartBlockTime, "1800000001"),
					sdk.NewAttribute(types.AttributeKeyEpochStartBlock, "1234"),
				),
			},
		},
		"error: epoch doesn't exist": {
			epochToCreate: types.EpochInfo{
				Name:     "different_name",
				Duration: 60,
				NextTick: 1800000060,
			},
			laterBlockTimeSec: 1800000059,
			expectedErr:       types.ErrEpochInfoNotFound,
		},
	}

	// Run tests.
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctx, keeper, _ := keepertest.EpochsKeeper(t)
			err := keeper.CreateEpochInfo(ctx, tc.epochToCreate)
			require.NoError(t, err)

			laterBlockTime := time.Unix(tc.laterBlockTimeSec, 0)
			laterCtx := ctx.WithBlockTime(laterBlockTime).WithBlockHeight(tc.laterBlockHeight)

			nextEpochStarted, err := keeper.MaybeStartNextEpoch(laterCtx, keepertest.TestEpochInfoName)
			if tc.expectedErr != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tc.expectedErr)
				return
			}

			require.NoError(t, err)
			require.Equal(t,
				tc.expectNextEpochTicked,
				nextEpochStarted,
			)
			epochInfo, found := keeper.GetEpochInfo(ctx, keepertest.TestEpochInfoName)
			require.True(t, found)

			if tc.expectedEpoch != nil {
				require.Equal(t,
					*tc.expectedEpoch,
					epochInfo,
				)
			} else {
				// New epoch not started, check epoch is unchanged.
				require.Equal(t,
					tc.epochToCreate,
					epochInfo,
				)
			}

			if tc.expectNextEpochTicked {
				require.ElementsMatch(t,
					sdk.Events(tc.expectedEvents).ToABCIEvents(),
					laterCtx.EventManager().ABCIEvents(),
				)
			} else {
				require.Equal(t,
					0,
					len(laterCtx.EventManager().ABCIEvents()),
				)
			}
		})
	}
}

func TestCreateEpochInfo(t *testing.T) {
	tests := map[string]struct {
		epochInfoToCreate   types.EpochInfo
		fastForwardNextTick bool
		blockTimeSec        int64
		blockHeight         int64
		expectedEpochInfo   *types.EpochInfo
		expectedErr         error
	}{
		"success - use given currentEpoch and currentEpochStartBlock": {
			epochInfoToCreate: types.EpochInfo{
				Name:                   keepertest.TestEpochInfoName,
				Duration:               60,
				CurrentEpoch:           123,
				CurrentEpochStartBlock: 4567,
			},
			blockTimeSec: 1800000000,
			expectedEpochInfo: &types.EpochInfo{
				Name:                   keepertest.TestEpochInfoName,
				Duration:               60,
				CurrentEpoch:           123,
				CurrentEpochStartBlock: 4567,
			},
		},
		"success": {
			epochInfoToCreate: types.EpochInfo{
				Name:     keepertest.TestEpochInfoName,
				NextTick: 1200,
				Duration: 60,
			},
			blockTimeSec: 1800000001,
			expectedEpochInfo: &types.EpochInfo{
				Name:     keepertest.TestEpochInfoName,
				NextTick: 1200,
				Duration: 60,
			},
		},
		"error - fails validation": {
			epochInfoToCreate: types.EpochInfo{
				Name: keepertest.TestEpochInfoName,
			},
			expectedErr: types.ErrDurationIsZero,
		},
	}

	// Run tests.
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctx, keeper, _ := keepertest.EpochsKeeper(t)
			blockTime := time.Unix(tc.blockTimeSec, 0)
			mockCtx := ctx.WithBlockTime(blockTime).WithBlockHeight(tc.blockHeight)
			err := keeper.CreateEpochInfo(mockCtx, tc.epochInfoToCreate)
			if tc.expectedErr != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tc.expectedErr)
				return
			}

			require.NoError(t, err)
			epochInfo, found := keeper.GetEpochInfo(ctx, keepertest.TestEpochInfoName)
			require.True(t, found)

			require.Equal(t,
				*tc.expectedEpochInfo,
				epochInfo,
			)
		})
	}
}

func TestEpochInfoGet(t *testing.T) {
	ctx, keeper, _ := keepertest.EpochsKeeper(t)
	items := createNEpochInfo(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetEpochInfo(ctx,
			item.GetEpochInfoName(),
		)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item), //nolint:staticcheck
			nullify.Fill(&rst),  //nolint:staticcheck
		)
	}
}

func TestMustGetFundingEpochInfo(t *testing.T) {
	tests := map[string]struct {
		epochInfoName        types.EpochInfoName
		mustGetEpochInfoFunc func(k *keeper.Keeper, ctx sdk.Context) types.EpochInfo
	}{
		"Must get funding-tick": {
			epochInfoName: types.FundingTickEpochInfoName,
			mustGetEpochInfoFunc: func(k *keeper.Keeper, ctx sdk.Context) types.EpochInfo {
				return k.MustGetFundingTickEpochInfo(ctx)
			},
		},
		"Must get funding-sample": {
			epochInfoName: types.FundingSampleEpochInfoName,
			mustGetEpochInfoFunc: func(k *keeper.Keeper, ctx sdk.Context) types.EpochInfo {
				return k.MustGetFundingSampleEpochInfo(ctx)
			},
		},
	}

	testDuration := uint32(10)

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctx, keeper, _ := keepertest.EpochsKeeper(t)

			// No epoch info created, should panic
			require.PanicsWithError(t, errorsmod.Wrapf(
				types.ErrEpochInfoNotFound,
				"name: %s",
				tc.epochInfoName,
			).Error(), func() {
				//nolint:errcheck
				tc.mustGetEpochInfoFunc(keeper, ctx)
			})

			// Create epoch info
			err := keeper.CreateEpochInfo(ctx, types.EpochInfo{
				Name:     string(tc.epochInfoName),
				Duration: testDuration,
			})
			require.NoError(t, err)

			epochInfo := tc.mustGetEpochInfoFunc(keeper, ctx)
			require.Equal(t,
				types.EpochInfo{
					Name:     string(tc.epochInfoName),
					Duration: testDuration,
				},
				epochInfo,
			)
		})
	}
}

func TestEpochInfoGetAll(t *testing.T) {
	ctx, keeper, _ := keepertest.EpochsKeeper(t)
	items := createNEpochInfo(keeper, ctx, 10)

	got := keeper.GetAllEpochInfo(ctx)
	require.ElementsMatch(t,
		nullify.Fill(items), //nolint:staticcheck
		nullify.Fill(got),   //nolint:staticcheck
	)
	require.True(t,
		sort.SliceIsSorted(got, func(p, q int) bool {
			return got[p].Name < got[q].Name
		}))
}

func TestNumBlocksSinceEpochStart(t *testing.T) {
	tests := map[string]struct {
		epochName       types.EpochInfoName
		epochStartBlock uint32
		currBlockHeight int64
		currentEpoch    uint32
		numBlocks       uint32
		expectedErr     error
	}{
		"success": {
			epochName:       keepertest.TestEpochInfoName,
			epochStartBlock: uint32(100),
			currentEpoch:    1,
			currBlockHeight: 123,
			numBlocks:       23,
		},
		"success with epoch processed": {
			epochName:       keepertest.TestEpochInfoName,
			epochStartBlock: uint32(0),
			currentEpoch:    0,
			currBlockHeight: 0,
			numBlocks:       0,
		},
		"success with same block height": {
			epochName:       keepertest.TestEpochInfoName,
			epochStartBlock: uint32(100),
			currBlockHeight: 100,
			currentEpoch:    1,
			numBlocks:       0,
		},
		"success with createBlockHeight = 0": {
			epochName:       keepertest.TestEpochInfoName,
			epochStartBlock: uint32(0),
			currBlockHeight: 23,
			numBlocks:       23,
		},
		"error - get non-existing epoch info Name": {
			epochName:       "11",
			epochStartBlock: uint32(0),
			currBlockHeight: 0,
			numBlocks:       0,
			expectedErr:     types.ErrEpochInfoNotFound,
		},
	}

	// Run tests.
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			epochInfo := types.EpochInfo{
				Name:                   keepertest.TestEpochInfoName,
				Duration:               keepertest.TestEpochDuration,
				CurrentEpoch:           tc.currentEpoch,
				CurrentEpochStartBlock: tc.epochStartBlock,
			}

			ctx, keeper, _ := keepertest.EpochsKeeper(t)
			require.NoError(t, keeper.CreateEpochInfo(ctx, epochInfo))

			numCtx := ctx.WithBlockHeight(tc.currBlockHeight)
			numBlocks, err := keeper.NumBlocksSinceEpochStart(numCtx, types.EpochInfoName(tc.epochName))
			if tc.expectedErr != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tc.expectedErr)
				return
			}

			require.NoError(t, err)
			require.Equal(t,
				uint32(tc.numBlocks),
				numBlocks,
			)
		})
	}
}
