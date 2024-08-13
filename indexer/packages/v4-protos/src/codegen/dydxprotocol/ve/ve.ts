import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
export interface DaemonVoteExtension_PricesEntry {
  key: number;
  value?: PricePair;
}
export interface DaemonVoteExtension_PricesEntrySDKType {
  key: number;
  value?: PricePairSDKType;
}
/** Daemon VoteExtension defines the vote extension structure for daemon prices. */

export interface DaemonVoteExtension {
  /** Prices defines a map of marketId -> PricePair. */
  prices?: {
    [key: number]: DaemonVoteExtension_PricePair;
  };
}
/** Daemon VoteExtension defines the vote extension structure for daemon prices. */

export interface DaemonVoteExtensionSDKType {
  /** Prices defines a map of marketId -> PricePair. */
  prices?: {
    [key: number]: DaemonVoteExtension_PricePairSDKType;
  };
}
/** PricePair defines a pair of prices for a market. */

export interface DaemonVoteExtension_PricePair {
  /** Plain oracle price (used for funding rates) */
  spotPrice: Uint8Array;
  /** Funding rate weighted price (used for pnl and liquidations) */

  pnlPrice: Uint8Array;
}
/** PricePair defines a pair of prices for a market. */

export interface DaemonVoteExtension_PricePairSDKType {
  /** Plain oracle price (used for funding rates) */
  spot_price: Uint8Array;
  /** Funding rate weighted price (used for pnl and liquidations) */

  pnl_price: Uint8Array;
}

function createBaseDaemonVoteExtension_PricesEntry(): DaemonVoteExtension_PricesEntry {
  return {
    key: 0,
    value: undefined
  };
}

export const DaemonVoteExtension_PricesEntry = {
  encode(message: DaemonVoteExtension_PricesEntry, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.key !== 0) {
      writer.uint32(8).uint32(message.key);
    }

    if (message.value !== undefined) {
      PricePair.encode(message.value, writer.uint32(18).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DaemonVoteExtension_PricesEntry {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDaemonVoteExtension_PricesEntry();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.key = reader.uint32();
          break;

        case 2:
          message.value = PricePair.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<DaemonVoteExtension_PricesEntry>): DaemonVoteExtension_PricesEntry {
    const message = createBaseDaemonVoteExtension_PricesEntry();
    message.key = object.key ?? 0;
    message.value = object.value !== undefined && object.value !== null ? PricePair.fromPartial(object.value) : undefined;
    return message;
  }

};

function createBaseDaemonVoteExtension(): DaemonVoteExtension {
  return {
    prices: {}
  };
}

export const DaemonVoteExtension = {
  encode(message: DaemonVoteExtension, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    Object.entries(message.prices).forEach(([key, value]) => {
      DaemonVoteExtension_PricesEntry.encode({
        key: (key as any),
        value
      }, writer.uint32(10).fork()).ldelim();
    });
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
          const entry1 = DaemonVoteExtension_PricesEntry.decode(reader, reader.uint32());

          if (entry1.value !== undefined) {
            message.prices[entry1.key] = entry1.value;
          }

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
    message.prices = Object.entries(object.prices ?? {}).reduce<{
      [key: number]: PricePair;
    }>((acc, [key, value]) => {
      if (value !== undefined) {
        acc[Number(key)] = PricePair.fromPartial(value);
      }

      return acc;
    }, {});
    return message;
  }

};

function createBaseDaemonVoteExtension_PricePair(): DaemonVoteExtension_PricePair {
  return {
    spotPrice: new Uint8Array(),
    pnlPrice: new Uint8Array()
  };
}

export const DaemonVoteExtension_PricePair = {
  encode(message: DaemonVoteExtension_PricePair, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.spotPrice.length !== 0) {
      writer.uint32(10).bytes(message.spotPrice);
    }

    if (message.pnlPrice.length !== 0) {
      writer.uint32(18).bytes(message.pnlPrice);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DaemonVoteExtension_PricePair {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDaemonVoteExtension_PricePair();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.spotPrice = reader.bytes();
          break;

        case 2:
          message.pnlPrice = reader.bytes();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<DaemonVoteExtension_PricePair>): DaemonVoteExtension_PricePair {
    const message = createBaseDaemonVoteExtension_PricePair();
    message.spotPrice = object.spotPrice ?? new Uint8Array();
    message.pnlPrice = object.pnlPrice ?? new Uint8Array();
    return message;
  }

};