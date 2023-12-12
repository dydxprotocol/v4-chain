import { PageRequest, PageRequestAmino, PageRequestSDKType, PageResponse, PageResponseAmino, PageResponseSDKType } from "../../cosmos/base/query/v1beta1/pagination";
import { MarketPrice, MarketPriceAmino, MarketPriceSDKType } from "./market_price";
import { MarketParam, MarketParamAmino, MarketParamSDKType } from "./market_param";
import { BinaryReader, BinaryWriter } from "../../binary";
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
export interface QueryMarketPriceRequestProtoMsg {
  typeUrl: "/dydxprotocol.prices.QueryMarketPriceRequest";
  value: Uint8Array;
}
/**
 * QueryMarketPriceRequest is request type for the Query/Params `MarketPrice`
 * RPC method.
 */
export interface QueryMarketPriceRequestAmino {
  /**
   * QueryMarketPriceRequest is request type for the Query/Params `MarketPrice`
   * RPC method.
   */
  id?: number;
}
export interface QueryMarketPriceRequestAminoMsg {
  type: "/dydxprotocol.prices.QueryMarketPriceRequest";
  value: QueryMarketPriceRequestAmino;
}
/**
 * QueryMarketPriceRequest is request type for the Query/Params `MarketPrice`
 * RPC method.
 */
export interface QueryMarketPriceRequestSDKType {
  id: number;
}
/**
 * QueryMarketPriceResponse is response type for the Query/Params `MarketPrice`
 * RPC method.
 */
export interface QueryMarketPriceResponse {
  marketPrice: MarketPrice;
}
export interface QueryMarketPriceResponseProtoMsg {
  typeUrl: "/dydxprotocol.prices.QueryMarketPriceResponse";
  value: Uint8Array;
}
/**
 * QueryMarketPriceResponse is response type for the Query/Params `MarketPrice`
 * RPC method.
 */
export interface QueryMarketPriceResponseAmino {
  market_price?: MarketPriceAmino;
}
export interface QueryMarketPriceResponseAminoMsg {
  type: "/dydxprotocol.prices.QueryMarketPriceResponse";
  value: QueryMarketPriceResponseAmino;
}
/**
 * QueryMarketPriceResponse is response type for the Query/Params `MarketPrice`
 * RPC method.
 */
export interface QueryMarketPriceResponseSDKType {
  market_price: MarketPriceSDKType;
}
/**
 * QueryAllMarketPricesRequest is request type for the Query/Params
 * `AllMarketPrices` RPC method.
 */
export interface QueryAllMarketPricesRequest {
  pagination?: PageRequest;
}
export interface QueryAllMarketPricesRequestProtoMsg {
  typeUrl: "/dydxprotocol.prices.QueryAllMarketPricesRequest";
  value: Uint8Array;
}
/**
 * QueryAllMarketPricesRequest is request type for the Query/Params
 * `AllMarketPrices` RPC method.
 */
