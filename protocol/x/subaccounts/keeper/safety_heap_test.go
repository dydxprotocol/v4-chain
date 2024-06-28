package keeper_test

import (
	"math/big"
	"math/rand"
	"testing"

	"github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/app/config"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	testutil "github.com/dydxprotocol/v4-chain/protocol/testutil/util"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/require"
)

func TestSafetyHeapInsertRemoval(t *testing.T) {

	allSubaccounts := make([]satypes.Subaccount, 0)
	for i := 0; i < 1000; i++ {
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
				big.NewInt(int64(i - 500)),
			),
		}

		// Handle special case.
		if i == 500 {
			subaccount.AssetPositions = nil
		}

		allSubaccounts = append(allSubaccounts, subaccount)
	}

	for iter := 0; iter < 100; iter++ {
		// Setup keeper state and test parameters.
		ctx, subaccountsKeeper, _, _, _, _, _, _, _, _ := keepertest.SubaccountsKeepers(t, false)

		rand.Shuffle(len(allSubaccounts), func(i, j int) {
			allSubaccounts[i], allSubaccounts[j] = allSubaccounts[j], allSubaccounts[i]
		})

		store := subaccountsKeeper.GetSafetyHeapStore(ctx, 0, satypes.Long)
		for _, subaccount := range allSubaccounts {
			subaccountsKeeper.SetSubaccount(ctx, subaccount)
			subaccountsKeeper.Insert(ctx, store, *subaccount.Id)
		}

		// Make sure subaccounts are sorted correctly.
		for i := 0; i < 1000; i++ {
			subaccountId := subaccountsKeeper.MustGetSubaccountAtIndex(store, uint32(0))

			// Subaccounts should be sorted by asset position balance.
			require.Equal(t, uint32(i), subaccountId.Number)

			// Remove the subaccount from the heap.
			subaccountsKeeper.MustRemoveElementAtIndex(ctx, store, 0)
		}
	}
}
