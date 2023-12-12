import { BinaryReader, BinaryWriter } from "../../binary";
import { bytesFromBase64, base64FromBytes } from "../../helpers";
/**
 * PerpetualPositions are an account’s positions of a `Perpetual`.
 * Therefore they hold any information needed to trade perpetuals.
 */
export interface PerpetualPosition {
  /** The `Id` of the `Perpetual`. */
  perpetualId: number;
  /** The size of the position in base quantums. */
  quantums: Uint8Array;
  /**
   * The funding_index of the `Perpetual` the last time this position was
   * settled.
   */
  fundingIndex: Uint8Array;
}
export interface PerpetualPositionProtoMsg {
  typeUrl: "/dydxprotocol.subaccounts.PerpetualPosition";
  value: Uint8Array;
}
/**
 * PerpetualPositions are an account’s positions of a `Perpetual`.
 * Therefore they hold any information needed to trade perpetuals.
 */
export interface PerpetualPositionAmino {
  /** The `Id` of the `Perpetual`. */
  perpetual_id?: number;
  /** The size of the position in base quantums. */
  quantums?: string;
  /**
   * The funding_index of the `Perpetual` the last time this position was
   * settled.
   */
  funding_index?: string;
}
export interface PerpetualPositionAminoMsg {
  type: "/dydxprotocol.subaccounts.PerpetualPosition";
  value: PerpetualPositionAmino;
}
/**
 * PerpetualPositions are an account’s positions of a `Perpetual`.
 * Therefore they hold any information needed to trade perpetuals.
 */
export interface PerpetualPositionSDKType {
  perpetual_id: number;
  quantums: Uint8Array;
  funding_index: Uint8Array;
}
function createBasePerpetualPosition(): PerpetualPosition {
  return {
    perpetualId: 0,
    quantums: new Uint8Array(),
    fundingIndex: new Uint8Array()
  };
}
export const PerpetualPosition = {
  typeUrl: "/dydxprotocol.subaccounts.PerpetualPosition",
  encode(message: PerpetualPosition, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.perpetualId !== 0) {
      writer.uint32(8).uint32(message.perpetualId);
    }
    if (message.quantums.length !== 0) {
      writer.uint32(18).bytes(message.quantums);
    }
    if (message.fundingIndex.length !== 0) {
      writer.uint32(26).bytes(message.fundingIndex);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): PerpetualPosition {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBasePerpetualPosition();
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
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<PerpetualPosition>): PerpetualPosition {
    const message = createBasePerpetualPosition();
    message.perpetualId = object.perpetualId ?? 0;
    message.quantums = object.quantums ?? new Uint8Array();
    message.fundingIndex = object.fundingIndex ?? new Uint8Array();
    return message;
  },
  fromAmino(object: PerpetualPositionAmino): PerpetualPosition {
    const message = createBasePerpetualPosition();
    if (object.perpetual_id !== undefined && object.perpetual_id !== null) {
      message.perpetualId = object.perpetual_id;
    }
    if (object.quantums !== undefined && object.quantums !== null) {
      message.quantums = bytesFromBase64(object.quantums);
    }
    if (object.funding_index !== undefined && object.funding_index !== null) {
      message.fundingIndex = bytesFromBase64(object.funding_index);
    }
    return message;
  },
  toAmino(message: PerpetualPosition): PerpetualPositionAmino {
    const obj: any = {};
    obj.perpetual_id = message.perpetualId;
    obj.quantums = message.quantums ? base64FromBytes(message.quantums) : undefined;
    obj.funding_index = message.fundingIndex ? base64FromBytes(message.fundingIndex) : undefined;
    return obj;
  },
  fromAminoMsg(object: PerpetualPositionAminoMsg): PerpetualPosition {
    return PerpetualPosition.fromAmino(object.value);
  },
  fromProtoMsg(message: PerpetualPositionProtoMsg): PerpetualPosition {
    return PerpetualPosition.decode(message.value);
  },
  toProto(message: PerpetualPosition): Uint8Array {
    return PerpetualPosition.encode(message).finish();
  },
  toProtoMsg(message: PerpetualPosition): PerpetualPositionProtoMsg {
    return {
      typeUrl: "/dydxprotocol.subaccounts.PerpetualPosition",
      value: PerpetualPosition.encode(message).finish()
    };
  }
};