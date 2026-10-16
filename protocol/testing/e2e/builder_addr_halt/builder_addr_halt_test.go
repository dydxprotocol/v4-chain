package builder_addr_halt_test

import (
	"testing"
	"time"

	"github.com/cometbft/cometbft/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	icatypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	"github.com/stretchr/testify/require"
)

var GenesisTime = time.Unix(1690000000, 0)

// TestBuilderAddrBlocked_AllBlockedModuleAccounts verifies that every blocked
// module account is rejected as a builder address at CheckTx. Each of these
// would cause a permanent chain halt if accepted. Regression test for SEC-76.
func TestBuilderAddrBlocked_AllBlockedModuleAccounts(t *testing.T) {
	blockedAddresses := map[string]string{
		"fee_collector":          authtypes.NewModuleAddress(authtypes.FeeCollectorName).String(),
		"distribution":           authtypes.NewModuleAddress(distrtypes.ModuleName).String(),
		"bonded_tokens_pool":     authtypes.NewModuleAddress(stakingtypes.BondedPoolName).String(),
		"not_bonded_tokens_pool": authtypes.NewModuleAddress(stakingtypes.NotBondedPoolName).String(),
		"ibc_transfer":           authtypes.NewModuleAddress(ibctransfertypes.ModuleName).String(),
		"interchain_accounts":    authtypes.NewModuleAddress(icatypes.ModuleName).String(),
	}

	for name, addr := range blockedAddresses {
		t.Run(name, func(t *testing.T) {
			gtbt := uint32(GenesisTime.Add(24 * time.Hour).Unix())

			tApp := testapp.NewTestAppBuilder(t).
				WithNonDeterminismChecksEnabled(false).
				WithGenesisDocFn(func() (genesis types.GenesisDoc) {
					genesis = testapp.DefaultGenesis()
					genesis.GenesisTime = GenesisTime
					return genesis
				}).Build()
			ctx := tApp.InitChain()

			order := clobtypes.Order{
				OrderId: clobtypes.OrderId{
					SubaccountId: constants.Alice_Num0,
					ClientId:     1,
					OrderFlags:   clobtypes.OrderIdFlags_LongTerm,
					ClobPairId:   0,
				},
				Side:         clobtypes.Order_SIDE_BUY,
				Quantums:     100_000_000,
				Subticks:     100_000_000_000,
				GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: gtbt},
				BuilderCodeParameters: &clobtypes.BuilderCodeParameters{
					BuilderAddress: addr,
					FeePpm:         10_000,
				},
			}

			resp := tApp.CheckTx(testapp.MustMakeCheckTxsWithClobMsg(
				ctx, tApp.App, *clobtypes.NewMsgPlaceOrder(order),
			)[0])

			require.Conditionf(t, resp.IsErr,
				"Order with blocked builder address %s (%s) should be rejected. Response: %+v",
				name, addr, resp)
			require.Contains(t, resp.Log, "blocked module account")
		})
	}
}

// TestBuilderAddrBlocked_OriginalPoC reproduces the exact attack from the bug
// report: an honest long-term maker sell is resting on the book, then a malicious
// long-term crossing taker buy with fee_collector as BuilderAddress is submitted.
// Pre-patch, the taker order would be committed and trigger a permanent chain-halt
// panic in PrepareCheckState at block 3. Post-patch, the taker is rejected at
// CheckTx and blocks advance normally. Regression test for SEC-76 / Finding #85.
func TestBuilderAddrBlocked_OriginalPoC(t *testing.T) {
	feeCollectorAddr := authtypes.NewModuleAddress(authtypes.FeeCollectorName).String()
	gtbt := uint32(GenesisTime.Add(24 * time.Hour).Unix())

	tApp := testapp.NewTestAppBuilder(t).
		WithNonDeterminismChecksEnabled(false).
		WithGenesisDocFn(func() (genesis types.GenesisDoc) {
			genesis = testapp.DefaultGenesis()
			genesis.GenesisTime = GenesisTime
			return genesis
		}).Build()
	ctx := tApp.InitChain()

	// Step 1: Honest long-term maker SELL at 50,000 (resting on the book).
	makerSell := clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: constants.Carl_Num0,
			ClientId:     0,
			OrderFlags:   clobtypes.OrderIdFlags_LongTerm,
			ClobPairId:   0,
		},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     100_000_000,
		Subticks:     50_000_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: gtbt},
	}
	makerResp := tApp.CheckTx(testapp.MustMakeCheckTxsWithClobMsg(
		ctx, tApp.App, *clobtypes.NewMsgPlaceOrder(makerSell),
	)[0])
	require.Conditionf(t, makerResp.IsOK, "Maker order should succeed. Response: %+v", makerResp)

	// Step 2: Malicious long-term crossing taker BUY at 100,000 with
	// BuilderAddress = fee_collector. This is the exact PoC from the bug report.
	takerBuy := clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: constants.Alice_Num0,
			ClientId:     1,
			OrderFlags:   clobtypes.OrderIdFlags_LongTerm,
			ClobPairId:   0,
		},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     100_000_000,
		Subticks:     100_000_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: gtbt},
		BuilderCodeParameters: &clobtypes.BuilderCodeParameters{
			BuilderAddress: feeCollectorAddr,
			FeePpm:         10_000,
		},
	}
	takerResp := tApp.CheckTx(testapp.MustMakeCheckTxsWithClobMsg(
		ctx, tApp.App, *clobtypes.NewMsgPlaceOrder(takerBuy),
	)[0])

	// Step 3: The malicious taker must be rejected at CheckTx.
	require.Conditionf(t, takerResp.IsErr,
		"Taker order with blocked builder address should be rejected. Response: %+v", takerResp)
	require.Contains(t, takerResp.Log, "blocked module account")

	// Step 4: Advance to block 3. Pre-patch this is where PrepareCheckState
	// would panic. Post-patch, blocks advance normally.
	tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{})
	tApp.AdvanceToBlock(3, testapp.AdvanceToBlockOptions{})
}

