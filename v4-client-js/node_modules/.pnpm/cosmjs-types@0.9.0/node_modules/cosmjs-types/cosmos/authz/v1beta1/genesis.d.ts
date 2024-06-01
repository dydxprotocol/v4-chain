import { GrantAuthorization } from "./authz";
import { BinaryReader, BinaryWriter } from "../../../binary";
export declare const protobufPackage = "cosmos.authz.v1beta1";
/** GenesisState defines the authz module's genesis state. */
export interface GenesisState {
    authorization: GrantAuthorization[];
}
export declare const GenesisState: {
    typeUrl: string;
    encode(message: GenesisState, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): GenesisState;
    fromJSON(object: any): GenesisState;
    toJSON(message: GenesisState): unknown;
    fromPartial<I extends {
        authorization?: {
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
        }[] | undefined;
    } & {
        authorization?: ({
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
        }[] & ({
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
            } & Record<Exclude<keyof I["authorization"][number]["authorization"], keyof import("../../../google/protobuf/any").Any>, never>) | undefined;
            expiration?: ({
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & Record<Exclude<keyof I["authorization"][number]["expiration"], keyof import("../../../google/protobuf/timestamp").Timestamp>, never>) | undefined;
        } & Record<Exclude<keyof I["authorization"][number], keyof GrantAuthorization>, never>)[] & Record<Exclude<keyof I["authorization"], keyof {
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
        }[]>, never>) | undefined;
    } & Record<Exclude<keyof I, "authorization">, never>>(object: I): GenesisState;
};
