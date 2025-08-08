package keeper_test

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	indexerevents "github.com/dydxprotocol/v4-chain/protocol/indexer/events"

	cmt "github.com/cometbft/cometbft/types"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	perptest "github.com/dydxprotocol/v4-chain/protocol/testutil/perpetuals"

	"github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	clobtest "github.com/dydxprotocol/v4-chain/protocol/testutil/clob"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	memclobtest "github.com/dydxprotocol/v4-chain/protocol/testutil/memclob"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/tracer"
	blocktimetypes "github.com/dydxprotocol/v4-chain/protocol/x/blocktime/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/memclob"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	feetypes "github.com/dydxprotocol/v4-chain/protocol/x/feetiers/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/perpetuals"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/prices"
	rewardtypes "github.com/dydxprotocol/v4-chain/protocol/x/rewards/types"
	statstypes "github.com/dydxprotocol/v4-chain/protocol/x/stats/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestPlaceShortTermOrder(t *testing.T) {
	tests := map[string]struct {
		// Perpetuals state.
		perpetuals []perptypes.Perpetual
		// Subaccount state.
		subaccounts []satypes.Subaccount
		// CLOB state.
		clobs          []types.ClobPair
		existingOrders []types.Order
		// Fee tier params.
		feeParams feetypes.PerpetualFeeParams

		// Parameters.
		order types.Order

		// Expectations.
		expectedMultiStoreWrites []string
		expectedOrderStatus      types.OrderStatus
		expectedFilledSize       satypes.BaseQuantums
		// Expected remaining OI after test.
		// The test initializes each perp with default open interest of 1 full coin.
		expectedOpenInterests map[uint32]*big.Int
		expectedErr           error
	}{
		"Can place an order on the orderbook closing a position": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_SmallMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short,
			},
			clobs:     []types.ClobPair{constants.ClobPair_Btc},
			feeParams: constants.PerpetualFeeParams,

			order: constants.Order_Carl_Num0_Id1_Clob0_Buy1BTC_Price49999,

			expectedOrderStatus: types.Success,
			expectedFilledSize:  0,
			expectedOpenInterests: map[uint32]*big.Int{
				// unchanged, no match happened
				constants.BtcUsd_SmallMarginRequirement.Params.Id: big.NewInt(100_000_000),
			},
		},
		"Can place an order with a valid order router address": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_SmallMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short,
			},
			clobs:     []types.ClobPair{constants.ClobPair_Btc},
			feeParams: constants.PerpetualFeeParams,

			order: constants.Order_Carl_Num0_Id1_Clob0_Buy1BTC_WithValidOrderRouter,

			expectedOrderStatus: types.Success,
			expectedFilledSize:  0,
			expectedOpenInterests: map[uint32]*big.Int{
				// unchanged, no match happened
				constants.BtcUsd_SmallMarginRequirement.Params.Id: big.NewInt(100_000_000),
			},
		},
		"Can place an order on the orderbook in a different market than their current perpetual position": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_SmallMarginRequirement,
				constants.EthUsd_20PercentInitial_10PercentMaintenance,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short,
			},
			clobs: []types.ClobPair{
				constants.ClobPair_Btc,
				constants.ClobPair_Eth,
			},
			feeParams: constants.PerpetualFeeParams,

			order: constants.Order_Carl_Num0_Id2_Clob1_Buy10ETH_Price3000,

			expectedOrderStatus: types.Success,
			expectedFilledSize:  0,
			expectedOpenInterests: map[uint32]*big.Int{
				// unchanged, no match happened
				constants.BtcUsd_SmallMarginRequirement.Params.Id: big.NewInt(100_000_000),
			},
		},
		"Can place an order and the order is fully matched": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_SmallMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short,
				constants.Dave_Num0_1BTC_Long_50000USD,
			},
			clobs: []types.ClobPair{
				constants.ClobPair_Btc,
			},
			existingOrders: []types.Order{
				constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000,
			},
			feeParams: constants.PerpetualFeeParams,

			order: constants.Order_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTB10,

			expectedOrderStatus: types.Success,
			expectedFilledSize:  constants.Order_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTB10.GetBaseQuantums(),
			expectedMultiStoreWrites: []string{
				// Update taker subaccount
				satypes.SubaccountKeyPrefix +
					string(constants.Carl_Num0.ToStateKey()),
				// Indexer event
				indexer_manager.IndexerEventsCountKey,
				// Update maker subaccount
				satypes.SubaccountKeyPrefix +
					string(constants.Dave_Num0.ToStateKey()),
				// Indexer event
				indexer_manager.IndexerEventsCountKey,
				// Update rewards
				rewardtypes.RewardShareKeyPrefix + constants.Carl_Num0.Owner,
				rewardtypes.RewardShareKeyPrefix + constants.Dave_Num0.Owner,
				// Update block stats
				statstypes.BlockStatsKey,
				// Update prunable block height for taker fill amount
				types.PrunableOrdersKeyPrefix,
				// Update taker order fill amount
				types.OrderAmountFilledKeyPrefix,
				// Update prunable block height for maker fill amount
				types.PrunableOrdersKeyPrefix,
				// Update maker order fill amount
				types.OrderAmountFilledKeyPrefix,
			},
			expectedOpenInterests: map[uint32]*big.Int{
				// positions fully closed
				constants.BtcUsd_SmallMarginRequirement.Params.Id: big.NewInt(0),
			},
		},
		"Can place an order on the orderbook if the subaccount is right at the initial margin ratio": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short,
			},
			clobs: []types.ClobPair{
				constants.ClobPair_Btc,
			},
			feeParams: constants.PerpetualFeeParamsNoFee,

			order: constants.Order_Carl_Num0_Id0_Clob0_Buy10QtBTC_Price100000QuoteQt,

			expectedOrderStatus: types.Success,
			expectedFilledSize:  0,
			expectedOpenInterests: map[uint32]*big.Int{
				// unchanged, no match happened
				constants.BtcUsd_SmallMarginRequirement.Params.Id: big.NewInt(100_000_000),
			},
		},
		"Can place an order on the orderbook if the account would be collateralized due to rebate": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short,
			},
			clobs: []types.ClobPair{
				constants.ClobPair_Btc,
			},
			// Same setup as the above two tests, but the order is for a slightly higher price that
			// cannot be collateralized without the rebate.
			feeParams: constants.PerpetualFeeParamsMakerRebate,

			order: constants.Order_Carl_Num0_Id0_Clob0_Buy10QtBTC_Price100001QuoteQt,

			expectedOrderStatus: types.Success,
			expectedFilledSize:  0,
			expectedOpenInterests: map[uint32]*big.Int{
				// unchanged, no match happened
				constants.BtcUsd_SmallMarginRequirement.Params.Id: big.NewInt(100_000_000),
			},
		},
		"Cannot open an order if it doesn't reference a valid CLOB": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_SmallMarginRequirement,
				constants.EthUsd_20PercentInitial_10PercentMaintenance,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_599USD,
			},
			clobs: []types.ClobPair{
				constants.ClobPair_Btc,
			},
			feeParams: constants.PerpetualFeeParams,

			order: constants.Order_Carl_Num0_Id3_Clob1_Buy1ETH_Price3000,

			expectedErr: types.ErrInvalidClob,
			expectedOpenInterests: map[uint32]*big.Int{
				// unchanged, no match happened
				constants.BtcUsd_SmallMarginRequirement.Params.Id: big.NewInt(100_000_000),
			},
		},
		"Cannot open an order if the subticks are invalid": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_SmallMarginRequirement,
				constants.EthUsd_20PercentInitial_10PercentMaintenance,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_599USD,
			},
			clobs: []types.ClobPair{
				constants.ClobPair_Btc,
			},
			feeParams: constants.PerpetualFeeParams,

			order: types.Order{
				OrderId: types.OrderId{
					SubaccountId: constants.Alice_Num1, ClientId: 1,
				},
				Subticks: 2,
			},

			expectedErr: types.ErrInvalidPlaceOrder,
			expectedOpenInterests: map[uint32]*big.Int{
				// unchanged, no match happened
				constants.BtcUsd_SmallMarginRequirement.Params.Id: big.NewInt(100_000_000),
			},
		},
		"Cannot open an order that is smaller than the minimum base quantums": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_SmallMarginRequirement,
				constants.EthUsd_20PercentInitial_10PercentMaintenance,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_599USD,
			},
			clobs: []types.ClobPair{
				constants.ClobPair_Btc,
			},
			feeParams: constants.PerpetualFeeParams,

			order: types.Order{
				OrderId: types.OrderId{
					SubaccountId: constants.Alice_Num1, ClientId: 1,
				},
				Quantums: 1,
			},

			expectedErr: types.ErrInvalidPlaceOrder,
			expectedOpenInterests: map[uint32]*big.Int{
				// unchanged, no match happened
				constants.BtcUsd_SmallMarginRequirement.Params.Id: big.NewInt(100_000_000),
			},
		},
		"Cannot open an order that is not divisible by the step size base quantums": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_SmallMarginRequirement,
				constants.EthUsd_20PercentInitial_10PercentMaintenance,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_599USD,
			},
			clobs: []types.ClobPair{
				constants.ClobPair_Btc,
			},
			feeParams: constants.PerpetualFeeParams,

			order: types.Order{
				OrderId:  types.OrderId{},
				Quantums: 11,
			},

			expectedErr: types.ErrInvalidPlaceOrder,
			expectedOpenInterests: map[uint32]*big.Int{
				// unchanged, no match happened
				constants.BtcUsd_SmallMarginRequirement.Params.Id: big.NewInt(100_000_000),
			},
		},
		"Cannot open an order with a GoodTilBlock in the past": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_SmallMarginRequirement,
				constants.EthUsd_20PercentInitial_10PercentMaintenance,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_599USD,
			},
			clobs: []types.ClobPair{
				constants.ClobPair_Btc,
			},
			feeParams: constants.PerpetualFeeParams,

			order: types.Order{
				OrderId:      types.OrderId{},
				Quantums:     10,
				GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 0},
			},

			expectedErr: types.ErrHeightExceedsGoodTilBlock,
			expectedOpenInterests: map[uint32]*big.Int{
				// unchanged, no match happened
				constants.BtcUsd_SmallMarginRequirement.Params.Id: big.NewInt(100_000_000),
			},
		},
		"Cannot open an order with a GoodTilBlock greater than ShortBlockWindow blocks in the future": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_SmallMarginRequirement,
				constants.EthUsd_20PercentInitial_10PercentMaintenance,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_599USD,
			},
			clobs: []types.ClobPair{
				constants.ClobPair_Btc,
			},
			feeParams: constants.PerpetualFeeParams,

			order: types.Order{
				OrderId:      types.OrderId{},
				Quantums:     10,
				GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 2 + types.ShortBlockWindow},
			},

			expectedErr: types.ErrGoodTilBlockExceedsShortBlockWindow,
			expectedOpenInterests: map[uint32]*big.Int{
				// unchanged, no match happened
				constants.BtcUsd_SmallMarginRequirement.Params.Id: big.NewInt(100_000_000),
			},
		},
		"Can open an order with an invalid order router address": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short,
			},
			clobs: []types.ClobPair{
				constants.ClobPair_Btc,
			},
			// Same setup as the above two tests, but the order is for a slightly higher price that
			// cannot be collateralized without the rebate.
			feeParams: constants.PerpetualFeeParamsMakerRebate,

			order: types.Order{
				OrderId:            types.OrderId{SubaccountId: constants.Carl_Num0, ClientId: 0, ClobPairId: 0},
				Side:               types.Order_SIDE_BUY,
				Quantums:           10,
				Subticks:           100_001_000_000,
				GoodTilOneof:       &types.Order_GoodTilBlock{GoodTilBlock: 20},
				OrderRouterAddress: constants.CarlAccAddress.String(),
			},

			expectedOrderStatus: types.Success,
			expectedFilledSize:  0,
			expectedOpenInterests: map[uint32]*big.Int{
				// unchanged, no match happened
				constants.BtcUsd_100PercentMarginRequirement.Params.Id: big.NewInt(100_000_000),
			},
		},
		`New order should be undercollateralized when matching when previous fills make it undercollateralized when using
				maker orders subticks, but would be collateralized if using taker order subticks`: {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_50PercentInitial_40PercentMaintenance,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num1_500USD,
				constants.Carl_Num0_10000USD,
			},
			clobs: []types.ClobPair{
				constants.ClobPair_Btc,
			},
			existingOrders: []types.Order{
				// Alice places sell order for $500 worth of BTC.
				constants.Order_Carl_Num0_Id0_Clob0_Sell1kQtBTC_Price50000,
				// Bob places a buy order to buy $500 worth of BTC at $600.
				// This completely fills Alice's order at the maker price of $500.
				constants.Order_Carl_Num1_Id0_Clob0_Buy1kQtBTC_Price60000,
				// Orderbook is now empty.
				// Alice places another sell order for $500 worth of BTC which now rests on the book.
				constants.Order_Carl_Num0_Id1_Clob0_Sell1kQtBTC_Price50000,
			},
			feeParams: constants.PerpetualFeeParamsNoFee,
			// Bob places a second order at a lower price than his first.
			// This should succeed, because Bob currently has a balance of $0 USDC,
			// and a perpetual of size 0.01 BTC ($500), and the perpetual has a 50% margin requirement.
			// This should bring Bob's balance to -500$ USD, and 0.02 BTC which is exactly at the 50% margin requirement.
			// If the Bob's previous order was viewed as being filled at the taker subticks (60_000_000_000) instead when
			// checking collateralization, the match would fail.
			order:               constants.Order_Carl_Num1_Id1_Clob0_Buy1kQtBTC_Price50000,
			expectedOrderStatus: types.Success,
			expectedFilledSize:  1_000_000,
			expectedMultiStoreWrites: []string{
				// Update taker subaccount
				satypes.SubaccountKeyPrefix +
					string(constants.Carl_Num1.ToStateKey()),
				indexer_manager.IndexerEventsCountKey,
				// Update maker subaccount
				satypes.SubaccountKeyPrefix +
					string(constants.Carl_Num0.ToStateKey()),
				indexer_manager.IndexerEventsCountKey,
				// Update block stats
				statstypes.BlockStatsKey,
				// Update prunable block height for taker fill amount
				types.PrunableOrdersKeyPrefix,
				// Update taker order fill amount
				types.OrderAmountFilledKeyPrefix,
				// Update prunable block height for maker fill amount
				types.PrunableOrdersKeyPrefix,
				// Update maker order fill amount
				types.OrderAmountFilledKeyPrefix,
			},
			expectedOpenInterests: map[uint32]*big.Int{
				// 1 BTC + 0.01 BTC + 0.01 BTC filled
				constants.BtcUsd_50PercentInitial_40PercentMaintenance.Params.Id: big.NewInt(102_000_000),
			},
		},
		// This is a regression test for an issue whereby orders that had been previously matched were being checked for
		// collateralization as if the CLOB pair ID of the order was `0`. This resulted in always using `0`
		// the quantum conversion exponent of the first CLOB when performing collateralization checks during `PlaceOrder`.
		// This lead to an incorrect amount of quote quantums being returned for all previously matched orders
		// that weren't placed on the first CLOB. If firstClobPair.QuantumConversionExponent >
		// expectedClobPair.QuantumConversionExponent, then sellers receive more quote quantums and buyers are charged
		// more. Vice versa if firstClobPair.QuantumConversionExponent < expectedClobPair.QuantumConversionExponent.
		// Context: https://github.com/dydxprotocol/v4-chain/protocol/pull/562#discussion_r1024319468
		`Regression: New order should be fully collateralized when matching with previous fills
				because the correct quantum conversion exponent was used`: {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_NoMarginRequirement,
				constants.EthUsd_20PercentInitial_10PercentMaintenance,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_660USD,
				constants.Dave_Num0_10000USD,
			},
			clobs: []types.ClobPair{
				{
					Id: 0,
					Metadata: &types.ClobPair_PerpetualClobMetadata{
						PerpetualClobMetadata: &types.PerpetualClobMetadata{
							PerpetualId: 0,
						},
					},
					Status:                    types.ClobPair_STATUS_ACTIVE,
					StepBaseQuantums:          10,
					SubticksPerTick:           100,
					QuantumConversionExponent: -1,
				},
				constants.ClobPair_Eth_No_Fee,
			},
			existingOrders: []types.Order{
				// Alice places buy order for $3,000 worth of ETH.
				constants.Order_Carl_Num0_Id3_Clob1_Buy1ETH_Price3000,
				// Bob places sell order for $3,000 worth of ETH.
				// This completely fills Alice's order at the maker price of $3,000.
				constants.Order_Dave_Num0_Id3_Clob1_Sell1ETH_Price3000,
			},
			feeParams: constants.PerpetualFeeParamsNoFee,
			// Bob places a second buy order at a the same price as the first.
			// This should bring his total initial margin requirement to $3,300 * 20%, or $660
			// which is exactly equal to the quote balance.
			// However, if the quantum conversion exponent of the first perpetual CLOB pair
			// is used for the collateralization check, it will over-report the amount of quote
			// quantums required to open the previous buy order and fail.
			order:               constants.Order_Carl_Num0_Id4_Clob1_Buy01ETH_Price3000,
			expectedOrderStatus: types.Success,
			expectedOpenInterests: map[uint32]*big.Int{
				// Unchanged, no BTC match happened
				constants.BtcUsd_NoMarginRequirement.Params.Id: big.NewInt(100_000_000),
				// 1 ETH + 1 ETH filled
				constants.EthUsd_20PercentInitial_10PercentMaintenance.Params.Id: big.NewInt(2_000_000_000),
			},
		},
		// <grouped tests: deprecating pessimistic value collateralization check -- BUY>
		// The following 3 tests are a group to test the deprecation of pessimistic collateralization check.
		// 1. The first should have a buy price well below the oracle price of 50,000. (success)
		// 2. The second should have a buy price above the oracle price of 50,000. (undercollateralized)
		// 3. The third should have the order in common with #1 and the subaccount in common with #2 and should succeed.
		`Subaccount can now place buy order that would have failed the 
				deprecated pessimistic collateralization check with the oracle price`: {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_50PercentInitial_40PercentMaintenance,
			},
			subaccounts: []satypes.Subaccount{constants.Carl_Num0_10000USD},
			clobs: []types.ClobPair{
				constants.ClobPair_Btc,
			},
			existingOrders:           []types.Order{},
			feeParams:                constants.PerpetualFeeParamsNoFee,
			order:                    constants.Order_Carl_Num0_Id0_Clob0_Buy1BTC_Price5subticks_GTB10,
			expectedOrderStatus:      types.Success,
			expectedFilledSize:       0,
			expectedMultiStoreWrites: []string{},
		},
		`Subaccount placed buy order passes collateralization check when using the maker price`: {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_50PercentInitial_40PercentMaintenance,
			},
			subaccounts: []satypes.Subaccount{constants.Carl_Num0_100000USD},
			clobs: []types.ClobPair{
				constants.ClobPair_Btc,
			},
			existingOrders:           []types.Order{},
			feeParams:                constants.PerpetualFeeParamsNoFee,
			order:                    constants.Order_Carl_Num0_Id0_Clob0_Buy1BTC_Price5subticks_GTB10,
			expectedOrderStatus:      types.Success,
			expectedFilledSize:       0,
			expectedMultiStoreWrites: []string{},
		},
		// <end of grouped tests: pessimistic value collateralization check -- BUY>
		// <grouped tests: deprecating pessimistic value collateralization check -- SELL>
		// The following 3 tests are a group to test the deprecation of pessimistic collateralization check.
		// 1. The first should have a sell price well above the oracle price of 50,000. (success)
		// 2. The second should have a sell price below the oracle price of 50,000 subticks. (undercollateralized)
		// 3. The third should have the order in common with #1 and the subaccount in common with #2 and should succeed.
		`Subaccount can now place sell order that would have failed the 
				deprecated pessimistic collateralization check with the oracle price`: {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_50PercentInitial_40PercentMaintenance,
			},
			subaccounts: []satypes.Subaccount{constants.Carl_Num0_10000USD},
			clobs: []types.ClobPair{
				constants.ClobPair_Btc,
			},
			existingOrders:           []types.Order{},
			feeParams:                constants.PerpetualFeeParamsNoFee,
			order:                    constants.Order_Carl_Num0_Id0_Clob0_Sell1BTC_Price500000_GTB10,
			expectedOrderStatus:      types.Success,
			expectedFilledSize:       0,
			expectedMultiStoreWrites: []string{},
		},
		`Subaccount placed sell order passes collateralization check when using the maker price`: {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_50PercentInitial_40PercentMaintenance,
			},
			subaccounts: []satypes.Subaccount{constants.Carl_Num0_50000USD},
			clobs: []types.ClobPair{
				constants.ClobPair_Btc,
			},
			existingOrders:           []types.Order{},
			feeParams:                constants.PerpetualFeeParamsNoFee,
			order:                    constants.Order_Carl_Num0_Id0_Clob0_Sell1BTC_Price500000_GTB10,
			expectedOrderStatus:      types.Success,
			expectedFilledSize:       0,
			expectedMultiStoreWrites: []string{},
		},
		// <end of grouped tests: pessimistic value collateralization check -- SELL>
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup keeper state.
			memClob := memclob.NewMemClobPriceTimePriority(false)
			mockBankKeeper := &mocks.BankKeeper{}
			mockBankKeeper.On(
				"SendCoins",
				mock.Anything,
				mock.Anything,
				mock.Anything,
				mock.Anything,
			).Return(nil)

			ks := keepertest.NewClobKeepersTestContext(t, memClob, mockBankKeeper, indexer_manager.NewIndexerEventManagerNoop())
			ctx := ks.Ctx.WithIsCheckTx(true)

			// Create the default markets.
			keepertest.CreateTestMarkets(t, ctx, ks.PricesKeeper)

			// Create liquidity tiers.
			keepertest.CreateTestLiquidityTiers(t, ctx, ks.PerpetualsKeeper)

			require.NoError(t, ks.FeeTiersKeeper.SetPerpetualFeeParams(ctx, tc.feeParams))

			// Set up USDC asset in assets module.
			err := keepertest.CreateUsdcAsset(ctx, ks.AssetsKeeper)
			require.NoError(t, err)

			err = ks.PricesKeeper.RevShareKeeper.SetOrderRouterRevShare(
				ctx, constants.AliceAccAddress.String(), 100_000)
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
			for _, clobPair := range tc.clobs {
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

			err = ks.ClobKeeper.InitializeEquityTierLimit(
				ctx,
				types.EquityTierLimitConfiguration{
					ShortTermOrderEquityTiers: []types.EquityTierLimit{
						{
							UsdTncRequired: dtypes.NewInt(20_000_000),
							Limit:          5,
						},
					},
					StatefulOrderEquityTiers: []types.EquityTierLimit{
						{
							UsdTncRequired: dtypes.NewInt(20_000_000),
							Limit:          5,
						},
					},
				},
			)
			require.NoError(t, err)

			// Create all existing orders.
			for _, order := range tc.existingOrders {
				msg := &types.MsgPlaceOrder{Order: order}

				txBuilder := constants.TestEncodingCfg.TxConfig.NewTxBuilder()
				err := txBuilder.SetMsgs(msg)
				require.NoError(t, err)
				bytes, err := constants.TestEncodingCfg.TxConfig.TxEncoder()(txBuilder.GetTx())
				require.NoError(t, err)
				ctx = ctx.WithTxBytes(bytes)

				_, _, err = ks.ClobKeeper.PlaceShortTermOrder(ctx, msg)
				require.NoError(t, err)
			}

			// Run the test.
			traceDecoder := &tracer.TraceDecoder{}
			ctx.MultiStore().SetTracer(traceDecoder)

			msg := &types.MsgPlaceOrder{Order: tc.order}

			txBuilder := constants.TestEncodingCfg.TxConfig.NewTxBuilder()
			err = txBuilder.SetMsgs(msg)
			require.NoError(t, err)
			bytes, err := constants.TestEncodingCfg.TxConfig.TxEncoder()(txBuilder.GetTx())
			require.NoError(t, err)
			ctx = ctx.WithTxBytes(bytes)

			orderSizeOptimisticallyFilledFromMatching,
				orderStatus,
				err := ks.ClobKeeper.PlaceShortTermOrder(ctx, msg)

			// Verify test expectations.
			require.ErrorIs(t, err, tc.expectedErr)
			if err == nil {
				require.Equal(t, tc.expectedOrderStatus, orderStatus)
				require.Equal(t, tc.expectedFilledSize, orderSizeOptimisticallyFilledFromMatching)
			}

			traceDecoder.RequireKeyPrefixesWritten(t, tc.expectedMultiStoreWrites)

			for _, perp := range tc.perpetuals {
				if expectedOI, exists := tc.expectedOpenInterests[perp.Params.Id]; exists {
					gotPerp, err := ks.PerpetualsKeeper.GetPerpetual(ks.Ctx, perp.Params.Id)
					require.NoError(t, err)
					require.Zero(t,
						expectedOI.Cmp(gotPerp.OpenInterest.BigInt()),
						"expected open interest %s, got %s",
						expectedOI.String(),
						gotPerp.OpenInterest.String(),
					)
				}
			}
		})
	}
}

