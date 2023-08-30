import { DowntimeParams, DowntimeParamsSDKType } from "./params";
import { Duration, DurationSDKType } from "../../google/protobuf/duration";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/** MsgUpdateDowntimeParams is the Msg/UpdateDowntimeParams request type. */

export interface MsgUpdateDowntimeParams {
  authority: string;
  /** Defines the parameters to update. All parameters must be supplied. */

  params?: DowntimeParams;
}
/** MsgUpdateDowntimeParams is the Msg/UpdateDowntimeParams request type. */

export interface MsgUpdateDowntimeParamsSDKType {
  authority: string;
  /** Defines the parameters to update. All parameters must be supplied. */

  params?: DowntimeParamsSDKType;
}
/**
 * MsgUpdateDowntimeParamsResponse is the Msg/UpdateDowntimeParams response
 * type.
 */

export interface MsgUpdateDowntimeParamsResponse {}
/**
 * MsgUpdateDowntimeParamsResponse is the Msg/UpdateDowntimeParams response
 * type.
 */

export interface MsgUpdateDowntimeParamsResponseSDKType {}
/** MsgIsDelayedBlock is the Msg/IsDelayedBlock request type. */

export interface MsgIsDelayedBlock {
  /**
   * The duration that the block is delayed by.
   * This value could possibly be negative in rare cases.
   */
  delayDuration?: Duration;
}
/** MsgIsDelayedBlock is the Msg/IsDelayedBlock request type. */

export interface MsgIsDelayedBlockSDKType {
  /**
   * The duration that the block is delayed by.
   * This value could possibly be negative in rare cases.
   */
  delay_duration?: DurationSDKType;
}
/** MsgIsDelayedBlock is the Msg/IsDelayedBlock response type. */

export interface MsgIsDelayedBlockResponse {}
/** MsgIsDelayedBlock is the Msg/IsDelayedBlock response type. */

export interface MsgIsDelayedBlockResponseSDKType {}

function createBaseMsgUpdateDowntimeParams(): MsgUpdateDowntimeParams {
  return {
    authority: "",
    params: undefined
  };
}

export const MsgUpdateDowntimeParams = {
  encode(message: MsgUpdateDowntimeParams, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.authority !== "") {
      writer.uint32(10).string(message.authority);
    }

    if (message.params !== undefined) {
      DowntimeParams.encode(message.params, writer.uint32(18).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdateDowntimeParams {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdateDowntimeParams();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.authority = reader.string();
          break;

        case 2:
          message.params = DowntimeParams.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<MsgUpdateDowntimeParams>): MsgUpdateDowntimeParams {
    const message = createBaseMsgUpdateDowntimeParams();
    message.authority = object.authority ?? "";
    message.params = object.params !== undefined && object.params !== null ? DowntimeParams.fromPartial(object.params) : undefined;
    return message;
  }

};

function createBaseMsgUpdateDowntimeParamsResponse(): MsgUpdateDowntimeParamsResponse {
  return {};
}

export const MsgUpdateDowntimeParamsResponse = {
  encode(_: MsgUpdateDowntimeParamsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdateDowntimeParamsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdateDowntimeParamsResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(_: DeepPartial<MsgUpdateDowntimeParamsResponse>): MsgUpdateDowntimeParamsResponse {
    const message = createBaseMsgUpdateDowntimeParamsResponse();
    return message;
  }

};

function createBaseMsgIsDelayedBlock(): MsgIsDelayedBlock {
  return {
    delayDuration: undefined
  };
}

export const MsgIsDelayedBlock = {
  encode(message: MsgIsDelayedBlock, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.delayDuration !== undefined) {
      Duration.encode(message.delayDuration, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgIsDelayedBlock {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgIsDelayedBlock();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.delayDuration = Duration.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<MsgIsDelayedBlock>): MsgIsDelayedBlock {
    const message = createBaseMsgIsDelayedBlock();
    message.delayDuration = object.delayDuration !== undefined && object.delayDuration !== null ? Duration.fromPartial(object.delayDuration) : undefined;
    return message;
  }

};

function createBaseMsgIsDelayedBlockResponse(): MsgIsDelayedBlockResponse {
  return {};
}

export const MsgIsDelayedBlockResponse = {
  encode(_: MsgIsDelayedBlockResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgIsDelayedBlockResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgIsDelayedBlockResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(_: DeepPartial<MsgIsDelayedBlockResponse>): MsgIsDelayedBlockResponse {
    const message = createBaseMsgIsDelayedBlockResponse();
    return message;
  }

};