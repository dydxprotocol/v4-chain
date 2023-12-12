import { BinaryReader, BinaryWriter } from "../../binary";
import { bytesFromBase64, base64FromBytes } from "../../helpers";
/**
 * AssetPositions define an account’s positions of an `Asset`.
 * Therefore they hold any information needed to trade on Spot and Margin.
 */
export interface AssetPosition {
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
export interface AssetPositionProtoMsg {
  typeUrl: "/dydxprotocol.subaccounts.AssetPosition";
  value: Uint8Array;
}
/**
 * AssetPositions define an account’s positions of an `Asset`.
 * Therefore they hold any information needed to trade on Spot and Margin.
 */
export interface AssetPositionAmino {
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
export interface AssetPositionAminoMsg {
  type: "/dydxprotocol.subaccounts.AssetPosition";
  value: AssetPositionAmino;
}
/**
 * AssetPositions define an account’s positions of an `Asset`.
 * Therefore they hold any information needed to trade on Spot and Margin.
 */
export interface AssetPositionSDKType {
  asset_id: number;
  quantums: Uint8Array;
  index: bigint;
}
function createBaseAssetPosition(): AssetPosition {
  return {
    assetId: 0,
    quantums: new Uint8Array(),
    index: BigInt(0)
  };
}
export const AssetPosition = {
  typeUrl: "/dydxprotocol.subaccounts.AssetPosition",
  encode(message: AssetPosition, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
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
  decode(input: BinaryReader | Uint8Array, length?: number): AssetPosition {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAssetPosition();
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
  fromPartial(object: Partial<AssetPosition>): AssetPosition {
    const message = createBaseAssetPosition();
    message.assetId = object.assetId ?? 0;
    message.quantums = object.quantums ?? new Uint8Array();
    message.index = object.index !== undefined && object.index !== null ? BigInt(object.index.toString()) : BigInt(0);
    return message;
  },
  fromAmino(object: AssetPositionAmino): AssetPosition {
    const message = createBaseAssetPosition();
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
  toAmino(message: AssetPosition): AssetPositionAmino {
    const obj: any = {};
    obj.asset_id = message.assetId;
    obj.quantums = message.quantums ? base64FromBytes(message.quantums) : undefined;
    obj.index = message.index ? message.index.toString() : undefined;
    return obj;
  },
  fromAminoMsg(object: AssetPositionAminoMsg): AssetPosition {
    return AssetPosition.fromAmino(object.value);
  },
  fromProtoMsg(message: AssetPositionProtoMsg): AssetPosition {
    return AssetPosition.decode(message.value);
  },
  toProto(message: AssetPosition): Uint8Array {
    return AssetPosition.encode(message).finish();
  },
  toProtoMsg(message: AssetPosition): AssetPositionProtoMsg {
    return {
      typeUrl: "/dydxprotocol.subaccounts.AssetPosition",
      value: AssetPosition.encode(message).finish()
    };
  }
};