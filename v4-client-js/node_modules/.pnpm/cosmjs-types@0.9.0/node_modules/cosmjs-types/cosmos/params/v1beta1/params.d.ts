import { BinaryReader, BinaryWriter } from "../../../binary";
export declare const protobufPackage = "cosmos.params.v1beta1";
/** ParameterChangeProposal defines a proposal to change one or more parameters. */
export interface ParameterChangeProposal {
    title: string;
    description: string;
    changes: ParamChange[];
}
/**
 * ParamChange defines an individual parameter change, for use in
 * ParameterChangeProposal.
 */
export interface ParamChange {
    subspace: string;
    key: string;
    value: string;
}
export declare const ParameterChangeProposal: {
    typeUrl: string;
    encode(message: ParameterChangeProposal, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): ParameterChangeProposal;
    fromJSON(object: any): ParameterChangeProposal;
    toJSON(message: ParameterChangeProposal): unknown;
    fromPartial<I extends {
        title?: string | undefined;
        description?: string | undefined;
        changes?: {
            subspace?: string | undefined;
            key?: string | undefined;
            value?: string | undefined;
        }[] | undefined;
    } & {
        title?: string | undefined;
        description?: string | undefined;
        changes?: ({
            subspace?: string | undefined;
            key?: string | undefined;
            value?: string | undefined;
        }[] & ({
            subspace?: string | undefined;
            key?: string | undefined;
            value?: string | undefined;
        } & {
            subspace?: string | undefined;
            key?: string | undefined;
            value?: string | undefined;
        } & Record<Exclude<keyof I["changes"][number], keyof ParamChange>, never>)[] & Record<Exclude<keyof I["changes"], keyof {
            subspace?: string | undefined;
            key?: string | undefined;
            value?: string | undefined;
        }[]>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof ParameterChangeProposal>, never>>(object: I): ParameterChangeProposal;
};
export declare const ParamChange: {
    typeUrl: string;
    encode(message: ParamChange, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): ParamChange;
    fromJSON(object: any): ParamChange;
    toJSON(message: ParamChange): unknown;
    fromPartial<I extends {
        subspace?: string | undefined;
        key?: string | undefined;
        value?: string | undefined;
    } & {
        subspace?: string | undefined;
        key?: string | undefined;
        value?: string | undefined;
    } & Record<Exclude<keyof I, keyof ParamChange>, never>>(object: I): ParamChange;
};
