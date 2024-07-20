import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
export interface DaemonVoteExtension_PricesEntry {
  key: number;
  value: Uint8Array;
}
export interface DaemonVoteExtension_PricesEntrySDKType {
  key: number;
  value: Uint8Array;
}
/** Daemon VoteExtension defines the vote extension structure for daemon prices. */

export interface DaemonVoteExtension {
  /**
   * Prices defines a map of marketId -> price.Bytes() . i.e. 1 ->
   * 0x123.. (bytes). Notice the `id` function is determined by the
   * `marketParams` used in the VoteExtensionHandler.
   */
  prices: {
    [key: number]: Uint8Array;
  };
}
/** Daemon VoteExtension defines the vote extension structure for daemon prices. */

export interface DaemonVoteExtensionSDKType {
  /**
   * Prices defines a map of marketId -> price.Bytes() . i.e. 1 ->
   * 0x123.. (bytes). Notice the `id` function is determined by the
   * `marketParams` used in the VoteExtensionHandler.
   */
  prices: {
    [key: number]: Uint8Array;
  };
}

function createBaseDaemonVoteExtension_PricesEntry(): DaemonVoteExtension_PricesEntry {
  return {
    key: 0,
    value: new Uint8Array()
  };
}

export const DaemonVoteExtension_PricesEntry = {
  encode(message: DaemonVoteExtension_PricesEntry, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.key !== 0) {
      writer.uint32(8).uint32(message.key);
    }

    if (message.value.length !== 0) {
      writer.uint32(18).bytes(message.value);
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
          message.value = reader.bytes();
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
    message.value = object.value ?? new Uint8Array();
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
      [key: number]: bytes;
    }>((acc, [key, value]) => {
      if (value !== undefined) {
        acc[Number(key)] = bytes.fromPartial(value);
      }

      return acc;
    }, {});
    return message;
  }

};