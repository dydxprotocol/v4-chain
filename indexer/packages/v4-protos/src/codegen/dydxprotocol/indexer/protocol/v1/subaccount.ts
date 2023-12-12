import { BinaryReader, BinaryWriter } from "../../../../binary";
import { bytesFromBase64, base64FromBytes } from "../../../../helpers";
/** IndexerSubaccountId defines a unique identifier for a Subaccount. */
export interface IndexerSubaccountId {
  /** The address of the wallet that owns this subaccount. */
  owner: string;
  /**
   * < 128 Since 128 should be enough to start and it fits within
   * 1 Byte (1 Bit needed to indicate that the first byte is the last).
   */
  number: number;
}
export interface IndexerSubaccountIdProtoMsg {
  typeUrl: "/dydxprotocol.indexer.protocol.v1.IndexerSubaccountId";
  value: Uint8Array;
}
/** IndexerSubaccountId defines a unique identifier for a Subaccount. */
export interface IndexerSubaccountIdAmino {
  /** The address of the wallet that owns this subaccount. */
  owner?: string;
  /**
   * < 128 Since 128 should be enough to start and it fits within
   * 1 Byte (1 Bit needed to indicate that the first byte is the last).
   */
  number?: number;
}
export interface IndexerSubaccountIdAminoMsg {
  type: "/dydxprotocol.indexer.protocol.v1.IndexerSubaccountId";
  value: IndexerSubaccountIdAmino;
}
/** IndexerSubaccountId defines a unique identifier for a Subaccount. */
export interface IndexerSubaccountIdSDKType {
  owner: string;
  number: number;
}
/**
 * IndexerPerpetualPosition are an account’s positions of a `Perpetual`.
 * Therefore they hold any information needed to trade perpetuals.
 */
export interface IndexerPerpetualPosition {
  /** The `Id` of the `Perpetual`. */
  perpetualId: number;
  /** The size of the position in base quantums. */
  quantums: Uint8Array;
  /**
   * The funding_index of the `Perpetual` the last time this position was
   * settled.
   */
  fundingIndex: Uint8Array;
  /**
   * Amount of funding payment (in quote quantums).
   * Note: 1. this field is not cumulative.
   * 2. a positive value means funding payment was paid out and
   * a negative value means funding payment was received.
   */
  fundingPayment: Uint8Array;
}
export interface IndexerPerpetualPositionProtoMsg {
  typeUrl: "/dydxprotocol.indexer.protocol.v1.IndexerPerpetualPosition";
  value: Uint8Array;
}
/**
 * IndexerPerpetualPosition are an account’s positions of a `Perpetual`.
 * Therefore they hold any information needed to trade perpetuals.
 */
export interface IndexerPerpetualPositionAmino {
  /** The `Id` of the `Perpetual`. */
  perpetual_id?: number;
  /** The size of the position in base quantums. */
  quantums?: string;
  /**
   * The funding_index of the `Perpetual` the last time this position was
   * settled.
   */
  funding_index?: string;
  /**
   * Amount of funding payment (in quote quantums).
   * Note: 1. this field is not cumulative.
   * 2. a positive value means funding payment was paid out and
   * a negative value means funding payment was received.
   */
  funding_payment?: string;
}
export interface IndexerPerpetualPositionAminoMsg {
  type: "/dydxprotocol.indexer.protocol.v1.IndexerPerpetualPosition";
  value: IndexerPerpetualPositionAmino;
}
/**
 * IndexerPerpetualPosition are an account’s positions of a `Perpetual`.
 * Therefore they hold any information needed to trade perpetuals.
 */
export interface IndexerPerpetualPositionSDKType {
  perpetual_id: number;
  quantums: Uint8Array;
  funding_index: Uint8Array;
  funding_payment: Uint8Array;
}
/**
 * IndexerAssetPosition define an account’s positions of an `Asset`.
 * Therefore they hold any information needed to trade on Spot and Margin.
 */
export interface IndexerAssetPosition {
  /** The `Id` of the `Asset`. */
  assetId: number;
  /** The absolute size of the position in base quantums. */
  quantums: Uint8Array;
  /**
   * The `Index` (either `LongIndex` or `ShortIndex`) of the `Asset` the last
   * time this position was settled
   * TODO(DEC-582): pending margin trading being added.
   */
  index: bigint;
}
export interface IndexerAssetPositionProtoMsg {
  typeUrl: "/dydxprotocol.indexer.protocol.v1.IndexerAssetPosition";
  value: Uint8Array;
}
/**
 * IndexerAssetPosition define an account’s positions of an `Asset`.
 * Therefore they hold any information needed to trade on Spot and Margin.
 */
