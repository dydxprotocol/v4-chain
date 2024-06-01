import { Grant } from "./authz";
import { Any } from "../../../google/protobuf/any";
import { BinaryReader, BinaryWriter } from "../../../binary";
import { Rpc } from "../../../helpers";
export declare const protobufPackage = "cosmos.authz.v1beta1";
/**
 * MsgGrant is a request type for Grant method. It declares authorization to the grantee
 * on behalf of the granter with the provided expiration time.
 */
export interface MsgGrant {
    granter: string;
    grantee: string;
    grant: Grant;
}
/** MsgExecResponse defines the Msg/MsgExecResponse response type. */
export interface MsgExecResponse {
    results: Uint8Array[];
}
/**
 * MsgExec attempts to execute the provided messages using
 * authorizations granted to the grantee. Each message should have only
 * one signer corresponding to the granter of the authorization.
 */
export interface MsgExec {
    grantee: string;
    /**
     * Execute Msg.
     * The x/authz will try to find a grant matching (msg.signers[0], grantee, MsgTypeURL(msg))
     * triple and validate it.
     */
    msgs: Any[];
}
/** MsgGrantResponse defines the Msg/MsgGrant response type. */
export interface MsgGrantResponse {
}
/**
 * MsgRevoke revokes any authorization with the provided sdk.Msg type on the
 * granter's account with that has been granted to the grantee.
 */
export interface MsgRevoke {
    granter: string;
    grantee: string;
    msgTypeUrl: string;
}
/** MsgRevokeResponse defines the Msg/MsgRevokeResponse response type. */
export interface MsgRevokeResponse {
}
export declare const MsgGrant: {
    typeUrl: string;
    encode(message: MsgGrant, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgGrant;
    fromJSON(object: any): MsgGrant;
    toJSON(message: MsgGrant): unknown;
    fromPartial<I extends {
        granter?: string | undefined;
        grantee?: string | undefined;
        grant?: {
            authorization?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } | undefined;
            expiration?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
        } | undefined;
    } & {
        granter?: string | undefined;
        grantee?: string | undefined;
        grant?: ({
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
            } & Record<Exclude<keyof I["grant"]["authorization"], keyof Any>, never>) | undefined;
            expiration?: ({
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & Record<Exclude<keyof I["grant"]["expiration"], keyof import("../../../google/protobuf/timestamp").Timestamp>, never>) | undefined;
        } & Record<Exclude<keyof I["grant"], keyof Grant>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof MsgGrant>, never>>(object: I): MsgGrant;
};
export declare const MsgExecResponse: {
    typeUrl: string;
    encode(message: MsgExecResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgExecResponse;
    fromJSON(object: any): MsgExecResponse;
    toJSON(message: MsgExecResponse): unknown;
    fromPartial<I extends {
        results?: Uint8Array[] | undefined;
    } & {
        results?: (Uint8Array[] & Uint8Array[] & Record<Exclude<keyof I["results"], keyof Uint8Array[]>, never>) | undefined;
    } & Record<Exclude<keyof I, "results">, never>>(object: I): MsgExecResponse;
};
export declare const MsgExec: {
    typeUrl: string;
    encode(message: MsgExec, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgExec;
    fromJSON(object: any): MsgExec;
    toJSON(message: MsgExec): unknown;
    fromPartial<I extends {
        grantee?: string | undefined;
        msgs?: {
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        }[] | undefined;
    } & {
        grantee?: string | undefined;
        msgs?: ({
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        }[] & ({
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } & {
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } & Record<Exclude<keyof I["msgs"][number], keyof Any>, never>)[] & Record<Exclude<keyof I["msgs"], keyof {
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        }[]>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof MsgExec>, never>>(object: I): MsgExec;
};
export declare const MsgGrantResponse: {
    typeUrl: string;
    encode(_: MsgGrantResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgGrantResponse;
    fromJSON(_: any): MsgGrantResponse;
    toJSON(_: MsgGrantResponse): unknown;
    fromPartial<I extends {} & {} & Record<Exclude<keyof I, never>, never>>(_: I): MsgGrantResponse;
};
export declare const MsgRevoke: {
    typeUrl: string;
    encode(message: MsgRevoke, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgRevoke;
    fromJSON(object: any): MsgRevoke;
    toJSON(message: MsgRevoke): unknown;
    fromPartial<I extends {
        granter?: string | undefined;
        grantee?: string | undefined;
        msgTypeUrl?: string | undefined;
    } & {
        granter?: string | undefined;
        grantee?: string | undefined;
        msgTypeUrl?: string | undefined;
    } & Record<Exclude<keyof I, keyof MsgRevoke>, never>>(object: I): MsgRevoke;
};
export declare const MsgRevokeResponse: {
    typeUrl: string;
    encode(_: MsgRevokeResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgRevokeResponse;
    fromJSON(_: any): MsgRevokeResponse;
    toJSON(_: MsgRevokeResponse): unknown;
    fromPartial<I extends {} & {} & Record<Exclude<keyof I, never>, never>>(_: I): MsgRevokeResponse;
};
/** Msg defines the authz Msg service. */
export interface Msg {
    /**
     * Grant grants the provided authorization to the grantee on the granter's
     * account with the provided expiration time. If there is already a grant
     * for the given (granter, grantee, Authorization) triple, then the grant
     * will be overwritten.
     */
    Grant(request: MsgGrant): Promise<MsgGrantResponse>;
    /**
     * Exec attempts to execute the provided messages using
     * authorizations granted to the grantee. Each message should have only
     * one signer corresponding to the granter of the authorization.
     */
    Exec(request: MsgExec): Promise<MsgExecResponse>;
    /**
     * Revoke revokes any authorization corresponding to the provided method name on the
     * granter's account that has been granted to the grantee.
     */
    Revoke(request: MsgRevoke): Promise<MsgRevokeResponse>;
}
export declare class MsgClientImpl implements Msg {
    private readonly rpc;
    constructor(rpc: Rpc);
    Grant(request: MsgGrant): Promise<MsgGrantResponse>;
    Exec(request: MsgExec): Promise<MsgExecResponse>;
    Revoke(request: MsgRevoke): Promise<MsgRevokeResponse>;
}
