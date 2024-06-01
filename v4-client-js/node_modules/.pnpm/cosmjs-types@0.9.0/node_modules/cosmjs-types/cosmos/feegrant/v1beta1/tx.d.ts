import { Any } from "../../../google/protobuf/any";
import { BinaryReader, BinaryWriter } from "../../../binary";
import { Rpc } from "../../../helpers";
export declare const protobufPackage = "cosmos.feegrant.v1beta1";
/**
 * MsgGrantAllowance adds permission for Grantee to spend up to Allowance
 * of fees from the account of Granter.
 */
export interface MsgGrantAllowance {
    /** granter is the address of the user granting an allowance of their funds. */
    granter: string;
    /** grantee is the address of the user being granted an allowance of another user's funds. */
    grantee: string;
    /** allowance can be any of basic, periodic, allowed fee allowance. */
    allowance?: Any;
}
/** MsgGrantAllowanceResponse defines the Msg/GrantAllowanceResponse response type. */
export interface MsgGrantAllowanceResponse {
}
/** MsgRevokeAllowance removes any existing Allowance from Granter to Grantee. */
export interface MsgRevokeAllowance {
    /** granter is the address of the user granting an allowance of their funds. */
    granter: string;
    /** grantee is the address of the user being granted an allowance of another user's funds. */
    grantee: string;
}
/** MsgRevokeAllowanceResponse defines the Msg/RevokeAllowanceResponse response type. */
export interface MsgRevokeAllowanceResponse {
}
export declare const MsgGrantAllowance: {
    typeUrl: string;
    encode(message: MsgGrantAllowance, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgGrantAllowance;
    fromJSON(object: any): MsgGrantAllowance;
    toJSON(message: MsgGrantAllowance): unknown;
    fromPartial<I extends {
        granter?: string | undefined;
        grantee?: string | undefined;
        allowance?: {
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } | undefined;
    } & {
        granter?: string | undefined;
        grantee?: string | undefined;
        allowance?: ({
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } & {
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } & Record<Exclude<keyof I["allowance"], keyof Any>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof MsgGrantAllowance>, never>>(object: I): MsgGrantAllowance;
};
export declare const MsgGrantAllowanceResponse: {
    typeUrl: string;
    encode(_: MsgGrantAllowanceResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgGrantAllowanceResponse;
    fromJSON(_: any): MsgGrantAllowanceResponse;
    toJSON(_: MsgGrantAllowanceResponse): unknown;
    fromPartial<I extends {} & {} & Record<Exclude<keyof I, never>, never>>(_: I): MsgGrantAllowanceResponse;
};
export declare const MsgRevokeAllowance: {
    typeUrl: string;
    encode(message: MsgRevokeAllowance, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgRevokeAllowance;
    fromJSON(object: any): MsgRevokeAllowance;
    toJSON(message: MsgRevokeAllowance): unknown;
    fromPartial<I extends {
        granter?: string | undefined;
        grantee?: string | undefined;
    } & {
        granter?: string | undefined;
        grantee?: string | undefined;
    } & Record<Exclude<keyof I, keyof MsgRevokeAllowance>, never>>(object: I): MsgRevokeAllowance;
};
export declare const MsgRevokeAllowanceResponse: {
    typeUrl: string;
    encode(_: MsgRevokeAllowanceResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgRevokeAllowanceResponse;
    fromJSON(_: any): MsgRevokeAllowanceResponse;
    toJSON(_: MsgRevokeAllowanceResponse): unknown;
    fromPartial<I extends {} & {} & Record<Exclude<keyof I, never>, never>>(_: I): MsgRevokeAllowanceResponse;
};
/** Msg defines the feegrant msg service. */
export interface Msg {
    /**
     * GrantAllowance grants fee allowance to the grantee on the granter's
     * account with the provided expiration time.
     */
    GrantAllowance(request: MsgGrantAllowance): Promise<MsgGrantAllowanceResponse>;
    /**
     * RevokeAllowance revokes any fee allowance of granter's account that
     * has been granted to the grantee.
     */
    RevokeAllowance(request: MsgRevokeAllowance): Promise<MsgRevokeAllowanceResponse>;
}
export declare class MsgClientImpl implements Msg {
    private readonly rpc;
    constructor(rpc: Rpc);
    GrantAllowance(request: MsgGrantAllowance): Promise<MsgGrantAllowanceResponse>;
    RevokeAllowance(request: MsgRevokeAllowance): Promise<MsgRevokeAllowanceResponse>;
}
