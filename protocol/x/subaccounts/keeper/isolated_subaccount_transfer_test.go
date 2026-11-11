package keeper_test

import (
	"errors"
	"math/big"
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/dydxprotocol/v4-chain/protocol/app/module"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	"github.com/dydxprotocol/v4-chain/protocol/streaming"
	bank_testutil "github.com/dydxprotocol/v4-chain/protocol/testutil/bank"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	testutil "github.com/dydxprotocol/v4-chain/protocol/testutil/util"
	asskeeper "github.com/dydxprotocol/v4-chain/protocol/x/assets/keeper"
	asstypes "github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
	perpskeeper "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/keeper"
	priceskeeper "github.com/dydxprotocol/v4-chain/protocol/x/prices/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

// isolatedPoolAddrForTest returns the isolated collateral pool address for the given perpetual id,
// mirroring Keeper.GetCollateralPoolFromPerpetualId for isolated perpetuals.
func isolatedPoolAddrForTest(perpetualId uint32) sdk.AccAddress {
	return authtypes.NewModuleAddress(types.ModuleName + ":" + lib.UintToString(perpetualId))
}

func setUpIsolatedPerpetualsForTest(
	t *testing.T,
	ctx sdk.Context,
	pricesKeeper *priceskeeper.Keeper,
	perpetualsKeeper *perpskeeper.Keeper,
	assetsKeeper *asskeeper.Keeper,
) {
	keepertest.CreateTestMarkets(t, ctx, pricesKeeper)
	keepertest.CreateTestLiquidityTiers(t, ctx, perpetualsKeeper)
	require.NoError(t, keepertest.CreateUsdcAsset(ctx, assetsKeeper))

	perpetualsKeeper.SetPerpetualForTest(ctx, constants.BtcUsd_NoMarginRequirement)
	perpetualsKeeper.SetPerpetualForTest(ctx, constants.IsoUsd_IsolatedMarket)
}

func TestUpdateSubaccounts_IsolatedCloseInsufficientPool_ReturnsWrappedError(t *testing.T) {
	ctx, k, pricesKeeper, perpetualsKeeper, _, bankKeeper, assetsKeeper, _, _, _, _ :=
		keepertest.SubaccountsKeepers(t, false)
	ctx = ctx.WithTxBytes(constants.TestTxBytes)
	setUpIsolatedPerpetualsForTest(t, ctx, pricesKeeper, perpetualsKeeper, assetsKeeper)

	isolatedPoolAddr := isolatedPoolAddrForTest(constants.IsoUsd_IsolatedMarket.Params.Id)
	crossPoolAddr := types.ModuleAddress

	const (
		crossPoolInitial    = int64(2_000_000_000_000) // $2,000,000
		isolatedPoolInitial = int64(500_000_000_000)   // $500,000: not enough to cover the $1,000,000 payout below.
	)
	require.NoError(t, bank_testutil.FundAccount(
		ctx,
		crossPoolAddr,
		sdk.Coins{sdk.NewCoin(asstypes.AssetUsdc.Denom, sdkmath.NewInt(crossPoolInitial))},
		*bankKeeper,
	))
	require.NoError(t, bank_testutil.FundAccount(
		ctx,
		isolatedPoolAddr,
		sdk.Coins{sdk.NewCoin(asstypes.AssetUsdc.Denom, sdkmath.NewInt(isolatedPoolInitial))},
		*bankKeeper,
	))

	subaccountId := types.SubaccountId{Owner: "isolated-close-insufficient-pool", Number: 0}
	subaccount := types.Subaccount{
		Id:                 &subaccountId,
		PerpetualPositions: []*types.PerpetualPosition{&constants.PerpetualPosition_OneISOLong},
		AssetPositions:     testutil.CreateUsdcAssetPositions(big.NewInt(999_900_000_000)), // $999,900
	}
	k.SetSubaccount(ctx, subaccount)

	updates := []types.Update{
		{
			SubaccountId: subaccountId,
			AssetUpdates: testutil.CreateUsdcAssetUpdates(big.NewInt(100_000_000)), // +$100
			PerpetualUpdates: []types.PerpetualUpdate{
				{
					PerpetualId:      constants.IsoUsd_IsolatedMarket.Params.Id,
					BigQuantumsDelta: big.NewInt(-1_000_000_000), // close the full 1 ISO position
				},
			},
		},
	}

	success, _, err := k.UpdateSubaccounts(ctx, updates, types.CollatCheck)

	require.False(t, success)
	require.Error(t, err)
	require.ErrorIs(t, err, types.ErrFailedToUpdateSubaccounts)

	// Subaccount state must be unchanged: the failure happens before any SetSubaccount call.
	gotSubaccount := k.GetSubaccount(ctx, subaccountId)
	require.Equal(t, subaccount, gotSubaccount)

	// Collateral pool balances must be unchanged: the failed SendCoins moves nothing.
	require.Equal(
		t,
		sdk.NewCoin(asstypes.AssetUsdc.Denom, sdkmath.NewInt(crossPoolInitial)),
		bankKeeper.GetBalance(ctx, crossPoolAddr, asstypes.AssetUsdc.Denom),
	)
	require.Equal(
		t,
		sdk.NewCoin(asstypes.AssetUsdc.Denom, sdkmath.NewInt(isolatedPoolInitial)),
		bankKeeper.GetBalance(ctx, isolatedPoolAddr, asstypes.AssetUsdc.Denom),
	)
}

// TestUpdateSubaccounts_MatchInflowBeforeOutflow_Succeeds reproduces the ordering regression at
// the keeper level: a single UpdateSubaccounts call carries a closer's update before an opener's
// update (matching the [taker, maker] order that persistMatchedOrders emits). The closer's payout
// requires the isolated pool's FULL resulting USDC balance ($1,000,000) to move out to the cross
// pool, but the pool only holds $700,000 up front. That $700,000 is exactly the closer's own
// pre-existing collateral sitting in the pool (unrelated to this call) -- it is less than the
// closer's payout ($1,000,000) but, together with the opener's $1,000,000 deposit, more than
// enough to cover it. So:
//   - Old order (Closed before Opened, i.e. processing updates in slice order [closer, opener]):
//     the closer's withdrawal of $1,000,000 is attempted first against a $700,000 balance and
//     fails with insufficient funds.
//   - New order (all Opened transitions before all Closed transitions): the opener's $1,000,000
//     deposit lands first, bringing the pool to $1,700,000, so the closer's $1,000,000 withdrawal
//     then succeeds, leaving the pool at $700,000 (unchanged net, since deposit == payout here).
func TestUpdateSubaccounts_MatchInflowBeforeOutflow_Succeeds(t *testing.T) {
	ctx, k, pricesKeeper, perpetualsKeeper, _, bankKeeper, assetsKeeper, _, _, _, _ :=
		keepertest.SubaccountsKeepers(t, false)
	ctx = ctx.WithTxBytes(constants.TestTxBytes)
	setUpIsolatedPerpetualsForTest(t, ctx, pricesKeeper, perpetualsKeeper, assetsKeeper)

	isolatedPoolAddr := isolatedPoolAddrForTest(constants.IsoUsd_IsolatedMarket.Params.Id)
	crossPoolAddr := types.ModuleAddress

	const (
		crossPoolInitial    = int64(2_000_000_000_000) // $2,000,000: ample funds to cover the opener's deposit.
		isolatedPoolInitial = int64(700_000_000_000)   // $700,000: the closer's own pre-existing collateral.
	)
	require.NoError(t, bank_testutil.FundAccount(
		ctx,
		crossPoolAddr,
		sdk.Coins{sdk.NewCoin(asstypes.AssetUsdc.Denom, sdkmath.NewInt(crossPoolInitial))},
		*bankKeeper,
	))
	require.NoError(t, bank_testutil.FundAccount(
		ctx,
		isolatedPoolAddr,
		sdk.Coins{sdk.NewCoin(asstypes.AssetUsdc.Denom, sdkmath.NewInt(isolatedPoolInitial))},
		*bankKeeper,
	))

	closerId := types.SubaccountId{Owner: "closer", Number: 0}
	closer := types.Subaccount{
		Id:                 &closerId,
		PerpetualPositions: []*types.PerpetualPosition{&constants.PerpetualPosition_OneISOLong},
		AssetPositions:     testutil.CreateUsdcAssetPositions(big.NewInt(999_900_000_000)), // $999,900
	}
	k.SetSubaccount(ctx, closer)

	openerId := types.SubaccountId{Owner: "opener", Number: 0}
	opener := types.Subaccount{
		Id:             &openerId,
		AssetPositions: testutil.CreateUsdcAssetPositions(big.NewInt(1_000_000_000_000)), // $1,000,000
	}
	k.SetSubaccount(ctx, opener)

	// persistMatchedOrders emits [taker, maker]; the closer (taker) is listed first, and the
	// opener (maker) second, to reproduce the exact slice order a real match produces.
	updates := []types.Update{
		{ // closer: fully closes the 1 ISO position for a $100 profit.
			SubaccountId: closerId,
			AssetUpdates: testutil.CreateUsdcAssetUpdates(big.NewInt(100_000_000)), // +$100
			PerpetualUpdates: []types.PerpetualUpdate{
				{
					PerpetualId:      constants.IsoUsd_IsolatedMarket.Params.Id,
					BigQuantumsDelta: big.NewInt(-1_000_000_000),
				},
			},
		},
		{ // opener: opens the opposing 1 ISO position, depositing its full pre-existing balance
			// (net of a $100 fee) into the isolated pool.
			SubaccountId: openerId,
			AssetUpdates: testutil.CreateUsdcAssetUpdates(big.NewInt(-100_000_000)), // -$100
			PerpetualUpdates: []types.PerpetualUpdate{
				{
					PerpetualId:      constants.IsoUsd_IsolatedMarket.Params.Id,
					BigQuantumsDelta: big.NewInt(1_000_000_000),
				},
			},
		},
	}

	success, successPerUpdate, err := k.UpdateSubaccounts(ctx, updates, types.CollatCheck)
	require.NoError(t, err)
	require.True(t, success)
	require.Equal(t, []types.UpdateResult{types.Success, types.Success}, successPerUpdate)

	gotCloser := k.GetSubaccount(ctx, closerId)
	require.Empty(t, gotCloser.PerpetualPositions)
	require.Equal(t, testutil.CreateUsdcAssetPositions(big.NewInt(1_000_000_000_000)), gotCloser.AssetPositions)

	gotOpener := k.GetSubaccount(ctx, openerId)
	require.Equal(
		t,
		[]*types.PerpetualPosition{
			testutil.CreateSinglePerpetualPosition(
				constants.IsoUsd_IsolatedMarket.Params.Id,
				big.NewInt(1_000_000_000),
				big.NewInt(0),
				big.NewInt(0),
			),
		},
		gotOpener.PerpetualPositions,
	)
	require.Equal(t, testutil.CreateUsdcAssetPositions(big.NewInt(999_900_000_000)), gotOpener.AssetPositions)

	// Pool final balance = initial + opener deposit ($1,000,000) - closer payout ($1,000,000).
	require.Equal(
		t,
		sdk.NewCoin(asstypes.AssetUsdc.Denom, sdkmath.NewInt(isolatedPoolInitial)),
		bankKeeper.GetBalance(ctx, isolatedPoolAddr, asstypes.AssetUsdc.Denom),
	)
	require.Equal(
		t,
		sdk.NewCoin(asstypes.AssetUsdc.Denom, sdkmath.NewInt(crossPoolInitial)),
		bankKeeper.GetBalance(ctx, crossPoolAddr, asstypes.AssetUsdc.Denom),
	)
}

// TestTransferCollateralForIsolatedPerpetual_SendCoinsErrorWrapped exercises the unexported
// transferCollateralForIsolatedPerpetual through the exported UpdateSubaccounts entry point, using
// a mocks.BankKeeper whose SendCoins fails. This asserts the raw SendCoins error is wrapped in
// types.ErrFailedToUpdateSubaccounts rather than bubbling up unwrapped.
func TestTransferCollateralForIsolatedPerpetual_SendCoinsErrorWrapped(t *testing.T) {
	ctx, origKeeper, pricesKeeper, perpetualsKeeper, _, _, assetsKeeper, blocktimeKeeper, _, _, storeKey :=
		keepertest.SubaccountsKeepers(t, false)
	ctx = ctx.WithTxBytes(constants.TestTxBytes)
	setUpIsolatedPerpetualsForTest(t, ctx, pricesKeeper, perpetualsKeeper, assetsKeeper)

	isolatedPoolAddr := isolatedPoolAddrForTest(constants.IsoUsd_IsolatedMarket.Params.Id)

	mockBankKeeper := &mocks.BankKeeper{}
	sendErr := errors.New("send coins boom")
	mockBankKeeper.On("SendCoins", mock.Anything, isolatedPoolAddr, types.ModuleAddress, mock.Anything).
		Return(sendErr)

	k := keeper.NewKeeper(
		codec.NewProtoCodec(module.InterfaceRegistry),
		storeKey,
		assetsKeeper,
		mockBankKeeper,
		perpetualsKeeper,
		blocktimeKeeper,
		origKeeper.GetIndexerEventManager(),
		streaming.NewNoopGrpcStreamingManager(),
	)

	subaccountId := types.SubaccountId{Owner: "mock-bank-sendcoins-error", Number: 0}
	k.SetSubaccount(ctx, types.Subaccount{
		Id:                 &subaccountId,
		PerpetualPositions: []*types.PerpetualPosition{&constants.PerpetualPosition_OneISOLong},
		AssetPositions:     testutil.CreateUsdcAssetPositions(big.NewInt(999_900_000_000)), // $999,900
	})

	updates := []types.Update{
		{
			SubaccountId: subaccountId,
			AssetUpdates: testutil.CreateUsdcAssetUpdates(big.NewInt(100_000_000)), // +$100
			PerpetualUpdates: []types.PerpetualUpdate{
				{
					PerpetualId:      constants.IsoUsd_IsolatedMarket.Params.Id,
					BigQuantumsDelta: big.NewInt(-1_000_000_000), // close the full 1 ISO position
				},
			},
		},
	}

	success, _, err := k.UpdateSubaccounts(ctx, updates, types.CollatCheck)

	require.False(t, success)
	require.Error(t, err)
	require.ErrorIs(t, err, types.ErrFailedToUpdateSubaccounts)
	require.ErrorContains(t, err, sendErr.Error())

	mockBankKeeper.AssertExpectations(t)
}
