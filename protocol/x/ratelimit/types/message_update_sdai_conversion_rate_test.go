package types_test

import (
	"testing"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/ratelimit/types"
	"github.com/stretchr/testify/require"
)

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
