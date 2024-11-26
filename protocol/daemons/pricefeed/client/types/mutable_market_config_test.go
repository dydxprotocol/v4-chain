package types_test

import (
	"errors"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/types"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	mmc := &types.MutableMarketConfig{
		Id:           constants.MarketId7,
		Pair:         "ABC-USD",
		Exponent:     -5,
		MinExchanges: 2,
	}

	mmcCopy := mmc.Copy()

	require.NotSame(t, mmc, mmcCopy)
	require.Equal(t, mmc, mmcCopy)
}

func TestValidate_Mixed(t *testing.T) {
	tests := map[string]struct {
		mmc           *types.MutableMarketConfig
		expectedError error
	}{
		"valid": {
			mmc: &types.MutableMarketConfig{
				Id:           constants.MarketId7,
				Pair:         "ABC-USD",
				Exponent:     -5,
				MinExchanges: 2,
			},
		},
		"invalid: empty pair": {
			mmc: &types.MutableMarketConfig{
				Id:           constants.MarketId7,
				Pair:         "",
				Exponent:     -5,
				MinExchanges: 2,
			},
			expectedError: errors.New("pair cannot be empty"),
		},
		"invalid: min exchanges 0": {
			mmc: &types.MutableMarketConfig{
				Id:           constants.MarketId7,
				Pair:         "ABC-USD",
				Exponent:     -5,
				MinExchanges: 0,
			},
			expectedError: errors.New("min exchanges cannot be 0"),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := tc.mmc.Validate()
			if tc.expectedError == nil {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tc.expectedError.Error())
			}
		})
	}
}
