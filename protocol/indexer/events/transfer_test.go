package events_test

import (
	"testing"

	"github.com/dydxprotocol/v4/indexer/events"
	"github.com/dydxprotocol/v4/indexer/protocol/v1"
	"github.com/dydxprotocol/v4/testutil/constants"
	satypes "github.com/dydxprotocol/v4/x/subaccounts/types"
	"github.com/stretchr/testify/require"
)

var (
	senderSubaccountId    = constants.Alice_Num0
	recipientSubaccountId = constants.Alice_Num1
	amount                = satypes.BaseQuantums(5)
	assetId               = uint32(0)
)

func TestNewTransferEvent_Success(t *testing.T) {
	transferEvent := events.NewTransferEvent(
		senderSubaccountId,
		recipientSubaccountId,
		assetId,
		amount,
	)
	expectedTransferEventProto := &events.TransferEventV1{
		SenderSubaccountId:    v1.SubaccountIdToIndexerSubaccountId(senderSubaccountId),
		RecipientSubaccountId: v1.SubaccountIdToIndexerSubaccountId(recipientSubaccountId),
		AssetId:               assetId,
		Amount:                amount.ToUint64(),
	}
	require.Equal(t, expectedTransferEventProto, transferEvent)
}
