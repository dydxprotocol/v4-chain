package keeper_test

import (
	"math/big"
	"testing"
	"time"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/dydxprotocol/v4-chain/protocol/app/process"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	clobtest "github.com/dydxprotocol/v4-chain/protocol/testutil/clob"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/memclob"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	feetierstypes "github.com/dydxprotocol/v4-chain/protocol/x/feetiers/types"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestRecordMevMetrics(t *testing.T) {
	// Set the maximum spread to 10%.
	keeper.MAX_SPREAD_BEFORE_FALLING_BACK_TO_ORACLE = new(big.Rat).SetFrac64(1, 10)

	tests := map[string]struct {
		// Setup.
		subaccounts       []satypes.Subaccount
		clobPairs         []types.ClobPair
		feeParams         feetierstypes.PerpetualFeeParams
		perpetuals        []perptypes.Perpetual
		restingOrders     []types.Order
		liquidationOrders []types.LiquidationOrder
		intrablockMsgs    []sdk.Msg

		// Mocks.
		setupPerpetualKeeperMocks func(perpKeeper *mocks.ProcessPerpetualKeeper)

		// Input.
		proposedOperations *types.MsgProposedOperations

		// Expectations.
		expectedMev                          float32
		expectedMidPrice                     uint64
		expectedValidatorNumFills            int
		expectedValidatorVolumeQuoteQuantums *big.Int
		expectedProposerNumFills             int
		expectedProposerVolumeQuoteQuantums  *big.Int
	}{
		"Case 1 - capturing the spread": {
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num0_10_000USD,
				constants.Bob_Num0_10_000USD,
				constants.Carl_Num0_10000USD,

				// Dave_Num1 is the subaccount controlled by the block proposer.
				constants.Dave_Num1_10_000USD,
			},
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc,
			},
			feeParams: constants.PerpetualFeeParamsNoFee,
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},
			restingOrders: []types.Order{
				constants.Order_Alice_Num0_Id0_Clob0_Sell200BTC_Price101_GTB20,
				constants.Order_Bob_Num0_Id1_Clob0_Buy100BTC_Price99_GTB20,
			},
			intrablockMsgs: []sdk.Msg{
				types.NewMsgPlaceOrder(
					constants.Order_Carl_Num0_Id1_Clob0_Buy100BTC_Price101_GTB20,
				),
			},
			proposedOperations: &types.MsgProposedOperations{
				OperationsQueue: []types.OperationRaw{
					// X sells 100 at $101 instead of A.
					clobtest.NewShortTermOrderPlacementOperationRaw(constants.Order_Dave_Num1_Id0_Clob0_Sell100BTC_Price101_GTB20),
					clobtest.NewShortTermOrderPlacementOperationRaw(constants.Order_Carl_Num0_Id1_Clob0_Buy100BTC_Price101_GTB20),
					clobtest.NewMatchOperationRaw(
						&constants.Order_Carl_Num0_Id1_Clob0_Buy100BTC_Price101_GTB20,
						[]types.MakerFill{
							{
								MakerOrderId: constants.Order_Dave_Num1_Id0_Clob0_Sell100BTC_Price101_GTB20.GetOrderId(),
								FillAmount:   10_000_000_000,
							},
						},
					),
				},
			},
			expectedMev:                          100_000_000, // $100
			expectedMidPrice:                     100_000_000,
			expectedValidatorNumFills:            1,
			expectedValidatorVolumeQuoteQuantums: big.NewInt(10_100_000_000),
			expectedProposerNumFills:             1,
			expectedProposerVolumeQuoteQuantums:  big.NewInt(10_100_000_000),
		},
		"Case 1 (liquidations) - capturing the spread": {
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num0_10_000USD,
				constants.Bob_Num0_10_000USD,
				// Carl_Num0's TNC is $100.
				// This means that closing 100 BTC short at $101 will not require any insurance fund payments.
				constants.Carl_Num0_100BTC_Short_10100USD,

				// Dave_Num1 is the subaccount controlled by the block proposer.
				constants.Dave_Num1_10_000USD,
			},
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc,
			},
			feeParams: constants.PerpetualFeeParamsNoFee,
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},
			restingOrders: []types.Order{
				constants.Order_Alice_Num0_Id0_Clob0_Sell200BTC_Price101_GTB20,
				constants.Order_Bob_Num0_Id1_Clob0_Buy100BTC_Price99_GTB20,
			},
			liquidationOrders: []types.LiquidationOrder{
				constants.LiquidationOrder_Carl_Num0_Clob0_Buy100BTC_Price101,
			},
			proposedOperations: &types.MsgProposedOperations{
				OperationsQueue: []types.OperationRaw{
					// X sells 100 at $101 instead of A.
					clobtest.NewShortTermOrderPlacementOperationRaw(constants.Order_Dave_Num1_Id0_Clob0_Sell100BTC_Price101_GTB20),
					clobtest.NewMatchOperationRaw(
						&constants.LiquidationOrder_Carl_Num0_Clob0_Buy100BTC_Price101,
						[]types.MakerFill{
							{
								MakerOrderId: constants.Order_Dave_Num1_Id0_Clob0_Sell100BTC_Price101_GTB20.GetOrderId(),
								FillAmount:   10_000_000_000,
							},
						},
					),
				},
			},
			expectedMev:                          100_000_000, // $100
			expectedMidPrice:                     100_000_000,
			expectedValidatorNumFills:            1,
			expectedValidatorVolumeQuoteQuantums: big.NewInt(10_100_000_000),
			expectedProposerNumFills:             1,
			expectedProposerVolumeQuoteQuantums:  big.NewInt(10_100_000_000),
		},
		// TODO(CLOB-742): re-enable deleveraging and funding in MEV calculation.
		// "Case 1 (positive funding) - capturing the spread": {
		// 	subaccounts: []satypes.Subaccount{
		// 		constants.Alice_Num0_10_000USD,
		// 		constants.Bob_Num0_10_000USD,
		// 		constants.Carl_Num0_10000USD,

		// 		// Dave_Num1 is the subaccount controlled by the block proposer.
		// 		constants.Dave_Num1_10_000USD,
		// 	},
		// 	clobPairs: []types.ClobPair{
		// 		constants.ClobPair_Btc,
		// 	},
		//  feeParams: constants.PerpetualFeeParamsNoFee,
		// 	perpetuals: []perptypes.Perpetual{
		// 		constants.BtcUsd_20PercentInitial_10PercentMaintenance,
		// 	},
		// 	restingOrders: []types.Order{
		// 		constants.Order_Alice_Num0_Id0_Clob0_Sell200BTC_Price101_GTB20,
		// 		constants.Order_Bob_Num0_Id1_Clob0_Buy100BTC_Price99_GTB20,
		// 	},
		// 	intrablockMsgs: []sdk.Msg{
		// 		types.NewMsgPlaceOrder(
		// 			constants.Order_Carl_Num0_Id1_Clob0_Buy100BTC_Price101_GTB20,
		// 		),
		// 	},
		// 	setupPerpetualKeeperMocks: func(perpKeeper *mocks.ProcessPerpetualKeeper) {
		// 		perpKeeper.On("MaybeProcessNewFundingTickEpoch", mock.Anything).Return()
		// 		perpKeeper.On("GetPerpetual", mock.Anything, mock.Anything).
		// 			Return(constants.BtcUsd_20PercentInitial_10PercentMaintenance, nil)
		// 		// Long positions pay $10.
		// 		perpKeeper.On("GetSettlement", mock.Anything, mock.Anything, big.NewInt(10_000_000_000), mock.Anything).
		// 			Return(big.NewInt(-10_000_000), new(big.Int), nil)
		// 		// Short positions receive $10.
		// 		perpKeeper.On("GetSettlement", mock.Anything, mock.Anything, big.NewInt(-10_000_000_000), mock.Anything).
		// 			Return(big.NewInt(10_000_000), new(big.Int), nil)
		// 	},
		// 	proposedOperations: &types.MsgProposedOperations{
		// 		OperationsQueue: []types.OperationRaw{
		// 			// X sells 100 at $101 instead of A.
		// 			clobtest.NewShortTermOrderPlacementOperationRaw(constants.Order_Dave_Num1_Id0_Clob0_Sell100BTC_Price101_GTB20),
		// 			clobtest.NewShortTermOrderPlacementOperationRaw(constants.Order_Carl_Num0_Id1_Clob0_Buy100BTC_Price101_GTB20),
		// 			clobtest.NewMatchOperationRaw(
		// 				&constants.Order_Carl_Num0_Id1_Clob0_Buy100BTC_Price101_GTB20,
		// 				[]types.MakerFill{
		// 					{
		// 						MakerOrderId: constants.Order_Dave_Num1_Id0_Clob0_Sell100BTC_Price101_GTB20.GetOrderId(),
		// 						FillAmount:   10_000_000_000,
		// 					},
		// 				},
		// 			),
		// 		},
		// 	},
		// 	// $100 from MEV, $10 from funding.
		// 	expectedMev:                          110_000_000, // $110
		// 	expectedMidPrice:                     100_000_000,
		// 	expectedValidatorNumFills:            1,
		// 	expectedValidatorVolumeQuoteQuantums: big.NewInt(10_100_000_000),
		// 	expectedProposerNumFills:             1,
		// 	expectedProposerVolumeQuoteQuantums:  big.NewInt(10_100_000_000),
		// },
		// "Case 1 (negative funding) - capturing the spread": {
		// 	subaccounts: []satypes.Subaccount{
		// 		constants.Alice_Num0_10_000USD,
		// 		constants.Bob_Num0_10_000USD,
		// 		constants.Carl_Num0_10000USD,

		// 		// Dave_Num1 is the subaccount controlled by the block proposer.
		// 		constants.Dave_Num1_10_000USD,
		// 	},
		// 	clobPairs: []types.ClobPair{
		// 		constants.ClobPair_Btc,
		// 	},
		//  feeParams: constants.PerpetualFeeParamsNoFee,
		// 	perpetuals: []perptypes.Perpetual{
		// 		constants.BtcUsd_20PercentInitial_10PercentMaintenance,
		// 	},
		// 	restingOrders: []types.Order{
		// 		constants.Order_Alice_Num0_Id0_Clob0_Sell200BTC_Price101_GTB20,
		// 		constants.Order_Bob_Num0_Id1_Clob0_Buy100BTC_Price99_GTB20,
		// 	},
		// 	intrablockMsgs: []sdk.Msg{
		// 		types.NewMsgPlaceOrder(
		// 			constants.Order_Carl_Num0_Id1_Clob0_Buy100BTC_Price101_GTB20,
		// 		),
		// 	},
		// 	setupPerpetualKeeperMocks: func(perpKeeper *mocks.ProcessPerpetualKeeper) {
		// 		perpKeeper.On("MaybeProcessNewFundingTickEpoch", mock.Anything).Return()
		// 		perpKeeper.On("GetPerpetual", mock.Anything, mock.Anything).
		// 			Return(constants.BtcUsd_20PercentInitial_10PercentMaintenance, nil)
		// 		// Long positions receive $10.
		// 		perpKeeper.On("GetSettlement", mock.Anything, mock.Anything, big.NewInt(10_000_000_000), mock.Anything).
		// 			Return(big.NewInt(10_000_000), new(big.Int), nil)
		// 		// Short positions pay $10.
		// 		perpKeeper.On("GetSettlement", mock.Anything, mock.Anything, big.NewInt(-10_000_000_000), mock.Anything).
		// 			Return(big.NewInt(-10_000_000), new(big.Int), nil)
		// 	},
		// 	proposedOperations: &types.MsgProposedOperations{
		// 		OperationsQueue: []types.OperationRaw{
		// 			// X sells 100 at $101 instead of A.
		// 			clobtest.NewShortTermOrderPlacementOperationRaw(constants.Order_Dave_Num1_Id0_Clob0_Sell100BTC_Price101_GTB20),
		// 			clobtest.NewShortTermOrderPlacementOperationRaw(constants.Order_Carl_Num0_Id1_Clob0_Buy100BTC_Price101_GTB20),
		// 			clobtest.NewMatchOperationRaw(
		// 				&constants.Order_Carl_Num0_Id1_Clob0_Buy100BTC_Price101_GTB20,
		// 				[]types.MakerFill{
		// 					{
		// 						MakerOrderId: constants.Order_Dave_Num1_Id0_Clob0_Sell100BTC_Price101_GTB20.GetOrderId(),
		// 						FillAmount:   10_000_000_000,
		// 					},
		// 				},
		// 			),
		// 		},
		// 	},
		// 	// $100 from MEV, minus $10 from funding.
		// 	expectedMev:                          90_000_000, // $90
		// 	expectedMidPrice:                     100_000_000,
		// 	expectedValidatorNumFills:            1,
		// 	expectedValidatorVolumeQuoteQuantums: big.NewInt(10_100_000_000),
		// 	expectedProposerNumFills:             1,
		// 	expectedProposerVolumeQuoteQuantums:  big.NewInt(10_100_000_000),
		// },
		"Case 2 - capturing utility from excess limit price": {
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num0_10_000USD,
				constants.Bob_Num0_10_000USD,
				constants.Carl_Num0_10000USD,
				constants.Dave_Num0_10000USD,

				// Dave_Num1 is the subaccount controlled by the block proposer.
				constants.Dave_Num1_10_000USD,
			},
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc,
			},
			feeParams: constants.PerpetualFeeParamsNoFee,
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},
			restingOrders: []types.Order{
				constants.Order_Alice_Num0_Id0_Clob0_Sell100BTC_Price102_GTB20,
				constants.Order_Bob_Num0_Id0_Clob0_Sell200BTC_Price101_GTB20,
				constants.Order_Carl_Num0_Id0_Clob0_Buy100BTC_Price99_GTB20,
			},
			intrablockMsgs: []sdk.Msg{
				types.NewMsgPlaceOrder(
					constants.Order_Dave_Num0_Id0_Clob0_Buy100BTC_Price102_GTB20,
				),
			},
			proposedOperations: &types.MsgProposedOperations{
				OperationsQueue: []types.OperationRaw{
					// X sells 100 at $102 to D.
					clobtest.NewShortTermOrderPlacementOperationRaw(constants.Order_Dave_Num1_Id0_Clob0_Sell100BTC_Price102_GTB20),
					clobtest.NewShortTermOrderPlacementOperationRaw(constants.Order_Dave_Num0_Id0_Clob0_Buy100BTC_Price102_GTB20),
					clobtest.NewMatchOperationRaw(
						&constants.Order_Dave_Num0_Id0_Clob0_Buy100BTC_Price102_GTB20,
						[]types.MakerFill{
							{
								MakerOrderId: constants.Order_Dave_Num1_Id0_Clob0_Sell100BTC_Price102_GTB20.GetOrderId(),
								FillAmount:   10_000_000_000,
							},
						},
					),
					// B sells 100 at $101 to X.
					clobtest.NewShortTermOrderPlacementOperationRaw(constants.Order_Dave_Num1_Id0_Clob0_Buy100BTC_Price101_GTB20),
					clobtest.NewShortTermOrderPlacementOperationRaw(constants.Order_Bob_Num0_Id0_Clob0_Sell200BTC_Price101_GTB20),
					clobtest.NewMatchOperationRaw(
						&constants.Order_Bob_Num0_Id0_Clob0_Sell200BTC_Price101_GTB20,
						[]types.MakerFill{
							{
								MakerOrderId: constants.Order_Dave_Num1_Id0_Clob0_Buy100BTC_Price101_GTB20.GetOrderId(),
								FillAmount:   10_000_000_000,
							},
						},
					),
				},
			},
			expectedMev:                          100_000_000, // $100
			expectedMidPrice:                     100_000_000,
			expectedValidatorNumFills:            1,
			expectedValidatorVolumeQuoteQuantums: big.NewInt(10_100_000_000),
			expectedProposerNumFills:             2,
			expectedProposerVolumeQuoteQuantums:  big.NewInt(20_300_000_000),
		},
		"Case 2 (liquidations) - capturing utility from excess limit price": {
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num0_10_000USD,
				constants.Bob_Num0_10_000USD,
				constants.Carl_Num0_10000USD,
				// Dave_Num0's TNC is $200.
				// This means that closing 100 BTC short at $102 will not require any insurance fund payments.
				constants.Dave_Num0_100BTC_Short_10200USD,

				// Dave_Num1 is the subaccount controlled by the block proposer.
				constants.Dave_Num1_10_000USD,
			},
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc,
			},
			feeParams: constants.PerpetualFeeParamsNoFee,
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},
			restingOrders: []types.Order{
				constants.Order_Alice_Num0_Id0_Clob0_Sell100BTC_Price102_GTB20,
				constants.Order_Bob_Num0_Id0_Clob0_Sell200BTC_Price101_GTB20,
				constants.Order_Carl_Num0_Id0_Clob0_Buy100BTC_Price99_GTB20,
			},
			liquidationOrders: []types.LiquidationOrder{
				constants.LiquidationOrder_Dave_Num0_Clob0_Buy100BTC_Price102,
			},
			proposedOperations: &types.MsgProposedOperations{
				OperationsQueue: []types.OperationRaw{
					// X sells 100 at $102 to D.
					clobtest.NewShortTermOrderPlacementOperationRaw(constants.Order_Dave_Num1_Id0_Clob0_Sell100BTC_Price102_GTB20),
					clobtest.NewMatchOperationRaw(
						&constants.LiquidationOrder_Dave_Num0_Clob0_Buy100BTC_Price102,
						[]types.MakerFill{
							{
								MakerOrderId: constants.Order_Dave_Num1_Id0_Clob0_Sell100BTC_Price102_GTB20.GetOrderId(),
								FillAmount:   10_000_000_000,
							},
						},
					),
					// B sells 100 at $101 to X.
					clobtest.NewShortTermOrderPlacementOperationRaw(constants.Order_Dave_Num1_Id0_Clob0_Buy100BTC_Price101_GTB20),
					clobtest.NewShortTermOrderPlacementOperationRaw(constants.Order_Bob_Num0_Id0_Clob0_Sell200BTC_Price101_GTB20),
					clobtest.NewMatchOperationRaw(
						&constants.Order_Bob_Num0_Id0_Clob0_Sell200BTC_Price101_GTB20,
						[]types.MakerFill{
							{
								MakerOrderId: constants.Order_Dave_Num1_Id0_Clob0_Buy100BTC_Price101_GTB20.GetOrderId(),
								FillAmount:   10_000_000_000,
							},
						},
					),
				},
			},
			// For the local validator, Dave_Num0's liquidation order filled at $101, which is better than Dave_Num0's
			// bankruptcy price. This means that Dave_Num0 will need to pay the insurance fund a fee.
			// In this case, liquidation fee is 0.5% of the liquidation order's notional value,
			// which is 0.5% * 100 * $101 = $50.5.
			// Using the MEV formula, we can calculate the expected MEV to be $74.75.
			expectedMev:                          74_750_000, // $74.75
			expectedMidPrice:                     100_000_000,
			expectedValidatorNumFills:            1,
			expectedValidatorVolumeQuoteQuantums: big.NewInt(10_100_000_000),
			expectedProposerNumFills:             2,
			expectedProposerVolumeQuoteQuantums:  big.NewInt(20_300_000_000),
		},
		"Case 3 - stale order snipe, cancel ignored": {
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num0_10_000USD,
				constants.Bob_Num0_10_000USD,
				constants.Carl_Num0_10000USD,
				constants.Dave_Num0_10000USD,

				// Dave_Num1 is the subaccount controlled by the block proposer.
				constants.Dave_Num1_10_000USD,
			},
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc,
			},
			feeParams: constants.PerpetualFeeParamsNoFee,
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},
			restingOrders: []types.Order{
				constants.Order_Alice_Num0_Id0_Clob0_Sell100BTC_Price106_GTB20,
				constants.Order_Bob_Num0_Id0_Clob0_Sell100BTC_Price101_GTB20,
				constants.Order_Carl_Num0_Id0_Clob0_Buy100BTC_Price99_GTB20,
			},
			intrablockMsgs: []sdk.Msg{
				types.NewMsgCancelOrderShortTerm(
					constants.Order_Bob_Num0_Id0_Clob0_Sell100BTC_Price101_GTB20.OrderId, // Cancellation.
					20,
				),
				types.NewMsgPlaceOrder(
					constants.Order_Dave_Num0_Id1_Clob0_Buy100BTC_Price104_GTB20,
				),
			},
			proposedOperations: &types.MsgProposedOperations{
				OperationsQueue: []types.OperationRaw{
					// B sells 100 at $101 to X.
					clobtest.NewShortTermOrderPlacementOperationRaw(constants.Order_Dave_Num1_Id0_Clob0_Buy100BTC_Price101_GTB20),
					clobtest.NewShortTermOrderPlacementOperationRaw(constants.Order_Bob_Num0_Id0_Clob0_Sell100BTC_Price101_GTB20),
					clobtest.NewMatchOperationRaw(
						&constants.Order_Bob_Num0_Id0_Clob0_Sell100BTC_Price101_GTB20,
						[]types.MakerFill{
							{
								MakerOrderId: constants.Order_Dave_Num1_Id0_Clob0_Buy100BTC_Price101_GTB20.GetOrderId(),
								FillAmount:   10_000_000_000,
							},
						},
					),
				},
			},
			expectedMev:                          400_000_000, // $400
			expectedMidPrice:                     105_000_000,
			expectedValidatorNumFills:            0,
			expectedValidatorVolumeQuoteQuantums: big.NewInt(0),
			expectedProposerNumFills:             1,
			expectedProposerVolumeQuoteQuantums:  big.NewInt(10_100_000_000),
		},
		"Case 4 - front running a stale order snipe": {
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num0_10_000USD,
				constants.Bob_Num0_10_000USD,
				constants.Carl_Num0_10000USD,
				constants.Dave_Num0_10000USD,

				// Dave_Num1 is the subaccount controlled by the block proposer.
				constants.Dave_Num1_10_000USD,
			},
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc,
			},
			feeParams: constants.PerpetualFeeParamsNoFee,
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},
			restingOrders: []types.Order{
				constants.Order_Alice_Num0_Id0_Clob0_Sell100BTC_Price106_GTB20,
				constants.Order_Bob_Num0_Id0_Clob0_Sell100BTC_Price101_GTB20,
				constants.Order_Carl_Num0_Id0_Clob0_Buy100BTC_Price99_GTB20,
			},
			intrablockMsgs: []sdk.Msg{
				types.NewMsgPlaceOrder(
					constants.Order_Dave_Num0_Id0_Clob0_Buy100BTC_Price101_GTB20,
				),
				types.NewMsgPlaceOrder(
					constants.Order_Dave_Num0_Id1_Clob0_Buy100BTC_Price104_GTB20,
				),
			},
			proposedOperations: &types.MsgProposedOperations{
				OperationsQueue: []types.OperationRaw{
					// B sells 100 at $101 to X.
					clobtest.NewShortTermOrderPlacementOperationRaw(constants.Order_Dave_Num1_Id0_Clob0_Buy100BTC_Price101_GTB20),
					clobtest.NewShortTermOrderPlacementOperationRaw(constants.Order_Bob_Num0_Id0_Clob0_Sell100BTC_Price101_GTB20),
					clobtest.NewMatchOperationRaw(
						&constants.Order_Bob_Num0_Id0_Clob0_Sell100BTC_Price101_GTB20,
						[]types.MakerFill{
							{
								MakerOrderId: constants.Order_Dave_Num1_Id0_Clob0_Buy100BTC_Price101_GTB20.GetOrderId(),
								FillAmount:   10_000_000_000,
							},
						},
					),
				},
			},
			expectedMev:                          400_000_000, // $400
			expectedMidPrice:                     105_000_000,
			expectedValidatorNumFills:            1,
			expectedValidatorVolumeQuoteQuantums: big.NewInt(10_100_000_000),
			expectedProposerNumFills:             1,
			expectedProposerVolumeQuoteQuantums:  big.NewInt(10_100_000_000),
		},
		"Case 4 (liquidations) - front running a stale order snipe": {
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num0_10_000USD,
				constants.Bob_Num0_10_000USD,
				// Carl_Num0's TNC is $100.
				// This means that closing 100 BTC short at $101 will not require any insurance fund payments.
				constants.Carl_Num0_100BTC_Short_10100USD,
				constants.Dave_Num0_10000USD,

				// Dave_Num1 is the subaccount controlled by the block proposer.
				constants.Dave_Num1_10_000USD,
			},
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc,
			},
			feeParams: constants.PerpetualFeeParamsNoFee,
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},
			restingOrders: []types.Order{
				// Asks.
				constants.Order_Alice_Num0_Id0_Clob0_Sell100BTC_Price106_GTB20,
				constants.Order_Bob_Num0_Id0_Clob0_Sell100BTC_Price101_GTB20,
				// Bids
				constants.Order_Bob_Num0_Id1_Clob0_Buy100BTC_Price99_GTB20,
			},
			liquidationOrders: []types.LiquidationOrder{
				constants.LiquidationOrder_Carl_Num0_Clob0_Buy100BTC_Price101,
			},
			intrablockMsgs: []sdk.Msg{
				types.NewMsgPlaceOrder(
					constants.Order_Dave_Num0_Id1_Clob0_Buy100BTC_Price104_GTB20,
				),
			},
			proposedOperations: &types.MsgProposedOperations{
				OperationsQueue: []types.OperationRaw{
					// B sells 100 at $101 to X.
					clobtest.NewShortTermOrderPlacementOperationRaw(constants.Order_Dave_Num1_Id0_Clob0_Buy100BTC_Price101_GTB20),
					clobtest.NewShortTermOrderPlacementOperationRaw(constants.Order_Bob_Num0_Id0_Clob0_Sell100BTC_Price101_GTB20),
					clobtest.NewMatchOperationRaw(
						&constants.Order_Bob_Num0_Id0_Clob0_Sell100BTC_Price101_GTB20,
						[]types.MakerFill{
							{
								MakerOrderId: constants.Order_Dave_Num1_Id0_Clob0_Buy100BTC_Price101_GTB20.GetOrderId(),
								FillAmount:   10_000_000_000,
							},
						},
					),
				},
			},
			expectedMev:                          400_000_000, // $400
			expectedMidPrice:                     105_000_000,
			expectedValidatorNumFills:            1,
			expectedValidatorVolumeQuoteQuantums: big.NewInt(10_100_000_000),
			expectedProposerNumFills:             1,
			expectedProposerVolumeQuoteQuantums:  big.NewInt(10_100_000_000),
		},
		"Case 5 - defending a stale order": {
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num0_10_000USD,
				constants.Bob_Num0_10_000USD,
				constants.Carl_Num0_10000USD,
				constants.Dave_Num0_10000USD,

				// Dave_Num1 is the subaccount controlled by the block proposer.
				constants.Dave_Num1_10_000USD,
			},
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc,
			},
			feeParams: constants.PerpetualFeeParamsNoFee,
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},
			restingOrders: []types.Order{
				constants.Order_Alice_Num0_Id0_Clob0_Sell100BTC_Price106_GTB20,
				constants.Order_Dave_Num1_Id0_Clob0_Sell100BTC_Price101_GTB20,
				constants.Order_Carl_Num0_Id0_Clob0_Buy100BTC_Price99_GTB20,
			},
			intrablockMsgs: []sdk.Msg{
				types.NewMsgPlaceOrder(
					constants.Order_Carl_Num0_Id1_Clob0_Buy100BTC_Price101_GTB20,
				),
				types.NewMsgPlaceOrder(
					constants.Order_Dave_Num0_Id1_Clob0_Buy100BTC_Price104_GTB20,
				),
			},
			proposedOperations: &types.MsgProposedOperations{
				OperationsQueue: []types.OperationRaw{
					// Empty.
				},
			},
			expectedMev:                          400_000_000, // $400
			expectedMidPrice:                     105_000_000,
			expectedValidatorNumFills:            1,
			expectedValidatorVolumeQuoteQuantums: big.NewInt(10_100_000_000),
			expectedProposerNumFills:             0,
			expectedProposerVolumeQuoteQuantums:  big.NewInt(0),
		},
		"Case 5 (liquidations) - defending a stale order": {
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num0_10_000USD,
				constants.Bob_Num0_10_000USD,
				constants.Carl_Num0_100BTC_Short_10100USD,
				constants.Dave_Num0_10000USD,

				// Dave_Num1 is the subaccount controlled by the block proposer.
				constants.Dave_Num1_10_000USD,
			},
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc,
			},
			feeParams: constants.PerpetualFeeParamsNoFee,
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},
			restingOrders: []types.Order{
				constants.Order_Alice_Num0_Id0_Clob0_Sell100BTC_Price106_GTB20,
				constants.Order_Dave_Num1_Id0_Clob0_Sell100BTC_Price101_GTB20,
				constants.Order_Bob_Num0_Id1_Clob0_Buy100BTC_Price99_GTB20,
			},
			liquidationOrders: []types.LiquidationOrder{
				constants.LiquidationOrder_Carl_Num0_Clob0_Buy100BTC_Price101,
			},
			intrablockMsgs: []sdk.Msg{
				types.NewMsgPlaceOrder(
					constants.Order_Dave_Num0_Id1_Clob0_Buy100BTC_Price104_GTB20,
				),
			},
			proposedOperations: &types.MsgProposedOperations{
				OperationsQueue: []types.OperationRaw{
					// Empty.
				},
			},
			expectedMev:                          400_000_000, // $400
			expectedMidPrice:                     105_000_000,
			expectedValidatorNumFills:            1,
			expectedValidatorVolumeQuoteQuantums: big.NewInt(10_100_000_000),
			expectedProposerNumFills:             0,
			expectedProposerVolumeQuoteQuantums:  big.NewInt(0),
		},
		"Case 6 - maker and taker switch": {
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num0_10_000USD,
				constants.Bob_Num0_10_000USD,
				constants.Carl_Num0_10000USD,
				constants.Dave_Num0_10000USD,

				// Dave_Num1 is the subaccount controlled by the block proposer.
				constants.Dave_Num1_10_000USD,
			},
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc, // Non zero fee.
			},
			feeParams: constants.PerpetualFeeParams,
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},
			restingOrders: []types.Order{
				constants.Order_Alice_Num0_Id0_Clob0_Sell100BTC_Price102_GTB20,
				constants.Order_Bob_Num0_Id0_Clob0_Sell100BTC_Price101_GTB20,
				constants.Order_Carl_Num0_Id0_Clob0_Buy100BTC_Price99_GTB20,
			},
			intrablockMsgs: []sdk.Msg{
				types.NewMsgPlaceOrder(
					constants.Order_Dave_Num1_Id0_Clob0_Buy100BTC_Price101_GTB20,
				),
				types.NewMsgPlaceOrder(
					constants.Order_Carl_Num0_Id1_Clob0_Buy100BTC_Price100_GTB20, // Make mid price $101.
				),
			},
			proposedOperations: &types.MsgProposedOperations{
				OperationsQueue: []types.OperationRaw{
					// B sells 100 at $101 to X.
					clobtest.NewShortTermOrderPlacementOperationRaw(constants.Order_Dave_Num1_Id0_Clob0_Buy100BTC_Price101_GTB20),
					clobtest.NewShortTermOrderPlacementOperationRaw(constants.Order_Bob_Num0_Id0_Clob0_Sell100BTC_Price101_GTB20),
					clobtest.NewMatchOperationRaw(
						&constants.Order_Bob_Num0_Id0_Clob0_Sell100BTC_Price101_GTB20,
						[]types.MakerFill{
							{
								MakerOrderId: constants.Order_Dave_Num1_Id0_Clob0_Buy100BTC_Price101_GTB20.GetOrderId(),
								FillAmount:   10_000_000_000,
							},
						},
					),
				},
			},
			expectedMev:                          3_030_000, // $3.03
			expectedMidPrice:                     101_000_000,
			expectedValidatorNumFills:            1,
			expectedValidatorVolumeQuoteQuantums: big.NewInt(10_100_000_000),
			expectedProposerNumFills:             1,
			expectedProposerVolumeQuoteQuantums:  big.NewInt(10_100_000_000),
		},
		"Case 7 (liquidations) - ignores liquidation": {
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num0_10_000USD,
				constants.Bob_Num0_10_000USD,
				constants.Carl_Num0_10000USD,
				constants.Dave_Num0_10000USD,

				// Dave_Num1 is the subaccount controlled by the block proposer.
				// Here Dave_Num1 is supposed to get liquidated.
				constants.Dave_Num1_100BTC_Short_10100USD,
			},
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc,
			},
			feeParams: constants.PerpetualFeeParamsNoFee,
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},
			restingOrders: []types.Order{
				constants.Order_Bob_Num0_Id0_Clob0_Sell200BTC_Price101_GTB20,
				constants.Order_Bob_Num0_Id1_Clob0_Buy100BTC_Price99_GTB20,
			},
			liquidationOrders: []types.LiquidationOrder{
				constants.LiquidationOrder_Dave_Num1_Clob0_Buy100BTC_Price101,
			},
			proposedOperations: &types.MsgProposedOperations{
				OperationsQueue: []types.OperationRaw{
					// Empty.
				},
			},
			expectedMev:                          100_000_000, // $100
			expectedMidPrice:                     100_000_000,
			expectedValidatorNumFills:            1,
			expectedValidatorVolumeQuoteQuantums: big.NewInt(10_100_000_000),
			expectedProposerNumFills:             0,
			expectedProposerVolumeQuoteQuantums:  big.NewInt(0),
		},
		// TODO(CLOB-742): re-enable deleveraging and funding in MEV calculation.
		// "Case 8 (deleveraging) - ignores deleveraging": {
		// 	subaccounts: []satypes.Subaccount{
		// 		constants.Alice_Num0_10_000USD,
		// 		constants.Bob_Num0_10_000USD,
		// 		constants.Carl_Num0_10000USD,
		// 		constants.Dave_Num0_100BTC_Long_9900USD_Short,

		// 		// Dave_Num1 is the subaccount controlled by the block proposer.
		// 		// Here Dave_Num1 is supposed to get deleveraged.
		// 		constants.Dave_Num1_100BTC_Short_10100USD,
		// 	},
		// 	clobPairs: []types.ClobPair{
		// 		constants.ClobPair_Btc,
		// 	},
		//  feeParams: constants.PerpetualFeeParamsNoFee,
		// 	perpetuals: []perptypes.Perpetual{
		// 		constants.BtcUsd_20PercentInitial_10PercentMaintenance,
		// 	},
		// 	restingOrders: []types.Order{
		// 		constants.Order_Alice_Num0_Id0_Clob0_Sell100BTC_Price102_GTB20,
		// 		constants.Order_Bob_Num0_Id1_Clob0_Buy100BTC_Price98_GTB20,
		// 	},
		// 	liquidationOrders: []types.LiquidationOrder{
		// 		constants.LiquidationOrder_Dave_Num1_Clob0_Buy100BTC_Price102,
		// 	},
		// 	proposedOperations: &types.MsgProposedOperations{
		// 		OperationsQueue: []types.OperationRaw{
		// 			// Empty.
		// 		},
		// 	},
		// 	expectedMev:                          100_000_000, // $100
		// 	expectedMidPrice:                     100_000_000,
		// 	expectedValidatorNumFills:            1,
		// 	expectedValidatorVolumeQuoteQuantums: big.NewInt(10_100_000_000),
		// 	expectedProposerNumFills:             0,
		// 	expectedProposerVolumeQuoteQuantums:  big.NewInt(0),
		// },
		// "Case 9 (deleveraging) - deleverage against different subaccount": {
		// 	subaccounts: []satypes.Subaccount{
		// 		constants.Alice_Num0_10_000USD,
		// 		constants.Bob_Num0_10_000USD,
		// 		constants.Carl_Num0_100BTC_Short_10100USD,
		// 		// Dave_Num0 is supposed to get deleveraged against Carl_Num0.
		// 		constants.Dave_Num0_100BTC_Long_9900USD_Short,

		// 		// Dave_Num1 is the subaccount controlled by the block proposer.
		// 		constants.Dave_Num1_100BTC_Short_10100USD,
		// 	},
		// 	clobPairs: []types.ClobPair{
		// 		constants.ClobPair_Btc,
		// 	},
		//  feeParams: constants.PerpetualFeeParamsNoFee,
		// 	perpetuals: []perptypes.Perpetual{
		// 		constants.BtcUsd_20PercentInitial_10PercentMaintenance,
		// 	},
		// 	restingOrders: []types.Order{
		// 		constants.Order_Alice_Num0_Id0_Clob0_Sell100BTC_Price102_GTB20,
		// 		constants.Order_Bob_Num0_Id1_Clob0_Buy100BTC_Price98_GTB20,
		// 	},
		// 	liquidationOrders: []types.LiquidationOrder{
		// 		constants.LiquidationOrder_Dave_Num0_Clob0_Sell100BTC_Price98,
		// 	},
		// 	proposedOperations: &types.MsgProposedOperations{
		// 		OperationsQueue: []types.OperationRaw{
		// 			clobtest.NewMatchOperationRawFromPerpetualDeleveragingLiquidation(
		// 				types.MatchPerpetualDeleveraging{
		// 					Liquidated:  constants.Dave_Num0,
		// 					PerpetualId: constants.ClobPair_Btc.MustGetPerpetualId(),
		// 					Fills: []types.MatchPerpetualDeleveraging_Fill{
		// 						{
		// 							// Deleverage against Dave_Num1 instead of Carl_Num0.
		// 							OffsettingSubaccountId: constants.Dave_Num1,
		// 							FillAmount:             10_000_000_000,
		// 						},
		// 					},
		// 				},
		// 			),
		// 		},
		// 	},
		// 	expectedMev:                          100_000_000, // $100
		// 	expectedMidPrice:                     100_000_000,
		// 	expectedValidatorNumFills:            1,
		// 	expectedValidatorVolumeQuoteQuantums: big.NewInt(9_900_000_000),
		// 	expectedProposerNumFills:             1,
		// 	expectedProposerVolumeQuoteQuantums:  big.NewInt(9_900_000_000),
		// },
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup keeper state.
			memClob := memclob.NewMemClobPriceTimePriority(false)
			mockBankKeeper := &mocks.BankKeeper{}
			mockBankKeeper.On(
				"SendCoinsFromModuleToModule",
				mock.Anything,
				mock.Anything,
				mock.Anything,
				mock.Anything,
			).Return(nil)
			mockBankKeeper.On(
				"SendCoins",
				mock.Anything,
				mock.Anything,
				mock.Anything,
				mock.Anything,
			).Return(nil)
			mockBankKeeper.On(
				"GetBalance",
				mock.Anything,
				perptypes.InsuranceFundModuleAddress,
				constants.Usdc.Denom,
			).Return(
				sdk.NewCoin(constants.Usdc.Denom, sdkmath.NewIntFromBigInt(new(big.Int))),
			)

			ks := keepertest.NewClobKeepersTestContext(
				t,
				memClob,
				mockBankKeeper,
				indexer_manager.NewIndexerEventManagerNoop(),
			)
			ctx := ks.Ctx.WithIsCheckTx(true).WithBlockTime(time.Unix(5, 0))

			// Create the default markets.
			keepertest.CreateTestMarkets(t, ctx, ks.PricesKeeper)
			err := ks.PricesKeeper.UpdateMarketPrices(
				ctx,
				[]*pricestypes.MsgUpdateMarketPrices_MarketPrice{
					{
						MarketId: 0,
						Price:    10_000_000, // $100
					},
				},
			)
			require.NoError(t, err)

			// Create liquidity tiers.
			keepertest.CreateTestLiquidityTiers(t, ctx, ks.PerpetualsKeeper)

			require.NoError(t, ks.FeeTiersKeeper.SetPerpetualFeeParams(ctx, tc.feeParams))

			// Set up USDC asset in assets module.
			err = keepertest.CreateUsdcAsset(ctx, ks.AssetsKeeper)
			require.NoError(t, err)

			// Create all perpetuals.
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
				)
				require.NoError(t, err)
			}

			// Create all subaccounts.
			for _, subaccount := range tc.subaccounts {
				ks.SubaccountsKeeper.SetSubaccount(ctx, subaccount)
			}

			// Create all CLOBs.
			for _, clobPair := range tc.clobPairs {
				_, err := ks.ClobKeeper.CreatePerpetualClobPairAndMemStructs(
					ctx,
					clobPair.Id,
					clobtest.MustPerpetualId(clobPair),
					satypes.BaseQuantums(clobPair.StepBaseQuantums),
					clobPair.QuantumConversionExponent,
					clobPair.SubticksPerTick,
					clobPair.Status,
				)
				require.NoError(t, err)
			}

			require.NoError(
				t,
				ks.ClobKeeper.InitializeLiquidationsConfig(ctx, constants.LiquidationsConfig_No_Limit),
			)

			// Create all existing orders.
			for _, order := range tc.restingOrders {
				// set bytes so operations proposed contains correct bytes
				shortTermOrderPlacement := clobtest.NewShortTermOrderPlacementOperationRaw(order)
				bytes := shortTermOrderPlacement.GetShortTermOrderPlacement()
				tempCtx := ctx.WithTxBytes(bytes)
				_, _, err := ks.ClobKeeper.PlaceShortTermOrder(tempCtx, &types.MsgPlaceOrder{Order: order})
				require.NoError(t, err)
			}

			// Use a cached context to simulate checkState.
			cachedCtx, _ := ctx.CacheContext()

			// Place the liquidation orders.
			for _, liquidationOrder := range tc.liquidationOrders {
				_, _, err := ks.ClobKeeper.PlacePerpetualLiquidation(
					cachedCtx,
					liquidationOrder,
				)
				require.NoError(t, err)
			}

			// Execute all intrablock messages.
			for _, msg := range tc.intrablockMsgs {
				if placeOrderMsg, ok := msg.(*types.MsgPlaceOrder); ok {
					// set bytes so operations proposed contains correct bytes
					shortTermOrderPlacement := clobtest.NewShortTermOrderPlacementOperationRaw(placeOrderMsg.Order)
					bytes := shortTermOrderPlacement.GetShortTermOrderPlacement()
					tempCtx := ctx.WithTxBytes(bytes)
					_, _, err := ks.ClobKeeper.PlaceShortTermOrder(tempCtx, placeOrderMsg)
					require.NoError(t, err)
				}
				if cancelOrderMsg, ok := msg.(*types.MsgCancelOrder); ok {
					err := ks.ClobKeeper.CancelShortTermOrder(
						cachedCtx,
						cancelOrderMsg,
					)
					require.NoError(t, err)
				}
			}

			// Run the test.
			ctx = ctx.WithValue(process.ConsensusRound, int64(0))
			ctx = ctx.WithProposer(constants.AliceConsAddress)
			aliceValidator, err := stakingtypes.NewValidator(
				constants.AliceValAddress.String(),
				constants.AlicePrivateKey.PubKey(),
				stakingtypes.Description{
					Moniker: "alice",
				},
			)
			require.NoError(t, err)
			mockStakingKeeper := &mocks.ProcessStakingKeeper{}
			mockStakingKeeper.On("GetValidatorByConsAddr", mock.Anything, constants.AliceConsAddress).
				Return(aliceValidator, nil)

			mockPerpetualKeeper := &mocks.ProcessPerpetualKeeper{}
			if tc.setupPerpetualKeeperMocks != nil {
				tc.setupPerpetualKeeperMocks(mockPerpetualKeeper)
			} else {
				mockPerpetualKeeper.On("MaybeProcessNewFundingTickEpoch", mock.Anything).Return()
				mockPerpetualKeeper.On("GetSettlementPpm", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(new(big.Int), new(big.Int), nil)
				for _, p := range tc.perpetuals {
					mockPerpetualKeeper.On("GetPerpetual", mock.Anything, p.Params.Id).Return(p, nil)
				}
			}

			mockLogger := &mocks.Logger{}
			mockLogger.On("With", mock.Anything, mock.Anything, mock.Anything, mock.Anything,
				mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(mockLogger)
			mockLogger.On(
				"Info",
				"Measuring MEV for proposed matches",
				metrics.Mev,
				tc.expectedMev,
				// Common metadata.
				metrics.BlockHeight,
				ctx.BlockHeight(),
				metrics.ConsensusRound,
				int64(0),
				metrics.Proposer,
				"alice",
				metrics.ClobPairId,
				uint32(0),
				metrics.MidPrice,
				tc.expectedMidPrice,
				metrics.OraclePrice,
				uint64(100_000_000),
				metrics.BestBid,
				mock.Anything,
				metrics.BestAsk,
				mock.Anything,
				// Validator stats.
				metrics.ValidatorNumFills,
				tc.expectedValidatorNumFills,
				metrics.ValidatorVolumeQuoteQuantums,
				tc.expectedValidatorVolumeQuoteQuantums.String(),
				// Proposer stats.
				metrics.ProposerNumFills,
				tc.expectedProposerNumFills,
				metrics.ProposerVolumeQuoteQuantums,
				tc.expectedProposerVolumeQuoteQuantums.String(),
			).Once()

			ctx = ctx.WithLogger(mockLogger)
			ks.ClobKeeper.RecordMevMetrics(ctx, mockStakingKeeper, mockPerpetualKeeper, tc.proposedOperations)

			mockLogger.AssertExpectations(t)
		})
	}
}

