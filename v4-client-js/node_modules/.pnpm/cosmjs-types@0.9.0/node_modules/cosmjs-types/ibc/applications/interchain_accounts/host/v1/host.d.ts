import { BinaryReader, BinaryWriter } from "../../../../../binary";
export declare const protobufPackage = "ibc.applications.interchain_accounts.host.v1";
/**
 * Params defines the set of on-chain interchain accounts parameters.
 * The following parameters may be used to disable the host submodule.
 */
export interface Params {
    /** host_enabled enables or disables the host submodule. */
    hostEnabled: boolean;
    /** allow_messages defines a list of sdk message typeURLs allowed to be executed on a host chain. */
    allowMessages: string[];
}
export declare const Params: {
    typeUrl: string;
    encode(message: Params, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): Params;
    fromJSON(object: any): Params;
    toJSON(message: Params): unknown;
    fromPartial<I extends {
        hostEnabled?: boolean | undefined;
        allowMessages?: string[] | undefined;
    } & {
        hostEnabled?: boolean | undefined;
        allowMessages?: (string[] & string[] & Record<Exclude<keyof I["allowMessages"], keyof string[]>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof Params>, never>>(object: I): Params;
};