// TestBuilderAddrBlocked_MakerSideVariant verifies that a resting maker order
// with a blocked builder address is also rejected. This prevents a delayed-
// detonation variant where the malicious order sits on the book until an
// honest taker crosses it. Regression test for SEC-76.
func TestBuilderAddrBlocked_MakerSideVariant(t *testing.T) {
	feeCollectorAddr := authtypes.NewModuleAddress(authtypes.FeeCollectorName).String()
	gtbt := uint32(GenesisTime.Add(24 * time.Hour).Unix())

	tApp := testapp.NewTestAppBuilder(t).
		WithNonDeterminismChecksEnabled(false).
		WithGenesisDocFn(func() (genesis types.GenesisDoc) {
			genesis = testapp.DefaultGenesis()
			genesis.GenesisTime = GenesisTime
			return genesis
		}).Build()
	ctx := tApp.InitChain()

	// Attempt to place a resting maker SELL with blocked builder address.
	makerSell := clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: constants.Carl_Num0,
			ClientId:     0,
			OrderFlags:   clobtypes.OrderIdFlags_LongTerm,
			ClobPairId:   0,
		},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     100_000_000,
		Subticks:     50_000_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: gtbt},
		BuilderCodeParameters: &clobtypes.BuilderCodeParameters{
			BuilderAddress: feeCollectorAddr,
			FeePpm:         10_000,
		},
	}
	makerResp := tApp.CheckTx(testapp.MustMakeCheckTxsWithClobMsg(
		ctx, tApp.App, *clobtypes.NewMsgPlaceOrder(makerSell),
	)[0])

	require.Conditionf(t, makerResp.IsErr,
		"Maker order with blocked builder address should be rejected. Response: %+v", makerResp)
	require.Contains(t, makerResp.Log, "blocked module account")

	// Place an honest taker BUY (no builder params) and advance blocks.
	// Confirm no panic — the malicious maker never entered state.
	takerBuy := clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: constants.Alice_Num0,
			ClientId:     1,
			OrderFlags:   clobtypes.OrderIdFlags_LongTerm,
			ClobPairId:   0,
		},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     100_000_000,
		Subticks:     100_000_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: gtbt},
	}
	takerResp := tApp.CheckTx(testapp.MustMakeCheckTxsWithClobMsg(
		ctx, tApp.App, *clobtypes.NewMsgPlaceOrder(takerBuy),
	)[0])
	require.Conditionf(t, takerResp.IsOK,
		"Honest taker order should succeed. Response: %+v", takerResp)

	tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{})
	tApp.AdvanceToBlock(3, testapp.AdvanceToBlockOptions{})
}

// TestBuilderAddrBlocked_ValidAddressSucceeds confirms that orders with
// non-blocked builder addresses pass validation and do not regress
// legitimate builder fee functionality.
func TestBuilderAddrBlocked_ValidAddressSucceeds(t *testing.T) {
	gtbt := uint32(GenesisTime.Add(24 * time.Hour).Unix())

	tApp := testapp.NewTestAppBuilder(t).
		WithNonDeterminismChecksEnabled(false).
		WithGenesisDocFn(func() (genesis types.GenesisDoc) {
			genesis = testapp.DefaultGenesis()
			genesis.GenesisTime = GenesisTime
			return genesis
		}).Build()
	ctx := tApp.InitChain()

	order := clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: constants.Alice_Num0,
			ClientId:     1,
			OrderFlags:   clobtypes.OrderIdFlags_LongTerm,
			ClobPairId:   0,
		},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     100_000_000,
		Subticks:     100_000_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: gtbt},
		BuilderCodeParameters: &clobtypes.BuilderCodeParameters{
			BuilderAddress: constants.Carl_Num0.Owner,
			FeePpm:         10_000,
		},
	}

	resp := tApp.CheckTx(testapp.MustMakeCheckTxsWithClobMsg(
		ctx, tApp.App, *clobtypes.NewMsgPlaceOrder(order),
	)[0])

	require.Conditionf(t, resp.IsOK,
		"Order with valid builder address should succeed. Response: %+v", resp)
}
