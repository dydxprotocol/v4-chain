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
			"/cosmos.staking.v1beta1.MsgDelegate",
			"/cosmos.staking.v1beta1.MsgBeginRedelegate",
			"/cosmos.staking.v1beta1.MsgUndelegate",
			"/cosmos.staking.v1beta1.MsgCancelUnbondingDelegation",
			"/cosmos.distribution.v1beta1.MsgSetWithdrawAddress",
			"/cosmos.distribution.v1beta1.MsgWithdrawDelegatorReward",
			"/cosmos.distribution.v1beta1.MsgFundCommunityPool",
			"/cosmos.gov.v1.MsgVote",
		},
		ICAHostAllowMessages,
	)
}
