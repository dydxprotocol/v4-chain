package revshare_test

import (
	"math/big"
	"testing"

	"cosmossdk.io/math"
	"github.com/stretchr/testify/require"

	abcitypes "github.com/cometbft/cometbft/abci/types"

	cometbfttypes "github.com/cometbft/cometbft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	"github.com/dydxprotocol/v4-chain/protocol/indexer"
	indexerevents "github.com/dydxprotocol/v4-chain/protocol/indexer/events"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/msgsender"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	clobtestutils "github.com/dydxprotocol/v4-chain/protocol/testutil/clob"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	testtx "github.com/dydxprotocol/v4-chain/protocol/testutil/tx"
	testutil "github.com/dydxprotocol/v4-chain/protocol/testutil/util"
	"github.com/dydxprotocol/v4-chain/protocol/x/affiliates/types"
	assettypes "github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	revsharetypes "github.com/dydxprotocol/v4-chain/protocol/x/revshare/types"
	statstypes "github.com/dydxprotocol/v4-chain/protocol/x/stats/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

func TestPlaceOrderWithAffiliate(t *testing.T) {
	testCases := []struct {
		name                                  string
		initialUserStateStats                 *statstypes.UserStats
		initialUnconditionalRevShareConfig    *revsharetypes.UnconditionalRevShareConfig
		initialMarketMapperRevShareParams     *revsharetypes.MarketMapperRevenueShareParams
		expectedTakerFeeQuantums              int64
		expectedMakerFeeQuantums              int64
		expectedAffiliateRevShareQuantums     int64
		expectedUnconditionalRevShareQuantums int64
		expectedMarketMapperRevShareQuantums  int64
	}{
		{
			name:                                  "affiliate revshare",
			initialUserStateStats:                 nil,
			initialUnconditionalRevShareConfig:    nil,
			initialMarketMapperRevShareParams:     nil,
			expectedTakerFeeQuantums:              2000,
			expectedMakerFeeQuantums:              550,
			expectedAffiliateRevShareQuantums:     300,
			expectedUnconditionalRevShareQuantums: 0,
			expectedMarketMapperRevShareQuantums:  0,
		},
		{
			name: "Affiliate over limit",
			initialUserStateStats: &statstypes.UserStats{
				TakerNotional: uint64(35_000_000_000_000),
				MakerNotional: uint64(35_000_000_000_000),
			},
			initialUnconditionalRevShareConfig:    nil,
			initialMarketMapperRevShareParams:     nil,
			expectedTakerFeeQuantums:              1750,
			expectedMakerFeeQuantums:              550,
			expectedAffiliateRevShareQuantums:     0,
			expectedUnconditionalRevShareQuantums: 0,
			expectedMarketMapperRevShareQuantums:  0,
		},
		{
			name:                  "affiliate revshare + unconditional revshare",
			initialUserStateStats: nil,
			initialUnconditionalRevShareConfig: &revsharetypes.UnconditionalRevShareConfig{
				Configs: []revsharetypes.UnconditionalRevShareConfig_RecipientConfig{
					{
						Address:  constants.Carl_Num0.Owner,
						SharePpm: 100_000,
					},
				},
			},
			initialMarketMapperRevShareParams:     nil,
			expectedTakerFeeQuantums:              2000,
			expectedMakerFeeQuantums:              550,
			expectedAffiliateRevShareQuantums:     300,
			expectedUnconditionalRevShareQuantums: 145,
			expectedMarketMapperRevShareQuantums:  0,
		},
		{
			name:                               "affiliate + market mapper revshare",
			initialUserStateStats:              nil,
			initialUnconditionalRevShareConfig: nil,
			initialMarketMapperRevShareParams: &revsharetypes.MarketMapperRevenueShareParams{
				Address:         constants.DaveAccAddress.String(),
				RevenueSharePpm: 100_000, // 10%
				ValidDays:       1,
			},
			expectedTakerFeeQuantums:              2000,
			expectedMakerFeeQuantums:              550,
			expectedAffiliateRevShareQuantums:     300,
			expectedUnconditionalRevShareQuantums: 0,
			expectedMarketMapperRevShareQuantums:  145,
		},
		{
			name:                  "affiliate + market mapper revshare + unconditional revshare",
			initialUserStateStats: nil,
			initialUnconditionalRevShareConfig: &revsharetypes.UnconditionalRevShareConfig{
				Configs: []revsharetypes.UnconditionalRevShareConfig_RecipientConfig{
					{
						Address:  constants.Carl_Num0.Owner,
						SharePpm: 100_000, // 10%
					},
				},
			},
			initialMarketMapperRevShareParams: &revsharetypes.MarketMapperRevenueShareParams{
				Address:         constants.DaveAccAddress.String(),
				RevenueSharePpm: 100_000, // 10%
				ValidDays:       1,
			},
			expectedTakerFeeQuantums:              2000,
			expectedMakerFeeQuantums:              550,
			expectedAffiliateRevShareQuantums:     300,
			expectedUnconditionalRevShareQuantums: 145,
			expectedMarketMapperRevShareQuantums:  145,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			tApp, ctx, msgSender := setupTest(t, tc.initialUserStateStats,
				tc.initialUnconditionalRevShareConfig, tc.initialMarketMapperRevShareParams)

			bankKeeper := tApp.App.BankKeeper
			carlBankBalance := bankKeeper.GetBalance(ctx, constants.CarlAccAddress, assettypes.AssetUsdc.GetDenom())
			aliceBankBalance := bankKeeper.GetBalance(ctx, constants.AliceAccAddress, assettypes.AssetUsdc.GetDenom())
			daveBankBalance := bankKeeper.GetBalance(ctx, constants.DaveAccAddress, assettypes.AssetUsdc.GetDenom())
			// Setup orders
			Clob_0, aliceOrder, bobOrder, aliceCheckTx, bobCheckTx, orders := setupOrders(ctx, tApp)

			// Get expected onchain messages
			expectedOnchainMessagesInNextBlock := getExpectedOnchainMessagesInNextBlock(
				ctx,
				Clob_0,
				aliceOrder,
				bobOrder,
				aliceCheckTx.Tx,
				bobCheckTx.Tx,
				big.NewInt(100000000005000000-tc.expectedTakerFeeQuantums),
				big.NewInt(99999999995000000+tc.expectedMakerFeeQuantums),
				tc.expectedTakerFeeQuantums,
				tc.expectedMakerFeeQuantums,
				tc.expectedAffiliateRevShareQuantums,
				4,
			)

			// Place orders
			for _, order := range orders {
				for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(ctx, tApp.App, order) {
					resp := tApp.CheckTx(checkTx)
					require.Conditionf(t, resp.IsOK, "Expected CheckTx to succeed. Response: %+v", resp)
				}
			}
			// Clear msgSender and advance to next block
			msgSender.Clear()
			ctx = tApp.AdvanceToBlock(4, testapp.AdvanceToBlockOptions{})

			// Check onchain messages
			require.ElementsMatch(t, expectedOnchainMessagesInNextBlock, msgSender.GetOnchainMessages())

			// Check bank balances
			carlBankBalanceAfter := bankKeeper.GetBalance(ctx,
				constants.CarlAccAddress,
				assettypes.AssetUsdc.GetDenom())
			aliceBankBalanceAfter := bankKeeper.GetBalance(ctx,
				constants.AliceAccAddress,
				assettypes.AssetUsdc.GetDenom())
			daveBankBalanceAfter := bankKeeper.GetBalance(ctx,
				constants.DaveAccAddress,
				assettypes.AssetUsdc.GetDenom())
			require.Equal(t, carlBankBalance.Add(sdk.NewCoin(assettypes.AssetUsdc.GetDenom(),
				math.NewInt(tc.expectedUnconditionalRevShareQuantums))), carlBankBalanceAfter)
			require.Equal(t, aliceBankBalance.Add(sdk.NewCoin(assettypes.AssetUsdc.GetDenom(),
				math.NewInt(tc.expectedAffiliateRevShareQuantums))), aliceBankBalanceAfter)
			require.Equal(t, daveBankBalance.Add(sdk.NewCoin(assettypes.AssetUsdc.GetDenom(),
				math.NewInt(tc.expectedMarketMapperRevShareQuantums))), daveBankBalanceAfter)
		})
	}
}

