package clob_test

import (
	"testing"

	"github.com/cometbft/cometbft/types"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	clobtestutils "github.com/dydxprotocol/v4-chain/protocol/testutil/clob"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	feetiertypes "github.com/dydxprotocol/v4-chain/protocol/x/feetiers/types"
	prices "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

func TestDeliverTxMatchValidation(t *testing.T) {
	tests := map[string]struct {
		subaccounts       []satypes.Subaccount
		blockAdvancements []testapp.BlockAdvancementWithErrors
	}{
		"Error: partially filled IOC taker order cannot be matched twice": {
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num1_10_000USD,
				constants.Bob_Num0_10_000USD,
			},
			blockAdvancements: []testapp.BlockAdvancementWithErrors{
				{
					BlockAdvancement: testapp.BlockAdvancement{
						ShortTermOrdersAndOperations: []interface{}{
							testapp.MustScaleOrder(constants.Order_Bob_Num0_Id11_Clob1_Buy5_Price40_GTB20, testapp.DefaultGenesis()),
							testapp.MustScaleOrder(constants.Order_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB20_IOC, testapp.DefaultGenesis()),
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
				},
				{
					BlockAdvancement: testapp.BlockAdvancement{
						ShortTermOrdersAndOperations: []interface{}{
							testapp.MustScaleOrder(constants.Order_Bob_Num0_Id12_Clob1_Buy5_Price40_GTB20, testapp.DefaultGenesis()),
							testapp.MustScaleOrder(constants.Order_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB20_IOC, testapp.DefaultGenesis()),
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
					},
					ExpectedDeliverTxErrors: testapp.TxIndexesToErrors{
						0: "IOC order is already filled, remaining size is cancelled.",
					},
				},
			},
		},
		"Error: cannot match partially filled conditional IOC order as taker": {
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num0_100_000USD,
				constants.Dave_Num0_500000USD,
			},
			blockAdvancements: []testapp.BlockAdvancementWithErrors{
				{
					// place stateful orders in state, trigger conditional order in EndBlocker
					BlockAdvancement: testapp.BlockAdvancement{
						StatefulOrders: []clobtypes.Order{
							constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49999_IOC,
							constants.LongTermOrder_Dave_Num0_Id0_Clob0_Sell025BTC_Price50000_GTBT10,
						},
					},
				},
				{
					// persist match that occurs with Alice as taker and Dave ID1 as maker
					BlockAdvancement: testapp.BlockAdvancement{},
				},
				{
					// match conditional order again, this will result in an error because conditional order
					// is removed from state after being partially filled.
					BlockAdvancement: testapp.BlockAdvancement{
						ShortTermOrdersAndOperations: []interface{}{
							constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000,
							clobtestutils.NewMatchOperationRaw(
								&constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49999_IOC,
								[]clobtypes.MakerFill{
									{
										MakerOrderId: constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000.OrderId,
										FillAmount:   10,
									},
								},
							),
						},
					},
					ExpectedDeliverTxErrors: testapp.TxIndexesToErrors{
						0: clobtypes.ErrStatefulOrderDoesNotExist.Error(),
					},
				},
			},
		},
		"Success: IOC order is taker with multiple maker fills": {
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num1_10_000USD,
				constants.Bob_Num0_10_000USD,
			},
			blockAdvancements: []testapp.BlockAdvancementWithErrors{
				{
					BlockAdvancement: testapp.BlockAdvancement{
						ShortTermOrdersAndOperations: []interface{}{
							testapp.MustScaleOrder(constants.Order_Bob_Num0_Id11_Clob1_Buy5_Price40_GTB20, testapp.DefaultGenesis()),
							testapp.MustScaleOrder(constants.Order_Bob_Num0_Id12_Clob1_Buy5_Price40_GTB20, testapp.DefaultGenesis()),
							testapp.MustScaleOrder(constants.Order_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB20_IOC, testapp.DefaultGenesis()),
							clobtestutils.NewMatchOperationRaw(
								&constants.Order_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB20_IOC,
								[]clobtypes.MakerFill{
									{
										FillAmount:   5_000, // step base quantums is 1000 for ETH/USDC (ClobPair 1)
										MakerOrderId: constants.Order_Bob_Num0_Id11_Clob1_Buy5_Price40_GTB20.OrderId,
									},
									{
										FillAmount:   5_000,
										MakerOrderId: constants.Order_Bob_Num0_Id12_Clob1_Buy5_Price40_GTB20.OrderId,
									},
								},
							),
						},
					},
				},
			},
		},
		"Error: IOC order is taker in multiple matches": {
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num1_10_000USD,
				constants.Bob_Num0_10_000USD,
			},
			blockAdvancements: []testapp.BlockAdvancementWithErrors{
				{
					BlockAdvancement: testapp.BlockAdvancement{
						ShortTermOrdersAndOperations: []interface{}{
							testapp.MustScaleOrder(constants.Order_Bob_Num0_Id11_Clob1_Buy5_Price40_GTB20, testapp.DefaultGenesis()),
							testapp.MustScaleOrder(constants.Order_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB20_IOC, testapp.DefaultGenesis()),
							clobtestutils.NewMatchOperationRaw(
								&constants.Order_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB20_IOC,
								[]clobtypes.MakerFill{
									{
										FillAmount:   5_000, // step base quantums is 1000 for ETH/USDC (ClobPair 1)
										MakerOrderId: constants.Order_Bob_Num0_Id11_Clob1_Buy5_Price40_GTB20.OrderId,
									},
								},
							),
							testapp.MustScaleOrder(constants.Order_Bob_Num0_Id12_Clob1_Buy5_Price40_GTB20, testapp.DefaultGenesis()),
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
					},
					ExpectedDeliverTxErrors: testapp.TxIndexesToErrors{
						0: "IOC order is already filled, remaining size is cancelled.",
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
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *feetiertypes.GenesisState) {
						genesisState.Params = constants.PerpetualFeeParamsNoFee
					},
				)
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *prices.GenesisState) {
						*genesisState = constants.TestPricesGenesisState
					},
				)
				return genesis
			}).Build()
			ctx := tApp.InitChain()

			for i, blockAdvancement := range tc.blockAdvancements {
				ctx = blockAdvancement.AdvanceToBlock(ctx, uint32(i+2), tApp, t)
			}
		})
	}
}
