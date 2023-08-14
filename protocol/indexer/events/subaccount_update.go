package events

import (
	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/protocol/v1"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

// NewSubaccountUpdateEvent creates a SubaccountUpdateEvent representing a subaccount update
// containing its updated perpetual/asset positions.
func NewSubaccountUpdateEvent(
	subaccountId *satypes.SubaccountId,
	updatedPerpetualPositions []*satypes.PerpetualPosition,
	updatedAssetPositions []*satypes.AssetPosition,
	fundingPayments map[uint32]dtypes.SerializableInt,
) *SubaccountUpdateEventV1 {
	indexerSubaccountId := v1.SubaccountIdToIndexerSubaccountId(*subaccountId)
	return &SubaccountUpdateEventV1{
		SubaccountId: &indexerSubaccountId,
		UpdatedPerpetualPositions: v1.PerpetualPositionsToIndexerPerpetualPositions(
			updatedPerpetualPositions,
			fundingPayments,
		),
		UpdatedAssetPositions: v1.AssetPositionsToIndexerAssetPositions(updatedAssetPositions),
	}
}
