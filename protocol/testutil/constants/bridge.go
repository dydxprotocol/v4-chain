package constants

import (
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
	coin = sdk.Coin{
		Denom:  "dv4tnt",
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

	// Eth Chain ID.
	EthChainId = 11155111
	// Eth Log of Bridge event ID 0 at block height 3872013 that bridges 42 tokens to address
	// `dydx1qy352euf40x77qfrg4ncn27dauqjx3t83x4ummcpydzk0zdtehhse25p74`.
	EthLog_Event0 = ethcoretypes.Log{
		Topics: []common.Hash{
			common.HexToHash("0xf8dd2841b36f876d311a264058cb076d68674181851a0688c405d2ae917a4fd2"),
			common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"),
		},
		Data: []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 42, 1, 35, 69, 103, 137, 171, 205, 239, 1, 35, 69, 103, 137,
			171, 205, 239, 1, 35, 69, 103, 137, 171, 205, 239, 1, 35, 69, 103, 137, 171, 205, 239,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 96, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 12, 117, 110, 107, 110, 111, 119, 110, 32, 109, 101, 109, 111, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		BlockNumber: 3872013,
	}
	// Eth Log of Bridge event ID 1 at block height 3969937 that bridges 222 tokens to address
	// `dydx1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqsnpjqx`.
	EthLog_Event1 = ethcoretypes.Log{
		Topics: []common.Hash{
			common.HexToHash("0xf8dd2841b36f876d311a264058cb076d68674181851a0688c405d2ae917a4fd2"),
			common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000001"),
		},
		Data: []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 222, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 96, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 128, 1, 129, 2, 130, 3, 131, 4, 132, 5,
			133, 6, 134, 7, 135, 8, 136, 9, 137, 10, 138, 11, 139, 12, 140, 13, 141, 14, 142, 15, 143,
			16, 144, 17, 145, 18, 146, 19, 147, 20, 148, 21, 149, 22, 150, 23, 151, 24, 152, 25, 153,
			26, 154, 27, 155, 28, 156, 29, 157, 30, 158, 31, 159, 32, 160, 33, 161, 34, 162, 35, 163,
			36, 164, 37, 165, 38, 166, 39, 167, 40, 168, 41, 169, 42, 170, 43, 171, 44, 172, 45, 173,
			46, 174, 47, 175, 48, 176, 49, 177, 50, 178, 51, 179, 52, 180, 53, 181, 54, 182, 55, 183,
			56, 184, 57, 185, 58, 186, 59, 187, 60, 188, 61, 189, 62, 190, 63, 191, 64, 192, 65, 193,
			66, 194, 67, 195, 68, 196, 69, 197, 70, 198, 71, 199, 72, 200, 73, 201, 74, 202, 75, 203,
			76, 204, 77, 205, 78, 206, 79, 207, 80, 208, 81, 209, 82, 210, 83, 211, 84, 212, 85, 213,
			86, 214, 87, 215, 88, 216, 89, 217, 90, 218, 91, 219, 92, 220, 93, 221, 94, 222, 95, 223,
			96, 224, 97, 225, 98, 226, 99, 227, 100, 228, 101, 229, 102, 230, 103, 231, 104, 232, 105,
			233, 106, 234, 107, 235, 108, 236, 109, 237, 110, 238, 111, 239, 112, 240, 113, 241, 114,
			242, 115, 243, 116, 244, 117, 245, 118, 246, 119, 247, 120, 248, 121, 249, 122, 250, 123,
			251, 124, 252, 125, 253, 126, 254, 127, 255},
		BlockNumber: 3969937,
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
