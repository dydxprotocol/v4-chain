package util

import (
	"fmt"

	"github.com/cosmos/gogoproto/proto"
	ocutypes "github.com/dydxprotocol/v4-chain/protocol/indexer/off_chain_updates/types"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
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

// TODO best practice for ensuring all cases are handled
// default error? default panic?
func GetOffChainUpdateV1SubaccountIdNumber(update ocutypes.OffChainUpdateV1) (uint32, error) {
	var orderSubaccountIdNumber uint32
	switch ocu1 := update.UpdateMessage.(type) {
	case *ocutypes.OffChainUpdateV1_OrderPlace:
		orderSubaccountIdNumber = ocu1.OrderPlace.Order.OrderId.SubaccountId.Number
	case *ocutypes.OffChainUpdateV1_OrderRemove:
		orderSubaccountIdNumber = ocu1.OrderRemove.RemovedOrderId.SubaccountId.Number
	case *ocutypes.OffChainUpdateV1_OrderUpdate:
		orderSubaccountIdNumber = ocu1.OrderUpdate.OrderId.SubaccountId.Number
	case *ocutypes.OffChainUpdateV1_OrderReplace:
		orderSubaccountIdNumber = ocu1.OrderReplace.Order.OrderId.SubaccountId.Number
	default:
		return 0, fmt.Errorf("UpdateMessage type not in {OrderPlace, OrderRemove, OrderUpdate, OrderReplace}")
	}
	return orderSubaccountIdNumber, nil
}
