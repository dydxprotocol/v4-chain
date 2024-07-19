import { Timestamp } from "../../google/protobuf/timestamp";
import { Duration, DurationSDKType } from "../../google/protobuf/duration";
import * as _m0 from "protobufjs/minimal";
import { toTimestamp, fromTimestamp, DeepPartial } from "../../helpers";
/** BlockInfo stores information about a block */

export interface BlockInfo {
  height: number;
  timestamp?: Date;
}
/** BlockInfo stores information about a block */

export interface BlockInfoSDKType {
  height: number;
  timestamp?: Date;
}
/** AllDowntimeInfo stores information for all downtime durations. */

export interface AllDowntimeInfo {
  /**
   * The downtime information for each tracked duration. Sorted by duration,
   * ascending. (i.e. the same order as they appear in DowntimeParams).
   */
  infos: AllDowntimeInfo_DowntimeInfo[];
}
/** AllDowntimeInfo stores information for all downtime durations. */

export interface AllDowntimeInfoSDKType {
  /**
   * The downtime information for each tracked duration. Sorted by duration,
   * ascending. (i.e. the same order as they appear in DowntimeParams).
   */
  infos: AllDowntimeInfo_DowntimeInfoSDKType[];
}
/**
 * Stores information about downtime. block_info corresponds to the most
 * recent block at which a downtime occurred.
 */

export interface AllDowntimeInfo_DowntimeInfo {
  duration?: Duration;
  blockInfo?: BlockInfo;
}
/**
 * Stores information about downtime. block_info corresponds to the most
 * recent block at which a downtime occurred.
 */

export interface AllDowntimeInfo_DowntimeInfoSDKType {
  duration?: DurationSDKType;
  block_info?: BlockInfoSDKType;
}

function createBaseBlockInfo(): BlockInfo {
  return {
    height: 0,
    timestamp: undefined
  };
}

export const BlockInfo = {
  encode(message: BlockInfo, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.height !== 0) {
      writer.uint32(8).uint32(message.height);
    }

    if (message.timestamp !== undefined) {
      Timestamp.encode(toTimestamp(message.timestamp), writer.uint32(18).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): BlockInfo {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseBlockInfo();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.height = reader.uint32();
          break;

        case 2:
          message.timestamp = fromTimestamp(Timestamp.decode(reader, reader.uint32()));
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<BlockInfo>): BlockInfo {
    const message = createBaseBlockInfo();
    message.height = object.height ?? 0;
    message.timestamp = object.timestamp ?? undefined;
    return message;
  }

};

function createBaseAllDowntimeInfo(): AllDowntimeInfo {
  return {
    infos: []
  };
}

export const AllDowntimeInfo = {
  encode(message: AllDowntimeInfo, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.infos) {
      AllDowntimeInfo_DowntimeInfo.encode(v!, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AllDowntimeInfo {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAllDowntimeInfo();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.infos.push(AllDowntimeInfo_DowntimeInfo.decode(reader, reader.uint32()));
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<AllDowntimeInfo>): AllDowntimeInfo {
    const message = createBaseAllDowntimeInfo();
    message.infos = object.infos?.map(e => AllDowntimeInfo_DowntimeInfo.fromPartial(e)) || [];
    return message;
  }

};

function createBaseAllDowntimeInfo_DowntimeInfo(): AllDowntimeInfo_DowntimeInfo {
  return {
    duration: undefined,
    blockInfo: undefined
  };
}

export const AllDowntimeInfo_DowntimeInfo = {
  encode(message: AllDowntimeInfo_DowntimeInfo, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.duration !== undefined) {
      Duration.encode(message.duration, writer.uint32(10).fork()).ldelim();
    }

    if (message.blockInfo !== undefined) {
      BlockInfo.encode(message.blockInfo, writer.uint32(18).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AllDowntimeInfo_DowntimeInfo {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAllDowntimeInfo_DowntimeInfo();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.duration = Duration.decode(reader, reader.uint32());
          break;

        case 2:
          message.blockInfo = BlockInfo.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<AllDowntimeInfo_DowntimeInfo>): AllDowntimeInfo_DowntimeInfo {
    const message = createBaseAllDowntimeInfo_DowntimeInfo();
    message.duration = object.duration !== undefined && object.duration !== null ? Duration.fromPartial(object.duration) : undefined;
    message.blockInfo = object.blockInfo !== undefined && object.blockInfo !== null ? BlockInfo.fromPartial(object.blockInfo) : undefined;
    return message;
  }

};