func TestPlaceShortTermOrder_PanicsOnStatefulOrder(t *testing.T) {
	memClob := memclob.NewMemClobPriceTimePriority(false)
	mockBankKeeper := &mocks.BankKeeper{}
	ks := keepertest.NewClobKeepersTestContext(t, memClob, mockBankKeeper, indexer_manager.NewIndexerEventManagerNoop())
	msgPlaceOrder := &types.MsgPlaceOrder{Order: constants.LongTermOrder_Bob_Num0_Id2_Clob0_Buy15_Price5_GTBT10}

	require.Panicsf(
		t,
		func() {
			//nolint:errcheck
			ks.ClobKeeper.PlaceShortTermOrder(ks.Ctx, msgPlaceOrder)
		},
		"MustBeShortTermOrder: called with stateful order ID (%+v)",
		msgPlaceOrder.Order.OrderId,
	)
}

// TODO(DEC-1648): Add test cases for additional order placement scenarios.
func TestAddPreexistingStatefulOrder(t *testing.T) {
	tests := map[string]struct {
		// Perpetuals state.
		perpetuals []perptypes.Perpetual
		// Subaccount state.
		subaccounts []satypes.Subaccount
		// CLOB state.
		clobs          []types.ClobPair
		existingOrders []types.Order

		// Parameters.
		order types.Order

		// Expectations.
		expectedMultiStoreWrites []string
		expectedOrderStatus      types.OrderStatus
		expectedFilledSize       satypes.BaseQuantums
		expectedErr              error
		expectedTransactionIndex uint32
	}{
		"Can place a stateful order on the orderbook": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_SmallMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short,
			},
			clobs: []types.ClobPair{constants.ClobPair_Btc},

			order: constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10,

			expectedOrderStatus:      types.Success,
			expectedFilledSize:       0,
			expectedTransactionIndex: 0,
		},
		`Can place multiple stateful orders on the orderbook, and the newly-placed stateful order
			matches`: {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_SmallMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num0_10_000USD,
				constants.Carl_Num0_1BTC_Short,
				constants.Dave_Num0_1BTC_Long_50000USD,
			},
			clobs: []types.ClobPair{constants.ClobPair_Btc},
			existingOrders: []types.Order{
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20,
				constants.LongTermOrder_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10,
			},

			order: constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10,

			expectedOrderStatus:      types.Success,
			expectedFilledSize:       100_000_000,
			expectedTransactionIndex: 2,
			expectedMultiStoreWrites: []string{
				// Update taker subaccount.
				satypes.SubaccountKeyPrefix +
					string(constants.Carl_Num0.ToStateKey()),
				indexer_manager.IndexerEventsCountKey,
				// Update maker subaccount.
				satypes.SubaccountKeyPrefix +
					string(constants.Dave_Num0.ToStateKey()),
				indexer_manager.IndexerEventsCountKey,
				// Update block stats
				statstypes.BlockStatsKey,
				// Update taker order fill amount to state.
				types.OrderAmountFilledKeyPrefix +
					string(constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10.OrderId.ToStateKey()),
				// Update maker order fill amount to state.
				types.OrderAmountFilledKeyPrefix +
					string(constants.LongTermOrder_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10.OrderId.ToStateKey()),
			},
		},
		`Can place a stateful post-only order that crosses the book`: {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_SmallMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num0_10_000USD,
				constants.Alice_Num1_10_000USD,
			},
			clobs: []types.ClobPair{constants.ClobPair_Btc},
			existingOrders: []types.Order{
				constants.LongTermOrder_Alice_Num1_Id4_Clob0_Buy10_Price45_GTBT20,
			},

			order: constants.LongTermOrder_Alice_Num0_Id2_Clob0_Sell65_Price10_GTBT25_PO,

			expectedOrderStatus:      types.Success,
			expectedErr:              types.ErrPostOnlyWouldCrossMakerOrder,
			expectedFilledSize:       0,
			expectedTransactionIndex: 0,
			// No multi store writes due to ErrPostOnlyWouldCrossMakerOrder error.
			expectedMultiStoreWrites: []string{},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup keeper state.
			memClob := memclob.NewMemClobPriceTimePriority(false)
			mockBankKeeper := &mocks.BankKeeper{}
			mockBankKeeper.On(
				"SendCoins",
				mock.Anything,
				mock.Anything,
				mock.Anything,
				mock.Anything,
			).Return(nil)

			ks := keepertest.NewClobKeepersTestContext(t, memClob, mockBankKeeper, indexer_manager.NewIndexerEventManagerNoop())
			ctx := ks.Ctx.WithIsCheckTx(true)

			// Create the default markets.
			keepertest.CreateTestMarkets(t, ctx, ks.PricesKeeper)

			// Create liquidity tiers.
			keepertest.CreateTestLiquidityTiers(t, ctx, ks.PerpetualsKeeper)

			require.NoError(t, ks.FeeTiersKeeper.SetPerpetualFeeParams(ctx, constants.PerpetualFeeParamsNoFee))

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
			for _, clobPair := range tc.clobs {
				_, err := ks.ClobKeeper.CreatePerpetualClobPairAndMemStructs(
					ctx,
					clobPair.Id,
					clobPair.GetPerpetualClobMetadata().PerpetualId,
					satypes.BaseQuantums(clobPair.StepBaseQuantums),
					clobPair.QuantumConversionExponent,
					clobPair.SubticksPerTick,
					clobPair.Status,
				)
				require.NoError(t, err)
			}

			// Set the block height and last committed block time.
			blockHeight := uint32(2)
			ctx = ctx.WithIsCheckTx(false).
				WithBlockHeight(int64(blockHeight)).
				WithBlockTime(time.Unix(5, 0))
			ks.BlockTimeKeeper.SetPreviousBlockInfo(ctx, &blocktimetypes.BlockInfo{
				Height:    2,
				Timestamp: time.Unix(int64(5), 0),
			})

			// Create all existing orders.
			for _, order := range tc.existingOrders {
				if order.IsStatefulOrder() {
					ks.ClobKeeper.SetLongTermOrderPlacement(ctx, order, blockHeight)
					ks.ClobKeeper.AddStatefulOrderIdExpiration(
						ctx,
						order.MustGetUnixGoodTilBlockTime(),
						order.GetOrderId(),
					)
				}

				_, orderStatus, _, err := ks.ClobKeeper.AddPreexistingStatefulOrder(
					ctx.WithIsCheckTx(true),
					&order,
					memClob,
				)
				require.NoError(t, err)
				require.True(t, orderStatus.IsSuccess())
			}

			// Stateful orders are not written to state in PlaceOrder, they are written to state in DeliverTx.
			// Write stateful order to state here, as DeliverTx would have executed for each stateful order
			// by the time PlaceOrder was called, as PlaceOrder is called in PrepareCheckState for stateful orders.
			if tc.order.IsStatefulOrder() {
				ks.ClobKeeper.SetLongTermOrderPlacement(ctx, tc.order, blockHeight)
				ks.ClobKeeper.AddStatefulOrderIdExpiration(
					ctx,
					tc.order.MustGetUnixGoodTilBlockTime(),
					tc.order.GetOrderId(),
				)
			}

			// Verify the order that will be placed is a Long-Term order.
			require.True(t, tc.order.OrderId.IsLongTermOrder())

			// Run the test.
			traceDecoder := &tracer.TraceDecoder{}
			ctx.MultiStore().SetTracer(traceDecoder)

			orderSizeOptimisticallyFilledFromMatching,
				orderStatus,
				_,
				err := ks.ClobKeeper.AddPreexistingStatefulOrder(ctx.WithIsCheckTx(true), &tc.order, memClob)

			// Verify test expectations.
			require.ErrorIs(t, err, tc.expectedErr)
			statefulOrderPlacement, _ := ks.ClobKeeper.GetLongTermOrderPlacement(ctx, tc.order.OrderId)
			statefulOrderIds := ks.ClobKeeper.GetStatefulOrderIdExpirations(ctx, tc.order.MustGetUnixGoodTilBlockTime())
			if err == nil {
				require.Equal(t, tc.expectedOrderStatus, orderStatus)
				require.Equal(t, tc.expectedFilledSize, orderSizeOptimisticallyFilledFromMatching)
				require.Equal(t, tc.order, statefulOrderPlacement.Order)
				require.Equal(t, blockHeight, statefulOrderPlacement.PlacementIndex.BlockHeight)
				require.Equal(t, tc.expectedTransactionIndex, statefulOrderPlacement.PlacementIndex.TransactionIndex)
				require.Contains(
					t,
					statefulOrderIds,
					tc.order.OrderId,
				)
			}

			traceDecoder.RequireKeyPrefixesWritten(t, tc.expectedMultiStoreWrites)
		})
	}
}

