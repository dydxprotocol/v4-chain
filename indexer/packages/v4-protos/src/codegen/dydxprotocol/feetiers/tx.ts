import { PerpetualFeeParams, PerpetualFeeParamsSDKType } from "./params";
import { FeeDiscountCampaignParams, FeeDiscountCampaignParamsSDKType } from "./fee_discount_campaign";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/** MsgUpdatePerpetualFeeParams is the Msg/UpdatePerpetualFeeParams request type. */

export interface MsgUpdatePerpetualFeeParams {
  authority: string;
  /** Defines the parameters to update. All parameters must be supplied. */

  params?: PerpetualFeeParams;
}
/** MsgUpdatePerpetualFeeParams is the Msg/UpdatePerpetualFeeParams request type. */

export interface MsgUpdatePerpetualFeeParamsSDKType {
  authority: string;
  /** Defines the parameters to update. All parameters must be supplied. */

  params?: PerpetualFeeParamsSDKType;
}
/**
 * MsgUpdatePerpetualFeeParamsResponse is the Msg/UpdatePerpetualFeeParams
 * response type.
 */

export interface MsgUpdatePerpetualFeeParamsResponse {}
/**
 * MsgUpdatePerpetualFeeParamsResponse is the Msg/UpdatePerpetualFeeParams
 * response type.
 */

export interface MsgUpdatePerpetualFeeParamsResponseSDKType {}
/**
 * MsgSetFeeDiscountCampaignParams is the Msg/SetFeeDiscountCampaignParams
 * request type.
 */

export interface MsgSetFeeDiscountCampaignParams {
  /** authority is the address that controls the module */
  authority: string;
  /** The fee discount campaigns to create or update */

  params: FeeDiscountCampaignParams[];
}
/**
 * MsgSetFeeDiscountCampaignParams is the Msg/SetFeeDiscountCampaignParams
 * request type.
 */

export interface MsgSetFeeDiscountCampaignParamsSDKType {
  /** authority is the address that controls the module */
  authority: string;
  /** The fee discount campaigns to create or update */

  params: FeeDiscountCampaignParamsSDKType[];
}
/**
 * MsgSetFeeDiscountCampaignParamsResponse is the
 * Msg/SetFeeDiscountCampaignParams response type.
 */

export interface MsgSetFeeDiscountCampaignParamsResponse {}
/**
 * MsgSetFeeDiscountCampaignParamsResponse is the
 * Msg/SetFeeDiscountCampaignParams response type.
 */

export interface MsgSetFeeDiscountCampaignParamsResponseSDKType {}

function createBaseMsgUpdatePerpetualFeeParams(): MsgUpdatePerpetualFeeParams {
  return {
    authority: "",
    params: undefined
  };
}

export const MsgUpdatePerpetualFeeParams = {
  encode(message: MsgUpdatePerpetualFeeParams, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.authority !== "") {
      writer.uint32(10).string(message.authority);
    }

    if (message.params !== undefined) {
      PerpetualFeeParams.encode(message.params, writer.uint32(18).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdatePerpetualFeeParams {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdatePerpetualFeeParams();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.authority = reader.string();
          break;

        case 2:
          message.params = PerpetualFeeParams.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<MsgUpdatePerpetualFeeParams>): MsgUpdatePerpetualFeeParams {
    const message = createBaseMsgUpdatePerpetualFeeParams();
    message.authority = object.authority ?? "";
    message.params = object.params !== undefined && object.params !== null ? PerpetualFeeParams.fromPartial(object.params) : undefined;
    return message;
  }

};

function createBaseMsgUpdatePerpetualFeeParamsResponse(): MsgUpdatePerpetualFeeParamsResponse {
  return {};
}

export const MsgUpdatePerpetualFeeParamsResponse = {
  encode(_: MsgUpdatePerpetualFeeParamsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdatePerpetualFeeParamsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdatePerpetualFeeParamsResponse();

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

  fromPartial(_: DeepPartial<MsgUpdatePerpetualFeeParamsResponse>): MsgUpdatePerpetualFeeParamsResponse {
    const message = createBaseMsgUpdatePerpetualFeeParamsResponse();
    return message;
  }

};

function createBaseMsgSetFeeDiscountCampaignParams(): MsgSetFeeDiscountCampaignParams {
  return {
    authority: "",
    params: []
  };
}

export const MsgSetFeeDiscountCampaignParams = {
  encode(message: MsgSetFeeDiscountCampaignParams, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.authority !== "") {
      writer.uint32(10).string(message.authority);
    }

    for (const v of message.params) {
      FeeDiscountCampaignParams.encode(v!, writer.uint32(18).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgSetFeeDiscountCampaignParams {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgSetFeeDiscountCampaignParams();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.authority = reader.string();
          break;

        case 2:
          message.params.push(FeeDiscountCampaignParams.decode(reader, reader.uint32()));
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<MsgSetFeeDiscountCampaignParams>): MsgSetFeeDiscountCampaignParams {
    const message = createBaseMsgSetFeeDiscountCampaignParams();
    message.authority = object.authority ?? "";
    message.params = object.params?.map(e => FeeDiscountCampaignParams.fromPartial(e)) || [];
    return message;
  }

};

function createBaseMsgSetFeeDiscountCampaignParamsResponse(): MsgSetFeeDiscountCampaignParamsResponse {
  return {};
}

export const MsgSetFeeDiscountCampaignParamsResponse = {
  encode(_: MsgSetFeeDiscountCampaignParamsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgSetFeeDiscountCampaignParamsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgSetFeeDiscountCampaignParamsResponse();

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

  fromPartial(_: DeepPartial<MsgSetFeeDiscountCampaignParamsResponse>): MsgSetFeeDiscountCampaignParamsResponse {
    const message = createBaseMsgSetFeeDiscountCampaignParamsResponse();
    return message;
  }

};