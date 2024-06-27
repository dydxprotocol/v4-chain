import { Grant } from "./feegrant";
import { BinaryReader, BinaryWriter } from "../../../binary";
export declare const protobufPackage = "cosmos.feegrant.v1beta1";
/** GenesisState contains a set of fee allowances, persisted from the store */
export interface GenesisState {
    allowances: Grant[];
}
export declare const GenesisState: {
    typeUrl: string;
    encode(message: GenesisState, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): GenesisState;
    fromJSON(object: any): GenesisState;
    toJSON(message: GenesisState): unknown;
    fromPartial<I extends {
        allowances?: {
            granter?: string | undefined;
            grantee?: string | undefined;
            allowance?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } | undefined;
        }[] | undefined;
    } & {
        allowances?: ({
            granter?: string | undefined;
            grantee?: string | undefined;
            allowance?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } | undefined;
        }[] & ({
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
            } & Record<Exclude<keyof I["allowances"][number]["allowance"], keyof import("../../../google/protobuf/any").Any>, never>) | undefined;
        } & Record<Exclude<keyof I["allowances"][number], keyof Grant>, never>)[] & Record<Exclude<keyof I["allowances"], keyof {
            granter?: string | undefined;
            grantee?: string | undefined;
            allowance?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } | undefined;
        }[]>, never>) | undefined;
    } & Record<Exclude<keyof I, "allowances">, never>>(object: I): GenesisState;
};
