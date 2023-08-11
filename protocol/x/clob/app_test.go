package clob_test

import (
	"fmt"
	"testing"

	"github.com/cometbft/cometbft/crypto/tmhash"
	"github.com/cometbft/cometbft/types"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4/dtypes"
	"github.com/dydxprotocol/v4/indexer"
	indexerevents "github.com/dydxprotocol/v4/indexer/events"
	"github.com/dydxprotocol/v4/indexer/indexer_manager"
	"github.com/dydxprotocol/v4/indexer/msgsender"
	"github.com/dydxprotocol/v4/indexer/off_chain_updates"
	"github.com/dydxprotocol/v4/lib"
	testapp "github.com/dydxprotocol/v4/testutil/app"
	"github.com/dydxprotocol/v4/testutil/constants"
	testmsgs "github.com/dydxprotocol/v4/testutil/msgs"
	clobtypes "github.com/dydxprotocol/v4/x/clob/types"
	satypes "github.com/dydxprotocol/v4/x/subaccounts/types"
	"github.com/stretchr/testify/require"
)

var (
	Clob_0                                       = MustGetClobPairsFromGenesis(testapp.DefaultGenesis())[0]
	PlaceOrder_Alice_Num0_Id0_Buy5_Price10_GTB20 = *clobtypes.NewMsgPlaceOrder(MustScaleOrder(
		clobtypes.Order{
			OrderId:      clobtypes.OrderId{SubaccountId: constants.Alice_Num0, ClientId: 0, ClobPairId: 0},
			Side:         clobtypes.Order_SIDE_BUY,
			Quantums:     5,
			Subticks:     10,
			GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
		},
		testapp.DefaultGenesis(),
	))
	PlaceOrder_Bob_Num0_Id0_Sell5_Price10_GTB20 = *clobtypes.NewMsgPlaceOrder(MustScaleOrder(
		clobtypes.Order{
			OrderId:      clobtypes.OrderId{SubaccountId: constants.Bob_Num0, ClientId: 0, ClobPairId: 0},
			Side:         clobtypes.Order_SIDE_SELL,
			Quantums:     5,
			Subticks:     10,
			GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
		},
		testapp.DefaultGenesis(),
	))

	LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB5 = *clobtypes.NewMsgPlaceOrder(MustScaleOrder(
		constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB5,
		testapp.DefaultGenesis(),
	))
	LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT5 = *clobtypes.NewMsgPlaceOrder(MustScaleOrder(
		constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT5,
		testapp.DefaultGenesis(),
	))
	LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15 = *clobtypes.NewMsgPlaceOrder(MustScaleOrder(
		constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
		testapp.DefaultGenesis(),
	))
	ConditionalPlaceOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTB15 = *clobtypes.NewMsgPlaceOrder(MustScaleOrder(
		constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTB15,
		testapp.DefaultGenesis(),
	))
	ConditionalPlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15 = *clobtypes.NewMsgPlaceOrder(MustScaleOrder(
		constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
		testapp.DefaultGenesis(),
	))
)

