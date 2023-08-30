package events_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/events"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/protocol/v1"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/require"
)

var (
	subaccountId              = constants.Alice_Num0
	indexerSubaccountId       = v1.SubaccountIdToIndexerSubaccountId(subaccountId)
	updatedPerpetualPositions = []*satypes.PerpetualPosition{
		&constants.Long_Perp_1BTC_PositiveFunding,
		&constants.Short_Perp_1ETH_NegativeFunding,
	}
	fundingPayments = map[uint32]dtypes.SerializableInt{
		constants.Long_Perp_1BTC_PositiveFunding.PerpetualId: dtypes.NewInt(500),
	}
	indexerPerpetualPositions = v1.PerpetualPositionsToIndexerPerpetualPositions(
		updatedPerpetualPositions,
		fundingPayments,
	)
	updatedAssetPositions = []*satypes.AssetPosition{
		&constants.Short_Asset_1BTC,
		&constants.Long_Asset_1ETH,
	}
	indexerAssetPositions = v1.AssetPositionsToIndexerAssetPositions(updatedAssetPositions)
)

func TestNewSubaccountUpdateEvent_Success(t *testing.T) {
	subaccountUpdateEvent := events.NewSubaccountUpdateEvent(
		&subaccountId,
		updatedPerpetualPositions,
		updatedAssetPositions,
		fundingPayments,
	)
	expectedSubaccountUpdateEventProto := &events.SubaccountUpdateEventV1{
		SubaccountId:              &indexerSubaccountId,
		UpdatedPerpetualPositions: indexerPerpetualPositions,
		UpdatedAssetPositions:     indexerAssetPositions,
	}
	require.Equal(t, expectedSubaccountUpdateEventProto, subaccountUpdateEvent)
}
