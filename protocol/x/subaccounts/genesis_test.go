package subaccounts_test

import (
	"math/big"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"

	keepertest "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/keeper"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/nullify"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/subaccounts"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/subaccounts/keeper"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Subaccounts: []types.Subaccount{
			{
				Id: &types.SubaccountId{
					Owner:  "foo",
					Number: uint32(0),
				},
				AssetPositions:  keepertest.CreateTDaiAssetPosition(big.NewInt(1_000)),
				AssetYieldIndex: big.NewRat(1, 1).String(),
			},
			{
				Id: &types.SubaccountId{
					Owner:  "bar",
					Number: uint32(99),
				},
				AssetPositions:  keepertest.CreateTDaiAssetPosition(big.NewInt(1_000)),
				AssetYieldIndex: big.NewRat(1, 1).String(),
			},
			{
				Id: &types.SubaccountId{
					Owner:  "bar",
					Number: uint32(101),
				},
				AssetPositions:  keepertest.CreateTDaiAssetPosition(big.NewInt(100)),
				AssetYieldIndex: big.NewRat(2, 1).String(),
			},
		},
	}

	ctx, k, _, _, _, _, _, _, _, _ := keepertest.SubaccountsKeepers(t, true)
	subaccounts.InitGenesis(ctx, *k, genesisState)
	assertSubaccountUpdateEventsInIndexerBlock(t, k, ctx, 3)
	got := subaccounts.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState) //nolint:staticcheck
	nullify.Fill(got)           //nolint:staticcheck

	require.ElementsMatch(t, genesisState.Subaccounts, got.Subaccounts)
}

// assertSubaccountUpdateEventsInIndexerBlock checks that the number of subaccount update events
// included in the Indexer block kafka message.
func assertSubaccountUpdateEventsInIndexerBlock(
	t *testing.T,
	k *keeper.Keeper,
	ctx sdk.Context,
	numSubaccounts int,
) {
	subaccountUpdates := keepertest.GetSubaccountUpdateEventsFromIndexerBlock(ctx, k)
	require.Len(t, subaccountUpdates, numSubaccounts)
}
