package clob_test

import (
	"testing"

	"github.com/cometbft/cometbft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/require"
)

func TestReduceOnlyOrders(t *testing.T) {
	tests := map[string]struct {
		subaccounts          []satypes.Subaccount
		ordersForFirstBlock  []clobtypes.Order
		ordersForSecondBlock []clobtypes.Order

		expectedOrderOnMemClob  map[clobtypes.OrderId]bool
		expectedOrderFillAmount map[clobtypes.OrderId]uint64
		expectedSubaccounts     []satypes.Subaccount
	}{
		"IOC Reduce only order partially matches short term order same block, maker order fully filled": {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_100000USD,
				constants.Alice_Num1_1BTC_Long_500_000USD,
			},
			ordersForFirstBlock: []clobtypes.Order{
				MustScaleOrder(
					constants.Order_Carl_Num0_Id0_Clob0_Buy10_Price500000_GTB20,
					testapp.DefaultGenesis(),
				),
				MustScaleOrder(
					constants.Order_Alice_Num1_Id1_Clob0_Sell15_Price500000_GTB20_IOC_RO,
					testapp.DefaultGenesis(),
				),
			},
			ordersForSecondBlock: []clobtypes.Order{},

			expectedOrderOnMemClob: map[clobtypes.OrderId]bool{
				constants.Order_Carl_Num0_Id0_Clob0_Buy10_Price500000_GTB20.OrderId:          false,
				constants.Order_Alice_Num1_Id1_Clob0_Sell15_Price500000_GTB20_IOC_RO.OrderId: false,
			},
			expectedOrderFillAmount: map[clobtypes.OrderId]uint64{
				constants.Order_Carl_Num0_Id0_Clob0_Buy10_Price500000_GTB20.OrderId:          100,
				constants.Order_Alice_Num1_Id1_Clob0_Sell15_Price500000_GTB20_IOC_RO.OrderId: 100,
			},
			expectedSubaccounts: []satypes.Subaccount{
				{
					Id: &constants.Carl_Num0,
					AssetPositions: []*satypes.AssetPosition{
						{
							AssetId:  0,
							Quantums: dtypes.NewInt(95_000_550_000),
						},
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						{
							PerpetualId:  0,
							Quantums:     dtypes.NewInt(100),
							FundingIndex: dtypes.NewInt(0),
						},
					},
				},
				{
					Id: &constants.Alice_Num1,
					AssetPositions: []*satypes.AssetPosition{
						{
							AssetId:  0,
							Quantums: dtypes.NewInt(504_997_500_000),
						},
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						{
							PerpetualId:  0,
							Quantums:     dtypes.NewInt(99_999_900),
							FundingIndex: dtypes.NewInt(0),
						},
					},
				},
			},
		},
		"IOC Reduce only order partially matches short term order second block, maker order fully filled": {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_100000USD,
				constants.Alice_Num1_1BTC_Long_500_000USD,
			},
			ordersForFirstBlock: []clobtypes.Order{
				MustScaleOrder(
					constants.Order_Carl_Num0_Id0_Clob0_Buy10_Price500000_GTB20,
					testapp.DefaultGenesis(),
				),
			},
			ordersForSecondBlock: []clobtypes.Order{
				MustScaleOrder(
					constants.Order_Alice_Num1_Id1_Clob0_Sell15_Price500000_GTB20_IOC_RO,
					testapp.DefaultGenesis(),
				),
			},

			expectedOrderOnMemClob: map[clobtypes.OrderId]bool{
				constants.Order_Carl_Num0_Id0_Clob0_Buy10_Price500000_GTB20.OrderId:          false,
				constants.Order_Alice_Num1_Id1_Clob0_Sell15_Price500000_GTB20_IOC_RO.OrderId: false,
			},
			expectedOrderFillAmount: map[clobtypes.OrderId]uint64{
				constants.Order_Carl_Num0_Id0_Clob0_Buy10_Price500000_GTB20.OrderId:          100,
				constants.Order_Alice_Num1_Id1_Clob0_Sell15_Price500000_GTB20_IOC_RO.OrderId: 100,
			},
			expectedSubaccounts: []satypes.Subaccount{
				{
					Id: &constants.Carl_Num0,
					AssetPositions: []*satypes.AssetPosition{
						{
							AssetId:  0,
							Quantums: dtypes.NewInt(95_000_550_000),
						},
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						{
							PerpetualId:  0,
							Quantums:     dtypes.NewInt(100),
							FundingIndex: dtypes.NewInt(0),
						},
					},
				},
				{
					Id: &constants.Alice_Num1,
					AssetPositions: []*satypes.AssetPosition{
						{
							AssetId:  0,
							Quantums: dtypes.NewInt(504_997_500_000),
						},
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						{
							PerpetualId:  0,
							Quantums:     dtypes.NewInt(99_999_900),
							FundingIndex: dtypes.NewInt(0),
						},
					},
				},
			},
		},
		"IOC Reduce only order partially matches short term order second block, maker order partially filled": {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_100000USD,
				constants.Alice_Num1_1BTC_Long_500_000USD,
			},
			ordersForFirstBlock: []clobtypes.Order{
				MustScaleOrder(
					constants.Order_Carl_Num0_Id0_Clob0_Buy80_Price500000_GTB20,
					testapp.DefaultGenesis(),
				),
			},
			ordersForSecondBlock: []clobtypes.Order{
				MustScaleOrder(
					constants.Order_Alice_Num1_Id1_Clob0_Sell15_Price500000_GTB20_IOC_RO,
					testapp.DefaultGenesis(),
				),
			},

			expectedOrderOnMemClob: map[clobtypes.OrderId]bool{
				constants.Order_Carl_Num0_Id0_Clob0_Buy80_Price500000_GTB20.OrderId:          true,
				constants.Order_Alice_Num1_Id1_Clob0_Sell15_Price500000_GTB20_IOC_RO.OrderId: false,
			},
			expectedOrderFillAmount: map[clobtypes.OrderId]uint64{
				constants.Order_Carl_Num0_Id0_Clob0_Buy80_Price500000_GTB20.OrderId:          150,
				constants.Order_Alice_Num1_Id1_Clob0_Sell15_Price500000_GTB20_IOC_RO.OrderId: 150,
			},
			expectedSubaccounts: []satypes.Subaccount{
				{
					Id: &constants.Carl_Num0,
					AssetPositions: []*satypes.AssetPosition{
						{
							AssetId:  0,
							Quantums: dtypes.NewInt(9_250_0825_000),
						},
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						{
							PerpetualId:  0,
							Quantums:     dtypes.NewInt(150),
							FundingIndex: dtypes.NewInt(0),
						},
					},
				},
				{
					Id: &constants.Alice_Num1,
					AssetPositions: []*satypes.AssetPosition{
						{
							AssetId:  0,
							Quantums: dtypes.NewInt(507_496_250_000),
						},
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						{
							PerpetualId:  0,
							Quantums:     dtypes.NewInt(99_999_850),
							FundingIndex: dtypes.NewInt(0),
						},
					},
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(func() (genesis types.GenesisDoc) {
				genesis = testapp.DefaultGenesis()
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *satypes.GenesisState) {
						genesisState.Subaccounts = tc.subaccounts
					},
				)
				return genesis
			}).Build()
			ctx := tApp.InitChain()

			// Create all orders.
			deliverTxsOverride := make([][]byte, 0)
			deliverTxsOverride = append(
				deliverTxsOverride,
				constants.ValidEmptyMsgProposedOperationsTxBytes,
			)

			for _, order := range tc.ordersForFirstBlock {
				for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(
					ctx,
					tApp.App,
					*clobtypes.NewMsgPlaceOrder(order),
				) {
					resp := tApp.CheckTx(checkTx)
					require.Conditionf(t, resp.IsOK, "Expected CheckTx to succeed. Response: %+v", resp)

					if order.IsStatefulOrder() {
						deliverTxsOverride = append(deliverTxsOverride, checkTx.Tx)
					}
				}
			}

			// Add an empty premium vote.
			deliverTxsOverride = append(deliverTxsOverride, constants.EmptyMsgAddPremiumVotesTxBytes)

			ctx = tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{
				DeliverTxsOverride: deliverTxsOverride,
			})

			// Place orders for second block
			for _, order := range tc.ordersForSecondBlock {
				for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(
					ctx,
					tApp.App,
					*clobtypes.NewMsgPlaceOrder(order),
				) {
					resp := tApp.CheckTx(checkTx)
					require.Conditionf(t, resp.IsOK, "Expected CheckTx to succeed. Response: %+v", resp)

					if order.IsStatefulOrder() {
						deliverTxsOverride = append(deliverTxsOverride, checkTx.Tx)
					}
				}
			}

			// Verify expectations.
			for orderId, exists := range tc.expectedOrderOnMemClob {
				_, existsOnMemclob := tApp.App.ClobKeeper.MemClob.GetOrder(ctx, orderId)
				// _, found := tApp.App.ClobKeeper.GetLongTermOrderPlacement(ctx, orderId)
				require.Equal(t, exists, existsOnMemclob)
			}

			for orderId, expectedFillAmount := range tc.expectedOrderFillAmount {
				exists, fillAmount, _ := tApp.App.ClobKeeper.GetOrderFillAmount(ctx, orderId)
				require.True(t, exists)
				require.Equal(t, expectedFillAmount, fillAmount.ToUint64())
			}

			for _, subaccount := range tc.expectedSubaccounts {
				actualSubaccount := tApp.App.SubaccountsKeeper.GetSubaccount(ctx, *subaccount.Id)
				require.Equal(t, subaccount, actualSubaccount)
			}
		})
	}
}

