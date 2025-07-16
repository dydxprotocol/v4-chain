package types_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	types "github.com/dydxprotocol/v4-chain/protocol/x/revshare/types"
	"github.com/stretchr/testify/require"
)

var (
	validAuthority = constants.AliceAccAddress.String()
)

func TestMsgSetOrderRouterRevShare_ValidateBasic(t *testing.T) {
	tests := map[string]struct {
		msg types.MsgSetOrderRouterRevShare
		err error
	}{
		"Valid": {
			msg: types.MsgSetOrderRouterRevShare{
				Authority: validAuthority,
				OrderRouterRevShare: types.OrderRouterRevShare{
					Address:  constants.AliceAccAddress.String(),
					SharePpm: 300_000,
				},
			},
		},
		"Invalid ppm": {
			msg: types.MsgSetOrderRouterRevShare{
				Authority: validAuthority,
				OrderRouterRevShare: types.OrderRouterRevShare{
					Address:  constants.AliceAccAddress.String(),
					SharePpm: 1_000_000,
				},
			},
			err: types.ErrInvalidRevenueSharePpm,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}