func setupTest(t *testing.T, initialUserStateStats *statstypes.UserStats,
	initialUnconditionalRevShareConfig *revsharetypes.UnconditionalRevShareConfig,
	initialMarketMapperRevShareParams *revsharetypes.MarketMapperRevenueShareParams) (*testapp.TestApp,
	sdk.Context, *msgsender.IndexerMessageSenderInMemoryCollector) {
	msgSender := msgsender.NewIndexerMessageSenderInMemoryCollector()
	appOpts := map[string]interface{}{
		indexer.MsgSenderInstanceForTest: msgSender,
	}

	tApp := testapp.NewTestAppBuilder(t).WithAppOptions(appOpts).WithGenesisDocFn(
		func() (genesis cometbfttypes.GenesisDoc) {
			genesis = testapp.DefaultGenesis()
			testapp.UpdateGenesisDocWithAppStateForModule(
				&genesis,
				func(genesisState *types.GenesisState) {
					genesisState.AffiliateTiers = types.DefaultAffiliateTiers
				},
			)
			if initialUserStateStats != nil {
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *statstypes.GenesisState) {
						genesisState.AddressToUserStats = []*statstypes.AddressToUserStats{
							{
								Address:   constants.Bob_Num0.Owner,
								UserStats: initialUserStateStats,
							},
						}
					},
				)
			}
			if initialUnconditionalRevShareConfig != nil || initialMarketMapperRevShareParams != nil {
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *revsharetypes.GenesisState) {
						if initialUnconditionalRevShareConfig != nil {
							genesisState.UnconditionalRevShareConfig = *initialUnconditionalRevShareConfig
						}
						if initialMarketMapperRevShareParams != nil {
							genesisState.Params = *initialMarketMapperRevShareParams
						}
					},
				)
			}
			return genesis
		}).Build()
	ctx := tApp.InitChain()

	return tApp, ctx, msgSender
}

