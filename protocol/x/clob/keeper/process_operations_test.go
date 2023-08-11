package keeper_test

import (
	"math/big"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4/dtypes"
	indexerevents "github.com/dydxprotocol/v4/indexer/events"
	"github.com/dydxprotocol/v4/indexer/indexer_manager"
	"github.com/dydxprotocol/v4/mocks"
	clobtest "github.com/dydxprotocol/v4/testutil/clob"
	"github.com/dydxprotocol/v4/testutil/constants"
	keepertest "github.com/dydxprotocol/v4/testutil/keeper"
	"github.com/dydxprotocol/v4/x/clob/memclob"
	"github.com/dydxprotocol/v4/x/clob/types"
	sakeeper "github.com/dydxprotocol/v4/x/subaccounts/keeper"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	perptypes "github.com/dydxprotocol/v4/x/perpetuals/types"
	satypes "github.com/dydxprotocol/v4/x/subaccounts/types"
)

type processProposerOperationsTestCase struct {
	// State
	perpetuals                              []*perptypes.Perpetual
	clobPairs                               []types.ClobPair
	subaccounts                             []satypes.Subaccount
	preExistingStatefulOrders               []types.Order
	operations                              []types.Operation
	addToOrderbookCollatCheckOrderHashesSet map[types.OrderHash]bool

	// Liquidation specific setup.
	liquidationConfig *types.LiquidationsConfig

	// Expectations.
	// Note that for expectedProcessProposerMatchesEvents, the OperationsProposedInLastBlock field is populated from
	// the operations field above.
	expectedProcessProposerMatchesEvents types.ProcessProposerMatchesEvents
	expectedMatches                      []*types.MatchWithOrders
	expectedQuoteBalances                map[satypes.SubaccountId]*big.Int
	expectedPerpetualPositions           map[satypes.SubaccountId][]*satypes.PerpetualPosition
	expectedStatefulOrderPlacements      []types.Order
	expectedStatefulOrderCancelations    []types.OrderId
}

