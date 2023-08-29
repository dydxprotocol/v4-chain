package delaymsg_test

import (
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	sdktest "github.com/dydxprotocol/v4-chain/protocol/testutil/sdk"
	"github.com/dydxprotocol/v4-chain/protocol/x/delaymsg"
	"testing"
	"time"
)

func TestEndBlocker(t *testing.T) {
	k := &mocks.DelayMsgKeeper{}
	ctx := sdktest.NewContextWithBlockHeightAndTime(0, time.Now())
	k.On("DispatchMessagesForBlock", ctx).Return().Once()
	delaymsg.EndBlocker(ctx, k)
	k.AssertExpectations(t)
}
