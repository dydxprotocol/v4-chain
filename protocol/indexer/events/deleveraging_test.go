package events_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/indexer/events"
	v1 "github.com/dydxprotocol/v4-chain/protocol/indexer/protocol/v1"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/stretchr/testify/require"
)

var (
	liquidatedSubaccountId = constants.Alice_Num0
	offsettingSubaccountId = constants.Bob_Num0
	clobPairId             = uint32(1)
	subticks               = uint64(1000)
	isBuy                  = true
)

func TestNewDeleveragingEvent_Success(t *testing.T) {
	deleveragingEvent := events.NewDeleveragingEvent(
		liquidatedSubaccountId,
		offsettingSubaccountId,
		clobPairId,
		fillAmount,
		subticks,
		isBuy,
	)
	indexerLiquidatedSubaccountId := v1.SubaccountIdToIndexerSubaccountId(liquidatedSubaccountId)
	indexerOffsettingSubaccountId := v1.SubaccountIdToIndexerSubaccountId(offsettingSubaccountId)
	expectedDeleveragingEventProto := &events.DeleveragingEventV1{
		Liquidated: indexerLiquidatedSubaccountId,
		Offsetting: indexerOffsettingSubaccountId,
		ClobPairId: clobPairId,
		FillAmount: fillAmount.ToUint64(),
		Subticks:   subticks,
		IsBuy:      isBuy,
	}
	require.Equal(t, expectedDeleveragingEventProto, deleveragingEvent)
}
