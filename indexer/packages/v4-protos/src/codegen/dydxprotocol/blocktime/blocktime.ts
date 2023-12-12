import { Timestamp } from "../../google/protobuf/timestamp";
import { Duration, DurationAmino, DurationSDKType } from "../../google/protobuf/duration";
import { BinaryReader, BinaryWriter } from "../../binary";
import { toTimestamp, fromTimestamp } from "../../helpers";
/** BlockInfo stores information about a block */
export interface BlockInfo {
  height: number;
  timestamp: Date;
}
export interface BlockInfoProtoMsg {
  typeUrl: "/dydxprotocol.blocktime.BlockInfo";
  value: Uint8Array;
}
/** BlockInfo stores information about a block */
export interface BlockInfoAmino {
  height?: number;
  timestamp?: string;
}
export interface BlockInfoAminoMsg {
  type: "/dydxprotocol.blocktime.BlockInfo";
  value: BlockInfoAmino;
}
/** BlockInfo stores information about a block */
export interface BlockInfoSDKType {
  height: number;
  timestamp: Date;
}
/** AllDowntimeInfo stores information for all downtime durations. */
export interface AllDowntimeInfo {
  /**
   * The downtime information for each tracked duration. Sorted by duration,
   * ascending. (i.e. the same order as they appear in DowntimeParams).
   */
  infos: AllDowntimeInfo_DowntimeInfo[];
}
export interface AllDowntimeInfoProtoMsg {
  typeUrl: "/dydxprotocol.blocktime.AllDowntimeInfo";
  value: Uint8Array;
}
/** AllDowntimeInfo stores information for all downtime durations. */
export interface AllDowntimeInfoAmino {
  /**
   * The downtime information for each tracked duration. Sorted by duration,
   * ascending. (i.e. the same order as they appear in DowntimeParams).
   */
  infos?: AllDowntimeInfo_DowntimeInfoAmino[];
}
export interface AllDowntimeInfoAminoMsg {
  type: "/dydxprotocol.blocktime.AllDowntimeInfo";
  value: AllDowntimeInfoAmino;
}
/** AllDowntimeInfo stores information for all downtime durations. */
export interface AllDowntimeInfoSDKType {
  infos: AllDowntimeInfo_DowntimeInfoSDKType[];
}
/**
 * Stores information about downtime. block_info corresponds to the most
 * recent block at which a downtime occurred.
 */
export interface AllDowntimeInfo_DowntimeInfo {
  duration: Duration;
  blockInfo: BlockInfo;
}
export interface AllDowntimeInfo_DowntimeInfoProtoMsg {
  typeUrl: "/dydxprotocol.blocktime.DowntimeInfo";
  value: Uint8Array;
}
/**
 * Stores information about downtime. block_info corresponds to the most
 * recent block at which a downtime occurred.
 */
export interface AllDowntimeInfo_DowntimeInfoAmino {
  duration?: DurationAmino;
  block_info?: BlockInfoAmino;
}
export interface AllDowntimeInfo_DowntimeInfoAminoMsg {
  type: "/dydxprotocol.blocktime.DowntimeInfo";
  value: AllDowntimeInfo_DowntimeInfoAmino;
}
/**
 * Stores information about downtime. block_info corresponds to the most
 * recent block at which a downtime occurred.
 */
