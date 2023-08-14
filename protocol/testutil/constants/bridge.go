package constants

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
)

func init() {
	_ = TestTxBuilder.SetMsgs(MsgAcknowledgeBridges_NoEvents)
	MsgAcknowledgeBridges_NoEvents_TxBytes, _ = TestEncodingCfg.TxConfig.TxEncoder()(TestTxBuilder.GetTx())

	_ = TestTxBuilder.SetMsgs(MsgAcknowledgeBridges_Id0_Height0)
	MsgAcknowledgeBridges_Id0_Height0_TxBytes, _ = TestEncodingCfg.TxConfig.TxEncoder()(TestTxBuilder.GetTx())

	_ = TestTxBuilder.SetMsgs(MsgAcknowledgeBridges_Id1_Height0)
	MsgAcknowledgeBridges_Id1_Height0_TxBytes, _ = TestEncodingCfg.TxConfig.TxEncoder()(TestTxBuilder.GetTx())

	_ = TestTxBuilder.SetMsgs(MsgAcknowledgeBridges_Id55_Height15)
	MsgAcknowledgeBridges_Id55_Height15_TxBytes, _ = TestEncodingCfg.TxConfig.TxEncoder()(TestTxBuilder.GetTx())

	_ = TestTxBuilder.SetMsgs(MsgAcknowledgeBridges_Ids0_1_Height0)
	MsgAcknowledgeBridges_Ids0_1_Height0_TxBytes, _ = TestEncodingCfg.TxConfig.TxEncoder()(TestTxBuilder.GetTx())

	_ = TestTxBuilder.SetMsgs(MsgAcknowledgeBridges_Ids0_55_Height0)
	MsgAcknowledgeBridges_Ids0_55_Height0_TxBytes, _ = TestEncodingCfg.TxConfig.TxEncoder()(TestTxBuilder.GetTx())
}

var (
	// Private
	coin = sdk.Coin{
		Denom:  "test-token",
		Amount: sdk.NewIntFromUint64(888),
	}

	// Public
	// Bridge Event.
	BridgeEvent_Id0_Height0 = types.BridgeEvent{
		Id:             0,
		Address:        AliceAccAddress.String(),
		Coin:           coin,
		EthBlockHeight: 0,
	}
	BridgeEvent_Id1_Height0 = types.BridgeEvent{
		Id:             1,
		Address:        BobAccAddress.String(),
		Coin:           coin,
		EthBlockHeight: 0,
	}
	BridgeEvent_Id2_Height1 = types.BridgeEvent{
		Id:             2,
		Address:        BobAccAddress.String(),
		Coin:           coin,
		EthBlockHeight: 1,
	}
	BridgeEvent_Id3_Height3 = types.BridgeEvent{
		Id:             3,
		Address:        CarlAccAddress.String(),
		Coin:           coin,
		EthBlockHeight: 3,
	}
	BridgeEvent_Id55_Height15 = types.BridgeEvent{
		Id:             55,
		Address:        DaveAccAddress.String(),
		Coin:           coin,
		EthBlockHeight: 15,
	}

	// Acknowledge Bridges Tx.
	MsgAcknowledgeBridges_NoEvents = &types.MsgAcknowledgeBridges{
		Events: []types.BridgeEvent{},
	}
	MsgAcknowledgeBridges_NoEvents_TxBytes []byte

	MsgAcknowledgeBridges_Id0_Height0 = &types.MsgAcknowledgeBridges{
		Events: []types.BridgeEvent{
			BridgeEvent_Id0_Height0,
		},
	}
	MsgAcknowledgeBridges_Id0_Height0_TxBytes []byte

	MsgAcknowledgeBridges_Id1_Height0 = &types.MsgAcknowledgeBridges{
		Events: []types.BridgeEvent{
			BridgeEvent_Id1_Height0,
		},
	}
	MsgAcknowledgeBridges_Id1_Height0_TxBytes []byte

	MsgAcknowledgeBridges_Id55_Height15 = &types.MsgAcknowledgeBridges{
		Events: []types.BridgeEvent{
			BridgeEvent_Id55_Height15,
		},
	}
	MsgAcknowledgeBridges_Id55_Height15_TxBytes []byte

	MsgAcknowledgeBridges_Ids0_1_Height0 = &types.MsgAcknowledgeBridges{
		Events: []types.BridgeEvent{
			BridgeEvent_Id0_Height0,
			BridgeEvent_Id1_Height0,
		},
	}
	MsgAcknowledgeBridges_Ids0_1_Height0_TxBytes []byte

	MsgAcknowledgeBridges_Ids0_55_Height0 = &types.MsgAcknowledgeBridges{
		Events: []types.BridgeEvent{
			BridgeEvent_Id0_Height0,
			BridgeEvent_Id55_Height15,
		},
	}
	MsgAcknowledgeBridges_Ids0_55_Height0_TxBytes []byte

	// Event Info.
	AcknowledgedEventInfo_Id0_Height0 = types.BridgeEventInfo{
		NextId:         0,
		EthBlockHeight: 0,
	}
	RecognizedEventInfo_Id2_Height0 = types.BridgeEventInfo{
		NextId:         2,
		EthBlockHeight: 0,
	}
)
