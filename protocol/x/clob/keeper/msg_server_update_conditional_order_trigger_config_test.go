package keeper_test

import (
	"testing"

	"github.com/cometbft/cometbft/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	"github.com/stretchr/testify/require"
)

// TestUpdateConditionalOrderTriggerConfig verifies the governance tuning lever for the
// conditional-order trigger mitigation: MsgUpdateConditionalOrderTriggerConfig is authority-gated,
// and a valid authority sets the budgets/caps in consensus state.
func TestUpdateConditionalOrderTriggerConfig(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(func() types.GenesisDoc {
		return testapp.DefaultGenesis()
	}).Build()
	ctx := tApp.InitChain()

	original := tApp.App.ClobKeeper.GetConditionalOrderTriggerConfig(ctx)

	handler := tApp.App.MsgServiceRouter().Handler(&clobtypes.MsgUpdateConditionalOrderTriggerConfig{})

	msg := clobtypes.MsgUpdateConditionalOrderTriggerConfig{
		Authority:                                    lib.GovModuleAddress.String(),
		MaxTriggersPerBlock:                          500,
		MaxRemovalsPerBlock:                          250,
		MaxUntriggeredConditionalOrdersGlobal:        123_456,
		MaxUntriggeredConditionalOrdersPerSubaccount: 77,
	}

	// Invalid authority is rejected and leaves the config unchanged.
	bad := msg
	bad.Authority = "fake authority"
	_, err := handler(ctx, &bad)
	require.Error(t, err, "invalid authority")
	require.Equal(t, original, tApp.App.ClobKeeper.GetConditionalOrderTriggerConfig(ctx))

	// Valid (gov) authority persists the budgets/caps.
	_, err = handler(ctx, &msg)
	require.NoError(t, err)
	got := tApp.App.ClobKeeper.GetConditionalOrderTriggerConfig(ctx)
	require.Equal(t, uint32(500), got.MaxTriggersPerBlock)
	require.Equal(t, uint32(250), got.MaxRemovalsPerBlock)
	require.Equal(t, uint32(123_456), got.MaxUntriggeredConditionalOrdersGlobal)
	require.Equal(t, uint32(77), got.MaxUntriggeredConditionalOrdersPerSubaccount)
}

// TestUpdateConditionalOrderTriggerConfig_ValidateBasicBounds covers:
// ValidateBasic rejects out-of-range budgets / caps (without an upper bound a proposal could set a
// MaxUint32 budget and erase the per-block work bound this feature exists to provide).
func TestUpdateConditionalOrderTriggerConfig_ValidateBasicBounds(t *testing.T) {
	base := clobtypes.MsgUpdateConditionalOrderTriggerConfig{
		Authority:                                    lib.GovModuleAddress.String(),
		MaxTriggersPerBlock:                          500,
		MaxRemovalsPerBlock:                          250,
		MaxUntriggeredConditionalOrdersGlobal:        100_000,
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

// TestUpdateConditionalOrderTriggerConfig_AllowsLargeLiveCount verifies a config update succeeds
// regardless of the current resting-set size (there is no one-shot cardinality ceiling; the index
// is already built by the upgrade and admission is per-placement bounded).
func TestUpdateConditionalOrderTriggerConfig_AllowsLargeLiveCount(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(func() types.GenesisDoc {
		return testapp.DefaultGenesis()
	}).Build()
	ctx := tApp.InitChain()
	handler := tApp.App.MsgServiceRouter().Handler(&clobtypes.MsgUpdateConditionalOrderTriggerConfig{})

	tApp.App.ClobKeeper.SetUntriggeredConditionalOrderCountGlobal(ctx, 1_000_000)

	msg := clobtypes.MsgUpdateConditionalOrderTriggerConfig{
		Authority:           lib.GovModuleAddress.String(),
		MaxTriggersPerBlock: 500,
	}
	_, err := handler(ctx, &msg)
	require.NoError(t, err)
	require.Equal(t, uint32(500), tApp.App.ClobKeeper.GetConditionalOrderTriggerConfig(ctx).MaxTriggersPerBlock)
}

// TestUpdateConditionalOrderTriggerConfig_AllowsCapBelowLiveCount preserves governance's ability to
// lower the admission cap below the current resting set (to stop admitting more while it drains).
func TestUpdateConditionalOrderTriggerConfig_AllowsCapBelowLiveCount(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(func() types.GenesisDoc {
		return testapp.DefaultGenesis()
	}).Build()
	ctx := tApp.InitChain()
	handler := tApp.App.MsgServiceRouter().Handler(&clobtypes.MsgUpdateConditionalOrderTriggerConfig{})

	// Live count is above the requested global cap.
	tApp.App.ClobKeeper.SetUntriggeredConditionalOrderCountGlobal(ctx, 100)

	msg := clobtypes.MsgUpdateConditionalOrderTriggerConfig{
		Authority:                             lib.GovModuleAddress.String(),
		MaxUntriggeredConditionalOrdersGlobal: 50, // below the live count of 100
	}
	_, err := handler(ctx, &msg)
	require.NoError(t, err)
	cfg := tApp.App.ClobKeeper.GetConditionalOrderTriggerConfig(ctx)
	require.Equal(t, uint32(50), cfg.MaxUntriggeredConditionalOrdersGlobal)
}