func TestPlaceOrder_SendOffchainMessages(t *testing.T) {
	indexerEventManager := &mocks.IndexerEventManager{}
	for _, message := range constants.TestOffchainMessages {
		indexerEventManager.On(
			"SendOffchainData",
			message.AddHeader(constants.TestTxHashHeader),
		).Return().Once()
	}

	memClob := &mocks.MemClob{}
	memClob.On("SetClobKeeper", mock.Anything).Return()

	ks := keepertest.NewClobKeepersTestContext(t, memClob, &mocks.BankKeeper{}, indexerEventManager)

	ks.MarketMapKeeper.InitGenesis(ks.Ctx, constants.MarketMap_DefaultGenesisState)
	prices.InitGenesis(ks.Ctx, *ks.PricesKeeper, constants.Prices_DefaultGenesisState)
	perpetuals.InitGenesis(ks.Ctx, *ks.PerpetualsKeeper, constants.Perpetuals_DefaultGenesisState)
	ctx := ks.Ctx.WithTxBytes(constants.TestTxBytes)
	ctx = ctx.WithIsCheckTx(true)

	memClob.On("CreateOrderbook", constants.ClobPair_Btc).Return()
	// PerpetualMarketCreateEvents are emitted when initializing the genesis state, so we need to mock
	// the indexer event manager to expect these events.
	indexerEventManager.On("AddTxnEvent",
		ctx,
		indexerevents.SubtypePerpetualMarket,
		indexerevents.PerpetualMarketEventVersion,
		indexer_manager.GetBytes(
			indexerevents.NewPerpetualMarketCreateEvent(
				0,
				0,
				constants.Perpetuals_DefaultGenesisState.Perpetuals[0].Params.Ticker,
				constants.Perpetuals_DefaultGenesisState.Perpetuals[0].Params.MarketId,
				constants.ClobPair_Btc.Status,
				constants.ClobPair_Btc.QuantumConversionExponent,
				constants.Perpetuals_DefaultGenesisState.Perpetuals[0].Params.AtomicResolution,
				constants.ClobPair_Btc.SubticksPerTick,
				constants.ClobPair_Btc.StepBaseQuantums,
				constants.Perpetuals_DefaultGenesisState.Perpetuals[0].Params.LiquidityTier,
				constants.Perpetuals_DefaultGenesisState.Perpetuals[0].Params.MarketType,
				constants.Perpetuals_DefaultGenesisState.Perpetuals[0].Params.DefaultFundingPpm,
			),
		),
	).Once().Return()
	_, err := ks.ClobKeeper.CreatePerpetualClobPairAndMemStructs(
		ctx,
		constants.ClobPair_Btc.Id,
		clobtest.MustPerpetualId(constants.ClobPair_Btc),
		satypes.BaseQuantums(constants.ClobPair_Btc.StepBaseQuantums),
		constants.ClobPair_Btc.QuantumConversionExponent,
		constants.ClobPair_Btc.SubticksPerTick,
		constants.ClobPair_Btc.Status,
	)
	require.NoError(t, err)

	order := constants.Order_Carl_Num0_Id5_Clob0_Buy2BTC_Price50000
	msgPlaceOrder := &types.MsgPlaceOrder{Order: order}
	memClob.On("PlaceOrder", ctx, order).
		Return(order.GetBaseQuantums(), types.OrderStatus(0), constants.TestOffchainUpdates, nil)

	_, _, err = ks.ClobKeeper.PlaceShortTermOrder(ctx, msgPlaceOrder)
	require.NoError(t, err)
	indexerEventManager.AssertNumberOfCalls(t, "SendOffchainData", len(constants.TestOffchainMessages))
	indexerEventManager.AssertExpectations(t)
	memClob.AssertExpectations(t)
}

