package keeper_test

import (
	"testing"

	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/x/revshare/types"
	"github.com/stretchr/testify/require"
)

func TestQueryUnconditionalRevShare(t *testing.T) {
	testCases := map[string]struct {
		config types.UnconditionalRevShareConfig
	}{
		"Single recipient": {
			config: types.UnconditionalRevShareConfig{
				Configs: []types.UnconditionalRevShareConfig_RecipientConfig{
					{
						Address:  constants.AliceAccAddress.String(),
						SharePpm: 100_000,
					},
				},
			},
		},
		"Multiple recipients": {
			config: types.UnconditionalRevShareConfig{
				Configs: []types.UnconditionalRevShareConfig_RecipientConfig{
					{
						Address:  constants.AliceAccAddress.String(),
						SharePpm: 50_000,
					},
					{
						Address:  constants.BobAccAddress.String(),
						SharePpm: 30_000,
					},
				},
			},
		},
		"Empty config": {
			config: types.UnconditionalRevShareConfig{},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).Build()
			ctx := tApp.InitChain()
			k := tApp.App.RevShareKeeper

			k.SetUnconditionalRevShareConfigParams(ctx, tc.config)

			resp, err := k.UnconditionalRevShareConfig(ctx, &types.QueryUnconditionalRevShareConfig{})
			require.NoError(t, err)
			require.Equal(t, tc.config, resp.Config)
		})
	}
}