export interface IndexerAssetPositionAmino {
  /** The `Id` of the `Asset`. */
  asset_id?: number;
  /** The absolute size of the position in base quantums. */
  quantums?: string;
  /**
   * The `Index` (either `LongIndex` or `ShortIndex`) of the `Asset` the last
   * time this position was settled
   * TODO(DEC-582): pending margin trading being added.
   */
  index?: string;
}
export interface IndexerAssetPositionAminoMsg {
  type: "/dydxprotocol.indexer.protocol.v1.IndexerAssetPosition";
  value: IndexerAssetPositionAmino;
}
/**
 * IndexerAssetPosition define an account’s positions of an `Asset`.
 * Therefore they hold any information needed to trade on Spot and Margin.
 */
export interface IndexerAssetPositionSDKType {
  asset_id: number;
  quantums: Uint8Array;
  index: bigint;
}
function createBaseIndexerSubaccountId(): IndexerSubaccountId {
  return {
    owner: "",
    number: 0
  };
}
export const IndexerSubaccountId = {
  typeUrl: "/dydxprotocol.indexer.protocol.v1.IndexerSubaccountId",
  encode(message: IndexerSubaccountId, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.owner !== "") {
      writer.uint32(10).string(message.owner);
    }
    if (message.number !== 0) {
      writer.uint32(16).uint32(message.number);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): IndexerSubaccountId {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseIndexerSubaccountId();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.owner = reader.string();
          break;
        case 2:
          message.number = reader.uint32();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<IndexerSubaccountId>): IndexerSubaccountId {
    const message = createBaseIndexerSubaccountId();
    message.owner = object.owner ?? "";
    message.number = object.number ?? 0;
    return message;
  },
  fromAmino(object: IndexerSubaccountIdAmino): IndexerSubaccountId {
    const message = createBaseIndexerSubaccountId();
    if (object.owner !== undefined && object.owner !== null) {
      message.owner = object.owner;
    }
    if (object.number !== undefined && object.number !== null) {
      message.number = object.number;
    }
    return message;
  },
  toAmino(message: IndexerSubaccountId): IndexerSubaccountIdAmino {
    const obj: any = {};
    obj.owner = message.owner;
    obj.number = message.number;
    return obj;
  },
  fromAminoMsg(object: IndexerSubaccountIdAminoMsg): IndexerSubaccountId {
    return IndexerSubaccountId.fromAmino(object.value);
  },
  fromProtoMsg(message: IndexerSubaccountIdProtoMsg): IndexerSubaccountId {
    return IndexerSubaccountId.decode(message.value);
  },
  toProto(message: IndexerSubaccountId): Uint8Array {
    return IndexerSubaccountId.encode(message).finish();
  },
  toProtoMsg(message: IndexerSubaccountId): IndexerSubaccountIdProtoMsg {
    return {
      typeUrl: "/dydxprotocol.indexer.protocol.v1.IndexerSubaccountId",
      value: IndexerSubaccountId.encode(message).finish()
    };
  }
};
function createBaseIndexerPerpetualPosition(): IndexerPerpetualPosition {
  return {
    perpetualId: 0,
    quantums: new Uint8Array(),
    fundingIndex: new Uint8Array(),
    fundingPayment: new Uint8Array()
  };
}
export const IndexerPerpetualPosition = {
  typeUrl: "/dydxprotocol.indexer.protocol.v1.IndexerPerpetualPosition",
  encode(message: IndexerPerpetualPosition, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.perpetualId !== 0) {
      writer.uint32(8).uint32(message.perpetualId);
    }
    if (message.quantums.length !== 0) {
      writer.uint32(18).bytes(message.quantums);
    }
    if (message.fundingIndex.length !== 0) {
      writer.uint32(26).bytes(message.fundingIndex);
    }
    if (message.fundingPayment.length !== 0) {
      writer.uint32(34).bytes(message.fundingPayment);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): IndexerPerpetualPosition {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseIndexerPerpetualPosition();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.perpetualId = reader.uint32();
          break;
        case 2:
          message.quantums = reader.bytes();
          break;
        case 3:
          message.fundingIndex = reader.bytes();
          break;
        case 4:
          message.fundingPayment = reader.bytes();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<IndexerPerpetualPosition>): IndexerPerpetualPosition {
    const message = createBaseIndexerPerpetualPosition();
    message.perpetualId = object.perpetualId ?? 0;
    message.quantums = object.quantums ?? new Uint8Array();
    message.fundingIndex = object.fundingIndex ?? new Uint8Array();
    message.fundingPayment = object.fundingPayment ?? new Uint8Array();
    return message;
  },
  fromAmino(object: IndexerPerpetualPositionAmino): IndexerPerpetualPosition {
    const message = createBaseIndexerPerpetualPosition();
    if (object.perpetual_id !== undefined && object.perpetual_id !== null) {
      message.perpetualId = object.perpetual_id;
    }
    if (object.quantums !== undefined && object.quantums !== null) {
      message.quantums = bytesFromBase64(object.quantums);
    }
    if (object.funding_index !== undefined && object.funding_index !== null) {
      message.fundingIndex = bytesFromBase64(object.funding_index);
    }
    if (object.funding_payment !== undefined && object.funding_payment !== null) {
      message.fundingPayment = bytesFromBase64(object.funding_payment);
    }
    return message;
  },
  toAmino(message: IndexerPerpetualPosition): IndexerPerpetualPositionAmino {
    const obj: any = {};
    obj.perpetual_id = message.perpetualId;
    obj.quantums = message.quantums ? base64FromBytes(message.quantums) : undefined;
    obj.funding_index = message.fundingIndex ? base64FromBytes(message.fundingIndex) : undefined;
    obj.funding_payment = message.fundingPayment ? base64FromBytes(message.fundingPayment) : undefined;
    return obj;
  },
  fromAminoMsg(object: IndexerPerpetualPositionAminoMsg): IndexerPerpetualPosition {
    return IndexerPerpetualPosition.fromAmino(object.value);
  },
  fromProtoMsg(message: IndexerPerpetualPositionProtoMsg): IndexerPerpetualPosition {
    return IndexerPerpetualPosition.decode(message.value);
  },
  toProto(message: IndexerPerpetualPosition): Uint8Array {
    return IndexerPerpetualPosition.encode(message).finish();
  },
  toProtoMsg(message: IndexerPerpetualPosition): IndexerPerpetualPositionProtoMsg {
    return {
      typeUrl: "/dydxprotocol.indexer.protocol.v1.IndexerPerpetualPosition",
      value: IndexerPerpetualPosition.encode(message).finish()
    };
  }
};
function createBaseIndexerAssetPosition(): IndexerAssetPosition {
  return {
    assetId: 0,
    quantums: new Uint8Array(),
    index: BigInt(0)
  };
}
export const IndexerAssetPosition = {
  typeUrl: "/dydxprotocol.indexer.protocol.v1.IndexerAssetPosition",
  encode(message: IndexerAssetPosition, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.assetId !== 0) {
      writer.uint32(8).uint32(message.assetId);
    }
    if (message.quantums.length !== 0) {
      writer.uint32(18).bytes(message.quantums);
    }
    if (message.index !== BigInt(0)) {
      writer.uint32(24).uint64(message.index);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): IndexerAssetPosition {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseIndexerAssetPosition();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.assetId = reader.uint32();
          break;
        case 2:
          message.quantums = reader.bytes();
          break;
        case 3:
          message.index = reader.uint64();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<IndexerAssetPosition>): IndexerAssetPosition {
    const message = createBaseIndexerAssetPosition();
    message.assetId = object.assetId ?? 0;
    message.quantums = object.quantums ?? new Uint8Array();
    message.index = object.index !== undefined && object.index !== null ? BigInt(object.index.toString()) : BigInt(0);
    return message;
  },
  fromAmino(object: IndexerAssetPositionAmino): IndexerAssetPosition {
    const message = createBaseIndexerAssetPosition();
    if (object.asset_id !== undefined && object.asset_id !== null) {
      message.assetId = object.asset_id;
    }
    if (object.quantums !== undefined && object.quantums !== null) {
      message.quantums = bytesFromBase64(object.quantums);
    }
    if (object.index !== undefined && object.index !== null) {
      message.index = BigInt(object.index);
    }
    return message;
  },
  toAmino(message: IndexerAssetPosition): IndexerAssetPositionAmino {
    const obj: any = {};
    obj.asset_id = message.assetId;
    obj.quantums = message.quantums ? base64FromBytes(message.quantums) : undefined;
    obj.index = message.index ? message.index.toString() : undefined;
    return obj;
  },
  fromAminoMsg(object: IndexerAssetPositionAminoMsg): IndexerAssetPosition {
    return IndexerAssetPosition.fromAmino(object.value);
  },
  fromProtoMsg(message: IndexerAssetPositionProtoMsg): IndexerAssetPosition {
    return IndexerAssetPosition.decode(message.value);
  },
  toProto(message: IndexerAssetPosition): Uint8Array {
    return IndexerAssetPosition.encode(message).finish();
  },
  toProtoMsg(message: IndexerAssetPosition): IndexerAssetPositionProtoMsg {
    return {
      typeUrl: "/dydxprotocol.indexer.protocol.v1.IndexerAssetPosition",
      value: IndexerAssetPosition.encode(message).finish()
    };
  }
};