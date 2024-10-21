package types_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/x/prices/client/testutil"

	types "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	"github.com/stretchr/testify/require"
)

func TestMsgCreateOracleMarket_ValidateBasic(t *testing.T) {
	tests := []struct {
		desc        string
		msg         types.MsgCreateOracleMarket
		expectedErr string
	}{
		{
			desc: "Empty authority",
			msg: types.MsgCreateOracleMarket{
				Authority: "",
			},
			expectedErr: "authority '' must be a valid bech32 address, but got error 'empty address string is not " +
				"allowed': Authority is invalid",
		},
		{
			desc: "Malformatted authority",
			msg: types.MsgCreateOracleMarket{
				Authority: "invalid",
			},
			expectedErr: "authority 'invalid' must be a valid bech32 address, but got error 'decoding bech32 " +
				"failed: invalid bech32 string length 7': Authority is invalid",
		},
		{
			desc: "Valid MsgCreateOracleMarket",
			msg: types.MsgCreateOracleMarket{
				Authority: testutil.ValidAuthority,
				Params: types.MarketParam{
					Pair:              "BTC-USD",
					MinPriceChangePpm: 1_000,
				},
			},
			expectedErr: "",
		},
		{
			desc: "Empty pair",
			msg: types.MsgCreateOracleMarket{
				Authority: testutil.ValidAuthority,
				Params: types.MarketParam{
					Pair:              "",
					MinPriceChangePpm: 1_000,
				},
			},
			expectedErr: "Pair cannot be empty",
		},
		{
			desc: "Invalid MinPriceChangePpm",
			msg: types.MsgCreateOracleMarket{
				Authority: testutil.ValidAuthority,
				Params: types.MarketParam{
					Pair:              "BTC-USD",
					MinPriceChangePpm: 0,
				},
			},
			expectedErr: "Min price change in parts-per-million must be greater than 0",
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.msg.ValidateBasic()
			if tc.expectedErr == "" {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, tc.expectedErr)
			}
		})
	}
}
