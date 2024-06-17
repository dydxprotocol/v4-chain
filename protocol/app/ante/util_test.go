package ante_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/app/ante"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/stretchr/testify/require"
)

type ValidationUtilTestCase struct {
	msgs                 []sdk.Msg
	shouldSkipValidation bool
}

func TestSkipSequenceValidation(t *testing.T) {
	testCases := map[string]ValidationUtilTestCase{
		"single place order message": {
			msgs: []sdk.Msg{
				constants.Msg_PlaceOrder,
			},
			shouldSkipValidation: true,
		},
		"single cancel order message": {
			msgs: []sdk.Msg{
				constants.Msg_CancelOrder,
			},
			shouldSkipValidation: true,
		},
		"single batch cancel message": {
			msgs: []sdk.Msg{
				constants.Msg_BatchCancel,
			},
			shouldSkipValidation: true,
		},
		"single transfer message": {
			msgs: []sdk.Msg{
				constants.Msg_Transfer,
			},
			shouldSkipValidation: false,
		},
		"single send message": {
			msgs: []sdk.Msg{
				constants.Msg_Send,
			},
			shouldSkipValidation: false,
		},
		"single long term order": {
			msgs: []sdk.Msg{
				constants.Msg_PlaceOrder_LongTerm,
			},
			shouldSkipValidation: false,
		},
		"single long term cancel": {
			msgs: []sdk.Msg{
				constants.Msg_CancelOrder_LongTerm,
			},
			shouldSkipValidation: false,
		},
		"single conditional order": {
			msgs: []sdk.Msg{
				constants.Msg_PlaceOrder_Conditional,
			},
			shouldSkipValidation: false,
		},
		"multiple GTB messages": {
			msgs: []sdk.Msg{
				constants.Msg_PlaceOrder,
				constants.Msg_CancelOrder,
			},
			shouldSkipValidation: true,
		},
		"mix of GTB messages and non-GTB messages": {
			msgs: []sdk.Msg{
				constants.Msg_Transfer,
				constants.Msg_Send,
			},
			shouldSkipValidation: false,
		},
		"mix of long term orders and short term orders": {
			msgs: []sdk.Msg{
				constants.Msg_PlaceOrder,
				constants.Msg_PlaceOrder_LongTerm,
			},
			shouldSkipValidation: false,
		},
		"mix of conditional orders and short term orders": {
			msgs: []sdk.Msg{
				constants.Msg_PlaceOrder,
				constants.Msg_PlaceOrder_Conditional,
			},
			shouldSkipValidation: false,
		},
		"mix of conditional orders and short term batch cancel orders": {
			msgs: []sdk.Msg{
				constants.Msg_BatchCancel,
				constants.Msg_PlaceOrder_Conditional,
			},
			shouldSkipValidation: false,
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			require.Equal(
				t,
				tc.shouldSkipValidation,
				ante.ShouldSkipSequenceValidation(tc.msgs),
			)
		})
	}
}
