package clob_test

import (
	"testing"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/types"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/app"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	clobtestutils "github.com/dydxprotocol/v4-chain/protocol/testutil/clob"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	testtx "github.com/dydxprotocol/v4-chain/protocol/testutil/tx"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	feetiertypes "github.com/dydxprotocol/v4-chain/protocol/x/feetiers/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/require"
)

// BlockAdvancement holds orders and matches to be placed in a block. Using this struct and building
// the ops queue with the getOperationsQueue helper function allows us to build the operations queue
// without going through CheckTx and, therefore, not affect the local memclob state. This also allows us to propose
// an invalid set of operations that an honest validator would not generate.
type BlockAdvancement struct {
	ordersAndMatches []interface{} // should hold Order and OperationRaw. OperationRaw are assumed to be matches.
	expectedError    string
}

func (b BlockAdvancement) getOperationsQueue(ctx sdktypes.Context, app *app.App) []clobtypes.OperationRaw {
	operationsQueue := make([]clobtypes.OperationRaw, len(b.ordersAndMatches))
	for i, orderOrMatch := range b.ordersAndMatches {
		switch castedValue := orderOrMatch.(type) {
		case clobtypes.Order:
			order := castedValue
			requestTxs := testapp.MustMakeCheckTxsWithClobMsg(
				ctx,
				app,
				*clobtypes.NewMsgPlaceOrder(order),
			)
			operationsQueue[i] = clobtypes.OperationRaw{
				Operation: &clobtypes.OperationRaw_ShortTermOrderPlacement{
					ShortTermOrderPlacement: requestTxs[0].Tx,
				},
			}
		case clobtypes.OperationRaw:
			operationsQueue[i] = castedValue
		default:
			panic("invalid type")
		}
	}

	return operationsQueue
}

func TestOrderMatches(t *testing.T) {
	tests := map[string]struct {
		subaccounts       []satypes.Subaccount
		blockAdvancements []BlockAdvancement
	}{
		"Error: partially filled IOC taker order cannot be matched twice": {
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num1_10_000USD,
				constants.Bob_Num0_10_000USD,
			},
			blockAdvancements: []BlockAdvancement{
				{
					ordersAndMatches: []interface{}{
						MustScaleOrder(constants.Order_Bob_Num0_Id11_Clob1_Buy5_Price40_GTB20, testapp.DefaultGenesis()),
						MustScaleOrder(constants.Order_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB20_IOC, testapp.DefaultGenesis()),
						clobtestutils.NewMatchOperationRaw(
							&constants.Order_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB20_IOC,
							[]clobtypes.MakerFill{
								{
									FillAmount:   5_000, // step base quantums is 1000 for ETH/USDC (ClobPair 1)
									MakerOrderId: constants.Order_Bob_Num0_Id11_Clob1_Buy5_Price40_GTB20.OrderId,
								},
							},
						),
					},
				},
				{
					ordersAndMatches: []interface{}{
						MustScaleOrder(constants.Order_Bob_Num0_Id12_Clob1_Buy5_Price40_GTB20, testapp.DefaultGenesis()),
						MustScaleOrder(constants.Order_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB20_IOC, testapp.DefaultGenesis()),
						clobtestutils.NewMatchOperationRaw(
							&constants.Order_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB20_IOC,
							[]clobtypes.MakerFill{
								{
									FillAmount:   5_000,
									MakerOrderId: constants.Order_Bob_Num0_Id12_Clob1_Buy5_Price40_GTB20.OrderId,
								},
							},
						),
					},
					expectedError: "ImmediateOrCancel order is already filled, remaining size is cancelled.",
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder().WithTesting(t).WithGenesisDocFn(func() (genesis types.GenesisDoc) {
				genesis = testapp.DefaultGenesis()
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *satypes.GenesisState) {
						genesisState.Subaccounts = tc.subaccounts
					},
				)
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *feetiertypes.GenesisState) {
						genesisState.Params = constants.PerpetualFeeParamsNoFee
					},
				)
				return genesis
			}).Build()
			ctx := tApp.InitChain()

			for i, blockAdvancement := range tc.blockAdvancements {
				msgProposedOperations := &clobtypes.MsgProposedOperations{
					OperationsQueue: blockAdvancement.getOperationsQueue(ctx, tApp.App),
				}
				ctx = tApp.AdvanceToBlock(uint32(2+i), testapp.AdvanceToBlockOptions{
					DeliverTxsOverride: [][]byte{testtx.MustGetTxBytes(msgProposedOperations)},
					ValidateDeliverTxs: func(
						ctx sdktypes.Context,
						request abcitypes.RequestDeliverTx,
						response abcitypes.ResponseDeliverTx,
					) (haltchain bool) {
						if blockAdvancement.expectedError != "" {
							require.True(t, response.IsErr())
							require.Contains(t, response.Log, blockAdvancement.expectedError)
						} else {
							require.True(t, response.IsOK())
						}
						return false
					},
				})
			}
		})
	}
}
