import { PageRequest, PageRequestSDKType, PageResponse, PageResponseSDKType } from "../../cosmos/base/query/v1beta1/pagination";
import { MarketPrice, MarketPriceSDKType } from "./market_price";
import { MarketParam, MarketParamSDKType } from "./market_param";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/**
 * QueryMarketPriceRequest is request type for the Query/Params `MarketPrice`
 * RPC method.
 */

export interface QueryMarketPriceRequest {
  /**
   * QueryMarketPriceRequest is request type for the Query/Params `MarketPrice`
   * RPC method.
   */
  id: number;
}
/**
 * QueryMarketPriceRequest is request type for the Query/Params `MarketPrice`
 * RPC method.
 */

export interface QueryMarketPriceRequestSDKType {
  /**
   * QueryMarketPriceRequest is request type for the Query/Params `MarketPrice`
   * RPC method.
   */
  id: number;
}
/**
 * QueryMarketPriceResponse is response type for the Query/Params `MarketPrice`
 * RPC method.
 */

export interface QueryMarketPriceResponse {
  marketPrice?: MarketPrice;
}
/**
 * QueryMarketPriceResponse is response type for the Query/Params `MarketPrice`
 * RPC method.
 */

export interface QueryMarketPriceResponseSDKType {
  market_price?: MarketPriceSDKType;
}
/**
 * QueryAllMarketPricesRequest is request type for the Query/Params
 * `AllMarketPrices` RPC method.
 */

export interface QueryAllMarketPricesRequest {
  pagination?: PageRequest;
}
/**
 * QueryAllMarketPricesRequest is request type for the Query/Params
 * `AllMarketPrices` RPC method.
 */

export interface QueryAllMarketPricesRequestSDKType {
  pagination?: PageRequestSDKType;
}
/**
 * QueryAllMarketPricesResponse is response type for the Query/Params
 * `AllMarketPrices` RPC method.
 */

export interface QueryAllMarketPricesResponse {
  marketPrices: MarketPrice[];
  pagination?: PageResponse;
}
/**
 * QueryAllMarketPricesResponse is response type for the Query/Params
 * `AllMarketPrices` RPC method.
 */

export interface QueryAllMarketPricesResponseSDKType {
  market_prices: MarketPriceSDKType[];
  pagination?: PageResponseSDKType;
}
/**
 * QueryMarketParamsRequest is request type for the Query/Params `MarketParams`
 * RPC method.
 */

export interface QueryMarketParamRequest {
  /**
   * QueryMarketParamsRequest is request type for the Query/Params `MarketParams`
   * RPC method.
   */
  id: number;
}
/**
 * QueryMarketParamsRequest is request type for the Query/Params `MarketParams`
 * RPC method.
 */

export interface QueryMarketParamRequestSDKType {
  /**
   * QueryMarketParamsRequest is request type for the Query/Params `MarketParams`
   * RPC method.
   */
  id: number;
}
/**
 * QueryMarketParamResponse is response type for the Query/Params `MarketParams`
 * RPC method.
 */

export interface QueryMarketParamResponse {
  marketParam?: MarketParam;
}
/**
 * QueryMarketParamResponse is response type for the Query/Params `MarketParams`
 * RPC method.
 */

export interface QueryMarketParamResponseSDKType {
  market_param?: MarketParamSDKType;
}
/**
 * QueryAllMarketParamsRequest is request type for the Query/Params
 * `AllMarketParams` RPC method.
 */

export interface QueryAllMarketParamsRequest {
  pagination?: PageRequest;
}
/**
 * QueryAllMarketParamsRequest is request type for the Query/Params
 * `AllMarketParams` RPC method.
 */

export interface QueryAllMarketParamsRequestSDKType {
  pagination?: PageRequestSDKType;
}
/**
 * QueryAllMarketParamsResponse is response type for the Query/Params
 * `AllMarketParams` RPC method.
 */

export interface QueryAllMarketParamsResponse {
  marketParams: MarketParam[];
  pagination?: PageResponse;
}
/**
 * QueryAllMarketParamsResponse is response type for the Query/Params
 * `AllMarketParams` RPC method.
 */

