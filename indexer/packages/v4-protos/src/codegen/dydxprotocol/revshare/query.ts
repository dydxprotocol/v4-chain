import { MarketMapperRevenueShareParams, MarketMapperRevenueShareParamsSDKType } from "./params";
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