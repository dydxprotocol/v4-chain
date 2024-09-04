package types_test

import (
	"testing"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/ratelimit/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestNewMsgUpdateSDAIConversionRate(t *testing.T) {
	tests := []struct {
		name           string
		sender         string
		conversionRate string
		expectedMsg    types.MsgUpdateSDAIConversionRate
	}{
		{
			name:           "Valid message",
			sender:         "cosmos139f7kncmglres2nf3h4hc4tade85ekfr8sulz5",
			conversionRate: "1",
			expectedMsg: types.MsgUpdateSDAIConversionRate{
				Sender:         "cosmos139f7kncmglres2nf3h4hc4tade85ekfr8sulz5",
				ConversionRate: "1",
			},
		},
		{
			name:           "Different conversion rate",
			sender:         "cosmos139f7kncmglres2nf3h4hc4tade85ekfr8sulz5",
			conversionRate: "123456789123456789123456789",
			expectedMsg: types.MsgUpdateSDAIConversionRate{
				Sender:         "cosmos139f7kncmglres2nf3h4hc4tade85ekfr8sulz5",
				ConversionRate: "123456789123456789123456789",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			senderAddr, err := sdk.AccAddressFromBech32(tc.sender)
			require.NoError(t, err)

			msg := types.NewMsgUpdateSDAIConversionRate(senderAddr, tc.conversionRate)
			require.Equal(t, &tc.expectedMsg, msg)
		})
	}
}

func TestMsgUpdateSDAIConversionRate_ValidateBasic(t *testing.T) {
	tests := map[string]struct {
		msg types.MsgUpdateSDAIConversionRate
		err bool
	}{
		"Valid": {
			msg: types.MsgUpdateSDAIConversionRate{
				Sender:         "cosmos139f7kncmglres2nf3h4hc4tade85ekfr8sulz5",
				ConversionRate: "1",
			},
			err: false,
		},
		"Invalid: empty sender": {
			msg: types.MsgUpdateSDAIConversionRate{
				Sender:         "",
				ConversionRate: "1",
			},
			err: true,
		},
		"Invalid: incorrect sender": {
			msg: types.MsgUpdateSDAIConversionRate{
				Sender:         "incorrect_sender",
				ConversionRate: "1",
			},
			err: true,
		},
		"Invalid: empty conversion rate": {
			msg: types.MsgUpdateSDAIConversionRate{
				Sender:         "dydx1eeeggku6dzk3mv7wph3zq035rhtd890smfq5z6",
				ConversionRate: "",
			},
			err: true,
		},
		"Invalid: negative conversion rate": {
			msg: types.MsgUpdateSDAIConversionRate{
				Sender:         "dydx1eeeggku6dzk3mv7wph3zq035rhtd890smfq5z6",
				ConversionRate: "-1",
			},
			err: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}
