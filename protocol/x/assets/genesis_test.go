package assets_test

import (
	"testing"

	keepertest "github.com/dydxprotocol/v4/testutil/keeper"
	"github.com/dydxprotocol/v4/x/assets"
	"github.com/dydxprotocol/v4/x/assets/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	expected := types.DefaultGenesis()
	ctx, k, _, _, _, _ := keepertest.AssetsKeepers(t)
	assets.InitGenesis(ctx, *k, *expected)
	actual := assets.ExportGenesis(ctx, *k)
	require.NotNil(t, actual)
	require.ElementsMatch(t, actual.Assets, expected.Assets)
}
