package keeper_test

import (
	"testing"

	abci "github.com/cometbft/cometbft/abci/types"
	testApp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	"github.com/stretchr/testify/require"
)

func TestEquityTierLimitConfiguration(
	t *testing.T,
) {
	tApp := testApp.NewTestAppBuilder().WithTesting(t).Build()
	ctx := tApp.InitChain()
	expected := types.QueryEquityTierLimitConfigurationResponse{
		EquityTierLimitConfig: tApp.App.ClobKeeper.GetEquityTierLimitConfiguration(ctx),
	}

	request := types.QueryEquityTierLimitConfigurationRequest{}
	abciResponse := tApp.App.Query(abci.RequestQuery{
		Path: "/dydxprotocol.clob.Query/EquityTierLimitConfiguration",
		Data: tApp.App.AppCodec().MustMarshal(&request),
	})
	require.True(t, abciResponse.IsOK())

	var actual types.QueryEquityTierLimitConfigurationResponse
	tApp.App.AppCodec().MustUnmarshal(abciResponse.Value, &actual)
	require.Equal(t, expected, actual)
}
