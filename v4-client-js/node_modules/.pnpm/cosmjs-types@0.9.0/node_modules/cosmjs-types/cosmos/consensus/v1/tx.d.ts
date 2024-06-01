import { BlockParams, EvidenceParams, ValidatorParams } from "../../../tendermint/types/params";
import { BinaryReader, BinaryWriter } from "../../../binary";
import { Rpc } from "../../../helpers";
export declare const protobufPackage = "cosmos.consensus.v1";
/** MsgUpdateParams is the Msg/UpdateParams request type. */
export interface MsgUpdateParams {
    /** authority is the address that controls the module (defaults to x/gov unless overwritten). */
    authority: string;
    /**
     * params defines the x/consensus parameters to update.
     * VersionsParams is not included in this Msg because it is tracked
     * separarately in x/upgrade.
     *
     * NOTE: All parameters must be supplied.
     */
    block?: BlockParams;
    evidence?: EvidenceParams;
    validator?: ValidatorParams;
}
/**
 * MsgUpdateParamsResponse defines the response structure for executing a
 * MsgUpdateParams message.
 */
export interface MsgUpdateParamsResponse {
}
export declare const MsgUpdateParams: {
    typeUrl: string;
    encode(message: MsgUpdateParams, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgUpdateParams;
    fromJSON(object: any): MsgUpdateParams;
    toJSON(message: MsgUpdateParams): unknown;
    fromPartial<I extends {
        authority?: string | undefined;
        block?: {
            maxBytes?: bigint | undefined;
            maxGas?: bigint | undefined;
        } | undefined;
        evidence?: {
            maxAgeNumBlocks?: bigint | undefined;
            maxAgeDuration?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            maxBytes?: bigint | undefined;
        } | undefined;
        validator?: {
            pubKeyTypes?: string[] | undefined;
        } | undefined;
    } & {
        authority?: string | undefined;
        block?: ({
            maxBytes?: bigint | undefined;
            maxGas?: bigint | undefined;
        } & {
            maxBytes?: bigint | undefined;
            maxGas?: bigint | undefined;
        } & Record<Exclude<keyof I["block"], keyof BlockParams>, never>) | undefined;
        evidence?: ({
            maxAgeNumBlocks?: bigint | undefined;
            maxAgeDuration?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            maxBytes?: bigint | undefined;
        } & {
            maxAgeNumBlocks?: bigint | undefined;
            maxAgeDuration?: ({
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & Record<Exclude<keyof I["evidence"]["maxAgeDuration"], keyof import("../../../google/protobuf/duration").Duration>, never>) | undefined;
            maxBytes?: bigint | undefined;
        } & Record<Exclude<keyof I["evidence"], keyof EvidenceParams>, never>) | undefined;
        validator?: ({
            pubKeyTypes?: string[] | undefined;
        } & {
            pubKeyTypes?: (string[] & string[] & Record<Exclude<keyof I["validator"]["pubKeyTypes"], keyof string[]>, never>) | undefined;
        } & Record<Exclude<keyof I["validator"], "pubKeyTypes">, never>) | undefined;
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
    /**
     * UpdateParams defines a governance operation for updating the x/consensus_param module parameters.
     * The authority is defined in the keeper.
     *
     * Since: cosmos-sdk 0.47
     */
    UpdateParams(request: MsgUpdateParams): Promise<MsgUpdateParamsResponse>;
}
export declare class MsgClientImpl implements Msg {
    private readonly rpc;
    constructor(rpc: Rpc);
    UpdateParams(request: MsgUpdateParams): Promise<MsgUpdateParamsResponse>;
}
