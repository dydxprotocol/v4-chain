package types_test

import (
	"testing"

	types "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	"github.com/stretchr/testify/require"
)

func TestMsgUpdatePerpetualParams_ValidateBasic(t *testing.T) {
	tests := map[string]struct {
		msg         types.MsgUpdatePerpetualParams
		expectedErr string
	}{
		"Success": {
			msg: types.MsgUpdatePerpetualParams{
				Authority: validAuthority,
				PerpetualParams: types.PerpetualParams{
					Ticker:            "test",
					DefaultFundingPpm: 217,
					MarketType:        types.PerpetualMarketType_PERPETUAL_MARKET_TYPE_CROSS,
				},
			},
		},
		"Failure: Invalid authority": {
			msg: types.MsgUpdatePerpetualParams{
				Authority: "",
			},
			expectedErr: "Authority is invalid",
		},
		"Failure: Empty ticker": {
			msg: types.MsgUpdatePerpetualParams{
				Authority: validAuthority,
				PerpetualParams: types.PerpetualParams{
					Ticker:     "",
					MarketType: types.PerpetualMarketType_PERPETUAL_MARKET_TYPE_CROSS,
				},
			},
			expectedErr: "Ticker must be non-empty string",
		},
		"Failure: DefaultFundingPpm >= MaxDefaultFundingPpmAbs": {
			msg: types.MsgUpdatePerpetualParams{
				Authority: validAuthority,
				PerpetualParams: types.PerpetualParams{
					Ticker:            "test",
					DefaultFundingPpm: 100_000_000,
					MarketType:        types.PerpetualMarketType_PERPETUAL_MARKET_TYPE_CROSS,
				},
			},
			expectedErr: "DefaultFundingPpm magnitude exceeds maximum value",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := tc.msg.ValidateBasic()
			if tc.expectedErr == "" {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, tc.expectedErr)
			}
		})
	}
}
