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