import { Timestamp } from "../../../../google/protobuf/timestamp";
import { BinaryReader, BinaryWriter } from "../../../../binary";
export declare const protobufPackage = "cosmos.base.store.v1beta1";
/**
 * CommitInfo defines commit information used by the multi-store when committing
 * a version/height.
 */
export interface CommitInfo {
    version: bigint;
    storeInfos: StoreInfo[];
    timestamp: Timestamp;
}
/**
 * StoreInfo defines store-specific commit information. It contains a reference
 * between a store name and the commit ID.
 */
export interface StoreInfo {
    name: string;
    commitId: CommitID;
}
/**
 * CommitID defines the commitment information when a specific store is
 * committed.
 */
export interface CommitID {
    version: bigint;
    hash: Uint8Array;
}
export declare const CommitInfo: {
    typeUrl: string;
    encode(message: CommitInfo, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): CommitInfo;
    fromJSON(object: any): CommitInfo;
    toJSON(message: CommitInfo): unknown;
    fromPartial<I extends {
        version?: bigint | undefined;
        storeInfos?: {
            name?: string | undefined;
            commitId?: {
                version?: bigint | undefined;
                hash?: Uint8Array | undefined;
            } | undefined;
        }[] | undefined;
        timestamp?: {
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } | undefined;
    } & {
        version?: bigint | undefined;
        storeInfos?: ({
            name?: string | undefined;
            commitId?: {
                version?: bigint | undefined;
                hash?: Uint8Array | undefined;
            } | undefined;
        }[] & ({
            name?: string | undefined;
            commitId?: {
                version?: bigint | undefined;
                hash?: Uint8Array | undefined;
            } | undefined;
        } & {
            name?: string | undefined;
            commitId?: ({
                version?: bigint | undefined;
                hash?: Uint8Array | undefined;
            } & {
                version?: bigint | undefined;
                hash?: Uint8Array | undefined;
            } & Record<Exclude<keyof I["storeInfos"][number]["commitId"], keyof CommitID>, never>) | undefined;
        } & Record<Exclude<keyof I["storeInfos"][number], keyof StoreInfo>, never>)[] & Record<Exclude<keyof I["storeInfos"], keyof {
            name?: string | undefined;
            commitId?: {
                version?: bigint | undefined;
                hash?: Uint8Array | undefined;
            } | undefined;
        }[]>, never>) | undefined;
        timestamp?: ({
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } & {
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } & Record<Exclude<keyof I["timestamp"], keyof Timestamp>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof CommitInfo>, never>>(object: I): CommitInfo;
};
export declare const StoreInfo: {
    typeUrl: string;
    encode(message: StoreInfo, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): StoreInfo;
    fromJSON(object: any): StoreInfo;
    toJSON(message: StoreInfo): unknown;
    fromPartial<I extends {
        name?: string | undefined;
        commitId?: {
            version?: bigint | undefined;
            hash?: Uint8Array | undefined;
        } | undefined;
    } & {
        name?: string | undefined;
        commitId?: ({
            version?: bigint | undefined;
            hash?: Uint8Array | undefined;
        } & {
            version?: bigint | undefined;
            hash?: Uint8Array | undefined;
        } & Record<Exclude<keyof I["commitId"], keyof CommitID>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof StoreInfo>, never>>(object: I): StoreInfo;
};
export declare const CommitID: {
    typeUrl: string;
    encode(message: CommitID, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): CommitID;
    fromJSON(object: any): CommitID;
    toJSON(message: CommitID): unknown;
    fromPartial<I extends {
        version?: bigint | undefined;
        hash?: Uint8Array | undefined;
    } & {
        version?: bigint | undefined;
        hash?: Uint8Array | undefined;
    } & Record<Exclude<keyof I, keyof CommitID>, never>>(object: I): CommitID;
};
