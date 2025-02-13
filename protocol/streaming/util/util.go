package util

import (
	"fmt"

	"github.com/cosmos/gogoproto/proto"
	ocutypes "github.com/dydxprotocol/v4-chain/protocol/indexer/off_chain_updates/types"
	v1types "github.com/dydxprotocol/v4-chain/protocol/indexer/protocol/v1/types"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

// GetOffchainUpdatesV1 unmarshals messages in offchain updates to OffchainUpdateV1.
func GetOffchainUpdatesV1(offchainUpdates *clobtypes.OffchainUpdates) []ocutypes.OffChainUpdateV1 {
	v1updates := make([]ocutypes.OffChainUpdateV1, 0)
	for _, message := range offchainUpdates.Messages {
		var update ocutypes.OffChainUpdateV1
		err := proto.Unmarshal(message.Message.Value, &update)
		if err != nil {
			panic(fmt.Sprintf("Failed to get OffchainUpdatesV1: %v", err))
		}
		v1updates = append(v1updates, update)
	}
	return v1updates
}

// Error expected if OffChainUpdateV1.UpdateMessage message type is extended to more order events
func GetOffChainUpdateV1SubaccountId(update ocutypes.OffChainUpdateV1) (satypes.SubaccountId, error) {
	var orderSubaccountId v1types.IndexerSubaccountId
	switch updateMessage := update.UpdateMessage.(type) {
	case *ocutypes.OffChainUpdateV1_OrderPlace:
		orderSubaccountId = updateMessage.OrderPlace.Order.OrderId.SubaccountId
	case *ocutypes.OffChainUpdateV1_OrderRemove:
		orderSubaccountId = updateMessage.OrderRemove.RemovedOrderId.SubaccountId
	case *ocutypes.OffChainUpdateV1_OrderUpdate:
		orderSubaccountId = updateMessage.OrderUpdate.OrderId.SubaccountId
	case *ocutypes.OffChainUpdateV1_OrderReplace:
		orderSubaccountId = updateMessage.OrderReplace.Order.OrderId.SubaccountId
	default:
		return satypes.SubaccountId{}, fmt.Errorf(
			"UpdateMessage type not in {OrderPlace, OrderRemove, OrderUpdate, OrderReplace}: %+v",
			updateMessage,
		)
	}
	return satypes.SubaccountId{Owner: orderSubaccountId.Owner, Number: orderSubaccountId.Number}, nil
}
