import { Timestamp } from "../../google/protobuf/timestamp";
import { BinaryReader, BinaryWriter } from "../../binary";
import { toTimestamp, fromTimestamp } from "../../helpers";
/** BlockStats is used to store stats transiently within the scope of a block. */
export interface BlockStats {
  /** The fills that occured on this block. */
  fills: BlockStats_Fill[];
}
export interface BlockStatsProtoMsg {
  typeUrl: "/dydxprotocol.stats.BlockStats";
  value: Uint8Array;
}
/** BlockStats is used to store stats transiently within the scope of a block. */
export interface BlockStatsAmino {
  /** The fills that occured on this block. */
  fills?: BlockStats_FillAmino[];
}
export interface BlockStatsAminoMsg {
  type: "/dydxprotocol.stats.BlockStats";
  value: BlockStatsAmino;
}
/** BlockStats is used to store stats transiently within the scope of a block. */
export interface BlockStatsSDKType {
  fills: BlockStats_FillSDKType[];
}
/** Fill records data about a fill on this block. */
export interface BlockStats_Fill {
  /** Taker wallet address */
  taker: string;
  /** Maker wallet address */
  maker: string;
  /** Notional USDC filled in quantums */
  notional: bigint;
}
export interface BlockStats_FillProtoMsg {
  typeUrl: "/dydxprotocol.stats.Fill";
  value: Uint8Array;
}
/** Fill records data about a fill on this block. */
export interface BlockStats_FillAmino {
  /** Taker wallet address */
  taker?: string;
  /** Maker wallet address */
  maker?: string;
  /** Notional USDC filled in quantums */
  notional?: string;
}
export interface BlockStats_FillAminoMsg {
  type: "/dydxprotocol.stats.Fill";
  value: BlockStats_FillAmino;
}
/** Fill records data about a fill on this block. */
export interface BlockStats_FillSDKType {
  taker: string;
  maker: string;
  notional: bigint;
}
/** StatsMetadata stores metadata for the x/stats module */
export interface StatsMetadata {
  /**
   * The oldest epoch that is included in the stats. The next epoch to be
   * removed from the window.
   */
  trailingEpoch: number;
}
export interface StatsMetadataProtoMsg {
  typeUrl: "/dydxprotocol.stats.StatsMetadata";
  value: Uint8Array;
}
/** StatsMetadata stores metadata for the x/stats module */
export interface StatsMetadataAmino {
  /**
   * The oldest epoch that is included in the stats. The next epoch to be
   * removed from the window.
   */
  trailing_epoch?: number;
}
export interface StatsMetadataAminoMsg {
  type: "/dydxprotocol.stats.StatsMetadata";
  value: StatsMetadataAmino;
}
/** StatsMetadata stores metadata for the x/stats module */
export interface StatsMetadataSDKType {
  trailing_epoch: number;
}
/** EpochStats stores stats for a particular epoch */
export interface EpochStats {
  /** Epoch end time */
  epochEndTime: Date;
  /** Stats for each user in this epoch. Sorted by user. */
  stats: EpochStats_UserWithStats[];
}
export interface EpochStatsProtoMsg {
  typeUrl: "/dydxprotocol.stats.EpochStats";
  value: Uint8Array;
}
/** EpochStats stores stats for a particular epoch */
export interface EpochStatsAmino {
  /** Epoch end time */
  epoch_end_time?: string;
  /** Stats for each user in this epoch. Sorted by user. */
  stats?: EpochStats_UserWithStatsAmino[];
}
export interface EpochStatsAminoMsg {
  type: "/dydxprotocol.stats.EpochStats";
  value: EpochStatsAmino;
}
/** EpochStats stores stats for a particular epoch */
export interface EpochStatsSDKType {
  epoch_end_time: Date;
  stats: EpochStats_UserWithStatsSDKType[];
}
/** A user and its associated stats */
export interface EpochStats_UserWithStats {
  user: string;
  stats?: UserStats;
}
export interface EpochStats_UserWithStatsProtoMsg {
  typeUrl: "/dydxprotocol.stats.UserWithStats";
  value: Uint8Array;
}
/** A user and its associated stats */
export interface EpochStats_UserWithStatsAmino {
  user?: string;
  stats?: UserStatsAmino;
}
export interface EpochStats_UserWithStatsAminoMsg {
  type: "/dydxprotocol.stats.UserWithStats";
  value: EpochStats_UserWithStatsAmino;
}
/** A user and its associated stats */
export interface EpochStats_UserWithStatsSDKType {
  user: string;
  stats?: UserStatsSDKType;
}
/** GlobalStats stores global stats */
export interface GlobalStats {
  /** Notional USDC traded in quantums */
  notionalTraded: bigint;
}
export interface GlobalStatsProtoMsg {
  typeUrl: "/dydxprotocol.stats.GlobalStats";
  value: Uint8Array;
}
/** GlobalStats stores global stats */
export interface GlobalStatsAmino {
  /** Notional USDC traded in quantums */
  notional_traded?: string;
}
export interface GlobalStatsAminoMsg {
  type: "/dydxprotocol.stats.GlobalStats";
  value: GlobalStatsAmino;
}
/** GlobalStats stores global stats */
export interface GlobalStatsSDKType {
  notional_traded: bigint;
}
/** UserStats stores stats for a User */
export interface UserStats {
  /** Taker USDC in quantums */
  takerNotional: bigint;
  /** Maker USDC in quantums */
  makerNotional: bigint;
}
export interface UserStatsProtoMsg {
  typeUrl: "/dydxprotocol.stats.UserStats";
  value: Uint8Array;
}
/** UserStats stores stats for a User */
export interface UserStatsAmino {
  /** Taker USDC in quantums */
  taker_notional?: string;
  /** Maker USDC in quantums */
  maker_notional?: string;
}
export interface UserStatsAminoMsg {
  type: "/dydxprotocol.stats.UserStats";
  value: UserStatsAmino;
}
/** UserStats stores stats for a User */
export interface UserStatsSDKType {
  taker_notional: bigint;
  maker_notional: bigint;
}
function createBaseBlockStats(): BlockStats {
  return {
    fills: []
  };
}
export const BlockStats = {
  typeUrl: "/dydxprotocol.stats.BlockStats",
  encode(message: BlockStats, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    for (const v of message.fills) {
      BlockStats_Fill.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): BlockStats {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseBlockStats();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.fills.push(BlockStats_Fill.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<BlockStats>): BlockStats {
    const message = createBaseBlockStats();
    message.fills = object.fills?.map(e => BlockStats_Fill.fromPartial(e)) || [];
    return message;
  },
  fromAmino(object: BlockStatsAmino): BlockStats {
    const message = createBaseBlockStats();
    message.fills = object.fills?.map(e => BlockStats_Fill.fromAmino(e)) || [];
    return message;
  },
  toAmino(message: BlockStats): BlockStatsAmino {
    const obj: any = {};
    if (message.fills) {
      obj.fills = message.fills.map(e => e ? BlockStats_Fill.toAmino(e) : undefined);
    } else {
      obj.fills = [];
    }
    return obj;
  },
  fromAminoMsg(object: BlockStatsAminoMsg): BlockStats {
    return BlockStats.fromAmino(object.value);
  },
  fromProtoMsg(message: BlockStatsProtoMsg): BlockStats {
    return BlockStats.decode(message.value);
  },
  toProto(message: BlockStats): Uint8Array {
    return BlockStats.encode(message).finish();
  },
  toProtoMsg(message: BlockStats): BlockStatsProtoMsg {
    return {
      typeUrl: "/dydxprotocol.stats.BlockStats",
      value: BlockStats.encode(message).finish()
    };
  }
};
function createBaseBlockStats_Fill(): BlockStats_Fill {
  return {
    taker: "",
    maker: "",
    notional: BigInt(0)
  };
}
export const BlockStats_Fill = {
  typeUrl: "/dydxprotocol.stats.Fill",
  encode(message: BlockStats_Fill, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.taker !== "") {
      writer.uint32(10).string(message.taker);
    }
    if (message.maker !== "") {
      writer.uint32(18).string(message.maker);
    }
    if (message.notional !== BigInt(0)) {
      writer.uint32(24).uint64(message.notional);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): BlockStats_Fill {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseBlockStats_Fill();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.taker = reader.string();
          break;
        case 2:
          message.maker = reader.string();
          break;
        case 3:
          message.notional = reader.uint64();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<BlockStats_Fill>): BlockStats_Fill {
    const message = createBaseBlockStats_Fill();
    message.taker = object.taker ?? "";
    message.maker = object.maker ?? "";
    message.notional = object.notional !== undefined && object.notional !== null ? BigInt(object.notional.toString()) : BigInt(0);
    return message;
  },
  fromAmino(object: BlockStats_FillAmino): BlockStats_Fill {
    const message = createBaseBlockStats_Fill();
    if (object.taker !== undefined && object.taker !== null) {
      message.taker = object.taker;
    }
    if (object.maker !== undefined && object.maker !== null) {
      message.maker = object.maker;
    }
    if (object.notional !== undefined && object.notional !== null) {
      message.notional = BigInt(object.notional);
    }
    return message;
  },
  toAmino(message: BlockStats_Fill): BlockStats_FillAmino {
    const obj: any = {};
    obj.taker = message.taker;
    obj.maker = message.maker;
    obj.notional = message.notional ? message.notional.toString() : undefined;
    return obj;
  },
  fromAminoMsg(object: BlockStats_FillAminoMsg): BlockStats_Fill {
    return BlockStats_Fill.fromAmino(object.value);
  },
  fromProtoMsg(message: BlockStats_FillProtoMsg): BlockStats_Fill {
    return BlockStats_Fill.decode(message.value);
  },
  toProto(message: BlockStats_Fill): Uint8Array {
    return BlockStats_Fill.encode(message).finish();
  },
  toProtoMsg(message: BlockStats_Fill): BlockStats_FillProtoMsg {
    return {
      typeUrl: "/dydxprotocol.stats.Fill",
      value: BlockStats_Fill.encode(message).finish()
    };
  }
};
function createBaseStatsMetadata(): StatsMetadata {
  return {
    trailingEpoch: 0
  };
}
export const StatsMetadata = {
  typeUrl: "/dydxprotocol.stats.StatsMetadata",
  encode(message: StatsMetadata, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.trailingEpoch !== 0) {
      writer.uint32(8).uint32(message.trailingEpoch);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): StatsMetadata {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseStatsMetadata();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.trailingEpoch = reader.uint32();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<StatsMetadata>): StatsMetadata {
    const message = createBaseStatsMetadata();
    message.trailingEpoch = object.trailingEpoch ?? 0;
    return message;
  },
  fromAmino(object: StatsMetadataAmino): StatsMetadata {
    const message = createBaseStatsMetadata();
    if (object.trailing_epoch !== undefined && object.trailing_epoch !== null) {
      message.trailingEpoch = object.trailing_epoch;
    }
    return message;
  },
  toAmino(message: StatsMetadata): StatsMetadataAmino {
    const obj: any = {};
    obj.trailing_epoch = message.trailingEpoch;
    return obj;
  },
  fromAminoMsg(object: StatsMetadataAminoMsg): StatsMetadata {
    return StatsMetadata.fromAmino(object.value);
  },
  fromProtoMsg(message: StatsMetadataProtoMsg): StatsMetadata {
    return StatsMetadata.decode(message.value);
  },
  toProto(message: StatsMetadata): Uint8Array {
    return StatsMetadata.encode(message).finish();
  },
  toProtoMsg(message: StatsMetadata): StatsMetadataProtoMsg {
    return {
      typeUrl: "/dydxprotocol.stats.StatsMetadata",
      value: StatsMetadata.encode(message).finish()
    };
  }
};
function createBaseEpochStats(): EpochStats {
  return {
    epochEndTime: new Date(),
    stats: []
  };
}
export const EpochStats = {
  typeUrl: "/dydxprotocol.stats.EpochStats",
  encode(message: EpochStats, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.epochEndTime !== undefined) {
      Timestamp.encode(toTimestamp(message.epochEndTime), writer.uint32(10).fork()).ldelim();
    }
    for (const v of message.stats) {
      EpochStats_UserWithStats.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): EpochStats {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseEpochStats();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.epochEndTime = fromTimestamp(Timestamp.decode(reader, reader.uint32()));
          break;
        case 2:
          message.stats.push(EpochStats_UserWithStats.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<EpochStats>): EpochStats {
    const message = createBaseEpochStats();
    message.epochEndTime = object.epochEndTime ?? undefined;
    message.stats = object.stats?.map(e => EpochStats_UserWithStats.fromPartial(e)) || [];
    return message;
  },
  fromAmino(object: EpochStatsAmino): EpochStats {
    const message = createBaseEpochStats();
    if (object.epoch_end_time !== undefined && object.epoch_end_time !== null) {
      message.epochEndTime = fromTimestamp(Timestamp.fromAmino(object.epoch_end_time));
    }
    message.stats = object.stats?.map(e => EpochStats_UserWithStats.fromAmino(e)) || [];
    return message;
  },
  toAmino(message: EpochStats): EpochStatsAmino {
    const obj: any = {};
    obj.epoch_end_time = message.epochEndTime ? Timestamp.toAmino(toTimestamp(message.epochEndTime)) : undefined;
    if (message.stats) {
      obj.stats = message.stats.map(e => e ? EpochStats_UserWithStats.toAmino(e) : undefined);
    } else {
      obj.stats = [];
    }
    return obj;
  },
  fromAminoMsg(object: EpochStatsAminoMsg): EpochStats {
    return EpochStats.fromAmino(object.value);
  },
  fromProtoMsg(message: EpochStatsProtoMsg): EpochStats {
    return EpochStats.decode(message.value);
  },
  toProto(message: EpochStats): Uint8Array {
    return EpochStats.encode(message).finish();
  },
  toProtoMsg(message: EpochStats): EpochStatsProtoMsg {
    return {
      typeUrl: "/dydxprotocol.stats.EpochStats",
      value: EpochStats.encode(message).finish()
    };
  }
};
function createBaseEpochStats_UserWithStats(): EpochStats_UserWithStats {
  return {
    user: "",
    stats: undefined
  };
}
export const EpochStats_UserWithStats = {
  typeUrl: "/dydxprotocol.stats.UserWithStats",
  encode(message: EpochStats_UserWithStats, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.user !== "") {
      writer.uint32(10).string(message.user);
    }
    if (message.stats !== undefined) {
      UserStats.encode(message.stats, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): EpochStats_UserWithStats {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseEpochStats_UserWithStats();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.user = reader.string();
          break;
        case 2:
          message.stats = UserStats.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<EpochStats_UserWithStats>): EpochStats_UserWithStats {
    const message = createBaseEpochStats_UserWithStats();
    message.user = object.user ?? "";
    message.stats = object.stats !== undefined && object.stats !== null ? UserStats.fromPartial(object.stats) : undefined;
    return message;
  },
  fromAmino(object: EpochStats_UserWithStatsAmino): EpochStats_UserWithStats {
    const message = createBaseEpochStats_UserWithStats();
    if (object.user !== undefined && object.user !== null) {
      message.user = object.user;
    }
    if (object.stats !== undefined && object.stats !== null) {
      message.stats = UserStats.fromAmino(object.stats);
    }
    return message;
  },
  toAmino(message: EpochStats_UserWithStats): EpochStats_UserWithStatsAmino {
    const obj: any = {};
    obj.user = message.user;
    obj.stats = message.stats ? UserStats.toAmino(message.stats) : undefined;
    return obj;
  },
  fromAminoMsg(object: EpochStats_UserWithStatsAminoMsg): EpochStats_UserWithStats {
    return EpochStats_UserWithStats.fromAmino(object.value);
  },
  fromProtoMsg(message: EpochStats_UserWithStatsProtoMsg): EpochStats_UserWithStats {
    return EpochStats_UserWithStats.decode(message.value);
  },
  toProto(message: EpochStats_UserWithStats): Uint8Array {
    return EpochStats_UserWithStats.encode(message).finish();
  },
  toProtoMsg(message: EpochStats_UserWithStats): EpochStats_UserWithStatsProtoMsg {
    return {
      typeUrl: "/dydxprotocol.stats.UserWithStats",
      value: EpochStats_UserWithStats.encode(message).finish()
    };
  }
};
function createBaseGlobalStats(): GlobalStats {
  return {
    notionalTraded: BigInt(0)
  };
}
export const GlobalStats = {
  typeUrl: "/dydxprotocol.stats.GlobalStats",
  encode(message: GlobalStats, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.notionalTraded !== BigInt(0)) {
      writer.uint32(8).uint64(message.notionalTraded);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): GlobalStats {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGlobalStats();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.notionalTraded = reader.uint64();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<GlobalStats>): GlobalStats {
    const message = createBaseGlobalStats();
    message.notionalTraded = object.notionalTraded !== undefined && object.notionalTraded !== null ? BigInt(object.notionalTraded.toString()) : BigInt(0);
    return message;
  },
  fromAmino(object: GlobalStatsAmino): GlobalStats {
    const message = createBaseGlobalStats();
    if (object.notional_traded !== undefined && object.notional_traded !== null) {
      message.notionalTraded = BigInt(object.notional_traded);
    }
    return message;
  },
  toAmino(message: GlobalStats): GlobalStatsAmino {
    const obj: any = {};
    obj.notional_traded = message.notionalTraded ? message.notionalTraded.toString() : undefined;
    return obj;
  },
  fromAminoMsg(object: GlobalStatsAminoMsg): GlobalStats {
    return GlobalStats.fromAmino(object.value);
  },
  fromProtoMsg(message: GlobalStatsProtoMsg): GlobalStats {
    return GlobalStats.decode(message.value);
  },
  toProto(message: GlobalStats): Uint8Array {
    return GlobalStats.encode(message).finish();
  },
  toProtoMsg(message: GlobalStats): GlobalStatsProtoMsg {
    return {
      typeUrl: "/dydxprotocol.stats.GlobalStats",
      value: GlobalStats.encode(message).finish()
    };
  }
};
function createBaseUserStats(): UserStats {
  return {
    takerNotional: BigInt(0),
    makerNotional: BigInt(0)
  };
}
export const UserStats = {
  typeUrl: "/dydxprotocol.stats.UserStats",
  encode(message: UserStats, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.takerNotional !== BigInt(0)) {
      writer.uint32(8).uint64(message.takerNotional);
    }
    if (message.makerNotional !== BigInt(0)) {
      writer.uint32(16).uint64(message.makerNotional);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): UserStats {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUserStats();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.takerNotional = reader.uint64();
          break;
        case 2:
          message.makerNotional = reader.uint64();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<UserStats>): UserStats {
    const message = createBaseUserStats();
    message.takerNotional = object.takerNotional !== undefined && object.takerNotional !== null ? BigInt(object.takerNotional.toString()) : BigInt(0);
    message.makerNotional = object.makerNotional !== undefined && object.makerNotional !== null ? BigInt(object.makerNotional.toString()) : BigInt(0);
    return message;
  },
  fromAmino(object: UserStatsAmino): UserStats {
    const message = createBaseUserStats();
    if (object.taker_notional !== undefined && object.taker_notional !== null) {
      message.takerNotional = BigInt(object.taker_notional);
    }
    if (object.maker_notional !== undefined && object.maker_notional !== null) {
      message.makerNotional = BigInt(object.maker_notional);
    }
    return message;
  },
  toAmino(message: UserStats): UserStatsAmino {
    const obj: any = {};
    obj.taker_notional = message.takerNotional ? message.takerNotional.toString() : undefined;
    obj.maker_notional = message.makerNotional ? message.makerNotional.toString() : undefined;
    return obj;
  },
  fromAminoMsg(object: UserStatsAminoMsg): UserStats {
    return UserStats.fromAmino(object.value);
  },
  fromProtoMsg(message: UserStatsProtoMsg): UserStats {
    return UserStats.decode(message.value);
  },
  toProto(message: UserStats): Uint8Array {
    return UserStats.encode(message).finish();
  },
  toProtoMsg(message: UserStats): UserStatsProtoMsg {
    return {
      typeUrl: "/dydxprotocol.stats.UserStats",
      value: UserStats.encode(message).finish()
    };
  }
};