func TestProcessProposerOperations_Success(t *testing.T) {
	blockHeight := uint32(5)
	tests := map[string]processProposerOperationsTestCase{
		"Succeeds no operations": {
			perpetuals:                              []*perptypes.Perpetual{},
			clobPairs:                               []types.ClobPair{},
			subaccounts:                             []satypes.Subaccount{},
			preExistingStatefulOrders:               []types.Order{},
			operations:                              []types.Operation{},
			addToOrderbookCollatCheckOrderHashesSet: map[types.OrderHash]bool{},

			expectedProcessProposerMatchesEvents: types.ProcessProposerMatchesEvents{
				BlockHeight: blockHeight,
			},
		},
		"Succeeds no operations with previous stateful orders": {
			perpetuals: []*perptypes.Perpetual{
				&constants.BtcUsd_100PercentMarginRequirement,
			},
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc,
			},
			subaccounts: []satypes.Subaccount{},
			preExistingStatefulOrders: []types.Order{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
			},
			operations: []types.Operation{},

			expectedProcessProposerMatchesEvents: types.ProcessProposerMatchesEvents{
				BlockHeight: blockHeight,
			},
		},
		"Succeeds with a newly placed long term order that does not match": {
			perpetuals: []*perptypes.Perpetual{
				&constants.BtcUsd_100PercentMarginRequirement,
			},
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc,
			},
			subaccounts: []satypes.Subaccount{
				{
					Id: &constants.Alice_Num0,
					AssetPositions: []*satypes.AssetPosition{
						&constants.Usdc_Asset_100_000,
					},
				},
			},
			preExistingStatefulOrders: []types.Order{},
			operations: []types.Operation{
				types.NewOrderPlacementOperation(
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
				),
			},
			expectedProcessProposerMatchesEvents: types.ProcessProposerMatchesEvents{
				PlacedStatefulOrders: []types.Order{
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
				},
				BlockHeight: blockHeight,
			},
			expectedStatefulOrderPlacements: []types.Order{
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
			},
		},
		"Succeeds with a newly placed long term order that does not match and is then cancelled": {
			perpetuals: []*perptypes.Perpetual{
				&constants.BtcUsd_100PercentMarginRequirement,
			},
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc,
			},
			subaccounts: []satypes.Subaccount{
				{
					Id: &constants.Alice_Num0,
					AssetPositions: []*satypes.AssetPosition{
						&constants.Usdc_Asset_100_000,
					},
				},
			},
			preExistingStatefulOrders: []types.Order{},
			operations: []types.Operation{
				types.NewOrderPlacementOperation(
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
				),
				types.NewOrderCancellationOperation(
					&constants.CancelLongTermOrder_Alice_Num0_Id0_Clob0_GTBT15,
				),
			},
			expectedProcessProposerMatchesEvents: types.ProcessProposerMatchesEvents{
				BlockHeight: blockHeight,
			},
			expectedStatefulOrderCancelations: []types.OrderId{
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId,
			},
			expectedStatefulOrderPlacements: []types.Order{
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
			},
		},
		"Succeeds with singular match of a short term maker and short term taker": {
			perpetuals: []*perptypes.Perpetual{
				&constants.BtcUsd_100PercentMarginRequirement,
			},
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
						{
							PerpetualId: 0,
							Quantums:    dtypes.NewInt(1_000_000_000), // 10 BTC
						},
					},
				},
				{
					Id: &constants.Bob_Num0,
					AssetPositions: []*satypes.AssetPosition{
						&constants.Usdc_Asset_100_000,
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						{
							PerpetualId: 0,
							Quantums:    dtypes.NewInt(1_000_000_000), // 10 BTC
						},
					},
				},
			},
			preExistingStatefulOrders: []types.Order{},
			operations: []types.Operation{
				types.NewOrderPlacementOperation(
					types.Order{
						OrderId:      types.OrderId{SubaccountId: constants.Alice_Num0, ClientId: 14, ClobPairId: 0},
						Side:         types.Order_SIDE_BUY,
						Quantums:     10,
						Subticks:     10,
						GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 25},
					},
				),
				types.NewOrderPlacementOperation(
					types.Order{
						OrderId:      types.OrderId{SubaccountId: constants.Bob_Num0, ClientId: 14, ClobPairId: 0},
						Side:         types.Order_SIDE_SELL,
						Quantums:     10,
						Subticks:     10,
						GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 25},
					},
				),
				types.NewMatchOperation(
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
							MakerOrderId: types.OrderId{SubaccountId: constants.Alice_Num0, ClientId: 14, ClobPairId: 0},
						},
					},
				),
			},
			expectedMatches: []*types.MatchWithOrders{
				{
					TakerOrder: &types.Order{
						OrderId:      types.OrderId{SubaccountId: constants.Bob_Num0, ClientId: 14, ClobPairId: 0},
						Side:         types.Order_SIDE_SELL,
						Quantums:     10,
						Subticks:     10,
						GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 25},
					},
					MakerOrder: &types.Order{
						OrderId:      types.OrderId{SubaccountId: constants.Alice_Num0, ClientId: 14, ClobPairId: 0},
						Side:         types.Order_SIDE_BUY,
						Quantums:     10,
						Subticks:     10,
						GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 25},
					},
					FillAmount: 10,
				},
			},
			expectedProcessProposerMatchesEvents: types.ProcessProposerMatchesEvents{
				OrdersIdsFilledInLastBlock: []types.OrderId{
					{SubaccountId: constants.Bob_Num0, ClientId: 14, ClobPairId: 0},
					{SubaccountId: constants.Alice_Num0, ClientId: 14, ClobPairId: 0},
				},
				BlockHeight: blockHeight,
			},
			expectedQuoteBalances: map[satypes.SubaccountId]*big.Int{
				constants.Alice_Num0: constants.Usdc_Asset_100_000.GetBigQuantums(),
				constants.Bob_Num0:   constants.Usdc_Asset_100_000.GetBigQuantums(),
			},
			expectedPerpetualPositions: map[satypes.SubaccountId][]*satypes.PerpetualPosition{
				constants.Bob_Num0: {
					{
						PerpetualId:  0,
						Quantums:     dtypes.NewInt(1_000_000_000 - 10),
						FundingIndex: dtypes.ZeroInt(),
					},
				},
				constants.Alice_Num0: {
					{
						PerpetualId:  0,
						Quantums:     dtypes.NewInt(1_000_000_000 + 10),
						FundingIndex: dtypes.ZeroInt(),
					},
				},
			},
		},
		"Succeeds with singular match of a preexisting maker and short term taker": {
			perpetuals: []*perptypes.Perpetual{
				&constants.BtcUsd_100PercentMarginRequirement,
			},
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
						{
							PerpetualId: 0,
							Quantums:    dtypes.NewInt(1_000_000_000), // 10 BTC
						},
					},
				},
				{
					Id: &constants.Bob_Num0,
					AssetPositions: []*satypes.AssetPosition{
						&constants.Usdc_Asset_100_000,
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						{
							PerpetualId: 0,
							Quantums:    dtypes.NewInt(1_000_000_000), // 10 BTC
						},
					},
				},
			},
			preExistingStatefulOrders: []types.Order{
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
			},
			operations: []types.Operation{
				types.NewPreexistingStatefulOrderPlacementOperation(
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
				),
				types.NewOrderPlacementOperation(
					types.Order{
						OrderId:      types.OrderId{SubaccountId: constants.Bob_Num0, ClientId: 14, ClobPairId: 0},
						Side:         types.Order_SIDE_SELL,
						Quantums:     10,
						Subticks:     10,
						GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 25},
					},
				),
				types.NewMatchOperation(
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
			expectedMatches: []*types.MatchWithOrders{
				{
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
			},
			expectedProcessProposerMatchesEvents: types.ProcessProposerMatchesEvents{
				OrdersIdsFilledInLastBlock: []types.OrderId{
					{SubaccountId: constants.Bob_Num0, ClientId: 14, ClobPairId: 0},
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.GetOrderId(),
				},
				BlockHeight: blockHeight,
			},
			expectedQuoteBalances: map[satypes.SubaccountId]*big.Int{
				constants.Alice_Num0: constants.Usdc_Asset_100_000.GetBigQuantums(),
				constants.Bob_Num0:   constants.Usdc_Asset_100_000.GetBigQuantums(),
			},
			expectedPerpetualPositions: map[satypes.SubaccountId][]*satypes.PerpetualPosition{
				constants.Bob_Num0: {
					{
						PerpetualId:  0,
						Quantums:     dtypes.NewInt(1_000_000_000 - 5),
						FundingIndex: dtypes.ZeroInt(),
					},
				},
				constants.Alice_Num0: {
					{
						PerpetualId:  0,
						Quantums:     dtypes.NewInt(1_000_000_000 + 5),
						FundingIndex: dtypes.ZeroInt(),
					},
				},
			},
		},
		"Succeeds with singular match of a preexisting maker and newly placed long term taker": {
			perpetuals: []*perptypes.Perpetual{
				&constants.BtcUsd_100PercentMarginRequirement,
			},
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
						{
							PerpetualId: 0,
							Quantums:    dtypes.NewInt(1_000_000_000), // 10 BTC
						},
					},
				},
				{
					Id: &constants.Bob_Num0,
					AssetPositions: []*satypes.AssetPosition{
						&constants.Usdc_Asset_100_000,
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						{
							PerpetualId: 0,
							Quantums:    dtypes.NewInt(1_000_000_000), // 10 BTC
						},
					},
				},
			},
			preExistingStatefulOrders: []types.Order{
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
			},
			operations: []types.Operation{
				types.NewPreexistingStatefulOrderPlacementOperation(
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
				),
				types.NewOrderPlacementOperation(
					constants.LongTermOrder_User2_Num0_Id1_Clob0_Sell50_Price10_GTBT15,
				),
				types.NewMatchOperation(
					&constants.LongTermOrder_User2_Num0_Id1_Clob0_Sell50_Price10_GTBT15,
					[]types.MakerFill{
						{
							FillAmount:   5,
							MakerOrderId: constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.GetOrderId(),
						},
					},
				),
			},
			expectedMatches: []*types.MatchWithOrders{
				{
					TakerOrder: &constants.LongTermOrder_User2_Num0_Id1_Clob0_Sell50_Price10_GTBT15,
					MakerOrder: &constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
					FillAmount: 5,
				},
			},

			expectedProcessProposerMatchesEvents: types.ProcessProposerMatchesEvents{
				PlacedStatefulOrders: []types.Order{
					constants.LongTermOrder_User2_Num0_Id1_Clob0_Sell50_Price10_GTBT15,
				},
				OrdersIdsFilledInLastBlock: []types.OrderId{
					constants.LongTermOrder_User2_Num0_Id1_Clob0_Sell50_Price10_GTBT15.GetOrderId(),
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.GetOrderId(),
				},
				BlockHeight: blockHeight,
			},
			expectedQuoteBalances: map[satypes.SubaccountId]*big.Int{
				constants.Alice_Num0: constants.Usdc_Asset_100_000.GetBigQuantums(),
				constants.Bob_Num0:   constants.Usdc_Asset_100_000.GetBigQuantums(),
			},
			expectedPerpetualPositions: map[satypes.SubaccountId][]*satypes.PerpetualPosition{
				constants.Bob_Num0: {
					{
						PerpetualId:  0,
						Quantums:     dtypes.NewInt(1_000_000_000 - 5),
						FundingIndex: dtypes.ZeroInt(),
					},
				},
				constants.Alice_Num0: {
					{
						PerpetualId:  0,
						Quantums:     dtypes.NewInt(1_000_000_000 + 5),
						FundingIndex: dtypes.ZeroInt(),
					},
				},
			},
			expectedStatefulOrderPlacements: []types.Order{
				constants.LongTermOrder_User2_Num0_Id1_Clob0_Sell50_Price10_GTBT15,
			},
		},
		"preexisting stateful maker order partially matches with 2 short term taker orders": {
			perpetuals: []*perptypes.Perpetual{
				&constants.BtcUsd_100PercentMarginRequirement,
			},
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
						{
							PerpetualId: 0,
							Quantums:    dtypes.NewInt(1_000_000_000), // 10 BTC
						},
					},
				},
				{
					Id: &constants.Bob_Num0,
					AssetPositions: []*satypes.AssetPosition{
						&constants.Usdc_Asset_100_000,
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						{
							PerpetualId: 0,
							Quantums:    dtypes.NewInt(1_000_000_000), // 10 BTC
						},
					},
				},
				{
					Id: &constants.Carl_Num0,
					AssetPositions: []*satypes.AssetPosition{
						&constants.Usdc_Asset_100_000,
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						{
							PerpetualId: 0,
							Quantums:    dtypes.NewInt(1_000_000_000), // 10 BTC
						},
					},
				},
			},
			preExistingStatefulOrders: []types.Order{
				constants.LongTermOrder_User1_Num0_Id0_Clob0_Buy100_Price10_GTBT15,
			},
			operations: []types.Operation{
				types.NewPreexistingStatefulOrderPlacementOperation(
					constants.LongTermOrder_User1_Num0_Id0_Clob0_Buy100_Price10_GTBT15,
				),
				types.NewOrderPlacementOperation(
					types.Order{
						OrderId:      types.OrderId{SubaccountId: constants.Bob_Num0, ClientId: 14, ClobPairId: 0},
						Side:         types.Order_SIDE_SELL,
						Quantums:     10,
						Subticks:     10,
						GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 25},
					},
				),
				types.NewMatchOperation(
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
							MakerOrderId: constants.LongTermOrder_User1_Num0_Id0_Clob0_Buy100_Price10_GTBT15.GetOrderId(),
						},
					},
				),
				types.NewOrderPlacementOperation(
					types.Order{
						OrderId:      types.OrderId{SubaccountId: constants.Carl_Num0, ClientId: 14, ClobPairId: 0},
						Side:         types.Order_SIDE_SELL,
						Quantums:     15,
						Subticks:     10,
						GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 25},
					},
				),
				types.NewMatchOperation(
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
							MakerOrderId: constants.LongTermOrder_User1_Num0_Id0_Clob0_Buy100_Price10_GTBT15.GetOrderId(),
						},
					},
				),
			},
			expectedMatches: []*types.MatchWithOrders{
				{
					TakerOrder: &types.Order{
						OrderId:      types.OrderId{SubaccountId: constants.Bob_Num0, ClientId: 14, ClobPairId: 0},
						Side:         types.Order_SIDE_SELL,
						Quantums:     10,
						Subticks:     10,
						GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 25},
					},
					MakerOrder: &constants.LongTermOrder_User1_Num0_Id0_Clob0_Buy100_Price10_GTBT15,
					FillAmount: 10,
				},
				{
					TakerOrder: &types.Order{
						OrderId:      types.OrderId{SubaccountId: constants.Carl_Num0, ClientId: 14, ClobPairId: 0},
						Side:         types.Order_SIDE_SELL,
						Quantums:     15,
						Subticks:     10,
						GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 25},
					},

					MakerOrder: &constants.LongTermOrder_User1_Num0_Id0_Clob0_Buy100_Price10_GTBT15,
					FillAmount: 15,
				},
			},

			expectedProcessProposerMatchesEvents: types.ProcessProposerMatchesEvents{
				OrdersIdsFilledInLastBlock: []types.OrderId{
					{SubaccountId: constants.Bob_Num0, ClientId: 14, ClobPairId: 0},
					constants.LongTermOrder_User1_Num0_Id0_Clob0_Buy100_Price10_GTBT15.GetOrderId(),
					{SubaccountId: constants.Carl_Num0, ClientId: 14, ClobPairId: 0},
				},
				BlockHeight: blockHeight,
			},
			expectedQuoteBalances: map[satypes.SubaccountId]*big.Int{
				constants.Alice_Num0: constants.Usdc_Asset_100_000.GetBigQuantums(),
				constants.Bob_Num0:   constants.Usdc_Asset_100_000.GetBigQuantums(),
				constants.Carl_Num0:  constants.Usdc_Asset_100_000.GetBigQuantums(),
			},
			expectedPerpetualPositions: map[satypes.SubaccountId][]*satypes.PerpetualPosition{
				constants.Bob_Num0: {
					{
						PerpetualId:  0,
						Quantums:     dtypes.NewInt(1_000_000_000 - 10),
						FundingIndex: dtypes.ZeroInt(),
					},
				},
				constants.Alice_Num0: {
					{
						PerpetualId:  0,
						Quantums:     dtypes.NewInt(1_000_000_000 + 10 + 15),
						FundingIndex: dtypes.ZeroInt(),
					},
				},
				constants.Carl_Num0: {
					{
						PerpetualId:  0,
						Quantums:     dtypes.NewInt(1_000_000_000 - 15),
						FundingIndex: dtypes.ZeroInt(),
					},
				},
			},
		},
		"preexisting stateful maker order partially matches with 2 short term taker orders and is then cancelled": {
			perpetuals: []*perptypes.Perpetual{
				&constants.BtcUsd_100PercentMarginRequirement,
			},
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
						{
							PerpetualId: 0,
							Quantums:    dtypes.NewInt(1_000_000_000), // 10 BTC
						},
					},
				},
				{
					Id: &constants.Bob_Num0,
					AssetPositions: []*satypes.AssetPosition{
						&constants.Usdc_Asset_100_000,
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						{
							PerpetualId: 0,
							Quantums:    dtypes.NewInt(1_000_000_000), // 10 BTC
						},
					},
				},
				{
					Id: &constants.Carl_Num0,
					AssetPositions: []*satypes.AssetPosition{
						&constants.Usdc_Asset_100_000,
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						{
							PerpetualId: 0,
							Quantums:    dtypes.NewInt(1_000_000_000), // 10 BTC
						},
					},
				},
			},
			preExistingStatefulOrders: []types.Order{
				constants.LongTermOrder_User1_Num0_Id0_Clob0_Buy100_Price10_GTBT15,
			},
			operations: []types.Operation{
				types.NewPreexistingStatefulOrderPlacementOperation(
					constants.LongTermOrder_User1_Num0_Id0_Clob0_Buy100_Price10_GTBT15,
				),
				types.NewOrderPlacementOperation(
					types.Order{
						OrderId:      types.OrderId{SubaccountId: constants.Bob_Num0, ClientId: 14, ClobPairId: 0},
						Side:         types.Order_SIDE_SELL,
						Quantums:     10,
						Subticks:     10,
						GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 25},
					},
				),
				types.NewMatchOperation(
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
							MakerOrderId: constants.LongTermOrder_User1_Num0_Id0_Clob0_Buy100_Price10_GTBT15.GetOrderId(),
						},
					},
				),
				types.NewOrderPlacementOperation(
					types.Order{
						OrderId:      types.OrderId{SubaccountId: constants.Carl_Num0, ClientId: 14, ClobPairId: 0},
						Side:         types.Order_SIDE_SELL,
						Quantums:     15,
						Subticks:     10,
						GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 25},
					},
				),
				types.NewMatchOperation(
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
							MakerOrderId: constants.LongTermOrder_User1_Num0_Id0_Clob0_Buy100_Price10_GTBT15.GetOrderId(),
						},
					},
				),
				types.NewOrderCancellationOperation(
					&constants.CancelLongTermOrder_Alice_Num0_Id0_Clob0_GTBT15,
				),
			},
			expectedStatefulOrderCancelations: []types.OrderId{
				constants.CancelLongTermOrder_Alice_Num0_Id0_Clob0_GTBT15.OrderId,
			},
			expectedMatches: []*types.MatchWithOrders{
				{
					TakerOrder: &types.Order{
						OrderId:      types.OrderId{SubaccountId: constants.Bob_Num0, ClientId: 14, ClobPairId: 0},
						Side:         types.Order_SIDE_SELL,
						Quantums:     10,
						Subticks:     10,
						GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 25},
					},
					MakerOrder: &constants.LongTermOrder_User1_Num0_Id0_Clob0_Buy100_Price10_GTBT15,
					FillAmount: 10,
				},
				{
					TakerOrder: &types.Order{
						OrderId:      types.OrderId{SubaccountId: constants.Carl_Num0, ClientId: 14, ClobPairId: 0},
						Side:         types.Order_SIDE_SELL,
						Quantums:     15,
						Subticks:     10,
						GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 25},
					},

					MakerOrder: &constants.LongTermOrder_User1_Num0_Id0_Clob0_Buy100_Price10_GTBT15,
					FillAmount: 15,
				},
			},

			expectedProcessProposerMatchesEvents: types.ProcessProposerMatchesEvents{
				OrdersIdsFilledInLastBlock: []types.OrderId{
					{SubaccountId: constants.Bob_Num0, ClientId: 14, ClobPairId: 0},
					constants.LongTermOrder_User1_Num0_Id0_Clob0_Buy100_Price10_GTBT15.GetOrderId(),
					{SubaccountId: constants.Carl_Num0, ClientId: 14, ClobPairId: 0},
				},
				BlockHeight: blockHeight,
			},
			expectedQuoteBalances: map[satypes.SubaccountId]*big.Int{
				constants.Alice_Num0: constants.Usdc_Asset_100_000.GetBigQuantums(),
				constants.Bob_Num0:   constants.Usdc_Asset_100_000.GetBigQuantums(),
				constants.Carl_Num0:  constants.Usdc_Asset_100_000.GetBigQuantums(),
			},
			expectedPerpetualPositions: map[satypes.SubaccountId][]*satypes.PerpetualPosition{
				constants.Alice_Num0: {
					{
						PerpetualId:  0,
						Quantums:     dtypes.NewInt(1_000_000_000 + 10 + 15),
						FundingIndex: dtypes.ZeroInt(),
					},
				},
				constants.Carl_Num0: {
					{
						PerpetualId:  0,
						Quantums:     dtypes.NewInt(1_000_000_000 - 15),
						FundingIndex: dtypes.ZeroInt(),
					},
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mockBankKeeper := &mocks.BankKeeper{}
			mockBankKeeper.On(
				"SendCoinsFromModuleToModule",
				mock.Anything,
				mock.Anything,
				mock.Anything,
				mock.Anything,
			).Return(nil)

			mockIndexerEventManager := &mocks.IndexerEventManager{}
			// This memclob is not used in the test since DeliverTx creates a new memclob to replay
			// operations on.
			ctx,
				clobKeeper,
				pricesKeeper,
				assetsKeeper,
				perpetualsKeeper,
				subaccountsKeeper,
				_,
				_ := keepertest.ClobKeepers(
				t,
				memclob.NewMemClobPriceTimePriority(false),
				mockBankKeeper,
				mockIndexerEventManager,
			)

			// set DeliverTx mode.
			ctx = ctx.WithIsCheckTx(false)

			// Assert Indexer messages
			setupNewMockEventManager(
				t,
				ctx,
				mockIndexerEventManager,
				tc.expectedMatches,
				tc.expectedStatefulOrderPlacements,
				tc.expectedStatefulOrderCancelations,
			)

			// Create the default markets.
			keepertest.CreateTestMarketsAndExchangeFeeds(t, ctx, pricesKeeper)

			// Create liquidity tiers.
			keepertest.CreateTestLiquidityTiers(t, ctx, perpetualsKeeper)

			err := keepertest.CreateUsdcAsset(ctx, assetsKeeper)
			require.NoError(t, err)

			// Create all perpetuals.
			for _, p := range tc.perpetuals {
				_, err := perpetualsKeeper.CreatePerpetual(
					ctx,
					p.Ticker,
					p.MarketId,
					p.AtomicResolution,
					p.DefaultFundingPpm,
					p.LiquidityTier,
				)
				require.NoError(t, err)
			}

			// Create all subaccounts.
			for _, subaccount := range tc.subaccounts {
				subaccountsKeeper.SetSubaccount(ctx, subaccount)
			}

			// Create all CLOBs.
			for _, clobPair := range tc.clobPairs {
				_, err = clobKeeper.CreatePerpetualClobPair(
					ctx,
					clobtest.MustPerpetualId(clobPair),
					satypes.BaseQuantums(clobPair.StepBaseQuantums),
					satypes.BaseQuantums(clobPair.MinOrderBaseQuantums),
					clobPair.QuantumConversionExponent,
					clobPair.SubticksPerTick,
					clobPair.Status,
					clobPair.MakerFeePpm,
					clobPair.TakerFeePpm,
				)
				require.NoError(t, err)
			}

			// Initialize the liquidations config.
			if tc.liquidationConfig != nil {
				require.NoError(t, clobKeeper.InitializeLiquidationsConfig(ctx, *tc.liquidationConfig))
			} else {
				require.NoError(t, clobKeeper.InitializeLiquidationsConfig(ctx, constants.LiquidationsConfig_No_Limit))
			}

			// Create all pre-existing stateful orders in state. Duplicate orders are not allowed.
			// We don't need to set the stateful order placement in memclob because the deliverTx flow
			// will create its own memclob.
			seenOrderIds := make(map[types.OrderId]struct{})
			for _, order := range tc.preExistingStatefulOrders {
				_, exists := seenOrderIds[order.GetOrderId()]
				require.Falsef(t, exists, "Duplicate pre-existing stateful order (+%v)", order)
				seenOrderIds[order.GetOrderId()] = struct{}{}
				clobKeeper.SetStatefulOrderPlacement(ctx, order, blockHeight)
				clobKeeper.MustAddOrderToStatefulOrdersTimeSlice(
					ctx,
					order.MustGetUnixGoodTilBlockTime(),
					order.OrderId,
				)
			}

			// Set the block time on the context and of the last committed block.
			ctx = ctx.WithBlockTime(time.Unix(5, 0)).WithBlockHeight(int64(blockHeight))
			clobKeeper.SetBlockTimeForLastCommittedBlock(ctx)

			// Run the DeliverTx ProcessProposerOperations flow.
			err = clobKeeper.ProcessProposerOperations(ctx, tc.operations, tc.addToOrderbookCollatCheckOrderHashesSet)
			require.NoError(t, err)

			// Verify that processProposerMatchesEvents is the same.
			processProposerMatchesEvents := clobKeeper.GetProcessProposerMatchesEvents(ctx)
			// Events operations proposed should directly match operations proposed.
			tc.expectedProcessProposerMatchesEvents.OperationsProposedInLastBlock = tc.operations
			// Workaround for the empty slice vs nil case. In this case, verify length assertion.
			if len(processProposerMatchesEvents.OperationsProposedInLastBlock) == 0 {
				require.Len(
					t,
					tc.expectedProcessProposerMatchesEvents.OperationsProposedInLastBlock,
					0,
				)
				tc.expectedProcessProposerMatchesEvents.OperationsProposedInLastBlock =
					processProposerMatchesEvents.OperationsProposedInLastBlock
			}
			require.Equal(t, tc.expectedProcessProposerMatchesEvents, processProposerMatchesEvents)

			// Verify that newly-placed stateful orders were written to state.
			for _, newlyPlacedStatefulOrder := range processProposerMatchesEvents.PlacedStatefulOrders {
				exists := clobKeeper.DoesStatefulOrderExistInState(ctx, newlyPlacedStatefulOrder)
				require.Truef(t, exists, "order (%+v) was not placed in state.", newlyPlacedStatefulOrder)
			}

			// Verify subaccount state.
			assertSubaccountState(t, ctx, subaccountsKeeper, tc.expectedQuoteBalances, tc.expectedPerpetualPositions)

			mockIndexerEventManager.AssertExpectations(t)

			// TODO(CLOB-230) Add more assertions.
		})
	}
}

