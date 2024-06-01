import { Plan } from "./upgrade";
import { BinaryReader, BinaryWriter } from "../../../binary";
import { Rpc } from "../../../helpers";
export declare const protobufPackage = "cosmos.upgrade.v1beta1";
/**
 * MsgSoftwareUpgrade is the Msg/SoftwareUpgrade request type.
 *
 * Since: cosmos-sdk 0.46
 */
export interface MsgSoftwareUpgrade {
    /** authority is the address that controls the module (defaults to x/gov unless overwritten). */
    authority: string;
    /** plan is the upgrade plan. */
    plan: Plan;
}
/**
 * MsgSoftwareUpgradeResponse is the Msg/SoftwareUpgrade response type.
 *
 * Since: cosmos-sdk 0.46
 */
export interface MsgSoftwareUpgradeResponse {
}
/**
 * MsgCancelUpgrade is the Msg/CancelUpgrade request type.
 *
 * Since: cosmos-sdk 0.46
 */
export interface MsgCancelUpgrade {
    /** authority is the address that controls the module (defaults to x/gov unless overwritten). */
    authority: string;
}
/**
 * MsgCancelUpgradeResponse is the Msg/CancelUpgrade response type.
 *
 * Since: cosmos-sdk 0.46
 */
export interface MsgCancelUpgradeResponse {
}
export declare const MsgSoftwareUpgrade: {
    typeUrl: string;
    encode(message: MsgSoftwareUpgrade, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgSoftwareUpgrade;
    fromJSON(object: any): MsgSoftwareUpgrade;
    toJSON(message: MsgSoftwareUpgrade): unknown;
    fromPartial<I extends {
        authority?: string | undefined;
        plan?: {
            name?: string | undefined;
            time?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            height?: bigint | undefined;
            info?: string | undefined;
            upgradedClientState?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } | undefined;
        } | undefined;
    } & {
        authority?: string | undefined;
        plan?: ({
            name?: string | undefined;
            time?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            height?: bigint | undefined;
            info?: string | undefined;
            upgradedClientState?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } | undefined;
        } & {
            name?: string | undefined;
            time?: ({
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & Record<Exclude<keyof I["plan"]["time"], keyof import("../../../google/protobuf/timestamp").Timestamp>, never>) | undefined;
            height?: bigint | undefined;
            info?: string | undefined;
            upgradedClientState?: ({
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } & {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } & Record<Exclude<keyof I["plan"]["upgradedClientState"], keyof import("../../../google/protobuf/any").Any>, never>) | undefined;
        } & Record<Exclude<keyof I["plan"], keyof Plan>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof MsgSoftwareUpgrade>, never>>(object: I): MsgSoftwareUpgrade;
};
export declare const MsgSoftwareUpgradeResponse: {
    typeUrl: string;
    encode(_: MsgSoftwareUpgradeResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgSoftwareUpgradeResponse;
    fromJSON(_: any): MsgSoftwareUpgradeResponse;
    toJSON(_: MsgSoftwareUpgradeResponse): unknown;
    fromPartial<I extends {} & {} & Record<Exclude<keyof I, never>, never>>(_: I): MsgSoftwareUpgradeResponse;
};
export declare const MsgCancelUpgrade: {
    typeUrl: string;
    encode(message: MsgCancelUpgrade, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgCancelUpgrade;
    fromJSON(object: any): MsgCancelUpgrade;
    toJSON(message: MsgCancelUpgrade): unknown;
    fromPartial<I extends {
        authority?: string | undefined;
    } & {
        authority?: string | undefined;
    } & Record<Exclude<keyof I, "authority">, never>>(object: I): MsgCancelUpgrade;
};
export declare const MsgCancelUpgradeResponse: {
    typeUrl: string;
    encode(_: MsgCancelUpgradeResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgCancelUpgradeResponse;
    fromJSON(_: any): MsgCancelUpgradeResponse;
    toJSON(_: MsgCancelUpgradeResponse): unknown;
    fromPartial<I extends {} & {} & Record<Exclude<keyof I, never>, never>>(_: I): MsgCancelUpgradeResponse;
};
/** Msg defines the upgrade Msg service. */
export interface Msg {
    /**
     * SoftwareUpgrade is a governance operation for initiating a software upgrade.
     *
     * Since: cosmos-sdk 0.46
     */
    SoftwareUpgrade(request: MsgSoftwareUpgrade): Promise<MsgSoftwareUpgradeResponse>;
    /**
     * CancelUpgrade is a governance operation for cancelling a previously
     * approved software upgrade.
     *
     * Since: cosmos-sdk 0.46
     */
    CancelUpgrade(request: MsgCancelUpgrade): Promise<MsgCancelUpgradeResponse>;
}
export declare class MsgClientImpl implements Msg {
    private readonly rpc;
    constructor(rpc: Rpc);
    SoftwareUpgrade(request: MsgSoftwareUpgrade): Promise<MsgSoftwareUpgradeResponse>;
    CancelUpgrade(request: MsgCancelUpgrade): Promise<MsgCancelUpgradeResponse>;
}
