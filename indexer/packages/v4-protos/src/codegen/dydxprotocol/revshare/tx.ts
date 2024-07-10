import { MarketMapperRevenueShareParams, MarketMapperRevenueShareParamsSDKType } from "./params";
import { MarketMapperRevShareDetails, MarketMapperRevShareDetailsSDKType } from "./revshare";
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
/**
 * Msg to set market mapper revenue share details (e.g. expiration timestamp)
 * for a specific market. To be used as an override for existing revenue share
 * settings set by the MsgSetMarketMapperRevenueShare msg
 */

export interface MsgSetMarketMapperRevShareDetailsForMarket {
  authority: string;
  /** The market ID for which to set the revenue share details */

  marketId: number;
  /** Parameters for the revenue share details */

  params?: MarketMapperRevShareDetails;
}
/**
 * Msg to set market mapper revenue share details (e.g. expiration timestamp)
 * for a specific market. To be used as an override for existing revenue share
 * settings set by the MsgSetMarketMapperRevenueShare msg
 */

export interface MsgSetMarketMapperRevShareDetailsForMarketSDKType {
  authority: string;
  /** The market ID for which to set the revenue share details */

  market_id: number;
  /** Parameters for the revenue share details */

  params?: MarketMapperRevShareDetailsSDKType;
}
/** Response to a MsgSetMarketMapperRevShareDetailsForMarket */

export interface MsgSetMarketMapperRevShareDetailsForMarketResponse {}
/** Response to a MsgSetMarketMapperRevShareDetailsForMarket */

export interface MsgSetMarketMapperRevShareDetailsForMarketResponseSDKType {}

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

function createBaseMsgSetMarketMapperRevShareDetailsForMarket(): MsgSetMarketMapperRevShareDetailsForMarket {
  return {
    authority: "",
    marketId: 0,
    params: undefined
  };
}

export const MsgSetMarketMapperRevShareDetailsForMarket = {
  encode(message: MsgSetMarketMapperRevShareDetailsForMarket, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.authority !== "") {
      writer.uint32(10).string(message.authority);
    }

    if (message.marketId !== 0) {
      writer.uint32(16).uint32(message.marketId);
    }

    if (message.params !== undefined) {
      MarketMapperRevShareDetails.encode(message.params, writer.uint32(26).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgSetMarketMapperRevShareDetailsForMarket {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgSetMarketMapperRevShareDetailsForMarket();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.authority = reader.string();
          break;

        case 2:
          message.marketId = reader.uint32();
          break;

        case 3:
          message.params = MarketMapperRevShareDetails.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<MsgSetMarketMapperRevShareDetailsForMarket>): MsgSetMarketMapperRevShareDetailsForMarket {
    const message = createBaseMsgSetMarketMapperRevShareDetailsForMarket();
    message.authority = object.authority ?? "";
    message.marketId = object.marketId ?? 0;
    message.params = object.params !== undefined && object.params !== null ? MarketMapperRevShareDetails.fromPartial(object.params) : undefined;
    return message;
  }

};

function createBaseMsgSetMarketMapperRevShareDetailsForMarketResponse(): MsgSetMarketMapperRevShareDetailsForMarketResponse {
  return {};
}

export const MsgSetMarketMapperRevShareDetailsForMarketResponse = {
  encode(_: MsgSetMarketMapperRevShareDetailsForMarketResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgSetMarketMapperRevShareDetailsForMarketResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgSetMarketMapperRevShareDetailsForMarketResponse();

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

  fromPartial(_: DeepPartial<MsgSetMarketMapperRevShareDetailsForMarketResponse>): MsgSetMarketMapperRevShareDetailsForMarketResponse {
    const message = createBaseMsgSetMarketMapperRevShareDetailsForMarketResponse();
    return message;
  }

};