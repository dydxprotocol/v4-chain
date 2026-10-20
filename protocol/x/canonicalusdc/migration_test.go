package canonicalusdc

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/dydxprotocol/v4-chain/protocol/x/canonicalusdc/types"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Outcome codes for propChain.settle.
const (
	migrationAckSuccess = 0
	migrationAckError   = 1
	migrationTimeout    = 2
)

// migrationLoop runs the ADR replacement loop: receive a batch of Injective
// USDC, send the same logical amount out over Noble and settle it. It stops
// once Noble backing is drained, checking the full ledger model after each step.
func migrationLoop(t *testing.T, p *propChain, migrator sdk.AccAddress, batch int64) {
	t.Helper()
	for i := 0; i < 100; i++ {
		noble, _, _, _ := p.ledger()
		if noble.IsZero() {
			return
		}
		amount := min(batch, noble.Int64())
		require.NoError(t, p.recv(p.r.controls.InjectiveChannel, p.r.controls.InjectivePacketDenom, migrator, amount))
		p.requireModel(t, "injective batch in")
		require.NoError(t, p.send(migrator, p.r.controls.NobleChannel, amount))
		require.NoError(t, p.settle(len(p.inflight)-1, migrationAckSuccess))
		p.requireModel(t, "noble batch out")
	}
	t.Fatal("noble backing did not drain within 100 batches")
}

func migrationFinalize(p *propChain) error {
	controls := p.r.keeper.GetControls(p.r.ctx)
	controls.NobleWithdrawalsEnabled = false
	return p.updateControls(controls)
}

// The full replacement: Noble backing drains to zero batch by batch, Noble
// withdrawals can only be turned off at zero, and afterwards logical USDC is
// fully Injective-backed and exits only via Injective.
func TestMigrationReplacesNobleBackingWithInjective(t *testing.T) {
	p := propSetup(t)
	migrator, holder := p.holders[0], p.holders[1]
	supply := p.r.bank.GetSupply(p.r.ctx, p.r.controls.LogicalDenom).Amount

	require.NoError(t, p.recv(p.r.controls.InjectiveChannel, p.r.controls.InjectivePacketDenom, migrator, 250))
	require.NoError(t, p.send(migrator, p.r.controls.NobleChannel, 250))
	require.NoError(t, p.settle(0, migrationAckSuccess))
	require.ErrorIs(t, migrationFinalize(p), types.ErrNobleBackingRemains)

	migrationLoop(t, p, migrator, 250)
	noble, injective, pendingInjective, pendingNoble := p.ledger()
	require.True(t, noble.IsZero())
	require.True(t, pendingInjective.IsZero())
	require.True(t, pendingNoble.IsZero())
	require.Equal(t, supply, injective, "supply is fully Injective-backed")
	require.Equal(t, supply, p.r.bank.GetSupply(p.r.ctx, p.r.controls.LogicalDenom).Amount)

	require.NoError(t, migrationFinalize(p))
	require.False(t, p.r.keeper.GetControls(p.r.ctx).NobleWithdrawalsEnabled)
	p.requireModel(t, "finalized")

	require.ErrorIs(t, p.send(holder, p.r.controls.NobleChannel, 10), types.ErrDisabled)
	require.Error(t, p.recv(p.r.controls.NobleChannel, p.r.controls.NoblePacketDenom, holder, 10),
		"Noble deposits stay rejected")
	require.NoError(t, p.send(holder, p.r.controls.InjectiveChannel, 100))
	require.NoError(t, p.settle(len(p.inflight)-1, migrationAckSuccess))
	_, injective, _, _ = p.ledger()
	require.Equal(t, supply.SubRaw(100), injective)
	p.requireModel(t, "exit via injective after finalization")
}

// Finalization waits for in-flight Noble settlements, and a failed last batch
// returns its amount to Noble backing so it must be drained again.
func TestMigrationFinalizationWaitsForNobleSettlement(t *testing.T) {
	p := propSetup(t)
	migrator := p.holders[0]
	require.NoError(t, p.recv(p.r.controls.InjectiveChannel, p.r.controls.InjectivePacketDenom, migrator, 950))
	require.NoError(t, p.send(migrator, p.r.controls.NobleChannel, 950))
	require.NoError(t, p.settle(len(p.inflight)-1, migrationAckSuccess))

	// The last batch is still in flight.
	require.NoError(t, p.recv(p.r.controls.InjectiveChannel, p.r.controls.InjectivePacketDenom, migrator, 50))
	require.NoError(t, p.send(migrator, p.r.controls.NobleChannel, 50))
	noble, _, _, pendingNoble := p.ledger()
	require.True(t, noble.IsZero())
	require.Equal(t, math.NewInt(50), pendingNoble)
	require.ErrorIs(t, migrationFinalize(p), types.ErrNobleBackingRemains, "pending Noble settlement blocks it")

	require.NoError(t, p.settle(len(p.inflight)-1, migrationTimeout))
	noble, _, _, _ = p.ledger()
	require.Equal(t, math.NewInt(50), noble, "refund returns to Noble backing")
	require.ErrorIs(t, migrationFinalize(p), types.ErrNobleBackingRemains)

	migrationLoop(t, p, migrator, 400)
	require.NoError(t, migrationFinalize(p))
	p.requireModel(t, "finalized after redrain")
}

// A Noble withdrawal sent while DISABLED and refunded after GRADUAL is credited
// back to Noble backing, so the migration still drains to exactly zero.
func TestMigrationDrainsLegacyNobleRefund(t *testing.T) {
	p := propSetup(t)
	migrator := p.holders[0]
	disabled := p.r.keeper.GetControls(p.r.ctx)
	disabled.Mode = types.Mode_MODE_DISABLED
	require.NoError(t, p.r.keeper.SetControls(p.r.ctx, disabled))
	require.NoError(t, p.r.keeper.SetLedger(p.r.ctx, types.DefaultLedger()))

	require.NoError(t, p.send(migrator, p.r.controls.NobleChannel, 30))
	require.Zero(t, p.r.keeper.PendingCount(p.r.ctx), "legacy send records no settlement")
	require.NoError(t, p.setMode(types.Mode_MODE_GRADUAL))
	noble, _, _, _ := p.ledger()
	require.Equal(t, int64(propInitialNoble-30), noble.Int64(), "snapshot excludes the burned amount")

	require.NoError(t, p.settle(len(p.inflight)-1, migrationAckError))
	noble, _, _, _ = p.ledger()
	require.Equal(t, int64(propInitialNoble), noble.Int64(), "legacy refund credited back")
	p.requireModel(t, "after legacy refund")

	migrationLoop(t, p, migrator, 300)
	require.NoError(t, migrationFinalize(p))
	_, injective, _, _ := p.ledger()
	require.Equal(t, p.r.bank.GetSupply(p.r.ctx, p.r.controls.LogicalDenom).Amount, injective)
}
