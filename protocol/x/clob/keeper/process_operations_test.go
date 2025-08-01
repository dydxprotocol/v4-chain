package keeper_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/dydxprotocol/v4-chain/protocol/app/module"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	indexerevents "github.com/dydxprotocol/v4-chain/protocol/indexer/events"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/shared"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	clobtest "github.com/dydxprotocol/v4-chain/protocol/testutil/clob"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	perptest "github.com/dydxprotocol/v4-chain/protocol/testutil/perpetuals"
	testutil "github.com/dydxprotocol/v4-chain/protocol/testutil/util"
	blocktimetypes "github.com/dydxprotocol/v4-chain/protocol/x/blocktime/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/memclob"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	feetierstypes "github.com/dydxprotocol/v4-chain/protocol/x/feetiers/types"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	sakeeper "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/keeper"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MatchWithOrdersForTesting represents a match which occurred between two orders and the amount that was matched.
type MatchWithOrdersForTesting struct {
	types.MatchWithOrders
	TotalFilledMaker satypes.BaseQuantums
	TotalFilledTaker satypes.BaseQuantums
}

type processProposerOperationsTestCase struct {
	// State
	perpetuals                    []perptypes.Perpetual
	perpetualFeeParams            *feetierstypes.PerpetualFeeParams
	clobPairs                     []types.ClobPair
	subaccounts                   []satypes.Subaccount
	preExistingStatefulOrders     []types.Order
	triggeredConditionalOrders    []types.Order
	marketIdToOraclePriceOverride map[uint32]uint64
	rawOperations                 []types.OperationRaw

	setupState          func(ctx sdk.Context, ks keepertest.ClobKeepersTestContext)
	setupMockBankKeeper func(bk *mocks.BankKeeper)

	// Liquidation specific setup.
	liquidationConfig    *types.LiquidationsConfig
	insuranceFundBalance uint64

	// Expectations.
	// Note that for expectedProcessProposerMatchesEvents, the OperationsProposedInLastBlock field is populated from
	// the operations field above.
	expectedProcessProposerMatchesEvents types.ProcessProposerMatchesEvents
	expectedMatches                      []*MatchWithOrdersForTesting
	expectedFillAmounts                  map[types.OrderId]satypes.BaseQuantums
	expectedQuoteBalances                map[satypes.SubaccountId]int64
	expectedPerpetualPositions           map[satypes.SubaccountId][]*satypes.PerpetualPosition
	expectedSubaccountLiquidationInfo    map[satypes.SubaccountId]types.SubaccountLiquidationInfo
	expectedNegativeTncSubaccountSeen    map[uint32]bool
	expectedError                        error
	expectedPanics                       string
}