func TestGetMidPrices(t *testing.T) {
	// Set the maximum spread to 1%.
	keeper.MAX_SPREAD_BEFORE_FALLING_BACK_TO_ORACLE = new(big.Rat).SetFrac64(1, 100)

	tests := map[string]struct {
		// Setup.
		perpetuals  []perptypes.Perpetual
		subaccounts []satypes.Subaccount
		clobPairs   []types.ClobPair
		orders      []types.Order

		// Expectations.
		expectedMidPrices map[types.ClobPairId]types.Subticks
	}{
		"can get mid prices from one orderbook": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_100000USD,
				constants.Dave_Num0_500000USD,
			},
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc,
			},
			orders: []types.Order{
				// Bid
				constants.Order_Carl_Num0_Id0_Clob0_Buy1BTC_Price49800_GTB10,
				// Ask
				constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10,
			},

			expectedMidPrices: map[types.ClobPairId]types.Subticks{
				constants.ClobPair_Btc.GetClobPairId(): 49_900_000_000, // 49800 + (50000 - 49800) / 2
			},
		},
		"can get mid prices from one orderbook when there are multiple orders on the same level": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_100000USD,
				constants.Dave_Num0_500000USD,
			},
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc,
			},
			orders: []types.Order{
				// Bid
				constants.Order_Carl_Num0_Id0_Clob0_Buy1BTC_Price49800_GTB10,
				constants.Order_Carl_Num0_Id3_Clob0_Buy025BTC_Price49800,
				// Ask
				constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10,
				constants.Order_Dave_Num0_Id1_Clob0_Sell025BTC_Price50000_GTB11,
			},

			expectedMidPrices: map[types.ClobPairId]types.Subticks{
				constants.ClobPair_Btc.GetClobPairId(): 49_900_000_000, // 49800 + (50000 - 49800) / 2
			},
		},
		"can get mid prices from one orderbook when there are multiple price levels": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_100000USD,
				constants.Dave_Num0_500000USD,
			},
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc,
			},
			orders: []types.Order{
				// Bid
				constants.Order_Carl_Num0_Id4_Clob0_Buy05BTC_Price40000,
				constants.Order_Carl_Num0_Id0_Clob0_Buy1BTC_Price49800_GTB10,
				// Ask
				constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10,
				constants.Order_Dave_Num0_Id1_Clob0_Sell025BTC_Price50500_GTB11,
			},

			expectedMidPrices: map[types.ClobPairId]types.Subticks{
				constants.ClobPair_Btc.GetClobPairId(): 49_900_000_000, // 49800 + (50000 - 49800) / 2
			},
		},
		"can get mid prices from multiple orderbooks": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
				constants.EthUsd_20PercentInitial_10PercentMaintenance,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_100000USD,
				constants.Dave_Num0_500000USD,
			},
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc,
				constants.ClobPair_Eth,
			},
			orders: []types.Order{
				// Bid
				constants.Order_Carl_Num0_Id0_Clob0_Buy1BTC_Price49800_GTB10,
				constants.Order_Carl_Num0_Id4_Clob1_Buy01ETH_Price3000,
				// Ask
				constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000,
				constants.Order_Dave_Num0_Id4_Clob1_Sell1ETH_Price3020,
			},

			expectedMidPrices: map[types.ClobPairId]types.Subticks{
				constants.ClobPair_Btc.GetClobPairId(): 49_900_000_000, // 49800 + (50000 - 49800) / 2
				constants.ClobPair_Eth.GetClobPairId(): 3_010_000_000,  // 3000 + (3020 - 3000) / 2
			},
		},
		"fallback to oracle price for orderbooks that are empty": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
				constants.EthUsd_20PercentInitial_10PercentMaintenance,
			},
			subaccounts: []satypes.Subaccount{},
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc,
			},
			orders: []types.Order{},

			expectedMidPrices: map[types.ClobPairId]types.Subticks{
				constants.ClobPair_Btc.GetClobPairId(): 50_000_000_000,
			},
		},
		"fallback to oracle price for orderbooks with missing best bid": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_100000USD,
				constants.Dave_Num0_500000USD,
			},
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc,
			},
			orders: []types.Order{
				// Bid
				constants.Order_Carl_Num0_Id3_Clob0_Buy025BTC_Price50000,
				// Ask
				constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000,
			},

			expectedMidPrices: map[types.ClobPairId]types.Subticks{
				constants.ClobPair_Btc.GetClobPairId(): 50_000_000_000,
			},
		},
		"fallback to oracle price for orderbooks with missing best ask": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_100000USD,
				constants.Dave_Num0_500000USD,
			},
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc,
			},
			orders: []types.Order{
				// Ask
				constants.Order_Dave_Num0_Id1_Clob0_Sell025BTC_Price50000_GTB11,
				// Bid
				constants.Order_Carl_Num0_Id0_Clob0_Buy1BTC_Price500000_GTB10,
			},

			expectedMidPrices: map[types.ClobPairId]types.Subticks{
				constants.ClobPair_Btc.GetClobPairId(): 50_000_000_000,
			},
		},
		"fallback to oracle price for orderbooks with spread >= 1%": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
				constants.EthUsd_20PercentInitial_10PercentMaintenance,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_100000USD,
				constants.Dave_Num0_500000USD,
			},
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc,
				constants.ClobPair_Eth,
			},
			orders: []types.Order{
				// Bid
				constants.Order_Carl_Num0_Id0_Clob0_Buy1BTC_Price49500_GTB10,
				constants.Order_Carl_Num0_Id4_Clob1_Buy01ETH_Price3000,
				// Ask
				constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000, // Spread > 1%
				constants.Order_Dave_Num0_Id4_Clob1_Sell1ETH_Price3030,  // Spread == 1%
			},

			expectedMidPrices: map[types.ClobPairId]types.Subticks{
				constants.ClobPair_Btc.GetClobPairId(): 50_000_000_000,
				constants.ClobPair_Eth.GetClobPairId(): 3_000_000_000,
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			memclob := memclob.NewMemClobPriceTimePriority(false)
			mockBankKeeper := &mocks.BankKeeper{}
			mockBankKeeper.On(
				"SendCoinsFromModuleToModule",
				mock.Anything,
				mock.Anything,
				mock.Anything,
				mock.Anything,
			).Return(nil)
			mockBankKeeper.On(
				"SendCoins",
				mock.Anything,
				mock.Anything,
				mock.Anything,
				mock.Anything,
			).Return(nil)

			ks := keepertest.NewClobKeepersTestContext(t, memclob, mockBankKeeper, indexer_manager.NewIndexerEventManagerNoop())
			ctx := ks.Ctx.WithIsCheckTx(true).WithBlockTime(time.Unix(5, 0))

			// Create the default markets.
			keepertest.CreateTestMarkets(t, ctx, ks.PricesKeeper)

			// Create liquidity tiers.
			keepertest.CreateTestLiquidityTiers(t, ctx, ks.PerpetualsKeeper)

			require.NoError(t, ks.FeeTiersKeeper.SetPerpetualFeeParams(ctx, constants.PerpetualFeeParams))

			// Set up USDC asset in assets module.
			err := keepertest.CreateUsdcAsset(ctx, ks.AssetsKeeper)
			require.NoError(t, err)

			// Create all perpetuals.
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
				)
				require.NoError(t, err)
			}

			// Create all subaccounts.
			for _, subaccount := range tc.subaccounts {
				ks.SubaccountsKeeper.SetSubaccount(ctx, subaccount)
			}

			// Create all CLOBs.
			for _, clobPair := range tc.clobPairs {
				_, err = ks.ClobKeeper.CreatePerpetualClobPairAndMemStructs(
					ctx,
					clobPair.Id,
					clobtest.MustPerpetualId(clobPair),
					satypes.BaseQuantums(clobPair.StepBaseQuantums),
					clobPair.QuantumConversionExponent,
					clobPair.SubticksPerTick,
					clobPair.Status,
				)
				require.NoError(t, err)
			}

			// Place all orders.
			for _, order := range tc.orders {
				_, _, err := ks.ClobKeeper.PlaceShortTermOrder(ctx, &types.MsgPlaceOrder{Order: order})
				require.NoError(t, err)
			}

			clobMetadata := ks.ClobKeeper.GetClobMetadata(ctx)
			blockProposerPnL, validatorPnL := ks.ClobKeeper.InitializeCumulativePnLs(
				ctx,
				ks.PerpetualsKeeper,
				clobMetadata,
			)

			for clobPairId, expectedMidPrice := range tc.expectedMidPrices {
				require.Equal(
					t,
					expectedMidPrice,
					blockProposerPnL[clobPairId].Metadata.MidPrice,
				)
				require.Equal(
					t,
					expectedMidPrice,
					validatorPnL[clobPairId].Metadata.MidPrice,
				)
			}
		})
	}
}

