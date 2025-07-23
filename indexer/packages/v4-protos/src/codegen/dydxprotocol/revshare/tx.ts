import { MarketMapperRevenueShareParams, MarketMapperRevenueShareParamsSDKType } from "./params";
import { MarketMapperRevShareDetails, MarketMapperRevShareDetailsSDKType, UnconditionalRevShareConfig, UnconditionalRevShareConfigSDKType, OrderRouterRevShare, OrderRouterRevShareSDKType } from "./revshare";
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
/** Message to update the unconditional revenue share config. */

export interface MsgUpdateUnconditionalRevShareConfig {
  authority: string;
  /** The config to update. */

  config?: UnconditionalRevShareConfig;
}
/** Message to update the unconditional revenue share config. */

export interface MsgUpdateUnconditionalRevShareConfigSDKType {
  authority: string;
  /** The config to update. */

  config?: UnconditionalRevShareConfigSDKType;
}
/** Response to MsgUpdateUnconditionalRevShareConfig */

export interface MsgUpdateUnconditionalRevShareConfigResponse {}
/** Response to MsgUpdateUnconditionalRevShareConfig */

export interface MsgUpdateUnconditionalRevShareConfigResponseSDKType {}
/** Governance message to create or update the order router revenue share */

export interface MsgSetOrderRouterRevShare {
  authority: string;
  /** The order router rev share to create or update. */

  orderRouterRevShare?: OrderRouterRevShare;
}
/** Governance message to create or update the order router revenue share */

export interface MsgSetOrderRouterRevShareSDKType {
  authority: string;
  /** The order router rev share to create or update. */

  order_router_rev_share?: OrderRouterRevShareSDKType;
}
/** Response to MsgSetOrderRouterRevShare */

export interface MsgSetOrderRouterRevShareResponse {}
/** Response to MsgSetOrderRouterRevShare */

export interface MsgSetOrderRouterRevShareResponseSDKType {}

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

function createBaseMsgUpdateUnconditionalRevShareConfig(): MsgUpdateUnconditionalRevShareConfig {
  return {
    authority: "",
    config: undefined
  };
}

export const MsgUpdateUnconditionalRevShareConfig = {
  encode(message: MsgUpdateUnconditionalRevShareConfig, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.authority !== "") {
      writer.uint32(10).string(message.authority);
    }

    if (message.config !== undefined) {
      UnconditionalRevShareConfig.encode(message.config, writer.uint32(18).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdateUnconditionalRevShareConfig {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdateUnconditionalRevShareConfig();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.authority = reader.string();
          break;

        case 2:
          message.config = UnconditionalRevShareConfig.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<MsgUpdateUnconditionalRevShareConfig>): MsgUpdateUnconditionalRevShareConfig {
    const message = createBaseMsgUpdateUnconditionalRevShareConfig();
    message.authority = object.authority ?? "";
    message.config = object.config !== undefined && object.config !== null ? UnconditionalRevShareConfig.fromPartial(object.config) : undefined;
    return message;
  }

};

function createBaseMsgUpdateUnconditionalRevShareConfigResponse(): MsgUpdateUnconditionalRevShareConfigResponse {
  return {};
}

export const MsgUpdateUnconditionalRevShareConfigResponse = {
  encode(_: MsgUpdateUnconditionalRevShareConfigResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdateUnconditionalRevShareConfigResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdateUnconditionalRevShareConfigResponse();

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

  fromPartial(_: DeepPartial<MsgUpdateUnconditionalRevShareConfigResponse>): MsgUpdateUnconditionalRevShareConfigResponse {
    const message = createBaseMsgUpdateUnconditionalRevShareConfigResponse();
    return message;
  }

};

function createBaseMsgSetOrderRouterRevShare(): MsgSetOrderRouterRevShare {
  return {
    authority: "",
    orderRouterRevShare: undefined
  };
}

export const MsgSetOrderRouterRevShare = {
  encode(message: MsgSetOrderRouterRevShare, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.authority !== "") {
      writer.uint32(10).string(message.authority);
    }

    if (message.orderRouterRevShare !== undefined) {
      OrderRouterRevShare.encode(message.orderRouterRevShare, writer.uint32(18).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgSetOrderRouterRevShare {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgSetOrderRouterRevShare();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.authority = reader.string();
          break;

        case 2:
          message.orderRouterRevShare = OrderRouterRevShare.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<MsgSetOrderRouterRevShare>): MsgSetOrderRouterRevShare {
    const message = createBaseMsgSetOrderRouterRevShare();
    message.authority = object.authority ?? "";
    message.orderRouterRevShare = object.orderRouterRevShare !== undefined && object.orderRouterRevShare !== null ? OrderRouterRevShare.fromPartial(object.orderRouterRevShare) : undefined;
    return message;
  }

};

function createBaseMsgSetOrderRouterRevShareResponse(): MsgSetOrderRouterRevShareResponse {
  return {};
}

export const MsgSetOrderRouterRevShareResponse = {
  encode(_: MsgSetOrderRouterRevShareResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgSetOrderRouterRevShareResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgSetOrderRouterRevShareResponse();

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

  fromPartial(_: DeepPartial<MsgSetOrderRouterRevShareResponse>): MsgSetOrderRouterRevShareResponse {
    const message = createBaseMsgSetOrderRouterRevShareResponse();
    return message;
  }

};