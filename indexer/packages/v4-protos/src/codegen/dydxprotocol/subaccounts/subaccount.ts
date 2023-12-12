import { AssetPosition, AssetPositionAmino, AssetPositionSDKType } from "./asset_position";
import { PerpetualPosition, PerpetualPositionAmino, PerpetualPositionSDKType } from "./perpetual_position";
import { BinaryReader, BinaryWriter } from "../../binary";
/** SubaccountId defines a unique identifier for a Subaccount. */
export interface SubaccountId {
  /** The address of the wallet that owns this subaccount. */
  owner: string;
  /**
   * < 128 Since 128 should be enough to start and it fits within
   * 1 Byte (1 Bit needed to indicate that the first byte is the last).
   */
  number: number;
}
export interface SubaccountIdProtoMsg {
  typeUrl: "/dydxprotocol.subaccounts.SubaccountId";
  value: Uint8Array;
}
/** SubaccountId defines a unique identifier for a Subaccount. */
export interface SubaccountIdAmino {
  /** The address of the wallet that owns this subaccount. */
  owner?: string;
  /**
   * < 128 Since 128 should be enough to start and it fits within
   * 1 Byte (1 Bit needed to indicate that the first byte is the last).
   */
  number?: number;
}
export interface SubaccountIdAminoMsg {
  type: "/dydxprotocol.subaccounts.SubaccountId";
  value: SubaccountIdAmino;
}
/** SubaccountId defines a unique identifier for a Subaccount. */
export interface SubaccountIdSDKType {
  owner: string;
  number: number;
}
/**
 * Subaccount defines a single sub-account for a given address.
 * Subaccounts are uniquely indexed by a subaccountNumber/owner pair.
 */
export interface Subaccount {
  /** The Id of the Subaccount */
  id?: SubaccountId;
  /**
   * All `AssetPosition`s associated with this subaccount.
   * Always sorted ascending by `asset_id`.
   */
  assetPositions: AssetPosition[];
  /**
   * All `PerpetualPosition`s associated with this subaccount.
   * Always sorted ascending by `perpetual_id.
   */
  perpetualPositions: PerpetualPosition[];
  /**
   * Set by the owner. If true, then margin trades can be made in this
   * subaccount.
   */
  marginEnabled: boolean;
}
export interface SubaccountProtoMsg {
  typeUrl: "/dydxprotocol.subaccounts.Subaccount";
  value: Uint8Array;
}
/**
 * Subaccount defines a single sub-account for a given address.
 * Subaccounts are uniquely indexed by a subaccountNumber/owner pair.
 */
export interface SubaccountAmino {
  /** The Id of the Subaccount */
  id?: SubaccountIdAmino;
  /**
   * All `AssetPosition`s associated with this subaccount.
   * Always sorted ascending by `asset_id`.
   */
  asset_positions?: AssetPositionAmino[];
  /**
   * All `PerpetualPosition`s associated with this subaccount.
   * Always sorted ascending by `perpetual_id.
   */
  perpetual_positions?: PerpetualPositionAmino[];
  /**
   * Set by the owner. If true, then margin trades can be made in this
   * subaccount.
   */
  margin_enabled?: boolean;
}
export interface SubaccountAminoMsg {
  type: "/dydxprotocol.subaccounts.Subaccount";
  value: SubaccountAmino;
}
/**
 * Subaccount defines a single sub-account for a given address.
 * Subaccounts are uniquely indexed by a subaccountNumber/owner pair.
 */
