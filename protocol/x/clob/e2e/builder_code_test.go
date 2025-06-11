package clob_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	clobtestutils "github.com/dydxprotocol/v4-chain/protocol/testutil/clob"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	testtx "github.com/dydxprotocol/v4-chain/protocol/testutil/tx"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

func TestBuilderCodeOrders(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()

	usdcDenom := "ibc/8E27BA2D5493AF5636760E354E46004562C46AB7EC0CC4C1CA14E9E20E2545B5"

	// Create test accounts with initial balances
	aliceSubaccount := tApp.App.SubaccountsKeeper.GetSubaccount(ctx, constants.Alice_Num0)
	bobSubaccount := tApp.App.SubaccountsKeeper.GetSubaccount(ctx, constants.Bob_Num0)
	builderAddress := tApp.App.SubaccountsKeeper.GetSubaccount(ctx, constants.Carl_Num0).Id.MustGetAccAddress()

	// Create orders with and without builder codes
	orderWithBuilderCode := *clobtypes.NewMsgPlaceOrder(
		clobtypes.Order{
			OrderId: clobtypes.OrderId{
				SubaccountId: *aliceSubaccount.Id,
				ClientId:     0,
				ClobPairId:   0,
			},
			Side:         clobtypes.Order_SIDE_BUY,
			Quantums:     10_000_000_000, // 1 BTC
			Subticks:     500_000_000,    // 50k USDC / BTC
			GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
			BuilderCodeParameters: &clobtypes.BuilderCodeParameters{
				BuilderAddress: builderAddress.String(),
				FeePpm:         1000, // 0.1% fee
			},
		},
	)

	orderWithoutBuilderCode := *clobtypes.NewMsgPlaceOrder(
		clobtypes.Order{
			OrderId: clobtypes.OrderId{
				SubaccountId: *bobSubaccount.Id,
				ClientId:     0,
				ClobPairId:   0,
			},
			Side:         clobtypes.Order_SIDE_SELL,
			Quantums:     10_000_000_000, // 1 BTC
			Subticks:     500_000_000,    // 50k USDC / BTC
			GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
		},
	)

	tests := map[string]struct {
		orderMsgs                         []clobtypes.MsgPlaceOrder
		expectedBuilderFee                int64
		shouldSucceed                     bool
		requestPrepareProposalTxsOverride [][]byte
		expectedFillAmount                uint64
	}{
		"Test order with builder code fills and fees are paid": {
			orderMsgs: []clobtypes.MsgPlaceOrder{
				orderWithBuilderCode,
				orderWithoutBuilderCode,
			},
			expectedBuilderFee: 50_000_000, // 0.1% of 50k USDC
			shouldSucceed:      true,
			expectedFillAmount: 10_000_000_000,
			requestPrepareProposalTxsOverride: [][]byte{
				testtx.MustGetTxBytes(&clobtypes.MsgProposedOperations{
					OperationsQueue: []clobtypes.OperationRaw{
						clobtestutils.NewMatchOperationRaw(
							&clobtypes.Order{
								OrderId:  orderWithoutBuilderCode.Order.OrderId,
								Side:     orderWithoutBuilderCode.Order.Side,
								Quantums: orderWithoutBuilderCode.Order.Quantums,
								Subticks: orderWithoutBuilderCode.Order.Subticks,
							},
							[]clobtypes.MakerFill{
								{
									FillAmount:   orderWithBuilderCode.Order.Quantums,
									MakerOrderId: orderWithBuilderCode.Order.OrderId,
								},
							},
						),
					},
				}),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// CheckTx
			builderBalancePreMatch := tApp.App.BankKeeper.GetBalance(ctx, builderAddress, usdcDenom)
			for _, order := range tc.orderMsgs {
				for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(
					ctx,
					tApp.App,
					order,
				) {
					resp := tApp.CheckTx(checkTx)
					require.True(t, resp.IsOK(), "Expected CheckTx to succeed. Response: %+v", resp)
				}
			}

			if tc.shouldSucceed {
				// Advance block to process orders
				ctx = tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{
					RequestPrepareProposalTxsOverride: [][]byte{
						testtx.MustGetTxBytes(&clobtypes.MsgProposedOperations{
							OperationsQueue: []clobtypes.OperationRaw{
								clobtestutils.NewMatchOperationRaw(
									&clobtypes.Order{
										OrderId:  orderWithoutBuilderCode.Order.OrderId,
										Side:     orderWithoutBuilderCode.Order.Side,
										Quantums: orderWithoutBuilderCode.Order.Quantums,
										Subticks: orderWithoutBuilderCode.Order.Subticks,
									},
									[]clobtypes.MakerFill{
										{
											FillAmount:   orderWithBuilderCode.Order.Quantums,
											MakerOrderId: orderWithBuilderCode.Order.OrderId,
										},
									},
								),
							},
						}),
					},
				})

				// Verify builder fee was paid
				if tc.expectedBuilderFee > 0 {
					builderBalance := tApp.App.BankKeeper.GetBalance(ctx, builderAddress, usdcDenom)
					balanaceDelta := builderBalance.Amount.Int64() - builderBalancePreMatch.Amount.Int64()
					require.Equal(t, tc.expectedBuilderFee, balanaceDelta)
				}

				// Verify orders were filled
				for _, order := range tc.orderMsgs {
					_, fillAmount, _ := tApp.App.ClobKeeper.GetOrderFillAmount(ctx, order.Order.OrderId)
					require.Equal(t, tc.expectedFillAmount, fillAmount.ToUint64())
				}
			}
		})
	}
}
