import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/**
 * PerpetualLeverageEntry represents a single perpetual leverage setting for
 * internal storage
 */

export interface PerpetualLeverageEntry {
  /** The perpetual ID (internal storage format) */
  perpetualId: number;
  /** The user selected IMF in parts per million */

  customImfPpm: number;
}
/**
 * PerpetualLeverageEntry represents a single perpetual leverage setting for
 * internal storage
 */

export interface PerpetualLeverageEntrySDKType {
  /** The perpetual ID (internal storage format) */
  perpetual_id: number;
  /** The user selected IMF in parts per million */

  custom_imf_ppm: number;
}
/** LeverageData represents the leverage settings for a subaccount */

export interface LeverageData {
  /** List of leverage entries for this subaccount */
  entries: PerpetualLeverageEntry[];
}
/** LeverageData represents the leverage settings for a subaccount */

export interface LeverageDataSDKType {
  /** List of leverage entries for this subaccount */
  entries: PerpetualLeverageEntrySDKType[];
}

function createBasePerpetualLeverageEntry(): PerpetualLeverageEntry {
  return {
    perpetualId: 0,
    customImfPpm: 0
  };
}

export const PerpetualLeverageEntry = {
  encode(message: PerpetualLeverageEntry, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.perpetualId !== 0) {
      writer.uint32(8).uint32(message.perpetualId);
    }

    if (message.customImfPpm !== 0) {
      writer.uint32(16).uint32(message.customImfPpm);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): PerpetualLeverageEntry {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBasePerpetualLeverageEntry();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.perpetualId = reader.uint32();
          break;

        case 2:
          message.customImfPpm = reader.uint32();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<PerpetualLeverageEntry>): PerpetualLeverageEntry {
    const message = createBasePerpetualLeverageEntry();
    message.perpetualId = object.perpetualId ?? 0;
    message.customImfPpm = object.customImfPpm ?? 0;
    return message;
  }

};

function createBaseLeverageData(): LeverageData {
  return {
    entries: []
  };
}

export const LeverageData = {
  encode(message: LeverageData, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.entries) {
      PerpetualLeverageEntry.encode(v!, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): LeverageData {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseLeverageData();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.entries.push(PerpetualLeverageEntry.decode(reader, reader.uint32()));
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<LeverageData>): LeverageData {
    const message = createBaseLeverageData();
    message.entries = object.entries?.map(e => PerpetualLeverageEntry.fromPartial(e)) || [];
    return message;
  }

};