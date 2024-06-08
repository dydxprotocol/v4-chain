package v_3_0_0

import (
	"github.com/stretchr/testify/require"

	"testing"
)

func TestICAHostAllowMessages(t *testing.T) {
	require.Equal(
		t,
		[]string{
			"/ibc.applications.transfer.v1.MsgTransfer",
			"/cosmos.bank.v1beta1.MsgSend",
		},
		ICAHostAllowMessages,
	)
}
