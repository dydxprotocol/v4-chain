package keeper_test

/*

func TestProcessProposerMatches_LongTerm_Success(t *testing.T) {
	tests := map[string]processProposerMatchesTestCase{
		"Succeeds with new maker Long-Term order": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short,
				constants.Dave_Num0_1BTC_Long,
			},
			setupMockBankKeeper: func(bk *mocks.BankKeeper) {
				bk.On(
					"SendCoinsFromModuleToModule",
					mock.Anything,
					satypes.ModuleName,
					authtypes.FeeCollectorName,
					mock.MatchedBy(
						testutil_bank.MatchUsdcOfAmount(
							25_000_000+10_000_000,
						),
					),
				).Return(nil).Once()
			},
			setupState: func(t *testing.T, ctx sdk.Context, k *keeper.Keeper) {
				k.SetBlockTimeForLastCommittedBlock(ctx.WithBlockTime(time.Unix(5, 0)))
			},
			placeOrders: []*types.MsgPlaceOrder{
				{
					Order: constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10,
				},
				{
					Order: constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10,
				},
			},
			clobMatches: []*types.ClobMatch{
				types.NewClobMatchFromMatchOrders(
					&types.MatchOrders{
						MakerOneof: &types.MatchOrders_MakerOrderId{
							MakerOrderId: &constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10.OrderId,
						},
						TakerOneof: &types.MatchOrders_TakerOrderIndex{TakerOrderIndex: 0},
						FillAmount: 100_000_000, // 1 BTC
					},
				),
			},
			expectedFillAmounts: map[types.OrderId]satypes.BaseQuantums{
				constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10.OrderId:         100_000_000,
				constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10.OrderId: 100_000_000,
			},
			expectedPruneableBlockHeights: map[uint32][]types.OrderId{
				10 + types.ShortBlockWindow: {
					constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10.OrderId,
				},
			},
			expectedQuoteBalances: map[satypes.SubaccountId]*big.Int{
				// $99,975
				constants.Dave_Num0: big.NewInt(100_000_000_000 - 25_000_000),
				// $49,990
				constants.Carl_Num0: big.NewInt(50_000_000_000 - 10_000_000),
			},
			expectedPerpetualPositions: map[satypes.SubaccountId][]*satypes.PerpetualPosition{
				constants.Carl_Num0: {},
				constants.Dave_Num0: {},
			},
		},
		"Succeeds with new taker Long-Term order": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short,
				constants.Dave_Num0_1BTC_Long,
			},
			setupMockBankKeeper: func(bk *mocks.BankKeeper) {
				bk.On(
					"SendCoinsFromModuleToModule",
					mock.Anything,
					satypes.ModuleName,
					authtypes.FeeCollectorName,
					mock.MatchedBy(
						testutil_bank.MatchUsdcOfAmount(
							25_000_000+10_000_000,
						),
					),
				).Return(nil).Once()
			},
			setupState: func(t *testing.T, ctx sdk.Context, k *keeper.Keeper) {
				k.SetBlockTimeForLastCommittedBlock(ctx.WithBlockTime(time.Unix(5, 0)))
			},
			placeOrders: []*types.MsgPlaceOrder{
				{
					Order: constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10,
				},
				{
					Order: constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10,
				},
			},
			clobMatches: []*types.ClobMatch{
				types.NewClobMatchFromMatchOrders(
					&types.MatchOrders{
						MakerOneof: &types.MatchOrders_MakerOrderIndex{MakerOrderIndex: 0},
						TakerOneof: &types.MatchOrders_TakerOrderId{
							TakerOrderId: &constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10.OrderId,
						},
						FillAmount: 100_000_000, // 1 BTC
					},
				),
			},
			expectedFillAmounts: map[types.OrderId]satypes.BaseQuantums{
				constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10.OrderId:         100_000_000,
				constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10.OrderId: 100_000_000,
			},
			expectedPruneableBlockHeights: map[uint32][]types.OrderId{
				10 + types.ShortBlockWindow: {
					constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10.OrderId,
				},
			},
			expectedQuoteBalances: map[satypes.SubaccountId]*big.Int{
				// $99,990
				constants.Dave_Num0: big.NewInt(100_000_000_000 - 10_000_000),
				// $49,975
				constants.Carl_Num0: big.NewInt(50_000_000_000 - 25_000_000),
			},
			expectedPerpetualPositions: map[satypes.SubaccountId][]*satypes.PerpetualPosition{
				constants.Carl_Num0: {},
				constants.Dave_Num0: {},
			},
		},
		"Succeeds with existing maker Long-Term order": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short,
				constants.Dave_Num0_1BTC_Long,
			},
			setupMockBankKeeper: func(bk *mocks.BankKeeper) {
				bk.On(
					"SendCoinsFromModuleToModule",
					mock.Anything,
					satypes.ModuleName,
					authtypes.FeeCollectorName,
					mock.MatchedBy(
						testutil_bank.MatchUsdcOfAmount(
							25_000_000+10_000_000,
						),
					),
				).Return(nil).Once()
			},
			setupState: func(t *testing.T, ctx sdk.Context, k *keeper.Keeper) {
				k.SetBlockTimeForLastCommittedBlock(ctx.WithBlockTime(time.Unix(5, 0)))
				k.SetLongTermOrderPlacement(
					ctx,
					constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10,
					4, // Placed in previous block.
				)
			},
			placeOrders: []*types.MsgPlaceOrder{
				{
					Order: constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10,
				},
			},
			clobMatches: []*types.ClobMatch{
				types.NewClobMatchFromMatchOrders(
					&types.MatchOrders{
						MakerOneof: &types.MatchOrders_MakerOrderId{
							MakerOrderId: &constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10.OrderId,
						},
						TakerOneof: &types.MatchOrders_TakerOrderIndex{TakerOrderIndex: 0},
						FillAmount: 100_000_000, // 1 BTC
					},
				),
			},
			expectedFillAmounts: map[types.OrderId]satypes.BaseQuantums{
				constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10.OrderId:         100_000_000,
				constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10.OrderId: 100_000_000,
			},
			expectedPruneableBlockHeights: map[uint32][]types.OrderId{
				10 + types.ShortBlockWindow: {
					constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10.OrderId,
				},
			},
			expectedQuoteBalances: map[satypes.SubaccountId]*big.Int{
				// $99,975
				constants.Dave_Num0: big.NewInt(100_000_000_000 - 25_000_000),
				// $49,990
				constants.Carl_Num0: big.NewInt(50_000_000_000 - 10_000_000),
			},
			expectedPerpetualPositions: map[satypes.SubaccountId][]*satypes.PerpetualPosition{
				constants.Carl_Num0: {},
				constants.Dave_Num0: {},
			},
		},
		"Succeeds with existing taker Long-Term order": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short,
				constants.Dave_Num0_1BTC_Long,
			},
			setupMockBankKeeper: func(bk *mocks.BankKeeper) {
				bk.On(
					"SendCoinsFromModuleToModule",
					mock.Anything,
					satypes.ModuleName,
					authtypes.FeeCollectorName,
					mock.MatchedBy(
						testutil_bank.MatchUsdcOfAmount(
							25_000_000+10_000_000,
						),
					),
				).Return(nil).Once()
			},
			setupState: func(t *testing.T, ctx sdk.Context, k *keeper.Keeper) {
				k.SetBlockTimeForLastCommittedBlock(ctx.WithBlockTime(time.Unix(5, 0)))
				k.SetLongTermOrderPlacement(
					ctx,
					constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10,
					4, // Placed in previous block.
				)
			},
			placeOrders: []*types.MsgPlaceOrder{
				{
					Order: constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10,
				},
			},
			clobMatches: []*types.ClobMatch{
				types.NewClobMatchFromMatchOrders(
					&types.MatchOrders{
						MakerOneof: &types.MatchOrders_MakerOrderIndex{MakerOrderIndex: 0},
						TakerOneof: &types.MatchOrders_TakerOrderId{
							TakerOrderId: &constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10.OrderId,
						},
						FillAmount: 100_000_000, // 1 BTC
					},
				),
			},
			expectedFillAmounts: map[types.OrderId]satypes.BaseQuantums{
				constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10.OrderId:         100_000_000,
				constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10.OrderId: 100_000_000,
			},
			expectedPruneableBlockHeights: map[uint32][]types.OrderId{
				10 + types.ShortBlockWindow: {
					constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10.OrderId,
				},
			},
			expectedQuoteBalances: map[satypes.SubaccountId]*big.Int{
				// $99,990
				constants.Dave_Num0: big.NewInt(100_000_000_000 - 10_000_000),
				// $49,975
				constants.Carl_Num0: big.NewInt(50_000_000_000 - 25_000_000),
			},
			expectedPerpetualPositions: map[satypes.SubaccountId][]*satypes.PerpetualPosition{
				constants.Carl_Num0: {},
				constants.Dave_Num0: {},
			},
		},
		"Succeeds with new maker and taker Long-Term orders completely filled": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short,
				constants.Dave_Num0_1BTC_Long,
			},
			setupMockBankKeeper: func(bk *mocks.BankKeeper) {
				bk.On(
					"SendCoinsFromModuleToModule",
					mock.Anything,
					satypes.ModuleName,
					authtypes.FeeCollectorName,
					mock.MatchedBy(
						testutil_bank.MatchUsdcOfAmount(
							25_000_000+10_000_000,
						),
					),
				).Return(nil).Once()
			},
			setupState: func(t *testing.T, ctx sdk.Context, k *keeper.Keeper) {
				k.SetBlockTimeForLastCommittedBlock(ctx.WithBlockTime(time.Unix(5, 0)))
			},
			placeOrders: []*types.MsgPlaceOrder{
				{
					Order: constants.LongTermOrder_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10,
				},
				{
					Order: constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10,
				},
			},
			clobMatches: []*types.ClobMatch{
				types.NewClobMatchFromMatchOrders(
					&types.MatchOrders{
						MakerOneof: &types.MatchOrders_MakerOrderId{
							MakerOrderId: &constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10.OrderId,
						},
						TakerOneof: &types.MatchOrders_TakerOrderId{
							TakerOrderId: &constants.LongTermOrder_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10.OrderId,
						},
						FillAmount: 100_000_000, // 1 BTC
					},
				),
			},
			expectedFillAmounts: map[types.OrderId]satypes.BaseQuantums{
				constants.LongTermOrder_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10.OrderId: 100_000_000,
				constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10.OrderId:  100_000_000,
			},
			expectedPruneableBlockHeights: map[uint32][]types.OrderId{},
			expectedQuoteBalances: map[satypes.SubaccountId]*big.Int{
				// $99,975
				constants.Dave_Num0: big.NewInt(100_000_000_000 - 25_000_000),
				// $49,990
				constants.Carl_Num0: big.NewInt(50_000_000_000 - 10_000_000),
			},
			expectedPerpetualPositions: map[satypes.SubaccountId][]*satypes.PerpetualPosition{
				constants.Carl_Num0: {},
				constants.Dave_Num0: {},
			},
		},
		"Succeeds with new maker and taker Long-Term orders partially filled": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short,
				constants.Dave_Num0_1BTC_Long,
			},
			setupMockBankKeeper: func(bk *mocks.BankKeeper) {
				bk.On(
					"SendCoinsFromModuleToModule",
					mock.Anything,
					satypes.ModuleName,
					authtypes.FeeCollectorName,
					mock.MatchedBy(
						testutil_bank.MatchUsdcOfAmount(
							12_500_000+5_000_000,
						),
					),
				).Return(nil).Once()
			},
			setupState: func(t *testing.T, ctx sdk.Context, k *keeper.Keeper) {
				k.SetBlockTimeForLastCommittedBlock(ctx.WithBlockTime(time.Unix(5, 0)))
			},
			placeOrders: []*types.MsgPlaceOrder{
				{
					Order: constants.LongTermOrder_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10,
				},
				{
					Order: constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10,
				},
			},
			clobMatches: []*types.ClobMatch{
				types.NewClobMatchFromMatchOrders(
					&types.MatchOrders{
						MakerOneof: &types.MatchOrders_MakerOrderId{
							MakerOrderId: &constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10.OrderId,
						},
						TakerOneof: &types.MatchOrders_TakerOrderId{
							TakerOrderId: &constants.LongTermOrder_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10.OrderId,
						},
						FillAmount: 50_000_000, // 0.5 BTC
					},
				),
			},
			expectedFillAmounts: map[types.OrderId]satypes.BaseQuantums{
				constants.LongTermOrder_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10.OrderId: 50_000_000,
				constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10.OrderId:  50_000_000,
			},
			expectedPruneableBlockHeights: map[uint32][]types.OrderId{},
			expectedQuoteBalances: map[satypes.SubaccountId]*big.Int{
				// $74,987.5
				constants.Dave_Num0: big.NewInt(75_000_000_000 - 12_500_000),
				// $74,995
				constants.Carl_Num0: big.NewInt(75_000_000_000 - 5_000_000),
			},
			expectedPerpetualPositions: map[satypes.SubaccountId][]*satypes.PerpetualPosition{
				constants.Carl_Num0: {
					{
						PerpetualId:  0,
						Quantums:     dtypes.NewInt(-50_000_000), // .5 BTC
						FundingIndex: dtypes.ZeroInt(),
					},
				},
				constants.Dave_Num0: {
					{
						PerpetualId:  0,
						Quantums:     dtypes.NewInt(50_000_000), // .5 BTC
						FundingIndex: dtypes.ZeroInt(),
					},
				},
			},
		},
		"Succeeds with Long-Term order and multiple fills": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short,
				constants.Dave_Num0_1BTC_Long,
			},
			setupMockBankKeeper: func(bk *mocks.BankKeeper) {
				bk.On(
					"SendCoinsFromModuleToModule",
					mock.Anything,
					satypes.ModuleName,
					authtypes.FeeCollectorName,
					mock.MatchedBy(
						testutil_bank.MatchUsdcOfAmount(
							12_500_000+5_000_000,
						),
					),
				).Return(nil).Twice()
			},
			setupState: func(t *testing.T, ctx sdk.Context, k *keeper.Keeper) {
				k.SetBlockTimeForLastCommittedBlock(ctx.WithBlockTime(time.Unix(5, 0)))
			},
			placeOrders: []*types.MsgPlaceOrder{
				{
					Order: constants.Order_Carl_Num0_Id0_Clob0_Buy05BTC_Price50000_GTB10_FOK,
				},
				{
					Order: constants.Order_Carl_Num0_Id2_Clob0_Buy05BTC_Price50000,
				},
				{
					Order: constants.LongTermOrder_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10,
				},
			},
			clobMatches: []*types.ClobMatch{
				types.NewClobMatchFromMatchOrders(
					&types.MatchOrders{
						MakerOneof: &types.MatchOrders_MakerOrderId{
							MakerOrderId: &constants.LongTermOrder_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10.OrderId,
						},
						TakerOneof: &types.MatchOrders_TakerOrderIndex{TakerOrderIndex: 0},
						FillAmount: 50_000_000, // 0.5 BTC
					},
				),
				types.NewClobMatchFromMatchOrders(
					&types.MatchOrders{
						MakerOneof: &types.MatchOrders_MakerOrderId{
							MakerOrderId: &constants.LongTermOrder_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10.OrderId,
						},
						TakerOneof: &types.MatchOrders_TakerOrderIndex{TakerOrderIndex: 1},
						FillAmount: 50_000_000, // 0.5 BTC
					},
				),
			},
			expectedFillAmounts: map[types.OrderId]satypes.BaseQuantums{
				constants.Order_Carl_Num0_Id0_Clob0_Buy05BTC_Price50000_GTB10_FOK.OrderId:      50_000_000,
				constants.Order_Carl_Num0_Id2_Clob0_Buy05BTC_Price50000.OrderId:                50_000_000,
				constants.LongTermOrder_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10.OrderId: 100_000_000,
			},
			expectedPruneableBlockHeights: map[uint32][]types.OrderId{
				10 + types.ShortBlockWindow: {
					constants.Order_Carl_Num0_Id0_Clob0_Buy05BTC_Price50000_GTB10_FOK.OrderId,
					constants.Order_Carl_Num0_Id2_Clob0_Buy05BTC_Price50000.OrderId,
				},
			},
			expectedQuoteBalances: map[satypes.SubaccountId]*big.Int{
				// $99,990
				constants.Dave_Num0: big.NewInt(100_000_000_000 - 10_000_000),
				// $49,975
				constants.Carl_Num0: big.NewInt(50_000_000_000 - 25_000_000),
			},
			expectedPerpetualPositions: map[satypes.SubaccountId][]*satypes.PerpetualPosition{
				constants.Carl_Num0: {},
				constants.Dave_Num0: {},
			},
		},
		"Succeeds with new maker Long-Term order in liquidation match": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_54999USD,
				constants.Dave_Num0_1BTC_Long,
			},
			setupMockBankKeeper: func(bk *mocks.BankKeeper) {
				bk.On(
					"SendCoinsFromModuleToModule",
					mock.Anything,
					satypes.ModuleName,
					authtypes.FeeCollectorName,
					mock.MatchedBy(testutil_bank.MatchUsdcOfAmount(10_000_000)),
				).Return(nil)
				bk.On(
					"SendCoinsFromModuleToModule",
					mock.Anything,
					satypes.ModuleName,
					types.InsuranceFundName,
					// Subaccount pays $250 to insurance fund for liquidating 1 BTC.
					mock.MatchedBy(testutil_bank.MatchUsdcOfAmount(250_000_000)),
				).Return(nil).Once()
			},
			setupState: func(t *testing.T, ctx sdk.Context, k *keeper.Keeper) {
				k.SetBlockTimeForLastCommittedBlock(ctx.WithBlockTime(time.Unix(5, 0)))
			},
			placeOrders: []*types.MsgPlaceOrder{
				{
					Order: constants.LongTermOrder_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10,
				},
			},
			clobMatches: []*types.ClobMatch{
				types.NewClobMatchFromMatchPerpetualLiquidation(
					&types.MatchPerpetualLiquidation{
						Liquidated:  constants.Carl_Num0,
						ClobPairId:  0,
						PerpetualId: 0,
						TotalSize:   100_000_000, // 1 BTC
						IsBuy:       true,
						Fills: []types.MatchPerpetualLiquidation_Fill{
							{
								MakerOneof: &types.MatchPerpetualLiquidation_Fill_MakerOrderId{
									MakerOrderId: &constants.LongTermOrder_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10.OrderId,
								},
								FillAmount: 100_000_000, // 1 BTC
							},
						},
					},
				),
			},
			expectedFillAmounts: map[types.OrderId]satypes.BaseQuantums{
				constants.LongTermOrder_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10.OrderId: 100_000_000,
			},
			expectedPruneableBlockHeights: map[uint32][]types.OrderId{},
			expectedQuoteBalances: map[satypes.SubaccountId]*big.Int{
				// $4,749, no taker fees, pays $250 insurance fee
				constants.Carl_Num0: big.NewInt(4_999_000_000 - 250_000_000),
				// $99,990
				constants.Dave_Num0: big.NewInt(100_000_000_000 - 10_000_000),
			},
			expectedPerpetualPositions: map[satypes.SubaccountId][]*satypes.PerpetualPosition{
				constants.Dave_Num0: {},
				constants.Carl_Num0: {},
			},
			expectedSubaccountLiquidationInfo: map[satypes.SubaccountId]types.SubaccountLiquidationInfo{
				constants.Carl_Num0: {
					PerpetualsLiquidated:  []uint32{0},
					NotionalLiquidated:    50_000_000_000, // Liquidated 1BTC at $50,000.
					QuantumsInsuranceLost: 0,
				},
				constants.Dave_Num0: {},
			},
		},
		"Succeeds with existing maker Long-Term order in liquidation match": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_54999USD,
				constants.Dave_Num0_1BTC_Long,
			},
			setupMockBankKeeper: func(bk *mocks.BankKeeper) {
				bk.On(
					"SendCoinsFromModuleToModule",
					mock.Anything,
					satypes.ModuleName,
					authtypes.FeeCollectorName,
					mock.MatchedBy(testutil_bank.MatchUsdcOfAmount(10_000_000)),
				).Return(nil)
				bk.On(
					"SendCoinsFromModuleToModule",
					mock.Anything,
					satypes.ModuleName,
					types.InsuranceFundName,
					// Subaccount pays $250 to insurance fund for liquidating 1 BTC.
					mock.MatchedBy(testutil_bank.MatchUsdcOfAmount(250_000_000)),
				).Return(nil).Once()
			},
			setupState: func(t *testing.T, ctx sdk.Context, k *keeper.Keeper) {
				k.SetBlockTimeForLastCommittedBlock(ctx.WithBlockTime(time.Unix(5, 0)))
				k.SetLongTermOrderPlacement(
					ctx,
					constants.LongTermOrder_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10,
					4, // Placed in previous block.
				)
			},
			placeOrders: []*types.MsgPlaceOrder{},
			clobMatches: []*types.ClobMatch{
				types.NewClobMatchFromMatchPerpetualLiquidation(
					&types.MatchPerpetualLiquidation{
						Liquidated:  constants.Carl_Num0,
						ClobPairId:  0,
						PerpetualId: 0,
						TotalSize:   100_000_000, // 1 BTC
						IsBuy:       true,
						Fills: []types.MatchPerpetualLiquidation_Fill{
							{
								MakerOneof: &types.MatchPerpetualLiquidation_Fill_MakerOrderId{
									MakerOrderId: &constants.LongTermOrder_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10.OrderId,
								},
								FillAmount: 100_000_000, // 1 BTC
							},
						},
					},
				),
			},
			expectedFillAmounts: map[types.OrderId]satypes.BaseQuantums{
				constants.LongTermOrder_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10.OrderId: 100_000_000,
			},
			expectedPruneableBlockHeights: map[uint32][]types.OrderId{},
			expectedQuoteBalances: map[satypes.SubaccountId]*big.Int{
				// $4,749, no taker fees, pays $250 insurance fee
				constants.Carl_Num0: big.NewInt(4_999_000_000 - 250_000_000),
				// $99,990
				constants.Dave_Num0: big.NewInt(100_000_000_000 - 10_000_000),
			},
			expectedPerpetualPositions: map[satypes.SubaccountId][]*satypes.PerpetualPosition{
				constants.Dave_Num0: {},
				constants.Carl_Num0: {},
			},
			expectedSubaccountLiquidationInfo: map[satypes.SubaccountId]types.SubaccountLiquidationInfo{
				constants.Carl_Num0: {
					PerpetualsLiquidated:  []uint32{0},
					NotionalLiquidated:    50_000_000_000, // Liquidated 1BTC at $50,000.
					QuantumsInsuranceLost: 0,
				},
				constants.Dave_Num0: {},
			},
		},
		"Succeeds with maker Long-Term order when considering state fill amount": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short,
				constants.Dave_Num0_1BTC_Long,
			},
			setupMockBankKeeper: func(bk *mocks.BankKeeper) {
				bk.On(
					"SendCoinsFromModuleToModule",
					mock.Anything,
					satypes.ModuleName,
					authtypes.FeeCollectorName,
					mock.MatchedBy(
						testutil_bank.MatchUsdcOfAmount(
							12_500_000+5_000_000,
						),
					),
				).Return(nil).Once()
			},
			setupState: func(t *testing.T, ctx sdk.Context, k *keeper.Keeper) {
				k.SetBlockTimeForLastCommittedBlock(ctx.WithBlockTime(time.Unix(5, 0)))
				k.SetLongTermOrderPlacement(
					ctx,
					constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10,
					4, // Placed in previous block.
				)
				k.SetOrderFillAmount(
					ctx,
					constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10.OrderId,
					50_000_000,
					math.MaxUint32,
				)
			},
			placeOrders: []*types.MsgPlaceOrder{
				{
					Order: constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10,
				},
			},
			clobMatches: []*types.ClobMatch{
				types.NewClobMatchFromMatchOrders(
					&types.MatchOrders{
						MakerOneof: &types.MatchOrders_MakerOrderId{
							MakerOrderId: &constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10.OrderId,
						},
						TakerOneof: &types.MatchOrders_TakerOrderIndex{TakerOrderIndex: 0},
						FillAmount: 50_000_000, // 0.5 BTC
					},
				),
			},
			expectedFillAmounts: map[types.OrderId]satypes.BaseQuantums{
				constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10.OrderId:         50_000_000,
				constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10.OrderId: 100_000_000,
			},
			expectedPruneableBlockHeights: map[uint32][]types.OrderId{
				10 + types.ShortBlockWindow: {
					constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10.OrderId,
				},
			},
			expectedQuoteBalances: map[satypes.SubaccountId]*big.Int{
				// $74,975
				constants.Dave_Num0: big.NewInt(75_000_000_000 - 12_500_000),
				// $74,990
				constants.Carl_Num0: big.NewInt(75_000_000_000 - 5_000_000),
			},
			expectedPerpetualPositions: map[satypes.SubaccountId][]*satypes.PerpetualPosition{
				constants.Carl_Num0: {
					{
						PerpetualId:  0,
						Quantums:     dtypes.NewInt(-50_000_000), // .5 BTC
						FundingIndex: dtypes.ZeroInt(),
					},
				},
				constants.Dave_Num0: {
					{
						PerpetualId:  0,
						Quantums:     dtypes.NewInt(50_000_000), // .5 BTC
						FundingIndex: dtypes.ZeroInt(),
					},
				},
			},
		},
		"Succeeds with taker Long-Term order when considering state fill amount": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short,
				constants.Dave_Num0_1BTC_Long,
			},
			setupMockBankKeeper: func(bk *mocks.BankKeeper) {
				bk.On(
					"SendCoinsFromModuleToModule",
					mock.Anything,
					satypes.ModuleName,
					authtypes.FeeCollectorName,
					mock.MatchedBy(
						testutil_bank.MatchUsdcOfAmount(
							12_500_000+5_000_000,
						),
					),
				).Return(nil).Once()
			},
			setupState: func(t *testing.T, ctx sdk.Context, k *keeper.Keeper) {
				k.SetBlockTimeForLastCommittedBlock(ctx.WithBlockTime(time.Unix(5, 0)))
				k.SetLongTermOrderPlacement(
					ctx,
					constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10,
					4, // Placed in previous block.
				)
				k.SetOrderFillAmount(
					ctx,
					constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10.OrderId,
					50_000_000,
					math.MaxUint32,
				)
			},
			placeOrders: []*types.MsgPlaceOrder{
				{
					Order: constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10,
				},
			},
			clobMatches: []*types.ClobMatch{
				types.NewClobMatchFromMatchOrders(
					&types.MatchOrders{
						MakerOneof: &types.MatchOrders_MakerOrderIndex{MakerOrderIndex: 0},
						TakerOneof: &types.MatchOrders_TakerOrderId{
							TakerOrderId: &constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10.OrderId,
						},
						FillAmount: 50_000_000, // 1 BTC
					},
				),
			},
			expectedFillAmounts: map[types.OrderId]satypes.BaseQuantums{
				constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10.OrderId:         50_000_000,
				constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10.OrderId: 100_000_000,
			},
			expectedPruneableBlockHeights: map[uint32][]types.OrderId{
				10 + types.ShortBlockWindow: {
					constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10.OrderId,
				},
			},
			expectedQuoteBalances: map[satypes.SubaccountId]*big.Int{
				// $74,990
				constants.Dave_Num0: big.NewInt(75_000_000_000 - 5_000_000),
				// $74,975
				constants.Carl_Num0: big.NewInt(75_000_000_000 - 12_500_000),
			},
			expectedPerpetualPositions: map[satypes.SubaccountId][]*satypes.PerpetualPosition{
				constants.Carl_Num0: {
					{
						PerpetualId:  0,
						Quantums:     dtypes.NewInt(-50_000_000), // .5 BTC
						FundingIndex: dtypes.ZeroInt(),
					},
				},
				constants.Dave_Num0: {
					{
						PerpetualId:  0,
						Quantums:     dtypes.NewInt(50_000_000), // .5 BTC
						FundingIndex: dtypes.ZeroInt(),
					},
				},
			},
		},
		"New maker order can overwrite existing order": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short,
				constants.Dave_Num0_1BTC_Long,
			},
			setupMockBankKeeper: func(bk *mocks.BankKeeper) {
				bk.On(
					"SendCoinsFromModuleToModule",
					mock.Anything,
					satypes.ModuleName,
					authtypes.FeeCollectorName,
					mock.MatchedBy(
						testutil_bank.MatchUsdcOfAmount(
							25_000_000+10_000_000,
						),
					),
				).Return(nil).Once()
			},
			setupState: func(t *testing.T, ctx sdk.Context, k *keeper.Keeper) {
				k.SetBlockTimeForLastCommittedBlock(ctx.WithBlockTime(time.Unix(5, 0)))
				k.SetLongTermOrderPlacement(
					ctx,
					types.Order{
						OrderId: types.OrderId{
							SubaccountId: constants.Carl_Num0,
							ClientId:     0,
							OrderFlags:   types.OrderIdFlags_LongTerm,
							ClobPairId:   0,
						},
						Side:         types.Order_SIDE_BUY,
						Quantums:     99_000_000,     // 0.99 BTC
						Subticks:     49_999_000_000, // $49,999
						GoodTilOneof: &types.Order_GoodTilBlockTime{GoodTilBlockTime: 9},
					},
					4, // Placed in a previous block.
				)
			},
			placeOrders: []*types.MsgPlaceOrder{
				{
					Order: constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10,
				},
				{
					Order: constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10,
				},
			},
			clobMatches: []*types.ClobMatch{
				types.NewClobMatchFromMatchOrders(
					&types.MatchOrders{
						MakerOneof: &types.MatchOrders_MakerOrderId{
							MakerOrderId: &constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10.OrderId,
						},
						TakerOneof: &types.MatchOrders_TakerOrderIndex{TakerOrderIndex: 0},
						FillAmount: 100_000_000, // 1 BTC
					},
				),
			},
			expectedFillAmounts: map[types.OrderId]satypes.BaseQuantums{
				constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10.OrderId:         100_000_000,
				constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10.OrderId: 100_000_000,
			},
			expectedPruneableBlockHeights: map[uint32][]types.OrderId{
				10 + types.ShortBlockWindow: {
					constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10.OrderId,
				},
			},
			expectedQuoteBalances: map[satypes.SubaccountId]*big.Int{
				// $99,975
				constants.Dave_Num0: big.NewInt(100_000_000_000 - 25_000_000),
				// $49,990
				constants.Carl_Num0: big.NewInt(50_000_000_000 - 10_000_000),
			},
			expectedPerpetualPositions: map[satypes.SubaccountId][]*satypes.PerpetualPosition{
				constants.Carl_Num0: {},
				constants.Dave_Num0: {},
			},
		},
		"New taker order can overwrite existing order": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short,
				constants.Dave_Num0_1BTC_Long,
			},
			setupMockBankKeeper: func(bk *mocks.BankKeeper) {
				bk.On(
					"SendCoinsFromModuleToModule",
					mock.Anything,
					satypes.ModuleName,
					authtypes.FeeCollectorName,
					mock.MatchedBy(
						testutil_bank.MatchUsdcOfAmount(
							25_000_000+10_000_000,
						),
					),
				).Return(nil).Once()
			},
			setupState: func(t *testing.T, ctx sdk.Context, k *keeper.Keeper) {
				k.SetBlockTimeForLastCommittedBlock(ctx.WithBlockTime(time.Unix(5, 0)))
				k.SetLongTermOrderPlacement(
					ctx,
					types.Order{
						OrderId: types.OrderId{
							SubaccountId: constants.Carl_Num0,
							ClientId:     0,
							OrderFlags:   types.OrderIdFlags_LongTerm,
							ClobPairId:   0,
						},
						Side:         types.Order_SIDE_BUY,
						Quantums:     99_000_000,     // 0.99 BTC
						Subticks:     49_999_000_000, // $49,999
						GoodTilOneof: &types.Order_GoodTilBlockTime{GoodTilBlockTime: 9},
					},
					4, // Placed in a previous block.
				)
			},
			placeOrders: []*types.MsgPlaceOrder{
				{
					Order: constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10,
				},
				{
					Order: constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10,
				},
			},
			clobMatches: []*types.ClobMatch{
				types.NewClobMatchFromMatchOrders(
					&types.MatchOrders{
						MakerOneof: &types.MatchOrders_MakerOrderIndex{MakerOrderIndex: 0},
						TakerOneof: &types.MatchOrders_TakerOrderId{
							TakerOrderId: &constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10.OrderId,
						},
						FillAmount: 100_000_000, // 1 BTC
					},
				),
			},
			expectedFillAmounts: map[types.OrderId]satypes.BaseQuantums{
				constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10.OrderId:         100_000_000,
				constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10.OrderId: 100_000_000,
			},
			expectedPruneableBlockHeights: map[uint32][]types.OrderId{
				10 + types.ShortBlockWindow: {
					constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10.OrderId,
				},
			},
			expectedQuoteBalances: map[satypes.SubaccountId]*big.Int{
				// $99,990
				constants.Dave_Num0: big.NewInt(100_000_000_000 - 10_000_000),
				// $49,975
				constants.Carl_Num0: big.NewInt(50_000_000_000 - 25_000_000),
			},
			expectedPerpetualPositions: map[satypes.SubaccountId][]*satypes.PerpetualPosition{
				constants.Carl_Num0: {},
				constants.Dave_Num0: {},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			runProcessProposerMatchSuccessTest(t, tc)
		})
	}
}

func TestProcessProposerMatches_LongTerm_Validation_Failure(t *testing.T) {
	tests := map[string]processProposerMatchesTestCase{
		`Stateful order validation: Long-term order GoodTilBlockTime is less than the
		block time of the previous block`: {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short,
			},
			setupState: func(t *testing.T, ctx sdk.Context, k *keeper.Keeper) {
				k.SetBlockTimeForLastCommittedBlock(ctx)
			},
			placeOrders: []*types.MsgPlaceOrder{
				{
					Order: types.Order{
						OrderId: types.OrderId{
							SubaccountId: constants.Carl_Num0,
							ClientId:     0,
							OrderFlags:   types.OrderIdFlags_LongTerm,
							ClobPairId:   0,
						},
						Side:     types.Order_SIDE_BUY,
						Quantums: 10,
						Subticks: 100,
						GoodTilOneof: &types.Order_GoodTilBlockTime{
							// Block-time of the previously committed block is 5.
							GoodTilBlockTime: 4,
						},
					},
				},
			},
			clobMatches:   []*types.ClobMatch{},
			expectedError: types.ErrTimeExceedsGoodTilBlockTime,
		},
		`Stateful order validation: Long-term order GoodTilBlockTime exceeds StatefulOrderTimeWindow`: {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short,
			},
			setupState: func(t *testing.T, ctx sdk.Context, k *keeper.Keeper) {
				k.SetBlockTimeForLastCommittedBlock(ctx)
			},
			placeOrders: []*types.MsgPlaceOrder{
				{
					Order: types.Order{
						OrderId: types.OrderId{
							SubaccountId: constants.Carl_Num0,
							ClientId:     0,
							OrderFlags:   types.OrderIdFlags_LongTerm,
							ClobPairId:   0,
						},
						Side:     types.Order_SIDE_BUY,
						Quantums: 10,
						Subticks: 100,
						GoodTilOneof: &types.Order_GoodTilBlockTime{
							GoodTilBlockTime: 5 + uint32(types.StatefulOrderTimeWindow.Seconds()) + 1,
						},
					},
				},
			},
			clobMatches:   []*types.ClobMatch{},
			expectedError: types.ErrGoodTilBlockTimeExceedsStatefulOrderTimeWindow,
		},
		`Stateful order validation: Long-term order already exists in state`: {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short,
			},
			setupState: func(t *testing.T, ctx sdk.Context, k *keeper.Keeper) {
				k.SetBlockTimeForLastCommittedBlock(ctx)
				k.SetLongTermOrderPlacement(
					ctx,
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
					lib.MustConvertIntegerToUint32(ctx.BlockHeight()),
				)
			},
			placeOrders: []*types.MsgPlaceOrder{
				{
					Order: constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
				},
			},
			clobMatches:   []*types.ClobMatch{},
			expectedError: types.ErrStatefulOrderAlreadyExists,
		},
		`Stateful order validation: GoodTilBlockTime of new order is less than existing order`: {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short,
			},
			setupState: func(t *testing.T, ctx sdk.Context, k *keeper.Keeper) {
				k.SetBlockTimeForLastCommittedBlock(ctx)
				k.SetLongTermOrderPlacement(
					ctx,
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20, // GTBT is 20.
					lib.MustConvertIntegerToUint32(ctx.BlockHeight()),
				)
			},
			placeOrders: []*types.MsgPlaceOrder{
				{
					// GTBT of the new order is 15.
					Order: constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
				},
			},
			clobMatches:   []*types.ClobMatch{},
			expectedError: types.ErrStatefulOrderAlreadyExists,
		},
		`Stateful order validation: referenced maker order does not exist in state`: {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short,
				constants.Dave_Num0_1BTC_Long,
			},
			placeOrders: []*types.MsgPlaceOrder{
				{
					Order: constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10,
				},
			},
			clobMatches: []*types.ClobMatch{
				types.NewClobMatchFromMatchOrders(
					&types.MatchOrders{
						// Maker order is a long-term order.
						MakerOneof: &types.MatchOrders_MakerOrderId{
							MakerOrderId: &constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10.OrderId,
						},
						// Taker order is a short-term order.
						TakerOneof: &types.MatchOrders_TakerOrderIndex{TakerOrderIndex: 0},
						FillAmount: 100_000_000, // 1 BTC
					},
				),
			},
			expectedError: types.ErrStatefulOrderDoesNotExist,
		},
		`Stateful order validation: referenced taker order does not exist in state`: {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short,
				constants.Dave_Num0_1BTC_Long,
			},
			placeOrders: []*types.MsgPlaceOrder{
				{
					Order: constants.Order_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTB10,
				},
			},
			clobMatches: []*types.ClobMatch{
				types.NewClobMatchFromMatchOrders(
					&types.MatchOrders{
						// Maker order is a short-term order.
						MakerOneof: &types.MatchOrders_MakerOrderIndex{MakerOrderIndex: 0},
						// Taker order is a long-term order.
						TakerOneof: &types.MatchOrders_TakerOrderId{
							TakerOrderId: &constants.LongTermOrder_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10.OrderId,
						},
						FillAmount: 100_000_000, // 1 BTC
					},
				),
			},
			expectedError: types.ErrStatefulOrderDoesNotExist,
		},
		`Stateful order validation: referenced maker order in liquidation match does not exist in state`: {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_54999USD,
				constants.Dave_Num0_1BTC_Long,
			},
			placeOrders: []*types.MsgPlaceOrder{},
			setupState: func(t *testing.T, ctx sdk.Context, k *keeper.Keeper) {
				require.NoError(
					t,
					k.InitializeLiquidationsConfig(
						ctx,
						constants.LiquidationsConfig_No_Limit,
					),
				)
			},
			clobMatches: []*types.ClobMatch{
				types.NewClobMatchFromMatchPerpetualLiquidation(
					&types.MatchPerpetualLiquidation{
						Liquidated:  constants.Carl_Num0,
						ClobPairId:  0,
						PerpetualId: 0,
						TotalSize:   100_000_000, // 1 BTC
						IsBuy:       true,
						Fills: []types.MatchPerpetualLiquidation_Fill{
							{
								// Maker order is a long-term order.
								MakerOneof: &types.MatchPerpetualLiquidation_Fill_MakerOrderId{
									MakerOrderId: &constants.LongTermOrder_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10.OrderId,
								},
								FillAmount: 100_000_000, // 1 BTC
							},
						},
					},
				),
			},
			expectedError: types.ErrStatefulOrderDoesNotExist,
		},
		`Stateful order validation: referenced long-term order is on the wrong side`: {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short,
				constants.Dave_Num0_1BTC_Long,
			},
			placeOrders: []*types.MsgPlaceOrder{
				{
					Order: constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10,
				},
			},
			setupState: func(t *testing.T, ctx sdk.Context, k *keeper.Keeper) {
				k.SetBlockTimeForLastCommittedBlock(ctx)
				k.SetLongTermOrderPlacement(
					ctx,
					types.Order{
						OrderId: types.OrderId{
							SubaccountId: constants.Carl_Num0,
							ClientId:     0,
							OrderFlags:   types.OrderIdFlags_LongTerm,
							ClobPairId:   0,
						},
						Side:         types.Order_SIDE_SELL, // This is a sell order instead of a buy order.
						Quantums:     100_000_000,
						Subticks:     50_000_000_000,
						GoodTilOneof: &types.Order_GoodTilBlockTime{GoodTilBlockTime: 10},
					},
					lib.MustConvertIntegerToUint32(ctx.BlockHeight()),
				)
			},
			clobMatches: []*types.ClobMatch{
				types.NewClobMatchFromMatchOrders(
					&types.MatchOrders{
						// Maker order is a long-term order.
						MakerOneof: &types.MatchOrders_MakerOrderId{
							MakerOrderId: &types.OrderId{
								SubaccountId: constants.Carl_Num0,
								ClientId:     0,
								OrderFlags:   types.OrderIdFlags_LongTerm,
							},
						},
						TakerOneof: &types.MatchOrders_TakerOrderIndex{TakerOrderIndex: 0},
						FillAmount: 100_000_000, // 1 BTC
					},
				),
			},
			expectedError: errors.New("Orders are not on opposing sides of the book in match"),
		},
		`Stateful order validation: referenced long-term order is for the wrong clob pair`: {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short,
				constants.Dave_Num0_1BTC_Long,
			},
			placeOrders: []*types.MsgPlaceOrder{
				{
					// This is a BTC order.
					Order: constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10,
				},
			},
			setupState: func(t *testing.T, ctx sdk.Context, k *keeper.Keeper) {
				k.SetBlockTimeForLastCommittedBlock(ctx)
				k.SetLongTermOrderPlacement(
					ctx,
					types.Order{
						OrderId: types.OrderId{
							SubaccountId: constants.Carl_Num0,
							ClientId:     0,
							OrderFlags:   types.OrderIdFlags_LongTerm,
							ClobPairId:   1, // ETH.
						},
						Side:         types.Order_SIDE_BUY,
						Quantums:     100_000_000,
						Subticks:     50_000_000_000,
						GoodTilOneof: &types.Order_GoodTilBlockTime{GoodTilBlockTime: 10},
					},
					lib.MustConvertIntegerToUint32(ctx.BlockHeight()),
				)
			},
			clobMatches: []*types.ClobMatch{
				types.NewClobMatchFromMatchOrders(
					&types.MatchOrders{
						// Maker order is a long-term order.
						MakerOneof: &types.MatchOrders_MakerOrderId{
							MakerOrderId: &types.OrderId{
								SubaccountId: constants.Carl_Num0,
								ClientId:     0,
								OrderFlags:   types.OrderIdFlags_LongTerm,
								ClobPairId:   1, // ETH.
							},
						},
						TakerOneof: &types.MatchOrders_TakerOrderIndex{TakerOrderIndex: 0},
						FillAmount: 100_000_000, // 1 BTC
					},
				),
			},
			expectedError: errors.New("ClobPairIds do not match in match"),
		},
		"Fails with Long-Term order when considering state fill amount": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short,
				constants.Dave_Num0_1BTC_Long,
			},
			setupMockBankKeeper: func(bk *mocks.BankKeeper) {
				bk.On(
					"SendCoinsFromModuleToModule",
					mock.Anything,
					satypes.ModuleName,
					authtypes.FeeCollectorName,
					mock.MatchedBy(
						testutil_bank.MatchUsdcOfAmount(
							12_500_000+5_000_000,
						),
					),
				).Return(nil).Once()
			},
			setupState: func(t *testing.T, ctx sdk.Context, k *keeper.Keeper) {
				k.SetBlockTimeForLastCommittedBlock(ctx.WithBlockTime(time.Unix(5, 0)))
				k.SetLongTermOrderPlacement(
					ctx,
					constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10,
					4, // Placed in previous block.
				)
				k.SetOrderFillAmount(
					ctx,
					constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10.OrderId,
					50_000_001,
					math.MaxUint32,
				)
			},
			placeOrders: []*types.MsgPlaceOrder{
				{
					Order: constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10,
				},
			},
			clobMatches: []*types.ClobMatch{
				types.NewClobMatchFromMatchOrders(
					&types.MatchOrders{
						MakerOneof: &types.MatchOrders_MakerOrderId{
							MakerOrderId: &constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10.OrderId,
						},
						TakerOneof: &types.MatchOrders_TakerOrderIndex{TakerOrderIndex: 0},
						FillAmount: 50_000_000,
					},
				),
			},
			expectedError: fmt.Errorf(
				"Match with Quantums 50000000 would exceed total Quantums 100000000 of "+
					"OrderId %v. New total filled quantums would be 100000001",
				constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10.OrderId,
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			runProcessProposerMatchFailureTest(t, tc)
		})
	}
}

func TestProcessProposerMatches_Conditional_Validation_Failure(t *testing.T) {
	tests := map[string]processProposerMatchesTestCase{
		`Stateful order validation: Conditional order GoodTilBlockTime is less than the
		block time of the previous block`: {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short,
			},
			setupState: func(t *testing.T, ctx sdk.Context, k *keeper.Keeper) {
				k.SetBlockTimeForLastCommittedBlock(ctx)
			},
			placeOrders: []*types.MsgPlaceOrder{
				{
					Order: types.Order{
						OrderId: types.OrderId{
							SubaccountId: constants.Carl_Num0,
							ClientId:     0,
							OrderFlags:   types.OrderIdFlags_Conditional,
							ClobPairId:   0,
						},
						Side:     types.Order_SIDE_BUY,
						Quantums: 10,
						Subticks: 100,
						GoodTilOneof: &types.Order_GoodTilBlockTime{
							// Block-time of the previously committed block is 5.
							GoodTilBlockTime: 4,
						},
					},
				},
			},
			clobMatches:   []*types.ClobMatch{},
			expectedError: types.ErrTimeExceedsGoodTilBlockTime,
		},
		`Stateful order validation: Conditional order GoodTilBlockTime exceeds StatefulOrderTimeWindow`: {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short,
			},
			setupState: func(t *testing.T, ctx sdk.Context, k *keeper.Keeper) {
				k.SetBlockTimeForLastCommittedBlock(ctx)
			},
			placeOrders: []*types.MsgPlaceOrder{
				{
					Order: types.Order{
						OrderId: types.OrderId{
							SubaccountId: constants.Carl_Num0,
							ClientId:     0,
							OrderFlags:   types.OrderIdFlags_Conditional,
							ClobPairId:   0,
						},
						Side:     types.Order_SIDE_BUY,
						Quantums: 10,
						Subticks: 100,
						GoodTilOneof: &types.Order_GoodTilBlockTime{
							GoodTilBlockTime: 5 + uint32(types.StatefulOrderTimeWindow.Seconds()) + 1,
						},
					},
				},
			},
			clobMatches:   []*types.ClobMatch{},
			expectedError: types.ErrGoodTilBlockTimeExceedsStatefulOrderTimeWindow,
		},
		`Stateful order validation: Conditional order already exists in state`: {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short,
			},
			setupState: func(t *testing.T, ctx sdk.Context, k *keeper.Keeper) {
				k.SetBlockTimeForLastCommittedBlock(ctx)
				k.SetLongTermOrderPlacement(
					ctx,
					constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
					lib.MustConvertIntegerToUint32(ctx.BlockHeight()),
				)
			},
			placeOrders: []*types.MsgPlaceOrder{
				{
					Order: constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
				},
			},
			clobMatches:   []*types.ClobMatch{},
			expectedError: types.ErrStatefulOrderAlreadyExists,
		},
		`Stateful order validation: GoodTilBlockTime of new order is less than existing order`: {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short,
			},
			setupState: func(t *testing.T, ctx sdk.Context, k *keeper.Keeper) {
				k.SetBlockTimeForLastCommittedBlock(ctx)
				k.SetLongTermOrderPlacement(
					ctx,
					constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20, // GTBT is 20.
					lib.MustConvertIntegerToUint32(ctx.BlockHeight()),
				)
			},
			placeOrders: []*types.MsgPlaceOrder{
				{
					// GTBT of the new order is 15.
					Order: constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
				},
			},
			clobMatches:   []*types.ClobMatch{},
			expectedError: types.ErrStatefulOrderAlreadyExists,
		},
		`Stateful order validation: referenced maker order does not exist in state`: {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short,
				constants.Dave_Num0_1BTC_Long,
			},
			placeOrders: []*types.MsgPlaceOrder{
				{
					Order: constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10,
				},
			},
			clobMatches: []*types.ClobMatch{
				types.NewClobMatchFromMatchOrders(
					&types.MatchOrders{
						// Maker order is a conditional order.
						MakerOneof: &types.MatchOrders_MakerOrderId{
							MakerOrderId: &constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10.OrderId,
						},
						// Taker order is a short-term order.
						TakerOneof: &types.MatchOrders_TakerOrderIndex{TakerOrderIndex: 0},
						FillAmount: 100_000_000, // 1 BTC
					},
				),
			},
			expectedError: types.ErrStatefulOrderDoesNotExist,
		},
		`Stateful order validation: referenced taker order does not exist in state`: {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short,
				constants.Dave_Num0_1BTC_Long,
			},
			placeOrders: []*types.MsgPlaceOrder{
				{
					Order: constants.Order_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTB10,
				},
			},
			clobMatches: []*types.ClobMatch{
				types.NewClobMatchFromMatchOrders(
					&types.MatchOrders{
						// Maker order is a short-term order.
						MakerOneof: &types.MatchOrders_MakerOrderIndex{MakerOrderIndex: 0},
						// Taker order is a conditional order.
						TakerOneof: &types.MatchOrders_TakerOrderId{
							TakerOrderId: &constants.ConditionalOrder_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10.OrderId,
						},
						FillAmount: 100_000_000, // 1 BTC
					},
				),
			},
			expectedError: types.ErrStatefulOrderDoesNotExist,
		},
		`Stateful order validation: referenced maker order in liquidation match does not exist in state`: {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_54999USD,
				constants.Dave_Num0_1BTC_Long,
			},
			placeOrders: []*types.MsgPlaceOrder{},
			setupState: func(t *testing.T, ctx sdk.Context, k *keeper.Keeper) {
				require.NoError(
					t,
					k.InitializeLiquidationsConfig(
						ctx,
						constants.LiquidationsConfig_No_Limit,
					),
				)
			},
			clobMatches: []*types.ClobMatch{
				types.NewClobMatchFromMatchPerpetualLiquidation(
					&types.MatchPerpetualLiquidation{
						Liquidated:  constants.Carl_Num0,
						ClobPairId:  0,
						PerpetualId: 0,
						TotalSize:   100_000_000, // 1 BTC
						IsBuy:       true,
						Fills: []types.MatchPerpetualLiquidation_Fill{
							{
								// Maker order is a conditional order.
								MakerOneof: &types.MatchPerpetualLiquidation_Fill_MakerOrderId{
									MakerOrderId: &constants.ConditionalOrder_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10.OrderId,
								},
								FillAmount: 100_000_000, // 1 BTC
							},
						},
					},
				),
			},
			expectedError: types.ErrStatefulOrderDoesNotExist,
		},
		`Stateful order validation: referenced conditional order is on the wrong side`: {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short,
				constants.Dave_Num0_1BTC_Long,
			},
			placeOrders: []*types.MsgPlaceOrder{
				{
					Order: constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10,
				},
			},
			setupState: func(t *testing.T, ctx sdk.Context, k *keeper.Keeper) {
				k.SetBlockTimeForLastCommittedBlock(ctx)
				k.SetLongTermOrderPlacement(
					ctx,
					types.Order{
						OrderId: types.OrderId{
							SubaccountId: constants.Carl_Num0,
							ClientId:     0,
							OrderFlags:   types.OrderIdFlags_Conditional,
							ClobPairId:   0,
						},
						Side:         types.Order_SIDE_SELL, // This is a sell order instead of a buy order.
						Quantums:     100_000_000,
						Subticks:     50_000_000_000,
						GoodTilOneof: &types.Order_GoodTilBlockTime{GoodTilBlockTime: 10},
					},
					lib.MustConvertIntegerToUint32(ctx.BlockHeight()),
				)
			},
			clobMatches: []*types.ClobMatch{
				types.NewClobMatchFromMatchOrders(
					&types.MatchOrders{
						// Maker order is a conditional order.
						MakerOneof: &types.MatchOrders_MakerOrderId{
							MakerOrderId: &types.OrderId{
								SubaccountId: constants.Carl_Num0,
								ClientId:     0,
								OrderFlags:   types.OrderIdFlags_Conditional,
							},
						},
						TakerOneof: &types.MatchOrders_TakerOrderIndex{TakerOrderIndex: 0},
						FillAmount: 100_000_000, // 1 BTC
					},
				),
			},
			expectedError: errors.New("Orders are not on opposing sides of the book in match"),
		},
		`Stateful order validation: referenced conditional order is for the wrong clob pair`: {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short,
				constants.Dave_Num0_1BTC_Long,
			},
			placeOrders: []*types.MsgPlaceOrder{
				{
					// This is a BTC order.
					Order: constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10,
				},
			},
			setupState: func(t *testing.T, ctx sdk.Context, k *keeper.Keeper) {
				k.SetBlockTimeForLastCommittedBlock(ctx)
				k.SetLongTermOrderPlacement(
					ctx,
					types.Order{
						OrderId: types.OrderId{
							SubaccountId: constants.Carl_Num0,
							ClientId:     0,
							OrderFlags:   types.OrderIdFlags_Conditional,
							ClobPairId:   1, // ETH.
						},
						Side:         types.Order_SIDE_BUY,
						Quantums:     100_000_000,
						Subticks:     50_000_000_000,
						GoodTilOneof: &types.Order_GoodTilBlockTime{GoodTilBlockTime: 10},
					},
					lib.MustConvertIntegerToUint32(ctx.BlockHeight()),
				)
			},
			clobMatches: []*types.ClobMatch{
				types.NewClobMatchFromMatchOrders(
					&types.MatchOrders{
						// Maker order is a conditional order.
						MakerOneof: &types.MatchOrders_MakerOrderId{
							MakerOrderId: &types.OrderId{
								SubaccountId: constants.Carl_Num0,
								ClientId:     0,
								OrderFlags:   types.OrderIdFlags_Conditional,
								ClobPairId:   1, // ETH.
							},
						},
						TakerOneof: &types.MatchOrders_TakerOrderIndex{TakerOrderIndex: 0},
						FillAmount: 100_000_000, // 1 BTC
					},
				),
			},
			expectedError: errors.New("ClobPairIds do not match in match"),
		},
		"Fails with conditional order when considering state fill amount": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short,
				constants.Dave_Num0_1BTC_Long,
			},
			setupMockBankKeeper: func(bk *mocks.BankKeeper) {
				bk.On(
					"SendCoinsFromModuleToModule",
					mock.Anything,
					satypes.ModuleName,
					authtypes.FeeCollectorName,
					mock.MatchedBy(
						testutil_bank.MatchUsdcOfAmount(
							12_500_000+5_000_000,
						),
					),
				).Return(nil).Once()
			},
			setupState: func(t *testing.T, ctx sdk.Context, k *keeper.Keeper) {
				k.SetBlockTimeForLastCommittedBlock(ctx.WithBlockTime(time.Unix(5, 0)))
				k.SetLongTermOrderPlacement(
					ctx,
					constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10,
					4, // Placed in previous block.
				)
				k.SetOrderFillAmount(
					ctx,
					constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10.OrderId,
					50_000_001,
					math.MaxUint32,
				)
			},
			placeOrders: []*types.MsgPlaceOrder{
				{
					Order: constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10,
				},
			},
			clobMatches: []*types.ClobMatch{
				types.NewClobMatchFromMatchOrders(
					&types.MatchOrders{
						MakerOneof: &types.MatchOrders_MakerOrderId{
							MakerOrderId: &constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10.OrderId,
						},
						TakerOneof: &types.MatchOrders_TakerOrderIndex{TakerOrderIndex: 0},
						FillAmount: 50_000_000,
					},
				),
			},
			expectedError: fmt.Errorf(
				"Match with Quantums 50000000 would exceed total Quantums 100000000 of "+
					"OrderId %v. New total filled quantums would be 100000001",
				constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10.OrderId,
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			runProcessProposerMatchFailureTest(t, tc)
		})
	}
}

*/
