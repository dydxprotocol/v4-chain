package keeper_test

import (
	"fmt"
	"math"
	"math/big"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/memclob"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
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
				k.UpdateSubaccountLiquidationInfo(
					ctx,
					constants.Alice_Num0,
					big.NewInt(5),
					big.NewInt(-10),
				)
			},

			subaccountId: constants.Alice_Num0,

			expectedSubaccountLiquidationInfo: types.SubaccountLiquidationInfo{
				PerpetualsLiquidated:  []uint32{2},
				NotionalLiquidated:    5,
				QuantumsInsuranceLost: 10,
			},
		},
		"can get and set subaccount liquidation info with positive insurance delta": {
			setupState: func(ctx sdk.Context, k *keeper.Keeper) {
				k.MustUpdateSubaccountPerpetualLiquidated(
					ctx,
					constants.Alice_Num0,
					2,
				)
				k.UpdateSubaccountLiquidationInfo(
					ctx,
					constants.Alice_Num0,
					big.NewInt(5),
					big.NewInt(10),
				)
			},

			subaccountId: constants.Alice_Num0,

			expectedSubaccountLiquidationInfo: types.SubaccountLiquidationInfo{
				PerpetualsLiquidated:  []uint32{2},
				NotionalLiquidated:    5,
				QuantumsInsuranceLost: 0,
			},
		},
		"can get and set subaccount liquidation info with negative notional liquidated": {
			setupState: func(ctx sdk.Context, k *keeper.Keeper) {
				k.MustUpdateSubaccountPerpetualLiquidated(
					ctx,
					constants.Alice_Num0,
					2,
				)
				k.UpdateSubaccountLiquidationInfo(
					ctx,
					constants.Alice_Num0,
					big.NewInt(-5),
					big.NewInt(10),
				)
			},

			subaccountId: constants.Alice_Num0,

			expectedSubaccountLiquidationInfo: types.SubaccountLiquidationInfo{
				PerpetualsLiquidated:  []uint32{2},
				NotionalLiquidated:    5,
				QuantumsInsuranceLost: 0,
			},
		},
		"can get and set subaccount liquidation info multiple times for one subaccount": {
			setupState: func(ctx sdk.Context, k *keeper.Keeper) {
				k.MustUpdateSubaccountPerpetualLiquidated(
					ctx,
					constants.Alice_Num0,
					2,
				)
				k.UpdateSubaccountLiquidationInfo(
					ctx,
					constants.Alice_Num0,
					big.NewInt(-5),
					big.NewInt(10),
				)
				k.MustUpdateSubaccountPerpetualLiquidated(
					ctx,
					constants.Alice_Num0,
					3,
				)
				k.UpdateSubaccountLiquidationInfo(
					ctx,
					constants.Alice_Num0,
					big.NewInt(-25),
					big.NewInt(-10),
				)
				k.MustUpdateSubaccountPerpetualLiquidated(
					ctx,
					constants.Alice_Num0,
					100,
				)
				k.UpdateSubaccountLiquidationInfo(
					ctx,
					constants.Alice_Num0,
					big.NewInt(25),
					big.NewInt(-50),
				)
			},

			subaccountId: constants.Alice_Num0,

			expectedSubaccountLiquidationInfo: types.SubaccountLiquidationInfo{
				PerpetualsLiquidated:  []uint32{2, 3, 100},
				NotionalLiquidated:    55,
				QuantumsInsuranceLost: 60,
			},
		},
		"can get and set subaccount liquidation info for multiple subaccounts": {
			setupState: func(ctx sdk.Context, k *keeper.Keeper) {
				k.MustUpdateSubaccountPerpetualLiquidated(
					ctx,
					constants.Alice_Num0,
					2,
				)
				k.UpdateSubaccountLiquidationInfo(
					ctx,
					constants.Alice_Num0,
					big.NewInt(-5),
					big.NewInt(10),
				)
				k.MustUpdateSubaccountPerpetualLiquidated(
					ctx,
					constants.Alice_Num1,
					3,
				)
				k.UpdateSubaccountLiquidationInfo(
					ctx,
					constants.Alice_Num1,
					big.NewInt(-25),
					big.NewInt(-10),
				)
				k.MustUpdateSubaccountPerpetualLiquidated(
					ctx,
					constants.Bob_Num0,
					100,
				)
				k.UpdateSubaccountLiquidationInfo(
					ctx,
					constants.Bob_Num0,
					big.NewInt(25),
					big.NewInt(-50),
				)
			},

			subaccountId: constants.Bob_Num0,

			expectedSubaccountLiquidationInfo: types.SubaccountLiquidationInfo{
				PerpetualsLiquidated:  []uint32{100},
				NotionalLiquidated:    25,
				QuantumsInsuranceLost: 50,
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

func TestUpdateSubaccountLiquidationInfo_NotionalLiquidatedOverflowPanics(t *testing.T) {
	memClob := memclob.NewMemClobPriceTimePriority(false)
	bankMock := &mocks.BankKeeper{}
	ks := keepertest.NewClobKeepersTestContext(
		t,
		memClob,
		bankMock,
		&mocks.IndexerEventManager{},
	)

	subaccountId := constants.Alice_Num0
	ks.ClobKeeper.MustUpdateSubaccountPerpetualLiquidated(
		ks.Ctx,
		subaccountId,
		0,
	)
	ks.ClobKeeper.UpdateSubaccountLiquidationInfo(
		ks.Ctx,
		subaccountId,
		new(big.Int).SetUint64(math.MaxUint64),
		big.NewInt(50),
	)

	require.PanicsWithError(
		t,
		fmt.Sprintf(
			"Notional liquidated update for subaccount %v overflows uint64: integer overflow",
			subaccountId,
		),
		func() {
			// Run the test and verify expectations.
			ks.ClobKeeper.MustUpdateSubaccountPerpetualLiquidated(
				ks.Ctx,
				subaccountId,
				1,
			)
			ks.ClobKeeper.UpdateSubaccountLiquidationInfo(
				ks.Ctx,
				subaccountId,
				big.NewInt(1),
				big.NewInt(50),
			)
		},
	)
}

func TestUpdateSubaccountLiquidationInfo_QuantumInsuranceLostOverflowPanics(t *testing.T) {
	memClob := memclob.NewMemClobPriceTimePriority(false)
	bankMock := &mocks.BankKeeper{}
	ks := keepertest.NewClobKeepersTestContext(t, memClob, bankMock, &mocks.IndexerEventManager{})

	subaccountId := constants.Alice_Num0
	ks.ClobKeeper.MustUpdateSubaccountPerpetualLiquidated(
		ks.Ctx,
		subaccountId,
		0,
	)
	ks.ClobKeeper.UpdateSubaccountLiquidationInfo(
		ks.Ctx,
		subaccountId,
		big.NewInt(50),
		constants.BigNegMaxUint64(),
	)

	require.PanicsWithError(
		t,
		fmt.Sprintf(
			"Quantums insurance lost update for subaccount %v overflows uint64: integer overflow",
			subaccountId,
		),
		func() {
			ks.ClobKeeper.MustUpdateSubaccountPerpetualLiquidated(
				ks.Ctx,
				subaccountId,
				1,
			)
			ks.ClobKeeper.UpdateSubaccountLiquidationInfo(
				ks.Ctx,
				subaccountId,
				big.NewInt(50),
				big.NewInt(-1),
			)
		},
	)
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

	ks.ClobKeeper.UpdateSubaccountLiquidationInfo(
		ks.Ctx,
		subaccountId,
		big.NewInt(50),
		big.NewInt(20),
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
			ks.ClobKeeper.UpdateSubaccountLiquidationInfo(
				ks.Ctx,
				subaccountId,
				big.NewInt(20),
				big.NewInt(-1),
			)
		},
	)
}
