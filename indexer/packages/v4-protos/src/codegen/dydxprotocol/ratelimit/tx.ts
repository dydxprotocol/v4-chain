import { LimitParams, LimitParamsSDKType } from "./limit_params";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/** MsgSetLimitParams is the Msg/SetLimitParams request type. */

export interface MsgSetLimitParams {
  authority: string;
  /** Defines the parameters to set. All parameters must be supplied. */

  limitParams?: LimitParams;
}
/** MsgSetLimitParams is the Msg/SetLimitParams request type. */

export interface MsgSetLimitParamsSDKType {
  authority: string;
  /** Defines the parameters to set. All parameters must be supplied. */

  limit_params?: LimitParamsSDKType;
}
/** MsgSetLimitParamsResponse is the Msg/SetLimitParams response type. */

export interface MsgSetLimitParamsResponse {}
/** MsgSetLimitParamsResponse is the Msg/SetLimitParams response type. */

export interface MsgSetLimitParamsResponseSDKType {}
/** MsgDeleteLimitParams is the Msg/SetLimitParams request type. */

export interface MsgDeleteLimitParams {
  authority: string;
  /** The denom for which the `LimitParams` should be deleted. */

  denom: string;
}
/** MsgDeleteLimitParams is the Msg/SetLimitParams request type. */

export interface MsgDeleteLimitParamsSDKType {
  authority: string;
  /** The denom for which the `LimitParams` should be deleted. */

  denom: string;
}
/** MsgDeleteLimitParamsResponse is the Msg/DeleteLimitParams response type. */

export interface MsgDeleteLimitParamsResponse {}
/** MsgDeleteLimitParamsResponse is the Msg/DeleteLimitParams response type. */

export interface MsgDeleteLimitParamsResponseSDKType {}

function createBaseMsgSetLimitParams(): MsgSetLimitParams {
  return {
    authority: "",
    limitParams: undefined
  };
}

export const MsgSetLimitParams = {
  encode(message: MsgSetLimitParams, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.authority !== "") {
      writer.uint32(10).string(message.authority);
    }

    if (message.limitParams !== undefined) {
      LimitParams.encode(message.limitParams, writer.uint32(18).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgSetLimitParams {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgSetLimitParams();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.authority = reader.string();
          break;

        case 2:
          message.limitParams = LimitParams.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<MsgSetLimitParams>): MsgSetLimitParams {
    const message = createBaseMsgSetLimitParams();
    message.authority = object.authority ?? "";
    message.limitParams = object.limitParams !== undefined && object.limitParams !== null ? LimitParams.fromPartial(object.limitParams) : undefined;
    return message;
  }

};

function createBaseMsgSetLimitParamsResponse(): MsgSetLimitParamsResponse {
  return {};
}

export const MsgSetLimitParamsResponse = {
  encode(_: MsgSetLimitParamsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgSetLimitParamsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgSetLimitParamsResponse();

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

  fromPartial(_: DeepPartial<MsgSetLimitParamsResponse>): MsgSetLimitParamsResponse {
    const message = createBaseMsgSetLimitParamsResponse();
    return message;
  }

};

function createBaseMsgDeleteLimitParams(): MsgDeleteLimitParams {
  return {
    authority: "",
    denom: ""
  };
}

export const MsgDeleteLimitParams = {
  encode(message: MsgDeleteLimitParams, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.authority !== "") {
      writer.uint32(10).string(message.authority);
    }

    if (message.denom !== "") {
      writer.uint32(18).string(message.denom);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgDeleteLimitParams {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgDeleteLimitParams();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.authority = reader.string();
          break;

        case 2:
          message.denom = reader.string();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<MsgDeleteLimitParams>): MsgDeleteLimitParams {
    const message = createBaseMsgDeleteLimitParams();
    message.authority = object.authority ?? "";
    message.denom = object.denom ?? "";
    return message;
  }

};

function createBaseMsgDeleteLimitParamsResponse(): MsgDeleteLimitParamsResponse {
  return {};
}

export const MsgDeleteLimitParamsResponse = {
  encode(_: MsgDeleteLimitParamsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgDeleteLimitParamsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgDeleteLimitParamsResponse();

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

  fromPartial(_: DeepPartial<MsgDeleteLimitParamsResponse>): MsgDeleteLimitParamsResponse {
    const message = createBaseMsgDeleteLimitParamsResponse();
    return message;
  }

};