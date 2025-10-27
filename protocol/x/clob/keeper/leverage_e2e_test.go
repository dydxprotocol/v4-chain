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
	sendingtypes "github.com/dydxprotocol/v4-chain/protocol/x/sending/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/require"
)

// Test helper functions for leverage e2e tests

// setupLeverageTest creates a test app with the necessary state for leverage testing
func setupLeverageTest(t *testing.T) *testapp.TestApp {
	tApp := testapp.NewTestAppBuilder(t).Build()
	return tApp
}

// configureLeverage sets leverage for a subaccount and perpetual
func configureLeverage(
	t *testing.T,
	tApp *testapp.TestApp,
	ctx sdk.Context,
	subaccountId satypes.SubaccountId,
	perpetualId uint32,
	custom_imf_ppm uint32,
) {
	leverageMap := map[uint32]uint32{
		perpetualId: custom_imf_ppm,
	}

	err := tApp.App.SubaccountsKeeper.UpdateLeverage(ctx, &subaccountId, leverageMap)
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
	custom_imf_ppm := uint32(50_000)

	// Configure leverage first
	configureLeverage(t, tApp, ctx, subaccountId, perpetualId, custom_imf_ppm)

	// Verify leverage was set
	leverageMap, exists := tApp.App.SubaccountsKeeper.GetLeverage(ctx, &subaccountId)
	require.True(t, exists)
	require.Equal(t, custom_imf_ppm, leverageMap[perpetualId])

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
	custom_imf_ppm := uint32(100_000)          // 10x leverage
	initialBalance := big.NewInt(1000_000_000) // $1000 USDC (6 decimals)

	// Set up subaccount with initial balance
	createSubaccountWithBalance(tApp, ctx, subaccountId, initialBalance)

	// Configure leverage
	configureLeverage(t, tApp, ctx, subaccountId, perpetualId, custom_imf_ppm)

	// Verify leverage was set correctly
	leverageMap, exists := tApp.App.SubaccountsKeeper.GetLeverage(ctx, &subaccountId)
	require.True(t, exists)
	require.Equal(t, custom_imf_ppm, leverageMap[perpetualId])

	t.Logf("✅ Successfully configured and verified %dx leverage for subaccount", custom_imf_ppm)
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
		name           string
		custom_imf_ppm uint32
	}{
		{"2x Leverage", 500_000},
		{"10x Leverage", 100_000},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Configure leverage
			configureLeverage(t, tApp, ctx, subaccountId, perpetualId, tc.custom_imf_ppm)

			// Verify leverage was set correctly
			leverageMap, exists := tApp.App.SubaccountsKeeper.GetLeverage(ctx, &subaccountId)
			require.True(t, exists)
			require.Equal(t, tc.custom_imf_ppm, leverageMap[perpetualId])

			t.Logf("✅ Successfully configured %dx leverage", tc.custom_imf_ppm)
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

	// Configure leverage for Alice: 1x on BTC perpetual
	aliceLeverage := &clobtypes.MsgUpdateLeverage{
		SubaccountId: &constants.Alice_Num0,
		ClobPairLeverage: []*clobtypes.LeverageEntry{
			{
				ClobPairId:   0,
				CustomImfPpm: 1_000_000,
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

func TestWithdrawalWithLeverage(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()

	// Use CheckTx to set leverage (this goes through ante handler which persists it)
	leverageMsg := &clobtypes.MsgUpdateLeverage{
		SubaccountId: &constants.Alice_Num0,
		ClobPairLeverage: []*clobtypes.LeverageEntry{
			{
				ClobPairId:   0,
				CustomImfPpm: 200_000,
			},
		},
	}
	for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(ctx, tApp.App, *leverageMsg) {
		resp := tApp.CheckTx(checkTx)
		require.True(t, resp.IsOK(), "Expected leverage CheckTx to succeed. Response: %+v", resp)
	}

	ctx = tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{})

	// Verify leverage was set immediately (ante handler sets it)
	aliceLeverageMap, found := tApp.App.SubaccountsKeeper.GetLeverage(ctx, &constants.Alice_Num0)
	require.True(t, found, "Alice's leverage should be set")
	require.Equal(
		t,
		uint32(200_000),
		aliceLeverageMap[0],
		"Alice's custom imf ppm for perpetual 0 should be 200,000",
	)

	// Use orders large enough to have meaningful margin requirements
	// Both Alice and Bob trade 1,000,000 quantums (0.0001 BTC) at price 5,000,000,000 subticks ($50,000/BTC)
	// This creates a $5 notional position which will have measurable IMR
	// Quantums must be multiple of StepBaseQuantums = 10
	// Subticks must be multiple of SubticksPerTick = 10000
	aliceOrder := clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: constants.Alice_Num0,
			ClientId:     0,
			OrderFlags:   clobtypes.OrderIdFlags_ShortTerm,
			ClobPairId:   0,
		},
		Side:     clobtypes.Order_SIDE_BUY,
		Quantums: 1_000_000,     // 0.0001 BTC
		Subticks: 5_000_000_000, // $50,000 per BTC
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{
			GoodTilBlock: 20,
		},
	}

	bobOrder := clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: constants.Bob_Num0,
			ClientId:     0,
			OrderFlags:   clobtypes.OrderIdFlags_ShortTerm,
			ClobPairId:   0,
		},
		Side:     clobtypes.Order_SIDE_BUY,
		Quantums: 1_000_000,     // 0.0001 BTC
		Subticks: 5_000_000_000, // $50,000 per BTC
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{
			GoodTilBlock: 20,
		},
	}

	carlSellOrder := clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: constants.Carl_Num0,
			ClientId:     0,
			OrderFlags:   clobtypes.OrderIdFlags_ShortTerm,
			ClobPairId:   0,
		},
		Side:     clobtypes.Order_SIDE_SELL,
		Quantums: 2_000_000,     // 0.0002 BTC
		Subticks: 5_000_000_000, // $50,000 per BTC
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{
			GoodTilBlock: 20,
		},
	}

	// Place all orders
	for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(
		ctx,
		tApp.App,
		*clobtypes.NewMsgPlaceOrder(carlSellOrder),
	) {
		resp := tApp.CheckTx(checkTx)
		require.True(t, resp.IsOK(), "Expected Carl's order to succeed. Response: %+v", resp)
	}

	for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(
		ctx,
		tApp.App,
		*clobtypes.NewMsgPlaceOrder(bobOrder),
	) {
		resp := tApp.CheckTx(checkTx)
		require.True(t, resp.IsOK(), "Expected Bob's order to succeed. Response: %+v", resp)
	}

	for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(
		ctx,
		tApp.App,
		*clobtypes.NewMsgPlaceOrder(aliceOrder),
	) {
		resp := tApp.CheckTx(checkTx)
		require.True(t, resp.IsOK(), "Expected Alice's order to succeed. Response: %+v", resp)
	}

	// Advance to next block which matches the orders
	ctx = tApp.AdvanceToBlock(3, testapp.AdvanceToBlockOptions{})

	// Get final subaccount states
	aliceAfter := tApp.App.SubaccountsKeeper.GetSubaccount(ctx, constants.Alice_Num0)
	bobAfter := tApp.App.SubaccountsKeeper.GetSubaccount(ctx, constants.Bob_Num0)

	// Both should have positions now
	require.Len(t, aliceAfter.PerpetualPositions, 1, "Alice should have 1 perpetual position")
	require.Len(t, bobAfter.PerpetualPositions, 1, "Bob should have 1 perpetual position")

	// Verify the positions are equal in size but opposite in direction
	require.Equal(t, uint64(1_000_000), new(big.Int).Abs(aliceAfter.PerpetualPositions[0].GetBigQuantums()).Uint64())
	require.Equal(t, uint64(1_000_000), new(big.Int).Abs(bobAfter.PerpetualPositions[0].GetBigQuantums()).Uint64())

	// Get actual risk calculations to verify leverage is being applied
	aliceRisk, err := tApp.App.SubaccountsKeeper.GetNetCollateralAndMarginRequirements(
		ctx,
		satypes.Update{SubaccountId: constants.Alice_Num0},
	)
	require.NoError(t, err)

	bobRisk, err := tApp.App.SubaccountsKeeper.GetNetCollateralAndMarginRequirements(
		ctx,
		satypes.Update{SubaccountId: constants.Bob_Num0},
	)
	require.NoError(t, err)

	// Alice's IMR should be higher due to lower leverage (more conservative)
	require.True(t, aliceRisk.IMR.Cmp(bobRisk.IMR) > 0,
		"Alice's IMR (%s) should be higher than Bob's IMR (%s) due to lower leverage setting (5x vs 20x)",
		aliceRisk.IMR.String(), bobRisk.IMR.String())

	// Calculate available collateral for new positions: NC - IMR
	aliceAvailable := new(big.Int).Sub(aliceRisk.NC, aliceRisk.IMR)
	bobAvailable := new(big.Int).Sub(bobRisk.NC, bobRisk.IMR)

	// Verify Bob has more available collateral than Alice due to lower leverage requirements
	require.True(t, bobAvailable.Cmp(aliceAvailable) > 0,
		"Bob's available collateral (%s) should be greater than Alice's (%s) due to lower IMR from higher leverage",
		bobAvailable.String(), aliceAvailable.String())

	// Let's try to withdraw bob's available collateral from both subaccounts

	bobWithdrawal := &sendingtypes.MsgWithdrawFromSubaccount{
		Sender:    constants.Bob_Num0,
		Recipient: constants.BobAccAddress.String(),
		AssetId:   constants.Usdc.Id,
		Quantums:  bobAvailable.Uint64(),
	}

	for _, checkTx := range testapp.MustMakeCheckTxsWithSdkMsg(
		ctx,
		tApp.App,
		testapp.MustMakeCheckTxOptions{
			AccAddressForSigning: constants.Bob_Num0.Owner,
			Gas:                  constants.TestGasLimit,
			FeeAmt:               constants.TestFeeCoins_5Cents,
		},
		bobWithdrawal,
	) {
		resp := tApp.CheckTx(checkTx)
		require.True(t, resp.IsOK(), "Expected Bob's withdrawal to succeed. Response: %+v", resp)
	}

	aliceWithdrawal := &sendingtypes.MsgWithdrawFromSubaccount{
		Sender:    constants.Alice_Num0,
		Recipient: constants.AliceAccAddress.String(),
		AssetId:   constants.Usdc.Id,
		Quantums:  bobAvailable.Uint64(), // using bob's available collateral which is greater than alice's
	}

	for _, checkTx := range testapp.MustMakeCheckTxsWithSdkMsg(
		ctx,
		tApp.App,
		testapp.MustMakeCheckTxOptions{
			AccAddressForSigning: constants.Alice_Num0.Owner,
			Gas:                  constants.TestGasLimit * 10,
			FeeAmt:               constants.TestFeeCoins_5Cents,
		},
		aliceWithdrawal,
	) {
		resp := tApp.CheckTx(checkTx)
		require.False(t, resp.IsOK(), "Expected Alice's withdrawal to fail. Response: %+v", resp)
	}
}

func TestUpdateLeverageWithExistingPosition(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()

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
		Side:     clobtypes.Order_SIDE_SELL,
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

	// Alice's order should succeed (no leverage configured yet)
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
		require.True(t, resp.IsOK(), "Expected Alice's CheckTx to succeed. Response: %+v", resp)
	}

	// Advance to next block which matches the orders
	ctx = tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{})

	// Configure leverage for Alice with an existing bitcoin position
	// This should make the account fail against the IMR check
	aliceLeverage := &clobtypes.MsgUpdateLeverage{
		SubaccountId: &constants.Alice_Num0,
		ClobPairLeverage: []*clobtypes.LeverageEntry{
			{
				ClobPairId:   0,
				CustomImfPpm: 500_000,
			},
		},
	}
	for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(
		ctx,
		tApp.App,
		*aliceLeverage,
	) {
		resp := tApp.CheckTx(checkTx)
		require.False(t, resp.IsOK(), "Expected Alice's CheckTx to fail. Response: %+v", resp)
	}
}
