import { Coin } from "../../base/v1beta1/coin";
import { BinaryReader, BinaryWriter } from "../../../binary";
import { Rpc } from "../../../helpers";
export declare const protobufPackage = "cosmos.crisis.v1beta1";
/** MsgVerifyInvariant represents a message to verify a particular invariance. */
export interface MsgVerifyInvariant {
    /** sender is the account address of private key to send coins to fee collector account. */
    sender: string;
    /** name of the invariant module. */
    invariantModuleName: string;
    /** invariant_route is the msg's invariant route. */
    invariantRoute: string;
}
/** MsgVerifyInvariantResponse defines the Msg/VerifyInvariant response type. */
export interface MsgVerifyInvariantResponse {
}
/**
 * MsgUpdateParams is the Msg/UpdateParams request type.
 *
 * Since: cosmos-sdk 0.47
 */
export interface MsgUpdateParams {
    /** authority is the address that controls the module (defaults to x/gov unless overwritten). */
    authority: string;
    /** constant_fee defines the x/crisis parameter. */
    constantFee: Coin;
}
/**
 * MsgUpdateParamsResponse defines the response structure for executing a
 * MsgUpdateParams message.
 *
 * Since: cosmos-sdk 0.47
 */
export interface MsgUpdateParamsResponse {
}
export declare const MsgVerifyInvariant: {
    typeUrl: string;
    encode(message: MsgVerifyInvariant, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgVerifyInvariant;
    fromJSON(object: any): MsgVerifyInvariant;
    toJSON(message: MsgVerifyInvariant): unknown;
    fromPartial<I extends {
        sender?: string | undefined;
        invariantModuleName?: string | undefined;
        invariantRoute?: string | undefined;
    } & {
        sender?: string | undefined;
        invariantModuleName?: string | undefined;
        invariantRoute?: string | undefined;
    } & Record<Exclude<keyof I, keyof MsgVerifyInvariant>, never>>(object: I): MsgVerifyInvariant;
};
export declare const MsgVerifyInvariantResponse: {
    typeUrl: string;
    encode(_: MsgVerifyInvariantResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgVerifyInvariantResponse;
    fromJSON(_: any): MsgVerifyInvariantResponse;
    toJSON(_: MsgVerifyInvariantResponse): unknown;
    fromPartial<I extends {} & {} & Record<Exclude<keyof I, never>, never>>(_: I): MsgVerifyInvariantResponse;
};
export declare const MsgUpdateParams: {
    typeUrl: string;
    encode(message: MsgUpdateParams, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgUpdateParams;
    fromJSON(object: any): MsgUpdateParams;
    toJSON(message: MsgUpdateParams): unknown;
    fromPartial<I extends {
        authority?: string | undefined;
        constantFee?: {
            denom?: string | undefined;
            amount?: string | undefined;
        } | undefined;
    } & {
        authority?: string | undefined;
        constantFee?: ({
            denom?: string | undefined;
            amount?: string | undefined;
        } & {
            denom?: string | undefined;
            amount?: string | undefined;
        } & Record<Exclude<keyof I["constantFee"], keyof Coin>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof MsgUpdateParams>, never>>(object: I): MsgUpdateParams;
};
export declare const MsgUpdateParamsResponse: {
    typeUrl: string;
    encode(_: MsgUpdateParamsResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgUpdateParamsResponse;
    fromJSON(_: any): MsgUpdateParamsResponse;
    toJSON(_: MsgUpdateParamsResponse): unknown;
    fromPartial<I extends {} & {} & Record<Exclude<keyof I, never>, never>>(_: I): MsgUpdateParamsResponse;
};
/** Msg defines the bank Msg service. */
export interface Msg {
    /** VerifyInvariant defines a method to verify a particular invariant. */
    VerifyInvariant(request: MsgVerifyInvariant): Promise<MsgVerifyInvariantResponse>;
    /**
     * UpdateParams defines a governance operation for updating the x/crisis module
     * parameters. The authority is defined in the keeper.
     *
     * Since: cosmos-sdk 0.47
     */
    UpdateParams(request: MsgUpdateParams): Promise<MsgUpdateParamsResponse>;
}
export declare class MsgClientImpl implements Msg {
    private readonly rpc;
    constructor(rpc: Rpc);
    VerifyInvariant(request: MsgVerifyInvariant): Promise<MsgVerifyInvariantResponse>;
    UpdateParams(request: MsgUpdateParams): Promise<MsgUpdateParamsResponse>;
}