func TestProcessProposerOperations(t *testing.T) {
	blockHeight := uint32(5)
	tests := map[string]processProposerOperationsTestCase{
		"Succeeds no operations": {
			perpetuals:                []perptypes.Perpetual{},
			perpetualFeeParams:        &constants.PerpetualFeeParams,
			clobPairs:                 []types.ClobPair{},
			subaccounts:               []satypes.Subaccount{},
			preExistingStatefulOrders: []types.Order{},
			rawOperations:             []types.OperationRaw{},

			expectedProcessProposerMatchesEvents: types.ProcessProposerMatchesEvents{
				BlockHeight: blockHeight,
			},
		},
		"Succeeds no operations with previous stateful orders": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			perpetualFeeParams: &constants.PerpetualFeeParams,
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc,
			},
			subaccounts: []satypes.Subaccount{},
			preExistingStatefulOrders: []types.Order{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20,
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
			},
			rawOperations: []types.OperationRaw{},

			expectedProcessProposerMatchesEvents: types.ProcessProposerMatchesEvents{
				BlockHeight: blockHeight,
			},
		},
		"Succeeds with singular match of a short term maker and short term taker": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			perpetualFeeParams: &constants.PerpetualFeeParams,
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc,
			},
			subaccounts: []satypes.Subaccount{
				{
					Id: &constants.Alice_Num0,
					AssetPositions: []*satypes.AssetPosition{
						&constants.Usdc_Asset_100_000,
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						testutil.CreateSinglePerpetualPosition(
							0,
							big.NewInt(1_000_000_000), // 10 BTC
							big.NewInt(0),
							big.NewInt(0),
						),
					},
				},
				{
					Id: &constants.Bob_Num0,
					AssetPositions: []*satypes.AssetPosition{
						&constants.Usdc_Asset_100_000,
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						testutil.CreateSinglePerpetualPosition(
							0,
							big.NewInt(1_000_000_000), // 10 BTC
							big.NewInt(0),
							big.NewInt(0),
						),
					},
				},
			},
			preExistingStatefulOrders: []types.Order{},
			rawOperations: []types.OperationRaw{
				clobtest.NewShortTermOrderPlacementOperationRaw(
					types.Order{
						OrderId:      types.OrderId{SubaccountId: constants.Alice_Num0, ClientId: 14, ClobPairId: 0},
						Side:         types.Order_SIDE_BUY,
						Quantums:     100_000_000, // 1 BTC
						Subticks:     50_000_000,
						GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 25},
					},
				),
				clobtest.NewShortTermOrderPlacementOperationRaw(
					types.Order{
						OrderId:      types.OrderId{SubaccountId: constants.Bob_Num0, ClientId: 14, ClobPairId: 0},
						Side:         types.Order_SIDE_SELL,
						Quantums:     100_000_000, // 1 BTC
						Subticks:     50_000_000,
						GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 25},
					},
				),
				clobtest.NewMatchOperationRaw(
					&types.Order{
						OrderId:      types.OrderId{SubaccountId: constants.Bob_Num0, ClientId: 14, ClobPairId: 0},
						Side:         types.Order_SIDE_SELL,
						Quantums:     100_000_000, // 1 BTC
						Subticks:     50_000_000,
						GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 25},
					},
					[]types.MakerFill{
						{
							FillAmount:   100_000_000,
							MakerOrderId: types.OrderId{SubaccountId: constants.Alice_Num0, ClientId: 14, ClobPairId: 0},
						},
					},
				),
			},
			expectedMatches: []*MatchWithOrdersForTesting{
				{
					MatchWithOrders: types.MatchWithOrders{
						TakerOrder: &types.Order{
							OrderId:      types.OrderId{SubaccountId: constants.Bob_Num0, ClientId: 14, ClobPairId: 0},
							Side:         types.Order_SIDE_SELL,
							Quantums:     100_000_000,
							Subticks:     50_000_000,
							GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 25},
						},
						MakerOrder: &types.Order{
							OrderId:      types.OrderId{SubaccountId: constants.Alice_Num0, ClientId: 14, ClobPairId: 0},
							Side:         types.Order_SIDE_BUY,
							Quantums:     100_000_000,
							Subticks:     50_000_000,
							GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 25},
						},
						FillAmount: 100_000_000,
						MakerFee:   10_000,
						TakerFee:   25_000,
					},
					TotalFilledMaker: 100_000_000,
					TotalFilledTaker: 100_000_000,
				},
			},
			expectedProcessProposerMatchesEvents: types.ProcessProposerMatchesEvents{
				OrderIdsFilledInLastBlock: []types.OrderId{
					{SubaccountId: constants.Bob_Num0, ClientId: 14, ClobPairId: 0},
					{SubaccountId: constants.Alice_Num0, ClientId: 14, ClobPairId: 0},
				},
				BlockHeight: blockHeight,
			},
			// Expected balances are initial balance + balance change due to order - fees
			expectedQuoteBalances: map[satypes.SubaccountId]int64{
				constants.Alice_Num0: constants.Usdc_Asset_100_000.GetBigQuantums().Int64() - 50_000_000 - 10_000,
				constants.Bob_Num0:   constants.Usdc_Asset_100_000.GetBigQuantums().Int64() + 50_000_000 - 25_000,
			},
			expectedPerpetualPositions: map[satypes.SubaccountId][]*satypes.PerpetualPosition{
				constants.Bob_Num0: {
					testutil.CreateSinglePerpetualPosition(
						0,
						big.NewInt(1_000_000_000-100_000_000),
						big.NewInt(0),
						big.NewInt(0),
					),
				},
				constants.Alice_Num0: {
					testutil.CreateSinglePerpetualPosition(
						0,
						big.NewInt(1_000_000_000+100_000_000),
						big.NewInt(0),
						big.NewInt(0),
					),
				},
			},
		},
		"Succeeds with maker rebate": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			perpetualFeeParams: &constants.PerpetualFeeParamsMakerRebate,
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc,
			},
			subaccounts: []satypes.Subaccount{
				{
					Id: &constants.Alice_Num0,
					AssetPositions: []*satypes.AssetPosition{
						&constants.Usdc_Asset_100_000,
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						testutil.CreateSinglePerpetualPosition(
							0,
							big.NewInt(1_000_000_000), // 10 BTC
							big.NewInt(0),
							big.NewInt(0),
						),
					},
				},
				{
					Id: &constants.Bob_Num0,
					AssetPositions: []*satypes.AssetPosition{
						&constants.Usdc_Asset_100_000,
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						testutil.CreateSinglePerpetualPosition(
							0,
							big.NewInt(1_000_000_000), // 10 BTC
							big.NewInt(0),
							big.NewInt(0),
						),
					},
				},
			},
			preExistingStatefulOrders: []types.Order{},
			rawOperations: []types.OperationRaw{
				clobtest.NewShortTermOrderPlacementOperationRaw(
					types.Order{
						OrderId:      types.OrderId{SubaccountId: constants.Alice_Num0, ClientId: 14, ClobPairId: 0},
						Side:         types.Order_SIDE_BUY,
						Quantums:     100_000_000, // 1 BTC
						Subticks:     50_000_000,
						GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 25},
					},
				),
				clobtest.NewShortTermOrderPlacementOperationRaw(
					types.Order{
						OrderId:      types.OrderId{SubaccountId: constants.Bob_Num0, ClientId: 14, ClobPairId: 0},
						Side:         types.Order_SIDE_SELL,
						Quantums:     100_000_000, // 1 BTC
						Subticks:     50_000_000,
						GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 25},
					},
				),
				clobtest.NewMatchOperationRaw(
					&types.Order{
						OrderId:      types.OrderId{SubaccountId: constants.Bob_Num0, ClientId: 14, ClobPairId: 0},
						Side:         types.Order_SIDE_SELL,
						Quantums:     100_000_000, // 1 BTC
						Subticks:     50_000_000,
						GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 25},
					},
					[]types.MakerFill{
						{
							FillAmount:   100_000_000,
							MakerOrderId: types.OrderId{SubaccountId: constants.Alice_Num0, ClientId: 14, ClobPairId: 0},
						},
					},
				),
			},
			expectedMatches: []*MatchWithOrdersForTesting{
				{
					MatchWithOrders: types.MatchWithOrders{
						TakerOrder: &types.Order{
							OrderId:      types.OrderId{SubaccountId: constants.Bob_Num0, ClientId: 14, ClobPairId: 0},
							Side:         types.Order_SIDE_SELL,
							Quantums:     100_000_000,
							Subticks:     50_000_000,
							GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 25},
						},
						MakerOrder: &types.Order{
							OrderId:      types.OrderId{SubaccountId: constants.Alice_Num0, ClientId: 14, ClobPairId: 0},
							Side:         types.Order_SIDE_BUY,
							Quantums:     100_000_000,
							Subticks:     50_000_000,
							GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 25},
						},
						FillAmount: 100_000_000,
						MakerFee:   -10_000,
						TakerFee:   25_000,
					},
					TotalFilledMaker: 100_000_000,
					TotalFilledTaker: 100_000_000,
				},
			},
			expectedProcessProposerMatchesEvents: types.ProcessProposerMatchesEvents{
				OrderIdsFilledInLastBlock: []types.OrderId{
					{SubaccountId: constants.Bob_Num0, ClientId: 14, ClobPairId: 0},
					{SubaccountId: constants.Alice_Num0, ClientId: 14, ClobPairId: 0},
				},
				BlockHeight: blockHeight,
			},
			// Expected balances are initial balance + balance change due to order - fees
			expectedQuoteBalances: map[satypes.SubaccountId]int64{
				constants.Alice_Num0: constants.Usdc_Asset_100_000.GetBigQuantums().Int64() - 50_000_000 + 10_000,
				constants.Bob_Num0:   constants.Usdc_Asset_100_000.GetBigQuantums().Int64() + 50_000_000 - 25_000,
			},
			expectedPerpetualPositions: map[satypes.SubaccountId][]*satypes.PerpetualPosition{
				constants.Bob_Num0: {
					testutil.CreateSinglePerpetualPosition(
						0,
						big.NewInt(1_000_000_000-100_000_000),
						big.NewInt(0),
						big.NewInt(0),
					),
				},
				constants.Alice_Num0: {
					testutil.CreateSinglePerpetualPosition(
						0,
						big.NewInt(1_000_000_000+100_000_000),
						big.NewInt(0),
						big.NewInt(0),
					),
				},
			},
		},
		"Succeeds with singular match of a preexisting maker and short term taker": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			perpetualFeeParams: &constants.PerpetualFeeParams,
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc,
			},
			subaccounts: []satypes.Subaccount{
				{
					Id: &constants.Alice_Num0,
					AssetPositions: []*satypes.AssetPosition{
						&constants.Usdc_Asset_100_000,
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						testutil.CreateSinglePerpetualPosition(
							0,
							big.NewInt(1_000_000_000), // 10 BTC
							big.NewInt(0),
							big.NewInt(0),
						),
					},
				},
				{
					Id: &constants.Bob_Num0,
					AssetPositions: []*satypes.AssetPosition{
						&constants.Usdc_Asset_100_000,
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						testutil.CreateSinglePerpetualPosition(
							0,
							big.NewInt(1_000_000_000), // 10 BTC
							big.NewInt(0),
							big.NewInt(0),
						),
					},
				},
			},
			preExistingStatefulOrders: []types.Order{
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
			},
			rawOperations: []types.OperationRaw{
				clobtest.NewShortTermOrderPlacementOperationRaw(
					types.Order{
						OrderId:      types.OrderId{SubaccountId: constants.Bob_Num0, ClientId: 14, ClobPairId: 0},
						Side:         types.Order_SIDE_SELL,
						Quantums:     10,
						Subticks:     10,
						GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 25},
					},
				),
				clobtest.NewMatchOperationRaw(
					&types.Order{
						OrderId:      types.OrderId{SubaccountId: constants.Bob_Num0, ClientId: 14, ClobPairId: 0},
						Side:         types.Order_SIDE_SELL,
						Quantums:     10,
						Subticks:     10,
						GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 25},
					},
					[]types.MakerFill{
						{
							FillAmount:   5,
							MakerOrderId: constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.GetOrderId(),
						},
					},
				),
			},
			expectedMatches: []*MatchWithOrdersForTesting{
				{
					MatchWithOrders: types.MatchWithOrders{
						TakerOrder: &types.Order{
							OrderId:      types.OrderId{SubaccountId: constants.Bob_Num0, ClientId: 14, ClobPairId: 0},
							Side:         types.Order_SIDE_SELL,
							Quantums:     10,
							Subticks:     10,
							GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 25},
						},
						MakerOrder: &constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
						FillAmount: 5,
					},
					TotalFilledMaker: 5,
					TotalFilledTaker: 5,
				},
			},
			expectedProcessProposerMatchesEvents: types.ProcessProposerMatchesEvents{
				OrderIdsFilledInLastBlock: []types.OrderId{
					{SubaccountId: constants.Bob_Num0, ClientId: 14, ClobPairId: 0},
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.GetOrderId(),
				},
				RemovedStatefulOrderIds: []types.OrderId{
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.GetOrderId(),
				},
				BlockHeight: blockHeight,
			},
			expectedQuoteBalances: map[satypes.SubaccountId]int64{
				constants.Alice_Num0: constants.Usdc_Asset_100_000.GetBigQuantums().Int64(),
				constants.Bob_Num0:   constants.Usdc_Asset_100_000.GetBigQuantums().Int64(),
			},
			expectedPerpetualPositions: map[satypes.SubaccountId][]*satypes.PerpetualPosition{
				constants.Bob_Num0: {
					testutil.CreateSinglePerpetualPosition(
						0,
						big.NewInt(1_000_000_000-5),
						big.NewInt(0),
						big.NewInt(0),
					),
				},
				constants.Alice_Num0: {
					testutil.CreateSinglePerpetualPosition(
						0,
						big.NewInt(1_000_000_000+5),
						big.NewInt(0),
						big.NewInt(0),
					),
				},
			},
		},
		"Succeeds with singular match of a preexisting maker and newly placed long term taker": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			perpetualFeeParams: &constants.PerpetualFeeParams,
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc,
			},
			subaccounts: []satypes.Subaccount{
				{
					Id: &constants.Alice_Num0,
					AssetPositions: []*satypes.AssetPosition{
						&constants.Usdc_Asset_100_000,
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						testutil.CreateSinglePerpetualPosition(
							0,
							big.NewInt(1_000_000_000), // 10 BTC
							big.NewInt(0),
							big.NewInt(0),
						),
					},
				},
				{
					Id: &constants.Bob_Num0,
					AssetPositions: []*satypes.AssetPosition{
						&constants.Usdc_Asset_100_000,
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						testutil.CreateSinglePerpetualPosition(
							0,
							big.NewInt(1_000_000_000), // 10 BTC
							big.NewInt(0),
							big.NewInt(0),
						),
					},
				},
			},
			preExistingStatefulOrders: []types.Order{
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
				constants.LongTermOrder_Bob_Num0_Id1_Clob0_Sell50_Price10_GTBT15,
			},
			rawOperations: []types.OperationRaw{
				clobtest.NewMatchOperationRaw(
					&constants.LongTermOrder_Bob_Num0_Id1_Clob0_Sell50_Price10_GTBT15,
					[]types.MakerFill{
						{
							FillAmount:   5,
							MakerOrderId: constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.GetOrderId(),
						},
					},
				),
			},
			expectedMatches: []*MatchWithOrdersForTesting{
				{
					MatchWithOrders: types.MatchWithOrders{
						TakerOrder: &constants.LongTermOrder_Bob_Num0_Id1_Clob0_Sell50_Price10_GTBT15,
						MakerOrder: &constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
						FillAmount: 5,
					},
					TotalFilledMaker: 5,
					TotalFilledTaker: 5,
				},
			},

			expectedProcessProposerMatchesEvents: types.ProcessProposerMatchesEvents{
				OrderIdsFilledInLastBlock: []types.OrderId{
					constants.LongTermOrder_Bob_Num0_Id1_Clob0_Sell50_Price10_GTBT15.GetOrderId(),
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.GetOrderId(),
				},
				RemovedStatefulOrderIds: []types.OrderId{
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.GetOrderId(),
				},
				BlockHeight: blockHeight,
			},
			expectedQuoteBalances: map[satypes.SubaccountId]int64{
				constants.Alice_Num0: constants.Usdc_Asset_100_000.GetBigQuantums().Int64(),
				constants.Bob_Num0:   constants.Usdc_Asset_100_000.GetBigQuantums().Int64(),
			},
			expectedPerpetualPositions: map[satypes.SubaccountId][]*satypes.PerpetualPosition{
				constants.Bob_Num0: {
					testutil.CreateSinglePerpetualPosition(
						0,
						big.NewInt(1_000_000_000-5),
						big.NewInt(0),
						big.NewInt(0),
					),
				},
				constants.Alice_Num0: {
					testutil.CreateSinglePerpetualPosition(
						0,
						big.NewInt(1_000_000_000+5),
						big.NewInt(0),
						big.NewInt(0),
					),
				},
			},
		},
		"Succeeds with singular match of a preexisting maker and short term taker with builder code": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			perpetualFeeParams: &constants.PerpetualFeeParams,
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc,
			},
			subaccounts: []satypes.Subaccount{
				{
					Id: &constants.Alice_Num0,
					AssetPositions: []*satypes.AssetPosition{
						&constants.Usdc_Asset_100_000,
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						testutil.CreateSinglePerpetualPosition(
							0,
							big.NewInt(1_000_000_000), // 10 BTC
							big.NewInt(0),
							big.NewInt(0),
						),
					},
				},
				{
					Id: &constants.Bob_Num0,
					AssetPositions: []*satypes.AssetPosition{
						&constants.Usdc_Asset_100_000,
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						testutil.CreateSinglePerpetualPosition(
							0,
							big.NewInt(1_000_000_000), // 10 BTC
							big.NewInt(0),
							big.NewInt(0),
						),
					},
				},
				{
					Id: &constants.Carl_Num0,
					AssetPositions: []*satypes.AssetPosition{
						&constants.Usdc_Asset_0,
					},
					PerpetualPositions: []*satypes.PerpetualPosition{},
				},
			},
			preExistingStatefulOrders: []types.Order{
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT15,
			},
			setupMockBankKeeper: func(bk *mocks.BankKeeper) {
				bk.On(
					"SendCoinsFromModuleToModule",
					mock.Anything,
					mock.Anything,
					mock.Anything,
					mock.Anything,
				).Return(nil)
				bk.On(
					"SendCoins",
					mock.Anything,
					mock.Anything,
					constants.CarlAccAddress,
					mock.MatchedBy(func(coins sdk.Coins) bool {
						// Carl should receive 50_000_000_000 * 10_000ppm (500_000_000) for
						// being the builder for the taker order
						return coins.AmountOf(
							"ibc/8E27BA2D5493AF5636760E354E46004562C46AB7EC0CC4C1CA14E9E20E2545B5",
						).Equal(sdkmath.NewInt(500_000_000))
					}),
				).Return(nil).Once()
				bk.On(
					"SendCoins",
					mock.Anything,
					mock.Anything,
					authtypes.NewModuleAddress(authtypes.FeeCollectorName),
					mock.Anything,
				).Return(nil).Once()
				bk.On(
					"GetBalance",
					mock.Anything,
					mock.Anything,
					mock.Anything,
				).Return(sdk.NewCoin("USDC", sdkmath.NewIntFromUint64(100_000_000_000)))
			},
			rawOperations: []types.OperationRaw{
				clobtest.NewShortTermOrderPlacementOperationRaw(
					types.Order{
						OrderId:      types.OrderId{SubaccountId: constants.Bob_Num0, ClientId: 14, ClobPairId: 0},
						Side:         types.Order_SIDE_SELL,
						Quantums:     100_000_000,    // 1 BTC
						Subticks:     50_000_000_000, // $50,000
						GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 25},
						BuilderCodeParameters: &types.BuilderCodeParameters{
							BuilderAddress: constants.Carl_Num0.Owner,
							FeePpm:         10_000,
						},
					},
				),
				clobtest.NewMatchOperationRaw(
					&types.Order{
						OrderId:      types.OrderId{SubaccountId: constants.Bob_Num0, ClientId: 14, ClobPairId: 0},
						Side:         types.Order_SIDE_SELL,
						Quantums:     100_000_000,
						Subticks:     50_000_000_000,
						GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 25},
						BuilderCodeParameters: &types.BuilderCodeParameters{
							BuilderAddress: constants.Carl_Num0.Owner,
							FeePpm:         10_000,
						},
					},
					[]types.MakerFill{
						{
							FillAmount:   100_000_000,
							MakerOrderId: constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT15.GetOrderId(),
						},
					},
				),
			},
			expectedMatches: []*MatchWithOrdersForTesting{
				{
					MatchWithOrders: types.MatchWithOrders{
						TakerOrder: &types.Order{
							OrderId:      types.OrderId{SubaccountId: constants.Bob_Num0, ClientId: 14, ClobPairId: 0},
							Side:         types.Order_SIDE_SELL,
							Quantums:     100_000_000,    // 1 BTC
							Subticks:     50_000_000_000, // $50,000
							GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 25},
							BuilderCodeParameters: &types.BuilderCodeParameters{
								BuilderAddress: constants.Carl_Num0.Owner,
								FeePpm:         10_000,
							},
						},
						MakerOrder:      &constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT15,
						FillAmount:      100_000_000,
						MakerFee:        10_000_000,
						TakerFee:        25_000_000,
						TakerBuilderFee: 500_000_000,
					},
					TotalFilledMaker: 100_000_000,
					TotalFilledTaker: 100_000_000,
				},
			},
			expectedProcessProposerMatchesEvents: types.ProcessProposerMatchesEvents{
				OrderIdsFilledInLastBlock: []types.OrderId{
					{SubaccountId: constants.Bob_Num0, ClientId: 14, ClobPairId: 0},
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT15.GetOrderId(),
				},
				RemovedStatefulOrderIds: []types.OrderId{
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT15.GetOrderId(),
				},
				BlockHeight: blockHeight,
			},
			expectedQuoteBalances: map[satypes.SubaccountId]int64{
				constants.Alice_Num0: big.NewInt(100_000_000_000 - 50_010_000_000).Int64(),
				constants.Bob_Num0: big.NewInt(100_000_000_000 +
					50_000_000_000 -
					500_000_000 - // builder fee applied
					25_000_000).Int64(),
			},
			expectedPerpetualPositions: map[satypes.SubaccountId][]*satypes.PerpetualPosition{
				constants.Bob_Num0: {
					testutil.CreateSinglePerpetualPosition(
						0,
						big.NewInt(1_000_000_000-100_000_000),
						big.NewInt(0),
						big.NewInt(0),
					),
				},
				constants.Alice_Num0: {
					testutil.CreateSinglePerpetualPosition(
						0,
						big.NewInt(1_000_000_000+100_000_000),
						big.NewInt(0),
						big.NewInt(0),
					),
				},
				constants.Carl_Num0: {},
			},
		},
		"preexisting stateful maker order partially matches with 2 short term taker orders": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			perpetualFeeParams: &constants.PerpetualFeeParams,
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc,
			},
			subaccounts: []satypes.Subaccount{
				{
					Id: &constants.Alice_Num0,
					AssetPositions: []*satypes.AssetPosition{
						&constants.Usdc_Asset_100_000,
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						testutil.CreateSinglePerpetualPosition(
							0,
							big.NewInt(1_000_000_000), // 10 BTC
							big.NewInt(0),
							big.NewInt(0),
						),
					},
				},
				{
					Id: &constants.Bob_Num0,
					AssetPositions: []*satypes.AssetPosition{
						&constants.Usdc_Asset_100_000,
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						testutil.CreateSinglePerpetualPosition(
							0,
							big.NewInt(1_000_000_000), // 10 BTC
							big.NewInt(0),
							big.NewInt(0),
						),
					},
				},
				{
					Id: &constants.Carl_Num0,
					AssetPositions: []*satypes.AssetPosition{
						&constants.Usdc_Asset_100_000,
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						testutil.CreateSinglePerpetualPosition(
							0,
							big.NewInt(1_000_000_000), // 10 BTC
							big.NewInt(0),
							big.NewInt(0),
						),
					},
				},
			},
			preExistingStatefulOrders: []types.Order{
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy100_Price10_GTBT15,
			},
			rawOperations: []types.OperationRaw{
				clobtest.NewShortTermOrderPlacementOperationRaw(
					types.Order{
						OrderId:      types.OrderId{SubaccountId: constants.Bob_Num0, ClientId: 14, ClobPairId: 0},
						Side:         types.Order_SIDE_SELL,
						Quantums:     10,
						Subticks:     10,
						GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 25},
					},
				),
				clobtest.NewMatchOperationRaw(
					&types.Order{
						OrderId:      types.OrderId{SubaccountId: constants.Bob_Num0, ClientId: 14, ClobPairId: 0},
						Side:         types.Order_SIDE_SELL,
						Quantums:     10,
						Subticks:     10,
						GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 25},
					},
					[]types.MakerFill{
						{
							FillAmount:   10,
							MakerOrderId: constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy100_Price10_GTBT15.GetOrderId(),
						},
					},
				),
				clobtest.NewShortTermOrderPlacementOperationRaw(
					types.Order{
						OrderId:      types.OrderId{SubaccountId: constants.Carl_Num0, ClientId: 14, ClobPairId: 0},
						Side:         types.Order_SIDE_SELL,
						Quantums:     15,
						Subticks:     10,
						GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 25},
					},
				),
				clobtest.NewMatchOperationRaw(
					&types.Order{
						OrderId:      types.OrderId{SubaccountId: constants.Carl_Num0, ClientId: 14, ClobPairId: 0},
						Side:         types.Order_SIDE_SELL,
						Quantums:     15,
						Subticks:     10,
						GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 25},
					},
					[]types.MakerFill{
						{
							FillAmount:   15,
							MakerOrderId: constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy100_Price10_GTBT15.GetOrderId(),
						},
					},
				),
			},
			expectedMatches: []*MatchWithOrdersForTesting{
				{
					MatchWithOrders: types.MatchWithOrders{
						TakerOrder: &types.Order{
							OrderId:      types.OrderId{SubaccountId: constants.Bob_Num0, ClientId: 14, ClobPairId: 0},
							Side:         types.Order_SIDE_SELL,
							Quantums:     10,
							Subticks:     10,
							GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 25},
						},
						MakerOrder: &constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy100_Price10_GTBT15,
						FillAmount: 10,
					},
					TotalFilledMaker: 10,
					TotalFilledTaker: 10,
				},
				{
					MatchWithOrders: types.MatchWithOrders{
						TakerOrder: &types.Order{
							OrderId:      types.OrderId{SubaccountId: constants.Carl_Num0, ClientId: 14, ClobPairId: 0},
							Side:         types.Order_SIDE_SELL,
							Quantums:     15,
							Subticks:     10,
							GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 25},
						},

						MakerOrder: &constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy100_Price10_GTBT15,
						FillAmount: 15,
					},
					TotalFilledMaker: 25,
					TotalFilledTaker: 15,
				},
			},

			expectedProcessProposerMatchesEvents: types.ProcessProposerMatchesEvents{
				OrderIdsFilledInLastBlock: []types.OrderId{
					{SubaccountId: constants.Bob_Num0, ClientId: 14, ClobPairId: 0},
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy100_Price10_GTBT15.GetOrderId(),
					{SubaccountId: constants.Carl_Num0, ClientId: 14, ClobPairId: 0},
				},
				BlockHeight: blockHeight,
			},
			expectedQuoteBalances: map[satypes.SubaccountId]int64{
				constants.Alice_Num0: constants.Usdc_Asset_100_000.GetBigQuantums().Int64(),
				constants.Bob_Num0:   constants.Usdc_Asset_100_000.GetBigQuantums().Int64(),
				constants.Carl_Num0:  constants.Usdc_Asset_100_000.GetBigQuantums().Int64(),
			},
			expectedPerpetualPositions: map[satypes.SubaccountId][]*satypes.PerpetualPosition{
				constants.Bob_Num0: {
					testutil.CreateSinglePerpetualPosition(
						0,
						big.NewInt(1_000_000_000-10),
						big.NewInt(0),
						big.NewInt(0),
					),
				},
				constants.Alice_Num0: {
					testutil.CreateSinglePerpetualPosition(
						0,
						big.NewInt(1_000_000_000+10+15),
						big.NewInt(0),
						big.NewInt(0),
					),
				},
				constants.Carl_Num0: {
					testutil.CreateSinglePerpetualPosition(
						0,
						big.NewInt(1_000_000_000-15),
						big.NewInt(0),
						big.NewInt(0),
					),
				},
			},
		},
		// This test matches a liquidation taker order with a short term maker order. The liquidation
		// order is fully filled at one dollar below the bankruptcy price ($49999 vs $50k). Carl's
		// $49,999 is transferred to Dave and Carl's $1 is paid to the insurance fund, leaving him
		// with nothing.
		"Succeeds with liquidation order": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},
			perpetualFeeParams: &constants.PerpetualFeeParams,
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc,
			},
			subaccounts: []satypes.Subaccount{
				// liquidatable: MMR = $5000, TNC = $0
				constants.Carl_Num0_1BTC_Short_50000USD,
				constants.Dave_Num0_1BTC_Long_50000USD,
			},
			rawOperations: []types.OperationRaw{
				clobtest.NewShortTermOrderPlacementOperationRaw(
					constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price49999_GTB10,
				),
				clobtest.NewMatchOperationRawFromPerpetualLiquidation(
					types.MatchPerpetualLiquidation{
						Liquidated:  constants.Carl_Num0,
						ClobPairId:  0,
						PerpetualId: 0,
						TotalSize:   100_000_000,
						IsBuy:       true,
						Fills: []types.MakerFill{
							{
								FillAmount:   100_000_000,
								MakerOrderId: constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price49999_GTB10.GetOrderId(),
							},
						},
					},
				),
			},
			expectedProcessProposerMatchesEvents: types.ProcessProposerMatchesEvents{
				OrderIdsFilledInLastBlock: []types.OrderId{
					constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price49999_GTB10.GetOrderId(),
				},
				BlockHeight: blockHeight,
			},
			expectedMatches: []*MatchWithOrdersForTesting{
				{
					MatchWithOrders: types.MatchWithOrders{
						TakerOrder: &constants.LiquidationOrder_Carl_Num0_Clob0_Buy1BTC_Price50500,
						MakerOrder: &constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price49999_GTB10,
						FillAmount: 100_000_000,
						MakerFee:   9_999_800,
						TakerFee:   1_000_000,
					},
					TotalFilledMaker: 100_000_000,
					TotalFilledTaker: 100_000_000,
				},
			},
			expectedQuoteBalances: map[satypes.SubaccountId]int64{
				constants.Carl_Num0: 0,
				constants.Dave_Num0: constants.Usdc_Asset_99_999.GetBigQuantums().Int64() - int64(9_999_800),
			},
			expectedPerpetualPositions: map[satypes.SubaccountId][]*satypes.PerpetualPosition{
				constants.Carl_Num0: {},
				constants.Dave_Num0: {},
			},
		},
		// This test proposes a set of operations where no liquidation match occurs before the
		// deleveraging match. This happens in the case where the liquidation taker order did
		// not match with any orders on the other side of the book, the subaccount total net collateral
		// is negative.
		// Deleveraging happens at the bankruptcy price ($50,499) so Dave ends up with all of Carl's money.
		"Succeeds with deleveraging with no liquidation order": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},
			perpetualFeeParams: &constants.PerpetualFeeParams,
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc,
			},
			subaccounts: []satypes.Subaccount{
				// liquidatable: MMR = $5000, TNC = $499
				constants.Carl_Num0_1BTC_Short_50499USD,
				constants.Dave_Num0_1BTC_Long_50000USD,
			},
			marketIdToOraclePriceOverride: map[uint32]uint64{
				constants.BtcUsd.MarketId: 5_050_000_000, // $50,500 / BTC
			},
			rawOperations: []types.OperationRaw{
				clobtest.NewMatchOperationRawFromPerpetualDeleveragingLiquidation(
					types.MatchPerpetualDeleveraging{
						Liquidated:  constants.Carl_Num0,
						PerpetualId: 0,
						Fills: []types.MatchPerpetualDeleveraging_Fill{
							{
								OffsettingSubaccountId: constants.Dave_Num0,
								FillAmount:             100_000_000,
							},
						},
					},
				),
			},
			expectedProcessProposerMatchesEvents: types.ProcessProposerMatchesEvents{
				BlockHeight: blockHeight,
			},
			expectedQuoteBalances: map[satypes.SubaccountId]int64{
				constants.Carl_Num0: 0,
				constants.Dave_Num0: constants.Usdc_Asset_100_499.GetBigQuantums().Int64(),
			},
			expectedPerpetualPositions: map[satypes.SubaccountId][]*satypes.PerpetualPosition{
				constants.Carl_Num0: {},
				constants.Dave_Num0: {},
			},
		},
		// This test proposes a set of operations where a liquidation taker order matches with a short
		// term maker order. A deleveraging match is also proposed. For this to happen, the liquidation
		// would have matched with this first order and then tried to match with a second order, resulting
		// in a match that requires insurance funds but the insurance funds are insufficient. When processing
		// the deleveraging operation, the validator will confirm that the subaccount in the deleveraging match
		// has negative TNC, confirming that this is a valid deleveraging match.
		// In this example, the liquidation and deleveraging
		// both happen at bankruptcy price resulting in all of Carl's funds being transferred to Dave.
		"Succeeds with deleveraging and partially filled liquidation": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},
			perpetualFeeParams: &constants.PerpetualFeeParams,
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc,
			},
			subaccounts: []satypes.Subaccount{
				// liquidatable: MMR = $5000.10, TNC = -$1.
				constants.Carl_Num0_1BTC_Short_50000USD,
				constants.Dave_Num0_1BTC_Long_50000USD,
			},
			marketIdToOraclePriceOverride: map[uint32]uint64{
				constants.BtcUsd.MarketId: 5_000_100_000, // $50,001 / BTC
			},
			rawOperations: []types.OperationRaw{
				clobtest.NewShortTermOrderPlacementOperationRaw(
					constants.Order_Dave_Num0_Id1_Clob0_Sell025BTC_Price50000_GTB11,
				),
				clobtest.NewMatchOperationRawFromPerpetualLiquidation(
					types.MatchPerpetualLiquidation{
						Liquidated:  constants.Carl_Num0,
						ClobPairId:  0,
						PerpetualId: 0,
						TotalSize:   100_000_000,
						IsBuy:       true,
						Fills: []types.MakerFill{
							{
								FillAmount:   25_000_000,
								MakerOrderId: constants.Order_Dave_Num0_Id1_Clob0_Sell025BTC_Price50000_GTB11.GetOrderId(),
							},
						},
					},
				),
				clobtest.NewMatchOperationRawFromPerpetualDeleveragingLiquidation(
					types.MatchPerpetualDeleveraging{
						Liquidated:  constants.Carl_Num0,
						PerpetualId: 0,
						Fills: []types.MatchPerpetualDeleveraging_Fill{
							{
								OffsettingSubaccountId: constants.Dave_Num0,
								FillAmount:             75_000_000,
							},
						},
					},
				),
			},
			expectedMatches: []*MatchWithOrdersForTesting{
				{
					MatchWithOrders: types.MatchWithOrders{
						TakerOrder: &constants.LiquidationOrder_Carl_Num0_Clob0_Buy1BTC_Price50501_01,
						MakerOrder: &constants.Order_Dave_Num0_Id1_Clob0_Sell025BTC_Price50000_GTB11,
						FillAmount: 25_000_000,
						MakerFee:   2_500_000,
						TakerFee:   0,
					},
					TotalFilledMaker: 25_000_000,
					TotalFilledTaker: 25_000_000,
				},
			},
			expectedProcessProposerMatchesEvents: types.ProcessProposerMatchesEvents{
				OrderIdsFilledInLastBlock: []types.OrderId{
					constants.Order_Dave_Num0_Id1_Clob0_Sell025BTC_Price50000_GTB11.GetOrderId(),
				},
				BlockHeight: blockHeight,
			},
			expectedQuoteBalances: map[satypes.SubaccountId]int64{
				constants.Carl_Num0: 0,
				constants.Dave_Num0: 100_000_000_000 - 2500000,
			},
			expectedPerpetualPositions: map[satypes.SubaccountId][]*satypes.PerpetualPosition{
				constants.Carl_Num0: {},
				constants.Dave_Num0: {},
			},
		},
		"Zero-fill deleveraging succeeds when the account is negative TNC and updates the last negative TNC subaccount " +
			"seen block number in state": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},
			perpetualFeeParams: &constants.PerpetualFeeParams,
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc,
			},
			subaccounts: []satypes.Subaccount{
				// deleverageable since TNC = -$1
				constants.Carl_Num0_1BTC_Short_50499USD,
			},
			marketIdToOraclePriceOverride: map[uint32]uint64{
				constants.BtcUsd.MarketId: 5_050_000_000, // $50,500 / BTC
			},
			rawOperations: []types.OperationRaw{
				clobtest.NewMatchOperationRawFromPerpetualDeleveragingLiquidation(
					types.MatchPerpetualDeleveraging{
						Liquidated:  constants.Carl_Num0,
						PerpetualId: 0,
						Fills:       []types.MatchPerpetualDeleveraging_Fill{},
					},
				),
			},

			expectedProcessProposerMatchesEvents: types.ProcessProposerMatchesEvents{
				BlockHeight: blockHeight,
			},
			expectedQuoteBalances: map[satypes.SubaccountId]int64{
				constants.Carl_Num0: constants.Carl_Num0_1BTC_Short_50499USD.GetUsdcPosition().Int64(),
			},
			expectedPerpetualPositions: map[satypes.SubaccountId][]*satypes.PerpetualPosition{
				constants.Carl_Num0: constants.Carl_Num0_1BTC_Short_50499USD.GetPerpetualPositions(),
			},
			expectedNegativeTncSubaccountSeen: map[uint32]bool{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance.Params.Id: true,
			},
		},
		"Zero-fill deleveraging succeeds when the account is negative TNC and has a position in final settlement" +
			" market. It updates the last negative TNC subaccount seen block number in state": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			perpetualFeeParams: &constants.PerpetualFeeParams,
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc_Final_Settlement,
			},
			subaccounts: []satypes.Subaccount{
				// liquidatable: MMR = $5000, TNC = -$1.
				constants.Carl_Num0_1BTC_Short_50499USD,
				constants.Dave_Num0_1BTC_Long_50000USD,
			},
			marketIdToOraclePriceOverride: map[uint32]uint64{
				constants.BtcUsd.MarketId: 5_050_000_000, // $50,500 / BTC
			},
			rawOperations: []types.OperationRaw{
				clobtest.NewMatchOperationRawFromPerpetualDeleveragingLiquidation(
					types.MatchPerpetualDeleveraging{
						Liquidated:  constants.Carl_Num0,
						PerpetualId: 0,
						Fills:       []types.MatchPerpetualDeleveraging_Fill{},
					},
				),
			},
			expectedProcessProposerMatchesEvents: types.ProcessProposerMatchesEvents{
				BlockHeight: blockHeight,
			},
			expectedQuoteBalances: map[satypes.SubaccountId]int64{
				constants.Carl_Num0: constants.Carl_Num0_1BTC_Short_50499USD.GetUsdcPosition().Int64(),
				constants.Dave_Num0: constants.Dave_Num0_1BTC_Long_50000USD.GetUsdcPosition().Int64(),
			},
			expectedPerpetualPositions: map[satypes.SubaccountId][]*satypes.PerpetualPosition{
				constants.Carl_Num0: constants.Carl_Num0_1BTC_Short_50499USD.GetPerpetualPositions(),
				constants.Dave_Num0: constants.Dave_Num0_1BTC_Long_50000USD.GetPerpetualPositions(),
			},
			expectedNegativeTncSubaccountSeen: map[uint32]bool{
				constants.BtcUsd_100PercentMarginRequirement.Params.Id: true,
			},
		},
		"Zero-fill deleveraging succeeds when there's multiple zero-fill deleveraging events for the same subaccount " +
			"and perpetual ID": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},
			perpetualFeeParams: &constants.PerpetualFeeParams,
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc,
			},
			subaccounts: []satypes.Subaccount{
				// deleverageable since TNC = -$1
				constants.Carl_Num0_1BTC_Short_50499USD,
			},
			marketIdToOraclePriceOverride: map[uint32]uint64{
				constants.BtcUsd.MarketId: 5_050_000_000, // $50,500 / BTC
			},
			rawOperations: []types.OperationRaw{
				clobtest.NewMatchOperationRawFromPerpetualDeleveragingLiquidation(
					types.MatchPerpetualDeleveraging{
						Liquidated:  constants.Carl_Num0,
						PerpetualId: 0,
						Fills:       []types.MatchPerpetualDeleveraging_Fill{},
					},
				),
				clobtest.NewMatchOperationRawFromPerpetualDeleveragingLiquidation(
					types.MatchPerpetualDeleveraging{
						Liquidated:  constants.Carl_Num0,
						PerpetualId: 0,
						Fills:       []types.MatchPerpetualDeleveraging_Fill{},
					},
				),
			},

			expectedProcessProposerMatchesEvents: types.ProcessProposerMatchesEvents{
				BlockHeight: blockHeight,
			},
			expectedQuoteBalances: map[satypes.SubaccountId]int64{
				constants.Carl_Num0: constants.Carl_Num0_1BTC_Short_50499USD.GetUsdcPosition().Int64(),
			},
			expectedPerpetualPositions: map[satypes.SubaccountId][]*satypes.PerpetualPosition{
				constants.Carl_Num0: constants.Carl_Num0_1BTC_Short_50499USD.GetPerpetualPositions(),
			},
			expectedNegativeTncSubaccountSeen: map[uint32]bool{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance.Params.Id: true,
			},
		},
		"Zero-fill deleverage succeeds after the same subaccount is partially deleveraged": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},
			perpetualFeeParams: &constants.PerpetualFeeParams,
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc,
			},
			subaccounts: []satypes.Subaccount{
				// deleveragable: TNC = -$1
				constants.Carl_Num0_1BTC_Short_50499USD,
				constants.Dave_Num0_1BTC_Long_50000USD,
			},
			marketIdToOraclePriceOverride: map[uint32]uint64{
				constants.BtcUsd.MarketId: 5_050_000_000, // $50,500 / BTC
			},
			rawOperations: []types.OperationRaw{
				clobtest.NewMatchOperationRawFromPerpetualDeleveragingLiquidation(
					types.MatchPerpetualDeleveraging{
						Liquidated:  constants.Carl_Num0,
						PerpetualId: 0,
						Fills: []types.MatchPerpetualDeleveraging_Fill{
							{
								OffsettingSubaccountId: constants.Dave_Num0,
								FillAmount:             50_000_000,
							},
						},
					},
				),
				clobtest.NewMatchOperationRawFromPerpetualDeleveragingLiquidation(
					types.MatchPerpetualDeleveraging{
						Liquidated:  constants.Carl_Num0,
						PerpetualId: 0,
						Fills:       []types.MatchPerpetualDeleveraging_Fill{},
					},
				),
			},
			expectedProcessProposerMatchesEvents: types.ProcessProposerMatchesEvents{
				BlockHeight: blockHeight,
			},
			expectedQuoteBalances: map[satypes.SubaccountId]int64{
				constants.Carl_Num0: constants.Carl_Num0_1BTC_Short_50499USD.GetUsdcPosition().Int64() - 25_249_500_000,
				constants.Dave_Num0: constants.Dave_Num0_1BTC_Long_50000USD.GetUsdcPosition().Int64() + 25_249_500_000,
			},
			expectedPerpetualPositions: map[satypes.SubaccountId][]*satypes.PerpetualPosition{
				constants.Carl_Num0: {
					testutil.CreateSinglePerpetualPosition(
						0,
						big.NewInt(-100_000_000+50_000_000),
						big.NewInt(0),
						big.NewInt(0),
					),
				},
				constants.Dave_Num0: {
					testutil.CreateSinglePerpetualPosition(
						0,
						big.NewInt(100_000_000-50_000_000),
						big.NewInt(0),
						big.NewInt(0),
					),
				},
			},
			expectedNegativeTncSubaccountSeen: map[uint32]bool{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance.Params.Id: true,
			},
		},
		"Zero-fill deleveraging succeeds when the account is negative TNC and updates the last negative TNC subaccount " +
			"seen block number in state for an isolated perpetual collateral pool if the subaccount is isolated to the " +
			"isolated perpetual": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_NoMarginRequirement,
				constants.IsoUsd_IsolatedMarket,
			},
			perpetualFeeParams: &constants.PerpetualFeeParams,
			clobPairs: []types.ClobPair{
				constants.ClobPair_3_Iso,
			},
			subaccounts: []satypes.Subaccount{
				// deleverageable since TNC = -$1
				constants.Carl_Num0_1ISO_Short_49USD,
			},
			marketIdToOraclePriceOverride: map[uint32]uint64{
				constants.IsoUsd_IsolatedMarket.Params.MarketId: 5_000_000_000, // $50 / ISO
			},
			rawOperations: []types.OperationRaw{
				clobtest.NewMatchOperationRawFromPerpetualDeleveragingLiquidation(
					types.MatchPerpetualDeleveraging{
						Liquidated:  constants.Carl_Num0,
						PerpetualId: 3,
						Fills:       []types.MatchPerpetualDeleveraging_Fill{},
					},
				),
			},

			expectedProcessProposerMatchesEvents: types.ProcessProposerMatchesEvents{
				BlockHeight: blockHeight,
			},
			expectedQuoteBalances: map[satypes.SubaccountId]int64{
				constants.Carl_Num0: constants.Carl_Num0_1ISO_Short_49USD.GetUsdcPosition().Int64(),
			},
			expectedPerpetualPositions: map[satypes.SubaccountId][]*satypes.PerpetualPosition{
				constants.Carl_Num0: constants.Carl_Num0_1ISO_Short_49USD.GetPerpetualPositions(),
			},
			expectedNegativeTncSubaccountSeen: map[uint32]bool{
				constants.BtcUsd_NoMarginRequirement.Params.Id: false,
				constants.IsoUsd_IsolatedMarket.Params.Id:      true,
			},
		},
		"Zero-fill deleveraging succeeds when the account is negative TNC and has a position in final settlement" +
			" market. It updates the last negative TNC subaccount seen block number in state for an isolated perpetual" +
			" collateral pool if the subaccount is isolated to the isolated perpetual": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
				constants.IsoUsd_IsolatedMarket,
			},
			perpetualFeeParams: &constants.PerpetualFeeParams,
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc,
				constants.ClobPair_3_Iso_Final_Settlement,
			},
			subaccounts: []satypes.Subaccount{
				// deleveragable: TNC = -$1.
				constants.Carl_Num0_1ISO_Short_49USD,
				constants.Dave_Num0_1BTC_Long_50000USD,
			},
			marketIdToOraclePriceOverride: map[uint32]uint64{
				constants.IsoUsd_IsolatedMarket.Params.MarketId: 5_000_000_000, // $50 / ISO
			},
			rawOperations: []types.OperationRaw{
				clobtest.NewMatchOperationRawFromPerpetualDeleveragingLiquidation(
					types.MatchPerpetualDeleveraging{
						Liquidated:  constants.Carl_Num0,
						PerpetualId: 3,
						Fills:       []types.MatchPerpetualDeleveraging_Fill{},
					},
				),
			},
			expectedProcessProposerMatchesEvents: types.ProcessProposerMatchesEvents{
				BlockHeight: blockHeight,
			},
			expectedQuoteBalances: map[satypes.SubaccountId]int64{
				constants.Carl_Num0: constants.Carl_Num0_1ISO_Short_49USD.GetUsdcPosition().Int64(),
				constants.Dave_Num0: constants.Dave_Num0_1BTC_Long_50000USD.GetUsdcPosition().Int64(),
			},
			expectedPerpetualPositions: map[satypes.SubaccountId][]*satypes.PerpetualPosition{
				constants.Carl_Num0: constants.Carl_Num0_1ISO_Short_49USD.GetPerpetualPositions(),
				constants.Dave_Num0: constants.Dave_Num0_1BTC_Long_50000USD.GetPerpetualPositions(),
			},
			expectedNegativeTncSubaccountSeen: map[uint32]bool{
				constants.BtcUsd_NoMarginRequirement.Params.Id: false,
				constants.IsoUsd_IsolatedMarket.Params.Id:      true,
			},
		},
		"Zero-fill deleveraging succeeds when there's multiple zero-fill deleveraging events for the different subaccount " +
			"and perpetual ID. It updates the last negative TNC subaccount seen block number in state for both isolated " +
			"perpetual collateral pools if the subaccounts are isolated to different isolated perpetuals": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_NoMarginRequirement,
				constants.IsoUsd_IsolatedMarket,
				constants.Iso2Usd_IsolatedMarket,
			},
			perpetualFeeParams: &constants.PerpetualFeeParams,
			clobPairs: []types.ClobPair{
				constants.ClobPair_3_Iso,
				constants.ClobPair_4_Iso2,
			},
			subaccounts: []satypes.Subaccount{
				// deleverageable since TNC = -$1
				constants.Carl_Num0_1ISO_Short_49USD,
				// deleverageable since TNC = -$1
				constants.Dave_Num0_1ISO2_Short_499USD,
			},
			marketIdToOraclePriceOverride: map[uint32]uint64{
				constants.IsoUsd_IsolatedMarket.Params.MarketId:  5_000_000_000, // $50 / ISO
				constants.Iso2Usd_IsolatedMarket.Params.MarketId: 5_000_000_000, // $500 / ISO2
			},
			rawOperations: []types.OperationRaw{
				clobtest.NewMatchOperationRawFromPerpetualDeleveragingLiquidation(
					types.MatchPerpetualDeleveraging{
						Liquidated:  constants.Carl_Num0,
						PerpetualId: 3,
						Fills:       []types.MatchPerpetualDeleveraging_Fill{},
					},
				),
				clobtest.NewMatchOperationRawFromPerpetualDeleveragingLiquidation(
					types.MatchPerpetualDeleveraging{
						Liquidated:  constants.Dave_Num0,
						PerpetualId: 4,
						Fills:       []types.MatchPerpetualDeleveraging_Fill{},
					},
				),
			},

			expectedProcessProposerMatchesEvents: types.ProcessProposerMatchesEvents{
				BlockHeight: blockHeight,
			},
			expectedQuoteBalances: map[satypes.SubaccountId]int64{
				constants.Carl_Num0: constants.Carl_Num0_1ISO_Short_49USD.GetUsdcPosition().Int64(),
				constants.Dave_Num0: constants.Dave_Num0_1ISO2_Short_499USD.GetUsdcPosition().Int64(),
			},
			expectedPerpetualPositions: map[satypes.SubaccountId][]*satypes.PerpetualPosition{
				constants.Carl_Num0: constants.Carl_Num0_1ISO_Short_49USD.GetPerpetualPositions(),
				constants.Dave_Num0: constants.Dave_Num0_1ISO2_Short_499USD.GetPerpetualPositions(),
			},
			expectedNegativeTncSubaccountSeen: map[uint32]bool{
				constants.BtcUsd_NoMarginRequirement.Params.Id: false,
				constants.IsoUsd_IsolatedMarket.Params.Id:      true,
				constants.Iso2Usd_IsolatedMarket.Params.Id:     true,
			},
		},
		"Succeeds order removal operations with previous stateful orders": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			perpetualFeeParams: &constants.PerpetualFeeParams,
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc,
			},
			subaccounts: []satypes.Subaccount{},
			preExistingStatefulOrders: []types.Order{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20,
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
				constants.LongTermOrder_Bob_Num0_Id0_Clob0_Sell10_Price10_GTBT10_PO,
			},
			rawOperations: []types.OperationRaw{
				clobtest.NewOrderRemovalOperationRaw(
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId,
					types.OrderRemoval_REMOVAL_REASON_INVALID_SELF_TRADE,
				),
				clobtest.NewOrderRemovalOperationRaw(
					constants.LongTermOrder_Bob_Num0_Id0_Clob0_Sell10_Price10_GTBT10_PO.OrderId,
					types.OrderRemoval_REMOVAL_REASON_POST_ONLY_WOULD_CROSS_MAKER_ORDER,
				),
			},

			expectedProcessProposerMatchesEvents: types.ProcessProposerMatchesEvents{
				BlockHeight: blockHeight,
				RemovedStatefulOrderIds: []types.OrderId{
					constants.LongTermOrder_Bob_Num0_Id0_Clob0_Sell10_Price10_GTBT10_PO.OrderId,
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId,
				},
			},
		},
		"Fails when attempting to match order with invalid order side": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			perpetualFeeParams: &constants.PerpetualFeeParams,
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc,
			},
			subaccounts: []satypes.Subaccount{
				{
					Id: &constants.Alice_Num0,
					AssetPositions: []*satypes.AssetPosition{
						&constants.Usdc_Asset_100_000,
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						testutil.CreateSinglePerpetualPosition(
							0,
							big.NewInt(1_000_000_000), // 10 BTC
							big.NewInt(0),
							big.NewInt(0),
						),
					},
				},
				{
					Id: &constants.Bob_Num0,
					AssetPositions: []*satypes.AssetPosition{
						&constants.Usdc_Asset_100_000,
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						testutil.CreateSinglePerpetualPosition(
							0,
							big.NewInt(1_000_000_000), // 10 BTC
							big.NewInt(0),
							big.NewInt(0),
						),
					},
				},
			},
			preExistingStatefulOrders: []types.Order{},
			rawOperations: []types.OperationRaw{
				clobtest.NewShortTermOrderPlacementOperationRaw(
					types.Order{
						OrderId:      types.OrderId{SubaccountId: constants.Alice_Num0, ClientId: 14, ClobPairId: 0},
						Side:         types.Order_SIDE_UNSPECIFIED, // Note this side is invalid.
						Quantums:     100_000_000,                  // 1 BTC
						Subticks:     50_000_000,
						GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 25},
					},
				),
				clobtest.NewShortTermOrderPlacementOperationRaw(
					types.Order{
						OrderId:      types.OrderId{SubaccountId: constants.Bob_Num0, ClientId: 14, ClobPairId: 0},
						Side:         types.Order_SIDE_SELL,
						Quantums:     100_000_000, // 1 BTC
						Subticks:     50_000_000,
						GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 25},
					},
				),
				clobtest.NewMatchOperationRaw(
					&types.Order{
						OrderId:      types.OrderId{SubaccountId: constants.Bob_Num0, ClientId: 14, ClobPairId: 0},
						Side:         types.Order_SIDE_SELL,
						Quantums:     100_000_000, // 1 BTC
						Subticks:     50_000_000,
						GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 25},
					},
					[]types.MakerFill{
						{
							FillAmount:   100_000_000,
							MakerOrderId: types.OrderId{SubaccountId: constants.Alice_Num0, ClientId: 14, ClobPairId: 0},
						},
					},
				),
			},

			expectedError: types.ErrInvalidOrderSide,
		},
		// This test proposes an invalid perpetual deleveraging liquidation match operation. The
		// subaccount is not liquidatable, so the match operation should be rejected.
		"Fails with deleveraging match for non-liquidatable subaccount": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},
			perpetualFeeParams: &constants.PerpetualFeeParams,
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_55000USD,
				constants.Dave_Num0_1BTC_Long_50000USD,
			},
			rawOperations: []types.OperationRaw{
				clobtest.NewMatchOperationRawFromPerpetualDeleveragingLiquidation(
					types.MatchPerpetualDeleveraging{
						Liquidated:  constants.Carl_Num0,
						PerpetualId: 0,
						Fills: []types.MatchPerpetualDeleveraging_Fill{
							{
								OffsettingSubaccountId: constants.Dave_Num0,
								FillAmount:             100_000_000,
							},
						},
					},
				),
			},
			expectedQuoteBalances: map[satypes.SubaccountId]int64{
				constants.Carl_Num0: constants.Carl_Num0_1BTC_Short_55000USD.GetUsdcPosition().Int64(),
				constants.Dave_Num0: constants.Usdc_Asset_50_000.GetBigQuantums().Int64(),
			},
			expectedPerpetualPositions: map[satypes.SubaccountId][]*satypes.PerpetualPosition{
				constants.Carl_Num0: constants.Carl_Num0_1BTC_Short_55000USD.GetPerpetualPositions(),
				constants.Dave_Num0: constants.Dave_Num0_1BTC_Long_50000USD.GetPerpetualPositions(),
			},
			expectedError: types.ErrInvalidDeleveragedSubaccount,
		},
		// This test proposes an invalid perpetual deleveraging liquidation match operation. The
		// subaccount has zero TNC, so the deleveraging operation should be rejected.
		"Fails with deleveraging match for subaccount with zero TNC": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},
			perpetualFeeParams: &constants.PerpetualFeeParams,
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_55000USD,
				constants.Dave_Num0_1BTC_Long_50000USD,
			},
			marketIdToOraclePriceOverride: map[uint32]uint64{
				constants.BtcUsd.MarketId: 5_500_000_000, // $55,000 / BTC
			},
			rawOperations: []types.OperationRaw{
				clobtest.NewMatchOperationRawFromPerpetualDeleveragingLiquidation(
					types.MatchPerpetualDeleveraging{
						Liquidated:  constants.Carl_Num0,
						PerpetualId: 0,
						Fills: []types.MatchPerpetualDeleveraging_Fill{
							{
								OffsettingSubaccountId: constants.Dave_Num0,
								FillAmount:             100_000_000,
							},
						},
					},
				),
			},
			expectedQuoteBalances: map[satypes.SubaccountId]int64{
				constants.Carl_Num0: constants.Carl_Num0_1BTC_Short_55000USD.GetUsdcPosition().Int64(),
				constants.Dave_Num0: constants.Usdc_Asset_50_000.GetBigQuantums().Int64(),
			},
			expectedPerpetualPositions: map[satypes.SubaccountId][]*satypes.PerpetualPosition{
				constants.Carl_Num0: constants.Carl_Num0_1BTC_Short_55000USD.GetPerpetualPositions(),
				constants.Dave_Num0: constants.Dave_Num0_1BTC_Long_50000USD.GetPerpetualPositions(),
			},
			expectedError: types.ErrInvalidDeleveragedSubaccount,
		},
		"Conditional: succeeds with singular match of a triggered conditional order": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			perpetualFeeParams: &constants.PerpetualFeeParams,
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc,
			},
			subaccounts: []satypes.Subaccount{
				{
					Id: &constants.Alice_Num0,
					AssetPositions: []*satypes.AssetPosition{
						&constants.Usdc_Asset_100_000,
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						testutil.CreateSinglePerpetualPosition(
							0,
							big.NewInt(1_000_000_000), // 10 BTC
							big.NewInt(0),
							big.NewInt(0),
						),
					},
				},
				{
					Id: &constants.Bob_Num0,
					AssetPositions: []*satypes.AssetPosition{
						&constants.Usdc_Asset_100_000,
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						testutil.CreateSinglePerpetualPosition(
							0,
							big.NewInt(1_000_000_000), // 10 BTC
							big.NewInt(0),
							big.NewInt(0),
						),
					},
				},
			},
			preExistingStatefulOrders: []types.Order{
				constants.LongTermOrder_Bob_Num0_Id1_Clob0_Sell50_Price10_GTBT15,
			},
			triggeredConditionalOrders: []types.Order{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20,
			},
			rawOperations: []types.OperationRaw{
				clobtest.NewMatchOperationRaw(
					&constants.LongTermOrder_Bob_Num0_Id1_Clob0_Sell50_Price10_GTBT15,
					[]types.MakerFill{
						{
							FillAmount:   5,
							MakerOrderId: constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20.GetOrderId(),
						},
					},
				),
			},
			expectedMatches: []*MatchWithOrdersForTesting{
				{
					MatchWithOrders: types.MatchWithOrders{
						TakerOrder: &constants.LongTermOrder_Bob_Num0_Id1_Clob0_Sell50_Price10_GTBT15,
						MakerOrder: &constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20,
						FillAmount: 5,
					},
					TotalFilledMaker: 5,
					TotalFilledTaker: 5,
				},
			},

			expectedProcessProposerMatchesEvents: types.ProcessProposerMatchesEvents{
				OrderIdsFilledInLastBlock: []types.OrderId{
					constants.LongTermOrder_Bob_Num0_Id1_Clob0_Sell50_Price10_GTBT15.GetOrderId(),
					constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20.GetOrderId(),
				},
				RemovedStatefulOrderIds: []types.OrderId{
					constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20.GetOrderId(),
				},
				BlockHeight: blockHeight,
			},
			expectedQuoteBalances: map[satypes.SubaccountId]int64{
				constants.Alice_Num0: constants.Usdc_Asset_100_000.GetBigQuantums().Int64(),
				constants.Bob_Num0:   constants.Usdc_Asset_100_000.GetBigQuantums().Int64(),
			},
			expectedPerpetualPositions: map[satypes.SubaccountId][]*satypes.PerpetualPosition{
				constants.Bob_Num0: {
					testutil.CreateSinglePerpetualPosition(
						0,
						big.NewInt(1_000_000_000-5),
						big.NewInt(0),
						big.NewInt(0),
					),
				},
				constants.Alice_Num0: {
					testutil.CreateSinglePerpetualPosition(
						0,
						big.NewInt(1_000_000_000+5),
						big.NewInt(0),
						big.NewInt(0),
					),
				},
			},
		},
		"Conditional: panics with a non-existent conditional order": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			perpetualFeeParams: &constants.PerpetualFeeParams,
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc,
			},
			subaccounts: []satypes.Subaccount{
				{
					Id: &constants.Alice_Num0,
					AssetPositions: []*satypes.AssetPosition{
						&constants.Usdc_Asset_100_000,
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						testutil.CreateSinglePerpetualPosition(
							0,
							big.NewInt(1_000_000_000), // 10 BTC
							big.NewInt(0),
							big.NewInt(0),
						),
					},
				},
				{
					Id: &constants.Bob_Num0,
					AssetPositions: []*satypes.AssetPosition{
						&constants.Usdc_Asset_100_000,
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						testutil.CreateSinglePerpetualPosition(
							0,
							big.NewInt(1_000_000_000), // 10 BTC
							big.NewInt(0),
							big.NewInt(0),
						),
					},
				},
			},
			preExistingStatefulOrders: []types.Order{
				constants.LongTermOrder_Bob_Num0_Id1_Clob0_Sell50_Price10_GTBT15,
			},
			rawOperations: []types.OperationRaw{
				clobtest.NewMatchOperationRaw(
					&constants.LongTermOrder_Bob_Num0_Id1_Clob0_Sell50_Price10_GTBT15,
					[]types.MakerFill{
						{
							FillAmount:   5,
							MakerOrderId: constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20.GetOrderId(),
						},
					},
				),
			},
			expectedError: errorsmod.Wrapf(
				types.ErrStatefulOrderDoesNotExist,
				"stateful conditional order id %+v does not exist in triggered conditional state.",
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20.OrderId,
			),
		},
		"Conditional: panics with an untriggered conditional order": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			perpetualFeeParams: &constants.PerpetualFeeParams,
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc,
			},
			subaccounts: []satypes.Subaccount{
				{
					Id: &constants.Alice_Num0,
					AssetPositions: []*satypes.AssetPosition{
						&constants.Usdc_Asset_100_000,
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						testutil.CreateSinglePerpetualPosition(
							0,
							big.NewInt(1_000_000_000), // 10 BTC
							big.NewInt(0),
							big.NewInt(0),
						),
					},
				},
				{
					Id: &constants.Bob_Num0,
					AssetPositions: []*satypes.AssetPosition{
						&constants.Usdc_Asset_100_000,
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						testutil.CreateSinglePerpetualPosition(
							0,
							big.NewInt(1_000_000_000), // 10 BTC
							big.NewInt(0),
							big.NewInt(0),
						),
					},
				},
			},
			preExistingStatefulOrders: []types.Order{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20,
				constants.LongTermOrder_Bob_Num0_Id1_Clob0_Sell50_Price10_GTBT15,
			},
			rawOperations: []types.OperationRaw{
				clobtest.NewMatchOperationRaw(
					&constants.LongTermOrder_Bob_Num0_Id1_Clob0_Sell50_Price10_GTBT15,
					[]types.MakerFill{
						{
							FillAmount:   5,
							MakerOrderId: constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20.GetOrderId(),
						},
					},
				),
			},
			expectedError: errorsmod.Wrapf(
				types.ErrStatefulOrderDoesNotExist,
				"stateful conditional order id %+v does not exist in triggered conditional state.",
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20.OrderId,
			),
		},
		"Fails with clob pair not found": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			perpetualFeeParams: &constants.PerpetualFeeParams,
			rawOperations: []types.OperationRaw{
				clobtest.NewMatchOperationRaw(
					&constants.LongTermOrder_Bob_Num0_Id1_Clob0_Sell50_Price10_GTBT15,
					[]types.MakerFill{
						{
							FillAmount:   5,
							MakerOrderId: constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20.GetOrderId(),
						},
					},
				),
			},
			expectedError: types.ErrInvalidClob,
		},
		"Panics with unsupported clob pair status": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			perpetualFeeParams: &constants.PerpetualFeeParams,
			rawOperations: []types.OperationRaw{
				clobtest.NewMatchOperationRaw(
					&constants.LongTermOrder_Bob_Num0_Id1_Clob0_Sell50_Price10_GTBT15,
					[]types.MakerFill{
						{
							FillAmount:   5,
							MakerOrderId: constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20.GetOrderId(),
						},
					},
				),
			},
			// write clob pair to state with unsupported status
			setupState: func(ctx sdk.Context, ks keepertest.ClobKeepersTestContext) {
				cdc := codec.NewProtoCodec(module.InterfaceRegistry)
				store := prefix.NewStore(ks.Ctx.KVStore(ks.StoreKey), []byte(types.ClobPairKeyPrefix))
				b := cdc.MustMarshal(&constants.ClobPair_Btc_Paused)
				store.Set(lib.Uint32ToKey(constants.ClobPair_Btc_Paused.Id), b)
			},
			expectedPanics: "validateInternalOperationAgainstClobPairStatus: ClobPair's status is not supported",
		},
		"Returns error if zero-fill deleveraging operation proposed for non-negative TNC subaccount in final settlement": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			perpetualFeeParams: &constants.PerpetualFeeParams,
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc_Final_Settlement,
			},
			subaccounts: []satypes.Subaccount{
				// both well-collateralized
				constants.Carl_Num0_1BTC_Short_100000USD,
				constants.Dave_Num0_1BTC_Long_50000USD,
			},
			rawOperations: []types.OperationRaw{
				clobtest.NewMatchOperationRawFromPerpetualDeleveragingLiquidation(
					types.MatchPerpetualDeleveraging{
						Liquidated:        constants.Carl_Num0,
						PerpetualId:       0,
						Fills:             []types.MatchPerpetualDeleveraging_Fill{},
						IsFinalSettlement: true,
					},
				),
			},
			expectedProcessProposerMatchesEvents: types.ProcessProposerMatchesEvents{
				BlockHeight: blockHeight,
			},
			expectedQuoteBalances: map[satypes.SubaccountId]int64{
				constants.Carl_Num0: constants.Carl_Num0_1BTC_Short_100000USD.GetUsdcPosition().Int64(),
				constants.Dave_Num0: constants.Dave_Num0_1BTC_Long_50000USD.GetUsdcPosition().Int64(),
			},
			expectedPerpetualPositions: map[satypes.SubaccountId][]*satypes.PerpetualPosition{
				constants.Carl_Num0: constants.Carl_Num0_1BTC_Short_100000USD.GetPerpetualPositions(),
				constants.Dave_Num0: constants.Dave_Num0_1BTC_Long_50000USD.GetPerpetualPositions(),
			},
			expectedError: types.ErrZeroFillDeleveragingForNonNegativeTncSubaccount,
		},
		"Fails with clob match for market in initializing mode": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			perpetualFeeParams: &constants.PerpetualFeeParams,
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc_Initializing,
			},
			rawOperations: []types.OperationRaw{
				clobtest.NewMatchOperationRaw(
					&constants.LongTermOrder_Bob_Num0_Id1_Clob0_Sell50_Price10_GTBT15,
					[]types.MakerFill{
						{
							FillAmount:   5,
							MakerOrderId: constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20.GetOrderId(),
						},
					},
				),
			},
			expectedError: types.ErrOperationConflictsWithClobPairStatus,
		},
		"Fails with short term order placement for market in initializing mode": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			perpetualFeeParams: &constants.PerpetualFeeParams,
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc_Initializing,
			},
			rawOperations: []types.OperationRaw{
				clobtest.NewShortTermOrderPlacementOperationRaw(
					constants.Order_Dave_Num0_Id1_Clob0_Sell025BTC_Price50000_GTB11,
				),
			},
			expectedError: types.ErrOperationConflictsWithClobPairStatus,
		},
		"Fails with order removal for market in initializing mode": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			perpetualFeeParams: &constants.PerpetualFeeParams,
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc_Initializing,
			},
			preExistingStatefulOrders: []types.Order{
				constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10,
			},
			rawOperations: []types.OperationRaw{
				clobtest.NewOrderRemovalOperationRaw(
					constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10.OrderId,
					types.OrderRemoval_REMOVAL_REASON_POST_ONLY_WOULD_CROSS_MAKER_ORDER,
				),
			},
			expectedError: types.ErrOperationConflictsWithClobPairStatus,
		},
		"Fails with order removal reason fully filled": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			perpetualFeeParams: &constants.PerpetualFeeParams,
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc,
			},
			preExistingStatefulOrders: []types.Order{
				constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10,
			},
			rawOperations: []types.OperationRaw{
				clobtest.NewOrderRemovalOperationRaw(
					constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10.OrderId,
					types.OrderRemoval_REMOVAL_REASON_FULLY_FILLED,
				),
			},
			expectedError: types.ErrInvalidOrderRemoval,
		},
		"Fails with order removal for market in final settlement": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			perpetualFeeParams: &constants.PerpetualFeeParams,
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc_Final_Settlement,
			},
			rawOperations: []types.OperationRaw{
				clobtest.NewOrderRemovalOperationRaw(
					constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10.OrderId,
					types.OrderRemoval_REMOVAL_REASON_INVALID_SELF_TRADE,
				),
			},
			expectedError: types.ErrOperationConflictsWithClobPairStatus,
		},
		"Fails with short-term order placement for market in final settlement": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			perpetualFeeParams: &constants.PerpetualFeeParams,
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc_Final_Settlement,
			},
			rawOperations: []types.OperationRaw{
				clobtest.NewShortTermOrderPlacementOperationRaw(
					constants.Order_Dave_Num0_Id1_Clob0_Sell025BTC_Price50000_GTB11,
				),
			},
			expectedError: types.ErrOperationConflictsWithClobPairStatus,
		},
		"Fails with ClobMatch_MatchOrders for market in final settlement": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			perpetualFeeParams: &constants.PerpetualFeeParams,
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc_Final_Settlement,
			},
			rawOperations: []types.OperationRaw{
				clobtest.NewMatchOperationRaw(
					&constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10,
					[]types.MakerFill{
						{
							FillAmount:   10,
							MakerOrderId: constants.LongTermOrder_Alice_Num0_Id1_Clob0_Sell20_Price10_GTBT10.OrderId,
						},
					},
				),
			},
			expectedError: types.ErrOperationConflictsWithClobPairStatus,
		},
		// Liquidations are disallowed for markets in final settlement because they may result
		// in a position increasing in size. This is not allowed for markets in final settlement.
		"Fails with ClobMatch_MatchPerpetualLiquidation for market in final settlement": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			perpetualFeeParams: &constants.PerpetualFeeParams,
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc_Final_Settlement,
			},
			rawOperations: []types.OperationRaw{
				clobtest.NewMatchOperationRawFromPerpetualLiquidation(
					types.MatchPerpetualLiquidation{
						Liquidated:  constants.Alice_Num0,
						ClobPairId:  0,
						PerpetualId: 0,
						TotalSize:   10,
						IsBuy:       false,
						Fills: []types.MakerFill{
							{
								FillAmount:   10,
								MakerOrderId: constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10.OrderId,
							},
						},
					},
				),
			},
			expectedError: types.ErrOperationConflictsWithClobPairStatus,
		},
		// Deleveraging is allowed for markets in final settlement to close out all open positions. A deleveraging
		// event with IsFinalSettlement set to false represents a negative TNC subaccount in the market getting deleveraged.
		"Succeeds with ClobMatch_MatchPerpetualDeleveraging, IsFinalSettlement is false for market in final settlement": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			perpetualFeeParams: &constants.PerpetualFeeParams,
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc_Final_Settlement,
			},
			subaccounts: []satypes.Subaccount{
				// liquidatable: MMR = $5000, TNC = -$1.
				constants.Carl_Num0_1BTC_Short_50499USD,
				constants.Dave_Num0_1BTC_Long_50000USD,
			},
			marketIdToOraclePriceOverride: map[uint32]uint64{
				constants.BtcUsd.MarketId: 5_050_000_000, // $50,500 / BTC
			},
			rawOperations: []types.OperationRaw{
				clobtest.NewMatchOperationRawFromPerpetualDeleveragingLiquidation(
					types.MatchPerpetualDeleveraging{
						Liquidated:  constants.Carl_Num0,
						PerpetualId: 0,
						Fills: []types.MatchPerpetualDeleveraging_Fill{
							{
								OffsettingSubaccountId: constants.Dave_Num0,
								FillAmount:             100_000_000,
							},
						},
					},
				),
			},
			expectedProcessProposerMatchesEvents: types.ProcessProposerMatchesEvents{
				BlockHeight: blockHeight,
			},
			expectedQuoteBalances: map[satypes.SubaccountId]int64{
				constants.Carl_Num0: 0,
				constants.Dave_Num0: constants.Usdc_Asset_100_499.GetBigQuantums().Int64(),
			},
			expectedPerpetualPositions: map[satypes.SubaccountId][]*satypes.PerpetualPosition{
				constants.Carl_Num0: {},
				constants.Dave_Num0: {},
			},
		},
		// Deleveraging is allowed for markets in final settlement to close out all open positions. A deleveraging
		// event with IsFinalSettlement set to true represents a non-negative TNC subaccount having its position closed
		// at the oracle price against other subaccounts with open positions on the opposing side of the book.
		"Succeeds with ClobMatch_MatchPerpetualDeleveraging, IsFinalSettlement is true for market in final settlement": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			perpetualFeeParams: &constants.PerpetualFeeParams,
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc_Final_Settlement,
			},
			subaccounts: []satypes.Subaccount{
				// both well-collateralized
				constants.Carl_Num0_1BTC_Short_100000USD,
				constants.Dave_Num0_1BTC_Long_50000USD,
			},
			rawOperations: []types.OperationRaw{
				clobtest.NewMatchOperationRawFromPerpetualDeleveragingLiquidation(
					types.MatchPerpetualDeleveraging{
						Liquidated:  constants.Carl_Num0,
						PerpetualId: 0,
						Fills: []types.MatchPerpetualDeleveraging_Fill{
							{
								OffsettingSubaccountId: constants.Dave_Num0,
								FillAmount:             100_000_000,
							},
						},
						IsFinalSettlement: true,
					},
				),
			},
			expectedProcessProposerMatchesEvents: types.ProcessProposerMatchesEvents{
				BlockHeight: blockHeight,
			},
			expectedQuoteBalances: map[satypes.SubaccountId]int64{
				constants.Carl_Num0: constants.Usdc_Asset_50_000.GetBigQuantums().Int64(),
				constants.Dave_Num0: constants.Usdc_Asset_100_000.GetBigQuantums().Int64(),
			},
			expectedPerpetualPositions: map[satypes.SubaccountId][]*satypes.PerpetualPosition{
				constants.Carl_Num0: {},
				constants.Dave_Num0: {},
			},
		},
		// This throws an error because the CanDeleverageSubaccount function will return false for
		// shouldFinalSettlePosition, but the IsFinalSettlement flag is set to true.
		`Fails with ClobMatch_MatchPerpetualDeleveraging for negative TNC subaccount,
			IsFinalSettlement is true for market not in final settlement`: {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			perpetualFeeParams: &constants.PerpetualFeeParams,
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_49999USD,
				constants.Dave_Num0_1BTC_Long_50000USD,
			},
			rawOperations: []types.OperationRaw{
				clobtest.NewMatchOperationRawFromPerpetualDeleveragingLiquidation(
					types.MatchPerpetualDeleveraging{
						Liquidated:  constants.Carl_Num0,
						PerpetualId: 0,
						Fills: []types.MatchPerpetualDeleveraging_Fill{
							{
								OffsettingSubaccountId: constants.Dave_Num0,
								FillAmount:             100_000_000,
							},
						},
						IsFinalSettlement: true,
					},
				),
			},
			expectedProcessProposerMatchesEvents: types.ProcessProposerMatchesEvents{
				BlockHeight: blockHeight,
			},
			expectedError: types.ErrDeleveragingIsFinalSettlementFlagMismatch,
		},
		// This test will fail because the CanDeleverageSubaccount function will return false for
		// shouldFinalSettlePosition, but the IsFinalSettlement flag is set to true. Negative TNC subaccounts
		// should never be deleveraged using final settlement (oracle price), and instead should be deleveraged
		// using the bankruptcy price.
		`Fails with ClobMatch_MatchPerpetualDeleveraging for negative TNC subaccount,
			IsFinalSettlement is true for market in final settlement`: {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			perpetualFeeParams: &constants.PerpetualFeeParams,
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc_Final_Settlement,
			},
			subaccounts: []satypes.Subaccount{
				// liquidatable: MMR = $5000, TNC = $499
				constants.Carl_Num0_1BTC_Short_50499USD,
				constants.Dave_Num0_1BTC_Long_50000USD,
			},
			marketIdToOraclePriceOverride: map[uint32]uint64{
				constants.BtcUsd.MarketId: 5_050_000_000, // $50,500 / BTC
			},
			rawOperations: []types.OperationRaw{
				clobtest.NewMatchOperationRawFromPerpetualDeleveragingLiquidation(
					types.MatchPerpetualDeleveraging{
						Liquidated:  constants.Carl_Num0,
						PerpetualId: 0,
						Fills: []types.MatchPerpetualDeleveraging_Fill{
							{
								OffsettingSubaccountId: constants.Dave_Num0,
								FillAmount:             100_000_000,
							},
						},
						IsFinalSettlement: true,
					},
				),
			},
			expectedProcessProposerMatchesEvents: types.ProcessProposerMatchesEvents{
				BlockHeight: blockHeight,
			},
			expectedError: types.ErrDeleveragingIsFinalSettlementFlagMismatch,
		},
		// This test will fail because the CanDeleverageSubaccount function will return false for
		// a non-negative TNC subaccount in a market not in final settlement.
		`Fails with ClobMatch_MatchPerpetualDeleveraging for non-negative TNC subaccount,
			IsFinalSettlement is true for market not in final settlement`: {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			perpetualFeeParams: &constants.PerpetualFeeParams,
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_100000USD,
				constants.Dave_Num0_1BTC_Long_50000USD,
			},
			rawOperations: []types.OperationRaw{
				clobtest.NewMatchOperationRawFromPerpetualDeleveragingLiquidation(
					types.MatchPerpetualDeleveraging{
						Liquidated:  constants.Carl_Num0,
						PerpetualId: 0,
						Fills: []types.MatchPerpetualDeleveraging_Fill{
							{
								OffsettingSubaccountId: constants.Dave_Num0,
								FillAmount:             100_000_000,
							},
						},
						IsFinalSettlement: true,
					},
				),
			},
			expectedProcessProposerMatchesEvents: types.ProcessProposerMatchesEvents{
				BlockHeight: blockHeight,
			},
			expectedError: types.ErrInvalidDeleveragedSubaccount,
		},
		`Fails with ClobMatch_MatchPerpetualDeleveraging for non-negative TNC subaccount,
			IsFinalSettlement is false for market in final settlement`: {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			perpetualFeeParams: &constants.PerpetualFeeParams,
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc_Final_Settlement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_100000USD,
				constants.Dave_Num0_1BTC_Long_50000USD,
			},
			rawOperations: []types.OperationRaw{
				clobtest.NewMatchOperationRawFromPerpetualDeleveragingLiquidation(
					types.MatchPerpetualDeleveraging{
						Liquidated:  constants.Carl_Num0,
						PerpetualId: 0,
						Fills: []types.MatchPerpetualDeleveraging_Fill{
							{
								OffsettingSubaccountId: constants.Dave_Num0,
								FillAmount:             100_000_000,
							},
						},
						IsFinalSettlement: false,
					},
				),
			},
			expectedProcessProposerMatchesEvents: types.ProcessProposerMatchesEvents{
				BlockHeight: blockHeight,
			},
			expectedError: types.ErrDeleveragingIsFinalSettlementFlagMismatch,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			runProcessProposerOperationsTestCase(t, tc)
		})
	}
}

