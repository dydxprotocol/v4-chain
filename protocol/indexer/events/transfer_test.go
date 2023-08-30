package events_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/indexer/events"
	v1 "github.com/dydxprotocol/v4-chain/protocol/indexer/protocol/v1"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/require"
)

var (
	senderSubaccountId    = constants.Alice_Num0
	recipientSubaccountId = constants.Alice_Num1
	senderAddress         = constants.AliceAccAddress
	recipientAddress      = constants.BobAccAddress
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
	indexerSenderSubaccountId := v1.SubaccountIdToIndexerSubaccountId(senderSubaccountId)
	indexerRecipientSubaccountId := v1.SubaccountIdToIndexerSubaccountId(recipientSubaccountId)
	expectedTransferEventProto := &events.TransferEventV1{
		SenderSubaccountId:    &indexerSenderSubaccountId,
		RecipientSubaccountId: &indexerRecipientSubaccountId,
		Sender: &events.SourceOfFunds{
			Source: &events.SourceOfFunds_SubaccountId{
				SubaccountId: &indexerSenderSubaccountId,
			},
		},
		Recipient: &events.SourceOfFunds{
			Source: &events.SourceOfFunds_SubaccountId{
				SubaccountId: &indexerRecipientSubaccountId,
			},
		},
		AssetId: assetId,
		Amount:  amount.ToUint64(),
	}
	require.Equal(t, expectedTransferEventProto, transferEvent)
}

func TestNewDepositEvent_Success(t *testing.T) {
	depositEvent := events.NewDepositEvent(
		senderAddress.String(),
		recipientSubaccountId,
		assetId,
		amount,
	)
	indexerRecipientSubaccountId := v1.SubaccountIdToIndexerSubaccountId(recipientSubaccountId)
	expectedDepositEventProto := &events.TransferEventV1{
		Sender: &events.SourceOfFunds{
			Source: &events.SourceOfFunds_Address{
				Address: senderAddress.String(),
			},
		},
		Recipient: &events.SourceOfFunds{
			Source: &events.SourceOfFunds_SubaccountId{
				SubaccountId: &indexerRecipientSubaccountId,
			},
		},
		AssetId: assetId,
		Amount:  amount.ToUint64(),
	}
	require.Equal(t, expectedDepositEventProto, depositEvent)
}

func TestNewWithdrawEvent_Success(t *testing.T) {
	withdrawEvent := events.NewWithdrawEvent(
		senderSubaccountId,
		recipientAddress.String(),
		assetId,
		amount,
	)
	indexerSenderSubaccountId := v1.SubaccountIdToIndexerSubaccountId(senderSubaccountId)
	expectedWithdrawEventProto := &events.TransferEventV1{
		Sender: &events.SourceOfFunds{
			Source: &events.SourceOfFunds_SubaccountId{
				SubaccountId: &indexerSenderSubaccountId,
			},
		},
		Recipient: &events.SourceOfFunds{
			Source: &events.SourceOfFunds_Address{
				Address: recipientAddress.String(),
			},
		},
		AssetId: assetId,
		Amount:  amount.ToUint64(),
	}
	require.Equal(t, expectedWithdrawEventProto, withdrawEvent)
}