export interface QueryAllMarketParamsResponseSDKType {
  market_params: MarketParamSDKType[];
  pagination?: PageResponseSDKType;
}
/** QueryNextMarketIdRequest is request type for the Query/Params `NextMarketId` */

export interface QueryNextMarketIdRequest {}
/** QueryNextMarketIdRequest is request type for the Query/Params `NextMarketId` */

export interface QueryNextMarketIdRequestSDKType {}
/**
 * QueryNextMarketIdResponse is response type for the Query/Params
 * `NextMarketId`
 */

export interface QueryNextMarketIdResponse {
  /**
   * QueryNextMarketIdResponse is response type for the Query/Params
   * `NextMarketId`
   */
  nextMarketId: number;
}
/**
 * QueryNextMarketIdResponse is response type for the Query/Params
 * `NextMarketId`
 */

export interface QueryNextMarketIdResponseSDKType {
  /**
   * QueryNextMarketIdResponse is response type for the Query/Params
   * `NextMarketId`
   */
  next_market_id: number;
}

function createBaseQueryMarketPriceRequest(): QueryMarketPriceRequest {
  return {
    id: 0
  };
}

export const QueryMarketPriceRequest = {
  encode(message: QueryMarketPriceRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== 0) {
      writer.uint32(8).uint32(message.id);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryMarketPriceRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryMarketPriceRequest();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.id = reader.uint32();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QueryMarketPriceRequest>): QueryMarketPriceRequest {
    const message = createBaseQueryMarketPriceRequest();
    message.id = object.id ?? 0;
    return message;
  }

};

function createBaseQueryMarketPriceResponse(): QueryMarketPriceResponse {
  return {
    marketPrice: undefined
  };
}