func TestGenerateProcessProposerMatchesEvents(t *testing.T) {
	blockHeight := uint32(5)
	tests := map[string]struct {
		// Params.
		operations []types.Operation

		// Expectations.
		expectedProcessProposerMatchesEvents types.ProcessProposerMatchesEvents
	}{
		"empty operations queue": {
			operations: []types.Operation{},
			expectedProcessProposerMatchesEvents: types.ProcessProposerMatchesEvents{
				PlacedStatefulOrders:       []types.Order{},
				ExpiredStatefulOrderIds:    []types.OrderId{},
				OrdersIdsFilledInLastBlock: []types.OrderId{},
				BlockHeight:                blockHeight,
			},
		},
		"does not include short-term orders": {
			operations: []types.Operation{
				types.NewOrderPlacementOperation(
					constants.Order_Alice_Num0_Id9_Clob1_Buy15_Price45_GTB19,
				),
				types.NewOrderPlacementOperation(
					constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10,
				),
				types.NewMatchOperation(
					&constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num0_Id9_Clob1_Buy15_Price45_GTB19.OrderId,
							FillAmount:   10,
						},
					},
				),
			},
			expectedProcessProposerMatchesEvents: types.ProcessProposerMatchesEvents{
				PlacedStatefulOrders: []types.Order{
					constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10,
				},
				ExpiredStatefulOrderIds: []types.OrderId{},
				OrdersIdsFilledInLastBlock: []types.OrderId{
					constants.Order_Alice_Num0_Id9_Clob1_Buy15_Price45_GTB19.OrderId,
					constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10.OrderId,
				},
				BlockHeight: blockHeight,
			},
		},
		"correctly handles order matches": {
			operations: []types.Operation{
				types.NewOrderPlacementOperation(
					constants.Order_Alice_Num0_Id9_Clob1_Buy15_Price45_GTB19,
				),
				types.NewOrderPlacementOperation(
					constants.Order_Bob_Num0_Id0_Clob1_Sell10_Price15_GTB20,
				),
				types.NewMatchOperation(
					&constants.Order_Bob_Num0_Id0_Clob1_Sell10_Price15_GTB20,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num0_Id9_Clob1_Buy15_Price45_GTB19.OrderId,
							FillAmount:   10,
						},
					},
				),
			},
			expectedProcessProposerMatchesEvents: types.ProcessProposerMatchesEvents{
				PlacedStatefulOrders:    []types.Order{},
				ExpiredStatefulOrderIds: []types.OrderId{},
				OrdersIdsFilledInLastBlock: []types.OrderId{
					constants.Order_Bob_Num0_Id0_Clob1_Sell10_Price15_GTB20.OrderId,
					constants.Order_Alice_Num0_Id9_Clob1_Buy15_Price45_GTB19.OrderId,
				},
				BlockHeight: blockHeight,
			},
		},
		"correctly handles liquidation matches": {
			operations: []types.Operation{
				types.NewOrderPlacementOperation(
					constants.Order_Alice_Num1_Id13_Clob0_Buy50_Price50_GTB30,
				),
				types.NewMatchOperation(
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
				PlacedStatefulOrders:    []types.Order{},
				ExpiredStatefulOrderIds: []types.OrderId{},
				OrdersIdsFilledInLastBlock: []types.OrderId{
					constants.Order_Alice_Num1_Id13_Clob0_Buy50_Price50_GTB30.OrderId,
				},
				BlockHeight: blockHeight,
			},
		},
		"correctly handles new stateful orders": {
			operations: []types.Operation{
				types.NewOrderPlacementOperation(
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
				),
				types.NewOrderPlacementOperation(
					constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10,
				),
			},
			expectedProcessProposerMatchesEvents: types.ProcessProposerMatchesEvents{
				PlacedStatefulOrders: []types.Order{
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
					constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10,
				},
				ExpiredStatefulOrderIds:    []types.OrderId{},
				OrdersIdsFilledInLastBlock: []types.OrderId{},
				BlockHeight:                blockHeight,
			},
		},
		"correctly handles stateful order cancellations": {
			operations: []types.Operation{
				types.NewOrderPlacementOperation(
					constants.Order_Alice_Num0_Id9_Clob1_Buy15_Price45_GTB19,
				),
				types.NewOrderPlacementOperation(
					constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10,
				),
				types.NewMatchOperation(
					&constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num0_Id9_Clob1_Buy15_Price45_GTB19.OrderId,
							FillAmount:   10,
						},
					},
				),
				types.NewOrderCancellationOperation(
					&types.MsgCancelOrder{
						OrderId:      constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10.OrderId,
						GoodTilOneof: &types.MsgCancelOrder_GoodTilBlockTime{GoodTilBlockTime: 20},
					},
				),
			},
			expectedProcessProposerMatchesEvents: types.ProcessProposerMatchesEvents{
				PlacedStatefulOrders:    []types.Order{},
				ExpiredStatefulOrderIds: []types.OrderId{},
				OrdersIdsFilledInLastBlock: []types.OrderId{
					constants.Order_Alice_Num0_Id9_Clob1_Buy15_Price45_GTB19.OrderId,
					constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10.OrderId,
				},
				BlockHeight: blockHeight,
			},
		},
		"correctly handles stateful order replacements": {
			operations: []types.Operation{
				types.NewOrderPlacementOperation(
					constants.Order_Alice_Num0_Id9_Clob1_Buy15_Price45_GTB19,
				),
				types.NewOrderPlacementOperation(
					constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10,
				),
				types.NewMatchOperation(
					&constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num0_Id9_Clob1_Buy15_Price45_GTB19.OrderId,
							FillAmount:   10,
						},
					},
				),
				types.NewOrderPlacementOperation(
					constants.LongTermOrder_User1_Num1_Id1_Clob0_Sell50_Price30_GTBT15,
				),
			},
			expectedProcessProposerMatchesEvents: types.ProcessProposerMatchesEvents{
				PlacedStatefulOrders: []types.Order{
					constants.LongTermOrder_User1_Num1_Id1_Clob0_Sell50_Price30_GTBT15,
				},
				ExpiredStatefulOrderIds: []types.OrderId{},
				OrdersIdsFilledInLastBlock: []types.OrderId{
					constants.Order_Alice_Num0_Id9_Clob1_Buy15_Price45_GTB19.OrderId,
					constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10.OrderId,
				},
				BlockHeight: blockHeight,
			},
		},
		"order is cancelled and then replaced": {
			operations: []types.Operation{
				types.NewOrderPlacementOperation(
					constants.Order_Alice_Num0_Id9_Clob1_Buy15_Price45_GTB19,
				),
				types.NewOrderPlacementOperation(
					constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10,
				),
				types.NewMatchOperation(
					&constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num0_Id9_Clob1_Buy15_Price45_GTB19.OrderId,
							FillAmount:   10,
						},
					},
				),
				types.NewOrderCancellationOperation(
					&types.MsgCancelOrder{
						OrderId:      constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10.OrderId,
						GoodTilOneof: &types.MsgCancelOrder_GoodTilBlockTime{GoodTilBlockTime: 10},
					},
				),
				types.NewOrderPlacementOperation(
					constants.LongTermOrder_User1_Num1_Id1_Clob0_Sell50_Price30_GTBT15,
				),
			},
			expectedProcessProposerMatchesEvents: types.ProcessProposerMatchesEvents{
				PlacedStatefulOrders: []types.Order{
					constants.LongTermOrder_User1_Num1_Id1_Clob0_Sell50_Price30_GTBT15,
				},
				ExpiredStatefulOrderIds: []types.OrderId{},
				OrdersIdsFilledInLastBlock: []types.OrderId{
					constants.Order_Alice_Num0_Id9_Clob1_Buy15_Price45_GTB19.OrderId,
					constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10.OrderId,
				},
				BlockHeight: blockHeight,
			},
		},
		"order is replaced and then cancelled": {
			operations: []types.Operation{
				types.NewOrderPlacementOperation(
					constants.Order_Alice_Num0_Id9_Clob1_Buy15_Price45_GTB19,
				),
				types.NewOrderPlacementOperation(
					constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10,
				),
				types.NewMatchOperation(
					&constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num0_Id9_Clob1_Buy15_Price45_GTB19.OrderId,
							FillAmount:   10,
						},
					},
				),
				types.NewOrderPlacementOperation(
					constants.LongTermOrder_User1_Num1_Id1_Clob0_Sell50_Price30_GTBT15,
				),
				types.NewOrderCancellationOperation(
					&types.MsgCancelOrder{
						OrderId:      constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10.OrderId,
						GoodTilOneof: &types.MsgCancelOrder_GoodTilBlockTime{GoodTilBlockTime: 20},
					},
				),
			},
			expectedProcessProposerMatchesEvents: types.ProcessProposerMatchesEvents{
				PlacedStatefulOrders:    []types.Order{},
				ExpiredStatefulOrderIds: []types.OrderId{},
				OrdersIdsFilledInLastBlock: []types.OrderId{
					constants.Order_Alice_Num0_Id9_Clob1_Buy15_Price45_GTB19.OrderId,
					constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10.OrderId,
				},
				BlockHeight: blockHeight,
			},
		},
		"pre-existing stateful orders are skipped": {
			operations: []types.Operation{
				types.NewOrderPlacementOperation(
					constants.Order_Alice_Num0_Id9_Clob1_Buy15_Price45_GTB19,
				),
				types.NewPreexistingStatefulOrderPlacementOperation(
					constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10,
				),
				types.NewMatchOperation(
					&constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num0_Id9_Clob1_Buy15_Price45_GTB19.OrderId,
							FillAmount:   10,
						},
					},
				),
			},
			expectedProcessProposerMatchesEvents: types.ProcessProposerMatchesEvents{
				PlacedStatefulOrders:    []types.Order{},
				ExpiredStatefulOrderIds: []types.OrderId{},
				OrdersIdsFilledInLastBlock: []types.OrderId{
					constants.Order_Alice_Num0_Id9_Clob1_Buy15_Price45_GTB19.OrderId,
					constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10.OrderId,
				},
				BlockHeight: blockHeight,
			},
		},
	}

	// Run tests.
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			memclob := memclob.NewMemClobPriceTimePriority(true)
			ctx,
				keeper,
				_,
				_,
				_,
				_,
				_,
				_ := keepertest.ClobKeepers(t, memclob, &mocks.BankKeeper{}, &mocks.IndexerEventManager{})
			ctx = ctx.WithBlockHeight(int64(blockHeight))

			processProposerMatchesEvents := keeper.GenerateProcessProposerMatchesEvents(ctx, tc.operations)
			tc.expectedProcessProposerMatchesEvents.OperationsProposedInLastBlock = tc.operations
			require.Equal(t, tc.expectedProcessProposerMatchesEvents, processProposerMatchesEvents)
		})
	}
}

