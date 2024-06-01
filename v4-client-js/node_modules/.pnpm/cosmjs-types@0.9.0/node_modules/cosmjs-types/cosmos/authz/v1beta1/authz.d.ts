import { Any } from "../../../google/protobuf/any";
import { Timestamp } from "../../../google/protobuf/timestamp";
import { BinaryReader, BinaryWriter } from "../../../binary";
export declare const protobufPackage = "cosmos.authz.v1beta1";
/**
 * GenericAuthorization gives the grantee unrestricted permissions to execute
 * the provided method on behalf of the granter's account.
 */
export interface GenericAuthorization {
    /** Msg, identified by it's type URL, to grant unrestricted permissions to execute */
    msg: string;
}
/**
 * Grant gives permissions to execute
 * the provide method with expiration time.
 */
export interface Grant {
    authorization?: Any;
    /**
     * time when the grant will expire and will be pruned. If null, then the grant
     * doesn't have a time expiration (other conditions  in `authorization`
     * may apply to invalidate the grant)
     */
    expiration?: Timestamp;
}
/**
 * GrantAuthorization extends a grant with both the addresses of the grantee and granter.
 * It is used in genesis.proto and query.proto
 */
export interface GrantAuthorization {
    granter: string;
    grantee: string;
    authorization?: Any;
    expiration?: Timestamp;
}
/** GrantQueueItem contains the list of TypeURL of a sdk.Msg. */
export interface GrantQueueItem {
    /** msg_type_urls contains the list of TypeURL of a sdk.Msg. */
    msgTypeUrls: string[];
}
export declare const GenericAuthorization: {
    typeUrl: string;
    encode(message: GenericAuthorization, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): GenericAuthorization;
    fromJSON(object: any): GenericAuthorization;
    toJSON(message: GenericAuthorization): unknown;
    fromPartial<I extends {
        msg?: string | undefined;
    } & {
        msg?: string | undefined;
    } & Record<Exclude<keyof I, "msg">, never>>(object: I): GenericAuthorization;
};
export declare const Grant: {
    typeUrl: string;
    encode(message: Grant, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): Grant;
    fromJSON(object: any): Grant;
    toJSON(message: Grant): unknown;
    fromPartial<I extends {
        authorization?: {
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } | undefined;
        expiration?: {
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } | undefined;
    } & {
        authorization?: ({
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } & {
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } & Record<Exclude<keyof I["authorization"], keyof Any>, never>) | undefined;
        expiration?: ({
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } & {
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } & Record<Exclude<keyof I["expiration"], keyof Timestamp>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof Grant>, never>>(object: I): Grant;
};
export declare const GrantAuthorization: {
    typeUrl: string;
    encode(message: GrantAuthorization, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): GrantAuthorization;
    fromJSON(object: any): GrantAuthorization;
    toJSON(message: GrantAuthorization): unknown;
    fromPartial<I extends {
        granter?: string | undefined;
        grantee?: string | undefined;
        authorization?: {
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } | undefined;
        expiration?: {
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } | undefined;
    } & {
        granter?: string | undefined;
        grantee?: string | undefined;
        authorization?: ({
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } & {
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } & Record<Exclude<keyof I["authorization"], keyof Any>, never>) | undefined;
        expiration?: ({
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } & {
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } & Record<Exclude<keyof I["expiration"], keyof Timestamp>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof GrantAuthorization>, never>>(object: I): GrantAuthorization;
};
export declare const GrantQueueItem: {
    typeUrl: string;
    encode(message: GrantQueueItem, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): GrantQueueItem;
    fromJSON(object: any): GrantQueueItem;
    toJSON(message: GrantQueueItem): unknown;
    fromPartial<I extends {
        msgTypeUrls?: string[] | undefined;
    } & {
        msgTypeUrls?: (string[] & string[] & Record<Exclude<keyof I["msgTypeUrls"], keyof string[]>, never>) | undefined;
    } & Record<Exclude<keyof I, "msgTypeUrls">, never>>(object: I): GrantQueueItem;
};
