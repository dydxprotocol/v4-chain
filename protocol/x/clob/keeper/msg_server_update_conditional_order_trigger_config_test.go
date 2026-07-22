package keeper_test

import (
	"testing"

	"github.com/cometbft/cometbft/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	"github.com/stretchr/testify/require"
)

// TestUpdateConditionalOrderTriggerConfig verifies the governance activation lever for the
// conditional-order EndBlocker work mitigation: MsgUpdateConditionalOrderTriggerConfig is authority-gated, and a valid
// authority flips the flag + sets the budgets/caps in consensus state.
func TestUpdateConditionalOrderTriggerConfig(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(func() types.GenesisDoc {
		return testapp.DefaultGenesis()
	}).Build()
	ctx := tApp.InitChain()

	// Default: disabled (legacy path).
	original := tApp.App.ClobKeeper.GetConditionalOrderTriggerConfig(ctx)
	require.False(t, original.Enabled)

	handler := tApp.App.MsgServiceRouter().Handler(&clobtypes.MsgUpdateConditionalOrderTriggerConfig{})

	msg := clobtypes.MsgUpdateConditionalOrderTriggerConfig{
		Authority:                             lib.GovModuleAddress.String(),
		Enabled:                               true,
		MaxTriggersPerBlock:                   500,
		MaxRemovalsPerBlock:                   250,
		MaxUntriggeredConditionalOrdersGlobal: 123_456,
		MaxUntriggeredConditionalOrdersPerSubaccount: 77,
	}

	// Invalid authority is rejected and leaves the config unchanged (still disabled).
	bad := msg
	bad.Authority = "fake authority"
	_, err := handler(ctx, &bad)
	require.Error(t, err, "invalid authority")
	require.False(t, tApp.App.ClobKeeper.GetConditionalOrderTriggerConfig(ctx).Enabled)

	// Valid (gov) authority activates the mitigation and persists the budgets/caps.
	_, err = handler(ctx, &msg)
	require.NoError(t, err)
	got := tApp.App.ClobKeeper.GetConditionalOrderTriggerConfig(ctx)
	require.True(t, got.Enabled)
	require.Equal(t, uint32(500), got.MaxTriggersPerBlock)
	require.Equal(t, uint32(250), got.MaxRemovalsPerBlock)
	require.Equal(t, uint32(123_456), got.MaxUntriggeredConditionalOrdersGlobal)
	require.Equal(t, uint32(77), got.MaxUntriggeredConditionalOrdersPerSubaccount)
}
