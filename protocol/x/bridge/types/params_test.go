package types_test

import (
	"testing"
	"time"

	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
	"github.com/stretchr/testify/require"
)

func TestEventParams_Validate(t *testing.T) {
	tests := map[string]struct {
		params *types.EventParams
		err    string
	}{
		"Valid: default": {
			params: &types.DefaultGenesis().EventParams,
		},
		"Invalid: Eth Address": {
			params: &types.EventParams{
				Denom:      "denom",
				EthChainId: 1,
				EthAddress: "",
			},
			err: types.ErrInvalidEthAddress.Error(),
		},
		"Invalid: denom": {
			params: &types.EventParams{
				Denom:      "7coin",
				EthChainId: 1,
				EthAddress: "test",
			},
			err: "invalid denom",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := tc.params.Validate()
			if tc.err == "" {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, tc.err)
			}
		})
	}
}

func TestProposeParams_Validate(t *testing.T) {
	tests := map[string]struct {
		params *types.ProposeParams
		err    error
	}{
		"default is valid": {
			params: &types.DefaultGenesis().ProposeParams,
			err:    nil,
		},
		"empty is valid": {
			params: &types.ProposeParams{},
			err:    nil,
		},
		"negative ProposeDelayDuration": {
			params: &types.ProposeParams{
				ProposeDelayDuration: time.Duration(-1),
			},
			err: types.ErrNegativeDuration,
		},
		"negative SkipIfBlockDelayedByDuration": {
			params: &types.ProposeParams{
				SkipIfBlockDelayedByDuration: time.Duration(-1),
			},
			err: types.ErrNegativeDuration,
		},
		"too large of SkipRatePpm": {
			params: &types.ProposeParams{
				SkipRatePpm: lib.OneMillion + 1,
			},
			err: types.ErrNegativeDuration,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := tc.params.Validate()
			if tc.err == nil {
				require.NoError(t, err)
			} else {
				require.Error(t, tc.err, err)
			}
		})
	}
}

func TestSafetyParams_Validate(t *testing.T) {
	tests := map[string]struct {
		params *types.SafetyParams
		err    error
	}{
		"default is valid": {
			params: &types.DefaultGenesis().SafetyParams,
			err:    nil,
		},
		"empty is valid": {
			params: &types.SafetyParams{},
			err:    nil,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := tc.params.Validate()
			if tc.err == nil {
				require.NoError(t, err)
			} else {
				require.Error(t, tc.err, err)
			}
		})
	}
}
