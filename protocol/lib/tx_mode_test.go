package lib_test

import (
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	testApp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAssertDeliverTxMode(t *testing.T) {
	tApp := testApp.NewTestAppBuilder().WithTesting(t).Build()
	// Initializing the chain returns a checkTx context so swap to a deliverTx context
	ctx := tApp.InitChain().WithIsCheckTx(false)

	require.NotPanics(t, func() {
		lib.AssertDeliverTxMode(ctx)
	})
	require.Panics(t, func() {
		lib.AssertDeliverTxMode(ctx.WithIsCheckTx(true))
	})
	require.Panics(t, func() {
		lib.AssertDeliverTxMode(ctx.WithIsReCheckTx(true))
	})
}

func TestIsDeliverTxMode(t *testing.T) {
	tApp := testApp.NewTestAppBuilder().WithTesting(t).Build()
	// Initializing the chain returns a checkTx context so swap to a deliverTx context
	ctx := tApp.InitChain().WithIsCheckTx(false)

	require.True(t, lib.IsDeliverTxMode(ctx))
	require.False(t, lib.IsDeliverTxMode(ctx.WithIsCheckTx(true)))
	require.False(t, lib.IsDeliverTxMode(ctx.WithIsReCheckTx(true)))
}

func TestAssertCheckTxMode(t *testing.T) {
	tApp := testApp.NewTestAppBuilder().WithTesting(t).Build()
	// Initializing the chain returns a checkTx context so swap to a deliverTx context
	ctx := tApp.InitChain().WithIsCheckTx(false)

	require.Panics(t, func() {
		lib.AssertCheckTxMode(ctx)
	})
	require.NotPanics(t, func() {
		lib.AssertCheckTxMode(ctx.WithIsCheckTx(true))
	})
	require.NotPanics(t, func() {
		lib.AssertCheckTxMode(ctx.WithIsReCheckTx(true))
	})
}

func TestTxMode(t *testing.T) {
	tApp := testApp.NewTestAppBuilder().WithTesting(t).Build()
	// Initializing the chain returns a checkTx context so swap to a deliverTx context
	ctx := tApp.InitChain().WithIsCheckTx(false)

	require.Equal(t, "DeliverTx", lib.TxMode(ctx))
	require.Equal(t, "CheckTx", lib.TxMode(ctx.WithIsCheckTx(true)))
	require.Equal(t, "ReCheckTx", lib.TxMode(ctx.WithIsReCheckTx(true)))
}
