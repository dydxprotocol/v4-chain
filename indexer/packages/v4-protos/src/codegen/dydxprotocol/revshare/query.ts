import { MarketMapperRevenueShareParams, MarketMapperRevenueShareParamsSDKType } from "./params";
import { MarketMapperRevShareDetails, MarketMapperRevShareDetailsSDKType, UnconditionalRevShareConfig, UnconditionalRevShareConfigSDKType, OrderRouterRevShare, OrderRouterRevShareSDKType } from "./revshare";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/** Queries for the default market mapper revenue share params */

export interface QueryMarketMapperRevenueShareParams {}
/** Queries for the default market mapper revenue share params */

export interface QueryMarketMapperRevenueShareParamsSDKType {}
/** Response type for QueryMarketMapperRevenueShareParams */

export interface QueryMarketMapperRevenueShareParamsResponse {
  params?: MarketMapperRevenueShareParams;
}
/** Response type for QueryMarketMapperRevenueShareParams */

export interface QueryMarketMapperRevenueShareParamsResponseSDKType {
  params?: MarketMapperRevenueShareParamsSDKType;
}
/** Queries market mapper revenue share details for a specific market */

export interface QueryMarketMapperRevShareDetails {
  /** Queries market mapper revenue share details for a specific market */
  marketId: number;
}
/** Queries market mapper revenue share details for a specific market */

export interface QueryMarketMapperRevShareDetailsSDKType {
  /** Queries market mapper revenue share details for a specific market */
  market_id: number;
}
/** Response type for QueryMarketMapperRevShareDetails */

export interface QueryMarketMapperRevShareDetailsResponse {
  details?: MarketMapperRevShareDetails;
}
/** Response type for QueryMarketMapperRevShareDetails */

export interface QueryMarketMapperRevShareDetailsResponseSDKType {
  details?: MarketMapperRevShareDetailsSDKType;
}
/** Queries unconditional revenue share details */

export interface QueryUnconditionalRevShareConfig {}
/** Queries unconditional revenue share details */

export interface QueryUnconditionalRevShareConfigSDKType {}
/** Response type for QueryUnconditionalRevShareConfig */

export interface QueryUnconditionalRevShareConfigResponse {
  config?: UnconditionalRevShareConfig;
}
/** Response type for QueryUnconditionalRevShareConfig */

export interface QueryUnconditionalRevShareConfigResponseSDKType {
  config?: UnconditionalRevShareConfigSDKType;
}
/** Queries order router rev shares */

export interface QueryOrderRouterRevShare {
  address: string;
}
/** Queries order router rev shares */

export interface QueryOrderRouterRevShareSDKType {
  address: string;
}
/** Response type for QueryOrderRouterRevShare */

export interface QueryOrderRouterRevShareResponse {
  orderRouterRevShare?: OrderRouterRevShare;
}
/** Response type for QueryOrderRouterRevShare */

export interface QueryOrderRouterRevShareResponseSDKType {
  order_router_rev_share?: OrderRouterRevShareSDKType;
}

function createBaseQueryMarketMapperRevenueShareParams(): QueryMarketMapperRevenueShareParams {
  return {};
}