export interface QueryAllMarketPricesRequestAmino {
  pagination?: PageRequestAmino;
}
export interface QueryAllMarketPricesRequestAminoMsg {
  type: "/dydxprotocol.prices.QueryAllMarketPricesRequest";
  value: QueryAllMarketPricesRequestAmino;
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
export interface QueryAllMarketPricesResponseProtoMsg {
  typeUrl: "/dydxprotocol.prices.QueryAllMarketPricesResponse";
  value: Uint8Array;
}
/**
 * QueryAllMarketPricesResponse is response type for the Query/Params
 * `AllMarketPrices` RPC method.
 */
export interface QueryAllMarketPricesResponseAmino {
  market_prices?: MarketPriceAmino[];
  pagination?: PageResponseAmino;
}
export interface QueryAllMarketPricesResponseAminoMsg {
  type: "/dydxprotocol.prices.QueryAllMarketPricesResponse";
  value: QueryAllMarketPricesResponseAmino;
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
export interface QueryMarketParamRequestProtoMsg {
  typeUrl: "/dydxprotocol.prices.QueryMarketParamRequest";
  value: Uint8Array;
}
/**
 * QueryMarketParamsRequest is request type for the Query/Params `MarketParams`
 * RPC method.
 */
export interface QueryMarketParamRequestAmino {
  /**
   * QueryMarketParamsRequest is request type for the Query/Params `MarketParams`
   * RPC method.
   */
  id?: number;
}
export interface QueryMarketParamRequestAminoMsg {
  type: "/dydxprotocol.prices.QueryMarketParamRequest";
  value: QueryMarketParamRequestAmino;
}
/**
 * QueryMarketParamsRequest is request type for the Query/Params `MarketParams`
 * RPC method.
 */
export interface QueryMarketParamRequestSDKType {
  id: number;
}
/**
 * QueryMarketParamResponse is response type for the Query/Params `MarketParams`
 * RPC method.
 */
export interface QueryMarketParamResponse {
  marketParam: MarketParam;
}
export interface QueryMarketParamResponseProtoMsg {
  typeUrl: "/dydxprotocol.prices.QueryMarketParamResponse";
  value: Uint8Array;
}
/**
 * QueryMarketParamResponse is response type for the Query/Params `MarketParams`
 * RPC method.
 */
export interface QueryMarketParamResponseAmino {
  market_param?: MarketParamAmino;
}
export interface QueryMarketParamResponseAminoMsg {
  type: "/dydxprotocol.prices.QueryMarketParamResponse";
  value: QueryMarketParamResponseAmino;
}
/**
 * QueryMarketParamResponse is response type for the Query/Params `MarketParams`
 * RPC method.
 */
export interface QueryMarketParamResponseSDKType {
  market_param: MarketParamSDKType;
}
/**
 * QueryAllMarketParamsRequest is request type for the Query/Params
 * `AllMarketParams` RPC method.
 */
export interface QueryAllMarketParamsRequest {
  pagination?: PageRequest;
}
export interface QueryAllMarketParamsRequestProtoMsg {
  typeUrl: "/dydxprotocol.prices.QueryAllMarketParamsRequest";
  value: Uint8Array;
}
/**
 * QueryAllMarketParamsRequest is request type for the Query/Params
 * `AllMarketParams` RPC method.
 */
export interface QueryAllMarketParamsRequestAmino {
  pagination?: PageRequestAmino;
}
export interface QueryAllMarketParamsRequestAminoMsg {
  type: "/dydxprotocol.prices.QueryAllMarketParamsRequest";
  value: QueryAllMarketParamsRequestAmino;
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
export interface QueryAllMarketParamsResponseProtoMsg {
  typeUrl: "/dydxprotocol.prices.QueryAllMarketParamsResponse";
  value: Uint8Array;
}
/**
 * QueryAllMarketParamsResponse is response type for the Query/Params
 * `AllMarketParams` RPC method.
 */
export interface QueryAllMarketParamsResponseAmino {
  market_params?: MarketParamAmino[];
  pagination?: PageResponseAmino;
}
export interface QueryAllMarketParamsResponseAminoMsg {
  type: "/dydxprotocol.prices.QueryAllMarketParamsResponse";
  value: QueryAllMarketParamsResponseAmino;
}
/**
 * QueryAllMarketParamsResponse is response type for the Query/Params
 * `AllMarketParams` RPC method.
 */
export interface QueryAllMarketParamsResponseSDKType {
  market_params: MarketParamSDKType[];
  pagination?: PageResponseSDKType;
}
function createBaseQueryMarketPriceRequest(): QueryMarketPriceRequest {
  return {
    id: 0
  };
}
export const QueryMarketPriceRequest = {
  typeUrl: "/dydxprotocol.prices.QueryMarketPriceRequest",
  encode(message: QueryMarketPriceRequest, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.id !== 0) {
      writer.uint32(8).uint32(message.id);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryMarketPriceRequest {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
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
  fromPartial(object: Partial<QueryMarketPriceRequest>): QueryMarketPriceRequest {
    const message = createBaseQueryMarketPriceRequest();
    message.id = object.id ?? 0;
    return message;
  },
  fromAmino(object: QueryMarketPriceRequestAmino): QueryMarketPriceRequest {
    const message = createBaseQueryMarketPriceRequest();
    if (object.id !== undefined && object.id !== null) {
      message.id = object.id;
    }
    return message;
  },
  toAmino(message: QueryMarketPriceRequest): QueryMarketPriceRequestAmino {
    const obj: any = {};
    obj.id = message.id;
    return obj;
  },
  fromAminoMsg(object: QueryMarketPriceRequestAminoMsg): QueryMarketPriceRequest {
    return QueryMarketPriceRequest.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryMarketPriceRequestProtoMsg): QueryMarketPriceRequest {
    return QueryMarketPriceRequest.decode(message.value);
  },
  toProto(message: QueryMarketPriceRequest): Uint8Array {
    return QueryMarketPriceRequest.encode(message).finish();
  },
  toProtoMsg(message: QueryMarketPriceRequest): QueryMarketPriceRequestProtoMsg {
    return {
      typeUrl: "/dydxprotocol.prices.QueryMarketPriceRequest",
      value: QueryMarketPriceRequest.encode(message).finish()
    };
  }
};
function createBaseQueryMarketPriceResponse(): QueryMarketPriceResponse {
  return {
    marketPrice: MarketPrice.fromPartial({})
  };
}
export const QueryMarketPriceResponse = {
  typeUrl: "/dydxprotocol.prices.QueryMarketPriceResponse",
  encode(message: QueryMarketPriceResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.marketPrice !== undefined) {
      MarketPrice.encode(message.marketPrice, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryMarketPriceResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
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
  fromPartial(object: Partial<QueryMarketPriceResponse>): QueryMarketPriceResponse {
    const message = createBaseQueryMarketPriceResponse();
    message.marketPrice = object.marketPrice !== undefined && object.marketPrice !== null ? MarketPrice.fromPartial(object.marketPrice) : undefined;
    return message;
  },
  fromAmino(object: QueryMarketPriceResponseAmino): QueryMarketPriceResponse {
    const message = createBaseQueryMarketPriceResponse();
    if (object.market_price !== undefined && object.market_price !== null) {
      message.marketPrice = MarketPrice.fromAmino(object.market_price);
    }
    return message;
  },
  toAmino(message: QueryMarketPriceResponse): QueryMarketPriceResponseAmino {
    const obj: any = {};
    obj.market_price = message.marketPrice ? MarketPrice.toAmino(message.marketPrice) : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryMarketPriceResponseAminoMsg): QueryMarketPriceResponse {
    return QueryMarketPriceResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryMarketPriceResponseProtoMsg): QueryMarketPriceResponse {
    return QueryMarketPriceResponse.decode(message.value);
  },
  toProto(message: QueryMarketPriceResponse): Uint8Array {
    return QueryMarketPriceResponse.encode(message).finish();
  },
  toProtoMsg(message: QueryMarketPriceResponse): QueryMarketPriceResponseProtoMsg {
    return {
      typeUrl: "/dydxprotocol.prices.QueryMarketPriceResponse",
      value: QueryMarketPriceResponse.encode(message).finish()
    };
  }
};
function createBaseQueryAllMarketPricesRequest(): QueryAllMarketPricesRequest {
  return {
    pagination: undefined
  };
}
export const QueryAllMarketPricesRequest = {
  typeUrl: "/dydxprotocol.prices.QueryAllMarketPricesRequest",
  encode(message: QueryAllMarketPricesRequest, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.pagination !== undefined) {
      PageRequest.encode(message.pagination, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryAllMarketPricesRequest {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
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
  fromPartial(object: Partial<QueryAllMarketPricesRequest>): QueryAllMarketPricesRequest {
    const message = createBaseQueryAllMarketPricesRequest();
    message.pagination = object.pagination !== undefined && object.pagination !== null ? PageRequest.fromPartial(object.pagination) : undefined;
    return message;
  },
  fromAmino(object: QueryAllMarketPricesRequestAmino): QueryAllMarketPricesRequest {
    const message = createBaseQueryAllMarketPricesRequest();
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageRequest.fromAmino(object.pagination);
    }
    return message;
  },
  toAmino(message: QueryAllMarketPricesRequest): QueryAllMarketPricesRequestAmino {
    const obj: any = {};
    obj.pagination = message.pagination ? PageRequest.toAmino(message.pagination) : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryAllMarketPricesRequestAminoMsg): QueryAllMarketPricesRequest {
    return QueryAllMarketPricesRequest.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryAllMarketPricesRequestProtoMsg): QueryAllMarketPricesRequest {
    return QueryAllMarketPricesRequest.decode(message.value);
  },
  toProto(message: QueryAllMarketPricesRequest): Uint8Array {
    return QueryAllMarketPricesRequest.encode(message).finish();
  },
  toProtoMsg(message: QueryAllMarketPricesRequest): QueryAllMarketPricesRequestProtoMsg {
    return {
      typeUrl: "/dydxprotocol.prices.QueryAllMarketPricesRequest",
      value: QueryAllMarketPricesRequest.encode(message).finish()
    };
  }
};
function createBaseQueryAllMarketPricesResponse(): QueryAllMarketPricesResponse {
  return {
    marketPrices: [],
    pagination: undefined
  };
}
export const QueryAllMarketPricesResponse = {
  typeUrl: "/dydxprotocol.prices.QueryAllMarketPricesResponse",
  encode(message: QueryAllMarketPricesResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    for (const v of message.marketPrices) {
      MarketPrice.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    if (message.pagination !== undefined) {
      PageResponse.encode(message.pagination, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryAllMarketPricesResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
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
  fromPartial(object: Partial<QueryAllMarketPricesResponse>): QueryAllMarketPricesResponse {
    const message = createBaseQueryAllMarketPricesResponse();
    message.marketPrices = object.marketPrices?.map(e => MarketPrice.fromPartial(e)) || [];
    message.pagination = object.pagination !== undefined && object.pagination !== null ? PageResponse.fromPartial(object.pagination) : undefined;
    return message;
  },
  fromAmino(object: QueryAllMarketPricesResponseAmino): QueryAllMarketPricesResponse {
    const message = createBaseQueryAllMarketPricesResponse();
    message.marketPrices = object.market_prices?.map(e => MarketPrice.fromAmino(e)) || [];
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageResponse.fromAmino(object.pagination);
    }
    return message;
  },
  toAmino(message: QueryAllMarketPricesResponse): QueryAllMarketPricesResponseAmino {
    const obj: any = {};
    if (message.marketPrices) {
      obj.market_prices = message.marketPrices.map(e => e ? MarketPrice.toAmino(e) : undefined);
    } else {
      obj.market_prices = [];
    }
    obj.pagination = message.pagination ? PageResponse.toAmino(message.pagination) : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryAllMarketPricesResponseAminoMsg): QueryAllMarketPricesResponse {
    return QueryAllMarketPricesResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryAllMarketPricesResponseProtoMsg): QueryAllMarketPricesResponse {
    return QueryAllMarketPricesResponse.decode(message.value);
  },
  toProto(message: QueryAllMarketPricesResponse): Uint8Array {
    return QueryAllMarketPricesResponse.encode(message).finish();
  },
  toProtoMsg(message: QueryAllMarketPricesResponse): QueryAllMarketPricesResponseProtoMsg {
    return {
      typeUrl: "/dydxprotocol.prices.QueryAllMarketPricesResponse",
      value: QueryAllMarketPricesResponse.encode(message).finish()
    };
  }
};
function createBaseQueryMarketParamRequest(): QueryMarketParamRequest {
  return {
    id: 0
  };
}
export const QueryMarketParamRequest = {
  typeUrl: "/dydxprotocol.prices.QueryMarketParamRequest",
  encode(message: QueryMarketParamRequest, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.id !== 0) {
      writer.uint32(8).uint32(message.id);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryMarketParamRequest {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
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
  fromPartial(object: Partial<QueryMarketParamRequest>): QueryMarketParamRequest {
    const message = createBaseQueryMarketParamRequest();
    message.id = object.id ?? 0;
    return message;
  },
  fromAmino(object: QueryMarketParamRequestAmino): QueryMarketParamRequest {
    const message = createBaseQueryMarketParamRequest();
    if (object.id !== undefined && object.id !== null) {
      message.id = object.id;
    }
    return message;
  },
  toAmino(message: QueryMarketParamRequest): QueryMarketParamRequestAmino {
    const obj: any = {};
    obj.id = message.id;
    return obj;
  },
  fromAminoMsg(object: QueryMarketParamRequestAminoMsg): QueryMarketParamRequest {
    return QueryMarketParamRequest.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryMarketParamRequestProtoMsg): QueryMarketParamRequest {
    return QueryMarketParamRequest.decode(message.value);
  },
  toProto(message: QueryMarketParamRequest): Uint8Array {
    return QueryMarketParamRequest.encode(message).finish();
  },
  toProtoMsg(message: QueryMarketParamRequest): QueryMarketParamRequestProtoMsg {
    return {
      typeUrl: "/dydxprotocol.prices.QueryMarketParamRequest",
      value: QueryMarketParamRequest.encode(message).finish()
    };
  }
};
function createBaseQueryMarketParamResponse(): QueryMarketParamResponse {
  return {
    marketParam: MarketParam.fromPartial({})
  };
}
export const QueryMarketParamResponse = {
  typeUrl: "/dydxprotocol.prices.QueryMarketParamResponse",
  encode(message: QueryMarketParamResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.marketParam !== undefined) {
      MarketParam.encode(message.marketParam, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryMarketParamResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
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
  fromPartial(object: Partial<QueryMarketParamResponse>): QueryMarketParamResponse {
    const message = createBaseQueryMarketParamResponse();
    message.marketParam = object.marketParam !== undefined && object.marketParam !== null ? MarketParam.fromPartial(object.marketParam) : undefined;
    return message;
  },
  fromAmino(object: QueryMarketParamResponseAmino): QueryMarketParamResponse {
    const message = createBaseQueryMarketParamResponse();
    if (object.market_param !== undefined && object.market_param !== null) {
      message.marketParam = MarketParam.fromAmino(object.market_param);
    }
    return message;
  },
  toAmino(message: QueryMarketParamResponse): QueryMarketParamResponseAmino {
    const obj: any = {};
    obj.market_param = message.marketParam ? MarketParam.toAmino(message.marketParam) : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryMarketParamResponseAminoMsg): QueryMarketParamResponse {
    return QueryMarketParamResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryMarketParamResponseProtoMsg): QueryMarketParamResponse {
    return QueryMarketParamResponse.decode(message.value);
  },
  toProto(message: QueryMarketParamResponse): Uint8Array {
    return QueryMarketParamResponse.encode(message).finish();
  },
  toProtoMsg(message: QueryMarketParamResponse): QueryMarketParamResponseProtoMsg {
    return {
      typeUrl: "/dydxprotocol.prices.QueryMarketParamResponse",
      value: QueryMarketParamResponse.encode(message).finish()
    };
  }
};
function createBaseQueryAllMarketParamsRequest(): QueryAllMarketParamsRequest {
  return {
    pagination: undefined
  };
}
export const QueryAllMarketParamsRequest = {
  typeUrl: "/dydxprotocol.prices.QueryAllMarketParamsRequest",
  encode(message: QueryAllMarketParamsRequest, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.pagination !== undefined) {
      PageRequest.encode(message.pagination, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryAllMarketParamsRequest {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
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
  fromPartial(object: Partial<QueryAllMarketParamsRequest>): QueryAllMarketParamsRequest {
    const message = createBaseQueryAllMarketParamsRequest();
    message.pagination = object.pagination !== undefined && object.pagination !== null ? PageRequest.fromPartial(object.pagination) : undefined;
    return message;
  },
  fromAmino(object: QueryAllMarketParamsRequestAmino): QueryAllMarketParamsRequest {
    const message = createBaseQueryAllMarketParamsRequest();
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageRequest.fromAmino(object.pagination);
    }
    return message;
  },
  toAmino(message: QueryAllMarketParamsRequest): QueryAllMarketParamsRequestAmino {
    const obj: any = {};
    obj.pagination = message.pagination ? PageRequest.toAmino(message.pagination) : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryAllMarketParamsRequestAminoMsg): QueryAllMarketParamsRequest {
    return QueryAllMarketParamsRequest.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryAllMarketParamsRequestProtoMsg): QueryAllMarketParamsRequest {
    return QueryAllMarketParamsRequest.decode(message.value);
  },
  toProto(message: QueryAllMarketParamsRequest): Uint8Array {
    return QueryAllMarketParamsRequest.encode(message).finish();
  },
  toProtoMsg(message: QueryAllMarketParamsRequest): QueryAllMarketParamsRequestProtoMsg {
    return {
      typeUrl: "/dydxprotocol.prices.QueryAllMarketParamsRequest",
      value: QueryAllMarketParamsRequest.encode(message).finish()
    };
  }
};
function createBaseQueryAllMarketParamsResponse(): QueryAllMarketParamsResponse {
  return {
    marketParams: [],
    pagination: undefined
  };
}
export const QueryAllMarketParamsResponse = {
  typeUrl: "/dydxprotocol.prices.QueryAllMarketParamsResponse",
  encode(message: QueryAllMarketParamsResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    for (const v of message.marketParams) {
      MarketParam.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    if (message.pagination !== undefined) {
      PageResponse.encode(message.pagination, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryAllMarketParamsResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
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
  fromPartial(object: Partial<QueryAllMarketParamsResponse>): QueryAllMarketParamsResponse {
    const message = createBaseQueryAllMarketParamsResponse();
    message.marketParams = object.marketParams?.map(e => MarketParam.fromPartial(e)) || [];
    message.pagination = object.pagination !== undefined && object.pagination !== null ? PageResponse.fromPartial(object.pagination) : undefined;
    return message;
  },
  fromAmino(object: QueryAllMarketParamsResponseAmino): QueryAllMarketParamsResponse {
    const message = createBaseQueryAllMarketParamsResponse();
    message.marketParams = object.market_params?.map(e => MarketParam.fromAmino(e)) || [];
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageResponse.fromAmino(object.pagination);
    }
    return message;
  },
  toAmino(message: QueryAllMarketParamsResponse): QueryAllMarketParamsResponseAmino {
    const obj: any = {};
    if (message.marketParams) {
      obj.market_params = message.marketParams.map(e => e ? MarketParam.toAmino(e) : undefined);
    } else {
      obj.market_params = [];
    }
    obj.pagination = message.pagination ? PageResponse.toAmino(message.pagination) : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryAllMarketParamsResponseAminoMsg): QueryAllMarketParamsResponse {
    return QueryAllMarketParamsResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryAllMarketParamsResponseProtoMsg): QueryAllMarketParamsResponse {
    return QueryAllMarketParamsResponse.decode(message.value);
  },
  toProto(message: QueryAllMarketParamsResponse): Uint8Array {
    return QueryAllMarketParamsResponse.encode(message).finish();
  },
  toProtoMsg(message: QueryAllMarketParamsResponse): QueryAllMarketParamsResponseProtoMsg {
    return {
      typeUrl: "/dydxprotocol.prices.QueryAllMarketParamsResponse",
      value: QueryAllMarketParamsResponse.encode(message).finish()
    };
  }
};