func TestGenerateProcessProposerMatchesEvents(t *testing.T) {
	blockHeight := uint32(5)
	tests := map[string]struct {
		// Params.
		operations []types.InternalOperation

		// Expectations.
		expectedProcessProposerMatchesEvents types.ProcessProposerMatchesEvents
	}{
		"empty operations queue": {
			operations: []types.InternalOperation{},
			expectedProcessProposerMatchesEvents: types.ProcessProposerMatchesEvents{
				ExpiredStatefulOrderIds:                 []types.OrderId{},
				OrderIdsFilledInLastBlock:               []types.OrderId{},
				RemovedStatefulOrderIds:                 []types.OrderId{},
				ConditionalOrderIdsTriggeredInLastBlock: []types.OrderId{},
				BlockHeight:                             blockHeight,
			},
		},
		"short term order matches": {
			operations: []types.InternalOperation{
				types.NewShortTermOrderPlacementInternalOperation(
					constants.Order_Alice_Num0_Id9_Clob1_Buy15_Price45_GTB19,
				),
				types.NewShortTermOrderPlacementInternalOperation(
					constants.Order_Bob_Num0_Id0_Clob1_Sell10_Price15_GTB20,
				),
				types.NewMatchOrdersInternalOperation(
					constants.Order_Bob_Num0_Id0_Clob1_Sell10_Price15_GTB20,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num0_Id9_Clob1_Buy15_Price45_GTB19.OrderId,
							FillAmount:   19,
						},
					},
				),
			},
			expectedProcessProposerMatchesEvents: types.ProcessProposerMatchesEvents{
				ExpiredStatefulOrderIds: []types.OrderId{},
				OrderIdsFilledInLastBlock: []types.OrderId{
					constants.Order_Bob_Num0_Id0_Clob1_Sell10_Price15_GTB20.OrderId,
					constants.Order_Alice_Num0_Id9_Clob1_Buy15_Price45_GTB19.OrderId,
				},
				RemovedStatefulOrderIds:                 []types.OrderId{},
				ConditionalOrderIdsTriggeredInLastBlock: []types.OrderId{},
				BlockHeight:                             blockHeight,
			},
		},
		"liquidation matches": {
			operations: []types.InternalOperation{
				types.NewShortTermOrderPlacementInternalOperation(
					constants.Order_Alice_Num1_Id13_Clob0_Buy50_Price50_GTB30,
				),
				types.NewMatchPerpetualLiquidationInternalOperation(
					&constants.LiquidationOrder_Alice_Num0_Clob0_Sell20_Price25_BTC,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num1_Id13_Clob0_Buy50_Price50_GTB30.OrderId,
							FillAmount:   20,
						},
					},
				),
			},
			expectedProcessProposerMatchesEvents: types.ProcessProposerMatchesEvents{
				ExpiredStatefulOrderIds: []types.OrderId{},
				OrderIdsFilledInLastBlock: []types.OrderId{
					constants.Order_Alice_Num1_Id13_Clob0_Buy50_Price50_GTB30.OrderId,
				},
				RemovedStatefulOrderIds:                 []types.OrderId{},
				ConditionalOrderIdsTriggeredInLastBlock: []types.OrderId{},
				BlockHeight:                             blockHeight,
			},
		},
		"stateful orders in matches": {
			operations: []types.InternalOperation{
				types.NewShortTermOrderPlacementInternalOperation(
					constants.Order_Alice_Num0_Id9_Clob1_Buy15_Price45_GTB19,
				),
				types.NewPreexistingStatefulOrderPlacementInternalOperation(
					constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10,
				),
				types.NewMatchOrdersInternalOperation(
					constants.Order_Alice_Num0_Id9_Clob1_Buy15_Price45_GTB19,
					[]types.MakerFill{
						{
							MakerOrderId: constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10.OrderId,
							FillAmount:   10,
						},
					},
				),
			},
			expectedProcessProposerMatchesEvents: types.ProcessProposerMatchesEvents{
				ExpiredStatefulOrderIds: []types.OrderId{},
				OrderIdsFilledInLastBlock: []types.OrderId{
					constants.Order_Alice_Num0_Id9_Clob1_Buy15_Price45_GTB19.OrderId,
					constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10.OrderId,
				},
				RemovedStatefulOrderIds:                 []types.OrderId{},
				ConditionalOrderIdsTriggeredInLastBlock: []types.OrderId{},
				BlockHeight:                             blockHeight,
			},
		},
		"skips pre existing stateful order operations": {
			operations: []types.InternalOperation{
				types.NewPreexistingStatefulOrderPlacementInternalOperation(
					constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10,
				),
			},
			expectedProcessProposerMatchesEvents: types.ProcessProposerMatchesEvents{
				ExpiredStatefulOrderIds:                 []types.OrderId{},
				OrderIdsFilledInLastBlock:               []types.OrderId{},
				RemovedStatefulOrderIds:                 []types.OrderId{},
				ConditionalOrderIdsTriggeredInLastBlock: []types.OrderId{},
				BlockHeight:                             blockHeight,
			},
		},
		"order removals": {
			operations: []types.InternalOperation{
				types.NewOrderRemovalInternalOperation(
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy100_Price10_GTBT15.OrderId,
					types.OrderRemoval_REMOVAL_REASON_INVALID_SELF_TRADE,
				),
				types.NewOrderRemovalInternalOperation(
					constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10.OrderId,
					types.OrderRemoval_REMOVAL_REASON_POST_ONLY_WOULD_CROSS_MAKER_ORDER,
				),
			},
			expectedProcessProposerMatchesEvents: types.ProcessProposerMatchesEvents{
				ExpiredStatefulOrderIds:   []types.OrderId{},
				OrderIdsFilledInLastBlock: []types.OrderId{},
				RemovedStatefulOrderIds: []types.OrderId{
					constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10.OrderId,
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy100_Price10_GTBT15.OrderId,
				},
				ConditionalOrderIdsTriggeredInLastBlock: []types.OrderId{},
				BlockHeight:                             blockHeight,
			},
		},
	}

	// Run tests.
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			memclob := memclob.NewMemClobPriceTimePriority(true)
			ks := keepertest.NewClobKeepersTestContext(t, memclob, &mocks.BankKeeper{}, &mocks.IndexerEventManager{})
			ctx := ks.Ctx.WithBlockHeight(int64(blockHeight))

			processProposerMatchesEvents := ks.ClobKeeper.GenerateProcessProposerMatchesEvents(ctx, tc.operations)
			require.Equal(t, tc.expectedProcessProposerMatchesEvents, processProposerMatchesEvents)
		})
	}
}

