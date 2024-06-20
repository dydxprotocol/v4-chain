import { MarketMapperRevenueShareParams, MarketMapperRevenueShareParamsSDKType } from "./params";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/** Message to set the market mapper revenue share */

export interface MsgSetMarketMapperRevenueShare {
  authority: string;
  /** Parameters for the revenue share */

  params?: MarketMapperRevenueShareParams;
}
/** Message to set the market mapper revenue share */

export interface MsgSetMarketMapperRevenueShareSDKType {
  authority: string;
  /** Parameters for the revenue share */

  params?: MarketMapperRevenueShareParamsSDKType;
}
/** Response to a MsgSetMarketMapperRevenueShare */

export interface MsgSetMarketMapperRevenueShareResponse {}
/** Response to a MsgSetMarketMapperRevenueShare */

export interface MsgSetMarketMapperRevenueShareResponseSDKType {}

function createBaseMsgSetMarketMapperRevenueShare(): MsgSetMarketMapperRevenueShare {
  return {
    authority: "",
    params: undefined
  };
}

export const MsgSetMarketMapperRevenueShare = {
  encode(message: MsgSetMarketMapperRevenueShare, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.authority !== "") {
      writer.uint32(10).string(message.authority);
    }

    if (message.params !== undefined) {
      MarketMapperRevenueShareParams.encode(message.params, writer.uint32(18).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgSetMarketMapperRevenueShare {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgSetMarketMapperRevenueShare();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.authority = reader.string();
          break;

        case 2:
          message.params = MarketMapperRevenueShareParams.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<MsgSetMarketMapperRevenueShare>): MsgSetMarketMapperRevenueShare {
    const message = createBaseMsgSetMarketMapperRevenueShare();
    message.authority = object.authority ?? "";
    message.params = object.params !== undefined && object.params !== null ? MarketMapperRevenueShareParams.fromPartial(object.params) : undefined;
    return message;
  }

};

function createBaseMsgSetMarketMapperRevenueShareResponse(): MsgSetMarketMapperRevenueShareResponse {
  return {};
}

export const MsgSetMarketMapperRevenueShareResponse = {
  encode(_: MsgSetMarketMapperRevenueShareResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgSetMarketMapperRevenueShareResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgSetMarketMapperRevenueShareResponse();

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

  fromPartial(_: DeepPartial<MsgSetMarketMapperRevenueShareResponse>): MsgSetMarketMapperRevenueShareResponse {
    const message = createBaseMsgSetMarketMapperRevenueShareResponse();
    return message;
  }

};