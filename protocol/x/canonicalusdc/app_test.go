package canonicalusdc_test

import (
	"testing"

	"github.com/cometbft/cometbft/types"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	assettypes "github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
	canonicalusdctypes "github.com/dydxprotocol/v4-chain/protocol/x/canonicalusdc/types"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

func TestAppConvertsBaseAccountAtModuleAddress(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(func() types.GenesisDoc {
		genesis := testapp.DefaultGenesis()
		testapp.UpdateGenesisDocWithAppStateForModule(&genesis, func(state *authtypes.GenesisState) {
			accounts, err := authtypes.UnpackAccounts(state.Accounts)
			require.NoError(t, err)
			accounts = append(accounts, authtypes.NewBaseAccountWithAddress(canonicalusdctypes.ModuleAddress))
			state.Accounts, err = authtypes.PackAccounts(authtypes.SanitizeGenesisAccounts(accounts))
			require.NoError(t, err)
		})
		return genesis
	}).Build()
	ctx := tApp.InitChain()

	account := tApp.App.AccountKeeper.GetAccount(ctx, canonicalusdctypes.ModuleAddress)
	moduleAccount, ok := account.(sdk.ModuleAccountI)
	require.True(t, ok, "got %T", account)
	require.Equal(t, canonicalusdctypes.ModuleName, moduleAccount.GetName())
	require.True(t, moduleAccount.HasPermission(authtypes.Minter))
	require.True(t, moduleAccount.HasPermission(authtypes.Burner))

	coins := sdk.NewCoins(sdk.NewInt64Coin(assettypes.UusdcDenom, 1))
	require.NoError(t, tApp.App.BankKeeper.MintCoins(ctx, canonicalusdctypes.ModuleName, coins))
	require.NoError(t, tApp.App.BankKeeper.BurnCoins(ctx, canonicalusdctypes.ModuleName, coins))
}