export const QueryMarketPriceResponse = {
  encode(message: QueryMarketPriceResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.marketPrice !== undefined) {
      MarketPrice.encode(message.marketPrice, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryMarketPriceResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryMarketPriceResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.marketPrice = MarketPrice.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QueryMarketPriceResponse>): QueryMarketPriceResponse {
    const message = createBaseQueryMarketPriceResponse();
    message.marketPrice = object.marketPrice !== undefined && object.marketPrice !== null ? MarketPrice.fromPartial(object.marketPrice) : undefined;
    return message;
  }

};

function createBaseQueryAllMarketPricesRequest(): QueryAllMarketPricesRequest {
  return {
    pagination: undefined
  };
}

export const QueryAllMarketPricesRequest = {
  encode(message: QueryAllMarketPricesRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.pagination !== undefined) {
      PageRequest.encode(message.pagination, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryAllMarketPricesRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAllMarketPricesRequest();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.pagination = PageRequest.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QueryAllMarketPricesRequest>): QueryAllMarketPricesRequest {
    const message = createBaseQueryAllMarketPricesRequest();
    message.pagination = object.pagination !== undefined && object.pagination !== null ? PageRequest.fromPartial(object.pagination) : undefined;
    return message;
  }

};

function createBaseQueryAllMarketPricesResponse(): QueryAllMarketPricesResponse {
  return {
    marketPrices: [],
    pagination: undefined
  };
}

export const QueryAllMarketPricesResponse = {
  encode(message: QueryAllMarketPricesResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.marketPrices) {
      MarketPrice.encode(v!, writer.uint32(10).fork()).ldelim();
    }

    if (message.pagination !== undefined) {
      PageResponse.encode(message.pagination, writer.uint32(18).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryAllMarketPricesResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAllMarketPricesResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.marketPrices.push(MarketPrice.decode(reader, reader.uint32()));
          break;

        case 2:
          message.pagination = PageResponse.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QueryAllMarketPricesResponse>): QueryAllMarketPricesResponse {
    const message = createBaseQueryAllMarketPricesResponse();
    message.marketPrices = object.marketPrices?.map(e => MarketPrice.fromPartial(e)) || [];
    message.pagination = object.pagination !== undefined && object.pagination !== null ? PageResponse.fromPartial(object.pagination) : undefined;
    return message;
  }

};

function createBaseQueryMarketParamRequest(): QueryMarketParamRequest {
  return {
    id: 0
  };
}

export const QueryMarketParamRequest = {
  encode(message: QueryMarketParamRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== 0) {
      writer.uint32(8).uint32(message.id);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryMarketParamRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryMarketParamRequest();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.id = reader.uint32();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QueryMarketParamRequest>): QueryMarketParamRequest {
    const message = createBaseQueryMarketParamRequest();
    message.id = object.id ?? 0;
    return message;
  }

};

function createBaseQueryMarketParamResponse(): QueryMarketParamResponse {
  return {
    marketParam: undefined
  };
}

export const QueryMarketParamResponse = {
  encode(message: QueryMarketParamResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.marketParam !== undefined) {
      MarketParam.encode(message.marketParam, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryMarketParamResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryMarketParamResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.marketParam = MarketParam.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QueryMarketParamResponse>): QueryMarketParamResponse {
    const message = createBaseQueryMarketParamResponse();
    message.marketParam = object.marketParam !== undefined && object.marketParam !== null ? MarketParam.fromPartial(object.marketParam) : undefined;
    return message;
  }

};

function createBaseQueryAllMarketParamsRequest(): QueryAllMarketParamsRequest {
  return {
    pagination: undefined
  };
}

export const QueryAllMarketParamsRequest = {
  encode(message: QueryAllMarketParamsRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.pagination !== undefined) {
      PageRequest.encode(message.pagination, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryAllMarketParamsRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAllMarketParamsRequest();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.pagination = PageRequest.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QueryAllMarketParamsRequest>): QueryAllMarketParamsRequest {
    const message = createBaseQueryAllMarketParamsRequest();
    message.pagination = object.pagination !== undefined && object.pagination !== null ? PageRequest.fromPartial(object.pagination) : undefined;
    return message;
  }

};

function createBaseQueryAllMarketParamsResponse(): QueryAllMarketParamsResponse {
  return {
    marketParams: [],
    pagination: undefined
  };
}

export const QueryAllMarketParamsResponse = {
  encode(message: QueryAllMarketParamsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.marketParams) {
      MarketParam.encode(v!, writer.uint32(10).fork()).ldelim();
    }

    if (message.pagination !== undefined) {
      PageResponse.encode(message.pagination, writer.uint32(18).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryAllMarketParamsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAllMarketParamsResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.marketParams.push(MarketParam.decode(reader, reader.uint32()));
          break;

        case 2:
          message.pagination = PageResponse.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QueryAllMarketParamsResponse>): QueryAllMarketParamsResponse {
    const message = createBaseQueryAllMarketParamsResponse();
    message.marketParams = object.marketParams?.map(e => MarketParam.fromPartial(e)) || [];
    message.pagination = object.pagination !== undefined && object.pagination !== null ? PageResponse.fromPartial(object.pagination) : undefined;
    return message;
  }

};

function createBaseQueryNextMarketIdRequest(): QueryNextMarketIdRequest {
  return {};
}

export const QueryNextMarketIdRequest = {
  encode(_: QueryNextMarketIdRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryNextMarketIdRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryNextMarketIdRequest();

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

  fromPartial(_: DeepPartial<QueryNextMarketIdRequest>): QueryNextMarketIdRequest {
    const message = createBaseQueryNextMarketIdRequest();
    return message;
  }

};

function createBaseQueryNextMarketIdResponse(): QueryNextMarketIdResponse {
  return {
    nextMarketId: 0
  };
}

export const QueryNextMarketIdResponse = {
  encode(message: QueryNextMarketIdResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.nextMarketId !== 0) {
      writer.uint32(8).uint32(message.nextMarketId);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryNextMarketIdResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryNextMarketIdResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.nextMarketId = reader.uint32();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QueryNextMarketIdResponse>): QueryNextMarketIdResponse {
    const message = createBaseQueryNextMarketIdResponse();
    message.nextMarketId = object.nextMarketId ?? 0;
    return message;
  }

};