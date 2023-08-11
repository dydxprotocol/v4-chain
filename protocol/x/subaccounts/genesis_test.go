package subaccounts_test

import (
	"testing"

	keepertest "github.com/dydxprotocol/v4/testutil/keeper"
	"github.com/dydxprotocol/v4/testutil/nullify"
	"github.com/dydxprotocol/v4/x/subaccounts"
	"github.com/dydxprotocol/v4/x/subaccounts/types"
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
			},
			{
				Id: &types.SubaccountId{
					Owner:  "bar",
					Number: uint32(99),
				},
			},
		},
	}

	ctx, k, _, _, _, _, _, _ := keepertest.SubaccountsKeepers(t, true)
	subaccounts.InitGenesis(ctx, *k, genesisState)
	got := subaccounts.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState) //nolint:staticcheck
	nullify.Fill(got)           //nolint:staticcheck

	require.ElementsMatch(t, genesisState.Subaccounts, got.Subaccounts)
}