func setupOrders(ctx sdk.Context, tApp *testapp.TestApp) (clobtypes.ClobPair,
	clobtypes.MsgPlaceOrder,
	clobtypes.MsgPlaceOrder,
	abcitypes.RequestCheckTx,
	abcitypes.RequestCheckTx,
	[]clobtypes.MsgPlaceOrder) {
	Clob_0 := testapp.MustGetClobPairsFromGenesis(testapp.DefaultGenesis())[0]
	msgRegisterAffiliate := types.MsgRegisterAffiliate{
		Referee:   constants.Bob_Num0.Owner,
		Affiliate: constants.Alice_Num0.Owner,
	}

	checkTxMsgRegisterAffiliate := testapp.MustMakeCheckTx(
		ctx,
		tApp.App,
		testapp.MustMakeCheckTxOptions{
			AccAddressForSigning: constants.Bob_Num0.Owner,
			Gas:                  constants.TestGasLimit,
			FeeAmt:               constants.TestFeeCoins_5Cents,
		},
		&msgRegisterAffiliate,
	)
	tApp.CheckTx(checkTxMsgRegisterAffiliate)

	ctx = tApp.AdvanceToBlock(3, testapp.AdvanceToBlockOptions{})

	aliceOrder := *clobtypes.NewMsgPlaceOrder(testapp.MustScaleOrder(
		clobtypes.Order{
			OrderId:      clobtypes.OrderId{SubaccountId: constants.Alice_Num0, ClientId: 0, ClobPairId: 0},
			Side:         clobtypes.Order_SIDE_BUY,
			Quantums:     500_000_000,
			Subticks:     10,
			GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
		},
		testapp.DefaultGenesis(),
	))

	bobOrder := *clobtypes.NewMsgPlaceOrder(testapp.MustScaleOrder(
		clobtypes.Order{
			OrderId:      clobtypes.OrderId{SubaccountId: constants.Bob_Num0, ClientId: 0, ClobPairId: 0},
			Side:         clobtypes.Order_SIDE_SELL,
			Quantums:     500_000_000,
			Subticks:     10,
			GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
		},
		testapp.DefaultGenesis(),
	))

	aliceCheckTx := testapp.MustMakeCheckTx(
		ctx,
		tApp.App,
		testapp.MustMakeCheckTxOptions{
			AccAddressForSigning: constants.Alice_Num0.Owner,
		},
		&aliceOrder,
	)

	bobCheckTx := testapp.MustMakeCheckTx(
		ctx,
		tApp.App,
		testapp.MustMakeCheckTxOptions{
			AccAddressForSigning: constants.Bob_Num0.Owner,
		},
		&bobOrder,
	)

	orders := []clobtypes.MsgPlaceOrder{
		aliceOrder,
		bobOrder,
	}
	return Clob_0, aliceOrder, bobOrder, aliceCheckTx, bobCheckTx, orders
}

