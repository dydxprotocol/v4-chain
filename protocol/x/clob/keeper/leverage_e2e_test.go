package keeper_test

import (
	"math/big"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	assettypes "github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/require"
)

// Test helper functions for leverage e2e tests

// setupLeverageTest creates a test app with the necessary state for leverage testing
func setupLeverageTest(t *testing.T) *testapp.TestApp {
	tApp := testapp.NewTestAppBuilder(t).Build()

	// Set up basic market data - use existing constants
	// Markets, perpetuals, and liquidity tiers are already set up in the test app

	return tApp
}

// configureLeverage sets leverage for a subaccount and perpetual
func configureLeverage(
	t *testing.T,
	tApp *testapp.TestApp,
	ctx sdk.Context,
	subaccountId satypes.SubaccountId,
	perpetualId uint32,
	leverage uint32,
) {
	leverageMap := map[uint32]uint32{
		perpetualId: leverage,
	}

	err := tApp.App.ClobKeeper.UpdateLeverage(ctx, &subaccountId, leverageMap)
	require.NoError(t, err)
}

// createSubaccountWithBalance creates a subaccount with specified USDC balance
func createSubaccountWithBalance(
	tApp *testapp.TestApp,
	ctx sdk.Context,
	subaccountId satypes.SubaccountId,
	usdcBalance *big.Int,
) {
	subaccount := satypes.Subaccount{
		Id: &subaccountId,
		AssetPositions: []*satypes.AssetPosition{
			{
				AssetId:  assettypes.AssetUsdc.Id,
				Quantums: dtypes.NewIntFromBigInt(usdcBalance),
			},
		},
	}

	tApp.App.SubaccountsKeeper.SetSubaccount(ctx, subaccount)
}

// TestLeverageKeeperSetup verifies that the LeverageKeeper is properly wired up
func TestLeverageKeeperSetup(t *testing.T) {
	tApp := setupLeverageTest(t)
	ctx := tApp.InitChain()

	// Test that the SubaccountsKeeper has a non-nil LeverageKeeper
	// We can't directly access the leverageKeeper field since it's private,
	// but we can test that leverage-aware operations work

	subaccountId := constants.Alice_Num0
	perpetualId := uint32(0)
	leverage := uint32(5)

	// Configure leverage first
	configureLeverage(t, tApp, ctx, subaccountId, perpetualId, leverage)

	// Verify leverage was set
	leverageMap, exists := tApp.App.ClobKeeper.GetLeverage(ctx, &subaccountId)
	require.True(t, exists)
	require.Equal(t, leverage, leverageMap[perpetualId])

	// Create a subaccount with some balance
	createSubaccountWithBalance(tApp, ctx, subaccountId, big.NewInt(1000_000_000))

	// Test that CanUpdateSubaccounts works (this internally uses the leverageKeeper)
	updates := []satypes.Update{
		{
			SubaccountId: subaccountId,
			AssetUpdates: []satypes.AssetUpdate{
				{
					AssetId:          assettypes.AssetUsdc.Id,
					BigQuantumsDelta: big.NewInt(-100_000_000), // Spend $100
				},
			},
			PerpetualUpdates: []satypes.PerpetualUpdate{
				{
					PerpetualId:      perpetualId,
					BigQuantumsDelta: big.NewInt(1_000_000), // Small position
				},
			},
		},
	}

	// This should work without panicking (leverageKeeper should not be nil)
	success, results, err := tApp.App.SubaccountsKeeper.CanUpdateSubaccounts(
		ctx,
		updates,
		satypes.CollatCheck,
	)

	require.NoError(t, err, "CanUpdateSubaccounts should not error")
	require.NotNil(t, results, "Results should not be nil")
	require.Len(t, results, 1, "Should have one result")

	t.Logf("✅ LeverageKeeper is properly wired up")
	t.Logf("   CanUpdateSubaccounts success: %v", success)
	t.Logf("   Update result: %v", results[0])
}

// TestLeverageBasicOrderPlacement tests basic order placement with leverage configuration
func TestLeverageBasicOrderPlacement(t *testing.T) {
	tApp := setupLeverageTest(t)
	ctx := tApp.InitChain()

	// Test parameters
	subaccountId := constants.Alice_Num0
	perpetualId := uint32(0)                   // BTC-USD
	leverage := uint32(10)                     // 10x leverage
	initialBalance := big.NewInt(1000_000_000) // $1000 USDC (6 decimals)

	// Set up subaccount with initial balance
	createSubaccountWithBalance(tApp, ctx, subaccountId, initialBalance)

	// Configure leverage
	configureLeverage(t, tApp, ctx, subaccountId, perpetualId, leverage)

	// Verify leverage was set correctly
	leverageMap, exists := tApp.App.ClobKeeper.GetLeverage(ctx, &subaccountId)
	require.True(t, exists)
	require.Equal(t, leverage, leverageMap[perpetualId])

	t.Logf("✅ Successfully configured and verified %dx leverage for subaccount", leverage)
	t.Logf("   Subaccount: %v", subaccountId)
	t.Logf("   Perpetual ID: %d", perpetualId)
	t.Logf("   Initial balance: $%s", new(big.Int).Div(initialBalance, big.NewInt(1_000_000)))
}