func TestPerformStatefulOrderValidation_PreExistingStatefulOrder(t *testing.T) {
	// Setup keeper state.
	memClob := &mocks.MemClob{}
	memClob.On("SetClobKeeper", mock.Anything).Return()
	indexerEventManager := &mocks.IndexerEventManager{}
	ks := keepertest.NewClobKeepersTestContext(t, memClob, &mocks.BankKeeper{}, indexerEventManager)

	ks.MarketMapKeeper.InitGenesis(ks.Ctx, constants.MarketMap_DefaultGenesisState)
	prices.InitGenesis(ks.Ctx, *ks.PricesKeeper, constants.Prices_DefaultGenesisState)
	perpetuals.InitGenesis(ks.Ctx, *ks.PerpetualsKeeper, constants.Perpetuals_DefaultGenesisState)

	memClob.On("CreateOrderbook", constants.ClobPair_Btc).Return()
	// PerpetualMarketCreateEvents are emitted when initializing the genesis state, so we need to mock
	// the indexer event manager to expect these events.
	indexerEventManager.On("AddTxnEvent",
		ks.Ctx,
		indexerevents.SubtypePerpetualMarket,
		indexerevents.PerpetualMarketEventVersion,
		indexer_manager.GetBytes(
			indexerevents.NewPerpetualMarketCreateEvent(
				0,
				0,
				constants.Perpetuals_DefaultGenesisState.Perpetuals[0].Params.Ticker,
				constants.Perpetuals_DefaultGenesisState.Perpetuals[0].Params.MarketId,
				constants.ClobPair_Btc.Status,
				constants.ClobPair_Btc.QuantumConversionExponent,
				constants.Perpetuals_DefaultGenesisState.Perpetuals[0].Params.AtomicResolution,
				constants.ClobPair_Btc.SubticksPerTick,
				constants.ClobPair_Btc.StepBaseQuantums,
				constants.Perpetuals_DefaultGenesisState.Perpetuals[0].Params.LiquidityTier,
				constants.Perpetuals_DefaultGenesisState.Perpetuals[0].Params.MarketType,
				constants.Perpetuals_DefaultGenesisState.Perpetuals[0].Params.DefaultFundingPpm,
			),
		),
	).Once().Return()
	_, err := ks.ClobKeeper.CreatePerpetualClobPairAndMemStructs(
		ks.Ctx,
		constants.ClobPair_Btc.Id,
		clobtest.MustPerpetualId(constants.ClobPair_Btc),
		satypes.BaseQuantums(constants.ClobPair_Btc.StepBaseQuantums),
		constants.ClobPair_Btc.QuantumConversionExponent,
		constants.ClobPair_Btc.SubticksPerTick,
		constants.ClobPair_Btc.Status,
	)
	require.NoError(t, err)
	ctx := ks.Ctx.WithBlockHeight(int64(100)).WithBlockTime(time.Unix(5, 0))
	ks.BlockTimeKeeper.SetPreviousBlockInfo(ctx, &blocktimetypes.BlockInfo{
		Height:    100,
		Timestamp: time.Unix(int64(5), 0),
	})
	order := constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15

	// Run the test if the preexisting order is not in state. Expected panic.
	require.Panicsf(
		t,
		func() {
			_ = ks.ClobKeeper.PerformStatefulOrderValidation(ctx, &order, 10, true)
		},
		fmt.Sprintf(
			"PerformStatefulOrderValidation: Expected pre-existing stateful order to exist in state "+
				"order: (%+v).",
			&order,
		),
	)

	// Run the test if the preexisting order is not in state. Expected no panic.
	err = ks.ClobKeeper.PerformStatefulOrderValidation(ctx, &order, 10, false)
	require.NoError(t, err)

	// Run the test if the preexisting order is in state. Expected no panic.
	ks.ClobKeeper.SetLongTermOrderPlacement(ctx, order, 10)
	err = ks.ClobKeeper.PerformStatefulOrderValidation(ctx, &order, 10, true)
	require.NoError(t, err)
}

