package rewards_test

import (
	"testing"

	testapp "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/app"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/rewards"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/rewards/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),
	}

	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.RewardsKeeper

	rewards.InitGenesis(ctx, k, genesisState)
	got := rewards.ExportGenesis(ctx, k)
	require.NotNil(t, got)
}
