import * as _m0 from "protobufjs/minimal";
import { DeepPartial, Long } from "../../../../helpers";
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
/** IndexerSubaccountId defines a unique identifier for a Subaccount. */

export interface IndexerSubaccountIdSDKType {
  /** The address of the wallet that owns this subaccount. */
  owner: string;
  /**
   * < 128 Since 128 should be enough to start and it fits within
   * 1 Byte (1 Bit needed to indicate that the first byte is the last).
   */

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
/**
 * IndexerPerpetualPosition are an account’s positions of a `Perpetual`.
 * Therefore they hold any information needed to trade perpetuals.
 */

export interface IndexerPerpetualPositionSDKType {
  /** The `Id` of the `Perpetual`. */
  perpetual_id: number;
  /** The size of the position in base quantums. */

  quantums: Uint8Array;
  /**
   * The funding_index of the `Perpetual` the last time this position was
   * settled.
   */

  funding_index: Uint8Array;
  /**
   * Amount of funding payment (in quote quantums).
   * Note: 1. this field is not cumulative.
   * 2. a positive value means funding payment was paid out and
   * a negative value means funding payment was received.
   */

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

  index: Long;
}
/**
 * IndexerAssetPosition define an account’s positions of an `Asset`.
 * Therefore they hold any information needed to trade on Spot and Margin.
 */

export interface IndexerAssetPositionSDKType {
  /** The `Id` of the `Asset`. */
  asset_id: number;
  /** The absolute size of the position in base quantums. */

  quantums: Uint8Array;
  /**
   * The `Index` (either `LongIndex` or `ShortIndex`) of the `Asset` the last
   * time this position was settled
   * TODO(DEC-582): pending margin trading being added.
   */

  index: Long;
}

function createBaseIndexerSubaccountId(): IndexerSubaccountId {
  return {
    owner: "",
    number: 0
  };
}

export const IndexerSubaccountId = {
  encode(message: IndexerSubaccountId, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.owner !== "") {
      writer.uint32(10).string(message.owner);
    }

    if (message.number !== 0) {
      writer.uint32(16).uint32(message.number);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): IndexerSubaccountId {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
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

  fromPartial(object: DeepPartial<IndexerSubaccountId>): IndexerSubaccountId {
    const message = createBaseIndexerSubaccountId();
    message.owner = object.owner ?? "";
    message.number = object.number ?? 0;
    return message;
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
  encode(message: IndexerPerpetualPosition, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
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

  decode(input: _m0.Reader | Uint8Array, length?: number): IndexerPerpetualPosition {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
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

  fromPartial(object: DeepPartial<IndexerPerpetualPosition>): IndexerPerpetualPosition {
    const message = createBaseIndexerPerpetualPosition();
    message.perpetualId = object.perpetualId ?? 0;
    message.quantums = object.quantums ?? new Uint8Array();
    message.fundingIndex = object.fundingIndex ?? new Uint8Array();
    message.fundingPayment = object.fundingPayment ?? new Uint8Array();
    return message;
  }

};

function createBaseIndexerAssetPosition(): IndexerAssetPosition {
  return {
    assetId: 0,
    quantums: new Uint8Array(),
    index: Long.UZERO
  };
}

export const IndexerAssetPosition = {
  encode(message: IndexerAssetPosition, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.assetId !== 0) {
      writer.uint32(8).uint32(message.assetId);
    }

    if (message.quantums.length !== 0) {
      writer.uint32(18).bytes(message.quantums);
    }

    if (!message.index.isZero()) {
      writer.uint32(24).uint64(message.index);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): IndexerAssetPosition {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
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
          message.index = (reader.uint64() as Long);
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<IndexerAssetPosition>): IndexerAssetPosition {
    const message = createBaseIndexerAssetPosition();
    message.assetId = object.assetId ?? 0;
    message.quantums = object.quantums ?? new Uint8Array();
    message.index = object.index !== undefined && object.index !== null ? Long.fromValue(object.index) : Long.UZERO;
    return message;
  }

};