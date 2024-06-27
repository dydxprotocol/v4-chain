/// <reference types="long" />
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
    index: Long;
}
/**
 * IndexerAssetPosition define an account’s positions of an `Asset`.
 * Therefore they hold any information needed to trade on Spot and Margin.
 */
export interface IndexerAssetPositionSDKType {
    asset_id: number;
    quantums: Uint8Array;
    index: Long;
}
export declare const IndexerSubaccountId: {
    encode(message: IndexerSubaccountId, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): IndexerSubaccountId;
    fromPartial(object: DeepPartial<IndexerSubaccountId>): IndexerSubaccountId;
};
export declare const IndexerPerpetualPosition: {
    encode(message: IndexerPerpetualPosition, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): IndexerPerpetualPosition;
    fromPartial(object: DeepPartial<IndexerPerpetualPosition>): IndexerPerpetualPosition;
};
export declare const IndexerAssetPosition: {
    encode(message: IndexerAssetPosition, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): IndexerAssetPosition;
    fromPartial(object: DeepPartial<IndexerAssetPosition>): IndexerAssetPosition;
};
