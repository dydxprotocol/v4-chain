import { Any, AnySDKType } from "../../google/protobuf/any";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial, Long } from "../../helpers";
/** MsgDelayMessage is a request type for the DelayMessage method. */

export interface MsgDelayMessage {
  authority: string;
  /** The message to be delayed. */

  msg?: Any;
  /** The number of blocks to delay the message for. */

  delayBlocks: number;
}
/** MsgDelayMessage is a request type for the DelayMessage method. */

export interface MsgDelayMessageSDKType {
  authority: string;
  /** The message to be delayed. */

  msg?: AnySDKType;
  /** The number of blocks to delay the message for. */

  delay_blocks: number;
}
/** MsgDelayMessageResponse is a response type for the DelayMessage method. */

export interface MsgDelayMessageResponse {
  /** The id of the created delayed message. */
  id: Long;
}
/** MsgDelayMessageResponse is a response type for the DelayMessage method. */

export interface MsgDelayMessageResponseSDKType {
  /** The id of the created delayed message. */
  id: Long;
}

function createBaseMsgDelayMessage(): MsgDelayMessage {
  return {
    authority: "",
    msg: undefined,
    delayBlocks: 0
  };
}

export const MsgDelayMessage = {
  encode(message: MsgDelayMessage, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.authority !== "") {
      writer.uint32(10).string(message.authority);
    }

    if (message.msg !== undefined) {
      Any.encode(message.msg, writer.uint32(18).fork()).ldelim();
    }

    if (message.delayBlocks !== 0) {
      writer.uint32(24).uint32(message.delayBlocks);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgDelayMessage {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgDelayMessage();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.authority = reader.string();
          break;

        case 2:
          message.msg = Any.decode(reader, reader.uint32());
          break;

        case 3:
          message.delayBlocks = reader.uint32();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<MsgDelayMessage>): MsgDelayMessage {
    const message = createBaseMsgDelayMessage();
    message.authority = object.authority ?? "";
    message.msg = object.msg !== undefined && object.msg !== null ? Any.fromPartial(object.msg) : undefined;
    message.delayBlocks = object.delayBlocks ?? 0;
    return message;
  }

};

function createBaseMsgDelayMessageResponse(): MsgDelayMessageResponse {
  return {
    id: Long.UZERO
  };
}

export const MsgDelayMessageResponse = {
  encode(message: MsgDelayMessageResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (!message.id.isZero()) {
      writer.uint32(8).uint64(message.id);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgDelayMessageResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgDelayMessageResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.id = (reader.uint64() as Long);
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<MsgDelayMessageResponse>): MsgDelayMessageResponse {
    const message = createBaseMsgDelayMessageResponse();
    message.id = object.id !== undefined && object.id !== null ? Long.fromValue(object.id) : Long.UZERO;
    return message;
  }

};