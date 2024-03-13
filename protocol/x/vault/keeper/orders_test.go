package keeper_test

import (
	"fmt"
	"testing"

	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
	"github.com/stretchr/testify/require"
)

func TestConstructVaultClobOrderss(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.VaultKeeper

	fmt.Println("vault orders", k.ConstructVaultClobOrders(
		ctx,
		types.VaultId{
			Type:   types.VaultType_VAULT_TYPE_CLOB,
			Number: 0,
		},
		tApp.App.ClobKeeper.GetAllClobPairs(ctx)[0],
		tApp.App.PerpetualsKeeper.GetAllPerpetuals(ctx)[0],
		tApp.App.PricesKeeper.GetAllMarketPrices(ctx)[0],
	))

	require.True(t, false)
}