func TestPerformStatefulOrderValidation(t *testing.T) {
	blockHeight := uint32(5)

	tests := map[string]struct {
		setupDeliverTxState func(ctx sdk.Context, k *keeper.Keeper)
		clobPairs           []types.ClobPair
		order               types.Order
		expectedErr         string
	}{
		"Succeeds with a GoodTilBlock of blockHeight": {
			order: types.Order{
				OrderId: types.OrderId{
					ClientId:     0,
					SubaccountId: constants.Alice_Num0,
					ClobPairId:   uint32(0),
				},
				Side:         types.Order_SIDE_BUY,
				Quantums:     600,
				Subticks:     78,
				GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: blockHeight},
			},
		},
		"Succeeds with a GoodTilBlock of blockHeight + ShortBlockWindow": {
			order: types.Order{
				OrderId: types.OrderId{
					ClientId:     0,
					SubaccountId: constants.Alice_Num0,
					ClobPairId:   uint32(0),
				},
				Side:         types.Order_SIDE_BUY,
				Quantums:     600,
				Subticks:     78,
				GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: blockHeight + types.ShortBlockWindow},
			},
		},
		"Fails with invalid ClobPairId": {
			order: types.Order{
				OrderId: types.OrderId{
					ClientId:     0,
					SubaccountId: constants.Alice_Num0,
					ClobPairId:   uint32(1),
				},
				Side:         types.Order_SIDE_BUY,
				Quantums:     600,
				Subticks:     78,
				GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: blockHeight + 5},
			},
			expectedErr: types.ErrInvalidClob.Error(),
		},
		"Fails if Subticks is not a multiple of SubticksPerTick": {
			order: types.Order{
				OrderId: types.OrderId{
					ClientId:     0,
					SubaccountId: constants.Alice_Num0,
					ClobPairId:   uint32(0),
				},
				Side:         types.Order_SIDE_BUY,
				Quantums:     600,
				Subticks:     77,
				GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: blockHeight + 5},
			},
			expectedErr: "must be a multiple of the ClobPair's SubticksPerTick",
		},
		"Fails if Quantums is not a multiple of StepBaseQuantums": {
			order: types.Order{
				OrderId: types.OrderId{
					ClientId:     0,
					SubaccountId: constants.Alice_Num0,
					ClobPairId:   uint32(0),
				},
				Side:         types.Order_SIDE_BUY,
				Quantums:     599,
				Subticks:     78,
				GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: blockHeight + 5},
			},
			expectedErr: "must be a multiple of the ClobPair's StepBaseQuantums",
		},
		"Still succeeds if Order router address is not found": {
			order: types.Order{
				OrderId: types.OrderId{
					ClientId:     0,
					SubaccountId: constants.Alice_Num0,
					ClobPairId:   uint32(0),
				},
				Side:               types.Order_SIDE_BUY,
				Quantums:           600,
				Subticks:           78,
				GoodTilOneof:       &types.Order_GoodTilBlock{GoodTilBlock: blockHeight + 5},
				OrderRouterAddress: constants.AliceAccAddress.String(),
			},
		},
		"Fails if GoodTilBlock is in the past": {
			order: types.Order{
				OrderId: types.OrderId{
					ClientId:     0,
					SubaccountId: constants.Alice_Num0,
					ClobPairId:   uint32(0),
				},
				Side:         types.Order_SIDE_BUY,
				Quantums:     600,
				Subticks:     78,
				GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: blockHeight - 1},
			},
			expectedErr: types.ErrHeightExceedsGoodTilBlock.Error(),
		},
		"Fails if GoodTilBlock Exceeds ShortBlockWindow": {
			order: types.Order{
				OrderId: types.OrderId{
					ClientId:     0,
					SubaccountId: constants.Alice_Num0,
					ClobPairId:   uint32(0),
				},
				Side:         types.Order_SIDE_BUY,
				Quantums:     600,
				Subticks:     78,
				GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: blockHeight + types.ShortBlockWindow + 1},
			},
			expectedErr: types.ErrGoodTilBlockExceedsShortBlockWindow.Error(),
		},
		"Stateful: Fails if GoodTilBlockTime is less than or equal to the block time of the previous block": {
			order: types.Order{
				OrderId: types.OrderId{
					ClientId:     0,
					SubaccountId: constants.Alice_Num0,
					OrderFlags:   types.OrderIdFlags_LongTerm,
					ClobPairId:   uint32(0),
				},
				Side:     types.Order_SIDE_BUY,
				Quantums: 600,
				Subticks: 78,
				GoodTilOneof: &types.Order_GoodTilBlockTime{
					GoodTilBlockTime: 4,
				},
			},
			expectedErr: types.ErrTimeExceedsGoodTilBlockTime.Error(),
		},
		"Stateful: Fails if GoodTilBlockTime Exceeds StatefulOrderTimeWindow": {
			order: types.Order{
				OrderId: types.OrderId{
					ClientId:     0,
					SubaccountId: constants.Alice_Num0,
					OrderFlags:   types.OrderIdFlags_LongTerm,
					ClobPairId:   uint32(0),
				},
				Side:     types.Order_SIDE_BUY,
				Quantums: 600,
				Subticks: 78,
				GoodTilOneof: &types.Order_GoodTilBlockTime{
					GoodTilBlockTime: lib.MustConvertIntegerToUint32(
						time.Unix(5, 0).Add(types.StatefulOrderTimeWindow).Unix() + 1,
					),
				},
			},
			expectedErr: types.ErrGoodTilBlockTimeExceedsStatefulOrderTimeWindow.Error(),
		},
		`Stateful: Returns error when order already exists in state`: {
			setupDeliverTxState: func(ctx sdk.Context, k *keeper.Keeper) {
				k.SetLongTermOrderPlacement(
					ctx,
					types.Order{
						OrderId: types.OrderId{
							SubaccountId: constants.Alice_Num0,
							ClientId:     0,
							OrderFlags:   types.OrderIdFlags_LongTerm,
						},
						Side:         types.Order_SIDE_BUY,
						Quantums:     600,
						Subticks:     78,
						GoodTilOneof: &types.Order_GoodTilBlockTime{GoodTilBlockTime: 20},
					},
					lib.MustConvertIntegerToUint32(ctx.BlockHeight()),
				)
			},
			order: types.Order{
				OrderId: types.OrderId{
					SubaccountId: constants.Alice_Num0,
					ClientId:     0,
					OrderFlags:   types.OrderIdFlags_LongTerm,
				},
				Side:         types.Order_SIDE_BUY,
				Quantums:     600,
				Subticks:     78,
				GoodTilOneof: &types.Order_GoodTilBlockTime{GoodTilBlockTime: 20},
			},
			expectedErr: types.ErrStatefulOrderAlreadyExists.Error(),
		},
		`Stateful: Returns error when order with same order id but higher priority exists in state`: {
			setupDeliverTxState: func(ctx sdk.Context, k *keeper.Keeper) {
				k.SetLongTermOrderPlacement(
					ctx,
					types.Order{
						OrderId: types.OrderId{
							SubaccountId: constants.Alice_Num0,
							ClientId:     0,
							OrderFlags:   types.OrderIdFlags_LongTerm,
						},
						Side:         types.Order_SIDE_BUY,
						Quantums:     600,
						Subticks:     78,
						GoodTilOneof: &types.Order_GoodTilBlockTime{GoodTilBlockTime: 15},
					},
					lib.MustConvertIntegerToUint32(ctx.BlockHeight()),
				)
			},
			order: types.Order{
				OrderId: types.OrderId{
					SubaccountId: constants.Alice_Num0,
					ClientId:     0,
					OrderFlags:   types.OrderIdFlags_LongTerm,
				},
				Side:         types.Order_SIDE_BUY,
				Quantums:     600,
				Subticks:     78,
				GoodTilOneof: &types.Order_GoodTilBlockTime{GoodTilBlockTime: 20},
			},
			expectedErr: types.ErrStatefulOrderAlreadyExists.Error(),
		},
		`Stateful: Returns error when order with same order id but lower priority exists in state`: {
			setupDeliverTxState: func(ctx sdk.Context, k *keeper.Keeper) {
				k.SetLongTermOrderPlacement(
					ctx,
					types.Order{
						OrderId: types.OrderId{
							SubaccountId: constants.Alice_Num0,
							ClientId:     0,
							OrderFlags:   types.OrderIdFlags_LongTerm,
						},
						Side:         types.Order_SIDE_BUY,
						Quantums:     600,
						Subticks:     78,
						GoodTilOneof: &types.Order_GoodTilBlockTime{GoodTilBlockTime: 20},
					},
					lib.MustConvertIntegerToUint32(ctx.BlockHeight()),
				)
			},
			order: types.Order{
				OrderId: types.OrderId{
					SubaccountId: constants.Alice_Num0,
					ClientId:     0,
					OrderFlags:   types.OrderIdFlags_LongTerm,
				},
				Side:         types.Order_SIDE_BUY,
				Quantums:     600,
				Subticks:     78,
				GoodTilOneof: &types.Order_GoodTilBlockTime{GoodTilBlockTime: 15},
			},
			expectedErr: types.ErrStatefulOrderAlreadyExists.Error(),
		},
		"Conditional: Fails if GoodTilBlockTime is in the past": {
			order: types.Order{
				OrderId: types.OrderId{
					ClientId:     0,
					SubaccountId: constants.Alice_Num0,
					OrderFlags:   types.OrderIdFlags_Conditional,
					ClobPairId:   uint32(0),
				},
				Side:     types.Order_SIDE_BUY,
				Quantums: 600,
				Subticks: 78,
				GoodTilOneof: &types.Order_GoodTilBlockTime{
					GoodTilBlockTime: 4,
				},
				ConditionType:                   types.Order_CONDITION_TYPE_TAKE_PROFIT,
				ConditionalOrderTriggerSubticks: uint64(100),
			},
			expectedErr: types.ErrTimeExceedsGoodTilBlockTime.Error(),
		},
		"Conditional: Fails if GoodTilBlockTime Exceeds StatefulOrderTimeWindow": {
			order: types.Order{
				OrderId: types.OrderId{
					ClientId:     0,
					SubaccountId: constants.Alice_Num0,
					OrderFlags:   types.OrderIdFlags_Conditional,
					ClobPairId:   uint32(0),
				},
				Side:     types.Order_SIDE_BUY,
				Quantums: 600,
				Subticks: 78,
				GoodTilOneof: &types.Order_GoodTilBlockTime{
					GoodTilBlockTime: lib.MustConvertIntegerToUint32(
						time.Unix(5, 0).Add(types.StatefulOrderTimeWindow).Unix() + 1,
					),
				},
				ConditionType:                   types.Order_CONDITION_TYPE_TAKE_PROFIT,
				ConditionalOrderTriggerSubticks: uint64(100),
			},
			expectedErr: types.ErrGoodTilBlockTimeExceedsStatefulOrderTimeWindow.Error(),
		},
		`Conditional: Returns error when order already exists in state`: {
			setupDeliverTxState: func(ctx sdk.Context, k *keeper.Keeper) {
				k.SetLongTermOrderPlacement(
					ctx,
					types.Order{
						OrderId: types.OrderId{
							SubaccountId: constants.Alice_Num0,
							ClientId:     0,
							OrderFlags:   types.OrderIdFlags_Conditional,
						},
						Side:                            types.Order_SIDE_BUY,
						Quantums:                        600,
						Subticks:                        78,
						GoodTilOneof:                    &types.Order_GoodTilBlockTime{GoodTilBlockTime: 20},
						ConditionType:                   types.Order_CONDITION_TYPE_TAKE_PROFIT,
						ConditionalOrderTriggerSubticks: uint64(100),
					},
					lib.MustConvertIntegerToUint32(ctx.BlockHeight()),
				)
			},
			order: types.Order{
				OrderId: types.OrderId{
					SubaccountId: constants.Alice_Num0,
					ClientId:     0,
					OrderFlags:   types.OrderIdFlags_Conditional,
				},
				Side:                            types.Order_SIDE_BUY,
				Quantums:                        600,
				Subticks:                        78,
				GoodTilOneof:                    &types.Order_GoodTilBlockTime{GoodTilBlockTime: 20},
				ConditionType:                   types.Order_CONDITION_TYPE_TAKE_PROFIT,
				ConditionalOrderTriggerSubticks: uint64(100),
			},
			expectedErr: types.ErrStatefulOrderAlreadyExists.Error(),
		},
		`Conditional: Returns error when order with same order id but lower priority exists in state`: {
			setupDeliverTxState: func(ctx sdk.Context, k *keeper.Keeper) {
				k.SetLongTermOrderPlacement(
					ctx,
					types.Order{
						OrderId: types.OrderId{
							SubaccountId: constants.Alice_Num0,
							ClientId:     0,
							OrderFlags:   types.OrderIdFlags_Conditional,
						},
						Side:                            types.Order_SIDE_BUY,
						Quantums:                        600,
						Subticks:                        78,
						GoodTilOneof:                    &types.Order_GoodTilBlockTime{GoodTilBlockTime: 20},
						ConditionType:                   types.Order_CONDITION_TYPE_TAKE_PROFIT,
						ConditionalOrderTriggerSubticks: uint64(100),
					},
					lib.MustConvertIntegerToUint32(ctx.BlockHeight()),
				)
			},
			order: types.Order{
				OrderId: types.OrderId{
					SubaccountId: constants.Alice_Num0,
					ClientId:     0,
					OrderFlags:   types.OrderIdFlags_Conditional,
				},
				Side:                            types.Order_SIDE_BUY,
				Quantums:                        600,
				Subticks:                        78,
				GoodTilOneof:                    &types.Order_GoodTilBlockTime{GoodTilBlockTime: 15},
				ConditionType:                   types.Order_CONDITION_TYPE_TAKE_PROFIT,
				ConditionalOrderTriggerSubticks: uint64(100),
			},
			expectedErr: types.ErrStatefulOrderAlreadyExists.Error(),
		},
		"Conditional: Fails if ConditionalOrderTriggerSubticks is not a multiple of clob pair subticks": {
			order: types.Order{
				OrderId: types.OrderId{
					ClientId:     0,
					SubaccountId: constants.Alice_Num0,
					OrderFlags:   types.OrderIdFlags_Conditional,
					ClobPairId:   uint32(0),
				},
				Side:     types.Order_SIDE_BUY,
				Quantums: 600,
				Subticks: 78,
				GoodTilOneof: &types.Order_GoodTilBlockTime{
					GoodTilBlockTime: lib.MustConvertIntegerToUint32(
						time.Unix(5, 0).Add(types.StatefulOrderTimeWindow).Unix(),
					),
				},
				ConditionType:                   types.Order_CONDITION_TYPE_TAKE_PROFIT,
				ConditionalOrderTriggerSubticks: uint64(101),
			},
			expectedErr: types.ErrInvalidPlaceOrder.Error(),
		},
		"Fails with long-term order and ClobPair_Status of INITIALIZING": {
			clobPairs: []types.ClobPair{
				{
					Metadata: &types.ClobPair_PerpetualClobMetadata{
						PerpetualClobMetadata: &types.PerpetualClobMetadata{
							PerpetualId: 0,
						},
					},
					Status:           types.ClobPair_STATUS_INITIALIZING,
					StepBaseQuantums: 10,
					SubticksPerTick:  10,
				},
			},
			order:       constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy100_Price10_GTBT15,
			expectedErr: "must not be stateful for clob pair with status",
		},
		"Fails with conditional order and ClobPair_Status of INITIALIZING": {
			clobPairs: []types.ClobPair{
				{
					Metadata: &types.ClobPair_PerpetualClobMetadata{
						PerpetualClobMetadata: &types.PerpetualClobMetadata{
							PerpetualId: 0,
						},
					},
					Status:           types.ClobPair_STATUS_INITIALIZING,
					StepBaseQuantums: 10,
					SubticksPerTick:  10,
				},
			},
			order:       constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_SL_50001,
			expectedErr: "must not be stateful for clob pair with status",
		},
		"Fails with short-term non-post-only order and ClobPair_Status of INITIALIZING": {
			clobPairs: []types.ClobPair{
				{
					Metadata: &types.ClobPair_PerpetualClobMetadata{
						PerpetualClobMetadata: &types.PerpetualClobMetadata{
							PerpetualId: 0,
						},
					},
					Status:           types.ClobPair_STATUS_INITIALIZING,
					StepBaseQuantums: 10,
					SubticksPerTick:  10,
				},
			},
			order:       constants.Order_Alice_Num0_Id0_Clob0_Buy10_Price10_GTB16,
			expectedErr: "must be post-only for clob pair with status",
		},
		"Fails with short-term post-only bid above oracle price and ClobPair_Status of INITIALIZING": {
			clobPairs: []types.ClobPair{
				{
					Metadata: &types.ClobPair_PerpetualClobMetadata{
						PerpetualClobMetadata: &types.PerpetualClobMetadata{
							PerpetualId: 0,
						},
					},
					Status:           types.ClobPair_STATUS_INITIALIZING,
					StepBaseQuantums: 10,
					SubticksPerTick:  1,
				},
			},
			order: types.Order{
				OrderId:      types.OrderId{SubaccountId: constants.Alice_Num0, ClientId: 1, ClobPairId: 0},
				Side:         types.Order_SIDE_BUY,
				Quantums:     20,
				Subticks:     3, // oracle price for btc is 2 subticks for default genesis
				GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 15},
				TimeInForce:  types.Order_TIME_IN_FORCE_POST_ONLY,
			},
			expectedErr: "must be less than or equal to oracle price subticks",
		},
		"Fails with short-term post-only ask below oracle price and ClobPair_Status of INITIALIZING": {
			clobPairs: []types.ClobPair{
				{
					Metadata: &types.ClobPair_PerpetualClobMetadata{
						PerpetualClobMetadata: &types.PerpetualClobMetadata{
							PerpetualId: 0,
						},
					},
					Status:           types.ClobPair_STATUS_INITIALIZING,
					StepBaseQuantums: 10,
					SubticksPerTick:  1,
				},
			},
			order: types.Order{
				OrderId:      types.OrderId{SubaccountId: constants.Alice_Num0, ClientId: 1, ClobPairId: 0},
				Side:         types.Order_SIDE_SELL,
				Quantums:     20,
				Subticks:     1, // oracle price for btc is 2 subticks for default genesis
				GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 15},
				TimeInForce:  types.Order_TIME_IN_FORCE_POST_ONLY,
			},
			expectedErr: "must be greater than or equal to oracle price subticks",
		},
		"Succeeds with short-term post-only bid below oracle price and ClobPair_Status of INITIALIZING": {
			clobPairs: []types.ClobPair{
				{
					Metadata: &types.ClobPair_PerpetualClobMetadata{
						PerpetualClobMetadata: &types.PerpetualClobMetadata{
							PerpetualId: 0,
						},
					},
					Status:           types.ClobPair_STATUS_INITIALIZING,
					StepBaseQuantums: 10,
					SubticksPerTick:  1,
				},
			},
			order: types.Order{
				OrderId:      types.OrderId{SubaccountId: constants.Alice_Num0, ClientId: 1, ClobPairId: 0},
				Side:         types.Order_SIDE_BUY,
				Quantums:     20,
				Subticks:     1, // oracle price for btc is 2 subticks for default genesis
				GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 15},
				TimeInForce:  types.Order_TIME_IN_FORCE_POST_ONLY,
			},
		},
		"Succeeds with short-term post-only ask above oracle price and ClobPair_Status of INITIALIZING": {
			clobPairs: []types.ClobPair{
				{
					Metadata: &types.ClobPair_PerpetualClobMetadata{
						PerpetualClobMetadata: &types.PerpetualClobMetadata{
							PerpetualId: 0,
						},
					},
					Status:           types.ClobPair_STATUS_INITIALIZING,
					StepBaseQuantums: 10,
					SubticksPerTick:  1,
				},
			},
			order: types.Order{
				OrderId:      types.OrderId{SubaccountId: constants.Alice_Num0, ClientId: 1, ClobPairId: 0},
				Side:         types.Order_SIDE_SELL,
				Quantums:     20,
				Subticks:     3, // oracle price for btc is 2 subticks for default genesis
				GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 15},
				TimeInForce:  types.Order_TIME_IN_FORCE_POST_ONLY,
			},
		},
		"Succeeds with short-term post-only bid equal to oracle price and ClobPair_Status of INITIALIZING": {
			clobPairs: []types.ClobPair{
				{
					Metadata: &types.ClobPair_PerpetualClobMetadata{
						PerpetualClobMetadata: &types.PerpetualClobMetadata{
							PerpetualId: 0,
						},
					},
					Status:           types.ClobPair_STATUS_INITIALIZING,
					StepBaseQuantums: 10,
					SubticksPerTick:  1,
				},
			},
			order: types.Order{
				OrderId:      types.OrderId{SubaccountId: constants.Alice_Num0, ClientId: 1, ClobPairId: 0},
				Side:         types.Order_SIDE_BUY,
				Quantums:     20,
				Subticks:     2, // oracle price for btc is 2 subticks for default genesis
				GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 15},
				TimeInForce:  types.Order_TIME_IN_FORCE_POST_ONLY,
			},
		},
		"Succeeds with short-term post-only ask equal to oracle price and ClobPair_Status of INITIALIZING": {
			clobPairs: []types.ClobPair{
				{
					Metadata: &types.ClobPair_PerpetualClobMetadata{
						PerpetualClobMetadata: &types.PerpetualClobMetadata{
							PerpetualId: 0,
						},
					},
					Status:           types.ClobPair_STATUS_INITIALIZING,
					StepBaseQuantums: 10,
					SubticksPerTick:  1,
				},
			},
			order: types.Order{
				OrderId:      types.OrderId{SubaccountId: constants.Alice_Num0, ClientId: 1, ClobPairId: 0},
				Side:         types.Order_SIDE_SELL,
				Quantums:     20,
				Subticks:     2, // oracle price for btc is 2 subticks for default genesis
				GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 15},
				TimeInForce:  types.Order_TIME_IN_FORCE_POST_ONLY,
			},
		},
		"Fails with short-term order and ClobPair_Status of FINAL_SETTLEMENT": {
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc_Final_Settlement,
			},
			order:       constants.Order_Alice_Num0_Id0_Clob0_Buy10_Price10_GTB16,
			expectedErr: "trading is disabled for clob pair",
		},
		"Fails with long-term order and ClobPair_Status of FINAL_SETTLEMENT": {
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc_Final_Settlement,
			},
			order:       constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy100_Price10_GTBT15,
			expectedErr: "trading is disabled for clob pair",
		},
		"Fails with conditional order and ClobPair_Status of FINAL_SETTLEMENT": {
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc_Final_Settlement,
			},
			order:       constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_SL_50001,
			expectedErr: "trading is disabled for clob pair",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).
				// Disable non-determinism checks since the tests update state via keeper directly.
				WithNonDeterminismChecksEnabled(false).
				WithGenesisDocFn(func() cmt.GenesisDoc {
					genesis := testapp.DefaultGenesis()
					clobPairs := []types.ClobPair{
						{
							Metadata: &types.ClobPair_PerpetualClobMetadata{
								PerpetualClobMetadata: &types.PerpetualClobMetadata{
									PerpetualId: 0,
								},
							},
							Status:           types.ClobPair_STATUS_ACTIVE,
							StepBaseQuantums: 12,
							SubticksPerTick:  39,
						},
					}
					if tc.clobPairs != nil {
						clobPairs = tc.clobPairs
					}
					testapp.UpdateGenesisDocWithAppStateForModule(&genesis, func(state *types.GenesisState) {
						state.ClobPairs = clobPairs
					})
					return genesis
				}).Build()

			ctx := tApp.AdvanceToBlock(
				// Stateful validation happens at blockHeight+1 for short term order placements.
				blockHeight-1,
				testapp.AdvanceToBlockOptions{BlockTime: time.Unix(5, 0)},
			)

			if tc.setupDeliverTxState != nil {
				tc.setupDeliverTxState(ctx.WithIsCheckTx(false), tApp.App.ClobKeeper)
			}

			resp := tApp.CheckTx(testapp.MustMakeCheckTxsWithClobMsg(
				ctx,
				tApp.App,
				*types.NewMsgPlaceOrder(tc.order))[0],
			)

			if tc.expectedErr != "" {
				require.Conditionf(t, resp.IsErr, "Expected CheckTx to error. Response: %+v", resp)
				require.Contains(t, resp.Log, tc.expectedErr)
			} else {
				require.Conditionf(t, resp.IsOK, "Expected CheckTx to succeed. Response: %+v", resp)
			}
		})
	}
}