export interface SubaccountSDKType {
  id?: SubaccountIdSDKType;
  asset_positions: AssetPositionSDKType[];
  perpetual_positions: PerpetualPositionSDKType[];
  margin_enabled: boolean;
}
function createBaseSubaccountId(): SubaccountId {
  return {
    owner: "",
    number: 0
  };
}
export const SubaccountId = {
  typeUrl: "/dydxprotocol.subaccounts.SubaccountId",
  encode(message: SubaccountId, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.owner !== "") {
      writer.uint32(10).string(message.owner);
    }
    if (message.number !== 0) {
      writer.uint32(16).uint32(message.number);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): SubaccountId {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSubaccountId();
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
  fromPartial(object: Partial<SubaccountId>): SubaccountId {
    const message = createBaseSubaccountId();
    message.owner = object.owner ?? "";
    message.number = object.number ?? 0;
    return message;
  },
  fromAmino(object: SubaccountIdAmino): SubaccountId {
    const message = createBaseSubaccountId();
    if (object.owner !== undefined && object.owner !== null) {
      message.owner = object.owner;
    }
    if (object.number !== undefined && object.number !== null) {
      message.number = object.number;
    }
    return message;
  },
  toAmino(message: SubaccountId): SubaccountIdAmino {
    const obj: any = {};
    obj.owner = message.owner;
    obj.number = message.number;
    return obj;
  },
  fromAminoMsg(object: SubaccountIdAminoMsg): SubaccountId {
    return SubaccountId.fromAmino(object.value);
  },
  fromProtoMsg(message: SubaccountIdProtoMsg): SubaccountId {
    return SubaccountId.decode(message.value);
  },
  toProto(message: SubaccountId): Uint8Array {
    return SubaccountId.encode(message).finish();
  },
  toProtoMsg(message: SubaccountId): SubaccountIdProtoMsg {
    return {
      typeUrl: "/dydxprotocol.subaccounts.SubaccountId",
      value: SubaccountId.encode(message).finish()
    };
  }
};
function createBaseSubaccount(): Subaccount {
  return {
    id: undefined,
    assetPositions: [],
    perpetualPositions: [],
    marginEnabled: false
  };
}
export const Subaccount = {
  typeUrl: "/dydxprotocol.subaccounts.Subaccount",
  encode(message: Subaccount, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.id !== undefined) {
      SubaccountId.encode(message.id, writer.uint32(10).fork()).ldelim();
    }
    for (const v of message.assetPositions) {
      AssetPosition.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    for (const v of message.perpetualPositions) {
      PerpetualPosition.encode(v!, writer.uint32(26).fork()).ldelim();
    }
    if (message.marginEnabled === true) {
      writer.uint32(32).bool(message.marginEnabled);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): Subaccount {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSubaccount();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.id = SubaccountId.decode(reader, reader.uint32());
          break;
        case 2:
          message.assetPositions.push(AssetPosition.decode(reader, reader.uint32()));
          break;
        case 3:
          message.perpetualPositions.push(PerpetualPosition.decode(reader, reader.uint32()));
          break;
        case 4:
          message.marginEnabled = reader.bool();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<Subaccount>): Subaccount {
    const message = createBaseSubaccount();
    message.id = object.id !== undefined && object.id !== null ? SubaccountId.fromPartial(object.id) : undefined;
    message.assetPositions = object.assetPositions?.map(e => AssetPosition.fromPartial(e)) || [];
    message.perpetualPositions = object.perpetualPositions?.map(e => PerpetualPosition.fromPartial(e)) || [];
    message.marginEnabled = object.marginEnabled ?? false;
    return message;
  },
  fromAmino(object: SubaccountAmino): Subaccount {
    const message = createBaseSubaccount();
    if (object.id !== undefined && object.id !== null) {
      message.id = SubaccountId.fromAmino(object.id);
    }
    message.assetPositions = object.asset_positions?.map(e => AssetPosition.fromAmino(e)) || [];
    message.perpetualPositions = object.perpetual_positions?.map(e => PerpetualPosition.fromAmino(e)) || [];
    if (object.margin_enabled !== undefined && object.margin_enabled !== null) {
      message.marginEnabled = object.margin_enabled;
    }
    return message;
  },
  toAmino(message: Subaccount): SubaccountAmino {
    const obj: any = {};
    obj.id = message.id ? SubaccountId.toAmino(message.id) : undefined;
    if (message.assetPositions) {
      obj.asset_positions = message.assetPositions.map(e => e ? AssetPosition.toAmino(e) : undefined);
    } else {
      obj.asset_positions = [];
    }
    if (message.perpetualPositions) {
      obj.perpetual_positions = message.perpetualPositions.map(e => e ? PerpetualPosition.toAmino(e) : undefined);
    } else {
      obj.perpetual_positions = [];
    }
    obj.margin_enabled = message.marginEnabled;
    return obj;
  },
  fromAminoMsg(object: SubaccountAminoMsg): Subaccount {
    return Subaccount.fromAmino(object.value);
  },
  fromProtoMsg(message: SubaccountProtoMsg): Subaccount {
    return Subaccount.decode(message.value);
  },
  toProto(message: Subaccount): Uint8Array {
    return Subaccount.encode(message).finish();
  },
  toProtoMsg(message: Subaccount): SubaccountProtoMsg {
    return {
      typeUrl: "/dydxprotocol.subaccounts.Subaccount",
      value: Subaccount.encode(message).finish()
    };
  }
};