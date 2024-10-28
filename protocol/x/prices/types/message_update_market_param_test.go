package types_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/x/prices/client/testutil"

	"github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	"github.com/stretchr/testify/require"
)

func TestMsgUpdateMarketParam_ValidateBasic(t *testing.T) {
	tests := map[string]struct {
		msg         types.MsgUpdateMarketParam
		expectedErr string
	}{
		"Success": {
			msg: types.MsgUpdateMarketParam{
				Authority: testutil.ValidAuthority,
				MarketParam: types.MarketParam{
					Pair:              "test",
					MinPriceChangePpm: 1_000,
				},
			},
		},
		"Failure: Empty authority": {
			msg: types.MsgUpdateMarketParam{
				Authority: "",
			},
			expectedErr: "authority '' must be a valid bech32 address, but got error 'empty address string is not " +
				"allowed': Authority is invalid",
		},
		"Failure: Invalid authority": {
			msg: types.MsgUpdateMarketParam{
				Authority: "invalid",
			},
			expectedErr: "authority 'invalid' must be a valid bech32 address, but got error 'decoding bech32 " +
				"failed: invalid bech32 string length 7': Authority is invalid",
		},
		"Failure: Empty pair": {
			msg: types.MsgUpdateMarketParam{
				Authority: testutil.ValidAuthority,
				MarketParam: types.MarketParam{
					Pair:              "",
					MinPriceChangePpm: 1_000,
				},
			},
			expectedErr: "Pair cannot be empty",
		},
		"Failure: 0 MinPriceChangePpm": {
			msg: types.MsgUpdateMarketParam{
				Authority: testutil.ValidAuthority,
				MarketParam: types.MarketParam{
					Pair:              "test",
					MinPriceChangePpm: 0,
				},
			},
			expectedErr: "Min price change in parts-per-million must be greater than 0",
		},
		"Failure: 10_000 MinPriceChangePpm": {
			msg: types.MsgUpdateMarketParam{
				Authority: testutil.ValidAuthority,
				MarketParam: types.MarketParam{
					Pair:              "test",
					MinPriceChangePpm: 10_000,
				},
			},
			expectedErr: "Min price change in parts-per-million must be greater than 0 and less than 10000",
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
