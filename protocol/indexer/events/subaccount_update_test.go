package events_test

import (
	"testing"

	"github.com/dydxprotocol/v4/indexer/events"
	"github.com/dydxprotocol/v4/testutil/constants"
	satypes "github.com/dydxprotocol/v4/x/subaccounts/types"
	"github.com/stretchr/testify/require"
)

var (
	subaccountId              = &constants.Alice_Num0
	updatedPerpetualPositions = []*satypes.PerpetualPosition{
		&constants.Long_Perp_1BTC_PositiveFunding,
		&constants.Short_Perp_1ETH_NegativeFunding,
	}
	updatedAssetPositions = []*satypes.AssetPosition{
		&constants.Short_Asset_1BTC,
		&constants.Long_Asset_1ETH,
	}
)

func TestNewSubaccountUpdateEvent_Success(t *testing.T) {
	subaccountUpdateEvent := events.NewSubaccountUpdateEvent(
		subaccountId,
		updatedPerpetualPositions,
		updatedAssetPositions,
	)
	expectedSubaccountUpdateEventProto := &events.SubaccountUpdateEvent{
		SubaccountId:              subaccountId,
		UpdatedPerpetualPositions: updatedPerpetualPositions,
		UpdatedAssetPositions:     updatedAssetPositions,
	}
	require.Equal(t, expectedSubaccountUpdateEventProto, subaccountUpdateEvent)
}
