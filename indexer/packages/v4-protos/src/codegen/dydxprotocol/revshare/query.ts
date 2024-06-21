import { MarketMapperRevShareDetails, MarketMapperRevShareDetailsSDKType } from "./revshare";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
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