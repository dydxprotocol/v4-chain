import { Coin, CoinAmino, CoinSDKType } from "../../cosmos/base/v1beta1/coin";
import { BinaryReader, BinaryWriter } from "../../binary";
/** BridgeEvent is a recognized event from the Ethereum blockchain. */
export interface BridgeEvent {
  /** The unique id of the Ethereum event log. */
  id: number;
  /** The tokens bridged. */
  coin: Coin;
  /** The account address or module address to bridge to. */
  address: string;
  /** The Ethereum block height of the event. */
  ethBlockHeight: bigint;
}
export interface BridgeEventProtoMsg {
  typeUrl: "/dydxprotocol.bridge.BridgeEvent";
  value: Uint8Array;
}
/** BridgeEvent is a recognized event from the Ethereum blockchain. */
export interface BridgeEventAmino {
  /** The unique id of the Ethereum event log. */
  id?: number;
  /** The tokens bridged. */
  coin?: CoinAmino;
  /** The account address or module address to bridge to. */
  address?: string;
  /** The Ethereum block height of the event. */
  eth_block_height?: string;
}
export interface BridgeEventAminoMsg {
  type: "/dydxprotocol.bridge.BridgeEvent";
  value: BridgeEventAmino;
}
/** BridgeEvent is a recognized event from the Ethereum blockchain. */
export interface BridgeEventSDKType {
  id: number;
  coin: CoinSDKType;
  address: string;
  eth_block_height: bigint;
}
function createBaseBridgeEvent(): BridgeEvent {
  return {
    id: 0,
    coin: Coin.fromPartial({}),
    address: "",
    ethBlockHeight: BigInt(0)
  };
}
export const BridgeEvent = {
  typeUrl: "/dydxprotocol.bridge.BridgeEvent",
  encode(message: BridgeEvent, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.id !== 0) {
      writer.uint32(8).uint32(message.id);
    }
    if (message.coin !== undefined) {
      Coin.encode(message.coin, writer.uint32(18).fork()).ldelim();
    }
    if (message.address !== "") {
      writer.uint32(26).string(message.address);
    }
    if (message.ethBlockHeight !== BigInt(0)) {
      writer.uint32(32).uint64(message.ethBlockHeight);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): BridgeEvent {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseBridgeEvent();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.id = reader.uint32();
          break;
        case 2:
          message.coin = Coin.decode(reader, reader.uint32());
          break;
        case 3:
          message.address = reader.string();
          break;
        case 4:
          message.ethBlockHeight = reader.uint64();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<BridgeEvent>): BridgeEvent {
    const message = createBaseBridgeEvent();
    message.id = object.id ?? 0;
    message.coin = object.coin !== undefined && object.coin !== null ? Coin.fromPartial(object.coin) : undefined;
    message.address = object.address ?? "";
    message.ethBlockHeight = object.ethBlockHeight !== undefined && object.ethBlockHeight !== null ? BigInt(object.ethBlockHeight.toString()) : BigInt(0);
    return message;
  },
  fromAmino(object: BridgeEventAmino): BridgeEvent {
    const message = createBaseBridgeEvent();
    if (object.id !== undefined && object.id !== null) {
      message.id = object.id;
    }
    if (object.coin !== undefined && object.coin !== null) {
      message.coin = Coin.fromAmino(object.coin);
    }
    if (object.address !== undefined && object.address !== null) {
      message.address = object.address;
    }
    if (object.eth_block_height !== undefined && object.eth_block_height !== null) {
      message.ethBlockHeight = BigInt(object.eth_block_height);
    }
    return message;
  },
  toAmino(message: BridgeEvent): BridgeEventAmino {
    const obj: any = {};
    obj.id = message.id;
    obj.coin = message.coin ? Coin.toAmino(message.coin) : undefined;
    obj.address = message.address;
    obj.eth_block_height = message.ethBlockHeight ? message.ethBlockHeight.toString() : undefined;
    return obj;
  },
  fromAminoMsg(object: BridgeEventAminoMsg): BridgeEvent {
    return BridgeEvent.fromAmino(object.value);
  },
  fromProtoMsg(message: BridgeEventProtoMsg): BridgeEvent {
    return BridgeEvent.decode(message.value);
  },
  toProto(message: BridgeEvent): Uint8Array {
    return BridgeEvent.encode(message).finish();
  },
  toProtoMsg(message: BridgeEvent): BridgeEventProtoMsg {
    return {
      typeUrl: "/dydxprotocol.bridge.BridgeEvent",
      value: BridgeEvent.encode(message).finish()
    };
  }
};