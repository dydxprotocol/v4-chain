package keeper_test

import (
	"testing"
	"time"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/dtypes"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/lib"
	testapp "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/app"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/ratelimit/keeper"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/ratelimit/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

var (
	altLimitParamsForSDai = types.LimitParams{
		Denom: types.SDaiDenom,
		Limiters: []types.Limiter{
			{
				Period:          3600 * time.Second,
				BaselineMinimum: dtypes.NewInt(1000),
				BaselineTvlPpm:  1000,
			},
		},
	}
	stakeLimitParams = types.LimitParams{
		Denom: lib.DefaultBaseDenom,
		Limiters: []types.Limiter{
			{
				Period:          3600 * time.Second,
				BaselineMinimum: dtypes.NewInt(1000),
				BaselineTvlPpm:  1000,
			}, {
				Period:          72 * time.Hour,
				BaselineMinimum: dtypes.NewInt(1000),
				BaselineTvlPpm:  1000,
			},
		},
	}
)

func TestMsgSetLimitParams(t *testing.T) {
	testCases := []struct {
		name                    string
		input                   *types.MsgSetLimitParams
		expectedLimitParamsList []types.LimitParams
		expErr                  bool
		expErrMsg               string
	}{
		{
			name: "Overwite default params with default params (no-op)",
			input: &types.MsgSetLimitParams{
				Authority:   lib.GovModuleAddress.String(),
				LimitParams: types.DefaultSDaiRateLimitParams(),
			},
			expectedLimitParamsList: []types.LimitParams{types.DefaultSDaiRateLimitParams()},
			expErr:                  false,
		},
		{
			name: "Overwrite default params",
			input: &types.MsgSetLimitParams{
				Authority:   lib.GovModuleAddress.String(),
				LimitParams: altLimitParamsForSDai,
			},
			expectedLimitParamsList: []types.LimitParams{altLimitParamsForSDai},
			expErr:                  false,
		},
		{
			name: "Add additional denom",
			input: &types.MsgSetLimitParams{
				Authority:   lib.GovModuleAddress.String(),
				LimitParams: stakeLimitParams,
			},
			expectedLimitParamsList: []types.LimitParams{
				stakeLimitParams,
				types.DefaultSDaiRateLimitParams(),
			},
			expErr: false,
		},
		{
			name: "Remove rate-limit for SDai",
			input: &types.MsgSetLimitParams{
				Authority: lib.GovModuleAddress.String(),
				LimitParams: types.LimitParams{
					Denom:    types.SDaiDenom,
					Limiters: []types.Limiter{}, // Empty list removes rate-limit
				},
			},
			expectedLimitParamsList: nil,
			expErr:                  false,
		},
		{
			name: "invalid authority",
			input: &types.MsgSetLimitParams{
				Authority:   "invalid",
				LimitParams: types.DefaultSDaiRateLimitParams(),
			},
			expErr:    true,
			expErrMsg: "invalid authority",
		},
		{
			name: "invalid params: invalid denom",
			input: &types.MsgSetLimitParams{
				Authority: lib.GovModuleAddress.String(),
				LimitParams: types.LimitParams{
					Denom: "",
					Limiters: []types.Limiter{
						{
							Period:          3600,
							BaselineMinimum: dtypes.NewInt(1000),
							BaselineTvlPpm:  1000,
						},
					},
				},
			},
			expErr:    true,
			expErrMsg: "invalid denom",
		},
		{
			name: "invalid params: negative baseline minimum",
			input: &types.MsgSetLimitParams{
				Authority: lib.GovModuleAddress.String(),
				LimitParams: types.LimitParams{
					Denom: "denom",
					Limiters: []types.Limiter{
						{
							Period:          3600,
							BaselineMinimum: dtypes.NewInt(-1000), // -1000, must be positive
							BaselineTvlPpm:  1000,
						},
					},
				},
			},
			expErr:    true,
			expErrMsg: types.ErrInvalidBaselineMinimum.Error(),
		},
		{
			name: "invalid params: negative baseline minimum",
			input: &types.MsgSetLimitParams{
				Authority: lib.GovModuleAddress.String(),
				LimitParams: types.LimitParams{
					Denom: "denom",
					Limiters: []types.Limiter{
						{
							Period:          3600,
							BaselineMinimum: dtypes.NewInt(1000),
							BaselineTvlPpm:  1_000_100, // 100.01%, must be in (0%, 100%)
						},
					},
				},
			},
			expErr:    true,
			expErrMsg: types.ErrInvalidBaselineTvlPpm.Error(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).Build()
			ctx := tApp.InitChain()
			k := tApp.App.RatelimitKeeper

			ms := keeper.NewMsgServerImpl(k)

			_, err := ms.SetLimitParams(ctx, tc.input)
			if tc.expErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expErrMsg)
			} else {
				require.NoError(t, err)
				sdkCtx := sdk.UnwrapSDKContext(ctx)
				require.Equal(t,
					tc.expectedLimitParamsList,
					k.GetAllLimitParams(sdkCtx),
				)
			}
		})
	}
}
