package clob_test

import (
	"testing"

	"github.com/cometbft/cometbft/types"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	clobtestutils "github.com/dydxprotocol/v4-chain/protocol/testutil/clob"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	feetiertypes "github.com/dydxprotocol/v4-chain/protocol/x/feetiers/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

func TestOrderMatches(t *testing.T) {
	tests := map[string]struct {
		subaccounts       []satypes.Subaccount
		blockAdvancements []clobtestutils.BlockAdvancementWithError
	}{
		"Error: partially filled IOC taker order cannot be matched twice": {
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num1_10_000USD,
				constants.Bob_Num0_10_000USD,
			},
			blockAdvancements: []clobtestutils.BlockAdvancementWithError{
				{
					BlockAdvancement: clobtestutils.BlockAdvancement{
						OrdersAndOperations: []interface{}{
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
				},
				{
					BlockAdvancement: clobtestutils.BlockAdvancement{
						OrdersAndOperations: []interface{}{
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
					},
					ExpectedDeliverTxError: "IOC/FOK order is already filled, remaining size is cancelled.",
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
				ctx = blockAdvancement.AdvanceToBlock(ctx, uint32(i+2), &tApp, t)
			}
		})
	}
}
