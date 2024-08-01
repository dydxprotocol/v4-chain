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
	require.Equal(t, "ExpHt:", types.LegacyBlockHeightToPotentiallyPrunableOrdersPrefix)
}

func TestStoreAndMemstoreKeys(t *testing.T) {
	require.Equal(t, "SO/P/T:", types.TriggeredConditionalOrderKeyPrefix)
	require.Equal(t, "SO/P/L:", types.LongTermOrderPlacementKeyPrefix)
	require.Equal(t, "SO/U:", types.UntriggeredConditionalOrderKeyPrefix)

	require.Equal(t, "NumSO:", types.StatefulOrderCountPrefix)
	require.Equal(t, "ProposerEvents", types.ProcessProposerMatchesEventsKey)
}

func TestTransientStoreKeys(t *testing.T) {
	require.Equal(t, "SaLiqInfo:", types.SubaccountLiquidationInfoKeyPrefix)
	require.Equal(t, "NextTxIdx", types.NextStatefulOrderBlockTransactionIndexKey)
	require.Equal(t, "UncmtSO:", types.UncommittedStatefulOrderPlacementKeyPrefix)
	require.Equal(t, "UncmtSOCxl:", types.UncommittedStatefulOrderCancellationKeyPrefix)
	require.Equal(t, "NumUncmtSO:", types.UncommittedStatefulOrderCountPrefix)
}
