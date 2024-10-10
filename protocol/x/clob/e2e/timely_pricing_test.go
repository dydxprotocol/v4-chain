package clob_test

import (
	"fmt"
	"math/big"
	"testing"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cometbft/cometbft/types"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve"
	sdaiservertypes "github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/server/types/sdaioracle"
	testapp "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/app"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/constants"
	vetesting "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/ve"
	clobtypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/clob/types"
	feetiertypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/feetiers/types"
	perptypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/perpetuals/types"
	prices "github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices/types"
	ratelimitkeeper "github.com/StreamFinance-Protocol/stream-chain/protocol/x/ratelimit/keeper"
	satypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/subaccounts/types"
	abcitypes "github.com/cometbft/cometbft/abci/types"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestTimelyPricing(t *testing.T) {
	tests := map[string]struct {
		subaccounts []satypes.Subaccount
		orders      []clobtypes.Order

		priceUpdate map[uint32]ve.VEPricePair

		expectedErr string
	}{
		"invalid order: depends on price update at end of block": {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_10000USD,
			},
			orders: []clobtypes.Order{
				constants.LongTermOrder_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10,
			},
			priceUpdate: map[uint32]ve.VEPricePair{
				0: {
					SpotPrice: 5_000_400_000,
					PnlPrice:  5_000_400_000,
				},
			},

			expectedErr: "Stateful order collateralization check failed",
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
					func(genesisState *prices.GenesisState) {
						*genesisState = constants.TestPricesGenesisState
					},
				)
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *perptypes.GenesisState) {
						genesisState.Params = constants.PerpetualsGenesisParams
						genesisState.LiquidityTiers = constants.LiquidityTiers
						genesisState.Perpetuals = []perptypes.Perpetual{
							constants.BtcUsd_20PercentInitial_10PercentMaintenance,
						}
					},
				)
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *clobtypes.GenesisState) {
						genesisState.ClobPairs = []clobtypes.ClobPair{
							constants.ClobPair_Btc,
						}
						genesisState.LiquidationsConfig = clobtypes.LiquidationsConfig_Default
						genesisState.EquityTierLimitConfig = clobtypes.EquityTierLimitConfiguration{}
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

			rateString := sdaiservertypes.TestSDAIEventRequest.ConversionRate
			rate, conversionErr := ratelimitkeeper.ConvertStringToBigInt(rateString)

			require.NoError(t, conversionErr)

			tApp.App.RatelimitKeeper.SetSDAIPrice(tApp.App.NewUncachedContext(false, tmproto.Header{}), rate)
			tApp.App.RatelimitKeeper.SetAssetYieldIndex(tApp.App.NewUncachedContext(false, tmproto.Header{}), big.NewRat(1, 1))

			tApp.CrashingApp.RatelimitKeeper.SetSDAIPrice(tApp.CrashingApp.NewUncachedContext(false, tmproto.Header{}), rate)
			tApp.CrashingApp.RatelimitKeeper.SetAssetYieldIndex(tApp.CrashingApp.NewUncachedContext(false, tmproto.Header{}), big.NewRat(1, 1))

			tApp.NoCheckTxApp.RatelimitKeeper.SetSDAIPrice(tApp.NoCheckTxApp.NewUncachedContext(false, tmproto.Header{}), rate)
			tApp.NoCheckTxApp.RatelimitKeeper.SetAssetYieldIndex(tApp.NoCheckTxApp.NewUncachedContext(false, tmproto.Header{}), big.NewRat(1, 1))

			tApp.ParallelApp.RatelimitKeeper.SetSDAIPrice(tApp.ParallelApp.NewUncachedContext(false, tmproto.Header{}), rate)
			tApp.ParallelApp.RatelimitKeeper.SetAssetYieldIndex(tApp.ParallelApp.NewUncachedContext(false, tmproto.Header{}), big.NewRat(1, 1))

			ctx := tApp.InitChain()

			// Create all orders and add to deliverTxsOverride
			deliverTxsOverride := make([][]byte, 0)

			for _, order := range tc.orders {
				for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(
					ctx,
					tApp.App,
					*clobtypes.NewMsgPlaceOrder(order),
				) {
					resp := tApp.CheckTx(checkTx)
					require.Conditionf(t, resp.IsErr, "Expected CheckTx to fail. Response: %+v", resp)
					deliverTxsOverride = append(deliverTxsOverride, checkTx.Tx)
				}
			}

			// Add the price update to deliverTxsOverride
			_, extCommitBz, err := vetesting.GetInjectedExtendedCommitInfoForTestApp(
				&tApp.App.ConsumerKeeper,
				ctx,
				tc.priceUpdate,
				"",
				tApp.GetHeader().Height,
			)
			require.NoError(t, err)

			deliverTxsOverride = append([][]byte{extCommitBz}, deliverTxsOverride...)

			// Advance to the next block, updating the price.
			tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{
				DeliverTxsOverride: deliverTxsOverride,
				ValidateFinalizeBlock: func(
					ctx sdktypes.Context,
					request abcitypes.RequestFinalizeBlock,
					response abcitypes.ResponseFinalizeBlock,
				) (haltchain bool) {
					execResult := response.TxResults[1]
					fmt.Println("ERR STATUEFUL ORDER COLLATERALIZAITON CHECK FAILED CODE", clobtypes.ErrStatefulOrderCollateralizationCheckFailed.ABCICode())
					fmt.Println("execResult.Code", execResult.Code)
					fmt.Println("execResult.Log", execResult.Log)
					fmt.Println("tc.expectedErr", tc.expectedErr)

					require.True(t, execResult.IsErr())
					require.Equal(t, clobtypes.ErrStatefulOrderCollateralizationCheckFailed.ABCICode(), execResult.Code)
					require.Contains(t, execResult.Log, tc.expectedErr)
					return false
				},
			})
		})
	}
}