func TestGetStatePosition_Success(t *testing.T) {
	tests := map[string]struct {
		// Subaccount state.
		subaccount *satypes.Subaccount
		// CLOB state.
		clob types.ClobPair

		// Parameters.
		subaccountId satypes.SubaccountId
		clobPairId   types.ClobPairId

		// Expectations.
		expectedPositionSize *big.Int
	}{
		`Can fetch the position size of a long position`: {
			subaccount: &constants.Dave_Num0_1BTC_Long_50000USD,

			subaccountId: constants.Dave_Num0,
			clobPairId:   types.ClobPairId(constants.ClobPair_Btc.Id),

			expectedPositionSize: constants.Dave_Num0_1BTC_Long_50000USD.PerpetualPositions[0].GetBigQuantums(),
		},
		`Can fetch the position size from multiple positions`: {
			subaccount: &constants.Dave_Num0_1BTC_Long_1ETH_Long_46000USD_Short,

			subaccountId: constants.Dave_Num0,
			clobPairId:   types.ClobPairId(constants.ClobPair_Eth.Id),

			expectedPositionSize: constants.Dave_Num0_1BTC_Long_1ETH_Long_46000USD_Short.PerpetualPositions[1].GetBigQuantums(),
		},
		`Can fetch the position size of a short position`: {
			subaccount: &constants.Carl_Num0_1BTC_Short,

			subaccountId: constants.Carl_Num0,
			clobPairId:   types.ClobPairId(constants.ClobPair_Btc.Id),

			expectedPositionSize: constants.Carl_Num0_1BTC_Short.PerpetualPositions[0].GetBigQuantums(),
		},
		`Fetching a non-existent subaccount returns 0`: {
			subaccountId: constants.Carl_Num0,
			clobPairId:   types.ClobPairId(constants.ClobPair_Btc.Id),

			expectedPositionSize: big.NewInt(0),
		},
		`Fetching a subaccount that doesn't have an open position corresponding to CLOB returns 0`: {
			subaccount: &constants.Dave_Num0_1BTC_Long_50000USD,

			subaccountId: constants.Dave_Num0,
			clobPairId:   types.ClobPairId(constants.ClobPair_Eth.Id),

			expectedPositionSize: big.NewInt(0),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup keeper state.
			memClob := memclob.NewMemClobPriceTimePriority(false)
			indexerEventManager := &mocks.IndexerEventManager{}
			ks := keepertest.NewClobKeepersTestContext(t, memClob, &mocks.BankKeeper{}, indexerEventManager)

			// Create subaccount if it's specified.
			if tc.subaccount != nil {
				ks.SubaccountsKeeper.SetSubaccount(ks.Ctx, *tc.subaccount)
			}

			ks.MarketMapKeeper.InitGenesis(ks.Ctx, constants.MarketMap_DefaultGenesisState)
			prices.InitGenesis(ks.Ctx, *ks.PricesKeeper, constants.Prices_DefaultGenesisState)
			perpetuals.InitGenesis(ks.Ctx, *ks.PerpetualsKeeper, constants.Perpetuals_DefaultGenesisState)

			// Create CLOB pairs.
			clobPairs := []types.ClobPair{constants.ClobPair_Btc, constants.ClobPair_Eth}
			for i, cp := range clobPairs {
				perpetualId := clobtest.MustPerpetualId(cp)
				indexerEventManager.On("AddTxnEvent",
					ks.Ctx,
					indexerevents.SubtypePerpetualMarket,
					indexerevents.PerpetualMarketEventVersion,
					indexer_manager.GetBytes(
						indexerevents.NewPerpetualMarketCreateEvent(
							perpetualId,
							uint32(i),
							constants.Perpetuals_DefaultGenesisState.Perpetuals[i].Params.Ticker,
							constants.Perpetuals_DefaultGenesisState.Perpetuals[i].Params.MarketId,
							cp.Status,
							cp.QuantumConversionExponent,
							constants.Perpetuals_DefaultGenesisState.Perpetuals[i].Params.AtomicResolution,
							cp.SubticksPerTick,
							cp.StepBaseQuantums,
							constants.Perpetuals_DefaultGenesisState.Perpetuals[i].Params.LiquidityTier,
							constants.Perpetuals_DefaultGenesisState.Perpetuals[i].Params.MarketType,
							constants.Perpetuals_DefaultGenesisState.Perpetuals[i].Params.DefaultFundingPpm,
						),
					),
				).Once().Return()
				_, err := ks.ClobKeeper.CreatePerpetualClobPairAndMemStructs(
					ks.Ctx,
					cp.Id,
					perpetualId,
					satypes.BaseQuantums(cp.StepBaseQuantums),
					cp.QuantumConversionExponent,
					cp.SubticksPerTick,
					cp.Status,
				)
				require.NoError(t, err)
			}

			// Run the test and verify expectations.
			positionSizeBig := ks.ClobKeeper.GetStatePosition(ks.Ctx, tc.subaccountId, tc.clobPairId)

			require.Equal(t, tc.expectedPositionSize, positionSizeBig)
		})
	}
}

func TestGetStatePosition_PanicsOnInvalidClob(t *testing.T) {
	// Setup keeper state.
	memClob := memclob.NewMemClobPriceTimePriority(false)
	ks := keepertest.NewClobKeepersTestContext(t, memClob, &mocks.BankKeeper{}, &mocks.IndexerEventManager{})

	ks.MarketMapKeeper.InitGenesis(ks.Ctx, constants.MarketMap_DefaultGenesisState)
	prices.InitGenesis(ks.Ctx, *ks.PricesKeeper, constants.Prices_DefaultGenesisState)
	perpetuals.InitGenesis(ks.Ctx, *ks.PerpetualsKeeper, constants.Perpetuals_DefaultGenesisState)

	// Run the test and verify expectations.
	clobPairId := types.ClobPairId(constants.ClobPair_Eth.Id)
	require.PanicsWithValue(
		t,
		fmt.Sprintf("GetStatePosition: CLOB pair %d not found", clobPairId),
		func() {
			ks.ClobKeeper.GetStatePosition(ks.Ctx, constants.Alice_Num0, clobPairId)
		},
	)
}

// TODO(DEC-1535): Uncomment this test when we can create spot CLOB pairs.
// func TestGetStatePosition_PanicsOnAssetClob(t *testing.T) {
// 	// Setup keeper state.
// 	MemClob := memclob.NewMemClobPriceTimePriority(false)
// 	ctx,
// 		clobKeeper,
// 		pricesKeeper,
// 		_,
// 		perpetualsKeeper,
// 		_,
// 		_,
// 		_ := keepertest.ClobKeepers(t, MemClob, &mocks.BankKeeper{}, &mocks.IndexerEventManager{})
// 	prices.InitGenesis(ctx, *pricesKeeper, constants.Prices_DefaultGenesisState)
// 	perpetuals.InitGenesis(ctx, *perpetualsKeeper, constants.Perpetuals_DefaultGenesisState)

// 	// Create CLOB pair.
// 	clobPair := constants.ClobPair_Asset
// 	clobKeeper.CreatePerpetualClobPairAndMemStructs(
// 		ctx,
// 		clobPair.Metadata.(*types.ClobPair_PerpetualClobMetadata),
// 		satypes.BaseQuantums(clobPair.StepBaseQuantums),
// 		satypes.BaseQuantums(clobPair.MinOrderBaseQuantums),
// 		clobPair.QuantumConversionExponent,
// 		clobPair.SubticksPerTick,
// 		clobPair.Status,
// 		clobPair.MakerFeePpm,
// 		clobPair.TakerFeePpm,
// 	)

// 	// Run the test and verify expectations.
// 	require.PanicsWithError(
// 		t,
// 		errorsmod.Wrap(
// 			types.ErrAssetOrdersNotImplemented,
// 			"GetStatePosition: Reduce-only orders for assets not implemented",
// 		).Error(),
// 		func() {
// 			clobKeeper.GetStatePosition(ctx, constants.Alice_Num0, clobPair.Id)
// 		},
// 	)
// }

