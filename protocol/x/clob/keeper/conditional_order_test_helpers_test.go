package keeper_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	indexer_manager "github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/memclob"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Shared fixtures for the conditional-order trigger tests.

const (
	condTestTriggerFarBelow = uint64(1)
	condTestGTBT            = uint32(1_700_000_000)
	condTestOwner           = "dydx199tqg4wdlnu4qjlxchpd7seg454937hjrknju4"
)

// Short aliases for the verbose order side / condition-type enum constants, used by the
// conditional-order index tests to keep table-driven order construction within line limits.
const (
	condSideBuy  = clobtypes.Order_SIDE_BUY
	condSideSell = clobtypes.Order_SIDE_SELL
	condTypeTP   = clobtypes.Order_CONDITION_TYPE_TAKE_PROFIT
	condTypeSL   = clobtypes.Order_CONDITION_TYPE_STOP_LOSS
)

// newCondTestKeepers builds a test ClobKeeper with the BTC clob pair created.
func newCondTestKeepers(t *testing.T) keepertest.ClobKeepersTestContext {
	t.Helper()

	bankKeeper := &mocks.BankKeeper{}
	bankKeeper.On("GetBalance", mock.Anything, mock.Anything, constants.Usdc.Denom).
		Return(sdk.NewCoin(constants.Usdc.Denom, sdkmath.ZeroInt())).
		Maybe()

	ks := keepertest.NewClobKeepersTestContext(
		t,
		memclob.NewMemClobPriceTimePriority(false),
		bankKeeper,
		indexer_manager.NewIndexerEventManagerNoop(),
	)
	ctx := ks.Ctx

	require.NoError(t, keepertest.CreateUsdcAsset(ctx, ks.AssetsKeeper))
	keepertest.CreateTestMarkets(t, ctx, ks.PricesKeeper)
	keepertest.CreateTestLiquidityTiers(t, ctx, ks.PerpetualsKeeper)
	_, err := ks.PerpetualsKeeper.CreatePerpetual(
		ctx,
		constants.BtcUsd_20PercentInitial_10PercentMaintenance.Params.Id,
		constants.BtcUsd_20PercentInitial_10PercentMaintenance.Params.Ticker,
		constants.BtcUsd_20PercentInitial_10PercentMaintenance.Params.MarketId,
		constants.BtcUsd_20PercentInitial_10PercentMaintenance.Params.AtomicResolution,
		constants.BtcUsd_20PercentInitial_10PercentMaintenance.Params.DefaultFundingPpm,
		constants.BtcUsd_20PercentInitial_10PercentMaintenance.Params.LiquidityTier,
		constants.BtcUsd_20PercentInitial_10PercentMaintenance.Params.MarketType,
	)
	require.NoError(t, err)

	_, err = ks.ClobKeeper.CreatePerpetualClobPairAndMemStructs(
		ctx,
		constants.ClobPair_Btc.Id,
		constants.ClobPair_Btc.MustGetPerpetualId(),
		satypes.BaseQuantums(constants.ClobPair_Btc.StepBaseQuantums),
		constants.ClobPair_Btc.QuantumConversionExponent,
		constants.ClobPair_Btc.SubticksPerTick,
		constants.ClobPair_Btc.Status,
	)
	require.NoError(t, err)

	return ks
}

// newCondTestOrder builds a take-profit buy conditional order with a unique client id.
func newCondTestOrder(i int, triggerSubticks uint64) clobtypes.Order {
	return clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: satypes.SubaccountId{
				Owner:  condTestOwner,
				Number: uint32(i / 10),
			},
			ClientId:   uint32(i),
			OrderFlags: clobtypes.OrderIdFlags_Conditional,
			ClobPairId: constants.ClobPair_Btc.Id,
		},
		Side:                            clobtypes.Order_SIDE_BUY,
		Quantums:                        1_000_000,
		Subticks:                        50_000_000_000,
		GoodTilOneof:                    &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: condTestGTBT},
		ConditionType:                   clobtypes.Order_CONDITION_TYPE_TAKE_PROFIT,
		ConditionalOrderTriggerSubticks: triggerSubticks,
	}
}
