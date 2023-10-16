package types_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	"github.com/stretchr/testify/require"
)

func TestModuleKeys(t *testing.T) {
	require.Equal(t, "clob", types.ModuleName)
	require.Equal(t, "clob", types.StoreKey)
	require.Equal(t, "mem_clob", types.MemStoreKey)
	require.Equal(t, "tmp_clob", types.TransientStoreKey)
}

func TestStateKeys(t *testing.T) {
	require.Equal(t, "SO/", types.StatefulOrderKeyPrefix)
	require.Equal(t, "SO/P/", types.PlacedStatefulOrderKeyPrefix)

	require.Equal(t, "LiqCfg", types.LiquidationsConfigKey)
	require.Equal(t, "EqTierCfg", types.EquityTierLimitConfigKey)
	require.Equal(t, "RateLimCfg", types.BlockRateLimitConfigKey)

	require.Equal(t, "Clob:", types.ClobPairKeyPrefix)
	require.Equal(t, "Fill:", types.OrderAmountFilledKeyPrefix)
	require.Equal(t, "ExpHt:", types.BlockHeightToPotentiallyPrunableOrdersPrefix)
	require.Equal(t, "ExpTm:", types.StatefulOrdersTimeSlicePrefix)
}

func TestStoreAndMemstoreKeys(t *testing.T) {
	require.Equal(t, "SO/P/T:", types.TriggeredConditionalOrderKeyPrefix)
	require.Equal(t, "SO/P/L:", types.LongTermOrderPlacementKeyPrefix)
	require.Equal(t, "SO/U:", types.UntriggeredConditionalOrderKeyPrefix)

	require.Equal(t, "ProposerEvents", types.ProcessProposerMatchesEventsKey)
}

func TestTransientStoreKeys(t *testing.T) {
	require.Equal(t, "SaLiqInfo:", types.SubaccountLiquidationInfoKeyPrefix)
	require.Equal(t, "NextTxIdx", types.NextStatefulOrderBlockTransactionIndexKey)
	require.Equal(t, "UncmtLT:", types.UncommittedStatefulOrderPlacementKeyPrefix)
	require.Equal(t, "UncmtLTCxl:", types.UncommittedStatefulOrderCancellationKeyPrefix)
	require.Equal(t, "NumUncmtLT:", types.UncommittedStatefulOrderCountPrefix)
	require.Equal(t, "NumLT:", types.StatefulOrderCountPrefix)
}

func TestModuleAccountKeys(t *testing.T) {
	require.Equal(t, "insurance_fund", types.InsuranceFundName)
}
