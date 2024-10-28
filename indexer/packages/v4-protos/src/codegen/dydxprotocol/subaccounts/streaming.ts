import { SubaccountId, SubaccountIdSDKType } from "./subaccount";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial, Long } from "../../helpers";
/**
 * StreamSubaccountUpdate provides information on a subaccount update. Used in
 * the full node GRPC stream.
 */

export interface StreamSubaccountUpdate {
  subaccountId?: SubaccountId;
  /** updated_perpetual_positions will each be for unique perpetuals. */

  updatedPerpetualPositions: SubaccountPerpetualPosition[];
  /** updated_asset_positions will each be for unique assets. */

  updatedAssetPositions: SubaccountAssetPosition[];
  /**
   * Snapshot indicates if the response is from a snapshot of the subaccount.
   * All updates should be ignored until snapshot is received.
   * If the snapshot is true, then all previous entries should be
   * discarded and the subaccount should be resynced.
   * For a snapshot subaccount update, the `updated_perpetual_positions` and
   * `updated_asset_positions` fields will contain the full state of the
   * subaccount.
   */

  snapshot: boolean;
}
/**
 * StreamSubaccountUpdate provides information on a subaccount update. Used in
 * the full node GRPC stream.
 */

export interface StreamSubaccountUpdateSDKType {
  subaccount_id?: SubaccountIdSDKType;
  /** updated_perpetual_positions will each be for unique perpetuals. */

  updated_perpetual_positions: SubaccountPerpetualPositionSDKType[];
  /** updated_asset_positions will each be for unique assets. */

  updated_asset_positions: SubaccountAssetPositionSDKType[];
  /**
   * Snapshot indicates if the response is from a snapshot of the subaccount.
   * All updates should be ignored until snapshot is received.
   * If the snapshot is true, then all previous entries should be
   * discarded and the subaccount should be resynced.
   * For a snapshot subaccount update, the `updated_perpetual_positions` and
   * `updated_asset_positions` fields will contain the full state of the
   * subaccount.
   */

  snapshot: boolean;
}
/**
 * SubaccountPerpetualPosition provides information on a subaccount's updated
 * perpetual positions.
 */

export interface SubaccountPerpetualPosition {
  /** The `Id` of the `Perpetual`. */
  perpetualId: number;
  /** The size of the position in base quantums. Negative means short. */

  quantums: Long;
}
/**
 * SubaccountPerpetualPosition provides information on a subaccount's updated
 * perpetual positions.
 */

export interface SubaccountPerpetualPositionSDKType {
  /** The `Id` of the `Perpetual`. */
  perpetual_id: number;
  /** The size of the position in base quantums. Negative means short. */

  quantums: Long;
}
/**
 * SubaccountAssetPosition provides information on a subaccount's updated asset
 * positions.
 */

export interface SubaccountAssetPosition {
  /** The `Id` of the `Asset`. */
  assetId: number;
  /** The absolute size of the position in base quantums. */

  quantums: Long;
}
/**
 * SubaccountAssetPosition provides information on a subaccount's updated asset
 * positions.
 */

export interface SubaccountAssetPositionSDKType {
  /** The `Id` of the `Asset`. */
  asset_id: number;
  /** The absolute size of the position in base quantums. */

  quantums: Long;
}

function createBaseStreamSubaccountUpdate(): StreamSubaccountUpdate {
  return {
    subaccountId: undefined,
    updatedPerpetualPositions: [],
    updatedAssetPositions: [],
    snapshot: false
  };
}

export const StreamSubaccountUpdate = {
  encode(message: StreamSubaccountUpdate, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.subaccountId !== undefined) {
      SubaccountId.encode(message.subaccountId, writer.uint32(10).fork()).ldelim();
    }

    for (const v of message.updatedPerpetualPositions) {
      SubaccountPerpetualPosition.encode(v!, writer.uint32(18).fork()).ldelim();
    }

    for (const v of message.updatedAssetPositions) {
      SubaccountAssetPosition.encode(v!, writer.uint32(26).fork()).ldelim();
    }

    if (message.snapshot === true) {
      writer.uint32(32).bool(message.snapshot);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): StreamSubaccountUpdate {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseStreamSubaccountUpdate();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.subaccountId = SubaccountId.decode(reader, reader.uint32());
          break;

        case 2:
          message.updatedPerpetualPositions.push(SubaccountPerpetualPosition.decode(reader, reader.uint32()));
          break;

        case 3:
          message.updatedAssetPositions.push(SubaccountAssetPosition.decode(reader, reader.uint32()));
          break;

        case 4:
          message.snapshot = reader.bool();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<StreamSubaccountUpdate>): StreamSubaccountUpdate {
    const message = createBaseStreamSubaccountUpdate();
    message.subaccountId = object.subaccountId !== undefined && object.subaccountId !== null ? SubaccountId.fromPartial(object.subaccountId) : undefined;
    message.updatedPerpetualPositions = object.updatedPerpetualPositions?.map(e => SubaccountPerpetualPosition.fromPartial(e)) || [];
    message.updatedAssetPositions = object.updatedAssetPositions?.map(e => SubaccountAssetPosition.fromPartial(e)) || [];
    message.snapshot = object.snapshot ?? false;
    return message;
  }

};

function createBaseSubaccountPerpetualPosition(): SubaccountPerpetualPosition {
  return {
    perpetualId: 0,
    quantums: Long.ZERO
  };
}

export const SubaccountPerpetualPosition = {
  encode(message: SubaccountPerpetualPosition, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.perpetualId !== 0) {
      writer.uint32(8).uint32(message.perpetualId);
    }

    if (!message.quantums.isZero()) {
      writer.uint32(16).int64(message.quantums);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SubaccountPerpetualPosition {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSubaccountPerpetualPosition();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.perpetualId = reader.uint32();
          break;

        case 2:
          message.quantums = (reader.int64() as Long);
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<SubaccountPerpetualPosition>): SubaccountPerpetualPosition {
    const message = createBaseSubaccountPerpetualPosition();
    message.perpetualId = object.perpetualId ?? 0;
    message.quantums = object.quantums !== undefined && object.quantums !== null ? Long.fromValue(object.quantums) : Long.ZERO;
    return message;
  }

};

function createBaseSubaccountAssetPosition(): SubaccountAssetPosition {
  return {
    assetId: 0,
    quantums: Long.UZERO
  };
}

export const SubaccountAssetPosition = {
  encode(message: SubaccountAssetPosition, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.assetId !== 0) {
      writer.uint32(8).uint32(message.assetId);
    }

    if (!message.quantums.isZero()) {
      writer.uint32(16).uint64(message.quantums);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SubaccountAssetPosition {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSubaccountAssetPosition();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.assetId = reader.uint32();
          break;

        case 2:
          message.quantums = (reader.uint64() as Long);
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<SubaccountAssetPosition>): SubaccountAssetPosition {
    const message = createBaseSubaccountAssetPosition();
    message.assetId = object.assetId ?? 0;
    message.quantums = object.quantums !== undefined && object.quantums !== null ? Long.fromValue(object.quantums) : Long.UZERO;
    return message;
  }

};