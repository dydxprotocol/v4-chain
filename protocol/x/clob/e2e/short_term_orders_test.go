package clob_test

import (
	"testing"

	"golang.org/x/exp/slices"

	"github.com/cometbft/cometbft/crypto/tmhash"
	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	"github.com/dydxprotocol/v4-chain/protocol/lib"

	"github.com/dydxprotocol/v4-chain/protocol/indexer"
	indexerevents "github.com/dydxprotocol/v4-chain/protocol/indexer/events"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/msgsender"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/off_chain_updates"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	clobtestutils "github.com/dydxprotocol/v4-chain/protocol/testutil/clob"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	testtx "github.com/dydxprotocol/v4-chain/protocol/testutil/tx"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/require"
)

func TestPlaceOrder(t *testing.T) {
	msgSender := msgsender.NewIndexerMessageSenderInMemoryCollector()
	appOpts := map[string]interface{}{
		indexer.MsgSenderInstanceForTest: msgSender,
	}
	tAppBuilder := testapp.NewTestAppBuilder().WithAppCreatorFn(testapp.DefaultTestAppCreatorFn(appOpts))
	tApp := tAppBuilder.Build()
	ctx := tApp.InitChain()

	aliceSubaccount := tApp.App.SubaccountsKeeper.GetSubaccount(ctx, constants.Alice_Num0)
	bobSubaccount := tApp.App.SubaccountsKeeper.GetSubaccount(ctx, constants.Bob_Num0)

	CheckTx_PlaceOrder_Alice_Num0_Id0_Buy5_Price10_GTB20 := testapp.MustMakeCheckTx(
		ctx,
		tApp.App,
		testapp.MustMakeCheckTxOptions{
			AccAddressForSigning: testtx.MustGetOnlySignerAddress(&PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB20),
		},
		&PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB20,
	)
	CheckTx_PlaceOrder_Alice_Num0_Id0_Buy6_Price10_GTB20 := testapp.MustMakeCheckTx(
		ctx,
		tApp.App,
		testapp.MustMakeCheckTxOptions{
			AccAddressForSigning: testtx.MustGetOnlySignerAddress(&PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20),
		},
		&PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20,
	)
	CheckTx_PlaceOrder_Bob_Num0_Id0_Sell5_Price10_GTB20 := testapp.MustMakeCheckTx(
		ctx,
		tApp.App,
		testapp.MustMakeCheckTxOptions{
			AccAddressForSigning: testtx.MustGetOnlySignerAddress(&PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20),
		},
		&PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20,
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
			orders: []clobtypes.MsgPlaceOrder{PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB20},
			expectedOffchainMessagesAfterPlaceOrder: []msgsender.Message{
				off_chain_updates.MustCreateOrderPlaceMessage(
					nil,
					PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB20.Order,
				).AddHeader(msgsender.MessageHeader{
					Key:   msgsender.TransactionHashHeaderKey,
					Value: tmhash.Sum(CheckTx_PlaceOrder_Alice_Num0_Id0_Buy5_Price10_GTB20.Tx),
				}),
				off_chain_updates.MustCreateOrderUpdateMessage(
					nil,
					PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB20.Order.OrderId,
					0,
				).AddHeader(msgsender.MessageHeader{
					Key:   msgsender.TransactionHashHeaderKey,
					Value: tmhash.Sum(CheckTx_PlaceOrder_Alice_Num0_Id0_Buy5_Price10_GTB20.Tx),
				}),
			},
			expectedOnchainMessagesInNextBlock: []msgsender.Message{indexer_manager.CreateIndexerBlockEventMessage(
				&indexer_manager.IndexerTendermintBlock{
					Height: 2,
					Time:   ctx.BlockTime(),
				})},
		},
		"Test matching an order fully": {
			orders: []clobtypes.MsgPlaceOrder{
				PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB20,
				PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20,
			},
			expectedOrdersFilled: []clobtypes.OrderId{
				PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB20.Order.OrderId,
				PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order.OrderId,
			},
			expectedOffchainMessagesAfterPlaceOrder: []msgsender.Message{
				off_chain_updates.MustCreateOrderPlaceMessage(
					nil,
					PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB20.Order,
				).AddHeader(msgsender.MessageHeader{
					Key:   msgsender.TransactionHashHeaderKey,
					Value: tmhash.Sum(CheckTx_PlaceOrder_Alice_Num0_Id0_Buy5_Price10_GTB20.Tx),
				}),
				off_chain_updates.MustCreateOrderUpdateMessage(
					nil,
					PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB20.Order.OrderId,
					0,
				).AddHeader(msgsender.MessageHeader{
					Key:   msgsender.TransactionHashHeaderKey,
					Value: tmhash.Sum(CheckTx_PlaceOrder_Alice_Num0_Id0_Buy5_Price10_GTB20.Tx),
				}),
				off_chain_updates.MustCreateOrderPlaceMessage(
					nil,
					PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order,
				).AddHeader(msgsender.MessageHeader{
					Key:   msgsender.TransactionHashHeaderKey,
					Value: tmhash.Sum(CheckTx_PlaceOrder_Bob_Num0_Id0_Sell5_Price10_GTB20.Tx),
				}),
				off_chain_updates.MustCreateOrderUpdateMessage(
					nil,
					PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB20.Order.OrderId,
					PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB20.Order.GetBaseQuantums(),
				).AddHeader(msgsender.MessageHeader{
					Key:   msgsender.TransactionHashHeaderKey,
					Value: tmhash.Sum(CheckTx_PlaceOrder_Bob_Num0_Id0_Sell5_Price10_GTB20.Tx),
				}),
				off_chain_updates.MustCreateOrderUpdateMessage(
					nil,
					PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order.OrderId,
					PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order.GetBaseQuantums(),
				).AddHeader(msgsender.MessageHeader{
					Key:   msgsender.TransactionHashHeaderKey,
					Value: tmhash.Sum(CheckTx_PlaceOrder_Bob_Num0_Id0_Sell5_Price10_GTB20.Tx),
				}),
			},
			expectedOnchainMessagesInNextBlock: []msgsender.Message{indexer_manager.CreateIndexerBlockEventMessage(
				&indexer_manager.IndexerTendermintBlock{
					Height: 2,
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
												PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB20.Order.GetQuantums())),
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
									nil, // no funding payments
								),
							),
							OrderingWithinBlock: &indexer_manager.IndexerTendermintEvent_TransactionIndex{},
							EventIndex:          0,
							Version:             indexerevents.SubaccountUpdateEventVersion,
							DataBytes: indexer_manager.GetBytes(
								indexerevents.NewSubaccountUpdateEvent(
									&constants.Bob_Num0,
									[]*satypes.PerpetualPosition{
										{
											PerpetualId: Clob_0.MustGetPerpetualId(),
											Quantums: dtypes.NewInt(-int64(
												PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB20.Order.GetQuantums())),
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
									nil, // no funding payments
								),
							),
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
												PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB20.Order.GetQuantums())),
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
									nil, // no funding payments
								),
							),
							OrderingWithinBlock: &indexer_manager.IndexerTendermintEvent_TransactionIndex{},
							EventIndex:          1,
							Version:             indexerevents.SubaccountUpdateEventVersion,
							DataBytes: indexer_manager.GetBytes(
								indexerevents.NewSubaccountUpdateEvent(
									&constants.Alice_Num0,
									[]*satypes.PerpetualPosition{
										{
											PerpetualId: Clob_0.MustGetPerpetualId(),
											Quantums: dtypes.NewInt(int64(
												PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB20.Order.GetQuantums())),
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
									nil, // no funding payments
								),
							),
						},
						{
							Subtype: indexerevents.SubtypeOrderFill,
							Data: indexer_manager.GetB64EncodedEventMessage(
								indexerevents.NewOrderFillEvent(
									PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB20.Order,
									PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order,
									PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB20.Order.GetBaseQuantums(),
									0, // Fees are 0 due to lost precision
									0,
									PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB20.Order.GetBaseQuantums(),
									PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB20.Order.GetBaseQuantums(),
								),
							),
							OrderingWithinBlock: &indexer_manager.IndexerTendermintEvent_TransactionIndex{},
							EventIndex:          2,
							Version:             indexerevents.OrderFillEventVersion,
							DataBytes: indexer_manager.GetBytes(
								indexerevents.NewOrderFillEvent(
									PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB20.Order,
									PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order,
									PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB20.Order.GetBaseQuantums(),
									0, // Fees are 0 due to lost precision
									0,
									PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB20.Order.GetBaseQuantums(),
									PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB20.Order.GetBaseQuantums(),
								),
							),
						},
					},
					TxHashes: []string{string(lib.GetTxHash(testtx.MustGetTxBytes(&clobtypes.MsgProposedOperations{
						OperationsQueue: []clobtypes.OperationRaw{
							{
								Operation: &clobtypes.OperationRaw_ShortTermOrderPlacement{
									ShortTermOrderPlacement: CheckTx_PlaceOrder_Alice_Num0_Id0_Buy5_Price10_GTB20.Tx,
								},
							},
							{
								Operation: &clobtypes.OperationRaw_ShortTermOrderPlacement{
									ShortTermOrderPlacement: CheckTx_PlaceOrder_Bob_Num0_Id0_Sell5_Price10_GTB20.Tx,
								},
							},
							clobtestutils.NewMatchOperationRaw(
								&PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order,
								[]clobtypes.MakerFill{
									{
										FillAmount: PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB20.
											Order.GetBaseQuantums().ToUint64(),
										MakerOrderId: PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB20.Order.OrderId,
									},
								},
							),
						},
					},
					)))},
				})},
		},
		"Test matching an order partially, maker order remains on book": {
			orders: []clobtypes.MsgPlaceOrder{
				PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20,
				PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20,
			},
			expectedOrdersFilled: []clobtypes.OrderId{
				PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order.OrderId,
			},
			expectedOffchainMessagesAfterPlaceOrder: []msgsender.Message{
				off_chain_updates.MustCreateOrderPlaceMessage(
					nil,
					PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20.Order,
				).AddHeader(msgsender.MessageHeader{
					Key:   msgsender.TransactionHashHeaderKey,
					Value: tmhash.Sum(CheckTx_PlaceOrder_Alice_Num0_Id0_Buy6_Price10_GTB20.Tx),
				}),
				off_chain_updates.MustCreateOrderUpdateMessage(
					nil,
					PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20.Order.OrderId,
					0,
				).AddHeader(msgsender.MessageHeader{
					Key:   msgsender.TransactionHashHeaderKey,
					Value: tmhash.Sum(CheckTx_PlaceOrder_Alice_Num0_Id0_Buy6_Price10_GTB20.Tx),
				}),
				off_chain_updates.MustCreateOrderPlaceMessage(
					nil,
					PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order,
				).AddHeader(msgsender.MessageHeader{
					Key:   msgsender.TransactionHashHeaderKey,
					Value: tmhash.Sum(CheckTx_PlaceOrder_Bob_Num0_Id0_Sell5_Price10_GTB20.Tx),
				}),
				off_chain_updates.MustCreateOrderUpdateMessage(
					nil,
					PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20.Order.OrderId,
					PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order.GetBaseQuantums(),
				).AddHeader(msgsender.MessageHeader{
					Key:   msgsender.TransactionHashHeaderKey,
					Value: tmhash.Sum(CheckTx_PlaceOrder_Bob_Num0_Id0_Sell5_Price10_GTB20.Tx),
				}),
				off_chain_updates.MustCreateOrderUpdateMessage(
					nil,
					PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order.OrderId,
					PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order.GetBaseQuantums(),
				).AddHeader(msgsender.MessageHeader{
					Key:   msgsender.TransactionHashHeaderKey,
					Value: tmhash.Sum(CheckTx_PlaceOrder_Bob_Num0_Id0_Sell5_Price10_GTB20.Tx),
				}),
			},
			expectedOffchainMessagesInNextBlock: []msgsender.Message{
				off_chain_updates.MustCreateOrderUpdateMessage(
					nil,
					PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20.Order.OrderId,
					PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order.GetBaseQuantums(),
				),
			},
			expectedOnchainMessagesInNextBlock: []msgsender.Message{indexer_manager.CreateIndexerBlockEventMessage(
				&indexer_manager.IndexerTendermintBlock{
					Height: 2,
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
												PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order.GetQuantums())),
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
									nil, // no funding payments
								),
							),
							OrderingWithinBlock: &indexer_manager.IndexerTendermintEvent_TransactionIndex{},
							EventIndex:          0,
							Version:             indexerevents.SubaccountUpdateEventVersion,
							DataBytes: indexer_manager.GetBytes(
								indexerevents.NewSubaccountUpdateEvent(
									&constants.Bob_Num0,
									[]*satypes.PerpetualPosition{
										{
											PerpetualId: Clob_0.MustGetPerpetualId(),
											Quantums: dtypes.NewInt(-int64(
												PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order.GetQuantums())),
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
									nil, // no funding payments
								),
							),
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
												PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order.GetQuantums())),
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
									nil, // no funding payments
								),
							),
							OrderingWithinBlock: &indexer_manager.IndexerTendermintEvent_TransactionIndex{},
							EventIndex:          1,
							Version:             indexerevents.SubaccountUpdateEventVersion,
							DataBytes: indexer_manager.GetBytes(
								indexerevents.NewSubaccountUpdateEvent(
									&constants.Alice_Num0,
									[]*satypes.PerpetualPosition{
										{
											PerpetualId: Clob_0.MustGetPerpetualId(),
											Quantums: dtypes.NewInt(int64(
												PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order.GetQuantums())),
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
									nil, // no funding payments
								),
							),
						},
						{
							Subtype: indexerevents.SubtypeOrderFill,
							Data: indexer_manager.GetB64EncodedEventMessage(
								indexerevents.NewOrderFillEvent(
									PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20.Order,
									PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order,
									PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order.GetBaseQuantums(),
									0, // Fees are 0 due to lost precision
									0,
									PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order.GetBaseQuantums(),
									PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order.GetBaseQuantums(),
								),
							),
							OrderingWithinBlock: &indexer_manager.IndexerTendermintEvent_TransactionIndex{},
							EventIndex:          2,
							Version:             indexerevents.OrderFillEventVersion,
							DataBytes: indexer_manager.GetBytes(
								indexerevents.NewOrderFillEvent(
									PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20.Order,
									PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order,
									PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order.GetBaseQuantums(),
									0, // Fees are 0 due to lost precision
									0,
									PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order.GetBaseQuantums(),
									PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order.GetBaseQuantums(),
								),
							),
						},
					},
					TxHashes: []string{string(lib.GetTxHash(testtx.MustGetTxBytes(&clobtypes.MsgProposedOperations{
						OperationsQueue: []clobtypes.OperationRaw{
							{
								Operation: &clobtypes.OperationRaw_ShortTermOrderPlacement{
									ShortTermOrderPlacement: CheckTx_PlaceOrder_Alice_Num0_Id0_Buy6_Price10_GTB20.Tx,
								},
							},
							{
								Operation: &clobtypes.OperationRaw_ShortTermOrderPlacement{
									ShortTermOrderPlacement: CheckTx_PlaceOrder_Bob_Num0_Id0_Sell5_Price10_GTB20.Tx,
								},
							},
							clobtestutils.NewMatchOperationRaw(
								&PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order,
								[]clobtypes.MakerFill{
									{
										FillAmount: PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.
											Order.GetBaseQuantums().ToUint64(),
										MakerOrderId: PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20.Order.OrderId,
									},
								},
							),
						},
					},
					)))},
				})},
		},
		"Test matching an order partially, taker order remains on book": {
			orders: []clobtypes.MsgPlaceOrder{
				PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20,
				PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20,
			},
			expectedOrdersFilled: []clobtypes.OrderId{
				PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order.OrderId,
			},
			expectedOffchainMessagesAfterPlaceOrder: []msgsender.Message{
				off_chain_updates.MustCreateOrderPlaceMessage(
					nil,
					PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order,
				).AddHeader(msgsender.MessageHeader{
					Key:   msgsender.TransactionHashHeaderKey,
					Value: tmhash.Sum(CheckTx_PlaceOrder_Bob_Num0_Id0_Sell5_Price10_GTB20.Tx),
				}),
				off_chain_updates.MustCreateOrderUpdateMessage(
					nil,
					PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order.OrderId,
					0,
				).AddHeader(msgsender.MessageHeader{
					Key:   msgsender.TransactionHashHeaderKey,
					Value: tmhash.Sum(CheckTx_PlaceOrder_Bob_Num0_Id0_Sell5_Price10_GTB20.Tx),
				}),
				off_chain_updates.MustCreateOrderPlaceMessage(
					nil,
					PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20.Order,
				).AddHeader(msgsender.MessageHeader{
					Key:   msgsender.TransactionHashHeaderKey,
					Value: tmhash.Sum(CheckTx_PlaceOrder_Alice_Num0_Id0_Buy6_Price10_GTB20.Tx),
				}),
				off_chain_updates.MustCreateOrderUpdateMessage(
					nil,
					PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order.OrderId,
					PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order.GetBaseQuantums(),
				).AddHeader(msgsender.MessageHeader{
					Key:   msgsender.TransactionHashHeaderKey,
					Value: tmhash.Sum(CheckTx_PlaceOrder_Alice_Num0_Id0_Buy6_Price10_GTB20.Tx),
				}),
				off_chain_updates.MustCreateOrderUpdateMessage(
					nil,
					PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20.Order.OrderId,
					PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order.GetBaseQuantums(),
				).AddHeader(msgsender.MessageHeader{
					Key:   msgsender.TransactionHashHeaderKey,
					Value: tmhash.Sum(CheckTx_PlaceOrder_Alice_Num0_Id0_Buy6_Price10_GTB20.Tx),
				}),
			},
			expectedOffchainMessagesInNextBlock: []msgsender.Message{
				off_chain_updates.MustCreateOrderUpdateMessage(
					nil,
					PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20.Order.OrderId,
					PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order.GetBaseQuantums(),
				),
			},
			expectedOnchainMessagesInNextBlock: []msgsender.Message{indexer_manager.CreateIndexerBlockEventMessage(
				&indexer_manager.IndexerTendermintBlock{
					Height: 2,
					Time:   ctx.BlockTime(),
					Events: []*indexer_manager.IndexerTendermintEvent{
						{
							Subtype: indexerevents.SubtypeSubaccountUpdate,
							Data: indexer_manager.GetB64EncodedEventMessage(
								indexerevents.NewSubaccountUpdateEvent(
									&constants.Alice_Num0,
									[]*satypes.PerpetualPosition{
										{
											PerpetualId: Clob_0.MustGetPerpetualId(),
											Quantums: dtypes.NewInt(int64(
												PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order.GetQuantums())),
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
									nil, // no funding payments
								),
							),
							OrderingWithinBlock: &indexer_manager.IndexerTendermintEvent_TransactionIndex{},
							EventIndex:          0,
							Version:             indexerevents.SubaccountUpdateEventVersion,
							DataBytes: indexer_manager.GetBytes(
								indexerevents.NewSubaccountUpdateEvent(
									&constants.Alice_Num0,
									[]*satypes.PerpetualPosition{
										{
											PerpetualId: Clob_0.MustGetPerpetualId(),
											Quantums: dtypes.NewInt(int64(
												PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order.GetQuantums())),
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
									nil, // no funding payments
								),
							),
						},
						{
							Subtype: indexerevents.SubtypeSubaccountUpdate,
							Data: indexer_manager.GetB64EncodedEventMessage(
								indexerevents.NewSubaccountUpdateEvent(
									&constants.Bob_Num0,
									[]*satypes.PerpetualPosition{
										{
											PerpetualId: Clob_0.MustGetPerpetualId(),
											Quantums: dtypes.NewInt(-int64(
												PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order.GetQuantums())),
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
									nil, // no funding payments
								),
							),
							OrderingWithinBlock: &indexer_manager.IndexerTendermintEvent_TransactionIndex{},
							EventIndex:          1,
							Version:             indexerevents.SubaccountUpdateEventVersion,
							DataBytes: indexer_manager.GetBytes(
								indexerevents.NewSubaccountUpdateEvent(
									&constants.Bob_Num0,
									[]*satypes.PerpetualPosition{
										{
											PerpetualId: Clob_0.MustGetPerpetualId(),
											Quantums: dtypes.NewInt(-int64(
												PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order.GetQuantums())),
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
									nil, // no funding payments
								),
							),
						},
						{
							Subtype: indexerevents.SubtypeOrderFill,
							Data: indexer_manager.GetB64EncodedEventMessage(
								indexerevents.NewOrderFillEvent(
									PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order,
									PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20.Order,
									PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order.GetBaseQuantums(),
									0, // Fees are 0 due to lost precision
									0,
									PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order.GetBaseQuantums(),
									PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order.GetBaseQuantums(),
								),
							),
							OrderingWithinBlock: &indexer_manager.IndexerTendermintEvent_TransactionIndex{},
							EventIndex:          2,
							Version:             indexerevents.OrderFillEventVersion,
							DataBytes: indexer_manager.GetBytes(
								indexerevents.NewOrderFillEvent(
									PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order,
									PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20.Order,
									PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order.GetBaseQuantums(),
									0, // Fees are 0 due to lost precision
									0,
									PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order.GetBaseQuantums(),
									PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order.GetBaseQuantums(),
								),
							),
						},
					},
					TxHashes: []string{string(lib.GetTxHash(testtx.MustGetTxBytes(&clobtypes.MsgProposedOperations{
						OperationsQueue: []clobtypes.OperationRaw{
							{
								Operation: &clobtypes.OperationRaw_ShortTermOrderPlacement{
									ShortTermOrderPlacement: CheckTx_PlaceOrder_Bob_Num0_Id0_Sell5_Price10_GTB20.Tx,
								},
							},
							{
								Operation: &clobtypes.OperationRaw_ShortTermOrderPlacement{
									ShortTermOrderPlacement: CheckTx_PlaceOrder_Alice_Num0_Id0_Buy6_Price10_GTB20.Tx,
								},
							},
							clobtestutils.NewMatchOperationRaw(
								&PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20.Order,
								[]clobtypes.MakerFill{
									{
										FillAmount: PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.
											Order.GetBaseQuantums().ToUint64(),
										MakerOrderId: PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order.OrderId,
									},
								},
							),
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

			ctx = tApp.InitChain()
			// Clear any messages produced prior to these checkTx calls.
			msgSender.Clear()
			for _, order := range tc.orders {
				for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(ctx, tApp.App, order) {
					resp := tApp.CheckTx(checkTx)
					require.Conditionf(t, resp.IsOK, "Expected CheckTx to succeed. Response: %+v", resp)
				}
			}
			require.ElementsMatch(
				t,
				tc.expectedOffchainMessagesAfterPlaceOrder,
				msgSender.GetOffchainMessages(),
			)
			require.ElementsMatch(t, tc.expectedOnchainMessagesAfterPlaceOrder, msgSender.GetOnchainMessages())

			// Clear the messages that we already matched prior to advancing to the next block.
			msgSender.Clear()
			ctx = tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{})
			require.ElementsMatch(t, tc.expectedOffchainMessagesInNextBlock, msgSender.GetOffchainMessages())
			require.ElementsMatch(t, tc.expectedOnchainMessagesInNextBlock, msgSender.GetOnchainMessages())
			for _, order := range tc.orders {
				if slices.Contains(tc.expectedOrdersFilled, order.Order.OrderId) {
					exists, fillAmount, _ := tApp.App.ClobKeeper.GetOrderFillAmount(
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

func TestShortTermOrderReplacements(t *testing.T) {
	type orderIdExpectations struct {
		shouldExistOnMemclob bool
		expectedOrder 	  clobtypes.Order
		expectedFillAmount uint64
	}
	type blockOrdersAndExpectations struct {
		ordersToPlace []clobtypes.MsgPlaceOrder
		orderIdsExpectations map[clobtypes.OrderId]orderIdExpectations
	}
	tests := map[string]struct {
		blocks []blockOrdersAndExpectations
	} {
		"Success: Replace in same block on same side": {
			blocks: []blockOrdersAndExpectations{
				{
					ordersToPlace: []clobtypes.MsgPlaceOrder{
						PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20,
						PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB21,
					},
					orderIdsExpectations: map[clobtypes.OrderId]orderIdExpectations{
						PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB21.Order.OrderId: {
							shouldExistOnMemclob: true,
							expectedOrder: PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB21.Order,
						},
					},
				},
			},
		},
		"Success: Replace in same block on opposite side": {
			blocks: []blockOrdersAndExpectations{
				{
					ordersToPlace: []clobtypes.MsgPlaceOrder{
						PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20,
						PlaceOrder_Alice_Num0_Id0_Clob0_Sell6_Price10_GTB21,
					},
					orderIdsExpectations: map[clobtypes.OrderId]orderIdExpectations{
						PlaceOrder_Alice_Num0_Id0_Clob0_Sell6_Price10_GTB21.Order.OrderId: {
							shouldExistOnMemclob: true,
							expectedOrder: PlaceOrder_Alice_Num0_Id0_Clob0_Sell6_Price10_GTB21.Order,
						},
					},
				},
			},
		},
		"Success: Replace in next block": {
			blocks: []blockOrdersAndExpectations{
				{
					ordersToPlace: []clobtypes.MsgPlaceOrder{
						PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20,
					},
					orderIdsExpectations: map[clobtypes.OrderId]orderIdExpectations{
						PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20.Order.OrderId: {
							shouldExistOnMemclob: true,
							expectedOrder: PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20.Order,
						},
					},
				},
				{
					ordersToPlace: []clobtypes.MsgPlaceOrder{
						PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB21,
					},
					orderIdsExpectations: map[clobtypes.OrderId]orderIdExpectations{
						PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB21.Order.OrderId: {
							shouldExistOnMemclob: true,
							expectedOrder: PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB21.Order,
						},
					},
				},
			},
		},
		"Fail: Replacement order has lower GTB than existing order": {
			blocks: []blockOrdersAndExpectations{
				{
					ordersToPlace: []clobtypes.MsgPlaceOrder{
						PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB21,
						PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20,
					},
					orderIdsExpectations: map[clobtypes.OrderId]orderIdExpectations{
						PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB21.Order.OrderId: {
							shouldExistOnMemclob: true,
							expectedOrder: PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB21.Order,
						},
					},
				},
			},
		},
		"Fail: Replacement order has equal GTB to existing order": {
			blocks: []blockOrdersAndExpectations{
				{
					ordersToPlace: []clobtypes.MsgPlaceOrder{
						PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20,
						PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB20,
					},
					orderIdsExpectations: map[clobtypes.OrderId]orderIdExpectations{
						PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20.Order.OrderId: {
							shouldExistOnMemclob: true,
							expectedOrder: PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20.Order,
						},
					},
				},
			},
		},
		"Success: Replacement order after partial match in same block": {
			blocks: []blockOrdersAndExpectations{
				{
					ordersToPlace: []clobtypes.MsgPlaceOrder{
						PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20,
						*clobtypes.NewMsgPlaceOrder(MustScaleOrder(
							clobtypes.Order{
								OrderId:      clobtypes.OrderId{SubaccountId: constants.Bob_Num0, ClientId: 0, ClobPairId: 0},
								Side:         clobtypes.Order_SIDE_SELL,
								Quantums:     3,
								Subticks:     10,
								GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
							},
							testapp.DefaultGenesis(),
						)),
						PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB21,
					},
					orderIdsExpectations: map[clobtypes.OrderId]orderIdExpectations{
						PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB21.Order.OrderId: {
							shouldExistOnMemclob: true,
							expectedOrder: PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB21.Order,
							expectedFillAmount: PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20.Order.Quantums / 2,
						},
					},
				},
			},
		},
		"Success: Replacement order increases size in next block after partial fill": {
			blocks: []blockOrdersAndExpectations{
				{
					ordersToPlace: []clobtypes.MsgPlaceOrder{
						PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20,
						*clobtypes.NewMsgPlaceOrder(MustScaleOrder(
							clobtypes.Order{
								OrderId:      clobtypes.OrderId{SubaccountId: constants.Bob_Num0, ClientId: 0, ClobPairId: 0},
								Side:         clobtypes.Order_SIDE_SELL,
								Quantums:     3,
								Subticks:     10,
								GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
							},
							testapp.DefaultGenesis(),
						)),
					},
					orderIdsExpectations: map[clobtypes.OrderId]orderIdExpectations{
						PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20.Order.OrderId: {
							shouldExistOnMemclob: true,
							expectedOrder: PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20.Order,
							expectedFillAmount: PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20.Order.Quantums / 2,
						},
					},
				},
				{
					ordersToPlace: []clobtypes.MsgPlaceOrder{
						PlaceOrder_Alice_Num0_Id0_Clob0_Buy7_Price10_GTB21,
					},
					orderIdsExpectations: map[clobtypes.OrderId]orderIdExpectations{
						PlaceOrder_Alice_Num0_Id0_Clob0_Buy7_Price10_GTB21.Order.OrderId: {
							shouldExistOnMemclob: true,
							expectedOrder: PlaceOrder_Alice_Num0_Id0_Clob0_Buy7_Price10_GTB21.Order,
							expectedFillAmount: PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20.Order.Quantums / 2,
						},
					},
				},
			},
		},
		"Success: Replacement order decreases size in next block after partial fill": {
			blocks: []blockOrdersAndExpectations{
				{
					ordersToPlace: []clobtypes.MsgPlaceOrder{
						PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20,
						*clobtypes.NewMsgPlaceOrder(MustScaleOrder(
							clobtypes.Order{
								OrderId:      clobtypes.OrderId{SubaccountId: constants.Bob_Num0, ClientId: 0, ClobPairId: 0},
								Side:         clobtypes.Order_SIDE_SELL,
								Quantums:     3,
								Subticks:     10,
								GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
							},
							testapp.DefaultGenesis(),
						)),
					},
					orderIdsExpectations: map[clobtypes.OrderId]orderIdExpectations{
						PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20.Order.OrderId: {
							shouldExistOnMemclob: true,
							expectedOrder: PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20.Order,
							expectedFillAmount: PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20.Order.Quantums / 2,
						},
					},
				},
				{
					ordersToPlace: []clobtypes.MsgPlaceOrder{
						PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB21,
					},
					orderIdsExpectations: map[clobtypes.OrderId]orderIdExpectations{
						PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB21.Order.OrderId: {
							shouldExistOnMemclob: true,
							expectedOrder: PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB21.Order,
							expectedFillAmount: PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20.Order.Quantums / 2,
						},
					},
				},
			},
		},
		"Fail: Replacement order attempts to decrease size such that the order would be fully filled": {
			blocks: []blockOrdersAndExpectations{
				{
					ordersToPlace: []clobtypes.MsgPlaceOrder{
						PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20,
						*clobtypes.NewMsgPlaceOrder(MustScaleOrder(
							clobtypes.Order{
								OrderId:      clobtypes.OrderId{SubaccountId: constants.Bob_Num0, ClientId: 0, ClobPairId: 0},
								Side:         clobtypes.Order_SIDE_SELL,
								Quantums:     3,
								Subticks:     10,
								GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
							},
							testapp.DefaultGenesis(),
						)),
					},
					orderIdsExpectations: map[clobtypes.OrderId]orderIdExpectations{
						PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20.Order.OrderId: {
							shouldExistOnMemclob: true,
							expectedOrder: PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20.Order,
							expectedFillAmount: PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20.Order.Quantums / 2,
						},
					},
				},
				{
					ordersToPlace: []clobtypes.MsgPlaceOrder{
						*clobtypes.NewMsgPlaceOrder(MustScaleOrder(
							clobtypes.Order{
								OrderId:      clobtypes.OrderId{SubaccountId: constants.Alice_Num0, ClientId: 0, ClobPairId: 0},
								Side:         clobtypes.Order_SIDE_BUY,
								Quantums:     3,
								Subticks:     10,
								GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
							},
							testapp.DefaultGenesis(),
						)),
					},
					orderIdsExpectations: map[clobtypes.OrderId]orderIdExpectations{
						PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20.Order.OrderId: {
							shouldExistOnMemclob: true,
							expectedOrder: PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20.Order,
							expectedFillAmount: PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20.Order.Quantums / 2,
						},
					},
				},
			},
		},
		"Fail: Replacement order attempts to decrease size below partially filled amount": {
			blocks: []blockOrdersAndExpectations{
				{
					ordersToPlace: []clobtypes.MsgPlaceOrder{
						PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20,
						*clobtypes.NewMsgPlaceOrder(MustScaleOrder(
							clobtypes.Order{
								OrderId:      clobtypes.OrderId{SubaccountId: constants.Bob_Num0, ClientId: 0, ClobPairId: 0},
								Side:         clobtypes.Order_SIDE_SELL,
								Quantums:     3,
								Subticks:     10,
								GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
							},
							testapp.DefaultGenesis(),
						)),
					},
					orderIdsExpectations: map[clobtypes.OrderId]orderIdExpectations{
						PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20.Order.OrderId: {
							shouldExistOnMemclob: true,
							expectedOrder: PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20.Order,
							expectedFillAmount: PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20.Order.Quantums / 2,
						},
					},
				},
				{
					ordersToPlace: []clobtypes.MsgPlaceOrder{
						*clobtypes.NewMsgPlaceOrder(MustScaleOrder(
							clobtypes.Order{
								OrderId:      clobtypes.OrderId{SubaccountId: constants.Alice_Num0, ClientId: 0, ClobPairId: 0},
								Side:         clobtypes.Order_SIDE_BUY,
								Quantums:     2,
								Subticks:     10,
								GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
							},
							testapp.DefaultGenesis(),
						)),
					},
					orderIdsExpectations: map[clobtypes.OrderId]orderIdExpectations{
						PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20.Order.OrderId: {
							shouldExistOnMemclob: true,
							expectedOrder: PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20.Order,
							expectedFillAmount: PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20.Order.Quantums / 2,
						},
					},
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder().Build()
			ctx := tApp.InitChain()

			for i, block := range tc.blocks {
				for _, order := range block.ordersToPlace {
					for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(ctx, tApp.App, order) {
						tApp.CheckTx(checkTx)
					}
				}

				ctx = tApp.AdvanceToBlock(uint32(i + 2), testapp.AdvanceToBlockOptions{})

				for orderId, expectations := range block.orderIdsExpectations {
					order, exists := tApp.App.ClobKeeper.MemClob.GetOrder(ctx, orderId)
					require.Equal(t, expectations.shouldExistOnMemclob, exists)
					if expectations.shouldExistOnMemclob {
						require.Equal(t, expectations.expectedOrder, order)
					}
					_, fillAmount, _ := tApp.App.ClobKeeper.GetOrderFillAmount(ctx, orderId)
					require.Equal(t, expectations.expectedFillAmount, uint64(fillAmount))
				}
			}
		})
	}
}
