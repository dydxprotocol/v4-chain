package keeper_test

import (
	"testing"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	"github.com/stretchr/testify/require"
)

func TestAllLiquidityTiers(
	t *testing.T,
) {
	tApp := testapp.NewTestAppBuilder(t).
		WithGenesisDocFn(func() types.GenesisDoc {
			genesis := testapp.DefaultGenesis()
			testapp.UpdateGenesisDocWithAppStateForModule(&genesis, func(state *perptypes.GenesisState) {
				state.LiquidityTiers = constants.LiquidityTiers
			})
			return genesis
		}).Build()

	tApp.InitChain()

	request := perptypes.QueryAllLiquidityTiersRequest{}
	abciResponse := tApp.App.Query(abci.RequestQuery{
		Path: "/dydxprotocol.perpetuals.Query/AllLiquidityTiers",
		Data: tApp.App.AppCodec().MustMarshal(&request),
	})
	require.True(t, abciResponse.IsOK())

	var actual perptypes.QueryAllLiquidityTiersResponse
	tApp.App.AppCodec().MustUnmarshal(abciResponse.Value, &actual)

	expected := perptypes.QueryAllLiquidityTiersResponse{
		LiquidityTiers: constants.LiquidityTiers,
		Pagination: &query.PageResponse{
			NextKey: nil,
			Total:   10,
		},
	}
	require.Equal(t, expected, actual)
}
