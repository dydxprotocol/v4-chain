package keeper_test

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/mocks"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/constants"
	keepertest "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/keeper"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/clob/keeper"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/clob/memclob"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/clob/types"
	perptypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/perpetuals/types"
	satypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/subaccounts/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestSetGetSubaccountLiquidationInfo(t *testing.T) {
	tests := map[string]struct {
		// Setup.
		setupState func(ctx sdk.Context, k *keeper.Keeper)

		// Parameters.
		subaccountId satypes.SubaccountId

		// Expectations.
		expectedSubaccountLiquidationInfo types.SubaccountLiquidationInfo
	}{
		"gets default subaccount liquidation info if it hasn't been created yet": {
			subaccountId: constants.Alice_Num0,

			expectedSubaccountLiquidationInfo: types.SubaccountLiquidationInfo{},
		},
		"can get and set subaccount liquidation info": {
			setupState: func(ctx sdk.Context, k *keeper.Keeper) {
				k.MustUpdateSubaccountPerpetualLiquidated(
					ctx,
					constants.Alice_Num0,
					2,
				)
			},

			subaccountId: constants.Alice_Num0,

			expectedSubaccountLiquidationInfo: types.SubaccountLiquidationInfo{
				PerpetualsLiquidated: []uint32{2},
			},
		},
		"can get and set subaccount liquidation info multiple times for one subaccount": {
			setupState: func(ctx sdk.Context, k *keeper.Keeper) {
				k.MustUpdateSubaccountPerpetualLiquidated(
					ctx,
					constants.Alice_Num0,
					2,
				)
				k.MustUpdateSubaccountPerpetualLiquidated(
					ctx,
					constants.Alice_Num0,
					3,
				)
				k.MustUpdateSubaccountPerpetualLiquidated(
					ctx,
					constants.Alice_Num0,
					100,
				)
			},

			subaccountId: constants.Alice_Num0,

			expectedSubaccountLiquidationInfo: types.SubaccountLiquidationInfo{
				PerpetualsLiquidated: []uint32{2, 3, 100},
			},
		},
		"can get and set subaccount liquidation info for multiple subaccounts": {
			setupState: func(ctx sdk.Context, k *keeper.Keeper) {
				k.MustUpdateSubaccountPerpetualLiquidated(
					ctx,
					constants.Alice_Num0,
					2,
				)

				k.MustUpdateSubaccountPerpetualLiquidated(
					ctx,
					constants.Alice_Num1,
					3,
				)
				k.MustUpdateSubaccountPerpetualLiquidated(
					ctx,
					constants.Bob_Num0,
					100,
				)
			},

			subaccountId: constants.Bob_Num0,

			expectedSubaccountLiquidationInfo: types.SubaccountLiquidationInfo{
				PerpetualsLiquidated: []uint32{100},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup keeper state.
			memClob := memclob.NewMemClobPriceTimePriority(false)
			bankMock := &mocks.BankKeeper{}
			ks := keepertest.NewClobKeepersTestContext(t, memClob, bankMock, &mocks.IndexerEventManager{})

			if tc.setupState != nil {
				tc.setupState(ks.Ctx, ks.ClobKeeper)
			}

			// Run the test and verify expectations.
			subaccountLiquidationInfo := ks.ClobKeeper.GetSubaccountLiquidationInfo(
				ks.Ctx,
				tc.subaccountId,
			)
			require.Equal(
				t,
				tc.expectedSubaccountLiquidationInfo,
				subaccountLiquidationInfo,
			)
		})
	}
}

func TestUpdateSubaccountLiquidationInfo_MultipleLiquidationsOfSubaccountAndPerpetualPanics(t *testing.T) {
	memClob := memclob.NewMemClobPriceTimePriority(false)
	bankMock := &mocks.BankKeeper{}
	ks := keepertest.NewClobKeepersTestContext(t, memClob, bankMock, &mocks.IndexerEventManager{})

	subaccountId := constants.Alice_Num0
	perpetualId := uint32(0)
	ks.ClobKeeper.MustUpdateSubaccountPerpetualLiquidated(
		ks.Ctx,
		subaccountId,
		perpetualId,
	)

	require.PanicsWithValue(
		t,
		fmt.Sprintf(
			"UpdateSubaccountLiquidationInfo: subaccount %v and perpetual %d have already been "+
				"liquidated within the last block",
			subaccountId,
			perpetualId,
		),
		func() {
			ks.ClobKeeper.MustUpdateSubaccountPerpetualLiquidated(
				ks.Ctx,
				subaccountId,
				perpetualId,
			)
		},
	)
}

