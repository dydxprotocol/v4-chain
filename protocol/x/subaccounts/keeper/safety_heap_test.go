package keeper_test

import (
	"math/big"
	"math/rand"
	"testing"

	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/app/config"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	testutil "github.com/dydxprotocol/v4-chain/protocol/testutil/util"
	"github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/keeper"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/require"
	"gopkg.in/typ.v4/slices"
)

func TestSafetyHeapInsertRemoveMin(t *testing.T) {
	perpetualId := uint32(0)
	side := satypes.Long
	totalSubaccounts := 1000

	// Create 1000 subaccounts with balances ranging from -500 to 500.
	// The subaccounts should be sorted by balance.
	allSubaccounts := make([]satypes.Subaccount, 0)
	for i := 0; i < totalSubaccounts; i++ {
		subaccount := satypes.Subaccount{
			Id: &satypes.SubaccountId{
				Owner: sdk.MustBech32ifyAddressBytes(
					config.Bech32PrefixAccAddr,
					constants.AliceAccAddress,
				),
				Number: uint32(i),
			},
			AssetPositions: testutil.CreateUsdcAssetPositions(
				// Create asset positions with balances ranging from -500 to 500.
				big.NewInt(int64(i - totalSubaccounts/2)),
			),
		}

		// Handle special case.
		if i-totalSubaccounts/2 == 0 {
			subaccount.AssetPositions = nil
		}

		allSubaccounts = append(allSubaccounts, subaccount)
	}

	for iter := 0; iter < 100; iter++ {
		// Setup keeper state and test parameters.
		ctx, subaccountsKeeper, _, _, _, _, _, _, _, _, _ := keepertest.SubaccountsKeepers(t, false)

		// Shuffle the subaccounts so that insertion order is random.
		slices.Shuffle(allSubaccounts)

		store := subaccountsKeeper.GetSafetyHeapStore(ctx, perpetualId, side)
		for i, subaccount := range allSubaccounts {
			subaccountsKeeper.SetSubaccount(ctx, subaccount)
			subaccountsKeeper.AddSubaccountToSafetyHeap(
				ctx,
				*subaccount.Id,
				perpetualId,
				side,
			)

			require.Equal(
				t,
				uint32(i+1),
				subaccountsKeeper.GetSafetyHeapLength(store),
			)
		}

		// Make sure subaccounts are sorted correctly.
		for i := 0; i < totalSubaccounts; i++ {
			// Get the subaccount with the lowest safety score.
			// In this case, the subaccount with the lowest USDC balance.
			subaccountId := subaccountsKeeper.MustGetSubaccountAtIndex(store, uint32(0))
			subaccount := subaccountsKeeper.GetSubaccount(ctx, subaccountId)

			// Subaccounts should be sorted by asset position balance.
			require.Equal(t, uint32(i), subaccountId.Number)
			require.Equal(
				t,
				big.NewInt(int64(i-totalSubaccounts/2)),
				subaccount.GetUsdcPosition(),
			)

			// Remove the subaccount from the heap.
			subaccountsKeeper.RemoveSubaccountFromSafetyHeap(
				ctx,
				subaccountId,
				perpetualId,
				side,
			)
			require.Equal(
				t,
				uint32(totalSubaccounts-i-1),
				subaccountsKeeper.GetSafetyHeapLength(store),
			)
		}
	}
}

func TestSafetyHeapInsertRemoveIndex(t *testing.T) {
	perpetualId := uint32(0)
	side := satypes.Long
	totalSubaccounts := 100

	// Create 1000 subaccounts with balances ranging from -500 to 500.
	// The subaccounts should be sorted by balance.
	allSubaccounts := make([]satypes.Subaccount, 0)
	for i := 0; i < totalSubaccounts; i++ {
		subaccount := satypes.Subaccount{
			Id: &satypes.SubaccountId{
				Owner: sdk.MustBech32ifyAddressBytes(
					config.Bech32PrefixAccAddr,
					constants.AliceAccAddress,
				),
				Number: uint32(i),
			},
			AssetPositions: testutil.CreateUsdcAssetPositions(
				// Create asset positions with balances ranging from -500 to 500.
				big.NewInt(int64(i - totalSubaccounts/2)),
			),
		}

		// Handle special case.
		if i-totalSubaccounts/2 == 0 {
			subaccount.AssetPositions = nil
		}

		allSubaccounts = append(allSubaccounts, subaccount)
	}

	for iter := 0; iter < 100; iter++ {
		// Setup keeper state and test parameters.
		ctx, subaccountsKeeper, _, _, _, _, _, _, _, _, _ := keepertest.SubaccountsKeepers(t, false)

		// Shuffle the subaccounts so that insertion order is random.
		slices.Shuffle(allSubaccounts)

		store := subaccountsKeeper.GetSafetyHeapStore(ctx, perpetualId, side)
		for i, subaccount := range allSubaccounts {
			subaccountsKeeper.SetSubaccount(ctx, subaccount)
			subaccountsKeeper.AddSubaccountToSafetyHeap(
				ctx,
				*subaccount.Id,
				perpetualId,
				side,
			)

			require.Equal(
				t,
				uint32(i+1),
				subaccountsKeeper.GetSafetyHeapLength(store),
			)
		}

		for i := totalSubaccounts; i > 0; i-- {
			// Remove a random subaccount from the heap.
			index := rand.Intn(i)

			subaccountId := subaccountsKeeper.MustGetSubaccountAtIndex(store, uint32(index))
			subaccountsKeeper.RemoveSubaccountFromSafetyHeap(
				ctx,
				subaccountId,
				perpetualId,
				side,
			)

			require.Equal(
				t,
				uint32(i-1),
				subaccountsKeeper.GetSafetyHeapLength(store),
			)

			// Verify that the heap property is restored.
			verifyHeapProperties(t, subaccountsKeeper, ctx, store, 0)
		}
	}
}

func verifyHeapProperties(t *testing.T, k *keeper.Keeper, ctx sdk.Context, store prefix.Store, index uint32) {
	length := k.GetSafetyHeapLength(store)
	leftChild, rightChild := 2*index+1, 2*index+2

	if leftChild < length {
		require.True(t, k.Less(ctx, store, index, leftChild))
		verifyHeapProperties(t, k, ctx, store, leftChild)
	}

	if rightChild < length {
		require.True(t, k.Less(ctx, store, index, rightChild))
		verifyHeapProperties(t, k, ctx, store, rightChild)
	}
}