func TestPlaceOrder(t *testing.T) {
	msgSender := msgsender.NewIndexerMessageSenderInMemoryCollector()
	appOpts := map[string]interface{}{
		indexer.MsgSenderInstanceForTest: msgSender,
	}
	tAppBuilder := testapp.NewTestAppBuilder().WithAppCreatorFn(testapp.DefaultTestAppCreatorFn(appOpts))
	tApp := tAppBuilder.Build()
	ctx, app := tApp.InitChain()

	aliceSubaccount := app.SubaccountsKeeper.GetSubaccount(ctx, constants.Alice_Num0)
	bobSubaccount := app.SubaccountsKeeper.GetSubaccount(ctx, constants.Bob_Num0)

	CheckTx_PlaceOrder_Alice_Num0_Id0_Buy5_Price10_GTB20 := testapp.MustMakeCheckTx(
		ctx,
		app,
		&PlaceOrder_Alice_Num0_Id0_Buy5_Price10_GTB20,
	)
	CheckTx_PlaceOrder_Bob_Num0_Id0_Sell5_Price10_GTB20 := testapp.MustMakeCheckTx(
		ctx,
		app,
		&PlaceOrder_Bob_Num0_Id0_Sell5_Price10_GTB20,
	)

	tests := map[string]struct {
		orders                                  []clobtypes.MsgPlaceOrder
		expectedOrdersFilled                    []clobtypes.OrderId
		expectedOffchainMessagesAfterPlaceOrder []msgsender.Message
		expectedOnchainMessagesAfterPlaceOrder  []msgsender.Message
		expectedOffchainMessagesInNextBlock     []msgsender.Message
		expectedOnchainMessagesInNextBlock      []msgsender.Message
	}{
		"Test placing an order": {
			orders: []clobtypes.MsgPlaceOrder{PlaceOrder_Alice_Num0_Id0_Buy5_Price10_GTB20},
			expectedOffchainMessagesAfterPlaceOrder: []msgsender.Message{
				off_chain_updates.MustCreateOrderPlaceMessage(
					nil,
					PlaceOrder_Alice_Num0_Id0_Buy5_Price10_GTB20.Order,
				).AddHeader(msgsender.MessageHeader{
					Key:   msgsender.TransactionHashHeaderKey,
					Value: tmhash.Sum(CheckTx_PlaceOrder_Alice_Num0_Id0_Buy5_Price10_GTB20.Tx),
				}),
				off_chain_updates.MustCreateOrderUpdateMessage(
					nil,
					PlaceOrder_Alice_Num0_Id0_Buy5_Price10_GTB20.Order.OrderId,
					0,
				).AddHeader(msgsender.MessageHeader{
					Key:   msgsender.TransactionHashHeaderKey,
					Value: tmhash.Sum(CheckTx_PlaceOrder_Alice_Num0_Id0_Buy5_Price10_GTB20.Tx),
				}),
			},
			expectedOnchainMessagesInNextBlock: []msgsender.Message{indexer_manager.CreateIndexerBlockEventMessage(
				&indexer_manager.IndexerTendermintBlock{
					Height: 3,
					Time:   ctx.BlockTime(),
				})},
		},
		"Test matching an order fully": {
			orders: []clobtypes.MsgPlaceOrder{
				PlaceOrder_Alice_Num0_Id0_Buy5_Price10_GTB20,
				PlaceOrder_Bob_Num0_Id0_Sell5_Price10_GTB20,
			},
			expectedOffchainMessagesAfterPlaceOrder: []msgsender.Message{
				off_chain_updates.MustCreateOrderPlaceMessage(
					nil,
					PlaceOrder_Alice_Num0_Id0_Buy5_Price10_GTB20.Order,
				).AddHeader(msgsender.MessageHeader{
					Key:   msgsender.TransactionHashHeaderKey,
					Value: tmhash.Sum(CheckTx_PlaceOrder_Alice_Num0_Id0_Buy5_Price10_GTB20.Tx),
				}),
				off_chain_updates.MustCreateOrderUpdateMessage(
					nil,
					PlaceOrder_Alice_Num0_Id0_Buy5_Price10_GTB20.Order.OrderId,
					0,
				).AddHeader(msgsender.MessageHeader{
					Key:   msgsender.TransactionHashHeaderKey,
					Value: tmhash.Sum(CheckTx_PlaceOrder_Alice_Num0_Id0_Buy5_Price10_GTB20.Tx),
				}),
				off_chain_updates.MustCreateOrderPlaceMessage(
					nil,
					PlaceOrder_Bob_Num0_Id0_Sell5_Price10_GTB20.Order,
				).AddHeader(msgsender.MessageHeader{
					Key:   msgsender.TransactionHashHeaderKey,
					Value: tmhash.Sum(CheckTx_PlaceOrder_Bob_Num0_Id0_Sell5_Price10_GTB20.Tx),
				}),
				off_chain_updates.MustCreateOrderUpdateMessage(
					nil,
					PlaceOrder_Alice_Num0_Id0_Buy5_Price10_GTB20.Order.OrderId,
					PlaceOrder_Alice_Num0_Id0_Buy5_Price10_GTB20.Order.GetBaseQuantums(),
				).AddHeader(msgsender.MessageHeader{
					Key:   msgsender.TransactionHashHeaderKey,
					Value: tmhash.Sum(CheckTx_PlaceOrder_Bob_Num0_Id0_Sell5_Price10_GTB20.Tx),
				}),
				off_chain_updates.MustCreateOrderUpdateMessage(
					nil,
					PlaceOrder_Bob_Num0_Id0_Sell5_Price10_GTB20.Order.OrderId,
					PlaceOrder_Bob_Num0_Id0_Sell5_Price10_GTB20.Order.GetBaseQuantums(),
				).AddHeader(msgsender.MessageHeader{
					Key:   msgsender.TransactionHashHeaderKey,
					Value: tmhash.Sum(CheckTx_PlaceOrder_Bob_Num0_Id0_Sell5_Price10_GTB20.Tx),
				}),
			},
			expectedOffchainMessagesInNextBlock: []msgsender.Message{
				off_chain_updates.MustCreateOrderUpdateMessage(
					nil,
					PlaceOrder_Alice_Num0_Id0_Buy5_Price10_GTB20.Order.OrderId,
					PlaceOrder_Alice_Num0_Id0_Buy5_Price10_GTB20.Order.GetBaseQuantums(),
				),
				off_chain_updates.MustCreateOrderUpdateMessage(
					nil,
					PlaceOrder_Bob_Num0_Id0_Sell5_Price10_GTB20.Order.OrderId,
					PlaceOrder_Bob_Num0_Id0_Sell5_Price10_GTB20.Order.GetBaseQuantums(),
				),
			},
			expectedOnchainMessagesInNextBlock: []msgsender.Message{indexer_manager.CreateIndexerBlockEventMessage(
				&indexer_manager.IndexerTendermintBlock{
					Height: 3,
					Time:   ctx.BlockTime(),
					Events: []*indexer_manager.IndexerTendermintEvent{
						{
							Subtype: indexerevents.SubtypeSubaccountUpdate,
							Data: indexer_manager.GetB64EncodedEventMessage(
								indexerevents.NewSubaccountUpdateEvent(
									&constants.Bob_Num0,
									[]*satypes.PerpetualPosition{
										{
											PerpetualId: Clob_0.MustGetPerpetualId(),
											Quantums: dtypes.NewInt(-int64(
												PlaceOrder_Alice_Num0_Id0_Buy5_Price10_GTB20.Order.GetQuantums())),
											FundingIndex: dtypes.NewInt(0),
										},
									},
									// Maker fees calculate to 0 so asset position doesn't change.
									[]*satypes.AssetPosition{
										{
											AssetId:  lib.UsdcAssetId,
											Quantums: dtypes.NewIntFromBigInt(bobSubaccount.GetUsdcPosition()),
										},
									},
								),
							),
							OrderingWithinBlock: &indexer_manager.IndexerTendermintEvent_TransactionIndex{},
							EventIndex:          0,
						},
						{
							Subtype: indexerevents.SubtypeSubaccountUpdate,
							Data: indexer_manager.GetB64EncodedEventMessage(
								indexerevents.NewSubaccountUpdateEvent(
									&constants.Alice_Num0,
									[]*satypes.PerpetualPosition{
										{
											PerpetualId: Clob_0.MustGetPerpetualId(),
											Quantums: dtypes.NewInt(int64(
												PlaceOrder_Alice_Num0_Id0_Buy5_Price10_GTB20.Order.GetQuantums())),
											FundingIndex: dtypes.NewInt(0),
										},
									},
									// Taker fees calculate to 0 so asset position doesn't change.
									[]*satypes.AssetPosition{
										{
											AssetId:  lib.UsdcAssetId,
											Quantums: dtypes.NewIntFromBigInt(aliceSubaccount.GetUsdcPosition()),
										},
									},
								),
							),
							OrderingWithinBlock: &indexer_manager.IndexerTendermintEvent_TransactionIndex{},
							EventIndex:          1,
						},
						{
							Subtype: indexerevents.SubtypeOrderFill,
							Data: indexer_manager.GetB64EncodedEventMessage(
								indexerevents.NewOrderFillEvent(
									PlaceOrder_Alice_Num0_Id0_Buy5_Price10_GTB20.Order,
									PlaceOrder_Bob_Num0_Id0_Sell5_Price10_GTB20.Order,
									PlaceOrder_Alice_Num0_Id0_Buy5_Price10_GTB20.Order.GetBaseQuantums(),
								),
							),
							OrderingWithinBlock: &indexer_manager.IndexerTendermintEvent_TransactionIndex{},
							EventIndex:          2,
						},
					},
					TxHashes: []string{string(lib.GetTxHash(MustGetTxBytes(&clobtypes.MsgProposedOperations{
						Proposer: constants.AliceAccAddress.String(),
						OperationsQueue: []clobtypes.Operation{
							clobtypes.NewOrderPlacementOperation(
								PlaceOrder_Alice_Num0_Id0_Buy5_Price10_GTB20.Order),
							clobtypes.NewOrderPlacementOperation(
								PlaceOrder_Bob_Num0_Id0_Sell5_Price10_GTB20.Order),
							clobtypes.NewMatchOperation(
								&PlaceOrder_Bob_Num0_Id0_Sell5_Price10_GTB20.Order,
								[]clobtypes.MakerFill{
									{
										FillAmount: PlaceOrder_Alice_Num0_Id0_Buy5_Price10_GTB20.
											Order.GetBaseQuantums().ToUint64(),
										MakerOrderId: PlaceOrder_Alice_Num0_Id0_Buy5_Price10_GTB20.Order.OrderId,
									},
								},
							),
						},
						AddToOrderbookCollatCheckOrderHashes: [][]byte{
							PlaceOrder_Alice_Num0_Id0_Buy5_Price10_GTB20.Order.GetOrderHash().ToBytes(),
							PlaceOrder_Bob_Num0_Id0_Sell5_Price10_GTB20.Order.GetOrderHash().ToBytes(),
						},
					},
					)))},
				})},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Reset for each iteration of the loop
			tApp.Reset()

			ctx, app = tApp.AdvanceToBlock(2)
			// Clear any messages produced prior to these checkTx calls.
			msgSender.Clear()
			for _, checkTx := range testapp.MustMakeCheckTxs(ctx, app, tc.orders...) {
				require.True(t, app.CheckTx(checkTx).IsOK())
			}
			require.ElementsMatch(
				t,
				tc.expectedOffchainMessagesAfterPlaceOrder,
				msgSender.GetOffchainMessages(),
			)
			require.ElementsMatch(t, tc.expectedOnchainMessagesAfterPlaceOrder, msgSender.GetOnchainMessages())

			// Clear the messages that we already matched prior to advancing to the next block.
			msgSender.Clear()
			ctx, app = tApp.AdvanceToBlock(3)
			require.ElementsMatch(t, tc.expectedOffchainMessagesInNextBlock, msgSender.GetOffchainMessages())
			require.ElementsMatch(t, tc.expectedOnchainMessagesInNextBlock, msgSender.GetOnchainMessages())
			for _, order := range tc.orders {
				if lib.ContainsValue(tc.expectedOrdersFilled, order.Order.OrderId) {
					exists, fillAmount, _ := app.ClobKeeper.GetOrderFillAmount(
						ctx,
						order.Order.OrderId,
					)

					require.True(t, exists)
					require.Equal(t, order.Order.GetBaseQuantums(), fillAmount)
				}
			}
		})
	}
}

