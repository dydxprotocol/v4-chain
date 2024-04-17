package constants

import (
	"fmt"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/indexer/msgsender"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/lib"
	clobtypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/clob/types"
	"github.com/cometbft/cometbft/crypto/tmhash"
)

// Used to construct the constants below.
var (
	testMessage = msgsender.Message{
		Key:   []byte("key"),
		Value: []byte("value"),
	}
	testOrderId = Order_Alice_Num0_Id7_Clob0_Sell25_Price15_GTB20.OrderId
)

var (
	TestTxBytes      = []byte{0x1, 0x2, 0x3}
	TestTxHashBytes  = tmhash.Sum(TestTxBytes)
	TestTxHashString = lib.TxHash(fmt.Sprintf("%X", TestTxHashBytes))
	TestTxHashHeader = msgsender.MessageHeader{
		Key:   msgsender.TransactionHashHeaderKey,
		Value: TestTxHashBytes,
	}
	TestTxBytes1        = []byte{0x4, 0x5, 0x6}
	TestTxHashBytes1    = tmhash.Sum(TestTxBytes1)
	TestTxHashString1   = lib.TxHash(fmt.Sprintf("%X", TestTxHashBytes1))
	TestOffchainUpdates = &clobtypes.OffchainUpdates{
		Messages: []clobtypes.OffchainUpdateMessage{
			{Type: clobtypes.PlaceMessageType, OrderId: testOrderId, Message: testMessage},
			{Type: clobtypes.UpdateMessageType, OrderId: testOrderId, Message: testMessage},
			{Type: clobtypes.RemoveMessageType, OrderId: testOrderId, Message: testMessage},
		},
	}
	TestOffchainMessages = TestOffchainUpdates.GetMessages()
)
