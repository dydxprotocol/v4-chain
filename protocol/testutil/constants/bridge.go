package constants

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
	"github.com/ethereum/go-ethereum/common"
	ethcoretypes "github.com/ethereum/go-ethereum/core/types"
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
	emptyCoin = sdk.Coin{
		Denom:  "adv4tnt",
		Amount: sdkmath.NewInt(0),
	}
	coin = sdk.Coin{
		Denom:  "adv4tnt",
		Amount: sdkmath.NewIntFromUint64(888),
	}

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
	BridgeEvent_Id4_Height0_EmptyCoin = types.BridgeEvent{
		Id:             0,
		Address:        AliceAccAddress.String(),
		Coin:           emptyCoin,
		EthBlockHeight: 0,
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

	// Eth Chain ID.
	EthChainId = 11155111
	// Eth Log of Bridge event ID 0 at block height 3872013 that bridges 12345 tokens to address
	// `dydx1qqgzqvzq2ps8pqys5zcvp58q7rluextx92xhln`.
	EthLog_Event0 = ethcoretypes.Log{
		Topics: []common.Hash{
			common.HexToHash("0x498a04382650bc110983392ed12ab27595af8ece270a344fc70d773d2481043a"),
			common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"),
		},
		Data: []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 48, 57, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 78, 213, 157, 76, 225, 209, 26,
			78, 111, 27, 28, 56, 208, 105, 160, 45, 130, 224, 44, 163, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 192, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 20, 0,
			16, 32, 48, 64, 80, 96, 112, 128, 144, 160, 176, 192, 208, 224, 240, 255, 204, 153, 102,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		BlockNumber: 3872013,
	}
	// Eth Log of Bridge event ID 1 at block height 3969937 that bridges 55 tokens to an empty address.
	EthLog_Event1 = ethcoretypes.Log{
		Topics: []common.Hash{
			common.HexToHash("0x498a04382650bc110983392ed12ab27595af8ece270a344fc70d773d2481043a"),
			common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000001"),
		},
		Data: []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 55, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 78, 213, 157, 76, 225, 209, 26,
			78, 111, 27, 28, 56, 208, 105, 160, 45, 130, 224, 44, 163, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 160, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 3,
			1, 2, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0},
		BlockNumber: 3969937,
	}
	// Eth Log of Bridge event ID 2 at block height 4139345 that bridges 777 tokens to address
	// `dydx1qqgzqvzq2ps8pqys5zcvp58q7rluextxzy3rx3z4vemc3xgq42as94fpcv` (32-byte address).
	EthLog_Event2 = ethcoretypes.Log{
		Topics: []common.Hash{
			common.HexToHash("0x498a04382650bc110983392ed12ab27595af8ece270a344fc70d773d2481043a"),
			common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000002"),
		},
		Data: []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 3, 9, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 78, 213, 157, 76, 225, 209, 26, 78, 111,
			27, 28, 56, 208, 105, 160, 45, 130, 224, 44, 163, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 192, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 32, 0, 16, 32, 48, 64, 80, 96,
			112, 128, 144, 160, 176, 192, 208, 224, 240, 255, 204, 153, 102, 17, 34, 51, 68, 85, 102, 119,
			136, 153, 0, 170, 187, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 2, 18, 52, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0},
		BlockNumber: 4139345,
	}
	// Eth Log of Bridge event ID 3 at block height 4139348 that bridges 888 tokens to address
	// `dydx124n92ej4ve2kv4tx24n92ej4ve2kv4tx24n92ej4ve2kv4tx24nyggjyyfzzy3pzgs3yggjyyfzzy3pzgs3ygg
	// jyyfzzy3pzgs3q8699x3` (62-byte address).
	EthLog_Event3 = ethcoretypes.Log{
		Topics: []common.Hash{
			common.HexToHash("0x498a04382650bc110983392ed12ab27595af8ece270a344fc70d773d2481043a"),
			common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000003"),
		},
		Data: []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 3, 120, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 78, 213, 157, 76, 225, 209, 26,
			78, 111, 27, 28, 56, 208, 105, 160, 45, 130, 224, 44, 163, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 224, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 62, 85,
			102, 85, 102, 85, 102, 85, 102, 85, 102, 85, 102, 85, 102, 85, 102, 85, 102, 85, 102, 85,
			102, 85, 102, 85, 102, 85, 102, 85, 102, 85, 102, 68, 34, 68, 34, 68, 34, 68, 34, 68, 34,
			68, 34, 68, 34, 68, 34, 68, 34, 68, 34, 68, 34, 68, 34, 68, 34, 68, 34, 68, 34, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			2, 255, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0},
		BlockNumber: 4139348,
	}
	// Eth Log of Bridge event ID 4 at block height 4139349 that bridges 1234123443214321 tokens to
	// address `dydx1zg6pydqhy4yy9`.
	EthLog_Event4 = ethcoretypes.Log{
		Topics: []common.Hash{
			common.HexToHash("0x498a04382650bc110983392ed12ab27595af8ece270a344fc70d773d2481043a"),
			common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000004"),
		},
		Data: []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4, 98,
			109, 193, 113, 23, 241, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 78, 213, 157, 76, 225, 209, 26,
			78, 111, 27, 28, 56, 208, 105, 160, 45, 130, 224, 44, 163, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 192, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4, 18, 52, 18, 52,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4, 67,
			33, 67, 33, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0},
		BlockNumber: 4139349,
	}

	// Params
	// Event Params.
	EventParams = types.EventParams{
		Denom:      coin.Denom,
		EthChainId: uint64(EthChainId),
		EthAddress: AliceAccAddress.String(),
	}
	// Propose Params.
	ProposeParams = types.ProposeParams{
		MaxBridgesPerBlock:           2,
		ProposeDelayDuration:         1,
		SkipRatePpm:                  0,
		SkipIfBlockDelayedByDuration: 1,
	}
)
