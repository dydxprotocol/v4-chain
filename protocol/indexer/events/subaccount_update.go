package events

import (
	satypes "github.com/dydxprotocol/v4/x/subaccounts/types"
)

// NewSubaccountUpdateEvent creates a SubaccountUpdateEvent representing a subaccount update
// containing its updated perpetual/asset positions.
func NewSubaccountUpdateEvent(
	subaccountId *satypes.SubaccountId,
	updatedPerpetualPositions []*satypes.PerpetualPosition,
	updatedAssetPositions []*satypes.AssetPosition,
) *SubaccountUpdateEvent {
	return &SubaccountUpdateEvent{
		SubaccountId:              subaccountId,
		UpdatedPerpetualPositions: updatedPerpetualPositions,
		UpdatedAssetPositions:     updatedAssetPositions,
	}
}