func TestReduceOnlyOrderFailure(t *testing.T) {
	tests := map[string]struct {
		subaccounts []satypes.Subaccount
		orders      []clobtypes.Order
		errorMsg    []string
	}{
		"Zero perpetual position subaccount position cannot place sell RO order": {
			orders: []clobtypes.Order{
				MustScaleOrder(
					constants.Order_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB20_FOK_RO,
					testapp.DefaultGenesis(),
				),
			},
			errorMsg: []string{
				clobtypes.ErrReduceOnlyWouldIncreasePositionSize.Error(),
			},
		},
		"Zero perpetual position subaccount position cannot place buy RO order": {
			orders: []clobtypes.Order{
				MustScaleOrder(
					constants.Order_Alice_Num1_Id1_Clob1_Buy10_Price15_GTB20_FOK_RO,
					testapp.DefaultGenesis(),
				),
			},
			errorMsg: []string{
				clobtypes.ErrReduceOnlyWouldIncreasePositionSize.Error(),
			},
		},
		"FOK Reduce only order is placed but does not match immediately and is cancelled.": {
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num1_1BTC_Short_100_000USD,
			},
			orders: []clobtypes.Order{
				MustScaleOrder(
					constants.Order_Alice_Num1_Id1_Clob0_Buy10_Price15_GTB20_FOK_RO,
					testapp.DefaultGenesis(),
				),
			},
			errorMsg: []string{
				clobtypes.ErrFokOrderCouldNotBeFullyFilled.Error(),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(func() (genesis types.GenesisDoc) {
				genesis = testapp.DefaultGenesis()
				if len(tc.subaccounts) > 0 {
					testapp.UpdateGenesisDocWithAppStateForModule(
						&genesis,
						func(genesisState *satypes.GenesisState) {
							genesisState.Subaccounts = tc.subaccounts
						},
					)
				}
				return genesis
			}).Build()
			ctx := tApp.InitChain()

			for idx, order := range tc.orders {
				for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(
					ctx,
					tApp.App,
					*clobtypes.NewMsgPlaceOrder(order),
				) {
					resp := tApp.CheckTx(checkTx)

					if tc.errorMsg[idx] == "" {
						require.Conditionf(t, resp.IsOK, "Expected CheckTx to succeed. Response: %+v", resp)
					} else {
						require.Conditionf(t, resp.IsErr, "Expected CheckTx to error. Response: %+v", resp)
						require.Contains(
							t,
							resp.Log,
							tc.errorMsg[idx],
						)
					}
				}
			}
		})
	}
}