func TestCancelStatefulOrder(t *testing.T) {
	tests := map[string]struct {
		blockWithMessages        []testmsgs.TestBlockWithMsgs
		checkCancelledPlaceOrder clobtypes.MsgPlaceOrder
		checkResultsBlock        uint32
	}{
		"Test stateful order is cancelled when placed and cancelled in the same block": {
			blockWithMessages: []testmsgs.TestBlockWithMsgs{
				{
					Block: 2,
					Msgs: []testmsgs.TestSdkMsg{
						{
							Msg:          &LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT5,
							ExpectedIsOk: true,
						},
						{
							Msg:          &constants.CancelLongTermOrder_Alice_Num0_Id0_Clob0_GTBT15,
							ExpectedIsOk: true,
						},
					},
				},
			},

			checkCancelledPlaceOrder: LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT5,
			checkResultsBlock:        4,
		},
		"Test stateful order is cancelled when placed then cancelled in a future block": {
			blockWithMessages: []testmsgs.TestBlockWithMsgs{
				{
					Block: 2,
					Msgs: []testmsgs.TestSdkMsg{{
						Msg:          &LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT5,
						ExpectedIsOk: true,
					}},
				},
				{
					Block: 3,
					Msgs: []testmsgs.TestSdkMsg{{
						Msg:          &constants.CancelLongTermOrder_Alice_Num0_Id0_Clob0_GTBT15,
						ExpectedIsOk: true,
					}},
				},
			},

			checkCancelledPlaceOrder: LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT5,
			checkResultsBlock:        4,
		},
		"Test stateful order is cancelled when placed, matched, and cancelled in the same block": {
			blockWithMessages: []testmsgs.TestBlockWithMsgs{
				{
					Block: 2,
					Msgs: []testmsgs.TestSdkMsg{
						{
							Msg:          &LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT5,
							ExpectedIsOk: true,
						},
						{
							Msg:          &PlaceOrder_Bob_Num0_Id0_Sell5_Price10_GTB20,
							ExpectedIsOk: true,
						},
						{
							Msg:          &constants.CancelLongTermOrder_Alice_Num0_Id0_Clob0_GTBT15,
							ExpectedIsOk: true,
						},
					},
				},
			},

			checkCancelledPlaceOrder: LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT5,
			checkResultsBlock:        4,
		},
		"Test stateful order is cancelled when placed, cancelled, then re-placed with the same order id": {
			blockWithMessages: []testmsgs.TestBlockWithMsgs{
				{
					Block: 2,
					Msgs: []testmsgs.TestSdkMsg{
						{
							Msg:          &LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT5,
							ExpectedIsOk: true,
						},
						{
							Msg:          &constants.CancelLongTermOrder_Alice_Num0_Id0_Clob0_GTBT15,
							ExpectedIsOk: true,
						},
					},
				},
				{
					Block: 3,
					Msgs: []testmsgs.TestSdkMsg{{
						Msg:          &LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
						ExpectedIsOk: true,
					}},
				},
			},

			checkCancelledPlaceOrder: LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT5,
			checkResultsBlock:        4,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder().WithTesting(t).Build()

			for _, blockWithMessages := range tc.blockWithMessages {
				ctx, app := tApp.AdvanceToBlock(blockWithMessages.Block)

				for _, testSdkMsg := range blockWithMessages.Msgs {
					result := app.CheckTx(testapp.MustMakeCheckTx(ctx, app, testSdkMsg.Msg)).IsOK()
					require.Equal(
						t,
						testSdkMsg.ExpectedIsOk,
						result,
					)
				}
			}

			ctx, app := tApp.AdvanceToBlockIfNecessary(tc.checkResultsBlock)
			require.True(t, app.ClobKeeper.HasSeenPlaceOrder(ctx, tc.checkCancelledPlaceOrder))
			exists, fillAmount, _ := app.ClobKeeper.GetOrderFillAmount(
				ctx,
				tc.checkCancelledPlaceOrder.Order.OrderId,
			)
			require.False(t, exists)
			require.Equal(t, uint64(0), fillAmount.ToUint64())
		})
	}
}