func setupNewMockEventManager(
	t *testing.T,
	ctx sdk.Context,
	mockIndexerEventManager *mocks.IndexerEventManager,
	matches []*types.MatchWithOrders,
	statefulOrderPlacements []types.Order,
	statefulOrderCancelations []types.OrderId,
) {
	if len(matches) > 0 {
		mockIndexerEventManager.On("Enabled").Return(true)
	}

	// TODO(CLOB-244) Support off-chain messages here as well.
	var statefulPlacementCallMap = make(map[types.OrderId]*mock.Call)
	for _, statefulOrder := range statefulOrderPlacements {
		call := mockIndexerEventManager.On("AddTxnEvent",
			mock.Anything,
			indexerevents.SubtypeStatefulOrder,
			indexer_manager.GetB64EncodedEventMessage(
				indexerevents.NewStatefulOrderPlacementEvent(
					statefulOrder,
				),
			),
		)
		statefulPlacementCallMap[statefulOrder.OrderId] = call
	}

	// Add an expectation to the mock for each expected message.
	var matchOrderCallMap = make(map[types.OrderId]*mock.Call)
	for _, match := range matches {
		if match.TakerOrder.IsLiquidation() {
			call := mockIndexerEventManager.On("AddTxnEvent",
				mock.Anything,
				indexerevents.SubtypeOrderFill,
				indexer_manager.GetB64EncodedEventMessage(
					indexerevents.NewLiquidationOrderFillEvent(
						match.MakerOrder.MustGetOrder(),
						match.TakerOrder,
						match.FillAmount,
					),
				),
			).Return()

			// Stateful orders should not emit match events before placement events.
			if placementCall, exists := statefulPlacementCallMap[match.MakerOrder.MustGetOrder().OrderId]; exists {
				call.NotBefore(placementCall)
			}

			matchOrderCallMap[match.MakerOrder.MustGetOrder().OrderId] = call
		} else {
			call := mockIndexerEventManager.On("AddTxnEvent",
				mock.Anything,
				indexerevents.SubtypeOrderFill,
				indexer_manager.GetB64EncodedEventMessage(
					indexerevents.NewOrderFillEvent(
						match.MakerOrder.MustGetOrder(),
						match.TakerOrder.MustGetOrder(),
						match.FillAmount,
					),
				),
			).Return()

			// Stateful orders should not emit match events before placement events.
			if placementCall, exists := statefulPlacementCallMap[match.MakerOrder.MustGetOrder().OrderId]; exists {
				call.NotBefore(placementCall)
			}

			// Stateful orders should not emit match events before placement events.
			if placementCall, exists := statefulPlacementCallMap[match.TakerOrder.MustGetOrder().OrderId]; exists {
				call.NotBefore(placementCall)
			}

			matchOrderCallMap[match.MakerOrder.MustGetOrder().OrderId] = call
			matchOrderCallMap[match.TakerOrder.MustGetOrder().OrderId] = call
		}
	}

	for _, statefulOrderId := range statefulOrderCancelations {
		call := mockIndexerEventManager.On("AddTxnEvent",
			mock.Anything,
			indexerevents.SubtypeStatefulOrder,
			indexer_manager.GetB64EncodedEventMessage(
				indexerevents.NewStatefulOrderCancelationEvent(
					statefulOrderId,
				),
			),
		)

		// Stateful orders should not emit cancel events before placement events.
		if placementCall, exists := statefulPlacementCallMap[statefulOrderId]; exists {
			call.NotBefore(placementCall)
		}
		// Stateful orders should not emit cancel events before match events.
		if matchCall, exists := matchOrderCallMap[statefulOrderId]; exists {
			call.NotBefore(matchCall)
		}
	}
}

func assertSubaccountState(
	t *testing.T,
	ctx sdk.Context,
	saKeeper *sakeeper.Keeper,
	expectedQuoteBalances map[satypes.SubaccountId]*big.Int,
	expectedPerpetualPositions map[satypes.SubaccountId][]*satypes.PerpetualPosition,
) {
	for subaccountId, quoteBalance := range expectedQuoteBalances {
		subaccount := saKeeper.GetSubaccount(ctx, subaccountId)
		require.Equal(t, quoteBalance, subaccount.GetUsdcPosition())
	}

	for subaccountId, perpetualPositions := range expectedPerpetualPositions {
		subaccount := saKeeper.GetSubaccount(ctx, subaccountId)
		require.ElementsMatch(t, subaccount.PerpetualPositions, perpetualPositions)
	}
}