export interface AllDowntimeInfo_DowntimeInfoSDKType {
  duration: DurationSDKType;
  block_info: BlockInfoSDKType;
}
function createBaseBlockInfo(): BlockInfo {
  return {
    height: 0,
    timestamp: new Date()
  };
}
export const BlockInfo = {
  typeUrl: "/dydxprotocol.blocktime.BlockInfo",
  encode(message: BlockInfo, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.height !== 0) {
      writer.uint32(8).uint32(message.height);
    }
    if (message.timestamp !== undefined) {
      Timestamp.encode(toTimestamp(message.timestamp), writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): BlockInfo {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
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
  fromPartial(object: Partial<BlockInfo>): BlockInfo {
    const message = createBaseBlockInfo();
    message.height = object.height ?? 0;
    message.timestamp = object.timestamp ?? undefined;
    return message;
  },
  fromAmino(object: BlockInfoAmino): BlockInfo {
    const message = createBaseBlockInfo();
    if (object.height !== undefined && object.height !== null) {
      message.height = object.height;
    }
    if (object.timestamp !== undefined && object.timestamp !== null) {
      message.timestamp = fromTimestamp(Timestamp.fromAmino(object.timestamp));
    }
    return message;
  },
  toAmino(message: BlockInfo): BlockInfoAmino {
    const obj: any = {};
    obj.height = message.height;
    obj.timestamp = message.timestamp ? Timestamp.toAmino(toTimestamp(message.timestamp)) : undefined;
    return obj;
  },
  fromAminoMsg(object: BlockInfoAminoMsg): BlockInfo {
    return BlockInfo.fromAmino(object.value);
  },
  fromProtoMsg(message: BlockInfoProtoMsg): BlockInfo {
    return BlockInfo.decode(message.value);
  },
  toProto(message: BlockInfo): Uint8Array {
    return BlockInfo.encode(message).finish();
  },
  toProtoMsg(message: BlockInfo): BlockInfoProtoMsg {
    return {
      typeUrl: "/dydxprotocol.blocktime.BlockInfo",
      value: BlockInfo.encode(message).finish()
    };
  }
};
function createBaseAllDowntimeInfo(): AllDowntimeInfo {
  return {
    infos: []
  };
}
export const AllDowntimeInfo = {
  typeUrl: "/dydxprotocol.blocktime.AllDowntimeInfo",
  encode(message: AllDowntimeInfo, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    for (const v of message.infos) {
      AllDowntimeInfo_DowntimeInfo.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): AllDowntimeInfo {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
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
  fromPartial(object: Partial<AllDowntimeInfo>): AllDowntimeInfo {
    const message = createBaseAllDowntimeInfo();
    message.infos = object.infos?.map(e => AllDowntimeInfo_DowntimeInfo.fromPartial(e)) || [];
    return message;
  },
  fromAmino(object: AllDowntimeInfoAmino): AllDowntimeInfo {
    const message = createBaseAllDowntimeInfo();
    message.infos = object.infos?.map(e => AllDowntimeInfo_DowntimeInfo.fromAmino(e)) || [];
    return message;
  },
  toAmino(message: AllDowntimeInfo): AllDowntimeInfoAmino {
    const obj: any = {};
    if (message.infos) {
      obj.infos = message.infos.map(e => e ? AllDowntimeInfo_DowntimeInfo.toAmino(e) : undefined);
    } else {
      obj.infos = [];
    }
    return obj;
  },
  fromAminoMsg(object: AllDowntimeInfoAminoMsg): AllDowntimeInfo {
    return AllDowntimeInfo.fromAmino(object.value);
  },
  fromProtoMsg(message: AllDowntimeInfoProtoMsg): AllDowntimeInfo {
    return AllDowntimeInfo.decode(message.value);
  },
  toProto(message: AllDowntimeInfo): Uint8Array {
    return AllDowntimeInfo.encode(message).finish();
  },
  toProtoMsg(message: AllDowntimeInfo): AllDowntimeInfoProtoMsg {
    return {
      typeUrl: "/dydxprotocol.blocktime.AllDowntimeInfo",
      value: AllDowntimeInfo.encode(message).finish()
    };
  }
};
function createBaseAllDowntimeInfo_DowntimeInfo(): AllDowntimeInfo_DowntimeInfo {
  return {
    duration: Duration.fromPartial({}),
    blockInfo: BlockInfo.fromPartial({})
  };
}
export const AllDowntimeInfo_DowntimeInfo = {
  typeUrl: "/dydxprotocol.blocktime.DowntimeInfo",
  encode(message: AllDowntimeInfo_DowntimeInfo, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.duration !== undefined) {
      Duration.encode(message.duration, writer.uint32(10).fork()).ldelim();
    }
    if (message.blockInfo !== undefined) {
      BlockInfo.encode(message.blockInfo, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): AllDowntimeInfo_DowntimeInfo {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
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
  fromPartial(object: Partial<AllDowntimeInfo_DowntimeInfo>): AllDowntimeInfo_DowntimeInfo {
    const message = createBaseAllDowntimeInfo_DowntimeInfo();
    message.duration = object.duration !== undefined && object.duration !== null ? Duration.fromPartial(object.duration) : undefined;
    message.blockInfo = object.blockInfo !== undefined && object.blockInfo !== null ? BlockInfo.fromPartial(object.blockInfo) : undefined;
    return message;
  },
  fromAmino(object: AllDowntimeInfo_DowntimeInfoAmino): AllDowntimeInfo_DowntimeInfo {
    const message = createBaseAllDowntimeInfo_DowntimeInfo();
    if (object.duration !== undefined && object.duration !== null) {
      message.duration = Duration.fromAmino(object.duration);
    }
    if (object.block_info !== undefined && object.block_info !== null) {
      message.blockInfo = BlockInfo.fromAmino(object.block_info);
    }
    return message;
  },
  toAmino(message: AllDowntimeInfo_DowntimeInfo): AllDowntimeInfo_DowntimeInfoAmino {
    const obj: any = {};
    obj.duration = message.duration ? Duration.toAmino(message.duration) : undefined;
    obj.block_info = message.blockInfo ? BlockInfo.toAmino(message.blockInfo) : undefined;
    return obj;
  },
  fromAminoMsg(object: AllDowntimeInfo_DowntimeInfoAminoMsg): AllDowntimeInfo_DowntimeInfo {
    return AllDowntimeInfo_DowntimeInfo.fromAmino(object.value);
  },
  fromProtoMsg(message: AllDowntimeInfo_DowntimeInfoProtoMsg): AllDowntimeInfo_DowntimeInfo {
    return AllDowntimeInfo_DowntimeInfo.decode(message.value);
  },
  toProto(message: AllDowntimeInfo_DowntimeInfo): Uint8Array {
    return AllDowntimeInfo_DowntimeInfo.encode(message).finish();
  },
  toProtoMsg(message: AllDowntimeInfo_DowntimeInfo): AllDowntimeInfo_DowntimeInfoProtoMsg {
    return {
      typeUrl: "/dydxprotocol.blocktime.DowntimeInfo",
      value: AllDowntimeInfo_DowntimeInfo.encode(message).finish()
    };
  }
};