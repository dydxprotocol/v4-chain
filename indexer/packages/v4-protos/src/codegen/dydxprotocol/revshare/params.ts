import { MarketMapperRevShareDetails, MarketMapperRevShareDetailsSDKType } from "./revshare";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/** MarketMappeRevenueShareParams represents params for the above message */

export interface MarketMapperRevenueShareParams {
  /** The address which will receive the revenue share payouts */
  address: string;
  /**
   * The fraction of the fees which will go to the above mentioned address.
   * In parts-per-million
   */

  revenueSharePpm: number;
  /**
   * This parameter defines how many days post market initiation will the
   * revenue share be applied for. After valid_days from market initiation
   * the revenue share goes down to 0
   */

  validDays: number;
}
/** MarketMappeRevenueShareParams represents params for the above message */

export interface MarketMapperRevenueShareParamsSDKType {
  /** The address which will receive the revenue share payouts */
  address: string;
  /**
   * The fraction of the fees which will go to the above mentioned address.
   * In parts-per-million
   */

  revenue_share_ppm: number;
  /**
   * This parameter defines how many days post market initiation will the
   * revenue share be applied for. After valid_days from market initiation
   * the revenue share goes down to 0
   */

  valid_days: number;
}
/**
 * MarketRevShareDetailsParams represents params for the
 * MsgSetMarketMapperRevShareDetailsForMarket message
 */

export interface MarketRevShareDetailsParams {
  /** Market Id */
  marketId: number;
  /** Market mapper rev share details for the market, e.g. expiration timestamp */

  marketMapperRevShareDetails?: MarketMapperRevShareDetails;
}
/**
 * MarketRevShareDetailsParams represents params for the
 * MsgSetMarketMapperRevShareDetailsForMarket message
 */

export interface MarketRevShareDetailsParamsSDKType {
  /** Market Id */
  market_id: number;
  /** Market mapper rev share details for the market, e.g. expiration timestamp */

  market_mapper_rev_share_details?: MarketMapperRevShareDetailsSDKType;
}

function createBaseMarketMapperRevenueShareParams(): MarketMapperRevenueShareParams {
  return {
    address: "",
    revenueSharePpm: 0,
    validDays: 0
  };
}

export const MarketMapperRevenueShareParams = {
  encode(message: MarketMapperRevenueShareParams, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.address !== "") {
      writer.uint32(10).string(message.address);
    }

    if (message.revenueSharePpm !== 0) {
      writer.uint32(16).uint32(message.revenueSharePpm);
    }

    if (message.validDays !== 0) {
      writer.uint32(24).uint32(message.validDays);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MarketMapperRevenueShareParams {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMarketMapperRevenueShareParams();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.address = reader.string();
          break;

        case 2:
          message.revenueSharePpm = reader.uint32();
          break;

        case 3:
          message.validDays = reader.uint32();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<MarketMapperRevenueShareParams>): MarketMapperRevenueShareParams {
    const message = createBaseMarketMapperRevenueShareParams();
    message.address = object.address ?? "";
    message.revenueSharePpm = object.revenueSharePpm ?? 0;
    message.validDays = object.validDays ?? 0;
    return message;
  }

};

function createBaseMarketRevShareDetailsParams(): MarketRevShareDetailsParams {
  return {
    marketId: 0,
    marketMapperRevShareDetails: undefined
  };
}

export const MarketRevShareDetailsParams = {
  encode(message: MarketRevShareDetailsParams, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.marketId !== 0) {
      writer.uint32(8).uint32(message.marketId);
    }

    if (message.marketMapperRevShareDetails !== undefined) {
      MarketMapperRevShareDetails.encode(message.marketMapperRevShareDetails, writer.uint32(18).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MarketRevShareDetailsParams {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMarketRevShareDetailsParams();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.marketId = reader.uint32();
          break;

        case 2:
          message.marketMapperRevShareDetails = MarketMapperRevShareDetails.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<MarketRevShareDetailsParams>): MarketRevShareDetailsParams {
    const message = createBaseMarketRevShareDetailsParams();
    message.marketId = object.marketId ?? 0;
    message.marketMapperRevShareDetails = object.marketMapperRevShareDetails !== undefined && object.marketMapperRevShareDetails !== null ? MarketMapperRevShareDetails.fromPartial(object.marketMapperRevShareDetails) : undefined;
    return message;
  }

};