import { BinaryReader, BinaryWriter } from "../../binary";
/** Status of the CLOB. */
export enum ClobPair_Status {
  /** STATUS_UNSPECIFIED - Default value. This value is invalid and unused. */
  STATUS_UNSPECIFIED = 0,
  /** STATUS_ACTIVE - STATUS_ACTIVE represents an active clob pair. */
  STATUS_ACTIVE = 1,
  /**
   * STATUS_PAUSED - STATUS_PAUSED behavior is unfinalized.
   * TODO(DEC-600): update this documentation.
   */
  STATUS_PAUSED = 2,
  /**
   * STATUS_CANCEL_ONLY - STATUS_CANCEL_ONLY behavior is unfinalized.
   * TODO(DEC-600): update this documentation.
   */
  STATUS_CANCEL_ONLY = 3,
  /**
   * STATUS_POST_ONLY - STATUS_POST_ONLY behavior is unfinalized.
   * TODO(DEC-600): update this documentation.
   */
  STATUS_POST_ONLY = 4,
  /**
   * STATUS_INITIALIZING - STATUS_INITIALIZING represents a newly-added clob pair.
   * Clob pairs in this state only accept orders which are
   * both short-term and post-only.
   */
  STATUS_INITIALIZING = 5,
  /**
   * STATUS_FINAL_SETTLEMENT - STATUS_FINAL_SETTLEMENT represents a clob pair which is deactivated
   * and trading has ceased. All open positions will be closed by the
   * protocol. Open stateful orders will be cancelled. Open short-term
   * orders will be left to expire.
   */
  STATUS_FINAL_SETTLEMENT = 6,
  UNRECOGNIZED = -1,
}
export const ClobPair_StatusSDKType = ClobPair_Status;
export const ClobPair_StatusAmino = ClobPair_Status;
export function clobPair_StatusFromJSON(object: any): ClobPair_Status {
  switch (object) {
    case 0:
    case "STATUS_UNSPECIFIED":
      return ClobPair_Status.STATUS_UNSPECIFIED;
    case 1:
    case "STATUS_ACTIVE":
      return ClobPair_Status.STATUS_ACTIVE;
    case 2:
    case "STATUS_PAUSED":
      return ClobPair_Status.STATUS_PAUSED;
    case 3:
    case "STATUS_CANCEL_ONLY":
      return ClobPair_Status.STATUS_CANCEL_ONLY;
    case 4:
    case "STATUS_POST_ONLY":
      return ClobPair_Status.STATUS_POST_ONLY;
    case 5:
    case "STATUS_INITIALIZING":
      return ClobPair_Status.STATUS_INITIALIZING;
    case 6:
    case "STATUS_FINAL_SETTLEMENT":
      return ClobPair_Status.STATUS_FINAL_SETTLEMENT;
    case -1:
    case "UNRECOGNIZED":
    default:
      return ClobPair_Status.UNRECOGNIZED;
  }
}
export function clobPair_StatusToJSON(object: ClobPair_Status): string {
  switch (object) {
    case ClobPair_Status.STATUS_UNSPECIFIED:
      return "STATUS_UNSPECIFIED";
    case ClobPair_Status.STATUS_ACTIVE:
      return "STATUS_ACTIVE";
    case ClobPair_Status.STATUS_PAUSED:
      return "STATUS_PAUSED";
    case ClobPair_Status.STATUS_CANCEL_ONLY:
      return "STATUS_CANCEL_ONLY";
    case ClobPair_Status.STATUS_POST_ONLY:
      return "STATUS_POST_ONLY";
    case ClobPair_Status.STATUS_INITIALIZING:
      return "STATUS_INITIALIZING";
    case ClobPair_Status.STATUS_FINAL_SETTLEMENT:
      return "STATUS_FINAL_SETTLEMENT";
    case ClobPair_Status.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}
/**
 * PerpetualClobMetadata contains metadata for a `ClobPair`
 * representing a Perpetual product.
 */
export interface PerpetualClobMetadata {
  /** Id of the Perpetual the CLOB allows trading of. */
  perpetualId: number;
}
export interface PerpetualClobMetadataProtoMsg {
  typeUrl: "/dydxprotocol.clob.PerpetualClobMetadata";
  value: Uint8Array;
}
/**
 * PerpetualClobMetadata contains metadata for a `ClobPair`
 * representing a Perpetual product.
 */
export interface PerpetualClobMetadataAmino {
  /** Id of the Perpetual the CLOB allows trading of. */
  perpetual_id?: number;
}
export interface PerpetualClobMetadataAminoMsg {
  type: "/dydxprotocol.clob.PerpetualClobMetadata";
  value: PerpetualClobMetadataAmino;
}
/**
 * PerpetualClobMetadata contains metadata for a `ClobPair`
 * representing a Perpetual product.
 */
export interface PerpetualClobMetadataSDKType {
  perpetual_id: number;
}
/**
 * PerpetualClobMetadata contains metadata for a `ClobPair`
 * representing a Spot product.
 */
export interface SpotClobMetadata {
  /** Id of the base Asset in the trading pair. */
  baseAssetId: number;
  /** Id of the quote Asset in the trading pair. */
  quoteAssetId: number;
}
export interface SpotClobMetadataProtoMsg {
  typeUrl: "/dydxprotocol.clob.SpotClobMetadata";
  value: Uint8Array;
}
/**
 * PerpetualClobMetadata contains metadata for a `ClobPair`
 * representing a Spot product.
 */
export interface SpotClobMetadataAmino {
  /** Id of the base Asset in the trading pair. */
  base_asset_id?: number;
  /** Id of the quote Asset in the trading pair. */
  quote_asset_id?: number;
}
export interface SpotClobMetadataAminoMsg {
  type: "/dydxprotocol.clob.SpotClobMetadata";
  value: SpotClobMetadataAmino;
}
/**
 * PerpetualClobMetadata contains metadata for a `ClobPair`
 * representing a Spot product.
 */
export interface SpotClobMetadataSDKType {
  base_asset_id: number;
  quote_asset_id: number;
}
/**
 * ClobPair represents a single CLOB pair for a given product
 * in state.
 */
export interface ClobPair {
  /** ID of the orderbook that stores all resting liquidity for this CLOB. */
  id: number;
  perpetualClobMetadata?: PerpetualClobMetadata;
  spotClobMetadata?: SpotClobMetadata;
  /** Minimum increment in the size of orders on the CLOB, in base quantums. */
  stepBaseQuantums: bigint;
  /**
   * Defines the tick size of the orderbook by defining how many subticks
   * are in one tick. That is, the subticks of any valid order must be a
   * multiple of this value. Generally this value should start `>= 100`to
   * allow room for decreasing it.
   */
  subticksPerTick: number;
  /**
   * `10^Exponent` gives the number of QuoteQuantums traded per BaseQuantum
   * per Subtick.
   */
  quantumConversionExponent: number;
  status: ClobPair_Status;
}
export interface ClobPairProtoMsg {
  typeUrl: "/dydxprotocol.clob.ClobPair";
  value: Uint8Array;
}
/**
 * ClobPair represents a single CLOB pair for a given product
 * in state.
 */
export interface ClobPairAmino {
  /** ID of the orderbook that stores all resting liquidity for this CLOB. */
  id?: number;
  perpetual_clob_metadata?: PerpetualClobMetadataAmino;
  spot_clob_metadata?: SpotClobMetadataAmino;
  /** Minimum increment in the size of orders on the CLOB, in base quantums. */
  step_base_quantums?: string;
  /**
   * Defines the tick size of the orderbook by defining how many subticks
   * are in one tick. That is, the subticks of any valid order must be a
   * multiple of this value. Generally this value should start `>= 100`to
   * allow room for decreasing it.
   */
  subticks_per_tick?: number;
  /**
   * `10^Exponent` gives the number of QuoteQuantums traded per BaseQuantum
   * per Subtick.
   */
  quantum_conversion_exponent?: number;
  status?: ClobPair_Status;
}
export interface ClobPairAminoMsg {
  type: "/dydxprotocol.clob.ClobPair";
  value: ClobPairAmino;
}
/**
 * ClobPair represents a single CLOB pair for a given product
 * in state.
 */
export interface ClobPairSDKType {
  id: number;
  perpetual_clob_metadata?: PerpetualClobMetadataSDKType;
  spot_clob_metadata?: SpotClobMetadataSDKType;
  step_base_quantums: bigint;
  subticks_per_tick: number;
  quantum_conversion_exponent: number;
  status: ClobPair_Status;
}
function createBasePerpetualClobMetadata(): PerpetualClobMetadata {
  return {
    perpetualId: 0
  };
}
export const PerpetualClobMetadata = {
  typeUrl: "/dydxprotocol.clob.PerpetualClobMetadata",
  encode(message: PerpetualClobMetadata, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.perpetualId !== 0) {
      writer.uint32(8).uint32(message.perpetualId);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): PerpetualClobMetadata {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBasePerpetualClobMetadata();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.perpetualId = reader.uint32();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<PerpetualClobMetadata>): PerpetualClobMetadata {
    const message = createBasePerpetualClobMetadata();
    message.perpetualId = object.perpetualId ?? 0;
    return message;
  },
  fromAmino(object: PerpetualClobMetadataAmino): PerpetualClobMetadata {
    const message = createBasePerpetualClobMetadata();
    if (object.perpetual_id !== undefined && object.perpetual_id !== null) {
      message.perpetualId = object.perpetual_id;
    }
    return message;
  },
  toAmino(message: PerpetualClobMetadata): PerpetualClobMetadataAmino {
    const obj: any = {};
    obj.perpetual_id = message.perpetualId;
    return obj;
  },
  fromAminoMsg(object: PerpetualClobMetadataAminoMsg): PerpetualClobMetadata {
    return PerpetualClobMetadata.fromAmino(object.value);
  },
  fromProtoMsg(message: PerpetualClobMetadataProtoMsg): PerpetualClobMetadata {
    return PerpetualClobMetadata.decode(message.value);
  },
  toProto(message: PerpetualClobMetadata): Uint8Array {
    return PerpetualClobMetadata.encode(message).finish();
  },
  toProtoMsg(message: PerpetualClobMetadata): PerpetualClobMetadataProtoMsg {
    return {
      typeUrl: "/dydxprotocol.clob.PerpetualClobMetadata",
      value: PerpetualClobMetadata.encode(message).finish()
    };
  }
};
function createBaseSpotClobMetadata(): SpotClobMetadata {
  return {
    baseAssetId: 0,
    quoteAssetId: 0
  };
}
export const SpotClobMetadata = {
  typeUrl: "/dydxprotocol.clob.SpotClobMetadata",
  encode(message: SpotClobMetadata, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.baseAssetId !== 0) {
      writer.uint32(8).uint32(message.baseAssetId);
    }
    if (message.quoteAssetId !== 0) {
      writer.uint32(16).uint32(message.quoteAssetId);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): SpotClobMetadata {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSpotClobMetadata();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.baseAssetId = reader.uint32();
          break;
        case 2:
          message.quoteAssetId = reader.uint32();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<SpotClobMetadata>): SpotClobMetadata {
    const message = createBaseSpotClobMetadata();
    message.baseAssetId = object.baseAssetId ?? 0;
    message.quoteAssetId = object.quoteAssetId ?? 0;
    return message;
  },
  fromAmino(object: SpotClobMetadataAmino): SpotClobMetadata {
    const message = createBaseSpotClobMetadata();
    if (object.base_asset_id !== undefined && object.base_asset_id !== null) {
      message.baseAssetId = object.base_asset_id;
    }
    if (object.quote_asset_id !== undefined && object.quote_asset_id !== null) {
      message.quoteAssetId = object.quote_asset_id;
    }
    return message;
  },
  toAmino(message: SpotClobMetadata): SpotClobMetadataAmino {
    const obj: any = {};
    obj.base_asset_id = message.baseAssetId;
    obj.quote_asset_id = message.quoteAssetId;
    return obj;
  },
  fromAminoMsg(object: SpotClobMetadataAminoMsg): SpotClobMetadata {
    return SpotClobMetadata.fromAmino(object.value);
  },
  fromProtoMsg(message: SpotClobMetadataProtoMsg): SpotClobMetadata {
    return SpotClobMetadata.decode(message.value);
  },
  toProto(message: SpotClobMetadata): Uint8Array {
    return SpotClobMetadata.encode(message).finish();
  },
  toProtoMsg(message: SpotClobMetadata): SpotClobMetadataProtoMsg {
    return {
      typeUrl: "/dydxprotocol.clob.SpotClobMetadata",
      value: SpotClobMetadata.encode(message).finish()
    };
  }
};
function createBaseClobPair(): ClobPair {
  return {
    id: 0,
    perpetualClobMetadata: undefined,
    spotClobMetadata: undefined,
    stepBaseQuantums: BigInt(0),
    subticksPerTick: 0,
    quantumConversionExponent: 0,
    status: 0
  };
}
export const ClobPair = {
  typeUrl: "/dydxprotocol.clob.ClobPair",
  encode(message: ClobPair, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.id !== 0) {
      writer.uint32(8).uint32(message.id);
    }
    if (message.perpetualClobMetadata !== undefined) {
      PerpetualClobMetadata.encode(message.perpetualClobMetadata, writer.uint32(18).fork()).ldelim();
    }
    if (message.spotClobMetadata !== undefined) {
      SpotClobMetadata.encode(message.spotClobMetadata, writer.uint32(26).fork()).ldelim();
    }
    if (message.stepBaseQuantums !== BigInt(0)) {
      writer.uint32(32).uint64(message.stepBaseQuantums);
    }
    if (message.subticksPerTick !== 0) {
      writer.uint32(40).uint32(message.subticksPerTick);
    }
    if (message.quantumConversionExponent !== 0) {
      writer.uint32(48).sint32(message.quantumConversionExponent);
    }
    if (message.status !== 0) {
      writer.uint32(56).int32(message.status);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): ClobPair {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseClobPair();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.id = reader.uint32();
          break;
        case 2:
          message.perpetualClobMetadata = PerpetualClobMetadata.decode(reader, reader.uint32());
          break;
        case 3:
          message.spotClobMetadata = SpotClobMetadata.decode(reader, reader.uint32());
          break;
        case 4:
          message.stepBaseQuantums = reader.uint64();
          break;
        case 5:
          message.subticksPerTick = reader.uint32();
          break;
        case 6:
          message.quantumConversionExponent = reader.sint32();
          break;
        case 7:
          message.status = (reader.int32() as any);
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<ClobPair>): ClobPair {
    const message = createBaseClobPair();
    message.id = object.id ?? 0;
    message.perpetualClobMetadata = object.perpetualClobMetadata !== undefined && object.perpetualClobMetadata !== null ? PerpetualClobMetadata.fromPartial(object.perpetualClobMetadata) : undefined;
    message.spotClobMetadata = object.spotClobMetadata !== undefined && object.spotClobMetadata !== null ? SpotClobMetadata.fromPartial(object.spotClobMetadata) : undefined;
    message.stepBaseQuantums = object.stepBaseQuantums !== undefined && object.stepBaseQuantums !== null ? BigInt(object.stepBaseQuantums.toString()) : BigInt(0);
    message.subticksPerTick = object.subticksPerTick ?? 0;
    message.quantumConversionExponent = object.quantumConversionExponent ?? 0;
    message.status = object.status ?? 0;
    return message;
  },
  fromAmino(object: ClobPairAmino): ClobPair {
    const message = createBaseClobPair();
    if (object.id !== undefined && object.id !== null) {
      message.id = object.id;
    }
    if (object.perpetual_clob_metadata !== undefined && object.perpetual_clob_metadata !== null) {
      message.perpetualClobMetadata = PerpetualClobMetadata.fromAmino(object.perpetual_clob_metadata);
    }
    if (object.spot_clob_metadata !== undefined && object.spot_clob_metadata !== null) {
      message.spotClobMetadata = SpotClobMetadata.fromAmino(object.spot_clob_metadata);
    }
    if (object.step_base_quantums !== undefined && object.step_base_quantums !== null) {
      message.stepBaseQuantums = BigInt(object.step_base_quantums);
    }
    if (object.subticks_per_tick !== undefined && object.subticks_per_tick !== null) {
      message.subticksPerTick = object.subticks_per_tick;
    }
    if (object.quantum_conversion_exponent !== undefined && object.quantum_conversion_exponent !== null) {
      message.quantumConversionExponent = object.quantum_conversion_exponent;
    }
    if (object.status !== undefined && object.status !== null) {
      message.status = clobPair_StatusFromJSON(object.status);
    }
    return message;
  },
  toAmino(message: ClobPair): ClobPairAmino {
    const obj: any = {};
    obj.id = message.id;
    obj.perpetual_clob_metadata = message.perpetualClobMetadata ? PerpetualClobMetadata.toAmino(message.perpetualClobMetadata) : undefined;
    obj.spot_clob_metadata = message.spotClobMetadata ? SpotClobMetadata.toAmino(message.spotClobMetadata) : undefined;
    obj.step_base_quantums = message.stepBaseQuantums ? message.stepBaseQuantums.toString() : undefined;
    obj.subticks_per_tick = message.subticksPerTick;
    obj.quantum_conversion_exponent = message.quantumConversionExponent;
    obj.status = clobPair_StatusToJSON(message.status);
    return obj;
  },
  fromAminoMsg(object: ClobPairAminoMsg): ClobPair {
    return ClobPair.fromAmino(object.value);
  },
  fromProtoMsg(message: ClobPairProtoMsg): ClobPair {
    return ClobPair.decode(message.value);
  },
  toProto(message: ClobPair): Uint8Array {
    return ClobPair.encode(message).finish();
  },
  toProtoMsg(message: ClobPair): ClobPairProtoMsg {
    return {
      typeUrl: "/dydxprotocol.clob.ClobPair",
      value: ClobPair.encode(message).finish()
    };
  }
};