// Returns the encoded msg as transaction. Will panic if encoding fails.
func MustGetTxBytes(msgs ...sdktypes.Msg) []byte {
	tx := constants.TestEncodingCfg.TxConfig.NewTxBuilder()
	err := tx.SetMsgs(msgs...)
	if err != nil {
		panic(err)
	}
	bz, err := constants.TestEncodingCfg.TxConfig.TxEncoder()(tx.GetTx())
	if err != nil {
		panic(err)
	}
	return bz
}

// MustScaleOrder scales clobtypes.Order and clobtypes.MsgPlaceorder based upon the clob information provided.
// Will panic if:
//   - OrderT is an unknown type.
//   - ClobPairSrcT is an unknown type.
//   - The clob pair id can't be used to look up the clob pair from the ClobPairSrcT.
func MustScaleOrder[
	OrderT clobtypes.Order | clobtypes.MsgPlaceOrder,
	ClobPairSrcT clobtypes.ClobPair | types.GenesisDoc](
	order OrderT,
	clobPairSrc ClobPairSrcT,
) OrderT {
	var msgPlaceOrder clobtypes.MsgPlaceOrder

	// Find the clob pair id based upon the type of order passed in.
	var clobPairId clobtypes.ClobPairId
	switch v := any(order).(type) {
	case clobtypes.MsgPlaceOrder:
		clobPairId = v.Order.GetClobPairId()
		msgPlaceOrder = v
	case clobtypes.Order:
		clobPairId = v.GetClobPairId()
		msgPlaceOrder = *clobtypes.NewMsgPlaceOrder(v)
	default:
		panic(fmt.Errorf("Unknown order type %T to get order", order))
	}

	// Find the clob pair based upon the clobPairSrc of the clob information passed in.
	var clobPair clobtypes.ClobPair
	switch v := any(clobPairSrc).(type) {
	case clobtypes.ClobPair:
		clobPair = v
	case types.GenesisDoc:
		clobPairs := MustGetClobPairsFromGenesis(v)
		if hasClobPair, ok := clobPairs[clobPairId]; ok {
			clobPair = hasClobPair
		} else {
			panic(fmt.Errorf("Clob not found in genesis doc for clob id %d", clobPairId))
		}
	default:
		panic(fmt.Errorf("Unknown source type %T to get clob pair", clobPairSrc))
	}

	// Scale the order based upon the quantums and subticks passed in.
	msgPlaceOrder.Order.Quantums = msgPlaceOrder.Order.Quantums * clobPair.MinOrderBaseQuantums
	msgPlaceOrder.Order.Subticks = msgPlaceOrder.Order.Subticks * uint64(clobPair.SubticksPerTick)

	// Return a type that matches what the user passed in for the order type.
	switch any(order).(type) {
	case clobtypes.MsgPlaceOrder:
		return any(msgPlaceOrder).(OrderT)
	case clobtypes.Order:
		return any(msgPlaceOrder.Order).(OrderT)
	default:
		panic(fmt.Errorf("Unable to convert to %T to %T", clobtypes.MsgPlaceOrder{}, order))
	}
}

// MustGetClobPairsFromGenesis unmarshals the initial genesis state and returns a map from clob pair id to clob pair.
func MustGetClobPairsFromGenesis(genesisDoc types.GenesisDoc) map[clobtypes.ClobPairId]clobtypes.ClobPair {
	var genesisState clobtypes.GenesisState
	testapp.UpdateGenesisDocWithAppStateForModule(&genesisDoc, func(genesisStatePtr *clobtypes.GenesisState) {
		genesisState = *genesisStatePtr
	})

	clobPairs := make(map[clobtypes.ClobPairId]clobtypes.ClobPair, len(genesisState.ClobPairs))
	for _, clobPair := range genesisState.ClobPairs {
		clobPairs[clobPair.GetClobPairId()] = clobPair
	}
	return clobPairs
}