func getExpectedOnchainMessagesInNextBlock(ctx sdk.Context, Clob_0 clobtypes.ClobPair,
	aliceOrder clobtypes.MsgPlaceOrder, bobOrder clobtypes.MsgPlaceOrder,
	aliceCheckTxHash []byte, bobCheckTxHash []byte, expectedAliceAssetQuantums *big.Int,
	expectedBobAssetQuantums *big.Int,
	expectedTakerFeeQuantums int64, expectedMakerFeeQuantums int64,
	expectedAffiliateRevShareQuantums int64, blockHeight uint32) []msgsender.Message {
	return []msgsender.Message{indexer_manager.CreateIndexerBlockEventMessage(
		&indexer_manager.IndexerTendermintBlock{
			Height: blockHeight,
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
										aliceOrder.Order.GetQuantums())),
									big.NewInt(0),
									big.NewInt(0),
								),
							},
							[]*satypes.AssetPosition{
								testutil.CreateSingleAssetPosition(
									assettypes.AssetUsdc.Id,
									expectedAliceAssetQuantums,
								),
							},
							nil,
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
										aliceOrder.Order.GetQuantums())),
									big.NewInt(0),
									big.NewInt(0),
								),
							},
							[]*satypes.AssetPosition{
								testutil.CreateSingleAssetPosition(
									assettypes.AssetUsdc.Id,
									expectedBobAssetQuantums,
								),
							},
							nil,
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
							aliceOrder.Order,
							bobOrder.Order,
							aliceOrder.Order.GetBaseQuantums(),
							-expectedMakerFeeQuantums,
							expectedTakerFeeQuantums,
							aliceOrder.Order.GetBaseQuantums(),
							aliceOrder.Order.GetBaseQuantums(),
							big.NewInt(expectedAffiliateRevShareQuantums),
						),
					),
				},
				{
					Subtype: indexerevents.SubtypeOpenInterestUpdate,
					OrderingWithinBlock: &indexer_manager.IndexerTendermintEvent_BlockEvent_{
						BlockEvent: indexer_manager.IndexerTendermintEvent_BLOCK_EVENT_END_BLOCK,
					},
					Version: indexerevents.OpenInterestUpdateVersion,
					DataBytes: indexer_manager.GetBytes(
						&indexerevents.OpenInterestUpdateEventV1{
							OpenInterestUpdates: []*indexerevents.OpenInterestUpdate{
								{
									PerpetualId: Clob_0.MustGetPerpetualId(),
									OpenInterest: dtypes.NewIntFromUint64(
										aliceOrder.Order.GetBigQuantums().Uint64(),
									),
								},
							},
						}),
				},
			},
			TxHashes: []string{string(lib.GetTxHash(testtx.MustGetTxBytes(&clobtypes.MsgProposedOperations{
				OperationsQueue: []clobtypes.OperationRaw{
					{
						Operation: &clobtypes.OperationRaw_ShortTermOrderPlacement{
							ShortTermOrderPlacement: aliceCheckTxHash,
						},
					},
					{
						Operation: &clobtypes.OperationRaw_ShortTermOrderPlacement{
							ShortTermOrderPlacement: bobCheckTxHash,
						},
					},
					clobtestutils.NewMatchOperationRaw(
						&bobOrder.Order,
						[]clobtypes.MakerFill{
							{
								FillAmount: aliceOrder.
									Order.GetBaseQuantums().ToUint64(),
								MakerOrderId: aliceOrder.Order.OrderId,
							},
						},
					),
				},
			},
			)))},
		})}
}