func setupProcessProposerOperationsTestCase(
	t *testing.T,
	tc processProposerOperationsTestCase,
) (
	ctx sdk.Context,
	ks keepertest.ClobKeepersTestContext,
	mockIndexerEventManager *mocks.IndexerEventManager,
) {
	blockHeight := tc.expectedProcessProposerMatchesEvents.BlockHeight

	mockBankKeeper := &mocks.BankKeeper{}
	if tc.setupMockBankKeeper != nil {
		tc.setupMockBankKeeper(mockBankKeeper)
	} else {
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
			mock.Anything,
			mock.Anything,
		).Return(sdk.NewCoin("USDC", sdkmath.NewIntFromUint64(tc.insuranceFundBalance)))
	}

	mockIndexerEventManager = &mocks.IndexerEventManager{}
	// This memclob is not used in the test since DeliverTx creates a new memclob to replay
	// operations on.
	ks = keepertest.NewClobKeepersTestContext(
		t,
		memclob.NewMemClobPriceTimePriority(false),
		mockBankKeeper,
		mockIndexerEventManager,
	)

	// set DeliverTx mode.
	ctx = ks.Ctx.WithIsCheckTx(false)

	// Assert Indexer messages
	if tc.expectedError == nil && tc.expectedPanics == "" && len(tc.expectedMatches) > 0 {
		setupNewMockEventManager(
			t,
			ctx,
			mockIndexerEventManager,
			tc.expectedMatches,
			tc.rawOperations,
		)
	} else {
		mockIndexerEventManager.On("AddTxnEvent",
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return().Maybe()
	}

	// Create the default markets.
	keepertest.CreateTestMarkets(t, ctx, ks.PricesKeeper)

	// Create liquidity tiers.
	keepertest.CreateTestLiquidityTiers(t, ctx, ks.PerpetualsKeeper)

	require.NotNil(t, tc.perpetualFeeParams)
	require.NoError(t, ks.FeeTiersKeeper.SetPerpetualFeeParams(ctx, *tc.perpetualFeeParams))

	err := keepertest.CreateUsdcAsset(ctx, ks.AssetsKeeper)
	require.NoError(t, err)

	// Create all perpetuals.
	for _, p := range tc.perpetuals {
		_, err := ks.PerpetualsKeeper.CreatePerpetual(
			ks.Ctx,
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

	perptest.SetUpDefaultPerpOIsForTest(
		t,
		ks.Ctx,
		ks.PerpetualsKeeper,
		tc.perpetuals,
	)

	// Create all subaccounts.
	for _, subaccount := range tc.subaccounts {
		ks.SubaccountsKeeper.SetSubaccount(ctx, subaccount)
	}

	// Create all CLOBs.
	for i, clobPair := range tc.clobPairs {
		perpetualId := clobtest.MustPerpetualId(clobPair)
		// PerpetualMarketCreateEvents are emitted when initializing the genesis state, so we need to mock
		// the indexer event manager to expect these events.
		if tc.expectedError == nil && tc.expectedPanics == "" && len(tc.expectedMatches) > 0 {
			mockIndexerEventManager.On("AddTxnEvent",
				mock.Anything,
				indexerevents.SubtypePerpetualMarket,
				indexerevents.PerpetualMarketEventVersion,
				indexer_manager.GetBytes(
					indexerevents.NewPerpetualMarketCreateEvent(
						perpetualId,
						uint32(i),
						tc.perpetuals[perpetualId].Params.Ticker,
						tc.perpetuals[perpetualId].Params.MarketId,
						clobPair.Status,
						clobPair.QuantumConversionExponent,
						tc.perpetuals[perpetualId].Params.AtomicResolution,
						clobPair.SubticksPerTick,
						clobPair.StepBaseQuantums,
						tc.perpetuals[perpetualId].Params.LiquidityTier,
						tc.perpetuals[perpetualId].Params.MarketType,
						tc.perpetuals[perpetualId].Params.DefaultFundingPpm,
					),
				),
			).Once().Return()
		}

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

	// Initialize the liquidations config.
	if tc.liquidationConfig != nil {
		require.NoError(t, ks.ClobKeeper.InitializeLiquidationsConfig(ctx, *tc.liquidationConfig))
	} else {
		require.NoError(t, ks.ClobKeeper.InitializeLiquidationsConfig(ctx, constants.LiquidationsConfig_No_Limit))
	}

	// Update the oracle prices.
	for marketId, oraclePrice := range tc.marketIdToOraclePriceOverride {
		err := ks.PricesKeeper.UpdateMarketPrices(
			ks.Ctx,
			[]*pricestypes.MsgUpdateMarketPrices_MarketPrice{
				{
					MarketId: marketId,
					Price:    oraclePrice,
				},
			},
		)
		require.NoError(t, err)
	}

	if tc.setupState != nil {
		tc.setupState(ctx, ks)
	}

	// Create all pre-existing stateful orders in state. Duplicate orders are not allowed.
	seenOrderIds := make(map[types.OrderId]struct{})
	for _, order := range tc.preExistingStatefulOrders {
		_, exists := seenOrderIds[order.GetOrderId()]
		require.Falsef(t, exists, "Duplicate pre-existing stateful order (%+v)", order)
		seenOrderIds[order.GetOrderId()] = struct{}{}
		ks.ClobKeeper.SetLongTermOrderPlacement(ctx, order, blockHeight)
		ks.ClobKeeper.AddStatefulOrderIdExpiration(
			ctx,
			order.MustGetUnixGoodTilBlockTime(),
			order.OrderId,
		)
	}

	for _, order := range tc.triggeredConditionalOrders {
		_, exists := seenOrderIds[order.GetOrderId()]
		require.Falsef(t, exists, "Duplicate pre-existing stateful order (%+v)", order)
		seenOrderIds[order.GetOrderId()] = struct{}{}
		ks.ClobKeeper.SetLongTermOrderPlacement(ctx, order, blockHeight)
		ks.ClobKeeper.AddStatefulOrderIdExpiration(
			ctx,
			order.MustGetUnixGoodTilBlockTime(),
			order.OrderId,
		)

		ks.ClobKeeper.MustTriggerConditionalOrder(ctx, order.OrderId)
	}

	// Set the block time on the context and of the last committed block.
	ctx = ctx.WithBlockTime(time.Unix(5, 0)).WithBlockHeight(int64(blockHeight))
	ks.BlockTimeKeeper.SetPreviousBlockInfo(ctx, &blocktimetypes.BlockInfo{
		Height:    blockHeight,
		Timestamp: time.Unix(int64(5), 0),
	})

	return ctx, ks, mockIndexerEventManager
}

func runProcessProposerOperationsTestCase(
	t *testing.T,
	tc processProposerOperationsTestCase,
) (
	ctx sdk.Context,
	ks keepertest.ClobKeepersTestContext,
) {
	ctx, ks, mockIndexerEventManager := setupProcessProposerOperationsTestCase(t, tc)

	if tc.expectedPanics != "" {
		require.PanicsWithValue(t, tc.expectedPanics, func() {
			_ = ks.ClobKeeper.ProcessProposerOperations(ctx, tc.rawOperations)
		})
		return ctx, ks
	}

	err := ks.ClobKeeper.ProcessProposerOperations(ctx, tc.rawOperations)
	if tc.expectedError != nil {
		require.ErrorContains(t, err, tc.expectedError.Error())
		return ctx, ks
	} else {
		require.NoError(t, err)
	}

	// Verify that processProposerMatchesEvents is the same.
	processProposerMatchesEvents := ks.ClobKeeper.GetProcessProposerMatchesEvents(ctx)
	require.Equal(t, tc.expectedProcessProposerMatchesEvents, processProposerMatchesEvents)

	// Verify that newly-placed stateful orders were written to state.
	for _, newlyPlacedStatefulOrderId := range processProposerMatchesEvents.PlacedLongTermOrderIds {
		_, exists := ks.ClobKeeper.GetLongTermOrderPlacement(ctx, newlyPlacedStatefulOrderId)
		require.Truef(t, exists, "order with ID (%+v) was not placed in state.", newlyPlacedStatefulOrderId)
	}

	// Verify that removed stateful orders were in fact removed from state.
	for _, removedStatefulOrderId := range processProposerMatchesEvents.RemovedStatefulOrderIds {
		_, exists := ks.ClobKeeper.GetLongTermOrderPlacement(ctx, removedStatefulOrderId)
		require.Falsef(t, exists, "order (%+v) was not removed from state.", removedStatefulOrderId)
	}

	// Verify subaccount liquidation info.
	for subaccountId, expected := range tc.expectedSubaccountLiquidationInfo {
		actual := ks.ClobKeeper.GetSubaccountLiquidationInfo(ctx, subaccountId)
		require.Equal(t, expected, actual)
	}

	// Verify subaccount state.
	assertSubaccountState(t, ctx, ks.SubaccountsKeeper, tc.expectedQuoteBalances, tc.expectedPerpetualPositions)

	for orderId, fillAmount := range tc.expectedFillAmounts {
		_, actualFillAmount, _ := ks.ClobKeeper.GetOrderFillAmount(ctx, orderId)
		require.Equal(t, fillAmount, actualFillAmount)
	}

	for perpetualId, expectedNegativeTncSubaccountSeen := range tc.expectedNegativeTncSubaccountSeen {
		// Verify the negative TNC subaccount seen block.
		seenNegativeTncSubaccountBlock, exists, err := ks.SubaccountsKeeper.GetNegativeTncSubaccountSeenAtBlock(
			ctx,
			perpetualId,
		)
		require.NoError(t, err)
		if expectedNegativeTncSubaccountSeen {
			require.True(t, exists)
			require.Equal(t, uint32(ctx.BlockHeight()), seenNegativeTncSubaccountBlock)
		} else {
			require.False(t, exists)
			require.Equal(t, uint32(0), seenNegativeTncSubaccountBlock)
		}
	}

	mockIndexerEventManager.AssertExpectations(t)

	// TODO(CLOB-230) Add more assertions.
	return ctx, ks
}

func setupNewMockEventManager(
	t *testing.T,
	ctx sdk.Context,
	mockIndexerEventManager *mocks.IndexerEventManager,
	matches []*MatchWithOrdersForTesting,
	rawOperations []types.OperationRaw,
) {
	// Add an expectation to the mock for each expected message.
	var matchOrderCallMap = make(map[types.OrderId]*mock.Call)
	for _, match := range matches {
		if match.TakerOrder.IsLiquidation() {
			call := mockIndexerEventManager.On("AddTxnEvent",
				mock.Anything,
				indexerevents.SubtypeOrderFill,
				indexerevents.OrderFillEventVersion,
				indexer_manager.GetBytes(
					indexerevents.NewLiquidationOrderFillEvent(
						match.MakerOrder.MustGetOrder(),
						match.TakerOrder,
						match.FillAmount,
						match.MakerFee,
						match.TakerFee,
						match.MakerBuilderFee,
						match.TotalFilledTaker,
						big.NewInt(0),
						match.MakerOrderRouterFee,
					),
				),
			).Return()

			matchOrderCallMap[match.MakerOrder.MustGetOrder().OrderId] = call
		} else {
			call := mockIndexerEventManager.On("AddTxnEvent",
				mock.Anything,
				indexerevents.SubtypeOrderFill,
				indexerevents.OrderFillEventVersion,
				indexer_manager.GetBytes(
					indexerevents.NewOrderFillEvent(
						match.MakerOrder.MustGetOrder(),
						match.TakerOrder.MustGetOrder(),
						match.FillAmount,
						match.MakerFee,
						match.TakerFee,
						match.MakerBuilderFee,
						match.TakerBuilderFee,
						match.TotalFilledMaker,
						match.TotalFilledTaker,
						big.NewInt(0),
						match.MakerOrderRouterFee,
						match.TakerOrderRouterFee,
					),
				),
			).Return()
			matchOrderCallMap[match.MakerOrder.MustGetOrder().OrderId] = call
			matchOrderCallMap[match.TakerOrder.MustGetOrder().OrderId] = call
		}
	}

	for _, operation := range rawOperations {
		if removal, ok := operation.Operation.(*types.OperationRaw_OrderRemoval); ok {
			mockIndexerEventManager.On("AddTxnEvent",
				mock.Anything,
				indexerevents.SubtypeStatefulOrder,
				indexerevents.StatefulOrderEventVersion,
				indexer_manager.GetBytes(
					indexerevents.NewStatefulOrderRemovalEvent(
						removal.OrderRemoval.OrderId,
						shared.ConvertOrderRemovalReasonToIndexerOrderRemovalReason(
							removal.OrderRemoval.RemovalReason,
						),
					),
				),
			).Once().Return()
		}
		if isClobMatchPerpetualDeleveraging(operation) {
			// Bankruptcy price in DeleveragingEvent is not exposed by API. It is also
			// being tested in other e2e tests. So we don't test it here.
			mockIndexerEventManager.On("AddTxnEvent",
				mock.Anything,
				indexerevents.SubtypeDeleveraging,
				indexerevents.DeleveragingEventVersion,
				mock.Anything,
			).Return()
		}
	}
}

// isClobMatchPerpetualDeleveraging checks if the Operation field is a ClobMatch with a MatchPerpetualDeleveraging.
// It returns true if it is, otherwise false.
func isClobMatchPerpetualDeleveraging(
	operationRaw types.OperationRaw,
) bool {
	matchOperation, ok := operationRaw.Operation.(*types.OperationRaw_Match)
	if !ok {
		return false
	}
	_, ok = matchOperation.Match.Match.(*types.ClobMatch_MatchPerpetualDeleveraging)
	return ok
}

func assertSubaccountState(
	t *testing.T,
	ctx sdk.Context,
	subaccountsKeeper *sakeeper.Keeper,
	expectedQuoteBalances map[satypes.SubaccountId]int64,
	expectedPerpetualPositions map[satypes.SubaccountId][]*satypes.PerpetualPosition,
) {
	for subaccountId, quoteBalance := range expectedQuoteBalances {
		subaccount := subaccountsKeeper.GetSubaccount(ctx, subaccountId)
		require.Equal(t, quoteBalance, subaccount.GetUsdcPosition().Int64())
	}

	for subaccountId, perpetualPositions := range expectedPerpetualPositions {
		subaccount := subaccountsKeeper.GetSubaccount(ctx, subaccountId)
		require.ElementsMatch(t, subaccount.PerpetualPositions, perpetualPositions)
	}
}
