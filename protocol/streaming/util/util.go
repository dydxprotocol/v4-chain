package util

import (
	"github.com/cosmos/gogoproto/proto"
	ocutypes "github.com/dydxprotocol/v4-chain/protocol/indexer/off_chain_updates/types"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

// GetOffchainUpdatesV1 unmarshals messages in offchain updates to OffchainUpdateV1.
func GetOffchainUpdatesV1(offchainUpdates *clobtypes.OffchainUpdates) ([]ocutypes.OffChainUpdateV1, error) {
	v1updates := make([]ocutypes.OffChainUpdateV1, 0)
	for _, message := range offchainUpdates.Messages {
		var update ocutypes.OffChainUpdateV1
		err := proto.Unmarshal(message.Message.Value, &update)
		if err != nil {
			return nil, err
		}
		v1updates = append(v1updates, update)
	}
	return v1updates, nil
}

func AggregateSubaccountUpdates(subaccountUpdates []satypes.StreamSubaccountUpdate) []satypes.StreamSubaccountUpdate {
	subaccounts := make(map[satypes.SubaccountId][]satypes.StreamSubaccountUpdate)

	for _, update := range subaccountUpdates {
		if update.SubaccountId == nil {
			continue
		}
		subaccountId := *update.SubaccountId

		if updates, exists := subaccounts[subaccountId]; exists {
			lastUpdate := updates[len(updates)-1]

			// If it's a snapshot, or the last update was a snapshot, just append
			if update.Snapshot || lastUpdate.Snapshot {
				subaccounts[subaccountId] = append(subaccounts[subaccountId], update)
			} else {
				// Merge with the last non-snapshot update
				lastUpdate.UpdatedPerpetualPositions = mergePerpetualPositions(
					lastUpdate.UpdatedPerpetualPositions, update.UpdatedPerpetualPositions)

				lastUpdate.UpdatedAssetPositions = mergeAssetPositions(
					lastUpdate.UpdatedAssetPositions, update.UpdatedAssetPositions)

				subaccounts[subaccountId][len(subaccounts[subaccountId])-1] = lastUpdate
			}
		} else {
			// First update for this subaccount, just append
			subaccounts[subaccountId] = []satypes.StreamSubaccountUpdate{update}
		}
	}

	// Convert the subaccounts map to a slice
	aggregatedUpdates := make([]satypes.StreamSubaccountUpdate, 0, len(subaccounts))
	for _, updates := range subaccounts {
		aggregatedUpdates = append(aggregatedUpdates, updates...)
	}

	return aggregatedUpdates
}

// Helper function to merge perpetual positions
func mergePerpetualPositions(existing, updates []*satypes.SubaccountPerpetualPosition) []*satypes.SubaccountPerpetualPosition {
	positionMap := make(map[uint32]*satypes.SubaccountPerpetualPosition)

	for _, pos := range existing {
		positionMap[pos.PerpetualId] = pos
	}

	for _, update := range updates {
		if existingPos, exists := positionMap[update.PerpetualId]; exists {
			existingPos.Quantums = update.Quantums
		} else {
			positionMap[update.PerpetualId] = update
		}
	}

	mergedPositions := make([]*satypes.SubaccountPerpetualPosition, 0, len(positionMap))
	for _, pos := range positionMap {
		mergedPositions = append(mergedPositions, pos)
	}

	return mergedPositions
}

// Helper function to merge asset positions
func mergeAssetPositions(existing, updates []*satypes.SubaccountAssetPosition) []*satypes.SubaccountAssetPosition {
	positionMap := make(map[uint32]*satypes.SubaccountAssetPosition)

	for _, pos := range existing {
		positionMap[pos.AssetId] = pos
	}

	for _, update := range updates {
		if existingPos, exists := positionMap[update.AssetId]; exists {
			existingPos.Quantums = update.Quantums
		} else {
			positionMap[update.AssetId] = update
		}
	}

	mergedPositions := make([]*satypes.SubaccountAssetPosition, 0, len(positionMap))
	for _, pos := range positionMap {
		mergedPositions = append(mergedPositions, pos)
	}

	return mergedPositions
}
