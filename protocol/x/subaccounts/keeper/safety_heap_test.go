package keeper_test

import (
	"math/big"
	"testing"

	"github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/app/config"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	testutil "github.com/dydxprotocol/v4-chain/protocol/testutil/util"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/require"
	"gopkg.in/typ.v4/slices"
)

func TestSafetyHeapInsertRemoval(t *testing.T) {
	totalSubaccounts := 1000

	allSubaccounts := make([]satypes.Subaccount, 0)
	for i := 0; i < totalSubaccounts; i++ {
		subaccount := satypes.Subaccount{
			Id: &satypes.SubaccountId{
				Owner: types.MustBech32ifyAddressBytes(
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
		ctx, subaccountsKeeper, _, _, _, _, _, _, _, _ := keepertest.SubaccountsKeepers(t, false)

		// Shuffle the subaccounts so that insertion order is random.
		slices.Shuffle(allSubaccounts)

		store := subaccountsKeeper.GetSafetyHeapStore(ctx, 0, satypes.Long)
		for i, subaccount := range allSubaccounts {
			subaccountsKeeper.SetSubaccount(ctx, subaccount)
			subaccountsKeeper.AddSubaccountToSafetyHeap(
				ctx,
				*subaccount.Id,
				0,
				satypes.Long,
			)

			require.Equal(
				t,
				uint32(i+1),
				subaccountsKeeper.GetSafetyHeapLength(store),
			)
		}

		// Make sure subaccounts are sorted correctly.
		for i := 0; i < totalSubaccounts; i++ {
			subaccountId := subaccountsKeeper.MustGetSubaccountAtIndex(store, uint32(0))

			// Subaccounts should be sorted by asset position balance.
			require.Equal(t, uint32(i), subaccountId.Number)

			// Remove the subaccount from the heap.
			subaccountsKeeper.RemoveSubaccountFromSafetyHeap(
				ctx,
				subaccountId,
				0,
				satypes.Long,
			)
			require.Equal(
				t,
				uint32(totalSubaccounts-i-1),
				subaccountsKeeper.GetSafetyHeapLength(store),
			)
		}
	}
}
