/// <reference types="long" />
import { Long, DeepPartial } from "../../../../helpers";
import * as _m0 from "protobufjs/minimal";
/** Snapshot contains Tendermint state sync snapshot info. */
export interface Snapshot {
    height: Long;
    format: number;
    chunks: number;
    hash: Uint8Array;
    metadata?: Metadata;
}
/** Snapshot contains Tendermint state sync snapshot info. */
export interface SnapshotSDKType {
    height: Long;
    format: number;
    chunks: number;
    hash: Uint8Array;
    metadata?: MetadataSDKType;
}
/** Metadata contains SDK-specific snapshot metadata. */
export interface Metadata {
    /** SHA-256 chunk hashes */
    chunkHashes: Uint8Array[];
}
/** Metadata contains SDK-specific snapshot metadata. */
export interface MetadataSDKType {
    chunk_hashes: Uint8Array[];
}
/**
 * SnapshotItem is an item contained in a rootmulti.Store snapshot.
 *
 * Since: cosmos-sdk 0.46
 */
export interface SnapshotItem {
    store?: SnapshotStoreItem;
    iavl?: SnapshotIAVLItem;
    extension?: SnapshotExtensionMeta;
    extensionPayload?: SnapshotExtensionPayload;
}
/**
 * SnapshotItem is an item contained in a rootmulti.Store snapshot.
 *
 * Since: cosmos-sdk 0.46
 */
export interface SnapshotItemSDKType {
    store?: SnapshotStoreItemSDKType;
    iavl?: SnapshotIAVLItemSDKType;
    extension?: SnapshotExtensionMetaSDKType;
    extension_payload?: SnapshotExtensionPayloadSDKType;
}
/**
 * SnapshotStoreItem contains metadata about a snapshotted store.
 *
 * Since: cosmos-sdk 0.46
 */
export interface SnapshotStoreItem {
    name: string;
}
/**
 * SnapshotStoreItem contains metadata about a snapshotted store.
 *
 * Since: cosmos-sdk 0.46
 */
export interface SnapshotStoreItemSDKType {
    name: string;
}
/**
 * SnapshotIAVLItem is an exported IAVL node.
 *
 * Since: cosmos-sdk 0.46
 */
export interface SnapshotIAVLItem {
    key: Uint8Array;
    value: Uint8Array;
    /** version is block height */
    version: Long;
    /** height is depth of the tree. */
    height: number;
}
/**
 * SnapshotIAVLItem is an exported IAVL node.
 *
 * Since: cosmos-sdk 0.46
 */
export interface SnapshotIAVLItemSDKType {
    key: Uint8Array;
    value: Uint8Array;
    version: Long;
    height: number;
}
/**
 * SnapshotExtensionMeta contains metadata about an external snapshotter.
 *
 * Since: cosmos-sdk 0.46
 */
export interface SnapshotExtensionMeta {
    name: string;
    format: number;
}
/**
 * SnapshotExtensionMeta contains metadata about an external snapshotter.
 *
 * Since: cosmos-sdk 0.46
 */
export interface SnapshotExtensionMetaSDKType {
    name: string;
    format: number;
}
/**
 * SnapshotExtensionPayload contains payloads of an external snapshotter.
 *
 * Since: cosmos-sdk 0.46
 */
export interface SnapshotExtensionPayload {
    payload: Uint8Array;
}
/**
 * SnapshotExtensionPayload contains payloads of an external snapshotter.
 *
 * Since: cosmos-sdk 0.46
 */
export interface SnapshotExtensionPayloadSDKType {
    payload: Uint8Array;
}
export declare const Snapshot: {
    encode(message: Snapshot, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): Snapshot;
    fromPartial(object: DeepPartial<Snapshot>): Snapshot;
};
export declare const Metadata: {
    encode(message: Metadata, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): Metadata;
    fromPartial(object: DeepPartial<Metadata>): Metadata;
};
export declare const SnapshotItem: {
    encode(message: SnapshotItem, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): SnapshotItem;
    fromPartial(object: DeepPartial<SnapshotItem>): SnapshotItem;
};
export declare const SnapshotStoreItem: {
    encode(message: SnapshotStoreItem, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): SnapshotStoreItem;
    fromPartial(object: DeepPartial<SnapshotStoreItem>): SnapshotStoreItem;
};
export declare const SnapshotIAVLItem: {
    encode(message: SnapshotIAVLItem, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): SnapshotIAVLItem;
    fromPartial(object: DeepPartial<SnapshotIAVLItem>): SnapshotIAVLItem;
};
export declare const SnapshotExtensionMeta: {
    encode(message: SnapshotExtensionMeta, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): SnapshotExtensionMeta;
    fromPartial(object: DeepPartial<SnapshotExtensionMeta>): SnapshotExtensionMeta;
};
export declare const SnapshotExtensionPayload: {
    encode(message: SnapshotExtensionPayload, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): SnapshotExtensionPayload;
    fromPartial(object: DeepPartial<SnapshotExtensionPayload>): SnapshotExtensionPayload;
};
