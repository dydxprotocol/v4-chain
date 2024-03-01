package events_test

import (
	"testing"

	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"

	"github.com/dydxprotocol/v4-chain/protocol/indexer/events"
	v1 "github.com/dydxprotocol/v4-chain/protocol/indexer/protocol/v1"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/stretchr/testify/require"
)

var (
	liquidatedSubaccountId = constants.Alice_Num0
	offsettingSubaccountId = constants.Bob_Num0
	perpetualId            = uint32(1)
	price     = satypes.BaseQuantums(1000)
	isBuy                  = true
)

func TestNewDeleveragingEvent_Success(t *testing.T) {
	deleveragingEvent := events.NewDeleveragingEvent(
		liquidatedSubaccountId,
		offsettingSubaccountId,
		perpetualId,
		fillAmount,
		price,
		isBuy,
		false,
	)
	indexerLiquidatedSubaccountId := v1.SubaccountIdToIndexerSubaccountId(liquidatedSubaccountId)
	indexerOffsettingSubaccountId := v1.SubaccountIdToIndexerSubaccountId(offsettingSubaccountId)
	expectedDeleveragingEventProto := &events.DeleveragingEventV2{
		Liquidated:         indexerLiquidatedSubaccountId,
		Offsetting:         indexerOffsettingSubaccountId,
		PerpetualId:        perpetualId,
		FillAmount:         fillAmount.ToUint64(),
		Price: price.ToUint64(),
		IsBuy:              isBuy,
	}
	require.Equal(t, expectedDeleveragingEventProto, deleveragingEvent)
}
