import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/** PricePair defines a pair of prices for a market. */

export interface PricePair {
  marketId: number;
  /** Plain oracle price (used for funding rates) */

  spotPrice: Uint8Array;
  /** Funding rate weighted price (used for pnl and liquidations) */

  pnlPrice: Uint8Array;
}
/** PricePair defines a pair of prices for a market. */

export interface PricePairSDKType {
  market_id: number;
  /** Plain oracle price (used for funding rates) */

  spot_price: Uint8Array;
  /** Funding rate weighted price (used for pnl and liquidations) */

  pnl_price: Uint8Array;
}
/** Daemon VoteExtension defines the vote extension structure for daemon prices. */

export interface DaemonVoteExtension {
  /** Prices defines a map of marketId -> PricePair. */
  prices: PricePair[];
  /** sDaiConversionRate defines the conversion rate for sDAI. */

  sDaiConversionRate: string;
}
/** Daemon VoteExtension defines the vote extension structure for daemon prices. */

export interface DaemonVoteExtensionSDKType {
  /** Prices defines a map of marketId -> PricePair. */
  prices: PricePairSDKType[];
  /** sDaiConversionRate defines the conversion rate for sDAI. */

  sDaiConversionRate: string;
}

function createBasePricePair(): PricePair {
  return {
    marketId: 0,
    spotPrice: new Uint8Array(),
    pnlPrice: new Uint8Array()
  };
}

export const PricePair = {
  encode(message: PricePair, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.marketId !== 0) {
      writer.uint32(8).uint32(message.marketId);
    }

    if (message.spotPrice.length !== 0) {
      writer.uint32(18).bytes(message.spotPrice);
    }

    if (message.pnlPrice.length !== 0) {
      writer.uint32(26).bytes(message.pnlPrice);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): PricePair {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBasePricePair();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.marketId = reader.uint32();
          break;

        case 2:
          message.spotPrice = reader.bytes();
          break;

        case 3:
          message.pnlPrice = reader.bytes();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<PricePair>): PricePair {
    const message = createBasePricePair();
    message.marketId = object.marketId ?? 0;
    message.spotPrice = object.spotPrice ?? new Uint8Array();
    message.pnlPrice = object.pnlPrice ?? new Uint8Array();
    return message;
  }

};

function createBaseDaemonVoteExtension(): DaemonVoteExtension {
  return {
    prices: [],
    sDaiConversionRate: ""
  };
}

export const DaemonVoteExtension = {
  encode(message: DaemonVoteExtension, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.prices) {
      PricePair.encode(v!, writer.uint32(10).fork()).ldelim();
    }

    if (message.sDaiConversionRate !== "") {
      writer.uint32(18).string(message.sDaiConversionRate);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DaemonVoteExtension {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDaemonVoteExtension();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.prices.push(PricePair.decode(reader, reader.uint32()));
          break;

        case 2:
          message.sDaiConversionRate = reader.string();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<DaemonVoteExtension>): DaemonVoteExtension {
    const message = createBaseDaemonVoteExtension();
    message.prices = object.prices?.map(e => PricePair.fromPartial(e)) || [];
    message.sDaiConversionRate = object.sDaiConversionRate ?? "";
    return message;
  }

};