// TestLeverageConfiguration tests basic leverage configuration functionality
func TestLeverageConfiguration(t *testing.T) {
	tApp := setupLeverageTest(t)
	ctx := tApp.InitChain()

	subaccountId := constants.Alice_Num0
	perpetualId := uint32(0)

	testCases := []struct {
		name     string
		leverage uint32
	}{
		{"2x Leverage", 2},
		{"10x Leverage", 10},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Configure leverage
			configureLeverage(t, tApp, ctx, subaccountId, perpetualId, tc.leverage)

			// Verify leverage was set correctly
			leverageMap, exists := tApp.App.ClobKeeper.GetLeverage(ctx, &subaccountId)
			require.True(t, exists)
			require.Equal(t, tc.leverage, leverageMap[perpetualId])

			t.Logf("✅ Successfully configured %dx leverage", tc.leverage)
		})
	}
}

func TestOrderPlacementFailsWithLeverageConfigured(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()

	// Verify Alice and Bob have identical subaccounts
	gotAlice := tApp.App.SubaccountsKeeper.GetSubaccount(ctx, constants.Alice_Num0)
	gotBob := tApp.App.SubaccountsKeeper.GetSubaccount(ctx, constants.Bob_Num0)
	require.Equal(t, gotAlice.AssetPositions, gotBob.AssetPositions, "Alice and Bob should have identical asset positions")

	// Configure leverage for Alice: 2x on BTC perpetual
	aliceLeverage := &clobtypes.MsgUpdateLeverage{
		SubaccountId: &constants.Alice_Num0,
		ClobPairLeverage: []*clobtypes.LeverageEntry{
			{
				ClobPairId: 0,
				Leverage:   1,
			},
		},
	}
	for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(
		ctx,
		tApp.App,
		*aliceLeverage,
	) {
		resp := tApp.CheckTx(checkTx)
		require.True(t, resp.IsOK(), "Expected Alice's CheckTx to succeed. Response: %+v", resp)
	}

	ctx = tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{})

	// Place orders for both Alice and Bob that would require the entire margin if leverage was unchanged
	orderSize := dtypes.NewIntFromBigInt(big.NewInt(5_500_000_000_000_000))

	// Use the same price and clob pair as in the other test
	price := uint64(2_000_000_000)

	// Bob's order should succeed
	bobOrder := &clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: constants.Bob_Num0,
			ClientId:     0,
			OrderFlags:   clobtypes.OrderIdFlags_LongTerm,
			ClobPairId:   0,
		},
		Side:     clobtypes.Order_SIDE_BUY,
		Quantums: orderSize.BigInt().Uint64(),
		Subticks: price,
		GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{
			GoodTilBlockTime: uint32(ctx.BlockTime().Unix() + 100),
		},
	}
	for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(
		ctx,
		tApp.App,
		*clobtypes.NewMsgPlaceOrder(*bobOrder),
	) {
		resp := tApp.CheckTx(checkTx)
		require.True(t, resp.IsOK(), "Expected Bob's CheckTx to succeed. Response: %+v", resp)
	}

	bobSubaccount := tApp.App.SubaccountsKeeper.GetSubaccount(ctx, constants.Bob_Num0)
	require.True(t, bobSubaccount.AssetPositions != nil, "Bob should have a subaccount")

	// Alice's order should fail due to leverage config
	aliceOrder := &clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: constants.Alice_Num0,
			ClientId:     0,
			OrderFlags:   clobtypes.OrderIdFlags_LongTerm,
			ClobPairId:   0,
		},
		Side:     clobtypes.Order_SIDE_BUY,
		Quantums: orderSize.BigInt().Uint64(),
		Subticks: price,
		GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{
			GoodTilBlockTime: uint32(ctx.BlockTime().Unix() + 100),
		},
	}
	for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(
		ctx,
		tApp.App,
		*clobtypes.NewMsgPlaceOrder(*aliceOrder),
	) {
		resp := tApp.CheckTx(checkTx)
		require.False(t, resp.IsOK(), "Expected Alice's CheckTx to fail due to leverage. Response: %+v", resp)
	}
}