func TestInitStatefulOrders(t *testing.T) {
	tests := map[string]struct {
		// CLOB module return values.
		statefulOrdersInState       []types.Order
		isConditionalOrderTriggered map[types.OrderId]bool
		orderPlacementErrors        map[types.OrderId]error
	}{
		`Can initialize 0 Long-Term orders or triggered conditional orders in the memclob with no errors`: {
			statefulOrdersInState:       []types.Order{},
			orderPlacementErrors:        map[types.OrderId]error{},
			isConditionalOrderTriggered: map[types.OrderId]bool{},
		},
		`Can initialize 0 Long-Term orders or triggered conditional orders in the memclob with no errors
			and does not place untriggered conditional orders`: {
			statefulOrdersInState: []types.Order{
				constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_TP_50005,
			},
			orderPlacementErrors:        map[types.OrderId]error{},
			isConditionalOrderTriggered: map[types.OrderId]bool{},
		},
		`Can initialize one Long-Term order in the memclob with no errors`: {
			statefulOrdersInState: []types.Order{
				constants.LongTermOrder_Alice_Num0_Id2_Clob0_Sell65_Price10_GTBT25,
			},
			orderPlacementErrors:        map[types.OrderId]error{},
			isConditionalOrderTriggered: map[types.OrderId]bool{},
		},
		`Can initialize one triggered conditional order in the memclob with no errors`: {
			statefulOrdersInState: []types.Order{
				constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_TP_50005,
			},
			orderPlacementErrors: map[types.OrderId]error{},
			isConditionalOrderTriggered: map[types.OrderId]bool{
				constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_TP_50005.OrderId: true,
			},
		},
		`Can initialize multiple Long-Term and triggered conditional orders in the memclob
			with no errors`: {
			statefulOrdersInState: []types.Order{
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
				constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_49995,
				constants.LongTermOrder_Alice_Num0_Id1_Clob0_Sell20_Price10_GTBT10,
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_TakeProfit10,
				constants.LongTermOrder_Alice_Num0_Id2_Clob0_Sell65_Price10_GTBT25,
			},
			orderPlacementErrors: map[types.OrderId]error{},
			isConditionalOrderTriggered: map[types.OrderId]bool{
				constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_49995.OrderId: true,
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_TakeProfit10.OrderId:  true,
			},
		},
		`Can initialize multiple Long-Term and triggered conditional orders in the memclob with errors`: {
			statefulOrdersInState: []types.Order{
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
				constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_49995,
				constants.LongTermOrder_Alice_Num0_Id1_Clob0_Sell20_Price10_GTBT10,
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_TakeProfit10,
				constants.LongTermOrder_Alice_Num0_Id2_Clob0_Sell65_Price10_GTBT25,
			},
			orderPlacementErrors: map[types.OrderId]error{
				constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_49995.
					OrderId: types.ErrInvalidStatefulOrderGoodTilBlockTime,
				constants.LongTermOrder_Alice_Num0_Id1_Clob0_Sell20_Price10_GTBT10.
					OrderId: types.ErrTimeExceedsGoodTilBlockTime,
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_TakeProfit10.
					OrderId: types.ErrStatefulOrderCollateralizationCheckFailed,
			},
			isConditionalOrderTriggered: map[types.OrderId]bool{
				constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_49995.OrderId: true,
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_TakeProfit10.OrderId:  true,
			},
		},
		`Can initialize multiple Long-Term and triggered conditional orders in the memclob where
			each order throws an error`: {
			statefulOrdersInState: []types.Order{
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
				constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_49995,
				constants.LongTermOrder_Alice_Num0_Id1_Clob0_Sell20_Price10_GTBT10,
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_TakeProfit10,
				constants.LongTermOrder_Alice_Num0_Id2_Clob0_Sell65_Price10_GTBT25,
			},
			orderPlacementErrors: map[types.OrderId]error{
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.
					OrderId: types.ErrStatefulOrderCollateralizationCheckFailed,
				constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_49995.
					OrderId: types.ErrInvalidStatefulOrderGoodTilBlockTime,
				constants.LongTermOrder_Alice_Num0_Id1_Clob0_Sell20_Price10_GTBT10.
					OrderId: types.ErrTimeExceedsGoodTilBlockTime,
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_TakeProfit10.
					OrderId: types.ErrStatefulOrderCollateralizationCheckFailed,
				constants.LongTermOrder_Alice_Num0_Id2_Clob0_Sell65_Price10_GTBT25.
					OrderId: types.ErrStatefulOrderCollateralizationCheckFailed,
			},
			isConditionalOrderTriggered: map[types.OrderId]bool{
				constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_49995.OrderId: true,
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_TakeProfit10.OrderId:  true,
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup state.
			memClob := &mocks.MemClob{}
			memClob.On("SetClobKeeper", mock.Anything).Return()

			indexerEventManager := &mocks.IndexerEventManager{}

			ks := keepertest.NewClobKeepersTestContext(t, memClob, &mocks.BankKeeper{}, indexerEventManager)

			ks.MarketMapKeeper.InitGenesis(ks.Ctx, constants.MarketMap_DefaultGenesisState)
			prices.InitGenesis(ks.Ctx, *ks.PricesKeeper, constants.Prices_DefaultGenesisState)
			perpetuals.InitGenesis(ks.Ctx, *ks.PerpetualsKeeper, constants.Perpetuals_DefaultGenesisState)

			// Create CLOB pair.
			memClob.On("CreateOrderbook", constants.ClobPair_Btc).Return()
			indexerEventManager.On("AddTxnEvent",
				ks.Ctx,
				indexerevents.SubtypePerpetualMarket,
				indexerevents.PerpetualMarketEventVersion,
				indexer_manager.GetBytes(
					indexerevents.NewPerpetualMarketCreateEvent(
						0,
						0,
						constants.Perpetuals_DefaultGenesisState.Perpetuals[0].Params.Ticker,
						constants.Perpetuals_DefaultGenesisState.Perpetuals[0].Params.MarketId,
						constants.ClobPair_Btc.Status,
						constants.ClobPair_Btc.QuantumConversionExponent,
						constants.Perpetuals_DefaultGenesisState.Perpetuals[0].Params.AtomicResolution,
						constants.ClobPair_Btc.SubticksPerTick,
						constants.ClobPair_Btc.StepBaseQuantums,
						constants.Perpetuals_DefaultGenesisState.Perpetuals[0].Params.LiquidityTier,
						constants.Perpetuals_DefaultGenesisState.Perpetuals[0].Params.MarketType,
						constants.Perpetuals_DefaultGenesisState.Perpetuals[0].Params.DefaultFundingPpm,
					),
				),
			).Once().Return()
			_, err := ks.ClobKeeper.CreatePerpetualClobPairAndMemStructs(
				ks.Ctx,
				constants.ClobPair_Btc.Id,
				clobtest.MustPerpetualId(constants.ClobPair_Btc),
				satypes.BaseQuantums(constants.ClobPair_Btc.StepBaseQuantums),
				constants.ClobPair_Btc.QuantumConversionExponent,
				constants.ClobPair_Btc.SubticksPerTick,
				constants.ClobPair_Btc.Status,
			)
			require.NoError(t, err)

			// Create each stateful order placement in state and properly mock the MemClob call.
			expectedPlacedOrders := make([]types.Order, 0)
			for i, order := range tc.statefulOrdersInState {
				require.True(t, order.IsStatefulOrder())

				// Write the stateful order placement to state.
				ks.ClobKeeper.SetLongTermOrderPlacement(ks.Ctx, order, uint32(i))
				// Clear the count since we expect InitStatefulOrders to initialize it.
				ks.ClobKeeper.SetStatefulOrderCount(ks.Ctx, order.OrderId.SubaccountId, 0)

				// No more state or memclob updates are required if this is an untriggered
				// conditional order.
				if order.IsConditionalOrder() && !tc.isConditionalOrderTriggered[order.OrderId] {
					require.NotContains(t, tc.orderPlacementErrors, order.OrderId)
					continue
				}

				// If it's a triggered conditional order, ensure it's triggered in state.
				if order.IsConditionalOrder() && tc.isConditionalOrderTriggered[order.OrderId] {
					ks.ClobKeeper.MustTriggerConditionalOrder(ks.Ctx, order.OrderId)
				}

				orderPlacementErr := tc.orderPlacementErrors[order.OrderId]
				memClob.On("PlaceOrder", mock.Anything, order).Return(
					satypes.BaseQuantums(0),
					types.Success,
					constants.TestOffchainUpdates,
					orderPlacementErr,
				).Once()

				for _, message := range constants.TestOffchainMessages {
					indexerEventManager.On("SendOffchainData", message).Return().Once()
				}

				expectedPlacedOrders = append(expectedPlacedOrders, order)
			}

			// Run the test and verify expectations.
			ks.ClobKeeper.InitStatefulOrders(ks.Ctx)
			indexerEventManager.AssertExpectations(t)
			indexerEventManager.AssertNumberOfCalls(
				t,
				"SendOffchainData",
				len(constants.TestOffchainMessages)*len(expectedPlacedOrders),
			)
			memClob.AssertExpectations(t)
		})
	}
}

func TestPlaceStatefulOrdersFromLastBlock(t *testing.T) {
	tests := map[string]struct {
		orders []types.Order

		expectedOrderPlacementCalls []types.Order
	}{
		"empty stateful orders": {
			orders: []types.Order{},

			expectedOrderPlacementCalls: []types.Order{},
		},
		"places stateful orders from last block": {
			orders: []types.Order{
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
				constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10,
				constants.LongTermOrder_Carl_Num0_Id0_Clob0_WithOrderRouterAddress,
			},
			expectedOrderPlacementCalls: []types.Order{
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
				constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10,
				constants.LongTermOrder_Carl_Num0_Id0_Clob0_WithOrderRouterAddress,
			},
		},
		"does not place orders with GTBT equal to block time": {
			orders: []types.Order{
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT5,
			},
			expectedOrderPlacementCalls: []types.Order{},
		},
		"does not place orders with GTBT less than block time": {
			orders: []types.Order{
				{
					OrderId: types.OrderId{
						SubaccountId: constants.Alice_Num0,
						ClientId:     0,
						OrderFlags:   types.OrderIdFlags_LongTerm,
						ClobPairId:   0,
					},
					Side:         types.Order_SIDE_BUY,
					Quantums:     5,
					Subticks:     10,
					GoodTilOneof: &types.Order_GoodTilBlockTime{GoodTilBlockTime: 2},
				},
			},
			expectedOrderPlacementCalls: []types.Order{},
		},
		"does not place orders with invalid clob pair id": {
			orders: []types.Order{
				{
					OrderId:      constants.InvalidClobPairId_Long_Term_Order,
					Side:         types.Order_SIDE_BUY,
					Quantums:     5,
					Subticks:     10,
					GoodTilOneof: &types.Order_GoodTilBlockTime{GoodTilBlockTime: 2},
				},
			},
			expectedOrderPlacementCalls: []types.Order{},
		},
		"does not place orders further than StatefulOrderTimeWindow in the future": {
			orders: []types.Order{
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
				{
					OrderId: types.OrderId{
						SubaccountId: constants.Alice_Num0,
						ClientId:     1,
						OrderFlags:   types.OrderIdFlags_LongTerm,
						ClobPairId:   0,
					},
					Side:     types.Order_SIDE_BUY,
					Quantums: 5,
					Subticks: 10,
					GoodTilOneof: &types.Order_GoodTilBlockTime{
						GoodTilBlockTime: lib.MustConvertIntegerToUint32(
							5 + int64(types.StatefulOrderTimeWindow.Seconds()+1),
						),
					},
				},
			},
			expectedOrderPlacementCalls: []types.Order{
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup state.
			memClob := &mocks.MemClob{}

			memClob.On("SetClobKeeper", mock.Anything).Return()

			ks := keepertest.NewClobKeepersTestContext(
				t,
				memClob,
				&mocks.BankKeeper{},
				indexer_manager.NewIndexerEventManagerNoop(),
			)

			ks.MarketMapKeeper.InitGenesis(ks.Ctx, constants.MarketMap_DefaultGenesisState)
			prices.InitGenesis(ks.Ctx, *ks.PricesKeeper, constants.Prices_DefaultGenesisState)
			perpetuals.InitGenesis(ks.Ctx, *ks.PerpetualsKeeper, constants.Perpetuals_DefaultGenesisState)
			err := ks.PricesKeeper.RevShareKeeper.SetOrderRouterRevShare(
				ks.Ctx,
				constants.AliceAccAddress.String(),
				100_000,
			)
			require.NoError(t, err)

			ctx := ks.Ctx.WithBlockHeight(int64(100)).WithBlockTime(time.Unix(5, 0))
			ctx = ctx.WithIsCheckTx(true)
			ks.BlockTimeKeeper.SetPreviousBlockInfo(ctx, &blocktimetypes.BlockInfo{
				Height:    100,
				Timestamp: time.Unix(int64(5), 0),
			})

			// Create CLOB pair.
			memClob.On("CreateOrderbook", constants.ClobPair_Btc).Return()
			_, err = ks.ClobKeeper.CreatePerpetualClobPairAndMemStructs(
				ctx,
				constants.ClobPair_Btc.Id,
				clobtest.MustPerpetualId(constants.ClobPair_Btc),
				satypes.BaseQuantums(constants.ClobPair_Btc.StepBaseQuantums),
				constants.ClobPair_Btc.QuantumConversionExponent,
				constants.ClobPair_Btc.SubticksPerTick,
				constants.ClobPair_Btc.Status,
			)
			require.NoError(t, err)

			// Create each stateful order placement in state
			for i, order := range tc.orders {
				require.True(t, order.IsStatefulOrder())

				ks.ClobKeeper.SetLongTermOrderPlacement(ctx.WithIsCheckTx(false), order, uint32(i))
			}

			// Assert expected order placement memclob calls.
			for _, order := range tc.expectedOrderPlacementCalls {
				memClob.On("PlaceOrder", mock.Anything, order).Return(
					satypes.BaseQuantums(0),
					types.Success,
					constants.TestOffchainUpdates,
					nil,
				).Once()
			}

			// Run the test and verify expectations.
			offchainUpdates := types.NewOffchainUpdates()
			orderIds := make([]types.OrderId, 0)
			for _, order := range tc.orders {
				orderIds = append(orderIds, order.OrderId)
			}
			ks.ClobKeeper.PlaceStatefulOrdersFromLastBlock(ctx, orderIds, offchainUpdates, true)
			ks.ClobKeeper.PlaceStatefulOrdersFromLastBlock(ctx, orderIds, offchainUpdates, false)

			// PlaceStatefulOrdersFromLastBlock utilizes the memclob's PlaceOrder flow, but we
			// do not want to emit PlaceMessages in offchain events for stateful orders. This assertion
			// verifies that we call `ClearPlaceMessages()` on the offchain updates before returning.
			require.Equal(t, 0, memclobtest.MessageCountOfType(offchainUpdates, types.PlaceMessageType))

			// Verify that all removed orders have an associated off-chain update.
			orderMap := make(map[types.OrderId]bool)
			for _, order := range tc.orders {
				orderMap[order.OrderId] = true
			}

			removedOrders := lib.FilterSlice(tc.expectedOrderPlacementCalls, func(order types.Order) bool {
				return !orderMap[order.OrderId]
			})

			for _, order := range removedOrders {
				require.True(
					t,
					memclobtest.HasMessage(offchainUpdates, order.OrderId, types.RemoveMessageType),
				)
			}

			memClob.AssertExpectations(t)
		})
	}
}

