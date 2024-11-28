import { DowntimeParams, DowntimeParamsSDKType, SynchronyParams, SynchronyParamsSDKType } from "./params";
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
/** MsgUpdateSynchronyParams is the Msg/UpdateSynchronyParams request type. */

export interface MsgUpdateSynchronyParams {
  authority: string;
  /** Defines the parameters to update. All parameters must be supplied. */

  params?: SynchronyParams;
}
/** MsgUpdateSynchronyParams is the Msg/UpdateSynchronyParams request type. */

export interface MsgUpdateSynchronyParamsSDKType {
  authority: string;
  /** Defines the parameters to update. All parameters must be supplied. */

  params?: SynchronyParamsSDKType;
}
/**
 * MsgUpdateSynchronyParamsResponse is the Msg/UpdateSynchronyParams response
 * type.
 */

export interface MsgUpdateSynchronyParamsResponse {}
/**
 * MsgUpdateSynchronyParamsResponse is the Msg/UpdateSynchronyParams response
 * type.
 */

export interface MsgUpdateSynchronyParamsResponseSDKType {}

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

function createBaseMsgUpdateSynchronyParams(): MsgUpdateSynchronyParams {
  return {
    authority: "",
    params: undefined
  };
}

export const MsgUpdateSynchronyParams = {
  encode(message: MsgUpdateSynchronyParams, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.authority !== "") {
      writer.uint32(10).string(message.authority);
    }

    if (message.params !== undefined) {
      SynchronyParams.encode(message.params, writer.uint32(18).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdateSynchronyParams {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdateSynchronyParams();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.authority = reader.string();
          break;

        case 2:
          message.params = SynchronyParams.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<MsgUpdateSynchronyParams>): MsgUpdateSynchronyParams {
    const message = createBaseMsgUpdateSynchronyParams();
    message.authority = object.authority ?? "";
    message.params = object.params !== undefined && object.params !== null ? SynchronyParams.fromPartial(object.params) : undefined;
    return message;
  }

};

function createBaseMsgUpdateSynchronyParamsResponse(): MsgUpdateSynchronyParamsResponse {
  return {};
}

export const MsgUpdateSynchronyParamsResponse = {
  encode(_: MsgUpdateSynchronyParamsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdateSynchronyParamsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdateSynchronyParamsResponse();

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

  fromPartial(_: DeepPartial<MsgUpdateSynchronyParamsResponse>): MsgUpdateSynchronyParamsResponse {
    const message = createBaseMsgUpdateSynchronyParamsResponse();
    return message;
  }

};