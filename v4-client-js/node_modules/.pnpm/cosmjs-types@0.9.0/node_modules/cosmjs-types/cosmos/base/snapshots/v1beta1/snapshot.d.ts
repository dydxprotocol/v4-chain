import { BinaryReader, BinaryWriter } from "../../../../binary";
export declare const protobufPackage = "cosmos.base.snapshots.v1beta1";
/** Snapshot contains Tendermint state sync snapshot info. */
export interface Snapshot {
    height: bigint;
    format: number;
    chunks: number;
    hash: Uint8Array;
    metadata: Metadata;
}
/** Metadata contains SDK-specific snapshot metadata. */
export interface Metadata {
    /** SHA-256 chunk hashes */
    chunkHashes: Uint8Array[];
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
    /** @deprecated */
    kv?: SnapshotKVItem;
    /** @deprecated */
    schema?: SnapshotSchema;
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
 * SnapshotIAVLItem is an exported IAVL node.
 *
 * Since: cosmos-sdk 0.46
 */
export interface SnapshotIAVLItem {
    key: Uint8Array;
    value: Uint8Array;
    /** version is block height */
    version: bigint;
    /** height is depth of the tree. */
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
 * SnapshotExtensionPayload contains payloads of an external snapshotter.
 *
 * Since: cosmos-sdk 0.46
 */
export interface SnapshotExtensionPayload {
    payload: Uint8Array;
}
/**
 * SnapshotKVItem is an exported Key/Value Pair
 *
 * Since: cosmos-sdk 0.46
 * Deprecated: This message was part of store/v2alpha1 which has been deleted from v0.47.
 */
/** @deprecated */
export interface SnapshotKVItem {
    key: Uint8Array;
    value: Uint8Array;
}
/**
 * SnapshotSchema is an exported schema of smt store
 *
 * Since: cosmos-sdk 0.46
 * Deprecated: This message was part of store/v2alpha1 which has been deleted from v0.47.
 */
/** @deprecated */
export interface SnapshotSchema {
    keys: Uint8Array[];
}
export declare const Snapshot: {
    typeUrl: string;
    encode(message: Snapshot, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): Snapshot;
    fromJSON(object: any): Snapshot;
    toJSON(message: Snapshot): unknown;
    fromPartial<I extends {
        height?: bigint | undefined;
        format?: number | undefined;
        chunks?: number | undefined;
        hash?: Uint8Array | undefined;
        metadata?: {
            chunkHashes?: Uint8Array[] | undefined;
        } | undefined;
    } & {
        height?: bigint | undefined;
        format?: number | undefined;
        chunks?: number | undefined;
        hash?: Uint8Array | undefined;
        metadata?: ({
            chunkHashes?: Uint8Array[] | undefined;
        } & {
            chunkHashes?: (Uint8Array[] & Uint8Array[] & Record<Exclude<keyof I["metadata"]["chunkHashes"], keyof Uint8Array[]>, never>) | undefined;
        } & Record<Exclude<keyof I["metadata"], "chunkHashes">, never>) | undefined;
    } & Record<Exclude<keyof I, keyof Snapshot>, never>>(object: I): Snapshot;
};
export declare const Metadata: {
    typeUrl: string;
    encode(message: Metadata, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): Metadata;
    fromJSON(object: any): Metadata;
    toJSON(message: Metadata): unknown;
    fromPartial<I extends {
        chunkHashes?: Uint8Array[] | undefined;
    } & {
        chunkHashes?: (Uint8Array[] & Uint8Array[] & Record<Exclude<keyof I["chunkHashes"], keyof Uint8Array[]>, never>) | undefined;
    } & Record<Exclude<keyof I, "chunkHashes">, never>>(object: I): Metadata;
};
export declare const SnapshotItem: {
    typeUrl: string;
    encode(message: SnapshotItem, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): SnapshotItem;
    fromJSON(object: any): SnapshotItem;
    toJSON(message: SnapshotItem): unknown;
    fromPartial<I extends {
        store?: {
            name?: string | undefined;
        } | undefined;
        iavl?: {
            key?: Uint8Array | undefined;
            value?: Uint8Array | undefined;
            version?: bigint | undefined;
            height?: number | undefined;
        } | undefined;
        extension?: {
            name?: string | undefined;
            format?: number | undefined;
        } | undefined;
        extensionPayload?: {
            payload?: Uint8Array | undefined;
        } | undefined;
        kv?: {
            key?: Uint8Array | undefined;
            value?: Uint8Array | undefined;
        } | undefined;
        schema?: {
            keys?: Uint8Array[] | undefined;
        } | undefined;
    } & {
        store?: ({
            name?: string | undefined;
        } & {
            name?: string | undefined;
        } & Record<Exclude<keyof I["store"], "name">, never>) | undefined;
        iavl?: ({
            key?: Uint8Array | undefined;
            value?: Uint8Array | undefined;
            version?: bigint | undefined;
            height?: number | undefined;
        } & {
            key?: Uint8Array | undefined;
            value?: Uint8Array | undefined;
            version?: bigint | undefined;
            height?: number | undefined;
        } & Record<Exclude<keyof I["iavl"], keyof SnapshotIAVLItem>, never>) | undefined;
        extension?: ({
            name?: string | undefined;
            format?: number | undefined;
        } & {
            name?: string | undefined;
            format?: number | undefined;
        } & Record<Exclude<keyof I["extension"], keyof SnapshotExtensionMeta>, never>) | undefined;
        extensionPayload?: ({
            payload?: Uint8Array | undefined;
        } & {
            payload?: Uint8Array | undefined;
        } & Record<Exclude<keyof I["extensionPayload"], "payload">, never>) | undefined;
        kv?: ({
            key?: Uint8Array | undefined;
            value?: Uint8Array | undefined;
        } & {
            key?: Uint8Array | undefined;
            value?: Uint8Array | undefined;
        } & Record<Exclude<keyof I["kv"], keyof SnapshotKVItem>, never>) | undefined;
        schema?: ({
            keys?: Uint8Array[] | undefined;
        } & {
            keys?: (Uint8Array[] & Uint8Array[] & Record<Exclude<keyof I["schema"]["keys"], keyof Uint8Array[]>, never>) | undefined;
        } & Record<Exclude<keyof I["schema"], "keys">, never>) | undefined;
    } & Record<Exclude<keyof I, keyof SnapshotItem>, never>>(object: I): SnapshotItem;
};
export declare const SnapshotStoreItem: {
    typeUrl: string;
    encode(message: SnapshotStoreItem, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): SnapshotStoreItem;
    fromJSON(object: any): SnapshotStoreItem;
    toJSON(message: SnapshotStoreItem): unknown;
    fromPartial<I extends {
        name?: string | undefined;
    } & {
        name?: string | undefined;
    } & Record<Exclude<keyof I, "name">, never>>(object: I): SnapshotStoreItem;
};
export declare const SnapshotIAVLItem: {
    typeUrl: string;
    encode(message: SnapshotIAVLItem, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): SnapshotIAVLItem;
    fromJSON(object: any): SnapshotIAVLItem;
    toJSON(message: SnapshotIAVLItem): unknown;
    fromPartial<I extends {
        key?: Uint8Array | undefined;
        value?: Uint8Array | undefined;
        version?: bigint | undefined;
        height?: number | undefined;
    } & {
        key?: Uint8Array | undefined;
        value?: Uint8Array | undefined;
        version?: bigint | undefined;
        height?: number | undefined;
    } & Record<Exclude<keyof I, keyof SnapshotIAVLItem>, never>>(object: I): SnapshotIAVLItem;
};
export declare const SnapshotExtensionMeta: {
    typeUrl: string;
    encode(message: SnapshotExtensionMeta, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): SnapshotExtensionMeta;
    fromJSON(object: any): SnapshotExtensionMeta;
    toJSON(message: SnapshotExtensionMeta): unknown;
    fromPartial<I extends {
        name?: string | undefined;
        format?: number | undefined;
    } & {
        name?: string | undefined;
        format?: number | undefined;
    } & Record<Exclude<keyof I, keyof SnapshotExtensionMeta>, never>>(object: I): SnapshotExtensionMeta;
};
export declare const SnapshotExtensionPayload: {
    typeUrl: string;
    encode(message: SnapshotExtensionPayload, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): SnapshotExtensionPayload;
    fromJSON(object: any): SnapshotExtensionPayload;
    toJSON(message: SnapshotExtensionPayload): unknown;
    fromPartial<I extends {
        payload?: Uint8Array | undefined;
    } & {
        payload?: Uint8Array | undefined;
    } & Record<Exclude<keyof I, "payload">, never>>(object: I): SnapshotExtensionPayload;
};
export declare const SnapshotKVItem: {
    typeUrl: string;
    encode(message: SnapshotKVItem, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): SnapshotKVItem;
    fromJSON(object: any): SnapshotKVItem;
    toJSON(message: SnapshotKVItem): unknown;
    fromPartial<I extends {
        key?: Uint8Array | undefined;
        value?: Uint8Array | undefined;
    } & {
        key?: Uint8Array | undefined;
        value?: Uint8Array | undefined;
    } & Record<Exclude<keyof I, keyof SnapshotKVItem>, never>>(object: I): SnapshotKVItem;
};
export declare const SnapshotSchema: {
    typeUrl: string;
    encode(message: SnapshotSchema, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): SnapshotSchema;
    fromJSON(object: any): SnapshotSchema;
    toJSON(message: SnapshotSchema): unknown;
    fromPartial<I extends {
        keys?: Uint8Array[] | undefined;
    } & {
        keys?: (Uint8Array[] & Uint8Array[] & Record<Exclude<keyof I["keys"], keyof Uint8Array[]>, never>) | undefined;
    } & Record<Exclude<keyof I, "keys">, never>>(object: I): SnapshotSchema;
};