func TestPlaceStatefulOrdersFromLastBlock_PostOnly(t *testing.T) {
	tests := map[string]struct {
		orders         []types.Order
		postOnlyFilter bool

		expectedOrderPlacementCalls []types.Order
	}{
		"places PO stateful orders from last block when postOnlyFilter = true": {
			postOnlyFilter: true,
			orders: []types.Order{
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT5_PO,
			},
			expectedOrderPlacementCalls: []types.Order{
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT5_PO,
			},
		},
		"does not places non-PO stateful orders when postOnlyFilter = true": {
			postOnlyFilter: true,
			orders: []types.Order{
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT5,
			},
			expectedOrderPlacementCalls: []types.Order{},
		},
		"does not places PO stateful orders from last block, when postOnlyFilter = false": {
			postOnlyFilter: false,
			orders: []types.Order{
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT5_PO,
			},
			expectedOrderPlacementCalls: []types.Order{},
		},
		"places non-PO stateful orders from last block when postOnlyFilter = false": {
			postOnlyFilter: false,
			orders: []types.Order{
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT5,
			},
			expectedOrderPlacementCalls: []types.Order{
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT5,
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup state.
			memClob := &mocks.MemClob{}

			memClob.On("SetClobKeeper", mock.Anything).Return()

			ks := keepertest.NewClobKeepersTestContext(
				t,
				memClob,
				&mocks.BankKeeper{},
				indexer_manager.NewIndexerEventManagerNoop(),
			)

			ks.MarketMapKeeper.InitGenesis(ks.Ctx, constants.MarketMap_DefaultGenesisState)
			prices.InitGenesis(ks.Ctx, *ks.PricesKeeper, constants.Prices_DefaultGenesisState)
			perpetuals.InitGenesis(ks.Ctx, *ks.PerpetualsKeeper, constants.Perpetuals_DefaultGenesisState)

			ctx := ks.Ctx.WithBlockHeight(int64(100)).WithBlockTime(time.Unix(5, 0))
			ctx = ctx.WithIsCheckTx(true)
			ks.BlockTimeKeeper.SetPreviousBlockInfo(ctx, &blocktimetypes.BlockInfo{
				Height:    100,
				Timestamp: time.Unix(int64(2), 0),
			})

			// Create CLOB pair.
			memClob.On("CreateOrderbook", constants.ClobPair_Btc).Return()
			_, err := ks.ClobKeeper.CreatePerpetualClobPairAndMemStructs(
				ctx,
				constants.ClobPair_Btc.Id,
				clobtest.MustPerpetualId(constants.ClobPair_Btc),
				satypes.BaseQuantums(constants.ClobPair_Btc.StepBaseQuantums),
				constants.ClobPair_Btc.QuantumConversionExponent,
				constants.ClobPair_Btc.SubticksPerTick,
				constants.ClobPair_Btc.Status,
			)
			require.NoError(t, err)

			// Create each stateful order placement in state
			for i, order := range tc.orders {
				require.True(t, order.IsStatefulOrder())

				ks.ClobKeeper.SetLongTermOrderPlacement(ctx.WithIsCheckTx(false), order, uint32(i))
			}

			// Assert expected order placement memclob calls.
			for _, order := range tc.expectedOrderPlacementCalls {
				memClob.On("PlaceOrder", mock.Anything, order).Return(
					satypes.BaseQuantums(0),
					types.Success,
					constants.TestOffchainUpdates,
					nil,
				).Once()
			}

			// Run the test and verify expectations.
			offchainUpdates := types.NewOffchainUpdates()
			orderIds := make([]types.OrderId, 0)
			for _, order := range tc.orders {
				orderIds = append(orderIds, order.OrderId)
			}
			ks.ClobKeeper.PlaceStatefulOrdersFromLastBlock(ctx, orderIds, offchainUpdates, tc.postOnlyFilter)

			// PlaceStatefulOrdersFromLastBlock utilizes the memclob's PlaceOrder flow, but we
			// do not want to emit PlaceMessages in offchain events for stateful orders. This assertion
			// verifies that we call `ClearPlaceMessages()` on the offchain updates before returning.
			require.Equal(t, 0, memclobtest.MessageCountOfType(offchainUpdates, types.PlaceMessageType))

			// Verify that all removed orders have an associated off-chain update.
			orderMap := make(map[types.OrderId]bool)
			for _, order := range tc.orders {
				orderMap[order.OrderId] = true
			}

			removedOrders := lib.FilterSlice(tc.expectedOrderPlacementCalls, func(order types.Order) bool {
				return !orderMap[order.OrderId]
			})

			for _, order := range removedOrders {
				require.True(
					t,
					memclobtest.HasMessage(offchainUpdates, order.OrderId, types.RemoveMessageType),
				)
			}

			memClob.AssertExpectations(t)
		})
	}
}

func TestPlaceConditionalOrdersTriggeredInLastBlock(t *testing.T) {
	tests := map[string]struct {
		triggeredOrders             []types.Order
		untriggeredOrders           []types.Order
		expectedOrderPlacementCalls []types.Order
		expectedPanic               string
	}{
		"empty conditional orders": {
			triggeredOrders:             []types.Order{},
			expectedOrderPlacementCalls: []types.Order{},
		},
		"places conditional orders triggered in last block": {
			triggeredOrders: []types.Order{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20,
				constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTBT15_StopLoss15,
			},
			expectedOrderPlacementCalls: []types.Order{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20,
				constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTBT15_StopLoss15,
			},
		},
		"does not place stateful order": {
			triggeredOrders: []types.Order{
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy100_Price10_GTBT15,
			},
			expectedPanic: fmt.Sprintf(
				"MustBeConditionalOrder: called with non-conditional order ID (%+v)",
				&constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy100_Price10_GTBT15.OrderId,
			),
		},
		"does not place conditional order if not in triggered state": {
			untriggeredOrders: []types.Order{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20,
			},
			expectedPanic: fmt.Sprintf(
				"PlaceConditionalOrdersTriggeredInLastBlock: Order with OrderId %+v is not in triggered state",
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20.OrderId,
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup state.
			memClob := &mocks.MemClob{}

			memClob.On("SetClobKeeper", mock.Anything).Return()

			ks := keepertest.NewClobKeepersTestContext(
				t,
				memClob,
				&mocks.BankKeeper{},
				indexer_manager.NewIndexerEventManagerNoop(),
			)
			ks.MarketMapKeeper.InitGenesis(ks.Ctx, constants.MarketMap_DefaultGenesisState)
			prices.InitGenesis(ks.Ctx, *ks.PricesKeeper, constants.Prices_DefaultGenesisState)
			perpetuals.InitGenesis(ks.Ctx, *ks.PerpetualsKeeper, constants.Perpetuals_DefaultGenesisState)

			ctx := ks.Ctx.WithBlockHeight(int64(100)).WithBlockTime(time.Unix(5, 0))
			ctx = ctx.WithIsCheckTx(true)
			ks.BlockTimeKeeper.SetPreviousBlockInfo(ctx, &blocktimetypes.BlockInfo{
				Height:    100,
				Timestamp: time.Unix(int64(5), 0),
			})

			// Create CLOB pair.
			memClob.On("CreateOrderbook", constants.ClobPair_Btc).Return()
			_, err := ks.ClobKeeper.CreatePerpetualClobPairAndMemStructs(
				ctx,
				constants.ClobPair_Btc.Id,
				clobtest.MustPerpetualId(constants.ClobPair_Btc),
				satypes.BaseQuantums(constants.ClobPair_Btc.StepBaseQuantums),
				constants.ClobPair_Btc.QuantumConversionExponent,
				constants.ClobPair_Btc.SubticksPerTick,
				constants.ClobPair_Btc.Status,
			)
			require.NoError(t, err)

			// Write to triggered orders state
			for _, order := range tc.triggeredOrders {
				longTermOrderPlacement := types.LongTermOrderPlacement{
					Order: order,
				}
				longTermOrderPlacementBytes := ks.Cdc.MustMarshal(&longTermOrderPlacement)

				store := ks.ClobKeeper.GetTriggeredConditionalOrderPlacementStore(ctx)

				orderKey := order.OrderId.ToStateKey()
				store.Set(orderKey, longTermOrderPlacementBytes)
			}

			// Write to untriggered orders state
			for _, order := range tc.untriggeredOrders {
				longTermOrderPlacement := types.LongTermOrderPlacement{
					Order: order,
				}
				longTermOrderPlacementBytes := ks.Cdc.MustMarshal(&longTermOrderPlacement)

				store := ks.ClobKeeper.GetUntriggeredConditionalOrderPlacementStore(ctx)

				orderKey := order.OrderId.ToStateKey()
				store.Set(orderKey, longTermOrderPlacementBytes)
			}

			// Assert expected order placement memclob calls.
			for _, order := range tc.expectedOrderPlacementCalls {
				memClob.On("PlaceOrder", mock.Anything, order).Return(
					satypes.BaseQuantums(0),
					types.Success,
					constants.TestOffchainUpdates,
					nil,
				).Once()
			}

			// Run the test and verify expectations.
			offchainUpdates := types.NewOffchainUpdates()
			orderIds := make([]types.OrderId, 0)
			for _, order := range tc.triggeredOrders {
				orderIds = append(orderIds, order.OrderId)
			}
			for _, order := range tc.untriggeredOrders {
				orderIds = append(orderIds, order.OrderId)
			}

			if tc.expectedPanic != "" {
				require.PanicsWithValue(
					t,
					tc.expectedPanic,
					func() {
						ks.ClobKeeper.PlaceConditionalOrdersTriggeredInLastBlock(ctx, orderIds, offchainUpdates, true)
						ks.ClobKeeper.PlaceConditionalOrdersTriggeredInLastBlock(ctx, orderIds, offchainUpdates, false)
					},
				)
				return
			}

			ks.ClobKeeper.PlaceConditionalOrdersTriggeredInLastBlock(ctx, orderIds, offchainUpdates, true)
			ks.ClobKeeper.PlaceConditionalOrdersTriggeredInLastBlock(ctx, orderIds, offchainUpdates, false)

			// PlaceStatefulOrdersFromLastBlock utilizes the memclob's PlaceOrder flow, but we
			// do not want to emit PlaceMessages in offchain events for stateful orders. This assertion
			// verifies that we call `ClearPlaceMessages()` on the offchain updates before returning.
			require.Equal(t, 0, memclobtest.MessageCountOfType(offchainUpdates, types.PlaceMessageType))

			// Verify that all removed orders have an associated off-chain update.
			orderMap := make(map[types.OrderId]bool)
			for _, order := range tc.triggeredOrders {
				orderMap[order.OrderId] = true
			}
			for _, order := range tc.untriggeredOrders {
				orderMap[order.OrderId] = true
			}

			removedOrders := lib.FilterSlice(tc.expectedOrderPlacementCalls, func(order types.Order) bool {
				return !orderMap[order.OrderId]
			})

			for _, order := range removedOrders {
				require.True(
					t,
					memclobtest.HasMessage(offchainUpdates, order.OrderId, types.RemoveMessageType),
				)
			}

			memClob.AssertExpectations(t)
		})
	}
}