func TestIncrementCumulativeInsuranceFundDelta(t *testing.T) {
	tests := map[string]struct {
		perpetuals      []perptypes.Perpetual
		initialDeltas   map[uint32]*big.Int
		incrementDeltas map[uint32]*big.Int
		expectedDeltas  map[uint32]*big.Int
		expectError     bool
	}{
		"increment from zero": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_SmallMarginRequirement,
			},
			initialDeltas: map[uint32]*big.Int{
				0: big.NewInt(0),
			},
			incrementDeltas: map[uint32]*big.Int{
				0: big.NewInt(1000),
			},
			expectedDeltas: map[uint32]*big.Int{
				0: big.NewInt(-1000),
			},
			expectError: false,
		},
		"increment positive delta": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_SmallMarginRequirement,
			},
			initialDeltas: map[uint32]*big.Int{
				0: big.NewInt(-500),
			},
			incrementDeltas: map[uint32]*big.Int{
				0: big.NewInt(500),
			},
			expectedDeltas: map[uint32]*big.Int{
				0: big.NewInt(0),
			},
			expectError: false,
		},
		"increment negative delta": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_SmallMarginRequirement,
			},
			initialDeltas: map[uint32]*big.Int{
				0: big.NewInt(-1000),
			},
			incrementDeltas: map[uint32]*big.Int{
				0: big.NewInt(-500),
			},
			expectedDeltas: map[uint32]*big.Int{
				0: big.NewInt(1500),
			},
			expectError: false,
		},
		"perpetual not set up": {
			perpetuals: []perptypes.Perpetual{},
			initialDeltas: map[uint32]*big.Int{
				0: big.NewInt(0),
			},
			incrementDeltas: map[uint32]*big.Int{
				0: big.NewInt(1000),
			},
			expectedDeltas: nil,
			expectError:    true,
		},
		"three perpetuals with shared and independent insurance funds": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_SmallMarginRequirement,
				constants.EthUsd_NoMarginRequirement,
				constants.IsoUsd_IsolatedMarket,
			},
			initialDeltas: map[uint32]*big.Int{
				0: big.NewInt(0),
				1: big.NewInt(0),
				3: big.NewInt(0),
			},
			incrementDeltas: map[uint32]*big.Int{
				0: big.NewInt(1000),
				1: big.NewInt(1000),
				3: big.NewInt(500),
			},
			expectedDeltas: map[uint32]*big.Int{
				0: big.NewInt(-2000), // Shared insurance fund for perpetuals 0 and 1
				1: big.NewInt(-2000), // Shared insurance fund for perpetuals 0 and 1
				3: big.NewInt(-500),  // Independent insurance fund for perpetual 3
			},
			expectError: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup keeper state.
			memClob := memclob.NewMemClobPriceTimePriority(false)
			bankMock := &mocks.BankKeeper{}
			ks := keepertest.NewClobKeepersTestContext(t, memClob, bankMock, &mocks.IndexerEventManager{})

			ctx := ks.Ctx.WithIsCheckTx(true)

			// Create the default markets.
			keepertest.CreateTestMarkets(t, ctx, ks.PricesKeeper)

			// Create liquidity tiers.
			keepertest.CreateTestLiquidityTiers(t, ctx, ks.PerpetualsKeeper)

			// Set up USDC asset in assets module.
			err := keepertest.CreateUsdcAsset(ctx, ks.AssetsKeeper)
			require.NoError(t, err)

			// Create perpetuals.
			for _, p := range tc.perpetuals {
				_, err := ks.PerpetualsKeeper.CreatePerpetual(
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
				)
				require.NoError(t, err)
			}

			// Set initial deltas.
			for perpId, initialDelta := range tc.initialDeltas {
				err = ks.ClobKeeper.IncrementCumulativeInsuranceFundDelta(ctx, perpId, initialDelta)
				if tc.expectError {
					require.Error(t, err)
					return
				}
				require.NoError(t, err)
			}

			// Increment deltas.
			for perpId, incrementDelta := range tc.incrementDeltas {
				err = ks.ClobKeeper.IncrementCumulativeInsuranceFundDelta(ctx, perpId, incrementDelta)
				if tc.expectError {
					require.Error(t, err)
					return
				}
				require.NoError(t, err)
			}

			// Get the updated deltas.
			for perpId, expectedDelta := range tc.expectedDeltas {
				delta, err := ks.ClobKeeper.GetCumulativeInsuranceFundDelta(ctx, perpId)
				require.NoError(t, err)
				require.Equal(t, expectedDelta, delta)
			}
		})
	}
}