export const QueryMarketMapperRevenueShareParams = {
  encode(_: QueryMarketMapperRevenueShareParams, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryMarketMapperRevenueShareParams {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryMarketMapperRevenueShareParams();

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

  fromPartial(_: DeepPartial<QueryMarketMapperRevenueShareParams>): QueryMarketMapperRevenueShareParams {
    const message = createBaseQueryMarketMapperRevenueShareParams();
    return message;
  }

};

function createBaseQueryMarketMapperRevenueShareParamsResponse(): QueryMarketMapperRevenueShareParamsResponse {
  return {
    params: undefined
  };
}

export const QueryMarketMapperRevenueShareParamsResponse = {
  encode(message: QueryMarketMapperRevenueShareParamsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.params !== undefined) {
      MarketMapperRevenueShareParams.encode(message.params, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryMarketMapperRevenueShareParamsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryMarketMapperRevenueShareParamsResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.params = MarketMapperRevenueShareParams.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QueryMarketMapperRevenueShareParamsResponse>): QueryMarketMapperRevenueShareParamsResponse {
    const message = createBaseQueryMarketMapperRevenueShareParamsResponse();
    message.params = object.params !== undefined && object.params !== null ? MarketMapperRevenueShareParams.fromPartial(object.params) : undefined;
    return message;
  }

};

function createBaseQueryMarketMapperRevShareDetails(): QueryMarketMapperRevShareDetails {
  return {
    marketId: 0
  };
}

export const QueryMarketMapperRevShareDetails = {
  encode(message: QueryMarketMapperRevShareDetails, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.marketId !== 0) {
      writer.uint32(8).uint32(message.marketId);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryMarketMapperRevShareDetails {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryMarketMapperRevShareDetails();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.marketId = reader.uint32();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QueryMarketMapperRevShareDetails>): QueryMarketMapperRevShareDetails {
    const message = createBaseQueryMarketMapperRevShareDetails();
    message.marketId = object.marketId ?? 0;
    return message;
  }

};

function createBaseQueryMarketMapperRevShareDetailsResponse(): QueryMarketMapperRevShareDetailsResponse {
  return {
    details: undefined
  };
}

export const QueryMarketMapperRevShareDetailsResponse = {
  encode(message: QueryMarketMapperRevShareDetailsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      MarketMapperRevShareDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryMarketMapperRevShareDetailsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryMarketMapperRevShareDetailsResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.details = MarketMapperRevShareDetails.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QueryMarketMapperRevShareDetailsResponse>): QueryMarketMapperRevShareDetailsResponse {
    const message = createBaseQueryMarketMapperRevShareDetailsResponse();
    message.details = object.details !== undefined && object.details !== null ? MarketMapperRevShareDetails.fromPartial(object.details) : undefined;
    return message;
  }

};

function createBaseQueryUnconditionalRevShareConfig(): QueryUnconditionalRevShareConfig {
  return {};
}

export const QueryUnconditionalRevShareConfig = {
  encode(_: QueryUnconditionalRevShareConfig, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryUnconditionalRevShareConfig {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryUnconditionalRevShareConfig();

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

  fromPartial(_: DeepPartial<QueryUnconditionalRevShareConfig>): QueryUnconditionalRevShareConfig {
    const message = createBaseQueryUnconditionalRevShareConfig();
    return message;
  }

};

function createBaseQueryUnconditionalRevShareConfigResponse(): QueryUnconditionalRevShareConfigResponse {
  return {
    config: undefined
  };
}

export const QueryUnconditionalRevShareConfigResponse = {
  encode(message: QueryUnconditionalRevShareConfigResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.config !== undefined) {
      UnconditionalRevShareConfig.encode(message.config, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryUnconditionalRevShareConfigResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryUnconditionalRevShareConfigResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.config = UnconditionalRevShareConfig.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QueryUnconditionalRevShareConfigResponse>): QueryUnconditionalRevShareConfigResponse {
    const message = createBaseQueryUnconditionalRevShareConfigResponse();
    message.config = object.config !== undefined && object.config !== null ? UnconditionalRevShareConfig.fromPartial(object.config) : undefined;
    return message;
  }

};

function createBaseQueryOrderRouterRevShare(): QueryOrderRouterRevShare {
  return {
    address: ""
  };
}

export const QueryOrderRouterRevShare = {
  encode(message: QueryOrderRouterRevShare, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.address !== "") {
      writer.uint32(10).string(message.address);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryOrderRouterRevShare {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryOrderRouterRevShare();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.address = reader.string();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QueryOrderRouterRevShare>): QueryOrderRouterRevShare {
    const message = createBaseQueryOrderRouterRevShare();
    message.address = object.address ?? "";
    return message;
  }

};

function createBaseQueryOrderRouterRevShareResponse(): QueryOrderRouterRevShareResponse {
  return {
    orderRouterRevShare: undefined
  };
}

export const QueryOrderRouterRevShareResponse = {
  encode(message: QueryOrderRouterRevShareResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.orderRouterRevShare !== undefined) {
      OrderRouterRevShare.encode(message.orderRouterRevShare, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryOrderRouterRevShareResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryOrderRouterRevShareResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.orderRouterRevShare = OrderRouterRevShare.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QueryOrderRouterRevShareResponse>): QueryOrderRouterRevShareResponse {
    const message = createBaseQueryOrderRouterRevShareResponse();
    message.orderRouterRevShare = object.orderRouterRevShare !== undefined && object.orderRouterRevShare !== null ? OrderRouterRevShare.fromPartial(object.orderRouterRevShare) : undefined;
    return message;
  }

};