func TestCumulativePnL_CalculateMev(t *testing.T) {
	tests := map[string]struct {
		blockProposerPnL *keeper.CumulativePnL
		validatorPnL     *keeper.CumulativePnL

		expectedResult *big.Float
	}{
		"can calculate difference for empty blocks": {
			blockProposerPnL: &keeper.CumulativePnL{
				SubaccountPnL: map[satypes.SubaccountId]*big.Int{},
			},
			validatorPnL: &keeper.CumulativePnL{
				SubaccountPnL: map[satypes.SubaccountId]*big.Int{},
			},

			expectedResult: big.NewFloat(0),
		},
		"can calculate difference for all positive PnLs": {
			blockProposerPnL: &keeper.CumulativePnL{
				SubaccountPnL: map[satypes.SubaccountId]*big.Int{
					constants.Alice_Num0: big.NewInt(100),
					constants.Bob_Num0:   big.NewInt(100),
				},
			},
			validatorPnL: &keeper.CumulativePnL{
				SubaccountPnL: map[satypes.SubaccountId]*big.Int{
					constants.Alice_Num0: big.NewInt(50),
					constants.Bob_Num0:   big.NewInt(50),
				},
			},

			expectedResult: big.NewFloat(50),
		},
		"can calculate difference for all negative PnLs": {
			blockProposerPnL: &keeper.CumulativePnL{
				SubaccountPnL: map[satypes.SubaccountId]*big.Int{
					constants.Alice_Num0: big.NewInt(100),
					constants.Bob_Num0:   big.NewInt(100),
				},
			},
			validatorPnL: &keeper.CumulativePnL{
				SubaccountPnL: map[satypes.SubaccountId]*big.Int{
					constants.Alice_Num0: big.NewInt(150),
					constants.Bob_Num0:   big.NewInt(150),
				},
			},

			expectedResult: big.NewFloat(50),
		},
		"can calculate difference between positive and negative PnLs": {
			blockProposerPnL: &keeper.CumulativePnL{
				SubaccountPnL: map[satypes.SubaccountId]*big.Int{
					constants.Alice_Num0: big.NewInt(100),
					constants.Bob_Num0:   big.NewInt(100),
					constants.Carl_Num0:  big.NewInt(-100),
					constants.Dave_Num0:  big.NewInt(-100),
				},
			},
			validatorPnL: &keeper.CumulativePnL{
				SubaccountPnL: map[satypes.SubaccountId]*big.Int{
					constants.Alice_Num0: big.NewInt(100),
					constants.Bob_Num0:   big.NewInt(-100),
					constants.Carl_Num0:  big.NewInt(100),
					constants.Dave_Num0:  big.NewInt(-100),
				},
			},

			expectedResult: big.NewFloat(200),
		},
		"can calculate difference for subaccounts only present in blockProposerPnL": {
			blockProposerPnL: &keeper.CumulativePnL{
				SubaccountPnL: map[satypes.SubaccountId]*big.Int{
					constants.Alice_Num0: big.NewInt(100),
					constants.Bob_Num0:   big.NewInt(100),
				},
			},
			validatorPnL: &keeper.CumulativePnL{
				SubaccountPnL: map[satypes.SubaccountId]*big.Int{
					constants.Alice_Num0: big.NewInt(100),
				},
			},

			expectedResult: big.NewFloat(50),
		},
		"can calculate difference for subaccounts only present in validatorPnL": {
			blockProposerPnL: &keeper.CumulativePnL{
				SubaccountPnL: map[satypes.SubaccountId]*big.Int{
					constants.Alice_Num0: big.NewInt(100),
				},
			},
			validatorPnL: &keeper.CumulativePnL{
				SubaccountPnL: map[satypes.SubaccountId]*big.Int{
					constants.Alice_Num0: big.NewInt(100),
					constants.Bob_Num0:   big.NewInt(-100),
				},
			},

			expectedResult: big.NewFloat(50),
		},
		"can calculate difference for completely disjoint PnLs": {
			blockProposerPnL: &keeper.CumulativePnL{
				SubaccountPnL: map[satypes.SubaccountId]*big.Int{
					constants.Alice_Num0: big.NewInt(100),
					constants.Bob_Num0:   big.NewInt(100),
				},
			},
			validatorPnL: &keeper.CumulativePnL{
				SubaccountPnL: map[satypes.SubaccountId]*big.Int{
					constants.Carl_Num0: big.NewInt(100),
					constants.Dave_Num0: big.NewInt(-100),
				},
			},

			expectedResult: big.NewFloat(200),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := tc.blockProposerPnL.CalculateMev(tc.validatorPnL)
			require.True(t, tc.expectedResult.Cmp(result) == 0)
		})
	}
}