func TestReduceOnlyOrderReplacement(t *testing.T) {
	tests := map[string]struct {
		subaccounts []satypes.Subaccount
		firstOrders []clobtypes.Order

		modifyFillAmountStateAfterFirstBlock func(ctx sdk.Context, tApp *testapp.TestApp)

		secondOrders       []clobtypes.Order
		secondOrdersErrors []string
	}{
		`A long position gets partially filled by sell order A to a short position.
		Order A is replaced by a BUY-side FOK reduce only order, which is an invalid replacement because
		original order was partially filled.`: {
			subaccounts: []satypes.Subaccount{
				{
					Id: &constants.Alice_Num1,
					AssetPositions: []*satypes.AssetPosition{
						&constants.Usdc_Asset_100_000,
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						{
							PerpetualId: 0,
							Quantums:    dtypes.NewInt(60), // 60 quantums of BTC long
						},
					},
				},
			},

			// Regular order on the opposite side of the following replacement RO order.
			firstOrders: []clobtypes.Order{
				constants.Order_Alice_Num1_Id0_Clob0_Sell100_Price50000_GTB20,
			},

			modifyFillAmountStateAfterFirstBlock: func(ctx sdk.Context, tApp *testapp.TestApp) {
				// Sell 80 quantums of the original order.
				// In state, the amount of BTC is still 60 quantums long.
				// After this partial fill, the net amount of btc is 60 quantums long - 80 quantums = 20 quantums short.
				// Thus, only buy reduce only orders are valid to reduce the short position size.
				tApp.App.ClobKeeper.SetOrderFillAmount(
					ctx,
					constants.Order_Alice_Num1_Id0_Clob0_Sell100_Price50000_GTB20.OrderId,
					satypes.BaseQuantums(80),
					uint32(22),
				)
			},

			secondOrders: []clobtypes.Order{
				// Currently, IOC/FOK replacement orders for orders that are partially filled are not allowed.
				// Because the order was partially filled, this results in a ErrInvalidReplacement error.
				// If IOC/FOK replacement orders were allowed, this buy RO order should fail because the
				// current position size is 20 quantums short even though in state, there is 60 quantums long.
				constants.Order_Alice_Num1_Id0_Clob0_Buy110_Price50000_GTB21_FOK_RO,
			},
			secondOrdersErrors: []string{
				clobtypes.ErrInvalidReplacement.Error(),
			},
		},
		`A long position and sell order A with 0 fills is replaced by a BUY-side FOK reduce only order,
		 which is a valid replacement. Fails validation because it would increase position size.`: {
			subaccounts: []satypes.Subaccount{
				{
					Id: &constants.Alice_Num1,
					AssetPositions: []*satypes.AssetPosition{
						&constants.Usdc_Asset_100_000,
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						{
							PerpetualId: 0,
							Quantums:    dtypes.NewInt(60), // 60 quantums of BTC long
						},
					},
				},
			},

			// Regular order on the opposite side of the following replacement RO order.
			firstOrders: []clobtypes.Order{
				constants.Order_Alice_Num1_Id0_Clob0_Sell100_Price50000_GTB20,
			},

			// No fills.
			modifyFillAmountStateAfterFirstBlock: func(ctx sdk.Context, tApp *testapp.TestApp) {},

			secondOrders: []clobtypes.Order{
				// Since the original order had no fills, this IOC/FOK order passes through.
				// But it is invalid because it would increase the positive perpetual position.
				constants.Order_Alice_Num1_Id0_Clob0_Buy110_Price50000_GTB21_FOK_RO,
			},
			secondOrdersErrors: []string{
				clobtypes.ErrReduceOnlyWouldIncreasePositionSize.Error(),
			},
		},
		`A long position and sell order A with 0 fills is replaced by a BUY-side FOK reduce only order,
		which is a valid replacement. Replacement order would decrease position size. Fails because FOK
		order is not fully filled.`: {
			subaccounts: []satypes.Subaccount{
				{
					Id: &constants.Alice_Num1,
					AssetPositions: []*satypes.AssetPosition{
						&constants.Usdc_Asset_100_000,
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						{
							PerpetualId: 0,
							Quantums:    dtypes.NewInt(60), // 60 quantums of BTC long
						},
					},
				},
				constants.Carl_Num0_100000USD,
			},

			// Regular order on the opposite side of the following replacement RO order.
			// Won't match Carl's order.
			firstOrders: []clobtypes.Order{
				constants.Order_Alice_Num1_Id0_Clob0_Sell100_Price100000_GTB20,
			},

			// No fills.
			modifyFillAmountStateAfterFirstBlock: func(ctx sdk.Context, tApp *testapp.TestApp) {},

			secondOrders: []clobtypes.Order{
				// An order that would match the below FOK order.
				constants.Order_Carl_Num0_Id0_Clob0_Buy110_Price50000_GTB20,
				// Since the original order had no fills, this IOC/FOK order passes through.
				// It is valid because it reduces the position size.
				// However, it is resized to 60 because the current position size is 60.
				// Thus, it fails because this FOK order is not fully filled.
				constants.Order_Alice_Num1_Id0_Clob0_Sell110_Price50000_GTB21_FOK_RO,
			},
			secondOrdersErrors: []string{
				"",
				clobtypes.ErrFokOrderCouldNotBeFullyFilled.Error(),
			},
		},
		`A long position and sell order A with 0 fills is replaced by a BUY-side IOC reduce only order,
		which is a valid replacement. Passes validation because it would decrease position size. IOC
		order is resized and succeeds to zero out all subaccount positions.`: {
			subaccounts: []satypes.Subaccount{
				{
					Id: &constants.Alice_Num1,
					AssetPositions: []*satypes.AssetPosition{
						&constants.Usdc_Asset_100_000,
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						{
							PerpetualId: 0,
							Quantums:    dtypes.NewInt(60), // 60 quantums of BTC long
						},
					},
				},
				constants.Carl_Num0_100000USD,
			},

			// Regular order on the opposite side of the following replacement RO order.
			// Won't match Carl's order.
			firstOrders: []clobtypes.Order{
				constants.Order_Alice_Num1_Id0_Clob0_Sell100_Price100000_GTB20,
			},

			// No fills.
			modifyFillAmountStateAfterFirstBlock: func(ctx sdk.Context, tApp *testapp.TestApp) {},

			secondOrders: []clobtypes.Order{
				// An order that would match the below FOK order.
				constants.Order_Carl_Num0_Id0_Clob0_Buy110_Price50000_GTB20,
				// Since the original order had no fills, this IOC/FOK order passes through.
				// It is valid because it reduces the position size.
				// However, it is resized to 60 because the current position size is 60.
				// Succeeds and matches above order.
				constants.Order_Alice_Num1_Id0_Clob0_Sell110_Price50000_GTB21_IOC_RO,
			},
			secondOrdersErrors: []string{
				"",
				"",
			},
		},
		`A long position gets partially filled by sell order A to a long position.
		Order A is replaced by a BUY-side FOK reduce only order, which is an invalid replacement because
		it would increase the long position size.`: {
			subaccounts: []satypes.Subaccount{
				{
					Id: &constants.Alice_Num1,
					AssetPositions: []*satypes.AssetPosition{
						&constants.Usdc_Asset_100_000,
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						{
							PerpetualId: 0,
							Quantums:    dtypes.NewInt(60), // 60 quantums of BTC long
						},
					},
				},
			},

			// Regular order on the opposite side of the following replacement RO order.
			firstOrders: []clobtypes.Order{
				constants.Order_Alice_Num1_Id0_Clob0_Sell100_Price50000_GTB20,
			},

			modifyFillAmountStateAfterFirstBlock: func(ctx sdk.Context, tApp *testapp.TestApp) {
				// Sell 20 quantums of the original order.
				// In state, the amount of BTC is still 60 quantums long.
				// After this partial fill, the net amount of btc is 60 quantums long - 20 quantums = 40 quantums long.
				// Thus, only sell reduce only orders are valid to reduce the short position size.
				tApp.App.ClobKeeper.SetOrderFillAmount(
					ctx,
					constants.Order_Alice_Num1_Id0_Clob0_Sell100_Price50000_GTB20.OrderId,
					satypes.BaseQuantums(20),
					uint32(22),
				)
			},

			secondOrders: []clobtypes.Order{
				// Invalid because it would increase the positive position size.
				constants.Order_Alice_Num1_Id0_Clob0_Buy110_Price50000_GTB21_FOK_RO,
			},
			secondOrdersErrors: []string{
				clobtypes.ErrReduceOnlyWouldIncreasePositionSize.Error(),
			},
		},
		`A long position gets partially filled by sell order A to a long position.
		Order A is replaced by a SELL-side FOK reduce only order, which would reduce position size.
		However, it is an invalid replacement because the original order was already partially filled.`: {
			subaccounts: []satypes.Subaccount{
				{
					Id: &constants.Alice_Num1,
					AssetPositions: []*satypes.AssetPosition{
						&constants.Usdc_Asset_100_000,
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						{
							PerpetualId: 0,
							Quantums:    dtypes.NewInt(60), // 60 quantums of BTC long
						},
					},
				},
			},

			// Regular order on the opposite side of the following replacement RO order.
			firstOrders: []clobtypes.Order{
				constants.Order_Alice_Num1_Id0_Clob0_Sell100_Price50000_GTB20,
			},

			modifyFillAmountStateAfterFirstBlock: func(ctx sdk.Context, tApp *testapp.TestApp) {
				// Sell 20 quantums of the original order.
				// In state, the amount of BTC is still 60 quantums long.
				// After this partial fill, the net amount of btc is 60 quantums long - 20 quantums = 40 quantums long.
				// Thus, only sell reduce only orders are valid to reduce the short position size.
				tApp.App.ClobKeeper.SetOrderFillAmount(
					ctx,
					constants.Order_Alice_Num1_Id0_Clob0_Sell100_Price50000_GTB20.OrderId,
					satypes.BaseQuantums(20),
					uint32(22),
				)
			},

			secondOrders: []clobtypes.Order{
				// Reduces position size, however is invalid replacement because of partial fills.
				constants.Order_Alice_Num1_Id0_Clob0_Sell110_Price50000_GTB21_FOK_RO,
			},
			secondOrdersErrors: []string{
				clobtypes.ErrInvalidReplacement.Error(),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).
				WithGenesisDocFn(func() (genesis types.GenesisDoc) {
					genesis = testapp.DefaultGenesis()
					if len(tc.subaccounts) > 0 {
						testapp.UpdateGenesisDocWithAppStateForModule(
							&genesis,
							func(genesisState *satypes.GenesisState) {
								genesisState.Subaccounts = tc.subaccounts
							},
						)
					}
					return genesis
				}).
				WithNonDeterminismChecksEnabled(false).
				Build()
			ctx := tApp.InitChain()

			// place first set of orders.
			for _, order := range tc.firstOrders {
				for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(
					ctx,
					tApp.App,
					*clobtypes.NewMsgPlaceOrder(order),
				) {
					resp := tApp.CheckTx(checkTx)
					require.Conditionf(t, resp.IsOK, "Expected CheckTx to succeed. Response: %+v", resp)
				}
			}
			// modify fill amounts
			if tc.modifyFillAmountStateAfterFirstBlock != nil {
				tc.modifyFillAmountStateAfterFirstBlock(ctx, tApp)
			}

			// place second set of orders.
			for idx, order := range tc.secondOrders {
				for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(
					ctx,
					tApp.App,
					*clobtypes.NewMsgPlaceOrder(order),
				) {
					resp := tApp.CheckTx(checkTx)

					if tc.secondOrdersErrors[idx] == "" {
						require.Conditionf(t, resp.IsOK, "Expected CheckTx to succeed. Response: %+v", resp)
					} else {
						require.Conditionf(t, resp.IsErr, "Expected CheckTx to error. Response: %+v", resp)
						require.Contains(
							t,
							resp.Log,
							tc.secondOrdersErrors[idx],
						)
					}
				}
			}
		})
	}
}
