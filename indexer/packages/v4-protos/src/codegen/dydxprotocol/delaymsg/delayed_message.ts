import { Any, AnyAmino, AnySDKType } from "../../google/protobuf/any";
import { BinaryReader, BinaryWriter } from "../../binary";
/** DelayedMessage is a message that is delayed until a certain block height. */
export interface DelayedMessage {
  /** The ID of the delayed message. */
  id: number;
  /** The message to be executed. */
  msg?: Any;
  /** The block height at which the message should be executed. */
  blockHeight: number;
}
export interface DelayedMessageProtoMsg {
  typeUrl: "/dydxprotocol.delaymsg.DelayedMessage";
  value: Uint8Array;
}
/** DelayedMessage is a message that is delayed until a certain block height. */
export interface DelayedMessageAmino {
  /** The ID of the delayed message. */
  id?: number;
  /** The message to be executed. */
  msg?: AnyAmino;
  /** The block height at which the message should be executed. */
  block_height?: number;
}
export interface DelayedMessageAminoMsg {
  type: "/dydxprotocol.delaymsg.DelayedMessage";
  value: DelayedMessageAmino;
}
/** DelayedMessage is a message that is delayed until a certain block height. */
export interface DelayedMessageSDKType {
  id: number;
  msg?: AnySDKType;
  block_height: number;
}
function createBaseDelayedMessage(): DelayedMessage {
  return {
    id: 0,
    msg: undefined,
    blockHeight: 0
  };
}
export const DelayedMessage = {
  typeUrl: "/dydxprotocol.delaymsg.DelayedMessage",
  encode(message: DelayedMessage, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.id !== 0) {
      writer.uint32(8).uint32(message.id);
    }
    if (message.msg !== undefined) {
      Any.encode(message.msg, writer.uint32(18).fork()).ldelim();
    }
    if (message.blockHeight !== 0) {
      writer.uint32(24).uint32(message.blockHeight);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): DelayedMessage {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDelayedMessage();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.id = reader.uint32();
          break;
        case 2:
          message.msg = Any.decode(reader, reader.uint32());
          break;
        case 3:
          message.blockHeight = reader.uint32();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<DelayedMessage>): DelayedMessage {
    const message = createBaseDelayedMessage();
    message.id = object.id ?? 0;
    message.msg = object.msg !== undefined && object.msg !== null ? Any.fromPartial(object.msg) : undefined;
    message.blockHeight = object.blockHeight ?? 0;
    return message;
  },
  fromAmino(object: DelayedMessageAmino): DelayedMessage {
    const message = createBaseDelayedMessage();
    if (object.id !== undefined && object.id !== null) {
      message.id = object.id;
    }
    if (object.msg !== undefined && object.msg !== null) {
      message.msg = Any.fromAmino(object.msg);
    }
    if (object.block_height !== undefined && object.block_height !== null) {
      message.blockHeight = object.block_height;
    }
    return message;
  },
  toAmino(message: DelayedMessage): DelayedMessageAmino {
    const obj: any = {};
    obj.id = message.id;
    obj.msg = message.msg ? Any.toAmino(message.msg) : undefined;
    obj.block_height = message.blockHeight;
    return obj;
  },
  fromAminoMsg(object: DelayedMessageAminoMsg): DelayedMessage {
    return DelayedMessage.fromAmino(object.value);
  },
  fromProtoMsg(message: DelayedMessageProtoMsg): DelayedMessage {
    return DelayedMessage.decode(message.value);
  },
  toProto(message: DelayedMessage): Uint8Array {
    return DelayedMessage.encode(message).finish();
  },
  toProtoMsg(message: DelayedMessage): DelayedMessageProtoMsg {
    return {
      typeUrl: "/dydxprotocol.delaymsg.DelayedMessage",
      value: DelayedMessage.encode(message).finish()
    };
  }
};