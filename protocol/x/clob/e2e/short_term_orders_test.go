package clob_test

import (
	"math/big"
	"testing"

	"github.com/cometbft/cometbft/crypto/tmhash"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/slices"

	"github.com/dydxprotocol/v4-chain/protocol/indexer"
	indexerevents "github.com/dydxprotocol/v4-chain/protocol/indexer/events"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/msgsender"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/off_chain_updates"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	clobtestutils "github.com/dydxprotocol/v4-chain/protocol/testutil/clob"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	testmsgs "github.com/dydxprotocol/v4-chain/protocol/testutil/msgs"
	testtx "github.com/dydxprotocol/v4-chain/protocol/testutil/tx"
	testutil "github.com/dydxprotocol/v4-chain/protocol/testutil/util"
	assettypes "github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

func TestPlaceOrder(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()

	aliceSubaccount := tApp.App.SubaccountsKeeper.GetSubaccount(ctx, constants.Alice_Num0)
	bobSubaccount := tApp.App.SubaccountsKeeper.GetSubaccount(ctx, constants.Bob_Num0)

	CheckTx_PlaceOrder_Alice_Num0_Id0_Buy5_Price10_GTB20 := testapp.MustMakeCheckTx(
		ctx,
		tApp.App,
		testapp.MustMakeCheckTxOptions{
			AccAddressForSigning: constants.Alice_Num0.Owner,
		},
		&PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB20,
	)
	CheckTx_PlaceOrder_Alice_Num0_Id0_Buy6_Price10_GTB20 := testapp.MustMakeCheckTx(
		ctx,
		tApp.App,
		testapp.MustMakeCheckTxOptions{
			AccAddressForSigning: constants.Alice_Num0.Owner,
		},
		&PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20,
	)
	CheckTx_PlaceOrder_Bob_Num0_Id0_Sell5_Price10_GTB20 := testapp.MustMakeCheckTx(
		ctx,
		tApp.App,
		testapp.MustMakeCheckTxOptions{
			AccAddressForSigning: constants.Bob_Num0.Owner,
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
					ctx,
					PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB20.Order,
				).AddHeader(msgsender.MessageHeader{
					Key:   msgsender.TransactionHashHeaderKey,
					Value: tmhash.Sum(CheckTx_PlaceOrder_Alice_Num0_Id0_Buy5_Price10_GTB20.Tx),
				}),
				off_chain_updates.MustCreateOrderUpdateMessage(
					ctx,
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
					ctx,
					PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB20.Order,
				).AddHeader(msgsender.MessageHeader{
					Key:   msgsender.TransactionHashHeaderKey,
					Value: tmhash.Sum(CheckTx_PlaceOrder_Alice_Num0_Id0_Buy5_Price10_GTB20.Tx),
				}),
				off_chain_updates.MustCreateOrderUpdateMessage(
					ctx,
					PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB20.Order.OrderId,
					0,
				).AddHeader(msgsender.MessageHeader{
					Key:   msgsender.TransactionHashHeaderKey,
					Value: tmhash.Sum(CheckTx_PlaceOrder_Alice_Num0_Id0_Buy5_Price10_GTB20.Tx),
				}),
				off_chain_updates.MustCreateOrderPlaceMessage(
					ctx,
					PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order,
				).AddHeader(msgsender.MessageHeader{
					Key:   msgsender.TransactionHashHeaderKey,
					Value: tmhash.Sum(CheckTx_PlaceOrder_Bob_Num0_Id0_Sell5_Price10_GTB20.Tx),
				}),
				off_chain_updates.MustCreateOrderUpdateMessage(
					ctx,
					PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB20.Order.OrderId,
					PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB20.Order.GetBaseQuantums(),
				).AddHeader(msgsender.MessageHeader{
					Key:   msgsender.TransactionHashHeaderKey,
					Value: tmhash.Sum(CheckTx_PlaceOrder_Bob_Num0_Id0_Sell5_Price10_GTB20.Tx),
				}),
				off_chain_updates.MustCreateOrderUpdateMessage(
					ctx,
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
							Subtype:             indexerevents.SubtypeSubaccountUpdate,
							OrderingWithinBlock: &indexer_manager.IndexerTendermintEvent_TransactionIndex{},
							EventIndex:          0,
							Version:             indexerevents.SubaccountUpdateEventVersion,
							DataBytes: indexer_manager.GetBytes(
								indexerevents.NewSubaccountUpdateEvent(
									&constants.Bob_Num0,
									[]*satypes.PerpetualPosition{
										testutil.CreateSinglePerpetualPosition(
											Clob_0.MustGetPerpetualId(),
											big.NewInt(-int64(
												PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB20.Order.GetQuantums())),
											big.NewInt(0),
											big.NewInt(0),
										),
									},
									// Maker fees calculate to 0 so asset position doesn't change.
									[]*satypes.AssetPosition{
										testutil.CreateSingleAssetPosition(
											assettypes.AssetUsdc.Id,
											bobSubaccount.GetUsdcPosition(),
										),
									},
									nil, // no funding payments
								),
							),
						},
						{
							Subtype:             indexerevents.SubtypeSubaccountUpdate,
							OrderingWithinBlock: &indexer_manager.IndexerTendermintEvent_TransactionIndex{},
							EventIndex:          1,
							Version:             indexerevents.SubaccountUpdateEventVersion,
							DataBytes: indexer_manager.GetBytes(
								indexerevents.NewSubaccountUpdateEvent(
									&constants.Alice_Num0,
									[]*satypes.PerpetualPosition{
										testutil.CreateSinglePerpetualPosition(
											Clob_0.MustGetPerpetualId(),
											big.NewInt(int64(
												PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB20.Order.GetQuantums())),
											big.NewInt(0),
											big.NewInt(0),
										),
									},
									// Taker fees calculate to 0 so asset position doesn't change.
									[]*satypes.AssetPosition{
										testutil.CreateSingleAssetPosition(
											assettypes.AssetUsdc.Id,
											aliceSubaccount.GetUsdcPosition(),
										),
									},
									nil, // no funding payments
								),
							),
						},
						{
							Subtype:             indexerevents.SubtypeOrderFill,
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
									0,
									0,
									PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB20.Order.GetBaseQuantums(),
									PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB20.Order.GetBaseQuantums(),
									big.NewInt(0),
									0,
									0,
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
					ctx,
					PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20.Order,
				).AddHeader(msgsender.MessageHeader{
					Key:   msgsender.TransactionHashHeaderKey,
					Value: tmhash.Sum(CheckTx_PlaceOrder_Alice_Num0_Id0_Buy6_Price10_GTB20.Tx),
				}),
				off_chain_updates.MustCreateOrderUpdateMessage(
					ctx,
					PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20.Order.OrderId,
					0,
				).AddHeader(msgsender.MessageHeader{
					Key:   msgsender.TransactionHashHeaderKey,
					Value: tmhash.Sum(CheckTx_PlaceOrder_Alice_Num0_Id0_Buy6_Price10_GTB20.Tx),
				}),
				off_chain_updates.MustCreateOrderPlaceMessage(
					ctx,
					PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order,
				).AddHeader(msgsender.MessageHeader{
					Key:   msgsender.TransactionHashHeaderKey,
					Value: tmhash.Sum(CheckTx_PlaceOrder_Bob_Num0_Id0_Sell5_Price10_GTB20.Tx),
				}),
				off_chain_updates.MustCreateOrderUpdateMessage(
					ctx,
					PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20.Order.OrderId,
					PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order.GetBaseQuantums(),
				).AddHeader(msgsender.MessageHeader{
					Key:   msgsender.TransactionHashHeaderKey,
					Value: tmhash.Sum(CheckTx_PlaceOrder_Bob_Num0_Id0_Sell5_Price10_GTB20.Tx),
				}),
				off_chain_updates.MustCreateOrderUpdateMessage(
					ctx,
					PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order.OrderId,
					PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order.GetBaseQuantums(),
				).AddHeader(msgsender.MessageHeader{
					Key:   msgsender.TransactionHashHeaderKey,
					Value: tmhash.Sum(CheckTx_PlaceOrder_Bob_Num0_Id0_Sell5_Price10_GTB20.Tx),
				}),
			},
			expectedOffchainMessagesInNextBlock: []msgsender.Message{
				off_chain_updates.MustCreateOrderUpdateMessage(
					ctx,
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
							Subtype:             indexerevents.SubtypeSubaccountUpdate,
							OrderingWithinBlock: &indexer_manager.IndexerTendermintEvent_TransactionIndex{},
							EventIndex:          0,
							Version:             indexerevents.SubaccountUpdateEventVersion,
							DataBytes: indexer_manager.GetBytes(
								indexerevents.NewSubaccountUpdateEvent(
									&constants.Bob_Num0,
									[]*satypes.PerpetualPosition{
										testutil.CreateSinglePerpetualPosition(
											Clob_0.MustGetPerpetualId(),
											big.NewInt(-int64(
												PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order.GetQuantums())),
											big.NewInt(0),
											big.NewInt(0),
										),
									},
									// Maker fees calculate to 0 so asset position doesn't change.
									[]*satypes.AssetPosition{
										testutil.CreateSingleAssetPosition(
											assettypes.AssetUsdc.Id,
											bobSubaccount.GetUsdcPosition(),
										),
									},
									nil, // no funding payments
								),
							),
						},
						{
							Subtype:             indexerevents.SubtypeSubaccountUpdate,
							OrderingWithinBlock: &indexer_manager.IndexerTendermintEvent_TransactionIndex{},
							EventIndex:          1,
							Version:             indexerevents.SubaccountUpdateEventVersion,
							DataBytes: indexer_manager.GetBytes(
								indexerevents.NewSubaccountUpdateEvent(
									&constants.Alice_Num0,
									[]*satypes.PerpetualPosition{
										testutil.CreateSinglePerpetualPosition(
											Clob_0.MustGetPerpetualId(),
											big.NewInt(int64(
												PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order.GetQuantums())),
											big.NewInt(0),
											big.NewInt(0),
										),
									},
									// Taker fees calculate to 0 so asset position doesn't change.
									[]*satypes.AssetPosition{
										testutil.CreateSingleAssetPosition(
											assettypes.AssetUsdc.Id,
											aliceSubaccount.GetUsdcPosition(),
										),
									},
									nil, // no funding payments
								),
							),
						},
						{
							Subtype:             indexerevents.SubtypeOrderFill,
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
									0,
									0,
									PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order.GetBaseQuantums(),
									PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order.GetBaseQuantums(),
									big.NewInt(0),
									0,
									0,
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
					ctx,
					PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order,
				).AddHeader(msgsender.MessageHeader{
					Key:   msgsender.TransactionHashHeaderKey,
					Value: tmhash.Sum(CheckTx_PlaceOrder_Bob_Num0_Id0_Sell5_Price10_GTB20.Tx),
				}),
				off_chain_updates.MustCreateOrderUpdateMessage(
					ctx,
					PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order.OrderId,
					0,
				).AddHeader(msgsender.MessageHeader{
					Key:   msgsender.TransactionHashHeaderKey,
					Value: tmhash.Sum(CheckTx_PlaceOrder_Bob_Num0_Id0_Sell5_Price10_GTB20.Tx),
				}),
				off_chain_updates.MustCreateOrderPlaceMessage(
					ctx,
					PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20.Order,
				).AddHeader(msgsender.MessageHeader{
					Key:   msgsender.TransactionHashHeaderKey,
					Value: tmhash.Sum(CheckTx_PlaceOrder_Alice_Num0_Id0_Buy6_Price10_GTB20.Tx),
				}),
				off_chain_updates.MustCreateOrderUpdateMessage(
					ctx,
					PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order.OrderId,
					PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order.GetBaseQuantums(),
				).AddHeader(msgsender.MessageHeader{
					Key:   msgsender.TransactionHashHeaderKey,
					Value: tmhash.Sum(CheckTx_PlaceOrder_Alice_Num0_Id0_Buy6_Price10_GTB20.Tx),
				}),
				off_chain_updates.MustCreateOrderUpdateMessage(
					ctx,
					PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20.Order.OrderId,
					PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order.GetBaseQuantums(),
				).AddHeader(msgsender.MessageHeader{
					Key:   msgsender.TransactionHashHeaderKey,
					Value: tmhash.Sum(CheckTx_PlaceOrder_Alice_Num0_Id0_Buy6_Price10_GTB20.Tx),
				}),
			},
			expectedOffchainMessagesInNextBlock: []msgsender.Message{
				off_chain_updates.MustCreateOrderUpdateMessage(
					ctx,
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
							Subtype:             indexerevents.SubtypeSubaccountUpdate,
							OrderingWithinBlock: &indexer_manager.IndexerTendermintEvent_TransactionIndex{},
							EventIndex:          0,
							Version:             indexerevents.SubaccountUpdateEventVersion,
							DataBytes: indexer_manager.GetBytes(
								indexerevents.NewSubaccountUpdateEvent(
									&constants.Alice_Num0,
									[]*satypes.PerpetualPosition{
										testutil.CreateSinglePerpetualPosition(
											Clob_0.MustGetPerpetualId(),
											big.NewInt(int64(
												PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order.GetQuantums())),
											big.NewInt(0),
											big.NewInt(0),
										),
									},
									// Taker fees calculate to 0 so asset position doesn't change.
									[]*satypes.AssetPosition{
										testutil.CreateSingleAssetPosition(
											assettypes.AssetUsdc.Id,
											aliceSubaccount.GetUsdcPosition(),
										),
									},
									nil, // no funding payments
								),
							),
						},
						{
							Subtype:             indexerevents.SubtypeSubaccountUpdate,
							OrderingWithinBlock: &indexer_manager.IndexerTendermintEvent_TransactionIndex{},
							EventIndex:          1,
							Version:             indexerevents.SubaccountUpdateEventVersion,
							DataBytes: indexer_manager.GetBytes(
								indexerevents.NewSubaccountUpdateEvent(
									&constants.Bob_Num0,
									[]*satypes.PerpetualPosition{
										testutil.CreateSinglePerpetualPosition(
											Clob_0.MustGetPerpetualId(),
											big.NewInt(-int64(
												PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order.GetQuantums())),
											big.NewInt(0),
											big.NewInt(0),
										),
									},
									// Maker fees calculate to 0 so asset position doesn't change.
									[]*satypes.AssetPosition{
										testutil.CreateSingleAssetPosition(
											assettypes.AssetUsdc.Id,
											bobSubaccount.GetUsdcPosition(),
										),
									},
									nil, // no funding payments
								),
							),
						},
						{
							Subtype:             indexerevents.SubtypeOrderFill,
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
									0,
									0,
									PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order.GetBaseQuantums(),
									PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order.GetBaseQuantums(),
									big.NewInt(0),
									0,
									0,
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
			msgSender := msgsender.NewIndexerMessageSenderInMemoryCollector()
			appOpts := map[string]interface{}{
				indexer.MsgSenderInstanceForTest: msgSender,
			}
			tApp = testapp.NewTestAppBuilder(t).WithAppOptions(appOpts).Build()
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
	order := PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20
	ioc_replacement := order
	ioc_replacement.Order.GoodTilOneof = &clobtypes.Order_GoodTilBlock{GoodTilBlock: 21}
	ioc_replacement.Order.TimeInForce = clobtypes.Order_TIME_IN_FORCE_IOC

	type orderIdExpectations struct {
		shouldExistOnMemclob bool
		expectedOrder        clobtypes.Order
		expectedFillAmount   uint64
	}
	type blockOrdersAndExpectations struct {
		ordersToPlace        []clobtypes.MsgPlaceOrder
		orderIdsExpectations map[clobtypes.OrderId]orderIdExpectations
	}
	tests := map[string]struct {
		blocks []blockOrdersAndExpectations
	}{
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
							expectedOrder:        PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB21.Order,
						},
					},
				},
			},
		},
		"Fail: Replace in same block on opposite side": {
			blocks: []blockOrdersAndExpectations{
				{
					ordersToPlace: []clobtypes.MsgPlaceOrder{
						PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20,
						PlaceOrder_Alice_Num0_Id0_Clob0_Sell6_Price10_GTB21,
					},
					orderIdsExpectations: map[clobtypes.OrderId]orderIdExpectations{
						PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20.Order.OrderId: {
							shouldExistOnMemclob: true,
							expectedOrder:        PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20.Order,
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
							expectedOrder:        PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20.Order,
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
							expectedOrder:        PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB21.Order,
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
							expectedOrder:        PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB21.Order,
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
							expectedOrder:        PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20.Order,
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
						*clobtypes.NewMsgPlaceOrder(testapp.MustScaleOrder(
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
							expectedOrder:        PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB21.Order,
							expectedFillAmount:   PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20.Order.Quantums / 2,
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
						*clobtypes.NewMsgPlaceOrder(testapp.MustScaleOrder(
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
							expectedOrder:        PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20.Order,
							expectedFillAmount:   PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20.Order.Quantums / 2,
						},
						PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20.Order.OrderId: {
							shouldExistOnMemclob: true,
							expectedOrder:        PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20.Order,
							expectedFillAmount:   PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20.Order.Quantums / 2,
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
							expectedOrder:        PlaceOrder_Alice_Num0_Id0_Clob0_Buy7_Price10_GTB21.Order,
							expectedFillAmount:   PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20.Order.Quantums / 2,
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
						*clobtypes.NewMsgPlaceOrder(testapp.MustScaleOrder(
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
							expectedOrder:        PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20.Order,
							expectedFillAmount:   PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20.Order.Quantums / 2,
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
							expectedOrder:        PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB21.Order,
							expectedFillAmount:   PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20.Order.Quantums / 2,
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
						*clobtypes.NewMsgPlaceOrder(testapp.MustScaleOrder(
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
							expectedOrder:        PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20.Order,
							expectedFillAmount:   PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20.Order.Quantums / 2,
						},
					},
				},
				{
					ordersToPlace: []clobtypes.MsgPlaceOrder{
						*clobtypes.NewMsgPlaceOrder(testapp.MustScaleOrder(
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
							expectedOrder:        PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20.Order,
							expectedFillAmount:   PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20.Order.Quantums / 2,
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
						*clobtypes.NewMsgPlaceOrder(testapp.MustScaleOrder(
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
							expectedOrder:        PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20.Order,
							expectedFillAmount:   PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20.Order.Quantums / 2,
						},
					},
				},
				{
					ordersToPlace: []clobtypes.MsgPlaceOrder{
						*clobtypes.NewMsgPlaceOrder(testapp.MustScaleOrder(
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
							expectedOrder:        PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20.Order,
							expectedFillAmount:   PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20.Order.Quantums / 2,
						},
					},
				},
			},
		},
		"Success: Replacing order with IOC which does not fully match results in order being removed from the book": {
			blocks: []blockOrdersAndExpectations{
				{
					ordersToPlace: []clobtypes.MsgPlaceOrder{
						PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20,
						ioc_replacement,
					},
					orderIdsExpectations: map[clobtypes.OrderId]orderIdExpectations{
						PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20.Order.OrderId: {
							shouldExistOnMemclob: true,
							expectedOrder:        PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20.Order,
						},
					},
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).Build()
			ctx := tApp.InitChain()

			for i, block := range tc.blocks {
				for _, order := range block.ordersToPlace {
					for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(ctx, tApp.App, order) {
						tApp.CheckTx(checkTx)
					}
				}

				for orderId, expectations := range block.orderIdsExpectations {
					order, exists := tApp.App.ClobKeeper.MemClob.GetOrder(orderId)
					require.Equal(t, expectations.shouldExistOnMemclob, exists)
					if expectations.shouldExistOnMemclob {
						require.Equal(t, expectations.expectedOrder, order)
					}
					_, fillAmount, _ := tApp.App.ClobKeeper.GetOrderFillAmount(ctx, orderId)
					require.Equal(t, expectations.expectedFillAmount, uint64(fillAmount))
				}

				ctx = tApp.AdvanceToBlock(uint32(i+2), testapp.AdvanceToBlockOptions{})
			}
		})
	}
}

func TestCancelShortTermOrder(t *testing.T) {
	tests := map[string]struct {
		firstBlockOrders   []clobtypes.MsgPlaceOrder
		firstBlockCancels  []clobtypes.MsgCancelOrder
		secondBlockOrders  []clobtypes.MsgPlaceOrder
		secondBlockCancels []clobtypes.MsgCancelOrder

		expectedOrderIdsInMemclob          map[clobtypes.OrderId]bool
		expectedCancelExpirationsInMemclob map[clobtypes.OrderId]uint32
		expectedOrderFillAmounts           map[clobtypes.OrderId]uint64
	}{
		"Cancel unfilled short term order": {
			firstBlockOrders: []clobtypes.MsgPlaceOrder{
				PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB5,
			},
			secondBlockCancels: []clobtypes.MsgCancelOrder{
				CancelOrder_Alice_Num0_Id0_Clob0_GTB5,
			},

			expectedOrderIdsInMemclob: map[clobtypes.OrderId]bool{
				PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB5.Order.OrderId: false,
			},
			expectedCancelExpirationsInMemclob: map[clobtypes.OrderId]uint32{
				CancelOrder_Alice_Num0_Id0_Clob0_GTB5.OrderId: 5,
			},
			expectedOrderFillAmounts: map[clobtypes.OrderId]uint64{
				PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB5.Order.OrderId: 0,
			},
		},
		"Cancel partially filled short term order in same block": {
			firstBlockOrders: []clobtypes.MsgPlaceOrder{
				PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB5,
				*clobtypes.NewMsgPlaceOrder(testapp.MustScaleOrder(
					clobtypes.Order{
						OrderId:      clobtypes.OrderId{SubaccountId: constants.Bob_Num0, ClientId: 0, ClobPairId: 0},
						Side:         clobtypes.Order_SIDE_SELL,
						Quantums:     4,
						Subticks:     10,
						GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
					},
					testapp.DefaultGenesis(),
				)),
			},
			firstBlockCancels: []clobtypes.MsgCancelOrder{
				CancelOrder_Alice_Num0_Id0_Clob0_GTB5,
			},

			expectedOrderIdsInMemclob: map[clobtypes.OrderId]bool{
				PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB5.Order.OrderId: false,
			},
			expectedCancelExpirationsInMemclob: map[clobtypes.OrderId]uint32{
				CancelOrder_Alice_Num0_Id0_Clob0_GTB5.OrderId: 5,
			},
			expectedOrderFillAmounts: map[clobtypes.OrderId]uint64{
				PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB5.Order.OrderId: 40,
			},
		},
		"Cancel partially filled short term order in next block": {
			firstBlockOrders: []clobtypes.MsgPlaceOrder{
				PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB5,
				*clobtypes.NewMsgPlaceOrder(testapp.MustScaleOrder(
					clobtypes.Order{
						OrderId:      clobtypes.OrderId{SubaccountId: constants.Bob_Num0, ClientId: 0, ClobPairId: 0},
						Side:         clobtypes.Order_SIDE_SELL,
						Quantums:     4,
						Subticks:     10,
						GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
					},
					testapp.DefaultGenesis(),
				)),
			},
			secondBlockCancels: []clobtypes.MsgCancelOrder{
				CancelOrder_Alice_Num0_Id0_Clob0_GTB5,
			},

			expectedOrderIdsInMemclob: map[clobtypes.OrderId]bool{
				PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB5.Order.OrderId: false,
			},
			expectedCancelExpirationsInMemclob: map[clobtypes.OrderId]uint32{
				CancelOrder_Alice_Num0_Id0_Clob0_GTB5.OrderId: 5,
			},
			expectedOrderFillAmounts: map[clobtypes.OrderId]uint64{
				PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB5.Order.OrderId: 40,
			},
		},
		"Cancel succeeds for fully-filled order": {
			firstBlockOrders: []clobtypes.MsgPlaceOrder{
				PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB5,
				PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20,
			},
			secondBlockCancels: []clobtypes.MsgCancelOrder{
				CancelOrder_Alice_Num0_Id0_Clob0_GTB5,
			},

			expectedOrderIdsInMemclob: map[clobtypes.OrderId]bool{
				PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB5.Order.OrderId: false,
			},
			expectedCancelExpirationsInMemclob: map[clobtypes.OrderId]uint32{
				CancelOrder_Alice_Num0_Id0_Clob0_GTB5.OrderId: 5,
			},
			expectedOrderFillAmounts: map[clobtypes.OrderId]uint64{
				PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB5.Order.OrderId: 50,
			},
		},
		"Cancel with GTB < existing order GTB does not remove order from memclob": {
			firstBlockOrders: []clobtypes.MsgPlaceOrder{
				PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB20,
			},
			secondBlockCancels: []clobtypes.MsgCancelOrder{
				CancelOrder_Alice_Num0_Id0_Clob0_GTB5,
			},

			expectedOrderIdsInMemclob: map[clobtypes.OrderId]bool{
				PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB20.Order.OrderId: true,
			},
			expectedCancelExpirationsInMemclob: map[clobtypes.OrderId]uint32{
				CancelOrder_Alice_Num0_Id0_Clob0_GTB5.OrderId: 5,
			},
		},
		"Cancel with GTB < existing cancel GTB is not placed on memclob": {
			firstBlockCancels: []clobtypes.MsgCancelOrder{
				CancelOrder_Alice_Num0_Id0_Clob0_GTB5,
				*clobtypes.NewMsgCancelOrderShortTerm(
					clobtypes.OrderId{
						SubaccountId: constants.Alice_Num0,
						ClientId:     0,
						ClobPairId:   0,
					},
					4,
				),
			},

			expectedCancelExpirationsInMemclob: map[clobtypes.OrderId]uint32{
				CancelOrder_Alice_Num0_Id0_Clob0_GTB5.OrderId: 5,
			},
		},
		"Cancel with GTB > existing cancel GTB is placed on memclob": {
			firstBlockCancels: []clobtypes.MsgCancelOrder{
				CancelOrder_Alice_Num0_Id0_Clob0_GTB5,
				*clobtypes.NewMsgCancelOrderShortTerm(
					clobtypes.OrderId{
						SubaccountId: constants.Alice_Num0,
						ClientId:     0,
						ClobPairId:   0,
					},
					6,
				),
			},

			expectedCancelExpirationsInMemclob: map[clobtypes.OrderId]uint32{
				CancelOrder_Alice_Num0_Id0_Clob0_GTB5.OrderId: 6,
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).Build()
			ctx := tApp.InitChain()

			// Place first block orders and cancels
			for _, order := range tc.firstBlockOrders {
				for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(ctx, tApp.App, order) {
					resp := tApp.CheckTx(checkTx)
					require.Conditionf(t, resp.IsOK, "Expected CheckTx to succeed. Response: %+v", resp)
				}
			}
			for _, cancel := range tc.firstBlockCancels {
				for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(ctx, tApp.App, cancel) {
					tApp.CheckTx(checkTx)
				}
			}

			// Advance block
			ctx = tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{})

			// Place second block orders and cancels
			for _, order := range tc.secondBlockOrders {
				for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(ctx, tApp.App, order) {
					resp := tApp.CheckTx(checkTx)
					require.Conditionf(t, resp.IsOK, "Expected CheckTx to succeed. Response: %+v", resp)
				}
			}
			for _, orderCancel := range tc.secondBlockCancels {
				for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(ctx, tApp.App, orderCancel) {
					resp := tApp.CheckTx(checkTx)
					require.Conditionf(t, resp.IsOK, "Expected CheckTx to succeed. Response: %+v", resp)
				}
			}

			// Verify expectations
			for orderId, shouldHaveOrder := range tc.expectedOrderIdsInMemclob {
				_, exists := tApp.App.ClobKeeper.MemClob.GetOrder(orderId)
				require.Equal(t, shouldHaveOrder, exists)
			}
			for orderId, expectedCancelExpirationBlock := range tc.expectedCancelExpirationsInMemclob {
				cancelExpirationBlock, exists := tApp.App.ClobKeeper.MemClob.GetCancelOrder(orderId)
				require.True(t, exists)
				require.Equal(t, expectedCancelExpirationBlock, cancelExpirationBlock)
			}
			for orderId, expectedFillAmount := range tc.expectedOrderFillAmounts {
				_, fillAmount, _ := tApp.App.ClobKeeper.GetOrderFillAmount(ctx, orderId)
				require.Equal(t, expectedFillAmount, fillAmount.ToUint64())
			}
		})
	}
}

func TestShortTermAdvancedOrders(t *testing.T) {
	tests := map[string]struct {
		blocks []testmsgs.TestBlockWithMsgs

		expectedOrderIdsInMemclob map[clobtypes.OrderId]bool
		expectedOrderFillAmounts  map[clobtypes.OrderId]uint64
	}{
		"IOC sell fully matches": {
			blocks: []testmsgs.TestBlockWithMsgs{
				{
					Block: 2,
					Msgs: []testmsgs.TestSdkMsg{
						{
							Msg: clobtypes.NewMsgPlaceOrder(
								testapp.MustScaleOrder(
									constants.Order_Bob_Num0_Id11_Clob1_Buy5_Price40_GTB20,
									testapp.DefaultGenesis(),
								),
							),
							ExpectedIsOk: true,
						},
						{
							Msg: clobtypes.NewMsgPlaceOrder(
								testapp.MustScaleOrder(
									constants.Order_Alice_Num0_Id1_Clob1_Sell5_Price15_GTB20_IOC,
									testapp.DefaultGenesis(),
								),
							),
							ExpectedIsOk: true,
						},
					},
				},
			},
			expectedOrderIdsInMemclob: map[clobtypes.OrderId]bool{
				constants.Order_Bob_Num0_Id11_Clob1_Buy5_Price40_GTB20.OrderId:       false,
				constants.Order_Alice_Num0_Id1_Clob1_Sell5_Price15_GTB20_IOC.OrderId: false,
			},
			expectedOrderFillAmounts: map[clobtypes.OrderId]uint64{
				constants.Order_Bob_Num0_Id11_Clob1_Buy5_Price40_GTB20.OrderId:       5000, // full size of scaled orders
				constants.Order_Alice_Num0_Id1_Clob1_Sell5_Price15_GTB20_IOC.OrderId: 5000,
			},
		},
		"IOC buy fully matches": {
			blocks: []testmsgs.TestBlockWithMsgs{
				{
					Block: 2,
					Msgs: []testmsgs.TestSdkMsg{
						{
							Msg: clobtypes.NewMsgPlaceOrder(
								testapp.MustScaleOrder(
									constants.Order_Bob_Num0_Id11_Clob1_Sell5_Price15_GTB20,
									testapp.DefaultGenesis(),
								),
							),
							ExpectedIsOk: true,
						},
						{
							Msg: clobtypes.NewMsgPlaceOrder(
								testapp.MustScaleOrder(
									constants.Order_Alice_Num0_Id1_Clob1_Buy5_Price15_GTB20_IOC,
									testapp.DefaultGenesis(),
								),
							),
							ExpectedIsOk: true,
						},
					},
				},
			},
			expectedOrderIdsInMemclob: map[clobtypes.OrderId]bool{
				constants.Order_Bob_Num0_Id11_Clob1_Sell5_Price15_GTB20.OrderId:     false,
				constants.Order_Alice_Num0_Id1_Clob1_Buy5_Price15_GTB20_IOC.OrderId: false,
			},
			expectedOrderFillAmounts: map[clobtypes.OrderId]uint64{
				constants.Order_Bob_Num0_Id11_Clob1_Sell5_Price15_GTB20.OrderId:     5000,
				constants.Order_Alice_Num0_Id1_Clob1_Buy5_Price15_GTB20_IOC.OrderId: 5000,
			},
		},
		"IOC sell partially matches and is not placed on the book": {
			blocks: []testmsgs.TestBlockWithMsgs{
				{
					Block: 2,
					Msgs: []testmsgs.TestSdkMsg{
						{
							Msg: clobtypes.NewMsgPlaceOrder(
								testapp.MustScaleOrder(
									constants.Order_Bob_Num0_Id11_Clob1_Buy5_Price40_GTB20,
									testapp.DefaultGenesis(),
								),
							),
							ExpectedIsOk: true,
						},
						{
							Msg: clobtypes.NewMsgPlaceOrder(
								testapp.MustScaleOrder(
									constants.Order_Alice_Num0_Id1_Clob1_Sell10_Price15_GTB20_IOC,
									testapp.DefaultGenesis(),
								),
							),
							ExpectedIsOk: true,
						},
					},
				},
			},
			expectedOrderIdsInMemclob: map[clobtypes.OrderId]bool{
				constants.Order_Bob_Num0_Id11_Clob1_Buy5_Price40_GTB20.OrderId:        false,
				constants.Order_Alice_Num0_Id1_Clob1_Sell10_Price15_GTB20_IOC.OrderId: false,
			},
			expectedOrderFillAmounts: map[clobtypes.OrderId]uint64{
				constants.Order_Bob_Num0_Id11_Clob1_Buy5_Price40_GTB20.OrderId:        5000,
				constants.Order_Alice_Num0_Id1_Clob1_Sell10_Price15_GTB20_IOC.OrderId: 5000,
			},
		},
		"IOC buy partially matches and is not placed on the book": {
			blocks: []testmsgs.TestBlockWithMsgs{
				{
					Block: 2,
					Msgs: []testmsgs.TestSdkMsg{
						{
							Msg: clobtypes.NewMsgPlaceOrder(
								testapp.MustScaleOrder(
									constants.Order_Bob_Num0_Id11_Clob1_Sell5_Price15_GTB20,
									testapp.DefaultGenesis(),
								),
							),
							ExpectedIsOk: true,
						},
						{
							Msg: clobtypes.NewMsgPlaceOrder(
								testapp.MustScaleOrder(
									constants.Order_Alice_Num0_Id1_Clob1_Buy10_Price15_GTB20_IOC,
									testapp.DefaultGenesis(),
								),
							),
							ExpectedIsOk: true,
						},
					},
				},
			},
			expectedOrderIdsInMemclob: map[clobtypes.OrderId]bool{
				constants.Order_Bob_Num0_Id11_Clob1_Sell5_Price15_GTB20.OrderId:      false,
				constants.Order_Alice_Num0_Id1_Clob1_Buy10_Price15_GTB20_IOC.OrderId: false,
			},
			expectedOrderFillAmounts: map[clobtypes.OrderId]uint64{
				constants.Order_Bob_Num0_Id11_Clob1_Sell5_Price15_GTB20.OrderId:      5000,
				constants.Order_Alice_Num0_Id1_Clob1_Buy10_Price15_GTB20_IOC.OrderId: 5000,
			},
		},
		"IOC fails CheckTx if previously filled": {
			blocks: []testmsgs.TestBlockWithMsgs{
				{
					Block: 2,
					Msgs: []testmsgs.TestSdkMsg{
						{
							Msg: clobtypes.NewMsgPlaceOrder(
								testapp.MustScaleOrder(
									constants.Order_Bob_Num0_Id11_Clob1_Buy5_Price40_GTB20,
									testapp.DefaultGenesis(),
								),
							),
							ExpectedIsOk: true,
						},
						{
							Msg: clobtypes.NewMsgPlaceOrder(
								testapp.MustScaleOrder(
									constants.Order_Alice_Num0_Id1_Clob1_Sell10_Price15_GTB20_IOC,
									testapp.DefaultGenesis(),
								),
							),
							ExpectedIsOk: true,
						},
					},
				},
				{
					Block: 3,
					Msgs: []testmsgs.TestSdkMsg{
						{
							Msg: clobtypes.NewMsgPlaceOrder(
								testapp.MustScaleOrder(
									constants.Order_Alice_Num0_Id1_Clob1_Sell10_Price15_GTB20_IOC,
									testapp.DefaultGenesis(),
								),
							),
							ExpectedIsOk:     false,
							ExpectedRespCode: clobtypes.ErrImmediateExecutionOrderAlreadyFilled.ABCICode(),
						},
					},
				},
			},
			expectedOrderIdsInMemclob: map[clobtypes.OrderId]bool{
				constants.Order_Bob_Num0_Id11_Clob1_Buy5_Price40_GTB20.OrderId:        false,
				constants.Order_Alice_Num0_Id1_Clob1_Sell10_Price15_GTB20_IOC.OrderId: false,
			},
			expectedOrderFillAmounts: map[clobtypes.OrderId]uint64{
				constants.Order_Bob_Num0_Id11_Clob1_Buy5_Price40_GTB20.OrderId:        5000,
				constants.Order_Alice_Num0_Id1_Clob1_Sell10_Price15_GTB20_IOC.OrderId: 5000,
			},
		},
		"Post-only buy does not cross and is placed on the book": {
			blocks: []testmsgs.TestBlockWithMsgs{
				{
					Block: 2,
					Msgs: []testmsgs.TestSdkMsg{
						{
							Msg: clobtypes.NewMsgPlaceOrder(
								testapp.MustScaleOrder(
									constants.Order_Bob_Num0_Id12_Clob0_Sell20_Price15_GTB20,
									testapp.DefaultGenesis(),
								),
							),
							ExpectedIsOk: true,
						},
						{
							Msg: clobtypes.NewMsgPlaceOrder(
								testapp.MustScaleOrder(
									constants.Order_Alice_Num0_Id1_Clob0_Buy15_Price10_GTB18_PO,
									testapp.DefaultGenesis(),
								),
							),
							ExpectedIsOk: true,
						},
					},
				},
			},
			expectedOrderIdsInMemclob: map[clobtypes.OrderId]bool{
				constants.Order_Bob_Num0_Id12_Clob0_Sell20_Price15_GTB20.OrderId:    true,
				constants.Order_Alice_Num0_Id1_Clob0_Buy15_Price10_GTB18_PO.OrderId: true,
			},
			expectedOrderFillAmounts: map[clobtypes.OrderId]uint64{
				constants.Order_Bob_Num0_Id12_Clob0_Sell20_Price15_GTB20.OrderId:    0,
				constants.Order_Alice_Num0_Id1_Clob0_Buy15_Price10_GTB18_PO.OrderId: 0,
			},
		},
		"Post-only sell does not cross and is placed on the book": {
			blocks: []testmsgs.TestBlockWithMsgs{
				{
					Block: 2,
					Msgs: []testmsgs.TestSdkMsg{
						{
							Msg: clobtypes.NewMsgPlaceOrder(
								testapp.MustScaleOrder(
									constants.Order_Bob_Num0_Id12_Clob0_Buy5_Price5_GTB20,
									testapp.DefaultGenesis(),
								),
							),
							ExpectedIsOk: true,
						},
						{
							Msg: clobtypes.NewMsgPlaceOrder(
								testapp.MustScaleOrder(
									constants.Order_Alice_Num0_Id1_Clob0_Sell15_Price10_GTB18_PO,
									testapp.DefaultGenesis(),
								),
							),
							ExpectedIsOk: true,
						},
					},
				},
			},
			expectedOrderIdsInMemclob: map[clobtypes.OrderId]bool{
				constants.Order_Bob_Num0_Id12_Clob0_Buy5_Price5_GTB20.OrderId:        true,
				constants.Order_Alice_Num0_Id1_Clob0_Sell15_Price10_GTB18_PO.OrderId: true,
			},
			expectedOrderFillAmounts: map[clobtypes.OrderId]uint64{
				constants.Order_Bob_Num0_Id12_Clob0_Buy5_Price5_GTB20.OrderId:        0,
				constants.Order_Alice_Num0_Id1_Clob0_Sell15_Price10_GTB18_PO.OrderId: 0,
			},
		},
		"Post-only buy crosses and is not placed on the book": {
			blocks: []testmsgs.TestBlockWithMsgs{
				{
					Block: 2,
					Msgs: []testmsgs.TestSdkMsg{
						{
							Msg: clobtypes.NewMsgPlaceOrder(
								testapp.MustScaleOrder(
									constants.Order_Bob_Num0_Id12_Clob0_Sell20_Price5_GTB20,
									testapp.DefaultGenesis(),
								),
							),
							ExpectedIsOk: true,
						},
						{
							Msg: clobtypes.NewMsgPlaceOrder(
								testapp.MustScaleOrder(
									constants.Order_Alice_Num0_Id1_Clob0_Buy15_Price10_GTB18_PO,
									testapp.DefaultGenesis(),
								),
							),
							ExpectedIsOk:     false,
							ExpectedRespCode: clobtypes.ErrPostOnlyWouldCrossMakerOrder.ABCICode(),
						},
					},
				},
			},
			expectedOrderIdsInMemclob: map[clobtypes.OrderId]bool{
				constants.Order_Bob_Num0_Id12_Clob0_Sell20_Price5_GTB20.OrderId:      true,
				constants.Order_Alice_Num0_Id1_Clob1_Sell5_Price15_GTB20_IOC.OrderId: false,
			},
			expectedOrderFillAmounts: map[clobtypes.OrderId]uint64{
				constants.Order_Bob_Num0_Id12_Clob0_Sell20_Price5_GTB20.OrderId:     0,
				constants.Order_Alice_Num0_Id1_Clob0_Buy15_Price10_GTB18_PO.OrderId: 0,
			},
		},
		"Post-only sell crosses and is not placed on the book": {
			blocks: []testmsgs.TestBlockWithMsgs{
				{
					Block: 2,
					Msgs: []testmsgs.TestSdkMsg{
						{
							Msg: clobtypes.NewMsgPlaceOrder(
								testapp.MustScaleOrder(
									constants.Order_Bob_Num0_Id12_Clob0_Buy5_Price40_GTB20,
									testapp.DefaultGenesis(),
								),
							),
							ExpectedIsOk: true,
						},
						{
							Msg: clobtypes.NewMsgPlaceOrder(
								testapp.MustScaleOrder(
									constants.Order_Alice_Num0_Id1_Clob0_Sell15_Price10_GTB18_PO,
									testapp.DefaultGenesis(),
								),
							),
							ExpectedIsOk:     false,
							ExpectedRespCode: clobtypes.ErrPostOnlyWouldCrossMakerOrder.ABCICode(),
						},
					},
				},
			},
			expectedOrderIdsInMemclob: map[clobtypes.OrderId]bool{
				constants.Order_Bob_Num0_Id12_Clob0_Buy5_Price40_GTB20.OrderId:       true,
				constants.Order_Alice_Num0_Id1_Clob0_Sell15_Price10_GTB18_PO.OrderId: false,
			},
			expectedOrderFillAmounts: map[clobtypes.OrderId]uint64{
				constants.Order_Bob_Num0_Id12_Clob0_Buy5_Price40_GTB20.OrderId:       0,
				constants.Order_Alice_Num0_Id1_Clob0_Sell15_Price10_GTB18_PO.OrderId: 0,
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).Build()
			ctx := tApp.InitChain()

			for _, block := range tc.blocks {
				for _, order := range block.Msgs {
					msgPlaceOrder, ok := order.Msg.(*clobtypes.MsgPlaceOrder)
					if !ok {
						t.Error("Expected MsgPlaceOrder")
					}
					for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(ctx, tApp.App, *msgPlaceOrder) {
						resp := tApp.CheckTx(checkTx)
						require.Equal(t, order.ExpectedIsOk, resp.IsOK(), "Response was not as expected: %+v", resp.Log)
						require.Equal(
							t,
							order.ExpectedRespCode,
							resp.Code,
							"Response code was not as expected",
						)
					}
				}
				ctx = tApp.AdvanceToBlock(block.Block, testapp.AdvanceToBlockOptions{})
			}

			for orderId, shouldHaveOrder := range tc.expectedOrderIdsInMemclob {
				_, exists := tApp.App.ClobKeeper.MemClob.GetOrder(orderId)
				require.Equal(t, shouldHaveOrder, exists)
			}

			for orderId, expectedFillAmount := range tc.expectedOrderFillAmounts {
				_, fillAmount, _ := tApp.App.ClobKeeper.GetOrderFillAmount(ctx, orderId)
				require.Equal(t, expectedFillAmount, fillAmount.ToUint64())
			}
		})
	}
}
