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

// TestUpdateConditionalOrderTriggerConfig_ValidateBasicBounds covers:
// ValidateBasic rejects out-of-range budgets / caps (previously the setter only normalized zeros
// and applied no upper bound, so a proposal could set a MaxUint32 budget and erase the work bound).
func TestUpdateConditionalOrderTriggerConfig_ValidateBasicBounds(t *testing.T) {
	base := clobtypes.MsgUpdateConditionalOrderTriggerConfig{
		Authority:                             lib.GovModuleAddress.String(),
		Enabled:                               true,
		MaxTriggersPerBlock:                   500,
		MaxRemovalsPerBlock:                   250,
		MaxUntriggeredConditionalOrdersGlobal: 100_000,
		MaxUntriggeredConditionalOrdersPerSubaccount: 200,
	}

	// Valid message passes.
	require.NoError(t, base.ValidateBasic())
	atBounds := base
	atBounds.MaxTriggersPerBlock = clobtypes.MaxConfigurableTriggersPerBlock
	atBounds.MaxRemovalsPerBlock = clobtypes.MaxConfigurableRemovalsPerBlock
	atBounds.MaxUntriggeredConditionalOrdersGlobal = clobtypes.MaxConfigurableUntriggeredGlobal
	atBounds.MaxUntriggeredConditionalOrdersPerSubaccount = clobtypes.MaxConfigurableUntriggeredPerSubaccount
	require.NoError(t, atBounds.ValidateBasic())

	// Zeros are allowed (setter normalizes to defaults).
	zeros := base
	zeros.MaxTriggersPerBlock = 0
	zeros.MaxRemovalsPerBlock = 0
	zeros.MaxUntriggeredConditionalOrdersGlobal = 0
	zeros.MaxUntriggeredConditionalOrdersPerSubaccount = 0
	require.NoError(t, zeros.ValidateBasic())

	// MaxUint32 trigger budget is rejected.
	tooManyTriggers := base
	tooManyTriggers.MaxTriggersPerBlock = ^uint32(0)
	require.ErrorIs(t, tooManyTriggers.ValidateBasic(), clobtypes.ErrInvalidConditionalOrderTriggerConfig)

	// MaxUint32 removal budget is rejected.
	tooManyRemovals := base
	tooManyRemovals.MaxRemovalsPerBlock = ^uint32(0)
	require.ErrorIs(t, tooManyRemovals.ValidateBasic(), clobtypes.ErrInvalidConditionalOrderTriggerConfig)

	tooManyGlobal := base
	tooManyGlobal.MaxUntriggeredConditionalOrdersGlobal = clobtypes.MaxConfigurableUntriggeredGlobal + 1
	require.ErrorIs(t, tooManyGlobal.ValidateBasic(), clobtypes.ErrInvalidConditionalOrderTriggerConfig)

	tooManyPerSubaccount := base
	tooManyPerSubaccount.MaxUntriggeredConditionalOrdersGlobal = clobtypes.MaxConfigurableUntriggeredGlobal
	tooManyPerSubaccount.MaxUntriggeredConditionalOrdersPerSubaccount =
		clobtypes.MaxConfigurableUntriggeredPerSubaccount + 1
	require.ErrorIs(t, tooManyPerSubaccount.ValidateBasic(), clobtypes.ErrInvalidConditionalOrderTriggerConfig)

	// Per-subaccount cap exceeding the global cap is rejected.
	inverted := base
	inverted.MaxUntriggeredConditionalOrdersGlobal = 100
	inverted.MaxUntriggeredConditionalOrdersPerSubaccount = 1_000
	require.ErrorIs(t, inverted.ValidateBasic(), clobtypes.ErrInvalidConditionalOrderTriggerConfig)

	// Invalid authority is rejected.
	badAuth := base
	badAuth.Authority = "not-a-bech32-address"
	require.Error(t, badAuth.ValidateBasic())
}

// TestUpdateConditionalOrderTriggerConfig_EnableAllowsLargeLiveCount covers incremental activation:
// enable no longer depends on a one-shot cardinality ceiling because reconciliation is per-block bounded.
func TestUpdateConditionalOrderTriggerConfig_EnableAllowsLargeLiveCount(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(func() types.GenesisDoc {
		return testapp.DefaultGenesis()
	}).Build()
	ctx := tApp.InitChain()
	handler := tApp.App.MsgServiceRouter().Handler(&clobtypes.MsgUpdateConditionalOrderTriggerConfig{})

	tApp.App.ClobKeeper.SetUntriggeredConditionalOrderCountGlobal(ctx, 1_000_000)

	msg := clobtypes.MsgUpdateConditionalOrderTriggerConfig{
		Authority: lib.GovModuleAddress.String(),
		Enabled:   true,
	}
	_, err := handler(ctx, &msg)
	require.NoError(t, err)
	require.True(t, tApp.App.ClobKeeper.GetConditionalOrderTriggerConfig(ctx).Enabled)
}

// TestUpdateConditionalOrderTriggerConfig_EnableAllowsCapBelowLiveCount preserves governance's
// ability to lower the cap immediately and drain an oversized resting set without admitting more.
func TestUpdateConditionalOrderTriggerConfig_EnableAllowsCapBelowLiveCount(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(func() types.GenesisDoc {
		return testapp.DefaultGenesis()
	}).Build()
	ctx := tApp.InitChain()
	handler := tApp.App.MsgServiceRouter().Handler(&clobtypes.MsgUpdateConditionalOrderTriggerConfig{})

	// Live count is above the requested global cap.
	tApp.App.ClobKeeper.SetUntriggeredConditionalOrderCountGlobal(ctx, 100)

	msg := clobtypes.MsgUpdateConditionalOrderTriggerConfig{
		Authority:                             lib.GovModuleAddress.String(),
		Enabled:                               true,
		MaxUntriggeredConditionalOrdersGlobal: 50, // below the live count of 100
	}
	_, err := handler(ctx, &msg)
	require.NoError(t, err)
	require.True(t, tApp.App.ClobKeeper.GetConditionalOrderTriggerConfig(ctx).Enabled)
}
