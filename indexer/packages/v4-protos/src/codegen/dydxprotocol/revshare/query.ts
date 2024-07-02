import { MarketMapperRevenueShareParams, MarketMapperRevenueShareParamsSDKType } from "./params";
import { MarketMapperRevShareDetails, MarketMapperRevShareDetailsSDKType } from "./revshare";
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