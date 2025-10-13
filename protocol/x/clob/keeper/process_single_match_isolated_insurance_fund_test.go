package keeper_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/cometbft/cometbft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	affiliatetypes "github.com/dydxprotocol/v4-chain/protocol/x/affiliates/types"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	feetierstypes "github.com/dydxprotocol/v4-chain/protocol/x/feetiers/types"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/require"
)

// TestProcessSingleMatch_IsolatedMarket_NegativeInsuranceFundDelta demonstrates the ordering bug
// where transferring collateral from the isolated pool fails because insurance fund payment hasn't
// been added to the pool yet.
// This test is expected to
// - fail before the fix that moves insurance payment before collateral transfers
// - succeed after the fix
func TestProcessSingleMatch_IsolatedMarket_NegativeInsuranceFundDelta(t *testing.T) {
	// SCENARIO:
	// - Carl: Short 1 ISO, quote +$49 (underwater after match)
	// - Alice: Long 1 ISO, quote +$100 (well-collateralized)
	// - Both have existing positions in the isolated market
	// - Match price: $100/ISO (higher than Carl can afford)
	//
	// ISOLATED POOL:
	// - Has $149 (Carl's $49 + Alice's $100)
	// - This is the correct amount for their collateral
	//
	// WHAT HAPPENS WHEN THEY CLOSE AT $100:
	// - Carl: Pays $100 to close short, has -$51, insurance covers → final USDC = $0
	// - Alice: Receives $100 from closing long → final USDC = $200
	// - System needs to transfer collateral from isolated pool to main pool when positions close
	//
	// INSURANCE FUND DELTA:
	// - Insurance fund has negative delta of -$51 (i.e. pays INTO pool)
	// - This payment would bring pool to $149 + $51 = $200
	//
	// THE BUG:
	// - UpdateSubaccounts() (line 450) tries to transfer $200 of collateral FIRST
	// - Pool only has $149 as insurance fund delta isn't transferred into the pool yet
	//   - TransferInsuranceFundPayments() (line 487) would add $51
	// Result: Bank transfer error "spendable balance 149000000 is smaller than 200000000"

	// Carl: short 1 ISO, quote $49
	subaccount1 := constants.Carl_Num0_1ISO_Short_49USD

	// Alice: long 1 ISO, quote $100
	subaccount2 := satypes.Subaccount{
		Id: &constants.Alice_Num0,
		AssetPositions: []*satypes.AssetPosition{
			{
				AssetId:  0,
				Quantums: dtypes.NewInt(100_000_000),
			},
		},
		PerpetualPositions: []*satypes.PerpetualPosition{
			{
				PerpetualId:  3,
				Quantums:     dtypes.NewInt(1_000_000_000),
				FundingIndex: dtypes.NewInt(0),
			},
		},
	}

	tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(func() (genesis types.GenesisDoc) {
		genesis = testapp.DefaultGenesis()
		testapp.UpdateGenesisDocWithAppStateForModule(
			&genesis,
			func(genesisState *satypes.GenesisState) {
				genesisState.Subaccounts = []satypes.Subaccount{subaccount1, subaccount2}
			},
		)
		// Set ALL fees to 0 to isolate the bug
		testapp.UpdateGenesisDocWithAppStateForModule(
			&genesis,
			func(genesisState *feetierstypes.GenesisState) {
				genesisState.Params = constants.PerpetualFeeParamsNoFee
			},
		)
		testapp.UpdateGenesisDocWithAppStateForModule(
			&genesis,
			func(genesisState *perptypes.GenesisState) {
				genesisState.Params = constants.PerpetualsGenesisParams
				genesisState.LiquidityTiers = constants.LiquidityTiers
				// Set open interest to 1 ISO (Carl short + Alice long)
				isoPerpetual := constants.IsoUsd_IsolatedMarket
				isoPerpetual.OpenInterest = dtypes.NewInt(1_000_000_000)
				genesisState.Perpetuals = []perptypes.Perpetual{
					constants.BtcUsd_20PercentInitial_10PercentMaintenance,
					constants.EthUsd_20PercentInitial_10PercentMaintenance,
					isoPerpetual,
				}
			},
		)
		testapp.UpdateGenesisDocWithAppStateForModule(
			&genesis,
			func(genesisState *clobtypes.GenesisState) {
				genesisState.ClobPairs = []clobtypes.ClobPair{
					constants.ClobPair_Btc,
					constants.ClobPair_Eth,
					constants.ClobPair_3_Iso,
				}
				genesisState.LiquidationsConfig = constants.LiquidationsConfig_No_Limit
			},
		)
		testapp.UpdateGenesisDocWithAppStateForModule(
			&genesis,
			func(genesisState *banktypes.GenesisState) {
				isolatedPoolAddr := authtypes.NewModuleAddress(
					satypes.ModuleName + ":" + lib.UintToString(constants.IsoUsd_IsolatedMarket.Params.Id),
				).String()
				insuranceFundAddr := authtypes.NewModuleAddress(
					perptypes.InsuranceFundName + ":" + lib.UintToString(constants.IsoUsd_IsolatedMarket.Params.Id),
				).String()

				foundIsolatedPool := false
				foundInsuranceFund := false

				for i, bal := range genesisState.Balances {
					if bal.Address == isolatedPoolAddr {
						genesisState.Balances[i] = banktypes.Balance{
							Address: bal.Address,
							Coins: sdk.NewCoins(
								// $149 = Carl's $49 + Alice's $100
								sdk.NewCoin(constants.Usdc.Denom, sdkmath.NewInt(149_000_000)),
							),
						}
						foundIsolatedPool = true
					}
					if bal.Address == insuranceFundAddr {
						genesisState.Balances[i] = banktypes.Balance{
							Address: bal.Address,
							Coins: sdk.NewCoins(
								sdk.NewCoin(constants.Usdc.Denom, sdkmath.NewInt(100_000_000_000)), // $100,000
							),
						}
						foundInsuranceFund = true
					}
				}

				if !foundIsolatedPool {
					genesisState.Balances = append(genesisState.Balances, banktypes.Balance{
						Address: isolatedPoolAddr,
						Coins: sdk.NewCoins(
							sdk.NewCoin(constants.Usdc.Denom, sdkmath.NewInt(149_000_000)), // $149
						),
					})
				}
				if !foundInsuranceFund {
					genesisState.Balances = append(genesisState.Balances, banktypes.Balance{
						Address: insuranceFundAddr,
						Coins: sdk.NewCoins(
							sdk.NewCoin(constants.Usdc.Denom, sdkmath.NewInt(100_000_000_000)), // $100,000
						),
					})
				}
			},
		)
		return genesis
	}).Build()
	ctx := tApp.InitChain()

	// Log initial state
	isolatedPoolAddr, err := tApp.App.SubaccountsKeeper.GetCollateralPoolFromPerpetualId(
		ctx,
		constants.IsoUsd_IsolatedMarket.Params.Id,
	)
	require.NoError(t, err)
	isolatedPoolBalance := tApp.App.BankKeeper.GetBalance(ctx, isolatedPoolAddr, constants.Usdc.Denom)

	insuranceFundBalance := tApp.App.SubaccountsKeeper.GetInsuranceFundBalance(
		ctx,
		constants.IsoUsd_IsolatedMarket.Params.Id,
	)

	carlInitial := tApp.App.SubaccountsKeeper.GetSubaccount(ctx, *subaccount1.Id)
	aliceInitial := tApp.App.SubaccountsKeeper.GetSubaccount(ctx, *subaccount2.Id)

	_, err = tApp.App.PerpetualsKeeper.GetPerpetual(ctx, constants.IsoUsd_IsolatedMarket.Params.Id)
	require.NoError(t, err)

	t.Logf("=== INITIAL STATE ===")
	t.Logf("Isolated pool: %s (Carl $49 + Alice $100 = $149)", isolatedPoolBalance.String())
	t.Logf("Insurance fund: %s", insuranceFundBalance.String())
	t.Logf("Carl: quote=%s, positions=%+v", carlInitial.GetUsdcPosition().String(), carlInitial.PerpetualPositions)
	t.Logf("Alice: quote=%s, positions=%+v", aliceInitial.GetUsdcPosition().String(), aliceInitial.PerpetualPositions)

	// Create liquidation match - Carl (short) and Alice (long) close positions
	// Price calculation: quoteQuantums = subticks × baseQuantums × 10^(quantumConversionExponent)
	// For $100 per ISO with 1 ISO (1,000,000,000 base quantums):
	//   100,000,000 = subticks × 1,000,000,000 × 10^(-8)
	//   100,000,000 = subticks × 10
	//   subticks = 10,000,000
	// At $100, Carl pays $100 to close his short:
	// - Carl final: $49 - $100 = -$51 (insurance covers $51)
	// - Insurance delta: -$51 (pays INTO pool)
	// - This creates larger shortfall requiring insurance payment
	liquidationOrder := clobtypes.NewLiquidationOrder(
		*subaccount1.Id,
		constants.ClobPair_3_Iso,
		true,                                // Carl buys to close short
		satypes.BaseQuantums(1_000_000_000), // 1 ISO
		clobtypes.Subticks(10_000_000),      // $100 per ISO
	)

	aliceOrder := clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: *subaccount2.Id,
			ClientId:     1,
			OrderFlags:   clobtypes.OrderIdFlags_ShortTerm,
			ClobPairId:   3,
		},
		Side:         clobtypes.Order_SIDE_SELL, // Alice sells to close long
		Quantums:     1_000_000_000,             // 1 ISO
		Subticks:     10_000_000,                // $100 per ISO
		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 100},
	}

	matchWithOrders := &clobtypes.MatchWithOrders{
		TakerOrder: liquidationOrder,
		MakerOrder: &aliceOrder,
		FillAmount: satypes.BaseQuantums(1_000_000_000),
	}

	// Verify fees are actually zero
	feeParams := tApp.App.FeeTiersKeeper.GetPerpetualFeeParams(ctx)
	t.Logf("\n=== FEE VERIFICATION ===")
	for i, tier := range feeParams.Tiers {
		t.Logf("Fee Tier %d: Maker=%d ppm, Taker=%d ppm", i, tier.MakerFeePpm, tier.TakerFeePpm)
	}

	t.Logf("\n=== PROCESSING MATCH (Carl closes short, Alice closes long) ===")
	t.Logf("Expected updates:")
	t.Logf("- Carl (taker): Buys 1 ISO at $100, pays $100, insurance covers $51")
	t.Logf("  - Quote delta: -$100 + $51 (insurance) = -$49")
	t.Logf("  - Perp delta: +1 ISO")
	t.Logf("- Alice (maker): Sells 1 ISO at $100, receives $100")
	t.Logf("  - Quote delta: +$100")
	t.Logf("  - Perp delta: -1 ISO")

	// Get pool balance before
	poolBalanceBefore := tApp.App.BankKeeper.GetBalance(ctx, isolatedPoolAddr, constants.Usdc.Denom)
	t.Logf("Pool balance BEFORE match: %s", poolBalanceBefore.String())

	success, takerResult, makerResult, _, err := tApp.App.ClobKeeper.ProcessSingleMatch(
		ctx,
		matchWithOrders,
		map[string]bool{},
		affiliatetypes.AffiliateParameters{},
	)

	// Get pool balance after
	poolBalanceAfter := tApp.App.BankKeeper.GetBalance(ctx, isolatedPoolAddr, constants.Usdc.Denom)
	t.Logf("Pool balance AFTER match attempt: %s", poolBalanceAfter.String())
	t.Logf("Pool balance CHANGED by: %s quantums", poolBalanceBefore.Amount.Sub(poolBalanceAfter.Amount).String())

	t.Logf("\n=== RESULT ===")
	t.Logf("Success: %v", success)
	t.Logf("Taker: %v, Maker: %v", takerResult, makerResult)
	if err != nil {
		t.Logf("Error: %v", err)
	}

	// DEMONSTRATE BUG
	// // Get subaccount states to see what partial updates happened
	// carlAfter := tApp.App.SubaccountsKeeper.GetSubaccount(ctx, *subaccount1.Id)
	// aliceAfter := tApp.App.SubaccountsKeeper.GetSubaccount(ctx, *subaccount2.Id)
	// t.Logf(
	// 	"\nCarl AFTER error: quote=%s, positions=%+v",
	// 	carlAfter.GetUsdcPosition().String(),
	// 	carlAfter.PerpetualPositions,
	// )
	// t.Logf(
	// 	"Alice AFTER error: quote=%s, positions=%+v",
	// 	aliceAfter.GetUsdcPosition().String(),
	// 	aliceAfter.PerpetualPositions,
	// )
	// t.Logf("(Both unchanged because the entire update was rolled back)")

	// // Analyze the error
	// t.Logf("\n=== ANALYSIS ===")
	// t.Logf("Pool balance: %s quantums ($%d)", poolBalanceBefore.Amount.String(), 149)
	// t.Logf("Insurance would add: $51 (bringing pool to $200)")
	// t.Logf("System tries to withdraw $200")
	// t.Logf("But UpdateSubaccounts() happens BEFORE TransferInsuranceFundPayments()")

	insuranceFundBalanceAfter := tApp.App.SubaccountsKeeper.GetInsuranceFundBalance(
		ctx,
		constants.IsoUsd_IsolatedMarket.Params.Id,
	)
	t.Logf("Insurance fund balance (unchanged due to error): %s", insuranceFundBalanceAfter.String())

	require.NoError(t, err, "Match should succeed after insurance payment is moved before collateral transfers")
	require.True(t, success, "Match should succeed")

	// Verify the final state after successful match
	carlFinal := tApp.App.SubaccountsKeeper.GetSubaccount(ctx, *subaccount1.Id)
	aliceFinal := tApp.App.SubaccountsKeeper.GetSubaccount(ctx, *subaccount2.Id)
	poolBalanceFinal := tApp.App.BankKeeper.GetBalance(ctx, isolatedPoolAddr, constants.Usdc.Denom)

	t.Logf("\n=== FINAL STATE (after successful match) ===")
	t.Logf("Carl final: quote=%s, positions=%+v", carlFinal.GetUsdcPosition().String(), carlFinal.PerpetualPositions)
	t.Logf("Alice final: quote=%s, positions=%+v", aliceFinal.GetUsdcPosition().String(), aliceFinal.PerpetualPositions)
	t.Logf("Pool final balance: %s", poolBalanceFinal.String())
	t.Logf("\n✓ Match succeeded because insurance payment happened before collateral transfers")
}
