package ante_test

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"

	appmsgs "github.com/StreamFinance-Protocol/stream-chain/protocol/app/msgs"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/lib/ante"
	testmsgs "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/msgs"
	"github.com/stretchr/testify/require"
)

func TestIsDisallowExternalSubmitMsg(t *testing.T) {
	// All disallow msgs should return true.
	disallowSampleMsgs := testmsgs.GetNonNilSampleMsgs(appmsgs.DisallowMsgs)
	fmt.Println("DISALLOW MSG SAMPLES")
	fmt.Println(disallowSampleMsgs)
	for _, sampleMsg := range disallowSampleMsgs {
		result := ante.IsDisallowExternalSubmitMsg(sampleMsg.Msg)
		fmt.Println("SAMPLE MSG")
		fmt.Println(sampleMsg.Name)
		fmt.Println(result)
		if ante.IsNestedMsg(sampleMsg.Msg) {
			// nested msgs are allowed as long as the inner msgs are allowed.
			require.False(t, result, sampleMsg.Name)
		} else {
			require.True(t, result, sampleMsg.Name)
		}
	}

	// All allow msgs should return false.
	allowSampleMsgs := testmsgs.GetNonNilSampleMsgs(appmsgs.AllowMsgs)
	fmt.Println("ALLOW MSG SAMPLES")
	fmt.Println(allowSampleMsgs)
	require.NotZero(t, len(allowSampleMsgs)) // checking just not zero is sufficient.
	for _, sampleMsg := range allowSampleMsgs {
		require.False(t, ante.IsDisallowExternalSubmitMsg(sampleMsg.Msg), sampleMsg.Name)
	}
}

func TestIsDisallowExternalSubmitMsg_InvalidInnerMsgs(t *testing.T) {
	containsInvalidInnerMsgs := []sdk.Msg{}

	for _, msg := range containsInvalidInnerMsgs {
		require.True(t, ante.IsDisallowExternalSubmitMsg(msg))
	}
}
